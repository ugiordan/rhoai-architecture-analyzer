package extractor

var networkPolicySearchPatterns = []string{
	"**/*networkpolicy*.yaml",
	"**/*networkpolicy*.yml",
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

// extractNetworkPolicies scans YAML files for NetworkPolicy definitions.
func extractNetworkPolicies(repoPath string) []NetworkPolicy {
	files := findYAMLFiles(repoPath, networkPolicySearchPatterns)
	var policies []NetworkPolicy

	for _, fpath := range files {
		for _, doc := range parseYAMLSafe(fpath) {
			kind, _ := doc["kind"].(string)
			if kind != "NetworkPolicy" {
				continue
			}
			metadata, _ := doc["metadata"].(map[string]interface{})
			if metadata == nil {
				metadata = map[string]interface{}{}
			}
			spec, ok := doc["spec"].(map[string]interface{})
			if !ok {
				continue
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

			policies = append(policies, NetworkPolicy{
				Name:         name,
				Source:       relativePath(repoPath, fpath),
				PodSelector:  matchLabels,
				PolicyTypes:  policyTypes,
				IngressRules: extractIngressRules(spec["ingress"]),
				EgressRules:  extractEgressRules(spec["egress"]),
			})
		}
	}

	if policies == nil {
		policies = []NetworkPolicy{}
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
			// In Kubernetes NetworkPolicy, a missing port means "all ports" (wildcard).
			// Preserve nil as-is so JSON output serializes as null, not 0.
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
