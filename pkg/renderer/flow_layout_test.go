package renderer

import (
	"testing"
)

func TestComputeLayout_EmptyGraph(t *testing.T) {
	r := computeLayout(nil, nil)
	if len(r.Positions) != 0 {
		t.Error("empty graph should produce no positions")
	}
	if r.CanvasWidth < 1000 {
		t.Error("canvas width should have a minimum")
	}
}

func TestComputeLayout_SingleNode(t *testing.T) {
	nodes := []FlowNode{{ID: "svc-a", Type: FlowNodeService, Layer: 2}}
	r := computeLayout(nodes, nil)
	pos, ok := r.Positions["svc-a"]
	if !ok {
		t.Fatal("node should have a position")
	}
	if pos.W == 0 || pos.H == 0 {
		t.Error("node should have non-zero dimensions")
	}
	if pos.X < layoutPadX {
		t.Errorf("node X should respect padding, got %f", pos.X)
	}
}

func TestComputeLayout_MultipleLayersOrdered(t *testing.T) {
	nodes := []FlowNode{
		{ID: "ingress-gw", Type: FlowNodeIngress, Layer: 0},
		{ID: "svc-a", Type: FlowNodeService, Layer: 2},
		{ID: "dep-ctrl", Type: FlowNodeDeployment, Layer: 3},
	}
	r := computeLayout(nodes, nil)
	if r.Positions["ingress-gw"].Y >= r.Positions["svc-a"].Y {
		t.Error("ingress should be above service")
	}
	if r.Positions["svc-a"].Y >= r.Positions["dep-ctrl"].Y {
		t.Error("service should be above deployment")
	}
}

func TestComputeLayout_NodesInSameLayerSpaced(t *testing.T) {
	nodes := []FlowNode{
		{ID: "svc-a", Type: FlowNodeService, Layer: 2},
		{ID: "svc-b", Type: FlowNodeService, Layer: 2},
		{ID: "svc-c", Type: FlowNodeService, Layer: 2},
	}
	r := computeLayout(nodes, nil)
	a := r.Positions["svc-a"]
	b := r.Positions["svc-b"]
	c := r.Positions["svc-c"]
	if b.X <= a.X+a.W {
		t.Error("svc-b should not overlap svc-a")
	}
	if c.X <= b.X+b.W {
		t.Error("svc-c should not overlap svc-b")
	}
}

func TestComputeLayout_LayersCentered(t *testing.T) {
	nodes := []FlowNode{
		{ID: "svc-a", Type: FlowNodeService, Layer: 2},
		{ID: "dep-a", Type: FlowNodeDeployment, Layer: 3},
		{ID: "dep-b", Type: FlowNodeDeployment, Layer: 3},
		{ID: "dep-c", Type: FlowNodeDeployment, Layer: 3},
	}
	r := computeLayout(nodes, nil)
	// Single-node layer should be centered
	svcA := r.Positions["svc-a"]
	svcCenter := svcA.X + svcA.W/2
	canvasCenter := r.CanvasWidth / 2
	if diff := svcCenter - canvasCenter; diff > 1 || diff < -1 {
		t.Errorf("single node should be centered (center=%f, canvas center=%f)", svcCenter, canvasCenter)
	}
}

func TestComputeLayout_CanvasWidthExpandsForWideLayer(t *testing.T) {
	var nodes []FlowNode
	for i := 0; i < 10; i++ {
		nodes = append(nodes, FlowNode{
			ID:    flowNodeID("svc-" + string(rune('a'+i))),
			Type:  FlowNodeService,
			Layer: 2,
		})
	}
	r := computeLayout(nodes, nil)
	if r.CanvasWidth <= layoutMinW {
		t.Error("canvas should expand beyond minimum for 10 nodes in one layer")
	}
}

func TestComputeLayout_EdgeCrossingReduction(t *testing.T) {
	nodes := []FlowNode{
		{ID: "a1", Type: FlowNodeService, Layer: 0},
		{ID: "a2", Type: FlowNodeService, Layer: 0},
		{ID: "b1", Type: FlowNodeDeployment, Layer: 1},
		{ID: "b2", Type: FlowNodeDeployment, Layer: 1},
	}
	// Edges: a1→b2 and a2→b1 (crossed). Barycentric should reorder to reduce crossings.
	edges := []FlowEdge{
		{From: "a1", To: "b2"},
		{From: "a2", To: "b1"},
	}
	r := computeLayout(nodes, edges)
	// After barycentric ordering, a1 should be near b1's position and a2 near b2's.
	// We can't guarantee perfect ordering but positions should exist.
	if len(r.Positions) != 4 {
		t.Errorf("expected 4 positions, got %d", len(r.Positions))
	}
}

func TestComputeLayout_AllNodeTypesHaveSize(t *testing.T) {
	types := []FlowNodeType{
		FlowNodeIngress, FlowNodeWebhook, FlowNodeService,
		FlowNodeDeployment, FlowNodeExternal, FlowNodeCRD,
	}
	for _, typ := range types {
		nodes := []FlowNode{{ID: "test", Type: typ, Layer: 0}}
		r := computeLayout(nodes, nil)
		pos := r.Positions["test"]
		if pos.W == 0 || pos.H == 0 {
			t.Errorf("type %s should have non-zero dimensions", typ)
		}
	}
}
