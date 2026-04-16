package renderer

import (
	"fmt"
	"sort"
	"strings"
)

// ReportRenderer generates a structured markdown report for a component.
type ReportRenderer struct{}

func (r *ReportRenderer) Filename() string { return "component-report.md" }

func (r *ReportRenderer) Render(data map[string]interface{}) string {
	var b strings.Builder

	component := getStr(data, "component", "unknown")
	repo := getStr(data, "repo", "")
	version := getStr(data, "analyzer_version", "")
	extractedAt := getStr(data, "extracted_at", "")

	// Header
	b.WriteString(fmt.Sprintf("# %s\n\n", component))
	b.WriteString(fmt.Sprintf("**Repository:** %s  \n", repo))
	b.WriteString(fmt.Sprintf("**Analyzer Version:** %s  \n", version))
	b.WriteString(fmt.Sprintf("**Extracted:** %s  \n\n", extractedAt))
	b.WriteString("---\n\n")

	// APIs Exposed
	b.WriteString("## APIs Exposed\n\n")
	renderCRDTable(&b, data)
	renderWebhookTable(&b, data)
	renderHTTPEndpointTable(&b, data)

	// Dependencies
	b.WriteString("## Dependencies\n\n")
	renderDependencySection(&b, data)

	// Network Architecture
	b.WriteString("## Network Architecture\n\n")
	renderServiceTable(&b, data)
	renderIngressTable(&b, data)
	renderNetworkPolicyTable(&b, data)

	// Security
	b.WriteString("## Security\n\n")
	renderRBACSection(&b, data)
	renderSecretTable(&b, data)
	renderSecurityContextSection(&b, data)

	// Configuration
	b.WriteString("## Configuration\n\n")
	renderConfigMapTable(&b, data)
	renderHelmSection(&b, data)

	// Build
	b.WriteString("## Build\n\n")
	renderDockerfileTable(&b, data)

	// Controller Watches
	b.WriteString("## Controller Watches\n\n")
	renderControllerWatchTable(&b, data)

	// Cache Architecture
	renderCacheSection(&b, data)

	return b.String()
}

func renderCRDTable(b *strings.Builder, data map[string]interface{}) {
	crds := getSlice(data, "crds")
	if len(crds) == 0 {
		b.WriteString("### CRDs\n\nNo CRDs defined.\n\n")
		return
	}
	b.WriteString("### CRDs\n\n")
	b.WriteString("| Group | Version | Kind | Scope | Fields | Validation Rules | Source |\n")
	b.WriteString("|-------|---------|------|-------|--------|------------------|--------|\n")
	for _, crd := range crds {
		rules := getStringSlice(crd, "validation_rules")
		ruleCount := len(rules)
		b.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %d | %d | `%s` |\n",
			getStr(crd, "group", ""), getStr(crd, "version", ""),
			getStr(crd, "kind", ""), getStr(crd, "scope", ""),
			getInt(crd, "fields_count"), ruleCount,
			getStr(crd, "source", "")))
	}
	b.WriteString("\n")
}

func renderWebhookTable(b *strings.Builder, data map[string]interface{}) {
	webhooks := getSlice(data, "webhooks")
	if len(webhooks) == 0 {
		return
	}
	b.WriteString("### Webhooks\n\n")
	b.WriteString("| Name | Type | Path | Failure Policy | Service | Source |\n")
	b.WriteString("|------|------|------|----------------|---------|--------|\n")
	for _, wh := range webhooks {
		b.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | `%s` |\n",
			getStr(wh, "name", ""), getStr(wh, "type", ""),
			getStr(wh, "path", ""), getStr(wh, "failure_policy", ""),
			getStr(wh, "service_ref", ""), getStr(wh, "source", "")))
	}
	b.WriteString("\n")
}

func renderHTTPEndpointTable(b *strings.Builder, data map[string]interface{}) {
	endpoints := getSlice(data, "http_endpoints")
	if len(endpoints) == 0 {
		return
	}
	b.WriteString("### HTTP Endpoints\n\n")
	b.WriteString("| Method | Path | Source |\n")
	b.WriteString("|--------|------|--------|\n")
	for _, ep := range endpoints {
		b.WriteString(fmt.Sprintf("| %s | %s | `%s` |\n",
			getStr(ep, "method", ""), getStr(ep, "path", ""),
			getStr(ep, "source", "")))
	}
	b.WriteString("\n")
}

