# Code Property Graph

The code property graph (CPG) is a unified representation of Go source code that supports cross-function analysis and security queries.

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

```mermaid
flowchart LR
    SRC["*.go files"] --> TS["Tree-sitter\nParser"]
    TS --> PR["ParseResult\n(functions, calls,\nparams, literals)"]
    PR --> BUILD["Builder"]
    BUILD --> CPG["Code Property\nGraph"]

    classDef parse fill:#3498db,stroke:#2980b9,color:#fff
    classDef build fill:#2ecc71,stroke:#27ae60,color:#fff
    class TS,PR parse
    class BUILD build
```

1. **Parser** (`pkg/parser/go_parser.go`): Tree-sitter parses each Go file, extracting:
    - Function declarations with parameters and return types
    - Function call expressions with arguments
    - Composite literals (struct instantiation)
    - Switch/case statements

2. **Builder** (`pkg/builder/builder.go`): Assembles parse results into the CPG:
    - Creates nodes for each function, parameter, call, literal
    - Creates edges (calls, contains, aliases)
    - Resolves cross-file references

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
