package extractor

import (
	"os"
	"path/filepath"
	"testing"
)

// setupTestRepo creates a temporary directory structure with sample fixtures
// that mirror the patterns the extractors look for.
func setupTestRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	// CRD fixture
	crdDir := filepath.Join(root, "config", "crd", "bases")
	mustMkdirAll(t, crdDir)
	mustWriteFile(t, filepath.Join(crdDir, "sample_crd.yaml"), sampleCRD)

	// RBAC fixture
	rbacDir := filepath.Join(root, "config", "rbac")
	mustMkdirAll(t, rbacDir)
	mustWriteFile(t, filepath.Join(rbacDir, "role.yaml"), sampleRBAC)

	// Deployment fixture
	deployDir := filepath.Join(root, "config", "manager")
	mustMkdirAll(t, deployDir)
	mustWriteFile(t, filepath.Join(deployDir, "manager.yaml"), sampleDeployment)

	// Service fixture
	svcDir := filepath.Join(root, "config", "webhook")
	mustMkdirAll(t, svcDir)
	mustWriteFile(t, filepath.Join(svcDir, "service.yaml"), sampleService)

	// NetworkPolicy fixture
	netpolDir := filepath.Join(root, "config", "network")
	mustMkdirAll(t, netpolDir)
	mustWriteFile(t, filepath.Join(netpolDir, "networkpolicy.yaml"), sampleNetworkPolicy)

	// Controller Go file fixture
	ctrlDir := filepath.Join(root, "controllers")
	mustMkdirAll(t, ctrlDir)
	mustWriteFile(t, filepath.Join(ctrlDir, "dsc_controller.go"), sampleController)

	// go.mod fixture
	mustWriteFile(t, filepath.Join(root, "go.mod"), sampleGoMod)

	// Dockerfile fixture
	mustWriteFile(t, filepath.Join(root, "Dockerfile"), sampleDockerfile)

	// Helm values fixture
	chartDir := filepath.Join(root, "charts", "operator")
	mustMkdirAll(t, chartDir)
	mustWriteFile(t, filepath.Join(chartDir, "values.yaml"), sampleValues)

	return root
}

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("failed to create dir %s: %v", path, err)
	}
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write %s: %v", path, err)
	}
}

// --- Test fixtures (embedded as string constants) ---

const sampleCRD = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: datascienceclusters.datasciencecluster.opendatahub.io
spec:
  group: datasciencecluster.opendatahub.io
  names:
    kind: DataScienceCluster
  scope: Cluster
  versions:
    - name: v1
      served: true
      storage: true
      schema:
        openAPIV3Schema:
          type: object
          properties:
            apiVersion:
              type: string
            kind:
              type: string
            spec:
              type: object
              properties:
                components:
                  type: object
              x-kubernetes-validations:
                - rule: "has(self.components)"
                  message: "components must be specified"
`

const sampleRBAC = `apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: opendatahub-operator-manager-role
rules:
  - apiGroups: [""]
    resources: ["configmaps", "secrets", "services"]
    verbs: ["get", "list", "watch"]
  - apiGroups: ["apps"]
    resources: ["deployments"]
    verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: opendatahub-operator-manager-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: opendatahub-operator-manager-role
subjects:
  - kind: ServiceAccount
    name: opendatahub-operator
    namespace: opendatahub-operator-system
`

const sampleDeployment = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: opendatahub-operator-controller-manager
spec:
  replicas: 1
  template:
    spec:
      serviceAccountName: opendatahub-operator
      automountServiceAccountToken: true
      containers:
        - name: manager
          image: quay.io/opendatahub/opendatahub-operator:v2.20.0
          ports:
            - name: https
              containerPort: 8443
              protocol: TCP
          env:
            - name: DB_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: db-credentials
                  key: password
          securityContext:
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
            runAsNonRoot: true
            capabilities:
              drop:
                - ALL
            seccompProfile:
              type: RuntimeDefault
          resources:
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 500m
              memory: 256Mi
          volumeMounts:
            - name: tls-certs
              mountPath: /tmp/k8s-webhook-server/serving-certs
      volumes:
        - name: tls-certs
          secret:
            secretName: webhook-server-cert
`

