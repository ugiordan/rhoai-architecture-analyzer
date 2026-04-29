package sarif_test

import (
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
	"github.com/ugiordan/architecture-analyzer/pkg/parser"
	"github.com/ugiordan/architecture-analyzer/pkg/sarif"
)

func makeCPGWithNodes(t *testing.T, nodes ...*graph.Node) *graph.CPG {
	t.Helper()
	cpg := graph.NewCPG()
	for _, n := range nodes {
		if err := cpg.AddNode(n); err != nil {
			t.Fatalf("AddNode: %v", err)
		}
	}
	return cpg
}

func makeReport(toolName, toolVersion string, results ...sarif.Result) *sarif.Report {
	return &sarif.Report{
		Version: "2.1.0",
		Runs: []sarif.Run{
			{
				Tool: sarif.Tool{
					Driver: sarif.ToolComponent{
						Name:    toolName,
						Version: toolVersion,
					},
				},
				Results: results,
			},
		},
	}
}

func makeResult(ruleID, level, msg, file string, line int) sarif.Result {
	return sarif.Result{
		RuleID:  ruleID,
		Level:   level,
		Message: sarif.Message{Text: msg},
		Locations: []sarif.Location{
			{
				PhysicalLocation: sarif.PhysicalLocation{
					ArtifactLocation: sarif.ArtifactLocation{URI: file},
					Region:           sarif.Region{StartLine: line, StartColumn: 1},
				},
			},
		},
	}
}

func TestIngestExactMatch(t *testing.T) {
	fnID := parser.NodeID(graph.NodeFunction, "handleRequest", "pkg/handler/serve.go", 40, 0)
	fn := &graph.Node{
		ID: fnID, Kind: graph.NodeFunction, Name: "handleRequest",
		File: "pkg/handler/serve.go", Line: 40, EndLine: 60,
	}
	cpg := makeCPGWithNodes(t, fn)

	report := makeReport("semgrep", "1.56.0",
		makeResult("xss-rule", "error", "XSS detected", "pkg/handler/serve.go", 40),
	)

	result, err := sarif.Ingest(cpg, report, "")
	if err != nil {
		t.Fatalf("Ingest error: %v", err)
	}
	if result.FindingsTotal != 1 {
		t.Errorf("FindingsTotal = %d, want 1", result.FindingsTotal)
	}
	if result.FindingsLinked != 1 {
		t.Errorf("FindingsLinked = %d, want 1", result.FindingsLinked)
	}
	if result.FindingsUnlinked != 0 {
		t.Errorf("FindingsUnlinked = %d, want 0", result.FindingsUnlinked)
	}
	if result.NodesCreated != 1 {
		t.Errorf("NodesCreated = %d, want 1", result.NodesCreated)
	}
	if result.EdgesCreated != 1 {
		t.Errorf("EdgesCreated = %d, want 1", result.EdgesCreated)
	}
	if len(result.ToolNames) != 1 || result.ToolNames[0] != "semgrep" {
		t.Errorf("ToolNames = %v, want [semgrep]", result.ToolNames)
	}

	// Verify the ExternalFinding node was added
	findings := cpg.NodesByKind(graph.NodeExternalFinding)
	if len(findings) != 1 {
		t.Fatalf("ExternalFinding nodes = %d, want 1", len(findings))
	}
	ef := findings[0]
	if ef.RuleID != "xss-rule" {
		t.Errorf("RuleID = %q", ef.RuleID)
	}
	if ef.ToolName != "semgrep" {
		t.Errorf("ToolName = %q", ef.ToolName)
	}
	if ef.Severity != "error" {
		t.Errorf("Severity = %q", ef.Severity)
	}

	// Verify REPORTED_BY edge: function -> finding
	edges := cpg.OutEdges(fnID)
	found := false
	for _, e := range edges {
		if e.Kind == graph.EdgeReportedBy && e.To == ef.ID {
			if e.Confidence != graph.ConfidenceCertain {
				t.Errorf("edge confidence = %q, want CERTAIN", e.Confidence)
			}
			found = true
		}
	}
	if !found {
		t.Error("missing REPORTED_BY edge from function to finding")
	}
}

