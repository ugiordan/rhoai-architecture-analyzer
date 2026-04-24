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
| spark-operator-webhook-svc | ClusterIP | 443/TCP | [`config/webhook/service.yaml`](https://github.com/kubeflow/spark-operator/blob/39b1d20a7fd4163c7c0efa15c3e0194942aa1df1/config/webhook/service.yaml) |

### Network Policies

| Name | Policy Types | Source |
|------|-------------|--------|
| spark-operator-allow-internal | Ingress | [`config/overlays/odh/networkpolicy.yaml`](https://github.com/kubeflow/spark-operator/blob/39b1d20a7fd4163c7c0efa15c3e0194942aa1df1/config/overlays/odh/networkpolicy.yaml) |
| spark-operator-allow-internal | Ingress | [`config/overlays/rhoai/networkpolicy.yaml`](https://github.com/kubeflow/spark-operator/blob/39b1d20a7fd4163c7c0efa15c3e0194942aa1df1/config/overlays/rhoai/networkpolicy.yaml) |

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
    np_1_spark_operator_allow_internal{{"spark-operator-allow-internal\nIngress"}}:::policy
    np_1_spark_operator_allow_internal --> spark_operator
```

