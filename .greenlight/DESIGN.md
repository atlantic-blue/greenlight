# System Design: Circuit Breaker Module

> **Project:** Greenlight
> **Scope:** Add circuit breaker module to prevent implementer death spirals -- attempt tracking, structured diagnostics, scope lock, manual override (`/gl-debug`), and rollback via git tags.
> **Stack:** Go 1.24, stdlib only. Deliverables are embedded markdown content plus two manifest entries.
> **Date:** 2026-02-18
> **Replaces:** Previous DESIGN.md (brownfield-and-docs -- complete, 415 tests passing)

---

## 1. Problem Statement

When the Greenlight implementer gets stuck on a failing test, it loops endlessly -- guessing at root causes, compounding bad changes, and eventually inventing phantom problems that don't exist. Real users lose full days to this. There is no detection mechanism, no structured diagnostic output, no rollback capability, and no manual pull cord.

**Real-world example:** A user asked the implementer to redesign campaign visualization from bars to card headers. After multiple attempts, the agent started claiming "the app was defaulting to stream view instead of calendar view" -- a fabricated root cause unrelated to the slice. Fresh slices, new contracts, and mockups all failed to break the cycle because the agent has no self-awareness that it's stuck.

### Current State Gap

The implementer's `<error_recovery>` section (gl-implementer.md lines 146-181) has a rudimentary "know when to stop" instruction: after 3 targeted attempts, document and report. This is:

1. **Aspirational, not enforced.** No state tracking structure -- the implementer is told to count, but there's no mechanism ensuring it does.
2. **No diagnostic format.** "Document what you've tried" produces inconsistent, unhelpful output.
3. **No rollback.** After 3 bad attempts, the codebase is full of debris with no clean state to return to.
4. **No scope awareness.** Nothing prevents the agent from modifying unrelated files when it starts grasping at straws.
5. **No manual override.** The user has to wait for the agent to finish spiraling -- no pull cord.

---

## 2. Requirements

### 2.1 Functional Requirements

**FR-1: Attempt Tracking.** The implementer must track every attempt per test per slice. State includes: slice_id, test_name, attempt_count, files touched per attempt, description of each attempt, and last known green reference (checkpoint tag).

**FR-2: Auto-Trip at 3 Per-Test Failures.** After 3 failed attempts on any single test, the circuit trips automatically. The implementer stops coding and switches to diagnostic mode. This is mandatory, not optional.

**FR-3: Slice-Level Ceiling at 7 Total Failures.** If total failed attempts across ALL tests in a slice exceed 7, the circuit trips regardless of per-test counts. This catches the pattern where the agent redistributes failures across different tests without making real progress.

**FR-4: Structured Diagnostic Report.** When the circuit trips, the implementer produces a structured diagnostic report in a prescribed format. The report includes: what the test expects, what actually happens (exact error output), an attempt log with files touched, cumulative file modifications, scope violations, best hypothesis, specific question for the user, and recovery options.

**FR-5: Scope Lock with Justify-or-Stop.** Before modifying any file, the implementer checks if it's within the slice's inferred scope (derived from contracts). Out-of-scope modifications require explicit, one-sentence justification tied to the current failing test. Unjustifiable out-of-scope touches count as scope violations and failed attempts.

**FR-6: Manual Override (/gl-debug).** A standalone command the user can run at any time to force the diagnostic report immediately, without waiting for 3 attempts. Reads current state, produces the diagnostic, and presents it. Does not control the pipeline.

**FR-7: Rollback via Lightweight Git Tags.** Before the implementer's first attempt on any test, a lightweight git tag is created as a checkpoint. On circuit break, the diagnostic report includes the exact rollback command. After human input and counter reset, the implementer starts from the tagged clean state.

**FR-8: Counter Reset on Human Input.** When the user provides guidance after a circuit break, the per-test counter for the affected test resets to 0. The slice-level accumulator also resets. The implementer rolls back to the checkpoint and tries again with the user's input as additional context.

### 2.2 Non-Functional Requirements

**NFR-1: Zero Go Code for Protocol.** The circuit breaker protocol is entirely embedded content (markdown files in `src/`). The only Go change is adding the new command and reference to the installer manifest.

**NFR-2: Context Budget.** The circuit breaker reference document must be concise enough that loading it does not push the implementer agent past the 50% context threshold. Target: under 200 lines.

**NFR-3: Structured Data First.** The diagnostic report is conceptually structured data rendered as markdown, not markdown that happens to contain data. Field names are consistent and machine-parseable, enabling future integration with pause/resume and automated tooling.

### 2.3 Constraints

