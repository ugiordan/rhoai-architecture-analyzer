package parser

import (
	"context"
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/typescript/tsx"

	"github.com/ugiordan/architecture-analyzer/pkg/dataflow"
	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

// TypeScriptSkipDirs lists directories that should be skipped when scanning TypeScript projects.
var TypeScriptSkipDirs = []string{"node_modules", "dist", "build", ".next", "coverage"}

// TypeScriptTestPatterns lists filename patterns that identify test files.
var TypeScriptTestPatterns = []string{"*.spec.ts", "*.test.ts", "*.spec.tsx", "*.test.tsx"}

// tsHTTPMethods maps Express-style HTTP method names to their uppercase equivalents.
var tsHTTPMethods = map[string]string{
	"get":    "GET",
	"post":   "POST",
	"put":    "PUT",
	"delete": "DELETE",
	"patch":  "PATCH",
	"use":    "USE",
}

// tsDBOps maps method names to their operation type (read/write) for DB detection.
var tsDBOps = map[string]string{
	"query":   "read",
	"findOne": "read",
	"find":    "read",
	"execute": "write",
	"save":    "write",
	"create":  "write",
	"update":  "write",
	"delete":  "write",
}

// TypeScriptParser extracts code property graph nodes from TypeScript/TSX source files
// using tree-sitter. Each goroutine MUST use its own TypeScriptParser instance
// (tree-sitter parsers are not thread-safe).
type TypeScriptParser struct {
	parser *sitter.Parser
}

// NewTypeScriptParser creates a parser for TypeScript source files backed by tree-sitter.
// Uses the TSX grammar for all files since it's a superset of TypeScript.
func NewTypeScriptParser() *TypeScriptParser {
	p := sitter.NewParser()
	p.SetLanguage(tsx.GetLanguage())
	return &TypeScriptParser{parser: p}
}

func (tp *TypeScriptParser) Language() string     { return "typescript" }
func (tp *TypeScriptParser) Extensions() []string { return []string{".ts", ".tsx"} }
func (tp *TypeScriptParser) Clone() Parser {
	p := sitter.NewParser()
	p.SetLanguage(tsx.GetLanguage())
	return &TypeScriptParser{parser: p}
}

// computeTypeScriptComplexity counts decision points in a TypeScript function body.
// Complexity = 1 (base) + count of: if, for, for-in, while, do, switch_case, switch_default, catch, &&, ternary.
func computeTypeScriptComplexity(node *sitter.Node) int {
	count := 1
	countTypeScriptDecisionPoints(node, &count)
	return count
}

func countTypeScriptDecisionPoints(node *sitter.Node, count *int) {
	if node == nil {
		return
	}
	switch node.Type() {
	case "if_statement":
		*count++
	case "for_statement", "for_in_statement":
		*count++
	case "while_statement", "do_statement":
		*count++
	case "switch_case", "switch_default":
		*count++
	case "catch_clause":
		*count++
	case "ternary_expression":
		*count++
	case "binary_expression":
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child != nil {
				op := child.Type()
				if op == "&&" {
					*count++
					break
				}
			}
		}
	}
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}
		// Handle else-if chains: when an else_clause contains an if_statement,
		// skip counting that if_statement itself (it's not a new branch, just a
		// continuation of the chain) but still recurse into its children.
		if node.Type() == "else_clause" && child.Type() == "if_statement" {
			countTypeScriptDecisionPointsSkipSelf(child, count)
			continue
		}
		countTypeScriptDecisionPoints(child, count)
	}
}

// countTypeScriptDecisionPointsSkipSelf recurses into a node's children without
// counting the node itself. Used for else-if chains to avoid double-counting.
func countTypeScriptDecisionPointsSkipSelf(node *sitter.Node, count *int) {
	if node == nil {
		return
	}
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil {
			countTypeScriptDecisionPoints(child, count)
		}
	}
}

// ParseFile parses a TypeScript/TSX source file and returns extracted nodes and edges.
// Declaration files (.d.ts) are skipped, returning an empty result.
func (tp *TypeScriptParser) ParseFile(path string, content []byte) (*ParseResult, error) {
	// Skip declaration files
	if strings.HasSuffix(path, ".d.ts") {
		return &ParseResult{}, nil
	}

	if len(content) > MaxFileSize {
		return nil, fmt.Errorf("file too large (%d bytes, max %d)", len(content), MaxFileSize)
	}

	tree, err := tp.parser.ParseCtx(context.Background(), nil, content)
	if err != nil {
		return nil, fmt.Errorf("tree-sitter parse failed: %w", err)
	}
	defer tree.Close()

	result := &ParseResult{}
	root := tree.RootNode()
	tp.walk(root, content, path, "", result)
	return result, nil
}

