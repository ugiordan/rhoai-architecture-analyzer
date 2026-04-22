package renderer

import (
	"fmt"
	"sort"
	"strings"
)

// DocsPage represents a single documentation page.
type DocsPage struct {
	Path    string // relative path, e.g. "index.md" or "components/kserve/index.md"
	Content string
}

// RenderDocs generates a complete set of documentation pages from architecture
// data. It auto-detects whether the input is a single component or aggregated
// platform data (by checking for the "platform" key). Returns a slice of pages
// with relative paths suitable for embedding in a docs site.
func RenderDocs(data map[string]interface{}) []DocsPage {
	if _, ok := data["platform"]; ok {
		return renderPlatformDocs(data)
	}
	return renderComponentDocs(data, "")
}

// NavSnippet generates a YAML navigation snippet for mkdocs.yml from the
// rendered docs pages. Component sub-pages are nested under their component
// directory to match the mkdocs nav hierarchy.
func NavSnippet(pages []DocsPage, prefix string) string {
	var b strings.Builder

	// Group pages by component directory (e.g., "components/kserve/")
	type pageEntry struct {
		path  string
		title string
	}
	topLevel := make([]pageEntry, 0)
	compPages := make(map[string][]pageEntry)
	var compOrder []string

	for _, p := range pages {
		path := p.Path
		if prefix != "" {
			path = prefix + "/" + path
		}
		title := extractTitle(p.Content)

		// Check if this is a component sub-page
		parts := strings.Split(p.Path, "/")
		if len(parts) >= 3 && parts[0] == "components" {
			compDir := parts[1]
			if _, seen := compPages[compDir]; !seen {
				compOrder = append(compOrder, compDir)
			}
			compPages[compDir] = append(compPages[compDir], pageEntry{path, title})
		} else {
			topLevel = append(topLevel, pageEntry{path, title})
		}
	}

	// Write top-level pages
	for _, p := range topLevel {
		b.WriteString(fmt.Sprintf("    - %s: %s\n", p.title, p.path))
	}

	// Write component groups
	if len(compOrder) > 0 {
		b.WriteString("    - Components:\n")
		for _, comp := range compOrder {
			// Use the component directory name as label
			b.WriteString(fmt.Sprintf("      - %s:\n", comp))
			for _, p := range compPages[comp] {
				b.WriteString(fmt.Sprintf("        - %s: %s\n", p.title, p.path))
			}
		}
	}

	return b.String()
}

func extractTitle(content string) string {
	for _, line := range strings.SplitN(content, "\n", 5) {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}
	return "Untitled"
}

// renderComponentDocs generates docs pages for a single component.
func renderComponentDocs(data map[string]interface{}, pathPrefix string) []DocsPage {
	component := getStr(data, "component", "unknown")
	var pages []DocsPage

	prefix := pathPrefix
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	pages = append(pages, DocsPage{
		Path:    prefix + "index.md",
		Content: renderComponentIndexPage(data, component),
	})
	pages = append(pages, DocsPage{
		Path:    prefix + "network.md",
		Content: renderComponentNetworkPage(data, component),
	})
	pages = append(pages, DocsPage{
		Path:    prefix + "rbac.md",
		Content: renderComponentRBACPage(data, component),
	})
	pages = append(pages, DocsPage{
		Path:    prefix + "security.md",
		Content: renderComponentSecurityPage(data, component),
	})
	// Cache page: only emit if cache data exists
	if getMap(data, "cache_config") != nil {
		pages = append(pages, DocsPage{
			Path:    prefix + "cache.md",
			Content: renderComponentCachePage(data, component),
		})
	}
	pages = append(pages, DocsPage{
		Path:    prefix + "dataflow.md",
		Content: renderComponentDataflowPage(data, component),
	})

	return pages
}

// --- Component index page ---

