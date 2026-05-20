# modelmesh-serving: Network

## Service Map

*7 unique services (9 total, duplicates from test fixtures collapsed).*

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    modelmesh_serving["modelmesh-serving"]:::component
    modelmesh_serving --> svc_0["cli-port-default\npython-source: 80/TCP"]:::svc
    modelmesh_serving --> svc_1["etcd\nClusterIP: 2379/TCP"]:::svc
    modelmesh_serving --> svc_2["kserve-controller-manager-service\nClusterIP: 8443/TCP"]:::svc
    modelmesh_serving --> svc_3["kserve-webhook-server-service\nClusterIP: 443/TCP"]:::svc
    modelmesh_serving --> svc_4["modelmesh-controller\nClusterIP: 8080/TCP"]:::svc
    modelmesh_serving --> svc_5["modelmesh-webhook-server-service\nClusterIP: 9443/TCP"]:::svc
    modelmesh_serving --> svc_6["models-server\npython-source: 8080/TCP"]:::svc
    modelmesh_serving -.-> ext_etcd[["etcd\ndatabase"]]:::ext
    modelmesh_serving -.-> ext_mysql[["mysql\ndatabase"]]:::ext
    modelmesh_serving -.-> ext_grpc[["grpc\ngrpc"]]:::ext
    modelmesh_serving -.-> ext_azure_blob[["azure-blob\nobject-storage"]]:::ext
    modelmesh_serving -.-> ext_gcs[["gcs\nobject-storage"]]:::ext
    modelmesh_serving -.-> ext_s3[["s3\nobject-storage"]]:::ext
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| cli-port-default | python-source | 80/TCP | [`.gomod-cache/github.com/kserve/kserve@v0.12.0/docs/samples/v1beta1/tensorflow/grpc_client.py:40`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/.gomod-cache/github.com/kserve/kserve@v0.12.0/docs/samples/v1beta1/tensorflow/grpc_client.py#L40) |
| etcd | ClusterIP | 2379/TCP | [`kustomize:config/overlays/odh`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/kustomize:config/overlays/odh) |
| kserve-controller-manager-service | ClusterIP | 8443/TCP | [`.gomod-cache/github.com/kserve/kserve@v0.12.0/config/manager/service.yaml`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/.gomod-cache/github.com/kserve/kserve@v0.12.0/config/manager/service.yaml) |
| kserve-controller-manager-service | ClusterIP | 8443/TCP | [`.gopath-loader/pkg/mod/github.com/kserve/kserve@v0.12.0/config/manager/service.yaml`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/.gopath-loader/pkg/mod/github.com/kserve/kserve@v0.12.0/config/manager/service.yaml) |
| kserve-webhook-server-service | ClusterIP | 443/TCP | [`.gomod-cache/github.com/kserve/kserve@v0.12.0/config/webhook/service.yaml`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/.gomod-cache/github.com/kserve/kserve@v0.12.0/config/webhook/service.yaml) |
| kserve-webhook-server-service | ClusterIP | 443/TCP | [`.gopath-loader/pkg/mod/github.com/kserve/kserve@v0.12.0/config/webhook/service.yaml`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/.gopath-loader/pkg/mod/github.com/kserve/kserve@v0.12.0/config/webhook/service.yaml) |
| modelmesh-controller | ClusterIP | 8080/TCP | [`kustomize:config/overlays/odh`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/kustomize:config/overlays/odh) |
| modelmesh-webhook-server-service | ClusterIP | 9443/TCP | [`kustomize:config/overlays/odh`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/kustomize:config/overlays/odh) |
| models-server | python-source | 8080/TCP | [`.gomod-cache/github.com/kserve/kserve@v0.12.0/docs/samples/fluid/docker/models.py:94`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/.gomod-cache/github.com/kserve/kserve@v0.12.0/docs/samples/fluid/docker/models.py#L94) |

### Network Policies

| Name | Policy Types | Source |
|------|-------------|--------|
| etcd | Ingress | [`kustomize:config/overlays/odh`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/kustomize:config/overlays/odh) |
| modelmesh-controller | Ingress | [`kustomize:config/overlays/odh`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/kustomize:config/overlays/odh) |
| modelmesh-runtimes | Ingress | [`config/rbac/common/networkpolicy-runtimes.yaml`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/config/rbac/common/networkpolicy-runtimes.yaml) |
| modelmesh-webhook | Ingress | [`kustomize:config/overlays/odh`](https://github.com/kserve/modelmesh-serving/blob/4e16417034a9fc02561b6cdb0356d337805589b1/kustomize:config/overlays/odh) |

## Network Policy Graph

Visual representation of NetworkPolicy rules. Ingress rules show what traffic is allowed into pods, egress rules show what traffic is allowed out.

```mermaid
graph LR
    classDef policy fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef pod fill:#3498db,stroke:#2980b9,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff

    modelmesh_serving["modelmesh-serving\nPods"]:::pod
    np_0_etcd{{"etcd\nIngress"}}:::policy
    np_0_etcd --> modelmesh_serving
    np_1_modelmesh_controller{{"modelmesh-controller\nIngress"}}:::policy
    np_1_modelmesh_controller --> modelmesh_serving
    np_2_modelmesh_runtimes{{"modelmesh-runtimes\nIngress"}}:::policy
    np_2_modelmesh_runtimes --> modelmesh_serving
    np_3_modelmesh_webhook{{"modelmesh-webhook\nIngress"}}:::policy
    np_3_modelmesh_webhook --> modelmesh_serving
```

