package extractor

import "sort"

var configmapYAMLPatterns = []string{
	"config/**/*configmap*.yaml",
	"config/**/*configmap*.yml",
	"**/configmap*.yaml",
	"**/configmap*.yml",
	"charts/**/templates/*configmap*.yaml",
	"manifests/**/*configmap*.yaml",
	"deploy/**/*configmap*.yaml",
}

// extractConfigMaps finds ConfigMap definitions and extracts their data keys.
func extractConfigMaps(repoPath string) []ConfigMapRef {
	files := findYAMLFiles(repoPath, configmapYAMLPatterns)
	var configmaps []ConfigMapRef

	for _, fpath := range files {
		for _, doc := range parseYAMLSafe(fpath) {
			kind, _ := doc["kind"].(string)
			if kind != "ConfigMap" {
				continue
			}

			metadata, _ := doc["metadata"].(map[string]interface{})
			name := ""
			if metadata != nil {
				name, _ = metadata["name"].(string)
			}

			var dataKeys []string
			if data, ok := doc["data"].(map[string]interface{}); ok {
				for k := range data {
					dataKeys = append(dataKeys, k)
				}
			}
			if binaryData, ok := doc["binaryData"].(map[string]interface{}); ok {
				for k := range binaryData {
					dataKeys = append(dataKeys, k)
				}
			}
			sort.Strings(dataKeys)
			if dataKeys == nil {
				dataKeys = []string{}
			}

			configmaps = append(configmaps, ConfigMapRef{
				Name:         name,
				DataKeys:     dataKeys,
				ReferencedBy: []string{},
				Source:       relativePath(repoPath, fpath),
			})
		}
	}

	if configmaps == nil {
		configmaps = []ConfigMapRef{}
	}
	return configmaps
}
