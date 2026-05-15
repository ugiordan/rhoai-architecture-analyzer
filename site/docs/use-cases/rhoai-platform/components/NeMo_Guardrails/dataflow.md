# NeMo-Guardrails: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

No controller watches found.

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for NeMo-Guardrails

    participant KubernetesAPI as Kubernetes API
    participant NeMo_Guardrails as NeMo-Guardrails


    Note over NeMo_Guardrails: Exposed Services
    Note right of NeMo_Guardrails: env-port-default:1235/TCP []
```

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

