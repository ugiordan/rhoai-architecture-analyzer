package parser

import (
	"context"
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/rust"

	"github.com/ugiordan/architecture-analyzer/pkg/dataflow"
	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

// RustSkipDirs lists directories that should be skipped when scanning Rust projects.
var RustSkipDirs = []string{"target"}

// RustTestPatterns lists filename patterns that identify test files.
var RustTestPatterns = []string{"*_test.rs", "test_*"}

// rustHTTPMethods maps attribute names to HTTP methods for actix-web/rocket-style route handlers.
var rustHTTPMethods = map[string]bool{
	"get": true, "post": true, "put": true,
	"delete": true, "patch": true, "head": true,
	"options": true, "route": true,
}

// rustDBPatterns maps Rust DB call/macro patterns to read/write classification.
var rustDBPatterns = map[string]string{
	"diesel::insert_into":   "write",
	"diesel::update":        "write",
	"diesel::delete":        "write",
	"sqlx::query":           "read",
	"sqlx::query_as":        "read",
	"sqlx::query_scalar":    "read",
}

// RustParser extracts code property graph nodes from Rust source files using tree-sitter.
// Each goroutine MUST use its own RustParser instance (tree-sitter parsers are not thread-safe).
type RustParser struct {
	parser *sitter.Parser
}

// NewRustParser creates a parser for Rust source files backed by tree-sitter.
func NewRustParser() *RustParser {
	p := sitter.NewParser()
	p.SetLanguage(rust.GetLanguage())
	return &RustParser{parser: p}
}

func (rp *RustParser) Language() string     { return "rust" }
func (rp *RustParser) Extensions() []string { return []string{".rs"} }
func (rp *RustParser) Clone() Parser {
	p := sitter.NewParser()
	p.SetLanguage(rust.GetLanguage())
	return &RustParser{parser: p}
}

// computeRustComplexity counts decision points in a Rust function body.
// Complexity = 1 (base) + count of: if, match arm, for, while, loop, &&, ? operator.
func computeRustComplexity(node *sitter.Node) int {
	count := 1
	countRustDecisionPoints(node, &count)
	return count
}

func countRustDecisionPoints(node *sitter.Node, count *int) {
	if node == nil {
		return
	}
	switch node.Type() {
	case "if_expression":
		// Skip if this if_expression is the direct child of an else_clause
		// (i.e. "else if"), because it is already counted by the parent if.
		parent := node.Parent()
		if parent == nil || parent.Type() != "else_clause" {
			*count++
		}
	case "for_expression":
		*count++
	case "while_expression":
		*count++
	case "loop_expression":
		*count++
	case "match_arm":
		*count++
	case "try_expression":
		*count++ // the ? operator
	case "binary_expression":
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child != nil && child.Type() == "&&" {
				*count++
				break
			}
		}
	}
	for i := 0; i < int(node.ChildCount()); i++ {
		countRustDecisionPoints(node.Child(i), count)
	}
}

// ParseFile parses a Rust source file and returns extracted nodes and edges.
func (rp *RustParser) ParseFile(path string, content []byte) (*ParseResult, error) {
	if len(content) > MaxFileSize {
		return nil, fmt.Errorf("file too large (%d bytes, max %d)", len(content), MaxFileSize)
	}
	tree, err := rp.parser.ParseCtx(context.Background(), nil, content)
	if err != nil {
		return nil, fmt.Errorf("tree-sitter parse failed: %w", err)
	}
	defer tree.Close()

	result := &ParseResult{}
	root := tree.RootNode()
	rp.walk(root, content, path, "", result)
	return result, nil
}

// walk recursively traverses the AST. implType tracks the enclosing impl target type.
func (rp *RustParser) walk(node *sitter.Node, src []byte, file, implType string, result *ParseResult) {
	switch node.Type() {
	case "function_item":
		rp.extractFunction(node, src, file, implType, result)
	case "call_expression":
		rp.extractCallSite(node, src, file, result)
	case "macro_invocation":
		rp.extractMacroInvocation(node, src, file, result)
	case "struct_expression":
		rp.extractStructExpression(node, src, file, result)
	case "impl_item":
		rp.extractImpl(node, src, file, result)
		return // children handled inside extractImpl
	}

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil {
			rp.walk(child, src, file, implType, result)
		}
	}
}

// extractImpl processes an impl_item, walking its body with the impl target type set.
func (rp *RustParser) extractImpl(node *sitter.Node, src []byte, file string, result *ParseResult) {
	typeNode := node.ChildByFieldName("type")
	typeName := ""
	if typeNode != nil {
		typeName = typeNode.Content(src)
	}

	body := node.ChildByFieldName("body")
	if body == nil {
		return
	}
	for i := 0; i < int(body.ChildCount()); i++ {
		child := body.Child(i)
		if child != nil {
			rp.walk(child, src, file, typeName, result)
		}
	}
}

