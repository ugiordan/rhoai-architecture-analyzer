package extractor

// ExtractOptions configures the extraction process.
type ExtractOptions struct {
	// Org overrides the auto-detected GitHub organization.
	Org string
	// ModulePrefixes lists Go module prefixes considered "internal" dependencies.
	// Defaults to ["github.com/opendatahub-io/", "github.com/red-hat-data-services/"].
	ModulePrefixes []string
}

// DefaultModulePrefixes returns the standard internal module prefixes for RHOAI.
func DefaultModulePrefixes() []string {
	return []string{
		"github.com/opendatahub-io/",
		"github.com/red-hat-data-services/",
	}
}

// ComponentArchitecture is the top-level extraction result.
type ComponentArchitecture struct {
	Component       string             `json:"component"`
	Repo            string             `json:"repo"`
	ExtractedAt     string             `json:"extracted_at"`
	AnalyzerVersion string             `json:"analyzer_version"`
	CRDs            []CRD              `json:"crds"`
	RBAC            *RBACData          `json:"rbac"`
	Services        []Service          `json:"services"`
	Deployments     []Deployment       `json:"deployments"`
	NetworkPolicies []NetworkPolicy    `json:"network_policies"`
	ControllerWatch []ControllerWatch  `json:"controller_watches"`
	Dependencies    *DependencyData    `json:"dependencies"`
	Secrets         []SecretRef        `json:"secrets_referenced"`
	Dockerfiles     []DockerfileInfo   `json:"dockerfiles"`
	Helm            *HelmData          `json:"helm"`
	Webhooks        []WebhookConfig    `json:"webhooks"`
	ConfigMaps      []ConfigMapRef     `json:"configmaps"`
	HTTPEndpoints   []HTTPEndpoint     `json:"http_endpoints"`
	IngressRouting      []IngressResource  `json:"ingress_routing"`
	ExternalConnections []ExternalConnection `json:"external_connections,omitempty"`
	FeatureGates        []FeatureGate        `json:"feature_gates,omitempty"`
	CacheConfig         *CacheConfig       `json:"cache_config,omitempty"`
}

// CRD represents a single version of a CustomResourceDefinition.
type CRD struct {
	Group           string   `json:"group"`
	Version         string   `json:"version"`
	Kind            string   `json:"kind"`
	Scope           string   `json:"scope"`
	FieldsCount     int      `json:"fields_count"`
	ValidationRules []string `json:"validation_rules"`
	Source          string   `json:"source"`
}

// RBACData holds all RBAC-related extractions.
type RBACData struct {
	ClusterRoles        []RBACRole        `json:"cluster_roles"`
	ClusterRoleBindings []RBACBinding     `json:"cluster_role_bindings"`
	Roles               []RBACRole        `json:"roles"`
	RoleBindings        []RBACBinding     `json:"role_bindings"`
	KubebuilderMarkers  []RBACMarker      `json:"kubebuilder_markers"`
}

// RBACRole is a ClusterRole or Role with its rules.
type RBACRole struct {
	Name   string     `json:"name"`
	Source string     `json:"source"`
	Rules  []RBACRule `json:"rules"`
}

// RBACRule is a single RBAC policy rule.
type RBACRule struct {
	APIGroups     []string `json:"apiGroups"`
	Resources     []string `json:"resources"`
	Verbs         []string `json:"verbs"`
	ResourceNames []string `json:"resourceNames"`
}

// RBACBinding is a ClusterRoleBinding or RoleBinding.
type RBACBinding struct {
	Name     string        `json:"name"`
	RoleRef  string        `json:"role_ref"`
	Subjects []RBACSubject `json:"subjects"`
	Source   string        `json:"source"`
}

