package extractor

import (
	"fmt"
	"sort"
	"strings"
)

const (
	// minSelectorLength is the minimum length for a selector value to be
	// considered for prefix matching against deployment names.
	minSelectorLength = 3
	// minSelectorCoveragePercent is the minimum percentage of the deployment
	// name that a selector value must cover to be a valid prefix match.
	// Prevents false positives from short generic labels.
	minSelectorCoveragePercent = 40
)

// buildCrossReferences links related resources after all extraction is complete.
// It matches Services to their target Deployments via label selectors, detects
// runtime service dependencies from environment variable values, enriches
// external connections with credential sources, and infers networking
// capabilities from RBAC permissions.
func buildCrossReferences(arch *ComponentArchitecture) {
	linkServicesToDeployments(arch)
	detectRuntimeDependencies(arch)
	enrichExternalConnections(arch)
	inferNetworkingCapabilities(arch)
}

// linkServicesToDeployments matches each Service's selector labels against
// Deployment pod template labels (extracted from kustomize-rendered or raw YAML).
// For raw YAML deployments, we match against the deployment name since pod
// template labels typically include app=<deployment-name>.
func linkServicesToDeployments(arch *ComponentArchitecture) {
	for i := range arch.Services {
		svc := &arch.Services[i]
		if svc.Selector == nil || len(svc.Selector) == 0 {
			continue
		}

		for _, dep := range arch.Deployments {
			if selectorMatchesDeployment(svc.Selector, dep) {
				svc.TargetDeployment = dep.Name
				break
			}
		}
	}
}

// selectorMatchesDeployment checks if any selector label value matches the
// deployment name. This is platform-agnostic: it doesn't hardcode label keys
// like "app" or "control-plane". Instead it checks every label value against
// the deployment name, accounting for kustomize namePrefix/nameSuffix transforms
// and partial name matches (e.g., value "operator" matching deployment
// "operator-controller-manager").
func selectorMatchesDeployment(selector map[string]interface{}, dep Deployment) bool {
	for _, v := range selector {
		val, ok := v.(string)
		if !ok || val == "" || val == "true" || val == "false" {
			continue
		}
		if val == dep.Name || namesMatch(val, dep.Name) || namesMatch(dep.Name, val) {
			return true
		}
		// Check if the deployment name starts with the selector value.
		// Require the value covers at least 40% of the deployment name to avoid
		// false positives from short generic labels like "component: model-registry"
		// matching "model-registry-operator-controller-manager".
		if len(val) >= minSelectorLength && len(dep.Name) > 0 && strings.HasPrefix(dep.Name, val+"-") && len(val)*100/len(dep.Name) >= minSelectorCoveragePercent {
			return true
		}
	}
	return false
}

// detectRuntimeDependencies scans environment variables in deployments for
// references to other services in the same component (via Kubernetes DNS names
// like <service>.svc or <service>.<namespace>.svc.cluster.local).
func detectRuntimeDependencies(arch *ComponentArchitecture) {
	serviceNames := make(map[string]bool)
	for _, svc := range arch.Services {
		serviceNames[svc.Name] = true
	}

	for i := range arch.Deployments {
		dep := &arch.Deployments[i]
		for j := range dep.Containers {
			c := &dep.Containers[j]
			if c.EnvVars == nil {
				continue
			}
			for envName, envVal := range c.EnvVars {
				for svcName := range serviceNames {
					if containsServiceRef(envVal, svcName) {
						dep.Issues = appendUniqueStr(dep.Issues,
							fmt.Sprintf("runtime dependency: env %s references service %s", envName, svcName))
					}
				}
			}
		}
	}
}

// containsServiceRef checks if a value contains a Kubernetes service DNS reference.
func containsServiceRef(value, serviceName string) bool {
	patterns := []string{
		serviceName + ".svc",
		serviceName + ":",
		"://" + serviceName + "/",
		"://" + serviceName + ":",
	}
	lower := strings.ToLower(value)
	for _, p := range patterns {
		if strings.Contains(lower, strings.ToLower(p)) {
			return true
		}
	}
	return false
}

// credentialKeywords maps external connection types to keywords that identify
// related environment variables and secrets.
var credentialKeywords = map[string][]string{
	"database": {
		"database", "db_", "_db", "mysql", "postgres", "mariadb",
		"dsn", "sql", "jdbc", "connection_string",
	},
	"object-storage": {
		"s3_", "_s3", "minio", "bucket", "object_storage",
		"aws_access", "aws_secret", "storage_endpoint",
	},
	"grpc": {"grpc_", "_grpc"},
	"messaging": {"kafka", "nats", "amqp", "rabbitmq", "broker"},
}

// enrichExternalConnections cross-references external connections with
// deployment secrets and env vars to identify credential sources.
func enrichExternalConnections(arch *ComponentArchitecture) {
	for i := range arch.ExternalConnections {
		conn := &arch.ExternalConnections[i]
		keywords := credentialKeywords[conn.Type]
		if len(keywords) == 0 {
			continue
		}

		// Also match on the specific service name (e.g., "mysql", "minio")
		keywords = append(keywords, conn.Service)

		var sources []string

		for _, dep := range arch.Deployments {
			allContainers := append(dep.Containers, dep.InitContainers...)
			for _, c := range allContainers {
				// Check env vars referencing secrets
				for _, secretName := range c.EnvFromSecrets {
					if nameMatchesKeywords(secretName, keywords) {
						sources = appendUniqueStr(sources, "secret/"+secretName)
					}
				}
				// Check env vars referencing configmaps
				for _, cmName := range c.EnvFromConfigmaps {
					if nameMatchesKeywords(cmName, keywords) {
						sources = appendUniqueStr(sources, "configmap/"+cmName)
					}
				}
				// Check env var names/values
				for envName, envVal := range c.EnvVars {
					if nameMatchesKeywords(envName, keywords) {
						if envVal != "" && !strings.Contains(envVal, "://") {
							// Env var name matches but value isn't a URI: likely
							// a host, port, or credential reference
							sources = appendUniqueStr(sources, "env/"+envName)
						}
					}
				}
				// Check volume mounts with secret sources
				for _, vm := range c.VolumeMounts {
					if sn, ok := vm["secret_name"].(string); ok && sn != "" && nameMatchesKeywords(sn, keywords) {
						sources = appendUniqueStr(sources, "secret/"+sn)
					}
				}
			}
		}

		// Also check the Secrets list for names matching connection type
		for _, secret := range arch.Secrets {
			if nameMatchesKeywords(secret.Name, keywords) {
				sources = appendUniqueStr(sources, "secret/"+secret.Name)
			}
		}

		if len(sources) > 0 {
			conn.CredentialSources = sources
		}
	}
}

