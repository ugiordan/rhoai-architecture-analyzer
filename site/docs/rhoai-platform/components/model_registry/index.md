# model-registry

> **Architecture snapshot: 2026-05-05** (2026-05-05)


**Repository:** kubeflow/model-registry  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-05T15:11:08Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 0 |
| Deployments | 3 |
| Services | 1 |
| Secrets | 3 |
| Cluster Roles | 6 |
| Controller Watches | 2 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for model-registry

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["model-registry Controller"]
        dep_1["controller-manager"]
        class dep_1 controller
        dep_2["minio"]
        class dep_2 controller
        dep_3["model-catalog-server"]
        class dep_3 controller
    end
```

### CRDs

No CRDs defined.

## Dependencies

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.3 |
| k8s.io/api | v0.34.4 |
| k8s.io/api | v0.35.4 |
| k8s.io/apimachinery | v0.35.4 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/apimachinery | v0.35.4 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/client-go | v0.34.4 |
| k8s.io/client-go | v0.35.3 |
| k8s.io/client-go | v0.34.4 |
| sigs.k8s.io/controller-runtime | v0.23.3 |
| sigs.k8s.io/controller-runtime | v0.22.4 |
| sigs.k8s.io/controller-runtime | v0.22.4 |

