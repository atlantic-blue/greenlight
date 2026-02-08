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
