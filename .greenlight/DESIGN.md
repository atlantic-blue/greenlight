# System Design: Verification Tiers

> **Project:** Greenlight
> **Scope:** Add verification tier system to close the gap between "tests pass" and "user got what they asked for" -- per-contract verification levels (auto/verify), human acceptance gates, rejection-to-TDD routing, and escalation.
> **Stack:** Go 1.24, stdlib only. Deliverables are embedded markdown content plus one manifest entry.
> **Date:** 2026-02-19
> **Replaces:** Previous DESIGN.md (circuit-breaker -- complete, 456 tests passing)

---

## 1. Problem Statement

Tests pass but output doesn't match user intent. The current system has a blind spot: after the verifier confirms contract coverage and test results, the slice is marked complete -- but "tests pass" doesn't mean "the user got what they asked for." The system is marking its own homework.

**Real-world failure mode:** User asks for a card-based campaign visualization. Tests pass (component renders, data flows, assertions pass). But the actual output is a table with card styling -- not the card layout the user intended. The verifier sees tests passing and marks the slice complete. The user only discovers the mismatch later.

### Current State Gap

The existing `/gl:slice` pipeline (Step 6 verification + Step 9 visual checkpoint) has four gaps:

1. **No human verification for non-UI slices.** The `visual_checkpoint` toggle only triggers for UI slices. Business logic, API behavior, and data processing slices skip human verification entirely. Tests pass and the slice is marked complete.

2. **Visual checkpoint is binary.** It either runs or it doesn't. There is no middle ground between "automated verification only" and "user must run the app and look at it." A structured acceptance review (checking criteria without running the app) is not supported.

3. **No rejection-to-TDD loop.** When the visual checkpoint surfaces a mismatch, the orchestrator spawns a debugger. This breaks TDD discipline -- the correct response is to write a test that captures the user's intent, then make the implementer pass it.

4. **No escalation on repeated rejection.** If the user keeps rejecting, the system keeps looping. There is no circuit-breaker equivalent for human verification -- no threshold where the system acknowledges "this slice needs re-scoping."

---

## 2. Requirements

### 2.1 Functional Requirements

**FR-1: Per-Contract Verification Tier.** Every contract has a `verification` field with value `auto` or `verify`. Default is `verify`. The tier determines what happens after tests pass and the verifier approves.

**FR-2: Auto Tier Behavior.** When all contracts in a slice are tier `auto`, the slice proceeds directly from verification (Step 6) to summary/docs (Step 7). No human checkpoint. This is the current default behavior, preserved for infrastructure, config, and internal plumbing slices.

**FR-3: Verify Tier Behavior.** When any contract in a slice has tier `verify`, the orchestrator presents a structured checkpoint combining acceptance criteria and optional steps. The user confirms the output matches intent or describes what doesn't match. This subsumes the existing `visual_checkpoint` functionality.

**FR-4: Tier Resolution.** A slice's effective verification tier is the highest tier among all its contracts: `verify` > `auto`. Acceptance criteria and steps from all `verify` contracts are aggregated into a single checkpoint. One approval per slice.

**FR-5: Rejection Flow to Test Writer.** When the user rejects (anything other than "approved"), the orchestrator collects the feedback, presents structured options that implicitly classify the gap (test gap, implementation gap, or contract gap), and routes accordingly:
- Test gap or implementation gap: spawn test writer with rejection feedback (behavioral, no implementation details), contract, and acceptance criteria. New tests written, then implementer makes them pass.
- Contract gap: present to user for contract revision, then restart slice.

**FR-6: Rejection Counter with Escalation.** Track rejections per slice. After 3 rejections on the same slice, trigger escalation with options: re-scope the slice, pair with user for detailed guidance, or skip verification (mark as auto and proceed with known mismatch).

**FR-7: Contract Schema Extension.** Add three optional fields to the contract format: `verification` (auto/verify, default verify), `acceptance_criteria` (list of behavioral criteria), `steps` (list of steps to verify, optional). If `verification` is `verify`, warn if both `acceptance_criteria` and `steps` are empty. If `verification` is `auto`, `acceptance_criteria` and `steps` are ignored.

**FR-8: Backward Compatibility.** Contracts without an explicit `verification` field default to `verify`. The existing `config.workflow.visual_checkpoint` toggle is deprecated with a warning -- tiers in contracts supersede it.

### 2.2 Non-Functional Requirements

**NFR-1: Zero Go Code for Protocol.** The verification tier protocol is entirely embedded content (markdown files in `src/`). The only Go change is adding the new reference file to the installer manifest.

