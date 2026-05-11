# Platform Aggregation

The aggregator merges multiple single-component analyses into a cross-repo platform view, revealing inter-component relationships that are invisible when analyzing repos in isolation.

## Usage

```bash
# Step 1: Analyze individual repos (using the batch script)
# Output is org-namespaced: results/<org>/<repo>/
scripts/analyze-repo.sh myorg/repo-a results/
scripts/analyze-repo.sh myorg/repo-b results/
scripts/analyze-repo.sh otherorg/repo-c results/

# Or manually:
arch-analyzer analyze /path/to/repo-a --output-dir results/myorg/repo-a
arch-analyzer analyze /path/to/repo-b --output-dir results/myorg/repo-b

# Step 2: Aggregate
arch-analyzer aggregate results/ --output-dir platform-output/
```

The aggregator recursively discovers all `component-architecture.json` files in the results directory and merges them. Both flat (`results/<repo>/`) and org-namespaced (`results/<org>/<repo>/`) layouts are supported.

## What aggregation reveals

### Cross-component CRD ownership

Which component owns each CRD (defines it) and which components watch or reference it:

```
CRD: myresource.example.io
  Owner: platform-operator
  Watchers: resource-controller, pipeline-operator
```

### Dependency graph

How components depend on each other through Go module imports:

```
platform-operator
  -> resource-controller (v0.12.0)
  -> registry-operator (v0.8.0)
  -> shared-library (v1.9.0)
```

### RBAC overlap

Multiple components requesting permissions on the same resources, which may indicate redundancy or over-permissioning:

```
Resource: secrets (core/v1)
  Reader: resource-controller
  Reader: registry-operator
  Admin: platform-operator
```

### Service mesh

Cross-component network policies showing which components can communicate and through which ports.

## Platform output

The aggregator produces:

- `platform-architecture.json`: Merged data with cross-references
- `PLATFORM.md`: Platform-level report with all cross-component findings
- Mermaid diagrams showing the full platform topology

### Generating browsable documentation

The `docs` command generates a full set of markdown pages with embedded mermaid diagrams from the aggregated JSON, ready to drop into any mkdocs site:

```bash
arch-analyzer docs --output-dir site/docs/my-platform platform-output/platform-architecture.json
```

This produces:

- Platform overview with dependency graph
- Network topology with service map diagrams
- RBAC surface with permission scope visualization
- Secrets inventory with distribution diagrams
- Per-component deep-dive pages (overview, network, RBAC, security, dataflow) with inline mermaid

See the [Use Cases](../use-cases/index.md) section for live examples generated from real Kubernetes platforms.

## Batch analysis with scan config

For automated platform analysis, define platforms and repos in `scan-config.yaml`:

```bash
# Analyze all repos from scan-config.yaml
# Results are org-namespaced automatically: results/<org>/<repo>/
for repo in $(yq '.repos[].name' scan-config.yaml); do
  ./scripts/analyze-repo.sh "$repo" results/
done

# Aggregate (recursively discovers all component JSONs)
arch-analyzer aggregate results/ --output-dir platform-output/
```

The `analyze-all.yml` GitHub Actions workflow does this automatically every Monday.

## Weekly analysis workflow

The `analyze-all.yml` workflow:

1. Runs every Monday at 06:00 UTC (or manual dispatch)
2. Builds the analyzer
3. Iterates repos from `scan-config.yaml`
4. Runs `scripts/analyze-repo.sh` for each repo
5. Aggregates all results into a platform view
6. Uploads component and platform results as artifacts (90-day retention)
