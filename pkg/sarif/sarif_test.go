package sarif_test

import (
	"strings"
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/sarif"
)

const validSARIF = `{
  "$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/main/sarif-2.1/schema/sarif-schema-2.1.0.json",
  "version": "2.1.0",
  "runs": [
    {
      "tool": {
        "driver": {
          "name": "semgrep",
          "version": "1.56.0",
          "rules": [
            {
              "id": "go.lang.security.audit.xss.direct-response-writer",
              "shortDescription": {"text": "Direct write to ResponseWriter"},
              "properties": {"tags": ["CWE-79", "security", "OWASP-A7"]}
            }
          ]
        }
      },
      "results": [
        {
          "ruleId": "go.lang.security.audit.xss.direct-response-writer",
          "ruleIndex": 0,
          "level": "error",
          "message": {"text": "Direct write to ResponseWriter detected"},
          "locations": [
            {
              "physicalLocation": {
                "artifactLocation": {"uri": "pkg/handler/serve.go"},
                "region": {"startLine": 45, "startColumn": 3, "endLine": 45, "endColumn": 40}
              }
            }
          ],
          "fingerprints": {"primaryLocationHash": "abc123"}
        }
      ]
    }
  ]
}`

func TestParseValidSARIF(t *testing.T) {
	report, err := sarif.Parse(strings.NewReader(validSARIF))
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if report.Version != "2.1.0" {
		t.Errorf("Version = %q, want 2.1.0", report.Version)
	}
	if len(report.Runs) != 1 {
		t.Fatalf("Runs = %d, want 1", len(report.Runs))
	}
	run := report.Runs[0]
	if run.Tool.Driver.Name != "semgrep" {
		t.Errorf("Tool.Driver.Name = %q, want semgrep", run.Tool.Driver.Name)
	}
	if run.Tool.Driver.Version != "1.56.0" {
		t.Errorf("Tool.Driver.Version = %q, want 1.56.0", run.Tool.Driver.Version)
	}
	if len(run.Tool.Driver.Rules) != 1 {
		t.Fatalf("Rules = %d, want 1", len(run.Tool.Driver.Rules))
	}
	if len(run.Results) != 1 {
		t.Fatalf("Results = %d, want 1", len(run.Results))
	}
	r := run.Results[0]
	if r.RuleID != "go.lang.security.audit.xss.direct-response-writer" {
		t.Errorf("RuleID = %q", r.RuleID)
	}
	if r.Level != "error" {
		t.Errorf("Level = %q, want error", r.Level)
	}
	if len(r.Locations) != 1 {
		t.Fatalf("Locations = %d, want 1", len(r.Locations))
	}
	loc := r.Locations[0]
	if loc.PhysicalLocation.ArtifactLocation.URI != "pkg/handler/serve.go" {
		t.Errorf("URI = %q", loc.PhysicalLocation.ArtifactLocation.URI)
	}
	if loc.PhysicalLocation.Region.StartLine != 45 {
		t.Errorf("StartLine = %d, want 45", loc.PhysicalLocation.Region.StartLine)
	}
}

func TestParseUnsupportedVersion(t *testing.T) {
	doc := `{"version": "2.0.0", "runs": []}`
	_, err := sarif.Parse(strings.NewReader(doc))
	if err == nil {
		t.Fatal("expected error for version 2.0.0")
	}
	if !strings.Contains(err.Error(), "2.1.0") {
		t.Errorf("error should mention 2.1.0, got: %v", err)
	}
}

func TestParseMalformedJSON(t *testing.T) {
	_, err := sarif.Parse(strings.NewReader("{not valid"))
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}

func TestParseEmptyRuns(t *testing.T) {
	doc := `{"version": "2.1.0", "runs": []}`
	report, err := sarif.Parse(strings.NewReader(doc))
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if len(report.Runs) != 0 {
		t.Errorf("Runs = %d, want 0", len(report.Runs))
	}
}

func TestParseWithLimitExceeded(t *testing.T) {
	_, err := sarif.ParseWithLimit(strings.NewReader(validSARIF), 10)
	if err == nil {
		t.Fatal("expected error when input exceeds size limit")
	}
}

func TestExtractCWEs(t *testing.T) {
	tags := []string{"CWE-79", "security", "OWASP-A7", "CWE-89"}
	cwes := sarif.ExtractCWEs(tags)
	if len(cwes) != 2 {
		t.Fatalf("ExtractCWEs = %d CWEs, want 2", len(cwes))
	}
	want := map[string]bool{"CWE-79": true, "CWE-89": true}
	for _, c := range cwes {
		if !want[c] {
			t.Errorf("unexpected CWE: %s", c)
		}
	}
}

func TestExtractCWEsEmpty(t *testing.T) {
	cwes := sarif.ExtractCWEs(nil)
	if len(cwes) != 0 {
		t.Errorf("ExtractCWEs(nil) = %d, want 0", len(cwes))
	}
	cwes = sarif.ExtractCWEs([]string{"security", "OWASP-A7"})
	if len(cwes) != 0 {
		t.Errorf("ExtractCWEs(no CWEs) = %d, want 0", len(cwes))
	}
}

func TestParseResultWithoutLocations(t *testing.T) {
	doc := `{
		"version": "2.1.0",
		"runs": [{
			"tool": {"driver": {"name": "test"}},
			"results": [{"ruleId": "R1", "message": {"text": "msg"}}]
		}]
	}`
	report, err := sarif.Parse(strings.NewReader(doc))
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if len(report.Runs[0].Results[0].Locations) != 0 {
		t.Error("expected empty locations")
	}
}

func TestParseMultipleRuns(t *testing.T) {
	doc := `{
		"version": "2.1.0",
		"runs": [
			{"tool": {"driver": {"name": "tool1"}}, "results": []},
			{"tool": {"driver": {"name": "tool2"}}, "results": []}
		]
	}`
	report, err := sarif.Parse(strings.NewReader(doc))
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if len(report.Runs) != 2 {
		t.Errorf("Runs = %d, want 2", len(report.Runs))
	}
}
