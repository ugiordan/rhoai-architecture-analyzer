package aggregator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// CPGSnapshot holds the serialized code property graph for a single component.
type CPGSnapshot struct {
	Component     string                   `json:"component"`
	SchemaVersion int                      `json:"schema_version"`
	Nodes         []map[string]interface{} `json:"nodes"`
	Edges         []map[string]interface{} `json:"edges"`
}

// PlatformCPG is the merged code property graph across all components.
type PlatformCPG struct {
	Components     []string                   `json:"components"`
	ComponentCount int                        `json:"component_count"`
	TotalNodes     int                        `json:"total_nodes"`
	TotalEdges     int                        `json:"total_edges"`
	CrossEdges     int                        `json:"cross_component_edges"`
	Nodes          []map[string]interface{}   `json:"nodes"`
	Edges          []map[string]interface{}   `json:"edges"`
	CrossLinks     []CrossComponentLink       `json:"cross_component_links"`
}

// Summary returns a lightweight version without full node/edge data,
// suitable for storage in git (the full CPG can be hundreds of MB).
func (p *PlatformCPG) Summary() *PlatformCPGSummary {
	return &PlatformCPGSummary{
		Components:     p.Components,
		ComponentCount: p.ComponentCount,
		TotalNodes:     p.TotalNodes,
		TotalEdges:     p.TotalEdges,
		CrossEdges:     p.CrossEdges,
		CrossLinks:     p.CrossLinks,
	}
}

// PlatformCPGSummary is a lightweight version of PlatformCPG without
// the full node/edge arrays (which can be hundreds of MB).
type PlatformCPGSummary struct {
	Components     []string               `json:"components"`
	ComponentCount int                    `json:"component_count"`
	TotalNodes     int                    `json:"total_nodes"`
	TotalEdges     int                    `json:"total_edges"`
	CrossEdges     int                    `json:"cross_component_edges"`
	CrossLinks     []CrossComponentLink   `json:"cross_component_links"`
}

// CrossComponentLink represents a detected dependency between two components
// at the code level (not just the architecture level).
type CrossComponentLink struct {
	From       string `json:"from_component"`
	To         string `json:"to_component"`
	Type       string `json:"link_type"` // api-call, grpc-call, shared-crd, shared-secret
	FromNode   string `json:"from_node,omitempty"`
	ToNode     string `json:"to_node,omitempty"`
	Evidence   string `json:"evidence,omitempty"`
}

// AggregateCPGs reads code-graph.json files from a results directory and
// merges them into a single platform-wide CPG with cross-component edges.
func AggregateCPGs(resultsDir string) (*PlatformCPG, error) {
	absDir, err := filepath.Abs(resultsDir)
	if err != nil {
		return nil, fmt.Errorf("resolving results dir: %w", err)
	}
	info, err := os.Stat(absDir)
	if err != nil || !info.IsDir() {
		return nil, fmt.Errorf("results directory does not exist: %s", absDir)
	}

	// Find all code-graph.json files
	var snapshots []CPGSnapshot
	err = filepath.Walk(absDir, func(path string, fi os.FileInfo, walkErr error) error {
		if walkErr != nil || fi.IsDir() {
			return nil
		}
		if fi.Name() != "code-graph.json" {
			return nil
		}
		if fi.Size() > 100*1024*1024 { // 100MB cap
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		var snap CPGSnapshot
		if err := json.Unmarshal(data, &snap); err != nil {
			return nil
		}

		// Derive component name from directory
		if snap.Component == "" {
			snap.Component = filepath.Base(filepath.Dir(path))
		}

		snapshots = append(snapshots, snap)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("scanning results: %w", err)
	}

	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].Component < snapshots[j].Component
	})

	return mergeCPGs(snapshots), nil
}

