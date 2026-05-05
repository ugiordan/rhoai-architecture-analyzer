package extractor

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var managerFilePatterns = []string{
	"cmd/main.go",
	"cmd/*/main.go",
	"main.go",
}

// Regex patterns for cache configuration detection.
var (
	ctrlNewManagerRE = regexp.MustCompile(`ctrl\.NewManager|manager\.New`)
	cacheOptionsRE   = regexp.MustCompile(`cache\.Options\s*\{`)
	byObjectRE       = regexp.MustCompile(`ByObject:\s*map`)
	defaultTransfRE  = regexp.MustCompile(`DefaultTransform\s*:`)
	defaultNsRE      = regexp.MustCompile(`DefaultNamespaces\s*:`)
	disableForRE     = regexp.MustCompile(`DisableFor\s*:\s*\[\]client\.Object\s*\{`)
	transformRE      = regexp.MustCompile(`Transform\s*:\s*(\w+)`)
	// Matches the start of a ByObject entry to locate it, then extractBlock gets the body.
	byObjectEntryStartRE   = regexp.MustCompile(`&(\w+)\.(\w+)\{\s*\}\s*:\s*cache\.ByObject\s*\{`)
	byObjectEntryInlineStartRE = regexp.MustCompile(`&(\w+)\.(\w+)\{\s*\}\s*:\s*\{`)
	// Matches variable ref: &pkg.Kind{}: varName,
	byObjectEntryVarRE = regexp.MustCompile(`&(\w+)\.(\w+)\{\s*\}\s*:\s*(\w+)\s*,`)
	// Matches variable definition start: varName := cache.ByObject{
	byObjectVarStartRE = regexp.MustCompile(`(\w+)\s*(?::=|=)\s*cache\.ByObject\s*\{`)
	labelFilterRE      = regexp.MustCompile(`Label\s*:`)
	fieldFilterRE      = regexp.MustCompile(`Field\s*:`)
	namespacesRE       = regexp.MustCompile(`Namespaces\s*:`)
	// Matches &pkg.Kind{} anywhere in content (used for multi-line Get/List scanning and DisableFor entries).
	clientGetTypeRE  = regexp.MustCompile(`&(\w+)\.(\w+)\{\s*\}`)
	clientListTypeRE = regexp.MustCompile(`&(\w+)\.(\w+)List\{\s*\}`)
	// getListRE matches .Get( and .List( call sites for implicit informer detection.
	getListRE = regexp.MustCompile(`\.(Get|List)\(`)
)

// disableEntryRE is an alias for clientGetTypeRE (same pattern: &pkg.Kind{}).
var disableEntryRE = clientGetTypeRE

// ignoredTypes contains non-API types that should not be treated as implicit informers.
var ignoredTypes = map[string]bool{
	"StatusWriter": true, "Client": true, "Scheme": true,
	"RESTMapper": true, "Manager": true, "Logger": true,
	"Options": true, "Config": true, "Builder": true,
}

