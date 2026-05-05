package extractor

import (
	"sort"
	"strings"
)

var resourceDefaultPatterns = []string{
	"config/configmap/*.yaml",
	"config/configmap/*.yml",
	"config/configmap/**/*.yaml",
	"config/default/*.yaml",
	"config/manager/*.yaml",
	"config/manager/*.yml",
	"manifests/**/configmap*.yaml",
}

// resourceDefaultKeys lists configmap data keys that contain resource defaults
// worth surfacing (inference config, deployment defaults, scaling config).
var resourceDefaultKeys = map[string]bool{
	"inferenceservice":  true,
	"deploy":            true,
	"defaults":          true,
	"resources":         true,
	"autoscaler":        true,
	"ingress":           true,
	"logger":            true,
	"router":            true,
	"agent":             true,
	"metricsAggregator": true,
	"batcher":           true,
	"storageInitializer": true,
	"explainers":        true,
}

// extractResourceDefaults scans configmaps for inference/deployment resource
// defaults. These are typically in config/configmap/inferenceservice.yaml and
// similar files used by KServe, KNative, and other serving frameworks.
func extractResourceDefaults(repoPath string) []ResourceDefault {
	files := findYAMLFiles(repoPath, resourceDefaultPatterns)
	var defaults []ResourceDefault

	for _, fpath := range files {
		rel := relativePath(repoPath, fpath)
		for _, doc := range parseYAMLSafe(fpath) {
			kind, _ := doc["kind"].(string)
			if kind != "ConfigMap" {
				continue
			}

			meta, _ := doc["metadata"].(map[string]interface{})
			name := ""
			if meta != nil {
				name, _ = meta["name"].(string)
			}
			if name == "" {
				continue
			}

			data, _ := doc["data"].(map[string]interface{})
			if data == nil {
				continue
			}

			for key, val := range data {
				if !isResourceDefaultKey(key) {
					continue
				}

				rd := ResourceDefault{
					Component: name,
					Key:       key,
					Source:    rel,
				}

				// Try to parse the value as YAML (configmap data values are
				// often embedded YAML strings containing resource specs)
				if strVal, ok := val.(string); ok {
					parsed := parseEmbeddedYAML(strVal)
					if parsed != nil {
						rd.Values = parsed
					}
				} else if mapVal, ok := val.(map[string]interface{}); ok {
					rd.Values = mapVal
				}

				if rd.Values != nil {
					defaults = append(defaults, rd)
				}
			}
		}
	}

	sort.Slice(defaults, func(i, j int) bool {
		if defaults[i].Component != defaults[j].Component {
			return defaults[i].Component < defaults[j].Component
		}
		return defaults[i].Key < defaults[j].Key
	})
	return defaults
}

// isResourceDefaultKey checks whether a configmap data key is likely to contain
// resource defaults worth extracting.
func isResourceDefaultKey(key string) bool {
	lower := strings.ToLower(key)
	for k := range resourceDefaultKeys {
		if strings.EqualFold(key, k) || strings.Contains(lower, strings.ToLower(k)) {
			return true
		}
	}
	return false
}

// parseEmbeddedYAML attempts to parse a string value as YAML. ConfigMap data
// values often contain embedded YAML with resource configurations.
func parseEmbeddedYAML(s string) map[string]interface{} {
	s = strings.TrimSpace(s)
	if s == "" || s == "{}" || s == "null" {
		return nil
	}

	// parseYAMLSafe works on files; we need to parse a string.
	// Use the same yaml decoder approach but from a string.
	docs := parseYAMLFromBytes([]byte(s))
	if len(docs) > 0 {
		return docs[0]
	}
	return nil
}
