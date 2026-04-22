# llama-stack-k8s-operator: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    llama_stack_k8s_operator["llama-stack-k8s-operator"]:::component
    llama_stack_k8s_operator --> svc_0["service\nClusterIP: 0/TCP"]:::svc
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| service | ClusterIP | 0/TCP | [`controllers/manifests/base/service.yaml`](https://github.com/llamastack/llama-stack-k8s-operator/blob/498ad3f2ac59e93ea7c1ebee43e4c9b27727ddea/controllers/manifests/base/service.yaml) |

### Network Policies

| Name | Policy Types | Source |
|------|-------------|--------|
| network-policy | Ingress | [`controllers/manifests/base/networkpolicy.yaml`](https://github.com/llamastack/llama-stack-k8s-operator/blob/498ad3f2ac59e93ea7c1ebee43e4c9b27727ddea/controllers/manifests/base/networkpolicy.yaml) |

## Network Policy Graph

Visual representation of NetworkPolicy rules. Ingress rules show what traffic is allowed into pods, egress rules show what traffic is allowed out.

```mermaid
graph LR
    classDef policy fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef pod fill:#3498db,stroke:#2980b9,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff

    llama_stack_k8s_operator["llama-stack-k8s-operator\nPods"]:::pod
    np_0_network_policy{{"network-policy\nIngress"}}:::policy
    np_0_network_policy --> llama_stack_k8s_operator
```

