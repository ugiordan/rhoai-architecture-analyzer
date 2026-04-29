package diff

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

func makeSnapshot(version int, nodes []graph.Node, edges []graph.Edge) GraphSnapshot {
	return GraphSnapshot{
		SchemaVersion: version,
		Nodes:         nodes,
		Edges:         edges,
	}
}

func TestIdenticalGraphsEmptyDiff(t *testing.T) {
	nodes := []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo", File: "a.go", Line: 1},
	}
	base := makeSnapshot(3, nodes, nil)
	head := makeSnapshot(3, nodes, nil)

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	if len(d.Nodes.Added) != 0 {
		t.Errorf("expected 0 added, got %d", len(d.Nodes.Added))
	}
	if len(d.Nodes.Removed) != 0 {
		t.Errorf("expected 0 removed, got %d", len(d.Nodes.Removed))
	}
	if len(d.Nodes.Modified) != 0 {
		t.Errorf("expected 0 modified, got %d", len(d.Nodes.Modified))
	}
	if d.Summary.NodesAdded != 0 || d.Summary.NodesRemoved != 0 || d.Summary.NodesModified != 0 {
		t.Errorf("summary should be all zeros: %+v", d.Summary)
	}
}

func TestNodeAdded(t *testing.T) {
	base := makeSnapshot(3, nil, nil)
	head := makeSnapshot(3, []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo", File: "a.go", Line: 1, Language: "go"},
	}, nil)

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	if len(d.Nodes.Added) != 1 {
		t.Fatalf("expected 1 added, got %d", len(d.Nodes.Added))
	}
	if d.Nodes.Added[0].Name != "foo" {
		t.Errorf("expected added node 'foo', got %q", d.Nodes.Added[0].Name)
	}
	if d.Summary.NodesAdded != 1 {
		t.Errorf("summary NodesAdded should be 1, got %d", d.Summary.NodesAdded)
	}
}

func TestNodeRemoved(t *testing.T) {
	base := makeSnapshot(3, []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo", File: "a.go", Line: 1, Language: "go"},
	}, nil)
	head := makeSnapshot(3, nil, nil)

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	if len(d.Nodes.Removed) != 1 {
		t.Fatalf("expected 1 removed, got %d", len(d.Nodes.Removed))
	}
	if d.Summary.NodesRemoved != 1 {
		t.Errorf("summary NodesRemoved should be 1, got %d", d.Summary.NodesRemoved)
	}
}

func TestNodeModified(t *testing.T) {
	base := makeSnapshot(3, []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo", File: "a.go", Line: 1, Complexity: 5},
	}, nil)
	head := makeSnapshot(3, []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo", File: "a.go", Line: 1, Complexity: 12},
	}, nil)

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	if len(d.Nodes.Modified) != 1 {
		t.Fatalf("expected 1 modified, got %d", len(d.Nodes.Modified))
	}
	mc := d.Nodes.Modified[0]
	if mc.ID != "fn_aaa" {
		t.Errorf("expected modified ID fn_aaa, got %q", mc.ID)
	}
	found := false
	for _, fc := range mc.Changes {
		if fc.Field == "complexity" {
			found = true
			if fc.OldValue != 5 || fc.NewValue != 12 {
				t.Errorf("expected complexity 5->12, got %v->%v", fc.OldValue, fc.NewValue)
			}
		}
	}
	if !found {
		t.Error("expected FieldChange for complexity")
	}
}

func TestEdgeAddedRemoved(t *testing.T) {
	baseEdges := []graph.Edge{
		{From: "fn_aaa", To: "call_bbb", Kind: graph.EdgeCalls, Label: "doStuff"},
	}
	headEdges := []graph.Edge{
		{From: "fn_aaa", To: "call_ccc", Kind: graph.EdgeCalls, Label: "doOther"},
	}
	nodes := []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo"},
	}
	base := makeSnapshot(3, nodes, baseEdges)
	head := makeSnapshot(3, nodes, headEdges)

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	if len(d.Edges.Added) != 1 {
		t.Errorf("expected 1 edge added, got %d", len(d.Edges.Added))
	}
	if len(d.Edges.Removed) != 1 {
		t.Errorf("expected 1 edge removed, got %d", len(d.Edges.Removed))
	}
}

