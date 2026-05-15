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
| model-catalog | ClusterIP | 8080/TCP | [`manifests/kustomize/options/catalog/base/service.yaml`](https://github.com/kubeflow/model-registry/blob/2bccb683c2f077c6d39db5588d7cb908885ac975/manifests/kustomize/options/catalog/base/service.yaml) |

!!! warning "No Network Policies"
    No NetworkPolicy resources found. All pod-to-pod traffic is allowed by default.

