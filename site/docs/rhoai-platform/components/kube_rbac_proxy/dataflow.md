# kube-rbac-proxy: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

No controller watches found.

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for kube-rbac-proxy

    participant KubernetesAPI as Kubernetes API
    participant kube_rbac_proxy as kube-rbac-proxy
```

### HTTP Endpoints

| Method | Path | Source |
|--------|------|--------|
| * | / | [`cmd/kube-rbac-proxy/app/kube-rbac-proxy.go:343`](https://github.com/brancz/kube-rbac-proxy/blob/c4cda2828b8cf1151e71fdb2b686b93b4e3c3c67/cmd/kube-rbac-proxy/app/kube-rbac-proxy.go#L343) |

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