// walk recursively traverses the AST. className tracks the enclosing class for methods.
func (tp *TypeScriptParser) walk(node *sitter.Node, src []byte, file, className string, result *ParseResult) {
	switch node.Type() {
	case "function_declaration":
		tp.extractFunction(node, src, file, result)
	case "method_definition":
		tp.extractMethod(node, src, file, className, result)
	case "class_declaration":
		tp.extractClass(node, src, file, result)
		return // children handled inside extractClass
	case "lexical_declaration", "variable_declaration":
		tp.extractArrowFunctions(node, src, file, result)
	case "call_expression":
		tp.extractCallSite(node, src, file, result)
	case "new_expression":
		tp.extractNewExpression(node, src, file, result)
	case "jsx_self_closing_element":
		tp.extractJSXRoute(node, src, file, result)
	case "jsx_opening_element":
		tp.extractJSXRoute(node, src, file, result)
	}

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil {
			tp.walk(child, src, file, className, result)
		}
	}
}

// extractFunction handles function_declaration nodes.
func (tp *TypeScriptParser) extractFunction(node *sitter.Node, src []byte, file string, result *ParseResult) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}
	line := int(node.StartPoint().Row) + 1
	col := int(node.StartPoint().Column)
	fn := &graph.Node{
		ID:          NodeID(graph.NodeFunction, nameNode.Content(src), file, line, col),
		Kind:        graph.NodeFunction,
		Name:        nameNode.Content(src),
		File:        file,
		Line:        line,
		Column:      col,
		EndLine:     int(node.EndPoint().Row) + 1,
		Language:    "typescript",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	// Extract parameter types and names for request handler detection and data flow
	params := node.ChildByFieldName("parameters")
	if params != nil {
		fn.ParamTypes = extractParamTypes(params, src)
		fn.ParamNames = extractTSParamNames(params, src)
	}
	body := node.ChildByFieldName("body")
	if body != nil {
		fn.Complexity = computeTypeScriptComplexity(body)
	}
	result.Functions = append(result.Functions, fn)

	// Extract intraprocedural data flow from function body
	if body != nil {
		tp.analyzeFunctionBody(body, src, file, fn, result)
	}
}

// extractMethod handles method_definition nodes inside classes.
func (tp *TypeScriptParser) extractMethod(node *sitter.Node, src []byte, file, className string, result *ParseResult) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}
	line := int(node.StartPoint().Row) + 1
	col := int(node.StartPoint().Column)
	fn := &graph.Node{
		ID:          NodeID(graph.NodeFunction, nameNode.Content(src), file, line, col),
		Kind:        graph.NodeFunction,
		Name:        nameNode.Content(src),
		File:        file,
		Line:        line,
		Column:      col,
		EndLine:     int(node.EndPoint().Row) + 1,
		Language:    "typescript",
		TypeName:    className,
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	// Extract parameter types and names for request handler detection and data flow
	params := node.ChildByFieldName("parameters")
	if params != nil {
		fn.ParamTypes = extractParamTypes(params, src)
		fn.ParamNames = extractTSParamNames(params, src)
	}
	body := node.ChildByFieldName("body")
	if body != nil {
		fn.Complexity = computeTypeScriptComplexity(body)
	}
	result.Functions = append(result.Functions, fn)

	// Extract intraprocedural data flow from function body
	if body != nil {
		tp.analyzeFunctionBody(body, src, file, fn, result)
	}
}

// extractClass handles class_declaration nodes, walking the body with the class name set.
func (tp *TypeScriptParser) extractClass(node *sitter.Node, src []byte, file string, result *ParseResult) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}
	clsName := nameNode.Content(src)

	body := node.ChildByFieldName("body")
	if body == nil {
		return
	}
	for i := 0; i < int(body.ChildCount()); i++ {
		child := body.Child(i)
		if child != nil {
			tp.walk(child, src, file, clsName, result)
		}
	}
}

// extractArrowFunctions checks lexical_declaration/variable_declaration children for
// variable_declarator nodes whose value is an arrow_function.
func (tp *TypeScriptParser) extractArrowFunctions(node *sitter.Node, src []byte, file string, result *ParseResult) {
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child == nil || child.Type() != "variable_declarator" {
			continue
		}
		nameNode := child.ChildByFieldName("name")
		valueNode := child.ChildByFieldName("value")
		if nameNode == nil || valueNode == nil {
			continue
		}
		// The value might be an arrow_function directly or wrapped in a type assertion/parenthesized expression.
		if isArrowFunction(valueNode) {
			line := int(child.StartPoint().Row) + 1
			col := int(child.StartPoint().Column)
			fn := &graph.Node{
				ID:          NodeID(graph.NodeFunction, nameNode.Content(src), file, line, col),
				Kind:        graph.NodeFunction,
				Name:        nameNode.Content(src),
				File:        file,
				Line:        line,
				Column:      col,
				EndLine:     int(child.EndPoint().Row) + 1,
				Language:    "typescript",
				Annotations: make(map[string]bool),
				Properties:  make(map[string]string),
			}
			// Find the arrow_function node to extract params and data flow
			arrowNode := findArrowFunction(valueNode)
			if arrowNode != nil {
				params := arrowNode.ChildByFieldName("parameters")
				if params != nil {
					fn.ParamTypes = extractParamTypes(params, src)
					fn.ParamNames = extractTSParamNames(params, src)
				}
				arrowBody := arrowNode.ChildByFieldName("body")
				if arrowBody != nil {
					fn.Complexity = computeTypeScriptComplexity(arrowBody)
				}
			}
			result.Functions = append(result.Functions, fn)

			// Extract intraprocedural data flow from arrow function body
			if arrowNode != nil {
				arrowBody := arrowNode.ChildByFieldName("body")
				if arrowBody != nil {
					tp.analyzeFunctionBody(arrowBody, src, file, fn, result)
				}
			}
		}
	}
}

