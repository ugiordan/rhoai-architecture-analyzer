package renderer

import (
	"strings"
	"testing"
)

// forbidden contains definitive negative claims that should never appear in output.
// "found none" != "none exist": the analyzer only sees static sources.
var forbidden = []string{
	"all traffic allowed by default",
	"All pod-to-pod traffic is allowed by default",
	"does not define any",
	"No CRDs defined",
	"No services defined",
	"No dependencies found.\n",
	"No Dockerfiles found.\n",
	"No controller watches found.\n",
	"(none found)\n",
	"none found - all traffic allowed",
}

// emptyComponentData returns component data with no extracted resources.
func emptyComponentData() map[string]interface{} {
	return map[string]interface{}{
		"component":        "empty-component",
		"repo":             "github.com/org/empty",
		"extracted_at":     "2026-05-18T00:00:00Z",
		"analyzer_version": "1.0.0",
	}
}

// ---- SecurityNetworkRenderer: empty data edge cases ----

func TestSecurityNetwork_EmptyServices(t *testing.T) {
	out := (&SecurityNetworkRenderer{}).Render(emptyComponentData())
	if !strings.Contains(out, "none found in analyzed sources") {
		t.Error("empty services should say 'none found in analyzed sources'")
	}
	assertNoForbidden(t, out, "SecurityNetwork empty services")
}

func TestSecurityNetwork_EmptyNetworkPolicies(t *testing.T) {
	out := (&SecurityNetworkRenderer{}).Render(emptyComponentData())
	assertNoForbidden(t, out, "SecurityNetwork empty network policies")
	if strings.Contains(out, "all traffic allowed") {
		t.Error("must not claim all traffic is allowed when policies are not found")
	}
}

func TestSecurityNetwork_EmptySecrets(t *testing.T) {
	out := (&SecurityNetworkRenderer{}).Render(emptyComponentData())
	if !strings.Contains(out, "no secret references found in analyzed sources") {
		t.Error("empty secrets should say 'no secret references found in analyzed sources'")
	}
	assertNoForbidden(t, out, "SecurityNetwork empty secrets")
}

func TestSecurityNetwork_EmptyDeployments(t *testing.T) {
	out := (&SecurityNetworkRenderer{}).Render(emptyComponentData())
	if !strings.Contains(out, "no deployments found in analyzed sources") {
		t.Error("empty deployments should say 'no deployments found in analyzed sources'")
	}
	assertNoForbidden(t, out, "SecurityNetwork empty deployments")
}

func TestSecurityNetwork_EmptyDockerfiles(t *testing.T) {
	out := (&SecurityNetworkRenderer{}).Render(emptyComponentData())
	if !strings.Contains(out, "no Dockerfiles found in analyzed sources") {
		t.Error("empty dockerfiles should say 'no Dockerfiles found in analyzed sources'")
	}
	assertNoForbidden(t, out, "SecurityNetwork empty dockerfiles")
}

func TestSecurityNetwork_WithData_SectionsPresent(t *testing.T) {
	out := (&SecurityNetworkRenderer{}).Render(sampleData())
	if strings.Contains(out, "no secret references found in analyzed sources") {
		t.Error("secrets exist in sample data, should not hedge secrets section")
	}
	if strings.Contains(out, "no deployments found in analyzed sources") {
		t.Error("deployments exist in sample data, should not hedge deployments section")
	}
	if strings.Contains(out, "no Dockerfiles found in analyzed sources") {
		t.Error("dockerfiles exist in sample data, should not hedge dockerfiles section")
	}
}

// ---- docs.go: network and RBAC page edge cases ----

func TestDocsNetworkPage_EmptyNetworkPolicies(t *testing.T) {
	data := emptyComponentData()
	out := renderComponentNetworkPage(data, "empty-component")
	if !strings.Contains(out, "No NetworkPolicy resources were found in the analyzed sources") {
		t.Error("empty network page should use hedged language for missing policies")
	}
	if strings.Contains(out, "All pod-to-pod traffic is allowed by default") {
		t.Error("must never claim all traffic is allowed by default")
	}
	if !strings.Contains(out, "overlays, Helm values, or cluster-level configurations") {
		t.Error("should explain why policies might be missing")
	}
	assertNoForbidden(t, out, "docs network page empty")
}

