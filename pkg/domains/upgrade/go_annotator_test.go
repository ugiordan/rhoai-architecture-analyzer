package upgrade

import (
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

func TestGoAnnotatorVersionConversion(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "ConvertTo",
		File:        "api/v1beta1/conversion.go",
		Line:        10,
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(fn)

	a := &GoAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("fn1").Annotations[AnnotVersionConversion] {
		t.Error("expected upgrade:version_conversion annotation")
	}
}

func TestGoAnnotatorFeatureGate(t *testing.T) {
	g := graph.NewCPG()
	cs := &graph.Node{
		ID:          "call1",
		Kind:        graph.NodeCallSite,
		Name:        "utilfeature.DefaultFeatureGate.Enabled",
		File:        "controller.go",
		Line:        25,
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(cs)

	a := &GoAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	if !g.GetNode("call1").Annotations[AnnotFeatureGate] {
		t.Error("expected upgrade:feature_gate annotation")
	}
}

func TestGoAnnotatorNoFalsePositives(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "Reconcile",
		File:        "controller.go",
		Line:        10,
		Annotations: make(map[string]bool),
		Properties:  make(map[string]string),
	}
	g.AddNode(fn)

	a := &GoAnnotator{}
	if err := a.Annotate(g, nil); err != nil {
		t.Fatalf("Annotate failed: %v", err)
	}

	n := g.GetNode("fn1")
	for _, ann := range []string{AnnotVersionConversion, AnnotFeatureGate, AnnotPreReleaseAPI, AnnotMigration, AnnotVersionCheck} {
		if n.Annotations[ann] {
			t.Errorf("unexpected annotation %q on regular function", ann)
		}
	}
}
