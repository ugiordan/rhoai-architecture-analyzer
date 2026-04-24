package security

import (
	"regexp"
	"strings"

	"github.com/ugiordan/architecture-analyzer/pkg/domains"
	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

type PythonAnnotator struct{}

var secretPattern = regexp.MustCompile(`(?i)(secret|password|token|api[_\-]?key|credential|auth|bearer)`)

func (a *PythonAnnotator) Annotate(g *graph.CPG, archData *domains.ArchitectureData) error {
	// First pass: annotate individual call sites
	for _, cs := range g.NodesByKind(graph.NodeCallSite) {
		if cs.Language != "python" {
			continue
		}
		a.annotateCallSite(g, cs)
	}

	// Second pass: annotate functions based on decorators, parameters, and propagation
	for _, fn := range g.NodesByKind(graph.NodeFunction) {
		if fn.Language != "python" {
			continue
		}
		a.annotateFunction(g, fn)
	}

	return nil
}

func (a *PythonAnnotator) annotateFunction(g *graph.CPG, fn *graph.Node) {
	// sec:handles_request: function with route decorator or request parameter
	for _, decorator := range fn.Decorators {
		if containsRouteDecorator(decorator) {
			g.SetAnnotation(fn.ID, AnnotHandlesRequest, true)
			break
		}
	}

	paramNames := strings.Join(fn.ParamNames, ",")
	if strings.Contains(paramNames, "request") {
		g.SetAnnotation(fn.ID, AnnotHandlesRequest, true)
	}

	// Propagate annotations from contained call sites
	for _, edge := range g.OutEdges(fn.ID) {
		if edge.Kind != graph.EdgeDataFlow {
			continue
		}
		target := g.GetNode(edge.To)
		if target == nil || target.Kind != graph.NodeCallSite {
			continue
		}

		// Propagate security annotations from call sites to function
		if target.Annotations[AnnotSubprocessCall] {
			g.SetAnnotation(fn.ID, AnnotSubprocessCall, true)
		}
		if target.Annotations[AnnotDeserializesInput] {
			g.SetAnnotation(fn.ID, AnnotDeserializesInput, true)
		}
		if target.Annotations[AnnotAccessesSecret] {
			g.SetAnnotation(fn.ID, AnnotAccessesSecret, true)
		}
		if target.Annotations[AnnotExecutesSQL] {
			g.SetAnnotation(fn.ID, AnnotExecutesSQL, true)
		}
		if target.Annotations[AnnotFileAccess] {
			g.SetAnnotation(fn.ID, AnnotFileAccess, true)
		}
		if target.Annotations[AnnotTemplateRender] {
			g.SetAnnotation(fn.ID, AnnotTemplateRender, true)
		}
	}
}

func (a *PythonAnnotator) annotateCallSite(g *graph.CPG, cs *graph.Node) {
	name := cs.Name

	// sec:subprocess_call
	if isSubprocessCall(name) {
		g.SetAnnotation(cs.ID, AnnotSubprocessCall, true)
	}

	// sec:deserializes_input
	if isDeserializationCall(name) {
		g.SetAnnotation(cs.ID, AnnotDeserializesInput, true)
	}

	// sec:accesses_secret
	if isEnvironAccessCall(name) {
		stringArgs := cs.Properties["string_args"]
		if secretPattern.MatchString(stringArgs) {
			g.SetAnnotation(cs.ID, AnnotAccessesSecret, true)
		}
	}

	// sec:executes_sql
	if isSQLExecuteCall(name) {
		g.SetAnnotation(cs.ID, AnnotExecutesSQL, true)
	}

	// sec:file_access
	if isFileAccessCall(name) {
		g.SetAnnotation(cs.ID, AnnotFileAccess, true)
	}

	// sec:template_render
	if isTemplateRenderCall(name) {
		g.SetAnnotation(cs.ID, AnnotTemplateRender, true)
	}
}

func containsRouteDecorator(decorator string) bool {
	for _, pattern := range []string{".route(", ".get(", ".post(", ".put(", ".delete(", ".patch("} {
		if strings.Contains(decorator, pattern) {
			return true
		}
	}
	return false
}

func isSubprocessCall(name string) bool {
	for _, call := range []string{"subprocess.run", "subprocess.Popen", "subprocess.call", "subprocess.check_output", "subprocess.check_call", "os.system", "os.popen"} {
		if name == call {
			return true
		}
	}
	return false
}

func isDeserializationCall(name string) bool {
	for _, call := range []string{"pickle.loads", "pickle.load", "yaml.load", "yaml.unsafe_load", "marshal.loads", "marshal.load"} {
		if name == call {
			return true
		}
	}
	return false
}

func isEnvironAccessCall(name string) bool {
	return strings.HasPrefix(name, "os.environ") || name == "os.getenv"
}

func isSQLExecuteCall(name string) bool {
	for _, call := range []string{"cursor.execute", "session.execute", "db.execute", "connection.execute"} {
		if name == call {
			return true
		}
	}
	return false
}

func isFileAccessCall(name string) bool {
	if name == "open" {
		return true
	}
	for _, suffix := range []string{".read_text", ".read_bytes"} {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}
	return false
}

func isTemplateRenderCall(name string) bool {
	return name == "render_template_string"
}
