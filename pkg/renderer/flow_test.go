package renderer

import (
	"strings"
	"testing"
)

// ---- flowNodeID tests ----

func TestFlowNodeID_BareAlphanumeric(t *testing.T) {
	if id := flowNodeID("my-service"); id != "my-service" {
		t.Errorf("expected 'my-service', got %q", id)
	}
}

func TestFlowNodeID_DotsBecomeDashes(t *testing.T) {
	if id := flowNodeID("serving.kserve.io"); id != "serving-kserve-io" {
		t.Errorf("expected 'serving-kserve-io', got %q", id)
	}
}

func TestFlowNodeID_SlashesBecomeDashes(t *testing.T) {
	if id := flowNodeID("ns/my-svc"); id != "ns-my-svc" {
		t.Errorf("expected 'ns-my-svc', got %q", id)
	}
}

func TestFlowNodeID_ColonsBecomeDashes(t *testing.T) {
	if id := flowNodeID("kube:system"); id != "kube-system" {
		t.Errorf("expected 'kube-system', got %q", id)
	}
}

func TestFlowNodeID_SpacesBecomeDashes(t *testing.T) {
	if id := flowNodeID("my service"); id != "my-service" {
		t.Errorf("expected 'my-service', got %q", id)
	}
}

func TestFlowNodeID_EmptyString(t *testing.T) {
	if id := flowNodeID(""); id != "node" {
		t.Errorf("expected 'node', got %q", id)
	}
}

func TestFlowNodeID_Underscores(t *testing.T) {
	if id := flowNodeID("my_service"); id != "my_service" {
		t.Errorf("expected 'my_service', got %q", id)
	}
}

// ---- stripNamespace tests ----

func TestStripNamespace_WithNamespace(t *testing.T) {
	if r := stripNamespace("opendatahub/my-svc"); r != "my-svc" {
		t.Errorf("expected 'my-svc', got %q", r)
	}
}

func TestStripNamespace_NoNamespace(t *testing.T) {
	if r := stripNamespace("my-svc"); r != "my-svc" {
		t.Errorf("expected 'my-svc', got %q", r)
	}
}

func TestStripNamespace_Empty(t *testing.T) {
	if r := stripNamespace(""); r != "" {
		t.Errorf("expected empty, got %q", r)
	}
}

func TestStripNamespace_JustSlash(t *testing.T) {
	if r := stripNamespace("/"); r != "" {
		t.Errorf("expected empty, got %q", r)
	}
}

func TestStripNamespace_MultiSlash(t *testing.T) {
	if r := stripNamespace("a/b/c"); r != "c" {
		t.Errorf("expected 'c', got %q", r)
	}
}

// ---- controllerDepID tests ----

func TestControllerDepID_MatchesComponentName(t *testing.T) {
	nodes := []FlowNode{
		{ID: "dep-webhook", Label: "webhook", Type: FlowNodeDeployment},
		{ID: "dep-kserve-ctrl", Label: "kserve-ctrl", Type: FlowNodeDeployment},
	}
	id := controllerDepID(nodes, "kserve")
	if id != "dep-kserve-ctrl" {
		t.Errorf("expected dep-kserve-ctrl, got %q", id)
	}
}

func TestControllerDepID_FallsBackToFirstDep(t *testing.T) {
	nodes := []FlowNode{
		{ID: "dep-alpha", Label: "alpha", Type: FlowNodeDeployment},
		{ID: "dep-beta", Label: "beta", Type: FlowNodeDeployment},
	}
	id := controllerDepID(nodes, "unrelated-component")
	if id != "dep-alpha" {
		t.Errorf("expected dep-alpha (first deployment), got %q", id)
	}
}

func TestControllerDepID_NoDeployments(t *testing.T) {
	nodes := []FlowNode{
		{ID: "svc-x", Label: "x", Type: FlowNodeService},
	}
	id := controllerDepID(nodes, "comp")
	if id != "" {
		t.Errorf("expected empty string, got %q", id)
	}
}

// ---- FlowGraph builder: network flow focus ----

func TestBuildFlowGraph_HasClientNode(t *testing.T) {
	g := BuildFlowGraph(emptyComponentData())
	var found bool
	for _, n := range g.Nodes {
		if n.ID == "client" && n.Label == "Client" {
			found = true
		}
	}
	if !found {
		t.Error("graph should always have a synthetic Client node")
	}
}

func TestBuildFlowGraph_ServicesBecomesNodes(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"services": []interface{}{
			map[string]interface{}{"name": "my-svc", "type": "ClusterIP"},
		},
	}
	g := BuildFlowGraph(data)
	var found bool
	for _, n := range g.Nodes {
		if n.ID == "svc-my-svc" && n.Type == FlowNodeService {
			found = true
		}
	}
	if !found {
		t.Error("service should produce a FlowNodeService node")
	}
}

func TestBuildFlowGraph_WebhooksCollapsed(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"webhooks": makeWebhooks(5, "validating"),
	}
	g := BuildFlowGraph(data)
	var whCount int
	for _, n := range g.Nodes {
		if n.Type == FlowNodeWebhook {
			whCount++
		}
	}
	if whCount > 2 {
		t.Errorf("5+ webhooks of same type should be collapsed, got %d webhook nodes", whCount)
	}
}

