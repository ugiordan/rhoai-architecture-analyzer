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
		fields := sl.Properties["fields"]
		hasIsCA := strings.Contains(fields, "IsCA")
		hasDNSNames := strings.Contains(fields, "DNSNames") || strings.Contains(fields, "IPAddresses")

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
		typeName := sl.Properties["type"]
		if !strings.Contains(typeName, "ByObject") {
			continue
		}
		fields := sl.Properties["fields"]
		hasFilter := strings.Contains(fields, "Field") || strings.Contains(fields, "Label")
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
