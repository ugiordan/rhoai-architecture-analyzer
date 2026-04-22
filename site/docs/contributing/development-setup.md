# Development Setup

## Prerequisites

- Go 1.25+
- git

## Clone and build

```bash
git clone https://github.com/ugiordan/architecture-analyzer.git
cd architecture-analyzer
go build -o arch-analyzer ./cmd/arch-analyzer/
```

## Project structure

```
architecture-analyzer/
  cmd/arch-analyzer/
    main.go                  # CLI entry point with subcommands
  pkg/
    extractor/               # 18 architecture extractors
    renderer/                # 7 diagram/report renderers
    aggregator/              # Platform-wide aggregation
    validator/               # CRD contract validation
    parser/                  # Tree-sitter Go parser
    builder/                 # Code property graph builder
    graph/                   # CPG data structure
    annotator/               # Annotation engine
    query/                   # Query engine + taint analysis
    domains/                 # Pluggable domain framework
      security/              # Security domain
      testing/               # Testing domain
      upgrade/               # Upgrade domain
    arch/                    # Architecture data types
    linker/                  # Storage linker
    config/                  # Configuration types
  scripts/
    analyze-repo.sh          # Clone + analyze + cleanup wrapper
  .github/workflows/
    analyze-all.yml          # Scheduled weekly analysis
    extract-schemas.yml      # CRD schema extraction
    validate-contracts.yml   # Contract validation on PR
  testdata/                  # Test fixtures
  scan-config.yaml           # RHOAI repos list
```

## Running tests

```bash
go test ./...
```

Key test files:

- `pkg/extractor/extract_test.go`: Extractor integration tests
- `pkg/graph/cpg_test.go`: CPG data structure tests
- `pkg/parser/go_parser_test.go`: Tree-sitter parser tests
- `pkg/renderer/renderer_test.go`: Renderer output tests
- `pkg/query/engine_test.go`: Query engine tests
- `pkg/domains/*/queries_test.go`: Per-domain query tests

## Testing against real repos

```bash
# Analyze a repo (e.g. an RHOAI component)
./arch-analyzer analyze /path/to/odh-model-controller --output-dir /tmp/analysis

# Run security scan
./arch-analyzer scan /path/to/odh-model-controller --format json --output /tmp/findings.json
```

## Adding new functionality

### New extractor

See [Adding Extractors](adding-extractors.md).

### New renderer

1. Create `pkg/renderer/myformat.go`
2. Implement the renderer function reading from `ComponentArchitecture`
3. Register in `pkg/renderer/renderer.go`
4. Add format name to CLI `--formats` flag
5. Add tests

### New domain

1. Create `pkg/domains/mydomain/` with analyzer, annotations, annotator, queries
2. Register in `main.go`
3. Add tests for annotator and queries

### New CLI command

1. Add command definition in `cmd/arch-analyzer/main.go`
2. Wire to appropriate pkg functions
3. Add usage documentation
