package security

import (
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

func TestGoAnnotatorHandlesAdmission(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "Handle",
		File:        "webhook.go",
		Line:        10,
		Language:    "go",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
		ParamTypes:  []string{"context.Context", "admission.Request"},
	}
	g.AddNode(fn)

	a := &GoAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("fn1").Annotations[AnnotHandlesAdmission] {
		t.Error("expected sec:handles_admission annotation")
	}
}

func TestGoAnnotatorCreatesRBAC(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "bindRole",
		File: "rbac.go", Line: 10, EndLine: 20, Language: "go",
		Annotations: make(map[string]bool), Properties: make(map[string]string),
	}
	g.AddNode(fn)

	cs := &graph.Node{
		ID: "call1", Kind: graph.NodeCallSite, Name: "c.Create",
		File: "rbac.go", Line: 15, Language: "go",
		Annotations: make(map[string]bool),
		Properties:  map[string]string{"arg_types": "&rbacv1.ClusterRoleBinding"},
	}
	g.AddNode(cs)
	g.AddEdge(&graph.Edge{From: "fn1", To: "call1", Kind: graph.EdgeDataFlow, Label: "contains_call"})

	a := &GoAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("call1").Annotations[AnnotCreatesRBAC] {
		t.Error("expected sec:creates_rbac on call site")
	}
	if !g.GetNode("fn1").Annotations[AnnotCreatesRBAC] {
		t.Error("expected sec:creates_rbac on containing function")
	}
}

func TestGoAnnotatorGeneratesCert(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "generateCert",
		File: "cert.go", Line: 10, EndLine: 30, Language: "go",
		Annotations: make(map[string]bool), Properties: make(map[string]string),
	}
	g.AddNode(fn)

	sl := &graph.Node{
		ID: "struct1", Kind: graph.NodeStructLiteral, Name: "x509.Certificate",
		File:        "cert.go",
		Line:        15,
		Language:    "go",
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
		StructType:  "x509.Certificate",
		FieldNames:  []string{"SerialNumber", "Subject", "IsCA", "KeyUsage", "DNSNames"},
	}
	g.AddNode(sl)
	g.AddEdge(&graph.Edge{From: "fn1", To: "struct1", Kind: graph.EdgeDataFlow, Label: "contains_struct"})

	a := &GoAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("struct1").Annotations[AnnotGeneratesCert] {
		t.Error("expected sec:generates_cert on struct literal")
	}
	if !g.GetNode("fn1").Annotations[AnnotGeneratesCert] {
		t.Error("expected sec:generates_cert on containing function")
	}
}

func TestGoAnnotatorBindsSubject(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "bindRole",
		File: "rbac.go", Line: 10, EndLine: 30, Language: "go",
		Annotations: make(map[string]bool), Properties: make(map[string]string),
	}
	g.AddNode(fn)

	cs := &graph.Node{
		ID: "call1", Kind: graph.NodeCallSite, Name: "c.Create",
		File: "rbac.go", Line: 15, Language: "go",
		Annotations: make(map[string]bool),
		Properties:  map[string]string{"arg_types": "&rbacv1.ClusterRoleBinding"},
	}
	g.AddNode(cs)
	g.AddEdge(&graph.Edge{From: "fn1", To: "call1", Kind: graph.EdgeDataFlow, Label: "contains_call"})

	sl := &graph.Node{
		ID: "struct1", Kind: graph.NodeStructLiteral, Name: "rbacv1.Subject",
		File: "rbac.go", Line: 18, Language: "go",
		Annotations: make(map[string]bool),
		Properties: map[string]string{
			"type":          "rbacv1.Subject",
			"fields":        "Kind,Name",
			"string_values": "Group,system:authenticated",
		},
	}
	g.AddNode(sl)
	g.AddEdge(&graph.Edge{From: "fn1", To: "struct1", Kind: graph.EdgeDataFlow, Label: "contains_struct"})

	a := &GoAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("fn1").Annotations[AnnotBindsSubject] {
		t.Error("expected sec:binds_subject on function")
	}
}

func TestGoAnnotatorSetsTrustLevel(t *testing.T) {
	g := graph.NewCPG()

	// Untrusted: HTTP handler without auth
	httpHandler := &graph.Node{
		ID:          "http1",
		Kind:        graph.NodeHTTPEndpoint,
		Name:        "publicHandler",
		Route:       "/public",
		HTTPMethod:  "GET",
		Language:    "go",
		Annotations: make(map[string]bool),
	}
	g.AddNode(httpHandler)

	// Semi-trusted: admission webhook handler
	admissionFn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "HandleAdmission",
		Language:    "go",
		ParamTypes:  []string{"admission.Request"},
		Annotations: make(map[string]bool),
	}
	g.AddNode(admissionFn)

	// Trusted: reconciler
	reconciler := &graph.Node{
		ID:          "fn2",
		Kind:        graph.NodeFunction,
		Name:        "Reconcile",
		Language:    "go",
		ParamTypes:  []string{"context.Context", "ctrl.Request"},
		Annotations: make(map[string]bool),
	}
	g.AddNode(reconciler)

	a := &GoAnnotator{}
	a.Annotate(g, nil)

	if httpHandler.TrustLevel != graph.TrustUntrusted {
		t.Errorf("HTTP handler trust = %q, want %q", httpHandler.TrustLevel, graph.TrustUntrusted)
	}
	if admissionFn.TrustLevel != graph.TrustSemiTrusted {
		t.Errorf("admission handler trust = %q, want %q", admissionFn.TrustLevel, graph.TrustSemiTrusted)
	}
	if reconciler.TrustLevel != graph.TrustTrusted {
		t.Errorf("reconciler trust = %q, want %q", reconciler.TrustLevel, graph.TrustTrusted)
	}
}

func TestGoAnnotatorNoFalsePositives(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "doStuff",
		File: "regular.go", Line: 10, Language: "go",
		Annotations: make(map[string]bool),
		Properties:  map[string]string{"param_types": "context.Context,string"},
	}
	g.AddNode(fn)

	a := &GoAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	n := g.GetNode("fn1")
	for _, ann := range []string{
		AnnotCreatesRBAC, AnnotHandlesAdmission, AnnotGeneratesCert,
		AnnotAccessesSecret, AnnotBindsSubject, AnnotConfiguresCache,
		AnnotWritesPlaintextSecret,
	} {
		if n.Annotations[ann] {
			t.Errorf("unexpected annotation %q on regular function", ann)
		}
	}
}
