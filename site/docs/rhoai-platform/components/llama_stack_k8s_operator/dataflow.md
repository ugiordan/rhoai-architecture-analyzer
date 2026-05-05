# llama-stack-k8s-operator: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

| Type | GVK | Source |
|------|-----|--------|
| For | api/v1alpha1/LlamaStackDistribution | [`controllers/llamastackdistribution_controller.go:590`](https://github.com/llamastack/llama-stack-k8s-operator/blob/916c672901f7e2fc091471677e219830761a532e/controllers/llamastackdistribution_controller.go#L590) |
| Owns | /v1/ConfigMap | [`controllers/llamastackdistribution_controller.go:597`](https://github.com/llamastack/llama-stack-k8s-operator/blob/916c672901f7e2fc091471677e219830761a532e/controllers/llamastackdistribution_controller.go#L597) |
| Owns | /v1/PersistentVolumeClaim | [`controllers/llamastackdistribution_controller.go:605`](https://github.com/llamastack/llama-stack-k8s-operator/blob/916c672901f7e2fc091471677e219830761a532e/controllers/llamastackdistribution_controller.go#L605) |
| Owns | /v1/Service | [`controllers/llamastackdistribution_controller.go:596`](https://github.com/llamastack/llama-stack-k8s-operator/blob/916c672901f7e2fc091471677e219830761a532e/controllers/llamastackdistribution_controller.go#L596) |
| Owns | apps/v1/Deployment | [`controllers/llamastackdistribution_controller.go:593`](https://github.com/llamastack/llama-stack-k8s-operator/blob/916c672901f7e2fc091471677e219830761a532e/controllers/llamastackdistribution_controller.go#L593) |
| Owns | autoscaling/v2/HorizontalPodAutoscaler | [`controllers/llamastackdistribution_controller.go:595`](https://github.com/llamastack/llama-stack-k8s-operator/blob/916c672901f7e2fc091471677e219830761a532e/controllers/llamastackdistribution_controller.go#L595) |
| Owns | networking.k8s.io/v1/Ingress | [`controllers/llamastackdistribution_controller.go:604`](https://github.com/llamastack/llama-stack-k8s-operator/blob/916c672901f7e2fc091471677e219830761a532e/controllers/llamastackdistribution_controller.go#L604) |
| Owns | networking.k8s.io/v1/NetworkPolicy | [`controllers/llamastackdistribution_controller.go:603`](https://github.com/llamastack/llama-stack-k8s-operator/blob/916c672901f7e2fc091471677e219830761a532e/controllers/llamastackdistribution_controller.go#L603) |
| Owns | policy/v1/PodDisruptionBudget | [`controllers/llamastackdistribution_controller.go:594`](https://github.com/llamastack/llama-stack-k8s-operator/blob/916c672901f7e2fc091471677e219830761a532e/controllers/llamastackdistribution_controller.go#L594) |

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for llama-stack-k8s-operator

    participant KubernetesAPI as Kubernetes API
    participant deployment as deployment
    participant llama_stack_k8s_operator_controller_manager as llama-stack-k8s-operator-controller-manager

    KubernetesAPI->>+deployment: Watch LlamaStackDistribution (reconcile)
    deployment->>KubernetesAPI: Create/Update ConfigMap
    deployment->>KubernetesAPI: Create/Update PersistentVolumeClaim
    deployment->>KubernetesAPI: Create/Update Service
    deployment->>KubernetesAPI: Create/Update Deployment
    deployment->>KubernetesAPI: Create/Update HorizontalPodAutoscaler
    deployment->>KubernetesAPI: Create/Update Ingress
    deployment->>KubernetesAPI: Create/Update NetworkPolicy
    deployment->>KubernetesAPI: Create/Update PodDisruptionBudget

    Note over deployment: Exposed Services
    Note right of deployment: llama-stack-k8s-operator-controller-manager-metrics-service:8443/TCP [https]

    Note over KubernetesAPI: Defined CRDs
    Note right of KubernetesAPI: LlamaStackDistribution (llamastack.io/v1alpha1)
```

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

