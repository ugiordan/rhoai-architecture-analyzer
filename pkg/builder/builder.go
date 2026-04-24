package builder

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
	"github.com/ugiordan/architecture-analyzer/pkg/parser"
)

// Builder constructs a code property graph from source files in a directory.
type Builder struct {
	parsers []parser.Parser
}

// NewBuilder creates a Builder with the default set of language parsers.
func NewBuilder() *Builder {
	return &Builder{
		parsers: []parser.Parser{
			parser.NewGoParser(),
			parser.NewPythonParser(),
			parser.NewTypeScriptParser(),
			parser.NewRustParser(),
		},
	}
}

// fileEntry holds a discovered source file path and its matched parser index.
type fileEntry struct {
	path     string
	relPath  string
	parserID int
}

// BuildFromDir walks a directory tree, parses supported source files,
// and returns a code property graph with resolved call edges.
func (b *Builder) BuildFromDir(dir string) (*graph.CPG, error) {
	cpg := graph.NewCPG()

	extMap := make(map[string]int)
	for i, p := range b.parsers {
		for _, ext := range p.Extensions() {
			extMap[ext] = i
		}
	}

	// Phase 1: collect file paths (fast, single-threaded walk)
	var files []fileEntry
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("WARN: walk error at %s: %v", path, err)
			return nil
		}
		if d.IsDir() {
			base := filepath.Base(path)
			if base == "vendor" || base == ".git" || base == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		ext := filepath.Ext(path)
		pid, ok := extMap[ext]
		if !ok {
			return nil
		}

		relPath, _ := filepath.Rel(dir, path)
		if relPath == "" {
			relPath = path
		}

		files = append(files, fileEntry{path: path, relPath: relPath, parserID: pid})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking directory: %w", err)
	}

	// Phase 2: parse files in parallel (tree-sitter requires one parser per goroutine)
	workers := runtime.NumCPU()
	if workers > 8 {
		workers = 8
	}
	if workers < 1 {
		workers = 1
	}

	type parseJob struct {
		entry  fileEntry
		result *parser.ParseResult
	}

	results := make([]parseJob, len(files))
	var wg sync.WaitGroup
	ch := make(chan int, len(files))

	for i := range files {
		ch <- i
	}
	close(ch)

	// Shared ID counter across all worker parsers to avoid node ID collisions
	var sharedSeq atomic.Int64

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Each goroutine gets its own parser instance with a shared ID counter
			localParsers := make([]parser.Parser, len(b.parsers))
			for i := range b.parsers {
				switch b.parsers[i].Language() {
				case "go":
					localParsers[i] = parser.NewGoParserWithSeq(&sharedSeq)
				case "python":
					localParsers[i] = parser.NewPythonParserWithSeq(&sharedSeq)
				case "typescript":
					localParsers[i] = parser.NewTypeScriptParserWithSeq(&sharedSeq)
				case "rust":
					localParsers[i] = parser.NewRustParserWithSeq(&sharedSeq)
				}
			}

			for idx := range ch {
				entry := files[idx]
				content, err := os.ReadFile(entry.path)
				if err != nil {
					log.Printf("WARN: failed to read %s: %v", entry.path, err)
					continue
				}

				result, err := localParsers[entry.parserID].ParseFile(entry.relPath, content)
				if err != nil {
					log.Printf("WARN: failed to parse %s: %v", entry.relPath, err)
					continue
				}

				results[idx] = parseJob{entry: entry, result: result}
			}
		}()
	}
	wg.Wait()

	// Phase 3: merge results (single-threaded, fast)
	for _, r := range results {
		if r.result != nil {
			b.mergeResult(cpg, r.result)
		}
	}

	b.resolveCallEdges(cpg)

	return cpg, nil
}

