package extractor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractPythonK8sCalls(t *testing.T) {
	tests := []struct {
		name          string
		files         map[string]string
		wantLen       int
		wantAPI       string
		wantOperation string
		wantResource  string
	}{
		{
			name: "kubernetes client import",
			files: map[string]string{
				"setup.py": `from kubernetes import client, config`,
			},
			wantLen:       1,
			wantAPI:       "kubernetes",
			wantOperation: "import",
		},
		{
			name: "CustomObjectsApi import",
			files: map[string]string{
				"k8s_utils.py": `from kubernetes.client import CustomObjectsApi`,
			},
			wantLen:       1,
			wantAPI:       "CustomObjectsApi",
			wantOperation: "import",
		},
		{
			name: "CoreV1Api import",
			files: map[string]string{
				"k8s_utils.py": `from kubernetes.client import CoreV1Api`,
			},
			wantLen:       1,
			wantAPI:       "CoreV1Api",
			wantOperation: "import",
		},
		{
			name: "DynamicClient import",
			files: map[string]string{
				"ocp.py": `from openshift.dynamic import DynamicClient`,
			},
			wantLen:       1,
			wantAPI:       "DynamicClient",
			wantOperation: "import",
		},
		{
			name: "create_namespaced_custom_object",
			files: map[string]string{
				"queue.py": `api.create_namespaced_custom_object(group, version, namespace, "localqueues", body)`,
			},
			wantLen:       1,
			wantAPI:       "CustomObjectsApi",
			wantOperation: "create",
			wantResource:  "custom_object",
		},
		{
			name: "list_namespaced_custom_object",
			files: map[string]string{
				"watcher.py": `api.list_namespaced_custom_object(group, version, namespace, "localqueues")`,
			},
			wantLen:       1,
			wantAPI:       "CustomObjectsApi",
			wantOperation: "list",
			wantResource:  "custom_object",
		},
		{
			name: "delete_cluster_custom_object",
			files: map[string]string{
				"cleanup.py": `api.delete_cluster_custom_object(group, version, "clusterqueues", name)`,
			},
			wantLen:       1,
			wantAPI:       "CustomObjectsApi",
			wantOperation: "delete",
			wantResource:  "custom_object",
		},
		{
			name: "create_namespaced_config_map",
			files: map[string]string{
				"config.py": `v1.create_namespaced_config_map(namespace, body)`,
			},
			wantLen:       1,
			wantAPI:       "CoreV1Api",
			wantOperation: "create",
			wantResource:  "config_map",
		},
		{
			name: "create_namespaced_deployment",
			files: map[string]string{
				"deploy.py": `apps_v1.create_namespaced_deployment(namespace, body)`,
			},
			wantLen:       1,
			wantAPI:       "AppsV1Api",
			wantOperation: "create",
			wantResource:  "deployment",
		},
		{
			name: "CRD kind reference LocalQueue",
			files: map[string]string{
				"queue.py": `body = {"kind": "LocalQueue", "apiVersion": "kueue.x-k8s.io/v1beta1"}`,
			},
			wantLen:       1,
			wantAPI:       "CRD",
			wantOperation: "kind_ref",
			wantResource:  "LocalQueue",
		},
		{
			name: "CRD kind via assignment",
			files: map[string]string{
				"inference.py": `kind = "InferenceService"`,
			},
			wantLen:       1,
			wantAPI:       "CRD",
			wantOperation: "kind_ref",
			wantResource:  "InferenceService",
		},
		{
			name: "CRD kind RayCluster",
			files: map[string]string{
				"ray.py": `kind = "RayCluster"`,
			},
			wantLen:       1,
			wantAPI:       "CRD",
			wantOperation: "kind_ref",
			wantResource:  "RayCluster",
		},
		{
			name: "unknown kind is ignored",
			files: map[string]string{
				"random.py": `kind = "SomethingRandom"`,
			},
			wantLen: 0,
		},
		{
			name: "skips comments",
			files: map[string]string{
				"commented.py": `# api.create_namespaced_custom_object(group, version, namespace, plural, body)`,
			},
			wantLen: 0,
		},
		{
			name: "skips test files",
			files: map[string]string{
				"test_k8s.py": `api.create_namespaced_custom_object(group, version, namespace, "localqueues", body)`,
			},
			wantLen: 0,
		},
		{
			name: "multiple patterns in one file",
			files: map[string]string{
				"k8s_ops.py": `from kubernetes.client import CustomObjectsApi, CoreV1Api
api = CustomObjectsApi()
api.create_namespaced_custom_object(group, version, namespace, "localqueues", body)
v1 = CoreV1Api()
v1.create_namespaced_config_map(namespace, cm_body)
body = {"kind": "LocalQueue", "apiVersion": "kueue.x-k8s.io/v1beta1"}`,
			},
			wantLen: 5, // 2 imports + 2 API calls + 1 kind ref
		},
		{
			name: "deduplication of same call",
			files: map[string]string{
				"dup.py": `api.create_namespaced_custom_object(group, version, ns1, "localqueues", body1)
api.create_namespaced_custom_object(group, version, ns2, "clusterqueues", body2)`,
			},
			wantLen: 1, // same api+operation+resource = deduplicated
		},
		{
			name: "create_namespaced_secret",
			files: map[string]string{
				"secrets.py": `v1.create_namespaced_secret(namespace, secret_body)`,
			},
			wantLen:       1,
			wantAPI:       "CoreV1Api",
			wantOperation: "create",
			wantResource:  "secret",
		},
		{
			name: "patch_namespaced_custom_object",
			files: map[string]string{
				"patch.py": `api.patch_namespaced_custom_object(group, version, namespace, "localqueues", name, patch_body)`,
			},
			wantLen:       1,
			wantAPI:       "CustomObjectsApi",
			wantOperation: "patch",
			wantResource:  "custom_object",
		},
		{
			name: "PyTorchJob kind ref",
			files: map[string]string{
				"training.py": `kind = "PyTorchJob"`,
			},
			wantLen:       1,
			wantAPI:       "CRD",
			wantOperation: "kind_ref",
			wantResource:  "PyTorchJob",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			for relPath, content := range tt.files {
				fullPath := filepath.Join(tmpDir, relPath)
				if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
					t.Fatal(err)
				}
			}

			calls := extractPythonK8sCalls(tmpDir)
			if len(calls) != tt.wantLen {
				t.Errorf("got %d calls, want %d", len(calls), tt.wantLen)
				for i, c := range calls {
					t.Logf("  call[%d]: api=%s op=%s resource=%s source=%s", i, c.API, c.Operation, c.Resource, c.Source)
				}
				return
			}

			if tt.wantLen == 0 {
				return
			}

			// For single-result tests, check the first (or only) call
			if tt.wantLen == 1 {
				c := calls[0]
				if tt.wantAPI != "" && c.API != tt.wantAPI {
					t.Errorf("API = %q, want %q", c.API, tt.wantAPI)
				}
				if tt.wantOperation != "" && c.Operation != tt.wantOperation {
					t.Errorf("Operation = %q, want %q", c.Operation, tt.wantOperation)
				}
				if tt.wantResource != "" && c.Resource != tt.wantResource {
					t.Errorf("Resource = %q, want %q", c.Resource, tt.wantResource)
				}
				if c.Source == "" {
					t.Error("Source should not be empty")
				}
			}
		})
	}
}

