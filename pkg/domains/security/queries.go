package security

import (
	"fmt"
	"strings"

	"github.com/ugiordan/architecture-analyzer/pkg/arch"
	"github.com/ugiordan/architecture-analyzer/pkg/graph"
	"github.com/ugiordan/architecture-analyzer/pkg/query"
)

func securityQueries() []query.Rule {
	return []query.Rule{
		{ID: "CGA-003", Name: "webhook-missing-update", Domain: "security", Severity: "high", Run: queryWebhookMissingUpdate},
		{ID: "CGA-004", Name: "rbac-precedence-bug", Domain: "security", Severity: "high", Run: queryRBACPrecedenceBug},
		{ID: "CGA-005", Name: "cert-as-ca", Domain: "security", Severity: "high", Run: queryCertAsCA},
		{ID: "CGA-006", Name: "cross-namespace-secret", Domain: "security", Severity: "high", Run: queryCrossNamespaceSecret},
		{ID: "CGA-007", Name: "unfiltered-cache", Domain: "security", Severity: "medium", Run: queryUnfilteredCache},
		{ID: "CGA-008", Name: "plaintext-secrets", Domain: "security", Severity: "medium", Run: queryPlaintextSecrets},
		{ID: "CGA-009", Name: "weak-serial-entropy", Domain: "security", Severity: "medium", Run: queryWeakSerialEntropy},
		{ID: "CGA-010", Name: "complexity-hotspot", Domain: "security", Severity: "medium", Run: queryComplexityHotspot},
		{ID: "CGA-011", Name: "untrusted-endpoint", Domain: "security", Severity: "informational", Run: queryUntrustedEndpoint},
		// Cross-domain queries: combine CPG code analysis with architecture deployment data.
		{ID: "CGA-012", Name: "unprotected-ingress", Domain: "security", Severity: "high", Run: queryUnprotectedIngress},
		{ID: "CGA-013", Name: "overprivileged-secret-access", Domain: "security", Severity: "medium", Run: queryOverprivilegedSecretAccess},
		{ID: "CGA-014", Name: "uncontrolled-egress", Domain: "security", Severity: "medium", Run: queryUncontrolledEgress},
	}
}

// webhookRef returns an ArchitectureRef string for webhooks, or empty if no arch data.
func webhookRef(g *graph.CPG) string {
	if g.ArchData == nil || len(g.ArchData.Webhooks) == 0 {
		return ""
	}
	var refs []string
	for _, wh := range g.ArchData.Webhooks {
		sources := make([]string, len(wh.Sources))
		for i, s := range wh.Sources {
			sources[i] = fmt.Sprintf("%s:%s", s.Type, s.File)
		}
		refs = append(refs, fmt.Sprintf("%s (path: %s, sources: [%s], failurePolicy: %s)", wh.Name, wh.Path, strings.Join(sources, ", "), wh.FailurePolicy))
	}
	return fmt.Sprintf("webhooks: %s", strings.Join(refs, "; "))
}

// secretsRef returns an ArchitectureRef string for secrets, or empty if no arch data.
func secretsRef(g *graph.CPG) string {
	if g.ArchData == nil || len(g.ArchData.Secrets) == 0 {
		return ""
	}
	var refs []string
	for _, s := range g.ArchData.Secrets {
		refs = append(refs, fmt.Sprintf("%s(%s, referenced by: %s)", s.Name, s.Type, strings.Join(s.ReferencedBy, ", ")))
	}
	return fmt.Sprintf("secrets: %s", strings.Join(refs, "; "))
}

// clusterRolesRef returns an ArchitectureRef string for RBAC cluster roles, or empty if no arch data.
func clusterRolesRef(g *graph.CPG) string {
	if g.ArchData == nil || len(g.ArchData.RBAC.ClusterRoles) == 0 {
		return ""
	}
	var refs []string
	for _, cr := range g.ArchData.RBAC.ClusterRoles {
		refs = append(refs, fmt.Sprintf("%s (source: %s)", cr.Name, cr.Source))
	}
	return fmt.Sprintf("cluster_roles: %s", strings.Join(refs, "; "))
}

