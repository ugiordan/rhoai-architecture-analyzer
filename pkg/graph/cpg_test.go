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

func TestTypedNodeFields(t *testing.T) {
	cpg := NewCPG()
	fn := &Node{
		ID:         "fn1",
		Kind:       NodeFunction,
		Name:       "handleRequest",
		File:       "server.go",
		Line:       10,
		Language:   "go",
		ParamTypes: []string{"string", "int"},
		ParamNames: []string{"name", "age"},
		ReturnType: "error",
		IsTest:     false,
		IsUnsafe:   false,
		Route:      "",
	}
	cpg.AddNode(fn)

	n := cpg.GetNode("fn1")
	if len(n.ParamTypes) != 2 || n.ParamTypes[0] != "string" {
		t.Errorf("expected ParamTypes [string int], got %v", n.ParamTypes)
	}
	if n.ReturnType != "error" {
		t.Errorf("expected ReturnType 'error', got %q", n.ReturnType)
	}
}

func TestTypedNodeHTTPEndpoint(t *testing.T) {
	cpg := NewCPG()
	ep := &Node{
		ID:         "ep1",
		Kind:       NodeHTTPEndpoint,
		Name:       "GET /users",
		File:       "routes.go",
		Line:       5,
		Route:      "/users",
		HTTPMethod: "GET",
	}
	cpg.AddNode(ep)

	n := cpg.GetNode("ep1")
	if n.Route != "/users" {
		t.Errorf("expected Route '/users', got %q", n.Route)
	}
	if n.HTTPMethod != "GET" {
		t.Errorf("expected HTTPMethod 'GET', got %q", n.HTTPMethod)
	}
}

func TestTypedNodeDBOperation(t *testing.T) {
	cpg := NewCPG()
	op := &Node{
		ID:        "db1",
		Kind:      NodeDBOperation,
		Name:      "db.Query",
		File:      "repo.go",
		Line:      20,
		Operation: "read",
		Table:     "users",
	}
	cpg.AddNode(op)

	n := cpg.GetNode("db1")
	if n.Operation != "read" {
		t.Errorf("expected Operation 'read', got %q", n.Operation)
	}
	if n.Table != "users" {
		t.Errorf("expected Table 'users', got %q", n.Table)
	}
}

func TestTypedNodeStructLiteral(t *testing.T) {
	cpg := NewCPG()
	sl := &Node{
		ID:         "sl1",
		Kind:       NodeStructLiteral,
		Name:       "Config",
		File:       "config.go",
		Line:       15,
		StructType: "Config",
		FieldNames: []string{"Host", "Port"},
	}
	cpg.AddNode(sl)

	n := cpg.GetNode("sl1")
	if n.StructType != "Config" {
		t.Errorf("expected StructType 'Config', got %q", n.StructType)
	}
	if len(n.FieldNames) != 2 {
		t.Errorf("expected 2 FieldNames, got %d", len(n.FieldNames))
	}
}

func TestEdgeConfidence(t *testing.T) {
	cpg := NewCPG()
	cpg.AddNode(&Node{ID: "cs1", Kind: NodeCallSite, Name: "doStuff"})
	cpg.AddNode(&Node{ID: "fn1", Kind: NodeFunction, Name: "doStuff"})

	cpg.AddEdge(&Edge{
		From:       "cs1",
		To:         "fn1",
		Kind:       EdgeCalls,
		Label:      "doStuff",
		Confidence: ConfidenceCertain,
	})

	edges := cpg.OutEdges("cs1")
	if len(edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(edges))
	}
	if edges[0].Confidence != ConfidenceCertain {
		t.Errorf("expected confidence CERTAIN, got %q", edges[0].Confidence)
	}
}

func TestEdgeConfidenceOmitsEmpty(t *testing.T) {
	e := &Edge{From: "a", To: "b", Kind: EdgeDataFlow, Label: "contains_call"}
	if e.Confidence != "" {
		t.Errorf("expected empty confidence for non-CALLS edge, got %q", e.Confidence)
	}
}
