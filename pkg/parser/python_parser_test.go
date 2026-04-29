package parser

import (
	"os"
	"testing"

	"github.com/ugiordan/architecture-analyzer/pkg/graph"
)

func TestPythonParserLanguageAndExtensions(t *testing.T) {
	p := NewPythonParser()
	if p.Language() != "python" {
		t.Errorf("expected language 'python', got %q", p.Language())
	}
	exts := p.Extensions()
	if len(exts) != 1 || exts[0] != ".py" {
		t.Errorf("expected extensions [.py], got %v", exts)
	}
}

func TestPythonParserFlaskApp(t *testing.T) {
	content, err := os.ReadFile("../../testdata/flask_app.py")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	p := NewPythonParser()
	result, err := p.ParseFile("testdata/flask_app.py", content)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	// --- Functions ---
	t.Run("functions", func(t *testing.T) {
		if len(result.Functions) < 5 {
			t.Errorf("expected at least 5 functions, got %d", len(result.Functions))
			for _, fn := range result.Functions {
				t.Logf("  function: %s (type=%s)", fn.Name, fn.TypeName)
			}
		}

		fnMap := make(map[string]*graph.Node)
		for _, fn := range result.Functions {
			fnMap[fn.Name] = fn
		}

		for _, expected := range []string{"get_users", "create_user", "delete_user", "run_backup", "run_migration"} {
			if fnMap[expected] == nil {
				t.Errorf("expected function %q not found", expected)
			}
		}

		// All functions should have language "python"
		for _, fn := range result.Functions {
			if fn.Language != "python" {
				t.Errorf("function %q has language %q, expected 'python'", fn.Name, fn.Language)
			}
		}
	})

	// --- Methods with TypeName ---
	t.Run("methods", func(t *testing.T) {
		fnMap := make(map[string]*graph.Node)
		for _, fn := range result.Functions {
			fnMap[fn.Name] = fn
		}

		for _, method := range []string{"get_all", "create"} {
			fn := fnMap[method]
			if fn == nil {
				t.Errorf("expected method %q not found", method)
				continue
			}
			if fn.TypeName != "UserService" {
				t.Errorf("method %q TypeName = %q, want 'UserService'", method, fn.TypeName)
			}
		}
	})

	// --- HTTP handlers ---
	t.Run("http_handlers", func(t *testing.T) {
		if len(result.HTTPHandlers) < 3 {
			t.Errorf("expected at least 3 HTTP handlers, got %d", len(result.HTTPHandlers))
			for _, h := range result.HTTPHandlers {
				t.Logf("  handler: %s route=%s", h.Name, h.Route)
			}
		}

		hasUsersRoute := false
		for _, h := range result.HTTPHandlers {
			if h.Route == "/users" {
				hasUsersRoute = true
				break
			}
		}
		if !hasUsersRoute {
			t.Error("expected an HTTP handler with route '/users'")
		}
	})

	// --- Call sites ---
	t.Run("call_sites", func(t *testing.T) {
		if len(result.CallSites) == 0 {
			t.Error("expected call sites, got 0")
		}

		hasSubprocessRun := false
		for _, cs := range result.CallSites {
			if cs.Name == "subprocess.run" {
				hasSubprocessRun = true
				break
			}
		}
		if !hasSubprocessRun {
			t.Error("expected call site 'subprocess.run'")
			for _, cs := range result.CallSites {
				t.Logf("  call: %s", cs.Name)
			}
		}
	})

	// --- DB operations ---
	t.Run("db_operations", func(t *testing.T) {
		if len(result.DBOperations) < 2 {
			t.Errorf("expected at least 2 DB operations, got %d", len(result.DBOperations))
			for _, op := range result.DBOperations {
				t.Logf("  db op: %s (op=%s)", op.Name, op.Operation)
			}
		}

		hasRead, hasWrite := false, false
		for _, op := range result.DBOperations {
			switch op.Operation {
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
	})

	// --- Struct literals (class instantiations) ---
	t.Run("struct_literals", func(t *testing.T) {
		if len(result.StructLiterals) == 0 {
			t.Error("expected struct literals (class instantiations), got 0")
		}

		names := make(map[string]bool)
		for _, sl := range result.StructLiterals {
			names[sl.Name] = true
		}

		if !names["UserService"] && !names["User"] {
			t.Error("expected UserService or User class instantiation")
			for _, sl := range result.StructLiterals {
				t.Logf("  struct literal: %s", sl.Name)
			}
		}
	})

	// --- Decorators ---
	t.Run("decorators", func(t *testing.T) {
		var getUsersFn *graph.Node
		for _, fn := range result.Functions {
			if fn.Name == "get_users" {
				getUsersFn = fn
				break
			}
		}
		if getUsersFn == nil {
			t.Fatal("expected to find get_users function")
		}
		if len(getUsersFn.Decorators) == 0 {
			t.Error("expected get_users to have decorators")
		}
	})
}

func TestPythonParserComputesComplexity(t *testing.T) {
	content, err := os.ReadFile("../../testdata/complexity_sample.py")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	p := NewPythonParser()
	result, err := p.ParseFile("testdata/complexity_sample.py", content)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	expected := map[string]int{
		"simple_func":        1,
		"complex_func":       6, // if + and + elif + for + if + base (or not counted)
		"loop_func":          4, // while + except + for + base
		"comprehension_func": 3, // if + comprehension-if + base
	}

	for _, fn := range result.Functions {
		if want, ok := expected[fn.Name]; ok {
			if fn.Complexity != want {
				t.Errorf("function %s: complexity = %d, want %d", fn.Name, fn.Complexity, want)
			}
			delete(expected, fn.Name)
		}
	}
	for name := range expected {
		t.Errorf("function %s not found in parse result", name)
	}
}

func TestPythonParserFastAPIApp(t *testing.T) {
	content, err := os.ReadFile("../../testdata/fastapi_app.py")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	p := NewPythonParser()
	result, err := p.ParseFile("testdata/fastapi_app.py", content)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	// --- HTTP handlers ---
	t.Run("http_handlers", func(t *testing.T) {
		if len(result.HTTPHandlers) < 2 {
			t.Errorf("expected at least 2 HTTP handlers, got %d", len(result.HTTPHandlers))
			for _, h := range result.HTTPHandlers {
				t.Logf("  handler: %s route=%s", h.Name, h.Route)
			}
		}
	})

	// --- Call sites include pickle.loads ---
	t.Run("pickle_loads", func(t *testing.T) {
		found := false
		for _, cs := range result.CallSites {
			if cs.Name == "pickle.loads" {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected call site 'pickle.loads'")
			for _, cs := range result.CallSites {
				t.Logf("  call: %s", cs.Name)
			}
		}
	})

	// --- Struct literal: DataProcessor ---
	t.Run("data_processor_instantiation", func(t *testing.T) {
		found := false
		for _, sl := range result.StructLiterals {
			if sl.Name == "DataProcessor" {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected DataProcessor class instantiation")
			for _, sl := range result.StructLiterals {
				t.Logf("  struct literal: %s", sl.Name)
			}
		}
	})
}
