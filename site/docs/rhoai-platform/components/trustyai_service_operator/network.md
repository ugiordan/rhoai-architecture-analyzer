# trustyai-service-operator: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    trustyai_service_operator["trustyai-service-operator"]:::component
    trustyai_service_operator --> svc_0["trustyai-service-operator-controller-manager-metrics-service\nClusterIP: 8443/TCP"]:::svc
    trustyai_service_operator --> svc_1["trustyai-service-operator-metrics-service\nClusterIP: 8080/TCP"]:::svc
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| trustyai-service-operator-controller-manager-metrics-service | ClusterIP | 8443/TCP | [`kustomize:config/overlays/odh`](https://github.com/trustyai-explainability/trustyai-service-operator/blob/6b52d04c51b89713876a2f783e3dd0729ad34065/kustomize:config/overlays/odh) |
| trustyai-service-operator-metrics-service | ClusterIP | 8080/TCP | [`kustomize:config/overlays/odh`](https://github.com/trustyai-explainability/trustyai-service-operator/blob/6b52d04c51b89713876a2f783e3dd0729ad34065/kustomize:config/overlays/odh) |

!!! warning "No Network Policies"
    No NetworkPolicy resources found. All pod-to-pod traffic is allowed by default.

