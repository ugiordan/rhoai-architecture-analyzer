package extractor

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var networkPolicySearchPatterns = []string{
	"**/*networkpolic*.yaml",
	"**/*networkpolic*.yml",
	"**/*network-polic*.yaml",
	"**/*network-polic*.yml",
	"**/*netpol*.yaml",
	"**/*netpol*.yml",
	// Directory-based discovery: repos that organize policies by directory
	"**/network-policies/**/*.yaml",
	"**/network-policies/**/*.yml",
	"**/networkpolicies/**/*.yaml",
	"**/networkpolicies/**/*.yml",
	"**/netpol/**/*.yaml",
	"**/netpol/**/*.yml",
}

// templateNetworkPolicyPatterns covers Go template files containing NetworkPolicy defs.
var templateNetworkPolicyPatterns = []string{
	"**/*networkpolicy*.tmpl",
	"**/*network-polic*.tmpl",
	"**/*netpol*.tmpl",
	"**/network-policies/**/*.tmpl",
}

// goTemplateDirectiveRE matches Go template directives like {{ .Foo }} or {{- if ... -}}
var goTemplateDirectiveRE = regexp.MustCompile(`\{\{-?\s*.*?\s*-?\}\}`)

// extractNetworkPolicies scans YAML files for NetworkPolicy definitions.
func extractNetworkPolicies(repoPath string) []NetworkPolicy {
	files := findYAMLFiles(repoPath, networkPolicySearchPatterns)
	var policies []NetworkPolicy

	for _, fpath := range files {
		for _, doc := range parseYAMLSafe(fpath) {
			if np := parseNetworkPolicyDoc(doc, repoPath, fpath); np != nil {
				policies = append(policies, *np)
			}
		}
	}

	// Also scan Go template files (.tmpl) for NetworkPolicy definitions.
	// Strip template directives before parsing as YAML.
	tmplFiles := findFiles(repoPath, templateNetworkPolicyPatterns)
	for _, fpath := range tmplFiles {
		for _, np := range parseTemplateNetworkPolicies(repoPath, fpath) {
			policies = append(policies, np)
		}
	}

	if policies == nil {
		policies = []NetworkPolicy{}
	}
	return policies
}

// parseNetworkPolicyDoc extracts a NetworkPolicy from a YAML document and runs
// security assessment on it.
func parseNetworkPolicyDoc(doc map[string]interface{}, repoPath, fpath string) *NetworkPolicy {
	kind, _ := doc["kind"].(string)
	if kind != "NetworkPolicy" {
		return nil
	}
	metadata, _ := doc["metadata"].(map[string]interface{})
	if metadata == nil {
		metadata = map[string]interface{}{}
	}
	spec, ok := doc["spec"].(map[string]interface{})
	if !ok {
		return nil
	}

	podSelector, _ := spec["podSelector"].(map[string]interface{})
	if podSelector == nil {
		podSelector = map[string]interface{}{}
	}
	matchLabels, _ := podSelector["matchLabels"].(map[string]interface{})
	if matchLabels == nil {
		matchLabels = map[string]interface{}{}
	}

	policyTypes := toStringSlice(spec["policyTypes"])
	name, _ := metadata["name"].(string)

	np := &NetworkPolicy{
		Name:         name,
		Source:       relativePath(repoPath, fpath),
		PodSelector:  matchLabels,
		PolicyTypes:  policyTypes,
		IngressRules: extractIngressRules(spec["ingress"]),
		EgressRules:  extractEgressRules(spec["egress"]),
	}

	// Security assessment
	np.Issues = assessNetworkPolicy(np, spec)
	return np
}

