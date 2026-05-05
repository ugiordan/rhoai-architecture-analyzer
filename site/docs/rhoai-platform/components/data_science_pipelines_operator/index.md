# data-science-pipelines-operator

> **Architecture snapshot: 2026-05-05** (2026-05-05)


**Repository:** opendatahub-io/data-science-pipelines-operator  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-05T15:10:08Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 4 |
| Deployments | 3 |
| Services | 11 |
| Secrets | 4 |
| Cluster Roles | 4 |
| Controller Watches | 11 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for data-science-pipelines-operator

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["data-science-pipelines-operator Controller"]
        dep_1["data-science-pipelines-operator-controller-manager"]
        class dep_1 controller
        dep_2["mariadb"]
        class dep_2 controller
        dep_3["minio"]
        class dep_3 controller
    end

    crd_DataSciencePipelinesApplication{{"DataSciencePipelinesApplication\ndatasciencepipelinesapplications.opendatahub.io/v1"}}
    class crd_DataSciencePipelinesApplication crd
    crd_DataSciencePipelinesApplication -->|"For (reconciles)"| controller
    crd_ScheduledWorkflow{{"ScheduledWorkflow\nkubeflow.org/v1beta1"}}
    class crd_ScheduledWorkflow crd
    crd_Pipeline{{"Pipeline\npipelines.kubeflow.org/v2beta1"}}
    class crd_Pipeline crd
    crd_PipelineVersion{{"PipelineVersion\npipelines.kubeflow.org/v2beta1"}}
    class crd_PipelineVersion crd
    controller -->|"Owns"| owned_4["ConfigMap"]
    class owned_4 owned
    controller -->|"Owns"| owned_5["Deployment"]
    class owned_5 owned
    controller -->|"Owns"| owned_6["NetworkPolicy"]
    class owned_6 owned
    controller -->|"Owns"| owned_7["PersistentVolumeClaim"]
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
```

### CRDs

| Group | Version | Kind | Scope | Fields | Validation Rules | Source |
|-------|---------|------|-------|--------|------------------|--------|
| datasciencepipelinesapplications.opendatahub.io | v1 | DataSciencePipelinesApplication | Namespaced | 205 | 2 | [`config/crd/bases/datasciencepipelinesapplications.opendatahub.io_datasciencepipelinesapplications.yaml`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/df94cb0eaab69dfb8c641ee8eef47a643921109f/config/crd/bases/datasciencepipelinesapplications.opendatahub.io_datasciencepipelinesapplications.yaml) |
| kubeflow.org | v1beta1 | ScheduledWorkflow | Namespaced | 5 | 0 | [`config/crd/bases/scheduledworkflows.yaml`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/df94cb0eaab69dfb8c641ee8eef47a643921109f/config/crd/bases/scheduledworkflows.yaml) |
| pipelines.kubeflow.org | v2beta1 | Pipeline | Namespaced | 7 | 0 | [`config/crd/bases/pipelines.kubeflow.org_pipelines.yaml`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/df94cb0eaab69dfb8c641ee8eef47a643921109f/config/crd/bases/pipelines.kubeflow.org_pipelines.yaml) |
| pipelines.kubeflow.org | v2beta1 | PipelineVersion | Namespaced | 18 | 0 | [`config/crd/bases/pipelines.kubeflow.org_pipelineversions.yaml`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/df94cb0eaab69dfb8c641ee8eef47a643921109f/config/crd/bases/pipelines.kubeflow.org_pipelineversions.yaml) |

## Dependencies

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.3 |
| github.com/prometheus/client_golang | v1.23.2 |
| k8s.io/api | v0.35.3 |
| k8s.io/apimachinery | v0.35.3 |
| k8s.io/client-go | v0.35.3 |
| sigs.k8s.io/controller-runtime | v0.23.3 |

