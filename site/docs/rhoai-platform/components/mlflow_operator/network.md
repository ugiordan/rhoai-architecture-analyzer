# mlflow-operator: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    mlflow_operator["mlflow-operator"]:::component
    mlflow_operator --> svc_0["minio-service\nClusterIP: 9000/TCP"]:::svc
    mlflow_operator --> svc_1["mlflow-operator-controller-manager-metrics-service\nClusterIP: 8443/TCP"]:::svc
    mlflow_operator --> svc_2["postgres-service\nClusterIP: 5432/TCP"]:::svc
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| minio-service | ClusterIP | 9000/TCP | [`config/seaweedfs/components/tls/service-tls-patch.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/seaweedfs/components/tls/service-tls-patch.yaml) |
| mlflow-operator-controller-manager-metrics-service | ClusterIP | 8443/TCP | [`kustomize:config/overlays/odh`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/kustomize:config/overlays/odh) |
| postgres-service | ClusterIP | 5432/TCP | [`config/postgres/base/service.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/postgres/base/service.yaml) |

### Ingress / Routing

| Kind | Name | Hosts | Paths | TLS | Source |
|------|------|-------|-------|-----|--------|
| HTTPRoute | rbac-inferred |  |  | no | [`rbac/manager-role`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/rbac/manager-role) |

### Network Policies

| Name | Policy Types | Source |
|------|-------------|--------|
| seaweedfs | Ingress | [`config/seaweedfs/base/seaweedfs-networkpolicy.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/seaweedfs/base/seaweedfs-networkpolicy.yaml) |

## Network Policy Graph

Visual representation of NetworkPolicy rules. Ingress rules show what traffic is allowed into pods, egress rules show what traffic is allowed out.

```mermaid
graph LR
    classDef policy fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef pod fill:#3498db,stroke:#2980b9,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff

    mlflow_operator["mlflow-operator\nPods"]:::pod
    np_0_seaweedfs{{"seaweedfs\nIngress"}}:::policy
    np_0_seaweedfs --> mlflow_operator
```

