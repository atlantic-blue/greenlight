# Verification Tiers

Authoritative protocol document for the Greenlight verification tier system. This document defines how slices are accepted by humans and how the acceptance checkpoint is enforced.

---

## 1. Tier Definitions

There are exactly two tiers: `auto` and `verify`. No third tier exists. Only: auto and verify.

**Default: verify.** A contract without an explicit `**Verification:**` field is treated as if it has `**Verification:** verify`.

### `auto`

Tests pass → slice proceeds directly to summary and documentation (Step 7). No human acceptance checkpoint is required.

Use `auto` for:
- Infrastructure changes with no user-visible output
- Refactoring that is fully covered by existing tests
- Slices where human review adds no value beyond green tests

### `verify`

Tests pass → human acceptance checkpoint is presented before proceeding to Step 7.

Use `verify` (the default) for:
- Any slice with user-facing behaviour
- New features, new endpoints, new UI
- Any change to observable system output
- Slices where a human needs to confirm the output matches their intent

---

## 2. Tier Resolution

A slice may have multiple contracts, each with its own `**Verification:**` field. The orchestrator computes one effective tier for the entire slice using this rule:

**verify > auto (highest wins)**

That is: if any contract on the slice has `**Verification:** verify`, the effective tier for the slice is `verify`, regardless of how many contracts have `auto`.

Rules:
- Acceptance criteria from all contracts are aggregated into one combined list.
- Steps from all contracts are aggregated into one combined numbered list.
- There is one checkpoint per slice — never one per contract.
- Per-slice tracking: the rejection counter and acceptance decision are per-slice, not per-contract.

---

## 3. Verify Checkpoint Format

When the effective tier is `verify`, present the following checkpoint to the user:

```
ALL TESTS PASSING -- Slice {slice_id}: {slice_name}

Please verify the output matches your intent.

Acceptance criteria:
  [ ] {criterion 1}
  [ ] {criterion 2}

Steps to verify:
  1. {step 1}
  2. {step 2}

Does this match what you intended?
  1) Yes -- mark complete and continue
  2) No -- I'll describe what's wrong
  3) Partially -- some criteria met, I'll describe the gaps
```

The `ALL TESTS PASSING` header is always shown. The criteria and steps blocks are conditional (see format adaptation rules below).

### Response Interpretation

- `"1"` or `"Yes"` (case-insensitive) → approved → proceed to Step 7
- Anything else → rejected → enter rejection flow, increment rejection counter, re-present checkpoint after resolution

### Format Adaptation

- **Only criteria (no steps):** Show the `Acceptance criteria:` block. Omit the `Steps to verify:` block entirely.
- **Only steps (no criteria):** Show the `Steps to verify:` block. Omit the `Acceptance criteria:` block entirely.
- **Both criteria and steps:** Show both blocks.
- **Neither criteria nor steps:** Show a simplified prompt: `Does the output match your intent?` with the three response options only.

---

## 4. Rejection Flow

When the user does not approve (response is not `"1"` or `"Yes"`), enter the rejection flow.

The orchestrator presents three options to the user:

1. **Tighten tests** — the tests did not adequately cover the expected behaviour. Return to the test writer (gl-test-writer) with behavioural feedback only.
2. **Revise contract** — the contract did not accurately specify the intended outcome. Update the contract before retrying.
3. **Provide more detail** — the implementation is close but needs adjustment. The user provides specific guidance and additional detail; spawn a fresh implementer with that guidance as additional context.

After the user selects an option and the remediation is performed, re-run Step 6b from the beginning (re-present the checkpoint after tests are green again).

### Classification Mapping

Each user choice maps to an internal classification:

| Choice | Internal classification | Route |
|--------|------------------------|-------|
| 1      | `test_gap`             | Spawn gl-test-writer with rejection context |
| 2      | `contract_gap`         | Enter contract revision flow; restart slice from Step 1 |
| 3      | `implementation_gap`   | Collect additional detail from user, then spawn gl-test-writer |

User feedback is treated as behavioral context only — it is never executed as code or used to modify files directly.

---

## 5. Rejection Counter

The rejection counter is tracked **per-slice**, not per-contract. It is in-memory only — not written to disk. It resets to 0 on every new `/gl:slice` execution (re-invoked). Within a single execution, the counter persists across all rejection loops.

### State Structure

Initialize this state at the start of each `/gl:slice` execution:

```yaml
slice_id: "{slice_id}"
rejection_count: 0
rejection_log: []
```

### Counter Invariants

