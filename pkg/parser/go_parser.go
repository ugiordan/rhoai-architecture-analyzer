package parser

import (
	"context"
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"

	"github.com/ugiordan/architecture-analyzer/pkg/dataflow"
	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

// MaxFileSize is the maximum file size (10 MB) that parsers will attempt to parse.
// Files larger than this are skipped to avoid excessive memory use in tree-sitter.
const MaxFileSize = 10 * 1024 * 1024

// GoParser extracts code property graph nodes from Go source files using tree-sitter.
// The underlying tree-sitter parser is NOT safe for concurrent use even with external
// locks: each goroutine MUST use its own GoParser instance (call NewGoParser per
// goroutine). Node IDs are deterministic content hashes, so no shared counter is needed.
type GoParser struct {
	parser *sitter.Parser
}

// NewGoParser creates a parser for Go source files backed by tree-sitter.
func NewGoParser() *GoParser {
	p := sitter.NewParser()
	p.SetLanguage(golang.GetLanguage())
	return &GoParser{parser: p}
}

func (gp *GoParser) Language() string     { return "go" }
func (gp *GoParser) Extensions() []string { return []string{".go"} }
func (gp *GoParser) Clone() Parser {
	p := sitter.NewParser()
	p.SetLanguage(golang.GetLanguage())
	return &GoParser{parser: p}
}

// computeComplexity counts decision points in a function body AST.
// Complexity = 1 (base) + count of: if, for/range, case (expression/type/default), &&.
func computeComplexity(node *sitter.Node) int {
	count := 1
	countDecisionPoints(node, &count)
	return count
}

func countDecisionPoints(node *sitter.Node, count *int) {
	if node == nil {
		return
	}
	switch node.Type() {
	case "if_statement":
		*count++
	case "for_statement":
		*count++
	case "expression_case", "type_case", "default_case":
		*count++
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
		child := node.Child(i)
		if child == nil {
			continue
		}
		// Skip "else if" chains: an if_statement that is the alternative branch
		// of a parent if_statement is not a new decision point.
		if node.Type() == "if_statement" && child.Type() == "if_statement" {
			// This is the else-if branch. Don't count it as a new if_statement,
			// but still recurse into its children.
			countDecisionPointsSkipSelf(child, count)
			continue
		}
		countDecisionPoints(child, count)
	}
}

// countDecisionPointsSkipSelf recurses into a node's children without counting
// the node itself as a decision point. Used for else-if chains.
func countDecisionPointsSkipSelf(node *sitter.Node, count *int) {
	if node == nil {
		return
	}
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil {
			countDecisionPoints(child, count)
		}
	}
}

// ParseFile parses a Go source file and returns extracted nodes and edges.
func (gp *GoParser) ParseFile(path string, content []byte) (*ParseResult, error) {
	if len(content) > MaxFileSize {
		return nil, fmt.Errorf("file too large (%d bytes, max %d)", len(content), MaxFileSize)
	}
	tree, err := gp.parser.ParseCtx(context.Background(), nil, content)
	if err != nil {
		return nil, fmt.Errorf("tree-sitter parse failed: %w", err)
	}
	defer tree.Close()

	result := &ParseResult{}
	root := tree.RootNode()
	gp.walk(root, content, path, result)
	return result, nil
}

func (gp *GoParser) walk(node *sitter.Node, src []byte, file string, result *ParseResult) {
	switch node.Type() {
	case "function_declaration", "method_declaration":
		gp.extractFunction(node, src, file, result)
	case "call_expression":
		gp.extractCallSite(node, src, file, result)
	case "composite_literal":
		gp.extractStructLiteral(node, src, file, result)
	}

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil {
			gp.walk(child, src, file, result)
		}
	}
}

func (gp *GoParser) extractFunction(node *sitter.Node, src []byte, file string, result *ParseResult) {
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
		Language:    "go",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}

	params := node.ChildByFieldName("parameters")
	if params != nil {
		paramText := params.Content(src)
		if strings.Contains(paramText, "http.ResponseWriter") && strings.Contains(paramText, "*http.Request") {
			fn.Annotations["handles_user_input"] = true
			fn.Properties["handler_type"] = "http"
		}

		// Extract individual parameter names and types
		var paramTypes []string
		var paramNames []string
		for i := 0; i < int(params.ChildCount()); i++ {
			child := params.Child(i)
			if child != nil && child.Type() == "parameter_declaration" {
				typeNode := child.ChildByFieldName("type")
				if typeNode != nil {
					paramTypes = append(paramTypes, typeNode.Content(src))
				}
				// Extract parameter names (identifiers before the type)
				for j := 0; j < int(child.ChildCount()); j++ {
					nameChild := child.Child(j)
					if nameChild != nil && nameChild.Type() == "identifier" {
						paramNames = append(paramNames, nameChild.Content(src))
					}
				}
			}
		}
		if len(paramTypes) > 0 {
			fn.Properties["param_types"] = strings.Join(paramTypes, ",")
			fn.ParamTypes = paramTypes
		}
		if len(paramNames) > 0 {
			fn.ParamNames = paramNames
		}
	}

	// Extract receiver type for method declarations
	if node.Type() == "method_declaration" {
		receiver := node.ChildByFieldName("receiver")
		if receiver != nil {
			fn.Properties["receiver"] = receiver.Content(src)
		}
	}

	// Compute cyclomatic complexity from function body
	body := node.ChildByFieldName("body")
	if body != nil {
		fn.Complexity = computeComplexity(body)
	}

	// Extract switch/case statements from function body
	gp.extractSwitchCases(node, src, fn)

	// Extract intraprocedural data flow from function body
	if body != nil {
		gp.analyzeFunctionBody(body, src, file, fn, result)
	}

	result.Functions = append(result.Functions, fn)
}

