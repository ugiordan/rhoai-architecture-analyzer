package extractor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractPythonDeps(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		wantMin  int
		validate func(t *testing.T, deps []PythonPackage)
	}{
		{
			name: "requirements.txt basic",
			files: map[string]string{
				"requirements.txt": "flask>=2.0\nrequests>=2.26.0\nnumpy\n# comment\n-r other.txt\n",
			},
			wantMin: 3,
			validate: func(t *testing.T, deps []PythonPackage) {
				found := map[string]bool{}
				for _, d := range deps {
					found[d.Name] = true
				}
				if !found["flask"] {
					t.Error("expected flask")
				}
				if !found["requests"] {
					t.Error("expected requests")
				}
				if !found["numpy"] {
					t.Error("expected numpy")
				}
			},
		},
		{
			name: "requirements dir with common.txt",
			files: map[string]string{
				"requirements/common.txt": "torch\ntransformers>=4.56.0\nfastapi>=0.120.1\n",
				"requirements/test.txt":   "pytest\ncoverage\n",
			},
			wantMin: 5,
			validate: func(t *testing.T, deps []PythonPackage) {
				for _, d := range deps {
					if d.Name == "torch" {
						if d.Category != "ml-framework" {
							t.Errorf("torch category: got %q, want ml-framework", d.Category)
						}
						if !d.Required {
							t.Error("torch should be required")
						}
					}
					if d.Name == "pytest" {
						if d.Required {
							t.Error("pytest should not be required (test dep)")
						}
					}
				}
			},
		},
		{
			name: "pyproject.toml",
			files: map[string]string{
				"pyproject.toml": `[project]
name = "mypackage"
version = "1.0"
dependencies = [
    "fastapi>=0.100.0",
    "uvicorn",
    "pydantic>=2.0",
]
`,
			},
			wantMin: 3,
			validate: func(t *testing.T, deps []PythonPackage) {
				found := map[string]bool{}
				for _, d := range deps {
					found[d.Name] = true
				}
				if !found["fastapi"] {
					t.Error("expected fastapi")
				}
				if !found["pydantic"] {
					t.Error("expected pydantic")
				}
			},
		},
		{
			name: "version constraint parsing",
			files: map[string]string{
				"requirements.txt": "openai>=1.99.1, <2.25.0\nprotobuf>=5.29.6, !=6.30.*\ngrpcio>=1.76.0\n",
			},
			wantMin: 3,
			validate: func(t *testing.T, deps []PythonPackage) {
				for _, d := range deps {
					if d.Name == "openai" {
						if d.Version == "" {
							t.Error("openai should have version constraint")
						}
						if d.Category != "llm-sdk" {
							t.Errorf("openai category: got %q, want llm-sdk", d.Category)
						}
					}
					if d.Name == "grpcio" {
						if d.Category != "grpc" {
							t.Errorf("grpcio category: got %q, want grpc", d.Category)
						}
					}
				}
			},
		},
		{
			name: "extras and markers stripped",
			files: map[string]string{
				"requirements.txt": "fastapi[standard]>=0.120.1\nnumba==0.61.2; platform_machine==\"x86_64\"\n",
			},
			wantMin: 2,
			validate: func(t *testing.T, deps []PythonPackage) {
				found := map[string]bool{}
				for _, d := range deps {
					found[d.Name] = true
				}
				if !found["fastapi"] {
					t.Error("expected fastapi (extras should be stripped)")
				}
				if !found["numba"] {
					t.Error("expected numba (markers should be stripped)")
				}
			},
		},
		{
			name: "deduplication across files",
			files: map[string]string{
				"requirements.txt":        "numpy>=1.20\n",
				"requirements/common.txt": "numpy>=1.24\n",
			},
			wantMin: 1,
			validate: func(t *testing.T, deps []PythonPackage) {
				count := 0
				for _, d := range deps {
					if d.Name == "numpy" {
						count++
					}
				}
				if count != 1 {
					t.Errorf("expected 1 numpy entry, got %d", count)
				}
			},
		},
		{
			name: "setup.py install_requires",
			files: map[string]string{
				"setup.py": `from setuptools import setup

setup(
    name="mypackage",
    install_requires=[
        "torch>=2.0",
        "fastapi",
        'uvicorn>=0.30',
    ],
)
`,
			},
			wantMin: 3,
			validate: func(t *testing.T, deps []PythonPackage) {
				found := map[string]bool{}
				for _, d := range deps {
					found[d.Name] = true
					if !d.Required {
						t.Errorf("%s should be required from setup.py", d.Name)
					}
				}
				if !found["torch"] {
					t.Error("expected torch")
				}
				if !found["uvicorn"] {
					t.Error("expected uvicorn")
				}
			},
		},
		{
			name: "setup.cfg install_requires",
			files: map[string]string{
				"setup.cfg": `[metadata]
name = mypackage
version = 1.0

[options]
install_requires =
    redis>=4.0
    boto3
    openai>=1.0

[options.extras_require]
test =
    pytest
`,
			},
			wantMin: 3,
			validate: func(t *testing.T, deps []PythonPackage) {
				found := map[string]bool{}
				for _, d := range deps {
					found[d.Name] = true
				}
				if !found["redis"] {
					t.Error("expected redis")
				}
				if !found["boto3"] {
					t.Error("expected boto3")
				}
				if !found["openai"] {
					t.Error("expected openai")
				}
			},
		},
		{
			name:  "empty repo no deps",
			files: map[string]string{},
			validate: func(t *testing.T, deps []PythonPackage) {
				if len(deps) != 0 {
					t.Errorf("expected 0 deps for empty repo, got %d", len(deps))
				}
			},
		},
		{
			name: "malformed and edge case lines",
			files: map[string]string{
				"requirements.txt": `# pure comment line
-r other.txt
--index-url https://pypi.org/simple

===invalid===
torch>=2.0
`,
			},
			wantMin: 1,
			validate: func(t *testing.T, deps []PythonPackage) {
				if len(deps) != 1 {
					t.Errorf("expected exactly 1 dep, got %d: %+v", len(deps), deps)
				}
				if deps[0].Name != "torch" {
					t.Errorf("expected torch, got %s", deps[0].Name)
				}
			},
		},
		{
			name: "name normalization underscores dots mixed case",
			files: map[string]string{
				"requirements.txt": "Pillow>=9.0\nscikit_learn>=1.0\nruamel.yaml>=0.17\n",
			},
			wantMin: 3,
			validate: func(t *testing.T, deps []PythonPackage) {
				found := map[string]bool{}
				for _, d := range deps {
					found[d.Name] = true
				}
				if !found["pillow"] {
					t.Error("expected 'pillow' (lowercased from Pillow)")
				}
				if !found["scikit-learn"] {
					t.Error("expected 'scikit-learn' (underscore normalized)")
				}
				if !found["ruamel-yaml"] {
					t.Error("expected 'ruamel-yaml' (dot normalized)")
				}
			},
		},
		{
			name: "merge prefers version over no version",
			files: map[string]string{
				"requirements.txt":        "numpy\n",
				"requirements/common.txt": "numpy>=1.24\n",
			},
			wantMin: 1,
			validate: func(t *testing.T, deps []PythonPackage) {
				for _, d := range deps {
					if d.Name == "numpy" {
						if d.Version == "" {
							t.Error("expected numpy to have version from merge")
						}
						return
					}
				}
				t.Error("numpy not found")
			},
		},
		{
			name: "merge prefers required over not required",
			files: map[string]string{
				"requirements-test.txt": "requests>=2.0\n",
				"requirements.txt":      "requests>=2.0\n",
			},
			wantMin: 1,
			validate: func(t *testing.T, deps []PythonPackage) {
				for _, d := range deps {
					if d.Name == "requests" {
						if !d.Required {
							t.Error("expected requests to be required after merge")
						}
						return
					}
				}
				t.Error("requests not found")
			},
		},
		{
			name: "pyproject.toml poetry format",
			files: map[string]string{
				"pyproject.toml": `[tool.poetry]
name = "mypackage"
version = "1.0.0"

[tool.poetry.dependencies]
python = "^3.10"
grpcio = ">=1.60"
prometheus-client = "^0.20"

[tool.poetry.dev-dependencies]
pytest = "^7.0"
`,
			},
			wantMin: 2,
			validate: func(t *testing.T, deps []PythonPackage) {
				found := map[string]bool{}
				for _, d := range deps {
					found[d.Name] = true
				}
				if !found["grpcio"] {
					t.Error("expected grpcio from poetry deps")
				}
			},
		},
		{
			name: "inline comments stripped",
			files: map[string]string{
				"requirements.txt": "torch>=2.0 # GPU required\nfastapi # web framework\n",
			},
			wantMin: 2,
			validate: func(t *testing.T, deps []PythonPackage) {
				for _, d := range deps {
					if d.Name == "torch" && d.Version == "" {
						t.Error("torch version should survive inline comment stripping")
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()

			for path, content := range tt.files {
				fullPath := filepath.Join(dir, path)
				if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
					t.Fatal(err)
				}
			}

			deps := extractPythonDeps(dir)
			if len(deps) < tt.wantMin {
				t.Errorf("expected at least %d deps, got %d: %+v", tt.wantMin, len(deps), deps)
			}
			if tt.validate != nil {
				tt.validate(t, deps)
			}
		})
	}
}
