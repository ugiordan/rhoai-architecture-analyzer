package extractor

import (
	"os"
	"regexp"
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

// templateServicePatterns covers Go template files that may contain Service definitions.
// Includes both service-prefixed and service-suffixed names (e.g., catalog-service.yaml.tmpl).
var templateServicePatterns = []string{
	"**/service*.tmpl",
	"**/*-service*.tmpl",
	"**/*_service*.tmpl",
	"**/templates/**/*.tmpl",
}

// goTemplateServiceRE matches Go template directives (reuses pattern from network_policies.go)
var goTemplateServiceRE = regexp.MustCompile(`\{\{-?\s*.*?\s*-?\}\}`)

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

	// Also scan Go template files (.tmpl) for Service definitions.
	tmplFiles := findFiles(repoPath, templateServicePatterns)
	for _, fpath := range tmplFiles {
		for _, svc := range parseTemplateServices(repoPath, fpath) {
			services = append(services, svc)
		}
	}

	if services == nil {
		services = []Service{}
	}
	return services
}

// templateConditionRE matches Go template if/else directives to extract conditions.
var templateConditionRE = regexp.MustCompile(`\{\{-?\s*(if|else if|else|end)\s*(.*?)\s*-?\}\}`)

// templatePortCondition maps a port name found in a .tmpl file to the Go
// template condition that governs it (e.g., ".Spec.KubeRBACProxy"). Ports
// in an else branch get the negated condition ("not .Spec.KubeRBACProxy").
type templatePortCondition struct {
	portName  string
	condition string
}

// extractPortConditions scans a Go template file for if/else blocks that
// surround port definitions. Returns a map of port-name to condition string.
func extractPortConditions(content string) map[string]string {
	lines := strings.Split(content, "\n")
	result := make(map[string]string)

	// Track the condition stack for nested if/else blocks.
	type condFrame struct {
		cond   string
		inElse bool
	}
	var stack []condFrame

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check for template condition directives on this line.
		matches := templateConditionRE.FindStringSubmatch(trimmed)
		if len(matches) >= 3 {
			keyword := matches[1]
			arg := strings.TrimSpace(matches[2])
			switch keyword {
			case "if":
				stack = append(stack, condFrame{cond: arg})
			case "else if":
				if len(stack) > 0 {
					stack[len(stack)-1] = condFrame{cond: arg}
				}
			case "else":
				if len(stack) > 0 {
					stack[len(stack)-1].inElse = true
				}
			case "end":
				if len(stack) > 0 {
					stack = stack[:len(stack)-1]
				}
			}
		}

		// Check if this line defines a port name (inside or outside a conditional).
		if strings.Contains(trimmed, "name:") && len(stack) > 0 {
			// Extract the port name value (strip template directives first).
			nameLine := goTemplateServiceRE.ReplaceAllString(trimmed, "template-value")
			parts := strings.SplitN(nameLine, "name:", 2)
			if len(parts) == 2 {
				portName := strings.TrimSpace(parts[1])
				if portName != "" && portName != "template-value" {
					// Use the innermost condition.
					frame := stack[len(stack)-1]
					if frame.inElse {
						result[portName] = "when not " + frame.cond
					} else {
						result[portName] = "when " + frame.cond
					}
				}
			}
		}
	}
	return result
}

// parseTemplateServices strips Go template directives from .tmpl files
// and parses the resulting YAML for Service definitions. Conditional
// blocks (if/else) are analyzed to annotate ports with their governing
// condition, so downstream consumers know which ports are alternatives
// rather than coexisting.
func parseTemplateServices(repoPath, fpath string) []Service {
	data, err := os.ReadFile(fpath)
	if err != nil {
		return nil
	}
	content := string(data)

	if !strings.Contains(content, "Service") {
		return nil
	}

	// Pre-scan: extract condition annotations for port names.
	portConditions := extractPortConditions(content)

	// Strip Go template directives
	cleaned := goTemplateServiceRE.ReplaceAllStringFunc(content, func(match string) string {
		trimmed := strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(match, "}}"), "{{"))
		trimmed = strings.TrimPrefix(trimmed, "-")
		trimmed = strings.TrimSuffix(trimmed, "-")
		trimmed = strings.TrimSpace(trimmed)
		if strings.HasPrefix(trimmed, "if ") || strings.HasPrefix(trimmed, "else") ||
			strings.HasPrefix(trimmed, "end") || strings.HasPrefix(trimmed, "range") ||
			strings.HasPrefix(trimmed, "define") || strings.HasPrefix(trimmed, "template") ||
			strings.HasPrefix(trimmed, "block") || strings.HasPrefix(trimmed, "with") ||
			trimmed == "-" {
			return ""
		}
		return "template-value"
	})

	docs := parseYAMLFromBytes([]byte(cleaned))
	var services []Service
	for _, doc := range docs {
		kind, _ := doc["kind"].(string)
		if kind != "Service" {
			continue
		}
		metadata, _ := doc["metadata"].(map[string]interface{})
		if metadata == nil {
			metadata = map[string]interface{}{}
		}
		name, _ := metadata["name"].(string)

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
			portName, _ := pm["name"].(string)
			protocol, _ := pm["protocol"].(string)
			if protocol == "" {
				protocol = "TCP"
			}
			portVal := pm["port"]
			if portVal == nil {
				portVal = 0
			}
			targetPortVal := pm["targetPort"]
			if targetPortVal == nil {
				targetPortVal = 0
			}
			sp := ServicePort{
				Name:       portName,
				Port:       portVal,
				TargetPort: targetPortVal,
				Protocol:   protocol,
			}
			// Annotate with the condition if this port is inside a conditional block.
			if cond, ok := portConditions[portName]; ok {
				sp.Condition = cond
			}
			ports = append(ports, sp)
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
