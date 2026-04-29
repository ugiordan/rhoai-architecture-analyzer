package parser

import (
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

func TestNodeIDDeterministic(t *testing.T) {
	id1 := NodeID(graph.NodeFunction, "handleRequest", "server.go", 42, 0)
	id2 := NodeID(graph.NodeFunction, "handleRequest", "server.go", 42, 0)
	if id1 != id2 {
		t.Errorf("same inputs should produce same ID: got %q and %q", id1, id2)
	}
}

func TestNodeIDFormat(t *testing.T) {
	id := NodeID(graph.NodeFunction, "foo", "bar.go", 1, 0)
	if !strings.HasPrefix(id, "fn_") {
		t.Errorf("Function node ID should have fn_ prefix, got %q", id)
	}
	// fn_ prefix + 16 hex chars = 19 chars total
	parts := strings.SplitN(id, "_", 2)
	if len(parts) != 2 {
		t.Fatalf("expected prefix_hash format, got %q", id)
	}
	hash := parts[1]
	if len(hash) != 16 {
		t.Errorf("expected 16 hex chars, got %d: %q", len(hash), hash)
	}
	for _, c := range hash {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("non-hex character %c in hash %q", c, hash)
		}
	}
}

func TestNodeIDDifferentColumns(t *testing.T) {
	id1 := NodeID(graph.NodeCallSite, "foo", "bar.go", 10, 0)
	id2 := NodeID(graph.NodeCallSite, "foo", "bar.go", 10, 5)
	if id1 == id2 {
		t.Errorf("different columns should produce different IDs: both are %q", id1)
	}
}

func TestNodeIDDifferentKinds(t *testing.T) {
	id1 := NodeID(graph.NodeCallSite, "foo", "bar.go", 10, 0)
	id2 := NodeID(graph.NodeHTTPEndpoint, "foo", "bar.go", 10, 0)
	if id1 == id2 {
		t.Errorf("different kinds should produce different IDs: both are %q", id1)
	}
	// Prefixes should differ
	if strings.HasPrefix(id1, "call_") && strings.HasPrefix(id2, "call_") {
		t.Error("different kinds should have different prefixes")
	}
}

func TestNodeIDPrefixes(t *testing.T) {
	tests := []struct {
		kind   graph.NodeKind
		prefix string
	}{
		{graph.NodeFunction, "fn_"},
		{graph.NodeCallSite, "call_"},
		{graph.NodeHTTPEndpoint, "http_"},
		{graph.NodeDBOperation, "db_"},
		{graph.NodeStructLiteral, "struct_"},
		{graph.NodeK8sResource, "k8s_"},
		{graph.NodeExternalCall, "ext_"},
		{graph.NodeParameter, "param_"},
		{graph.NodeVariable, "var_"},
		{graph.NodeLiteral, "lit_"},
	}
	for _, tt := range tests {
		id := NodeID(tt.kind, "test", "test.go", 1, 0)
		if !strings.HasPrefix(id, tt.prefix) {
			t.Errorf("NodeKind %q: expected prefix %q, got ID %q", tt.kind, tt.prefix, id)
		}
	}
}

func collectIDs(r *ParseResult) []string {
	var ids []string
	for _, n := range r.Functions {
		ids = append(ids, n.ID)
	}
	for _, n := range r.CallSites {
		ids = append(ids, n.ID)
	}
	for _, n := range r.HTTPHandlers {
		ids = append(ids, n.ID)
	}
	for _, n := range r.DBOperations {
		ids = append(ids, n.ID)
	}
	for _, n := range r.StructLiterals {
		ids = append(ids, n.ID)
	}
	sort.Strings(ids)
	return ids
}

func TestPythonParserStableIDs(t *testing.T) {
	content, err := os.ReadFile("../../testdata/flask_app.py")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	p1 := NewPythonParser()
	r1, err := p1.ParseFile("testdata/flask_app.py", content)
	if err != nil {
		t.Fatalf("first parse failed: %v", err)
	}

	p2 := NewPythonParser()
	r2, err := p2.ParseFile("testdata/flask_app.py", content)
	if err != nil {
		t.Fatalf("second parse failed: %v", err)
	}

	ids1 := collectIDs(r1)
	ids2 := collectIDs(r2)

	if len(ids1) != len(ids2) {
		t.Fatalf("different number of IDs: %d vs %d", len(ids1), len(ids2))
	}
	for i, id := range ids1 {
		if id != ids2[i] {
			t.Errorf("ID mismatch at index %d: %q vs %q", i, id, ids2[i])
		}
	}
}

func TestTypeScriptParserStableIDs(t *testing.T) {
	content, err := os.ReadFile("../../testdata/express_server.ts")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	p1 := NewTypeScriptParser()
	r1, err := p1.ParseFile("testdata/express_server.ts", content)
	if err != nil {
		t.Fatalf("first parse failed: %v", err)
	}

	p2 := NewTypeScriptParser()
	r2, err := p2.ParseFile("testdata/express_server.ts", content)
	if err != nil {
		t.Fatalf("second parse failed: %v", err)
	}

	ids1 := collectIDs(r1)
	ids2 := collectIDs(r2)

	if len(ids1) != len(ids2) {
		t.Fatalf("different number of IDs: %d vs %d", len(ids1), len(ids2))
	}
	for i, id := range ids1 {
		if id != ids2[i] {
			t.Errorf("ID mismatch at index %d: %q vs %q", i, id, ids2[i])
		}
	}
}

func TestRustParserStableIDs(t *testing.T) {
	content, err := os.ReadFile("../../testdata/actix_handler.rs")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	p1 := NewRustParser()
	r1, err := p1.ParseFile("testdata/actix_handler.rs", content)
	if err != nil {
		t.Fatalf("first parse failed: %v", err)
	}

	p2 := NewRustParser()
	r2, err := p2.ParseFile("testdata/actix_handler.rs", content)
	if err != nil {
		t.Fatalf("second parse failed: %v", err)
	}

	ids1 := collectIDs(r1)
	ids2 := collectIDs(r2)

	if len(ids1) != len(ids2) {
		t.Fatalf("different number of IDs: %d vs %d", len(ids1), len(ids2))
	}
	for i, id := range ids1 {
		if id != ids2[i] {
			t.Errorf("ID mismatch at index %d: %q vs %q", i, id, ids2[i])
		}
	}
}

func TestGoParserStableIDs(t *testing.T) {
	content, err := os.ReadFile("../../testdata/simple_http_server.go")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	p1 := NewGoParser()
	r1, err := p1.ParseFile("testdata/simple_http_server.go", content)
	if err != nil {
		t.Fatalf("first parse failed: %v", err)
	}

	p2 := NewGoParser()
	r2, err := p2.ParseFile("testdata/simple_http_server.go", content)
	if err != nil {
		t.Fatalf("second parse failed: %v", err)
	}

	ids1 := collectIDs(r1)
	ids2 := collectIDs(r2)

	if len(ids1) != len(ids2) {
		t.Fatalf("different number of IDs: %d vs %d", len(ids1), len(ids2))
	}
	for i, id := range ids1 {
		if id != ids2[i] {
			t.Errorf("ID mismatch at index %d: %q vs %q", i, id, ids2[i])
		}
	}
}
