package domains

import (
	"fmt"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
	"github.com/ugiordan/architecture-analyzer/pkg/query"
)

// DomainResult holds the output from running a single domain analyzer.
type DomainResult struct {
	Domain           string
	AnnotationsAdded int
	Findings         []query.Finding
}

// Orchestrator manages domain analyzer execution order and results.
type Orchestrator struct {
	domains     []DomainAnalyzer
	sorted      []DomainAnalyzer
	annotCounts map[string]int
}

// NewOrchestrator creates an orchestrator for the given domain analyzers.
func NewOrchestrator(domains []DomainAnalyzer) *Orchestrator {
	return &Orchestrator{domains: domains}
}

// AnnotateAll sorts domains by dependency and runs Annotate on each.
// Must be called before RunQueries.
func (o *Orchestrator) AnnotateAll(cpg *graph.CPG, lang string, archData *ArchitectureData) error {
	sorted, err := sortByDependency(o.domains)
	if err != nil {
		return err
	}
	o.sorted = sorted
	o.annotCounts = make(map[string]int, len(sorted))

	for _, d := range o.sorted {
		before := countAnnotations(cpg)

		if err := d.Annotate(cpg, lang, archData); err != nil {
			return fmt.Errorf("domain %q annotate: %w", d.Name(), err)
		}

		o.annotCounts[d.Name()] = countAnnotations(cpg) - before
	}

	return nil
}

// RunQueries runs query rules for each domain and returns results.
// Must be called after AnnotateAll.
func (o *Orchestrator) RunQueries(cpg *graph.CPG) ([]DomainResult, error) {
	if o.sorted == nil {
		return nil, fmt.Errorf("RunQueries called before AnnotateAll")
	}

	var results []DomainResult
	for _, d := range o.sorted {
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
			AnnotationsAdded: o.annotCounts[d.Name()],
			Findings:         findings,
		})
	}

	return results, nil
}

// Run executes all domain analyzers in dependency order, then runs their queries.
// Backwards-compatible wrapper around AnnotateAll + RunQueries.
// NOTE: this does not run the TaintEngine between annotation and query phases.
// For taint-aware analysis, call AnnotateAll, then TaintEngine.Run, then RunQueries.
func (o *Orchestrator) Run(cpg *graph.CPG, lang string, archData *ArchitectureData) ([]DomainResult, error) {
	if err := o.AnnotateAll(cpg, lang, archData); err != nil {
		return nil, err
	}
	return o.RunQueries(cpg)
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
