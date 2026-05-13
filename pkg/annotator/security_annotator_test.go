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
	if err := cpg.AddNode(fn); err != nil { t.Fatal(err) }

	dbWrite := &graph.Node{
		ID:         "db1",
		Kind:       graph.NodeDBOperation,
		Name:       "db.Exec",
		File:       "server.go",
		Line:       15,
		Operation:  "write",
		Properties: map[string]string{"operation": "write"},
	}
	if err := cpg.AddNode(dbWrite); err != nil { t.Fatal(err) }
	cpg.AddEdge(&graph.Edge{From: "fn1", To: "db1", Kind: graph.EdgeDataFlow})

	a := NewSecurityAnnotator()
	a.Annotate(cpg)

	updated := cpg.GetNode("fn1")
	if !updated.Annotations["handles_user_input"] {
		t.Error("expected handles_user_input annotation")
	}
	if !updated.Annotations["writes_storage"] {
		t.Error("expected writes_storage propagated from DBOperation to parent function")
	}
	if !updated.Annotations["mutates_state"] {
		t.Error("expected mutates_state propagated from DBOperation to parent function")
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
	if err := cpg.AddNode(call); err != nil { t.Fatal(err) }

	a := NewSecurityAnnotator()
	a.Annotate(cpg)

	updated := cpg.GetNode("call1")
	if !updated.Annotations[annotCallsExternal] {
		t.Error("expected annotCallsExternal annotation on http.Post call site")
	}
}

func TestAnnotateExternalCall_Propagation(t *testing.T) {
	cpg := graph.NewCPG()

	fn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "doRequest",
		File:        "client.go",
		Line:        1,
		Properties:  make(map[string]string),
		Annotations: make(map[string]bool),
	}
	call := &graph.Node{
		ID:          "call1",
		Kind:        graph.NodeCallSite,
		Name:        "http.Post",
		File:        "client.go",
		Line:        5,
		Properties:  make(map[string]string),
		Annotations: make(map[string]bool),
	}
	if err := cpg.AddNode(fn); err != nil { t.Fatal(err) }
	if err := cpg.AddNode(call); err != nil { t.Fatal(err) }
	cpg.AddEdge(&graph.Edge{From: "fn1", To: "call1", Kind: graph.EdgeDataFlow})

	a := NewSecurityAnnotator()
	a.Annotate(cpg)

	updatedFn := cpg.GetNode("fn1")
	if !updatedFn.Annotations[annotCallsExternal] {
		t.Error("expected annotCallsExternal propagated to parent function")
	}

	updatedCall := cpg.GetNode("call1")
	if !updatedCall.Annotations[annotCallsExternal] {
		t.Error("expected annotCallsExternal on call site")
	}
}

func TestAnnotateCrossesNamespace_Propagation(t *testing.T) {
	cpg := graph.NewCPG()

	fn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "reconcile",
		File:        "controller.go",
		Line:        1,
		Properties:  make(map[string]string),
		Annotations: make(map[string]bool),
	}
	call := &graph.Node{
		ID:          "call1",
		Kind:        graph.NodeCallSite,
		Name:        "client.get.namespace",
		File:        "controller.go",
		Line:        10,
		Properties:  make(map[string]string),
		Annotations: make(map[string]bool),
	}
	if err := cpg.AddNode(fn); err != nil { t.Fatal(err) }
	if err := cpg.AddNode(call); err != nil { t.Fatal(err) }
	cpg.AddEdge(&graph.Edge{From: "fn1", To: "call1", Kind: graph.EdgeDataFlow})

	a := NewSecurityAnnotator()
	a.Annotate(cpg)

	updatedFn := cpg.GetNode("fn1")
	if !updatedFn.Annotations[annotCrossesNamespace] {
		t.Error("expected annotCrossesNamespace propagated to parent function")
	}
}
