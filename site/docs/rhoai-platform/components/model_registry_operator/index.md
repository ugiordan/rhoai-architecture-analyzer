# model-registry-operator

> **Architecture snapshot: 2026-05-05** (2026-05-05)


**Repository:** opendatahub-io/model-registry-operator  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-05T15:09:25Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 1 |
| Deployments | 1 |
| Services | 6 |
| Secrets | 2 |
| Cluster Roles | 6 |
| Controller Watches | 17 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for model-registry-operator

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["model-registry-operator Controller"]
        dep_1["model-registry-operator-controller-manager"]
        class dep_1 controller
    end

    crd_ModelRegistry{{"ModelRegistry\nmodelregistry.opendatahub.io/v1beta1"}}
    class crd_ModelRegistry crd
    crd_ModelRegistry -->|"For (reconciles)"| controller
    controller -->|"Owns"| owned_2["Deployment"]
    class owned_2 owned
    controller -->|"Owns"| owned_3["NetworkPolicy"]
    class owned_3 owned
    controller -->|"Owns"| owned_4["Role"]
    class owned_4 owned
    controller -->|"Owns"| owned_5["RoleBinding"]
    class owned_5 owned
    controller -->|"Owns"| owned_6["Service"]
    class owned_6 owned
    controller -->|"Owns"| owned_7["ServiceAccount"]
    class owned_7 owned
    watch_8["ClusterRoleBinding"] -->|"Watches"| controller
    class watch_8 external
    watch_9["ConfigMap"] -->|"Watches"| controller
    class watch_9 external
    watch_10["Deployment"] -->|"Watches"| controller
    class watch_10 external
    watch_11["NetworkPolicy"] -->|"Watches"| controller
    class watch_11 external
    watch_12["PersistentVolumeClaim"] -->|"Watches"| controller
    class watch_12 external
    watch_13["Role"] -->|"Watches"| controller
    class watch_13 external
    watch_14["RoleBinding"] -->|"Watches"| controller
    class watch_14 external
    watch_15["Secret"] -->|"Watches"| controller
    class watch_15 external
    watch_16["Service"] -->|"Watches"| controller
    class watch_16 external
    watch_17["ServiceAccount"] -->|"Watches"| controller
    class watch_17 external
```

### CRDs

| Group | Version | Kind | Scope | Fields | Validation Rules | Source |
|-------|---------|------|-------|--------|------------------|--------|
| modelregistry.opendatahub.io | v1beta1 | ModelRegistry | Namespaced | 113 | 6 | [`config/crd/bases/modelregistry.opendatahub.io_modelregistries.yaml`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/config/crd/bases/modelregistry.opendatahub.io_modelregistries.yaml) |

## Dependencies

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.3 |
| k8s.io/api | v0.35.4 |
| k8s.io/apiextensions-apiserver | v0.35.4 |
| k8s.io/apimachinery | v0.35.4 |
| k8s.io/client-go | v0.35.4 |
| sigs.k8s.io/controller-runtime | v0.23.3 |

