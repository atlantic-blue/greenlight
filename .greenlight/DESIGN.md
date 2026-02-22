# System Design: File-Per-Slice State (Parallel Session Safety)

> **Project:** Greenlight
> **Scope:** Fix concurrent session state corruption by replacing the single STATE.md with per-slice state files. Each session writes only to its own slice's file, eliminating write conflicts entirely. STATE.md becomes a generated summary view.
> **Stack:** Go 1.24, stdlib only. Deliverables are embedded markdown content plus manifest entries.
> **Date:** 2026-02-22
> **Replaces:** Previous DESIGN.md (verification-tiers -- complete, 710 tests passing)

---

## 1. Problem Statement

When multiple Claude Code sessions run `/gl:slice` in parallel, they all perform read-modify-write cycles on a single STATE.md file without coordination. Last write wins, previous updates are lost.

**Initial failure mode (Session 1):** Session A starts S2.5, Session B starts S2.6, Session C starts S3.1 -- all reading STATE.md before any writes complete. Sessions write back in arbitrary order. Result: only the last writer's state survives. S3.1 was actively being built (test file existed on disk) but was not listed in STATE.md's Parallel section. The orchestrator in another session had no way to know S3.1 was taken, risking duplicate work or conflicts.

### Confirmed Failure Modes (Session 2, 2026-02-22)

A second investigation with 4 concurrent terminals (S2.5, S2.6, S2.7, S3.1; later S2.8, S3.4) confirmed the bug is ongoing and identified failure modes beyond lost writes:

**1. Invisible work.** Sessions have no mechanism to signal they've started a slice before their first commit. A session mid-slice writing tests has produced zero commits. The only evidence it exists is a running process in another terminal. There is no file, no lock, no signal. Other sessions cannot determine which slices are taken.

**2. False availability.** Slices that are actively being worked on appear as `pending` in STATE.md, inviting duplicate work. S3.4 (ClassificationBreakdown) had test and implementation commits in git (`40ce481`, `0b2ab6f`) but STATE.md showed it as `pending` with `0/0` tests because the session running S3.4 hadn't written back yet.

**3. Stale pointers.** The Current/Parallel sections point to the wrong terminal. STATE.md showed `Current: S2.8 — Recording completion` in a terminal that was idle. The session actually working on S2.8 was in a different terminal. The Parallel section read "None" throughout despite 4 active sessions -- every session overwrites the entire file, sets itself as Current, and doesn't know about other sessions.

**4. Status field corruption.** S3.2 (SnoreScoreGauge) had full git history through security fix and docs commit (`ebeb1a5`), but STATE.md showed it as `pending` with `0/0` tests. A later write from the S3.2 session corrected it, but during the window between completion and that write, any other session reading STATE.md would see S3.2 as unstarted.

**5. Three-source inconsistency.** At any given moment, STATE.md, git log, and terminal reality all disagree:

| Source | Shows S3.4 as | Shows S2.8 as |
|--------|--------------|--------------|
| STATE.md | pending (wrong) | current/tests (partially right) |
| Git log | in progress -- has test + impl commits | pending -- no commits |
| Reality | in progress in terminal 2 | in progress in terminal 1 |

No single source gives the correct picture. STATE.md lies due to write races. Git only shows committed work. The actual terminal processes are invisible to each other.

### Root Cause

STATE.md is a single mutable file used as the source of truth for all slice state. Multiple concurrent sessions perform uncoordinated read-modify-write cycles on it. There is no file locking, no atomic updates, and no conflict detection.

### Why File-Per-Slice Solves This

If each slice's state lives in its own file (`S-28.md`), sessions writing to different slices never touch the same file. The race condition is structurally impossible -- not prevented by locks, but eliminated by design.

The file-per-slice solution addresses all five failure modes:

| Failure Mode | Fix |
|---|---|
| Invisible work | Session writes `status: in_progress` + `session` field to slice file **immediately** on `/gl:slice` start, before any agent work |
| False availability | Each session owns its slice file -- no other session can overwrite its status to `pending` |
| Stale pointers | No shared Current/Parallel sections. Active slices derived from individual file status fields |
| Status field corruption | One file per slice, one writer per file. No read-modify-write on a shared monolith |
| Three-source inconsistency | Slice file becomes the single source of truth. Written before first commit, updated throughout |

---

## 2. Requirements

### 2.1 Functional Requirements

**FR-1: File-per-slice state storage.** Each slice gets its own state file at `.greenlight/slices/{slice-id}.md`. The file is the authoritative source of truth for that slice's status, step, test counts, timestamps, decisions, and files touched.

**FR-2: Conflict-free concurrent writes.** Multiple sessions running `/gl:slice` in parallel write only to their own slice's file. No shared mutable file for slice state. The concurrent write conflict is eliminated by design, not by locking.

**FR-3: Generated summary view.** `/gl:status` reads all slice files from `.greenlight/slices/` and computes a summary. STATE.md becomes generated output (marked with a comment header), not a source of truth.

**FR-4: Project-level state separation.** Non-slice state (session metadata, active blockers, project overview) lives in `.greenlight/project-state.json`. This file is read/written by at most one session at a time (the orchestrator).

**FR-5: Backward compatibility.** The system detects which state format is active:
- `.greenlight/slices/` directory exists -> file-per-slice format
- `.greenlight/STATE.md` exists (no slices/) -> legacy format
- Neither exists -> no state, suggest `/gl:init`
Old projects continue working with STATE.md until explicitly migrated.

**FR-6: Migration command.** `/gl:migrate-state` converts an existing STATE.md project to file-per-slice format. Migration is one-way, creates a backup of STATE.md, and is all-or-nothing.

**FR-7: Slice file self-documentation.** Each slice file contains structured frontmatter (machine-readable) and markdown body sections (human-readable): Why, What, Dependencies, Contracts, Decisions, Files. The file is a living record of the slice, not just status tracking.

**FR-8: State abstraction in all commands.** All state-reading commands check format detection before reading. Commands work identically regardless of which format is active. Detection logic is documented in a reference file so all agents implement it consistently.

**FR-9: Advisory session tracking.** When a session starts working on a slice, it records a session identifier (ISO timestamp + random suffix) in the slice file's frontmatter `session` field. Other sessions reading the slices directory can see which slices are actively being worked on and warn before claiming a slice in progress.

**FR-10: STATE.md regeneration.** STATE.md is regenerated after every state write operation, keeping it in sync as a convenience view without requiring explicit generation commands.

### 2.2 Non-Functional Requirements

**NFR-1: No file locking required.** The design eliminates the need for OS-level file locks by ensuring each session writes only to its own slice file. The only coordination is advisory (session tracking in frontmatter).

**NFR-2: Performance.** Reading 50+ slice files for `/gl:status` must complete in under 1 second. Go's `os.ReadDir` plus sequential file reads handles this easily for the expected scale (10-50 slices per milestone).

**NFR-3: Zero external dependencies.** All Go code changes must use stdlib only. Frontmatter is parsed by Claude Code agents (which read markdown natively), not by the Go binary. The flat key-value frontmatter format is trivially parseable with line-by-line string splitting if the Go CLI ever needs to parse slice files (e.g., `doctor` command).

**NFR-4: Crash safety.** Writes to slice files use write-to-temp-then-rename pattern. Temp files are created in `.greenlight/slices/` (same directory, same filesystem) to guarantee atomic rename on POSIX systems.

### 2.3 Constraints

- Go 1.24 stdlib only (no external dependencies)
- This is primarily embedded content changes (markdown in `src/`), not Go library code
- Must be backward compatible with existing projects using STATE.md
- Agents (Claude Code sessions) parse the files, not the Go binary
- Must work on macOS and Linux (POSIX `rename` atomicity)
- Must not break existing 710 tests

### 2.4 Out of Scope