func TestSchemaVersionTooOld(t *testing.T) {
	base := makeSnapshot(2, nil, nil)
	head := makeSnapshot(3, nil, nil)

	_, err := Compare(base, head)
	if err == nil {
		t.Fatal("expected error for schema_version < 3")
	}
}

func TestDuplicateIDInInput(t *testing.T) {
	nodes := []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo"},
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "bar"},
	}
	base := makeSnapshot(3, nodes, nil)
	head := makeSnapshot(3, nil, nil)

	_, err := Compare(base, head)
	if err == nil {
		t.Fatal("expected error for duplicate ID")
	}
}

func TestEmptyVsPopulated(t *testing.T) {
	nodes := []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo", Language: "go"},
		{ID: "fn_bbb", Kind: graph.NodeFunction, Name: "bar", Language: "go"},
	}
	base := makeSnapshot(3, nil, nil)
	head := makeSnapshot(3, nodes, nil)

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	if len(d.Nodes.Added) != 2 {
		t.Errorf("expected 2 added, got %d", len(d.Nodes.Added))
	}
	if d.Summary.NodesAdded != 2 {
		t.Errorf("summary NodesAdded should be 2, got %d", d.Summary.NodesAdded)
	}
}

func TestAnnotationMapChanges(t *testing.T) {
	base := makeSnapshot(3, []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo", Annotations: map[string]bool{"auth": true}},
	}, nil)
	head := makeSnapshot(3, []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo", Annotations: map[string]bool{"auth": false, "new_key": true}},
	}, nil)

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	if len(d.Nodes.Modified) != 1 {
		t.Fatalf("expected 1 modified, got %d", len(d.Nodes.Modified))
	}
	found := false
	for _, fc := range d.Nodes.Modified[0].Changes {
		if fc.Field == "annotations" {
			found = true
		}
	}
	if !found {
		t.Error("expected FieldChange for annotations")
	}
}

func TestSliceComparisonOrderIndependent(t *testing.T) {
	base := makeSnapshot(3, []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo", ParamTypes: []string{"int", "string"}},
	}, nil)
	head := makeSnapshot(3, []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo", ParamTypes: []string{"string", "int"}},
	}, nil)

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	if len(d.Nodes.Modified) != 0 {
		t.Errorf("expected 0 modified (order-independent), got %d", len(d.Nodes.Modified))
	}
}

func TestSummaryByKind(t *testing.T) {
	base := makeSnapshot(3, nil, nil)
	head := makeSnapshot(3, []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo", Language: "go"},
		{ID: "http_bbb", Kind: graph.NodeHTTPEndpoint, Name: "bar", Language: "go"},
	}, nil)

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	fnCounts, ok := d.Summary.ByKind["Function"]
	if !ok {
		t.Fatal("expected Function in ByKind")
	}
	if fnCounts.Added != 1 {
		t.Errorf("expected Function added=1, got %d", fnCounts.Added)
	}
}

func TestGraphDiffJSONRoundTrip(t *testing.T) {
	base := makeSnapshot(3, []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo", Complexity: 5},
	}, nil)
	head := makeSnapshot(3, []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo", Complexity: 10},
		{ID: "fn_bbb", Kind: graph.NodeFunction, Name: "bar"},
	}, nil)

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}

	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var d2 GraphDiff
	if err := json.Unmarshal(data, &d2); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if d2.Summary.NodesAdded != 1 || d2.Summary.NodesModified != 1 {
		t.Errorf("round-trip summary mismatch: %+v", d2.Summary)
	}
}

// --- Edge cases from spec ---

