package security

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/arch"
	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

func TestQueryWebhookMissingUpdate(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "Handle",
		File: "webhook.go", Line: 10,
		Annotations: map[string]bool{AnnotHandlesAdmission: true},
		Properties:  map[string]string{"case_values": "admissionv1.Create,admissionv1.Delete"},
	}
	if err := g.AddNode(fn); err != nil { t.Fatal(err) }

	findings := queryWebhookMissingUpdate(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].RuleID != "CGA-003" {
		t.Errorf("expected CGA-003, got %s", findings[0].RuleID)
	}
}

func TestQueryWebhookWithUpdate(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "Handle",
		File: "webhook.go", Line: 10,
		Annotations: map[string]bool{AnnotHandlesAdmission: true},
		Properties:  map[string]string{"case_values": "admissionv1.Create,admissionv1.Delete,admissionv1.Update"},
	}
	if err := g.AddNode(fn); err != nil { t.Fatal(err) }

	findings := queryWebhookMissingUpdate(g)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for webhook with UPDATE, got %d", len(findings))
	}
}

func TestQueryCertAsCA(t *testing.T) {
	g := graph.NewCPG()
	sl := &graph.Node{
		ID: "struct1", Kind: graph.NodeStructLiteral, Name: "x509.Certificate",
		File: "cert.go", Line: 15,
		Annotations: map[string]bool{AnnotGeneratesCert: true},
		StructType:  "x509.Certificate",
		FieldNames:  []string{"SerialNumber", "Subject", "IsCA", "KeyUsage", "DNSNames"},
	}
	if err := g.AddNode(sl); err != nil { t.Fatal(err) }

	findings := queryCertAsCA(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].RuleID != "CGA-005" {
		t.Errorf("expected CGA-005, got %s", findings[0].RuleID)
	}
}

func TestQueryUnfilteredCache(t *testing.T) {
	g := graph.NewCPG()
	sl := &graph.Node{
		ID: "struct1", Kind: graph.NodeStructLiteral, Name: "cache.ByObject",
		File: "cache.go", Line: 10,
		Annotations: map[string]bool{AnnotConfiguresCache: true},
		StructType:  "cache.ByObject",
		FieldNames:  []string{},
	}
	if err := g.AddNode(sl); err != nil { t.Fatal(err) }

	findings := queryUnfilteredCache(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].RuleID != "CGA-007" {
		t.Errorf("expected CGA-007, got %s", findings[0].RuleID)
	}
}

func TestQueryUnfilteredCacheWithFilter(t *testing.T) {
	g := graph.NewCPG()
	sl := &graph.Node{
		ID: "struct1", Kind: graph.NodeStructLiteral, Name: "cache.ByObject",
		File: "cache.go", Line: 10,
		Annotations: map[string]bool{AnnotConfiguresCache: true},
		StructType: "cache.ByObject",
		FieldNames: []string{"Field", "Label"},
	}
	if err := g.AddNode(sl); err != nil { t.Fatal(err) }

	findings := queryUnfilteredCache(g)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for cache with filters, got %d", len(findings))
	}
}

func TestQueryRBACPrecedenceBug(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "createRoleBinding",
		File: "rbac.go", Line: 20,
		Annotations: map[string]bool{
			AnnotCreatesRBAC:  true,
			AnnotBindsSubject: true,
		},
		Properties: map[string]string{},
	}
	if err := g.AddNode(fn); err != nil { t.Fatal(err) }

	findings := queryRBACPrecedenceBug(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].RuleID != "CGA-004" {
		t.Errorf("expected CGA-004, got %s", findings[0].RuleID)
	}
}

func TestQueryCrossNamespaceSecret(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "copySecret",
		File: "secret.go", Line: 30,
		Annotations: map[string]bool{
			AnnotAccessesSecret:   true,
			AnnotCrossesNamespace: true,
		},
		Properties: map[string]string{},
	}
	if err := g.AddNode(fn); err != nil { t.Fatal(err) }

	findings := queryCrossNamespaceSecret(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].RuleID != "CGA-006" {
		t.Errorf("expected CGA-006, got %s", findings[0].RuleID)
	}
}

