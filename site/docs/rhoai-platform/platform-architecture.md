# Platform Architecture

## CRD Ownership Map

The platform defines 48 CRDs. Each CRD is owned by the component that declares it.

| Owner | CRDs | Count |
|-------|------|-------|
| **codeflare-operator** | AppWrapper | 1 |
| **data-science-pipelines** | CompositeController, ControllerRevision, DecoratorController | 3 |
| **data-science-pipelines-operator** | DataSciencePipelinesApplication, Pipeline, PipelineVersion, ScheduledWorkflow | 4 |
| **kserve** | ClusterServingRuntime, ClusterStorageContainer, InferenceGraph, InferenceService, LLMInferenceService, LLMInferenceServiceConfig, LocalModelCache, LocalModelNamespaceCache, LocalModelNode, LocalModelNodeGroup, ServingRuntime, TrainedModel | 12 |
| **llama-stack-k8s-operator** | LlamaStackDistribution | 1 |
| **mlflow-operator** | MLflow, MLflowConfig | 2 |
| **model-registry-operator** | ModelRegistry | 1 |
| **modelmesh-serving** | ClusterServingRuntime, InferenceService, Predictor, ServingRuntime | 4 |
| **odh-model-controller** | Account | 1 |
| **opendatahub-operator** | DSCInitialization, DataScienceCluster, OdhApplication, OdhDashboardConfig, OdhDocument, OdhQuickStart | 6 |
| **spark-operator** | ScheduledSparkApplication, SparkApplication, SparkConnect | 3 |
| **trainer** | ClusterTrainingRuntime, TrainJob, TrainingRuntime | 3 |
| **training-operator** | JAXJob, MPIJob, PaddleJob, PyTorchJob, TFJob, XGBoostJob | 6 |
| **workload-variant-autoscaler** | VariantAutoscaling | 1 |

## Cross-Component Dependencies

Relationships detected through Go module imports and CRD watch patterns.

| From | To | Type |
|------|----|------|
| codeflare-operator | opendatahub-operator | go-module |
| kubeflow | data-science-pipelines-operator | go-module |
| mlflow-operator | mlflow-operator | go-module |
| model-registry | kserve | watches-crd:InferenceService |
| modelmesh-serving | kserve | watches-crd:ServingRuntime |
| models-as-a-service | kserve | go-module |
| odh-dashboard | llama-stack-k8s-operator | go-module |
| odh-dashboard | mlflow-go | go-module |
| odh-model-controller | kserve | go-module |
| odh-model-controller | kserve | watches-crd:InferenceGraph |
| odh-model-controller | kserve | watches-crd:ServingRuntime |
| odh-model-controller | kserve | watches-crd:LLMInferenceService |
| odh-model-controller | kserve | watches-crd:InferenceService |
| opendatahub-operator | opendatahub-operator | go-module |

**Tightest coupling:** `odh-model-controller -> kserve` (5 dependency edges).

