# mlflow-operator: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

| Type | GVK | Source |
|------|-----|--------|
| For | api/v1/MLflow | [`internal/controller/mlflow_controller.go:408`](https://github.com/opendatahub-io/mlflow-operator/blob/f753d470caec527a7f134dec1863ddfa8fd975e5/internal/controller/mlflow_controller.go#L408) |
| Owns | /v1/PersistentVolumeClaim | [`internal/controller/mlflow_controller.go:414`](https://github.com/opendatahub-io/mlflow-operator/blob/f753d470caec527a7f134dec1863ddfa8fd975e5/internal/controller/mlflow_controller.go#L414) |
| Owns | /v1/Secret | [`internal/controller/mlflow_controller.go:411`](https://github.com/opendatahub-io/mlflow-operator/blob/f753d470caec527a7f134dec1863ddfa8fd975e5/internal/controller/mlflow_controller.go#L411) |
| Owns | /v1/Service | [`internal/controller/mlflow_controller.go:412`](https://github.com/opendatahub-io/mlflow-operator/blob/f753d470caec527a7f134dec1863ddfa8fd975e5/internal/controller/mlflow_controller.go#L412) |
| Owns | /v1/ServiceAccount | [`internal/controller/mlflow_controller.go:413`](https://github.com/opendatahub-io/mlflow-operator/blob/f753d470caec527a7f134dec1863ddfa8fd975e5/internal/controller/mlflow_controller.go#L413) |
| Owns | apis/v1/HTTPRoute | [`internal/controller/mlflow_controller.go:443`](https://github.com/opendatahub-io/mlflow-operator/blob/f753d470caec527a7f134dec1863ddfa8fd975e5/internal/controller/mlflow_controller.go#L443) |
| Owns | apps/v1/Deployment | [`internal/controller/mlflow_controller.go:409`](https://github.com/opendatahub-io/mlflow-operator/blob/f753d470caec527a7f134dec1863ddfa8fd975e5/internal/controller/mlflow_controller.go#L409) |
| Owns | batch/v1/CronJob | [`internal/controller/mlflow_controller.go:410`](https://github.com/opendatahub-io/mlflow-operator/blob/f753d470caec527a7f134dec1863ddfa8fd975e5/internal/controller/mlflow_controller.go#L410) |
| Owns | console/v1/ConsoleLink | [`internal/controller/mlflow_controller.go:435`](https://github.com/opendatahub-io/mlflow-operator/blob/f753d470caec527a7f134dec1863ddfa8fd975e5/internal/controller/mlflow_controller.go#L435) |
| Owns | monitoring/v1/ServiceMonitor | [`internal/controller/mlflow_controller.go:451`](https://github.com/opendatahub-io/mlflow-operator/blob/f753d470caec527a7f134dec1863ddfa8fd975e5/internal/controller/mlflow_controller.go#L451) |
| Owns | rbac.authorization.k8s.io/v1/ClusterRoleBinding | [`internal/controller/mlflow_controller.go:420`](https://github.com/opendatahub-io/mlflow-operator/blob/f753d470caec527a7f134dec1863ddfa8fd975e5/internal/controller/mlflow_controller.go#L420) |
| Watches | rbac.authorization.k8s.io/v1/ClusterRole | [`internal/controller/mlflow_controller.go:419`](https://github.com/opendatahub-io/mlflow-operator/blob/f753d470caec527a7f134dec1863ddfa8fd975e5/internal/controller/mlflow_controller.go#L419) |

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for mlflow-operator

    participant KubernetesAPI as Kubernetes API
    participant controller_manager as controller-manager
    participant mlflow_operator_controller_manager as mlflow-operator-controller-manager
    participant postgres_deployment as postgres-deployment

    KubernetesAPI->>+controller_manager: Watch MLflow (reconcile)
    controller_manager->>KubernetesAPI: Create/Update PersistentVolumeClaim
    controller_manager->>KubernetesAPI: Create/Update Secret
    controller_manager->>KubernetesAPI: Create/Update Service
    controller_manager->>KubernetesAPI: Create/Update ServiceAccount
    controller_manager->>KubernetesAPI: Create/Update HTTPRoute
    controller_manager->>KubernetesAPI: Create/Update Deployment
    controller_manager->>KubernetesAPI: Create/Update CronJob
    controller_manager->>KubernetesAPI: Create/Update ConsoleLink
    controller_manager->>KubernetesAPI: Create/Update ServiceMonitor
    controller_manager->>KubernetesAPI: Create/Update ClusterRoleBinding
    KubernetesAPI-->>+controller_manager: Watch ClusterRole (informer)

    Note over controller_manager: Exposed Services
    Note right of controller_manager: minio-service:9000/TCP [https]
    Note right of controller_manager: postgres-service:5432/TCP []

    Note over KubernetesAPI: Defined CRDs
    Note right of KubernetesAPI: MLflowConfig (mlflow.kubeflow.org/v1)
    Note right of KubernetesAPI: MLflow (mlflow.opendatahub.io/v1)
```

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

### Helm

**Chart:** mlflow v0.1.0