func TestPopulatedVsEmpty(t *testing.T) {
	nodes := []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo", Language: "go"},
		{ID: "fn_bbb", Kind: graph.NodeFunction, Name: "bar", Language: "go"},
		{ID: "call_ccc", Kind: graph.NodeCallSite, Name: "baz", Language: "go"},
	}
	base := makeSnapshot(3, nodes, nil)
	head := makeSnapshot(3, nil, nil)

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	if len(d.Nodes.Removed) != 3 {
		t.Errorf("expected 3 removed, got %d", len(d.Nodes.Removed))
	}
	if len(d.Nodes.Added) != 0 {
		t.Errorf("expected 0 added, got %d", len(d.Nodes.Added))
	}
	if d.Summary.NodesRemoved != 3 {
		t.Errorf("summary NodesRemoved should be 3, got %d", d.Summary.NodesRemoved)
	}
}

func TestNodeEveryFieldChanged(t *testing.T) {
	base := makeSnapshot(3, []graph.Node{
		{
			ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo", File: "a.go", Line: 1,
			EndLine: 10, Language: "go", TypeName: "Foo", Complexity: 3,
			ReturnType: "int", IsTest: false, IsUnsafe: false, IsExtern: false,
			CallTarget: "", IsMacro: false, Route: "", HTTPMethod: "",
			Operation: "", Table: "", StructType: "", TrustLevel: "",
			Decorators: []string{"old_deco"},
			ParamNames: []string{"a"},
			ParamTypes: []string{"int"},
			FieldNames: []string{"x"},
			Annotations: map[string]bool{"auth": true},
			Properties:  map[string]string{"key": "val"},
		},
	}, nil)
	head := makeSnapshot(3, []graph.Node{
		{
			ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo", File: "a.go", Line: 1,
			EndLine: 20, Language: "python", TypeName: "Bar", Complexity: 15,
			ReturnType: "string", IsTest: true, IsUnsafe: true, IsExtern: true,
			CallTarget: "remote", IsMacro: true, Route: "/api", HTTPMethod: "POST",
			Operation: "write", Table: "users", StructType: "MyStruct", TrustLevel: "untrusted",
			Decorators: []string{"new_deco"},
			ParamNames: []string{"b"},
			ParamTypes: []string{"string"},
			FieldNames: []string{"y"},
			Annotations: map[string]bool{"auth": false},
			Properties:  map[string]string{"key": "new_val"},
		},
	}, nil)

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	if len(d.Nodes.Modified) != 1 {
		t.Fatalf("expected 1 modified, got %d", len(d.Nodes.Modified))
	}
	mc := d.Nodes.Modified[0]
	// Should have changes for: end_line, language, type_name, complexity, return_type,
	// is_test, is_unsafe, is_extern, call_target, is_macro, route, http_method,
	// operation, table, struct_type, trust_level, decorators, param_names, param_types,
	// field_names, annotations, properties = 22 fields
	if len(mc.Changes) < 20 {
		t.Errorf("expected at least 20 field changes, got %d", len(mc.Changes))
		for _, fc := range mc.Changes {
			t.Logf("  changed: %s", fc.Field)
		}
	}
}

func TestNodeSameIDIdenticalContentNotModified(t *testing.T) {
	node := graph.Node{
		ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo", File: "a.go", Line: 1,
		Complexity: 5, Language: "go", TrustLevel: "trusted",
		Annotations: map[string]bool{"auth": true},
		Properties:  map[string]string{"key": "val"},
		ParamTypes:  []string{"int", "string"},
	}
	base := makeSnapshot(3, []graph.Node{node}, nil)
	head := makeSnapshot(3, []graph.Node{node}, nil)

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	if len(d.Nodes.Modified) != 0 {
		t.Errorf("identical nodes should not appear as modified, got %d changes", len(d.Nodes.Modified))
		for _, mc := range d.Nodes.Modified {
			for _, fc := range mc.Changes {
				t.Logf("  unexpected change: %s %v -> %v", fc.Field, fc.OldValue, fc.NewValue)
			}
		}
	}
}

