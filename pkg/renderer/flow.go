package renderer

import (
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"strings"
	"unicode"
)

// FlowNodeType controls the visual style and layer assignment of a node.
type FlowNodeType string

const (
	FlowNodeIngress    FlowNodeType = "ingress"    // layer 0 — entry points (HTTPRoute, Gateway, Route)
	FlowNodeWebhook    FlowNodeType = "webhook"    // layer 1 — API intercepts
	FlowNodeService    FlowNodeType = "service"    // layer 2 — Kubernetes Services
	FlowNodeDeployment FlowNodeType = "deployment" // layer 3 — workloads
	FlowNodeExternal   FlowNodeType = "external"   // layer 4 — out-of-cluster connections
	FlowNodeCRD        FlowNodeType = "crd"        // layer 5 — managed custom resources
)

// FlowNode is a node in the architecture flow graph.
type FlowNode struct {
	ID    string            `json:"id"`
	Label string            `json:"label"`
	Type  FlowNodeType      `json:"type"`
	Layer int               `json:"layer"`
	Meta  map[string]string `json:"meta,omitempty"`
}

// FlowEdge is a directed edge between two nodes.
type FlowEdge struct {
	ID    string `json:"id"`
	From  string `json:"from"`
	To    string `json:"to"`
	Type  string `json:"type"` // route, intercept, target, watches, creates, external
	Label string `json:"label,omitempty"`
}

// FlowPath is a pre-computed sequence of edge IDs for animated dot traversal.
type FlowPath struct {
	Name  string   `json:"name"`
	Edges []string `json:"edges"` // edge IDs in traversal order
	Color string   `json:"color"` // dot color
}

// FlowGraph is the complete data model passed to the HTML template.
type FlowGraph struct {
	Component string     `json:"component"`
	Nodes     []FlowNode `json:"nodes"`
	Edges     []FlowEdge `json:"edges"`
	Paths     []FlowPath `json:"paths"`
}

// FlowRenderer generates a self-contained interactive HTML flow diagram.
type FlowRenderer struct{}

func (r *FlowRenderer) Filename() string { return "flow.html" }

func (r *FlowRenderer) Render(data map[string]interface{}) string {
	diagram := buildFlowDiagram(data)
	diagramJSON, err := json.Marshal(diagram)
	if err != nil {
		diagramJSON = []byte(`{"meta":{"title":"error"},"canvas":{"width":800,"height":400},"nodes":{},"tooltips":{},"flows":{},"legend":[],"flowOrder":[],"defaultFlow":""}`)
	}

	tmpl := template.Must(template.New("flow").Parse(flowCanvasTemplate))
	var b strings.Builder
	if err := tmpl.Execute(&b, map[string]interface{}{
		"Title":       diagram.Meta.Title,
		"DiagramJSON": template.JS(diagramJSON), //nolint:gosec // generated from struct, not user input; json.Marshal HTML-escapes </>
	}); err != nil {
		return "<!DOCTYPE html><html><body>Flow render error: " + html.EscapeString(err.Error()) + "</body></html>"
	}
	return b.String()
}

// flowNodeID returns a safe HTML element ID suffix for a node label.
// Letters, digits, hyphens, and underscores are preserved; everything else
// (dots, slashes, colons, spaces) becomes a hyphen.
func flowNodeID(s string) string {
	var b strings.Builder
	for _, ch := range s {
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '-' || ch == '_' {
			b.WriteRune(ch)
		} else {
			b.WriteByte('-')
		}
	}
	result := b.String()
	if result == "" {
		return "node"
	}
	return result
}

// stripNamespace removes a "namespace/" prefix from a service reference.
func stripNamespace(ref string) string {
	if idx := strings.LastIndex(ref, "/"); idx >= 0 {
		return ref[idx+1:]
	}
	return ref
}