func TestQueryPlaintextSecrets(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "writeConfig",
		File: "config.go", Line: 40,
		Annotations: map[string]bool{AnnotWritesPlaintextSecret: true},
		Properties:  map[string]string{},
	}
	if err := g.AddNode(fn); err != nil { t.Fatal(err) }

	findings := queryPlaintextSecrets(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].RuleID != "CGA-008" {
		t.Errorf("expected CGA-008, got %s", findings[0].RuleID)
	}
}

func TestQueryWeakSerialEntropy(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "generateCert",
		File: "cert.go", Line: 50,
		Annotations: map[string]bool{AnnotGeneratesCert: true},
		Properties:  map[string]string{},
	}
	callSite := &graph.Node{
		ID:   "call1",
		Kind: graph.NodeCallSite,
		Name: "rand.Int",
		File: "cert.go",
		Line: 55,
		Properties: map[string]string{
			"arg_types": "time.Now().UnixNano()",
		},
	}
	if err := g.AddNode(fn); err != nil { t.Fatal(err) }
	if err := g.AddNode(callSite); err != nil { t.Fatal(err) }
	g.AddEdge(&graph.Edge{From: fn.ID, To: callSite.ID, Kind: graph.EdgeCalls})

	findings := queryWeakSerialEntropy(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].RuleID != "CGA-009" {
		t.Errorf("expected CGA-009, got %s", findings[0].RuleID)
	}
}

func TestSecurityQueriesReturnsAllRules(t *testing.T) {
	rules := securityQueries()
	if len(rules) != 12 {
		t.Fatalf("expected 12 rules, got %d", len(rules))
	}

	expectedRuleIDs := []string{"CGA-003", "CGA-004", "CGA-005", "CGA-006", "CGA-007", "CGA-008", "CGA-009", "CGA-010", "CGA-011", "CGA-012", "CGA-013", "CGA-014"}
	for i, rule := range rules {
		if rule.ID != expectedRuleIDs[i] {
			t.Errorf("rule %d: expected ID %s, got %s", i, expectedRuleIDs[i], rule.ID)
		}
		if rule.Domain != "security" {
			t.Errorf("rule %s: expected domain 'security', got %s", rule.ID, rule.Domain)
		}
		if rule.Run == nil {
			t.Errorf("rule %s: Run function is nil", rule.ID)
		}
	}
}

func TestQueryWebhookMissingUpdate_WithArchRef(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{
		Webhooks: []arch.Webhook{
			{
				Name:          "mwidget.kb.io",
				Type:          "validating",
				Path:          "/validate-widget",
				Sources:       []arch.SourceRef{{Type: "yaml_manifest", File: "config/webhook/manifests.yaml"}},
				FailurePolicy: "Fail",
			},
		},
	}

	fn := &graph.Node{
		ID:   "fn-handler",
		Kind: graph.NodeFunction,
		Name: "Handle",
		File: "pkg/webhook/handler.go",
		Line: 15,
		Properties: map[string]string{
			"case_values": "Create,Delete",
			"param_types": "admission.Request",
		},
		Annotations: map[string]bool{
			AnnotHandlesAdmission: true,
		},
	}
	if err := g.AddNode(fn); err != nil { t.Fatal(err) }

	findings := queryWebhookMissingUpdate(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if !strings.Contains(findings[0].ArchitectureRef, "mwidget.kb.io") {
		t.Errorf("ArchitectureRef should contain webhook name, got %q", findings[0].ArchitectureRef)
	}
}

func TestQueryWebhookMissingUpdate_NoArchData(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID:   "fn-handler",
		Kind: graph.NodeFunction,
		Name: "Handle",
		File: "pkg/webhook/handler.go",
		Line: 15,
		Properties: map[string]string{
			"case_values": "Create,Delete",
			"param_types": "admission.Request",
		},
		Annotations: map[string]bool{
			AnnotHandlesAdmission: true,
		},
	}
	if err := g.AddNode(fn); err != nil { t.Fatal(err) }

	findings := queryWebhookMissingUpdate(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].ArchitectureRef != "" {
		t.Error("expected empty ArchitectureRef without arch data")
	}
}

