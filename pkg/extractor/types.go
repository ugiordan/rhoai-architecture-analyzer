package extractor

// ExtractOptions configures the extraction process.
type ExtractOptions struct {
	// Org overrides the auto-detected GitHub organization.
	Org string
	// ModulePrefixes lists Go module prefixes considered "internal" dependencies.
	// Defaults to ["github.com/opendatahub-io/", "github.com/red-hat-data-services/"].
	ModulePrefixes []string
	// OverlayPreference lists kustomize overlay directory names in priority order.
	// Defaults to DefaultPreferredOverlays if empty.
	OverlayPreference []string
}

// DefaultModulePrefixes returns the standard internal module prefixes for the analyzed platform.
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
	CommitSHA       string             `json:"commit_sha,omitempty"`
	ExtractedAt     string             `json:"extracted_at"`
	AnalyzerVersion string             `json:"analyzer_version"`
	CRDs            []CRD              `json:"crds,omitempty"`
	RBAC            *RBACData          `json:"rbac,omitempty"`
	Services        []Service          `json:"services,omitempty"`
	Deployments     []Deployment       `json:"deployments,omitempty"`
	NetworkPolicies []NetworkPolicy    `json:"network_policies,omitempty"`
	ControllerWatch []ControllerWatch  `json:"controller_watches,omitempty"`
	Dependencies    *DependencyData    `json:"dependencies,omitempty"`
	Secrets         []SecretRef        `json:"secrets_referenced,omitempty"`
	Dockerfiles     []DockerfileInfo   `json:"dockerfiles,omitempty"`
	Helm            *HelmData          `json:"helm,omitempty"`
	Webhooks        []WebhookConfig    `json:"webhooks,omitempty"`
	ConfigMaps      []ConfigMapRef     `json:"configmaps,omitempty"`
	HTTPEndpoints   []HTTPEndpoint     `json:"http_endpoints,omitempty"`
	IngressRouting      []IngressResource  `json:"ingress_routing,omitempty"`
	ExternalConnections []ExternalConnection `json:"external_connections,omitempty"`
	FeatureGates        []FeatureGate        `json:"feature_gates,omitempty"`
	CacheConfig         *CacheConfig       `json:"cache_config,omitempty"`
	KustomizeComponents []KustomizeComponent `json:"kustomize_components,omitempty"`
	ServingRuntimes     []ServingRuntime     `json:"serving_runtimes,omitempty"`
	ResourceDefaults    []ResourceDefault    `json:"resource_defaults,omitempty"`
	PodDisruptionBudgets    []PodDisruptionBudget    `json:"pod_disruption_budgets,omitempty"`
	HorizontalPodAutoscalers []HorizontalPodAutoscaler `json:"horizontal_pod_autoscalers,omitempty"`
	APITypes                []APITypeDefinition        `json:"api_types,omitempty"`
	OperatorConfig       []OperatorConstant     `json:"operator_config,omitempty"`
	ReconcileSequences   []ReconcileSequence    `json:"reconcile_sequences,omitempty"`
	PrometheusMetrics    []PrometheusMetric     `json:"prometheus_metrics,omitempty"`
	StatusConditions     []StatusCondition      `json:"status_conditions,omitempty"`
	PlatformDetection    *PlatformDetection     `json:"platform_detection,omitempty"`
	TemplateFiles           []TemplateFile             `json:"template_files,omitempty"`
	DataCoverage            map[string]string          `json:"data_coverage,omitempty"`
	Summary                 string                    `json:"summary,omitempty"`
}

// CRD represents a CustomResourceDefinition with all its versions.
type CRD struct {
	Group           string       `json:"group"`
	Version         string       `json:"version"`
	Kind            string       `json:"kind"`
	Scope           string       `json:"scope,omitempty"`
	Versions        []CRDVersion `json:"versions,omitempty"`
	FieldsCount     int          `json:"fields_count,omitempty"`
	ValidationRules []string     `json:"validation_rules,omitempty"`
	Source          string       `json:"source"`
}

// CRDVersion represents a single served version of a CRD.
type CRDVersion struct {
	Name    string `json:"name"`
	Served  bool   `json:"served"`
	Storage bool   `json:"storage"`
}

