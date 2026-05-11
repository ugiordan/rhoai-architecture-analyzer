package extractor

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// serviceClassification maps an env var name or prefix to a runtime service.
type serviceClassification struct {
	Name string // human-readable service name
	Type string // database, object-storage, cache, service, messaging, observability
}

// envVarClassification maps well-known environment variable names (uppercased)
// to the runtime service they imply.
var envVarClassification = map[string]serviceClassification{
	// Databases
	"DATABASE_HOST":     {"PostgreSQL", "database"},
	"DATABASE_URL":      {"PostgreSQL", "database"},
	"DATABASE_PORT":     {"PostgreSQL", "database"},
	"DB_HOST":           {"PostgreSQL", "database"},
	"DB_URL":            {"PostgreSQL", "database"},
	"POSTGRES_HOST":     {"PostgreSQL", "database"},
	"POSTGRES_URL":      {"PostgreSQL", "database"},
	"POSTGRES_PORT":     {"PostgreSQL", "database"},
	"PGHOST":            {"PostgreSQL", "database"},
	"PGPORT":            {"PostgreSQL", "database"},
	"PGDATABASE":        {"PostgreSQL", "database"},
	"MYSQL_HOST":        {"MySQL", "database"},
	"MYSQL_PORT":        {"MySQL", "database"},
	"MYSQL_DATABASE":    {"MySQL", "database"},
	"MYSQL_URL":         {"MySQL", "database"},
	"MONGO_HOST":        {"MongoDB", "database"},
	"MONGO_URL":         {"MongoDB", "database"},
	"MONGODB_URI":       {"MongoDB", "database"},
	"MONGODB_HOST":      {"MongoDB", "database"},

	// Cache
	"REDIS_URL":  {"Redis", "cache"},
	"REDIS_HOST": {"Redis", "cache"},
	"REDIS_PORT": {"Redis", "cache"},
	"REDIS_ADDR": {"Redis", "cache"},

	// Object storage
	"S3_ENDPOINT":          {"S3", "object-storage"},
	"S3_BUCKET":            {"S3", "object-storage"},
	"AWS_S3_ENDPOINT":      {"S3", "object-storage"},
	"AWS_S3_BUCKET":        {"S3", "object-storage"},
	"MINIO_ENDPOINT":       {"MinIO/S3", "object-storage"},
	"MINIO_HOST":           {"MinIO/S3", "object-storage"},
	"OBJECT_STORE_ENDPOINT": {"S3", "object-storage"},

	// ML/AI services
	"MLFLOW_TRACKING_URI":      {"MLflow", "service"},
	"MLFLOW_TRACKING_URL":      {"MLflow", "service"},
	"MLFLOW_S3_ENDPOINT_URL":   {"MLflow", "service"},

	// Observability
	"OTEL_EXPORTER_OTLP_ENDPOINT":       {"OpenTelemetry Collector", "observability"},
	"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT": {"OpenTelemetry Collector", "observability"},
	"JAEGER_ENDPOINT":                    {"Jaeger", "observability"},
	"JAEGER_AGENT_HOST":                  {"Jaeger", "observability"},
	"PROMETHEUS_PUSHGATEWAY_URL":         {"Prometheus Pushgateway", "observability"},

	// Messaging
	"KAFKA_BOOTSTRAP_SERVERS": {"Kafka", "messaging"},
	"KAFKA_BROKERS":           {"Kafka", "messaging"},
	"KAFKA_BROKER_URL":        {"Kafka", "messaging"},
	"NATS_URL":                {"NATS", "messaging"},
	"NATS_HOST":               {"NATS", "messaging"},
	"RABBITMQ_HOST":           {"RabbitMQ", "messaging"},
	"AMQP_URL":                {"RabbitMQ", "messaging"},

	// Auth
	"OAUTH2_PROXY_UPSTREAM":     {"OAuth2 Proxy", "service"},
	"KEYCLOAK_URL":              {"Keycloak", "service"},
	"OIDC_ISSUER":               {"OIDC Provider", "service"},

	// General services
	"ETCD_ENDPOINTS":     {"etcd", "database"},
	"VAULT_ADDR":         {"HashiCorp Vault", "service"},
	"CONSUL_HTTP_ADDR":   {"HashiCorp Consul", "service"},
}

