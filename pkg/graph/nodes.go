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

// TrustLevel classifies HTTP handlers and entrypoints by trust level.
// Populated in Phase A2 by domain annotators.
type TrustLevel string

const (
	TrustUntrusted   TrustLevel = "untrusted"
	TrustSemiTrusted TrustLevel = "semi_trusted"
	TrustTrusted     TrustLevel = "trusted"
)

// Node represents a vertex in the code property graph (function, call site, HTTP endpoint, etc).
type Node struct {
	// Core fields (all node kinds)
	ID          string          `json:"id"`
	Kind        NodeKind        `json:"kind"`
	Name        string          `json:"name"`
	File        string          `json:"file"`
	Line        int             `json:"line"`
	EndLine     int             `json:"end_line,omitempty"`
	Language    string          `json:"language,omitempty"`
	TypeName    string          `json:"type_name,omitempty"`
	Decorators  []string        `json:"decorators,omitempty"`
	Annotations map[string]bool `json:"annotations,omitempty"`

	// Function fields
	Complexity int      `json:"complexity,omitempty"`
	ParamNames []string `json:"param_names,omitempty"`
	ParamTypes []string `json:"param_types,omitempty"`
	ReturnType string   `json:"return_type,omitempty"`
	IsTest     bool     `json:"is_test,omitempty"`
	IsUnsafe   bool     `json:"is_unsafe,omitempty"`
	IsExtern   bool     `json:"is_extern,omitempty"`

	// Call site fields
	CallTarget string `json:"call_target,omitempty"`
	IsMacro    bool   `json:"is_macro,omitempty"`

	// HTTP endpoint fields
	Route      string `json:"route,omitempty"`
	HTTPMethod string `json:"http_method,omitempty"`

	// DB operation fields
	Operation string `json:"operation,omitempty"`
	Table     string `json:"table,omitempty"`

	// Struct literal fields
	StructType string   `json:"struct_type,omitempty"`
	FieldNames []string `json:"field_names,omitempty"`

	// Entrypoint classification (populated in A2)
	TrustLevel TrustLevel `json:"trust_level,omitempty"`

	// Preserved for language-specific edge cases
	Properties map[string]string `json:"properties,omitempty"`
}