func TestExtractPythonK8sCalls_SkipsVenvAndPycache(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file in venv that should be skipped
	venvDir := filepath.Join(tmpDir, "venv", "lib", "python3.11")
	os.MkdirAll(venvDir, 0o755)
	os.WriteFile(filepath.Join(venvDir, "k8s.py"),
		[]byte(`api.create_namespaced_custom_object(g, v, ns, "queues", body)`), 0o644)

	// Create a file in __pycache__ that should be skipped
	cacheDir := filepath.Join(tmpDir, "__pycache__")
	os.MkdirAll(cacheDir, 0o755)
	os.WriteFile(filepath.Join(cacheDir, "k8s.py"),
		[]byte(`api.create_namespaced_custom_object(g, v, ns, "queues", body)`), 0o644)

	// Create a file in site-packages that should be skipped
	siteDir := filepath.Join(tmpDir, "lib", "site-packages")
	os.MkdirAll(siteDir, 0o755)
	os.WriteFile(filepath.Join(siteDir, "k8s.py"),
		[]byte(`api.create_namespaced_custom_object(g, v, ns, "queues", body)`), 0o644)

	calls := extractPythonK8sCalls(tmpDir)
	if len(calls) != 0 {
		t.Errorf("expected 0 calls from skipped dirs, got %d", len(calls))
		for i, c := range calls {
			t.Logf("  call[%d]: %+v", i, c)
		}
	}
}

func TestInferAPIFromResource(t *testing.T) {
	tests := []struct {
		resource string
		wantAPI  string
	}{
		{"custom_object", "CustomObjectsApi"},
		{"deployment", "AppsV1Api"},
		{"stateful_set", "AppsV1Api"},
		{"daemon_set", "AppsV1Api"},
		{"config_map", "CoreV1Api"},
		{"secret", "CoreV1Api"},
		{"service", "CoreV1Api"},
		{"pod", "CoreV1Api"},
		{"job", "BatchV1Api"},
		{"cron_job", "BatchV1Api"},
		{"ingress", "NetworkingV1Api"},
		{"network_policy", "NetworkingV1Api"},
	}

	for _, tt := range tests {
		t.Run(tt.resource, func(t *testing.T) {
			got := inferAPIFromResource(tt.resource)
			if got != tt.wantAPI {
				t.Errorf("inferAPIFromResource(%q) = %q, want %q", tt.resource, got, tt.wantAPI)
			}
		})
	}
}
