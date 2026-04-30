package validator

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// KubeVersion represents a parsed Kubernetes version (major.minor).
type KubeVersion struct {
	Major int
	Minor int
}

// String returns "major.minor".
func (v KubeVersion) String() string {
	return fmt.Sprintf("%d.%d", v.Major, v.Minor)
}

// AtLeast returns true if v >= other.
func (v KubeVersion) AtLeast(other KubeVersion) bool {
	if v.Major != other.Major {
		return v.Major > other.Major
	}
	return v.Minor >= other.Minor
}

// ParseKubeVersion parses a version string like "1.27", "v1.27.3", or "4.14"
// (OCP versions are translated to kube versions).
func ParseKubeVersion(s string) (KubeVersion, error) {
	s = strings.TrimPrefix(s, "v")
	parts := strings.Split(s, ".")
	if len(parts) < 2 {
		return KubeVersion{}, fmt.Errorf("invalid version %q: need at least major.minor", s)
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return KubeVersion{}, fmt.Errorf("invalid major version in %q: %w", s, err)
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return KubeVersion{}, fmt.Errorf("invalid minor version in %q: %w", s, err)
	}
	return KubeVersion{Major: major, Minor: minor}, nil
}

// OCPToKubeVersion converts an OCP version to the corresponding Kubernetes
// version. OCP 4.x maps to Kube 1.(x+13), e.g. OCP 4.14 -> Kube 1.27.
func OCPToKubeVersion(ocpVersion string) (KubeVersion, error) {
	v, err := ParseKubeVersion(ocpVersion)
	if err != nil {
		return KubeVersion{}, err
	}
	if v.Major == 4 {
		return KubeVersion{Major: 1, Minor: v.Minor + 13}, nil
	}
	if v.Major == 1 {
		return v, nil // already a Kubernetes version
	}
	return KubeVersion{}, fmt.Errorf("unsupported version %q: only OCP 4.x and Kubernetes 1.x are supported", ocpVersion)
}

// APIDeprecation records when a Kubernetes API was deprecated and removed.
type APIDeprecation struct {
	Group          string
	Version        string
	Kind           string
	DeprecatedIn   KubeVersion // version where API was deprecated
	RemovedIn      KubeVersion // version where API was removed (0.0 if not yet removed)
	Replacement    string      // the replacement API version
}

// knownDeprecations lists well-known Kubernetes API deprecations relevant to
// operator development.
var knownDeprecations = []APIDeprecation{
	// CRD v1beta1 -> v1
	{
		Group: "apiextensions.k8s.io", Version: "v1beta1", Kind: "CustomResourceDefinition",
		DeprecatedIn: KubeVersion{1, 16}, RemovedIn: KubeVersion{1, 22},
		Replacement: "apiextensions.k8s.io/v1",
	},
	// Webhook v1beta1 -> v1
	{
		Group: "admissionregistration.k8s.io", Version: "v1beta1", Kind: "ValidatingWebhookConfiguration",
		DeprecatedIn: KubeVersion{1, 16}, RemovedIn: KubeVersion{1, 22},
		Replacement: "admissionregistration.k8s.io/v1",
	},
	{
		Group: "admissionregistration.k8s.io", Version: "v1beta1", Kind: "MutatingWebhookConfiguration",
		DeprecatedIn: KubeVersion{1, 16}, RemovedIn: KubeVersion{1, 22},
		Replacement: "admissionregistration.k8s.io/v1",
	},
	// PodSecurityPolicy
	{
		Group: "policy", Version: "v1beta1", Kind: "PodSecurityPolicy",
		DeprecatedIn: KubeVersion{1, 21}, RemovedIn: KubeVersion{1, 25},
		Replacement: "Pod Security Admission (built-in)",
	},
	// Ingress v1beta1 -> v1
	{
		Group: "networking.k8s.io", Version: "v1beta1", Kind: "Ingress",
		DeprecatedIn: KubeVersion{1, 19}, RemovedIn: KubeVersion{1, 22},
		Replacement: "networking.k8s.io/v1",
	},
	{
		Group: "extensions", Version: "v1beta1", Kind: "Ingress",
		DeprecatedIn: KubeVersion{1, 14}, RemovedIn: KubeVersion{1, 22},
		Replacement: "networking.k8s.io/v1",
	},
	// HPA v2beta2 -> v2
	{
		Group: "autoscaling", Version: "v2beta2", Kind: "HorizontalPodAutoscaler",
		DeprecatedIn: KubeVersion{1, 23}, RemovedIn: KubeVersion{1, 26},
		Replacement: "autoscaling/v2",
	},
	// FlowSchema/PriorityLevelConfiguration v1beta2 -> v1beta3
	{
		Group: "flowcontrol.apiserver.k8s.io", Version: "v1beta2", Kind: "FlowSchema",
		DeprecatedIn: KubeVersion{1, 26}, RemovedIn: KubeVersion{1, 29},
		Replacement: "flowcontrol.apiserver.k8s.io/v1",
	},
}

// VersionCompatIssue reports an API usage that is incompatible with a target version.
type VersionCompatIssue struct {
	Severity    string `json:"severity"` // "error" (removed), "warning" (deprecated)
	APIGroup    string `json:"api_group"`
	APIVersion  string `json:"api_version"`
	Kind        string `json:"kind"`
	Source      string `json:"source"`
	Message     string `json:"message"`
	Replacement string `json:"replacement,omitempty"`
}