- OS-level file locking mechanisms (design eliminates the need)
- Real-time state synchronization between sessions
- Database-backed state storage
- Distributed state across machines
- Dual-write transition period (detect-and-migrate is sufficient)
- Windows-specific atomicity guarantees
- Automatic conflict resolution if two sessions claim the same slice
- Slice archival for completed milestones

---

## 3. Technical Decisions

| # | Decision | Chosen | Rejected | Rationale |
|---|----------|--------|----------|-----------|
| D-30 | Frontmatter format | Flat key-value between `---` delimiters (YAML-like) | Full YAML (requires external dep `gopkg.in/yaml.v3`); JSON frontmatter (less readable in markdown) | Go stdlib has no YAML parser. Flat key-value is parseable by both Claude (reads markdown natively) and Go (line-by-line string split). All needed fields are flat -- no nesting required for slice metadata. |
| D-31 | State detection strategy | Check for `.greenlight/slices/` directory existence | Version field in config.json; Magic comment in STATE.md | Directory existence is the simplest, most reliable signal. No config migration needed. Works even if config.json is missing or corrupt. |
| D-32 | Migration approach | One-way explicit `/gl:migrate-state` command | Automatic migration on first access; Dual-write transition period | Explicit migration is safer -- user controls when it happens. Automatic migration risks corrupting state during a critical operation. Dual-write adds complexity with minimal benefit since detect-and-migrate is sufficient. |
| D-33 | Session tracking | Advisory: ISO timestamp + random suffix in frontmatter `session` field | No tracking (blind); OS-level file locks; PID-based tracking | Advisory tracking lets other sessions warn without blocking. PID is meaningless across machines or after crashes. File locks are fragile and leave stale locks after crashes. Timestamp + random suffix is unique enough and human-readable. |
| D-34 | STATE.md regeneration trigger | After every state write operation | Only on `/gl:status` (requires explicit generation); Never (deprecate entirely) | Keeps STATE.md in sync as a convenience view. Minimal cost (read all slice files + one write per state change). Humans and tools that grep STATE.md continue working without behavior changes. |
| D-35 | project-state.json contents | Session metadata + active blockers + project overview | Session metadata only (too minimal); Everything from STATE.md (duplicates slice data) | Session metadata and blockers are the only non-slice state that changes during execution. Overview (value prop, stack, mode) is stable context that belongs here. Decisions have their own file (DECISIONS.md). |
| D-36 | Slice file naming | `{slice-id}.md` (e.g., `S-28.md`) | Slugified name (`file-per-slice-state.md`); Sequential number; UUID | Slice ID is already unique, human-readable, and matches GRAPH.json entries. No translation or lookup needed. |
| D-37 | Backward compatibility duration | Indefinite (both formats supported forever) | Deprecation timeline; Force migration after N versions | No cost to supporting both formats. Detection is a single directory existence check. Removing legacy support would break existing projects for no benefit. |
| D-38 | Dual-write period | No dual-write | Optional dual-write safety net during migration | Dual-write adds complexity and introduces its own bugs (e.g., partial writes to one format). Clean cutover via `/gl:migrate-state` with backup is simpler and more reliable. |

---

## 4. Architecture

### 4.1 New File Structure

```
.greenlight/
  slices/                    # NEW: per-slice state directory
    S-01.md                  # One file per slice
    S-02.md
    ...
    S-28.md
  project-state.json         # NEW: non-slice state
  STATE.md                   # CHANGED: now generated, not source of truth
  GRAPH.json                 # Unchanged
  CONTRACTS.md               # Unchanged
  config.json                # Unchanged
  DESIGN.md                  # Unchanged
  ROADMAP.md                 # Unchanged
  DECISIONS.md               # Unchanged
  summaries/                 # Unchanged
```

### 4.2 Slice State File Schema

