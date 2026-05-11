package extractor

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// LabelContract represents a well-known Kubernetes label or annotation that
// implies a cross-component integration contract (e.g., Kueue queue assignment,
// KServe inference service binding).
type LabelContract struct {
	Label       string `json:"label"`       // the label/annotation key
	Integration string `json:"integration"` // target component/system (e.g., "Kueue", "KServe")
	Type        string `json:"type"`        // "label" or "annotation"
	Context     string `json:"context"`     // where found: "job-template", "deployment", "pod-template", "source-code", etc.
	Source      string `json:"source"`      // file:line
}

// wellKnownLabelEntry defines a well-known label/annotation and the integration it implies.
type wellKnownLabelEntry struct {
	key         string
	integration string
}

// wellKnownLabels maps label keys to their integration targets.
// These are labels that, when present on a Kubernetes resource, indicate
// the component relies on an external system for scheduling, serving, etc.
var wellKnownLabels = []wellKnownLabelEntry{
	{"kueue.x-k8s.io/queue-name", "Kueue"},
	{"kueue.x-k8s.io/priority-class", "Kueue"},
	{"serving.kserve.io/inferenceservice", "KServe"},
	{"serving.knative.dev/service", "Knative Serving"},
	{"app.kubernetes.io/managed-by", "management relationship"},
	{"istio.io/rev", "Istio"},
	{"sidecar.istio.io/inject", "Istio"},
	{"prometheus.io/scrape", "Prometheus"},
	{"prometheus.io/port", "Prometheus"},
	{"cert-manager.io/cluster-issuer", "cert-manager"},
	{"cert-manager.io/issuer", "cert-manager"},
}

// wellKnownLabelSet builds a fast lookup set from wellKnownLabels.
var wellKnownLabelSet map[string]string

func init() {
	wellKnownLabelSet = make(map[string]string, len(wellKnownLabels))
	for _, entry := range wellKnownLabels {
		wellKnownLabelSet[entry.key] = entry.integration
	}
}

// labelLiteralRE matches well-known label key string literals in Go source.
// Built dynamically from wellKnownLabels.
var labelLiteralRE *regexp.Regexp

func init() {
	var escaped []string
	for _, entry := range wellKnownLabels {
		escaped = append(escaped, regexp.QuoteMeta(entry.key))
	}
	labelLiteralRE = regexp.MustCompile(`"(` + strings.Join(escaped, "|") + `)"`)
}

// yamlPatternsForLabels covers all manifest types that can carry labels/annotations.
var yamlPatternsForLabels = []string{
	"**/*.yaml",
	"**/*.yml",
	"**/*.yaml.tmpl",
	"**/*.yml.tmpl",
}

// extractLabelContracts scans YAML manifests and Go source for well-known
// Kubernetes labels and annotations that create cross-component integration contracts.
func extractLabelContracts(repoPath string) []LabelContract {
	var contracts []LabelContract

	// Pass 1: YAML manifests
	contracts = append(contracts, scanYAMLForLabels(repoPath)...)

	// Pass 2: Go source string literals
	contracts = append(contracts, scanGoForLabelLiterals(repoPath)...)

	// Deduplicate: same label + source is redundant
	contracts = deduplicateLabelContracts(contracts)

	return contracts
}

// scanYAMLForLabels parses YAML manifests and checks metadata.labels and
// metadata.annotations (plus spec.template.metadata for workload resources)
// against the well-known label set.
func scanYAMLForLabels(repoPath string) []LabelContract {
	files := findYAMLFiles(repoPath, yamlPatternsForLabels)
	var contracts []LabelContract

	for _, fpath := range files {
		relPath := relativePath(repoPath, fpath)

		// Handle .tmpl files separately
		var docs []map[string]interface{}
		if strings.HasSuffix(fpath, ".tmpl") {
			docs = parseTemplateYAML(fpath)
		} else {
			docs = parseYAMLSafe(fpath)
		}

		for _, doc := range docs {
			kind, _ := doc["kind"].(string)
			if kind == "" {
				continue
			}

			// Check top-level metadata
			contracts = append(contracts, checkMetadataForLabels(doc, kind, relPath)...)

			// Check spec.template.metadata for workload resources
			spec, _ := doc["spec"].(map[string]interface{})
			if spec == nil {
				continue
			}
			template, _ := spec["template"].(map[string]interface{})
			if template == nil {
				// Also check spec.jobTemplate.spec.template for CronJobs
				jobTemplate, _ := spec["jobTemplate"].(map[string]interface{})
				if jobTemplate != nil {
					jobSpec, _ := jobTemplate["spec"].(map[string]interface{})
					if jobSpec != nil {
						template, _ = jobSpec["template"].(map[string]interface{})
					}
				}
			}
			if template != nil {
				context := inferTemplateContext(kind)
				contracts = append(contracts, checkMapForLabels(template, context, relPath)...)
			}
		}
	}

	return contracts
}

