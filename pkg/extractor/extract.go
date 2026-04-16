package extractor

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// AnalyzerVersion is the version of the analyzer, set by the CLI entry point.
// Defaults to "dev" if not set.
var AnalyzerVersion = "dev"

// ExtractAll runs all extractors on the given repository path and returns the
// combined ComponentArchitecture. Pass nil for opts to use defaults.
func ExtractAll(repoPath string, opts *ExtractOptions) (*ComponentArchitecture, error) {
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, fmt.Errorf("invalid repo path: %w", err)
	}
	info, err := os.Stat(absPath)
	if err != nil || !info.IsDir() {
		return nil, fmt.Errorf("repository path does not exist or is not a directory: %s", absPath)
	}

	if opts == nil {
		opts = &ExtractOptions{}
	}

	componentName := filepath.Base(absPath)

	org := opts.Org
	if org == "" {
		org = detectOrg(absPath)
	}

	modulePrefixes := opts.ModulePrefixes
	if len(modulePrefixes) == 0 {
		modulePrefixes = DefaultModulePrefixes()
	}

	arch := &ComponentArchitecture{
		Component:       componentName,
		Repo:            fmt.Sprintf("%s/%s", org, componentName),
		ExtractedAt:     time.Now().UTC().Format(time.RFC3339),
		AnalyzerVersion: AnalyzerVersion,
		CRDs:            extractCRDs(absPath),
		RBAC:            extractRBAC(absPath),
		Services:        extractServices(absPath),
		Deployments:     extractDeployments(absPath),
		NetworkPolicies: extractNetworkPolicies(absPath),
		ControllerWatch: extractControllerWatches(absPath),
		Dependencies:    extractDependencies(absPath, modulePrefixes),
		Secrets:         extractSecrets(absPath),
		Dockerfiles:     extractDockerfiles(absPath),
		Helm:            extractHelm(absPath),
		Webhooks:        extractWebhooks(absPath),
		ConfigMaps:      extractConfigMaps(absPath),
		HTTPEndpoints:   extractHTTPEndpoints(absPath),
		IngressRouting:      extractIngress(absPath),
		ExternalConnections: extractExternalConnections(absPath),
	}

	// Cache analysis runs after watches and deployments are extracted
	arch.CacheConfig = extractCacheConfig(absPath, arch.ControllerWatch, arch.Deployments)

	return arch, nil
}

// detectOrg tries to determine the GitHub organization from the repo's go.mod
// module path, then from .git/config remote origin, then falls back to
// "opendatahub-io".
func detectOrg(repoPath string) string {
	// Try go.mod first
	goModPath := filepath.Join(repoPath, "go.mod")
	if f, err := os.Open(goModPath); err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "module ") {
				module := strings.TrimPrefix(line, "module ")
				module = strings.TrimSpace(module)
				// Parse github.com/org/repo format
				parts := strings.Split(module, "/")
				if len(parts) >= 2 && parts[0] == "github.com" {
					return parts[1]
				}
			}
		}
	}

	// Try .git/config remote origin
	gitConfigPath := filepath.Join(repoPath, ".git", "config")
	if f, err := os.Open(gitConfigPath); err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		inOrigin := false
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			// Track which [remote] section we're in
			if strings.HasPrefix(line, "[") {
				inOrigin = line == `[remote "origin"]`
				continue
			}
			if !inOrigin {
				continue
			}
			if !strings.HasPrefix(line, "url = ") {
				continue
			}
			url := strings.TrimPrefix(line, "url = ")
			url = strings.TrimSpace(url)
			if org := orgFromGitURL(url); org != "" {
				return org
			}
		}
	}

	return "opendatahub-io"
}

// orgFromGitURL extracts the GitHub organization from a git remote URL.
// Supports HTTPS (https://github.com/org/repo.git) and SSH (git@github.com:org/repo.git).
func orgFromGitURL(url string) string {
	if !strings.Contains(url, "github.com") {
		return ""
	}
	url = strings.TrimSuffix(url, ".git")

	// SSH format: git@github.com:org/repo
	if strings.Contains(url, ":") && !strings.Contains(url, "://") {
		colonParts := strings.SplitN(url, ":", 2)
		if len(colonParts) == 2 {
			orgRepo := strings.SplitN(colonParts[1], "/", 2)
			if len(orgRepo) >= 1 && orgRepo[0] != "" {
				return orgRepo[0]
			}
		}
		return ""
	}

	// HTTPS format: https://github.com/org/repo
	parts := strings.Split(url, "/")
	// Find "github.com" in parts and return the next segment
	for i, part := range parts {
		if part == "github.com" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}