func mergeCPGs(snapshots []CPGSnapshot) *PlatformCPG {
	platform := &PlatformCPG{
		Nodes: make([]map[string]interface{}, 0),
		Edges: make([]map[string]interface{}, 0),
	}

	// Index for cross-component linking
	httpEndpoints := make(map[string]endpointInfo)  // "method:path" -> info
	crdOwners := make(map[string]string)             // CRD kind -> component
	secretUsers := make(map[string][]string)          // secret name -> []component

	for _, snap := range snapshots {
		platform.Components = append(platform.Components, snap.Component)

		// Prefix node IDs to avoid collisions
		prefix := snap.Component + "::"

		for _, node := range snap.Nodes {
			// Clone and prefix
			n := copyMapIface(node)
			if id, ok := n["id"].(string); ok {
				n["id"] = prefix + id
			}
			n["component"] = snap.Component
			platform.Nodes = append(platform.Nodes, n)

			// Index for cross-linking
			kind, _ := n["kind"].(string)
			switch kind {
			case "HTTPEndpoint":
				method, _ := n["method"].(string)
				path, _ := n["path"].(string)
				name, _ := n["name"].(string)
				if path != "" {
					key := strings.ToUpper(method) + ":" + path
					httpEndpoints[key] = endpointInfo{
						component: snap.Component,
						nodeID:    prefix + getString(node, "id", ""),
						name:      name,
					}
				}
			case "CRDDefinition":
				crdKind, _ := n["name"].(string)
				if crdKind != "" {
					crdOwners[crdKind] = snap.Component
				}
			case "SecretRef", "SecretMount", "SecretEnvFrom":
				secretName, _ := n["name"].(string)
				if secretName != "" {
					secretUsers[secretName] = append(secretUsers[secretName], snap.Component)
				}
			}
		}

		for _, edge := range snap.Edges {
			e := copyMapIface(edge)
			if from, ok := e["from"].(string); ok {
				e["from"] = prefix + from
			}
			if to, ok := e["to"].(string); ok {
				e["to"] = prefix + to
			}
			e["component"] = snap.Component
			platform.Edges = append(platform.Edges, e)
		}
	}

	// Detect cross-component links
	var crossLinks []CrossComponentLink

	// 1. HTTP API calls across components
	for _, snap := range snapshots {
		prefix := snap.Component + "::"
		for _, node := range snap.Nodes {
			kind, _ := node["kind"].(string)
			if kind != "ExternalCall" && kind != "CallSite" {
				continue
			}
			target, _ := node["target"].(string)
			if target == "" {
				continue
			}
			// Check if target matches an HTTP endpoint in another component
			for epKey, ep := range httpEndpoints {
				if ep.component == snap.Component {
					continue
				}
				if matchesEndpoint(target, epKey) {
					crossLinks = append(crossLinks, CrossComponentLink{
						From:     snap.Component,
						To:       ep.component,
						Type:     "api-call",
						FromNode: prefix + getString(node, "id", ""),
						ToNode:   ep.nodeID,
						Evidence: fmt.Sprintf("%s calls %s endpoint %s", snap.Component, ep.component, ep.name),
					})
				}
			}
		}
	}

	// 2. Shared CRD watchers (component A defines CRD, component B watches it)
	for _, snap := range snapshots {
		prefix := snap.Component + "::"
		for _, node := range snap.Nodes {
			kind, _ := node["kind"].(string)
			if kind != "ControllerWatch" && kind != "CRDWatch" {
				continue
			}
			watchedKind, _ := node["name"].(string)
			if watchedKind == "" {
				gvk, _ := node["gvk"].(string)
				if idx := strings.LastIndex(gvk, "/"); idx >= 0 {
					watchedKind = gvk[idx+1:]
				} else {
					watchedKind = gvk
				}
			}
			if owner, ok := crdOwners[watchedKind]; ok && owner != snap.Component {
				crossLinks = append(crossLinks, CrossComponentLink{
					From:     snap.Component,
					To:       owner,
					Type:     "shared-crd",
					FromNode: prefix + getString(node, "id", ""),
					Evidence: fmt.Sprintf("%s watches CRD %s owned by %s", snap.Component, watchedKind, owner),
				})
			}
		}
	}

	// 3. Shared secrets
	for secret, users := range secretUsers {
		if len(users) > 1 {
			for i := 0; i < len(users); i++ {
				for j := i + 1; j < len(users); j++ {
					crossLinks = append(crossLinks, CrossComponentLink{
						From:     users[i],
						To:       users[j],
						Type:     "shared-secret",
						Evidence: fmt.Sprintf("both %s and %s reference secret %s", users[i], users[j], secret),
					})
				}
			}
		}
	}

	// Deduplicate cross links (include FromNode and ToNode for fine-grained dedup)
	seen := make(map[string]bool)
	var dedupedLinks []CrossComponentLink
	for _, link := range crossLinks {
		key := link.From + "|" + link.To + "|" + link.Type + "|" + link.FromNode + "|" + link.ToNode
		if !seen[key] {
			seen[key] = true
			dedupedLinks = append(dedupedLinks, link)
		}
	}

	// Add cross-component edges to the merged graph (only for node-level links)
	for _, link := range dedupedLinks {
		if link.FromNode == "" || link.ToNode == "" {
			continue // component-level links (shared-secret, shared-crd) lack node IDs
		}
		platform.Edges = append(platform.Edges, map[string]interface{}{
			"from":           link.FromNode,
			"to":             link.ToNode,
			"kind":           "CrossComponent",
			"link_type":      link.Type,
			"from_component": link.From,
			"to_component":   link.To,
		})
	}

	sort.Slice(dedupedLinks, func(i, j int) bool {
		if dedupedLinks[i].From != dedupedLinks[j].From {
			return dedupedLinks[i].From < dedupedLinks[j].From
		}
		return dedupedLinks[i].To < dedupedLinks[j].To
	})

	platform.CrossLinks = dedupedLinks
	platform.ComponentCount = len(platform.Components)
	platform.TotalNodes = len(platform.Nodes)
	platform.TotalEdges = len(platform.Edges)
	platform.CrossEdges = len(dedupedLinks)

	return platform
}

type endpointInfo struct {
	component string
	nodeID    string
	name      string
}

func matchesEndpoint(target, endpointKey string) bool {
	// endpointKey is "METHOD:/path"
	parts := strings.SplitN(endpointKey, ":", 2)
	if len(parts) != 2 {
		return false
	}
	path := parts[1]
	// Check if the call target references this path
	return strings.Contains(target, path)
}

func copyMapIface(m map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
