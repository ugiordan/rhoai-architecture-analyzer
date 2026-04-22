# llm-d-inference-scheduler: Dataflow

## Controller Watches

Kubernetes resources this controller monitors for changes. Each watch triggers reconciliation when the watched resource is created, updated, or deleted.

No controller watches found.

## Reconciliation Flow

How the controller interacts with the Kubernetes API during reconciliation.

```mermaid
sequenceDiagram
    %% Static dataflow for llm-d-inference-scheduler

    participant KubernetesAPI as Kubernetes API
    participant n___EPP_NAME_ as ${EPP_NAME}
    participant n___MODEL_NAME_SAFE__vllm_sim as ${MODEL_NAME_SAFE}-vllm-sim
    participant n_0 as 0
    participant e2e_epp as e2e-epp
    participant istiod_llm_d_gateway as istiod-llm-d-gateway
    participant vllm_sim_d as vllm-sim-d
    participant vllm_sim_p as vllm-sim-p


    Note over n___EPP_NAME_: Exposed Services
    Note right of n___EPP_NAME_: ${EPP_NAME}:9002/TCP [default]
    Note right of n___EPP_NAME_: ${EPP_NAME}:5557/TCP [zmq]
    Note right of n___EPP_NAME_: ${EPP_NAME}:9090/TCP [metrics]
    Note right of n___EPP_NAME_: e2e-epp:9002/TCP [ext-proc]
    Note right of n___EPP_NAME_: e2e-epp:5557/TCP [zmq]
    Note right of n___EPP_NAME_: e2e-epp-health:9003/TCP [health]
    Note right of n___EPP_NAME_: e2e-epp-metrics:9090/TCP [metrics]
    Note right of n___EPP_NAME_: inference-gateway-istio-nodeport:15021/TCP [status-port]
    Note right of n___EPP_NAME_: inference-gateway-istio-nodeport:80/TCP [default]
    Note right of n___EPP_NAME_: istiod-llm-d-gateway:15010/TCP [grpc-xds]
    Note right of n___EPP_NAME_: istiod-llm-d-gateway:15012/TCP [https-dns]
    Note right of n___EPP_NAME_: istiod-llm-d-gateway:443/TCP [https-webhook]
    Note right of n___EPP_NAME_: istiod-llm-d-gateway:15014/TCP [http-monitoring]
    Note right of n___EPP_NAME_: service:8080/TCP []
```

### Webhooks

| Name | Type | Path | Failure Policy | Service | Source |
|------|------|------|----------------|---------|--------|
| namespace.sidecar-injector.istio.io | mutating | /inject | Fail | llm-d-istio-system/istiod-llm-d-gateway | [`deploy/components/istio-control-plane/webhooks.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/c7c0201b58d76321e79e12446a5e8d1397e8dcf0/deploy/components/istio-control-plane/webhooks.yaml) |
| object.sidecar-injector.istio.io | mutating | /inject | Fail | llm-d-istio-system/istiod-llm-d-gateway | [`deploy/components/istio-control-plane/webhooks.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/c7c0201b58d76321e79e12446a5e8d1397e8dcf0/deploy/components/istio-control-plane/webhooks.yaml) |
| rev.namespace.sidecar-injector.istio.io | mutating | /inject | Fail | llm-d-istio-system/istiod-llm-d-gateway | [`deploy/components/istio-control-plane/webhooks.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/c7c0201b58d76321e79e12446a5e8d1397e8dcf0/deploy/components/istio-control-plane/webhooks.yaml) |
| rev.object.sidecar-injector.istio.io | mutating | /inject | Fail | llm-d-istio-system/istiod-llm-d-gateway | [`deploy/components/istio-control-plane/webhooks.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/c7c0201b58d76321e79e12446a5e8d1397e8dcf0/deploy/components/istio-control-plane/webhooks.yaml) |
| rev.validation.istio.io | validating | /validate | Ignore | llm-d-istio-system/istiod-llm-d-gateway | [`deploy/components/istio-control-plane/webhooks.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/c7c0201b58d76321e79e12446a5e8d1397e8dcf0/deploy/components/istio-control-plane/webhooks.yaml) |

### HTTP Endpoints

| Method | Path | Source |
|--------|------|--------|
| * | / | [`pkg/sidecar/proxy/proxy.go:310`](https://github.com/llm-d/llm-d-inference-scheduler/blob/c7c0201b58d76321e79e12446a5e8d1397e8dcf0/pkg/sidecar/proxy/proxy.go#L310) |
| * | GET /health | [`pkg/sidecar/proxy/proxy.go:302`](https://github.com/llm-d/llm-d-inference-scheduler/blob/c7c0201b58d76321e79e12446a5e8d1397e8dcf0/pkg/sidecar/proxy/proxy.go#L302) |
| * | POST  | [`pkg/sidecar/proxy/proxy.go:305`](https://github.com/llm-d/llm-d-inference-scheduler/blob/c7c0201b58d76321e79e12446a5e8d1397e8dcf0/pkg/sidecar/proxy/proxy.go#L305) |
| * | POST  | [`pkg/sidecar/proxy/proxy.go:306`](https://github.com/llm-d/llm-d-inference-scheduler/blob/c7c0201b58d76321e79e12446a5e8d1397e8dcf0/pkg/sidecar/proxy/proxy.go#L306) |

## Configuration

ConfigMaps and Helm values that control this component's runtime behavior.

### ConfigMaps

| Name | Data Keys | Source |
|------|-----------|--------|
| istio-llm-d-gateway | mesh, meshNetworks | [`deploy/components/istio-control-plane/configmaps.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/c7c0201b58d76321e79e12446a5e8d1397e8dcf0/deploy/components/istio-control-plane/configmaps.yaml) |
| istio-sidecar-injector-llm-d-gateway | config, values | [`deploy/components/istio-control-plane/configmaps.yaml`](https://github.com/llm-d/llm-d-inference-scheduler/blob/c7c0201b58d76321e79e12446a5e8d1397e8dcf0/deploy/components/istio-control-plane/configmaps.yaml) |

