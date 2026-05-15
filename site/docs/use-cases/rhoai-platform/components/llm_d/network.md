# llm-d: Network

### Services

No services defined.

### Ingress / Routing

| Kind | Name | Hosts | Paths | TLS | Source |
|------|------|-------|-------|-----|--------|
| Gateway | llm-d-inference-gateway |  |  | no | [`guides/recipes/gateway/agentgateway/gateway.yaml`](https://github.com/llm-d/llm-d/blob/3d04e73d481491695c0ffcdfee300628afb3f404/guides/recipes/gateway/agentgateway/gateway.yaml) |
| Gateway | llm-d-inference-gateway |  |  | no | [`guides/recipes/gateway/base/gateway.yaml`](https://github.com/llm-d/llm-d/blob/3d04e73d481491695c0ffcdfee300628afb3f404/guides/recipes/gateway/base/gateway.yaml) |
| Gateway | llm-d-inference-gateway |  |  | no | [`guides/recipes/gateway/istio/gateway.yaml`](https://github.com/llm-d/llm-d/blob/3d04e73d481491695c0ffcdfee300628afb3f404/guides/recipes/gateway/istio/gateway.yaml) |
| Gateway | llm-d-inference-gateway |  |  | no | [`guides/recipes/gateway/kgateway/gateway.yaml`](https://github.com/llm-d/llm-d/blob/3d04e73d481491695c0ffcdfee300628afb3f404/guides/recipes/gateway/kgateway/gateway.yaml) |

!!! warning "No Network Policies"
    No NetworkPolicy resources found. All pod-to-pod traffic is allowed by default.

