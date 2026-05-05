package extractor

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// KustomizeComponent represents a component discovered from kustomize overlay
// support files (e.g. *_support.go) in an operator repository.
type KustomizeComponent struct {
	Name           string            `json:"name"`
	SupportFile    string            `json:"support_file"`
	OverlayPaths   []string          `json:"overlay_paths,omitempty"`
	ImageParams    []ImageParam      `json:"image_params,omitempty"`
	ManagedCRDs    []string          `json:"managed_crds,omitempty"`
	FeatureFlags   []string          `json:"feature_flags,omitempty"`
}

// ImageParam maps a params.env placeholder to a RELATED_IMAGE environment
// variable. The chain is: env var name -> params.env key -> container image.
type ImageParam struct {
	EnvVar      string `json:"env_var"`
	ParamsKey   string `json:"params_key"`
	DefaultImage string `json:"default_image,omitempty"`
}

// ParamsEnv holds parsed params.env key-value pairs.
type ParamsEnv struct {
	Source string            `json:"source"`
	Params map[string]string `json:"params"`
}

// Regular expressions for parsing *_support.go and *_component.go files.
// These target operator conventions common in kubebuilder/operator-sdk repos
// (GetComponentName methods, RELATED_IMAGE params, overlay paths).
// Regex is used instead of go/ast because the patterns are simple string
// literal extractions from well-structured operator scaffold files. Limitations:
//   - componentReturnRe only checks the first 5 lines after GetComponentName()
//   - featureFlagRe may match in comments/strings (mitigated by false positive filter)
//   - Multi-line or computed return values are not captured
var (
	componentNameRe   = regexp.MustCompile(`func\s+\([^)]+\)\s+GetComponentName\(\)\s+string\s*\{`)
	componentReturnRe = regexp.MustCompile(`return\s+"([^"]+)"`)
	imageParamRe      = regexp.MustCompile(`"([^"]+)":\s*"(RELATED_IMAGE[^"]*)"`)
	overlayPathRe     = regexp.MustCompile(`"([^"]*(?:overlay|kustomize|manifests)[^"]*)"`)
	managedResourceRe = regexp.MustCompile(`GroupVersionKind\{[^}]*Kind:\s*"([^"]+)"`)
	featureFlagRe     = regexp.MustCompile(`features?\.(\w+)`)
)

// extractKustomizeComponents scans for *_support.go files in the repo and
// extracts component metadata from each. Returns nil if no support files found.
func extractKustomizeComponents(repoPath string) []KustomizeComponent {
	supportFiles := findSupportFiles(repoPath)
	if len(supportFiles) == 0 {
		return nil
	}

	var components []KustomizeComponent
	for _, sf := range supportFiles {
		comp := parseComponentSupportFile(repoPath, sf)
		if comp.Name != "" {
			components = append(components, comp)
		}
	}

	sort.Slice(components, func(i, j int) bool {
		return components[i].Name < components[j].Name
	})
	return components
}

// findSupportFiles searches the repo for Go files matching the pattern
// *_support.go or *_component.go, which typically define component metadata.
func findSupportFiles(repoPath string) []string {
	var files []string
	_ = filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if isExcludedDir(filepath.Base(path), nil) {
				return filepath.SkipDir
			}
			return nil
		}
		name := info.Name()
		if strings.HasSuffix(name, "_support.go") || strings.HasSuffix(name, "_component.go") {
			files = append(files, path)
		}
		return nil
	})
	return files
}

// parseComponentSupportFile extracts component metadata from a single support file.
func parseComponentSupportFile(repoPath, filePath string) KustomizeComponent {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return KustomizeComponent{}
	}
	content := string(data)
	relPath, _ := filepath.Rel(repoPath, filePath)

	comp := KustomizeComponent{
		SupportFile: relPath,
	}

	// Extract component name from GetComponentName() method
	comp.Name = extractComponentName(content)
	if comp.Name == "" {
		// Fallback: derive from filename (dashboard_support.go -> dashboard)
		base := filepath.Base(filePath)
		comp.Name = strings.TrimSuffix(strings.TrimSuffix(base, "_support.go"), "_component.go")
	}

	// Extract image parameter mappings
	comp.ImageParams = extractImageParams(content)

	// Extract overlay paths
	comp.OverlayPaths = extractOverlayPaths(content)

	// Extract managed CRDs/resources
	comp.ManagedCRDs = extractManagedResources(content)

	// Extract feature flag references
	comp.FeatureFlags = extractFeatureFlags(content)

	return comp
}

// extractComponentName looks for a GetComponentName() method and returns its value.
func extractComponentName(content string) string {
	loc := componentNameRe.FindStringIndex(content)
	if loc == nil {
		return ""
	}
	// Scan forward from the opening brace to find the return statement
	rest := content[loc[1]:]
	lines := strings.SplitN(rest, "\n", 5)
	for _, line := range lines {
		matches := componentReturnRe.FindStringSubmatch(line)
		if len(matches) >= 2 {
			return matches[1]
		}
	}
	return ""
}

// extractImageParams finds all imageParamMap entries mapping params.env keys
// to RELATED_IMAGE environment variables.
func extractImageParams(content string) []ImageParam {
	matches := imageParamRe.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	var params []ImageParam
	for _, m := range matches {
		key := m[1]
		envVar := m[2]
		if seen[key] {
			continue
		}
		seen[key] = true
		params = append(params, ImageParam{
			EnvVar:    envVar,
			ParamsKey: key,
		})
	}
	return params
}

