# RBAC Surface

50 cluster roles across the platform.

## Permission Scope by Component

How many distinct Kubernetes resource types can each component's most powerful ClusterRole access? A wider scope means the component can read/write more types of resources, which increases its blast radius if compromised. Color: 🔴 wide (>30 types), 🟠 medium (10-30), 🟢 narrow (<10).

<div markdown class="bar-chart-container" style="margin: 1em 0; padding: 1em; border: 1px solid var(--md-default-fg-color--lightest); border-radius: 8px;">

**Widest Role Scope (resource types)**

<div style="display: flex; flex-direction: column; gap: 6px;">
<div style="display: flex; align-items: center; gap: 8px;">
  <span style="min-width: 220px; text-align: right; font-size: 0.85em; white-space: nowrap;">argo-workflows</span>
  <div style="flex: 1; background: var(--md-default-fg-color--lightest); border-radius: 4px; height: 22px; position: relative;">
    <div style="width: 38%; background: #f39c12; height: 100%; border-radius: 4px; min-width: 20px;"></div>
  </div>
  <span style="min-width: 30px; font-size: 0.85em; font-weight: 600;">21</span>
</div>
<div style="display: flex; align-items: center; gap: 8px;">
  <span style="min-width: 220px; text-align: right; font-size: 0.85em; white-space: nowrap;">codeflare-operator</span>
  <div style="flex: 1; background: var(--md-default-fg-color--lightest); border-radius: 4px; height: 22px; position: relative;">
    <div style="width: 61%; background: #e74c3c; height: 100%; border-radius: 4px; min-width: 20px;"></div>
  </div>
  <span style="min-width: 30px; font-size: 0.85em; font-weight: 600;">34</span>
</div>
<div style="display: flex; align-items: center; gap: 8px;">
  <span style="min-width: 220px; text-align: right; font-size: 0.85em; white-space: nowrap;">data-science-pipelines-operator</span>
  <div style="flex: 1; background: var(--md-default-fg-color--lightest); border-radius: 4px; height: 22px; position: relative;">
    <div style="width: 100%; background: #e74c3c; height: 100%; border-radius: 4px; min-width: 20px;"></div>
  </div>
  <span style="min-width: 30px; font-size: 0.85em; font-weight: 600;">55</span>
</div>
<div style="display: flex; align-items: center; gap: 8px;">
  <span style="min-width: 220px; text-align: right; font-size: 0.85em; white-space: nowrap;">llama-stack-k8s-operator</span>
  <div style="flex: 1; background: var(--md-default-fg-color--lightest); border-radius: 4px; height: 22px; position: relative;">
    <div style="width: 30%; background: #f39c12; height: 100%; border-radius: 4px; min-width: 20px;"></div>
  </div>
  <span style="min-width: 30px; font-size: 0.85em; font-weight: 600;">17</span>
</div>
<div style="display: flex; align-items: center; gap: 8px;">
  <span style="min-width: 220px; text-align: right; font-size: 0.85em; white-space: nowrap;">model-registry</span>
  <div style="flex: 1; background: var(--md-default-fg-color--lightest); border-radius: 4px; height: 22px; position: relative;">
    <div style="width: 5%; background: #27ae60; height: 100%; border-radius: 4px; min-width: 20px;"></div>
  </div>
  <span style="min-width: 30px; font-size: 0.85em; font-weight: 600;">3</span>
</div>
<div style="display: flex; align-items: center; gap: 8px;">
  <span style="min-width: 220px; text-align: right; font-size: 0.85em; white-space: nowrap;">model-registry-operator</span>
  <div style="flex: 1; background: var(--md-default-fg-color--lightest); border-radius: 4px; height: 22px; position: relative;">
    <div style="width: 49%; background: #f39c12; height: 100%; border-radius: 4px; min-width: 20px;"></div>
  </div>
  <span style="min-width: 30px; font-size: 0.85em; font-weight: 600;">27</span>
</div>
<div style="display: flex; align-items: center; gap: 8px;">
  <span style="min-width: 220px; text-align: right; font-size: 0.85em; white-space: nowrap;">odh-dashboard</span>
  <div style="flex: 1; background: var(--md-default-fg-color--lightest); border-radius: 4px; height: 22px; position: relative;">
    <div style="width: 72%; background: #e74c3c; height: 100%; border-radius: 4px; min-width: 20px;"></div>
  </div>
  <span style="min-width: 30px; font-size: 0.85em; font-weight: 600;">40</span>
</div>
<div style="display: flex; align-items: center; gap: 8px;">
  <span style="min-width: 220px; text-align: right; font-size: 0.85em; white-space: nowrap;">odh-model-controller</span>
  <div style="flex: 1; background: var(--md-default-fg-color--lightest); border-radius: 4px; height: 22px; position: relative;">
    <div style="width: 81%; background: #e74c3c; height: 100%; border-radius: 4px; min-width: 20px;"></div>
  </div>
  <span style="min-width: 30px; font-size: 0.85em; font-weight: 600;">45</span>
</div>
<div style="display: flex; align-items: center; gap: 8px;">
  <span style="min-width: 220px; text-align: right; font-size: 0.85em; white-space: nowrap;">opendatahub-operator</span>
  <div style="flex: 1; background: var(--md-default-fg-color--lightest); border-radius: 4px; height: 22px; position: relative;">
    <div style="width: 54%; background: #f39c12; height: 100%; border-radius: 4px; min-width: 20px;"></div>
  </div>
  <span style="min-width: 30px; font-size: 0.85em; font-weight: 600;">30</span>
