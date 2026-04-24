package annotator

import (
	"strings"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

type SecurityAnnotator struct{}

func NewSecurityAnnotator() *SecurityAnnotator {
	return &SecurityAnnotator{}
}

func (sa *SecurityAnnotator) Annotate(cpg *graph.CPG) {
	for _, node := range cpg.Nodes() {
		cpg.EnsureAnnotations(node.ID)

		switch node.Kind {
		case graph.NodeFunction:
			sa.annotateFunction(node, cpg)
		case graph.NodeCallSite:
			sa.annotateCallSite(node, cpg)
		case graph.NodeDBOperation:
			sa.annotateDBOp(node, cpg)
		}
	}
}

func (sa *SecurityAnnotator) annotateFunction(fn *graph.Node, cpg *graph.CPG) {
	for _, dec := range fn.Decorators {
		lower := strings.ToLower(dec)
		if strings.Contains(lower, "auth") || strings.Contains(lower, "login_required") ||
			strings.Contains(lower, "require_admin") || strings.Contains(lower, "authenticated") {
			cpg.SetAnnotation(fn.ID, "has_auth", true)
		}
		if strings.Contains(lower, "rate_limit") || strings.Contains(lower, "limiter") {
			cpg.SetAnnotation(fn.ID, "has_rate_limit", true)
		}
	}

	// Check direct edges and one level of transitive edges (function -> call_site -> db_op)
	for _, edge := range cpg.OutEdges(fn.ID) {
		target := cpg.GetNode(edge.To)
		if target == nil {
			continue
		}
		sa.checkDBAnnotation(fn.ID, target, cpg)
		// Follow one more hop for call sites that contain DB operations
		if target.Kind == graph.NodeCallSite {
			for _, innerEdge := range cpg.OutEdges(target.ID) {
				innerTarget := cpg.GetNode(innerEdge.To)
				if innerTarget != nil {
					sa.checkDBAnnotation(fn.ID, innerTarget, cpg)
				}
			}
		}
	}
}

func (sa *SecurityAnnotator) checkDBAnnotation(fnID string, target *graph.Node, cpg *graph.CPG) {
	if target.Kind == graph.NodeDBOperation {
		op := target.Operation
		if op == "write" {
			cpg.SetAnnotation(fnID, "writes_storage", true)
			cpg.SetAnnotation(fnID, "mutates_state", true)
		} else if op == "read" {
			cpg.SetAnnotation(fnID, "reads_storage", true)
		}
	}
}

func (sa *SecurityAnnotator) annotateCallSite(cs *graph.Node, cpg *graph.CPG) {
	name := strings.ToLower(cs.Name)
	if strings.HasPrefix(name, "http.") && (strings.Contains(name, "post") ||
		strings.Contains(name, "get") || strings.Contains(name, "do")) {
		cpg.SetAnnotation(cs.ID, "calls_external", true)
	}
	if strings.Contains(name, "client.do") || strings.Contains(name, "client.post") ||
		strings.Contains(name, "client.get") {
		cpg.SetAnnotation(cs.ID, "calls_external", true)
	}
}

func (sa *SecurityAnnotator) annotateDBOp(op *graph.Node, cpg *graph.CPG) {
	switch op.Operation {
	case "write":
		cpg.SetAnnotation(op.ID, "writes_storage", true)
	case "read":
		cpg.SetAnnotation(op.ID, "reads_storage", true)
	}
}
