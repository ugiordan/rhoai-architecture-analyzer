package extractor

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// servingRuntimeRefKinds are the Kubernetes kinds we extract image refs from.
var servingRuntimeRefKinds = map[string]bool{
	"ServingRuntime":        true,
	"ClusterServingRuntime": true,
	"InferenceService":      true,
}

// extractServingRuntimeRefs scans all YAML files in the repo for
// ServingRuntime, ClusterServingRuntime, and InferenceService resources,
// extracts container image references from each, and also scans Go source
// for string literals that reference these kinds (indicating programmatic
// creation of serving resources).
func extractServingRuntimeRefs(repoPath string) []ServingRuntimeRef {
	var refs []ServingRuntimeRef

	// Phase 1: scan YAML manifests for serving runtime and inference service CRs
	refs = append(refs, extractServingRuntimeRefsFromYAML(repoPath)...)

	// Phase 2: scan Go source for programmatic Kind references
	refs = append(refs, extractServingRuntimeRefsFromGo(repoPath)...)

	// Dedup by (name, kind, image, source)
	refs = dedupServingRuntimeRefs(refs)

	sort.Slice(refs, func(i, j int) bool {
		if refs[i].Kind != refs[j].Kind {
			return refs[i].Kind < refs[j].Kind
		}
		if refs[i].Name != refs[j].Name {
			return refs[i].Name < refs[j].Name
		}
		return refs[i].ContainerImage < refs[j].ContainerImage
	})

	return refs
}

// extractServingRuntimeRefsFromYAML walks all YAML files and extracts
// container image references from serving runtime and inference service CRs.
func extractServingRuntimeRefsFromYAML(repoPath string) []ServingRuntimeRef {
	var refs []ServingRuntimeRef

	_ = filepath.WalkDir(repoPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if isExcludedDir(d.Name(), nil) {
				return fs.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}

		for _, doc := range parseYAMLSafe(path) {
			kind, _ := doc["kind"].(string)
			if !servingRuntimeRefKinds[kind] {
				continue
			}

			meta, _ := doc["metadata"].(map[string]interface{})
			name := ""
			if meta != nil {
				name, _ = meta["name"].(string)
			}
			if name == "" {
				name = "unnamed"
			}

			rel := relativePath(repoPath, path)
			images := extractImagesFromServingDoc(kind, doc)
			for _, img := range images {
				refs = append(refs, ServingRuntimeRef{
					Name:           name,
					Kind:           kind,
					ContainerImage: img,
					Source:         rel,
				})
			}
		}
		return nil
	})

	return refs
}

// extractImagesFromServingDoc extracts container image strings from a
// parsed YAML document based on the resource kind.
func extractImagesFromServingDoc(kind string, doc map[string]interface{}) []string {
	var images []string

	spec, _ := doc["spec"].(map[string]interface{})
	if spec == nil {
		return nil
	}

	switch kind {
	case "ServingRuntime", "ClusterServingRuntime":
		// spec.containers[].image
		images = append(images, extractContainerImages(spec)...)

	case "InferenceService":
		// spec.predictor.containers[].image
		// spec.predictor.model.runtime (references a ServingRuntime name, not an image)
		// spec.transformer.containers[].image
		// spec.explainer.containers[].image
		for _, section := range []string{"predictor", "transformer", "explainer"} {
			sectionMap, _ := spec[section].(map[string]interface{})
			if sectionMap == nil {
				continue
			}
			images = append(images, extractContainerImages(sectionMap)...)
		}
	}

	return images
}

// extractContainerImages pulls image strings from a map's "containers" and
// "initContainers" fields.
func extractContainerImages(parent map[string]interface{}) []string {
	var images []string
	for _, field := range []string{"containers", "initContainers"} {
		containers, _ := parent[field].([]interface{})
		for _, c := range containers {
			cm, ok := c.(map[string]interface{})
			if !ok {
				continue
			}
			img, _ := cm["image"].(string)
			if img != "" && !hasUnresolvedTemplates(img) {
				images = append(images, img)
			}
		}
	}
	return images
}

// goServingKindRE matches Go string literals like `Kind: "ServingRuntime"` or
// `kind = "InferenceService"`.
var goServingKindRE = regexp.MustCompile(
	`(?:Kind:\s*"(ServingRuntime|ClusterServingRuntime|InferenceService)"|` +
		`kind\s*[:=]\s*"(ServingRuntime|ClusterServingRuntime|InferenceService)")`,
)

// goImageLiteralRE matches image reference patterns commonly found in Go
// source (registry/org/image:tag or registry/org/image@sha256:...).
var goImageLiteralRE = regexp.MustCompile(
	`"([a-zA-Z0-9][a-zA-Z0-9._-]*/[a-zA-Z0-9._/-]+:[a-zA-Z0-9._-]+)"`,
)

// extractServingRuntimeRefsFromGo scans Go source files for string literals
// that reference ServingRuntime/InferenceService kinds near image literals,
// capturing programmatic CR creation patterns.
func extractServingRuntimeRefsFromGo(repoPath string) []ServingRuntimeRef {
	var refs []ServingRuntimeRef

	_ = filepath.WalkDir(repoPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if isExcludedDir(d.Name(), nil) {
				return fs.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		f, ferr := os.Open(path)
		if ferr != nil {
			return nil
		}
		defer f.Close()

		rel := relativePath(repoPath, path)
		scanner := bufio.NewScanner(f)
		scanner.Buffer(make([]byte, 0, 256*1024), 256*1024)
		lineNum := 0

		// Track the most recent kind reference within a window of lines
		var lastKind string
		var lastKindLine int
		const kindWindowLines = 20

		for scanner.Scan() {
			lineNum++
			line := scanner.Text()

			// Check for Kind references
			if matches := goServingKindRE.FindStringSubmatch(line); len(matches) > 0 {
				kind := matches[1]
				if kind == "" {
					kind = matches[2]
				}
				lastKind = kind
				lastKindLine = lineNum
			}

			// Check for image literals near a kind reference
			if lastKind != "" && lineNum-lastKindLine <= kindWindowLines {
				if imgMatches := goImageLiteralRE.FindAllStringSubmatch(line, -1); len(imgMatches) > 0 {
					for _, m := range imgMatches {
						img := m[1]
						// Skip obviously non-image strings (e.g., module paths)
						if strings.Contains(img, "github.com") || strings.Contains(img, "golang.org") {
							continue
						}
						refs = append(refs, ServingRuntimeRef{
							Name:           "go-source",
							Kind:           lastKind,
							ContainerImage: img,
							Source:         fmt.Sprintf("%s:%d", rel, lineNum),
						})
					}
				}
			}
		}
		return nil
	})

	return refs
}

// dedupServingRuntimeRefs removes duplicates based on (kind, name, image).
// When duplicates exist, prefer the YAML source over Go source.
func dedupServingRuntimeRefs(refs []ServingRuntimeRef) []ServingRuntimeRef {
	type key struct {
		kind, name, image string
	}
	seen := make(map[key]int) // maps to index in result
	var result []ServingRuntimeRef

	for _, ref := range refs {
		k := key{ref.Kind, ref.Name, ref.ContainerImage}
		if idx, ok := seen[k]; ok {
			// Prefer YAML sources over Go sources
			if strings.HasSuffix(result[idx].Source, ".go") && !strings.HasSuffix(ref.Source, ".go") {
				result[idx] = ref
			}
			continue
		}
		seen[k] = len(result)
		result = append(result, ref)
	}
	return result
}
