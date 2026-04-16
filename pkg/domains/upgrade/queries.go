package upgrade

import (
	"fmt"
	"strings"

	"github.com/ugiordan/rhoai-architecture-analyzer/pkg/graph"
	"github.com/ugiordan/rhoai-architecture-analyzer/pkg/query"
)

func upgradeQueries() []query.Rule {
	return []query.Rule{
		{ID: "CGA-U01", Name: "unconverted-crd", Domain: "upgrade", Severity: "medium", Run: queryUnconvertedCRD},
		{ID: "CGA-U02", Name: "deprecated-api-usage", Domain: "upgrade", Severity: "low", Run: queryDeprecatedAPIUsage},
		{ID: "CGA-U03", Name: "ungated-feature", Domain: "upgrade", Severity: "medium", Run: queryUngatedFeature},
		{ID: "CGA-U04", Name: "unchecked-version-access", Domain: "upgrade", Severity: "high", Run: queryUncheckedVersionAccess},
	}
}

// CGA-U01: CRD types with multiple versions but no conversion functions.
// Requires --with-arch to provide CRD inventory from architecture extraction.
func queryUnconvertedCRD(g *graph.CPG) []query.Finding {
	if g.ArchData == nil {
		return nil
	}

	// Group CRDs by group+kind to find multi-version types
	type crdKey struct{ group, kind string }
	crdVersions := make(map[crdKey][]string)
	crdSource := make(map[crdKey]string)
	for _, crd := range g.ArchData.CRDs {
		key := crdKey{crd.Group, crd.Kind}
		crdVersions[key] = append(crdVersions[key], crd.Version)
		if crdSource[key] == "" {
			crdSource[key] = crd.Source
		}
	}

	// Build set of CRD kinds that have conversion functions.
	hasConversion := make(map[string]bool)
	for _, fn := range g.NodesByKind(graph.NodeFunction) {
		if !fn.Annotations[AnnotVersionConversion] {
			continue
		}
		receiver := fn.Properties["receiver"]
		typeName := extractReceiverType(receiver)
		if typeName != "" {
			hasConversion[typeName] = true
		}
	}

	var findings []query.Finding
	for key, versions := range crdVersions {
		if len(versions) < 2 {
			continue
		}
		if hasConversion[key.kind] {
			continue
		}
		findings = append(findings, query.Finding{
			RuleID:          "CGA-U01",
			Severity:        "medium",
			Message:         fmt.Sprintf("CRD %s.%s has versions %v but no conversion webhook (ConvertTo/ConvertFrom) found in code", key.kind, key.group, versions),
			ArchitectureRef: fmt.Sprintf("crd_source: %s", crdSource[key]),
		})
	}
	return findings
}

// extractReceiverType parses the type name from a Go receiver string.
// Input examples: "(w *Widget)", "(r Widget)", "(w *v1alpha1.Widget)".
// Returns just the unqualified type name, e.g., "Widget".
func extractReceiverType(receiver string) string {
	receiver = strings.TrimSpace(receiver)
	receiver = strings.Trim(receiver, "()")
	parts := strings.Fields(receiver)
	if len(parts) == 0 {
		return ""
	}
	typeName := parts[len(parts)-1]
	typeName = strings.TrimPrefix(typeName, "*")
	if idx := strings.LastIndex(typeName, "."); idx >= 0 {
		typeName = typeName[idx+1:]
	}
	return typeName
}

// CGA-U02: Deprecated API usage in non-test files
func queryDeprecatedAPIUsage(g *graph.CPG) []query.Finding {
	var findings []query.Finding
	for _, cs := range g.NodesByKind(graph.NodeCallSite) {
		if !cs.Annotations[AnnotDeprecatedAPI] {
			continue
		}
		if strings.HasSuffix(cs.File, "_test.go") || strings.Contains(cs.File, "vendor/") {
			continue
		}
		findings = append(findings, query.Finding{
			RuleID:   "CGA-U02",
			Severity: "low",
			Message:  fmt.Sprintf("Deprecated API usage: %s at %s:%d", cs.Name, cs.File, cs.Line),
			File:     cs.File,
			Line:     cs.Line,
			NodeID:   cs.ID,
		})
	}
	return findings
}

// CGA-U03: Functions without feature gate checks.
// TODO: Implement once architecture data provides feature gate inventory.
// Requires matching function names against feature-gated component list from arch data.
func queryUngatedFeature(g *graph.CPG) []query.Finding {
	return nil
}

// CGA-U04: Version check functions that may access slices without bounds check
func queryUncheckedVersionAccess(g *graph.CPG) []query.Finding {
	var findings []query.Finding
	for _, cs := range g.NodesByKind(graph.NodeCallSite) {
		if !cs.Annotations[AnnotVersionCheck] {
			continue
		}
		for _, edge := range g.InEdges(cs.ID) {
			fn := g.GetNode(edge.From)
			if fn == nil || fn.Kind != graph.NodeFunction {
				continue
			}
			findings = append(findings, query.Finding{
				RuleID:   "CGA-U04",
				Severity: "high",
				Message:  fmt.Sprintf("Version check in %s: verify slice access has bounds check", fn.Name),
				File:     fn.File,
				Line:     cs.Line,
				NodeID:   cs.ID,
			})
		}
	}
	return findings
}
