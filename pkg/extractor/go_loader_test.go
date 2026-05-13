package extractor

import (
	"path/filepath"
	"runtime"
	"testing"
)

func fixtureDir() string {
	_, f, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(f), "..", "..", "testdata", "go-ast-fixture")
}

func TestLoadGoPackages_FullMode(t *testing.T) {
	pkgs := loadGoPackages(fixtureDir())
	if pkgs == nil {
		t.Fatal("expected non-nil GoPackageSet")
	}
	if pkgs.Mode != "full" {
		t.Errorf("expected mode=full, got %s (warning: %s)", pkgs.Mode, pkgs.Warning)
	}
	if len(pkgs.Packages) == 0 {
		t.Fatal("expected at least one package loaded")
	}
}

func TestLoadGoPackages_NotGoRepo(t *testing.T) {
	pkgs := loadGoPackages(t.TempDir())
	if pkgs != nil {
		t.Error("expected nil for non-Go repo")
	}
}

func TestLoadGoPackages_FindStructsWithMarker(t *testing.T) {
	pkgs := loadGoPackages(fixtureDir())
	if pkgs == nil {
		t.Fatal("expected non-nil GoPackageSet")
	}
	structs := pkgs.FindStructsWithMarker("+kubebuilder:object:root=true")
	if len(structs) == 0 {
		t.Fatal("expected to find Widget root type")
	}
	found := false
	for _, s := range structs {
		if s.Name == "Widget" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected to find Widget, got %v", structs)
	}
}
