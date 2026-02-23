# System Design: Greenlight CLI with Parallel Execution

> **Project:** Greenlight
> **Scope:** Extend the Go CLI binary into a full `gl` orchestrator that runs `/gl:slice` sessions autonomously, executes multiple slices in parallel via tmux, provides local commands (`gl status`, `gl roadmap`, `gl changelog`) that work without Claude, and handles interactive commands (`gl init`, `gl design`) that need user input.
> **Stack:** Go 1.24, stdlib only. New internal packages for frontmatter parsing, state reading, tmux management, and process spawning.
> **Date:** 2026-02-22
> **Replaces:** Previous DESIGN.md (parallel-state -- complete, 1050 tests passing)

---

## 1. Problem Statement

Running parallel slices today requires manually opening terminals, picking slices, and tracking what's running. Parallelism should be a first-class Greenlight feature, not manual coordination.

### Root Cause

There is no CLI orchestrator that can:
- Detect which slices are ready (status=pending, all deps=complete)
- Spawn Claude processes headlessly
- Manage tmux sessions for parallel execution
- Monitor progress and auto-refill completed slots

Users manually run `/gl:slice` in separate terminals, manually check GRAPH.json for dependencies, and manually track what's running.

### What This Design Solves

The `gl` CLI becomes the single entry point for all operations -- both from the shell and from inside Claude. The CLI handles orchestration, state reading, and tmux management. Claude skills (`/gl:slice`, `/gl:status`, etc.) call the CLI under the hood.

---

## 2. Requirements

### 2.1 Functional Requirements

**FR-1: CLI command dispatch.** The existing Go binary gains new subcommands (`slice`, `status`, `help`, `roadmap`, `changelog`, `init`, `design`, etc.) alongside existing ones (`install`, `check`, `uninstall`, `version`).

**FR-2: Three command categories.** Every command falls into one of three execution modes:

| Mode | What happens | Needs Claude? | Examples |
|------|-------------|---------------|----------|
| **Interactive** | Runs claude in conversational mode -- needs human input | Yes (interactive) | `gl init`, `gl design` |
| **Autonomous** | Runs claude headlessly -- no human input needed | Yes (headless) | `gl slice`, `gl quick`, `gl wrap`, `gl assess` |
| **Local** | CLI handles it directly -- reads files, prints output | No | `gl status`, `gl roadmap`, `gl changelog`, `gl help` |

**FR-3: Context detection.** `$CLAUDE_CODE` env var distinguishes execution context:
- Set -> running inside Claude -> commands behave as agent tools
- Unset -> running from shell -> commands behave as CLI tools

**FR-4: Parallel slice execution.** When 2+ slices are ready and tmux is available, create a tmux session with tiled panes, one per independent slice.

**FR-5: Watch mode.** `gl slice --watch` polls every 30s, detects completed slices, and auto-fills empty slots with newly ready slices.

**FR-6: Sequential fallback.** When tmux is unavailable, run slices one at a time, re-scanning for ready slices after each completion.

**FR-7: Local commands without Claude.** `gl status`, `gl help`, `gl roadmap`, `gl changelog` read files directly and print output. No Claude process needed.

**FR-8: Interactive command launching.** `gl init` and `gl design` launch Claude in interactive mode (user present, not headless).

**FR-9: Frontmatter parsing.** The Go binary can parse flat key-value YAML frontmatter from `.greenlight/slices/*.md` files to determine slice status and dependencies.

**FR-10: State reading.** The Go binary can read GRAPH.json dependencies and compute which slices are ready (pending + all deps complete + not in_progress).

### 2.2 Non-Functional Requirements

**NFR-1: stdlib only.** All Go code uses stdlib only. No external dependencies.

**NFR-2: Performance.** `gl status` completes in under 1 second for 50+ slice files.

**NFR-3: Graceful shutdown.** SIGINT/SIGTERM handling for clean tmux session teardown.

**NFR-4: Error messages.** Clear actionable errors: tmux not installed, claude not in PATH, no slices ready, etc.

### 2.3 Constraints

- Go 1.24 stdlib only
- Must not break existing 1050 tests
- Extends existing binary (no separate CLI)
- tmux for parallelism (no custom terminal multiplexer)

### 2.4 Out of Scope

- Custom terminal multiplexer (tmux only)
- Remote execution (local machine only)
- Web dashboard for monitoring
- Windows support for tmux features

---

## 3. Technical Decisions

