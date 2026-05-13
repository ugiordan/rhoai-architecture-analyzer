package extractor

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"unicode"
)

// WebhookBehavior captures the mutations and validations performed by a webhook handler.
type WebhookBehavior struct {
	TargetType  string    `json:"target_type"`
	Mutations   []FieldOp `json:"mutations,omitempty"`
	Validations []FieldOp `json:"validations,omitempty"`
}

// extractWebhookBehavior scans all packages for kubebuilder webhook markers and
// analyzes the corresponding Default/Validate* methods to extract field-level
// mutations and validations.
func extractWebhookBehavior(pkgs *GoPackageSet) map[string]WebhookBehavior {
	if pkgs == nil {
		return nil
	}
	result := make(map[string]WebhookBehavior)

	for _, pkg := range pkgs.Packages {
		for _, file := range pkg.Syntax {
			// Scan all comments in the file for webhook markers.
			webhooks := findWebhookMarkers(file)
			if len(webhooks) == 0 {
				continue
			}

			for _, wh := range webhooks {
				behavior := WebhookBehavior{}

				// Determine the target type by finding the closest method
				// declaration after the webhook marker. Kubebuilder convention
				// places the marker immediately before the annotated method.
				typeName := findReceiverTypeNearMarker(file, wh.pos, pkg.Fset)

				if typeName == "" {
					continue
				}
				behavior.TargetType = typeName

				// Get all methods for this type across the package
				allMethods := FindMethodsOnType(pkg, typeName)
				methodMap := make(map[string]*ast.FuncDecl)
				for _, m := range allMethods {
					methodMap[m.Name.Name] = m
				}

				if wh.mutating {
					behavior.Mutations = extractMutations(methodMap, pkg.Fset)
				} else {
					behavior.Validations = extractValidations(methodMap, pkg.Fset)
				}

				behavior.Mutations = dedupeFieldOps(behavior.Mutations)
				behavior.Validations = dedupeFieldOps(behavior.Validations)

				result[wh.path] = behavior
			}
		}
	}
	return result
}

// webhookMarkerInfo holds parsed info from a +kubebuilder:webhook comment.
type webhookMarkerInfo struct {
	path     string
	mutating bool
	pos      token.Pos
}

// findWebhookMarkers scans all comments in a file for +kubebuilder:webhook markers.
func findWebhookMarkers(file *ast.File) []webhookMarkerInfo {
	var result []webhookMarkerInfo

	for _, cg := range file.Comments {
		for _, c := range cg.List {
			text := c.Text
			// Strip leading // or /* */
			if strings.HasPrefix(text, "//") {
				text = strings.TrimPrefix(text, "//")
			} else if strings.HasPrefix(text, "/*") {
				text = strings.TrimPrefix(text, "/*")
				text = strings.TrimSuffix(text, "*/")
			}
			text = strings.TrimSpace(text)

			if !strings.HasPrefix(text, "+kubebuilder:webhook:") {
				continue
			}

			info := parseWebhookMarker(text)
			if info.path != "" {
				info.pos = c.Pos()
				result = append(result, info)
			}
		}
	}
	return result
}

// parseWebhookMarker parses a +kubebuilder:webhook:key=val,... marker line.
func parseWebhookMarker(marker string) webhookMarkerInfo {
	info := webhookMarkerInfo{}
	// Remove the prefix
	content := strings.TrimPrefix(marker, "+kubebuilder:webhook:")

	for _, part := range strings.Split(content, ",") {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])
		switch key {
		case "path":
			info.path = val
		case "mutating":
			info.mutating = val == "true"
		}
	}
	return info
}

// findReceiverTypeNearMarker finds the receiver type of the closest method
// declaration after the given position. Kubebuilder places webhook markers
// immediately before the method they annotate.
func findReceiverTypeNearMarker(file *ast.File, markerPos token.Pos, fset *token.FileSet) string {
	var closest *ast.FuncDecl
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv == nil || len(fn.Recv.List) == 0 {
			continue
		}
		if fn.Pos() < markerPos {
			continue
		}
		if closest == nil || fn.Pos() < closest.Pos() {
			closest = fn
		}
	}
	if closest == nil {
		return ""
	}
	return receiverTypeName(closest.Recv.List[0].Type)
}

