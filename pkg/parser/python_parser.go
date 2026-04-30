package parser

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/python"

	"github.com/ugiordan/architecture-analyzer/pkg/dataflow"
	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

// PythonSkipDirs lists directories that should be skipped when scanning Python projects.
var PythonSkipDirs = []string{"__pycache__", "migrations", ".tox", "venv", "site-packages", ".venv", "env"}

// PythonTestPatterns lists filename patterns that identify test files.
var PythonTestPatterns = []string{"test_*", "*_test.py"}

// pythonBuiltins is the set of PascalCase names that are builtins or standard
// exception types, not user-defined class instantiations.
var pythonBuiltins = map[string]bool{
	"True": true, "False": true, "None": true,
	"Exception": true, "BaseException": true, "ValueError": true,
	"TypeError": true, "KeyError": true, "IndexError": true,
	"AttributeError": true, "RuntimeError": true, "OSError": true,
	"IOError": true, "ImportError": true, "StopIteration": true,
	"NotImplementedError": true, "PermissionError": true,
	"FileNotFoundError": true, "ConnectionError": true,
	"TimeoutError": true, "UnicodeError": true,
}

// httpMethods is the set of decorator method names that indicate HTTP route handlers.
var httpMethods = map[string]bool{
	"route": true, "get": true, "post": true,
	"put": true, "delete": true, "patch": true,
}

// dbCallOps maps Python DB call patterns to their operation type (read/write).
var dbCallOps = map[string]string{
	"session.query":      "read",
	"session.execute":    "write",
	"session.add":        "write",
	"session.commit":     "write",
	"cursor.execute":     "write",
	"db.execute":         "write",
	"connection.execute": "write",
}

// PythonParser extracts code property graph nodes from Python source files using tree-sitter.
// Each goroutine MUST use its own PythonParser instance (tree-sitter parsers are not thread-safe).
type PythonParser struct {
	parser *sitter.Parser
}

// NewPythonParser creates a parser for Python source files backed by tree-sitter.
func NewPythonParser() *PythonParser {
	p := sitter.NewParser()
	p.SetLanguage(python.GetLanguage())
	return &PythonParser{parser: p}
}

func (pp *PythonParser) Language() string     { return "python" }
func (pp *PythonParser) Extensions() []string { return []string{".py"} }
func (pp *PythonParser) Clone() Parser {
	p := sitter.NewParser()
	p.SetLanguage(python.GetLanguage())
	return &PythonParser{parser: p}
}

// computePythonComplexity counts decision points in a Python function body.
// Complexity = 1 (base) + count of: if, elif, for, while, except, and, ternary, comprehension if.
func computePythonComplexity(node *sitter.Node) int {
	count := 1
	countPythonDecisionPoints(node, &count)
	return count
}

func countPythonDecisionPoints(node *sitter.Node, count *int) {
	if node == nil {
		return
	}
	switch node.Type() {
	case "if_statement":
		*count++
	case "elif_clause":
		*count++
	case "for_statement":
		*count++
	case "while_statement":
		*count++
	case "except_clause":
		*count++
	case "boolean_operator":
		// Only count "and" (short-circuit logical AND), not "or".
		// This aligns with Go/Rust/TypeScript which count only &&.
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child != nil && child.Type() == "and" {
				*count++
				break
			}
		}
	case "conditional_expression":
		*count++ // ternary: x if cond else y
	case "if_clause":
		*count++ // list comprehension if
	}
	for i := 0; i < int(node.ChildCount()); i++ {
		countPythonDecisionPoints(node.Child(i), count)
	}
}

// ParseFile parses a Python source file and returns extracted nodes and edges.
func (pp *PythonParser) ParseFile(path string, content []byte) (*ParseResult, error) {
	if len(content) > MaxFileSize {
		return nil, fmt.Errorf("file too large (%d bytes, max %d)", len(content), MaxFileSize)
	}
	tree, err := pp.parser.ParseCtx(context.Background(), nil, content)
	if err != nil {
		return nil, fmt.Errorf("tree-sitter parse failed: %w", err)
	}
	defer tree.Close()

	result := &ParseResult{}
	root := tree.RootNode()
	pp.walk(root, content, path, "", result)
	return result, nil
}

