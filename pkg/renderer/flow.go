package renderer

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
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
	Type  string `json:"type"`  // route, intercept, target, watches, creates, external
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
	g := buildFlowGraph(data)
	graphJSON, _ := json.Marshal(g)

	tmpl := template.Must(template.New("flow").Parse(flowHTMLTemplate))
	var b strings.Builder
	_ = tmpl.Execute(&b, map[string]interface{}{
		"Component": g.Component,
		"GraphJSON": template.JS(graphJSON), //nolint:gosec // content is generated, not user input
	})
	return b.String()
}

// flowNodeID returns a safe HTML element ID suffix for a node label.
// Hyphens are preserved (valid in HTML IDs); spaces become underscores.
func flowNodeID(s string) string {
	var b strings.Builder
	for _, ch := range s {
		switch {
		case ch == ' ' || ch == '\t':
			b.WriteByte('_')
		default:
			b.WriteRune(ch)
		}
	}
	if b.Len() == 0 {
		return "node"
	}
	return b.String()
}

// buildFlowGraph converts component architecture data into a FlowGraph.
func buildFlowGraph(data map[string]interface{}) FlowGraph {
	g := FlowGraph{
		Component: getStr(data, "component", "unknown"),
	}

	seen := map[string]bool{}
	addNode := func(n FlowNode) {
		if seen[n.ID] {
			return
		}
		seen[n.ID] = true
		g.Nodes = append(g.Nodes, n)
	}

	// Layer 0: ingress routing
	for _, item := range getSlice(data, "ingress_routing") {
		name := getStr(item, "name", "ingress")
		addNode(FlowNode{
			ID:    "ingress-" + flowNodeID(name),
			Label: name,
			Type:  FlowNodeIngress,
			Layer: 0,
			Meta:  map[string]string{"kind": getStr(item, "kind", "")},
		})
	}

	// Layer 1: webhooks
	for _, item := range getSlice(data, "webhooks") {
		name := getStr(item, "name", "webhook")
		addNode(FlowNode{
			ID:    "wh-" + flowNodeID(name),
			Label: name,
			Type:  FlowNodeWebhook,
			Layer: 1,
			Meta:  map[string]string{"type": getStr(item, "type", "")},
		})
	}

	// Layer 2: services
	for _, item := range getSlice(data, "services") {
		name := getStr(item, "name", "service")
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
		addNode(FlowNode{
			ID:    "svc-" + flowNodeID(name),
			Label: name,
			Type:  FlowNodeService,
			Layer: 2,
			Meta:  meta,
		})
	}

	// Layer 3: deployments
	for _, item := range getSlice(data, "deployments") {
		name := getStr(item, "name", "deployment")
		addNode(FlowNode{
			ID:    "dep-" + flowNodeID(name),
			Label: name,
			Type:  FlowNodeDeployment,
			Layer: 3,
		})
	}

	// Layer 4: external connections
	for _, item := range getSlice(data, "external_connections") {
		target := getStr(item, "target", "external")
		addNode(FlowNode{
			ID:    "ext-" + flowNodeID(target),
			Label: target,
			Type:  FlowNodeExternal,
			Layer: 4,
			Meta:  map[string]string{"type": getStr(item, "type", "")},
		})
	}

	// Layer 5: CRDs
	for _, item := range getSlice(data, "crds") {
		kind := getStr(item, "kind", "CRD")
		addNode(FlowNode{
			ID:    "crd-" + flowNodeID(kind),
			Label: kind,
			Type:  FlowNodeCRD,
			Layer: 5,
			Meta:  map[string]string{"group": getStr(item, "group", "")},
		})
	}

	// Build lookup maps for edge wiring
	serviceByName := map[string]string{}
	crdByKind := map[string]string{}
	var firstDepID string
	var firstSvcID string
	for _, n := range g.Nodes {
		switch n.Type {
		case FlowNodeService:
			serviceByName[n.Label] = n.ID
			if firstSvcID == "" {
				firstSvcID = n.ID
			}
		case FlowNodeCRD:
			crdByKind[n.Label] = n.ID
		case FlowNodeDeployment:
			if firstDepID == "" {
				firstDepID = n.ID
			}
		}
	}

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

	// Ingress → service (by service_ref, or single-service heuristic)
	for _, item := range getSlice(data, "ingress_routing") {
		ingressID := "ingress-" + flowNodeID(getStr(item, "name", "ingress"))
		svcRef := getStr(item, "service_ref", "")
		if svcID, ok := serviceByName[svcRef]; ok {
			addEdge(ingressID, svcID, "route", "")
		} else if len(serviceByName) == 1 {
			addEdge(ingressID, firstSvcID, "route", "")
		}
	}

	// Webhook → service (by service_ref, which may be "namespace/name" or bare "name")
	for _, item := range getSlice(data, "webhooks") {
		whID := "wh-" + flowNodeID(getStr(item, "name", "webhook"))
		svcRef := getStr(item, "service_ref", "")
		// Strip namespace prefix if present (e.g. "ns/svc-name" → "svc-name")
		if idx := strings.LastIndex(svcRef, "/"); idx >= 0 {
			svcRef = svcRef[idx+1:]
		}
		if svcID, ok := serviceByName[svcRef]; ok {
			addEdge(whID, svcID, "intercept", "")
		} else if len(serviceByName) == 1 {
			addEdge(whID, firstSvcID, "intercept", "")
		}
	}

	// Service → deployment (by target_deployment, or single-deployment heuristic)
	for _, item := range getSlice(data, "services") {
		svcID := "svc-" + flowNodeID(getStr(item, "name", "service"))
		target := getStr(item, "target_deployment", "")
		depID := "dep-" + flowNodeID(target)
		if target != "" {
			addEdge(svcID, depID, "target", "")
		} else if firstDepID != "" {
			addEdge(svcID, firstDepID, "target", "")
		}
	}

	// Controller watches → CRDs (or built-in k8s resources for Owns)
	for _, item := range getSlice(data, "controller_watches") {
		watchType := getStr(item, "type", "")
		gvk := getStr(item, "gvk", "")
		parts := strings.Split(gvk, "/")
		kind := parts[len(parts)-1]
		if kind == "" {
			continue
		}
		if crdID, ok := crdByKind[kind]; ok {
			if watchType == "Owns" {
				addEdge(firstDepID, crdID, "creates", "")
			} else {
				addEdge(firstDepID, crdID, "watches", "")
			}
		} else if watchType == "Owns" {
			// Owns a built-in k8s resource — add it as a CRD node if not seen
			nodeID := "crd-" + flowNodeID(kind)
			addNode(FlowNode{ID: nodeID, Label: kind, Type: FlowNodeCRD, Layer: 5})
			crdByKind[kind] = nodeID
			addEdge(firstDepID, nodeID, "creates", "")
		}
	}

	// Deployment → external connections
	for _, item := range getSlice(data, "external_connections") {
		extID := "ext-" + flowNodeID(getStr(item, "target", "external"))
		addEdge(firstDepID, extID, "external", "")
	}

	// Assign stable edge IDs
	for i := range g.Edges {
		g.Edges[i].ID = fmt.Sprintf("e%d", i)
	}

	g.Paths = buildFlowPaths(g)
	return g
}

