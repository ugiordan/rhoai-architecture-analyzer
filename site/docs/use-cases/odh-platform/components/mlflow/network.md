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
| env-port-default | python-source | 9137/TCP | [`dev/benchmarks/gateway/fake_server.py:66`](https://github.com/opendatahub-io/mlflow/blob/9700acd222c8f8b4300ff2bebf4966f448ebadda/dev/benchmarks/gateway/fake_server.py#L66) |

!!! warning "No Network Policies"
    No NetworkPolicy resources were found in the analyzed sources. Network policies may exist in overlays, Helm values, or cluster-level configurations not captured by static analysis.

