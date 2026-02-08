# DESIGN.md -- Greenlight Brownfield Support

> **Project:** Greenlight
> **Scope:** Add brownfield support: `/gl:assess` (gap analysis) and `/gl:wrap` (locking tests + contract extraction). Update existing commands for brownfield awareness.
> **Stack:** Go 1.24, stdlib only. But the deliverables are markdown prompt files, not Go code.
> **Date:** 2026-02-08
> **Replaces:** Previous DESIGN.md (CLI stabilisation -- complete, 228 tests passing)

---

## 1. Requirements

### 1.1 Functional Requirements -- /gl:assess

#### FR-1: Prerequisites and Context Loading

| ID | Requirement |
|----|-------------|
| FR-1.1 | MUST read `.greenlight/codebase/` docs produced by `/gl:map` as primary input |
| FR-1.2 | MUST read `.greenlight/config.json` for project context (stack, test commands, project directories) |
| FR-1.3 | MUST work without `/gl:map` having been run; warn that results will be shallow and recommend running `/gl:map` first |
| FR-1.4 | MUST NOT require any other command to have been run first (assess is always available) |
| FR-1.5 | MUST read `CLAUDE.md` engineering standards for gap comparison baseline |

#### FR-2: Test Coverage Analysis

| ID | Requirement |
|----|-------------|
| FR-2.1 | MUST map source files to corresponding test files using stack-specific naming conventions |
| FR-2.2 | MUST detect test file patterns per stack: Go (`*_test.go` co-located), Python (`test_*.py` / `*_test.py` in `tests/` or co-located), JS/TS (`*.test.{js,ts}` / `*.spec.{js,ts}` in `__tests__/` or co-located), Rust (`#[cfg(test)]` inline + `tests/` dir) |
| FR-2.3 | MUST calculate source-to-test file ratio (e.g., "47 source files, 12 test files, 25% coverage by file count") |
| FR-2.4 | MUST flag source files with zero test coverage (no corresponding test file found) |
| FR-2.5 | IF `test.coverage_command` is configured in `config.json`, MUST run it and parse output for line/branch coverage percentages |
| FR-2.6 | IF `test.coverage_command` is NOT configured, MUST note in output that coverage percentages are unavailable; file mapping still runs |
| FR-2.7 | MUST classify each source directory/module as: tested (>50% files have tests), partially tested (1-50%), untested (0%) |

#### FR-3: Contract Inventory

| ID | Requirement |
|----|-------------|
| FR-3.1 | MUST scan for typed schemas, interfaces, validation logic, and API route definitions |
| FR-3.2 | MUST classify each boundary as: **explicit** (typed interface/schema/contract exists), **implicit** (behaviour exists but no formal type definition), **none** (untyped, unvalidated boundary) |
| FR-3.3 | MUST identify external boundaries: API endpoints, database queries, third-party service calls, webhook handlers, message consumers |
| FR-3.4 | MUST identify internal module boundaries: exported functions, public interfaces, cross-package imports |
| FR-3.5 | MUST record the source file and line range for each identified boundary |

#### FR-4: Risk Assessment

| ID | Requirement |
|----|-------------|
| FR-4.1 | MUST spawn `gl-security` agent in `full-audit` mode for security scanning |
| FR-4.2 | MUST identify fragile areas: high cyclomatic complexity, deep nesting (>4 levels), long functions (>50 lines), high fan-in/fan-out |
| FR-4.3 | MUST identify critical paths: authentication flows, authorization checks, payment processing, data mutation endpoints, admin operations |
| FR-4.4 | MUST assess tech debt hotspots: `TODO`/`FIXME`/`HACK` comments with file locations, deprecated API usage, outdated dependency patterns |
| FR-4.5 | SHOULD check dependency health: outdated versions, known vulnerabilities (via `npm audit` / `pip audit` / `go vet` / `cargo audit` as appropriate) |
| FR-4.6 | Security findings from gl-security MUST be included in ASSESS.md under a dedicated section |

#### FR-5: Architecture Gap Analysis

| ID | Requirement |
|----|-------------|
| FR-5.1 | MUST compare existing codebase against CLAUDE.md engineering standards (error handling, naming, functions, security, API design, database, testing, logging, file structure, git, performance) |
| FR-5.2 | MUST identify specific standard violations with file paths and line numbers where possible |
| FR-5.3 | MUST categorize gaps by severity: CRITICAL (security/correctness risk), HIGH (maintainability risk), MEDIUM (quality concern), LOW (style/convention) |
| FR-5.4 | MUST produce a standards compliance summary table showing pass/fail per CLAUDE.md section |

#### FR-6: Output (ASSESS.md)

| ID | Requirement |
|----|-------------|
| FR-6.1 | MUST produce `.greenlight/ASSESS.md` as structured output following the schema in section 4.1 |
| FR-6.2 | MUST include a prioritized wrap recommendation as a priority-tiered list: Critical, High, Medium tiers |
| FR-6.3 | Each recommended boundary MUST include: name, type (external/internal), current contract status (explicit/implicit/none), test status, estimated complexity, risk level |
| FR-6.4 | MUST commit ASSESS.md with conventional commit: `docs: greenlight codebase assessment` |
| FR-6.5 | MUST report summary to user after completion: boundary count, coverage stats, critical findings count, recommended next action |

### 1.2 Functional Requirements -- /gl:wrap

#### FR-7: Prerequisites

| ID | Requirement |
|----|-------------|
| FR-7.1 | MUST read `.greenlight/ASSESS.md` if it exists; use it as primary input for prioritization |
| FR-7.2 | MUST work without ASSESS.md; user can choose what to wrap manually |
| FR-7.3 | MUST read `.greenlight/config.json` for test commands and project context |
| FR-7.4 | MUST read `.greenlight/codebase/` docs for codebase understanding |
| FR-7.5 | MUST read existing `.greenlight/CONTRACTS.md` if it exists; do not duplicate existing contracts |
| FR-7.6 | MUST read `CLAUDE.md` engineering standards |

#### FR-8: Prioritize What to Wrap