// isArrowFunction checks if a node is or contains an arrow function.
// Handles cases like: `const x = () => ...` and `const x: React.FC = () => ...`
func isArrowFunction(node *sitter.Node) bool {
	if node.Type() == "arrow_function" {
		return true
	}
	// Check children for arrow function (handles type assertions, parenthesized expressions)
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil && child.Type() == "arrow_function" {
			return true
		}
	}
	return false
}

// extractCallSite creates a CallSite node and checks for Express HTTP handlers and DB operations.
func (tp *TypeScriptParser) extractCallSite(node *sitter.Node, src []byte, file string, result *ParseResult) {
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
		Language:   "typescript",
		Properties: make(map[string]string),
	}
	result.CallSites = append(result.CallSites, cs)

	// Check for Express HTTP handler patterns: app.get, router.post, etc.
	if fnNode.Type() == "member_expression" {
		propNode := fnNode.ChildByFieldName("property")
		if propNode != nil {
			method := propNode.Content(src)
			if httpMethod, ok := tsHTTPMethods[method]; ok {
				tp.maybeExtractExpressHandler(node, src, file, line, col, callText, httpMethod, result)
			} else if op, ok := tsDBOps[method]; ok {
				// Only treat as DB operation if not already matched as HTTP handler
				dbOp := &graph.Node{
					ID:         NodeID(graph.NodeDBOperation, callText, file, line, col),
					Kind:       graph.NodeDBOperation,
					Name:       callText,
					File:       file,
					Line:       line,
					Column:     col,
					Language:   "typescript",
					Operation:  op,
					Properties: map[string]string{"operation": op},
				}
				result.DBOperations = append(result.DBOperations, dbOp)
			}
		}
	}
}

// maybeExtractExpressHandler creates an HTTP endpoint node for Express-style route registrations.
func (tp *TypeScriptParser) maybeExtractExpressHandler(node *sitter.Node, src []byte, file string, line, col int, callText, httpMethod string, result *ParseResult) {
	args := node.ChildByFieldName("arguments")
	if args == nil {
		return
	}

	// The first string argument is the route path
	route := ""
	for i := 0; i < int(args.ChildCount()); i++ {
		arg := args.Child(i)
		if arg == nil {
			continue
		}
		if arg.Type() == "string" || arg.Type() == "template_string" {
			route = stripQuotes(arg.Content(src))
			break
		}
	}

	handler := &graph.Node{
		ID:       NodeID(graph.NodeHTTPEndpoint, callText, file, line, col),
		Kind:     graph.NodeHTTPEndpoint,
		Name:     callText,
		File:     file,
		Line:     line,
		Column:   col,
		Language: "typescript",
		Annotations: make(map[string]bool),
		Properties:  map[string]string{
			"method": httpMethod,
		},
	}
	handler.HTTPMethod = httpMethod
	if route != "" {
		handler.Properties["route"] = route
		handler.Route = route
	}
	result.HTTPHandlers = append(result.HTTPHandlers, handler)
}

// extractNewExpression creates a StructLiteral node for `new X()` expressions.
func (tp *TypeScriptParser) extractNewExpression(node *sitter.Node, src []byte, file string, result *ParseResult) {
	// The constructor is the first named child after "new"
	ctorNode := node.ChildByFieldName("constructor")
	if ctorNode == nil {
		return
	}
	name := ctorNode.Content(src)

	line := int(node.StartPoint().Row) + 1
	col := int(node.StartPoint().Column)
	sl := &graph.Node{
		ID:         NodeID(graph.NodeStructLiteral, name, file, line, col),
		Kind:       graph.NodeStructLiteral,
		Name:       name,
		File:       file,
		Line:       line,
		Column:     col,
		Language:   "typescript",
		Properties: make(map[string]string),
	}
	result.StructLiterals = append(result.StructLiterals, sl)
}

