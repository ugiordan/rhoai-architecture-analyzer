package extractor

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var singleRequireRE = regexp.MustCompile(`^require\s+(\S+)\s+(\S+)`)
var singleReplaceRE = regexp.MustCompile(`^replace\s+(\S+)\s+\S*\s*=>\s*(\S+)\s+(\S+)`)


// extractDependencies parses go.mod files and extracts module dependencies,
// identifying internal dependencies based on the provided module prefixes.
func extractDependencies(repoPath string, modulePrefixes []string) *DependencyData {
	goModFiles := findFiles(repoPath, []string{"go.mod", "**/go.mod"})

	var goModules []GoModule
	var internalODH []InternalODH
	var replaceDirectives []ReplaceDirective
	var goVersion, toolchain string

	for _, fpath := range goModFiles {
		data, err := os.ReadFile(fpath)
		if err != nil {
			log.Printf("warning: skipping %s: %v", fpath, err)
			continue
		}

		content := string(data)

		// First pass: collect replace directives
		replaces := parseReplaceDirectives(content)

		// Surface replace directives in output
		for orig, rep := range replaces {
			replaceDirectives = append(replaceDirectives, ReplaceDirective{
				Original:    orig,
				Replacement: rep.module,
				Version:     rep.version,
			})
		}

		// Second pass: collect require directives and metadata
		inRequire := false
		for _, line := range strings.Split(content, "\n") {
			stripped := strings.TrimSpace(line)

			// Skip comments (go.mod uses // for comments)
			if strings.HasPrefix(stripped, "//") {
				continue
			}

			// Extract go version directive (e.g., "go 1.21")
			if goVersion == "" && strings.HasPrefix(stripped, "go ") && !strings.HasPrefix(stripped, "go.") {
				fields := strings.Fields(stripped)
				if len(fields) == 2 {
					goVersion = fields[1]
				}
				continue
			}

			// Extract toolchain directive (e.g., "toolchain go1.21.0")
			if toolchain == "" && strings.HasPrefix(stripped, "toolchain ") {
				fields := strings.Fields(stripped)
				if len(fields) == 2 {
					toolchain = fields[1]
				}
				continue
			}

			if strings.HasPrefix(stripped, "require (") || strings.HasPrefix(stripped, "require(") {
				inRequire = true
				continue
			}
			if inRequire && stripped == ")" {
				inRequire = false
				continue
			}

			// Single-line require
			if match := singleRequireRE.FindStringSubmatch(stripped); match != nil {
				if strings.HasSuffix(stripped, "// indirect") {
					continue
				}
				module, version := match[1], match[2]
				// Apply replace directive if present
				if rep, ok := replaces[module]; ok {
					module = rep.module
					version = rep.version
				}
				goModules = append(goModules, GoModule{Module: module, Version: version})
				if comp, ok := matchInternalModule(module, modulePrefixes); ok {
					internalODH = append(internalODH, InternalODH{
						Component:   comp,
						Interaction: fmt.Sprintf("Go module dependency: %s", module),
					})
				}
				continue
			}

			// Lines inside require block
			if inRequire {
				if strings.HasSuffix(stripped, "// indirect") {
					continue
				}
				parts := strings.Fields(stripped)
				if len(parts) >= 2 {
					module, version := parts[0], parts[1]
					if rep, ok := replaces[module]; ok {
						module = rep.module
						version = rep.version
					}
					goModules = append(goModules, GoModule{Module: module, Version: version})
					if comp, ok := matchInternalModule(module, modulePrefixes); ok {
						internalODH = append(internalODH, InternalODH{
							Component:   comp,
							Interaction: fmt.Sprintf("Go module dependency: %s", module),
						})
					}
				}
			}
		}
	}

	if goModules == nil {
		goModules = []GoModule{}
	}
	if internalODH == nil {
		internalODH = []InternalODH{}
	}

	// Categorize modules and detect compatibility issues
	categorizeModules(goModules)
	issues := detectDependencyIssues(goModules)

	return &DependencyData{
		GoVersion:         goVersion,
		Toolchain:         toolchain,
		GoModules:         goModules,
		ReplaceDirectives: replaceDirectives,
		InternalODH:       internalODH,
		Issues:            issues,
	}
}

type replaceTarget struct {
	module  string
	version string
}

