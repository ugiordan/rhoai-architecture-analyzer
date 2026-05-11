// pkg/renderer/interactions_test.go
package renderer

import (
	"strings"
	"testing"
)

func TestRenderInteractions(t *testing.T) {
	platformData := map[string]interface{}{
		"dependency_graph": []interface{}{
			map[string]interface{}{"from": "codeflare-operator", "to": "kuberay", "type": "go-module"},
			map[string]interface{}{"from": "kserve", "to": "kube-rbac-proxy", "type": "uses-image"},
			map[string]interface{}{"from": "notebooks", "to": "opendatahub-operator", "type": "go-module"},
		},
		"component_data": []interface{}{
			map[string]interface{}{
				"component": "kserve",
				"component_refs": []interface{}{
					map[string]interface{}{"target": "vllm", "type": "provider", "source": "providers/vllm/"},
				},
			},
			map[string]interface{}{
				"component": "codeflare-operator",
			},
		},
	}

	result := RenderInteractions(platformData)

	if !strings.Contains(result, "# Cross-Component Interactions") {
		t.Error("missing heading")
	}
	if !strings.Contains(result, "## All Interactions") {
		t.Error("missing All Interactions table")
	}
	if !strings.Contains(result, "codeflare-operator") {
		t.Error("missing codeflare-operator in interactions")
	}
	if !strings.Contains(result, "kuberay") {
		t.Error("missing kuberay in interactions")
	}
	if !strings.Contains(result, "## Per-Component View") {
		t.Error("missing Per-Component View")
	}
	if !strings.Contains(result, "vllm") {
		t.Error("missing vllm from component_refs in interactions")
	}
}

func TestRenderInteractions_Empty(t *testing.T) {
	platformData := map[string]interface{}{}
	result := RenderInteractions(platformData)
	if !strings.Contains(result, "0 interactions") {
		t.Error("should show 0 interactions for empty data")
	}
}
