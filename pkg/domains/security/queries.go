package security

import (
	"fmt"
	"strings"

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
	}
}

// webhookRef returns an ArchitectureRef string for webhooks, or empty if no arch data.
func webhookRef(g *graph.CPG) string {
	if g.ArchData == nil || len(g.ArchData.Webhooks) == 0 {
		return ""
	}
	var refs []string
	for _, wh := range g.ArchData.Webhooks {
		refs = append(refs, fmt.Sprintf("%s (path: %s, source: %s, failurePolicy: %s)", wh.Name, wh.Path, wh.Source, wh.FailurePolicy))
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
