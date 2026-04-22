package query

import "github.com/ugiordan/architecture-analyzer/pkg/graph"

// maxPaths caps the number of taint paths returned to prevent combinatorial explosion.
const maxPaths = 100

// maxTotalVisits caps the total number of node visits across all DFS branches
// to prevent exponential exploration in highly connected graphs.
const maxTotalVisits = 10000

// maxNodeRevisits caps how many times a single node can be revisited across
// different DFS branches. This replaces boolean backtracking to prevent
// exponential re-exploration while still finding multiple distinct paths.
const maxNodeRevisits = 3

func traceToExternalSink(cpg *graph.CPG, startID string, maxDepth int) [][]string {
	var results [][]string
	visitCount := make(map[string]int)
	totalVisits := 0

	var dfs func(nodeID string, path []string, depth int, hasStorageLink bool)
	dfs = func(nodeID string, path []string, depth int, hasStorageLink bool) {
		if len(results) >= maxPaths {
			return
		}
		totalVisits++
		if totalVisits > maxTotalVisits {
			return
		}
		if depth > maxDepth {
			return
		}
		if visitCount[nodeID] >= maxNodeRevisits {
			return
		}
		visitCount[nodeID]++

		currentPath := make([]string, len(path)+1)
		copy(currentPath, path)
		currentPath[len(path)] = nodeID

		node := cpg.GetNode(nodeID)
		if node != nil && node.Annotations != nil && node.Annotations["calls_external"] && hasStorageLink {
			pathCopy := make([]string, len(currentPath))
			copy(pathCopy, currentPath)
			results = append(results, pathCopy)
			return
		}

		for _, edge := range cpg.OutEdges(nodeID) {
			if edge.Kind == graph.EdgeDataFlow || edge.Kind == graph.EdgeStorageLink || edge.Kind == graph.EdgeCalls {
				nextHasStorage := hasStorageLink || edge.Kind == graph.EdgeStorageLink
				dfs(edge.To, currentPath, depth+1, nextHasStorage)
			}
		}
	}

	dfs(startID, nil, 0, false)
	return results
}
