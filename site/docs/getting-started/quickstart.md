# Quick Start

After [installing](installation.md) the analyzer, here's how to analyze your first repository.

## Analyze a repository

```bash
arch-analyzer analyze /path/to/your/repo --output-dir output/
```

This runs all 22 extractors and all 7 renderers, producing:

```
output/
  component-architecture.json    # Extracted architecture data
  diagrams/
    rbac.mmd                     # Mermaid RBAC graph
    component.mmd                # Mermaid component diagram
    dependencies.mmd             # Mermaid dependency graph
    dataflow.mmd                 # Mermaid sequence diagram
    security-network.txt         # ASCII security/network view
    c4-context.dsl               # Structurizr C4 diagram
    report.md                    # Structured markdown report
```

## Run a security scan

```bash
arch-analyzer scan /path/to/your/repo --format sarif --output findings.sarif
```

Builds a code property graph from Go source files and runs security queries: taint analysis, SQL injection detection, hardcoded secrets, missing authentication.

## Full analysis (everything at once)

```bash
arch-analyzer full-analysis /path/to/your/repo --output-dir output/
```

Runs architecture extraction, diagram rendering, code graph scanning, and schema extraction in one pass.

## View the results

### Mermaid diagrams

Open any `.mmd` file in a Mermaid-compatible viewer or paste into the [Mermaid Live Editor](https://mermaid.live/).

### C4 diagrams

Load `c4-context.dsl` into [Structurizr](https://structurizr.com/) or the VS Code Structurizr extension.

### Markdown report

View `report.md` directly in GitHub or any markdown renderer. It contains tables for all extracted data plus cache analysis findings.

### Security findings

SARIF output can be loaded into GitHub Code Scanning, VS Code SARIF Viewer, or any SARIF-compatible tool.

## What it extracts

The analyzer reads:

- Kubernetes YAML manifests (deployments, services, RBAC, network policies, etc.)
- Go source code (controller watches, HTTP endpoints, cache config, operator constants, reconcile sequences, Prometheus metrics, status conditions, platform detection)
- Dockerfiles (base images, security settings)
- Helm charts (metadata, security defaults)
- go.mod (dependencies, internal ODH modules)

It never modifies any files. Read-only static analysis.
