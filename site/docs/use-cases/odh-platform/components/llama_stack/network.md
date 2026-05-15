# llama-stack: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    llama_stack["llama-stack"]:::component
    llama_stack --> svc_0["cli-port-default\npython-source: 8081/TCP"]:::svc
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| cli-port-default | python-source | 8081/TCP | [`benchmarking/k8s-benchmark/openai-mock-server.py:191`](https://github.com/opendatahub-io/llama-stack/blob/2dbe07923373ed11299112e7097e91a882b3e5be/benchmarking/k8s-benchmark/openai-mock-server.py#L191) |

!!! warning "No Network Policies"
    No NetworkPolicy resources found. All pod-to-pod traffic is allowed by default.

