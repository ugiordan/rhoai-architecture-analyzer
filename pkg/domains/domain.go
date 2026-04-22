// Package domains provides pluggable domain analysis capabilities for the code
// property graph. Each domain adds annotations and queries specific to its concern
// (security, testing, upgrade).
package domains

import (
	"github.com/ugiordan/architecture-analyzer/pkg/graph"
	"github.com/ugiordan/architecture-analyzer/pkg/query"
)

// ArchitectureData holds extracted architecture information that domain analyzers
// can cross-reference. Nil when --with-arch is not specified.
type ArchitectureData struct {
	Raw map[string]interface{}
}

// DomainAnalyzer provides domain-specific code analysis capabilities.
// Domains add annotations and edges to a shared CPG, then provide
// query rules that detect patterns in the annotated graph.
type DomainAnalyzer interface {
	// Name returns the domain identifier (e.g., "security", "testing", "upgrade").
	Name() string

	// SupportedLanguages returns which languages this analyzer handles.
	SupportedLanguages() []string

	// Annotate adds domain-specific annotations and edges to the graph.
	// lang indicates the source language of the code in the graph.
	// archData is nil if no architecture extraction was performed.
	Annotate(g *graph.CPG, lang string, archData *ArchitectureData) error

	// Dependencies returns names of domains that must run before this one.
	// Returns nil if no dependencies.
	Dependencies() []string

	// Queries returns domain-specific query rules to run against the annotated graph.
	Queries() []query.Rule
}

// Annotator is a language-specific annotator used by domain analyzers.
type Annotator interface {
	Annotate(g *graph.CPG, archData *ArchitectureData) error
}