// extractFunction creates a Function node from a function_item.
// It also checks preceding siblings for attribute_items to detect HTTP handlers and test markers.
func (rp *RustParser) extractFunction(node *sitter.Node, src []byte, file, implType string, result *ParseResult) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}
	name := nameNode.Content(src)
	line := int(node.StartPoint().Row) + 1
	col := int(node.StartPoint().Column)

	fn := &graph.Node{
		ID:          NodeID(graph.NodeFunction, name, file, line, col),
		Kind:        graph.NodeFunction,
		Name:        name,
		File:        file,
		Line:        line,
		Column:      col,
		EndLine:     int(node.EndPoint().Row) + 1,
		Language:    "rust",
		TypeName:    implType,
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}

	// Check function modifiers for unsafe/extern
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}
		if child.Type() == "function_modifiers" {
			modText := child.Content(src)
			if strings.Contains(modText, "unsafe") {
				fn.Properties["is_unsafe"] = "true"
				fn.IsUnsafe = true
			}
			if strings.Contains(modText, "extern") {
				fn.Properties["is_extern"] = "true"
				fn.IsExtern = true
			}
		}
	}

	// Extract parameter names from the parameters node
	params := node.ChildByFieldName("parameters")
	if params != nil {
		fn.ParamNames = extractRustParamNames(params, src)
	}

	// Compute cyclomatic complexity from function body.
	body := node.ChildByFieldName("body")
	if body != nil {
		fn.Complexity = computeRustComplexity(body)
	}

	// Extract intraprocedural data flow from function body
	if body != nil {
		rp.analyzeFunctionBody(body, src, file, fn, result)
	}

	// Check preceding siblings for attribute_items (HTTP routes, #[test], etc.)
	rp.checkPrecedingAttributes(node, src, fn, file, result)

	result.Functions = append(result.Functions, fn)
}

// checkPrecedingAttributes looks at preceding siblings of a function_item for attribute_items.
// It detects HTTP route attributes (#[get("/path")], #[post("/path")], etc.) and #[test].
func (rp *RustParser) checkPrecedingAttributes(fnNode *sitter.Node, src []byte, fn *graph.Node, file string, result *ParseResult) {
	parent := fnNode.Parent()
	if parent == nil {
		return
	}

	// Find this node's index in parent
	nodeIdx := -1
	for i := 0; i < int(parent.ChildCount()); i++ {
		if parent.Child(i) == fnNode {
			nodeIdx = i
			break
		}
	}
	if nodeIdx < 0 {
		return
	}

	// Walk backwards from the function to collect preceding attribute_items
	for i := nodeIdx - 1; i >= 0; i-- {
		sibling := parent.Child(i)
		if sibling == nil || sibling.Type() != "attribute_item" {
			break
		}

		attrText := sibling.Content(src)
		fn.Decorators = append(fn.Decorators, attrText)

		// Extract the attribute content: find the attribute child
		for j := 0; j < int(sibling.ChildCount()); j++ {
			attrNode := sibling.Child(j)
			if attrNode == nil || attrNode.Type() != "attribute" {
				continue
			}

			// The attribute has an identifier (method) and optionally a token_tree (args)
			attrName := ""
			var tokenTree *sitter.Node
			for k := 0; k < int(attrNode.ChildCount()); k++ {
				child := attrNode.Child(k)
				if child == nil {
					continue
				}
				if child.Type() == "identifier" {
					attrName = child.Content(src)
				}
				if child.Type() == "token_tree" {
					tokenTree = child
				}
			}

			// Check for #[test]
			if attrName == "test" {
				fn.Properties["is_test"] = "true"
				fn.IsTest = true
			}

			// Check for HTTP route attributes
			if rustHTTPMethods[attrName] && tokenTree != nil {
				route := rp.extractRouteFromTokenTree(tokenTree, src)
				method := strings.ToUpper(attrName)
				handler := &graph.Node{
					ID:          NodeID(graph.NodeHTTPEndpoint, fn.Name, file, fn.Line, fn.Column),
					Kind:        graph.NodeHTTPEndpoint,
					Name:        fn.Name,
					File:        file,
					Line:        fn.Line,
					Column:      fn.Column,
					Language:    "rust",
					Annotations: make(map[string]bool),
					Properties:  make(map[string]string),
					Route:       route,
					HTTPMethod:  method,
				}
				if route != "" {
					handler.Properties["route"] = route
				}
				handler.Properties["method"] = method
				result.HTTPHandlers = append(result.HTTPHandlers, handler)
			}
		}
	}
}

// extractRouteFromTokenTree extracts the route path string from a token_tree like ("/health").
func (rp *RustParser) extractRouteFromTokenTree(tokenTree *sitter.Node, src []byte) string {
	for i := 0; i < int(tokenTree.ChildCount()); i++ {
		child := tokenTree.Child(i)
		if child == nil {
			continue
		}
		if child.Type() == "string_literal" {
			// Extract the content between quotes
			for j := 0; j < int(child.ChildCount()); j++ {
				gc := child.Child(j)
				if gc != nil && gc.Type() == "string_content" {
					return gc.Content(src)
				}
			}
			// Fallback: strip quotes
			text := child.Content(src)
			return strings.Trim(text, "\"")
		}
	}
	return ""
}

// extractCallSite creates a CallSite node from a call_expression and detects DB operations.
func (rp *RustParser) extractCallSite(node *sitter.Node, src []byte, file string, result *ParseResult) {
	fnNode := node.ChildByFieldName("function")
	if fnNode == nil {
		return
	}
	callText := fnNode.Content(src)
	line := int(node.StartPoint().Row) + 1
	col := int(node.StartPoint().Column)

	cs := &graph.Node{
		ID:         NodeID(graph.NodeCallSite, callText, file, line, col),
		Kind:       graph.NodeCallSite,
		Name:       callText,
		File:       file,
		Line:       line,
		Column:     col,
		Language:   "rust",
		Properties: make(map[string]string),
	}
	result.CallSites = append(result.CallSites, cs)

	// Extract string arguments for secret detection
	args := node.ChildByFieldName("arguments")
	if args != nil {
		var stringArgs []string
		for i := 0; i < int(args.ChildCount()); i++ {
			arg := args.Child(i)
			if arg == nil {
				continue
			}
			if arg.Type() == "string_literal" {
				for j := 0; j < int(arg.ChildCount()); j++ {
					ch := arg.Child(j)
					if ch != nil && ch.Type() == "string_content" {
						stringArgs = append(stringArgs, ch.Content(src))
					}
				}
			}
		}
		if len(stringArgs) > 0 {
			cs.Properties["string_args"] = strings.Join(stringArgs, ",")
		}
	}

	// Check for DB operations by looking at the full call chain text
	rp.maybeExtractDBFromCall(callText, node, src, file, line, col, result)
}

