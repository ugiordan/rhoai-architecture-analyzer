# kueue: Cache Architecture

Controller-runtime cache configuration controls which Kubernetes resources are cached in-memory. Misconfigured caches (cluster-wide watches on high-cardinality types without filters) are a primary cause of operator OOM kills.

## Cache Architecture

### Manager Configuration

| Property | Value |
|----------|-------|
| Manager file | `cmd/kueue/main.go` |
| Cache scope | cluster-wide |
| DefaultTransform | no |
| Memory limit | 512Mi |

### Implicit Informers (OOM Risk)

| Type | Source | Risk |
|------|--------|------|
| client.ListOptions | `pkg/util/testing/core.go:66` | medium |

### Issues

- No GOMEMLIMIT set in deployment (Go GC cannot pressure-tune). Set GOMEMLIMIT to 80-90% of container memory limit for optimal GC behavior
- No cache configuration: all informers are cluster-wide (OOM risk). See https://book.kubebuilder.io/reference/watching-resources/filtering for cache filtering patterns
- Type AdmissionCheck is watched but has no cache filter (cluster-wide informer)
- Type ClusterQueue is watched but has no cache filter (cluster-wide informer)
- Type LeaderWorkerSet is watched but has no cache filter (cluster-wide informer)
- Type LimitRange is watched but has no cache filter (cluster-wide informer)
- Type LocalQueue is watched but has no cache filter (cluster-wide informer)
- Type Namespace is watched but has no cache filter (cluster-wide informer)
- Type Pod is watched but has no cache filter (cluster-wide informer)
- Type ProvisioningRequest is watched but has no cache filter (cluster-wide informer)
- Type ProvisioningRequestConfig is watched but has no cache filter (cluster-wide informer)
- Type ResourceFlavor is watched but has no cache filter (cluster-wide informer)
- Type RuntimeClass is watched but has no cache filter (cluster-wide informer)
- Type StatefulSet is watched but has no cache filter (cluster-wide informer)
- Type Workload is watched but has no cache filter (cluster-wide informer)

