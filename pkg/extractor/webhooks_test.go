package extractor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractWebhooks_YAMLSource(t *testing.T) {
	// Create a temporary directory with a ValidatingWebhookConfiguration YAML
	tmpDir := t.TempDir()
	yamlContent := `apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
metadata:
  name: test-webhook-config
webhooks:
  - name: test-webhook.example.com
    failurePolicy: Fail
    sideEffects: None
    clientConfig:
      service:
        name: test-service
        namespace: test-ns
        path: /validate
    rules:
      - apiGroups: ["apps"]
        apiVersions: ["v1"]
        resources: ["deployments"]
        operations: ["CREATE", "UPDATE"]
`
	yamlPath := filepath.Join(tmpDir, "webhook-config.yaml")
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write test YAML: %v", err)
	}

	// Run extractWebhooks
	webhooks := extractWebhooks(tmpDir)

	// Verify results
	if len(webhooks) != 1 {
		t.Fatalf("expected 1 webhook, got %d", len(webhooks))
	}

	wh := webhooks[0]
	if wh.Name != "test-webhook.example.com" {
		t.Errorf("expected name 'test-webhook.example.com', got '%s'", wh.Name)
	}
	if wh.Type != "validating" {
		t.Errorf("expected type 'validating', got '%s'", wh.Type)
	}
	if len(wh.Sources) != 1 {
		t.Fatalf("expected 1 source, got %d", len(wh.Sources))
	}
	if wh.Sources[0].Type != "yaml_manifest" {
		t.Errorf("expected source type 'yaml_manifest', got '%s'", wh.Sources[0].Type)
	}
	if wh.Sources[0].File != "webhook-config.yaml" {
		t.Errorf("expected source file 'webhook-config.yaml', got '%s'", wh.Sources[0].File)
	}
	if wh.FailurePolicy != "Fail" {
		t.Errorf("expected failurePolicy 'Fail', got '%s'", wh.FailurePolicy)
	}
	if wh.ServiceRef != "test-ns/test-service" {
		t.Errorf("expected serviceRef 'test-ns/test-service', got '%s'", wh.ServiceRef)
	}
	if wh.Path != "/validate" {
		t.Errorf("expected path '/validate', got '%s'", wh.Path)
	}
	if len(wh.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(wh.Rules))
	}
}

func TestExtractWebhooks_KubebuilderMarkerSource(t *testing.T) {
	// Create a temporary directory with a Go file containing kubebuilder webhook marker
	tmpDir := t.TempDir()
	goContent := `package main

//+kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,failurePolicy=fail,sideEffects=None,name=mpod.kb.io

func main() {
	// webhook handler code
}
`
	goPath := filepath.Join(tmpDir, "pod_webhook.go")
	if err := os.WriteFile(goPath, []byte(goContent), 0644); err != nil {
		t.Fatalf("failed to write test Go file: %v", err)
	}

	// Run extractWebhooks
	webhooks := extractWebhooks(tmpDir)

	// Verify results
	if len(webhooks) != 1 {
		t.Fatalf("expected 1 webhook, got %d", len(webhooks))
	}

	wh := webhooks[0]
	if wh.Name != "mpod.kb.io" {
		t.Errorf("expected name 'mpod.kb.io', got '%s'", wh.Name)
	}
	if wh.Type != "mutating" {
		t.Errorf("expected type 'mutating', got '%s'", wh.Type)
	}
	if len(wh.Sources) != 1 {
		t.Fatalf("expected 1 source, got %d", len(wh.Sources))
	}
	if wh.Sources[0].Type != "kubebuilder_marker" {
		t.Errorf("expected source type 'kubebuilder_marker', got '%s'", wh.Sources[0].Type)
	}
	if wh.Sources[0].File != "pod_webhook.go" {
		t.Errorf("expected source file 'pod_webhook.go', got '%s'", wh.Sources[0].File)
	}
	if wh.Path != "/mutate-v1-pod" {
		t.Errorf("expected path '/mutate-v1-pod', got '%s'", wh.Path)
	}
	if wh.FailurePolicy != "fail" {
		t.Errorf("expected failurePolicy 'fail', got '%s'", wh.FailurePolicy)
	}
}

