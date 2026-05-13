package extractor

import (
	"fmt"
	"go/ast"
	"os"
	"regexp"
	"strings"

	"golang.org/x/tools/go/packages"
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

// --- AST-based resource operation extraction ---

// clientVerbs maps controller-runtime client method names to CRUD verbs.
var clientVerbs = map[string]string{
	"Create": "create",
	"Update": "update",
	"Patch":  "patch",
	"Delete": "delete",
}

// k8sAPIGroupMap maps known k8s.io/api import paths to Kubernetes API groups.
var k8sAPIGroupMap = map[string]string{
	"k8s.io/api/core/v1":                     "",
	"k8s.io/api/apps/v1":                     "apps",
	"k8s.io/api/batch/v1":                    "batch",
	"k8s.io/api/networking/v1":               "networking.k8s.io",
	"k8s.io/api/rbac/v1":                     "rbac.authorization.k8s.io",
	"k8s.io/api/policy/v1":                   "policy",
	"k8s.io/api/autoscaling/v1":              "autoscaling",
	"k8s.io/api/autoscaling/v2":              "autoscaling",
	"k8s.io/api/admissionregistration/v1":    "admissionregistration.k8s.io",
	"k8s.io/api/storage/v1":                  "storage.k8s.io",
	"k8s.io/api/coordination/v1":             "coordination.k8s.io",
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1": "apiextensions.k8s.io",
}

// extractResourceOps scans the loaded Go packages for programmatic Kubernetes
// resource operations (Create, Update, Patch, Delete) inside Reconcile methods,
// resolving argument types via go/packages type information.
func extractResourceOps(pkgs *GoPackageSet) []ResourceOp {
	if pkgs == nil {
		return nil
	}

	var ops []ResourceOp

	for _, pkg := range pkgs.Packages {
		if !isControllerPackage(pkg.PkgPath) {
			continue
		}

		for _, file := range pkg.Syntax {
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok {
					continue
				}
				if !isReconcileMethod(fn) {
					continue
				}

				found := extractOpsFromFunc(fn, pkg, pkgs)
				ops = append(ops, found...)
			}
		}
	}

	return dedupeResourceOps(ops)
}

// isControllerPackage returns true if the package path suggests it contains
// controller/reconciler logic.
func isControllerPackage(path string) bool {
	return strings.Contains(path, "controller") ||
		strings.Contains(path, "reconciler") ||
		strings.Contains(path, "internal/")
}

// isReconcileMethod returns true if the function declaration is a method named
// Reconcile or starting with Reconcile/reconcile (and has a receiver).
func isReconcileMethod(fn *ast.FuncDecl) bool {
	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return false
	}
	name := fn.Name.Name
	return name == "Reconcile" ||
		strings.HasPrefix(name, "Reconcile") ||
		strings.HasPrefix(name, "reconcile")
}

// extractOpsFromFunc walks a function's AST looking for client.Create/Update/Patch/Delete
// calls and resolves the resource type of each call's object argument.
func extractOpsFromFunc(fn *ast.FuncDecl, pkg *packages.Package, pkgs *GoPackageSet) []ResourceOp {
	var ops []ResourceOp

	ast.Inspect(fn.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		verb, isClientVerb := clientVerbs[sel.Sel.Name]
		if !isClientVerb {
			return true
		}

		// The object argument is typically the second argument (after ctx).
		// Create(ctx, obj, ...opts), Update(ctx, obj, ...opts), etc.
		if len(call.Args) < 2 {
			return true
		}

		objArg := call.Args[1]
		fullType := ResolveType(objArg, pkg)

		// Strip pointer prefix
		fullType = strings.TrimPrefix(fullType, "*")

		kind, group, isK8s := splitTypeToKindGroup(fullType)
		if kind == "" || !isK8s {
			return true
		}

		pos := pkgs.Fset.Position(call.Pos())
		source := SourceRef{
			Type: "go_ast",
			File: pos.Filename,
			Line: pos.Line,
		}

		ops = append(ops, ResourceOp{
			Kind:   kind,
			Group:  group,
			Verb:   verb,
			Source: source,
		})

		return true
	})

	return ops
}

// splitTypeToKindGroup splits a fully qualified Go type name like
// "k8s.io/api/core/v1.Service" into kind="Service", group="", and isK8s=true.
// The third return value indicates whether the type is from a known Kubernetes
// API package (prevents false positives from non-client methods like cache.Delete).
func splitTypeToKindGroup(fullType string) (string, string, bool) {
	dotIdx := strings.LastIndex(fullType, ".")
	if dotIdx < 0 {
		return fullType, "", false
	}

	kind := fullType[dotIdx+1:]
	pkgPath := fullType[:dotIdx]

	if group, ok := k8sAPIGroupMap[pkgPath]; ok {
		return kind, group, true
	}

	// Accept types from packages with k8s-style API paths
	isK8sStyle := strings.Contains(pkgPath, "k8s.io/api") ||
		strings.Contains(pkgPath, "/api/") ||
		strings.Contains(pkgPath, "/apis/")

	if !isK8sStyle {
		return kind, "", false
	}

	parts := strings.Split(pkgPath, "/")
	if len(parts) >= 2 {
		lastSeg := parts[len(parts)-1]
		if versionRE.MatchString(lastSeg) && len(parts) >= 3 {
			return kind, parts[len(parts)-2], true
		}
		return kind, lastSeg, true
	}

	return kind, pkgPath, true
}

// dedupeResourceOps removes duplicate ResourceOps based on kind+group+verb.
func dedupeResourceOps(ops []ResourceOp) []ResourceOp {
	type key struct {
		Kind  string
		Group string
		Verb  string
	}
	seen := make(map[key]bool)
	var result []ResourceOp

	for _, op := range ops {
		k := key{Kind: op.Kind, Group: op.Group, Verb: op.Verb}
		if seen[k] {
			continue
		}
		seen[k] = true
		result = append(result, op)
	}

	return result
}
