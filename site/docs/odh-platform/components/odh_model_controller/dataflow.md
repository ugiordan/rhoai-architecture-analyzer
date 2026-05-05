# odh-model-controller: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

| Type | GVK | Source |
|------|-----|--------|
| For | /v1/ConfigMap | [`internal/controller/core/configmap_controller.go:157`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/core/configmap_controller.go#L157) |
| For | /v1/Pod | [`internal/controller/core/pod_controller.go:53`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/core/pod_controller.go#L53) |
| For | /v1/Secret | [`internal/controller/core/secret_controller.go:251`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/core/secret_controller.go#L251) |
| For | apis/v1/Gateway | [`internal/controller/serving/llm/gateway_controller.go:801`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/llm/gateway_controller.go#L801) |
| For | nim/v1/Account | [`internal/controller/nim/account_controller.go:87`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/nim/account_controller.go#L87) |
| For | serving/v1alpha1/InferenceGraph | [`internal/controller/serving/inferencegraph_controller.go:74`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/inferencegraph_controller.go#L74) |
| For | serving/v1alpha1/ServingRuntime | [`internal/controller/serving/servingruntime_controller.go:418`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/servingruntime_controller.go#L418) |
| For | serving/v1alpha2/LLMInferenceService | [`internal/controller/serving/llm/llm_inferenceservice_controller.go:149`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/llm/llm_inferenceservice_controller.go#L149) |
| For | serving/v1beta1/InferenceService | [`internal/controller/serving/inferenceservice_controller.go:205`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/inferenceservice_controller.go#L205) |
| Owns | /v1/ConfigMap | [`internal/controller/nim/account_controller.go:88`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/nim/account_controller.go#L88) |
| Owns | /v1/ConfigMap | [`internal/controller/serving/inferenceservice_controller.go:211`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/inferenceservice_controller.go#L211) |
| Owns | /v1/Namespace | [`internal/controller/serving/inferenceservice_controller.go:207`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/inferenceservice_controller.go#L207) |
| Owns | /v1/Secret | [`internal/controller/nim/account_controller.go:89`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/nim/account_controller.go#L89) |
| Owns | /v1/Secret | [`internal/controller/serving/inferenceservice_controller.go:212`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/inferenceservice_controller.go#L212) |
| Owns | /v1/Service | [`internal/controller/serving/inferenceservice_controller.go:210`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/inferenceservice_controller.go#L210) |
| Owns | /v1/ServiceAccount | [`internal/controller/serving/inferenceservice_controller.go:209`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/inferenceservice_controller.go#L209) |
| Owns | api/v1/AuthPolicy | [`internal/controller/serving/llm/gateway_controller.go:807`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/llm/gateway_controller.go#L807) |
| Owns | keda/v1alpha1/TriggerAuthentication | [`internal/controller/serving/inferenceservice_controller.go:278`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/inferenceservice_controller.go#L278) |
| Owns | monitoring/v1/PodMonitor | [`internal/controller/serving/inferenceservice_controller.go:216`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/inferenceservice_controller.go#L216) |
| Owns | monitoring/v1/ServiceMonitor | [`internal/controller/serving/inferenceservice_controller.go:215`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/inferenceservice_controller.go#L215) |
| Owns | networking.k8s.io/v1/NetworkPolicy | [`internal/controller/serving/inferenceservice_controller.go:214`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/inferenceservice_controller.go#L214) |
| Owns | networking/v1alpha3/EnvoyFilter | [`internal/controller/serving/llm/gateway_controller.go:804`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/llm/gateway_controller.go#L804) |
| Owns | rbac.authorization.k8s.io/v1/ClusterRoleBinding | [`internal/controller/serving/inferenceservice_controller.go:213`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/inferenceservice_controller.go#L213) |
| Owns | rbac.authorization.k8s.io/v1/Role | [`internal/controller/serving/inferenceservice_controller.go:217`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/inferenceservice_controller.go#L217) |
| Owns | rbac.authorization.k8s.io/v1/RoleBinding | [`internal/controller/serving/servingruntime_controller.go:420`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/servingruntime_controller.go#L420) |
| Owns | rbac.authorization.k8s.io/v1/RoleBinding | [`internal/controller/serving/inferenceservice_controller.go:218`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/inferenceservice_controller.go#L218) |
| Owns | route/v1/Route | [`internal/controller/serving/inferenceservice_controller.go:208`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/inferenceservice_controller.go#L208) |
| Owns | serving/v1alpha1/ServingRuntime | [`internal/controller/serving/inferenceservice_controller.go:206`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/inferenceservice_controller.go#L206) |
| Owns | template/v1/Template | [`internal/controller/nim/account_controller.go:90`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/nim/account_controller.go#L90) |
| Watches | /v1/ConfigMap | [`internal/controller/nim/account_controller.go:104`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/nim/account_controller.go#L104) |
| Watches | /v1/Namespace | [`internal/controller/serving/servingruntime_controller.go:437`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/servingruntime_controller.go#L437) |
| Watches | /v1/Namespace | [`internal/controller/serving/llm/gateway_controller.go:841`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/llm/gateway_controller.go#L841) |
| Watches | /v1/Secret | [`internal/controller/serving/inferenceservice_controller.go:250`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/inferenceservice_controller.go#L250) |
| Watches | /v1/Secret | [`internal/controller/nim/account_controller.go:91`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/nim/account_controller.go#L91) |
| Watches | api/v1/AuthPolicy | [`internal/controller/serving/llm/llm_inferenceservice_controller.go:157`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/llm/llm_inferenceservice_controller.go#L157) |
| Watches | api/v1beta1/Authorino | [`internal/controller/serving/llm/llm_inferenceservice_controller.go:181`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/llm/llm_inferenceservice_controller.go#L181) |
| Watches | api/v1beta1/Kuadrant | [`internal/controller/serving/llm/llm_inferenceservice_controller.go:175`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/llm/llm_inferenceservice_controller.go#L175) |
| Watches | rbac.authorization.k8s.io/v1/RoleBinding | [`internal/controller/serving/servingruntime_controller.go:452`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/servingruntime_controller.go#L452) |
| Watches | serving/v1alpha1/ServingRuntime | [`internal/controller/serving/inferenceservice_controller.go:220`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/inferenceservice_controller.go#L220) |
| Watches | serving/v1alpha2/LLMInferenceService | [`internal/controller/serving/llm/gateway_controller.go:811`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/llm/gateway_controller.go#L811) |
| Watches | serving/v1alpha2/LLMInferenceServiceConfig | [`internal/controller/serving/llm/gateway_controller.go:826`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/llm/gateway_controller.go#L826) |

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for odh-model-controller

    participant KubernetesAPI as Kubernetes API
    participant odh_model_controller as odh-model-controller

    KubernetesAPI->>+odh_model_controller: Watch ConfigMap (reconcile)
    KubernetesAPI->>+odh_model_controller: Watch Pod (reconcile)
    KubernetesAPI->>+odh_model_controller: Watch Secret (reconcile)
    KubernetesAPI->>+odh_model_controller: Watch Gateway (reconcile)
    KubernetesAPI->>+odh_model_controller: Watch Account (reconcile)
    KubernetesAPI->>+odh_model_controller: Watch InferenceGraph (reconcile)
    KubernetesAPI->>+odh_model_controller: Watch ServingRuntime (reconcile)
    KubernetesAPI->>+odh_model_controller: Watch LLMInferenceService (reconcile)
    KubernetesAPI->>+odh_model_controller: Watch InferenceService (reconcile)
    odh_model_controller->>KubernetesAPI: Create/Update ConfigMap
    odh_model_controller->>KubernetesAPI: Create/Update ConfigMap
    odh_model_controller->>KubernetesAPI: Create/Update Namespace
    odh_model_controller->>KubernetesAPI: Create/Update Secret
    odh_model_controller->>KubernetesAPI: Create/Update Secret
    odh_model_controller->>KubernetesAPI: Create/Update Service
    odh_model_controller->>KubernetesAPI: Create/Update ServiceAccount
    odh_model_controller->>KubernetesAPI: Create/Update AuthPolicy
    odh_model_controller->>KubernetesAPI: Create/Update TriggerAuthentication
    odh_model_controller->>KubernetesAPI: Create/Update PodMonitor
    odh_model_controller->>KubernetesAPI: Create/Update ServiceMonitor
    odh_model_controller->>KubernetesAPI: Create/Update NetworkPolicy
    odh_model_controller->>KubernetesAPI: Create/Update EnvoyFilter
    odh_model_controller->>KubernetesAPI: Create/Update ClusterRoleBinding
    odh_model_controller->>KubernetesAPI: Create/Update Role
    odh_model_controller->>KubernetesAPI: Create/Update RoleBinding
    odh_model_controller->>KubernetesAPI: Create/Update RoleBinding
    odh_model_controller->>KubernetesAPI: Create/Update Route
    odh_model_controller->>KubernetesAPI: Create/Update ServingRuntime
    odh_model_controller->>KubernetesAPI: Create/Update Template
    KubernetesAPI-->>+odh_model_controller: Watch ConfigMap (informer)
    KubernetesAPI-->>+odh_model_controller: Watch Namespace (informer)
    KubernetesAPI-->>+odh_model_controller: Watch Namespace (informer)
    KubernetesAPI-->>+odh_model_controller: Watch Secret (informer)
    KubernetesAPI-->>+odh_model_controller: Watch Secret (informer)
    KubernetesAPI-->>+odh_model_controller: Watch AuthPolicy (informer)
    KubernetesAPI-->>+odh_model_controller: Watch Authorino (informer)
    KubernetesAPI-->>+odh_model_controller: Watch Kuadrant (informer)
    KubernetesAPI-->>+odh_model_controller: Watch RoleBinding (informer)
    KubernetesAPI-->>+odh_model_controller: Watch ServingRuntime (informer)
    KubernetesAPI-->>+odh_model_controller: Watch LLMInferenceService (informer)
    KubernetesAPI-->>+odh_model_controller: Watch LLMInferenceServiceConfig (informer)

    Note over odh_model_controller: Exposed Services
    Note right of odh_model_controller: model-serving-api:443/TCP [https]
    Note right of odh_model_controller: model-serving-api:8080/TCP [metrics]
    Note right of odh_model_controller: odh-model-controller-webhook-service:443/TCP []

    Note over KubernetesAPI: Defined CRDs
    Note right of KubernetesAPI: Account (nim.opendatahub.io/v1)
```

### Webhooks

| Name | Type | Path | Failure Policy | Service | Source |
|------|------|------|----------------|---------|--------|
| minferencegraph-v1alpha1.odh-model-controller.opendatahub.io | mutating | /mutate-serving-kserve-io-v1alpha1-inferencegraph | Fail | system/webhook-service | [`config/webhook/manifests.yaml`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/config/webhook/manifests.yaml) |
| minferenceservice-v1beta1.odh-model-controller.opendatahub.io | mutating | /mutate-serving-kserve-io-v1beta1-inferenceservice | Fail | system/webhook-service | [`config/webhook/manifests.yaml`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/config/webhook/manifests.yaml) |
| mutating.pod.odh-model-controller.opendatahub.io | mutating | /mutate--v1-pod | Fail | system/webhook-service | [`config/webhook/manifests.yaml`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/config/webhook/manifests.yaml) |
| validating.isvc.odh-model-controller.opendatahub.io | validating | /validate-serving-kserve-io-v1beta1-inferenceservice | Fail | system/webhook-service | [`config/webhook/manifests.yaml`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/config/webhook/manifests.yaml) |
| validating.nim.account.odh-model-controller.opendatahub.io | validating | /validate-nim-opendatahub-io-v1-account | Fail | system/webhook-service | [`config/webhook/manifests.yaml`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/config/webhook/manifests.yaml) |
| vinferencegraph-v1alpha1.odh-model-controller.opendatahub.io | validating | /validate-serving-kserve-io-v1alpha1-inferencegraph | Fail | system/webhook-service | [`config/webhook/manifests.yaml`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/config/webhook/manifests.yaml) |

### HTTP Endpoints

| Method | Path | Source |
|--------|------|--------|
| * | /api/v1/gateways | [`server/server.go:23`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/server/server.go#L23) |
| * | /healthz | [`server/server.go:19`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/server/server.go#L19) |
| * | /metrics | [`server/observability/observability.go:87`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/server/observability/observability.go#L87) |
| * | /readyz | [`server/server.go:20`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/server/server.go#L20) |
| * | gateway.networking.k8s.io | [`internal/controller/serving/llm/fixture/gwapi_builders.go:258`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/llm/fixture/gwapi_builders.go#L258) |
| * | gateway.networking.k8s.io | [`internal/controller/serving/llm/fixture/gwapi_builders.go:276`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/llm/fixture/gwapi_builders.go#L276) |
| * | gateway.networking.k8s.io | [`internal/controller/serving/llm/fixture/gwapi_builders.go:448`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/llm/fixture/gwapi_builders.go#L448) |
| * | inference.networking.x-k8s.io | [`internal/controller/serving/llm/fixture/gwapi_builders.go:350`](https://github.com/opendatahub-io/odh-model-controller/blob/6546a54fc9bdb8f1702596ef91ecfe8d93403e5f/internal/controller/serving/llm/fixture/gwapi_builders.go#L350) |

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