```markdown
---
id: S-28
status: implementing
step: security
milestone: parallel-state
started: 2026-02-22
updated: 2026-02-22T14:30:00Z
tests: 12
security_tests: 2
session: 2026-02-22T14:00:00Z-a7f3
deps: S-26,S-27
---

# S-28: File-per-slice state storage

## Why
Concurrent sessions corrupt STATE.md because of uncoordinated read-modify-write cycles.

## What
Each slice gets its own state file. Sessions write only to their own slice file, eliminating write conflicts by design.

## Dependencies
- S-26: Verification tier documentation (complete)
- S-27: Architect integration (complete)

## Contracts
- C-50: SliceStateReader
- C-51: SliceStateWriter

## Decisions
- 2026-02-22: Used flat key-value frontmatter instead of nested YAML to maintain zero external dependencies

## Files
- src/templates/slice-state.md (new)
- src/references/state-format.md (new)
- src/commands/gl/migrate-state.md (new)
```

### 4.3 Frontmatter Field Definitions

| Field | Type | Values | Required | Description |
|-------|------|--------|----------|-------------|
| `id` | string | `S-{N}` or `S-{NN}` | yes | Slice identifier, matches GRAPH.json |
| `status` | enum | `pending`, `ready`, `tests`, `implementing`, `security`, `fixing`, `verifying`, `complete` | yes | Current slice status |
| `step` | string | `none`, `tests`, `implementing`, `security`, `fixing`, `verifying`, `complete` | yes | Current step within the TDD loop |
| `milestone` | string | milestone slug | yes | Which milestone this slice belongs to |
| `started` | ISO date | `YYYY-MM-DD` or empty | no | When work began on this slice |
| `updated` | ISO timestamp | `YYYY-MM-DDTHH:MM:SSZ` | yes | Last modification time |
| `tests` | integer | >= 0 | yes | Number of passing functional tests |
| `security_tests` | integer | >= 0 | yes | Number of passing security tests |
| `session` | string | `{ISO-timestamp}-{random}` or empty | no | Advisory: which session is actively working on this slice |
| `deps` | comma-separated | `S-01,S-02` or empty | no | Slice dependencies (references other slice IDs) |

### 4.4 project-state.json Schema

```json
{
  "overview": {
    "value_prop": "TDD-first development system for Claude Code",
    "stack": "Go 1.24 (stdlib only)",
    "mode": "yolo"
  },
  "session": {
    "last_session": "2026-02-22T14:30:00Z",
    "resume_file": null
  },
  "blockers": []
}
```

### 4.5 State Detection Logic

Every command that reads state must follow this detection flow:

```
ReadState():
  if directoryExists(".greenlight/slices/"):
    return readSliceFiles(".greenlight/slices/")
  else if fileExists(".greenlight/STATE.md"):
    return parseLegacyState(".greenlight/STATE.md")
  else:
    return NoStateError("Run /gl:init to get started")
```

This logic is documented once in `references/state-format.md` and referenced by all commands.

### 4.6 Generated STATE.md Format

When using file-per-slice format, STATE.md is regenerated after every state write:

```markdown
<!-- GENERATED by greenlight -- source of truth is .greenlight/slices/*.md -->
<!-- Do not edit this file directly. Changes will be overwritten. -->
# Project State

## Overview
[from project-state.json overview section]
Stack: [from project-state.json]
Mode: [from project-state.json]

## Slices

| ID | Name | Status | Tests | Security | Deps |
|----|------|--------|-------|----------|------|
[computed by reading all .greenlight/slices/*.md frontmatter]

Progress: [computed progress bar] done/total slices

## Current

[list of slices where status is not pending and not complete, with their current step]

## Test Summary

Total: N passing, N failing, N security
[sum of tests and security_tests from all slice files]

## Blockers

[from project-state.json blockers array]

## Session

Last session: [from project-state.json]
Resume file: [from project-state.json]
```

### 4.7 Migration Flow (/gl:migrate-state)

