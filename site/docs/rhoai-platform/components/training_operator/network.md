# training-operator: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    training_operator["training-operator"]:::component
    training_operator --> svc_0["training-operator\nClusterIP: 443/TCP,8080/TCP"]:::svc
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| training-operator | ClusterIP | 8080/TCP, 443/TCP | [`manifests/base/service.yaml`](https://github.com/kubeflow/training-operator/blob/8582a4b2a238e3552c6b726764580295303a3414/manifests/base/service.yaml) |

!!! warning "No Network Policies"
    No NetworkPolicy resources found. All pod-to-pod traffic is allowed by default.

