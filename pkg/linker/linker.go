package linker

import "github.com/ugiordan/architecture-analyzer/pkg/graph"

type Linker interface {
	Link(cpg *graph.CPG) int
}
