# Greenlight Agent Architecture

## System Overview

Greenlight is a TDD-first, contract-driven development framework for Claude Code. It orchestrates specialised agents through slash commands, enforcing strict isolation boundaries to prevent agents from seeing their own tests or implementation details.

```mermaid
graph TB
    User([User])

    subgraph Commands["Orchestrator Commands"]
        init["gl:init"]
        design["gl:design"]
        slice["gl:slice"]
        wrap["gl:wrap"]
        assess["gl:assess"]
        map["gl:map"]
        quick["gl:quick"]
        ship["gl:ship"]
        roadmap["gl:roadmap"]
        addslice["gl:add-slice"]
    end

    subgraph Agents["Specialised Agents"]
        designer["gl-designer"]
        architect["gl-architect"]
        testwriter["gl-test-writer"]
        implementer["gl-implementer"]
        security["gl-security"]
        verifier["gl-verifier"]
        assessor["gl-assessor"]
        wrapper["gl-wrapper"]
        mapper["gl-codebase-mapper"]
        debugger["gl-debugger"]
    end

    subgraph Artifacts[".greenlight/ Artifacts"]
        config["config.json"]
        interview["INTERVIEW.md"]
        designdoc["DESIGN.md"]
        contracts["CONTRACTS.md"]
        graph["GRAPH.json"]
        state["STATE.md"]
        assessdoc["ASSESS.md"]
        roadmapdoc["ROADMAP.md"]
        decisions["DECISIONS.md"]
        summaries["summaries/"]
        codebase["codebase/"]
    end

    subgraph ReadOnly["Read-Only Commands"]
        status["gl:status"]
        help["gl:help"]
        changelog["gl:changelog"]
        pause["gl:pause"]
        resume["gl:resume"]
        settings["gl:settings"]
    end

    User --> Commands
    User --> ReadOnly

    %% Command -> Agent spawning
    design -->|spawns| designer
    init -->|Phase 3| architect
    slice -->|1. tests| testwriter
    slice -->|2. implement| implementer
    slice -->|3. security| security
    slice -->|4. verify| verifier
    wrap -->|extract contracts| wrapper
    wrap -->|security scan| security
    assess -->|analyse| assessor
    map -->|4 parallel| mapper
    quick -->|debug| debugger
    ship -->|full audit| security
    ship -->|verify| verifier
    roadmap -->|milestone planning| designer
    addslice -->|design slice| architect

    %% Agent -> Artifact writes
    designer -->|writes| designdoc
    designer -->|writes| roadmapdoc
    designer -->|writes| decisions
    architect -->|writes| contracts
    architect -->|writes| graph
    testwriter -->|writes| testfiles["test files"]
    implementer -->|writes| prodcode["production code"]
    security -->|writes| sectests["security tests"]
    assessor -->|writes| assessdoc
    wrapper -->|writes| contracts
    wrapper -->|writes| lockingtests["tests/locking/"]
    mapper -->|writes| codebase

    %% Read-only commands -> Artifacts
    status -->|reads| state
    status -->|reads| contracts
    changelog -->|reads| summaries
    roadmap -->|reads| roadmapdoc

    %% Orchestrator state updates
    slice -->|updates| state
    slice -->|updates| summaries
    slice -->|updates| decisions
    wrap -->|updates| state
    wrap -->|updates| summaries
```

## Agent Descriptions

### gl-designer
**Spawned by:** `/gl:design`, `/gl:roadmap milestone`
**Tools:** Read, Write, Bash, Glob, Grep, WebSearch, WebFetch, AskUserQuestion

Runs interactive system design sessions. Gathers requirements through conversation, researches technical unknowns via web search, proposes architecture, and gets user sign-off. Produces DESIGN.md, ROADMAP.md, and DECISIONS.md. In brownfield projects, factors in risk tiers from ASSESS.md and [WRAPPED] boundaries from CONTRACTS.md. Supports a lightweight `milestone_planning` mode that skips init phases.

### gl-architect
**Spawned by:** `/gl:init` (Phase 3), `/gl:add-slice`
**Tools:** Read, Write, Bash, Glob, Grep

Converts approved designs into typed contracts, API schemas, and dependency graphs. Reads DESIGN.md and produces CONTRACTS.md + GRAPH.json. Understands [WRAPPED] contracts from brownfield wrapping. Never writes implementation code.

### gl-test-writer
**Spawned by:** `/gl:slice` (Step 1)
**Tools:** Read, Write, Bash, Glob, Grep

Generates integration tests from contracts before any implementation exists. Tests behaviour and outputs, never implementation details. Has awareness of locking test names from wrapped boundaries but cannot see implementation code.

### gl-implementer
**Spawned by:** `/gl:slice` (Step 2)
**Tools:** Read, Write, Edit, Bash, Glob, Grep

Makes failing tests pass. Sees contracts and test names but never test source code. Follows CLAUDE.md engineering standards. Handles deviations (auto-fixes bugs, adds missing functionality, stops for architectural changes). Cannot modify test files.

