package diff_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/diff"
	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

func writeSnapshotFile(t *testing.T, dir, name string, snap diff.GraphSnapshot) string {
	t.Helper()
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write snapshot: %v", err)
	}
	return path
}

func buildBinary(t *testing.T) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "arch-analyzer")
	cmd := exec.Command("go", "build", "-o", binary, "./cmd/arch-analyzer/")
	cmd.Dir = filepath.Join("..", "..")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build binary: %v\n%s", err, out)
	}
	return binary
}

func TestCLIDiffExitCode0NoDifferences(t *testing.T) {
	binary := buildBinary(t)
	dir := t.TempDir()

	snap := diff.GraphSnapshot{
		SchemaVersion: 3,
		Nodes: []graph.Node{
			{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo", File: "a.go", Line: 1},
		},
	}
	base := writeSnapshotFile(t, dir, "base.json", snap)
	head := writeSnapshotFile(t, dir, "head.json", snap)

	cmd := exec.Command(binary, "diff", base, head)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("expected exit code 0 for identical snapshots, got error: %v\n%s", err, out)
	}
}

func TestCLIDiffExitCode1DifferencesFound(t *testing.T) {
	binary := buildBinary(t)
	dir := t.TempDir()

	base := writeSnapshotFile(t, dir, "base.json", diff.GraphSnapshot{
		SchemaVersion: 3,
	})
	head := writeSnapshotFile(t, dir, "head.json", diff.GraphSnapshot{
		SchemaVersion: 3,
		Nodes: []graph.Node{
			{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo"},
		},
	})

	cmd := exec.Command(binary, "diff", base, head)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit for differences found")
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected ExitError, got %T: %v", err, err)
	}
	if exitErr.ExitCode() != 1 {
		t.Errorf("expected exit code 1, got %d\n%s", exitErr.ExitCode(), out)
	}
}

func TestCLIDiffExitCode2OnError(t *testing.T) {
	binary := buildBinary(t)

	cmd := exec.Command(binary, "diff", "/nonexistent/base.json", "/nonexistent/head.json")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit for missing files")
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected ExitError, got %T: %v", err, err)
	}
	if exitErr.ExitCode() != 2 {
		t.Errorf("expected exit code 2, got %d\n%s", exitErr.ExitCode(), out)
	}
}

func TestCLIDiffMalformedJSON(t *testing.T) {
	binary := buildBinary(t)
	dir := t.TempDir()

	malformed := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(malformed, []byte("{not valid json}"), 0o644); err != nil {
		t.Fatal(err)
	}
	good := writeSnapshotFile(t, dir, "good.json", diff.GraphSnapshot{SchemaVersion: 3})

	cmd := exec.Command(binary, "diff", malformed, good)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected ExitError, got %T: %v", err, err)
	}
	if exitErr.ExitCode() != 2 {
		t.Errorf("expected exit code 2 for malformed JSON, got %d\n%s", exitErr.ExitCode(), out)
	}
}

func TestCLIDiffSchemaVersionError(t *testing.T) {
	binary := buildBinary(t)
	dir := t.TempDir()

	base := writeSnapshotFile(t, dir, "v2.json", diff.GraphSnapshot{SchemaVersion: 2})
	head := writeSnapshotFile(t, dir, "v3.json", diff.GraphSnapshot{SchemaVersion: 3})

	cmd := exec.Command(binary, "diff", base, head)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error for schema v2")
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected ExitError, got %T: %v", err, err)
	}
	if exitErr.ExitCode() != 2 {
		t.Errorf("expected exit code 2, got %d\n%s", exitErr.ExitCode(), out)
	}
}

func TestCLIDiffTextFormat(t *testing.T) {
	binary := buildBinary(t)
	dir := t.TempDir()

	base := writeSnapshotFile(t, dir, "base.json", diff.GraphSnapshot{
		SchemaVersion: 3,
	})
	head := writeSnapshotFile(t, dir, "head.json", diff.GraphSnapshot{
		SchemaVersion: 3,
		Nodes: []graph.Node{
			{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo", File: "a.go", Line: 1, Language: "go"},
		},
	})

	cmd := exec.Command(binary, "diff", "--format", "text", base, head)
	out, err := cmd.CombinedOutput()
	// Exit code 1 expected (differences found)
	if err == nil {
		t.Fatal("expected exit code 1")
	}
	output := string(out)
	if len(output) == 0 {
		t.Fatal("expected text output")
	}
	// Should contain the summary line
	if !contains(output, "Nodes:") {
		t.Errorf("text output should contain 'Nodes:' summary, got:\n%s", output)
	}
	if !contains(output, "Added:") {
		t.Errorf("text output should contain 'Added:' section, got:\n%s", output)
	}
}

func TestCLIDiffOutputFile(t *testing.T) {
	binary := buildBinary(t)
	dir := t.TempDir()

	snap := diff.GraphSnapshot{
		SchemaVersion: 3,
		Nodes: []graph.Node{
			{ID: "fn_aaa", Kind: graph.NodeFunction, Name: "foo"},
		},
	}
	base := writeSnapshotFile(t, dir, "base.json", snap)
	head := writeSnapshotFile(t, dir, "head.json", snap)
	outFile := filepath.Join(dir, "result.json")

	cmd := exec.Command(binary, "diff", "--output", outFile, base, head)
	if err := cmd.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	var result diff.GraphDiff
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("output file is not valid JSON: %v", err)
	}
	if result.HasDifferences() {
		t.Error("identical snapshots should produce no differences")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