// extractJSXRoute handles <Route path="/x" component={Y} /> elements.
func (tp *TypeScriptParser) extractJSXRoute(node *sitter.Node, src []byte, file string, result *ParseResult) {
	// Get the tag name
	var tagName string
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}
		if child.Type() == "identifier" || child.Type() == "jsx_namespace_name" {
			tagName = child.Content(src)
			break
		}
	}
	if tagName != "Route" {
		return
	}

	// Extract path and component attributes
	route := ""
	componentName := ""
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child == nil || child.Type() != "jsx_attribute" {
			continue
		}
		attrName := ""
		attrValue := ""
		for j := 0; j < int(child.ChildCount()); j++ {
			attrChild := child.Child(j)
			if attrChild == nil {
				continue
			}
			switch attrChild.Type() {
			case "property_identifier":
				attrName = attrChild.Content(src)
			case "jsx_attribute_name":
				attrName = attrChild.Content(src)
			case "string", "jsx_attribute_value":
				attrValue = stripQuotes(attrChild.Content(src))
			case "jsx_expression":
				// {ComponentName} or {<Element />}
				for k := 0; k < int(attrChild.ChildCount()); k++ {
					inner := attrChild.Child(k)
					if inner != nil && inner.Type() == "identifier" {
						attrValue = inner.Content(src)
					}
				}
			}
		}
		switch attrName {
		case "path":
			route = attrValue
		case "component", "element":
			componentName = attrValue
		}
	}

	line := int(node.StartPoint().Row) + 1
	col := int(node.StartPoint().Column)
	handler := &graph.Node{
		ID:       NodeID(graph.NodeHTTPEndpoint, "Route", file, line, col),
		Kind:     graph.NodeHTTPEndpoint,
		Name:     "Route",
		File:     file,
		Line:     line,
		Column:   col,
		Language: "typescript",
		Annotations: make(map[string]bool),
		Properties:  map[string]string{
			"component": "true",
		},
	}
	if route != "" {
		handler.Properties["route"] = route
		handler.Route = route
	}
	if componentName != "" {
		handler.Properties["component_name"] = componentName
	}
	result.HTTPHandlers = append(result.HTTPHandlers, handler)
}

// stripQuotes removes surrounding quotes (single, double, or backtick) from a string.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		first := s[0]
		last := s[len(s)-1]
		if (first == '"' && last == '"') || (first == '\'' && last == '\'') || (first == '`' && last == '`') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

// extractParamTypes extracts type annotation strings from a formal_parameters node.
func extractParamTypes(params *sitter.Node, src []byte) []string {
	var types []string
	for i := 0; i < int(params.ChildCount()); i++ {
		param := params.Child(i)
		if param == nil {
			continue
		}
		switch param.Type() {
		case "required_parameter", "optional_parameter":
			typeNode := param.ChildByFieldName("type")
			if typeNode != nil {
				// type_annotation contains ": Type", extract just the type part
				typeText := typeNode.Content(src)
				typeText = strings.TrimPrefix(typeText, ":")
				typeText = strings.TrimSpace(typeText)
				if typeText != "" {
					types = append(types, typeText)
				}
			} else {
				// No type annotation, use the parameter name as fallback
				nameNode := param.ChildByFieldName("pattern")
				if nameNode != nil {
					types = append(types, nameNode.Content(src))
				}
			}
		}
	}
	return types
}

// findArrowFunction finds an arrow_function node within a node tree (handles type assertions, etc.)
func findArrowFunction(node *sitter.Node) *sitter.Node {
	if node.Type() == "arrow_function" {
		return node
	}
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil && child.Type() == "arrow_function" {
			return child
		}
	}
	return nil
}

// extractTSParamNames extracts parameter names from a formal_parameters node.
func extractTSParamNames(params *sitter.Node, src []byte) []string {
	var names []string
	for i := 0; i < int(params.ChildCount()); i++ {
		param := params.Child(i)
		if param == nil {
			continue
		}
		switch param.Type() {
		case "required_parameter", "optional_parameter":
			nameNode := param.ChildByFieldName("pattern")
			if nameNode != nil {
				names = append(names, nameNode.Content(src))
			}
		}
	}
	return names
}

// ---------------------------------------------------------------------------
// Intraprocedural data flow analysis
// ---------------------------------------------------------------------------

// analyzeFunctionBody creates variable/parameter nodes and data flow edges
// for a single TypeScript function body. One SymbolTable + FlowBuilder per function.
func (tp *TypeScriptParser) analyzeFunctionBody(body *sitter.Node, src []byte, file string, fn *graph.Node, result *ParseResult) {
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
			Language: "typescript",
		}
		fb.AddParameter(paramNode)
		fb.AddDeclares(fn.ID, paramID)
		st.Define(pName, paramID)
	}

	// Walk body statements sequentially
	for i := 0; i < int(body.ChildCount()); i++ {
		child := body.Child(i)
		if child == nil {
			continue
		}
		tp.analyzeStatementDF(child, src, file, fn, st, fb)
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
	tp.buildCFG(body, src, file, fn, posMap, result)
}

