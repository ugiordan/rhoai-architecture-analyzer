package parser

import (
	"os"
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

func TestTypeScriptParserLanguageAndExtensions(t *testing.T) {
	p := NewTypeScriptParser()
	if p.Language() != "typescript" {
		t.Errorf("expected language 'typescript', got %q", p.Language())
	}
	exts := p.Extensions()
	if len(exts) != 2 || exts[0] != ".ts" || exts[1] != ".tsx" {
		t.Errorf("expected extensions [.ts .tsx], got %v", exts)
	}
}

func TestTypeScriptParserExpressServer(t *testing.T) {
	content, err := os.ReadFile("../../testdata/express_server.ts")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	p := NewTypeScriptParser()
	result, err := p.ParseFile("testdata/express_server.ts", content)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	// Functions: at least 2 (authMiddleware, startServer)
	if len(result.Functions) < 2 {
		t.Errorf("expected at least 2 functions, got %d", len(result.Functions))
	}
	funcNames := make(map[string]bool)
	for _, fn := range result.Functions {
		funcNames[fn.Name] = true
		t.Logf("  function: %s (kind=%s, type=%s)", fn.Name, fn.Kind, fn.TypeName)
	}
	for _, expected := range []string{"authMiddleware", "startServer"} {
		if !funcNames[expected] {
			t.Errorf("expected function %q not found", expected)
		}
	}

	// All functions should have correct language
	for _, fn := range result.Functions {
		if fn.Language != "typescript" {
			t.Errorf("function %q has language %q, expected 'typescript'", fn.Name, fn.Language)
		}
		if fn.Kind != graph.NodeFunction {
			t.Errorf("function %q has kind %q, expected NodeFunction", fn.Name, fn.Kind)
		}
	}

	// HTTP handlers: at least 4 with routes
	if len(result.HTTPHandlers) < 4 {
		t.Errorf("expected at least 4 HTTP handlers, got %d", len(result.HTTPHandlers))
	}
	routes := make(map[string]bool)
	for _, h := range result.HTTPHandlers {
		routes[h.Route] = true
		t.Logf("  http handler: %s route=%s", h.Name, h.Route)
	}
	for _, expected := range []string{"/users", "/search"} {
		if !routes[expected] {
			t.Errorf("expected HTTP route %q not found, got routes: %v", expected, routes)
		}
	}

	// DB operations: at least 3 (pool.query x3)
	if len(result.DBOperations) < 3 {
		t.Errorf("expected at least 3 DB operations, got %d", len(result.DBOperations))
	}
	for _, op := range result.DBOperations {
		t.Logf("  db op: %s (op=%s)", op.Name, op.Operation)
	}

	// Call sites: non-zero
	if len(result.CallSites) == 0 {
		t.Error("expected call sites, got 0")
	}
	for _, cs := range result.CallSites {
		t.Logf("  call site: %s", cs.Name)
	}

	// Struct literals: new Pool() should be detected
	if len(result.StructLiterals) == 0 {
		t.Error("expected struct literals (new Pool()), got 0")
	}
	foundPool := false
	for _, sl := range result.StructLiterals {
		t.Logf("  struct literal: %s", sl.Name)
		if sl.Name == "Pool" {
			foundPool = true
		}
	}
	if !foundPool {
		t.Error("expected struct literal 'Pool' from new Pool()")
	}
}

func TestTypeScriptParserReactComponent(t *testing.T) {
	content, err := os.ReadFile("../../testdata/react_component.tsx")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	p := NewTypeScriptParser()
	result, err := p.ParseFile("testdata/react_component.tsx", content)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	// Arrow functions: Dashboard, UserProfile, AppRoutes
	funcNames := make(map[string]bool)
	for _, fn := range result.Functions {
		funcNames[fn.Name] = true
		t.Logf("  function: %s (kind=%s)", fn.Name, fn.Kind)
	}
	for _, expected := range []string{"Dashboard", "UserProfile", "AppRoutes"} {
		if !funcNames[expected] {
			t.Errorf("expected arrow function %q not found", expected)
		}
	}

	// React Router: at least 2 Route elements with paths
	if len(result.HTTPHandlers) < 2 {
		t.Errorf("expected at least 2 React Router handlers, got %d", len(result.HTTPHandlers))
	}
	routes := make(map[string]bool)
	for _, h := range result.HTTPHandlers {
		routes[h.Route] = true
		t.Logf("  react route: name=%s route=%s component=%s", h.Name, h.Route, h.Properties["component"])
	}
	if !routes["/dashboard"] {
		t.Errorf("expected route /dashboard not found, got routes: %v", routes)
	}
	if !routes["/users/:id"] {
		t.Errorf("expected route /users/:id not found, got routes: %v", routes)
	}
}

func TestTypeScriptParserSkipsDeclarationFiles(t *testing.T) {
	content := []byte(`
declare module 'express' {
    export interface Request {}
    export interface Response {}
}
`)
	p := NewTypeScriptParser()
	result, err := p.ParseFile("types/express.d.ts", content)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}
	if len(result.Functions) != 0 {
		t.Errorf("expected 0 functions for .d.ts file, got %d", len(result.Functions))
	}
	if len(result.CallSites) != 0 {
		t.Errorf("expected 0 call sites for .d.ts file, got %d", len(result.CallSites))
	}
}