// extractMutations analyzes the Default() method and any helper methods it
// calls on the same receiver to find field assignments (mutations).
func extractMutations(methods map[string]*ast.FuncDecl, fset *token.FileSet) []FieldOp {
	defaultFn, ok := methods["Default"]
	if !ok {
		return nil
	}

	var ops []FieldOp
	visited := make(map[string]bool)
	ops = collectMutationsFromFunc(defaultFn, methods, fset, visited)
	return ops
}

// collectMutationsFromFunc walks a function body looking for field assignments
// and follows same-receiver method calls.
func collectMutationsFromFunc(fn *ast.FuncDecl, methods map[string]*ast.FuncDecl, fset *token.FileSet, visited map[string]bool) []FieldOp {
	if fn == nil || fn.Body == nil {
		return nil
	}
	name := fn.Name.Name
	if visited[name] {
		return nil
	}
	visited[name] = true

	var ops []FieldOp
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		switch stmt := n.(type) {
		case *ast.AssignStmt:
			// Look for field assignments like w.Spec.Image = "..."
			enclosingCond := findEnclosingCondition(fn.Body, stmt, fset)
			for _, lhs := range stmt.Lhs {
				path := selectorPath(lhs)
				if path == "" {
					continue
				}
				jsonPath := goPathToJSON(path)
				if jsonPath == "" {
					continue
				}
				pos := fset.Position(stmt.Pos())
				ops = append(ops, FieldOp{
					Field:     jsonPath,
					Operation: "set",
					Condition: enclosingCond,
					Source: SourceRef{
						Type: "go_handler",
						File: pos.Filename,
						Line: pos.Line,
					},
				})
			}
		case *ast.ExprStmt:
			// Look for same-receiver method calls like r.setGPUDefaults(w)
			call, ok := stmt.X.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			calledName := sel.Sel.Name
			if calledFn, exists := methods[calledName]; exists {
				subOps := collectMutationsFromFunc(calledFn, methods, fset, visited)
				ops = append(ops, subOps...)
			}
		}
		return true
	})
	return ops
}

// extractValidations analyzes ValidateCreate/ValidateUpdate/ValidateDelete methods
// to find field.Invalid/Required/Forbidden/Duplicate calls.
func extractValidations(methods map[string]*ast.FuncDecl, fset *token.FileSet) []FieldOp {
	var ops []FieldOp
	visited := make(map[string]bool)

	for _, name := range []string{"ValidateCreate", "ValidateUpdate", "ValidateDelete"} {
		fn, ok := methods[name]
		if !ok {
			continue
		}
		subOps := collectValidationsFromFunc(fn, methods, fset, visited)
		ops = append(ops, subOps...)
	}
	return ops
}

// collectValidationsFromFunc walks a function body looking for field.Invalid/Required/etc.
// calls and follows same-receiver method calls.
func collectValidationsFromFunc(fn *ast.FuncDecl, methods map[string]*ast.FuncDecl, fset *token.FileSet, visited map[string]bool) []FieldOp {
	if fn == nil || fn.Body == nil {
		return nil
	}
	name := fn.Name.Name
	if visited[name] {
		return nil
	}
	visited[name] = true

	var ops []FieldOp
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Check for field.Invalid, field.Required, field.Forbidden, field.Duplicate
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		opName := ""
		switch sel.Sel.Name {
		case "Invalid":
			opName = "invalid"
		case "Required":
			opName = "required"
		case "Forbidden":
			opName = "forbidden"
		case "Duplicate":
			opName = "duplicate"
		default:
			// Check if this is a same-receiver method call
			calledName := sel.Sel.Name
			if calledFn, exists := methods[calledName]; exists && !visited[calledName] {
				subOps := collectValidationsFromFunc(calledFn, methods, fset, visited)
				ops = append(ops, subOps...)
			}
			return true
		}

		// Verify caller is "field" package
		if ident, ok := sel.X.(*ast.Ident); ok {
			if ident.Name != "field" {
				return true
			}
		} else {
			return true
		}

		// First arg should be field.NewPath(...)
		if len(call.Args) < 1 {
			return true
		}
		fieldPath := extractFieldNewPath(call.Args[0])
		if fieldPath == "" {
			return true
		}

		pos := fset.Position(call.Pos())
		ops = append(ops, FieldOp{
			Field:     fieldPath,
			Operation: opName,
			Source: SourceRef{
				Type: "go_handler",
				File: pos.Filename,
				Line: pos.Line,
			},
		})
		return true
	})
	return ops
}