func (gp *GoParser) extractCallSite(node *sitter.Node, src []byte, file string, result *ParseResult) {
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
		Language:   "go",
		Properties: make(map[string]string),
	}
	result.CallSites = append(result.CallSites, cs)

	// Extract argument types and string values
	args := node.ChildByFieldName("arguments")
	if args != nil {
		var argTypes []string
		var stringArgs []string
		for i := 0; i < int(args.ChildCount()); i++ {
			arg := args.Child(i)
			if arg == nil {
				continue
			}
			switch arg.Type() {
			case "unary_expression":
				// &pkg.Type{} -> extract the type
				operand := arg.ChildByFieldName("operand")
				if operand != nil && operand.Type() == "composite_literal" {
					typeNode := operand.ChildByFieldName("type")
					if typeNode != nil {
						argTypes = append(argTypes, "&"+typeNode.Content(src))
					}
				}
			case "interpreted_string_literal", "raw_string_literal":
				stringArgs = append(stringArgs, strings.Trim(arg.Content(src), "\"`"))
			}
		}
		if len(argTypes) > 0 {
			cs.Properties["arg_types"] = strings.Join(argTypes, ",")
		}
		if len(stringArgs) > 0 {
			cs.Properties["string_args"] = strings.Join(stringArgs, ",")
		}
	}

	// Detect HTTP handler registrations
	if callText == "http.HandleFunc" || callText == "mux.HandleFunc" || callText == "router.HandleFunc" {
		if args != nil {
			handler := &graph.Node{
				ID:         NodeID(graph.NodeHTTPEndpoint, callText, file, line, col),
				Kind:       graph.NodeHTTPEndpoint,
				Name:       callText,
				File:       file,
				Line:       line,
				Column:     col,
				Language:   "go",
				Properties: make(map[string]string),
			}
			if args.ChildCount() > 1 {
				firstArg := args.Child(1)
				if firstArg != nil {
					route := strings.Trim(firstArg.Content(src), "\"")
					handler.Properties["route"] = route
					handler.Route = route
				}
			}
			result.HTTPHandlers = append(result.HTTPHandlers, handler)
		}
	}

	// Detect DB operations
	if isDBWrite(callText) {
		dbOp := &graph.Node{
			ID:         NodeID(graph.NodeDBOperation, callText, file, line, col),
			Kind:       graph.NodeDBOperation,
			Name:       callText,
			File:       file,
			Line:       line,
			Column:     col,
			Language:   "go",
			Properties: map[string]string{"operation": "write"},
			Operation:  "write",
		}
		gp.extractTableName(node, src, dbOp)
		result.DBOperations = append(result.DBOperations, dbOp)
	} else if isDBRead(callText) {
		dbOp := &graph.Node{
			ID:         NodeID(graph.NodeDBOperation, callText, file, line, col),
			Kind:       graph.NodeDBOperation,
			Name:       callText,
			File:       file,
			Line:       line,
			Column:     col,
			Language:   "go",
			Properties: map[string]string{"operation": "read"},
			Operation:  "read",
		}
		gp.extractTableName(node, src, dbOp)
		result.DBOperations = append(result.DBOperations, dbOp)
	}
}

// extractTableName attempts to identify the database table from SQL keywords
// in string literal arguments. This is a best-effort heuristic: it won't catch
// dynamically constructed queries, ORM model references, or multi-statement SQL.
func (gp *GoParser) extractTableName(node *sitter.Node, src []byte, dbOp *graph.Node) {
	args := node.ChildByFieldName("arguments")
	if args == nil {
		return
	}

	// First pass: look for SQL keywords in string literals
	for j := 0; j < int(args.ChildCount()); j++ {
		arg := args.Child(j)
		if arg == nil {
			continue
		}
		// Walk into call_expression args (e.g., fmt.Sprintf("INSERT INTO %s..."))
		gp.findTableInNode(arg, src, dbOp)
	}

	// Second pass: if no SQL table found, use first plain string as table name
	if dbOp.Properties["table"] == "" {
		for j := 0; j < int(args.ChildCount()); j++ {
			arg := args.Child(j)
			if arg != nil && arg.Type() == "interpreted_string_literal" {
				tableName := strings.Trim(arg.Content(src), "\"")
				if tableName != "" && !strings.Contains(tableName, " ") && tableName != "%s" {
					dbOp.Properties["table"] = tableName
					dbOp.Table = tableName
					break
				}
			}
		}
	}
}

