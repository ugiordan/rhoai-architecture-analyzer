package extractor

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"strings"
)

// apiTypesPatterns locates Go files that define Kubernetes API types.
// Covers the standard kubebuilder/operator-sdk layout (api/<version>/*_types.go)
// and upstream Kubernetes API conventions (pkg/apis/**/*_types.go).
var apiTypesPatterns = []string{
	"api/**/*_types.go",
	"apis/**/*_types.go",
	"pkg/apis/**/*_types.go",
	"pkg/api/**/*_types.go",
}

// extractAPITypes parses *_types.go files using go/ast and returns struct
// definitions with fields, doc comments, JSON tags, and kubebuilder markers.
func extractAPITypes(repoPath string) []APITypeDefinition {
	files := findFiles(repoPath, apiTypesPatterns)
	if len(files) == 0 {
		return nil
	}

	var result []APITypeDefinition
	seen := make(map[string]bool)

	for _, fpath := range files {
		defs := parseAPITypesFile(fpath, repoPath)
		for _, d := range defs {
			key := d.Name + ":" + d.Source
			if seen[key] {
				continue
			}
			seen[key] = true
			result = append(result, d)
		}
	}

	return result
}

// parseAPITypesFile parses a single Go file and extracts struct type definitions
// that look like Kubernetes API types (have JSON tags or kubebuilder markers).
func parseAPITypesFile(path, repoPath string) []APITypeDefinition {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		log.Printf("warning: cannot parse %s: %v", path, err)
		return nil
	}

	relPath := relativePath(repoPath, path)
	var defs []APITypeDefinition

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			// Extract doc comment from the GenDecl (appears before the type keyword)
			// or from the TypeSpec itself (inline).
			doc := extractDocComment(genDecl.Doc, typeSpec.Doc)

			// Extract kubebuilder markers from doc comments
			markers := extractMarkersFromCommentGroup(genDecl.Doc)

			fields := extractStructFields(structType, fset)

			// Skip structs with no JSON-tagged fields (not API types)
			if !hasJSONFields(fields) && len(markers) == 0 {
				continue
			}

			name := typeSpec.Name.Name
			line := fset.Position(typeSpec.Pos()).Line

			def := APITypeDefinition{
				Name:    name,
				Doc:     doc,
				Fields:  fields,
				Markers: markers,
				Source:  fmt.Sprintf("%s:%d", relPath, line),
				IsSpec:  strings.HasSuffix(name, "Spec"),
				IsStatus: strings.HasSuffix(name, "Status"),
			}

			defs = append(defs, def)
		}
	}

	return defs
}

// extractStructFields extracts field metadata from a struct type.
func extractStructFields(st *ast.StructType, fset *token.FileSet) []APIField {
	if st.Fields == nil {
		return nil
	}

	var fields []APIField
	for _, field := range st.Fields.List {
		f := parseStructField(field, fset)
		if f != nil {
			fields = append(fields, *f)
		}
	}
	return fields
}

// parseStructField parses a single struct field into an APIField.
func parseStructField(field *ast.Field, fset *token.FileSet) *APIField {
	goType := typeExprToString(field.Type)

	// Get field name (embedded fields have no names)
	name := ""
	embedded := false
	if len(field.Names) > 0 {
		name = field.Names[0].Name
	} else {
		// Embedded field: use the type name
		name = embeddedTypeName(field.Type)
		embedded = true
	}

	if name == "" {
		return nil
	}

	// Parse JSON tag
	jsonTag := ""
	if field.Tag != nil {
		jsonTag = extractJSONTag(field.Tag.Value)
	}

	// Skip fields with json:"-" (not serialized)
	if jsonTag == "-" {
		return nil
	}

	// Extract doc comment
	doc := extractDocComment(field.Doc, field.Comment)

	// Extract kubebuilder markers
	markers := extractMarkersFromCommentGroup(field.Doc)

	// Determine if required (from markers or no omitempty)
	required := false
	hasDefault := ""
	for _, m := range markers {
		if strings.Contains(m, "validation:Required") || strings.Contains(m, "validation:required") {
			required = true
		}
		if strings.Contains(m, "default:") || strings.Contains(m, "default=") {
			// Extract default value
			if idx := strings.Index(m, "default:"); idx >= 0 {
				hasDefault = strings.TrimSpace(m[idx+len("default:"):])
			} else if idx := strings.Index(m, "default="); idx >= 0 {
				hasDefault = strings.TrimSpace(m[idx+len("default="):])
			}
		}
	}

	// Also mark required if JSON tag has no omitempty
	if jsonTag != "" && !strings.Contains(jsonTag, ",omitempty") && !embedded {
		required = true
	}

	// Detect secret references from type name
	secretRef := isSecretRefType(goType) || isSecretRefType(name)

	return &APIField{
		Name:      name,
		GoType:    goType,
		JSONTag:   jsonTag,
		Doc:       doc,
		Markers:   markers,
		Required:  required,
		Default:   hasDefault,
		SecretRef: secretRef,
		Embedded:  embedded,
	}
}