// RBACData holds all RBAC-related extractions.
type RBACData struct {
	ClusterRoles        []RBACRole        `json:"cluster_roles,omitempty"`
	ClusterRoleBindings []RBACBinding     `json:"cluster_role_bindings,omitempty"`
	Roles               []RBACRole        `json:"roles,omitempty"`
	RoleBindings        []RBACBinding     `json:"role_bindings,omitempty"`
	KubebuilderMarkers  []RBACMarker      `json:"kubebuilder_markers,omitempty"`
}

// RBACRole is a ClusterRole or Role with its rules.
type RBACRole struct {
	Name            string            `json:"name"`
	Source          string            `json:"source"`
	Rules           []RBACRule        `json:"rules,omitempty"`
	AggregationRule map[string]string `json:"aggregation_rule,omitempty"`
}

// RBACRule is a single RBAC policy rule.
type RBACRule struct {
	APIGroups     []string `json:"apiGroups"`
	Resources     []string `json:"resources"`
	Verbs         []string `json:"verbs"`
	ResourceNames []string `json:"resourceNames,omitempty"`
}

// RBACBinding is a ClusterRoleBinding or RoleBinding.
type RBACBinding struct {
	Name     string        `json:"name"`
	RoleRef  string        `json:"role_ref,omitempty"`
	Subjects []RBACSubject `json:"subjects,omitempty"`
	Source   string        `json:"source"`
}

// RBACSubject is a subject in a role binding.
type RBACSubject struct {
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
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
	Name             string                 `json:"name"`
	Source           string                 `json:"source"`
	Type             string                 `json:"type,omitempty"`
	Ports            []ServicePort          `json:"ports,omitempty"`
	Selector         map[string]interface{} `json:"selector,omitempty"`
	TargetDeployment string                 `json:"target_deployment,omitempty"`
}

// ServicePort is a single port entry in a Service.
type ServicePort struct {
	Name       string      `json:"name,omitempty"`
	Port       interface{} `json:"port"`
	TargetPort interface{} `json:"targetPort,omitempty"`
	Protocol   string      `json:"protocol,omitempty"`
	Condition  string      `json:"condition,omitempty"`
}

// Deployment represents a Deployment or StatefulSet.
type Deployment struct {
	Name                         string      `json:"name"`
	Kind                         string      `json:"kind"`
	Source                       string      `json:"source"`
	Replicas                     interface{} `json:"replicas,omitempty"`
	ServiceAccount               string      `json:"service_account,omitempty"`
	AutomountServiceAccountToken interface{} `json:"automount_service_account_token,omitempty"`
	Containers                   []Container `json:"containers,omitempty"`
	InitContainers               []Container `json:"init_containers,omitempty"`
	Issues                       []string    `json:"issues,omitempty"`
}

// Container is a container spec within a Deployment.
type Container struct {
	Name              string                   `json:"name"`
	Image             string                   `json:"image"`
	Ports             []ContainerPort          `json:"ports,omitempty"`
	SecurityContext   map[string]interface{}   `json:"security_context,omitempty"`
	EnvFromSecrets    []string                 `json:"env_from_secrets,omitempty"`
	EnvFromConfigmaps []string                 `json:"env_from_configmaps,omitempty"`
	VolumeMounts      []map[string]interface{} `json:"volume_mounts,omitempty"`
	Resources         map[string]interface{}   `json:"resources,omitempty"`
	EnvVars           map[string]string        `json:"env_vars,omitempty"`
	LivenessProbe     *ProbeInfo               `json:"liveness_probe,omitempty"`
	ReadinessProbe    *ProbeInfo               `json:"readiness_probe,omitempty"`
	StartupProbe      *ProbeInfo               `json:"startup_probe,omitempty"`
	Issues            []string                 `json:"issues,omitempty"`
}

// ProbeInfo holds liveness/readiness/startup probe metadata.
type ProbeInfo struct {
	Type string      `json:"type"`
	Path string      `json:"path,omitempty"`
	Port interface{} `json:"port,omitempty"`
}

// ContainerPort is a port exposed by a container.
type ContainerPort struct {
	Name          string `json:"name,omitempty"`
	ContainerPort int    `json:"containerPort"`
	Protocol      string `json:"protocol,omitempty"`
}

