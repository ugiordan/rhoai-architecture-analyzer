# notebooks: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

No controller watches found.

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for notebooks

    participant KubernetesAPI as Kubernetes API
    participant notebook as notebook


    Note over notebook: Exposed Services
    Note right of notebook: notebook:8888/TCP []
    Note right of notebook: notebook:8888/TCP []
    Note right of notebook: notebook:8888/TCP []
    Note right of notebook: notebook:8888/TCP []
    Note right of notebook: notebook:8888/TCP []
    Note right of notebook: notebook:8888/TCP []
    Note right of notebook: notebook:8888/TCP []
    Note right of notebook: notebook:8888/TCP []
```

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

