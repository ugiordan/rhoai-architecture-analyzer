package extractor

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var webhookYAMLPatterns = []string{
	"config/webhook/*.yaml",
	"config/webhook/*.yml",
	"**/webhook*.yaml",
	"**/webhook*.yml",
	"**/mutating*.yaml",
	"**/validating*.yaml",
	"**/mutating*.yaml.tmpl",
	"**/validating*.yaml.tmpl",
	"**/webhook*.yaml.tmpl",
	"charts/**/templates/*webhook*.yaml",
	"manifests/**/webhook*.yaml",
}

var kubebuilderWebhookRE = regexp.MustCompile(`//\+kubebuilder:webhook:(.+)`)

// hookServerRegisterRE matches webhook server registration calls like:
//   hookServer.Register("/path", ...)
//   GetWebhookServer().Register("/path", ...)
var hookServerRegisterRE = regexp.MustCompile(`(?:hookServer|WebhookServer\(\))\.Register\("(/[^"]+)"`)

// clientGetListRE matches client.Get/List, r.Get/List, r.Client.Get/List, r.client.Get/List
// with a type argument like &corev1.Secret{} or &servingv1beta1.InferenceService{}
var clientGetListRE = regexp.MustCompile(`\.\s*(?:Get|List)\s*\([^)]*&(\w+)\.(\w+?)(?:List)?\{\}`)

// packageGroupMap maps common package aliases to their Kubernetes API groups.
var packageGroupMap = map[string]string{
	"v1":                 "",
	"corev1":             "",
	"appsv1":             "apps",
	"batchv1":            "batch",
	"networkingv1":       "networking.k8s.io",
	"rbacv1":             "rbac.authorization.k8s.io",
	"storagev1":          "storage.k8s.io",
	"policyv1":           "policy",
	"autoscalingv1":      "autoscaling",
	"autoscalingv2":      "autoscaling",
	"coordinationv1":     "coordination.k8s.io",
	"admissionv1":        "admissionregistration.k8s.io",
	"apiextensionsv1":    "apiextensions.k8s.io",
	"servingv1":          "serving.knative.dev",
	"servingv1beta1":     "serving.kserve.io",
	"servingv1alpha1":    "serving.kserve.io",
	"inferencev1alpha1":  "inference.networking.x-k8s.io",
}

// skipPackages lists package aliases that should be ignored for data_read extraction.
var skipPackages = map[string]bool{
	"metav1":        true,
	"types":         true,
	"ctrl":          true,
	"unstructured":  true,
}

// handlerFileGlobs lists patterns for Go files that may contain webhook handlers.
var handlerFileGlobs = []string{
	"**/*_webhook.go",
	"**/webhook*.go",
	"**/webhooks.go",
	"**/*_validation.go",
	"**/*_defaults.go",
	"pkg/webhook/**/*.go",
	"internal/webhook/**/*.go",
}

// webhookPortRE matches webhook server port configuration in controller-runtime main.go:
//   webhook.NewServer(webhook.Options{Port: 9443})
//   webhook.Options{Port: 9443}
//   WebhookServer: webhook.NewServer(webhook.Options{Port: 9443})
var webhookPortRE = regexp.MustCompile(`webhook\.Options\{[^}]*Port:\s*(\d+)`)

// webhookSetupFuncRE matches functions that set up webhooks.
var webhookSetupFuncRE = regexp.MustCompile(`func\s+(?:\([^)]+\)\s+)?(SetupWithManager|setupWebhooks?|registerWebhooks?|SetupWebhookWithManager)\s*\(`)

// enableConditionPatterns lists regexes for webhook guard conditions.
var enableConditionPatterns = []*regexp.Regexp{
	regexp.MustCompile(`os\.Getenv\([^)]+\)\s*[!=]=\s*"[^"]*"`),
	regexp.MustCompile(`[fF]eature[Gg]ate\.Enabled\([^)]+\)`),
	regexp.MustCompile(`\.ManagementState\s*==\s*\S+`),
	regexp.MustCompile(`(?:^|[^.\w])!?\s*(?:enable|disable)[Ww]ebhooks?\s*[{&|]`),
}

