package extractor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractLabelContractsYAML(t *testing.T) {
	dir := t.TempDir()

	// Create a Job manifest with kueue labels
	jobYAML := `apiVersion: batch/v1
kind: Job
metadata:
  name: training-job
  labels:
    kueue.x-k8s.io/queue-name: user-queue
spec:
  template:
    metadata:
      labels:
        kueue.x-k8s.io/queue-name: user-queue
        kueue.x-k8s.io/priority-class: high-priority
    spec:
      containers:
      - name: trainer
        image: trainer:latest
`
	os.MkdirAll(filepath.Join(dir, "config", "jobs"), 0o755)
	os.WriteFile(filepath.Join(dir, "config", "jobs", "training-job.yaml"), []byte(jobYAML), 0o644)

	// Create a Deployment with Istio annotation
	depYAML := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
  annotations:
    sidecar.istio.io/inject: "true"
spec:
  template:
    metadata:
      labels:
        app.kubernetes.io/managed-by: kustomize
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "8080"
    spec:
      containers:
      - name: app
        image: app:latest
`
	os.WriteFile(filepath.Join(dir, "config", "jobs", "deployment.yaml"), []byte(depYAML), 0o644)

	contracts := extractLabelContracts(dir)

	if len(contracts) == 0 {
		t.Fatal("expected label contracts, got none")
	}

	// Build lookup map
	type key struct {
		label   string
		context string
	}
	found := make(map[key]LabelContract)
	for _, c := range contracts {
		found[key{c.Label, c.Context}] = c
	}

	// Verify kueue queue-name on the Job resource itself
	if c, ok := found[key{"kueue.x-k8s.io/queue-name", "job"}]; !ok {
		t.Error("missing kueue queue-name contract on job metadata")
	} else if c.Integration != "Kueue" {
		t.Errorf("expected integration Kueue, got %q", c.Integration)
	} else if c.Type != "label" {
		t.Errorf("expected type label, got %q", c.Type)
	}

	// Verify kueue labels on job template
	if _, ok := found[key{"kueue.x-k8s.io/queue-name", "job-template"}]; !ok {
		t.Error("missing kueue queue-name contract on job-template")
	}
	if _, ok := found[key{"kueue.x-k8s.io/priority-class", "job-template"}]; !ok {
		t.Error("missing kueue priority-class contract on job-template")
	}

	// Verify Istio annotation on deployment
	if c, ok := found[key{"sidecar.istio.io/inject", "deployment"}]; !ok {
		t.Error("missing istio sidecar inject contract on deployment")
	} else if c.Type != "annotation" {
		t.Errorf("expected type annotation, got %q", c.Type)
	} else if c.Integration != "Istio" {
		t.Errorf("expected integration Istio, got %q", c.Integration)
	}

	// Verify prometheus annotations on deployment template
	if c, ok := found[key{"prometheus.io/scrape", "deployment-template"}]; !ok {
		t.Error("missing prometheus scrape contract on deployment-template")
	} else if c.Type != "annotation" {
		t.Errorf("expected type annotation, got %q", c.Type)
	}

	// Verify managed-by label on deployment template
	if _, ok := found[key{"app.kubernetes.io/managed-by", "deployment-template"}]; !ok {
		t.Error("missing managed-by contract on deployment-template")
	}
}

func TestExtractLabelContractsGoSource(t *testing.T) {
	dir := t.TempDir()

	goSource := `package controller

const (
	queueLabel = "kueue.x-k8s.io/queue-name"
	inferenceLabel = "serving.kserve.io/inferenceservice"
)

func setLabels(labels map[string]string) {
	labels["kueue.x-k8s.io/priority-class"] = "default"
}
`
	os.MkdirAll(filepath.Join(dir, "pkg", "controller"), 0o755)
	os.WriteFile(filepath.Join(dir, "pkg", "controller", "labels.go"), []byte(goSource), 0o644)

	contracts := extractLabelContracts(dir)

	if len(contracts) == 0 {
		t.Fatal("expected label contracts from Go source, got none")
	}

	foundLabels := make(map[string]bool)
	for _, c := range contracts {
		if c.Context != "source-code" {
			t.Errorf("expected context source-code, got %q", c.Context)
		}
		foundLabels[c.Label] = true
	}

	for _, expected := range []string{
		"kueue.x-k8s.io/queue-name",
		"serving.kserve.io/inferenceservice",
		"kueue.x-k8s.io/priority-class",
	} {
		if !foundLabels[expected] {
			t.Errorf("missing expected label %q from Go source scan", expected)
		}
	}
}

func TestExtractLabelContractsTestFilesExcluded(t *testing.T) {
	dir := t.TempDir()

	// Test files should be skipped
	testSource := `package controller_test

func TestFoo(t *testing.T) {
	labels := map[string]string{
		"kueue.x-k8s.io/queue-name": "test-queue",
	}
	_ = labels
}
`
	os.MkdirAll(filepath.Join(dir, "pkg", "controller"), 0o755)
	os.WriteFile(filepath.Join(dir, "pkg", "controller", "labels_test.go"), []byte(testSource), 0o644)

	contracts := extractLabelContracts(dir)

	for _, c := range contracts {
		if c.Context == "source-code" {
			t.Errorf("should not have found source-code contract from test file: %+v", c)
		}
	}
}

func TestExtractLabelContractsCronJob(t *testing.T) {
	dir := t.TempDir()

	cronJobYAML := `apiVersion: batch/v1
kind: CronJob
metadata:
  name: scheduled-training
spec:
  schedule: "0 2 * * *"
  jobTemplate:
    spec:
      template:
        metadata:
          labels:
            kueue.x-k8s.io/queue-name: batch-queue
        spec:
          containers:
          - name: trainer
            image: trainer:latest
`
	os.WriteFile(filepath.Join(dir, "cronjob.yaml"), []byte(cronJobYAML), 0o644)

	contracts := extractLabelContracts(dir)

	found := false
	for _, c := range contracts {
		if c.Label == "kueue.x-k8s.io/queue-name" && c.Context == "cronjob-template" {
			found = true
			break
		}
	}
	if !found {
		t.Error("missing kueue queue-name contract on cronjob-template")
		for _, c := range contracts {
			t.Logf("  found: label=%s context=%s source=%s", c.Label, c.Context, c.Source)
		}
	}
}
