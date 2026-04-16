package extractor

import (
	"os"
	"path/filepath"
	"testing"
)

const sampleFeatureGateFile = `package features

import (
	"k8s.io/component-base/featuregate"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
)

const (
	// PipelineReuse enables pipeline object reuse across runs.
	PipelineReuse featuregate.Feature = "PipelineReuse"

	// CacheOptimization enables the optimized cache layer.
	CacheOptimization featuregate.Feature = "CacheOptimization"

	// ExperimentalAPI enables experimental API endpoints.
	ExperimentalAPI featuregate.Feature = "ExperimentalAPI"
)

func init() {
	utilfeature.DefaultMutableFeatureGate.Add(map[featuregate.Feature]featuregate.FeatureSpec{
		PipelineReuse:     {Default: true, PreRelease: featuregate.Beta},
		CacheOptimization: {Default: false, PreRelease: featuregate.Alpha},
		ExperimentalAPI:   {Default: false, PreRelease: featuregate.Alpha},
	})
}
`

const sampleFeatureGateNoConst = `package features

import (
	"k8s.io/component-base/featuregate"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
)

func init() {
	utilfeature.DefaultMutableFeatureGate.Add(map[featuregate.Feature]featuregate.FeatureSpec{
		"InlineGate": {Default: true, PreRelease: featuregate.GA},
	})
}
`

const sampleFeatureGateRuntimeSet = `package main

import (
	utilfeature "k8s.io/apiserver/pkg/util/feature"
)

func enableDebugMode() {
	utilfeature.DefaultMutableFeatureGate.Set("DebugMode=true")
}
`

func TestExtractFeatureGates(t *testing.T) {
	root := t.TempDir()
	pkgDir := filepath.Join(root, "pkg", "features")
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pkgDir, "gates.go"), []byte(sampleFeatureGateFile), 0o644); err != nil {
		t.Fatal(err)
	}

	gates := extractFeatureGates(root)
	if len(gates) != 3 {
		t.Fatalf("expected 3 feature gates, got %d", len(gates))
	}

	gateMap := make(map[string]FeatureGate)
	for _, g := range gates {
		gateMap[g.Name] = g
	}

	// PipelineReuse: default=true, Beta
	pr, ok := gateMap["PipelineReuse"]
	if !ok {
		t.Fatal("missing PipelineReuse gate")
	}
	if !pr.Default {
		t.Error("PipelineReuse should default to true")
	}
	if pr.PreRelease != "Beta" {
		t.Errorf("PipelineReuse PreRelease = %q, want Beta", pr.PreRelease)
	}

	// CacheOptimization: default=false, Alpha
	co, ok := gateMap["CacheOptimization"]
	if !ok {
		t.Fatal("missing CacheOptimization gate")
	}
	if co.Default {
		t.Error("CacheOptimization should default to false")
	}
	if co.PreRelease != "Alpha" {
		t.Errorf("CacheOptimization PreRelease = %q, want Alpha", co.PreRelease)
	}

	// ExperimentalAPI: default=false, Alpha
	ea, ok := gateMap["ExperimentalAPI"]
	if !ok {
		t.Fatal("missing ExperimentalAPI gate")
	}
	if ea.Default {
		t.Error("ExperimentalAPI should default to false")
	}
}

func TestExtractFeatureGates_RuntimeSet(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "main.go"), []byte(sampleFeatureGateRuntimeSet), 0o644); err != nil {
		t.Fatal(err)
	}

	gates := extractFeatureGates(root)
	if len(gates) != 1 {
		t.Fatalf("expected 1 feature gate, got %d", len(gates))
	}
	if gates[0].Name != "DebugMode" {
		t.Errorf("gate name = %q, want DebugMode", gates[0].Name)
	}
	if !gates[0].RuntimeSet {
		t.Error("expected RuntimeSet=true for Set() call")
	}
}

func TestExtractFeatureGates_NoGates(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\nfunc main() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	gates := extractFeatureGates(root)
	if len(gates) != 0 {
		t.Errorf("expected 0 gates, got %d", len(gates))
	}
}

func TestExtractFeatureGates_SkipsTestFiles(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "gates_test.go"), []byte(sampleFeatureGateFile), 0o644); err != nil {
		t.Fatal(err)
	}

	gates := extractFeatureGates(root)
	if len(gates) != 0 {
		t.Errorf("expected 0 gates from test file, got %d", len(gates))
	}
}

func TestExtractFeatureGates_Dedup(t *testing.T) {
	root := t.TempDir()
	pkgDir := filepath.Join(root, "pkg", "features")
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Write the same gate definitions in two files
	if err := os.WriteFile(filepath.Join(pkgDir, "gates.go"), []byte(sampleFeatureGateFile), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pkgDir, "gates2.go"), []byte(sampleFeatureGateFile), 0o644); err != nil {
		t.Fatal(err)
	}

	gates := extractFeatureGates(root)
	// Should still be 3 unique gates, not 6
	if len(gates) != 3 {
		t.Errorf("expected 3 deduplicated gates, got %d", len(gates))
	}
}

