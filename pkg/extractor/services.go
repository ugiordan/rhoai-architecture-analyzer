package extractor

import (
	"strings"
)

var serviceSearchPatterns = []string{
	"**/service.yaml",
	"**/service.yml",
	"**/service*.yaml",
	"**/service*.yml",
	"config/**/service*.yaml",
	"config/**/service*.yml",
}

// extractServices scans YAML files for Kubernetes Service definitions.
func extractServices(repoPath string) []Service {
	files := findYAMLFiles(repoPath, serviceSearchPatterns)
	var services []Service

	for _, fpath := range files {
		for _, doc := range parseYAMLSafe(fpath) {
			kind, _ := doc["kind"].(string)
			if kind != "Service" {
				continue
			}
			metadata, _ := doc["metadata"].(map[string]interface{})
			if metadata == nil {
				metadata = map[string]interface{}{}
			}
			name, _ := metadata["name"].(string)
			if hasUnresolvedTemplates(name) {
				continue
			}

			spec, _ := doc["spec"].(map[string]interface{})
			if spec == nil {
				spec = map[string]interface{}{}
			}

			rawPorts, _ := spec["ports"].([]interface{})
			var ports []ServicePort
			for _, p := range rawPorts {
				pm, ok := p.(map[string]interface{})
				if !ok {
					continue
				}
				portVal := pm["port"]
				targetPortVal := pm["targetPort"]
				if hasUnresolvedTemplatesAny(portVal) || hasUnresolvedTemplatesAny(targetPortVal) {
					continue
				}
				portName, _ := pm["name"].(string)
				protocol, _ := pm["protocol"].(string)
				if protocol == "" {
					protocol = "TCP"
				}
				if portVal == nil {
					portVal = 0
				}
				if targetPortVal == nil {
					targetPortVal = 0
				}
				ports = append(ports, ServicePort{
					Name:       portName,
					Port:       portVal,
					TargetPort: targetPortVal,
					Protocol:   protocol,
				})
			}
			if ports == nil {
				ports = []ServicePort{}
			}

			selector, _ := spec["selector"].(map[string]interface{})
			if selector == nil {
				selector = map[string]interface{}{}
			}

			svcType, _ := spec["type"].(string)
			if svcType == "" {
				svcType = "ClusterIP"
			}

			services = append(services, Service{
				Name:     name,
				Source:   relativePath(repoPath, fpath),
				Type:     svcType,
				Ports:    ports,
				Selector: selector,
			})
		}
	}

	if services == nil {
		services = []Service{}
	}
	return services
}

// hasUnresolvedTemplates checks if a string contains Helm template syntax.
func hasUnresolvedTemplates(s string) bool {
	return strings.Contains(s, "{{") && strings.Contains(s, "}}")
}

// hasUnresolvedTemplatesAny checks if an interface value contains Helm templates.
func hasUnresolvedTemplatesAny(v interface{}) bool {
	if s, ok := v.(string); ok {
		return hasUnresolvedTemplates(s)
	}
	return false
}