func TestIngestRangeProximity(t *testing.T) {
	fnID := parser.NodeID(graph.NodeFunction, "processData", "pkg/data.go", 10, 0)
	fn := &graph.Node{
		ID: fnID, Kind: graph.NodeFunction, Name: "processData",
		File: "pkg/data.go", Line: 10, EndLine: 50,
	}
	cpg := makeCPGWithNodes(t, fn)

	report := makeReport("codeql", "2.16.0",
		makeResult("sql-injection", "warning", "SQL injection", "pkg/data.go", 30),
	)

	result, err := sarif.Ingest(cpg, report, "")
	if err != nil {
		t.Fatalf("Ingest error: %v", err)
	}
	if result.FindingsLinked != 1 {
		t.Errorf("FindingsLinked = %d, want 1 (range proximity)", result.FindingsLinked)
	}

	edges := cpg.OutEdges(fnID)
	for _, e := range edges {
		if e.Kind == graph.EdgeReportedBy {
			if e.Confidence != graph.ConfidenceInferred {
				t.Errorf("edge confidence = %q, want INFERRED for range proximity", e.Confidence)
			}
		}
	}
}

func TestIngestUnlinked(t *testing.T) {
	cpg := graph.NewCPG()

	report := makeReport("gosec", "2.19.0",
		makeResult("G101", "warning", "Hardcoded credential", "config/secrets.go", 5),
	)

	result, err := sarif.Ingest(cpg, report, "")
	if err != nil {
		t.Fatalf("Ingest error: %v", err)
	}
	if result.FindingsTotal != 1 {
		t.Errorf("FindingsTotal = %d, want 1", result.FindingsTotal)
	}
	if result.FindingsLinked != 0 {
		t.Errorf("FindingsLinked = %d, want 0", result.FindingsLinked)
	}
	if result.FindingsUnlinked != 1 {
		t.Errorf("FindingsUnlinked = %d, want 1", result.FindingsUnlinked)
	}
	if result.NodesCreated != 1 {
		t.Errorf("NodesCreated = %d, want 1 (unlinked findings still create nodes)", result.NodesCreated)
	}
	if result.EdgesCreated != 0 {
		t.Errorf("EdgesCreated = %d, want 0", result.EdgesCreated)
	}
}

func TestIngestIdempotent(t *testing.T) {
	cpg := graph.NewCPG()

	report := makeReport("semgrep", "1.56.0",
		makeResult("xss", "error", "XSS", "a.go", 10),
	)

	result1, err := sarif.Ingest(cpg, report, "")
	if err != nil {
		t.Fatalf("first Ingest: %v", err)
	}
	if result1.NodesCreated != 1 {
		t.Errorf("first: NodesCreated = %d, want 1", result1.NodesCreated)
	}

	result2, err := sarif.Ingest(cpg, report, "")
	if err != nil {
		t.Fatalf("second Ingest: %v", err)
	}
	if result2.NodesCreated != 0 {
		t.Errorf("second: NodesCreated = %d, want 0 (idempotent)", result2.NodesCreated)
	}
	if result2.FindingsTotal != 1 {
		t.Errorf("second: FindingsTotal = %d, want 1", result2.FindingsTotal)
	}
}

