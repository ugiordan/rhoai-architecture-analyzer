package upgrade

import (
	"testing"

	"github.com/ugiordan/rhoai-architecture-analyzer/pkg/arch"
	"github.com/ugiordan/rhoai-architecture-analyzer/pkg/graph"
)

func TestQueryDeprecatedAPIUsage(t *testing.T) {
	g := graph.NewCPG()
	cs := &graph.Node{
		ID:          "call1",
		Kind:        graph.NodeCallSite,
		Name:        "v1alpha1.NewInferenceService",
		File:        "controller.go",
		Line:        20,
		Annotations: map[string]bool{AnnotDeprecatedAPI: true},
		Properties:  make(map[string]string),
	}
	g.AddNode(cs)

	findings := queryDeprecatedAPIUsage(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].RuleID != "CGA-U02" {
		t.Errorf("expected CGA-U02, got %s", findings[0].RuleID)
	}
}

func TestQueryDeprecatedAPISkipsTests(t *testing.T) {
	g := graph.NewCPG()
	cs := &graph.Node{
		ID:          "call1",
		Kind:        graph.NodeCallSite,
		Name:        "v1alpha1.NewInferenceService",
		File:        "controller_test.go",
		Line:        20,
		Annotations: map[string]bool{AnnotDeprecatedAPI: true},
		Properties:  make(map[string]string),
	}
	g.AddNode(cs)

	findings := queryDeprecatedAPIUsage(g)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for test file, got %d", len(findings))
	}
}

func TestQueryUnconvertedCRD_NoArchData(t *testing.T) {
	g := graph.NewCPG()
	findings := queryUnconvertedCRD(g)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings without arch data, got %d", len(findings))
	}
}

func TestQueryUnconvertedCRD_SingleVersionCRD(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{
		CRDs: []arch.CRD{
			{Group: "apps.example.io", Version: "v1", Kind: "Widget", Source: "config/crd/widgets.yaml"},
		},
	}
	findings := queryUnconvertedCRD(g)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for single-version CRD, got %d", len(findings))
	}
}

func TestQueryUnconvertedCRD_MultiVersionWithConversion(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{
		CRDs: []arch.CRD{
			{Group: "apps.example.io", Version: "v1alpha1", Kind: "Widget", Source: "config/crd/widgets.yaml"},
			{Group: "apps.example.io", Version: "v1", Kind: "Widget", Source: "config/crd/widgets.yaml"},
		},
	}
	fn := &graph.Node{
		ID:   "fn-convert-to",
		Kind: graph.NodeFunction,
		Name: "ConvertTo",
		File: "api/v1alpha1/widget_conversion.go",
		Line: 10,
		Properties: map[string]string{
			"receiver": "(w *Widget)",
		},
		Annotations: map[string]bool{
			AnnotVersionConversion: true,
		},
	}
	g.AddNode(fn)

	findings := queryUnconvertedCRD(g)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for multi-version CRD with conversion, got %d", len(findings))
	}
}

func TestQueryUnconvertedCRD_MultiVersionWithoutConversion(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{
		CRDs: []arch.CRD{
			{Group: "apps.example.io", Version: "v1alpha1", Kind: "Widget", Source: "config/crd/widgets.yaml"},
			{Group: "apps.example.io", Version: "v1", Kind: "Widget", Source: "config/crd/widgets.yaml"},
		},
	}
	findings := queryUnconvertedCRD(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for multi-version CRD without conversion, got %d", len(findings))
	}
	if findings[0].RuleID != "CGA-U01" {
		t.Errorf("RuleID = %q, want CGA-U01", findings[0].RuleID)
	}
	if findings[0].ArchitectureRef == "" {
		t.Error("expected ArchitectureRef to be set")
	}
}

func TestQueryUnconvertedCRD_PackageQualifiedReceiver(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{
		CRDs: []arch.CRD{
			{Group: "apps.example.io", Version: "v1alpha1", Kind: "Widget", Source: "config/crd/widgets.yaml"},
			{Group: "apps.example.io", Version: "v1", Kind: "Widget", Source: "config/crd/widgets.yaml"},
		},
	}
	fn := &graph.Node{
		ID:   "fn-convert-to",
		Kind: graph.NodeFunction,
		Name: "ConvertTo",
		File: "api/v1alpha1/widget_conversion.go",
		Line: 10,
		Properties: map[string]string{
			"receiver": "(w *v1alpha1.Widget)",
		},
		Annotations: map[string]bool{
			AnnotVersionConversion: true,
		},
	}
	g.AddNode(fn)

	findings := queryUnconvertedCRD(g)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for package-qualified receiver, got %d", len(findings))
	}
}

func TestExtractReceiverType(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"(w *Widget)", "Widget"},
		{"(r Widget)", "Widget"},
		{"(w *v1alpha1.Widget)", "Widget"},
		{"(w *pkg.Widget)", "Widget"},
		{"(*Widget)", "Widget"},
		{"", ""},
		{"()", ""},
	}
	for _, tt := range tests {
		got := extractReceiverType(tt.input)
		if got != tt.want {
			t.Errorf("extractReceiverType(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
