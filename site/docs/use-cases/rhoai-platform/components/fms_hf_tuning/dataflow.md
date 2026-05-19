# fms-hf-tuning: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

No controller watches found in analyzed sources.

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for fms-hf-tuning

    participant KubernetesAPI as Kubernetes API
    participant fms_hf_tuning as fms-hf-tuning
```

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

