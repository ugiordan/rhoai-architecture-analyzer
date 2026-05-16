# spark-operator: Network

## Service Map

*2 unique services (3 total, duplicates from test fixtures collapsed).*

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    spark_operator["spark-operator"]:::component
    spark_operator --> svc_0["spark-operator-webhook-svc\nClusterIP: 443/TCP"]:::svc
    spark_operator --> svc_1["the-service\nLoadBalancer: 8666/TCP"]:::svc
    spark_operator -.-> ext_postgres[["postgres\ndatabase"]]:::ext
    spark_operator -.-> ext_sqlite[["sqlite\ndatabase"]]:::ext
    spark_operator -.-> ext_grpc[["grpc\ngrpc"]]:::ext
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| spark-operator-webhook-svc | ClusterIP | 443/TCP | [`kustomize:config/overlays/odh`](https://github.com/kubeflow/spark-operator/blob/e88e255636428e903c0beb8854a8a2870dedc2fd/kustomize:config/overlays/odh) |
| the-service | LoadBalancer | 8666/TCP | [`.gomod-cache/k8s.io/cli-runtime@v0.32.5/artifacts/kustomization/service.yaml`](https://github.com/kubeflow/spark-operator/blob/e88e255636428e903c0beb8854a8a2870dedc2fd/.gomod-cache/k8s.io/cli-runtime@v0.32.5/artifacts/kustomization/service.yaml) |
| the-service | LoadBalancer | 8666/TCP | [`.gopath-loader/pkg/mod/k8s.io/cli-runtime@v0.32.5/artifacts/kustomization/service.yaml`](https://github.com/kubeflow/spark-operator/blob/e88e255636428e903c0beb8854a8a2870dedc2fd/.gopath-loader/pkg/mod/k8s.io/cli-runtime@v0.32.5/artifacts/kustomization/service.yaml) |

### Ingress / Routing

| Kind | Name | Hosts | Paths | TLS | Source |
|------|------|-------|-------|-----|--------|
| Ingress | rbac-inferred |  |  | no | [`rbac/spark-operator-controller`](https://github.com/kubeflow/spark-operator/blob/e88e255636428e903c0beb8854a8a2870dedc2fd/rbac/spark-operator-controller) |

### Network Policies

| Name | Policy Types | Source |
|------|-------------|--------|
| spark-operator-allow-internal | Ingress | [`kustomize:config/overlays/odh`](https://github.com/kubeflow/spark-operator/blob/e88e255636428e903c0beb8854a8a2870dedc2fd/kustomize:config/overlays/odh) |

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

