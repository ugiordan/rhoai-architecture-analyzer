package extractor

import (
	"sort"
)

var servingRuntimePatterns = []string{
	"config/runtimes/*.yaml",
	"config/runtimes/*.yml",
	"config/runtimes/**/*.yaml",
	"config/servingruntime*.yaml",
	"config/servingruntime*.yml",
	"manifests/**/servingruntime*.yaml",
	"manifests/**/clusterservingruntime*.yaml",
	"charts/**/templates/servingruntime*.yaml",
}

// extractServingRuntimes scans for KServe/ModelMesh ServingRuntime and
// ClusterServingRuntime definitions and returns structured metadata
// including containers, supported model formats, and GPU requirements.
func extractServingRuntimes(repoPath string) []ServingRuntime {
	files := findYAMLFiles(repoPath, servingRuntimePatterns)
	var runtimes []ServingRuntime

	for _, fpath := range files {
		rel := relativePath(repoPath, fpath)
		for _, doc := range parseYAMLSafe(fpath) {
			kind, _ := doc["kind"].(string)
			if kind != "ServingRuntime" && kind != "ClusterServingRuntime" {
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

			spec, _ := doc["spec"].(map[string]interface{})
			if spec == nil {
				continue
			}

			rt := ServingRuntime{
				Name:   name,
				Kind:   kind,
				Source: rel,
			}

			// Multi-model flag
			if mm, ok := spec["multiModel"].(bool); ok {
				rt.MultiModel = mm
			}

			// Disabled flag
			if d, ok := spec["disabled"].(bool); ok {
				rt.Disabled = d
			}

			// Containers
			if containers, ok := spec["containers"].([]interface{}); ok {
				for _, c := range containers {
					cm, ok := c.(map[string]interface{})
					if !ok {
						continue
					}
					sc := ServingContainer{
						Name: getString(cm, "name"),
					}
					if img, ok := cm["image"].(string); ok && !hasUnresolvedTemplates(img) {
						sc.Image = img
					}

					// Resources
					if res, ok := cm["resources"].(map[string]interface{}); ok {
						sc.Resources = res
					}

					// Ports
					if ports, ok := cm["ports"].([]interface{}); ok {
						for _, p := range ports {
							pm, ok := p.(map[string]interface{})
							if !ok {
								continue
							}
							cp := ContainerPort{
								Name:     getString(pm, "name"),
								Protocol: getString(pm, "protocol"),
							}
							if v, ok := pm["containerPort"].(int); ok {
								cp.ContainerPort = v
							} else if v, ok := pm["containerPort"].(float64); ok {
								cp.ContainerPort = int(v)
							}
							sc.Ports = append(sc.Ports, cp)
						}
					}

					// Args (useful for understanding runtime configuration)
					if args, ok := cm["args"].([]interface{}); ok {
						for _, a := range args {
							if s, ok := a.(string); ok {
								sc.Args = append(sc.Args, s)
							}
						}
					}

					rt.Containers = append(rt.Containers, sc)
				}
			}

			// Supported model formats
			if formats, ok := spec["supportedModelFormats"].([]interface{}); ok {
				for _, f := range formats {
					fm, ok := f.(map[string]interface{})
					if !ok {
						continue
					}
					smf := SupportedModelFormat{
						Name: getString(fm, "name"),
					}
					if v, ok := fm["version"].(string); ok {
						smf.Version = v
					}
					if v, ok := fm["autoSelect"].(bool); ok {
						smf.AutoSelect = v
					}
					if v, ok := fm["priority"].(int); ok {
						smf.Priority = v
					} else if v, ok := fm["priority"].(float64); ok {
						smf.Priority = int(v)
					}
					rt.SupportedFormats = append(rt.SupportedFormats, smf)
				}
			}

			runtimes = append(runtimes, rt)
		}
	}

	sort.Slice(runtimes, func(i, j int) bool {
		return runtimes[i].Name < runtimes[j].Name
	})
	return runtimes
}

func getString(m map[string]interface{}, key string) string {
	v, _ := m[key].(string)
	return v
}