| # | Decision | Chosen | Rejected | Rationale |
|---|----------|--------|----------|-----------|
| D-40 | CLI language | Extend existing Go binary | Separate Node.js CLI | Single binary, already has CLI dispatch, flag parsing, and test patterns. No new runtime dependency. |
| D-41 | Parallelism mechanism | tmux sessions | Go goroutines with terminal mux; GNU screen | tmux is ubiquitous, provides visual monitoring, and users can attach/detach. Screen is less common. Goroutines would hide output. |
| D-42 | Claude process spawning | `os/exec.Command` with stdout/stderr capture | Shell scripts; embedded Claude SDK | Direct process control, Go stdlib, full lifecycle management. |
| D-43 | Frontmatter parsing | Simple key-value line parser (stdlib) | External YAML library; regex parsing | Flat key-value only (no nesting). Line-by-line string splitting is reliable and trivially testable. |
| D-44 | Context detection | `$CLAUDE_CODE` env var | Config flag; CLI flag | Environment variable is set automatically by Claude Code. No user action needed. Reliable detection. |
| D-45 | Watch interval | 30 seconds configurable | Real-time filesystem watching; 5 second polling | 30s balances responsiveness with CPU usage. File watching adds inotify/kqueue complexity. |
| D-46 | Command routing | Extend existing `cli.Run` switch statement | Separate command registry; plugin system | Consistent with existing pattern. Simple, tested, no abstraction needed for ~15 commands. |

---

## 4. Architecture

### 4.1 The Greenlight CLI

```
gl <command> [options]

Project lifecycle:
  gl init                 Interactive project setup -- interview, design, contracts, graph
  gl design               Run system design session -- produces DESIGN.md
  gl roadmap              Display/manage project roadmap

Building:
  gl slice [id]           Build slices -- auto-parallelises when multiple are ready
  gl quick [description]  Ad-hoc task with test-first guarantees

State & progress:
  gl status               Show project progress from slice state files
  gl debug [slice_id]     Diagnostic report for a slice

Pre-existing code:
  gl assess               Analyse codebase for gaps, risks, and wrap priorities
  gl map                  Analyse codebase with parallel mapper agents
  gl wrap [boundary]      Wrap existing boundary with contracts and locking tests

Release:
  gl ship                 Final verification -- all tests green, security audit

Session management:
  gl pause                Create handoff file when stopping work
  gl resume               Resume from previous session

Admin:
  gl migrate              Migrate STATE.md -> file-per-slice state
  gl settings             Configure model profiles, workflow options
  gl add-slice [name]     Add a new slice to the dependency graph
  gl changelog            Display project changelog
  gl help                 Show all commands and current state

Global options:
  --max N                 Maximum parallel sessions (default: 4)
  --watch                 Auto-refill -- start new slices as others complete
  --dry-run               Preview what would happen without doing it
  --sequential            Force sequential even if tmux is available
  --verbose               Show detailed output
```

### 4.2 New Internal Packages

```
internal/
  frontmatter/            # YAML frontmatter parser (stdlib only)
    frontmatter.go        # Parse/write flat key-value frontmatter
    frontmatter_test.go
  state/                  # State reader/writer
    state.go              # Read slices/*.md, GRAPH.json, compute ready slices
    state_test.go
  tmux/                   # tmux session management
    tmux.go               # Create/manage tmux sessions and panes
    tmux_test.go
  process/                # Claude process spawning
    process.go            # Spawn and manage claude processes
    process_test.go
```

### 4.3 Command Execution Flow

#### How interactive commands work

```bash
# gl init -- launches claude in interactive mode
$ gl init
# Launches interactive claude session with the init skill loaded
# No --dangerously-skip-permissions -- the user is present
```

For interactive commands, the CLI is a thin launcher:
1. Check `.greenlight/` state to provide context
2. Launch `claude` with the appropriate skill prompt
3. No `--dangerously-skip-permissions` -- the user is present

#### How autonomous commands work

```bash
# gl slice -- headless, can parallelise
$ gl slice
# 1. CLI reads state, finds ready slices
# 2. If 1: exec claude -p "/gl:slice {id}" --dangerously-skip-permissions
# 3. If 2+: spawn tmux, one pane per slice
```

#### How local commands work

```bash
# gl status -- no claude, instant
$ gl status
# CLI reads slices/*.md, GRAPH.json, computes and prints

# gl roadmap -- no claude, instant
$ gl roadmap
# CLI reads ROADMAP.md and prints

# gl changelog -- no claude, instant
$ gl changelog
# CLI reads summaries/*.md and prints
```

