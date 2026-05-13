# Go AST Extraction

## Overview

Many Kubernetes operators `.gitignore` their generated YAML manifests (CRDs, webhooks, RBAC). The `opendatahub-operator` is a typical example: CRD YAML files are generated at build time and never committed to the repository. Standard YAML-based extractors find nothing in these repos.

Go AST extraction solves this by analyzing the Go source directly. Three extractors work together:

1. **Go CRD Extraction** (`go_crds.go`): finds CRD types from kubebuilder markers
2. **Webhook Behavioral Analysis** (`go_webhooks.go`): extracts field-level mutations and validations from webhook method bodies
3. **Programmatic Resource Operations** (`controller_watches.go`): detects `client.Create/Update/Patch/Delete` calls in reconcile methods

All three use `go/packages` for type-resolved loading. When `go/packages` fails (missing module dependencies, non-Go repos), they fall back to `go/parser` with reduced accuracy.

## CRD Extraction from Go Source

### What it extracts

For each CRD type found in Go source:

| Field | Source |
|-------|--------|
| Group | `SchemeBuilder`, `GroupVersion` var, or package path |
| Version | `GroupVersion` var or package directory name (e.g., `v1alpha1`) |
| Kind | Go type name with `+kubebuilder:object:root=true` |
| Scope | `+kubebuilder:resource:scope=Cluster` marker (default: Namespaced) |
| Storage version | `+kubebuilder:storageversion` marker |
| Hub/spoke | Conversion hub/spoke markers for multi-version CRDs |
| Field count | Recursive field count from struct definition |
| CEL rules | `+kubebuilder:validation:XValidation` rule expressions |

### How it finds CRD types

The extractor looks for Go struct types annotated with `+kubebuilder:object:root=true` in their doc comments. This is the standard kubebuilder marker that identifies a type as a CRD root object.

```go
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
type Widget struct {
    metav1.TypeMeta   `json:"..."`
    metav1.ObjectMeta `json:"metadata,omitempty"`
    Spec   WidgetSpec   `json:"spec,omitempty"`
    Status WidgetStatus `json:"status,omitempty"`
}
```

### GroupVersion resolution

The group and version are resolved by searching the same package for:

1. A `SchemeBuilder` registration with `GroupVersion`
2. A `GroupVersion = schema.GroupVersion{Group: "...", Version: "..."}` variable
3. Package path heuristic (last two segments as `group/version`)

## Webhook Behavioral Analysis

Standard webhook extraction tells you which resources a webhook intercepts and whether it's mutating or validating. Webhook behavioral analysis goes further: it reports *what* the webhook actually does to each field.

### Mutating webhooks

For each `Default()` method implementation, the extractor walks the AST to find:

- **Field assignments**: `w.Spec.Image = "default-image"` becomes a mutation on `spec.image`
- **Conditional mutations**: wrapping `if w.Spec.Image == ""` is captured as the condition
- **Helper method calls**: `r.setGPUDefaults()` on the same receiver is followed to find mutations inside

Example output:

```
Webhook /mutate-v1alpha1-widget (mutating) target=Widget
  MUTATES: spec.image (when w.Spec.Image == "")
  MUTATES: spec.gpu (via setGPUDefaults)
```

### Validating webhooks

For `ValidateCreate()`, `ValidateUpdate()`, and `ValidateDelete()` methods:

- **Field validations**: detects checks against specific fields
- **Validation classification**: categorizes as "invalid check", "required check", etc.
- **Helper method following**: same-receiver calls are traced for nested validations

Example output:

```
Webhook /validate-v1alpha1-widget (validating) target=Widget
  VALIDATES: spec.replicas (invalid check)
```

### Why this matters

Knowing that a webhook "intercepts Widget CREATE/UPDATE" is useful. Knowing that it "sets `spec.image` to a default when empty and validates that `spec.replicas` is positive" is actionable. This data feeds into security queries that check whether webhook validations are sufficient for the fields they protect.

## Programmatic Resource Operations

Controllers often create Kubernetes resources programmatically rather than through declarative manifests. A reconciler might construct a `Service` or `Deployment` in code and call `client.Create()` to apply it.

### Detection

The extractor scans reconcile methods for calls to:

- `client.Create(ctx, obj)`
- `client.Update(ctx, obj)`
- `client.Patch(ctx, obj, patch)`
- `client.Delete(ctx, obj)`

