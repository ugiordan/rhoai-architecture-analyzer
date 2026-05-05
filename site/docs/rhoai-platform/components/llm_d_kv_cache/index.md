# llm-d-kv-cache

> **Architecture snapshot: 2026-05-05** (2026-05-05)


**Repository:** llm-d/llm-d-kv-cache  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-05T15:10:13Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 0 |
| Deployments | 1 |
| Services | 0 |
| Secrets | 0 |
| Cluster Roles | 0 |
| Controller Watches | 0 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for llm-d-kv-cache

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["llm-d-kv-cache Controller"]
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
| github.com/go-logr/logr | v1.4.3 |
| github.com/prometheus/client_golang | v1.22.0 |
| github.com/prometheus/client_model | v0.6.1 |
| google.golang.org/grpc | v1.77.0 |
| k8s.io/api | v0.33.8 |
| k8s.io/apimachinery | v0.33.8 |
| k8s.io/client-go | v0.33.8 |
| sigs.k8s.io/controller-runtime | v0.21.0 |