// cacheIssuesRef returns an ArchitectureRef string for cache issues, or empty if no arch data.
func cacheIssuesRef(g *graph.CPG) string {
	if g.ArchData == nil || len(g.ArchData.Cache.Issues) == 0 {
		return ""
	}
	return fmt.Sprintf("cache_issues: %s", strings.Join(g.ArchData.Cache.Issues, "; "))
}

func queryWebhookMissingUpdate(g *graph.CPG) []query.Finding {
	var findings []query.Finding
	for _, fn := range g.NodesByKind(graph.NodeFunction) {
		if !fn.Annotations[AnnotHandlesAdmission] {
			continue
		}
		caseValues := fn.Properties["case_values"]
		if caseValues == "" {
			continue
		}
		hasCreate := strings.Contains(caseValues, "Create")
		hasDelete := strings.Contains(caseValues, "Delete")
		hasUpdate := strings.Contains(caseValues, "Update")

		if hasCreate && hasDelete && !hasUpdate {
			findings = append(findings, query.Finding{
				RuleID:           "CGA-003",
				Severity:         "high",
				Message:          fmt.Sprintf("Webhook %s handles CREATE/DELETE but not UPDATE (falls through to default allow)", fn.Name),
				File:             fn.File,
				Line:             fn.Line,
				NodeID:           fn.ID,
				ArchitectureRef:  webhookRef(g),
			})
		}
	}
	return findings
}

func queryRBACPrecedenceBug(g *graph.CPG) []query.Finding {
	var findings []query.Finding
	for _, fn := range g.NodesByKind(graph.NodeFunction) {
		if !fn.Annotations[AnnotCreatesRBAC] || !fn.Annotations[AnnotBindsSubject] {
			continue
		}
		findings = append(findings, query.Finding{
			RuleID:           "CGA-004",
			Severity:         "high",
			Message:          fmt.Sprintf("Function %s creates RBAC bindings with subject strings: check for operator precedence bugs", fn.Name),
			File:             fn.File,
			Line:             fn.Line,
			NodeID:           fn.ID,
			ArchitectureRef:  clusterRolesRef(g),
		})
	}
	return findings
}

func queryCertAsCA(g *graph.CPG) []query.Finding {
	var findings []query.Finding
	for _, sl := range g.NodesByKind(graph.NodeStructLiteral) {
		if !sl.Annotations[AnnotGeneratesCert] {
			continue
		}
		hasIsCA := containsField(sl.FieldNames, "IsCA")
		hasDNSNames := containsField(sl.FieldNames, "DNSNames") || containsField(sl.FieldNames, "IPAddresses")

		if hasIsCA && hasDNSNames {
			findings = append(findings, query.Finding{
				RuleID:   "CGA-005",
				Severity: "high",
				Message:  fmt.Sprintf("Certificate template at %s:%d has IsCA:true but also sets DNSNames/IPAddresses (server cert indicators)", sl.File, sl.Line),
				File:     sl.File,
				Line:     sl.Line,
				NodeID:   sl.ID,
			})
		}
	}
	return findings
}

func containsField(fields []string, target string) bool {
	for _, f := range fields {
		if f == target {
			return true
		}
	}
	return false
}

func queryCrossNamespaceSecret(g *graph.CPG) []query.Finding {
	var findings []query.Finding
	for _, fn := range g.NodesByKind(graph.NodeFunction) {
		if fn.Annotations[AnnotAccessesSecret] && fn.Annotations[AnnotCrossesNamespace] {
			findings = append(findings, query.Finding{
				RuleID:           "CGA-006",
				Severity:         "high",
				Message:          fmt.Sprintf("Function %s accesses secrets and crosses namespace boundary", fn.Name),
				File:             fn.File,
				Line:             fn.Line,
				NodeID:           fn.ID,
				ArchitectureRef:  secretsRef(g),
			})
		}
	}
	return findings
}

