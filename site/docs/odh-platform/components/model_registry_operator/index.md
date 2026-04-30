# model-registry-operator

> **Architecture snapshot: 2026-04-30** (2026-04-30)


**Repository:** opendatahub-io/model-registry-operator  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-04-30T15:06:03Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 2 |
| Deployments | 8 |
| Services | 1 |
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
        dep_1["controller-manager"]
        class dep_1 controller
        dep_2["controller-manager"]
        class dep_2 controller
        dep_3["controller-manager"]
        class dep_3 controller
        dep_4["controller-manager"]
        class dep_4 controller
        dep_5["controller-manager"]
        class dep_5 controller
        dep_6["controller-manager"]
        class dep_6 controller
        dep_7["controller-manager"]
        class dep_7 controller
        dep_8["controller-manager"]
        class dep_8 controller
    end

    crd_ModelRegistry{{"ModelRegistry\nmodelregistry.opendatahub.io/v1alpha1"}}
    class crd_ModelRegistry crd
    crd_ModelRegistry -->|"For (reconciles)"| controller
    crd_ModelRegistry{{"ModelRegistry\nmodelregistry.opendatahub.io/v1beta1"}}
    class crd_ModelRegistry crd
    crd_ModelRegistry -->|"For (reconciles)"| controller
    controller -->|"Owns"| owned_9["Deployment"]
    class owned_9 owned
    controller -->|"Owns"| owned_10["NetworkPolicy"]
    class owned_10 owned
    controller -->|"Owns"| owned_11["Role"]
    class owned_11 owned
    controller -->|"Owns"| owned_12["RoleBinding"]
    class owned_12 owned
    controller -->|"Owns"| owned_13["Service"]
    class owned_13 owned
    controller -->|"Owns"| owned_14["ServiceAccount"]
    class owned_14 owned
    watch_15["ClusterRoleBinding"] -->|"Watches"| controller
    class watch_15 external
    watch_16["ConfigMap"] -->|"Watches"| controller
    class watch_16 external
    watch_17["Deployment"] -->|"Watches"| controller
    class watch_17 external
    watch_18["NetworkPolicy"] -->|"Watches"| controller
    class watch_18 external
    watch_19["PersistentVolumeClaim"] -->|"Watches"| controller
    class watch_19 external
    watch_20["Role"] -->|"Watches"| controller
    class watch_20 external
    watch_21["RoleBinding"] -->|"Watches"| controller
    class watch_21 external
    watch_22["Secret"] -->|"Watches"| controller
    class watch_22 external
    watch_23["Service"] -->|"Watches"| controller
    class watch_23 external
    watch_24["ServiceAccount"] -->|"Watches"| controller
    class watch_24 external
```

### CRDs

| Group | Version | Kind | Scope | Fields | Validation Rules | Source |
|-------|---------|------|-------|--------|------------------|--------|
| modelregistry.opendatahub.io | v1alpha1 | ModelRegistry | Namespaced | 120 | 2 | [`config/crd/bases/modelregistry.opendatahub.io_modelregistries.yaml`](https://github.com/opendatahub-io/model-registry-operator/blob/d56c75fadb1ee4aa2b162859055bf91734084a03/config/crd/bases/modelregistry.opendatahub.io_modelregistries.yaml) |
| modelregistry.opendatahub.io | v1beta1 | ModelRegistry | Namespaced | 113 | 6 | [`config/crd/bases/modelregistry.opendatahub.io_modelregistries.yaml`](https://github.com/opendatahub-io/model-registry-operator/blob/d56c75fadb1ee4aa2b162859055bf91734084a03/config/crd/bases/modelregistry.opendatahub.io_modelregistries.yaml) |

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

