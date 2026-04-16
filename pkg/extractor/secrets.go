package extractor

import (
	"fmt"
	"strings"
)

var secretDeploymentPatterns = []string{
	"**/deployment.yaml",
	"**/deployment*.yaml",
	"**/statefulset.yaml",
	"**/manager*.yaml",
	"charts/**/templates/deployment*.yaml",
}

var secretServicePatterns = []string{
	"**/service.yaml",
	"**/service*.yaml",
}

// secretEntry tracks a secret reference during extraction.
type secretEntry struct {
	secretType    string
	referencedBy  []string
	provisionedBy string
}

// extractSecrets scans YAML files for references to Kubernetes Secrets (names
// and types only, never values).
func extractSecrets(repoPath string) []SecretRef {
	secretsMap := make(map[string]*secretEntry)

	// Scan deployments for secret references
	depFiles := findYAMLFiles(repoPath, secretDeploymentPatterns)
	for _, fpath := range depFiles {
		for _, doc := range parseYAMLSafe(fpath) {
			kind, _ := doc["kind"].(string)
			if kind != "Deployment" && kind != "StatefulSet" {
				continue
			}
			metadata, _ := doc["metadata"].(map[string]interface{})
			if metadata == nil {
				metadata = map[string]interface{}{}
			}
			depName, _ := metadata["name"].(string)
			refLabel := fmt.Sprintf("%s/%s", strings.ToLower(kind), depName)

			spec, _ := doc["spec"].(map[string]interface{})
			if spec == nil {
				continue
			}
			template, _ := spec["template"].(map[string]interface{})
			if template == nil {
				continue
			}
			podSpec, _ := template["spec"].(map[string]interface{})
			if podSpec == nil {
				continue
			}

			// Volumes with secrets
			for _, vol := range toSliceOfMaps(podSpec["volumes"]) {
				if secret, ok := vol["secret"].(map[string]interface{}); ok {
					if secretName, ok := secret["secretName"].(string); ok {
						addSecret(secretsMap, secretName, refLabel, "volume-mounted", "Opaque")
					}
				}
			}

			// Containers env and envFrom
			allContainers := append(
				toSliceOfMaps(podSpec["containers"]),
				toSliceOfMaps(podSpec["initContainers"])...,
			)
			for _, container := range allContainers {
				// env[].valueFrom.secretKeyRef
				for _, envVar := range toSliceOfMaps(container["env"]) {
					valueFrom, ok := envVar["valueFrom"].(map[string]interface{})
					if !ok {
						continue
					}
					if secretRef, ok := valueFrom["secretKeyRef"].(map[string]interface{}); ok {
						if name, ok := secretRef["name"].(string); ok {
							addSecret(secretsMap, name, refLabel, "env-var", "Opaque")
						}
					}
				}
				// envFrom[].secretRef
				for _, envFrom := range toSliceOfMaps(container["envFrom"]) {
					if secretRef, ok := envFrom["secretRef"].(map[string]interface{}); ok {
						if name, ok := secretRef["name"].(string); ok {
							addSecret(secretsMap, name, refLabel, "envFrom", "Opaque")
						}
					}
				}
			}
		}
	}

	// Scan services for cert annotations
	svcFiles := findYAMLFiles(repoPath, secretServicePatterns)
	for _, fpath := range svcFiles {
		for _, doc := range parseYAMLSafe(fpath) {
			if kind, _ := doc["kind"].(string); kind != "Service" {
				continue
			}
			metadata, _ := doc["metadata"].(map[string]interface{})
			if metadata == nil {
				continue
			}
			annotations, _ := metadata["annotations"].(map[string]interface{})
			if annotations == nil {
				continue
			}
			certSecret, ok := annotations["service.beta.openshift.io/serving-cert-secret-name"].(string)
			if !ok {
				continue
			}
			svcName, _ := metadata["name"].(string)
			addSecret(secretsMap, certSecret, fmt.Sprintf("service/%s", svcName), "OpenShift serving cert", "kubernetes.io/tls")
		}
	}

	var refs []SecretRef
	for name, entry := range secretsMap {
		refs = append(refs, SecretRef{
			Name:          name,
			Type:          entry.secretType,
			ReferencedBy:  entry.referencedBy,
			ProvisionedBy: entry.provisionedBy,
		})
	}
	if refs == nil {
		refs = []SecretRef{}
	}
	return refs
}

// addSecret adds or updates a secret reference in the map.
func addSecret(m map[string]*secretEntry, name, referencedBy, provisionedBy, secretType string) {
	entry, exists := m[name]
	if !exists {
		entry = &secretEntry{
			secretType:    secretType,
			provisionedBy: provisionedBy,
		}
		m[name] = entry
	}
	// Check for duplicate referencedBy
	for _, ref := range entry.referencedBy {
		if ref == referencedBy {
			return
		}
	}
	entry.referencedBy = append(entry.referencedBy, referencedBy)
	// Upgrade type if more specific
	if secretType != "Opaque" {
		entry.secretType = secretType
	}
}
