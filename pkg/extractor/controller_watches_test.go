package extractor

import "testing"

func TestExtractResourceOps(t *testing.T) {
	pkgs := loadGoPackages(fixtureDir())
	if pkgs == nil {
		t.Fatal("failed to load fixture packages")
	}
	ops := extractResourceOps(pkgs)
	if len(ops) == 0 {
		t.Fatal("expected ResourceOps from fixture controller")
	}
	kinds := make(map[string]bool)
	for _, op := range ops {
		kinds[op.Kind] = true
	}
	if !kinds["Service"] {
		t.Error("expected ResourceOp for Service")
	}
	if !kinds["Deployment"] {
		t.Error("expected ResourceOp for Deployment")
	}
	for _, op := range ops {
		if op.Verb != "create" {
			t.Errorf("expected verb=create, got %s", op.Verb)
		}
	}
}
