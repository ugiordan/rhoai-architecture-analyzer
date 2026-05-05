package extractor

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var rbacYAMLPatterns = []string{
	"config/rbac/*.yaml",
	"config/rbac/*.yml",
	"charts/**/templates/*role*.yaml",
	"charts/**/templates/*role*.yml",
	"deploy/rbac/*.yaml",
	"deploy/rbac/*.yml",
	"manifests/**/*role*.yaml",
	"manifests/**/*role*.yml",
}

var kubebuilderRBACRE = regexp.MustCompile(`//\+kubebuilder:rbac:(.+)`)

// extractRBAC scans YAML files for RBAC resources and Go files for kubebuilder
// RBAC markers.
func extractRBAC(repoPath string) *RBACData {
	files := findYAMLFiles(repoPath, rbacYAMLPatterns)

	var clusterRoles []RBACRole
	var clusterRoleBindings []RBACBinding
	var roles []RBACRole
	var roleBindings []RBACBinding

	for _, fpath := range files {
		for _, doc := range parseYAMLSafe(fpath) {
			kind, _ := doc["kind"].(string)
			metadata, _ := doc["metadata"].(map[string]interface{})
			if metadata == nil {
				metadata = map[string]interface{}{}
			}
			name, _ := metadata["name"].(string)
			source := relativePath(repoPath, fpath)

			switch kind {
			case "ClusterRole":
				clusterRoles = append(clusterRoles, RBACRole{
					Name:            name,
					Source:          source,
					Rules:           extractRBACRules(doc),
					AggregationRule: extractAggregationRule(doc),
				})
			case "ClusterRoleBinding":
				roleRef, _ := doc["roleRef"].(map[string]interface{})
				refName := ""
				if roleRef != nil {
					refName, _ = roleRef["name"].(string)
				}
				clusterRoleBindings = append(clusterRoleBindings, RBACBinding{
					Name:     name,
					RoleRef:  refName,
					Subjects: extractRBACSubjects(doc),
					Source:   source,
				})
			case "Role":
				roles = append(roles, RBACRole{
					Name:   name,
					Source: source,
					Rules:  extractRBACRules(doc),
				})
			case "RoleBinding":
				roleRef, _ := doc["roleRef"].(map[string]interface{})
				refName := ""
				if roleRef != nil {
					refName, _ = roleRef["name"].(string)
				}
				roleBindings = append(roleBindings, RBACBinding{
					Name:     name,
					RoleRef:  refName,
					Subjects: extractRBACSubjects(doc),
					Source:   source,
				})
			}
		}
	}

	// Scan Go files for kubebuilder RBAC markers
	goFiles := findGoFiles(repoPath)
	var markers []RBACMarker

	for _, fpath := range goFiles {
		data, err := os.ReadFile(fpath)
		if err != nil {
			continue
		}
		lines := strings.Split(string(data), "\n")
		for lineNo, line := range lines {
			match := kubebuilderRBACRE.FindStringSubmatch(line)
			if match == nil {
				continue
			}
			markerText := strings.TrimLeft(match[0], "/")
			markers = append(markers, RBACMarker{
				File:   relativePath(repoPath, fpath),
				Line:   lineNo + 1,
				Marker: markerText,
				Parsed: parseKubebuilderMarker(match[1]),
			})
		}
	}

	if clusterRoles == nil {
		clusterRoles = []RBACRole{}
	}
	if clusterRoleBindings == nil {
		clusterRoleBindings = []RBACBinding{}
	}
	if roles == nil {
		roles = []RBACRole{}
	}
	if roleBindings == nil {
		roleBindings = []RBACBinding{}
	}
	if markers == nil {
		markers = []RBACMarker{}
	}

	return &RBACData{
		ClusterRoles:        clusterRoles,
		ClusterRoleBindings: clusterRoleBindings,
		Roles:               roles,
		RoleBindings:        roleBindings,
		KubebuilderMarkers:  markers,
	}
}

