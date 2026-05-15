# vllm-cpu: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    vllm_cpu["vllm-cpu"]:::component
    vllm_cpu --> svc_0["cli-port-default\npython-source: 8000/TCP"]:::svc
    vllm_cpu --> svc_1["cli-port-default\npython-source: 8006/TCP"]:::svc
    vllm_cpu --> svc_2["cli-port-default\npython-source: 8001/TCP"]:::svc
    vllm_cpu --> svc_3["disagg_proxy_p2p_nccl_xpyd-server\npython-source: 10001/TCP"]:::svc
    vllm_cpu --> svc_4["moriio_toy_proxy_server-server\npython-source: 10001/TCP"]:::svc
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| cli-port-default | python-source | 8000/TCP | [`benchmarks/benchmark_serving_structured_output.py:869`](https://github.com/red-hat-data-services/vllm-cpu/blob/4a21dc6fdc261bc6cd2b1200af5c3a495c5fc29b/benchmarks/benchmark_serving_structured_output.py#L869) |
| cli-port-default | python-source | 8006/TCP | [`examples/online_serving/elastic_ep/scale.py:41`](https://github.com/red-hat-data-services/vllm-cpu/blob/4a21dc6fdc261bc6cd2b1200af5c3a495c5fc29b/examples/online_serving/elastic_ep/scale.py#L41) |
| cli-port-default | python-source | 8001/TCP | [`examples/online_serving/gradio_openai_chatbot_webserver.py:75`](https://github.com/red-hat-data-services/vllm-cpu/blob/4a21dc6fdc261bc6cd2b1200af5c3a495c5fc29b/examples/online_serving/gradio_openai_chatbot_webserver.py#L75) |
| disagg_proxy_p2p_nccl_xpyd-server | python-source | 10001/TCP | [`examples/online_serving/disaggregated_serving_p2p_nccl_xpyd/disagg_proxy_p2p_nccl_xpyd.py:189`](https://github.com/red-hat-data-services/vllm-cpu/blob/4a21dc6fdc261bc6cd2b1200af5c3a495c5fc29b/examples/online_serving/disaggregated_serving_p2p_nccl_xpyd/disagg_proxy_p2p_nccl_xpyd.py#L189) |
| moriio_toy_proxy_server-server | python-source | 10001/TCP | [`examples/online_serving/disaggregated_serving/moriio_toy_proxy_server.py:305`](https://github.com/red-hat-data-services/vllm-cpu/blob/4a21dc6fdc261bc6cd2b1200af5c3a495c5fc29b/examples/online_serving/disaggregated_serving/moriio_toy_proxy_server.py#L305) |

!!! warning "No Network Policies"
    No NetworkPolicy resources found. All pod-to-pod traffic is allowed by default.

