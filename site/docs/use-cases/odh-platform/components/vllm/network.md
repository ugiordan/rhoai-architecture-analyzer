# vllm: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    vllm["vllm"]:::component
    vllm --> svc_0["cli-port-default\npython-source: 8000/TCP"]:::svc
    vllm --> svc_1["cli-port-default\npython-source: 8001/TCP"]:::svc
    vllm --> svc_2["disagg_prefill_proxy_server-server\npython-source: 8000/TCP"]:::svc
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| cli-port-default | python-source | 8000/TCP | [`benchmarks/benchmark_serving.py:1048`](https://github.com/red-hat-data-services/vllm/blob/7f15f7870a38aa1d06906c25156de8b100da140a/benchmarks/benchmark_serving.py#L1048) |
| cli-port-default | python-source | 8001/TCP | [`examples/online_serving/gradio_openai_chatbot_webserver.py:29`](https://github.com/red-hat-data-services/vllm/blob/7f15f7870a38aa1d06906c25156de8b100da140a/examples/online_serving/gradio_openai_chatbot_webserver.py#L29) |
| disagg_prefill_proxy_server-server | python-source | 8000/TCP | [`benchmarks/disagg_benchmarks/disagg_prefill_proxy_server.py:63`](https://github.com/red-hat-data-services/vllm/blob/7f15f7870a38aa1d06906c25156de8b100da140a/benchmarks/disagg_benchmarks/disagg_prefill_proxy_server.py#L63) |

!!! warning "No Network Policies"
    No NetworkPolicy resources found. All pod-to-pod traffic is allowed by default.