func TestQueryUnfilteredCache_WithArchRef(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{
		Cache: arch.CacheConfig{
			Issues: []string{
				"Type Widget is watched but has no cache filter",
			},
		},
	}

	sl := &graph.Node{
		ID:         "sl-byobject",
		Kind:       graph.NodeStructLiteral,
		Name:       "ByObject",
		File:       "pkg/controller/setup.go",
		Line:       42,
		StructType: "cache.ByObject",
		FieldNames: []string{"Object"},
		Annotations: map[string]bool{
			AnnotConfiguresCache: true,
		},
	}
	if err := g.AddNode(sl); err != nil { t.Fatal(err) }

	findings := queryUnfilteredCache(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if !strings.Contains(findings[0].ArchitectureRef, "cache_issues:") {
		t.Errorf("ArchitectureRef should contain cache_issues, got %q", findings[0].ArchitectureRef)
	}
}

func TestQueryCrossNamespaceSecret_WithArchRef(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{
		Secrets: []arch.Secret{
			{Name: "webhook-cert", Type: "kubernetes.io/tls", ReferencedBy: []string{"deployment/controller"}},
		},
	}

	fn := &graph.Node{
		ID:   "fn-secret",
		Kind: graph.NodeFunction,
		Name: "fetchSecret",
		File: "pkg/controller/secret.go",
		Line: 20,
		Annotations: map[string]bool{
			AnnotAccessesSecret:   true,
			AnnotCrossesNamespace: true,
		},
	}
	if err := g.AddNode(fn); err != nil { t.Fatal(err) }

	findings := queryCrossNamespaceSecret(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if !strings.Contains(findings[0].ArchitectureRef, "webhook-cert") {
		t.Errorf("ArchitectureRef should contain secret name, got %q", findings[0].ArchitectureRef)
	}
}

func TestQueryComplexityHotspot(t *testing.T) {
	g := graph.NewCPG()

	// High complexity + security annotation -> finding
	fn1 := &graph.Node{
		ID:          "fn1",
		Kind:        graph.NodeFunction,
		Name:        "dangerousFunc",
		File:        "handler.go",
		Line:        10,
		Complexity:  15,
		Annotations: map[string]bool{AnnotHandlesRequest: true},
	}
	if err := g.AddNode(fn1); err != nil { t.Fatal(err) }

	// High complexity, no security annotation -> no finding
	fn2 := &graph.Node{
		ID:          "fn2",
		Kind:        graph.NodeFunction,
		Name:        "internalFunc",
		File:        "internal.go",
		Line:        20,
		Complexity:  20,
		Annotations: make(map[string]bool),
	}
	if err := g.AddNode(fn2); err != nil { t.Fatal(err) }

	// Low complexity + security annotation -> no finding
	fn3 := &graph.Node{
		ID:          "fn3",
		Kind:        graph.NodeFunction,
		Name:        "simpleHandler",
		File:        "handler.go",
		Line:        30,
		Complexity:  3,
		Annotations: map[string]bool{AnnotHandlesRequest: true},
	}
	if err := g.AddNode(fn3); err != nil { t.Fatal(err) }

	findings := queryComplexityHotspot(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].NodeID != "fn1" {
		t.Errorf("expected finding on fn1, got %s", findings[0].NodeID)
	}
}

func TestQueryUntrustedEndpoint(t *testing.T) {
	g := graph.NewCPG()

	// Untrusted endpoint -> finding
	ep1 := &graph.Node{
		ID:          "ep1",
		Kind:        graph.NodeHTTPEndpoint,
		Name:        "publicHandler",
		File:        "api.go",
		Line:        10,
		TrustLevel:  graph.TrustUntrusted,
		Annotations: make(map[string]bool),
	}
	if err := g.AddNode(ep1); err != nil { t.Fatal(err) }

	// Semi-trusted endpoint -> no finding
	ep2 := &graph.Node{
		ID:          "ep2",
		Kind:        graph.NodeHTTPEndpoint,
		Name:        "authHandler",
		File:        "api.go",
		Line:        20,
		TrustLevel:  graph.TrustSemiTrusted,
		Annotations: make(map[string]bool),
	}
	if err := g.AddNode(ep2); err != nil { t.Fatal(err) }

	findings := queryUntrustedEndpoint(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].NodeID != "ep1" {
		t.Errorf("expected finding on ep1, got %s", findings[0].NodeID)
	}
	if findings[0].Severity != "informational" {
		t.Errorf("expected severity informational, got %s", findings[0].Severity)
	}
}

// Cross-domain query tests (CGA-012, CGA-013, CGA-014)

func TestQueryUnprotectedIngress_NoArchData(t *testing.T) {
	g := graph.NewCPG()
	ep := &graph.Node{
		ID: "ep1", Kind: graph.NodeHTTPEndpoint, Name: "publicAPI",
		File: "api.go", Line: 10, TrustLevel: graph.TrustUntrusted,
		Annotations: make(map[string]bool),
	}
	if err := g.AddNode(ep); err != nil { t.Fatal(err) }

	findings := queryUnprotectedIngress(g)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings without arch data, got %d", len(findings))
	}
}

func TestQueryUnprotectedIngress_NoNetworkPolicy(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{NetworkPolicies: []arch.NetworkPolicy{}}

	ep := &graph.Node{
		ID: "ep1", Kind: graph.NodeHTTPEndpoint, Name: "publicAPI",
		File: "api.go", Line: 10, TrustLevel: graph.TrustUntrusted,
		HTTPMethod: "POST", Route: "/api/v1/predict",
		Annotations: make(map[string]bool),
	}
	if err := g.AddNode(ep); err != nil { t.Fatal(err) }

	findings := queryUnprotectedIngress(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].RuleID != "CGA-012" {
		t.Errorf("expected CGA-012, got %s", findings[0].RuleID)
	}
	if !strings.Contains(findings[0].ArchitectureRef, "none defined") {
		t.Errorf("expected 'none defined' in ArchitectureRef, got %q", findings[0].ArchitectureRef)
	}
}

func TestQueryUnprotectedIngress_WithNetworkPolicy(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{
		NetworkPolicies: []arch.NetworkPolicy{
			{
				Name:         "allow-ingress",
				PolicyTypes:  []string{"Ingress"},
				IngressRules: []json.RawMessage{[]byte(`{"from": []}`)},
			},
		},
	}

	ep := &graph.Node{
		ID: "ep1", Kind: graph.NodeHTTPEndpoint, Name: "publicAPI",
		File: "api.go", Line: 10, TrustLevel: graph.TrustUntrusted,
		Annotations: make(map[string]bool),
	}
	if err := g.AddNode(ep); err != nil { t.Fatal(err) }

	findings := queryUnprotectedIngress(g)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings with ingress policy, got %d", len(findings))
	}
}

