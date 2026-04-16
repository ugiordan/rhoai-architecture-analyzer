# Adding Extractors

This guide covers adding a new architecture extractor to the analyzer.

## Extractor types

There are three categories:

1. **YAML extractors**: Parse Kubernetes manifests (most extractors)
2. **Go source extractors**: Parse Go code with tree-sitter
3. **File extractors**: Parse specific file formats (Dockerfiles, Helm, go.mod)

## Step 1: Define the data type

Add your extracted data type to `pkg/extractor/types.go`:

```go
type MyResource struct {
    Name       string   `json:"name"`
    Namespace  string   `json:"namespace,omitempty"`
    Properties []string `json:"properties"`
    SourceFile string   `json:"source_file"`
}
```

Add the field to `ComponentArchitecture`:

```go
type ComponentArchitecture struct {
    // ... existing fields
    MyResources []MyResource `json:"my_resources,omitempty"`
}
```

## Step 2: Write the extractor

Create `pkg/extractor/my_resource.go`:

```go
package extractor

import (
    "path/filepath"
)

func extractMyResources(repoPath string) ([]MyResource, error) {
    var results []MyResource

    // Find relevant files
    files, err := findYAMLFiles(repoPath, "**/my-resource*.yaml")
    if err != nil {
        return nil, err
    }

    for _, file := range files {
        // Parse YAML
        resources, err := parseMyResourceYAML(file)
        if err != nil {
            // Log warning, don't fail
            continue
        }

        relPath, _ := filepath.Rel(repoPath, file)
        for _, r := range resources {
            r.SourceFile = relPath
            results = append(results, r)
        }
    }

    return results, nil
}
```

## Step 3: Register in ExtractAll

Add your extractor to `pkg/extractor/extract.go`:

```go
func ExtractAll(repoPath string) (*ComponentArchitecture, error) {
    arch := &ComponentArchitecture{}

    // ... existing extractors

    // My resources
    myResources, err := extractMyResources(repoPath)
    if err != nil {
        log.Printf("WARNING: my-resource extraction failed: %v", err)
    }
    arch.MyResources = myResources

    return arch, nil
}
```

!!! note "Resilience"
    Extractors should log warnings on failure, not return errors. A failed extractor should not prevent other extractors from running.

## Step 4: Add a renderer (optional)

If your data needs visualization, create `pkg/renderer/my_resource.go`:

```go
func renderMyResources(arch *ComponentArchitecture) string {
    // Build Mermaid diagram, markdown table, etc.
}
```

Register in `pkg/renderer/renderer.go`.

## Step 5: Add tests

Create `pkg/extractor/my_resource_test.go` with test fixtures in `testdata/`.

Test at minimum:

- Valid YAML parsing
- Missing/empty files (should not error)
- Malformed YAML (should skip, not crash)
- Relative path calculation

## Guidelines

- **Read-only**: Never modify files in the repository
- **Resilient**: Skip bad files, don't crash
- **Source-traceable**: Every extracted fact must include `source_file`
- **No secrets**: Never extract secret values (only names and references)
- **Minimal dependencies**: Prefer parsing YAML/Go directly over importing heavy libraries

## YAML file discovery

Use the existing `findYAMLFiles` helper in `pkg/extractor/yaml.go`:

```go
files, err := findYAMLFiles(repoPath, "config/my-resource/**/*.yaml")
```

Supports glob patterns and walks the repository filesystem.

## Go source extraction

For extracting information from Go source code, use the tree-sitter parser in `pkg/parser/`:

```go
parser := parser.NewGoParser()
result, err := parser.ParseFile(filePath)
// result.Functions, result.Calls, result.StructLiterals, etc.
```

See the controller watches extractor (`controller_watches.go`) for a complete example of Go source extraction.