**NFR-2: Context Budget.** The verification tiers reference document must be concise enough that loading it does not push the orchestrator agent past the 30% context threshold. Target: under 150 lines.

**NFR-3: Agent Isolation Preserved.** The test writer receives the user's rejection feedback verbatim (behavioral description), the contract, and acceptance criteria. No implementation code, no test source code from the current cycle. The implementer still receives test names only.

**NFR-4: Checkpoint Consistency.** Verification tier checkpoints follow the same format patterns as existing checkpoint types (Visual, Decision, External Action, Circuit Break) in `references/checkpoint-protocol.md`.

### 2.3 Constraints

- Go 1.24, stdlib only. No external dependencies.
- All content embedded via `go:embed` from `src/` directory.
- Must integrate with existing agent isolation rules.
- Must not break existing 456 tests.
- Must preserve existing circuit breaker protocol -- verification tiers are additive.
- Must preserve existing deviation rules protocol.
- Rejection routing must go through the test writer first (TDD-correct approach).

### 2.4 Out of Scope

- **Automated acceptance detection.** The system does not auto-detect whether acceptance criteria are met. A human decides.
- **Screenshot comparison.** No visual diffing or pixel-level verification. Human eyes only.
- **Partial approval.** The user approves or rejects the entire slice checkpoint, not individual criteria. Future work could support per-criterion approval.
- **Rejection history persistence.** Rejection feedback is passed to the test writer in the current cycle. It is not persisted to `.greenlight/` for cross-session reference. Future work could save rejection history.
- **Configurable rejection threshold.** 3 rejections per slice is fixed. Matches the circuit breaker's per-test threshold. Configurability deferred until real-world data calibrates the right default.
- **AI-assisted gap classification.** The orchestrator presents options to the user. It does not use AI to auto-classify the gap type. The user's choice is the classification.

---

## 3. Technical Decisions

| # | Decision | Chosen | Rejected | Rationale |
|---|----------|--------|----------|-----------|
| 1 | Verification tier count | Two tiers: `auto` and `verify` | Three tiers (auto/review/demo) -- the distinction between review and demo is artificial. In both cases the user opens the app, looks at output, and confirms intent. The contract author doesn't need to choose between "check a list" and "walk through steps." | Simpler is better. One decision: "Can tests alone capture my intent for this slice?" Two tiers, not three. |
| 2 | Default verification tier | `verify` -- forgetting to set a tier gives a human checkpoint | `auto` (current behavior, but unsafe) | Safe default. Forgetting to annotate a contract gives you human verification, not silent auto-approve. The cost of an unnecessary verify checkpoint is low (user types "approved"). The cost of a missing checkpoint is a completed slice that doesn't match intent. |
| 3 | Tier location | In the contract, not in GRAPH.json or config.json | GRAPH.json (per-slice only, loses contract granularity); config.json (global only, no per-boundary control) | The contract is the source of truth for what a boundary does and how it's verified. Tier is a property of the boundary, not the build graph or the project config. |
| 4 | Rejection routing | Always through test writer first, even for implementation gaps | Direct to implementer with "fix this" instructions; Direct to debugger | TDD-correct. If the implementation is wrong and tests pass, the tests weren't tight enough. Adding a test that specifically asserts the user's intent, then making the implementer pass it, is the right fix. The implementer never gets "fix this" -- they get new failing tests. |
| 5 | Tier resolution across contracts | Highest tier wins + aggregation. verify > auto. One checkpoint per slice. | Per-contract checkpoints (too many interruptions); Per-slice config only (loses contract granularity) | Minimizes user interruptions. A slice with mixed tiers gets the highest tier with all criteria aggregated. |
| 6 | Rejection counter scope | Per-slice, escalation at 3 | Per-contract (each contract gets 3 rejections independently) | A slice is the unit of work. If a user has rejected 3 times, the whole slice needs re-scoping. Mirrors the circuit breaker's slice-level ceiling concept. |
| 7 | Gap classification UX | Orchestrator presents actionable options that implicitly map to gap types | Ask user to classify directly ("Is this a test gap?"); Auto-classify with heuristics | Users should not need to understand Greenlight's internal taxonomy. Present options like "tighten the tests," "revise the contract," or "provide more detail." The user picks what feels right; the orchestrator routes. |
| 8 | Rejection feedback to test writer | User's verbatim behavioral feedback + contract + acceptance criteria. No implementation details. | Include implementation code (breaks isolation); Include test source (breaks TDD) | Preserves agent isolation. The user's rejection is behavioral ("I expected X, I got Y"). The test writer uses this behavioral description plus the contract to write tighter tests. No implementation leakage. |
| 9 | visual_checkpoint backward compatibility | Deprecate with warning. Tiers in contracts supersede it. Keep in config but ignore. | Remove from config (breaking change); Add new config field (unnecessary -- tiers live in contracts) | The default `verify` tier already provides human verification for all non-auto slices. The `visual_checkpoint` toggle is redundant. Deprecation with warning is the cleanest migration path. |
| 10 | New reference file | New `references/verification-tiers.md` (follows circuit-breaker.md pattern) | Extend `references/verification-patterns.md` (different concerns -- automated vs human) | Existing verification-patterns.md covers automated verification (stubs, wiring, test quality). Verification tiers cover human acceptance. Different concerns, different audiences. |
| 11 | Verify checkpoint content | Combined: `acceptance_criteria` (what to check) + `steps` (how to check, optional) under one tier | Separate review and demo tiers with different required fields | The user does the same thing in both cases: look at the output and confirm intent. Criteria describe what to verify. Steps describe how, when the how isn't obvious. Both optional but warn if neither present. |

