package graph

// EdgeKind classifies the relationship between two graph nodes.
type EdgeKind string

const (
	EdgeCalls         EdgeKind = "CALLS"
	EdgeDataFlow      EdgeKind = "DATA_FLOW"
	EdgeControlFlow   EdgeKind = "CONTROL_FLOW"
	EdgeStorageLink   EdgeKind = "STORAGE_LINK"
	EdgeAnnotatedWith EdgeKind = "ANNOTATED_WITH"
)

// EdgeConfidence expresses how reliable a call resolution is.
type EdgeConfidence string

const (
	ConfidenceCertain   EdgeConfidence = "CERTAIN"
	ConfidenceInferred  EdgeConfidence = "INFERRED"
	ConfidenceUncertain EdgeConfidence = "UNCERTAIN"
)

// Edge represents a directed relationship between two nodes in the code property graph.
type Edge struct {
	From       string         `json:"from"`
	To         string         `json:"to"`
	Kind       EdgeKind       `json:"kind"`
	Label      string         `json:"label,omitempty"`
	Confidence EdgeConfidence `json:"confidence,omitempty"`
}
