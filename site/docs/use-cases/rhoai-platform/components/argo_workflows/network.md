# argo-workflows: Network

## Service Map

*1 unique services (2 total, duplicates from test fixtures collapsed).*

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    argo_workflows["argo-workflows"]:::component
    argo_workflows --> svc_0["the-service\nLoadBalancer: 8666/TCP"]:::svc
    argo_workflows -.-> ext_mysql[["mysql\ndatabase"]]:::ext
    argo_workflows -.-> ext_postgres[["postgres\ndatabase"]]:::ext
    argo_workflows -.-> ext_redis[["redis\ndatabase"]]:::ext
    argo_workflows -.-> ext_sqlite[["sqlite\ndatabase"]]:::ext
    argo_workflows -.-> ext_grpc[["grpc\ngrpc"]]:::ext
    argo_workflows -.-> ext_kafka[["kafka\nmessaging"]]:::ext
    argo_workflows -.-> ext_azure_blob[["azure-blob\nobject-storage"]]:::ext
    argo_workflows -.-> ext_gcs[["gcs\nobject-storage"]]:::ext
    argo_workflows -.-> ext_minio[["minio\nobject-storage"]]:::ext
    argo_workflows -.-> ext_s3[["s3\nobject-storage"]]:::ext
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| the-service | LoadBalancer | 8666/TCP | [`.gomod-cache/k8s.io/cli-runtime@v0.32.2/artifacts/kustomization/service.yaml`](https://github.com/argoproj/argo-workflows/blob/003ed2b35a398772211441cb7c866c51f6f87e2d/.gomod-cache/k8s.io/cli-runtime@v0.32.2/artifacts/kustomization/service.yaml) |
| the-service | LoadBalancer | 8666/TCP | [`.gopath-loader/pkg/mod/k8s.io/cli-runtime@v0.32.2/artifacts/kustomization/service.yaml`](https://github.com/argoproj/argo-workflows/blob/003ed2b35a398772211441cb7c866c51f6f87e2d/.gopath-loader/pkg/mod/k8s.io/cli-runtime@v0.32.2/artifacts/kustomization/service.yaml) |

!!! warning "No Network Policies"
    No NetworkPolicy resources were found in the analyzed sources. Network policies may exist in overlays, Helm values, or cluster-level configurations not captured by static analysis.

