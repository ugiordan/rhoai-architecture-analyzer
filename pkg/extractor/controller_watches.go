package extractor

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var controllerGoPatterns = []string{
	"**/*_controller.go",
	"**/*_reconciler.go",
	"**/setup.go",
	"**/controller.go",
	"**/reconciler.go",
}

var (
	importAliasRE    = regexp.MustCompile(`(\w+)\s+"([^"]+)"`)
	forRE            = regexp.MustCompile(`\.?For\(\s*&(\w+)\.(\w+)\{`)
	ownsRE           = regexp.MustCompile(`\.?Owns\(\s*&(\w+)\.(\w+)\{`)
	watchesRE        = regexp.MustCompile(`\.?Watches\(\s*&?(?:source\.Kind\{Type:\s*&)?(\w+)\.(\w+)\{`)
	setupFuncRE      = regexp.MustCompile(`func\s+\(\s*\w+\s+\*(\w+)\)\s+SetupWithManager`)
	reconcilerNameRE = regexp.MustCompile(`func\s+\(\s*\w+\s+\*(\w+)\)\s+Reconcile\b`)
)

// knownGroups maps Go import paths to Kubernetes API group/version strings.
var knownGroups = map[string]string{
	"k8s.io/api/core/v1":                     "/v1",
	"k8s.io/api/apps/v1":                     "apps/v1",
	"k8s.io/api/batch/v1":                    "batch/v1",
	"k8s.io/api/networking/v1":               "networking.k8s.io/v1",
	"k8s.io/api/rbac/v1":                     "rbac.authorization.k8s.io/v1",
	"k8s.io/api/policy/v1":                   "policy/v1",
	"k8s.io/api/autoscaling/v1":              "autoscaling/v1",
	"k8s.io/api/autoscaling/v2":              "autoscaling/v2",
	"k8s.io/api/admissionregistration/v1":    "admissionregistration.k8s.io/v1",
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1": "apiextensions.k8s.io/v1",
}

// extractControllerWatches scans Go controller files for For/Owns/Watches
// patterns and resolves the import aliases to API group/version/kind.
func extractControllerWatches(repoPath string) []ControllerWatch {
	files := findFiles(repoPath, controllerGoPatterns)
	var watches []ControllerWatch

	for _, fpath := range files {
		data, err := os.ReadFile(fpath)
		if err != nil {
			continue
		}
		content := string(data)
		imports := parseImports(content)
		lines := strings.Split(content, "\n")

		relPath := relativePath(repoPath, fpath)

		// Detect controller name from SetupWithManager or Reconcile method receivers
		controllerName := detectControllerName(content)

		for lineNo, line := range lines {
			source := fmt.Sprintf("%s:%d", relPath, lineNo+1)

			for _, match := range forRE.FindAllStringSubmatch(line, -1) {
				alias, kind := match[1], match[2]
				gv := resolveImportAlias(alias, imports)
				watches = append(watches, ControllerWatch{
					Type:       "For",
					GVK:        fmt.Sprintf("%s/%s", gv, kind),
					Controller: controllerName,
					Source:     source,
				})
			}

			for _, match := range ownsRE.FindAllStringSubmatch(line, -1) {
				alias, kind := match[1], match[2]
				gv := resolveImportAlias(alias, imports)
				watches = append(watches, ControllerWatch{
					Type:       "Owns",
					GVK:        fmt.Sprintf("%s/%s", gv, kind),
					Controller: controllerName,
					Source:     source,
				})
			}

			for _, match := range watchesRE.FindAllStringSubmatch(line, -1) {
				alias, kind := match[1], match[2]
				gv := resolveImportAlias(alias, imports)
				watches = append(watches, ControllerWatch{
					Type:       "Watches",
					GVK:        fmt.Sprintf("%s/%s", gv, kind),
					Controller: controllerName,
					Source:     source,
				})
			}
		}
	}

	if watches == nil {
		watches = []ControllerWatch{}
	}
	return watches
}

// detectControllerName extracts the reconciler struct name from a Go source file
// by looking for SetupWithManager or Reconcile method receivers.
func detectControllerName(content string) string {
	// Prefer SetupWithManager since it's where For/Owns/Watches live
	if m := setupFuncRE.FindStringSubmatch(content); m != nil {
		return m[1]
	}
	if m := reconcilerNameRE.FindStringSubmatch(content); m != nil {
		return m[1]
	}
	return ""
}

var (
	importBlockRE = regexp.MustCompile(`(?s)import\s*\((.*?)\)`)
	pathOnlyRE    = regexp.MustCompile(`^"([^"]+)"$`)
	versionRE     = regexp.MustCompile(`^v\d+`)
)

// parseImports extracts Go import alias to path mappings from source content.
func parseImports(content string) map[string]string {
	imports := make(map[string]string)

	blocks := importBlockRE.FindAllStringSubmatch(content, -1)
	for _, block := range blocks {
		for _, line := range strings.Split(block[1], "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "//") {
				continue
			}
			if match := importAliasRE.FindStringSubmatch(line); match != nil {
				imports[match[1]] = match[2]
			} else if match := pathOnlyRE.FindStringSubmatch(line); match != nil {
				path := match[1]
				parts := strings.Split(strings.TrimRight(path, "/"), "/")
				last := parts[len(parts)-1]
				imports[last] = path
			}
		}
	}
	return imports
}

// resolveImportAlias resolves an import alias to an API group/version string.
func resolveImportAlias(alias string, imports map[string]string) string {
	importPath, ok := imports[alias]
	if !ok {
		return alias
	}

	if gv, exists := knownGroups[importPath]; exists {
		return gv
	}

	// Try to infer group/version from import path
	parts := strings.Split(strings.TrimRight(importPath, "/"), "/")
	if len(parts) >= 2 {
		version := parts[len(parts)-1]
		matched := versionRE.MatchString(version)
		if matched {
			groupPart := ""
			if len(parts) >= 3 {
				groupPart = parts[len(parts)-2]
			}
			if groupPart != "" {
				return fmt.Sprintf("%s/%s", groupPart, version)
			}
			return fmt.Sprintf("/%s", version)
		}
	}

	return importPath
}
