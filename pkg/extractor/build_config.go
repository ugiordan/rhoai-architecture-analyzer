package extractor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// BuildConfig represents metadata extracted from a RHOAI-Build-Config style
// repository: OCP version ranges, CPU architectures, OLM features, and
// image-to-source-repo mappings.
type BuildConfig struct {
	OCPVersions    OCPVersionRange    `json:"ocp_versions"`
	Architectures  []string           `json:"architectures"`
	OLMFeatures    []string           `json:"olm_features,omitempty"`
	ImageMappings  []ImageMapping     `json:"image_mappings,omitempty"`
	OperatorConfig *OperatorConfig    `json:"operator_config,omitempty"`
	Source         string             `json:"source"`
}

// OCPVersionRange captures minimum and maximum supported OCP versions.
type OCPVersionRange struct {
	Min string `json:"min,omitempty"`
	Max string `json:"max,omitempty"`
}

// ImageMapping maps a container image to its source repository and commit.
type ImageMapping struct {
	Image      string `json:"image"`
	SourceRepo string `json:"source_repo,omitempty"`
	Component  string `json:"component,omitempty"`
	Tag        string `json:"tag,omitempty"`
}

// OperatorConfig holds OLM operator-specific build configuration.
type OperatorConfig struct {
	PackageName       string   `json:"package_name,omitempty"`
	DefaultChannel    string   `json:"default_channel,omitempty"`
	Channels          []string `json:"channels,omitempty"`
	SkipRange         string   `json:"skip_range,omitempty"`
	InstallModes      []string `json:"install_modes,omitempty"`
	MinKubeVersion    string   `json:"min_kube_version,omitempty"`
	Replaces          string   `json:"replaces,omitempty"`
}

// ParseBuildConfig reads build configuration from a directory that follows the
// RHOAI-Build-Config layout. It looks for:
//   - config.yaml / build.yaml / build-config.yaml (OCP versions, arches)
//   - Makefile / Dockerfile (architecture hints)
//   - CSV template files (OLM metadata)
//   - image-references or image mapping files
func ParseBuildConfig(dir string) (*BuildConfig, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("resolving path: %w", err)
	}
	info, err := os.Stat(absDir)
	if err != nil || !info.IsDir() {
		return nil, fmt.Errorf("build config directory does not exist: %s", absDir)
	}

	bc := &BuildConfig{Source: absDir}

	// Parse YAML config files
	for _, name := range []string{"config.yaml", "build.yaml", "build-config.yaml", "config.json"} {
		path := filepath.Join(absDir, name)
		if _, err := os.Stat(path); err == nil {
			parseBuildConfigFile(path, bc)
			break
		}
	}

	// Parse Makefile for architecture hints
	parseMakefileArches(absDir, bc)

	// Parse CSV template for OLM metadata
	parseCSVTemplate(absDir, bc)

	// Parse image references
	parseImageReferences(absDir, bc)

	// Deduplicate
	bc.Architectures = dedupStrings(bc.Architectures)
	bc.OLMFeatures = dedupStrings(bc.OLMFeatures)

	sort.Strings(bc.Architectures)
	sort.Strings(bc.OLMFeatures)

	return bc, nil
}

var (
	ocpVersionRe = regexp.MustCompile(`(?i)(?:ocp|openshift)[_\-\s]*(?:version|ver)?[_\-\s]*(?:min|minimum)[:\s=]+["']?([0-9]+\.[0-9]+)["']?`)
	ocpMaxRe     = regexp.MustCompile(`(?i)(?:ocp|openshift)[_\-\s]*(?:version|ver)?[_\-\s]*(?:max|maximum)[:\s=]+["']?([0-9]+\.[0-9]+)["']?`)
	archRe       = regexp.MustCompile(`(?:amd64|x86_64|arm64|aarch64|ppc64le|s390x)`)
	kubeVerRe    = regexp.MustCompile(`(?i)min[_\-]?kube[_\-]?version[:\s=]+["']?([0-9]+\.[0-9]+(?:\.[0-9]+)?)["']?`)
	channelRe    = regexp.MustCompile(`(?i)(?:default[_\-]?)?channel[:\s=]+["']?([a-zA-Z0-9._-]+)["']?`)
	packageRe    = regexp.MustCompile(`(?i)package[_\-]?name[:\s=]+["']?([a-zA-Z0-9._-]+)["']?`)
	skipRangeRe  = regexp.MustCompile(`(?i)skip[_\-]?range[:\s=]+["']?([^"'\s]+)["']?`)
)

func parseBuildConfigFile(path string, bc *BuildConfig) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	content := string(data)

	// Try JSON first
	if strings.HasSuffix(path, ".json") {
		var raw map[string]interface{}
		if json.Unmarshal(data, &raw) == nil {
			parseBuildConfigJSON(raw, bc)
			return
		}
	}

	// Regex-based extraction from YAML or text
	if m := ocpVersionRe.FindStringSubmatch(content); len(m) > 1 {
		bc.OCPVersions.Min = m[1]
	}
	if m := ocpMaxRe.FindStringSubmatch(content); len(m) > 1 {
		bc.OCPVersions.Max = m[1]
	}
	if m := kubeVerRe.FindStringSubmatch(content); len(m) > 1 {
		if bc.OperatorConfig == nil {
			bc.OperatorConfig = &OperatorConfig{}
		}
		bc.OperatorConfig.MinKubeVersion = m[1]
	}

	arches := archRe.FindAllString(content, -1)
	for _, a := range arches {
		bc.Architectures = append(bc.Architectures, normalizeArch(a))
	}
}

