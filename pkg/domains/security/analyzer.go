package security

import (
	"github.com/ugiordan/rhoai-architecture-analyzer/pkg/domains"
	"github.com/ugiordan/rhoai-architecture-analyzer/pkg/graph"
	"github.com/ugiordan/rhoai-architecture-analyzer/pkg/query"
)

type Analyzer struct {
	annotators map[string]domains.Annotator
}

func New() *Analyzer {
	return &Analyzer{
		annotators: map[string]domains.Annotator{
			"go": &GoAnnotator{},
		},
	}
}

func (a *Analyzer) Name() string                { return "security" }
func (a *Analyzer) SupportedLanguages() []string { return []string{"go"} }
func (a *Analyzer) Dependencies() []string       { return nil }

func (a *Analyzer) Annotate(g *graph.CPG, lang string, archData *domains.ArchitectureData) error {
	ann, ok := a.annotators[lang]
	if !ok {
		return nil
	}
	return ann.Annotate(g, archData)
}

func (a *Analyzer) Queries() []query.Rule {
	return securityQueries()
}
