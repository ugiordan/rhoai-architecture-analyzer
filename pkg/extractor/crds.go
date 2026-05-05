package extractor

import (
	"fmt"
	"strings"
)

var crdSearchPatterns = []string{
	"config/crd/bases/*.yaml",
	"config/crd/bases/*.yml",
	"config/crd/*.yaml",
	"config/crd/*.yml",
	"config/crd/**/*.yaml",
	"config/crd/**/*.yml",
	"deploy/crds/*.yaml",
	"deploy/crds/*.yml",
	"charts/**/crds/*.yaml",
	"charts/**/crds/*.yml",
	"charts/**/files/*.yaml",
	"manifests/**/crd*.yaml",
	"manifests/**/crd*.yml",
	"manifests/**/crds/*.yaml",
	"manifests/**/base/crds/*.yaml",
}

// crdSkipPaths filters out external and bundle CRD files that should not
// be included in the component's own CRD inventory.
// Directory-level exclusions (test, tests, testdata) are handled by
// DefaultExcludedDirs in yaml.go and no longer need to be listed here.
var crdSkipPaths = []string{
	"/external/",
	"_test",
	"-bundle/",
	"/opt/manifests/",
}

// crdVersionSelection holds the result of selecting the best version from a CRD spec.
type crdVersionSelection struct {
	versions       []CRDVersion
	storageVersion string
	bestFields     int
	bestRules      []string
}

// selectBestCRDVersion analyzes CRD versions and returns the storage version
// (or the version with the most fields if no storage version is marked).
func selectBestCRDVersion(rawVersions []interface{}) crdVersionSelection {
	var result crdVersionSelection
	for _, v := range rawVersions {
		ver, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		verName, _ := ver["name"].(string)
		served, _ := ver["served"].(bool)
		storage, _ := ver["storage"].(bool)

		result.versions = append(result.versions, CRDVersion{
			Name:    verName,
			Served:  served,
			Storage: storage,
		})

		var openAPI map[string]interface{}
		if schema, ok := ver["schema"].(map[string]interface{}); ok {
			openAPI, _ = schema["openAPIV3Schema"].(map[string]interface{})
		}
		fieldsCount := countFields(openAPI, 0)
		celRules := extractCELRules(openAPI, 0)

		// Prefer the storage version; among non-storage versions,
		// pick the one with the most fields.
		if storage || (result.storageVersion == "" && fieldsCount > result.bestFields) {
			result.storageVersion = verName
			result.bestFields = fieldsCount
			result.bestRules = celRules
		}
	}
	if result.storageVersion == "" && len(result.versions) > 0 {
		result.storageVersion = result.versions[0].Name
	}
	return result
}

// parseCRDFromDoc extracts a CRD struct from a parsed YAML document.
// Returns nil if the document is not a valid CRD.
func parseCRDFromDoc(doc map[string]interface{}, source string) *CRD {
	kind, _ := doc["kind"].(string)
	if kind != "CustomResourceDefinition" {
		return nil
	}
	spec, ok := doc["spec"].(map[string]interface{})
	if !ok {
		return nil
	}
	group, _ := spec["group"].(string)
	names, _ := spec["names"].(map[string]interface{})
	crdKind := ""
	if names != nil {
		crdKind, _ = names["kind"].(string)
	}
	if group == "" || crdKind == "" {
		return nil
	}
	scope, _ := spec["scope"].(string)
	rawVersions, _ := spec["versions"].([]interface{})

	vs := selectBestCRDVersion(rawVersions)

	return &CRD{
		Group:           group,
		Version:         vs.storageVersion,
		Kind:            crdKind,
		Scope:           scope,
		Versions:        vs.versions,
		FieldsCount:     vs.bestFields,
		ValidationRules: vs.bestRules,
		Source:          source,
	}
}

// extractCRDs scans YAML files for CustomResourceDefinition documents and
// returns a slice of CRD structs with schema statistics and CEL rules.
// Deduplicates by (group, kind), keeping the entry with the most fields
// (fullest schema) when the same CRD appears in multiple directories.
func extractCRDs(repoPath string) []CRD {
	files := findYAMLFiles(repoPath, crdSearchPatterns)

	seen := make(map[string]int) // "group/kind" -> index in crds slice
	var crds []CRD

	for _, fpath := range files {
		rel := relativePath(repoPath, fpath)
		if shouldSkipCRDPath(rel) {
			continue
		}
		for _, doc := range parseYAMLSafe(fpath) {
			crd := parseCRDFromDoc(doc, rel)
			if crd == nil {
				continue
			}
			gk := fmt.Sprintf("%s/%s", crd.Group, crd.Kind)
			if idx, exists := seen[gk]; exists {
				if crd.FieldsCount > crds[idx].FieldsCount {
					crds[idx] = *crd
				}
			} else {
				seen[gk] = len(crds)
				crds = append(crds, *crd)
			}
		}
	}

	if crds == nil {
		crds = []CRD{}
	}
	return crds
}

// countFields recursively counts fields in an OpenAPI v3 schema.
func countFields(schema map[string]interface{}, depth int) int {
	if schema == nil || depth > 50 {
		return 0
	}
	count := 0
	if props, ok := schema["properties"].(map[string]interface{}); ok && len(props) > 0 {
		count += len(props)
		for _, propSchema := range props {
			if ps, ok := propSchema.(map[string]interface{}); ok && len(ps) > 0 {
				count += countFields(ps, depth+1)
			}
		}
	}
	if items, ok := schema["items"].(map[string]interface{}); ok && len(items) > 0 {
		count += countFields(items, depth+1)
	}
	if additional, ok := schema["additionalProperties"].(map[string]interface{}); ok && len(additional) > 0 {
		count += countFields(additional, depth+1)
	}
	return count
}

// extractCELRules recursively extracts x-kubernetes-validations CEL rules.
func extractCELRules(schema map[string]interface{}, depth int) []string {
	if schema == nil || depth > 50 {
		return nil
	}
	var rules []string
	if validations, ok := schema["x-kubernetes-validations"].([]interface{}); ok {
		for _, v := range validations {
			if vm, ok := v.(map[string]interface{}); ok {
				if rule, ok := vm["rule"].(string); ok {
					rules = append(rules, rule)
				}
			}
		}
	}
	if props, ok := schema["properties"].(map[string]interface{}); ok && len(props) > 0 {
		for _, propSchema := range props {
			if ps, ok := propSchema.(map[string]interface{}); ok && len(ps) > 0 {
				rules = append(rules, extractCELRules(ps, depth+1)...)
			}
		}
	}
	if items, ok := schema["items"].(map[string]interface{}); ok && len(items) > 0 {
		rules = append(rules, extractCELRules(items, depth+1)...)
	}
	if rules == nil {
		rules = []string{}
	}
	return rules
}

func shouldSkipCRDPath(rel string) bool {
	for _, skip := range crdSkipPaths {
		if strings.Contains(rel, skip) {
			return true
		}
	}
	return false
}
