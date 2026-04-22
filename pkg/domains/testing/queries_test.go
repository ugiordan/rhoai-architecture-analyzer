package testing

import (
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/domains/security"
	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

func TestQueryUntestedSecurityFunc(t *testing.T) {
	g := graph.NewCPG()

	fn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "Handle",
		File:        "webhook.go",
		Line:        10,
		Annotations: map[string]bool{security.AnnotHandlesAdmission: true},
		Properties:  make(map[string]string),
	}
	g.AddNode(fn)

	findings := queryUntestedSecurityFunc(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].RuleID != "CGA-T01" {
		t.Errorf("expected CGA-T01, got %s", findings[0].RuleID)
	}
}
