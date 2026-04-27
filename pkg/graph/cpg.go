package graph

import (
	"sync"

	"github.com/ugiordan/architecture-analyzer/pkg/arch"
)

// CPG is a thread-safe code property graph containing nodes and directed edges.
// ArchData is an optional enrichment sidecar, set once before queries run when
// --with-arch is specified. It is not core graph data.
type CPG struct {
	mu        sync.RWMutex
	nodes     map[string]*Node
	kindIndex map[NodeKind][]*Node
	outEdges  map[string][]*Edge
	inEdges   map[string][]*Edge
	ArchData  *arch.Data
}

// NewCPG creates an empty code property graph.
func NewCPG() *CPG {
	return &CPG{
		nodes:     make(map[string]*Node),
		kindIndex: make(map[NodeKind][]*Node),
		outEdges:  make(map[string][]*Edge),
		inEdges:   make(map[string][]*Edge),
	}
}

func (g *CPG) AddNode(n *Node) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.nodes[n.ID] = n
	g.kindIndex[n.Kind] = append(g.kindIndex[n.Kind], n)
}

func (g *CPG) GetNode(id string) *Node {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.nodes[id]
}

func (g *CPG) Nodes() []*Node {
	g.mu.RLock()
	defer g.mu.RUnlock()
	result := make([]*Node, 0, len(g.nodes))
	for _, n := range g.nodes {
		result = append(result, n)
	}
	return result
}

func (g *CPG) NodesByKind(kind NodeKind) []*Node {
	g.mu.RLock()
	defer g.mu.RUnlock()
	src := g.kindIndex[kind]
	result := make([]*Node, len(src))
	copy(result, src)
	return result
}

func (g *CPG) AddEdge(e *Edge) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.outEdges[e.From] = append(g.outEdges[e.From], e)
	g.inEdges[e.To] = append(g.inEdges[e.To], e)
}

func (g *CPG) OutEdges(nodeID string) []*Edge {
	g.mu.RLock()
	defer g.mu.RUnlock()
	src := g.outEdges[nodeID]
	result := make([]*Edge, len(src))
	copy(result, src)
	return result
}

func (g *CPG) InEdges(nodeID string) []*Edge {
	g.mu.RLock()
	defer g.mu.RUnlock()
	src := g.inEdges[nodeID]
	result := make([]*Edge, len(src))
	copy(result, src)
	return result
}

// SetAnnotation sets an annotation on a node in a thread-safe manner.
func (g *CPG) SetAnnotation(nodeID string, key string, value bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	n, ok := g.nodes[nodeID]
	if !ok {
		return
	}
	if n.Annotations == nil {
		n.Annotations = make(map[string]bool)
	}
	n.Annotations[key] = value
}

// SetProperty sets a property on a node in a thread-safe manner.
func (g *CPG) SetProperty(nodeID string, key string, value string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	n, ok := g.nodes[nodeID]
	if !ok {
		return
	}
	if n.Properties == nil {
		n.Properties = make(map[string]string)
	}
	n.Properties[key] = value
}

// EnsureAnnotations initializes the Annotations map on a node if nil, thread-safely.
func (g *CPG) EnsureAnnotations(nodeID string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	n, ok := g.nodes[nodeID]
	if !ok {
		return
	}
	if n.Annotations == nil {
		n.Annotations = make(map[string]bool)
	}
}

func (g *CPG) Edges() []*Edge {
	g.mu.RLock()
	defer g.mu.RUnlock()
	total := 0
	for _, edges := range g.outEdges {
		total += len(edges)
	}
	result := make([]*Edge, 0, total)
	for _, edges := range g.outEdges {
		result = append(result, edges...)
	}
	return result
}