| ID | Requirement |
|----|-------------|
| FR-8.1 | IF ASSESS.md exists, MUST present the priority-tiered wrap recommendation from the assessment |
| FR-8.2 | IF ASSESS.md does not exist, MUST scan the codebase for boundaries and present them to the user |
| FR-8.3 | MUST let the user pick which boundary/module to wrap |
| FR-8.4 | MUST show estimated complexity for each candidate: file count, function count, dependency count |
| FR-8.5 | MUST prevent wrapping the entire codebase at once; enforce one boundary per wrap invocation |
| FR-8.6 | Wrapper agent MUST assess whether the selected boundary fits within 50% context budget; if too large, suggest splitting into sub-boundaries and present them for selection |

#### FR-9: Extract Contracts

| ID | Requirement |
|----|-------------|
| FR-9.1 | MUST analyse existing implementation code to infer implicit contracts |
| FR-9.2 | MUST identify: function signatures, parameter types, return types, error types/patterns, validation rules, invariants observable from code |
| FR-9.3 | MUST present extracted contracts to the user for review and confirmation before writing |
| FR-9.4 | MUST write approved contracts to `.greenlight/CONTRACTS.md` with `[WRAPPED]` tag |
| FR-9.5 | Contracts MUST be descriptive (what code DOES), not prescriptive (what it SHOULD do) |
| FR-9.6 | `[WRAPPED]` tag MUST include: `Source: {file}:{lines}`, `Wrapped on: {date}`, `Locking tests: tests/locking/{boundary}.test.{ext}` |
| FR-9.7 | MUST follow existing contract format from `gl-architect.md`: input, output, errors, invariants, security, dependencies |

#### FR-10: Write Locking Tests

| ID | Requirement |
|----|-------------|
| FR-10.1 | MUST generate tests that verify existing behaviour (locking tests, not specification tests) |
| FR-10.2 | Locking tests MUST pass against existing code without any source code changes |
| FR-10.3 | Locking tests MUST go in `tests/locking/{boundary-name}.test.{ext}` -- one file per boundary |
| FR-10.4 | MUST test both happy paths and observable error paths of existing code |
| FR-10.5 | MUST NOT test implementation details; test observable behaviour at the boundary (inputs/outputs) |
| FR-10.6 | MUST handle non-deterministic behaviour automatically: timestamps (freeze time), random IDs (use matchers), environment-dependent values (mock env) |
| FR-10.7 | IF a locking test cannot be made to pass after 2 attempts due to non-determinism or complexity, document the specific assertion in ASSESS.md as non-deterministic and skip it |
| FR-10.8 | Test names MUST use descriptive `[LOCK]` prefix: `[LOCK] should return user object when valid email provided` |

#### FR-11: Run and Verify

| ID | Requirement |
|----|-------------|
| FR-11.1 | MUST run locking tests after generation using `config.test.command` |
| FR-11.2 | ALL locking tests MUST pass; if they don't, the tests are wrong (not the code) |
| FR-11.3 | IF any locking test fails, the wrapper MUST fix the test (not the code) and re-run |
| FR-11.4 | Maximum 3 test-fix-rerun cycles per boundary before escalating to user |
| FR-11.5 | MUST report test results: total tests, passing, any skipped due to non-determinism |
| FR-11.6 | MUST run full test suite after locking tests pass to ensure no regressions |

#### FR-12: Security Baseline

