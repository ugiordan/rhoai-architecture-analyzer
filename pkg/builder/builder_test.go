package builder

import (
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

func TestBuildFromDirectory(t *testing.T) {
	b := NewBuilder()
	cpg, err := b.BuildFromDir("../../testdata")
	if err != nil {
		t.Fatalf("BuildFromDir failed: %v", err)
	}

	if len(cpg.Nodes()) == 0 {
		t.Error("expected nodes, got 0")
	}

	fns := cpg.NodesByKind(graph.NodeFunction)
	if len(fns) < 3 {
		t.Errorf("expected at least 3 functions, got %d", len(fns))
	}

	callEdges := 0
	for _, e := range cpg.Edges() {
		if e.Kind == graph.EdgeCalls {
			callEdges++
		}
	}
	if callEdges == 0 {
		t.Error("expected CALLS edges, got 0")
	}
}

func TestBuilderLinksStructLiteralsToFunctions(t *testing.T) {
	b := NewBuilder()
	cpg, err := b.BuildFromDir("../../testdata")
	if err != nil {
		t.Fatalf("BuildFromDir failed: %v", err)
	}

	structLiterals := cpg.NodesByKind(graph.NodeStructLiteral)
	if len(structLiterals) == 0 {
		t.Fatal("expected struct literals from testdata fixtures")
	}

	// Check that at least one struct literal has an incoming DATA_FLOW edge with "contains_struct"
	hasContainment := false
	for _, sl := range structLiterals {
		for _, edge := range cpg.InEdges(sl.ID) {
			if edge.Kind == graph.EdgeDataFlow && edge.Label == "contains_struct" {
				hasContainment = true
				break
			}
		}
		if hasContainment {
			break
		}
	}
	if !hasContainment {
		t.Error("expected at least one struct literal to be linked to its containing function")
	}
}

func TestResolveCallEdgesConfidence(t *testing.T) {
	b := NewBuilder()
	cpg := graph.NewCPG()

	// Same-package function and call site
	cpg.AddNode(&graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "doStuff",
		File: "pkg/handler/handler.go", Line: 10, EndLine: 20, Language: "go",
	})
	cpg.AddNode(&graph.Node{
		ID: "cs1", Kind: graph.NodeCallSite, Name: "doStuff",
		File: "pkg/handler/handler.go", Line: 15, Language: "go",
	})
	// Cross-package function (same short name, different dir)
	cpg.AddNode(&graph.Node{
		ID: "fn2", Kind: graph.NodeFunction, Name: "doStuff",
		File: "pkg/util/util.go", Line: 5, EndLine: 15, Language: "go",
	})
	// Qualified call to cross-package
	cpg.AddNode(&graph.Node{
		ID: "cs2", Kind: graph.NodeCallSite, Name: "util.doStuff",
		File: "pkg/handler/handler.go", Line: 18, Language: "go",
	})

	b.resolveCallEdges(cpg)

	// cs1 -> fn1 should be CERTAIN (same package, exact name)
	edges1 := cpg.OutEdges("cs1")
	callEdges1 := filterCallEdges(edges1)
	if len(callEdges1) == 0 {
		t.Fatal("expected call edge from cs1")
	}
	if callEdges1[0].Confidence != graph.ConfidenceCertain {
		t.Errorf("expected CERTAIN for same-package call, got %q", callEdges1[0].Confidence)
	}

	// cs2 -> fn2 should be INFERRED (cross-package, short-name match)
	edges2 := cpg.OutEdges("cs2")
	callEdges2 := filterCallEdges(edges2)
	if len(callEdges2) == 0 {
		t.Fatal("expected call edge from cs2")
	}
	for _, e := range callEdges2 {
		if e.To == "fn2" && e.Confidence != graph.ConfidenceInferred {
			t.Errorf("expected INFERRED for cross-package call, got %q", e.Confidence)
		}
	}
}

func filterCallEdges(edges []*graph.Edge) []*graph.Edge {
	var result []*graph.Edge
	for _, e := range edges {
		if e.Kind == graph.EdgeCalls {
			result = append(result, e)
		}
	}
	return result
}
