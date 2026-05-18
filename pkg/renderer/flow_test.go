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
		if n.ID == "svc-my-svc" && n.Type == FlowNodeService && n.Layer == 2 {
			found = true
		}
	}
	if !found {
		t.Error("service should produce a FlowNodeService node with id 'svc-my-svc' at layer 2")
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
		if n.Type == FlowNodeWebhook && n.Layer == 1 && n.Meta["type"] == "validating" {
			found = true
		}
	}
	if !found {
		t.Error("webhook should produce a FlowNodeWebhook node at layer 1 with type meta")
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
		if n.Type == FlowNodeExternal && n.Layer == 4 {
			found = true
		}
	}
	if !found {
		t.Error("external connection should produce FlowNodeExternal node at layer 4")
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
		if n.Type == FlowNodeCRD && n.ID == "crd-MyKind" && n.Layer == 5 {
			found = true
		}
	}
	if !found {
		t.Error("CRD should produce a FlowNodeCRD node with id 'crd-MyKind' at layer 5")
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
		if n.Type == FlowNodeDeployment && n.ID == "dep-my-controller" && n.Layer == 3 {
			found = true
		}
	}
	if !found {
		t.Error("deployment should produce a FlowNodeDeployment node at layer 3")
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
		if n.Type == FlowNodeIngress && n.Layer == 0 {
			found = true
		}
	}
	if !found {
		t.Error("ingress routing should produce a FlowNodeIngress node at layer 0")
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
		if e.Type == "watches" && e.To == "crd-MyKind" {
			found = true
		}
	}
	if !found {
		t.Error("controller watch should produce a 'watches' edge to the CRD")
	}
}

func TestBuildFlowGraph_OwnsCreatesEdge(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"deployments": []interface{}{
			map[string]interface{}{"name": "my-controller"},
		},
		"controller_watches": []interface{}{
			map[string]interface{}{"type": "Owns", "gvk": "/v1/ConfigMap"},
		},
	}
	g := buildFlowGraph(data)
	var foundEdge, foundNode bool
	for _, e := range g.Edges {
		if e.Type == "creates" && e.To == "crd-ConfigMap" {
			foundEdge = true
		}
	}
	for _, n := range g.Nodes {
		if n.ID == "crd-ConfigMap" && n.Type == FlowNodeCRD {
			foundNode = true
		}
	}
	if !foundEdge {
		t.Error("Owns watch should produce a 'creates' edge")
	}
	if !foundNode {
		t.Error("Owns watch for built-in resource should create a synthetic CRD node")
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
		t.Fatal("should build at least one flow path when ingress+service+deployment are present")
	}
	var requestFlow *FlowPath
	for i := range g.Paths {
		if g.Paths[i].Name == "Request Flow" {
			requestFlow = &g.Paths[i]
		}
	}
	if requestFlow == nil {
		t.Fatal("should have a 'Request Flow' path")
	}
	if len(requestFlow.Edges) == 0 {
		t.Error("Request Flow path should have at least one edge")
	}
	// Verify edge IDs in path actually exist in g.Edges
	edgeIDs := map[string]bool{}
	for _, e := range g.Edges {
		edgeIDs[e.ID] = true
	}
	for _, eid := range requestFlow.Edges {
		if !edgeIDs[eid] {
			t.Errorf("path references edge ID %q which doesn't exist in g.Edges", eid)
		}
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

// ---- Edge wiring tests ----

func TestBuildFlowGraph_WebhookServiceRefWithNamespace(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"services": []interface{}{
			map[string]interface{}{"name": "webhook-svc"},
		},
		"webhooks": []interface{}{
			map[string]interface{}{
				"name":        "validate-thing",
				"type":        "validating",
				"service_ref": "opendatahub/webhook-svc",
			},
		},
	}
	g := buildFlowGraph(data)
	var found bool
	for _, e := range g.Edges {
		if e.Type == "intercept" && e.To == "svc-webhook-svc" {
			found = true
		}
	}
	if !found {
		t.Error("webhook with namespace-qualified service_ref should produce intercept edge after stripping namespace")
	}
}