func TestAnnotationDistinctChanges(t *testing.T) {
	base := makeSnapshot(3, []graph.Node{
		{
			ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo",
			Annotations: map[string]bool{"kept": true, "removed_key": true, "flipped": true},
		},
	}, nil)
	head := makeSnapshot(3, []graph.Node{
		{
			ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo",
			Annotations: map[string]bool{"kept": true, "added_key": true, "flipped": false},
		},
	}, nil)

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	if len(d.Nodes.Modified) != 1 {
		t.Fatalf("expected 1 modified, got %d", len(d.Nodes.Modified))
	}
	found := false
	for _, fc := range d.Nodes.Modified[0].Changes {
		if fc.Field == "annotations" {
			found = true
			// The old and new values should be the full maps
			oldMap, ok1 := fc.OldValue.(map[string]bool)
			newMap, ok2 := fc.NewValue.(map[string]bool)
			if !ok1 || !ok2 {
				t.Error("annotations change values should be map[string]bool")
				break
			}
			// removed_key should be absent from new
			if _, exists := newMap["removed_key"]; exists {
				t.Error("removed_key should be absent from new annotations")
			}
			// added_key should be absent from old
			if _, exists := oldMap["added_key"]; exists {
				t.Error("added_key should be absent from old annotations")
			}
			// flipped should be true in old, false in new
			if !oldMap["flipped"] {
				t.Error("flipped should be true in old")
			}
			if newMap["flipped"] {
				t.Error("flipped should be false in new")
			}
		}
	}
	if !found {
		t.Error("expected FieldChange for annotations")
	}
}

func TestPropertiesMapChanges(t *testing.T) {
	base := makeSnapshot(3, []graph.Node{
		{
			ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo",
			Properties: map[string]string{"method": "GET", "old_prop": "val"},
		},
	}, nil)
	head := makeSnapshot(3, []graph.Node{
		{
			ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo",
			Properties: map[string]string{"method": "POST", "new_prop": "val"},
		},
	}, nil)

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	if len(d.Nodes.Modified) != 1 {
		t.Fatalf("expected 1 modified, got %d", len(d.Nodes.Modified))
	}
	found := false
	for _, fc := range d.Nodes.Modified[0].Changes {
		if fc.Field == "properties" {
			found = true
		}
	}
	if !found {
		t.Error("expected FieldChange for properties")
	}
}

func TestDuplicateIDInHeadSnapshot(t *testing.T) {
	base := makeSnapshot(3, nil, nil)
	head := makeSnapshot(3, []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo"},
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "bar"},
	}, nil)

	_, err := Compare(base, head)
	if err == nil {
		t.Fatal("expected error for duplicate ID in head snapshot")
	}
}

func TestBothSnapshotsEmpty(t *testing.T) {
	base := makeSnapshot(3, nil, nil)
	head := makeSnapshot(3, nil, nil)

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	if d.HasDifferences() {
		t.Error("two empty snapshots should have no differences")
	}
	if d.Summary.NodesAdded != 0 || d.Summary.NodesRemoved != 0 ||
		d.Summary.NodesModified != 0 || d.Summary.EdgesAdded != 0 ||
		d.Summary.EdgesRemoved != 0 {
		t.Errorf("all summary counts should be zero: %+v", d.Summary)
	}
}

func TestSchemaVersionHeadTooOld(t *testing.T) {
	base := makeSnapshot(3, nil, nil)
	head := makeSnapshot(1, nil, nil)

	_, err := Compare(base, head)
	if err == nil {
		t.Fatal("expected error for head schema_version < 3")
	}
}

func TestSchemaVersionBothTooOld(t *testing.T) {
	base := makeSnapshot(2, nil, nil)
	head := makeSnapshot(2, nil, nil)

	_, err := Compare(base, head)
	if err == nil {
		t.Fatal("expected error when both schema_version < 3")
	}
}

