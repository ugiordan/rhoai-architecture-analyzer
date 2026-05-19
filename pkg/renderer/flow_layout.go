package renderer

import (
	"sort"
)

// NodeLayout holds the computed position and size for a node.
type NodeLayout struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	W float64 `json:"w"`
	H float64 `json:"h"`
}

// LayoutResult contains positions for all nodes and the canvas dimensions.
type LayoutResult struct {
	Positions    map[string]NodeLayout
	CanvasWidth  float64
	CanvasHeight float64
}

// Node size presets by type.
var nodeSize = map[FlowNodeType]struct{ w, h float64 }{
	FlowNodeIngress:    {140, 44},
	FlowNodeWebhook:    {130, 44},
	FlowNodeService:    {140, 44},
	FlowNodeDeployment: {160, 50},
	FlowNodeExternal:   {130, 44},
	FlowNodeCRD:        {130, 40},
}

const (
	layoutPadX    = 60.0  // horizontal padding from canvas edge
	layoutPadY    = 50.0  // vertical padding from canvas edge
	layoutGapX    = 50.0  // horizontal gap between nodes in same layer
	layoutGapY    = 90.0  // vertical gap between layers
	layoutMinW    = 1200.0
	layoutDefaultW = 140.0
	layoutDefaultH = 44.0
)

// computeLayout assigns x/y/w/h to every node using a layered graph layout.
// Nodes are grouped by Layer, ordered within each layer to reduce edge
// crossings (barycentric heuristic), then positioned with even spacing.
func computeLayout(nodes []FlowNode, edges []FlowEdge) LayoutResult {
	if len(nodes) == 0 {
		return LayoutResult{
			Positions:    map[string]NodeLayout{},
			CanvasWidth:  layoutMinW,
			CanvasHeight: 400,
		}
	}

	// Group nodes by layer.
	layers := map[int][]FlowNode{}
	for _, n := range nodes {
		layers[n.Layer] = append(layers[n.Layer], n)
	}

	layerKeys := make([]int, 0, len(layers))
	for k := range layers {
		layerKeys = append(layerKeys, k)
	}
	sort.Ints(layerKeys)

	// Build adjacency for barycentric ordering.
	adj := map[string][]string{} // node ID → connected node IDs
	for _, e := range edges {
		adj[e.From] = append(adj[e.From], e.To)
		adj[e.To] = append(adj[e.To], e.From)
	}

	// Order nodes within each layer using barycentric method:
	// position each node at the average x of its neighbors in adjacent layers.
	// Run a few passes top-down then bottom-up.
	posIndex := map[string]float64{} // node ID → position index within its layer
	for _, lk := range layerKeys {
		for i, n := range layers[lk] {
			posIndex[n.ID] = float64(i)
		}
	}

	for pass := 0; pass < 4; pass++ {
		keys := layerKeys
		if pass%2 == 1 {
			keys = make([]int, len(layerKeys))
			copy(keys, layerKeys)
			for i, j := 0, len(keys)-1; i < j; i, j = i+1, j-1 {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
		for _, lk := range keys {
			type ranked struct {
				node FlowNode
				bary float64
			}
			var items []ranked
			for _, n := range layers[lk] {
				neighbors := adj[n.ID]
				if len(neighbors) == 0 {
					items = append(items, ranked{n, posIndex[n.ID]})
					continue
				}
				sum := 0.0
				for _, nb := range neighbors {
					sum += posIndex[nb]
				}
				items = append(items, ranked{n, sum / float64(len(neighbors))})
			}
			sort.Slice(items, func(i, j int) bool {
				return items[i].bary < items[j].bary
			})
			reordered := make([]FlowNode, len(items))
			for i, it := range items {
				reordered[i] = it.node
				posIndex[it.node.ID] = float64(i)
			}
			layers[lk] = reordered
		}
	}

	// Compute node sizes and canvas dimensions.
	positions := map[string]NodeLayout{}
	maxLayerWidth := 0.0

	for _, lk := range layerKeys {
		layerNodes := layers[lk]
		totalW := 0.0
		for _, n := range layerNodes {
			sz := nodeSize[n.Type]
			if sz.w == 0 {
				sz.w = layoutDefaultW
				sz.h = layoutDefaultH
			}
			totalW += sz.w
		}
		totalW += float64(len(layerNodes)-1) * layoutGapX
		if totalW > maxLayerWidth {
			maxLayerWidth = totalW
		}
	}

	canvasW := maxLayerWidth + 2*layoutPadX
	if canvasW < layoutMinW {
		canvasW = layoutMinW
	}
	canvasH := float64(len(layerKeys))*(layoutDefaultH+layoutGapY) + 2*layoutPadY

	// Assign positions: center each layer horizontally.
	for li, lk := range layerKeys {
		layerNodes := layers[lk]
		y := layoutPadY + float64(li)*(layoutDefaultH+layoutGapY)

		// Compute total width of this layer.
		sizes := make([]struct{ w, h float64 }, len(layerNodes))
		totalW := 0.0
		for i, n := range layerNodes {
			sz := nodeSize[n.Type]
			if sz.w == 0 {
				sz.w = layoutDefaultW
				sz.h = layoutDefaultH
			}
			sizes[i] = sz
			totalW += sz.w
		}
		totalW += float64(len(layerNodes)-1) * layoutGapX

		// Center the layer.
		startX := (canvasW - totalW) / 2
		x := startX

		for i, n := range layerNodes {
			positions[n.ID] = NodeLayout{
				X: x,
				Y: y,
				W: sizes[i].w,
				H: sizes[i].h,
			}
			x += sizes[i].w + layoutGapX
		}
	}

	return LayoutResult{
		Positions:    positions,
		CanvasWidth:  canvasW,
		CanvasHeight: canvasH,
	}
}
