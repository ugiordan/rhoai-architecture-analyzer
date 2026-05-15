# modelmesh

> **Architecture snapshot: 2026-05-15** (2026-05-15)


**Repository:** opendatahub-io/modelmesh  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-15T11:36:25Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 0 |
| Deployments | 1 |
| Services | 1 |
| Secrets | 0 |
| Cluster Roles | 0 |
| Controller Watches | 0 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for modelmesh

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["modelmesh Controller"]
        dep_1["model-mesh"]
        class dep_1 controller
    end
```

### CRDs

No CRDs defined.

## Dependencies

