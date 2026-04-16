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

// crdSkipPaths filters out test, external, and bundle CRD files that should not
// be included in the component's own CRD inventory.
var crdSkipPaths = []string{
	"/test/",
	"/tests/",
	"/external/",
	"/testdata/",
	"_test",
	"-bundle/",
	"/opt/manifests/",
}

// extractCRDs scans YAML files for CustomResourceDefinition documents and
// returns a slice of CRD structs with schema statistics and CEL rules.
// Deduplicates by (group, version, kind), keeping the entry with the most fields
// (fullest schema) when the same CRD appears in multiple directories.
func extractCRDs(repoPath string) []CRD {
	files := findYAMLFiles(repoPath, crdSearchPatterns)

	// Map from "group/version/kind" to best CRD entry (highest field count)
	seen := make(map[string]int) // gvk -> index in crds slice
	var crds []CRD

	for _, fpath := range files {
		rel := relativePath(repoPath, fpath)
		if shouldSkipCRDPath(rel) {
			continue
		}
		for _, doc := range parseYAMLSafe(fpath) {
			kind, _ := doc["kind"].(string)
			if kind != "CustomResourceDefinition" {
				continue
			}
			spec, ok := doc["spec"].(map[string]interface{})
			if !ok {
				continue
			}
			group, _ := spec["group"].(string)
			names, _ := spec["names"].(map[string]interface{})
			crdKind := ""
			if names != nil {
				crdKind, _ = names["kind"].(string)
			}
			scope, _ := spec["scope"].(string)
			versions, _ := spec["versions"].([]interface{})

			for _, v := range versions {
				ver, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				verName, _ := ver["name"].(string)
				var openAPI map[string]interface{}
				if schema, ok := ver["schema"].(map[string]interface{}); ok {
					openAPI, _ = schema["openAPIV3Schema"].(map[string]interface{})
				}
				fieldsCount := countFields(openAPI, 0)
				celRules := extractCELRules(openAPI, 0)

				gvk := fmt.Sprintf("%s/%s/%s", group, verName, crdKind)
				crd := CRD{
					Group:           group,
					Version:         verName,
					Kind:            crdKind,
					Scope:           scope,
					FieldsCount:     fieldsCount,
					ValidationRules: celRules,
					Source:          relativePath(repoPath, fpath),
				}

				if idx, exists := seen[gvk]; exists {
					// Keep the entry with more fields (fullest schema)
					if fieldsCount > crds[idx].FieldsCount {
						crds[idx] = crd
					}
				} else {
					seen[gvk] = len(crds)
					crds = append(crds, crd)
				}
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