// extractOverlayPaths finds kustomize overlay directory references.
func extractOverlayPaths(content string) []string {
	matches := overlayPathRe.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	var paths []string
	for _, m := range matches {
		p := m[1]
		if seen[p] {
			continue
		}
		seen[p] = true
		paths = append(paths, p)
	}
	sort.Strings(paths)
	return paths
}

// extractManagedResources finds GVK references indicating managed resources.
func extractManagedResources(content string) []string {
	matches := managedResourceRe.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	var resources []string
	for _, m := range matches {
		r := m[1]
		if seen[r] {
			continue
		}
		seen[r] = true
		resources = append(resources, r)
	}
	sort.Strings(resources)
	return resources
}

// extractFeatureFlags finds feature gate references in the file.
func extractFeatureFlags(content string) []string {
	matches := featureFlagRe.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	var flags []string
	for _, m := range matches {
		f := m[1]
		// Skip common false positives
		if f == "Feature" || f == "FeatureGate" || f == "FeatureTracker" || f == "String" {
			continue
		}
		if seen[f] {
			continue
		}
		seen[f] = true
		flags = append(flags, f)
	}
	sort.Strings(flags)
	return flags
}

// parseParamsEnv reads a params.env file and returns key-value pairs.
// params.env uses KEY=VALUE format, one per line. Lines starting with # are comments.
func parseParamsEnv(path string) (*ParamsEnv, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := &ParamsEnv{
		Source: path,
		Params: make(map[string]string),
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		result.Params[parts[0]] = parts[1]
	}
	return result, scanner.Err()
}

// resolveImageDefaults populates DefaultImage on each ImageParam by looking up
// the ParamsKey in the given params.env data.
func resolveImageDefaults(components []KustomizeComponent, paramsEnv *ParamsEnv) {
	if paramsEnv == nil {
		return
	}
	for i := range components {
		for j := range components[i].ImageParams {
			key := components[i].ImageParams[j].ParamsKey
			if val, ok := paramsEnv.Params[key]; ok {
				components[i].ImageParams[j].DefaultImage = val
			}
		}
	}
}

// findParamsEnv searches for params.env files in common kustomize locations.
func findParamsEnv(repoPath string) []string {
	candidates := []string{
		"config/params.env",
		"config/manager/params.env",
		"manifests/params.env",
	}

	var found []string
	for _, c := range candidates {
		p := filepath.Join(repoPath, c)
		if _, err := os.Stat(p); err == nil {
			found = append(found, p)
		}
	}

	// Also search more broadly
	_ = filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if isExcludedDir(filepath.Base(path), nil) {
				return filepath.SkipDir
			}
			return nil
		}
		if info.Name() == "params.env" {
			for _, existing := range found {
				if existing == path {
					return nil
				}
			}
			found = append(found, path)
		}
		return nil
	})

	return found
}

// DiscoverPlatformComponents discovers all components from an operator repo
// by scanning support files and resolving params.env defaults.
// This is the main entry point for kustomize component discovery.
func DiscoverPlatformComponents(repoPath string) (*PlatformDiscovery, error) {
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, fmt.Errorf("invalid repo path: %w", err)
	}

	components := extractKustomizeComponents(absPath)
	if len(components) == 0 {
		return &PlatformDiscovery{
			Source:     absPath,
			Components: nil,
		}, nil
	}

	// Resolve image defaults from params.env files
	paramsFiles := findParamsEnv(absPath)
	var allParams *ParamsEnv
	for _, pf := range paramsFiles {
		pe, err := parseParamsEnv(pf)
		if err != nil {
			continue
		}
		if allParams == nil {
			allParams = pe
		} else {
			// Merge: later files override earlier ones
			for k, v := range pe.Params {
				allParams.Params[k] = v
			}
		}
	}
	resolveImageDefaults(components, allParams)

	return &PlatformDiscovery{
		Source:     absPath,
		Components: components,
		ParamsEnv:  paramsFiles,
	}, nil
}

// PlatformDiscovery holds the results of kustomize component discovery.
type PlatformDiscovery struct {
	Source     string               `json:"source"`
	Components []KustomizeComponent `json:"components"`
	ParamsEnv  []string             `json:"params_env_files,omitempty"`
}

// ComponentNames returns the names of all discovered components.
func (pd *PlatformDiscovery) ComponentNames() []string {
	names := make([]string, len(pd.Components))
	for i, c := range pd.Components {
		names[i] = c.Name
	}
	return names
}

// TotalImageParams returns the total count of image parameter mappings.
func (pd *PlatformDiscovery) TotalImageParams() int {
	total := 0
	for _, c := range pd.Components {
		total += len(c.ImageParams)
	}
	return total
}

// FormatSummary returns a human-readable summary of the discovery results.
func (pd *PlatformDiscovery) FormatSummary() string {
	if len(pd.Components) == 0 {
		return "No kustomize components discovered."
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Discovered %d component(s) with %d image parameter(s):\n",
		len(pd.Components), pd.TotalImageParams())
	for _, c := range pd.Components {
		fmt.Fprintf(&b, "  %-30s images: %d, overlays: %d, crds: %d\n",
			c.Name, len(c.ImageParams), len(c.OverlayPaths), len(c.ManagedCRDs))
	}
	return b.String()
}
