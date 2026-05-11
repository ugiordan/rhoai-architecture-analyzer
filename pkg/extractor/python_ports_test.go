package extractor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractPythonPorts(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string // relative path -> content
		wantLen  int
		wantPort int
		wantName string
	}{
		{
			name: "uvicorn.run with port",
			files: map[string]string{
				"server.py": `import uvicorn
uvicorn.run(app, host="0.0.0.0", port=8000)`,
			},
			wantLen:  1,
			wantPort: 8000,
			wantName: "uvicorn-server",
		},
		{
			name: "uvicorn.run with string app",
			files: map[string]string{
				"main.py": `uvicorn.run("app:app", port=8080)`,
			},
			wantLen:  1,
			wantPort: 8080,
			wantName: "uvicorn-server",
		},
		{
			name: "flask app.run",
			files: map[string]string{
				"app.py": `app = Flask(__name__)
app.run(host="0.0.0.0", port=5000)`,
			},
			wantLen:  1,
			wantPort: 5000,
			wantName: "app-server",
		},
		{
			name: "grpc insecure port",
			files: map[string]string{
				"grpc_server.py": `server = grpc.server(futures.ThreadPoolExecutor())
server.add_insecure_port('[::]:50051')`,
			},
			wantLen:  1,
			wantPort: 50051,
			wantName: "grpc-server",
		},
		{
			name: "grpc secure port",
			files: map[string]string{
				"grpc_server.py": `server.add_secure_port('[::]:50051', server_credentials)`,
			},
			wantLen:  1,
			wantPort: 50051,
			wantName: "grpc-server-tls",
		},
		{
			name: "gunicorn bind config",
			files: map[string]string{
				"gunicorn_conf.py": `bind = "0.0.0.0:8000"
workers = 4`,
			},
			wantLen:  1,
			wantPort: 8000,
			wantName: "gunicorn-server",
		},
		{
			name: "gunicorn bind port only",
			files: map[string]string{
				"gunicorn.conf.py": `bind = ":8080"`,
			},
			wantLen:  1,
			wantPort: 8080,
			wantName: "gunicorn-server",
		},
		{
			name: "argparse port default",
			files: map[string]string{
				"cli.py": `parser.add_argument("--port", type=int, default=8080)`,
			},
			wantLen:  1,
			wantPort: 8080,
			wantName: "cli-port-default",
		},
		{
			name: "argparse short and long port flag",
			files: map[string]string{
				"cli.py": `parser.add_argument("-p", "--port", type=int, default=8000)`,
			},
			wantLen:  1,
			wantPort: 8000,
			wantName: "cli-port-default",
		},
		{
			name: "os.environ.get PORT default",
			files: map[string]string{
				"config.py": `port = int(os.environ.get("PORT", 8000))`,
			},
			wantLen:  1,
			wantPort: 8000,
			wantName: "env-port-default",
		},
		{
			name: "os.getenv PORT string default",
			files: map[string]string{
				"config.py": `port = int(os.getenv("PORT", "8080"))`,
			},
			wantLen:  1,
			wantPort: 8080,
			wantName: "env-port-default",
		},
		{
			name: "no python files",
			files: map[string]string{
				"main.go": `package main`,
			},
			wantLen: 0,
		},
		{
			name: "test files skipped",
			files: map[string]string{
				"test_server.py": `uvicorn.run(app, port=8000)`,
			},
			wantLen: 0,
		},
		{
			name: "commented line skipped",
			files: map[string]string{
				"server.py": `# uvicorn.run(app, port=8000)`,
			},
			wantLen: 0,
		},
		{
			name: "multiple patterns in one file",
			files: map[string]string{
				"server.py": `uvicorn.run(app, port=8000)
server.add_insecure_port('[::]:50051')`,
			},
			wantLen: 2,
		},
		{
			name: "dedup same port and name",
			files: map[string]string{
				"a.py": `uvicorn.run(app, port=8000)`,
				"b.py": `uvicorn.run(app, port=8000)`,
			},
			wantLen:  1,
			wantPort: 8000,
			wantName: "uvicorn-server",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			for relPath, content := range tc.files {
				full := filepath.Join(dir, relPath)
				if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(full, []byte(content), 0644); err != nil {
					t.Fatal(err)
				}
			}

			services := extractPythonPorts(dir)
			if len(services) != tc.wantLen {
				t.Fatalf("got %d services, want %d: %+v", len(services), tc.wantLen, services)
			}

			if tc.wantLen == 0 {
				return
			}

			if tc.wantPort > 0 {
				// Check first service when expecting a single result
				svc := services[0]
				if svc.Type != "python-source" {
					t.Errorf("type = %q, want %q", svc.Type, "python-source")
				}
				if tc.wantName != "" && svc.Name != tc.wantName {
					t.Errorf("name = %q, want %q", svc.Name, tc.wantName)
				}
				if len(svc.Ports) != 1 {
					t.Fatalf("got %d ports, want 1", len(svc.Ports))
				}
				if svc.Ports[0].Port != tc.wantPort {
					t.Errorf("port = %v, want %d", svc.Ports[0].Port, tc.wantPort)
				}
				if svc.Ports[0].Protocol != "TCP" {
					t.Errorf("protocol = %q, want TCP", svc.Ports[0].Protocol)
				}
				if svc.Source == "" {
					t.Error("source is empty")
				}
			}
		})
	}
}

func TestExtractPortFromAddr(t *testing.T) {
	tests := []struct {
		addr string
		want int
	}{
		{"[::]:50051", 50051},
		{"0.0.0.0:8080", 8080},
		{":8000", 8000},
		{"localhost:9090", 9090},
		{"unix:/tmp/grpc.sock", 0},
		{"invalid", 0},
		{"", 0},
	}
	for _, tc := range tests {
		t.Run(tc.addr, func(t *testing.T) {
			got := extractPortFromAddr(tc.addr)
			if got != tc.want {
				t.Errorf("extractPortFromAddr(%q) = %d, want %d", tc.addr, got, tc.want)
			}
		})
	}
}

func TestFindPythonFiles_SkipsExcludedDirs(t *testing.T) {
	dir := t.TempDir()

	// Create files in various locations
	paths := []struct {
		rel    string
		expect bool
	}{
		{"src/server.py", true},
		{"app.py", true},
		{"venv/lib/site.py", false},
		{"__pycache__/cached.py", false},
		{"tests/test_server.py", false},
		{"test_main.py", false},
		{"src/utils_test.py", false},
		{".git/hooks/pre-commit.py", false},
	}

	for _, p := range paths {
		full := filepath.Join(dir, p.rel)
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte("# python\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	found := findPythonFiles(dir)
	foundSet := make(map[string]bool)
	for _, f := range found {
		rel, _ := filepath.Rel(dir, f)
		foundSet[rel] = true
	}

	for _, p := range paths {
		_, inSet := foundSet[p.rel]
		if p.expect && !inSet {
			t.Errorf("expected %s to be found, but it wasn't", p.rel)
		}
		if !p.expect && inSet {
			t.Errorf("expected %s to be excluded, but it was found", p.rel)
		}
	}
}