### gl-security
**Spawned by:** `/gl:slice` (Step 3), `/gl:wrap`, `/gl:ship`, `/gl:assess` (internal)
**Tools:** Read, Write, Bash, Glob, Grep

Reviews implementation diffs for OWASP Top 10 vulnerabilities. Produces failing test cases for each issue found — the implementer fixes them. Never fixes code directly. In `/gl:ship`, runs a full audit across the entire codebase.

### gl-verifier
**Spawned by:** `/gl:slice` (Step 4), `/gl:ship`
**Tools:** Read, Bash, Glob, Grep

Read-only agent that verifies slice completion through goal-backward analysis. Checks that contracts are actually satisfied, not just that tests pass. Can see everything but cannot modify any code.

### gl-assessor
**Spawned by:** `/gl:assess`
**Tools:** Read, Bash, Glob, Grep

Read-only analytical agent for brownfield codebases. Maps test coverage, inventories contracts, identifies risks, compares against CLAUDE.md standards. Produces ASSESS.md with priority-tiered wrap recommendations (Critical/High/Medium risk).

### gl-wrapper
**Spawned by:** `/gl:wrap`
**Tools:** Read, Write, Bash, Glob, Grep

Extracts contracts from existing code boundaries and writes locking tests. **Deliberate isolation exception**: sees implementation code AND writes tests, because locking tests must verify what code currently does. Scoped to `tests/locking/` only. When boundaries are later refactored via `/gl:slice`, locking tests are deleted and proper integration tests replace them.

### gl-codebase-mapper
**Spawned by:** `/gl:map` (4 parallel instances)
**Tools:** Read, Bash, Glob, Grep

Analyses existing codebases before initialisation. Four instances run in parallel with different focus areas (tech stack, architecture, code quality, concerns). Produces structured documentation in `.greenlight/codebase/`.

### gl-debugger
**Spawned by:** `/gl:quick`
**Tools:** Read, Write, Edit, Bash, Grep, Glob

Investigates bugs using the scientific method with hypothesis testing. Reproduces the issue, writes a failing test for the root cause, then optionally fixes the code. Can see everything but cannot modify tests without approval.

## Agent Isolation Matrix

```
                    Can See                      Cannot See              Cannot Do
┌──────────────────┬──────────────────────────┬────────────────────────┬────────────────────────┐
│ gl-designer      │ Requirements, context    │ Implementation code    │ Write code/contracts   │
│ gl-architect     │ Requirements, constraints│ Implementation code    │ Write production code  │
│ gl-test-writer   │ Contracts, test patterns │ Implementation code    │ Write production code  │
│ gl-implementer   │ Contracts, test names    │ Test source code       │ Modify test files      │
│ gl-security      │ Diffs, contracts         │ Test impl details      │ Fix production code    │
│ gl-verifier      │ Everything               │ N/A                    │ Modify any code        │
│ gl-assessor      │ Codebase, tests, stds    │ N/A (read-only)        │ Modify any code        │
│ gl-wrapper       │ Impl code, existing tests│ N/A                    │ Modify production code │
│ gl-debugger      │ Everything               │ N/A                    │ Modify tests (w/o OK)  │
│ gl-codebase-mapper│ Full codebase           │ N/A (read-only)        │ Modify any code        │
└──────────────────┴──────────────────────────┴────────────────────────┴────────────────────────┘
```

## Workflows

### Greenfield (new project)
```
gl:init → gl:design → gl:init Phase 3 → gl:slice (repeat per slice)
  interview   designer     architect      test-writer → implementer → security → verifier
```

### Brownfield (existing codebase)
```
gl:map → gl:assess → gl:wrap (per boundary) → gl:slice (refactor wrapped boundaries)
  mapper   assessor    wrapper + security      full TDD loop (locking tests deleted)
```

### Milestone Planning
```
gl:roadmap milestone → gl-designer (lightweight session, brownfield-aware)
```

### Ship
```
gl:ship → gl-security (full audit) + gl-verifier (goal-backward check)
```

## Data Flow

```
INTERVIEW.md ──→ DESIGN.md ──→ CONTRACTS.md ──→ GRAPH.json
  (gl:init)      (designer)     (architect)      (architect)
                                     │
                                     ▼
                    ┌─── test-writer writes tests
                    │         │
                    │    implementer makes them pass
                    │         │
                    │    security scans diffs
                    │         │
                    │    verifier checks contracts
                    │         │
                    └──→ STATE.md updated ──→ summaries/ ──→ ROADMAP.md
                                                            DECISIONS.md

Brownfield data flow:
  codebase/ ──→ ASSESS.md ──→ CONTRACTS.md ([WRAPPED]) ──→ STATE.md
  (mapper)      (assessor)    (wrapper)                    (wrap progress)
       │                           │
       └───────────────────────────┴──→ designer (risk tiers, wrap awareness)
```