// findGoFiles finds all .go files under the repo path, skipping non-source directories.
func findGoFiles(repoPath string) []string {
	var files []string
	_ = filepath.WalkDir(repoPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() && isExcludedDir(d.Name(), nil) {
			return filepath.SkipDir
		}
		if !d.IsDir() && strings.HasSuffix(path, ".go") {
			files = append(files, path)
		}
		return nil
	})
	return files
}

// parseKubebuilderMarker parses key=value pairs from a kubebuilder RBAC marker.
// Multi-value fields use semicolons (e.g. verbs=get;list;watch).
// Limitation: commas inside quoted values are not supported. Standard
// kubebuilder markers never use this format, so this is safe in practice.
func parseKubebuilderMarker(body string) map[string]interface{} {
	result := make(map[string]interface{})
	for _, part := range strings.Split(body, ",") {
		part = strings.TrimSpace(part)
		if idx := strings.Index(part, "="); idx >= 0 {
			key := strings.TrimSpace(part[:idx])
			value := strings.TrimSpace(part[idx+1:])
			if strings.Contains(value, ";") {
				result[key] = strings.Split(value, ";")
			} else {
				result[key] = []string{value}
			}
		}
	}
	return result
}

// extractRBACRules extracts RBAC rules from a Role/ClusterRole document.
func extractRBACRules(doc map[string]interface{}) []RBACRule {
	rawRules, _ := doc["rules"].([]interface{})
	var rules []RBACRule
	for _, r := range rawRules {
		rule, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		rules = append(rules, RBACRule{
			APIGroups:     toStringSlice(rule["apiGroups"]),
			Resources:     toStringSlice(rule["resources"]),
			Verbs:         toStringSlice(rule["verbs"]),
			ResourceNames: toStringSlice(rule["resourceNames"]),
		})
	}
	if rules == nil {
		rules = []RBACRule{}
	}
	return rules
}

// extractRBACSubjects extracts subjects from a RoleBinding/ClusterRoleBinding.
func extractRBACSubjects(doc map[string]interface{}) []RBACSubject {
	raw, _ := doc["subjects"].([]interface{})
	var subjects []RBACSubject
	for _, s := range raw {
		subj, ok := s.(map[string]interface{})
		if !ok {
			continue
		}
		kind, _ := subj["kind"].(string)
		name, _ := subj["name"].(string)
		namespace, _ := subj["namespace"].(string)
		subjects = append(subjects, RBACSubject{
			Kind:      kind,
			Name:      name,
			Namespace: namespace,
		})
	}
	if subjects == nil {
		subjects = []RBACSubject{}
	}
	return subjects
}

// extractAggregationRule extracts the aggregationRule.clusterRoleSelectors
// matchLabels from a ClusterRole. Returns nil if no aggregation rule is present.
func extractAggregationRule(doc map[string]interface{}) map[string]string {
	aggRule, ok := doc["aggregationRule"].(map[string]interface{})
	if !ok {
		return nil
	}
	selectors, ok := aggRule["clusterRoleSelectors"].([]interface{})
	if !ok || len(selectors) == 0 {
		return nil
	}
	labels := make(map[string]string)
	for _, sel := range selectors {
		selMap, ok := sel.(map[string]interface{})
		if !ok {
			continue
		}
		matchLabels, ok := selMap["matchLabels"].(map[string]interface{})
		if !ok {
			continue
		}
		for k, v := range matchLabels {
			if vs, ok := v.(string); ok {
				labels[k] = vs
			}
		}
	}
	if len(labels) == 0 {
		return nil
	}
	return labels
}

// toStringSlice converts an interface{} (expected []interface{} of strings)
// to []string.
func toStringSlice(v interface{}) []string {
	items, ok := v.([]interface{})
	if !ok {
		return []string{}
	}
	result := make([]string, 0, len(items))
	for _, item := range items {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}
	return result
}