// maybeExtractDBFromCall checks if a call expression's direct function name matches a known DB pattern.
// Only matches on the innermost call to avoid duplicates from chained method calls.
func (rp *RustParser) maybeExtractDBFromCall(callText string, _ *sitter.Node, _ []byte, file string, line, col int, result *ParseResult) {
	// Check if the direct call name matches a DB pattern exactly.
	// callText is the function field content. For chained calls like
	// diesel::insert_into(items::table).values(...).execute(...), the outer calls have
	// callText that includes the full chain. We only want the innermost match.
	for pattern, op := range rustDBPatterns {
		if callText == pattern {
			dbOp := &graph.Node{
				ID:         NodeID(graph.NodeDBOperation, pattern, file, line, col),
				Kind:       graph.NodeDBOperation,
				Name:       pattern,
				File:       file,
				Line:       line,
				Column:     col,
				Language:   "rust",
				Properties: map[string]string{"operation": op},
				Operation:  op,
			}
			result.DBOperations = append(result.DBOperations, dbOp)
			return
		}
	}
}

// extractMacroInvocation creates a CallSite node from a macro_invocation and detects DB macros.
func (rp *RustParser) extractMacroInvocation(node *sitter.Node, src []byte, file string, result *ParseResult) {
	// Get the macro name: could be an identifier or scoped_identifier
	macroName := ""
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}
		switch child.Type() {
		case "identifier":
			macroName = child.Content(src)
		case "scoped_identifier":
			macroName = child.Content(src)
		}
		if macroName != "" {
			break
		}
	}
	if macroName == "" {
		return
	}

	// Append "!" for macro naming convention
	displayName := macroName + "!"
	line := int(node.StartPoint().Row) + 1
	col := int(node.StartPoint().Column)

	cs := &graph.Node{
		ID:         NodeID(graph.NodeCallSite, displayName, file, line, col),
		Kind:       graph.NodeCallSite,
		Name:       displayName,
		File:       file,
		Line:       line,
		Column:     col,
		Language:   "rust",
		Properties: map[string]string{"is_macro": "true"},
		IsMacro:    true,
	}
	result.CallSites = append(result.CallSites, cs)

	// Check if this macro is a DB operation (e.g., sqlx::query_as!)
	if op, ok := rustDBPatterns[macroName]; ok {
		// Check token tree content for write indicators
		actualOp := op
		fullText := node.Content(src)
		if strings.Contains(fullText, "INSERT") || strings.Contains(fullText, "UPDATE") || strings.Contains(fullText, "DELETE") {
			actualOp = "write"
		}

		dbOp := &graph.Node{
			ID:         NodeID(graph.NodeDBOperation, displayName, file, line, col),
			Kind:       graph.NodeDBOperation,
			Name:       displayName,
			File:       file,
			Line:       line,
			Column:     col,
			Language:   "rust",
			Properties: map[string]string{"operation": actualOp},
			Operation:  actualOp,
		}
		result.DBOperations = append(result.DBOperations, dbOp)
	}
}

// extractStructExpression creates a StructLiteral node from a struct_expression (e.g., Foo { bar: 1 }).
func (rp *RustParser) extractStructExpression(node *sitter.Node, src []byte, file string, result *ParseResult) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}
	name := nameNode.Content(src)

	var fieldNames []string
	body := node.ChildByFieldName("body")
	if body != nil {
		for i := 0; i < int(body.ChildCount()); i++ {
			child := body.Child(i)
			if child == nil {
				continue
			}
			if child.Type() == "field_initializer" {
				fieldNode := child.ChildByFieldName("name")
				if fieldNode == nil {
					fieldNode = child.ChildByFieldName("field")
				}
				if fieldNode != nil {
					fieldNames = append(fieldNames, fieldNode.Content(src))
				}
			}
		}
	}

	line := int(node.StartPoint().Row) + 1
	col := int(node.StartPoint().Column)

	sl := &graph.Node{
		ID:         NodeID(graph.NodeStructLiteral, name, file, line, col),
		Kind:       graph.NodeStructLiteral,
		Name:       name,
		File:       file,
		Line:       line,
		Column:     col,
		EndLine:    int(node.EndPoint().Row) + 1,
		Language:   "rust",
		Properties: map[string]string{
			"type":   name,
			"fields": strings.Join(fieldNames, ","),
		},
		StructType: name,
		FieldNames: fieldNames,
	}
	result.StructLiterals = append(result.StructLiterals, sl)
}

// ---------------------------------------------------------------------------
// Intraprocedural data flow analysis
// ---------------------------------------------------------------------------

// extractRustParamNames extracts parameter names from a Rust parameters node.
// It excludes self/&self/&mut self parameters (similar to Python's self exclusion).
func extractRustParamNames(params *sitter.Node, src []byte) []string {
	var names []string
	for i := 0; i < int(params.ChildCount()); i++ {
		param := params.Child(i)
		if param == nil {
			continue
		}
		switch param.Type() {
		case "parameter":
			// Regular parameter: identifier child is the name
			for j := 0; j < int(param.ChildCount()); j++ {
				child := param.Child(j)
				if child != nil && child.Type() == "identifier" {
					names = append(names, child.Content(src))
					break
				}
			}
		case "self_parameter":
			// Skip self/&self/&mut self
			continue
		}
	}
	return names
}