// analyzeStatementDF dispatches to handlers based on tree-sitter node type.
func (tp *TypeScriptParser) analyzeStatementDF(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	switch node.Type() {
	case "lexical_declaration":
		tp.analyzeLexicalDeclDF(node, src, file, fn, st, fb)
	case "variable_declaration":
		tp.analyzeLexicalDeclDF(node, src, file, fn, st, fb)
	case "expression_statement":
		// Check for assignment_expression or standalone call_expression
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child == nil {
				continue
			}
			switch child.Type() {
			case "assignment_expression":
				tp.analyzeAssignmentDF(child, src, file, fn, st, fb)
			case "call_expression":
				tp.analyzeCallArgsDF(child, src, file, st, fb)
			}
		}
	case "return_statement":
		tp.analyzeReturnDF(node, src, file, fn, st, fb)
	default:
		// Recurse into block-level structures (if, for, while, etc.)
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child == nil {
				continue
			}
			tp.analyzeStatementDF(child, src, file, fn, st, fb)
			if fb.VariableCount() >= dataflow.MaxVariablesPerFunction {
				return
			}
		}
	}
}

// analyzeLexicalDeclDF handles `const x = expr` / `let x = expr` / `var x = expr`.
// These are lexical_declaration or variable_declaration containing variable_declarator children.
func (tp *TypeScriptParser) analyzeLexicalDeclDF(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child == nil || child.Type() != "variable_declarator" {
			continue
		}

		nameNode := child.ChildByFieldName("name")
		if nameNode == nil {
			continue
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
			Language: "typescript",
		}
		fb.AddVariable(varNode)
		st.Define(name, varID)

		// Resolve RHS source and create assigns edge
		valueNode := child.ChildByFieldName("value")
		if valueNode != nil {
			sourceID := tp.resolveRHSSourceDF(valueNode, src, file, fn, st, fb)
			if sourceID != "" {
				fb.AddAssign(sourceID, varID)
			}
			// Emit reads edges for binary expressions
			tp.emitReadsDF(valueNode, src, varID, st, fb)
		}
	}
}

// analyzeAssignmentDF handles `x = expr` assignment expressions.
func (tp *TypeScriptParser) analyzeAssignmentDF(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
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

	sourceID := tp.resolveRHSSourceDF(right, src, file, fn, st, fb)
	if sourceID != "" {
		fb.AddAssign(sourceID, targetID)
	}
	tp.emitReadsDF(right, src, targetID, st, fb)
}

// resolveRHSSourceDF resolves the primary data source from an RHS expression.
func (tp *TypeScriptParser) resolveRHSSourceDF(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) string {
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
		// Also process call arguments
		tp.analyzeCallArgsDF(node, src, file, st, fb)
		return callID

	case "identifier":
		name := node.Content(src)
		if varID, ok := st.Resolve(name); ok {
			return varID
		}
		return ""

	case "member_expression":
		return tp.resolveMemberExprDF(node, src, file, fn, st, fb)

	case "subscript_expression":
		// e.g., payload["name"]: resolve the object
		obj := node.ChildByFieldName("object")
		if obj != nil && obj.Type() == "identifier" {
			name := obj.Content(src)
			if varID, ok := st.Resolve(name); ok {
				return varID
			}
		}
		return ""

	case "binary_expression":
		// For binary expressions like "str" + name, resolve first meaningful operand
		left := node.ChildByFieldName("left")
		right := node.ChildByFieldName("right")
		if leftID := tp.resolveRHSSourceDF(left, src, file, fn, st, fb); leftID != "" {
			return leftID
		}
		return tp.resolveRHSSourceDF(right, src, file, fn, st, fb)

	case "template_string":
		// Template literals: walk children for template_substitution
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child != nil && child.Type() == "template_substitution" {
				for j := 0; j < int(child.ChildCount()); j++ {
					inner := child.Child(j)
					if inner != nil {
						if id := tp.resolveRHSSourceDF(inner, src, file, fn, st, fb); id != "" {
							return id
						}
					}
				}
			}
		}
		return ""

	default:
		return ""
	}
}

// resolveMemberExprDF handles member_expression chains like user.email or req.body.
// Creates a single synthetic field variable with the full text and emits one field_access
// edge from the root identifier to the synthetic variable.
func (tp *TypeScriptParser) resolveMemberExprDF(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) string {
	if node == nil || node.Type() != "member_expression" {
		return ""
	}

	fullText := node.Content(src)

	// Walk to the leftmost identifier to find the root
	var rootID string
	current := node
	for {
		obj := current.ChildByFieldName("object")
		if obj == nil {
			break
		}
		if obj.Type() == "identifier" {
			rootName := obj.Content(src)
			if id, ok := st.Resolve(rootName); ok {
				rootID = id
			}
			break
		} else if obj.Type() == "member_expression" {
			current = obj
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
		Language: "typescript",
	}
	fb.AddVariable(fieldVarNode)
	fb.AddFieldAccess(rootID, fieldVarID)

	return fieldVarID
}

// analyzeCallArgsDF processes arguments to a call expression, creating passes_to edges.
func (tp *TypeScriptParser) analyzeCallArgsDF(callNode *sitter.Node, src []byte, file string, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
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
		case "member_expression":
			// e.g., user.name: passes_to from root identifier
			tp.passMemberToCall(arg, src, callID, st, fb)
		}
	}
}