func TestExtractWebhooks_MultiSource(t *testing.T) {
	// Create a temporary directory with both YAML and Go kubebuilder marker
	// for the same webhook (same name) to test multi-source tracking
	tmpDir := t.TempDir()

	// YAML manifest
	yamlContent := `apiVersion: admissionregistration.k8s.io/v1
kind: MutatingWebhookConfiguration
metadata:
  name: multi-source-config
webhooks:
  - name: multi.example.com
    failurePolicy: Ignore
    sideEffects: None
    clientConfig:
      service:
        name: webhook-service
        namespace: default
        path: /mutate
    rules:
      - apiGroups: [""]
        apiVersions: ["v1"]
        resources: ["pods"]
        operations: ["CREATE"]
`
	yamlPath := filepath.Join(tmpDir, "webhook.yaml")
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write test YAML: %v", err)
	}

	// Go file with kubebuilder marker (different webhook name for this test)
	goContent := `package main

//+kubebuilder:webhook:path=/validate,mutating=false,failurePolicy=fail,name=validator.example.com

func main() {}
`
	goPath := filepath.Join(tmpDir, "validator_webhook.go")
	if err := os.WriteFile(goPath, []byte(goContent), 0644); err != nil {
		t.Fatalf("failed to write test Go file: %v", err)
	}

	// Run extractWebhooks
	webhooks := extractWebhooks(tmpDir)

	// Verify results: should have 2 webhooks with different names
	if len(webhooks) != 2 {
		t.Fatalf("expected 2 webhooks, got %d", len(webhooks))
	}

	// Find the YAML-sourced webhook
	var yamlWebhook *WebhookConfig
	var goWebhook *WebhookConfig
	for i := range webhooks {
		if webhooks[i].Name == "multi.example.com" {
			yamlWebhook = &webhooks[i]
		} else if webhooks[i].Name == "validator.example.com" {
			goWebhook = &webhooks[i]
		}
	}

	if yamlWebhook == nil {
		t.Fatal("YAML webhook not found")
	}
	if goWebhook == nil {
		t.Fatal("Go webhook not found")
	}

	// Verify YAML webhook
	if len(yamlWebhook.Sources) != 1 {
		t.Errorf("expected 1 source for YAML webhook, got %d", len(yamlWebhook.Sources))
	}
	if yamlWebhook.Sources[0].Type != "yaml_manifest" {
		t.Errorf("expected YAML source type 'yaml_manifest', got '%s'", yamlWebhook.Sources[0].Type)
	}

	// Verify Go webhook
	if len(goWebhook.Sources) != 1 {
		t.Errorf("expected 1 source for Go webhook, got %d", len(goWebhook.Sources))
	}
	if goWebhook.Sources[0].Type != "kubebuilder_marker" {
		t.Errorf("expected Go source type 'kubebuilder_marker', got '%s'", goWebhook.Sources[0].Type)
	}
}

func TestExtractWebhooks_CRDConversion(t *testing.T) {
	// Create a temporary directory with a CRD that has conversion webhook
	tmpDir := t.TempDir()
	crdContent := `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: inferenceservices.serving.kserve.io
spec:
  group: serving.kserve.io
  names:
    kind: InferenceService
  scope: Namespaced
  conversion:
    strategy: Webhook
    webhook:
      clientConfig:
        service:
          name: kserve-webhook-server-service
          namespace: kserve
          path: /convert
      conversionReviewVersions:
      - v1
      - v1beta1
  versions:
  - name: v1
    served: true
    storage: true
`
	crdPath := filepath.Join(tmpDir, "config", "crd", "bases")
	if err := os.MkdirAll(crdPath, 0755); err != nil {
		t.Fatalf("failed to create CRD directory: %v", err)
	}
	crdFile := filepath.Join(crdPath, "inferenceservice.yaml")
	if err := os.WriteFile(crdFile, []byte(crdContent), 0644); err != nil {
		t.Fatalf("failed to write test CRD: %v", err)
	}

	// Run extractWebhooks
	webhooks := extractWebhooks(tmpDir)

	// Verify results
	if len(webhooks) != 1 {
		t.Fatalf("expected 1 webhook, got %d", len(webhooks))
	}

	wh := webhooks[0]
	if wh.Name != "conversion-inferenceservice" {
		t.Errorf("expected name 'conversion-inferenceservice', got '%s'", wh.Name)
	}
	if wh.Type != "conversion" {
		t.Errorf("expected type 'conversion', got '%s'", wh.Type)
	}
	if wh.ConversionCRD != "InferenceService" {
		t.Errorf("expected conversion_crd 'InferenceService', got '%s'", wh.ConversionCRD)
	}
	if wh.Path != "/convert" {
		t.Errorf("expected path '/convert', got '%s'", wh.Path)
	}
	if wh.ServiceRef != "kserve/kserve-webhook-server-service" {
		t.Errorf("expected serviceRef 'kserve/kserve-webhook-server-service', got '%s'", wh.ServiceRef)
	}
	if len(wh.Sources) != 1 {
		t.Fatalf("expected 1 source, got %d", len(wh.Sources))
	}
	if wh.Sources[0].Type != "yaml_manifest" {
		t.Errorf("expected source type 'yaml_manifest', got '%s'", wh.Sources[0].Type)
	}
	expectedSourceFile := filepath.Join("config", "crd", "bases", "inferenceservice.yaml")
	if wh.Sources[0].File != expectedSourceFile {
		t.Errorf("expected source file '%s', got '%s'", expectedSourceFile, wh.Sources[0].File)
	}
}

