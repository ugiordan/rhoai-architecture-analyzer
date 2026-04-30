package sarif

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

var sarifAnnotationSanitizer = regexp.MustCompile(`[^a-zA-Z0-9_-]`)

// sanitizeAnnotationKey replaces characters outside [a-zA-Z0-9_-] with underscores.
func sanitizeAnnotationKey(s string) string {
	return sarifAnnotationSanitizer.ReplaceAllString(s, "_")
}

// DefaultMaxResults is the maximum number of SARIF results to ingest.
// Prevents memory exhaustion from crafted SARIF files within the 100MB size budget.
const DefaultMaxResults = 100_000

// maxFieldLen defines truncation limits for untrusted SARIF string fields.
var maxFieldLen = struct {
	Message     int
	RuleID      int
	ToolName    int
	ToolVersion int
}{
	Message:     4096,
	RuleID:      256,
	ToolName:    128,
	ToolVersion: 64,
}

func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max])
}

// ToolInfo pairs a tool name with its version.
type ToolInfo struct {
	Name    string
	Version string
}

// IngestResult holds statistics from a SARIF ingestion.
type IngestResult struct {
	FindingsTotal    int
	FindingsLinked   int
	FindingsUnlinked int
	NodesCreated     int
	EdgesCreated     int
	Tools            []ToolInfo
}

// ToolSummary returns a human-readable summary of tools, e.g. "semgrep v1.56.0, codeql v2.16.0".
func (r *IngestResult) ToolSummary() string {
	var parts []string
	for _, t := range r.Tools {
		s := t.Name
		if t.Version != "" {
			s += " v" + t.Version
		}
		parts = append(parts, s)
	}
	return strings.Join(parts, ", ")
}