const sampleService = `apiVersion: v1
kind: Service
metadata:
  name: opendatahub-operator-webhook-service
  annotations:
    service.beta.openshift.io/serving-cert-secret-name: webhook-server-cert
spec:
  type: ClusterIP
  selector:
    app: opendatahub-operator
  ports:
    - name: https
      port: 8443
      targetPort: 8443
      protocol: TCP
`

const sampleNetworkPolicy = `apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: opendatahub-operator-netpol
spec:
  podSelector:
    matchLabels:
      app: opendatahub-operator
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - ports:
        - port: 8443
          protocol: TCP
      from:
        - namespaceSelector:
            matchLabels:
              kubernetes.io/metadata.name: openshift-monitoring
  egress:
    - to:
        - namespaceSelector: {}
      ports:
        - port: 6443
          protocol: TCP
`

const sampleController = `package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dscv1 "github.com/opendatahub-io/opendatahub-operator/v2/api/datasciencecluster/v1"
)

//+kubebuilder:rbac:groups=datasciencecluster.opendatahub.io,resources=datascienceclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

type DataScienceClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *DataScienceClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dscv1.DataScienceCluster{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Watches(&corev1.ConfigMap{}).
		Complete(r)
}
`

const sampleGoMod = `module github.com/opendatahub-io/opendatahub-operator/v2

go 1.22

require (
	github.com/opendatahub-io/model-registry v0.2.3
	github.com/opendatahub-io/odh-model-controller v0.12.0
	sigs.k8s.io/controller-runtime v0.19.0
	k8s.io/api v0.31.0
	github.com/prometheus/client_golang v1.20.0 // indirect
)
`

const sampleDockerfile = `FROM golang:1.22 AS builder

WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o manager main.go

FROM registry.access.redhat.com/ubi9/ubi-minimal:9.4

WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532
EXPOSE 8443 8080

ENTRYPOINT ["/manager"]
`

const sampleValues = `replicaCount: 1

securityContext:
  allowPrivilegeEscalation: false
  readOnlyRootFilesystem: true
  runAsNonRoot: true
  capabilities:
    drop:
      - ALL

tls:
  enabled: true
  certManager: true
  secretName: webhook-server-cert

networkPolicy:
  enabled: true
  ingressPorts:
    - 8443
    - 8080
`

// --- Tests ---

func TestExtractAll(t *testing.T) {
	root := setupTestRepo(t)
	arch, err := ExtractAll(root, nil)
	if err != nil {
		t.Fatalf("ExtractAll failed: %v", err)
	}

	if arch.Component == "" {
		t.Error("expected non-empty component name")
	}
	if arch.AnalyzerVersion != AnalyzerVersion {
		t.Errorf("expected analyzer version %s, got %s", AnalyzerVersion, arch.AnalyzerVersion)
	}
	if arch.ExtractedAt == "" {
		t.Error("expected non-empty extracted_at timestamp")
	}

	// Verify that each extractor produced some data
	if len(arch.CRDs) == 0 {
		t.Error("expected at least one CRD")
	}
	if arch.RBAC == nil || len(arch.RBAC.ClusterRoles) == 0 {
		t.Error("expected at least one ClusterRole")
	}
	if len(arch.Services) == 0 {
		t.Error("expected at least one Service")
	}
	if len(arch.Deployments) == 0 {
		t.Error("expected at least one Deployment")
	}
	if len(arch.NetworkPolicies) == 0 {
		t.Error("expected at least one NetworkPolicy")
	}
	if len(arch.ControllerWatch) == 0 {
		t.Error("expected at least one controller watch")
	}
	if arch.Dependencies == nil || len(arch.Dependencies.GoModules) == 0 {
		t.Error("expected at least one Go module dependency")
	}
	if len(arch.Secrets) == 0 {
		t.Error("expected at least one secret reference")
	}
	if len(arch.Dockerfiles) == 0 {
		t.Error("expected at least one Dockerfile")
	}
}

