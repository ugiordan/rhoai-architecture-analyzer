# opendatahub-operator

> **Architecture snapshot: 2026-05-05** (2026-05-05)


**Repository:** opendatahub-io/opendatahub-operator  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-05T15:10:56Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 0 |
| Deployments | 19 |
| Services | 3 |
| Secrets | 2 |
| Cluster Roles | 23 |
| Controller Watches | 194 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for opendatahub-operator

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["opendatahub-operator Controller"]
        dep_1["azure-cloud-manager-operator"]
        class dep_1 controller
        dep_2["azure-cloud-manager-operator"]
        class dep_2 controller
        dep_3["azure-cloud-manager-operator"]
        class dep_3 controller
        dep_4["controller-manager"]
        class dep_4 controller
        dep_5["controller-manager"]
        class dep_5 controller
        dep_6["controller-manager"]
        class dep_6 controller
        dep_7["controller-manager"]
        class dep_7 controller
        dep_8["controller-manager"]
        class dep_8 controller
        dep_9["controller-manager"]
        class dep_9 controller
        dep_10["controller-manager"]
        class dep_10 controller
        dep_11["coreweave-cloud-manager-operator"]
        class dep_11 controller
        dep_12["coreweave-cloud-manager-operator"]
        class dep_12 controller
        dep_13["coreweave-cloud-manager-operator"]
        class dep_13 controller
        dep_14["rhods-operator"]
        class dep_14 controller
        dep_15["rhods-operator"]
        class dep_15 controller
        dep_16["rhods-operator"]
        class dep_16 controller
        dep_17["rhods-operator"]
        class dep_17 controller
        dep_18["rhods-operator"]
        class dep_18 controller
        dep_19["rhods-operator"]
        class dep_19 controller
    end

    controller -->|"Owns"| owned_20["ClusterRole"]
    class owned_20 owned
    controller -->|"Owns"| owned_21["ClusterRoleBinding"]
    class owned_21 owned
    controller -->|"Owns"| owned_22["ConfigMap"]
    class owned_22 owned
    controller -->|"Owns"| owned_23["ConsoleLink"]
    class owned_23 owned
    controller -->|"Owns"| owned_24["Dashboard"]
    class owned_24 owned
    controller -->|"Owns"| owned_25["DataSciencePipelines"]
    class owned_25 owned
    controller -->|"Owns"| owned_26["Deployment"]
    class owned_26 owned
    controller -->|"Owns"| owned_27["FeastOperator"]
    class owned_27 owned
    controller -->|"Owns"| owned_28["HTTPRoute"]
    class owned_28 owned
    controller -->|"Owns"| owned_29["Job"]
    class owned_29 owned
    controller -->|"Owns"| owned_30["Kserve"]
    class owned_30 owned
    controller -->|"Owns"| owned_31["Kueue"]
    class owned_31 owned
    controller -->|"Owns"| owned_32["LlamaStackOperator"]
    class owned_32 owned
    controller -->|"Owns"| owned_33["MLflowOperator"]
    class owned_33 owned
    controller -->|"Owns"| owned_34["ModelController"]
    class owned_34 owned
    controller -->|"Owns"| owned_35["ModelRegistry"]
    class owned_35 owned
    controller -->|"Owns"| owned_36["MutatingWebhookConfiguration"]
    class owned_36 owned
    controller -->|"Owns"| owned_37["NetworkPolicy"]
    class owned_37 owned
    controller -->|"Owns"| owned_38["PodDisruptionBudget"]
    class owned_38 owned
    controller -->|"Owns"| owned_39["PodMonitor"]
    class owned_39 owned
    controller -->|"Owns"| owned_40["PrometheusRule"]
    class owned_40 owned
    controller -->|"Owns"| owned_41["Ray"]
    class owned_41 owned
    controller -->|"Owns"| owned_42["Role"]
    class owned_42 owned
    controller -->|"Owns"| owned_43["RoleBinding"]
    class owned_43 owned
    controller -->|"Owns"| owned_44["Route"]
    class owned_44 owned
    controller -->|"Owns"| owned_45["Secret"]
    class owned_45 owned
    controller -->|"Owns"| owned_46["SecurityContextConstraints"]
    class owned_46 owned
    controller -->|"Owns"| owned_47["Service"]
    class owned_47 owned
    controller -->|"Owns"| owned_48["ServiceAccount"]
    class owned_48 owned
    controller -->|"Owns"| owned_49["ServiceMonitor"]
    class owned_49 owned
    controller -->|"Owns"| owned_50["SparkOperator"]
    class owned_50 owned
    controller -->|"Owns"| owned_51["Template"]
    class owned_51 owned
    controller -->|"Owns"| owned_52["Trainer"]
    class owned_52 owned
    controller -->|"Owns"| owned_53["TrainingOperator"]
    class owned_53 owned
    controller -->|"Owns"| owned_54["TrustyAI"]
    class owned_54 owned
    controller -->|"Owns"| owned_55["ValidatingAdmissionPolicy"]
    class owned_55 owned
    controller -->|"Owns"| owned_56["ValidatingAdmissionPolicyBinding"]
    class owned_56 owned
    controller -->|"Owns"| owned_57["ValidatingWebhookConfiguration"]
    class owned_57 owned
    controller -->|"Owns"| owned_58["Workbenches"]
    class owned_58 owned
    watch_59["Auth"] -->|"Watches"| controller
    class watch_59 external
    watch_60["ClusterRole"] -->|"Watches"| controller
    class watch_60 external
    watch_61["HTTPRoute"] -->|"Watches"| controller
    class watch_61 external
    watch_62["Namespace"] -->|"Watches"| controller
    class watch_62 external
    controller -.->|"depends on"| odh_63["models-as-a-service"]
    class odh_63 dep
    controller -.->|"depends on"| odh_64["opendatahub-operator"]
    class odh_64 dep
    controller -.->|"depends on"| odh_65["opendatahub-operator"]
    class odh_65 dep
    controller -.->|"depends on"| odh_66["opendatahub-operator"]
    class odh_66 dep
    controller -.->|"depends on"| odh_67["opendatahub-operator"]
    class odh_67 dep
    controller -.->|"depends on"| odh_68["opendatahub-operator"]
    class odh_68 dep
    controller -.->|"depends on"| odh_69["opendatahub-operator"]
    class odh_69 dep
