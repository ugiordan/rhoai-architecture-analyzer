package extractor

import (
	"encoding/json"
	"testing"
)

func TestSortOutput_Deterministic(t *testing.T) {
	// Build an architecture with intentionally unordered slices
	arch := &ComponentArchitecture{
		Component: "test",
		CRDs: []CRD{
			{Group: "z.io", Version: "v1", Kind: "Zebra"},
			{Group: "a.io", Version: "v1", Kind: "Alpha"},
			{Group: "a.io", Version: "v1", Kind: "Beta"},
		},
		Services: []Service{
			{Name: "svc-z"},
			{Name: "svc-a"},
			{Name: "svc-m"},
		},
		Deployments: []Deployment{
			{Name: "deploy-b"},
			{Name: "deploy-a"},
		},
		Secrets: []SecretRef{
			{Name: "secret-z", ReferencedBy: []string{"deploy/b", "deploy/a"}},
			{Name: "secret-a", ReferencedBy: []string{"deploy/c"}},
		},
		HTTPEndpoints: []HTTPEndpoint{
			{Method: "POST", Path: "/api/v1/models"},
			{Method: "GET", Path: "/api/v1/models"},
			{Method: "GET", Path: "/api/v1/health"},
		},
		FeatureGates: []FeatureGate{
			{Name: "ZetaFeature"},
			{Name: "AlphaFeature"},
		},
		ExternalConnections: []ExternalConnection{
			{Type: "database", Service: "postgres", Source: "pkg/db.go:10"},
			{Type: "api", Service: "auth", Source: "pkg/auth.go:5"},
			{Type: "database", Service: "mysql", Source: "pkg/db.go:20"},
		},
		RBAC: &RBACData{
			ClusterRoles: []RBACRole{
				{Name: "role-z"},
				{Name: "role-a"},
			},
		},
	}

	SortOutput(arch)

	// Verify ordering
	if arch.CRDs[0].Group != "a.io" || arch.CRDs[0].Kind != "Alpha" {
		t.Errorf("CRDs not sorted: got %s/%s first", arch.CRDs[0].Group, arch.CRDs[0].Kind)
	}
	if arch.CRDs[1].Kind != "Beta" {
		t.Errorf("CRDs not sorted within group: got %s second", arch.CRDs[1].Kind)
	}
	if arch.Services[0].Name != "svc-a" {
		t.Errorf("Services not sorted: got %s first", arch.Services[0].Name)
	}
	if arch.Deployments[0].Name != "deploy-a" {
		t.Errorf("Deployments not sorted: got %s first", arch.Deployments[0].Name)
	}
	if arch.Secrets[0].Name != "secret-a" {
		t.Errorf("Secrets not sorted: got %s first", arch.Secrets[0].Name)
	}
	if arch.Secrets[1].ReferencedBy[0] != "deploy/a" {
		t.Errorf("Secret ReferencedBy not sorted: got %s first", arch.Secrets[1].ReferencedBy[0])
	}
	if arch.HTTPEndpoints[0].Path != "/api/v1/health" {
		t.Errorf("HTTPEndpoints not sorted: got %s first", arch.HTTPEndpoints[0].Path)
	}
	if arch.FeatureGates[0].Name != "AlphaFeature" {
		t.Errorf("FeatureGates not sorted: got %s first", arch.FeatureGates[0].Name)
	}
	if arch.ExternalConnections[0].Type != "api" {
		t.Errorf("ExternalConnections not sorted: got %s first", arch.ExternalConnections[0].Type)
	}
	if arch.RBAC.ClusterRoles[0].Name != "role-a" {
		t.Errorf("RBAC ClusterRoles not sorted: got %s first", arch.RBAC.ClusterRoles[0].Name)
	}
}

func TestSortOutput_IdempotentJSON(t *testing.T) {
	arch := &ComponentArchitecture{
		Component: "test",
		Secrets: []SecretRef{
			{Name: "z-secret", ReferencedBy: []string{"b", "a"}},
			{Name: "a-secret", ReferencedBy: []string{"c"}},
		},
		HTTPEndpoints: []HTTPEndpoint{
			{Method: "POST", Path: "/b"},
			{Method: "GET", Path: "/a"},
		},
	}

	SortOutput(arch)
	json1, _ := json.MarshalIndent(arch, "", "  ")

	// Sort again, should produce identical output
	SortOutput(arch)
	json2, _ := json.MarshalIndent(arch, "", "  ")

	if string(json1) != string(json2) {
		t.Error("SortOutput is not idempotent: second sort produced different JSON")
	}
}

func TestSortOutput_NilSafe(t *testing.T) {
	// Should not panic on nil
	SortOutput(nil)

	// Should not panic on empty
	SortOutput(&ComponentArchitecture{})

	// Should not panic on nil nested structs
	SortOutput(&ComponentArchitecture{
		RBAC:         nil,
		Dependencies: nil,
		CacheConfig:  nil,
	})
}
