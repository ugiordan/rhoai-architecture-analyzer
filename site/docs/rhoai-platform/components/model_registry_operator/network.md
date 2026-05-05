# model-registry-operator: Network

## Service Map

*4 unique services (6 total, duplicates from test fixtures collapsed).*

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    model_registry_operator["model-registry-operator"]:::component
    model_registry_operator --> svc_0["model-registry-operator-controller-manager-metrics-service\nClusterIP: 8443/TCP"]:::svc
    model_registry_operator --> svc_1["model-registry-operator-webhook-service\nClusterIP: 443/TCP"]:::svc
    model_registry_operator --> svc_2["template-value\nClusterIP: 0/TCP,0/TCP"]:::svc
    model_registry_operator --> svc_3["template-value-postgres\nClusterIP: 5432/TCP"]:::svc
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| model-registry-operator-controller-manager-metrics-service | ClusterIP | 8443/TCP | [`kustomize:config/overlays/odh`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/kustomize:config/overlays/odh) |
| model-registry-operator-webhook-service | ClusterIP | 443/TCP | [`kustomize:config/overlays/odh`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/kustomize:config/overlays/odh) |
| template-value | ClusterIP | 0/TCP, 0/TCP | [`internal/controller/config/templates/service.yaml.tmpl`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/internal/controller/config/templates/service.yaml.tmpl) |
| template-value | ClusterIP | 0/TCP, 0/TCP | [`internal/controller/config/templates/catalog/catalog-service.yaml.tmpl`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/internal/controller/config/templates/catalog/catalog-service.yaml.tmpl) |
| template-value-postgres | ClusterIP | 5432/TCP | [`internal/controller/config/templates/catalog/catalog-postgres-service.yaml.tmpl`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/internal/controller/config/templates/catalog/catalog-postgres-service.yaml.tmpl) |
| template-value-postgres | ClusterIP | 5432/TCP | [`internal/controller/config/templates/postgres-service.yaml.tmpl`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/internal/controller/config/templates/postgres-service.yaml.tmpl) |

### Ingress / Routing

| Kind | Name | Hosts | Paths | TLS | Source |
|------|------|-------|-------|-----|--------|
| Route | rbac-inferred |  |  | no | [`rbac/manager-role`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/rbac/manager-role) |

### Network Policies

| Name | Policy Types | Source |
|------|-------------|--------|
| template-value-https-route | Ingress | [`internal/controller/config/templates/catalog/catalog-kube-rbac-proxy-network-policy.yaml.tmpl`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/internal/controller/config/templates/catalog/catalog-kube-rbac-proxy-network-policy.yaml.tmpl) |
| template-value-https-route | Ingress | [`internal/controller/config/templates/kube-rbac-proxy/kube-rbac-proxy-network-policy.yaml.tmpl`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/internal/controller/config/templates/kube-rbac-proxy/kube-rbac-proxy-network-policy.yaml.tmpl) |
| template-value-postgres | Ingress | [`internal/controller/config/templates/catalog/catalog-postgres-network-policy.yaml.tmpl`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/internal/controller/config/templates/catalog/catalog-postgres-network-policy.yaml.tmpl) |
| template-value-postgres | Ingress | [`internal/controller/config/templates/postgres-network-policy.yaml.tmpl`](https://github.com/opendatahub-io/model-registry-operator/blob/34bd7584e8cbe37e2f767baa95883d5f3774ca51/internal/controller/config/templates/postgres-network-policy.yaml.tmpl) |

## Network Policy Graph

Visual representation of NetworkPolicy rules. Ingress rules show what traffic is allowed into pods, egress rules show what traffic is allowed out.

```mermaid
graph LR
    classDef policy fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef pod fill:#3498db,stroke:#2980b9,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff

    model_registry_operator["model-registry-operator\nPods"]:::pod
    np_0_template_value_https_route{{"template-value-https-route\nIngress"}}:::policy
    np_0_template_value_https_route --> model_registry_operator
    np_1_template_value_https_route{{"template-value-https-route\nIngress"}}:::policy
    np_1_template_value_https_route --> model_registry_operator
    np_2_template_value_postgres{{"template-value-postgres\nIngress"}}:::policy
    np_2_template_value_postgres --> model_registry_operator
    np_3_template_value_postgres{{"template-value-postgres\nIngress"}}:::policy
    np_3_template_value_postgres --> model_registry_operator
```

