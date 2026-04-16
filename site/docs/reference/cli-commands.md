# CLI Commands

## rhoai-analyzer analyze

Extract architecture data and render diagrams.

```bash
rhoai-analyzer analyze <repo-path> --output-dir <dir>
```

| Flag | Description |
|------|-------------|
| `--output-dir` | Directory for all output files |

**Output**: `component-architecture.json` + all diagram/report files.

## rhoai-analyzer extract

Extract architecture data only (no rendering).

```bash
rhoai-analyzer extract <repo-path> --output <file>
```

| Flag | Description |
|------|-------------|
| `--output` | Output JSON file path |

## rhoai-analyzer render

Render diagrams from existing JSON.

```bash
rhoai-analyzer render <json-file> --output-dir <dir> [--formats <list>]
```

| Flag | Description |
|------|-------------|
| `--output-dir` | Directory for diagram files |
| `--formats` | Comma-separated list: `rbac`, `component`, `security`, `dependencies`, `c4`, `dataflow`, `report` |

## rhoai-analyzer scan

Build code property graph and run security queries.

```bash
rhoai-analyzer scan <repo-path> [flags]
```

| Flag | Description |
|------|-------------|
| `--format` | Output format: `json` or `sarif` |
| `--output` | Output file path |
| `--domains` | Comma-separated domain list: `security`, `testing`, `upgrade` |
| `--with-arch` | Path to `component-architecture.json` for enrichment |

## rhoai-analyzer full-analysis

Run everything: extract + render + scan + extract-schema.

```bash
rhoai-analyzer full-analysis <repo-path> --output-dir <dir>
```

| Flag | Description |
|------|-------------|
| `--output-dir` | Directory for all output files |

## rhoai-analyzer aggregate

Merge multiple component analyses into platform view.

```bash
rhoai-analyzer aggregate <results-dir> --output-dir <dir>
```

| Flag | Description |
|------|-------------|
| `--output-dir` | Directory for platform output |

Discovers all `component-architecture.json` files in the results directory recursively.

## rhoai-analyzer extract-schema

Extract CRD JSON schemas for contract validation.

```bash
rhoai-analyzer extract-schema <repo-path> --output-dir <dir>
```

| Flag | Description |
|------|-------------|
| `--output-dir` | Directory for schema JSON files |

## rhoai-analyzer validate

Validate CRD changes against baseline contracts.

```bash
rhoai-analyzer validate <repo-path> --contracts-dir <dir>
```

| Flag | Description |
|------|-------------|
| `--contracts-dir` | Directory containing baseline schemas |

Exit code 1 if breaking changes detected.

## rhoai-analyzer graph

Export the code property graph.

```bash
rhoai-analyzer graph <repo-path> [flags]
```

Exports the raw CPG for inspection or custom analysis.

## rhoai-analyzer domains

List all registered analysis domains.

```bash
rhoai-analyzer domains
```

Output includes domain name, description, and available queries.

## rhoai-analyzer version

Print version information.

```bash
rhoai-analyzer version
```
