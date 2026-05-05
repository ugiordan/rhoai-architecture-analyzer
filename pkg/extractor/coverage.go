package extractor

import (
	"strings"
)

// computeDataCoverage assesses the richness of each extraction section
// and returns a map of section name to coverage level.
// Levels: "none" (0 items), "sparse" (1-2 or all placeholders),
// "moderate" (3-5), "rich" (6+).
func computeDataCoverage(arch *ComponentArchitecture) map[string]string {
	cov := make(map[string]string)

	cov["crds"] = classifyCoverage(len(arch.CRDs), 0)
	cov["services"] = classifyWithTemplates(countServices(arch))
	cov["deployments"] = classifyCoverage(len(arch.Deployments), 0)
	cov["network_policies"] = classifyCoverage(len(arch.NetworkPolicies), 0)
	cov["ingress_routing"] = classifyCoverage(len(arch.IngressRouting), 0)
	cov["webhooks"] = classifyCoverage(len(arch.Webhooks), 0)
	cov["external_connections"] = classifyCoverage(len(arch.ExternalConnections), 0)
	cov["secrets"] = classifyCoverage(len(arch.Secrets), 0)
	cov["configmaps"] = classifyCoverage(len(arch.ConfigMaps), 0)
	cov["dockerfiles"] = classifyCoverage(len(arch.Dockerfiles), 0)
	cov["api_types"] = classifyCoverage(len(arch.APITypes), 0)

	// RBAC: count cluster roles + roles
	rbacCount := 0
	if arch.RBAC != nil {
		rbacCount = len(arch.RBAC.ClusterRoles) + len(arch.RBAC.Roles)
	}
	cov["rbac"] = classifyCoverage(rbacCount, 0)

	return cov
}

// countServices returns (total, templateCount) for services.
func countServices(arch *ComponentArchitecture) (int, int) {
	total := len(arch.Services)
	templates := 0
	for _, svc := range arch.Services {
		if isTemplateItem(svc.Name) {
			templates++
			continue
		}
		for _, p := range svc.Ports {
			if isTemplateItem(p.Name) || isTemplateItem(portToString(p.Port)) {
				templates++
				break
			}
		}
	}
	return total, templates
}

func portToString(port interface{}) string {
	if s, ok := port.(string); ok {
		return s
	}
	return ""
}

// classifyCoverage returns the coverage level based on item count.
func classifyCoverage(count, templateCount int) string {
	if count == 0 {
		return "none"
	}
	if count <= 2 || (templateCount > 0 && templateCount == count) {
		return "sparse"
	}
	if count <= 5 {
		return "moderate"
	}
	return "rich"
}

// classifyWithTemplates classifies using both total and template counts.
func classifyWithTemplates(total, templates int) string {
	return classifyCoverage(total, templates)
}

// isTemplateItem checks if a string value looks like a placeholder.
func isTemplateItem(v string) bool {
	if v == "" {
		return false
	}
	if v == "template-value" {
		return true
	}
	if strings.Contains(v, "$(") {
		return true
	}
	if strings.Contains(v, "${") {
		return true
	}
	return false
}
