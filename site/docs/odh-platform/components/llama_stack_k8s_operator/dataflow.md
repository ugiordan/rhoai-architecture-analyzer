# llama-stack-k8s-operator: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

| Type | GVK | Source |
|------|-----|--------|
| For | api/v1alpha1/LlamaStackDistribution | [`controllers/llamastackdistribution_controller.go:590`](https://github.com/llamastack/llama-stack-k8s-operator/blob/ba8020a4fc5b6ac86e14aea251992ee2ccdde5ef/controllers/llamastackdistribution_controller.go#L590) |
| Owns | /v1/ConfigMap | [`controllers/llamastackdistribution_controller.go:597`](https://github.com/llamastack/llama-stack-k8s-operator/blob/ba8020a4fc5b6ac86e14aea251992ee2ccdde5ef/controllers/llamastackdistribution_controller.go#L597) |
| Owns | /v1/PersistentVolumeClaim | [`controllers/llamastackdistribution_controller.go:605`](https://github.com/llamastack/llama-stack-k8s-operator/blob/ba8020a4fc5b6ac86e14aea251992ee2ccdde5ef/controllers/llamastackdistribution_controller.go#L605) |
| Owns | /v1/Service | [`controllers/llamastackdistribution_controller.go:596`](https://github.com/llamastack/llama-stack-k8s-operator/blob/ba8020a4fc5b6ac86e14aea251992ee2ccdde5ef/controllers/llamastackdistribution_controller.go#L596) |
| Owns | apps/v1/Deployment | [`controllers/llamastackdistribution_controller.go:593`](https://github.com/llamastack/llama-stack-k8s-operator/blob/ba8020a4fc5b6ac86e14aea251992ee2ccdde5ef/controllers/llamastackdistribution_controller.go#L593) |
| Owns | autoscaling/v2/HorizontalPodAutoscaler | [`controllers/llamastackdistribution_controller.go:595`](https://github.com/llamastack/llama-stack-k8s-operator/blob/ba8020a4fc5b6ac86e14aea251992ee2ccdde5ef/controllers/llamastackdistribution_controller.go#L595) |
| Owns | networking.k8s.io/v1/Ingress | [`controllers/llamastackdistribution_controller.go:604`](https://github.com/llamastack/llama-stack-k8s-operator/blob/ba8020a4fc5b6ac86e14aea251992ee2ccdde5ef/controllers/llamastackdistribution_controller.go#L604) |
| Owns | networking.k8s.io/v1/NetworkPolicy | [`controllers/llamastackdistribution_controller.go:603`](https://github.com/llamastack/llama-stack-k8s-operator/blob/ba8020a4fc5b6ac86e14aea251992ee2ccdde5ef/controllers/llamastackdistribution_controller.go#L603) |
| Owns | policy/v1/PodDisruptionBudget | [`controllers/llamastackdistribution_controller.go:594`](https://github.com/llamastack/llama-stack-k8s-operator/blob/ba8020a4fc5b6ac86e14aea251992ee2ccdde5ef/controllers/llamastackdistribution_controller.go#L594) |

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for llama-stack-k8s-operator

    participant KubernetesAPI as Kubernetes API
    participant controller_manager as controller-manager
    participant deployment as deployment

    KubernetesAPI->>+controller_manager: Watch LlamaStackDistribution (reconcile)
    controller_manager->>KubernetesAPI: Create/Update ConfigMap
    controller_manager->>KubernetesAPI: Create/Update PersistentVolumeClaim
    controller_manager->>KubernetesAPI: Create/Update Service
    controller_manager->>KubernetesAPI: Create/Update Deployment
    controller_manager->>KubernetesAPI: Create/Update HorizontalPodAutoscaler
    controller_manager->>KubernetesAPI: Create/Update Ingress
    controller_manager->>KubernetesAPI: Create/Update NetworkPolicy
    controller_manager->>KubernetesAPI: Create/Update PodDisruptionBudget

    Note over controller_manager: Exposed Services
    Note right of controller_manager: service:0/TCP [http]

    Note over KubernetesAPI: Defined CRDs
    Note right of KubernetesAPI: LlamaStackDistribution (llamastack.io/v1alpha1)
```

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

### ConfigMaps

| Name | Data Keys | Source |
|------|-----------|--------|
| llama-stack-config | config.yaml | [`config/samples/example-with-configmap.yaml`](https://github.com/llamastack/llama-stack-k8s-operator/blob/ba8020a4fc5b6ac86e14aea251992ee2ccdde5ef/config/samples/example-with-configmap.yaml) |