func TestMapGoHandlers_KubebuilderCorrelation(t *testing.T) {
	// Create a temporary directory with a Go file containing a kubebuilder marker.
	// After extractWebhooks + enrichWebhooks, the webhook should get a go_handler source
	// pointing to the same file as the kubebuilder_marker source.
	tmpDir := t.TempDir()
	goContent := `package main

//+kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,failurePolicy=fail,sideEffects=None,name=mpod.kb.io

func main() {}
`
	goPath := filepath.Join(tmpDir, "pod_webhook.go")
	if err := os.WriteFile(goPath, []byte(goContent), 0644); err != nil {
		t.Fatalf("failed to write test Go file: %v", err)
	}

	webhooks := extractWebhooks(tmpDir)
	enrichWebhooks(webhooks, tmpDir)

	if len(webhooks) != 1 {
		t.Fatalf("expected 1 webhook, got %d", len(webhooks))
	}

	wh := webhooks[0]
	// Should have both kubebuilder_marker and go_handler sources.
	var hasMarker, hasHandler bool
	for _, src := range wh.Sources {
		switch src.Type {
		case "kubebuilder_marker":
			hasMarker = true
		case "go_handler":
			hasHandler = true
			if src.File != "pod_webhook.go" {
				t.Errorf("expected go_handler file 'pod_webhook.go', got '%s'", src.File)
			}
		}
	}
	if !hasMarker {
		t.Error("missing kubebuilder_marker source")
	}
	if !hasHandler {
		t.Error("missing go_handler source after enrichment")
	}
}

func TestMapGoHandlers_HookServerRegister(t *testing.T) {
	// Create a webhook slice manually, then a Go file with hookServer.Register.
	// mapGoHandlers should add a go_handler source.
	tmpDir := t.TempDir()

	goContent := `package webhook

import "sigs.k8s.io/controller-runtime/pkg/webhook"

func SetupWebhooks(mgr ctrl.Manager) {
	hookServer := mgr.GetWebhookServer()
	hookServer.Register("/mutate-pods", &handler{})
}
`
	goPath := filepath.Join(tmpDir, "webhook_setup.go")
	if err := os.WriteFile(goPath, []byte(goContent), 0644); err != nil {
		t.Fatalf("failed to write test Go file: %v", err)
	}

	webhooks := []WebhookConfig{
		{
			Name: "pod-mutator",
			Type: "mutating",
			Path: "/mutate-pods",
			Sources: []SourceRef{
				{Type: "yaml_manifest", File: "config/webhook/manifests.yaml"},
			},
			Rules: []WebhookRule{},
		},
	}

	mapGoHandlers(webhooks, tmpDir)

	if len(webhooks[0].Sources) != 2 {
		t.Fatalf("expected 2 sources, got %d: %+v", len(webhooks[0].Sources), webhooks[0].Sources)
	}

	goHandler := webhooks[0].Sources[1]
	if goHandler.Type != "go_handler" {
		t.Errorf("expected source type 'go_handler', got '%s'", goHandler.Type)
	}
	if goHandler.File != "webhook_setup.go" {
		t.Errorf("expected source file 'webhook_setup.go', got '%s'", goHandler.File)
	}
}

