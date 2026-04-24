package parser

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/typescript/tsx"

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
// (tree-sitter parsers are not thread-safe). When running multiple instances in
// parallel, pass a shared counter via NewTypeScriptParserWithSeq.
type TypeScriptParser struct {
	parser *sitter.Parser
	idSeq  *atomic.Int64
}

// NewTypeScriptParser creates a parser for TypeScript source files backed by tree-sitter.
// Uses the TSX grammar for all files since it's a superset of TypeScript.
func NewTypeScriptParser() *TypeScriptParser {
	p := sitter.NewParser()
	p.SetLanguage(tsx.GetLanguage())
	return &TypeScriptParser{parser: p, idSeq: &atomic.Int64{}}
}

// NewTypeScriptParserWithSeq creates a parser that shares an ID counter with other instances.
// Use this when running multiple parsers in parallel to avoid node ID collisions.
func NewTypeScriptParserWithSeq(seq *atomic.Int64) *TypeScriptParser {
	p := sitter.NewParser()
	p.SetLanguage(tsx.GetLanguage())
	return &TypeScriptParser{parser: p, idSeq: seq}
}

func (tp *TypeScriptParser) Language() string     { return "typescript" }
func (tp *TypeScriptParser) Extensions() []string { return []string{".ts", ".tsx"} }

func (tp *TypeScriptParser) nextID(prefix string) string {
	id := tp.idSeq.Add(1)
	return fmt.Sprintf("%s_%d", prefix, id)
}

// ParseFile parses a TypeScript/TSX source file and returns extracted nodes and edges.
// Declaration files (.d.ts) are skipped, returning an empty result.
func (tp *TypeScriptParser) ParseFile(path string, content []byte) (*ParseResult, error) {
	// Skip declaration files
	if strings.HasSuffix(path, ".d.ts") {
		return &ParseResult{}, nil
	}

	if len(content) > maxFileSize {
		return nil, fmt.Errorf("file too large (%d bytes, max %d)", len(content), maxFileSize)
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
	fn := &graph.Node{
		ID:          tp.nextID("fn"),
		Kind:        graph.NodeFunction,
		Name:        nameNode.Content(src),
		File:        file,
		Line:        int(node.StartPoint().Row) + 1,
		EndLine:     int(node.EndPoint().Row) + 1,
		Language:    "typescript",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	result.Functions = append(result.Functions, fn)
}

// extractMethod handles method_definition nodes inside classes.
func (tp *TypeScriptParser) extractMethod(node *sitter.Node, src []byte, file, className string, result *ParseResult) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}
	fn := &graph.Node{
		ID:          tp.nextID("fn"),
		Kind:        graph.NodeFunction,
		Name:        nameNode.Content(src),
		File:        file,
		Line:        int(node.StartPoint().Row) + 1,
		EndLine:     int(node.EndPoint().Row) + 1,
		Language:    "typescript",
		TypeName:    className,
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	result.Functions = append(result.Functions, fn)
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
			fn := &graph.Node{
				ID:          tp.nextID("fn"),
				Kind:        graph.NodeFunction,
				Name:        nameNode.Content(src),
				File:        file,
				Line:        int(child.StartPoint().Row) + 1,
				EndLine:     int(child.EndPoint().Row) + 1,
				Language:    "typescript",
				Annotations: make(map[string]bool),
				Properties:  make(map[string]string),
			}
			result.Functions = append(result.Functions, fn)
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

	cs := &graph.Node{
		ID:         tp.nextID("call"),
		Kind:       graph.NodeCallSite,
		Name:       callText,
		File:       file,
		Line:       line,
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
				tp.maybeExtractExpressHandler(node, src, file, line, callText, httpMethod, result)
			} else if op, ok := tsDBOps[method]; ok {
				// Only treat as DB operation if not already matched as HTTP handler
				dbOp := &graph.Node{
					ID:         tp.nextID("db"),
					Kind:       graph.NodeDBOperation,
					Name:       callText,
					File:       file,
					Line:       line,
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
func (tp *TypeScriptParser) maybeExtractExpressHandler(node *sitter.Node, src []byte, file string, line int, callText, httpMethod string, result *ParseResult) {
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
		ID:       tp.nextID("http"),
		Kind:     graph.NodeHTTPEndpoint,
		Name:     callText,
		File:     file,
		Line:     line,
		Language: "typescript",
		Annotations: make(map[string]bool),
		Properties:  map[string]string{
			"method": httpMethod,
		},
	}
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

	sl := &graph.Node{
		ID:         tp.nextID("struct"),
		Kind:       graph.NodeStructLiteral,
		Name:       name,
		File:       file,
		Line:       int(node.StartPoint().Row) + 1,
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

	handler := &graph.Node{
		ID:       tp.nextID("http"),
		Kind:     graph.NodeHTTPEndpoint,
		Name:     "Route",
		File:     file,
		Line:     int(node.StartPoint().Row) + 1,
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
