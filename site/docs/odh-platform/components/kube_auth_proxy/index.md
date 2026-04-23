# kube-auth-proxy

> **Architecture snapshot: 2026-04-23** (2026-04-23)


**Repository:** opendatahub-io/kube-auth-proxy  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-04-23T08:07:18Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 0 |
| Deployments | 29 |
| Services | 11 |
| Secrets | 2 |
| Cluster Roles | 0 |
| Controller Watches | 0 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for kube-auth-proxy

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["kube-auth-proxy Controller"]
        dep_1["example-app"]
        class dep_1 controller
        dep_2["example-app"]
        class dep_2 controller
        dep_3["example-app"]
        class dep_3 controller
        dep_4["kube-auth-proxy"]
        class dep_4 controller
        dep_5["kube-auth-proxy"]
        class dep_5 controller
        dep_6["kube-auth-proxy"]
        class dep_6 controller
        dep_7["kube-rbac-proxy"]
        class dep_7 controller
        dep_8["kube-rbac-proxy"]
        class dep_8 controller
        dep_9["kube-rbac-proxy"]
        class dep_9 controller
        dep_10["kube-rbac-proxy"]
        class dep_10 controller
        dep_11["kube-rbac-proxy"]
        class dep_11 controller
        dep_12["kube-rbac-proxy"]
        class dep_12 controller
        dep_13["kube-rbac-proxy"]
        class dep_13 controller
        dep_14["kube-rbac-proxy"]
        class dep_14 controller
        dep_15["kube-rbac-proxy"]
        class dep_15 controller
        dep_16["kube-rbac-proxy"]
        class dep_16 controller
        dep_17["kube-rbac-proxy"]
        class dep_17 controller
        dep_18["kube-rbac-proxy"]
        class dep_18 controller
        dep_19["kube-rbac-proxy"]
        class dep_19 controller
        dep_20["kube-rbac-proxy"]
        class dep_20 controller
        dep_21["kube-rbac-proxy"]
        class dep_21 controller
        dep_22["kube-rbac-proxy"]
        class dep_22 controller
        dep_23["kube-rbac-proxy"]
        class dep_23 controller
        dep_24["kube-rbac-proxy"]
        class dep_24 controller
        dep_25["kube-rbac-proxy"]
        class dep_25 controller
        dep_26["kube-rbac-proxy"]
        class dep_26 controller
        dep_27["kube-rbac-proxy"]
        class dep_27 controller
        dep_28["kube-rbac-proxy"]
        class dep_28 controller
        dep_29["kube-rbac-proxy-verb-override"]
        class dep_29 controller
    end
```

### CRDs

No CRDs defined.

## Dependencies

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/prometheus/client_golang | v1.23.2 |
| k8s.io/api | v0.35.3 |
| k8s.io/api | v0.35.3 |
| k8s.io/apimachinery | v0.35.3 |
| k8s.io/apimachinery | v0.35.3 |
| k8s.io/apiserver | v0.35.3 |
| k8s.io/client-go | v0.35.3 |
| k8s.io/client-go | v0.35.3 |

