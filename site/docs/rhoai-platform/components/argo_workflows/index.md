# argo-workflows

> **Architecture snapshot: 2026-05-05** (2026-05-05)


**Repository:** argoproj/argo-workflows  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-05T15:10:46Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 0 |
| Deployments | 0 |
| Services | 0 |
| Secrets | 0 |
| Cluster Roles | 5 |
| Controller Watches | 0 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for argo-workflows

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["argo-workflows Controller"]
        ctrl_1["Controller"]
        class ctrl_1 controller
    end
```

### CRDs

No CRDs defined.

## Dependencies

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/prometheus/client_golang | v1.21.1 |
| github.com/prometheus/common | v0.62.0 |
| google.golang.org/grpc | v1.71.1 |
| k8s.io/api | v0.32.2 |
| k8s.io/apimachinery | v0.32.2 |
| k8s.io/client-go | v0.32.2 |