</div>
<div style="display: flex; align-items: center; gap: 8px;">
  <span style="min-width: 220px; text-align: right; font-size: 0.85em; white-space: nowrap;">training-operator</span>
  <div style="flex: 1; background: var(--md-default-fg-color--lightest); border-radius: 4px; height: 22px; position: relative;">
    <div style="width: 56%; background: #e74c3c; height: 100%; border-radius: 4px; min-width: 20px;"></div>
  </div>
  <span style="min-width: 30px; font-size: 0.85em; font-weight: 600;">31</span>
</div>
</div>
</div>

## RBAC Binding Graph

Subject-to-role bindings across all platform components. Edge direction shows who has access to what.

```mermaid
graph LR
    classDef role fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef subject fill:#3498db,stroke:#2980b9,color:#fff

    sa_argo["argo\nServiceAccount"]:::subject
    role_argo_cluster_role["argo-cluster-role"]:::role
    sa_argo -->|argo-workflows| role_argo_cluster_role
    sa_argo_server["argo-server\nServiceAccount"]:::subject
    role_argo_server_cluster_role["argo-server-cluster-role"]:::role
    sa_argo_server -->|argo-workflows| role_argo_server_cluster_role
    sa_controller_manager["controller-manager\nServiceAccount"]:::subject
    role_manager_role["manager-role"]:::role
    sa_controller_manager -->|codeflare-operator| role_manager_role
    role_manager_argo_role["manager-argo-role"]:::role
    sa_controller_manager -->|data-science-pipelines-operator| role_manager_argo_role
    sa_controller_manager -->|data-science-pipelines-operator| role_manager_role
    sa_controller_manager -->|llama-stack-k8s-operator| role_manager_role
    role_proxy_role["proxy-role"]:::role
    sa_controller_manager -->|llama-stack-k8s-operator| role_proxy_role
    sa_controller_manager -->|model-registry-operator| role_manager_role
    sa_controller_manager -->|model-registry-operator| role_proxy_role
    role_metrics_auth_role["metrics-auth-role"]:::role
    sa_controller_manager -->|model-registry| role_metrics_auth_role
    sa_model_registry_ui["model-registry-ui\nServiceAccount"]:::subject
    role_model_registry_create_sars["model-registry-create-sars"]:::role
    sa_model_registry_ui -->|model-registry| role_model_registry_create_sars
    role_model_registry_manager_role["model-registry-manager-role"]:::role
    sa_controller_manager -->|model-registry| role_model_registry_manager_role
    role_model_registry_retrieve_clusterrolebindings["model-registry-retrieve-clusterrolebindings"]:::role
    sa_model_registry_ui -->|model-registry| role_model_registry_retrieve_clusterrolebindings
    role_model_registry_ui_services_reader["model-registry-ui-services-reader"]:::role
    sa_model_registry_ui -->|model-registry| role_model_registry_ui_services_reader
    sa_odh_dashboard["odh-dashboard\nServiceAccount"]:::subject
    role_odh_dashboard["odh-dashboard"]:::role
    sa_odh_dashboard -->|odh-dashboard| role_odh_dashboard
    role_system_auth_delegator["system:auth-delegator"]:::role
    sa_odh_dashboard -->|odh-dashboard| role_system_auth_delegator
    role_cluster_monitoring_view["cluster-monitoring-view"]:::role
    sa_odh_dashboard -->|odh-dashboard| role_cluster_monitoring_view
    sa_controller_manager -->|odh-model-controller| role_metrics_auth_role
    sa_odh_model_controller["odh-model-controller\nServiceAccount"]:::subject
    role_odh_model_controller_role["odh-model-controller-role"]:::role
    sa_odh_model_controller -->|odh-model-controller| role_odh_model_controller_role
    sa_odh_model_controller -->|odh-model-controller| role_proxy_role
    role_controller_manager_role["controller-manager-role"]:::role
    sa_controller_manager -->|opendatahub-operator| role_controller_manager_role
    sa_controller_manager -->|opendatahub-operator| role_proxy_role
    sa_training_operator["training-operator\nServiceAccount"]:::subject
    role_training_operator["training-operator"]:::role
    sa_training_operator -->|training-operator| role_training_operator
```

## Roles by Component

| Component | Roles | Widest Role | Resources | Scope |
|-----------|-------|-------------|-----------|-------|
| argo-workflows | 5 | argo-cluster-role | 21 | medium |
| codeflare-operator | 3 | manager-role | 34 | **wide** |
| data-science-pipelines-operator | 4 | manager-role | 55 | **wide** |
| llama-stack-k8s-operator | 5 | manager-role | 17 | medium |
| model-registry | 6 | model-registry-manager-role | 3 | narrow |
| model-registry-operator | 6 | manager-role | 27 | medium |
| odh-dashboard | 1 | odh-dashboard | 40 | **wide** |
| odh-model-controller | 7 | odh-model-controller-role | 45 | **wide** |
| opendatahub-operator | 7 | controller-manager-role | 30 | medium |
| training-operator | 6 | training-operator | 31 | **wide** |