// analyzeFunctionBody creates variable/parameter nodes and data flow edges
// for a single Rust function body. One SymbolTable + FlowBuilder per function.
func (rp *RustParser) analyzeFunctionBody(body *sitter.Node, src []byte, file string, fn *graph.Node, result *ParseResult) {
	st := dataflow.NewSymbolTable()
	fb := dataflow.NewFlowBuilder()

	// Create parameter nodes and register them in the symbol table
	for _, pName := range fn.ParamNames {
		paramID := NodeID(graph.NodeParameter, pName, file, fn.Line, 0)
		paramNode := &graph.Node{
			ID:       paramID,
			Kind:     graph.NodeParameter,
			Name:     pName,
			File:     file,
			Line:     fn.Line,
			Language: "rust",
		}
		fb.AddParameter(paramNode)
		fb.AddDeclares(fn.ID, paramID)
		st.Define(pName, paramID)
	}

	// Walk body statements sequentially
	childCount := int(body.ChildCount())
	for i := 0; i < childCount; i++ {
		child := body.Child(i)
		if child == nil {
			continue
		}

		// Check for implicit return: last non-brace expression in block
		isLastExpr := false
		if i == childCount-1 || (i == childCount-2 && body.Child(childCount-1) != nil && body.Child(childCount-1).Type() == "}") {
			if child.Type() == "identifier" || child.Type() == "field_expression" ||
				child.Type() == "call_expression" || child.Type() == "macro_invocation" {
				isLastExpr = true
			}
		}

		if isLastExpr {
			rp.analyzeImplicitReturn(child, src, file, fn, st, fb)
		} else {
			rp.analyzeStatementDF(child, src, file, fn, st, fb)
		}

		if fb.VariableCount() >= dataflow.MaxVariablesPerFunction {
			fn.Annotations["dataflow:truncated"] = true
			break
		}
	}

	// Merge results into ParseResult
	nodes, edges := fb.Result()
	for _, n := range nodes {
		switch n.Kind {
		case graph.NodeVariable:
			result.Variables = append(result.Variables, n)
		case graph.NodeParameter:
			result.Parameters = append(result.Parameters, n)
		}
	}
	result.Edges = append(result.Edges, edges...)

	// --- CFG construction pass ---
	posMap := buildPositionMap(nodes)
	for _, p := range result.Parameters {
		if p.File == file {
			k := posKey{Line: p.Line, Col: p.Column}
			posMap[k] = append(posMap[k], p.ID)
		}
	}
	rp.buildCFG(body, src, file, fn, posMap, result)
}

// analyzeStatementDF dispatches to handlers based on tree-sitter node type.
func (rp *RustParser) analyzeStatementDF(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	switch node.Type() {
	case "let_declaration":
		rp.analyzeLetDeclDF(node, src, file, fn, st, fb)
	case "assignment_expression":
		rp.analyzeAssignmentDF(node, src, file, fn, st, fb)
	case "expression_statement":
		// Check for return_expression, assignment_expression, or standalone call
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child == nil {
				continue
			}
			switch child.Type() {
			case "return_expression":
				rp.analyzeReturnDF(child, src, file, fn, st, fb)
			case "assignment_expression":
				rp.analyzeAssignmentDF(child, src, file, fn, st, fb)
			case "call_expression":
				rp.analyzeCallArgsDF(child, src, file, st, fb)
			case "macro_invocation":
				rp.analyzeMacroArgsDF(child, src, file, st, fb)
			}
		}
	case "return_expression":
		rp.analyzeReturnDF(node, src, file, fn, st, fb)
	default:
		// Recurse into block-level structures (if, for, match, loop, etc.)
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child == nil {
				continue
			}
			rp.analyzeStatementDF(child, src, file, fn, st, fb)
			if fb.VariableCount() >= dataflow.MaxVariablesPerFunction {
				return
			}
		}
	}
}

// analyzeLetDeclDF handles `let x = expr` and `let x: Type = expr`.
func (rp *RustParser) analyzeLetDeclDF(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	// In Rust tree-sitter, let_declaration has children: let, identifier (pattern), optional type, =, value
	// The pattern is typically an identifier
	var nameNode *sitter.Node
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil && child.Type() == "identifier" {
			nameNode = child
			break
		}
	}
	if nameNode == nil {
		return
	}

	name := nameNode.Content(src)
	line := int(nameNode.StartPoint().Row) + 1
	col := int(nameNode.StartPoint().Column)
	varID := NodeID(graph.NodeVariable, name, file, line, col)
	varNode := &graph.Node{
		ID:       varID,
		Kind:     graph.NodeVariable,
		Name:     name,
		File:     file,
		Line:     line,
		Column:   col,
		Language: "rust",
	}
	fb.AddVariable(varNode)
	st.Define(name, varID)

	// Find the value expression (after the = sign)
	valueNode := rp.findLetValue(node)
	if valueNode != nil {
		sourceID := rp.resolveRHSSourceDF(valueNode, src, file, fn, st, fb)
		if sourceID != "" {
			fb.AddAssign(sourceID, varID)
		}
		rp.emitReadsDF(valueNode, src, varID, st, fb)
	}
}

// findLetValue finds the value expression in a let_declaration.
// It looks for the expression after the '=' token.
func (rp *RustParser) findLetValue(node *sitter.Node) *sitter.Node {
	foundEq := false
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}
		if child.Type() == "=" {
			foundEq = true
			continue
		}
		if foundEq && child.Type() != ";" {
			return child
		}
	}
	return nil
}

