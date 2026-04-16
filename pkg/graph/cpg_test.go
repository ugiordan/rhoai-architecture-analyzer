package graph

import "testing"

func TestNewCPG(t *testing.T) {
	cpg := NewCPG()
	if cpg == nil {
		t.Fatal("NewCPG returned nil")
	}
	if len(cpg.Nodes()) != 0 {
		t.Errorf("expected 0 nodes, got %d", len(cpg.Nodes()))
	}
	if len(cpg.Edges()) != 0 {
		t.Errorf("expected 0 edges, got %d", len(cpg.Edges()))
	}
}

func TestAddNodeAndEdge(t *testing.T) {
	cpg := NewCPG()
	fn := &Node{ID: "fn1", Kind: NodeFunction, Name: "handleRequest", File: "server.go", Line: 10, Language: "go"}
	cpg.AddNode(fn)
	param := &Node{ID: "param1", Kind: NodeParameter, Name: "r", File: "server.go", Line: 10, Language: "go"}
	cpg.AddNode(param)
	cpg.AddEdge(&Edge{From: "fn1", To: "param1", Kind: EdgeDataFlow, Label: "parameter"})

	if len(cpg.Nodes()) != 2 { t.Errorf("expected 2 nodes, got %d", len(cpg.Nodes())) }
	if len(cpg.Edges()) != 1 { t.Errorf("expected 1 edge, got %d", len(cpg.Edges())) }

	n := cpg.GetNode("fn1")
	if n == nil || n.Name != "handleRequest" { t.Error("GetNode failed") }

	out := cpg.OutEdges("fn1")
	if len(out) != 1 || out[0].To != "param1" { t.Error("OutEdges failed") }

	in := cpg.InEdges("param1")
	if len(in) != 1 || in[0].From != "fn1" { t.Error("InEdges failed") }
}

func TestNodesByKind(t *testing.T) {
	cpg := NewCPG()
	cpg.AddNode(&Node{ID: "fn1", Kind: NodeFunction, Name: "a"})
	cpg.AddNode(&Node{ID: "fn2", Kind: NodeFunction, Name: "b"})
	cpg.AddNode(&Node{ID: "p1", Kind: NodeParameter, Name: "x"})

	fns := cpg.NodesByKind(NodeFunction)
	if len(fns) != 2 { t.Errorf("expected 2 functions, got %d", len(fns)) }
}

func TestSecurityAnnotations(t *testing.T) {
	cpg := NewCPG()
	fn := &Node{ID: "fn1", Kind: NodeFunction, Name: "handler", Annotations: make(map[string]bool)}
	fn.Annotations["has_auth"] = true
	fn.Annotations["handles_user_input"] = true
	cpg.AddNode(fn)

	n := cpg.GetNode("fn1")
	if !n.Annotations["has_auth"] { t.Error("expected has_auth annotation") }
}
