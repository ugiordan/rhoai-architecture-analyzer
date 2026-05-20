package extractor

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/tools/go/packages"
)

// GoPackageSet holds loaded Go packages with type information.
type GoPackageSet struct {
	Packages []*packages.Package
	Fset     *token.FileSet
	Mode     string // "full" or "syntax-only"
	Warning  string
}

// MarkedStruct is a struct type whose doc comment contains a kubebuilder marker.
type MarkedStruct struct {
	Name     string
	Doc      string
	Markers  []string
	Pkg      *packages.Package
	TypeSpec *ast.TypeSpec
	File     *ast.File
}

// loadGoPackages loads all Go packages under repoPath using go/packages.
// Returns nil if repoPath has no go.mod (not a Go project).
func loadGoPackages(repoPath string) *GoPackageSet {
	// Quick check: bail if no go.mod
	if _, err := os.Stat(filepath.Join(repoPath, "go.mod")); err != nil {
		return nil
	}

	absRepo, err := filepath.Abs(repoPath)
	if err != nil {
		return nil
	}

	// Run go mod download with security-hardened env and 2min timeout.
	// Note: go mod download needs -mod=mod (not readonly) to actually fetch.
	dlCtx, dlCancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer dlCancel()

	dlCmd := exec.CommandContext(dlCtx, "go", "mod", "download")
	dlCmd.Dir = absRepo
	dlCmd.Env = loaderEnv(absRepo, false)
	if out, err := dlCmd.CombinedOutput(); err != nil {
		outStr := string(out)
		if strings.Contains(outStr, "SECURITY ERROR") || strings.Contains(outStr, "checksum mismatch") {
			log.Printf("[go_loader] checksum verification failed for %s, skipping Go AST", absRepo)
			return nil
		}
		log.Printf("[go_loader] go mod download failed (non-security): %v", err)
	}

	// Re-check for symlinks after go mod download (TOCTOU mitigation).
	// A malicious go.mod replace directive could create symlinks outside the repo.
	if err := checkRepoSymlinks(absRepo); err != nil {
		log.Printf("[go_loader] symlink boundary violation after go mod download: %v", err)
		return nil
	}

	// Load packages with 5min timeout
	loadCtx, loadCancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer loadCancel()

	fset := token.NewFileSet()
	cfg := &packages.Config{
		Mode: packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo |
			packages.NeedFiles |
			packages.NeedName |
			packages.NeedImports |
			packages.NeedCompiledGoFiles,
		Dir:     absRepo,
		Fset:    fset,
		Context: loadCtx,
		Env:     loaderEnv(absRepo, true),
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		// Total failure to load
		return nil
	}
	if len(pkgs) == 0 {
		return nil
	}

	// Filter out test packages and count errors
	var valid []*packages.Package
	errCount := 0
	for _, pkg := range pkgs {
		if strings.HasSuffix(pkg.PkgPath, "_test") || strings.HasSuffix(pkg.Name, "_test") {
			continue
		}
		if len(pkg.Errors) > 0 {
			errCount++
		}
		valid = append(valid, pkg)
	}

	if len(valid) == 0 {
		return nil
	}

	// If >50% of packages have errors, fall back to syntax-only mode
	if errCount > 0 && errCount*2 > len(valid) {
		return &GoPackageSet{
			Packages: valid,
			Fset:     fset,
			Mode:     "syntax-only",
			Warning:  "majority of packages had type-checking errors, using syntax-only mode",
		}
	}

	return &GoPackageSet{
		Packages: valid,
		Fset:     fset,
		Mode:     "full",
	}
}

// loaderEnv builds a sanitized environment for go commands.
func loaderEnv(absRepoPath string, readonly bool) []string {
	env := os.Environ()
	gopath := filepath.Join(absRepoPath, ".gopath-loader")
	overrides := map[string]string{
		"CGO_ENABLED":  "0",
		"GONOSUMCHECK": "",
		"GONOSUMDB":    "",
		"GOMAXPROCS":   "2",
		"GOPATH":       gopath,
		"GOMODCACHE":   filepath.Join(gopath, "pkg", "mod"),
		"GOCACHE":      filepath.Join(gopath, "cache"),
		"GOPRIVATE":    "",
		"GONOPROXY":    "",
	}
	if readonly {
		overrides["GOFLAGS"] = "-mod=readonly"
	} else {
		overrides["GOFLAGS"] = ""
	}
	// Remove existing keys, then append overrides
	filtered := make([]string, 0, len(env)+len(overrides))
	for _, e := range env {
		key := strings.SplitN(e, "=", 2)[0]
		if _, skip := overrides[key]; skip {
			continue
		}
		filtered = append(filtered, e)
	}
	for k, v := range overrides {
		filtered = append(filtered, k+"="+v)
	}
	return filtered
}

