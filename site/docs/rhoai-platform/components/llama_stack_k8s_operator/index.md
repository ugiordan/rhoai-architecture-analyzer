# llama-stack-k8s-operator

> **Architecture snapshot: 2026-05-05** (2026-05-05)


**Repository:** llamastack/llama-stack-k8s-operator  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-05T15:09:11Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 1 |
| Deployments | 2 |
| Services | 1 |
| Secrets | 0 |
| Cluster Roles | 5 |
| Controller Watches | 9 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for llama-stack-k8s-operator

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["llama-stack-k8s-operator Controller"]
        dep_1["deployment"]
        class dep_1 controller
        dep_2["llama-stack-k8s-operator-controller-manager"]
        class dep_2 controller
    end

    crd_LlamaStackDistribution{{"LlamaStackDistribution\nllamastack.io/v1alpha1"}}
    class crd_LlamaStackDistribution crd
    crd_LlamaStackDistribution -->|"For (reconciles)"| controller
    controller -->|"Owns"| owned_3["ConfigMap"]
    class owned_3 owned
    controller -->|"Owns"| owned_4["Deployment"]
    class owned_4 owned
    controller -->|"Owns"| owned_5["HorizontalPodAutoscaler"]
    class owned_5 owned
    controller -->|"Owns"| owned_6["Ingress"]
    class owned_6 owned
    controller -->|"Owns"| owned_7["NetworkPolicy"]
    class owned_7 owned
    controller -->|"Owns"| owned_8["PersistentVolumeClaim"]
    class owned_8 owned
    controller -->|"Owns"| owned_9["PodDisruptionBudget"]
    class owned_9 owned
    controller -->|"Owns"| owned_10["Service"]
    class owned_10 owned
```

### CRDs

| Group | Version | Kind | Scope | Fields | Validation Rules | Source |
|-------|---------|------|-------|--------|------------------|--------|
| llamastack.io | v1alpha1 | LlamaStackDistribution | Namespaced | 371 | 1 | [`config/crd/bases/llamastack.io_llamastackdistributions.yaml`](https://github.com/llamastack/llama-stack-k8s-operator/blob/916c672901f7e2fc091471677e219830761a532e/config/crd/bases/llamastack.io_llamastackdistributions.yaml) |

## Dependencies

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.3 |
| k8s.io/api | v0.34.3 |
| k8s.io/apiextensions-apiserver | v0.34.3 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/client-go | v0.34.3 |
| sigs.k8s.io/controller-runtime | v0.22.4 |

