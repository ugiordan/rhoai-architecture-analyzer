package annotator

import (
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

func TestAnnotateHTTPHandler(t *testing.T) {
	cpg := graph.NewCPG()
	fn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "handleRequest",
		File:        "server.go",
		Line:        10,
		Annotations: map[string]bool{"handles_user_input": true},
		Properties:  map[string]string{"handler_type": "http"},
	}
	cpg.AddNode(fn)

	dbWrite := &graph.Node{
		ID:         "db1",
		Kind:       graph.NodeDBOperation,
		Name:       "db.Exec",
		File:       "server.go",
		Line:       15,
		Operation:  "write",
		Properties: map[string]string{"operation": "write"},
	}
	cpg.AddNode(dbWrite)
	cpg.AddEdge(&graph.Edge{From: "fn1", To: "db1", Kind: graph.EdgeDataFlow})

	a := NewSecurityAnnotator()
	a.Annotate(cpg)

	updated := cpg.GetNode("fn1")
	if !updated.Annotations["handles_user_input"] {
		t.Error("expected handles_user_input annotation")
	}
}

func TestAnnotateExternalCall(t *testing.T) {
	cpg := graph.NewCPG()

	call := &graph.Node{
		ID:          "call1",
		Kind:        graph.NodeCallSite,
		Name:        "http.Post",
		File:        "client.go",
		Line:        5,
		Properties:  make(map[string]string),
		Annotations: make(map[string]bool),
	}
	cpg.AddNode(call)

	a := NewSecurityAnnotator()
	a.Annotate(cpg)

	updated := cpg.GetNode("call1")
	if !updated.Annotations["calls_external"] {
		t.Error("expected calls_external annotation on http.Post")
	}
}
