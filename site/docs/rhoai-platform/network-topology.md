# Network Topology

29 Kubernetes services across the platform.

## Network Topology Graph

Interactive service mesh view of the platform. Drag nodes to rearrange, hover to highlight connections, click for details. Double-click background to fit all.

<div class="topology-toolbar">
  <button data-action="fit" title="Fit to view">Fit</button>
  <button data-action="zoom-in" title="Zoom in">+</button>
  <button data-action="zoom-out" title="Zoom out">&minus;</button>
  <button data-action="relayout" title="Re-run layout">Relayout</button>
</div>
<div class="cytoscape-topology">
  <script type="application/json">
  {"components":[{"id":"data_science_pipelines_operator","name":"data-science-pipelines-operator","serviceCount":3,"netpolCount":0,"hasIngress":false},{"id":"kserve","name":"kserve","serviceCount":6,"netpolCount":1,"hasIngress":true},{"id":"kube_auth_proxy","name":"kube-auth-proxy","serviceCount":0,"netpolCount":0,"hasIngress":false},{"id":"kube_rbac_proxy","name":"kube-rbac-proxy","serviceCount":1,"netpolCount":0,"hasIngress":false},{"id":"kuberay","name":"kuberay","serviceCount":2,"netpolCount":0,"hasIngress":true},{"id":"model_registry_operator","name":"model-registry-operator","serviceCount":1,"netpolCount":0,"hasIngress":false},{"id":"notebooks","name":"notebooks","serviceCount":1,"netpolCount":0,"hasIngress":false},{"id":"odh_dashboard","name":"odh-dashboard","serviceCount":5,"netpolCount":5,"hasIngress":true},{"id":"odh_model_controller","name":"odh-model-controller","serviceCount":1,"netpolCount":0,"hasIngress":false},{"id":"opendatahub_operator","name":"opendatahub-operator","serviceCount":8,"netpolCount":2,"hasIngress":true},{"id":"trustyai_service_operator","name":"trustyai-service-operator","serviceCount":0,"netpolCount":0,"hasIngress":false}],"services":[{"id":"data_science_pipelines_operator_svc_0","name":"mariadb","parent":"data_science_pipelines_operator","ports":"3306"},{"id":"data_science_pipelines_operator_svc_1","name":"minio","parent":"data_science_pipelines_operator","ports":"9000,9001"},{"id":"data_science_pipelines_operator_svc_2","name":"pypi-server","parent":"data_science_pipelines_operator","ports":"8080"},{"id":"kserve_svc_0","name":"llmisvc-controller-manager-service","parent":"kserve","ports":"8443"},{"id":"kserve_svc_1","name":"kserve-controller-manager-service","parent":"kserve","ports":"8443"},{"id":"kserve_svc_2","name":"llmisvc-webhook-server-service","parent":"kserve","ports":"443"},{"id":"kserve_svc_3","name":"localmodel-webhook-server-service","parent":"kserve","ports":"443"},{"id":"kserve_svc_4","name":"kserve-webhook-server-service","parent":"kserve","ports":"443"},{"id":"kserve_svc_5","name":"webhook-service","parent":"kserve","ports":"443"},{"id":"kube_rbac_proxy_svc_0","name":"kube-rbac-proxy","parent":"kube_rbac_proxy","ports":"8443"},{"id":"kuberay_svc_0","name":"kuberay-operator","parent":"kuberay","ports":"8080"},{"id":"kuberay_svc_1","name":"webhook-service","parent":"kuberay","ports":"443"},{"id":"model_registry_operator_svc_0","name":"webhook-service","parent":"model_registry_operator","ports":"443"},{"id":"notebooks_svc_0","name":"notebook","parent":"notebooks","ports":"8888"},{"id":"odh_dashboard_svc_0","name":"odh-dashboard","parent":"odh_dashboard","ports":"8443"},{"id":"odh_dashboard_svc_1","name":"workspaces-backend","parent":"odh_dashboard","ports":"4000"},{"id":"odh_dashboard_svc_2","name":"workspaces-webhook-service","parent":"odh_dashboard","ports":"443"},{"id":"odh_dashboard_svc_3","name":"workspaces-controller-metrics-service","parent":"odh_dashboard","ports":"8080"},{"id":"odh_dashboard_svc_4","name":"workspaces-frontend","parent":"odh_dashboard","ports":"8080"},{"id":"odh_model_controller_svc_0","name":"odh-model-controller-webhook-service","parent":"odh_model_controller","ports":"443"},{"id":"opendatahub_operator_svc_0","name":"webhook-service","parent":"opendatahub_operator","ports":"443"},{"id":"opendatahub_operator_svc_1","name":"odh-dashboard","parent":"opendatahub_operator","ports":"8443"},{"id":"opendatahub_operator_svc_2","name":"kserve-controller-manager-service","parent":"opendatahub_operator","ports":"8443"},{"id":"opendatahub_operator_svc_3","name":"kserve-webhook-server-service","parent":"opendatahub_operator","ports":"443"},{"id":"opendatahub_operator_svc_4","name":"odh-model-controller-webhook-service","parent":"opendatahub_operator","ports":"443"},{"id":"opendatahub_operator_svc_5","name":"kuberay-operator","parent":"opendatahub_operator","ports":"8080"},{"id":"opendatahub_operator_svc_6","name":"training-operator","parent":"opendatahub_operator","ports":"8080,443"},{"id":"opendatahub_operator_svc_7","name":"service","parent":"opendatahub_operator","ports":"443"}],"externals":[{"id":"ext_mysql","name":"mysql","type":"database"},{"id":"ext_minio","name":"minio","type":"object-storage"},{"id":"ext_azure_blob","name":"azure-blob","type":"object-storage"},{"id":"ext_gcs","name":"gcs","type":"object-storage"},{"id":"ext_s3","name":"s3","type":"object-storage"},{"id":"ext_redis","name":"redis","type":"database"},{"id":"ext_grpc","name":"grpc","type":"grpc"}],"edges":[{"from":"odh_dashboard","to":"llama_stack_k8s_operator","type":"module"},{"from":"odh_dashboard","to":"mlflow_go","type":"module"},{"from":"odh_model_controller","to":"kserve","type":"watches"},{"from":"opendatahub_operator","to":"models_as_a_service","type":"module"},{"from":"kserve","to":"kube_rbac_proxy","type":"sidecar"},{"from":"kube_auth_proxy","to":"kube_rbac_proxy","type":"sidecar"},{"from":"odh_dashboard","to":"kube_rbac_proxy","type":"sidecar"},{"from":"opendatahub_operator","to":"kube_rbac_proxy","type":"sidecar"},{"from":"opendatahub_operator","to":"kserve","type":"sidecar"},{"from":"data_science_pipelines_operator","to":"ext_mysql","type":"external"},{"from":"data_science_pipelines_operator","to":"ext_minio","type":"external"},{"from":"kserve","to":"ext_azure_blob","type":"external"},{"from":"kserve","to":"ext_gcs","type":"external"},{"from":"kserve","to":"ext_s3","type":"external"},{"from":"kube_auth_proxy","to":"ext_redis","type":"external"},{"from":"kuberay","to":"ext_grpc","type":"external"},{"from":"odh_dashboard","to":"ext_s3","type":"external"},{"from":"opendatahub_operator","to":"ext_s3","type":"external"}]}
  </script>
