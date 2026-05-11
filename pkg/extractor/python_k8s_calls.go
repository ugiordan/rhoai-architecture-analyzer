package extractor

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// PythonK8sCall represents a Kubernetes API call detected in Python source code.
// These reveal cross-component relationships (e.g., codeflare-sdk creating
// LocalQueue objects implies Kueue integration).
type PythonK8sCall struct {
	API       string `json:"api"`       // "CustomObjectsApi", "CoreV1Api", "AppsV1Api", "DynamicClient", etc.
	Operation string `json:"operation"` // "create", "get", "list", "patch", "delete", "import", "kind_ref"
	Resource  string `json:"resource"`  // "custom_object", "config_map", "deployment", or CRD kind like "LocalQueue"
	Source    string `json:"source"`    // file:line
}

// k8s Python client import patterns.
var pyK8sImportPatterns = []struct {
	re  *regexp.Regexp
	api string
}{
	// from kubernetes import client, config
	{regexp.MustCompile(`^\s*from\s+kubernetes\s+import\s+(.+)`), "kubernetes"},
	// from kubernetes.client import CustomObjectsApi, CoreV1Api, ...
	{regexp.MustCompile(`^\s*from\s+kubernetes\.client\s+import\s+(.+)`), "kubernetes.client"},
	// import kubernetes
	{regexp.MustCompile(`^\s*import\s+kubernetes\b`), "kubernetes"},
	// from openshift.dynamic import DynamicClient
	{regexp.MustCompile(`^\s*from\s+openshift\.dynamic\s+import\s+(.+)`), "openshift.dynamic"},
}

// k8s Python client method call patterns. Each captures the operation and resource
// from method names like create_namespaced_custom_object, list_namespaced_config_map, etc.
var pyK8sMethodPatterns = []struct {
	re        *regexp.Regexp
	apiGroup  int // capture group for API variable (0 = skip)
	operation int // capture group for operation
	resource  int // capture group for resource
}{
	// CustomObjectsApi: verb_namespaced_custom_object or verb_cluster_custom_object
	// e.g., client.create_namespaced_custom_object(...), api.list_cluster_custom_object(...)
	{
		regexp.MustCompile(`\b\w+\.(create|get|list|patch|delete|replace|watch)_(namespaced|cluster)_custom_object\b`),
		0, 1, 0,
	},
	// CoreV1Api / AppsV1Api / etc.: verb_namespaced_resource or verb_resource
	// e.g., v1.create_namespaced_config_map(...), apps_v1.create_namespaced_deployment(...)
	{
		regexp.MustCompile(`\b\w+\.(create|get|list|patch|delete|replace|watch)_(?:namespaced_)?(\w+)\b`),
		0, 1, 2,
	},
}

// Known CRD kind values commonly set in Python code interacting with k8s.
var pyKnownCRDKinds = map[string]bool{
	"LocalQueue":       true,
	"ClusterQueue":     true,
	"ResourceFlavor":   true,
	"Workload":         true,
	"InferenceService": true,
	"ServingRuntime":   true,
	"RayCluster":       true,
	"RayJob":           true,
	"RayService":       true,
	"PyTorchJob":       true,
	"TFJob":            true,
	"MPIJob":           true,
	"XGBoostJob":       true,
	"PaddleJob":        true,
	"AppWrapper":       true,
	"ScheduledWorkflow": true,
	"Notebook":         true,
	"DataScienceCluster": true,
	"DSCInitialization": true,
	"Pipeline":         true,
	"PipelineRun":      true,
	"Task":             true,
	"TaskRun":          true,
	"IsvcDeployment":   true,
}

// Regex to detect kind = "SomeKind" or "kind": "SomeKind" patterns in Python.
var pyKindRefPatterns = []*regexp.Regexp{
	// kind = "LocalQueue" or kind="LocalQueue"
	regexp.MustCompile(`\bkind\s*=\s*["']([A-Z][a-zA-Z]+)["']`),
	// "kind": "LocalQueue"
	regexp.MustCompile(`["']kind["']\s*:\s*["']([A-Z][a-zA-Z]+)["']`),
}

// Directories to skip when scanning Python files for k8s calls.
// Extends pythonSkipDirs from python_ports.go with additional exclusions.
var pyK8sSkipDirs = map[string]bool{
	"venv":           true,
	".venv":          true,
	"env":            true,
	".env":           true,
	"__pycache__":    true,
	".tox":           true,
	".eggs":          true,
	"node_modules":   true,
	"vendor":         true,
	".git":           true,
	"site-packages":  true,
}

