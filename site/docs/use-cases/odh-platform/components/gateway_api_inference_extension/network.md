# gateway-api-inference-extension: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    gateway_api_inference_extension["gateway-api-inference-extension"]:::component
    gateway_api_inference_extension --> svc_0["uvicorn-server\npython-source: 8000/TCP"]:::svc
    gateway_api_inference_extension -.-> ext_mysql[["mysql\ndatabase"]]:::ext
    gateway_api_inference_extension -.-> ext_sqlite[["sqlite\ndatabase"]]:::ext
    gateway_api_inference_extension -.-> ext_grpc[["grpc\ngrpc"]]:::ext
    gateway_api_inference_extension -.-> ext_kafka[["kafka\nmessaging"]]:::ext
    gateway_api_inference_extension -.-> ext_gcs[["gcs\nobject-storage"]]:::ext
    gateway_api_inference_extension -.-> ext_s3[["s3\nobject-storage"]]:::ext
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| uvicorn-server | python-source | 8000/TCP | [`latencypredictor/training_server.py:2171`](https://github.com/kubernetes-sigs/gateway-api-inference-extension/blob/c4c8fef6438746226ed1b7d3cab210229d687f2c/latencypredictor/training_server.py#L2171) |

### Ingress / Routing

| Kind | Name | Hosts | Paths | TLS | Source |
|------|------|-------|-------|-----|--------|
| Gateway | inference-gateway |  |  | no | [`config/manifests/gateway/agentgateway/gateway.yaml`](https://github.com/kubernetes-sigs/gateway-api-inference-extension/blob/c4c8fef6438746226ed1b7d3cab210229d687f2c/config/manifests/gateway/agentgateway/gateway.yaml) |
| Gateway | inference-gateway |  |  | no | [`config/manifests/gateway/envoyaigateway/gateway.yaml`](https://github.com/kubernetes-sigs/gateway-api-inference-extension/blob/c4c8fef6438746226ed1b7d3cab210229d687f2c/config/manifests/gateway/envoyaigateway/gateway.yaml) |
| Gateway | inference-gateway |  |  | no | [`config/manifests/gateway/gke/gateway.yaml`](https://github.com/kubernetes-sigs/gateway-api-inference-extension/blob/c4c8fef6438746226ed1b7d3cab210229d687f2c/config/manifests/gateway/gke/gateway.yaml) |
| Gateway | inference-gateway |  |  | no | [`config/manifests/gateway/istio/gateway.yaml`](https://github.com/kubernetes-sigs/gateway-api-inference-extension/blob/c4c8fef6438746226ed1b7d3cab210229d687f2c/config/manifests/gateway/istio/gateway.yaml) |
| Gateway | inference-gateway |  |  | no | [`config/manifests/gateway/nginxgatewayfabric/gateway.yaml`](https://github.com/kubernetes-sigs/gateway-api-inference-extension/blob/c4c8fef6438746226ed1b7d3cab210229d687f2c/config/manifests/gateway/nginxgatewayfabric/gateway.yaml) |

!!! warning "No Network Policies"
    No NetworkPolicy resources were found in the analyzed sources. Network policies may exist in overlays, Helm values, or cluster-level configurations not captured by static analysis.

