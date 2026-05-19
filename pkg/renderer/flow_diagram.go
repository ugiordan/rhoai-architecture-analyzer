package renderer

import (
	"fmt"
	"strings"
)

// DiagramMeta holds metadata for the diagram.
type DiagramMeta struct {
	Title string `json:"title"`
}

// CanvasSize defines the logical canvas dimensions.
type CanvasSize struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// DiagramNode is a positioned node for the Canvas renderer.
type DiagramNode struct {
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	W        float64 `json:"w"`
	H        float64 `json:"h"`
	Label    string  `json:"label"`
	Sublabel string  `json:"sublabel,omitempty"`
	Type     string  `json:"type,omitempty"` // boundary, container, icon, plugin, or "" (default box)
	Color    string  `json:"color,omitempty"`
	FontSize int     `json:"fontSize,omitempty"`
}

// DiagramStep is a single step in an animated flow.
type DiagramStep struct {
	Mode   string `json:"mode"`             // "arrow" or "lightup"
	From   string `json:"from,omitempty"`   // node key (arrow)
	To     string `json:"to,omitempty"`     // node key (arrow)
	Target string `json:"target,omitempty"` // node key (lightup)
	Text   string `json:"text"`
	Color  string `json:"color,omitempty"`
	Num    int    `json:"num,omitempty"`   // badge number (arrow)
	Badge  string `json:"badge,omitempty"` // badge text (lightup)
}

// DiagramFlow is a named sequence of steps.
type DiagramFlow struct {
	Label string        `json:"label"`
	Steps []DiagramStep `json:"steps"`
}

// DiagramTooltip holds rich tooltip data for a node.
type DiagramTooltip struct {
	Title       string      `json:"title"`
	Description string      `json:"description,omitempty"`
	Details     [][2]string `json:"details,omitempty"`
}

// LegendEntry defines a color legend item.
type LegendEntry struct {
	Label string `json:"label"`
	Color string `json:"color"`
}

// InspectorLine is a single line in the inspector panel.
type InspectorLine struct {
	Value string `json:"value"`
	Style string `json:"style"` // keep, add, highlight, err
	ID    string `json:"id"`
}

// InspectorState represents the inspector panel state.
type InspectorState struct {
	Phase   string          `json:"phase"`
	Headers []InspectorLine `json:"headers"`
	Body    []InspectorLine `json:"body"`
}

// InspectorMutation describes how the inspector changes at a step.
type InspectorMutation struct {
	Step           int             `json:"step"`
	Label          string          `json:"label"`
	Phase          string          `json:"phase,omitempty"`
	ReplaceHeaders []InspectorLine `json:"replaceHeaders,omitempty"`
	ReplaceBody    []InspectorLine `json:"replaceBody,omitempty"`
}

// DiagramInspector holds the inspector panel configuration.
type DiagramInspector struct {
	InitialState InspectorState                `json:"initialState"`
	Mutations    map[string][]InspectorMutation `json:"mutations"`
}

// FlowDiagram is the complete data model for the Canvas renderer.
type FlowDiagram struct {
	Meta        DiagramMeta                `json:"meta"`
	Canvas      CanvasSize                 `json:"canvas"`
	Nodes       map[string]DiagramNode     `json:"nodes"`
	Tooltips    map[string]DiagramTooltip  `json:"tooltips"`
	Flows       map[string]DiagramFlow     `json:"flows"`
	Inspector   *DiagramInspector          `json:"inspector,omitempty"`
	Legend      []LegendEntry              `json:"legend"`
	FlowOrder   []string                   `json:"flowOrder"`
	DefaultFlow string                     `json:"defaultFlow"`
}

// Node color presets by type.
var nodeColors = map[FlowNodeType]string{
	FlowNodeIngress:    "#2980b9",
	FlowNodeWebhook:    "#e67e22",
	FlowNodeService:    "#27ae60",
	FlowNodeDeployment: "#2471a3",
	FlowNodeExternal:   "#7f8c8d",
	FlowNodeCRD:        "#8e44ad",
}

// Flow path color presets.
var flowColors = map[string]string{
	"Request Flow":      "#3498db",
	"Controller Flow":   "#9b59b6",
	"External Calls":    "#e74c3c",
	"Webhook Intercept": "#e67e22",
}

