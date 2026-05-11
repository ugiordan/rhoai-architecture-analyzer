package aggregator

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestPathSegments(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"src/providers/vllm/init.py", []string{"src", "providers", "vllm", "init.py", "init", "py"}},
		{"evalhub/adapter.py", []string{"evalhub", "adapter.py", "adapter", "py"}},
		{"", nil},
		{"single", []string{"single"}},
	}
	for _, tt := range tests {
		got := pathSegments(tt.input)
		if len(got) != len(tt.want) {
			t.Errorf("pathSegments(%q) = %v, want %v", tt.input, got, tt.want)
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("pathSegments(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
			}
		}
	}
}

func TestExtractCodeGraphFilePaths(t *testing.T) {
	// Create a minimal code-graph.json
	tmpDir := t.TempDir()
	cgPath := filepath.Join(tmpDir, "code-graph.json")

	cg := map[string]interface{}{
		"nodes": []map[string]interface{}{
			{"id": "a", "kind": "FUNCTION", "name": "foo", "file": "src/providers/vllm/adapter.py", "line": 1},
			{"id": "b", "kind": "FUNCTION", "name": "bar", "file": "src/providers/vllm/adapter.py", "line": 10},
			{"id": "c", "kind": "CLASS", "name": "Baz", "file": "src/core/main.py", "line": 1},
		},
		"edges": []map[string]interface{}{
			{"from": "a", "to": "c", "kind": "CALLS", "label": "foo"},
		},
		"schema_version": "1",
	}
	raw, _ := json.Marshal(cg)
	if err := os.WriteFile(cgPath, raw, 0o644); err != nil {
		t.Fatal(err)
	}

	files, err := extractCodeGraphFilePaths(cgPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 {
		t.Errorf("expected 2 unique files, got %d: %v", len(files), files)
	}
}

func TestDetectCodeGraphRefs(t *testing.T) {
	tmpDir := t.TempDir()

	// Create two component directories with architecture and code-graph JSONs
	llamaDir := filepath.Join(tmpDir, "llama-stack")
	vllmDir := filepath.Join(tmpDir, "vllm")
	os.MkdirAll(llamaDir, 0o755)
	os.MkdirAll(vllmDir, 0o755)

	// llama-stack component-architecture.json
	llamaArch := map[string]interface{}{"component": "llama-stack"}
	raw, _ := json.Marshal(llamaArch)
	llamaArchPath := filepath.Join(llamaDir, "component-architecture.json")
	os.WriteFile(llamaArchPath, raw, 0o644)

	// llama-stack code-graph.json with vllm reference in file paths
	llamaCG := map[string]interface{}{
		"nodes": []map[string]interface{}{
			{"id": "a", "file": "src/ogx/providers/remote/inference/vllm/vllm.py"},
			{"id": "b", "file": "src/ogx/core/main.py"},
			{"id": "c", "file": "tests/unit/providers/inference/test_remote_vllm.py"},
		},
		"edges":         []interface{}{},
		"schema_version": "1",
	}
	raw, _ = json.Marshal(llamaCG)
	os.WriteFile(filepath.Join(llamaDir, "code-graph.json"), raw, 0o644)

	// vllm component-architecture.json
	vllmArch := map[string]interface{}{"component": "vllm"}
	raw, _ = json.Marshal(vllmArch)
	vllmArchPath := filepath.Join(vllmDir, "component-architecture.json")
	os.WriteFile(vllmArchPath, raw, 0o644)

	// vllm code-graph.json (no llama-stack references)
	vllmCG := map[string]interface{}{
		"nodes": []map[string]interface{}{
			{"id": "x", "file": "src/vllm/engine.py"},
		},
		"edges":         []interface{}{},
		"schema_version": "1",
	}
	raw, _ = json.Marshal(vllmCG)
	os.WriteFile(filepath.Join(vllmDir, "code-graph.json"), raw, 0o644)

	jsonPaths := []string{llamaArchPath, vllmArchPath}
	componentNames := []string{"llama-stack", "vllm"}

	refs := detectCodeGraphRefs(jsonPaths, componentNames)

	// Should find llama-stack -> vllm reference
	found := false
	for _, ref := range refs {
		if ref.From == "llama-stack" && ref.To == "vllm" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected llama-stack -> vllm code-ref, got refs: %+v", refs)
	}

	// Should NOT find vllm -> llama-stack (no such file paths)
	for _, ref := range refs {
		if ref.From == "vllm" && ref.To == "llama-stack" {
			t.Errorf("unexpected vllm -> llama-stack code-ref")
		}
	}
}
