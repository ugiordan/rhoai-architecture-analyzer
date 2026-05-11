package extractor

import (
	"fmt"
	"sort"
	"strings"
)

// generateSummary produces a natural-language summary of the component architecture.
// This helps LLMs understand the component at a glance without parsing the full JSON.
func generateSummary(arch *ComponentArchitecture) string {
	var parts []string

	// Component identity with aliases
	if len(arch.Aliases) > 0 {
		parts = append(parts, fmt.Sprintf("%s (also known as: %s) is a Kubernetes component",
			arch.Component, strings.Join(arch.Aliases, ", ")))
	} else {
		parts = append(parts, fmt.Sprintf("%s is a Kubernetes component", arch.Component))
	}

	// CRDs
	if len(arch.CRDs) > 0 {
		kinds := make([]string, 0, len(arch.CRDs))
		seen := make(map[string]bool)
		for _, crd := range arch.CRDs {
			if !seen[crd.Kind] {
				kinds = append(kinds, crd.Kind)
				seen[crd.Kind] = true
			}
		}
		parts = append(parts, fmt.Sprintf("managing %d CRD(s): %s", len(kinds), strings.Join(kinds, ", ")))
	}

	// Deployments
	if len(arch.Deployments) > 0 {
		names := make([]string, 0, len(arch.Deployments))
		totalContainers := 0
		for _, dep := range arch.Deployments {
			names = append(names, dep.Name)
			totalContainers += len(dep.Containers)
		}
		parts = append(parts, fmt.Sprintf("deploying %d workload(s) (%s) with %d container(s)",
			len(arch.Deployments), strings.Join(names, ", "), totalContainers))
	}

	// Controller watches
	if len(arch.ControllerWatch) > 0 {
		watchTypes := make(map[string]bool)
		for _, w := range arch.ControllerWatch {
			watchTypes[w.GVK] = true
		}
		parts = append(parts, fmt.Sprintf("watching %d resource type(s)", len(watchTypes)))
	}

	// Services
	if len(arch.Services) > 0 {
		parts = append(parts, fmt.Sprintf("exposing %d service(s)", len(arch.Services)))
	}

	// Webhooks
	if len(arch.Webhooks) > 0 {
		validating := 0
		mutating := 0
		for _, wh := range arch.Webhooks {
			switch wh.Type {
			case "validating":
				validating++
			case "mutating":
				mutating++
			}
		}
		whParts := []string{}
		if validating > 0 {
			whParts = append(whParts, fmt.Sprintf("%d validating", validating))
		}
		if mutating > 0 {
			whParts = append(whParts, fmt.Sprintf("%d mutating", mutating))
		}
		parts = append(parts, fmt.Sprintf("with %d webhook(s) (%s)", len(arch.Webhooks), strings.Join(whParts, ", ")))
	}

	// Network policies
	if len(arch.NetworkPolicies) > 0 {
		parts = append(parts, fmt.Sprintf("protected by %d network policy/policies", len(arch.NetworkPolicies)))
	}

	// PDB/HPA
	if len(arch.PodDisruptionBudgets) > 0 {
		parts = append(parts, fmt.Sprintf("with %d PodDisruptionBudget(s)", len(arch.PodDisruptionBudgets)))
	}
	if len(arch.HorizontalPodAutoscalers) > 0 {
		parts = append(parts, fmt.Sprintf("with %d HorizontalPodAutoscaler(s)", len(arch.HorizontalPodAutoscalers)))
	}

	// External connections
	if len(arch.ExternalConnections) > 0 {
		types := make(map[string]bool)
		for _, ec := range arch.ExternalConnections {
			types[ec.Service] = true
		}
		services := make([]string, 0, len(types))
		for s := range types {
			services = append(services, s)
		}
		sort.Strings(services)
		parts = append(parts, fmt.Sprintf("connecting to external services: %s", strings.Join(services, ", ")))
	}

	// Runtime dependencies
	if len(arch.RuntimeDependencies) > 0 {
		depTypes := make(map[string]bool)
		for _, rd := range arch.RuntimeDependencies {
			depTypes[rd.Name] = true
		}
		depNames := make([]string, 0, len(depTypes))
		for n := range depTypes {
			depNames = append(depNames, n)
		}
		sort.Strings(depNames)
		parts = append(parts, fmt.Sprintf("runtime dependencies: %s", strings.Join(depNames, ", ")))
	}

	// API types
	if len(arch.APITypes) > 0 {
		specCount := 0
		totalFields := 0
		for _, t := range arch.APITypes {
			if t.IsSpec {
				specCount++
			}
			totalFields += len(t.Fields)
		}
		if specCount > 0 {
			parts = append(parts, fmt.Sprintf("defining %d API spec type(s) across %d struct(s) with %d total fields",
				specCount, len(arch.APITypes), totalFields))
		} else {
			parts = append(parts, fmt.Sprintf("defining %d API type struct(s) with %d total fields",
				len(arch.APITypes), totalFields))
		}
	}

	// Operator config
	if len(arch.OperatorConfig) > 0 {
		imageCount := countConstantsByCategory(arch.OperatorConfig, "image")
		envCount := countConstantsByCategory(arch.OperatorConfig, "env_var")
		configParts := []string{}
		if imageCount > 0 {
			configParts = append(configParts, fmt.Sprintf("%d image references", imageCount))
		}
		if envCount > 0 {
			configParts = append(configParts, fmt.Sprintf("%d env vars", envCount))
		}
		if len(configParts) > 0 {
			parts = append(parts, fmt.Sprintf("with %d operator config constants (%s)",
				len(arch.OperatorConfig), strings.Join(configParts, ", ")))
		} else {
			parts = append(parts, fmt.Sprintf("with %d operator config constants", len(arch.OperatorConfig)))
		}
	}

	// Reconcile sequences
	if len(arch.ReconcileSequences) > 0 {
		totalSteps := countTotalReconcileSteps(arch.ReconcileSequences)
		parts = append(parts, fmt.Sprintf("following %d reconciliation steps across %d controller(s)",
			totalSteps, len(arch.ReconcileSequences)))
	}

	// Prometheus metrics
	if len(arch.PrometheusMetrics) > 0 {
		parts = append(parts, fmt.Sprintf("exposing %d Prometheus metrics", len(arch.PrometheusMetrics)))
	}

	// Status conditions
	if len(arch.StatusConditions) > 0 {
		parts = append(parts, fmt.Sprintf("reporting %d status condition types", len(arch.StatusConditions)))
	}

	// Platform detection
	if arch.PlatformDetection != nil && len(arch.PlatformDetection.Conditionals) > 0 {
		parts = append(parts, fmt.Sprintf("with %d platform-conditional resource paths",
			len(arch.PlatformDetection.Conditionals)))
	}

	// Ingress
	if len(arch.IngressRouting) > 0 {
		parts = append(parts, fmt.Sprintf("with %d ingress route(s)", len(arch.IngressRouting)))
	}

	// Python k8s API calls
	if len(arch.PythonK8sCalls) > 0 {
		kindRefs := 0
		apiCalls := 0
		for _, c := range arch.PythonK8sCalls {
			if c.Operation == "kind_ref" {
				kindRefs++
			} else if c.Operation != "import" {
				apiCalls++
			}
		}
		if apiCalls > 0 || kindRefs > 0 {
			callParts := []string{}
			if apiCalls > 0 {
				callParts = append(callParts, fmt.Sprintf("%d API call(s)", apiCalls))
			}
			if kindRefs > 0 {
				callParts = append(callParts, fmt.Sprintf("%d CRD kind reference(s)", kindRefs))
			}
			parts = append(parts, fmt.Sprintf("Python k8s client with %s", strings.Join(callParts, " and ")))
		}
	}

	// Label contracts: surface cross-component integration signals
	if len(arch.LabelContracts) > 0 {
		integrations := make(map[string]bool)
		for _, lc := range arch.LabelContracts {
			integrations[lc.Integration] = true
		}
		intNames := make([]string, 0, len(integrations))
		for name := range integrations {
			intNames = append(intNames, name)
		}
		sort.Strings(intNames)
		parts = append(parts, fmt.Sprintf("integrates with %s via labels/annotations", strings.Join(intNames, ", ")))
	}

	// Sidecar/proxy pattern detection from deployments
	sidecarPatterns := detectSidecarPatterns(arch)
	if len(sidecarPatterns) > 0 {
		parts = append(parts, strings.Join(sidecarPatterns, ". "))
	}

	// RBAC binding creation patterns: detect when this component creates bindings for other components
	rbacPatterns := detectRBACBindingPatterns(arch)
	if len(rbacPatterns) > 0 {
		parts = append(parts, strings.Join(rbacPatterns, ". "))
	}

	// Kustomize overlay structure
	if len(arch.KustomizeOverlayRefs) > 0 {
		overlayNames := make([]string, 0, len(arch.KustomizeOverlayRefs))
		for _, ref := range arch.KustomizeOverlayRefs {
			overlayNames = append(overlayNames, ref.Overlay)
		}
		parts = append(parts, fmt.Sprintf("kustomize overlays: %s", strings.Join(overlayNames, ", ")))
	}

	// Component references (provider/adapter patterns)
	if len(arch.ComponentRefs) > 0 {
		targets := make(map[string][]string) // target -> types
		for _, ref := range arch.ComponentRefs {
			targets[ref.Target] = append(targets[ref.Target], ref.Type)
		}
		refParts := make([]string, 0, len(targets))
		for target, types := range targets {
			refParts = append(refParts, target+" ("+strings.Join(types, ", ")+")")
		}
		sort.Strings(refParts)
		parts = append(parts, fmt.Sprintf("references components: %s", strings.Join(refParts, ", ")))
	}

	// ConfigMap volume mounts
	if len(arch.ConfigMapVolumes) > 0 {
		cmNames := make(map[string]bool)
		for _, vol := range arch.ConfigMapVolumes {
			cmNames[vol.ConfigMapName] = true
		}
		names := make([]string, 0, len(cmNames))
		for n := range cmNames {
			names = append(names, n)
		}
		sort.Strings(names)
		parts = append(parts, fmt.Sprintf("mounts %d ConfigMap(s) as volumes: %s", len(cmNames), strings.Join(names, ", ")))
	}

	return strings.Join(parts, ". ") + "."
}

