# vllm: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

No controller watches found.

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for vllm

    participant KubernetesAPI as Kubernetes API
    participant vllm as vllm


    Note over vllm: Exposed Services
    Note right of vllm: cli-port-default:8000/TCP []
    Note right of vllm: cli-port-default:8001/TCP []
    Note right of vllm: disagg_prefill_proxy_server-server:8000/TCP []
```

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

