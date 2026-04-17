# Platform Architecture

## CRD Ownership Map

The platform defines 50 CRDs. Each CRD is owned by the component that declares it.

| Owner | CRDs | Count |
|-------|------|-------|
| **data-science-pipelines-operator** | DataSciencePipelinesApplication, Pipeline, PipelineVersion, ScheduledWorkflow | 4 |
| **kserve** | ClusterServingRuntime, ClusterStorageContainer, InferenceGraph, InferenceService, LLMInferenceService, LLMInferenceServiceConfig, LocalModelCache, LocalModelNamespaceCache, LocalModelNode, LocalModelNodeGroup, ServingRuntime, TrainedModel | 12 |
| **model-registry-operator** | ModelRegistry | 1 |
| **odh-model-controller** | Account | 1 |
| **opendatahub-operator** | Auth, AzureKubernetesEngine, CoreWeaveKubernetesEngine, Dashboard, DataSciencePipelines, FeastOperator, GatewayConfig, HardwareProfile, Kserve, Kueue, LlamaStackOperator, MLflowOperator, ModelController, ModelRegistry, ModelsAsService, Monitoring, Ray, SparkOperator, Trainer, TrainingOperator, TrustyAI, Workbenches | 22 |
| **trustyai-service-operator** | EvalHub, GuardrailsOrchestrator, LMEvalJob, NemoGuardrails, TrustyAIService | 5 |

## Cross-Component Dependencies

Relationships detected through Go module imports and CRD watch patterns.

| From | To | Type |
|------|----|------|
| odh-dashboard | llama-stack-k8s-operator | go-module |
| odh-dashboard | mlflow-go | go-module |
| odh-model-controller | kserve | go-module |
| odh-model-controller | kserve | watches-crd:InferenceGraph |
| odh-model-controller | kserve | watches-crd:InferenceService |
| odh-model-controller | kserve | watches-crd:LLMInferenceService |
| odh-model-controller | kserve | watches-crd:ServingRuntime |
| opendatahub-operator | models-as-a-service | go-module |
| opendatahub-operator | opendatahub-operator | go-module |
| kserve | kube-rbac-proxy | uses-image |
| kube-auth-proxy | kube-rbac-proxy | uses-image |
| odh-dashboard | kube-rbac-proxy | uses-image |
| opendatahub-operator | kube-rbac-proxy | uses-image |
| opendatahub-operator | kserve | uses-image |

**Tightest coupling:** `odh-model-controller -> kserve` (5 dependency edges).

