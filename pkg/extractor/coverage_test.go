package extractor

import "testing"

func TestClassifyCoverage(t *testing.T) {
	tests := []struct {
		count    int
		tplCount int
		want     string
	}{
		{0, 0, "none"},
		{1, 0, "sparse"},
		{2, 0, "sparse"},
		{3, 0, "moderate"},
		{5, 0, "moderate"},
		{6, 0, "rich"},
		{10, 0, "rich"},
		// All items are templates -> sparse regardless of count
		{5, 5, "sparse"},
		{10, 10, "sparse"},
		// Some templates but not all -> normal classification
		{5, 2, "moderate"},
	}
	for _, tt := range tests {
		got := classifyCoverage(tt.count, tt.tplCount)
		if got != tt.want {
			t.Errorf("classifyCoverage(%d, %d) = %q, want %q", tt.count, tt.tplCount, got, tt.want)
		}
	}
}

func TestIsTemplateItem(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"", false},
		{"my-service", false},
		{"template-value", true},
		{"$(SERVICE_PORT)", true},
		{"some-$(VAR)-name", true},
		{"${ENV_VAR}", true},
		{"prefix-${SUFFIX}", true},
		{"normal-value", false},
	}
	for _, tt := range tests {
		got := isTemplateItem(tt.input)
		if got != tt.want {
			t.Errorf("isTemplateItem(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestComputeDataCoverage(t *testing.T) {
	arch := &ComponentArchitecture{
		CRDs:            []CRD{{Kind: "Foo"}, {Kind: "Bar"}, {Kind: "Baz"}},
		Services:        []Service{{Name: "svc-1"}, {Name: "template-value"}},
		Deployments:     []Deployment{},
		NetworkPolicies: []NetworkPolicy{},
		IngressRouting:  []IngressResource{{Kind: "HTTPRoute"}, {Kind: "VirtualService"}, {Kind: "Ingress"}, {Kind: "Route"}},
		Webhooks:        []WebhookConfig{{Name: "wh1"}},
		RBAC:            &RBACData{ClusterRoles: []RBACRole{{Name: "role1"}}},
	}

	cov := computeDataCoverage(arch)

	expectations := map[string]string{
		"crds":                "moderate",
		"services":            "sparse",
		"deployments":         "none",
		"network_policies":    "none",
		"ingress_routing":     "moderate",
		"webhooks":            "sparse",
		"rbac":                "sparse",
		"external_connections": "none",
	}

	for section, want := range expectations {
		got := cov[section]
		if got != want {
			t.Errorf("data_coverage[%q] = %q, want %q", section, got, want)
		}
	}
}
