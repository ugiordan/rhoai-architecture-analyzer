# Scan Configuration

The `scan-config.yaml` file defines which repositories to analyze in batch mode and CI workflows.

## Format

```yaml
repos:
  - name: odh-model-controller
    url: https://github.com/opendatahub-io/odh-model-controller
  - name: model-registry-operator
    url: https://github.com/opendatahub-io/model-registry-operator
  - name: data-science-pipelines-operator
    url: https://github.com/opendatahub-io/data-science-pipelines-operator
  - name: kubeflow
    url: https://github.com/opendatahub-io/kubeflow
  # ... more repos
```

## Usage

The scan config is used by:

- `scripts/analyze-repo.sh`: Wrapper script that clones, analyzes, and cleans up
- `analyze-all.yml` workflow: Weekly scheduled analysis of all RHOAI repos
- `extract-schemas.yml` workflow: CRD schema extraction across repos

## Batch analysis

Using the wrapper script:

```bash
# Analyze a single repo from config
./scripts/analyze-repo.sh odh-model-controller output/

# The script handles:
# 1. Clone to temp directory
# 2. Run rhoai-analyzer analyze
# 3. Copy results to output directory
# 4. Clean up temp clone
```

## Adding a repository

Add an entry to `scan-config.yaml`:

```yaml
  - name: my-new-operator
    url: https://github.com/opendatahub-io/my-new-operator
```

The name must match the repository name (used for output directory naming and artifact labeling).

## Platform aggregation

After analyzing multiple repos, aggregate results into a platform-wide view:

```bash
# Analyze each repo
for repo in $(yq '.repos[].name' scan-config.yaml); do
  ./scripts/analyze-repo.sh "$repo" results/"$repo"/
done

# Aggregate all results
rhoai-analyzer aggregate results/ --output-dir platform-output/
```

See [Platform Aggregation](../guides/platform-aggregation.md) for details.
