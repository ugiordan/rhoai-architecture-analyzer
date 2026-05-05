# trustyai-service-operator

> **Architecture snapshot: 2026-05-05** (2026-05-05)


**Repository:** trustyai-explainability/trustyai-service-operator  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-05T15:10:22Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 0 |
| Deployments | 1 |
| Services | 2 |
| Secrets | 0 |
| Cluster Roles | 0 |
| Controller Watches | 14 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for trustyai-service-operator

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["trustyai-service-operator Controller"]
        dep_1["trustyai-service-operator-controller-manager"]
        class dep_1 controller
    end

    controller -->|"Owns"| owned_2["ConfigMap"]
    class owned_2 owned
    controller -->|"Owns"| owned_3["Deployment"]
    class owned_3 owned
    controller -->|"Owns"| owned_4["Service"]
    class owned_4 owned
    watch_5["InferenceService"] -->|"Watches"| controller
    class watch_5 external
    watch_6["Namespace"] -->|"Watches"| controller
    class watch_6 external
```

### CRDs

No CRDs defined.

## Dependencies

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.3 |
| github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring | v0.64.1 |
| github.com/prometheus/client_golang | v1.18.0 |
| k8s.io/api | v0.29.2 |
| k8s.io/apiextensions-apiserver | v0.29.0 |
| k8s.io/apimachinery | v0.29.2 |
| k8s.io/client-go | v0.29.2 |
| sigs.k8s.io/controller-runtime | v0.17.0 |