func renderDependencySection(b *strings.Builder, data map[string]interface{}) {
	deps := getMap(data, "dependencies")
	if deps == nil {
		b.WriteString("No dependencies found.\n\n")
		return
	}

	internalODH := getSlice(deps, "internal_odh")
	if len(internalODH) > 0 {
		b.WriteString("### Internal RHOAI Dependencies\n\n")
		b.WriteString("| Component | Interaction |\n")
		b.WriteString("|-----------|-------------|\n")
		for _, dep := range internalODH {
			b.WriteString(fmt.Sprintf("| %s | %s |\n",
				getStr(dep, "component", ""), getStr(dep, "interaction", "")))
		}
		b.WriteString("\n")
	}

	goModules := getSlice(deps, "go_modules")
	if len(goModules) > 0 {
		b.WriteString("### Key External Dependencies\n\n")
		b.WriteString("| Module | Version |\n")
		b.WriteString("|--------|---------|\n")
		notablePrefixes := NotableDependencyPrefixes
		for _, mod := range goModules {
			module := getStr(mod, "module", "")
			notable := false
			for _, prefix := range notablePrefixes {
				if strings.HasPrefix(module, prefix) {
					notable = true
					break
				}
			}
			if !notable {
				continue
			}
			b.WriteString(fmt.Sprintf("| %s | %s |\n",
				module, getStr(mod, "version", "")))
		}
		b.WriteString("\n")
	}
}

func renderServiceTable(b *strings.Builder, data map[string]interface{}) {
	services := getSlice(data, "services")
	if len(services) == 0 {
		b.WriteString("### Services\n\nNo services defined.\n\n")
		return
	}
	b.WriteString("### Services\n\n")
	b.WriteString("| Name | Type | Ports | Source |\n")
	b.WriteString("|------|------|-------|--------|\n")
	for _, svc := range services {
		ports := getSlice(svc, "ports")
		var portParts []string
		for _, p := range ports {
			portParts = append(portParts, fmt.Sprintf("%v/%s",
				p["port"], getStr(p, "protocol", "TCP")))
		}
		b.WriteString(fmt.Sprintf("| %s | %s | %s | `%s` |\n",
			getStr(svc, "name", ""), getStr(svc, "type", "ClusterIP"),
			strings.Join(portParts, ", "), getStr(svc, "source", "")))
	}
	b.WriteString("\n")
}

func renderIngressTable(b *strings.Builder, data map[string]interface{}) {
	ingress := getSlice(data, "ingress_routing")
	if len(ingress) == 0 {
		return
	}
	b.WriteString("### Ingress / Routing\n\n")
	b.WriteString("| Kind | Name | Hosts | Paths | TLS | Source |\n")
	b.WriteString("|------|------|-------|-------|-----|--------|\n")
	for _, res := range ingress {
		hosts := getStringSlice(res, "hosts")
		paths := getStringSlice(res, "paths")
		tls := getBool(res, "tls", false)
		tlsStr := "no"
		if tls {
			tlsStr = "yes"
		}
		b.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | `%s` |\n",
			getStr(res, "kind", ""), getStr(res, "name", ""),
			strings.Join(hosts, ", "), strings.Join(paths, ", "),
			tlsStr, getStr(res, "source", "")))
	}
	b.WriteString("\n")
}

func renderNetworkPolicyTable(b *strings.Builder, data map[string]interface{}) {
	netpols := getSlice(data, "network_policies")
	if len(netpols) == 0 {
		return
	}
	b.WriteString("### Network Policies\n\n")
	b.WriteString("| Name | Policy Types | Source |\n")
	b.WriteString("|------|-------------|--------|\n")
	for _, np := range netpols {
		policyTypes := getStringSlice(np, "policy_types")
		b.WriteString(fmt.Sprintf("| %s | %s | `%s` |\n",
			getStr(np, "name", ""), strings.Join(policyTypes, ", "),
			getStr(np, "source", "")))
	}
	b.WriteString("\n")
}

func renderRBACSection(b *strings.Builder, data map[string]interface{}) {
	rbac := getMap(data, "rbac")
	if rbac == nil {
		return
	}

	clusterRoles := getSlice(rbac, "cluster_roles")
	if len(clusterRoles) > 0 {
		b.WriteString("### Cluster Roles\n\n")
		b.WriteString("| Name | Resources | Verbs | Source |\n")
		b.WriteString("|------|-----------|-------|--------|\n")
		for _, role := range clusterRoles {
			roleName := getStr(role, "name", "")
			rules := getSlice(role, "rules")
			for _, rule := range rules {
				resources := getStringSlice(rule, "resources")
				verbs := getStringSlice(rule, "verbs")
				b.WriteString(fmt.Sprintf("| %s | %s | %s | `%s` |\n",
					roleName, strings.Join(resources, ", "),
					strings.Join(verbs, ", "), getStr(role, "source", "")))
			}
		}
		b.WriteString("\n")
	}

	markers := getSlice(rbac, "kubebuilder_markers")
	if len(markers) > 0 {
		b.WriteString("### Kubebuilder RBAC Markers\n\n")
		b.WriteString(fmt.Sprintf("%d markers found in source code.\n\n", len(markers)))
	}
}

