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
    participant kube_rbac_proxy_verb_override as kube-rbac-proxy-verb-override


    Note over kube_rbac_proxy: Exposed Services
    Note right of kube_rbac_proxy: kube-rbac-proxy:8443/TCP [https]
```

### HTTP Endpoints

| Method | Path | Source |
|--------|------|--------|
| * | / | [`cmd/kube-rbac-proxy/app/kube-rbac-proxy.go:343`](https://github.com/brancz/kube-rbac-proxy/blob/ad283cb8e4e6fe5751f83cbef520c641e74cbf20/cmd/kube-rbac-proxy/app/kube-rbac-proxy.go#L343) |

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