</div>
<div class="topology-legend">
  <span><span class="swatch" style="background:#3498db"></span> Component</span>
  <span><span class="swatch" style="background:#2ecc71"></span> Service</span>
  <span><span class="swatch" style="background:#e74c3c"></span> External</span>
  <span><span class="line-swatch" style="background:#e74c3c"></span> CRD Watch</span>
  <span><span class="line-swatch" style="background:#9b59b6"></span> Sidecar</span>
  <span><span class="line-swatch" style="background:#95a5a6;border-top:2px dashed #95a5a6;height:0"></span> Module</span>
  <span><span class="line-swatch" style="background:#e67e22;border-top:2px dotted #e67e22;height:0"></span> External Conn</span>
</div>

## Cross-Component Service References

Services referenced across component boundaries. When component A defines a service that component B also references, it indicates a deployment dependency.

```mermaid
graph LR
    classDef comp fill:#3498db,stroke:#2980b9,color:#fff

    kserve["kserve"]:::comp
    kuberay["kuberay"]:::comp
    odh_dashboard["odh-dashboard"]:::comp
    opendatahub_operator["opendatahub-operator"]:::comp

    opendatahub_operator -.->|"kserve-controller-manager-service"| kserve
    opendatahub_operator -.->|"kuberay-operator"| kuberay
    opendatahub_operator -.->|"odh-dashboard"| odh_dashboard
```

## Services by Component