// walk recursively traverses the AST. className tracks the enclosing class for method extraction.
func (pp *PythonParser) walk(node *sitter.Node, src []byte, file, className string, result *ParseResult) {
	switch node.Type() {
	case "function_definition":
		pp.extractFunction(node, src, file, className, nil, result)
		return // children handled inside extractFunction
	case "decorated_definition":
		pp.extractDecorated(node, src, file, className, result)
		return // children handled inside extractDecorated
	case "class_definition":
		pp.extractClass(node, src, file, result)
		return // children handled inside extractClass
	case "call":
		pp.extractCallSite(node, src, file, result)
	}

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil {
			pp.walk(child, src, file, className, result)
		}
	}
}

// extractClass processes a class_definition node, walking its body with the class name set.
func (pp *PythonParser) extractClass(node *sitter.Node, src []byte, file string, result *ParseResult) {
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
			pp.walk(child, src, file, clsName, result)
		}
	}
}

// extractDecorated handles a decorated_definition node: collects decorators, then
// delegates to extractFunction for the inner function_definition.
func (pp *PythonParser) extractDecorated(node *sitter.Node, src []byte, file, className string, result *ParseResult) {
	var decorators []string
	var fnNode *sitter.Node

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}
		switch child.Type() {
		case "decorator":
			decorators = append(decorators, child.Content(src))
		case "function_definition":
			fnNode = child
		}
	}

	if fnNode != nil {
		pp.extractFunction(fnNode, src, file, className, decorators, result)
	}
}

// extractFunction creates a Function node and checks decorators for HTTP route patterns.
func (pp *PythonParser) extractFunction(node *sitter.Node, src []byte, file, className string, decorators []string, result *ParseResult) {
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
		Language:    "python",
		TypeName:    className,
		Decorators:  decorators,
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	// Compute complexity from function body
	body := node.ChildByFieldName("body")
	if body != nil {
		fn.Complexity = computePythonComplexity(body)
	}
	result.Functions = append(result.Functions, fn)

	// Extract parameter names
	params := node.ChildByFieldName("parameters")
	if params != nil {
		for i := 0; i < int(params.ChildCount()); i++ {
			param := params.Child(i)
			if param == nil {
				continue
			}
			var paramName string
			switch param.Type() {
			case "identifier":
				paramName = param.Content(src)
			case "typed_parameter", "typed_default_parameter", "default_parameter":
				nameNode := param.ChildByFieldName("name")
				if nameNode == nil {
					// fallback: first identifier child
					for j := 0; j < int(param.ChildCount()); j++ {
						ch := param.Child(j)
						if ch != nil && ch.Type() == "identifier" {
							paramName = ch.Content(src)
							break
						}
					}
				} else {
					paramName = nameNode.Content(src)
				}
			case "list_splat_pattern", "dictionary_splat_pattern":
				for j := 0; j < int(param.ChildCount()); j++ {
					ch := param.Child(j)
					if ch != nil && ch.Type() == "identifier" {
						paramName = ch.Content(src)
						break
					}
				}
			}
			if paramName != "" && paramName != "self" && paramName != "cls" {
				fn.ParamNames = append(fn.ParamNames, paramName)
			}
		}
	}

	// Check decorators for HTTP route patterns
	for _, dec := range decorators {
		pp.maybeExtractHTTPHandler(dec, fn, file, result)
	}

	// Extract intraprocedural data flow from function body
	if body != nil {
		pp.analyzeFunctionBody(body, src, file, fn, result)
	}

	// Walk function body for call sites, etc.
	if body != nil {
		for i := 0; i < int(body.ChildCount()); i++ {
			child := body.Child(i)
			if child != nil {
				pp.walk(child, src, file, className, result)
			}
		}
	}
}

