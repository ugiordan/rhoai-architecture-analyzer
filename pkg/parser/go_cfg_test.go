package parser

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

// Helper functions for CFG testing

func filterCFEdges(edges []*graph.Edge) []*graph.Edge {
	var out []*graph.Edge
	for _, e := range edges {
		if e.Kind == graph.EdgeControlFlow {
			out = append(out, e)
		}
	}
	return out
}

func filterCFEdgesByLabel(edges []*graph.Edge, label string) []*graph.Edge {
	var out []*graph.Edge
	for _, e := range edges {
		if e.Kind == graph.EdgeControlFlow && e.Label == label {
			out = append(out, e)
		}
	}
	return out
}

func findBlockByName(blocks []*graph.Node, parentID, name string) *graph.Node {
	for _, b := range blocks {
		if b.Kind == graph.NodeBasicBlock && b.ParentID == parentID && b.Name == name {
			return b
		}
	}
	return nil
}

func findFunctionByName(fns []*graph.Node, name string) *graph.Node {
	for _, fn := range fns {
		if fn.Name == name {
			return fn
		}
	}
	return nil
}

func parseCFGSample(t *testing.T) *ParseResult {
	t.Helper()
	content, err := os.ReadFile("../../testdata/cfg_sample.go")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}
	p := NewGoParser()
	result, err := p.ParseFile("testdata/cfg_sample.go", content)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}
	return result
}

// Test 1: Verify blocks exist with correct Kind and ParentID
func TestGoCFGBasicBlocksCreated(t *testing.T) {
	result := parseCFGSample(t)

	if len(result.BasicBlocks) == 0 {
		t.Fatal("expected BasicBlocks, got 0")
	}

	for _, b := range result.BasicBlocks {
		if b.Kind != graph.NodeBasicBlock {
			t.Errorf("block %q has Kind=%s, want NodeBasicBlock", b.Name, b.Kind)
		}
		if b.ParentID == "" {
			t.Errorf("block %q has empty ParentID", b.Name)
		}
	}
}

// Test 2: Verify entry and exit edges exist
func TestGoCFGEntryAndExitEdges(t *testing.T) {
	result := parseCFGSample(t)

	entryEdges := filterCFEdgesByLabel(result.Edges, "entry")
	if len(entryEdges) == 0 {
		t.Fatal("expected entry edges, got 0")
	}

	exitEdges := filterCFEdgesByLabel(result.Edges, "exit")
	if len(exitEdges) == 0 {
		t.Fatal("expected exit edges, got 0")
	}

	// Verify entry edge goes from function to entry block
	for _, e := range entryEdges {
		if e.Kind != graph.EdgeControlFlow {
			t.Errorf("entry edge has Kind=%s, want EdgeControlFlow", e.Kind)
		}
	}
}

// Test 3: Verify true_branch and false_branch edges for IfElse function
func TestGoCFGIfElse(t *testing.T) {
	result := parseCFGSample(t)

	trueBranch := filterCFEdgesByLabel(result.Edges, "true_branch")
	if len(trueBranch) == 0 {
		t.Error("expected true_branch edges from IfElse function")
	}

	falseBranch := filterCFEdgesByLabel(result.Edges, "false_branch")
	if len(falseBranch) == 0 {
		t.Error("expected false_branch edges from IfElse function")
	}
}

// Test 4: Verify loop_back and loop_exit edges for ForLoop function
func TestGoCFGForLoop(t *testing.T) {
	result := parseCFGSample(t)

	loopBack := filterCFEdgesByLabel(result.Edges, "loop_back")
	if len(loopBack) == 0 {
		t.Error("expected loop_back edges from ForLoop function")
	}

	loopExit := filterCFEdgesByLabel(result.Edges, "loop_exit")
	if len(loopExit) == 0 {
		t.Error("expected loop_exit edges from ForLoop function")
	}
}

// Test 5: Verify fallthrough edges exist
func TestGoCFGFallthrough(t *testing.T) {
	result := parseCFGSample(t)

	fallthroughEdges := filterCFEdgesByLabel(result.Edges, "fallthrough")
	if len(fallthroughEdges) == 0 {
		t.Error("expected fallthrough edges")
	}
}

// Test 6: Verify minimal CFG for empty function
func TestGoCFGEmptyFunction(t *testing.T) {
	result := parseCFGSample(t)

	fn := findFunctionByName(result.Functions, "EmptyFunction")
	if fn == nil {
		t.Fatal("EmptyFunction not found")
	}

	// Should have entry and exit blocks
	entry := findBlockByName(result.BasicBlocks, fn.ID, "entry")
	if entry == nil {
		t.Error("EmptyFunction should have entry block")
	}

	exit := findBlockByName(result.BasicBlocks, fn.ID, "exit")
	if exit == nil {
		t.Error("EmptyFunction should have exit block")
	}
}