func (b *Builder) mergeResult(cpg *graph.CPG, result *parser.ParseResult) {
	for _, n := range result.Functions {
		cpg.AddNode(n)
	}
	for _, n := range result.CallSites {
		cpg.AddNode(n)
	}
	for _, n := range result.HTTPHandlers {
		cpg.AddNode(n)
	}
	for _, n := range result.DBOperations {
		cpg.AddNode(n)
	}
	for _, n := range result.StructLiterals {
		cpg.AddNode(n)
	}
	for _, e := range result.Edges {
		cpg.AddEdge(e)
	}
}

func (b *Builder) resolveCallEdges(cpg *graph.CPG) {
	callSites := cpg.NodesByKind(graph.NodeCallSite)
	functions := cpg.NodesByKind(graph.NodeFunction)

	// Build function name index keyed by short name
	fnByName := make(map[string][]*graph.Node)
	for _, n := range functions {
		fnByName[n.Name] = append(fnByName[n.Name], n)
	}

	// Build package-qualified index: dir -> funcName -> []*Node
	pkgIndex := make(map[string]map[string][]*graph.Node)
	for _, fn := range functions {
		dir := filepath.Dir(fn.File)
		if pkgIndex[dir] == nil {
			pkgIndex[dir] = make(map[string][]*graph.Node)
		}
		pkgIndex[dir][fn.Name] = append(pkgIndex[dir][fn.Name], fn)
	}

	// Build file-based function index for containment checks
	type fnRange struct {
		node    *graph.Node
		line    int
		endLine int
	}
	fileFns := make(map[string][]fnRange)
	for _, fn := range functions {
		fileFns[fn.File] = append(fileFns[fn.File], fnRange{node: fn, line: fn.Line, endLine: fn.EndLine})
	}

	for _, cs := range callSites {
		callName := cs.Name
		parts := strings.Split(callName, ".")
		shortName := parts[len(parts)-1]
		isQualified := len(parts) > 1

		csDir := filepath.Dir(cs.File)

		var matched []*graph.Node

		if !isQualified {
			// Unqualified call (e.g., "doStuff"): only match within same package
			if pkgFns, ok := pkgIndex[csDir]; ok {
				matched = pkgFns[shortName]
			}
		} else {
			// Qualified call (e.g., "pkg.Func" or "obj.Method"):
			// prefer same-package matches, fall back to cross-package only if none found
			if pkgFns, ok := pkgIndex[csDir]; ok {
				matched = pkgFns[shortName]
			}
			if len(matched) == 0 {
				matched = fnByName[shortName]
			}
		}

		for _, target := range matched {
			confidence := graph.ConfidenceCertain
			targetDir := filepath.Dir(target.File)

			if isQualified && targetDir != csDir {
				// Cross-package match via short name fallback
				confidence = graph.ConfidenceInferred
			}
			if len(matched) > 1 {
				// Multiple candidates: ambiguous resolution
				confidence = graph.ConfidenceUncertain
			}

			cpg.AddEdge(&graph.Edge{
				From:       cs.ID,
				To:         target.ID,
				Kind:       graph.EdgeCalls,
				Label:      callName,
				Confidence: confidence,
			})
		}

		// Find containing function using file index
		for _, fr := range fileFns[cs.File] {
			if cs.Line >= fr.line && cs.Line <= fr.endLine {
				cpg.AddEdge(&graph.Edge{
					From:  fr.node.ID,
					To:    cs.ID,
					Kind:  graph.EdgeDataFlow,
					Label: "contains_call",
				})
				break
			}
		}
	}

	// Link struct literals to containing functions
	structLiterals := cpg.NodesByKind(graph.NodeStructLiteral)
	for _, sl := range structLiterals {
		for _, fr := range fileFns[sl.File] {
			if sl.Line >= fr.line && sl.Line <= fr.endLine {
				cpg.AddEdge(&graph.Edge{
					From:  fr.node.ID,
					To:    sl.ID,
					Kind:  graph.EdgeDataFlow,
					Label: "contains_struct",
				})
				break
			}
		}
	}
}