func TestEdgeSameEndpointsDifferentKind(t *testing.T) {
	base := makeSnapshot(3, []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo"},
		{ID: "fn_bbb", Kind: graph.NodeFunction, Name: "bar"},
	}, []graph.Edge{
		{From: "fn_aaa", To: "fn_bbb", Kind: graph.EdgeCalls, Label: "bar"},
	})
	head := makeSnapshot(3, []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo"},
		{ID: "fn_bbb", Kind: graph.NodeFunction, Name: "bar"},
	}, []graph.Edge{
		{From: "fn_aaa", To: "fn_bbb", Kind: graph.EdgeDataFlow, Label: "bar"},
	})

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	// Same From/To/Label but different Kind means the CALLS edge is removed and DATA_FLOW is added
	if len(d.Edges.Added) != 1 {
		t.Errorf("expected 1 edge added, got %d", len(d.Edges.Added))
	}
	if len(d.Edges.Removed) != 1 {
		t.Errorf("expected 1 edge removed, got %d", len(d.Edges.Removed))
	}
	if len(d.Edges.Added) > 0 && d.Edges.Added[0].Kind != graph.EdgeDataFlow {
		t.Errorf("added edge should be DATA_FLOW, got %s", d.Edges.Added[0].Kind)
	}
	if len(d.Edges.Removed) > 0 && d.Edges.Removed[0].Kind != graph.EdgeCalls {
		t.Errorf("removed edge should be CALLS, got %s", d.Edges.Removed[0].Kind)
	}
}

func TestHasDifferencesTrue(t *testing.T) {
	base := makeSnapshot(3, nil, nil)
	head := makeSnapshot(3, []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo"},
	}, nil)

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	if !d.HasDifferences() {
		t.Error("HasDifferences should be true when nodes are added")
	}
}

func TestHasDifferencesEdgesOnly(t *testing.T) {
	nodes := []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo"},
		{ID: "fn_bbb", Kind: graph.NodeFunction, Name: "bar"},
	}
	base := makeSnapshot(3, nodes, nil)
	head := makeSnapshot(3, nodes, []graph.Edge{
		{From: "fn_aaa", To: "fn_bbb", Kind: graph.EdgeCalls, Label: "bar"},
	})

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	if !d.HasDifferences() {
		t.Error("HasDifferences should be true when only edges differ")
	}
}

func TestSummaryByLanguage(t *testing.T) {
	base := makeSnapshot(3, nil, nil)
	head := makeSnapshot(3, []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo", Language: "go"},
		{ID: "fn_bbb", Kind: graph.NodeFunction, Name: "bar", Language: "python"},
		{ID: "fn_ccc", Kind: graph.NodeFunction, Name: "baz", Language: "go"},
	}, nil)

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	goCounts, ok := d.Summary.ByLanguage["go"]
	if !ok {
		t.Fatal("expected 'go' in ByLanguage")
	}
	if goCounts.Added != 2 {
		t.Errorf("expected go added=2, got %d", goCounts.Added)
	}
	pyCounts, ok := d.Summary.ByLanguage["python"]
	if !ok {
		t.Fatal("expected 'python' in ByLanguage")
	}
	if pyCounts.Added != 1 {
		t.Errorf("expected python added=1, got %d", pyCounts.Added)
	}
}

func TestSummaryMixedOperations(t *testing.T) {
	base := makeSnapshot(3, []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "removed_fn", Language: "go"},
		{ID: "fn_bbb", Kind: graph.NodeFunction, Name: "modified_fn", Language: "go", Complexity: 3},
	}, nil)
	head := makeSnapshot(3, []graph.Node{
		{ID: "fn_bbb", Kind: graph.NodeFunction, Name: "modified_fn", Language: "go", Complexity: 15},
		{ID: "fn_ccc", Kind: graph.NodeFunction, Name: "added_fn", Language: "go"},
	}, nil)

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	fnCounts := d.Summary.ByKind["Function"]
	if fnCounts.Added != 1 {
		t.Errorf("expected Function added=1, got %d", fnCounts.Added)
	}
	if fnCounts.Removed != 1 {
		t.Errorf("expected Function removed=1, got %d", fnCounts.Removed)
	}
	if fnCounts.Modified != 1 {
		t.Errorf("expected Function modified=1, got %d", fnCounts.Modified)
	}
}

