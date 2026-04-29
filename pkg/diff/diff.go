package diff

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

// GraphSnapshot represents a serialized code-graph.json file.
type GraphSnapshot struct {
	SchemaVersion int          `json:"schema_version"`
	Nodes         []graph.Node `json:"nodes"`
	Edges         []graph.Edge `json:"edges"`
}

// GraphDiff is the result of comparing two code-graph.json snapshots.
type GraphDiff struct {
	SchemaVersion int         `json:"schema_version"`
	BaseVersion   string      `json:"base_version,omitempty"`
	HeadVersion   string      `json:"head_version,omitempty"`
	Nodes         NodeDiff    `json:"nodes"`
	Edges         EdgeDiff    `json:"edges"`
	Summary       DiffSummary `json:"summary"`
}

// NodeDiff holds added, removed, and modified nodes.
type NodeDiff struct {
	Added    []graph.Node `json:"added"`
	Removed  []graph.Node `json:"removed"`
	Modified []NodeChange `json:"modified"`
}

// NodeChange describes a node that exists in both snapshots but has field differences.
type NodeChange struct {
	ID      string        `json:"id"`
	Before  graph.Node    `json:"before"`
	After   graph.Node    `json:"after"`
	Changes []FieldChange `json:"changes"`
}

// FieldChange describes a single field difference between two versions of a node.
type FieldChange struct {
	Field    string `json:"field"`
	OldValue any    `json:"old_value"`
	NewValue any    `json:"new_value"`
}

// EdgeDiff holds added and removed edges.
type EdgeDiff struct {
	Added   []graph.Edge `json:"added"`
	Removed []graph.Edge `json:"removed"`
}

// DiffSummary provides aggregate counts of changes.
type DiffSummary struct {
	NodesAdded    int                   `json:"nodes_added"`
	NodesRemoved  int                   `json:"nodes_removed"`
	NodesModified int                   `json:"nodes_modified"`
	EdgesAdded    int                   `json:"edges_added"`
	EdgesRemoved  int                   `json:"edges_removed"`
	ByKind        map[string]KindCounts `json:"by_kind"`
	ByLanguage    map[string]KindCounts `json:"by_language"`
}

// KindCounts tracks add/remove/modify counts for a single category.
type KindCounts struct {
	Added    int `json:"added"`
	Removed  int `json:"removed"`
	Modified int `json:"modified"`
}

