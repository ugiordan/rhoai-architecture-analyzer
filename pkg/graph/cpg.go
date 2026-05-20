package graph

import (
	"fmt"
	"sync"

	"github.com/ugiordan/architecture-analyzer/pkg/arch"
)

// SchemaVersion is the current graph schema version used when serializing.
const SchemaVersion = 3

// CPG is a thread-safe code property graph containing nodes and directed edges.
// ArchData is an optional enrichment sidecar, set once before queries run when
// --with-arch is specified. It is not core graph data.
type CPG struct {
	mu              sync.RWMutex
	nodes           map[string]*Node
	kindIndex       map[NodeKind][]*Node
	outEdges        map[string][]*Edge
	inEdges         map[string][]*Edge
	outEdgesByKind  map[string]map[EdgeKind][]*Edge // nodeID → kind → edges
	ArchData        *arch.Data
}

// NewCPG creates an empty code property graph.
func NewCPG() *CPG {
	return &CPG{
		nodes:          make(map[string]*Node),
		kindIndex:      make(map[NodeKind][]*Node),
		outEdges:       make(map[string][]*Edge),
		inEdges:        make(map[string][]*Edge),
		outEdgesByKind: make(map[string]map[EdgeKind][]*Edge),
	}
}

func (g *CPG) AddNode(n *Node) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if existing, ok := g.nodes[n.ID]; ok {
		// Skip silently when the collision is an identical node (same name, file, line).
		// Python/TypeScript parsers emit duplicate Variable nodes for re-assignments,
		// augmented assignments, and list comprehensions at the same source location.
		if existing.Kind == n.Kind && existing.Name == n.Name &&
			existing.File == n.File && existing.Line == n.Line {
			return nil
		}
		return fmt.Errorf("duplicate node ID %q: existing node %q (%s:%d) collides with %q (%s:%d)",
			n.ID, existing.Name, existing.File, existing.Line, n.Name, n.File, n.Line)
	}
	g.nodes[n.ID] = n
	g.kindIndex[n.Kind] = append(g.kindIndex[n.Kind], n)
	return nil
}

func (g *CPG) GetNode(id string) *Node {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.nodes[id]
}

// NodeCount returns the number of nodes without allocating a copy.
func (g *CPG) NodeCount() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.nodes)
}

// KindCount returns the number of nodes of a specific kind without allocating.
func (g *CPG) KindCount(kind NodeKind) int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.kindIndex[kind])
}

// AnnotationCount returns the total number of annotations without allocating a node copy.
func (g *CPG) AnnotationCount() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	count := 0
	for _, n := range g.nodes {
		count += len(n.Annotations)
	}
	return count
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
	if g.outEdgesByKind[e.From] == nil {
		g.outEdgesByKind[e.From] = make(map[EdgeKind][]*Edge)
	}
	g.outEdgesByKind[e.From][e.Kind] = append(g.outEdgesByKind[e.From][e.Kind], e)
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

// EdgesByKindFrom returns all edges of the given kind originating from a specific node.
// Uses a compound index for O(1) lookup instead of linear scan.
func (g *CPG) EdgesByKindFrom(kind EdgeKind, fromID string) []*Edge {
	g.mu.RLock()
	defer g.mu.RUnlock()
	kindMap := g.outEdgesByKind[fromID]
	if kindMap == nil {
		return nil
	}
	src := kindMap[kind]
	result := make([]*Edge, len(src))
	copy(result, src)
	return result
}

// SetTrustLevel sets the trust level on a node in a thread-safe manner.
func (g *CPG) SetTrustLevel(nodeID string, level TrustLevel) {
	g.mu.Lock()
	defer g.mu.Unlock()
	n, ok := g.nodes[nodeID]
	if !ok {
		return
	}
	n.TrustLevel = level
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

// SARIFFindings returns all ExternalFinding nodes in the graph, providing
// a flat accessor for SARIF-ingested findings without traversing edges.
func (g *CPG) SARIFFindings() []*Node {
	return g.NodesByKind(NodeExternalFinding)
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
