package extractor

import "testing"

func TestExtractResourceOps(t *testing.T) {
	pkgs := loadGoPackages(fixtureDir())
	if pkgs == nil {
		t.Fatal("failed to load fixture packages")
	}
	ops := extractResourceOps(pkgs)
	if len(ops) == 0 {
		t.Fatal("expected ResourceOps from fixture controller")
	}
	kinds := make(map[string]bool)
	for _, op := range ops {
		kinds[op.Kind] = true
	}
	if !kinds["Service"] {
		t.Error("expected ResourceOp for Service")
	}
	if !kinds["Deployment"] {
		t.Error("expected ResourceOp for Deployment")
	}
	for _, op := range ops {
		if op.Verb != "create" {
			t.Errorf("expected verb=create, got %s", op.Verb)
		}
	}
}

func TestExtractResourceOps_NilPackages(t *testing.T) {
	ops := extractResourceOps(nil)
	if len(ops) != 0 {
		t.Error("expected empty for nil packages")
	}
}

func TestExtractResourceOps_FallbackMode(t *testing.T) {
	// An empty GoPackageSet (no packages) should return empty ops
	pkgs := &GoPackageSet{Mode: "fallback"}
	ops := extractResourceOps(pkgs)
	if len(ops) != 0 {
		t.Error("expected empty in fallback mode with no packages")
	}
}

func TestExtractResourceOps_VerbDetection(t *testing.T) {
	pkgs := loadGoPackages(fixtureDir())
	if pkgs == nil {
		t.Fatal("failed to load fixture packages")
	}
	ops := extractResourceOps(pkgs)
	for _, op := range ops {
		if op.Verb == "" {
			t.Error("expected non-empty verb")
		}
		if op.Kind == "" {
			t.Error("expected non-empty kind")
		}
		if op.Source.Type != "go_ast" {
			t.Errorf("expected source type go_ast, got %s", op.Source.Type)
		}
	}
}

func TestSplitTypeToKindGroup_Known(t *testing.T) {
	tests := []struct {
		fullType  string
		wantKind  string
		wantGroup string
		wantK8s   bool
	}{
		{"k8s.io/api/core/v1.Service", "Service", "", true},
		{"k8s.io/api/apps/v1.Deployment", "Deployment", "apps", true},
		{"k8s.io/api/batch/v1.Job", "Job", "batch", true},
		{"k8s.io/api/networking/v1.Ingress", "Ingress", "networking.k8s.io", true},
	}
	for _, tt := range tests {
		kind, group, isK8s := splitTypeToKindGroup(tt.fullType)
		if kind != tt.wantKind || group != tt.wantGroup || isK8s != tt.wantK8s {
			t.Errorf("splitTypeToKindGroup(%q) = (%q,%q,%v), want (%q,%q,%v)",
				tt.fullType, kind, group, isK8s, tt.wantKind, tt.wantGroup, tt.wantK8s)
		}
	}
}

func TestSplitTypeToKindGroup_NonK8s(t *testing.T) {
	_, _, isK8s := splitTypeToKindGroup("github.com/some/pkg.MyType")
	if isK8s {
		t.Error("expected isK8s=false for non-k8s type")
	}
}

func TestSplitTypeToKindGroup_CustomAPI(t *testing.T) {
	kind, _, isK8s := splitTypeToKindGroup("example.com/widget-operator/api/v1alpha1.Widget")
	if !isK8s {
		t.Error("expected isK8s=true for type with /api/ in path")
	}
	if kind != "Widget" {
		t.Errorf("expected kind=Widget, got %s", kind)
	}
}

func TestDedupeResourceOps(t *testing.T) {
	ops := []ResourceOp{
		{Kind: "Service", Group: "", Verb: "create"},
		{Kind: "Service", Group: "", Verb: "create"},
		{Kind: "Deployment", Group: "apps", Verb: "create"},
	}
	deduped := dedupeResourceOps(ops)
	if len(deduped) != 2 {
		t.Errorf("expected 2 after dedup, got %d", len(deduped))
	}
}