func TestQueryUnprotectedIngress_SkipsTrustedEndpoints(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{NetworkPolicies: []arch.NetworkPolicy{}}

	ep := &graph.Node{
		ID: "ep1", Kind: graph.NodeHTTPEndpoint, Name: "healthz",
		File: "api.go", Line: 10, TrustLevel: graph.TrustTrusted,
		Annotations: make(map[string]bool),
	}
	if err := g.AddNode(ep); err != nil { t.Fatal(err) }

	findings := queryUnprotectedIngress(g)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for trusted endpoint, got %d", len(findings))
	}
}

func TestQueryOverprivilegedSecretAccess_NoArchData(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "readSecret",
		File: "secret.go", Line: 10,
		Annotations: map[string]bool{AnnotAccessesSecret: true},
	}
	if err := g.AddNode(fn); err != nil { t.Fatal(err) }

	findings := queryOverprivilegedSecretAccess(g)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings without arch data, got %d", len(findings))
	}
}

func TestQueryOverprivilegedSecretAccess_WildcardVerbs(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{
		RBAC: arch.RBAC{
			ClusterRoles: []arch.ClusterRole{
				{
					Name: "controller-role",
					Rules: []arch.RBACRule{
						{Resources: []string{"secrets"}, Verbs: []string{"*"}},
					},
				},
			},
		},
	}

	fn := &graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "readSecret",
		File: "secret.go", Line: 10,
		Annotations: map[string]bool{AnnotAccessesSecret: true},
	}
	if err := g.AddNode(fn); err != nil { t.Fatal(err) }

	findings := queryOverprivilegedSecretAccess(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].RuleID != "CGA-013" {
		t.Errorf("expected CGA-013, got %s", findings[0].RuleID)
	}
}