func TestIngestCrossToolDistinct(t *testing.T) {
	cpg := graph.NewCPG()

	r1 := makeReport("semgrep", "1.0", makeResult("rule-1", "error", "msg", "a.go", 10))
	r2 := makeReport("codeql", "2.0", makeResult("rule-1", "error", "msg", "a.go", 10))

	sarif.Ingest(cpg, r1, "")
	sarif.Ingest(cpg, r2, "")

	findings := cpg.NodesByKind(graph.NodeExternalFinding)
	if len(findings) != 2 {
		t.Errorf("ExternalFinding nodes = %d, want 2 (distinct tools)", len(findings))
	}
}

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		uri      string
		repoRoot string
		want     string
	}{
		{"./pkg/foo.go", "", "pkg/foo.go"},
		{"pkg/foo.go", "", "pkg/foo.go"},
		{"file:///home/user/repo/pkg/foo.go", "/home/user/repo", "pkg/foo.go"},
		{"pkg/foo%20bar.go", "", "pkg/foo bar.go"},
		{"./pkg/../pkg/foo.go", "", "pkg/foo.go"},
		{"pkg//foo.go", "", "pkg/foo.go"},
	}
	for _, tt := range tests {
		got := sarif.NormalizePath(tt.uri, tt.repoRoot)
		if got != tt.want {
			t.Errorf("NormalizePath(%q, %q) = %q, want %q", tt.uri, tt.repoRoot, got, tt.want)
		}
	}
}

func TestIngestWithCWEs(t *testing.T) {
	cpg := graph.NewCPG()

	report := &sarif.Report{
		Version: "2.1.0",
		Runs: []sarif.Run{
			{
				Tool: sarif.Tool{
					Driver: sarif.ToolComponent{
						Name: "semgrep",
						Rules: []sarif.Rule{
							{
								ID:         "xss-rule",
								Properties: sarif.RuleProperties{Tags: []string{"CWE-79", "security"}},
							},
						},
					},
				},
				Results: []sarif.Result{
					makeResult("xss-rule", "error", "XSS", "a.go", 10),
				},
			},
		},
	}

	sarif.Ingest(cpg, report, "")

	findings := cpg.NodesByKind(graph.NodeExternalFinding)
	if len(findings) != 1 {
		t.Fatalf("findings = %d, want 1", len(findings))
	}
	if len(findings[0].CWEs) != 1 || findings[0].CWEs[0] != "CWE-79" {
		t.Errorf("CWEs = %v, want [CWE-79]", findings[0].CWEs)
	}
}

func TestIngestResultWithoutLocations(t *testing.T) {
	cpg := graph.NewCPG()

	report := makeReport("test", "1.0",
		sarif.Result{RuleID: "R1", Message: sarif.Message{Text: "msg"}},
	)

	result, err := sarif.Ingest(cpg, report, "")
	if err != nil {
		t.Fatalf("Ingest error: %v", err)
	}
	if result.FindingsTotal != 0 {
		t.Errorf("FindingsTotal = %d, want 0 (no locations)", result.FindingsTotal)
	}
}

func TestIngestMultipleRuns(t *testing.T) {
	cpg := graph.NewCPG()

	report := &sarif.Report{
		Version: "2.1.0",
		Runs: []sarif.Run{
			{
				Tool:    sarif.Tool{Driver: sarif.ToolComponent{Name: "tool1", Version: "1.0"}},
				Results: []sarif.Result{makeResult("R1", "error", "msg1", "a.go", 1)},
			},
			{
				Tool:    sarif.Tool{Driver: sarif.ToolComponent{Name: "tool2", Version: "2.0"}},
				Results: []sarif.Result{makeResult("R2", "warning", "msg2", "b.go", 2)},
			},
		},
	}

	result, err := sarif.Ingest(cpg, report, "")
	if err != nil {
		t.Fatalf("Ingest error: %v", err)
	}
	if result.FindingsTotal != 2 {
		t.Errorf("FindingsTotal = %d, want 2", result.FindingsTotal)
	}
	if result.NodesCreated != 2 {
		t.Errorf("NodesCreated = %d, want 2", result.NodesCreated)
	}
	findings := cpg.NodesByKind(graph.NodeExternalFinding)
	if len(findings) != 2 {
		t.Errorf("ExternalFinding nodes = %d, want 2", len(findings))
	}
	// Verify both tools are tracked
	if len(result.ToolNames) != 2 {
		t.Errorf("ToolNames = %v, want 2 entries", result.ToolNames)
	}
	summary := result.ToolSummary()
	if summary != "tool1 v1.0, tool2 v2.0" {
		t.Errorf("ToolSummary = %q, want 'tool1 v1.0, tool2 v2.0'", summary)
	}
}

