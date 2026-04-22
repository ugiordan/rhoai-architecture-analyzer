package testing

import (
	"github.com/ugiordan/architecture-analyzer/pkg/domains"
	"github.com/ugiordan/architecture-analyzer/pkg/graph"
	"github.com/ugiordan/architecture-analyzer/pkg/query"
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

func (a *Analyzer) Name() string                  { return "testing" }
func (a *Analyzer) SupportedLanguages() []string   { return []string{"go"} }
func (a *Analyzer) Dependencies() []string          { return []string{"security"} }

func (a *Analyzer) Annotate(g *graph.CPG, lang string, archData *domains.ArchitectureData) error {
	ann, ok := a.annotators[lang]
	if !ok {
		return nil
	}
	return ann.Annotate(g, archData)
}

func (a *Analyzer) Queries() []query.Rule {
	return testingQueries()
}