// extractCacheConfig analyzes the controller-runtime cache configuration from
// Go source and cross-references with watches and deployments to detect OOM risks.
func extractCacheConfig(repoPath string, watches []ControllerWatch, deployments []Deployment) *CacheConfig {
	managerFiles := findFiles(repoPath, managerFilePatterns)
	if len(managerFiles) == 0 {
		return nil
	}

	// Also scan files that might contain cache configuration (referenced from main)
	cacheFiles := findFiles(repoPath, []string{
		"cmd/**/cache*.go",
		"internal/**/cache*.go",
		"pkg/**/cache*.go",
	})

	config := &CacheConfig{
		FilteredTypes:     []CacheFilteredType{},
		TransformTypes:    []CacheTransform{},
		DisabledTypes:     []string{},
		ImplicitInformers: []ImplicitInformer{},
		Issues:            []string{},
	}

	// Shared dedup set across all parseCacheOptions calls to prevent duplicates
	// when the same type appears in both the manager file and a cache helper.
	seenByObject := make(map[string]bool)

	// Find and parse manager entry points
	foundManager := false
	for _, fpath := range managerFiles {
		data, err := os.ReadFile(fpath)
		if err != nil {
			continue
		}
		content := string(data)
		if !ctrlNewManagerRE.MatchString(content) {
			continue
		}
		foundManager = true
		config.ManagerFile = relativePath(repoPath, fpath)
		parseCacheOptions(content, config, seenByObject)
	}

	if !foundManager {
		return nil
	}

	// Also parse dedicated cache config files
	for _, fpath := range cacheFiles {
		data, err := os.ReadFile(fpath)
		if err != nil {
			continue
		}
		content := string(data)
		if cacheOptionsRE.MatchString(content) || byObjectRE.MatchString(content) {
			parseCacheOptions(content, config, seenByObject)
		}
	}

	// Determine cache scope
	if config.CacheScope == "" {
		config.CacheScope = "cluster-wide"
	}

	// Extract GOMEMLIMIT and memory limits from deployments
	extractMemoryConfig(config, deployments)

	// Compute GOMEMLIMIT ratio if both values are present
	computeMemLimitRatio(config)

	// Cross-reference watches with cache filters to find unfiltered informers
	crossReferenceWatches(config, watches)

	// Scan for implicit informers (client.Get for unwatched types)
	scanImplicitInformers(repoPath, config, watches)

	// Generate issues
	generateCacheIssues(config, watches)

	sort.Strings(config.Issues)

	return config
}

