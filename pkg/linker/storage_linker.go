package linker

import (
	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

type StorageLinker struct{}

func NewStorageLinker() *StorageLinker {
	return &StorageLinker{}
}

func (sl *StorageLinker) Link(cpg *graph.CPG) int {
	dbOps := cpg.NodesByKind(graph.NodeDBOperation)

	writes := make(map[string][]*graph.Node)
	reads := make(map[string][]*graph.Node)

	for _, op := range dbOps {
		table := op.Table
		if table == "" {
			continue
		}
		switch op.Operation {
		case "write":
			writes[table] = append(writes[table], op)
		case "read":
			reads[table] = append(reads[table], op)
		}
	}

	linked := 0
	for table, writeOps := range writes {
		readOps, ok := reads[table]
		if !ok {
			continue
		}
		for _, w := range writeOps {
			for _, r := range readOps {
				cpg.AddEdge(&graph.Edge{
					From:  w.ID,
					To:    r.ID,
					Kind:  graph.EdgeStorageLink,
					Label: "table:" + table,
				})
				linked++
			}
		}
	}

	return linked
}
