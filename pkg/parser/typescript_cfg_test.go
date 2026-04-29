package parser

import (
	"os"
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

func parseTSCFGSample(t *testing.T) *ParseResult {
	t.Helper()
	content, err := os.ReadFile("../../testdata/cfg_sample.ts")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}
	p := NewTypeScriptParser()
	result, err := p.ParseFile("testdata/cfg_sample.ts", content)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}
	return result
}

func TestTSCFGBasicBlocksCreated(t *testing.T) {
	result := parseTSCFGSample(t)
	if len(result.BasicBlocks) == 0 {
		t.Fatal("expected BasicBlocks, got 0")
	}
	for _, bb := range result.BasicBlocks {
		if bb.Kind != graph.NodeBasicBlock {
			t.Errorf("block %q has Kind=%s, want NodeBasicBlock", bb.Name, bb.Kind)
		}
	}
}

func TestTSCFGEdgeLabels(t *testing.T) {
	result := parseTSCFGSample(t)
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

func TestTSCFGBlockMembers(t *testing.T) {
	result := parseTSCFGSample(t)
	totalMembers := 0
	for _, bb := range result.BasicBlocks {
		totalMembers += len(bb.Members)
	}
	if totalMembers == 0 {
		t.Error("expected at least some blocks to have B1 node IDs as members")
	}
}
