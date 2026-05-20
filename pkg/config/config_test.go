package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadScanConfigOrgGrouped(t *testing.T) {
	dir := t.TempDir()
	content := `
orgs:
  opendatahub-io:
    repos:
      - opendatahub-operator
      - odh-dashboard
  kserve:
    repos:
      - kserve
`
	path := filepath.Join(dir, "scan-config.yaml")
	os.WriteFile(path, []byte(content), 0o644)

	specs, err := LoadScanConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(specs) != 3 {
		t.Fatalf("expected 3 specs, got %d", len(specs))
	}
	// Should be sorted by org
	if specs[0].Org != "kserve" {
		t.Errorf("expected first org 'kserve' (sorted), got %q", specs[0].Org)
	}
}

func TestLoadScanConfigLegacy(t *testing.T) {
	dir := t.TempDir()
	content := `
repos:
  - opendatahub-io/dashboard
  - kserve/kserve
`
	path := filepath.Join(dir, "scan-config.yaml")
	os.WriteFile(path, []byte(content), 0o644)

	specs, err := LoadScanConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(specs) != 2 {
		t.Fatalf("expected 2 specs, got %d", len(specs))
	}
	if specs[0].FullName() != "opendatahub-io/dashboard" {
		t.Errorf("unexpected full name: %q", specs[0].FullName())
	}
}

func TestLoadPlatformConfig(t *testing.T) {
	dir := t.TempDir()
	content := `
platforms:
  odh:
    name: "Open Data Hub"
    orgs:
      opendatahub-io:
        - opendatahub-operator
        - odh-dashboard
        - notebooks
      kserve:
        - kserve
    exclude_globs:
      - "notebooks"
    ocp_versions:
      min: "4.14"
      max: "4.16"
    repo_overrides:
      kserve:
        tier: ml
        version_pin: "v0.12.1"
  rhoai:
    name: "RHOAI"
    orgs:
      red-hat-data-services:
        - opendatahub-operator
`
	path := filepath.Join(dir, "scan-config.yaml")
	os.WriteFile(path, []byte(content), 0o644)

	specs, platCfg, err := LoadPlatformConfig(path, "odh")
	if err != nil {
		t.Fatal(err)
	}

	// "notebooks" should be excluded by glob
	for _, s := range specs {
		if s.Repo == "notebooks" {
			t.Error("notebooks should be excluded by glob")
		}
	}

	// Should have 3 repos (4 minus 1 excluded)
	if len(specs) != 3 {
		t.Errorf("expected 3 specs after exclusion, got %d", len(specs))
	}

	// kserve should have overrides
	for _, s := range specs {
		if s.Repo == "kserve" {
			if s.Overrides == nil {
				t.Error("expected overrides for kserve")
			} else {
				if s.Overrides.Tier != "ml" {
					t.Errorf("expected tier 'ml', got %q", s.Overrides.Tier)
				}
				if s.Overrides.VersionPin != "v0.12.1" {
					t.Errorf("expected version_pin 'v0.12.1', got %q", s.Overrides.VersionPin)
				}
			}
		}
	}

	// OCP versions
	if platCfg.OCPVersions == nil {
		t.Fatal("expected OCP versions config")
	}
	if platCfg.OCPVersions.Min != "4.14" {
		t.Errorf("expected min 4.14, got %q", platCfg.OCPVersions.Min)
	}
}

func TestLoadPlatformConfigMissingPlatform(t *testing.T) {
	dir := t.TempDir()
	content := `
platforms:
  odh:
    name: "ODH"
    orgs:
      test:
        - repo1
`
	path := filepath.Join(dir, "scan-config.yaml")
	os.WriteFile(path, []byte(content), 0o644)

	_, _, err := LoadPlatformConfig(path, "nonexistent")
	if err == nil {
		t.Error("expected error for missing platform")
	}
}

func TestLoadPlatformConfigExcludeOverride(t *testing.T) {
	dir := t.TempDir()
	content := `
platforms:
  test:
    name: "Test"
    orgs:
      org:
        - keep
        - drop
    repo_overrides:
      drop:
        exclude: true
`
	path := filepath.Join(dir, "scan-config.yaml")
	os.WriteFile(path, []byte(content), 0o644)

	specs, _, err := LoadPlatformConfig(path, "test")
	if err != nil {
		t.Fatal(err)
	}
	if len(specs) != 1 {
		t.Errorf("expected 1 spec (drop excluded), got %d", len(specs))
	}
	if specs[0].Repo != "keep" {
		t.Errorf("expected 'keep', got %q", specs[0].Repo)
	}
}

func TestListPlatforms(t *testing.T) {
	dir := t.TempDir()
	content := `
platforms:
  rhoai:
    name: "RHOAI"
    orgs: {}
  odh:
    name: "ODH"
    orgs: {}
`
	path := filepath.Join(dir, "scan-config.yaml")
	os.WriteFile(path, []byte(content), 0o644)

	names, err := ListPlatforms(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 platforms, got %d", len(names))
	}
	// Sorted
	if names[0] != "odh" {
		t.Errorf("expected 'odh' first (sorted), got %q", names[0])
	}
}

func TestSplitOrgRepo(t *testing.T) {
	tests := []struct {
		input string
		org   string
		repo  string
	}{
		{"opendatahub-io/dashboard", "opendatahub-io", "dashboard"},
		{"kserve/kserve", "kserve", "kserve"},
		{"standalone-repo", "", "standalone-repo"},
	}
	for _, tt := range tests {
		got := splitOrgRepo(tt.input)
		if got.Org != tt.org || got.Repo != tt.repo {
			t.Errorf("splitOrgRepo(%q) = {%q, %q}, want {%q, %q}", tt.input, got.Org, got.Repo, tt.org, tt.repo)
		}
	}
}

func TestLoadScanConfigMissingFile(t *testing.T) {
	_, err := LoadScanConfig("/nonexistent/config.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