// webhookCallRE matches webhook registration calls inside functions.
var webhookCallRE = regexp.MustCompile(`(?:WebhookManagedBy|NewWebhookManagedBy|WithValidator|WithDefaulter|\.Register\("/)`)

// extractWebhooks finds ValidatingWebhookConfiguration and
// MutatingWebhookConfiguration resources, plus kubebuilder webhook markers.
func extractWebhooks(repoPath string) []WebhookConfig {
	var webhooks []WebhookConfig

	// Extract from YAML manifests (including Go template files)
	files := findYAMLFiles(repoPath, webhookYAMLPatterns)
	for _, fpath := range files {
		var docs []map[string]interface{}
		if strings.HasSuffix(fpath, ".tmpl") {
			docs = parseTemplateYAML(fpath)
		} else {
			docs = parseYAMLSafe(fpath)
		}
		for _, doc := range docs {
			kind, _ := doc["kind"].(string)
			if kind != "ValidatingWebhookConfiguration" && kind != "MutatingWebhookConfiguration" {
				continue
			}

			whType := "validating"
			if kind == "MutatingWebhookConfiguration" {
				whType = "mutating"
			}

			metadata, _ := doc["metadata"].(map[string]interface{})
			name := ""
			if metadata != nil {
				name, _ = metadata["name"].(string)
			}

			webhooksField, _ := doc["webhooks"].([]interface{})
			for _, wh := range webhooksField {
				whMap, ok := wh.(map[string]interface{})
				if !ok {
					continue
				}

				whName, _ := whMap["name"].(string)
				if whName == "" {
					whName = name
				}

				failurePolicy, _ := whMap["failurePolicy"].(string)
				sideEffects, _ := whMap["sideEffects"].(string)
				timeoutSeconds := 0
				if ts, ok := whMap["timeoutSeconds"].(float64); ok {
					timeoutSeconds = int(ts)
				} else if ts, ok := whMap["timeoutSeconds"].(int); ok {
					timeoutSeconds = ts
				}

				// Extract service ref and path from clientConfig
				serviceRef := ""
				path := ""
				if clientConfig, ok := whMap["clientConfig"].(map[string]interface{}); ok {
					if svc, ok := clientConfig["service"].(map[string]interface{}); ok {
						svcName, _ := svc["name"].(string)
						svcNs, _ := svc["namespace"].(string)
						if svcName != "" {
							serviceRef = svcNs + "/" + svcName
						}
						p, _ := svc["path"].(string)
						path = p
					}
				}

				// Extract rules
				var rules []WebhookRule
				rulesField, _ := whMap["rules"].([]interface{})
				for _, r := range rulesField {
					rMap, ok := r.(map[string]interface{})
					if !ok {
						continue
					}
					rules = append(rules, WebhookRule{
						APIGroups:   toStringSlice(rMap["apiGroups"]),
						APIVersions: toStringSlice(rMap["apiVersions"]),
						Resources:   toStringSlice(rMap["resources"]),
						Operations:  toStringSlice(rMap["operations"]),
					})
				}

				if rules == nil {
					rules = []WebhookRule{}
				}

				webhooks = append(webhooks, WebhookConfig{
					Name:           whName,
					Type:           whType,
					ServiceRef:     serviceRef,
					Path:           path,
					FailurePolicy:  failurePolicy,
					SideEffects:    sideEffects,
					TimeoutSeconds: timeoutSeconds,
					Rules:          rules,
					Sources:        []SourceRef{{Type: "yaml_manifest", File: relativePath(repoPath, fpath)}},
				})
			}
		}
	}

	// Extract from kubebuilder markers in Go files
	goFiles := findFiles(repoPath, []string{"**/*_webhook.go", "**/*_webhook_test.go", "**/webhook.go", "**/webhooks.go"})
	for _, fpath := range goFiles {
		data, err := os.ReadFile(fpath)
		if err != nil {
			log.Printf("warning: skipping %s: %v", fpath, err)
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			match := kubebuilderWebhookRE.FindStringSubmatch(strings.TrimSpace(line))
			if match == nil {
				continue
			}
			attrs := parseKubebuilderAttrs(match[1])
			whType := "validating"
			if v, ok := attrs["mutating"]; ok && v == "true" {
				whType = "mutating"
			}
			webhooks = append(webhooks, WebhookConfig{
				Name:          attrs["name"],
				Type:          whType,
				Path:          attrs["path"],
				FailurePolicy: attrs["failurePolicy"],
				SideEffects:   attrs["sideEffects"],
				Rules:         []WebhookRule{},
				Sources:       []SourceRef{{Type: "kubebuilder_marker", File: relativePath(repoPath, fpath)}},
			})
		}
	}

	// Extract conversion webhooks from CRD definitions
	webhooks = append(webhooks, extractConversionWebhooks(repoPath)...)

	if webhooks == nil {
		webhooks = []WebhookConfig{}
	}
	return webhooks
}

