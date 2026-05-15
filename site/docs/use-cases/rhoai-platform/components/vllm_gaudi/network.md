# vllm-gaudi: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    vllm_gaudi["vllm-gaudi"]:::component
    vllm_gaudi --> svc_0["cli-port-default\npython-source: 8000/TCP"]:::svc
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| cli-port-default | python-source | 8000/TCP | [`examples/lmcache/disagg_prefill_lmcache_v1/disagg_proxy_server.py:68`](https://github.com/red-hat-data-services/vllm-gaudi/blob/913c8137213a04b27d908d59479f1394ce8e5e9f/examples/lmcache/disagg_prefill_lmcache_v1/disagg_proxy_server.py#L68) |

!!! warning "No Network Policies"
    No NetworkPolicy resources found. All pod-to-pod traffic is allowed by default.