// envVarPrefixClassification matches env var prefixes (checked with HasPrefix).
// Checked after exact match fails.
var envVarPrefixClassification = []struct {
	prefix string
	class  serviceClassification
}{
	{"POSTGRES_", serviceClassification{"PostgreSQL", "database"}},
	{"MYSQL_", serviceClassification{"MySQL", "database"}},
	{"REDIS_", serviceClassification{"Redis", "cache"}},
	{"KAFKA_", serviceClassification{"Kafka", "messaging"}},
	{"NATS_", serviceClassification{"NATS", "messaging"}},
	{"MINIO_", serviceClassification{"MinIO/S3", "object-storage"}},
	{"MLFLOW_", serviceClassification{"MLflow", "service"}},
	{"OTEL_", serviceClassification{"OpenTelemetry Collector", "observability"}},
	{"JAEGER_", serviceClassification{"Jaeger", "observability"}},
	{"MONGO_", serviceClassification{"MongoDB", "database"}},
	{"MONGODB_", serviceClassification{"MongoDB", "database"}},
}

// envVarContainsClassification matches substrings anywhere in the env var name.
// Checked after exact and prefix matches fail. Catches patterns like
// TEST_DATA_S3_BUCKET or MY_APP_REDIS_HOST.
var envVarContainsClassification = []struct {
	substring string
	class     serviceClassification
}{
	{"_S3_", serviceClassification{"S3", "object-storage"}},
	{"_POSTGRES", serviceClassification{"PostgreSQL", "database"}},
	{"_MYSQL", serviceClassification{"MySQL", "database"}},
	{"_REDIS", serviceClassification{"Redis", "cache"}},
	{"_KAFKA", serviceClassification{"Kafka", "messaging"}},
	{"_MINIO", serviceClassification{"MinIO/S3", "object-storage"}},
	{"_MLFLOW", serviceClassification{"MLflow", "service"}},
	{"_MONGO", serviceClassification{"MongoDB", "database"}},
	{"_OTEL_", serviceClassification{"OpenTelemetry Collector", "observability"}},
}

// svcDNSPattern matches Kubernetes service DNS references:
// <svc>.<ns>.svc.cluster.local or <svc>.<ns>.svc
var svcDNSPattern = regexp.MustCompile(`([a-z0-9][-a-z0-9]*)\.[a-z0-9][-a-z0-9]*\.svc(?:\.cluster\.local)?`)

// connectionURISchemes maps URI scheme prefixes in env var values to services.
var connectionURISchemes = map[string]serviceClassification{
	"postgresql://": {"PostgreSQL", "database"},
	"postgres://":   {"PostgreSQL", "database"},
	"mysql://":      {"MySQL", "database"},
	"mongodb://":    {"MongoDB", "database"},
	"mongodb+srv://": {"MongoDB", "database"},
	"redis://":      {"Redis", "cache"},
	"rediss://":     {"Redis", "cache"},
	"amqp://":       {"RabbitMQ", "messaging"},
	"amqps://":      {"RabbitMQ", "messaging"},
	"kafka://":      {"Kafka", "messaging"},
	"nats://":       {"NATS", "messaging"},
}

// osGetenvPattern matches os.Getenv("VAR_NAME") calls in Go source (literal strings).
var osGetenvPattern = regexp.MustCompile(`os\.Getenv\(\s*"([A-Z][A-Z0-9_]+)"\s*\)`)

// osGetenvIndirectPattern matches os.Getenv(varName) with non-literal identifiers.
var osGetenvIndirectPattern = regexp.MustCompile(`os\.Getenv\(\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*\)`)

// goConstEnvPattern matches const declarations like: envBucket = "TEST_DATA_S3_BUCKET"
// Handles both typed and untyped constants, with = or := assignment.
var goConstEnvPattern = regexp.MustCompile(`(?:const\s+)?([a-zA-Z_][a-zA-Z0-9_]*)\s*(?::?=)\s*"([A-Z][A-Z0-9_]+)"`)

// resolveGoConstants does a two-pass scan of Go source files:
// 1. Collect all constant/var declarations that assign uppercase string literals
// 2. For each os.Getenv(constName), resolve constName to the actual env var name
func resolveGoConstants(repoPath string, goFiles []string) map[string]string {
	constants := make(map[string]string) // constName -> "ENV_VAR_VALUE"

	for _, fpath := range goFiles {
		if strings.Contains(fpath, "_test.go") || strings.Contains(fpath, "/vendor/") {
			continue
		}
		info, err := os.Lstat(fpath)
		if err != nil || info.Mode()&os.ModeSymlink != 0 || info.Size() > maxFileSize {
			continue
		}
		f, err := os.Open(fpath)
		if err != nil {
			continue
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			stripped := strings.TrimSpace(line)
			if strings.HasPrefix(stripped, "//") {
				continue
			}
			matches := goConstEnvPattern.FindAllStringSubmatch(line, -1)
			for _, m := range matches {
				if len(m) >= 3 {
					constants[m[1]] = m[2]
				}
			}
		}
		f.Close()
	}
	return constants
}

