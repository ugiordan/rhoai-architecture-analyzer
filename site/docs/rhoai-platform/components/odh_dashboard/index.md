# odh-dashboard

> **Architecture snapshot: 2026-05-05** (2026-05-05)


**Repository:** red-hat-data-services/odh-dashboard  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-05T15:09:05Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 0 |
| Deployments | 7 |
| Services | 5 |
| Secrets | 2 |
| Cluster Roles | 1 |
| Controller Watches | 5 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for odh-dashboard

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["odh-dashboard Controller"]
        dep_1["odh-dashboard"]
        class dep_1 controller
        dep_2["workspaces-backend"]
        class dep_2 controller
        dep_3["workspaces-controller"]
        class dep_3 controller
        dep_4["workspaces-controller"]
        class dep_4 controller
        dep_5["workspaces-controller"]
        class dep_5 controller
        dep_6["workspaces-controller"]
        class dep_6 controller
        dep_7["workspaces-frontend"]
        class dep_7 controller
    end

    controller -->|"Owns"| owned_8["Service"]
    class owned_8 owned
    controller -->|"Owns"| owned_9["StatefulSet"]
    class owned_9 owned
    controller -->|"Owns"| owned_10["VirtualService"]
    class owned_10 owned
    controller -.->|"depends on"| odh_11["llama-stack-k8s-operator"]
    class odh_11 dep
    controller -.->|"depends on"| odh_12["mlflow-go"]
    class odh_12 dep
    controller -.->|"depends on"| odh_13["mlflow-go"]
    class odh_13 dep
```

### CRDs

No CRDs defined.

## Dependencies

### Internal Platform Dependencies

| Component | Interaction |
|-----------|-------------|
| llama-stack-k8s-operator | Go module dependency: github.com/opendatahub-io/llama-stack-k8s-operator |
| mlflow-go | Go module dependency: github.com/opendatahub-io/mlflow-go |
| mlflow-go | Go module dependency: github.com/opendatahub-io/mlflow-go |

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.2 |
| k8s.io/api | v0.31.0 |
| k8s.io/api | v0.34.1 |
| k8s.io/api | v0.34.3 |
| k8s.io/api | v0.35.4 |
| k8s.io/api | v0.31.0 |
| k8s.io/api | v0.34.3 |
| k8s.io/api | v0.34.1 |
| k8s.io/api | v0.34.1 |
| k8s.io/api | v0.34.3 |
| k8s.io/apiextensions-apiserver | v0.34.3 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apimachinery | v0.34.1 |
| k8s.io/apimachinery | v0.31.0 |
| k8s.io/apimachinery | v0.35.4 |
| k8s.io/apimachinery | v0.31.0 |
| k8s.io/apimachinery | v0.34.1 |
| k8s.io/apimachinery | v0.34.1 |
| k8s.io/apiserver | v0.31.0 |
| k8s.io/client-go | v0.34.1 |
| k8s.io/client-go | v0.34.1 |
| k8s.io/client-go | v0.31.0 |
| k8s.io/client-go | v0.34.1 |
| k8s.io/client-go | v0.35.3 |
| k8s.io/client-go | v0.31.0 |
| k8s.io/client-go | v0.34.3 |
| k8s.io/client-go | v0.34.3 |
| k8s.io/client-go | v0.34.3 |
| sigs.k8s.io/controller-runtime | v0.22.4 |
| sigs.k8s.io/controller-runtime | v0.22.3 |
| sigs.k8s.io/controller-runtime | v0.19.1 |
| sigs.k8s.io/controller-runtime | v0.22.4 |
| sigs.k8s.io/controller-runtime | v0.22.3 |
| sigs.k8s.io/controller-runtime | v0.23.3 |
| sigs.k8s.io/controller-runtime | v0.22.3 |
| sigs.k8s.io/controller-runtime | v0.22.3 |
| sigs.k8s.io/controller-runtime | v0.19.1 |