func (gp *GoParser) findTableInNode(node *sitter.Node, src []byte, dbOp *graph.Node) {
	if node == nil {
		return
	}
	if node.Type() == "interpreted_string_literal" {
		argText := strings.Trim(node.Content(src), "\"")
		upper := strings.ToUpper(argText)
		for _, keyword := range []string{"FROM ", "INTO ", "UPDATE ", "DELETE FROM "} {
			if idx := strings.Index(upper, keyword); idx >= 0 {
				rest := argText[idx+len(keyword):]
				fields := strings.Fields(rest)
				if len(fields) > 0 {
					tableName := strings.Trim(fields[0], "\"'`()")
					if tableName != "" && tableName != "%s" {
						table := strings.ToLower(tableName)
						dbOp.Properties["table"] = table
						dbOp.Table = table
						return
					}
				}
			}
		}
	}
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil {
			gp.findTableInNode(child, src, dbOp)
			if dbOp.Properties["table"] != "" {
				return
			}
		}
	}
}

func isDBWrite(call string) bool {
	// Must have a receiver (contain a dot)
	if !strings.Contains(call, ".") {
		return false
	}
	// Check for known DB patterns
	dbReceivers := []string{"db.", "tx.", "conn.", "stmt.", "gorm.", "sqlx."}
	hasDBReceiver := false
	lower := strings.ToLower(call)
	for _, r := range dbReceivers {
		if strings.Contains(lower, r) {
			hasDBReceiver = true
			break
		}
	}
	if !hasDBReceiver {
		return false
	}
	writes := []string{".Exec", ".Create", ".Save", ".Insert", ".Update", ".Delete", ".Set"}
	for _, w := range writes {
		if strings.HasSuffix(call, w) {
			return true
		}
	}
	return false
}

func isDBRead(call string) bool {
	if !strings.Contains(call, ".") {
		return false
	}
	dbReceivers := []string{"db.", "tx.", "conn.", "stmt.", "gorm.", "sqlx."}
	hasDBReceiver := false
	lower := strings.ToLower(call)
	for _, r := range dbReceivers {
		if strings.Contains(lower, r) {
			hasDBReceiver = true
			break
		}
	}
	if !hasDBReceiver {
		return false
	}
	reads := []string{".Query", ".QueryRow", ".Find", ".First", ".Where", ".Get", ".Select"}
	for _, r := range reads {
		if strings.HasSuffix(call, r) {
			return true
		}
	}
	return false
}

func (gp *GoParser) collectStrings(node *sitter.Node, src []byte, out *[]string) {
	if node == nil {
		return
	}
	if node.Type() == "interpreted_string_literal" || node.Type() == "raw_string_literal" {
		*out = append(*out, strings.Trim(node.Content(src), "\"`"))
	}
	for i := 0; i < int(node.ChildCount()); i++ {
		gp.collectStrings(node.Child(i), src, out)
	}
}

func (gp *GoParser) extractStructLiteral(node *sitter.Node, src []byte, file string, result *ParseResult) {
	typeNode := node.ChildByFieldName("type")
	if typeNode == nil {
		return
	}
	typeName := typeNode.Content(src)

	body := node.ChildByFieldName("body")
	if body == nil {
		return
	}

	var fieldNames []string
	for i := 0; i < int(body.ChildCount()); i++ {
		elem := body.Child(i)
		if elem == nil || elem.Type() != "keyed_element" {
			continue
		}
		// In tree-sitter Go grammar, keyed_element contains literal_element nodes
		// The first literal_element typically contains the field name
		for j := 0; j < int(elem.ChildCount()); j++ {
			child := elem.Child(j)
			if child == nil {
				continue
			}
			// Try field_identifier first (older grammar versions)
			if child.Type() == "field_identifier" {
				fieldNames = append(fieldNames, child.Content(src))
				break
			}
			// Try literal_element > identifier (current grammar)
			if child.Type() == "literal_element" && child.ChildCount() > 0 {
				grandchild := child.Child(0)
				if grandchild != nil && grandchild.Type() == "identifier" {
					fieldNames = append(fieldNames, grandchild.Content(src))
					break
				}
			}
		}
	}

	var stringValues []string
	gp.collectStrings(body, src, &stringValues)

	properties := map[string]string{
		"type":   typeName,
		"fields": strings.Join(fieldNames, ","),
	}
	if len(stringValues) > 0 {
		properties["string_values"] = strings.Join(stringValues, ",")
	}

	line := int(node.StartPoint().Row) + 1
	col := int(node.StartPoint().Column)
	sl := &graph.Node{
		ID:         NodeID(graph.NodeStructLiteral, typeName, file, line, col),
		Kind:       graph.NodeStructLiteral,
		Name:       typeName,
		File:       file,
		Line:       line,
		Column:     col,
		EndLine:    int(node.EndPoint().Row) + 1,
		Language:   "go",
		Properties: properties,
		StructType: typeName,
		FieldNames: fieldNames,
	}
	result.StructLiterals = append(result.StructLiterals, sl)
}

