# spark-operator

> **Architecture snapshot: 2026-05-05** (2026-05-05)


**Repository:** kubeflow/spark-operator  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-05T15:10:40Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 3 |
| Deployments | 3 |
| Services | 1 |
| Secrets | 1 |
| Cluster Roles | 5 |
| Controller Watches | 1 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for spark-operator

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["spark-operator Controller"]
        dep_1["controller-manager"]
        class dep_1 controller
        dep_2["spark-operator-controller"]
        class dep_2 controller
        dep_3["spark-operator-webhook"]
        class dep_3 controller
    end

    crd_SparkConnect{{"SparkConnect\nsparkoperator.k8s.io/v1alpha1"}}
    class crd_SparkConnect crd
    crd_SparkConnect -->|"For (reconciles)"| controller
    crd_ScheduledSparkApplication{{"ScheduledSparkApplication\nsparkoperator.k8s.io/v1beta2"}}
    class crd_ScheduledSparkApplication crd
    crd_SparkApplication{{"SparkApplication\nsparkoperator.k8s.io/v1beta2"}}
    class crd_SparkApplication crd
```

### CRDs

| Group | Version | Kind | Scope | Fields | Validation Rules | Source |
|-------|---------|------|-------|--------|------------------|--------|
| sparkoperator.k8s.io | v1alpha1 | SparkConnect | Namespaced | 95 | 0 | [`config/crd/bases/sparkoperator.k8s.io_sparkconnects.yaml`](https://github.com/kubeflow/spark-operator/blob/b8a995788a0bd700354170600d0813db8a597241/config/crd/bases/sparkoperator.k8s.io_sparkconnects.yaml) |
| sparkoperator.k8s.io | v1beta2 | ScheduledSparkApplication | Namespaced | 1676 | 0 | [`config/crd/bases/sparkoperator.k8s.io_scheduledsparkapplications.yaml`](https://github.com/kubeflow/spark-operator/blob/b8a995788a0bd700354170600d0813db8a597241/config/crd/bases/sparkoperator.k8s.io_scheduledsparkapplications.yaml) |
| sparkoperator.k8s.io | v1beta2 | SparkApplication | Namespaced | 1679 | 0 | [`config/crd/bases/sparkoperator.k8s.io_sparkapplications.yaml`](https://github.com/kubeflow/spark-operator/blob/b8a995788a0bd700354170600d0813db8a597241/config/crd/bases/sparkoperator.k8s.io_sparkapplications.yaml) |

## Dependencies

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.3 |
| github.com/prometheus/client_golang | v1.23.2 |
| k8s.io/api | v0.32.5 |
| k8s.io/apiextensions-apiserver | v0.32.5 |
| k8s.io/apimachinery | v0.32.5 |
| k8s.io/apiserver | v0.32.5 |
| k8s.io/client-go | v0.32.5 |
| sigs.k8s.io/controller-runtime | v0.20.4 |

