package parser

import (
	"os"
	"strings"
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

func TestRustParserLanguageAndExtensions(t *testing.T) {
	p := NewRustParser()
	if p.Language() != "rust" {
		t.Errorf("expected language 'rust', got %q", p.Language())
	}
	exts := p.Extensions()
	if len(exts) != 1 || exts[0] != ".rs" {
		t.Errorf("expected extensions [.rs], got %v", exts)
	}
}

func TestRustParserActixHandler(t *testing.T) {
	content, err := os.ReadFile("../../testdata/actix_handler.rs")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	p := NewRustParser()
	result, err := p.ParseFile("testdata/actix_handler.rs", content)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	// Debug: log all extracted nodes
	t.Logf("Functions (%d):", len(result.Functions))
	for _, fn := range result.Functions {
		t.Logf("  %s (type=%s, props=%v)", fn.Name, fn.TypeName, fn.Properties)
	}
	t.Logf("HTTP handlers (%d):", len(result.HTTPHandlers))
	for _, h := range result.HTTPHandlers {
		t.Logf("  %s route=%s method=%s", h.Name, h.Route, h.HTTPMethod)
	}
	t.Logf("Call sites (%d):", len(result.CallSites))
	for _, cs := range result.CallSites {
		t.Logf("  %s (props=%v)", cs.Name, cs.Properties)
	}
	t.Logf("DB operations (%d):", len(result.DBOperations))
	for _, db := range result.DBOperations {
		t.Logf("  %s (props=%v)", db.Name, db.Properties)
	}

	// At least 4 functions: health_check, create_item, list_items, get_api_key
	if len(result.Functions) < 4 {
		t.Errorf("expected at least 4 functions, got %d", len(result.Functions))
	}

	fnNames := make(map[string]bool)
	for _, fn := range result.Functions {
		fnNames[fn.Name] = true
	}
	for _, expected := range []string{"health_check", "create_item", "list_items", "get_api_key"} {
		if !fnNames[expected] {
			t.Errorf("expected function %q not found", expected)
		}
	}

	// At least 3 HTTP handlers with routes
	if len(result.HTTPHandlers) < 3 {
		t.Errorf("expected at least 3 HTTP handlers, got %d", len(result.HTTPHandlers))
	}

	routes := make(map[string]bool)
	for _, h := range result.HTTPHandlers {
		if h.Route != "" {
			routes[h.Route] = true
		}
	}
	if !routes["/health"] {
		t.Error("expected HTTP handler with route /health")
	}
	if !routes["/items"] {
		t.Error("expected HTTP handler with route /items")
	}

	// Non-zero call sites
	if len(result.CallSites) == 0 {
		t.Error("expected call sites, got 0")
	}

	// At least 1 DB operation (diesel::insert_into)
	if len(result.DBOperations) < 1 {
		t.Errorf("expected at least 1 DB operation, got %d", len(result.DBOperations))
	}
}

func TestRustParserAxumRouter(t *testing.T) {
	content, err := os.ReadFile("../../testdata/axum_router.rs")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	p := NewRustParser()
	result, err := p.ParseFile("testdata/axum_router.rs", content)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	// Debug: log all extracted nodes
	t.Logf("Functions (%d):", len(result.Functions))
	for _, fn := range result.Functions {
		t.Logf("  %s (type=%s, props=%v)", fn.Name, fn.TypeName, fn.Properties)
	}
	t.Logf("Call sites (%d):", len(result.CallSites))
	for _, cs := range result.CallSites {
		t.Logf("  %s (props=%v)", cs.Name, cs.Properties)
	}
	t.Logf("DB operations (%d):", len(result.DBOperations))
	for _, db := range result.DBOperations {
		t.Logf("  %s (props=%v)", db.Name, db.Properties)
	}
	t.Logf("Struct literals (%d):", len(result.StructLiterals))
	for _, sl := range result.StructLiterals {
		t.Logf("  %s (props=%v)", sl.Name, sl.Properties)
	}

	// At least 4 functions: list_users, create_user, raw_pointer_op, ffi_callback, build_router
	if len(result.Functions) < 4 {
		t.Errorf("expected at least 4 functions, got %d", len(result.Functions))
	}

	fnMap := make(map[string]*graph.Node)
	for _, fn := range result.Functions {
		fnMap[fn.Name] = fn
	}

	// Unsafe function detection
	if fn, ok := fnMap["raw_pointer_op"]; ok {
		if !fn.IsUnsafe {
			t.Errorf("expected raw_pointer_op to have IsUnsafe=true, got %v", fn.IsUnsafe)
		}
	} else {
		t.Error("expected to find function raw_pointer_op")
	}

	// Extern function detection
	if fn, ok := fnMap["ffi_callback"]; ok {
		if !fn.IsExtern {
			t.Errorf("expected ffi_callback to have IsExtern=true, got %v", fn.IsExtern)
		}
	} else {
		t.Error("expected to find function ffi_callback")
	}

	// Macro invocations (sqlx::query_as!) should be detected as call sites with IsMacro
	hasMacro := false
	for _, cs := range result.CallSites {
		if cs.IsMacro {
			hasMacro = true
			break
		}
	}
	if !hasMacro {
		t.Error("expected at least one macro invocation with IsMacro=true")
	}

	// At least 2 DB operations from sqlx macros
	if len(result.DBOperations) < 2 {
		t.Errorf("expected at least 2 DB operations, got %d", len(result.DBOperations))
	}

	// DB operations should have read/write classification
	hasRead, hasWrite := false, false
	for _, db := range result.DBOperations {
		switch db.Operation {
		case "read":
			hasRead = true
		case "write":
			hasWrite = true
		}
	}
	if !hasRead {
		t.Error("expected at least one DB read operation")
	}
	if !hasWrite {
		t.Error("expected at least one DB write operation")
	}

	// Struct literals: log whatever we find (axum_router.rs may not have struct expressions)
	t.Logf("Struct literals found: %d", len(result.StructLiterals))
	for _, sl := range result.StructLiterals {
		t.Logf("  struct: %s (props=%v)", sl.Name, sl.Properties)
	}
}

func TestRustParserSkipDirs(t *testing.T) {
	if !ShouldSkipDir("target", RustSkipDirs) {
		t.Error("expected 'target' to be skipped for Rust")
	}
	if ShouldSkipDir("src", RustSkipDirs) {
		t.Error("did not expect 'src' to be skipped for Rust")
	}
}

func TestRustParserAllFunctionsHaveLanguage(t *testing.T) {
	content, err := os.ReadFile("../../testdata/actix_handler.rs")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	p := NewRustParser()
	result, err := p.ParseFile("testdata/actix_handler.rs", content)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	for _, fn := range result.Functions {
		if fn.Language != "rust" {
			t.Errorf("function %s has language %q, expected 'rust'", fn.Name, fn.Language)
		}
	}
}

func TestRustParserFileTooLarge(t *testing.T) {
	p := NewRustParser()
	bigContent := make([]byte, maxFileSize+1)
	_, err := p.ParseFile("big.rs", bigContent)
	if err == nil {
		t.Error("expected error for oversized file")
	}
	if !strings.Contains(err.Error(), "too large") {
		t.Errorf("expected 'too large' in error, got %q", err.Error())
	}
}
