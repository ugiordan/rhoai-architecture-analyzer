package parser

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/rust"

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
// When running multiple instances in parallel, pass a shared counter via NewRustParserWithSeq.
type RustParser struct {
	parser *sitter.Parser
	idSeq  *atomic.Int64
}

// NewRustParser creates a parser for Rust source files backed by tree-sitter.
func NewRustParser() *RustParser {
	p := sitter.NewParser()
	p.SetLanguage(rust.GetLanguage())
	return &RustParser{parser: p, idSeq: &atomic.Int64{}}
}

// NewRustParserWithSeq creates a parser that shares an ID counter with other instances.
func NewRustParserWithSeq(seq *atomic.Int64) *RustParser {
	p := sitter.NewParser()
	p.SetLanguage(rust.GetLanguage())
	return &RustParser{parser: p, idSeq: seq}
}

func (rp *RustParser) Language() string     { return "rust" }
func (rp *RustParser) Extensions() []string { return []string{".rs"} }

func (rp *RustParser) nextID(prefix string) string {
	id := rp.idSeq.Add(1)
	return fmt.Sprintf("%s_%d", prefix, id)
}

// ParseFile parses a Rust source file and returns extracted nodes and edges.
func (rp *RustParser) ParseFile(path string, content []byte) (*ParseResult, error) {
	if len(content) > maxFileSize {
		return nil, fmt.Errorf("file too large (%d bytes, max %d)", len(content), maxFileSize)
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

	fn := &graph.Node{
		ID:          rp.nextID("fn"),
		Kind:        graph.NodeFunction,
		Name:        name,
		File:        file,
		Line:        int(node.StartPoint().Row) + 1,
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
					ID:          rp.nextID("http"),
					Kind:        graph.NodeHTTPEndpoint,
					Name:        fn.Name,
					File:        file,
					Line:        fn.Line,
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

	cs := &graph.Node{
		ID:         rp.nextID("call"),
		Kind:       graph.NodeCallSite,
		Name:       callText,
		File:       file,
		Line:       line,
		Language:   "rust",
		Properties: make(map[string]string),
	}
	result.CallSites = append(result.CallSites, cs)

	// Check for DB operations by looking at the full call chain text
	rp.maybeExtractDBFromCall(callText, node, src, file, line, result)
}

// maybeExtractDBFromCall checks if a call expression's direct function name matches a known DB pattern.
// Only matches on the innermost call to avoid duplicates from chained method calls.
func (rp *RustParser) maybeExtractDBFromCall(callText string, _ *sitter.Node, _ []byte, file string, line int, result *ParseResult) {
	// Check if the direct call name matches a DB pattern exactly.
	// callText is the function field content. For chained calls like
	// diesel::insert_into(items::table).values(...).execute(...), the outer calls have
	// callText that includes the full chain. We only want the innermost match.
	for pattern, op := range rustDBPatterns {
		if callText == pattern {
			dbOp := &graph.Node{
				ID:         rp.nextID("db"),
				Kind:       graph.NodeDBOperation,
				Name:       pattern,
				File:       file,
				Line:       line,
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

	cs := &graph.Node{
		ID:         rp.nextID("call"),
		Kind:       graph.NodeCallSite,
		Name:       displayName,
		File:       file,
		Line:       line,
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
			ID:         rp.nextID("db"),
			Kind:       graph.NodeDBOperation,
			Name:       displayName,
			File:       file,
			Line:       line,
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

	sl := &graph.Node{
		ID:         rp.nextID("struct"),
		Kind:       graph.NodeStructLiteral,
		Name:       name,
		File:       file,
		Line:       int(node.StartPoint().Row) + 1,
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