- Go 1.24, stdlib only. No external dependencies.
- All content embedded via `go:embed` from `src/` directory.
- Must integrate with existing agent isolation rules (implementer cannot see test source code).
- Must not break existing 415 tests.
- Must preserve existing deviation rules protocol -- circuit breaker is additive, not a replacement.

### 2.4 Out of Scope

- **Automated debugging.** The circuit breaker stops and reports; it does not attempt to fix. The existing debugger agent handles investigation when the user routes to it.
- **Cross-slice learning.** The circuit breaker does not learn from previous slice failures. Each slice starts fresh.
- **UI/dashboard.** No visual interface for circuit breaker state. Output is markdown in the terminal.
- **Automatic re-routing to debugger.** The user decides whether to route to the debugger, re-scope, or retry. Future work could automate this.
- **Test modification.** The circuit breaker never suggests or performs test changes. Tests are the source of truth.
- **Configurable thresholds.** 3 per-test and 7 per-slice are fixed. Configurability deferred until real-world data calibrates the right defaults.

---

## 3. Technical Decisions

| # | Decision | Chosen | Rejected | Rationale |
|---|----------|--------|----------|-----------|
| 1 | Scope lock source | Inferred from contracts (default) with optional `files_in_scope` override in GRAPH.json | Manual allowlist maintained by architect; Pure inference with no override | Contracts define what you're building -- scope is inferrable (component, types, service, styles, parent for wiring). Optional override handles edge cases (shared utility, config file) without architect busywork. |
| 2 | Attempt counter granularity | Per-test tracking with slice-level ceiling (7 total failures) | Per-test only; Per-slice only | Per-test catches single-test spirals. Slice ceiling catches failure redistribution where agent appears to make progress but just breaks different things each time. |
| 3 | Rollback mechanism | Lightweight git tags (`greenlight/checkpoint/{slice_id}`) | Git commits with checkpoint prefix; Git stash | Tags give a named, findable reference to a known-good state without polluting commit history. Easy to create, find, roll back to, and clean up when done. Stashes are a stack, unnamed by default, and get lost. |
| 4 | /gl-debug integration level | Standalone diagnostic command, architecturally designed for future pause/resume integration | Integrated with /gl:pause pipeline; Full pipeline controller | Simplicity now. /gl-debug reads state, produces report, user decides. Structured report format enables future automation without redesign. |
| 5 | CLAUDE.md integration pattern | 5-line hard rule in CLAUDE.md + full protocol in `references/circuit-breaker.md` | Full protocol in CLAUDE.md (~50 lines); Brief reference only in CLAUDE.md | Follows existing deviation-rules.md pattern. Keeps CLAUDE.md focused on universal standards. Full protocol with examples lives in reference doc. Hard rule in CLAUDE.md ensures compliance -- agents cannot claim ignorance. |
| 6 | Diagnostic report structure | Structured fields with consistent naming, rendered as markdown | Free-form markdown; JSON output | Structured data rendered as markdown is human-readable now and machine-parseable later. Field names like `slice`, `test`, `trip_reason`, `recovery_options` can be extracted programmatically for future tooling. |
| 7 | Slice-level ceiling threshold | 7 total failures across all tests | 5 (too aggressive -- trips on 2 tests); 10 (too permissive -- 3+ tests fully exhausted) | With 3 attempts per test, 7 means the agent has failed on at least 3 different tests (or exhausted one and partially exhausted others). Enough signal that the problem is systemic, not isolated. |
| 8 | Tag cleanup timing | Tags cleaned up at slice completion (Step 10 of /gl:slice) | Immediately after test passes; Never cleaned up | Keeping tags until slice completion means rollback is always available if a later test's fix breaks an earlier test. Cleaning at completion keeps the tag namespace tidy. |

---

## 4. Architecture

### 4.1 Component Overview

The circuit breaker is five components distributed across existing system files plus two new files:

```
src/
  CLAUDE.md                          # +5 lines: hard rule reference
  references/
    circuit-breaker.md               # NEW: full protocol (~180 lines)
                                     #   - attempt tracking rules
                                     #   - diagnostic report format
                                     #   - scope lock rules
                                     #   - rollback protocol
                                     #   - good vs bad examples
  agents/
    gl-implementer.md                # REWRITE: <error_recovery> replaced
                                     #   with circuit breaker integration
  commands/
    gl/
      debug.md                       # NEW: /gl-debug command definition
      slice.md                       # MODIFY: checkpoint tagging before
                                     #   implementation, handle diagnostic
                                     #   output, tag cleanup at completion

internal/
  installer/
    installer.go                     # +2 manifest entries:
                                     #   "commands/gl/debug.md"
                                     #   "references/circuit-breaker.md"
```

