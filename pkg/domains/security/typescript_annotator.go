package security

import (
	"strings"

	"github.com/ugiordan/architecture-analyzer/pkg/domains"
	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

type TypeScriptAnnotator struct{}

func (a *TypeScriptAnnotator) Annotate(g *graph.CPG, archData *domains.ArchitectureData) error {
	// First pass: annotate call sites
	for _, cs := range g.NodesByKind(graph.NodeCallSite) {
		if cs.Language != "typescript" {
			continue
		}
		a.annotateCallSite(g, cs)
	}

	// Second pass: annotate functions based on parameters and contained nodes
	for _, fn := range g.NodesByKind(graph.NodeFunction) {
		if fn.Language != "typescript" {
			continue
		}
		a.annotateFunction(g, fn)
	}

	return nil
}

func (a *TypeScriptAnnotator) annotateFunction(g *graph.CPG, fn *graph.Node) {
	// sec:handles_request: function with (req, res) or (request, response) parameters
	paramTypes := strings.ToLower(strings.Join(fn.ParamTypes, ","))
	if (strings.Contains(paramTypes, "req") && strings.Contains(paramTypes, "res")) ||
		(strings.Contains(paramTypes, "request") && strings.Contains(paramTypes, "response")) {
		g.SetAnnotation(fn.ID, AnnotHandlesRequest, true)
	}

	// Propagate annotations from contained call sites via EdgeDataFlow
	for _, edge := range g.OutEdges(fn.ID) {
		if edge.Kind != graph.EdgeDataFlow {
			continue
		}
		target := g.GetNode(edge.To)
		if target == nil || target.Kind != graph.NodeCallSite {
			continue
		}

		// Propagate security-relevant annotations
		if target.Annotations[AnnotAccessesSecret] {
			g.SetAnnotation(fn.ID, AnnotAccessesSecret, true)
		}
		if target.Annotations[AnnotRendersHTML] {
			g.SetAnnotation(fn.ID, AnnotRendersHTML, true)
		}
		if target.Annotations[AnnotEvalUsage] {
			g.SetAnnotation(fn.ID, AnnotEvalUsage, true)
		}
		if target.Annotations[AnnotExecutesSQL] {
			g.SetAnnotation(fn.ID, AnnotExecutesSQL, true)
		}
		if target.Annotations[AnnotRedirect] {
			g.SetAnnotation(fn.ID, AnnotRedirect, true)
		}
	}
}

func (a *TypeScriptAnnotator) annotateCallSite(g *graph.CPG, cs *graph.Node) {
	name := cs.Name

	// sec:accesses_secret: process.env.* where the env var matches secret pattern
	if strings.HasPrefix(name, "process.env.") {
		envVar := strings.TrimPrefix(name, "process.env.")
		if secretPattern.MatchString(envVar) {
			g.SetAnnotation(cs.ID, AnnotAccessesSecret, true)
		}
	}

	// sec:renders_html: res.send, .innerHTML, dangerouslySetInnerHTML
	if name == "res.send" || strings.HasSuffix(name, ".innerHTML") || strings.Contains(name, "dangerouslySetInnerHTML") {
		g.SetAnnotation(cs.ID, AnnotRendersHTML, true)
	}

	// sec:executes_sql: calls ending with .query or .execute
	if strings.HasSuffix(name, ".query") || strings.HasSuffix(name, ".execute") {
		g.SetAnnotation(cs.ID, AnnotExecutesSQL, true)
	}

	// sec:eval_usage: eval, Function constructor
	if name == "eval" || name == "Function" {
		g.SetAnnotation(cs.ID, AnnotEvalUsage, true)
	}

	// sec:redirect: res.redirect
	if name == "res.redirect" {
		g.SetAnnotation(cs.ID, AnnotRedirect, true)
	}
}
