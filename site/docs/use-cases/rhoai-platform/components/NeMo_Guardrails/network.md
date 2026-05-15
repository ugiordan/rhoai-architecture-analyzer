# NeMo-Guardrails: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    NeMo_Guardrails["NeMo-Guardrails"]:::component
    NeMo_Guardrails --> svc_0["env-port-default\npython-source: 1235/TCP"]:::svc
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| env-port-default | python-source | 1235/TCP | [`examples/deployment/gliner_server/src/gliner_server/server.py:51`](https://github.com/red-hat-data-services/NeMo-Guardrails/blob/488345c4a7b04768550f4da8b42493ce77acb6fc/examples/deployment/gliner_server/src/gliner_server/server.py#L51) |

!!! warning "No Network Policies"
    No NetworkPolicy resources found. All pod-to-pod traffic is allowed by default.

