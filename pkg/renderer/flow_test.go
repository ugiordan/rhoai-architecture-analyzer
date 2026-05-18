package renderer

import (
	"strings"
	"testing"
)

// ---- FlowGraph builder tests ----

func TestBuildFlowGraph_EmptyData(t *testing.T) {
	g := buildFlowGraph(emptyComponentData())
	if g.Component == "" {
		t.Error("component name should be set")
	}
}

func TestBuildFlowGraph_ServicesBecomesNodes(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"services": []interface{}{
			map[string]interface{}{"name": "my-svc", "type": "ClusterIP"},
		},
	}
	g := buildFlowGraph(data)
	var found bool
	for _, n := range g.Nodes {
		if n.ID == "svc-my-svc" && n.Type == FlowNodeService {
			found = true
		}
	}
	if !found {
		t.Error("service should produce a FlowNodeService node with id 'svc-my-svc'")
	}
}

func TestBuildFlowGraph_WebhooksBecomesNodes(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"webhooks": []interface{}{
			map[string]interface{}{"name": "validate-nb", "type": "validating"},
		},
	}
	g := buildFlowGraph(data)
	var found bool
	for _, n := range g.Nodes {
		if n.Type == FlowNodeWebhook {
			found = true
		}
	}
	if !found {
		t.Error("webhook should produce a FlowNodeWebhook node")
	}
}

func TestBuildFlowGraph_ExternalConnectionsBecomesNodes(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"external_connections": []interface{}{
			map[string]interface{}{"target": "postgres", "type": "database"},
		},
	}
	g := buildFlowGraph(data)
	var found bool
	for _, n := range g.Nodes {
		if n.Type == FlowNodeExternal {
			found = true
		}
	}
	if !found {
		t.Error("external connection should produce FlowNodeExternal node")
	}
}

func TestBuildFlowGraph_CRDsBecomesNodes(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"crds": []interface{}{
			map[string]interface{}{"kind": "MyKind", "group": "mygroup", "version": "v1"},
		},
	}
	g := buildFlowGraph(data)
	var found bool
	for _, n := range g.Nodes {
		if n.Type == FlowNodeCRD && n.ID == "crd-MyKind" {
			found = true
		}
	}
	if !found {
		t.Error("CRD should produce a FlowNodeCRD node with id 'crd-MyKind'")
	}
}

func TestBuildFlowGraph_DeploymentsBecomesNodes(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"deployments": []interface{}{
			map[string]interface{}{"name": "my-controller"},
		},
	}
	g := buildFlowGraph(data)
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

func TestBuildFlowGraph_IngressBecomesNodes(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"ingress_routing": []interface{}{
			map[string]interface{}{"name": "my-gateway", "kind": "HTTPRoute"},
		},
	}
	g := buildFlowGraph(data)
	var found bool
	for _, n := range g.Nodes {
		if n.Type == FlowNodeIngress {
			found = true
		}
	}
	if !found {
		t.Error("ingress routing should produce a FlowNodeIngress node")
	}
}

func TestBuildFlowGraph_ControllerWatchesBecomesEdges(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"deployments": []interface{}{
			map[string]interface{}{"name": "my-controller"},
		},
		"crds": []interface{}{
			map[string]interface{}{"kind": "MyKind", "group": "mygroup", "version": "v1"},
		},
		"controller_watches": []interface{}{
			map[string]interface{}{"type": "For", "gvk": "mygroup/v1/MyKind"},
		},
	}
	g := buildFlowGraph(data)
	var found bool
	for _, e := range g.Edges {
		if e.Type == "watches" {
			found = true
		}
	}
	if !found {
		t.Error("controller watch should produce a 'watches' edge")
	}
}

func TestBuildFlowGraph_ExternalEdgesFromDeployment(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"deployments": []interface{}{
			map[string]interface{}{"name": "my-app"},
		},
		"external_connections": []interface{}{
			map[string]interface{}{"target": "redis", "type": "database"},
		},
	}
	g := buildFlowGraph(data)
	var found bool
	for _, e := range g.Edges {
		if e.Type == "external" && e.To == "ext-redis" {
			found = true
		}
	}
	if !found {
		t.Error("external connection should produce an 'external' edge from deployment")
	}
}

