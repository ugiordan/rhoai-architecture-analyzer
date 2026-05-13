package extractor

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestExtractWebhookBehavior_Mutations(t *testing.T) {
	pkgs := loadGoPackages(fixtureDir())
	if pkgs == nil {
		t.Fatal("failed to load fixture packages")
	}
	behaviors := extractWebhookBehavior(pkgs)
	if len(behaviors) == 0 {
		t.Fatal("expected webhook behaviors")
	}
	b, ok := behaviors["/mutate-v1alpha1-widget"]
	if !ok {
		keys := make([]string, 0, len(behaviors))
		for k := range behaviors {
			keys = append(keys, k)
		}
		t.Fatalf("expected /mutate-v1alpha1-widget, got keys: %v", keys)
	}
	if len(b.Mutations) == 0 {
		t.Fatal("expected mutations from Default() method")
	}
	fields := make(map[string]bool)
	for _, m := range b.Mutations {
		fields[m.Field] = true
	}
	if !fields["spec.image"] {
		t.Errorf("expected mutation on spec.image, got fields: %v", fields)
	}
}

func TestExtractWebhookBehavior_Validations(t *testing.T) {
	pkgs := loadGoPackages(fixtureDir())
	if pkgs == nil {
		t.Fatal("failed to load fixture packages")
	}
	behaviors := extractWebhookBehavior(pkgs)
	b, ok := behaviors["/validate-v1alpha1-widget"]
	if !ok {
		keys := make([]string, 0, len(behaviors))
		for k := range behaviors {
			keys = append(keys, k)
		}
		t.Fatalf("expected /validate-v1alpha1-widget, got keys: %v", keys)
	}
	if len(b.Validations) == 0 {
		t.Fatal("expected validations from ValidateCreate()")
	}
	found := false
	for _, v := range b.Validations {
		if v.Field == "spec.replicas" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected validation on spec.replicas, got %v", b.Validations)
	}
}

func TestExtractWebhookBehavior_TargetType(t *testing.T) {
	pkgs := loadGoPackages(fixtureDir())
	if pkgs == nil {
		t.Fatal("failed to load fixture packages")
	}
	behaviors := extractWebhookBehavior(pkgs)
	for _, b := range behaviors {
		if b.TargetType != "Widget" {
			t.Errorf("expected TargetType=Widget, got %s", b.TargetType)
		}
	}
}

func TestCamelToJSON(t *testing.T) {
	tests := []struct{ input, want string }{
		{"Image", "image"},
		{"GPU", "gpu"},
		{"APIVersion", "apiVersion"},
		{"TLSConfig", "tlsConfig"},
		{"IPAddresses", "ipAddresses"},
		{"URL", "url"},
		{"HTTPSPort", "httpsPort"},
		{"", ""},
		{"a", "a"},
		{"A", "a"},
		{"Replicas", "replicas"},
	}
	for _, tt := range tests {
		got := camelToJSON(tt.input)
		if got != tt.want {
			t.Errorf("camelToJSON(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestGoPathToJSON(t *testing.T) {
	tests := []struct{ input, want string }{
		{"w.Spec.Image", "spec.image"},
		{"w.Spec.GPU", "spec.gpu"},
		{"w.Status", "status"},
		{"w.Spec.APIVersion", "spec.apiVersion"},
		{"w", ""},                                       // single segment
		{"", ""},                                        // empty
		{"w.Spec.TLSConfig.CertFile", "spec.tlsConfig.certFile"},
	}
	for _, tt := range tests {
		got := goPathToJSON(tt.input)
		if got != tt.want {
			t.Errorf("goPathToJSON(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatCaseValues(t *testing.T) {
	// Empty list = default
	got := formatCaseValues(nil)
	if got != "default" {
		t.Errorf("formatCaseValues(nil) = %q, want %q", got, "default")
	}

	// Single value: BasicLit.Value includes quotes, formatExpr returns it as-is
	got = formatCaseValues([]ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: `"create"`}})
	if got != `"create"` {
		t.Errorf("single value: got %q, want %q", got, `"create"`)
	}

	// Multiple values
	got = formatCaseValues([]ast.Expr{
		&ast.BasicLit{Kind: token.STRING, Value: `"create"`},
		&ast.BasicLit{Kind: token.STRING, Value: `"update"`},
	})
	if got != `"create", "update"` {
		t.Errorf("multi value: got %q, want %q", got, `"create", "update"`)
	}
}

func TestExtractWebhookBehavior_NilPackages(t *testing.T) {
	behaviors := extractWebhookBehavior(nil)
	if behaviors != nil {
		t.Error("expected nil for nil packages")
	}
}

func TestExtractWebhookBehavior_FallbackMode(t *testing.T) {
	// An empty GoPackageSet (no packages) should return empty behaviors
	pkgs := &GoPackageSet{Mode: "fallback"}
	behaviors := extractWebhookBehavior(pkgs)
	if len(behaviors) != 0 {
		t.Errorf("expected empty behaviors for fallback mode, got %d", len(behaviors))
	}
}

func TestExtractWebhookBehavior_MethodFollowing(t *testing.T) {
	// The fixture's Default() calls r.setGPUDefaults(w) which sets w.Spec.GPU
	pkgs := loadGoPackages(fixtureDir())
	if pkgs == nil {
		t.Fatal("failed to load fixture packages")
	}
	behaviors := extractWebhookBehavior(pkgs)
	b, ok := behaviors["/mutate-v1alpha1-widget"]
	if !ok {
		t.Fatal("expected mutating webhook behavior")
	}
	// Should have mutation from the followed method setGPUDefaults
	fields := make(map[string]bool)
	for _, m := range b.Mutations {
		fields[m.Field] = true
	}
	if !fields["spec.gpu"] {
		t.Errorf("expected spec.gpu mutation from followed setGPUDefaults method, got fields: %v", fields)
	}
}

func TestExtractWebhookBehavior_MutationCondition(t *testing.T) {
	pkgs := loadGoPackages(fixtureDir())
	if pkgs == nil {
		t.Fatal("failed to load fixture packages")
	}
	behaviors := extractWebhookBehavior(pkgs)
	b := behaviors["/mutate-v1alpha1-widget"]
	for _, m := range b.Mutations {
		if m.Field == "spec.image" {
			if m.Condition == "" {
				t.Error("expected non-empty condition for spec.image (guarded by if w.Spec.Image == \"\")")
			}
			return
		}
	}
	t.Error("spec.image mutation not found")
}

func TestCamelToJSON_EdgeCases(t *testing.T) {
	tests := []struct{ input, want string }{
		{"ID", "id"},
		{"XMLParser", "xmlParser"},
		{"HTMLElement", "htmlElement"},
		{"AB", "ab"},
		{"ABc", "aBc"},
	}
	for _, tt := range tests {
		got := camelToJSON(tt.input)
		if got != tt.want {
			t.Errorf("camelToJSON(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestGoPathToJSON_SingleSegment(t *testing.T) {
	got := goPathToJSON("w")
	if got != "" {
		t.Errorf("single segment should return empty, got %q", got)
	}
}

func TestDedupeFieldOps(t *testing.T) {
	ops := []FieldOp{
		{Field: "spec.image", Operation: "set", Condition: ""},
		{Field: "spec.image", Operation: "set", Condition: ""},
		{Field: "spec.gpu", Operation: "set", Condition: ""},
	}
	deduped := dedupeFieldOps(ops)
	if len(deduped) != 2 {
		t.Errorf("expected 2 after dedup, got %d", len(deduped))
	}
}