// analyzeAssignmentDF handles `x = expr`.
func (rp *RustParser) analyzeAssignmentDF(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	left := node.ChildByFieldName("left")
	right := node.ChildByFieldName("right")
	if left == nil || right == nil {
		return
	}

	var targetName string
	if left.Type() == "identifier" {
		targetName = left.Content(src)
	}
	if targetName == "" {
		return
	}

	targetID, ok := st.Resolve(targetName)
	if !ok {
		return
	}

	sourceID := rp.resolveRHSSourceDF(right, src, file, fn, st, fb)
	if sourceID != "" {
		fb.AddAssign(sourceID, targetID)
	}
	rp.emitReadsDF(right, src, targetID, st, fb)
}

// analyzeReturnDF handles explicit `return expr`.
func (rp *RustParser) analyzeReturnDF(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}
		switch child.Type() {
		case "identifier":
			name := child.Content(src)
			if varID, ok := st.Resolve(name); ok {
				fb.AddReturn(varID, fn.ID)
			}
		case "field_expression":
			rootID := rp.resolveFieldExprDF(child, src, file, fn, st, fb)
			if rootID != "" {
				fb.AddReturn(rootID, fn.ID)
			}
		default:
			rp.walkReturnExprDF(child, src, file, fn, st, fb)
		}
	}
}

// analyzeImplicitReturn handles implicit returns (last expression in block body).
func (rp *RustParser) analyzeImplicitReturn(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	switch node.Type() {
	case "identifier":
		name := node.Content(src)
		if varID, ok := st.Resolve(name); ok {
			fb.AddReturn(varID, fn.ID)
		}
	case "field_expression":
		rootID := rp.resolveFieldExprDF(node, src, file, fn, st, fb)
		if rootID != "" {
			fb.AddReturn(rootID, fn.ID)
		}
	case "call_expression":
		// Implicit return of a call result
		fnNode := node.ChildByFieldName("function")
		if fnNode != nil {
			callText := fnNode.Content(src)
			line := int(node.StartPoint().Row) + 1
			col := int(node.StartPoint().Column)
			callID := NodeID(graph.NodeCallSite, callText, file, line, col)
			rp.analyzeCallArgsDF(node, src, file, st, fb)
			fb.AddReturn(callID, fn.ID)
		}
	default:
		rp.walkReturnExprDF(node, src, file, fn, st, fb)
	}
}

// walkReturnExprDF recursively finds identifiers in return expressions.
func (rp *RustParser) walkReturnExprDF(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	if node == nil {
		return
	}
	if node.Type() == "identifier" {
		name := node.Content(src)
		if varID, ok := st.Resolve(name); ok {
			fb.AddReturn(varID, fn.ID)
		}
		return
	}
	if node.Type() == "field_expression" {
		rootID := rp.resolveFieldExprDF(node, src, file, fn, st, fb)
		if rootID != "" {
			fb.AddReturn(rootID, fn.ID)
		}
		return
	}
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil {
			rp.walkReturnExprDF(child, src, file, fn, st, fb)
		}
	}
}

// resolveRHSSourceDF resolves the primary data source from an RHS expression.
func (rp *RustParser) resolveRHSSourceDF(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) string {
	if node == nil {
		return ""
	}

	switch node.Type() {
	case "call_expression":
		fnNode := node.ChildByFieldName("function")
		if fnNode == nil {
			return ""
		}
		callText := fnNode.Content(src)
		line := int(node.StartPoint().Row) + 1
		col := int(node.StartPoint().Column)
		callID := NodeID(graph.NodeCallSite, callText, file, line, col)
		rp.analyzeCallArgsDF(node, src, file, st, fb)
		// For method calls like user.email.clone(), resolve the field chain
		// to emit field_access edges from the receiver object
		if fnNode.Type() == "field_expression" {
			rp.resolveFieldExprDF(fnNode, src, file, fn, st, fb)
		}
		return callID

	case "identifier":
		name := node.Content(src)
		if varID, ok := st.Resolve(name); ok {
			return varID
		}
		return ""

	case "field_expression":
		return rp.resolveFieldExprDF(node, src, file, fn, st, fb)

	case "index_expression":
		// e.g., payload["name"]: resolve the object
		if node.ChildCount() > 0 {
			obj := node.Child(0)
			if obj != nil && obj.Type() == "identifier" {
				name := obj.Content(src)
				if varID, ok := st.Resolve(name); ok {
					return varID
				}
			}
		}
		return ""

	case "reference_expression":
		// e.g., &body: resolve the inner expression
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child != nil && child.Type() != "&" {
				return rp.resolveRHSSourceDF(child, src, file, fn, st, fb)
			}
		}
		return ""

	case "try_expression":
		// e.g., expr? : resolve the inner expression
		if node.ChildCount() > 0 {
			return rp.resolveRHSSourceDF(node.Child(0), src, file, fn, st, fb)
		}
		return ""

	case "macro_invocation":
		// e.g., format!(...): return the macro call site ID
		macroName := ""
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child == nil {
				continue
			}
			if child.Type() == "identifier" || child.Type() == "scoped_identifier" {
				macroName = child.Content(src)
				break
			}
		}
		if macroName == "" {
			return ""
		}
		displayName := macroName + "!"
		line := int(node.StartPoint().Row) + 1
		col := int(node.StartPoint().Column)
		callID := NodeID(graph.NodeCallSite, displayName, file, line, col)
		// Process macro arguments for passes_to edges
		rp.analyzeMacroArgsDF(node, src, file, st, fb)
		return callID

	case "binary_expression":
		left := node.ChildByFieldName("left")
		right := node.ChildByFieldName("right")
		if leftID := rp.resolveRHSSourceDF(left, src, file, fn, st, fb); leftID != "" {
			return leftID
		}
		return rp.resolveRHSSourceDF(right, src, file, fn, st, fb)

	default:
		return ""
	}
}

