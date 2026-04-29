# workload-variant-autoscaler

> **Architecture snapshot: 2026-04-29** (2026-04-29)


**Repository:** llm-d/workload-variant-autoscaler  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-04-29T11:06:28Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 1 |
| Deployments | 1 |
| Services | 0 |
| Secrets | 2 |
| Cluster Roles | 7 |
| Controller Watches | 4 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for workload-variant-autoscaler

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["workload-variant-autoscaler Controller"]
        dep_1["controller-manager"]
        class dep_1 controller
    end

    crd_VariantAutoscaling{{"VariantAutoscaling\nllmd.ai/v1alpha1"}}
    class crd_VariantAutoscaling crd
    crd_VariantAutoscaling -->|"For (reconciles)"| controller
```

### CRDs

| Group | Version | Kind | Scope | Fields | Validation Rules | Source |
|-------|---------|------|-------|--------|------------------|--------|
| llmd.ai | v1alpha1 | VariantAutoscaling | Namespaced | 26 | 1 | [`config/crd/bases/llmd.ai_variantautoscalings.yaml`](https://github.com/llm-d/workload-variant-autoscaler/blob/958740f7b8a7ac2e9de13b1fbd1edfbfc8b6b782/config/crd/bases/llmd.ai_variantautoscalings.yaml) |

## Dependencies

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.3 |
| github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring | v0.89.0 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/common | v0.67.5 |
| k8s.io/api | v0.34.5 |
| k8s.io/apimachinery | v0.34.5 |
| k8s.io/client-go | v0.34.5 |
| sigs.k8s.io/controller-runtime | v0.22.5 |

