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
			{Group: "a.io", Version: "v2", Kind: "Alpha"},
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
			{Method: "DELETE", Path: "/api/v1/models"},
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
			ClusterRoleBindings: []RBACBinding{
				{Name: "binding-z"},
				{Name: "binding-a"},
			},
			Roles: []RBACRole{
				{Name: "ns-role-z"},
				{Name: "ns-role-a"},
			},
			RoleBindings: []RBACBinding{
				{Name: "ns-binding-z"},
				{Name: "ns-binding-a"},
			},
			KubebuilderMarkers: []RBACMarker{
				{File: "z.go", Line: 10},
				{File: "a.go", Line: 5},
				{File: "a.go", Line: 3},
			},
		},
		NetworkPolicies: []NetworkPolicy{
			{Name: "netpol-z"},
			{Name: "netpol-a"},
		},
		ControllerWatch: []ControllerWatch{
			{Type: "Watches", GVK: "v1/Pod"},
			{Type: "For", GVK: "v1/Service"},
			{Type: "For", GVK: "v1/Deployment"},
		},
		Dockerfiles: []DockerfileInfo{
			{Path: "z/Dockerfile"},
			{Path: "a/Dockerfile"},
		},
		Webhooks: []WebhookConfig{
			{Name: "webhook-z"},
			{Name: "webhook-a"},
		},
		ConfigMaps: []ConfigMapRef{
			{Name: "cm-z"},
			{Name: "cm-a"},
		},
		IngressRouting: []IngressResource{
			{Kind: "Route", Name: "route-z"},
			{Kind: "Gateway", Name: "gw-a"},
			{Kind: "Route", Name: "route-a"},
		},
		Dependencies: &DependencyData{
			GoModules: []GoModule{
				{Module: "z-lib", Version: "v1.0"},
				{Module: "a-lib", Version: "v2.0"},
			},
			InternalODH: []InternalODH{
				{Component: "z-comp"},
				{Component: "a-comp"},
			},
			ReplaceDirectives: []ReplaceDirective{
				{Original: "z-orig"},
				{Original: "a-orig"},
			},
		},
		CacheConfig: &CacheConfig{
			FilteredTypes: []CacheFilteredType{
				{Type: "Secret"},
				{Type: "ConfigMap"},
			},
			TransformTypes:  []CacheTransform{{Type: "z-type"}, {Type: "a-type"}},
			DisabledTypes:   []string{"z-disabled", "a-disabled"},
			ImplicitInformers: []ImplicitInformer{
				{Type: "z-informer"},
				{Type: "a-informer"},
			},
			Issues: []string{"z-issue", "a-issue"},
		},
	}

	SortOutput(arch)

	// Verify ordering: top-level slices
	if arch.CRDs[0].Group != "a.io" || arch.CRDs[0].Version != "v1" || arch.CRDs[0].Kind != "Alpha" {
		t.Errorf("CRDs not sorted: got %s/%s/%s first", arch.CRDs[0].Group, arch.CRDs[0].Version, arch.CRDs[0].Kind)
	}
	if arch.CRDs[1].Version != "v1" || arch.CRDs[1].Kind != "Beta" {
		t.Errorf("CRDs not sorted within group: got %s/%s second", arch.CRDs[1].Version, arch.CRDs[1].Kind)
	}
	if arch.CRDs[2].Version != "v2" || arch.CRDs[2].Kind != "Alpha" {
		t.Errorf("CRDs version tiebreaker not sorted: got %s/%s third", arch.CRDs[2].Version, arch.CRDs[2].Kind)
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
	// Verify method tiebreaker within same path (/api/v1/models: DELETE < GET < POST)
	if arch.HTTPEndpoints[1].Method != "DELETE" {
		t.Errorf("HTTPEndpoints method tiebreaker failed: got %s second", arch.HTTPEndpoints[1].Method)
	}
	if arch.HTTPEndpoints[2].Method != "GET" {
		t.Errorf("HTTPEndpoints method tiebreaker failed: got %s third", arch.HTTPEndpoints[2].Method)
	}
	if arch.FeatureGates[0].Name != "AlphaFeature" {
		t.Errorf("FeatureGates not sorted: got %s first", arch.FeatureGates[0].Name)
	}
	if arch.ExternalConnections[0].Type != "api" {
		t.Errorf("ExternalConnections not sorted: got %s first", arch.ExternalConnections[0].Type)
	}
	// Verify service tiebreaker within same type (database: mysql < postgres)
	if arch.ExternalConnections[1].Service != "mysql" {
		t.Errorf("ExternalConnections service tiebreaker failed: got %s second", arch.ExternalConnections[1].Service)
	}
	if arch.ExternalConnections[2].Service != "postgres" {
		t.Errorf("ExternalConnections service tiebreaker failed: got %s third", arch.ExternalConnections[2].Service)
	}

	// RBAC sub-types
	if arch.RBAC.ClusterRoles[0].Name != "role-a" {
		t.Errorf("RBAC ClusterRoles not sorted: got %s first", arch.RBAC.ClusterRoles[0].Name)
	}
	if arch.RBAC.ClusterRoleBindings[0].Name != "binding-a" {
		t.Errorf("RBAC ClusterRoleBindings not sorted: got %s first", arch.RBAC.ClusterRoleBindings[0].Name)
	}
	if arch.RBAC.KubebuilderMarkers[0].File != "a.go" || arch.RBAC.KubebuilderMarkers[0].Line != 3 {
		t.Errorf("RBAC KubebuilderMarkers not sorted: got %s:%d first", arch.RBAC.KubebuilderMarkers[0].File, arch.RBAC.KubebuilderMarkers[0].Line)
	}
	if arch.RBAC.Roles[0].Name != "ns-role-a" {
		t.Errorf("RBAC Roles not sorted: got %s first", arch.RBAC.Roles[0].Name)
	}
	if arch.RBAC.RoleBindings[0].Name != "ns-binding-a" {
		t.Errorf("RBAC RoleBindings not sorted: got %s first", arch.RBAC.RoleBindings[0].Name)
	}

	// Multi-field composite sorts
	if arch.ControllerWatch[0].Type != "For" || arch.ControllerWatch[0].GVK != "v1/Deployment" {
		t.Errorf("ControllerWatch not sorted: got %s/%s first", arch.ControllerWatch[0].Type, arch.ControllerWatch[0].GVK)
	}
	if arch.IngressRouting[0].Kind != "Gateway" {
		t.Errorf("IngressRouting not sorted: got %s first", arch.IngressRouting[0].Kind)
	}
	if arch.IngressRouting[1].Name != "route-a" {
		t.Errorf("IngressRouting not sorted within kind: got %s second", arch.IngressRouting[1].Name)
	}

	// Single-field sorts
	if arch.NetworkPolicies[0].Name != "netpol-a" {
		t.Errorf("NetworkPolicies not sorted: got %s first", arch.NetworkPolicies[0].Name)
	}
	if arch.Dockerfiles[0].Path != "a/Dockerfile" {
		t.Errorf("Dockerfiles not sorted: got %s first", arch.Dockerfiles[0].Path)
	}
	if arch.Webhooks[0].Name != "webhook-a" {
		t.Errorf("Webhooks not sorted: got %s first", arch.Webhooks[0].Name)
	}
	if arch.ConfigMaps[0].Name != "cm-a" {
		t.Errorf("ConfigMaps not sorted: got %s first", arch.ConfigMaps[0].Name)
	}

	// Dependencies sub-slices
	if arch.Dependencies.GoModules[0].Module != "a-lib" {
		t.Errorf("GoModules not sorted: got %s first", arch.Dependencies.GoModules[0].Module)
	}
	if arch.Dependencies.InternalODH[0].Component != "a-comp" {
		t.Errorf("InternalODH not sorted: got %s first", arch.Dependencies.InternalODH[0].Component)
	}
	if arch.Dependencies.ReplaceDirectives[0].Original != "a-orig" {
		t.Errorf("ReplaceDirectives not sorted: got %s first", arch.Dependencies.ReplaceDirectives[0].Original)
	}

	// CacheConfig sub-slices
	if arch.CacheConfig.FilteredTypes[0].Type != "ConfigMap" {
		t.Errorf("CacheConfig FilteredTypes not sorted: got %s first", arch.CacheConfig.FilteredTypes[0].Type)
	}
	if arch.CacheConfig.TransformTypes[0].Type != "a-type" {
		t.Errorf("CacheConfig TransformTypes not sorted: got %s first", arch.CacheConfig.TransformTypes[0].Type)
	}
	if arch.CacheConfig.DisabledTypes[0] != "a-disabled" {
		t.Errorf("CacheConfig DisabledTypes not sorted: got %s first", arch.CacheConfig.DisabledTypes[0])
	}
	if arch.CacheConfig.ImplicitInformers[0].Type != "a-informer" {
		t.Errorf("CacheConfig ImplicitInformers not sorted: got %s first", arch.CacheConfig.ImplicitInformers[0].Type)
	}
	if arch.CacheConfig.Issues[0] != "a-issue" {
		t.Errorf("CacheConfig Issues not sorted: got %s first", arch.CacheConfig.Issues[0])
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
