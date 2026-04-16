package renderer

import (
	"strings"
	"testing"
)

// sampleData returns minimal architecture data for testing.
func sampleData() map[string]interface{} {
	return map[string]interface{}{
		"component": "test-controller",
		"crds": []interface{}{
			map[string]interface{}{
				"kind":    "Notebook",
				"group":   "kubeflow.org",
				"version": "v1",
			},
		},
		"deployments": []interface{}{
			map[string]interface{}{
				"name":            "test-controller-manager",
				"service_account": "test-sa",
				"containers": []interface{}{
					map[string]interface{}{
						"name":  "manager",
						"image": "quay.io/test:v1",
						"security_context": map[string]interface{}{
							"runAsNonRoot": true,
						},
					},
				},
			},
		},
		"controller_watches": []interface{}{
			map[string]interface{}{"gvk": "kubeflow.org/v1/Notebook", "type": "For"},
			map[string]interface{}{"gvk": "apps/v1/Deployment", "type": "Owns"},
			map[string]interface{}{"gvk": "v1/ConfigMap", "type": "Watches"},
		},
		"services": []interface{}{
			map[string]interface{}{
				"name": "webhook-service",
				"type": "ClusterIP",
				"ports": []interface{}{
					map[string]interface{}{
						"name":       "https",
						"port":       443,
						"targetPort": 9443,
						"protocol":   "TCP",
					},
				},
			},
		},
		"rbac": map[string]interface{}{
			"cluster_roles": []interface{}{
				map[string]interface{}{
					"name": "test-role",
					"rules": []interface{}{
						map[string]interface{}{
							"apiGroups": []interface{}{"kubeflow.org"},
							"resources": []interface{}{"notebooks"},
							"verbs":     []interface{}{"get", "list", "watch"},
						},
					},
				},
			},
			"cluster_role_bindings": []interface{}{
				map[string]interface{}{
					"name":     "test-binding",
					"role_ref": "test-role",
					"subjects": []interface{}{
						map[string]interface{}{
							"kind":      "ServiceAccount",
							"name":      "test-sa",
							"namespace": "test-ns",
						},
					},
				},
			},
			"roles":                []interface{}{},
			"role_bindings":        []interface{}{},
			"kubebuilder_markers":  []interface{}{},
		},
		"secrets_referenced": []interface{}{
			map[string]interface{}{
				"name":          "webhook-cert",
				"type":          "kubernetes.io/tls",
				"referenced_by": []interface{}{"test-controller-manager"},
				"provisioned_by": "cert-manager",
			},
		},
		"dependencies": map[string]interface{}{
			"internal_odh": []interface{}{
				map[string]interface{}{
					"component":   "odh-dashboard",
					"interaction": "watches CRDs",
				},
			},
			"go_modules": []interface{}{
				map[string]interface{}{
					"module":  "sigs.k8s.io/controller-runtime",
					"version": "v0.17.0",
				},
			},
		},
		"dockerfiles": []interface{}{
			map[string]interface{}{
				"path":       "Dockerfile",
				"base_image": "golang:1.21",
				"user":       "65532",
				"stages":     2,
				"issues":     []interface{}{},
			},
		},
	}
}

