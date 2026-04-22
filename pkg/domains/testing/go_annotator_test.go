package testing

import (
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

func TestGoAnnotatorTestFunc(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "TestReconcile",
		File:        "controller_test.go",
		Line:        10,
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(fn)

	a := &GoAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("fn1").Annotations[AnnotIsTestFunc] {
		t.Error("expected test:is_test_func annotation")
	}
}

func TestGoAnnotatorNonTestFile(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "TestReconcile",
		File:        "controller.go",
		Line:        10,
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(fn)

	a := &GoAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if g.GetNode("fn1").Annotations[AnnotIsTestFunc] {
		t.Error("should not annotate test func in non-test file")
	}
}

func TestGoAnnotatorFakeClient(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "TestReconcile",
		File:        "controller_test.go",
		Line:        10,
		EndLine:     30,
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(fn)

	cs := &graph.Node{
		ID:          "call1",
		Kind:        graph.NodeCallSite,
		Name:        "fake.NewClientBuilder",
		File:        "controller_test.go",
		Line:        15,
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(cs)
	g.AddEdge(&graph.Edge{From: "fn1", To: "call1", Kind: graph.EdgeDataFlow, Label: "contains_call"})

	a := &GoAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("fn1").Annotations[AnnotUsesFakeClient] {
		t.Error("expected test:uses_fake_client annotation")
	}
}
