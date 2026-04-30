package dataflow

import (
	"fmt"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

const (
	MaxTaintPaths  = 100   // per-source path limit
	MaxTaintVisits = 10000 // total DFS visit limit
	MaxTaintDepth  = 20    // interprocedural call chain depth limit
)

const maxNodeRevisits = 3

// FlowKind classifies how tainted data flows within a function summary.
type FlowKind string

const (
	FlowKindReturn  FlowKind = "return"
	FlowKindCallArg FlowKind = "call_arg"
	FlowKindSink    FlowKind = "sink"
	FlowKindStorage FlowKind = "storage"
)

// defaultSources returns the default source annotation names.
func defaultSources() []string {
	return []string{
		"handles_user_input",
		"sec:handles_request",
		"sec:deserializes_input",
	}
}

// defaultSinks returns the default sink annotation names.
func defaultSinks() []string {
	return []string{
		"sec:executes_sql",
		"sec:subprocess_call",
		"sec:command_execution",
		"sec:template_render",
		"sec:renders_html",
		"sec:file_access",
		"sec:eval_usage",
		"calls_external",
	}
}

// FlowTarget describes where tainted data flows within a function summary.
type FlowTarget struct {
	Kind     FlowKind // FlowKindReturn, FlowKindCallArg, FlowKindSink, FlowKindStorage
	TargetID string
	ArgIndex int
	Path     []string
}

// TaintSummary captures how a function's parameters propagate taint.
type TaintSummary struct {
	FuncID     string
	ParamFlows map[int][]FlowTarget
}

type taintSource struct {
	nodeID     string
	annotation string
	paramIndex int // -1 if not a parameter
}

// funcContext caches per-function data to avoid redundant computation.
type funcContext struct {
	blocks     []*graph.Node
	nodeToBlock map[string]string
	cfgAdj     map[string][]string
}

// TaintOption configures the TaintEngine.
type TaintOption func(*TaintEngine)

// WithSources sets custom source annotation names.
func WithSources(sources []string) TaintOption {
	return func(te *TaintEngine) { te.sources = sources }
}

// WithSinks sets custom sink annotation names.
func WithSinks(sinks []string) TaintOption {
	return func(te *TaintEngine) { te.sinks = sinks }
}

// WithMaxPaths sets the maximum taint paths per source.
func WithMaxPaths(n int) TaintOption {
	return func(te *TaintEngine) { te.maxPaths = n }
}

// WithMaxVisits sets the total DFS visit limit.
func WithMaxVisits(n int) TaintOption {
	return func(te *TaintEngine) { te.maxVisits = n }
}

// WithMaxDepth sets the interprocedural call chain depth limit.
func WithMaxDepth(n int) TaintOption {
	return func(te *TaintEngine) { te.maxDepth = n }
}

// TaintEngine runs taint propagation analysis over a CPG.
type TaintEngine struct {
	sources   []string
	sinks     []string
	maxPaths  int
	maxVisits int
	maxDepth  int
}

// NewTaintEngine creates an engine with default source/sink annotations and limits.
// Use functional options to customize: NewTaintEngine(WithSources([]string{"custom"})).
func NewTaintEngine(opts ...TaintOption) *TaintEngine {
	te := &TaintEngine{
		sources:   defaultSources(),
		sinks:     defaultSinks(),
		maxPaths:  MaxTaintPaths,
		maxVisits: MaxTaintVisits,
		maxDepth:  MaxTaintDepth,
	}
	for _, opt := range opts {
		opt(te)
	}
	return te
}

// Run executes Phase A (intraprocedural) and Phase B (interprocedural).
func (te *TaintEngine) Run(cpg *graph.CPG) []*graph.Edge {
	// Build per-function context index once, shared by all phases
	funcCtxs := te.buildFuncContexts(cpg)

	edges := te.phaseA(cpg, funcCtxs)
	summaries := te.buildSummaries(cpg, funcCtxs)
	phaseBEdges := te.phaseB(cpg, summaries)
	if phaseBEdges != nil {
		edges = append(edges, phaseBEdges...)
	}
	return edges
}

// buildFuncContexts builds a funcID -> funcContext index from a single scan of all BasicBlocks.
func (te *TaintEngine) buildFuncContexts(cpg *graph.CPG) map[string]*funcContext {
	// Single scan of all basic blocks, grouped by parent function
	blocksByFunc := make(map[string][]*graph.Node)
	for _, b := range cpg.NodesByKind(graph.NodeBasicBlock) {
		blocksByFunc[b.ParentID] = append(blocksByFunc[b.ParentID], b)
	}

	ctxs := make(map[string]*funcContext, len(blocksByFunc))
	for funcID, blocks := range blocksByFunc {
		ntb := make(map[string]string)
		for _, b := range blocks {
			for _, memberID := range b.Members {
				ntb[memberID] = b.ID
			}
		}

		ctxs[funcID] = &funcContext{
			blocks:      blocks,
			nodeToBlock: ntb,
			cfgAdj:      te.buildCFGAdjacency(cpg, blocks),
		}
	}
	return ctxs
}

// buildSummaries constructs TaintSummary for each function by tracing data flow
// from parameters through the function body.
func (te *TaintEngine) buildSummaries(cpg *graph.CPG, funcCtxs map[string]*funcContext) map[string]*TaintSummary {
	summaries := make(map[string]*TaintSummary)

	for _, fn := range cpg.NodesByKind(graph.NodeFunction) {
		ctx, ok := funcCtxs[fn.ID]
		if !ok || len(ctx.blocks) == 0 {
			continue
		}

		summary := &TaintSummary{
			FuncID:     fn.ID,
			ParamFlows: make(map[int][]FlowTarget),
		}

		for _, e := range cpg.EdgesByKindFrom(graph.EdgeDataFlow, fn.ID) {
			if e.Label != "declares" {
				continue
			}
			paramNode := cpg.GetNode(e.To)
			if paramNode == nil || paramNode.Kind != graph.NodeParameter {
				continue
			}
			paramIdx := te.resolveParamIndex(fn, paramNode)
			if paramIdx < 0 {
				continue
			}

			srcBlock, ok := ctx.nodeToBlock[paramNode.ID]
			if !ok {
				continue
			}
			reachable := te.bfsReachable(ctx.cfgAdj, srcBlock)

			targets := te.traceFlowTargets(cpg, fn.ID, paramNode.ID, reachable, ctx.nodeToBlock)
			if len(targets) > 0 {
				summary.ParamFlows[paramIdx] = targets
			}
		}

		if len(summary.ParamFlows) > 0 {
			summaries[fn.ID] = summary
		}
	}
	return summaries
}

// traceFlowTargets traces data flow from a parameter node to find sinks, call args, and returns.
func (te *TaintEngine) traceFlowTargets(
	cpg *graph.CPG,
	funcID string,
	startID string,
	reachable map[string]bool,
	nodeToBlock map[string]string,
) []FlowTarget {
	var targets []FlowTarget
	visited := make(map[string]int)
	totalVisits := 0

	var dfs func(nodeID string, path []string)
	dfs = func(nodeID string, path []string) {
		if totalVisits >= te.maxVisits {
			return
		}
		visited[nodeID]++
		totalVisits++

		node := cpg.GetNode(nodeID)
		if node == nil {
			return
		}

		// Check for sink
		sinkAnn := te.matchSink(cpg, nodeID)
		if sinkAnn != "" {
			pathCopy := make([]string, len(path))
			copy(pathCopy, path)
			targets = append(targets, FlowTarget{
				Kind:     FlowKindSink,
				TargetID: nodeID,
				Path:     pathCopy,
			})
			return
		}

		// Check for storage-linked nodes (DB write operations)
		if len(cpg.EdgesByKindFrom(graph.EdgeStorageLink, nodeID)) > 0 {
			pathCopy := make([]string, len(path))
			copy(pathCopy, path)
			targets = append(targets, FlowTarget{
				Kind:     FlowKindStorage,
				TargetID: nodeID,
				Path:     pathCopy,
			})
		}

		// Check for call site (call_arg target)
		if node.Kind == graph.NodeCallSite {
			argIdx := te.resolveCallSiteArgIndex(cpg, nodeID, startID, path)

			pathCopy := make([]string, len(path))
			copy(pathCopy, path)
			targets = append(targets, FlowTarget{
				Kind:     FlowKindCallArg,
				TargetID: nodeID,
				ArgIndex: argIdx,
				Path:     pathCopy,
			})
		}

		// Check for return to function
		for _, e := range cpg.EdgesByKindFrom(graph.EdgeDataFlow, nodeID) {
			if e.Label == "returns" && e.To == funcID {
				pathCopy := make([]string, len(path))
				copy(pathCopy, path)
				targets = append(targets, FlowTarget{
					Kind:     FlowKindReturn,
					TargetID: funcID,
					Path:     pathCopy,
				})
			}
		}

		// Follow data flow edges
		for _, e := range cpg.EdgesByKindFrom(graph.EdgeDataFlow, nodeID) {
			if !te.isTaintPropagatingLabel(e.Label) {
				continue
			}
			targetBlock, inBlock := nodeToBlock[e.To]
			if inBlock && !reachable[targetBlock] {
				continue
			}
			if visited[e.To] >= maxNodeRevisits {
				continue
			}
			// Explicit path copy to avoid slice aliasing across sibling DFS branches
			newPath := make([]string, len(path)+1)
			copy(newPath, path)
			newPath[len(path)] = e.To
			dfs(e.To, newPath)
		}
	}

	dfs(startID, []string{startID})
	return targets
}

// resolveCallSiteArgIndex determines the argument position for a call site node,
// trying multiple resolution strategies and returning -1 if unresolvable.
func (te *TaintEngine) resolveCallSiteArgIndex(cpg *graph.CPG, callSiteID, startID string, path []string) int {
	// Strategy 1: direct resolution from start node
	argIdx := te.resolveArgIndex(cpg, callSiteID, startID)
	if argIdx >= 0 {
		return argIdx
	}

	// Strategy 2: resolution from path predecessor
	if len(path) >= 2 {
		argIdx = te.resolveArgIndex(cpg, callSiteID, path[len(path)-2])
		if argIdx >= 0 {
			return argIdx
		}
	}

	// Strategy 3: scan in-edges for any path node with passes_to
	inEdges := cpg.InEdges(callSiteID)
	for _, ie := range inEdges {
		if ie.Kind == graph.EdgeDataFlow && ie.Label == "passes_to" {
			for _, pid := range path {
				if ie.From == pid {
					return te.countPassesToBefore(cpg, callSiteID, ie.From)
				}
			}
		}
	}

	// Unresolvable: return -1 to signal callers should skip or handle gracefully
	return -1
}

// countPassesToBefore counts how many passes_to in-edges appear before the given fromID.
func (te *TaintEngine) countPassesToBefore(cpg *graph.CPG, callSiteID, fromID string) int {
	idx := 0
	for _, e := range cpg.InEdges(callSiteID) {
		if e.Kind == graph.EdgeDataFlow && e.Label == "passes_to" {
			if e.From == fromID {
				return idx
			}
			idx++
		}
	}
	return 0
}

// phaseB performs interprocedural taint propagation using function summaries.
func (te *TaintEngine) phaseB(cpg *graph.CPG, summaries map[string]*TaintSummary) []*graph.Edge {
	if len(summaries) == 0 {
		return nil
	}

	// Build call graph index from all call site nodes
	calleeOf := make(map[string]string)                  // call site ID -> target function ID
	callConf := make(map[string]graph.EdgeConfidence)     // call site ID -> edge confidence

	for _, cs := range cpg.NodesByKind(graph.NodeCallSite) {
		for _, ce := range cpg.EdgesByKindFrom(graph.EdgeCalls, cs.ID) {
			calleeOf[cs.ID] = ce.To
			callConf[cs.ID] = ce.Confidence
		}
	}

	var result []*graph.Edge
	totalVisits := 0

	// Walk summaries for interprocedural call_arg propagation
	for _, summary := range summaries {
		fn := cpg.GetNode(summary.FuncID)
		if fn == nil {
			continue
		}

		sourceAnn := te.matchSource(fn)
		if sourceAnn == "" {
			continue
		}

		for paramIdx, flows := range summary.ParamFlows {
			for _, flow := range flows {
				if flow.Kind != FlowKindCallArg {
					continue
				}

				// Track visited (funcID, paramIdx) pairs to detect cycles
				callStack := make(map[string]bool)

				edges := te.followCallChain(
					cpg, summaries, calleeOf, callConf,
					summary.FuncID, paramIdx, sourceAnn,
					flow, []string{}, 0, &totalVisits,
					graph.ConfidenceCertain, callStack,
				)
				result = append(result, edges...)

				if totalVisits >= te.maxVisits {
					return result
				}
			}
		}
	}

	// Follow StorageLink edges
	storageTaint := te.followStorageLinks(cpg, summaries)
	result = append(result, storageTaint...)

	return result
}

// callChainKey uniquely identifies a position in a call chain for cycle detection.
func callChainKey(funcID string, paramIdx int) string {
	return fmt.Sprintf("%s:%d", funcID, paramIdx)
}

// followCallChain recursively traces interprocedural call chains up to maxDepth,
// with cycle detection and call edge confidence propagation.
func (te *TaintEngine) followCallChain(
	cpg *graph.CPG,
	summaries map[string]*TaintSummary,
	calleeOf map[string]string,
	callConf map[string]graph.EdgeConfidence,
	sourceFuncID string,
	sourceParamIdx int,
	sourceAnn string,
	flow FlowTarget,
	prefixPath []string,
	depth int,
	totalVisits *int,
	minConfidence graph.EdgeConfidence,
	callStack map[string]bool,
) []*graph.Edge {
	if *totalVisits >= te.maxVisits {
		return nil
	}

	if depth >= te.maxDepth {
		return nil
	}

	// Build current path segment
	currentPath := make([]string, 0, len(prefixPath)+len(flow.Path))
	currentPath = append(currentPath, prefixPath...)
	for _, p := range flow.Path {
		if len(currentPath) > 0 && currentPath[len(currentPath)-1] == p {
			continue
		}
		currentPath = append(currentPath, p)
	}

	var result []*graph.Edge

	switch flow.Kind {
	case FlowKindSink:
		*totalVisits++

		sourceNodeID := te.resolveSourceNodeID(cpg, sourceFuncID, sourceParamIdx, prefixPath)
		sinkAnn := te.matchSink(cpg, flow.TargetID)
		if sinkAnn == "" {
			sinkAnn = "unknown_sink"
		}

		// Confidence: interprocedural paths are at least INFERRED,
		// downgraded to UNCERTAIN if any call in the chain was not CERTAIN.
		conf := graph.ConfidenceInferred
		if minConfidence == graph.ConfidenceUncertain || depth >= 2 {
			conf = graph.ConfidenceUncertain
		}

		result = append(result, &graph.Edge{
			From:       sourceNodeID,
			To:         flow.TargetID,
			Kind:       graph.EdgeTaint,
			Label:      fmt.Sprintf("%s->%s", sourceAnn, sinkAnn),
			Confidence: conf,
			Path:       currentPath,
		})

	case FlowKindCallArg:
		targetFuncID, ok := calleeOf[flow.TargetID]
		if !ok {
			return nil
		}

		// Cycle detection: skip if we've already visited this (func, param) in the current chain
		key := callChainKey(targetFuncID, flow.ArgIndex)
		if callStack[key] {
			return nil
		}

		targetSummary, ok := summaries[targetFuncID]
		if !ok {
			return nil
		}

		targetFlows, ok := targetSummary.ParamFlows[flow.ArgIndex]
		if !ok {
			return nil
		}

		// Propagate minimum confidence through the chain
		edgeConf := callConf[flow.TargetID]
		chainConf := minConfidence
		if edgeConf == graph.ConfidenceInferred || edgeConf == graph.ConfidenceUncertain {
			chainConf = graph.ConfidenceUncertain
		}

		// Mark this position as visited on the call stack
		callStack[key] = true

		for _, tf := range targetFlows {
			edges := te.followCallChain(
				cpg, summaries, calleeOf, callConf,
				sourceFuncID, sourceParamIdx, sourceAnn,
				tf, currentPath, depth+1, totalVisits,
				chainConf, callStack,
			)
			result = append(result, edges...)

			if *totalVisits >= te.maxVisits {
				delete(callStack, key)
				return result
			}
		}

		delete(callStack, key)

	case FlowKindReturn, FlowKindStorage:
		return nil
	}

	return result
}

// resolveSourceNodeID determines the source node ID for a taint edge.
func (te *TaintEngine) resolveSourceNodeID(cpg *graph.CPG, funcID string, paramIdx int, prefixPath []string) string {
	if len(prefixPath) > 0 {
		return prefixPath[0]
	}
	return te.resolveParamNodeID(cpg, funcID, paramIdx)
}

// matchSource returns the first source annotation found on a function node, or "".
func (te *TaintEngine) matchSource(fn *graph.Node) string {
	for _, s := range te.sources {
		if fn.Annotations[s] {
			return s
		}
	}
	return ""
}

// resolveParamNodeID finds the parameter node ID for a given function and param index.
func (te *TaintEngine) resolveParamNodeID(cpg *graph.CPG, funcID string, paramIdx int) string {
	fn := cpg.GetNode(funcID)
	if fn == nil || paramIdx < 0 || paramIdx >= len(fn.ParamNames) {
		return funcID
	}
	paramName := fn.ParamNames[paramIdx]
	for _, e := range cpg.EdgesByKindFrom(graph.EdgeDataFlow, funcID) {
		if e.Label == "declares" {
			pn := cpg.GetNode(e.To)
			if pn != nil && pn.Kind == graph.NodeParameter && pn.Name == paramName {
				return pn.ID
			}
		}
	}
	return funcID
}

// followStorageLinks traces taint through storage link edges (DB write -> DB read).
func (te *TaintEngine) followStorageLinks(cpg *graph.CPG, summaries map[string]*TaintSummary) []*graph.Edge {
	var result []*graph.Edge

	for _, summary := range summaries {
		fn := cpg.GetNode(summary.FuncID)
		if fn == nil {
			continue
		}

		sourceAnn := te.matchSource(fn)
		if sourceAnn == "" {
			continue
		}

		for paramIdx, flows := range summary.ParamFlows {
			for _, flow := range flows {
				if flow.Kind != FlowKindStorage {
					continue
				}
				if len(flow.Path) == 0 {
					continue
				}
				storageLinks := cpg.EdgesByKindFrom(graph.EdgeStorageLink, flow.TargetID)

				for _, sl := range storageLinks {
					readNodeID := sl.To

					visited := make(map[string]bool)
					sinkEdges := te.storageDFS(cpg, readNodeID, visited, 0)

					for _, sinkID := range sinkEdges {
						sinkAnn := te.matchSink(cpg, sinkID)
						if sinkAnn == "" {
							continue
						}

						sourceNodeID := flow.Path[0]
						if sourceNodeID == summary.FuncID {
							sourceNodeID = te.resolveParamNodeID(cpg, summary.FuncID, paramIdx)
						}

						combinedPath := make([]string, 0, len(flow.Path)+2)
						combinedPath = append(combinedPath, flow.Path...)
						combinedPath = append(combinedPath, readNodeID, sinkID)

						result = append(result, &graph.Edge{
							From:       sourceNodeID,
							To:         sinkID,
							Kind:       graph.EdgeTaint,
							Label:      fmt.Sprintf("%s->%s", sourceAnn, sinkAnn),
							Confidence: graph.ConfidenceUncertain,
							Path:       combinedPath,
						})
					}
				}
			}
		}
	}
	return result
}

// storageDFS follows data flow edges from a node to find all reachable sink nodes,
// bounded by maxDepth to prevent stack overflow on deep linear chains.
func (te *TaintEngine) storageDFS(cpg *graph.CPG, nodeID string, visited map[string]bool, depth int) []string {
	if visited[nodeID] || depth >= te.maxDepth {
		return nil
	}
	visited[nodeID] = true

	var sinks []string

	if te.matchSink(cpg, nodeID) != "" {
		sinks = append(sinks, nodeID)
		return sinks
	}

	for _, e := range cpg.EdgesByKindFrom(graph.EdgeDataFlow, nodeID) {
		if !te.isTaintPropagatingLabel(e.Label) {
			continue
		}
		sinks = append(sinks, te.storageDFS(cpg, e.To, visited, depth+1)...)
	}
	return sinks
}

// phaseA runs intraprocedural taint propagation per function.
func (te *TaintEngine) phaseA(cpg *graph.CPG, funcCtxs map[string]*funcContext) []*graph.Edge {
	var result []*graph.Edge

	for _, fn := range cpg.NodesByKind(graph.NodeFunction) {
		edges := te.analyzeFunction(cpg, fn, funcCtxs)
		result = append(result, edges...)
	}
	return result
}

// analyzeFunction performs intraprocedural taint analysis on a single function.
func (te *TaintEngine) analyzeFunction(cpg *graph.CPG, fn *graph.Node, funcCtxs map[string]*funcContext) []*graph.Edge {
	ctx, ok := funcCtxs[fn.ID]
	if !ok || len(ctx.blocks) == 0 {
		return nil
	}

	sources := te.findSources(cpg, fn, ctx.blocks)
	if len(sources) == 0 {
		return nil
	}

	var result []*graph.Edge
	for _, src := range sources {
		srcBlock, ok := ctx.nodeToBlock[src.nodeID]
		if !ok {
			continue
		}
		reachable := te.bfsReachable(ctx.cfgAdj, srcBlock)

		edges := te.dfsPropagate(cpg, src, reachable, ctx.nodeToBlock)
		result = append(result, edges...)
	}
	return result
}

// buildCFGAdjacency builds a directed adjacency map from EdgeControlFlow edges between blocks.
func (te *TaintEngine) buildCFGAdjacency(cpg *graph.CPG, blocks []*graph.Node) map[string][]string {
	blockSet := make(map[string]bool, len(blocks))
	for _, b := range blocks {
		blockSet[b.ID] = true
	}

	adj := make(map[string][]string)
	for _, b := range blocks {
		for _, e := range cpg.EdgesByKindFrom(graph.EdgeControlFlow, b.ID) {
			if blockSet[e.To] {
				adj[b.ID] = append(adj[b.ID], e.To)
			}
		}
	}
	return adj
}

// findSources identifies taint origins in a function.
// Accepts pre-computed blocks to avoid redundant scans.
func (te *TaintEngine) findSources(cpg *graph.CPG, fn *graph.Node, blocks []*graph.Node) []taintSource {
	var sources []taintSource
	seen := make(map[string]bool) // deduplicate by nodeID

	// Check if function itself has a source annotation
	for _, ann := range te.sources {
		if fn.Annotations[ann] {
			// Function-level source: parameters become taint origins
			for _, e := range cpg.EdgesByKindFrom(graph.EdgeDataFlow, fn.ID) {
				if e.Label == "declares" {
					paramNode := cpg.GetNode(e.To)
					if paramNode == nil {
						continue
					}
					if seen[paramNode.ID] {
						continue
					}
					seen[paramNode.ID] = true
					paramIdx := te.resolveParamIndex(fn, paramNode)
					sources = append(sources, taintSource{
						nodeID:     paramNode.ID,
						annotation: ann,
						paramIndex: paramIdx,
					})
				}
			}
			break // one matching annotation is enough for function-level sources
		}
	}

	// Check individual nodes in function's blocks for source annotations
	for _, b := range blocks {
		for _, memberID := range b.Members {
			if seen[memberID] {
				continue
			}
			node := cpg.GetNode(memberID)
			if node == nil || node.Kind == graph.NodeParameter {
				continue
			}
			for _, ann := range te.sources {
				if node.Annotations[ann] {
					seen[memberID] = true
					sources = append(sources, taintSource{
						nodeID:     node.ID,
						annotation: ann,
						paramIndex: -1,
					})
					break // one annotation per node is enough
				}
			}
		}
	}

	return sources
}

// resolveParamIndex finds the parameter's position among the function's declared params.
func (te *TaintEngine) resolveParamIndex(fn *graph.Node, param *graph.Node) int {
	for i, name := range fn.ParamNames {
		if name == param.Name {
			return i
		}
	}
	return -1
}

// bfsReachable computes the set of block IDs reachable from startBlock via CFG edges.
func (te *TaintEngine) bfsReachable(adj map[string][]string, startBlock string) map[string]bool {
	reachable := map[string]bool{startBlock: true}
	queue := []string{startBlock}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, next := range adj[cur] {
			if !reachable[next] {
				reachable[next] = true
				queue = append(queue, next)
			}
		}
	}
	return reachable
}