func TestSanitizeID(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "hello"},
		{"hello-world", "hello_world"},
		{"123abc", "n_123abc"},
		{"", "node"},
		{"a.b/c:d", "a_b_c_d"},
	}
	for _, tt := range tests {
		got := sanitizeID(tt.input)
		if got != tt.want {
			t.Errorf("sanitizeID(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestEscapeLabel(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`hello "world"`, "hello 'world'"},
		{"a<b>c", "a&lt;b&gt;c"},
		{"plain", "plain"},
	}
	for _, tt := range tests {
		got := escapeLabel(tt.input)
		if got != tt.want {
			t.Errorf("escapeLabel(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestRBACRenderer(t *testing.T) {
	r := &RBACRenderer{}
	if r.Filename() != "rbac.mmd" {
		t.Errorf("unexpected filename: %s", r.Filename())
	}
	out := r.Render(sampleData())
	if !strings.HasPrefix(out, "graph TD") {
		t.Error("RBAC output should start with 'graph TD'")
	}
	if !strings.Contains(out, "test-role") {
		t.Error("RBAC output should contain role name")
	}
	if !strings.Contains(out, "classDef sa") {
		t.Error("RBAC output should contain style classes")
	}
}

func TestRBACRenderer_Empty(t *testing.T) {
	out := (&RBACRenderer{}).Render(map[string]interface{}{"component": "x"})
	if !strings.Contains(out, "No RBAC data") {
		t.Error("empty RBAC should show no-data note")
	}
}

func TestComponentRenderer(t *testing.T) {
	r := &ComponentRenderer{}
	if r.Filename() != "component.mmd" {
		t.Errorf("unexpected filename: %s", r.Filename())
	}
	out := r.Render(sampleData())
	if !strings.HasPrefix(out, "graph LR") {
		t.Error("component output should start with 'graph LR'")
	}
	if !strings.Contains(out, "subgraph controller") {
		t.Error("should contain controller subgraph")
	}
	if !strings.Contains(out, "For (reconciles)") {
		t.Error("should contain For watch type")
	}
	if !strings.Contains(out, "Owns") {
		t.Error("should contain Owns relationship")
	}
}

func TestSecurityNetworkRenderer(t *testing.T) {
	r := &SecurityNetworkRenderer{}
	if r.Filename() != "security-network.txt" {
		t.Errorf("unexpected filename: %s", r.Filename())
	}
	out := r.Render(sampleData())
	if !strings.Contains(out, "SECURITY & NETWORK ARCHITECTURE") {
		t.Error("should contain title")
	}
	if !strings.Contains(out, "NETWORK TOPOLOGY") {
		t.Error("should contain network topology section")
	}
	if !strings.Contains(out, "RBAC SUMMARY") {
		t.Error("should contain RBAC summary")
	}
	if !strings.Contains(out, "SECRETS INVENTORY") {
		t.Error("should contain secrets inventory")
	}
	if !strings.Contains(out, "webhook-cert") {
		t.Error("should contain secret name")
	}
	if !strings.Contains(out, "DEPLOYMENT SECURITY CONTROLS") {
		t.Error("should contain deployment security section")
	}
	if !strings.Contains(out, "DOCKERFILE SECURITY") {
		t.Error("should contain dockerfile section")
	}
}

func TestDependencyRenderer(t *testing.T) {
	r := &DependencyRenderer{}
	if r.Filename() != "dependencies.mmd" {
		t.Errorf("unexpected filename: %s", r.Filename())
	}
	out := r.Render(sampleData())
	if !strings.HasPrefix(out, "graph LR") {
		t.Error("dependency output should start with 'graph LR'")
	}
	if !strings.Contains(out, "odh-dashboard") {
		t.Error("should contain internal ODH dependency")
	}
	if !strings.Contains(out, "controller-runtime") {
		t.Error("should contain notable external dependency")
	}
}

func TestC4Renderer(t *testing.T) {
	r := &C4Renderer{}
	if r.Filename() != "c4-context.dsl" {
		t.Errorf("unexpected filename: %s", r.Filename())
	}
	out := r.Render(sampleData())
	if !strings.HasPrefix(out, "workspace {") {
		t.Error("C4 output should start with 'workspace {'")
	}
	if !strings.Contains(out, "model {") {
		t.Error("should contain model block")
	}
	if !strings.Contains(out, "views {") {
		t.Error("should contain views block")
	}
	if !strings.Contains(out, "Platform Admin") {
		t.Error("should contain admin person")
	}
	if !strings.Contains(out, "softwareSystem") {
		t.Error("should contain software system")
	}
}

func TestDataflowRenderer(t *testing.T) {
	r := &DataflowRenderer{}
	if r.Filename() != "dataflow.mmd" {
		t.Errorf("unexpected filename: %s", r.Filename())
	}
	out := r.Render(sampleData())
	if !strings.HasPrefix(out, "sequenceDiagram") {
		t.Error("dataflow output should start with 'sequenceDiagram'")
	}
	if !strings.Contains(out, "participant KubernetesAPI") {
		t.Error("should contain KubernetesAPI participant")
	}
	if !strings.Contains(out, "Watch Notebook (reconcile)") {
		t.Error("should contain For watch")
	}
	if !strings.Contains(out, "Create/Update Deployment") {
		t.Error("should contain Owns create/update")
	}
}

func TestRenderAll_Empty(t *testing.T) {
	results := RenderAll(sampleData(), nil)
	if len(results) != 7 {
		t.Errorf("RenderAll with nil formats should return 7 renderers, got %d", len(results))
	}
}

func TestRenderAll_Selective(t *testing.T) {
	results := RenderAll(sampleData(), []string{"rbac", "c4"})
	if len(results) != 2 {
		t.Errorf("RenderAll with 2 formats should return 2 results, got %d", len(results))
	}
	if _, ok := results["rbac.mmd"]; !ok {
		t.Error("should contain rbac.mmd")
	}
	if _, ok := results["c4-context.dsl"]; !ok {
		t.Error("should contain c4-context.dsl")
	}
}

func TestRenderPlatformAll(t *testing.T) {
	platformData := map[string]interface{}{
		"components": []interface{}{"comp-a", "comp-b"},
		"dependency_graph": []interface{}{
			map[string]interface{}{
				"from": "comp-a",
				"to":   "comp-b",
				"type": "go-module",
			},
			map[string]interface{}{
				"from": "comp-a",
				"to":   "comp-b",
				"type": "watches-crd:Notebook",
			},
		},
		"crd_ownership": map[string]interface{}{
			"Notebook": "comp-a",
		},
		"rbac_cluster_roles": []interface{}{
			map[string]interface{}{
				"owner": "comp-a",
				"name":  "comp-a-role",
				"rules": []interface{}{
					map[string]interface{}{
						"resources": []interface{}{"pods", "services"},
					},
				},
			},
		},
		"services": []interface{}{
			map[string]interface{}{
				"owner": "comp-a",
				"name":  "svc-a",
				"type":  "ClusterIP",
				"ports": []interface{}{
					map[string]interface{}{"port": 8080, "protocol": "TCP"},
				},
			},
		},
	}

	results := RenderPlatformAll(platformData)
	if len(results) != 5 {
		t.Errorf("expected 5 platform outputs, got %d", len(results))
	}

	depGraph := results["platform-dependencies.mmd"]
	if !strings.Contains(depGraph, "graph LR") {
		t.Error("platform deps should start with graph LR")
	}
	if !strings.Contains(depGraph, "watches Notebook") {
		t.Error("should contain CRD watch edge")
	}

	crdMap := results["platform-crd-ownership.mmd"]
	if !strings.Contains(crdMap, "defines") {
		t.Error("CRD ownership should contain 'defines'")
	}

	rbac := results["platform-rbac-overview.mmd"]
	if !strings.Contains(rbac, "comp-a-role") {
		t.Error("RBAC overview should contain role name")
	}

	network := results["platform-network-topology.mmd"]
	if !strings.Contains(network, "svc-a") {
		t.Error("network topology should contain service name")
	}
}
