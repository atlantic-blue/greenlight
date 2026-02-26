# Checkpoint Protocol

When and how to pause for human input during automated execution.

## Checkpoint Types

### 1. Visual Checkpoint (deprecated)

> **Deprecated.** Use `**Verification: verify**` in contract definitions instead. The Acceptance checkpoint (Step 6b) now handles all human acceptance for slice output. The `visual_checkpoint` config key is also deprecated. See `references/verification-tiers.md` for the replacement.

**When:** A slice includes user-facing UI that automated tests can't fully verify.

**Trigger:** Slice contracts include any of:
- UI components, pages, or views
- CSS/layout changes
- Animations or transitions
- Responsive behaviour
- Accessibility requirements beyond automated checks

**Format:**
```
VISUAL CHECK

What was built: {slice name} — {description}

How to verify:
1. {command to start dev server or navigate to URL}
2. {what to look for}
3. {expected behaviour / appearance}

Type "approved" to continue, or describe issues.
```

**Rules:**
- Provide a single command to see the result (e.g., `npm run dev` → `http://localhost:3000/signup`)
- Be specific about what to look for — not "check if it works" but "form should have email + password fields, submit button should be disabled until both are filled"
- If user reports issues → spawn debugger to investigate, then re-implement

### 1a. Acceptance Checkpoint

**When:** A slice has one or more contracts with `**Verification:** verify` (or contracts without an explicit tier, which defaults to `verify`). This is the Step 6b gate in `/gl:slice`.

**Trigger:** The effective verification tier for the slice is `verify`. See `references/verification-tiers.md` for tier resolution rules.

**Rules:**
- This checkpoint always pauses. It cannot be skipped, even in yolo mode.
- The gate is blocking — the pipeline does not continue to Step 7 until the user approves.
- Rejection feedback routes to the test writer first (see Step 6b rejection flow).

**Format:**
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

Full protocol: see Step 6b in `commands/gl/slice.md` and `references/verification-tiers.md`.

### 2. Decision Checkpoint (rare)

**When:** An architectural decision (Rule 4 deviation) is needed during execution.

**Trigger:** Implementer or debugger reports a Rule 4 stop.

**Format:**
```
DECISION NEEDED

What happened: {description of what was found}
Why it matters: {impact on current slice and others}

Options:
A) {option}: {trade-offs}
B) {option}: {trade-offs}
C) Skip and revisit later: {what we lose}

Which approach?
```

**Rules:**
- Always provide at least 2 concrete options plus "skip and revisit"
- Explain trade-offs honestly — no default recommendation unless one option is clearly better
- After decision, update contracts if needed, then resume

### 3. External Action Checkpoint (very rare)

**When:** Something requires human action outside Claude Code.

**Trigger:** The system needs:
- Environment variables set (API keys, database URLs)
- External service configuration (OAuth app creation, DNS records)
- Manual process completion (App Store submission, domain purchase)
- Local tool installation (database server, runtime)

**Format:**
```
ACTION NEEDED

What's required: {specific action}
Why: {what depends on this}

Steps:
1. {step-by-step instructions}
2. {what value to copy/set}
3. {where to put it — e.g., .env file}

Type "done" when complete, or "skip" to continue without.
```

**Rules:**
- Provide exact steps — don't assume the user knows how to create an OAuth app
- Specify exactly what value/credential is needed and where to put it
- If the action can be deferred, say so
- After "done", verify the action worked (e.g., test the API key)

## Checkpoint Behaviour by Mode

| Mode | Visual (deprecated) | Acceptance | Decision | External Action |
|------|---------------------|------------|----------|-----------------|
| interactive | Pause and ask | Always pause | Pause and ask | Pause and ask |
| yolo | Skip (log warning) | Always pause | Pause and ask | Pause and ask |

**Even in YOLO mode, Acceptance checkpoints always pause.** Acceptance checkpoints are blocking regardless of mode — they cannot be skipped. Only visual checks (deprecated) can be skipped in yolo mode — and they're logged so the user can review later.

## Anti-Patterns

### Too Many Checkpoints

```
BAD: Pausing after every file is written
BAD: Asking "should I continue?" after each step
BAD: Checkpoint for non-visual internal code changes
```

Checkpoints should be rare. The test runner handles verification. Only pause for things tests can't verify.

### Vague Checkpoints

```
BAD: "Check if it looks right"
BAD: "Verify the page works"
BAD: "Make sure everything is correct"
```

Be specific. The user should know exactly what to look at and what "correct" means.

### Checkpoints That Block Unnecessarily

```
BAD: Pausing for a missing env var that has a sensible default
BAD: Pausing for a visual check on a non-visual slice
BAD: Pausing to ask about a decision the contracts already answer
```

If the answer is already in the contracts, don't ask again.

## Checkpoint Logging

All checkpoints are logged in the slice summary:

```markdown
## Checkpoints
| Type | Description | Resolution | Duration |
|------|-------------|------------|----------|
| visual | Registration form layout | approved | 2m |
| decision | Auth: JWT vs session | JWT chosen (stateless) | 5m |
```

This helps track where time is spent and whether checkpoints are adding value.

## File-Per-Slice State Context

When saving or restoring checkpoint state context, the state format determines which files to read and write:

- In **file-per-slice mode**: read current slice state from `.greenlight/slices/{id}.md` and project context from `project-state.json`. Write state updates to the individual slice file, not STATE.md.
- In **legacy mode**: read and write STATE.md as before (no change).

State format detection (see `references/state-format.md`) determines which path to use. Checkpoint save/restore works with both formats — the checkpoint protocol itself is format-agnostic.
