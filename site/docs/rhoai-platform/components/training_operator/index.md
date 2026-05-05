# training-operator

> **Architecture snapshot: 2026-05-05** (2026-05-05)


**Repository:** kubeflow/training-operator  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-05T15:09:06Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 6 |
| Deployments | 3 |
| Services | 1 |
| Secrets | 2 |
| Cluster Roles | 6 |
| Controller Watches | 0 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for training-operator

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["training-operator Controller"]
        dep_1["training-operator"]
        class dep_1 controller
        dep_2["training-operator"]
        class dep_2 controller
        dep_3["training-operator"]
        class dep_3 controller
    end

    crd_JAXJob{{"JAXJob\nkubeflow.org/v1"}}
    class crd_JAXJob crd
    crd_MPIJob{{"MPIJob\nkubeflow.org/v1"}}
    class crd_MPIJob crd
    crd_PaddleJob{{"PaddleJob\nkubeflow.org/v1"}}
    class crd_PaddleJob crd
    crd_PyTorchJob{{"PyTorchJob\nkubeflow.org/v1"}}
    class crd_PyTorchJob crd
    crd_TFJob{{"TFJob\nkubeflow.org/v1"}}
    class crd_TFJob crd
    crd_XGBoostJob{{"XGBoostJob\nkubeflow.org/v1"}}
    class crd_XGBoostJob crd
```

### CRDs

| Group | Version | Kind | Scope | Fields | Validation Rules | Source |
|-------|---------|------|-------|--------|------------------|--------|
| kubeflow.org | v1 | JAXJob | Namespaced | 1073 | 1 | [`manifests/base/crds/kubeflow.org_jaxjobs.yaml`](https://github.com/kubeflow/training-operator/blob/8582a4b2a238e3552c6b726764580295303a3414/manifests/base/crds/kubeflow.org_jaxjobs.yaml) |
| kubeflow.org | v1 | MPIJob | Namespaced | 1076 | 1 | [`manifests/base/crds/kubeflow.org_mpijobs.yaml`](https://github.com/kubeflow/training-operator/blob/8582a4b2a238e3552c6b726764580295303a3414/manifests/base/crds/kubeflow.org_mpijobs.yaml) |
| kubeflow.org | v1 | PaddleJob | Namespaced | 1140 | 1 | [`manifests/base/crds/kubeflow.org_paddlejobs.yaml`](https://github.com/kubeflow/training-operator/blob/8582a4b2a238e3552c6b726764580295303a3414/manifests/base/crds/kubeflow.org_paddlejobs.yaml) |
| kubeflow.org | v1 | PyTorchJob | Namespaced | 1150 | 1 | [`manifests/base/crds/kubeflow.org_pytorchjobs.yaml`](https://github.com/kubeflow/training-operator/blob/8582a4b2a238e3552c6b726764580295303a3414/manifests/base/crds/kubeflow.org_pytorchjobs.yaml) |
| kubeflow.org | v1 | TFJob | Namespaced | 1075 | 1 | [`manifests/base/crds/kubeflow.org_tfjobs.yaml`](https://github.com/kubeflow/training-operator/blob/8582a4b2a238e3552c6b726764580295303a3414/manifests/base/crds/kubeflow.org_tfjobs.yaml) |
| kubeflow.org | v1 | XGBoostJob | Namespaced | 1073 | 1 | [`manifests/base/crds/kubeflow.org_xgboostjobs.yaml`](https://github.com/kubeflow/training-operator/blob/8582a4b2a238e3552c6b726764580295303a3414/manifests/base/crds/kubeflow.org_xgboostjobs.yaml) |

## Dependencies

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.2 |
| github.com/prometheus/client_golang | v1.20.2 |
| k8s.io/api | v0.31.3 |
| k8s.io/apimachinery | v0.31.3 |
| k8s.io/client-go | v0.31.3 |
| sigs.k8s.io/controller-runtime | v0.19.1 |

