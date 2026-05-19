package sbom

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestGenerate_EmptyData(t *testing.T) {
	bom := Generate(map[string]interface{}{})
	if bom.BOMFormat != "CycloneDX" {
		t.Error("BOMFormat should be CycloneDX")
	}
	if bom.SpecVersion != "1.5" {
		t.Error("SpecVersion should be 1.5")
	}
	if bom.Metadata.Component.Name != "unknown" {
		t.Error("empty data should produce 'unknown' component name")
	}
}

func TestGenerate_GoModules(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"dependencies": map[string]interface{}{
			"go_modules": []interface{}{
				map[string]interface{}{"module": "github.com/example/lib", "version": "v1.2.3"},
			},
		},
	}
	bom := Generate(data)
	var found bool
	for _, c := range bom.Components {
		if c.Name == "github.com/example/lib" && c.Version == "v1.2.3" && c.Type == "library" {
			if !strings.Contains(c.PURL, "pkg:golang/") {
				t.Error("Go module should have golang PURL")
			}
			found = true
		}
	}
	if !found {
		t.Error("should include Go module as library component")
	}
}

func TestGenerate_Dockerfiles(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"dockerfiles": []interface{}{
			map[string]interface{}{
				"path":       "Dockerfile",
				"base_image": "registry.access.redhat.com/ubi9/ubi-minimal:latest",
				"stages":     2,
				"user":       "1001",
			},
		},
	}
	bom := Generate(data)
	var found bool
	for _, c := range bom.Components {
		if c.Type == "container" && strings.Contains(c.Name, "ubi-minimal") {
			found = true
			if c.Version != "latest" {
				t.Errorf("tag should be 'latest', got %q", c.Version)
			}
		}
	}
	if !found {
		t.Error("should include Dockerfile base image as container component")
	}
}

func TestGenerate_DeploymentContainers(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"deployments": []interface{}{
			map[string]interface{}{
				"name": "my-app",
				"containers": []interface{}{
					map[string]interface{}{
						"name":  "main",
						"image": "quay.io/org/app:v1.0",
						"security_context": map[string]interface{}{
							"runAsNonRoot": true,
						},
						"resources": map[string]interface{}{
							"requests": map[string]interface{}{"cpu": "100m", "memory": "256Mi"},
							"limits":   map[string]interface{}{"cpu": "500m", "memory": "512Mi"},
						},
					},
				},
			},
		},
	}
	bom := Generate(data)
	var found bool
	for _, c := range bom.Components {
		if c.Type == "container" && c.Name == "quay.io/org/app" {
			found = true
			hasDeployment := false
			hasCPU := false
			hasRunAs := false
			for _, p := range c.Properties {
				if p.Name == "arch-analyzer:deployment" && p.Value == "my-app" {
					hasDeployment = true
				}
				if p.Name == "arch-analyzer:cpu-request" && p.Value == "100m" {
					hasCPU = true
				}
				if p.Name == "arch-analyzer:runAsNonRoot" && p.Value == "true" {
					hasRunAs = true
				}
			}
			if !hasDeployment {
				t.Error("should include deployment name in properties")
			}
			if !hasCPU {
				t.Error("should include CPU request in properties")
			}
			if !hasRunAs {
				t.Error("should include runAsNonRoot in properties")
			}
		}
	}
	if !found {
		t.Error("should include deployment container image as container component")
	}
}

func TestGenerate_OperatorImageConstants(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"operator_config": []interface{}{
			map[string]interface{}{
				"name":     "DefaultImage",
				"value":    "quay.io/org/operator:latest",
				"category": "image",
				"source":   "pkg/config.go",
			},
		},
	}
	bom := Generate(data)
	var found bool
	for _, c := range bom.Components {
		if c.Type == "container" && c.Name == "quay.io/org/operator" {
			found = true
		}
	}
	if !found {
		t.Error("should include operator image constant as container component")
	}
}

func TestGenerate_SkipsVariableImages(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"dockerfiles": []interface{}{
			map[string]interface{}{
				"path":       "Dockerfile",
				"base_image": "${BASE_IMAGE}",
			},
		},
	}
	bom := Generate(data)
	for _, c := range bom.Components {
		if c.Type == "container" && strings.HasPrefix(c.Name, "$") {
			t.Error("should skip variable-based image references")
		}
	}
}

func TestGenerate_ImageDigest(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"dockerfiles": []interface{}{
			map[string]interface{}{
				"path":       "Dockerfile",
				"base_image": "registry.redhat.io/ubi9/ubi-minimal@sha256:abc123def456",
			},
		},
	}
	bom := Generate(data)
	for _, c := range bom.Components {
		if c.Type == "container" && strings.Contains(c.Name, "ubi-minimal") {
			if len(c.Hashes) == 0 {
				t.Error("should include hash for digest-pinned images")
			}
		}
	}
}

func TestGenerateJSON_ValidJSON(t *testing.T) {
	data := map[string]interface{}{"component": "test"}
	out, err := GenerateJSON(data)
	if err != nil {
		t.Fatalf("GenerateJSON failed: %v", err)
	}
	var bom BOM
	if err := json.Unmarshal(out, &bom); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if bom.BOMFormat != "CycloneDX" {
		t.Error("should produce CycloneDX format")
	}
}

func TestParseImageRef(t *testing.T) {
	tests := []struct {
		input  string
		name   string
		tag    string
		digest string
	}{
		{"nginx:latest", "nginx", "latest", ""},
		{"registry.io/org/app:v1.0", "registry.io/org/app", "v1.0", ""},
		{"img@sha256:abc123", "img", "", "sha256:abc123"},
		{"registry.io/app:v1@sha256:abc", "registry.io/app", "v1", "sha256:abc"},
		{"registry.io:5000/app:v1", "registry.io:5000/app", "v1", ""},
	}
	for _, tc := range tests {
		name, tag, digest := parseImageRef(tc.input)
		if name != tc.name {
			t.Errorf("parseImageRef(%q): name=%q, want %q", tc.input, name, tc.name)
		}
		if tag != tc.tag {
			t.Errorf("parseImageRef(%q): tag=%q, want %q", tc.input, tag, tc.tag)
		}
		if digest != tc.digest {
			t.Errorf("parseImageRef(%q): digest=%q, want %q", tc.input, digest, tc.digest)
		}
	}
}
