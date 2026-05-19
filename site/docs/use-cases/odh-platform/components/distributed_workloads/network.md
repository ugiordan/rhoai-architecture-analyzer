# distributed-workloads: Network

## Service Map

*4 unique services (12 total, duplicates from test fixtures collapsed).*

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    distributed_workloads["distributed-workloads"]:::component
    distributed_workloads --> svc_0["kuberay-operator\nClusterIP: 8080/TCP"]:::svc
    distributed_workloads --> svc_1["training-operator\nClusterIP: 8080/TCP"]:::svc
    distributed_workloads --> svc_2["visibility-server\nClusterIP: 443/TCP"]:::svc
    distributed_workloads --> svc_3["webhook-service\nClusterIP: 443/TCP"]:::svc
    distributed_workloads -.-> ext_grpc[["grpc\ngrpc"]]:::ext
    distributed_workloads -.-> ext_azure_blob[["azure-blob\nobject-storage"]]:::ext
    distributed_workloads -.-> ext_minio[["minio\nobject-storage"]]:::ext
    distributed_workloads -.-> ext_s3[["s3\nobject-storage"]]:::ext
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| kuberay-operator | ClusterIP | 8080/TCP | [`.gomod-cache/github.com/ray-project/kuberay/ray-operator@v1.5.1/config/manager/service.yaml`](https://github.com/opendatahub-io/distributed-workloads/blob/c968c77b6e79b132962256e9759655e9173d9dd7/.gomod-cache/github.com/ray-project/kuberay/ray-operator@v1.5.1/config/manager/service.yaml) |
| kuberay-operator | ClusterIP | 8080/TCP | [`.gopath-loader/pkg/mod/github.com/ray-project/kuberay/ray-operator@v1.5.1/config/manager/service.yaml`](https://github.com/opendatahub-io/distributed-workloads/blob/c968c77b6e79b132962256e9759655e9173d9dd7/.gopath-loader/pkg/mod/github.com/ray-project/kuberay/ray-operator@v1.5.1/config/manager/service.yaml) |
| training-operator | ClusterIP | 8080/TCP | [`.gomod-cache/github.com/kubeflow/training-operator@v1.7.0/manifests/base/service.yaml`](https://github.com/opendatahub-io/distributed-workloads/blob/c968c77b6e79b132962256e9759655e9173d9dd7/.gomod-cache/github.com/kubeflow/training-operator@v1.7.0/manifests/base/service.yaml) |
| training-operator | ClusterIP | 8080/TCP | [`.gopath-loader/pkg/mod/github.com/kubeflow/training-operator@v1.7.0/manifests/base/service.yaml`](https://github.com/opendatahub-io/distributed-workloads/blob/c968c77b6e79b132962256e9759655e9173d9dd7/.gopath-loader/pkg/mod/github.com/kubeflow/training-operator@v1.7.0/manifests/base/service.yaml) |
| visibility-server | ClusterIP | 443/TCP | [`.gomod-cache/sigs.k8s.io/kueue@v0.15.7/config/components/visibility/service.yaml`](https://github.com/opendatahub-io/distributed-workloads/blob/c968c77b6e79b132962256e9759655e9173d9dd7/.gomod-cache/sigs.k8s.io/kueue@v0.15.7/config/components/visibility/service.yaml) |
| visibility-server | ClusterIP | 443/TCP | [`.gopath-loader/pkg/mod/sigs.k8s.io/kueue@v0.15.7/config/components/visibility/service.yaml`](https://github.com/opendatahub-io/distributed-workloads/blob/c968c77b6e79b132962256e9759655e9173d9dd7/.gopath-loader/pkg/mod/sigs.k8s.io/kueue@v0.15.7/config/components/visibility/service.yaml) |
| webhook-service | ClusterIP | 443/TCP | [`.gomod-cache/github.com/ray-project/kuberay/ray-operator@v1.5.1/config/webhook/service.yaml`](https://github.com/opendatahub-io/distributed-workloads/blob/c968c77b6e79b132962256e9759655e9173d9dd7/.gomod-cache/github.com/ray-project/kuberay/ray-operator@v1.5.1/config/webhook/service.yaml) |
| webhook-service | ClusterIP | 443/TCP | [`.gomod-cache/sigs.k8s.io/jobset@v0.10.1/config/components/webhook/service.yaml`](https://github.com/opendatahub-io/distributed-workloads/blob/c968c77b6e79b132962256e9759655e9173d9dd7/.gomod-cache/sigs.k8s.io/jobset@v0.10.1/config/components/webhook/service.yaml) |
| webhook-service | ClusterIP | 443/TCP | [`.gomod-cache/sigs.k8s.io/kueue@v0.15.7/config/components/webhook/service.yaml`](https://github.com/opendatahub-io/distributed-workloads/blob/c968c77b6e79b132962256e9759655e9173d9dd7/.gomod-cache/sigs.k8s.io/kueue@v0.15.7/config/components/webhook/service.yaml) |
| webhook-service | ClusterIP | 443/TCP | [`.gopath-loader/pkg/mod/github.com/ray-project/kuberay/ray-operator@v1.5.1/config/webhook/service.yaml`](https://github.com/opendatahub-io/distributed-workloads/blob/c968c77b6e79b132962256e9759655e9173d9dd7/.gopath-loader/pkg/mod/github.com/ray-project/kuberay/ray-operator@v1.5.1/config/webhook/service.yaml) |
| webhook-service | ClusterIP | 443/TCP | [`.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/config/components/webhook/service.yaml`](https://github.com/opendatahub-io/distributed-workloads/blob/c968c77b6e79b132962256e9759655e9173d9dd7/.gopath-loader/pkg/mod/sigs.k8s.io/jobset@v0.10.1/config/components/webhook/service.yaml) |
| webhook-service | ClusterIP | 443/TCP | [`.gopath-loader/pkg/mod/sigs.k8s.io/kueue@v0.15.7/config/components/webhook/service.yaml`](https://github.com/opendatahub-io/distributed-workloads/blob/c968c77b6e79b132962256e9759655e9173d9dd7/.gopath-loader/pkg/mod/sigs.k8s.io/kueue@v0.15.7/config/components/webhook/service.yaml) |

!!! warning "No Network Policies"
    No NetworkPolicy resources were found in the analyzed sources. Network policies may exist in overlays, Helm values, or cluster-level configurations not captured by static analysis.

