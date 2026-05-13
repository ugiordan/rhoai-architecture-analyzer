package extractor

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

// extractCRDsFromGo scans loaded Go packages for types annotated with
// +kubebuilder:object:root=true and returns a CRD for each one (skipping
// List types). Group/version are resolved from the package's GroupVersion
// variable, scope from the resource marker, storage version from the
// storageversion marker, and field counts are computed recursively via
// type info when available.
func extractCRDsFromGo(pkgs *GoPackageSet) []CRD {
	if pkgs == nil {
		return nil
	}

	roots := pkgs.FindStructsWithMarker("+kubebuilder:object:root=true")

	var crds []CRD
	for _, root := range roots {
		if strings.HasSuffix(root.Name, "List") {
			continue
		}

		group, version := resolveGroupVersion(root.Pkg)
		scope := "Namespaced"
		isStorage := false

		for _, m := range root.Markers {
			if strings.Contains(m, "resource:") && strings.Contains(m, "scope=Cluster") {
				scope = "Cluster"
			}
			if strings.Contains(m, "storageversion") {
				isStorage = true
			}
		}

		celRules := extractCELRulesFromDoc(root.Doc)
		fieldCount := countStructFields(root.TypeSpec, root.Pkg)

		// Detect hub/spoke conversion methods
		hubVersion := ""
		var spokeVersions []string
		methods := FindMethodsOnType(root.Pkg, root.Name)
		for _, m := range methods {
			switch m.Name.Name {
			case "Hub":
				hubVersion = version
			case "ConvertTo", "ConvertFrom":
				if !containsStr(spokeVersions, version) {
					spokeVersions = append(spokeVersions, version)
				}
			}
		}

		source := pkgs.Fset.Position(root.TypeSpec.Pos()).Filename

		crd := CRD{
			Group:           group,
			Version:         version,
			Kind:            root.Name,
			Scope:           scope,
			FieldsCount:     fieldCount,
			ValidationRules: celRules,
			Source:          source,
			GoSource:        "go_ast",
			HubVersion:      hubVersion,
			SpokeVersions:   spokeVersions,
			Versions: []CRDVersion{{
				Name:    version,
				Served:  true,
				Storage: isStorage,
			}},
		}
		crds = append(crds, crd)
	}

	return crds
}

// resolveGroupVersion extracts the API group and version from a Go package.
// It searches all files in the package for a GroupVersion or SchemeGroupVersion
// variable declaration with a composite literal containing a "Group" field.
// The version is inferred from the package path (e.g. "api/v1alpha1" -> "v1alpha1").
func resolveGroupVersion(pkg *packages.Package) (string, string) {
	group := ""
	version := ""

	// Extract version from package path segments
	parts := strings.Split(pkg.PkgPath, "/")
	for _, p := range parts {
		if len(p) > 1 && p[0] == 'v' && p[1] >= '0' && p[1] <= '9' {
			version = p
		}
	}

	// Search all files for GroupVersion or SchemeGroupVersion var
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range gd.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for _, name := range vs.Names {
					if name.Name == "GroupVersion" || name.Name == "SchemeGroupVersion" {
						group = extractGroupFromComposite(vs, pkg)
						if group != "" {
							return group, version
						}
					}
				}
			}
		}
	}

	// Fallback: try to infer group from SchemeBuilder AddToScheme registration
	// or from type info if available
	if group == "" && pkg.TypesInfo != nil {
		for _, file := range pkg.Syntax {
			for _, decl := range file.Decls {
				gd, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}
				for _, spec := range gd.Specs {
					vs, ok := spec.(*ast.ValueSpec)
					if !ok || len(vs.Values) == 0 {
						continue
					}
					for _, name := range vs.Names {
						if name.Name == "SchemeBuilder" || name.Name == "AddToScheme" {
							group = extractGroupFromSchemeBuilder(vs, pkg)
							if group != "" {
								return group, version
							}
						}
					}
				}
			}
		}
	}

	return group, version
}

func extractGroupFromSchemeBuilder(vs *ast.ValueSpec, pkg *packages.Package) string {
	if len(vs.Values) == 0 {
		return ""
	}
	var found string
	ast.Inspect(vs.Values[0], func(n ast.Node) bool {
		if found != "" {
			return false
		}
		if ident, ok := n.(*ast.Ident); ok && pkg.TypesInfo != nil {
			if obj, use := pkg.TypesInfo.Uses[ident]; use {
				if c, ok := obj.(*types.Const); ok {
					val := strings.Trim(c.Val().String(), `"`)
					if strings.Contains(val, ".") && !strings.Contains(val, "/") {
						found = val
						return false
					}
				}
			}
		}
		return true
	})
	return found
}

// extractGroupFromComposite extracts the "Group" field value from a composite
// literal in a var declaration (e.g. schema.GroupVersion{Group: "apps.example.com", ...}).
func extractGroupFromComposite(vs *ast.ValueSpec, pkg *packages.Package) string {
	if len(vs.Values) == 0 {
		return ""
	}
	comp, ok := vs.Values[0].(*ast.CompositeLit)
	if !ok {
		return ""
	}
	for _, elt := range comp.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok || key.Name != "Group" {
			continue
		}
		// Direct string literal
		if lit, ok := kv.Value.(*ast.BasicLit); ok {
			return strings.Trim(lit.Value, `"`)
		}
		// Constant reference: resolve via type info
		if ident, ok := kv.Value.(*ast.Ident); ok && pkg.TypesInfo != nil {
			if obj, found := pkg.TypesInfo.Uses[ident]; found {
				if c, ok := obj.(*types.Const); ok {
					return strings.Trim(c.Val().String(), `"`)
				}
			}
		}
	}
	return ""
}

// extractCELRulesFromDoc pulls CEL validation rules from doc comments.
func extractCELRulesFromDoc(doc string) []string {
	var rules []string
	for _, line := range strings.Split(doc, "\n") {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "XValidation") || strings.Contains(line, "x-kubernetes-validations") {
			rules = append(rules, line)
		}
	}
	return rules
}

// countStructFields counts the fields in a struct type, recursively following
// embedded types via the package's type info when available. Falls back to
// counting AST fields directly when type info is nil.
func countStructFields(ts *ast.TypeSpec, pkg *packages.Package) int {
	st, ok := ts.Type.(*ast.StructType)
	if !ok {
		return 0
	}
	count := 0
	for _, f := range st.Fields.List {
		if len(f.Names) == 0 {
			// Embedded field: resolve through type info if available
			if pkg.TypesInfo != nil {
				t := pkg.TypesInfo.TypeOf(f.Type)
				if t != nil {
					count += countTypeFields(derefPtr(t))
					continue
				}
			}
			// Fallback: count embedded as 1
			count++
			continue
		}
		count += len(f.Names)
	}
	return count
}

// countTypeFields recursively counts fields in a resolved type, expanding
// embedded struct fields.
func countTypeFields(t types.Type) int {
	switch u := t.Underlying().(type) {
	case *types.Struct:
		count := 0
		for i := 0; i < u.NumFields(); i++ {
			f := u.Field(i)
			if f.Embedded() {
				count += countTypeFields(f.Type())
			} else {
				count++
			}
		}
		return count
	default:
		return 1
	}
}

// containsStr checks if a string slice contains the given string.
func containsStr(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}
