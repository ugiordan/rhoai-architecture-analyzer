# llama-stack-k8s-operator: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    llama_stack_k8s_operator["llama-stack-k8s-operator"]:::component
    llama_stack_k8s_operator --> svc_0["llama-stack-k8s-operator-controller-manager-metrics-service\nClusterIP: 8443/TCP"]:::svc
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| llama-stack-k8s-operator-controller-manager-metrics-service | ClusterIP | 8443/TCP | [`kustomize:config/overlays/odh`](https://github.com/llamastack/llama-stack-k8s-operator/blob/521ca25391e1deca8e192b010c16f86b3c97fbf8/kustomize:config/overlays/odh) |

### Network Policies

| Name | Policy Types | Source |
|------|-------------|--------|
| network-policy | Ingress | [`controllers/manifests/base/networkpolicy.yaml`](https://github.com/llamastack/llama-stack-k8s-operator/blob/521ca25391e1deca8e192b010c16f86b3c97fbf8/controllers/manifests/base/networkpolicy.yaml) |

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

