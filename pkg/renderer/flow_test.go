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

func TestBuildFlowGraph_ServiceTargetDeploymentHappyPath(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"services": []interface{}{
			map[string]interface{}{"name": "svc", "target_deployment": "my-dep"},
		},
		"deployments": []interface{}{
			map[string]interface{}{"name": "my-dep"},
			map[string]interface{}{"name": "other-dep"},
		},
	}
	g := buildFlowGraph(data)
	var found bool
	for _, e := range g.Edges {
		if e.Type == "target" && e.From == "svc-svc" && e.To == "dep-my-dep" {
			found = true
		}
	}
	if !found {
		t.Error("explicit target_deployment should wire to the named deployment, not fallback")
	}
}

func TestBuildFlowGraph_CRDCrossGroupLookup(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"deployments": []interface{}{
			map[string]interface{}{"name": "ctrl"},
		},
		"crds": []interface{}{
			map[string]interface{}{"kind": "Certificate", "group": "cert-manager.io"},
			map[string]interface{}{"kind": "Certificate", "group": "networking.knative.dev"},
		},
		"controller_watches": []interface{}{
			map[string]interface{}{"type": "For", "gvk": "networking.knative.dev/v1alpha1/Certificate"},
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
		t.Fatal("should produce a watches edge for cross-group CRD")
	}
	// The watch should target the knative Certificate, not the cert-manager one
	certManagerID := ""
	knativeID := ""
	for _, n := range g.Nodes {
		if n.Type == FlowNodeCRD && n.Meta["group"] == "cert-manager.io" {
			certManagerID = n.ID
		}
		if n.Type == FlowNodeCRD && n.Meta["group"] == "networking.knative.dev" {
			knativeID = n.ID
		}
	}
	if knativeID == "" {
		t.Fatal("should have a knative Certificate CRD node")
	}
	if watchEdge.To == certManagerID {
		t.Error("watch for networking.knative.dev/Certificate should NOT wire to the cert-manager CRD")
	}
	if watchEdge.To != knativeID {
		t.Errorf("watch should wire to knative Certificate node %q, got %q", knativeID, watchEdge.To)
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

// ---- edge wiring with uniqueName tests ----

func TestBuildFlowGraph_UniqueNameEdgeSync(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"services": []interface{}{
			map[string]interface{}{"type": "ClusterIP"},
			map[string]interface{}{"type": "NodePort"},
		},
		"deployments": []interface{}{
			map[string]interface{}{"name": "app"},
		},
	}
	g := buildFlowGraph(data)
	nodeIDs := map[string]bool{}
	for _, n := range g.Nodes {
		nodeIDs[n.ID] = true
	}
	for _, e := range g.Edges {
		if !nodeIDs[e.From] {
			t.Errorf("edge From %q does not correspond to any node", e.From)
		}
		if !nodeIDs[e.To] {
			t.Errorf("edge To %q does not correspond to any node", e.To)
		}
	}
}

func TestBuildFlowGraph_UniqueNameExternalEdgeSync(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"deployments": []interface{}{
			map[string]interface{}{"name": "app"},
		},
		"external_connections": []interface{}{
			map[string]interface{}{"type": "database"},
			map[string]interface{}{"type": "cache"},
		},
	}
	g := buildFlowGraph(data)
	nodeIDs := map[string]bool{}
	for _, n := range g.Nodes {
		nodeIDs[n.ID] = true
	}
	for _, e := range g.Edges {
		if e.Type == "external" && !nodeIDs[e.To] {
			t.Errorf("external edge To %q does not correspond to any node", e.To)
		}
	}
}

// ---- single-service heuristic test ----

func TestBuildFlowGraph_SingleServiceHeuristic(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"services": []interface{}{
			map[string]interface{}{"name": "only-svc"},
		},
		"ingress_routing": []interface{}{
			map[string]interface{}{"name": "gw", "service_ref": "nonexistent"},
		},
	}
	g := buildFlowGraph(data)
	var found bool
	for _, e := range g.Edges {
		if e.Type == "route" && e.To == "svc-only-svc" {
			found = true
		}
	}
	if !found {
		t.Error("with one service and non-matching service_ref, should fall back to single-service heuristic")
	}
}

func TestBuildFlowGraph_MultiServiceNoFallback(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"services": []interface{}{
			map[string]interface{}{"name": "svc-a"},
			map[string]interface{}{"name": "svc-b"},
		},
		"ingress_routing": []interface{}{
			map[string]interface{}{"name": "gw", "service_ref": "nonexistent"},
		},
	}
	g := buildFlowGraph(data)
	for _, e := range g.Edges {
		if e.Type == "route" {
			t.Error("with multiple services and non-matching service_ref, should NOT create a route edge")
		}
	}
}

// ---- flow path type tests ----

func TestBuildFlowPaths_ControllerFlow(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"deployments": []interface{}{
			map[string]interface{}{"name": "ctrl"},
		},
		"crds": []interface{}{
			map[string]interface{}{"kind": "Widget", "group": "g"},
		},
		"controller_watches": []interface{}{
			map[string]interface{}{"type": "For", "gvk": "g/v1/Widget"},
		},
	}
	g := buildFlowGraph(data)
	var found *FlowPath
	for i := range g.Paths {
		if g.Paths[i].Name == "Controller Flow" {
			found = &g.Paths[i]
		}
	}
	if found == nil {
		t.Fatal("should have a 'Controller Flow' path")
	}
	if len(found.Edges) == 0 {
		t.Error("Controller Flow should have at least one edge")
	}
}

