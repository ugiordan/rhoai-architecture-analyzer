// Package config handles scan configuration parsing.
package config

import (
	"fmt"
	"os"
	"sort"

	"gopkg.in/yaml.v3"
)

// RepoSpec identifies a repository to analyze.
type RepoSpec struct {
	Org  string
	Repo string
}

// FullName returns the "org/repo" string.
func (r RepoSpec) FullName() string {
	return r.Org + "/" + r.Repo
}

// ScanConfig represents the scan-config.yaml file.
type ScanConfig struct {
	// Orgs maps organization names to their repo lists (new format).
	Orgs map[string]OrgConfig `yaml:"orgs,omitempty"`
	// Repos is the legacy flat list format ("org/repo" strings).
	Repos []string `yaml:"repos,omitempty"`
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

func splitOrgRepo(fullRepo string) RepoSpec {
	for i, c := range fullRepo {
		if c == '/' {
			return RepoSpec{Org: fullRepo[:i], Repo: fullRepo[i+1:]}
		}
	}
	return RepoSpec{Org: "opendatahub-io", Repo: fullRepo}
}