// parseCacheOptions extracts cache configuration details from Go source content.
// seenByObject is shared across calls to prevent duplicate entries when the same
// type appears in both the manager file and a cache helper file.
func parseCacheOptions(content string, config *CacheConfig, seenByObject map[string]bool) {
	if !cacheOptionsRE.MatchString(content) && !byObjectRE.MatchString(content) {
		return
	}

	if defaultTransfRE.MatchString(content) {
		config.DefaultTransform = true
		if m := transformRE.FindStringSubmatch(content[defaultTransfRE.FindStringIndex(content)[0]:]); m != nil {
			config.DefaultTransformFunc = m[1]
		}
	}

	if defaultNsRE.MatchString(content) {
		config.CacheScope = "namespace-scoped"
	}

	// Build variable definitions map for resolving variable references.
	// Uses extractBlock for brace-balanced body extraction.
	varDefs := make(map[string]string) // varName -> body content
	for _, idx := range byObjectVarStartRE.FindAllStringIndex(content, -1) {
		// Extract the variable name from the match
		match := byObjectVarStartRE.FindStringSubmatch(content[idx[0]:idx[1]])
		if match == nil {
			continue
		}
		// Find the opening brace position and extract the balanced body
		bracePos := strings.Index(content[idx[0]:], "{")
		if bracePos < 0 {
			continue
		}
		body := extractBlock(content[idx[0]+bracePos:], '{', '}')
		varDefs[match[1]] = body
	}

	// Helper to classify a ByObject body
	classifyBody := func(typeName, body string) {
		filterKind := "none"
		filter := ""
		if labelFilterRE.MatchString(body) {
			filterKind = "label"
			filter = extractLabelValue(body)
		} else if fieldFilterRE.MatchString(body) {
			filterKind = "field"
			filter = extractFieldValue(body)
		} else if namespacesRE.MatchString(body) {
			filterKind = "namespace"
			filter = "namespace-scoped"
		}

		if filterKind != "none" {
			config.FilteredTypes = append(config.FilteredTypes, CacheFilteredType{
				Type:       typeName,
				FilterKind: filterKind,
				Filter:     filter,
			})
		}

		if m := transformRE.FindStringSubmatch(body); m != nil {
			config.TransformTypes = append(config.TransformTypes, CacheTransform{
				Type:     typeName,
				Function: m[1],
			})
		}
	}

	// Pattern 1: &pkg.Kind{}: cache.ByObject{body} (brace-balanced extraction)
	for _, idx := range byObjectEntryStartRE.FindAllStringIndex(content, -1) {
		match := byObjectEntryStartRE.FindStringSubmatch(content[idx[0]:idx[1]])
		if match == nil {
			continue
		}
		typeName := match[1] + "." + match[2]
		if seenByObject[typeName] {
			continue
		}
		seenByObject[typeName] = true
		// Find the opening brace of cache.ByObject{ and extract balanced body
		braceStart := strings.LastIndex(content[idx[0]:idx[1]], "{")
		if braceStart < 0 {
			continue
		}
		body := extractBlock(content[idx[0]+braceStart:], '{', '}')
		classifyBody(typeName, body)
	}

	// Pattern 2: &pkg.Kind{}: {body} (inline without cache.ByObject prefix, brace-balanced)
	for _, idx := range byObjectEntryInlineStartRE.FindAllStringIndex(content, -1) {
		match := byObjectEntryInlineStartRE.FindStringSubmatch(content[idx[0]:idx[1]])
		if match == nil {
			continue
		}
		typeName := match[1] + "." + match[2]
		if seenByObject[typeName] {
			continue
		}
		seenByObject[typeName] = true
		braceStart := strings.LastIndex(content[idx[0]:idx[1]], "{")
		if braceStart < 0 {
			continue
		}
		body := extractBlock(content[idx[0]+braceStart:], '{', '}')
		classifyBody(typeName, body)
	}

	// Pattern 3: &pkg.Kind{}: varName, (variable reference)
	for _, match := range byObjectEntryVarRE.FindAllStringSubmatch(content, -1) {
		typeName := match[1] + "." + match[2]
		varName := match[3]
		if seenByObject[typeName] {
			continue
		}
		seenByObject[typeName] = true
		if body, ok := varDefs[varName]; ok {
			classifyBody(typeName, body)
		} else {
			// Variable defined elsewhere, mark as unknown for manual review
			config.FilteredTypes = append(config.FilteredTypes, CacheFilteredType{
				Type:       typeName,
				FilterKind: "unknown",
				Filter:     "unresolved variable: " + varName,
			})
		}
	}

	// Parse DisableFor entries. A type can appear in both ByObject (cached with
	// filters/transforms) and DisableFor (client reads bypass cache), so we use
	// a separate dedup set instead of seenByObject.
	if disableForRE.MatchString(content) {
		idx := disableForRE.FindStringIndex(content)
		if idx != nil {
			seenDisabled := make(map[string]bool)
			block := extractBlock(content[idx[0]:], '{', '}')
			for _, entry := range disableEntryRE.FindAllStringSubmatch(block, -1) {
				typeName := entry[1] + "." + entry[2]
				if !seenDisabled[typeName] {
					seenDisabled[typeName] = true
					config.DisabledTypes = append(config.DisabledTypes, typeName)
				}
			}
		}
	}
}

// extractBlock extracts text from the first open to its matching close,
// skipping braces inside string literals, backtick strings, rune literals,
// and comments.
func extractBlock(s string, open, close byte) string {
	depth := 0
	start := -1
	n := len(s)
	for i := 0; i < n; i++ {
		ch := s[i]

		// Skip double-quoted strings
		if ch == '"' {
			i++
			for i < n && s[i] != '"' {
				if s[i] == '\\' {
					i++ // skip escaped char
				}
				i++
			}
			continue
		}

		// Skip backtick (raw) strings
		if ch == '`' {
			i++
			for i < n && s[i] != '`' {
				i++
			}
			continue
		}

		// Skip rune literals (handles multi-char escapes like '\x41', '\u0041')
		if ch == '\'' {
			i++
			for i < n && s[i] != '\'' {
				if s[i] == '\\' {
					i++
				}
				i++
			}
			continue
		}

		// Skip single-line comments
		if ch == '/' && i+1 < n && s[i+1] == '/' {
			i += 2
			for i < n && s[i] != '\n' {
				i++
			}
			continue
		}

		// Skip multi-line comments
		if ch == '/' && i+1 < n && s[i+1] == '*' {
			i += 2
			for i+1 < n {
				if s[i] == '*' && s[i+1] == '/' {
					i++ // position on '/', loop increment moves past
					break
				}
				i++
			}
			continue
		}

		if ch == open {
			if start == -1 {
				start = i
			}
			depth++
		} else if ch == close {
			depth--
			if depth == 0 {
				return s[start : i+1]
			}
		}
	}
	if start >= 0 {
		return s[start:]
	}
	return ""
}

