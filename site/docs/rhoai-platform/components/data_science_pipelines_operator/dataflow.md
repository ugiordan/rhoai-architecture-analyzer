# data-science-pipelines-operator: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

| Type | GVK | Source |
|------|-----|--------|
| For | api/v1/DataSciencePipelinesApplication | [`controllers/dspipeline_controller.go:796`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/df94cb0eaab69dfb8c641ee8eef47a643921109f/controllers/dspipeline_controller.go#L796) |
| Owns | /v1/ConfigMap | [`controllers/dspipeline_controller.go:799`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/df94cb0eaab69dfb8c641ee8eef47a643921109f/controllers/dspipeline_controller.go#L799) |
| Owns | /v1/PersistentVolumeClaim | [`controllers/dspipeline_controller.go:802`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/df94cb0eaab69dfb8c641ee8eef47a643921109f/controllers/dspipeline_controller.go#L802) |
| Owns | /v1/Secret | [`controllers/dspipeline_controller.go:798`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/df94cb0eaab69dfb8c641ee8eef47a643921109f/controllers/dspipeline_controller.go#L798) |
| Owns | /v1/Service | [`controllers/dspipeline_controller.go:800`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/df94cb0eaab69dfb8c641ee8eef47a643921109f/controllers/dspipeline_controller.go#L800) |
| Owns | /v1/ServiceAccount | [`controllers/dspipeline_controller.go:801`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/df94cb0eaab69dfb8c641ee8eef47a643921109f/controllers/dspipeline_controller.go#L801) |
| Owns | apps/v1/Deployment | [`controllers/dspipeline_controller.go:797`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/df94cb0eaab69dfb8c641ee8eef47a643921109f/controllers/dspipeline_controller.go#L797) |
| Owns | networking.k8s.io/v1/NetworkPolicy | [`controllers/dspipeline_controller.go:803`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/df94cb0eaab69dfb8c641ee8eef47a643921109f/controllers/dspipeline_controller.go#L803) |
| Owns | rbac.authorization.k8s.io/v1/Role | [`controllers/dspipeline_controller.go:804`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/df94cb0eaab69dfb8c641ee8eef47a643921109f/controllers/dspipeline_controller.go#L804) |
| Owns | rbac.authorization.k8s.io/v1/RoleBinding | [`controllers/dspipeline_controller.go:805`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/df94cb0eaab69dfb8c641ee8eef47a643921109f/controllers/dspipeline_controller.go#L805) |
| Owns | route/v1/Route | [`controllers/dspipeline_controller.go:806`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/df94cb0eaab69dfb8c641ee8eef47a643921109f/controllers/dspipeline_controller.go#L806) |

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for data-science-pipelines-operator

    participant KubernetesAPI as Kubernetes API
    participant data_science_pipelines_operator_controller_manager as data-science-pipelines-operator-controller-manager
    participant mariadb as mariadb
    participant minio as minio

    KubernetesAPI->>+data_science_pipelines_operator_controller_manager: Watch DataSciencePipelinesApplication (reconcile)
    data_science_pipelines_operator_controller_manager->>KubernetesAPI: Create/Update ConfigMap
    data_science_pipelines_operator_controller_manager->>KubernetesAPI: Create/Update PersistentVolumeClaim
    data_science_pipelines_operator_controller_manager->>KubernetesAPI: Create/Update Secret
    data_science_pipelines_operator_controller_manager->>KubernetesAPI: Create/Update Service
    data_science_pipelines_operator_controller_manager->>KubernetesAPI: Create/Update ServiceAccount
    data_science_pipelines_operator_controller_manager->>KubernetesAPI: Create/Update Deployment
    data_science_pipelines_operator_controller_manager->>KubernetesAPI: Create/Update NetworkPolicy
    data_science_pipelines_operator_controller_manager->>KubernetesAPI: Create/Update Role
    data_science_pipelines_operator_controller_manager->>KubernetesAPI: Create/Update RoleBinding
    data_science_pipelines_operator_controller_manager->>KubernetesAPI: Create/Update Route

    Note over data_science_pipelines_operator_controller_manager: Exposed Services
    Note right of data_science_pipelines_operator_controller_manager: data-science-pipelines-operator-service:8080/TCP [metrics]
    Note right of data_science_pipelines_operator_controller_manager: ds-pipeline-workflow-controller-metrics-template-value:9090/TCP [metrics]
    Note right of data_science_pipelines_operator_controller_manager: mariadb:3306/TCP []
    Note right of data_science_pipelines_operator_controller_manager: mariadb-template-value:3306/TCP []
    Note right of data_science_pipelines_operator_controller_manager: minio:9000/TCP [https]
    Note right of data_science_pipelines_operator_controller_manager: minio:9001/TCP [console]
    Note right of data_science_pipelines_operator_controller_manager: minio-service:9000/TCP [http]
    Note right of data_science_pipelines_operator_controller_manager: minio-template-value:9000/TCP [http]
    Note right of data_science_pipelines_operator_controller_manager: minio-template-value:80/TCP [kfp-ui-http]
    Note right of data_science_pipelines_operator_controller_manager: ml-pipeline:8443/TCP [proxy]
    Note right of data_science_pipelines_operator_controller_manager: ml-pipeline:8888/TCP [http]
    Note right of data_science_pipelines_operator_controller_manager: ml-pipeline:8887/TCP [grpc]
    Note right of data_science_pipelines_operator_controller_manager: pypi-server:8080/TCP [pypi-server]
    Note right of data_science_pipelines_operator_controller_manager: template-value:8443/TCP [proxy]
    Note right of data_science_pipelines_operator_controller_manager: template-value:8888/TCP [http]
    Note right of data_science_pipelines_operator_controller_manager: template-value:8887/TCP [grpc]
    Note right of data_science_pipelines_operator_controller_manager: template-value:8443/TCP [webhook]

    Note over KubernetesAPI: Defined CRDs
    Note right of KubernetesAPI: DataSciencePipelinesApplication (datasciencepipelinesapplications.opendatahub.io/v1)
    Note right of KubernetesAPI: ScheduledWorkflow (kubeflow.org/v1beta1)
    Note right of KubernetesAPI: Pipeline (pipelines.kubeflow.org/v2beta1)
    Note right of KubernetesAPI: PipelineVersion (pipelines.kubeflow.org/v2beta1)
```

### Webhooks

| Name | Type | Path | Failure Policy | Service | Source |
|------|------|------|----------------|---------|--------|
| pipelineversions.pipelines.kubeflow.org | mutating | /webhooks/mutate-pipelineversion | Fail | template-value/template-value | [`config/internal/webhook/mutating_webhook.yaml.tmpl`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/df94cb0eaab69dfb8c641ee8eef47a643921109f/config/internal/webhook/mutating_webhook.yaml.tmpl) |
| pipelineversions.pipelines.kubeflow.org | validating | /webhooks/validate-pipelineversion | Fail | template-value/template-value | [`config/internal/webhook/validating_webhook.yaml.tmpl`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/df94cb0eaab69dfb8c641ee8eef47a643921109f/config/internal/webhook/validating_webhook.yaml.tmpl) |

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

### ConfigMaps

| Name | Data Keys | Source |
|------|-----------|--------|
| workflow-controller-configmap |  | [`config/argo/configmap.workflow-controller-configmap.yaml`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/df94cb0eaab69dfb8c641ee8eef47a643921109f/config/argo/configmap.workflow-controller-configmap.yaml) |

