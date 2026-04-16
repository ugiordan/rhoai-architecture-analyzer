// Package arch provides typed representations of architecture extraction data.
package arch

// Data holds parsed architecture information from the extractor JSON.
type Data struct {
	Component string      `json:"component"`
	CRDs      []CRD       `json:"crds"`
	RBAC      RBAC        `json:"rbac"`
	Webhooks  []Webhook   `json:"webhooks"`
	Secrets   []Secret    `json:"secrets_referenced"`
	Cache     CacheConfig `json:"cache_config"`
}

// CRD represents a CustomResourceDefinition from the architecture extraction.
type CRD struct {
	Group           string   `json:"group"`
	Version         string   `json:"version"`
	Kind            string   `json:"kind"`
	Scope           string   `json:"scope"`
	FieldsCount     int      `json:"fields_count"`
	ValidationRules []string `json:"validation_rules"`
	Source          string   `json:"source"`
}

// RBAC holds role-based access control information.
type RBAC struct {
	ClusterRoles       []ClusterRole       `json:"cluster_roles"`
	KubebuilderMarkers []KubebuilderMarker `json:"kubebuilder_markers"`
}

// ClusterRole represents a Kubernetes ClusterRole.
type ClusterRole struct {
	Name   string     `json:"name"`
	Source string     `json:"source"`
	Rules  []RBACRule `json:"rules"`
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
	Parsed map[string]interface{} `json:"parsed"`
}

// Webhook represents a Kubernetes admission webhook.
type Webhook struct {
	Name          string        `json:"name"`
	Type          string        `json:"type"`
	ServiceRef    string        `json:"service_ref,omitempty"`
	Path          string        `json:"path,omitempty"`
	FailurePolicy string        `json:"failure_policy,omitempty"`
	Rules         []WebhookRule `json:"rules"`
	Source        string        `json:"source"`
}

// WebhookRule defines what resources a webhook intercepts.
// JSON tags use camelCase to match Kubernetes convention and extractor output.
type WebhookRule struct {
	Operations  []string `json:"operations"`
	APIGroups   []string `json:"apiGroups"`
	APIVersions []string `json:"apiVersions"`
	Resources   []string `json:"resources"`
}

// Secret represents a referenced Kubernetes secret.
type Secret struct {
	Name          string   `json:"name"`
	Type          string   `json:"type"`
	ReferencedBy  []string `json:"referenced_by"`
	ProvisionedBy string   `json:"provisioned_by"`
}

// CacheConfig holds controller-runtime cache configuration details.
type CacheConfig struct {
	FilteredTypes []CacheFilteredType `json:"filtered_types"`
	Issues        []string            `json:"issues"`
}

// CacheFilteredType is a type with a cache filter (label, field, or namespace).
type CacheFilteredType struct {
	Type       string `json:"type"`
	FilterKind string `json:"filter_kind"`
	Filter     string `json:"filter"`
}
