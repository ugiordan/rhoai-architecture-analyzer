# kubeflow: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

| Type | GVK | Source |
|------|-----|--------|
| For | api/v1/Notebook | [`components/odh-notebook-controller/controllers/notebook_controller.go:737`](https://github.com/red-hat-data-services/kubeflow/blob/2afcff72f7d4be00ddb2cae2f27a36cebc6d1213/components/odh-notebook-controller/controllers/notebook_controller.go#L737) |
| For | api/v1beta1/Notebook | [`components/notebook-controller/controllers/culling_controller.go:580`](https://github.com/red-hat-data-services/kubeflow/blob/2afcff72f7d4be00ddb2cae2f27a36cebc6d1213/components/notebook-controller/controllers/culling_controller.go#L580) |
| For | api/v1beta1/Notebook | [`components/notebook-controller/controllers/notebook_controller.go:801`](https://github.com/red-hat-data-services/kubeflow/blob/2afcff72f7d4be00ddb2cae2f27a36cebc6d1213/components/notebook-controller/controllers/notebook_controller.go#L801) |
| Owns | /v1/ConfigMap | [`components/odh-notebook-controller/controllers/notebook_controller.go:741`](https://github.com/red-hat-data-services/kubeflow/blob/2afcff72f7d4be00ddb2cae2f27a36cebc6d1213/components/odh-notebook-controller/controllers/notebook_controller.go#L741) |
| Owns | /v1/Secret | [`components/odh-notebook-controller/controllers/notebook_controller.go:740`](https://github.com/red-hat-data-services/kubeflow/blob/2afcff72f7d4be00ddb2cae2f27a36cebc6d1213/components/odh-notebook-controller/controllers/notebook_controller.go#L740) |
| Owns | /v1/Service | [`components/odh-notebook-controller/controllers/notebook_controller.go:739`](https://github.com/red-hat-data-services/kubeflow/blob/2afcff72f7d4be00ddb2cae2f27a36cebc6d1213/components/odh-notebook-controller/controllers/notebook_controller.go#L739) |
| Owns | /v1/Service | [`components/notebook-controller/controllers/notebook_controller.go:803`](https://github.com/red-hat-data-services/kubeflow/blob/2afcff72f7d4be00ddb2cae2f27a36cebc6d1213/components/notebook-controller/controllers/notebook_controller.go#L803) |
| Owns | /v1/ServiceAccount | [`components/odh-notebook-controller/controllers/notebook_controller.go:738`](https://github.com/red-hat-data-services/kubeflow/blob/2afcff72f7d4be00ddb2cae2f27a36cebc6d1213/components/odh-notebook-controller/controllers/notebook_controller.go#L738) |
| Owns | apps/v1/StatefulSet | [`components/notebook-controller/controllers/notebook_controller.go:802`](https://github.com/red-hat-data-services/kubeflow/blob/2afcff72f7d4be00ddb2cae2f27a36cebc6d1213/components/notebook-controller/controllers/notebook_controller.go#L802) |
| Owns | networking.k8s.io/v1/NetworkPolicy | [`components/odh-notebook-controller/controllers/notebook_controller.go:742`](https://github.com/red-hat-data-services/kubeflow/blob/2afcff72f7d4be00ddb2cae2f27a36cebc6d1213/components/odh-notebook-controller/controllers/notebook_controller.go#L742) |
| Owns | rbac.authorization.k8s.io/v1/RoleBinding | [`components/odh-notebook-controller/controllers/notebook_controller.go:743`](https://github.com/red-hat-data-services/kubeflow/blob/2afcff72f7d4be00ddb2cae2f27a36cebc6d1213/components/odh-notebook-controller/controllers/notebook_controller.go#L743) |
| Watches | /v1/ConfigMap | [`components/odh-notebook-controller/controllers/notebook_controller.go:818`](https://github.com/red-hat-data-services/kubeflow/blob/2afcff72f7d4be00ddb2cae2f27a36cebc6d1213/components/odh-notebook-controller/controllers/notebook_controller.go#L818) |
| Watches | apis/v1/HTTPRoute | [`components/odh-notebook-controller/controllers/notebook_controller.go:747`](https://github.com/red-hat-data-services/kubeflow/blob/2afcff72f7d4be00ddb2cae2f27a36cebc6d1213/components/odh-notebook-controller/controllers/notebook_controller.go#L747) |
| Watches | apis/v1beta1/ReferenceGrant | [`components/odh-notebook-controller/controllers/notebook_controller.go:777`](https://github.com/red-hat-data-services/kubeflow/blob/2afcff72f7d4be00ddb2cae2f27a36cebc6d1213/components/odh-notebook-controller/controllers/notebook_controller.go#L777) |

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for kubeflow

    participant KubernetesAPI as Kubernetes API
    participant controller_manager as controller-manager
    participant deployment as deployment
    participant manager as manager

    KubernetesAPI->>+controller_manager: Watch Notebook (reconcile)
    KubernetesAPI->>+controller_manager: Watch Notebook (reconcile)
    KubernetesAPI->>+controller_manager: Watch Notebook (reconcile)
    controller_manager->>KubernetesAPI: Create/Update ConfigMap
    controller_manager->>KubernetesAPI: Create/Update Secret
    controller_manager->>KubernetesAPI: Create/Update Service
    controller_manager->>KubernetesAPI: Create/Update Service
    controller_manager->>KubernetesAPI: Create/Update ServiceAccount
    controller_manager->>KubernetesAPI: Create/Update StatefulSet
    controller_manager->>KubernetesAPI: Create/Update NetworkPolicy
    controller_manager->>KubernetesAPI: Create/Update RoleBinding
    KubernetesAPI-->>+controller_manager: Watch ConfigMap (informer)
    KubernetesAPI-->>+controller_manager: Watch HTTPRoute (informer)
    KubernetesAPI-->>+controller_manager: Watch ReferenceGrant (informer)

    Note over controller_manager: Exposed Services
    Note right of controller_manager: service:443/TCP []
    Note right of controller_manager: service:8080/TCP [metrics]
    Note right of controller_manager: webhook-service:443/TCP [webhook]
```

### Webhooks

| Name | Type | Path | Failure Policy | Service | Source |
|------|------|------|----------------|---------|--------|
| notebooks-validation.opendatahub.io | validating | /validate-notebook-v1 | fail |  | [`components/odh-notebook-controller/controllers/notebook_validating_webhook.go`](https://github.com/red-hat-data-services/kubeflow/blob/2afcff72f7d4be00ddb2cae2f27a36cebc6d1213/components/odh-notebook-controller/controllers/notebook_validating_webhook.go) |
| notebooks.opendatahub.io | mutating | /mutate-notebook-v1 | fail |  | [`components/odh-notebook-controller/controllers/notebook_mutating_webhook.go`](https://github.com/red-hat-data-services/kubeflow/blob/2afcff72f7d4be00ddb2cae2f27a36cebc6d1213/components/odh-notebook-controller/controllers/notebook_mutating_webhook.go) |

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