// controllerDepID picks the best deployment for controller watch edges.
// Prefers a deployment whose name contains the component name.
func controllerDepID(nodes []FlowNode, component string) string {
	var firstDep string
	comp := strings.ToLower(component)
	for _, n := range nodes {
		if n.Type != FlowNodeDeployment {
			continue
		}
		if firstDep == "" {
			firstDep = n.ID
		}
		if strings.Contains(strings.ToLower(n.Label), comp) {
			return n.ID
		}
	}
	return firstDep
}

// buildFlowGraph converts component architecture data into a FlowGraph.
func buildFlowGraph(data map[string]interface{}) FlowGraph {
	g := FlowGraph{
		Component: getStr(data, "component", "unknown"),
		Nodes:     []FlowNode{},
		Edges:     []FlowEdge{},
		Paths:     []FlowPath{},
	}

	seen := map[string]bool{}
	nameCounter := map[string]int{}
	addNode := func(n FlowNode) {
		if seen[n.ID] {
			return
		}
		seen[n.ID] = true
		g.Nodes = append(g.Nodes, n)
	}

	// uniqueName deduplicates fallback default names
	uniqueName := func(name, prefix string) string {
		key := prefix + ":" + name
		nameCounter[key]++
		if nameCounter[key] > 1 {
			return fmt.Sprintf("%s-%d", name, nameCounter[key])
		}
		return name
	}

	// nodeRef ties a raw data item to its created node ID.
	type nodeRef struct {
		id   string
		item map[string]interface{}
	}

	// Phase 1: create all nodes, storing refs for edge wiring.
	var ingressRefs, webhookRefs, serviceRefs, externalRefs []nodeRef

	// Layer 0: ingress routing
	for _, item := range getSlice(data, "ingress_routing") {
		name := getStr(item, "name", "ingress")
		name = uniqueName(name, "ingress")
		id := "ingress-" + flowNodeID(name)
		addNode(FlowNode{ID: id, Label: name, Type: FlowNodeIngress, Layer: 0,
			Meta: map[string]string{"kind": getStr(item, "kind", "")}})
		ingressRefs = append(ingressRefs, nodeRef{id, item})
	}

	// Layer 1: webhooks
	for _, item := range getSlice(data, "webhooks") {
		name := getStr(item, "name", "webhook")
		name = uniqueName(name, "webhook")
		id := "wh-" + flowNodeID(name)
		addNode(FlowNode{ID: id, Label: name, Type: FlowNodeWebhook, Layer: 1,
			Meta: map[string]string{"type": getStr(item, "type", "")}})
		webhookRefs = append(webhookRefs, nodeRef{id, item})
	}

	// Layer 2: services
	for _, item := range getSlice(data, "services") {
		name := getStr(item, "name", "service")
		name = uniqueName(name, "service")
		svcType := getStr(item, "type", "ClusterIP")
		ports := getSlice(item, "ports")
		var portParts []string
		for _, p := range ports {
			portParts = append(portParts, fmt.Sprintf("%d", getInt(p, "port")))
		}
		meta := map[string]string{"type": svcType}
		if len(portParts) > 0 {
			meta["ports"] = strings.Join(portParts, ", ")
		}
		id := "svc-" + flowNodeID(name)
		addNode(FlowNode{ID: id, Label: name, Type: FlowNodeService, Layer: 2, Meta: meta})
		serviceRefs = append(serviceRefs, nodeRef{id, item})
	}

	// Layer 3: deployments
	for _, item := range getSlice(data, "deployments") {
		name := getStr(item, "name", "deployment")
		name = uniqueName(name, "deployment")
		id := "dep-" + flowNodeID(name)
		addNode(FlowNode{ID: id, Label: name, Type: FlowNodeDeployment, Layer: 3})
	}

	// Layer 4: external connections
	for _, item := range getSlice(data, "external_connections") {
		target := getStr(item, "target", "")
		connType := getStr(item, "type", "external")
		if target == "" || strings.Contains(target, "%s") || strings.Contains(target, "%d") {
			target = connType
		}
		target = uniqueName(target, "external")
		id := "ext-" + flowNodeID(target)
		addNode(FlowNode{ID: id, Label: target, Type: FlowNodeExternal, Layer: 4,
			Meta: map[string]string{"type": connType}})
		externalRefs = append(externalRefs, nodeRef{id, item})
	}

	// Layer 5: CRDs (group-qualified to avoid kind collisions across API groups)
	// Store the original kind in Meta so crdLookup can use the un-mangled name.
	for _, item := range getSlice(data, "crds") {
		origKind := getStr(item, "kind", "CRD")
		group := getStr(item, "group", "")
		displayKind := uniqueName(origKind, "crd")
		id := "crd-" + flowNodeID(displayKind)
		addNode(FlowNode{ID: id, Label: displayKind, Type: FlowNodeCRD, Layer: 5,
			Meta: map[string]string{"group": group, "kind": origKind}})
	}

	// Phase 2: build lookup maps for edge wiring.
	// crdLookup maps both bare "Kind" and "group/Kind" to the node ID.
	// Uses Meta["kind"] (the original, un-mangled kind) for the lookup key
	// so controller watches can match by their GVK-extracted kind.
	serviceByName := map[string]string{}
	crdLookup := map[string]string{}
	var firstSvcID string
	for _, n := range g.Nodes {
		switch n.Type {
		case FlowNodeService:
			serviceByName[n.Label] = n.ID
			if firstSvcID == "" {
				firstSvcID = n.ID
			}
		case FlowNodeCRD:
			origKind := n.Meta["kind"]
			if origKind == "" {
				origKind = n.Label
			}
			if _, exists := crdLookup[origKind]; !exists {
				crdLookup[origKind] = n.ID
			}
			if group := n.Meta["group"]; group != "" {
				crdLookup[group+"/"+origKind] = n.ID
			}
		}
	}

	ctrlDepID := controllerDepID(g.Nodes, g.Component)

	edgeSeen := map[string]bool{}
	addEdge := func(from, to, edgeType, label string) {
		if from == "" || to == "" {
			return
		}
		key := from + "|" + to + "|" + edgeType
		if edgeSeen[key] {
			return
		}
		edgeSeen[key] = true
		g.Edges = append(g.Edges, FlowEdge{
			From:  from,
			To:    to,
			Type:  edgeType,
			Label: label,
		})
	}

	// Phase 3: wire edges using stored node refs (not re-reading raw data).

	// Ingress → service
	for _, ref := range ingressRefs {
		svcRef := stripNamespace(getStr(ref.item, "service_ref", ""))
		if svcID, ok := serviceByName[svcRef]; ok {
			addEdge(ref.id, svcID, "route", "")
		} else if len(serviceByName) == 1 {
			addEdge(ref.id, firstSvcID, "route", "")
		}
	}

	// Webhook → service
	for _, ref := range webhookRefs {
		svcRef := stripNamespace(getStr(ref.item, "service_ref", ""))
		if svcID, ok := serviceByName[svcRef]; ok {
			addEdge(ref.id, svcID, "intercept", "")
		} else if len(serviceByName) == 1 {
			addEdge(ref.id, firstSvcID, "intercept", "")
		}
	}

	// Service → deployment
	for _, ref := range serviceRefs {
		target := getStr(ref.item, "target_deployment", "")
		if target != "" {
			depID := "dep-" + flowNodeID(target)
			if seen[depID] {
				addEdge(ref.id, depID, "target", "")
			}
		} else if ctrlDepID != "" {
			addEdge(ref.id, ctrlDepID, "target", "")
		}
	}

	// Controller watches → CRDs (or built-in k8s resources for Owns)
	// Try group-qualified lookup first ("group/Kind"), then bare kind.
	for _, item := range getSlice(data, "controller_watches") {
		watchType := getStr(item, "type", "")
		gvk := strings.TrimRight(getStr(item, "gvk", ""), "/")
		parts := strings.Split(gvk, "/")
		kind := parts[len(parts)-1]
		if kind == "" {
			continue
		}
		group := ""
		if len(parts) >= 3 {
			group = parts[0]
		}
		crdID := ""
		if group != "" {
			crdID = crdLookup[group+"/"+kind]
		}
		if crdID == "" {
			crdID = crdLookup[kind]
		}
		if crdID != "" {
			if watchType == "Owns" {
				addEdge(ctrlDepID, crdID, "creates", "")
			} else {
				addEdge(ctrlDepID, crdID, "watches", "")
			}
		} else if watchType == "Owns" {
			nodeID := "crd-" + flowNodeID(kind)
			addNode(FlowNode{ID: nodeID, Label: kind, Type: FlowNodeCRD, Layer: 5})
			crdLookup[kind] = nodeID
			addEdge(ctrlDepID, nodeID, "creates", "")
		}
	}

	// Deployment → external connections
	for _, ref := range externalRefs {
		addEdge(ctrlDepID, ref.id, "external", "")
	}

	// Collapse similar nodes to reduce visual clutter.
	g = collapseFlowGraph(g)

	// Assign stable edge IDs
	for i := range g.Edges {
		g.Edges[i].ID = fmt.Sprintf("e%d", i)
	}

	if paths := buildFlowPaths(g); len(paths) > 0 {
		g.Paths = paths
	}
	return g
}

