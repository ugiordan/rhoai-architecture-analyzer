# data-science-pipelines-operator: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    data_science_pipelines_operator["data-science-pipelines-operator"]:::component
    data_science_pipelines_operator --> svc_0["mariadb\nClusterIP: 3306/TCP"]:::svc
    data_science_pipelines_operator --> svc_1["minio\nClusterIP: 9000/TCP,9001/TCP"]:::svc
    data_science_pipelines_operator --> svc_2["pypi-server\nClusterIP: 8080/TCP"]:::svc
    data_science_pipelines_operator -.-> ext_mysql[["mysql\ndatabase"]]:::ext
    data_science_pipelines_operator -.-> ext_minio[["minio\nobject-storage"]]:::ext
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| mariadb | ClusterIP | 3306/TCP | [`.github/resources/mariadb/service.yaml`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/b0e56fd54c1f852519894fc869c367bc6818118f/.github/resources/mariadb/service.yaml) |
| minio | ClusterIP | 9000/TCP, 9001/TCP | [`.github/resources/minio/service.yaml`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/b0e56fd54c1f852519894fc869c367bc6818118f/.github/resources/minio/service.yaml) |
| pypi-server | ClusterIP | 8080/TCP | [`.github/resources/pypiserver/base/service.yaml`](https://github.com/opendatahub-io/data-science-pipelines-operator/blob/b0e56fd54c1f852519894fc869c367bc6818118f/.github/resources/pypiserver/base/service.yaml) |

!!! warning "No Network Policies"
    No NetworkPolicy resources found. All pod-to-pod traffic is allowed by default.

