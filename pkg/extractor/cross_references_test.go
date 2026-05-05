package extractor

import (
	"testing"
)

func TestInferNetworkingCapabilities_NoRBAC(t *testing.T) {
	arch := &ComponentArchitecture{}
	inferNetworkingCapabilities(arch)
	if len(arch.IngressRouting) != 0 {
		t.Errorf("expected 0 ingress entries, got %d", len(arch.IngressRouting))
	}
}

func TestInferNetworkingCapabilities_SingleAPIGroup(t *testing.T) {
	arch := &ComponentArchitecture{
		RBAC: &RBACData{
			ClusterRoles: []RBACRole{{
				Name: "test-role",
				Rules: []RBACRule{{
					APIGroups: []string{"gateway.networking.k8s.io"},
					Resources: []string{"httproutes"},
					Verbs:     []string{"get", "list", "watch"},
				}},
			}},
		},
	}
	inferNetworkingCapabilities(arch)
	if len(arch.IngressRouting) != 1 {
		t.Fatalf("expected 1 ingress entry, got %d", len(arch.IngressRouting))
	}
	entry := arch.IngressRouting[0]
	if entry.Kind != "HTTPRoute" {
		t.Errorf("expected Kind=HTTPRoute, got %q", entry.Kind)
	}
	if entry.Name != "rbac-inferred" {
		t.Errorf("expected Name=rbac-inferred, got %q", entry.Name)
	}
	if entry.Source != "rbac/test-role" {
		t.Errorf("expected Source=rbac/test-role, got %q", entry.Source)
	}
	if len(entry.RBACVerbs) != 3 {
		t.Errorf("expected 3 verbs, got %d", len(entry.RBACVerbs))
	}
	if entry.Note == "" {
		t.Error("expected Note to be set")
	}
}

func TestInferNetworkingCapabilities_MultipleResources(t *testing.T) {
	arch := &ComponentArchitecture{
		RBAC: &RBACData{
			ClusterRoles: []RBACRole{{
				Name: "manager-role",
				Rules: []RBACRule{
					{
						APIGroups: []string{"gateway.networking.k8s.io"},
						Resources: []string{"httproutes"},
						Verbs:     []string{"get", "create"},
					},
					{
						APIGroups: []string{"networking.istio.io"},
						Resources: []string{"virtualservices", "virtualservices/finalizers"},
						Verbs:     []string{"get", "list", "create"},
					},
					{
						APIGroups: []string{"networking.k8s.io"},
						Resources: []string{"ingresses"},
						Verbs:     []string{"get"},
					},
				},
			}},
		},
	}
	inferNetworkingCapabilities(arch)
	if len(arch.IngressRouting) != 3 {
		t.Fatalf("expected 3 ingress entries, got %d", len(arch.IngressRouting))
	}

	kinds := make(map[string]bool)
	for _, entry := range arch.IngressRouting {
		kinds[entry.Kind] = true
	}
	for _, expected := range []string{"HTTPRoute", "VirtualService", "Ingress"} {
		if !kinds[expected] {
			t.Errorf("missing expected Kind %q", expected)
		}
	}
}

func TestInferNetworkingCapabilities_Wildcard(t *testing.T) {
	arch := &ComponentArchitecture{
		RBAC: &RBACData{
			ClusterRoles: []RBACRole{{
				Name: "admin-role",
				Rules: []RBACRule{{
					APIGroups: []string{"gateway.networking.k8s.io"},
					Resources: []string{"*"},
					Verbs:     []string{"*"},
				}},
			}},
		},
	}
	inferNetworkingCapabilities(arch)
	// Should emit both HTTPRoute and Gateway
	if len(arch.IngressRouting) != 2 {
		t.Fatalf("expected 2 ingress entries for wildcard, got %d", len(arch.IngressRouting))
	}
	kinds := make(map[string]bool)
	for _, e := range arch.IngressRouting {
		kinds[e.Kind] = true
	}
	if !kinds["HTTPRoute"] || !kinds["Gateway"] {
		t.Errorf("expected HTTPRoute and Gateway, got %v", kinds)
	}
}

func TestInferNetworkingCapabilities_EnrichExisting(t *testing.T) {
	arch := &ComponentArchitecture{
		IngressRouting: []IngressResource{{
			Kind:   "HTTPRoute",
			Name:   "my-route",
			Source: "config/networking/route.yaml",
			Hosts:  []string{"example.com"},
		}},
		RBAC: &RBACData{
			ClusterRoles: []RBACRole{{
				Name: "test-role",
				Rules: []RBACRule{{
					APIGroups: []string{"gateway.networking.k8s.io"},
					Resources: []string{"httproutes"},
					Verbs:     []string{"get", "create"},
				}},
			}},
		},
	}
	inferNetworkingCapabilities(arch)
	// Should enrich existing, not create new
	if len(arch.IngressRouting) != 1 {
		t.Fatalf("expected 1 entry (enriched), got %d", len(arch.IngressRouting))
	}
	entry := arch.IngressRouting[0]
	if entry.Name != "my-route" {
		t.Errorf("expected Name=my-route (preserved), got %q", entry.Name)
	}
	if entry.Source != "config/networking/route.yaml" {
		t.Errorf("expected Source preserved, got %q", entry.Source)
	}
	if len(entry.Hosts) != 1 || entry.Hosts[0] != "example.com" {
		t.Errorf("expected Hosts preserved, got %v", entry.Hosts)
	}
	if len(entry.RBACVerbs) != 2 {
		t.Errorf("expected 2 RBACVerbs, got %d", len(entry.RBACVerbs))
	}
	if entry.Note == "" {
		t.Error("expected Note to be set on enriched entry")
	}
}

func TestInferNetworkingCapabilities_VerbMerging(t *testing.T) {
	// Two rules for the same resource should merge verbs
	arch := &ComponentArchitecture{
		RBAC: &RBACData{
			ClusterRoles: []RBACRole{{
				Name: "test-role",
				Rules: []RBACRule{
					{
						APIGroups: []string{"networking.istio.io"},
						Resources: []string{"virtualservices"},
						Verbs:     []string{"get", "list"},
					},
					{
						APIGroups: []string{"networking.istio.io"},
						Resources: []string{"virtualservices"},
						Verbs:     []string{"get", "patch", "update"},
					},
				},
			}},
		},
	}
	inferNetworkingCapabilities(arch)
	if len(arch.IngressRouting) != 1 {
		t.Fatalf("expected 1 entry (merged), got %d", len(arch.IngressRouting))
	}
	entry := arch.IngressRouting[0]
	if entry.Kind != "VirtualService" {
		t.Errorf("expected Kind=VirtualService, got %q", entry.Kind)
	}
	// Should have merged verbs: get, list, patch, update (deduplicated)
	if len(entry.RBACVerbs) != 4 {
		t.Errorf("expected 4 deduplicated verbs, got %d: %v", len(entry.RBACVerbs), entry.RBACVerbs)
	}
}