// maybeExtractHTTPHandler checks if a decorator string matches an HTTP route pattern
// like @app.route("/path"), @app.get("/path"), @router.post("/path"), etc.
func (pp *PythonParser) maybeExtractHTTPHandler(decorator string, fn *graph.Node, file string, result *ParseResult) {
	// Strip the leading @
	dec := strings.TrimPrefix(decorator, "@")

	// We need: <identifier>.<method>(...) where method is in httpMethods
	// Find the method call part
	parenIdx := strings.Index(dec, "(")
	if parenIdx < 0 {
		return
	}
	callPart := dec[:parenIdx]
	dotIdx := strings.LastIndex(callPart, ".")
	if dotIdx < 0 {
		return
	}
	method := callPart[dotIdx+1:]
	if !httpMethods[method] {
		return
	}

	// Extract route path from first string argument
	argPart := dec[parenIdx+1:]
	route := extractStringArg(argPart)

	handler := &graph.Node{
		ID:          NodeID(graph.NodeHTTPEndpoint, fn.Name, file, fn.Line, fn.Column),
		Kind:        graph.NodeHTTPEndpoint,
		Name:        fn.Name,
		File:        file,
		Line:        fn.Line,
		Column:      fn.Column,
		Language:    "python",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	if route != "" {
		handler.Properties["route"] = route
		handler.Route = route
	}
	httpMethod := strings.ToUpper(method)
	handler.Properties["method"] = httpMethod
	handler.HTTPMethod = httpMethod
	result.HTTPHandlers = append(result.HTTPHandlers, handler)
}

// extractStringArg extracts the first quoted string from a decorator argument list.
// Input example: `"/users", methods=["GET"])`
func extractStringArg(s string) string {
	for _, q := range []byte{'"', '\''} {
		start := strings.IndexByte(s, q)
		if start < 0 {
			continue
		}
		end := strings.IndexByte(s[start+1:], q)
		if end < 0 {
			continue
		}
		return s[start+1 : start+1+end]
	}
	return ""
}

// classifyDBCall checks if a call text matches a known DB operation pattern.
// It checks the exact text first, then for dotted calls checks the last two segments.
func classifyDBCall(callText string) (op string, matched bool) {
	if op, ok := dbCallOps[callText]; ok {
		return op, true
	}
	if strings.Contains(callText, ".") {
		parts := strings.Split(callText, ".")
		if len(parts) >= 2 {
			suffix := parts[len(parts)-2] + "." + parts[len(parts)-1]
			if op, ok := dbCallOps[suffix]; ok {
				return op, true
			}
		}
	}
	return "", false
}

// extractCallSite creates a CallSite node from a call expression, and detects
// DB operations and class instantiations (struct literals).
func (pp *PythonParser) extractCallSite(node *sitter.Node, src []byte, file string, result *ParseResult) {
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
		Language:   "python",
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
			if arg.Type() == "string" {
				// Extract string_content from string node
				for j := 0; j < int(arg.ChildCount()); j++ {
					ch := arg.Child(j)
					if ch != nil && ch.Type() == "string_content" {
						stringArgs = append(stringArgs, ch.Content(src))
					}
				}
				if len(stringArgs) == 0 {
					// Fallback: strip quotes from full string content
					text := arg.Content(src)
					text = strings.Trim(text, "\"'")
					if text != "" {
						stringArgs = append(stringArgs, text)
					}
				}
			}
		}
		if len(stringArgs) > 0 {
			cs.Properties["string_args"] = strings.Join(stringArgs, ",")
		}
	}

	// Check for DB operations
	if op, ok := classifyDBCall(callText); ok {
		dbOp := &graph.Node{
			ID:         NodeID(graph.NodeDBOperation, callText, file, line, col),
			Kind:       graph.NodeDBOperation,
			Name:       callText,
			File:       file,
			Line:       line,
			Column:     col,
			Language:   "python",
			Properties: map[string]string{"operation": op},
			Operation:  op,
		}
		result.DBOperations = append(result.DBOperations, dbOp)
		return
	}

	// For non-DB, non-dotted calls, check if it's a PascalCase class instantiation
	if !strings.Contains(callText, ".") {
		pp.maybeExtractStructLiteral(callText, node, src, file, line, col, result)
	}
}

// maybeExtractStructLiteral checks if a simple (non-dotted) call name looks like
// a PascalCase class instantiation and, if so, creates a StructLiteral node.
func (pp *PythonParser) maybeExtractStructLiteral(name string, node *sitter.Node, src []byte, file string, line, col int, result *ParseResult) {
	if len(name) == 0 {
		return
	}
	// Must start with an uppercase letter
	if !unicode.IsUpper(rune(name[0])) {
		return
	}
	// Skip Python builtins
	if pythonBuiltins[name] {
		return
	}

	sl := &graph.Node{
		ID:         NodeID(graph.NodeStructLiteral, name, file, line, col),
		Kind:       graph.NodeStructLiteral,
		Name:       name,
		File:       file,
		Line:       line,
		Column:     col,
		Language:   "python",
		Properties: make(map[string]string),
	}
	result.StructLiterals = append(result.StructLiterals, sl)
}

// ---------------------------------------------------------------------------
// Intraprocedural data flow analysis
// ---------------------------------------------------------------------------