// collapseFlowGraph merges nodes of the same type that share the same
// meta characteristics into a single group node, rewiring all edges.
// This reduces visual clutter (e.g. 12 webhooks → "Webhooks (12)").
const collapseThreshold = 3

func collapseFlowGraph(g FlowGraph) FlowGraph {
	// Group nodes by type+meta key for potential collapsing.
	type groupKey struct {
		nodeType FlowNodeType
		metaKey  string // the distinguishing meta value
	}
	groups := map[groupKey][]FlowNode{}
	for _, n := range g.Nodes {
		mk := ""
		switch n.Type {
		case FlowNodeWebhook:
			mk = n.Meta["type"] // group by validating/mutating
		case FlowNodeExternal:
			mk = n.Meta["type"] // group by connection type
		case FlowNodeCRD:
			mk = "crd"
		default:
			// Deployments, services, ingress: never collapse
			mk = n.ID
		}
		gk := groupKey{n.Type, mk}
		groups[gk] = append(groups[gk], n)
	}

	// Build replacement map: old node ID → new collapsed node ID.
	replace := map[string]string{}
	var newNodes []FlowNode

	for gk, nodes := range groups {
		if len(nodes) < collapseThreshold {
			newNodes = append(newNodes, nodes...)
			continue
		}
		// Collapse this group into a single node.
		var label string
		switch gk.nodeType {
		case FlowNodeWebhook:
			label = fmt.Sprintf("%s webhooks (%d)", gk.metaKey, len(nodes))
		case FlowNodeExternal:
			label = fmt.Sprintf("%s (%d)", gk.metaKey, len(nodes))
		case FlowNodeCRD:
			label = fmt.Sprintf("CRDs (%d)", len(nodes))
		default:
			label = fmt.Sprintf("%s (%d)", gk.metaKey, len(nodes))
		}
		groupID := flowNodeID(label)
		collapsed := FlowNode{
			ID:    groupID,
			Label: label,
			Type:  gk.nodeType,
			Layer: nodes[0].Layer,
			Meta:  map[string]string{"count": fmt.Sprintf("%d", len(nodes))},
		}
		// Preserve first node's meta for tooltips.
		for k, v := range nodes[0].Meta {
			if collapsed.Meta[k] == "" {
				collapsed.Meta[k] = v
			}
		}
		newNodes = append(newNodes, collapsed)
		for _, n := range nodes {
			replace[n.ID] = groupID
		}
	}

	// Rewire edges to point to collapsed nodes; dedup.
	edgeSeen := map[string]bool{}
	var newEdges []FlowEdge
	for _, e := range g.Edges {
		from := e.From
		if r, ok := replace[from]; ok {
			from = r
		}
		to := e.To
		if r, ok := replace[to]; ok {
			to = r
		}
		key := from + "|" + to + "|" + e.Type
		if edgeSeen[key] || from == to {
			continue
		}
		edgeSeen[key] = true
		newEdges = append(newEdges, FlowEdge{From: from, To: to, Type: e.Type, Label: e.Label})
	}

	if newNodes == nil {
		newNodes = []FlowNode{}
	}
	if newEdges == nil {
		newEdges = []FlowEdge{}
	}
	return FlowGraph{
		Component: g.Component,
		Nodes:     newNodes,
		Edges:     newEdges,
		Paths:     []FlowPath{},
	}
}