func Ingest(cpg *graph.CPG, report *Report, repoRoot string) (*IngestResult, error) {
	result := &IngestResult{}

	type fileLineKey struct {
		File string
		Line int
	}
	exactIndex := make(map[fileLineKey][]*graph.Node)
	for _, n := range cpg.Nodes() {
		if n.File != "" && n.Line > 0 {
			k := fileLineKey{File: n.File, Line: n.Line}
			exactIndex[k] = append(exactIndex[k], n)
		}
	}
	// Sort targets by ID for deterministic matching when multiple nodes share a file+line
	for k, targets := range exactIndex {
		sort.Slice(targets, func(i, j int) bool { return targets[i].ID < targets[j].ID })
		exactIndex[k] = targets
	}

	type funcRange struct {
		Node    *graph.Node
		EndLine int
	}
	funcIndex := make(map[string][]funcRange)
	for _, n := range cpg.NodesByKind(graph.NodeFunction) {
		if n.File != "" && n.EndLine > 0 {
			funcIndex[n.File] = append(funcIndex[n.File], funcRange{Node: n, EndLine: n.EndLine})
		}
	}
	// PERF-002: Sort per-file function ranges by StartLine for binary search
	for file, funcs := range funcIndex {
		sort.Slice(funcs, func(i, j int) bool { return funcs[i].Node.Line < funcs[j].Node.Line })
		funcIndex[file] = funcs
	}

	totalResults := 0

	for _, run := range report.Runs {
		toolName := truncate(run.Tool.Driver.Name, maxFieldLen.ToolName)
		toolVersion := truncate(run.Tool.Driver.Version, maxFieldLen.ToolVersion)
		result.Tools = append(result.Tools, ToolInfo{Name: toolName, Version: toolVersion})

		// CORR-003: Build rule map per-run to avoid cross-run metadata contamination
		rules := make(map[string]Rule)
		for _, rule := range run.Tool.Driver.Rules {
			rules[rule.ID] = rule
		}

		// PERF-003: Cache annotation keys and edge labels per (toolName, ruleID) pair
		annotationCache := make(map[string]string)
		labelCache := make(map[string]string)

		for _, r := range run.Results {
			if len(r.Locations) == 0 {
				continue
			}

			// SEC-002: Enforce maximum result count at result level.
			// Check before iterating locations so multi-location results
			// are either fully ingested or fully rejected.
			if totalResults+len(r.Locations) > DefaultMaxResults {
				return result, fmt.Errorf("SARIF result count exceeds limit of %d", DefaultMaxResults)
			}

			for _, loc := range r.Locations {
				file := NormalizePath(loc.PhysicalLocation.ArtifactLocation.URI, repoRoot)
				line := loc.PhysicalLocation.Region.StartLine
				col := loc.PhysicalLocation.Region.StartColumn

				var cwes []string
				if rule, ok := rules[r.RuleID]; ok {
					cwes = ExtractCWEs(rule.Properties.Tags)
				}

				// SEC-003: Truncate untrusted string fields
				ruleID := truncate(r.RuleID, maxFieldLen.RuleID)
				message := truncate(r.Message.Text, maxFieldLen.Message)

				idInput := toolName + "/" + ruleID
				nodeID := graph.NodeID(graph.NodeExternalFinding, idInput, file, line, col)

				efNode := &graph.Node{
					ID:          nodeID,
					Kind:        graph.NodeExternalFinding,
					Name:        ruleID,
					File:        file,
					Line:        line,
					Column:      col,
					RuleID:      ruleID,
					Severity:    NormalizedSeverity(r.Level),
					Message:     message,
					ToolName:    toolName,
					ToolVersion: toolVersion,
					CWEs:        cwes,
				}

				// Check for existing node first (idempotent re-ingestion)
				if cpg.GetNode(nodeID) != nil {
					continue
				}
				if err := cpg.AddNode(efNode); err != nil {
					continue
				}
				result.FindingsTotal++
				result.NodesCreated++
				totalResults++

				// PERF-003: Use cached annotation key and edge label
				cacheKey := ruleID
				annotationKey, ok := annotationCache[cacheKey]
				if !ok {
					annotationKey = fmt.Sprintf("sarif:%s:%s",
						sanitizeAnnotationKey(toolName),
						sanitizeAnnotationKey(ruleID))
					annotationCache[cacheKey] = annotationKey
				}
				edgeLabel, ok := labelCache[cacheKey]
				if !ok {
					edgeLabel = toolName + ":" + ruleID
					labelCache[cacheKey] = edgeLabel
				}

				k := fileLineKey{File: file, Line: line}
				if targets, ok := exactIndex[k]; ok && len(targets) > 0 {
					cpg.AddEdge(&graph.Edge{
						From:       targets[0].ID,
						To:         nodeID,
						Kind:       graph.EdgeReportedBy,
						Label:      edgeLabel,
						Confidence: graph.ConfidenceCertain,
					})
					cpg.SetAnnotation(targets[0].ID, annotationKey, true)
					result.EdgesCreated++
					result.FindingsLinked++
					continue
				}

				linked := false
				if funcs, ok := funcIndex[file]; ok {
					// PERF-002: Binary search for the rightmost function with StartLine <= line
					idx := sort.Search(len(funcs), func(i int) bool {
						return funcs[i].Node.Line > line
					})
					// idx is the first function with StartLine > line
					// Walk backward to find enclosing functions and pick the tightest span
					var bestFunc *graph.Node
					bestSpan := int(^uint(0) >> 1)
					for i := idx - 1; i >= 0; i-- {
						fr := funcs[i]
						if fr.Node.Line <= line && line <= fr.EndLine {
							span := fr.EndLine - fr.Node.Line
							if span < bestSpan {
								bestSpan = span
								bestFunc = fr.Node
							}
						}
					}
					if bestFunc != nil {
						cpg.AddEdge(&graph.Edge{
							From:       bestFunc.ID,
							To:         nodeID,
							Kind:       graph.EdgeReportedBy,
							Label:      edgeLabel,
							Confidence: graph.ConfidenceInferred,
						})
						cpg.SetAnnotation(bestFunc.ID, annotationKey, true)
						result.EdgesCreated++
						result.FindingsLinked++
						linked = true
					}
				}

				if !linked {
					result.FindingsUnlinked++
				}
			}
		}
	}

	return result, nil
}

func NormalizePath(uri string, repoRoot string) string {
	decoded, err := url.PathUnescape(uri)
	if err != nil {
		decoded = uri
	}

	if strings.HasPrefix(decoded, "file://") {
		decoded = strings.TrimPrefix(decoded, "file://")
	}

	if repoRoot != "" && filepath.IsAbs(decoded) {
		rel, err := filepath.Rel(repoRoot, decoded)
		if err == nil {
			decoded = rel
		}
	}

	decoded = strings.TrimPrefix(decoded, "./")
	decoded = filepath.ToSlash(decoded)
	decoded = path.Clean(decoded)

	// Reject paths that escape the repository root via parent directory traversal.
	if !filepath.IsAbs(decoded) && strings.HasPrefix(decoded, "..") {
		return filepath.Base(decoded)
	}

	// SEC-001: Strip absolute paths to base name when no repoRoot is set,
	// preventing arbitrary filesystem paths from persisting in CPG nodes.
	if filepath.IsAbs(decoded) && repoRoot == "" {
		return filepath.Base(decoded)
	}

	return decoded
}
