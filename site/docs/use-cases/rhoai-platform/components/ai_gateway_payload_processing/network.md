# ai-gateway-payload-processing: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    ai_gateway_payload_processing["ai-gateway-payload-processing"]:::component
    ai_gateway_payload_processing --> svc_0["uvicorn-server\npython-source: 8000/TCP"]:::svc
    ai_gateway_payload_processing -.-> ext_grpc[["grpc\ngrpc"]]:::ext
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| uvicorn-server | python-source | 8000/TCP | [`.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v0.0.0-20260429190324-8ed5a0cd5d11/latencypredictor/training_server.py:2171`](https://github.com/opendatahub-io/ai-gateway-payload-processing/blob/cfe4df8b2bcfba56ffff6332e46e31a842ec0639/.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v0.0.0-20260429190324-8ed5a0cd5d11/latencypredictor/training_server.py#L2171) |

### Ingress / Routing

| Kind | Name | Hosts | Paths | TLS | Source |
|------|------|-------|-------|-----|--------|
| Gateway | inference-gateway |  |  | no | [`.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v0.0.0-20260429190324-8ed5a0cd5d11/config/manifests/gateway/agentgateway/gateway.yaml`](https://github.com/opendatahub-io/ai-gateway-payload-processing/blob/cfe4df8b2bcfba56ffff6332e46e31a842ec0639/.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v0.0.0-20260429190324-8ed5a0cd5d11/config/manifests/gateway/agentgateway/gateway.yaml) |
| Gateway | inference-gateway |  |  | no | [`.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v0.0.0-20260429190324-8ed5a0cd5d11/config/manifests/gateway/envoyaigateway/gateway.yaml`](https://github.com/opendatahub-io/ai-gateway-payload-processing/blob/cfe4df8b2bcfba56ffff6332e46e31a842ec0639/.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v0.0.0-20260429190324-8ed5a0cd5d11/config/manifests/gateway/envoyaigateway/gateway.yaml) |
| Gateway | inference-gateway |  |  | no | [`.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v0.0.0-20260429190324-8ed5a0cd5d11/config/manifests/gateway/gke/gateway.yaml`](https://github.com/opendatahub-io/ai-gateway-payload-processing/blob/cfe4df8b2bcfba56ffff6332e46e31a842ec0639/.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v0.0.0-20260429190324-8ed5a0cd5d11/config/manifests/gateway/gke/gateway.yaml) |
| Gateway | inference-gateway |  |  | no | [`.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v0.0.0-20260429190324-8ed5a0cd5d11/config/manifests/gateway/istio/gateway.yaml`](https://github.com/opendatahub-io/ai-gateway-payload-processing/blob/cfe4df8b2bcfba56ffff6332e46e31a842ec0639/.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v0.0.0-20260429190324-8ed5a0cd5d11/config/manifests/gateway/istio/gateway.yaml) |
| Gateway | inference-gateway |  |  | no | [`.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v0.0.0-20260429190324-8ed5a0cd5d11/config/manifests/gateway/nginxgatewayfabric/gateway.yaml`](https://github.com/opendatahub-io/ai-gateway-payload-processing/blob/cfe4df8b2bcfba56ffff6332e46e31a842ec0639/.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v0.0.0-20260429190324-8ed5a0cd5d11/config/manifests/gateway/nginxgatewayfabric/gateway.yaml) |
| Gateway | inference-gateway |  |  | no | [`.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v0.0.0-20260429190324-8ed5a0cd5d11/config/manifests/gateway/agentgateway/gateway.yaml`](https://github.com/opendatahub-io/ai-gateway-payload-processing/blob/cfe4df8b2bcfba56ffff6332e46e31a842ec0639/.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v0.0.0-20260429190324-8ed5a0cd5d11/config/manifests/gateway/agentgateway/gateway.yaml) |
| Gateway | inference-gateway |  |  | no | [`.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v0.0.0-20260429190324-8ed5a0cd5d11/config/manifests/gateway/envoyaigateway/gateway.yaml`](https://github.com/opendatahub-io/ai-gateway-payload-processing/blob/cfe4df8b2bcfba56ffff6332e46e31a842ec0639/.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v0.0.0-20260429190324-8ed5a0cd5d11/config/manifests/gateway/envoyaigateway/gateway.yaml) |
| Gateway | inference-gateway |  |  | no | [`.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v0.0.0-20260429190324-8ed5a0cd5d11/config/manifests/gateway/gke/gateway.yaml`](https://github.com/opendatahub-io/ai-gateway-payload-processing/blob/cfe4df8b2bcfba56ffff6332e46e31a842ec0639/.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v0.0.0-20260429190324-8ed5a0cd5d11/config/manifests/gateway/gke/gateway.yaml) |
| Gateway | inference-gateway |  |  | no | [`.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v0.0.0-20260429190324-8ed5a0cd5d11/config/manifests/gateway/istio/gateway.yaml`](https://github.com/opendatahub-io/ai-gateway-payload-processing/blob/cfe4df8b2bcfba56ffff6332e46e31a842ec0639/.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v0.0.0-20260429190324-8ed5a0cd5d11/config/manifests/gateway/istio/gateway.yaml) |
| Gateway | inference-gateway |  |  | no | [`.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v0.0.0-20260429190324-8ed5a0cd5d11/config/manifests/gateway/nginxgatewayfabric/gateway.yaml`](https://github.com/opendatahub-io/ai-gateway-payload-processing/blob/cfe4df8b2bcfba56ffff6332e46e31a842ec0639/.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v0.0.0-20260429190324-8ed5a0cd5d11/config/manifests/gateway/nginxgatewayfabric/gateway.yaml) |

!!! warning "No Network Policies"
    No NetworkPolicy resources were found in the analyzed sources. Network policies may exist in overlays, Helm values, or cluster-level configurations not captured by static analysis.