// assessNetworkPolicy checks a NetworkPolicy for common security issues.
func assessNetworkPolicy(np *NetworkPolicy, spec map[string]interface{}) []string {
	var issues []string

	// Check for permissive ingress (empty from or missing ingress entirely)
	ingress, hasIngress := spec["ingress"]
	if hasIngress {
		rules, ok := ingress.([]interface{})
		if ok {
			for i, r := range rules {
				rule, ok := r.(map[string]interface{})
				if !ok {
					continue
				}
				from, hasFrom := rule["from"]
				ports, hasPorts := rule["ports"]

				if !hasFrom || from == nil {
					issues = append(issues, fmt.Sprintf(
						"ingress rule %d allows traffic from all sources (no 'from' selector)", i+1))
				} else if fromList, ok := from.([]interface{}); ok && len(fromList) == 0 {
					issues = append(issues, fmt.Sprintf(
						"ingress rule %d has empty 'from' list (allows all sources)", i+1))
				}

				if !hasPorts || ports == nil {
					issues = append(issues, fmt.Sprintf(
						"ingress rule %d allows all ports (no port restriction)", i+1))
				} else if portList, ok := ports.([]interface{}); ok && len(portList) == 0 {
					issues = append(issues, fmt.Sprintf(
						"ingress rule %d has empty ports list (allows all ports)", i+1))
				}
			}
		}
	}

	// Check for missing policyTypes
	if len(np.PolicyTypes) == 0 {
		issues = append(issues, "no policyTypes specified (defaults to Ingress only, egress is unrestricted)")
	} else {
		hasEgress := false
		for _, pt := range np.PolicyTypes {
			if pt == "Egress" {
				hasEgress = true
			}
		}
		if !hasEgress {
			issues = append(issues, "policyTypes does not include Egress (all egress traffic is allowed)")
		}
	}

	// Check for empty pod selector (applies to all pods in namespace)
	if len(np.PodSelector) == 0 {
		issues = append(issues, "empty podSelector (applies to all pods in namespace)")
	}

	return issues
}

// parseTemplateNetworkPolicies strips Go template directives from .tmpl files
// and parses the resulting YAML for NetworkPolicy definitions.
func parseTemplateNetworkPolicies(repoPath, fpath string) []NetworkPolicy {
	data, err := os.ReadFile(fpath)
	if err != nil {
		return nil
	}
	content := string(data)

	// Only proceed if it looks like a NetworkPolicy template
	if !strings.Contains(content, "NetworkPolicy") {
		return nil
	}

	// Strip Go template directives, replacing with empty strings or placeholder values
	cleaned := goTemplateDirectiveRE.ReplaceAllStringFunc(content, func(match string) string {
		// For template values used in string contexts, replace with a placeholder
		trimmed := strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(match, "}}"), "{{"))
		trimmed = strings.TrimPrefix(trimmed, "-")
		trimmed = strings.TrimSuffix(trimmed, "-")
		trimmed = strings.TrimSpace(trimmed)
		// Control flow directives (if, range, end, define, etc.) -> empty
		if strings.HasPrefix(trimmed, "if ") || strings.HasPrefix(trimmed, "else") ||
			strings.HasPrefix(trimmed, "end") || strings.HasPrefix(trimmed, "range") ||
			strings.HasPrefix(trimmed, "define") || strings.HasPrefix(trimmed, "template") ||
			strings.HasPrefix(trimmed, "block") || strings.HasPrefix(trimmed, "with") ||
			trimmed == "-" {
			return ""
		}
		// Value expressions (.Name, .Spec.Port) -> placeholder
		return "template-value"
	})

	docs := parseYAMLFromBytes([]byte(cleaned))
	var policies []NetworkPolicy
	for _, doc := range docs {
		if np := parseNetworkPolicyDoc(doc, repoPath, fpath); np != nil {
			policies = append(policies, *np)
		}
	}
	return policies
}

// extractIngressRules extracts ingress rules from a NetworkPolicy spec.
func extractIngressRules(v interface{}) []map[string]interface{} {
	return extractDirectionalRules(v, "from")
}

// extractEgressRules extracts egress rules from a NetworkPolicy spec.
func extractEgressRules(v interface{}) []map[string]interface{} {
	return extractDirectionalRules(v, "to")
}

// extractDirectionalRules extracts ingress or egress rules from a NetworkPolicy
// spec, using directionKey ("from" for ingress, "to" for egress) to extract the
// peer selectors.
func extractDirectionalRules(v interface{}, directionKey string) []map[string]interface{} {
	rules, ok := v.([]interface{})
	if !ok {
		return []map[string]interface{}{}
	}
	var result []map[string]interface{}
	for _, r := range rules {
		rule, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		var ports []map[string]interface{}
		rawPorts := toSliceOfMaps(rule["ports"])
		for _, p := range rawPorts {
			port := p["port"]
			protocol, _ := p["protocol"].(string)
			if protocol == "" {
				protocol = "TCP"
			}
			ports = append(ports, map[string]interface{}{
				"port":     port,
				"protocol": protocol,
			})
		}
		if ports == nil {
			ports = []map[string]interface{}{}
		}
		selectors, _ := rule[directionKey].([]interface{})
		if selectors == nil {
			selectors = []interface{}{}
		}
		result = append(result, map[string]interface{}{
			"ports":      ports,
			directionKey: selectors,
		})
	}
	if result == nil {
		result = []map[string]interface{}{}
	}
	return result
}
