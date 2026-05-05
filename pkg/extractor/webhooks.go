package extractor

import (
	"log"
	"os"
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

// webhookPortRE matches webhook server port configuration in controller-runtime main.go:
//   webhook.NewServer(webhook.Options{Port: 9443})
//   webhook.Options{Port: 9443}
//   WebhookServer: webhook.NewServer(webhook.Options{Port: 9443})
var webhookPortRE = regexp.MustCompile(`webhook\.Options\{[^}]*Port:\s*(\d+)`)

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
					Source:         relativePath(repoPath, fpath),
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
				Source:        relativePath(repoPath, fpath),
			})
		}
	}

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

// toStringSlice is defined in rbac.go
