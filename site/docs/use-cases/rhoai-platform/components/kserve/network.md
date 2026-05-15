# kserve: Network

## Service Map

*12 unique services (18 total, duplicates from test fixtures collapsed).*

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    kserve["kserve"]:::component
    kserve --> svc_0["cli-port-default\npython-source: 80/TCP"]:::svc
    kserve --> svc_1["keda-admission-webhooks\nClusterIP: 443/TCP,8080/TCP"]:::svc
    kserve --> svc_2["keda-metrics-apiserver\nClusterIP: 443/TCP,8080/TCP"]:::svc
    kserve --> svc_3["keda-operator\nClusterIP: 8080/TCP,9666/TCP"]:::svc
    kserve --> svc_4["kserve-controller-manager-metrics-service\nClusterIP: 8443/TCP"]:::svc
    kserve --> svc_5["kserve-controller-manager-service\nClusterIP: 8443/TCP"]:::svc
    kserve --> svc_6["kserve-webhook-server-service\nClusterIP: 443/TCP"]:::svc
    kserve --> svc_7["llmisvc-controller-manager-service\nClusterIP: 8443/TCP"]:::svc
    kserve --> svc_8["llmisvc-webhook-server-service\nClusterIP: 443/TCP"]:::svc
    kserve --> svc_9["localmodel-webhook-server-service\nClusterIP: 443/TCP"]:::svc
    kserve --> svc_10["uvicorn-server\npython-source: 8000/TCP"]:::svc
    kserve --> svc_11["webhook-service\nClusterIP: 443/TCP"]:::svc
    kserve -.-> ext_etcd[["etcd\ndatabase"]]:::ext
    kserve -.-> ext_mongodb[["mongodb\ndatabase"]]:::ext
    kserve -.-> ext_mysql[["mysql\ndatabase"]]:::ext
    kserve -.-> ext_redis[["redis\ndatabase"]]:::ext
    kserve -.-> ext_grpc[["grpc\ngrpc"]]:::ext
    kserve -.-> ext_azure_blob[["azure-blob\nobject-storage"]]:::ext
    kserve -.-> ext_gcs[["gcs\nobject-storage"]]:::ext
    kserve -.-> ext_s3[["s3\nobject-storage"]]:::ext
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| cli-port-default | python-source | 80/TCP | [`docs/samples/v1beta1/tensorflow/grpc_client.py:48`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/docs/samples/v1beta1/tensorflow/grpc_client.py#L48) |
| keda-admission-webhooks | ClusterIP | 443/TCP, 8080/TCP | [`.gopath-loader/pkg/mod/github.com/kedacore/keda/v2@v2.17.3/config/webhooks/service.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gopath-loader/pkg/mod/github.com/kedacore/keda/v2@v2.17.3/config/webhooks/service.yaml) |
| keda-admission-webhooks | ClusterIP | 443/TCP, 8080/TCP | [`.gomod-cache/github.com/kedacore/keda/v2@v2.17.3/config/webhooks/service.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gomod-cache/github.com/kedacore/keda/v2@v2.17.3/config/webhooks/service.yaml) |
| keda-metrics-apiserver | ClusterIP | 443/TCP, 8080/TCP | [`.gomod-cache/github.com/kedacore/keda/v2@v2.17.3/config/metrics-server/service.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gomod-cache/github.com/kedacore/keda/v2@v2.17.3/config/metrics-server/service.yaml) |
| keda-metrics-apiserver | ClusterIP | 443/TCP, 8080/TCP | [`.gopath-loader/pkg/mod/github.com/kedacore/keda/v2@v2.17.3/config/metrics-server/service.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gopath-loader/pkg/mod/github.com/kedacore/keda/v2@v2.17.3/config/metrics-server/service.yaml) |
| keda-operator | ClusterIP | 9666/TCP, 8080/TCP | [`.gomod-cache/github.com/kedacore/keda/v2@v2.17.3/config/manager/service.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gomod-cache/github.com/kedacore/keda/v2@v2.17.3/config/manager/service.yaml) |
| keda-operator | ClusterIP | 9666/TCP, 8080/TCP | [`.gopath-loader/pkg/mod/github.com/kedacore/keda/v2@v2.17.3/config/manager/service.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gopath-loader/pkg/mod/github.com/kedacore/keda/v2@v2.17.3/config/manager/service.yaml) |
| kserve-controller-manager-metrics-service | ClusterIP | 8443/TCP | [`kustomize:config/overlays/odh`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/kustomize:config/overlays/odh) |
| kserve-controller-manager-service | ClusterIP | 8443/TCP | [`kustomize:config/overlays/odh`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/kustomize:config/overlays/odh) |
| kserve-webhook-server-service | ClusterIP | 443/TCP | [`kustomize:config/overlays/odh`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/kustomize:config/overlays/odh) |
| llmisvc-controller-manager-service | ClusterIP | 8443/TCP | [`kustomize:config/overlays/odh`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/kustomize:config/overlays/odh) |
| llmisvc-webhook-server-service | ClusterIP | 443/TCP | [`kustomize:config/overlays/odh`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/kustomize:config/overlays/odh) |
| localmodel-webhook-server-service | ClusterIP | 443/TCP | [`config/webhook/localmodel/service.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/config/webhook/localmodel/service.yaml) |
| uvicorn-server | python-source | 8000/TCP | [`.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/latencypredictor/training_server.py:1866`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/latencypredictor/training_server.py#L1866) |
| webhook-service | ClusterIP | 443/TCP | [`.gomod-cache/sigs.k8s.io/lws@v0.8.0/config/webhook/service.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gomod-cache/sigs.k8s.io/lws@v0.8.0/config/webhook/service.yaml) |
| webhook-service | ClusterIP | 443/TCP | [`.gopath-loader/pkg/mod/sigs.k8s.io/lws@v0.8.0/config/webhook/service.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gopath-loader/pkg/mod/sigs.k8s.io/lws@v0.8.0/config/webhook/service.yaml) |
| webhook-service | ClusterIP | 443/TCP | [`.gopath-loader/pkg/mod/github.com/open-telemetry/opentelemetry-operator@v0.113.0/config/webhook/service.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gopath-loader/pkg/mod/github.com/open-telemetry/opentelemetry-operator@v0.113.0/config/webhook/service.yaml) |
| webhook-service | ClusterIP | 443/TCP | [`.gomod-cache/github.com/open-telemetry/opentelemetry-operator@v0.113.0/config/webhook/service.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gomod-cache/github.com/open-telemetry/opentelemetry-operator@v0.113.0/config/webhook/service.yaml) |

### Ingress / Routing

| Kind | Name | Hosts | Paths | TLS | Source |
|------|------|-------|-------|-----|--------|
| Gateway | inference-gateway |  |  | no | [`.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/gateway/gke/gateway.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/gateway/gke/gateway.yaml) |
| Gateway | inference-gateway |  |  | no | [`.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/gateway/istio/gateway.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/gateway/istio/gateway.yaml) |
| Gateway | inference-gateway |  |  | no | [`.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/gateway/kgateway/gateway.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/gateway/kgateway/gateway.yaml) |
| Gateway | inference-gateway |  |  | no | [`.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/gateway/nginxgatewayfabric/gateway.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/gateway/nginxgatewayfabric/gateway.yaml) |
| Gateway | inference-gateway |  |  | no | [`.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/gateway/envoyaigateway/gateway.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/gateway/envoyaigateway/gateway.yaml) |
| Gateway | inference-gateway |  |  | no | [`.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/gateway/gke/gateway.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/gateway/gke/gateway.yaml) |
| Gateway | inference-gateway |  |  | no | [`.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/gateway/istio/gateway.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/gateway/istio/gateway.yaml) |
| Gateway | inference-gateway |  |  | no | [`.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/gateway/kgateway/gateway.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/gateway/kgateway/gateway.yaml) |
| Gateway | inference-gateway |  |  | no | [`.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/gateway/nginxgatewayfabric/gateway.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/gateway/nginxgatewayfabric/gateway.yaml) |
| Gateway | inference-gateway |  |  | no | [`.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/gateway/envoyaigateway/gateway.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/gateway/envoyaigateway/gateway.yaml) |
| HTTPRoute | llm-deepseek-route |  | /, /, / | no | [`.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/bbr-example/httproute_bbr_lora.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/bbr-example/httproute_bbr_lora.yaml) |
| HTTPRoute | llm-deepseek-route |  | /, /, / | no | [`.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/bbr-example/httproute_bbr_lora.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/bbr-example/httproute_bbr_lora.yaml) |
| HTTPRoute | llm-llama-route |  | / | no | [`.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/bbr-example/httproute_bbr.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/bbr-example/httproute_bbr.yaml) |
| HTTPRoute | llm-llama-route |  | /, / | no | [`.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/bbr-example/httproute_bbr_lora.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/bbr-example/httproute_bbr_lora.yaml) |
| HTTPRoute | llm-llama-route |  | / | no | [`.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/bbr-example/httproute_bbr.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/bbr-example/httproute_bbr.yaml) |
| HTTPRoute | llm-llama-route |  | /, / | no | [`.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/bbr-example/httproute_bbr_lora.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/bbr-example/httproute_bbr_lora.yaml) |
| HTTPRoute | llm-phi4-route |  | / | no | [`.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/bbr-example/httproute_bbr.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gopath-loader/pkg/mod/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/bbr-example/httproute_bbr.yaml) |
| HTTPRoute | llm-phi4-route |  | / | no | [`.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/bbr-example/httproute_bbr.yaml`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/.gomod-cache/sigs.k8s.io/gateway-api-inference-extension@v1.3.1/config/manifests/bbr-example/httproute_bbr.yaml) |
| Ingress | rbac-inferred |  |  | no | [`rbac/kserve-manager-role`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/rbac/kserve-manager-role) |
| Route | rbac-inferred |  |  | no | [`rbac/kserve-manager-role`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/rbac/kserve-manager-role) |
| VirtualService | rbac-inferred |  |  | no | [`rbac/kserve-manager-role`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/rbac/kserve-manager-role) |

### Network Policies

| Name | Policy Types | Source |
|------|-------------|--------|
| kserve-controller-manager |  | [`kustomize:config/overlays/odh`](https://github.com/kserve/kserve/blob/ee2590545dbe0990eeca74e1918657aeb7b7d7e5/kustomize:config/overlays/odh) |

## Network Policy Graph

Visual representation of NetworkPolicy rules. Ingress rules show what traffic is allowed into pods, egress rules show what traffic is allowed out.

```mermaid
graph LR
    classDef policy fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef pod fill:#3498db,stroke:#2980b9,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff

    kserve["kserve\nPods"]:::pod
    np_0_kserve_controller_manager{{"kserve-controller-manager\nIngress"}}:::policy
    np_0_kserve_controller_manager --> kserve
```

