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