func parseBuildConfigJSON(raw map[string]interface{}, bc *BuildConfig) {
	if ocp, ok := raw["ocp_versions"].(map[string]interface{}); ok {
		if v, ok := ocp["min"].(string); ok {
			bc.OCPVersions.Min = v
		}
		if v, ok := ocp["max"].(string); ok {
			bc.OCPVersions.Max = v
		}
	}
	if arches, ok := raw["architectures"].([]interface{}); ok {
		for _, a := range arches {
			if s, ok := a.(string); ok {
				bc.Architectures = append(bc.Architectures, normalizeArch(s))
			}
		}
	}
	if features, ok := raw["olm_features"].([]interface{}); ok {
		for _, f := range features {
			if s, ok := f.(string); ok {
				bc.OLMFeatures = append(bc.OLMFeatures, s)
			}
		}
	}
}

func parseMakefileArches(dir string, bc *BuildConfig) {
	for _, name := range []string{"Makefile", "makefile", "GNUmakefile"} {
		path := filepath.Join(dir, name)
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			// Look for PLATFORMS or GOARCH variable assignments
			if strings.Contains(line, "PLATFORM") || strings.Contains(line, "GOARCH") || strings.Contains(line, "ARCH") {
				arches := archRe.FindAllString(line, -1)
				for _, a := range arches {
					bc.Architectures = append(bc.Architectures, normalizeArch(a))
				}
			}
		}
		f.Close()
		break
	}
}

func parseCSVTemplate(dir string, bc *BuildConfig) {
	// Walk directory tree for ClusterServiceVersion YAML files.
	// Note: filepath.Glob does not support "**" recursive patterns,
	// so we use filepath.WalkDir instead.
	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		name := d.Name()
		if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
			return nil
		}
		isCSVCandidate := strings.Contains(name, "csv") ||
			strings.Contains(name, "clusterserviceversion") ||
			filepath.Base(filepath.Dir(path)) == "manifests"
		if !isCSVCandidate {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		content := string(data)
		if !strings.Contains(content, "ClusterServiceVersion") && !strings.Contains(content, "olm.skipRange") {
			return nil
		}

		if bc.OperatorConfig == nil {
			bc.OperatorConfig = &OperatorConfig{}
		}

		if m := kubeVerRe.FindStringSubmatch(content); len(m) > 1 {
			bc.OperatorConfig.MinKubeVersion = m[1]
		}
		if m := channelRe.FindStringSubmatch(content); len(m) > 1 {
			bc.OperatorConfig.DefaultChannel = m[1]
		}
		if m := packageRe.FindStringSubmatch(content); len(m) > 1 {
			bc.OperatorConfig.PackageName = m[1]
		}
		if m := skipRangeRe.FindStringSubmatch(content); len(m) > 1 {
			bc.OperatorConfig.SkipRange = m[1]
		}

		// Extract install modes
		if strings.Contains(content, "OwnNamespace") {
			bc.OLMFeatures = append(bc.OLMFeatures, "OwnNamespace")
		}
		if strings.Contains(content, "SingleNamespace") {
			bc.OLMFeatures = append(bc.OLMFeatures, "SingleNamespace")
		}
		if strings.Contains(content, "AllNamespaces") {
			bc.OLMFeatures = append(bc.OLMFeatures, "AllNamespaces")
		}
		if strings.Contains(content, "MultiNamespace") {
			bc.OLMFeatures = append(bc.OLMFeatures, "MultiNamespace")
		}

		return filepath.SkipAll // only process first match
	})
}

func parseImageReferences(dir string, bc *BuildConfig) {
	// Check for image-references file (Konflux/OSBS format)
	for _, name := range []string{"image-references", "image_references.yaml", "image-references.yaml"} {
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		// Parse line-by-line: "name=image:tag" or YAML-ish format
		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				im := ImageMapping{
					Component: strings.TrimSpace(parts[0]),
					Image:     strings.TrimSpace(parts[1]),
				}
				if idx := strings.LastIndex(im.Image, ":"); idx > 0 {
					im.Tag = im.Image[idx+1:]
				}
				bc.ImageMappings = append(bc.ImageMappings, im)
			}
		}
		break
	}
}

func normalizeArch(arch string) string {
	switch strings.ToLower(arch) {
	case "x86_64", "amd64":
		return "amd64"
	case "aarch64", "arm64":
		return "arm64"
	default:
		return strings.ToLower(arch)
	}
}

func dedupStrings(ss []string) []string {
	if len(ss) == 0 {
		return nil
	}
	seen := make(map[string]bool)
	var out []string
	for _, s := range ss {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}
