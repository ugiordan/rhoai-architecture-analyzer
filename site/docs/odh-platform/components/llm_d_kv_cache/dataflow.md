# llm-d-kv-cache: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

| Type | GVK | Source |
|------|-----|--------|
| For | /v1/Pod | [`examples/kv_events/pod_reconciler/pod_reconciler.go:180`](https://github.com/llm-d/llm-d-kv-cache/blob/ba3a65b3c2a1ef22caa32583f390dcf1970115b2/examples/kv_events/pod_reconciler/pod_reconciler.go#L180) |

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for llm-d-kv-cache

    participant KubernetesAPI as Kubernetes API
    participant n_0 as 0

    KubernetesAPI->>+n_0: Watch Pod (reconcile)
```

### HTTP Endpoints

| Method | Path | Source |
|--------|------|--------|
| * | /metrics | [`examples/kv_events/online/main.go:243`](https://github.com/llm-d/llm-d-kv-cache/blob/ba3a65b3c2a1ef22caa32583f390dcf1970115b2/examples/kv_events/online/main.go#L243) |
| * | /score_chat_completions | [`examples/kv_events/online/main.go:273`](https://github.com/llm-d/llm-d-kv-cache/blob/ba3a65b3c2a1ef22caa32583f390dcf1970115b2/examples/kv_events/online/main.go#L273) |
| * | /score_completions | [`examples/kv_events/online/main.go:247`](https://github.com/llm-d/llm-d-kv-cache/blob/ba3a65b3c2a1ef22caa32583f390dcf1970115b2/examples/kv_events/online/main.go#L247) |

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

### Helm

**Chart:** pvc-evictor v0.1.0