```
1. Verify .greenlight/STATE.md exists
2. Verify .greenlight/slices/ does NOT exist (prevent double migration)
3. Parse STATE.md:
   a. Extract slice table rows (ID, Name, Status, Tests, Security, Deps)
   b. Extract Current section (active slice, step)
   c. Extract Decisions section
   d. Extract Blockers section
   e. Extract Session section
   f. Extract Overview section (value prop, stack, mode)
4. Create .greenlight/slices/ directory
5. For each slice row:
   a. Create .greenlight/slices/{id}.md with frontmatter from table data
   b. Populate minimal body sections (name in heading, deps listed)
   c. If this is the current slice, set step from Current section
6. Create .greenlight/project-state.json from non-slice sections
7. Rename .greenlight/STATE.md to .greenlight/STATE.md.backup
8. Generate new .greenlight/STATE.md (generated format with header comment)
9. Report: "Migrated N slices to file-per-slice format. Backup: STATE.md.backup"
```

### 4.8 Command Impact Matrix

| Command | Current Reads | New Reads | Current Writes | New Writes |
|---------|--------------|-----------|----------------|------------|
| /gl:slice | STATE.md (pre-flight, Step 4, Step 10) | Own slice file + project-state.json | STATE.md (Step 4, Step 10) | Own slice file + regenerate STATE.md |
| /gl:status | STATE.md (full read) | All slice files + project-state.json | Never | Regenerate STATE.md |
| /gl:pause | STATE.md (current section) | Own slice file + project-state.json | STATE.md (session section) | Own slice file + project-state.json |
| /gl:resume | STATE.md (full read) | All slice files + project-state.json | STATE.md (session section) | Own slice file + project-state.json |
| /gl:ship | STATE.md (pre-check: all complete?) | All slice files | Never | Never |
| /gl:init | N/A (creates STATE.md) | N/A (creates) | Creates STATE.md | Create slices/ dir + slice files + project-state.json |
| /gl:add-slice | STATE.md (slice list) | All slice files | STATE.md (add row) | Create new slice file + regenerate STATE.md |
| /gl:quick | STATE.md (test summary) | Relevant slice file | STATE.md (test counts) | Relevant slice file + regenerate STATE.md |
| /gl:migrate-state | N/A (new command) | STATE.md (parse) | N/A (new) | Create slices/ + slice files + project-state.json |

---

## 5. File Changes in Codebase

### 5.1 New Embedded Content (3 files)

| File | Purpose | Est. Lines |
|------|---------|-----------|
| `src/templates/slice-state.md` | Template and schema for per-slice state files. Documents frontmatter fields, body sections, lifecycle, status values, and examples. | ~120 |
| `src/references/state-format.md` | Reference doc for state format detection, migration protocol, backward compatibility rules, concurrent access patterns, and advisory session tracking. | ~100 |
| `src/commands/gl/migrate-state.md` | Command definition for `/gl:migrate-state`. Parses legacy STATE.md, creates slices/ directory, writes individual slice files, creates project-state.json, backs up STATE.md. | ~80 |

### 5.2 Modified Embedded Content (11 files)

| File | Change | Est. Lines Changed |
|------|--------|--------------------|
| `src/commands/gl/slice.md` | Update pre-flight (state detection), Step 4 (write to slice file), Step 10 (write to slice file + regenerate STATE.md). Add state format detection at start. | ~40 |
| `src/commands/gl/status.md` | Read from slices/ directory instead of STATE.md. Compute summary from all slice files. Generate STATE.md. Add state detection fallback. | ~30 |
| `src/commands/gl/pause.md` | Write session info to slice file and project-state.json instead of STATE.md. | ~15 |
| `src/commands/gl/resume.md` | Read from slice files and project-state.json instead of STATE.md. State detection for both formats. | ~20 |
| `src/commands/gl/ship.md` | Read all slice files for pre-check instead of STATE.md. State detection fallback. | ~15 |
| `src/commands/gl/init.md` | Phase 6: create slices/ directory, write individual slice files, create project-state.json instead of single STATE.md. | ~30 |
| `src/commands/gl/add-slice.md` | Step 6: create new slice file instead of updating STATE.md row. Regenerate STATE.md. | ~15 |
| `src/commands/gl/quick.md` | Update test counts in relevant slice file instead of STATE.md. Regenerate STATE.md. | ~10 |
| `src/templates/state.md` | Document both formats. Explain generated nature of STATE.md in file-per-slice mode. Add migration instructions. | ~30 |
| `src/CLAUDE.md` | Add state format awareness rule: "Check `.greenlight/slices/` before reading STATE.md directly." | ~5 |
| `src/references/checkpoint-protocol.md` | Update checkpoint save/restore to reference slice files for state context. | ~10 |