func TestExtractFeatureGates_InlineKeys(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "gates.go"), []byte(sampleFeatureGateNoConst), 0o644); err != nil {
		t.Fatal(err)
	}

	gates := extractFeatureGates(root)
	if len(gates) != 1 {
		t.Fatalf("expected 1 feature gate from inline key, got %d", len(gates))
	}
	if gates[0].Name != "InlineGate" {
		t.Errorf("gate name = %q, want InlineGate", gates[0].Name)
	}
	if !gates[0].Default {
		t.Error("InlineGate should default to true")
	}
	if gates[0].PreRelease != "GA" {
		t.Errorf("InlineGate PreRelease = %q, want GA", gates[0].PreRelease)
	}
}

func TestExtractFeatureGates_MixedConstAndInline(t *testing.T) {
	// Edge case: same file has both const-referenced and inline string keys
	mixed := `package features

import (
	"k8s.io/component-base/featuregate"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
)

const ConstGate featuregate.Feature = "ConstGate"

func init() {
	utilfeature.DefaultMutableFeatureGate.Add(map[featuregate.Feature]featuregate.FeatureSpec{
		ConstGate:      {Default: true, PreRelease: featuregate.Beta},
		"InlineGate":   {Default: false, PreRelease: featuregate.Alpha},
	})
}
`
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "gates.go"), []byte(mixed), 0o644); err != nil {
		t.Fatal(err)
	}

	gates := extractFeatureGates(root)
	if len(gates) != 2 {
		t.Fatalf("expected 2 gates (const + inline), got %d", len(gates))
	}
	gateMap := make(map[string]FeatureGate)
	for _, g := range gates {
		gateMap[g.Name] = g
	}
	if _, ok := gateMap["ConstGate"]; !ok {
		t.Error("missing ConstGate")
	}
	if _, ok := gateMap["InlineGate"]; !ok {
		t.Error("missing InlineGate")
	}
}

func TestExtractFeatureGates_SkipsSymlinks(t *testing.T) {
	root := t.TempDir()
	realFile := filepath.Join(root, "real.go")
	if err := os.WriteFile(realFile, []byte(sampleFeatureGateNoConst), 0o644); err != nil {
		t.Fatal(err)
	}
	symlink := filepath.Join(root, "link.go")
	if err := os.Symlink(realFile, symlink); err != nil {
		t.Skip("symlinks not supported on this platform")
	}

	gates := extractFeatureGates(root)
	// Should find the gate from real.go but not double-count from link.go
	if len(gates) != 1 {
		t.Errorf("expected 1 gate (symlink should be skipped), got %d", len(gates))
	}
}

func TestExtractFeatureGates_SkipsVendor(t *testing.T) {
	root := t.TempDir()
	vendorDir := filepath.Join(root, "vendor", "k8s.io", "features")
	if err := os.MkdirAll(vendorDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(vendorDir, "gates.go"), []byte(sampleFeatureGateFile), 0o644); err != nil {
		t.Fatal(err)
	}

	gates := extractFeatureGates(root)
	if len(gates) != 0 {
		t.Errorf("expected 0 gates from vendor dir, got %d", len(gates))
	}
}

func TestExtractFeatureGates_RuntimeSetAndRegistration(t *testing.T) {
	// Edge case: same gate registered and also Set() at runtime
	combined := `package features

import (
	"k8s.io/component-base/featuregate"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
)

const DebugMode featuregate.Feature = "DebugMode"

func init() {
	utilfeature.DefaultMutableFeatureGate.Add(map[featuregate.Feature]featuregate.FeatureSpec{
		DebugMode: {Default: false, PreRelease: featuregate.Alpha},
	})
}

func enableDebug() {
	utilfeature.DefaultMutableFeatureGate.Set("DebugMode=true")
}
`
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "gates.go"), []byte(combined), 0o644); err != nil {
		t.Fatal(err)
	}

	gates := extractFeatureGates(root)
	// DebugMode should appear only once (dedup), not twice
	if len(gates) != 1 {
		t.Fatalf("expected 1 deduplicated gate, got %d", len(gates))
	}
	if gates[0].Name != "DebugMode" {
		t.Errorf("gate name = %q, want DebugMode", gates[0].Name)
	}
}

func TestExtractFeatureGates_SourceLocation(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "features.go"), []byte(sampleFeatureGateFile), 0o644); err != nil {
		t.Fatal(err)
	}

	gates := extractFeatureGates(root)
	for _, g := range gates {
		if g.Source == "" {
			t.Errorf("gate %s has empty source", g.Name)
		}
	}
}