func queryUnfilteredCache(g *graph.CPG) []query.Finding {
	var findings []query.Finding
	for _, sl := range g.NodesByKind(graph.NodeStructLiteral) {
		if !sl.Annotations[AnnotConfiguresCache] {
			continue
		}
		if !strings.Contains(sl.StructType, "ByObject") {
			continue
		}
		hasFilter := containsField(sl.FieldNames, "Field") || containsField(sl.FieldNames, "Label")
		if !hasFilter {
			findings = append(findings, query.Finding{
				RuleID:           "CGA-007",
				Severity:         "medium",
				Message:          fmt.Sprintf("ByObject{} at %s:%d has no Field or Label selector (unfiltered cache may cause OOM)", sl.File, sl.Line),
				File:             sl.File,
				Line:             sl.Line,
				NodeID:           sl.ID,
				ArchitectureRef:  cacheIssuesRef(g),
			})
		}
	}
	return findings
}

func queryPlaintextSecrets(g *graph.CPG) []query.Finding {
	var findings []query.Finding
	for _, fn := range g.NodesByKind(graph.NodeFunction) {
		if fn.Annotations[AnnotWritesPlaintextSecret] {
			findings = append(findings, query.Finding{
				RuleID:           "CGA-008",
				Severity:         "medium",
				Message:          fmt.Sprintf("Function %s writes secret-like values to disk in plaintext", fn.Name),
				File:             fn.File,
				Line:             fn.Line,
				NodeID:           fn.ID,
				ArchitectureRef:  secretsRef(g),
			})
		}
	}
	return findings
}

func queryWeakSerialEntropy(g *graph.CPG) []query.Finding {
	var findings []query.Finding
	for _, fn := range g.NodesByKind(graph.NodeFunction) {
		if !fn.Annotations[AnnotGeneratesCert] {
			continue
		}
		for _, edge := range g.OutEdges(fn.ID) {
			target := g.GetNode(edge.To)
			if target == nil || target.Kind != graph.NodeCallSite {
				continue
			}
			if strings.Contains(target.Name, "rand.Int") {
				argTypes := target.Properties["arg_types"]
				stringArgs := target.Properties["string_args"]
				if strings.Contains(argTypes, "UnixNano") || strings.Contains(stringArgs, "UnixNano") {
					findings = append(findings, query.Finding{
						RuleID:   "CGA-009",
						Severity: "medium",
						Message:  fmt.Sprintf("Function %s generates cert with serial number bounded by time (weak entropy)", fn.Name),
						File:     fn.File,
						Line:     fn.Line,
						NodeID:   fn.ID,
					})
				}
			}
		}
	}
	return findings
}

// queryComplexityHotspot finds functions with complexity > 10 that have security-relevant annotations.
// Any annotation with the "sec:" prefix is considered security-relevant, so new annotations
// added to future annotators are automatically included.
func queryComplexityHotspot(g *graph.CPG) []query.Finding {
	var findings []query.Finding

	for _, fn := range g.NodesByKind(graph.NodeFunction) {
		if fn.Complexity <= 10 {
			continue
		}
		hasSecAnnotation := false
		for ann := range fn.Annotations {
			if strings.HasPrefix(ann, SecurityAnnotationPrefix) {
				hasSecAnnotation = true
				break
			}
		}
		if !hasSecAnnotation {
			continue
		}
		findings = append(findings, query.Finding{
			RuleID:   "CGA-010",
			Domain:   "security",
			Severity: "medium",
			Message:  fmt.Sprintf("Security-sensitive function %s has high cyclomatic complexity (%d), increasing review difficulty", fn.Name, fn.Complexity),
			File:     fn.File,
			Line:     fn.Line,
			NodeID:   fn.ID,
		})
	}
	return findings
}

