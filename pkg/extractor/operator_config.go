package extractor

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// operatorConfigSearchPaths lists directories to scan for operator configuration constants.
var operatorConfigSearchPaths = []string{
	"controllers/",
	"internal/controller/",
	"pkg/config/",
	"config/",
}

// imageRegistryPrefixes identifies container image references.
var imageRegistryPrefixes = []string{
	"quay.io/", "registry.redhat.io/", "gcr.io/", "docker.io/",
	"ghcr.io/", "registry.access.redhat.com/", "registry.k8s.io/",
}

// upperSnakeCaseRE matches UPPER_SNAKE_CASE names (3+ chars).
var upperSnakeCaseRE = regexp.MustCompile(`^[A-Z][A-Z0-9_]{2,}$`)

// extractOperatorConfig scans Go source for const/var blocks that define operator
// configuration. statusConditionNames is the dedup set from the status conditions
// extractor; constants matching those names are skipped.
func extractOperatorConfig(repoPath string, statusConditionNames map[string]bool) []OperatorConstant {
	var goFiles []string

	// Scan configured directories
	for _, dir := range operatorConfigSearchPaths {
		fullDir := filepath.Join(repoPath, dir)
		if info, err := os.Stat(fullDir); err != nil || !info.IsDir() {
			continue
		}
		filepath.Walk(fullDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") &&
				!strings.Contains(path, "/vendor/") {
				goFiles = append(goFiles, path)
			}
			return nil
		})
	}

	// Also scan root-level .go files
	entries, _ := os.ReadDir(repoPath)
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".go") && !strings.HasSuffix(e.Name(), "_test.go") {
			goFiles = append(goFiles, filepath.Join(repoPath, e.Name()))
		}
	}

	if statusConditionNames == nil {
		statusConditionNames = map[string]bool{}
	}

	var constants []OperatorConstant
	seen := make(map[string]bool)

	for _, fpath := range goFiles {
		parsed := parseOperatorConfigFile(fpath, repoPath, statusConditionNames)
		for _, c := range parsed {
			key := c.Name + ":" + c.Source
			if seen[key] {
				continue
			}
			seen[key] = true
			constants = append(constants, c)
		}
	}

	if constants == nil {
		constants = []OperatorConstant{}
	}
	return constants
}

// parseOperatorConfigFile parses a single Go file for const/var declarations.
func parseOperatorConfigFile(fpath, repoPath string, dedupNames map[string]bool) []OperatorConstant {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, fpath, nil, parser.ParseComments)
	if err != nil {
		return nil
	}

	relPath := relativePath(repoPath, fpath)
	var constants []OperatorConstant

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || (genDecl.Tok != token.CONST && genDecl.Tok != token.VAR) {
			continue
		}

		// Check if this is an iota-only enum block
		if isIotaOnlyBlock(genDecl) {
			continue
		}

		for _, spec := range genDecl.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok || len(vs.Names) == 0 {
				continue
			}
			name := vs.Names[0].Name

			// Skip single-char unexported names
			if len(name) == 1 && unicode.IsLower(rune(name[0])) {
				continue
			}

			// Skip if captured by status conditions
			if dedupNames[name] {
				continue
			}

			// Extract value
			value := extractConstValue(vs, fset)
			if value == "" {
				continue
			}

			// Extract type
			goType := ""
			if vs.Type != nil {
				goType = typeIdentName(vs.Type)
			}

			// Extract doc comment
			doc := ""
			if vs.Doc != nil {
				doc = strings.TrimSpace(vs.Doc.Text())
			} else if genDecl.Doc != nil && len(genDecl.Specs) == 1 {
				doc = strings.TrimSpace(genDecl.Doc.Text())
			}

			category := classifyConstant(name, value, goType)

			constants = append(constants, OperatorConstant{
				Name:     name,
				Value:    value,
				GoType:   goType,
				Category: category,
				Doc:      doc,
				Source:   relPath,
			})
		}
	}

	return constants
}

// classifyConstant determines the category of a constant using precedence order.
func classifyConstant(name, value, goType string) string {
	// Priority 1: image
	for _, prefix := range imageRegistryPrefixes {
		if strings.Contains(value, prefix) {
			return "image"
		}
	}
	if strings.Contains(name, "Image") {
		return "image"
	}

	// Priority 2: port
	if strings.Contains(name, "Port") && isNumericString(value) {
		return "port"
	}

	// Priority 3: timeout
	if goType == "time.Duration" || goType == "Duration" {
		return "timeout"
	}
	for _, kw := range []string{"Timeout", "Interval", "Requeue", "Expiry"} {
		if strings.Contains(name, kw) {
			return "timeout"
		}
	}

	// Priority 4: env_var
	if upperSnakeCaseRE.MatchString(name) {
		return "env_var"
	}

	// Priority 5: resource
	for _, kw := range []string{"Size", "Memory", "CPU"} {
		if strings.Contains(name, kw) {
			return "resource"
		}
	}
	if isK8sQuantity(value) {
		return "resource"
	}

	// Priority 6: name_pattern
	if strings.HasSuffix(value, "-") {
		return "name_pattern"
	}

	// Priority 7: general
	return "general"
}

// isIotaOnlyBlock checks if a const block only contains iota values (no string/numeric).
func isIotaOnlyBlock(genDecl *ast.GenDecl) bool {
	if genDecl.Tok != token.CONST {
		return false
	}
	hasIota := false
	hasValue := false
	for _, spec := range genDecl.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		for _, v := range vs.Values {
			if ident, ok := v.(*ast.Ident); ok && ident.Name == "iota" {
				hasIota = true
			} else if _, ok := v.(*ast.BasicLit); ok {
				hasValue = true
			} else if _, ok := v.(*ast.BinaryExpr); ok {
				// Could be iota + 1, etc.
				hasValue = true
			}
		}
		if len(vs.Values) == 0 && hasIota {
			// Implicit iota continuation
			continue
		}
	}
	return hasIota && !hasValue
}

// extractConstValue extracts the string representation of a const/var value.
func extractConstValue(vs *ast.ValueSpec, fset *token.FileSet) string {
	if len(vs.Values) == 0 {
		return ""
	}
	val := vs.Values[0]

	// BasicLit: string, int, float
	if lit, ok := val.(*ast.BasicLit); ok {
		switch lit.Kind {
		case token.STRING:
			return strings.Trim(lit.Value, `"`+"`")
		case token.INT, token.FLOAT:
			return lit.Value
		}
	}

	// For other expressions, use go/printer to render source
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, val); err == nil {
		rendered := buf.String()
		if len(rendered) < 200 {
			return rendered
		}
	}

	return ""
}

// isNumericString checks if a string is a numeric value.
func isNumericString(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// isK8sQuantity checks if a string looks like a Kubernetes resource quantity.
func isK8sQuantity(s string) bool {
	if s == "" {
		return false
	}
	suffixes := []string{"Gi", "Mi", "Ki", "Ti", "m", "k", "M", "G", "T"}
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			prefix := strings.TrimSuffix(s, suffix)
			return isNumericString(prefix)
		}
	}
	return false
}
