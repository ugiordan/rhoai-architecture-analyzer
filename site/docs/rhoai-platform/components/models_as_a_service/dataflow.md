# models-as-a-service: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

| Type | GVK | Source |
|------|-----|--------|
| For | apps/v1/Deployment | [`maas-controller/pkg/controller/maas/self_deployment_controller.go:222`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-controller/pkg/controller/maas/self_deployment_controller.go#L222) |
| For | maas/v1alpha1/ExternalModel | [`maas-controller/pkg/reconciler/externalmodel/reconciler.go:299`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-controller/pkg/reconciler/externalmodel/reconciler.go#L299) |
| For | maas/v1alpha1/MaaSAuthPolicy | [`maas-controller/pkg/controller/maas/maasauthpolicy_controller.go:1188`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-controller/pkg/controller/maas/maasauthpolicy_controller.go#L1188) |
| For | maas/v1alpha1/MaaSModelRef | [`maas-controller/pkg/controller/maas/maasmodelref_controller.go:337`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-controller/pkg/controller/maas/maasmodelref_controller.go#L337) |
| For | maas/v1alpha1/MaaSSubscription | [`maas-controller/pkg/controller/maas/maassubscription_controller.go:983`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-controller/pkg/controller/maas/maassubscription_controller.go#L983) |
| For | maas/v1alpha1/Tenant | [`maas-controller/pkg/controller/maas/tenant_controller.go:189`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-controller/pkg/controller/maas/tenant_controller.go#L189) |
| Watches | apis/v1/HTTPRoute | [`maas-controller/pkg/controller/maas/maasauthpolicy_controller.go:1194`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-controller/pkg/controller/maas/maasauthpolicy_controller.go#L1194) |
| Watches | apis/v1/HTTPRoute | [`maas-controller/pkg/controller/maas/maasmodelref_controller.go:343`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-controller/pkg/controller/maas/maasmodelref_controller.go#L343) |
| Watches | apis/v1/HTTPRoute | [`maas-controller/pkg/controller/maas/maassubscription_controller.go:996`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-controller/pkg/controller/maas/maassubscription_controller.go#L996) |
| Watches | maas/v1alpha1/MaaSModelRef | [`maas-controller/pkg/controller/maas/maasauthpolicy_controller.go:1198`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-controller/pkg/controller/maas/maasauthpolicy_controller.go#L1198) |
| Watches | maas/v1alpha1/MaaSModelRef | [`maas-controller/pkg/controller/maas/maassubscription_controller.go:1000`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-controller/pkg/controller/maas/maassubscription_controller.go#L1000) |
| Watches | serving/v1alpha1/LLMInferenceService | [`maas-controller/pkg/controller/maas/maasmodelref_controller.go:348`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-controller/pkg/controller/maas/maasmodelref_controller.go#L348) |

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for models-as-a-service

    participant KubernetesAPI as Kubernetes API
    participant maas_api as maas-api
    participant maas_controller as maas-controller
    participant payload_processing as payload-processing

    KubernetesAPI->>+maas_api: Watch Deployment (reconcile)
    KubernetesAPI->>+maas_api: Watch ExternalModel (reconcile)
    KubernetesAPI->>+maas_api: Watch MaaSAuthPolicy (reconcile)
    KubernetesAPI->>+maas_api: Watch MaaSModelRef (reconcile)
    KubernetesAPI->>+maas_api: Watch MaaSSubscription (reconcile)
    KubernetesAPI->>+maas_api: Watch Tenant (reconcile)
    KubernetesAPI-->>+maas_api: Watch HTTPRoute (informer)
    KubernetesAPI-->>+maas_api: Watch HTTPRoute (informer)
    KubernetesAPI-->>+maas_api: Watch HTTPRoute (informer)
    KubernetesAPI-->>+maas_api: Watch MaaSModelRef (informer)
    KubernetesAPI-->>+maas_api: Watch MaaSModelRef (informer)
    KubernetesAPI-->>+maas_api: Watch LLMInferenceService (informer)

    Note over maas_api: Exposed Services
    Note right of maas_api: maas-api:8080/TCP [http]
    Note right of maas_api: maas-api:9090/TCP [metrics]
    Note right of maas_api: maas-api:0/TCP []
    Note right of maas_api: maas-api:8443/TCP [https]
    Note right of maas_api: payload-processing:9004/TCP []
```

### HTTP Endpoints

| Method | Path | Source |
|--------|------|--------|
| OPTIONS | /*path | [`maas-api/cmd/main.go:112`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-api/cmd/main.go#L112) |
| DELETE | /:id | [`maas-api/cmd/main.go:213`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-api/cmd/main.go#L213) |
| GET | /:id | [`maas-api/cmd/main.go:212`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-api/cmd/main.go#L212) |
| * | /api-keys | [`maas-api/cmd/main.go:208`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-api/cmd/main.go#L208) |
| POST | /api-keys/cleanup | [`maas-api/cmd/main.go:218`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-api/cmd/main.go#L218) |
| POST | /api-keys/validate | [`maas-api/cmd/main.go:217`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-api/cmd/main.go#L217) |
| POST | /bulk-revoke | [`maas-api/cmd/main.go:211`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-api/cmd/main.go#L211) |
| GET | /health | [`maas-api/cmd/main.go:179`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-api/cmd/main.go#L179) |
| * | /internal/v1 | [`maas-api/cmd/main.go:216`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-api/cmd/main.go#L216) |
| * | /metrics | [`maas-api/internal/metrics/server.go:19`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-api/internal/metrics/server.go#L19) |
| GET | /model/:model-id/subscriptions | [`maas-api/cmd/main.go:205`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-api/cmd/main.go#L205) |
| GET | /models | [`maas-api/cmd/main.go:201`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-api/cmd/main.go#L201) |
| POST | /search | [`maas-api/cmd/main.go:210`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-api/cmd/main.go#L210) |
| GET | /subscriptions | [`maas-api/cmd/main.go:204`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-api/cmd/main.go#L204) |
| POST | /subscriptions/select | [`maas-api/cmd/main.go:219`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-api/cmd/main.go#L219) |
| * | /v1 | [`maas-api/cmd/main.go:185`](https://github.com/red-hat-data-services/models-as-a-service/blob/72c9ad90167f991c099e22d84293e6253495d234/maas-api/cmd/main.go#L185) |

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