func TestMapGoHandlers_ExcludesTestFiles(t *testing.T) {
	// A _test.go file with hookServer.Register should not produce a go_handler source.
	tmpDir := t.TempDir()

	testContent := `package webhook

func TestWebhookSetup(t *testing.T) {
	hookServer.Register("/mutate-pods", &handler{})
}
`
	testPath := filepath.Join(tmpDir, "webhook_setup_test.go")
	if err := os.WriteFile(testPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	webhooks := []WebhookConfig{
		{
			Name: "pod-mutator",
			Type: "mutating",
			Path: "/mutate-pods",
			Sources: []SourceRef{
				{Type: "yaml_manifest", File: "config/webhook/manifests.yaml"},
			},
			Rules: []WebhookRule{},
		},
	}

	mapGoHandlers(webhooks, tmpDir)

	if len(webhooks[0].Sources) != 1 {
		t.Fatalf("expected 1 source (no go_handler from test file), got %d: %+v",
			len(webhooks[0].Sources), webhooks[0].Sources)
	}
}

func TestExtractDataRead_BasicPattern(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a handler file with multiple client.Get/List calls.
	handlerContent := `package webhook

import (
	corev1 "k8s.io/api/core/v1"
	servingv1beta1 "github.com/kserve/kserve/pkg/apis/serving/v1beta1"
)

func (h *Handler) Handle(ctx context.Context, req admission.Request) admission.Response {
	if err := h.client.Get(ctx, key, &corev1.Secret{}); err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	if err := h.Client.Get(ctx, key, &servingv1beta1.InferenceService{}); err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	if err := h.client.List(ctx, &corev1.ServiceList{}); err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.Allowed("")
}
`
	handlerPath := filepath.Join(tmpDir, "webhook_handler.go")
	if err := os.WriteFile(handlerPath, []byte(handlerContent), 0644); err != nil {
		t.Fatalf("failed to write handler file: %v", err)
	}

	webhooks := []WebhookConfig{
		{
			Name: "test-webhook",
			Type: "validating",
			Path: "/validate",
			Sources: []SourceRef{
				{Type: "yaml_manifest", File: "config/webhook/manifests.yaml"},
				{Type: "go_handler", File: "webhook_handler.go"},
			},
			Rules: []WebhookRule{},
		},
	}

	extractDataRead(webhooks, tmpDir)

	if len(webhooks[0].DataRead) != 3 {
		t.Fatalf("expected 3 TypeRefs, got %d: %+v", len(webhooks[0].DataRead), webhooks[0].DataRead)
	}

	// Verify Secret has Group="" and GroupKnown=true.
	var foundSecret, foundInferenceService, foundService bool
	for _, tr := range webhooks[0].DataRead {
		if tr.Kind == "Secret" {
			foundSecret = true
			if tr.Group != "" {
				t.Errorf("expected Secret Group to be empty, got '%s'", tr.Group)
			}
			if !tr.GroupKnown {
				t.Error("expected Secret GroupKnown to be true")
			}
		}
		if tr.Kind == "InferenceService" {
			foundInferenceService = true
			if tr.Group != "serving.kserve.io" {
				t.Errorf("expected InferenceService Group to be 'serving.kserve.io', got '%s'", tr.Group)
			}
			if !tr.GroupKnown {
				t.Error("expected InferenceService GroupKnown to be true")
			}
		}
		if tr.Kind == "Service" {
			foundService = true
			if tr.Group != "" {
				t.Errorf("expected Service Group to be empty, got '%s'", tr.Group)
			}
			if !tr.GroupKnown {
				t.Error("expected Service GroupKnown to be true")
			}
		}
	}

	if !foundSecret {
		t.Error("Secret not found in DataRead")
	}
	if !foundInferenceService {
		t.Error("InferenceService not found in DataRead")
	}
	if !foundService {
		t.Error("Service not found in DataRead")
	}
}

func TestExtractDataRead_SkipsComments(t *testing.T) {
	tmpDir := t.TempDir()

	handlerContent := `package webhook

import (
	corev1 "k8s.io/api/core/v1"
)

func (h *Handler) Handle(ctx context.Context, req admission.Request) admission.Response {
	// h.client.Get(ctx, key, &corev1.Secret{})
	if err := h.client.Get(ctx, key, &corev1.ConfigMap{}); err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	return admission.Allowed("")
}
`
	handlerPath := filepath.Join(tmpDir, "webhook_handler.go")
	if err := os.WriteFile(handlerPath, []byte(handlerContent), 0644); err != nil {
		t.Fatalf("failed to write handler file: %v", err)
	}

	webhooks := []WebhookConfig{
		{
			Name: "test-webhook",
			Type: "validating",
			Path: "/validate",
			Sources: []SourceRef{
				{Type: "go_handler", File: "webhook_handler.go"},
			},
			Rules: []WebhookRule{},
		},
	}

	extractDataRead(webhooks, tmpDir)

	if len(webhooks[0].DataRead) != 1 {
		t.Fatalf("expected 1 TypeRef, got %d: %+v", len(webhooks[0].DataRead), webhooks[0].DataRead)
	}

	if webhooks[0].DataRead[0].Kind != "ConfigMap" {
		t.Errorf("expected Kind 'ConfigMap', got '%s'", webhooks[0].DataRead[0].Kind)
	}
}

func TestExtractDataRead_UnknownPackage(t *testing.T) {
	tmpDir := t.TempDir()

	handlerContent := `package webhook

import (
	customv1 "example.com/custom/v1"
)

func (h *Handler) Handle(ctx context.Context, req admission.Request) admission.Response {
	if err := h.client.Get(ctx, key, &customv1.MyResource{}); err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	return admission.Allowed("")
}
`
	handlerPath := filepath.Join(tmpDir, "webhook_handler.go")
	if err := os.WriteFile(handlerPath, []byte(handlerContent), 0644); err != nil {
		t.Fatalf("failed to write handler file: %v", err)
	}

	webhooks := []WebhookConfig{
		{
			Name: "test-webhook",
			Type: "validating",
			Path: "/validate",
			Sources: []SourceRef{
				{Type: "go_handler", File: "webhook_handler.go"},
			},
			Rules: []WebhookRule{},
		},
	}

	extractDataRead(webhooks, tmpDir)

	if len(webhooks[0].DataRead) != 1 {
		t.Fatalf("expected 1 TypeRef, got %d: %+v", len(webhooks[0].DataRead), webhooks[0].DataRead)
	}

	tr := webhooks[0].DataRead[0]
	if tr.Kind != "MyResource" {
		t.Errorf("expected Kind 'MyResource', got '%s'", tr.Kind)
	}
	if tr.GroupKnown {
		t.Error("expected GroupKnown to be false for unknown package")
	}
	if tr.Group != "" {
		t.Errorf("expected Group to be empty for unknown package, got '%s'", tr.Group)
	}
}

func TestExtractEnableConditions_EnvVar(t *testing.T) {
	// Create a temporary directory with a Go file containing webhook setup with env var guard.
	tmpDir := t.TempDir()
	goContent := `package main

import (
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
)

func setupWebhooks(mgr ctrl.Manager) {
	if os.Getenv("ENABLE_WEBHOOKS") == "false" {
		return
	}
	mgr.GetWebhookServer().Register("/validate-model", &handler{})
}
`
	goPath := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(goPath, []byte(goContent), 0644); err != nil {
		t.Fatalf("failed to write test Go file: %v", err)
	}

	webhooks := []WebhookConfig{
		{
			Name: "model-validator",
			Type: "validating",
			Path: "/validate-model",
			Sources: []SourceRef{
				{Type: "yaml_manifest", File: "config/webhook/manifests.yaml"},
			},
			Rules: []WebhookRule{},
		},
	}

	extractEnableConditions(webhooks, tmpDir)

	if webhooks[0].EnableCondition == "" {
		t.Error("expected EnableCondition to be set, got empty string")
	}
	if !strings.Contains(webhooks[0].EnableCondition, "os.Getenv") {
		t.Errorf("expected EnableCondition to contain 'os.Getenv', got: %s", webhooks[0].EnableCondition)
	}
}

func TestExtractEnableConditions_FeatureGate(t *testing.T) {
	// Create a temporary directory with a Go file containing feature gate guard.
	tmpDir := t.TempDir()
	goContent := `package main

import ctrl "sigs.k8s.io/controller-runtime"

func SetupWithManager(mgr ctrl.Manager) error {
	if featureGate.Enabled(features.WebhookValidation) {
		mgr.GetWebhookServer().Register("/validate-isvc", &handler{})
	}
	return nil
}
`
	goPath := filepath.Join(tmpDir, "webhook_setup.go")
	if err := os.WriteFile(goPath, []byte(goContent), 0644); err != nil {
		t.Fatalf("failed to write test Go file: %v", err)
	}

	webhooks := []WebhookConfig{
		{
			Name: "isvc-validator",
			Type: "validating",
			Path: "/validate-isvc",
			Sources: []SourceRef{
				{Type: "yaml_manifest", File: "config/webhook/manifests.yaml"},
			},
			Rules: []WebhookRule{},
		},
	}

	extractEnableConditions(webhooks, tmpDir)

	if webhooks[0].EnableCondition == "" {
		t.Error("expected EnableCondition to be set, got empty string")
	}
	if !strings.Contains(webhooks[0].EnableCondition, "featureGate.Enabled") {
		t.Errorf("expected EnableCondition to contain 'featureGate.Enabled', got: %s", webhooks[0].EnableCondition)
	}
}

func TestExtractEnableConditions_ManagementState(t *testing.T) {
	// Create a temporary directory with a Go file containing management state guard.
	tmpDir := t.TempDir()
	goContent := `package main

import ctrl "sigs.k8s.io/controller-runtime"

func setupWebhooks(mgr ctrl.Manager, dsc *dscv1.DataScienceCluster) {
	if dsc.Spec.Components.Kserve.ManagementState == operatorv1.Managed {
		mgr.GetWebhookServer().Register("/mutate-isvc", &handler{})
	}
}
`
	goPath := filepath.Join(tmpDir, "webhooks.go")
	if err := os.WriteFile(goPath, []byte(goContent), 0644); err != nil {
		t.Fatalf("failed to write test Go file: %v", err)
	}

	webhooks := []WebhookConfig{
		{
			Name: "isvc-mutator",
			Type: "mutating",
			Path: "/mutate-isvc",
			Sources: []SourceRef{
				{Type: "yaml_manifest", File: "config/webhook/manifests.yaml"},
			},
			Rules: []WebhookRule{},
		},
	}

	extractEnableConditions(webhooks, tmpDir)

	if webhooks[0].EnableCondition == "" {
		t.Error("expected EnableCondition to be set, got empty string")
	}
	if !strings.Contains(webhooks[0].EnableCondition, "ManagementState") {
		t.Errorf("expected EnableCondition to contain 'ManagementState', got: %s", webhooks[0].EnableCondition)
	}
}

func TestExtractEnableConditions_BooleanFlag(t *testing.T) {
	tmpDir := t.TempDir()
	goContent := `package main

import ctrl "sigs.k8s.io/controller-runtime"

func setupWebhooks(mgr ctrl.Manager, enableWebhooks bool) {
	if enableWebhooks {
		mgr.GetWebhookServer().Register("/validate-model", &handler{})
	}
}
`
	goPath := filepath.Join(tmpDir, "webhook_setup.go")
	if err := os.WriteFile(goPath, []byte(goContent), 0644); err != nil {
		t.Fatalf("failed to write test Go file: %v", err)
	}

	webhooks := []WebhookConfig{
		{
			Name:    "model-validator",
			Type:    "validating",
			Path:    "/validate-model",
			Sources: []SourceRef{{Type: "yaml_manifest", File: "config/webhook/manifests.yaml"}},
			Rules:   []WebhookRule{},
		},
	}

	extractEnableConditions(webhooks, tmpDir)

	if webhooks[0].EnableCondition == "" {
		t.Error("expected EnableCondition to be set, got empty string")
	}
	// Must not contain trailing punctuation like { or &
	cond := webhooks[0].EnableCondition
	if strings.ContainsAny(cond, "{&|") {
		t.Errorf("EnableCondition should not contain trailing punctuation, got: %q", cond)
	}
	if !strings.Contains(cond, "enableWebhooks") {
		t.Errorf("expected EnableCondition to contain 'enableWebhooks', got: %s", cond)
	}
}
