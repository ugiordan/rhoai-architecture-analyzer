# data-science-pipelines-operator

> **Architecture snapshot: 2026-05-15** (2026-05-15)


**Repository:** opendatahub-io/data-science-pipelines-operator  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-15T11:38:34Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 4 |
| Deployments | 5 |
| Services | 13 |
| Secrets | 4 |
| Cluster Roles | 4 |
| Controller Watches | 12 |

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
        dep_4["the-deployment"]
        class dep_4 controller
        dep_5["the-deployment"]
        class dep_5 controller
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
    controller -->|"Owns"| owned_6["ConfigMap"]
    class owned_6 owned
    controller -->|"Owns"| owned_7["Deployment"]
    class owned_7 owned
    controller -->|"Owns"| owned_8["NetworkPolicy"]
    class owned_8 owned
    controller -->|"Owns"| owned_9["PersistentVolumeClaim"]
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
```

### CRDs

| Group | Version | Kind | Scope | Fields | Validation Rules | Discovery | Source |
|-------|---------|------|-------|--------|------------------|-----------|--------|
| datasciencepipelinesapplications.opendatahub.io | v1 | DataSciencePipelinesApplication | Namespaced | 205 | 2 | YAML | [`config/crd/bases/datasciencepipelinesapplications.opendatahub.io_datasciencepipelinesapplications.yaml`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/1e6ce36e03d4d7e8ca3ab52e0026842e035c1ad2/config/crd/bases/datasciencepipelinesapplications.opendatahub.io_datasciencepipelinesapplications.yaml) |
| kubeflow.org | v1beta1 | ScheduledWorkflow | Namespaced | 5 | 0 | YAML | [`config/crd/bases/scheduledworkflows.yaml`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/1e6ce36e03d4d7e8ca3ab52e0026842e035c1ad2/config/crd/bases/scheduledworkflows.yaml) |
| pipelines.kubeflow.org | v2beta1 | Pipeline | Namespaced | 7 | 0 | YAML | [`config/crd/bases/pipelines.kubeflow.org_pipelines.yaml`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/1e6ce36e03d4d7e8ca3ab52e0026842e035c1ad2/config/crd/bases/pipelines.kubeflow.org_pipelines.yaml) |
| pipelines.kubeflow.org | v2beta1 | PipelineVersion | Namespaced | 18 | 0 | YAML | [`config/crd/bases/pipelines.kubeflow.org_pipelineversions.yaml`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/1e6ce36e03d4d7e8ca3ab52e0026842e035c1ad2/config/crd/bases/pipelines.kubeflow.org_pipelineversions.yaml) |

## Dependencies

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.1 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.2.2 |
| github.com/go-logr/logr | v1.3.0 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.3.0 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.1 |
| github.com/go-logr/logr | v1.2.2 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/common | v0.66.1 |
| github.com/prometheus/common | v0.66.1 |
| github.com/prometheus/procfs | v0.16.1 |
| github.com/prometheus/procfs | v0.16.1 |
| google.golang.org/grpc | v1.72.2 |
| google.golang.org/grpc | v1.72.2 |
| k8s.io/api | v0.21.3 |
| k8s.io/api | v0.22.5 |
| k8s.io/api | v0.35.3 |
| k8s.io/api | v0.35.0 |
| k8s.io/api | v0.35.0 |
| k8s.io/api | v0.35.3 |
| k8s.io/api | v0.35.1 |
| k8s.io/api | v0.35.3 |
| k8s.io/api | v0.35.1 |
| k8s.io/api | v0.35.3 |
| k8s.io/api | v0.21.3 |
| k8s.io/api | v0.22.5 |
| k8s.io/api | v0.35.3 |
| k8s.io/api | v0.35.3 |
| k8s.io/api | v0.35.3 |
| k8s.io/apiextensions-apiserver | v0.35.0 |
| k8s.io/apiextensions-apiserver | v0.35.0 |
| k8s.io/apimachinery | v0.35.1 |
| k8s.io/apimachinery | v0.35.0 |
| k8s.io/apimachinery | v0.35.3 |
| k8s.io/apimachinery | v0.22.5 |
| k8s.io/apimachinery | v0.35.3 |
| k8s.io/apimachinery | v0.35.1 |
| k8s.io/apimachinery | v0.22.5 |
| k8s.io/apimachinery | v0.35.3 |
| k8s.io/apimachinery | v0.35.3 |
| k8s.io/apimachinery | v0.35.3 |
| k8s.io/apimachinery | v0.35.3 |
| k8s.io/apimachinery | v0.19.7 |
| k8s.io/apimachinery | v0.21.3 |
| k8s.io/apimachinery | v0.35.3 |
| k8s.io/apimachinery | v0.35.3 |
| k8s.io/apimachinery | v0.35.3 |
| k8s.io/apimachinery | v0.19.7 |
| k8s.io/apimachinery | v0.21.3 |
| k8s.io/apimachinery | v0.35.0 |
| k8s.io/apiserver | v0.35.0 |
| k8s.io/apiserver | v0.35.3 |
| k8s.io/apiserver | v0.35.0 |
| k8s.io/apiserver | v0.35.3 |
| k8s.io/client-go | v0.22.5 |
| k8s.io/client-go | v0.35.0 |
| k8s.io/client-go | v0.21.3 |
| k8s.io/client-go | v0.35.3 |
| k8s.io/client-go | v0.35.3 |
| k8s.io/client-go | v0.21.3 |
| k8s.io/client-go | v0.35.0 |
| k8s.io/client-go | v0.35.3 |
| k8s.io/client-go | v0.35.3 |
| k8s.io/client-go | v0.35.3 |
| k8s.io/client-go | v0.22.5 |
| sigs.k8s.io/controller-runtime | v0.23.3 |
| sigs.k8s.io/controller-runtime | v0.7.2 |
| sigs.k8s.io/controller-runtime | v0.7.2 |

