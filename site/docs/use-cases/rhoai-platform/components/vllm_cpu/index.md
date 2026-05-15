# vllm-cpu

> **Architecture snapshot: 2026-05-15** (2026-05-15)


**Repository:** red-hat-data-services/vllm-cpu  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-15T11:49:02Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 0 |
| Deployments | 0 |
| Services | 5 |
| Secrets | 0 |
| Cluster Roles | 0 |
| Controller Watches | 0 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for vllm-cpu

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["vllm-cpu Controller"]
        ctrl_1["Controller"]
        class ctrl_1 controller
    end
```

### CRDs

No CRDs defined.

## Dependencies

