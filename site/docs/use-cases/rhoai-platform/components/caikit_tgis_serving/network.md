# caikit-tgis-serving: Network

### Services

No services defined.

### Ingress / Routing

| Kind | Name | Hosts | Paths | TLS | Source |
|------|------|-------|-------|-----|--------|
| Gateway | knative-ingress-gateway |  |  | no | [`demo/kserve/custom-manifests/serverless/gateways.yaml`](https://github.com/red-hat-data-services/caikit-tgis-serving/blob/27e5ef01c74822e835e3ae7d55c69d747be718fd/demo/kserve/custom-manifests/serverless/gateways.yaml) |
| Gateway | knative-local-gateway |  |  | no | [`demo/kserve/custom-manifests/serverless/gateways.yaml`](https://github.com/red-hat-data-services/caikit-tgis-serving/blob/27e5ef01c74822e835e3ae7d55c69d747be718fd/demo/kserve/custom-manifests/serverless/gateways.yaml) |

### Network Policies

| Name | Policy Types | Source |
|------|-------------|--------|
| allow-from-openshift-monitoring-ns | Ingress | [`demo/kserve/custom-manifests/metrics/networkpolicy-uwm.yaml`](https://github.com/red-hat-data-services/caikit-tgis-serving/blob/27e5ef01c74822e835e3ae7d55c69d747be718fd/demo/kserve/custom-manifests/metrics/networkpolicy-uwm.yaml) |

## Network Policy Graph

Visual representation of NetworkPolicy rules. Ingress rules show what traffic is allowed into pods, egress rules show what traffic is allowed out.

```mermaid
graph LR
    classDef policy fill:#e74c3c,stroke:#c0392b,color:#fff
    classDef pod fill:#3498db,stroke:#2980b9,color:#fff
    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff

    caikit_tgis_serving["caikit-tgis-serving\nPods"]:::pod
    np_0_allow_from_openshift_monitoring_ns{{"allow-from-openshift-monitoring-ns\nIngress"}}:::policy
    np_0_allow_from_openshift_monitoring_ns --> caikit_tgis_serving
```

