package parser

import (
	"os"
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

func parseTSDataFlowSample(t *testing.T) *ParseResult {
	t.Helper()
	content, err := os.ReadFile("../../testdata/dataflow_sample.ts")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}
	p := NewTypeScriptParser()
	result, err := p.ParseFile("testdata/dataflow_sample.ts", content)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}
	return result
}

func TestTypeScriptDataFlowParameters(t *testing.T) {
	result := parseTSDataFlowSample(t)

	if len(result.Parameters) == 0 {
		t.Fatal("expected Parameters, got 0")
	}

	names := make(map[string]bool)
	for _, p := range result.Parameters {
		names[p.Name] = true
		if p.Kind != graph.NodeParameter {
			t.Errorf("parameter %q has Kind=%s, want NodeParameter", p.Name, p.Kind)
		}
	}

	// handleRequest has "req" and "res"
	if !names["req"] {
		t.Error("expected parameter 'req' from handleRequest")
	}
	if !names["res"] {
		t.Error("expected parameter 'res' from handleRequest")
	}
}

func TestTypeScriptDataFlowVariables(t *testing.T) {
	result := parseTSDataFlowSample(t)

	if len(result.Variables) < 4 {
		t.Errorf("expected at least 4 variables, got %d", len(result.Variables))
		for _, v := range result.Variables {
			t.Logf("  variable: %s", v.Name)
		}
	}

	names := make(map[string]bool)
	for _, v := range result.Variables {
		names[v.Name] = true
		if v.Kind != graph.NodeVariable {
			t.Errorf("variable %q has Kind=%s, want NodeVariable", v.Name, v.Kind)
		}
	}

	for _, expected := range []string{"body", "payload", "name", "query"} {
		if !names[expected] {
			t.Errorf("expected variable %q not found", expected)
		}
	}
}

func TestTypeScriptDataFlowEdgeLabels(t *testing.T) {
	result := parseTSDataFlowSample(t)

	labels := make(map[string]bool)
	for _, e := range result.Edges {
		labels[e.Label] = true
	}

	for _, expected := range []string{"declares", "assigns", "passes_to", "field_access", "returns"} {
		if !labels[expected] {
			t.Errorf("expected edge label %q not found in edges", expected)
			t.Log("  found labels:")
			for l := range labels {
				t.Logf("    %s", l)
			}
		}
	}

	// Verify all data flow edges have correct Kind (filter out control flow edges)
	for _, e := range result.Edges {
		// Skip control flow edges (they're tested separately in CFG tests)
		if e.Kind == graph.EdgeControlFlow {
			continue
		}
		if e.Kind != graph.EdgeDataFlow {
			t.Errorf("edge %s (%s -> %s) has Kind=%s, want EdgeDataFlow", e.Label, e.From, e.To, e.Kind)
		}
	}
}
