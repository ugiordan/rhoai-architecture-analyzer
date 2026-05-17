# feast

> **Architecture snapshot: 2026-05-17** (2026-05-17)


**Repository:** feast-dev/feast  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-17T04:12:42Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 0 |
| Deployments | 2 |
| Services | 1 |
| Secrets | 0 |
| Cluster Roles | 0 |
| Controller Watches | 13 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for feast

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["feast Controller"]
        dep_1["controller-manager"]
        class dep_1 controller
        dep_2["controller-manager"]
        class dep_2 controller
    end

    controller -->|"Owns"| owned_3["ConfigMap"]
    class owned_3 owned
    controller -->|"Owns"| owned_4["CronJob"]
    class owned_4 owned
    controller -->|"Owns"| owned_5["Deployment"]
    class owned_5 owned
    controller -->|"Owns"| owned_6["HorizontalPodAutoscaler"]
    class owned_6 owned
    controller -->|"Owns"| owned_7["PersistentVolumeClaim"]
    class owned_7 owned
    controller -->|"Owns"| owned_8["PodDisruptionBudget"]
    class owned_8 owned
    controller -->|"Owns"| owned_9["Role"]
    class owned_9 owned
    controller -->|"Owns"| owned_10["RoleBinding"]
    class owned_10 owned
    controller -->|"Owns"| owned_11["Route"]
    class owned_11 owned
    controller -->|"Owns"| owned_12["Service"]
    class owned_12 owned
    controller -->|"Owns"| owned_13["ServiceAccount"]
    class owned_13 owned
    watch_14["FeatureStore"] -->|"Watches"| controller
    class watch_14 external
```

### CRDs

No CRDs defined.

## Dependencies

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.2.2 |
| github.com/go-logr/logr | v1.2.2 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/stdr | v1.2.2 |
| github.com/go-logr/stdr | v1.2.2 |
| github.com/go-logr/stdr | v1.2.2 |
| github.com/go-logr/stdr | v1.2.2 |
| github.com/prometheus-operator/prometheus-operator/pkg/client | v0.75.0 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_golang | v1.14.0 |
| github.com/prometheus/client_golang | v1.14.0 |
| github.com/prometheus/client_golang | v1.0.0 |
| github.com/prometheus/client_golang | v1.0.0 |
| github.com/prometheus/client_model | v0.3.0 |
| github.com/prometheus/client_model | v0.6.1 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.3.0 |
| github.com/prometheus/client_model | v0.6.1 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/common | v0.66.1 |
| github.com/prometheus/common | v0.66.1 |
| github.com/prometheus/procfs | v0.16.1 |
| github.com/prometheus/procfs | v0.16.1 |
| google.golang.org/grpc | v1.70.0 |
| google.golang.org/grpc | v1.74.2 |
| google.golang.org/grpc | v1.74.2 |
| google.golang.org/grpc | v1.75.0 |
| google.golang.org/grpc | v1.59.0 |
| google.golang.org/grpc | v1.75.0 |
| google.golang.org/grpc | v1.76.0 |
| google.golang.org/grpc | v1.63.2 |
| google.golang.org/grpc | v1.71.1 |
| google.golang.org/grpc | v1.74.2 |
| google.golang.org/grpc | v1.74.2 |
| google.golang.org/grpc | v1.63.2 |
| google.golang.org/grpc | v1.74.2 |
| google.golang.org/grpc | v1.74.2 |
| google.golang.org/grpc | v1.76.0 |
| google.golang.org/grpc | v1.74.2 |
| google.golang.org/grpc | v1.61.1 |
| google.golang.org/grpc | v1.70.0 |
| google.golang.org/grpc | v1.76.0 |
| google.golang.org/grpc | v1.74.2 |
| google.golang.org/grpc | v1.75.0 |
| google.golang.org/grpc | v1.63.2 |
| google.golang.org/grpc | v1.59.0 |
| google.golang.org/grpc | v1.73.0 |
| google.golang.org/grpc | v1.71.0 |
| google.golang.org/grpc | v1.71.1 |
| google.golang.org/grpc | v1.63.2 |
| google.golang.org/grpc | v1.75.0 |
| google.golang.org/grpc | v1.58.2 |
| google.golang.org/grpc | v1.75.0 |
| google.golang.org/grpc | v1.73.0 |
| google.golang.org/grpc | v1.63.2 |
| google.golang.org/grpc | v1.61.1 |
| google.golang.org/grpc | v1.63.2 |
| google.golang.org/grpc | v1.70.0 |
| google.golang.org/grpc | v1.70.0 |
| google.golang.org/grpc | v1.75.0 |
| google.golang.org/grpc | v1.74.2 |
| google.golang.org/grpc | v1.76.0 |
| google.golang.org/grpc | v1.75.1 |
| google.golang.org/grpc | v1.75.1 |
| google.golang.org/grpc | v1.76.0 |
| google.golang.org/grpc | v1.74.2 |
| google.golang.org/grpc | v1.71.0 |
| google.golang.org/grpc | v1.58.2 |
| google.golang.org/grpc/examples | v0.0.0-20230224211313-3775f633ce20 |
| google.golang.org/grpc/examples | v0.0.0-20230224211313-3775f633ce20 |
| k8s.io/api | v0.33.0 |
| k8s.io/apiextensions-apiserver | v0.33.0 |
| k8s.io/apimachinery | v0.33.0 |
| k8s.io/client-go | v0.33.0 |
| sigs.k8s.io/controller-runtime | v0.21.0 |

