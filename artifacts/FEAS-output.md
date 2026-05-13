===REVIEW_TARGET_356d451dea5b9d195a00e13ecc0cf05a_START===
# FEAS Review: Go AST Extraction Design

## Findings

```
Finding ID: FEAS-001
Specialist: Feasibility Analyst
Severity: Important
Confidence: High
Document: 2026-05-12-go-ast-extraction-design.md
Citation: Section "Decisions" table row "CI strategy" and Section "7. CI Changes"
Title: go mod download adds unquantified CI cost across 50+ repos, some of which are not Go repos or have private/replace dependencies
Evidence: The design states "Always run go mod download for Go repos" and "Verified on 12+ repos, adds ~30-60s per job, parallel execution." The CI pipeline runs against 50+ repos (Section "Prior Art" mentions "we correlate across 50+ repos"). The 30-60s estimate was verified on only 12 repos. Repos with replace directives pointing to private modules, large vendored dependencies, or complex dependency trees (e.g., operator repos pulling in controller-runtime's full transitive closure) can take well beyond 60 seconds. The 2-minute timeout mitigates worst-case hangs but a timeout failure still means the CI job spent 2 minutes producing nothing beyond what the fallback would have provided. Additionally, the design does not address whether Go module cache persistence is available in CI (GitHub Actions cache). Without module caching, every CI run re-downloads all dependencies from scratch, compounding the time cost.
Recommended fix: Add explicit guidance on Go module caching in CI (e.g., actions/cache for GOMODCACHE). Clarify whether the 30-60s estimate includes the full transitive dependency download or just go mod download with warm cache. Consider running go mod download only when YAML CRD/webhook files are absent (the actual problem being solved) rather than unconditionally for all Go repos.
Verdict: Revise
```

```
Finding ID: FEAS-002
Specialist: Feasibility Analyst
Severity: Important
Confidence: High
Document: 2026-05-12-go-ast-extraction-design.md
Citation: Section "1. Go Package Loading", packages.Load config block
Title: go/packages full type loading (NeedTypes + NeedTypesInfo + NeedDeps) requires successful compilation, which many operator repos may not achieve in CI without build tags and CGO setup
Evidence: The design specifies loading with packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedDeps | packages.NeedFiles | packages.NeedName | packages.NeedImports | packages.NeedCompiledGoFiles on "./...". The NeedTypes and NeedTypesInfo modes require the Go type checker to run, which in turn requires all imported packages to be available and type-checkable. Operator repos commonly import C-dependent packages (e.g., prometheus client with CGO, certain crypto libraries), packages with build tags (e.g., "//go:build linux"), or generated code (protobuf, deepcopy-gen output). On macOS CI runners analyzing linux-targeted operator code, or on any runner without the correct build constraints, packages.Load with these modes will produce type-checking errors for those packages. While the design states "On any error: log warning, return GoPackageSet{Mode: 'fallback'}", the packages.Load API does not return a clean error for partial type failures. It returns packages with Errors fields populated. The design does not specify how to distinguish "partial type errors in non-critical packages" from "total loading failure", meaning the fallback may trigger too aggressively (losing all Go AST data) or not aggressively enough (producing incorrect type resolutions).
Recommended fix: Specify a strategy for handling partial packages.Load errors. Consider using packages.NeedSyntax | packages.NeedName | packages.NeedFiles | packages.NeedImports (without NeedTypes/NeedTypesInfo) as a middle ground that still provides AST and import resolution without requiring full type checking. Only escalate to NeedTypes for the specific packages where webhook/CRD types live (the api/ and internal/ subtrees). Document the error-handling heuristic explicitly, e.g., "if >50% of packages have errors, fall back; otherwise use the successfully loaded packages."
Verdict: Revise
```

```
Finding ID: FEAS-003
Specialist: Feasibility Analyst
Severity: Important
Confidence: High
Document: 2026-05-12-go-ast-extraction-design.md
Citation: Section "Architecture" > "New Files" table and "Modified Files" table, Risks table row "Large golang.org/x/tools dependency"
Title: golang.org/x/tools dependency is not currently in go.mod and will significantly increase the dependency tree
Evidence: The design states in the Risks table: "Large golang.org/x/tools dependency - Already common in Go tooling, well-maintained." However, the current go.mod at the project root shows no golang.org/x/tools dependency. The project currently depends on go-tree-sitter, gopkg.in/yaml.v3, sigs.k8s.io/kustomize, and sigs.k8s.io/yaml as direct dependencies. Adding golang.org/x/tools/go/packages will pull in golang.org/x/tools and its transitive dependencies (including golang.org/x/mod, additional golang.org/x/ packages). This is a material increase in the project's dependency surface. The design's claim that it's "already common in Go tooling" is about the ecosystem, not this specific project. The architecture-analyzer binary is distributed and built in CI. A larger dependency tree increases build time, binary size, and supply chain attack surface.
Recommended fix: Acknowledge the dependency addition explicitly (not just as a risk dismissal). Quantify the expected binary size increase. Consider whether go/packages loading could be implemented as a separate binary or plugin to keep the core analyzer lean for repos where Go AST extraction is not needed.
Verdict: Revise
```

