# Platform Aggregation

The aggregator merges multiple single-component analyses into a cross-repo platform view, revealing inter-component relationships that are invisible when analyzing repos in isolation.

## Usage

```bash
# Step 1: Analyze individual repos
rhoai-analyzer analyze /path/to/repo-a --output-dir results/repo-a
rhoai-analyzer analyze /path/to/repo-b --output-dir results/repo-b
rhoai-analyzer analyze /path/to/repo-c --output-dir results/repo-c

# Step 2: Aggregate
rhoai-analyzer aggregate results/ --output-dir platform-output/
```

The aggregator discovers all `component-architecture.json` files in the results directory and merges them.

## What aggregation reveals

### Cross-component CRD ownership

Which component owns each CRD (defines it) and which components watch or reference it:

```
CRD: datasciencecluster.datasciencecluster.opendatahub.io
  Owner: opendatahub-operator
  Watchers: odh-model-controller, data-science-pipelines-operator
```

### Dependency graph

How internal ODH components depend on each other through Go module imports:

```
opendatahub-operator
  -> odh-model-controller (v0.12.0)
  -> model-registry-operator (v0.8.0)
  -> kubeflow (v1.9.0)
```

### RBAC overlap

Multiple components requesting permissions on the same resources, which may indicate redundancy or over-permissioning:

```
Resource: secrets (core/v1)
  Reader: odh-model-controller
  Reader: model-registry-operator
  Admin: opendatahub-operator
```

### Service mesh

Cross-component network policies showing which components can communicate and through which ports.

## Platform output

The aggregator produces:

- `platform-architecture.json`: Merged data with cross-references
- `platform-report.md`: Platform-level report with all cross-component findings
- Mermaid diagrams showing the full platform topology

## Batch analysis with scan config

For automated platform analysis:

```bash
# Analyze all repos from scan-config.yaml
for repo in $(yq '.repos[].name' scan-config.yaml); do
  ./scripts/analyze-repo.sh "$repo" results/"$repo"/
done

# Aggregate
rhoai-analyzer aggregate results/ --output-dir platform-output/
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