// extractSwitchCases scans a function body for switch statements and stores
// the switch expression and case values as properties on the function node.
func (gp *GoParser) extractSwitchCases(fnNode *sitter.Node, src []byte, fn *graph.Node) {
	body := fnNode.ChildByFieldName("body")
	if body == nil {
		return
	}
	var switchExprs []string
	var caseVals []string
	gp.findSwitchStatements(body, src, &switchExprs, &caseVals)
	if len(switchExprs) > 0 {
		fn.Properties["switch_expr"] = strings.Join(switchExprs, ";")
	}
	if len(caseVals) > 0 {
		fn.Properties["case_values"] = strings.Join(caseVals, ",")
	}
}

func (gp *GoParser) findSwitchStatements(node *sitter.Node, src []byte, exprs *[]string, vals *[]string) {
	if node.Type() == "expression_switch_statement" || node.Type() == "type_switch_statement" {
		// Extract the switch value/expression
		valueNode := node.ChildByFieldName("value")
		if valueNode != nil {
			*exprs = append(*exprs, valueNode.Content(src))
		}
		// Extract case values from expression_case children
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child == nil {
				continue
			}
			if child.Type() == "expression_case" {
				valueList := child.ChildByFieldName("value")
				if valueList != nil {
					*vals = append(*vals, valueList.Content(src))
				}
			} else if child.Type() == "type_case" {
				typeNode := child.ChildByFieldName("type")
				if typeNode != nil {
					*vals = append(*vals, typeNode.Content(src))
				}
			}
		}
		return
	}
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil {
			gp.findSwitchStatements(child, src, exprs, vals)
		}
	}
}

// ---------------------------------------------------------------------------
// Intraprocedural data flow analysis
// ---------------------------------------------------------------------------

// analyzeFunctionBody creates variable/parameter nodes and data flow edges
// for a single function body. One SymbolTable + FlowBuilder per function.
func (gp *GoParser) analyzeFunctionBody(body *sitter.Node, src []byte, file string, fn *graph.Node, result *ParseResult) {
	st := dataflow.NewSymbolTable()
	fb := dataflow.NewFlowBuilder()

	// Create parameter nodes and register them in the symbol table
	for _, pName := range fn.ParamNames {
		if pName == "_" {
			continue
		}
		paramID := NodeID(graph.NodeParameter, pName, file, fn.Line, 0)
		paramNode := &graph.Node{
			ID:       paramID,
			Kind:     graph.NodeParameter,
			Name:     pName,
			File:     file,
			Line:     fn.Line,
			Language: "go",
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
		gp.analyzeStatement(child, src, file, fn, st, fb)
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

	// Build CFG from function body
	posMap := buildPositionMap(nodes)
	blocks, cfEdges := gp.buildCFG(body, src, file, fn, posMap)
	result.BasicBlocks = append(result.BasicBlocks, blocks...)
	result.Edges = append(result.Edges, cfEdges...)
}

// analyzeStatement dispatches to handlers based on tree-sitter node type.
func (gp *GoParser) analyzeStatement(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	switch node.Type() {
	case "short_var_declaration":
		gp.analyzeShortVarDecl(node, src, file, fn, st, fb)
	case "var_declaration":
		gp.analyzeVarDecl(node, src, file, fn, st, fb)
	case "assignment_statement":
		gp.analyzeAssignment(node, src, file, fn, st, fb)
	case "expression_statement":
		// Check for standalone call expressions (e.g., json.Unmarshal(body, &review))
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child != nil && child.Type() == "call_expression" {
				gp.analyzeCallArgs(child, src, file, st, fb)
			}
		}
	case "return_statement":
		gp.analyzeReturn(node, src, file, fn, st, fb)
	default:
		// Recurse into block-level structures (if, for, etc.)
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child == nil {
				continue
			}
			gp.analyzeStatement(child, src, file, fn, st, fb)
			if fb.VariableCount() >= dataflow.MaxVariablesPerFunction {
				return
			}
		}
	}
}

// analyzeShortVarDecl handles `x := expr` and `a, b := expr1, expr2`.
func (gp *GoParser) analyzeShortVarDecl(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	left := node.ChildByFieldName("left")
	right := node.ChildByFieldName("right")
	if left == nil {
		return
	}

	// Collect LHS names
	var lhsNames []*sitter.Node
	if left.Type() == "expression_list" {
		for i := 0; i < int(left.ChildCount()); i++ {
			child := left.Child(i)
			if child != nil && child.Type() == "identifier" {
				lhsNames = append(lhsNames, child)
			}
		}
	} else if left.Type() == "identifier" {
		lhsNames = append(lhsNames, left)
	}

	// Collect RHS expressions
	var rhsExprs []*sitter.Node
	if right != nil {
		if right.Type() == "expression_list" {
			for i := 0; i < int(right.ChildCount()); i++ {
				child := right.Child(i)
				if child != nil && child.Type() != "," {
					rhsExprs = append(rhsExprs, child)
				}
			}
		} else {
			rhsExprs = append(rhsExprs, right)
		}
	}

	for idx, nameNode := range lhsNames {
		name := nameNode.Content(src)
		if name == "_" {
			continue
		}

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
			Language: "go",
		}
		fb.AddVariable(varNode)
		st.Define(name, varID)

		// Resolve RHS source and create assigns edge
		if idx < len(rhsExprs) {
			rhsNode := rhsExprs[idx]
			sourceID := gp.resolveRHSSource(rhsNode, src, file, fn, st, fb)
			if sourceID != "" {
				fb.AddAssign(sourceID, varID)
			}
			// Emit reads edges for compound expressions
			gp.emitReads(rhsNode, src, file, varID, st, fb)
		}
	}
}

