package testing

import (
	"fmt"

	"github.com/ugiordan/architecture-analyzer/pkg/domains/security"
	"github.com/ugiordan/architecture-analyzer/pkg/graph"
	"github.com/ugiordan/architecture-analyzer/pkg/query"
)

func testingQueries() []query.Rule {
	return []query.Rule{
		{ID: "CGA-T01", Name: "untested-security-func", Domain: "testing", Severity: "medium", Run: queryUntestedSecurityFunc},
		{ID: "CGA-T02", Name: "fake-only-integration", Domain: "testing", Severity: "low", Run: queryFakeOnlyIntegration},
		{ID: "CGA-T03", Name: "missing-error-paths", Domain: "testing", Severity: "medium", Run: queryMissingErrorPaths},
		{ID: "CGA-T04", Name: "consolidation-opportunity", Domain: "testing", Severity: "low", Run: queryConsolidationOpportunity},
	}
}

// CGA-T01: Security-annotated functions with no test calling them
func queryUntestedSecurityFunc(g *graph.CPG) []query.Finding {
	testedFns := make(map[string]bool)
	for _, fn := range g.NodesByKind(graph.NodeFunction) {
		if !fn.Annotations[AnnotIsTestFunc] {
			continue
		}
		for _, edge := range g.OutEdges(fn.ID) {
			if edge.Kind == graph.EdgeCalls {
				target := g.GetNode(edge.To)
				if target != nil && target.Kind == graph.NodeCallSite {
					for _, callEdge := range g.OutEdges(target.ID) {
						if callEdge.Kind == graph.EdgeCalls {
							testedFns[callEdge.To] = true
						}
					}
				}
			}
		}
	}

	var findings []query.Finding
	for _, fn := range g.NodesByKind(graph.NodeFunction) {
		isSecurityFunc := fn.Annotations[security.AnnotHandlesAdmission] || fn.Annotations[security.AnnotCreatesRBAC]
		if isSecurityFunc && !testedFns[fn.ID] {
			findings = append(findings, query.Finding{
				RuleID:   "CGA-T01",
				Severity: "medium",
				Message:  fmt.Sprintf("Security function %s has no test calling it", fn.Name),
				File:     fn.File,
				Line:     fn.Line,
				NodeID:   fn.ID,
			})
		}
	}
	return findings
}

// CGA-T02: Reconcile functions tested only with fake client, never envtest.
// Checks which test functions actually call the target function (via call edges)
// and whether those specific tests use fake client vs envtest.
func queryFakeOnlyIntegration(g *graph.CPG) []query.Finding {
	// Build map: function ID -> set of test function IDs that call it
	calledBy := make(map[string][]string)
	for _, testFn := range g.NodesByKind(graph.NodeFunction) {
		if !testFn.Annotations[AnnotIsTestFunc] {
			continue
		}
		for _, edge := range g.OutEdges(testFn.ID) {
			if edge.Kind != graph.EdgeCalls {
				continue
			}
			target := g.GetNode(edge.To)
			if target == nil {
				continue
			}
			if target.Kind == graph.NodeCallSite {
				for _, callEdge := range g.OutEdges(target.ID) {
					if callEdge.Kind == graph.EdgeCalls {
						calledBy[callEdge.To] = append(calledBy[callEdge.To], testFn.ID)
					}
				}
			}
		}
	}

	var findings []query.Finding
	for _, fn := range g.NodesByKind(graph.NodeFunction) {
		if fn.Name != "Reconcile" && fn.Name != "SetupWithManager" {
			continue
		}
		testIDs := calledBy[fn.ID]
		if len(testIDs) == 0 {
			continue
		}
		hasFakeTest := false
		hasEnvtest := false
		for _, testID := range testIDs {
			testFn := g.GetNode(testID)
			if testFn == nil {
				continue
			}
			if testFn.Annotations[AnnotUsesFakeClient] {
				hasFakeTest = true
			}
			if testFn.Annotations[AnnotUsesEnvtest] {
				hasEnvtest = true
			}
		}
		if hasFakeTest && !hasEnvtest {
			findings = append(findings, query.Finding{
				RuleID:   "CGA-T02",
				Severity: "low",
				Message:  fmt.Sprintf("Function %s is tested with fake client but never with envtest", fn.Name),
				File:     fn.File,
				Line:     fn.Line,
				NodeID:   fn.ID,
			})
		}
	}
	return findings
}

// CGA-T03: Test functions without error path assertions
func queryMissingErrorPaths(g *graph.CPG) []query.Finding {
	var findings []query.Finding
	for _, fn := range g.NodesByKind(graph.NodeFunction) {
		if !fn.Annotations[AnnotIsTestFunc] {
			continue
		}
		if !fn.Annotations[AnnotErrorPath] && !fn.Annotations[AnnotIsTestHelper] {
			findings = append(findings, query.Finding{
				RuleID:   "CGA-T03",
				Severity: "medium",
				Message:  fmt.Sprintf("Test function %s has no error path assertions", fn.Name),
				File:     fn.File,
				Line:     fn.Line,
				NodeID:   fn.ID,
			})
		}
	}
	return findings
}

// CGA-T04: Test files with many non-table-driven tests
func queryConsolidationOpportunity(g *graph.CPG) []query.Finding {
	testsByFile := make(map[string]int)
	tableByFile := make(map[string]int)

	for _, fn := range g.NodesByKind(graph.NodeFunction) {
		if !fn.Annotations[AnnotIsTestFunc] {
			continue
		}
		testsByFile[fn.File]++
		if fn.Annotations[AnnotTableDriven] {
			tableByFile[fn.File]++
		}
	}

	var findings []query.Finding
	for file, count := range testsByFile {
		nonTable := count - tableByFile[file]
		if nonTable > 5 {
			findings = append(findings, query.Finding{
				RuleID:   "CGA-T04",
				Severity: "low",
				Message:  fmt.Sprintf("File %s has %d non-table-driven tests (consolidation opportunity)", file, nonTable),
				File:     file,
				Line:     1,
			})
		}
	}
	return findings
}
