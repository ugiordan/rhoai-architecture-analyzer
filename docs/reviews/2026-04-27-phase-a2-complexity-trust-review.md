# Phase A2: Complexity & Trust Classification Review

## Executive Summary

**Review Date:** 2026-04-27
**Specialists:** SEC (Security Auditor), PERF (Performance Analyst), QUAL (Code Quality Reviewer), CORR (Correctness Verifier), ARCH (Architecture Reviewer)
**Files Reviewed:** 19 files (803 lines added, 4 removed)
**Agreement Level:** Strong Agreement (6/6 findings with majority or consensus, 0 escalated, 0 dismissed)
**Configuration:** 1 iteration (convergence achieved via cross-specialist agreement), budget ~120K/350K consumed (34%)

| Severity | Count |
|----------|-------|
| Critical | 0     |
| Important| 3     |
| Minor    | 3     |
| **Total**| **6** |

## Review Configuration

- **Date:** 2026-04-27
- **Scope:** `pkg/parser/`, `pkg/domains/security/`, `testdata/` (19 files, ~803 lines changed)
- **Specialists:** SEC, PERF, QUAL, CORR, ARCH
- **Mode flags:** `--diff --range 7796f12..HEAD`
- **Iterations:** All specialists: 1 (strong cross-specialist convergence on iteration 1)
- **Budget:** ~120K / 350K consumed (34%)
- **Reference modules:** 0 loaded
- **Constraints:** none

## Consensus Findings

All specialists agree on the following findings.

### PHASE-A2-001: Cross-language complexity operator inconsistency

- **Severity:** Important
- **Confidence:** High
- **File:** `pkg/parser/go_parser.go` (lines 59-107), `pkg/parser/rust_parser.go` (lines 68-107), `pkg/parser/python_parser.go` (lines 82-113), `pkg/parser/typescript_parser.go` (lines 76-138)
- **Evidence:** Go and Rust count only `&&` as a complexity decision point in binary expressions. Python counts both `and` and `or` (via `boolean_operator` node type). TypeScript counts `&&`, `||`, and `??`. This means identical control flow patterns produce different complexity scores across languages. For example, `if a || b { ... }` yields complexity 1 in Go/Rust but 2 in TypeScript. This undermines cross-language comparisons and any threshold-based queries (like CGA-010).
- **Recommended fix:** Document the operator counting policy explicitly. Either: (a) count only short-circuit operators (`&&`/`and`) consistently across all languages, or (b) count all boolean operators (`&&`, `||`, `and`, `or`, `??`) consistently. Option (a) is simpler and matches the "decision point" definition more closely. Whichever is chosen, update all 4 parsers to match.
- **Agreement:** Unanimous (4/5 specialists raised independently: SEC, QUAL, CORR, ARCH; PERF concurred)

### PHASE-A2-002: Go security annotator missing language filter on function trust classification

- **Severity:** Important
- **Confidence:** High
- **File:** `pkg/domains/security/go_annotator.go` (lines 38-58)
- **Evidence:** `classifyTrust` filters HTTP endpoints by language (line 33: `if ep.Language != "" && ep.Language != "go" { continue }`) but the function trust classification loop at line 38 has no such filter. If a Python function named `Reconcile` with a parameter type containing `ctrl.Request` appears in the graph, the Go annotator will incorrectly set it to `TrustTrusted`. Similarly, any non-Go function named `init` or `main` gets marked trusted.
- **Recommended fix:** Add the same language guard to the function loop:
  ```go
  for _, fn := range g.NodesByKind(graph.NodeFunction) {
      if fn.Language != "" && fn.Language != "go" {
          continue
      }
      // ... existing trust logic
  }
  ```
- **Agreement:** Unanimous (5/5 specialists raised independently)

### PHASE-A2-003: CGA-011 query is tautological

- **Severity:** Important
- **Confidence:** High
- **File:** `pkg/domains/security/queries.go` (lines for `queryUntrustedWithoutValidation`)
- **Evidence:** CGA-011 ("Untrusted endpoint without input validation") fires on every node with `TrustLevel == TrustUntrusted`. There is no check for whether the endpoint actually lacks validation. The query name says "without validation" but the implementation checks only trust level, not validation presence. Every untrusted HTTP endpoint triggers this finding, making it a noise generator rather than a meaningful security query.
- **Recommended fix:** Either: (a) rename to "Untrusted endpoint" (removing the "without validation" claim) and accept it as a low-priority informational query, or (b) add actual validation detection: check if the function body contains calls to validation libraries, input sanitization functions, or schema validators, and only fire when those are absent. Option (b) requires parser-level support for tracing validation calls, which may be future work. Option (a) is immediate.
- **Agreement:** Unanimous (4/5 specialists raised independently: SEC, QUAL, CORR, ARCH; PERF concurred)

## Majority Findings

The following findings achieved majority agreement.

### PHASE-A2-004: Unused `src []byte` parameter in Go complexity functions