// analyzeVarDecl handles `var x Type = expr`.
func (gp *GoParser) analyzeVarDecl(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	// var_declaration contains var_spec children
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child == nil || child.Type() != "var_spec" {
			continue
		}
		gp.analyzeVarSpec(child, src, file, fn, st, fb)
	}
}

// analyzeVarSpec handles a single var_spec inside a var_declaration.
func (gp *GoParser) analyzeVarSpec(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		// Try iterating children for identifiers
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child != nil && child.Type() == "identifier" {
				nameNode = child
				break
			}
		}
	}
	if nameNode == nil {
		return
	}

	name := nameNode.Content(src)
	if name == "_" {
		return
	}

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
		Language: "go",
	}
	fb.AddVariable(varNode)
	st.Define(name, varID)

	// Check for value expression
	valueNode := node.ChildByFieldName("value")
	if valueNode != nil {
		sourceID := gp.resolveRHSSource(valueNode, src, file, fn, st, fb)
		if sourceID != "" {
			fb.AddAssign(sourceID, varID)
		}
	}
}

// analyzeAssignment handles `x = expr`.
func (gp *GoParser) analyzeAssignment(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	left := node.ChildByFieldName("left")
	right := node.ChildByFieldName("right")
	if left == nil || right == nil {
		return
	}

	// Get the target variable
	var targetName string
	if left.Type() == "identifier" {
		targetName = left.Content(src)
	} else if left.Type() == "expression_list" {
		for i := 0; i < int(left.ChildCount()); i++ {
			child := left.Child(i)
			if child != nil && child.Type() == "identifier" {
				targetName = child.Content(src)
				break
			}
		}
	}

	if targetName == "" || targetName == "_" {
		return
	}

	targetID, ok := st.Resolve(targetName)
	if !ok {
		return
	}

	sourceID := gp.resolveRHSSource(right, src, file, fn, st, fb)
	if sourceID != "" {
		fb.AddAssign(sourceID, targetID)
	}
	gp.emitReads(right, src, file, targetID, st, fb)
}

// analyzeReturn handles `return expr`.
func (gp *GoParser) analyzeReturn(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
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
		case "selector_expression":
			// e.g., return r.URL.Path
			rootID := gp.resolveSelectorChain(child, src, file, fn, st, fb)
			if rootID != "" {
				fb.AddReturn(rootID, fn.ID)
			}
		default:
			// Walk into compound expressions to find identifiers
			gp.walkReturnExpr(child, src, file, fn, st, fb)
		}
	}
}

// walkReturnExpr recursively finds identifiers in return expressions.
func (gp *GoParser) walkReturnExpr(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
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
	if node.Type() == "selector_expression" {
		rootID := gp.resolveSelectorChain(node, src, file, fn, st, fb)
		if rootID != "" {
			fb.AddReturn(rootID, fn.ID)
		}
		return
	}
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil {
			gp.walkReturnExpr(child, src, file, fn, st, fb)
		}
	}
}

// resolveRHSSource resolves the primary data source from an RHS expression.
// Returns a node ID that the LHS variable gets its value from.
func (gp *GoParser) resolveRHSSource(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) string {
	if node == nil {
		return ""
	}

	switch node.Type() {
	case "call_expression":
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
		gp.analyzeCallArgs(node, src, file, st, fb)
		return callID

	case "identifier":
		name := node.Content(src)
		if varID, ok := st.Resolve(name); ok {
			return varID
		}
		return ""

	case "selector_expression":
		return gp.resolveSelectorChain(node, src, file, fn, st, fb)

	case "index_expression":
		// e.g., review["name"]: resolve the object (whole collection taint)
		obj := node.ChildByFieldName("operand")
		if obj == nil {
			// Try first child
			if node.ChildCount() > 0 {
				obj = node.Child(0)
			}
		}
		if obj != nil && obj.Type() == "identifier" {
			name := obj.Content(src)
			if varID, ok := st.Resolve(name); ok {
				return varID
			}
		}
		return ""

	case "expression_list":
		// Use first element
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child != nil && child.Type() != "," {
				return gp.resolveRHSSource(child, src, file, fn, st, fb)
			}
		}
		return ""

	case "binary_expression":
		// For binary expressions like "str" + name.(string), resolve first meaningful operand
		left := node.ChildByFieldName("left")
		right := node.ChildByFieldName("right")
		// Try left first, then right
		if leftID := gp.resolveRHSSource(left, src, file, fn, st, fb); leftID != "" {
			return leftID
		}
		return gp.resolveRHSSource(right, src, file, fn, st, fb)

	case "type_assertion_expression":
		// e.g., name.(string): resolve the operand
		operand := node.ChildByFieldName("operand")
		if operand != nil {
			return gp.resolveRHSSource(operand, src, file, fn, st, fb)
		}
		return ""

	case "unary_expression":
		// e.g., &http.Request{}: resolve the operand
		operand := node.ChildByFieldName("operand")
		if operand != nil {
			return gp.resolveRHSSource(operand, src, file, fn, st, fb)
		}
		return ""

	default:
		return ""
	}
}

