# llama-stack-k8s-operator: Cache Architecture

Controller-runtime cache configuration controls which Kubernetes resources are cached in-memory. Misconfigured caches (cluster-wide watches on high-cardinality types without filters) are a primary cause of operator OOM kills.

## Cache Architecture

### Manager Configuration

| Property | Value |
|----------|-------|
| Manager file | `main.go` |
| Cache scope | cluster-wide |
| DefaultTransform | yes |
| GOMEMLIMIT | 800MiB |
| Memory limit | 1Gi |

### Filtered Types

| Type | Filter Kind | Filter |
|------|-------------|--------|
| appsv1.Deployment | label | label selector |
| autoscalingv2.HorizontalPodAutoscaler | label | label selector |
| corev1.ConfigMap | label | controllers.WatchLabelKey=controllers.WatchLabelValue (constants, resolved at runtime) |
| corev1.PersistentVolumeClaim | label | label selector |
| corev1.Service | label | label selector |
| networkingv1.Ingress | label | label selector |
| networkingv1.NetworkPolicy | label | label selector |
| policyv1.PodDisruptionBudget | label | label selector |

### Issues

- GOMEMLIMIT ratio 78.1% is below recommended 80% minimum (GC cannot pressure-tune effectively)
- Type LlamaStackDistribution is watched but has no cache filter (cluster-wide informer)

