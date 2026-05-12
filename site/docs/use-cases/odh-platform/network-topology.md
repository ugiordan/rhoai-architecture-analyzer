# Network Topology

22 Kubernetes services across the platform.

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
  {"components":[{"id":"data_science_pipelines","name":"data-science-pipelines","serviceCount":2,"netpolCount":1,"hasIngress":true},{"id":"data_science_pipelines_operator","name":"data-science-pipelines-operator","serviceCount":10,"netpolCount":2,"hasIngress":true},{"id":"kserve","name":"kserve","serviceCount":7,"netpolCount":0,"hasIngress":true},{"id":"modelmesh_serving","name":"modelmesh-serving","serviceCount":3,"netpolCount":4,"hasIngress":false}],"services":[{"id":"data_science_pipelines_svc_0","name":"kubeflow-pipelines-profile-controller","parent":"data_science_pipelines","ports":"80"},{"id":"data_science_pipelines_svc_1","name":"squid","parent":"data_science_pipelines","ports":"3128"},{"id":"data_science_pipelines_operator_svc_0","name":"data-science-pipelines-operator-service","parent":"data_science_pipelines_operator","ports":"8080"},{"id":"data_science_pipelines_operator_svc_1","name":"ds-pipeline-workflow-controller-metrics-template-value","parent":"data_science_pipelines_operator","ports":"9090"},{"id":"data_science_pipelines_operator_svc_2","name":"mariadb","parent":"data_science_pipelines_operator","ports":"3306"},{"id":"data_science_pipelines_operator_svc_3","name":"mariadb-template-value","parent":"data_science_pipelines_operator","ports":"3306"},{"id":"data_science_pipelines_operator_svc_4","name":"minio","parent":"data_science_pipelines_operator","ports":"9000,9001"},{"id":"data_science_pipelines_operator_svc_5","name":"minio-service","parent":"data_science_pipelines_operator","ports":"9000"},{"id":"data_science_pipelines_operator_svc_6","name":"minio-template-value","parent":"data_science_pipelines_operator","ports":"9000,80"},{"id":"data_science_pipelines_operator_svc_7","name":"ml-pipeline","parent":"data_science_pipelines_operator","ports":"8443,8888,8887"},{"id":"data_science_pipelines_operator_svc_8","name":"pypi-server","parent":"data_science_pipelines_operator","ports":"8080"},{"id":"data_science_pipelines_operator_svc_9","name":"template-value","parent":"data_science_pipelines_operator","ports":"8443,8888,8887"},{"id":"kserve_svc_0","name":"cli-port-default","parent":"kserve","ports":"80"},{"id":"kserve_svc_1","name":"kserve-controller-manager-metrics-service","parent":"kserve","ports":"8443"},{"id":"kserve_svc_2","name":"kserve-controller-manager-service","parent":"kserve","ports":"8443"},{"id":"kserve_svc_3","name":"kserve-webhook-server-service","parent":"kserve","ports":"443"},{"id":"kserve_svc_4","name":"llmisvc-controller-manager-service","parent":"kserve","ports":"8443"},{"id":"kserve_svc_5","name":"llmisvc-webhook-server-service","parent":"kserve","ports":"443"},{"id":"kserve_svc_6","name":"localmodel-webhook-server-service","parent":"kserve","ports":"443"},{"id":"modelmesh_serving_svc_0","name":"etcd","parent":"modelmesh_serving","ports":"2379"},{"id":"modelmesh_serving_svc_1","name":"modelmesh-controller","parent":"modelmesh_serving","ports":"8080"},{"id":"modelmesh_serving_svc_2","name":"modelmesh-webhook-server-service","parent":"modelmesh_serving","ports":"9443"}],"externals":[{"id":"ext_grpc","name":"grpc","type":"grpc"},{"id":"ext_minio","name":"minio","type":"object-storage"},{"id":"ext_s3","name":"s3","type":"object-storage"},{"id":"ext_mysql","name":"mysql","type":"database"},{"id":"ext_azure_blob","name":"azure-blob","type":"object-storage"},{"id":"ext_gcs","name":"gcs","type":"object-storage"}],"edges":[{"from":"modelmesh_serving","to":"kserve","type":"watches"},{"from":"data_science_pipelines","to":"kserve","type":"module"},{"from":"data_science_pipelines","to":"ext_grpc","type":"external"},{"from":"data_science_pipelines","to":"ext_minio","type":"external"},{"from":"data_science_pipelines","to":"ext_s3","type":"external"},{"from":"data_science_pipelines_operator","to":"ext_mysql","type":"external"},{"from":"data_science_pipelines_operator","to":"ext_minio","type":"external"},{"from":"kserve","to":"ext_azure_blob","type":"external"},{"from":"kserve","to":"ext_gcs","type":"external"},{"from":"kserve","to":"ext_s3","type":"external"},{"from":"modelmesh_serving","to":"ext_azure_blob","type":"external"}]}
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

## Services by Component

| Component | Services | Webhook (443) | Metrics (8443) | Data |
|-----------|----------|---------------|----------------|------|
| data-science-pipelines | 2 | 0 | 0 | 2 |
| data-science-pipelines-operator | 10 | 0 | 2 | 8 |
| kserve | 7 | 3 | 3 | 1 |
| modelmesh-serving | 3 | 0 | 0 | 3 |

## Service Detail

Per-component service breakdown with exact port numbers and protocols.

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

### kserve (7 services)

| Service | Type | Ports |
|---------|------|-------|
| cli-port-default | python-source | 80/TCP |
| kserve-controller-manager-metrics-service | ClusterIP | 8443/TCP |
| kserve-controller-manager-service | ClusterIP | 8443/TCP |
| kserve-webhook-server-service | ClusterIP | 443/TCP |
| llmisvc-controller-manager-service | ClusterIP | 8443/TCP |
| llmisvc-webhook-server-service | ClusterIP | 443/TCP |
| localmodel-webhook-server-service | ClusterIP | 443/TCP |

### modelmesh-serving (3 services)

| Service | Type | Ports |
|---------|------|-------|
| etcd | ClusterIP | 2379/TCP |
| modelmesh-controller | ClusterIP | 8080/TCP |
| modelmesh-webhook-server-service | ClusterIP | 9443/TCP |

## Port Patterns

- **2379/TCP**: etcd
- **3128/TCP**: squid
- **3306/TCP**: mariadb, mariadb-template-value
- **443/TCP**: kserve-webhook-server-service, llmisvc-webhook-server-service, localmodel-webhook-server-service
- **80/TCP**: cli-port-default, minio-template-value, kubeflow-pipelines-profile-controller
- **8080/TCP**: data-science-pipelines-operator-service, pypi-server, modelmesh-controller
- **8443/TCP**: kserve-controller-manager-metrics-service, kserve-controller-manager-service, llmisvc-controller-manager-service, ml-pipeline, template-value
- **8887/TCP**: ml-pipeline, template-value
- **8888/TCP**: ml-pipeline, template-value
- **9000/TCP**: minio, minio-service, minio-template-value
- **9001/TCP**: minio
- **9090/TCP**: ds-pipeline-workflow-controller-metrics-template-value
- **9443/TCP**: modelmesh-webhook-server-service

