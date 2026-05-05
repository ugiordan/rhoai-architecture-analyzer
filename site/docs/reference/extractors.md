# Extractors

The analyzer includes 22 extractors that parse Kubernetes manifests, Go source code, Dockerfiles, and Helm charts.

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
| 16 | Feature Gates | Go source (`DefaultMutableFeatureGate.Add`, `featuregate.Feature` consts) | Gate name, default state, pre-release stage, source location |
| 17 | Cache Config | Go source (`ctrl.NewManager`, `cache.Options`) | Cache scope, filtered types, disabled types, implicit informers, GOMEMLIMIT |
| 18 | Operator Config | Go source (const/var blocks in controllers, pkg/config) | Classified constants: images, ports, timeouts, env vars, resources, name patterns |
| 19 | Reconcile Sequences | Go source (`Reconcile()` methods) | Ordered sub-resource reconciliation steps with conditional guards |
| 20 | Prometheus Metrics | Go source (`prometheus.New*`, `promauto.New*`) | Metric name, type (gauge/counter/histogram/summary), help, labels, namespace |
| 21 | Status Conditions | Go source (const blocks in controllers, API types) | Condition type constants, associated reason constants, source location |
| 22 | Platform Detection | Go source (controllers, reconcilers, config packages) | Capability structs (IsOpenShift, HasRoute), API discovery checks, conditional resource creation |

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

Extractors 6, 13, 15-22 parse Go source code.

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

### Feature Gates extractor

Scans Go source for feature gate definitions and registrations. Detects:

- **Gate registrations**: `DefaultMutableFeatureGate.Add()`, `MutableFeatureGate.Add()`, `featuregate.NewFeatureGate()`
- **Gate constants**: `const MyFeature featuregate.Feature = "MyFeature"` declarations
- **Runtime overrides**: `DefaultMutableFeatureGate.Set("GateName=true")` calls

For each gate, extracts:

- Gate name (resolved from constant if defined)
- Default enabled/disabled state
- Pre-release stage (Alpha, Beta, GA, Deprecated)
- Source file and line number
- Whether the gate is set via runtime `Set()` rather than static `Add()`

Deduplicates gates by name across files. Test files and vendor directories are skipped.

This data feeds into the CPG upgrade domain query **CGA-U03** (ungated-feature), which detects functions that check unregistered gates or feature-related functions that lack gate checks entirely.

### Cache Config extractor

Analyzes controller-runtime cache configuration. See [Cache Analysis](../architecture/cache-analysis.md) for details.

### Operator Config extractor

Scans Go const/var blocks in controller, config, and root-level source files. Classifies each constant by category using a precedence-based heuristic:

- **image**: Container image references (registry paths, `:tag` suffixes)
- **port**: Port numbers (numeric values 1-65535 with port-related names)
- **timeout**: Duration values (`time.Second`, `time.Minute`, etc.)
- **env_var**: Environment variable names (UPPER_SNAKE_CASE)
- **resource**: Kubernetes resource quantities (`100m`, `256Mi`, etc.)
- **name_pattern**: Kubernetes object names (lowercase with hyphens)
- **general**: Everything else

Deduplicates against status condition constants to avoid overlap. Skips iota-only blocks, test files, and vendor directories.

### Reconcile Sequences extractor

Parses `Reconcile()` and `reconcile*()` methods using go/ast to extract the ordered sequence of sub-resource reconciliation steps. For each step, captures:

- Method name and derived component (e.g., `ReconcileDatabase` -> `Database`)
- Conditional guards (if-blocks wrapping the call)
- Source location

This reveals the reconciliation order and which steps are conditional on configuration.

### Prometheus Metrics extractor

Regex-based extraction of Prometheus metric registrations. Matches `prometheus.New{Gauge,Counter,Histogram,Summary}(Vec)?` and `promauto.New*` patterns. For each metric:

- Composes the full metric name from Namespace + Subsystem + Name fields
- Extracts help text
- Captures labels from `[]string{...}` literals for Vec types
- Reports metric type and source location

### Status Conditions extractor

Parses Go const blocks using go/ast to extract status condition type and reason constants. Detects conditions by:

- Explicit type annotations (`ConditionType`, `StatusConditionType`)
- Name suffixes (`Available`, `Ready`, `Degraded`, `Progressing`, `Reconciled`)
- Name prefixes (`Condition*`, `Reason*`)

Associates reason constants with their preceding condition type within the same const block. Returns Go constant names for dedup with the operator config extractor.

### Platform Detection extractor

Three-pass detection of platform-specific behavior:

1. **Capability structs**: Finds struct types with boolean fields matching platform patterns (`IsOpenShift`, `HasRoute`, `HasIstio`, etc.) with a blocklist for common non-platform booleans (`IsDeleted`, `IsReady`, `HasError`)
2. **API discovery calls**: Matches `discovery.ServerResourcesForGroupVersion`, `RESTMapper().ResourcesFor`, and similar runtime capability checks
3. **Conditional resource creation**: Scans reconciler functions for if-blocks guarded by capability checks, extracting the resource kind and action from the body

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
