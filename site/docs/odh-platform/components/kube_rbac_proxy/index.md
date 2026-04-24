# kube-rbac-proxy

> **Architecture snapshot: 2026-04-24** (2026-04-24)


**Repository:** brancz/kube-rbac-proxy  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-04-24T08:14:03Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 0 |
| Deployments | 7 |
| Services | 1 |
| Secrets | 0 |
| Cluster Roles | 0 |
| Controller Watches | 0 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for kube-rbac-proxy

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["kube-rbac-proxy Controller"]
        dep_1["kube-rbac-proxy"]
        class dep_1 controller
        dep_2["kube-rbac-proxy"]
        class dep_2 controller
        dep_3["kube-rbac-proxy"]
        class dep_3 controller
        dep_4["kube-rbac-proxy"]
        class dep_4 controller
        dep_5["kube-rbac-proxy"]
        class dep_5 controller
        dep_6["kube-rbac-proxy"]
        class dep_6 controller
        dep_7["kube-rbac-proxy"]
        class dep_7 controller
    end
```

### CRDs

No CRDs defined.

## Dependencies

### Key External Dependencies

| Module | Version |
|--------|---------|
| k8s.io/api | v0.35.4 |
| k8s.io/apimachinery | v0.35.4 |
| k8s.io/apiserver | v0.35.4 |
| k8s.io/client-go | v0.35.4 |