// resolveSelectorChain handles selector_expression chains like r.URL.Path.
// Creates a single synthetic field variable with the full text (e.g., "r.URL.Path")
// and emits one field_access edge from the root identifier to the synthetic variable.
// The synthetic variable is NOT registered in the symbol table (per spec).
func (gp *GoParser) resolveSelectorChain(node *sitter.Node, src []byte, file string, fn *graph.Node, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) string {
	if node == nil || node.Type() != "selector_expression" {
		return ""
	}

	// Get the full text of the selector expression (e.g., "r.URL.Path")
	fullText := node.Content(src)

	// Walk to the leftmost identifier to find the root
	var rootID string
	current := node
	for {
		operand := current.ChildByFieldName("operand")
		if operand == nil {
			break
		}

		if operand.Type() == "identifier" {
			// Found the root identifier
			rootName := operand.Content(src)
			if id, ok := st.Resolve(rootName); ok {
				rootID = id
			}
			break
		} else if operand.Type() == "selector_expression" {
			// Continue walking left
			current = operand
		} else {
			// Unexpected operand type, stop
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
		Language: "go",
	}
	fb.AddVariable(fieldVarNode)
	// Emit ONE field_access edge from root to synthetic variable
	fb.AddFieldAccess(rootID, fieldVarID)
	// Do NOT register in symbol table (per spec)

	return fieldVarID
}

// analyzeCallArgs processes arguments to a call expression, creating
// passes_to edges and mutates edges for &x patterns.
func (gp *GoParser) analyzeCallArgs(callNode *sitter.Node, src []byte, file string, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
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

		case "unary_expression":
			// Check for & operator (address-of)
			operand := arg.ChildByFieldName("operand")
			operator := arg.ChildByFieldName("operator")
			isAddressOf := false
			if operator != nil && operator.Content(src) == "&" {
				isAddressOf = true
			} else {
				// Check for "&" as a direct child
				for j := 0; j < int(arg.ChildCount()); j++ {
					child := arg.Child(j)
					if child != nil && child.Type() == "&" {
						isAddressOf = true
						break
					}
				}
			}

			if isAddressOf && operand != nil {
				// Handle &identifier
				if operand.Type() == "identifier" {
					name := operand.Content(src)
					if varID, ok := st.Resolve(name); ok {
						fb.AddPassesTo(varID, callID)
						fb.AddMutates(callID, varID)
					}
				} else if operand.Type() == "selector_expression" {
					// Handle &x.Field or &x[0]: resolve root identifier and apply both edges to it
					rootID := gp.findRootIdentifier(operand, src, st)
					if rootID != "" {
						fb.AddPassesTo(rootID, callID)
						fb.AddMutates(callID, rootID)
					}
				} else if operand.Type() == "index_expression" {
					// Handle &x[0]: resolve the array/slice variable
					indexOperand := operand.ChildByFieldName("operand")
					if indexOperand != nil && indexOperand.Type() == "identifier" {
						name := indexOperand.Content(src)
						if varID, ok := st.Resolve(name); ok {
							fb.AddPassesTo(varID, callID)
							fb.AddMutates(callID, varID)
						}
					}
				}
			}

		case "selector_expression":
			// e.g., r.Body: passes_to from root identifier
			gp.passSelectorToCall(arg, src, callID, st, fb)
		}
	}
}

// passSelectorToCall finds the root identifier of a selector expression
// and emits a passes_to edge.
func (gp *GoParser) passSelectorToCall(node *sitter.Node, src []byte, callID string, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	if node == nil {
		return
	}
	switch node.Type() {
	case "identifier":
		name := node.Content(src)
		if varID, ok := st.Resolve(name); ok {
			fb.AddPassesTo(varID, callID)
		}
	case "selector_expression":
		operand := node.ChildByFieldName("operand")
		if operand != nil {
			gp.passSelectorToCall(operand, src, callID, st, fb)
		}
	}
}

// findRootIdentifier walks a selector expression to find the leftmost identifier
// and resolves it in the symbol table. Used for &x.Field patterns.
func (gp *GoParser) findRootIdentifier(node *sitter.Node, src []byte, st *dataflow.SymbolTable) string {
	if node == nil {
		return ""
	}

	switch node.Type() {
	case "identifier":
		name := node.Content(src)
		if varID, ok := st.Resolve(name); ok {
			return varID
		}
		return ""

	case "selector_expression":
		operand := node.ChildByFieldName("operand")
		if operand != nil {
			return gp.findRootIdentifier(operand, src, st)
		}
		return ""

	default:
		return ""
	}
}

