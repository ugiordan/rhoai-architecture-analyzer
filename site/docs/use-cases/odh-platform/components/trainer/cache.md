# trainer: Cache Architecture

Controller-runtime cache configuration controls which Kubernetes resources are cached in-memory. Misconfigured caches (cluster-wide watches on high-cardinality types without filters) are a primary cause of operator OOM kills.

## Cache Architecture

### Manager Configuration

| Property | Value |
|----------|-------|
| Manager file | `cmd/trainer-controller-manager/main.go` |
| Cache scope | cluster-wide |
| DefaultTransform | no |

### Implicit Informers (OOM Risk)

| Type | Source | Risk |
|------|--------|------|
| corev1.Secret | `pkg/runtime/framework/plugins/mpi/mpi.go:256` | **HIGH** |

### Issues

- Implicit informer for corev1.Secret via client.Get at pkg/runtime/framework/plugins/mpi/mpi.go:256 (cluster-wide, OOM risk). This bypasses cache filters and creates a full cluster-wide watch