// resolveFieldExprDF handles field_expression chains like user.email or request.to_string.
// Creates a single synthetic field variable with the full text and emits one field_access
// edge from the root identifier to the synthetic variable.
func (rp *RustParser) resolveFieldExprDF(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) string {
	if node == nil || node.Type() != "field_expression" {
		return ""
	}

	fullText := node.Content(src)

	// Walk to the leftmost identifier to find the root
	var rootID string
	current := node
	for {
		// In Rust field_expression, the first child is the value (object)
		if current.ChildCount() == 0 {
			break
		}
		obj := current.Child(0)
		if obj == nil {
			break
		}
		if obj.Type() == "identifier" {
			rootName := obj.Content(src)
			if id, ok := st.Resolve(rootName); ok {
				rootID = id
			}
			break
		} else if obj.Type() == "field_expression" {
			current = obj
		} else if obj.Type() == "call_expression" {
			// e.g., payload["name"].as_str().unwrap - walk into call's function
			fnChild := obj.ChildByFieldName("function")
			if fnChild != nil && fnChild.Type() == "field_expression" {
				current = fnChild
			} else {
				break
			}
		} else if obj.Type() == "index_expression" {
			// e.g., payload["name"]: resolve the indexed object
			if obj.ChildCount() > 0 {
				first := obj.Child(0)
				if first != nil && first.Type() == "identifier" {
					rootName := first.Content(src)
					if id, ok := st.Resolve(rootName); ok {
						rootID = id
					}
				}
			}
			break
		} else {
			break
		}
	}

	if rootID == "" {
		return ""
	}

	// Create ONE synthetic field variable with the full text
	line := int(node.StartPoint().Row) + 1
	col := int(node.StartPoint().Column)
	fieldVarID := NodeID(graph.NodeVariable, fullText, file, line, col)
	fieldVarNode := &graph.Node{
		ID:       fieldVarID,
		Kind:     graph.NodeVariable,
		Name:     fullText,
		File:     file,
		Line:     line,
		Column:   col,
		Language: "rust",
	}
	fb.AddVariable(fieldVarNode)
	fb.AddFieldAccess(rootID, fieldVarID)

	return fieldVarID
}

// analyzeCallArgsDF processes arguments to a call expression, creating
// passes_to edges and mutates edges for &x reference_expression patterns.
func (rp *RustParser) analyzeCallArgsDF(callNode *sitter.Node, src []byte, file string, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	fnNode := callNode.ChildByFieldName("function")
	if fnNode == nil {
		return
	}
	callText := fnNode.Content(src)
	line := int(callNode.StartPoint().Row) + 1
	col := int(callNode.StartPoint().Column)
	callID := NodeID(graph.NodeCallSite, callText, file, line, col)

	args := callNode.ChildByFieldName("arguments")
	if args == nil {
		return
	}

	for i := 0; i < int(args.ChildCount()); i++ {
		arg := args.Child(i)
		if arg == nil {
			continue
		}

		switch arg.Type() {
		case "identifier":
			name := arg.Content(src)
			if varID, ok := st.Resolve(name); ok {
				fb.AddPassesTo(varID, callID)
			}

		case "reference_expression":
			// &x: passes_to + mutates
			for j := 0; j < int(arg.ChildCount()); j++ {
				child := arg.Child(j)
				if child == nil || child.Type() == "&" {
					continue
				}
				if child.Type() == "identifier" {
					name := child.Content(src)
					if varID, ok := st.Resolve(name); ok {
						fb.AddPassesTo(varID, callID)
						fb.AddMutates(callID, varID)
					}
				}
			}

		case "field_expression":
			// e.g., user.name: passes_to from root identifier
			rp.passFieldToCall(arg, src, callID, st, fb)
		}
	}
}

// analyzeMacroArgsDF processes arguments inside a macro_invocation's token_tree.
// Finds identifiers and emits passes_to edges.
func (rp *RustParser) analyzeMacroArgsDF(node *sitter.Node, src []byte, file string, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	macroName := ""
	var tokenTree *sitter.Node
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}
		if child.Type() == "identifier" || child.Type() == "scoped_identifier" {
			macroName = child.Content(src)
		}
		if child.Type() == "token_tree" {
			tokenTree = child
		}
	}
	if macroName == "" || tokenTree == nil {
		return
	}

	displayName := macroName + "!"
	line := int(node.StartPoint().Row) + 1
	col := int(node.StartPoint().Column)
	callID := NodeID(graph.NodeCallSite, displayName, file, line, col)

	// Walk token_tree for identifiers
	rp.walkTokenTreeForArgs(tokenTree, src, callID, st, fb)
}

// walkTokenTreeForArgs walks a token_tree looking for identifier arguments
// and emits passes_to edges.
func (rp *RustParser) walkTokenTreeForArgs(node *sitter.Node, src []byte, callID string, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}
		switch child.Type() {
		case "identifier":
			name := child.Content(src)
			if varID, ok := st.Resolve(name); ok {
				fb.AddPassesTo(varID, callID)
			}
		case "reference_expression":
			// &x in macro args
			for j := 0; j < int(child.ChildCount()); j++ {
				gc := child.Child(j)
				if gc != nil && gc.Type() == "identifier" {
					name := gc.Content(src)
					if varID, ok := st.Resolve(name); ok {
						fb.AddPassesTo(varID, callID)
						fb.AddMutates(callID, varID)
					}
				}
			}
		case "token_tree":
			rp.walkTokenTreeForArgs(child, src, callID, st, fb)
		}
	}
}