func TestQueryOverprivilegedSecretAccess_ScopedVerbs(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{
		RBAC: arch.RBAC{
			ClusterRoles: []arch.ClusterRole{
				{
					Name: "controller-role",
					Rules: []arch.RBACRule{
						{Resources: []string{"secrets"}, Verbs: []string{"get", "list", "watch"}},
					},
				},
			},
		},
	}

	fn := &graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "readSecret",
		File: "secret.go", Line: 10,
		Annotations: map[string]bool{AnnotAccessesSecret: true},
	}
	if err := g.AddNode(fn); err != nil { t.Fatal(err) }

	findings := queryOverprivilegedSecretAccess(g)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings with scoped verbs, got %d", len(findings))
	}
}

func TestQueryUncontrolledEgress_NoArchData(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "callAPI",
		File: "client.go", Line: 10,
		Annotations: map[string]bool{AnnotCallsExternal: true},
	}
	if err := g.AddNode(fn); err != nil { t.Fatal(err) }

	findings := queryUncontrolledEgress(g)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings without arch data, got %d", len(findings))
	}
}

func TestQueryUncontrolledEgress_NoEgressPolicy(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{NetworkPolicies: []arch.NetworkPolicy{}}

	fn := &graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "callAPI",
		File: "client.go", Line: 10,
		Annotations: map[string]bool{AnnotCallsExternal: true},
	}
	if err := g.AddNode(fn); err != nil { t.Fatal(err) }

	findings := queryUncontrolledEgress(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].RuleID != "CGA-014" {
		t.Errorf("expected CGA-014, got %s", findings[0].RuleID)
	}
}

func TestQueryUncontrolledEgress_WithEgressPolicy(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{
		NetworkPolicies: []arch.NetworkPolicy{
			{
				Name:        "restrict-egress",
				PolicyTypes: []string{"Egress"},
				EgressRules: []json.RawMessage{[]byte(`{"to": []}`)},
			},
		},
	}

	fn := &graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "callAPI",
		File: "client.go", Line: 10,
		Annotations: map[string]bool{AnnotCallsExternal: true},
	}
	if err := g.AddNode(fn); err != nil { t.Fatal(err) }

	findings := queryUncontrolledEgress(g)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings with egress policy, got %d", len(findings))
	}
}

func TestQueryUncontrolledEgress_ExternalCallNodes(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{NetworkPolicies: []arch.NetworkPolicy{}}

	ec := &graph.Node{
		ID: "ec1", Kind: graph.NodeExternalCall, Name: "sql.Open",
		File: "db.go", Line: 15,
		Annotations: make(map[string]bool),
	}
	if err := g.AddNode(ec); err != nil { t.Fatal(err) }

	findings := queryUncontrolledEgress(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for ExternalCall node, got %d", len(findings))
	}
}

func TestQueryOverprivilegedSecretAccess_WildcardResource(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{
		RBAC: arch.RBAC{
			ClusterRoles: []arch.ClusterRole{
				{
					Name: "admin-role",
					Rules: []arch.RBACRule{
						{APIGroups: []string{""}, Resources: []string{"*"}, Verbs: []string{"*"}},
					},
				},
			},
		},
	}

	fn := &graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "readSecret",
		File: "secret.go", Line: 10,
		Annotations: map[string]bool{AnnotAccessesSecret: true},
	}
	if err := g.AddNode(fn); err != nil { t.Fatal(err) }

	findings := queryOverprivilegedSecretAccess(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for wildcard resource, got %d", len(findings))
	}
}

