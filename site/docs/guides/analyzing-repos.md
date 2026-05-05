# Analyzing Repositories

## Single repository analysis

### Extract + render (recommended)

```bash
arch-analyzer analyze /path/to/repo --output-dir output/
```

This is the most common operation. It:

1. Runs all 22 extractors against the repository
2. Produces `component-architecture.json` with all extracted data
3. Renders all 7 diagram/report formats

### Extract only

If you only need the JSON data (for custom processing or aggregation):

```bash
arch-analyzer extract /path/to/repo --output component-architecture.json
```

### Render from existing JSON

If you already have extracted JSON and want different renderers or formats:

```bash
# All formats
arch-analyzer render component-architecture.json --output-dir diagrams/

# Specific formats only
arch-analyzer render component-architecture.json --formats rbac,component
```

Available format names: `rbac`, `component`, `security`, `dependencies`, `c4`, `dataflow`, `report`.

## Full analysis

Combines architecture extraction, diagram rendering, code graph scanning, and schema extraction:

```bash
arch-analyzer full-analysis /path/to/repo --output-dir output/
```

Output includes everything from `analyze` plus:

- Security findings from code property graph queries
- CRD JSON schemas for contract validation

## What each extractor produces

The analyzer walks the repository looking for specific file patterns. Each extractor operates independently:

- **YAML extractors** (CRDs, RBAC, services, deployments, etc.) parse Kubernetes manifests
- **Go source extractors** (controller watches, HTTP endpoints, cache config, operator constants, reconcile sequences, Prometheus metrics, status conditions, platform detection) use go/ast parsing
- **File extractors** (Dockerfiles, Helm charts, go.mod) parse specialized formats

Extractors are designed to be resilient: if a file doesn't match the expected format, the extractor skips it and logs a warning instead of failing.

## Understanding the output

### component-architecture.json

The core data structure containing all extracted information:

```json
{
  "component": "my-operator",
  "repo": "github.com/org/my-operator",
  "extracted_at": "2026-04-14T10:30:00Z",
  "analyzer_version": "0.2.0",
  "crds": [...],
  "rbac": { "cluster_roles": [...], "role_bindings": [...] },
  "services": [...],
  "deployments": [...],
  "network_policies": [...],
  "controller_watches": { ... },
  "dependencies": [...],
  "secrets": [...],
  "dockerfiles": [...],
  "helm": { ... },
  "webhooks": [...],
  "config_maps": [...],
  "http_endpoints": [...],
  "ingress_routing": [...],
  "cache_config": { ... },
  "operator_config": [...],
  "reconcile_sequences": [...],
  "prometheus_metrics": [...],
  "status_conditions": [...],
  "platform_detection": { ... }
}
```

### Diagrams

Each `.mmd` file is a self-contained Mermaid diagram. View in:

- GitHub (renders Mermaid natively in markdown)
- [Mermaid Live Editor](https://mermaid.live/)
- VS Code with Mermaid extension

### report.md

The structured markdown report contains tables for every extracted category plus cache analysis findings with severity ratings.

## Tips

- Run against a clean checkout for most accurate results
- The analyzer reads from local filesystem, so submodules need to be initialized first
- For large repos with many YAML files, extraction typically takes under 10 seconds
- The code property graph (for security scanning) adds a few seconds for tree-sitter parsing