// pythonOsEnvironPattern matches os.environ["VAR"], os.environ.get("VAR"),
// and os.environ.get("VAR", default) in Python.
var pythonOsEnvironPattern = regexp.MustCompile(`os\.environ(?:\.get\s*\(\s*["']([A-Z][A-Z0-9_]+)["']|\[["']([A-Z][A-Z0-9_]+)["']\])`)

// dockerfileEnvPattern matches ENV directives in Dockerfiles.
var dockerfileEnvPattern = regexp.MustCompile(`^ENV\s+([A-Z][A-Z0-9_]+)\s*[=\s]`)

// extractRuntimeDependencies scans deployment manifests, Go/Python source,
// and Dockerfiles for environment variables and patterns that indicate
// runtime service dependencies.
func extractRuntimeDependencies(repoPath string) []RuntimeDependency {
	seen := make(map[string]RuntimeDependency) // keyed by Name to deduplicate

	// 1. Scan YAML deployment manifests for env vars and configMapRef/secretRef
	scanYAMLManifests(repoPath, seen)

	// 2. Scan Go source for os.Getenv patterns
	scanGoSourceForEnvVars(repoPath, seen)

	// 3. Scan Python source for os.environ patterns
	scanPythonSourceForEnvVars(repoPath, seen)

	// 4. Scan Dockerfiles for ENV directives
	scanDockerfilesForEnvVars(repoPath, seen)

	// 5. Scan all source for Kubernetes service DNS patterns
	scanForServiceDNS(repoPath, seen)

	// 6. Scan all source for connection URI schemes in string literals
	scanForConnectionURIs(repoPath, seen)

	result := make([]RuntimeDependency, 0, len(seen))
	for _, dep := range seen {
		result = append(result, dep)
	}
	if result == nil {
		result = []RuntimeDependency{}
	}
	return result
}

// recordDependency adds a runtime dependency, deduplicating by service name.
// If the same service was already detected, it keeps the first occurrence
// (which is typically the most authoritative source like a deployment manifest).
func recordDependency(seen map[string]RuntimeDependency, dep RuntimeDependency) {
	if existing, ok := seen[dep.Name]; ok {
		// If existing is not required but this one is, upgrade
		if dep.Required && !existing.Required {
			existing.Required = true
			seen[dep.Name] = existing
		}
		return
	}
	seen[dep.Name] = dep
}

// classifyEnvVar checks if an env var name maps to a known runtime service.
func classifyEnvVar(envName string) (serviceClassification, bool) {
	upper := strings.ToUpper(envName)

	// Exact match first
	if cls, ok := envVarClassification[upper]; ok {
		return cls, true
	}

	// Prefix match
	for _, pc := range envVarPrefixClassification {
		if strings.HasPrefix(upper, pc.prefix) {
			return pc.class, true
		}
	}

	// Contains match (catches TEST_DATA_S3_BUCKET, MY_APP_REDIS_HOST, etc.)
	for _, cc := range envVarContainsClassification {
		if strings.Contains(upper, cc.substring) {
			return cc.class, true
		}
	}

	return serviceClassification{}, false
}

// classifyEnvValue checks if an env var value contains a connection URI scheme.
func classifyEnvValue(value string) (serviceClassification, bool) {
	lower := strings.ToLower(value)
	for scheme, cls := range connectionURISchemes {
		if strings.HasPrefix(lower, scheme) {
			return cls, true
		}
	}
	return serviceClassification{}, false
}