// FindStructsWithMarker scans all packages for struct types whose doc comment
// contains the given marker string (e.g., "+kubebuilder:object:root=true").
func (gps *GoPackageSet) FindStructsWithMarker(marker string) []MarkedStruct {
	var result []MarkedStruct
	for _, pkg := range gps.Packages {
		for _, file := range pkg.Syntax {
			for _, decl := range file.Decls {
				gd, ok := decl.(*ast.GenDecl)
				if !ok || gd.Tok != token.TYPE {
					continue
				}
				for _, spec := range gd.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					// Only care about struct types
					if _, isStruct := ts.Type.(*ast.StructType); !isStruct {
						continue
					}
					// Collect doc from the GenDecl (group doc) and the TypeSpec doc
					doc := docText(gd.Doc) + docText(ts.Doc)
					if !strings.Contains(doc, marker) {
						continue
					}
					markers := extractMarkers(doc)
					result = append(result, MarkedStruct{
						Name:     ts.Name.Name,
						Doc:      strings.TrimSpace(doc),
						Markers:  markers,
						Pkg:      pkg,
						TypeSpec: ts,
						File:     file,
					})
				}
			}
		}
	}
	return result
}

// ResolveType uses the package's type info to resolve an ast.Expr to its
// fully qualified type name. Falls back to printing the expression if
// type info is unavailable.
func ResolveType(expr ast.Expr, pkg *packages.Package) string {
	if pkg.TypesInfo != nil {
		if t := pkg.TypesInfo.TypeOf(expr); t != nil {
			return t.String()
		}
	}
	return exprString(expr)
}

// FindMethodsOnType finds all methods declared on the given type name
// within the package (including pointer receivers).
func FindMethodsOnType(pkg *packages.Package, typeName string) []*ast.FuncDecl {
	var methods []*ast.FuncDecl
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || len(fn.Recv.List) == 0 {
				continue
			}
			recvName := receiverTypeName(fn.Recv.List[0].Type)
			if recvName == typeName {
				methods = append(methods, fn)
			}
		}
	}
	return methods
}

// receiverTypeName extracts the type name from a receiver expression,
// handling both value and pointer receivers.
func receiverTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return receiverTypeName(t.X)
	case *ast.Ident:
		return t.Name
	case *ast.IndexExpr: // generic type T[X]
		return receiverTypeName(t.X)
	case *ast.IndexListExpr: // generic type T[X, Y]
		return receiverTypeName(t.X)
	default:
		return ""
	}
}

// derefPtr unwraps pointer types to get the underlying named type.
func derefPtr(t types.Type) types.Type {
	for {
		p, ok := t.(*types.Pointer)
		if !ok {
			return t
		}
		t = p.Elem()
	}
}

// docText returns the text of a comment group, or empty string if nil.
func docText(cg *ast.CommentGroup) string {
	if cg == nil {
		return ""
	}
	return cg.Text()
}

// extractMarkers pulls all lines starting with "+" from a doc string.
func extractMarkers(doc string) []string {
	var markers []string
	for _, line := range strings.Split(doc, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "+") {
			markers = append(markers, line)
		}
	}
	return markers
}

// exprString is a simple fallback for printing ast.Expr when type info is unavailable.
func exprString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return exprString(t.X) + "." + t.Sel.Name
	case *ast.StarExpr:
		return "*" + exprString(t.X)
	case *ast.ArrayType:
		return "[]" + exprString(t.Elt)
	case *ast.MapType:
		return "map[" + exprString(t.Key) + "]" + exprString(t.Value)
	default:
		return "unknown"
	}
}

// checkRepoSymlinks walks the repo directory and returns an error if any
// symlink points outside the repo boundary.
func checkRepoSymlinks(repoPath string) error {
	return filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.Mode()&os.ModeSymlink != 0 {
			target, evalErr := filepath.EvalSymlinks(path)
			if evalErr != nil {
				return fmt.Errorf("cannot resolve symlink %s: %w", path, evalErr)
			}
			if !strings.HasPrefix(target, repoPath) {
				return fmt.Errorf("symlink %s escapes repo boundary: %s", path, target)
			}
		}
		return nil
	})
}
