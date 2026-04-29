package graph

import (
	"crypto/sha256"
	"fmt"
)

// kindPrefix maps NodeKind to the short prefix used in node IDs.
var kindPrefix = map[NodeKind]string{
	NodeFunction:        "fn",
	NodeParameter:       "param",
	NodeVariable:        "var",
	NodeCallSite:        "call",
	NodeLiteral:         "lit",
	NodeHTTPEndpoint:    "http",
	NodeDBOperation:     "db",
	NodeExternalCall:    "ext",
	NodeK8sResource:     "k8s",
	NodeStructLiteral:   "struct",
	NodeExternalFinding: "extf",
	NodeBasicBlock:      "bb",
}

// NodeID produces a deterministic, stable node ID from the node's identity fields.
// The ID is a SHA-256 hash of (Kind, Name, File, Line, Column), truncated to
// 16 hex characters (64 bits), prefixed with a kind-specific short name.
// Format: fn_a3b2c1d4e5f67890
func NodeID(kind NodeKind, name, file string, line, column int) string {
	prefix := kindPrefix[kind]
	if prefix == "" {
		prefix = "node"
	}
	input := fmt.Sprintf("%s:%s:%s:%d:%d", kind, name, file, line, column)
	sum := sha256.Sum256([]byte(input))
	return fmt.Sprintf("%s_%x", prefix, sum[:8])
}
