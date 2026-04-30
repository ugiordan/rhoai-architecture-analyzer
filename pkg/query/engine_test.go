package query

import (
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

func TestQueryMissingAuth(t *testing.T) {
	cpg := graph.NewCPG()

	noAuth := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "handleDelete",
		File:        "handler.go",
		Line:        10,
		Annotations: map[string]bool{
			"handles_user_input": true,
			"mutates_state":     true,
		},
	}
	if err := cpg.AddNode(noAuth); err != nil { t.Fatal(err) }

	withAuth := &graph.Node{
		ID:          "fn2",
		Kind:        graph.NodeFunction,
		Name:        "handleSafeDelete",
		File:        "handler.go",
		Line:        20,
		Annotations: map[string]bool{
			"handles_user_input": true,
			"mutates_state":     true,
			"has_auth":          true,
		},
	}
	if err := cpg.AddNode(withAuth); err != nil { t.Fatal(err) }

	engine := NewEngine()
	findings := engine.QueryMissingAuth(cpg)

	if len(findings) != 1 {
		t.Errorf("expected 1 finding, got %d", len(findings))
	}
	if len(findings) > 0 && findings[0].NodeID != "fn1" {
		t.Errorf("expected finding for fn1, got %s", findings[0].NodeID)
	}
}

func TestQueryTaintToExternalSink(t *testing.T) {
	cpg := graph.NewCPG()

	handler := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "handler",
		File:        "handler.go",
		Line:        5,
		Annotations: map[string]bool{"handles_user_input": true},
	}
	if err := cpg.AddNode(handler); err != nil { t.Fatal(err) }

	extCall := &graph.Node{
		ID:          "ext1",
		Kind:        graph.NodeCallSite,
		Name:        "llmClient.Complete",
		Annotations: map[string]bool{"calls_external": true},
	}
	if err := cpg.AddNode(extCall); err != nil { t.Fatal(err) }

	// Pre-computed EdgeTaint edge (produced by TaintEngine)
	cpg.AddEdge(&graph.Edge{
		From:  "fn1",
		To:    "ext1",
		Kind:  graph.EdgeTaint,
		Label: "calls_external",
		Path:  []string{"fn1", "db_w", "db_r", "ext1"},
	})

	engine := NewEngine()
	findings := engine.QueryCrossStorageTaint(cpg)

	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].RuleID != "CGA-002" {
		t.Errorf("expected rule CGA-002, got %s", findings[0].RuleID)
	}
	if findings[0].File != "handler.go" {
		t.Errorf("expected file handler.go, got %s", findings[0].File)
	}
	if len(findings[0].Path) != 4 {
		t.Errorf("expected path length 4, got %d", len(findings[0].Path))
	}
}

func TestRunRules(t *testing.T) {
	cpg := graph.NewCPG()
	fn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "testFunc",
		File:        "test.go",
		Line:        1,
		Annotations: map[string]bool{"test_annotation": true},
	}
	if err := cpg.AddNode(fn); err != nil { t.Fatal(err) }

	rules := []Rule{
		{
			ID:       "TEST-001",
			Name:     "test-rule",
			Domain:   "test",
			Severity: "medium",
			Run: func(g *graph.CPG) []Finding {
				var findings []Finding
				for _, n := range g.NodesByKind(graph.NodeFunction) {
					if n.Annotations["test_annotation"] {
						findings = append(findings, Finding{
							RuleID:   "TEST-001",
							Severity: "medium",
							Message:  "test finding: " + n.Name,
							File:     n.File,
							Line:     n.Line,
							NodeID:   n.ID,
						})
					}
				}
				return findings
			},
		},
	}

	e := NewEngine()
	findings := e.RunRules(cpg, rules)
	if len(findings) != 1 {
		t.Errorf("expected 1 finding, got %d", len(findings))
	}
	if len(findings) > 0 && findings[0].RuleID != "TEST-001" {
		t.Errorf("expected rule ID TEST-001, got %s", findings[0].RuleID)
	}
}
