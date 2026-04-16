package renderer

import (
	"fmt"
	"strings"
)

// RBACRenderer renders RBAC relationships as a Mermaid graph.
type RBACRenderer struct{}

func (r *RBACRenderer) Filename() string { return "rbac.mmd" }

func (r *RBACRenderer) Render(data map[string]interface{}) string {
	component := getStr(data, "component", "unknown")
	rbac := getMap(data, "rbac")
	if rbac == nil || len(rbac) == 0 {
		return fmt.Sprintf("graph TD\n    note[No RBAC data for %s]", escapeLabel(component))
	}

	var b strings.Builder
	b.WriteString("graph TD\n")
	b.WriteString(fmt.Sprintf("    %%%% RBAC hierarchy for %s\n", component))
	b.WriteString("    classDef sa fill:#4a90d9,stroke:#2c5f8a,color:#fff\n")
	b.WriteString("    classDef role fill:#e8a838,stroke:#b07828,color:#fff\n")
	b.WriteString("    classDef resource fill:#5cb85c,stroke:#3d8b3d,color:#fff\n")
	b.WriteString("\n")

	counter := 0
	nextID := func(prefix string) string {
		counter++
		return fmt.Sprintf("%s_%d", prefix, counter)
	}

	// ClusterRoleBindings
	for _, binding := range getSlice(rbac, "cluster_role_bindings") {
		bindingName := getStr(binding, "name", "")
		roleRef := getStr(binding, "role_ref", "")
		bindingID := nextID("crb")
		roleID := sanitizeID("cr_" + roleRef)

		for _, subject := range getSlice(binding, "subjects") {
			subjKind := getStr(subject, "kind", "")
			subjName := getStr(subject, "name", "")
			subjNS := getStr(subject, "namespace", "")
			saLabel := subjName
			if subjNS != "" {
				saLabel = fmt.Sprintf("%s (%s)", subjName, subjNS)
			}
			saID := nextID("sa")
			b.WriteString(fmt.Sprintf("    %s[\"%s: %s\"] -->|bound via %s| %s[\"%s\"]\n",
				saID, escapeLabel(subjKind), escapeLabel(saLabel),
				escapeLabel(bindingName), bindingID, escapeLabel(bindingName)))
			b.WriteString(fmt.Sprintf("    class %s sa\n", saID))
		}

		b.WriteString(fmt.Sprintf("    %s -->|grants| %s[\"CR: %s\"]\n",
			bindingID, roleID, escapeLabel(roleRef)))
		b.WriteString(fmt.Sprintf("    class %s role\n", roleID))
	}

	// RoleBindings
	for _, binding := range getSlice(rbac, "role_bindings") {
		bindingName := getStr(binding, "name", "")
		roleRef := getStr(binding, "role_ref", "")
		bindingID := nextID("rb")
		roleID := sanitizeID("r_" + roleRef)

		for _, subject := range getSlice(binding, "subjects") {
			subjName := getStr(subject, "name", "")
			subjKind := getStr(subject, "kind", "")
			saID := nextID("sa")
			b.WriteString(fmt.Sprintf("    %s[\"%s: %s\"] -->|bound via %s| %s[\"%s\"]\n",
				saID, escapeLabel(subjKind), escapeLabel(subjName),
				escapeLabel(bindingName), bindingID, escapeLabel(bindingName)))
			b.WriteString(fmt.Sprintf("    class %s sa\n", saID))
		}

		b.WriteString(fmt.Sprintf("    %s -->|grants| %s[\"Role: %s\"]\n",
			bindingID, roleID, escapeLabel(roleRef)))
		b.WriteString(fmt.Sprintf("    class %s role\n", roleID))
	}

	// ClusterRole rules
	for _, role := range getSlice(rbac, "cluster_roles") {
		roleName := getStr(role, "name", "")
		roleID := sanitizeID("cr_" + roleName)
		for _, rule := range getSlice(role, "rules") {
			apiGroups := getStringSlice(rule, "apiGroups")
			resources := getStringSlice(rule, "resources")
			verbs := getStringSlice(rule, "verbs")
			for _, res := range resources {
				group := "core"
				if len(apiGroups) > 0 && apiGroups[0] != "" {
					group = apiGroups[0]
				}
				resID := nextID("res")
				verbStr := strings.Join(verbs, ", ")
				b.WriteString(fmt.Sprintf("    %s -->|%s| %s[\"%s: %s\"]\n",
					roleID, escapeLabel(verbStr), resID, escapeLabel(group), escapeLabel(res)))
				b.WriteString(fmt.Sprintf("    class %s resource\n", resID))
			}
		}
	}

	// Role rules
	for _, role := range getSlice(rbac, "roles") {
		roleName := getStr(role, "name", "")
		roleID := sanitizeID("r_" + roleName)
		for _, rule := range getSlice(role, "rules") {
			apiGroups := getStringSlice(rule, "apiGroups")
			resources := getStringSlice(rule, "resources")
			verbs := getStringSlice(rule, "verbs")
			for _, res := range resources {
				group := "core"
				if len(apiGroups) > 0 && apiGroups[0] != "" {
					group = apiGroups[0]
				}
				resID := nextID("res")
				verbStr := strings.Join(verbs, ", ")
				b.WriteString(fmt.Sprintf("    %s -->|%s| %s[\"%s: %s\"]\n",
					roleID, escapeLabel(verbStr), resID, escapeLabel(group), escapeLabel(res)))
				b.WriteString(fmt.Sprintf("    class %s resource\n", resID))
			}
		}
	}

	return strings.TrimRight(b.String(), "\n")
}