### 4.2 Component 1: Attempt Tracker

**Lives in:** `references/circuit-breaker.md` (protocol) + `gl-implementer.md` (integration)

The implementer maintains attempt state as a mental model. The state structure:

```
Circuit Breaker State:
  slice_id: {from orchestrator input}
  checkpoint_tag: greenlight/checkpoint/{slice_id}
  per_test_attempts:
    {test_name}:
      count: N
      changes:
        - attempt: 1
          files_touched: [path1, path2]
          description: "Added UserService with email validation"
          error: "expected 409, got 500: missing unique constraint check"
        - attempt: 2
          ...
  total_slice_failures: N
  scope_violations: [{file, justification, test_name, verdict}]
```

On each test failure:
1. Increment `per_test_attempts[test_name].count`
2. Record `{attempt_number, files_touched, description_of_change, exact_error}`
3. Increment `total_slice_failures`
4. Check trip conditions: `per_test count >= 3` OR `total_slice_failures >= 7`
5. If tripped: **stop coding immediately**, produce diagnostic report

The orchestrator (`/gl:slice`) passes attempt history when spawning a fresh implementer after counter reset, so the fresh agent knows what was already tried and does not repeat the same approaches.

### 4.3 Component 2: Diagnostic Report Generator

**Lives in:** `references/circuit-breaker.md`

When the circuit trips (auto or manual), the implementer produces this exact format:

```markdown
## Circuit Breaker Tripped

**Slice:** {slice_id} -- {slice_name}
**Test:** {test_name}
**Trip reason:** {per-test limit (3/3) | slice ceiling (N/7)}

### What the test expects
{Exact assertion from test name and contract -- not interpretation}

### What actually happens
{Exact error output -- copy/paste, no summarizing}

### Attempts log
1. **{description}** -> files: [{list}] -> result: {exact error}
2. **{description}** -> files: [{list}] -> result: {exact error}
3. **{description}** -> files: [{list}] -> result: {exact error}

### Files modified (cumulative)
| File | Change Type | Attempts |
|------|-------------|----------|
| {path} | created | 1, 2 |
| {path} | modified | 1, 2, 3 |

### Scope violations
| File | Justification Given | Test | Verdict |
|------|---------------------|------|---------|
| {path} | {reason} | {test} | justified / violation |

{Or: "None"}

### Best hypothesis
{One clear hypothesis based on the pattern across attempts. Derived from
what changed between attempts and how the error shifted. Not a guess.}

### What I need from you
{Specific question or decision. Must be answerable. Examples:
  - "Is the contract correct? The test expects X but the component API provides Y."
  - "Should this component own its own state or receive it as props?"
  - "The error suggests a dependency on {file} -- should I modify it?"
NOT: "help me" or "I'm stuck" or "I don't know what's wrong"}

### Recovery options
- [ ] Roll back to checkpoint: `git checkout greenlight/checkpoint/{slice_id}`
- [ ] Provide guidance and reset counter (I'll retry with your input)
- [ ] Re-scope this slice (the contract may need adjustment)
- [ ] Skip this test and continue with remaining tests
```

### 4.4 Component 3: Scope Lock

**Lives in:** `references/circuit-breaker.md` (rules) + `gl-implementer.md` (enforcement)

Before modifying any file, the implementer performs scope inference:

**In-scope (no justification needed):**
- Files directly implied by the contract (the component, its types, its service, its styles)
- Files in the same directory/module as contract entities
- Files explicitly listed in `files_in_scope` in GRAPH.json (if present)

**Requires justification (one sentence, must reference current failing test):**
- Parent components (to pass props or wire up)
- Shared utilities referenced by contract entities
- Configuration files needed for new functionality

**Automatic scope violation -- do NOT proceed:**
- Routing changes for a UI component slice
- Database schema changes not in the contract
- Unrelated view/page modifications
- Global state changes not implied by the contract
- Any file where the justification references concerns outside the current slice

The heuristic: **if you cannot explain, in one sentence tied to the current failing test, why modifying this file makes that test pass, do not touch it.**

When a scope violation is detected:
1. Log it in the attempt record
2. Count it as a failed attempt
3. Do NOT make the change
4. If this triggers the circuit breaker, include it in the diagnostic report

### 4.5 Component 4: /gl-debug Command

**Lives in:** `commands/gl/debug.md` (new file)

Standalone diagnostic command. No pipeline integration.

**Trigger:** User runs `/gl-debug` at any time during implementation.

