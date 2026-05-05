package sarif

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"sort"
)

const DefaultMaxSARIFSize = 100 * 1024 * 1024 // 100 MB

type Report struct {
	Schema  string `json:"$schema"`
	Version string `json:"version"`
	Runs    []Run  `json:"runs"`
}

type Run struct {
	Tool    Tool     `json:"tool"`
	Results []Result `json:"results"`
}

type Tool struct {
	Driver ToolComponent `json:"driver"`
}

type ToolComponent struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
	Rules   []Rule `json:"rules,omitempty"`
}

type Rule struct {
	ID               string         `json:"id"`
	ShortDescription *Message       `json:"shortDescription,omitempty"`
	Properties       RuleProperties `json:"properties,omitempty"`
}

type RuleProperties struct {
	Tags []string `json:"tags,omitempty"`
}

type Result struct {
	RuleID       string            `json:"ruleId"`
	RuleIndex    int               `json:"ruleIndex,omitempty"`
	Level        string            `json:"level,omitempty"`
	Message      Message           `json:"message"`
	Locations    []Location        `json:"locations,omitempty"`
	Fingerprints map[string]string `json:"fingerprints,omitempty"`
}


type Message struct {
	Text string `json:"text"`
}

type Location struct {
	PhysicalLocation PhysicalLocation `json:"physicalLocation"`
}

type PhysicalLocation struct {
	ArtifactLocation ArtifactLocation `json:"artifactLocation"`
	Region           Region           `json:"region,omitempty"`
}

type ArtifactLocation struct {
	URI string `json:"uri"`
}

type Region struct {
	StartLine   int `json:"startLine,omitempty"`
	StartColumn int `json:"startColumn,omitempty"`
	EndLine     int `json:"endLine,omitempty"`
	EndColumn   int `json:"endColumn,omitempty"`
}

func Parse(r io.Reader) (*Report, error) {
	return ParseWithLimit(r, DefaultMaxSARIFSize)
}

func ParseWithLimit(r io.Reader, maxBytes int64) (*Report, error) {
	// Read into a size-limited buffer first, then unmarshal. This ensures
	// the size limit is enforced before any JSON parsing/allocation.
	lr := io.LimitReader(r, maxBytes+1)
	data, err := io.ReadAll(lr)
	if err != nil {
		return nil, fmt.Errorf("reading SARIF: %w", err)
	}
	if int64(len(data)) > maxBytes {
		return nil, fmt.Errorf("SARIF input exceeds size limit of %d bytes", maxBytes)
	}

	var report Report
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("parsing SARIF JSON: %w", err)
	}

	if report.Version != "2.1.0" {
		return nil, fmt.Errorf("unsupported SARIF version %q: only 2.1.0 is supported", report.Version)
	}

	return &report, nil
}

var cwePattern = regexp.MustCompile(`^CWE-\d+$`)

func ExtractCWEs(tags []string) []string {
	var cwes []string
	for _, tag := range tags {
		if cwePattern.MatchString(tag) {
			cwes = append(cwes, tag)
		}
	}
	sort.Strings(cwes)
	return cwes
}
