# Code Property Graph

The code property graph (CPG) is a unified representation of Go source code that supports cross-function analysis and security queries.

## What is a CPG?

When you write Go code like `func main() { fmt.Println("hello") }`, the compiler first turns it into an **Abstract Syntax Tree (AST)**: a tree structure where each node represents a syntactic element (function declaration, call expression, string literal, etc.). The AST captures the structure of the code but not the relationships between different parts of the program.

A **Code Property Graph** goes further: it takes the AST and adds cross-references, data flow edges, and semantic annotations. This means you can ask questions like "which functions call `exec.Command` with user-supplied input?" or "does any controller watch a CRD it doesn't own?" by traversing the graph rather than doing text search.

The analyzer builds its CPG using [tree-sitter](https://tree-sitter.github.io/tree-sitter/), a fast incremental parsing library, via the Go binding [`go-tree-sitter`](https://github.com/smacker/go-tree-sitter). Tree-sitter produces a concrete syntax tree for each Go source file, which the builder then transforms into the CPG by resolving cross-file references and adding semantic edges.

## Structure

```mermaid
graph TD
    FILE["File Node"] --> FUNC["Function Node"]
    FUNC --> PARAM["Parameter Nodes"]
    FUNC --> CALL["Call Nodes"]
    FUNC --> STRUCT["Struct Literal Nodes"]

    FUNC -->|"EdgeCalls"| FUNC2["Called Function"]
    PARAM -->|"EdgeDataFlow"| CALL
    STRUCT -->|"EdgeContains"| FUNC

    classDef node fill:#9b59b6,stroke:#8e44ad,color:#fff
    class FILE,FUNC,PARAM,CALL,STRUCT,FUNC2 node
```

### Node kinds

| Kind | Description |
|------|-------------|
| `File` | Source file |
| `Function` | Function or method declaration |
| `Parameter` | Function parameter with type |
| `Call` | Function call expression |
| `StructLiteral` | Composite literal (struct instantiation) |

### Edge kinds

| Kind | Description |
|------|-------------|
| `EdgeCalls` | Function A calls function B |
| `EdgeAliases` | Type alias relationship |
| `EdgeContains` | Containment (file contains function, function contains literal) |

### Properties

Each node carries:

- **Name**: Identifier
- **File**: Source file path
- **Line**: Line number
- **Properties**: Key-value metadata (e.g., parameter type, return type)
- **Annotations**: Domain-specific metadata added by annotators

## Thread safety

The CPG implementation (`pkg/graph/cpg.go`) is thread-safe:

- `sync.RWMutex` protects all node and edge operations
- Multiple annotators can read concurrently
- Write operations (adding nodes/edges) are serialized

## Building the CPG

The CPG is built in two stages: parse, then assemble.

```mermaid
flowchart LR
    SRC["*.go files"] --> TS["Tree-sitter\nParser"]
    TS --> AST["Syntax Trees\n(per file)"]
    AST --> EXT["Extractor\n(functions, calls,\nparams, literals)"]
    EXT --> PR["ParseResult"]
    PR --> BUILD["Builder\n(cross-file\nresolution)"]
    BUILD --> CPG["Code Property\nGraph"]

    classDef parse fill:#3498db,stroke:#2980b9,color:#fff
    classDef build fill:#2ecc71,stroke:#27ae60,color:#fff
    class TS,AST,EXT,PR parse
    class BUILD build
```

1. **Parser** (`pkg/parser/go_parser.go`): Tree-sitter parses each `.go` file into a syntax tree, then the parser walks the tree extracting:
    - Function declarations with parameters and return types
    - Function call expressions with arguments (including detecting sensitive sinks like `exec.Command`, `sql.Query`)
    - Composite literals (struct instantiation with field values)
    - Switch/case statements (used for detecting unhandled CRD versions)

2. **Builder** (`pkg/builder/builder.go`): Assembles per-file parse results into the unified CPG:
    - Creates nodes for each function, parameter, call, literal
    - Creates edges (calls, contains, aliases)
    - Resolves cross-file references (function A in file X calls function B in file Y)

## Architecture enrichment

When `--with-arch` is provided, the CPG gains an `ArchData` sidecar:

```go
type CPG struct {
    nodes    map[string]*Node
    edges    map[string][]*Edge
    ArchData *arch.ArchitectureData  // Optional
}
```

Architecture data enables queries that cross-reference code against extracted architecture:

- CGA-U01: Compare CRD version references in code against extracted CRD schemas
- Architecture-aware taint analysis: Follow data through known API boundaries
- Finding enrichment: Add `ArchRef` to findings linking code to architecture components

## Query execution

Queries traverse the CPG looking for patterns:

```mermaid
flowchart LR
    QUERY["Query\n(e.g., CGA-004)"] --> WALK["Walk nodes\nmatching criteria"]
    WALK --> CHECK["Check annotations\nand properties"]
    CHECK --> EMIT["Emit finding\nwith file:line"]

    classDef query fill:#e74c3c,stroke:#c0392b,color:#fff
    class QUERY,WALK,CHECK query
```

The query engine (`pkg/query/engine.go`) provides:

- `RunDomain(domain)`: Execute all queries for a domain
- Results grouped by domain with deduplication
- Per-domain SARIF output support

## Taint analysis

Taint analysis (`pkg/query/taint.go`) traces data flow from sources to sinks:

1. **Sources**: Functions that introduce untrusted data (HTTP request params, env vars, user input)
2. **Propagators**: Functions that pass data through (assignments, returns, function calls)
3. **Sinks**: Functions that perform sensitive operations (SQL queries, command execution, file writes)

A taint finding is emitted when data flows from a source to a sink without passing through a sanitizer.
