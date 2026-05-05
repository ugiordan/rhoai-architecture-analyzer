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
)

// reconcileSearchPaths lists directories to scan for reconciler implementations.
var reconcileSearchPaths = []string{
	"controllers/",
	"internal/controller/",
	"pkg/controller/",
	"pkg/reconciler/",
}

// reconcileMethodPrefixes are method name prefixes that indicate sub-resource reconciliation.
var reconcileMethodPrefixes = []string{
	"Reconcile", "Deploy", "deploy", "Create", "create",
	"ensure", "Ensure", "apply", "Apply", "setup", "Setup",
	"sync", "Sync", "handle", "Handle", "manage", "Manage",
}

// reconcileMethodRE matches method names that start with a reconcile-like prefix
// followed by at least one more uppercase letter (to avoid matching generic methods).
var reconcileMethodRE = regexp.MustCompile(
	`^(?:Reconcile|Deploy|deploy|Create|create|ensure|Ensure|apply|Apply|setup|Setup|sync|Sync|handle|Handle|manage|Manage)[A-Z]\w*$`,
)

// extractReconcileSequences scans Go source for Reconcile() method bodies
// and extracts the ordered sequence of sub-resource reconciliation calls.
func extractReconcileSequences(repoPath string) []ReconcileSequence {
	var goFiles []string
	for _, dir := range reconcileSearchPaths {
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

	var sequences []ReconcileSequence
	for _, fpath := range goFiles {
		seqs := parseReconcileFile(fpath, repoPath)
		sequences = append(sequences, seqs...)
	}

	if sequences == nil {
		sequences = []ReconcileSequence{}
	}
	return sequences
}

// parseReconcileFile parses a Go file for Reconcile methods and extracts call sequences.
func parseReconcileFile(fpath, repoPath string) []ReconcileSequence {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, fpath, nil, parser.ParseComments)
	if err != nil {
		return nil
	}

	relPath := relativePath(repoPath, fpath)
	var sequences []ReconcileSequence

	for _, decl := range node.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}

		// Check if this is a Reconcile method on a controller type
		methodName := funcDecl.Name.Name
		if methodName != "Reconcile" && methodName != "ReconcileAll" && methodName != "reconcile" {
			continue
		}
		if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
			continue
		}

		controllerName := extractReceiverTypeName(funcDecl.Recv.List[0].Type)
		if controllerName == "" {
			continue
		}
		if !strings.HasSuffix(controllerName, "Reconciler") && !strings.HasSuffix(controllerName, "Controller") {
			continue
		}

		// Extract call sequence from function body
		steps := extractStepsFromBody(funcDecl.Body.List, fset, relPath, "")
		if len(steps) == 0 {
			continue
		}

		sequences = append(sequences, ReconcileSequence{
			Controller: controllerName,
			Source:     relPath,
			Steps:      steps,
		})
	}

	return sequences
}

// extractStepsFromBody walks a statement list and extracts reconcile-like calls.
func extractStepsFromBody(stmts []ast.Stmt, fset *token.FileSet, source, parentCondition string) []ReconcileStep {
	var steps []ReconcileStep
	order := 1

	for _, stmt := range stmts {
		switch s := stmt.(type) {
		case *ast.ExprStmt:
			if method := extractReconcileCall(s.X); method != "" {
				steps = append(steps, ReconcileStep{
					Order:     order,
					Method:    method,
					Component: deriveComponent(method),
					Condition: parentCondition,
					Source:    source,
				})
				order++
			}
		case *ast.AssignStmt:
			for _, rhs := range s.Rhs {
				if method := extractReconcileCall(rhs); method != "" {
					steps = append(steps, ReconcileStep{
						Order:     order,
						Method:    method,
						Component: deriveComponent(method),
						Condition: parentCondition,
						Source:    source,
					})
					order++
				}
			}
		case *ast.IfStmt:
			// Extract condition string
			condition := renderExpr(s.Cond, fset)
			if parentCondition != "" {
				condition = parentCondition + " && " + condition
			}
			// Scan the if body (one level deep)
			ifSteps := extractStepsFromBody(s.Body.List, fset, source, condition)
			for _, step := range ifSteps {
				step.Order = order
				steps = append(steps, step)
				order++
			}
			// Scan else body if present
			if s.Else != nil {
				if elseBlock, ok := s.Else.(*ast.BlockStmt); ok {
					negCondition := "!(" + renderExpr(s.Cond, fset) + ")"
					if parentCondition != "" {
						negCondition = parentCondition + " && " + negCondition
					}
					elseSteps := extractStepsFromBody(elseBlock.List, fset, source, negCondition)
					for _, step := range elseSteps {
						step.Order = order
						steps = append(steps, step)
						order++
					}
				}
			}
		}
	}

	return steps
}

// extractReconcileCall checks if an expression is a call to a reconcile-like method.
func extractReconcileCall(expr ast.Expr) string {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return ""
	}

	var methodName string
	switch fn := call.Fun.(type) {
	case *ast.SelectorExpr:
		methodName = fn.Sel.Name
	case *ast.Ident:
		methodName = fn.Name
	default:
		return ""
	}

	if reconcileMethodRE.MatchString(methodName) {
		return methodName
	}
	return ""
}

// deriveComponent strips the verb prefix from a method name to get the component name.
func deriveComponent(method string) string {
	for _, prefix := range reconcileMethodPrefixes {
		if strings.HasPrefix(method, prefix) && len(method) > len(prefix) {
			rest := method[len(prefix):]
			// Ensure the first char after prefix is uppercase (component name, not a common word)
			if len(rest) > 0 && rest[0] >= 'A' && rest[0] <= 'Z' {
				return rest
			}
		}
	}
	return method
}

// extractReceiverTypeName gets the type name from a receiver expression.
func extractReceiverTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return extractReceiverTypeName(t.X)
	case *ast.Ident:
		return t.Name
	}
	return ""
}

// renderExpr renders an AST expression back to source code string.
func renderExpr(expr ast.Expr, fset *token.FileSet) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, expr); err != nil {
		return ""
	}
	result := buf.String()
	if len(result) > 200 {
		return result[:200] + "..."
	}
	return result
}
