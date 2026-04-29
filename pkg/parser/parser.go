package parser

import (
	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

// ParseResult holds the nodes and edges extracted from a single source file.
type ParseResult struct {
	Functions      []*graph.Node
	CallSites      []*graph.Node
	Edges          []*graph.Edge
	HTTPHandlers   []*graph.Node
	DBOperations   []*graph.Node
	StructLiterals []*graph.Node
	Variables      []*graph.Node
	Parameters     []*graph.Node
	BasicBlocks    []*graph.Node
}

// Parser extracts code property graph nodes and edges from source files.
type Parser interface {
	ParseFile(path string, content []byte) (*ParseResult, error)
	Language() string
	Extensions() []string
	// Clone returns a new parser instance for use in a separate goroutine.
	// Each goroutine needs its own parser (tree-sitter requirement).
	Clone() Parser
}

// NodeID delegates to graph.NodeID. The function conceptually belongs to the
// graph package (it operates on graph.NodeKind and produces graph node IDs).
// Kept here as a convenience so parser-internal callers don't need updating.
var NodeID = graph.NodeID