// emitReads walks a compound expression (binary_expression, etc.) to find all
// identifier references and emits reads edges.
func (gp *GoParser) emitReads(node *sitter.Node, src []byte, file string, targetID string, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
	if node == nil {
		return
	}

	if node.Type() != "binary_expression" {
		return
	}

	// Walk all children to find identifier reads
	gp.walkForReads(node, src, targetID, st, fb)
}

// walkForReads recursively finds identifiers in an expression tree and emits reads edges.
func (gp *GoParser) walkForReads(node *sitter.Node, src []byte, targetID string, st *dataflow.SymbolTable, fb *dataflow.FlowBuilder) {
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

	// Recurse into type_assertion_expression operands
	if node.Type() == "type_assertion_expression" {
		operand := node.ChildByFieldName("operand")
		if operand != nil {
			gp.walkForReads(operand, src, targetID, st, fb)
		}
		return
	}

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil {
			gp.walkForReads(child, src, targetID, st, fb)
		}
	}
}

// ---------------------------------------------------------------------------
// Control Flow Graph Construction
// ---------------------------------------------------------------------------

// posKey represents a source code position for indexing B1 nodes.
type posKey struct {
	Line int
	Col  int
}

// buildPositionMap indexes B1 nodes (variables, parameters) by their position.
// Used by CFG construction to collect node IDs at AST positions.
func buildPositionMap(nodes []*graph.Node) map[posKey][]string {
	m := make(map[posKey][]string)
	for _, n := range nodes {
		key := posKey{Line: n.Line, Col: n.Column}
		m[key] = append(m[key], n.ID)
	}
	return m
}

// collectNodeIDs walks an AST node recursively and collects all B1 node IDs
// at positions within it, using the position map.
func collectNodeIDs(node *sitter.Node, posMap map[posKey][]string) []string {
	if node == nil {
		return nil
	}

	var ids []string
	line := int(node.StartPoint().Row) + 1
	col := int(node.StartPoint().Column)
	key := posKey{Line: line, Col: col}
	if nodeIDs, ok := posMap[key]; ok {
		ids = append(ids, nodeIDs...)
	}

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil {
			ids = append(ids, collectNodeIDs(child, posMap)...)
		}
	}

	return ids
}

// buildCFG constructs a control flow graph for a function body.
// Creates entry/exit blocks, connects control flow edges, and populates
// block members with B1 node IDs.
func (gp *GoParser) buildCFG(body *sitter.Node, src []byte, file string, fn *graph.Node, posMap map[posKey][]string) ([]*graph.Node, []*graph.Edge) {
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
	lastBlock := gp.buildCFGStatements(body, src, file, fn, entry, exit, posMap, cb)

	// Connect last block to exit (if it didn't terminate early)
	if lastBlock != nil {
		cb.AddEdge(lastBlock.ID, exit.ID, "exit")
	}

	return cb.Result()
}

