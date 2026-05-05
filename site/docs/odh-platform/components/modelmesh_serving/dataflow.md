# modelmesh-serving: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

| Type | GVK | Source |
|------|-----|--------|
| For | /v1/Namespace | [`controllers/service_controller.go:476`](https://github.com/kserve/modelmesh-serving/blob/1fcf541d867ceb459fbc76aa1e2bef102c4816db/controllers/service_controller.go#L476) |
| For | apps/v1/Deployment | [`controllers/service_controller.go:449`](https://github.com/kserve/modelmesh-serving/blob/1fcf541d867ceb459fbc76aa1e2bef102c4816db/controllers/service_controller.go#L449) |
| For | serving/v1alpha1/Predictor | [`controllers/predictor_controller.go:594`](https://github.com/kserve/modelmesh-serving/blob/1fcf541d867ceb459fbc76aa1e2bef102c4816db/controllers/predictor_controller.go#L594) |
| For | serving/v1alpha1/ServingRuntime | [`controllers/servingruntime_controller.go:607`](https://github.com/kserve/modelmesh-serving/blob/1fcf541d867ceb459fbc76aa1e2bef102c4816db/controllers/servingruntime_controller.go#L607) |
| Owns | /v1/Service | [`controllers/service_controller.go:433`](https://github.com/kserve/modelmesh-serving/blob/1fcf541d867ceb459fbc76aa1e2bef102c4816db/controllers/service_controller.go#L433) |
| Owns | apps/v1/Deployment | [`controllers/servingruntime_controller.go:608`](https://github.com/kserve/modelmesh-serving/blob/1fcf541d867ceb459fbc76aa1e2bef102c4816db/controllers/servingruntime_controller.go#L608) |
| Watches | /v1/ConfigMap | [`controllers/servingruntime_controller.go:610`](https://github.com/kserve/modelmesh-serving/blob/1fcf541d867ceb459fbc76aa1e2bef102c4816db/controllers/servingruntime_controller.go#L610) |
| Watches | /v1/ConfigMap | [`controllers/service_controller.go:454`](https://github.com/kserve/modelmesh-serving/blob/1fcf541d867ceb459fbc76aa1e2bef102c4816db/controllers/service_controller.go#L454) |
| Watches | /v1/ConfigMap | [`controllers/service_controller.go:477`](https://github.com/kserve/modelmesh-serving/blob/1fcf541d867ceb459fbc76aa1e2bef102c4816db/controllers/service_controller.go#L477) |
| Watches | /v1/Namespace | [`controllers/servingruntime_controller.go:626`](https://github.com/kserve/modelmesh-serving/blob/1fcf541d867ceb459fbc76aa1e2bef102c4816db/controllers/servingruntime_controller.go#L626) |
| Watches | /v1/Secret | [`controllers/servingruntime_controller.go:650`](https://github.com/kserve/modelmesh-serving/blob/1fcf541d867ceb459fbc76aa1e2bef102c4816db/controllers/servingruntime_controller.go#L650) |
| Watches | monitoring/v1/ServiceMonitor | [`controllers/service_controller.go:499`](https://github.com/kserve/modelmesh-serving/blob/1fcf541d867ceb459fbc76aa1e2bef102c4816db/controllers/service_controller.go#L499) |
| Watches | monitoring/v1/ServiceMonitor | [`controllers/service_controller.go:465`](https://github.com/kserve/modelmesh-serving/blob/1fcf541d867ceb459fbc76aa1e2bef102c4816db/controllers/service_controller.go#L465) |
| Watches | serving/v1alpha1/ClusterServingRuntime | [`controllers/servingruntime_controller.go:643`](https://github.com/kserve/modelmesh-serving/blob/1fcf541d867ceb459fbc76aa1e2bef102c4816db/controllers/servingruntime_controller.go#L643) |
| Watches | serving/v1alpha1/Predictor | [`controllers/servingruntime_controller.go:619`](https://github.com/kserve/modelmesh-serving/blob/1fcf541d867ceb459fbc76aa1e2bef102c4816db/controllers/servingruntime_controller.go#L619) |
| Watches | serving/v1beta1/InferenceService | [`controllers/servingruntime_controller.go:633`](https://github.com/kserve/modelmesh-serving/blob/1fcf541d867ceb459fbc76aa1e2bef102c4816db/controllers/servingruntime_controller.go#L633) |
| Watches | serving/v1beta1/InferenceService | [`controllers/predictor_controller.go:602`](https://github.com/kserve/modelmesh-serving/blob/1fcf541d867ceb459fbc76aa1e2bef102c4816db/controllers/predictor_controller.go#L602) |

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for modelmesh-serving

    participant KubernetesAPI as Kubernetes API
    participant controller_manager as controller-manager
    participant etcd as etcd
    participant modelmesh_controller as modelmesh-controller

    KubernetesAPI->>+controller_manager: Watch Namespace (reconcile)
    KubernetesAPI->>+controller_manager: Watch Deployment (reconcile)
    KubernetesAPI->>+controller_manager: Watch Predictor (reconcile)
    KubernetesAPI->>+controller_manager: Watch ServingRuntime (reconcile)
    controller_manager->>KubernetesAPI: Create/Update Service
    controller_manager->>KubernetesAPI: Create/Update Deployment
    KubernetesAPI-->>+controller_manager: Watch ConfigMap (informer)
    KubernetesAPI-->>+controller_manager: Watch ConfigMap (informer)
    KubernetesAPI-->>+controller_manager: Watch ConfigMap (informer)
    KubernetesAPI-->>+controller_manager: Watch Namespace (informer)
    KubernetesAPI-->>+controller_manager: Watch Secret (informer)
    KubernetesAPI-->>+controller_manager: Watch ServiceMonitor (informer)
    KubernetesAPI-->>+controller_manager: Watch ServiceMonitor (informer)
    KubernetesAPI-->>+controller_manager: Watch ClusterServingRuntime (informer)
    KubernetesAPI-->>+controller_manager: Watch Predictor (informer)
    KubernetesAPI-->>+controller_manager: Watch InferenceService (informer)
    KubernetesAPI-->>+controller_manager: Watch InferenceService (informer)

    Note over controller_manager: Exposed Services
    Note right of controller_manager: etcd:2379/TCP [etcd-client-port]
    Note right of controller_manager: modelmesh-controller:8080/TCP []
    Note right of controller_manager: modelmesh-webhook-server-service:9443/TCP []

    Note over KubernetesAPI: Defined CRDs
    Note right of KubernetesAPI: ClusterServingRuntime (serving.kserve.io/v1alpha1)
    Note right of KubernetesAPI: Predictor (serving.kserve.io/v1alpha1)
    Note right of KubernetesAPI: ServingRuntime (serving.kserve.io/v1alpha1)
    Note right of KubernetesAPI: InferenceService (serving.kserve.io/v1beta1)
```

### Webhooks

| Name | Type | Path | Failure Policy | Service | Source |
|------|------|------|----------------|---------|--------|
| servingruntime.modelmesh-webhook-server.default | validating | /validate-serving-modelmesh-io-v1alpha1-servingruntime | Fail | opendatahub/modelmesh-webhook-server-service | [`kustomize:config/overlays/odh (modelmesh-servingruntime.serving.kserve.io)`](https://github.com/kserve/modelmesh-serving/blob/1fcf541d867ceb459fbc76aa1e2bef102c4816db/kustomize:config/overlays/odh (modelmesh-servingruntime.serving.kserve.io)) |

### HTTP Endpoints

| Method | Path | Source |
|--------|------|--------|
| * | /debug/ | [`main.go:339`](https://github.com/kserve/modelmesh-serving/blob/1fcf541d867ceb459fbc76aa1e2bef102c4816db/main.go#L339) |

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

