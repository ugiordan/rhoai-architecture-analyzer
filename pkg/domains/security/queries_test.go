package security

import (
	"strings"
	"testing"

	"github.com/ugiordan/rhoai-architecture-analyzer/pkg/arch"
	"github.com/ugiordan/rhoai-architecture-analyzer/pkg/graph"
)

func TestQueryWebhookMissingUpdate(t *testing.T) {
	g := graph.NewCPG()
	fn := &graph.Node{
		ID: "fn1", Kind: graph.NodeFunction, Name: "Handle",
		File: "webhook.go", Line: 10,
		Annotations: map[string]bool{AnnotHandlesAdmission: true},
		Properties:  map[string]string{"case_values": "admissionv1.Create,admissionv1.Delete"},
	}
	g.AddNode(fn)

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
	g.AddNode(fn)

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
		Properties: map[string]string{
			"type":   "x509.Certificate",
			"fields": "SerialNumber,Subject,IsCA,KeyUsage,DNSNames",
		},
	}
	g.AddNode(sl)

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
		Properties:  map[string]string{"type": "cache.ByObject", "fields": ""},
	}
	g.AddNode(sl)

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
		Properties:  map[string]string{"type": "cache.ByObject", "fields": "Field,Label"},
	}
	g.AddNode(sl)

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
	g.AddNode(fn)

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
	g.AddNode(fn)

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
	g.AddNode(fn)

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
	g.AddNode(fn)
	g.AddNode(callSite)
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
	if len(rules) != 7 {
		t.Fatalf("expected 7 rules, got %d", len(rules))
	}

	expectedRuleIDs := []string{"CGA-003", "CGA-004", "CGA-005", "CGA-006", "CGA-007", "CGA-008", "CGA-009"}
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
				Source:        "config/webhook/manifests.yaml",
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
	g.AddNode(fn)

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
	g.AddNode(fn)

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
		ID:   "sl-byobject",
		Kind: graph.NodeStructLiteral,
		Name: "ByObject",
		File: "pkg/controller/setup.go",
		Line: 42,
		Properties: map[string]string{
			"type":   "cache.ByObject",
			"fields": "Object",
		},
		Annotations: map[string]bool{
			AnnotConfiguresCache: true,
		},
	}
	g.AddNode(sl)

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
	g.AddNode(fn)

	findings := queryCrossNamespaceSecret(g)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if !strings.Contains(findings[0].ArchitectureRef, "webhook-cert") {
		t.Errorf("ArchitectureRef should contain secret name, got %q", findings[0].ArchitectureRef)
	}
}
