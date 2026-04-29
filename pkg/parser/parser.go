package parser

import (
	"crypto/sha256"
	"fmt"

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

// kindPrefix maps NodeKind to the short prefix used in node IDs.
var kindPrefix = map[graph.NodeKind]string{
	graph.NodeFunction:      "fn",
	graph.NodeParameter:     "param",
	graph.NodeVariable:      "var",
	graph.NodeCallSite:      "call",
	graph.NodeLiteral:       "lit",
	graph.NodeHTTPEndpoint:  "http",
	graph.NodeDBOperation:   "db",
	graph.NodeExternalCall:  "ext",
	graph.NodeK8sResource:   "k8s",
	graph.NodeStructLiteral:   "struct",
	graph.NodeExternalFinding: "extf",
}

// NodeID produces a deterministic, stable node ID from the node's identity fields.
// The ID is a SHA-256 hash of (Kind, Name, File, Line, Column), truncated to
// 16 hex characters (64 bits), prefixed with a kind-specific short name.
// Format: fn_a3b2c1d4e5f67890
func NodeID(kind graph.NodeKind, name, file string, line, column int) string {
	prefix := kindPrefix[kind]
	if prefix == "" {
		prefix = "node"
	}
	input := fmt.Sprintf("%s:%s:%s:%d:%d", kind, name, file, line, column)
	sum := sha256.Sum256([]byte(input))
	return fmt.Sprintf("%s_%x", prefix, sum[:8])
}
