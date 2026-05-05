# data-science-pipelines: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    data_science_pipelines["data-science-pipelines"]:::component
    data_science_pipelines --> svc_0["kubeflow-pipelines-profile-controller\nClusterIP: 80/TCP"]:::svc
    data_science_pipelines --> svc_1["squid\nClusterIP: 3128/TCP"]:::svc
    data_science_pipelines -.-> ext_grpc[["grpc\ngrpc"]]:::ext
    data_science_pipelines -.-> ext_minio[["minio\nobject-storage"]]:::ext
    data_science_pipelines -.-> ext_s3[["s3\nobject-storage"]]:::ext
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| kubeflow-pipelines-profile-controller | ClusterIP | 80/TCP | [`manifests/kustomize/base/installs/multi-user/pipelines-profile-controller/service.yaml`](https://github.com/kubeflow/data-science-pipelines/blob/1e2007f4374655ad9e06fdcfb68a36d0a6fc2d0f/manifests/kustomize/base/installs/multi-user/pipelines-profile-controller/service.yaml) |
| squid | ClusterIP | 3128/TCP | [`.github/resources/squid/manifests/service.yaml`](https://github.com/kubeflow/data-science-pipelines/blob/1e2007f4374655ad9e06fdcfb68a36d0a6fc2d0f/.github/resources/squid/manifests/service.yaml) |

### Ingress / Routing

| Kind | Name | Hosts | Paths | TLS | Source |
|------|------|-------|-------|-----|--------|
| DestinationRule | rbac-inferred |  |  | no | [`rbac/kubeflow-metacontroller`](https://github.com/kubeflow/data-science-pipelines/blob/1e2007f4374655ad9e06fdcfb68a36d0a6fc2d0f/rbac/kubeflow-metacontroller) |
| Route | ml-pipeline-ui |  |  | yes | [`manifests/kustomize/env/openshift/base/route.yaml`](https://github.com/kubeflow/data-science-pipelines/blob/1e2007f4374655ad9e06fdcfb68a36d0a6fc2d0f/manifests/kustomize/env/openshift/base/route.yaml) |

### Network Policies

| Name | Policy Types | Source |
|------|-------------|--------|
| seaweedfs | Ingress | [`manifests/kustomize/third-party/seaweedfs/base/seaweedfs/seaweedfs-networkpolicy.yaml`](https://github.com/kubeflow/data-science-pipelines/blob/1e2007f4374655ad9e06fdcfb68a36d0a6fc2d0f/manifests/kustomize/third-party/seaweedfs/base/seaweedfs/seaweedfs-networkpolicy.yaml) |

## Network Policy Graph

Visual representation of NetworkPolicy rules. Ingress rules show what traffic is allowed into pods, egress rules show what traffic is allowed out.

```mermaid
graph LR
    classDef policy fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef pod fill:#3498db,stroke:#2980b9,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff

    data_science_pipelines["data-science-pipelines\nPods"]:::pod
    np_0_seaweedfs{{"seaweedfs\nIngress"}}:::policy
    np_0_seaweedfs --> data_science_pipelines
```

