---
hide:
  - navigation
  - toc
---

<div style="text-align: center; padding: 40px 0 20px;">
  <img src="images/logo.svg" alt="Architecture Analyzer" style="width: 180px; height: 180px; margin-bottom: 16px;">
  <h1 style="margin: 0; font-size: 2.4em;">Architecture Analyzer</h1>
  <p style="font-size: 1.3em; color: #666; margin-top: 8px;">
    Static analysis for Kubernetes/OpenShift architecture.<br>
    25 extractors. 7 renderers. Code property graph with security queries.
  </p>
  <p style="margin-top: 20px;">
    <a href="getting-started/installation/" class="md-button md-button--primary">Get Started</a>
    <a href="https://github.com/ugiordan/architecture-analyzer" class="md-button">GitHub</a>
  </p>
</div>

## What Is This?

A Go-based static analysis tool that extracts architecture data from Kubernetes/OpenShift component repositories and produces diagrams, security reports, and code property graphs. Works with any Go-based K8s operator ecosystem. Currently deployed for OpenShift AI (RHOAI) and Open Data Hub (ODH) analysis.

Zero LLM involvement. Deterministic, reproducible, and free to run.

## Architecture

```mermaid
graph LR
    subgraph Inputs
        REPO[Git Repository]
    end

    subgraph "Extractors (22)"
        E1[CRDs]
        E2[RBAC]
        E3[Services]
        E4[Deployments]
        E5[Network Policies]
        E6[Controller Watches]
        E7[Dependencies]
        E8[Secrets]
        E9[Helm]
        E10[Dockerfiles]
        E11[Webhooks]
        E12[ConfigMaps]
        E13[HTTP Endpoints]
        E14[Ingress]
        E15[External Connections]
        E16[Feature Gates]
        E17[Cache Config]
        E18[Operator Config]
        E19[Reconcile Sequences]
        E20[Prometheus Metrics]
        E21[Status Conditions]
        E22[Platform Detection]
    end

    subgraph Data
        JSON[component-architecture.json]
    end

    subgraph "Renderers (7)"
        R1[Mermaid RBAC]
        R2[Mermaid Component]
        R3[ASCII Security]
        R4[Mermaid Dependencies]
        R5[C4 DSL]
        R6[Mermaid Dataflow]
        R7[Markdown Report]
    end

    subgraph "Code Graph"
        CPG["Multi-Language CPG\n(Go, Python, TS, Rust)"]
        TAINT[Taint Propagation]
        SEC["20 Domain Queries"]
    end

    subgraph Aggregator
        AGG[Platform Aggregator]
    end

    REPO --> E1 & E2 & E3 & E4 & E5 & E6 & E7 & E8 & E9 & E10 & E11 & E12 & E13 & E14 & E15 & E16 & E17 & E18 & E19 & E20 & E21 & E22
    E1 & E2 & E3 & E4 & E5 & E6 & E7 & E8 & E9 & E10 & E11 & E12 & E13 & E14 & E15 & E16 & E17 & E18 & E19 & E20 & E21 & E22 --> JSON
    JSON --> R1 & R2 & R3 & R4 & R5 & R6 & R7
    JSON --> AGG
    REPO --> CPG --> TAINT --> SEC

    classDef extractor fill:#3498db,stroke:#2980b9,color:#fff
    classDef renderer fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef data fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef agg fill:#f39c12,stroke:#e67e22,color:#fff
    classDef cpg fill:#9b59b6,stroke:#8e44ad,color:#fff

    class E1,E2,E3,E4,E5,E6,E7,E8,E9,E10,E11,E12,E13,E14,E15,E16,E17,E18,E19,E20,E21,E22 extractor
    class R1,R2,R3,R4,R5,R6,R7 renderer
    class JSON data
    class AGG agg
    class CPG,TAINT,SEC cpg
```

## Key Features

<div class="grid cards" markdown>

- **22 Architecture Extractors**

    ---

    CRDs, RBAC, deployments, services, network policies, controller watches, dependencies, secrets, Helm charts, Dockerfiles, webhooks, ConfigMaps, HTTP endpoints, ingress, external connections (database, gRPC, messaging), feature gates, cache architecture, operator config constants, reconciliation sequences, Prometheus metrics, status conditions, and platform detection.

    [:octicons-arrow-right-24: Extractors reference](reference/extractors.md)

- **Code Property Graph**

    ---

    Multi-language CPG (Go, Python, TypeScript, Rust) via tree-sitter. Typed nodes with edge confidence, intraprocedural data flow, control flow graphs, two-phase taint propagation, SARIF ingestion, and structural diff. 20 security queries across 3 domains (security, testing, upgrade).

    [:octicons-arrow-right-24: CPG architecture](architecture/cpg.md)

- **OOM Risk Detection**

    ---

    Cross-references controller-runtime cache config against watches and deployment memory limits. Catches real production bugs.

    [:octicons-arrow-right-24: Cache analysis](architecture/cache-analysis.md)

- **CRD Contract Validation**

    ---

    Detects breaking schema changes across repos. Runs on every PR that modifies CRD definitions.

    [:octicons-arrow-right-24: CRD validation guide](guides/crd-validation.md)

</div>

## Output Formats

| Format | File | Description |
|--------|------|-------------|
| Mermaid RBAC | `rbac.mmd` | ServiceAccounts, bindings, roles, resources |
| Mermaid Component | `component.mmd` | CRDs watched/owned, dependencies |
| ASCII Security | `security-network.txt` | Layered network, RBAC, secrets view |
| Mermaid Dependencies | `dependencies.mmd` | Go module graph (internal deps highlighted) |
| C4 DSL | `c4-context.dsl` | Structurizr C4 context diagram |
| Mermaid Dataflow | `dataflow.mmd` | Controller watches and service connections |
| Markdown Report | `report.md` | Structured tables for all extracted data |
| JSON | `component-architecture.json` | Machine-readable extracted data |
| Code Graph | `code-graph.json` | CPG nodes, edges, basic blocks, taint findings |
| Security JSON | `security-findings.json` | Domain query findings with file:line references |
| SARIF | `findings.sarif` | Security findings in SARIF 2.1.0 format |

## Real-world impact

The cache analysis has caught real production bugs:

- [opendatahub-io/data-science-pipelines-operator#992](https://github.com/opendatahub-io/data-science-pipelines-operator/issues/992): OOM from cluster-wide informers
- [opendatahub-io/model-registry-operator#457](https://github.com/opendatahub-io/model-registry-operator/issues/457): Missing cache filters on watched types