### 5.3 Go Code Changes (1 file)

| File | Change |
|------|--------|
| `internal/installer/installer.go` | Add 3 entries to Manifest: `"templates/slice-state.md"`, `"references/state-format.md"`, `"commands/gl/migrate-state.md"` |

### 5.4 No Change Required

- `main.go` -- existing `go:embed` glob patterns (`src/templates/*.md`, `src/references/*.md`, `src/commands/gl/*.md`) already cover the new file paths.
- `internal/cli/cli.go` -- no new subcommands in the Go binary.
- Agent definitions (`src/agents/*.md`) -- agents read state through commands, not directly. No agent file changes needed.

**Total new content:** ~300 lines of markdown across 3 new files, ~220 lines modified across 11 existing files, 3 lines added to Go manifest.

---

## 6. Data Model

### 6.1 Entities

**SliceState** (one per slice, stored as `.greenlight/slices/{id}.md`)
- id: string (S-{N})
- status: enum (pending, ready, tests, implementing, security, fixing, verifying, complete)
- step: string (current position in TDD loop)
- milestone: string (milestone slug)
- started: date (optional, when work began)
- updated: timestamp (last modification)
- tests: integer (passing functional test count)
- security_tests: integer (passing security test count)
- session: string (optional, advisory session identifier)
- deps: string[] (slice ID references)
- body: markdown sections (why, what, dependencies, contracts, decisions, files)

**ProjectState** (singleton, stored as `.greenlight/project-state.json`)
- overview: { value_prop: string, stack: string, mode: string }
- session: { last_session: timestamp, resume_file: string|null }
- blockers: string[]

### 6.2 Relationships

- SliceState.deps references other SliceState.id values (dependency graph)
- SliceState.milestone groups slices into milestones
- ProjectState is independent of SliceState (no foreign keys between them)
- GRAPH.json remains the canonical dependency graph; slice file `deps` field is denormalized for convenience

---

## 7. Security

**Input validation on slice IDs:** Slice IDs must match pattern `S-{digits}` or `S-{digits}.{digits}`. File paths are constructed from validated IDs only. This prevents path traversal via malicious slice IDs (e.g., `S-../../etc/passwd`).

**File permissions:** Slice files and project-state.json follow existing conventions: directories `0o755`, files `0o644`.

**No sensitive data:** Slice files contain only project structure information -- status, test counts, file lists. No secrets, tokens, or PII.

**Session tracking is advisory only:** The `session` field is informational. It does not provide access control or locking. A session cannot prevent another session from writing to a slice file -- it can only warn. This is intentional: blocking locks leave stale state after crashes.

**Migration safety:** `/gl:migrate-state` creates a backup before modifying anything. If migration fails partway through, the backup preserves the original STATE.md. The slices/ directory is created atomically (all files written, then STATE.md renamed to backup).

---

## 8. Deployment

No deployment changes. This is embedded content installed by the CLI binary. The binary size increases negligibly (3 new markdown files, ~300 lines total). No new Go packages, no new imports, no new binary dependencies.

The manifest grows from 33 to 36 entries. The `go:embed` directive in `main.go` does not need updating because existing glob patterns already cover the new file locations.

---

## 9. Proposed Build Order

These are logical groupings for the architect to refine into slices with contracts:

1. **Slice state template and reference** -- Create `templates/slice-state.md` (schema, lifecycle, examples) and `references/state-format.md` (detection logic, concurrent access patterns, backward compatibility). These are the foundation documents that all other changes reference.

2. **State detection and /gl:init update** -- Update `/gl:init` to create slices/ directory and individual slice files instead of single STATE.md. Implement state format detection logic. Create project-state.json.