// VersionCompatResult holds all compatibility issues for a target version.
type VersionCompatResult struct {
	TargetVersion string               `json:"target_version"`
	KubeVersion   string               `json:"kube_version"`
	Issues        []VersionCompatIssue `json:"issues"`
	Compatible    bool                 `json:"compatible"`
}

// CheckVersionCompat checks CRD and API usage against a target OCP or Kubernetes
// version. archData should contain the architecture extraction JSON
// (parsed as map[string]interface{}).
func CheckVersionCompat(archData map[string]interface{}, targetVersion string) (*VersionCompatResult, error) {
	target, err := resolveTargetVersion(targetVersion)
	if err != nil {
		return nil, err
	}

	result := &VersionCompatResult{
		TargetVersion: targetVersion,
		KubeVersion:   target.String(),
		Compatible:    true,
	}

	// Check CRDs for deprecated API versions
	crds := getSlice(archData, "crds")
	for _, raw := range crds {
		crd, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		source, _ := crd["source"].(string)
		checkCRDCompat(crd, target, source, result)
	}

	// Check webhooks
	webhooks := getSlice(archData, "webhooks")
	for _, raw := range webhooks {
		wh, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		source, _ := wh["source"].(string)
		checkWebhookCompat(wh, target, source, result)
	}

	// Check ingress routing for deprecated API versions
	ingress := getSlice(archData, "ingress_routing")
	for _, raw := range ingress {
		ing, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		source, _ := ing["source"].(string)
		checkIngressCompat(ing, target, source, result)
	}

	sort.Slice(result.Issues, func(i, j int) bool {
		ri := severityRank(result.Issues[i].Severity)
		rj := severityRank(result.Issues[j].Severity)
		if ri != rj {
			return ri < rj
		}
		return result.Issues[i].Source < result.Issues[j].Source
	})

	return result, nil
}

func resolveTargetVersion(ver string) (KubeVersion, error) {
	return OCPToKubeVersion(ver)
}

var apiVersionRe = regexp.MustCompile(`^([a-z0-9.-]+)/([a-z0-9]+)$`)

func checkCRDCompat(crd map[string]interface{}, target KubeVersion, source string, result *VersionCompatResult) {
	group, _ := crd["group"].(string)
	version, _ := crd["version"].(string)
	kind, _ := crd["kind"].(string)

	for _, dep := range knownDeprecations {
		if dep.Group == group && dep.Version == version {
			if dep.Kind != "" && kind != "" && dep.Kind != kind {
				continue
			}
			addIssue(dep, target, source, kind, result)
		}
	}
}

func checkWebhookCompat(wh map[string]interface{}, target KubeVersion, source string, result *VersionCompatResult) {
	whType, _ := wh["type"].(string)
	kind := "ValidatingWebhookConfiguration"
	if whType == "mutating" {
		kind = "MutatingWebhookConfiguration"
	}
	// Webhooks defined in YAML manifests may reference v1beta1
	for _, dep := range knownDeprecations {
		if dep.Kind == kind && dep.Version == "v1beta1" {
			// Check if the source file hints at v1beta1 usage
			if source != "" && (strings.Contains(source, "v1beta1") || strings.Contains(source, "beta")) {
				addIssue(dep, target, source, kind, result)
			}
		}
	}
}

func checkIngressCompat(ing map[string]interface{}, target KubeVersion, source string, result *VersionCompatResult) {
	ingressKind, _ := ing["kind"].(string)
	if ingressKind != "Ingress" {
		return
	}
	for _, dep := range knownDeprecations {
		if dep.Kind == "Ingress" {
			if source != "" && (strings.Contains(source, dep.Group) || strings.Contains(source, dep.Version)) {
				addIssue(dep, target, source, "Ingress", result)
			}
		}
	}
}

func addIssue(dep APIDeprecation, target KubeVersion, source, kind string, result *VersionCompatResult) {
	if dep.RemovedIn.Major > 0 && target.AtLeast(dep.RemovedIn) {
		result.Compatible = false
		result.Issues = append(result.Issues, VersionCompatIssue{
			Severity:    "error",
			APIGroup:    dep.Group,
			APIVersion:  dep.Version,
			Kind:        kind,
			Source:      source,
			Message:     fmt.Sprintf("%s/%s %s was removed in Kubernetes %s", dep.Group, dep.Version, kind, dep.RemovedIn),
			Replacement: dep.Replacement,
		})
	} else if target.AtLeast(dep.DeprecatedIn) {
		result.Issues = append(result.Issues, VersionCompatIssue{
			Severity:    "warning",
			APIGroup:    dep.Group,
			APIVersion:  dep.Version,
			Kind:        kind,
			Source:      source,
			Message:     fmt.Sprintf("%s/%s %s is deprecated since Kubernetes %s", dep.Group, dep.Version, kind, dep.DeprecatedIn),
			Replacement: dep.Replacement,
		})
	}
}

func severityRank(s string) int {
	switch s {
	case "error":
		return 0
	case "warning":
		return 1
	case "info":
		return 2
	default:
		return 3
	}
}

func getSlice(m map[string]interface{}, key string) []interface{} {
	v, ok := m[key]
	if !ok {
		return nil
	}
	s, ok := v.([]interface{})
	if !ok {
		return nil
	}
	return s
}