func TestOutputDeterministic(t *testing.T) {
	// Run Compare multiple times, verify JSON output is identical
	base := makeSnapshot(3, []graph.Node{
		{ID: "fn_zzz", Kind: graph.NodeFunction, Name: "z"},
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "a"},
	}, nil)
	head := makeSnapshot(3, []graph.Node{
		{ID: "fn_mmm", Kind: graph.NodeFunction, Name: "m"},
		{ID: "fn_bbb", Kind: graph.NodeFunction, Name: "b"},
	}, nil)

	var outputs []string
	for i := 0; i < 5; i++ {
		d, err := Compare(base, head)
		if err != nil {
			t.Fatalf("Compare failed on iteration %d: %v", i, err)
		}
		data, err := json.Marshal(d)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}
		outputs = append(outputs, string(data))
	}

	for i := 1; i < len(outputs); i++ {
		if outputs[i] != outputs[0] {
			t.Errorf("non-deterministic output on iteration %d", i)
		}
	}
}

func TestSliceFieldActualDifference(t *testing.T) {
	// Slices with genuinely different content (not just order)
	base := makeSnapshot(3, []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo", ParamTypes: []string{"int", "string"}},
	}, nil)
	head := makeSnapshot(3, []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo", ParamTypes: []string{"int", "bool"}},
	}, nil)

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	if len(d.Nodes.Modified) != 1 {
		t.Fatalf("expected 1 modified, got %d", len(d.Nodes.Modified))
	}
	found := false
	for _, fc := range d.Nodes.Modified[0].Changes {
		if fc.Field == "param_types" {
			found = true
		}
	}
	if !found {
		t.Error("expected FieldChange for param_types when actual content differs")
	}
}

func TestSliceEmptyVsNil(t *testing.T) {
	// Empty slice and nil should be treated as equal
	base := makeSnapshot(3, []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo", ParamTypes: nil},
	}, nil)
	head := makeSnapshot(3, []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo", ParamTypes: []string{}},
	}, nil)

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	if len(d.Nodes.Modified) != 0 {
		t.Errorf("nil and empty slice should be equal, got %d modifications", len(d.Nodes.Modified))
	}
}

func TestManyNodesPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping benchmark-style test in short mode")
	}

	const n = 10000
	baseNodes := make([]graph.Node, n)
	headNodes := make([]graph.Node, n)
	for i := 0; i < n; i++ {
		id := fmt.Sprintf("fn_%016x", i)
		baseNodes[i] = graph.Node{
			ID: id, Kind: graph.NodeFunction,
			Name: fmt.Sprintf("func_%d", i), File: "big.go", Line: i + 1, Language: "go",
		}
		headNodes[i] = baseNodes[i]
	}
	// Modify first 100 nodes
	for i := 0; i < 100; i++ {
		headNodes[i].Complexity = 99
	}
	// Remove last 100 original nodes (indices n-100..n-1)
	headNodes = headNodes[:n-100]
	// Add 100 new nodes
	for i := 0; i < 100; i++ {
		headNodes = append(headNodes, graph.Node{
			ID: fmt.Sprintf("fn_new_%016x", i), Kind: graph.NodeFunction,
			Name: fmt.Sprintf("new_func_%d", i), File: "big.go", Line: n + i + 1, Language: "go",
		})
	}

	base := makeSnapshot(3, baseNodes, nil)
	head := makeSnapshot(3, headNodes, nil)

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare failed with %d nodes: %v", n, err)
	}
	if d.Summary.NodesModified != 100 {
		t.Errorf("expected 100 modified, got %d", d.Summary.NodesModified)
	}
	if d.Summary.NodesAdded != 100 {
		t.Errorf("expected 100 added, got %d", d.Summary.NodesAdded)
	}
	if d.Summary.NodesRemoved != 100 {
		t.Errorf("expected 100 removed, got %d", d.Summary.NodesRemoved)
	}
}

