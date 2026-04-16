package renderer

import (
	"fmt"
	"strings"
)

// DataflowRenderer renders a Mermaid sequence diagram showing static connections.
type DataflowRenderer struct{}

func (r *DataflowRenderer) Filename() string { return "dataflow.mmd" }

func (r *DataflowRenderer) Render(data map[string]interface{}) string {
	component := getStr(data, "component", "unknown")

	var b strings.Builder
	b.WriteString("sequenceDiagram\n")
	b.WriteString(fmt.Sprintf("    %%%% Static dataflow for %s\n", component))
	b.WriteString("\n")

	deployments := getSlice(data, "deployments")
	services := getSlice(data, "services")
	watches := getSlice(data, "controller_watches")
	crds := getSlice(data, "crds")

	// Collect participants
	participants := []string{"KubernetesAPI"}
	b.WriteString("    participant KubernetesAPI as Kubernetes API\n")

	// Add deployments as participants
	for _, dep := range deployments {
		name := getStr(dep, "name", "")
		pid := sanitizeID(name)
		if !containsStr(participants, pid) {
			participants = append(participants, pid)
			b.WriteString(fmt.Sprintf("    participant %s as %s\n", pid, escapeLabel(name)))
		}
	}

	// Add generic controller if no deployments
	if len(deployments) == 0 {
		ctrlID := sanitizeID(component)
		participants = append(participants, ctrlID)
		b.WriteString(fmt.Sprintf("    participant %s as %s\n", ctrlID, escapeLabel(component)))
	}

	b.WriteString("\n")

	// Determine main controller ID
	controllerID := sanitizeID(component)
	if len(deployments) > 0 {
		controllerID = sanitizeID(getStr(deployments[0], "name", ""))
	}

	// Controller watches
	for _, w := range watches {
		wtype := getStr(w, "type", "")
		gvk := getStr(w, "gvk", "")
		kind := gvk
		if idx := strings.LastIndex(gvk, "/"); idx >= 0 {
			kind = gvk[idx+1:]
		}
		switch wtype {
		case "For":
			b.WriteString(fmt.Sprintf("    KubernetesAPI->>+%s: Watch %s (reconcile)\n",
				controllerID, escapeLabel(kind)))
		case "Owns":
			b.WriteString(fmt.Sprintf("    %s->>KubernetesAPI: Create/Update %s\n",
				controllerID, escapeLabel(kind)))
		case "Watches":
			b.WriteString(fmt.Sprintf("    KubernetesAPI-->>+%s: Watch %s (informer)\n",
				controllerID, escapeLabel(kind)))
		}
	}

	// Service connections
	if len(services) > 0 {
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("    Note over %s: Exposed Services\n", controllerID))
		for _, svc := range services {
			name := getStr(svc, "name", "")
			for _, port := range getSlice(svc, "ports") {
				portNum := getInt(port, "port")
				portName := getStr(port, "name", "")
				protocol := getStr(port, "protocol", "TCP")
				b.WriteString(fmt.Sprintf("    Note right of %s: %s:%d/%s [%s]\n",
					controllerID, escapeLabel(name), portNum, protocol, escapeLabel(portName)))
			}
		}
	}

	// CRD definitions
	if len(crds) > 0 {
		b.WriteString("\n")
		b.WriteString("    Note over KubernetesAPI: Defined CRDs\n")
		for _, crd := range crds {
			kind := getStr(crd, "kind", "")
			group := getStr(crd, "group", "")
			version := getStr(crd, "version", "")
			b.WriteString(fmt.Sprintf("    Note right of KubernetesAPI: %s (%s/%s)\n",
				escapeLabel(kind), escapeLabel(group), escapeLabel(version)))
		}
	}

	return strings.TrimRight(b.String(), "\n")
}

func containsStr(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}
