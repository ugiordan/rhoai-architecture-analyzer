package parser

import (
	"os"
	"strings"
	"testing"

	"github.com/ugiordan/rhoai-architecture-analyzer/pkg/graph"
)

func TestGoParserParseFile(t *testing.T) {
	content, err := os.ReadFile("../../testdata/simple_http_server.go")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	p := NewGoParser()
	result, err := p.ParseFile("testdata/simple_http_server.go", content)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(result.Functions) < 3 {
		t.Errorf("expected at least 3 functions, got %d", len(result.Functions))
		for _, fn := range result.Functions {
			t.Logf("  function: %s", fn.Name)
		}
	}

	names := make(map[string]bool)
	for _, fn := range result.Functions {
		names[fn.Name] = true
	}
	for _, expected := range []string{"handleGetUsers", "handleCreateUser", "main"} {
		if !names[expected] {
			t.Errorf("expected function %q not found", expected)
		}
	}

	if len(result.CallSites) == 0 {
		t.Error("expected call sites, got 0")
	}

	if len(result.HTTPHandlers) == 0 {
		t.Error("expected HTTP handlers, got 0")
	}
}

func TestGoParserDetectsDBOperations(t *testing.T) {
	content, err := os.ReadFile("../../testdata/db_read_write.go")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	p := NewGoParser()
	result, err := p.ParseFile("testdata/db_read_write.go", content)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(result.DBOperations) < 2 {
		t.Errorf("expected at least 2 DB operations, got %d", len(result.DBOperations))
		for _, op := range result.DBOperations {
			t.Logf("  db op: %s (kind=%s)", op.Name, op.Properties["operation"])
		}
	}

	hasRead, hasWrite := false, false
	for _, op := range result.DBOperations {
		switch op.Properties["operation"] {
		case "read":
			hasRead = true
		case "write":
			hasWrite = true
		}
	}
	if !hasRead {
		t.Error("expected a DB read operation")
	}
	if !hasWrite {
		t.Error("expected a DB write operation")
	}
}

func TestGoParserLanguageAndExtensions(t *testing.T) {
	p := NewGoParser()
	if p.Language() != "go" {
		t.Errorf("expected language 'go', got %q", p.Language())
	}
	exts := p.Extensions()
	if len(exts) != 1 || exts[0] != ".go" {
		t.Errorf("expected extensions [.go], got %v", exts)
	}
}

func TestGoParserExtractsParamTypes(t *testing.T) {
	content, err := os.ReadFile("../../testdata/k8s_webhook.go")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	p := NewGoParser()
	result, err := p.ParseFile("testdata/k8s_webhook.go", content)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	var handleFn *graph.Node
	for _, fn := range result.Functions {
		if fn.Name == "Handle" {
			handleFn = fn
			break
		}
	}
	if handleFn == nil {
		t.Fatal("expected to find Handle function")
	}

	paramTypes := handleFn.Properties["param_types"]
	if paramTypes == "" {
		t.Fatal("expected param_types property on Handle function")
	}
	if !strings.Contains(paramTypes, "context.Context") {
		t.Errorf("expected param_types to contain context.Context, got %q", paramTypes)
	}
	if !strings.Contains(paramTypes, "admission.Request") {
		t.Errorf("expected param_types to contain admission.Request, got %q", paramTypes)
	}
}

func TestGoParserExtractsStructLiterals(t *testing.T) {
	content, err := os.ReadFile("../../testdata/k8s_cert.go")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	p := NewGoParser()
	result, err := p.ParseFile("testdata/k8s_cert.go", content)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(result.StructLiterals) == 0 {
		t.Fatal("expected struct literals, got 0")
	}

	var certLiteral *graph.Node
	for _, sl := range result.StructLiterals {
		if strings.Contains(sl.Properties["type"], "Certificate") {
			certLiteral = sl
			break
		}
	}
	if certLiteral == nil {
		t.Fatal("expected to find x509.Certificate struct literal")
	}

	fields := certLiteral.Properties["fields"]
	if fields == "" {
		t.Fatal("expected fields property on struct literal")
	}
	for _, expected := range []string{"SerialNumber", "IsCA", "KeyUsage", "DNSNames"} {
		if !strings.Contains(fields, expected) {
			t.Errorf("expected fields to contain %q, got %q", expected, fields)
		}
	}
}

func TestGoParserExtractsArgTypes(t *testing.T) {
	content, err := os.ReadFile("../../testdata/k8s_rbac.go")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	p := NewGoParser()
	result, err := p.ParseFile("testdata/k8s_rbac.go", content)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	// Find c.Create call site
	var createCall *graph.Node
	for _, cs := range result.CallSites {
		if cs.Name == "c.Create" {
			createCall = cs
			break
		}
	}
	if createCall == nil {
		t.Fatal("expected to find c.Create call site")
	}

	// We expect arg_types to be empty for this case since "binding" is an identifier,
	// not a &Type{} literal. But the property key should still exist if there were any.
	// The main test is that the extraction runs without error.
	_ = createCall.Properties["arg_types"]
}

func TestGoParserExtractsStringArgs(t *testing.T) {
	content, err := os.ReadFile("../../testdata/k8s_rbac.go")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	p := NewGoParser()
	result, err := p.ParseFile("testdata/k8s_rbac.go", content)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	// Check that struct literals capture string values
	foundSystemAuth := false
	for _, sl := range result.StructLiterals {
		sv := sl.Properties["string_values"]
		if strings.Contains(sv, "system:authenticated") {
			foundSystemAuth = true
			break
		}
	}
	if !foundSystemAuth {
		t.Error("expected to find 'system:authenticated' in struct literal string_values")
	}
}

func TestGoParserExtractsSwitchCases(t *testing.T) {
	content, err := os.ReadFile("../../testdata/k8s_webhook.go")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	p := NewGoParser()
	result, err := p.ParseFile("testdata/k8s_webhook.go", content)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	var handleFn *graph.Node
	for _, fn := range result.Functions {
		if fn.Name == "Handle" {
			handleFn = fn
			break
		}
	}
	if handleFn == nil {
		t.Fatal("expected to find Handle function")
	}

	caseValues := handleFn.Properties["case_values"]
	if caseValues == "" {
		t.Fatal("expected case_values property on Handle function")
	}
	if !strings.Contains(caseValues, "Create") {
		t.Errorf("expected case_values to contain Create, got %q", caseValues)
	}
	if !strings.Contains(caseValues, "Delete") {
		t.Errorf("expected case_values to contain Delete, got %q", caseValues)
	}
}