// dfsPropagate follows data flow edges from a taint source, constrained by CFG reachability.
func (te *TaintEngine) dfsPropagate(
	cpg *graph.CPG,
	src taintSource,
	reachable map[string]bool,
	nodeToBlock map[string]string,
) []*graph.Edge {
	var result []*graph.Edge

	visited := make(map[string]int) // nodeID -> visit count
	totalVisits := 0
	pathCount := 0

	var dfs func(nodeID string, path []string)
	dfs = func(nodeID string, path []string) {
		if totalVisits >= te.maxVisits || pathCount >= te.maxPaths {
			return
		}

		visited[nodeID]++
		totalVisits++

		// Check if current node is a sink
		sinkAnn := te.matchSink(cpg, nodeID)
		if sinkAnn != "" {
			pathCount++
			pathCopy := make([]string, len(path))
			copy(pathCopy, path)
			result = append(result, &graph.Edge{
				From:       src.nodeID,
				To:         nodeID,
				Kind:       graph.EdgeTaint,
				Label:      fmt.Sprintf("%s->%s", src.annotation, sinkAnn),
				Confidence: graph.ConfidenceCertain,
				Path:       pathCopy,
			})
			return
		}

		// Follow data flow edges
		for _, e := range cpg.EdgesByKindFrom(graph.EdgeDataFlow, nodeID) {
			if !te.isTaintPropagatingLabel(e.Label) {
				continue
			}

			targetID := e.To
			targetBlock, inBlock := nodeToBlock[targetID]
			if inBlock && !reachable[targetBlock] {
				continue
			}

			if visited[targetID] >= maxNodeRevisits {
				continue
			}

			// Explicit path copy to avoid slice aliasing across sibling DFS branches
			newPath := make([]string, len(path)+1)
			copy(newPath, path)
			newPath[len(path)] = targetID
			dfs(targetID, newPath)
		}
	}

	dfs(src.nodeID, []string{src.nodeID})
	return result
}

// isTaintPropagatingLabel returns true for data flow labels that propagate taint.
func (te *TaintEngine) isTaintPropagatingLabel(label string) bool {
	switch label {
	case "assigns", "reads", "field_access", "mutates", "passes_to", "returns":
		return true
	}
	return false
}

// matchSink checks if a node has any sink annotation and returns the first match.
func (te *TaintEngine) matchSink(cpg *graph.CPG, nodeID string) string {
	node := cpg.GetNode(nodeID)
	if node == nil {
		return ""
	}
	for _, s := range te.sinks {
		if node.Annotations[s] {
			return s
		}
	}
	return ""
}

// resolveArgIndex determines argument position by counting passes_to InEdges to the call site.
func (te *TaintEngine) resolveArgIndex(cpg *graph.CPG, callSiteID string, argNodeID string) int {
	idx := 0
	for _, e := range cpg.InEdges(callSiteID) {
		if e.Kind == graph.EdgeDataFlow && e.Label == "passes_to" {
			if e.From == argNodeID {
				return idx
			}
			idx++
		}
	}
	return -1
}
