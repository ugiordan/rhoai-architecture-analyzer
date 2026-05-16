# feast: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    feast["feast"]:::component
    feast --> svc_0["uvicorn-server\npython-source: 6566/TCP"]:::svc
    feast -.-> ext_postgres[["postgres\ndatabase"]]:::ext
    feast -.-> ext_redis[["redis\ndatabase"]]:::ext
    feast -.-> ext_sqlite[["sqlite\ndatabase"]]:::ext
    feast -.-> ext_grpc[["grpc\ngrpc"]]:::ext
    feast -.-> ext_gcs[["gcs\nobject-storage"]]:::ext
    feast -.-> ext_s3[["s3\nobject-storage"]]:::ext
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| uvicorn-server | python-source | 6566/TCP | [`infra/scripts/feature_server_docker_smoke.py:38`](https://github.com/feast-dev/feast/blob/0ab134e67b808322415520a6f071e722ef5a9b45/infra/scripts/feature_server_docker_smoke.py#L38) |

!!! warning "No Network Policies"
    No NetworkPolicy resources found. All pod-to-pod traffic is allowed by default.

