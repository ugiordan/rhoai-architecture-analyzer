package extractor

import "sort"

// SortOutput normalizes all slices in a ComponentArchitecture for deterministic
// JSON output. Call after ExtractAll to ensure re-running produces identical
// results regardless of filesystem walk order or map iteration order.
func SortOutput(arch *ComponentArchitecture) {
	if arch == nil {
		return
	}

	// CRDs: sort by group/version/kind
	sort.Slice(arch.CRDs, func(i, j int) bool {
		a, b := arch.CRDs[i], arch.CRDs[j]
		if a.Group != b.Group {
			return a.Group < b.Group
		}
		if a.Version != b.Version {
			return a.Version < b.Version
		}
		return a.Kind < b.Kind
	})

	// RBAC
	if arch.RBAC != nil {
		sortRBACRoles(arch.RBAC.ClusterRoles)
		sortRBACBindings(arch.RBAC.ClusterRoleBindings)
		sortRBACRoles(arch.RBAC.Roles)
		sortRBACBindings(arch.RBAC.RoleBindings)
		sort.Slice(arch.RBAC.KubebuilderMarkers, func(i, j int) bool {
			a, b := arch.RBAC.KubebuilderMarkers[i], arch.RBAC.KubebuilderMarkers[j]
			if a.File != b.File {
				return a.File < b.File
			}
			return a.Line < b.Line
		})
	}

	// Services: sort by name
	sort.Slice(arch.Services, func(i, j int) bool {
		return arch.Services[i].Name < arch.Services[j].Name
	})

	// Deployments: sort by name
	sort.Slice(arch.Deployments, func(i, j int) bool {
		return arch.Deployments[i].Name < arch.Deployments[j].Name
	})

	// NetworkPolicies: sort by name
	sort.Slice(arch.NetworkPolicies, func(i, j int) bool {
		return arch.NetworkPolicies[i].Name < arch.NetworkPolicies[j].Name
	})

	// ControllerWatches: sort by type, then GVK
	sort.Slice(arch.ControllerWatch, func(i, j int) bool {
		a, b := arch.ControllerWatch[i], arch.ControllerWatch[j]
		if a.Type != b.Type {
			return a.Type < b.Type
		}
		return a.GVK < b.GVK
	})

	// Dependencies
	if arch.Dependencies != nil {
		sort.Slice(arch.Dependencies.GoModules, func(i, j int) bool {
			return arch.Dependencies.GoModules[i].Module < arch.Dependencies.GoModules[j].Module
		})
		sort.Slice(arch.Dependencies.InternalODH, func(i, j int) bool {
			return arch.Dependencies.InternalODH[i].Component < arch.Dependencies.InternalODH[j].Component
		})
		sort.Slice(arch.Dependencies.ReplaceDirectives, func(i, j int) bool {
			return arch.Dependencies.ReplaceDirectives[i].Original < arch.Dependencies.ReplaceDirectives[j].Original
		})
		sort.Strings(arch.Dependencies.Issues)
	}

	// Secrets: sort by name
	sort.Slice(arch.Secrets, func(i, j int) bool {
		return arch.Secrets[i].Name < arch.Secrets[j].Name
	})
	for idx := range arch.Secrets {
		sort.Strings(arch.Secrets[idx].ReferencedBy)
	}

	// Dockerfiles: sort by path
	sort.Slice(arch.Dockerfiles, func(i, j int) bool {
		return arch.Dockerfiles[i].Path < arch.Dockerfiles[j].Path
	})

	// Webhooks: sort by name
	sort.Slice(arch.Webhooks, func(i, j int) bool {
		return arch.Webhooks[i].Name < arch.Webhooks[j].Name
	})

	// ConfigMaps: sort by name
	sort.Slice(arch.ConfigMaps, func(i, j int) bool {
		return arch.ConfigMaps[i].Name < arch.ConfigMaps[j].Name
	})

	// HTTPEndpoints: sort by path, then method
	sort.Slice(arch.HTTPEndpoints, func(i, j int) bool {
		a, b := arch.HTTPEndpoints[i], arch.HTTPEndpoints[j]
		if a.Path != b.Path {
			return a.Path < b.Path
		}
		return a.Method < b.Method
	})

	// IngressRouting: sort by kind, then name
	sort.Slice(arch.IngressRouting, func(i, j int) bool {
		a, b := arch.IngressRouting[i], arch.IngressRouting[j]
		if a.Kind != b.Kind {
			return a.Kind < b.Kind
		}
		return a.Name < b.Name
	})

	// ExternalConnections: sort by type, then service, then source
	sort.Slice(arch.ExternalConnections, func(i, j int) bool {
		a, b := arch.ExternalConnections[i], arch.ExternalConnections[j]
		if a.Type != b.Type {
			return a.Type < b.Type
		}
		if a.Service != b.Service {
			return a.Service < b.Service
		}
		return a.Source < b.Source
	})

	// FeatureGates: sort by name
	sort.Slice(arch.FeatureGates, func(i, j int) bool {
		return arch.FeatureGates[i].Name < arch.FeatureGates[j].Name
	})

	// KustomizeComponents: already sorted by extractKustomizeComponents, but ensure consistency
	sort.Slice(arch.KustomizeComponents, func(i, j int) bool {
		return arch.KustomizeComponents[i].Name < arch.KustomizeComponents[j].Name
	})

	// ServingRuntimes: sort by name
	sort.Slice(arch.ServingRuntimes, func(i, j int) bool {
		return arch.ServingRuntimes[i].Name < arch.ServingRuntimes[j].Name
	})

	// ResourceDefaults: sort by component, then key
	sort.Slice(arch.ResourceDefaults, func(i, j int) bool {
		a, b := arch.ResourceDefaults[i], arch.ResourceDefaults[j]
		if a.Component != b.Component {
			return a.Component < b.Component
		}
		return a.Key < b.Key
	})

	// PodDisruptionBudgets: sort by name
	sort.Slice(arch.PodDisruptionBudgets, func(i, j int) bool {
		return arch.PodDisruptionBudgets[i].Name < arch.PodDisruptionBudgets[j].Name
	})

	// HorizontalPodAutoscalers: sort by name
	sort.Slice(arch.HorizontalPodAutoscalers, func(i, j int) bool {
		return arch.HorizontalPodAutoscalers[i].Name < arch.HorizontalPodAutoscalers[j].Name
	})

	// APITypes: sort by source (preserves file order), then name
	sort.Slice(arch.APITypes, func(i, j int) bool {
		a, b := arch.APITypes[i], arch.APITypes[j]
		if a.Source != b.Source {
			return a.Source < b.Source
		}
		return a.Name < b.Name
	})

	// CacheConfig: sort filtered types and implicit informers
	if arch.CacheConfig != nil {
		sort.Slice(arch.CacheConfig.FilteredTypes, func(i, j int) bool {
			return arch.CacheConfig.FilteredTypes[i].Type < arch.CacheConfig.FilteredTypes[j].Type
		})
		sort.Slice(arch.CacheConfig.TransformTypes, func(i, j int) bool {
			return arch.CacheConfig.TransformTypes[i].Type < arch.CacheConfig.TransformTypes[j].Type
		})
		sort.Strings(arch.CacheConfig.DisabledTypes)
		sort.Slice(arch.CacheConfig.ImplicitInformers, func(i, j int) bool {
			return arch.CacheConfig.ImplicitInformers[i].Type < arch.CacheConfig.ImplicitInformers[j].Type
		})
		sort.Strings(arch.CacheConfig.Issues)
	}

	// OperatorConfig: sort by source (file locality), then name
	sort.Slice(arch.OperatorConfig, func(i, j int) bool {
		if arch.OperatorConfig[i].Source != arch.OperatorConfig[j].Source {
			return arch.OperatorConfig[i].Source < arch.OperatorConfig[j].Source
		}
		return arch.OperatorConfig[i].Name < arch.OperatorConfig[j].Name
	})

	// ReconcileSequences: sort by controller name, steps by order
	sort.Slice(arch.ReconcileSequences, func(i, j int) bool {
		return arch.ReconcileSequences[i].Controller < arch.ReconcileSequences[j].Controller
	})
	for idx := range arch.ReconcileSequences {
		sort.Slice(arch.ReconcileSequences[idx].Steps, func(i, j int) bool {
			return arch.ReconcileSequences[idx].Steps[i].Order < arch.ReconcileSequences[idx].Steps[j].Order
		})
	}

	// PrometheusMetrics: sort by name
	sort.Slice(arch.PrometheusMetrics, func(i, j int) bool {
		return arch.PrometheusMetrics[i].Name < arch.PrometheusMetrics[j].Name
	})

	// StatusConditions: sort by type, then reasons within each
	sort.Slice(arch.StatusConditions, func(i, j int) bool {
		return arch.StatusConditions[i].Type < arch.StatusConditions[j].Type
	})
	for idx := range arch.StatusConditions {
		sort.Strings(arch.StatusConditions[idx].Reasons)
	}

	// PlatformDetection: sort nested slices
	if arch.PlatformDetection != nil {
		sort.Slice(arch.PlatformDetection.Capabilities, func(i, j int) bool {
			return arch.PlatformDetection.Capabilities[i].Name < arch.PlatformDetection.Capabilities[j].Name
		})
		sort.Slice(arch.PlatformDetection.Conditionals, func(i, j int) bool {
			if arch.PlatformDetection.Conditionals[i].Condition != arch.PlatformDetection.Conditionals[j].Condition {
				return arch.PlatformDetection.Conditionals[i].Condition < arch.PlatformDetection.Conditionals[j].Condition
			}
			return arch.PlatformDetection.Conditionals[i].Source < arch.PlatformDetection.Conditionals[j].Source
		})
	}
}

func sortRBACRoles(roles []RBACRole) {
	sort.Slice(roles, func(i, j int) bool {
		return roles[i].Name < roles[j].Name
	})
}

func sortRBACBindings(bindings []RBACBinding) {
	sort.Slice(bindings, func(i, j int) bool {
		return bindings[i].Name < bindings[j].Name
	})
}
