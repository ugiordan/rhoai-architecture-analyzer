package parser

import (
	"os"
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

func parseRustCFGSample(t *testing.T) *ParseResult {
	t.Helper()
	content, err := os.ReadFile("../../testdata/cfg_sample.rs")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}
	p := NewRustParser()
	result, err := p.ParseFile("testdata/cfg_sample.rs", content)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}
	return result
}

func TestRustCFGBasicBlocksCreated(t *testing.T) {
	result := parseRustCFGSample(t)

	if len(result.BasicBlocks) == 0 {
		t.Fatal("expected BasicBlocks, got 0")
	}

	for _, bb := range result.BasicBlocks {
		if bb.Kind != graph.NodeBasicBlock {
			t.Errorf("block %q has Kind=%s, want NodeBasicBlock", bb.Name, bb.Kind)
		}
	}
}

func TestRustCFGEdgeLabels(t *testing.T) {
	result := parseRustCFGSample(t)

	cfEdges := filterCFEdges(result.Edges)
	labels := make(map[string]bool)
	for _, e := range cfEdges {
		labels[e.Label] = true
	}

	expected := []string{"entry", "exit", "true_branch", "false_branch", "fallthrough", "loop_back", "loop_exit"}
	for _, label := range expected {
		if !labels[label] {
			t.Errorf("missing EdgeControlFlow label %q", label)
		}
	}
}

func TestRustCFGBlockMembers(t *testing.T) {
	result := parseRustCFGSample(t)

	totalMembers := 0
	for _, bb := range result.BasicBlocks {
		totalMembers += len(bb.Members)
	}

	if totalMembers == 0 {
		t.Error("expected at least some blocks to have B1 node IDs as members")
	}
}

func TestRustCFGImplicitReturn(t *testing.T) {
	result := parseRustCFGSample(t)

	fn := findFunctionByName(result.Functions, "linear_function")
	if fn == nil {
		t.Fatal("function linear_function not found")
	}

	// Should have exit edge from entry block (linear, no branches)
	cfEdges := filterCFEdges(result.Edges)
	hasExit := false
	for _, e := range cfEdges {
		if e.Label == "exit" {
			for _, bb := range result.BasicBlocks {
				if bb.ID == e.From && bb.ParentID == fn.ID {
					hasExit = true
					break
				}
			}
		}
	}

	if !hasExit {
		t.Error("linear_function: expected exit edge from entry block")
	}
}