func renderComponentIndexPage(data map[string]interface{}, component string) string {
	var b strings.Builder

	repo := getStr(data, "repo", "")
	version := getStr(data, "analyzer_version", "")
	extractedAt := getStr(data, "extracted_at", "")

	b.WriteString(fmt.Sprintf("# %s\n\n", component))
	if repo != "" {
		b.WriteString(fmt.Sprintf("**Repository:** %s  \n", repo))
	}
	b.WriteString(fmt.Sprintf("**Analyzer:** arch-analyzer %s  \n", version))
	b.WriteString(fmt.Sprintf("**Extracted:** %s\n\n", extractedAt))

	// Summary table
	crds := getSlice(data, "crds")
	services := getSlice(data, "services")
	secrets := getSlice(data, "secrets_referenced")
	rbac := getMap(data, "rbac")
	deployments := getSlice(data, "deployments")
	watches := getSlice(data, "controller_watches")

	clusterRoleCount := 0
	if rbac != nil {
		clusterRoleCount = len(getSlice(rbac, "cluster_roles"))
	}

	b.WriteString("## Summary\n\n")
	b.WriteString("| Metric | Count |\n")
	b.WriteString("|--------|-------|\n")
	b.WriteString(fmt.Sprintf("| CRDs | %d |\n", len(crds)))
	b.WriteString(fmt.Sprintf("| Deployments | %d |\n", len(deployments)))
	b.WriteString(fmt.Sprintf("| Services | %d |\n", len(services)))
	b.WriteString(fmt.Sprintf("| Secrets | %d |\n", len(secrets)))
	b.WriteString(fmt.Sprintf("| Cluster Roles | %d |\n", clusterRoleCount))
	b.WriteString(fmt.Sprintf("| Controller Watches | %d |\n", len(watches)))
	b.WriteString("\n")

	// Component architecture diagram (inline mermaid)
	b.WriteString("## Component Architecture\n\n")
	b.WriteString("CRDs, controllers, and owned Kubernetes resources.\n\n")
	componentDiagram := (&ComponentRenderer{}).Render(data)
	b.WriteString("```mermaid\n")
	b.WriteString(componentDiagram)
	b.WriteString("\n```\n\n")

	// CRD table
	renderCRDTable(&b, data)

	// Dependencies
	b.WriteString("## Dependencies\n\n")
	renderDependencySection(&b, data)

	return b.String()
}

// --- Component network page ---

func renderComponentNetworkPage(data map[string]interface{}, component string) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# %s: Network\n\n", component))

	services := getSlice(data, "services")
	netpols := getSlice(data, "network_policies")

	if len(services) > 0 {
		// Dedup services by (name, type, ports) to avoid test fixture duplicates
		type svcKey struct {
			name, svcType, ports string
		}
		type dedupedSvc struct {
			svc           map[string]interface{}
			count         int
			hasProduction bool
			isTest        bool
		}
		seen := make(map[svcKey]*dedupedSvc)
		var order []svcKey
		for _, svc := range services {
			name := getStr(svc, "name", "")
			svcType := getStr(svc, "type", "ClusterIP")
			ports := getSlice(svc, "ports")
			var portParts []string
			for _, p := range ports {
				portParts = append(portParts, fmt.Sprintf("%d/%s", getInt(p, "port"), getStr(p, "protocol", "TCP")))
			}
			sort.Strings(portParts)
			key := svcKey{name, svcType, strings.Join(portParts, ",")}
			source := getStr(svc, "source", "")
			isTest := strings.Contains(source, "/test/") || strings.Contains(source, "/e2e/") || strings.Contains(source, "/testdata/")
			if existing, ok := seen[key]; ok {
				existing.count++
				if isTest {
					existing.isTest = true
				} else {
					existing.hasProduction = true
				}
			} else {
				seen[key] = &dedupedSvc{svc: svc, count: 1, hasProduction: !isTest, isTest: isTest}
				order = append(order, key)
			}
		}

		// Only show unique services in the diagram
		b.WriteString("## Service Map\n\n")
		if len(order) < len(services) {
			b.WriteString(fmt.Sprintf("*%d unique services (%d total, duplicates from test fixtures collapsed).*\n\n", len(order), len(services)))
		}
		b.WriteString("```mermaid\n")
		b.WriteString("graph LR\n")
		b.WriteString("    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff\n")
		b.WriteString("    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff\n")
		b.WriteString("    classDef component fill:#3498db,stroke:#2980b9,color:#fff\n")
		b.WriteString("    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff\n\n")
		compID := sanitizeID(component)
		b.WriteString(fmt.Sprintf("    %s[\"%s\"]:::%s\n", compID, escapeLabel(component), "component"))
		for i, key := range order {
			ds := seen[key]
			svcID := fmt.Sprintf("svc_%d", i)
			cls := "svc"
			if ds.isTest && !ds.hasProduction {
				cls = "test"
			}
			b.WriteString(fmt.Sprintf("    %s --> %s[\"%s\\n%s: %s\"]:::%s\n",
				compID, svcID, escapeLabel(key.name), key.svcType, key.ports, cls))
		}
		// External connections
		extConns := getSlice(data, "external_connections")
		extSeen := make(map[string]bool)
		for _, ec := range extConns {
			service := getStr(ec, "service", "")
			ecType := getStr(ec, "type", "")
			if service == "" || extSeen[service] {
				continue
			}
			extSeen[service] = true
			extID := sanitizeID("ext_" + service)
			b.WriteString(fmt.Sprintf("    %s -.-> %s[[\"%s\\n%s\"]]:::ext\n",
				compID, extID, escapeLabel(service), ecType))
		}
		b.WriteString("```\n\n")
	}

	renderServiceTable(&b, data)
	renderIngressTable(&b, data)
	renderNetworkPolicyTable(&b, data)

	if len(netpols) == 0 {
		b.WriteString("!!! warning \"No Network Policies\"\n")
		b.WriteString("    No NetworkPolicy resources found. All pod-to-pod traffic is allowed by default.\n\n")
	} else {
		// NetworkPolicy graph showing ingress/egress flows
		renderNetworkPolicyGraph(&b, netpols, component)
	}

	return b.String()
}