// buildFlowDiagram converts component architecture data into a FlowDiagram
// ready for the Canvas renderer.
func buildFlowDiagram(data map[string]interface{}) FlowDiagram {
	graph := buildFlowGraph(data)
	layout := computeLayout(graph.Nodes, graph.Edges)

	component := graph.Component

	// Build positioned nodes.
	nodes := make(map[string]DiagramNode, len(graph.Nodes))
	for _, n := range graph.Nodes {
		pos := layout.Positions[n.ID]
		sublabel := buildSublabel(n)
		color := nodeColors[n.Type]
		if color == "" {
			color = "#7f8c8d"
		}
		nodes[n.ID] = DiagramNode{
			X:        pos.X,
			Y:        pos.Y,
			W:        pos.W,
			H:        pos.H,
			Label:    n.Label,
			Sublabel: sublabel,
			Color:    color,
		}
	}

	// Build tooltips.
	tooltips := buildTooltips(graph.Nodes)

	// Build flows from paths.
	flows, flowOrder := buildFlowSteps(graph.Paths, graph.Edges, graph.Nodes)

	// Build inspector.
	inspector := buildInspector(graph, flows, flowOrder)

	// Legend.
	legend := []LegendEntry{
		{Label: "Ingress", Color: "#2980b9"},
		{Label: "Webhook", Color: "#e67e22"},
		{Label: "Service", Color: "#27ae60"},
		{Label: "Deployment", Color: "#2471a3"},
		{Label: "External", Color: "#7f8c8d"},
		{Label: "CRD", Color: "#8e44ad"},
	}

	defaultFlow := ""
	if len(flowOrder) > 0 {
		defaultFlow = flowOrder[0]
	}

	return FlowDiagram{
		Meta:        DiagramMeta{Title: component + " Architecture Flow"},
		Canvas:      CanvasSize{Width: layout.CanvasWidth, Height: layout.CanvasHeight},
		Nodes:       nodes,
		Tooltips:    tooltips,
		Flows:       flows,
		Inspector:   inspector,
		Legend:      legend,
		FlowOrder:   flowOrder,
		DefaultFlow: defaultFlow,
	}
}

// buildSublabel creates a descriptive sublabel from node metadata.
func buildSublabel(n FlowNode) string {
	switch n.Type {
	case FlowNodeService:
		parts := []string{}
		if t := n.Meta["type"]; t != "" {
			parts = append(parts, t)
		}
		if p := n.Meta["ports"]; p != "" {
			parts = append(parts, p)
		}
		return strings.Join(parts, " : ")
	case FlowNodeWebhook:
		return n.Meta["type"]
	case FlowNodeExternal:
		return n.Meta["type"]
	case FlowNodeCRD:
		return n.Meta["group"]
	case FlowNodeIngress:
		return n.Meta["kind"]
	default:
		return ""
	}
}

// buildTooltips generates rich tooltips for each node.
func buildTooltips(nodes []FlowNode) map[string]DiagramTooltip {
	tooltips := make(map[string]DiagramTooltip, len(nodes))
	for _, n := range nodes {
		tt := DiagramTooltip{
			Title: n.Label,
		}
		switch n.Type {
		case FlowNodeService:
			tt.Description = "Kubernetes Service"
			if t := n.Meta["type"]; t != "" {
				tt.Details = append(tt.Details, [2]string{"Type", t})
			}
			if p := n.Meta["ports"]; p != "" {
				tt.Details = append(tt.Details, [2]string{"Ports", p})
			}
		case FlowNodeDeployment:
			tt.Description = "Deployment workload"
		case FlowNodeWebhook:
			tt.Description = "Admission webhook"
			if t := n.Meta["type"]; t != "" {
				tt.Details = append(tt.Details, [2]string{"Type", t})
			}
		case FlowNodeIngress:
			tt.Description = "Ingress / Gateway routing"
			if k := n.Meta["kind"]; k != "" {
				tt.Details = append(tt.Details, [2]string{"Kind", k})
			}
		case FlowNodeExternal:
			tt.Description = "External connection"
			if t := n.Meta["type"]; t != "" {
				tt.Details = append(tt.Details, [2]string{"Type", t})
			}
		case FlowNodeCRD:
			tt.Description = "Custom Resource Definition"
			if g := n.Meta["group"]; g != "" {
				tt.Details = append(tt.Details, [2]string{"API Group", g})
			}
			if k := n.Meta["kind"]; k != "" {
				tt.Details = append(tt.Details, [2]string{"Kind", k})
			}
		}
		tooltips[n.ID] = tt
	}
	return tooltips
}

