# Network Topology

55 Kubernetes services across the platform.

## Network Topology Graph

Interactive service mesh view of the platform. Drag nodes to rearrange, hover to highlight connections, click for details. Double-click background to fit all.

<div class="topology-toolbar">
  <button data-action="fit" title="Fit to view">Fit</button>
  <button data-action="zoom-in" title="Zoom in">+</button>
  <button data-action="zoom-out" title="Zoom out">&minus;</button>
  <button data-action="relayout" title="Re-run layout">Relayout</button>
  <button data-action="fullscreen" title="Toggle fullscreen">Fullscreen</button>
</div>
<div class="cytoscape-topology">
  <script type="application/json">
  {"components":[{"id":"argo_workflows","name":"argo-workflows","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"batch_gateway","name":"batch-gateway","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"codeflare_operator","name":"codeflare-operator","serviceCount":1,"netpolCount":0,"hasIngress":true},{"id":"data_science_pipelines","name":"data-science-pipelines","serviceCount":2,"netpolCount":1,"hasIngress":true},{"id":"data_science_pipelines_operator","name":"data-science-pipelines-operator","serviceCount":10,"netpolCount":2,"hasIngress":true},{"id":"distributed_workloads","name":"distributed-workloads","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"eval_hub","name":"eval-hub","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"fms_guardrails_orchestrator","name":"fms-guardrails-orchestrator","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"kserve","name":"kserve","serviceCount":6,"netpolCount":1,"hasIngress":true},{"id":"kube_auth_proxy","name":"kube-auth-proxy","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"kube_rbac_proxy","name":"kube-rbac-proxy","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"kubeflow","name":"kubeflow","serviceCount":2,"netpolCount":0,"hasIngress":false},{"id":"kuberay","name":"kuberay","serviceCount":2,"netpolCount":0,"hasIngress":false},{"id":"kueue","name":"kueue","serviceCount":2,"netpolCount":0,"hasIngress":false},{"id":"llama_stack_k8s_operator","name":"llama-stack-k8s-operator","serviceCount":1,"netpolCount":1,"hasIngress":true},{"id":"llm_d_inference_scheduler","name":"llm-d-inference-scheduler","serviceCount":4,"netpolCount":0,"hasIngress":true},{"id":"llm_d_kv_cache","name":"llm-d-kv-cache","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"mlflow_operator","name":"mlflow-operator","serviceCount":3,"netpolCount":1,"hasIngress":true},{"id":"model_registry","name":"model-registry","serviceCount":1,"netpolCount":0,"hasIngress":false},{"id":"model_registry_operator","name":"model-registry-operator","serviceCount":4,"netpolCount":4,"hasIngress":true},{"id":"modelmesh_serving","name":"modelmesh-serving","serviceCount":3,"netpolCount":4,"hasIngress":false},{"id":"models_as_a_service","name":"models-as-a-service","serviceCount":2,"netpolCount":4,"hasIngress":true},{"id":"notebooks","name":"notebooks","serviceCount":1,"netpolCount":0,"hasIngress":false},{"id":"odh_dashboard","name":"odh-dashboard","serviceCount":5,"netpolCount":5,"hasIngress":true},{"id":"odh_model_controller","name":"odh-model-controller","serviceCount":2,"netpolCount":0,"hasIngress":true},{"id":"opendatahub_operator","name":"opendatahub-operator","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"spark_operator","name":"spark-operator","serviceCount":1,"netpolCount":1,"hasIngress":true},{"id":"trainer","name":"trainer","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"training_operator","name":"training-operator","serviceCount":1,"netpolCount":0,"hasIngress":false},{"id":"trustyai_service_operator","name":"trustyai-service-operator","serviceCount":2,"netpolCount":0,"hasIngress":false},{"id":"workload_variant_autoscaler","name":"workload-variant-autoscaler","serviceCount":0,"netpolCount":0,"hasIngress":false}],"services":[{"id":"codeflare_operator_svc_0","name":"webhook-service","parent":"codeflare_operator","ports":"443"},{"id":"data_science_pipelines_svc_0","name":"kubeflow-pipelines-profile-controller","parent":"data_science_pipelines","ports":"80"},{"id":"data_science_pipelines_svc_1","name":"squid","parent":"data_science_pipelines","ports":"3128"},{"id":"data_science_pipelines_operator_svc_0","name":"data-science-pipelines-operator-service","parent":"data_science_pipelines_operator","ports":"8080"},{"id":"data_science_pipelines_operator_svc_1","name":"ds-pipeline-workflow-controller-metrics-template-value","parent":"data_science_pipelines_operator","ports":"9090"},{"id":"data_science_pipelines_operator_svc_2","name":"mariadb","parent":"data_science_pipelines_operator","ports":"3306"},{"id":"data_science_pipelines_operator_svc_3","name":"mariadb-template-value","parent":"data_science_pipelines_operator","ports":"3306"},{"id":"data_science_pipelines_operator_svc_4","name":"minio","parent":"data_science_pipelines_operator","ports":"9000,9001"},{"id":"data_science_pipelines_operator_svc_5","name":"minio-service","parent":"data_science_pipelines_operator","ports":"9000"},{"id":"data_science_pipelines_operator_svc_6","name":"minio-template-value","parent":"data_science_pipelines_operator","ports":"9000,80"},{"id":"data_science_pipelines_operator_svc_7","name":"ml-pipeline","parent":"data_science_pipelines_operator","ports":"8443,8888,8887"},{"id":"data_science_pipelines_operator_svc_8","name":"pypi-server","parent":"data_science_pipelines_operator","ports":"8080"},{"id":"data_science_pipelines_operator_svc_9","name":"template-value","parent":"data_science_pipelines_operator","ports":"8443,8888,8887"},{"id":"kserve_svc_0","name":"kserve-controller-manager-metrics-service","parent":"kserve","ports":"8443"},{"id":"kserve_svc_1","name":"kserve-controller-manager-service","parent":"kserve","ports":"8443"},{"id":"kserve_svc_2","name":"kserve-webhook-server-service","parent":"kserve","ports":"443"},{"id":"kserve_svc_3","name":"llmisvc-controller-manager-service","parent":"kserve","ports":"8443"},{"id":"kserve_svc_4","name":"llmisvc-webhook-server-service","parent":"kserve","ports":"443"},{"id":"kserve_svc_5","name":"localmodel-webhook-server-service","parent":"kserve","ports":"443"},{"id":"kubeflow_svc_0","name":"service","parent":"kubeflow","ports":"443"},{"id":"kubeflow_svc_1","name":"webhook-service","parent":"kubeflow","ports":"443"},{"id":"kuberay_svc_0","name":"kuberay-operator","parent":"kuberay","ports":"8080"},{"id":"kuberay_svc_1","name":"webhook-service","parent":"kuberay","ports":"443"},{"id":"kueue_svc_0","name":"visibility-server","parent":"kueue","ports":"443"},{"id":"kueue_svc_1","name":"webhook-service","parent":"kueue","ports":"443"},{"id":"llama_stack_k8s_operator_svc_0","name":"llama-stack-k8s-operator-controller-manager-metrics-service","parent":"llama_stack_k8s_operator","ports":"8443"},{"id":"llm_d_inference_scheduler_svc_0","name":"${EPP_NAME}","parent":"llm_d_inference_scheduler","ports":"9002,5557,9090"},{"id":"llm_d_inference_scheduler_svc_1","name":"inference-gateway-istio-nodeport","parent":"llm_d_inference_scheduler","ports":"15021,80"},{"id":"llm_d_inference_scheduler_svc_2","name":"istiod-llm-d-gateway","parent":"llm_d_inference_scheduler","ports":"15010,15012,443,15014"},{"id":"llm_d_inference_scheduler_svc_3","name":"service","parent":"llm_d_inference_scheduler","ports":"8080"},{"id":"mlflow_operator_svc_0","name":"minio-service","parent":"mlflow_operator","ports":"9000"},{"id":"mlflow_operator_svc_1","name":"mlflow-operator-controller-manager-metrics-service","parent":"mlflow_operator","ports":"8443"},{"id":"mlflow_operator_svc_2","name":"postgres-service","parent":"mlflow_operator","ports":"5432"},{"id":"model_registry_svc_0","name":"model-catalog","parent":"model_registry","ports":"8080"},{"id":"model_registry_operator_svc_0","name":"model-registry-operator-controller-manager-metrics-service","parent":"model_registry_operator","ports":"8443"},{"id":"model_registry_operator_svc_1","name":"model-registry-operator-webhook-service","parent":"model_registry_operator","ports":"443"},{"id":"model_registry_operator_svc_2","name":"template-value","parent":"model_registry_operator","ports":"0,0"},{"id":"model_registry_operator_svc_3","name":"template-value-postgres","parent":"model_registry_operator","ports":"5432"},{"id":"modelmesh_serving_svc_0","name":"etcd","parent":"modelmesh_serving","ports":"2379"},{"id":"modelmesh_serving_svc_1","name":"modelmesh-controller","parent":"modelmesh_serving","ports":"8080"},{"id":"modelmesh_serving_svc_2","name":"modelmesh-webhook-server-service","parent":"modelmesh_serving","ports":"9443"},{"id":"models_as_a_service_svc_0","name":"maas-api","parent":"models_as_a_service","ports":"8080"},{"id":"models_as_a_service_svc_1","name":"payload-processing","parent":"models_as_a_service","ports":"9004"},{"id":"notebooks_svc_0","name":"notebook","parent":"notebooks","ports":"8888"},{"id":"odh_dashboard_svc_0","name":"odh-dashboard","parent":"odh_dashboard","ports":"8443"},{"id":"odh_dashboard_svc_1","name":"workspaces-backend","parent":"odh_dashboard","ports":"4000"},{"id":"odh_dashboard_svc_2","name":"workspaces-controller-metrics-service","parent":"odh_dashboard","ports":"8080"},{"id":"odh_dashboard_svc_3","name":"workspaces-frontend","parent":"odh_dashboard","ports":"8080"},{"id":"odh_dashboard_svc_4","name":"workspaces-webhook-service","parent":"odh_dashboard","ports":"443"},{"id":"odh_model_controller_svc_0","name":"model-serving-api","parent":"odh_model_controller","ports":"443,9090"},{"id":"odh_model_controller_svc_1","name":"odh-model-controller-webhook-service","parent":"odh_model_controller","ports":"443"},{"id":"spark_operator_svc_0","name":"spark-operator-webhook-svc","parent":"spark_operator","ports":"443"},{"id":"training_operator_svc_0","name":"training-operator","parent":"training_operator","ports":"8080,443"},{"id":"trustyai_service_operator_svc_0","name":"trustyai-service-operator-controller-manager-metrics-service","parent":"trustyai_service_operator","ports":"8443"},{"id":"trustyai_service_operator_svc_1","name":"trustyai-service-operator-metrics-service","parent":"trustyai_service_operator","ports":"8080"}],"externals":[{"id":"ext_grpc","name":"grpc","type":"grpc"},{"id":"ext_azure_blob","name":"azure-blob","type":"object-storage"},{"id":"ext_gcs","name":"gcs","type":"object-storage"},{"id":"ext_redis","name":"redis","type":"database"},{"id":"ext_s3","name":"s3","type":"object-storage"},{"id":"ext_minio","name":"minio","type":"object-storage"},{"id":"ext_mysql","name":"mysql","type":"database"}],"edges":[{"from":"codeflare_operator","to":"opendatahub_operator","type":"module"},{"from":"kubeflow","to":"data_science_pipelines_operator","type":"module"},{"from":"model_registry","to":"kserve","type":"watches"},{"from":"modelmesh_serving","to":"kserve","type":"watches"},{"from":"models_as_a_service","to":"kserve","type":"module"},{"from":"odh_dashboard","to":"llama_stack_k8s_operator","type":"module"},{"from":"odh_dashboard","to":"mlflow_go","type":"module"},{"from":"odh_model_controller","to":"kserve","type":"watches"},{"from":"argo_workflows","to":"ext_grpc","type":"external"},{"from":"argo_workflows","to":"ext_azure_blob","type":"external"},{"from":"argo_workflows","to":"ext_gcs","type":"external"},{"from":"batch_gateway","to":"ext_redis","type":"external"},{"from":"batch_gateway","to":"ext_s3","type":"external"},{"from":"data_science_pipelines","to":"ext_grpc","type":"external"},{"from":"data_science_pipelines","to":"ext_minio","type":"external"},{"from":"data_science_pipelines","to":"ext_s3","type":"external"},{"from":"data_science_pipelines_operator","to":"ext_mysql","type":"external"},{"from":"data_science_pipelines_operator","to":"ext_minio","type":"external"},{"from":"eval_hub","to":"ext_s3","type":"external"},{"from":"kserve","to":"ext_azure_blob","type":"external"},{"from":"kserve","to":"ext_gcs","type":"external"},{"from":"kserve","to":"ext_s3","type":"external"},{"from":"kube_auth_proxy","to":"ext_redis","type":"external"},{"from":"kuberay","to":"ext_grpc","type":"external"},{"from":"llm_d_kv_cache","to":"ext_redis","type":"external"},{"from":"modelmesh_serving","to":"ext_azure_blob","type":"external"},{"from":"odh_dashboard","to":"ext_s3","type":"external"}]}
  </script>
