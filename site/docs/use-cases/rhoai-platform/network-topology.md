# Network Topology

81 Kubernetes services across the platform.

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
  {"components":[{"id":"MLServer","name":"MLServer","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"NeMo_Guardrails","name":"NeMo-Guardrails","serviceCount":1,"netpolCount":0,"hasIngress":false},{"id":"ai_gateway_payload_processing","name":"ai-gateway-payload-processing","serviceCount":1,"netpolCount":0,"hasIngress":true},{"id":"ai4rag","name":"ai4rag","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"argo_workflows","name":"argo-workflows","serviceCount":1,"netpolCount":0,"hasIngress":false},{"id":"batch_gateway","name":"batch-gateway","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"caikit_nlp","name":"caikit-nlp","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"caikit_tgis_serving","name":"caikit-tgis-serving","serviceCount":0,"netpolCount":1,"hasIngress":true},{"id":"codeflare_sdk","name":"codeflare-sdk","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"data_science_pipelines","name":"data-science-pipelines","serviceCount":2,"netpolCount":1,"hasIngress":true},{"id":"data_science_pipelines_operator","name":"data-science-pipelines-operator","serviceCount":11,"netpolCount":2,"hasIngress":true},{"id":"distributed_workloads","name":"distributed-workloads","serviceCount":4,"netpolCount":0,"hasIngress":false},{"id":"eval_hub","name":"eval-hub","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"feast","name":"feast","serviceCount":1,"netpolCount":0,"hasIngress":false},{"id":"fms_guardrails_orchestrator","name":"fms-guardrails-orchestrator","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"fms_hf_tuning","name":"fms-hf-tuning","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"gateway_api_inference_extension","name":"gateway-api-inference-extension","serviceCount":1,"netpolCount":0,"hasIngress":true},{"id":"guardrails_regex_detector","name":"guardrails-regex-detector","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"kserve","name":"kserve","serviceCount":12,"netpolCount":1,"hasIngress":true},{"id":"kube_auth_proxy","name":"kube-auth-proxy","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"kube_rbac_proxy","name":"kube-rbac-proxy","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"kubeflow","name":"kubeflow","serviceCount":2,"netpolCount":0,"hasIngress":false},{"id":"kuberay","name":"kuberay","serviceCount":3,"netpolCount":0,"hasIngress":false},{"id":"kueue","name":"kueue","serviceCount":5,"netpolCount":0,"hasIngress":true},{"id":"llama_stack","name":"llama-stack","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"llama_stack_k8s_operator","name":"llama-stack-k8s-operator","serviceCount":2,"netpolCount":1,"hasIngress":true},{"id":"llm_d","name":"llm-d","serviceCount":0,"netpolCount":0,"hasIngress":true},{"id":"llm_d_inference_scheduler","name":"llm-d-inference-scheduler","serviceCount":5,"netpolCount":0,"hasIngress":true},{"id":"llm_d_kv_cache","name":"llm-d-kv-cache","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"llm_d_routing_sidecar","name":"llm-d-routing-sidecar","serviceCount":1,"netpolCount":0,"hasIngress":true},{"id":"lm_evaluation_harness","name":"lm-evaluation-harness","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"ml_metadata","name":"ml-metadata","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"mlflow_operator","name":"mlflow-operator","serviceCount":3,"netpolCount":1,"hasIngress":true},{"id":"model_metadata_collection","name":"model-metadata-collection","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"model_registry","name":"model-registry","serviceCount":1,"netpolCount":0,"hasIngress":false},{"id":"modelmesh_runtime_adapter","name":"modelmesh-runtime-adapter","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"modelmesh_serving","name":"modelmesh-serving","serviceCount":7,"netpolCount":4,"hasIngress":false},{"id":"models_as_a_service","name":"models-as-a-service","serviceCount":2,"netpolCount":6,"hasIngress":true},{"id":"notebooks","name":"notebooks","serviceCount":1,"netpolCount":0,"hasIngress":false},{"id":"notebooks_downstream","name":"notebooks-downstream","serviceCount":1,"netpolCount":0,"hasIngress":false},{"id":"odh_cli","name":"odh-cli","serviceCount":1,"netpolCount":0,"hasIngress":false},{"id":"odh_deployer","name":"odh-deployer","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"openvino_model_server","name":"openvino_model_server","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"pipelines_components","name":"pipelines-components","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"rest_proxy","name":"rest-proxy","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"rhds_llama_stack_distribution","name":"rhds-llama-stack-distribution","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"spark_operator","name":"spark-operator","serviceCount":2,"netpolCount":1,"hasIngress":true},{"id":"text_generation_inference","name":"text-generation-inference","serviceCount":1,"netpolCount":0,"hasIngress":false},{"id":"trainer","name":"trainer","serviceCount":1,"netpolCount":0,"hasIngress":false},{"id":"trustyai_explainability","name":"trustyai-explainability","serviceCount":0,"netpolCount":0,"hasIngress":true},{"id":"vllm_cpu","name":"vllm-cpu","serviceCount":3,"netpolCount":0,"hasIngress":false},{"id":"vllm_gaudi","name":"vllm-gaudi","serviceCount":1,"netpolCount":0,"hasIngress":false},{"id":"vllm_orchestrator_gateway","name":"vllm-orchestrator-gateway","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"vllm_rocm","name":"vllm-rocm","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"vllm_spyre","name":"vllm-spyre","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"workload_variant_autoscaler","name":"workload-variant-autoscaler","serviceCount":5,"netpolCount":0,"hasIngress":true}],"services":[{"id":"NeMo_Guardrails_svc_0","name":"env-port-default","parent":"NeMo_Guardrails","ports":"1235"},{"id":"ai_gateway_payload_processing_svc_0","name":"uvicorn-server","parent":"ai_gateway_payload_processing","ports":"8000"},{"id":"argo_workflows_svc_0","name":"the-service","parent":"argo_workflows","ports":"8666"},{"id":"data_science_pipelines_svc_0","name":"kubeflow-pipelines-profile-controller","parent":"data_science_pipelines","ports":"80"},{"id":"data_science_pipelines_svc_1","name":"squid","parent":"data_science_pipelines","ports":"3128"},{"id":"data_science_pipelines_operator_svc_0","name":"data-science-pipelines-operator-service","parent":"data_science_pipelines_operator","ports":"8080"},{"id":"data_science_pipelines_operator_svc_1","name":"ds-pipeline-workflow-controller-metrics-template-value","parent":"data_science_pipelines_operator","ports":"9090"},{"id":"data_science_pipelines_operator_svc_2","name":"mariadb","parent":"data_science_pipelines_operator","ports":"3306"},{"id":"data_science_pipelines_operator_svc_3","name":"mariadb-template-value","parent":"data_science_pipelines_operator","ports":"3306"},{"id":"data_science_pipelines_operator_svc_4","name":"minio","parent":"data_science_pipelines_operator","ports":"9000,9001"},{"id":"data_science_pipelines_operator_svc_5","name":"minio-service","parent":"data_science_pipelines_operator","ports":"9000"},{"id":"data_science_pipelines_operator_svc_6","name":"minio-template-value","parent":"data_science_pipelines_operator","ports":"9000,80"},{"id":"data_science_pipelines_operator_svc_7","name":"ml-pipeline","parent":"data_science_pipelines_operator","ports":"8443,8888,8887"},{"id":"data_science_pipelines_operator_svc_8","name":"pypi-server","parent":"data_science_pipelines_operator","ports":"8080"},{"id":"data_science_pipelines_operator_svc_9","name":"template-value","parent":"data_science_pipelines_operator","ports":"8443,8888,8887"},{"id":"data_science_pipelines_operator_svc_10","name":"the-service","parent":"data_science_pipelines_operator","ports":"8666"},{"id":"distributed_workloads_svc_0","name":"kuberay-operator","parent":"distributed_workloads","ports":"8080"},{"id":"distributed_workloads_svc_1","name":"training-operator","parent":"distributed_workloads","ports":"8080"},{"id":"distributed_workloads_svc_2","name":"visibility-server","parent":"distributed_workloads","ports":"443"},{"id":"distributed_workloads_svc_3","name":"webhook-service","parent":"distributed_workloads","ports":"443"},{"id":"feast_svc_0","name":"uvicorn-server","parent":"feast","ports":"6566"},{"id":"gateway_api_inference_extension_svc_0","name":"uvicorn-server","parent":"gateway_api_inference_extension","ports":"8000"},{"id":"kserve_svc_0","name":"cli-port-default","parent":"kserve","ports":"80"},{"id":"kserve_svc_1","name":"keda-admission-webhooks","parent":"kserve","ports":"443,8080"},{"id":"kserve_svc_2","name":"keda-metrics-apiserver","parent":"kserve","ports":"443,8080"},{"id":"kserve_svc_3","name":"keda-operator","parent":"kserve","ports":"9666,8080"},{"id":"kserve_svc_4","name":"kserve-controller-manager-metrics-service","parent":"kserve","ports":"8443"},{"id":"kserve_svc_5","name":"kserve-controller-manager-service","parent":"kserve","ports":"8443"},{"id":"kserve_svc_6","name":"kserve-webhook-server-service","parent":"kserve","ports":"443"},{"id":"kserve_svc_7","name":"llmisvc-controller-manager-service","parent":"kserve","ports":"8443"},{"id":"kserve_svc_8","name":"llmisvc-webhook-server-service","parent":"kserve","ports":"443"},{"id":"kserve_svc_9","name":"localmodel-webhook-server-service","parent":"kserve","ports":"443"},{"id":"kserve_svc_10","name":"uvicorn-server","parent":"kserve","ports":"8000"},{"id":"kserve_svc_11","name":"webhook-service","parent":"kserve","ports":"443"},{"id":"kubeflow_svc_0","name":"service","parent":"kubeflow","ports":"443"},{"id":"kubeflow_svc_1","name":"webhook-service","parent":"kubeflow","ports":"443"},{"id":"kuberay_svc_0","name":"kuberay-operator","parent":"kuberay","ports":"8080"},{"id":"kuberay_svc_1","name":"the-service","parent":"kuberay","ports":"8666"},{"id":"kuberay_svc_2","name":"webhook-service","parent":"kuberay","ports":"443"},{"id":"kueue_svc_0","name":"kuberay-operator","parent":"kueue","ports":"8080"},{"id":"kueue_svc_1","name":"the-service","parent":"kueue","ports":"8666"},{"id":"kueue_svc_2","name":"training-operator","parent":"kueue","ports":"8080,443"},{"id":"kueue_svc_3","name":"visibility-server","parent":"kueue","ports":"443"},{"id":"kueue_svc_4","name":"webhook-service","parent":"kueue","ports":"443"},{"id":"llama_stack_k8s_operator_svc_0","name":"ogx-k8s-operator-controller-manager-metrics-service","parent":"llama_stack_k8s_operator","ports":"8443"},{"id":"llama_stack_k8s_operator_svc_1","name":"ogx-k8s-operator-webhook-service","parent":"llama_stack_k8s_operator","ports":"443"},{"id":"llm_d_inference_scheduler_svc_0","name":"${EPP_NAME}","parent":"llm_d_inference_scheduler","ports":"9002,5557,9090"},{"id":"llm_d_inference_scheduler_svc_1","name":"inference-gateway-istio-nodeport","parent":"llm_d_inference_scheduler","ports":"15021,80"},{"id":"llm_d_inference_scheduler_svc_2","name":"istiod-llm-d-gateway","parent":"llm_d_inference_scheduler","ports":"15010,15012,443,15014"},{"id":"llm_d_inference_scheduler_svc_3","name":"service","parent":"llm_d_inference_scheduler","ports":"8080"},{"id":"llm_d_inference_scheduler_svc_4","name":"uvicorn-server","parent":"llm_d_inference_scheduler","ports":"8000"},{"id":"llm_d_routing_sidecar_svc_0","name":"service","parent":"llm_d_routing_sidecar","ports":"8080"},{"id":"mlflow_operator_svc_0","name":"minio-service","parent":"mlflow_operator","ports":"9000"},{"id":"mlflow_operator_svc_1","name":"mlflow-operator-controller-manager-metrics-service","parent":"mlflow_operator","ports":"8443"},{"id":"mlflow_operator_svc_2","name":"postgres-service","parent":"mlflow_operator","ports":"5432"},{"id":"model_registry_svc_0","name":"model-catalog","parent":"model_registry","ports":"8080"},{"id":"modelmesh_serving_svc_0","name":"cli-port-default","parent":"modelmesh_serving","ports":"80"},{"id":"modelmesh_serving_svc_1","name":"etcd","parent":"modelmesh_serving","ports":"2379"},{"id":"modelmesh_serving_svc_2","name":"kserve-controller-manager-service","parent":"modelmesh_serving","ports":"8443"},{"id":"modelmesh_serving_svc_3","name":"kserve-webhook-server-service","parent":"modelmesh_serving","ports":"443"},{"id":"modelmesh_serving_svc_4","name":"modelmesh-controller","parent":"modelmesh_serving","ports":"8080"},{"id":"modelmesh_serving_svc_5","name":"modelmesh-webhook-server-service","parent":"modelmesh_serving","ports":"9443"},{"id":"modelmesh_serving_svc_6","name":"models-server","parent":"modelmesh_serving","ports":"8080"},{"id":"models_as_a_service_svc_0","name":"maas-api","parent":"models_as_a_service","ports":"8080,9090"},{"id":"models_as_a_service_svc_1","name":"payload-processing","parent":"models_as_a_service","ports":"9004"},{"id":"notebooks_svc_0","name":"notebook","parent":"notebooks","ports":"8888"},{"id":"notebooks_downstream_svc_0","name":"notebook","parent":"notebooks_downstream","ports":"8888"},{"id":"odh_cli_svc_0","name":"the-service","parent":"odh_cli","ports":"8666"},{"id":"spark_operator_svc_0","name":"spark-operator-webhook-svc","parent":"spark_operator","ports":"443"},{"id":"spark_operator_svc_1","name":"the-service","parent":"spark_operator","ports":"8666"},{"id":"text_generation_inference_svc_0","name":"inference-server","parent":"text_generation_inference","ports":"8033"},{"id":"trainer_svc_0","name":"webhook-service","parent":"trainer","ports":"443"},{"id":"vllm_cpu_svc_0","name":"cli-port-default","parent":"vllm_cpu","ports":"8000"},{"id":"vllm_cpu_svc_1","name":"disagg_proxy_p2p_nccl_xpyd-server","parent":"vllm_cpu","ports":"10001"},{"id":"vllm_cpu_svc_2","name":"moriio_toy_proxy_server-server","parent":"vllm_cpu","ports":"10001"},{"id":"vllm_gaudi_svc_0","name":"cli-port-default","parent":"vllm_gaudi","ports":"8000"},{"id":"workload_variant_autoscaler_svc_0","name":"keda-admission-webhooks","parent":"workload_variant_autoscaler","ports":"443,8080"},{"id":"workload_variant_autoscaler_svc_1","name":"keda-metrics-apiserver","parent":"workload_variant_autoscaler","ports":"443,8080"},{"id":"workload_variant_autoscaler_svc_2","name":"keda-operator","parent":"workload_variant_autoscaler","ports":"9666,8080"},{"id":"workload_variant_autoscaler_svc_3","name":"uvicorn-server","parent":"workload_variant_autoscaler","ports":"8000"},{"id":"workload_variant_autoscaler_svc_4","name":"webhook-service","parent":"workload_variant_autoscaler","ports":"443"}],"externals":[{"id":"ext_grpc","name":"grpc","type":"grpc"},{"id":"ext_mysql","name":"mysql","type":"database"},{"id":"ext_postgres","name":"postgres","type":"database"},{"id":"ext_redis","name":"redis","type":"database"},{"id":"ext_sqlite","name":"sqlite","type":"database"},{"id":"ext_kafka","name":"kafka","type":"messaging"},{"id":"ext_azure_blob","name":"azure-blob","type":"object-storage"},{"id":"ext_gcs","name":"gcs","type":"object-storage"},{"id":"ext_minio","name":"minio","type":"object-storage"},{"id":"ext_s3","name":"s3","type":"object-storage"},{"id":"ext_etcd","name":"etcd","type":"database"},{"id":"ext_mongodb","name":"mongodb","type":"database"},{"id":"ext_rabbitmq","name":"rabbitmq","type":"messaging"}],"edges":[{"from":"ai_gateway_payload_processing","to":"gateway_api_inference_extension","type":"watches"},{"from":"kserve","to":"gateway_api_inference_extension","type":"watches"},{"from":"kubeflow","to":"data_science_pipelines_operator","type":"module"},{"from":"llm_d_inference_scheduler","to":"gateway_api_inference_extension","type":"watches"},{"from":"model_registry","to":"kserve","type":"watches"},{"from":"modelmesh_serving","to":"kserve","type":"watches"},{"from":"models_as_a_service","to":"kserve","type":"module"},{"from":"models_as_a_service","to":"ai_gateway_payload_processing","type":"watches"},{"from":"odh_cli","to":"opendatahub_operator","type":"module"},{"from":"workload_variant_autoscaler","to":"gateway_api_inference_extension","type":"watches"},{"from":"kserve","to":"modelmesh_serving","type":"module"},{"from":"modelmesh_serving","to":"workload_variant_autoscaler","type":"module"},{"from":"spark_operator","to":"workload_variant_autoscaler","type":"module"},{"from":"ai_gateway_payload_processing","to":"ext_grpc","type":"external"},{"from":"argo_workflows","to":"ext_mysql","type":"external"},{"from":"argo_workflows","to":"ext_postgres","type":"external"},{"from":"argo_workflows","to":"ext_redis","type":"external"},{"from":"argo_workflows","to":"ext_sqlite","type":"external"},{"from":"argo_workflows","to":"ext_grpc","type":"external"},{"from":"argo_workflows","to":"ext_kafka","type":"external"},{"from":"argo_workflows","to":"ext_azure_blob","type":"external"},{"from":"argo_workflows","to":"ext_gcs","type":"external"},{"from":"argo_workflows","to":"ext_minio","type":"external"},{"from":"argo_workflows","to":"ext_s3","type":"external"},{"from":"batch_gateway","to":"ext_postgres","type":"external"},{"from":"batch_gateway","to":"ext_redis","type":"external"},{"from":"batch_gateway","to":"ext_grpc","type":"external"},{"from":"batch_gateway","to":"ext_s3","type":"external"},{"from":"data_science_pipelines","to":"ext_postgres","type":"external"},{"from":"data_science_pipelines","to":"ext_sqlite","type":"external"},{"from":"data_science_pipelines","to":"ext_grpc","type":"external"},{"from":"data_science_pipelines","to":"ext_azure_blob","type":"external"},{"from":"data_science_pipelines","to":"ext_gcs","type":"external"},{"from":"data_science_pipelines","to":"ext_minio","type":"external"},{"from":"data_science_pipelines","to":"ext_s3","type":"external"},{"from":"data_science_pipelines_operator","to":"ext_mysql","type":"external"},{"from":"data_science_pipelines_operator","to":"ext_azure_blob","type":"external"},{"from":"data_science_pipelines_operator","to":"ext_minio","type":"external"},{"from":"data_science_pipelines_operator","to":"ext_s3","type":"external"},{"from":"distributed_workloads","to":"ext_grpc","type":"external"},{"from":"distributed_workloads","to":"ext_azure_blob","type":"external"},{"from":"distributed_workloads","to":"ext_minio","type":"external"},{"from":"distributed_workloads","to":"ext_s3","type":"external"},{"from":"eval_hub","to":"ext_etcd","type":"external"},{"from":"eval_hub","to":"ext_postgres","type":"external"},{"from":"eval_hub","to":"ext_grpc","type":"external"},{"from":"eval_hub","to":"ext_s3","type":"external"},{"from":"feast","to":"ext_postgres","type":"external"},{"from":"feast","to":"ext_redis","type":"external"},{"from":"feast","to":"ext_sqlite","type":"external"},{"from":"feast","to":"ext_grpc","type":"external"},{"from":"feast","to":"ext_gcs","type":"external"},{"from":"feast","to":"ext_s3","type":"external"},{"from":"gateway_api_inference_extension","to":"ext_etcd","type":"external"},{"from":"gateway_api_inference_extension","to":"ext_mongodb","type":"external"},{"from":"gateway_api_inference_extension","to":"ext_mysql","type":"external"},{"from":"gateway_api_inference_extension","to":"ext_sqlite","type":"external"},{"from":"gateway_api_inference_extension","to":"ext_grpc","type":"external"},{"from":"gateway_api_inference_extension","to":"ext_kafka","type":"external"},{"from":"gateway_api_inference_extension","to":"ext_rabbitmq","type":"external"},{"from":"gateway_api_inference_extension","to":"ext_gcs","type":"external"},{"from":"gateway_api_inference_extension","to":"ext_s3","type":"external"},{"from":"kserve","to":"ext_etcd","type":"external"},{"from":"kserve","to":"ext_mongodb","type":"external"},{"from":"kserve","to":"ext_mysql","type":"external"},{"from":"kserve","to":"ext_redis","type":"external"},{"from":"kserve","to":"ext_grpc","type":"external"},{"from":"kserve","to":"ext_azure_blob","type":"external"},{"from":"kserve","to":"ext_gcs","type":"external"},{"from":"kserve","to":"ext_s3","type":"external"},{"from":"kube_auth_proxy","to":"ext_etcd","type":"external"},{"from":"kube_auth_proxy","to":"ext_redis","type":"external"},{"from":"kube_auth_proxy","to":"ext_grpc","type":"external"},{"from":"kube_rbac_proxy","to":"ext_etcd","type":"external"},{"from":"kube_rbac_proxy","to":"ext_grpc","type":"external"},{"from":"kuberay","to":"ext_grpc","type":"external"},{"from":"kuberay","to":"ext_azure_blob","type":"external"},{"from":"kueue","to":"ext_etcd","type":"external"},{"from":"kueue","to":"ext_grpc","type":"external"},{"from":"llm_d_inference_scheduler","to":"ext_redis","type":"external"},{"from":"llm_d_inference_scheduler","to":"ext_grpc","type":"external"},{"from":"llm_d_kv_cache","to":"ext_redis","type":"external"},{"from":"llm_d_kv_cache","to":"ext_grpc","type":"external"},{"from":"mlflow_operator","to":"ext_grpc","type":"external"},{"from":"mlflow_operator","to":"ext_azure_blob","type":"external"},{"from":"model_metadata_collection","to":"ext_sqlite","type":"external"},{"from":"model_metadata_collection","to":"ext_grpc","type":"external"},{"from":"model_metadata_collection","to":"ext_s3","type":"external"},{"from":"model_registry","to":"ext_mongodb","type":"external"},{"from":"model_registry","to":"ext_mysql","type":"external"},{"from":"model_registry","to":"ext_postgres","type":"external"},{"from":"model_registry","to":"ext_sqlite","type":"external"},{"from":"model_registry","to":"ext_gcs","type":"external"},{"from":"modelmesh_runtime_adapter","to":"ext_mysql","type":"external"},{"from":"modelmesh_runtime_adapter","to":"ext_grpc","type":"external"},{"from":"modelmesh_runtime_adapter","to":"ext_azure_blob","type":"external"},{"from":"modelmesh_runtime_adapter","to":"ext_gcs","type":"external"},{"from":"modelmesh_runtime_adapter","to":"ext_s3","type":"external"},{"from":"modelmesh_serving","to":"ext_etcd","type":"external"},{"from":"modelmesh_serving","to":"ext_mysql","type":"external"},{"from":"modelmesh_serving","to":"ext_grpc","type":"external"},{"from":"modelmesh_serving","to":"ext_azure_blob","type":"external"},{"from":"modelmesh_serving","to":"ext_gcs","type":"external"},{"from":"modelmesh_serving","to":"ext_s3","type":"external"},{"from":"odh_cli","to":"ext_grpc","type":"external"},{"from":"openvino_model_server","to":"ext_grpc","type":"external"},{"from":"rest_proxy","to":"ext_grpc","type":"external"},{"from":"spark_operator","to":"ext_postgres","type":"external"},{"from":"spark_operator","to":"ext_sqlite","type":"external"},{"from":"spark_operator","to":"ext_grpc","type":"external"},{"from":"workload_variant_autoscaler","to":"ext_etcd","type":"external"},{"from":"workload_variant_autoscaler","to":"ext_mongodb","type":"external"},{"from":"workload_variant_autoscaler","to":"ext_mysql","type":"external"},{"from":"workload_variant_autoscaler","to":"ext_redis","type":"external"},{"from":"workload_variant_autoscaler","to":"ext_grpc","type":"external"},{"from":"workload_variant_autoscaler","to":"ext_azure_blob","type":"external"},{"from":"workload_variant_autoscaler","to":"ext_gcs","type":"external"}]}
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

    ai_gateway_payload_processing["ai-gateway-payload-processing"]:::comp
    argo_workflows["argo-workflows"]:::comp
    data_science_pipelines_operator["data-science-pipelines-operator"]:::comp
    distributed_workloads["distributed-workloads"]:::comp
    feast["feast"]:::comp
    gateway_api_inference_extension["gateway-api-inference-extension"]:::comp
    kserve["kserve"]:::comp
    kuberay["kuberay"]:::comp
    kueue["kueue"]:::comp
    llm_d_inference_scheduler["llm-d-inference-scheduler"]:::comp
    mlflow_operator["mlflow-operator"]:::comp
    modelmesh_serving["modelmesh-serving"]:::comp
    notebooks["notebooks"]:::comp
    notebooks_downstream["notebooks-downstream"]:::comp
    odh_cli["odh-cli"]:::comp
    spark_operator["spark-operator"]:::comp
    vllm_cpu["vllm-cpu"]:::comp
    vllm_gaudi["vllm-gaudi"]:::comp
    workload_variant_autoscaler["workload-variant-autoscaler"]:::comp

    kueue -.->|"visibility-server"| distributed_workloads
    modelmesh_serving -.->|"kserve-controller-manager-service"| kserve
    modelmesh_serving -.->|"cli-port-default"| kserve
    vllm_cpu -.->|"cli-port-default"| kserve
    vllm_gaudi -.->|"cli-port-default"| kserve
    data_science_pipelines_operator -.->|"the-service"| argo_workflows
    kuberay -.->|"the-service"| argo_workflows
    kueue -.->|"the-service"| argo_workflows
    odh_cli -.->|"the-service"| argo_workflows
    spark_operator -.->|"the-service"| argo_workflows
    notebooks_downstream -.->|"notebook"| notebooks
    kueue -.->|"training-operator"| distributed_workloads
    feast -.->|"uvicorn-server"| ai_gateway_payload_processing
    gateway_api_inference_extension -.->|"uvicorn-server"| ai_gateway_payload_processing
    kserve -.->|"uvicorn-server"| ai_gateway_payload_processing
    llm_d_inference_scheduler -.->|"uvicorn-server"| ai_gateway_payload_processing
    workload_variant_autoscaler -.->|"uvicorn-server"| ai_gateway_payload_processing
    workload_variant_autoscaler -.->|"keda-operator"| kserve
    mlflow_operator -.->|"minio-service"| data_science_pipelines_operator
    kuberay -.->|"kuberay-operator"| distributed_workloads
    kueue -.->|"kuberay-operator"| distributed_workloads
```

## Services by Component

| Component | Services | Webhook (443) | Metrics (8443) | Data |
|-----------|----------|---------------|----------------|------|
| NeMo-Guardrails | 1 | 0 | 0 | 1 |
| ai-gateway-payload-processing | 1 | 0 | 0 | 1 |
| argo-workflows | 1 | 0 | 0 | 1 |
| data-science-pipelines | 2 | 0 | 0 | 2 |
| data-science-pipelines-operator | 11 | 0 | 2 | 9 |
| distributed-workloads | 4 | 2 | 0 | 2 |
| feast | 1 | 0 | 0 | 1 |
| gateway-api-inference-extension | 1 | 0 | 0 | 1 |
| kserve | 12 | 6 | 3 | 3 |
| kubeflow | 2 | 2 | 0 | 0 |
| kuberay | 3 | 1 | 0 | 2 |
| kueue | 5 | 3 | 0 | 2 |
| llama-stack-k8s-operator | 2 | 1 | 1 | 0 |
| llm-d-inference-scheduler | 5 | 1 | 0 | 4 |
| llm-d-routing-sidecar | 1 | 0 | 0 | 1 |
| mlflow-operator | 3 | 0 | 1 | 2 |
| model-registry | 1 | 0 | 0 | 1 |
| modelmesh-serving | 7 | 1 | 1 | 5 |
| models-as-a-service | 2 | 0 | 0 | 2 |
| notebooks | 1 | 0 | 0 | 1 |
| notebooks-downstream | 1 | 0 | 0 | 1 |
| odh-cli | 1 | 0 | 0 | 1 |
| spark-operator | 2 | 1 | 0 | 1 |
| text-generation-inference | 1 | 0 | 0 | 1 |
| trainer | 1 | 1 | 0 | 0 |
| vllm-cpu | 3 | 0 | 0 | 3 |
| vllm-gaudi | 1 | 0 | 0 | 1 |
| workload-variant-autoscaler | 5 | 3 | 0 | 2 |

## Service Detail

Per-component service breakdown with exact port numbers and protocols.

### NeMo-Guardrails (1 services)

| Service | Type | Ports |
|---------|------|-------|
| env-port-default | python-source | 1235/TCP |

### ai-gateway-payload-processing (1 services)

| Service | Type | Ports |
|---------|------|-------|
| uvicorn-server | python-source | 8000/TCP |

### argo-workflows (1 services)

| Service | Type | Ports |
|---------|------|-------|
| the-service | LoadBalancer | 8666/TCP |

### data-science-pipelines (2 services)

| Service | Type | Ports |
|---------|------|-------|
| kubeflow-pipelines-profile-controller | ClusterIP | 80/TCP |
| squid | ClusterIP | 3128/TCP |

### data-science-pipelines-operator (11 services)

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
| the-service | LoadBalancer | 8666/TCP |

### distributed-workloads (4 services)

| Service | Type | Ports |
|---------|------|-------|
| kuberay-operator | ClusterIP | 8080/TCP |
| training-operator | ClusterIP | 8080/TCP |
| visibility-server | ClusterIP | 443/TCP |
| webhook-service | ClusterIP | 443/TCP |

### feast (1 services)

| Service | Type | Ports |
|---------|------|-------|
| uvicorn-server | python-source | 6566/TCP |

### gateway-api-inference-extension (1 services)

| Service | Type | Ports |
|---------|------|-------|
| uvicorn-server | python-source | 8000/TCP |

### kserve (12 services)

| Service | Type | Ports |
|---------|------|-------|
| cli-port-default | python-source | 80/TCP |
| keda-admission-webhooks | ClusterIP | 443/TCP, 8080/TCP |
| keda-metrics-apiserver | ClusterIP | 443/TCP, 8080/TCP |
| keda-operator | ClusterIP | 9666/TCP, 8080/TCP |
| kserve-controller-manager-metrics-service | ClusterIP | 8443/TCP |
| kserve-controller-manager-service | ClusterIP | 8443/TCP |
| kserve-webhook-server-service | ClusterIP | 443/TCP |
| llmisvc-controller-manager-service | ClusterIP | 8443/TCP |
| llmisvc-webhook-server-service | ClusterIP | 443/TCP |
| localmodel-webhook-server-service | ClusterIP | 443/TCP |
| uvicorn-server | python-source | 8000/TCP |
| webhook-service | ClusterIP | 443/TCP |

### kubeflow (2 services)

| Service | Type | Ports |
|---------|------|-------|
| service | ClusterIP | 443/TCP |
| webhook-service | ClusterIP | 443/TCP |

### kuberay (3 services)

| Service | Type | Ports |
|---------|------|-------|
| kuberay-operator | ClusterIP | 8080/TCP |
| the-service | LoadBalancer | 8666/TCP |
| webhook-service | ClusterIP | 443/TCP |

### kueue (5 services)

| Service | Type | Ports |
|---------|------|-------|
| kuberay-operator | ClusterIP | 8080/TCP |
| the-service | LoadBalancer | 8666/TCP |
| training-operator | ClusterIP | 8080/TCP, 443/TCP |
| visibility-server | ClusterIP | 443/TCP |
| webhook-service | ClusterIP | 443/TCP |

### llama-stack-k8s-operator (2 services)

| Service | Type | Ports |
|---------|------|-------|
| ogx-k8s-operator-controller-manager-metrics-service | ClusterIP | 8443/TCP |
| ogx-k8s-operator-webhook-service | ClusterIP | 443/TCP |

### llm-d-inference-scheduler (5 services)

| Service | Type | Ports |
|---------|------|-------|
| ${EPP_NAME} | ClusterIP | 9002/TCP, 5557/TCP, 9090/TCP |
| inference-gateway-istio-nodeport | NodePort | 15021/TCP, 80/TCP |
| istiod-llm-d-gateway | ClusterIP | 15010/TCP, 15012/TCP, 443/TCP, 15014/TCP |
| service | ClusterIP | 8080/TCP |
| uvicorn-server | python-source | 8000/TCP |

### llm-d-routing-sidecar (1 services)

| Service | Type | Ports |
|---------|------|-------|
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

### modelmesh-serving (7 services)

| Service | Type | Ports |
|---------|------|-------|
| cli-port-default | python-source | 80/TCP |
| etcd | ClusterIP | 2379/TCP |
| kserve-controller-manager-service | ClusterIP | 8443/TCP |
| kserve-webhook-server-service | ClusterIP | 443/TCP |
| modelmesh-controller | ClusterIP | 8080/TCP |
| modelmesh-webhook-server-service | ClusterIP | 9443/TCP |
| models-server | python-source | 8080/TCP |

### models-as-a-service (2 services)

| Service | Type | Ports |
|---------|------|-------|
| maas-api | ClusterIP | 8080/TCP, 9090/TCP |
| payload-processing | ClusterIP | 9004/TCP |

### notebooks (1 services)

| Service | Type | Ports |
|---------|------|-------|
| notebook | ClusterIP | 8888/TCP |

### notebooks-downstream (1 services)

| Service | Type | Ports |
|---------|------|-------|
| notebook | ClusterIP | 8888/TCP |

### odh-cli (1 services)

| Service | Type | Ports |
|---------|------|-------|
| the-service | LoadBalancer | 8666/TCP |

### spark-operator (2 services)

| Service | Type | Ports |
|---------|------|-------|
| spark-operator-webhook-svc | ClusterIP | 443/TCP |
| the-service | LoadBalancer | 8666/TCP |

### text-generation-inference (1 services)

| Service | Type | Ports |
|---------|------|-------|
| inference-server | ClusterIP | 8033/TCP |

### trainer (1 services)

| Service | Type | Ports |
|---------|------|-------|
| webhook-service | ClusterIP | 443/TCP |

### vllm-cpu (3 services)

| Service | Type | Ports |
|---------|------|-------|
| cli-port-default | python-source | 8000/TCP |
| disagg_proxy_p2p_nccl_xpyd-server | python-source | 10001/TCP |
| moriio_toy_proxy_server-server | python-source | 10001/TCP |

### vllm-gaudi (1 services)

| Service | Type | Ports |
|---------|------|-------|
| cli-port-default | python-source | 8000/TCP |

### workload-variant-autoscaler (5 services)

| Service | Type | Ports |
|---------|------|-------|
| keda-admission-webhooks | ClusterIP | 443/TCP, 8080/TCP |
| keda-metrics-apiserver | ClusterIP | 443/TCP, 8080/TCP |
| keda-operator | ClusterIP | 9666/TCP, 8080/TCP |
| uvicorn-server | python-source | 8000/TCP |
| webhook-service | ClusterIP | 443/TCP |

## Port Patterns

- **10001/TCP**: disagg_proxy_p2p_nccl_xpyd-server, moriio_toy_proxy_server-server
- **1235/TCP**: env-port-default
- **15010/TCP**: istiod-llm-d-gateway
- **15012/TCP**: istiod-llm-d-gateway
- **15014/TCP**: istiod-llm-d-gateway
- **15021/TCP**: inference-gateway-istio-nodeport
- **2379/TCP**: etcd
- **3128/TCP**: squid
- **3306/TCP**: mariadb, mariadb-template-value
- **443/TCP**: visibility-server, webhook-service, keda-admission-webhooks, keda-metrics-apiserver, kserve-webhook-server-service, llmisvc-webhook-server-service, localmodel-webhook-server-service, webhook-service, service, webhook-service, webhook-service, training-operator, visibility-server, webhook-service, ogx-k8s-operator-webhook-service, istiod-llm-d-gateway, kserve-webhook-server-service, spark-operator-webhook-svc, webhook-service, keda-admission-webhooks, keda-metrics-apiserver, webhook-service
- **5432/TCP**: postgres-service
- **5557/TCP**: ${EPP_NAME}
- **6566/TCP**: uvicorn-server
- **80/TCP**: minio-template-value, kubeflow-pipelines-profile-controller, cli-port-default, inference-gateway-istio-nodeport, cli-port-default
- **8000/TCP**: uvicorn-server, uvicorn-server, uvicorn-server, uvicorn-server, cli-port-default, cli-port-default, uvicorn-server
- **8033/TCP**: inference-server
- **8080/TCP**: data-science-pipelines-operator-service, pypi-server, kuberay-operator, training-operator, keda-admission-webhooks, keda-metrics-apiserver, keda-operator, kuberay-operator, kuberay-operator, training-operator, service, service, model-catalog, modelmesh-controller, models-server, maas-api, keda-admission-webhooks, keda-metrics-apiserver, keda-operator
- **8443/TCP**: ml-pipeline, template-value, kserve-controller-manager-metrics-service, kserve-controller-manager-service, llmisvc-controller-manager-service, ogx-k8s-operator-controller-manager-metrics-service, mlflow-operator-controller-manager-metrics-service, kserve-controller-manager-service
- **8666/TCP**: the-service, the-service, the-service, the-service, the-service, the-service
- **8887/TCP**: ml-pipeline, template-value
- **8888/TCP**: ml-pipeline, template-value, notebook, notebook
- **9000/TCP**: minio, minio-service, minio-template-value, minio-service
- **9001/TCP**: minio
- **9002/TCP**: ${EPP_NAME}
- **9004/TCP**: payload-processing
- **9090/TCP**: ds-pipeline-workflow-controller-metrics-template-value, ${EPP_NAME}, maas-api
- **9443/TCP**: modelmesh-webhook-server-service
- **9666/TCP**: keda-operator, keda-operator

