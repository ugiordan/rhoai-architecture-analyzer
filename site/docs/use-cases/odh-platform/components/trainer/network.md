# trainer: Network

## Service Map

*1 unique services (2 total, duplicates from test fixtures collapsed).*

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    trainer["trainer"]:::component
    trainer --> svc_0["webhook-service\nClusterIP: 443/TCP"]:::svc
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| webhook-service | ClusterIP | 443/TCP | [`.gomod-cache/sigs.k8s.io/jobset@v0.10.1/config/components/webhook/service.yaml`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gomod-cache/sigs.k8s.io/jobset@v0.10.1/config/components/webhook/service.yaml) |
| webhook-service | ClusterIP | 443/TCP | [`.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/config/components/webhook/service.yaml`](https://github.com/kubeflow/trainer/blob/51baadf644cd5d2c1672f1c658be46ad82f01b44/.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/config/components/webhook/service.yaml) |

!!! warning "No Network Policies"
    No NetworkPolicy resources found. All pod-to-pod traffic is allowed by default.