Also matches the `r.Client.Create()` and `r.client.Create()` patterns common in controller-runtime code.

### Type resolution

Using `go/packages` type information, the extractor resolves the concrete type of the object argument to determine the target Kind and API group. For example:

```go
svc := &corev1.Service{...}
r.Client.Create(ctx, svc)
```

Resolves to: `create Service` (group: `core/v1`).

### Output

Resource operations appear in the output as:

```json
{
  "resource_ops": [
    {"operation": "create", "kind": "Service", "group": "core/v1"},
    {"operation": "create", "kind": "Deployment", "group": "apps/v1"}
  ]
}
```

## Merge Strategy

Go AST extraction follows a YAML-authoritative merge strategy:

1. **YAML first**: if a CRD, webhook, or resource is found in YAML manifests, that data is authoritative
2. **Go supplements**: Go-extracted data fills gaps where YAML is absent
3. **No overrides**: Go data never overwrites YAML data for the same resource

Each CRD in the output carries a discovery badge indicating its source:

| Badge | Meaning |
|-------|---------|
| `YAML` | Discovered from YAML manifests only |
| `Go AST` | Discovered from Go source only (no YAML present) |
| `YAML + Go AST` | Found in both, YAML data used with Go supplementing |

The `go_source` field on each CRD indicates `go_ast` when the CRD was extracted from Go types.

## Security Hardening

Go AST extraction loads and analyzes code from untrusted repositories. The following hardening measures are in place:

| Measure | Purpose |
|---------|---------|
| `CGO_ENABLED=0` | Prevents native code execution during `go/packages` loading |
| `GOMODCACHE` isolation | Uses a temporary module cache directory, cleaned up after analysis |
| `GOPRIVATE` cleared | Prevents module loading from pulling from private registries |
| Symlink boundary checks | Prevents path traversal via symlinks that escape the repo root |
| `boundedFileSystem` | Kustomize file operations are confined to the repository directory |
| Checksum verification | `GONOSUMCHECK` is not set; module checksums are verified normally |

These measures ensure that analyzing a malicious repository cannot execute arbitrary code, access the host filesystem outside the repo, or leak credentials.

## Fallback Behavior

When `go/packages` loading fails (common reasons: missing Go toolchain, unresolvable dependencies, non-Go repo), the extractors degrade gracefully:

1. `go/packages` load is attempted first
2. On failure, extractors fall back to `go/parser` (AST-only, no type resolution)
3. `go/parser` extractors still find CRD types and webhook methods but cannot resolve cross-package types
4. If both fail, the extractors produce no output (not an error)

The `go_ast_mode` field in the output indicates the resolution level:

| Mode | Meaning |
|------|---------|
| `full` | `go/packages` loaded successfully, full type resolution available |
| `syntax` | Fell back to `go/parser`, AST-only analysis |
| (absent) | Go AST extraction was not attempted or produced no results |

## Example Output

### CRD from Go source

```json
{
  "crds": [
    {
      "group": "apps.example.com",
      "version": "v1alpha1",
      "kind": "Widget",
      "scope": "Namespaced",
      "field_count": 19,
      "go_source": "go_ast",
      "storage_version": true,
      "cel_rules": [
        "self.spec.replicas <= 100"
      ]
    }
  ]
}
```

### Webhook with behavioral analysis

```json
{
  "webhooks": [
    {
      "path": "/mutate-v1alpha1-widget",
      "type": "mutating",
      "target": "Widget",
      "mutations": [
        {
          "field": "spec.image",
          "condition": "w.Spec.Image == \"\""
        },
        {
          "field": "spec.gpu",
          "via": "setGPUDefaults"
        }
      ]
    },
    {
      "path": "/validate-v1alpha1-widget",
      "type": "validating",
      "target": "Widget",
      "validations": [
        {
          "field": "spec.replicas",
          "check": "invalid check"
        }
      ]
    }
  ]
}
```

### Programmatic resource operations

```json
{
  "resource_ops": [
    {"operation": "create", "kind": "Service", "group": "core/v1"},
    {"operation": "create", "kind": "Deployment", "group": "apps/v1"}
  ]
}
```

### GoASTMode in top-level output

```json
{
  "go_ast_mode": "full",
  "crds": [...],
  "webhooks": [...],
  "resource_ops": [...]
}
```