// extractWebhookServerPort scans main.go and cmd/ for the webhook server port
// configured via controller-runtime's webhook.Options{Port: N}.
// Returns 0 if not found.
func extractWebhookServerPort(repoPath string) int {
	mainFiles := findFiles(repoPath, []string{"main.go", "cmd/**/main.go", "cmd/**/*.go"})
	for _, fpath := range mainFiles {
		data, err := os.ReadFile(fpath)
		if err != nil {
			continue
		}
		match := webhookPortRE.FindSubmatch(data)
		if match != nil {
			port := 0
			for _, b := range match[1] {
				port = port*10 + int(b-'0')
			}
			if port > 0 {
				return port
			}
		}
	}
	return 0
}

func parseKubebuilderAttrs(s string) map[string]string {
	attrs := make(map[string]string)
	for _, part := range strings.Split(s, ",") {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) == 2 {
			attrs[kv[0]] = kv[1]
		}
	}
	return attrs
}

// extractConversionWebhooks scans CRD YAML files for conversion webhooks.
// Returns WebhookConfig entries with Type="conversion" for CRDs that use
// spec.conversion.strategy=Webhook.
func extractConversionWebhooks(repoPath string) []WebhookConfig {
	var webhooks []WebhookConfig

	// Reuse CRD search patterns from crds.go to stay in sync
	files := findYAMLFiles(repoPath, crdSearchPatterns)
	for _, fpath := range files {
		docs := parseYAMLSafe(fpath)
		for _, doc := range docs {
			kind, _ := doc["kind"].(string)
			if kind != "CustomResourceDefinition" {
				continue
			}

			spec, ok := doc["spec"].(map[string]interface{})
			if !ok {
				continue
			}

			// Check if conversion strategy is Webhook
			conversion, ok := spec["conversion"].(map[string]interface{})
			if !ok {
				continue
			}
			strategy, _ := conversion["strategy"].(string)
			if strategy != "Webhook" {
				continue
			}

			// Extract CRD kind
			names, _ := spec["names"].(map[string]interface{})
			crdKind := ""
			if names != nil {
				crdKind, _ = names["kind"].(string)
			}

			// Extract service ref and path from webhook.clientConfig.service
			serviceRef := ""
			path := ""
			if webhook, ok := conversion["webhook"].(map[string]interface{}); ok {
				if clientConfig, ok := webhook["clientConfig"].(map[string]interface{}); ok {
					if svc, ok := clientConfig["service"].(map[string]interface{}); ok {
						svcName, _ := svc["name"].(string)
						svcNs, _ := svc["namespace"].(string)
						if svcName != "" {
							serviceRef = svcNs + "/" + svcName
						}
						p, _ := svc["path"].(string)
						path = p
					}
				}
			}

			name := "conversion-" + strings.ToLower(crdKind)
			if crdKind == "" {
				name = "conversion-unknown"
			}
			webhooks = append(webhooks, WebhookConfig{
				Name:          name,
				Type:          "conversion",
				ConversionCRD: crdKind,
				ServiceRef:    serviceRef,
				Path:          path,
				Rules:         []WebhookRule{},
				Sources:       []SourceRef{{Type: "yaml_manifest", File: relativePath(repoPath, fpath)}},
			})
		}
	}

	if webhooks == nil {
		webhooks = []WebhookConfig{}
	}
	return webhooks
}

