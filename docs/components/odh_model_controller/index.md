# odh-model-controller

**Repository:** opendatahub-io/odh-model-controller  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-04-16T15:36:18Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 1 |
| Deployments | 2 |
| Services | 1 |
| Secrets | 1 |
| Cluster Roles | 7 |
| Controller Watches | 39 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for odh-model-controller

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["odh-model-controller Controller"]
        dep_1["odh-model-controller"]
        class dep_1 controller
        dep_2["odh-model-controller"]
        class dep_2 controller
    end

    crd_Account{{"Account\nnim.opendatahub.io/v1"}}
    class crd_Account crd
    crd_Account -->|"For (reconciles)"| controller
    controller -->|"Owns"| owned_3["ClusterRoleBinding"]
    class owned_3 owned
    controller -->|"Owns"| owned_4["ConfigMap"]
    class owned_4 owned
    controller -->|"Owns"| owned_5["Namespace"]
    class owned_5 owned
    controller -->|"Owns"| owned_6["NetworkPolicy"]
    class owned_6 owned
    controller -->|"Owns"| owned_7["PodMonitor"]
    class owned_7 owned
    controller -->|"Owns"| owned_8["Role"]
    class owned_8 owned
    controller -->|"Owns"| owned_9["RoleBinding"]
    class owned_9 owned
    controller -->|"Owns"| owned_10["Route"]
    class owned_10 owned
    controller -->|"Owns"| owned_11["Secret"]
    class owned_11 owned
    controller -->|"Owns"| owned_12["Service"]
    class owned_12 owned
    controller -->|"Owns"| owned_13["ServiceAccount"]
    class owned_13 owned
    controller -->|"Owns"| owned_14["ServiceMonitor"]
    class owned_14 owned
    controller -->|"Owns"| owned_15["ServingRuntime"]
    class owned_15 owned
    controller -->|"Owns"| owned_16["Template"]
    class owned_16 owned
    controller -->|"Owns"| owned_17["TriggerAuthentication"]
    class owned_17 owned
    watch_18["AuthPolicy"] -->|"Watches"| controller
    class watch_18 external
    watch_19["Authorino"] -->|"Watches"| controller
    class watch_19 external
    watch_20["ConfigMap"] -->|"Watches"| controller
    class watch_20 external
    watch_21["EnvoyFilter"] -->|"Watches"| controller
    class watch_21 external
    watch_22["Gateway"] -->|"Watches"| controller
    class watch_22 external
    watch_23["Kuadrant"] -->|"Watches"| controller
    class watch_23 external
    watch_24["Namespace"] -->|"Watches"| controller
    class watch_24 external
    watch_25["Role"] -->|"Watches"| controller
    class watch_25 external
    watch_26["RoleBinding"] -->|"Watches"| controller
    class watch_26 external
    watch_27["Secret"] -->|"Watches"| controller
    class watch_27 external
    watch_28["ServingRuntime"] -->|"Watches"| controller
    class watch_28 external
    controller -.->|"depends on"| odh_29["kserve"]
    class odh_29 dep
```

### CRDs

| Group | Version | Kind | Scope | Fields | Validation Rules | Source |
|-------|---------|------|-------|--------|------------------|--------|
| nim.opendatahub.io | v1 | Account | Namespaced | 57 | 0 | `config/crd/bases/nim.opendatahub.io_accounts.yaml` |

## Dependencies

### Internal RHOAI Dependencies

| Component | Interaction |
|-----------|-------------|
| kserve | Go module dependency: github.com/opendatahub-io/kserve |

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring | v0.76.2 |
| k8s.io/api | v0.33.1 |
| k8s.io/apiextensions-apiserver | v0.33.1 |
| k8s.io/apimachinery | v0.33.1 |
| k8s.io/apiserver | v0.33.1 |
| k8s.io/client-go | v0.33.1 |
| sigs.k8s.io/controller-runtime | v0.19.1 |