// --- Edge Case Tests ---

func TestIngestEmptyReport(t *testing.T) {
	cpg := graph.NewCPG()
	report := &sarif.Report{Version: "2.1.0", Runs: []sarif.Run{}}

	result, err := sarif.Ingest(cpg, report, "")
	if err != nil {
		t.Fatalf("Ingest error: %v", err)
	}
	if result.FindingsTotal != 0 {
		t.Errorf("FindingsTotal = %d, want 0", result.FindingsTotal)
	}
	if result.NodesCreated != 0 {
		t.Errorf("NodesCreated = %d, want 0", result.NodesCreated)
	}
	if len(result.ToolNames) != 0 {
		t.Errorf("ToolNames = %v, want empty", result.ToolNames)
	}
}

func TestIngestRunWithEmptyResults(t *testing.T) {
	cpg := graph.NewCPG()
	report := makeReport("semgrep", "1.0")

	result, err := sarif.Ingest(cpg, report, "")
	if err != nil {
		t.Fatalf("Ingest error: %v", err)
	}
	if result.FindingsTotal != 0 {
		t.Errorf("FindingsTotal = %d, want 0", result.FindingsTotal)
	}
	if len(result.ToolNames) != 1 || result.ToolNames[0] != "semgrep" {
		t.Errorf("ToolNames = %v, want [semgrep] even with empty results", result.ToolNames)
	}
}

func TestIngestMultipleLocationsPerResult(t *testing.T) {
	cpg := graph.NewCPG()

	// One result with two locations (some scanners report this for taint flows)
	report := makeReport("semgrep", "1.0", sarif.Result{
		RuleID:  "taint-flow",
		Level:   "error",
		Message: sarif.Message{Text: "data flows from source to sink"},
		Locations: []sarif.Location{
			{PhysicalLocation: sarif.PhysicalLocation{
				ArtifactLocation: sarif.ArtifactLocation{URI: "src.go"},
				Region:           sarif.Region{StartLine: 10, StartColumn: 1},
			}},
			{PhysicalLocation: sarif.PhysicalLocation{
				ArtifactLocation: sarif.ArtifactLocation{URI: "sink.go"},
				Region:           sarif.Region{StartLine: 20, StartColumn: 1},
			}},
		},
	})

	result, err := sarif.Ingest(cpg, report, "")
	if err != nil {
		t.Fatalf("Ingest error: %v", err)
	}
	// Each location creates a separate finding
	if result.FindingsTotal != 2 {
		t.Errorf("FindingsTotal = %d, want 2 (one per location)", result.FindingsTotal)
	}
	if result.NodesCreated != 2 {
		t.Errorf("NodesCreated = %d, want 2", result.NodesCreated)
	}
}

func TestIngestNestedFunctionRanges(t *testing.T) {
	// Outer function spans lines 1-100, inner function spans lines 40-60
	// Finding at line 50 should link to the inner (tighter) function
	outerID := parser.NodeID(graph.NodeFunction, "outer", "nested.go", 1, 0)
	innerID := parser.NodeID(graph.NodeFunction, "inner", "nested.go", 40, 0)
	outer := &graph.Node{ID: outerID, Kind: graph.NodeFunction, Name: "outer", File: "nested.go", Line: 1, EndLine: 100}
	inner := &graph.Node{ID: innerID, Kind: graph.NodeFunction, Name: "inner", File: "nested.go", Line: 40, EndLine: 60}
	cpg := makeCPGWithNodes(t, outer, inner)

	report := makeReport("scanner", "1.0",
		makeResult("R1", "warning", "issue", "nested.go", 50),
	)

	result, err := sarif.Ingest(cpg, report, "")
	if err != nil {
		t.Fatalf("Ingest error: %v", err)
	}
	if result.FindingsLinked != 1 {
		t.Fatalf("FindingsLinked = %d, want 1", result.FindingsLinked)
	}

	// Should link to inner function (tighter range), not outer
	innerEdges := cpg.OutEdges(innerID)
	outerEdges := cpg.OutEdges(outerID)

	foundInner := false
	for _, e := range innerEdges {
		if e.Kind == graph.EdgeReportedBy {
			foundInner = true
		}
	}
	foundOuter := false
	for _, e := range outerEdges {
		if e.Kind == graph.EdgeReportedBy {
			foundOuter = true
		}
	}

	if !foundInner {
		t.Error("expected REPORTED_BY edge from inner function (tighter range)")
	}
	if foundOuter {
		t.Error("unexpected REPORTED_BY edge from outer function (should use tighter range)")
	}
}