// parseReplaceDirectives extracts replace directives from go.mod content.
// Handles both single-line (replace A => B v1) and block (replace (...)) forms.
func parseReplaceDirectives(content string) map[string]replaceTarget {
	replaces := make(map[string]replaceTarget)
	inReplace := false

	for _, line := range strings.Split(content, "\n") {
		stripped := strings.TrimSpace(line)

		if strings.HasPrefix(stripped, "replace (") || strings.HasPrefix(stripped, "replace(") {
			inReplace = true
			continue
		}
		if inReplace && stripped == ")" {
			inReplace = false
			continue
		}

		// Single-line replace
		if match := singleReplaceRE.FindStringSubmatch(stripped); match != nil {
			original, replacement, version := match[1], match[2], match[3]
			// Skip local path replacements (no version, starts with . or /)
			if !strings.HasPrefix(replacement, ".") && !strings.HasPrefix(replacement, "/") {
				replaces[original] = replaceTarget{module: replacement, version: version}
			}
			continue
		}

		// Lines inside replace block: A v1 => B v2
		if inReplace {
			parts := strings.Fields(stripped)
			// Format: module [version] => replacement version
			arrowIdx := -1
			for i, p := range parts {
				if p == "=>" {
					arrowIdx = i
					break
				}
			}
			if arrowIdx >= 1 && arrowIdx+2 <= len(parts) {
				original := parts[0]
				replacement := parts[arrowIdx+1]
				version := ""
				if arrowIdx+2 < len(parts) {
					version = parts[arrowIdx+2]
				}
				if !strings.HasPrefix(replacement, ".") && !strings.HasPrefix(replacement, "/") {
					replaces[original] = replaceTarget{module: replacement, version: version}
				}
			}
		}
	}

	return replaces
}

// matchInternalModule checks if a module matches any of the internal prefixes.
// Returns the component name and true if matched.
func matchInternalModule(module string, prefixes []string) (string, bool) {
	for _, prefix := range prefixes {
		if strings.HasPrefix(module, prefix) {
			component := strings.SplitN(module[len(prefix):], "/", 2)[0]
			return component, true
		}
	}
	return "", false
}

// moduleCategoryRules maps module path prefixes to categories and purposes.
// Order matters: first match wins.
var moduleCategoryRules = []struct {
	prefix   string
	category string
	purpose  string
}{
	// Controller frameworks
	{"sigs.k8s.io/controller-runtime", "controller-framework", "Kubernetes controller framework"},
	{"sigs.k8s.io/controller-tools", "controller-framework", "Controller code generation tools"},

	// Kubernetes clients and API
	{"k8s.io/client-go", "k8s-client", "Kubernetes API client"},
	{"k8s.io/api", "k8s-api", "Kubernetes API types"},
	{"k8s.io/apimachinery", "k8s-api", "Kubernetes API machinery (runtime, schema, types)"},
	{"k8s.io/apiextensions-apiserver", "k8s-api", "CRD API types"},
	{"k8s.io/kube-openapi", "k8s-api", "Kubernetes OpenAPI schema generation"},
	{"k8s.io/utils", "k8s-api", "Kubernetes utility library"},
	{"k8s.io/apiserver", "k8s-api", "Kubernetes API server library"},
	{"k8s.io/kubectl", "k8s-client", "Kubernetes CLI library"},
	{"k8s.io/klog", "k8s-api", "Kubernetes logging library"},
	{"k8s.io/component-base", "k8s-api", "Kubernetes component base library"},

	// OpenShift
	{"github.com/openshift/api", "openshift", "OpenShift API types"},
	{"github.com/openshift/client-go", "openshift", "OpenShift API client"},
	{"github.com/openshift/custom-resource-status", "openshift", "OpenShift CR status helpers"},
	{"github.com/openshift/library-go", "openshift", "OpenShift library utilities"},

	// Cloud providers
	{"github.com/aws/aws-sdk-go", "cloud-provider", "AWS SDK"},
	{"github.com/Azure/azure-sdk-for-go", "cloud-provider", "Azure SDK"},
	{"cloud.google.com/go", "cloud-provider", "Google Cloud SDK"},
	{"google.golang.org/api", "cloud-provider", "Google API client"},

	// Storage and databases
	{"github.com/go-sql-driver/mysql", "storage", "MySQL database driver"},
	{"github.com/lib/pq", "storage", "PostgreSQL database driver"},
	{"github.com/jackc/pgx", "storage", "PostgreSQL database driver"},
	{"github.com/minio/minio-go", "storage", "MinIO/S3-compatible object storage client"},
	{"go.etcd.io/etcd", "storage", "etcd client"},
	{"github.com/redis/go-redis", "storage", "Redis client"},

	// Networking and service mesh
	{"istio.io/", "networking", "Istio service mesh"},
	{"knative.dev/", "networking", "Knative serverless framework"},
	{"sigs.k8s.io/gateway-api", "networking", "Kubernetes Gateway API"},

	// Observability
	{"github.com/prometheus/", "observability", "Prometheus metrics"},
	{"go.opentelemetry.io/", "observability", "OpenTelemetry tracing/metrics"},

	// Operator tooling
	{"github.com/operator-framework/", "operator-tooling", "Operator Framework (OLM, SDK)"},
	{"sigs.k8s.io/kustomize", "operator-tooling", "Kustomize configuration management"},
	{"github.com/manifestival/manifestival", "operator-tooling", "Manifest manipulation library"},

	// Container and image
	{"github.com/google/go-containerregistry", "container", "Container image registry client"},
	{"github.com/containers/image", "container", "Container image tools"},
	{"oras.land/oras-go", "container", "OCI registry artifact storage"},

	// Auth
	{"github.com/coreos/go-oidc", "auth", "OIDC authentication"},
	{"golang.org/x/oauth2", "auth", "OAuth2 client"},

	// Serialization
	{"google.golang.org/protobuf", "serialization", "Protocol Buffers"},
	{"google.golang.org/grpc", "serialization", "gRPC framework"},
	{"gopkg.in/yaml.v", "serialization", "YAML parser"},
	{"sigs.k8s.io/yaml", "serialization", "Kubernetes YAML parser"},
	{"github.com/ghodss/yaml", "serialization", "JSON/YAML converter"},

	// Cert management
	{"github.com/cert-manager/", "cert-management", "Certificate management"},

	// Testing
	{"github.com/onsi/ginkgo", "testing", "BDD testing framework"},
	{"github.com/onsi/gomega", "testing", "Test matcher library"},
	{"github.com/stretchr/testify", "testing", "Test assertions library"},

	// ML/AI specific
	{"github.com/kubeflow/", "ml-platform", "Kubeflow ML platform"},
	{"github.com/kserve/", "ml-platform", "KServe model serving"},
}

