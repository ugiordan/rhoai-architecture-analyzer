package security

import (
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

func TestPythonAnnotatorHandlesRequest(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "get_users",
		File:        "app.py",
		Line:        10,
		EndLine:     20,
		Language:    "python",
		Decorators:  []string{"app.route(\"/users\")"},
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(fn)

	a := &PythonAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("fn1").Annotations[AnnotHandlesRequest] {
		t.Error("expected sec:handles_request on decorated route handler")
	}
}

func TestPythonAnnotatorHandlesRequestWithRequestParam(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "handle_webhook",
		File:        "webhook.py",
		Line:        10,
		EndLine:     20,
		Language:    "python",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
		ParamNames:  []string{"request", "data"},
	}
	g.AddNode(fn)

	a := &PythonAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("fn1").Annotations[AnnotHandlesRequest] {
		t.Error("expected sec:handles_request on function with request parameter")
	}
}

func TestPythonAnnotatorAccessesSecret(t *testing.T) {
	g := graph.NewCPG()
	cs := &graph.Node{
		ID:          "call1",
		Kind:        graph.NodeCallSite,
		Name:        "os.environ",
		File:        "config.py",
		Line:        15,
		Language:    "python",
		Annotations: make(map[string]bool),
		Properties:  map[string]string{"string_args": "SECRET_KEY"},
	}
	g.AddNode(cs)

	a := &PythonAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("call1").Annotations[AnnotAccessesSecret] {
		t.Error("expected sec:accesses_secret on os.environ call with SECRET_KEY")
	}
}

func TestPythonAnnotatorAccessesSecretCaseInsensitive(t *testing.T) {
	g := graph.NewCPG()
	cs := &graph.Node{
		ID:          "call1",
		Kind:        graph.NodeCallSite,
		Name:        "os.getenv",
		File:        "config.py",
		Line:        15,
		Language:    "python",
		Annotations: make(map[string]bool),
		Properties:  map[string]string{"string_args": "api-key"},
	}
	g.AddNode(cs)

	a := &PythonAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("call1").Annotations[AnnotAccessesSecret] {
		t.Error("expected sec:accesses_secret on os.getenv call with api-key")
	}
}

func TestPythonAnnotatorSubprocessCall(t *testing.T) {
	g := graph.NewCPG()
	cs := &graph.Node{
		ID:          "call1",
		Kind:        graph.NodeCallSite,
		Name:        "subprocess.run",
		File:        "runner.py",
		Line:        20,
		Language:    "python",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(cs)

	a := &PythonAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("call1").Annotations[AnnotSubprocessCall] {
		t.Error("expected sec:subprocess_call on subprocess.run")
	}
}

func TestPythonAnnotatorDeserializesInput(t *testing.T) {
	g := graph.NewCPG()
	cs := &graph.Node{
		ID:          "call1",
		Kind:        graph.NodeCallSite,
		Name:        "pickle.loads",
		File:        "serializer.py",
		Line:        25,
		Language:    "python",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(cs)

	a := &PythonAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("call1").Annotations[AnnotDeserializesInput] {
		t.Error("expected sec:deserializes_input on pickle.loads")
	}
}

func TestPythonAnnotatorExecutesSQL(t *testing.T) {
	g := graph.NewCPG()
	cs := &graph.Node{
		ID:          "call1",
		Kind:        graph.NodeCallSite,
		Name:        "cursor.execute",
		File:        "database.py",
		Line:        30,
		Language:    "python",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(cs)

	a := &PythonAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("call1").Annotations[AnnotExecutesSQL] {
		t.Error("expected sec:executes_sql on cursor.execute")
	}
}

func TestPythonAnnotatorFileAccess(t *testing.T) {
	g := graph.NewCPG()
	cs := &graph.Node{
		ID:          "call1",
		Kind:        graph.NodeCallSite,
		Name:        "open",
		File:        "reader.py",
		Line:        35,
		Language:    "python",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(cs)

	a := &PythonAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("call1").Annotations[AnnotFileAccess] {
		t.Error("expected sec:file_access on open")
	}
}

func TestPythonAnnotatorTemplateRender(t *testing.T) {
	g := graph.NewCPG()
	cs := &graph.Node{
		ID:          "call1",
		Kind:        graph.NodeCallSite,
		Name:        "render_template_string",
		File:        "views.py",
		Line:        40,
		Language:    "python",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(cs)

	a := &PythonAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("call1").Annotations[AnnotTemplateRender] {
		t.Error("expected sec:template_render on render_template_string")
	}
}

func TestPythonAnnotatorPropagateToFunction(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "process_data",
		File:        "handler.py",
		Line:        10,
		EndLine:     30,
		Language:    "python",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(fn)

	cs := &graph.Node{
		ID:          "call1",
		Kind:        graph.NodeCallSite,
		Name:        "subprocess.run",
		File:        "handler.py",
		Line:        20,
		Language:    "python",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(cs)
	g.AddEdge(&graph.Edge{From: "fn1", To: "call1", Kind: graph.EdgeDataFlow, Label: "contains_call"})

	a := &PythonAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("call1").Annotations[AnnotSubprocessCall] {
		t.Error("expected sec:subprocess_call on call site")
	}
	if !g.GetNode("fn1").Annotations[AnnotSubprocessCall] {
		t.Error("expected sec:subprocess_call propagated to containing function")
	}
}

func TestPythonAnnotatorNoFalsePositives(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "helper_function",
		File:        "utils.py",
		Line:        10,
		Language:    "python",
		Annotations: make(map[string]bool),
		Properties:  map[string]string{"param_names": "data,config"},
	}
	g.AddNode(fn)

	a := &PythonAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	n := g.GetNode("fn1")
	for _, ann := range []string{
		AnnotHandlesRequest, AnnotAccessesSecret, AnnotExecutesSQL,
		AnnotDeserializesInput, AnnotSubprocessCall, AnnotFileAccess,
		AnnotTemplateRender,
	} {
		if n.Annotations[ann] {
			t.Errorf("unexpected annotation %q on regular function", ann)
		}
	}
}

func TestPythonAnnotatorIgnoresOtherLanguages(t *testing.T) {
	g := graph.NewCPG()
	cs := &graph.Node{
		ID:          "call1",
		Kind:        graph.NodeCallSite,
		Name:        "subprocess.run",
		File:        "test.go",
		Line:        10,
		Language:    "go",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(cs)

	a := &PythonAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if g.GetNode("call1").Annotations[AnnotSubprocessCall] {
		t.Error("should not annotate Go nodes with Python security annotations")
	}
}
