package renderer

import (
	"fmt"
	"strings"
)

// DependencyRenderer renders a dependency graph in Mermaid.
type DependencyRenderer struct{}

func (r *DependencyRenderer) Filename() string { return "dependencies.mmd" }

func (r *DependencyRenderer) Render(data map[string]interface{}) string {
	component := getStr(data, "component", "unknown")
	deps := getMap(data, "dependencies")
	if deps == nil {
		deps = map[string]interface{}{}
	}
	goModules := getSlice(deps, "go_modules")
	internalODH := getSlice(deps, "internal_odh")

	var b strings.Builder
	b.WriteString("graph LR\n")
	b.WriteString(fmt.Sprintf("    %%%% Dependency graph for %s\n", component))
	b.WriteString("\n")
	b.WriteString("    classDef component fill:#3498db,stroke:#2980b9,color:#fff,stroke-width:3px\n")
	b.WriteString("    classDef internal fill:#2ecc71,stroke:#27ae60,color:#fff\n")
	b.WriteString("    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#333\n")
	b.WriteString("\n")

	compID := sanitizeID("comp_" + component)
	b.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", compID, escapeLabel(component)))
	b.WriteString(fmt.Sprintf("    class %s component\n", compID))
	b.WriteString("\n")

	counter := 0
	nextID := func(prefix string) string {
		counter++
		return fmt.Sprintf("%s_%d", prefix, counter)
	}

	// Internal ODH dependencies
	if len(internalODH) > 0 {
		b.WriteString("    %% Internal ODH dependencies\n")
		for _, odh := range internalODH {
			compName := getStr(odh, "component", "")
			interaction := getStr(odh, "interaction", "")
			depID := nextID("odh")
			b.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", depID, escapeLabel(compName)))
			b.WriteString(fmt.Sprintf("    class %s internal\n", depID))
			b.WriteString(fmt.Sprintf("    %s -->|\"%s\"| %s\n", compID, escapeLabel(interaction), depID))
		}
		b.WriteString("\n")
	}

	// Key external dependencies
	notablePrefixes := NotableDependencyPrefixes

	if len(goModules) > 0 {
		b.WriteString("    %% Key external dependencies\n")
		seen := make(map[string]bool)
		for _, mod := range goModules {
			module := getStr(mod, "module", "")
			version := getStr(mod, "version", "")

			// Skip internal RHOAI dependencies (any known org)
			if isInternalModule(module, data) {
				continue
			}

			// Only show notable dependencies
			notable := false
			for _, prefix := range notablePrefixes {
				if strings.HasPrefix(module, prefix) {
					notable = true
					break
				}
			}
			if !notable {
				continue
			}

			// Deduplicate by base module (first two path segments)
			parts := strings.SplitN(module, "/", 3)
			base := module
			if len(parts) >= 2 {
				base = parts[0] + "/" + parts[1]
			}
			if seen[base] {
				continue
			}
			seen[base] = true

			depID := nextID("ext")
			shortName := module
			if idx := strings.LastIndex(module, "/"); idx >= 0 {
				shortName = module[idx+1:]
			}
			b.WriteString(fmt.Sprintf("    %s[\"%s\\n%s\"]\n",
				depID, escapeLabel(shortName), escapeLabel(version)))
			b.WriteString(fmt.Sprintf("    class %s external\n", depID))
			b.WriteString(fmt.Sprintf("    %s --> %s\n", compID, depID))
		}
	}

	return strings.TrimRight(b.String(), "\n")
}

// isInternalModule checks if a Go module belongs to one of the known RHOAI orgs.
// It derives the org from the component data's "repo" field, and also checks
// known RHOAI organizations.
func isInternalModule(module string, data map[string]interface{}) bool {
	knownPrefixes := make([]string, len(KnownInternalPrefixes))
	copy(knownPrefixes, KnownInternalPrefixes)
	// Also derive from the repo field if present
	if repo := getStr(data, "repo", ""); repo != "" {
		parts := strings.SplitN(repo, "/", 2)
		if len(parts) >= 1 {
			prefix := "github.com/" + parts[0] + "/"
			found := false
			for _, kp := range knownPrefixes {
				if kp == prefix {
					found = true
					break
				}
			}
			if !found {
				knownPrefixes = append(knownPrefixes, prefix)
			}
		}
	}
	for _, prefix := range knownPrefixes {
		if strings.HasPrefix(module, prefix) {
			return true
		}
	}
	return false
}