// nameMatchesKeywords checks if a resource name contains any of the keywords
// (case-insensitive).
func nameMatchesKeywords(name string, keywords []string) bool {
	lower := strings.ToLower(name)
	for _, kw := range keywords {
		if strings.Contains(lower, strings.ToLower(kw)) {
			return true
		}
	}
	return false
}

// appendUniqueStr appends s to slice only if not already present.
func appendUniqueStr(slice []string, s string) []string {
	for _, existing := range slice {
		if existing == s {
			return slice
		}
	}
	return append(slice, s)
}

// networkingResources maps RBAC apiGroup -> resource name (plural, lowercase) -> IngressResource Kind.
var networkingResources = map[string]map[string]string{
	"gateway.networking.k8s.io": {
		"httproutes": "HTTPRoute",
		"gateways":   "Gateway",
	},
	"networking.istio.io": {
		"virtualservices":  "VirtualService",
		"destinationrules": "DestinationRule",
	},
	"networking.k8s.io": {
		"ingresses": "Ingress",
	},
	"route.openshift.io": {
		"routes": "Route",
	},
}

// rbacMatch holds a single RBAC permission match for a networking resource kind.
type rbacMatch struct {
	kind     string
	verbs    []string
	roleName string
}

// kindInfo holds merged RBAC verbs and role name for a networking resource kind.
type kindInfo struct {
	verbs    []string
	roleName string
}

// inferNetworkingCapabilities scans RBAC cluster roles for networking-related
// API groups and either enriches existing IngressResource entries or creates
// new RBAC-inferred ones.
func inferNetworkingCapabilities(arch *ComponentArchitecture) {
	if arch.RBAC == nil {
		return
	}

	matches := collectNetworkingRBACMatches(arch.RBAC.ClusterRoles)
	if len(matches) == 0 {
		return
	}

	merged := mergeRBACMatchesByKind(matches)
	arch.IngressRouting = enrichIngressWithRBAC(arch.IngressRouting, merged)
}

// collectNetworkingRBACMatches scans RBAC cluster roles for networking-related
// API groups and returns all (Kind, verbs, roleName) tuples found.
func collectNetworkingRBACMatches(roles []RBACRole) []rbacMatch {
	var matches []rbacMatch
	for _, role := range roles {
		for _, rule := range role.Rules {
			for _, apiGroup := range rule.APIGroups {
				resourceMap, ok := networkingResources[apiGroup]
				if !ok {
					continue
				}
				hasWildcard := false
				for _, res := range rule.Resources {
					if res == "*" {
						hasWildcard = true
						break
					}
				}
				if hasWildcard {
					for _, kind := range resourceMap {
						matches = append(matches, rbacMatch{
							kind:     kind,
							verbs:    rule.Verbs,
							roleName: role.Name,
						})
					}
				} else {
					for _, res := range rule.Resources {
						resLower := strings.ToLower(res)
						if kind, ok := resourceMap[resLower]; ok {
							matches = append(matches, rbacMatch{
								kind:     kind,
								verbs:    rule.Verbs,
								roleName: role.Name,
							})
						}
					}
				}
			}
		}
	}
	return matches
}

// mergeRBACMatchesByKind deduplicates matches, merging verbs per Kind.
func mergeRBACMatchesByKind(matches []rbacMatch) map[string]*kindInfo {
	merged := make(map[string]*kindInfo)
	for _, m := range matches {
		ki, ok := merged[m.kind]
		if !ok {
			ki = &kindInfo{roleName: m.roleName}
			merged[m.kind] = ki
		}
		for _, v := range m.verbs {
			ki.verbs = appendUniqueStr(ki.verbs, v)
		}
	}
	return merged
}

// enrichIngressWithRBAC enriches existing IngressResource entries with RBAC
// verbs if they share the same Kind, or creates new rbac-inferred entries.
// Returns the updated slice. Keys are iterated in sorted order for
// deterministic output.
func enrichIngressWithRBAC(routing []IngressResource, merged map[string]*kindInfo) []IngressResource {
	kinds := make([]string, 0, len(merged))
	for k := range merged {
		kinds = append(kinds, k)
	}
	sort.Strings(kinds)

	for _, kind := range kinds {
		ki := merged[kind]
		enriched := false
		for i := range routing {
			if routing[i].Kind == kind {
				routing[i].RBACVerbs = ki.verbs
				routing[i].Note = "component has RBAC permissions to manage this resource type"
				enriched = true
				break
			}
		}
		if !enriched {
			routing = append(routing, IngressResource{
				Kind:      kind,
				Name:      "rbac-inferred",
				Source:    fmt.Sprintf("rbac/%s", ki.roleName),
				Note:      "component has RBAC permissions for this resource type but no static manifests found",
				RBACVerbs: ki.verbs,
			})
		}
	}
	return routing
}
