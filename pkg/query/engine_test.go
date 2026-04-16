package query

import (
	"testing"

	"github.com/ugiordan/rhoai-architecture-analyzer/pkg/graph"
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
	cpg.AddNode(noAuth)

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
	cpg.AddNode(withAuth)

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
		Annotations: map[string]bool{"handles_user_input": true},
	}
	cpg.AddNode(handler)

	dbWrite := &graph.Node{
		ID:         "db_w",
		Kind:       graph.NodeDBOperation,
		Properties: map[string]string{"operation": "write", "table": "templates"},
	}
	cpg.AddNode(dbWrite)

	dbRead := &graph.Node{
		ID:         "db_r",
		Kind:       graph.NodeDBOperation,
		Properties: map[string]string{"operation": "read", "table": "templates"},
	}
	cpg.AddNode(dbRead)

	extCall := &graph.Node{
		ID:          "ext1",
		Kind:        graph.NodeCallSite,
		Name:        "llmClient.Complete",
		Annotations: map[string]bool{"calls_external": true},
	}
	cpg.AddNode(extCall)

	cpg.AddEdge(&graph.Edge{From: "fn1", To: "db_w", Kind: graph.EdgeDataFlow})
	cpg.AddEdge(&graph.Edge{From: "db_w", To: "db_r", Kind: graph.EdgeStorageLink})
	cpg.AddEdge(&graph.Edge{From: "db_r", To: "ext1", Kind: graph.EdgeDataFlow})

	engine := NewEngine()
	findings := engine.QueryCrossStorageTaint(cpg)

	if len(findings) == 0 {
		t.Error("expected cross-storage taint finding, got 0")
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
	cpg.AddNode(fn)

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