// extractFieldNewPath extracts the field path from a field.NewPath("spec", "replicas") call.
func extractFieldNewPath(expr ast.Expr) string {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return ""
	}

	// Handle field.NewPath("spec", "replicas") or field.NewPath("spec").Child("replicas")
	// First, check if it's a method call chain: NewPath(...).Child(...)
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if sel.Sel.Name == "Child" {
			// Recursively get the base path
			basePath := extractFieldNewPath(sel.X)
			childSegments := collectPathSegments(call.Args)
			if basePath != "" && len(childSegments) > 0 {
				return basePath + "." + strings.Join(childSegments, ".")
			}
			if basePath != "" {
				return basePath
			}
			return strings.Join(childSegments, ".")
		}

		// Check for field.NewPath(...)
		if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "field" && sel.Sel.Name == "NewPath" {
			segments := collectPathSegments(call.Args)
			return strings.Join(segments, ".")
		}
	}

	return ""
}

// collectPathSegments extracts string literal arguments from a function call.
func collectPathSegments(args []ast.Expr) []string {
	var segments []string
	for _, arg := range args {
		lit, ok := arg.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			continue
		}
		// Remove quotes
		val := strings.Trim(lit.Value, `"`)
		segments = append(segments, val)
	}
	return segments
}

// selectorPath builds a dotted path from a selector expression.
// e.g., w.Spec.Image -> "w.Spec.Image"
func selectorPath(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.SelectorExpr:
		base := selectorPath(e.X)
		if base == "" {
			return ""
		}
		return base + "." + e.Sel.Name
	case *ast.Ident:
		return e.Name
	default:
		return ""
	}
}

// goPathToJSON converts a Go field path like "w.Spec.Image" to a JSON path like "spec.image".
// It strips the first segment (variable name) and lowercases the first character of each remaining segment.
func goPathToJSON(path string) string {
	parts := strings.Split(path, ".")
	if len(parts) < 2 {
		return ""
	}
	// Skip the first segment (variable name like "w" or "r")
	jsonParts := make([]string, 0, len(parts)-1)
	for _, p := range parts[1:] {
		jsonParts = append(jsonParts, camelToJSON(p))
	}
	return strings.Join(jsonParts, ".")
}

// camelToJSON converts a Go field name to its JSON equivalent.
// Handles Go's acronym convention: APIVersion -> apiVersion, TLSConfig -> tlsConfig, GPU -> gpu.
func camelToJSON(s string) string {
	if s == "" {
		return s
	}
	if strings.ToUpper(s) == s {
		return strings.ToLower(s)
	}
	runes := []rune(s)
	i := 0
	for i < len(runes) && unicode.IsUpper(runes[i]) {
		i++
	}
	if i == 0 {
		return s
	}
	if i == 1 {
		runes[0] = unicode.ToLower(runes[0])
		return string(runes)
	}
	for j := 0; j < i-1; j++ {
		runes[j] = unicode.ToLower(runes[j])
	}
	return string(runes)
}

// findEnclosingCondition walks the AST to find the innermost if statement
// that encloses the target node and returns a human-readable condition string.
func findEnclosingCondition(body *ast.BlockStmt, target ast.Node, fset *token.FileSet) string {
	if body == nil {
		return ""
	}
	return findConditionInBlock(body.List, target)
}