| ID | Requirement |
|----|-------------|
| FR-12.1 | MUST spawn `gl-security` agent in `slice` mode scoped to the wrapped boundary's files |
| FR-12.2 | Security issues MUST be documented in ASSESS.md (or created if it doesn't exist), NOT fixed during wrap |
| FR-12.3 | Security findings MUST be recorded as known issues in the Wrapped Boundaries table of STATE.md |
| FR-12.4 | Security agent MUST NOT write failing tests during wrap (unlike in /gl:slice); document only |

#### FR-13: Commit and Track

| ID | Requirement |
|----|-------------|
| FR-13.1 | MUST commit locking tests and extracted contracts atomically |
| FR-13.2 | Commit MUST use conventional format: `test(wrap): lock {boundary-name}` |
| FR-13.3 | Commit body MUST list: contracts extracted (count), locking tests written (count), known security issues (count) |
| FR-13.4 | MUST update `.greenlight/STATE.md` Wrapped Boundaries section with the new boundary |
| FR-13.5 | MUST update `.greenlight/CONTRACTS.md` with wrapped contracts |

#### FR-14: Next Action

| ID | Requirement |
|----|-------------|
| FR-14.1 | MUST guide user to either wrap another boundary or start building new features via `/gl:slice` |
| FR-14.2 | MUST show wrap progress: `{N} of {M} boundaries wrapped` (where M comes from ASSESS.md priority list, or is omitted if no ASSESS.md) |
| FR-14.3 | IF all Critical-tier boundaries are wrapped, MUST suggest moving to `/gl:design` or `/gl:slice` for new features |

### 1.3 Functional Requirements -- Updates to Existing Commands

#### FR-15: /gl:slice Updates (wraps field, locking-to-integration transition)

| ID | Requirement |
|----|-------------|
| FR-15.1 | GRAPH.json slices MAY include a `wraps` field referencing wrapped boundary names |
| FR-15.2 | When a slice has a `wraps` field, /gl:slice MUST read the existing locking tests for that boundary |
| FR-15.3 | Test writer MUST receive locking test names (not code) as context: "these are the existing locked behaviours" |
| FR-15.4 | After new integration tests pass and verification succeeds, locking tests for the wrapped boundary MUST be deleted |
| FR-15.5 | After locking tests are deleted, the `[WRAPPED]` tag MUST be removed from corresponding contracts in CONTRACTS.md |
| FR-15.6 | STATE.md Wrapped Boundaries status MUST update to `refactored` when locking tests are replaced |

#### FR-16: /gl:status Updates

| ID | Requirement |
|----|-------------|
| FR-16.1 | MUST display Wrapped Boundaries table if any wrapped boundaries exist in STATE.md |
| FR-16.2 | Table columns: Boundary, Contracts, Locking Tests (count), Known Issues (count), Status |
| FR-16.3 | Status values for wrapped boundaries: `wrapped` (locking tests in place), `refactored` (replaced by integration tests) |

#### FR-17: /gl:help Updates

| ID | Requirement |
|----|-------------|
| FR-17.1 | MUST add BROWNFIELD section between SETUP and BUILD sections |
| FR-17.2 | Commands listed: `/gl:assess` (Gap analysis and risk assessment), `/gl:wrap` (Extract contracts + locking tests) |
| FR-17.3 | FLOW line MUST update to: `map? -> assess? -> init -> design -> wrap? -> slice 1 -> ... -> ship` |

#### FR-18: /gl:settings Updates

| ID | Requirement |
|----|-------------|
| FR-18.1 | MUST display `assessor` and `wrapper` agent models in settings table |
| FR-18.2 | Valid agents list MUST include `assessor` and `wrapper` |

#### FR-19: /gl:design Updates (brownfield awareness)

| ID | Requirement |
|----|-------------|
| FR-19.1 | IF `.greenlight/ASSESS.md` exists, designer MUST receive assessment context |
| FR-19.2 | Designer MUST be aware of wrapped boundaries and their contracts when planning new features |
| FR-19.3 | Designer SHOULD suggest which wrapped boundaries to refactor as part of new feature slices |

#### FR-20: CLAUDE.md Updates

| ID | Requirement |
|----|-------------|
| FR-20.1 | Agent Isolation Rules table MUST add `gl-assessor` row: Can See (codebase docs, test results, standards), Cannot See (N/A -- read-only analytical agent), Cannot Do (modify any code) |
| FR-20.2 | Agent Isolation Rules table MUST add `gl-wrapper` row: Can See (implementation code, existing tests), Cannot See (N/A), Cannot Do (modify production code -- only writes contracts and locking tests) |
| FR-20.3 | MUST add a clearly marked exception note under Agent Isolation Rules: "gl-wrapper is a deliberate exception. It sees implementation code AND writes locking tests. This is necessary because locking tests must verify what code currently does, not what contracts say it should do. This exception is scoped: only applies to tests in `tests/locking/`. When a boundary is later refactored via /gl:slice, locking tests are deleted and proper integration tests are written under strict isolation." |

#### FR-21: gl-architect.md Updates

| ID | Requirement |
|----|-------------|
| FR-21.1 | MUST recognise `[WRAPPED]` contracts in CONTRACTS.md |
| FR-21.2 | MUST NOT redefine wrapped contracts; treat them as existing boundaries |
| FR-21.3 | When adding new slices, MUST be able to reference wrapped contracts as dependencies |
| FR-21.4 | When a slice's `wraps` field targets a wrapped boundary, architect MUST plan the contract transition: wrapped contract becomes a proper contract (tag removed) |

#### FR-22: gl-test-writer.md Updates

| ID | Requirement |
|----|-------------|
| FR-22.1 | MUST check for existing locking tests when writing tests for a slice that wraps a boundary |
| FR-22.2 | MUST receive locking test NAMES (not source code) as context for understanding existing locked behaviour |
| FR-22.3 | Integration tests MUST cover at least all behaviours that locking tests covered (superset) |
| FR-22.4 | MUST NOT be aware of locking test implementation -- only names/descriptions |

### 1.4 Non-Functional Requirements

| ID | Category | Requirement |
|----|----------|-------------|
| NFR-1 | Context budget | gl-assessor MUST complete within 50% context window. If codebase is too large, split analysis by directory/module and aggregate |
| NFR-2 | Context budget | gl-wrapper MUST complete one boundary wrap within 50% context window. If a boundary is too large, it MUST suggest splitting before proceeding |
| NFR-3 | Idempotency | Running `/gl:assess` multiple times MUST overwrite ASSESS.md (latest assessment replaces previous) |
| NFR-4 | Idempotency | Running `/gl:wrap` on an already-wrapped boundary MUST warn and ask before overwriting existing locking tests |
| NFR-5 | Safety | `/gl:wrap` MUST NEVER modify production source code. It writes contracts and locking tests only |
| NFR-6 | Safety | `/gl:assess` is entirely read-only except for writing ASSESS.md |
| NFR-7 | Compatibility | Both commands MUST work with any stack Greenlight supports (Go, Python, JS/TS, Rust, Swift) |
| NFR-8 | Fault tolerance | If gl-security agent fails during assess/wrap, the command MUST continue without security findings (warn user, note in output) |

### 1.5 Constraints

| Constraint | Detail |
|------------|--------|
| Deliverables | Markdown prompt files in `src/agents/` and `src/commands/gl/`. NOT Go code |
| Embedding | New `.md` files must be added to Go manifest and `go:embed` directive |
| Agent models | gl-assessor defaults to sonnet (analytical, not decision-making). gl-wrapper defaults to sonnet (follows extracted contracts) |
| Existing flow | New commands MUST NOT break existing greenfield flow. `/gl:assess` and `/gl:wrap` are optional |
| Commit format | All commits use conventional format enforced by lefthook/commitlint |

### 1.6 Out of Scope

| Item | Rationale |
|------|-----------|
| Auto-refactoring during wrap | Wrap locks existing behaviour. Refactoring happens in subsequent slices |
| Whole-codebase contract generation | One boundary at a time. Prevents context overflow and ensures quality |
| Failing wrap on poor code quality | Wrap works with any code. That is the point |
| Making /gl:assess mandatory | Assessment is recommended, not required. Users can wrap manually |
| Migration from GSD's `.planning/` folder | Different system, different data model. Not compatible |
| AST parsing or static analysis tooling | Claude reads code natively. No external tooling needed |
| Automated dependency updates | Assess identifies outdated deps. Fixing them is a separate task |
| Test generation for untested code that is NOT being wrapped | Wrap is opt-in per boundary. Assess identifies gaps; wrap addresses them one at a time |

---

## 2. Technical Decisions

| # | Decision | Chosen | Rejected | Rationale |
|---|----------|--------|----------|-----------|
| TD-1 | Where extracted contracts go | **CONTRACTS.md with `[WRAPPED]` tag** | Separate CONTRACTS-WRAPPED.md; Inline in ASSESS.md | One source of truth. Tag distinguishes wrapped from greenfield. Architect, test writer, and implementer already read CONTRACTS.md. Tag lifecycle: removed when boundary is refactored via /gl:slice |
| TD-2 | gl-wrapper isolation | **Deliberate exception: sees code AND writes locking tests** | Split into two agents (extractor + test writer); Strict isolation | Locking tests must verify what code DOES. You cannot write characterization tests without seeing the implementation. Exception is scoped to `tests/locking/` only. Greenfield tests via /gl:slice maintain strict isolation |
| TD-3 | Test coverage detection | **File mapping (always) + coverage command (optional)** | Coverage only; File mapping only | File mapping gives 80% of value with zero side effects. Running tests on unknown codebase can have side effects (database writes, network calls). Coverage runs only when explicitly configured |
| TD-4 | Wrapped boundary tracking in STATE.md | **Separate "Wrapped Boundaries" section** | Mixed into slices table; Track in ASSESS.md only | Different lifecycle (wrapped vs. pending/implementing/complete). Separate section keeps slice table clean. Shows at a glance: "4/7 boundaries wrapped, 3/5 feature slices complete" |
| TD-5 | Assess output format | **Priority-tiered list (Critical/High/Medium)** | Sequenced order with dependencies | Wrap order rarely has real dependencies. What matters is risk priority. Users pick from highest tier. Less overhead than maintaining a dependency graph for wrapping |
| TD-6 | Locking test location | **`tests/locking/{boundary}.test.{ext}` -- one file per boundary** | One file per contract; Mirror source structure | Clean directory, obvious which boundaries are wrapped, easy to delete when refactored. Matches how boundaries are tracked in STATE.md |
| TD-7 | Non-deterministic test handling | **Auto-handle with fallback to documenting** | Require user intervention; Skip all non-deterministic code | Wrapper sees the code and can spot `Date.now()`, `uuid()`, etc. Auto-adapts (matchers, time freezing). After 2 failed attempts, documents in ASSESS.md and moves on. Pragmatic balance |
| TD-8 | Boundary scope assessment | **Wrapper agent decides, suggests splitting if too large** | Always wrap one file; Always wrap entire module | Agent assesses complexity against 50% context budget. Small boundaries wrap at once. Large boundaries get split into sub-boundaries presented to user. Context-aware, not arbitrary |
| TD-9 | gl-assessor model | **Sonnet (default in balanced profile)** | Opus | Analytical, read-only work. No architectural decisions. Sonnet handles pattern matching, file scanning, and reporting reliably |
| TD-10 | gl-wrapper model | **Sonnet (default in balanced profile)** | Opus | Follows contracts and patterns. Work is constrained (extract what exists, write tests that pass). TDD loop catches quality issues. Upgrade to opus for complex legacy codebases |
| TD-11 | Security during wrap | **Document only, no failing tests** | Write failing tests (like /gl:slice) | Wrap locks existing behaviour without changing it. Writing failing security tests would require code changes, which violates the wrap principle. Security issues are documented as known gaps for future slices to address |
| TD-12 | GRAPH.json wraps field | **Optional `wraps` field on slices** | Separate wraps graph; No formal link | Minimal addition to existing data model. Slices can reference what they refactor. Dependencies resolve naturally: slice depends on understanding what the wrapped boundary does |

---

## 3. Architecture

### 3.1 Brownfield Flow Integration

The brownfield commands slot into the existing Greenlight flow as optional steps:

```
GREENFIELD FLOW (existing):
  map? -> init -> design -> slice 1 -> slice 2 -> ... -> ship

BROWNFIELD FLOW (new):
  map -> assess? -> init -> design -> wrap? -> slice 1 -> ... -> ship
         ^^^^^                        ^^^^
         NEW                          NEW

Detailed brownfield sequence:
  /gl:map        -> understand the codebase (EXISTING)
  /gl:assess     -> identify gaps, risks, untested boundaries (NEW)
  /gl:design     -> plan what to fix and build, informed by assessment (EXISTING, updated)
  /gl:wrap       -> extract contracts + write locking tests for existing code (NEW)
  /gl:slice 1    -> new features or refactored features, TDD as normal (EXISTING, updated)
```

Both `/gl:assess` and `/gl:wrap` are optional. A user can:
- Skip assess entirely and wrap manually
- Skip wrap entirely and build on top of existing code
- Assess, wrap critical boundaries, then slice new features
- Use any combination

### 3.2 Agent Architecture

#### New Agents

```
gl-assessor (NEW)
  Role: Analytical, read-only codebase assessment
  Tools: Read, Bash, Glob, Grep
  Model: sonnet (balanced profile)
  Spawned by: /gl:assess
  Can see: codebase docs, source code, test files, config, CLAUDE.md standards
  Cannot do: modify any code or write any files except ASSESS.md
  Isolation: read-only analytical agent, no special exceptions

gl-wrapper (NEW)
  Role: Contract extraction from existing code + locking test generation
  Tools: Read, Write, Bash, Glob, Grep
  Model: sonnet (balanced profile)
  Spawned by: /gl:wrap
  Can see: implementation code, existing tests, contracts, config
  Cannot do: modify production source code
  Isolation: DELIBERATE EXCEPTION -- sees implementation AND writes locking tests
  Exception scope: only tests/locking/ directory
```

#### Updated Agent Interactions

```
/gl:assess orchestrator
  |
  +--> gl-assessor (parallel focus areas, similar to /gl:map pattern)
  |      +--> coverage analysis
  |      +--> contract inventory
  |      +--> architecture gaps
  |
  +--> gl-security (full-audit mode, spawned by assessor or orchestrator)
  |
  +--> writes ASSESS.md (aggregated from all findings)


/gl:wrap orchestrator
  |
  +--> present boundary selection (from ASSESS.md or manual scan)
  +--> user picks boundary
  |
  +--> gl-wrapper (fresh context for each boundary)
  |      +--> read implementation code
  |      +--> extract contracts
  |      +--> present contracts to user for confirmation
  |      +--> write locking tests
  |      +--> run and fix tests (max 3 cycles)
  |
  +--> gl-security (slice mode, scoped to boundary files -- document only)
  |
  +--> commit (tests + contracts atomically)
  +--> update STATE.md
```

### 3.3 Data Flow

```
                    .greenlight/codebase/
                    (from /gl:map)
                          |
                          v
  CLAUDE.md -------> /gl:assess -------> .greenlight/ASSESS.md
  config.json ----/        |
                           v
                      gl-security
                      (full-audit)
                           |
                           v
                    Security findings
                    merged into ASSESS.md


  ASSESS.md -------> /gl:wrap ---------> .greenlight/CONTRACTS.md ([WRAPPED] entries)
  config.json ----/      |               tests/locking/{boundary}.test.{ext}
  source code --------/  |               .greenlight/STATE.md (Wrapped Boundaries section)
                         v
                    gl-security
                    (document only)
                         |
                         v
                    Known issues
                    in STATE.md


  CONTRACTS.md -----> /gl:slice --------> Refactored code
  (with [WRAPPED])       |               tests/integration/ (new tests)
  STATE.md ----------/   |               tests/locking/ (deleted)
  GRAPH.json --------/   |               CONTRACTS.md ([WRAPPED] tag removed)
  (with wraps field)     v               STATE.md (boundary status -> refactored)
                    Normal TDD loop
                    (strict isolation)
```

### 3.4 Command Orchestration Patterns

Both new commands follow the existing Greenlight orchestration pattern:

1. **Orchestrator reads state** (config.json, ASSESS.md, STATE.md)
2. **Orchestrator resolves models** from config.json profiles
3. **Orchestrator spawns agent** via Task with structured context (XML blocks)
4. **Agent writes output** directly to files (agent doesn't return content to orchestrator -- saves context)
5. **Orchestrator verifies output** (file exists, non-empty)
6. **Orchestrator commits** with conventional format
7. **Orchestrator reports** summary and next action

This mirrors `/gl:map` (parallel agents writing directly) and `/gl:slice` (sequential agent pipeline).

---

## 4. Data Model

### 4.1 ASSESS.md Structure

```markdown
# Codebase Assessment

Generated: {YYYY-MM-DD}
Project: {project name}
Stack: {stack from config.json}

## Summary

| Metric | Value |
|--------|-------|
| Source files | {N} |
| Test files | {N} |
| File coverage | {N}% |
| Line coverage | {N}% or "not configured" |
| Boundaries identified | {N} |
| Explicit contracts | {N} |
| Implicit contracts | {N} |
| No contract | {N} |
| Security findings | {N} (C:{N} H:{N} M:{N} L:{N}) |
| Standards compliance | {N}/{M} sections passing |

## Test Coverage

### By Module

| Module | Source Files | Test Files | Coverage | Status |
|--------|-------------|------------|----------|--------|
| {module} | {N} | {N} | {N}% | tested / partial / untested |

### Untested Files

| File | Type | Risk | Recommended Priority |
|------|------|------|---------------------|
| {path} | {endpoint/service/util} | {high/medium/low} | Critical / High / Medium |

## Contract Inventory

### Boundaries

| # | Boundary | Type | Contract Status | Location | Tests |
|---|----------|------|----------------|----------|-------|
| 1 | {name} | external/internal | explicit/implicit/none | {file}:{lines} | {yes/no} |

### Summary by Status

| Status | Count | Percentage |
|--------|-------|------------|
| Explicit | {N} | {N}% |
| Implicit | {N} | {N}% |
| None | {N} | {N}% |

## Risk Assessment

### Security Findings

{Output from gl-security full-audit, or "Security scan not performed" if agent failed}

| # | Severity | Category | Location | Description |
|---|----------|----------|----------|-------------|
| 1 | {CRITICAL/HIGH/MEDIUM/LOW} | {category} | {file}:{line} | {description} |

### Fragile Areas

| File | Concern | Severity | Detail |
|------|---------|----------|--------|
| {path} | complexity / nesting / length / coupling | {severity} | {specifics} |

### Tech Debt

| File | Type | Detail |
|------|------|--------|
| {path} | TODO / FIXME / HACK / deprecated / outdated | {comment or description} |

## Architecture Gaps

### Standards Compliance

| CLAUDE.md Section | Status | Key Gaps |
|-------------------|--------|----------|
| Error Handling | pass/partial/fail | {brief description if not pass} |
| Naming | pass/partial/fail | {brief description} |
| Functions | pass/partial/fail | {brief description} |
| Security | pass/partial/fail | {brief description} |
| API Design | pass/partial/fail | {brief description} |
| Database | pass/partial/fail | {brief description} |
| Testing | pass/partial/fail | {brief description} |
| Logging | pass/partial/fail | {brief description} |
| File Structure | pass/partial/fail | {brief description} |
| Git | pass/partial/fail | {brief description} |
| Performance | pass/partial/fail | {brief description} |

### Specific Violations

| # | Standard | Violation | Location | Severity |
|---|----------|-----------|----------|----------|
| 1 | {section} | {what's wrong} | {file}:{line} | {severity} |

## Wrap Recommendations

Boundaries recommended for wrapping, grouped by priority tier.

### Critical -- Wrap These First

| # | Boundary | Type | Why Critical | Estimated Complexity |
|---|----------|------|-------------|---------------------|
| 1 | {name} | {type} | {rationale} | {low/medium/high} |

### High -- Wrap Before New Features

| # | Boundary | Type | Why High | Estimated Complexity |
|---|----------|------|----------|---------------------|
| 1 | {name} | {type} | {rationale} | {low/medium/high} |

### Medium -- Wrap When Convenient

| # | Boundary | Type | Why Medium | Estimated Complexity |
|---|----------|------|-----------|---------------------|
| 1 | {name} | {type} | {rationale} | {low/medium/high} |

## Next Steps

1. {Primary recommendation -- usually "wrap Critical boundaries"}
2. {Secondary recommendation -- usually "run /gl:design for new features"}
3. {Tertiary recommendation}
```

### 4.2 CONTRACTS.md [WRAPPED] Format

Wrapped contracts follow the standard contract format from `gl-architect.md` with an additional metadata block:

```markdown
### Contract: {BoundaryName} [WRAPPED]

**Source:** `{file}:{start_line}-{end_line}`
**Wrapped on:** {YYYY-MM-DD}
**Locking tests:** `tests/locking/{boundary-name}.test.{ext}`

**Boundary:** {what talks to what}
**Slice:** wrappable (available for refactoring via /gl:slice with wraps field)

**Input:**
```{language}
{inferred input type/interface}
```

**Output:**
```{language}
{inferred output type/interface}
```

**Errors:**
| Error | Status/Type | When |
|-------|-------------|------|
| {error} | {code/type} | {condition observed in code} |

**Invariants:**
- {observed invariant from existing code behaviour}

**Security:**
- Known issues: {list from security baseline, or "none identified"}

**Dependencies:** {other contracts this uses}
```

**Agent behaviour rules for `[WRAPPED]` contracts:**

| Agent | Behaviour with [WRAPPED] contracts |
|-------|-----------------------------------|
| gl-architect | Do NOT redefine. Reference as existing. Can add `wraps` field to new slices targeting this boundary |
| gl-test-writer | Check for existing locking tests. When writing tests for a slice with `wraps` field, receive locking test NAMES as context. Integration tests must be a superset of locked behaviours |
| gl-implementer | Build on top of existing code. Use wrapped contract as constraint. When refactoring a `wraps` slice, existing code is the starting point |
| gl-security | Note known issues from wrapped contract. Check if issues persist after refactoring |
| gl-verifier | Verify that locking tests are removed after successful refactoring. Verify `[WRAPPED]` tag is removed |

**Lifecycle:**

```
[WRAPPED] contract created by /gl:wrap
  |
  v
Slice with wraps field targets this boundary (/gl:slice)
  |
  v
Test writer receives locking test names as context
Integration tests written (superset of locked behaviours)
  |
  v
Implementer refactors code, making integration tests pass
  |
  v
Verification: both locking tests AND integration tests pass
  |
  v
Locking tests deleted from tests/locking/
[WRAPPED] tag removed from contract
STATE.md boundary status -> refactored
```

### 4.3 STATE.md Wrapped Boundaries Section

Added below the existing Slices section:

```markdown
## Wrapped Boundaries

| Boundary | Contracts | Locking Tests | Known Issues | Status |
|----------|-----------|---------------|--------------|--------|
| {name} | {N} | {N} | {N} | wrapped / refactored |

Wrap progress: {N}/{M} boundaries wrapped
```

**Status values:**

| Status | Meaning |
|--------|---------|
| wrapped | Locking tests in place, contracts extracted, existing behaviour locked |
| refactored | Replaced by integration tests via /gl:slice. Locking tests deleted |

**Size constraint:** This section counts toward STATE.md's 80-line budget. Keep entries concise. When a boundary is refactored, it can be compressed to a single summary line or removed entirely.

### 4.4 GRAPH.json wraps Field

Optional field on slice objects:

```json
{
  "id": "refactor-auth",
  "name": "Refactor authentication with proper contracts",
  "description": "Replace implicit auth behaviour with explicit contracts and integration tests",
  "contracts": ["AuthenticateUser", "ValidateToken"],
  "depends_on": [],
  "wraps": ["auth"],
  "priority": 2,
  "estimated_tests": 12,
  "boundaries": ["Client -> API", "API -> Database"]
}
```

**`wraps` field rules:**
- Array of boundary names matching entries in STATE.md Wrapped Boundaries table
- When present, `/gl:slice` knows to read locking tests and handle the transition
- A slice can wrap multiple boundaries if they're closely related
- The `wraps` field does NOT create a dependency on the boundary being wrapped first (the boundary is already wrapped by /gl:wrap before the slice is planned)

### 4.5 config.json Updates

New agent entries in profiles:

```json
{
  "profiles": {
    "quality": {
      "assessor": "opus",
      "wrapper": "opus"
    },
    "balanced": {
      "assessor": "sonnet",
      "wrapper": "sonnet"
    },
    "budget": {
      "assessor": "haiku",
      "wrapper": "sonnet"
    }
  }
}
```

**Note:** These are added to the profile definitions in `templates/config.md` and `commands/gl/init.md`. The Go CLI does not need code changes -- config.json is a data file read by Claude Code agents, not parsed by Go.

---

## 5. API Surface

### 5.1 /gl:assess Command

**Invocation:** `/gl:assess`

**Arguments:** None. Assesses the entire codebase.

**Prerequisites:**
- `.greenlight/config.json` must exist (run `/gl:init` first)

**Optional prior step:**
- `/gl:map` (recommended but not required)

**Output files:**
- `.greenlight/ASSESS.md` (created or overwritten)

**Commit:** `docs: greenlight codebase assessment`

**User interaction:**
- Minimal. Assessment is mostly automated
- Shows progress as each analysis phase completes
- Presents final summary with recommended next actions

**Report format:**
```
Assessment complete.

Source files: {N}
Test files: {N} ({coverage}% file coverage)
Boundaries: {N} ({explicit} explicit, {implicit} implicit, {none} no contract)
Security findings: {N} (CRITICAL: {N}, HIGH: {N}, MEDIUM: {N}, LOW: {N})
Standards compliance: {pass}/{total} sections passing

Wrap recommendations:
  Critical: {N} boundaries
  High: {N} boundaries
  Medium: {N} boundaries

ASSESS.md written to .greenlight/ASSESS.md

Next: Run /gl:wrap to lock existing boundaries with tests.
      Or /gl:design to plan new features informed by this assessment.
```

### 5.2 /gl:wrap Command

**Invocation:** `/gl:wrap`

**Arguments:** None. Interactive boundary selection.

**Prerequisites:**
- `.greenlight/config.json` must exist

**Optional prior steps:**
- `/gl:assess` (provides prioritized boundary list)
- `/gl:map` (provides codebase understanding)

**Output files:**
- `.greenlight/CONTRACTS.md` (created or appended with `[WRAPPED]` contracts)
- `tests/locking/{boundary-name}.test.{ext}` (locking test file)
- `.greenlight/STATE.md` (Wrapped Boundaries section updated)
- `.greenlight/ASSESS.md` (known issues updated, if file exists)

**Commit:** `test(wrap): lock {boundary-name}`

**User interaction:**
1. Present boundary candidates with priority (from ASSESS.md or fresh scan)
2. User picks boundary to wrap
3. Wrapper extracts contracts, presents for user confirmation
4. Wrapper writes locking tests, runs them
5. Report results, suggest next action

**Wrap session flow:**
```
/gl:wrap

Reading assessment... (.greenlight/ASSESS.md)

Recommended boundaries to wrap:

CRITICAL
  1. auth — external, no contract, untested, high complexity
  2. payments — external, implicit contract, untested

HIGH
  3. users — external, implicit contract, partial tests
  4. api-routes — external, no contract, tested

MEDIUM
  5. utils — internal, no contract, tested
  6. config — internal, implicit contract, tested

Which boundary? > 1

Analysing auth boundary...
  Files: src/auth/login.ts, src/auth/middleware.ts, src/auth/token.ts
  Functions: 12
  Estimated contracts: 4

Scope fits within context budget. Proceeding.

Extracted contracts:

  1. AuthenticateUser [WRAPPED]
     Input: { email: string, password: string }
     Output: { token: string, user: { id, email, name } }
     Errors: InvalidCredentials (401), ValidationError (400)
     Source: src/auth/login.ts:15-47

  2. ValidateToken [WRAPPED]
     ...

Accept these contracts? [y/N/edit] > y

Writing locking tests... tests/locking/auth.test.ts

Running locking tests...
  12 tests passing

Running full suite...
  All {N} tests passing (no regressions)

Security baseline...
  2 known issues documented (1 HIGH, 1 MEDIUM)

Committed: test(wrap): lock auth

Wrap progress: 1/6 boundaries wrapped

Next:
  /gl:wrap — wrap next boundary (recommended: payments)
  /gl:design — plan new features
  /gl:slice 1 — start building
```

### 5.3 Updated /gl:help Output

```
GREENLIGHT v1.x — TDD-first development for Claude Code

SETUP
  /gl:init              Brief interview + project config
  /gl:design            System design session -> DESIGN.md
  /gl:map               Analyse existing codebase first
  /gl:settings          Configure models, mode, options

BROWNFIELD
  /gl:assess            Gap analysis + risk assessment -> ASSESS.md
  /gl:wrap              Extract contracts + locking tests

BUILD
  /gl:slice <N>         TDD loop: test -> implement ->
                        security -> verify -> commit
  /gl:quick             Ad-hoc task with test guarantees
  /gl:add-slice         Add new slice to graph

MONITOR
  /gl:status            Real progress from test results
  /gl:pause             Save state for next session
  /gl:resume            Restore and continue

SHIP
  /gl:ship              Full audit + deploy readiness

FLOW
  map? -> assess? -> init -> design -> wrap? -> slice 1 -> ... -> ship

Tests are the source of truth. Green means done.
Security is built in, not bolted on.
```

### 5.4 Updated /gl:status Output (with wrapped boundaries)

```
GREENLIGHT STATUS

Tests:  {pass} passing  {fail} failing  {skip} skipped
Slices: {done}/{total}  [{progress_bar}]

1. {name}                complete  {N} tests ({S} sec)
2. {name}                failing   {N}/{M} passing
3. {name}                blocked   (needs 2)

Wrapped Boundaries:
  auth                   wrapped   12 locking tests  2 known issues
  payments               wrapped    8 locking tests  0 known issues
  users                  refactored (replaced by slice 1)

Wrap: 2/6 boundaries wrapped, 1 refactored

Next: /gl:slice 2 (fix failing)
  or: /gl:wrap (wrap next boundary)
```

---

## 6. Security

### 6.1 Security Agent in /gl:assess

The assessor spawns `gl-security` in `full-audit` mode. This is the same security agent used by `/gl:slice` and `/gl:ship`, with the same vulnerability checklist. The difference:

| Context | /gl:assess | /gl:slice | /gl:ship |
|---------|-----------|-----------|----------|
| Mode | full-audit | slice (diff only) | full-audit |
| Scope | Entire codebase | One slice's changes | Entire codebase |
| Output | Findings in ASSESS.md | Failing tests | Failing tests |
| Action | Document only | Tests must pass | Tests must pass |

During assessment, security findings are **documented, not enforced**. No failing tests are written. This is intentional -- the codebase may have known security issues that the team is aware of. Assessment provides a baseline; wrap and slice enforce fixes incrementally.

### 6.2 Security Agent in /gl:wrap

The wrapper spawns `gl-security` in `slice` mode scoped to the files in the wrapped boundary. Again, **document only** -- no failing tests.

Rationale: Writing failing security tests during wrap would require code changes to make them pass, which violates the wrap principle (lock existing behaviour, don't change it). Security issues are recorded as "known issues" in STATE.md's Wrapped Boundaries table and in the wrapped contract's Security section.

When the boundary is later refactored via `/gl:slice`, the security agent runs in its normal mode and writes failing tests. At that point, the code is being actively modified, so security fixes are part of the refactoring work.

### 6.3 Security During Locking-to-Integration Transition

When `/gl:slice` processes a slice with `wraps` field:

1. Normal TDD loop applies (test writer, implementer, security, verifier)
2. Security agent reviews the diff (including any refactored code)
3. Security agent writes failing tests for NEW vulnerabilities introduced
4. Security agent checks whether KNOWN issues from the wrapped contract's Security section have been addressed
5. If known issues persist, they are flagged but do not block the slice (they were pre-existing)

---

## 7. Deployment

### 7.1 New Files to Create

All files go in `src/` directory and are markdown prompt files:

| File | Type | Description |
|------|------|-------------|
| `src/agents/gl-assessor.md` | Agent definition | Analytical codebase assessment agent |
| `src/agents/gl-wrapper.md` | Agent definition | Contract extraction + locking test agent |
| `src/commands/gl/assess.md` | Command definition | /gl:assess orchestrator |
| `src/commands/gl/wrap.md` | Command definition | /gl:wrap orchestrator |

### 7.2 Existing Files to Update

| File | Change |
|------|--------|
| `src/agents/gl-architect.md` | Add `[WRAPPED]` contract recognition rules |
| `src/agents/gl-test-writer.md` | Add locking test awareness for wraps-field slices |
| `src/commands/gl/slice.md` | Add wraps field handling, locking-to-integration transition |
| `src/commands/gl/status.md` | Add Wrapped Boundaries table display |
| `src/commands/gl/help.md` | Add BROWNFIELD section, update FLOW line |
| `src/commands/gl/settings.md` | Add assessor and wrapper to agent list |
| `src/CLAUDE.md` | Add gl-assessor and gl-wrapper to isolation table, add locking test exception note |
| `README.md` (project root) | Add brownfield section with flow diagram |

### 7.3 Go CLI Changes (minimal)

These are the only Go code changes needed:

| File | Change |
|------|--------|
| `main.go` | No change needed -- `go:embed src/agents/*.md src/commands/gl/*.md` already uses wildcards that will pick up new `.md` files |
| `internal/installer/installer.go` | Add 4 new entries to `Manifest` slice: `agents/gl-assessor.md`, `agents/gl-wrapper.md`, `commands/gl/assess.md`, `commands/gl/wrap.md` |

The `go:embed` directive already uses `src/agents/*.md` and `src/commands/gl/*.md` globs, so new markdown files in those directories are automatically embedded. Only the Manifest slice (which controls install/uninstall/check) needs updating.

### 7.4 Test Updates

Existing Go tests that validate the manifest (if any) will need the expected file count updated from the current count to current + 4.

### 7.5 config.json Template Update

`src/templates/config.md` needs the profiles updated to include `assessor` and `wrapper` agent model entries. The init command template in `src/commands/gl/init.md` also needs the default config JSON updated.

---

## 8. Deferred

| Item | Source | Notes |
|------|--------|-------|
| Auto-refactoring during wrap | Spec (explicit exclusion) | Wrap locks behaviour. Refactoring is a subsequent slice |
| Whole-codebase contract generation | Spec (explicit exclusion) | One boundary at a time. Context budget enforcement |
| Wrap failure on poor code quality | Spec (explicit exclusion) | Wrap works with any code. That is the point |
| Mandatory /gl:assess | Spec (explicit exclusion) | Assessment is recommended, not required |
| GSD migration | Spec (explicit exclusion) | Different system, different data model |
| Automated wrap dependency graph | Design (TD-5) | Priority tiers are sufficient. Wrap order rarely has real dependencies |
| Coverage threshold enforcement | Design | Assess reports coverage. Enforcing thresholds is a future feature |
| Locking test mutation testing | Design | Could verify locking tests are meaningful (not tautological). Complex. Future consideration |
| Cross-boundary locking tests | Design | Current model wraps one boundary at a time. Cross-boundary interactions are tested via /gl:slice integration tests |
| Visual diff for contract extraction | Design | Could show side-by-side code vs. extracted contract. Nice-to-have UX improvement |
| Wrap undo / unwrap command | Design | Currently you'd delete locking tests and [WRAPPED] contracts manually. Formal undo could be a future command |

---

## 9. User Decisions

Decisions locked during the design session. These are final and should not be revisited without explicit user request.

### Gap Resolutions

| # | Gap | Decision | Rationale |
|---|-----|----------|-----------|
| UD-1 | Where extracted contracts go | **CONTRACTS.md with `[WRAPPED]` tag.** Same file as greenfield contracts, one source of truth. Each wrapped contract gets `[WRAPPED]` tag with Source, Wrapped on, and Locking tests fields. Tag tells agents how to handle: test writer checks for existing locking tests, implementer builds on top, architect doesn't redefine, security notes known gaps. Lifecycle: when refactored via /gl:slice, `[WRAPPED]` tag removed, locking tests replaced by integration tests. | Single source of truth. Agents already read CONTRACTS.md. Tag provides clean lifecycle |
| UD-2 | gl-wrapper isolation | **Deliberate exception: wrapper sees implementation code AND writes locking tests.** This is intentional -- you NEED to see what code does to lock its behaviour. Exception is scoped: only applies to locking tests in `tests/locking/`. When boundary is later refactored, locking tests are deleted and proper integration tests written via /gl:slice with strict isolation. CLAUDE.md needs a note documenting this exception. | Locking tests by definition test what code does. Cannot write them without seeing the code |
| UD-3 | Test coverage detection | **Both file mapping (always) and coverage command (optional).** File existence mapping always runs (fast, no config needed). Actual test coverage runs only if `test.coverage_command` is configured in config.json. If not configured, note in ASSESS.md that coverage percentages are unavailable. | Running tests on unknown codebase can have side effects. File mapping gives 80% of value |
| UD-4 | Wrapped boundary tracking | **Separate "Wrapped Boundaries" section in STATE.md.** Different lifecycle from slices: next -> wrapped (not the full slice status flow). Slices can depend on wrapped boundaries via `wraps` field in GRAPH.json. When a slice refactors a wrapped boundary: read locking tests, write integration tests, verify both pass, then remove locking tests and `[WRAPPED]` tag. | Different lifecycle. Keeps slice table clean. At-a-glance brownfield progress |

### Gray Area Decisions

| # | Gray Area | Decision | Rationale |
|---|-----------|----------|-----------|
| UD-5 | Locking test failure on first run (non-determinism) | **Wrapper auto-handles non-determinism (time freezing, matchers for random IDs, env mocking). Falls back to documenting assertion in ASSESS.md as non-deterministic after 2 failed attempts.** | Pragmatic. Wrapper sees the code and can identify non-deterministic patterns. Documenting unresolvable cases prevents blocking the entire wrap |
| UD-6 | Multi-file boundary scope | **Wrapper agent decides scope based on complexity vs. 50% context budget. If too large, suggests splitting into sub-boundaries and presents them for user selection.** | Context-aware, not arbitrary. Small boundaries wrap at once. Large ones split naturally |
| UD-7 | Assess output ordering | **Priority-tiered list (Critical / High / Medium), not sequenced order.** No dependency tracking between wrap targets. | Wrap order rarely has real dependencies. Priority tiers are simpler and sufficient |
| UD-8 | Locking test file organization | **One file per boundary: `tests/locking/{boundary-name}.test.{ext}`.** | Clean directory structure. Easy to see what's wrapped. One file to delete when refactored |

### Files Specification

| # | Category | Files |
|---|----------|-------|
| UD-9 | New files | `src/agents/gl-assessor.md`, `src/agents/gl-wrapper.md`, `src/commands/gl/assess.md`, `src/commands/gl/wrap.md` |
| UD-10 | Updated files | `src/agents/gl-architect.md`, `src/agents/gl-test-writer.md`, `src/commands/gl/slice.md`, `src/commands/gl/status.md`, `src/commands/gl/help.md`, `src/commands/gl/settings.md`, `src/CLAUDE.md`, `README.md` |
| UD-11 | Go CLI changes | `internal/installer/installer.go` (manifest update only). No new Go code |
