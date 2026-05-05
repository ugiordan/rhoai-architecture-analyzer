# mlflow-operator

> **Architecture snapshot: 2026-05-05** (2026-05-05)


**Repository:** opendatahub-io/mlflow-operator  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-05T15:09:47Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 2 |
| Deployments | 2 |
| Services | 3 |
| Secrets | 2 |
| Cluster Roles | 6 |
| Controller Watches | 12 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for mlflow-operator

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["mlflow-operator Controller"]
        dep_1["mlflow-operator-controller-manager"]
        class dep_1 controller
        dep_2["postgres-deployment"]
        class dep_2 controller
    end

    crd_MLflowConfig{{"MLflowConfig\nmlflow.kubeflow.org/v1"}}
    class crd_MLflowConfig crd
    crd_MLflow{{"MLflow\nmlflow.opendatahub.io/v1"}}
    class crd_MLflow crd
    crd_MLflow -->|"For (reconciles)"| controller
    controller -->|"Owns"| owned_3["ClusterRoleBinding"]
    class owned_3 owned
    controller -->|"Owns"| owned_4["ConsoleLink"]
    class owned_4 owned
    controller -->|"Owns"| owned_5["CronJob"]
    class owned_5 owned
    controller -->|"Owns"| owned_6["Deployment"]
    class owned_6 owned
    controller -->|"Owns"| owned_7["HTTPRoute"]
    class owned_7 owned
    controller -->|"Owns"| owned_8["PersistentVolumeClaim"]
    class owned_8 owned
    controller -->|"Owns"| owned_9["Secret"]
    class owned_9 owned
    controller -->|"Owns"| owned_10["Service"]
    class owned_10 owned
    controller -->|"Owns"| owned_11["ServiceAccount"]
    class owned_11 owned
    controller -->|"Owns"| owned_12["ServiceMonitor"]
    class owned_12 owned
    watch_13["ClusterRole"] -->|"Watches"| controller
    class watch_13 external
    controller -.->|"depends on"| odh_14["mlflow-operator"]
    class odh_14 dep
```

### CRDs

| Group | Version | Kind | Scope | Fields | Validation Rules | Source |
|-------|---------|------|-------|--------|------------------|--------|
| mlflow.kubeflow.org | v1 | MLflowConfig | Namespaced | 6 | 4 | [`config/crd/mlflow.kubeflow.org_mlflowconfigs.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/c47bea907957a9eeb35ee4ed3b0403e855d096cc/config/crd/mlflow.kubeflow.org_mlflowconfigs.yaml) |
| mlflow.opendatahub.io | v1 | MLflow | Cluster | 296 | 16 | [`config/crd/bases/mlflow.opendatahub.io_mlflows.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/c47bea907957a9eeb35ee4ed3b0403e855d096cc/config/crd/bases/mlflow.opendatahub.io_mlflows.yaml) |

## Dependencies

### Internal Platform Dependencies

| Component | Interaction |
|-----------|-------------|
| mlflow-operator | Go module dependency: github.com/opendatahub-io/mlflow-operator/api |

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring | v0.89.0 |
| k8s.io/api | v0.34.3 |
| k8s.io/api | v0.34.2 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apimachinery | v0.34.2 |
| k8s.io/client-go | v0.34.3 |
| sigs.k8s.io/controller-runtime | v0.22.4 |
| sigs.k8s.io/controller-runtime | v0.22.4 |

