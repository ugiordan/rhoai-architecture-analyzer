package extractor

import (
	"sort"
	"strings"
)

// ComponentTier classifies a component's importance in the platform.
type ComponentTier string

const (
	TierCore       ComponentTier = "core"       // platform infrastructure (operator, dashboard, etc.)
	TierML         ComponentTier = "ml"         // ML serving/training (kserve, codeflare, ray, etc.)
	TierData       ComponentTier = "data"       // data pipeline (data-science-pipelines, etc.)
	TierMonitoring ComponentTier = "monitoring" // observability (trustyai, model-mesh-monitoring, etc.)
	TierIntegration ComponentTier = "integration" // third-party integrations
	TierUnknown    ComponentTier = "unknown"
)

// ComponentType classifies a component's function.
type ComponentType string

const (
	TypeOperator   ComponentType = "operator"
	TypeController ComponentType = "controller"
	TypeServer     ComponentType = "server"
	TypeAPI        ComponentType = "api"
	TypeUI         ComponentType = "ui"
	TypeLibrary    ComponentType = "library"
	TypeUnknown    ComponentType = "unknown"
)

// ComponentMapEntry represents a single component in the platform component map.
type ComponentMapEntry struct {
	Name        string        `json:"name"`
	Repo        string        `json:"repo,omitempty"`
	Tier        ComponentTier `json:"tier"`
	Type        ComponentType `json:"type"`
	ImageCount  int           `json:"image_count"`
	CRDCount    int           `json:"crd_count"`
	HasOverlays bool          `json:"has_overlays"`
	Features    []string      `json:"features,omitempty"`
}

// ComponentMap is a structured catalog of all platform components.
type ComponentMap struct {
	Platform   string              `json:"platform"`
	Components []ComponentMapEntry `json:"components"`
	Summary    ComponentMapSummary `json:"summary"`
}

// ComponentMapSummary provides aggregate statistics about the component map.
type ComponentMapSummary struct {
	TotalComponents int            `json:"total_components"`
	TotalImages     int            `json:"total_images"`
	TotalCRDs       int            `json:"total_crds"`
	ByTier          map[string]int `json:"by_tier"`
	ByType          map[string]int `json:"by_type"`
}

// BuildComponentMap creates a ComponentMap from discovered kustomize components.
// If org is provided, it's used to construct repo URLs.
func BuildComponentMap(discovery *PlatformDiscovery, org string) *ComponentMap {
	if discovery == nil || len(discovery.Components) == 0 {
		return &ComponentMap{
			Summary: ComponentMapSummary{
				ByTier: make(map[string]int),
				ByType: make(map[string]int),
			},
		}
	}

	var entries []ComponentMapEntry
	for _, comp := range discovery.Components {
		entry := ComponentMapEntry{
			Name:        comp.Name,
			Tier:        classifyTier(comp.Name),
			Type:        classifyType(comp),
			ImageCount:  len(comp.ImageParams),
			CRDCount:    len(comp.ManagedCRDs),
			HasOverlays: len(comp.OverlayPaths) > 0,
			Features:    comp.FeatureFlags,
		}
		if org != "" {
			entry.Repo = org + "/" + comp.Name
		}
		entries = append(entries, entry)
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Tier != entries[j].Tier {
			return tierOrder(entries[i].Tier) < tierOrder(entries[j].Tier)
		}
		return entries[i].Name < entries[j].Name
	})

	cm := &ComponentMap{
		Components: entries,
		Summary:    computeSummary(entries),
	}
	return cm
}

// tierOrder returns a sort priority for tiers (core first, unknown last).
func tierOrder(t ComponentTier) int {
	switch t {
	case TierCore:
		return 0
	case TierML:
		return 1
	case TierData:
		return 2
	case TierMonitoring:
		return 3
	case TierIntegration:
		return 4
	default:
		return 5
	}
}

// classifyTier assigns a tier based on component name patterns.
// This is heuristic-based and can be overridden via configuration.
// More specific patterns are checked first to avoid false matches
// (e.g. "codeflare-operator" should be ML, not core).
func classifyTier(name string) ComponentTier {
	lower := strings.ToLower(name)

	// ML serving/training (check before core, since some have "operator" in name)
	mlPatterns := []string{"kserve", "modelmesh", "model-mesh", "codeflare", "ray", "kueue",
		"training", "vllm", "tgi", "serving"}
	for _, p := range mlPatterns {
		if strings.Contains(lower, p) {
			return TierML
		}
	}

	// Data pipeline (check before core, since some have "operator" in name)
	dataPatterns := []string{"pipeline", "data-science", "argo", "tekton", "elyra"}
	for _, p := range dataPatterns {
		if strings.Contains(lower, p) {
			return TierData
		}
	}

	// Monitoring/observability
	monPatterns := []string{"trusty", "monitor", "metric", "observ", "prometheus", "grafana"}
	for _, p := range monPatterns {
		if strings.Contains(lower, p) {
			return TierMonitoring
		}
	}

	// Core platform components (checked last to avoid shadowing more specific tiers)
	corePatterns := []string{"operator", "dashboard", "notebook", "controller-manager"}
	for _, p := range corePatterns {
		if strings.Contains(lower, p) {
			return TierCore
		}
	}

	return TierUnknown
}

// classifyType assigns a component type based on its characteristics.
func classifyType(comp KustomizeComponent) ComponentType {
	lower := strings.ToLower(comp.Name)

	if strings.Contains(lower, "operator") || strings.Contains(lower, "controller") {
		if len(comp.ManagedCRDs) > 0 {
			return TypeOperator
		}
		return TypeController
	}
	if strings.Contains(lower, "dashboard") || strings.Contains(lower, "ui") || strings.Contains(lower, "console") {
		return TypeUI
	}
	if strings.Contains(lower, "api") || strings.Contains(lower, "server") || strings.Contains(lower, "gateway") {
		return TypeServer
	}

	// If it has CRDs, it's likely an operator
	if len(comp.ManagedCRDs) > 0 {
		return TypeOperator
	}

	// If it has image params, it's a deployable component
	if len(comp.ImageParams) > 0 {
		return TypeServer
	}

	return TypeUnknown
}

func computeSummary(entries []ComponentMapEntry) ComponentMapSummary {
	s := ComponentMapSummary{
		TotalComponents: len(entries),
		ByTier:          make(map[string]int),
		ByType:          make(map[string]int),
	}
	for _, e := range entries {
		s.TotalImages += e.ImageCount
		s.TotalCRDs += e.CRDCount
		s.ByTier[string(e.Tier)]++
		s.ByType[string(e.Type)]++
	}
	return s
}