// NetworkPolicy represents a Kubernetes NetworkPolicy.
type NetworkPolicy struct {
	Name         string                   `json:"name"`
	Source       string                   `json:"source"`
	PodSelector  map[string]interface{}   `json:"pod_selector,omitempty"`
	PolicyTypes  []string                 `json:"policy_types,omitempty"`
	IngressRules []map[string]interface{} `json:"ingress_rules,omitempty"`
	EgressRules  []map[string]interface{} `json:"egress_rules,omitempty"`
	Issues       []string                 `json:"issues,omitempty"`
}

// ControllerWatch represents a For/Owns/Watches call in controller code.
type ControllerWatch struct {
	Type       string `json:"type"`
	GVK        string `json:"gvk"`
	Controller string `json:"controller,omitempty"`
	Source     string `json:"source"`
}

// DependencyData holds Go module dependencies.
type DependencyData struct {
	GoVersion         string              `json:"go_version,omitempty"`
	Toolchain         string              `json:"toolchain,omitempty"`
	GoModules         []GoModule          `json:"go_modules,omitempty"`
	ReplaceDirectives []ReplaceDirective  `json:"replace_directives,omitempty"`
	InternalODH       []InternalODH       `json:"internal_odh,omitempty"`
	Issues            []string            `json:"issues,omitempty"`
}

// ReplaceDirective represents a go.mod replace directive.
type ReplaceDirective struct {
	Original    string `json:"original"`
	Replacement string `json:"replacement"`
	Version     string `json:"version,omitempty"`
}

// GoModule is a single Go module dependency.
type GoModule struct {
	Module   string `json:"module"`
	Version  string `json:"version"`
	Category string `json:"category,omitempty"`
	Purpose  string `json:"purpose,omitempty"`
}

// InternalODH is an internal OpenDataHub dependency.
type InternalODH struct {
	Component   string `json:"component"`
	Interaction string `json:"interaction"`
}

// SecretRef is a reference to a Kubernetes Secret (name only, never values).
type SecretRef struct {
	Name           string   `json:"name"`
	Type           string   `json:"type,omitempty"`
	ReferencedBy   []string `json:"referenced_by,omitempty"`
	ProvisionedBy  string   `json:"provisioned_by,omitempty"`
}

// DockerfileInfo holds metadata extracted from a Dockerfile.
type DockerfileInfo struct {
	Path             string            `json:"path"`
	BaseImage        string            `json:"base_image"`
	BuildStageImages []string          `json:"build_stage_images,omitempty"`
	Stages           int               `json:"stages,omitempty"`
	User             string            `json:"user,omitempty"`
	ExposedPorts     []int             `json:"exposed_ports,omitempty"`
	Issues           []string          `json:"issues,omitempty"`
	Architectures    []string          `json:"architectures,omitempty"`
	FIPSEnabled      bool              `json:"fips_enabled,omitempty"`
	BuildArgs        map[string]string `json:"build_args,omitempty"`
}

// HelmData holds Helm chart metadata and security-relevant value defaults.
type HelmData struct {
	ChartName      string                 `json:"chart_name,omitempty"`
	ChartVersion   string                 `json:"chart_version,omitempty"`
	ValuesDefaults map[string]interface{} `json:"values_defaults,omitempty"`
}


// WebhookConfig represents a Kubernetes webhook configuration.
type WebhookConfig struct {
	Name           string        `json:"name"`
	Type           string        `json:"type"` // "validating" or "mutating"
	ServiceRef     string        `json:"service_ref,omitempty"`
	Path           string        `json:"path,omitempty"`
	Port           int           `json:"port,omitempty"`
	FailurePolicy  string        `json:"failure_policy,omitempty"`
	SideEffects    string        `json:"side_effects,omitempty"`
	TimeoutSeconds int           `json:"timeout_seconds,omitempty"`
	Rules          []WebhookRule `json:"rules,omitempty"`
	Source         string        `json:"source"`
}

// WebhookRule defines what resources a webhook intercepts.
type WebhookRule struct {
	APIGroups   []string `json:"apiGroups,omitempty"`
	APIVersions []string `json:"apiVersions,omitempty"`
	Resources   []string `json:"resources,omitempty"`
	Operations  []string `json:"operations,omitempty"`
}

// ConfigMapRef represents a ConfigMap definition with its data keys.
type ConfigMapRef struct {
	Name         string   `json:"name"`
	DataKeys     []string `json:"data_keys,omitempty"`
	ReferencedBy []string `json:"referenced_by,omitempty"`
	Source       string   `json:"source"`
}

