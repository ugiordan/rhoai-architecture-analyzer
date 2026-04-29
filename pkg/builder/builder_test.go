package builder

import (
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
	"github.com/ugiordan/architecture-analyzer/pkg/parser"
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

func TestMergeResultIncludesVariablesAndParameters(t *testing.T) {
	cpg := graph.NewCPG()
	b := NewBuilder()

	result := &parser.ParseResult{
		Variables: []*graph.Node{
			{ID: "var_test1", Kind: graph.NodeVariable, Name: "x", File: "test.go", Line: 5},
			{ID: "var_test2", Kind: graph.NodeVariable, Name: "y", File: "test.go", Line: 6},
		},
		Parameters: []*graph.Node{
			{ID: "param_test1", Kind: graph.NodeParameter, Name: "r", File: "test.go", Line: 3},
		},
	}

	if err := b.mergeResult(cpg, result); err != nil {
		t.Fatalf("mergeResult failed: %v", err)
	}

	vars := cpg.NodesByKind(graph.NodeVariable)
	if len(vars) != 2 {
		t.Errorf("got %d variables, want 2", len(vars))
	}

	params := cpg.NodesByKind(graph.NodeParameter)
	if len(params) != 1 {
		t.Errorf("got %d parameters, want 1", len(params))
	}
}

func TestBuildFromDirectoryIncludesDataFlow(t *testing.T) {
	b := NewBuilder()
	cpg, err := b.BuildFromDir("../../testdata")
	if err != nil {
		t.Fatalf("BuildFromDir failed: %v", err)
	}

	// Verify NodeVariable and NodeParameter nodes appear in CPG
	vars := cpg.NodesByKind(graph.NodeVariable)
	if len(vars) == 0 {
		t.Error("expected NodeVariable nodes in CPG, got 0")
	}

	params := cpg.NodesByKind(graph.NodeParameter)
	if len(params) == 0 {
		t.Error("expected NodeParameter nodes in CPG, got 0")
	}

	// Verify EdgeDataFlow edges with data flow labels appear
	dataFlowLabels := make(map[string]bool)
	for _, e := range cpg.Edges() {
		if e.Kind == graph.EdgeDataFlow {
			dataFlowLabels[e.Label] = true
		}
	}

	for _, label := range []string{"declares", "assigns"} {
		if !dataFlowLabels[label] {
			t.Errorf("expected EdgeDataFlow with label %q in CPG", label)
		}
	}

	t.Logf("CPG contains %d variables, %d parameters", len(vars), len(params))
	t.Logf("Data flow labels found: %v", dataFlowLabels)
}

func TestMergeResultIncludesBasicBlocks(t *testing.T) {
	b := NewBuilder()
	cpg := graph.NewCPG()

	result := &parser.ParseResult{
		BasicBlocks: []*graph.Node{
			{ID: "bb_entry", Kind: graph.NodeBasicBlock, Name: "entry", File: "main.go", Line: 10, ParentID: "fn_abc"},
			{ID: "bb_exit", Kind: graph.NodeBasicBlock, Name: "exit", File: "main.go", Line: 0, ParentID: "fn_abc"},
		},
		Edges: []*graph.Edge{
			{From: "fn_abc", To: "bb_entry", Kind: graph.EdgeControlFlow, Label: "entry", Confidence: graph.ConfidenceCertain},
		},
	}

	if err := b.mergeResult(cpg, result); err != nil {
		t.Fatalf("mergeResult failed: %v", err)
	}

	bbs := cpg.NodesByKind(graph.NodeBasicBlock)
	if len(bbs) != 2 {
		t.Errorf("got %d BasicBlock nodes, want 2", len(bbs))
	}

	cfEdges := 0
	for _, e := range cpg.Edges() {
		if e.Kind == graph.EdgeControlFlow {
			cfEdges++
		}
	}
	if cfEdges != 1 {
		t.Errorf("got %d CONTROL_FLOW edges, want 1", cfEdges)
	}
}

func TestBuildFromDirectoryIncludesCFG(t *testing.T) {
	b := NewBuilder()
	cpg, err := b.BuildFromDir("../../testdata")
	if err != nil {
		t.Fatalf("BuildFromDir failed: %v", err)
	}

	// Verify BasicBlock nodes exist in the CPG
	bbs := cpg.NodesByKind(graph.NodeBasicBlock)
	if len(bbs) == 0 {
		t.Error("expected BasicBlock nodes in CPG, got 0")
	}

	// Verify all basic blocks have ParentID set
	for _, bb := range bbs {
		if bb.ParentID == "" {
			t.Errorf("BasicBlock %q has empty ParentID", bb.Name)
		}
	}

	// Verify EdgeControlFlow edges exist
	cfEdges := 0
	cfLabels := make(map[string]bool)
	for _, e := range cpg.Edges() {
		if e.Kind == graph.EdgeControlFlow {
			cfEdges++
			cfLabels[e.Label] = true
		}
	}
	if cfEdges == 0 {
		t.Error("expected CONTROL_FLOW edges in CPG, got 0")
	}

	// Verify at least entry and exit labels exist
	if !cfLabels["entry"] {
		t.Error("expected 'entry' label in CONTROL_FLOW edges")
	}
	if !cfLabels["exit"] {
		t.Error("expected 'exit' label in CONTROL_FLOW edges")
	}

	t.Logf("CPG has %d BasicBlock nodes, %d CONTROL_FLOW edges, labels: %v", len(bbs), cfEdges, cfLabels)
}
