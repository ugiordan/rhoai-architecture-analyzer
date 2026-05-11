package extractor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractComponentRefs(t *testing.T) {
	tests := []struct {
		name       string
		selfName   string
		components []string
		dirs       []string // directories to create
		files      map[string]string // file path -> content
		wantLen    int
		validate   func(t *testing.T, refs []ComponentRef)
	}{
		{
			name:       "provider directory pattern",
			selfName:   "llama-stack",
			components: []string{"vllm", "kserve", "llama-stack"},
			dirs: []string{
				"src/providers/remote/inference/vllm",
			},
			files: map[string]string{
				"src/providers/remote/inference/vllm/vllm.py": "class VLLMInferenceAdapter:\n    pass\n",
			},
			wantLen: 1,
			validate: func(t *testing.T, refs []ComponentRef) {
				if refs[0].Target != "vllm" {
					t.Errorf("expected target 'vllm', got %q", refs[0].Target)
				}
				if refs[0].Type != "provider" {
					t.Errorf("expected type 'provider', got %q", refs[0].Type)
				}
			},
		},
		{
			name:       "adapter file pattern",
			selfName:   "my-service",
			components: []string{"kserve", "my-service"},
			files: map[string]string{
				"pkg/kserve_adapter.go": "package pkg\n",
			},
			wantLen: 1,
			validate: func(t *testing.T, refs []ComponentRef) {
				if refs[0].Target != "kserve" {
					t.Errorf("expected target 'kserve', got %q", refs[0].Target)
				}
				if refs[0].Type != "adapter" {
					t.Errorf("expected type 'adapter', got %q", refs[0].Type)
				}
			},
		},
		{
			name:       "self-references excluded",
			selfName:   "vllm",
			components: []string{"vllm", "kserve"},
			dirs: []string{
				"src/providers/vllm",
			},
			wantLen: 0,
		},
		{
			name:       "multiple refs across patterns",
			selfName:   "orchestrator",
			components: []string{"vllm", "kserve", "modelmesh", "orchestrator"},
			dirs: []string{
				"backends/vllm",
				"plugins/kserve",
			},
			files: map[string]string{
				"backends/vllm/config.py": "",
				"plugins/kserve/handler.go": "package kserve\n",
				"pkg/modelmesh_client.go": "package pkg\n",
			},
			wantLen: 3,
		},
		{
			name:       "underscore variant matches",
			selfName:   "my-service",
			components: []string{"model-registry"},
			dirs: []string{
				"adapters/model_registry",
			},
			files: map[string]string{
				"adapters/model_registry/__init__.py": "",
			},
			wantLen: 1,
			validate: func(t *testing.T, refs []ComponentRef) {
				if refs[0].Target != "model-registry" {
					t.Errorf("expected target 'model-registry', got %q", refs[0].Target)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()

			for _, d := range tt.dirs {
				if err := os.MkdirAll(filepath.Join(dir, d), 0o755); err != nil {
					t.Fatal(err)
				}
			}
			for path, content := range tt.files {
				fullPath := filepath.Join(dir, path)
				if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
					t.Fatal(err)
				}
			}

			refs := extractComponentRefs(dir, tt.selfName, tt.components)
			if len(refs) != tt.wantLen {
				t.Errorf("expected %d refs, got %d: %+v", tt.wantLen, len(refs), refs)
			}
			if tt.validate != nil && len(refs) == tt.wantLen {
				tt.validate(t, refs)
			}
		})
	}
}
