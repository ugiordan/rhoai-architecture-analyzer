package parser

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

// MaxFileSize is the maximum file size (10 MB) that parsers will attempt to parse.
// Files larger than this are skipped to avoid excessive memory use in tree-sitter.
const MaxFileSize = 10 * 1024 * 1024

// GoParser extracts code property graph nodes from Go source files using tree-sitter.
// The idSeq counter is safe for concurrent use via atomic operations. The underlying
// tree-sitter parser is NOT safe for concurrent use even with external locks: each
// goroutine MUST use its own GoParser instance (call NewGoParser per goroutine).
// When running multiple GoParser instances in parallel, pass a shared counter via
// NewGoParserWithSeq to avoid node ID collisions.
type GoParser struct {
	parser *sitter.Parser
	idSeq  *atomic.Int64
}

// NewGoParser creates a parser for Go source files backed by tree-sitter.
func NewGoParser() *GoParser {
	p := sitter.NewParser()
	p.SetLanguage(golang.GetLanguage())
	return &GoParser{parser: p, idSeq: &atomic.Int64{}}
}

// NewGoParserWithSeq creates a parser that shares an ID counter with other instances.
// Use this when running multiple parsers in parallel to avoid node ID collisions.
func NewGoParserWithSeq(seq *atomic.Int64) *GoParser {
	p := sitter.NewParser()
	p.SetLanguage(golang.GetLanguage())
	return &GoParser{parser: p, idSeq: seq}
}

func (gp *GoParser) Language() string     { return "go" }
func (gp *GoParser) Extensions() []string { return []string{".go"} }
func (gp *GoParser) CloneWithSeq(seq *atomic.Int64) Parser {
	return NewGoParserWithSeq(seq)
}

func (gp *GoParser) nextID(prefix string) string {
	id := gp.idSeq.Add(1)
	return fmt.Sprintf("%s_%d", prefix, id)
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
	fn := &graph.Node{
		ID:          gp.nextID("fn"),
		Kind:        graph.NodeFunction,
		Name:        name,
		File:        file,
		Line:        int(node.StartPoint().Row) + 1,
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

		// Extract individual parameter types
		var paramTypes []string
		for i := 0; i < int(params.ChildCount()); i++ {
			child := params.Child(i)
			if child != nil && child.Type() == "parameter_declaration" {
				typeNode := child.ChildByFieldName("type")
				if typeNode != nil {
					paramTypes = append(paramTypes, typeNode.Content(src))
				}
			}
		}
		if len(paramTypes) > 0 {
			fn.Properties["param_types"] = strings.Join(paramTypes, ",")
			fn.ParamTypes = paramTypes
		}
	}

	// Extract receiver type for method declarations
	if node.Type() == "method_declaration" {
		receiver := node.ChildByFieldName("receiver")
		if receiver != nil {
			fn.Properties["receiver"] = receiver.Content(src)
		}
	}

	// Extract switch/case statements from function body
	gp.extractSwitchCases(node, src, fn)

	result.Functions = append(result.Functions, fn)
}

func (gp *GoParser) extractCallSite(node *sitter.Node, src []byte, file string, result *ParseResult) {
	fnNode := node.ChildByFieldName("function")
	if fnNode == nil {
		return
	}
	callText := fnNode.Content(src)
	line := int(node.StartPoint().Row) + 1

	cs := &graph.Node{
		ID:         gp.nextID("call"),
		Kind:       graph.NodeCallSite,
		Name:       callText,
		File:       file,
		Line:       line,
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
				ID:         gp.nextID("http"),
				Kind:       graph.NodeHTTPEndpoint,
				Name:       callText,
				File:       file,
				Line:       line,
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
			ID:         gp.nextID("db"),
			Kind:       graph.NodeDBOperation,
			Name:       callText,
			File:       file,
			Line:       line,
			Language:   "go",
			Properties: map[string]string{"operation": "write"},
			Operation:  "write",
		}
		gp.extractTableName(node, src, dbOp)
		result.DBOperations = append(result.DBOperations, dbOp)
	} else if isDBRead(callText) {
		dbOp := &graph.Node{
			ID:         gp.nextID("db"),
			Kind:       graph.NodeDBOperation,
			Name:       callText,
			File:       file,
			Line:       line,
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

	sl := &graph.Node{
		ID:         gp.nextID("struct"),
		Kind:       graph.NodeStructLiteral,
		Name:       typeName,
		File:       file,
		Line:       int(node.StartPoint().Row) + 1,
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
