# modelmesh: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

No controller watches found.

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for modelmesh

    participant KubernetesAPI as Kubernetes API
    participant model_mesh as model-mesh


    Note over model_mesh: Exposed Services
    Note right of model_mesh: model-mesh:8033/TCP [grpc]
```

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

