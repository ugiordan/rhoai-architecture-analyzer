package builder

import (
	"testing"

	"github.com/ugiordan/rhoai-architecture-analyzer/pkg/graph"
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
