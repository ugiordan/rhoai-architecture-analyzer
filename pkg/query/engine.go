package query

import "github.com/ugiordan/architecture-analyzer/pkg/graph"

// Finding represents a security issue detected by a query rule.
type Finding struct {
	RuleID          string   `json:"rule_id"`
	Severity        string   `json:"severity"`
	Message         string   `json:"message"`
	File            string   `json:"file"`
	Line            int      `json:"line"`
	NodeID          string   `json:"node_id"`
	Path            []string `json:"path,omitempty"`
	Domain          string   `json:"domain,omitempty"`
	ArchitectureRef string   `json:"architecture_ref,omitempty"`
}

// Rule is a named query that can detect patterns in the annotated graph.
type Rule struct {
	ID       string
	Name     string
	Domain   string
	Severity string
	Run      func(cpg *graph.CPG) []Finding
}

// Engine runs security analysis queries against a code property graph.
type Engine struct{}

// NewEngine creates a query engine.
func NewEngine() *Engine { return &Engine{} }

// QueryMissingAuth finds HTTP handlers that accept user input and mutate state without auth checks.
func (e *Engine) QueryMissingAuth(cpg *graph.CPG) []Finding {
	var findings []Finding
	for _, fn := range cpg.NodesByKind(graph.NodeFunction) {
		if fn.Annotations == nil {
			continue
		}
		if fn.Annotations["handles_user_input"] &&
			fn.Annotations["mutates_state"] &&
			!fn.Annotations["has_auth"] {
			findings = append(findings, Finding{
				RuleID:   "CGA-001",
				Severity: "high",
				Message:  "HTTP handler accepts user input and mutates state without authentication: " + fn.Name,
				File:     fn.File,
				Line:     fn.Line,
				NodeID:   fn.ID,
			})
		}
	}
	return findings
}

// QueryCrossStorageTaint traces data flow from user input through storage boundaries to external sinks.
func (e *Engine) QueryCrossStorageTaint(cpg *graph.CPG) []Finding {
	var findings []Finding

	var inputNodes []*graph.Node
	for _, n := range cpg.Nodes() {
		if n.Annotations != nil && n.Annotations["handles_user_input"] {
			inputNodes = append(inputNodes, n)
		}
	}

	for _, input := range inputNodes {
		paths := traceToExternalSink(cpg, input.ID, 10)
		for _, path := range paths {
			msg := "User input flows through storage boundary to external sink"
			if len(paths) >= maxPaths {
				msg += " (results truncated, more paths may exist)"
			}
			findings = append(findings, Finding{
				RuleID:   "CGA-002",
				Severity: "critical",
				Message:  msg,
				File:     input.File,
				Line:     input.Line,
				NodeID:   input.ID,
				Path:     path,
			})
		}
	}

	return findings
}

// RunRules executes a set of domain-specific query rules against the graph.
func (e *Engine) RunRules(cpg *graph.CPG, rules []Rule) []Finding {
	var all []Finding
	for _, rule := range rules {
		all = append(all, rule.Run(cpg)...)
	}
	return all
}

// RunAll executes all security queries and returns combined findings.
func (e *Engine) RunAll(cpg *graph.CPG) []Finding {
	var all []Finding
	all = append(all, e.QueryMissingAuth(cpg)...)
	all = append(all, e.QueryCrossStorageTaint(cpg)...)
	return all
}
