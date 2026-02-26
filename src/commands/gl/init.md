---
name: gl:init
description: Initialize project with brief interview and config, or generate contracts from an existing design
allowed-tools: [Read, Write, Bash, Glob, Grep, Task, AskUserQuestion]
---

# Greenlight: Project Initialization

You are the Greenlight orchestrator. Guide the user through project setup.

**Read CLAUDE.md first.** Internalise the engineering standards.
**Read templates/config.md** for config schema.
**Read templates/state.md** for STATE.md lifecycle.

## Route: New Project or Continue

Check for existing state:

```bash
cat .greenlight/config.json 2>/dev/null
cat .greenlight/DESIGN.md 2>/dev/null
```

**If DESIGN.md exists and no contracts yet:** Jump to [Phase 3: Contracts](#phase-3-contracts).
**If contracts exist:** "Project already initialized. Run /gl:slice 1 to start building."
**Otherwise:** Start from Phase 1.

## Phase 1: Interview

Ask the user about their project. Be conversational, not interrogative.

<essential_context>
- What is this? (one sentence value proposition)
- Who uses it? (end users, developers, internal)
- What's the first thing a user should be able to do? (this becomes slice 1)
- What does "done" look like for an MVP? (max 5 user actions)
- What stack? (or "you pick" — recommend based on project type)
- Hard constraints? (specific DB, platform, existing codebase, deployment target)
- Existing code? (if yes, suggest `/gl:map` first)
</essential_context>

Push back if scope is too large. The entire MVP should fit in 5 user actions or fewer. If the user describes more, help them identify the core 5 and defer the rest.

**Extract the thinnest first slice.** Not "set up auth" — that's infrastructure. The first slice is the first thing a user does that delivers value.

## Phase 2: Config

Create `.greenlight/config.json` using the schema from `templates/config.md`:

```bash
mkdir -p .greenlight
```

```json
{
  "version": "1.0.0",
  "mode": "interactive",
  "model_profile": "balanced",
  "model_overrides": {},
  "profiles": {
    "quality": {
      "architect": "opus", "designer": "opus", "test_writer": "opus",
      "implementer": "opus", "security": "opus", "debugger": "opus",
      "verifier": "opus", "codebase_mapper": "opus",
      "assessor": "opus", "wrapper": "opus"
    },
    "balanced": {
      "architect": "opus", "designer": "opus", "test_writer": "sonnet",
      "implementer": "sonnet", "security": "sonnet", "debugger": "sonnet",
      "verifier": "sonnet", "codebase_mapper": "sonnet",
      "assessor": "sonnet", "wrapper": "sonnet"
    },
    "budget": {
      "architect": "sonnet", "designer": "sonnet", "test_writer": "sonnet",
      "implementer": "sonnet", "security": "haiku", "debugger": "sonnet",
      "verifier": "haiku", "codebase_mapper": "haiku",
      "assessor": "haiku", "wrapper": "sonnet"
    }
  },
  "workflow": {
    "security_scan": true,
    "visual_checkpoint": true,
    "auto_parallel": true,
    "max_implementation_retries": 3,
    "max_security_retries": 2,
    "run_full_suite_after_slice": true
  },
  "test": {
    "command": "{detected or chosen}",
    "filter_flag": "{detected or chosen}",
    "coverage_command": "{detected or chosen}",
    "security_filter": "security"
  },
  "project": {
    "name": "{from user}",
    "stack": "{from user}",
    "src_dir": "src",
    "test_dir": "tests"
  }
}
```

Ask: "Want to run in interactive mode (confirm each step) or YOLO mode (auto-approve non-critical steps)?"

**Stack-specific test config detection:**

| Stack | test.command | test.filter_flag |
|-------|-------------|-----------------|
| Node + Vitest | `npx vitest run` | `--filter` |
| Node + Jest | `npx jest` | `--testPathPattern` |
| Python + pytest | `pytest` | `-k` |
| Go | `go test ./...` | `-run` |
| Swift + XCTest | `swift test` | `--filter` |

Save a summary of the interview for the designer:

Write `.greenlight/INTERVIEW.md`:

```markdown
# Project Interview

## Value Proposition
{one sentence}

## Users
{who uses it}

## MVP Scope
{3-5 user actions, priority order}

## Stack
{chosen stack}

## Constraints
{hard constraints}

## Deferred Ideas
{anything the user mentioned but was pushed to post-MVP}
```

### Report

```
Project initialized.

Name: {project name}
Stack: {stack}
MVP: {N} user actions
Mode: {interactive/yolo}

Next: Run /gl:design to start the system design session.
The designer will gather detailed requirements, research technical
decisions, and produce an architecture for your project.
```

## Model Resolution

Before spawning any agent, resolve its model from `.greenlight/config.json`:

1. Check `model_overrides[agent_name]` — if set, use it
2. Else check `profiles[model_profile][agent_name]` — use profile default
3. Else fall back to `sonnet`

Agents spawned by this command: `architect`.

## Phase 3: Contracts

This phase runs when DESIGN.md exists (after `/gl:design`).

Read the design document:

```bash
cat .greenlight/DESIGN.md
```

Spawn the architect with the full design context:

```
Task(prompt="
Read agents/gl-architect.md
Read CLAUDE.md

<project_context>
{value prop, users, MVP scope, stack, constraints from config}
</project_context>

<user_actions>
{3-5 things a user can do, priority order}
</user_actions>

<design>
{full contents of DESIGN.md}
</design>

<stack>
{chosen stack with versions}
</stack>

<existing_code>
{if brownfield: summary from .greenlight/codebase/ docs. Otherwise: 'Greenfield project'}
</existing_code>

Produce:
1. Typed contracts for every boundary in the MVP
2. Dependency graph as GRAPH.json
3. Each slice mapped to user actions it enables

The design document contains requirements, architecture, data model,
API surface, security approach, and locked technical decisions. Use
these as constraints — do not contradict design decisions.
", subagent_type="gl-architect", model="{resolved_model.architect}", description="Generate contracts and dependency graph")
```

## Phase 4: User Review

Present contracts and dependency graph clearly:

1. **Slices in priority order** — "Here's what we build and when"
2. **Dependency graph** — "These can run in parallel, these depend on each other"
3. **Contract summary** — for each contract: boundary, input, output, errors

Ask: "Does this match? Anything missing? Anything wrong?"

Iterate until approved. Each revision = fresh architect agent with updated context.

## Phase 5: Scaffold

Generate project structure adapted to the chosen stack:

```bash
mkdir -p {config.project.src_dir}
mkdir -p {config.project.test_dir}/integration
mkdir -p {config.project.test_dir}/security
mkdir -p {config.project.test_dir}/fixtures
mkdir -p .greenlight
```

### Stack-Specific Setup

**Node/TypeScript:**
- `package.json` with test framework, linter, TypeScript deps
- `tsconfig.json` with strict mode
- Test framework config (vitest.config.ts or jest.config.ts)
- `.eslintrc.json` or `eslint.config.js`
- `.gitignore` (node_modules, dist, .env, coverage)

**Python:**
- `pyproject.toml` with pytest, ruff/flake8, type checker
- `requirements.txt` or `pyproject.toml` dependencies
- `conftest.py` with test fixtures
- `.gitignore` (venv, __pycache__, .env, .coverage)

**General:**
- `.env.example` with placeholder values
- `.gitignore` appropriate for stack
- CI pipeline (GitHub Actions):
  ```yaml
  # .github/workflows/ci.yml
  # lint → test → build on every push
  ```

### Write Contract Files

Write the approved contracts to `.greenlight/CONTRACTS.md`.
Write the dependency graph to `.greenlight/GRAPH.json`.

### Test Fixture Scaffolding

Create base factory file:
```
{config.project.test_dir}/fixtures/factories.{ext}
```

Create test setup file:
```
{config.project.test_dir}/fixtures/setup.{ext}
```

### Commit

```bash
git init  # if not already a git repo
git add package.json tsconfig.json .gitignore .github/ .greenlight/ tests/fixtures/ # etc.
git commit -m "chore: greenlight scaffold

Contracts: {N} defined
Slices: {N} in dependency graph
Stack: {stack description}
CI: GitHub Actions (lint → test → build)
"
```

## Phase 6: Ready State

Write `.greenlight/STATE.md` using the template from `templates/state.md`:

```markdown
# Project State

## Overview
{one line value prop}
Stack: {stack}
Mode: {interactive/yolo}

## Slices

| ID | Name | Status | Tests | Security | Deps |
|----|------|--------|-------|----------|------|
{for each slice from GRAPH.json}

Progress: [░░░░░░░░░░] 0/{N} slices

## Current

Slice: 1 — {first slice name}
Step: pending
Last activity: {today's date} — Project initialized

## Test Summary

Total: 0 passing, 0 failing, 0 security
Last run: never

## Decisions

{decisions from DESIGN.md}

## Blockers

None

## Session

Last session: {now}
Resume file: None
```

### Report

```
Greenlight initialized.

Contracts: {N} defined
Slices: {N} in dependency graph
Tests: scaffold ready, 0 written
CI: configured
Mode: {interactive/yolo}

Run /gl:slice 1 to start: "{first slice name}"
```

## Phase 6a: Parallel State (File-Per-Slice Format)

After the legacy STATE.md is written (Phase 6), also create the file-per-slice state format.
This enables per-slice state tracking with individual slice files.

### Create .greenlight/slices/ Directory

Create the `.greenlight/slices/` directory with permissions 0o755:

```bash
mkdir -m 0755 -p .greenlight/slices/
```

**Error handling:**

| Error State | When | Behaviour |
|-------------|------|-----------|
| SlicesDirExists | `.greenlight/slices/` already exists | Warn user. Offer to overwrite or skip. Do not silently overwrite |
| DirectoryCreateFailure | Cannot create `.greenlight/slices/` directory | Report error. Abort init. Suggest checking permissions |

### Create Individual Slice Files from GRAPH.json

For each slice defined in GRAPH.json, create a corresponding file at
`.greenlight/slices/{id}.md`.

**Slice ID validation:** Each slice ID must match the pattern `S-{digits}` (e.g. S-1, S-29).
IDs not matching `S-{digits}` are invalid — skip them with a warning:
"Skipping invalid slice ID: {id}". This validation also prevents path traversal attacks.

**Slice file frontmatter schema:**

```yaml
---
id: {slice_id}
status: pending
step: none
milestone: {milestone_name}
started: ""
updated: {current_ISO_timestamp}
tests: 0
security_tests: 0
session: ""
deps: {deps_from_GRAPH_json}
---
```

Fields:
- `id` — slice identifier (e.g. S-1)
- `status` — initial value: `pending`
- `step` — initial value: `none`
- `milestone` — milestone name from the design phase
- `started` — empty string on init
- `updated` — current ISO timestamp
- `tests` — integer, initial value: `0`
- `security_tests` — integer, initial value: `0`
- `session` — empty string on init
- `deps` — dependency list from GRAPH.json

**Slice file body:**

```markdown
# {Slice Name}

## Summary

## Notes
```

**Crash safety (atomic writes):** Use write-to-temp-then-rename for every slice file.
Write to `.greenlight/slices/{id}.md.tmp` first, then rename to the final path.
This ensures no partial state on crash (NFR-4).

**File permissions:** Set each slice file to 0o644 after writing.

**Error handling:**

| Error State | When | Behaviour |
|-------------|------|-----------|
| InvalidSliceIdInGraph | GRAPH.json contains slice ID not matching `S-{digits}` | Skip invalid slice. Warn: "Skipping invalid slice ID: {id}" |
| FileWriteFailure | Cannot write a slice file | Report error for that file. Continue with remaining slices. Warn user of partial state |

### Generate STATE.md as Summary View

After all individual slice files are written, regenerate `.greenlight/STATE.md` as a
summary computed from all slice files. Add the header comment:

```markdown
<!-- GENERATED by greenlight -->
```

The GENERATED comment signals that STATE.md is derived from the slice files, not the
source of truth. The format follows DESIGN.md section 4.6.

### Invariants

- Directory permissions are 0o755, file permissions are 0o644
- Each slice file follows the schema in `templates/slice-state.md` (C-76)
- Slice IDs are validated before file creation (path traversal prevention)
- All slice files are created atomically (write-to-temp-then-rename)
- STATE.md is generated after all slice files are written
- Existing /gl:init behaviour for non-state operations is unchanged
- If /gl:init was already run with legacy format, this does NOT auto-migrate
  existing state. Legacy format and file-per-slice format coexist; state detection
  (see below) determines which format is active.

## Phase 6b: Project State File

Create `.greenlight/project-state.json` to store non-slice project state.

**Error handling:**

| Error State | When | Behaviour |
|-------------|------|-----------|
| ProjectStateExists | `.greenlight/project-state.json` already exists | Overwrite with fresh state. Warn user |
| WriteFailure | Cannot write `project-state.json` | Report error. Abort |
| MissingDesignContext | Design phase did not provide `value_prop` or `stack` | Use placeholder values |

**Schema:**

```json
{
  "overview": {
    "value_prop": "{from design phase}",
    "stack": "{from design phase}",
    "mode": "{from config.json or default yolo}"
  },
  "session": {
    "last_session": "{current ISO timestamp}",
    "resume_file": null
  },
  "blockers": []
}
```

Fields:
- `overview.value_prop` — one-line value proposition from the design phase
- `overview.stack` — technology stack from the design phase
- `overview.mode` — run mode; defaults to `yolo` if not specified in config
- `session.last_session` — current ISO timestamp; always a valid ISO timestamp
- `session.resume_file` — `null` on init
- `blockers` — always an array (never null); initialises as an empty array `"blockers": []`

**Invariants:**
- `overview.mode` defaults to `"yolo"` if not specified in config.json
- `blockers` is always an array (never null)
- `session.last_session` is always a valid ISO timestamp

## State Detection

State detection resolves which format is active for all commands. The detection flow
is documented as the single source of truth in `references/state-format.md`.

Detection is read-only — it does not modify any files.

### Detection Flow

1. Check if `.greenlight/slices/` directory exists:
   - **Yes** → file-per-slice format. Read individual slice files from `.greenlight/slices/`.
   - **No** → continue to step 2.

2. Check if `.greenlight/STATE.md` exists:
   - **Yes** → legacy format. Parse STATE.md as the source of truth.
   - **No** → continue to step 3.

3. No state found:
   - Return `NoStateFound` error.
   - Display: "No project state found. Run /gl:init to get started."

**Format values:** `"file-per-slice"` | `"legacy"` | `"none"`

**Error handling:**

| Error State | When | Behaviour |
|-------------|------|-----------|
| NoStateFound | Neither `.greenlight/slices/` nor `.greenlight/STATE.md` exists | Display suggestion to run /gl:init |
| SlicesDirUnreadable | `.greenlight/slices/` exists but cannot be read | Report permission error |
| LegacyStateMalformed | STATE.md exists but cannot be parsed | Report parse error |

**Invariants:**
- Detection is identical across all commands (single source of truth in `state-format.md`)
- The `.greenlight/slices/` directory existence check is the primary signal
- Detection is read-only (does not modify files)