// --- Component RBAC page ---

func renderComponentRBACPage(data map[string]interface{}, component string) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# %s: RBAC\n\n", component))

	b.WriteString("ServiceAccount bindings, roles, and resource permissions.\n\n")

	// Check if RBAC data has actual content
	rbac := getMap(data, "rbac")
	hasBindings := rbac != nil && (len(getSlice(rbac, "cluster_role_bindings")) > 0 || len(getSlice(rbac, "role_bindings")) > 0)
	hasRoles := rbac != nil && (len(getSlice(rbac, "cluster_roles")) > 0 || len(getSlice(rbac, "roles")) > 0)

	if !hasBindings && !hasRoles {
		b.WriteString("!!! info \"No RBAC Resources\"\n")
		b.WriteString("    This component does not define any ClusterRoles, Roles, or RoleBindings.\n")
		b.WriteString("    It may operate as a sidecar or library used by other components that define their own RBAC.\n\n")
		return b.String()
	}

	// Always render an RBAC graph. For large surfaces, use a clustered/simplified graph.
	rbacDiagram := (&RBACRenderer{}).Render(data)
	lineCount := strings.Count(rbacDiagram, "\n")

	if lineCount > 60 {
		// Large RBAC: show a clustered summary graph grouping roles by scope
		b.WriteString("## RBAC Overview\n\n")
		b.WriteString(fmt.Sprintf("This component defines a large RBAC surface (%d diagram lines). The graph below groups roles by permission scope.\n\n", lineCount))
		renderClusteredRBACGraph(&b, rbac)
	} else {
		// Small enough for full inline diagram
		b.WriteString("## RBAC Hierarchy\n\n")
		b.WriteString("```mermaid\n")
		b.WriteString(rbacDiagram)
		b.WriteString("\n```\n\n")
	}

	// Bindings table
	if hasBindings {
		b.WriteString("## Bindings\n\n")
		b.WriteString("Subject-to-role mappings defining who has access to what.\n\n")
		b.WriteString("| Binding | Type | Role | Subject |\n")
		b.WriteString("|---------|------|------|---------|\n")
		for _, binding := range getSlice(rbac, "cluster_role_bindings") {
			name := getStr(binding, "name", "")
			roleRef := getStr(binding, "role_ref", "")
			for _, subj := range getSlice(binding, "subjects") {
				b.WriteString(fmt.Sprintf("| %s | ClusterRoleBinding | %s | %s/%s |\n",
					name, roleRef, getStr(subj, "kind", ""), getStr(subj, "name", "")))
			}
		}
		for _, binding := range getSlice(rbac, "role_bindings") {
			name := getStr(binding, "name", "")
			roleRef := getStr(binding, "role_ref", "")
			for _, subj := range getSlice(binding, "subjects") {
				b.WriteString(fmt.Sprintf("| %s | RoleBinding | %s | %s/%s |\n",
					name, roleRef, getStr(subj, "kind", ""), getStr(subj, "name", "")))
			}
		}
		b.WriteString("\n")
	}

	// Full per-rule detail table (the only table, not duplicated)
	if hasRoles {
		b.WriteString("## Role Details\n\n")
		b.WriteString("Per-rule breakdown of API groups, resources, and verbs for each role.\n\n")
		b.WriteString("| Role | Kind | API Groups | Resources | Verbs |\n")
		b.WriteString("|------|------|------------|-----------|-------|\n")
		for _, role := range getSlice(rbac, "cluster_roles") {
			roleName := getStr(role, "name", "")
			for _, rule := range getSlice(role, "rules") {
				apiGroups := getStringSlice(rule, "api_groups")
				resources := getStringSlice(rule, "resources")
				verbs := getStringSlice(rule, "verbs")
				b.WriteString(fmt.Sprintf("| %s | ClusterRole | %s | %s | %s |\n",
					escapeMdCell(roleName),
					escapeMdCell(strings.Join(apiGroups, ", ")),
					escapeMdCell(strings.Join(resources, ", ")),
					escapeMdCell(strings.Join(verbs, ", "))))
			}
		}
		for _, role := range getSlice(rbac, "roles") {
			roleName := getStr(role, "name", "")
			for _, rule := range getSlice(role, "rules") {
				apiGroups := getStringSlice(rule, "api_groups")
				resources := getStringSlice(rule, "resources")
				verbs := getStringSlice(rule, "verbs")
				b.WriteString(fmt.Sprintf("| %s | Role | %s | %s | %s |\n",
					escapeMdCell(roleName),
					escapeMdCell(strings.Join(apiGroups, ", ")),
					escapeMdCell(strings.Join(resources, ", ")),
					escapeMdCell(strings.Join(verbs, ", "))))
			}
		}
		b.WriteString("\n")
	}

	renderRBACSection(&b, data)
	return b.String()
}

