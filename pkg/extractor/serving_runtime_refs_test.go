package extractor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractServingRuntimeRefsFromYAML(t *testing.T) {
	dir := t.TempDir()

	// Create a ServingRuntime manifest
	srYAML := `apiVersion: serving.kserve.io/v1alpha1
kind: ServingRuntime
metadata:
  name: vllm-runtime
spec:
  containers:
    - name: kserve-container
      image: quay.io/opendatahub/vllm:latest
`
	os.MkdirAll(filepath.Join(dir, "config", "runtimes"), 0o755)
	os.WriteFile(filepath.Join(dir, "config", "runtimes", "vllm.yaml"), []byte(srYAML), 0o644)

	// Create a ClusterServingRuntime manifest
	csrYAML := `apiVersion: serving.kserve.io/v1alpha1
kind: ClusterServingRuntime
metadata:
  name: modelmesh-runtime
spec:
  containers:
    - name: kserve-container
      image: quay.io/modelmesh-serving/modelmesh:latest
`
	os.WriteFile(filepath.Join(dir, "config", "runtimes", "modelmesh.yaml"), []byte(csrYAML), 0o644)

	// Create an InferenceService manifest
	isYAML := `apiVersion: serving.kserve.io/v1beta1
kind: InferenceService
metadata:
  name: my-model
spec:
  predictor:
    containers:
      - name: predictor
        image: quay.io/opendatahub/vllm:latest
  transformer:
    containers:
      - name: transformer
        image: quay.io/opendatahub/transformer:v1
`
	os.MkdirAll(filepath.Join(dir, "manifests"), 0o755)
	os.WriteFile(filepath.Join(dir, "manifests", "inference.yaml"), []byte(isYAML), 0o644)

	refs := extractServingRuntimeRefs(dir)

	// Expect at least 4 refs: vllm, modelmesh, predictor image, transformer image
	if len(refs) < 4 {
		t.Errorf("expected at least 4 refs, got %d", len(refs))
		for _, r := range refs {
			t.Logf("  %s/%s: %s (%s)", r.Kind, r.Name, r.ContainerImage, r.Source)
		}
		return
	}

	// Verify we have the expected kinds
	kindCounts := map[string]int{}
	for _, r := range refs {
		kindCounts[r.Kind]++
	}

	if kindCounts["ServingRuntime"] < 1 {
		t.Error("expected at least 1 ServingRuntime ref")
	}
	if kindCounts["ClusterServingRuntime"] < 1 {
		t.Error("expected at least 1 ClusterServingRuntime ref")
	}
	if kindCounts["InferenceService"] < 2 {
		t.Error("expected at least 2 InferenceService refs (predictor + transformer)")
	}

	// Verify specific images are captured
	imageFound := map[string]bool{}
	for _, r := range refs {
		imageFound[r.ContainerImage] = true
	}
	for _, expected := range []string{
		"quay.io/opendatahub/vllm:latest",
		"quay.io/modelmesh-serving/modelmesh:latest",
		"quay.io/opendatahub/transformer:v1",
	} {
		if !imageFound[expected] {
			t.Errorf("expected image %q not found in refs", expected)
		}
	}
}

func TestExtractServingRuntimeRefsSkipsTemplates(t *testing.T) {
	dir := t.TempDir()

	// YAML with unresolved Helm template should be skipped
	srYAML := `apiVersion: serving.kserve.io/v1alpha1
kind: ServingRuntime
metadata:
  name: templated-runtime
spec:
  containers:
    - name: kserve-container
      image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
`
	os.WriteFile(filepath.Join(dir, "runtime.yaml"), []byte(srYAML), 0o644)

	refs := extractServingRuntimeRefs(dir)

	// The templated image should be skipped (not parseable as YAML due to {{)
	for _, r := range refs {
		if r.Name == "templated-runtime" {
			t.Errorf("expected templated runtime image to be skipped, got: %+v", r)
		}
	}
}

func TestExtractServingRuntimeRefsFromGoSource(t *testing.T) {
	dir := t.TempDir()

	goSource := `package main

func createRuntime() {
	runtime := &ServingRuntime{
		Kind: "ServingRuntime",
		Spec: ServingRuntimeSpec{
			Containers: []Container{
				{
					Image: "quay.io/opendatahub/vllm-cpu:v0.6.2",
				},
			},
		},
	}
	return runtime
}
`
	os.WriteFile(filepath.Join(dir, "controller.go"), []byte(goSource), 0o644)

	refs := extractServingRuntimeRefsFromGo(dir)

	if len(refs) == 0 {
		t.Fatal("expected at least 1 ref from Go source")
	}

	found := false
	for _, r := range refs {
		if r.Kind == "ServingRuntime" && r.ContainerImage == "quay.io/opendatahub/vllm-cpu:v0.6.2" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected Go source ref with Kind=ServingRuntime and image=quay.io/opendatahub/vllm-cpu:v0.6.2")
		for _, r := range refs {
			t.Logf("  %s/%s: %s (%s)", r.Kind, r.Name, r.ContainerImage, r.Source)
		}
	}
}

func TestExtractServingRuntimeRefsDedup(t *testing.T) {
	refs := []ServingRuntimeRef{
		{Name: "vllm", Kind: "ServingRuntime", ContainerImage: "quay.io/vllm:latest", Source: "a.yaml"},
		{Name: "vllm", Kind: "ServingRuntime", ContainerImage: "quay.io/vllm:latest", Source: "b.go:10"},
	}

	deduped := dedupServingRuntimeRefs(refs)
	if len(deduped) != 1 {
		t.Errorf("expected 1 deduped ref, got %d", len(deduped))
	}
	// Should prefer YAML source over Go source
	if len(deduped) > 0 && deduped[0].Source != "a.yaml" {
		t.Errorf("expected YAML source to be preferred, got %q", deduped[0].Source)
	}
}

func TestExtractServingRuntimeRefsExcludesTestDirs(t *testing.T) {
	dir := t.TempDir()

	// Put a valid ServingRuntime in a test directory (should be excluded)
	os.MkdirAll(filepath.Join(dir, "test"), 0o755)
	testYAML := `apiVersion: serving.kserve.io/v1alpha1
kind: ServingRuntime
metadata:
  name: test-runtime
spec:
  containers:
    - name: kserve-container
      image: quay.io/test/runtime:latest
`
	os.WriteFile(filepath.Join(dir, "test", "runtime.yaml"), []byte(testYAML), 0o644)

	refs := extractServingRuntimeRefs(dir)
	for _, r := range refs {
		if r.Name == "test-runtime" {
			t.Error("expected test directory content to be excluded")
		}
	}
}
