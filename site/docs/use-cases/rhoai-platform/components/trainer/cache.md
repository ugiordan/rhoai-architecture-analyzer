# trainer: Cache Architecture

Controller-runtime cache configuration controls which Kubernetes resources are cached in-memory. Misconfigured caches (cluster-wide watches on high-cardinality types without filters) are a primary cause of operator OOM kills.

## Cache Architecture

### Manager Configuration

| Property | Value |
|----------|-------|
| Manager file | `cmd/trainer-controller-manager/main.go` |
| Cache scope | cluster-wide |
| DefaultTransform | no |
| Memory limit | 4Gi |

### Implicit Informers (OOM Risk)

| Type | Source | Risk |
|------|--------|------|
| corev1.Secret | `pkg/runtime/framework/plugins/mpi/mpi.go:256` | **HIGH** |

### Issues

- Implicit informer for corev1.Secret via client.Get at pkg/runtime/framework/plugins/mpi/mpi.go:256 (cluster-wide, OOM risk). This bypasses cache filters and creates a full cluster-wide watch
- No GOMEMLIMIT set in deployment (Go GC cannot pressure-tune). Set GOMEMLIMIT to 80-90% of container memory limit for optimal GC behavior
- No cache configuration: all informers are cluster-wide (OOM risk). See https://book.kubebuilder.io/reference/watching-resources/filtering for cache filtering patterns
- Type ElasticQuota is watched but has no cache filter (cluster-wide informer)
- Type Job is watched but has no cache filter (cluster-wide informer)
- Type JobSet is watched but has no cache filter (cluster-wide informer)
- Type Pod is watched but has no cache filter (cluster-wide informer)
- Type PodGroup is watched but has no cache filter (cluster-wide informer)
- Type Service is watched but has no cache filter (cluster-wide informer)

