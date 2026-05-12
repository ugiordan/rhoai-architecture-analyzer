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
		{
			name:       "Go import scanning",
			selfName:   "my-operator",
			components: []string{"kserve", "model-registry", "my-operator"},
			files: map[string]string{
				"pkg/controller/reconciler.go": `package controller

import (
	"context"

	kservev1beta1 "github.com/kserve/kserve/pkg/apis/serving/v1beta1"
	"github.com/opendatahub-io/model-registry/pkg/api"
)

func Reconcile(ctx context.Context) error {
	_ = kservev1beta1.InferenceService{}
	return nil
}
`,
			},
			wantLen: 2,
			validate: func(t *testing.T, refs []ComponentRef) {
				found := map[string]bool{}
				for _, r := range refs {
					found[r.Target] = true
					if r.Type != "import" {
						t.Errorf("expected type 'import', got %q for %s", r.Type, r.Target)
					}
				}
				if !found["kserve"] {
					t.Error("expected kserve import ref")
				}
				if !found["model-registry"] {
					t.Error("expected model-registry import ref")
				}
			},
		},
		{
			name:       "Python import scanning",
			selfName:   "llama-stack",
			components: []string{"vllm", "kserve", "llama-stack"},
			files: map[string]string{
				"src/inference/adapter.py": `import os
from vllm import LLM, SamplingParams
from kserve.protocol import grpc as kserve_grpc

class InferenceAdapter:
    pass
`,
			},
			wantLen: 2,
			validate: func(t *testing.T, refs []ComponentRef) {
				found := map[string]bool{}
				for _, r := range refs {
					found[r.Target] = true
					if r.Type != "import" {
						t.Errorf("expected type 'import', got %q for %s", r.Type, r.Target)
					}
				}
				if !found["vllm"] {
					t.Error("expected vllm import ref")
				}
				if !found["kserve"] {
					t.Error("expected kserve import ref")
				}
			},
		},
		{
			name:       "Go import with underscore component name",
			selfName:   "dashboard",
			components: []string{"data-science-pipelines", "dashboard"},
			files: map[string]string{
				"pkg/client/dsp.go": `package client

import (
	dspv1 "github.com/opendatahub-io/data_science_pipelines/api/v1alpha1"
)

func GetPipeline() {}
`,
			},
			wantLen: 1,
			validate: func(t *testing.T, refs []ComponentRef) {
				if refs[0].Target != "data-science-pipelines" {
					t.Errorf("expected target 'data-science-pipelines', got %q", refs[0].Target)
				}
			},
		},
		{
			name:       "test files excluded from import scanning",
			selfName:   "my-service",
			components: []string{"vllm"},
			files: map[string]string{
				"pkg/handler_test.go": `package pkg

import "github.com/vllm-project/vllm/client"

func TestHandler() {}
`,
				"tests/test_integration.py": `import vllm

def test_integration():
    pass
`,
			},
			wantLen: 0,
		},
		{
			name:       "deduplication across directory and import refs",
			selfName:   "orchestrator",
			components: []string{"vllm", "orchestrator"},
			dirs: []string{
				"providers/vllm",
			},
			files: map[string]string{
				"providers/vllm/config.py": "",
				"pkg/inference.go": `package pkg

import vllmclient "github.com/vllm-project/vllm/client"

func Infer() {}
`,
			},
			wantLen: 2, // provider + import (different types)
			validate: func(t *testing.T, refs []ComponentRef) {
				types := map[string]bool{}
				for _, r := range refs {
					if r.Target != "vllm" {
						t.Errorf("expected target 'vllm', got %q", r.Target)
					}
					types[r.Type] = true
				}
				if !types["provider"] {
					t.Error("expected provider type")
				}
				if !types["import"] {
					t.Error("expected import type")
				}
			},
		},
		{
			name:       "hidden and vendor dirs skipped",
			selfName:   "my-service",
			components: []string{"vllm"},
			dirs: []string{
				".hidden/providers/vllm",
				"vendor/providers/vllm",
				"node_modules/providers/vllm",
			},
			files: map[string]string{
				".hidden/providers/vllm/main.go":       "package vllm\n",
				"vendor/providers/vllm/main.go":        "package vllm\n",
				"node_modules/providers/vllm/index.js": "",
			},
			wantLen: 0,
		},
		{
			name:       "directory matches component but no provider ancestor",
			selfName:   "my-service",
			components: []string{"vllm"},
			dirs:       []string{"src/vllm"},
			files: map[string]string{
				"src/vllm/main.py": "# not under a provider dir\n",
			},
			wantLen: 0,
		},
		{
			name:       "empty component list uses defaults",
			selfName:   "my-service",
			components: []string{}, // triggers defaultKnownComponents
			dirs:       []string{"adapters/vllm"},
			files: map[string]string{
				"adapters/vllm/__init__.py": "",
			},
			wantLen: 1,
			validate: func(t *testing.T, refs []ComponentRef) {
				if refs[0].Target != "vllm" {
					t.Errorf("expected target 'vllm', got %q", refs[0].Target)
				}
			},
		},
		{
			name:       "deeply nested provider directory",
			selfName:   "orchestrator",
			components: []string{"kserve"},
			dirs:       []string{"src/inference/providers/remote/kserve"},
			files: map[string]string{
				"src/inference/providers/remote/kserve/config.py": "",
			},
			wantLen: 1,
			validate: func(t *testing.T, refs []ComponentRef) {
				if refs[0].Type != "provider" {
					t.Errorf("expected type 'provider', got %q", refs[0].Type)
				}
			},
		},
		{
			name:       "file suffix detection multiple suffixes",
			selfName:   "controller",
			components: []string{"kserve", "modelmesh", "vllm"},
			files: map[string]string{
				"pkg/kserve_client.go":    "package pkg\n",
				"pkg/modelmesh_plugin.go": "package pkg\n",
				"pkg/vllm_provider.py":    "# provider\n",
			},
			wantLen: 3,
			validate: func(t *testing.T, refs []ComponentRef) {
				typeMap := map[string]string{}
				for _, r := range refs {
					typeMap[r.Target] = r.Type
				}
				if typeMap["kserve"] != "client" {
					t.Errorf("kserve type: got %q, want 'client'", typeMap["kserve"])
				}
				if typeMap["modelmesh"] != "plugin" {
					t.Errorf("modelmesh type: got %q, want 'plugin'", typeMap["modelmesh"])
				}
				if typeMap["vllm"] != "provider" {
					t.Errorf("vllm type: got %q, want 'provider'", typeMap["vllm"])
				}
			},
		},
		{
			name:       "Python conftest and test prefix excluded",
			selfName:   "my-service",
			components: []string{"vllm"},
			files: map[string]string{
				"tests/test_vllm_integration.py": "import vllm\ndef test_it(): pass\n",
			},
			wantLen: 0,
		},
		{
			name:       "same type deduplication",
			selfName:   "orchestrator",
			components: []string{"vllm"},
			dirs: []string{
				"providers/inference/vllm",
				"providers/embeddings/vllm",
			},
			files: map[string]string{
				"providers/inference/vllm/main.py":  "",
				"providers/embeddings/vllm/main.py": "",
			},
			wantLen: 1, // same target+type deduped
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