</div>
<div class="topology-legend">
  <span><span class="swatch" style="background:#3498db"></span> Component</span>
  <span><span class="swatch" style="background:#27ae60"></span> Has Ingress</span>
  <span><span class="swatch" style="background:#3498db;border:2px solid #f39c12"></span> Has NetworkPolicy</span>
  <span><span class="swatch" style="background:#e74c3c;border-radius:2px;transform:rotate(45deg)"></span> External</span>
  <span><span class="line-swatch" style="background:#e74c3c"></span> CRD Watch</span>
  <span><span class="line-swatch" style="background:#9b59b6"></span> Sidecar</span>
  <span><span class="line-swatch" style="background:#95a5a6;border-top:2px dashed #95a5a6;height:0"></span> Module</span>
  <span><span class="line-swatch" style="background:#e67e22;border-top:2px dotted #e67e22;height:0"></span> External</span>
</div>

## Cross-Component Service References

Services referenced across component boundaries. When component A defines a service that component B also references, it indicates a deployment dependency.

```mermaid
graph LR
    classDef comp fill:#3498db,stroke:#2980b9,color:#fff

    data_science_pipelines_operator["data-science-pipelines-operator"]:::comp
    mlflow_operator["mlflow-operator"]:::comp
    model_registry_operator["model-registry-operator"]:::comp

    model_registry_operator -.->|"template-value"| data_science_pipelines_operator
    mlflow_operator -.->|"minio-service"| data_science_pipelines_operator
```

