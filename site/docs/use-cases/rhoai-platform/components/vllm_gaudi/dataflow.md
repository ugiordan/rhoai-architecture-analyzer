# vllm-gaudi: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

No controller watches found.

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for vllm-gaudi

    participant KubernetesAPI as Kubernetes API
    participant vllm_gaudi as vllm-gaudi


    Note over vllm_gaudi: Exposed Services
    Note right of vllm_gaudi: cli-port-default:8000/TCP []
```

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

