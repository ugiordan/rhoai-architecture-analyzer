# llm-d-inference-scheduler: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    llm_d_inference_scheduler["llm-d-inference-scheduler"]:::component
    llm_d_inference_scheduler --> svc_0["${EPP_NAME}\nClusterIP: 5557/TCP,9002/TCP,9090/TCP"]:::svc
    llm_d_inference_scheduler --> svc_1["inference-gateway-istio-nodeport\nNodePort: 15021/TCP,80/TCP"]:::svc
    llm_d_inference_scheduler --> svc_2["istiod-llm-d-gateway\nClusterIP: 15010/TCP,15012/TCP,15014/TCP,443/TCP"]:::svc
    llm_d_inference_scheduler --> svc_3["service\nClusterIP: 8080/TCP"]:::svc
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| ${EPP_NAME} | ClusterIP | 9002/TCP, 5557/TCP, 9090/TCP | [`deploy/components/inference-gateway/services.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/eb2ef5d06644cdf1726fcbc3276d41d8f91f70eb/deploy/components/inference-gateway/services.yaml) |
| inference-gateway-istio-nodeport | NodePort | 15021/TCP, 80/TCP | [`deploy/environments/dev/base-kind-istio/services.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/eb2ef5d06644cdf1726fcbc3276d41d8f91f70eb/deploy/environments/dev/base-kind-istio/services.yaml) |
| istiod-llm-d-gateway | ClusterIP | 15010/TCP, 15012/TCP, 443/TCP, 15014/TCP | [`deploy/components/istio-control-plane/services.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/eb2ef5d06644cdf1726fcbc3276d41d8f91f70eb/deploy/components/istio-control-plane/services.yaml) |
| service | ClusterIP | 8080/TCP | [`deploy/environments/kubernetes-base/common/service.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/eb2ef5d06644cdf1726fcbc3276d41d8f91f70eb/deploy/environments/kubernetes-base/common/service.yaml) |

### Ingress / Routing

| Kind | Name | Hosts | Paths | TLS | Source |
|------|------|-------|-------|-----|--------|
| Gateway | inference-gateway |  |  | no | [`deploy/components/inference-gateway/gateways.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/eb2ef5d06644cdf1726fcbc3276d41d8f91f70eb/deploy/components/inference-gateway/gateways.yaml) |
| HTTPRoute | ${POOL_NAME}-inference-route |  | / | no | [`deploy/components/inference-gateway/httproutes.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/eb2ef5d06644cdf1726fcbc3276d41d8f91f70eb/deploy/components/inference-gateway/httproutes.yaml) |
| Route | route |  |  | yes | [`deploy/environments/kubernetes-base/openshift/route.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/eb2ef5d06644cdf1726fcbc3276d41d8f91f70eb/deploy/environments/kubernetes-base/openshift/route.yaml) |

!!! warning "No Network Policies"
    No NetworkPolicy resources found. All pod-to-pod traffic is allowed by default.