func TestBuildFlowGraph_IngressServiceRefWithNamespace(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"services": []interface{}{
			map[string]interface{}{"name": "my-svc"},
		},
		"ingress_routing": []interface{}{
			map[string]interface{}{
				"name":        "my-route",
				"service_ref": "default/my-svc",
			},
		},
	}
	g := buildFlowGraph(data)
	var found bool
	for _, e := range g.Edges {
		if e.Type == "route" && e.To == "svc-my-svc" {
			found = true
		}
	}
	if !found {
		t.Error("ingress with namespace-qualified service_ref should produce route edge after stripping namespace")
	}
}

func TestBuildFlowGraph_MultiDeployment_ControllerPreference(t *testing.T) {
	data := map[string]interface{}{
		"component": "kserve",
		"deployments": []interface{}{
			map[string]interface{}{"name": "webhook-server"},
			map[string]interface{}{"name": "kserve-controller"},
		},
		"crds": []interface{}{
			map[string]interface{}{"kind": "InferenceService", "group": "serving.kserve.io"},
		},
		"controller_watches": []interface{}{
			map[string]interface{}{"type": "For", "gvk": "serving.kserve.io/v1/InferenceService"},
		},
	}
	g := buildFlowGraph(data)
	var watchEdge *FlowEdge
	for i := range g.Edges {
		if g.Edges[i].Type == "watches" {
			watchEdge = &g.Edges[i]
		}
	}
	if watchEdge == nil {
		t.Fatal("should produce a watches edge")
	}
	if watchEdge.From != "dep-kserve-controller" {
		t.Errorf("watches edge should come from 'dep-kserve-controller' (matches component name 'kserve'), got %q", watchEdge.From)
	}
}

func TestBuildFlowGraph_ServiceTargetDeploymentMissing(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"services": []interface{}{
			map[string]interface{}{"name": "svc", "target_deployment": "nonexistent"},
		},
		"deployments": []interface{}{
			map[string]interface{}{"name": "real-dep"},
		},
	}
	g := buildFlowGraph(data)
	for _, e := range g.Edges {
		if e.Type == "target" && e.To == "dep-nonexistent" {
			t.Error("should not create edge to nonexistent deployment")
		}
	}
}

func TestBuildFlowGraph_GVKTrailingSlash(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"deployments": []interface{}{
			map[string]interface{}{"name": "ctrl"},
		},
		"crds": []interface{}{
			map[string]interface{}{"kind": "MyKind", "group": "g"},
		},
		"controller_watches": []interface{}{
			map[string]interface{}{"type": "For", "gvk": "g/v1/MyKind/"},
		},
	}
	g := buildFlowGraph(data)
	var found bool
	for _, e := range g.Edges {
		if e.Type == "watches" && e.To == "crd-MyKind" {
			found = true
		}
	}
	if !found {
		t.Error("GVK with trailing slash should still produce a watches edge after trimming")
	}
}

func TestBuildFlowGraph_DefaultNameCollision(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"services": []interface{}{
			map[string]interface{}{"type": "ClusterIP"},
			map[string]interface{}{"type": "NodePort"},
		},
	}
	g := buildFlowGraph(data)
	svcCount := 0
	for _, n := range g.Nodes {
		if n.Type == FlowNodeService {
			svcCount++
		}
	}
	if svcCount < 2 {
		t.Errorf("two unnamed services should both appear (got %d), not silently dropped by ID collision", svcCount)
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
	if !strings.Contains(out, "<title>my-special-component") {
		t.Error("component name should appear in the HTML title")
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

func TestFlowRenderer_HTMLContainsPlainHref(t *testing.T) {
	out := (&FlowRenderer{}).Render(sampleData())
	if !strings.Contains(out, "mpath.setAttribute('href'") {
		t.Error("mpath should use plain href for SVG2/Firefox compatibility")
	}
}

func TestFlowRenderer_HTMLContainsNodeByIdLookup(t *testing.T) {
	out := (&FlowRenderer{}).Render(sampleData())
	if !strings.Contains(out, "nodeLabel(") {
		t.Error("showDetail should use nodeLabel() to display human-readable names")
	}
}

func TestFlowRenderer_HTMLContainsEmptyMessage(t *testing.T) {
	out := (&FlowRenderer{}).Render(emptyComponentData())
	if !strings.Contains(out, "No architecture data found in analyzed sources") {
		t.Error("empty data should include a user-facing empty-state message in the template")
	}
}
