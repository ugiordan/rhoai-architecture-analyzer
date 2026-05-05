# models-as-a-service

> **Architecture snapshot: 2026-05-05** (2026-05-05)


**Repository:** opendatahub-io/models-as-a-service  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-05T15:09:04Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 0 |
| Deployments | 4 |
| Services | 3 |
| Secrets | 1 |
| Cluster Roles | 0 |
| Controller Watches | 11 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for models-as-a-service

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["models-as-a-service Controller"]
        dep_1["maas-api"]
        class dep_1 controller
        dep_2["maas-api"]
        class dep_2 controller
        dep_3["maas-controller"]
        class dep_3 controller
        dep_4["payload-processing"]
        class dep_4 controller
    end

    watch_5["HTTPRoute"] -->|"Watches"| controller
    class watch_5 external
    watch_6["LLMInferenceService"] -->|"Watches"| controller
    class watch_6 external
    watch_7["MaaSModelRef"] -->|"Watches"| controller
    class watch_7 external
    controller -.->|"depends on"| odh_8["kserve"]
    class odh_8 dep
```

### CRDs

No CRDs defined.

## Dependencies

### Internal Platform Dependencies

| Component | Interaction |
|-----------|-------------|
| kserve | Go module dependency: github.com/opendatahub-io/kserve |

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.3 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_model | v0.6.2 |
| k8s.io/api | v0.34.1 |
| k8s.io/api | v0.33.1 |
| k8s.io/apiextensions-apiserver | v0.33.1 |
| k8s.io/apimachinery | v0.33.1 |
| k8s.io/apimachinery | v0.34.1 |
| k8s.io/client-go | v0.34.1 |
| k8s.io/client-go | v0.33.1 |
| sigs.k8s.io/controller-runtime | v0.20.4 |

