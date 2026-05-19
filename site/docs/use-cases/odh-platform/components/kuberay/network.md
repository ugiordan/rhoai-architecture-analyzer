# kuberay: Network

## Service Map

*3 unique services (4 total, duplicates from test fixtures collapsed).*

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    kuberay["kuberay"]:::component
    kuberay --> svc_0["kuberay-operator\nClusterIP: 8080/TCP"]:::svc
    kuberay --> svc_1["the-service\nLoadBalancer: 8666/TCP"]:::svc
    kuberay --> svc_2["webhook-service\nClusterIP: 443/TCP"]:::svc
    kuberay -.-> ext_grpc[["grpc\ngrpc"]]:::ext
    kuberay -.-> ext_azure_blob[["azure-blob\nobject-storage"]]:::ext
    kuberay -.-> ext_gcs[["gcs\nobject-storage"]]:::ext
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| kuberay-operator | ClusterIP | 8080/TCP | [`ray-operator/config/manager/service.yaml`](https://github.com/ray-project/kuberay/blob/1466289a1bcff3df4b0be5f0c804e178b7aa8e05/ray-operator/config/manager/service.yaml) |
| the-service | LoadBalancer | 8666/TCP | [`.gomod-cache/k8s.io/cli-runtime@v0.36.0/artifacts/kustomization/service.yaml`](https://github.com/ray-project/kuberay/blob/1466289a1bcff3df4b0be5f0c804e178b7aa8e05/.gomod-cache/k8s.io/cli-runtime@v0.36.0/artifacts/kustomization/service.yaml) |
| the-service | LoadBalancer | 8666/TCP | [`.gopath-loader/pkg/mod/k8s.io/cli-runtime@v0.36.0/artifacts/kustomization/service.yaml`](https://github.com/ray-project/kuberay/blob/1466289a1bcff3df4b0be5f0c804e178b7aa8e05/.gopath-loader/pkg/mod/k8s.io/cli-runtime@v0.36.0/artifacts/kustomization/service.yaml) |
| webhook-service | ClusterIP | 443/TCP | [`ray-operator/config/webhook/service.yaml`](https://github.com/ray-project/kuberay/blob/1466289a1bcff3df4b0be5f0c804e178b7aa8e05/ray-operator/config/webhook/service.yaml) |

!!! warning "No Network Policies"
    No NetworkPolicy resources were found in the analyzed sources. Network policies may exist in overlays, Helm values, or cluster-level configurations not captured by static analysis.

