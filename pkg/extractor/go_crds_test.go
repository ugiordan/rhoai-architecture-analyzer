package extractor

import (
	"testing"
)

func TestExtractCRDsFromGo(t *testing.T) {
	pkgs := loadGoPackages(fixtureDir())
	if pkgs == nil {
		t.Fatal("failed to load fixture packages")
	}
	crds := extractCRDsFromGo(pkgs)
	if len(crds) == 0 {
		t.Fatal("expected at least one CRD from Go source")
	}
	var widget *CRD
	for i := range crds {
		if crds[i].Kind == "Widget" {
			widget = &crds[i]
			break
		}
	}
	if widget == nil {
		t.Fatal("expected Widget CRD")
	}
	if widget.Group != "apps.example.com" {
		t.Errorf("expected group apps.example.com, got %s", widget.Group)
	}
	if widget.Version != "v1alpha1" {
		t.Errorf("expected version v1alpha1, got %s", widget.Version)
	}
	if widget.Scope != "Namespaced" {
		t.Errorf("expected scope Namespaced, got %s", widget.Scope)
	}
	if widget.GoSource != "go_ast" {
		t.Errorf("expected GoSource=go_ast, got %s", widget.GoSource)
	}
	if widget.FieldsCount == 0 {
		t.Error("expected FieldsCount > 0")
	}
}

func TestExtractCRDsFromGo_SkipsList(t *testing.T) {
	pkgs := loadGoPackages(fixtureDir())
	if pkgs == nil {
		t.Fatal("failed to load fixture packages")
	}
	crds := extractCRDsFromGo(pkgs)
	for _, crd := range crds {
		if crd.Kind == "WidgetList" {
			t.Error("should not extract List types as CRDs")
		}
	}
}

func TestExtractCRDsFromGo_NilPackages(t *testing.T) {
	crds := extractCRDsFromGo(nil)
	if len(crds) != 0 {
		t.Error("expected empty CRDs for nil packages")
	}
}

func TestExtractCRDsFromGo_StorageVersion(t *testing.T) {
	pkgs := loadGoPackages(fixtureDir())
	if pkgs == nil {
		t.Fatal("failed to load fixture packages")
	}
	crds := extractCRDsFromGo(pkgs)
	for _, crd := range crds {
		if crd.Kind == "Widget" {
			found := false
			for _, v := range crd.Versions {
				if v.Name == "v1alpha1" && v.Storage {
					found = true
				}
			}
			if !found {
				t.Error("expected v1alpha1 to be marked as storage version")
			}
			return
		}
	}
	t.Fatal("Widget CRD not found")
}

func TestExtractCRDsFromGo_FallbackMode(t *testing.T) {
	// A GoPackageSet with no packages should return empty CRDs
	pkgs := &GoPackageSet{Mode: "fallback", Warning: "test"}
	crds := extractCRDsFromGo(pkgs)
	if len(crds) != 0 {
		t.Error("expected empty CRDs when no packages loaded")
	}
}

func TestExtractCRDsFromGo_FieldCount(t *testing.T) {
	pkgs := loadGoPackages(fixtureDir())
	if pkgs == nil {
		t.Fatal("failed to load fixture packages")
	}
	crds := extractCRDsFromGo(pkgs)
	for _, crd := range crds {
		if crd.Kind == "Widget" {
			// Widget has: TypeMeta (embedded), ObjectMeta (embedded), Spec (with Replicas, Image, GPU), Status (with Ready)
			if crd.FieldsCount < 5 {
				t.Errorf("expected FieldsCount >= 5 for Widget, got %d", crd.FieldsCount)
			}
			return
		}
	}
	t.Fatal("Widget not found")
}

func TestExtractCRDsFromGo_VersionFromPath(t *testing.T) {
	pkgs := loadGoPackages(fixtureDir())
	if pkgs == nil {
		t.Fatal("failed to load fixture packages")
	}
	crds := extractCRDsFromGo(pkgs)
	for _, crd := range crds {
		if crd.Kind == "Widget" {
			if crd.Version != "v1alpha1" {
				t.Errorf("expected version v1alpha1 from package path, got %s", crd.Version)
			}
			return
		}
	}
	t.Fatal("Widget not found")
}

func TestMergeCRDs_GoOnly(t *testing.T) {
	goCRDs := []CRD{{Group: "test.io", Kind: "Foo", Version: "v1", GoSource: "go_ast"}}
	merged := mergeCRDs(nil, goCRDs)
	if len(merged) != 1 {
		t.Fatalf("expected 1 CRD, got %d", len(merged))
	}
	if merged[0].GoSource != "go_ast" {
		t.Errorf("expected GoSource=go_ast, got %s", merged[0].GoSource)
	}
}

func TestMergeCRDs_YAMLOnly(t *testing.T) {
	yamlCRDs := []CRD{{Group: "test.io", Kind: "Foo", Version: "v1", Source: "config/crd"}}
	merged := mergeCRDs(yamlCRDs, nil)
	if len(merged) != 1 {
		t.Fatalf("expected 1 CRD, got %d", len(merged))
	}
	if merged[0].GoSource != "" {
		t.Errorf("expected empty GoSource for YAML-only, got %s", merged[0].GoSource)
	}
}

func TestMergeCRDs_Overlap(t *testing.T) {
	yamlCRDs := []CRD{{Group: "test.io", Kind: "Foo", Version: "v1", FieldsCount: 10, Source: "config/crd"}}
	goCRDs := []CRD{{Group: "test.io", Kind: "Foo", Version: "v1", GoSource: "go_ast", HubVersion: "v1"}}
	merged := mergeCRDs(yamlCRDs, goCRDs)
	if len(merged) != 1 {
		t.Fatalf("expected 1 merged CRD, got %d", len(merged))
	}
	if merged[0].GoSource != "go_ast_enriched" {
		t.Errorf("expected GoSource=go_ast_enriched for overlap, got %s", merged[0].GoSource)
	}
	if merged[0].HubVersion != "v1" {
		t.Errorf("expected HubVersion=v1 from Go source, got %s", merged[0].HubVersion)
	}
	if merged[0].FieldsCount != 10 {
		t.Errorf("expected FieldsCount=10 from YAML (authoritative), got %d", merged[0].FieldsCount)
	}
}

func TestMergeCRDs_Disjoint(t *testing.T) {
	yamlCRDs := []CRD{{Group: "a.io", Kind: "A", Version: "v1"}}
	goCRDs := []CRD{{Group: "b.io", Kind: "B", Version: "v1", GoSource: "go_ast"}}
	merged := mergeCRDs(yamlCRDs, goCRDs)
	if len(merged) != 2 {
		t.Fatalf("expected 2 CRDs (disjoint), got %d", len(merged))
	}
}