## Services by Component

| Component | Services | Webhook (443) | Metrics (8443) | Data |
|-----------|----------|---------------|----------------|------|
| codeflare-operator | 1 | 1 | 0 | 0 |
| data-science-pipelines | 2 | 0 | 0 | 2 |
| data-science-pipelines-operator | 10 | 0 | 2 | 8 |
| kserve | 6 | 3 | 3 | 0 |
| kubeflow | 2 | 2 | 0 | 0 |
| kuberay | 2 | 1 | 0 | 1 |
| kueue | 2 | 2 | 0 | 0 |
| llama-stack-k8s-operator | 1 | 0 | 1 | 0 |
| llm-d-inference-scheduler | 4 | 1 | 0 | 3 |
| mlflow-operator | 3 | 0 | 1 | 2 |
| model-registry | 1 | 0 | 0 | 1 |
| model-registry-operator | 4 | 1 | 1 | 2 |
| modelmesh-serving | 3 | 0 | 0 | 3 |
| models-as-a-service | 2 | 0 | 0 | 2 |
| notebooks | 1 | 0 | 0 | 1 |
| odh-dashboard | 5 | 1 | 1 | 3 |
| odh-model-controller | 2 | 2 | 0 | 0 |
| spark-operator | 1 | 1 | 0 | 0 |
| training-operator | 1 | 1 | 0 | 0 |
| trustyai-service-operator | 2 | 0 | 1 | 1 |

