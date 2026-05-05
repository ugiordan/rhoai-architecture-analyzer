package extractor

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// prometheusRegistrationRE matches prometheus.New* and promauto.New* calls.
var prometheusRegistrationRE = regexp.MustCompile(
	`(?:prometheus|promauto)\.New(Gauge|Counter|Histogram|Summary)(Vec)?\(`,
)

// prometheusOptsFieldRE extracts a field value from a prometheus.*Opts struct literal.
func prometheusOptsFieldRE(field string) *regexp.Regexp {
	return regexp.MustCompile(field + `:\s*"([^"]*)"`)
}

var (
	promNameRE      = prometheusOptsFieldRE("Name")
	promHelpRE      = prometheusOptsFieldRE("Help")
	promNamespaceRE = prometheusOptsFieldRE("Namespace")
	promSubsystemRE = prometheusOptsFieldRE("Subsystem")
)

// prometheusLabelsRE extracts the labels slice from a Vec registration.
var prometheusLabelsRE = regexp.MustCompile(`\[\]string\{([^}]*)\}`)

// extractPrometheusMetrics scans Go source for Prometheus metric registrations.
func extractPrometheusMetrics(repoPath string) []PrometheusMetric {
	var metrics []PrometheusMetric
	seen := make(map[string]bool)

	filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && strings.Contains(path, "/vendor") {
			return filepath.SkipDir
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil || len(data) > 500*1024 {
			return nil
		}
		content := string(data)

		if !strings.Contains(content, "prometheus.") && !strings.Contains(content, "promauto.") {
			return nil
		}

		relPath := relativePath(repoPath, path)
		parsed := parsePrometheusMetrics(content, relPath)
		for _, m := range parsed {
			key := m.Name + ":" + m.Source
			if !seen[key] && m.Name != "" {
				seen[key] = true
				metrics = append(metrics, m)
			}
		}
		return nil
	})

	if metrics == nil {
		metrics = []PrometheusMetric{}
	}
	return metrics
}

// parsePrometheusMetrics extracts metric definitions from file content.
func parsePrometheusMetrics(content, source string) []PrometheusMetric {
	var metrics []PrometheusMetric

	// Find all registration calls and extract the surrounding block
	locs := prometheusRegistrationRE.FindAllStringIndex(content, -1)
	for _, loc := range locs {
		// Find the matching registration type
		match := prometheusRegistrationRE.FindStringSubmatch(content[loc[0]:loc[1]])
		if len(match) < 2 {
			continue
		}
		metricType := strings.ToLower(match[1]) // Gauge -> gauge, Counter -> counter, etc.
		isVec := len(match) >= 3 && match[2] == "Vec"

		// Extract a chunk of text after the registration call (opts struct + labels)
		end := loc[1] + 1000
		if end > len(content) {
			end = len(content)
		}
		chunk := content[loc[0]:end]

		// Extract opts fields
		name := extractRegexField(promNameRE, chunk)
		help := extractRegexField(promHelpRE, chunk)
		namespace := extractRegexField(promNamespaceRE, chunk)
		subsystem := extractRegexField(promSubsystemRE, chunk)

		if name == "" {
			continue
		}

		// Compose the full metric name
		composedName := composeMetricName(name, namespace, subsystem)

		// Extract labels for Vec types
		var labels []string
		if isVec {
			labels = extractLabels(chunk)
		}

		metrics = append(metrics, PrometheusMetric{
			Name:      composedName,
			Type:      metricType,
			Help:      help,
			Labels:    labels,
			Subsystem: subsystem,
			Namespace: namespace,
			Source:    source,
		})
	}

	return metrics
}

// composeMetricName builds the full metric name following Prometheus conventions.
func composeMetricName(name, namespace, subsystem string) string {
	// If name already starts with namespace or subsystem, use as-is
	if namespace != "" && strings.HasPrefix(name, namespace) {
		return name
	}
	if subsystem != "" && strings.HasPrefix(name, subsystem) {
		return name
	}

	parts := []string{}
	if namespace != "" {
		parts = append(parts, namespace)
	}
	if subsystem != "" {
		parts = append(parts, subsystem)
	}
	parts = append(parts, name)
	return strings.Join(parts, "_")
}

// extractRegexField extracts the first capture group from a regex match.
func extractRegexField(re *regexp.Regexp, text string) string {
	match := re.FindStringSubmatch(text)
	if len(match) >= 2 {
		return match[1]
	}
	return ""
}

// extractLabels extracts label names from a []string{...} literal.
func extractLabels(text string) []string {
	match := prometheusLabelsRE.FindStringSubmatch(text)
	if len(match) < 2 {
		return nil
	}
	raw := match[1]
	var labels []string
	for _, part := range strings.Split(raw, ",") {
		label := strings.TrimSpace(part)
		label = strings.Trim(label, `"`)
		if label != "" {
			labels = append(labels, label)
		}
	}
	return labels
}
