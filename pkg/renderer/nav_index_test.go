package renderer

import (
	"strings"
	"testing"
)

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"kserve", "kserve"},
		{"vllm-cpu", "vllm-cpu"},
		{"openvino_model_server", "openvino_model_server"},
		{"My/Component", "my-component"},
		{"UPPER-case", "upper-case"},
		{"dots.in.name", "dots.in.name"},
		{"special@chars!", "special-chars"},
		{"multi--dash", "multi-dash"},
		{"", "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := SanitizeFilename(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeFilename(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestRenderNavIndex(t *testing.T) {
	platformData := map[string]interface{}{
		"platform":        "OpenShift AI",
		"aggregated_at":   "2026-05-08T00:00:00Z",
		"component_count": float64(2),
		"components":      []interface{}{"kserve", "vllm"},
		"component_data": []interface{}{
			map[string]interface{}{
				"component": "kserve",
				"aliases":   []interface{}{"KServe"},
				"crds":      []interface{}{map[string]interface{}{"kind": "InferenceService"}},
				"services":  []interface{}{map[string]interface{}{"name": "svc1"}, map[string]interface{}{"name": "svc2"}},
			},
			map[string]interface{}{
				"component": "vllm",
			},
		},
	}

	result := RenderNavIndex(platformData)

	if !strings.Contains(result, "# Architecture Analyzer Output") {
		t.Error("missing heading")
	}
	if !strings.Contains(result, "## How to Find Information") {
		t.Error("missing routing table")
	}
	if !strings.Contains(result, "## Component Index") {
		t.Error("missing component index")
	}
	if !strings.Contains(result, "kserve") {
		t.Error("missing kserve in index")
	}
	if !strings.Contains(result, "KServe") {
		t.Error("missing alias in index")
	}
	if !strings.Contains(result, autoGenHeader) {
		t.Error("missing auto-generated header")
	}
}