// HTTPEndpoint represents an HTTP route registration found in source code.
type HTTPEndpoint struct {
	Method  string `json:"method,omitempty"`
	Path    string `json:"path"`
	Handler string `json:"handler,omitempty"`
	Source  string `json:"source"`
}

// CacheConfig represents the controller-runtime cache configuration extracted
// from Go source, used to detect OOM risks from unfiltered informers.
type CacheConfig struct {
	ManagerFile       string              `json:"manager_file"`
	CacheScope        string              `json:"cache_scope,omitempty"`
	FilteredTypes     []CacheFilteredType `json:"filtered_types,omitempty"`
	TransformTypes    []CacheTransform    `json:"transform_types,omitempty"`
	DefaultTransform  bool                `json:"default_transform,omitempty"`
	DefaultTransformFunc string           `json:"default_transform_func,omitempty"`
	DisabledTypes     []string            `json:"disabled_types,omitempty"`
	GoMemLimit        string              `json:"gomemlimit,omitempty"`
	MemoryLimit       string              `json:"memory_limit,omitempty"`
	GoMemLimitRatio   float64             `json:"gomemlimit_ratio,omitempty"`
	ImplicitInformers []ImplicitInformer  `json:"implicit_informers,omitempty"`
	Issues            []string            `json:"issues,omitempty"`
}

// CacheTransform is a type with a custom cache transform function.
type CacheTransform struct {
	Type     string `json:"type"`
	Function string `json:"function"`
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
	Type              string   `json:"type"`                         // database, object-storage, grpc, messaging, api
	Service           string   `json:"service"`                      // postgres, mysql, redis, s3, kafka, etc.
	Target            string   `json:"target"`                       // redacted connection target
	Source            string   `json:"source"`                       // file:line
	Function          string   `json:"function,omitempty"`           // enclosing function name
	CredentialSources []string `json:"credential_sources,omitempty"` // secrets/configmaps providing credentials
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

// ServingRuntime represents a KServe/ModelMesh serving runtime definition.
type ServingRuntime struct {
	Name             string                `json:"name"`
	Kind             string                `json:"kind"` // ServingRuntime or ClusterServingRuntime
	Containers       []ServingContainer    `json:"containers,omitempty"`
	SupportedFormats []SupportedModelFormat `json:"supported_formats,omitempty"`
	MultiModel       bool                  `json:"multi_model,omitempty"`
	Disabled         bool                  `json:"disabled,omitempty"`
	Source           string                `json:"source"`
}

// ServingContainer is a container within a serving runtime.
type ServingContainer struct {
	Name      string                 `json:"name"`
	Image     string                 `json:"image,omitempty"`
	Resources map[string]interface{} `json:"resources,omitempty"`
	Ports     []ContainerPort        `json:"ports,omitempty"`
	Args      []string               `json:"args,omitempty"`
}

// SupportedModelFormat describes a model format supported by a serving runtime.
type SupportedModelFormat struct {
	Name       string `json:"name"`
	Version    string `json:"version,omitempty"`
	AutoSelect bool   `json:"auto_select,omitempty"`
	Priority   int    `json:"priority,omitempty"`
}

// ResourceDefault represents a resource default value extracted from configmaps.
type ResourceDefault struct {
	Component string                 `json:"component"`
	Key       string                 `json:"key"`
	Values    map[string]interface{} `json:"values,omitempty"`
	Source    string                 `json:"source"`
}

// IngressResource represents a Gateway API, Istio, or Kubernetes ingress resource.
// Resources can be extracted from YAML manifests or inferred from RBAC permissions.
type IngressResource struct {
	Kind      string   `json:"kind"`
	Name      string   `json:"name"`
	Hosts     []string `json:"hosts,omitempty"`
	Paths     []string `json:"paths,omitempty"`
	Backend   string   `json:"backend,omitempty"`
	TLS       bool     `json:"tls"`
	Source    string   `json:"source"`
	Note      string   `json:"note,omitempty"`
	RBACVerbs []string `json:"rbac_verbs,omitempty"`
}

// PodDisruptionBudget represents a Kubernetes PDB.
type PodDisruptionBudget struct {
	Name           string                 `json:"name"`
	MinAvailable   interface{}            `json:"min_available,omitempty"`
	MaxUnavailable interface{}            `json:"max_unavailable,omitempty"`
	Selector       map[string]interface{} `json:"selector,omitempty"`
	Source         string                 `json:"source"`
}

// HorizontalPodAutoscaler represents a Kubernetes HPA.
type HorizontalPodAutoscaler struct {
	Name        string      `json:"name"`
	TargetRef   string      `json:"target_ref,omitempty"`
	MinReplicas interface{} `json:"min_replicas,omitempty"`
	MaxReplicas interface{} `json:"max_replicas,omitempty"`
	Metrics     []string    `json:"metrics,omitempty"`
	Source      string      `json:"source"`
}

// APITypeDefinition represents a Go struct definition found in *_types.go files,
// typically defining a Kubernetes Custom Resource spec or its sub-components.
type APITypeDefinition struct {
	Name       string     `json:"name"`
	Doc        string     `json:"doc,omitempty"`
	Fields     []APIField `json:"fields,omitempty"`
	Markers    []string   `json:"markers,omitempty"`
	Source     string     `json:"source"`
	IsSpec     bool       `json:"is_spec,omitempty"`
	IsStatus   bool       `json:"is_status,omitempty"`
}

// APIField represents a single field within an API type struct.
type APIField struct {
	Name        string   `json:"name"`
	GoType      string   `json:"type"`
	JSONTag     string   `json:"json_tag,omitempty"`
	Doc         string   `json:"doc,omitempty"`
	Markers     []string `json:"markers,omitempty"`
	Required    bool     `json:"required,omitempty"`
	Default     string   `json:"default,omitempty"`
	SecretRef   bool     `json:"secret_ref,omitempty"`
	Embedded    bool     `json:"embedded,omitempty"`
}

// OperatorConstant represents a const or var declaration extracted from Go source
// that defines operator configuration (images, ports, timeouts, env vars, etc.).
type OperatorConstant struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	GoType   string `json:"type,omitempty"`
	Category string `json:"category"` // image, port, timeout, env_var, resource, name_pattern, general
	Doc      string `json:"doc,omitempty"`
	Source   string `json:"source"`
}

