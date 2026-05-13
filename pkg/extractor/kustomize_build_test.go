package extractor

import (
	"os"
	"path/filepath"
	"testing"

	"sigs.k8s.io/kustomize/kyaml/filesys"
)

func TestHasSeparatedAffix(t *testing.T) {
	tests := []struct {
		longer, shorter string
		want            bool
	}{
		{"controller-manager-metrics", "controller-manager", true},  // suffix match
		{"controller-manager-metrics", "metrics", true},             // suffix match
		{"controller-manager", "controller", true},                  // prefix match
		{"controller-manager", "controller-manager", false},         // equal length
		{"controller", "controller-manager", false},                 // shorter is longer
		{"controllermanager", "controller", false},                  // no separator
		{"my-controller", "controller", true},                       // suffix match
		{"controller-my", "controller", true},                       // prefix match
		{"xcontroller", "controller", false},                        // substring but no separator
	}
	for _, tt := range tests {
		got := hasSeparatedAffix(tt.longer, tt.shorter)
		if got != tt.want {
			t.Errorf("hasSeparatedAffix(%q, %q) = %v, want %v", tt.longer, tt.shorter, got, tt.want)
		}
	}
}

func TestBoundedFileSystem_CheckBound(t *testing.T) {
	rawRoot := t.TempDir()
	// Resolve symlinks so checkBound comparisons work on macOS (/var -> /private/var)
	root, err := filepath.EvalSymlinks(rawRoot)
	if err != nil {
		t.Fatal(err)
	}
	// Create a sibling dir to test the separator check
	sibling := root + "-sibling"
	if err := os.MkdirAll(sibling, 0755); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(sibling)

	bfs := &boundedFileSystem{inner: filesys.MakeFsOnDisk(), root: root}

	// Should allow paths inside root
	subdir := filepath.Join(root, "subdir")
	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := bfs.checkBound(subdir); err != nil {
		t.Errorf("expected subdir to pass: %v", err)
	}

	// Should allow root itself
	if err := bfs.checkBound(root); err != nil {
		t.Errorf("expected root itself to pass: %v", err)
	}

	// Should block sibling directory (the separator bug)
	if err := bfs.checkBound(sibling); err == nil {
		t.Error("expected sibling directory to be blocked")
	}

	// Should block parent
	if err := bfs.checkBound(filepath.Dir(root)); err == nil {
		t.Error("expected parent to be blocked")
	}

	// Should block /etc
	if err := bfs.checkBound("/etc"); err == nil {
		t.Error("expected /etc to be blocked")
	}
}
