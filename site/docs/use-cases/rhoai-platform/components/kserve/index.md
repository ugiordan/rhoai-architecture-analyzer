# kserve

> **Architecture snapshot: 2026-05-19** (2026-05-19)


**Repository:** kserve/kserve  
**Analyzer:** arch-analyzer 0.2.0  
**Extracted:** 2026-05-19T04:09:30Z

## Summary

| Metric | Count |
|--------|-------|
| CRDs | 26 |
| Deployments | 9 |
| Services | 18 |
| Secrets | 10 |
| Cluster Roles | 2 |
| Controller Watches | 152 |

## Component Architecture

CRDs, controllers, and owned Kubernetes resources.

```mermaid
graph LR
    %% Component architecture for kserve

    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef controller fill:#3498db,stroke:#2980b9,color:#fff
    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff

    subgraph controller["kserve Controller"]
        dep_1["keda-metrics-apiserver"]
        class dep_1 controller
        dep_2["keda-metrics-apiserver"]
        class dep_2 controller
        dep_3["keda-operator"]
        class dep_3 controller
        dep_4["keda-operator"]
        class dep_4 controller
        dep_5["kserve-controller-manager"]
        class dep_5 controller
        dep_6["kserve-localmodel-controller-manager"]
        class dep_6 controller
        dep_7["llama-deployment"]
        class dep_7 controller
        dep_8["llama-deployment"]
        class dep_8 controller
        dep_9["llmisvc-controller-manager"]
        class dep_9 controller
    end

    crd_ClusterServingRuntime{{"ClusterServingRuntime\n/v1alpha1"}}
    class crd_ClusterServingRuntime crd
    crd_ClusterStorageContainer{{"ClusterStorageContainer\n/v1alpha1"}}
    class crd_ClusterStorageContainer crd
    crd_InferenceGraph{{"InferenceGraph\n/v1alpha1"}}
    class crd_InferenceGraph crd
    crd_InferenceGraph -->|"For (reconciles)"| controller
    crd_LLMInferenceService{{"LLMInferenceService\n/v1alpha1"}}
    class crd_LLMInferenceService crd
    crd_LLMInferenceService -->|"For (reconciles)"| controller
    crd_LLMInferenceServiceConfig{{"LLMInferenceServiceConfig\n/v1alpha1"}}
    class crd_LLMInferenceServiceConfig crd
    crd_LocalModelCache{{"LocalModelCache\n/v1alpha1"}}
    class crd_LocalModelCache crd
    crd_LocalModelCache -->|"For (reconciles)"| controller
    crd_LocalModelNamespaceCache{{"LocalModelNamespaceCache\n/v1alpha1"}}
    class crd_LocalModelNamespaceCache crd
    crd_LocalModelNamespaceCache -->|"For (reconciles)"| controller
    crd_LocalModelNode{{"LocalModelNode\n/v1alpha1"}}
    class crd_LocalModelNode crd
    crd_LocalModelNode -->|"For (reconciles)"| controller
    crd_LocalModelNodeGroup{{"LocalModelNodeGroup\n/v1alpha1"}}
    class crd_LocalModelNodeGroup crd
    crd_ServingRuntime{{"ServingRuntime\n/v1alpha1"}}
    class crd_ServingRuntime crd
    crd_TrainedModel{{"TrainedModel\n/v1alpha1"}}
    class crd_TrainedModel crd
    crd_TrainedModel -->|"For (reconciles)"| controller
    crd_LLMInferenceService{{"LLMInferenceService\n/v1alpha2"}}
    class crd_LLMInferenceService crd
    crd_LLMInferenceService -->|"For (reconciles)"| controller
    crd_LLMInferenceServiceConfig{{"LLMInferenceServiceConfig\n/v1alpha2"}}
    class crd_LLMInferenceServiceConfig crd
    crd_InferenceService{{"InferenceService\n/v1beta1"}}
    class crd_InferenceService crd
    crd_InferenceService -->|"For (reconciles)"| controller
    crd_ClusterServingRuntime{{"ClusterServingRuntime\nserving.kserve.io/v1alpha1"}}
    class crd_ClusterServingRuntime crd
    crd_ClusterStorageContainer{{"ClusterStorageContainer\nserving.kserve.io/v1alpha1"}}
    class crd_ClusterStorageContainer crd
    crd_InferenceGraph{{"InferenceGraph\nserving.kserve.io/v1alpha1"}}
    class crd_InferenceGraph crd
    crd_InferenceGraph -->|"For (reconciles)"| controller
    crd_LocalModelCache{{"LocalModelCache\nserving.kserve.io/v1alpha1"}}
    class crd_LocalModelCache crd
    crd_LocalModelCache -->|"For (reconciles)"| controller
    crd_LocalModelNamespaceCache{{"LocalModelNamespaceCache\nserving.kserve.io/v1alpha1"}}
    class crd_LocalModelNamespaceCache crd
    crd_LocalModelNamespaceCache -->|"For (reconciles)"| controller
    crd_LocalModelNode{{"LocalModelNode\nserving.kserve.io/v1alpha1"}}
    class crd_LocalModelNode crd
    crd_LocalModelNode -->|"For (reconciles)"| controller
    crd_LocalModelNodeGroup{{"LocalModelNodeGroup\nserving.kserve.io/v1alpha1"}}
    class crd_LocalModelNodeGroup crd
    crd_ServingRuntime{{"ServingRuntime\nserving.kserve.io/v1alpha1"}}
    class crd_ServingRuntime crd
    crd_TrainedModel{{"TrainedModel\nserving.kserve.io/v1alpha1"}}
    class crd_TrainedModel crd
    crd_TrainedModel -->|"For (reconciles)"| controller
    crd_LLMInferenceService{{"LLMInferenceService\nserving.kserve.io/v1alpha2"}}
    class crd_LLMInferenceService crd
    crd_LLMInferenceService -->|"For (reconciles)"| controller
    crd_LLMInferenceServiceConfig{{"LLMInferenceServiceConfig\nserving.kserve.io/v1alpha2"}}
    class crd_LLMInferenceServiceConfig crd
    crd_InferenceService{{"InferenceService\nserving.kserve.io/v1beta1"}}
    class crd_InferenceService crd
    crd_InferenceService -->|"For (reconciles)"| controller
    controller -->|"Owns"| owned_10["ClusterRole"]
    class owned_10 owned
    controller -->|"Owns"| owned_11["ClusterRoleBinding"]
    class owned_11 owned
    controller -->|"Owns"| owned_12["ConfigMap"]
    class owned_12 owned
    controller -->|"Owns"| owned_13["DaemonSet"]
    class owned_13 owned
    controller -->|"Owns"| owned_14["Deployment"]
    class owned_14 owned
    controller -->|"Owns"| owned_15["HTTPRoute"]
    class owned_15 owned
    controller -->|"Owns"| owned_16["HorizontalPodAutoscaler"]
    class owned_16 owned
    controller -->|"Owns"| owned_17["InferencePool"]
    class owned_17 owned
    controller -->|"Owns"| owned_18["Ingress"]
    class owned_18 owned
    controller -->|"Owns"| owned_19["Job"]
    class owned_19 owned
    controller -->|"Owns"| owned_20["LeaderWorkerSet"]
    class owned_20 owned
    controller -->|"Owns"| owned_21["OpenTelemetryCollector"]
    class owned_21 owned
    controller -->|"Owns"| owned_22["PersistentVolume"]
    class owned_22 owned
    controller -->|"Owns"| owned_23["PersistentVolumeClaim"]
    class owned_23 owned
    controller -->|"Owns"| owned_24["PodDisruptionBudget"]
    class owned_24 owned
    controller -->|"Owns"| owned_25["PodMonitor"]
    class owned_25 owned
    controller -->|"Owns"| owned_26["Route"]
    class owned_26 owned
    controller -->|"Owns"| owned_27["ScaledObject"]
    class owned_27 owned
    controller -->|"Owns"| owned_28["Secret"]
    class owned_28 owned
    controller -->|"Owns"| owned_29["Service"]
    class owned_29 owned
    controller -->|"Owns"| owned_30["ServiceAccount"]
    class owned_30 owned
    controller -->|"Owns"| owned_31["ServiceMonitor"]
    class owned_31 owned
    controller -->|"Owns"| owned_32["StatefulSet"]
    class owned_32 owned
    controller -->|"Owns"| owned_33["VariantAutoscaling"]
    class owned_33 owned
    controller -->|"Owns"| owned_34["VirtualService"]
    class owned_34 owned
    watch_35["ClusterServingRuntime"] -->|"Watches"| controller
    class watch_35 external
    watch_36["ConfigMap"] -->|"Watches"| controller
    class watch_36 external
    watch_37["Gateway"] -->|"Watches"| controller
    class watch_37 external
    watch_38["HTTPRoute"] -->|"Watches"| controller
    class watch_38 external
    watch_39["InferenceService"] -->|"Watches"| controller
    class watch_39 external
    watch_40["LLMInferenceServiceConfig"] -->|"Watches"| controller
    class watch_40 external
    watch_41["LocalModelNode"] -->|"Watches"| controller
    class watch_41 external
    watch_42["Node"] -->|"Watches"| controller
    class watch_42 external
    watch_43["Pod"] -->|"Watches"| controller
    class watch_43 external
    watch_44["ServingRuntime"] -->|"Watches"| controller
    class watch_44 external
    watch_45["StatefulSet"] -->|"Watches"| controller
    class watch_45 external
```