func TestBuildFlowGraph_ExternalCollapsedByType(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"external_connections": []interface{}{
			map[string]interface{}{"target": "", "type": "object-storage"},
			map[string]interface{}{"target": "", "type": "object-storage"},
			map[string]interface{}{"target": "", "type": "object-storage"},
			map[string]interface{}{"target": "redis", "type": "cache"},
		},
	}
	g := BuildFlowGraph(data)
	var extLabels []string
	for _, n := range g.Nodes {
		if n.Type == FlowNodeExternal {
			extLabels = append(extLabels, n.Label)
		}
	}
	if len(extLabels) != 2 {
		t.Errorf("expected 2 external nodes (object-storage collapsed, cache separate), got %v", extLabels)
	}
}

func TestBuildFlowGraph_NoCRDNodes(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"crds": []interface{}{
			map[string]interface{}{"kind": "MyKind", "group": "g"},
		},
	}
	g := BuildFlowGraph(data)
	for _, n := range g.Nodes {
		if n.Type == FlowNodeCRD {
			t.Error("network flow graph should not contain CRD nodes")
		}
	}
}

func TestBuildFlowGraph_DeploymentsBecomesNodes(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"deployments": []interface{}{
			map[string]interface{}{"name": "my-controller"},
		},
	}
	g := BuildFlowGraph(data)
	var found bool
	for _, n := range g.Nodes {
		if n.Type == FlowNodeDeployment && n.ID == "dep-my-controller" {
			found = true
		}
	}
	if !found {
		t.Error("deployment should produce a FlowNodeDeployment node")
	}
}

func TestBuildFlowGraph_ClientToIngressEdge(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"ingress_routing": []interface{}{
			map[string]interface{}{"name": "gw"},
		},
	}
	g := BuildFlowGraph(data)
	var found bool
	for _, e := range g.Edges {
		if e.From == "client" && e.Type == "route" {
			found = true
		}
	}
	if !found {
		t.Error("should have an edge from client to ingress")
	}
}

func TestBuildFlowGraph_RequestFlowPathExists(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"ingress_routing": []interface{}{
			map[string]interface{}{"name": "gw"},
		},
		"services": []interface{}{
			map[string]interface{}{"name": "svc"},
		},
		"deployments": []interface{}{
			map[string]interface{}{"name": "app"},
		},
	}
	g := BuildFlowGraph(data)
	var found bool
	for _, p := range g.Paths {
		if p.Name == "Request Flow" && len(p.Edges) >= 2 {
			found = true
		}
	}
	if !found {
		t.Fatalf("should have a Request Flow path with 2+ edges, got: %v", pathNames(g.Paths))
	}
}

func TestBuildFlowGraph_ExternalEdgesFromDeployment(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"deployments": []interface{}{
			map[string]interface{}{"name": "my-app"},
		},
		"external_connections": []interface{}{
			map[string]interface{}{"target": "redis", "type": "cache"},
		},
	}
	g := BuildFlowGraph(data)
	var found bool
	for _, e := range g.Edges {
		if e.Type == "external" {
			found = true
		}
	}
	if !found {
		t.Error("external connection should produce an 'external' edge")
	}
}

func TestBuildFlowGraph_NilSliceSafety(t *testing.T) {
	g := BuildFlowGraph(emptyComponentData())
	if g.Nodes == nil {
		t.Error("Nodes should be empty slice, not nil")
	}
	if g.Edges == nil {
		t.Error("Edges should be empty slice, not nil")
	}
	if g.Paths == nil {
		t.Error("Paths should be empty slice, not nil")
	}
}

func TestBuildFlowGraph_WebhookServiceRefWithNamespace(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"services": []interface{}{
			map[string]interface{}{"name": "webhook-svc"},
		},
		"webhooks": []interface{}{
			map[string]interface{}{
				"name": "validate", "type": "validating",
				"service_ref": "opendatahub/webhook-svc",
			},
		},
	}
	g := BuildFlowGraph(data)
	var found bool
	for _, e := range g.Edges {
		if e.Type == "intercept" {
			found = true
		}
	}
	if !found {
		t.Error("webhook with namespace-qualified service_ref should produce intercept edge")
	}
}

func TestBuildFlowGraph_IngressServiceRefWithNamespace(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"services": []interface{}{
			map[string]interface{}{"name": "my-svc"},
		},
		"ingress_routing": []interface{}{
			map[string]interface{}{"name": "route", "service_ref": "default/my-svc"},
		},
	}
	g := BuildFlowGraph(data)
	var found bool
	for _, e := range g.Edges {
		if e.Type == "route" && e.To == "svc-my-svc" {
			found = true
		}
	}
	if !found {
		t.Error("ingress with namespace-qualified service_ref should produce route edge")
	}
}

// FlowRenderer has been removed. Flow visualization is now handled by
// a separate frontend project (viz/) using proper JS tooling, not Go templates.
// The graph builder (buildFlowGraph, buildFlowPaths) is kept as exported
// functions for the /flow-diagram skill to use.

// ---- helpers ----

func makeWebhooks(n int, whType string) []interface{} {
	var whs []interface{}
	for i := 0; i < n; i++ {
		whs = append(whs, map[string]interface{}{
			"name": strings.Repeat("a", i+1), "type": whType,
		})
	}
	return whs
}

func pathNames(paths []FlowPath) []string {
	var names []string
	for _, p := range paths {
		names = append(names, p.Name)
	}
	return names
}
