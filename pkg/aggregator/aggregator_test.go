package aggregator

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestAggregate_NoDir(t *testing.T) {
	_, err := Aggregate("/nonexistent/path")
	if err == nil {
		t.Error("expected error for nonexistent directory")
	}
}

func TestAggregate_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	result, err := Aggregate(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["platform"] != "OpenShift AI" {
		t.Error("platform should be 'OpenShift AI'")
	}
	if result["component_count"] != 0 {
		t.Error("component_count should be 0 for empty dir")
	}
}

func TestAggregate_SingleComponent(t *testing.T) {
	dir := t.TempDir()
	compDir := filepath.Join(dir, "comp-a")
	if err := os.MkdirAll(compDir, 0o755); err != nil {
		t.Fatal(err)
	}

	compData := map[string]interface{}{
		"component": "comp-a",
		"crds": []interface{}{
			map[string]interface{}{
				"kind":    "Notebook",
				"group":   "kubeflow.org",
				"version": "v1",
			},
		},
		"services": []interface{}{
			map[string]interface{}{
				"name": "svc-a",
				"type": "ClusterIP",
				"ports": []interface{}{
					map[string]interface{}{"port": 8080, "protocol": "TCP"},
				},
			},
		},
		"secrets_referenced": []interface{}{},
		"rbac": map[string]interface{}{
			"cluster_roles": []interface{}{
				map[string]interface{}{
					"name": "comp-a-role",
					"rules": []interface{}{
						map[string]interface{}{
							"resources": []interface{}{"pods"},
						},
					},
				},
			},
		},
		"dependencies": map[string]interface{}{
			"internal_odh": []interface{}{
				map[string]interface{}{
					"component":   "comp-b",
					"interaction": "watches CRDs",
				},
			},
		},
		"controller_watches": []interface{}{},
	}

	raw, _ := json.Marshal(compData)
	if err := os.WriteFile(filepath.Join(compDir, "component-architecture.json"), raw, 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := Aggregate(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["component_count"] != 1 {
		t.Errorf("expected component_count=1, got %v", result["component_count"])
	}

	crds, ok := result["crds"].([]interface{})
	if !ok || len(crds) != 1 {
		t.Error("expected 1 CRD")
	}

	services, ok := result["services"].([]interface{})
	if !ok || len(services) != 1 {
		t.Error("expected 1 service")
	}

	depGraph, ok := result["dependency_graph"].([]interface{})
	if !ok || len(depGraph) != 1 {
		t.Error("expected 1 dependency edge")
	}

	crdOwnership, ok := result["crd_ownership"].(map[string]interface{})
	if !ok {
		t.Fatal("crd_ownership should be a map")
	}
	if crdOwnership["Notebook"] != "comp-a" {
		t.Error("Notebook CRD should be owned by comp-a")
	}
}

func TestAggregate_MultipleComponents(t *testing.T) {
	dir := t.TempDir()

	// comp-a defines Notebook CRD
	compA := map[string]interface{}{
		"component": "comp-a",
		"crds": []interface{}{
			map[string]interface{}{"kind": "Notebook", "group": "kubeflow.org", "version": "v1"},
		},
		"services":           []interface{}{},
		"secrets_referenced": []interface{}{},
		"rbac":               map[string]interface{}{"cluster_roles": []interface{}{}},
		"dependencies":       map[string]interface{}{"internal_odh": []interface{}{}},
		"controller_watches": []interface{}{},
	}

	// comp-b watches Notebook (cross-component)
	compB := map[string]interface{}{
		"component":          "comp-b",
		"crds":               []interface{}{},
		"services":           []interface{}{},
		"secrets_referenced": []interface{}{},
		"rbac":               map[string]interface{}{"cluster_roles": []interface{}{}},
		"dependencies":       map[string]interface{}{"internal_odh": []interface{}{}},
		"controller_watches": []interface{}{
			map[string]interface{}{"gvk": "kubeflow.org/v1/Notebook", "type": "For"},
		},
	}

	for name, data := range map[string]map[string]interface{}{"comp-a": compA, "comp-b": compB} {
		d := filepath.Join(dir, name)
		os.MkdirAll(d, 0o755)
		raw, _ := json.Marshal(data)
		os.WriteFile(filepath.Join(d, "component-architecture.json"), raw, 0o644)
		_ = name
	}

	result, err := Aggregate(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["component_count"] != 2 {
		t.Errorf("expected 2 components, got %v", result["component_count"])
	}

	// Check cross-component watch creates a dependency
	depGraph, ok := result["dependency_graph"].([]interface{})
	if !ok {
		t.Fatal("dependency_graph should be a slice")
	}
	found := false
	for _, d := range depGraph {
		dm, ok := d.(map[string]interface{})
		if !ok {
			continue
		}
		if dm["from"] == "comp-b" && dm["to"] == "comp-a" && dm["type"] == "watches-crd:Notebook" {
			found = true
		}
	}
	if !found {
		t.Error("expected cross-component CRD watch dependency from comp-b to comp-a")
	}
}