// scanYAMLManifests scans YAML files for deployment/statefulset env vars
// and configMapRef/secretRef patterns.
func scanYAMLManifests(repoPath string, seen map[string]RuntimeDependency) {
	yamlFiles := findYAMLFiles(repoPath, []string{
		"**/*.yaml",
		"**/*.yml",
	})

	for _, fpath := range yamlFiles {
		if strings.Contains(fpath, "/vendor/") || strings.Contains(fpath, "/testdata/") {
			continue
		}

		docs := parseYAMLSafe(fpath)
		source := relativePath(repoPath, fpath)

		for _, doc := range docs {
			kind, _ := doc["kind"].(string)
			if kind != "Deployment" && kind != "StatefulSet" && kind != "DaemonSet" && kind != "Job" && kind != "CronJob" {
				continue
			}

			// Walk the spec to find containers
			spec, _ := doc["spec"].(map[string]interface{})
			if spec == nil {
				continue
			}

			// Handle CronJob nesting: spec.jobTemplate.spec.template.spec
			if kind == "CronJob" {
				if jt, ok := spec["jobTemplate"].(map[string]interface{}); ok {
					if jtSpec, ok := jt["spec"].(map[string]interface{}); ok {
						spec = jtSpec
					}
				}
			}

			template, _ := spec["template"].(map[string]interface{})
			if template == nil {
				continue
			}
			podSpec, _ := template["spec"].(map[string]interface{})
			if podSpec == nil {
				continue
			}

			// Scan containers and initContainers
			for _, containerKey := range []string{"containers", "initContainers"} {
				containers, _ := podSpec[containerKey].([]interface{})
				for _, c := range containers {
					cm, ok := c.(map[string]interface{})
					if !ok {
						continue
					}
					scanContainerEnvVars(cm, source, seen)
					scanContainerEnvFrom(cm, source, seen)
				}
			}
		}
	}
}

// scanContainerEnvVars extracts env var names/values from a container map.
func scanContainerEnvVars(container map[string]interface{}, source string, seen map[string]RuntimeDependency) {
	envList, _ := container["env"].([]interface{})
	for _, e := range envList {
		envMap, ok := e.(map[string]interface{})
		if !ok {
			continue
		}
		envName, _ := envMap["name"].(string)
		if envName == "" {
			continue
		}

		// Check the env var name
		if cls, ok := classifyEnvVar(envName); ok {
			// Check if there's a default value (which might make it optional)
			value, _ := envMap["value"].(string)
			required := value == "" // no inline default means likely required
			recordDependency(seen, RuntimeDependency{
				Name:     cls.Name,
				Type:     cls.Type,
				Source:   source,
				Evidence: "env:" + envName,
				Required: required,
			})
			continue
		}

		// Check the env var value for connection URI schemes
		value, _ := envMap["value"].(string)
		if value != "" {
			if cls, ok := classifyEnvValue(value); ok {
				recordDependency(seen, RuntimeDependency{
					Name:     cls.Name,
					Type:     cls.Type,
					Source:   source,
					Evidence: "env:" + envName + "=" + redactConnectionString(value),
					Required: true,
				})
			}
		}
	}
}

// scanContainerEnvFrom detects configMapRef and secretRef references that
// suggest external service configuration.
func scanContainerEnvFrom(container map[string]interface{}, source string, seen map[string]RuntimeDependency) {
	envFrom, _ := container["envFrom"].([]interface{})
	for _, ef := range envFrom {
		efMap, ok := ef.(map[string]interface{})
		if !ok {
			continue
		}

		// configMapRef
		if cmRef, ok := efMap["configMapRef"].(map[string]interface{}); ok {
			name, _ := cmRef["name"].(string)
			classifyRefName(name, source, "configMapRef:"+name, seen)
		}

		// secretRef
		if secRef, ok := efMap["secretRef"].(map[string]interface{}); ok {
			name, _ := secRef["name"].(string)
			classifyRefName(name, source, "secretRef:"+name, seen)
		}
	}
}

// refNameKeywords maps substrings in configMapRef/secretRef names to services.
var refNameKeywords = []struct {
	keyword string
	class   serviceClassification
}{
	{"postgres", serviceClassification{"PostgreSQL", "database"}},
	{"mysql", serviceClassification{"MySQL", "database"}},
	{"mongo", serviceClassification{"MongoDB", "database"}},
	{"redis", serviceClassification{"Redis", "cache"}},
	{"s3", serviceClassification{"S3", "object-storage"}},
	{"minio", serviceClassification{"MinIO/S3", "object-storage"}},
	{"kafka", serviceClassification{"Kafka", "messaging"}},
	{"nats", serviceClassification{"NATS", "messaging"}},
	{"rabbitmq", serviceClassification{"RabbitMQ", "messaging"}},
	{"amqp", serviceClassification{"RabbitMQ", "messaging"}},
	{"mlflow", serviceClassification{"MLflow", "service"}},
	{"otel", serviceClassification{"OpenTelemetry Collector", "observability"}},
	{"jaeger", serviceClassification{"Jaeger", "observability"}},
	{"vault", serviceClassification{"HashiCorp Vault", "service"}},
	{"keycloak", serviceClassification{"Keycloak", "service"}},
	{"etcd", serviceClassification{"etcd", "database"}},
}

