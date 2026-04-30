package domains

import (
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
	"github.com/ugiordan/architecture-analyzer/pkg/query"
)

// testDomain is a minimal DomainAnalyzer for testing.
type testDomain struct {
	name      string
	deps      []string
	annotated bool
}

func (d *testDomain) Name() string                { return d.name }
func (d *testDomain) SupportedLanguages() []string { return []string{"go"} }
func (d *testDomain) Dependencies() []string       { return d.deps }
func (d *testDomain) Queries() []query.Rule        { return nil }
func (d *testDomain) Annotate(g *graph.CPG, lang string, archData *ArchitectureData) error {
	d.annotated = true
	return nil
}

func TestOrchestratorSort(t *testing.T) {
	sec := &testDomain{name: "security", deps: nil}
	test := &testDomain{name: "testing", deps: []string{"security"}}
	upg := &testDomain{name: "upgrade", deps: nil}

	sorted, err := sortByDependency([]DomainAnalyzer{test, upg, sec})
	if err != nil {
		t.Fatalf("sortByDependency failed: %v", err)
	}

	// security must come before testing
	secIdx, testIdx := -1, -1
	for i, d := range sorted {
		switch d.Name() {
		case "security":
			secIdx = i
		case "testing":
			testIdx = i
		}
	}
	if secIdx >= testIdx {
		t.Errorf("security (idx=%d) should come before testing (idx=%d)", secIdx, testIdx)
	}
}

func TestOrchestratorCyclicDeps(t *testing.T) {
	a := &testDomain{name: "a", deps: []string{"b"}}
	b := &testDomain{name: "b", deps: []string{"a"}}

	_, err := sortByDependency([]DomainAnalyzer{a, b})
	if err == nil {
		t.Error("expected error for cyclic dependencies")
	}
}

func TestOrchestratorMissingDependency(t *testing.T) {
	test := &testDomain{name: "testing", deps: []string{"security"}}

	_, err := sortByDependency([]DomainAnalyzer{test})
	if err == nil {
		t.Error("expected error for missing dependency")
	}
}

func TestResolveDependencies(t *testing.T) {
	// Save and restore registry state
	origRegistry := registry
	registry = map[string]DomainAnalyzer{}
	defer func() { registry = origRegistry }()

	sec := &testDomain{name: "security", deps: nil}
	test := &testDomain{name: "testing", deps: []string{"security"}}
	upg := &testDomain{name: "upgrade", deps: nil}
	Register(sec)
	Register(test)
	Register(upg)

	// Requesting "testing" should auto-include "security"
	resolved, err := ResolveDependencies([]string{"testing"})
	if err != nil {
		t.Fatalf("ResolveDependencies failed: %v", err)
	}
	if len(resolved) != 2 {
		t.Fatalf("expected 2 resolved domains, got %d: %v", len(resolved), resolved)
	}

	// Requesting unknown domain should error
	_, err = ResolveDependencies([]string{"nonexistent"})
	if err == nil {
		t.Error("expected error for unknown domain")
	}
}

func TestOrchestratorRunAll(t *testing.T) {
	sec := &testDomain{name: "security"}
	test := &testDomain{name: "testing", deps: []string{"security"}}

	o := NewOrchestrator([]DomainAnalyzer{test, sec})
	cpg := graph.NewCPG()

	results, err := o.Run(cpg, "go", nil)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if !sec.annotated {
		t.Error("security domain was not annotated")
	}
	if !test.annotated {
		t.Error("testing domain was not annotated")
	}

	if len(results) != 2 {
		t.Errorf("expected 2 domain results, got %d", len(results))
	}
}

// mockDomain is a DomainAnalyzer that sets an annotation and returns a finding.
type mockDomain struct {
	name      string
	annotated bool
}

