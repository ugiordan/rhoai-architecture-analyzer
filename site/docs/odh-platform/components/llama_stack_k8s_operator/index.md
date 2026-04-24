# llama-stack-k8s-operator

> **Architecture snapshot: 2026-04-24** (2026-04-24)


**Repository:** llamastack/llama-stack-k8s-operator  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-04-24T08:14:45Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 1 |
| Deployments | 4 |
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
        dep_1["controller-manager"]
        class dep_1 controller
        dep_2["controller-manager"]
        class dep_2 controller
        dep_3["controller-manager"]
        class dep_3 controller
        dep_4["deployment"]
        class dep_4 controller
    end

    crd_LlamaStackDistribution{{"LlamaStackDistribution\nllamastack.io/v1alpha1"}}
    class crd_LlamaStackDistribution crd
    crd_LlamaStackDistribution -->|"For (reconciles)"| controller
    controller -->|"Owns"| owned_5["ConfigMap"]
    class owned_5 owned
    controller -->|"Owns"| owned_6["Deployment"]
    class owned_6 owned
    controller -->|"Owns"| owned_7["HorizontalPodAutoscaler"]
    class owned_7 owned
    controller -->|"Owns"| owned_8["Ingress"]
    class owned_8 owned
    controller -->|"Owns"| owned_9["NetworkPolicy"]
    class owned_9 owned
    controller -->|"Owns"| owned_10["PersistentVolumeClaim"]
    class owned_10 owned
    controller -->|"Owns"| owned_11["PodDisruptionBudget"]
    class owned_11 owned
    controller -->|"Owns"| owned_12["Service"]
    class owned_12 owned
```

### CRDs

| Group | Version | Kind | Scope | Fields | Validation Rules | Source |
|-------|---------|------|-------|--------|------------------|--------|
| llamastack.io | v1alpha1 | LlamaStackDistribution | Namespaced | 371 | 1 | [`config/crd/bases/llamastack.io_llamastackdistributions.yaml`](https://github.com/llamastack/llama-stack-k8s-operator/blob/ba8020a4fc5b6ac86e14aea251992ee2ccdde5ef/config/crd/bases/llamastack.io_llamastackdistributions.yaml) |

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

