package parser

import (
	"os"
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

func parsePythonCFGSample(t *testing.T) *ParseResult {
	t.Helper()
	content, err := os.ReadFile("../../testdata/cfg_sample.py")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}
	p := NewPythonParser()
	result, err := p.ParseFile("testdata/cfg_sample.py", content)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}
	return result
}

func TestPythonCFGBasicBlocksCreated(t *testing.T) {
	result := parsePythonCFGSample(t)
	if len(result.BasicBlocks) == 0 {
		t.Fatal("expected BasicBlocks, got 0")
	}
	for _, bb := range result.BasicBlocks {
		if bb.Kind != graph.NodeBasicBlock {
			t.Errorf("block %q has Kind=%s, want NodeBasicBlock", bb.Name, bb.Kind)
		}
	}
}

func TestPythonCFGEdgeLabels(t *testing.T) {
	result := parsePythonCFGSample(t)
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

func TestPythonCFGBlockMembers(t *testing.T) {
	result := parsePythonCFGSample(t)
	totalMembers := 0
	for _, bb := range result.BasicBlocks {
		totalMembers += len(bb.Members)
	}
	if totalMembers == 0 {
		t.Error("expected at least some blocks to have B1 node IDs as members")
	}
}
