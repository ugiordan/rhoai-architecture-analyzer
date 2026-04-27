package security

import (
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
		StructType:  "x509.Certificate",
		FieldNames:  []string{"SerialNumber", "Subject", "IsCA", "KeyUsage", "DNSNames"},
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
		StructType:  "cache.ByObject",
		FieldNames:  []string{},
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
		StructType: "cache.ByObject",
		FieldNames: []string{"Field", "Label"},
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
	if len(rules) != 9 {
		t.Fatalf("expected 9 rules, got %d", len(rules))
	}

	expectedRuleIDs := []string{"CGA-003", "CGA-004", "CGA-005", "CGA-006", "CGA-007", "CGA-008", "CGA-009", "CGA-010", "CGA-011"}
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
	g.AddNode(fn1)

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
	g.AddNode(fn2)

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
	g.AddNode(fn3)

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
	g.AddNode(ep1)

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
	g.AddNode(ep2)

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

func TestAllAnnotationsHaveSecurityPrefix(t *testing.T) {
	annotations := []string{
		AnnotCreatesRBAC, AnnotHandlesAdmission, AnnotGeneratesCert,
		AnnotAccessesSecret, AnnotCrossesNamespace, AnnotConfiguresCache,
		AnnotBindsSubject, AnnotWritesPlaintextSecret,
		AnnotHandlesRequest, AnnotExecutesSQL, AnnotDeserializesInput,
		AnnotSubprocessCall, AnnotFileAccess, AnnotTemplateRender,
		AnnotRendersHTML, AnnotEvalUsage, AnnotRedirect,
		AnnotUnsafeBlock, AnnotFFICall, AnnotCommandExecution,
	}
	for _, ann := range annotations {
		if !strings.HasPrefix(ann, SecurityAnnotationPrefix) {
			t.Errorf("annotation %q does not start with %q", ann, SecurityAnnotationPrefix)
		}
	}
}
