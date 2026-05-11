// Package config handles scan configuration parsing.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// RepoSpec identifies a repository to analyze with optional overrides.
type RepoSpec struct {
	Org       string         `json:"org"`
	Repo      string         `json:"repo"`
	// Overrides are platform-level configuration for this repo.
	Overrides *RepoOverrides `json:"overrides,omitempty"`
}

// FullName returns the "org/repo" string.
func (r RepoSpec) FullName() string {
	return r.Org + "/" + r.Repo
}

// RepoOverrides holds per-repo configuration overrides set in the platform config.
type RepoOverrides struct {
	Tier       string   `yaml:"tier,omitempty" json:"tier,omitempty"`
	Type       string   `yaml:"type,omitempty" json:"type,omitempty"`
	VersionPin string   `yaml:"version_pin,omitempty" json:"version_pin,omitempty"`
	Branch     string   `yaml:"branch,omitempty" json:"branch,omitempty"`
	Exclude    bool     `yaml:"exclude,omitempty" json:"exclude,omitempty"`
	Aliases    []string `yaml:"aliases,omitempty" json:"aliases,omitempty"`
}

// ScanConfig represents the scan-config.yaml file.
type ScanConfig struct {
	// Orgs maps organization names to their repo lists (new format).
	Orgs map[string]OrgConfig `yaml:"orgs,omitempty"`
	// Repos is the legacy flat list format ("org/repo" strings).
	Repos []string `yaml:"repos,omitempty"`
}

// PlatformScanConfig represents the top-level scan-config.yaml with platform sections.
type PlatformScanConfig struct {
	Platforms map[string]PlatformConfig `yaml:"platforms"`
}

// PlatformConfig holds config for a single platform (e.g., odh, rhoai).
type PlatformConfig struct {
	Name           string                     `yaml:"name"`
	Description    string                     `yaml:"description,omitempty"`
	Orgs           map[string]interface{}      `yaml:"orgs"`
	ExcludeGlobs   []string                   `yaml:"exclude_globs,omitempty"`
	OCPVersions    *OCPVersionConfig          `yaml:"ocp_versions,omitempty"`
	RepoOverrides  map[string]*RepoOverrides  `yaml:"repo_overrides,omitempty"`
}

// OCPVersionConfig holds OCP version constraints for a platform.
type OCPVersionConfig struct {
	Min string `yaml:"min,omitempty" json:"min,omitempty"`
	Max string `yaml:"max,omitempty" json:"max,omitempty"`
}

// OrgConfig holds the repos for a single organization.
type OrgConfig struct {
	Repos []string `yaml:"repos"`
}

// LoadScanConfig reads and parses scan-config.yaml, supporting both the new
// org-grouped format and the legacy flat list format.
func LoadScanConfig(path string) ([]RepoSpec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading scan config: %w", err)
	}

	var cfg ScanConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing scan config: %w", err)
	}

	var specs []RepoSpec

	// New format: org-grouped
	if len(cfg.Orgs) > 0 {
		orgs := make([]string, 0, len(cfg.Orgs))
		for org := range cfg.Orgs {
			orgs = append(orgs, org)
		}
		sort.Strings(orgs)
		for _, org := range orgs {
			orgCfg := cfg.Orgs[org]
			for _, repo := range orgCfg.Repos {
				specs = append(specs, RepoSpec{Org: org, Repo: repo})
			}
		}
	}

	// Legacy format: flat list of "org/repo"
	for _, fullRepo := range cfg.Repos {
		parts := splitOrgRepo(fullRepo)
		specs = append(specs, parts)
	}

	return specs, nil
}

// LoadPlatformConfig reads scan-config.yaml and returns repos for a specific
// platform with overrides applied.
func LoadPlatformConfig(path, platform string) ([]RepoSpec, *PlatformConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("reading scan config: %w", err)
	}

	var psc PlatformScanConfig
	if err := yaml.Unmarshal(data, &psc); err != nil {
		return nil, nil, fmt.Errorf("parsing platform scan config: %w", err)
	}

	platCfg, ok := psc.Platforms[platform]
	if !ok {
		return nil, nil, fmt.Errorf("platform %q not found in config (available: %s)",
			platform, joinKeys(psc.Platforms))
	}

	specs, err := expandPlatformOrgs(platCfg)
	if err != nil {
		return nil, nil, err
	}

	// Apply exclusion globs
	if len(platCfg.ExcludeGlobs) > 0 {
		specs = applyExcludeGlobs(specs, platCfg.ExcludeGlobs)
	}

	// Apply repo overrides (prefer org/repo key, fall back to bare repo name)
	if len(platCfg.RepoOverrides) > 0 {
		for i := range specs {
			fullKey := specs[i].FullName()
			if ov, ok := platCfg.RepoOverrides[fullKey]; ok {
				specs[i].Overrides = ov
			} else if ov, ok := platCfg.RepoOverrides[specs[i].Repo]; ok {
				specs[i].Overrides = ov
			}
		}
		// Remove excluded repos
		var filtered []RepoSpec
		for _, s := range specs {
			if s.Overrides != nil && s.Overrides.Exclude {
				continue
			}
			filtered = append(filtered, s)
		}
		specs = filtered
	}

	return specs, &platCfg, nil
}

// ListPlatforms returns all platform names defined in a scan config file.
func ListPlatforms(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading scan config: %w", err)
	}

	var psc PlatformScanConfig
	if err := yaml.Unmarshal(data, &psc); err != nil {
		return nil, fmt.Errorf("parsing scan config: %w", err)
	}

	names := make([]string, 0, len(psc.Platforms))
	for name := range psc.Platforms {
		names = append(names, name)
	}
	sort.Strings(names)
	return names, nil
}

func expandPlatformOrgs(platCfg PlatformConfig) ([]RepoSpec, error) {
	var specs []RepoSpec

	orgs := make([]string, 0, len(platCfg.Orgs))
	for org := range platCfg.Orgs {
		orgs = append(orgs, org)
	}
	sort.Strings(orgs)

	for _, org := range orgs {
		raw := platCfg.Orgs[org]
		switch v := raw.(type) {
		case []interface{}:
			for _, item := range v {
				if s, ok := item.(string); ok {
					specs = append(specs, RepoSpec{Org: org, Repo: s})
				}
			}
		default:
			// Try marshaling back to YAML and parsing as string list
			yamlBytes, _ := yaml.Marshal(raw)
			var repos []string
			if yaml.Unmarshal(yamlBytes, &repos) == nil {
				for _, repo := range repos {
					specs = append(specs, RepoSpec{Org: org, Repo: repo})
				}
			}
		}
	}

	return specs, nil
}

func applyExcludeGlobs(specs []RepoSpec, globs []string) []RepoSpec {
	var out []RepoSpec
	for _, s := range specs {
		excluded := false
		for _, pattern := range globs {
			matched, _ := filepath.Match(pattern, s.Repo)
			if matched {
				excluded = true
				break
			}
			// Also try matching full name
			matched, _ = filepath.Match(pattern, s.FullName())
			if matched {
				excluded = true
				break
			}
		}
		if !excluded {
			out = append(out, s)
		}
	}
	return out
}

func joinKeys(m map[string]PlatformConfig) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, ", ")
}

func splitOrgRepo(fullRepo string) RepoSpec {
	for i, c := range fullRepo {
		if c == '/' {
			return RepoSpec{Org: fullRepo[:i], Repo: fullRepo[i+1:]}
		}
	}
	return RepoSpec{Org: "opendatahub-io", Repo: fullRepo}
}