// buildCFGStatements walks container children and builds CFG.
// Returns the last open block (nil if all paths terminated).
func (gp *GoParser) buildCFGStatements(container *sitter.Node, src []byte, file string, fn *graph.Node, current, exit *graph.Node, posMap map[posKey][]string, cb *dataflow.CFGBuilder) *graph.Node {
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

		switch child.Type() {
		case "if_statement":
			current = gp.buildCFGIf(child, src, file, fn, current, exit, posMap, cb)
			if current == nil {
				return nil
			}

		case "for_statement":
			current = gp.buildCFGFor(child, src, file, fn, current, exit, posMap, cb)
			if current == nil {
				return nil
			}

		case "expression_switch_statement":
			current = gp.buildCFGSwitch(child, src, file, fn, current, exit, posMap, cb)
			if current == nil {
				return nil
			}

		case "return_statement":
			// Collect node IDs from return statement and add to current block
			nodeIDs := collectNodeIDs(child, posMap)
			for _, id := range nodeIDs {
				cb.AddMember(current, id)
			}
			// Connect to exit and terminate this path
			cb.AddEdge(current.ID, exit.ID, "exit")
			return nil

		case "expression_statement":
			// Check for panic() calls
			if gp.containsPanic(child, src) {
				// Collect node IDs and add to current block
				nodeIDs := collectNodeIDs(child, posMap)
				for _, id := range nodeIDs {
					cb.AddMember(current, id)
				}
				// panic terminates control flow
				cb.AddEdge(current.ID, exit.ID, "exit")
				return nil
			}
			// Regular statement: add node IDs to current block
			nodeIDs := collectNodeIDs(child, posMap)
			for _, id := range nodeIDs {
				cb.AddMember(current, id)
			}

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

// containsPanic checks if an expression_statement contains a panic() call.
func (gp *GoParser) containsPanic(node *sitter.Node, src []byte) bool {
	if node == nil {
		return false
	}

	if node.Type() == "call_expression" {
		fnNode := node.ChildByFieldName("function")
		if fnNode != nil && fnNode.Content(src) == "panic" {
			return true
		}
	}

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil && gp.containsPanic(child, src) {
			return true
		}
	}

	return false
}

// buildCFGIf handles if_statement nodes.
// Creates condition block, then-block, else-block (if present), and merge block.
// Returns merge block or nil if all branches terminate.
func (gp *GoParser) buildCFGIf(node *sitter.Node, src []byte, file string, fn *graph.Node, current, exit *graph.Node, posMap map[posKey][]string, cb *dataflow.CFGBuilder) *graph.Node {
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

	// Process then-branch
	consequence := node.ChildByFieldName("consequence")
	thenEnd := gp.buildCFGStatements(consequence, src, file, fn, thenBlock, exit, posMap, cb)

	// Check for else-branch
	alternative := node.ChildByFieldName("alternative")
	var elseEnd *graph.Node

	if alternative != nil {
		// Check if alternative is an else-if (another if_statement)
		if alternative.Type() == "if_statement" {
			// else if: recurse into buildCFGIf from current block (false_branch)
			// Create a block for the else-if condition
			elseIfBlock := cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), int(alternative.StartPoint().Row)+1)
			cb.AddEdge(current.ID, elseIfBlock.ID, "false_branch")
			elseEnd = gp.buildCFGIf(alternative, src, file, fn, elseIfBlock, exit, posMap, cb)
		} else {
			// Regular else block
			elseBlock := cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), int(alternative.StartPoint().Row)+1)
			cb.AddEdge(current.ID, elseBlock.ID, "false_branch")
			elseEnd = gp.buildCFGStatements(alternative, src, file, fn, elseBlock, exit, posMap, cb)
		}
	} else {
		// No else: false_branch goes to merge block
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
		cb.AddEdge(elseEnd.ID, mergeBlock.ID, "fallthrough")
	} else if alternative == nil {
		// No else branch: connect current to merge via false_branch
		cb.AddEdge(current.ID, mergeBlock.ID, "false_branch")
	}

	return mergeBlock
}

// buildCFGFor handles for_statement nodes.
// Creates header-block, body-block, loop_back, and loop_exit edges.
func (gp *GoParser) buildCFGFor(node *sitter.Node, src []byte, file string, fn *graph.Node, current, exit *graph.Node, posMap map[posKey][]string, cb *dataflow.CFGBuilder) *graph.Node {
	line := int(node.StartPoint().Row) + 1

	// Create header block (loop condition evaluation)
	headerBlock := cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), line)
	cb.AddEdge(current.ID, headerBlock.ID, "fallthrough")

	// Create body block
	bodyBlock := cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), line)
	cb.AddEdge(headerBlock.ID, bodyBlock.ID, "true_branch")

	// Process loop body
	body := node.ChildByFieldName("body")
	if body == nil {
		// Try to find block as last child
		for i := int(node.ChildCount()) - 1; i >= 0; i-- {
			child := node.Child(i)
			if child != nil && child.Type() == "block" {
				body = child
				break
			}
		}
	}

	bodyEnd := gp.buildCFGStatements(body, src, file, fn, bodyBlock, exit, posMap, cb)

	// Loop back from body end to header
	if bodyEnd != nil {
		cb.AddEdge(bodyEnd.ID, headerBlock.ID, "loop_back")
	}

	// Create after-loop block (loop exit)
	afterBlock := cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), line)
	cb.AddEdge(headerBlock.ID, afterBlock.ID, "loop_exit")

	return afterBlock
}

// buildCFGSwitch handles expression_switch_statement nodes.
// Creates one case block per case, connects with true_branch/false_branch, merges.
func (gp *GoParser) buildCFGSwitch(node *sitter.Node, src []byte, file string, fn *graph.Node, current, exit *graph.Node, posMap map[posKey][]string, cb *dataflow.CFGBuilder) *graph.Node {
	line := int(node.StartPoint().Row) + 1

	// Switch value evaluation happens in current block
	value := node.ChildByFieldName("value")
	if value != nil {
		nodeIDs := collectNodeIDs(value, posMap)
		for _, id := range nodeIDs {
			cb.AddMember(current, id)
		}
	}

	// Create merge block
	mergeBlock := cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), line)

	// Process each case
	hasDefault := false
	var caseEnds []*graph.Node

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}

		var caseBlock *graph.Node
		var label string

		if child.Type() == "expression_case" {
			caseBlock = cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), int(child.StartPoint().Row)+1)
			label = "true_branch"
		} else if child.Type() == "default_case" {
			caseBlock = cb.NewBlock(fmt.Sprintf("bb%d", cb.BlockCount()-2), int(child.StartPoint().Row)+1)
			label = "false_branch"
			hasDefault = true
		} else {
			continue
		}

		cb.AddEdge(current.ID, caseBlock.ID, label)

		// Process case body
		caseEnd := gp.buildCFGStatements(child, src, file, fn, caseBlock, exit, posMap, cb)
		caseEnds = append(caseEnds, caseEnd)
	}

	// Connect case ends to merge block
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
