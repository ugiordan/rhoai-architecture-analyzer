package renderer

import (
	"fmt"
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

// BuildFlowGraph converts component architecture data into a network-flow
// focused FlowGraph. Only includes nodes relevant to request paths:
// Client → Ingress → Webhooks → Services → Deployments → External.
// CRDs and controller watches are excluded (better suited for static diagrams).
// Exported for use by the /flow-diagram skill.
func BuildFlowGraph(data map[string]interface{}) FlowGraph {
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

	uniqueName := func(name, prefix string) string {
		key := prefix + ":" + name
		nameCounter[key]++
		if nameCounter[key] > 1 {
			return fmt.Sprintf("%s-%d", name, nameCounter[key])
		}
		return name
	}

	type nodeRef struct {
		id   string
		item map[string]interface{}
	}

	// Synthetic "Client" entry point (layer 0).
	addNode(FlowNode{
		ID: "client", Label: "Client", Type: FlowNodeIngress, Layer: 0,
		Meta: map[string]string{"kind": "icon"},
	})

	// Phase 1: extract network-relevant nodes only.
	var ingressRefs, webhookRefs, serviceRefs, externalRefs []nodeRef

	// Layer 1: ingress/gateway routing
	for _, item := range getSlice(data, "ingress_routing") {
		name := getStr(item, "name", "ingress")
		name = uniqueName(name, "ingress")
		id := "ingress-" + flowNodeID(name)
		addNode(FlowNode{ID: id, Label: name, Type: FlowNodeIngress, Layer: 1,
			Meta: map[string]string{"kind": getStr(item, "kind", "")}})
		ingressRefs = append(ingressRefs, nodeRef{id, item})
	}

	// Layer 2: webhooks (collapsed if many)
	webhooks := getSlice(data, "webhooks")
	if len(webhooks) > 3 {
		// Group by type: validating vs mutating
		valCount, mutCount := 0, 0
		var valRef, mutRef map[string]interface{}
		for _, item := range webhooks {
			whType := getStr(item, "type", "validating")
			if whType == "mutating" {
				mutCount++
				if mutRef == nil {
					mutRef = item
				}
			} else {
				valCount++
				if valRef == nil {
					valRef = item
				}
			}
		}
		if valCount > 0 {
			label := fmt.Sprintf("Validating Webhooks (%d)", valCount)
			id := "wh-validating"
			addNode(FlowNode{ID: id, Label: label, Type: FlowNodeWebhook, Layer: 2,
				Meta: map[string]string{"type": "validating", "count": fmt.Sprintf("%d", valCount)}})
			if valRef != nil {
				webhookRefs = append(webhookRefs, nodeRef{id, valRef})
			}
		}
		if mutCount > 0 {
			label := fmt.Sprintf("Mutating Webhooks (%d)", mutCount)
			id := "wh-mutating"
			addNode(FlowNode{ID: id, Label: label, Type: FlowNodeWebhook, Layer: 2,
				Meta: map[string]string{"type": "mutating", "count": fmt.Sprintf("%d", mutCount)}})
			if mutRef != nil {
				webhookRefs = append(webhookRefs, nodeRef{id, mutRef})
			}
		}
	} else {
		for _, item := range webhooks {
			name := getStr(item, "name", "webhook")
			name = uniqueName(name, "webhook")
			id := "wh-" + flowNodeID(name)
			addNode(FlowNode{ID: id, Label: name, Type: FlowNodeWebhook, Layer: 2,
				Meta: map[string]string{"type": getStr(item, "type", "")}})
			webhookRefs = append(webhookRefs, nodeRef{id, item})
		}
	}

	// Layer 3: services
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
		addNode(FlowNode{ID: id, Label: name, Type: FlowNodeService, Layer: 3, Meta: meta})
		serviceRefs = append(serviceRefs, nodeRef{id, item})
	}

	// Layer 4: deployments
	for _, item := range getSlice(data, "deployments") {
		name := getStr(item, "name", "deployment")
		name = uniqueName(name, "deployment")
		id := "dep-" + flowNodeID(name)
		addNode(FlowNode{ID: id, Label: name, Type: FlowNodeDeployment, Layer: 4})
	}

	// Layer 5: external connections (collapsed by type)
	extByType := map[string]int{}
	extFirstItem := map[string]map[string]interface{}{}
	for _, item := range getSlice(data, "external_connections") {
		target := getStr(item, "target", "")
		connType := getStr(item, "type", "external")
		if target == "" || strings.Contains(target, "%s") || strings.Contains(target, "%d") {
			target = connType
		}
		extByType[connType]++
		if extFirstItem[connType] == nil {
			extFirstItem[connType] = item
		}
	}
	for connType, count := range extByType {
		label := connType
		if count > 1 {
			label = fmt.Sprintf("%s (%d)", connType, count)
		}
		id := "ext-" + flowNodeID(connType)
		addNode(FlowNode{ID: id, Label: label, Type: FlowNodeExternal, Layer: 5,
			Meta: map[string]string{"type": connType, "count": fmt.Sprintf("%d", count)}})
		externalRefs = append(externalRefs, nodeRef{id, extFirstItem[connType]})
	}

	// Phase 2: build lookup maps.
	serviceByName := map[string]string{}
	var firstSvcID string
	for _, n := range g.Nodes {
		if n.Type == FlowNodeService {
			serviceByName[n.Label] = n.ID
			if firstSvcID == "" {
				firstSvcID = n.ID
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
		g.Edges = append(g.Edges, FlowEdge{From: from, To: to, Type: edgeType, Label: label})
	}

	// Phase 3: wire edges for network flow.

	// Client → first ingress (or first service if no ingress)
	if len(ingressRefs) > 0 {
		addEdge("client", ingressRefs[0].id, "route", "")
	} else if firstSvcID != "" {
		addEdge("client", firstSvcID, "route", "")
	}

	// Ingress → service
	for _, ref := range ingressRefs {
		svcRef := stripNamespace(getStr(ref.item, "service_ref", ""))
		if svcID, ok := serviceByName[svcRef]; ok {
			addEdge(ref.id, svcID, "route", "")
		} else if len(serviceByName) == 1 {
			addEdge(ref.id, firstSvcID, "route", "")
		}
	}

	// Webhook → service (intercept on the request path)
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

	// Deployment → external connections
	for _, ref := range externalRefs {
		addEdge(ctrlDepID, ref.id, "external", "")
	}

	// Assign stable edge IDs.
	for i := range g.Edges {
		g.Edges[i].ID = fmt.Sprintf("e-%s-%s-%s", g.Edges[i].From, g.Edges[i].To, g.Edges[i].Type)
	}

	if paths := buildFlowPaths(g); len(paths) > 0 {
		g.Paths = paths
	}
	return g
}

// buildFlowPaths generates a single end-to-end request flow tracing
// Client → Ingress → Webhook → Service → Deployment → External,
// following actual edges in the graph. Also generates per-deployment
// external connection flows if there are external nodes.
func buildFlowPaths(g FlowGraph) []FlowPath {
	edgesByFrom := map[string][]FlowEdge{}
	for _, e := range g.Edges {
		edgesByFrom[e.From] = append(edgesByFrom[e.From], e)
	}

	var paths []FlowPath

	// Main request flow: trace from "client" through the network.
	var requestEdges []string
	cursor := "client"
	visited := map[string]bool{"client": true}
	for i := 0; i < 10; i++ { // safety limit
		outEdges := edgesByFrom[cursor]
		if len(outEdges) == 0 {
			break
		}
		// Pick the best next edge: prefer route > intercept > target > external
		var best *FlowEdge
		for _, pref := range []string{"route", "intercept", "target", "external"} {
			for j := range outEdges {
				if outEdges[j].Type == pref && !visited[outEdges[j].To] {
					best = &outEdges[j]
					break
				}
			}
			if best != nil {
				break
			}
		}
		if best == nil {
			break
		}
		requestEdges = append(requestEdges, best.ID)
		visited[best.To] = true
		cursor = best.To
	}
	if len(requestEdges) > 0 {
		paths = append(paths, FlowPath{
			Name:  "Request Flow",
			Edges: requestEdges,
			Color: "#3498db",
		})
	}

	// Webhook intercept flow (if webhooks exist and have intercept edges)
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
				Name:  n.Label + " validation",
				Edges: edges,
				Color: "#e67e22",
			})
		}
	}

	// External connection flows (deployment → external)
	for _, n := range g.Nodes {
		if n.Type != FlowNodeDeployment {
			continue
		}
		var edges []string
		for _, e := range edgesByFrom[n.ID] {
			if e.Type == "external" {
				edges = append(edges, e.ID)
			}
		}
		if len(edges) > 0 {
			paths = append(paths, FlowPath{
				Name:  n.Label + " external calls",
				Edges: edges,
				Color: "#e74c3c",
			})
		}
	}

	return paths
}