---

## 4. Architecture

### 4.1 Component Overview

The verification tier system is five components distributed across existing system files plus one new file:

```
src/
  CLAUDE.md                          # +4 lines: hard rule reference
  references/
    verification-tiers.md            # NEW: full protocol (~130 lines)
                                     #   - tier definitions and defaults
                                     #   - verify checkpoint format
                                     #   - rejection flow
                                     #   - rejection counter + escalation
                                     #   - gap classification routing
    checkpoint-protocol.md           # MODIFY: add Acceptance checkpoint type
                                     #   deprecate Visual, update mode table
    verification-patterns.md         # MODIFY: add cross-reference to tiers
  agents/
    gl-architect.md                  # MODIFY: add verification/acceptance_criteria/
                                     #   steps to contract format
    gl-verifier.md                   # MODIFY: report tier in verification output
  commands/
    gl/
      slice.md                       # MODIFY: add Step 6b (verification tier gate)
                                     #   modify Step 9 (deprecate, reference tiers)
                                     #   add rejection flow handling
  templates/
    config.md                        # MODIFY: deprecation note on visual_checkpoint

internal/
  installer/
    installer.go                     # +1 manifest entry:
                                     #   "references/verification-tiers.md"
```

### 4.2 Component 1: Verification Tier Gate (Step 6b)

**Lives in:** `commands/gl/slice.md` (new step) + `references/verification-tiers.md` (protocol)

After Step 6 (verification passes) and Step 6a (locking-to-integration transition, if applicable), the orchestrator reads the verification tier for each contract in the slice and resolves the effective tier.

**Tier resolution:**

```
For each contract in the slice:
  Read contract.verification (default: "verify")

Effective tier = highest tier among all contracts:
  verify > auto

If effective tier is auto:
  Skip to Step 7 (summary/docs)

If effective tier is verify:
  Aggregate acceptance_criteria from all verify contracts
  Aggregate steps from all verify contracts
  Present Verify Checkpoint
  Wait for user response
```

### 4.3 Component 2: Verify Checkpoint

**Lives in:** `references/verification-tiers.md`

Format presented to the user:

```
ALL TESTS PASSING — Slice {slice_id}: {slice_name}

Please verify the output matches your intent.

Acceptance criteria:
  [ ] {criterion 1 from contract A}
  [ ] {criterion 2 from contract A}
  [ ] {criterion 3 from contract B}

Steps to verify:
  1. {step 1 from contract A}
  2. {step 2 from contract A}
  3. {step 3 from contract B}

Does this match what you intended?
  1) Yes — mark complete and continue
  2) No — I'll describe what's wrong
  3) Partially — some criteria met, I'll describe the gaps
```