### 4.4 `gl slice` -- from the shell

```
gl slice [id] [--max N] [--watch] [--dry-run] [--sequential]
  |
  +-- Read .greenlight/slices/*.md frontmatter
  +-- Read GRAPH.json for deps
  +-- Find ready slices (status=pending, all deps=complete, not in_progress)
  |
  +-- id provided -> run single slice directly:
  |    exec claude -p "/gl:slice {id}" $CLAUDE_FLAGS
  |
  +-- 0 ready -> print what's blocked and why, exit
  |
  +-- 1 ready -> run directly, no tmux:
  |    exec claude -p "/gl:slice {id}" $CLAUDE_FLAGS
  |
  +-- 2+ ready:
       |
       +-- --sequential OR no tmux -> sequential mode
       |    +-- pick first ready slice by graph order
       |    +-- exec claude -p "/gl:slice {id}" $CLAUDE_FLAGS
       |    +-- on exit, re-read slice state, find next ready
       |    +-- repeat until no more ready or all blocked
       |
       +-- tmux available -> parallel mode
            +-- create tmux session "{prefix}-{project}" (e.g. "gl-greenlight")
            +-- one window per slice (up to --max)
            +-- each window runs:
            |    claude -p "/gl:slice {id}" $CLAUDE_FLAGS
            +-- tile layout
            +-- set status bar
            +-- attach to session
```

Where `$CLAUDE_FLAGS` comes from config:
```
--dangerously-skip-permissions --max-turns 200
```

### 4.5 `gl slice` -- from inside Claude

When `/gl:slice` calls `gl slice` and the CLI detects `$CLAUDE_CODE` is set:

```
gl slice [id]
  |
  +-- id provided -> output instructions for claude to build it
  |    (claude skill reads this and proceeds with the slice)
  |
  +-- 0 ready -> print blocked status
  |
  +-- 1 ready -> output the slice id for claude to build
  |
  +-- 2+ ready:
       +-- Output the first ready slice id for claude to build
       +-- Print to stderr (visible to user):
            "4 more slices ready. Run in a new terminal:
             gl slice --max 4"
```

The CLI never tries to spawn tmux from inside Claude. It handles one slice and surfaces the parallel opportunity to the user.

### 4.6 Watch Mode

When `gl slice --watch` is passed from the shell:

```
while true:
  sleep $WATCH_INTERVAL (default 30s)
  read all .greenlight/slices/*.md frontmatter

  for each in_progress slice:
    check if tmux pane is still alive
    if status changed to "complete":
      log "S-{id} complete ({test_count} tests)"

  ready = slices where status=pending AND all deps=complete
  running = count of in_progress slices
  slots = max - running

  if slots > 0 AND ready is not empty:
    for each ready slice (up to slots):
      spawn new tmux window: claude -p "/gl:slice {id}" $CLAUDE_FLAGS
      log "launched {id} ({name})"

  if no in_progress AND no ready:
    log "All slices complete or blocked"
    print summary (total done, total tests, remaining blocked)
    break
```

Run `gl slice --watch` once and walk away. It drains the entire dependency graph.

### 4.7 No-tmux Fallback

| Ready slices | With tmux | Without tmux |
|-------------|-----------|--------------|
| 0 | Print blocked status | Print blocked status |
| 1 | Run directly | Run directly |
| 2+ | Parallel tmux panes | Sequential: run one, re-scan, run next |
| 2+ with `--watch` | Parallel + auto-refill | Sequential + auto-refill |

Without tmux, the sequential fallback still auto-detects ready slices -- no manual picking. It just processes them one at a time.

### 4.8 tmux Session Layout

Named `{prefix}-{project}` from config (e.g. `gl-greenlight`), tiled automatically:

```
+-------------------------+-------------------------+
| S4.1                    | S4.2                    |
| SessionList component   | TrendChart component    |
|                         |                         |
| claude -p "/gl:slice    | claude -p "/gl:slice    |
|   S4.1" --dangerously.. |   S4.2" --dangerously.. |
+-------------------------+-------------------------+
| S4.3                    | S5.1                    |
| RemedyComparison        | Basic settings          |
|                         |                         |
| claude -p "/gl:slice    | claude -p "/gl:slice    |
|   S4.3" --dangerously.. |   S5.1" --dangerously.. |
+-------------------------+-------------------------+
```

### 4.9 Status Bar

