# codeflare-operator: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    codeflare_operator["codeflare-operator"]:::component
    codeflare_operator --> svc_0["webhook-service\nClusterIP: 443/TCP"]:::svc
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| webhook-service | ClusterIP | 443/TCP | [`config/webhook/service.yaml`](https://github.com/project-codeflare/codeflare-operator/blob/fb0d403419a114d26adcf65215b6a89e723667d8/config/webhook/service.yaml) |

### Ingress / Routing

| Kind | Name | Hosts | Paths | TLS | Source |
|------|------|-------|-------|-----|--------|
| Ingress | rbac-inferred |  |  | no | [`rbac/manager-role`](https://github.com/project-codeflare/codeflare-operator/blob/fb0d403419a114d26adcf65215b6a89e723667d8/rbac/manager-role) |
| Route | rbac-inferred |  |  | no | [`rbac/manager-role`](https://github.com/project-codeflare/codeflare-operator/blob/fb0d403419a114d26adcf65215b6a89e723667d8/rbac/manager-role) |

!!! warning "No Network Policies"
    No NetworkPolicy resources found. All pod-to-pod traffic is allowed by default.

