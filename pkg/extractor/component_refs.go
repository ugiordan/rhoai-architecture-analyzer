package extractor

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// ComponentRef represents a detected reference to another known component
// within the source tree. These capture cross-component integration patterns
// like provider adapters, client wrappers, or direct module references.
type ComponentRef struct {
	Target      string `json:"target"`       // Referenced component name (e.g., "vllm")
	Type        string `json:"type"`         // "provider", "adapter", "client", "plugin", "import"
	Source      string `json:"source"`       // File or directory path where reference was found
	Evidence    string `json:"evidence"`     // What triggered the detection
}

// providerDirPatterns lists directory name patterns that indicate a provider/adapter
// relationship when they contain a subdirectory matching a component name.
var providerDirPatterns = []string{
	"providers",
	"adapters",
	"backends",
	"plugins",
	"clients",
	"connectors",
	"drivers",
	"components",
}

// extractComponentRefs scans the repository for directory structures and file names
// that reference other known components. This detects cross-component integration
// patterns like llama-stack's providers/remote/inference/vllm/ adapter.
//
// knownComponents should be the list of all component names from scan-config.
// If empty, falls back to well-known RHOAI/ODH component names.
func extractComponentRefs(repoPath string, selfName string, knownComponents []string) []ComponentRef {
	if len(knownComponents) == 0 {
		knownComponents = defaultKnownComponents()
	}

	// Build lookup: both hyphenated and underscore variants
	componentSet := make(map[string]string) // normalized -> canonical
	for _, name := range knownComponents {
		if name == selfName {
			continue // skip self-references
		}
		lower := strings.ToLower(name)
		componentSet[lower] = name
		componentSet[strings.ReplaceAll(lower, "-", "_")] = name
		componentSet[strings.ReplaceAll(lower, "-", "")] = name
	}

	type refKey struct {
		target string
		typ    string
	}
	seen := make(map[refKey]bool)
	var refs []ComponentRef

	addRef := func(target, typ, source, evidence string) {
		key := refKey{target, typ}
		if seen[key] {
			return
		}
		seen[key] = true
		refs = append(refs, ComponentRef{
			Target:   target,
			Type:     typ,
			Source:   relativePath(repoPath, source),
			Evidence: evidence,
		})
	}

	// Walk the directory tree looking for component name references
	filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		name := info.Name()
		nameLower := strings.ToLower(name)

		// Skip hidden dirs, vendor, generated code
		if info.IsDir() {
			if strings.HasPrefix(name, ".") || name == "vendor" || name == "node_modules" ||
				name == "__pycache__" || name == "site-packages" || name == "venv" ||
				name == ".venv" || name == "testdata" {
				return filepath.SkipDir
			}

			// Check if this directory name matches a known component
			// AND any ancestor is a provider-like directory
			if canonical, ok := componentSet[nameLower]; ok {
				if hasProviderAncestor(path, repoPath) {
					addRef(canonical, "provider", path, "directory: "+relativePath(repoPath, path))
				}
			}
			return nil
		}

		// Check file names (e.g., vllm_adapter.py, kserve_client.go)
		if !strings.HasSuffix(nameLower, ".py") && !strings.HasSuffix(nameLower, ".go") {
			return nil
		}

		baseName := strings.TrimSuffix(strings.TrimSuffix(nameLower, ".py"), ".go")
		// Strip common suffixes to isolate the component reference
		for _, suffix := range []string{"_adapter", "_client", "_provider", "_plugin", "_connector", "_driver", "_backend", "_config", "_utils"} {
			stripped := strings.TrimSuffix(baseName, suffix)
			if stripped != baseName {
				if canonical, ok := componentSet[stripped]; ok {
					typ := strings.Trim(suffix, "_")
					addRef(canonical, typ, path, "file: "+relativePath(repoPath, path))
				}
			}
		}

		return nil
	})

	// Scan Go and Python source files for import statements referencing components
	scanImportsForComponentRefs(repoPath, componentSet, addRef)

	sort.Slice(refs, func(i, j int) bool {
		if refs[i].Target != refs[j].Target {
			return refs[i].Target < refs[j].Target
		}
		return refs[i].Type < refs[j].Type
	})

	return refs
}

// hasProviderAncestor checks if any ancestor directory of path (up to repoPath)
// matches a provider directory pattern.
func hasProviderAncestor(path, repoPath string) bool {
	current := filepath.Dir(path)
	for current != repoPath && current != "." && current != "/" {
		dirName := strings.ToLower(filepath.Base(current))
		for _, pattern := range providerDirPatterns {
			if dirName == pattern || strings.Contains(dirName, pattern) {
				return true
			}
		}
		current = filepath.Dir(current)
	}
	return false
}

// goImportRE matches Go import lines: "github.com/org/repo/..."
var goImportRE = regexp.MustCompile(`^\s*(?:\w+\s+)?"([^"]+)"`)

// pyImportRE matches Python import lines: from X import Y or import X
var pyImportRE = regexp.MustCompile(`^\s*(?:from|import)\s+([a-zA-Z0-9_.]+)`)

