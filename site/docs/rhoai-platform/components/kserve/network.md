# kserve: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    kserve["kserve"]:::component
    kserve --> svc_0["kserve-controller-manager-service\nClusterIP: 8443/TCP"]:::svc
    kserve --> svc_1["kserve-webhook-server-service\nClusterIP: 443/TCP"]:::svc
    kserve --> svc_2["llmisvc-controller-manager-service\nClusterIP: 8443/TCP"]:::svc
    kserve --> svc_3["llmisvc-webhook-server-service\nClusterIP: 443/TCP"]:::svc
    kserve --> svc_4["localmodel-webhook-server-service\nClusterIP: 443/TCP"]:::svc
    kserve --> svc_5["webhook-service\nClusterIP: 443/TCP"]:::svc
    kserve -.-> ext_azure_blob[["azure-blob\nobject-storage"]]:::ext
    kserve -.-> ext_gcs[["gcs\nobject-storage"]]:::ext
    kserve -.-> ext_s3[["s3\nobject-storage"]]:::ext
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| kserve-controller-manager-service | ClusterIP | 8443/TCP | [`config/manager/service.yaml`](https://github.com/kserve/kserve/blob/5d509f23f903a2657dbe2394e785b3bd84c4c40d/config/manager/service.yaml) |
| kserve-webhook-server-service | ClusterIP | 443/TCP | [`config/webhook/service.yaml`](https://github.com/kserve/kserve/blob/5d509f23f903a2657dbe2394e785b3bd84c4c40d/config/webhook/service.yaml) |
| llmisvc-controller-manager-service | ClusterIP | 8443/TCP | [`config/llmisvc/service.yaml`](https://github.com/kserve/kserve/blob/5d509f23f903a2657dbe2394e785b3bd84c4c40d/config/llmisvc/service.yaml) |
| llmisvc-webhook-server-service | ClusterIP | 443/TCP | [`config/webhook/llmisvc/service.yaml`](https://github.com/kserve/kserve/blob/5d509f23f903a2657dbe2394e785b3bd84c4c40d/config/webhook/llmisvc/service.yaml) |
| localmodel-webhook-server-service | ClusterIP | 443/TCP | [`config/webhook/localmodel/service.yaml`](https://github.com/kserve/kserve/blob/5d509f23f903a2657dbe2394e785b3bd84c4c40d/config/webhook/localmodel/service.yaml) |
| webhook-service | ClusterIP | 443/TCP | [`test/webhooks/service.yaml`](https://github.com/kserve/kserve/blob/5d509f23f903a2657dbe2394e785b3bd84c4c40d/test/webhooks/service.yaml) |

### Ingress / Routing

| Kind | Name | Hosts | Paths | TLS | Source |
|------|------|-------|-------|-----|--------|
| Gateway | ai-gateway |  |  | no | [`docs/samples/llmisvc/e2e-gpt-oss/gateway.yaml`](https://github.com/kserve/kserve/blob/5d509f23f903a2657dbe2394e785b3bd84c4c40d/docs/samples/llmisvc/e2e-gpt-oss/gateway.yaml) |
| Gateway | knative-ingress-gateway |  |  | no | [`docs/openshift/serverless/gateways.yaml`](https://github.com/kserve/kserve/blob/5d509f23f903a2657dbe2394e785b3bd84c4c40d/docs/openshift/serverless/gateways.yaml) |
| Gateway | knative-local-gateway |  |  | no | [`docs/openshift/serverless/gateways.yaml`](https://github.com/kserve/kserve/blob/5d509f23f903a2657dbe2394e785b3bd84c4c40d/docs/openshift/serverless/gateways.yaml) |

### Network Policies

| Name | Policy Types | Source |
|------|-------------|--------|
| kserve-controller-manager |  | [`config/overlays/odh/network-policies.yaml`](https://github.com/kserve/kserve/blob/5d509f23f903a2657dbe2394e785b3bd84c4c40d/config/overlays/odh/network-policies.yaml) |

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

