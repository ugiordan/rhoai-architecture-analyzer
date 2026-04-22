# kubeflow: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    kubeflow["kubeflow"]:::component
    kubeflow --> svc_0["service\nClusterIP: 443/TCP"]:::svc
    kubeflow --> svc_1["service\nClusterIP: 8080/TCP"]:::svc
    kubeflow --> svc_2["webhook-service\nClusterIP: 443/TCP"]:::svc
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| service | ClusterIP | 443/TCP | [`components/notebook-controller/config/manager/service.yaml`](https://github.com/opendatahub-io/kubeflow/blob/ec116787be830a287a314e370ff732cdfffda873/components/notebook-controller/config/manager/service.yaml) |
| service | ClusterIP | 8080/TCP | [`components/odh-notebook-controller/config/manager/service.yaml`](https://github.com/opendatahub-io/kubeflow/blob/ec116787be830a287a314e370ff732cdfffda873/components/odh-notebook-controller/config/manager/service.yaml) |
| webhook-service | ClusterIP | 443/TCP | [`components/odh-notebook-controller/config/webhook/service.yaml`](https://github.com/opendatahub-io/kubeflow/blob/ec116787be830a287a314e370ff732cdfffda873/components/odh-notebook-controller/config/webhook/service.yaml) |

!!! warning "No Network Policies"
    No NetworkPolicy resources found. All pod-to-pod traffic is allowed by default.

