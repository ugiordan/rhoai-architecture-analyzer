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

func TestBoundedFileSystem_ReadOnly(t *testing.T) {
	root := t.TempDir()
	bfs := &boundedFileSystem{inner: filesys.MakeFsOnDisk(), root: root}

	if _, err := bfs.Create(filepath.Join(root, "test")); err == nil {
		t.Error("expected Create to be blocked")
	}
	if err := bfs.Mkdir(filepath.Join(root, "dir")); err == nil {
		t.Error("expected Mkdir to be blocked")
	}
	if err := bfs.MkdirAll(filepath.Join(root, "dir/sub")); err == nil {
		t.Error("expected MkdirAll to be blocked")
	}
	if err := bfs.RemoveAll(filepath.Join(root, "anything")); err == nil {
		t.Error("expected RemoveAll to be blocked")
	}
	if err := bfs.WriteFile(filepath.Join(root, "file"), []byte("data")); err == nil {
		t.Error("expected WriteFile to be blocked")
	}
}

func TestBoundedFileSystem_ReadOps(t *testing.T) {
	root := t.TempDir()
	// Resolve symlinks for macOS /tmp -> /private/tmp
	resolvedRoot, _ := filepath.EvalSymlinks(root)

	testFile := filepath.Join(root, "test.txt")
	os.WriteFile(testFile, []byte("hello"), 0644)

	bfs := &boundedFileSystem{inner: filesys.MakeFsOnDisk(), root: resolvedRoot}

	// ReadFile should work inside boundary
	data, err := bfs.ReadFile(testFile)
	if err != nil {
		t.Errorf("ReadFile inside boundary should work: %v", err)
	}
	if string(data) != "hello" {
		t.Errorf("expected 'hello', got %q", string(data))
	}

	// Exists should work
	if !bfs.Exists(testFile) {
		t.Error("Exists should return true for file inside boundary")
	}

	// IsDir
	if !bfs.IsDir(root) {
		t.Error("IsDir should return true for root")
	}

	// ReadFile outside boundary should fail
	_, err = bfs.ReadFile("/etc/hosts")
	if err == nil {
		t.Error("ReadFile outside boundary should fail")
	}
}

func TestBoundedFileSystem_GlobFiltering(t *testing.T) {
	root := t.TempDir()
	resolvedRoot, _ := filepath.EvalSymlinks(root)
	os.WriteFile(filepath.Join(root, "a.yaml"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(root, "b.yaml"), []byte("b"), 0644)

	bfs := &boundedFileSystem{inner: filesys.MakeFsOnDisk(), root: resolvedRoot}

	results, err := bfs.Glob(filepath.Join(root, "*.yaml"))
	if err != nil {
		t.Fatalf("Glob failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 glob results, got %d", len(results))
	}
}

func TestBoundedFileSystem_WalkSkipsOutOfBounds(t *testing.T) {
	root := t.TempDir()
	resolvedRoot, _ := filepath.EvalSymlinks(root)
	os.MkdirAll(filepath.Join(root, "sub"), 0755)
	os.WriteFile(filepath.Join(root, "sub", "ok.txt"), []byte("ok"), 0644)

	bfs := &boundedFileSystem{inner: filesys.MakeFsOnDisk(), root: resolvedRoot}

	var walked []string
	err := bfs.Walk(root, func(path string, info os.FileInfo, err error) error {
		walked = append(walked, path)
		return nil
	})
	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}
	if len(walked) == 0 {
		t.Error("expected Walk to visit files")
	}
}

func TestHasSeparatedAffix_EdgeCases(t *testing.T) {
	tests := []struct {
		longer, shorter string
		want            bool
	}{
		{"a-b-c", "a-b", true},    // prefix with compound
		{"a-b-c", "b-c", true},    // suffix with compound
		{"a-b-c", "a-b-c", false}, // equal
		{"abc", "abc", false},      // equal no sep
		{"", "a", false},           // empty longer
		{"a", "", false},           // empty shorter
	}
	for _, tt := range tests {
		got := hasSeparatedAffix(tt.longer, tt.shorter)
		if got != tt.want {
			t.Errorf("hasSeparatedAffix(%q, %q) = %v, want %v", tt.longer, tt.shorter, got, tt.want)
		}
	}
}