// analyzeFunctionBody creates variable/parameter nodes and data flow edges
// for a single Python function body. One SymbolTable + FlowBuilder per function.
func (pp *PythonParser) analyzeFunctionBody(body *sitter.Node, src []byte, file string, fn *graph.Node, result *ParseResult) {
	st := dataflow.NewSymbolTable()
	fb := dataflow.NewFlowBuilder()

	// Create parameter nodes and register them in the symbol table.
	// ParamNames already excludes self/cls (handled in extractFunction).
	for _, pName := range fn.ParamNames {
		paramID := NodeID(graph.NodeParameter, pName, file, fn.Line, 0)
		paramNode := &graph.Node{
			ID:       paramID,
			Kind:     graph.NodeParameter,
			Name:     pName,
			File:     file,
			Line:     fn.Line,
			Language: "python",
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
		pp.analyzeStatementDF(child, src, file, fn, st, fb)
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
	blocks, cfEdges := pp.buildCFG(body, src, file, fn, posMap)
	result.BasicBlocks = append(result.BasicBlocks, blocks...)
	result.Edges = append(result.Edges, cfEdges...)
}

// analyzeStatementDF dispatches to handlers based on tree-sitter node type.
func (pp *PythonParser) analyzeStatementDF(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	switch node.Type() {
	case "assignment":
		pp.analyzeAssignmentDF(node, src, file, fn, st, fb)
	case "expression_statement":
		// Assignments in Python are wrapped in expression_statement
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child == nil {
				continue
			}
			switch child.Type() {
			case "assignment":
				pp.analyzeAssignmentDF(child, src, file, fn, st, fb)
			case "call":
				pp.analyzeCallArgsDF(child, src, file, st, fb)
			}
		}
	case "return_statement":
		pp.analyzeReturnDF(node, src, file, fn, st, fb)
	default:
		// Recurse into block-level structures (if, for, while, etc.)
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child == nil {
				continue
			}
			pp.analyzeStatementDF(child, src, file, fn, st, fb)
			if fb.VariableCount() >= dataflow.MaxVariablesPerFunction {
				return
			}
		}
	}
}

// analyzeAssignmentDF handles `x = expr` assignments.
func (pp *PythonParser) analyzeAssignmentDF(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	left := node.ChildByFieldName("left")
	right := node.ChildByFieldName("right")
	if left == nil {
		return
	}

	// Get LHS variable name
	var targetName string
	if left.Type() == "identifier" {
		targetName = left.Content(src)
	} else if left.Type() == "pattern_list" {
		// Tuple unpacking: take the first identifier
		for i := 0; i < int(left.ChildCount()); i++ {
			child := left.Child(i)
			if child != nil && child.Type() == "identifier" {
				targetName = child.Content(src)
				break
			}
		}
	}

	if targetName == "" {
		return
	}

	// Python assignments always create/rebind, so create a new variable node
	line := int(left.StartPoint().Row) + 1
	col := int(left.StartPoint().Column)
	varID := NodeID(graph.NodeVariable, targetName, file, line, col)
	varNode := &graph.Node{
		ID:       varID,
		Kind:     graph.NodeVariable,
		Name:     targetName,
		File:     file,
		Line:     line,
		Column:   col,
		Language: "python",
	}
	fb.AddVariable(varNode)
	st.Define(targetName, varID)

	// Resolve RHS source and create assigns edge
	if right != nil {
		sourceID := pp.resolveRHSSourceDF(right, src, file, fn, st, fb)
		if sourceID != "" {
			fb.AddAssign(sourceID, varID)
		}
		// Emit reads edges for binary expressions
		pp.emitReadsDF(right, src, varID, st, fb)
	}
}

