package diff_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/diff"
	"github.com/ugiordan/architecture-analyzer/pkg/graph"
	"github.com/ugiordan/architecture-analyzer/pkg/parser"
)

func TestIntegrationGoParserDiff(t *testing.T) {
	// Parse the base fixture
	baseContent, err := os.ReadFile("../../testdata/simple_http_server.go")
	if err != nil {
		t.Fatalf("Failed to read base fixture: %v", err)
	}

	p := parser.NewGoParser()
	baseResult, err := p.ParseFile("testdata/simple_http_server.go", baseContent)
	if err != nil {
		t.Fatalf("Base parse failed: %v", err)
	}

	baseSnap := resultToSnapshot(baseResult)

	// Parse the complexity sample (different file, different functions)
	headContent, err := os.ReadFile("../../testdata/complexity_sample.go")
	if err != nil {
		t.Fatalf("Failed to read head fixture: %v", err)
	}

	p2 := parser.NewGoParser()
	headResult, err := p2.ParseFile("testdata/complexity_sample.go", headContent)
	if err != nil {
		t.Fatalf("Head parse failed: %v", err)
	}

	headSnap := resultToSnapshot(headResult)

	// Diff them
	d, err := diff.Compare(baseSnap, headSnap)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}

	// Different files mean all base nodes are removed, all head nodes are added
	if len(d.Nodes.Removed) == 0 {
		t.Error("expected removed nodes from base fixture")
	}
	if len(d.Nodes.Added) == 0 {
		t.Error("expected added nodes from head fixture")
	}

	// Verify JSON serialization works
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty JSON output")
	}

	// Verify summary is consistent
	if d.Summary.NodesAdded != len(d.Nodes.Added) {
		t.Errorf("summary NodesAdded=%d but len(Added)=%d", d.Summary.NodesAdded, len(d.Nodes.Added))
	}
	if d.Summary.NodesRemoved != len(d.Nodes.Removed) {
		t.Errorf("summary NodesRemoved=%d but len(Removed)=%d", d.Summary.NodesRemoved, len(d.Nodes.Removed))
	}
}

func TestIntegrationSameFileDiff(t *testing.T) {
	content, err := os.ReadFile("../../testdata/simple_http_server.go")
	if err != nil {
		t.Fatalf("Failed to read fixture: %v", err)
	}

	p1 := parser.NewGoParser()
	r1, err := p1.ParseFile("testdata/simple_http_server.go", content)
	if err != nil {
		t.Fatalf("Parse 1 failed: %v", err)
	}

	p2 := parser.NewGoParser()
	r2, err := p2.ParseFile("testdata/simple_http_server.go", content)
	if err != nil {
		t.Fatalf("Parse 2 failed: %v", err)
	}

	s1 := resultToSnapshot(r1)
	s2 := resultToSnapshot(r2)

	d, err := diff.Compare(s1, s2)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}

	if d.HasDifferences() {
		t.Errorf("same file parsed twice should produce empty diff, got: added=%d removed=%d modified=%d",
			d.Summary.NodesAdded, d.Summary.NodesRemoved, d.Summary.NodesModified)
	}
}

func resultToSnapshot(r *parser.ParseResult) diff.GraphSnapshot {
	var nodes []graph.Node
	for _, n := range r.Functions {
		nodes = append(nodes, *n)
	}
	for _, n := range r.CallSites {
		nodes = append(nodes, *n)
	}
	for _, n := range r.HTTPHandlers {
		nodes = append(nodes, *n)
	}
	for _, n := range r.DBOperations {
		nodes = append(nodes, *n)
	}
	for _, n := range r.StructLiterals {
		nodes = append(nodes, *n)
	}

	var edges []graph.Edge
	for _, e := range r.Edges {
		edges = append(edges, *e)
	}

	return diff.GraphSnapshot{
		SchemaVersion: 3,
		Nodes:         nodes,
		Edges:         edges,
	}
}
