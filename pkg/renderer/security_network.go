package renderer

import (
	"fmt"
	"sort"
	"strings"
)

// SecurityNetworkRenderer renders a layered ASCII security and network diagram.
type SecurityNetworkRenderer struct{}

func (r *SecurityNetworkRenderer) Filename() string { return "security-network.txt" }

func (r *SecurityNetworkRenderer) Render(data map[string]interface{}) string {
	component := getStr(data, "component", "unknown")
	width := 80

	var b strings.Builder

	header := func(title string) {
		b.WriteString("\n")
		b.WriteString(strings.Repeat("=", width) + "\n")
		b.WriteString(fmt.Sprintf("  %s\n", title))
		b.WriteString(strings.Repeat("=", width) + "\n")
	}

	section := func(title string) {
		dashLen := width - len(title) - 5
		if dashLen < 1 {
			dashLen = 1
		}
		b.WriteString(fmt.Sprintf("\n--- %s %s\n", title, strings.Repeat("-", dashLen)))
	}

	// Title box
	b.WriteString("+" + strings.Repeat("=", width-2) + "+\n")
	b.WriteString("|" + centerPad("SECURITY & NETWORK ARCHITECTURE", width-2) + "|\n")
	b.WriteString("|" + centerPad(component, width-2) + "|\n")
	b.WriteString("+" + strings.Repeat("=", width-2) + "+\n")

	// Network Topology
	header("NETWORK TOPOLOGY")
	services := getSlice(data, "services")
	section("Services")
	if len(services) > 0 {
		for _, svc := range services {
			name := getStr(svc, "name", "")
			svcType := getStr(svc, "type", "ClusterIP")
			b.WriteString(fmt.Sprintf("  [%s] %s\n", svcType, name))
			for _, port := range getSlice(svc, "ports") {
				portName := getStr(port, "name", "")
				portNum := getInt(port, "port")
				target := getInt(port, "targetPort")
				protocol := getStr(port, "protocol", "TCP")
				tlsHint := ""
				if strings.Contains(strings.ToLower(portName), "https") || portNum == 443 {
					tlsHint = " (TLS)"
				}
				b.WriteString(fmt.Sprintf("    Port: %d -> %d/%s [%s]%s\n",
					portNum, target, protocol, portName, tlsHint))
			}
		}
	} else {
		b.WriteString("  (none found)\n")
	}

	// Network Policies
	policies := getSlice(data, "network_policies")
	section("Network Policies")
	if len(policies) > 0 {
		for _, pol := range policies {
			name := getStr(pol, "name", "")
			types := getStringSlice(pol, "policy_types")
			selector := getMap(pol, "pod_selector")
			selStr := "(all pods)"
			if selector != nil && len(selector) > 0 {
				parts := make([]string, 0, len(selector))
				keys := make([]string, 0, len(selector))
				for k := range selector {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					parts = append(parts, fmt.Sprintf("%s=%v", k, selector[k]))
				}
				selStr = strings.Join(parts, ", ")
			}
			b.WriteString(fmt.Sprintf("  Policy: %s\n", name))
			b.WriteString(fmt.Sprintf("    Selector: %s\n", selStr))
			b.WriteString(fmt.Sprintf("    Types: %s\n", strings.Join(types, ", ")))

			for _, ing := range getSlice(pol, "ingress_rules") {
				for _, p := range getSlice(ing, "ports") {
					b.WriteString(fmt.Sprintf("    INGRESS: port %v/%s\n",
						p["port"], getStr(p, "protocol", "TCP")))
				}
				fromList := getSlice(ing, "from")
				if len(fromList) == 0 {
					b.WriteString("    INGRESS: from ALL sources\n")
				}
			}
			for _, eg := range getSlice(pol, "egress_rules") {
				for _, p := range getSlice(eg, "ports") {
					b.WriteString(fmt.Sprintf("    EGRESS: port %v/%s\n",
						p["port"], getStr(p, "protocol", "TCP")))
				}
			}
		}
	} else {
		b.WriteString("  (none found - all traffic allowed by default)\n")
	}

	// RBAC Summary
	header("RBAC SUMMARY")
	rbac := getMap(data, "rbac")
	if rbac == nil {
		rbac = map[string]interface{}{}
	}
	crCount := len(getSlice(rbac, "cluster_roles"))
	crbCount := len(getSlice(rbac, "cluster_role_bindings"))
	rCount := len(getSlice(rbac, "roles"))
	rbCount := len(getSlice(rbac, "role_bindings"))
	markerCount := len(getSlice(rbac, "kubebuilder_markers"))

	b.WriteString(fmt.Sprintf("  ClusterRoles:        %d\n", crCount))
	b.WriteString(fmt.Sprintf("  ClusterRoleBindings: %d\n", crbCount))
	b.WriteString(fmt.Sprintf("  Roles:               %d\n", rCount))
	b.WriteString(fmt.Sprintf("  RoleBindings:        %d\n", rbCount))
	b.WriteString(fmt.Sprintf("  Kubebuilder markers: %d\n", markerCount))

	for _, cr := range getSlice(rbac, "cluster_roles") {
		name := getStr(cr, "name", "")
		rules := getSlice(cr, "rules")
		totalResources := 0
		for _, rule := range rules {
			totalResources += len(getStringSlice(rule, "resources"))
		}
		b.WriteString(fmt.Sprintf("    CR: %s (%d resource types)\n", name, totalResources))
	}

	// Secrets Inventory
	header("SECRETS INVENTORY")
	secrets := getSlice(data, "secrets_referenced")
	if len(secrets) > 0 {
		for _, secret := range secrets {
			name := getStr(secret, "name", "")
			stype := getStr(secret, "type", "Opaque")
			refs := getStringSlice(secret, "referenced_by")
			prov := getStr(secret, "provisioned_by", "")
			b.WriteString(fmt.Sprintf("  Secret: %s\n", name))
			b.WriteString(fmt.Sprintf("    Type: %s\n", stype))
			if prov == "" {
				prov = "unknown"
			}
			b.WriteString(fmt.Sprintf("    Provisioned by: %s\n", prov))
			for _, ref := range refs {
				b.WriteString(fmt.Sprintf("    Referenced by: %s\n", ref))
			}
		}
	} else {
		b.WriteString("  (no secret references found)\n")
	}

	// Deployment Security Controls
	header("DEPLOYMENT SECURITY CONTROLS")
	deployments := getSlice(data, "deployments")
	if len(deployments) > 0 {
		for _, dep := range deployments {
			name := getStr(dep, "name", "")
			sa := getStr(dep, "service_account", "")
			automount := getBool(dep, "automount_service_account_token", true)
			b.WriteString(fmt.Sprintf("  Deployment: %s\n", name))
			if sa == "" {
				sa = "(default)"
			}
			b.WriteString(fmt.Sprintf("    Service Account: %s\n", sa))
			b.WriteString(fmt.Sprintf("    Automount SA Token: %v\n", automount))

			for _, container := range getSlice(dep, "containers") {
				cName := getStr(container, "name", "")
				sc := getMap(container, "security_context")
				b.WriteString(fmt.Sprintf("    Container: %s\n", cName))
				if sc != nil && len(sc) > 0 {
					keys := make([]string, 0, len(sc))
					for k := range sc {
						keys = append(keys, k)
					}
					sort.Strings(keys)
					for _, key := range keys {
						val := sc[key]
						var valStr string
						switch v := val.(type) {
						case map[string]interface{}:
							parts := make([]string, 0, len(v))
							vkeys := make([]string, 0, len(v))
							for k := range v {
								vkeys = append(vkeys, k)
							}
							sort.Strings(vkeys)
							for _, k := range vkeys {
								parts = append(parts, fmt.Sprintf("%s=%v", k, v[k]))
							}
							valStr = strings.Join(parts, ", ")
						case []interface{}:
							parts := make([]string, 0, len(v))
							for _, item := range v {
								parts = append(parts, fmt.Sprintf("%v", item))
							}
							valStr = strings.Join(parts, ", ")
						default:
							valStr = fmt.Sprintf("%v", val)
						}
						b.WriteString(fmt.Sprintf("      %s: %s\n", key, valStr))
					}
				} else {
					b.WriteString("      (no security context)\n")
				}
			}
		}
	} else {
		b.WriteString("  (no deployments found)\n")
	}

	// Dockerfile Security
	header("DOCKERFILE SECURITY")
	dockerfiles := getSlice(data, "dockerfiles")
	if len(dockerfiles) > 0 {
		for _, df := range dockerfiles {
			path := getStr(df, "path", "")
			base := getStr(df, "base_image", "")
			user := getStr(df, "user", "")
			stages := getInt(df, "stages")
			issues := getStringSlice(df, "issues")
			b.WriteString(fmt.Sprintf("  %s\n", path))
			b.WriteString(fmt.Sprintf("    Base image: %s\n", base))
			b.WriteString(fmt.Sprintf("    Stages: %d\n", stages))
			if user == "" {
				user = "(not set)"
			}
			b.WriteString(fmt.Sprintf("    User: %s\n", user))
			for _, issue := range issues {
				b.WriteString(fmt.Sprintf("    [!] %s\n", issue))
			}
		}
	} else {
		b.WriteString("  (no Dockerfiles found)\n")
	}

	b.WriteString("\n+" + strings.Repeat("=", width-2) + "+")

	return b.String()
}

// centerPad centers text within the given width, safe for multi-byte UTF-8.
func centerPad(text string, width int) string {
	runes := []rune(text)
	if len(runes) >= width {
		return string(runes[:width])
	}
	totalPad := width - len(runes)
	left := totalPad / 2
	right := totalPad - left
	return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
}