// --- Component security page ---

func renderComponentSecurityPage(data map[string]interface{}, component string) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# %s: Security\n\n", component))

	// Secrets
	b.WriteString("## Secrets\n\n")
	b.WriteString("Kubernetes secrets referenced by this component. Only names and types are shown, not values.\n\n")
	renderSecretTable(&b, data)

	// Deployment security
	b.WriteString("## Deployment Security Controls\n\n")
	b.WriteString("SecurityContext settings on pod and container specs. These control privilege escalation, filesystem access, and user identity.\n\n")
	renderSecurityContextSection(&b, data)

	// Dockerfiles
	b.WriteString("## Build Security\n\n")
	b.WriteString("Dockerfile patterns and base image analysis. Covers supply chain security: base images, build stages, runtime user, FIPS compliance.\n\n")
	renderDockerfileTable(&b, data)

	return b.String()
}

// --- Component cache page ---

func renderComponentCachePage(data map[string]interface{}, component string) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# %s: Cache Architecture\n\n", component))
	b.WriteString("Controller-runtime cache configuration controls which Kubernetes resources are cached in-memory. ")
	b.WriteString("Misconfigured caches (cluster-wide watches on high-cardinality types without filters) are a primary cause of operator OOM kills.\n\n")

	renderCacheSection(&b, data)

	return b.String()
}

// --- Component dataflow page ---

func renderComponentDataflowPage(data map[string]interface{}, component string) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# %s: Dataflow\n\n", component))

	// Controller watch table
	b.WriteString("## Controller Watches\n\n")
	b.WriteString("Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.\n\n")
	renderControllerWatchTable(&b, data)

	// Dataflow diagram
	dataflowDiagram := (&DataflowRenderer{}).Render(data)
	if dataflowDiagram != "" {
		b.WriteString("## Reconciliation Flow\n\n")
		b.WriteString("How the controller interacts with the Kubernetes API during reconciliation.\n\n")
		b.WriteString("```mermaid\n")
		b.WriteString(dataflowDiagram)
		b.WriteString("\n```\n\n")
	}

	// Webhooks
	renderWebhookTable(&b, data)
	renderHTTPEndpointTable(&b, data)

	// Configuration
	b.WriteString("## Configuration\n\n")
	b.WriteString("ConfigMaps and Helm values that control this component's runtime behavior.\n\n")
	renderConfigMapTable(&b, data)
	renderHelmSection(&b, data)

	return b.String()
}

