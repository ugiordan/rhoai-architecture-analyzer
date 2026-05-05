# odh-model-controller: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    odh_model_controller["odh-model-controller"]:::component
    odh_model_controller --> svc_0["model-serving-api\nClusterIP: 443/TCP,9090/TCP"]:::svc
    odh_model_controller --> svc_1["odh-model-controller-webhook-service\nClusterIP: 443/TCP"]:::svc
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| model-serving-api | ClusterIP | 443/TCP, 9090/TCP | [`config/server/service.yaml`](https://github.com/opendatahub-io/odh-model-controller/blob/55c98bf18a3fa4334d31305c836593a7f6dc4d6d/config/server/service.yaml) |
| odh-model-controller-webhook-service | ClusterIP | 443/TCP | [`config/webhook/service.yaml`](https://github.com/opendatahub-io/odh-model-controller/blob/55c98bf18a3fa4334d31305c836593a7f6dc4d6d/config/webhook/service.yaml) |

### Ingress / Routing

| Kind | Name | Hosts | Paths | TLS | Source |
|------|------|-------|-------|-----|--------|
| Gateway | rbac-inferred |  |  | no | [`rbac/odh-model-controller-role`](https://github.com/opendatahub-io/odh-model-controller/blob/55c98bf18a3fa4334d31305c836593a7f6dc4d6d/rbac/odh-model-controller-role) |
| HTTPRoute | rbac-inferred |  |  | no | [`rbac/odh-model-controller-role`](https://github.com/opendatahub-io/odh-model-controller/blob/55c98bf18a3fa4334d31305c836593a7f6dc4d6d/rbac/odh-model-controller-role) |
| Ingress | rbac-inferred |  |  | no | [`rbac/odh-model-controller-role`](https://github.com/opendatahub-io/odh-model-controller/blob/55c98bf18a3fa4334d31305c836593a7f6dc4d6d/rbac/odh-model-controller-role) |
| Route | rbac-inferred |  |  | no | [`rbac/odh-model-controller-role`](https://github.com/opendatahub-io/odh-model-controller/blob/55c98bf18a3fa4334d31305c836593a7f6dc4d6d/rbac/odh-model-controller-role) |

!!! warning "No Network Policies"
    No NetworkPolicy resources found. All pod-to-pod traffic is allowed by default.

