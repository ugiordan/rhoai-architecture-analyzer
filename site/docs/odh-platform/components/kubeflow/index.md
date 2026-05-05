# kubeflow

> **Architecture snapshot: 2026-05-05** (2026-05-05)


**Repository:** opendatahub-io/kubeflow  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-05T15:09:07Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 0 |
| Deployments | 7 |
| Services | 3 |
| Secrets | 2 |
| Cluster Roles | 0 |
| Controller Watches | 14 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for kubeflow

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["kubeflow Controller"]
        dep_1["controller-manager"]
        class dep_1 controller
        dep_2["controller-manager"]
        class dep_2 controller
        dep_3["controller-manager"]
        class dep_3 controller
        dep_4["controller-manager"]
        class dep_4 controller
        dep_5["deployment"]
        class dep_5 controller
        dep_6["deployment"]
        class dep_6 controller
        dep_7["manager"]
        class dep_7 controller
    end

    controller -->|"Owns"| owned_8["ConfigMap"]
    class owned_8 owned
    controller -->|"Owns"| owned_9["NetworkPolicy"]
    class owned_9 owned
    controller -->|"Owns"| owned_10["RoleBinding"]
    class owned_10 owned
    controller -->|"Owns"| owned_11["Secret"]
    class owned_11 owned
    controller -->|"Owns"| owned_12["Service"]
    class owned_12 owned
    controller -->|"Owns"| owned_13["ServiceAccount"]
    class owned_13 owned
    controller -->|"Owns"| owned_14["StatefulSet"]
    class owned_14 owned
    watch_15["ConfigMap"] -->|"Watches"| controller
    class watch_15 external
    watch_16["HTTPRoute"] -->|"Watches"| controller
    class watch_16 external
    watch_17["ReferenceGrant"] -->|"Watches"| controller
    class watch_17 external
    controller -.->|"depends on"| odh_18["data-science-pipelines-operator"]
    class odh_18 dep
```

### CRDs

No CRDs defined.

## Dependencies

### Internal Platform Dependencies

| Component | Interaction |
|-----------|-------------|
| data-science-pipelines-operator | Go module dependency: github.com/opendatahub-io/data-science-pipelines-operator |

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/prometheus/client_golang | v1.23.2 |
| k8s.io/api | v0.33.7 |
| k8s.io/api | v0.33.7 |
| k8s.io/apiextensions-apiserver | v0.33.7 |
| k8s.io/apimachinery | v0.33.7 |
| k8s.io/apimachinery | v0.33.7 |
| k8s.io/client-go | v0.33.7 |
| k8s.io/client-go | v0.33.7 |
| sigs.k8s.io/controller-runtime | v0.21.0 |
| sigs.k8s.io/controller-runtime | v0.21.0 |