func TestIngestFindingOutsideFunctionRange(t *testing.T) {
	// Function spans lines 10-50, finding at line 60 (outside range)
	fnID := parser.NodeID(graph.NodeFunction, "handler", "out.go", 10, 0)
	fn := &graph.Node{ID: fnID, Kind: graph.NodeFunction, Name: "handler", File: "out.go", Line: 10, EndLine: 50}
	cpg := makeCPGWithNodes(t, fn)

	report := makeReport("scanner", "1.0",
		makeResult("R1", "warning", "issue", "out.go", 60),
	)

	result, err := sarif.Ingest(cpg, report, "")
	if err != nil {
		t.Fatalf("Ingest error: %v", err)
	}
	if result.FindingsLinked != 0 {
		t.Errorf("FindingsLinked = %d, want 0 (finding outside function range)", result.FindingsLinked)
	}
	if result.FindingsUnlinked != 1 {
		t.Errorf("FindingsUnlinked = %d, want 1", result.FindingsUnlinked)
	}
}

func TestIngestDeterministicMatching(t *testing.T) {
	// Two nodes at the same file+line (e.g., function and parameter)
	// Should deterministically pick the same one every time (sorted by ID)
	id1 := parser.NodeID(graph.NodeFunction, "handler", "det.go", 10, 0)
	id2 := parser.NodeID(graph.NodeParameter, "req", "det.go", 10, 5)
	n1 := &graph.Node{ID: id1, Kind: graph.NodeFunction, Name: "handler", File: "det.go", Line: 10}
	n2 := &graph.Node{ID: id2, Kind: graph.NodeParameter, Name: "req", File: "det.go", Line: 10}
	cpg := makeCPGWithNodes(t, n1, n2)

	report := makeReport("scanner", "1.0", makeResult("R1", "error", "issue", "det.go", 10))

	result, err := sarif.Ingest(cpg, report, "")
	if err != nil {
		t.Fatalf("Ingest error: %v", err)
	}
	if result.FindingsLinked != 1 {
		t.Fatalf("FindingsLinked = %d, want 1", result.FindingsLinked)
	}

	// Run again on a fresh CPG with nodes added in reverse order
	cpg2 := makeCPGWithNodes(t, n2, n1)
	_, err = sarif.Ingest(cpg2, report, "")
	if err != nil {
		t.Fatalf("Ingest error: %v", err)
	}

	// Both should pick the same target (sorted by ID)
	findings1 := cpg.NodesByKind(graph.NodeExternalFinding)
	findings2 := cpg2.NodesByKind(graph.NodeExternalFinding)
	if len(findings1) != 1 || len(findings2) != 1 {
		t.Fatal("expected exactly 1 finding in each CPG")
	}

	edges1 := cpg.Edges()
	edges2 := cpg2.Edges()
	var from1, from2 string
	for _, e := range edges1 {
		if e.Kind == graph.EdgeReportedBy {
			from1 = e.From
		}
	}
	for _, e := range edges2 {
		if e.Kind == graph.EdgeReportedBy {
			from2 = e.From
		}
	}
	if from1 != from2 {
		t.Errorf("non-deterministic matching: first run linked from %q, second from %q", from1, from2)
	}
}