func TestQueryOverprivilegedSecretAccess_NonCoreAPIGroup(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{
		RBAC: arch.RBAC{
			ClusterRoles: []arch.ClusterRole{
				{
					Name: "custom-role",
					Rules: []arch.RBACRule{
						{APIGroups: []string{"apps"}, Resources: []string{"secrets"}, Verbs: []string{"*"}},
					},
				},
			},
		},
	}

	fn := &graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "readSecret",
		File: "secret.go", Line: 10,
		Annotations: map[string]bool{AnnotAccessesSecret: true},
	}
	if err := g.AddNode(fn); err != nil { t.Fatal(err) }

	findings := queryOverprivilegedSecretAccess(g)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for non-core API group, got %d", len(findings))
	}
}

func TestQueryUnprotectedIngress_EmptyPolicyTypes(t *testing.T) {
	g := graph.NewCPG()
	// Per K8s spec, empty PolicyTypes with ingress rules defaults to Ingress.
	g.ArchData = &arch.Data{
		NetworkPolicies: []arch.NetworkPolicy{
			{
				Name:         "legacy-policy",
				PolicyTypes:  []string{},
				IngressRules: []json.RawMessage{[]byte(`{"from": []}`)},
			},
		},
	}

	ep := &graph.Node{
		ID: "ep1", Kind: graph.NodeHTTPEndpoint, Name: "publicAPI",
		File: "api.go", Line: 10, TrustLevel: graph.TrustUntrusted,
		Annotations: make(map[string]bool),
	}
	if err := g.AddNode(ep); err != nil { t.Fatal(err) }

	findings := queryUnprotectedIngress(g)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings when empty PolicyTypes has ingress rules (K8s default), got %d", len(findings))
	}
}

func TestQueryUncontrolledEgress_DedupAnnotationAndExternalCall(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{NetworkPolicies: []arch.NetworkPolicy{}}

	// Function with calls_external annotation
	fn := &graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "callAPI",
		File: "client.go", Line: 10,
		Annotations: map[string]bool{AnnotCallsExternal: true},
	}
	// ExternalCall node with same file+name (should be deduped)
	ec := &graph.Node{
		ID: "ec1", Kind: graph.NodeExternalCall, Name: "callAPI",
		File: "client.go", Line: 12,
		Annotations: make(map[string]bool),
	}
	if err := g.AddNode(fn); err != nil { t.Fatal(err) }
	if err := g.AddNode(ec); err != nil { t.Fatal(err) }

	findings := queryUncontrolledEgress(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding (deduped), got %d", len(findings))
	}
	if findings[0].NodeID != "fn1" {
		t.Errorf("expected finding from function annotation, got node %s", findings[0].NodeID)
	}
}

func TestQueryUncontrolledEgress_DifferentExternalCalls(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{NetworkPolicies: []arch.NetworkPolicy{}}

	ec1 := &graph.Node{
		ID: "ec1", Kind: graph.NodeExternalCall, Name: "sql.Open",
		File: "db.go", Line: 15, Annotations: make(map[string]bool),
	}
	ec2 := &graph.Node{
		ID: "ec2", Kind: graph.NodeExternalCall, Name: "http.Get",
		File: "api.go", Line: 20, Annotations: make(map[string]bool),
	}
	if err := g.AddNode(ec1); err != nil { t.Fatal(err) }
	if err := g.AddNode(ec2); err != nil { t.Fatal(err) }

	findings := queryUncontrolledEgress(g)
	if len(findings) != 2 {
		t.Fatalf("expected 2 findings for different external calls, got %d", len(findings))
	}
}

