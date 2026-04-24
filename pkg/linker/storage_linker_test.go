package linker

import (
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

func TestStorageLinkerLinksWriteToRead(t *testing.T) {
	cpg := graph.NewCPG()

	writeOp := &graph.Node{
		ID:         "db_write_1",
		Kind:       graph.NodeDBOperation,
		Name:       "db.Exec",
		File:       "handler.go",
		Line:       10,
		Table:      "users",
		Operation:  "write",
		Properties: map[string]string{"operation": "write", "table": "users"},
	}
	cpg.AddNode(writeOp)

	readOp := &graph.Node{
		ID:         "db_read_1",
		Kind:       graph.NodeDBOperation,
		Name:       "db.Query",
		File:       "handler.go",
		Line:       20,
		Table:      "users",
		Operation:  "read",
		Properties: map[string]string{"operation": "read", "table": "users"},
	}
	cpg.AddNode(readOp)

	otherRead := &graph.Node{
		ID:         "db_read_2",
		Kind:       graph.NodeDBOperation,
		Name:       "db.Query",
		File:       "handler.go",
		Line:       30,
		Table:      "products",
		Operation:  "read",
		Properties: map[string]string{"operation": "read", "table": "products"},
	}
	cpg.AddNode(otherRead)

	linker := NewStorageLinker()
	linked := linker.Link(cpg)

	if linked != 1 {
		t.Errorf("expected 1 storage link, got %d", linked)
	}

	edges := cpg.OutEdges("db_write_1")
	found := false
	for _, e := range edges {
		if e.Kind == graph.EdgeStorageLink && e.To == "db_read_1" {
			found = true
		}
	}
	if !found {
		t.Error("expected STORAGE_LINK edge from write to read")
	}

	for _, e := range cpg.OutEdges("db_write_1") {
		if e.To == "db_read_2" {
			t.Error("should not link write(users) to read(products)")
		}
	}
}

func TestStorageLinkerNoTableProperty(t *testing.T) {
	cpg := graph.NewCPG()

	cpg.AddNode(&graph.Node{
		ID:         "db1",
		Kind:       graph.NodeDBOperation,
		Name:       "db.Exec",
		Operation:  "write",
		Properties: map[string]string{"operation": "write"},
	})
	cpg.AddNode(&graph.Node{
		ID:         "db2",
		Kind:       graph.NodeDBOperation,
		Name:       "db.Query",
		Operation:  "read",
		Properties: map[string]string{"operation": "read"},
	})

	linker := NewStorageLinker()
	linked := linker.Link(cpg)

	if linked != 0 {
		t.Errorf("expected 0 storage links without table info, got %d", linked)
	}
}