### CRDs

| Group | Version | Kind | Scope | Fields | Validation Rules | Discovery | Source |
|-------|---------|------|-------|--------|------------------|-----------|--------|
|  | v1alpha1 | ClusterServingRuntime | Namespaced | 19 | 0 | Go AST | [`/home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1alpha1/servingruntime_types.go`](https://github.com/kserve/kserve/blob/c053aa6f71aafc2b91d87553f2df3cad02b5f0d3//home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1alpha1/servingruntime_types.go) |
|  | v1alpha1 | ClusterStorageContainer | Namespaced | 19 | 0 | Go AST | [`/home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1alpha1/storage_container_types.go`](https://github.com/kserve/kserve/blob/c053aa6f71aafc2b91d87553f2df3cad02b5f0d3//home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1alpha1/storage_container_types.go) |
|  | v1alpha1 | InferenceGraph | Namespaced | 19 | 0 | Go AST | [`/home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1alpha1/inference_graph.go`](https://github.com/kserve/kserve/blob/c053aa6f71aafc2b91d87553f2df3cad02b5f0d3//home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1alpha1/inference_graph.go) |
|  | v1alpha1 | LLMInferenceService | Namespaced | 19 | 0 | Go AST | [`/home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1alpha1/llm_inference_service_types.go`](https://github.com/kserve/kserve/blob/c053aa6f71aafc2b91d87553f2df3cad02b5f0d3//home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1alpha1/llm_inference_service_types.go) |
|  | v1alpha1 | LLMInferenceServiceConfig | Namespaced | 18 | 0 | Go AST | [`/home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1alpha1/llm_inference_service_types.go`](https://github.com/kserve/kserve/blob/c053aa6f71aafc2b91d87553f2df3cad02b5f0d3//home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1alpha1/llm_inference_service_types.go) |
|  | v1alpha1 | LocalModelCache | Namespaced | 19 | 0 | Go AST | [`/home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1alpha1/local_model_cache_types.go`](https://github.com/kserve/kserve/blob/c053aa6f71aafc2b91d87553f2df3cad02b5f0d3//home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1alpha1/local_model_cache_types.go) |
|  | v1alpha1 | LocalModelNamespaceCache | Namespaced | 19 | 0 | Go AST | [`/home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1alpha1/local_model_namespace_cache_types.go`](https://github.com/kserve/kserve/blob/c053aa6f71aafc2b91d87553f2df3cad02b5f0d3//home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1alpha1/local_model_namespace_cache_types.go) |
|  | v1alpha1 | LocalModelNode | Namespaced | 19 | 0 | Go AST | [`/home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1alpha1/local_model_node_types.go`](https://github.com/kserve/kserve/blob/c053aa6f71aafc2b91d87553f2df3cad02b5f0d3//home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1alpha1/local_model_node_types.go) |
|  | v1alpha1 | LocalModelNodeGroup | Namespaced | 19 | 0 | Go AST | [`/home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1alpha1/local_model_node_group_types.go`](https://github.com/kserve/kserve/blob/c053aa6f71aafc2b91d87553f2df3cad02b5f0d3//home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1alpha1/local_model_node_group_types.go) |
|  | v1alpha1 | ServingRuntime | Namespaced | 19 | 0 | Go AST | [`/home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1alpha1/servingruntime_types.go`](https://github.com/kserve/kserve/blob/c053aa6f71aafc2b91d87553f2df3cad02b5f0d3//home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1alpha1/servingruntime_types.go) |
|  | v1alpha1 | TrainedModel | Namespaced | 19 | 0 | Go AST | [`/home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1alpha1/trained_model.go`](https://github.com/kserve/kserve/blob/c053aa6f71aafc2b91d87553f2df3cad02b5f0d3//home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1alpha1/trained_model.go) |
|  | v1alpha2 | LLMInferenceService | Namespaced | 19 | 0 | Go AST | [`/home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1alpha2/llm_inference_service_types.go`](https://github.com/kserve/kserve/blob/c053aa6f71aafc2b91d87553f2df3cad02b5f0d3//home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1alpha2/llm_inference_service_types.go) |
|  | v1alpha2 | LLMInferenceServiceConfig | Namespaced | 18 | 0 | Go AST | [`/home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1alpha2/llm_inference_service_types.go`](https://github.com/kserve/kserve/blob/c053aa6f71aafc2b91d87553f2df3cad02b5f0d3//home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1alpha2/llm_inference_service_types.go) |
|  | v1beta1 | InferenceService | Namespaced | 19 | 0 | Go AST | [`/home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1beta1/inference_service.go`](https://github.com/kserve/kserve/blob/c053aa6f71aafc2b91d87553f2df3cad02b5f0d3//home/runner/work/_temp/arch-analyzer-repos/kserve/pkg/apis/serving/v1beta1/inference_service.go) |
| serving.kserve.io | v1alpha1 | ClusterServingRuntime | Cluster | 1183 | 0 | YAML | [`config/crd/full/serving.kserve.io_clusterservingruntimes.yaml`](https://github.com/kserve/kserve/blob/c053aa6f71aafc2b91d87553f2df3cad02b5f0d3/config/crd/full/serving.kserve.io_clusterservingruntimes.yaml) |
| serving.kserve.io | v1alpha1 | ClusterStorageContainer | Cluster | 216 | 0 | YAML | [`config/crd/full/clusterstoragecontainer/serving.kserve.io_clusterstoragecontainers.yaml`](https://github.com/kserve/kserve/blob/c053aa6f71aafc2b91d87553f2df3cad02b5f0d3/config/crd/full/clusterstoragecontainer/serving.kserve.io_clusterstoragecontainers.yaml) |
| serving.kserve.io | v1alpha1 | InferenceGraph | Namespaced | 150 | 0 | YAML | [`config/crd/full/serving.kserve.io_inferencegraphs.yaml`](https://github.com/kserve/kserve/blob/c053aa6f71aafc2b91d87553f2df3cad02b5f0d3/config/crd/full/serving.kserve.io_inferencegraphs.yaml) |
| serving.kserve.io | v1alpha1 | LocalModelCache | Cluster | 20 | 1 | YAML | [`config/crd/full/localmodel/serving.kserve.io_localmodelcaches.yaml`](https://github.com/kserve/kserve/blob/c053aa6f71aafc2b91d87553f2df3cad02b5f0d3/config/crd/full/localmodel/serving.kserve.io_localmodelcaches.yaml) |
| serving.kserve.io | v1alpha1 | LocalModelNamespaceCache | Namespaced | 20 | 1 | YAML | [`config/crd/full/localmodel/serving.kserve.io_localmodelnamespacecaches.yaml`](https://github.com/kserve/kserve/blob/c053aa6f71aafc2b91d87553f2df3cad02b5f0d3/config/crd/full/localmodel/serving.kserve.io_localmodelnamespacecaches.yaml) |
| serving.kserve.io | v1alpha1 | LocalModelNode | Cluster | 15 | 0 | YAML | [`config/crd/full/localmodel/serving.kserve.io_localmodelnodes.yaml`](https://github.com/kserve/kserve/blob/c053aa6f71aafc2b91d87553f2df3cad02b5f0d3/config/crd/full/localmodel/serving.kserve.io_localmodelnodes.yaml) |
| serving.kserve.io | v1alpha1 | LocalModelNodeGroup | Cluster | 220 | 0 | YAML | [`config/crd/full/localmodel/serving.kserve.io_localmodelnodegroups.yaml`](https://github.com/kserve/kserve/blob/c053aa6f71aafc2b91d87553f2df3cad02b5f0d3/config/crd/full/localmodel/serving.kserve.io_localmodelnodegroups.yaml) |
| serving.kserve.io | v1alpha1 | ServingRuntime | Namespaced | 1183 | 0 | YAML | [`config/crd/full/serving.kserve.io_servingruntimes.yaml`](https://github.com/kserve/kserve/blob/c053aa6f71aafc2b91d87553f2df3cad02b5f0d3/config/crd/full/serving.kserve.io_servingruntimes.yaml) |
| serving.kserve.io | v1alpha1 | TrainedModel | Namespaced | 25 | 0 | YAML | [`config/crd/full/serving.kserve.io_trainedmodels.yaml`](https://github.com/kserve/kserve/blob/c053aa6f71aafc2b91d87553f2df3cad02b5f0d3/config/crd/full/serving.kserve.io_trainedmodels.yaml) |
| serving.kserve.io | v1alpha2 | LLMInferenceService | Namespaced | 5732 | 110 | YAML | [`config/crd/full/llmisvc/serving.kserve.io_llminferenceservices.yaml`](https://github.com/kserve/kserve/blob/c053aa6f71aafc2b91d87553f2df3cad02b5f0d3/config/crd/full/llmisvc/serving.kserve.io_llminferenceservices.yaml) |
| serving.kserve.io | v1alpha2 | LLMInferenceServiceConfig | Namespaced | 5711 | 95 | YAML | [`config/crd/full/llmisvc/serving.kserve.io_llminferenceserviceconfigs.yaml`](https://github.com/kserve/kserve/blob/c053aa6f71aafc2b91d87553f2df3cad02b5f0d3/config/crd/full/llmisvc/serving.kserve.io_llminferenceserviceconfigs.yaml) |
| serving.kserve.io | v1beta1 | InferenceService | Namespaced | 6547 | 0 | YAML | [`config/crd/full/serving.kserve.io_inferenceservices.yaml`](https://github.com/kserve/kserve/blob/c053aa6f71aafc2b91d87553f2df3cad02b5f0d3/config/crd/full/serving.kserve.io_inferenceservices.yaml) |

## Dependencies

### Key External Dependencies

| Module | Version |
|--------|---------|
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.2 |
| github.com/go-logr/logr | v1.4.2 |
| github.com/go-logr/logr | v1.4.2 |
| github.com/go-logr/logr | v1.2.2 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.2 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.3.0 |
| github.com/go-logr/logr | v1.3.0 |
| github.com/go-logr/logr | v1.4.2 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.1 |
| github.com/go-logr/logr | v1.2.2 |
| github.com/go-logr/logr | v1.4.2 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.1 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.2 |
| github.com/go-logr/logr | v1.4.2 |
| github.com/go-logr/logr | v1.4.2 |
| github.com/go-logr/logr | v1.4.2 |
| github.com/go-logr/logr | v1.4.3 |
| github.com/go-logr/logr | v1.4.2 |
| github.com/go-logr/logr | v1.4.2 |
| github.com/go-logr/stdr | v1.2.2 |
| github.com/go-logr/stdr | v1.2.2 |
| github.com/go-logr/stdr | v1.2.2 |
| github.com/go-logr/stdr | v1.2.2 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/go-logr/zapr | v1.3.0 |
| github.com/operator-framework/api | v0.27.0 |
| github.com/operator-framework/api | v0.27.0 |
| github.com/operator-framework/operator-lib | v0.15.0 |
| github.com/operator-framework/operator-lib | v0.15.0 |
| github.com/prometheus-operator/prometheus-operator | v0.76.2 |
| github.com/prometheus-operator/prometheus-operator | v0.76.2 |
| github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring | v0.76.2 |
| github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring | v0.89.0 |
| github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring | v0.76.2 |
| github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring | v0.89.0 |
| github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring | v0.89.0 |
| github.com/prometheus-operator/prometheus-operator/pkg/client | v0.76.2 |
| github.com/prometheus-operator/prometheus-operator/pkg/client | v0.76.2 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_golang | v1.11.1 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_golang | v1.19.1 |
| github.com/prometheus/client_golang | v1.19.1 |
| github.com/prometheus/client_golang | v1.19.1 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_golang | v1.11.1 |
| github.com/prometheus/client_golang | v1.22.0 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_golang | v1.20.5 |
| github.com/prometheus/client_golang | v1.23.2 |
| github.com/prometheus/client_golang | v1.22.0 |
| github.com/prometheus/client_golang | v1.19.1 |
| github.com/prometheus/client_golang | v1.20.5 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.1 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.1 |
| github.com/prometheus/client_model | v0.6.1 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.1 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.1 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/client_model | v0.6.1 |
| github.com/prometheus/client_model | v0.6.2 |
| github.com/prometheus/common | v0.65.0 |
| github.com/prometheus/common | v0.65.0 |
| github.com/prometheus/common | v0.67.4 |
| github.com/prometheus/common | v0.67.5 |
| github.com/prometheus/common | v0.55.0 |
| github.com/prometheus/common | v0.62.0 |
| github.com/prometheus/common | v0.55.0 |
| github.com/prometheus/common | v0.67.5 |
| github.com/prometheus/common | v0.67.5 |
| github.com/prometheus/common | v0.66.1 |
| github.com/prometheus/common | v0.60.1 |
| github.com/prometheus/common | v0.67.5 |
| github.com/prometheus/common | v0.67.4 |
| github.com/prometheus/common | v0.60.1 |
| github.com/prometheus/common | v0.67.4 |
| github.com/prometheus/common | v0.62.0 |
| github.com/prometheus/common | v0.66.1 |
| github.com/prometheus/otlptranslator | v1.0.0 |
| github.com/prometheus/otlptranslator | v1.0.0 |
| github.com/prometheus/otlptranslator | v1.0.0 |
| github.com/prometheus/otlptranslator | v1.0.0 |
| github.com/prometheus/procfs | v0.16.1 |
| github.com/prometheus/procfs | v0.15.1 |
| github.com/prometheus/procfs | v0.15.1 |
| github.com/prometheus/procfs | v0.16.1 |
| github.com/prometheus/prometheus | v0.308.1 |
| github.com/prometheus/prometheus | v0.308.1 |
| github.com/prometheus/prometheus | v0.54.0 |
| github.com/prometheus/prometheus | v0.55.0 |
| github.com/prometheus/prometheus | v0.55.0 |
| github.com/prometheus/prometheus | v0.54.0 |
| google.golang.org/grpc | v1.77.0 |
| google.golang.org/grpc | v1.74.2 |
| google.golang.org/grpc | v1.75.1 |
| google.golang.org/grpc | v1.72.1 |
| google.golang.org/grpc | v1.75.1 |
| google.golang.org/grpc | v1.72.1 |
| google.golang.org/grpc | v1.71.0 |
| google.golang.org/grpc | v1.74.2 |
| google.golang.org/grpc | v1.71.0 |
| google.golang.org/grpc | v1.72.1 |
| google.golang.org/grpc | v1.75.1 |
| google.golang.org/grpc | v1.72.1 |
| google.golang.org/grpc | v1.73.0 |
| google.golang.org/grpc | v1.71.0 |
| google.golang.org/grpc | v1.71.1 |
| google.golang.org/grpc | v1.72.1 |
| google.golang.org/grpc | v1.58.2 |
| google.golang.org/grpc | v1.77.0 |
| google.golang.org/grpc | v1.69.4 |
| google.golang.org/grpc | v1.73.0 |
| google.golang.org/grpc | v1.75.1 |
| google.golang.org/grpc | v1.71.1 |
| google.golang.org/grpc | v1.72.1 |
| google.golang.org/grpc | v1.75.1 |
| google.golang.org/grpc | v1.77.0 |
| google.golang.org/grpc | v1.75.1 |
| google.golang.org/grpc | v1.75.1 |
| google.golang.org/grpc | v1.75.1 |
| google.golang.org/grpc | v1.78.0 |
| google.golang.org/grpc | v1.71.0 |
| google.golang.org/grpc | v1.69.4 |
| google.golang.org/grpc | v1.75.1 |
| google.golang.org/grpc | v1.56.3 |
| google.golang.org/grpc | v1.75.0 |
| google.golang.org/grpc | v1.67.3 |
| google.golang.org/grpc | v1.78.0 |
| google.golang.org/grpc | v1.71.1 |
| google.golang.org/grpc | v1.71.1 |
| google.golang.org/grpc | v1.63.2 |
| google.golang.org/grpc | v1.77.0 |
| google.golang.org/grpc | v1.67.3 |
| google.golang.org/grpc | v1.75.0 |
| google.golang.org/grpc | v1.67.0 |
| google.golang.org/grpc | v1.71.0 |
| google.golang.org/grpc | v1.71.0 |
| google.golang.org/grpc | v1.58.2 |
| google.golang.org/grpc | v1.71.0 |
| google.golang.org/grpc | v1.72.1 |
| google.golang.org/grpc | v1.71.1 |
| google.golang.org/grpc | v1.75.1 |
| google.golang.org/grpc | v1.72.1 |
| google.golang.org/grpc | v1.71.1 |
| google.golang.org/grpc | v1.56.3 |
| google.golang.org/grpc | v1.67.0 |
| google.golang.org/grpc | v1.77.0 |
| google.golang.org/grpc | v1.77.0 |
| google.golang.org/grpc | v1.71.0 |
| google.golang.org/grpc | v1.63.2 |
| google.golang.org/grpc/cmd/protoc-gen-go-grpc | v1.5.1 |
| google.golang.org/grpc/cmd/protoc-gen-go-grpc | v1.5.1 |
| google.golang.org/grpc/cmd/protoc-gen-go-grpc | v1.5.1 |
| google.golang.org/grpc/cmd/protoc-gen-go-grpc | v1.5.1 |
| google.golang.org/grpc/examples | v0.0.0-20250407062114-b368379ef8f6 |
| google.golang.org/grpc/examples | v0.0.0-20250407062114-b368379ef8f6 |
| k8s.io/api | v0.34.1 |
| k8s.io/api | v0.34.3 |
| k8s.io/api | v0.34.3 |
| k8s.io/api | v0.34.5 |
| k8s.io/api | v0.34.3 |
| k8s.io/api | v0.34.3 |
| k8s.io/api | v0.34.3 |
| k8s.io/api | v0.34.3 |
| k8s.io/api | v0.31.2 |
| k8s.io/api | v0.34.1 |
| k8s.io/api | v0.34.3 |
| k8s.io/api | v0.34.3 |
| k8s.io/api | v0.34.3 |
| k8s.io/api | v0.29.0 |
| k8s.io/api | v0.31.0 |
| k8s.io/api | v0.34.5 |
| k8s.io/api | v0.34.5 |
| k8s.io/api | v0.34.3 |
| k8s.io/api | v0.34.3 |
| k8s.io/api | v0.31.7 |
| k8s.io/api | v0.34.5 |
| k8s.io/api | v0.34.3 |
| k8s.io/api | v0.34.3 |
| k8s.io/api | v0.34.3 |
| k8s.io/api | v0.31.7 |
| k8s.io/api | v0.34.3 |
| k8s.io/api | v0.34.3 |
| k8s.io/api | v0.34.5 |
| k8s.io/api | v0.29.0 |
| k8s.io/api | v0.34.3 |
| k8s.io/api | v0.34.3 |
| k8s.io/api | v0.31.2 |
| k8s.io/api | v0.31.0 |
| k8s.io/apiextensions-apiserver | v0.31.0 |
| k8s.io/apiextensions-apiserver | v0.34.3 |
| k8s.io/apiextensions-apiserver | v0.34.3 |
| k8s.io/apiextensions-apiserver | v0.31.2 |
| k8s.io/apiextensions-apiserver | v0.31.0 |
| k8s.io/apiextensions-apiserver | v0.34.3 |
| k8s.io/apiextensions-apiserver | v0.34.3 |
| k8s.io/apiextensions-apiserver | v0.34.3 |
| k8s.io/apiextensions-apiserver | v0.34.3 |
| k8s.io/apiextensions-apiserver | v0.31.2 |
| k8s.io/apiextensions-apiserver | v0.34.1 |
| k8s.io/apiextensions-apiserver | v0.34.3 |
| k8s.io/apiextensions-apiserver | v0.34.3 |
| k8s.io/apiextensions-apiserver | v0.34.3 |
| k8s.io/apiextensions-apiserver | v0.34.3 |
| k8s.io/apiextensions-apiserver | v0.34.1 |
| k8s.io/apiextensions-apiserver | v0.34.3 |
| k8s.io/apimachinery | v0.29.0 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apimachinery | v0.34.5 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apimachinery | v0.32.1 |
| k8s.io/apimachinery | v0.31.0 |
| k8s.io/apimachinery | v0.31.7 |
| k8s.io/apimachinery | v0.34.5 |
| k8s.io/apimachinery | v0.34.5 |
| k8s.io/apimachinery | v0.32.1 |
| k8s.io/apimachinery | v0.34.1 |
| k8s.io/apimachinery | v0.34.5 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apimachinery | v0.31.2 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apimachinery | v0.31.7 |
| k8s.io/apimachinery | v0.31.0 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apimachinery | v0.34.5 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apimachinery | v0.34.1 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apimachinery | v0.34.5 |
| k8s.io/apimachinery | v0.34.5 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apimachinery | v0.31.2 |
| k8s.io/apimachinery | v0.29.0 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apimachinery | v0.34.3 |
| k8s.io/apiserver | v0.34.3 |
| k8s.io/apiserver | v0.31.0 |
| k8s.io/apiserver | v0.34.3 |
| k8s.io/apiserver | v0.34.3 |
| k8s.io/apiserver | v0.31.7 |
| k8s.io/apiserver | v0.31.7 |
| k8s.io/apiserver | v0.31.0 |
| k8s.io/apiserver | v0.34.3 |
| k8s.io/client-go | v0.34.3 |
| k8s.io/client-go | v0.34.5 |
| k8s.io/client-go | v0.32.1 |
| k8s.io/client-go | v0.34.3 |
| k8s.io/client-go | v0.34.3 |
| k8s.io/client-go | v0.34.5 |
| k8s.io/client-go | v0.34.3 |
| k8s.io/client-go | v0.34.3 |
| k8s.io/client-go | v0.31.0 |
| k8s.io/client-go | v0.34.3 |
| k8s.io/client-go | v0.34.3 |
| k8s.io/client-go | v0.34.3 |
| k8s.io/client-go | v0.34.3 |
| k8s.io/client-go | v0.31.2 |
| k8s.io/client-go | v0.31.0 |
| k8s.io/client-go | v0.34.3 |
| k8s.io/client-go | v0.34.1 |
| k8s.io/client-go | v0.34.3 |
| k8s.io/client-go | v0.34.3 |
| k8s.io/client-go | v0.34.3 |
| k8s.io/client-go | v0.34.3 |
| k8s.io/client-go | v0.34.3 |
| k8s.io/client-go | v0.31.2 |
| k8s.io/client-go | v0.34.3 |
| k8s.io/client-go | v0.32.1 |
| k8s.io/client-go | v0.34.3 |
| k8s.io/client-go | v0.31.7 |
| k8s.io/client-go | v0.34.5 |
| k8s.io/client-go | v0.31.7 |
| k8s.io/client-go | v0.34.3 |
| k8s.io/client-go | v0.34.1 |
| sigs.k8s.io/controller-runtime | v0.19.7 |
| sigs.k8s.io/controller-runtime | v0.22.4 |
| sigs.k8s.io/controller-runtime | v0.22.5 |
| sigs.k8s.io/controller-runtime | v0.22.4 |
| sigs.k8s.io/controller-runtime | v0.22.3 |
| sigs.k8s.io/controller-runtime | v0.22.4 |
| sigs.k8s.io/controller-runtime | v0.22.5 |
| sigs.k8s.io/controller-runtime | v0.22.1 |
| sigs.k8s.io/controller-runtime | v0.19.7 |
| sigs.k8s.io/controller-runtime | v0.22.4 |
| sigs.k8s.io/controller-runtime | v0.19.1 |
| sigs.k8s.io/controller-runtime | v0.22.3 |
| sigs.k8s.io/controller-runtime | v0.22.1 |
| sigs.k8s.io/controller-runtime | v0.19.1 |
| sigs.k8s.io/controller-runtime | v0.19.7 |
| sigs.k8s.io/controller-runtime/tools/setup-envtest | v0.0.0-20240804232438-89b5deec030c |
| sigs.k8s.io/controller-runtime/tools/setup-envtest | v0.0.0-20240804232438-89b5deec030c |

