# spark-operator: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    spark_operator["spark-operator"]:::component
    spark_operator --> svc_0["spark-operator-webhook-svc\nClusterIP: 443/TCP"]:::svc
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| spark-operator-webhook-svc | ClusterIP | 443/TCP | [`kustomize:config/overlays/odh`](https://github.com/kubeflow/spark-operator/blob/16adb437ef96672ef47603845e2078e899f3edbe/kustomize:config/overlays/odh) |

### Ingress / Routing

| Kind | Name | Hosts | Paths | TLS | Source |
|------|------|-------|-------|-----|--------|
| Ingress | rbac-inferred |  |  | no | [`rbac/spark-operator-controller`](https://github.com/kubeflow/spark-operator/blob/16adb437ef96672ef47603845e2078e899f3edbe/rbac/spark-operator-controller) |

### Network Policies

| Name | Policy Types | Source |
|------|-------------|--------|
| spark-operator-allow-internal | Ingress | [`kustomize:config/overlays/odh`](https://github.com/kubeflow/spark-operator/blob/16adb437ef96672ef47603845e2078e899f3edbe/kustomize:config/overlays/odh) |

## Network Policy Graph

Visual representation of NetworkPolicy rules. Ingress rules show what traffic is allowed into pods, egress rules show what traffic is allowed out.

```mermaid
graph LR
    classDef policy fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef pod fill:#3498db,stroke:#2980b9,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff

    spark_operator["spark-operator\nPods"]:::pod
    np_0_spark_operator_allow_internal{{"spark-operator-allow-internal\nIngress"}}:::policy
    np_0_spark_operator_allow_internal --> spark_operator
```