// Compare produces a GraphDiff from two snapshots.
func Compare(base, head GraphSnapshot) (*GraphDiff, error) {
	if base.SchemaVersion < 3 {
		return nil, fmt.Errorf("base snapshot has schema_version %d, but stable IDs require schema_version >= 3", base.SchemaVersion)
	}
	if head.SchemaVersion < 3 {
		return nil, fmt.Errorf("head snapshot has schema_version %d, but stable IDs require schema_version >= 3", head.SchemaVersion)
	}

	baseNodes, err := buildNodeMap(base.Nodes)
	if err != nil {
		return nil, fmt.Errorf("base snapshot: %w", err)
	}
	headNodes, err := buildNodeMap(head.Nodes)
	if err != nil {
		return nil, fmt.Errorf("head snapshot: %w", err)
	}

	d := &GraphDiff{
		SchemaVersion: 1,
		Nodes: NodeDiff{
			Added:    make([]graph.Node, 0),
			Removed:  make([]graph.Node, 0),
			Modified: make([]NodeChange, 0),
		},
		Edges: EdgeDiff{
			Added:   make([]graph.Edge, 0),
			Removed: make([]graph.Edge, 0),
		},
		Summary: DiffSummary{
			ByKind:     make(map[string]KindCounts),
			ByLanguage: make(map[string]KindCounts),
		},
	}

	// Added nodes: in head but not base
	for id, node := range headNodes {
		if _, ok := baseNodes[id]; !ok {
			d.Nodes.Added = append(d.Nodes.Added, node)
		}
	}

	// Removed nodes: in base but not head
	for id, node := range baseNodes {
		if _, ok := headNodes[id]; !ok {
			d.Nodes.Removed = append(d.Nodes.Removed, node)
		}
	}

	// Modified nodes: in both, compare fields
	for id, baseNode := range baseNodes {
		headNode, ok := headNodes[id]
		if !ok {
			continue
		}
		changes := compareNodes(baseNode, headNode)
		if len(changes) > 0 {
			d.Nodes.Modified = append(d.Nodes.Modified, NodeChange{
				ID:      id,
				Before:  baseNode,
				After:   headNode,
				Changes: changes,
			})
		}
	}

	// Edge diff
	baseEdgeSet := buildEdgeSet(base.Edges)
	headEdgeSet := buildEdgeSet(head.Edges)

	for key, edge := range headEdgeSet {
		if _, ok := baseEdgeSet[key]; !ok {
			d.Edges.Added = append(d.Edges.Added, edge)
		}
	}
	for key, edge := range baseEdgeSet {
		if _, ok := headEdgeSet[key]; !ok {
			d.Edges.Removed = append(d.Edges.Removed, edge)
		}
	}

	// Sort for deterministic output
	sort.Slice(d.Nodes.Added, func(i, j int) bool { return d.Nodes.Added[i].ID < d.Nodes.Added[j].ID })
	sort.Slice(d.Nodes.Removed, func(i, j int) bool { return d.Nodes.Removed[i].ID < d.Nodes.Removed[j].ID })
	sort.Slice(d.Nodes.Modified, func(i, j int) bool { return d.Nodes.Modified[i].ID < d.Nodes.Modified[j].ID })
	sortEdges := func(edges []graph.Edge) {
		sort.Slice(edges, func(i, j int) bool {
			if edges[i].From != edges[j].From {
				return edges[i].From < edges[j].From
			}
			if edges[i].To != edges[j].To {
				return edges[i].To < edges[j].To
			}
			return edges[i].Label < edges[j].Label
		})
	}
	sortEdges(d.Edges.Added)
	sortEdges(d.Edges.Removed)

	// Build summary
	d.Summary.NodesAdded = len(d.Nodes.Added)
	d.Summary.NodesRemoved = len(d.Nodes.Removed)
	d.Summary.NodesModified = len(d.Nodes.Modified)
	d.Summary.EdgesAdded = len(d.Edges.Added)
	d.Summary.EdgesRemoved = len(d.Edges.Removed)

	for _, n := range d.Nodes.Added {
		incKind(d.Summary.ByKind, string(n.Kind), "added")
		incLang(d.Summary.ByLanguage, n.Language, "added")
	}
	for _, n := range d.Nodes.Removed {
		incKind(d.Summary.ByKind, string(n.Kind), "removed")
		incLang(d.Summary.ByLanguage, n.Language, "removed")
	}
	for _, mc := range d.Nodes.Modified {
		incKind(d.Summary.ByKind, string(mc.After.Kind), "modified")
		incLang(d.Summary.ByLanguage, mc.After.Language, "modified")
	}

	return d, nil
}

// HasDifferences returns true if the diff contains any changes.
func (d *GraphDiff) HasDifferences() bool {
	return d.Summary.NodesAdded > 0 || d.Summary.NodesRemoved > 0 ||
		d.Summary.NodesModified > 0 || d.Summary.EdgesAdded > 0 ||
		d.Summary.EdgesRemoved > 0
}

func buildNodeMap(nodes []graph.Node) (map[string]graph.Node, error) {
	m := make(map[string]graph.Node, len(nodes))
	for _, n := range nodes {
		if _, exists := m[n.ID]; exists {
			return nil, fmt.Errorf("duplicate node ID %q", n.ID)
		}
		m[n.ID] = n
	}
	return m, nil
}

type edgeKey struct {
	From  string
	To    string
	Kind  graph.EdgeKind
	Label string
}

func buildEdgeSet(edges []graph.Edge) map[edgeKey]graph.Edge {
	m := make(map[edgeKey]graph.Edge, len(edges))
	for _, e := range edges {
		k := edgeKey{From: e.From, To: e.To, Kind: e.Kind, Label: e.Label}
		m[k] = e
	}
	return m
}