// resolveRHSSourceDF resolves the primary data source from an RHS expression.
func (pp *PythonParser) resolveRHSSourceDF(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) string {
	if node == nil {
		return ""
	}

	switch node.Type() {
	case "call":
		// Return the call site ID (must match what extractCallSite generates)
		fnNode := node.ChildByFieldName("function")
		if fnNode == nil {
			return ""
		}
		callText := fnNode.Content(src)
		line := int(node.StartPoint().Row) + 1
		col := int(node.StartPoint().Column)
		callID := NodeID(graph.NodeCallSite, callText, file, line, col)
		// Also process call arguments
		pp.analyzeCallArgsDF(node, src, file, st, fb)
		return callID

	case "identifier":
		name := node.Content(src)
		if varID, ok := st.Resolve(name); ok {
			return varID
		}
		return ""

	case "attribute":
		return pp.resolveFieldAccessDF(node, src, file, fn, st, fb)

	case "subscript":
		// e.g., payload["name"]: resolve the object (whole collection taint)
		obj := node.ChildByFieldName("value")
		if obj != nil && obj.Type() == "identifier" {
			name := obj.Content(src)
			if varID, ok := st.Resolve(name); ok {
				return varID
			}
		}
		return ""

	case "binary_operator":
		// For binary ops like "str" + name, resolve first meaningful operand
		left := node.ChildByFieldName("left")
		right := node.ChildByFieldName("right")
		if leftID := pp.resolveRHSSourceDF(left, src, file, fn, st, fb); leftID != "" {
			return leftID
		}
		return pp.resolveRHSSourceDF(right, src, file, fn, st, fb)

	default:
		return ""
	}
}

// resolveFieldAccessDF handles attribute access chains like user.email or self.db.execute.
// Creates a single synthetic field variable with the full text and emits one field_access
// edge from the root identifier to the synthetic variable.
// The synthetic variable is NOT registered in the symbol table (per spec).
func (pp *PythonParser) resolveFieldAccessDF(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) string {
	if node == nil || node.Type() != "attribute" {
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
			// Skip "self" as root: resolve from the next level
			if rootName == "self" || rootName == "cls" {
				break
			}
			if id, ok := st.Resolve(rootName); ok {
				rootID = id
			}
			break
		} else if obj.Type() == "attribute" {
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
		Language: "python",
	}
	fb.AddVariable(fieldVarNode)
	fb.AddFieldAccess(rootID, fieldVarID)

	return fieldVarID
}

// analyzeCallArgsDF processes arguments to a call expression, creating passes_to edges.
func (pp *PythonParser) analyzeCallArgsDF(callNode *sitter.Node, src []byte, file string, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
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
		case "attribute":
			// e.g., user.name: passes_to from root identifier
			pp.passAttributeToCall(arg, src, callID, st, fb)
		}
	}
}

// passAttributeToCall finds the root identifier of an attribute expression
// and emits a passes_to edge.
func (pp *PythonParser) passAttributeToCall(node *sitter.Node, src []byte, callID string, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	if node == nil {
		return
	}
	switch node.Type() {
	case "identifier":
		name := node.Content(src)
		if name == "self" || name == "cls" {
			return
		}
		if varID, ok := st.Resolve(name); ok {
			fb.AddPassesTo(varID, callID)
		}
	case "attribute":
		obj := node.ChildByFieldName("object")
		if obj != nil {
			pp.passAttributeToCall(obj, src, callID, st, fb)
		}
	}
}

// analyzeReturnDF handles `return expr`.
func (pp *PythonParser) analyzeReturnDF(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
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
		case "attribute":
			rootID := pp.resolveFieldAccessDF(child, src, file, fn, st, fb)
			if rootID != "" {
				fb.AddReturn(rootID, fn.ID)
			}
		default:
			// Walk into compound expressions to find identifiers
			pp.walkReturnExprDF(child, src, file, fn, st, fb)
		}
	}
}

// walkReturnExprDF recursively finds identifiers in return expressions.
func (pp *PythonParser) walkReturnExprDF(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
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
	if node.Type() == "attribute" {
		rootID := pp.resolveFieldAccessDF(node, src, file, fn, st, fb)
		if rootID != "" {
			fb.AddReturn(rootID, fn.ID)
		}
		return
	}
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil {
			pp.walkReturnExprDF(child, src, file, fn, st, fb)
		}
	}
}

// emitReadsDF walks a compound expression (binary_operator) to find all
// identifier references and emits reads edges.
func (pp *PythonParser) emitReadsDF(node *sitter.Node, src []byte, targetID string, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	if node == nil {
		return
	}
	if node.Type() != "binary_operator" {
		return
	}
	pp.walkForReadsDF(node, src, targetID, st, fb)
}

// walkForReadsDF recursively finds identifiers in an expression tree and emits reads edges.
func (pp *PythonParser) walkForReadsDF(node *sitter.Node, src []byte, targetID string, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
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
			pp.walkForReadsDF(child, src, targetID, st, fb)
		}
	}
}

// ---------------------------------------------------------------------------
// Control Flow Graph construction
// ---------------------------------------------------------------------------

