package graph

// NodeKind classifies the type of a code property graph node.
type NodeKind string

const (
	NodeFunction      NodeKind = "Function"
	NodeParameter     NodeKind = "Parameter"
	NodeVariable      NodeKind = "Variable"
	NodeCallSite      NodeKind = "CallSite"
	NodeLiteral       NodeKind = "Literal"
	NodeHTTPEndpoint  NodeKind = "HTTPEndpoint"
	NodeDBOperation   NodeKind = "DBOperation"
	NodeExternalCall  NodeKind = "ExternalCall"
	NodeK8sResource   NodeKind = "K8sResource"
	NodeStructLiteral NodeKind = "StructLiteral"
)

// Node represents a vertex in the code property graph (function, call site, HTTP endpoint, etc).
type Node struct {
	ID          string            `json:"id"`
	Kind        NodeKind          `json:"kind"`
	Name        string            `json:"name"`
	File        string            `json:"file"`
	Line        int               `json:"line"`
	EndLine     int               `json:"end_line,omitempty"`
	Language    string            `json:"language,omitempty"`
	TypeName    string            `json:"type_name,omitempty"`
	Decorators  []string          `json:"decorators,omitempty"`
	Annotations map[string]bool   `json:"annotations,omitempty"`
	Properties  map[string]string `json:"properties,omitempty"`
}
