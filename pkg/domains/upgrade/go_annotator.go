package upgrade

import (
	"strings"

	"github.com/ugiordan/architecture-analyzer/pkg/domains"
	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

type GoAnnotator struct{}

func (a *GoAnnotator) Annotate(g *graph.CPG, archData *domains.ArchitectureData) error {
	for _, fn := range g.NodesByKind(graph.NodeFunction) {
		a.annotateFunction(g, fn)
	}
	for _, cs := range g.NodesByKind(graph.NodeCallSite) {
		a.annotateCallSite(g, cs)
	}
	return nil
}

func (a *GoAnnotator) annotateFunction(g *graph.CPG, fn *graph.Node) {
	name := fn.Name
	nameLower := strings.ToLower(name)

	// upgrade:version_conversion
	if name == "ConvertTo" || name == "ConvertFrom" || name == "Hub" {
		g.SetAnnotation(fn.ID, AnnotVersionConversion, true)
	}

	// upgrade:migration
	if strings.Contains(nameLower, "migrate") || strings.Contains(nameLower, "upgrade") || strings.Contains(nameLower, "convert") {
		g.SetAnnotation(fn.ID, AnnotMigration, true)
	}
}

func (a *GoAnnotator) annotateCallSite(g *graph.CPG, cs *graph.Node) {
	name := cs.Name

	// upgrade:feature_gate
	if strings.Contains(name, "featuregate.Enabled") ||
		strings.Contains(name, "DefaultFeatureGate.Enabled") ||
		strings.Contains(name, "FeatureGate.Enabled") {
		g.SetAnnotation(cs.ID, AnnotFeatureGate, true)
	}

	// upgrade:version_check
	if strings.Contains(name, "getOCPVersion") ||
		strings.Contains(name, "semver.Compare") ||
		strings.Contains(name, "version.MustParseSemantic") {
		g.SetAnnotation(cs.ID, AnnotVersionCheck, true)
	}

	// upgrade:pre_release_api — flags pre-release API version usage.
	// Note: pre-release (alpha/beta) is not the same as deprecated. A project's
	// own v1alpha1 types that are the current and only version are not deprecated.
	if strings.Contains(name, "v1alpha1.") || strings.Contains(name, "v1beta1.") {
		g.SetAnnotation(cs.ID, AnnotPreReleaseAPI, true)
	}
}