// queryUntrustedEndpoint identifies HTTP endpoints classified as untrusted (public-facing).
// This is informational: it surfaces endpoints that receive unvalidated external input
// so they can be prioritized for manual review.
func queryUntrustedEndpoint(g *graph.CPG) []query.Finding {
	var findings []query.Finding
	for _, ep := range g.NodesByKind(graph.NodeHTTPEndpoint) {
		if ep.TrustLevel != graph.TrustUntrusted {
			continue
		}
		findings = append(findings, query.Finding{
			RuleID:   "CGA-011",
			Domain:   "security",
			Severity: "informational",
			Message:  fmt.Sprintf("Untrusted HTTP endpoint %s (%s %s) receives unvalidated external input", ep.Name, ep.HTTPMethod, ep.Route),
			File:     ep.File,
			Line:     ep.Line,
			NodeID:   ep.ID,
		})
	}
	return findings
}

// hasIngressNetworkPolicy checks whether the architecture data contains at least
// one NetworkPolicy with Ingress policy type and at least one ingress rule.
// An empty PodSelector ({}) means the policy applies to all pods in the namespace.
// A non-empty PolicyTypes list is required; policies with empty PolicyTypes are
// treated as Ingress-only per Kubernetes spec, but we require explicit declaration.
func hasIngressNetworkPolicy(g *graph.CPG) bool {
	if g.ArchData == nil {
		return false
	}
	for _, np := range g.ArchData.NetworkPolicies {
		if len(np.PolicyTypes) == 0 {
			// Per K8s spec, empty PolicyTypes defaults to Ingress-only.
			// An empty IngressRules means "deny all ingress" which is still a policy.
			return true
		}
		for _, pt := range np.PolicyTypes {
			// A NetworkPolicy with Ingress type and zero rules is a valid
			// "deny all ingress" policy per Kubernetes spec.
			if pt == "Ingress" {
				return true
			}
		}
	}
	return false
}

// hasEgressNetworkPolicy checks whether the architecture data contains at least
// one NetworkPolicy with Egress policy type. A policy with zero egress rules
// is a valid "deny all egress" policy per Kubernetes spec.
func hasEgressNetworkPolicy(g *graph.CPG) bool {
	if g.ArchData == nil {
		return false
	}
	for _, np := range g.ArchData.NetworkPolicies {
		for _, pt := range np.PolicyTypes {
			if pt == "Egress" {
				return true
			}
		}
	}
	return false
}

// networkPoliciesRef returns an ArchitectureRef string for network policies.
func networkPoliciesRef(g *graph.CPG) string {
	if g.ArchData == nil || len(g.ArchData.NetworkPolicies) == 0 {
		return "network_policies: none defined"
	}
	var refs []string
	for _, np := range g.ArchData.NetworkPolicies {
		refs = append(refs, fmt.Sprintf("%s (types: %s, source: %s)",
			np.Name, strings.Join(np.PolicyTypes, "/"), np.Source))
	}
	return fmt.Sprintf("network_policies: %s", strings.Join(refs, "; "))
}

// queryUnprotectedIngress finds HTTP endpoints that handle user input but have
// no NetworkPolicy restricting ingress traffic. Requires --with-arch data.
func queryUnprotectedIngress(g *graph.CPG) []query.Finding {
	if g.ArchData == nil {
		return nil
	}
	if hasIngressNetworkPolicy(g) {
		return nil
	}

	var findings []query.Finding
	for _, ep := range g.NodesByKind(graph.NodeHTTPEndpoint) {
		if ep.TrustLevel != graph.TrustUntrusted {
			continue
		}
		findings = append(findings, query.Finding{
			RuleID:          "CGA-012",
			Domain:          "security",
			Severity:        "high",
			Message:         fmt.Sprintf("Untrusted endpoint %s (%s %s) has no NetworkPolicy restricting ingress", ep.Name, ep.HTTPMethod, ep.Route),
			File:            ep.File,
			Line:            ep.Line,
			NodeID:          ep.ID,
			ArchitectureRef: networkPoliciesRef(g),
		})
	}
	return findings
}

