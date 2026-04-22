package parser

import "github.com/ugiordan/architecture-analyzer/pkg/graph"

// ParseResult holds the nodes and edges extracted from a single source file.
type ParseResult struct {
	Functions      []*graph.Node
	CallSites      []*graph.Node
	Edges          []*graph.Edge
	HTTPHandlers   []*graph.Node
	DBOperations   []*graph.Node
	StructLiterals []*graph.Node
}

// Parser extracts code property graph nodes and edges from source files.
type Parser interface {
	ParseFile(path string, content []byte) (*ParseResult, error)
	Language() string
	Extensions() []string
}
