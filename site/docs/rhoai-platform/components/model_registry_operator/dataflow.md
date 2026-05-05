# model-registry-operator: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

| Type | GVK | Source |
|------|-----|--------|
| For | api/v1beta1/ModelRegistry | [`internal/controller/modelregistry_controller.go:258`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/internal/controller/modelregistry_controller.go#L258) |
| Owns | /v1/Service | [`internal/controller/modelregistry_controller.go:259`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/internal/controller/modelregistry_controller.go#L259) |
| Owns | /v1/ServiceAccount | [`internal/controller/modelregistry_controller.go:260`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/internal/controller/modelregistry_controller.go#L260) |
| Owns | apps/v1/Deployment | [`internal/controller/modelregistry_controller.go:261`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/internal/controller/modelregistry_controller.go#L261) |
| Owns | networking.k8s.io/v1/NetworkPolicy | [`internal/controller/modelregistry_controller.go:263`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/internal/controller/modelregistry_controller.go#L263) |
| Owns | rbac.authorization.k8s.io/v1/Role | [`internal/controller/modelregistry_controller.go:262`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/internal/controller/modelregistry_controller.go#L262) |
| Owns | rbac.authorization.k8s.io/v1/RoleBinding | [`internal/controller/modelregistry_controller.go:265`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/internal/controller/modelregistry_controller.go#L265) |
| Watches | /v1/ConfigMap | [`internal/controller/modelcatalog_controller.go:1335`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/internal/controller/modelcatalog_controller.go#L1335) |
| Watches | /v1/PersistentVolumeClaim | [`internal/controller/modelcatalog_controller.go:1339`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/internal/controller/modelcatalog_controller.go#L1339) |
| Watches | /v1/Secret | [`internal/controller/modelcatalog_controller.go:1336`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/internal/controller/modelcatalog_controller.go#L1336) |
| Watches | /v1/Service | [`internal/controller/modelcatalog_controller.go:1338`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/internal/controller/modelcatalog_controller.go#L1338) |
| Watches | /v1/ServiceAccount | [`internal/controller/modelcatalog_controller.go:1337`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/internal/controller/modelcatalog_controller.go#L1337) |
| Watches | apps/v1/Deployment | [`internal/controller/modelcatalog_controller.go:1334`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/internal/controller/modelcatalog_controller.go#L1334) |
| Watches | networking.k8s.io/v1/NetworkPolicy | [`internal/controller/modelcatalog_controller.go:1343`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/internal/controller/modelcatalog_controller.go#L1343) |
| Watches | rbac.authorization.k8s.io/v1/ClusterRoleBinding | [`internal/controller/modelcatalog_controller.go:1340`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/internal/controller/modelcatalog_controller.go#L1340) |
| Watches | rbac.authorization.k8s.io/v1/Role | [`internal/controller/modelcatalog_controller.go:1341`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/internal/controller/modelcatalog_controller.go#L1341) |
| Watches | rbac.authorization.k8s.io/v1/RoleBinding | [`internal/controller/modelcatalog_controller.go:1342`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/internal/controller/modelcatalog_controller.go#L1342) |

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for model-registry-operator

    participant KubernetesAPI as Kubernetes API
    participant model_registry_operator_controller_manager as model-registry-operator-controller-manager

    KubernetesAPI->>+model_registry_operator_controller_manager: Watch ModelRegistry (reconcile)
    model_registry_operator_controller_manager->>KubernetesAPI: Create/Update Service
    model_registry_operator_controller_manager->>KubernetesAPI: Create/Update ServiceAccount
    model_registry_operator_controller_manager->>KubernetesAPI: Create/Update Deployment
    model_registry_operator_controller_manager->>KubernetesAPI: Create/Update NetworkPolicy
    model_registry_operator_controller_manager->>KubernetesAPI: Create/Update Role
    model_registry_operator_controller_manager->>KubernetesAPI: Create/Update RoleBinding
    KubernetesAPI-->>+model_registry_operator_controller_manager: Watch ConfigMap (informer)
    KubernetesAPI-->>+model_registry_operator_controller_manager: Watch PersistentVolumeClaim (informer)
    KubernetesAPI-->>+model_registry_operator_controller_manager: Watch Secret (informer)
    KubernetesAPI-->>+model_registry_operator_controller_manager: Watch Service (informer)
    KubernetesAPI-->>+model_registry_operator_controller_manager: Watch ServiceAccount (informer)
    KubernetesAPI-->>+model_registry_operator_controller_manager: Watch Deployment (informer)
    KubernetesAPI-->>+model_registry_operator_controller_manager: Watch NetworkPolicy (informer)
    KubernetesAPI-->>+model_registry_operator_controller_manager: Watch ClusterRoleBinding (informer)
    KubernetesAPI-->>+model_registry_operator_controller_manager: Watch Role (informer)
    KubernetesAPI-->>+model_registry_operator_controller_manager: Watch RoleBinding (informer)

    Note over model_registry_operator_controller_manager: Exposed Services
    Note right of model_registry_operator_controller_manager: model-registry-operator-controller-manager-metrics-service:8443/TCP [https]
    Note right of model_registry_operator_controller_manager: model-registry-operator-webhook-service:443/TCP []
    Note right of model_registry_operator_controller_manager: template-value:0/TCP [https-api]
    Note right of model_registry_operator_controller_manager: template-value:0/TCP [http-api]
    Note right of model_registry_operator_controller_manager: template-value:0/TCP [https-api]
    Note right of model_registry_operator_controller_manager: template-value:0/TCP [http-api]
    Note right of model_registry_operator_controller_manager: template-value-postgres:5432/TCP [postgresql]
    Note right of model_registry_operator_controller_manager: template-value-postgres:5432/TCP []

    Note over KubernetesAPI: Defined CRDs
    Note right of KubernetesAPI: ModelRegistry (modelregistry.opendatahub.io/v1beta1)
```

### Webhooks

| Name | Type | Path | Failure Policy | Service | Source |
|------|------|------|----------------|---------|--------|
| mmodelregistry.opendatahub.io | mutating | /mutate-modelregistry-opendatahub-io-modelregistry | Fail | opendatahub/model-registry-operator-webhook-service | [`kustomize:config/overlays/odh (model-registry-operator-mutating-webhook-configuration)`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/kustomize:config/overlays/odh (model-registry-operator-mutating-webhook-configuration)) |
| vmodelregistry.opendatahub.io | validating | /validate-modelregistry-opendatahub-io-modelregistry | Fail | opendatahub/model-registry-operator-webhook-service | [`kustomize:config/overlays/odh (model-registry-operator-validating-webhook-configuration)`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/kustomize:config/overlays/odh (model-registry-operator-validating-webhook-configuration)) |

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

