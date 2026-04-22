# modelmesh-serving: Cache Architecture

Controller-runtime cache configuration controls which Kubernetes resources are cached in-memory. Misconfigured caches (cluster-wide watches on high-cardinality types without filters) are a primary cause of operator OOM kills.

## Cache Architecture

### Manager Configuration

| Property | Value |
|----------|-------|
| Manager file | `main.go` |
| Cache scope | namespace-scoped |
| DefaultTransform | no |
| Memory limit | 512Mi |

### Issues

- No GOMEMLIMIT set in deployment (Go GC cannot pressure-tune)
- No cache configuration: all informers are cluster-wide (OOM risk)
- Type ClusterServingRuntime is watched but has no cache filter (cluster-wide informer)
- Type ConfigMap is watched but has no cache filter (cluster-wide informer)
- Type Deployment is watched but has no cache filter (cluster-wide informer)
- Type InferenceService is watched but has no cache filter (cluster-wide informer)
- Type Namespace is watched but has no cache filter (cluster-wide informer)
- Type Predictor is watched but has no cache filter (cluster-wide informer)
- Type Secret is watched but has no cache filter (cluster-wide informer)
- Type Service is watched but has no cache filter (cluster-wide informer)
- Type ServiceMonitor is watched but has no cache filter (cluster-wide informer)
- Type ServingRuntime is watched but has no cache filter (cluster-wide informer)