func compareNodes(base, head graph.Node) []FieldChange {
	var changes []FieldChange

	addChange := func(field string, old, new any) {
		changes = append(changes, FieldChange{Field: field, OldValue: old, NewValue: new})
	}

	if base.EndLine != head.EndLine {
		addChange("end_line", base.EndLine, head.EndLine)
	}
	if base.Language != head.Language {
		addChange("language", base.Language, head.Language)
	}
	if base.TypeName != head.TypeName {
		addChange("type_name", base.TypeName, head.TypeName)
	}
	if base.Complexity != head.Complexity {
		addChange("complexity", base.Complexity, head.Complexity)
	}
	if base.ReturnType != head.ReturnType {
		addChange("return_type", base.ReturnType, head.ReturnType)
	}
	if base.IsTest != head.IsTest {
		addChange("is_test", base.IsTest, head.IsTest)
	}
	if base.IsUnsafe != head.IsUnsafe {
		addChange("is_unsafe", base.IsUnsafe, head.IsUnsafe)
	}
	if base.IsExtern != head.IsExtern {
		addChange("is_extern", base.IsExtern, head.IsExtern)
	}
	if base.CallTarget != head.CallTarget {
		addChange("call_target", base.CallTarget, head.CallTarget)
	}
	if base.IsMacro != head.IsMacro {
		addChange("is_macro", base.IsMacro, head.IsMacro)
	}
	if base.Route != head.Route {
		addChange("route", base.Route, head.Route)
	}
	if base.HTTPMethod != head.HTTPMethod {
		addChange("http_method", base.HTTPMethod, head.HTTPMethod)
	}
	if base.Operation != head.Operation {
		addChange("operation", base.Operation, head.Operation)
	}
	if base.Table != head.Table {
		addChange("table", base.Table, head.Table)
	}
	if base.StructType != head.StructType {
		addChange("struct_type", base.StructType, head.StructType)
	}
	if base.TrustLevel != head.TrustLevel {
		addChange("trust_level", base.TrustLevel, head.TrustLevel)
	}

	// External finding fields
	if base.RuleID != head.RuleID {
		addChange("rule_id", base.RuleID, head.RuleID)
	}
	if base.Severity != head.Severity {
		addChange("severity", base.Severity, head.Severity)
	}
	if base.Message != head.Message {
		addChange("message", base.Message, head.Message)
	}
	if base.ToolName != head.ToolName {
		addChange("tool_name", base.ToolName, head.ToolName)
	}
	if base.ToolVersion != head.ToolVersion {
		addChange("tool_version", base.ToolVersion, head.ToolVersion)
	}
	if !sortedEqual(base.CWEs, head.CWEs) {
		addChange("cwes", base.CWEs, head.CWEs)
	}

	// Slice comparisons (order-independent)
	if !sortedEqual(base.Decorators, head.Decorators) {
		addChange("decorators", base.Decorators, head.Decorators)
	}
	if !sortedEqual(base.ParamNames, head.ParamNames) {
		addChange("param_names", base.ParamNames, head.ParamNames)
	}
	if !sortedEqual(base.ParamTypes, head.ParamTypes) {
		addChange("param_types", base.ParamTypes, head.ParamTypes)
	}
	if !sortedEqual(base.FieldNames, head.FieldNames) {
		addChange("field_names", base.FieldNames, head.FieldNames)
	}

	// Map comparisons
	if !reflect.DeepEqual(base.Annotations, head.Annotations) {
		addChange("annotations", base.Annotations, head.Annotations)
	}
	if !reflect.DeepEqual(base.Properties, head.Properties) {
		addChange("properties", base.Properties, head.Properties)
	}

	return changes
}

func sortedEqual(a, b []string) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	if len(a) != len(b) {
		return false
	}
	aCopy := make([]string, len(a))
	bCopy := make([]string, len(b))
	copy(aCopy, a)
	copy(bCopy, b)
	sort.Strings(aCopy)
	sort.Strings(bCopy)
	for i := range aCopy {
		if aCopy[i] != bCopy[i] {
			return false
		}
	}
	return true
}

func incKind(m map[string]KindCounts, kind, op string) {
	if kind == "" {
		return
	}
	c := m[kind]
	switch op {
	case "added":
		c.Added++
	case "removed":
		c.Removed++
	case "modified":
		c.Modified++
	}
	m[kind] = c
}

func incLang(m map[string]KindCounts, lang, op string) {
	if lang == "" {
		return
	}
	c := m[lang]
	switch op {
	case "added":
		c.Added++
	case "removed":
		c.Removed++
	case "modified":
		c.Modified++
	}
	m[lang] = c
}
