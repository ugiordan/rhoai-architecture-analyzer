# ai4rag: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

No controller watches found in analyzed sources.

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for ai4rag

    participant KubernetesAPI as Kubernetes API
    participant ai4rag as ai4rag
```

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

