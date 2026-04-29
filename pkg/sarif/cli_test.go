package sarif_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

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

func TestCLIIngestStandalone(t *testing.T) {
	binary := buildBinary(t)
	dir := t.TempDir()

	sarifContent := `{
		"version": "2.1.0",
		"runs": [{
			"tool": {"driver": {"name": "test-scanner", "version": "1.0"}},
			"results": [{
				"ruleId": "R1",
				"level": "error",
				"message": {"text": "test finding"},
				"locations": [{
					"physicalLocation": {
						"artifactLocation": {"uri": "main.go"},
						"region": {"startLine": 10, "startColumn": 1}
					}
				}]
			}]
		}]
	}`
	sarifFile := filepath.Join(dir, "results.sarif")
	os.WriteFile(sarifFile, []byte(sarifContent), 0o644)

	cmd := exec.Command(binary, "ingest", sarifFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("ingest failed: %v\n%s", err, out)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, out)
	}
	nodes, ok := result["nodes"].([]interface{})
	if !ok || len(nodes) != 1 {
		t.Fatalf("expected 1 node in output, got: %v", result)
	}
}

func TestCLIIngestWithGraph(t *testing.T) {
	binary := buildBinary(t)
	dir := t.TempDir()

	graphContent := map[string]interface{}{
		"schema_version": 3,
		"nodes": []map[string]interface{}{
			{"id": "fn_abc", "kind": "Function", "name": "handler", "file": "main.go", "line": 10, "end_line": 20},
		},
		"edges": []interface{}{},
	}
	graphData, _ := json.MarshalIndent(graphContent, "", "  ")
	graphFile := filepath.Join(dir, "code-graph.json")
	os.WriteFile(graphFile, graphData, 0o644)

	sarifContent := `{
		"version": "2.1.0",
		"runs": [{
			"tool": {"driver": {"name": "semgrep"}},
			"results": [{
				"ruleId": "xss",
				"level": "error",
				"message": {"text": "XSS"},
				"locations": [{"physicalLocation": {"artifactLocation": {"uri": "main.go"}, "region": {"startLine": 15}}}]
			}]
		}]
	}`
	sarifFile := filepath.Join(dir, "results.sarif")
	os.WriteFile(sarifFile, []byte(sarifContent), 0o644)

	outFile := filepath.Join(dir, "enriched.json")
	cmd := exec.Command(binary, "ingest", "--graph", graphFile, "--output", outFile, sarifFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("ingest failed: %v\nOutput:\n%s", err, out)
	}
	if len(out) > 0 {
		t.Logf("Command output:\n%s", out)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	var enriched map[string]interface{}
	if err := json.Unmarshal(data, &enriched); err != nil {
		t.Fatalf("failed to parse output JSON: %v", err)
	}

	nodes, ok := enriched["nodes"].([]interface{})
	if !ok {
		t.Fatalf("nodes field missing or wrong type in output: %+v", enriched)
	}
	if len(nodes) != 2 {
		t.Errorf("expected 2 nodes (function + finding), got %d", len(nodes))
	}
}

func TestCLIIngestInvalidSARIF(t *testing.T) {
	binary := buildBinary(t)
	dir := t.TempDir()

	sarifFile := filepath.Join(dir, "bad.sarif")
	os.WriteFile(sarifFile, []byte("{invalid}"), 0o644)

	cmd := exec.Command(binary, "ingest", sarifFile)
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error for invalid SARIF")
	}
}