// buildFlowSteps converts FlowPaths into ordered DiagramFlows with
// alternating arrow and lightup steps.
func buildFlowSteps(paths []FlowPath, edges []FlowEdge, nodes []FlowNode) (map[string]DiagramFlow, []string) {
	edgeByID := map[string]FlowEdge{}
	for _, e := range edges {
		edgeByID[e.ID] = e
	}

	nodeByID := map[string]FlowNode{}
	for _, n := range nodes {
		nodeByID[n.ID] = n
	}

	flows := map[string]DiagramFlow{}
	var flowOrder []string

	for _, path := range paths {
		flowKey := sanitizeID(path.Name)
		color := path.Color
		if c, ok := flowColors[path.Name]; ok {
			color = c
		}

		var steps []DiagramStep
		stepNum := 1

		for _, edgeID := range path.Edges {
			edge, ok := edgeByID[edgeID]
			if !ok {
				continue
			}
			fromNode := nodeByID[edge.From]
			toNode := nodeByID[edge.To]

			// Arrow step: animated dot from source to target.
			arrowText := fmt.Sprintf("%s → %s", fromNode.Label, toNode.Label)
			if edge.Type != "" {
				arrowText += " (" + edge.Type + ")"
			}
			steps = append(steps, DiagramStep{
				Mode:  "arrow",
				From:  edge.From,
				To:    edge.To,
				Text:  arrowText,
				Color: color,
				Num:   stepNum,
			})
			stepNum++

			// Lightup step: highlight the target node.
			steps = append(steps, DiagramStep{
				Mode:   "lightup",
				Target: edge.To,
				Text:   toNode.Label + " receives",
				Badge:  fmt.Sprintf("%d", stepNum),
			})
			stepNum++
		}

		if len(steps) > 0 {
			flows[flowKey] = DiagramFlow{
				Label: path.Name,
				Steps: steps,
			}
			flowOrder = append(flowOrder, flowKey)
		}
	}

	return flows, flowOrder
}

// buildInspector creates architecture-context inspector data for each flow.
func buildInspector(graph FlowGraph, flows map[string]DiagramFlow, flowOrder []string) *DiagramInspector {
	if len(flows) == 0 {
		return nil
	}

	nodeByID := map[string]FlowNode{}
	for _, n := range graph.Nodes {
		nodeByID[n.ID] = n
	}

	initial := InspectorState{
		Phase: "architecture",
		Headers: []InspectorLine{
			{Value: "Component: " + graph.Component, Style: "keep", ID: "h-component"},
			{Value: fmt.Sprintf("Nodes: %d", len(graph.Nodes)), Style: "keep", ID: "h-nodes"},
			{Value: fmt.Sprintf("Edges: %d", len(graph.Edges)), Style: "keep", ID: "h-edges"},
		},
		Body: []InspectorLine{},
	}

	mutations := map[string][]InspectorMutation{}

	for _, flowKey := range flowOrder {
		flow := flows[flowKey]
		var muts []InspectorMutation

		for i, step := range flow.Steps {
			mut := InspectorMutation{
				Step:  i + 1,
				Label: step.Text,
			}

			// Show context for the target node.
			targetID := step.To
			if targetID == "" {
				targetID = step.Target
			}
			if targetID != "" {
				n := nodeByID[targetID]
				var headers []InspectorLine
				headers = append(headers, InspectorLine{
					Value: "Layer: " + string(n.Type),
					Style: "highlight",
					ID:    "ctx-layer",
				})
				headers = append(headers, InspectorLine{
					Value: "Node: " + n.Label,
					Style: "keep",
					ID:    "ctx-node",
				})
				for k, v := range n.Meta {
					if v != "" {
						headers = append(headers, InspectorLine{
							Value: k + ": " + v,
							Style: "add",
							ID:    "ctx-" + k,
						})
					}
				}
				mut.ReplaceHeaders = headers
			}

			muts = append(muts, mut)
		}
		mutations[flowKey] = muts
	}

	return &DiagramInspector{
		InitialState: initial,
		Mutations:    mutations,
	}
}
