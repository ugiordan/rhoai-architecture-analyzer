# Extractors

The analyzer includes 16 extractors that parse Kubernetes manifests, Go source code, Dockerfiles, and Helm charts.

## Extractor reference

| # | Extractor | Source Patterns | Data Extracted |
|---|-----------|----------------|----------------|
| 1 | CRDs | `config/crd/**`, `deploy/crds/`, `charts/**/crds/`, `manifests/**/crd*` | Group, version, kind, scope, field count, CEL rules |
| 2 | RBAC | `config/rbac/`, `deploy/rbac/`, Go kubebuilder markers | ClusterRoles, bindings, rules, kubebuilder RBAC markers |
| 3 | Services | `**/service*.yaml` | Name, type, ports, selector |
| 4 | Deployments | `**/deployment*.yaml`, `**/manager*.yaml`, `**/statefulset*.yaml` | Containers, security context, env vars, volumes, resources, probes |
| 5 | Network Policies | `**/*networkpolicy*`, `**/*network-polic*`, `**/*netpol*`, `**/network-policies/**` | Pod selector, ingress/egress rules |
| 6 | Controller Watches | `**/*_controller.go`, `**/setup.go`, `**/*reconciler*.go` | For/Owns/Watches with GVK resolution |
| 7 | Dependencies | `go.mod` | Go version, toolchain, modules (direct only), internal ODH deps, replace directives |
| 8 | Secrets | Deployments, services | Secret names, types, references (never values) |
| 9 | Helm | `Chart.yaml`, `values.yaml` | Chart metadata, security-relevant defaults |
| 10 | Dockerfiles | `Dockerfile*`, `Containerfile*` | Base image, stages, USER, EXPOSE, FIPS indicators |
| 11 | Webhooks | `**/webhook*.yaml`, `**/mutating*`, `**/validating*` | Webhook rules, failure policy, side effects |
| 12 | ConfigMaps | `**/configmap*.yaml` | ConfigMap names, data keys |
| 13 | HTTP Endpoints | Go source (`http.HandleFunc`, `mux.Route`, `gin.Engine`) | Method, path, handler, middleware |
| 14 | Ingress | `**/ingress*`, `**/virtualservice*`, `**/httproute*` | Gateway API, Istio, K8s Ingress resources |
| 15 | External Connections | Go source (`sql.Open`, `redis.NewClient`, `grpc.Dial`, `sarama.New*`) | Database, object storage, gRPC, messaging references with credential redaction |
| 16 | Cache Config | Go source (`ctrl.NewManager`, `cache.Options`) | Cache scope, filtered types, disabled types, implicit informers, GOMEMLIMIT |

## YAML extractors

Extractors 1-5 and 8-14 parse Kubernetes YAML manifests. They:

- Walk the repository for matching file patterns
- Parse YAML into typed Kubernetes structs
- Extract relevant fields into the `ComponentArchitecture` struct
- Skip files that don't parse as valid Kubernetes resources

### CRDs extractor

Extracts Custom Resource Definitions including:

- Group, version, kind (GVK)
- Scope (Namespaced/Cluster)
- Field count per version
- CEL validation rules
- Conversion webhook configuration

### RBAC extractor

Combines two sources:

1. **YAML files**: ClusterRoles, Roles, ClusterRoleBindings, RoleBindings from config/rbac/
2. **Go source**: Kubebuilder `// +kubebuilder:rbac:` markers from controller files

The RBAC renderer uses this data to produce the ServiceAccount -> Binding -> Role -> Resource graph.

### Deployments extractor

Extracts detailed deployment information:

- Container images, commands, args
- Security contexts (privileged, capabilities, runAsUser)
- Environment variables (names and sources, never secret values)
- Volume mounts
- Resource requests and limits
- Readiness and liveness probes
- Replica count

### Secrets extractor

Scans deployments and services for secret references. Extracts secret names and types but never extracts secret values.

## Go source extractors

Extractors 6, 13, 15, and 16 parse Go source code.

### Controller Watches extractor

Parses controller `SetupWithManager` functions to extract:

- `For()` (primary watched type)
- `Owns()` (owned child resources)
- `Watches()` (secondary watched types)

Resolves GVKs by following type imports and package aliases.

### HTTP Endpoints extractor

Detects HTTP route registrations from multiple frameworks:

- `net/http`: `HandleFunc`, `Handle`
- `gorilla/mux`: `HandleFunc`, `Methods().Path()`
- `gin`: `GET`, `POST`, etc.
- `chi`: `Get`, `Post`, etc.

Extracts method, path pattern, handler function name, and middleware.

### External Connections extractor

Scans Go source files for references to external services. Detects four categories:

- **Database**: PostgreSQL (`pgx.Connect`, `sql.Open("postgres")`), MySQL, MongoDB (`mongo.Connect`), Redis (`redis.NewClient`), SQLite, etcd (`clientv3.New`)
- **Object storage**: AWS S3 (`s3.NewClient`), MinIO (`minio.New`), GCS (`storage.NewClient`), Azure Blob (`azblob.NewClient`)
- **gRPC**: `grpc.Dial`, `grpc.NewClient`, `grpc.DialContext`
- **Messaging**: Kafka (Sarama, confluent-kafka-go, segmentio/kafka-go), NATS (`nats.Connect`), RabbitMQ (`amqp.Dial`)

Connection strings are automatically redacted: credentials are stripped from URIs (`postgres://***@host:5432/db`), sensitive query parameters are masked, and `user:pass@` patterns are replaced in non-URI targets. Environment variable references (`os.Getenv`, `${}`) are preserved as-is.

Each match includes the enclosing function name (via brace-counting heuristic) and source location.

### Cache Config extractor

Analyzes controller-runtime cache configuration. See [Cache Analysis](../architecture/cache-analysis.md) for details.

## File extractors

### Dockerfiles extractor

Parses Dockerfile/Containerfile syntax:

- Multi-stage build stages and base images
- `USER` directive (non-root detection)
- `EXPOSE` ports
- FIPS indicators (build args, base image names)

### Helm extractor

Reads `Chart.yaml` and `values.yaml`:

- Chart name, version, appVersion
- Security-relevant default values
- Dependencies

### Dependencies extractor

Parses `go.mod`:

- Go version directive (`go 1.25`)
- Toolchain directive (`toolchain go1.25.0`)
- Direct dependencies only (not transitive)
- Internal ODH dependencies (highlighted in diagrams)
- Replace directives (original module, replacement, version)
- Comment lines are filtered to avoid false matches
