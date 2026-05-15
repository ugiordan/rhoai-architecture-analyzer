# text-generation-inference: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    text_generation_inference["text-generation-inference"]:::component
    text_generation_inference --> svc_0["inference-server\nClusterIP: 8033/TCP"]:::svc
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| inference-server | ClusterIP | 8033/TCP | [`deployment/base/service.yaml`](https://github.com/red-hat-data-services/text-generation-inference/blob/fded01861025fff09ba5f9a49cda710fcfd3ca93/deployment/base/service.yaml) |

!!! warning "No Network Policies"
    No NetworkPolicy resources found. All pod-to-pod traffic is allowed by default.

