package security

import (
	"strings"

	"github.com/ugiordan/architecture-analyzer/pkg/domains"
	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

type RustAnnotator struct{}

func (a *RustAnnotator) Annotate(g *graph.CPG, archData *domains.ArchitectureData) error {
	// Build set of HTTP endpoint names for handles_request matching
	endpointNames := make(map[string]bool)
	for _, ep := range g.NodesByKind(graph.NodeHTTPEndpoint) {
		if ep.Language == "rust" {
			endpointNames[ep.Name] = true
		}
	}

	// First pass: annotate individual call sites
	for _, cs := range g.NodesByKind(graph.NodeCallSite) {
		if cs.Language != "rust" {
			continue
		}
		a.annotateCallSite(g, cs)
	}

	// Second pass: annotate functions based on properties and contained nodes
	for _, fn := range g.NodesByKind(graph.NodeFunction) {
		if fn.Language != "rust" {
			continue
		}
		a.annotateFunction(g, fn, endpointNames)
	}

	return nil
}

func (a *RustAnnotator) annotateCallSite(g *graph.CPG, cs *graph.Node) {
	name := cs.Name
	stringArgs := cs.Properties["string_args"]

	// sec:accesses_secret: env::var with secret pattern
	if strings.Contains(name, "env::var") && secretPattern.MatchString(stringArgs) {
		g.SetAnnotation(cs.ID, AnnotAccessesSecret, true)
	}

	// sec:deserializes_input: serde_json deserialization
	if strings.Contains(name, "serde_json::from_str") || strings.Contains(name, "serde_json::from_slice") {
		g.SetAnnotation(cs.ID, AnnotDeserializesInput, true)
	}

	// sec:command_execution: process command execution
	if strings.Contains(name, "Command::new") || strings.Contains(name, "process::Command") {
		g.SetAnnotation(cs.ID, AnnotCommandExecution, true)
	}
}

func (a *RustAnnotator) annotateFunction(g *graph.CPG, fn *graph.Node, endpointNames map[string]bool) {
	// sec:unsafe_block: functions with is_unsafe property
	if fn.IsUnsafe {
		g.SetAnnotation(fn.ID, AnnotUnsafeBlock, true)
	}

	// sec:ffi_call: functions with is_extern property
	if fn.IsExtern {
		g.SetAnnotation(fn.ID, AnnotFFICall, true)
	}

	// sec:handles_request: function name matches HTTP endpoint
	if endpointNames[fn.Name] {
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

		if target.Annotations[AnnotAccessesSecret] {
			g.SetAnnotation(fn.ID, AnnotAccessesSecret, true)
		}
		if target.Annotations[AnnotCommandExecution] {
			g.SetAnnotation(fn.ID, AnnotCommandExecution, true)
		}
		if target.Annotations[AnnotDeserializesInput] {
			g.SetAnnotation(fn.ID, AnnotDeserializesInput, true)
		}
	}
}
