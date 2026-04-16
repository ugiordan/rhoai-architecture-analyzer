package arch

import (
	"testing"
)

func TestParse(t *testing.T) {
	raw := map[string]interface{}{
		"component": "test-operator",
		"crds": []interface{}{
			map[string]interface{}{
				"group":   "apps.example.io",
				"version": "v1alpha1",
				"kind":    "Widget",
				"scope":   "Namespaced",
				"fields_count": float64(5),
				"validation_rules": []interface{}{"self.metadata.name == 'default'"},
				"source":  "config/crd/widgets.yaml",
			},
			map[string]interface{}{
				"group":   "apps.example.io",
				"version": "v1",
				"kind":    "Widget",
				"scope":   "Namespaced",
				"fields_count": float64(8),
				"source":  "config/crd/widgets.yaml",
			},
		},
		"rbac": map[string]interface{}{
			"cluster_roles": []interface{}{
				map[string]interface{}{
					"name":   "manager-role",
					"source": "config/rbac/role.yaml",
					"rules": []interface{}{
						map[string]interface{}{
							"apiGroups": []interface{}{"apps.example.io"},
							"resources": []interface{}{"widgets"},
							"verbs":     []interface{}{"get", "list", "create"},
						},
					},
				},
			},
			"kubebuilder_markers": []interface{}{
				map[string]interface{}{
					"file":   "controllers/widget_controller.go",
					"line":   float64(15),
					"marker": "//+kubebuilder:rbac:groups=apps.example.io,resources=widgets,verbs=get;list",
					"parsed": map[string]interface{}{"groups": "apps.example.io"},
				},
			},
		},
		"webhooks": []interface{}{
			map[string]interface{}{
				"name":           "mwidget.kb.io",
				"type":           "mutating",
				"path":           "/mutate-apps-v1-widget",
				"failure_policy": "Fail",
				"rules": []interface{}{
					map[string]interface{}{
						"operations":  []interface{}{"CREATE", "UPDATE"},
						"apiGroups":   []interface{}{"apps.example.io"},
						"apiVersions": []interface{}{"v1"},
						"resources":   []interface{}{"widgets"},
					},
				},
				"source": "config/webhook/manifests.yaml",
			},
		},
		"secrets_referenced": []interface{}{
			map[string]interface{}{
				"name":            "webhook-cert",
				"type":            "kubernetes.io/tls",
				"referenced_by":   []interface{}{"deployment/controller-manager"},
				"provisioned_by":  "cert-manager",
			},
		},
		"cache_config": map[string]interface{}{
			"filtered_types": []interface{}{
				map[string]interface{}{
					"type":        "Secret",
					"filter_kind": "label",
					"filter":      "app=myapp",
				},
			},
			"issues": []interface{}{
				"Type Widget is watched but has no cache filter",
			},
		},
	}

	data, err := Parse(raw)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	if data.Component != "test-operator" {
		t.Errorf("Component = %q, want %q", data.Component, "test-operator")
	}
	if len(data.CRDs) != 2 {
		t.Fatalf("CRDs count = %d, want 2", len(data.CRDs))
	}
	if data.CRDs[0].Kind != "Widget" || data.CRDs[0].Version != "v1alpha1" {
		t.Errorf("CRD[0] = %+v, want Widget/v1alpha1", data.CRDs[0])
	}
	if data.CRDs[1].Version != "v1" {
		t.Errorf("CRD[1].Version = %q, want %q", data.CRDs[1].Version, "v1")
	}
	if len(data.CRDs[0].ValidationRules) != 1 {
		t.Errorf("CRD[0].ValidationRules count = %d, want 1", len(data.CRDs[0].ValidationRules))
	}

	// RBAC
	if len(data.RBAC.ClusterRoles) != 1 {
		t.Fatalf("ClusterRoles count = %d, want 1", len(data.RBAC.ClusterRoles))
	}
	if data.RBAC.ClusterRoles[0].Name != "manager-role" {
		t.Errorf("ClusterRole name = %q, want %q", data.RBAC.ClusterRoles[0].Name, "manager-role")
	}
	if data.RBAC.ClusterRoles[0].Source != "config/rbac/role.yaml" {
		t.Errorf("ClusterRole source = %q, want %q", data.RBAC.ClusterRoles[0].Source, "config/rbac/role.yaml")
	}
	if len(data.RBAC.ClusterRoles[0].Rules) != 1 {
		t.Fatalf("RBAC rules count = %d, want 1", len(data.RBAC.ClusterRoles[0].Rules))
	}
	if len(data.RBAC.ClusterRoles[0].Rules[0].APIGroups) != 1 || data.RBAC.ClusterRoles[0].Rules[0].APIGroups[0] != "apps.example.io" {
		t.Errorf("RBAC rule APIGroups = %v, want [apps.example.io]", data.RBAC.ClusterRoles[0].Rules[0].APIGroups)
	}

	// KubebuilderMarkers
	if len(data.RBAC.KubebuilderMarkers) != 1 {
		t.Fatalf("KubebuilderMarkers count = %d, want 1", len(data.RBAC.KubebuilderMarkers))
	}
	if data.RBAC.KubebuilderMarkers[0].File != "controllers/widget_controller.go" {
		t.Errorf("Marker file = %q, want %q", data.RBAC.KubebuilderMarkers[0].File, "controllers/widget_controller.go")
	}

	// Webhooks
	if len(data.Webhooks) != 1 {
		t.Fatalf("Webhooks count = %d, want 1", len(data.Webhooks))
	}
	if data.Webhooks[0].Path != "/mutate-apps-v1-widget" {
		t.Errorf("Webhook path = %q, want %q", data.Webhooks[0].Path, "/mutate-apps-v1-widget")
	}
	if len(data.Webhooks[0].Rules) != 1 {
		t.Fatalf("Webhook rules count = %d, want 1", len(data.Webhooks[0].Rules))
	}
	if len(data.Webhooks[0].Rules[0].APIGroups) != 1 || data.Webhooks[0].Rules[0].APIGroups[0] != "apps.example.io" {
		t.Errorf("Webhook rule APIGroups = %v, want [apps.example.io]", data.Webhooks[0].Rules[0].APIGroups)
	}
	if len(data.Webhooks[0].Rules[0].APIVersions) != 1 {
		t.Errorf("Webhook rule APIVersions count = %d, want 1", len(data.Webhooks[0].Rules[0].APIVersions))
	}

	// Secrets
	if len(data.Secrets) != 1 {
		t.Fatalf("Secrets count = %d, want 1", len(data.Secrets))
	}
	if data.Secrets[0].Name != "webhook-cert" {
		t.Errorf("Secret name = %q, want %q", data.Secrets[0].Name, "webhook-cert")
	}
	if len(data.Secrets[0].ReferencedBy) != 1 {
		t.Errorf("Secret.ReferencedBy count = %d, want 1", len(data.Secrets[0].ReferencedBy))
	}
	if data.Secrets[0].ProvisionedBy != "cert-manager" {
		t.Errorf("Secret.ProvisionedBy = %q, want %q", data.Secrets[0].ProvisionedBy, "cert-manager")
	}

	// Cache
	if len(data.Cache.Issues) != 1 {
		t.Errorf("Cache issues count = %d, want 1", len(data.Cache.Issues))
	}
	if len(data.Cache.FilteredTypes) != 1 {
		t.Fatalf("Cache filtered types count = %d, want 1", len(data.Cache.FilteredTypes))
	}
	if data.Cache.FilteredTypes[0].Type != "Secret" {
		t.Errorf("FilteredType[0].Type = %q, want %q", data.Cache.FilteredTypes[0].Type, "Secret")
	}
	if data.Cache.FilteredTypes[0].FilterKind != "label" {
		t.Errorf("FilteredType[0].FilterKind = %q, want %q", data.Cache.FilteredTypes[0].FilterKind, "label")
	}
}

func TestParseNilRaw(t *testing.T) {
	data, err := Parse(nil)
	if err != nil {
		t.Fatalf("Parse(nil) error: %v", err)
	}
	if data.Component != "" {
		t.Errorf("Component = %q, want empty", data.Component)
	}
	if len(data.CRDs) != 0 {
		t.Errorf("CRDs = %d, want 0", len(data.CRDs))
	}
}

func TestParseEmptyRaw(t *testing.T) {
	data, err := Parse(map[string]interface{}{})
	if err != nil {
		t.Fatalf("Parse({}) error: %v", err)
	}
	if data.Component != "" {
		t.Errorf("Component = %q, want empty", data.Component)
	}
}