// passMemberToCall finds the root identifier of a member expression
// and emits a passes_to edge.
func (tp *TypeScriptParser) passMemberToCall(node *sitter.Node, src []byte, callID string, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	if node == nil {
		return
	}
	switch node.Type() {
	case "identifier":
		name := node.Content(src)
		if varID, ok := st.Resolve(name); ok {
			fb.AddPassesTo(varID, callID)
		}
	case "member_expression":
		obj := node.ChildByFieldName("object")
		if obj != nil {
			tp.passMemberToCall(obj, src, callID, st, fb)
		}
	}
}

// analyzeReturnDF handles `return expr`.
func (tp *TypeScriptParser) analyzeReturnDF(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
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
		case "member_expression":
			rootID := tp.resolveMemberExprDF(child, src, file, fn, st, fb)
			if rootID != "" {
				fb.AddReturn(rootID, fn.ID)
			}
		default:
			tp.walkReturnExprDF(child, src, file, fn, st, fb)
		}
	}
}

// walkReturnExprDF recursively finds identifiers in return expressions.
func (tp *TypeScriptParser) walkReturnExprDF(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
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
	if node.Type() == "member_expression" {
		rootID := tp.resolveMemberExprDF(node, src, file, fn, st, fb)
		if rootID != "" {
			fb.AddReturn(rootID, fn.ID)
		}
		return
	}
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil {
			tp.walkReturnExprDF(child, src, file, fn, st, fb)
		}
	}
}

// emitReadsDF walks a compound expression (binary_expression) to find all
// identifier references and emits reads edges.
func (tp *TypeScriptParser) emitReadsDF(node *sitter.Node, src []byte, targetID string, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	if node == nil {
		return
	}
	if node.Type() != "binary_expression" {
		return
	}
	tp.walkForReadsDF(node, src, targetID, st, fb)
}

// walkForReadsDF recursively finds identifiers in an expression tree and emits reads edges.
func (tp *TypeScriptParser) walkForReadsDF(node *sitter.Node, src []byte, targetID string, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
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
			tp.walkForReadsDF(child, src, targetID, st, fb)
		}
	}
}

// ---------------------------------------------------------------------------
// Control flow graph construction
// ---------------------------------------------------------------------------

// buildCFG constructs a control flow graph for a function body.
// Creates entry/exit blocks, connects control flow edges, and populates
// block members with B1 node IDs.
func (tp *TypeScriptParser) buildCFG(body *sitter.Node, src []byte, file string, fn *graph.Node, posMap map[posKey][]string, result *ParseResult) {
	cb := dataflow.NewCFGBuilder(fn.ID, file)

	// Create entry and exit blocks
	entry := cb.NewBlock("entry", fn.Line)
	exit := cb.NewBlock("exit", 0)

	// Connect function to entry block
	cb.AddEdge(fn.ID, entry.ID, "entry")

	// Add parameters to entry block
	for _, pName := range fn.ParamNames {
		if pName == "_" {
			continue
		}
		paramID := NodeID(graph.NodeParameter, pName, file, fn.Line, 0)
		cb.AddMember(entry, paramID)
	}

	// Process function body statements
	lastBlock := tp.buildCFGStatementsTS(body, src, file, fn, entry, exit, posMap, cb)

	// Connect last block to exit (if it didn't terminate early)
	if lastBlock != nil {
		cb.AddEdge(lastBlock.ID, exit.ID, "exit")
	}

	// Merge results into ParseResult
	blocks, cfEdges := cb.Result()
	result.BasicBlocks = append(result.BasicBlocks, blocks...)
	result.Edges = append(result.Edges, cfEdges...)
}