// passFieldToCall finds the root identifier of a field expression
// and emits a passes_to edge.
func (rp *RustParser) passFieldToCall(node *sitter.Node, src []byte, callID string, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	if node == nil {
		return
	}
	switch node.Type() {
	case "identifier":
		name := node.Content(src)
		if varID, ok := st.Resolve(name); ok {
			fb.AddPassesTo(varID, callID)
		}
	case "field_expression":
		if node.ChildCount() > 0 {
			rp.passFieldToCall(node.Child(0), src, callID, st, fb)
		}
	}
}

// emitReadsDF walks a compound expression (binary_expression, etc.) to find all
// identifier references and emits reads edges.
func (rp *RustParser) emitReadsDF(node *sitter.Node, src []byte, targetID string, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	if node == nil {
		return
	}

	if node.Type() != "binary_expression" {
		return
	}

	rp.walkForReadsDF(node, src, targetID, st, fb)
}

// walkForReadsDF recursively finds identifiers in an expression tree and emits reads edges.
func (rp *RustParser) walkForReadsDF(node *sitter.Node, src []byte, targetID string, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	if node == nil {
		return
	}

	if node.Type() == "identifier" {
		name := node.Content(src)
		if varID, ok := st.Resolve(name); ok {
			fb.AddRead(varID, targetID)
		}
		return
	}

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil {
			rp.walkForReadsDF(child, src, targetID, st, fb)
		}
	}
}

// ---------------------------------------------------------------------------
// Control flow graph construction
// ---------------------------------------------------------------------------

func (rp *RustParser) buildCFG(body *sitter.Node, src []byte, file string, fn *graph.Node, posMap map[posKey][]string, result *ParseResult) {
	cb := dataflow.NewCFGBuilder(fn.ID, file)

	entryBlock := cb.NewBlock("entry", fn.Line)
	exitBlock := cb.NewBlock("exit", 0)
	cb.AddEdge(fn.ID, entryBlock.ID, "entry")

	for _, p := range result.Parameters {
		if p.File == file && p.Line == fn.Line {
			cb.AddMember(entryBlock, p.ID)
		}
	}

	lastBlock := rp.buildCFGStatementsRust(body, src, file, fn, entryBlock, exitBlock, cb, posMap)
	if lastBlock != nil {
		cb.AddEdge(lastBlock.ID, exitBlock.ID, "exit")
	}

	blocks, edges := cb.Result()
	result.BasicBlocks = append(result.BasicBlocks, blocks...)
	result.Edges = append(result.Edges, edges...)
}

func (rp *RustParser) buildCFGStatementsRust(container *sitter.Node, src []byte, file string, fn *graph.Node, currentBlock *graph.Node, exitBlock *graph.Node, cb *dataflow.CFGBuilder, posMap map[posKey][]string) *graph.Node {
	if container == nil {
		return currentBlock
	}
	for i := 0; i < int(container.ChildCount()); i++ {
		if cb.BlockCount() >= dataflow.MaxBlocksPerFunction {
			if fn.Annotations == nil {
				fn.Annotations = make(map[string]bool)
			}
			fn.Annotations["cfg:truncated"] = true
			return currentBlock
		}

		child := container.Child(i)
		if child == nil {
			continue
		}

		switch child.Type() {
		case "if_expression":
			currentBlock = rp.buildCFGIfRust(child, src, file, fn, currentBlock, exitBlock, cb, posMap)
			if currentBlock == nil {
				return nil
			}

		case "for_expression", "while_expression", "loop_expression":
			currentBlock = rp.buildCFGForRust(child, src, file, fn, currentBlock, exitBlock, cb, posMap)
			if currentBlock == nil {
				return nil
			}

		case "match_expression":
			currentBlock = rp.buildCFGMatchRust(child, src, file, fn, currentBlock, exitBlock, cb, posMap)
			if currentBlock == nil {
				return nil
			}

		case "return_expression":
			for _, id := range collectNodeIDs(child, posMap) {
				cb.AddMember(currentBlock, id)
			}
			cb.AddEdge(currentBlock.ID, exitBlock.ID, "exit")
			return nil

		case "expression_statement":
			// Check for inner control flow expressions
			for j := 0; j < int(child.ChildCount()); j++ {
				inner := child.Child(j)
				if inner == nil {
					continue
				}
				switch inner.Type() {
				case "if_expression":
					currentBlock = rp.buildCFGIfRust(inner, src, file, fn, currentBlock, exitBlock, cb, posMap)
					if currentBlock == nil {
						return nil
					}
				case "for_expression", "while_expression", "loop_expression":
					currentBlock = rp.buildCFGForRust(inner, src, file, fn, currentBlock, exitBlock, cb, posMap)
					if currentBlock == nil {
						return nil
					}
				case "match_expression":
					currentBlock = rp.buildCFGMatchRust(inner, src, file, fn, currentBlock, exitBlock, cb, posMap)
					if currentBlock == nil {
						return nil
					}
				case "return_expression":
					for _, id := range collectNodeIDs(inner, posMap) {
						cb.AddMember(currentBlock, id)
					}
					cb.AddEdge(currentBlock.ID, exitBlock.ID, "exit")
					return nil
				default:
					for _, id := range collectNodeIDs(inner, posMap) {
						cb.AddMember(currentBlock, id)
					}
				}
			}

		default:
			for _, id := range collectNodeIDs(child, posMap) {
				cb.AddMember(currentBlock, id)
			}
		}
	}
	return currentBlock
}

