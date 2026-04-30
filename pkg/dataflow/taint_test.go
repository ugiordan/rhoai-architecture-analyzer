package dataflow

import (
	"fmt"
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

// --- test helpers ---

func buildTestCPG() *graph.CPG { return graph.NewCPG() }

func addFunc(cpg *graph.CPG, id, name string) *graph.Node {
	n := &graph.Node{
		ID:          id,
		Kind:        graph.NodeFunction,
		Name:        name,
		Annotations: make(map[string]bool),
	}
	if err := cpg.AddNode(n); err != nil {
		panic(err)
	}
	return n
}

func addParam(cpg *graph.CPG, funcID, paramID, name string) *graph.Node {
	n := &graph.Node{
		ID:          paramID,
		Kind:        graph.NodeParameter,
		Name:        name,
		Annotations: make(map[string]bool),
	}
	if err := cpg.AddNode(n); err != nil {
		panic(err)
	}
	cpg.AddEdge(&graph.Edge{
		From:       funcID,
		To:         paramID,
		Kind:       graph.EdgeDataFlow,
		Label:      "declares",
		Confidence: graph.ConfidenceCertain,
	})
	return n
}

func addVar(cpg *graph.CPG, id, name string) *graph.Node {
	n := &graph.Node{
		ID:          id,
		Kind:        graph.NodeVariable,
		Name:        name,
		Annotations: make(map[string]bool),
	}
	if err := cpg.AddNode(n); err != nil {
		panic(err)
	}
	return n
}

func addCallSite(cpg *graph.CPG, funcID, csID, name string) *graph.Node {
	n := &graph.Node{
		ID:          csID,
		Kind:        graph.NodeCallSite,
		Name:        name,
		Annotations: make(map[string]bool),
	}
	if err := cpg.AddNode(n); err != nil {
		panic(err)
	}
	cpg.AddEdge(&graph.Edge{
		From:       funcID,
		To:         csID,
		Kind:       graph.EdgeDataFlow,
		Label:      "contains_call",
		Confidence: graph.ConfidenceCertain,
	})
	return n
}

func addBlock(cpg *graph.CPG, funcID, blockID, name string, members []string) *graph.Node {
	n := &graph.Node{
		ID:          blockID,
		Kind:        graph.NodeBasicBlock,
		Name:        name,
		ParentID:    funcID,
		Members:     members,
		Annotations: make(map[string]bool),
	}
	if err := cpg.AddNode(n); err != nil {
		panic(err)
	}
	return n
}

func dataFlowEdge(from, to, label string) *graph.Edge {
	return &graph.Edge{
		From:       from,
		To:         to,
		Kind:       graph.EdgeDataFlow,
		Label:      label,
		Confidence: graph.ConfidenceCertain,
	}
}

func cfgEdge(from, to, label string) *graph.Edge {
	return &graph.Edge{
		From:       from,
		To:         to,
		Kind:       graph.EdgeControlFlow,
		Label:      label,
		Confidence: graph.ConfidenceCertain,
	}
}

// --- tests ---

// TestTaintIntraSimple: function with handles_user_input, param -> assigns -> var -> passes_to -> sink call site.
// All in one basic block. Expect 1 EdgeTaint CERTAIN.
func TestTaintIntraSimple(t *testing.T) {
	cpg := buildTestCPG()

	fn := addFunc(cpg, "func1", "handleRequest")
	fn.Annotations["handles_user_input"] = true

	param := addParam(cpg, "func1", "param1", "input")
	_ = param

	query := addVar(cpg, "var1", "query")
	_ = query

	cs := addCallSite(cpg, "func1", "cs1", "db.Query")
	cs.Annotations["sec:executes_sql"] = true

	// single basic block containing all nodes
	addBlock(cpg, "func1", "block1", "entry", []string{"param1", "var1", "cs1"})

	// data flow: param1 -> assigns -> var1 -> passes_to -> cs1
	cpg.AddEdge(dataFlowEdge("param1", "var1", "assigns"))
	cpg.AddEdge(dataFlowEdge("var1", "cs1", "passes_to"))

	te := NewTaintEngine()
	edges := te.Run(cpg)

	if len(edges) != 1 {
		t.Fatalf("expected 1 taint edge, got %d", len(edges))
	}
	e := edges[0]
	if e.Kind != graph.EdgeTaint {
		t.Errorf("expected EdgeTaint, got %s", e.Kind)
	}
	if e.Confidence != graph.ConfidenceCertain {
		t.Errorf("expected CERTAIN confidence, got %s", e.Confidence)
	}
	if e.From != "param1" {
		t.Errorf("expected from=param1, got %s", e.From)
	}
	if e.To != "cs1" {
		t.Errorf("expected to=cs1, got %s", e.To)
	}
	if e.Label != "handles_user_input->sec:executes_sql" {
		t.Errorf("unexpected label: %s", e.Label)
	}
}

// TestTaintIntraCFGFilter: source param in entry block, v1 in then-block, v2 in else-block with sink.
// No data flow from param to v2. Expect 0 taint edges.
func TestTaintIntraCFGFilter(t *testing.T) {
	cpg := buildTestCPG()

	fn := addFunc(cpg, "func2", "branchingFunc")
	fn.Annotations["handles_user_input"] = true

	param := addParam(cpg, "func2", "p2", "input")
	_ = param

	v1 := addVar(cpg, "v1", "v1")
	_ = v1

	v2 := addVar(cpg, "v2", "v2")
	_ = v2

	cs := addCallSite(cpg, "func2", "cs2", "exec")
	cs.Annotations["sec:subprocess_call"] = true

	// blocks
	addBlock(cpg, "func2", "entry2", "entry", []string{"p2"})
	addBlock(cpg, "func2", "then2", "then", []string{"v1"})
	addBlock(cpg, "func2", "else2", "else", []string{"v2", "cs2"})

	// CFG: entry -> then (true_branch), entry -> else (false_branch)
	cpg.AddEdge(cfgEdge("entry2", "then2", "true_branch"))
	cpg.AddEdge(cfgEdge("entry2", "else2", "false_branch"))

	// data flow: param assigns to v1 (in then block), but no flow to v2/cs2
	cpg.AddEdge(dataFlowEdge("p2", "v1", "assigns"))
	// v2 passes_to cs2 but v2 has no taint source flowing into it
	cpg.AddEdge(dataFlowEdge("v2", "cs2", "passes_to"))

	te := NewTaintEngine()
	edges := te.Run(cpg)

	if len(edges) != 0 {
		t.Fatalf("expected 0 taint edges (no data flow to sink), got %d", len(edges))
	}
}

// TestTaintNoSources: function with no source annotations. Has a sink. Expect 0 edges, no panic.
func TestTaintNoSources(t *testing.T) {
	cpg := buildTestCPG()

	addFunc(cpg, "func3", "cleanFunc")
	addParam(cpg, "func3", "p3", "x")

	cs := addCallSite(cpg, "func3", "cs3", "exec")
	cs.Annotations["sec:subprocess_call"] = true

	addBlock(cpg, "func3", "block3", "entry", []string{"p3", "cs3"})
	cpg.AddEdge(dataFlowEdge("p3", "cs3", "passes_to"))

	te := NewTaintEngine()
	edges := te.Run(cpg)

	if len(edges) != 0 {
		t.Fatalf("expected 0 taint edges, got %d", len(edges))
	}
}

// TestTaintNoSinks: function with source annotation but no sinks. Expect 0 edges.
func TestTaintNoSinks(t *testing.T) {
	cpg := buildTestCPG()

	fn := addFunc(cpg, "func4", "sourceOnly")
	fn.Annotations["handles_user_input"] = true

	addParam(cpg, "func4", "p4", "data")
	v := addVar(cpg, "v4", "temp")
	_ = v

	addBlock(cpg, "func4", "block4", "entry", []string{"p4", "v4"})
	cpg.AddEdge(dataFlowEdge("p4", "v4", "assigns"))

	te := NewTaintEngine()
	edges := te.Run(cpg)

	if len(edges) != 0 {
		t.Fatalf("expected 0 taint edges, got %d", len(edges))
	}
}

func addDBOp(cpg *graph.CPG, id, name, operation, table string) *graph.Node {
	n := &graph.Node{
		ID:          id,
		Kind:        graph.NodeDBOperation,
		Name:        name,
		Operation:   operation,
		Table:       table,
		Annotations: make(map[string]bool),
	}
	if err := cpg.AddNode(n); err != nil {
		panic(err)
	}
	return n
}

// TestTaintInterCall: caller "handleRequest" with handles_user_input, param "input" passes_to
// call site "cs_query" which calls callee "query". Callee has param "sql" that passes_to
// sink "db.Exec" with sec:executes_sql. EdgeCalls from cs_query to callee with CERTAIN confidence.
// Expect interprocedural EdgeTaint from caller's param to callee's sink with INFERRED confidence.
func TestTaintInterCall(t *testing.T) {
	cpg := buildTestCPG()

	// Caller function
	caller := addFunc(cpg, "fn_caller", "handleRequest")
	caller.Annotations["handles_user_input"] = true
	caller.ParamNames = []string{"input"}

	callerParam := addParam(cpg, "fn_caller", "p_caller", "input")
	_ = callerParam

	csQuery := addCallSite(cpg, "fn_caller", "cs_query", "query")
	_ = csQuery

	// Caller basic block
	addBlock(cpg, "fn_caller", "blk_caller", "entry", []string{"p_caller", "cs_query"})

	// Data flow: caller param -> passes_to -> call site
	cpg.AddEdge(dataFlowEdge("p_caller", "cs_query", "passes_to"))

	// Callee function
	callee := addFunc(cpg, "fn_callee", "query")
	callee.ParamNames = []string{"sql"}

	calleeParam := addParam(cpg, "fn_callee", "p_callee", "sql")
	_ = calleeParam

	sink := addCallSite(cpg, "fn_callee", "cs_exec", "db.Exec")
	sink.Annotations["sec:executes_sql"] = true

	// Callee basic block
	addBlock(cpg, "fn_callee", "blk_callee", "entry", []string{"p_callee", "cs_exec"})

	// Data flow: callee param -> passes_to -> sink
	cpg.AddEdge(dataFlowEdge("p_callee", "cs_exec", "passes_to"))

	// Call edge: cs_query -> fn_callee
	cpg.AddEdge(&graph.Edge{
		From:       "cs_query",
		To:         "fn_callee",
		Kind:       graph.EdgeCalls,
		Label:      "calls",
		Confidence: graph.ConfidenceCertain,
	})

	te := NewTaintEngine()
	edges := te.Run(cpg)

	// Should have interprocedural taint edge
	var interEdges []*graph.Edge
	for _, e := range edges {
		if e.Kind == graph.EdgeTaint && e.Confidence == graph.ConfidenceInferred {
			interEdges = append(interEdges, e)
		}
	}

	if len(interEdges) != 1 {
		t.Fatalf("expected 1 interprocedural taint edge, got %d (total edges: %d)", len(interEdges), len(edges))
	}
	e := interEdges[0]
	if e.From != "p_caller" {
		t.Errorf("expected from=p_caller, got %s", e.From)
	}
	if e.To != "cs_exec" {
		t.Errorf("expected to=cs_exec, got %s", e.To)
	}
}

// TestTaintInterReturn: caller has param "input" passes_to call site "cs_sanitize" which calls
// "sanitize". Call site result assigns to var "result" which passes_to sink "cs_exec".
// The intraprocedural path should find the sink via Phase A (p->cs_sanitize->v_result->cs_exec).
func TestTaintInterReturn(t *testing.T) {
	cpg := buildTestCPG()

	fn := addFunc(cpg, "fn_ret", "process")
	fn.Annotations["handles_user_input"] = true
	fn.ParamNames = []string{"input"}

	param := addParam(cpg, "fn_ret", "p_input", "input")
	_ = param

	csSanitize := addCallSite(cpg, "fn_ret", "cs_sanitize", "sanitize")
	_ = csSanitize

	vResult := addVar(cpg, "v_result", "result")
	_ = vResult

	csExec := addCallSite(cpg, "fn_ret", "cs_exec", "exec")
	csExec.Annotations["sec:subprocess_call"] = true

	addBlock(cpg, "fn_ret", "blk_ret", "entry", []string{"p_input", "cs_sanitize", "v_result", "cs_exec"})

	// Data flow: param -> passes_to -> cs_sanitize -> assigns -> v_result -> passes_to -> cs_exec
	cpg.AddEdge(dataFlowEdge("p_input", "cs_sanitize", "passes_to"))
	cpg.AddEdge(dataFlowEdge("cs_sanitize", "v_result", "assigns"))
	cpg.AddEdge(dataFlowEdge("v_result", "cs_exec", "passes_to"))

	// Callee function (sanitize)
	callee := addFunc(cpg, "fn_sanitize", "sanitize")
	callee.ParamNames = []string{"s"}
	addParam(cpg, "fn_sanitize", "p_s", "s")
	addBlock(cpg, "fn_sanitize", "blk_sanitize", "entry", []string{"p_s"})
	cpg.AddEdge(dataFlowEdge("p_s", "fn_sanitize", "returns"))

	// Call edge
	cpg.AddEdge(&graph.Edge{
		From:       "cs_sanitize",
		To:         "fn_sanitize",
		Kind:       graph.EdgeCalls,
		Label:      "calls",
		Confidence: graph.ConfidenceCertain,
	})

	te := NewTaintEngine()
	edges := te.Run(cpg)

	// Phase A should find: p_input -> cs_sanitize -> v_result -> cs_exec (intraprocedural)
	found := false
	for _, e := range edges {
		if e.Kind == graph.EdgeTaint && e.To == "cs_exec" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected taint edge reaching cs_exec, got %d edges total", len(edges))
	}
}

// TestTaintStorageLink: writer function with handles_user_input, param passes_to DB write.
// Reader function has DB read linked via EdgeStorageLink. Read assigns to var, var passes_to sink.
// Expect EdgeTaint with UNCERTAIN confidence.
func TestTaintStorageLink(t *testing.T) {
	cpg := buildTestCPG()

	// Writer function
	writer := addFunc(cpg, "fn_writer", "writeData")
	writer.Annotations["handles_user_input"] = true
	writer.ParamNames = []string{"data"}

	wParam := addParam(cpg, "fn_writer", "p_wdata", "data")
	_ = wParam

	_ = addDBOp(cpg, "db_write", "db.Insert", "write", "users")

	addBlock(cpg, "fn_writer", "blk_writer", "entry", []string{"p_wdata", "db_write"})
	cpg.AddEdge(dataFlowEdge("p_wdata", "db_write", "passes_to"))

	// Reader function
	reader := addFunc(cpg, "fn_reader", "readData")
	_ = reader

	_ = addDBOp(cpg, "db_read", "db.Select", "read", "users")

	vUser := addVar(cpg, "v_user", "user")
	_ = vUser

	renderCS := addCallSite(cpg, "fn_reader", "cs_render", "template.Render")
	renderCS.Annotations["sec:template_render"] = true

	addBlock(cpg, "fn_reader", "blk_reader", "entry", []string{"db_read", "v_user", "cs_render"})
	cpg.AddEdge(dataFlowEdge("db_read", "v_user", "assigns"))
	cpg.AddEdge(dataFlowEdge("v_user", "cs_render", "passes_to"))

	// Storage link: db_write -> db_read
	cpg.AddEdge(&graph.Edge{
		From:       "db_write",
		To:         "db_read",
		Kind:       graph.EdgeStorageLink,
		Label:      "storage_link",
		Confidence: graph.ConfidenceCertain,
	})

	te := NewTaintEngine()
	edges := te.Run(cpg)

	var storageEdges []*graph.Edge
	for _, e := range edges {
		if e.Kind == graph.EdgeTaint && e.Confidence == graph.ConfidenceUncertain {
			storageEdges = append(storageEdges, e)
		}
	}

	if len(storageEdges) != 1 {
		t.Fatalf("expected 1 storage-linked taint edge, got %d (total edges: %d)", len(storageEdges), len(edges))
	}
	e := storageEdges[0]
	if e.To != "cs_render" {
		t.Errorf("expected to=cs_render, got %s", e.To)
	}
}

// TestTaintBoundsEnforced: function with 150 parallel paths, expect <= MaxTaintPaths edges.
func TestTaintBoundsEnforced(t *testing.T) {
	cpg := buildTestCPG()

	fn := addFunc(cpg, "fn_bounds", "bigFunc")
	fn.Annotations["handles_user_input"] = true
	fn.ParamNames = []string{"input"}

	addParam(cpg, "fn_bounds", "p_bounds", "input")

	members := []string{"p_bounds"}
	for i := 0; i < 150; i++ {
		varID := fmt.Sprintf("v_%d", i)
		sinkID := fmt.Sprintf("sink_%d", i)

		v := addVar(cpg, varID, fmt.Sprintf("var%d", i))
		_ = v

		s := addCallSite(cpg, "fn_bounds", sinkID, fmt.Sprintf("db.Query%d", i))
		s.Annotations["sec:executes_sql"] = true

		members = append(members, varID, sinkID)

		cpg.AddEdge(dataFlowEdge("p_bounds", varID, "assigns"))
		cpg.AddEdge(dataFlowEdge(varID, sinkID, "passes_to"))
	}

	addBlock(cpg, "fn_bounds", "blk_bounds", "entry", members)

	te := NewTaintEngine()
	edges := te.Run(cpg)

	if len(edges) == 0 {
		t.Fatal("expected at least some taint edges, got 0")
	}
	if len(edges) > MaxTaintPaths {
		t.Fatalf("expected at most %d taint edges, got %d", MaxTaintPaths, len(edges))
	}
}

// TestTaintInterDepthChain: A calls B calls C, C has sink. Tests multi-hop interprocedural
// taint and verifies depth limiting works.
func TestTaintInterDepthChain(t *testing.T) {
	cpg := buildTestCPG()

	// Function A: source, param passes_to call site calling B
	fnA := addFunc(cpg, "fn_a", "handler")
	fnA.Annotations["handles_user_input"] = true
	fnA.ParamNames = []string{"input"}
	addParam(cpg, "fn_a", "p_a", "input")
	csB := addCallSite(cpg, "fn_a", "cs_b", "processInput")
	_ = csB
	addBlock(cpg, "fn_a", "blk_a", "entry", []string{"p_a", "cs_b"})
	cpg.AddEdge(dataFlowEdge("p_a", "cs_b", "passes_to"))

	// Function B: param passes_to call site calling C
	fnB := addFunc(cpg, "fn_b", "processInput")
	fnB.ParamNames = []string{"data"}
	addParam(cpg, "fn_b", "p_b", "data")
	csC := addCallSite(cpg, "fn_b", "cs_c", "executeQuery")
	_ = csC
	addBlock(cpg, "fn_b", "blk_b", "entry", []string{"p_b", "cs_c"})
	cpg.AddEdge(dataFlowEdge("p_b", "cs_c", "passes_to"))

	// Function C: param passes_to sink
	fnC := addFunc(cpg, "fn_c", "executeQuery")
	fnC.ParamNames = []string{"sql"}
	addParam(cpg, "fn_c", "p_c", "sql")
	sinkCS := addCallSite(cpg, "fn_c", "cs_sink", "db.Exec")
	sinkCS.Annotations["sec:executes_sql"] = true
	addBlock(cpg, "fn_c", "blk_c", "entry", []string{"p_c", "cs_sink"})
	cpg.AddEdge(dataFlowEdge("p_c", "cs_sink", "passes_to"))

	// Call edges: A->B, B->C
	cpg.AddEdge(&graph.Edge{From: "cs_b", To: "fn_b", Kind: graph.EdgeCalls, Label: "calls", Confidence: graph.ConfidenceCertain})
	cpg.AddEdge(&graph.Edge{From: "cs_c", To: "fn_c", Kind: graph.EdgeCalls, Label: "calls", Confidence: graph.ConfidenceCertain})

	te := NewTaintEngine()
	edges := te.Run(cpg)

	// Should have a taint edge from A's param to C's sink with UNCERTAIN (depth >= 2)
	var taintEdges []*graph.Edge
	for _, e := range edges {
		if e.Kind == graph.EdgeTaint && e.To == "cs_sink" {
			taintEdges = append(taintEdges, e)
		}
	}

	if len(taintEdges) != 1 {
		t.Fatalf("expected 1 taint edge to cs_sink, got %d (total: %d)", len(taintEdges), len(edges))
	}
	e := taintEdges[0]
	if e.Confidence != graph.ConfidenceUncertain {
		t.Errorf("expected UNCERTAIN confidence for 2-hop chain, got %s", e.Confidence)
	}
	if e.From != "p_a" {
		t.Errorf("expected from=p_a, got %s", e.From)
	}
}

// TestTaintMaxDepthEnforced: chain of depth > maxDepth. Verify engine terminates and
// does not produce edges beyond the limit.
func TestTaintMaxDepthEnforced(t *testing.T) {
	cpg := buildTestCPG()

	depth := 5
	// Build a chain: fn_0 -> fn_1 -> fn_2 -> fn_3 -> fn_4 -> fn_5 (sink)
	for i := 0; i <= depth; i++ {
		fnID := fmt.Sprintf("fn_%d", i)
		fn := addFunc(cpg, fnID, fmt.Sprintf("func%d", i))
		fn.ParamNames = []string{"x"}
		addParam(cpg, fnID, fmt.Sprintf("p_%d", i), "x")

		if i == 0 {
			fn.Annotations["handles_user_input"] = true
		}

		if i == depth {
			// Last function has a sink
			sinkCS := addCallSite(cpg, fnID, "cs_sink_deep", "dangerous")
			sinkCS.Annotations["sec:subprocess_call"] = true
			addBlock(cpg, fnID, fmt.Sprintf("blk_%d", i), "entry", []string{fmt.Sprintf("p_%d", i), "cs_sink_deep"})
			cpg.AddEdge(dataFlowEdge(fmt.Sprintf("p_%d", i), "cs_sink_deep", "passes_to"))
		} else {
			// Intermediate function: param passes_to call site calling next function
			csID := fmt.Sprintf("cs_%d", i)
			addCallSite(cpg, fnID, csID, fmt.Sprintf("func%d", i+1))
			addBlock(cpg, fnID, fmt.Sprintf("blk_%d", i), "entry", []string{fmt.Sprintf("p_%d", i), csID})
			cpg.AddEdge(dataFlowEdge(fmt.Sprintf("p_%d", i), csID, "passes_to"))
			cpg.AddEdge(&graph.Edge{
				From: csID, To: fmt.Sprintf("fn_%d", i+1),
				Kind: graph.EdgeCalls, Label: "calls", Confidence: graph.ConfidenceCertain,
			})
		}
	}

	// With maxDepth=3, the chain of 5 should be truncated
	te := NewTaintEngine(WithMaxDepth(3))
	edges := te.Run(cpg)

	// Should NOT find the deep sink because depth > 3
	for _, e := range edges {
		if e.Kind == graph.EdgeTaint && e.To == "cs_sink_deep" {
			t.Fatalf("should not find taint edge to deep sink with maxDepth=3, but found one")
		}
	}
}

// TestTaintRecursiveCycleDetection: self-recursive function should not produce
// duplicate taint edges for the same logical flow.
func TestTaintRecursiveCycleDetection(t *testing.T) {
	cpg := buildTestCPG()

	// Self-recursive function: param passes_to call site that calls itself,
	// and also passes_to a sink
	fn := addFunc(cpg, "fn_rec", "process")
	fn.Annotations["handles_user_input"] = true
	fn.ParamNames = []string{"data"}
	addParam(cpg, "fn_rec", "p_rec", "data")

	addCallSite(cpg, "fn_rec", "cs_self", "process")
	sinkCS := addCallSite(cpg, "fn_rec", "cs_sink_rec", "exec")
	sinkCS.Annotations["sec:subprocess_call"] = true

	addBlock(cpg, "fn_rec", "blk_rec", "entry", []string{"p_rec", "cs_self", "cs_sink_rec"})
	cpg.AddEdge(dataFlowEdge("p_rec", "cs_self", "passes_to"))
	cpg.AddEdge(dataFlowEdge("p_rec", "cs_sink_rec", "passes_to"))

	// Call edge: cs_self -> fn_rec (self-recursive)
	cpg.AddEdge(&graph.Edge{From: "cs_self", To: "fn_rec", Kind: graph.EdgeCalls, Label: "calls", Confidence: graph.ConfidenceCertain})

	te := NewTaintEngine()
	edges := te.Run(cpg)

	// Count taint edges to the sink:
	// - 1 from Phase A (intraprocedural: p_rec -> cs_sink_rec)
	// - 1 from Phase B (interprocedural: p_rec -> cs_self -> fn_rec summary -> cs_sink_rec at depth 1)
	// Cycle detection prevents additional edges from deeper recursive unfolding.
	// Without cycle detection, we'd get up to maxDepth (20) edges.
	var sinkEdges []*graph.Edge
	for _, e := range edges {
		if e.Kind == graph.EdgeTaint && e.To == "cs_sink_rec" {
			sinkEdges = append(sinkEdges, e)
		}
	}

	if len(sinkEdges) > 2 {
		t.Fatalf("expected at most 2 taint edges to sink (1 intra + 1 inter before cycle detection), got %d", len(sinkEdges))
	}
	if len(sinkEdges) == 0 {
		t.Fatal("expected at least 1 taint edge to sink, got 0")
	}
}

// TestTaintEngineOptions: verify functional options work correctly.
func TestTaintEngineOptions(t *testing.T) {
	te := NewTaintEngine(
		WithSources([]string{"custom_source"}),
		WithSinks([]string{"custom_sink"}),
		WithMaxPaths(50),
		WithMaxVisits(500),
		WithMaxDepth(5),
	)

	if len(te.sources) != 1 || te.sources[0] != "custom_source" {
		t.Errorf("expected custom source, got %v", te.sources)
	}
	if len(te.sinks) != 1 || te.sinks[0] != "custom_sink" {
		t.Errorf("expected custom sink, got %v", te.sinks)
	}
	if te.maxPaths != 50 {
		t.Errorf("expected maxPaths=50, got %d", te.maxPaths)
	}
	if te.maxVisits != 500 {
		t.Errorf("expected maxVisits=500, got %d", te.maxVisits)
	}
	if te.maxDepth != 5 {
		t.Errorf("expected maxDepth=5, got %d", te.maxDepth)
	}
}
