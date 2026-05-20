# fms-hf-tuning

> **Architecture snapshot: 2026-05-20** (2026-05-20)


**Repository:** red-hat-data-services/fms-hf-tuning  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-20T04:16:47Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 0 |
| Deployments | 0 |
| Services | 0 |
| Secrets | 0 |
| Cluster Roles | 0 |
| Controller Watches | 0 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for fms-hf-tuning

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["fms-hf-tuning Controller"]
        ctrl_1["Controller"]
        class ctrl_1 controller
    end
```

### CRDs

No CRDs found in analyzed sources.

## Dependencies

