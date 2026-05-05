# opendatahub-operator

> **Architecture snapshot: 2026-05-05** (2026-05-05)


**Repository:** opendatahub-io/opendatahub-operator  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-05T15:10:47Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 6 |
| Deployments | 3 |
| Services | 0 |
| Secrets | 0 |
| Cluster Roles | 7 |
| Controller Watches | 28 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for opendatahub-operator

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["opendatahub-operator Controller"]
        dep_1["controller-manager"]
        class dep_1 controller
        dep_2["controller-manager"]
        class dep_2 controller
        dep_3["controller-manager"]
        class dep_3 controller
    end

    crd_OdhQuickStart{{"OdhQuickStart\nconsole.openshift.io/v1"}}
    class crd_OdhQuickStart crd
    crd_OdhApplication{{"OdhApplication\ndashboard.opendatahub.io/v1"}}
    class crd_OdhApplication crd
    crd_OdhDocument{{"OdhDocument\ndashboard.opendatahub.io/v1"}}
    class crd_OdhDocument crd
    crd_DataScienceCluster{{"DataScienceCluster\ndatasciencecluster.opendatahub.io/v1alpha1"}}
    class crd_DataScienceCluster crd
    crd_DataScienceCluster -->|"For (reconciles)"| controller
    crd_DSCInitialization{{"DSCInitialization\ndscinitialization.opendatahub.io/v1alpha1"}}
    class crd_DSCInitialization crd
    crd_DSCInitialization -->|"For (reconciles)"| controller
    crd_OdhDashboardConfig{{"OdhDashboardConfig\nopendatahub.io/v1alpha"}}
    class crd_OdhDashboardConfig crd
    controller -->|"Owns"| owned_4["ClusterRole"]
    class owned_4 owned
    controller -->|"Owns"| owned_5["ClusterRoleBinding"]
    class owned_5 owned
    controller -->|"Owns"| owned_6["ConfigMap"]
    class owned_6 owned
    controller -->|"Owns"| owned_7["Deployment"]
    class owned_7 owned
    controller -->|"Owns"| owned_8["Namespace"]
    class owned_8 owned
    controller -->|"Owns"| owned_9["NetworkPolicy"]
    class owned_9 owned
    controller -->|"Owns"| owned_10["Pod"]
    class owned_10 owned
    controller -->|"Owns"| owned_11["ReplicaSet"]
    class owned_11 owned
    controller -->|"Owns"| owned_12["Role"]
    class owned_12 owned
    controller -->|"Owns"| owned_13["RoleBinding"]
    class owned_13 owned
    controller -->|"Owns"| owned_14["Secret"]
    class owned_14 owned
    controller -->|"Owns"| owned_15["Service"]
    class owned_15 owned
    controller -->|"Owns"| owned_16["ServiceAccount"]
    class owned_16 owned
    watch_17["Secret"] -->|"Watches"| controller
    class watch_17 external
    controller -.->|"depends on"| odh_18["opendatahub-operator"]
    class odh_18 dep
```

### CRDs

| Group | Version | Kind | Scope | Fields | Validation Rules | Source |
|-------|---------|------|-------|--------|------------------|--------|
| console.openshift.io | v1 | OdhQuickStart | Namespaced | 31 | 0 | [`config/crd/bases/odhquickstarts.console.openshift.io_odhquickstarts.yaml`](https://github.com/opendatahub-io/opendatahub-operator/blob/fc3568b08335435af8f5ca135376f7793c260b43/config/crd/bases/odhquickstarts.console.openshift.io_odhquickstarts.yaml) |
| dashboard.opendatahub.io | v1 | OdhApplication | Namespaced | 52 | 0 | [`config/crd/bases/odhapplications.dashboard.opendatahub.io_odhapplications.yaml`](https://github.com/opendatahub-io/opendatahub-operator/blob/fc3568b08335435af8f5ca135376f7793c260b43/config/crd/bases/odhapplications.dashboard.opendatahub.io_odhapplications.yaml) |
| dashboard.opendatahub.io | v1 | OdhDocument | Namespaced | 16 | 0 | [`config/crd/bases/odhdocuments.dashboard.opendatahub.io_odhdocuments.yaml`](https://github.com/opendatahub-io/opendatahub-operator/blob/fc3568b08335435af8f5ca135376f7793c260b43/config/crd/bases/odhdocuments.dashboard.opendatahub.io_odhdocuments.yaml) |
| datasciencecluster.opendatahub.io | v1alpha1 | DataScienceCluster | Cluster | 38 | 0 | [`config/crd/bases/datasciencecluster.opendatahub.io_datascienceclusters.yaml`](https://github.com/opendatahub-io/opendatahub-operator/blob/fc3568b08335435af8f5ca135376f7793c260b43/config/crd/bases/datasciencecluster.opendatahub.io_datascienceclusters.yaml) |
| dscinitialization.opendatahub.io | v1alpha1 | DSCInitialization | Cluster | 27 | 0 | [`config/crd/bases/dscinitialization.opendatahub.io_dscinitializations.yaml`](https://github.com/opendatahub-io/opendatahub-operator/blob/fc3568b08335435af8f5ca135376f7793c260b43/config/crd/bases/dscinitialization.opendatahub.io_dscinitializations.yaml) |
| opendatahub.io | v1alpha | OdhDashboardConfig | Namespaced | 54 | 0 | [`config/crd/bases/odhdashboardconfigs.opendatahub.io_odhdashboardconfigs.yaml`](https://github.com/opendatahub-io/opendatahub-operator/blob/fc3568b08335435af8f5ca135376f7793c260b43/config/crd/bases/odhdashboardconfigs.opendatahub.io_odhdashboardconfigs.yaml) |

## Dependencies

### Internal Platform Dependencies

| Component | Interaction |
|-----------|-------------|
| opendatahub-operator | Go module dependency: github.com/opendatahub-io/opendatahub-operator |

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.2.4 |
| github.com/operator-framework/api | v0.17.6 |
| k8s.io/api | v0.26.0 |
| k8s.io/apiextensions-apiserver | v0.27.2 |
| k8s.io/apimachinery | v0.27.2 |
| k8s.io/client-go | v0.26.0 |
| sigs.k8s.io/controller-runtime | v0.14.4 |