// categorizeModules assigns category and purpose to each module.
func categorizeModules(modules []GoModule) {
	for i := range modules {
		for _, rule := range moduleCategoryRules {
			if strings.HasPrefix(modules[i].Module, rule.prefix) {
				modules[i].Category = rule.category
				modules[i].Purpose = rule.purpose
				break
			}
		}
	}
}

// detectDependencyIssues checks for version compatibility issues between
// related Kubernetes modules (controller-runtime, client-go, k8s.io/api).
func detectDependencyIssues(modules []GoModule) []string {
	versions := make(map[string]string) // module -> version
	for _, m := range modules {
		versions[m.Module] = m.Version
	}

	var issues []string

	// Extract k8s minor version from module versions.
	// controller-runtime v0.X maps to k8s 1.(X+15) approximately:
	//   CR v0.15 -> k8s 1.27, CR v0.16 -> k8s 1.28, CR v0.17 -> k8s 1.29, etc.
	// client-go v0.X.Y maps directly to k8s 1.X.
	clientGoVer := ""
	crVer := ""
	for mod, ver := range versions {
		if mod == "k8s.io/client-go" {
			clientGoVer = ver
		}
		if mod == "sigs.k8s.io/controller-runtime" {
			crVer = ver
		}
	}

	if clientGoVer != "" && crVer != "" {
		cgMinor := extractMinorVersion(clientGoVer)
		crMinor := extractMinorVersion(crVer)

		// controller-runtime v0.X pairs with client-go v0.(X+12) approximately
		// e.g., CR v0.15 -> client-go v0.27 (k8s 1.27), CR v0.18 -> client-go v0.30
		// The offset varies slightly across releases, so use tolerance of 4.
		if cgMinor > 0 && crMinor > 0 {
			expectedCGMinor := crMinor + 12
			if abs(cgMinor-expectedCGMinor) > 4 {
				issues = append(issues, fmt.Sprintf(
					"controller-runtime %s and client-go %s may be from different Kubernetes release cycles",
					crVer, clientGoVer))
			}
		}
	}

	// Check k8s.io/* modules are from the same minor version.
	// Exclude modules with independent versioning (klog, utils, component-base).
	k8sIndependentVersioning := map[string]bool{
		"k8s.io/klog":           true,
		"k8s.io/klog/v2":        true,
		"k8s.io/utils":          true,
		"k8s.io/kube-openapi":   true,
		"k8s.io/component-base": true,
	}
	k8sMinors := make(map[int][]string) // minor -> list of modules
	for _, m := range modules {
		if !strings.HasPrefix(m.Module, "k8s.io/") {
			continue
		}
		if k8sIndependentVersioning[m.Module] {
			continue
		}
		minor := extractMinorVersion(m.Version)
		if minor > 0 {
			k8sMinors[minor] = append(k8sMinors[minor], m.Module)
		}
	}
	if len(k8sMinors) > 1 {
		var parts []string
		for minor, mods := range k8sMinors {
			parts = append(parts, fmt.Sprintf("v0.%d (%s)", minor, strings.Join(mods, ", ")))
		}
		issues = append(issues, "k8s.io modules span multiple minor versions: "+strings.Join(parts, "; "))
	}

	return issues
}

// extractMinorVersion extracts the minor version number from a semver string.
// e.g., "v0.27.2" -> 27, "v0.15.0" -> 15.
func extractMinorVersion(ver string) int {
	ver = strings.TrimPrefix(ver, "v")
	parts := strings.Split(ver, ".")
	if len(parts) >= 2 {
		if n, err := fmt.Sscanf(parts[1], "%d", new(int)); n == 1 && err == nil {
			v := 0
			fmt.Sscanf(parts[1], "%d", &v)
			return v
		}
	}
	return 0
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