- **Severity:** Minor
- **Confidence:** High
- **File:** `pkg/parser/go_parser.go` (lines 53-108)
- **Evidence:** `computeComplexity(node *sitter.Node, src []byte)` and `countDecisionPoints(node *sitter.Node, src []byte, count *int)` accept a `src []byte` parameter that is never read. The parameter is only passed through recursively. No other language's complexity function requires source bytes.
- **Recommended fix:** Remove the `src` parameter from both functions and update the single call site in `extractFunction`.
- **Agreement:** 3/5 specialists (PERF, QUAL, ARCH)

### PHASE-A2-005: CGA-010 hardcoded annotation list

- **Severity:** Minor
- **Confidence:** Medium
- **File:** `pkg/domains/security/queries.go` (lines for `queryComplexityHotspot`)
- **Evidence:** CGA-010 checks complexity > 10 against 9 hardcoded security annotation keys. Adding new security annotations (e.g., from new annotators for Rust or TypeScript) requires manually updating this list. The annotation check uses `hasAnyAnnotation` with a string slice, so missing annotations silently reduce recall.
- **Recommended fix:** Consider using a registry pattern or iterating all annotations with a `sec:` prefix, so new annotations are automatically included. Alternatively, document the annotation list explicitly so future annotator authors know to update CGA-010.
- **Agreement:** 2/5 specialists (SEC, QUAL); others noted but assessed as low priority

### PHASE-A2-006: Python annotator O(E*F) nested loop in trust classification

- **Severity:** Minor
- **Confidence:** Medium
- **File:** `pkg/domains/security/python_annotator.go` (trust classification section)
- **Evidence:** The Python annotator iterates over all HTTP endpoints and then over all functions to match auth decorators, producing O(endpoints * functions) complexity. For typical project sizes this is negligible, but the pattern differs from the Go annotator's single-pass approach.
- **Recommended fix:** Low priority. For consistency, consider building an index of decorated functions first, then iterating endpoints. Not worth fixing unless profiling shows it matters.
- **Agreement:** 2/5 specialists (PERF, ARCH); others assessed as negligible for real-world sizes

## Escalated Disagreements

None.

## Dismissed Findings

None.

## Co-located Findings

### Co-location Group: `pkg/domains/security/queries.go`

| Finding ID | Specialist | Severity | Title |
|------------|-----------|----------|-------|
| PHASE-A2-003 | SEC, QUAL, CORR, ARCH | Important | CGA-011 tautological query |
| PHASE-A2-005 | SEC, QUAL | Minor | CGA-010 hardcoded annotation list |

**Interaction notes:** Both findings target the security queries file. CGA-011 and CGA-010 were added together in this phase. Fixing CGA-011 (either renaming or adding validation detection) should be done before extending CGA-010's annotation list, since the query semantics may shift.

### Co-location Group: `pkg/parser/*_parser.go` (all 4 parsers)

| Finding ID | Specialist | Severity | Title |
|------------|-----------|----------|-------|
| PHASE-A2-001 | SEC, QUAL, CORR, ARCH | Important | Cross-language operator inconsistency |
| PHASE-A2-004 | PERF, QUAL, ARCH | Minor | Unused `src` parameter (Go only) |

**Interaction notes:** Both findings target the parser complexity computation. PHASE-A2-001 requires touching all 4 parsers to align operator counting. PHASE-A2-004 is Go-only and can be fixed independently.

## Remediation Summary

### All Findings by Severity

| ID | Severity | Area | File | Title |
|----|----------|------|------|-------|
| PHASE-A2-001 | Important | Parser | `pkg/parser/*_parser.go` | Cross-language complexity operator inconsistency |
| PHASE-A2-002 | Important | Security | `pkg/domains/security/go_annotator.go:38-58` | Missing language filter on function trust |
| PHASE-A2-003 | Important | Security | `pkg/domains/security/queries.go` | CGA-011 tautological query |
| PHASE-A2-004 | Minor | Parser | `pkg/parser/go_parser.go:53-108` | Unused `src` parameter |
| PHASE-A2-005 | Minor | Security | `pkg/domains/security/queries.go` | CGA-010 hardcoded annotation list |
| PHASE-A2-006 | Minor | Security | `pkg/domains/security/python_annotator.go` | O(E*F) nested loop |

### Remediation Roadmap

| Category | Count | Description |
|----------|-------|-------------|
| Actionable (Chore) | 4 | Self-contained fixes for direct PR (001, 002, 003, 004) |
| Blocked/Deferred | 1 | CGA-010 annotation registry depends on annotator design decisions (005) |
| Low Priority | 1 | Python loop optimization, not worth fixing unless profiling shows need (006) |

### Top Priorities

1. **PHASE-A2-002** - Go annotator sets trust on non-Go functions. One-line fix, highest confidence, highest impact.
2. **PHASE-A2-001** - Complexity operator inconsistency undermines cross-language comparisons. Needs a design decision (count `||` or not), then straightforward implementation.
3. **PHASE-A2-003** - CGA-011 fires on all untrusted endpoints unconditionally. Rename at minimum, or add validation detection for real signal.

