package extractor

import (
	"os"
	"path/filepath"
	"strings"

	"sigs.k8s.io/yaml"
)

// KustomizeOverlayRef captures the structure of a kustomization.yaml file,
// exposing what resources, patches, and generators an overlay references.
// This gives LLMs explicit cross-reference data about overlay composition.
type KustomizeOverlayRef struct {
	Overlay          string   `json:"overlay"`
	Resources        []string `json:"resources,omitempty"`
	Patches          []string `json:"patches,omitempty"`
	ConfigMapGens    []string `json:"configmap_generators,omitempty"`
	SecretGens       []string `json:"secret_generators,omitempty"`
	Components       []string `json:"components,omitempty"`
	Images           []string `json:"images,omitempty"`
	Namespace        string   `json:"namespace,omitempty"`
	NamePrefix       string   `json:"name_prefix,omitempty"`
	NameSuffix       string   `json:"name_suffix,omitempty"`
	CommonLabels     []string `json:"common_labels,omitempty"`
	Replacements     int      `json:"replacements,omitempty"`
	Source           string   `json:"source"`
}

// extractKustomizeOverlayRefs scans all kustomization.yaml files in the repo
// and extracts their structure (resources, patches, generators, transforms).
func extractKustomizeOverlayRefs(repoPath string) []KustomizeOverlayRef {
	var refs []KustomizeOverlayRef

	// Walk config/ directory for kustomization files
	configDir := filepath.Join(repoPath, "config")
	if _, err := os.Stat(configDir); err != nil {
		return nil
	}

	_ = filepath.Walk(configDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		name := info.Name()
		if name != "kustomization.yaml" && name != "kustomization.yml" && name != "Kustomization" {
			return nil
		}

		ref := parseKustomizationFile(repoPath, path)
		if ref != nil {
			refs = append(refs, *ref)
		}
		return nil
	})

	return refs
}

// parseKustomizationFile reads a kustomization.yaml and extracts cross-reference data.
func parseKustomizationFile(repoPath, path string) *KustomizeOverlayRef {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var kustomization struct {
		Resources    []interface{} `json:"resources" yaml:"resources"`
		Bases        []string      `json:"bases" yaml:"bases"`
		Patches      []interface{} `json:"patches" yaml:"patches"`
		PatchesSM    []interface{} `json:"patchesStrategicMerge" yaml:"patchesStrategicMerge"`
		PatchesJSON  []interface{} `json:"patchesJson6902" yaml:"patchesJson6902"`
		Components   []string      `json:"components" yaml:"components"`
		Namespace    string        `json:"namespace" yaml:"namespace"`
		NamePrefix   string        `json:"namePrefix" yaml:"namePrefix"`
		NameSuffix   string        `json:"nameSuffix" yaml:"nameSuffix"`
		CommonLabels map[string]string `json:"commonLabels" yaml:"commonLabels"`
		Images       []struct {
			Name    string `json:"name" yaml:"name"`
			NewName string `json:"newName" yaml:"newName"`
			NewTag  string `json:"newTag" yaml:"newTag"`
			Digest  string `json:"digest" yaml:"digest"`
		} `json:"images" yaml:"images"`
		ConfigMapGenerator []struct {
			Name string `json:"name" yaml:"name"`
		} `json:"configMapGenerator" yaml:"configMapGenerator"`
		SecretGenerator []struct {
			Name string `json:"name" yaml:"name"`
		} `json:"secretGenerator" yaml:"secretGenerator"`
		Replacements []interface{} `json:"replacements" yaml:"replacements"`
	}

	if err := yaml.Unmarshal(data, &kustomization); err != nil {
		return nil
	}

	relPath, _ := filepath.Rel(repoPath, path)
	overlayName := filepath.Dir(relPath)

	ref := &KustomizeOverlayRef{
		Overlay:    overlayName,
		Namespace:  kustomization.Namespace,
		NamePrefix: kustomization.NamePrefix,
		NameSuffix: kustomization.NameSuffix,
		Source:     relPath,
	}

	// Resources (can be strings or maps with path field)
	for _, r := range kustomization.Resources {
		switch v := r.(type) {
		case string:
			ref.Resources = append(ref.Resources, v)
		case map[string]interface{}:
			if p, ok := v["path"].(string); ok {
				ref.Resources = append(ref.Resources, p)
			}
		}
	}
	// Legacy bases field
	ref.Resources = append(ref.Resources, kustomization.Bases...)

	// Patches: extract target file paths from various patch formats
	for _, p := range kustomization.Patches {
		switch v := p.(type) {
		case string:
			ref.Patches = append(ref.Patches, v)
		case map[string]interface{}:
			if path, ok := v["path"].(string); ok {
				ref.Patches = append(ref.Patches, path)
			} else if target, ok := v["target"].(map[string]interface{}); ok {
				// Inline patch with target selector
				kind, _ := target["kind"].(string)
				name, _ := target["name"].(string)
				if kind != "" || name != "" {
					ref.Patches = append(ref.Patches, strings.TrimSpace(kind+" "+name))
				}
			}
		}
	}
	for _, p := range kustomization.PatchesSM {
		if s, ok := p.(string); ok {
			ref.Patches = append(ref.Patches, s)
		}
	}
	for _, p := range kustomization.PatchesJSON {
		if m, ok := p.(map[string]interface{}); ok {
			if path, ok := m["path"].(string); ok {
				ref.Patches = append(ref.Patches, path)
			}
		}
	}

	// Components
	ref.Components = kustomization.Components

	// Images
	for _, img := range kustomization.Images {
		target := img.Name
		if img.NewName != "" {
			target += " -> " + img.NewName
		}
		if img.NewTag != "" {
			target += ":" + img.NewTag
		} else if img.Digest != "" {
			target += "@" + img.Digest
		}
		ref.Images = append(ref.Images, target)
	}

	// ConfigMap generators
	for _, g := range kustomization.ConfigMapGenerator {
		if g.Name != "" {
			ref.ConfigMapGens = append(ref.ConfigMapGens, g.Name)
		}
	}

	// Secret generators
	for _, g := range kustomization.SecretGenerator {
		if g.Name != "" {
			ref.SecretGens = append(ref.SecretGens, g.Name)
		}
	}

	// Common labels
	for k, v := range kustomization.CommonLabels {
		ref.CommonLabels = append(ref.CommonLabels, k+"="+v)
	}

	// Replacements count
	ref.Replacements = len(kustomization.Replacements)

	// Skip empty refs (no meaningful data)
	if len(ref.Resources) == 0 && len(ref.Patches) == 0 &&
		len(ref.ConfigMapGens) == 0 && len(ref.SecretGens) == 0 &&
		len(ref.Components) == 0 && len(ref.Images) == 0 {
		return nil
	}

	return ref
}