// buildCFGStatementsTS walks container children and builds CFG for TypeScript statements.
// Returns the last open block (nil if all paths terminated).
func (tp *TypeScriptParser) buildCFGStatementsTS(container *sitter.Node, src []byte, file string, fn *graph.Node, current, exit *graph.Node, posMap map[posKey][]string, cb *dataflow.CFGBuilder) *graph.Node {
	if container == nil {
		return current
	}

	for i := 0; i < int(container.ChildCount()); i++ {
		if cb.BlockCount() >= dataflow.MaxBlocksPerFunction {
			fn.Annotations["cfg:truncated"] = true
			return current
		}

		child := container.Child(i)
		if child == nil {
			continue
		}

		// If current is nil, control flow has terminated, skip remaining statements
		if current == nil {
			return nil
		}

		switch child.Type() {
		case "if_statement":
			current = tp.buildCFGIfTS(child, src, file, fn, current, exit, posMap, cb)

		case "for_statement", "for_in_statement", "while_statement":
			current = tp.buildCFGForTS(child, src, file, fn, current, exit, posMap, cb)

		case "switch_statement":
			current = tp.buildCFGSwitchTS(child, src, file, fn, current, exit, posMap, cb)

		case "return_statement":
			// Collect node IDs from return statement and add to current block
			nodeIDs := collectNodeIDs(child, posMap)
			for _, id := range nodeIDs {
				cb.AddMember(current, id)
			}
			// Connect to exit and terminate this path
			cb.AddEdge(current.ID, exit.ID, "exit")
			return nil

		case "throw_statement":
			// Collect node IDs from throw statement and add to current block
			nodeIDs := collectNodeIDs(child, posMap)
			for _, id := range nodeIDs {
				cb.AddMember(current, id)
			}
			// throw terminates control flow
			cb.AddEdge(current.ID, exit.ID, "exit")
			return nil

		default:
			// Sequential statement: add node IDs to current block
			nodeIDs := collectNodeIDs(child, posMap)
			for _, id := range nodeIDs {
				cb.AddMember(current, id)
			}
		}
	}

	return current
}

// buildCFGIfTS handles if_statement nodes in TypeScript.
// Creates condition block, then-block, else-block (if present), and merge block.
// Returns merge block or nil if all branches terminate.
func (tp *TypeScriptParser) buildCFGIfTS(node *sitter.Node, src []byte, file string, fn *graph.Node, current, exit *graph.Node, posMap map[posKey][]string, cb *dataflow.CFGBuilder) *graph.Node {
	if current == nil {
		return nil
	}

	// Condition evaluation happens in current block
	condition := node.ChildByFieldName("condition")
	if condition != nil {
		nodeIDs := collectNodeIDs(condition, posMap)
		for _, id := range nodeIDs {
			cb.AddMember(current, id)
		}
	}

	// Create then-block
	line := int(node.StartPoint().Row) + 1
	thenBlock := cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), line)
	cb.AddEdge(current.ID, thenBlock.ID, "true_branch")

	// Process then-branch (consequence)
	consequence := node.ChildByFieldName("consequence")
	thenEnd := tp.buildCFGStatementsTS(consequence, src, file, fn, thenBlock, exit, posMap, cb)

	// Check for else-branch
	alternative := node.ChildByFieldName("alternative")
	var elseEnd *graph.Node

	if alternative != nil {
		// Check if alternative is else_clause containing an if_statement (else-if)
		if alternative.Type() == "else_clause" {
			// Look for if_statement child
			var elseIfNode *sitter.Node
			for i := 0; i < int(alternative.ChildCount()); i++ {
				child := alternative.Child(i)
				if child != nil && child.Type() == "if_statement" {
					elseIfNode = child
					break
				}
			}

			if elseIfNode != nil {
				// else if: create block and recurse
				elseIfBlock := cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), int(elseIfNode.StartPoint().Row)+1)
				cb.AddEdge(current.ID, elseIfBlock.ID, "false_branch")
				elseEnd = tp.buildCFGIfTS(elseIfNode, src, file, fn, elseIfBlock, exit, posMap, cb)
			} else {
				// Regular else clause
				elseBlock := cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), int(alternative.StartPoint().Row)+1)
				cb.AddEdge(current.ID, elseBlock.ID, "false_branch")
				elseEnd = tp.buildCFGStatementsTS(alternative, src, file, fn, elseBlock, exit, posMap, cb)
			}
		} else {
			// Unexpected alternative structure, treat as regular block
			elseBlock := cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), int(alternative.StartPoint().Row)+1)
			cb.AddEdge(current.ID, elseBlock.ID, "false_branch")
			elseEnd = tp.buildCFGStatementsTS(alternative, src, file, fn, elseBlock, exit, posMap, cb)
		}
	} else {
		// No else branch, false_branch goes to merge
		elseEnd = current
	}

	// Create merge block if at least one branch continues
	if thenEnd == nil && elseEnd == nil {
		// Both branches terminated
		return nil
	}

	mergeBlock := cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), line)

	if thenEnd != nil {
		cb.AddEdge(thenEnd.ID, mergeBlock.ID, "fallthrough")
	}
	if elseEnd != nil && elseEnd != current {
		// Only add edge if elseEnd is not the condition block
		cb.AddEdge(elseEnd.ID, mergeBlock.ID, "fallthrough")
	} else if elseEnd == current {
		// No else branch: connect current to merge via false_branch
		cb.AddEdge(current.ID, mergeBlock.ID, "false_branch")
	}

	return mergeBlock
}

