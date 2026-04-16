package extractor

import "fmt"

var chartPatterns = []string{
	"**/Chart.yaml",
	"**/Chart.yml",
}

var valuesPatterns = []string{
	"**/values.yaml",
	"**/values.yml",
}

var securityKeys = map[string]bool{
	"securityContext":    true,
	"podSecurityContext": true,
	"tls":               true,
	"networkPolicy":     true,
	"serviceAccount":    true,
	"rbac":              true,
	"auth":              true,
}

// extractHelm extracts Helm chart metadata and security-relevant value defaults.
func extractHelm(repoPath string) *HelmData {
	chartFiles := findYAMLFiles(repoPath, chartPatterns)
	valuesFiles := findYAMLFiles(repoPath, valuesPatterns)

	chartName := ""
	chartVersion := ""

	// Parse Chart.yaml for metadata
	for _, fpath := range chartFiles {
		docs := parseYAMLSafe(fpath)
		for _, doc := range docs {
			if n, ok := doc["name"].(string); ok && n != "" {
				chartName = n
			}
			if v, ok := doc["version"].(string); ok && v != "" {
				chartVersion = v
			}
			break // Use first doc in first Chart.yaml
		}
		if chartName != "" {
			break
		}
	}

	// Parse values.yaml for security-relevant defaults
	valuesDefaults := make(map[string]interface{})
	for _, fpath := range valuesFiles {
		for _, doc := range parseYAMLSafe(fpath) {
			for key := range doc {
				if securityKeys[key] {
					flattened := flattenDict(doc[key], key)
					for k, v := range flattened {
						valuesDefaults[k] = v
					}
				}
			}
		}
	}

	if chartName == "" && len(valuesDefaults) == 0 {
		return &HelmData{}
	}

	return &HelmData{
		ChartName:      chartName,
		ChartVersion:   chartVersion,
		ValuesDefaults: valuesDefaults,
	}
}

// flattenDict flattens a nested map to dot-notation keys.
func flattenDict(v interface{}, prefix string) map[string]interface{} {
	return flattenDictDepth(v, prefix, 0)
}

func flattenDictDepth(v interface{}, prefix string, depth int) map[string]interface{} {
	if depth > 50 {
		if prefix != "" {
			return map[string]interface{}{prefix: v}
		}
		return map[string]interface{}{}
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		if prefix != "" {
			return map[string]interface{}{prefix: v}
		}
		return map[string]interface{}{}
	}
	items := make(map[string]interface{})
	for key, value := range m {
		newKey := key
		if prefix != "" {
			newKey = fmt.Sprintf("%s.%s", prefix, key)
		}
		if _, isMap := value.(map[string]interface{}); isMap {
			for fk, fv := range flattenDictDepth(value, newKey, depth+1) {
				items[fk] = fv
			}
		} else {
			items[newKey] = value
		}
	}
	return items
}