// renderClusteredRBACGraph generates a simplified RBAC graph for large surfaces.
// Groups roles by permission scope (wide/medium/narrow) and shows subject bindings.
func renderClusteredRBACGraph(b *strings.Builder, rbac map[string]interface{}) {
	var roles []RoleSummary
	for _, role := range getSlice(rbac, "cluster_roles") {
		roles = append(roles, computeRoleSummary(role, "ClusterRole"))
	}
	for _, role := range getSlice(rbac, "roles") {
		roles = append(roles, computeRoleSummary(role, "Role"))
	}

	// Classify by scope
	var wide, medium, narrow []RoleSummary
	for _, r := range roles {
		if r.ResourceCount > 30 || r.HasWildcard {
			wide = append(wide, r)
		} else if r.ResourceCount > 10 {
			medium = append(medium, r)
		} else {
			narrow = append(narrow, r)
		}
	}

	b.WriteString("```mermaid\n")
	b.WriteString("graph LR\n")
	b.WriteString("    classDef wide fill:#e74c3c,stroke:#c0392b,color:#fff\n")
	b.WriteString("    classDef medium fill:#f39c12,stroke:#d68910,color:#fff\n")
	b.WriteString("    classDef narrow fill:#2ecc71,stroke:#27ae60,color:#fff\n")
	b.WriteString("    classDef subject fill:#3498db,stroke:#2980b9,color:#fff\n\n")

	writeRoleNodes := func(roles []RoleSummary, cls string) {
		for _, r := range roles {
			id := sanitizeID(r.Name)
			label := escapeLabel(r.Name)
			if r.ResourceCount > 0 {
				label += fmt.Sprintf("\\n%d resources", r.ResourceCount)
			}
			if r.HasWildcard {
				label += "\\n!! wildcard"
			}
			b.WriteString(fmt.Sprintf("    %s[\"%s\"]:::%s\n", id, label, cls))
		}
	}

	if len(wide) > 0 {
		b.WriteString("    subgraph wide[\"Wide Scope (>30 resources)\"]\n")
		writeRoleNodes(wide, "wide")
		b.WriteString("    end\n")
	}
	if len(medium) > 0 {
		b.WriteString("    subgraph med[\"Medium Scope (10-30)\"]\n")
		writeRoleNodes(medium, "medium")
		b.WriteString("    end\n")
	}
	if len(narrow) > 0 {
		b.WriteString("    subgraph nar[\"Narrow Scope (<10)\"]\n")
		writeRoleNodes(narrow, "narrow")
		b.WriteString("    end\n")
	}

	// Show bindings as edges from subjects to roles
	b.WriteString("\n")
	subjectsSeen := make(map[string]bool)
	for _, binding := range getSlice(rbac, "cluster_role_bindings") {
		roleRef := getStr(binding, "role_ref", "")
		for _, subj := range getSlice(binding, "subjects") {
			subjName := getStr(subj, "name", "")
			subjKind := getStr(subj, "kind", "")
			subjID := sanitizeID("subj_" + subjName)
			if !subjectsSeen[subjID] {
				subjectsSeen[subjID] = true
				b.WriteString(fmt.Sprintf("    %s[\"%s\\n%s\"]:::subject\n", subjID, escapeLabel(subjName), subjKind))
			}
			b.WriteString(fmt.Sprintf("    %s -->|binds| %s\n", subjID, sanitizeID(roleRef)))
		}
	}
	for _, binding := range getSlice(rbac, "role_bindings") {
		roleRef := getStr(binding, "role_ref", "")
		for _, subj := range getSlice(binding, "subjects") {
			subjName := getStr(subj, "name", "")
			subjKind := getStr(subj, "kind", "")
			subjID := sanitizeID("subj_" + subjName)
			if !subjectsSeen[subjID] {
				subjectsSeen[subjID] = true
				b.WriteString(fmt.Sprintf("    %s[\"%s\\n%s\"]:::subject\n", subjID, escapeLabel(subjName), subjKind))
			}
			b.WriteString(fmt.Sprintf("    %s -->|binds| %s\n", subjID, sanitizeID(roleRef)))
		}
	}

	b.WriteString("```\n\n")
}

