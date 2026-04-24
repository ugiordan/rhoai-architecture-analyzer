package security

import (
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

func TestRustAnnotatorUnsafeBlock(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "raw_op",
		File:        "lib.rs",
		Line:        10,
		EndLine:     20,
		Language:    "rust",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
		IsUnsafe:    true,
	}
	g.AddNode(fn)

	a := &RustAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("fn1").Annotations[AnnotUnsafeBlock] {
		t.Error("expected sec:unsafe_block on unsafe function")
	}
}

func TestRustAnnotatorFFICall(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "c_wrapper",
		File:        "ffi.rs",
		Line:        5,
		EndLine:     15,
		Language:    "rust",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
		IsExtern:    true,
	}
	g.AddNode(fn)

	a := &RustAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("fn1").Annotations[AnnotFFICall] {
		t.Error("expected sec:ffi_call on extern function")
	}
}

func TestRustAnnotatorHandlesRequest(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "list_items",
		File:        "handler.rs",
		Line:        10,
		EndLine:     30,
		Language:    "rust",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(fn)

	handler := &graph.Node{
		ID:          "http1",
		Kind:        graph.NodeHTTPEndpoint,
		Name:        "list_items",
		File:        "handler.rs",
		Language:    "rust",
		Annotations: make(map[string]bool),
		Properties:  map[string]string{"route": "/items"},
	}
	g.AddNode(handler)

	a := &RustAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("fn1").Annotations[AnnotHandlesRequest] {
		t.Error("expected sec:handles_request on function matching HTTP endpoint")
	}
}

func TestRustAnnotatorAccessesSecret(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "get_config",
		File:        "config.rs",
		Line:        5,
		EndLine:     15,
		Language:    "rust",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(fn)

	call := &graph.Node{
		ID:          "call1",
		Kind:        graph.NodeCallSite,
		Name:        "std::env::var",
		File:        "config.rs",
		Line:        8,
		Language:    "rust",
		Annotations: make(map[string]bool),
		Properties:  map[string]string{"string_args": "API_KEY"},
	}
	g.AddNode(call)

	g.AddEdge(&graph.Edge{
		Kind: graph.EdgeDataFlow,
		From: fn.ID,
		To:   call.ID,
	})

	a := &RustAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("call1").Annotations[AnnotAccessesSecret] {
		t.Error("expected sec:accesses_secret on env::var call with API_KEY")
	}

	if !g.GetNode("fn1").Annotations[AnnotAccessesSecret] {
		t.Error("expected sec:accesses_secret propagated to containing function")
	}
}

func TestRustAnnotatorCommandExecution(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "run_task",
		File:        "executor.rs",
		Line:        10,
		EndLine:     25,
		Language:    "rust",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(fn)

	call := &graph.Node{
		ID:          "call1",
		Kind:        graph.NodeCallSite,
		Name:        "std::process::Command::new",
		File:        "executor.rs",
		Line:        15,
		Language:    "rust",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(call)

	g.AddEdge(&graph.Edge{
		Kind: graph.EdgeDataFlow,
		From: fn.ID,
		To:   call.ID,
	})

	a := &RustAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("call1").Annotations[AnnotCommandExecution] {
		t.Error("expected sec:command_execution on Command::new call")
	}

	if !g.GetNode("fn1").Annotations[AnnotCommandExecution] {
		t.Error("expected sec:command_execution propagated to containing function")
	}
}

func TestRustAnnotatorDeserializesInput(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "parse_json",
		File:        "parser.rs",
		Line:        5,
		EndLine:     20,
		Language:    "rust",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(fn)

	call := &graph.Node{
		ID:          "call1",
		Kind:        graph.NodeCallSite,
		Name:        "serde_json::from_str",
		File:        "parser.rs",
		Line:        10,
		Language:    "rust",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(call)

	g.AddEdge(&graph.Edge{
		Kind: graph.EdgeDataFlow,
		From: fn.ID,
		To:   call.ID,
	})

	a := &RustAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("call1").Annotations[AnnotDeserializesInput] {
		t.Error("expected sec:deserializes_input on serde_json::from_str call")
	}

	if !g.GetNode("fn1").Annotations[AnnotDeserializesInput] {
		t.Error("expected sec:deserializes_input propagated to containing function")
	}
}

func TestRustAnnotatorNoFalsePositives(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "helper_function",
		File:        "utils.rs",
		Line:        10,
		EndLine:     20,
		Language:    "rust",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(fn)

	call := &graph.Node{
		ID:          "call1",
		Kind:        graph.NodeCallSite,
		Name:        "println!",
		File:        "utils.rs",
		Line:        15,
		Language:    "rust",
		Annotations: make(map[string]bool),
		Properties:  map[string]string{"string_args": "Hello world"},
	}
	g.AddNode(call)

	g.AddEdge(&graph.Edge{
		Kind: graph.EdgeDataFlow,
		From: fn.ID,
		To:   call.ID,
	})

	a := &RustAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	node := g.GetNode("fn1")
	if node.Annotations[AnnotUnsafeBlock] ||
		node.Annotations[AnnotFFICall] ||
		node.Annotations[AnnotHandlesRequest] ||
		node.Annotations[AnnotAccessesSecret] ||
		node.Annotations[AnnotCommandExecution] ||
		node.Annotations[AnnotDeserializesInput] {
		t.Error("expected no security annotations on plain helper function")
	}

	callNode := g.GetNode("call1")
	if callNode.Annotations[AnnotAccessesSecret] ||
		callNode.Annotations[AnnotCommandExecution] ||
		callNode.Annotations[AnnotDeserializesInput] {
		t.Error("expected no security annotations on println! call")
	}
}