func TestExtractAll_InvalidPath(t *testing.T) {
	_, err := ExtractAll("/nonexistent/path/that/does/not/exist", nil)
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestExtractCRDs(t *testing.T) {
	root := setupTestRepo(t)
	crds := extractCRDs(root)

	if len(crds) != 1 {
		t.Fatalf("expected 1 CRD, got %d", len(crds))
	}
	crd := crds[0]
	if crd.Group != "datasciencecluster.opendatahub.io" {
		t.Errorf("unexpected group: %s", crd.Group)
	}
	if crd.Kind != "DataScienceCluster" {
		t.Errorf("unexpected kind: %s", crd.Kind)
	}
	if crd.Version != "v1" {
		t.Errorf("unexpected version: %s", crd.Version)
	}
	if crd.Scope != "Cluster" {
		t.Errorf("unexpected scope: %s", crd.Scope)
	}
	if crd.FieldsCount == 0 {
		t.Error("expected non-zero fields count")
	}
	if len(crd.ValidationRules) == 0 {
		t.Error("expected at least one CEL validation rule")
	}
	if crd.ValidationRules[0] != "has(self.components)" {
		t.Errorf("unexpected CEL rule: %s", crd.ValidationRules[0])
	}
}

func TestExtractRBAC(t *testing.T) {
	root := setupTestRepo(t)
	rbac := extractRBAC(root)

	if len(rbac.ClusterRoles) != 1 {
		t.Fatalf("expected 1 ClusterRole, got %d", len(rbac.ClusterRoles))
	}
	cr := rbac.ClusterRoles[0]
	if cr.Name != "opendatahub-operator-manager-role" {
		t.Errorf("unexpected ClusterRole name: %s", cr.Name)
	}
	if len(cr.Rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(cr.Rules))
	}

	if len(rbac.ClusterRoleBindings) != 1 {
		t.Fatalf("expected 1 ClusterRoleBinding, got %d", len(rbac.ClusterRoleBindings))
	}
	crb := rbac.ClusterRoleBindings[0]
	if crb.RoleRef != "opendatahub-operator-manager-role" {
		t.Errorf("unexpected roleRef: %s", crb.RoleRef)
	}
	if len(crb.Subjects) != 1 {
		t.Errorf("expected 1 subject, got %d", len(crb.Subjects))
	}

	// Check kubebuilder markers from controller Go file
	if len(rbac.KubebuilderMarkers) < 2 {
		t.Errorf("expected at least 2 kubebuilder markers, got %d", len(rbac.KubebuilderMarkers))
	}
}

func TestExtractDeployments(t *testing.T) {
	root := setupTestRepo(t)
	deps := extractDeployments(root)

	if len(deps) != 1 {
		t.Fatalf("expected 1 deployment, got %d", len(deps))
	}
	dep := deps[0]
	if dep.Name != "opendatahub-operator-controller-manager" {
		t.Errorf("unexpected name: %s", dep.Name)
	}
	if dep.Kind != "Deployment" {
		t.Errorf("unexpected kind: %s", dep.Kind)
	}
	if dep.ServiceAccount != "opendatahub-operator" {
		t.Errorf("unexpected service account: %s", dep.ServiceAccount)
	}
	if len(dep.Containers) != 1 {
		t.Fatalf("expected 1 container, got %d", len(dep.Containers))
	}
	c := dep.Containers[0]
	if c.Name != "manager" {
		t.Errorf("unexpected container name: %s", c.Name)
	}
	if len(c.EnvFromSecrets) == 0 {
		t.Error("expected at least one env-from-secret reference")
	}
}

func TestExtractDockerfiles(t *testing.T) {
	root := setupTestRepo(t)
	dfs := extractDockerfiles(root)

	if len(dfs) != 1 {
		t.Fatalf("expected 1 Dockerfile, got %d", len(dfs))
	}
	df := dfs[0]
	if df.BaseImage != "registry.access.redhat.com/ubi9/ubi-minimal:9.4" {
		t.Errorf("unexpected base image (should be runtime stage): %s", df.BaseImage)
	}
	if df.Stages != 2 {
		t.Errorf("expected 2 stages, got %d", df.Stages)
	}
	if df.User != "65532:65532" {
		t.Errorf("unexpected user: %s", df.User)
	}
	if len(df.ExposedPorts) != 2 {
		t.Errorf("expected 2 exposed ports, got %d", len(df.ExposedPorts))
	}
}

