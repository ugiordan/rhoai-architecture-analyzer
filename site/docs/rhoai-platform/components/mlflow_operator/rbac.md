# mlflow-operator: RBAC

ServiceAccount bindings, roles, and resource permissions.

## RBAC Overview

This component defines a large RBAC surface (119 diagram lines). The graph below groups roles by permission scope.

```mermaid
graph LR
    classDef wide fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef medium fill:#f39c12,stroke:#d68910,color:#fff
    classDef narrow fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef subject fill:#3498db,stroke:#2980b9,color:#fff

    subgraph med["Medium Scope (10-30)"]
    mlflow_edit["mlflow-edit\n13 resources"]:::medium
    end
    subgraph nar["Narrow Scope (<10)"]
    manager_role["manager-role\n10 resources"]:::narrow
    metrics_auth_role["metrics-auth-role\n2 resources"]:::narrow
    metrics_reader["metrics-reader"]:::narrow
    mlflow_integration["mlflow-integration\n5 resources"]:::narrow
    mlflow_view["mlflow-view\n8 resources"]:::narrow
    leader_election_role["leader-election-role\n3 resources"]:::narrow
    manager_role["manager-role\n8 resources"]:::narrow
    end

    subj_controller_manager["controller-manager\nServiceAccount"]:::subject
    subj_controller_manager -->|binds| manager_role
    subj_controller_manager -->|binds| metrics_auth_role
    subj_controller_manager -->|binds| leader_election_role
    subj_controller_manager -->|binds| manager_role
```

## Bindings

Subject-to-role mappings defining who has access to what.

| Binding | Type | Role | Subject |
|---------|------|------|---------|
| manager-rolebinding | ClusterRoleBinding | manager-role | ServiceAccount/controller-manager |
| metrics-auth-rolebinding | ClusterRoleBinding | metrics-auth-role | ServiceAccount/controller-manager |
| leader-election-rolebinding | RoleBinding | leader-election-role | ServiceAccount/controller-manager |
| manager-rolebinding | RoleBinding | manager-role | ServiceAccount/controller-manager |

## Role Details

Per-rule breakdown of API groups, resources, and verbs for each role.

| Role | Kind | API Groups | Resources | Verbs |
|------|------|------------|-----------|-------|
| manager-role | ClusterRole |  | namespaces | get, list, watch |
| manager-role | ClusterRole |  | secrets | get, list, watch |
| manager-role | ClusterRole |  | consolelinks | create, delete, get, list, patch, update, watch |
| manager-role | ClusterRole |  | httproutes | create, delete, get, list, patch, update, watch |
| manager-role | ClusterRole |  | mlflowconfigs | get, list, watch |
| manager-role | ClusterRole |  | mlflows | create, delete, get, list, patch, update, watch |
| manager-role | ClusterRole |  | mlflows/finalizers | update |
| manager-role | ClusterRole |  | mlflows/status | get, patch, update |
| manager-role | ClusterRole |  | clusterrolebindings, clusterroles | create, delete, get, list, patch, update, watch |
| metrics-auth-role | ClusterRole |  | tokenreviews | create |
| metrics-auth-role | ClusterRole |  | subjectaccessreviews | create |
| metrics-reader | ClusterRole |  |  | get |
| mlflow-edit | ClusterRole |  | mlflows | get, list, watch, create, delete, deletecollection, patch, update |
| mlflow-edit | ClusterRole |  | mlflows/finalizers | patch, update |
| mlflow-edit | ClusterRole |  | mlflowconfigs | create, delete, deletecollection, patch, update |
| mlflow-edit | ClusterRole |  | gatewaysecrets | get, list |
| mlflow-edit | ClusterRole |  | datasets, experiments, registeredmodels, gatewaysecrets, gatewayendpoints, gatewaymodeldefinitions | create, update, delete |
| mlflow-edit | ClusterRole |  | gatewaysecrets/use, gatewayendpoints/use, gatewaymodeldefinitions/use | create |
| mlflow-integration | ClusterRole |  | datasets, experiments, registeredmodels | get, list, create, update |
| mlflow-integration | ClusterRole |  | gatewayendpoints | get, list |
| mlflow-integration | ClusterRole |  | gatewayendpoints/use | create |
| mlflow-view | ClusterRole |  | mlflows | get, list, watch |
| mlflow-view | ClusterRole |  | mlflowconfigs | get, list, watch |
| mlflow-view | ClusterRole |  | mlflows/status | get, list, watch |
| mlflow-view | ClusterRole |  | datasets, experiments, registeredmodels, gatewayendpoints, gatewaymodeldefinitions | get, list |
| leader-election-role | Role |  | configmaps | get, list, watch, create, update, patch, delete |
| leader-election-role | Role |  | leases | get, list, watch, create, update, patch, delete |
| leader-election-role | Role |  | events | create, patch |
| manager-role | Role |  | configmaps, persistentvolumeclaims, secrets, serviceaccounts, services | create, delete, get, list, patch, update, watch |
| manager-role | Role |  | deployments | create, delete, get, list, patch, update, watch |
| manager-role | Role |  | networkpolicies | create, delete, get, list, patch, update, watch |
| manager-role | Role |  | servicemonitors | create, delete, get, list, patch, update, watch |

