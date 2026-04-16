# Installation

## Requirements

- Go 1.25+
- git

## Build from source

```bash
git clone https://github.com/ugiordan/rhoai-architecture-analyzer.git
cd rhoai-architecture-analyzer
go build -o rhoai-analyzer ./cmd/rhoai-analyzer/
```

The binary `rhoai-analyzer` is now ready to use.

## Verify installation

```bash
./rhoai-analyzer version
```

Expected output:

```
rhoai-architecture-analyzer v0.2.0
```

## Optional: Add to PATH

```bash
# Move to a directory in your PATH
sudo mv rhoai-analyzer /usr/local/bin/

# Or add the project directory to PATH
export PATH="$PATH:$(pwd)"
```

## Tree-sitter dependency

The code property graph and security scanning features use tree-sitter for Go source parsing. Tree-sitter is included as a Go dependency and compiled automatically during `go build`. No additional installation needed.

## GitHub Actions

For CI/CD usage, the analyzer builds in the workflow and runs against repos. See [CI Integration](../guides/ci-integration.md) for workflow examples.
