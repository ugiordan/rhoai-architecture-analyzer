# llama-stack-provider-trustyai-garak: Network

### Services

No services defined.

### Network Policies

| Name | Policy Types | Source |
|------|-------------|--------|
| allow-kfp-to-llamastack | Ingress | [`lsd_remote/kfp-setup/kfp-networkpolicy.yaml`](https://github.com/red-hat-data-services/llama-stack-provider-trustyai-garak/blob/37e9b4476992e8313fbeaa0541867097b3d5e9cc/lsd_remote/kfp-setup/kfp-networkpolicy.yaml) |

## Network Policy Graph

Visual representation of NetworkPolicy rules. Ingress rules show what traffic is allowed into pods, egress rules show what traffic is allowed out.

```mermaid
graph LR
    classDef policy fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef pod fill:#3498db,stroke:#2980b9,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff

    llama_stack_provider_trustyai_garak["llama-stack-provider-trustyai-garak\nPods"]:::pod
    np_0_allow_kfp_to_llamastack{{"allow-kfp-to-llamastack\nIngress"}}:::policy
    np_0_allow_kfp_to_llamastack --> llama_stack_provider_trustyai_garak
```