### Cluster Roles

| Name | Resources | Verbs | Source |
|------|-----------|-------|--------|
| manager-role | namespaces | get, list, watch | [`config/rbac/role.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/rbac/role.yaml) |
| manager-role | secrets | get, list, watch | [`config/rbac/role.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/rbac/role.yaml) |
| manager-role | consolelinks | create, delete, get, list, patch, update, watch | [`config/rbac/role.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/rbac/role.yaml) |
| manager-role | httproutes | create, delete, get, list, patch, update, watch | [`config/rbac/role.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/rbac/role.yaml) |
| manager-role | mlflowconfigs | get, list, watch | [`config/rbac/role.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/rbac/role.yaml) |
| manager-role | mlflows | create, delete, get, list, patch, update, watch | [`config/rbac/role.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/rbac/role.yaml) |
| manager-role | mlflows/finalizers | update | [`config/rbac/role.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/rbac/role.yaml) |
| manager-role | mlflows/status | get, patch, update | [`config/rbac/role.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/rbac/role.yaml) |
| manager-role | clusterrolebindings, clusterroles | create, delete, get, list, patch, update, watch | [`config/rbac/role.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/rbac/role.yaml) |
| metrics-auth-role | tokenreviews | create | [`config/rbac/metrics_auth_role.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/rbac/metrics_auth_role.yaml) |
| metrics-auth-role | subjectaccessreviews | create | [`config/rbac/metrics_auth_role.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/rbac/metrics_auth_role.yaml) |
| metrics-reader |  | get | [`config/rbac/metrics_reader_role.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/rbac/metrics_reader_role.yaml) |
| mlflow-edit | mlflows | get, list, watch, create, delete, deletecollection, patch, update | [`config/rbac/mlflow_aggregate_roles.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/rbac/mlflow_aggregate_roles.yaml) |
| mlflow-edit | mlflows/finalizers | patch, update | [`config/rbac/mlflow_aggregate_roles.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/rbac/mlflow_aggregate_roles.yaml) |
| mlflow-edit | mlflowconfigs | create, delete, deletecollection, patch, update | [`config/rbac/mlflow_aggregate_roles.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/rbac/mlflow_aggregate_roles.yaml) |
| mlflow-edit | gatewaysecrets | get, list | [`config/rbac/mlflow_aggregate_roles.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/rbac/mlflow_aggregate_roles.yaml) |
| mlflow-edit | datasets, experiments, registeredmodels, gatewaysecrets, gatewayendpoints, gatewaymodeldefinitions | create, update, delete | [`config/rbac/mlflow_aggregate_roles.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/rbac/mlflow_aggregate_roles.yaml) |
| mlflow-edit | gatewaysecrets/use, gatewayendpoints/use, gatewaymodeldefinitions/use | create | [`config/rbac/mlflow_aggregate_roles.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/rbac/mlflow_aggregate_roles.yaml) |
| mlflow-integration | datasets, experiments, registeredmodels | get, list, create, update | [`config/rbac/mlflow_integration_role.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/rbac/mlflow_integration_role.yaml) |
| mlflow-integration | gatewayendpoints | get, list | [`config/rbac/mlflow_integration_role.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/rbac/mlflow_integration_role.yaml) |
| mlflow-integration | gatewayendpoints/use | create | [`config/rbac/mlflow_integration_role.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/rbac/mlflow_integration_role.yaml) |
| mlflow-view | mlflows | get, list, watch | [`config/rbac/mlflow_aggregate_roles.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/rbac/mlflow_aggregate_roles.yaml) |
| mlflow-view | mlflowconfigs | get, list, watch | [`config/rbac/mlflow_aggregate_roles.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/rbac/mlflow_aggregate_roles.yaml) |
| mlflow-view | mlflows/status | get, list, watch | [`config/rbac/mlflow_aggregate_roles.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/rbac/mlflow_aggregate_roles.yaml) |
| mlflow-view | datasets, experiments, registeredmodels, gatewayendpoints, gatewaymodeldefinitions | get, list | [`config/rbac/mlflow_aggregate_roles.yaml`](https://github.com/opendatahub-io/mlflow-operator/blob/682055b5deae5d1cc92c0a24270aee8400704084/config/rbac/mlflow_aggregate_roles.yaml) |

