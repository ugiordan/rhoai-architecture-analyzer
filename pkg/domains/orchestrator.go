package domains

import (
	"fmt"

	"github.com/ugiordan/rhoai-architecture-analyzer/pkg/graph"
	"github.com/ugiordan/rhoai-architecture-analyzer/pkg/query"
)

// DomainResult holds the output from running a single domain analyzer.
type DomainResult struct {
	Domain           string
	AnnotationsAdded int
	Findings         []query.Finding
}

// Orchestrator manages domain analyzer execution order and results.
type Orchestrator struct {
	domains []DomainAnalyzer
}

// NewOrchestrator creates an orchestrator for the given domain analyzers.
func NewOrchestrator(domains []DomainAnalyzer) *Orchestrator {
	return &Orchestrator{domains: domains}
}

// Run executes all domain analyzers in dependency order, then runs their queries.
func (o *Orchestrator) Run(cpg *graph.CPG, lang string, archData *ArchitectureData) ([]DomainResult, error) {
	sorted, err := sortByDependency(o.domains)
	if err != nil {
		return nil, err
	}

	var results []DomainResult
	for _, d := range sorted {
		annotationsBefore := countAnnotations(cpg)

		if err := d.Annotate(cpg, lang, archData); err != nil {
			return nil, fmt.Errorf("domain %q annotate: %w", d.Name(), err)
		}

		annotationsAfter := countAnnotations(cpg)

		var findings []query.Finding
		for _, rule := range d.Queries() {
			ruleFindings := rule.Run(cpg)
			for i := range ruleFindings {
				ruleFindings[i].Domain = d.Name()
			}
			findings = append(findings, ruleFindings...)
		}

		results = append(results, DomainResult{
			Domain:           d.Name(),
			AnnotationsAdded: annotationsAfter - annotationsBefore,
			Findings:         findings,
		})
	}

	return results, nil
}

func countAnnotations(cpg *graph.CPG) int {
	count := 0
	for _, n := range cpg.Nodes() {
		count += len(n.Annotations)
	}
	return count
}

// sortByDependency performs a topological sort of domain analyzers using Kahn's algorithm.
func sortByDependency(domains []DomainAnalyzer) ([]DomainAnalyzer, error) {
	byName := make(map[string]DomainAnalyzer)
	for _, d := range domains {
		byName[d.Name()] = d
	}

	// Validate all dependencies are present
	for _, d := range domains {
		for _, dep := range d.Dependencies() {
			if _, ok := byName[dep]; !ok {
				return nil, fmt.Errorf("domain %q requires missing dependency %q", d.Name(), dep)
			}
		}
	}

	// Compute in-degrees
	inDegree := make(map[string]int)
	for _, d := range domains {
		if _, ok := inDegree[d.Name()]; !ok {
			inDegree[d.Name()] = 0
		}
		inDegree[d.Name()] += len(d.Dependencies())
	}

	// Start with zero in-degree nodes
	var queue []string
	for name, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, name)
		}
	}

	var sorted []DomainAnalyzer
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		sorted = append(sorted, byName[name])

		for _, d := range domains {
			for _, dep := range d.Dependencies() {
				if dep == name {
					inDegree[d.Name()]--
					if inDegree[d.Name()] == 0 {
						queue = append(queue, d.Name())
					}
				}
			}
		}
	}

	if len(sorted) != len(domains) {
		return nil, fmt.Errorf("cyclic dependency detected among domains")
	}
	return sorted, nil
}