## Service Detail

Per-component service breakdown with exact port numbers and protocols.

### codeflare-operator (1 services)

| Service | Type | Ports |
|---------|------|-------|
| webhook-service | ClusterIP | 443/TCP |

### data-science-pipelines (2 services)

| Service | Type | Ports |
|---------|------|-------|
| kubeflow-pipelines-profile-controller | ClusterIP | 80/TCP |
| squid | ClusterIP | 3128/TCP |

### data-science-pipelines-operator (10 services)

| Service | Type | Ports |
|---------|------|-------|
| data-science-pipelines-operator-service | ClusterIP | 8080/TCP |
| ds-pipeline-workflow-controller-metrics-template-value | ClusterIP | 9090/TCP |
| mariadb | ClusterIP | 3306/TCP |
| mariadb-template-value | ClusterIP | 3306/TCP |
| minio | ClusterIP | 9000/TCP, 9001/TCP |
| minio-service | ClusterIP | 9000/TCP |
| minio-template-value | ClusterIP | 9000/TCP, 80/TCP |
| ml-pipeline | ClusterIP | 8443/TCP, 8888/TCP, 8887/TCP |
| pypi-server | ClusterIP | 8080/TCP |
| template-value | ClusterIP | 8443/TCP, 8888/TCP, 8887/TCP |

