# mlflow: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    mlflow["mlflow"]:::component
    mlflow --> svc_0["env-port-default\npython-source: 9137/TCP"]:::svc
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| env-port-default | python-source | 9137/TCP | [`dev/benchmarks/gateway/fake_server.py:66`](https://github.com/opendatahub-io/mlflow/blob/b4d3741f3c446a3087b1b159bba8303eeb9552bb/dev/benchmarks/gateway/fake_server.py#L66) |

!!! warning "No Network Policies"
    No NetworkPolicy resources found. All pod-to-pod traffic is allowed by default.

