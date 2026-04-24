package security

import (
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

func TestTypeScriptAnnotatorHandlesRequest(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "getUsers",
		File:        "server.ts",
		Line:        10,
		EndLine:     20,
		Language:    "typescript",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
		ParamTypes:  []string{"req: Request", "res: Response"},
	}
	g.AddNode(fn)

	a := &TypeScriptAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("fn1").Annotations[AnnotHandlesRequest] {
		t.Error("expected sec:handles_request annotation")
	}
}

func TestTypeScriptAnnotatorAccessesSecret(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "connect",
		File: "db.ts", Line: 10, EndLine: 20,
		Language:    "typescript",
		Annotations: make(map[string]bool), Properties: make(map[string]string),
	}
	g.AddNode(fn)

	cs := &graph.Node{
		ID: "call1", Kind: graph.NodeCallSite, Name: "process.env.DB_PASSWORD",
		File:        "db.ts",
		Line:        15,
		Language:    "typescript",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(cs)
	g.AddEdge(&graph.Edge{From: "fn1", To: "call1", Kind: graph.EdgeDataFlow, Label: "contains_call"})

	a := &TypeScriptAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("call1").Annotations[AnnotAccessesSecret] {
		t.Error("expected sec:accesses_secret on call site")
	}
	if !g.GetNode("fn1").Annotations[AnnotAccessesSecret] {
		t.Error("expected sec:accesses_secret on containing function")
	}
}

func TestTypeScriptAnnotatorRendersHTML(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "sendResponse",
		File: "handler.ts", Line: 10, EndLine: 20,
		Language:    "typescript",
		Annotations: make(map[string]bool), Properties: make(map[string]string),
	}
	g.AddNode(fn)

	cs := &graph.Node{
		ID: "call1", Kind: graph.NodeCallSite, Name: "res.send",
		File:        "handler.ts",
		Line:        15,
		Language:    "typescript",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(cs)
	g.AddEdge(&graph.Edge{From: "fn1", To: "call1", Kind: graph.EdgeDataFlow, Label: "contains_call"})

	a := &TypeScriptAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("call1").Annotations[AnnotRendersHTML] {
		t.Error("expected sec:renders_html on call site")
	}
	if !g.GetNode("fn1").Annotations[AnnotRendersHTML] {
		t.Error("expected sec:renders_html on containing function")
	}
}

func TestTypeScriptAnnotatorEvalUsage(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "runCode",
		File: "eval.ts", Line: 10, EndLine: 20,
		Language:    "typescript",
		Annotations: make(map[string]bool), Properties: make(map[string]string),
	}
	g.AddNode(fn)

	cs := &graph.Node{
		ID: "call1", Kind: graph.NodeCallSite, Name: "eval",
		File:        "eval.ts",
		Line:        15,
		Language:    "typescript",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(cs)
	g.AddEdge(&graph.Edge{From: "fn1", To: "call1", Kind: graph.EdgeDataFlow, Label: "contains_call"})

	a := &TypeScriptAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("call1").Annotations[AnnotEvalUsage] {
		t.Error("expected sec:eval_usage on call site")
	}
	if !g.GetNode("fn1").Annotations[AnnotEvalUsage] {
		t.Error("expected sec:eval_usage on containing function")
	}
}

func TestTypeScriptAnnotatorNoFalsePositives(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "helper",
		File:        "utils.ts",
		Line:        10,
		Language:    "typescript",
		Annotations: make(map[string]bool),
		Properties:  map[string]string{"param_types": "(x: number)"},
	}
	g.AddNode(fn)

	a := &TypeScriptAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	n := g.GetNode("fn1")
	for _, ann := range []string{
		AnnotHandlesRequest, AnnotAccessesSecret, AnnotRendersHTML,
		AnnotEvalUsage, AnnotRedirect, AnnotExecutesSQL,
	} {
		if n.Annotations[ann] {
			t.Errorf("unexpected annotation %q on regular function", ann)
		}
	}
}
