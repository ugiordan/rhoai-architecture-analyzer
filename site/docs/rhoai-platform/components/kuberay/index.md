# kuberay

> **Architecture snapshot: 2026-05-05** (2026-05-05)


**Repository:** ray-project/kuberay  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-05T15:09:52Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 0 |
| Deployments | 2 |
| Services | 2 |
| Secrets | 1 |
| Cluster Roles | 0 |
| Controller Watches | 17 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for kuberay

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["kuberay Controller"]
        dep_1["kuberay-operator"]
        class dep_1 controller
        dep_2["kuberay-operator"]
        class dep_2 controller
    end

    controller -->|"Owns"| owned_3["Job"]
    class owned_3 owned
    controller -->|"Owns"| owned_4["NetworkPolicy"]
    class owned_4 owned
    controller -->|"Owns"| owned_5["Pod"]
    class owned_5 owned
    controller -->|"Owns"| owned_6["RayCluster"]
    class owned_6 owned
    controller -->|"Owns"| owned_7["Route"]
    class owned_7 owned
    controller -->|"Owns"| owned_8["Service"]
    class owned_8 owned
    controller -->|"Owns"| owned_9["ServiceAccount"]
    class owned_9 owned
```

### CRDs

No CRDs defined.

## Dependencies

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/go-logr/zerologr | v1.2.3 |
| github.com/prometheus/client_golang | v1.23.0 |
| github.com/prometheus/client_golang | v1.23.0 |
| google.golang.org/grpc | v1.72.1 |
| k8s.io/api | v0.34.1 |
| k8s.io/api | v0.34.1 |
| k8s.io/apiextensions-apiserver | v0.34.1 |
| k8s.io/apimachinery | v0.34.1 |
| k8s.io/apimachinery | v0.34.1 |
| k8s.io/apiserver | v0.34.1 |
| k8s.io/client-go | v0.34.1 |
| k8s.io/client-go | v0.34.1 |
| sigs.k8s.io/controller-runtime | v0.22.1 |
| sigs.k8s.io/controller-runtime | v0.22.1 |