func TestIngestFindingAtLine0(t *testing.T) {
	cpg := graph.NewCPG()

	// Some scanners may report line 0 for file-level findings
	report := makeReport("scanner", "1.0",
		makeResult("R1", "warning", "file-level issue", "config.yaml", 0),
	)

	result, err := sarif.Ingest(cpg, report, "")
	if err != nil {
		t.Fatalf("Ingest error: %v", err)
	}
	// Line 0 should still be processed (creates node, stays unlinked)
	if result.FindingsTotal != 1 {
		t.Errorf("FindingsTotal = %d, want 1", result.FindingsTotal)
	}
	if result.NodesCreated != 1 {
		t.Errorf("NodesCreated = %d, want 1", result.NodesCreated)
	}
}

func TestIngestEmptyLevel(t *testing.T) {
	cpg := graph.NewCPG()

	// Level is optional in SARIF - some results omit it
	report := makeReport("scanner", "1.0", sarif.Result{
		RuleID:  "R1",
		Message: sarif.Message{Text: "no level specified"},
		Locations: []sarif.Location{
			{PhysicalLocation: sarif.PhysicalLocation{
				ArtifactLocation: sarif.ArtifactLocation{URI: "a.go"},
				Region:           sarif.Region{StartLine: 5, StartColumn: 1},
			}},
		},
	})

	result, err := sarif.Ingest(cpg, report, "")
	if err != nil {
		t.Fatalf("Ingest error: %v", err)
	}
	if result.NodesCreated != 1 {
		t.Errorf("NodesCreated = %d, want 1 (empty level is valid)", result.NodesCreated)
	}
	findings := cpg.NodesByKind(graph.NodeExternalFinding)
	if findings[0].Severity != "" {
		t.Errorf("Severity = %q, want empty", findings[0].Severity)
	}
}

func TestIngestSameRuleDifferentFiles(t *testing.T) {
	cpg := graph.NewCPG()

	// Same rule reported in two different files should create distinct nodes
	report := makeReport("semgrep", "1.0",
		makeResult("xss", "error", "XSS in file1", "a.go", 10),
		makeResult("xss", "error", "XSS in file2", "b.go", 10),
	)

	result, err := sarif.Ingest(cpg, report, "")
	if err != nil {
		t.Fatalf("Ingest error: %v", err)
	}
	if result.NodesCreated != 2 {
		t.Errorf("NodesCreated = %d, want 2 (same rule, different files)", result.NodesCreated)
	}
}

func TestIngestSameRuleSameFileDifferentLines(t *testing.T) {
	cpg := graph.NewCPG()

	report := makeReport("semgrep", "1.0",
		makeResult("xss", "error", "XSS at line 10", "a.go", 10),
		makeResult("xss", "error", "XSS at line 20", "a.go", 20),
	)

	result, err := sarif.Ingest(cpg, report, "")
	if err != nil {
		t.Fatalf("Ingest error: %v", err)
	}
	if result.NodesCreated != 2 {
		t.Errorf("NodesCreated = %d, want 2 (same rule, same file, different lines)", result.NodesCreated)
	}
}

func TestNormalizePathEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		uri      string
		repoRoot string
		want     string
	}{
		{"empty string", "", "", "."},
		{"just a filename", "main.go", "", "main.go"},
		{"deeply nested", "a/b/c/d/e/f.go", "", "a/b/c/d/e/f.go"},
		{"trailing slash", "pkg/foo/", "", "pkg/foo"},
		{"dot only", ".", "", "."},
		{"multiple parent refs", "a/b/../../c.go", "", "c.go"},
		{"encoded special chars", "pkg/%E4%B8%AD%E6%96%87.go", "", "pkg/\u4e2d\u6587.go"},
		{"absolute path no repo root", "/absolute/path/file.go", "", "/absolute/path/file.go"},
		{"absolute path with repo root", "/repo/root/pkg/file.go", "/repo/root", "pkg/file.go"},
		{"file URI with spaces", "file:///home/user/my%20project/file.go", "/home/user/my project", "file.go"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sarif.NormalizePath(tt.uri, tt.repoRoot)
			if got != tt.want {
				t.Errorf("NormalizePath(%q, %q) = %q, want %q", tt.uri, tt.repoRoot, got, tt.want)
			}
		})
	}
}
