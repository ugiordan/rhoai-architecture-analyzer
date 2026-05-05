# odh-dashboard: Network

## Service Map

```mermaid
graph LR
    classDef svc fill:#2ecc71,stroke:#27ae60,color:#fff
    classDef test fill:#95a5a6,stroke:#7f8c8d,color:#fff
    classDef component fill:#3498db,stroke:#2980b9,color:#fff
    classDef ext fill:#e74c3c,stroke:#c0392b,color:#fff

    odh_dashboard["odh-dashboard"]:::component
    odh_dashboard --> svc_0["odh-dashboard\nClusterIP: 8443/TCP"]:::svc
    odh_dashboard --> svc_1["workspaces-backend\nClusterIP: 4000/TCP"]:::svc
    odh_dashboard --> svc_2["workspaces-controller-metrics-service\nClusterIP: 8080/TCP"]:::svc
    odh_dashboard --> svc_3["workspaces-frontend\nClusterIP: 8080/TCP"]:::svc
    odh_dashboard --> svc_4["workspaces-webhook-service\nClusterIP: 443/TCP"]:::svc
    odh_dashboard -.-> ext_s3[["s3\nobject-storage"]]:::ext
```

### Services

| Name | Type | Ports | Source |
|------|------|-------|--------|
| odh-dashboard | ClusterIP | 8443/TCP | [`manifests/core-bases/base/service.yaml`](https://github.com/red-hat-data-services/odh-dashboard/blob/9f2858e35f91324c8d5f4021189b10a82fa78147/manifests/core-bases/base/service.yaml) |
| workspaces-backend | ClusterIP | 4000/TCP | [`packages/notebooks/upstream/workspaces/backend/manifests/kustomize/base/service.yaml`](https://github.com/red-hat-data-services/odh-dashboard/blob/9f2858e35f91324c8d5f4021189b10a82fa78147/packages/notebooks/upstream/workspaces/backend/manifests/kustomize/base/service.yaml) |
| workspaces-controller-metrics-service | ClusterIP | 8080/TCP | [`packages/notebooks/upstream/workspaces/controller/manifests/kustomize/components/prometheus/service.yaml`](https://github.com/red-hat-data-services/odh-dashboard/blob/9f2858e35f91324c8d5f4021189b10a82fa78147/packages/notebooks/upstream/workspaces/controller/manifests/kustomize/components/prometheus/service.yaml) |
| workspaces-frontend | ClusterIP | 8080/TCP | [`packages/notebooks/upstream/workspaces/frontend/manifests/kustomize/base/service.yaml`](https://github.com/red-hat-data-services/odh-dashboard/blob/9f2858e35f91324c8d5f4021189b10a82fa78147/packages/notebooks/upstream/workspaces/frontend/manifests/kustomize/base/service.yaml) |
| workspaces-webhook-service | ClusterIP | 443/TCP | [`packages/notebooks/upstream/workspaces/controller/manifests/kustomize/base/webhook/service.yaml`](https://github.com/red-hat-data-services/odh-dashboard/blob/9f2858e35f91324c8d5f4021189b10a82fa78147/packages/notebooks/upstream/workspaces/controller/manifests/kustomize/base/webhook/service.yaml) |

### Ingress / Routing

| Kind | Name | Hosts | Paths | TLS | Source |
|------|------|-------|-------|-----|--------|
| Gateway | kubeflow-gateway |  |  | no | [`packages/notebooks/upstream/developing/manifests/istio-gateway/gateway.yaml`](https://github.com/red-hat-data-services/odh-dashboard/blob/9f2858e35f91324c8d5f4021189b10a82fa78147/packages/notebooks/upstream/developing/manifests/istio-gateway/gateway.yaml) |
| HTTPRoute | odh-dashboard |  | / | no | [`manifests/core-bases/base/httproute.yaml`](https://github.com/red-hat-data-services/odh-dashboard/blob/9f2858e35f91324c8d5f4021189b10a82fa78147/manifests/core-bases/base/httproute.yaml) |
| Route | odh-dashboard |  |  | yes | [`manifests/core-bases/base/routes.yaml`](https://github.com/red-hat-data-services/odh-dashboard/blob/9f2858e35f91324c8d5f4021189b10a82fa78147/manifests/core-bases/base/routes.yaml) |

### Network Policies

| Name | Policy Types | Source |
|------|-------------|--------|
| allow-perses-operator-access | Ingress | [`packages/observability/setup/network-policy-perses-operator-access.yaml`](https://github.com/red-hat-data-services/odh-dashboard/blob/9f2858e35f91324c8d5f4021189b10a82fa78147/packages/observability/setup/network-policy-perses-operator-access.yaml) |
| dashboard-perses-access | Ingress | [`manifests/observability/odh/network-policy.yaml`](https://github.com/red-hat-data-services/odh-dashboard/blob/9f2858e35f91324c8d5f4021189b10a82fa78147/manifests/observability/odh/network-policy.yaml) |
| dashboard-perses-access | Ingress | [`manifests/observability/rhoai/network-policy.yaml`](https://github.com/red-hat-data-services/odh-dashboard/blob/9f2858e35f91324c8d5f4021189b10a82fa78147/manifests/observability/rhoai/network-policy.yaml) |
| odh-dashboard-allow-ports | Ingress, Egress | [`manifests/modular-architecture/networkpolicy.yaml`](https://github.com/red-hat-data-services/odh-dashboard/blob/9f2858e35f91324c8d5f4021189b10a82fa78147/manifests/modular-architecture/networkpolicy.yaml) |
| workspaces-controller | Ingress | [`packages/notebooks/upstream/workspaces/controller/manifests/kustomize/components/istio/network-policy.yaml`](https://github.com/red-hat-data-services/odh-dashboard/blob/9f2858e35f91324c8d5f4021189b10a82fa78147/packages/notebooks/upstream/workspaces/controller/manifests/kustomize/components/istio/network-policy.yaml) |

## Network Policy Graph

Visual representation of NetworkPolicy rules. Ingress rules show what traffic is allowed into pods, egress rules show what traffic is allowed out.

```mermaid
graph LR
    classDef policy fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef pod fill:#3498db,stroke:#2980b9,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff

    odh_dashboard["odh-dashboard\nPods"]:::pod
    np_0_allow_perses_operator_access{{"allow-perses-operator-access\nIngress"}}:::policy
    np_0_allow_perses_operator_access --> odh_dashboard
    np_1_dashboard_perses_access{{"dashboard-perses-access\nIngress"}}:::policy
    np_1_dashboard_perses_access --> odh_dashboard
    np_2_dashboard_perses_access{{"dashboard-perses-access\nIngress"}}:::policy
    np_2_dashboard_perses_access --> odh_dashboard
    np_3_odh_dashboard_allow_ports{{"odh-dashboard-allow-ports\nIngress, Egress"}}:::policy
    np_3_odh_dashboard_allow_ports --> odh_dashboard
    np_4_workspaces_controller{{"workspaces-controller\nIngress"}}:::policy
    np_4_workspaces_controller --> odh_dashboard
```