// detectSidecarPatterns scans deployments for well-known sidecar containers
// (kube-rbac-proxy, kube-auth-proxy, envoy, oauth-proxy) and reports them
// as interaction patterns in the summary.
func detectSidecarPatterns(arch *ComponentArchitecture) []string {
	knownSidecars := map[string]string{
		"kube-rbac-proxy": "deploys kube-rbac-proxy sidecar for RBAC-based authorization",
		"kube-auth-proxy": "deploys kube-auth-proxy for authentication",
		"oauth-proxy":     "deploys OAuth proxy sidecar for authentication",
		"envoy":           "deploys Envoy sidecar for traffic management",
		"istio-proxy":     "uses Istio sidecar injection",
	}

	found := make(map[string]bool)
	for _, dep := range arch.Deployments {
		for _, c := range dep.Containers {
			for prefix, desc := range knownSidecars {
				if strings.Contains(strings.ToLower(c.Name), prefix) ||
					strings.Contains(strings.ToLower(c.Image), prefix) {
					if !found[prefix] {
						found[prefix] = true
						_ = desc // used below
					}
				}
			}
		}
	}

	// Also check reconciliation sequences for sidecar deployment steps
	for _, seq := range arch.ReconcileSequences {
		for _, step := range seq.Steps {
			methodLower := strings.ToLower(step.Method)
			for prefix, desc := range knownSidecars {
				if strings.Contains(methodLower, strings.ReplaceAll(prefix, "-", "")) {
					if !found[prefix] {
						found[prefix] = true
						_ = desc
					}
				}
			}
		}
	}

	// Also check operator config constants for sidecar image references
	for _, oc := range arch.OperatorConfig {
		if oc.Category != "image" {
			continue
		}
		valLower := strings.ToLower(oc.Value)
		for prefix, desc := range knownSidecars {
			if strings.Contains(valLower, prefix) {
				if !found[prefix] {
					found[prefix] = true
					_ = desc
				}
			}
		}
	}

	var patterns []string
	for prefix := range found {
		patterns = append(patterns, knownSidecars[prefix])
	}
	sort.Strings(patterns)
	return patterns
}

