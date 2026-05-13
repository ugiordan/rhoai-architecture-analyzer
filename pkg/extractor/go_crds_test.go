package extractor

import "testing"

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
