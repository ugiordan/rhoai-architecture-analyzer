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
	"charts/**/templates/*webhook*.yaml",
	"manifests/**/webhook*.yaml",
}

var kubebuilderWebhookRE = regexp.MustCompile(`//\+kubebuilder:webhook:(.+)`)

// extractWebhooks finds ValidatingWebhookConfiguration and
// MutatingWebhookConfiguration resources, plus kubebuilder webhook markers.
func extractWebhooks(repoPath string) []WebhookConfig {
	var webhooks []WebhookConfig

	// Extract from YAML manifests
	files := findYAMLFiles(repoPath, webhookYAMLPatterns)
	for _, fpath := range files {
		for _, doc := range parseYAMLSafe(fpath) {
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
					Name:          whName,
					Type:          whType,
					ServiceRef:    serviceRef,
					Path:          path,
					FailurePolicy: failurePolicy,
					Rules:         rules,
					Source:        relativePath(repoPath, fpath),
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
