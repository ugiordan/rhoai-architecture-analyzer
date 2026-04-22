# CI Integration

The analyzer provides three GitHub Actions workflows for automated analysis.

## Workflows overview

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| `analyze-all.yml` | Weekly Monday 06:00 UTC, manual | Full platform analysis |
| `extract-schemas.yml` | Weekly Monday 06:00 UTC | CRD schema extraction and PR creation |
| `validate-contracts.yml` | Push/PR to `contracts/` | Breaking change detection |

## analyze-all.yml

Scheduled weekly analysis of all configured platform repos (e.g. RHOAI, ODH):

```yaml
name: Analyze All Repos

on:
  schedule:
    - cron: '0 6 * * 1'  # Monday 06:00 UTC
  workflow_dispatch:

jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.25.0"

      - name: Build analyzer
        run: go build -o arch-analyzer ./cmd/arch-analyzer/

      - name: Analyze repos
        run: |
          for repo in $(yq '.repos[].name' scan-config.yaml); do
            ./scripts/analyze-repo.sh "$repo" results/"$repo"/
          done

      - name: Aggregate platform
        run: ./arch-analyzer aggregate results/ --output-dir platform-output/

      - uses: actions/upload-artifact@v4
        with:
          name: platform-analysis
          path: |
            results/
            platform-output/
          retention-days: 90
```

## Adding to your own repo

### Basic analysis on PR

```yaml
name: Architecture Analysis

on:
  pull_request:
    paths:
      - 'config/**'
      - 'pkg/**'
      - 'cmd/**'
      - 'go.mod'

jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.25.0"

      - name: Install analyzer
        run: |
          git clone https://github.com/ugiordan/architecture-analyzer.git /tmp/analyzer
          cd /tmp/analyzer && go build -o /usr/local/bin/arch-analyzer ./cmd/arch-analyzer/

      - name: Run analysis
        run: arch-analyzer analyze . --output-dir analysis/

      - uses: actions/upload-artifact@v4
        with:
          name: architecture-analysis
          path: analysis/
```

### Security scan with SARIF upload

```yaml
      - name: Security scan
        run: arch-analyzer scan . --format sarif --output findings.sarif

      - name: Upload SARIF
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: findings.sarif
```

### CRD contract validation on PR

```yaml
name: Validate CRD Contracts

on:
  pull_request:
    paths: ['config/crd/**']

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.25.0"

      - name: Install analyzer
        run: |
          git clone https://github.com/ugiordan/architecture-analyzer.git /tmp/analyzer
          cd /tmp/analyzer && go build -o /usr/local/bin/arch-analyzer ./cmd/arch-analyzer/

      - name: Validate schemas
        run: arch-analyzer validate . --contracts-dir contracts
```

## Artifacts

All workflows upload results as GitHub Actions artifacts:

- **Component results**: Per-repo JSON, diagrams, and reports
- **Platform results**: Aggregated platform view
- **Retention**: 90 days by default

Access artifacts from the Actions tab in your repository.
