# CLI Commands

## arch-analyzer analyze

Extract architecture data and render diagrams.

```bash
arch-analyzer analyze <repo-path> --output-dir <dir>
```

| Flag | Description |
|------|-------------|
| `--output-dir` | Directory for all output files |

**Output**: `component-architecture.json` + all diagram/report files.

## arch-analyzer extract

Extract architecture data only (no rendering).

```bash
arch-analyzer extract <repo-path> --output <file>
```

| Flag | Description |
|------|-------------|
| `--output` | Output JSON file path |

## arch-analyzer render

Render diagrams from existing JSON.

```bash
arch-analyzer render <json-file> --output-dir <dir> [--formats <list>]
```

| Flag | Description |
|------|-------------|
| `--output-dir` | Directory for diagram files |
| `--formats` | Comma-separated list: `rbac`, `component`, `security`, `dependencies`, `c4`, `dataflow`, `report` |

## arch-analyzer scan

Build code property graph and run security queries.

```bash
arch-analyzer scan <repo-path> [flags]
```

| Flag | Description |
|------|-------------|
| `--format` | Output format: `json` or `sarif` |
| `--output` | Output file path |
| `--domains` | Comma-separated domain list: `security`, `testing`, `upgrade` |
| `--with-arch` | Path to `component-architecture.json` for enrichment |

## arch-analyzer full-analysis

Run everything: extract + render + scan + extract-schema.

```bash
arch-analyzer full-analysis <repo-path> --output-dir <dir>
```

| Flag | Description |
|------|-------------|
| `--output-dir` | Directory for all output files |

## arch-analyzer aggregate

Merge multiple component analyses into platform view.

```bash
arch-analyzer aggregate <results-dir> --output-dir <dir>
```

| Flag | Description |
|------|-------------|
| `--output-dir` | Directory for platform output |

Discovers all `component-architecture.json` files in the results directory recursively.

## arch-analyzer docs

Generate browsable documentation pages from architecture JSON.

```bash
arch-analyzer docs --output-dir <dir> <json-file>
```

| Flag | Description |
|------|-------------|
| `--output-dir` | Directory for generated markdown pages (default: `docs`) |
| `--prefix` | Path prefix for the nav snippet output |

Auto-detects whether the input is a single component or aggregated platform JSON. For platform data, generates per-component deep-dive pages under `components/`. Outputs a mkdocs.yml nav snippet for integration.

## arch-analyzer extract-schema

Extract CRD JSON schemas for contract validation.

```bash
arch-analyzer extract-schema <repo-path> --output-dir <dir>
```

| Flag | Description |
|------|-------------|
| `--output-dir` | Directory for schema JSON files |

## arch-analyzer validate

Validate CRD changes against baseline contracts.

```bash
arch-analyzer validate <repo-path> --contracts-dir <dir>
```

| Flag | Description |
|------|-------------|
| `--contracts-dir` | Directory containing baseline schemas |

Exit code 1 if breaking changes detected.

## arch-analyzer graph

Export the code property graph.

```bash
arch-analyzer graph <repo-path> [flags]
```

Exports the raw CPG for inspection or custom analysis.

## arch-analyzer domains

List all registered analysis domains.

```bash
arch-analyzer domains
```

Output includes domain name, description, and available queries.

## arch-analyzer version

Print version information.

```bash
arch-analyzer version
```
