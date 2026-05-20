# llm-d: Network

### Services

No services found in analyzed sources.

### Ingress / Routing

| Kind | Name | Hosts | Paths | TLS | Source |
|------|------|-------|-------|-----|--------|
| Gateway | llm-d-inference-gateway |  |  | no | [`guides/recipes/gateway/agentgateway/gateway.yaml`](https://github.com/llm-d/llm-d/blob/5bc8871217b23586fb778f24bfbcf41bacc7ec4b/guides/recipes/gateway/agentgateway/gateway.yaml) |
| Gateway | llm-d-inference-gateway |  |  | no | [`guides/recipes/gateway/base/gateway.yaml`](https://github.com/llm-d/llm-d/blob/5bc8871217b23586fb778f24bfbcf41bacc7ec4b/guides/recipes/gateway/base/gateway.yaml) |
| Gateway | llm-d-inference-gateway |  |  | no | [`guides/recipes/gateway/istio/gateway.yaml`](https://github.com/llm-d/llm-d/blob/5bc8871217b23586fb778f24bfbcf41bacc7ec4b/guides/recipes/gateway/istio/gateway.yaml) |
| Gateway | llm-d-inference-gateway |  |  | no | [`guides/recipes/gateway/kgateway/gateway.yaml`](https://github.com/llm-d/llm-d/blob/5bc8871217b23586fb778f24bfbcf41bacc7ec4b/guides/recipes/gateway/kgateway/gateway.yaml) |

!!! warning "No Network Policies"
    No NetworkPolicy resources were found in the analyzed sources. Network policies may exist in overlays, Helm values, or cluster-level configurations not captured by static analysis.