// enrichWebhooks runs all enrichment passes on the webhook slice in place.
// Called after YAML/marker extraction and kustomize merge.
func enrichWebhooks(webhooks []WebhookConfig, repoRoot string) {
	mapGoHandlers(webhooks, repoRoot)
	extractDataRead(webhooks, repoRoot)
	extractEnableConditions(webhooks, repoRoot)
}

// mapGoHandlers correlates webhooks with their Go handler source files.
//
// Strategy A (primary): if a webhook already has a kubebuilder_marker source,
// the same Go file is the handler. A go_handler SourceRef is added.
//
// Strategy B (fallback): for webhooks not matched by Strategy A, scan Go files
// for hookServer.Register("/path", ...) literals and match on path.
func mapGoHandlers(webhooks []WebhookConfig, repoRoot string) {
	// Build path -> indices map for Strategy B lookups.
	pathIndex := make(map[string][]int)
	for i := range webhooks {
		if webhooks[i].Path != "" {
			pathIndex[webhooks[i].Path] = append(pathIndex[webhooks[i].Path], i)
		}
	}

	// Track which indices were handled by Strategy A so B can skip them.
	handledByA := make(map[int]bool)

	// Strategy A: kubebuilder marker correlation.
	for i := range webhooks {
		for _, src := range webhooks[i].Sources {
			if src.Type == "kubebuilder_marker" {
				if !hasSource(webhooks[i].Sources, "go_handler", src.File) {
					webhooks[i].Sources = append(webhooks[i].Sources, SourceRef{
						Type: "go_handler",
						File: src.File,
					})
				}
				handledByA[i] = true
				break
			}
		}
	}

	// Strategy B: hookServer.Register scanning.
	goFiles := findFiles(repoRoot, handlerFileGlobs)
	for _, fpath := range goFiles {
		// Skip test files.
		if strings.HasSuffix(fpath, "_test.go") {
			continue
		}

		data, err := os.ReadFile(fpath)
		if err != nil {
			log.Printf("warning: skipping %s: %v", fpath, err)
			continue
		}

		matches := hookServerRegisterRE.FindAllStringSubmatch(string(data), -1)
		if len(matches) == 0 {
			continue
		}

		relFile := relativePath(repoRoot, fpath)
		for _, m := range matches {
			registeredPath := m[1]
			indices, ok := pathIndex[registeredPath]
			if !ok {
				continue
			}
			for _, idx := range indices {
				if handledByA[idx] {
					continue
				}
				if !hasSource(webhooks[idx].Sources, "go_handler", relFile) {
					webhooks[idx].Sources = append(webhooks[idx].Sources, SourceRef{
						Type: "go_handler",
						File: relFile,
					})
				}
			}
		}
	}
}

// hasSource checks whether a source with the given type and file already exists.
func hasSource(sources []SourceRef, srcType, file string) bool {
	for _, s := range sources {
		if s.Type == srcType && s.File == file {
			return true
		}
	}
	return false
}

// extractDataRead scans Go handler files for client.Get/List calls to determine
// which Kubernetes types each webhook reads.
func extractDataRead(webhooks []WebhookConfig, repoRoot string) {
	for i := range webhooks {
		// Collect all go_handler source files for this webhook.
		var handlerFiles []string
		for _, src := range webhooks[i].Sources {
			if src.Type == "go_handler" {
				handlerFiles = append(handlerFiles, src.File)
			}
		}
		if len(handlerFiles) == 0 {
			continue
		}

		// Deduplicate by "pkg.Kind" key.
		seen := make(map[string]bool)
		var dataRead []TypeRef

		for _, relFile := range handlerFiles {
			fullPath := filepath.Join(repoRoot, relFile)
			data, err := os.ReadFile(fullPath)
			if err != nil {
				log.Printf("warning: skipping %s: %v", fullPath, err)
				continue
			}

			for _, line := range strings.Split(string(data), "\n") {
				// Strip // comments.
				if idx := strings.Index(line, "//"); idx != -1 {
					line = line[:idx]
				}

				matches := clientGetListRE.FindAllStringSubmatch(line, -1)
				for _, m := range matches {
					pkg := m[1]
					kind := m[2]

					// Skip known non-resource packages.
					if skipPackages[pkg] {
						continue
					}

					// Dedup key.
					key := pkg + "." + kind
					if seen[key] {
						continue
					}
					seen[key] = true

					// Look up group from packageGroupMap.
					group, known := packageGroupMap[pkg]
					dataRead = append(dataRead, TypeRef{
						Kind:       kind,
						Group:      group,
						GroupKnown: known,
					})
				}
			}
		}

		if len(dataRead) > 0 {
			webhooks[i].DataRead = dataRead
		}
	}
}

