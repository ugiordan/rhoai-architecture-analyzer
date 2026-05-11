package extractor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractKustomizeOverlayRefs(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		wantLen  int
		validate func(t *testing.T, refs []KustomizeOverlayRef)
	}{
		{
			name: "basic overlay with resources and patches",
			files: map[string]string{
				"config/overlays/default/kustomization.yaml": `
resources:
  - ../../base
  - ../../crd
patches:
  - path: manager_patch.yaml
namespace: system
namePrefix: my-operator-
`,
			},
			wantLen: 1,
			validate: func(t *testing.T, refs []KustomizeOverlayRef) {
				ref := refs[0]
				if ref.Namespace != "system" {
					t.Errorf("expected namespace 'system', got %q", ref.Namespace)
				}
				if ref.NamePrefix != "my-operator-" {
					t.Errorf("expected namePrefix 'my-operator-', got %q", ref.NamePrefix)
				}
				if len(ref.Resources) != 2 {
					t.Errorf("expected 2 resources, got %d", len(ref.Resources))
				}
				if len(ref.Patches) != 1 {
					t.Errorf("expected 1 patch, got %d", len(ref.Patches))
				}
			},
		},
		{
			name: "overlay with generators and images",
			files: map[string]string{
				"config/overlays/production/kustomization.yaml": `
resources:
  - ../default
configMapGenerator:
  - name: manager-config
secretGenerator:
  - name: webhook-cert
images:
  - name: controller
    newName: registry.example.com/controller
    newTag: v1.2.3
commonLabels:
  app.kubernetes.io/part-of: my-platform
`,
			},
			wantLen: 1,
			validate: func(t *testing.T, refs []KustomizeOverlayRef) {
				ref := refs[0]
				if len(ref.ConfigMapGens) != 1 || ref.ConfigMapGens[0] != "manager-config" {
					t.Errorf("expected configMapGenerator 'manager-config', got %v", ref.ConfigMapGens)
				}
				if len(ref.SecretGens) != 1 || ref.SecretGens[0] != "webhook-cert" {
					t.Errorf("expected secretGenerator 'webhook-cert', got %v", ref.SecretGens)
				}
				if len(ref.Images) != 1 {
					t.Errorf("expected 1 image transform, got %d", len(ref.Images))
				}
				if len(ref.CommonLabels) != 1 {
					t.Errorf("expected 1 common label, got %d", len(ref.CommonLabels))
				}
			},
		},
		{
			name: "multiple overlays",
			files: map[string]string{
				"config/overlays/odh/kustomization.yaml": `
resources:
  - ../default
`,
				"config/overlays/rhoai/kustomization.yaml": `
resources:
  - ../default
patches:
  - path: rhoai_patch.yaml
`,
			},
			wantLen: 2,
		},
		{
			name:    "no config directory",
			files:   map[string]string{},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			for path, content := range tt.files {
				fullPath := filepath.Join(dir, path)
				if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
					t.Fatal(err)
				}
			}

			refs := extractKustomizeOverlayRefs(dir)
			if len(refs) != tt.wantLen {
				t.Errorf("expected %d refs, got %d", tt.wantLen, len(refs))
			}
			if tt.validate != nil && len(refs) == tt.wantLen {
				tt.validate(t, refs)
			}
		})
	}
}