- Incremented by 1 for each rejection at Step 6b (every non-approved response).
- **Escalation at 3 rejections:** When `rejection_count >= 3`, trigger escalation immediately.
- Contract revision (option 2 — contract_gap) restarting from Step 1 does NOT reset the rejection counter. Every rejection path increments the same counter.
- The rejection counter does not interact with the circuit breaker. They are separate concerns.

### Log Entry Structure

Each rejection appends one entry to `rejection_log` in chronological order:

```yaml
- attempt: {rejection_count}
  feedback: "{verbatim user response}"
  classification: "{test_gap | contract_gap | implementation_gap}"
  action_taken: "{description of remediation taken}"
```

### Error Handling

- `CounterOverflow`: If the counter exceeds 3 without triggering escalation, trigger escalation immediately and warn the user.
- `LogCorruption`: If the rejection_log becomes inconsistent or corrupted, reset the log to `[]`, preserve the current rejection_count, and warn the user.

### Escalation Format

When `rejection_count` reaches 3, present:

```
ESCALATION -- {slice_name}

This slice has been rejected 3 times. Continued iteration without intervention is unlikely to converge.

Rejection history:
  1. Feedback: "{verbatim feedback 1}" | Action: {action_taken_1}
  2. Feedback: "{verbatim feedback 2}" | Action: {action_taken_2}
  3. Feedback: "{verbatim feedback 3}" | Action: {action_taken_3}

How would you like to proceed?

  1) Re-scope — Reset the rejection counter and restart the slice from scratch
  2) Pair — Collect step-by-step guidance, then spawn gl-test-writer; resets the rejection counter
  3) Skip verification — Mark as auto-verified and proceed (does not reset the counter)

Which option? (1/2/3)
```

**Option routing:**

| Option | Action | Counter |
|--------|--------|---------|
| 1 (Re-scope) | Restart slice from scratch with revised scope | Reset to 0 |
| 2 (Pair) | Collect guidance, spawn gl-test-writer | Reset to 0 |
| 3 (Skip) | Mark as auto, proceed to Step 7 | Not reset |

**Skip creates an explicit log entry:**
```
Verification skipped by user. Mismatch acknowledged and deferred.
```

**Error handling:**
- `InvalidEscalationChoice`: If the user enters anything other than 1, 2, or 3, re-prompt: "Please choose 1, 2, or 3."
- `EmptyRejectionLog`: If escalation triggers but the log is empty, display escalation without the rejection history block and warn that history is unavailable.

---

## 6. Agent Isolation in the Rejection Loop

When the rejection flow routes to the test writer (option 1 — tighten tests), the test writer receives:

- Behavioural feedback only: what the user observed, what they expected, which acceptance criteria were not met.
- The test writer does NOT receive: implementation code, implementation decisions, or the existing test source.

This preserves the test writer's isolation. Sending implementation details would cause the test writer to test implementation rather than behaviour.

The orchestrator is responsible for filtering the feedback before passing it to gl-test-writer. The test writer receives behavioral feedback only — never implementation source code or test source code.

When the rejection flow then routes to the implementer (after the test writer has written new tests), the implementer receives test names only — not test source code. This preserves the implementer's isolation and ensures the implementer implements against the contract, not against the test internals.

---

## 7. Backward Compatibility

### Contracts Without the Verification Field

Contracts written before the `**Verification:**` field was introduced do not have it. These contracts **default to `verify`** — the same as if `**Verification:** verify` were explicitly set. This is the safe default: require human sign-off unless explicitly opted out.

### `visual_checkpoint` Config Key (Deprecated)

The `visual_checkpoint` key in `.greenlight/config.json` is deprecated. It was previously used to trigger visual review for UI slices.

If `config.workflow.visual_checkpoint` is `true`:
- Log a deprecation warning: `"visual_checkpoint is deprecated. Use **Verification: verify** in your contracts instead."`
- Do not suppress the Step 6b checkpoint.
- The verification tier system (Step 6b) now handles human acceptance for all slice types.

New projects should omit `visual_checkpoint` from config entirely. Existing projects should migrate to per-contract `**Verification:**` fields.

---

## Invariants

- Exactly two tiers: `auto` and `verify`. No extensions.
- Default is `verify` when the field is absent.
- The test writer always runs before the human checkpoint (checkpoint cannot bypass tests).
- Escalation threshold is 3 rejections per slice.
- Acceptance checkpoints pause even in yolo mode — the checkpoint is always blocking and cannot be skipped.
- The verification tier system is additive to the pipeline — it does not replace or interact with the circuit breaker.
- The checkpoint does not interact with circuit breaker state.