// classifyRefName checks if a configMapRef/secretRef name implies a service dependency.
func classifyRefName(name, source, evidence string, seen map[string]RuntimeDependency) {
	if name == "" {
		return
	}
	lower := strings.ToLower(name)
	for _, kw := range refNameKeywords {
		if strings.Contains(lower, kw.keyword) {
			recordDependency(seen, RuntimeDependency{
				Name:     kw.class.Name,
				Type:     kw.class.Type,
				Source:   source,
				Evidence: evidence,
				Required: true,
			})
			return
		}
	}
}

// scanGoSourceForEnvVars scans Go files for os.Getenv calls with known env vars.
// Uses a two-pass approach: first resolves constants, then matches both literal
// and indirect os.Getenv calls.
func scanGoSourceForEnvVars(repoPath string, seen map[string]RuntimeDependency) {
	goFiles := findFiles(repoPath, []string{"**/*.go"})

	// Pass 1: resolve literal os.Getenv("VAR") calls
	filteredFiles := make([]string, 0, len(goFiles))
	for _, fpath := range goFiles {
		if strings.Contains(fpath, "_test.go") || strings.Contains(fpath, "/vendor/") {
			continue
		}
		filteredFiles = append(filteredFiles, fpath)
		scanSourceFileForPatterns(repoPath, fpath, osGetenvPattern, seen)
	}

	// Pass 2: resolve indirect os.Getenv(constName) via constant lookup
	constants := resolveGoConstants(repoPath, filteredFiles)
	if len(constants) == 0 {
		return
	}

	for _, fpath := range filteredFiles {
		info, err := os.Lstat(fpath)
		if err != nil || info.Mode()&os.ModeSymlink != 0 || info.Size() > maxFileSize {
			continue
		}
		f, err := os.Open(fpath)
		if err != nil {
			continue
		}

		source := relativePath(repoPath, fpath)
		scanner := bufio.NewScanner(f)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()
			stripped := strings.TrimSpace(line)
			if strings.HasPrefix(stripped, "//") {
				continue
			}

			matches := osGetenvIndirectPattern.FindAllStringSubmatch(line, -1)
			for _, m := range matches {
				if len(m) < 2 {
					continue
				}
				constName := m[1]
				// Skip if it's a literal string match (already handled in pass 1)
				if strings.HasPrefix(constName, "\"") {
					continue
				}
				envName, ok := constants[constName]
				if !ok {
					continue
				}
				cls, matched := classifyEnvVar(envName)
				if !matched {
					continue
				}
				recordDependency(seen, RuntimeDependency{
					Name:     cls.Name,
					Type:     cls.Type,
					Source:   source + ":" + strconv.Itoa(lineNum),
					Evidence: "go:os.Getenv(" + constName + "→" + envName + ")",
					Required: true,
				})
			}
		}
		f.Close()
	}
}

// scanPythonSourceForEnvVars scans Python files for os.environ patterns.
func scanPythonSourceForEnvVars(repoPath string, seen map[string]RuntimeDependency) {
	pyFiles := findFiles(repoPath, []string{"**/*.py"})
	for _, fpath := range pyFiles {
		if strings.Contains(fpath, "/vendor/") || strings.Contains(fpath, "/venv/") ||
			strings.Contains(fpath, "/.venv/") || strings.Contains(fpath, "/site-packages/") {
			continue
		}
		scanSourceFileForPatterns(repoPath, fpath, pythonOsEnvironPattern, seen)
	}
}

// scanDockerfilesForEnvVars scans Dockerfiles for ENV directives.
func scanDockerfilesForEnvVars(repoPath string, seen map[string]RuntimeDependency) {
	dockerFiles := findFiles(repoPath, []string{"**/Dockerfile", "**/Dockerfile.*", "**/Containerfile", "**/Containerfile.*"})
	for _, fpath := range dockerFiles {
		scanSourceFileForPatterns(repoPath, fpath, dockerfileEnvPattern, seen)
	}
}