// RBACSubject is a subject in a role binding.
type RBACSubject struct {
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// RBACMarker is a kubebuilder RBAC annotation found in Go source.
type RBACMarker struct {
	File   string                 `json:"file"`
	Line   int                    `json:"line"`
	Marker string                 `json:"marker"`
	Parsed map[string]interface{} `json:"parsed"`
}

// Service represents a Kubernetes Service definition.
type Service struct {
	Name     string                 `json:"name"`
	Source   string                 `json:"source"`
	Type     string                 `json:"type"`
	Ports    []ServicePort          `json:"ports"`
	Selector map[string]interface{} `json:"selector"`
}

// ServicePort is a single port entry in a Service.
type ServicePort struct {
	Name       string      `json:"name"`
	Port       interface{} `json:"port"`
	TargetPort interface{} `json:"targetPort"`
	Protocol   string      `json:"protocol"`
}

// Deployment represents a Deployment or StatefulSet.
type Deployment struct {
	Name                         string      `json:"name"`
	Kind                         string      `json:"kind"`
	Source                       string      `json:"source"`
	Replicas                     interface{} `json:"replicas"`
	ServiceAccount               string      `json:"service_account"`
	AutomountServiceAccountToken interface{} `json:"automount_service_account_token"`
	Containers                   []Container `json:"containers"`
}

// Container is a container spec within a Deployment.
type Container struct {
	Name              string                   `json:"name"`
	Image             string                   `json:"image"`
	Ports             []ContainerPort          `json:"ports"`
	SecurityContext   map[string]interface{}   `json:"security_context"`
	EnvFromSecrets    []string                 `json:"env_from_secrets"`
	EnvFromConfigmaps []string                 `json:"env_from_configmaps"`
	VolumeMounts      []map[string]interface{} `json:"volume_mounts"`
	Resources         map[string]interface{}   `json:"resources"`
	EnvVars           map[string]string        `json:"env_vars,omitempty"`
	LivenessProbe     *ProbeInfo               `json:"liveness_probe,omitempty"`
	ReadinessProbe    *ProbeInfo               `json:"readiness_probe,omitempty"`
	StartupProbe      *ProbeInfo               `json:"startup_probe,omitempty"`
}

// ProbeInfo holds liveness/readiness/startup probe metadata.
type ProbeInfo struct {
	Type string      `json:"type"`
	Path string      `json:"path,omitempty"`
	Port interface{} `json:"port,omitempty"`
}

// ContainerPort is a port exposed by a container.
type ContainerPort struct {
	Name          string `json:"name"`
	ContainerPort int    `json:"containerPort"`
	Protocol      string `json:"protocol"`
}

// NetworkPolicy represents a Kubernetes NetworkPolicy.
type NetworkPolicy struct {
	Name         string                   `json:"name"`
	Source       string                   `json:"source"`
	PodSelector  map[string]interface{}   `json:"pod_selector"`
	PolicyTypes  []string                 `json:"policy_types"`
	IngressRules []map[string]interface{} `json:"ingress_rules"`
	EgressRules  []map[string]interface{} `json:"egress_rules"`
}

// ControllerWatch represents a For/Owns/Watches call in controller code.
type ControllerWatch struct {
	Type   string `json:"type"`
	GVK    string `json:"gvk"`
	Source string `json:"source"`
}

// DependencyData holds Go module dependencies.
type DependencyData struct {
	GoVersion         string              `json:"go_version,omitempty"`
	Toolchain         string              `json:"toolchain,omitempty"`
	GoModules         []GoModule          `json:"go_modules"`
	ReplaceDirectives []ReplaceDirective  `json:"replace_directives,omitempty"`
	InternalODH       []InternalODH       `json:"internal_odh"`
}

// ReplaceDirective represents a go.mod replace directive.
type ReplaceDirective struct {
	Original    string `json:"original"`
	Replacement string `json:"replacement"`
	Version     string `json:"version"`
}

// GoModule is a single Go module dependency.
type GoModule struct {
	Module  string `json:"module"`
	Version string `json:"version"`
}

// InternalODH is an internal OpenDataHub dependency.
type InternalODH struct {
	Component   string `json:"component"`
	Interaction string `json:"interaction"`
}

// SecretRef is a reference to a Kubernetes Secret (name only, never values).
type SecretRef struct {
	Name           string   `json:"name"`
	Type           string   `json:"type"`
	ReferencedBy   []string `json:"referenced_by"`
	ProvisionedBy  string   `json:"provisioned_by"`
}

// DockerfileInfo holds metadata extracted from a Dockerfile.
type DockerfileInfo struct {
	Path           string            `json:"path"`
	BaseImage      string            `json:"base_image"`
	Stages         int               `json:"stages"`
	User           string            `json:"user"`
	ExposedPorts   []int             `json:"exposed_ports"`
	Issues         []string          `json:"issues"`
	Architectures  []string          `json:"architectures,omitempty"`
	FIPSEnabled    bool              `json:"fips_enabled,omitempty"`
	BuildArgs      map[string]string `json:"build_args,omitempty"`
}

// HelmData holds Helm chart metadata and security-relevant value defaults.
type HelmData struct {
	ChartName      string                 `json:"chart_name,omitempty"`
	ChartVersion   string                 `json:"chart_version,omitempty"`
	ValuesDefaults map[string]interface{} `json:"values_defaults,omitempty"`
}

// WebhookConfig represents a Kubernetes webhook configuration.
type WebhookConfig struct {
	Name          string        `json:"name"`
	Type          string        `json:"type"` // "validating" or "mutating"
	ServiceRef    string        `json:"service_ref,omitempty"`
	Path          string        `json:"path,omitempty"`
	FailurePolicy string       `json:"failure_policy,omitempty"`
	Rules         []WebhookRule `json:"rules"`
	Source        string        `json:"source"`
}

// WebhookRule defines what resources a webhook intercepts.
type WebhookRule struct {
	APIGroups   []string `json:"apiGroups"`
	APIVersions []string `json:"apiVersions"`
	Resources   []string `json:"resources"`
	Operations  []string `json:"operations"`
}

// ConfigMapRef represents a ConfigMap definition with its data keys.
type ConfigMapRef struct {
	Name         string   `json:"name"`
	DataKeys     []string `json:"data_keys"`
	ReferencedBy []string `json:"referenced_by"`
	Source       string   `json:"source"`
}

// HTTPEndpoint represents an HTTP route registration found in source code.
type HTTPEndpoint struct {
	Method  string `json:"method"`
	Path    string `json:"path"`
	Handler string `json:"handler,omitempty"`
	Source  string `json:"source"`
}

// CacheConfig represents the controller-runtime cache configuration extracted
// from Go source, used to detect OOM risks from unfiltered informers.
type CacheConfig struct {
	ManagerFile       string              `json:"manager_file"`
	CacheScope        string              `json:"cache_scope"`
	FilteredTypes     []CacheFilteredType `json:"filtered_types"`
	TransformTypes    []string            `json:"transform_types"`
	DefaultTransform  bool                `json:"default_transform"`
	DisabledTypes     []string            `json:"disabled_types"`
	GoMemLimit        string              `json:"gomemlimit"`
	MemoryLimit       string              `json:"memory_limit"`
	ImplicitInformers []ImplicitInformer  `json:"implicit_informers"`
	Issues            []string            `json:"issues"`
}

// CacheFilteredType is a type with a cache filter (label, field, or namespace).
type CacheFilteredType struct {
	Type       string `json:"type"`
	FilterKind string `json:"filter_kind"`
	Filter     string `json:"filter"`
}

// ImplicitInformer is a client.Get call that silently creates a cluster-wide
// informer for a type not in the watch set or DisableFor list.
type ImplicitInformer struct {
	Type   string `json:"type"`
	Source string `json:"source"`
	Risk   string `json:"risk"`
}

// ExternalConnection represents a reference to an external service found in Go source.
type ExternalConnection struct {
	Type     string `json:"type"`               // database, object-storage, grpc, messaging, api
	Service  string `json:"service"`             // postgres, mysql, redis, s3, kafka, etc.
	Target   string `json:"target"`              // redacted connection target
	Source   string `json:"source"`              // file:line
	Function string `json:"function,omitempty"`   // enclosing function name
}

// FeatureGate represents a feature gate definition found in Go source.
type FeatureGate struct {
	Name          string `json:"name"`                      // feature gate name (string value)
	Default       bool   `json:"default"`                   // default enabled state
	PreRelease    string `json:"pre_release,omitempty"`     // Alpha, Beta, GA, Deprecated
	LockToDefault bool   `json:"lock_to_default,omitempty"` // locked to default value
	Source        string `json:"source"`                    // file:line of registration
	RuntimeSet    bool   `json:"runtime_set,omitempty"`     // set via Set() at runtime rather than Add()
}

// IngressResource represents a Gateway API, Istio, or Kubernetes ingress resource.
type IngressResource struct {
	Kind    string   `json:"kind"`
	Name    string   `json:"name"`
	Hosts   []string `json:"hosts,omitempty"`
	Paths   []string `json:"paths,omitempty"`
	Backend string   `json:"backend,omitempty"`
	TLS     bool     `json:"tls"`
	Source  string   `json:"source"`
}
