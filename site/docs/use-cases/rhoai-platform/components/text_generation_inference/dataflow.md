# text-generation-inference: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

No controller watches found.

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for text-generation-inference

    participant KubernetesAPI as Kubernetes API
    participant inference_server as inference-server


    Note over inference_server: Exposed Services
    Note right of inference_server: inference-server:8033/TCP [grpc]
```

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

