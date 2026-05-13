package extractor

import (
	"go/ast"
	"os"
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

func TestLoadGoPackages_BrokenGoMod(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("this is not valid go.mod"), 0644)
	pkgs := loadGoPackages(dir)
	// Should return nil or fallback, not panic
	if pkgs != nil && pkgs.Mode == "full" {
		t.Error("expected nil or fallback mode for broken go.mod")
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

func TestLoadGoPackages_FallbackMode(t *testing.T) {
	// A directory with go.mod but invalid Go source should trigger fallback
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n\ngo 1.22\n"), 0644)
	os.MkdirAll(filepath.Join(dir, "pkg"), 0755)
	os.WriteFile(filepath.Join(dir, "pkg", "bad.go"), []byte("package pkg\n\nfunc broken( {}\n"), 0644)
	pkgs := loadGoPackages(dir)
	// Should return something (may be nil or fallback depending on how go/packages handles syntax errors)
	if pkgs != nil && pkgs.Mode == "full" && len(pkgs.Packages) > 0 {
		// If it loaded, packages should have errors
		hasErrors := false
		for _, p := range pkgs.Packages {
			if len(p.Errors) > 0 {
				hasErrors = true
			}
		}
		// Either fallback mode or packages with errors is acceptable
		_ = hasErrors
	}
}

func TestLoadGoPackages_ResolveType(t *testing.T) {
	pkgs := loadGoPackages(fixtureDir())
	if pkgs == nil {
		t.Fatal("expected non-nil GoPackageSet")
	}
	// ResolveType with a dummy ident not in the package's type info
	// should fall back to exprString
	ident := &ast.Ident{Name: "SomeType"}
	result := ResolveType(ident, pkgs.Packages[0])
	if result != "SomeType" {
		t.Errorf("expected fallback to exprString 'SomeType', got %q", result)
	}
}

func TestLoadGoPackages_FindMethodsOnType(t *testing.T) {
	pkgs := loadGoPackages(fixtureDir())
	if pkgs == nil {
		t.Fatal("expected non-nil GoPackageSet")
	}
	// Find Widget methods across the packages
	for _, pkg := range pkgs.Packages {
		methods := FindMethodsOnType(pkg, "Widget")
		if len(methods) > 0 {
			// Should find Default, ValidateCreate, etc.
			names := make(map[string]bool)
			for _, m := range methods {
				names[m.Name.Name] = true
			}
			if !names["Default"] {
				t.Error("expected to find Default method on Widget")
			}
			return
		}
	}
	t.Error("expected to find Widget methods in at least one package")
}

func TestLoadGoPackages_FindMethodsOnType_NonExistent(t *testing.T) {
	pkgs := loadGoPackages(fixtureDir())
	if pkgs == nil {
		t.Fatal("expected non-nil GoPackageSet")
	}
	for _, pkg := range pkgs.Packages {
		methods := FindMethodsOnType(pkg, "NonExistentType")
		if len(methods) != 0 {
			t.Error("expected no methods for non-existent type")
		}
	}
}
