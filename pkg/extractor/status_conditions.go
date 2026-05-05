package extractor

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// statusConditionSearchPaths lists directories to scan for status condition definitions.
var statusConditionSearchPaths = []string{
	"controllers/",
	"internal/controller/",
	"pkg/controller/",
	"api/",
	"apis/",
	"pkg/apis/",
	"pkg/api/",
}

// conditionTypeSuffixes matches constants that are status condition types by name suffix.
var conditionTypeSuffixes = []string{"Available", "Ready", "Degraded", "Progressing", "Reconciled"}

// conditionTypeTypes matches constants with explicit condition type identifiers.
var conditionTypeTypes = []string{"ConditionType", "StatusConditionType", "conditionType"}

// conditionReasonTypes matches constants with explicit reason type identifiers.
var conditionReasonTypes = []string{"ConditionReason", "StatusReason", "conditionReason"}

// extractStatusConditions scans Go source for status condition type and reason constants.
// It also returns a set of Go constant names (not values) for dedup with operator config.
func extractStatusConditions(repoPath string) ([]StatusCondition, map[string]bool) {
	var goFiles []string
	for _, dir := range statusConditionSearchPaths {
		fullDir := filepath.Join(repoPath, dir)
		if info, err := os.Stat(fullDir); err != nil || !info.IsDir() {
			continue
		}
		filepath.Walk(fullDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() && strings.Contains(path, "/vendor") {
				return filepath.SkipDir
			}
			if info.IsDir() {
				return nil
			}
			if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
				goFiles = append(goFiles, path)
			}
			return nil
		})
	}

	var conditions []StatusCondition
	seen := make(map[string]bool)
	constNames := make(map[string]bool)

	for _, fpath := range goFiles {
		parsed, names := parseStatusConditionsFile(fpath, repoPath)
		for _, c := range parsed {
			if c.Type != "" && !seen[c.Type] {
				seen[c.Type] = true
				conditions = append(conditions, c)
			}
		}
		for name := range names {
			constNames[name] = true
		}
	}

	if conditions == nil {
		conditions = []StatusCondition{}
	}
	return conditions, constNames
}

// parseStatusConditionsFile parses a single Go file for status condition constants.
// Returns the conditions and a set of Go constant names matched.
func parseStatusConditionsFile(fpath, repoPath string) ([]StatusCondition, map[string]bool) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, fpath, nil, parser.ParseComments)
	if err != nil {
		return nil, nil
	}

	relPath := relativePath(repoPath, fpath)
	var conditions []StatusCondition
	constNames := make(map[string]bool)

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}

		// Process each const block: track current condition type for reason association.
		var currentCondType string
		var currentReasons []string

		for _, spec := range genDecl.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok || len(vs.Names) == 0 {
				continue
			}
			name := vs.Names[0].Name
			value := extractConstStringValue(vs)

			isType := isConditionType(name, vs)
			isReason := isConditionReason(name, vs)

			if isType {
				constNames[name] = true
				// Flush previous condition if any
				if currentCondType != "" {
					conditions = append(conditions, StatusCondition{
						Type:    currentCondType,
						Reasons: currentReasons,
						Source:  relPath,
					})
				}
				if value != "" {
					currentCondType = value
				} else {
					currentCondType = name
				}
				currentReasons = nil
			} else if isReason {
				constNames[name] = true
				if currentCondType != "" {
					if value != "" {
						currentReasons = append(currentReasons, value)
					} else {
						currentReasons = append(currentReasons, name)
					}
				}
			}
		}

		// Flush last condition in block
		if currentCondType != "" {
			conditions = append(conditions, StatusCondition{
				Type:    currentCondType,
				Reasons: currentReasons,
				Source:  relPath,
			})
		}
	}

	return conditions, constNames
}

// isConditionType checks if a constant is a status condition type.
func isConditionType(name string, vs *ast.ValueSpec) bool {
	// Check explicit type annotation
	if vs.Type != nil {
		typeName := typeIdentName(vs.Type)
		for _, t := range conditionTypeTypes {
			if typeName == t || strings.HasSuffix(typeName, "."+t) {
				return true
			}
		}
	}
	// Check name suffix
	for _, suffix := range conditionTypeSuffixes {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}
	// Check name prefix
	if strings.HasPrefix(name, "Condition") {
		return true
	}
	return false
}

// isConditionReason checks if a constant is a status condition reason.
func isConditionReason(name string, vs *ast.ValueSpec) bool {
	// Check explicit type annotation
	if vs.Type != nil {
		typeName := typeIdentName(vs.Type)
		for _, t := range conditionReasonTypes {
			if typeName == t || strings.HasSuffix(typeName, "."+t) {
				return true
			}
		}
	}
	// Check name prefix
	return strings.HasPrefix(name, "Reason")
}

// typeIdentName extracts the type name string from an ast.Expr.
func typeIdentName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name + "." + t.Sel.Name
		}
		return t.Sel.Name
	}
	return ""
}

// extractConstStringValue extracts the string value of a const declaration.
func extractConstStringValue(vs *ast.ValueSpec) string {
	if len(vs.Values) == 0 {
		return ""
	}
	lit, ok := vs.Values[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return ""
	}
	return strings.Trim(lit.Value, `"`)
}