// Test 7: Verify no-branch function has entry block with members
func TestGoCFGLinearFunction(t *testing.T) {
	result := parseCFGSample(t)

	fn := findFunctionByName(result.Functions, "LinearFunction")
	if fn == nil {
		t.Fatal("LinearFunction not found")
	}

	entry := findBlockByName(result.BasicBlocks, fn.ID, "entry")
	if entry == nil {
		t.Fatal("LinearFunction should have entry block")
	}

	// Entry block should have members (the variable declarations and return)
	if len(entry.Members) == 0 {
		t.Error("LinearFunction entry block should have members")
	}
}

// Test 8: Verify nested control flow produces both loop and branch edges
func TestGoCFGNestedIfInFor(t *testing.T) {
	result := parseCFGSample(t)

	fn := findFunctionByName(result.Functions, "NestedIfInFor")
	if fn == nil {
		t.Fatal("NestedIfInFor not found")
	}

	// Should have both loop edges and branch edges
	cfEdges := filterCFEdges(result.Edges)
	var hasLoopBack, hasTrueBranch bool
	for _, e := range cfEdges {
		// Check if edge belongs to this function by checking From block's ParentID
		fromBlock := findBlockByID(result.BasicBlocks, e.From)
		if fromBlock == nil || fromBlock.ParentID != fn.ID {
			continue
		}
		if e.Label == "loop_back" {
			hasLoopBack = true
		}
		if e.Label == "true_branch" {
			hasTrueBranch = true
		}
	}

	if !hasLoopBack {
		t.Error("NestedIfInFor should have loop_back edge")
	}
	if !hasTrueBranch {
		t.Error("NestedIfInFor should have true_branch edge")
	}
}

func findBlockByID(blocks []*graph.Node, id string) *graph.Node {
	for _, b := range blocks {
		if b.ID == id {
			return b
		}
	}
	return nil
}

// Test 9: Verify all 7 labels appear
func TestGoCFGAllEdgeLabels(t *testing.T) {
	result := parseCFGSample(t)

	cfEdges := filterCFEdges(result.Edges)
	labels := make(map[string]bool)
	for _, e := range cfEdges {
		labels[e.Label] = true
	}

	expectedLabels := []string{
		"entry", "exit", "true_branch", "false_branch",
		"fallthrough", "loop_back", "loop_exit",
	}

	for _, label := range expectedLabels {
		if !labels[label] {
			t.Errorf("expected label %q not found in CFG edges", label)
		}
	}
}

// Test 10: Verify at least some blocks have node IDs as members
func TestGoCFGBlockMembers(t *testing.T) {
	result := parseCFGSample(t)

	hasMembers := false
	for _, b := range result.BasicBlocks {
		if len(b.Members) > 0 {
			hasMembers = true
			// Verify members are B1 node IDs (should be non-empty strings)
			for _, m := range b.Members {
				if m == "" {
					t.Errorf("block %q has empty member ID", b.Name)
				}
			}
		}
	}

	if !hasMembers {
		t.Error("expected at least some blocks to have members")
	}
}

// Test 11: Verify all control flow edges use ConfidenceCertain
func TestGoCFGEdgesConfidence(t *testing.T) {
	result := parseCFGSample(t)

	cfEdges := filterCFEdges(result.Edges)
	for _, e := range cfEdges {
		if e.Confidence != graph.ConfidenceCertain {
			t.Errorf("CF edge %s->%s has Confidence=%s, want CERTAIN",
				e.From, e.To, e.Confidence)
		}
	}
}

// Test 14: Verify MaxBlocksPerFunction limit is enforced
func TestGoCFGMaxBlocksLimit(t *testing.T) {
	// Generate a Go source file with many nested if statements to exceed MaxBlocksPerFunction
	var sb strings.Builder
	sb.WriteString("package testdata\n\nfunc ManyBlocks(x int) int {\n")
	// Each if/else creates ~3-4 blocks. 200 / 3 ≈ 67 ifs should exceed the limit.
	for i := 0; i < 80; i++ {
		sb.WriteString(fmt.Sprintf("\tif x > %d {\n\t\tx = x + 1\n\t} else {\n\t\tx = x - 1\n\t}\n", i))
	}
	sb.WriteString("\treturn x\n}\n")

	p := NewGoParser()
	result, err := p.ParseFile("testdata/many_blocks.go", []byte(sb.String()))
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	fn := findFunctionByName(result.Functions, "ManyBlocks")
	if fn == nil {
		t.Fatal("function ManyBlocks not found")
	}

	// Should have the truncation annotation
	if fn.Annotations == nil || !fn.Annotations["cfg:truncated"] {
		t.Error("expected cfg:truncated annotation on function with too many blocks")
	}

	// Should still have some blocks (the ones created before the limit)
	if len(result.BasicBlocks) == 0 {
		t.Error("expected some BasicBlocks even when truncated")
	}
}
