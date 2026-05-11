package extractor

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestClassifyEnvVar(t *testing.T) {
	tests := []struct {
		envName  string
		wantName string
		wantType string
		wantOK   bool
	}{
		{"DATABASE_HOST", "PostgreSQL", "database", true},
		{"REDIS_URL", "Redis", "cache", true},
		{"S3_ENDPOINT", "S3", "object-storage", true},
		{"KAFKA_BOOTSTRAP_SERVERS", "Kafka", "messaging", true},
		{"OTEL_EXPORTER_OTLP_ENDPOINT", "OpenTelemetry Collector", "observability", true},
		{"MLFLOW_TRACKING_URI", "MLflow", "service", true},
		{"MY_APP_PORT", "", "", false},
		// prefix match
		{"POSTGRES_PASSWORD", "PostgreSQL", "database", true},
		{"REDIS_SENTINEL_HOST", "Redis", "cache", true},
		{"KAFKA_TOPIC", "Kafka", "messaging", true},
	}
	for _, tt := range tests {
		cls, ok := classifyEnvVar(tt.envName)
		if ok != tt.wantOK {
			t.Errorf("classifyEnvVar(%q): got ok=%v, want %v", tt.envName, ok, tt.wantOK)
			continue
		}
		if !ok {
			continue
		}
		if cls.Name != tt.wantName {
			t.Errorf("classifyEnvVar(%q): got name=%q, want %q", tt.envName, cls.Name, tt.wantName)
		}
		if cls.Type != tt.wantType {
			t.Errorf("classifyEnvVar(%q): got type=%q, want %q", tt.envName, cls.Type, tt.wantType)
		}
	}
}

func TestClassifyEnvValue(t *testing.T) {
	tests := []struct {
		value    string
		wantName string
		wantOK   bool
	}{
		{"postgresql://user:pass@host:5432/db", "PostgreSQL", true},
		{"redis://localhost:6379", "Redis", true},
		{"amqp://guest:guest@rabbit:5672", "RabbitMQ", true},
		{"https://example.com", "", false},
		{"some-plain-value", "", false},
	}
	for _, tt := range tests {
		cls, ok := classifyEnvValue(tt.value)
		if ok != tt.wantOK {
			t.Errorf("classifyEnvValue(%q): got ok=%v, want %v", tt.value, ok, tt.wantOK)
			continue
		}
		if ok && cls.Name != tt.wantName {
			t.Errorf("classifyEnvValue(%q): got name=%q, want %q", tt.value, cls.Name, tt.wantName)
		}
	}
}

func TestDeduplication(t *testing.T) {
	seen := make(map[string]RuntimeDependency)

	// Record PostgreSQL twice from different sources
	recordDependency(seen, RuntimeDependency{
		Name:     "PostgreSQL",
		Type:     "database",
		Source:   "deploy.yaml",
		Evidence: "env:DATABASE_HOST",
		Required: false,
	})
	recordDependency(seen, RuntimeDependency{
		Name:     "PostgreSQL",
		Type:     "database",
		Source:   "main.go:42",
		Evidence: "env:POSTGRES_HOST",
		Required: true,
	})

	if len(seen) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(seen))
	}

	dep := seen["PostgreSQL"]
	// The Required flag should be upgraded to true
	if !dep.Required {
		t.Error("expected Required to be upgraded to true")
	}
	// Source should be from the first detection (deploy.yaml)
	if dep.Source != "deploy.yaml" {
		t.Errorf("expected source from first detection, got %q", dep.Source)
	}
}

