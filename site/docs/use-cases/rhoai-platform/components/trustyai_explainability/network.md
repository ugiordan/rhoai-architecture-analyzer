# trustyai-explainability: Network

### Services

No services defined.

### Ingress / Routing

| Kind | Name | Hosts | Paths | TLS | Source |
|------|------|-------|-------|-----|--------|
| Route | trustyai |  |  | yes | [`explainability-service/manifests/base/route.yaml`](https://github.com/red-hat-data-services/trustyai-explainability/blob/ee0b2b22cb42f0f60431192315488f3b195137d9/explainability-service/manifests/base/route.yaml) |

!!! warning "No Network Policies"
    No NetworkPolicy resources found. All pod-to-pod traffic is allowed by default.

