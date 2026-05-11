# modelmesh: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    modelmesh["modelmesh"]:::component
    modelmesh --> svc_0["model-mesh\nClusterIP: 8033/TCP"]:::svc
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| model-mesh | ClusterIP | 8033/TCP | [`config/base/service.yaml`](https://github.com/red-hat-data-services/modelmesh/blob/663e9404150dc48010c5e9263bdbdfd24a561f65/config/base/service.yaml) |

### Network Policies

| Name | Policy Types | Source |
|------|-------------|--------|
| model-mesh | Ingress | [`config/base/networkpolicy.yaml`](https://github.com/red-hat-data-services/modelmesh/blob/663e9404150dc48010c5e9263bdbdfd24a561f65/config/base/networkpolicy.yaml) |

## Network Policy Graph

Visual representation of NetworkPolicy rules. Ingress rules show what traffic is allowed into pods, egress rules show what traffic is allowed out.

```mermaid
graph LR
    classDef policy fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef pod fill:#3498db,stroke:#2980b9,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff

    modelmesh["modelmesh\nPods"]:::pod
    np_0_model_mesh{{"model-mesh\nIngress"}}:::policy
    np_0_model_mesh --> modelmesh
```