**Behavior:**
1. Read `.greenlight/STATE.md` to determine current slice and step
2. Read current git status and recent commits
3. Run the test suite to capture current failures
4. Gather whatever attempt data is available from the current conversation context
5. Produce the diagnostic report format (same as Component 2) with available data
6. Present to user
7. Done -- user decides next action

**What it does NOT do:**
- Modify any files
- Reset any counters
- Interact with the /gl:slice pipeline
- Spawn other agents
- Create or delete git tags

**Design for future:** The report uses the same structured format as the auto-trip report. A future version could write this to `.continue-here.md` for `/gl:resume` integration.

### 4.6 Component 5: Rollback Integration

**Lives in:** `commands/gl/slice.md` (modifications to Steps 3 and 10)

**Before implementation (new sub-step in Step 3 of /gl:slice):**

```bash
# Create checkpoint tag before implementer starts
git tag greenlight/checkpoint/{slice_id}
```

If the tag already exists (resumed slice or retry), verify it points to a clean state and skip re-tagging.

**On circuit break:**

The diagnostic report includes the rollback command:
```
Roll back to checkpoint: git checkout greenlight/checkpoint/{slice_id}
```

**On counter reset (user provides input after circuit break):**

The orchestrator:
1. Rolls back to the checkpoint: `git checkout greenlight/checkpoint/{slice_id}`
2. Spawns a fresh implementer with:
   - The user's guidance as additional context
   - The diagnostic report as "what was already tried" (so the agent doesn't repeat approaches)
   - Reset counters (per-test = 0, slice total = 0)
3. Re-tags if the rollback moved HEAD

**On slice completion (added to Step 10 of /gl:slice):**

```bash
# Clean up checkpoint tags for this slice
git tag -d greenlight/checkpoint/{slice_id} 2>/dev/null
```

### 4.7 CLAUDE.md Addition

Added to `src/CLAUDE.md` after the "Deviation Rules" section:

```markdown
## Circuit Breaker Protocol

**MANDATORY.** The implementer MUST follow `references/circuit-breaker.md`. No exceptions.
Count every attempt. Stop at 3 per test or 7 per slice. Produce the structured diagnostic.
Justify every out-of-scope file touch. Roll back to checkpoint on reset.
This protocol is not optional and cannot be overridden by any agent.
```

---

## 5. Data Model

No database entities. No persistent state files. The circuit breaker operates within:

| Storage | What | Lifetime |
|---------|------|----------|
| Agent mental state | Attempt counters, scope tracking, file touch log | Single agent spawn (not persisted between spawns) |
| Git tags | `greenlight/checkpoint/{slice_id}` | Created before implementation, deleted at slice completion |
| Orchestrator context | Attempt history passed to fresh implementer spawns | Single /gl:slice execution |
| Diagnostic report | Markdown output | Displayed to user (ephemeral, optionally saved by user) |

---

## 6. Integration with Existing Systems

### 6.1 Agent Isolation (CLAUDE.md)

No changes to agent isolation rules. The implementer still cannot see test source code. The circuit breaker works with test names and error output only.

The agent isolation table gains no new rows. `gl-debugger` already has full read access -- when the user routes a circuit-break diagnostic to the debugger, the debugger can investigate without restriction.

### 6.2 Deviation Rules (references/deviation-rules.md)

Circuit breaker is additive to deviation rules, not a replacement:

- **Rules 1-3 (auto-fix)** still apply. If the implementer encounters a bug, missing functionality, or blocker, it fixes immediately per deviation rules.
- **Rule 4 (architectural stop)** still takes priority. An architectural issue stops execution before the circuit breaker can trip.
- **Circuit breaker** catches a different failure mode: the agent is stuck on a test, not encountering a deviation. It's trying to implement correctly but failing. Deviation rules handle unplanned work; circuit breaker handles unproductive work.

### 6.3 Checkpoint Protocol (references/checkpoint-protocol.md)

The circuit-break diagnostic is a new checkpoint type. It fits alongside the existing types:

| Checkpoint Type | Trigger | When to Pause |
|-----------------|---------|---------------|
| Visual | UI slice needs human eyes | interactive mode only |
| Decision | Rule 4 architectural stop | always |
| External Action | Human action needed outside Claude | always |
| **Circuit Break** | **3 per-test or 7 per-slice failures** | **always** |

Circuit break checkpoints always pause, even in yolo mode. A stuck agent cannot self-recover.

### 6.4 Existing /gl:slice Pipeline

Changes to the pipeline:

| Step | Current | After Circuit Breaker |
|------|---------|----------------------|
| Step 3 (Implement) | Spawn implementer, handle pass/fail | Add: create checkpoint tag before spawn |
| Step 3 failure handling | "Max 3 implementation attempts, then pause" | Replace: circuit breaker diagnostic report, user chooses recovery option |
| Step 10 (Complete) | Update state, final report, suggest next | Add: clean up checkpoint tags |

The existing "max 3 implementation attempts" at the orchestrator level is superseded by the circuit breaker's per-test and slice-level tracking. The orchestrator no longer needs its own retry counter -- the implementer self-reports when it's stuck.

---

## 7. File Changes Summary

| File | Change | Lines Added/Modified |
|------|--------|---------------------|
| `src/references/circuit-breaker.md` | **NEW** | ~180 lines |
| `src/commands/gl/debug.md` | **NEW** | ~80 lines |
| `src/agents/gl-implementer.md` | Rewrite `<error_recovery>` section | ~40 lines replaced |
| `src/commands/gl/slice.md` | Add checkpoint tag creation (Step 3) and cleanup (Step 10) | ~25 lines added |
| `src/CLAUDE.md` | Add Circuit Breaker Protocol section | 5 lines added |
| `internal/installer/installer.go` | Add 2 manifest entries | 2 lines added |
| `internal/installer/manifest_docs_test.go` | Update manifest count in tests | ~2 lines modified |

**Total new content:** ~260 lines of markdown, 4 lines of Go.

---

## 8. Proposed Build Order

These are logical groupings for the architect to refine into slices with contracts:

1. **Attempt Tracker + Diagnostic Report** -- Core protocol in `references/circuit-breaker.md` and `gl-implementer.md` rewrite. This is the foundation everything else depends on.

2. **Scope Lock** -- Add scope inference and justification rules to the circuit breaker reference and implementer agent. Depends on the attempt tracker (scope violations count as failed attempts).

3. **Rollback Integration** -- Modify `/gl:slice` to create checkpoint tags before implementation and clean up at completion. Depends on the attempt tracker (rollback happens on circuit break).

4. **/gl-debug Command** -- New standalone command. Depends on the diagnostic report format being defined (shares the same format).

5. **CLAUDE.md + Manifest Integration** -- Add the hard rule to CLAUDE.md, add manifest entries to installer.go, update tests. Final slice that ties everything together.

6. **End-to-End Verification** -- Verify the full circuit breaker flow works: implementer tracks attempts, trips at 3, produces diagnostic, user provides input, rollback to tag, fresh implementer with guidance, test passes.

---

## 9. Security

No new security surface. The circuit breaker:
- Does not expose any external interfaces
- Does not handle user input beyond the existing slash command pattern
- Does not store credentials or sensitive data
- Git tags are local only (not pushed to remote unless user explicitly pushes tags)

The scope lock component has a **positive security effect**: it prevents the implementer from making unauthorized changes to files outside the slice boundary, reducing the blast radius of implementation errors.

---

## 10. Deferred

| Item | Why Deferred | When to Revisit |
|------|-------------|-----------------|
| Automated debugger routing | Adds pipeline complexity; user decision after circuit break is valuable signal | After circuit breaker has real-world usage data showing users always route to debugger |
| /gl-debug + /gl:pause integration | Need structured report format stable first | Next milestone after circuit-breaker |
| Cross-slice failure learning | Requires persistent state between slices; unclear value | When users report recurring failure patterns across slices |
| Configurable thresholds (3/7) | May not be right for all projects; need calibration data | After 20+ real-world circuit breaks provide calibration data |
| Diagnostic report persistence | Currently ephemeral; could save to `.greenlight/diagnostics/` | When users request history of circuit breaks for post-mortem analysis |

---

## 11. User Decisions (Locked)

| # | Gray Area | Decision | Rationale |
|---|-----------|----------|-----------|
| 1 | Scope lock source | Inferred from contracts (default) with optional `files_in_scope` override in GRAPH.json | Contracts define what you're building; scope is inferrable. Optional override handles edge cases without architect busywork. |
| 2 | Attempt counter granularity | Per-test with slice-level ceiling (7) | Per-test catches single-test spirals. Slice ceiling catches failure redistribution. Both failure modes covered. |
| 3 | Rollback mechanism | Lightweight git tags | Named reference to known-good state. Zero commit history pollution. Easy to create, find, roll back to, and clean up. |
| 4 | /gl-debug integration | Standalone now, structured for future pause/resume integration | Simple pull cord. Structured report format enables future automation without redesign. |
| 5 | CLAUDE.md integration | 5-line hard rule in CLAUDE.md + full protocol in references/circuit-breaker.md | Follows existing deviation-rules.md pattern. CLAUDE.md stays lean. Hard rule ensures compliance. |
