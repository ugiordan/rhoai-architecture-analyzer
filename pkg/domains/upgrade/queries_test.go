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
		Annotations: map[string]bool{AnnotPreReleaseAPI: true},
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
		Annotations: map[string]bool{AnnotPreReleaseAPI: true},
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

// --- CGA-U03 tests ---

func TestQueryUngatedFeature_NoArchData(t *testing.T) {
	g := graph.NewCPG()
	findings := queryUngatedFeature(g)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings without arch data, got %d", len(findings))
	}
}

func TestQueryUngatedFeature_NoGatesRegistered(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{
		FeatureGates: []arch.FeatureGate{},
	}
	// Even with a feature gate call site, if no gates are registered,
	// the project likely doesn't use feature gates so we skip.
	cs := &graph.Node{
		ID:          "call1",
		Kind:        graph.NodeCallSite,
		Name:        "utilfeature.DefaultFeatureGate.Enabled",
		File:        "controller.go",
		Line:        50,
		Annotations: map[string]bool{AnnotFeatureGate: true},
		Properties:  make(map[string]string),
	}
	g.AddNode(cs)

	findings := queryUngatedFeature(g)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings when no gates registered, got %d", len(findings))
	}
}

func TestQueryUngatedFeature_UnregisteredGateCheck(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{
		FeatureGates: []arch.FeatureGate{
			{Name: "PipelineReuse", Default: true, PreRelease: "Beta"},
		},
	}
	// Call site referencing a gate NOT in the inventory
	cs := &graph.Node{
		ID:          "call1",
		Kind:        graph.NodeCallSite,
		Name:        "utilfeature.DefaultFeatureGate.Enabled",
		File:        "controller.go",
		Line:        50,
		Annotations: map[string]bool{AnnotFeatureGate: true},
		Properties:  make(map[string]string),
	}
	g.AddNode(cs)

	findings := queryUngatedFeature(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for unregistered gate check, got %d", len(findings))
	}
	if findings[0].RuleID != "CGA-U03" {
		t.Errorf("RuleID = %q, want CGA-U03", findings[0].RuleID)
	}
}

func TestQueryUngatedFeature_RegisteredGateInProperties(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{
		FeatureGates: []arch.FeatureGate{
			{Name: "PipelineReuse", Default: true, PreRelease: "Beta"},
		},
	}
	// Call site with gate name in properties (argument context)
	cs := &graph.Node{
		ID:          "call1",
		Kind:        graph.NodeCallSite,
		Name:        "utilfeature.DefaultFeatureGate.Enabled",
		File:        "controller.go",
		Line:        50,
		Annotations: map[string]bool{AnnotFeatureGate: true},
		Properties:  map[string]string{"arg": "PipelineReuse"},
	}
	g.AddNode(cs)

	findings := queryUngatedFeature(g)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings when gate is registered and in properties, got %d", len(findings))
	}
}

func TestQueryUngatedFeature_FeatureNamedFunctionWithoutGate(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{
		FeatureGates: []arch.FeatureGate{
			{Name: "PipelineReuse", Default: true},
		},
	}
	// Function with feature-related name but no gate check
	fn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "enableFeatureX",
		File:        "controller.go",
		Line:        10,
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(fn)

	findings := queryUngatedFeature(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for ungated feature function, got %d", len(findings))
	}
	if findings[0].RuleID != "CGA-U03" {
		t.Errorf("RuleID = %q, want CGA-U03", findings[0].RuleID)
	}
}

func TestQueryUngatedFeature_FeatureNamedFunctionWithGate(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{
		FeatureGates: []arch.FeatureGate{
			{Name: "PipelineReuse", Default: true},
		},
	}
	// Function with feature-related name
	fn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "enableFeatureX",
		File:        "controller.go",
		Line:        10,
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(fn)

	// Gate check call inside that function
	cs := &graph.Node{
		ID:          "call1",
		Kind:        graph.NodeCallSite,
		Name:        "utilfeature.DefaultFeatureGate.Enabled",
		File:        "controller.go",
		Line:        15,
		Annotations: map[string]bool{AnnotFeatureGate: true},
		Properties:  map[string]string{"arg": "PipelineReuse"},
	}
	g.AddNode(cs)

	// Edge from function to call site (function contains the call)
	g.AddEdge(&graph.Edge{From: fn.ID, To: cs.ID, Kind: "contains"})

	findings := queryUngatedFeature(g)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings when feature function has gate check, got %d", len(findings))
	}
}

func TestQueryUngatedFeature_SkipsTestFiles(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{
		FeatureGates: []arch.FeatureGate{
			{Name: "PipelineReuse", Default: true},
		},
	}
	cs := &graph.Node{
		ID:          "call1",
		Kind:        graph.NodeCallSite,
		Name:        "utilfeature.DefaultFeatureGate.Enabled",
		File:        "controller_test.go",
		Line:        50,
		Annotations: map[string]bool{AnnotFeatureGate: true},
		Properties:  make(map[string]string),
	}
	g.AddNode(cs)

	findings := queryUngatedFeature(g)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for test file, got %d", len(findings))
	}
}

func TestIsFeatureRelatedFunction(t *testing.T) {
	positives := []string{
		"enablefeaturex", "disablefeaturey",
		"setupfeaturegatez", "initfeatureflags",
		"togglefeature", "activatefeaturemode",
		"configurefeaturepipeline", "registerfeaturegate",
		"turnonfeaturetest", "turnofffeatureold",
		"withfeatureflag", "handlefeaturerequest",
		"featureflagcheck", "featuregateinit",
	}
	for _, name := range positives {
		if !isFeatureRelatedFunction(name) {
			t.Errorf("isFeatureRelatedFunction(%q) = false, want true", name)
		}
	}

	negatives := []string{
		"reconcile", "handlerequest", "setupcontroller",
		"initmanager", "enablelogging", "configureauth",
	}
	for _, name := range negatives {
		if isFeatureRelatedFunction(name) {
			t.Errorf("isFeatureRelatedFunction(%q) = true, want false", name)
		}
	}
}

func TestQueryUncheckedVersionAccess_SeverityIsLow(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "checkVersion",
		File:        "upgrade.go",
		Line:        10,
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(fn)
	cs := &graph.Node{
		ID:          "call1",
		Kind:        graph.NodeCallSite,
		Name:        "semver.Compare",
		File:        "upgrade.go",
		Line:        15,
		Annotations: map[string]bool{AnnotVersionCheck: true},
		Properties:  make(map[string]string),
	}
	g.AddNode(cs)
	g.AddEdge(&graph.Edge{From: fn.ID, To: cs.ID, Kind: "contains"})

	findings := queryUncheckedVersionAccess(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != "low" {
		t.Errorf("severity = %q, want low (advisory)", findings[0].Severity)
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