### kserve (6 services)

| Service | Type | Ports |
|---------|------|-------|
| kserve-controller-manager-metrics-service | ClusterIP | 8443/TCP |
| kserve-controller-manager-service | ClusterIP | 8443/TCP |
| kserve-webhook-server-service | ClusterIP | 443/TCP |
| llmisvc-controller-manager-service | ClusterIP | 8443/TCP |
| llmisvc-webhook-server-service | ClusterIP | 443/TCP |
| localmodel-webhook-server-service | ClusterIP | 443/TCP |

### kubeflow (2 services)

| Service | Type | Ports |
|---------|------|-------|
| service | ClusterIP | 443/TCP |
| webhook-service | ClusterIP | 443/TCP |

### kuberay (2 services)

| Service | Type | Ports |
|---------|------|-------|
| kuberay-operator | ClusterIP | 8080/TCP |
| webhook-service | ClusterIP | 443/TCP |

### kueue (2 services)

| Service | Type | Ports |
|---------|------|-------|
| visibility-server | ClusterIP | 443/TCP |
| webhook-service | ClusterIP | 443/TCP |

### llama-stack-k8s-operator (1 services)

| Service | Type | Ports |
|---------|------|-------|
| llama-stack-k8s-operator-controller-manager-metrics-service | ClusterIP | 8443/TCP |

### llm-d-inference-scheduler (4 services)

| Service | Type | Ports |
|---------|------|-------|
| ${EPP_NAME} | ClusterIP | 9002/TCP, 5557/TCP, 9090/TCP |
| inference-gateway-istio-nodeport | NodePort | 15021/TCP, 80/TCP |
| istiod-llm-d-gateway | ClusterIP | 15010/TCP, 15012/TCP, 443/TCP, 15014/TCP |
| service | ClusterIP | 8080/TCP |

### mlflow-operator (3 services)

| Service | Type | Ports |
|---------|------|-------|
| minio-service | ClusterIP | 9000/TCP |
| mlflow-operator-controller-manager-metrics-service | ClusterIP | 8443/TCP |
| postgres-service | ClusterIP | 5432/TCP |

### model-registry (1 services)

| Service | Type | Ports |
|---------|------|-------|
| model-catalog | ClusterIP | 8080/TCP |

### model-registry-operator (4 services)

| Service | Type | Ports |
|---------|------|-------|
| model-registry-operator-controller-manager-metrics-service | ClusterIP | 8443/TCP |
| model-registry-operator-webhook-service | ClusterIP | 443/TCP |
| template-value | ClusterIP | 0/TCP, 0/TCP |
| template-value-postgres | ClusterIP | 5432/TCP |

### modelmesh-serving (3 services)