func TestBuildFlowPaths_ExternalCalls(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"deployments": []interface{}{
			map[string]interface{}{"name": "app"},
		},
		"external_connections": []interface{}{
			map[string]interface{}{"target": "redis", "type": "cache"},
		},
	}
	g := buildFlowGraph(data)
	var found *FlowPath
	for i := range g.Paths {
		if g.Paths[i].Name == "External Calls" {
			found = &g.Paths[i]
		}
	}
	if found == nil {
		t.Fatal("should have an 'External Calls' path")
	}
	if len(found.Edges) == 0 {
		t.Error("External Calls should have at least one edge")
	}
}

func TestBuildFlowPaths_WebhookIntercept(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"services": []interface{}{
			map[string]interface{}{"name": "webhook-svc"},
		},
		"webhooks": []interface{}{
			map[string]interface{}{"name": "validate", "service_ref": "webhook-svc"},
		},
	}
	g := buildFlowGraph(data)
	var found *FlowPath
	for i := range g.Paths {
		if g.Paths[i].Name == "Webhook Intercept" {
			found = &g.Paths[i]
		}
	}
	if found == nil {
		t.Fatal("should have a 'Webhook Intercept' path")
	}
	if len(found.Edges) == 0 {
		t.Error("Webhook Intercept should have at least one edge")
	}
}

func TestBuildFlowPaths_RequestFlowFollowsActualRoute(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"services": []interface{}{
			map[string]interface{}{"name": "frontend"},
			map[string]interface{}{"name": "backend"},
		},
		"ingress_routing": []interface{}{
			map[string]interface{}{"name": "gw", "service_ref": "backend"},
		},
		"deployments": []interface{}{
			map[string]interface{}{"name": "test-app"},
		},
	}
	g := buildFlowGraph(data)
	var routeEdge *FlowEdge
	for i := range g.Edges {
		if g.Edges[i].Type == "route" {
			routeEdge = &g.Edges[i]
		}
	}
	if routeEdge == nil {
		t.Fatal("should have a route edge")
	}
	if routeEdge.To != "svc-backend" {
		t.Errorf("route edge should go to svc-backend, got %q", routeEdge.To)
	}
	var reqFlow *FlowPath
	for i := range g.Paths {
		if g.Paths[i].Name == "Request Flow" {
			reqFlow = &g.Paths[i]
		}
	}
	if reqFlow == nil {
		t.Fatal("should have a Request Flow path")
	}
	if len(reqFlow.Edges) < 1 {
		t.Fatal("Request Flow should have edges")
	}
	// Verify the path follows the actual route, not the first service
	edgeByID := map[string]FlowEdge{}
	for _, e := range g.Edges {
		edgeByID[e.ID] = e
	}
	firstEdge := edgeByID[reqFlow.Edges[0]]
	if firstEdge.To != "svc-backend" {
		t.Errorf("Request Flow first edge should go to svc-backend (the routed service), got %q", firstEdge.To)
	}
}

// ---- controllerDepID edge From verification ----

func TestBuildFlowGraph_ControllerWatchEdgeFrom(t *testing.T) {
	data := map[string]interface{}{
		"component": "test",
		"deployments": []interface{}{
			map[string]interface{}{"name": "test-controller"},
		},
		"crds": []interface{}{
			map[string]interface{}{"kind": "MyKind", "group": "g"},
		},
		"controller_watches": []interface{}{
			map[string]interface{}{"type": "For", "gvk": "g/v1/MyKind"},
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
		t.Fatal("should have a watches edge")
	}
	if watchEdge.From != "dep-test-controller" {
		t.Errorf("watches edge From should be dep-test-controller, got %q", watchEdge.From)
	}
}

// ---- nil slice safety test ----

func TestBuildFlowGraph_EmptyData_NilSafety(t *testing.T) {
	g := buildFlowGraph(emptyComponentData())
	if g.Nodes == nil {
		t.Error("Nodes should be empty slice, not nil (JSON serializes nil as null)")
	}
	if g.Edges == nil {
		t.Error("Edges should be empty slice, not nil")
	}
	if g.Paths == nil {
		t.Error("Paths should be empty slice, not nil")
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
	if !strings.Contains(out, "var D =") {
		t.Error("output should embed diagram JSON")
	}
	if !strings.Contains(out, "<canvas") {
		t.Error("output should contain canvas element")
	}
}

func TestFlowRenderer_EmptyData_ProducesValidHTML(t *testing.T) {
	out := (&FlowRenderer{}).Render(emptyComponentData())
	if !strings.Contains(out, "<!DOCTYPE html>") {
		t.Error("even empty data should produce valid HTML")
	}
	if !strings.Contains(out, "var D =") {
		t.Error("empty data should still embed diagram JSON")
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

func TestFlowRenderer_HTMLContainsCanvasEngine(t *testing.T) {
	out := (&FlowRenderer{}).Render(sampleData())
	if !strings.Contains(out, "requestAnimationFrame") {
		t.Error("should contain Canvas render loop")
	}
	if !strings.Contains(out, "renderAll") {
		t.Error("should contain renderAll function")
	}
}

func TestFlowRenderer_HTMLContainsPlayback(t *testing.T) {
	out := (&FlowRenderer{}).Render(sampleData())
	if !strings.Contains(out, "execStep") {
		t.Error("should contain step execution logic")
	}
	if !strings.Contains(out, "btn-play") {
		t.Error("should contain play button")
	}
}

func TestFlowRenderer_HTMLContainsInspector(t *testing.T) {
	out := (&FlowRenderer{}).Render(sampleData())
	if !strings.Contains(out, "inspector-panel") {
		t.Error("should contain inspector panel")
	}
}

func TestFlowRenderer_DiagramHasFlows(t *testing.T) {
	out := (&FlowRenderer{}).Render(sampleData())
	if !strings.Contains(out, "flowOrder") {
		t.Error("diagram JSON should contain flowOrder")
	}
}