// buildFlowPaths infers animated flow paths from the graph structure.
func buildFlowPaths(g FlowGraph) []FlowPath {
	// Index edges by from node and by type
	edgesByFrom := map[string][]FlowEdge{}
	for _, e := range g.Edges {
		edgesByFrom[e.From] = append(edgesByFrom[e.From], e)
	}

	firstEdgeOfType := func(fromID, edgeType string) (FlowEdge, bool) {
		for _, e := range edgesByFrom[fromID] {
			if e.Type == edgeType {
				return e, true
			}
		}
		return FlowEdge{}, false
	}

	firstNodeOfType := func(t FlowNodeType) string {
		for _, n := range g.Nodes {
			if n.Type == t {
				return n.ID
			}
		}
		return ""
	}

	ingressID := firstNodeOfType(FlowNodeIngress)
	svcID := firstNodeOfType(FlowNodeService)
	depID := firstNodeOfType(FlowNodeDeployment)

	var paths []FlowPath

	// Request Flow: ingress → service → deployment
	if ingressID != "" && svcID != "" && depID != "" {
		var edges []string
		if e, ok := firstEdgeOfType(ingressID, "route"); ok {
			edges = append(edges, e.ID)
		}
		if e, ok := firstEdgeOfType(svcID, "target"); ok {
			edges = append(edges, e.ID)
		}
		if len(edges) > 0 {
			paths = append(paths, FlowPath{
				Name:  "Request Flow",
				Edges: edges,
				Color: "#3498db",
			})
		}
	}

	// Controller Flow: deployment → CRD watches/creates
	if depID != "" {
		var edges []string
		for _, e := range edgesByFrom[depID] {
			if e.Type == "watches" || e.Type == "creates" {
				edges = append(edges, e.ID)
			}
		}
		if len(edges) > 0 {
			paths = append(paths, FlowPath{
				Name:  "Controller Flow",
				Edges: edges,
				Color: "#9b59b6",
			})
		}
	}

	// External Calls: deployment → external connections
	if depID != "" {
		var edges []string
		for _, e := range edgesByFrom[depID] {
			if e.Type == "external" {
				edges = append(edges, e.ID)
			}
		}
		if len(edges) > 0 {
			paths = append(paths, FlowPath{
				Name:  "External Calls",
				Edges: edges,
				Color: "#e74c3c",
			})
		}
	}

	// Webhook Intercept: webhook → service
	whID := firstNodeOfType(FlowNodeWebhook)
	if whID != "" {
		var edges []string
		for _, e := range edgesByFrom[whID] {
			if e.Type == "intercept" {
				edges = append(edges, e.ID)
			}
		}
		if len(edges) > 0 {
			paths = append(paths, FlowPath{
				Name:  "Webhook Intercept",
				Edges: edges,
				Color: "#e67e22",
			})
		}
	}

	return paths
}

