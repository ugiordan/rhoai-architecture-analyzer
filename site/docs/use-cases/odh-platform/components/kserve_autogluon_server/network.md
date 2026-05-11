# kserve-autogluon-server: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    kserve_autogluon_server["kserve-autogluon-server"]:::component
    kserve_autogluon_server --> svc_0["cli-port-default\npython-source: 80/TCP"]:::svc
    kserve_autogluon_server --> svc_1["kserve-controller-manager-metrics-service\nClusterIP: 8443/TCP"]:::svc
    kserve_autogluon_server --> svc_2["kserve-controller-manager-service\nClusterIP: 8443/TCP"]:::svc
    kserve_autogluon_server --> svc_3["kserve-webhook-server-service\nClusterIP: 443/TCP"]:::svc
    kserve_autogluon_server --> svc_4["llmisvc-controller-manager-service\nClusterIP: 8443/TCP"]:::svc
    kserve_autogluon_server --> svc_5["llmisvc-webhook-server-service\nClusterIP: 443/TCP"]:::svc
    kserve_autogluon_server --> svc_6["localmodel-webhook-server-service\nClusterIP: 443/TCP"]:::svc
    kserve_autogluon_server -.-> ext_azure_blob[["azure-blob\nobject-storage"]]:::ext
    kserve_autogluon_server -.-> ext_gcs[["gcs\nobject-storage"]]:::ext
    kserve_autogluon_server -.-> ext_s3[["s3\nobject-storage"]]:::ext
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| cli-port-default | python-source | 80/TCP | [`docs/samples/v1beta1/tensorflow/grpc_client.py:48`](https://github.com/red-hat-data-services/kserve-autogluon-server/blob/047a7264a84c9bc5c2932db3d0e91a02838a4443/docs/samples/v1beta1/tensorflow/grpc_client.py#L48) |
| kserve-controller-manager-metrics-service | ClusterIP | 8443/TCP | [`kustomize:config/overlays/all`](https://github.com/red-hat-data-services/kserve-autogluon-server/blob/047a7264a84c9bc5c2932db3d0e91a02838a4443/kustomize:config/overlays/all) |
| kserve-controller-manager-service | ClusterIP | 8443/TCP | [`kustomize:config/overlays/all`](https://github.com/red-hat-data-services/kserve-autogluon-server/blob/047a7264a84c9bc5c2932db3d0e91a02838a4443/kustomize:config/overlays/all) |
| kserve-webhook-server-service | ClusterIP | 443/TCP | [`kustomize:config/overlays/all`](https://github.com/red-hat-data-services/kserve-autogluon-server/blob/047a7264a84c9bc5c2932db3d0e91a02838a4443/kustomize:config/overlays/all) |
| llmisvc-controller-manager-service | ClusterIP | 8443/TCP | [`kustomize:config/overlays/all`](https://github.com/red-hat-data-services/kserve-autogluon-server/blob/047a7264a84c9bc5c2932db3d0e91a02838a4443/kustomize:config/overlays/all) |
| llmisvc-webhook-server-service | ClusterIP | 443/TCP | [`kustomize:config/overlays/all`](https://github.com/red-hat-data-services/kserve-autogluon-server/blob/047a7264a84c9bc5c2932db3d0e91a02838a4443/kustomize:config/overlays/all) |
| localmodel-webhook-server-service | ClusterIP | 443/TCP | [`kustomize:config/overlays/all`](https://github.com/red-hat-data-services/kserve-autogluon-server/blob/047a7264a84c9bc5c2932db3d0e91a02838a4443/kustomize:config/overlays/all) |

### Ingress / Routing

| Kind | Name | Hosts | Paths | TLS | Source |
|------|------|-------|-------|-----|--------|
| HTTPRoute | rbac-inferred |  |  | no | [`rbac/kserve-manager-role`](https://github.com/red-hat-data-services/kserve-autogluon-server/blob/047a7264a84c9bc5c2932db3d0e91a02838a4443/rbac/kserve-manager-role) |
| Ingress | rbac-inferred |  |  | no | [`rbac/kserve-manager-role`](https://github.com/red-hat-data-services/kserve-autogluon-server/blob/047a7264a84c9bc5c2932db3d0e91a02838a4443/rbac/kserve-manager-role) |
| VirtualService | rbac-inferred |  |  | no | [`rbac/kserve-manager-role`](https://github.com/red-hat-data-services/kserve-autogluon-server/blob/047a7264a84c9bc5c2932db3d0e91a02838a4443/rbac/kserve-manager-role) |

!!! warning "No Network Policies"
    No NetworkPolicy resources found. All pod-to-pod traffic is allowed by default.

