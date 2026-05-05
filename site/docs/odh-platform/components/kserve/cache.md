# kserve: Cache Architecture

Controller-runtime cache configuration controls which Kubernetes resources are cached in-memory. Misconfigured caches (cluster-wide watches on high-cardinality types without filters) are a primary cause of operator OOM kills.

## Cache Architecture

### Manager Configuration

| Property | Value |
|----------|-------|
| Manager file | `cmd/manager/main.go` |
| Cache scope | cluster-wide |
| DefaultTransform | no |
| Memory limit | 300Mi |

### Filtered Types

| Type | Filter Kind | Filter |
|------|-------------|--------|
| appsv1.Deployment | label | label selector |
| autoscalingv2.HorizontalPodAutoscaler | label | label selector |
| corev1.ConfigMap | label | label selector |
| corev1.Pod | label | label selector |
| corev1.Secret | label | label selector |

### Issues

- No DefaultTransform: managedFields cached for all objects (wasted memory). Add cache.DefaultTransform to strip managedFields and reduce memory footprint
- No GOMEMLIMIT set in deployment (Go GC cannot pressure-tune). Set GOMEMLIMIT to 80-90% of container memory limit for optimal GC behavior
- Type ClusterServingRuntime is watched but has no cache filter (cluster-wide informer)
- Type Gateway is watched but has no cache filter (cluster-wide informer)
- Type HTTPRoute is watched but has no cache filter (cluster-wide informer)
- Type InferenceGraph is watched but has no cache filter (cluster-wide informer)
- Type InferencePool is watched but has no cache filter (cluster-wide informer)
- Type InferenceService is watched but has no cache filter (cluster-wide informer)
- Type Ingress is watched but has no cache filter (cluster-wide informer)
- Type Job is watched but has no cache filter (cluster-wide informer)
- Type LLMInferenceService is watched but has no cache filter (cluster-wide informer)
- Type LLMInferenceServiceConfig is watched but has no cache filter (cluster-wide informer)
- Type LeaderWorkerSet is watched but has no cache filter (cluster-wide informer)
- Type LocalModelCache is watched but has no cache filter (cluster-wide informer)
- Type LocalModelNamespaceCache is watched but has no cache filter (cluster-wide informer)
- Type LocalModelNode is watched but has no cache filter (cluster-wide informer)
- Type Node is watched but has no cache filter (cluster-wide informer)
- Type OpenTelemetryCollector is watched but has no cache filter (cluster-wide informer)
- Type PersistentVolume is watched but has no cache filter (cluster-wide informer)
- Type PersistentVolumeClaim is watched but has no cache filter (cluster-wide informer)
- Type ScaledObject is watched but has no cache filter (cluster-wide informer)
- Type Service is watched but has no cache filter (cluster-wide informer)
- Type ServingRuntime is watched but has no cache filter (cluster-wide informer)
- Type TrainedModel is watched but has no cache filter (cluster-wide informer)
- Type VariantAutoscaling is watched but has no cache filter (cluster-wide informer)
- Type VirtualService is watched but has no cache filter (cluster-wide informer)