3. **Core command updates (/gl:slice)** -- Update `/gl:slice` to read from and write to individual slice files. Add session tracking on slice claim. Regenerate STATE.md after writes.

4. **Supporting command updates** -- Update `/gl:status`, `/gl:pause`, `/gl:resume`, `/gl:ship`, `/gl:add-slice`, `/gl:quick` to use state detection and file-per-slice reads/writes.

5. **Migration command** -- Create `/gl:migrate-state` command. Parse legacy STATE.md, create slice files, create project-state.json, backup STATE.md.

6. **Documentation and CLAUDE.md** -- Update `src/CLAUDE.md` with state format awareness rule. Update `templates/state.md` to document both formats. Update `references/checkpoint-protocol.md`.

7. **Manifest and integration** -- Add 3 entries to Go manifest. Verify all embedded content is installed correctly. End-to-end verification.

---

## 10. Deferred

| Item | Why Deferred | When to Revisit |
|------|-------------|-----------------|
| OS-level file locking | Design eliminates the need by ensuring each session writes to its own file. Advisory session tracking is sufficient for 1-3 concurrent sessions. | If users report frequent same-slice contention despite advisory warnings. |
| Slice file validation in `doctor` command | The `doctor` command is in the cli-hardening milestone (pending). Slice file validation can be added as a doctor check there. | When cli-hardening milestone begins. |
| Automatic conflict resolution | If two sessions accidentally claim the same slice, they get a warning. Automatic merge/rebase of slice state is premature complexity. | When users report this as a frequent problem. |
| Slice archival | Completed slices accumulate in `.greenlight/slices/`. Archiving old milestone slices to a subdirectory could reduce clutter. | When projects exceed 100 slice files and users report navigation difficulty. |
| Windows atomicity | `os.Rename` is not atomic on Windows. POSIX systems (macOS, Linux) are the primary target. | If Windows becomes a supported platform. |
| Nested frontmatter (full YAML) | Flat key-value is sufficient for current fields. If future fields need nesting, would need either custom parser or external dependency. | If slice metadata requirements grow beyond flat key-value. |
| Cross-session messaging | Sessions can see each other's slice state but cannot send messages. A notification mechanism (e.g., `.greenlight/messages/`) could help coordination. | If users report needing inter-session communication beyond advisory warnings. |

---

## 11. User Decisions (Locked)

| # | Gray Area | Decision | Rationale |
|---|-----------|----------|-----------|
| 1 | Dual-write period | No dual-write. Detect-and-migrate only. | Dual-write adds complexity and its own bugs (partial writes). Clean cutover via `/gl:migrate-state` with backup is simpler and more reliable. |
| 2 | Session tracking approach | Advisory: ISO timestamp + random suffix in frontmatter. Warn, don't block. | Blocking locks are fragile -- stale locks after crashes require manual cleanup. Advisory warnings give developers enough information to coordinate without the risk of deadlocks. |
| 3 | project-state.json scope | Session metadata + active blockers + project overview (value prop, stack, mode). | These are the only non-slice state fields. Decisions already have DECISIONS.md. Slice data lives in slice files. |
| 4 | STATE.md regeneration | After every state write operation. | Keeps convenience view in sync. Minimal cost. Humans and tools that grep STATE.md continue working. |
| 5 | Frontmatter format | Flat key-value between `---` delimiters (YAML-like but no nesting). | Go stdlib has no YAML parser. Flat format is parseable by both Claude (markdown reader) and Go (string split). All current fields are flat. |
| 6 | Migration trigger | Explicit `/gl:migrate-state` command, not automatic. | User controls when state format changes. Prevents surprise migrations during critical slice operations. |
| 7 | Backward compatibility | Indefinite support for both formats. | No cost to detection (one directory check). Removing legacy support would break existing projects for zero benefit. |
| 8 | Slice file naming convention | `{slice-id}.md` directly (e.g., `S-28.md`). | Slice ID is already unique, human-readable, and matches GRAPH.json. No translation layer needed. |
