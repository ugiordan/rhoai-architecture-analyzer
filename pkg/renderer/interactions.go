// pkg/renderer/interactions.go
package renderer

import (
	"fmt"
	"sort"
	"strings"
)

// interaction represents a single directed relationship between components.
type interaction struct {
	From     string
	To       string
	Type     string
	Evidence string
}

// RenderInteractions produces cross-component/interactions.md from the
// aggregated dependency_graph plus component_refs extracted from component_data.
func RenderInteractions(platformData map[string]interface{}) string {
	var b strings.Builder

	b.WriteString(autoGenHeader)
	b.WriteString("# Cross-Component Interactions\n\n")

	var interactions []interaction
	seen := make(map[string]bool)

	for _, dep := range getSlice(platformData, "dependency_graph") {
		from := getStr(dep, "from", "")
		to := getStr(dep, "to", "")
		typ := getStr(dep, "type", "")
		key := from + "|" + to + "|" + typ
		if seen[key] {
			continue
		}
		seen[key] = true
		evidence := typ
		if strings.HasPrefix(typ, "watches-crd:") {
			evidence = "controller watch"
		} else if typ == "go-module" {
			evidence = "import dependency"
		} else if typ == "uses-image" {
			evidence = "sidecar container"
		}
		interactions = append(interactions, interaction{From: from, To: to, Type: typ, Evidence: evidence})
	}

	for _, cd := range getSlice(platformData, "component_data") {
		compName := getStr(cd, "component", "")
		for _, ref := range getSlice(cd, "component_refs") {
			target := getStr(ref, "target", "")
			typ := getStr(ref, "type", "")
			key := compName + "|" + target + "|component-ref:" + typ
			if seen[key] {
				continue
			}
			seen[key] = true
			interactions = append(interactions, interaction{
				From:     compName,
				To:       target,
				Type:     "component-ref:" + typ,
				Evidence: getStr(ref, "evidence", typ),
			})
		}
	}

	sort.Slice(interactions, func(i, j int) bool {
		if interactions[i].From != interactions[j].From {
			return interactions[i].From < interactions[j].From
		}
		if interactions[i].To != interactions[j].To {
			return interactions[i].To < interactions[j].To
		}
		return interactions[i].Type < interactions[j].Type
	})

	compSet := make(map[string]bool)
	for _, ix := range interactions {
		compSet[ix.From] = true
		compSet[ix.To] = true
	}

	b.WriteString(fmt.Sprintf("## Summary\n\n%d interactions across %d components.\n\n", len(interactions), len(compSet)))

	b.WriteString("## All Interactions\n\n")
	b.WriteString("| From | To | Type | Evidence |\n")
	b.WriteString("|------|-----|------|----------|\n")
	for _, ix := range interactions {
		b.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
			escapeMdCell(ix.From), escapeMdCell(ix.To), escapeMdCell(ix.Type), escapeMdCell(ix.Evidence)))
	}
	b.WriteString("\n")

	b.WriteString("## Per-Component View\n\n")
	dependsOn := make(map[string][]interaction)
	usedBy := make(map[string][]interaction)
	for _, ix := range interactions {
		dependsOn[ix.From] = append(dependsOn[ix.From], ix)
		usedBy[ix.To] = append(usedBy[ix.To], ix)
	}
	compNames := make([]string, 0, len(compSet))
	for c := range compSet {
		compNames = append(compNames, c)
	}
	sort.Strings(compNames)
	for _, comp := range compNames {
		deps := dependsOn[comp]
		uses := usedBy[comp]
		if len(deps) == 0 && len(uses) == 0 {
			continue
		}
		b.WriteString(fmt.Sprintf("### %s\n\n", comp))
		if len(deps) > 0 {
			parts := make([]string, 0, len(deps))
			for _, d := range deps {
				parts = append(parts, fmt.Sprintf("%s (%s)", d.To, d.Type))
			}
			b.WriteString(fmt.Sprintf("**Depends on:** %s  \n", strings.Join(parts, ", ")))
		}
		if len(uses) > 0 {
			parts := make([]string, 0, len(uses))
			for _, u := range uses {
				parts = append(parts, fmt.Sprintf("%s (%s)", u.From, u.Type))
			}
			b.WriteString(fmt.Sprintf("**Used by:** %s  \n", strings.Join(parts, ", ")))
		}
		b.WriteString("\n")
	}

	return b.String()
}