// scanSourceFileForPatterns reads a file and applies a regex to extract env var names.
func scanSourceFileForPatterns(repoPath, fpath string, pattern *regexp.Regexp, seen map[string]RuntimeDependency) {
	info, err := os.Lstat(fpath)
	if err != nil || info.Mode()&os.ModeSymlink != 0 {
		return
	}
	if info.Size() > maxFileSize {
		return
	}

	f, err := os.Open(fpath)
	if err != nil {
		log.Printf("warning: skipping %s: %v", fpath, err)
		return
	}
	defer f.Close()

	source := relativePath(repoPath, fpath)
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		stripped := strings.TrimSpace(line)
		if strings.HasPrefix(stripped, "//") || strings.HasPrefix(stripped, "#") {
			continue
		}

		matches := pattern.FindAllStringSubmatch(line, -1)
		for _, m := range matches {
			if len(m) < 2 {
				continue
			}
			// Pick the first non-empty capture group (handles regexes
			// with multiple alternation groups like the Python pattern).
			envName := ""
			for _, g := range m[1:] {
				if g != "" {
					envName = g
					break
				}
			}
			if envName == "" {
				continue
			}
			if cls, ok := classifyEnvVar(envName); ok {
				recordDependency(seen, RuntimeDependency{
					Name:     cls.Name,
					Type:     cls.Type,
					Source:   source + ":" + strconv.Itoa(lineNum),
					Evidence: "env:" + envName,
					Required: true, // source code lookups typically mean it's needed
				})
			}
		}
	}
}

// scanForServiceDNS scans Go and Python source for Kubernetes service DNS patterns.
func scanForServiceDNS(repoPath string, seen map[string]RuntimeDependency) {
	sourceFiles := findFiles(repoPath, []string{"**/*.go", "**/*.py"})
	for _, fpath := range sourceFiles {
		if strings.Contains(fpath, "_test.go") || strings.Contains(fpath, "/vendor/") ||
			strings.Contains(fpath, "/venv/") || strings.Contains(fpath, "/.venv/") ||
			strings.Contains(fpath, "/site-packages/") {
			continue
		}

		info, err := os.Lstat(fpath)
		if err != nil || info.Mode()&os.ModeSymlink != 0 {
			continue
		}
		if info.Size() > maxFileSize {
			continue
		}

		data, err := os.ReadFile(fpath)
		if err != nil {
			continue
		}

		source := relativePath(repoPath, fpath)
		lines := strings.Split(string(data), "\n")
		for lineNum, line := range lines {
			stripped := strings.TrimSpace(line)
			if strings.HasPrefix(stripped, "//") || strings.HasPrefix(stripped, "#") {
				continue
			}

			matches := svcDNSPattern.FindAllStringSubmatch(line, -1)
			for _, m := range matches {
				if len(m) < 2 {
					continue
				}
				svcName := m[0] // full match like "mlflow-server.namespace.svc.cluster.local"
				recordDependency(seen, RuntimeDependency{
					Name:     m[1], // just the service name portion
					Type:     "service",
					Source:   source + ":" + strconv.Itoa(lineNum+1),
					Evidence: "dns:" + svcName,
					Required: true,
				})
			}
		}
	}
}

// scanForConnectionURIs scans source files for connection URI schemes in string literals.
func scanForConnectionURIs(repoPath string, seen map[string]RuntimeDependency) {
	sourceFiles := findFiles(repoPath, []string{"**/*.go", "**/*.py", "**/*.yaml", "**/*.yml"})
	for _, fpath := range sourceFiles {
		if strings.Contains(fpath, "_test.go") || strings.Contains(fpath, "/vendor/") ||
			strings.Contains(fpath, "/venv/") || strings.Contains(fpath, "/.venv/") ||
			strings.Contains(fpath, "/site-packages/") || strings.Contains(fpath, "/testdata/") {
			continue
		}

		info, err := os.Lstat(fpath)
		if err != nil || info.Mode()&os.ModeSymlink != 0 {
			continue
		}
		if info.Size() > maxFileSize {
			continue
		}

		data, err := os.ReadFile(fpath)
		if err != nil {
			continue
		}

		source := relativePath(repoPath, fpath)
		content := string(data)
		lines := strings.Split(content, "\n")
		for lineNum, line := range lines {
			stripped := strings.TrimSpace(line)
			if strings.HasPrefix(stripped, "//") || strings.HasPrefix(stripped, "#") {
				continue
			}

			lower := strings.ToLower(line)
			for scheme, cls := range connectionURISchemes {
				if strings.Contains(lower, scheme) {
					recordDependency(seen, RuntimeDependency{
						Name:     cls.Name,
						Type:     cls.Type,
						Source:   source + ":" + strconv.Itoa(lineNum+1),
						Evidence: "uri:" + scheme,
						Required: true,
					})
				}
			}
		}
	}
}
