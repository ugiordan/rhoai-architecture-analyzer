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


    Note over kube_rbac_proxy: Exposed Services
    Note right of kube_rbac_proxy: kube-rbac-proxy:8443/TCP [https]
```

### HTTP Endpoints

| Method | Path | Source |
|--------|------|--------|
| * | / | [`cmd/kube-rbac-proxy/app/kube-rbac-proxy.go:317`](https://github.com/brancz/kube-rbac-proxy/blob/d1195a883e22af75d26a0dd7e31e6172c659f81c/cmd/kube-rbac-proxy/app/kube-rbac-proxy.go#L317) |

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

