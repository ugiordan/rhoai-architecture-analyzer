# modelmesh-serving: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    modelmesh_serving["modelmesh-serving"]:::component
    modelmesh_serving --> svc_0["etcd\nClusterIP: 2379/TCP"]:::svc
    modelmesh_serving --> svc_1["modelmesh-controller\nClusterIP: 8080/TCP"]:::svc
    modelmesh_serving --> svc_2["modelmesh-webhook-server-service\nClusterIP: 9443/TCP"]:::svc
    modelmesh_serving -.-> ext_azure_blob[["azure-blob\nobject-storage"]]:::ext
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| etcd | ClusterIP | 2379/TCP | [`kustomize:config/overlays/odh`](https://github.com/kserve/modelmesh-serving/blob/056bc2e855779c02536db9ef786b26cc73c63f20/kustomize:config/overlays/odh) |
| modelmesh-controller | ClusterIP | 8080/TCP | [`kustomize:config/overlays/odh`](https://github.com/kserve/modelmesh-serving/blob/056bc2e855779c02536db9ef786b26cc73c63f20/kustomize:config/overlays/odh) |
| modelmesh-webhook-server-service | ClusterIP | 9443/TCP | [`kustomize:config/overlays/odh`](https://github.com/kserve/modelmesh-serving/blob/056bc2e855779c02536db9ef786b26cc73c63f20/kustomize:config/overlays/odh) |

### Network Policies

| Name | Policy Types | Source |
|------|-------------|--------|
| etcd | Ingress | [`kustomize:config/overlays/odh`](https://github.com/kserve/modelmesh-serving/blob/056bc2e855779c02536db9ef786b26cc73c63f20/kustomize:config/overlays/odh) |
| modelmesh-controller | Ingress | [`kustomize:config/overlays/odh`](https://github.com/kserve/modelmesh-serving/blob/056bc2e855779c02536db9ef786b26cc73c63f20/kustomize:config/overlays/odh) |
| modelmesh-runtimes | Ingress | [`config/rbac/common/networkpolicy-runtimes.yaml`](https://github.com/kserve/modelmesh-serving/blob/056bc2e855779c02536db9ef786b26cc73c63f20/config/rbac/common/networkpolicy-runtimes.yaml) |
| modelmesh-webhook | Ingress | [`kustomize:config/overlays/odh`](https://github.com/kserve/modelmesh-serving/blob/056bc2e855779c02536db9ef786b26cc73c63f20/kustomize:config/overlays/odh) |

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