var labelValueRE = regexp.MustCompile(`"([^"]+)"\s*:\s*"([^"]*)"`)

// selectorFromSetRE matches labels.SelectorFromSet(labels.Set{"key": "value"})
var selectorFromSetRE = regexp.MustCompile(`SelectorFromSet\s*\(\s*labels\.Set\s*\{([^}]*)\}`)

// fieldSelectorRE matches fields.SelectorFromSet(fields.Set{"key": "value"}) or
// fields.ParseSelector("key=value") or fields.OneTermEqualSelector("key", "value")
var fieldSelectorRE = regexp.MustCompile(`fields\.(?:SelectorFromSet|ParseSelector|OneTermEqualSelector)\s*\(([^)]+)\)`)

// constantRefRE matches Go constant references like pkg.ConstName or just ConstName
var constantRefRE = regexp.MustCompile(`(\w+(?:\.\w+)?)\s*:\s*(\w+(?:\.\w+)?)`)

// extractLabelValue tries to extract a label key=value from a ByObject body.
// Handles both map literal {"key": "value"} and SelectorFromSet patterns.
// When label values use Go constants, includes the constant names as hints.
func extractLabelValue(body string) string {
	// Try SelectorFromSet pattern first
	if m := selectorFromSetRE.FindStringSubmatch(body); m != nil {
		setBody := m[1]
		// Try string literal pairs
		pairs := labelValueRE.FindAllStringSubmatch(setBody, -1)
		var parts []string
		for _, p := range pairs {
			parts = append(parts, p[1]+"="+p[2])
		}
		if len(parts) > 0 {
			return strings.Join(parts, ", ")
		}
		// Try constant references (config.Label: config.Value)
		constPairs := constantRefRE.FindAllStringSubmatch(setBody, -1)
		for _, p := range constPairs {
			parts = append(parts, p[1]+"="+p[2])
		}
		if len(parts) > 0 {
			return strings.Join(parts, ", ") + " (constants, resolved at runtime)"
		}
	}

	// Fall back to direct map literal matching
	matches := labelValueRE.FindAllStringSubmatch(body, -1)
	var parts []string
	for _, m := range matches {
		parts = append(parts, m[1]+"="+m[2])
	}
	if len(parts) > 0 {
		return strings.Join(parts, ", ")
	}
	return "label selector"
}

// extractFieldValue tries to extract a field selector expression from a ByObject body.
func extractFieldValue(body string) string {
	if m := fieldSelectorRE.FindStringSubmatch(body); m != nil {
		args := m[1]
		// Try to extract key=value pairs from the arguments
		pairs := labelValueRE.FindAllStringSubmatch(args, -1)
		var parts []string
		for _, p := range pairs {
			parts = append(parts, p[1]+"="+p[2])
		}
		if len(parts) > 0 {
			return strings.Join(parts, ", ")
		}
		// For ParseSelector("key=value") style
		if strings.Contains(args, `"`) {
			// Extract the string argument
			start := strings.Index(args, `"`)
			end := strings.LastIndex(args, `"`)
			if start >= 0 && end > start {
				return args[start+1 : end]
			}
		}
	}
	return "field selector"
}