// buildCFG constructs a control flow graph for a Python function body.
// Creates entry/exit blocks, connects control flow edges, and populates
// block members with B1 node IDs.
func (pp *PythonParser) buildCFG(body *sitter.Node, src []byte, file string, fn *graph.Node, posMap map[posKey][]string) ([]*graph.Node, []*graph.Edge) {
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

	// Process function body statements (no loop context at top level)
	lastBlock := pp.buildCFGStatementsPy(body, src, file, fn, entry, exit, posMap, cb, nil)

	// Connect last block to exit (if it didn't terminate early)
	if lastBlock != nil {
		cb.AddEdge(lastBlock.ID, exit.ID, "exit")
	}

	return cb.Result()
}

// buildCFGStatementsPy walks container children and builds CFG for Python statements.
// Returns the last open block (nil if all paths terminated).
// lc is non-nil when inside a loop body (enables break/continue edges).
func (pp *PythonParser) buildCFGStatementsPy(container *sitter.Node, src []byte, file string, fn *graph.Node, current, exit *graph.Node, posMap map[posKey][]string, cb *dataflow.CFGBuilder, lc *loopContext) *graph.Node {
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

		if current == nil {
			return nil
		}

		switch child.Type() {
		case "if_statement":
			current = pp.buildCFGIfPy(child, src, file, fn, current, exit, posMap, cb, lc)

		case "for_statement", "while_statement":
			current = pp.buildCFGForPy(child, src, file, fn, current, exit, posMap, cb)

		case "return_statement":
			nodeIDs := collectNodeIDs(child, posMap)
			for _, id := range nodeIDs {
				cb.AddMember(current, id)
			}
			cb.AddEdge(current.ID, exit.ID, "exit")
			return nil

		case "raise_statement":
			nodeIDs := collectNodeIDs(child, posMap)
			for _, id := range nodeIDs {
				cb.AddMember(current, id)
			}
			cb.AddEdge(current.ID, exit.ID, "exit")
			return nil

		case "break_statement":
			nodeIDs := collectNodeIDs(child, posMap)
			for _, id := range nodeIDs {
				cb.AddMember(current, id)
			}
			if lc != nil {
				cb.AddEdge(current.ID, lc.afterBlock.ID, "fallthrough")
			}
			return nil

		case "continue_statement":
			nodeIDs := collectNodeIDs(child, posMap)
			for _, id := range nodeIDs {
				cb.AddMember(current, id)
			}
			if lc != nil {
				cb.AddEdge(current.ID, lc.headerBlock.ID, "loop_back")
			}
			return nil

		case "pass_statement":
			continue

		default:
			nodeIDs := collectNodeIDs(child, posMap)
			for _, id := range nodeIDs {
				cb.AddMember(current, id)
			}
		}
	}

	return current
}

