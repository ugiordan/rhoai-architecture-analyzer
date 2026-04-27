// Package main implements the arch-analyzer CLI. As the tool grows, consider
// splitting subcommands into dedicated packages under internal/cmd/.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/ugiordan/architecture-analyzer/pkg/aggregator"
	"github.com/ugiordan/architecture-analyzer/pkg/annotator"
	"github.com/ugiordan/architecture-analyzer/pkg/arch"
	"github.com/ugiordan/architecture-analyzer/pkg/builder"
	"github.com/ugiordan/architecture-analyzer/pkg/domains"
	"github.com/ugiordan/architecture-analyzer/pkg/domains/security"
	testingdomain "github.com/ugiordan/architecture-analyzer/pkg/domains/testing"
	"github.com/ugiordan/architecture-analyzer/pkg/domains/upgrade"
	"github.com/ugiordan/architecture-analyzer/pkg/extractor"
	"github.com/ugiordan/architecture-analyzer/pkg/graph"
	"github.com/ugiordan/architecture-analyzer/pkg/linker"
	"github.com/ugiordan/architecture-analyzer/pkg/query"
	"github.com/ugiordan/architecture-analyzer/pkg/renderer"
	"github.com/ugiordan/architecture-analyzer/pkg/validator"
)

const version = "0.2.0"

func init() {
	extractor.AnalyzerVersion = version
	domains.Register(security.New())
	domains.Register(testingdomain.New())
	domains.Register(upgrade.New())
}

// versionLabelRe validates version labels for snapshot output directories.
var versionLabelRe = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{0,63}$`)

// validateVersionLabel checks that a version label is safe for use in file paths and git tags.
func validateVersionLabel(label string) error {
	if !versionLabelRe.MatchString(label) {
		return fmt.Errorf("invalid version label %q: must match %s", label, versionLabelRe.String())
	}
	return nil
}

// applyVersion adjusts outputDir to include a version subdirectory when version is non-empty.
// e.g. applyVersion("output", "v2.15.0") returns "output/v2.15.0".
func applyVersion(outputDir, ver string) string {
	if ver == "" {
		return outputDir
	}
	return filepath.Join(outputDir, ver)
}

// snapshotMetadata holds version and provenance information written alongside output.
type snapshotMetadata struct {
	Version         string            `json:"version"`
	Timestamp       string            `json:"timestamp"`
	AnalyzerVersion string            `json:"analyzer_version"`
	ReposAnalyzed   map[string]string `json:"repos_analyzed"`
	Platform        string            `json:"platform,omitempty"`
}

// writeSnapshotMetadata writes snapshot-metadata.json to the given directory.
func writeSnapshotMetadata(dir, ver string, repos map[string]string, platform string) error {
	meta := snapshotMetadata{
		Version:         ver,
		Timestamp:       time.Now().UTC().Format(time.RFC3339),
		AnalyzerVersion: version,
		ReposAnalyzed:   repos,
		Platform:        platform,
	}
	path := filepath.Join(dir, "snapshot-metadata.json")
	return writeJSON(path, meta)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	var err error
	switch cmd {
	case "extract":
		err = cmdExtract(args)
	case "render":
		err = cmdRender(args)
	case "analyze":
		err = cmdAnalyze(args)
	case "aggregate":
		err = cmdAggregate(args)
	case "extract-schema":
		err = cmdExtractSchema(args)
	case "validate":
		err = cmdValidate(args)
	case "scan":
		err = cmdScan(args)
	case "graph":
		err = cmdGraph(args)
	case "domains":
		err = cmdDomains()
	case "docs":
		err = cmdDocs(args)
	case "full-analysis":
		err = cmdFullAnalysis(args)
	case "version":
		fmt.Printf("arch-analyzer %s\n", version)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmd)
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`arch-analyzer: Architecture Analyzer and Code Graph Security Scanner

Usage: arch-analyzer <command> [options]

Architecture commands:
  extract <repo-path>                  Extract architecture data from a repository
  render <json-file>                   Render diagrams from architecture JSON
  analyze <repo-path>                  Extract + render in one step
  aggregate <results-dir>              Aggregate multiple component JSONs into platform view
  docs <json-file>                     Generate browsable documentation site from architecture JSON