// extractMemoryConfig pulls GOMEMLIMIT and memory limits from deployment specs.
// Prefers the "manager" container if multiple deployments/containers exist.
func extractMemoryConfig(config *CacheConfig, deployments []Deployment) {
	// First pass: look for a container named "manager" (convention for controller-runtime)
	for _, dep := range deployments {
		for _, c := range dep.Containers {
			if c.Name != "manager" {
				continue
			}
			if c.EnvVars != nil {
				if val, ok := c.EnvVars["GOMEMLIMIT"]; ok {
					config.GoMemLimit = val
				}
			}
			if c.Resources != nil {
				if limits, ok := c.Resources["limits"].(map[string]interface{}); ok {
					if mem, ok := limits["memory"].(string); ok && mem != "" {
						config.MemoryLimit = mem
					}
				}
			}
			if config.GoMemLimit != "" || config.MemoryLimit != "" {
				return
			}
		}
	}

	// Fallback: take first deployment/container with relevant values
	for _, dep := range deployments {
		for _, c := range dep.Containers {
			if c.EnvVars != nil {
				if val, ok := c.EnvVars["GOMEMLIMIT"]; ok && config.GoMemLimit == "" {
					config.GoMemLimit = val
				}
			}
			if c.Resources != nil {
				if limits, ok := c.Resources["limits"].(map[string]interface{}); ok {
					if mem, ok := limits["memory"].(string); ok && mem != "" && config.MemoryLimit == "" {
						config.MemoryLimit = mem
					}
				}
			}
		}
	}
}

// computeMemLimitRatio calculates the GOMEMLIMIT-to-memory-limit ratio and
// adds an issue if the ratio is outside the recommended 80-90% range.
func computeMemLimitRatio(config *CacheConfig) {
	if config.GoMemLimit == "" || config.MemoryLimit == "" {
		return
	}
	goMem := parseMemoryBytes(config.GoMemLimit)
	limit := parseMemoryBytes(config.MemoryLimit)
	if goMem <= 0 || limit <= 0 {
		return
	}
	ratio := float64(goMem) / float64(limit) * 100
	// Round to one decimal
	config.GoMemLimitRatio = float64(int(ratio*10+0.5)) / 10

	if ratio < 80 {
		config.Issues = append(config.Issues,
			fmt.Sprintf("GOMEMLIMIT ratio %.1f%% is below recommended 80%% minimum (GC cannot pressure-tune effectively)", ratio))
	} else if ratio > 90 {
		config.Issues = append(config.Issues,
			fmt.Sprintf("GOMEMLIMIT ratio %.1f%% exceeds recommended 90%% maximum (insufficient headroom for non-Go memory)", ratio))
	}
}

// parseMemoryBytes converts Kubernetes-style memory strings (e.g., "4Gi", "3600MiB",
// "512Mi", "1G", "256M") to bytes. Returns 0 if the format is unrecognized.
func parseMemoryBytes(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}

	// Try binary suffixes first (Ki, Mi, Gi, Ti, KiB, MiB, GiB, TiB)
	binarySuffixes := []struct {
		suffix     string
		multiplier int64
	}{
		{"TiB", 1 << 40}, {"Ti", 1 << 40},
		{"GiB", 1 << 30}, {"Gi", 1 << 30},
		{"MiB", 1 << 20}, {"Mi", 1 << 20},
		{"KiB", 1 << 10}, {"Ki", 1 << 10},
	}
	for _, bs := range binarySuffixes {
		if strings.HasSuffix(s, bs.suffix) {
			numStr := strings.TrimSuffix(s, bs.suffix)
			if v, err := strconv.ParseFloat(numStr, 64); err == nil {
				return int64(v * float64(bs.multiplier))
			}
		}
	}

	// Decimal suffixes (T, G, M, K)
	decimalSuffixes := []struct {
		suffix     string
		multiplier int64
	}{
		{"T", 1e12}, {"G", 1e9}, {"M", 1e6}, {"K", 1e3},
	}
	for _, ds := range decimalSuffixes {
		if strings.HasSuffix(s, ds.suffix) {
			numStr := strings.TrimSuffix(s, ds.suffix)
			if v, err := strconv.ParseFloat(numStr, 64); err == nil {
				return int64(v * float64(ds.multiplier))
			}
		}
	}

	// Plain number (bytes)
	if v, err := strconv.ParseInt(s, 10, 64); err == nil {
		return v
	}
	return 0
}