// flowHTMLTemplate is the self-contained HTML/JS flow diagram template.
// Uses html/template syntax: {{.Component}} and {{.GraphJSON}}.
// The JS uses no backtick template literals to stay compatible with Go raw strings.
const flowHTMLTemplate = `<!DOCTYPE html>
<html data-theme="dark">
<head>
<meta charset="utf-8">
<title>{{.Component}} - Architecture Flow</title>
<style>
:root[data-theme="dark"] {
  --bg: #0d1117; --fg: #c9d1d9; --panel-bg: #161b22; --panel-border: #30363d;
  --ctrl-bg: #161b22; --ctrl-border: #30363d; --btn-bg: #21262d; --btn-fg: #c9d1d9;
  --btn-hover: #30363d; --legend-fg: #8b949e;
}
:root[data-theme="light"] {
  --bg: #ffffff; --fg: #24292f; --panel-bg: #f6f8fa; --panel-border: #d0d7de;
  --ctrl-bg: #f6f8fa; --ctrl-border: #d0d7de; --btn-bg: #f6f8fa; --btn-fg: #24292f;
  --btn-hover: #eaeef2; --legend-fg: #57606a;
}
* { box-sizing: border-box; margin: 0; padding: 0; }
body { background: var(--bg); color: var(--fg); font-family: 'Segoe UI', monospace, sans-serif; overflow: hidden; }
#controls {
  display: flex; align-items: center; gap: 8px; padding: 8px 12px;
  background: var(--ctrl-bg); border-bottom: 1px solid var(--ctrl-border);
  height: 48px; flex-wrap: nowrap; overflow-x: auto;
}
#comp-title { font-weight: 600; font-size: 13px; white-space: nowrap; margin-right: 4px; }
button {
  background: var(--btn-bg); color: var(--btn-fg); border: 1px solid var(--ctrl-border);
  border-radius: 6px; padding: 4px 10px; cursor: pointer; font-size: 12px; white-space: nowrap;
}
button:hover { background: var(--btn-hover); }
button:disabled { opacity: 0.4; cursor: default; }
select {
  background: var(--btn-bg); color: var(--btn-fg); border: 1px solid var(--ctrl-border);
  border-radius: 6px; padding: 4px 8px; font-size: 12px; cursor: pointer;
}
label { font-size: 12px; color: var(--legend-fg); white-space: nowrap; }
input[type=range] { width: 80px; cursor: pointer; vertical-align: middle; }
#sep { flex: 1; }
#diagram { width: 100vw; height: calc(100vh - 48px); cursor: grab; display: block; }
#diagram:active { cursor: grabbing; }
#detail-panel {
  position: fixed; right: 0; top: 48px; width: 260px; height: calc(100vh - 48px);
  background: var(--panel-bg); border-left: 1px solid var(--panel-border);
  padding: 16px; overflow-y: auto; display: none; z-index: 100; font-size: 13px;
}
#detail-panel h3 { font-size: 14px; margin-bottom: 8px; word-break: break-all; }
#detail-panel .meta-row { color: var(--legend-fg); margin: 4px 0; font-size: 12px; }
#detail-panel .meta-row b { color: var(--fg); }
#detail-panel .section-title { font-size: 11px; color: var(--legend-fg); margin: 12px 0 4px; text-transform: uppercase; letter-spacing: 0.5px; }
#detail-panel .close-btn { float: right; background: none; border: none; cursor: pointer; font-size: 16px; color: var(--legend-fg); }
#legend { font-size: 11px; color: var(--legend-fg); white-space: nowrap; }
#legend span { margin-right: 10px; }
.node-label { pointer-events: none; }
</style>
</head>
<body>
<div id="controls">
  <span id="comp-title">{{.Component}}</span>
  <button id="btn-play">&#9654; Play</button>
  <button id="btn-pause">&#9646;&#9646; Pause</button>
  <select id="flow-select" title="Select flow to animate"></select>
  <label>Speed: <input type="range" id="speed" min="0.5" max="4" step="0.5" value="1"></label>
  <span id="sep"></span>
  <span id="legend"></span>
  <button id="btn-theme">&#9728; Light</button>
</div>
<svg id="diagram" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <marker id="arrow-dark" markerWidth="8" markerHeight="6" refX="7" refY="3" orient="auto">
      <path d="M0,0 L0,6 L8,3 z" fill="#8b949e"/>
    </marker>
    <marker id="arrow-route"    markerWidth="8" markerHeight="6" refX="7" refY="3" orient="auto"><path d="M0,0 L0,6 L8,3 z" fill="#2980b9"/></marker>
    <marker id="arrow-intercept" markerWidth="8" markerHeight="6" refX="7" refY="3" orient="auto"><path d="M0,0 L0,6 L8,3 z" fill="#e67e22"/></marker>
    <marker id="arrow-target"   markerWidth="8" markerHeight="6" refX="7" refY="3" orient="auto"><path d="M0,0 L0,6 L8,3 z" fill="#27ae60"/></marker>
    <marker id="arrow-watches"  markerWidth="8" markerHeight="6" refX="7" refY="3" orient="auto"><path d="M0,0 L0,6 L8,3 z" fill="#9b59b6"/></marker>
    <marker id="arrow-creates"  markerWidth="8" markerHeight="6" refX="7" refY="3" orient="auto"><path d="M0,0 L0,6 L8,3 z" fill="#a569bd"/></marker>
    <marker id="arrow-external" markerWidth="8" markerHeight="6" refX="7" refY="3" orient="auto"><path d="M0,0 L0,6 L8,3 z" fill="#e74c3c"/></marker>
  </defs>
  <g id="g-edges"></g>
  <g id="g-nodes"></g>
  <g id="g-dots"></g>
</svg>
<div id="detail-panel"></div>
<script>
const GRAPH = {{.GraphJSON}};
var SVG_NS = 'http://www.w3.org/2000/svg';
var XLINK_NS = 'http://www.w3.org/1999/xlink';

var NODE_STYLE = {
  ingress:    { w: 130, h: 38, rx: 19, fill: '#1a5276', stroke: '#2980b9' },
  webhook:    { w: 120, h: 38, rx:  5, fill: '#784212', stroke: '#e67e22' },
  service:    { w: 120, h: 38, rx:  8, fill: '#1e8449', stroke: '#27ae60' },
  deployment: { w: 140, h: 38, rx:  5, fill: '#154360', stroke: '#2471a3' },
  external:   { w: 120, h: 38, rx: 19, fill: '#424949', stroke: '#7f8c8d' },
  crd:        { w: 120, h: 38, rx:  5, fill: '#4a235a', stroke: '#8e44ad' }
};

var EDGE_COLOR = {
  route:     '#2980b9',
  intercept: '#e67e22',
  target:    '#27ae60',
  watches:   '#9b59b6',
  creates:   '#a569bd',
  external:  '#e74c3c'
};

var positions = {};
var playing = false;
var currentPathIdx = 0;
var svg = document.getElementById('diagram');

function svgEl(tag, attrs) {
  var el = document.createElementNS(SVG_NS, tag);
  for (var k in attrs) { el.setAttribute(k, attrs[k]); }
  return el;
}

function computeLayout(nodes) {
  var layers = {};
  for (var i = 0; i < nodes.length; i++) {
    var n = nodes[i];
    if (!layers[n.layer]) layers[n.layer] = [];
    layers[n.layer].push(n);
  }
  var W = svg.clientWidth || window.innerWidth;
  var H = (svg.clientHeight || window.innerHeight) - 48;
  var layerKeys = Object.keys(layers).map(Number).sort(function(a,b){return a-b;});
  var layerCount = layerKeys.length || 1;
  var result = {};
  for (var li = 0; li < layerKeys.length; li++) {
    var layer = layerKeys[li];
    var layerNodes = layers[layer];
    var y = 60 + (li + 0.5) * ((H - 80) / layerCount);
    for (var ni = 0; ni < layerNodes.length; ni++) {
      var x = (ni + 1) * W / (layerNodes.length + 1);
      result[layerNodes[ni].id] = { x: x, y: y };
    }
  }
  return result;
}

function makeNodeEl(node, pos) {
  var s = NODE_STYLE[node.type] || NODE_STYLE.service;
  var g = svgEl('g', { id: 'node-' + node.id, transform: 'translate(' + pos.x + ',' + pos.y + ')' });
  g.style.cursor = 'pointer';

  var rect = svgEl('rect', {
    x: -(s.w/2), y: -(s.h/2), width: s.w, height: s.h, rx: s.rx,
    fill: s.fill, stroke: s.stroke, 'stroke-width': '1.5'
  });
  g.appendChild(rect);

  var label = node.label.length > 17 ? node.label.slice(0, 15) + '…' : node.label;
  var text = svgEl('text', {
    'text-anchor': 'middle', dy: '0.35em',
    fill: '#ecf0f1', 'font-size': '11', 'font-family': 'monospace',
    'class': 'node-label'
  });
  text.textContent = label;
  g.appendChild(text);

  g.addEventListener('click', function() { showDetail(node); });
  return g;
}

function makeEdgeEl(edge) {
  var from = positions[edge.from], to = positions[edge.to];
  if (!from || !to) return null;
  var cx = (from.x + to.x) / 2;
  var cy1 = from.y + 30, cy2 = to.y - 30;
  var d = 'M' + from.x + ',' + from.y +
          ' C' + cx + ',' + cy1 + ' ' + cx + ',' + cy2 +
          ' ' + to.x + ',' + to.y;
  var color = EDGE_COLOR[edge.type] || '#8b949e';
  var marker = 'url(#arrow-' + (EDGE_COLOR[edge.type] ? edge.type : 'dark') + ')';
  return svgEl('path', {
    id: 'edge-' + edge.id, d: d, fill: 'none',
    stroke: color, 'stroke-width': '1.5', opacity: '0.75',
    'marker-end': marker
  });
}

function clearDots() {
  var el = document.getElementById('g-dots');
  while (el.firstChild) el.removeChild(el.firstChild);
}

function startAnimation() {
  clearDots();
  if (!GRAPH.paths || GRAPH.paths.length === 0) return;
  var flowPath = GRAPH.paths[currentPathIdx];
  var speed = parseFloat(document.getElementById('speed').value) || 1;
  var dur = (3 / speed).toFixed(1) + 's';
  var dotsEl = document.getElementById('g-dots');

  for (var i = 0; i < flowPath.edges.length; i++) {
    var edgeId = flowPath.edges[i];
    var pathEl = document.getElementById('edge-' + edgeId);
    if (!pathEl) continue;

    var circle = svgEl('circle', { r: '6', fill: flowPath.color });
    circle.setAttribute('filter', 'drop-shadow(0 0 4px ' + flowPath.color + ')');

    var motion = svgEl('animateMotion', {
      dur: dur,
      begin: (i * parseFloat(dur) * 0.3).toFixed(1) + 's',
      repeatCount: 'indefinite',
      rotate: 'auto'
    });
    var mpath = document.createElementNS(SVG_NS, 'mpath');
    mpath.setAttributeNS(XLINK_NS, 'xlink:href', '#edge-' + edgeId);
    motion.appendChild(mpath);
    circle.appendChild(motion);
    dotsEl.appendChild(circle);
  }
  playing = true;
}

function showDetail(node) {
  var panel = document.getElementById('detail-panel');
  panel.style.display = 'block';

  var incoming = GRAPH.edges.filter(function(e) { return e.to === node.id; });
  var outgoing = GRAPH.edges.filter(function(e) { return e.from === node.id; });

  var html = '<button class="close-btn" onclick="document.getElementById(\'detail-panel\').style.display=\'none\'">&#10005;</button>';
  html += '<h3>' + escHtml(node.label) + '</h3>';
  html += '<p class="meta-row" style="margin-bottom:8px">' + escHtml(node.type) + '</p>';

  var meta = node.meta || {};
  for (var k in meta) {
    if (meta[k]) html += '<p class="meta-row"><b>' + escHtml(k) + ':</b> ' + escHtml(meta[k]) + '</p>';
  }

  if (incoming.length > 0) {
    html += '<p class="section-title">Receives from</p>';
    for (var i = 0; i < incoming.length; i++) {
      html += '<p class="meta-row">← ' + escHtml(incoming[i].from) + ' <span style="color:var(--legend-fg)">(' + escHtml(incoming[i].type) + ')</span></p>';
    }
  }
  if (outgoing.length > 0) {
    html += '<p class="section-title">Sends to</p>';
    for (var j = 0; j < outgoing.length; j++) {
      html += '<p class="meta-row">→ ' + escHtml(outgoing[j].to) + ' <span style="color:var(--legend-fg)">(' + escHtml(outgoing[j].type) + ')</span></p>';
    }
  }
  panel.innerHTML = html;
}

function escHtml(s) {
  return String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');
}

function renderLegend() {
  var types = [
    ['ingress','#2980b9'], ['webhook','#e67e22'], ['service','#27ae60'],
    ['deployment','#2471a3'], ['external','#7f8c8d'], ['crd','#8e44ad']
  ];
  var html = '';
  for (var i = 0; i < types.length; i++) {
    html += '<span><svg width="10" height="10" style="vertical-align:middle;margin-right:3px">' +
            '<rect width="10" height="10" rx="2" fill="' + types[i][1] + '"/></svg>' +
            types[i][0] + '</span>';
  }
  document.getElementById('legend').innerHTML = html;
}

function populateFlowSelect() {
  var sel = document.getElementById('flow-select');
  sel.innerHTML = '';
  if (!GRAPH.paths || GRAPH.paths.length === 0) {
    var opt = document.createElement('option');
    opt.textContent = '(no flows)';
    sel.appendChild(opt);
    document.getElementById('btn-play').disabled = true;
    document.getElementById('btn-pause').disabled = true;
    return;
  }
  for (var i = 0; i < GRAPH.paths.length; i++) {
    var opt = document.createElement('option');
    opt.value = i;
    opt.textContent = GRAPH.paths[i].name;
    sel.appendChild(opt);
  }
}

function drawAll() {
  positions = computeLayout(GRAPH.nodes);
  var edgesEl = document.getElementById('g-edges');
  edgesEl.innerHTML = '';
  for (var i = 0; i < GRAPH.edges.length; i++) {
    var el = makeEdgeEl(GRAPH.edges[i]);
    if (el) edgesEl.appendChild(el);
  }
  var nodesEl = document.getElementById('g-nodes');
  nodesEl.innerHTML = '';
  for (var j = 0; j < GRAPH.nodes.length; j++) {
    var pos = positions[GRAPH.nodes[j].id];
    if (pos) nodesEl.appendChild(makeNodeEl(GRAPH.nodes[j], pos));
  }
}

// Pan and zoom
var vb = { x: 0, y: 0, w: 0, h: 0 };
var drag = false, dragOrigin = null, vbOrigin = null;

function initViewBox() {
  vb.x = 0; vb.y = 0;
  vb.w = svg.clientWidth || window.innerWidth;
  vb.h = (svg.clientHeight || window.innerHeight) - 48;
  svg.setAttribute('viewBox', vb.x + ' ' + vb.y + ' ' + vb.w + ' ' + vb.h);
}

svg.addEventListener('mousedown', function(e) {
  drag = true;
  dragOrigin = { x: e.clientX, y: e.clientY };
  vbOrigin = { x: vb.x, y: vb.y, w: vb.w, h: vb.h };
});
svg.addEventListener('mousemove', function(e) {
  if (!drag) return;
  var sx = vb.w / (svg.clientWidth || 1);
  var sy = vb.h / (svg.clientHeight || 1);
  vb.x = vbOrigin.x - (e.clientX - dragOrigin.x) * sx;
  vb.y = vbOrigin.y - (e.clientY - dragOrigin.y) * sy;
  svg.setAttribute('viewBox', vb.x + ' ' + vb.y + ' ' + vb.w + ' ' + vb.h);
});
svg.addEventListener('mouseup', function() { drag = false; });
svg.addEventListener('mouseleave', function() { drag = false; });
svg.addEventListener('wheel', function(e) {
  e.preventDefault();
  var scale = e.deltaY > 0 ? 1.15 : 0.87;
  var mx = e.offsetX * (vb.w / (svg.clientWidth || 1)) + vb.x;
  var my = e.offsetY * (vb.h / (svg.clientHeight || 1)) + vb.y;
  vb.x = mx - (mx - vb.x) * scale;
  vb.y = my - (my - vb.y) * scale;
  vb.w *= scale; vb.h *= scale;
  svg.setAttribute('viewBox', vb.x + ' ' + vb.y + ' ' + vb.w + ' ' + vb.h);
}, { passive: false });

// Controls
document.getElementById('btn-play').addEventListener('click', function() {
  playing = true; startAnimation();
});
document.getElementById('btn-pause').addEventListener('click', function() {
  playing = false; clearDots();
});
document.getElementById('speed').addEventListener('input', function() {
  if (playing) startAnimation();
});
document.getElementById('flow-select').addEventListener('change', function(e) {
  currentPathIdx = parseInt(e.target.value) || 0;
  if (playing) startAnimation();
});
document.getElementById('btn-theme').addEventListener('click', function() {
  var html = document.documentElement;
  if (html.getAttribute('data-theme') === 'dark') {
    html.setAttribute('data-theme', 'light');
    this.innerHTML = '&#127769; Dark';
  } else {
    html.setAttribute('data-theme', 'dark');
    this.innerHTML = '&#9728; Light';
  }
});

window.addEventListener('load', function() {
  initViewBox();
  drawAll();
  populateFlowSelect();
  renderLegend();
});
window.addEventListener('resize', function() {
  clearDots();
  initViewBox();
  drawAll();
  if (playing) startAnimation();
});
</script>
</body>
</html>`