func TestBuildFlowGraph_EdgeIDsAssigned(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"services": []interface{}{
			map[string]interface{}{"name": "svc", "type": "ClusterIP"},
		},
		"deployments": []interface{}{
			map[string]interface{}{"name": "app"},
		},
	}
	g := buildFlowGraph(data)
	for _, e := range g.Edges {
		if e.ID == "" {
			t.Errorf("edge from %s to %s has no ID", e.From, e.To)
		}
	}
}

func TestBuildFlowGraph_NodeLayerAssignments(t *testing.T) {
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
		"external_connections": []interface{}{
			map[string]interface{}{"target": "db", "type": "database"},
		},
	}
	g := buildFlowGraph(data)
	layers := map[FlowNodeType]int{}
	for _, n := range g.Nodes {
		layers[n.Type] = n.Layer
	}
	if layers[FlowNodeIngress] >= layers[FlowNodeService] {
		t.Error("ingress should be at a lower layer number than service")
	}
	if layers[FlowNodeService] >= layers[FlowNodeDeployment] {
		t.Error("service should be at a lower layer number than deployment")
	}
	if layers[FlowNodeDeployment] >= layers[FlowNodeExternal] {
		t.Error("deployment should be at a lower layer number than external")
	}
}

func TestBuildFlowGraph_PathsBuiltWhenDataPresent(t *testing.T) {
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
	g := buildFlowGraph(data)
	if len(g.Paths) == 0 {
		t.Error("should build at least one flow path when ingress+service+deployment are present")
	}
}

func TestBuildFlowGraph_NodesDeduped(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"services": []interface{}{
			map[string]interface{}{"name": "dup-svc"},
			map[string]interface{}{"name": "dup-svc"},
		},
	}
	g := buildFlowGraph(data)
	count := 0
	for _, n := range g.Nodes {
		if n.ID == "svc-dup-svc" {
			count++
		}
	}
	if count > 1 {
		t.Errorf("duplicate nodes should be deduped, got %d nodes with id 'svc-dup-svc'", count)
	}
}

// ---- FlowRenderer tests ----

func TestFlowRenderer_Filename(t *testing.T) {
	if (&FlowRenderer{}).Filename() != "flow.html" {
		t.Error("filename should be flow.html")
	}
}

func TestFlowRenderer_ProducesHTML(t *testing.T) {
	out := (&FlowRenderer{}).Render(sampleData())
	if !strings.Contains(out, "<!DOCTYPE html>") {
		t.Error("output should be HTML")
	}
	if !strings.Contains(out, "const GRAPH =") {
		t.Error("output should embed graph JSON")
	}
	if !strings.Contains(out, "<svg") {
		t.Error("output should contain SVG element")
	}
}

func TestFlowRenderer_EmptyData_ProducesValidHTML(t *testing.T) {
	out := (&FlowRenderer{}).Render(emptyComponentData())
	if !strings.Contains(out, "<!DOCTYPE html>") {
		t.Error("even empty data should produce valid HTML")
	}
	if !strings.Contains(out, "const GRAPH =") {
		t.Error("empty data should still embed graph JSON")
	}
}

func TestFlowRenderer_EmptyData_NoForbiddenClaims(t *testing.T) {
	out := (&FlowRenderer{}).Render(emptyComponentData())
	assertNoForbidden(t, out, "FlowRenderer empty data")
}

func TestFlowRenderer_FullData_NoForbiddenClaims(t *testing.T) {
	out := (&FlowRenderer{}).Render(sampleData())
	assertNoForbidden(t, out, "FlowRenderer full data")
}

func TestFlowRenderer_ComponentNameInTitle(t *testing.T) {
	data := map[string]interface{}{"component": "my-special-component"}
	out := (&FlowRenderer{}).Render(data)
	if !strings.Contains(out, "my-special-component") {
		t.Error("component name should appear in the HTML output")
	}
}

func TestFlowRenderer_GraphJSONContainsNodes(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"services": []interface{}{
			map[string]interface{}{"name": "my-svc", "type": "ClusterIP"},
		},
	}
	out := (&FlowRenderer{}).Render(data)
	if !strings.Contains(out, "svc-my-svc") {
		t.Error("rendered HTML should contain node ID from graph JSON")
	}
}