tmux status line showing live progress from slice state files:

```bash
# Set by the CLI when creating the tmux session
tmux set -g status-right '#(gl status --compact)'
```

Where `gl status --compact` outputs: `18/36 done | 4 running`

### 4.10 Two Execution Contexts

The `gl` CLI behaves differently depending on where it's called:

| Context | Detection | Parallel strategy |
|---------|-----------|-------------------|
| **Shell** | `$CLAUDE_CODE` env var is NOT set | Spawns tmux + claude instances |
| **Inside Claude** | `$CLAUDE_CODE` env var IS set | Builds one slice, prints hint for the rest |

This means `/gl:slice` (the claude skill) can simply shell out to `gl slice` via the Bash tool, and the CLI detects it's inside Claude and does the right thing.

### 4.11 Relationship Between CLI and Claude Skills

```
+----------------------------------------------------------+
|  User's shell                                            |
|                                                          |
|  $ gl slice --watch                                      |
|    |                                                     |
|    +-- CLI reads state, finds ready slices               |
|    +-- CLI spawns tmux session                           |
|    +-- Each pane runs:                                   |
|    |    claude -p "/gl:slice S4.1" --dangerously-skip..  |
|    |      |                                              |
|    |      +-- Claude loads /gl:slice skill               |
|    |      +-- Skill calls: gl slice S4.1 (detects inside |
|    |      |   claude, returns slice info)                 |
|    |      +-- Skill builds the slice (tests -> impl)     |
|    |      +-- On completion, writes slices/S4.1.md       |
|    |                                                     |
|    +-- CLI polls slices/*.md every 30s                   |
|    +-- CLI detects S4.1 complete, spawns S4.4            |
|    +-- Repeat until done                                 |
+----------------------------------------------------------+
```

The CLI is the orchestrator. Claude is the worker. The slice state files are the communication channel between them.

---

## 5. Config Additions

```json
{
  "parallel": {
    "max": 4,
    "claude_flags": "--dangerously-skip-permissions --max-turns 200",
    "tmux_session_prefix": "gl",
    "watch_interval_seconds": 30
  }
}
```

Users can override `claude_flags` to use `--allowedTools` instead of `--dangerously-skip-permissions` if they want a safer setup.

---

## 6. Error Handling

| Scenario | Behaviour |
|----------|-----------|
| Slice fails (non-zero exit) | Mark `status: failed` in frontmatter, do NOT retry, log to pane |
| tmux not installed | Sequential fallback with install hint |
| claude not in PATH | Error and exit with install instructions |
| No slices ready | Print blocked slices and their unmet deps |
| Git conflict between panes | Committing session retries once, then marks failed |
| `--max-turns` exceeded | claude exits, slice stays `in_progress`, user investigates |
| All deps blocked | Print dependency chain, suggest which slice to unblock |
| Pane crashes mid-slice | Watch loop detects dead pane + `in_progress` status, logs warning |

---

## 7. Examples

```bash
# From shell -- auto-detects and parallelises
$ gl slice
-> 5 slices ready, launching 4 in tmux session "gl-greenlight"

# Watch mode -- fire and forget
$ gl slice --watch
-> 13 slices ready, launching 4... (9 queued)
S4.1 complete (42 tests) -> launched S4.4
S4.2 complete (38 tests) -> launched S5.2
...
All slices complete or blocked.
  28/36 done | 1,847 tests | 8 blocked on S6.1

# Specific slice
$ gl slice S2.8

# Preview
$ gl slice --dry-run
Ready (5):   S4.1, S4.2, S4.3, S5.1, S5.4
Running (2): S2.8, S3.4
Blocked (8): S3.5 (needs S3.4), S4.4 (needs S4.1,S4.2,S4.3), ...
Would launch: S4.1, S4.2, S4.3, S5.1

# Project status
$ gl status
Progress: [########..........] 18/36 slices
Running:  S2.8 (implementing), S3.4 (security)
Ready:    S3.6, S4.1, S4.2, S4.3, S5.1
Tests:    1,247 passing, 0 failing, 134 security
```

---

## 8. File Changes in Codebase

### 8.1 New Go Packages (4 packages)