Contract validation commands:
  extract-schema <repo-path>           Extract CRD JSON schemas from a repository
  validate <repo-path>                 Validate CRD changes against stored contracts

Code graph commands:
  scan <repo-path>                     Build code graph and run security queries
                                       [--domains d1,d2] [--with-arch]
  graph <repo-path>                    Export code property graph (JSON or DOT)
  domains                              List registered analysis domains

Combined:
  full-analysis <repo-path>            Run architecture extraction + code graph scan

Other:
  version                              Print version
  help                                 Show this help`)
}

// cmdExtract extracts architecture data from a repo and writes JSON.
func cmdExtract(args []string) error {
	fs := flag.NewFlagSet("extract", flag.ExitOnError)
	output := fs.String("output", "component-architecture.json", "Output JSON file")
	org := fs.String("org", "", "GitHub organization (auto-detected from go.mod if empty)")
	ver := fs.String("version", "", "Version label for snapshot output (e.g. v2.15.0)")
	fs.Parse(args)

	if fs.NArg() < 1 {
		return fmt.Errorf("usage: arch-analyzer extract <repo-path> [--output file.json] [--org org] [--version label]")
	}

	if *ver != "" {
		if err := validateVersionLabel(*ver); err != nil {
			return err
		}
	}

	opts := &extractor.ExtractOptions{Org: *org}
	arch, err := extractor.ExtractAll(fs.Arg(0), opts)
	if err != nil {
		return err
	}

	outPath := *output
	if *ver != "" {
		// Place output in a versioned subdirectory
		dir := filepath.Dir(outPath)
		base := filepath.Base(outPath)
		outPath = filepath.Join(dir, *ver, base)
	}
	return writeJSON(outPath, arch)
}

// cmdRender renders diagrams from architecture JSON.
func cmdRender(args []string) error {
	fs := flag.NewFlagSet("render", flag.ExitOnError)
	outputDir := fs.String("output-dir", "", "Output directory (default: <json-dir>/diagrams)")
	formats := fs.String("formats", "", "Comma-separated formats: rbac,component,security_network,dependencies,c4,dataflow (default: all)")
	fs.Parse(args)

	if fs.NArg() < 1 {
		return fmt.Errorf("usage: arch-analyzer render <json-file> [--output-dir dir] [--formats fmt1,fmt2]")
	}

	jsonPath := fs.Arg(0)
	data, err := loadJSON(jsonPath)
	if err != nil {
		return err
	}

	outDir := *outputDir
	if outDir == "" {
		outDir = filepath.Join(filepath.Dir(jsonPath), "diagrams")
	}

	var fmts []string
	if *formats != "" {
		fmts = strings.Split(*formats, ",")
	}

	diagrams := renderer.RenderAll(data, fmts)
	return writeDiagrams(outDir, diagrams)
}

// cmdAnalyze runs extract + render in one step.
func cmdAnalyze(args []string) error {
	fs := flag.NewFlagSet("analyze", flag.ExitOnError)
	outputDir := fs.String("output-dir", "output", "Output directory")
	org := fs.String("org", "", "GitHub organization (auto-detected from go.mod if empty)")
	ver := fs.String("version", "", "Version label for snapshot output (e.g. v2.15.0)")
	fs.Parse(args)

	if fs.NArg() < 1 {
		return fmt.Errorf("usage: arch-analyzer analyze <repo-path> [--output-dir dir] [--org org] [--version label]")
	}

	if *ver != "" {
		if err := validateVersionLabel(*ver); err != nil {
			return err
		}
	}

	repoPath := fs.Arg(0)
	opts := &extractor.ExtractOptions{Org: *org}
	arch, err := extractor.ExtractAll(repoPath, opts)
	if err != nil {
		return err
	}

	outDir := applyVersion(*outputDir, *ver)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	jsonPath := filepath.Join(outDir, "component-architecture.json")
	if err := writeJSON(jsonPath, arch); err != nil {
		return err
	}
	fmt.Printf("Extracted architecture to: %s\n", jsonPath)

	data, err := loadJSON(jsonPath)
	if err != nil {
		return err
	}

	diagramsDir := filepath.Join(outDir, "diagrams")
	diagrams := renderer.RenderAll(data, nil)
	if err := writeDiagrams(diagramsDir, diagrams); err != nil {
		return err
	}
	fmt.Printf("Rendered %d diagram(s) to: %s\n", len(diagrams), diagramsDir)

	if *ver != "" {
		repos := map[string]string{arch.Repo: arch.CommitSHA}
		if err := writeSnapshotMetadata(outDir, *ver, repos, ""); err != nil {
			return fmt.Errorf("writing snapshot metadata: %w", err)
		}
		fmt.Printf("Snapshot metadata written to: %s/snapshot-metadata.json\n", outDir)
	}

	return nil
}

// cmdAggregate merges multiple component JSONs into platform view.
func cmdAggregate(args []string) error {
	fs := flag.NewFlagSet("aggregate", flag.ExitOnError)
	outputDir := fs.String("output-dir", "platform-output", "Output directory")
	ver := fs.String("version", "", "Version label for snapshot output (e.g. v2.15.0)")
	platform := fs.String("platform", "", "Platform name for snapshot metadata (e.g. rhoai, odh)")
	fs.Parse(args)

	if fs.NArg() < 1 {
		return fmt.Errorf("usage: arch-analyzer aggregate <results-dir> [--output-dir dir] [--version label] [--platform name]")
	}

	if *ver != "" {
		if err := validateVersionLabel(*ver); err != nil {
			return err
		}
	}

	platformData, err := aggregator.Aggregate(fs.Arg(0))
	if err != nil {
		return err
	}

	outDir := applyVersion(*outputDir, *ver)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	jsonPath := filepath.Join(outDir, "platform-architecture.json")
	if err := writeJSON(jsonPath, platformData); err != nil {
		return err
	}
	fmt.Printf("Aggregated platform architecture to: %s\n", jsonPath)

	diagramsDir := filepath.Join(outDir, "diagrams")
	diagrams := renderer.RenderPlatformAll(platformData)
	if err := writeDiagrams(diagramsDir, diagrams); err != nil {
		return err
	}
	fmt.Printf("Rendered %d platform diagram(s) to: %s\n", len(diagrams), diagramsDir)

	if *ver != "" {
		repos, err := collectRepoSHAs(fs.Arg(0))
		if err != nil {
			return fmt.Errorf("collecting repo SHAs: %w", err)
		}
		if err := writeSnapshotMetadata(outDir, *ver, repos, *platform); err != nil {
			return fmt.Errorf("writing snapshot metadata: %w", err)
		}
		fmt.Printf("Snapshot metadata written to: %s/snapshot-metadata.json\n", outDir)
	}

	return nil
}

// cmdDocs generates a browsable documentation site from architecture JSON.
func cmdDocs(args []string) error {
	fs := flag.NewFlagSet("docs", flag.ExitOnError)
	outputDir := fs.String("output-dir", "docs", "Output directory for generated docs")
	prefix := fs.String("prefix", "", "Path prefix for nav snippet (e.g. 'rhoai-platform')")
	ver := fs.String("version", "", "Version label for snapshot output (e.g. v2.15.0)")
	fs.Parse(args)

	if fs.NArg() < 1 {
		return fmt.Errorf("usage: arch-analyzer docs <json-file> [--output-dir dir] [--prefix path] [--version label]")
	}

	if *ver != "" {
		if err := validateVersionLabel(*ver); err != nil {
			return err
		}
	}

	data, err := loadJSON(fs.Arg(0))
	if err != nil {
		return err
	}

	pages := renderer.RenderDocs(data)

	// Look for snapshot metadata next to the input JSON file
	var banner string
	metaPath := filepath.Join(filepath.Dir(fs.Arg(0)), "snapshot-metadata.json")
	if metaData, err := loadJSON(metaPath); err == nil {
		snapVer, _ := metaData["version"].(string)
		snapTS, _ := metaData["timestamp"].(string)
		if snapTS != "" {
			// Truncate to date only for display
			if len(snapTS) >= 10 {
				snapTS = snapTS[:10]
			}
		}
		if snapVer != "" {
			banner = fmt.Sprintf("> **Architecture snapshot: %s** (%s)\n\n", snapVer, snapTS)
		}
	}

	outDir := applyVersion(*outputDir, *ver)
	for _, page := range pages {
		content := page.Content
		if banner != "" && strings.HasSuffix(page.Path, "index.md") {
			// Inject version banner after the first heading
			if idx := strings.Index(content, "\n"); idx > 0 {
				content = content[:idx+1] + "\n" + banner + content[idx+1:]
			}
		}
		outPath := filepath.Join(outDir, page.Path)
		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return fmt.Errorf("creating directory for %s: %w", outPath, err)
		}
		if err := os.WriteFile(outPath, []byte(content), 0o644); err != nil {
			return fmt.Errorf("writing %s: %w", outPath, err)
		}
	}

	fmt.Printf("Generated %d documentation pages to: %s\n", len(pages), outDir)

	// Print nav snippet
	navSnippet := renderer.NavSnippet(pages, *prefix)
	if navSnippet != "" {
		fmt.Println("\nmkdocs.yml nav snippet:")
		fmt.Println(navSnippet)
	}

	return nil
}

// cmdExtractSchema extracts CRD JSON schemas from a repo.
func cmdExtractSchema(args []string) error {
	fs := flag.NewFlagSet("extract-schema", flag.ExitOnError)
	outputDir := fs.String("output-dir", "contracts/schemas", "Output directory for schemas")
	repoName := fs.String("repo-name", "", "Repository name (default: directory name)")
	fs.Parse(args)

	if fs.NArg() < 1 {
		return fmt.Errorf("usage: arch-analyzer extract-schema <repo-path> [--output-dir dir] [--repo-name name]")
	}

	repoPath := fs.Arg(0)
	schemas, err := validator.ExtractSchemasFromDir(repoPath)
	if err != nil {
		return err
	}

	if len(schemas) == 0 {
		fmt.Printf("No CRD schemas found in %s\n", repoPath)
		return nil
	}

	name := *repoName
	if name == "" {
		name = filepath.Base(repoPath)
	}

	schemaDir := filepath.Join(*outputDir, name)
	if err := os.MkdirAll(schemaDir, 0o755); err != nil {
		return fmt.Errorf("creating schema dir: %w", err)
	}

	for _, s := range schemas {
		outPath := filepath.Join(schemaDir, s.ResourceKey+".json")
		if err := writeJSON(outPath, s.Schema); err != nil {
			return fmt.Errorf("writing schema %s: %w", s.ResourceKey, err)
		}
		fmt.Printf("Extracted: %s -> %s\n", s.ResourceKey, outPath)
	}

	fmt.Printf("Extracted %d schema(s) to %s\n", len(schemas), schemaDir)
	return nil
}

// cmdValidate validates CRD changes against stored contracts.
func cmdValidate(args []string) error {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)
	contractsDir := fs.String("contracts-dir", "contracts", "Path to contracts directory")
	repoName := fs.String("repo-name", "", "Repository name (default: directory name)")
	fs.Parse(args)

	if fs.NArg() < 1 {
		return fmt.Errorf("usage: arch-analyzer validate <repo-path> [--contracts-dir dir] [--repo-name name]")
	}

	repoPath := fs.Arg(0)
	schemas, err := validator.ExtractSchemasFromDir(repoPath)
	if err != nil {
		return err
	}

	if len(schemas) == 0 {
		fmt.Printf("No CRD schemas found in %s, nothing to validate\n", repoPath)
		return nil
	}

	name := *repoName
	if name == "" {
		name = filepath.Base(repoPath)
	}

	result, err := validator.CheckContract(name, schemas, *contractsDir)
	if err != nil {
		return err
	}

	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
	fmt.Printf("Contract Validation: %s\n", name)
	fmt.Printf("%s\n", strings.Repeat("=", 60))

	for _, check := range result.Checks {
		symbol := "v"
		status := "PASS"
		if !check.IsCompatible {
			symbol = "X"
			status = "FAIL"
		}
		fmt.Printf("  [%s] %s: %s\n", symbol, check.Resource, status)
		for _, d := range check.Details {
			fmt.Printf("      - %s\n", d.Description)
		}
		if len(check.Consumers) > 0 {
			fmt.Printf("      Consumers: %s\n", strings.Join(check.Consumers, ", "))
		}
	}

	if len(result.AffectedConsumers) > 0 {
		fmt.Printf("\nAFFECTED CONSUMERS:\n")
		for _, c := range result.AffectedConsumers {
			fmt.Printf("  - %s: %s\n", c.Repo, c.Usage)
			for _, bc := range c.BreakingChanges {
				fmt.Printf("      Breaking: %s\n", bc)
			}
		}
	}

	if result.IsCompatible {
		fmt.Printf("\nResult: COMPATIBLE\n")
		return nil
	}
	fmt.Printf("\nResult: BREAKING CHANGES DETECTED\n")
	return fmt.Errorf("breaking changes detected")
}

// cmdScan builds a code property graph and runs security queries.
func cmdScan(args []string) error {
	fs := flag.NewFlagSet("scan", flag.ExitOnError)
	outputFile := fs.String("output", "", "Output findings JSON file (default: stdout)")
	format := fs.String("format", "text", "Output format: text, json, sarif")
	domainList := fs.String("domains", "", "Comma-separated domains to run (default: all registered)")
	withArch := fs.Bool("with-arch", false, "Cross-reference with architecture data")
	fs.Parse(args)

	if fs.NArg() < 1 {
		return fmt.Errorf("usage: arch-analyzer scan <repo-path> [--output file] [--format text|json|sarif] [--domains sec,test] [--with-arch]")
	}

	repoPath := fs.Arg(0)

	cpg, err := buildCPG(repoPath)
	if err != nil {
		return err
	}

	var archData *domains.ArchitectureData
	if *withArch {
		archData = prepareArchData(repoPath, cpg, "")
	}

	findings, err := runSecurityScan(cpg, *domainList, archData)
	if err != nil {
		return err
	}

	registeredDomains := domains.Names()
	switch *format {
	case "text":
		printFindings(cpg, findings)
	case "json":
		if len(registeredDomains) > 0 {
			return outputJSON(*outputFile, domainGroupedJSON(findings))
		}
		return outputJSON(*outputFile, findings)
	case "sarif":
		return outputSARIF(*outputFile, findings)
	default:
		return fmt.Errorf("unknown format: %s", *format)
	}
	return nil
}

// cmdDomains lists all registered analysis domains.
func cmdDomains() error {
	registered := domains.All()
	if len(registered) == 0 {
		fmt.Println("No domains registered.")
		return nil
	}
	fmt.Printf("%d registered domain(s):\n", len(registered))
	for _, d := range registered {
		fmt.Printf("  %-12s languages: %s", d.Name(), strings.Join(d.SupportedLanguages(), ", "))
		deps := d.Dependencies()
		if len(deps) > 0 {
			fmt.Printf("  deps: %s", strings.Join(deps, ", "))
		}
		fmt.Printf("  queries: %d\n", len(d.Queries()))
	}
	return nil
}

// cmdGraph exports the code property graph.
func cmdGraph(args []string) error {
	fs := flag.NewFlagSet("graph", flag.ExitOnError)
	format := fs.String("format", "json", "Output format: json, dot")
	outputFile := fs.String("output", "", "Output file (default: stdout)")
	fs.Parse(args)

	if fs.NArg() < 1 {
		return fmt.Errorf("usage: arch-analyzer graph <repo-path> [--format json|dot] [--output file]")
	}

	cpg, err := buildCPG(fs.Arg(0))
	if err != nil {
		return err
	}

	var content []byte
	switch *format {
	case "json":
		output := map[string]interface{}{
			"schema_version": 2,
			"nodes":          cpg.Nodes(),
			"edges":          cpg.Edges(),
		}
		content, err = json.MarshalIndent(output, "", "  ")
		if err != nil {
			return err
		}
		content = append(content, '\n')
	case "dot":
		content = []byte(renderDOT(cpg))
	default:
		return fmt.Errorf("unknown format: %s", *format)
	}

	if *outputFile != "" {
		return os.WriteFile(*outputFile, content, 0o644)
	}
	_, err = os.Stdout.Write(content)
	return err
}

// cmdFullAnalysis runs architecture extraction + code graph scan.
func cmdFullAnalysis(args []string) error {
	fs := flag.NewFlagSet("full-analysis", flag.ExitOnError)
	outputDir := fs.String("output-dir", "output", "Output directory")
	org := fs.String("org", "", "GitHub organization (auto-detected from go.mod if empty)")
	domainList := fs.String("domains", "", "Comma-separated domains (default: all)")
	ver := fs.String("version", "", "Version label for snapshot output (e.g. v2.15.0)")
	fs.Parse(args)

	if fs.NArg() < 1 {
		return fmt.Errorf("usage: arch-analyzer full-analysis <repo-path> [--output-dir dir] [--org org] [--domains sec,test] [--version label]")
	}

	if *ver != "" {
		if err := validateVersionLabel(*ver); err != nil {
			return err
		}
	}

	repoPath := fs.Arg(0)
	outDir := applyVersion(*outputDir, *ver)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	// Architecture extraction
	fmt.Println("=== Architecture Extraction ===")
	extractOpts := &extractor.ExtractOptions{Org: *org}
	archResult, err := extractor.ExtractAll(repoPath, extractOpts)
	var archData *domains.ArchitectureData
	var parsedArch *arch.Data
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: architecture extraction failed: %v\n", err)
	} else {
		jsonPath := filepath.Join(outDir, "component-architecture.json")
		if err := writeJSON(jsonPath, archResult); err != nil {
			return err
		}
		fmt.Printf("Extracted architecture to: %s\n", jsonPath)

		// Prepare arch data for domain analyzers
		raw, _ := json.Marshal(archResult)
		var data map[string]interface{}
		json.Unmarshal(raw, &data)
		archData = &domains.ArchitectureData{Raw: data}

		parsed, parseErr := arch.Parse(data)
		if parseErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: architecture data parsing failed: %v\n", parseErr)
		} else {
			parsedArch = parsed
		}

		data2, loadErr := loadJSON(jsonPath)
		if loadErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to load architecture JSON: %v\n", loadErr)
		} else if data2 != nil {
			diagramsDir := filepath.Join(outDir, "diagrams")
			diagrams := renderer.RenderAll(data2, nil)
			if wErr := writeDiagrams(diagramsDir, diagrams); wErr != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to write diagrams: %v\n", wErr)
			} else {
				fmt.Printf("Rendered %d diagram(s) to: %s\n", len(diagrams), diagramsDir)
			}
		}
	}

	// Code graph scan with domains
	fmt.Println("\n=== Code Graph Security Scan ===")
	cpg, err := buildCPG(repoPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: code graph build failed: %v\n", err)
	} else {
		if parsedArch != nil {
			cpg.ArchData = parsedArch
		}

		findings, scanErr := runSecurityScan(cpg, *domainList, archData)
		if scanErr != nil {
			return scanErr
		}

		printFindings(cpg, findings)

		findingsPath := filepath.Join(outDir, "security-findings.json")
		if wErr := outputJSON(findingsPath, findings); wErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to write findings: %v\n", wErr)
		} else {
			fmt.Printf("Findings written to: %s\n", findingsPath)
		}

		graphPath := filepath.Join(outDir, "code-graph.json")
		graphData := map[string]interface{}{
			"schema_version": 2,
			"nodes":          cpg.Nodes(),
			"edges":          cpg.Edges(),
		}
		if wErr := writeJSON(graphPath, graphData); wErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to write code graph: %v\n", wErr)
		} else {
			fmt.Printf("Code graph written to: %s\n", graphPath)
		}
	}

	// Schema extraction
	fmt.Println("\n=== CRD Schema Extraction ===")
	schemas, err := validator.ExtractSchemasFromDir(repoPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: schema extraction failed: %v\n", err)
	} else if len(schemas) > 0 {
		schemaDir := filepath.Join(outDir, "schemas")
		if mkErr := os.MkdirAll(schemaDir, 0o755); mkErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to create schema dir: %v\n", mkErr)
		} else {
			for _, s := range schemas {
				outPath := filepath.Join(schemaDir, s.ResourceKey+".json")
				if wErr := writeJSON(outPath, s.Schema); wErr != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to write schema %s: %v\n", s.ResourceKey, wErr)
				}
			}
			fmt.Printf("Extracted %d CRD schema(s) to: %s\n", len(schemas), schemaDir)
		}
	} else {
		fmt.Println("No CRD schemas found")
	}

	if *ver != "" && archResult != nil {
		repos := map[string]string{archResult.Repo: archResult.CommitSHA}
		if err := writeSnapshotMetadata(outDir, *ver, repos, ""); err != nil {
			return fmt.Errorf("writing snapshot metadata: %w", err)
		}
		fmt.Printf("Snapshot metadata written to: %s/snapshot-metadata.json\n", outDir)
	}

	return nil
}

// collectRepoSHAs scans a results directory for component-architecture.json files
// and extracts repo -> commit_sha pairs for snapshot metadata.
func collectRepoSHAs(resultsDir string) (map[string]string, error) {
	repos := make(map[string]string)
	entries, err := os.ReadDir(resultsDir)
	if err != nil {
		return repos, fmt.Errorf("reading results directory %s: %w", resultsDir, err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		jsonPath := filepath.Join(resultsDir, entry.Name(), "component-architecture.json")
		data, err := loadJSON(jsonPath)
		if err != nil {
			fmt.Printf("WARN: failed to load %s: %v\n", jsonPath, err)
			continue
		}
		repo, _ := data["repo"].(string)
		sha, _ := data["commit_sha"].(string)
		if repo != "" {
			repos[repo] = sha
		}
	}
	return repos, nil
}

// runSecurityScan runs legacy queries and domain analyzers against a CPG,
// returning all findings. If domainList is non-empty, only those domains run.
func runSecurityScan(cpg *graph.CPG, domainList string, archData *domains.ArchitectureData) ([]query.Finding, error) {
	engine := query.NewEngine()
	findings := engine.RunAll(cpg)

	var analyzers []domains.DomainAnalyzer
	if domainList != "" {
		names := strings.Split(domainList, ",")
		resolved, resolveErr := domains.ResolveDependencies(names)
		if resolveErr != nil {
			return nil, resolveErr
		}
		var err error
		analyzers, err = domains.Get(resolved)
		if err != nil {
			return nil, err
		}
	} else {
		analyzers = domains.All()
	}

	if len(analyzers) > 0 {
		orch := domains.NewOrchestrator(analyzers)
		results, runErr := orch.Run(cpg, "go", archData)
		if runErr != nil {
			return nil, fmt.Errorf("domain analysis: %w", runErr)
		}
		for _, dr := range results {
			fmt.Printf("Domain %s: %d annotations, %d findings\n",
				dr.Domain, dr.AnnotationsAdded, len(dr.Findings))
			findings = append(findings, dr.Findings...)
		}
	}

	return findings, nil
}

// prepareArchData extracts architecture data from a repo and sets it on the CPG.
func prepareArchData(repoPath string, cpg *graph.CPG, org string) *domains.ArchitectureData {
	opts := &extractor.ExtractOptions{Org: org}
	archResult, err := extractor.ExtractAll(repoPath, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: architecture extraction failed: %v\n", err)
		return nil
	}

	raw, _ := json.Marshal(archResult)
	var data map[string]interface{}
	json.Unmarshal(raw, &data)
	archData := &domains.ArchitectureData{Raw: data}

	parsed, parseErr := arch.Parse(data)
	if parseErr != nil {
		fmt.Fprintf(os.Stderr, "Warning: architecture data parsing failed: %v\n", parseErr)
	} else {
		cpg.ArchData = parsed
	}

	return archData
}

// buildCPG constructs a Code Property Graph from a repo directory.
func buildCPG(repoPath string) (*graph.CPG, error) {
	b := builder.NewBuilder()
	cpg, err := b.BuildFromDir(repoPath)
	if err != nil {
		return nil, fmt.Errorf("building code graph: %w", err)
	}

	sl := linker.NewStorageLinker()
	linked := sl.Link(cpg)

	sa := annotator.NewSecurityAnnotator()
	sa.Annotate(cpg)

	fmt.Printf("Graph: %d nodes, %d edges, %d storage links\n",
		len(cpg.Nodes()), len(cpg.Edges()), linked)
	fmt.Printf("  Functions: %d, Call sites: %d, HTTP handlers: %d, DB ops: %d\n",
		len(cpg.NodesByKind(graph.NodeFunction)),
		len(cpg.NodesByKind(graph.NodeCallSite)),
		len(cpg.NodesByKind(graph.NodeHTTPEndpoint)),
		len(cpg.NodesByKind(graph.NodeDBOperation)))

	return cpg, nil
}

func printFindings(cpg *graph.CPG, findings []query.Finding) {
	if len(findings) == 0 {
		fmt.Println("No security findings.")
		return
	}
	fmt.Printf("\n%d security finding(s):\n", len(findings))
	for _, f := range findings {
		fmt.Printf("  [%s] %s: %s (%s:%d)\n", f.Severity, f.RuleID, f.Message, f.File, f.Line)
	}
}

func renderDOT(cpg *graph.CPG) string {
	var b strings.Builder
	b.WriteString("digraph CPG {\n")
	b.WriteString("  rankdir=LR;\n")
	for _, n := range cpg.Nodes() {
		label := fmt.Sprintf("%s\\n(%s)", n.Name, n.Kind)
		fmt.Fprintf(&b, "  %q [label=%q];\n", n.ID, label)
	}
	for _, e := range cpg.Edges() {
		fmt.Fprintf(&b, "  %q -> %q [label=%q];\n", e.From, e.To, e.Kind)
	}
	b.WriteString("}\n")
	return b.String()
}

func outputSARIF(path string, findings []query.Finding) error {
	// Group findings by domain for per-domain SARIF runs
	grouped := make(map[string][]query.Finding)
	for _, f := range findings {
		domain := f.Domain
		if domain == "" {
			domain = "legacy"
		}
		grouped[domain] = append(grouped[domain], f)
	}

	var runs []map[string]interface{}
	for domain, domainFindings := range grouped {
		runs = append(runs, map[string]interface{}{
			"tool": map[string]interface{}{
				"driver": map[string]interface{}{
					"name":    "arch-analyzer/" + domain,
					"version": version,
				},
			},
			"results": sarifResults(domainFindings),
		})
	}

	sarif := map[string]interface{}{
		"$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/main/sarif-2.1/schema/sarif-schema-2.1.0.json",
		"version": "2.1.0",
		"runs":    runs,
	}
	if path != "" {
		return writeJSON(path, sarif)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(sarif)
}

func domainGroupedJSON(findings []query.Finding) map[string]interface{} {
	grouped := make(map[string][]query.Finding)
	for _, f := range findings {
		domain := f.Domain
		if domain == "" {
			domain = "legacy"
		}
		grouped[domain] = append(grouped[domain], f)
	}

	domainResults := make(map[string]interface{})
	for domain, domainFindings := range grouped {
		domainResults[domain] = map[string]interface{}{
			"findings": domainFindings,
			"count":    len(domainFindings),
		}
	}

	return map[string]interface{}{
		"domains":        domainResults,
		"total_findings": len(findings),
	}
}

func sarifResults(findings []query.Finding) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(findings))
	for _, f := range findings {
		r := map[string]interface{}{
			"ruleId": f.RuleID,
			"level":  sarifLevel(f.Severity),
			"message": map[string]string{
				"text": f.Message,
			},
			"locations": []map[string]interface{}{
				{
					"physicalLocation": map[string]interface{}{
						"artifactLocation": map[string]string{
							"uri": f.File,
						},
						"region": map[string]int{
							"startLine": f.Line,
						},
					},
				},
			},
		}
		results = append(results, r)
	}
	return results
}

func sarifLevel(severity string) string {
	switch strings.ToLower(severity) {
	case "critical", "high":
		return "error"
	case "medium":
		return "warning"
	default:
		return "note"
	}
}

func outputJSON(path string, data interface{}) error {
	if path != "" {
		return writeJSON(path, data)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func writeJSON(path string, data interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return fmt.Errorf("creating directory for %s: %w", path, err)
	}
	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(raw, '\n'), 0o640)
}

func loadJSON(path string) (map[string]interface{}, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	return data, nil
}

func writeDiagrams(dir string, diagrams map[string]string) error {
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("creating diagram dir: %w", err)
	}
	for filename, content := range diagrams {
		path := filepath.Join(dir, filename)
		if err := os.WriteFile(path, []byte(content), 0o640); err != nil {
			return fmt.Errorf("writing diagram %s: %w", path, err)
		}
	}
	return nil
}