## Change Impact Summary

Changed symbols and files in `7796f12..HEAD`:

**Parsers (complexity computation):**
- `computeComplexity`, `countDecisionPoints`, `countDecisionPointsSkipSelf` added to `go_parser.go`
- `computePythonComplexity`, `countPythonDecisionPoints` added to `python_parser.go`
- `computeTypeScriptComplexity`, `countTypeScriptDecisionPoints`, `countTypeScriptDecisionPointsSkipSelf` added to `typescript_parser.go`
- `computeRustComplexity`, `countRustDecisionPoints` added to `rust_parser.go`

**Security annotators (trust classification):**
- `classifyTrust` added to `go_annotator.go`, `python_annotator.go`, `typescript_annotator.go`, `rust_annotator.go`

**Security queries:**
- `queryComplexityHotspot` (CGA-010) and `queryUntrustedWithoutValidation` (CGA-011) added to `queries.go`

**Test fixtures:**
- `complexity_sample.{go,py,rs,ts}` created in `testdata/`

**Callers affected:** `Annotate()` methods in all 4 annotators now call `classifyTrust` as a third pass. `extractFunction` in all 4 parsers now computes complexity before appending to results.

> Advisory: impact graph is based on diff analysis and may be incomplete.

## Review Metrics

- Findings raised: 27 (raw, across all 5 specialists)
- Findings surviving challenge: 6 (deduplicated)
- Findings dismissed: 0 (0%)
- Consensus rate: 100% (all 6 findings had majority or unanimous agreement)
- Forced convergence: 0 agents

## Guardrails Triggered

None.

## Audit Log

No external actions taken.

<!-- REVIEW METADATA
timestamp: 2026-04-27T19:30:00Z
commit_sha: 9b752c597c19e95aa39eb0d2d4097cac47307d75
reviewed_files:
  pkg/domains/security/go_annotator.go c90261216116dfb182f965f22f5b17b8a8a2ce316d3ccb6e535c18b92e225d77
  pkg/domains/security/go_annotator_test.go ec3ca6513e5075111b55e003b8aa1edbc31076df7ec048f6315c8ee8d9aabb20
  pkg/domains/security/python_annotator.go a7b679c9ea8c48b0095bc2aa37d70932cffe81a72c2c7161d8c55e22d795ae47
  pkg/domains/security/queries.go db2128797d53a8360cdbc75d535ede278bf14f82643cb85e9d63d0bc347c533c
  pkg/domains/security/queries_test.go b59e79d2054d7bd2545005f281a46410e18e8426f3ab9b0545c9228625759e8e
  pkg/domains/security/rust_annotator.go 67dd6114bfc4bfdea4a4ba43ebfe470134ffdde3f1cbdfc75caa9a054841025b
  pkg/domains/security/typescript_annotator.go 254558662a096c0e34830b3721b24267dd55e132a846456a450c86c5eee44865
  pkg/parser/go_parser.go 396daafcd90a04de05b65c371f6aa2e22f3fba08e31feec499f0b60f894ca11c
  pkg/parser/go_parser_test.go 13a670e70981a1f67036c901edf14e37dbbe89310b551a9943ea5662cd1e1a89
  pkg/parser/python_parser.go 8d5506ad280f930808b43a1d247d87305790e2ce1bd023ae6412a81b17598808
  pkg/parser/python_parser_test.go c705ba5d13ffd7bb5059e386c32d857873e7b646e04badaece32c754ce5953a7
  pkg/parser/rust_parser.go e2dea5448897312d268b68e1343f041d96d0aefd6ff27cf074434c355e291c57
  pkg/parser/rust_parser_test.go 0f68b806928749710774284093649ee416abeb3266edec0793696a3b33cc32d1
  pkg/parser/typescript_parser.go f2a9634c8889748343b7e4b41919576b9d7fb6bb9f7b6f9858b207cc6bb1c917
  pkg/parser/typescript_parser_test.go cc5726fb4845da93b42da617b75ab9a898395d34c42391261d619b3c40ec8efb
  testdata/complexity_sample.go 4ba674bf3dd7da7b6775de0ca496de914a8c68886fbe39fdd7d29284f2234f8b
  testdata/complexity_sample.py 31a30a86dce08fd6cb9e157667c7ed2faecb3621b470588fb19ad235a851f88d
  testdata/complexity_sample.rs f6bb0579a14e0af7708a5c7692b2957f03ece31c4eeeab02a22b7407b6d4b0c3
  testdata/complexity_sample.ts 3b1ca68acaec73036b42563e2bdb2de76706c1c44b5295c009f2e5586269fba6
specialists: SEC, PERF, QUAL, CORR, ARCH
configuration: 1 iteration, convergence achieved (cross-specialist agreement), flags: --diff --range 7796f12..HEAD
-->
