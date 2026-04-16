package renderer

import (
	"fmt"
	"sort"
	"strings"
)

// ComponentRenderer renders component architecture as a Mermaid graph showing
// CRDs, controllers, owned/watched resources, and internal dependencies.
type ComponentRenderer struct{}

func (r *ComponentRenderer) Filename() string { return "component.mmd" }

func (r *ComponentRenderer) Render(data map[string]interface{}) string {
	component := getStr(data, "component", "unknown")

	var b strings.Builder
	b.WriteString("graph LR\n")
	b.WriteString(fmt.Sprintf("    %%%% Component architecture for %s\n", component))
	b.WriteString("\n")
	b.WriteString("    classDef crd fill:#e74c3c,stroke:#c0392b,color:#fff\n")
	b.WriteString("    classDef controller fill:#3498db,stroke:#2980b9,color:#fff\n")
	b.WriteString("    classDef owned fill:#2ecc71,stroke:#27ae60,color:#fff\n")
	b.WriteString("    classDef external fill:#95a5a6,stroke:#7f8c8d,color:#fff\n")
	b.WriteString("    classDef dep fill:#f39c12,stroke:#e67e22,color:#fff\n")
	b.WriteString("\n")

	counter := 0
	nextID := func(prefix string) string {
		counter++
		return fmt.Sprintf("%s_%d", prefix, counter)
	}

	crds := getSlice(data, "crds")
	watches := getSlice(data, "controller_watches")
	deployments := getSlice(data, "deployments")

	// Build kind sets by watch type
	forKinds := make(map[string]bool)
	ownsKinds := make(map[string]bool)
	watchesKinds := make(map[string]bool)
	for _, w := range watches {
		gvk := getStr(w, "gvk", "")
		kind := gvk
		if idx := strings.LastIndex(gvk, "/"); idx >= 0 {
			kind = gvk[idx+1:]
		}
		wtype := getStr(w, "type", "")
		switch wtype {
		case "For":
			forKinds[kind] = true
		case "Owns":
			ownsKinds[kind] = true
		case "Watches":
			watchesKinds[kind] = true
		}
	}

	// Controller subgraph
	b.WriteString(fmt.Sprintf("    subgraph controller[\"%s Controller\"]\n", escapeLabel(component)))
	if len(deployments) > 0 {
		for _, dep := range deployments {
			depName := getStr(dep, "name", "")
			depID := nextID("dep")
			b.WriteString(fmt.Sprintf("        %s[\"%s\"]\n", depID, escapeLabel(depName)))
			b.WriteString(fmt.Sprintf("        class %s controller\n", depID))
		}
	} else {
		ctrlID := nextID("ctrl")
		b.WriteString(fmt.Sprintf("        %s[\"Controller\"]\n", ctrlID))
		b.WriteString(fmt.Sprintf("        class %s controller\n", ctrlID))
	}
	b.WriteString("    end\n")
	b.WriteString("\n")

	// CRDs watched via For
	for _, crd := range crds {
		kind := getStr(crd, "kind", "")
		group := getStr(crd, "group", "")
		version := getStr(crd, "version", "")
		crdID := sanitizeID("crd_" + kind)
		label := fmt.Sprintf("%s\\n%s/%s", kind, group, version)
		b.WriteString(fmt.Sprintf("    %s{{\"%s\"}}\n", crdID, escapeLabel(label)))
		b.WriteString(fmt.Sprintf("    class %s crd\n", crdID))
		if forKinds[kind] {
			b.WriteString(fmt.Sprintf("    %s -->|\"For (reconciles)\"| controller\n", crdID))
		}
	}

	// Owned resources
	sortedOwns := sortedKeys(ownsKinds)
	for _, kind := range sortedOwns {
		resID := nextID("owned")
		b.WriteString(fmt.Sprintf("    controller -->|\"Owns\"| %s[\"%s\"]\n", resID, escapeLabel(kind)))
		b.WriteString(fmt.Sprintf("    class %s owned\n", resID))
	}

	// Watched resources
	sortedWatches := sortedKeys(watchesKinds)
	for _, kind := range sortedWatches {
		resID := nextID("watch")
		b.WriteString(fmt.Sprintf("    %s[\"%s\"] -->|\"Watches\"| controller\n", resID, escapeLabel(kind)))
		b.WriteString(fmt.Sprintf("    class %s external\n", resID))
	}

	// Internal ODH dependencies
	deps := getMap(data, "dependencies")
	if deps != nil {
		for _, odhDep := range getSlice(deps, "internal_odh") {
			comp := getStr(odhDep, "component", "")
			depID := nextID("odh")
			b.WriteString(fmt.Sprintf("    controller -.->|\"depends on\"| %s[\"%s\"]\n", depID, escapeLabel(comp)))
			b.WriteString(fmt.Sprintf("    class %s dep\n", depID))
		}
	}

	return strings.TrimRight(b.String(), "\n")
}

func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
