# trustyai-service-operator

**Repository:** trustyai-explainability/trustyai-service-operator  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-04-16T15:34:14Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 6 |
| Deployments | 1 |
| Services | 0 |
| Secrets | 0 |
| Cluster Roles | 9 |
| Controller Watches | 11 |

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
        dep_1["controller-manager"]
        class dep_1 controller
    end

    crd_EvalHub{{"EvalHub\ntrustyai.opendatahub.io/v1alpha1"}}
    class crd_EvalHub crd
    crd_EvalHub -->|"For (reconciles)"| controller
    crd_GuardrailsOrchestrator{{"GuardrailsOrchestrator\ntrustyai.opendatahub.io/v1alpha1"}}
    class crd_GuardrailsOrchestrator crd
    crd_GuardrailsOrchestrator -->|"For (reconciles)"| controller
    crd_LMEvalJob{{"LMEvalJob\ntrustyai.opendatahub.io/v1alpha1"}}
    class crd_LMEvalJob crd
    crd_LMEvalJob -->|"For (reconciles)"| controller
    crd_NemoGuardrails{{"NemoGuardrails\ntrustyai.opendatahub.io/v1alpha1"}}
    class crd_NemoGuardrails crd
    crd_NemoGuardrails -->|"For (reconciles)"| controller
    crd_TrustyAIService{{"TrustyAIService\ntrustyai.opendatahub.io/v1"}}
    class crd_TrustyAIService crd
    crd_TrustyAIService -->|"For (reconciles)"| controller
    crd_TrustyAIService{{"TrustyAIService\ntrustyai.opendatahub.io/v1alpha1"}}
    class crd_TrustyAIService crd
    crd_TrustyAIService -->|"For (reconciles)"| controller
    controller -->|"Owns"| owned_2["ConfigMap"]
    class owned_2 owned
    controller -->|"Owns"| owned_3["Deployment"]
    class owned_3 owned
    controller -->|"Owns"| owned_4["Service"]
    class owned_4 owned
    watch_5["InferenceService"] -->|"Watches"| controller
    class watch_5 external
```

### CRDs

| Group | Version | Kind | Scope | Fields | Validation Rules | Source |
|-------|---------|------|-------|--------|------------------|--------|
| trustyai.opendatahub.io | v1alpha1 | EvalHub | Namespaced | 38 | 0 | `config/crd/bases/trustyai.opendatahub.io_evalhubs.yaml` |
| trustyai.opendatahub.io | v1alpha1 | GuardrailsOrchestrator | Namespaced | 70 | 0 | `config/crd/bases/trustyai.opendatahub.io_guardrailsorchestrators.yaml` |
| trustyai.opendatahub.io | v1alpha1 | LMEvalJob | Namespaced | 740 | 0 | `config/crd/bases/trustyai.opendatahub.io_lmevaljobs.yaml` |
| trustyai.opendatahub.io | v1alpha1 | NemoGuardrails | Namespaced | 46 | 0 | `config/crd/bases/trustyai.opendatahub.io_nemoguardrails.yaml` |
| trustyai.opendatahub.io | v1 | TrustyAIService | Namespaced | 26 | 0 | `config/crd/bases/trustyai.opendatahub.io_trustyaiservices.yaml` |
| trustyai.opendatahub.io | v1alpha1 | TrustyAIService | Namespaced | 26 | 0 | `config/crd/bases/trustyai.opendatahub.io_trustyaiservices.yaml` |

## Dependencies

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring | v0.64.1 |
| k8s.io/api | v0.29.2 |
| k8s.io/apimachinery | v0.29.2 |
| k8s.io/client-go | v0.29.2 |
| sigs.k8s.io/controller-runtime | v0.17.0 |
| github.com/go-logr/logr | v1.4.2 |
| github.com/prometheus/client_golang | v1.18.0 |
| k8s.io/apiextensions-apiserver | v0.29.0 |

