# odh-cli

> **Architecture snapshot: 2026-05-19** (2026-05-19)


**Repository:** opendatahub-io/odh-cli  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-19T04:17:07Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 0 |
| Deployments | 2 |
| Services | 2 |
| Secrets | 0 |
| Cluster Roles | 0 |
| Controller Watches | 82 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for odh-cli

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["odh-cli Controller"]
        dep_1["the-deployment"]
        class dep_1 controller
        dep_2["the-deployment"]
        class dep_2 controller
    end

    watch_3["APIService"] -->|"Watches"| controller
    class watch_3 external
    watch_4["ClusterRole"] -->|"Watches"| controller
    class watch_4 external
    watch_5["ClusterRoleBinding"] -->|"Watches"| controller
    class watch_5 external
    watch_6["ClusterServiceVersion"] -->|"Watches"| controller
    class watch_6 external
    watch_7["ConfigMap"] -->|"Watches"| controller
    class watch_7 external
    watch_8["CustomResourceDefinition"] -->|"Watches"| controller
    class watch_8 external
    watch_9["Deployment"] -->|"Watches"| controller
    class watch_9 external
    watch_10["InstallPlan"] -->|"Watches"| controller
    class watch_10 external
    watch_11["Namespace"] -->|"Watches"| controller
    class watch_11 external
    watch_12["OperatorCondition"] -->|"Watches"| controller
    class watch_12 external
    watch_13["Role"] -->|"Watches"| controller
    class watch_13 external
    watch_14["RoleBinding"] -->|"Watches"| controller
    class watch_14 external
    watch_15["Secret"] -->|"Watches"| controller
    class watch_15 external
    watch_16["Service"] -->|"Watches"| controller
    class watch_16 external
    watch_17["ServiceAccount"] -->|"Watches"| controller
    class watch_17 external
    watch_18["Subscription"] -->|"Watches"| controller
    class watch_18 external
    controller -.->|"depends on"| odh_19["opendatahub-operator"]
    class odh_19 dep
```

### CRDs

No CRDs found in analyzed sources.

## Dependencies

### Internal Platform Dependencies

| Component | Interaction |
|-----------|-------------|
| opendatahub-operator | Go module dependency: github.com/opendatahub-io/opendatahub-operator/pkg/clusterhealth |

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.1 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.1 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/operator-framework/api | v0.39.0 |
| github.com/operator-framework/api | v0.39.0 |
| github.com/operator-framework/api | v0.39.0 |
| github.com/operator-framework/operator-lifecycle-manager | v0.40.0 |
| github.com/operator-framework/operator-registry | v1.63.0 |
| github.com/operator-framework/operator-registry | v1.63.0 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/common | v0.67.5 |
| github.com/prometheus/common | v0.67.5 |
| google.golang.org/grpc | v1.72.2 |
| google.golang.org/grpc | v1.72.2 |
| google.golang.org/grpc | v1.72.2 |
| google.golang.org/grpc | v1.78.0 |
| google.golang.org/grpc | v1.78.0 |
| google.golang.org/grpc | v1.72.2 |
| k8s.io/api | v0.35.2 |
| k8s.io/api | v0.35.2 |
| k8s.io/api | v0.35.2 |
| k8s.io/api | v0.35.2 |
| k8s.io/api | v0.35.0 |
| k8s.io/api | v0.35.0 |
| k8s.io/api | v0.35.0 |
| k8s.io/api | v0.35.2 |
| k8s.io/api | v0.35.0 |
| k8s.io/api | v0.35.2 |
| k8s.io/api | v0.35.2 |
| k8s.io/api | v0.35.2 |
| k8s.io/api | v0.35.2 |
| k8s.io/api | v0.35.2 |
| k8s.io/api | v0.35.0 |
| k8s.io/api | v0.35.2 |
| k8s.io/api | v0.35.0 |
| k8s.io/apiextensions-apiserver | v0.35.0 |
| k8s.io/apiextensions-apiserver | v0.35.0 |
| k8s.io/apiextensions-apiserver | v0.35.0 |
| k8s.io/apiextensions-apiserver | v0.35.0 |
| k8s.io/apiextensions-apiserver | v0.35.0 |
| k8s.io/apiextensions-apiserver | v0.35.2 |
| k8s.io/apiextensions-apiserver | v0.35.0 |
| k8s.io/apimachinery | v0.35.0 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/apimachinery | v0.35.0 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/apimachinery | v0.35.0 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/apimachinery | v0.35.0 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/apimachinery | v0.35.0 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/apimachinery | v0.35.0 |
| k8s.io/apiserver | v0.35.2 |
| k8s.io/apiserver | v0.35.0 |
| k8s.io/apiserver | v0.35.0 |
| k8s.io/apiserver | v0.35.2 |
| k8s.io/apiserver | v0.35.0 |
| k8s.io/apiserver | v0.35.0 |
| k8s.io/client-go | v0.35.0 |
| k8s.io/client-go | v0.35.2 |
| k8s.io/client-go | v0.35.2 |
| k8s.io/client-go | v0.35.2 |
| k8s.io/client-go | v0.35.2 |
| k8s.io/client-go | v0.35.2 |
| k8s.io/client-go | v0.35.0 |
| k8s.io/client-go | v0.35.2 |
| k8s.io/client-go | v0.35.2 |
| k8s.io/client-go | v0.35.0 |
| k8s.io/client-go | v0.35.0 |
| k8s.io/client-go | v0.35.0 |
| k8s.io/client-go | v0.35.2 |
| k8s.io/client-go | v0.35.2 |
| k8s.io/client-go | v0.35.0 |
| sigs.k8s.io/controller-runtime | v0.23.1 |
| sigs.k8s.io/controller-runtime | v0.23.3 |
| sigs.k8s.io/controller-runtime | v0.22.4 |
| sigs.k8s.io/controller-runtime | v0.23.1 |
| sigs.k8s.io/controller-runtime | v0.22.4 |
| sigs.k8s.io/controller-runtime | v0.23.1 |
| sigs.k8s.io/controller-runtime | v0.23.1 |