func TestDocsNetworkPage_WithNetworkPolicies(t *testing.T) {
	data := map[string]interface{}{
		"component": "secured",
		"network_policies": []interface{}{
			map[string]interface{}{
				"name": "allow-ingress",
				"pod_selector": map[string]interface{}{
					"app": "secured",
				},
				"policy_types": []interface{}{"Ingress"},
				"ingress_rules": []interface{}{
					map[string]interface{}{
						"from": []interface{}{
							map[string]interface{}{
								"podSelector": map[string]interface{}{
									"app": "frontend",
								},
							},
						},
					},
				},
			},
		},
	}
	out := renderComponentNetworkPage(data, "secured")
	if strings.Contains(out, "No NetworkPolicy") {
		t.Error("should not show 'No NetworkPolicy' warning when policies exist")
	}
	if strings.Contains(out, "No NetworkPolicy resources were found in the analyzed sources") {
		t.Error("network policy hedging should not appear when policies exist")
	}
}

func TestDocsRBACPage_EmptyRBAC(t *testing.T) {
	data := emptyComponentData()
	out := renderComponentRBACPage(data, "empty-component")
	if !strings.Contains(out, "No ClusterRoles, Roles, or RoleBindings were found in the analyzed sources") {
		t.Error("empty RBAC page should use hedged language")
	}
	if strings.Contains(out, "does not define any") {
		t.Error("must never say 'does not define any'")
	}
	if !strings.Contains(out, "Kustomize overlays, Helm templates") {
		t.Error("should explain alternative sources of RBAC")
	}
	assertNoForbidden(t, out, "docs RBAC page empty")
}

func TestDocsRBACPage_NilRBAC(t *testing.T) {
	data := map[string]interface{}{
		"component": "nil-rbac",
		"rbac":      nil,
	}
	out := renderComponentRBACPage(data, "nil-rbac")
	if !strings.Contains(out, "No ClusterRoles, Roles, or RoleBindings were found in the analyzed sources") {
		t.Error("nil RBAC should use hedged language")
	}
	assertNoForbidden(t, out, "docs RBAC page nil")
}

func TestDocsRBACPage_EmptySlicesRBAC(t *testing.T) {
	data := map[string]interface{}{
		"component": "empty-slices",
		"rbac": map[string]interface{}{
			"cluster_roles":         []interface{}{},
			"cluster_role_bindings": []interface{}{},
			"roles":                 []interface{}{},
			"role_bindings":         []interface{}{},
		},
	}
	out := renderComponentRBACPage(data, "empty-slices")
	if !strings.Contains(out, "No ClusterRoles, Roles, or RoleBindings were found in the analyzed sources") {
		t.Error("RBAC with all empty slices should use hedged language")
	}
	assertNoForbidden(t, out, "docs RBAC page empty slices")
}

func TestDocsRBACPage_WithData(t *testing.T) {
	out := renderComponentRBACPage(sampleData(), "test-controller")
	if strings.Contains(out, "No ClusterRoles") {
		t.Error("should not show empty message when RBAC data exists")
	}
	if strings.Contains(out, "in analyzed sources") {
		t.Error("hedging language should not appear when data exists")
	}
}

// ---- report.go: empty section edge cases ----

func TestReport_EmptyCRDs(t *testing.T) {
	data := emptyComponentData()
	out := (&ReportRenderer{}).Render(data)
	if !strings.Contains(out, "No CRDs found in analyzed sources") {
		t.Error("empty CRDs should say 'No CRDs found in analyzed sources'")
	}
	if strings.Contains(out, "No CRDs defined") {
		t.Error("must not say 'No CRDs defined'")
	}
	assertNoForbidden(t, out, "report empty CRDs")
}

func TestReport_EmptyServices(t *testing.T) {
	data := emptyComponentData()
	out := (&ReportRenderer{}).Render(data)
	if !strings.Contains(out, "No services found in analyzed sources") {
		t.Error("empty services should use hedged language")
	}
	assertNoForbidden(t, out, "report empty services")
}

func TestReport_EmptyDependencies(t *testing.T) {
	data := emptyComponentData()
	out := (&ReportRenderer{}).Render(data)
	if !strings.Contains(out, "No dependencies found in analyzed sources") {
		t.Error("empty dependencies should use hedged language")
	}
	assertNoForbidden(t, out, "report empty dependencies")
}

func TestReport_EmptyDockerfiles(t *testing.T) {
	data := emptyComponentData()
	out := (&ReportRenderer{}).Render(data)
	if !strings.Contains(out, "No Dockerfiles found in analyzed sources") {
		t.Error("empty dockerfiles should use hedged language")
	}
	assertNoForbidden(t, out, "report empty dockerfiles")
}