func TestQueryUnprotectedIngress_MultipleUntrustedEndpoints(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{NetworkPolicies: []arch.NetworkPolicy{}}

	ep1 := &graph.Node{
		ID: "ep1", Kind: graph.NodeHTTPEndpoint, Name: "predict",
		File: "api.go", Line: 10, TrustLevel: graph.TrustUntrusted,
		HTTPMethod: "POST", Route: "/api/v1/predict",
		Annotations: make(map[string]bool),
	}
	ep2 := &graph.Node{
		ID: "ep2", Kind: graph.NodeHTTPEndpoint, Name: "upload",
		File: "api.go", Line: 20, TrustLevel: graph.TrustUntrusted,
		HTTPMethod: "POST", Route: "/api/v1/upload",
		Annotations: make(map[string]bool),
	}
	if err := g.AddNode(ep1); err != nil { t.Fatal(err) }
	if err := g.AddNode(ep2); err != nil { t.Fatal(err) }

	findings := queryUnprotectedIngress(g)
	if len(findings) != 2 {
		t.Fatalf("expected 2 findings for 2 untrusted endpoints, got %d", len(findings))
	}
}

func TestQueryOverprivilegedSecretAccess_EmptyAPIGroups(t *testing.T) {
	// Empty APIGroups list implies core group per K8s convention
	g := graph.NewCPG()
	g.ArchData = &arch.Data{
		RBAC: arch.RBAC{
			ClusterRoles: []arch.ClusterRole{
				{
					Name: "legacy-role",
					Rules: []arch.RBACRule{
						{Resources: []string{"secrets"}, Verbs: []string{"*"}},
					},
				},
			},
		},
	}

	fn := &graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "readSecret",
		File: "secret.go", Line: 10,
		Annotations: map[string]bool{AnnotAccessesSecret: true},
	}
	if err := g.AddNode(fn); err != nil { t.Fatal(err) }

	findings := queryOverprivilegedSecretAccess(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for empty APIGroups (implies core), got %d", len(findings))
	}
}

func TestQueryOverprivilegedSecretAccess_WildcardAPIGroup(t *testing.T) {
	g := graph.NewCPG()
	g.ArchData = &arch.Data{
		RBAC: arch.RBAC{
			ClusterRoles: []arch.ClusterRole{
				{
					Name: "super-admin",
					Rules: []arch.RBACRule{
						{APIGroups: []string{"*"}, Resources: []string{"secrets"}, Verbs: []string{"*"}},
					},
				},
			},
		},
	}

	fn := &graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "readSecret",
		File: "secret.go", Line: 10,
		Annotations: map[string]bool{AnnotAccessesSecret: true},
	}
	if err := g.AddNode(fn); err != nil { t.Fatal(err) }

	findings := queryOverprivilegedSecretAccess(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for wildcard APIGroup, got %d", len(findings))
	}
}

func TestHasEgressNetworkPolicy_EmptyPolicyTypes(t *testing.T) {
	// Unlike ingress, empty PolicyTypes does NOT default to Egress
	g := graph.NewCPG()
	g.ArchData = &arch.Data{
		NetworkPolicies: []arch.NetworkPolicy{
			{
				Name:        "legacy",
				PolicyTypes: []string{},
				EgressRules: []json.RawMessage{[]byte(`{"to": []}`)},
			},
		},
	}
	if hasEgressNetworkPolicy(g) {
		t.Error("empty PolicyTypes should NOT match egress (K8s defaults to Ingress only)")
	}
}

func TestAllAnnotationsHaveSecurityPrefix(t *testing.T) {
	annotations := []string{
		AnnotCreatesRBAC, AnnotHandlesAdmission, AnnotGeneratesCert,
		AnnotAccessesSecret, AnnotCrossesNamespace, AnnotConfiguresCache,
		AnnotBindsSubject, AnnotWritesPlaintextSecret,
		AnnotHandlesRequest, AnnotExecutesSQL, AnnotDeserializesInput,
		AnnotSubprocessCall, AnnotFileAccess, AnnotTemplateRender,
		AnnotRendersHTML, AnnotEvalUsage, AnnotRedirect,
		AnnotUnsafeBlock, AnnotFFICall, AnnotCommandExecution,
		AnnotCallsExternal,
	}
	for _, ann := range annotations {
		if !strings.HasPrefix(ann, SecurityAnnotationPrefix) {
			t.Errorf("annotation %q does not start with %q", ann, SecurityAnnotationPrefix)
		}
	}
}