// crossReferenceWatches checks which watched types have cache filters.
// Uses both Kind-only and full type matching to handle different naming between
// cache config (pkg.Kind like "corev1.Secret") and watches (GVK like "/v1/Secret").
func crossReferenceWatches(config *CacheConfig, watches []ControllerWatch) {
	// Build sets of filtered and disabled types using both full type and Kind-only.
	// Also track "unknown" filter kinds separately: they should not suppress warnings.
	filteredKinds := make(map[string]bool)
	unknownFilterKinds := make(map[string]bool)
	for _, ft := range config.FilteredTypes {
		parts := strings.SplitN(ft.Type, ".", 2)
		kind := ft.Type
		if len(parts) == 2 {
			kind = parts[1]
		}
		if ft.FilterKind == "unknown" {
			unknownFilterKinds[kind] = true
			unknownFilterKinds[ft.Type] = true
		} else {
			filteredKinds[kind] = true
			filteredKinds[ft.Type] = true
		}
	}
	// Types with transforms (even without label/field filters) are still configured
	// in ByObject and should not be flagged as unfiltered.
	for _, tt := range config.TransformTypes {
		parts := strings.SplitN(tt.Type, ".", 2)
		kind := tt.Type
		if len(parts) == 2 {
			kind = parts[1]
		}
		filteredKinds[kind] = true
		filteredKinds[tt.Type] = true
	}
	disabledKinds := make(map[string]bool)
	for _, dt := range config.DisabledTypes {
		parts := strings.SplitN(dt, ".", 2)
		if len(parts) == 2 {
			disabledKinds[parts[1]] = true
		}
		disabledKinds[dt] = true
	}

	// Check each watched type
	watchedKinds := make(map[string]bool)
	for _, w := range watches {
		parts := strings.Split(w.GVK, "/")
		kind := parts[len(parts)-1]
		if watchedKinds[kind] {
			continue
		}
		watchedKinds[kind] = true

		if disabledKinds[kind] {
			continue
		}
		if filteredKinds[kind] {
			continue
		}
		suffix := ""
		if unknownFilterKinds[kind] {
			suffix = " (cache variable unresolved, manual review needed)"
		}
		config.Issues = append(config.Issues,
			fmt.Sprintf("Type %s is watched but has no cache filter (cluster-wide informer)%s", kind, suffix))
	}
}

// scanImplicitInformers finds client.Get/List calls for types not in the watch
// set, which silently create cluster-wide informers. Scans full file content to
// catch multi-line Get/List calls where the type argument is on a separate line.
func scanImplicitInformers(repoPath string, config *CacheConfig, watches []ControllerWatch) {
	watchedKinds := make(map[string]bool)
	for _, w := range watches {
		parts := strings.Split(w.GVK, "/")
		kind := parts[len(parts)-1]
		watchedKinds[kind] = true
	}

	disabledKinds := make(map[string]bool)
	for _, dt := range config.DisabledTypes {
		parts := strings.SplitN(dt, ".", 2)
		if len(parts) == 2 {
			disabledKinds[parts[1]] = true
		}
	}

	goFiles := findFiles(repoPath, []string{
		"internal/**/*.go",
		"pkg/**/*.go",
		"cmd/**/*.go",
	})

	seen := make(map[string]bool)
	for _, fpath := range goFiles {
		if strings.Contains(fpath, "_test.go") || strings.Contains(fpath, "/vendor/") {
			continue
		}
		data, err := os.ReadFile(fpath)
		if err != nil {
			continue
		}
		content := string(data)
		relPath := relativePath(repoPath, fpath)

		// Find all .Get( and .List( call regions, then scan for &pkg.Type{} within
		scanGetListCalls(content, relPath, watchedKinds, disabledKinds, seen, config)
	}
}