// buildFlowPaths generates animated flow paths for each deployment
// that has outgoing edges, not just the first one.
func buildFlowPaths(g FlowGraph) []FlowPath {
	edgesByFrom := map[string][]FlowEdge{}
	for _, e := range g.Edges {
		edgesByFrom[e.From] = append(edgesByFrom[e.From], e)
	}

	nodeByID := map[string]FlowNode{}
	for _, n := range g.Nodes {
		nodeByID[n.ID] = n
	}

	firstEdgeOfType := func(fromID, edgeType string) (FlowEdge, bool) {
		for _, e := range edgesByFrom[fromID] {
			if e.Type == edgeType {
				return e, true
			}
		}
		return FlowEdge{}, false
	}

	var paths []FlowPath

	// Request Flow: for each ingress, trace ingress → service → deployment
	for _, n := range g.Nodes {
		if n.Type != FlowNodeIngress {
			continue
		}
		var edges []string
		if routeEdge, ok := firstEdgeOfType(n.ID, "route"); ok {
			edges = append(edges, routeEdge.ID)
			if targetEdge, ok := firstEdgeOfType(routeEdge.To, "target"); ok {
				edges = append(edges, targetEdge.ID)
			}
		}
		if len(edges) > 0 {
			name := "Request Flow"
			if n.Label != "" {
				name = n.Label + " request"
			}
			paths = append(paths, FlowPath{Name: name, Edges: edges, Color: "#3498db"})
		}
	}

	// Per-deployment flows: controller watches, external calls
	for _, n := range g.Nodes {
		if n.Type != FlowNodeDeployment {
			continue
		}
		outEdges := edgesByFrom[n.ID]
		if len(outEdges) == 0 {
			continue
		}

		// Collect edges by type for this deployment.
		var watchEdges, externalEdges []string
		for _, e := range outEdges {
			switch e.Type {
			case "watches", "creates":
				watchEdges = append(watchEdges, e.ID)
			case "external":
				externalEdges = append(externalEdges, e.ID)
			}
		}

		depLabel := n.Label
		if watchEdges != nil {
			paths = append(paths, FlowPath{
				Name:  depLabel + " reconcile",
				Edges: watchEdges,
				Color: "#9b59b6",
			})
		}
		if externalEdges != nil {
			paths = append(paths, FlowPath{
				Name:  depLabel + " external",
				Edges: externalEdges,
				Color: "#e74c3c",
			})
		}
	}

	// Webhook intercepts: for each webhook with outgoing intercept edges
	for _, n := range g.Nodes {
		if n.Type != FlowNodeWebhook {
			continue
		}
		var edges []string
		for _, e := range edgesByFrom[n.ID] {
			if e.Type == "intercept" {
				edges = append(edges, e.ID)
			}
		}
		if len(edges) > 0 {
			paths = append(paths, FlowPath{
				Name:  n.Label + " intercept",
				Edges: edges,
				Color: "#e67e22",
			})
		}
	}

	return paths
}