// renderNetworkPolicyGraph generates a Mermaid graph showing NetworkPolicy ingress/egress rules.
func renderNetworkPolicyGraph(b *strings.Builder, netpols []map[string]interface{}, component string) {
	b.WriteString("## Network Policy Graph\n\n")
	b.WriteString("Visual representation of NetworkPolicy rules. Ingress rules show what traffic is allowed into pods, egress rules show what traffic is allowed out.\n\n")

	b.WriteString("```mermaid\n")
	b.WriteString("graph LR\n")
	b.WriteString("    classDef policy fill:#e74c3c,stroke:#c0392b,color:#fff\n")
	b.WriteString("    classDef pod fill:#3498db,stroke:#2980b9,color:#fff\n")
	b.WriteString("    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff\n\n")

	compID := sanitizeID(component)
	b.WriteString(fmt.Sprintf("    %s[\"%s\\nPods\"]:::pod\n", compID, escapeLabel(component)))

	for i, np := range netpols {
		npName := getStr(np, "name", fmt.Sprintf("policy-%d", i))
		npID := sanitizeID(fmt.Sprintf("np_%d_%s", i, npName))
		policyTypes := getStringSlice(np, "policy_types")
		typeLabel := strings.Join(policyTypes, ", ")
		if typeLabel == "" {
			typeLabel = "Ingress"
		}
		b.WriteString(fmt.Sprintf("    %s{{\"%s\\n%s\"}}:::policy\n", npID, escapeLabel(npName), typeLabel))

		// Show ingress sources
		ingressRules := getSlice(np, "ingress")
		for j, rule := range ingressRules {
			fromPeers := getSlice(rule, "from")
			if len(fromPeers) == 0 {
				srcID := fmt.Sprintf("src_%d_%d_all", i, j)
				b.WriteString(fmt.Sprintf("    %s[\"all pods\"]:::external\n", srcID))
				b.WriteString(fmt.Sprintf("    %s -->|ingress| %s\n", srcID, npID))
			}
			for k, peer := range fromPeers {
				srcID := fmt.Sprintf("src_%d_%d_%d", i, j, k)
				label := describeNetworkPeer(peer)
				b.WriteString(fmt.Sprintf("    %s[\"%s\"]:::external\n", srcID, escapeLabel(label)))
				b.WriteString(fmt.Sprintf("    %s -->|ingress| %s\n", srcID, npID))
			}
		}
		b.WriteString(fmt.Sprintf("    %s --> %s\n", npID, compID))

		// Show egress destinations
		egressRules := getSlice(np, "egress")
		for j, rule := range egressRules {
			toPeers := getSlice(rule, "to")
			if len(toPeers) == 0 {
				dstID := fmt.Sprintf("dst_%d_%d_all", i, j)
				b.WriteString(fmt.Sprintf("    %s[\"all destinations\"]:::external\n", dstID))
				b.WriteString(fmt.Sprintf("    %s -->|egress| %s\n", compID, dstID))
			}
			for k, peer := range toPeers {
				dstID := fmt.Sprintf("dst_%d_%d_%d", i, j, k)
				label := describeNetworkPeer(peer)
				b.WriteString(fmt.Sprintf("    %s[\"%s\"]:::external\n", dstID, escapeLabel(label)))
				b.WriteString(fmt.Sprintf("    %s -->|egress| %s\n", compID, dstID))
			}
		}
	}

	b.WriteString("```\n\n")
}

// describeNetworkPeer generates a human-readable label for a NetworkPolicy peer.
func describeNetworkPeer(peer map[string]interface{}) string {
	if ns := getMap(peer, "namespaceSelector"); ns != nil {
		labels := getMap(ns, "matchLabels")
		if labels != nil {
			var parts []string
			for k, v := range labels {
				parts = append(parts, fmt.Sprintf("%s=%v", k, v))
			}
			sort.Strings(parts)
			return "ns: " + strings.Join(parts, ", ")
		}
		return "all namespaces"
	}
	if pod := getMap(peer, "podSelector"); pod != nil {
		labels := getMap(pod, "matchLabels")
		if labels != nil {
			var parts []string
			for k, v := range labels {
				parts = append(parts, fmt.Sprintf("%s=%v", k, v))
			}
			sort.Strings(parts)
			return "pods: " + strings.Join(parts, ", ")
		}
		return "all pods"
	}
	if ipBlock := getMap(peer, "ipBlock"); ipBlock != nil {
		cidr := getStr(ipBlock, "cidr", "")
		return "CIDR: " + cidr
	}
	return "peer"
}
