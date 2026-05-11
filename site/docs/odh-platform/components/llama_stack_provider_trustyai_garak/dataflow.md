# llama-stack-provider-trustyai-garak: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

No controller watches found.

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for llama-stack-provider-trustyai-garak

    participant KubernetesAPI as Kubernetes API
    participant llama_stack_provider_trustyai_garak as llama-stack-provider-trustyai-garak
```

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