// buildCFGIfPy handles if_statement nodes in Python.
// Creates condition block, then-block, elif/else-blocks (if present), and merge block.
// Returns merge block or nil if all branches terminate.
func (pp *PythonParser) buildCFGIfPy(node *sitter.Node, src []byte, file string, fn *graph.Node, current, exit *graph.Node, posMap map[posKey][]string, cb *dataflow.CFGBuilder, lc *loopContext) *graph.Node {
	if current == nil {
		return nil
	}

	condition := node.ChildByFieldName("condition")
	if condition != nil {
		nodeIDs := collectNodeIDs(condition, posMap)
		for _, id := range nodeIDs {
			cb.AddMember(current, id)
		}
	}

	line := int(node.StartPoint().Row) + 1
	thenBlock := cb.NewBlock("if-then", line)
	cb.AddEdge(current.ID, thenBlock.ID, "true_branch")

	consequence := node.ChildByFieldName("consequence")
	thenEnd := pp.buildCFGStatementsPy(consequence, src, file, fn, thenBlock, exit, posMap, cb, lc)

	// Track all elif/else branch endpoints for merge
	var branchEnds []*graph.Node
	if thenEnd != nil {
		branchEnds = append(branchEnds, thenEnd)
	}

	hasElse := false

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}

		if child.Type() == "elif_clause" {
			elifBlock := cb.NewBlock("elif", int(child.StartPoint().Row)+1)
			cb.AddEdge(current.ID, elifBlock.ID, "false_branch")

			elifCondition := child.ChildByFieldName("condition")
			if elifCondition != nil {
				nodeIDs := collectNodeIDs(elifCondition, posMap)
				for _, id := range nodeIDs {
					cb.AddMember(elifBlock, id)
				}
			}

			elifBodyBlock := cb.NewBlock("elif-body", int(child.StartPoint().Row)+1)
			cb.AddEdge(elifBlock.ID, elifBodyBlock.ID, "true_branch")

			elifConsequence := child.ChildByFieldName("consequence")
			elifEnd := pp.buildCFGStatementsPy(elifConsequence, src, file, fn, elifBodyBlock, exit, posMap, cb, lc)
			if elifEnd != nil {
				branchEnds = append(branchEnds, elifEnd)
			}

			// Chain: this elif's false_branch becomes next clause's source
			current = elifBlock

		} else if child.Type() == "else_clause" {
			elseBlock := cb.NewBlock("else", int(child.StartPoint().Row)+1)
			cb.AddEdge(current.ID, elseBlock.ID, "false_branch")

			elseBody := child.ChildByFieldName("body")
			elseEnd := pp.buildCFGStatementsPy(elseBody, src, file, fn, elseBlock, exit, posMap, cb, lc)
			if elseEnd != nil {
				branchEnds = append(branchEnds, elseEnd)
			}
			hasElse = true
		}
	}

	// CORR-001 fix: if no else clause, the last condition's false_branch must
	// connect to merge. Track current as a "pass-through" branch endpoint.
	if !hasElse {
		branchEnds = append(branchEnds, nil) // sentinel: need false_branch edge
	}

	// All branches terminated and no implicit fall-through needed
	if len(branchEnds) == 0 {
		return nil
	}

	// Lazy merge: only create if at least one branch continues
	hasLiveBranch := false
	for _, b := range branchEnds {
		if b != nil {
			hasLiveBranch = true
			break
		}
	}
	if !hasLiveBranch && hasElse {
		return nil
	}

	mergeBlock := cb.NewBlock("merge", line)

	for _, b := range branchEnds {
		if b != nil {
			cb.AddEdge(b.ID, mergeBlock.ID, "fallthrough")
		}
	}

	if !hasElse {
		// Connect last condition's false_branch to merge
		cb.AddEdge(current.ID, mergeBlock.ID, "false_branch")
	}

	return mergeBlock
}

// buildCFGForPy handles for_statement and while_statement nodes in Python.
// Creates header-block, body-block, loop_back, and loop_exit edges.
func (pp *PythonParser) buildCFGForPy(node *sitter.Node, src []byte, file string, fn *graph.Node, current, exit *graph.Node, posMap map[posKey][]string, cb *dataflow.CFGBuilder) *graph.Node {
	if current == nil {
		return nil
	}

	line := int(node.StartPoint().Row) + 1

	headerBlock := cb.NewBlock("loop-header", line)
	cb.AddEdge(current.ID, headerBlock.ID, "fallthrough")

	if node.Type() == "while_statement" {
		condition := node.ChildByFieldName("condition")
		if condition != nil {
			nodeIDs := collectNodeIDs(condition, posMap)
			for _, id := range nodeIDs {
				cb.AddMember(headerBlock, id)
			}
		}
	}

	if node.Type() == "for_statement" {
		left := node.ChildByFieldName("left")
		right := node.ChildByFieldName("right")
		if left != nil {
			nodeIDs := collectNodeIDs(left, posMap)
			for _, id := range nodeIDs {
				cb.AddMember(headerBlock, id)
			}
		}
		if right != nil {
			nodeIDs := collectNodeIDs(right, posMap)
			for _, id := range nodeIDs {
				cb.AddMember(headerBlock, id)
			}
		}
	}

	bodyBlock := cb.NewBlock("loop-body", line)
	cb.AddEdge(headerBlock.ID, bodyBlock.ID, "true_branch")

	afterBlock := cb.NewBlock("loop-after", line)
	cb.AddEdge(headerBlock.ID, afterBlock.ID, "loop_exit")

	// Pass loop context so break/continue in the body create correct edges
	lc := &loopContext{headerBlock: headerBlock, afterBlock: afterBlock}

	body := node.ChildByFieldName("body")
	bodyEnd := pp.buildCFGStatementsPy(body, src, file, fn, bodyBlock, exit, posMap, cb, lc)

	if bodyEnd != nil {
		cb.AddEdge(bodyEnd.ID, headerBlock.ID, "loop_back")
	}

	return afterBlock
}
