# pipelines-components: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

No controller watches found.

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for pipelines-components

    participant KubernetesAPI as Kubernetes API
    participant pipelines_components as pipelines-components
```

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