If only criteria exist, show criteria. If only steps exist, show steps. If both exist, show both. If neither exists (shouldn't happen due to validation warning, but handle gracefully), just ask "Does the output match your intent?"

### 4.4 Component 3: Rejection Flow

**Lives in:** `references/verification-tiers.md` (protocol) + `commands/gl/slice.md` (orchestrator integration)

When the user types anything other than "Yes" (option 1):

1. **Capture feedback.** Store the user's verbatim response.

2. **Present classification options.** The orchestrator presents:

```
Your feedback: "{user's response}"

How should we address this?

1) Tighten the tests -- the tests aren't specific enough to catch this mismatch
   (routes to: test writer adds more precise assertions, then implementer passes them)

2) Revise the contract -- the contract doesn't capture what I actually want
   (routes to: you update the contract, then the slice restarts)

3) Provide more detail -- I'll describe exactly what I expect
   (routes to: test writer uses your detail to write targeted tests, then implementer passes them)

Which option? (1/2/3)
```

3. **Route based on choice:**

| Choice | Internal classification | Action |
|--------|----------------------|--------|
| 1 (tighten tests) | Test gap | Spawn test writer with: rejection feedback (verbatim), contract, acceptance criteria. Test writer adds/tightens tests. Then spawn implementer to pass them. Re-run verification tier gate. |
| 2 (revise contract) | Contract gap | Present contract to user for revision. User edits acceptance criteria or contract definition. Restart slice from Step 1. |
| 3 (provide detail) | Implementation gap | Collect detailed description from user. Spawn test writer with: detailed description, rejection feedback, contract, acceptance criteria. Test writer writes targeted tests. Then spawn implementer. Re-run verification tier gate. |

4. **Increment rejection counter.** After any rejection, increment the per-slice rejection count.

### 4.5 Component 4: Rejection Counter and Escalation

**Lives in:** `references/verification-tiers.md`

The orchestrator tracks rejections per slice:

```yaml
slice_id: S-{N}
rejection_count: 0          # Increments on each non-"approved" response
rejection_log:
  - feedback: "{user's words}"
    classification: "test_gap"
    action_taken: "spawned test writer with feedback"
  - feedback: "{user's words}"
    classification: "implementation_gap"
    action_taken: "spawned test writer with detailed description"
```

When `rejection_count` reaches 3:

```
ESCALATION: {slice_name}

This slice has been rejected 3 times. The verification criteria may not match
what the contracts and tests can deliver.

Rejection history:
1. "{feedback 1}" -> {action taken}
2. "{feedback 2}" -> {action taken}
3. "{feedback 3}" -> {action taken}

Options:
1) Re-scope -- the contract is fundamentally wrong. Revise contracts and restart.
2) Pair -- provide detailed, step-by-step guidance for exactly what you want.
3) Skip verification -- mark this slice as auto-verified and proceed.
   (The mismatch is acknowledged but deferred.)

Which option? (1/2/3)
```

### 4.6 Component 5: Contract Schema Extension

**Lives in:** `agents/gl-architect.md` (contract format addition)

Three new optional fields added after the Security section in the contract format:

```markdown
**Verification:** verify
**Acceptance Criteria:**
- User can see campaign cards in a grid layout
- Each card shows campaign name, status, and key metric
- Cards are clickable and navigate to campaign detail

**Steps:**
- Run `npm run dev` and open localhost:3000
- Navigate to /campaigns
- Verify cards render in a 3-column grid
```

Field rules:
- `verification`: Optional. Values: `auto`, `verify`. Default: `verify`.
- `acceptance_criteria`: Optional under `verify`. List of behavioral statements the user can verify. Warn if empty when tier is `verify`.
- `steps`: Optional under `verify`. List of steps the user runs to verify the feature. Include when the how-to-verify isn't obvious. Warn if both `acceptance_criteria` and `steps` are empty when tier is `verify`.

### 4.7 CLAUDE.md Addition

Added to `src/CLAUDE.md` after the "Circuit Breaker" section:

```markdown
### Verification Tiers
- Every contract has a verification tier: `auto` or `verify` (default)
- After tests pass and verifier approves, the tier gate determines if human acceptance is required
- Rejection feedback routes to the test writer first -- if the implementation is wrong, the tests weren't tight enough
- Full protocol: `references/verification-tiers.md`
```

---

## 5. Data Model

No database entities. No persistent state files beyond what the orchestrator already maintains. The verification tier system operates within:

| Storage | What | Lifetime |
|---------|------|----------|
| Contract fields | verification, acceptance_criteria, steps | Permanent (in CONTRACTS.md) |
| Orchestrator context | Rejection count, rejection log, effective tier | Single /gl:slice execution |
| User feedback | Verbatim rejection text | Passed to test writer, then discarded |
| Checkpoint output | Verify checkpoint markdown | Displayed to user (ephemeral) |

---

## 6. Integration with Existing Systems

### 6.1 /gl:slice Pipeline

Changes to the pipeline:

| Step | Current | After Verification Tiers |
|------|---------|--------------------------|
| Step 6 (Verification) | Verifier checks contract coverage, stubs, wiring | Unchanged -- verifier also reports tier from contracts |
| Step 6b (NEW) | Does not exist | Verification tier gate: resolve tier, present checkpoint if verify, handle approval/rejection |
| Step 7 (Summary/docs) | Runs after verification passes | Runs after verification tier gate passes (approval received or tier is auto) |
| Step 9 (Visual checkpoint) | Reads config.workflow.visual_checkpoint | Deprecated -- log warning if visual_checkpoint is true, reference verification tiers instead |

### 6.2 Circuit Breaker (references/circuit-breaker.md)

No interaction. The circuit breaker handles implementation death spirals (Step 3). Verification tiers handle post-verification human acceptance (Step 6b). They operate on different steps in the pipeline and track different counters (attempt_count vs rejection_count).

### 6.3 Checkpoint Protocol (references/checkpoint-protocol.md)

The checkpoint type table gains one new type:

| Checkpoint Type | Trigger | When to Pause |
|-----------------|---------|---------------|
| Visual | ~~UI slice needs human eyes~~ Deprecated -- use verify tier | interactive mode only |
| Decision | Rule 4 architectural stop | always |
| External Action | Human action needed outside Claude | always |
| Circuit Break | 3 per-test or 7 per-slice failures | always |
| **Acceptance** | **Slice verification tier is verify** | **always (even in yolo mode)** |

Acceptance checkpoints always pause, even in yolo mode. The whole point of verification tiers is human confirmation -- skipping it defeats the purpose.

### 6.4 Agent Isolation (CLAUDE.md)

No changes to the isolation table. The test writer already cannot see implementation code. The rejection feedback is behavioral (user's words about what they expected vs. what they observed). The implementer still receives test names only.

### 6.5 Deviation Rules (references/deviation-rules.md)

No interaction. Deviation rules handle unplanned work during implementation. Verification tiers handle post-implementation acceptance. The rejection flow spawns agents through the normal /gl:slice pipeline, which already integrates deviation rules.

### 6.6 gl-verifier Agent

The verifier gains awareness of verification tiers:
- Read the `verification` field from each contract
- Include the effective tier in the verification report
- Flag contracts that are missing both `acceptance_criteria` and `steps` when tier is `verify`

This is informational -- the verifier reports tier status, it does not enforce the gate. The orchestrator enforces the gate.

---

## 7. File Changes Summary

| File | Change | Lines Added/Modified |
|------|--------|---------------------|
| `src/references/verification-tiers.md` | **NEW** | ~130 lines |
| `src/commands/gl/slice.md` | Add Step 6b, modify Step 9, add rejection handling | ~80 lines added/modified |
| `src/agents/gl-architect.md` | Add verification/acceptance_criteria/steps to contract format | ~25 lines added |
| `src/agents/gl-verifier.md` | Add tier awareness to verification report | ~15 lines added |
| `src/references/checkpoint-protocol.md` | Add Acceptance checkpoint type, deprecate Visual, update mode table | ~20 lines modified |
| `src/references/verification-patterns.md` | Add cross-reference to verification-tiers.md | ~5 lines added |
| `src/templates/config.md` | Add deprecation note on visual_checkpoint | ~5 lines added |
| `src/CLAUDE.md` | Add Verification Tiers rule | 4 lines added |
| `internal/installer/installer.go` | Add 1 manifest entry | 1 line added |
| `internal/installer/manifest_docs_test.go` | Update manifest count in tests | ~1 line modified |

**Total new content:** ~280 lines of markdown, 2 lines of Go.

---

## 8. Proposed Build Order

These are logical groupings for the architect to refine into slices with contracts:

1. **Schema Extension** -- Add `verification`, `acceptance_criteria`, `steps` fields to the contract format in `gl-architect.md`. Update the contract format template. Update `gl-verifier.md` to report tier in verification output. This is the foundation -- contracts must support tiers before anything else.

2. **Verification Gate** -- Add Step 6b to `/gl:slice`. Read tier from contracts, resolve effective tier (verify > auto, aggregation), present Verify checkpoint, handle "approved" response. Simple gate -- no rejection handling yet. Create `references/verification-tiers.md` with tier definitions and checkpoint format.

3. **Rejection Flow** -- Handle non-"approved" responses in the verification gate. Present structured classification options. Route to test writer (with behavioral feedback, contract, acceptance criteria) or to user for contract revision. Resume TDD loop after new tests.

4. **Rejection Counter** -- Track rejections per slice. Increment on each rejection. Escalation at 3 with options: re-scope, pair, skip. Add escalation format to `references/verification-tiers.md`.

5. **Documentation and Deprecation** -- Update `CLAUDE.md` with verification tier rule. Update `references/checkpoint-protocol.md` with Acceptance checkpoint type and Visual deprecation. Update `references/verification-patterns.md` with cross-reference. Add deprecation note to `templates/config.md` for `visual_checkpoint`. Modify Step 9 of `/gl:slice` to log deprecation warning.

6. **Architect Integration** -- Update `gl-architect.md` to generate `acceptance_criteria` and `steps` from requirements. Add guidance for tier selection: `auto` for pure data/logic/infrastructure, `verify` for everything else (it's the default anyway). Update `gl-designer.md` to capture verification preferences during design sessions.

---

## 9. Security

No new security surface. The verification tier system:
- Does not expose any external interfaces
- Does not handle user input beyond the existing slash command and "approved"/free-text response pattern
- Does not store credentials or sensitive data
- Does not modify agent isolation boundaries

The rejection flow has a **positive security property**: it prevents auto-completion of slices that may have security-relevant behavioral mismatches. A user reviewing acceptance criteria may catch issues that automated tests miss (e.g., "the API returns sensitive fields that shouldn't be visible").

---

## 10. Deferred

| Item | Why Deferred | When to Revisit |
|------|-------------|-----------------|
| Partial approval (per-criterion) | Adds UX complexity; one approval per slice is simpler and sufficient for MVP | When users report needing granular approval on large slices |
| Rejection history persistence | Currently ephemeral per-slice execution; no cross-session need identified | When users request post-mortem analysis of rejection patterns |
| AI-assisted gap classification | Orchestrator presents options, user chooses; AI classification adds unreliable complexity | When classification accuracy can be measured against user choices |
| Configurable rejection threshold (3) | Need real-world data to calibrate; matches circuit breaker threshold for consistency | After 20+ real-world rejections provide calibration data |
| Screenshot/visual diffing | Requires tooling outside Go CLI scope; human eyes are more reliable for MVP | When visual regression testing tools are integrated |
| Tier inference from contract content | Could auto-suggest verify for UI contracts, auto for infra | When architect consistently forgets to set tiers |

---

## 11. User Decisions (Locked)

| # | Gray Area | Decision | Rationale |
|---|-----------|----------|-----------|
| 1 | Tier count | Two tiers: `auto` and `verify`. No `demo` or `review` distinction. | The distinction between "check a list" and "walk through steps" is artificial. Both involve looking at output and confirming intent. Two tiers, one decision: "Can tests alone capture intent?" |
| 2 | Default tier | `verify` -- forgetting to set a tier gives a human checkpoint, not auto-approve | Safe default. The cost of an unnecessary verify is low. The cost of a missing verify is a completed slice that doesn't match intent. |
| 3 | Tier location | In the contract, not in GRAPH.json or config.json | The contract is the source of truth for what a boundary does and how it's verified. |
| 4 | Rejection routing | Always through test writer first, even for implementation gaps | TDD-correct. If the implementation is wrong and tests pass, the tests weren't tight enough. |
| 5 | Tier resolution | Highest tier wins + aggregation. verify > auto. One checkpoint per slice. | Minimizes user interruptions. A slice with mixed tiers gets the highest tier with all criteria aggregated. |
| 6 | Rejection counter | Per-slice, escalation at 3. | A slice is the unit of work. Mirrors the circuit breaker's slice-level ceiling. |
| 7 | visual_checkpoint deprecation | Keep in config but deprecate with warning. Tiers in contracts supersede it. | Default `verify` tier already provides human verification. The toggle is redundant but harmless. |
| 8 | Rejection feedback isolation | Test writer receives user's verbatim behavioral feedback, contract, and acceptance criteria. No implementation details. | Preserves agent isolation. User's words are about behavior, not implementation. |
| 9 | Gap classification UX | Orchestrator presents actionable options that implicitly map to gap types. User picks; orchestrator routes. | Users should not need to understand Greenlight's internal taxonomy. |
| 10 | Reference file location | New `references/verification-tiers.md` following the circuit-breaker.md pattern. | Different concern from existing verification-patterns.md. Automated vs human verification are separate topics. |
| 11 | Verify checkpoint content | Combined: `acceptance_criteria` + `steps` under one tier. Criteria = what to check. Steps = how, when not obvious. Both optional but warn if neither present. | Eliminates the artificial review/demo split. One checkpoint format handles both structured acceptance and interactive demos. |
