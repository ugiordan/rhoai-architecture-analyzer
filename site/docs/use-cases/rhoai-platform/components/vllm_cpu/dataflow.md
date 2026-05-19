# vllm-cpu: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

No controller watches found in analyzed sources.

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for vllm-cpu

    participant KubernetesAPI as Kubernetes API
    participant vllm_cpu as vllm-cpu


    Note over vllm_cpu: Exposed Services
    Note right of vllm_cpu: cli-port-default:8000/TCP []
    Note right of vllm_cpu: cli-port-default:8006/TCP []
    Note right of vllm_cpu: cli-port-default:8001/TCP []
    Note right of vllm_cpu: disagg_proxy_p2p_nccl_xpyd-server:10001/TCP []
    Note right of vllm_cpu: moriio_toy_proxy_server-server:10001/TCP []
```

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