```
Finding ID: FEAS-004
Specialist: Feasibility Analyst
Severity: Important
Confidence: Medium
Document: 2026-05-12-go-ast-extraction-design.md
Citation: Section "3. Webhook Behavioral Analysis", Step 2 and Step 3
Title: Webhook behavioral analysis depth (one-level method following, field assignment tracking) may produce incomplete results for real-world webhook handlers that use helper packages, interfaces, or table-driven patterns
Evidence: The design states "Follow one level of method calls on the same receiver: if Default() calls r.setGPUDefaults(), analyze that method body too." Real-world webhook handlers in the Kubernetes ecosystem frequently delegate to utility functions in separate packages (e.g., a defaults package or webhookutil package), use interface-based dispatch (e.g., a Defaulter interface chain), or use table-driven patterns where defaults are stored in a map/slice and applied in a loop. The one-level, same-receiver limitation means: (1) helper functions in other packages are missed, (2) defaults applied through iteration over a config struct are missed, (3) validation logic delegated to a shared validation package is missed. The design acknowledges "Level 2 only (field ops)" but does not quantify what percentage of real webhook handlers would produce meaningful output vs. empty/minimal results. If the hit rate is low, the feature's value proposition weakens.
Recommended fix: Before implementation, audit 5-10 actual webhook handlers from opendatahub-operator and other target repos to measure what percentage of field mutations and validations would be captured by the one-level, same-receiver heuristic. If the hit rate is below ~60%, consider expanding to cross-package method following (at least for methods on the same type defined in other files of the same package) or documenting the expected coverage gap in the output.
Verdict: Revise
```

```
Finding ID: FEAS-005
Specialist: Feasibility Analyst
Severity: Minor
Confidence: High
Document: 2026-05-12-go-ast-extraction-design.md
Citation: Section "Architecture" > "Pipeline Integration", specifically the fallback path "existing go/parser extractors (unchanged)"
Title: Design conflates tree-sitter Go parser (pkg/parser/go_parser.go) with stdlib go/parser; fallback path description is imprecise
Evidence: The design states the fallback path uses "existing go/parser extractors (unchanged)". In the codebase, there are two distinct Go parsing systems: (1) The CPG builder in pkg/parser/go_parser.go uses go-tree-sitter (github.com/smacker/go-tree-sitter), not go/parser, and (2) several extractors in pkg/extractor/ (api_types.go, platform_detection.go, reconcile_sequence.go, operator_config.go, status_conditions.go) use stdlib go/parser. The design's fallback description does not distinguish between these two systems. Since the new Go AST extractors (go_crds.go, go_webhooks.go) are in pkg/extractor/, the fallback would mean the existing pkg/extractor/ go/parser-based extractors continue to run (which is correct), but the phrasing could mislead implementers about what "existing go/parser extractors" covers.
Recommended fix: Clarify in the Pipeline Integration section that the fallback means "the existing stdlib go/parser-based extractors in pkg/extractor/ continue to run for their respective domains (api_types, platform_detection, etc.), while the new go/packages-dependent extractors (CRD-from-Go, webhook behavior, enhanced controller watches) simply do not run." This avoids confusion with the tree-sitter-based CPG parser in pkg/parser/.
Verdict: Approve
```

```
Finding ID: FEAS-006
Specialist: Feasibility Analyst
Severity: Minor
Confidence: Medium
Document: 2026-05-12-go-ast-extraction-design.md
Citation: Section "2. CRD Extraction from Go", step 2, bullet "Group: find GroupVersion var in same package's groupversion_info.go"
Title: GroupVersion variable location assumption is fragile across operator codebases
Evidence: The design states: "Group: find GroupVersion var in same package's groupversion_info.go, extract Group field." While this is the kubebuilder scaffold convention, not all operator repos follow this pattern. Some repos define GroupVersion in a different file name (e.g., register.go, types.go, doc.go), use a const block instead of a var, or define the group string inline in the SchemeBuilder registration. The design assumes a single file naming convention without specifying a fallback search strategy.
Recommended fix: Search for GroupVersion (or SchemeGroupVersion) variable/constant declarations across all Go files in the same package, not just groupversion_info.go. This is a minor robustness improvement that aligns with the design's general philosophy of handling real-world code variations.
Verdict: Approve
```

## Overall Verdict

```
OVERALL_VERDICT: REVISE
Justification: Four findings require revision (FEAS-001 through FEAS-004) covering CI cost modeling, partial type-loading error handling, dependency impact acknowledgment, and webhook analysis coverage validation. None are individually blocking, but together they represent underspecified feasibility dimensions that should be addressed before implementation begins.
```
===REVIEW_TARGET_356d451dea5b9d195a00e13ecc0cf05a_END===