func TestReport_EmptyControllerWatches(t *testing.T) {
	data := emptyComponentData()
	out := (&ReportRenderer{}).Render(data)
	if !strings.Contains(out, "No controller watches found in analyzed sources") {
		t.Error("empty controller watches should use hedged language")
	}
	assertNoForbidden(t, out, "report empty controller watches")
}

func TestReport_WithData_NoHedging(t *testing.T) {
	out := (&ReportRenderer{}).Render(sampleData())
	if strings.Contains(out, "No CRDs found in analyzed sources") {
		t.Error("should not hedge when CRDs exist")
	}
	if strings.Contains(out, "No services found in analyzed sources") {
		t.Error("should not hedge when services exist")
	}
	if strings.Contains(out, "No controller watches found in analyzed sources") {
		t.Error("should not hedge when controller watches exist")
	}
}

// ---- platform renderers: empty data edge cases ----

func TestPlatformReport_EmptyCRDs(t *testing.T) {
	data := map[string]interface{}{
		"components":       []interface{}{},
		"dependency_graph": []interface{}{},
		"crd_ownership":    map[string]interface{}{},
	}
	out := renderPlatformReport(data)
	if !strings.Contains(out, "No CRDs found in analyzed sources") {
		t.Error("platform report empty CRDs should use hedged language")
	}
	assertNoForbidden(t, out, "platform report empty CRDs")
}

func TestPlatformReport_EmptyDependencies(t *testing.T) {
	data := map[string]interface{}{
		"components":       []interface{}{},
		"dependency_graph": []interface{}{},
		"crd_ownership":    map[string]interface{}{},
	}
	out := renderPlatformReport(data)
	if !strings.Contains(out, "No cross-component dependencies found in analyzed sources") {
		t.Error("platform report empty deps should use hedged language")
	}
	assertNoForbidden(t, out, "platform report empty deps")
}

func TestPlatformDocs_EmptyDependencies(t *testing.T) {
	data := map[string]interface{}{
		"components":       []interface{}{},
		"dependency_graph": []interface{}{},
		"crd_ownership":    map[string]interface{}{},
	}
	pages := renderPlatformDocs(data)
	var found bool
	for _, page := range pages {
		if strings.Contains(page.Content, "No cross-component dependencies detected in analyzed sources") {
			found = true
		}
		assertNoForbidden(t, page.Content, "platform docs page: "+page.Path)
	}
	if !found {
		t.Error("platform docs should use hedged language for empty dependencies")
	}
}

func TestPlatformDocs_EmptyRBACBindings(t *testing.T) {
	data := map[string]interface{}{
		"components":       []interface{}{"comp-a"},
		"dependency_graph": []interface{}{},
		"crd_ownership":    map[string]interface{}{},
		"rbac_cluster_roles": []interface{}{
			map[string]interface{}{
				"owner": "comp-a",
				"name":  "comp-a-role",
				"rules": []interface{}{
					map[string]interface{}{
						"resources": []interface{}{"pods"},
					},
				},
			},
		},
		"component_data": []interface{}{
			map[string]interface{}{
				"component": "comp-a",
				"rbac": map[string]interface{}{
					"cluster_role_bindings": []interface{}{},
				},
			},
		},
	}
	pages := renderPlatformDocs(data)
	for _, page := range pages {
		if strings.Contains(page.Path, "rbac") {
			if !strings.Contains(page.Content, "No RBAC bindings found across analyzed components") {
				t.Error("platform RBAC page should use hedged language for empty bindings")
			}
			assertNoForbidden(t, page.Content, "platform RBAC page")
		}
	}
}

// ---- comprehensive sweep: no forbidden claims in any renderer output ----

func TestAllRenderers_EmptyData_NoForbiddenClaims(t *testing.T) {
	data := emptyComponentData()

	outputs := map[string]string{
		"SecurityNetwork": (&SecurityNetworkRenderer{}).Render(data),
		"Report":          (&ReportRenderer{}).Render(data),
		"RBAC":            (&RBACRenderer{}).Render(data),
		"Component":       (&ComponentRenderer{}).Render(data),
		"C4":              (&C4Renderer{}).Render(data),
		"Dependency":      (&DependencyRenderer{}).Render(data),
		"Dataflow":        (&DataflowRenderer{}).Render(data),
		"Markdown":        RenderComponentMarkdown(data),
	}

	for name, out := range outputs {
		assertNoForbidden(t, out, name+" renderer with empty data")
	}
}