// checkMetadataForLabels checks the top-level metadata.labels and metadata.annotations
// of a YAML document for well-known entries.
func checkMetadataForLabels(doc map[string]interface{}, kind, source string) []LabelContract {
	metadata, _ := doc["metadata"].(map[string]interface{})
	if metadata == nil {
		return nil
	}
	context := strings.ToLower(kind)
	return checkMapForLabels(metadata, context, source)
}

// checkMapForLabels scans labels and annotations maps in a metadata-like
// object for well-known keys.
func checkMapForLabels(obj map[string]interface{}, context, source string) []LabelContract {
	metadata, _ := obj["metadata"].(map[string]interface{})
	if metadata == nil {
		// obj might be the metadata itself
		metadata = obj
	}

	var contracts []LabelContract

	if labels, ok := metadata["labels"].(map[string]interface{}); ok {
		for key := range labels {
			if integration, found := wellKnownLabelSet[key]; found {
				contracts = append(contracts, LabelContract{
					Label:       key,
					Integration: integration,
					Type:        "label",
					Context:     context,
					Source:      source,
				})
			}
		}
	}

	if annotations, ok := metadata["annotations"].(map[string]interface{}); ok {
		for key := range annotations {
			if integration, found := wellKnownLabelSet[key]; found {
				contracts = append(contracts, LabelContract{
					Label:       key,
					Integration: integration,
					Type:        "annotation",
					Context:     context,
					Source:      source,
				})
			}
		}
	}

	return contracts
}

// inferTemplateContext returns the context string for a pod template based
// on the parent resource kind.
func inferTemplateContext(kind string) string {
	switch kind {
	case "Job":
		return "job-template"
	case "CronJob":
		return "cronjob-template"
	case "Deployment":
		return "deployment-template"
	case "StatefulSet":
		return "statefulset-template"
	case "DaemonSet":
		return "daemonset-template"
	case "ReplicaSet":
		return "replicaset-template"
	default:
		return "pod-template"
	}
}

// scanGoForLabelLiterals scans Go source files for string literals matching
// well-known label keys.
func scanGoForLabelLiterals(repoPath string) []LabelContract {
	goFiles := findFiles(repoPath, []string{"**/*.go"})
	var contracts []LabelContract

	for _, fpath := range goFiles {
		if strings.Contains(fpath, "_test.go") || strings.Contains(fpath, "/vendor/") {
			continue
		}

		info, err := os.Lstat(fpath)
		if err != nil || info.Mode()&os.ModeSymlink != 0 {
			continue
		}
		if info.Size() > maxFileSize {
			continue
		}

		data, err := os.ReadFile(fpath)
		if err != nil {
			continue
		}

		content := string(data)
		lines := strings.Split(content, "\n")
		relPath := relativePath(repoPath, fpath)

		for lineNum, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "//") {
				continue
			}

			matches := labelLiteralRE.FindAllStringSubmatch(line, -1)
			for _, match := range matches {
				if len(match) < 2 {
					continue
				}
				key := match[1]
				integration, found := wellKnownLabelSet[key]
				if !found {
					continue
				}
				contracts = append(contracts, LabelContract{
					Label:       key,
					Integration: integration,
					Type:        "label", // in source code we can't always distinguish, default to label
					Context:     "source-code",
					Source:      fmt.Sprintf("%s:%d", relPath, lineNum+1),
				})
			}
		}
	}

	return contracts
}

// deduplicateLabelContracts removes entries with the same label+source combination,
// preferring YAML-sourced entries (which have accurate type info) over source-code ones.
func deduplicateLabelContracts(contracts []LabelContract) []LabelContract {
	type dedupKey struct {
		label   string
		source  string
		context string
	}
	seen := make(map[dedupKey]bool)
	var result []LabelContract

	for _, c := range contracts {
		key := dedupKey{label: c.Label, source: c.Source, context: c.Context}
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, c)
	}

	return result
}
