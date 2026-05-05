# trainer

> **Architecture snapshot: 2026-05-05** (2026-05-05)


**Repository:** kubeflow/trainer  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-05T15:09:15Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 3 |
| Deployments | 4 |
| Services | 0 |
| Secrets | 1 |
| Cluster Roles | 8 |
| Controller Watches | 0 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for trainer

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["trainer Controller"]
        dep_1["kubeflow-trainer-controller-manager"]
        class dep_1 controller
        dep_2["kubeflow-trainer-controller-manager"]
        class dep_2 controller
        dep_3["kubeflow-trainer-controller-manager"]
        class dep_3 controller
        dep_4["kubeflow-trainer-controller-manager"]
        class dep_4 controller
    end

    crd_ClusterTrainingRuntime{{"ClusterTrainingRuntime\ntrainer.kubeflow.org/v1alpha1"}}
    class crd_ClusterTrainingRuntime crd
    crd_TrainJob{{"TrainJob\ntrainer.kubeflow.org/v1alpha1"}}
    class crd_TrainJob crd
    crd_TrainingRuntime{{"TrainingRuntime\ntrainer.kubeflow.org/v1alpha1"}}
    class crd_TrainingRuntime crd
```

### CRDs

| Group | Version | Kind | Scope | Fields | Validation Rules | Source |
|-------|---------|------|-------|--------|------------------|--------|
| trainer.kubeflow.org | v1alpha1 | ClusterTrainingRuntime | Cluster | 1246 | 9 | [`manifests/base/crds/trainer.kubeflow.org_clustertrainingruntimes.yaml`](https://github.com/kubeflow/trainer/blob/5adde88079bb88d4fcb58072110bbbbd9c8692f7/manifests/base/crds/trainer.kubeflow.org_clustertrainingruntimes.yaml) |
| trainer.kubeflow.org | v1alpha1 | TrainJob | Namespaced | 562 | 5 | [`manifests/base/crds/trainer.kubeflow.org_trainjobs.yaml`](https://github.com/kubeflow/trainer/blob/5adde88079bb88d4fcb58072110bbbbd9c8692f7/manifests/base/crds/trainer.kubeflow.org_trainjobs.yaml) |
| trainer.kubeflow.org | v1alpha1 | TrainingRuntime | Namespaced | 1246 | 9 | [`manifests/base/crds/trainer.kubeflow.org_trainingruntimes.yaml`](https://github.com/kubeflow/trainer/blob/5adde88079bb88d4fcb58072110bbbbd9c8692f7/manifests/base/crds/trainer.kubeflow.org_trainingruntimes.yaml) |

## Dependencies

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.3 |
| k8s.io/api | v0.34.1 |
| k8s.io/apimachinery | v0.34.1 |
| k8s.io/client-go | v0.34.1 |
| sigs.k8s.io/controller-runtime | v0.22.3 |