// buildCFGForTS handles for_statement, for_in_statement, and while_statement nodes.
// Creates header block (condition), body block, loop_back and loop_exit edges.
// Returns after-loop block.
func (tp *TypeScriptParser) buildCFGForTS(node *sitter.Node, src []byte, file string, fn *graph.Node, current, exit *graph.Node, posMap map[posKey][]string, cb *dataflow.CFGBuilder) *graph.Node {
	if current == nil {
		return nil
	}

	line := int(node.StartPoint().Row) + 1

	// Create header block (condition evaluation)
	headerBlock := cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), line)

	// For for_in_statement, condition is the iterator expression
	// For for_statement, condition is in the initializer/condition/update
	// For while_statement, condition is explicit
	var conditionNode *sitter.Node
	switch node.Type() {
	case "while_statement":
		conditionNode = node.ChildByFieldName("condition")
	case "for_in_statement":
		conditionNode = node.ChildByFieldName("right")
	case "for_statement":
		conditionNode = node.ChildByFieldName("condition")
	}

	if conditionNode != nil {
		nodeIDs := collectNodeIDs(conditionNode, posMap)
		for _, id := range nodeIDs {
			cb.AddMember(headerBlock, id)
		}
	}

	// Connect current to header
	cb.AddEdge(current.ID, headerBlock.ID, "fallthrough")

	// Create body block
	bodyBlock := cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), line)
	cb.AddEdge(headerBlock.ID, bodyBlock.ID, "true_branch")

	// Process loop body
	body := node.ChildByFieldName("body")
	bodyEnd := tp.buildCFGStatementsTS(body, src, file, fn, bodyBlock, exit, posMap, cb)

	// Loop back from body end to header
	if bodyEnd != nil {
		cb.AddEdge(bodyEnd.ID, headerBlock.ID, "loop_back")
	}

	// Create after-loop block (loop exit)
	afterBlock := cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), line)
	cb.AddEdge(headerBlock.ID, afterBlock.ID, "loop_exit")

	return afterBlock
}

// buildCFGSwitchTS handles switch_statement nodes in TypeScript.
// Creates one case block per case, connects with true_branch/false_branch, handles breaks.
// Returns merge block.
func (tp *TypeScriptParser) buildCFGSwitchTS(node *sitter.Node, src []byte, file string, fn *graph.Node, current, exit *graph.Node, posMap map[posKey][]string, cb *dataflow.CFGBuilder) *graph.Node {
	if current == nil {
		return nil
	}

	line := int(node.StartPoint().Row) + 1

	// Switch value evaluation happens in current block
	value := node.ChildByFieldName("value")
	if value != nil {
		nodeIDs := collectNodeIDs(value, posMap)
		for _, id := range nodeIDs {
			cb.AddMember(current, id)
		}
	}

	// Find the switch body
	var switchBody *sitter.Node
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil && child.Type() == "switch_body" {
			switchBody = child
			break
		}
	}

	if switchBody == nil {
		return current
	}

	// Create merge block
	mergeBlock := cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), line)

	// Process each case
	hasDefault := false
	var caseEnds []*graph.Node

	for i := 0; i < int(switchBody.ChildCount()); i++ {
		child := switchBody.Child(i)
		if child == nil {
			continue
		}

		var caseBlock *graph.Node
		var label string

		if child.Type() == "switch_case" {
			caseBlock = cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), int(child.StartPoint().Row)+1)
			label = "true_branch"
		} else if child.Type() == "switch_default" {
			caseBlock = cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), int(child.StartPoint().Row)+1)
			label = "false_branch"
			hasDefault = true
		} else {
			continue
		}

		cb.AddEdge(current.ID, caseBlock.ID, label)

		// Process case body (everything except the case/default keyword)
		// Walk children to find statements
		caseEnd := caseBlock
		for j := 0; j < int(child.ChildCount()); j++ {
			stmt := child.Child(j)
			if stmt == nil {
				continue
			}

			// Skip the case value and colon
			if stmt.Type() == ":" || stmt.Type() == "case" || stmt.Type() == "default" {
				continue
			}

			// Check for break statement
			if stmt.Type() == "break_statement" {
				nodeIDs := collectNodeIDs(stmt, posMap)
				for _, id := range nodeIDs {
					cb.AddMember(caseEnd, id)
				}
				// break goes to merge
				cb.AddEdge(caseEnd.ID, mergeBlock.ID, "fallthrough")
				caseEnd = nil
				break
			}

			// Regular statement
			if caseEnd != nil {
				nodeIDs := collectNodeIDs(stmt, posMap)
				for _, id := range nodeIDs {
					cb.AddMember(caseEnd, id)
				}
			}
		}

		caseEnds = append(caseEnds, caseEnd)
	}

	// Connect case ends to merge block (fallthrough cases)
	for _, caseEnd := range caseEnds {
		if caseEnd != nil {
			cb.AddEdge(caseEnd.ID, mergeBlock.ID, "fallthrough")
		}
	}

	// If no default case, connect current to merge (implicit default)
	if !hasDefault {
		cb.AddEdge(current.ID, mergeBlock.ID, "false_branch")
	}

	return mergeBlock
}

