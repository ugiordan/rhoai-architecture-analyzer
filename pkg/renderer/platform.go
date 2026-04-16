package renderer

import (
	"fmt"
	"sort"
	"strings"
)

// RenderPlatformAll renders all platform-wide diagrams from aggregated data,
// returning a filename->content map.
func RenderPlatformAll(data map[string]interface{}) map[string]string {
	return map[string]string{
		"platform-dependencies.mmd":     renderPlatformDependencyGraph(data),
		"platform-crd-ownership.mmd":    renderCRDOwnership(data),
		"platform-rbac-overview.mmd":    renderRBACOverview(data),
		"platform-network-topology.mmd": renderNetworkTopology(data),
		"PLATFORM.md":                   renderPlatformReport(data),
	}
}

func renderPlatformDependencyGraph(data map[string]interface{}) string {
	var b strings.Builder
	b.WriteString("graph LR\n")
	b.WriteString("    %% Platform-wide dependency graph\n")
	b.WriteString("\n")
	b.WriteString("    classDef component fill:#3498db,stroke:#2980b9,color:#fff,stroke-width:2px\n")
	b.WriteString("    classDef gomod fill:#2ecc71,stroke:#27ae60,color:#fff\n")
	b.WriteString("    classDef crdwatch fill:#e74c3c,stroke:#c0392b,color:#fff\n")
	b.WriteString("\n")

	components := getStringSlice(data, "components")
	for _, comp := range components {
		cid := sanitizeID(comp)
		b.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", cid, escapeLabel(comp)))
		b.WriteString(fmt.Sprintf("    class %s component\n", cid))
	}
	b.WriteString("\n")

	depGraph := getSlice(data, "dependency_graph")
	for _, dep := range depGraph {
		fromID := sanitizeID(getStr(dep, "from", ""))
		toID := sanitizeID(getStr(dep, "to", ""))
		depType := getStr(dep, "type", "")
		if strings.HasPrefix(depType, "watches-crd:") {
			crdKind := depType[len("watches-crd:"):]
			b.WriteString(fmt.Sprintf("    %s -->|\"watches %s\"| %s\n",
				fromID, escapeLabel(crdKind), toID))
		} else {
			b.WriteString(fmt.Sprintf("    %s -.->|\"%s\"| %s\n",
				fromID, escapeLabel(depType), toID))
		}
	}

	return strings.TrimRight(b.String(), "\n")
}

func renderCRDOwnership(data map[string]interface{}) string {
	var b strings.Builder
	b.WriteString("graph TD\n")
	b.WriteString("    %% CRD ownership map\n")
	b.WriteString("\n")
	b.WriteString("    classDef component fill:#3498db,stroke:#2980b9,color:#fff\n")
	b.WriteString("    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff\n")
	b.WriteString("\n")

	crdOwners := getMap(data, "crd_ownership")
	if crdOwners == nil {
		return b.String()
	}

	// Sort keys for deterministic output
	kinds := make([]string, 0, len(crdOwners))
	for k := range crdOwners {
		kinds = append(kinds, k)
	}
	sort.Strings(kinds)

	counter := 0
	for _, kind := range kinds {
		owner, ok := crdOwners[kind].(string)
		if !ok {
			continue
		}
		counter++
		ownerID := sanitizeID(owner)
		crdID := fmt.Sprintf("crd_%d", counter)
		b.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", ownerID, escapeLabel(owner)))
		b.WriteString(fmt.Sprintf("    class %s component\n", ownerID))
		b.WriteString(fmt.Sprintf("    %s -->|\"defines\"| %s{{\"%s\"}}\n",
			ownerID, crdID, escapeLabel(kind)))
		b.WriteString(fmt.Sprintf("    class %s crd\n", crdID))
	}

	return strings.TrimRight(b.String(), "\n")
}

func renderRBACOverview(data map[string]interface{}) string {
	var b strings.Builder
	b.WriteString("graph TD\n")
	b.WriteString("    %% Platform RBAC overview\n")
	b.WriteString("\n")
	b.WriteString("    classDef component fill:#3498db,stroke:#2980b9,color:#fff\n")
	b.WriteString("    classDef role fill:#e8a838,stroke:#b07828,color:#fff\n")
	b.WriteString("\n")

	roles := getSlice(data, "rbac_cluster_roles")
	counter := 0
	for _, role := range roles {
		counter++
		owner := getStr(role, "owner", "")
		name := getStr(role, "name", "")
		rules := getSlice(role, "rules")
		totalResources := 0
		for _, rule := range rules {
			totalResources += len(getStringSlice(rule, "resources"))
		}
		ownerID := sanitizeID(owner)
		roleID := fmt.Sprintf("role_%d", counter)
		b.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", ownerID, escapeLabel(owner)))
		b.WriteString(fmt.Sprintf("    class %s component\n", ownerID))
		b.WriteString(fmt.Sprintf("    %s --> %s[\"%s\\n(%d resources)\"]\n",
			ownerID, roleID, escapeLabel(name), totalResources))
		b.WriteString(fmt.Sprintf("    class %s role\n", roleID))
	}

	return strings.TrimRight(b.String(), "\n")
}

func renderNetworkTopology(data map[string]interface{}) string {
	var b strings.Builder
	b.WriteString("graph LR\n")
	b.WriteString("    %% Platform network topology\n")
	b.WriteString("\n")
	b.WriteString("    classDef component fill:#3498db,stroke:#2980b9,color:#fff\n")
	b.WriteString("    classDef service fill:#2ecc71,stroke:#27ae60,color:#fff\n")
	b.WriteString("\n")

	services := getSlice(data, "services")
	counter := 0
	for _, svc := range services {
		counter++
		owner := getStr(svc, "owner", "")
		name := getStr(svc, "name", "")
		svcType := getStr(svc, "type", "ClusterIP")
		ports := getSlice(svc, "ports")
		portParts := make([]string, 0, len(ports))
		for _, p := range ports {
			portParts = append(portParts, fmt.Sprintf("%d/%s",
				getInt(p, "port"), getStr(p, "protocol", "TCP")))
		}
		portStr := strings.Join(portParts, ", ")

		ownerID := sanitizeID(owner)
		svcID := fmt.Sprintf("svc_%d", counter)
		b.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", ownerID, escapeLabel(owner)))
		b.WriteString(fmt.Sprintf("    class %s component\n", ownerID))
		b.WriteString(fmt.Sprintf("    %s --> %s[\"%s\\n%s: %s\"]\n",
			ownerID, svcID, escapeLabel(name), escapeLabel(svcType), escapeLabel(portStr)))
		b.WriteString(fmt.Sprintf("    class %s service\n", svcID))
	}

	return strings.TrimRight(b.String(), "\n")
}