func TestAllRenderers_FullData_NoForbiddenClaims(t *testing.T) {
	data := sampleData()

	outputs := map[string]string{
		"SecurityNetwork": (&SecurityNetworkRenderer{}).Render(data),
		"Report":          (&ReportRenderer{}).Render(data),
		"RBAC":            (&RBACRenderer{}).Render(data),
		"Component":       (&ComponentRenderer{}).Render(data),
		"C4":              (&C4Renderer{}).Render(data),
		"Dependency":      (&DependencyRenderer{}).Render(data),
		"Dataflow":        (&DataflowRenderer{}).Render(data),
		"Markdown":        RenderComponentMarkdown(data),
	}

	for name, out := range outputs {
		assertNoForbidden(t, out, name+" renderer with full data")
	}
}

func TestDocsPages_EmptyData_NoForbiddenClaims(t *testing.T) {
	data := emptyComponentData()
	pages := renderComponentDocs(data, "test/")
	for _, page := range pages {
		assertNoForbidden(t, page.Content, "docs page: "+page.Path)
	}
}

func TestDocsPages_FullData_NoForbiddenClaims(t *testing.T) {
	data := sampleData()
	pages := renderComponentDocs(data, "test/")
	for _, page := range pages {
		assertNoForbidden(t, page.Content, "docs page: "+page.Path)
	}
}

func TestPlatformRenderers_EmptyData_NoForbiddenClaims(t *testing.T) {
	data := map[string]interface{}{
		"components":       []interface{}{},
		"dependency_graph": []interface{}{},
		"crd_ownership":    map[string]interface{}{},
	}

	report := renderPlatformReport(data)
	assertNoForbidden(t, report, "PlatformReport empty")

	pages := renderPlatformDocs(data)
	for _, page := range pages {
		assertNoForbidden(t, page.Content, "PlatformDocs page: "+page.Path)
	}

	diagrams := RenderPlatformAll(data)
	for name, content := range diagrams {
		assertNoForbidden(t, content, "PlatformAll diagram: "+name)
	}
}

// ---- edge cases: partial data ----

func TestDocsNetworkPage_ServicesButNoNetpols(t *testing.T) {
	data := map[string]interface{}{
		"component": "partial",
		"services": []interface{}{
			map[string]interface{}{
				"name": "my-svc",
				"type": "ClusterIP",
				"ports": []interface{}{
					map[string]interface{}{"port": 8080, "protocol": "TCP"},
				},
			},
		},
	}
	out := renderComponentNetworkPage(data, "partial")
	if !strings.Contains(out, "my-svc") {
		t.Error("should render existing services")
	}
	if !strings.Contains(out, "No NetworkPolicy resources were found in the analyzed sources") {
		t.Error("should hedge about missing network policies")
	}
	assertNoForbidden(t, out, "network page with services but no netpols")
}

func TestDocsRBACPage_BindingsButNoRoles(t *testing.T) {
	data := map[string]interface{}{
		"component": "bindings-only",
		"rbac": map[string]interface{}{
			"cluster_roles": []interface{}{},
			"cluster_role_bindings": []interface{}{
				map[string]interface{}{
					"name":     "external-binding",
					"role_ref": "external-role",
					"subjects": []interface{}{
						map[string]interface{}{"kind": "ServiceAccount", "name": "sa"},
					},
				},
			},
			"roles":         []interface{}{},
			"role_bindings": []interface{}{},
		},
	}
	out := renderComponentRBACPage(data, "bindings-only")
	if strings.Contains(out, "No ClusterRoles, Roles, or RoleBindings were found") {
		t.Error("should not show empty message when bindings exist")
	}
	assertNoForbidden(t, out, "RBAC page with bindings only")
}

func TestReport_PartialData_SomeSectionsEmpty(t *testing.T) {
	data := map[string]interface{}{
		"component": "partial",
		"crds": []interface{}{
			map[string]interface{}{"kind": "MyKind", "group": "g", "version": "v1"},
		},
	}
	out := (&ReportRenderer{}).Render(data)
	if strings.Contains(out, "No CRDs found in analyzed sources") {
		t.Error("CRDs exist, should not hedge")
	}
	if !strings.Contains(out, "No services found in analyzed sources") {
		t.Error("services are empty, should hedge")
	}
	if !strings.Contains(out, "No controller watches found in analyzed sources") {
		t.Error("watches are empty, should hedge")
	}
	assertNoForbidden(t, out, "report partial data")
}

// assertNoForbidden checks that output contains none of the forbidden definitive claims.
func assertNoForbidden(t *testing.T, output, context string) {
	t.Helper()
	lower := strings.ToLower(output)
	for _, phrase := range forbidden {
		if strings.Contains(lower, strings.ToLower(phrase)) {
			t.Errorf("[%s] contains forbidden definitive claim: %q", context, phrase)
		}
	}
}
