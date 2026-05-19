# model-registry: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    model_registry["model-registry"]:::component
    model_registry --> svc_0["model-catalog\nClusterIP: 8080/TCP"]:::svc
    model_registry -.-> ext_mongodb[["mongodb\ndatabase"]]:::ext
    model_registry -.-> ext_mysql[["mysql\ndatabase"]]:::ext
    model_registry -.-> ext_postgres[["postgres\ndatabase"]]:::ext
    model_registry -.-> ext_sqlite[["sqlite\ndatabase"]]:::ext
    model_registry -.-> ext_gcs[["gcs\nobject-storage"]]:::ext
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| model-catalog | ClusterIP | 8080/TCP | [`manifests/kustomize/options/catalog/base/service.yaml`](https://github.com/kubeflow/model-registry/blob/fd68a656951df0e3e5b24b3d5b3489326d8b3c26/manifests/kustomize/options/catalog/base/service.yaml) |

!!! warning "No Network Policies"
    No NetworkPolicy resources were found in the analyzed sources. Network policies may exist in overlays, Helm values, or cluster-level configurations not captured by static analysis.

