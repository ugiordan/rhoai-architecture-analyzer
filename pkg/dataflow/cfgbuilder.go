package dataflow

import "github.com/ugiordan/architecture-analyzer/pkg/graph"

// MaxBlocksPerFunction is the limit on NodeBasicBlock nodes per function.
// When exceeded, CFG construction stops for that function.
const MaxBlocksPerFunction = 200

// CFGBuilder accumulates basic blocks and control flow edges during
// CFG construction. One instance per function body analysis.
type CFGBuilder struct {
	blocks []*graph.Node
	edges  []*graph.Edge
	funcID string
	file   string
}

// NewCFGBuilder creates a CFGBuilder for the given function.
func NewCFGBuilder(funcID, file string) *CFGBuilder {
	return &CFGBuilder{
		funcID: funcID,
		file:   file,
	}
}

// NewBlock creates a NodeBasicBlock node. The block ID includes the funcID
// to ensure uniqueness across functions in the same file.
func (cb *CFGBuilder) NewBlock(name string, line int) *graph.Node {
	block := &graph.Node{
		ID:       graph.NodeID(graph.NodeBasicBlock, cb.funcID+"/"+name, cb.file, line, 0),
		Kind:     graph.NodeBasicBlock,
		Name:     name,
		File:     cb.file,
		Line:     line,
		ParentID: cb.funcID,
	}
	cb.blocks = append(cb.blocks, block)
	return block
}

// AddMember appends a node ID to a block's Members list.
func (cb *CFGBuilder) AddMember(block *graph.Node, memberID string) {
	block.Members = append(block.Members, memberID)
}

// AddEdge creates an EdgeControlFlow edge between two blocks.
func (cb *CFGBuilder) AddEdge(fromBlockID, toBlockID, label string) {
	cb.edges = append(cb.edges, &graph.Edge{
		From:       fromBlockID,
		To:         toBlockID,
		Kind:       graph.EdgeControlFlow,
		Label:      label,
		Confidence: graph.ConfidenceCertain,
	})
}

// BlockCount returns the number of blocks created so far.
func (cb *CFGBuilder) BlockCount() int {
	return len(cb.blocks)
}

// Result returns all accumulated blocks and edges.
func (cb *CFGBuilder) Result() ([]*graph.Node, []*graph.Edge) {
	return cb.blocks, cb.edges
}