// typeExprToString converts a Go AST type expression to a readable string.
func typeExprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return typeExprToString(t.X) + "." + t.Sel.Name
	case *ast.StarExpr:
		return "*" + typeExprToString(t.X)
	case *ast.ArrayType:
		return "[]" + typeExprToString(t.Elt)
	case *ast.MapType:
		return "map[" + typeExprToString(t.Key) + "]" + typeExprToString(t.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.StructType:
		return "struct{}"
	default:
		return "unknown"
	}
}

// embeddedTypeName extracts the base type name from an embedded field.
func embeddedTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return t.Sel.Name
	case *ast.StarExpr:
		return embeddedTypeName(t.X)
	default:
		return ""
	}
}

// extractJSONTag extracts the JSON field name from a struct tag string.
// Input is the raw tag including backticks, e.g. `json:"name,omitempty"`.
func extractJSONTag(rawTag string) string {
	tag := strings.Trim(rawTag, "`")
	// Find json:"..."
	idx := strings.Index(tag, `json:"`)
	if idx < 0 {
		return ""
	}
	rest := tag[idx+len(`json:"`):]
	end := strings.Index(rest, `"`)
	if end < 0 {
		return ""
	}
	return rest[:end]
}

// extractDocComment returns the cleaned text from doc comment groups.
func extractDocComment(groups ...*ast.CommentGroup) string {
	for _, g := range groups {
		if g == nil {
			continue
		}
		var lines []string
		for _, c := range g.List {
			text := c.Text
			// Strip comment prefix
			text = strings.TrimPrefix(text, "//")
			text = strings.TrimPrefix(text, " ")
			// Skip kubebuilder markers (they go in Markers)
			if strings.HasPrefix(text, "+kubebuilder:") || strings.HasPrefix(text, "+optional") || strings.HasPrefix(text, "+required") {
				continue
			}
			text = strings.TrimSpace(text)
			if text != "" {
				lines = append(lines, text)
			}
		}
		if len(lines) > 0 {
			return strings.Join(lines, " ")
		}
	}
	return ""
}

// extractMarkersFromCommentGroup pulls kubebuilder markers from comments.
func extractMarkersFromCommentGroup(group *ast.CommentGroup) []string {
	if group == nil {
		return nil
	}
	var markers []string
	for _, c := range group.List {
		text := strings.TrimPrefix(c.Text, "//")
		text = strings.TrimSpace(text)
		if strings.HasPrefix(text, "+kubebuilder:") {
			markers = append(markers, text)
		}
	}
	return markers
}

// isSecretRefType returns true if a type or field name indicates a Secret reference.
func isSecretRefType(name string) bool {
	lower := strings.ToLower(name)
	return strings.Contains(lower, "secret") ||
		strings.Contains(lower, "password") ||
		strings.Contains(lower, "credential")
}

// hasJSONFields returns true if at least one field has a JSON tag.
func hasJSONFields(fields []APIField) bool {
	for _, f := range fields {
		if f.JSONTag != "" {
			return true
		}
	}
	return false
}

// apiVersionFromPath extracts the API version directory from a file path.
// e.g., "api/v1/types.go" -> "v1", "pkg/apis/serving/v1alpha1/types.go" -> "v1alpha1"
func apiVersionFromPath(path string) string {
	dir := filepath.Dir(path)
	base := filepath.Base(dir)
	if strings.HasPrefix(base, "v") {
		return base
	}
	return ""
}
