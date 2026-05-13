package extractor

import "testing"

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