```

### CRDs

No CRDs defined.

## Dependencies

### Internal Platform Dependencies

| Component | Interaction |
|-----------|-------------|
| models-as-a-service | Go module dependency: github.com/opendatahub-io/models-as-a-service/maas-controller |
| opendatahub-operator | Go module dependency: github.com/opendatahub-io/opendatahub-operator/pkg/clusterhealth |
| opendatahub-operator | Go module dependency: github.com/opendatahub-io/opendatahub-operator/v2/pkg/failureclassifier |
| opendatahub-operator | Go module dependency: github.com/opendatahub-io/opendatahub-operator/pkg/clusterhealth |
| opendatahub-operator | Go module dependency: github.com/opendatahub-io/opendatahub-operator/pkg/clusterhealth |
| opendatahub-operator | Go module dependency: github.com/opendatahub-io/opendatahub-operator/v2/pkg/failureclassifier |
| opendatahub-operator | Go module dependency: github.com/opendatahub-io/opendatahub-operator/pkg/clusterhealth |

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.3 |
| github.com/operator-framework/api | v0.42.0 |
| github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring | v0.68.0 |
| github.com/prometheus/client_golang | v1.23.2 |
| k8s.io/api | v0.35.2 |
| k8s.io/api | v0.35.2 |
| k8s.io/api | v0.35.3 |
| k8s.io/apiextensions-apiserver | v0.35.3 |
| k8s.io/apimachinery | v0.35.3 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/apimachinery | v0.35.2 |
| k8s.io/client-go | v0.35.2 |
| k8s.io/client-go | v0.35.2 |
| k8s.io/client-go | v0.35.3 |
| k8s.io/client-go | v0.35.2 |
| sigs.k8s.io/controller-runtime | v0.23.3 |
| sigs.k8s.io/controller-runtime | v0.22.4 |
| sigs.k8s.io/controller-runtime | v0.22.4 |
| sigs.k8s.io/controller-runtime | v0.22.4 |

