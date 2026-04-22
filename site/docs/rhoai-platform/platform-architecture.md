# Platform Architecture

## CRD Ownership Map

The platform defines 25 CRDs. Each CRD is owned by the component that declares it.

| Owner | CRDs | Count |
|-------|------|-------|
| **codeflare-operator** | AppWrapper | 1 |
| **data-science-pipelines-operator** | DataSciencePipelinesApplication, Pipeline, PipelineVersion, ScheduledWorkflow | 4 |
| **llama-stack-k8s-operator** | LlamaStackDistribution | 1 |
| **model-registry-operator** | ModelRegistry | 1 |
| **modelmesh-serving** | ClusterServingRuntime, InferenceService, Predictor, ServingRuntime | 4 |
| **odh-model-controller** | Account | 1 |
| **opendatahub-operator** | DSCInitialization, DataScienceCluster, OdhApplication, OdhDashboardConfig, OdhDocument, OdhQuickStart | 6 |
| **training-operator** | JAXJob, MPIJob, PaddleJob, PyTorchJob, TFJob, XGBoostJob | 6 |

## Cross-Component Dependencies

Relationships detected through Go module imports and CRD watch patterns.

| From | To | Type |
|------|----|------|
| codeflare-operator | opendatahub-operator | go-module |
| kubeflow | data-science-pipelines-operator | go-module |
| models-as-a-service | kserve | go-module |
| odh-dashboard | llama-stack-k8s-operator | go-module |
| odh-dashboard | mlflow-go | go-module |
| odh-model-controller | kserve | go-module |
| odh-model-controller | modelmesh-serving | watches-crd:ServingRuntime |
| odh-model-controller | modelmesh-serving | watches-crd:InferenceService |
| opendatahub-operator | opendatahub-operator | go-module |
| kube-auth-proxy | kube-rbac-proxy | uses-image |
| kubeflow | kube-rbac-proxy | uses-image |
| llama-stack-k8s-operator | kube-rbac-proxy | uses-image |
| modelmesh-serving | kube-rbac-proxy | uses-image |
| odh-dashboard | kube-rbac-proxy | uses-image |

**Tightest coupling:** `odh-model-controller -> modelmesh-serving` (2 dependency edges).

