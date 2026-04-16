package renderer

import (
	"fmt"
	"strings"
)

// C4Renderer renders a C4 context diagram in Structurizr DSL.
type C4Renderer struct{}

func (r *C4Renderer) Filename() string { return "c4-context.dsl" }

func (r *C4Renderer) Render(data map[string]interface{}) string {
	component := getStr(data, "component", "unknown")
	sanitizedComponent := sanitizeID(component)

	var b strings.Builder
	b.WriteString("workspace {\n")
	b.WriteString("\n")
	b.WriteString("    model {\n")

	// Persons
	b.WriteString("        admin = person \"Platform Admin\" \"Manages the OpenShift AI platform\"\n")
	b.WriteString("        user = person \"Data Scientist\" \"Uses ML tools and services\"\n")
	b.WriteString("\n")

	// System boundary
	b.WriteString(fmt.Sprintf("        %s_system = softwareSystem \"%s\" {\n",
		sanitizedComponent, dslEscape(component)))

	// Containers from deployments
	deployments := getSlice(data, "deployments")
	var containerIDs []string
	for _, dep := range deployments {
		name := getStr(dep, "name", "")
		depID := sanitizeID("container_" + name)
		containerIDs = append(containerIDs, depID)

		var image string
		var portsStr string
		for _, c := range getSlice(dep, "containers") {
			if img := getStr(c, "image", ""); img != "" {
				image = img
			}
			for _, p := range getSlice(c, "ports") {
				port := getInt(p, "containerPort")
				proto := getStr(p, "protocol", "TCP")
				portsStr += fmt.Sprintf(" port %d/%s", port, proto)
			}
		}

		desc := "Controller container"
		if image != "" {
			desc = "Container image: " + image
		}
		tech := "Kubernetes"
		if strings.TrimSpace(portsStr) != "" {
			tech = strings.TrimSpace(portsStr)
		}

		b.WriteString(fmt.Sprintf("            %s = container \"%s\" \"%s\" \"%s\"\n",
			depID, dslEscape(name), dslEscape(desc), dslEscape(tech)))
	}

	if len(deployments) == 0 {
		ctrlID := sanitizeID("container_" + component)
		containerIDs = append(containerIDs, ctrlID)
		b.WriteString(fmt.Sprintf("            %s = container \"%s\" \"Main controller\" \"Kubernetes\"\n",
			ctrlID, dslEscape(component)))
	}

	b.WriteString("        }\n")
	b.WriteString("\n")

	// External systems from dependencies
	deps := getMap(data, "dependencies")
	var extIDs []string
	if deps != nil {
		for _, odh := range getSlice(deps, "internal_odh") {
			compName := getStr(odh, "component", "")
			extID := sanitizeID("ext_" + compName)
			extIDs = append(extIDs, extID)
			b.WriteString(fmt.Sprintf("        %s = softwareSystem \"%s\" \"ODH component\" \"Existing System\"\n",
				extID, dslEscape(compName)))
		}
	}

	// Kubernetes API
	k8sID := "kubernetes_api"
	b.WriteString(fmt.Sprintf("        %s = softwareSystem \"Kubernetes API\" \"Cluster API server\" \"Existing System\"\n", k8sID))
	b.WriteString("\n")

	// Relationships
	b.WriteString(fmt.Sprintf("        admin -> %s_system \"Manages\"\n", sanitizedComponent))
	b.WriteString(fmt.Sprintf("        user -> %s_system \"Uses\"\n", sanitizedComponent))

	for _, cid := range containerIDs {
		b.WriteString(fmt.Sprintf("        %s -> %s \"API calls\" \"HTTPS\"\n", cid, k8sID))
	}

	for _, extID := range extIDs {
		b.WriteString(fmt.Sprintf("        %s_system -> %s \"Depends on\"\n", sanitizedComponent, extID))
	}

	// CRD relationships (just the first one for context diagram)
	crds := getSlice(data, "crds")
	for _, crd := range crds {
		kind := getStr(crd, "kind", "")
		group := getStr(crd, "group", "")
		b.WriteString(fmt.Sprintf("        admin -> %s_system \"Creates %s (%s)\" \"YAML/kubectl\"\n",
			sanitizedComponent, dslEscape(kind), dslEscape(group)))
		break
	}

	b.WriteString("    }\n")
	b.WriteString("\n")

	// Views
	b.WriteString("    views {\n")
	b.WriteString(fmt.Sprintf("        systemContext %s_system \"SystemContext\" {\n", sanitizedComponent))
	b.WriteString("            include *\n")
	b.WriteString("            autoLayout\n")
	b.WriteString("        }\n")
	b.WriteString(fmt.Sprintf("        container %s_system \"Containers\" {\n", sanitizedComponent))
	b.WriteString("            include *\n")
	b.WriteString("            autoLayout\n")
	b.WriteString("        }\n")
	b.WriteString("    }\n")
	b.WriteString("\n")
	b.WriteString("}")

	return b.String()
}

// dslEscape escapes special characters for Structurizr DSL strings.
func dslEscape(text string) string {
	text = strings.ReplaceAll(text, `"`, `\"`)
	text = strings.ReplaceAll(text, "\n", " ")
	return text
}