func TestEdgeIdenticalNotReported(t *testing.T) {
	nodes := []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo"},
		{ID: "fn_bbb", Kind: graph.NodeFunction, Name: "bar"},
	}
	edges := []graph.Edge{
		{From: "fn_aaa", To: "fn_bbb", Kind: graph.EdgeCalls, Label: "bar"},
	}
	base := makeSnapshot(3, nodes, edges)
	head := makeSnapshot(3, nodes, edges)

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	if len(d.Edges.Added) != 0 || len(d.Edges.Removed) != 0 {
		t.Errorf("identical edges should not appear as added/removed: added=%d removed=%d",
			len(d.Edges.Added), len(d.Edges.Removed))
	}
}

func TestMultipleEdgesSameNodes(t *testing.T) {
	nodes := []graph.Node{
		{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo"},
		{ID: "fn_bbb", Kind: graph.NodeFunction, Name: "bar"},
	}
	base := makeSnapshot(3, nodes, []graph.Edge{
		{From: "fn_aaa", To: "fn_bbb", Kind: graph.EdgeCalls, Label: "bar"},
		{From: "fn_aaa", To: "fn_bbb", Kind: graph.EdgeDataFlow, Label: "contains_call"},
	})
	head := makeSnapshot(3, nodes, []graph.Edge{
		{From: "fn_aaa", To: "fn_bbb", Kind: graph.EdgeCalls, Label: "bar"},
		// DATA_FLOW edge removed, new CONTROL_FLOW added
		{From: "fn_aaa", To: "fn_bbb", Kind: graph.EdgeControlFlow, Label: "branch"},
	})

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	if len(d.Edges.Added) != 1 {
		t.Errorf("expected 1 edge added, got %d", len(d.Edges.Added))
	}
	if len(d.Edges.Removed) != 1 {
		t.Errorf("expected 1 edge removed, got %d", len(d.Edges.Removed))
	}
}

func TestCompareExternalFindingFields(t *testing.T) {
	base := GraphSnapshot{
		SchemaVersion: 3,
		Nodes: []graph.Node{
			{
				ID: "extf_abc123", Kind: graph.NodeExternalFinding, Name: "xss-rule",
				File: "a.go", Line: 10, RuleID: "xss-rule", Severity: "warning",
				Message: "XSS detected", ToolName: "semgrep", ToolVersion: "1.0",
				CWEs: []string{"CWE-79"},
			},
		},
	}
	head := GraphSnapshot{
		SchemaVersion: 3,
		Nodes: []graph.Node{
			{
				ID: "extf_abc123", Kind: graph.NodeExternalFinding, Name: "xss-rule",
				File: "a.go", Line: 10, RuleID: "xss-rule", Severity: "error",
				Message: "XSS detected (updated)", ToolName: "semgrep", ToolVersion: "1.1",
				CWEs: []string{"CWE-79", "CWE-80"},
			},
		},
	}

	d, err := Compare(base, head)
	if err != nil {
		t.Fatalf("Compare error: %v", err)
	}
	if len(d.Nodes.Modified) != 1 {
		t.Fatalf("Modified = %d, want 1", len(d.Nodes.Modified))
	}

	mc := d.Nodes.Modified[0]
	changedFields := make(map[string]bool)
	for _, fc := range mc.Changes {
		changedFields[fc.Field] = true
	}

	for _, field := range []string{"severity", "message", "tool_version", "cwes"} {
		if !changedFields[field] {
			t.Errorf("expected change for field %q", field)
		}
	}
	// rule_id and tool_name are unchanged, should NOT appear
	for _, field := range []string{"rule_id", "tool_name"} {
		if changedFields[field] {
			t.Errorf("unexpected change for field %q (value unchanged)", field)
		}
	}
}