// extractEnableConditions scans webhook setup functions for if-blocks guarding
// webhook registration, extracting enable conditions from env vars, feature gates,
// or management state checks.
func extractEnableConditions(webhooks []WebhookConfig, repoRoot string) {
	// Find Go files that may contain webhook setup functions.
	patterns := []string{
		"main.go",
		"cmd/**/main.go",
		"cmd/**/*.go",
		"**/*_webhook.go",
		"**/webhook*.go",
		"**/webhooks.go",
		"internal/**/*.go",
		"pkg/**/*.go",
	}
	files := findFiles(repoRoot, patterns)

	// Collect all unique conditions found across all webhook setup files.
	var conditions []string
	conditionSet := make(map[string]bool)

	for _, fpath := range files {
		// Skip test files.
		if strings.HasSuffix(filepath.Base(fpath), "_test.go") {
			continue
		}

		data, err := os.ReadFile(fpath)
		if err != nil {
			continue
		}

		content := string(data)

		// Check if this file contains both a webhook setup function AND webhook registration calls.
		if !webhookSetupFuncRE.MatchString(content) || !webhookCallRE.MatchString(content) {
			continue
		}

		// Scan for if-blocks and try to match condition patterns.
		for _, line := range strings.Split(content, "\n") {
			trimmed := strings.TrimSpace(line)
			if !strings.HasPrefix(trimmed, "if ") && !strings.HasPrefix(trimmed, "if\t") {
				continue
			}

			// Try each condition pattern.
			for _, pattern := range enableConditionPatterns {
				if match := pattern.FindString(trimmed); match != "" {
					match = strings.TrimSpace(match)
					match = strings.TrimRight(match, "{&| ")
					if !conditionSet[match] {
						conditionSet[match] = true
						conditions = append(conditions, match)
					}
				}
			}
		}
	}

	// Apply the collected conditions to all webhooks that don't already have one.
	if len(conditions) == 0 {
		return
	}

	combinedCondition := strings.Join(conditions, " && ")
	for i := range webhooks {
		if webhooks[i].EnableCondition == "" {
			webhooks[i].EnableCondition = combinedCondition
		}
	}
}

// mergeWebhookBehavior merges Go AST-extracted webhook mutations and validations
// into the existing webhook configs. Matches by path. For webhooks found in the
// AST but not in YAML/markers, a new entry is appended with a "go_ast" source.
func mergeWebhookBehavior(arch *ComponentArchitecture, behaviors map[string]WebhookBehavior) {
	if len(behaviors) == 0 {
		return
	}
	pathIndex := make(map[string]int)
	for i, wh := range arch.Webhooks {
		if wh.Path != "" {
			pathIndex[wh.Path] = i
		}
	}
	for path, b := range behaviors {
		if idx, exists := pathIndex[path]; exists {
			arch.Webhooks[idx].Mutations = b.Mutations
			arch.Webhooks[idx].Validations = b.Validations
			arch.Webhooks[idx].TargetType = b.TargetType
		} else {
			whType := "validating"
			if strings.Contains(path, "mutate") {
				whType = "mutating"
			}
			arch.Webhooks = append(arch.Webhooks, WebhookConfig{
				Name:        b.TargetType + "-webhook",
				Type:        whType,
				Path:        path,
				Mutations:   b.Mutations,
				Validations: b.Validations,
				TargetType:  b.TargetType,
				Sources:     []SourceRef{{Type: "go_ast"}},
			})
		}
	}
}

// toStringSlice is defined in rbac.go