// rbacHasWildcardSecretVerbs checks whether any RBAC ClusterRole grants wildcard
// verbs on secrets resources. Only matches rules whose APIGroups include the core
// group ("" or "*"), since secrets are a core API resource.
func rbacHasWildcardSecretVerbs(g *graph.CPG) bool {
	if g.ArchData == nil {
		return false
	}
	for _, cr := range g.ArchData.RBAC.ClusterRoles {
		for _, rule := range cr.Rules {
			if !ruleMatchesCoreGroup(rule) {
				continue
			}
			hasSecrets := false
			for _, r := range rule.Resources {
				if r == "secrets" || r == "*" {
					hasSecrets = true
					break
				}
			}
			if !hasSecrets {
				continue
			}
			for _, v := range rule.Verbs {
				if v == "*" {
					return true
				}
			}
		}
	}
	return false
}

// ruleMatchesCoreGroup returns true if the RBAC rule applies to the core API group
// (empty string) or uses a wildcard ("*"). Secrets are a core API resource.
func ruleMatchesCoreGroup(rule arch.RBACRule) bool {
	if len(rule.APIGroups) == 0 {
		return true // empty APIGroups list implies core group
	}
	for _, g := range rule.APIGroups {
		if g == "" || g == "*" {
			return true
		}
	}
	return false
}

// queryOverprivilegedSecretAccess finds functions that access secrets while RBAC
// grants wildcard verbs on secrets. Requires --with-arch data.
func queryOverprivilegedSecretAccess(g *graph.CPG) []query.Finding {
	if g.ArchData == nil {
		return nil
	}
	if !rbacHasWildcardSecretVerbs(g) {
		return nil
	}

	var findings []query.Finding
	for _, fn := range g.NodesByKind(graph.NodeFunction) {
		if !fn.Annotations[AnnotAccessesSecret] {
			continue
		}
		findings = append(findings, query.Finding{
			RuleID:          "CGA-013",
			Domain:          "security",
			Severity:        "medium",
			Message:         fmt.Sprintf("Function %s accesses secrets and RBAC grants wildcard verbs on secrets resource", fn.Name),
			File:            fn.File,
			Line:            fn.Line,
			NodeID:          fn.ID,
			ArchitectureRef: clusterRolesRef(g),
		})
	}
	return findings
}

// queryUncontrolledEgress finds functions that make external connections while
// no NetworkPolicy restricts egress traffic. Requires --with-arch data.
func queryUncontrolledEgress(g *graph.CPG) []query.Finding {
	if g.ArchData == nil {
		return nil
	}
	if hasEgressNetworkPolicy(g) {
		return nil
	}

	var findings []query.Finding
	reported := make(map[string]bool)

	// Check CPG nodes annotated as calling external services.
	for _, fn := range g.NodesByKind(graph.NodeFunction) {
		if !fn.Annotations[AnnotCallsExternal] {
			continue
		}
		key := fn.File + ":" + fn.Name
		reported[key] = true
		findings = append(findings, query.Finding{
			RuleID:          "CGA-014",
			Domain:          "security",
			Severity:        "medium",
			Message:         fmt.Sprintf("Function %s makes external connections with no NetworkPolicy restricting egress", fn.Name),
			File:            fn.File,
			Line:            fn.Line,
			NodeID:          fn.ID,
			ArchitectureRef: networkPoliciesRef(g),
		})
	}
	// Also check ExternalCall nodes from the CPG builder.
	// Skip if the parent function was already reported via annotation.
	for _, ec := range g.NodesByKind(graph.NodeExternalCall) {
		key := ec.File + ":" + ec.Name
		if reported[key] {
			continue
		}
		findings = append(findings, query.Finding{
			RuleID:          "CGA-014",
			Domain:          "security",
			Severity:        "medium",
			Message:         fmt.Sprintf("External call to %s with no NetworkPolicy restricting egress", ec.Name),
			File:            ec.File,
			Line:            ec.Line,
			NodeID:          ec.ID,
			ArchitectureRef: networkPoliciesRef(g),
		})
	}
	return findings
}
