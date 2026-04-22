# kube-rbac-proxy: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    kube_rbac_proxy["kube-rbac-proxy"]:::component
    kube_rbac_proxy --> svc_0["kube-rbac-proxy\nClusterIP: 8443/TCP"]:::svc
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| kube-rbac-proxy | ClusterIP | 8443/TCP | [`test/kubetest/testtemplates/data/service.yaml`](https://github.com/brancz/kube-rbac-proxy/blob/d1195a883e22af75d26a0dd7e31e6172c659f81c/test/kubetest/testtemplates/data/service.yaml) |

!!! warning "No Network Policies"
    No NetworkPolicy resources found. All pod-to-pod traffic is allowed by default.