func findConditionInBlock(stmts []ast.Stmt, target ast.Node) string {
	targetPos := target.Pos()
	targetEnd := target.End()

	for _, stmt := range stmts {
		switch s := stmt.(type) {
		case *ast.IfStmt:
			if targetPos >= s.Body.Pos() && targetEnd <= s.Body.End() {
				inner := findConditionInBlock(s.Body.List, target)
				cond := formatExpr(s.Cond)
				if inner != "" {
					return cond + " && " + inner
				}
				return cond
			}
			if s.Else != nil {
				if elseBlock, ok := s.Else.(*ast.BlockStmt); ok {
					if targetPos >= elseBlock.Pos() && targetEnd <= elseBlock.End() {
						inner := findConditionInBlock(elseBlock.List, target)
						cond := "!(" + formatExpr(s.Cond) + ")"
						if inner != "" {
							return cond + " && " + inner
						}
						return cond
					}
				}
				if elseIf, ok := s.Else.(*ast.IfStmt); ok {
					if targetPos >= elseIf.Pos() && targetEnd <= elseIf.End() {
						inner := findConditionInBlock([]ast.Stmt{elseIf}, target)
						negated := "!(" + formatExpr(s.Cond) + ")"
						if inner != "" {
							return negated + " && " + inner
						}
						return negated
					}
				}
			}
		case *ast.ForStmt:
			if s.Body != nil && targetPos >= s.Body.Pos() && targetEnd <= s.Body.End() {
				return findConditionInBlock(s.Body.List, target)
			}
		case *ast.RangeStmt:
			if s.Body != nil && targetPos >= s.Body.Pos() && targetEnd <= s.Body.End() {
				return findConditionInBlock(s.Body.List, target)
			}
		case *ast.SwitchStmt:
			if s.Body != nil && targetPos >= s.Body.Pos() && targetEnd <= s.Body.End() {
				for _, cc := range s.Body.List {
					clause, ok := cc.(*ast.CaseClause)
					if !ok {
						continue
					}
					for _, bodyStmt := range clause.Body {
						if targetPos >= bodyStmt.Pos() && targetEnd <= bodyStmt.End() {
							caseLabel := formatCaseValues(clause.List)
							inner := findConditionInBlock(clause.Body, target)
							if inner != "" {
								return caseLabel + " && " + inner
							}
							return caseLabel
						}
					}
				}
			}
		case *ast.TypeSwitchStmt:
			if s.Body != nil && targetPos >= s.Body.Pos() && targetEnd <= s.Body.End() {
				for _, cc := range s.Body.List {
					clause, ok := cc.(*ast.CaseClause)
					if !ok {
						continue
					}
					for _, bodyStmt := range clause.Body {
						if targetPos >= bodyStmt.Pos() && targetEnd <= bodyStmt.End() {
							caseLabel := formatCaseValues(clause.List)
							inner := findConditionInBlock(clause.Body, target)
							if inner != "" {
								return caseLabel + " && " + inner
							}
							return caseLabel
						}
					}
				}
			}
		case *ast.BlockStmt:
			if targetPos >= s.Pos() && targetEnd <= s.End() {
				return findConditionInBlock(s.List, target)
			}
		}
	}
	return ""
}

func formatCaseValues(exprs []ast.Expr) string {
	if len(exprs) == 0 {
		return "default"
	}
	parts := make([]string, len(exprs))
	for i, e := range exprs {
		parts[i] = formatExpr(e)
	}
	return strings.Join(parts, ", ")
}

// formatExpr renders an ast.Expr as a readable string.
func formatExpr(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		return formatExpr(e.X) + " " + e.Op.String() + " " + formatExpr(e.Y)
	case *ast.UnaryExpr:
		return e.Op.String() + formatExpr(e.X)
	case *ast.SelectorExpr:
		return formatExpr(e.X) + "." + e.Sel.Name
	case *ast.Ident:
		return e.Name
	case *ast.BasicLit:
		return e.Value
	case *ast.CallExpr:
		return formatExpr(e.Fun) + "(...)"
	case *ast.ParenExpr:
		return "(" + formatExpr(e.X) + ")"
	case *ast.StarExpr:
		return "*" + formatExpr(e.X)
	default:
		return fmt.Sprintf("%T", expr)
	}
}

// dedupeFieldOps removes duplicate FieldOp entries by field+operation key.
func dedupeFieldOps(ops []FieldOp) []FieldOp {
	if len(ops) == 0 {
		return nil
	}
	seen := make(map[string]bool)
	var result []FieldOp
	for _, op := range ops {
		key := op.Field + "|" + op.Operation
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, op)
	}
	return result
}
