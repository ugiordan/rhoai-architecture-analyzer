# codeflare-operator

> **Architecture snapshot: 2026-05-05** (2026-05-05)


**Repository:** project-codeflare/codeflare-operator  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-05T15:09:23Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 1 |
| Deployments | 2 |
| Services | 1 |
| Secrets | 1 |
| Cluster Roles | 3 |
| Controller Watches | 8 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for codeflare-operator

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["codeflare-operator Controller"]
        dep_1["manager"]
        class dep_1 controller
        dep_2["manager"]
        class dep_2 controller
    end

    crd_AppWrapper{{"AppWrapper\nworkload.codeflare.dev/v1beta2"}}
    class crd_AppWrapper crd
    controller -->|"Owns"| owned_3["Ingress"]
    class owned_3 owned
    controller -->|"Owns"| owned_4["NetworkPolicy"]
    class owned_4 owned
    controller -->|"Owns"| owned_5["Route"]
    class owned_5 owned
    controller -->|"Owns"| owned_6["Secret"]
    class owned_6 owned
    controller -->|"Owns"| owned_7["Service"]
    class owned_7 owned
    controller -->|"Owns"| owned_8["ServiceAccount"]
    class owned_8 owned
    watch_9["ClusterRoleBinding"] -->|"Watches"| controller
    class watch_9 external
    controller -.->|"depends on"| odh_10["opendatahub-operator"]
    class odh_10 dep
```

### CRDs

| Group | Version | Kind | Scope | Fields | Validation Rules | Source |
|-------|---------|------|-------|--------|------------------|--------|
| workload.codeflare.dev | v1beta2 | AppWrapper | Namespaced | 50 | 0 | [`config/crd/crd-appwrapper.yml`](https://github.com/project-codeflare/codeflare-operator/blob/fb0d403419a114d26adcf65215b6a89e723667d8/config/crd/crd-appwrapper.yml) |

## Dependencies

### Internal Platform Dependencies

| Component | Interaction |
|-----------|-------------|
| opendatahub-operator | Go module dependency: github.com/opendatahub-io/opendatahub-operator/v2 |

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.2 |
| k8s.io/api | v0.31.4 |
| k8s.io/apiextensions-apiserver | v0.31.2 |
| k8s.io/apimachinery | v0.31.4 |
| k8s.io/client-go | v0.31.4 |
| sigs.k8s.io/controller-runtime | v0.19.3 |

