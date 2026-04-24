// Package validator extracts, diffs, and validates OpenAPI JSON schemas from Kubernetes CRD YAML files.
package validator

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// SchemaInfo holds the extracted schema for a single CRD version.
type SchemaInfo struct {
	Group       string
	Kind        string
	Version     string
	ResourceKey string
	Schema      map[string]interface{}
	SourceFile  string
}

// crdDoc is the minimal structure we need from a CRD YAML document.
type crdDoc struct {
	Kind string `yaml:"kind"`
	Spec struct {
		Group string `yaml:"group"`
		Names struct {
			Kind string `yaml:"kind"`
		} `yaml:"names"`
		Versions []struct {
			Name   string `yaml:"name"`
			Schema struct {
				OpenAPIV3Schema map[string]interface{} `yaml:"openAPIV3Schema"`
			} `yaml:"schema"`
		} `yaml:"versions"`
	} `yaml:"spec"`
}

// looksLikeTemplate returns true if the file content contains Go template
// directives ({{ ... }}), which means it's a Helm template, not raw YAML.
func looksLikeTemplate(data []byte) bool {
	return bytes.Contains(data, []byte("{{"))
}

// ExtractSchemasFromCRD parses a CRD YAML file (possibly multi-doc) and returns
// the openAPIV3Schema for each version found.
func ExtractSchemasFromCRD(path string) ([]SchemaInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Skip Helm template files: they contain Go template directives that
	// aren't valid YAML and can never be CRD definitions.
	if looksLikeTemplate(data) {
		return nil, nil
	}

	var results []SchemaInfo
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	for {
		var doc crdDoc
		err := decoder.Decode(&doc)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("warning: failed to decode document in %s: %v", path, err)
			break
		}
		if doc.Kind != "CustomResourceDefinition" {
			continue
		}

		group := doc.Spec.Group
		kind := doc.Spec.Names.Kind
		kindLower := strings.ToLower(kind)

		for _, ver := range doc.Spec.Versions {
			if ver.Schema.OpenAPIV3Schema == nil {
				continue
			}
			results = append(results, SchemaInfo{
				Group:       group,
				Kind:        kind,
				Version:     ver.Name,
				ResourceKey: kindLower + "." + ver.Name,
				Schema:      ver.Schema.OpenAPIV3Schema,
				SourceFile:  path,
			})
		}
	}

	return results, nil
}

// ExtractSchemasFromDir searches common CRD directories under the given root,
// falling back to a recursive scan if nothing is found in the well-known paths.
func ExtractSchemasFromDir(dir string) ([]SchemaInfo, error) {
	searchDirs := []string{
		"config/crd/bases",
		"config/crd",
		"deploy/crds",
		"charts",
		"manifests",
	}

	var results []SchemaInfo
	searched := map[string]bool{}
	seen := map[string]bool{}

	for _, sub := range searchDirs {
		d := filepath.Join(dir, sub)
		info, err := os.Stat(d)
		if err != nil || !info.IsDir() {
			continue
		}
		searched[d] = true
		found, err := scanDirDedup(d, seen)
		if err != nil {
			return nil, err
		}
		results = append(results, found...)
	}

	// Fallback: recursive scan of the whole directory, skipping already-searched paths.
	// Capped at 500 YAML files to avoid spending 30+ minutes on large repos
	// (e.g. kserve, kuberay, opendatahub-operator) that have thousands of
	// test fixtures, Helm charts, and vendored manifests.
	if len(results) == 0 {
		yamlCount := 0
		const maxYAMLFiles = 500
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				base := filepath.Base(path)
				// Skip directories that never contain CRD definitions.
				switch base {
				case ".git", "vendor", "node_modules", "testdata", "test", "tests",
					"_output", "hack", "docs", "examples", "samples", ".github":
					return filepath.SkipDir
				}
				return nil
			}
			ext := filepath.Ext(path)
			if ext != ".yaml" && ext != ".yml" {
				return nil
			}
			if seen[path] {
				return nil
			}
			// Skip files under already-searched directories.
			for s := range searched {
				rel, relErr := filepath.Rel(s, path)
				if relErr == nil && !strings.HasPrefix(rel, "..") {
					return nil
				}
			}
			yamlCount++
			if yamlCount > maxYAMLFiles {
				return filepath.SkipAll
			}
			schemas, extractErr := ExtractSchemasFromCRD(path)
			if extractErr != nil {
				return nil
			}
			results = append(results, schemas...)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

// scanDirDedup recursively finds .yaml and .yml files in d and extracts CRD schemas,
// tracking already-processed files in the seen map to avoid duplicates.
func scanDirDedup(d string, seen map[string]bool) ([]SchemaInfo, error) {
	var results []SchemaInfo
	err := filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			// Skip Helm template directories: rendered templates, not raw CRDs.
			if filepath.Base(path) == "templates" {
				return filepath.SkipDir
			}
			return nil
		}
		ext := filepath.Ext(path)
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}
		if seen[path] {
			return nil
		}
		seen[path] = true
		schemas, extractErr := ExtractSchemasFromCRD(path)
		if extractErr != nil {
			return nil
		}
		results = append(results, schemas...)
		return nil
	})
	return results, err
}
