// Package arch provides typed representations of architecture extraction data.
package arch

import "encoding/json"

// Data holds parsed architecture information from the extractor JSON.
type Data struct {
	Component string      `json:"component"`
	CRDs      []CRD       `json:"crds,omitempty"`
	RBAC      RBAC        `json:"rbac"`
	Webhooks  []Webhook   `json:"webhooks,omitempty"`
	Secrets   []Secret    `json:"secrets_referenced,omitempty"`
	Cache        CacheConfig  `json:"cache_config"`
	FeatureGates []FeatureGate `json:"feature_gates,omitempty"`

	// Cross-domain fields: deployment context for security posture queries.
	NetworkPolicies     []NetworkPolicy      `json:"network_policies,omitempty"`
	HTTPEndpoints       []HTTPEndpoint       `json:"http_endpoints,omitempty"`
	ExternalConnections []ExternalConnection `json:"external_connections,omitempty"`
	Deployments         []Deployment         `json:"deployments,omitempty"`
}

// CRD represents a CustomResourceDefinition from the architecture extraction.
type CRD struct {
	Group           string   `json:"group"`
	Version         string   `json:"version"`
	Kind            string   `json:"kind"`
	Scope           string   `json:"scope,omitempty"`
	FieldsCount     int      `json:"fields_count,omitempty"`
	ValidationRules []string `json:"validation_rules,omitempty"`
	Source          string   `json:"source"`
}

// RBAC holds role-based access control information.
type RBAC struct {
	ClusterRoles       []ClusterRole       `json:"cluster_roles,omitempty"`
	KubebuilderMarkers []KubebuilderMarker `json:"kubebuilder_markers,omitempty"`
}

// ClusterRole represents a Kubernetes ClusterRole.
type ClusterRole struct {
	Name   string     `json:"name"`
	Source string     `json:"source"`
	Rules  []RBACRule `json:"rules,omitempty"`
}

// RBACRule represents a single RBAC policy rule.
// JSON tags use camelCase to match Kubernetes RBAC convention and extractor output.
type RBACRule struct {
	APIGroups     []string `json:"apiGroups"`
	Resources     []string `json:"resources"`
	Verbs         []string `json:"verbs"`
	ResourceNames []string `json:"resourceNames,omitempty"`
}

// KubebuilderMarker represents a kubebuilder RBAC marker annotation found in Go source.
// Matches extractor.RBACMarker field layout.
type KubebuilderMarker struct {
	File   string                 `json:"file"`
	Line   int                    `json:"line"`
	Marker string                 `json:"marker"`
	Parsed map[string]interface{} `json:"parsed,omitempty"`
}

// Webhook represents a Kubernetes admission webhook.
type Webhook struct {
	Name          string        `json:"name"`
	Type          string        `json:"type"`
	ServiceRef    string        `json:"service_ref,omitempty"`
	Path          string        `json:"path,omitempty"`
	FailurePolicy string        `json:"failure_policy,omitempty"`
	Rules         []WebhookRule `json:"rules,omitempty"`
	Source        string        `json:"source"`
}

// WebhookRule defines what resources a webhook intercepts.
// JSON tags use camelCase to match Kubernetes convention and extractor output.
type WebhookRule struct {
	Operations  []string `json:"operations,omitempty"`
	APIGroups   []string `json:"apiGroups,omitempty"`
	APIVersions []string `json:"apiVersions,omitempty"`
	Resources   []string `json:"resources,omitempty"`
}

// Secret represents a referenced Kubernetes secret.
type Secret struct {
	Name          string   `json:"name"`
	Type          string   `json:"type,omitempty"`
	ReferencedBy  []string `json:"referenced_by,omitempty"`
	ProvisionedBy string   `json:"provisioned_by,omitempty"`
}

// CacheConfig holds controller-runtime cache configuration details.
type CacheConfig struct {
	FilteredTypes []CacheFilteredType `json:"filtered_types,omitempty"`
	Issues        []string            `json:"issues,omitempty"`
}

// FeatureGate represents a feature gate definition from the architecture extraction.
type FeatureGate struct {
	Name          string `json:"name"`
	Default       bool   `json:"default"`
	PreRelease    string `json:"pre_release,omitempty"`
	LockToDefault bool   `json:"lock_to_default,omitempty"`
	Source        string `json:"source"`
	RuntimeSet    bool   `json:"runtime_set,omitempty"`
}

// CacheFilteredType is a type with a cache filter (label, field, or namespace).
type CacheFilteredType struct {
	Type       string `json:"type"`
	FilterKind string `json:"filter_kind"`
	Filter     string `json:"filter"`
}

// NetworkPolicy represents a Kubernetes NetworkPolicy from the architecture extraction.
// IngressRules and EgressRules are raw JSON arrays; queries only need len() > 0.
// PodSelector maps to the extractor's pod_selector field for per-pod policy matching.
type NetworkPolicy struct {
	Name         string                 `json:"name"`
	Source       string                 `json:"source"`
	PodSelector  map[string]interface{} `json:"pod_selector,omitempty"`
	PolicyTypes  []string               `json:"policy_types,omitempty"`
	IngressRules []json.RawMessage      `json:"ingress_rules,omitempty"`
	EgressRules  []json.RawMessage      `json:"egress_rules,omitempty"`
}

// HTTPEndpoint represents an HTTP route registration found by the extractor.
type HTTPEndpoint struct {
	Method  string `json:"method,omitempty"`
	Path    string `json:"path"`
	Handler string `json:"handler,omitempty"`
	Source  string `json:"source"`
}

// ExternalConnection represents a reference to an external service found in source code.
type ExternalConnection struct {
	Type     string `json:"type"`
	Service  string `json:"service"`
	Target   string `json:"target"`
	Source   string `json:"source"`
	Function string `json:"function,omitempty"`
}

// Deployment represents a Deployment or StatefulSet from the architecture extraction.
type Deployment struct {
	Name           string `json:"name"`
	Kind           string `json:"kind"`
	Source         string `json:"source"`
	ServiceAccount string `json:"service_account,omitempty"`
}
