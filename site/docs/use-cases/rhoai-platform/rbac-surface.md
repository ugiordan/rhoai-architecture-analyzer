# RBAC Surface

19 cluster roles across the platform.

## Permission Scope by Component

How many distinct Kubernetes resource types can each component's most powerful ClusterRole access? A wider scope means the component can read/write more types of resources, which increases its blast radius if compromised. Color: 🔴 wide (>30 types), 🟠 medium (10-30), 🟢 narrow (<10).

<div markdown class="bar-chart-container" style="margin: 1em 0; padding: 1em; border: 1px solid var(--md-default-fg-color--lightest); border-radius: 8px;">

**Widest Role Scope (resource types)**

<div style="display: flex; flex-direction: column; gap: 6px;">
<div style="display: flex; align-items: center; gap: 8px;">
  <span style="min-width: 220px; text-align: right; font-size: 0.85em; white-space: nowrap;">data-science-pipelines</span>
  <div style="flex: 1; background: var(--md-default-fg-color--lightest); border-radius: 4px; height: 22px; position: relative;">
    <div style="width: 23%; background: #f39c12; height: 100%; border-radius: 4px; min-width: 20px;"></div>
  </div>
  <span style="min-width: 30px; font-size: 0.85em; font-weight: 600;">13</span>
</div>
<div style="display: flex; align-items: center; gap: 8px;">
  <span style="min-width: 220px; text-align: right; font-size: 0.85em; white-space: nowrap;">data-science-pipelines-operator</span>
  <div style="flex: 1; background: var(--md-default-fg-color--lightest); border-radius: 4px; height: 22px; position: relative;">
    <div style="width: 100%; background: #e74c3c; height: 100%; border-radius: 4px; min-width: 20px;"></div>
  </div>
  <span style="min-width: 30px; font-size: 0.85em; font-weight: 600;">55</span>
</div>
<div style="display: flex; align-items: center; gap: 8px;">
  <span style="min-width: 220px; text-align: right; font-size: 0.85em; white-space: nowrap;">kserve</span>
  <div style="flex: 1; background: var(--md-default-fg-color--lightest); border-radius: 4px; height: 22px; position: relative;">
    <div style="width: 81%; background: #e74c3c; height: 100%; border-radius: 4px; min-width: 20px;"></div>
  </div>
  <span style="min-width: 30px; font-size: 0.85em; font-weight: 600;">45</span>
</div>
</div>
</div>

## RBAC Binding Graph

Subject-to-role bindings across all platform components. Edge direction shows who has access to what.

```mermaid
graph LR
    classDef role fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef subject fill:#3498db,stroke:#2980b9,color:#fff

    sa_controller_manager["controller-manager\nServiceAccount"]:::subject
    role_manager_argo_role["manager-argo-role"]:::role
    sa_controller_manager -->|data-science-pipelines-operator| role_manager_argo_role
    role_manager_role["manager-role"]:::role
    sa_controller_manager -->|data-science-pipelines-operator| role_manager_role
    sa_kubeflow_pipelines_cache["kubeflow-pipelines-cache\nServiceAccount"]:::subject
    role_kubeflow_pipelines_cache_role["kubeflow-pipelines-cache-role"]:::role
    sa_kubeflow_pipelines_cache -->|data-science-pipelines| role_kubeflow_pipelines_cache_role
    sa_kubeflow_pipelines_cache_deployer_sa["kubeflow-pipelines-cache-deployer-sa\nServiceAccount"]:::subject
    role_kubeflow_pipelines_cache_deployer_clusterrole["kubeflow-pipelines-cache-deployer-clusterrole"]:::role
    sa_kubeflow_pipelines_cache_deployer_sa -->|data-science-pipelines| role_kubeflow_pipelines_cache_deployer_clusterrole
    sa_kubeflow_pipelines_metadata_writer["kubeflow-pipelines-metadata-writer\nServiceAccount"]:::subject
    role_kubeflow_pipelines_metadata_writer_role["kubeflow-pipelines-metadata-writer-role"]:::role
    sa_kubeflow_pipelines_metadata_writer -->|data-science-pipelines| role_kubeflow_pipelines_metadata_writer_role
    sa_meta_controller_service["meta-controller-service\nServiceAccount"]:::subject
    role_kubeflow_metacontroller["kubeflow-metacontroller"]:::role
    sa_meta_controller_service -->|data-science-pipelines| role_kubeflow_metacontroller
    sa_ml_pipeline["ml-pipeline\nServiceAccount"]:::subject
    role_ml_pipeline["ml-pipeline"]:::role
    sa_ml_pipeline -->|data-science-pipelines| role_ml_pipeline
    sa_ml_pipeline_persistenceagent["ml-pipeline-persistenceagent\nServiceAccount"]:::subject
    role_ml_pipeline_persistenceagent_role["ml-pipeline-persistenceagent-role"]:::role
    sa_ml_pipeline_persistenceagent -->|data-science-pipelines| role_ml_pipeline_persistenceagent_role
    sa_ml_pipeline_scheduledworkflow["ml-pipeline-scheduledworkflow\nServiceAccount"]:::subject
    role_ml_pipeline_scheduledworkflow_role["ml-pipeline-scheduledworkflow-role"]:::role
    sa_ml_pipeline_scheduledworkflow -->|data-science-pipelines| role_ml_pipeline_scheduledworkflow_role
    sa_ml_pipeline_ui["ml-pipeline-ui\nServiceAccount"]:::subject
    role_ml_pipeline_ui["ml-pipeline-ui"]:::role
    sa_ml_pipeline_ui -->|data-science-pipelines| role_ml_pipeline_ui
    sa_ml_pipeline_viewer_crd_service_account["ml-pipeline-viewer-crd-service-account\nServiceAccount"]:::subject
    role_ml_pipeline_viewer_controller_role["ml-pipeline-viewer-controller-role"]:::role
    sa_ml_pipeline_viewer_crd_service_account -->|data-science-pipelines| role_ml_pipeline_viewer_controller_role
    sa_kserve_controller_manager["kserve-controller-manager\nServiceAccount"]:::subject
    role_kserve_manager_role["kserve-manager-role"]:::role
    sa_kserve_controller_manager -->|kserve| role_kserve_manager_role
    role_kserve_proxy_role["kserve-proxy-role"]:::role
    sa_kserve_controller_manager -->|kserve| role_kserve_proxy_role
```

## Roles by Component

| Component | Roles | Widest Role | Resources | Scope |
|-----------|-------|-------------|-----------|-------|
| data-science-pipelines | 13 | aggregate-to-kubeflow-pipelines-edit | 13 | medium |
| data-science-pipelines-operator | 4 | manager-role | 55 | **wide** |
| kserve | 2 | kserve-manager-role | 45 | **wide** |

