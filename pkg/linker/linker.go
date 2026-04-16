package linker

import "github.com/ugiordan/rhoai-architecture-analyzer/pkg/graph"

type Linker interface {
	Link(cpg *graph.CPG) int
}
