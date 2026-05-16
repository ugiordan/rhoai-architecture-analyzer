# llm-d-routing-sidecar

> **Architecture snapshot: 2026-05-16** (2026-05-16)


**Repository:** llm-d/llm-d-routing-sidecar  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-16T03:47:38Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 0 |
| Deployments | 1 |
| Services | 1 |
| Secrets | 0 |
| Cluster Roles | 0 |
| Controller Watches | 0 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for llm-d-routing-sidecar

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["llm-d-routing-sidecar Controller"]
        dep_1["0"]
        class dep_1 controller
    end
```

### CRDs

No CRDs defined.

## Dependencies

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.2 |
| github.com/go-logr/logr | v1.4.2 |
| github.com/go-logr/logr | v1.4.1 |
| github.com/go-logr/logr | v1.4.2 |
| github.com/go-logr/logr | v1.4.1 |
| k8s.io/api | v0.31.3 |
| k8s.io/api | v0.31.3 |
| k8s.io/apimachinery | v0.31.3 |
| k8s.io/apimachinery | v0.31.3 |
| k8s.io/apimachinery | v0.31.3 |
| k8s.io/client-go | v0.31.3 |

