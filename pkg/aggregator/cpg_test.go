package aggregator

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestAggregateCPGsEmpty(t *testing.T) {
	dir := t.TempDir()
	platform, err := AggregateCPGs(dir)
	if err != nil {
		t.Fatal(err)
	}
	if platform.ComponentCount != 0 {
		t.Errorf("expected 0 components, got %d", platform.ComponentCount)
	}
}

func TestAggregateCPGsSingleComponent(t *testing.T) {
	dir := t.TempDir()
	compDir := filepath.Join(dir, "dashboard")
	os.MkdirAll(compDir, 0o755)

	cpg := map[string]interface{}{
		"schema_version": 3,
		"component":      "dashboard",
		"nodes": []interface{}{
			map[string]interface{}{
				"id":   "func1",
				"kind": "Function",
				"name": "handleRequest",
				"file": "main.go",
			},
			map[string]interface{}{
				"id":     "ep1",
				"kind":   "HTTPEndpoint",
				"name":   "GET /api/v1/dashboard",
				"method": "GET",
				"path":   "/api/v1/dashboard",
			},
		},
		"edges": []interface{}{
			map[string]interface{}{
				"from": "func1",
				"to":   "ep1",
				"kind": "Serves",
			},
		},
	}
	data, _ := json.MarshalIndent(cpg, "", "  ")
	os.WriteFile(filepath.Join(compDir, "code-graph.json"), data, 0o644)

	platform, err := AggregateCPGs(dir)
	if err != nil {
		t.Fatal(err)
	}
	if platform.ComponentCount != 1 {
		t.Errorf("expected 1 component, got %d", platform.ComponentCount)
	}
	if platform.TotalNodes != 2 {
		t.Errorf("expected 2 nodes, got %d", platform.TotalNodes)
	}
	if platform.TotalEdges < 1 {
		t.Errorf("expected at least 1 edge, got %d", platform.TotalEdges)
	}
}

func TestAggregateCPGsCrossComponent(t *testing.T) {
	dir := t.TempDir()

	// Component A: serves an HTTP endpoint
	aDir := filepath.Join(dir, "model-registry")
	os.MkdirAll(aDir, 0o755)
	cpgA := map[string]interface{}{
		"schema_version": 3,
		"component":      "model-registry",
		"nodes": []interface{}{
			map[string]interface{}{
				"id":     "ep1",
				"kind":   "HTTPEndpoint",
				"name":   "GET /api/model_registry/v1alpha3/registered_models",
				"method": "GET",
				"path":   "/api/model_registry/v1alpha3/registered_models",
			},
			map[string]interface{}{
				"id":   "crd1",
				"kind": "CRDDefinition",
				"name": "ModelRegistry",
			},
		},
		"edges": []interface{}{},
	}
	dataA, _ := json.MarshalIndent(cpgA, "", "  ")
	os.WriteFile(filepath.Join(aDir, "code-graph.json"), dataA, 0o644)

	// Component B: watches CRD from A
	bDir := filepath.Join(dir, "dashboard")
	os.MkdirAll(bDir, 0o755)
	cpgB := map[string]interface{}{
		"schema_version": 3,
		"component":      "dashboard",
		"nodes": []interface{}{
			map[string]interface{}{
				"id":   "watch1",
				"kind": "ControllerWatch",
				"name": "ModelRegistry",
				"gvk":  "modelregistry.opendatahub.io/v1alpha1/ModelRegistry",
			},
		},
		"edges": []interface{}{},
	}
	dataB, _ := json.MarshalIndent(cpgB, "", "  ")
	os.WriteFile(filepath.Join(bDir, "code-graph.json"), dataB, 0o644)

	platform, err := AggregateCPGs(dir)
	if err != nil {
		t.Fatal(err)
	}
	if platform.ComponentCount != 2 {
		t.Errorf("expected 2 components, got %d", platform.ComponentCount)
	}

	// Should detect shared-crd cross-link
	if len(platform.CrossLinks) == 0 {
		t.Error("expected cross-component links for shared CRD")
	}

	foundCRDLink := false
	for _, link := range platform.CrossLinks {
		if link.Type == "shared-crd" && link.From == "dashboard" && link.To == "model-registry" {
			foundCRDLink = true
		}
	}
	if !foundCRDLink {
		t.Error("expected shared-crd link from dashboard to model-registry")
	}
}

func TestAggregateCPGsBadDir(t *testing.T) {
	_, err := AggregateCPGs("/nonexistent/dir")
	if err == nil {
		t.Error("expected error for nonexistent directory")
	}
}

func TestMergeCPGsNodeIDPrefixed(t *testing.T) {
	snapshots := []CPGSnapshot{
		{
			Component: "a",
			Nodes: []map[string]interface{}{
				{"id": "f1", "kind": "Function", "name": "foo"},
			},
			Edges: []map[string]interface{}{},
		},
		{
			Component: "b",
			Nodes: []map[string]interface{}{
				{"id": "f1", "kind": "Function", "name": "bar"},
			},
			Edges: []map[string]interface{}{},
		},
	}

	platform := mergeCPGs(snapshots)

	// Both nodes should exist with different prefixed IDs
	if platform.TotalNodes != 2 {
		t.Errorf("expected 2 nodes, got %d", platform.TotalNodes)
	}

	ids := make(map[string]bool)
	for _, n := range platform.Nodes {
		if id, ok := n["id"].(string); ok {
			ids[id] = true
		}
	}
	if !ids["a::f1"] || !ids["b::f1"] {
		t.Errorf("expected prefixed IDs a::f1 and b::f1, got %v", ids)
	}
}