| Component | Services | Webhook (443) | Metrics (8443) | Data |
|-----------|----------|---------------|----------------|------|
| data-science-pipelines-operator | 3 | 0 | 0 | 3 |
| kserve | 6 | 4 | 2 | 0 |
| kube-auth-proxy | 1 | 0 | 1 | 0 |
| kube-rbac-proxy | 1 | 0 | 1 | 0 |
| kuberay | 2 | 1 | 0 | 1 |
| model-registry-operator | 1 | 1 | 0 | 0 |
| notebooks | 1 | 0 | 0 | 1 |
| odh-dashboard | 5 | 1 | 1 | 3 |
| odh-model-controller | 1 | 1 | 0 | 0 |
| opendatahub-operator | 8 | 5 | 2 | 1 |

## Service Detail

Per-component service breakdown with exact port numbers and protocols.

### data-science-pipelines-operator (3 services)

| Service | Type | Ports |
|---------|------|-------|
| mariadb | ClusterIP | 3306/TCP |
| minio | ClusterIP | 9000/TCP, 9001/TCP |
| pypi-server | ClusterIP | 8080/TCP |

### kserve (6 services)

| Service | Type | Ports |
|---------|------|-------|
| llmisvc-controller-manager-service | ClusterIP | 8443/TCP |
| kserve-controller-manager-service | ClusterIP | 8443/TCP |
| llmisvc-webhook-server-service | ClusterIP | 443/TCP |
| localmodel-webhook-server-service | ClusterIP | 443/TCP |
| kserve-webhook-server-service | ClusterIP | 443/TCP |
| webhook-service | ClusterIP | 443/TCP |

### kube-auth-proxy (1 services)

| Service | Type | Ports |
|---------|------|-------|
| kube-rbac-proxy | ClusterIP | 8443/TCP |

### kube-rbac-proxy (1 services)

| Service | Type | Ports |
|---------|------|-------|
| kube-rbac-proxy | ClusterIP | 8443/TCP |

### kuberay (2 services)

| Service | Type | Ports |
|---------|------|-------|
| kuberay-operator | ClusterIP | 8080/TCP |
| webhook-service | ClusterIP | 443/TCP |

### model-registry-operator (1 services)

| Service | Type | Ports |
|---------|------|-------|
| webhook-service | ClusterIP | 443/TCP |

### notebooks (1 services)

| Service | Type | Ports |
|---------|------|-------|
| notebook | ClusterIP | 8888/TCP |

### odh-dashboard (5 services)

| Service | Type | Ports |
|---------|------|-------|
| odh-dashboard | ClusterIP | 8443/TCP |
| workspaces-backend | ClusterIP | 4000/TCP |
| workspaces-webhook-service | ClusterIP | 443/TCP |
| workspaces-controller-metrics-service | ClusterIP | 8080/TCP |
| workspaces-frontend | ClusterIP | 8080/TCP |

### odh-model-controller (1 services)

| Service | Type | Ports |
|---------|------|-------|
| odh-model-controller-webhook-service | ClusterIP | 443/TCP |

### opendatahub-operator (8 services)

| Service | Type | Ports |
|---------|------|-------|
| webhook-service | ClusterIP | 443/TCP |
| odh-dashboard | ClusterIP | 8443/TCP |
| kserve-controller-manager-service | ClusterIP | 8443/TCP |
| kserve-webhook-server-service | ClusterIP | 443/TCP |
| odh-model-controller-webhook-service | ClusterIP | 443/TCP |
| kuberay-operator | ClusterIP | 8080/TCP |
| training-operator | ClusterIP | 8080/TCP, 443/TCP |
| service | ClusterIP | 443/TCP |

## Port Patterns

- **3306/TCP**: mariadb
- **4000/TCP**: workspaces-backend
- **443/TCP**: llmisvc-webhook-server-service, localmodel-webhook-server-service, kserve-webhook-server-service, webhook-service, webhook-service, webhook-service, workspaces-webhook-service, odh-model-controller-webhook-service, webhook-service, kserve-webhook-server-service, odh-model-controller-webhook-service, training-operator, service
- **8080/TCP**: pypi-server, kuberay-operator, workspaces-controller-metrics-service, workspaces-frontend, kuberay-operator, training-operator
- **8443/TCP**: llmisvc-controller-manager-service, kserve-controller-manager-service, kube-rbac-proxy, kube-rbac-proxy, odh-dashboard, odh-dashboard, kserve-controller-manager-service
- **8888/TCP**: notebook
- **9000/TCP**: minio
- **9001/TCP**: minio

