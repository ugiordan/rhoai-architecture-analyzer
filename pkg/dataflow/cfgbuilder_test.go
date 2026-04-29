package dataflow

import (
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

func TestCFGBuilderNewBlock(t *testing.T) {
	cb := NewCFGBuilder("fn_abc123", "main.go")

	entry := cb.NewBlock("entry", 10)
	if entry.Kind != graph.NodeBasicBlock {
		t.Errorf("Kind = %s, want BasicBlock", entry.Kind)
	}
	if entry.Name != "entry" {
		t.Errorf("Name = %q, want entry", entry.Name)
	}
	if entry.File != "main.go" {
		t.Errorf("File = %q, want main.go", entry.File)
	}
	if entry.Line != 10 {
		t.Errorf("Line = %d, want 10", entry.Line)
	}
	if entry.ParentID != "fn_abc123" {
		t.Errorf("ParentID = %q, want fn_abc123", entry.ParentID)
	}
	if entry.ID == "" {
		t.Error("ID should not be empty")
	}
}

func TestCFGBuilderAddMember(t *testing.T) {
	cb := NewCFGBuilder("fn_abc123", "main.go")
	block := cb.NewBlock("entry", 10)

	cb.AddMember(block, "var_001")
	cb.AddMember(block, "call_002")

	if len(block.Members) != 2 {
		t.Fatalf("got %d members, want 2", len(block.Members))
	}
	if block.Members[0] != "var_001" {
		t.Errorf("Members[0] = %q, want var_001", block.Members[0])
	}
	if block.Members[1] != "call_002" {
		t.Errorf("Members[1] = %q, want call_002", block.Members[1])
	}
}

func TestCFGBuilderAddEdge(t *testing.T) {
	cb := NewCFGBuilder("fn_abc123", "main.go")

	cb.AddEdge("block_a", "block_b", "true_branch")

	_, edges := cb.Result()
	if len(edges) != 1 {
		t.Fatalf("got %d edges, want 1", len(edges))
	}
	e := edges[0]
	if e.From != "block_a" || e.To != "block_b" {
		t.Errorf("edge = %s -> %s, want block_a -> block_b", e.From, e.To)
	}
	if e.Kind != graph.EdgeControlFlow {
		t.Errorf("kind = %s, want CONTROL_FLOW", e.Kind)
	}
	if e.Label != "true_branch" {
		t.Errorf("label = %q, want true_branch", e.Label)
	}
	if e.Confidence != graph.ConfidenceCertain {
		t.Errorf("confidence = %s, want CERTAIN", e.Confidence)
	}
}

func TestCFGBuilderResult(t *testing.T) {
	cb := NewCFGBuilder("fn_abc123", "main.go")

	entry := cb.NewBlock("entry", 10)
	exit := cb.NewBlock("exit", 0)
	cb.AddMember(entry, "var_001")
	cb.AddEdge(entry.ID, exit.ID, "exit")

	blocks, edges := cb.Result()
	if len(blocks) != 2 {
		t.Fatalf("got %d blocks, want 2", len(blocks))
	}
	if len(edges) != 1 {
		t.Fatalf("got %d edges, want 1", len(edges))
	}
}

func TestCFGBuilderBlockCount(t *testing.T) {
	cb := NewCFGBuilder("fn_abc123", "main.go")
	if cb.BlockCount() != 0 {
		t.Errorf("initial BlockCount = %d, want 0", cb.BlockCount())
	}
	cb.NewBlock("entry", 10)
	cb.NewBlock("bb0", 15)
	if cb.BlockCount() != 2 {
		t.Errorf("BlockCount = %d, want 2", cb.BlockCount())
	}
}

func TestCFGBuilderBlockIDUniqueness(t *testing.T) {
	cb1 := NewCFGBuilder("fn_aaa", "main.go")
	cb2 := NewCFGBuilder("fn_bbb", "main.go")

	exit1 := cb1.NewBlock("exit", 0)
	exit2 := cb2.NewBlock("exit", 0)

	if exit1.ID == exit2.ID {
		t.Errorf("exit blocks from different functions should have different IDs, both got %s", exit1.ID)
	}
}

func TestCFGBuilderEdgeLabels(t *testing.T) {
	cb := NewCFGBuilder("fn_abc", "main.go")
	labels := []string{"entry", "exit", "true_branch", "false_branch", "fallthrough", "loop_back", "loop_exit"}

	for _, label := range labels {
		cb.AddEdge("from", "to", label)
	}
	_, edges := cb.Result()
	if len(edges) != len(labels) {
		t.Fatalf("got %d edges, want %d", len(edges), len(labels))
	}
	for i, label := range labels {
		if edges[i].Label != label {
			t.Errorf("edge %d label = %q, want %q", i, edges[i].Label, label)
		}
	}
}