// ReconcileSequence represents the ordered reconciliation steps from a controller's
// Reconcile() method, with conditional guards.
type ReconcileSequence struct {
	Controller string          `json:"controller"`
	Source     string          `json:"source"`
	Steps      []ReconcileStep `json:"steps"`
}

// ReconcileStep is a single sub-resource reconciliation call within a controller.
type ReconcileStep struct {
	Order     int    `json:"order"`
	Method    string `json:"method"`
	Component string `json:"component,omitempty"`
	Condition string `json:"condition,omitempty"`
	Source    string `json:"source"`
}

// PrometheusMetric represents a Prometheus metric registration found in Go source.
type PrometheusMetric struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"` // gauge, counter, histogram, summary
	Help      string   `json:"help,omitempty"`
	Labels    []string `json:"labels,omitempty"`
	Subsystem string   `json:"subsystem,omitempty"`
	Namespace string   `json:"namespace,omitempty"`
	Source    string   `json:"source"`
}

// StatusCondition represents a status condition type and its associated reasons,
// defining the operator's observable state machine.
type StatusCondition struct {
	Type    string   `json:"type"`
	Reasons []string `json:"reasons,omitempty"`
	Source  string   `json:"source"`
}

// PlatformDetection holds platform capability checks and conditional resource
// creation patterns (e.g., OpenShift vs vanilla K8s).
type PlatformDetection struct {
	Capabilities []PlatformCapability  `json:"capabilities,omitempty"`
	Conditionals []PlatformConditional `json:"conditionals,omitempty"`
}

// PlatformCapability represents a detected platform capability field (e.g., IsOpenShift).
type PlatformCapability struct {
	Name   string `json:"name"`
	Check  string `json:"check,omitempty"`
	Source string `json:"source"`
}

// PlatformConditional represents a conditional resource creation guarded by a platform check.
type PlatformConditional struct {
	Condition    string `json:"condition"`
	ResourceKind string `json:"resource_kind,omitempty"`
	Action       string `json:"action"` // create, deploy, ensure, setup, allocate, watch
	Source       string `json:"source"`
}

// TemplateFile represents a Go template file that defines Kubernetes resources.
type TemplateFile struct {
	Path          string   `json:"path"`
	ResourceKinds []string `json:"resource_kinds,omitempty"`
	Conditionals  []string `json:"conditionals,omitempty"`
}
