# mlflow-operator: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

| Type | GVK | Source |
|------|-----|--------|
| For | api/v1/MLflow | [`internal/controller/mlflow_controller.go:414`](https://github.com/opendatahub-io/mlflow-operator/blob/1fed87d8872e24b4a28bcb5e2a2d3e6e3d7f57ff/internal/controller/mlflow_controller.go#L414) |
| Owns | /v1/PersistentVolumeClaim | [`internal/controller/mlflow_controller.go:421`](https://github.com/opendatahub-io/mlflow-operator/blob/1fed87d8872e24b4a28bcb5e2a2d3e6e3d7f57ff/internal/controller/mlflow_controller.go#L421) |
| Owns | /v1/Secret | [`internal/controller/mlflow_controller.go:418`](https://github.com/opendatahub-io/mlflow-operator/blob/1fed87d8872e24b4a28bcb5e2a2d3e6e3d7f57ff/internal/controller/mlflow_controller.go#L418) |
| Owns | /v1/Service | [`internal/controller/mlflow_controller.go:419`](https://github.com/opendatahub-io/mlflow-operator/blob/1fed87d8872e24b4a28bcb5e2a2d3e6e3d7f57ff/internal/controller/mlflow_controller.go#L419) |
| Owns | /v1/ServiceAccount | [`internal/controller/mlflow_controller.go:420`](https://github.com/opendatahub-io/mlflow-operator/blob/1fed87d8872e24b4a28bcb5e2a2d3e6e3d7f57ff/internal/controller/mlflow_controller.go#L420) |
| Owns | apis/v1/HTTPRoute | [`internal/controller/mlflow_controller.go:450`](https://github.com/opendatahub-io/mlflow-operator/blob/1fed87d8872e24b4a28bcb5e2a2d3e6e3d7f57ff/internal/controller/mlflow_controller.go#L450) |
| Owns | apps/v1/Deployment | [`internal/controller/mlflow_controller.go:415`](https://github.com/opendatahub-io/mlflow-operator/blob/1fed87d8872e24b4a28bcb5e2a2d3e6e3d7f57ff/internal/controller/mlflow_controller.go#L415) |
| Owns | batch/v1/CronJob | [`internal/controller/mlflow_controller.go:417`](https://github.com/opendatahub-io/mlflow-operator/blob/1fed87d8872e24b4a28bcb5e2a2d3e6e3d7f57ff/internal/controller/mlflow_controller.go#L417) |
| Owns | batch/v1/Job | [`internal/controller/mlflow_controller.go:416`](https://github.com/opendatahub-io/mlflow-operator/blob/1fed87d8872e24b4a28bcb5e2a2d3e6e3d7f57ff/internal/controller/mlflow_controller.go#L416) |
| Owns | console/v1/ConsoleLink | [`internal/controller/mlflow_controller.go:442`](https://github.com/opendatahub-io/mlflow-operator/blob/1fed87d8872e24b4a28bcb5e2a2d3e6e3d7f57ff/internal/controller/mlflow_controller.go#L442) |
| Owns | monitoring/v1/ServiceMonitor | [`internal/controller/mlflow_controller.go:458`](https://github.com/opendatahub-io/mlflow-operator/blob/1fed87d8872e24b4a28bcb5e2a2d3e6e3d7f57ff/internal/controller/mlflow_controller.go#L458) |
| Owns | rbac.authorization.k8s.io/v1/ClusterRoleBinding | [`internal/controller/mlflow_controller.go:427`](https://github.com/opendatahub-io/mlflow-operator/blob/1fed87d8872e24b4a28bcb5e2a2d3e6e3d7f57ff/internal/controller/mlflow_controller.go#L427) |
| Watches | rbac.authorization.k8s.io/v1/ClusterRole | [`internal/controller/mlflow_controller.go:426`](https://github.com/opendatahub-io/mlflow-operator/blob/1fed87d8872e24b4a28bcb5e2a2d3e6e3d7f57ff/internal/controller/mlflow_controller.go#L426) |

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for mlflow-operator

    participant KubernetesAPI as Kubernetes API
    participant mlflow_operator_controller_manager as mlflow-operator-controller-manager
    participant postgres_deployment as postgres-deployment

    KubernetesAPI->>+mlflow_operator_controller_manager: Watch MLflow (reconcile)
    mlflow_operator_controller_manager->>KubernetesAPI: Create/Update PersistentVolumeClaim
    mlflow_operator_controller_manager->>KubernetesAPI: Create/Update Secret
    mlflow_operator_controller_manager->>KubernetesAPI: Create/Update Service
    mlflow_operator_controller_manager->>KubernetesAPI: Create/Update ServiceAccount
    mlflow_operator_controller_manager->>KubernetesAPI: Create/Update HTTPRoute
    mlflow_operator_controller_manager->>KubernetesAPI: Create/Update Deployment
    mlflow_operator_controller_manager->>KubernetesAPI: Create/Update CronJob
    mlflow_operator_controller_manager->>KubernetesAPI: Create/Update Job
    mlflow_operator_controller_manager->>KubernetesAPI: Create/Update ConsoleLink
    mlflow_operator_controller_manager->>KubernetesAPI: Create/Update ServiceMonitor
    mlflow_operator_controller_manager->>KubernetesAPI: Create/Update ClusterRoleBinding
    KubernetesAPI-->>+mlflow_operator_controller_manager: Watch ClusterRole (informer)

    Note over mlflow_operator_controller_manager: Exposed Services
    Note right of mlflow_operator_controller_manager: minio-service:9000/TCP [https]
    Note right of mlflow_operator_controller_manager: mlflow-operator-controller-manager-metrics-service:8443/TCP [https]
    Note right of mlflow_operator_controller_manager: postgres-service:5432/TCP []

    Note over KubernetesAPI: Defined CRDs
    Note right of KubernetesAPI: MLflowConfig (mlflow.kubeflow.org/v1)
    Note right of KubernetesAPI: MLflow (mlflow.opendatahub.io/v1)
```

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

### Helm

**Chart:** mlflow v0.1.0

