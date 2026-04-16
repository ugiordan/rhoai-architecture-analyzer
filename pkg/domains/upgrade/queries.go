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
		{ID: "CGA-U02", Name: "pre-release-api-usage", Domain: "upgrade", Severity: "low", Run: queryDeprecatedAPIUsage},
		{ID: "CGA-U03", Name: "ungated-feature", Domain: "upgrade", Severity: "medium", Run: queryUngatedFeature},
		{ID: "CGA-U04", Name: "unchecked-version-access", Domain: "upgrade", Severity: "low", Run: queryUncheckedVersionAccess},
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

// CGA-U02: Pre-release API usage in non-test files.
// Flags v1alpha1/v1beta1 API references. These are pre-release versions that
// may change or be removed. Not necessarily deprecated: a project's own
// v1alpha1 types that are the current and only version are not deprecated.
func queryDeprecatedAPIUsage(g *graph.CPG) []query.Finding {
	var findings []query.Finding
	for _, cs := range g.NodesByKind(graph.NodeCallSite) {
		if !cs.Annotations[AnnotPreReleaseAPI] {
			continue
		}
		if strings.HasSuffix(cs.File, "_test.go") || strings.Contains(cs.File, "vendor/") {
			continue
		}
		findings = append(findings, query.Finding{
			RuleID:   "CGA-U02",
			Severity: "low",
			Message:  fmt.Sprintf("Pre-release API usage: %s at %s:%d (v1alpha1/v1beta1 may change between versions)", cs.Name, cs.File, cs.Line),
			File:     cs.File,
			Line:     cs.Line,
			NodeID:   cs.ID,
		})
	}
	return findings
}

// CGA-U03: Functions that check feature gates at runtime, but the gate they
// reference is not registered in the architecture extraction's feature gate
// inventory. This catches typos, removed gates, and gates defined in external
// packages that may not be available at runtime.
//
// Additionally, detects functions that contain feature-related naming patterns
// (e.g., "enableFeatureX", "handleNewCapability") but lack any feature gate
// check, suggesting the feature may be ungated.
func queryUngatedFeature(g *graph.CPG) []query.Finding {
	if g.ArchData == nil {
		return nil
	}

	// Build set of registered gate names
	registeredGates := make(map[string]bool)
	for _, fg := range g.ArchData.FeatureGates {
		registeredGates[fg.Name] = true
	}

	// If no gates are registered at all, skip: the project likely doesn't use feature gates
	if len(registeredGates) == 0 {
		return nil
	}

	var findings []query.Finding

	// Check 1: Find feature gate Enabled() calls referencing unregistered gates.
	// The annotator marks call sites with AnnotFeatureGate when they call
	// featuregate.Enabled / DefaultFeatureGate.Enabled.
	for _, cs := range g.NodesByKind(graph.NodeCallSite) {
		if !cs.Annotations[AnnotFeatureGate] {
			continue
		}
		if strings.HasSuffix(cs.File, "_test.go") || strings.Contains(cs.File, "vendor/") {
			continue
		}

		// Extract the gate name from the call. The call site name looks like:
		// "utilfeature.DefaultFeatureGate.Enabled" with the argument being
		// the gate constant. We check if the function's enclosing context
		// references any registered gate name.
		gateReferenced := false
		for gateName := range registeredGates {
			if strings.Contains(cs.Name, gateName) {
				gateReferenced = true
				break
			}
		}

		// If none of the registered gates appear in the call, check surrounding
		// lines for the gate constant. This is a best-effort heuristic since
		// the CPG doesn't carry argument details for call sites.
		if !gateReferenced {
			// Look for any registered gate name in the properties
			for _, prop := range cs.Properties {
				for gateName := range registeredGates {
					if strings.Contains(prop, gateName) {
						gateReferenced = true
						break
					}
				}
				if gateReferenced {
					break
				}
			}
		}

		// If we still can't confirm a registered gate is referenced, flag it
		if !gateReferenced {
			findings = append(findings, query.Finding{
				RuleID:   "CGA-U03",
				Severity: "medium",
				Message:  fmt.Sprintf("Feature gate check at %s:%d references a gate not found in the registered feature gate inventory", cs.File, cs.Line),
				File:     cs.File,
				Line:     cs.Line,
				NodeID:   cs.ID,
			})
		}
	}

	// Check 2: Find functions whose names suggest feature-gating but that
	// contain no feature gate check call within their body.
	gatedFunctions := make(map[string]bool)
	for _, cs := range g.NodesByKind(graph.NodeCallSite) {
		if !cs.Annotations[AnnotFeatureGate] {
			continue
		}
		// Find the parent function via incoming edges
		for _, edge := range g.InEdges(cs.ID) {
			fn := g.GetNode(edge.From)
			if fn != nil && fn.Kind == graph.NodeFunction {
				gatedFunctions[fn.ID] = true
			}
		}
	}

	for _, fn := range g.NodesByKind(graph.NodeFunction) {
		if strings.HasSuffix(fn.File, "_test.go") || strings.Contains(fn.File, "vendor/") {
			continue
		}
		if gatedFunctions[fn.ID] {
			continue
		}
		nameLower := strings.ToLower(fn.Name)
		if isFeatureRelatedFunction(nameLower) {
			findings = append(findings, query.Finding{
				RuleID:   "CGA-U03",
				Severity: "medium",
				Message:  fmt.Sprintf("Function %s appears feature-related but contains no feature gate check (project has %d registered gates)", fn.Name, len(registeredGates)),
				File:     fn.File,
				Line:     fn.Line,
				NodeID:   fn.ID,
			})
		}
	}

	return findings
}

// isFeatureRelatedFunction checks if a function name suggests it implements
// a feature that should be gated.
func isFeatureRelatedFunction(nameLower string) bool {
	featurePatterns := []string{
		"enablefeature", "disablefeature",
		"setupfeature", "initfeature",
		"togglefeature", "activatefeature",
		"configurefeature", "registerfeature",
		"turnonfeature", "turnofffeature",
		"withfeature", "handlefeature",
		"featureflag", "featuregate",
	}
	for _, p := range featurePatterns {
		if strings.Contains(nameLower, p) {
			return true
		}
	}
	return false
}

// CGA-U04: Advisory rule. Flags version comparison call sites for manual review.
// No structural verification is performed: the finding indicates a version check
// exists, not that a bounds check is necessarily missing. Severity is low (informational).
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
				Severity: "low",
				Message:  fmt.Sprintf("Advisory: version comparison in %s may need bounds check on subsequent slice access (manual review)", fn.Name),
				File:     fn.File,
				Line:     cs.Line,
				NodeID:   cs.ID,
			})
		}
	}
	return findings
}