| Service | Type | Ports |
|---------|------|-------|
| etcd | ClusterIP | 2379/TCP |
| modelmesh-controller | ClusterIP | 8080/TCP |
| modelmesh-webhook-server-service | ClusterIP | 9443/TCP |

### models-as-a-service (2 services)

| Service | Type | Ports |
|---------|------|-------|
| maas-api | ClusterIP | 8080/TCP |
| payload-processing | ClusterIP | 9004/TCP |

### notebooks (1 services)

| Service | Type | Ports |
|---------|------|-------|
| notebook | ClusterIP | 8888/TCP |

### odh-dashboard (5 services)

| Service | Type | Ports |
|---------|------|-------|
| odh-dashboard | ClusterIP | 8443/TCP |
| workspaces-backend | ClusterIP | 4000/TCP |
| workspaces-controller-metrics-service | ClusterIP | 8080/TCP |
| workspaces-frontend | ClusterIP | 8080/TCP |
| workspaces-webhook-service | ClusterIP | 443/TCP |

### odh-model-controller (2 services)

| Service | Type | Ports |
|---------|------|-------|
| model-serving-api | ClusterIP | 443/TCP, 9090/TCP |
| odh-model-controller-webhook-service | ClusterIP | 443/TCP |

### spark-operator (1 services)

| Service | Type | Ports |
|---------|------|-------|
| spark-operator-webhook-svc | ClusterIP | 443/TCP |

### training-operator (1 services)

| Service | Type | Ports |
|---------|------|-------|
| training-operator | ClusterIP | 8080/TCP, 443/TCP |

### trustyai-service-operator (2 services)

| Service | Type | Ports |
|---------|------|-------|
| trustyai-service-operator-controller-manager-metrics-service | ClusterIP | 8443/TCP |
| trustyai-service-operator-metrics-service | ClusterIP | 8080/TCP |

## Port Patterns

- **0/TCP**: template-value, template-value
- **15010/TCP**: istiod-llm-d-gateway
- **15012/TCP**: istiod-llm-d-gateway
- **15014/TCP**: istiod-llm-d-gateway
- **15021/TCP**: inference-gateway-istio-nodeport
- **2379/TCP**: etcd
- **3128/TCP**: squid
- **3306/TCP**: mariadb, mariadb-template-value
- **4000/TCP**: workspaces-backend
- **443/TCP**: webhook-service, kserve-webhook-server-service, llmisvc-webhook-server-service, localmodel-webhook-server-service, service, webhook-service, webhook-service, visibility-server, webhook-service, istiod-llm-d-gateway, model-registry-operator-webhook-service, workspaces-webhook-service, model-serving-api, odh-model-controller-webhook-service, spark-operator-webhook-svc, training-operator
- **5432/TCP**: postgres-service, template-value-postgres
- **5557/TCP**: ${EPP_NAME}
- **80/TCP**: minio-template-value, kubeflow-pipelines-profile-controller, inference-gateway-istio-nodeport
- **8080/TCP**: data-science-pipelines-operator-service, pypi-server, kuberay-operator, service, model-catalog, modelmesh-controller, maas-api, workspaces-controller-metrics-service, workspaces-frontend, training-operator, trustyai-service-operator-metrics-service
- **8443/TCP**: ml-pipeline, template-value, kserve-controller-manager-metrics-service, kserve-controller-manager-service, llmisvc-controller-manager-service, llama-stack-k8s-operator-controller-manager-metrics-service, mlflow-operator-controller-manager-metrics-service, model-registry-operator-controller-manager-metrics-service, odh-dashboard, trustyai-service-operator-controller-manager-metrics-service
- **8887/TCP**: ml-pipeline, template-value
- **8888/TCP**: ml-pipeline, template-value, notebook
- **9000/TCP**: minio, minio-service, minio-template-value, minio-service
- **9001/TCP**: minio
- **9002/TCP**: ${EPP_NAME}
- **9004/TCP**: payload-processing
- **9090/TCP**: ds-pipeline-workflow-controller-metrics-template-value, ${EPP_NAME}, model-serving-api
- **9443/TCP**: modelmesh-webhook-server-service

