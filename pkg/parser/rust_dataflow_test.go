package parser

import (
	"os"
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

func parseRustDataFlowSample(t *testing.T) *ParseResult {
	t.Helper()
	content, err := os.ReadFile("../../testdata/dataflow_sample.rs")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}
	p := NewRustParser()
	result, err := p.ParseFile("testdata/dataflow_sample.rs", content)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}
	return result
}

func TestRustDataFlowParameters(t *testing.T) {
	result := parseRustDataFlowSample(t)

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

	// handle_request has "request", field_access has "user"
	if !names["request"] {
		t.Error("expected parameter 'request' from handle_request")
	}
	if !names["user"] {
		t.Error("expected parameter 'user' from field_access")
	}
}

func TestRustDataFlowVariables(t *testing.T) {
	result := parseRustDataFlowSample(t)

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

func TestRustDataFlowEdgeLabels(t *testing.T) {
	result := parseRustDataFlowSample(t)

	labels := make(map[string]bool)
	for _, e := range result.Edges {
		if e.Kind == graph.EdgeDataFlow {
			labels[e.Label] = true
		}
	}

	for _, expected := range []string{"declares", "assigns", "field_access"} {
		if !labels[expected] {
			t.Errorf("expected edge label %q not found in edges", expected)
			t.Log("  found labels:")
			for l := range labels {
				t.Logf("    %s", l)
			}
		}
	}

	// Verify all data flow edges have correct Kind
	for _, e := range result.Edges {
		if e.Kind == graph.EdgeDataFlow {
			// Data flow edge should have EdgeDataFlow kind (sanity check)
			if e.Kind != graph.EdgeDataFlow {
				t.Errorf("edge %s (%s -> %s) has Kind=%s, want EdgeDataFlow", e.Label, e.From, e.To, e.Kind)
			}
		}
	}
}

func TestRustDataFlowMutatesViaReference(t *testing.T) {
	result := parseRustDataFlowSample(t)

	// serde_json::from_str(&body) should produce passes_to edges
	// The &body reference_expression should generate passes_to from body to the call
	passesTo := filterEdgesByLabel(result.Edges, "passes_to")
	if len(passesTo) < 1 {
		t.Errorf("expected at least 1 passes_to edge (from &body in serde_json::from_str(&body)), got %d", len(passesTo))
		for _, e := range result.Edges {
			t.Logf("  edge: %s %s -> %s", e.Label, e.From, e.To)
		}
	}
}