func TestExtractRuntimeDependenciesFromYAML(t *testing.T) {
	tmp := t.TempDir()

	// Create a deployment YAML with known env vars
	deployYAML := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-app
spec:
  template:
    spec:
      containers:
        - name: app
          image: test:latest
          env:
            - name: DATABASE_HOST
              value: "postgres.default.svc.cluster.local"
            - name: REDIS_URL
              value: "redis://redis:6379"
            - name: S3_ENDPOINT
            - name: MY_CUSTOM_VAR
              value: "something"
          envFrom:
            - secretRef:
                name: postgres-credentials
            - configMapRef:
                name: kafka-config
`
	err := os.WriteFile(filepath.Join(tmp, "deployment.yaml"), []byte(deployYAML), 0644)
	if err != nil {
		t.Fatal(err)
	}

	deps := extractRuntimeDependencies(tmp)

	// Build a map for easier assertions
	byName := make(map[string]RuntimeDependency)
	for _, d := range deps {
		byName[d.Name] = d
	}

	// PostgreSQL should be detected from DATABASE_HOST env var
	if _, ok := byName["PostgreSQL"]; !ok {
		t.Error("expected PostgreSQL dependency to be detected")
	}

	// Redis from REDIS_URL
	if _, ok := byName["Redis"]; !ok {
		t.Error("expected Redis dependency to be detected")
	}

	// S3 from S3_ENDPOINT (no value = required)
	if s3, ok := byName["S3"]; !ok {
		t.Error("expected S3 dependency to be detected")
	} else if !s3.Required {
		t.Error("expected S3 to be required (no default value)")
	}

	// Kafka from configMapRef name
	if _, ok := byName["Kafka"]; !ok {
		t.Error("expected Kafka dependency to be detected from configMapRef name")
	}

	// MY_CUSTOM_VAR should NOT produce a dependency
	for _, d := range deps {
		if d.Evidence == "env:MY_CUSTOM_VAR" {
			t.Error("MY_CUSTOM_VAR should not produce a runtime dependency")
		}
	}
}

func TestExtractRuntimeDependenciesFromGoSource(t *testing.T) {
	tmp := t.TempDir()

	goSource := `package main

import "os"

func main() {
	dbHost := os.Getenv("DATABASE_HOST")
	redisAddr := os.Getenv("REDIS_ADDR")
	appPort := os.Getenv("PORT")
	_ = dbHost
	_ = redisAddr
	_ = appPort
}
`
	err := os.WriteFile(filepath.Join(tmp, "main.go"), []byte(goSource), 0644)
	if err != nil {
		t.Fatal(err)
	}

	deps := extractRuntimeDependencies(tmp)

	byName := make(map[string]RuntimeDependency)
	for _, d := range deps {
		byName[d.Name] = d
	}

	if _, ok := byName["PostgreSQL"]; !ok {
		t.Error("expected PostgreSQL from os.Getenv(\"DATABASE_HOST\")")
	}
	if _, ok := byName["Redis"]; !ok {
		t.Error("expected Redis from os.Getenv(\"REDIS_ADDR\")")
	}

	// PORT should not match any known service
	for _, d := range deps {
		if d.Evidence == "env:PORT" {
			t.Error("PORT should not produce a runtime dependency")
		}
	}
}

func TestExtractRuntimeDependenciesFromPython(t *testing.T) {
	tmp := t.TempDir()

	pySource := `import os

db_url = os.environ.get("POSTGRES_HOST", "localhost")
mlflow_uri = os.environ["MLFLOW_TRACKING_URI"]
`
	err := os.WriteFile(filepath.Join(tmp, "app.py"), []byte(pySource), 0644)
	if err != nil {
		t.Fatal(err)
	}

	deps := extractRuntimeDependencies(tmp)

	byName := make(map[string]RuntimeDependency)
	for _, d := range deps {
		byName[d.Name] = d
	}

	if _, ok := byName["PostgreSQL"]; !ok {
		t.Error("expected PostgreSQL from os.environ.get(\"POSTGRES_HOST\")")
	}
	if _, ok := byName["MLflow"]; !ok {
		t.Error("expected MLflow from os.environ[\"MLFLOW_TRACKING_URI\"]")
	}
}

func TestSortRuntimeDependencies(t *testing.T) {
	deps := []RuntimeDependency{
		{Name: "Redis", Type: "cache"},
		{Name: "PostgreSQL", Type: "database"},
		{Name: "S3", Type: "object-storage"},
		{Name: "Kafka", Type: "messaging"},
	}

	sort.Slice(deps, func(i, j int) bool {
		if deps[i].Type != deps[j].Type {
			return deps[i].Type < deps[j].Type
		}
		return deps[i].Name < deps[j].Name
	})

	expected := []string{"Redis", "PostgreSQL", "Kafka", "S3"}
	for i, d := range deps {
		if d.Name != expected[i] {
			t.Errorf("position %d: got %q, want %q", i, d.Name, expected[i])
		}
	}
}
