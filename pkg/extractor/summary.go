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

	// Component identity
	parts = append(parts, fmt.Sprintf("%s is a Kubernetes component", arch.Component))

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

	return strings.Join(parts, ". ") + "."
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