func (m *mockDomain) Name() string                { return m.name }
func (m *mockDomain) SupportedLanguages() []string { return []string{"go"} }
func (m *mockDomain) Dependencies() []string       { return nil }
func (m *mockDomain) Annotate(g *graph.CPG, lang string, archData *ArchitectureData) error {
	m.annotated = true
	g.SetAnnotation("fn1", "test_annotation", true)
	return nil
}
func (m *mockDomain) Queries() []query.Rule {
	return []query.Rule{{
		ID: "TEST-001", Name: "test-rule", Domain: m.name, Severity: "low",
		Run: func(cpg *graph.CPG) []query.Finding {
			return []query.Finding{{RuleID: "TEST-001", Message: "test"}}
		},
	}}
}

func TestOrchestratorAnnotateAll(t *testing.T) {
	cpg := graph.NewCPG()
	_ = cpg.AddNode(&graph.Node{ID: "fn1", Kind: graph.NodeFunction, Name: "myFunc"})

	md := &mockDomain{name: "mock"}
	o := NewOrchestrator([]DomainAnalyzer{md})

	if err := o.AnnotateAll(cpg, "go", nil); err != nil {
		t.Fatalf("AnnotateAll failed: %v", err)
	}

	if !md.annotated {
		t.Error("mock domain was not annotated")
	}

	node := cpg.GetNode("fn1")
	if node == nil {
		t.Fatal("node fn1 not found")
	}
	if !node.Annotations["test_annotation"] {
		t.Error("expected test_annotation to be set on fn1")
	}
}

func TestOrchestratorRunQueries(t *testing.T) {
	cpg := graph.NewCPG()
	_ = cpg.AddNode(&graph.Node{ID: "fn1", Kind: graph.NodeFunction, Name: "myFunc"})

	md := &mockDomain{name: "mock"}
	o := NewOrchestrator([]DomainAnalyzer{md})

	if err := o.AnnotateAll(cpg, "go", nil); err != nil {
		t.Fatalf("AnnotateAll failed: %v", err)
	}

	results, err := o.RunQueries(cpg)
	if err != nil {
		t.Fatalf("RunQueries failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	r := results[0]
	if r.Domain != "mock" {
		t.Errorf("expected domain 'mock', got %q", r.Domain)
	}
	if r.AnnotationsAdded != 1 {
		t.Errorf("expected 1 annotation added, got %d", r.AnnotationsAdded)
	}
	if len(r.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(r.Findings))
	}
	if r.Findings[0].Domain != "mock" {
		t.Errorf("expected finding domain 'mock', got %q", r.Findings[0].Domain)
	}
	if r.Findings[0].RuleID != "TEST-001" {
		t.Errorf("expected finding rule ID 'TEST-001', got %q", r.Findings[0].RuleID)
	}
}

func TestOrchestratorRunQueriesBeforeAnnotate(t *testing.T) {
	cpg := graph.NewCPG()

	md := &mockDomain{name: "mock"}
	o := NewOrchestrator([]DomainAnalyzer{md})

	_, err := o.RunQueries(cpg)
	if err == nil {
		t.Error("expected error when calling RunQueries before AnnotateAll")
	}
}

func TestOrchestratorRunBackwardsCompat(t *testing.T) {
	cpg := graph.NewCPG()
	_ = cpg.AddNode(&graph.Node{ID: "fn1", Kind: graph.NodeFunction, Name: "myFunc"})

	md := &mockDomain{name: "mock"}
	o := NewOrchestrator([]DomainAnalyzer{md})

	results, err := o.Run(cpg, "go", nil)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if !md.annotated {
		t.Error("mock domain was not annotated")
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	r := results[0]
	if r.Domain != "mock" {
		t.Errorf("expected domain 'mock', got %q", r.Domain)
	}
	if len(r.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(r.Findings))
	}
	if r.Findings[0].Domain != "mock" {
		t.Errorf("expected finding domain 'mock', got %q", r.Findings[0].Domain)
	}

	node := cpg.GetNode("fn1")
	if !node.Annotations["test_annotation"] {
		t.Error("expected test_annotation to be set on fn1")
	}
}