// detectRBACBindingPatterns looks for RBAC bindings that reference external
// ClusterRoles (roles not defined in this component), which indicate
// cross-component RBAC integration.
func detectRBACBindingPatterns(arch *ComponentArchitecture) []string {
	if arch.RBAC == nil {
		return nil
	}

	// Collect roles defined by this component
	ownRoles := make(map[string]bool)
	for _, r := range arch.RBAC.ClusterRoles {
		ownRoles[r.Name] = true
	}
	for _, r := range arch.RBAC.Roles {
		ownRoles[r.Name] = true
	}

	// Find bindings referencing external roles
	var external []string
	seen := make(map[string]bool)
	allBindings := append(arch.RBAC.ClusterRoleBindings, arch.RBAC.RoleBindings...)
	for _, b := range allBindings {
		if b.RoleRef != "" && !ownRoles[b.RoleRef] && !seen[b.RoleRef] {
			seen[b.RoleRef] = true
			external = append(external, b.RoleRef)
		}
	}

	if len(external) == 0 {
		return nil
	}

	sort.Strings(external)
	return []string{fmt.Sprintf("creates RBAC bindings to external roles: %s", strings.Join(external, ", "))}
}

// countConstantsByCategory counts constants matching a specific category.
func countConstantsByCategory(constants []OperatorConstant, category string) int {
	count := 0
	for _, c := range constants {
		if c.Category == category {
			count++
		}
	}
	return count
}
