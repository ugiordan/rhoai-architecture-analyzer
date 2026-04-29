package sarif

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
	"github.com/ugiordan/architecture-analyzer/pkg/parser"
)

// IngestResult holds statistics from a SARIF ingestion.
type IngestResult struct {
	FindingsTotal    int
	FindingsLinked   int
	FindingsUnlinked int
	NodesCreated     int
	EdgesCreated     int
	ToolNames        []string
	ToolVersions     []string
}

// ToolSummary returns a human-readable summary of tools, e.g. "semgrep v1.56.0, codeql v2.16.0".
func (r *IngestResult) ToolSummary() string {
	var parts []string
	for i, name := range r.ToolNames {
		s := name
		if i < len(r.ToolVersions) && r.ToolVersions[i] != "" {
			s += " v" + r.ToolVersions[i]
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

	ruleMap := make(map[string]map[string]Rule)
	for _, run := range report.Runs {
		toolName := run.Tool.Driver.Name
		if ruleMap[toolName] == nil {
			ruleMap[toolName] = make(map[string]Rule)
		}
		for _, rule := range run.Tool.Driver.Rules {
			ruleMap[toolName][rule.ID] = rule
		}
	}

	for _, run := range report.Runs {
		toolName := run.Tool.Driver.Name
		toolVersion := run.Tool.Driver.Version
		result.ToolNames = append(result.ToolNames, toolName)
		result.ToolVersions = append(result.ToolVersions, toolVersion)

		for _, r := range run.Results {
			if len(r.Locations) == 0 {
				continue
			}

			for _, loc := range r.Locations {
				result.FindingsTotal++

				file := NormalizePath(loc.PhysicalLocation.ArtifactLocation.URI, repoRoot)
				line := loc.PhysicalLocation.Region.StartLine
				col := loc.PhysicalLocation.Region.StartColumn

				var cwes []string
				if rules, ok := ruleMap[toolName]; ok {
					if rule, ok := rules[r.RuleID]; ok {
						cwes = ExtractCWEs(rule.Properties.Tags)
					}
				}

				idInput := fmt.Sprintf("%s/%s", toolName, r.RuleID)
				nodeID := parser.NodeID(graph.NodeExternalFinding, idInput, file, line, col)

				efNode := &graph.Node{
					ID:          nodeID,
					Kind:        graph.NodeExternalFinding,
					Name:        r.RuleID,
					File:        file,
					Line:        line,
					Column:      col,
					RuleID:      r.RuleID,
					Severity:    r.Level,
					Message:     r.Message.Text,
					ToolName:    toolName,
					ToolVersion: toolVersion,
					CWEs:        cwes,
				}

				if err := cpg.AddNode(efNode); err != nil {
					continue
				}
				result.NodesCreated++

				k := fileLineKey{File: file, Line: line}
				if targets, ok := exactIndex[k]; ok && len(targets) > 0 {
					cpg.AddEdge(&graph.Edge{
						From:       targets[0].ID,
						To:         nodeID,
						Kind:       graph.EdgeReportedBy,
						Label:      fmt.Sprintf("%s:%s", toolName, r.RuleID),
						Confidence: graph.ConfidenceCertain,
					})
					result.EdgesCreated++
					result.FindingsLinked++
					continue
				}

				linked := false
				if funcs, ok := funcIndex[file]; ok {
					var bestFunc *graph.Node
					bestSpan := int(^uint(0) >> 1)
					for _, fr := range funcs {
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
							Label:      fmt.Sprintf("%s:%s", toolName, r.RuleID),
							Confidence: graph.ConfidenceInferred,
						})
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

	return decoded
}