// extractPythonK8sCalls scans Python source files for Kubernetes API interactions
// and returns structured call metadata.
func extractPythonK8sCalls(repoPath string) []PythonK8sCall {
	pyFiles := findPythonFilesForK8s(repoPath)
	if len(pyFiles) == 0 {
		return nil
	}

	// Deduplicate by api+operation+resource (keep first occurrence).
	type callKey struct {
		api       string
		operation string
		resource  string
	}
	seen := make(map[callKey]bool)
	var calls []PythonK8sCall

	for _, fpath := range pyFiles {
		info, err := os.Lstat(fpath)
		if err != nil || info.Mode()&os.ModeSymlink != 0 {
			continue
		}
		if info.Size() > maxFileSize {
			log.Printf("skipping oversized Python file %s: %d bytes", fpath, info.Size())
			continue
		}

		data, err := os.ReadFile(fpath)
		if err != nil {
			log.Printf("warning: skipping Python file %s: %v", fpath, err)
			continue
		}

		lines := strings.Split(string(data), "\n")
		source := relativePath(repoPath, fpath)

		for lineNum, line := range lines {
			stripped := strings.TrimSpace(line)
			// Skip comments and blank lines
			if stripped == "" || strings.HasPrefix(stripped, "#") {
				continue
			}

			loc := source + ":" + strconv.Itoa(lineNum+1)

			// Check import patterns
			for _, pat := range pyK8sImportPatterns {
				if m := pat.re.FindStringSubmatch(line); m != nil {
					importedNames := ""
					if len(m) > 1 {
						importedNames = m[1]
					}
					apis := inferAPIsFromImport(pat.api, importedNames)
					for _, api := range apis {
						key := callKey{api, "import", ""}
						if !seen[key] {
							seen[key] = true
							calls = append(calls, PythonK8sCall{
								API:       api,
								Operation: "import",
								Resource:  "",
								Source:    loc,
							})
						}
					}
				}
			}

			// Check method call patterns.
			// The first pattern (custom object) is more specific, so if it
			// matches we skip the generic second pattern to avoid duplicates.
			customObjectMatched := false
			if m := pyK8sMethodPatterns[0].re.FindStringSubmatch(stripped); m != nil {
				customObjectMatched = true
				operation := m[pyK8sMethodPatterns[0].operation]
				key := callKey{"CustomObjectsApi", operation, "custom_object"}
				if !seen[key] {
					seen[key] = true
					calls = append(calls, PythonK8sCall{
						API:       "CustomObjectsApi",
						Operation: operation,
						Resource:  "custom_object",
						Source:    loc,
					})
				}
			}
			if !customObjectMatched {
				pat := pyK8sMethodPatterns[1]
				if m := pat.re.FindStringSubmatch(stripped); m != nil {
					operation := m[pat.operation]
					resource := m[pat.resource]
					api := inferAPIFromResource(resource)
					key := callKey{api, operation, resource}
					if !seen[key] {
						seen[key] = true
						calls = append(calls, PythonK8sCall{
							API:       api,
							Operation: operation,
							Resource:  resource,
							Source:    loc,
						})
					}
				}
			}

			// Check CRD kind references
			for _, re := range pyKindRefPatterns {
				if m := re.FindStringSubmatch(stripped); m != nil {
					kind := m[1]
					if pyKnownCRDKinds[kind] {
						key := callKey{"CRD", "kind_ref", kind}
						if !seen[key] {
							seen[key] = true
							calls = append(calls, PythonK8sCall{
								API:       "CRD",
								Operation: "kind_ref",
								Resource:  kind,
								Source:    loc,
							})
						}
					}
				}
			}
		}
	}

	return calls
}

// knownK8sAPIClasses is the set of recognized kubernetes client API class names.
var knownK8sAPIClasses = map[string]bool{
	"CustomObjectsApi":         true,
	"CoreV1Api":                true,
	"AppsV1Api":                true,
	"BatchV1Api":               true,
	"NetworkingV1Api":          true,
	"RbacAuthorizationV1Api":   true,
}

// inferAPIsFromImport determines one or more API classes from import context.
// Returns a slice because a single import line can import multiple API classes.
func inferAPIsFromImport(module, importedNames string) []string {
	if module == "openshift.dynamic" {
		return []string{"DynamicClient"}
	}
	if module == "kubernetes.client" {
		var apis []string
		for _, name := range strings.Split(importedNames, ",") {
			name = strings.TrimSpace(name)
			if knownK8sAPIClasses[name] {
				apis = append(apis, name)
			}
		}
		if len(apis) > 0 {
			return apis
		}
		return []string{"kubernetes.client"}
	}
	return []string{module}
}

// inferAPIFromResource determines the likely API class from the resource name
// in a method call.
func inferAPIFromResource(resource string) string {
	switch {
	case resource == "custom_object":
		return "CustomObjectsApi"
	case strings.HasSuffix(resource, "deployment") ||
		strings.HasSuffix(resource, "stateful_set") ||
		strings.HasSuffix(resource, "daemon_set") ||
		strings.HasSuffix(resource, "replica_set"):
		return "AppsV1Api"
	case strings.HasSuffix(resource, "job") && !strings.Contains(resource, "cron"):
		return "BatchV1Api"
	case strings.HasSuffix(resource, "cron_job"):
		return "BatchV1Api"
	case strings.HasSuffix(resource, "ingress") ||
		strings.HasSuffix(resource, "network_policy"):
		return "NetworkingV1Api"
	default:
		// config_map, secret, service, pod, namespace, etc.
		return "CoreV1Api"
	}
}

// findPythonFilesForK8s returns all .py files under repoPath suitable for k8s
// call scanning. Skips venv, __pycache__, test files, site-packages.
func findPythonFilesForK8s(repoPath string) []string {
	var result []string
	walkFiles(repoPath, func(path string, info os.FileInfo) {
		if !strings.HasSuffix(info.Name(), ".py") {
			return
		}
		// Skip test files
		name := info.Name()
		if strings.HasPrefix(name, "test_") || strings.HasSuffix(name, "_test.py") {
			return
		}
		// Skip files inside test/tests directories and site-packages
		rel, _ := strings.CutPrefix(path, repoPath+"/")
		if strings.Contains(rel, "/tests/") || strings.Contains(rel, "/test/") ||
			strings.Contains(rel, "/site-packages/") {
			return
		}
		result = append(result, path)
	})
	return result
}

// walkFiles walks the directory tree, skipping directories in pyK8sSkipDirs.
func walkFiles(root string, fn func(path string, info os.FileInfo)) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return
	}
	for _, e := range entries {
		path := root + "/" + e.Name()
		if e.IsDir() {
			if pyK8sSkipDirs[e.Name()] {
				continue
			}
			walkFiles(path, fn)
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		fn(path, info)
	}
}