| Package | Purpose | Est. Lines |
|---------|---------|-----------|
| `internal/frontmatter/` | Parse/write flat key-value YAML frontmatter from slice state files | ~150 |
| `internal/state/` | Read slices/*.md and GRAPH.json, compute ready/blocked/running slices | ~200 |
| `internal/tmux/` | Create/manage tmux sessions and panes via os/exec | ~200 |
| `internal/process/` | Spawn and manage Claude processes via os/exec | ~150 |

### 8.2 Modified Go Files

| File | Change |
|------|--------|
| `internal/cli/cli.go` | Add new subcommand cases to `Run()` switch: slice, status, help, roadmap, changelog, init, design, etc. |
| `internal/cli/cli.go:printUsage()` | Update help text with new commands |
| `internal/cmd/` | New command handler files: status.go, help.go, slice.go, roadmap.go, changelog.go, init.go, design.go |

### 8.3 No Change Required

- `main.go` -- args already forwarded to `cli.Run`, contentFS already available
- `internal/installer/` -- no manifest changes for CLI commands (they're Go code, not embedded content)
- `internal/version/` -- no changes

### 8.4 npm Wrapper Update

| File | Change |
|------|--------|
| `npm/bin/index.js` | Forward new subcommands to binary (should already work since it passes all args through) |

---

## 9. Proposed Build Order

### Phase 1: Foundation
1. **Frontmatter parser** (`internal/frontmatter`) -- Parse/write flat key-value frontmatter from `.greenlight/slices/*.md` files
2. **State reader** (`internal/state`) -- Read slice files, GRAPH.json deps, compute ready/blocked/running
3. **CLI dispatch refactoring** -- Add new subcommand cases to `cli.Run()`, update `printUsage()`

### Phase 2: Local Commands (no Claude needed)
4. **`gl status`** -- Read all slice files, compute summary, display progress
5. **`gl status --compact`** -- One-liner for tmux status bar
6. **`gl help`** -- List commands, detect project state, show context
7. **`gl roadmap`** -- Read and display ROADMAP.md
8. **`gl changelog`** -- Read and display from summaries/*.md

### Phase 3: Autonomous Commands (headless Claude)
9. **Process spawner** (`internal/process`) -- Spawn claude processes with configurable flags
10. **`gl slice {id}`** -- Single slice, headless Claude
11. **`gl slice` auto-detect** -- Find ready slices, run one
12. **tmux manager** (`internal/tmux`) -- Create/manage tmux sessions
13. **`gl slice` parallel** -- tmux spawning for 2+ ready slices
14. **`gl slice --watch`** -- Poll loop, auto-refill completed slots

### Phase 4: Interactive Commands
15. **`gl init`** -- Launch interactive Claude with init skill
16. **`gl design`** -- Launch interactive Claude with design skill

### Phase 5: Integration
17. **Signal handling** -- SIGINT/SIGTERM graceful shutdown for tmux sessions
18. **npm wrapper verification** -- Ensure new subcommands pass through correctly

---

## 10. Deferred

| Item | Why Deferred | When to Revisit |
|------|-------------|-----------------|
| `gl wrap` parallel | Parallelising boundary wrapping adds complexity | When users wrap 5+ boundaries regularly |
| `gl assess` / `gl map` CLI | These are less frequently used and work fine as Claude skills | When users run these from shell regularly |
| `gl ship` CLI | Final verification is a one-time operation per milestone | When the workflow is mature |
| Custom terminal mux | tmux is sufficient and ubiquitous | If tmux is unavailable on target platforms |
| Web dashboard | CLI monitoring is sufficient for 1-10 parallel sessions | When teams need shared visibility |
| Windows tmux support | macOS and Linux are primary targets | If Windows becomes a supported platform |

---

## 11. User Decisions (Locked)

| # | Gray Area | Decision | Rationale |
|---|-----------|----------|-----------|
| 1 | CLI language | Extend existing Go binary | Single binary, already has CLI dispatch, flag parsing, test patterns. No new runtime dependency. |
| 2 | Workflow approach | Full Greenlight workflow (/gl:design -> /gl:init -> /gl:slice) | Follows established pattern, generates contracts and dependency graph properly. |
| 3 | Parallelism mechanism | tmux | Ubiquitous, visual monitoring, attach/detach support. |
| 4 | Context detection | $CLAUDE_CODE env var | Automatic, reliable, no user action needed. |
| 5 | Watch interval | 30 seconds (configurable) | Balances responsiveness with CPU usage. |
| 6 | Sequential fallback | Auto-detect and run one-at-a-time | No manual picking, graceful degradation. |
| 7 | Claude flags | Configurable in config.json parallel section | Users can switch between --dangerously-skip-permissions and --allowedTools. |