// scanGetListCalls finds client.Get/List calls in file content and extracts the
// type argument, handling multi-line calls correctly.
func scanGetListCalls(content, relPath string, watchedKinds, disabledKinds map[string]bool, seen map[string]bool, config *CacheConfig) {
	// Find all .Get( and .List( call sites
	for _, idx := range getListRE.FindAllStringIndex(content, -1) {
		// Extract the balanced parenthesized call arguments
		parenStart := idx[1] - 1 // position of '('
		callBody := extractBlock(content[parenStart:], '(', ')')
		if callBody == "" {
			continue
		}

		// Find line number for the start of the call
		lineNo := strings.Count(content[:idx[0]], "\n") + 1

		// Look for &pkg.Type{} or &pkg.TypeList{} within the call arguments
		for _, re := range []*regexp.Regexp{clientGetTypeRE, clientListTypeRE} {
			for _, match := range re.FindAllStringSubmatch(callBody, -1) {
				pkg, kind := match[1], match[2]
				if watchedKinds[kind] || disabledKinds[kind] {
					continue
				}
				if ignoredTypes[kind] {
					continue
				}
				key := pkg + "." + kind
				if seen[key] {
					continue
				}
				seen[key] = true

				risk := "medium"
				if kind == "Secret" || kind == "ConfigMap" || kind == "Pod" {
					risk = "high"
				}

				config.ImplicitInformers = append(config.ImplicitInformers, ImplicitInformer{
					Type:   key,
					Source: fmt.Sprintf("%s:%d", relPath, lineNo),
					Risk:   risk,
				})
			}
		}
	}
}

// generateCacheIssues produces warning/info messages based on cache configuration.
func generateCacheIssues(config *CacheConfig, watches []ControllerWatch) {
	if len(watches) > 0 && config.ManagerFile != "" {
		// Check if any cache configuration exists at all
		if len(config.FilteredTypes) == 0 && len(config.DisabledTypes) == 0 && !config.DefaultTransform {
			config.Issues = append(config.Issues,
				"No cache configuration: all informers are cluster-wide (OOM risk). "+
					"See https://book.kubebuilder.io/reference/watching-resources/filtering for cache filtering patterns")
		}

		if !config.DefaultTransform && (len(config.FilteredTypes) > 0 || len(config.DisabledTypes) > 0) {
			config.Issues = append(config.Issues,
				"No DefaultTransform: managedFields cached for all objects (wasted memory). "+
					"Add cache.DefaultTransform to strip managedFields and reduce memory footprint")
		}
	}

	if config.GoMemLimit == "" && config.MemoryLimit != "" {
		config.Issues = append(config.Issues,
			"No GOMEMLIMIT set in deployment (Go GC cannot pressure-tune). "+
				"Set GOMEMLIMIT to 80-90% of container memory limit for optimal GC behavior")
	}

	// Enrich disabled type messages with context on why they were disabled
	for _, dt := range config.DisabledTypes {
		parts := strings.SplitN(dt, ".", 2)
		kind := dt
		if len(parts) == 2 {
			kind = parts[1]
		}
		if kind == "Secret" || kind == "ConfigMap" {
			config.Issues = append(config.Issues,
				fmt.Sprintf("Cache bypass (DisableFor) configured for %s. "+
					"This is a common fix for OOM caused by informer cache flooding from high-cardinality types "+
					"(e.g., opendatahub-io/model-registry-operator#457)", dt))
		}
	}

	for _, imp := range config.ImplicitInformers {
		if imp.Risk == "high" {
			config.Issues = append(config.Issues,
				fmt.Sprintf("Implicit informer for %s via client.Get at %s (cluster-wide, OOM risk). "+
					"This bypasses cache filters and creates a full cluster-wide watch", imp.Type, imp.Source))
		}
	}
}