func renderSecretTable(b *strings.Builder, data map[string]interface{}) {
	secrets := getSlice(data, "secrets_referenced")
	if len(secrets) == 0 {
		return
	}
	b.WriteString("### Secrets Referenced\n\n")
	b.WriteString("| Name | Type | Referenced By |\n")
	b.WriteString("|------|------|---------------|\n")
	for _, s := range secrets {
		refs := getStringSlice(s, "referenced_by")
		b.WriteString(fmt.Sprintf("| %s | %s | %s |\n",
			getStr(s, "name", ""), getStr(s, "type", ""),
			strings.Join(refs, ", ")))
	}
	b.WriteString("\n")
}

func renderSecurityContextSection(b *strings.Builder, data map[string]interface{}) {
	deployments := getSlice(data, "deployments")
	if len(deployments) == 0 {
		return
	}
	b.WriteString("### Container Security Contexts\n\n")
	b.WriteString("| Deployment | Container | RunAsNonRoot | ReadOnlyFS | Privileged | Source |\n")
	b.WriteString("|------------|-----------|--------------|------------|------------|--------|\n")
	for _, dep := range deployments {
		depName := getStr(dep, "name", "")
		source := getStr(dep, "source", "")
		containers := getSlice(dep, "containers")
		for _, c := range containers {
			cName := getStr(c, "name", "")
			sc := getMap(c, "security_context")
			runAsNonRoot := "?"
			readOnlyFS := "?"
			privileged := "?"
			if sc != nil {
				if v, ok := sc["runAsNonRoot"]; ok {
					runAsNonRoot = fmt.Sprintf("%v", v)
				}
				if v, ok := sc["readOnlyRootFilesystem"]; ok {
					readOnlyFS = fmt.Sprintf("%v", v)
				}
				if v, ok := sc["privileged"]; ok {
					privileged = fmt.Sprintf("%v", v)
				}
			}
			b.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | `%s` |\n",
				depName, cName, runAsNonRoot, readOnlyFS, privileged, source))
		}
	}
	b.WriteString("\n")
}

func renderConfigMapTable(b *strings.Builder, data map[string]interface{}) {
	configmaps := getSlice(data, "configmaps")
	if len(configmaps) == 0 {
		return
	}
	b.WriteString("### ConfigMaps\n\n")
	b.WriteString("| Name | Data Keys | Source |\n")
	b.WriteString("|------|-----------|--------|\n")
	for _, cm := range configmaps {
		keys := getStringSlice(cm, "data_keys")
		b.WriteString(fmt.Sprintf("| %s | %s | `%s` |\n",
			getStr(cm, "name", ""), strings.Join(keys, ", "),
			getStr(cm, "source", "")))
	}
	b.WriteString("\n")
}

func renderHelmSection(b *strings.Builder, data map[string]interface{}) {
	helm := getMap(data, "helm")
	if helm == nil {
		return
	}
	chartName := getStr(helm, "chart_name", "")
	if chartName == "" {
		return
	}
	b.WriteString("### Helm\n\n")
	b.WriteString(fmt.Sprintf("**Chart:** %s v%s\n\n",
		chartName, getStr(helm, "chart_version", "")))
}

func renderDockerfileTable(b *strings.Builder, data map[string]interface{}) {
	dockerfiles := getSlice(data, "dockerfiles")
	if len(dockerfiles) == 0 {
		b.WriteString("No Dockerfiles found.\n\n")
		return
	}
	b.WriteString("| Path | Base Image | Stages | User | Ports | Architectures | FIPS | Issues |\n")
	b.WriteString("|------|------------|--------|------|-------|---------------|------|--------|\n")
	for _, df := range dockerfiles {
		ports := make([]string, 0)
		if portSlice := getSlice(df, "exposed_ports"); len(portSlice) > 0 {
			for _, p := range portSlice {
				ports = append(ports, fmt.Sprintf("%v", p))
			}
		}
		archs := getStringSlice(df, "architectures")
		fips := getBool(df, "fips_enabled", false)
		fipsStr := ""
		if fips {
			fipsStr = "yes"
		}
		issues := getStringSlice(df, "issues")
		b.WriteString(fmt.Sprintf("| `%s` | %s | %d | %s | %s | %s | %s | %s |\n",
			getStr(df, "path", ""), getStr(df, "base_image", ""),
			getInt(df, "stages"), getStr(df, "user", ""),
			strings.Join(ports, ", "), strings.Join(archs, ", "),
			fipsStr, strings.Join(issues, "; ")))
	}
	b.WriteString("\n")
}

