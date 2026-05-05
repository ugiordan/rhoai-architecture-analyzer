# odh-model-controller

> **Architecture snapshot: 2026-05-05** (2026-05-05)


**Repository:** opendatahub-io/odh-model-controller  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-05T15:10:49Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 1 |
| Deployments | 2 |
| Services | 2 |
| Secrets | 2 |
| Cluster Roles | 7 |
| Controller Watches | 41 |

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
    controller -->|"Owns"| owned_3["AuthPolicy"]
    class owned_3 owned
    controller -->|"Owns"| owned_4["ClusterRoleBinding"]
    class owned_4 owned
    controller -->|"Owns"| owned_5["ConfigMap"]
    class owned_5 owned
    controller -->|"Owns"| owned_6["EnvoyFilter"]
    class owned_6 owned
    controller -->|"Owns"| owned_7["Namespace"]
    class owned_7 owned
    controller -->|"Owns"| owned_8["NetworkPolicy"]
    class owned_8 owned
    controller -->|"Owns"| owned_9["PodMonitor"]
    class owned_9 owned
    controller -->|"Owns"| owned_10["Role"]
    class owned_10 owned
    controller -->|"Owns"| owned_11["RoleBinding"]
    class owned_11 owned
    controller -->|"Owns"| owned_12["Route"]
    class owned_12 owned
    controller -->|"Owns"| owned_13["Secret"]
    class owned_13 owned
    controller -->|"Owns"| owned_14["Service"]
    class owned_14 owned
    controller -->|"Owns"| owned_15["ServiceAccount"]
    class owned_15 owned
    controller -->|"Owns"| owned_16["ServiceMonitor"]
    class owned_16 owned
    controller -->|"Owns"| owned_17["ServingRuntime"]
    class owned_17 owned
    controller -->|"Owns"| owned_18["Template"]
    class owned_18 owned
    controller -->|"Owns"| owned_19["TriggerAuthentication"]
    class owned_19 owned
    watch_20["AuthPolicy"] -->|"Watches"| controller
    class watch_20 external
    watch_21["Authorino"] -->|"Watches"| controller
    class watch_21 external
    watch_22["ConfigMap"] -->|"Watches"| controller
    class watch_22 external
    watch_23["Kuadrant"] -->|"Watches"| controller
    class watch_23 external
    watch_24["LLMInferenceService"] -->|"Watches"| controller
    class watch_24 external
    watch_25["LLMInferenceServiceConfig"] -->|"Watches"| controller
    class watch_25 external
    watch_26["Namespace"] -->|"Watches"| controller
    class watch_26 external
    watch_27["RoleBinding"] -->|"Watches"| controller
    class watch_27 external
    watch_28["Secret"] -->|"Watches"| controller
    class watch_28 external
    watch_29["ServingRuntime"] -->|"Watches"| controller
    class watch_29 external
    controller -.->|"depends on"| odh_30["kserve"]
    class odh_30 dep
```

### CRDs

| Group | Version | Kind | Scope | Fields | Validation Rules | Source |
|-------|---------|------|-------|--------|------------------|--------|
| nim.opendatahub.io | v1 | Account | Namespaced | 57 | 0 | [`config/crd/bases/nim.opendatahub.io_accounts.yaml`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/config/crd/bases/nim.opendatahub.io_accounts.yaml) |

## Dependencies

### Internal Platform Dependencies

| Component | Interaction |
|-----------|-------------|
| kserve | Go module dependency: github.com/opendatahub-io/kserve |

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring | v0.89.0 |
| github.com/prometheus/client_golang | v1.23.2 |
| k8s.io/api | v0.34.3 |
| k8s.io/apiextensions-apiserver | v0.34.3 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apiserver | v0.34.3 |
| k8s.io/client-go | v0.34.3 |
| sigs.k8s.io/controller-runtime | v0.19.7 |