func TestExtractControllerWatches(t *testing.T) {
	root := setupTestRepo(t)
	watches := extractControllerWatches(root)

	if len(watches) < 4 {
		t.Fatalf("expected at least 4 watches, got %d", len(watches))
	}

	// Check that we found the For, Owns, and Watches calls
	typeCount := map[string]int{}
	for _, w := range watches {
		typeCount[w.Type]++
	}
	if typeCount["For"] < 1 {
		t.Error("expected at least one For watch")
	}
	if typeCount["Owns"] < 2 {
		t.Errorf("expected at least 2 Owns watches, got %d", typeCount["Owns"])
	}
	if typeCount["Watches"] < 1 {
		t.Error("expected at least one Watches watch")
	}
}

func TestExtractDependencies(t *testing.T) {
	root := setupTestRepo(t)
	deps := extractDependencies(root, DefaultModulePrefixes())

	if len(deps.GoModules) == 0 {
		t.Fatal("expected at least one Go module")
	}
	if len(deps.InternalODH) < 2 {
		t.Errorf("expected at least 2 internal ODH deps, got %d", len(deps.InternalODH))
	}

	// Verify indirect deps are filtered
	for _, m := range deps.GoModules {
		if m.Module == "github.com/prometheus/client_golang" {
			t.Error("indirect dependency should have been filtered out")
		}
	}
}

func TestExtractServices(t *testing.T) {
	root := setupTestRepo(t)
	svcs := extractServices(root)

	if len(svcs) != 1 {
		t.Fatalf("expected 1 service, got %d", len(svcs))
	}
	svc := svcs[0]
	if svc.Name != "opendatahub-operator-webhook-service" {
		t.Errorf("unexpected service name: %s", svc.Name)
	}
	if svc.Type != "ClusterIP" {
		t.Errorf("unexpected service type: %s", svc.Type)
	}
	if len(svc.Ports) != 1 {
		t.Errorf("expected 1 port, got %d", len(svc.Ports))
	}
}

func TestExtractNetworkPolicies(t *testing.T) {
	root := setupTestRepo(t)
	policies := extractNetworkPolicies(root)

	if len(policies) != 1 {
		t.Fatalf("expected 1 network policy, got %d", len(policies))
	}
	pol := policies[0]
	if pol.Name != "opendatahub-operator-netpol" {
		t.Errorf("unexpected name: %s", pol.Name)
	}
	if len(pol.PolicyTypes) != 2 {
		t.Errorf("expected 2 policy types, got %d", len(pol.PolicyTypes))
	}
	if len(pol.IngressRules) == 0 {
		t.Error("expected at least one ingress rule")
	}
	if len(pol.EgressRules) == 0 {
		t.Error("expected at least one egress rule")
	}
}

func TestExtractSecrets(t *testing.T) {
	root := setupTestRepo(t)
	secrets := extractSecrets(root)

	if len(secrets) == 0 {
		t.Fatal("expected at least one secret reference")
	}

	secretNames := map[string]bool{}
	for _, s := range secrets {
		secretNames[s.Name] = true
	}
	if !secretNames["db-credentials"] {
		t.Error("expected db-credentials secret reference")
	}
	if !secretNames["webhook-server-cert"] {
		t.Error("expected webhook-server-cert secret reference")
	}
}

func TestExtractHelm(t *testing.T) {
	root := setupTestRepo(t)
	helm := extractHelm(root)

	if len(helm.ValuesDefaults) == 0 {
		t.Error("expected non-empty values defaults")
	}
	// Check that securityContext values were flattened
	if _, ok := helm.ValuesDefaults["securityContext.runAsNonRoot"]; !ok {
		t.Error("expected securityContext.runAsNonRoot in values defaults")
	}
	if _, ok := helm.ValuesDefaults["tls.enabled"]; !ok {
		t.Error("expected tls.enabled in values defaults")
	}
}