func renderControllerWatchTable(b *strings.Builder, data map[string]interface{}) {
	watches := getSlice(data, "controller_watches")
	if len(watches) == 0 {
		b.WriteString("No controller watches found.\n\n")
		return
	}

	// Group by type
	typeOrder := []string{"For", "Owns", "Watches"}
	grouped := make(map[string][]map[string]interface{})
	for _, w := range watches {
		t := getStr(w, "type", "")
		grouped[t] = append(grouped[t], w)
	}

	b.WriteString("| Type | GVK | Source |\n")
	b.WriteString("|------|-----|--------|\n")
	for _, t := range typeOrder {
		items := grouped[t]
		sort.Slice(items, func(i, j int) bool {
			return getStr(items[i], "gvk", "") < getStr(items[j], "gvk", "")
		})
		for _, w := range items {
			b.WriteString(fmt.Sprintf("| %s | %s | `%s` |\n",
				t, getStr(w, "gvk", ""), getStr(w, "source", "")))
		}
	}
	b.WriteString("\n")
}

func renderCacheSection(b *strings.Builder, data map[string]interface{}) {
	cache := getMap(data, "cache_config")
	if cache == nil {
		return
	}

	b.WriteString("## Cache Architecture\n\n")

	b.WriteString("### Manager Configuration\n\n")
	b.WriteString("| Property | Value |\n")
	b.WriteString("|----------|-------|\n")
	b.WriteString(fmt.Sprintf("| Manager file | `%s` |\n", getStr(cache, "manager_file", "")))
	b.WriteString(fmt.Sprintf("| Cache scope | %s |\n", getStr(cache, "cache_scope", "cluster-wide")))
	defaultTransform := getBool(cache, "default_transform", false)
	transformStr := "no"
	if defaultTransform {
		transformStr = "yes"
	}
	b.WriteString(fmt.Sprintf("| DefaultTransform | %s |\n", transformStr))
	goMemLimit := getStr(cache, "gomemlimit", "")
	if goMemLimit != "" {
		b.WriteString(fmt.Sprintf("| GOMEMLIMIT | %s |\n", goMemLimit))
	}
	memLimit := getStr(cache, "memory_limit", "")
	if memLimit != "" {
		b.WriteString(fmt.Sprintf("| Memory limit | %s |\n", memLimit))
	}
	b.WriteString("\n")

	filtered := getSlice(cache, "filtered_types")
	if len(filtered) > 0 {
		b.WriteString("### Filtered Types\n\n")
		b.WriteString("| Type | Filter Kind | Filter |\n")
		b.WriteString("|------|-------------|--------|\n")
		for _, ft := range filtered {
			b.WriteString(fmt.Sprintf("| %s | %s | %s |\n",
				getStr(ft, "type", ""), getStr(ft, "filter_kind", ""),
				getStr(ft, "filter", "")))
		}
		b.WriteString("\n")
	}

	disabled := getStringSlice(cache, "disabled_types")
	if len(disabled) > 0 {
		b.WriteString("### Cache-Bypassed Types (DisableFor)\n\n")
		for _, dt := range disabled {
			b.WriteString(fmt.Sprintf("- %s\n", dt))
		}
		b.WriteString("\n")
	}

	transforms := getStringSlice(cache, "transform_types")
	if len(transforms) > 0 {
		b.WriteString("### Transform-Stripped Types\n\n")
		for _, tt := range transforms {
			b.WriteString(fmt.Sprintf("- %s\n", tt))
		}
		b.WriteString("\n")
	}

	implicit := getSlice(cache, "implicit_informers")
	if len(implicit) > 0 {
		b.WriteString("### Implicit Informers (OOM Risk)\n\n")
		b.WriteString("| Type | Source | Risk |\n")
		b.WriteString("|------|--------|------|\n")
		for _, imp := range implicit {
			risk := getStr(imp, "risk", "")
			riskStr := risk
			if risk == "high" {
				riskStr = "**HIGH**"
			}
			b.WriteString(fmt.Sprintf("| %s | `%s` | %s |\n",
				getStr(imp, "type", ""), getStr(imp, "source", ""), riskStr))
		}
		b.WriteString("\n")
	}

	issues := getStringSlice(cache, "issues")
	if len(issues) > 0 {
		b.WriteString("### Issues\n\n")
		for _, issue := range issues {
			b.WriteString(fmt.Sprintf("- %s\n", issue))
		}
		b.WriteString("\n")
	}
}
