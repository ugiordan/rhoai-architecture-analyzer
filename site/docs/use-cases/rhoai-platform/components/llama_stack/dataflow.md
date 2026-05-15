# llama-stack: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

No controller watches found.

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for llama-stack

    participant KubernetesAPI as Kubernetes API
    participant llama_stack as llama-stack
```

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

