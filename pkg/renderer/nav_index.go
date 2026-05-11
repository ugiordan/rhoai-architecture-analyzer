package renderer

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// SanitizeFilename converts a component name to a safe filename.
// Lowercase, replace non-[a-z0-9-_.] with dash, collapse consecutive dashes,
// strip leading/trailing dashes.
func SanitizeFilename(name string) string {
	if name == "" {
		return "unknown"
	}
	name = strings.ToLower(name)
	var b strings.Builder
	prevDash := false
	for _, ch := range name {
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' || ch == '.' {
			b.WriteRune(ch)
			prevDash = false
		} else {
			if !prevDash {
				b.WriteByte('-')
			}
			prevDash = true
		}
	}
	result := strings.Trim(b.String(), "-")
	if result == "" {
		return "unknown"
	}
	return result
}

// RenderNavIndex produces INDEX.md with a query routing table and component
// stats table from aggregated platform data.
func RenderNavIndex(platformData map[string]interface{}) string {
	var b strings.Builder

	b.WriteString(autoGenHeader)
	b.WriteString("# Architecture Analyzer Output\n\n")

	platform := getStr(platformData, "platform", "")
	if platform != "" {
		b.WriteString(fmt.Sprintf("**Platform**: %s  \n", platform))
	}
	aggAt := getStr(platformData, "aggregated_at", "")
	if aggAt != "" {
		b.WriteString(fmt.Sprintf("**Generated**: %s  \n", aggAt))
	}

	componentData := getSlice(platformData, "component_data")
	b.WriteString(fmt.Sprintf("**Components**: %d  \n\n", len(componentData)))

	// Query routing table
	b.WriteString("## How to Find Information\n\n")
	b.WriteString("| Question Type | Where to Look |\n")
	b.WriteString("|---------------|---------------|\n")
	b.WriteString("| \"What CRDs does X manage?\" | `components/<name>.md` -> CRDs section |\n")
	b.WriteString("| \"What ports does X use?\" | `components/<name>.md` -> Services section |\n")
	b.WriteString("| \"How does X interact with Y?\" | `cross-component/interactions.md` or `components/<name>.md` -> Interactions section |\n")
	b.WriteString("| \"Where is the X component doc?\" | `components/<name>.md` |\n")
	b.WriteString("| \"What dependencies does X have?\" | `components/<name>.md` -> Dependencies section |\n")
	b.WriteString("| \"What RBAC does X need?\" | `components/<name>.md` -> RBAC section |\n")
	b.WriteString("| Platform-level overview | `PLATFORM.md` |\n")
	b.WriteString("| Full structured data | `platform-architecture.json` |\n")
	b.WriteString("| Detailed component data (JSON) | `<name>/component-architecture.json` |\n")
	b.WriteString("| Code property graph (functions, calls, file paths) | `<name>/code-graph.json` |\n")
	b.WriteString("| Cross-component interactions (deep dive) | Grep `<name>/code-graph.json` for other component names in file paths |\n\n")

	// Per-component raw data section
	b.WriteString("## Per-Component Raw Data\n\n")
	b.WriteString("Each component has a subdirectory with structured JSON extracted by static analysis:\n\n")
	b.WriteString("```\n")
	b.WriteString("<component-name>/\n")
	b.WriteString("  component-architecture.json   # Full extracted data: CRDs, services, RBAC,\n")
	b.WriteString("                                # deployments, operator config, reconcile sequences,\n")
	b.WriteString("                                # controller watches, and more\n")
	b.WriteString("  code-graph.json               # Code property graph: functions, classes, calls,\n")
	b.WriteString("                                # file paths. Use file paths to discover cross-\n")
	b.WriteString("                                # component integrations (e.g., a path like\n")
	b.WriteString("                                # providers/remote/inference/vllm/ indicates\n")
	b.WriteString("                                # integration with the vllm component)\n")
	b.WriteString("```\n\n")
	b.WriteString("Use these for deep-dive queries when the markdown docs don't have enough detail.\n\n")

	// Component index table
	b.WriteString("## Component Index\n\n")
	b.WriteString("| Component | Aliases | CRDs | Services | Deps | File |\n")
	b.WriteString("|-----------|---------|------|----------|------|------|\n")

	type compEntry struct {
		name     string
		aliases  string
		crds     int
		services int
		deps     int
		file     string
	}
	entries := make([]compEntry, 0, len(componentData))
	for _, cd := range componentData {
		name := getStr(cd, "component", "unknown")
		aliases := strings.Join(getStringSlice(cd, "aliases"), ", ")
		crdCount := len(getSlice(cd, "crds"))
		svcCount := len(getSlice(cd, "services"))
		depCount := 0
		if d := getMap(cd, "dependencies"); d != nil {
			depCount += len(getSlice(d, "internal_odh"))
			depCount += len(getSlice(d, "go_modules"))
			depCount += len(getSlice(d, "python_packages"))
		}
		depCount += len(getSlice(cd, "external_connections"))
		depCount += len(getSlice(cd, "runtime_dependencies"))
		filename := SanitizeFilename(name) + ".md"
		entries = append(entries, compEntry{name, aliases, crdCount, svcCount, depCount, filename})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].name < entries[j].name })
	for _, e := range entries {
		b.WriteString(fmt.Sprintf("| %s | %s | %d | %d | %d | [%s](components/%s) |\n",
			escapeMdCell(e.name), escapeMdCell(e.aliases), e.crds, e.services, e.deps, e.file, e.file))
	}
	b.WriteString("\n")

	return b.String()
}
