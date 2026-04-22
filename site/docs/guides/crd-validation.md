# CRD Contract Validation

The analyzer extracts CRD JSON schemas and validates them against a baseline to detect breaking changes.

## How it works

1. **Extract schemas**: Pull CRD JSON schemas from repository manifests
2. **Store as baseline**: Commit schemas to `contracts/schemas/` directory
3. **Validate on change**: When CRD definitions change, compare new schemas against baseline
4. **Report breaks**: Flag removed fields, type changes, and other breaking changes

## Usage

### Extract schemas

```bash
arch-analyzer extract-schema /path/to/repo --output-dir contracts/schemas
```

This produces one JSON schema file per CRD found in the repository.

### Validate changes

```bash
arch-analyzer validate /path/to/repo --contracts-dir contracts
```

The validator checks for:

| Change type | Severity | Example |
|-------------|----------|---------|
| Field removed | Breaking | `spec.replicas` removed from CRD |
| Type changed | Breaking | `spec.timeout` changed from string to integer |
| Required field added | Breaking | New required field `spec.config` without default |
| Enum value removed | Breaking | Allowed value `v1alpha1` removed from version enum |
| Optional field added | Non-breaking | New optional `spec.annotations` field |
| Description changed | Non-breaking | Updated field description |

## CI workflow

The `validate-contracts.yml` workflow runs automatically on PRs that modify the `contracts/` directory:

```yaml
# Triggered by changes to contracts/
on:
  push:
    paths: ['contracts/**']
  pull_request:
    paths: ['contracts/**']
```

Workflow steps:

1. Detect which schema files changed
2. For each changed repo, extract current schemas
3. Compare current vs. baseline schemas
4. Report breaking changes
5. Exit with code 1 if breaking changes detected (blocks PR merge)

## Schema extraction workflow

The `extract-schemas.yml` workflow runs weekly:

1. Clones all repos from `scan-config.yaml`
2. Extracts CRD JSON schemas from each
3. If schemas changed, creates a PR with the updates
4. The PR triggers `validate-contracts.yml` for review

## Setting up contract validation

### For a single repo

```bash
# Initial baseline extraction
arch-analyzer extract-schema /path/to/repo --output-dir contracts/schemas

# Commit the baseline
git add contracts/
git commit -m "Add CRD schema baseline for contract validation"

# On subsequent changes, validate
arch-analyzer validate /path/to/repo --contracts-dir contracts
```

### For multiple repos (platform-wide)

Use `scan-config.yaml` to define all repos, then run the extraction workflow. Each repo's schemas are stored in `contracts/schemas/<repo-name>/`.

## Consumer tracking

The validator also identifies affected consumers: other repos that reference the changed CRD. This is done by checking controller watches and RBAC definitions across all repos in the scan config.