func (rp *RustParser) buildCFGIfRust(node *sitter.Node, src []byte, file string, fn *graph.Node, currentBlock *graph.Node, exitBlock *graph.Node, cb *dataflow.CFGBuilder, posMap map[posKey][]string) *graph.Node {
	cond := node.ChildByFieldName("condition")
	if cond != nil {
		for _, id := range collectNodeIDs(cond, posMap) {
			cb.AddMember(currentBlock, id)
		}
	}
	condBlock := currentBlock

	consequence := node.ChildByFieldName("consequence")
	thenLine := int(node.StartPoint().Row) + 1
	if consequence != nil {
		thenLine = int(consequence.StartPoint().Row) + 1
	}
	thenBlock := cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), thenLine)
	cb.AddEdge(condBlock.ID, thenBlock.ID, "true_branch")
	thenEnd := rp.buildCFGStatementsRust(consequence, src, file, fn, thenBlock, exitBlock, cb, posMap)

	alternative := node.ChildByFieldName("alternative")
	var elseEnd *graph.Node
	if alternative != nil {
		elseLine := int(alternative.StartPoint().Row) + 1
		elseBlock := cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), elseLine)
		cb.AddEdge(condBlock.ID, elseBlock.ID, "false_branch")

		// Check for else if
		hasElseIf := false
		for j := 0; j < int(alternative.ChildCount()); j++ {
			c := alternative.Child(j)
			if c != nil && c.Type() == "if_expression" {
				elseEnd = rp.buildCFGIfRust(c, src, file, fn, elseBlock, exitBlock, cb, posMap)
				hasElseIf = true
				break
			}
		}
		if !hasElseIf {
			elseEnd = rp.buildCFGStatementsRust(alternative, src, file, fn, elseBlock, exitBlock, cb, posMap)
		}
	}

	if thenEnd != nil || elseEnd != nil || alternative == nil {
		mergeLine := int(node.EndPoint().Row) + 1
		mergeBlock := cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), mergeLine)
		if thenEnd != nil {
			cb.AddEdge(thenEnd.ID, mergeBlock.ID, "fallthrough")
		}
		if alternative == nil {
			cb.AddEdge(condBlock.ID, mergeBlock.ID, "false_branch")
		} else if elseEnd != nil {
			cb.AddEdge(elseEnd.ID, mergeBlock.ID, "fallthrough")
		}
		return mergeBlock
	}
	return nil
}

func (rp *RustParser) buildCFGForRust(node *sitter.Node, src []byte, file string, fn *graph.Node, currentBlock *graph.Node, exitBlock *graph.Node, cb *dataflow.CFGBuilder, posMap map[posKey][]string) *graph.Node {
	headerLine := int(node.StartPoint().Row) + 1
	headerBlock := cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), headerLine)
	cb.AddEdge(currentBlock.ID, headerBlock.ID, "fallthrough")

	body := node.ChildByFieldName("body")
	bodyLine := headerLine + 1
	if body != nil {
		bodyLine = int(body.StartPoint().Row) + 1
	}
	bodyBlock := cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), bodyLine)
	cb.AddEdge(headerBlock.ID, bodyBlock.ID, "true_branch")

	var bodyEnd *graph.Node
	if body != nil {
		bodyEnd = rp.buildCFGStatementsRust(body, src, file, fn, bodyBlock, exitBlock, cb, posMap)
	} else {
		bodyEnd = bodyBlock
	}

	if bodyEnd != nil {
		cb.AddEdge(bodyEnd.ID, headerBlock.ID, "loop_back")
	}

	afterLine := int(node.EndPoint().Row) + 1
	afterBlock := cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), afterLine)
	cb.AddEdge(headerBlock.ID, afterBlock.ID, "loop_exit")
	return afterBlock
}

func (rp *RustParser) buildCFGMatchRust(node *sitter.Node, src []byte, file string, fn *graph.Node, currentBlock *graph.Node, exitBlock *graph.Node, cb *dataflow.CFGBuilder, posMap map[posKey][]string) *graph.Node {
	condBlock := currentBlock
	body := node.ChildByFieldName("body")
	if body == nil {
		return currentBlock
	}

	mergeLine := int(node.EndPoint().Row) + 1
	mergeBlock := cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), mergeLine)

	hasWildcard := false

	for i := 0; i < int(body.ChildCount()); i++ {
		child := body.Child(i)
		if child == nil || child.Type() != "match_arm" {
			continue
		}

		// Check for wildcard pattern (_)
		pattern := child.ChildByFieldName("pattern")
		if pattern != nil && pattern.Content(src) == "_" {
			hasWildcard = true
		}

		armLine := int(child.StartPoint().Row) + 1
		armBlock := cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), armLine)

		label := "true_branch"
		if hasWildcard && pattern != nil && pattern.Content(src) == "_" {
			label = "false_branch"
		}
		cb.AddEdge(condBlock.ID, armBlock.ID, label)

		for _, id := range collectNodeIDs(child, posMap) {
			cb.AddMember(armBlock, id)
		}

		// Check if arm contains return
		hasReturn := false
		for j := 0; j < int(child.ChildCount()); j++ {
			c := child.Child(j)
			if c != nil && c.Type() == "return_expression" {
				cb.AddEdge(armBlock.ID, exitBlock.ID, "exit")
				hasReturn = true
				break
			}
		}

		if !hasReturn {
			cb.AddEdge(armBlock.ID, mergeBlock.ID, "fallthrough")
		}
	}

	if !hasWildcard {
		cb.AddEdge(condBlock.ID, mergeBlock.ID, "false_branch")
	}

	return mergeBlock
}