// scanImportsForComponentRefs scans Go and Python source files for import
// statements that reference known component names. This detects cross-component
// dependencies that aren't visible from directory structure alone.
func scanImportsForComponentRefs(repoPath string, componentSet map[string]string, addRef func(string, string, string, string)) {
	filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			if info != nil && info.IsDir() {
				name := info.Name()
				if strings.HasPrefix(name, ".") || name == "vendor" || name == "node_modules" ||
					name == "__pycache__" || name == "site-packages" || name == "venv" ||
					name == ".venv" || name == "testdata" {
					return filepath.SkipDir
				}
			}
			return nil
		}

		nameLower := strings.ToLower(info.Name())
		isGo := strings.HasSuffix(nameLower, ".go")
		isPy := strings.HasSuffix(nameLower, ".py")
		if !isGo && !isPy {
			return nil
		}

		// Skip test files to reduce noise
		if isGo && strings.HasSuffix(nameLower, "_test.go") {
			return nil
		}
		if isPy && (strings.HasPrefix(nameLower, "test_") || strings.HasSuffix(nameLower, "_test.py")) {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		relPath := relativePath(repoPath, path)

		if isGo {
			scanGoImports(scanner, relPath, componentSet, addRef)
		} else {
			scanPyImports(scanner, relPath, componentSet, addRef)
		}
		return nil
	})
}

// scanGoImports parses Go import statements and checks for component references.
func scanGoImports(scanner *bufio.Scanner, relPath string, componentSet map[string]string, addRef func(string, string, string, string)) {
	inImportBlock := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "import (" {
			inImportBlock = true
			continue
		}
		if inImportBlock && line == ")" {
			inImportBlock = false
			continue
		}

		// Single-line import
		if strings.HasPrefix(line, "import ") && !strings.Contains(line, "(") {
			if match := goImportRE.FindStringSubmatch(line[7:]); match != nil {
				checkGoImportPath(match[1], relPath, componentSet, addRef)
			}
			continue
		}

		if inImportBlock {
			if match := goImportRE.FindStringSubmatch(line); match != nil {
				checkGoImportPath(match[1], relPath, componentSet, addRef)
			}
		}

		// Stop scanning after imports (Go imports must be at top of file)
		if !inImportBlock && !strings.HasPrefix(line, "import") &&
			!strings.HasPrefix(line, "package") && !strings.HasPrefix(line, "//") &&
			line != "" {
			break
		}
	}
}

// checkGoImportPath checks if a Go import path references a known component.
func checkGoImportPath(importPath, relPath string, componentSet map[string]string, addRef func(string, string, string, string)) {
	// Split import path into segments and check each against component names
	parts := strings.Split(strings.ToLower(importPath), "/")
	for _, part := range parts {
		if canonical, ok := componentSet[part]; ok {
			addRef(canonical, "import", relPath, "go import: "+importPath)
			return
		}
		// Also check with hyphens replaced (Go packages use underscores)
		hyphenated := strings.ReplaceAll(part, "_", "-")
		if hyphenated != part {
			if canonical, ok := componentSet[hyphenated]; ok {
				addRef(canonical, "import", relPath, "go import: "+importPath)
				return
			}
		}
	}
}

// scanPyImports parses Python import statements and checks for component references.
func scanPyImports(scanner *bufio.Scanner, relPath string, componentSet map[string]string, addRef func(string, string, string, string)) {
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Stop at class/function definitions (imports are at the top)
		if strings.HasPrefix(line, "class ") || strings.HasPrefix(line, "def ") ||
			strings.HasPrefix(line, "async def ") {
			break
		}

		if match := pyImportRE.FindStringSubmatch(line); match != nil {
			modulePath := match[1]
			// Split module path and check each component against known names
			parts := strings.Split(strings.ToLower(modulePath), ".")
			for _, part := range parts {
				// Normalize: Python uses underscores, component names use hyphens
				hyphenated := strings.ReplaceAll(part, "_", "-")
				if canonical, ok := componentSet[hyphenated]; ok {
					addRef(canonical, "import", relPath, "python import: "+modulePath)
					break
				}
				if canonical, ok := componentSet[part]; ok {
					addRef(canonical, "import", relPath, "python import: "+modulePath)
					break
				}
			}
		}
	}
}

// defaultKnownComponents returns a static list of well-known RHOAI/ODH component
// names used as fallback when scan-config component list isn't provided.
func defaultKnownComponents() []string {
	return []string{
		"vllm", "vllm-cpu",
		"kserve", "kserve-autogluon-server",
		"modelmesh", "modelmesh-serving",
		"model-registry", "model-registry-operator",
		"data-science-pipelines", "data-science-pipelines-operator",
		"codeflare-operator", "codeflare-sdk",
		"kuberay", "kueue", "kubeflow",
		"training-operator", "trainer",
		"odh-dashboard", "odh-model-controller",
		"opendatahub-operator",
		"trustyai-service-operator",
		"lm-evaluation-harness",
		"llama-stack", "llama-stack-k8s-operator",
		"mlflow", "mlflow-operator",
		"notebooks",
		"spark-operator",
		"fms-guardrails-orchestrator", "guardrails-detectors",
		"kube-rbac-proxy", "kube-auth-proxy",
		"distributed-workloads",
		"argo-workflows",
		"batch-gateway",
		"llm-d-inference-scheduler", "llm-d-kv-cache",
		"workload-variant-autoscaler",
		"eval-hub",
		"fms-hf-tuning",
		"models-as-a-service",
	}
}
