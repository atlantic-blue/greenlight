---
name: gl:slice
description: Build a vertical slice using TDD — tests → implement → security scan → verify → green → commit
argument-hint: "<slice-number>"
allowed-tools: [Read, Write, Bash, Glob, Grep, Task, AskUserQuestion]
---

# Greenlight: Build Slice

Execute a single vertical slice through the full TDD loop:
contracts → tests → implementation → security → verification → complete.

**Read first:**
- `CLAUDE.md` — engineering standards
- `.greenlight/STATE.md` — current state
- `.greenlight/GRAPH.json` — dependency graph
- `.greenlight/CONTRACTS.md` — contracts
- `.greenlight/config.json` — settings (model profiles, workflow toggles)

**Read references:**
- `references/deviation-rules.md` — for deviation handling context
- `references/checkpoint-protocol.md` — for checkpoint decisions
- `references/verification-patterns.md` — for verification context

## Model Resolution

Before spawning any agent, resolve its model from `.greenlight/config.json`:

1. Check `model_overrides[agent_name]` — if set, use it
2. Else check `profiles[model_profile][agent_name]` — use profile default
3. Else fall back to `sonnet`

Agents spawned by this command: `test_writer`, `implementer`, `security`, `verifier`.

## Pre-flight

### 1. Validate Slice

```bash
# Read GRAPH.json, find the requested slice
cat .greenlight/GRAPH.json
```

Checks:
- Slice number/ID exists in GRAPH.json
- All dependencies are status "complete" in STATE.md
- Slice is not already complete

If any check fails → stop and explain:
```
Cannot start slice {N}:
- Dependency "{dep_name}" (slice {dep_id}) is not complete
- Complete it first: /gl:slice {dep_id}
```

### 2. Load Context

From GRAPH.json:
- Slice ID, name, description
- Contract names for this slice
- Dependencies (already verified complete)
- **`wraps` field** (optional) — array of boundary names from STATE.md Wrapped Boundaries table

From CONTRACTS.md:
- Full contract definitions for this slice's contracts
- Input/output types, error states, invariants, security requirements
- If slice has `wraps` field: read the `[WRAPPED]` contracts for those boundaries

From config.json:
- Model assignments for each agent
- Workflow toggles (security_scan, visual_checkpoint, etc.)
- Test commands

### 2a. Detect Wraps Field (Locking-to-Integration)

```bash
# Check if slice has a wraps field
# Example GRAPH.json entry:
# { "id": "S-XX", "wraps": ["auth", "payments"], ... }
```

If the slice has a `wraps` field:

1. **Read locking tests** for each wrapped boundary:
   ```bash
   ls tests/locking/{boundary-name}.test.* 2>/dev/null
   ```

2. **Extract locking test names** (NOT source code — names only):
   ```bash
   grep -n "func Test\|it('\|it(\"" tests/locking/{boundary-name}.test.* | sed 's/.*func //' | sed 's/(.*//'
   ```
   These become context for the test writer: "these are the existing locked behaviours."

3. **Read [WRAPPED] contracts** for each wrapped boundary from CONTRACTS.md.

4. **Verify wrapped boundaries exist** in STATE.md with status `wrapped`:
   - If boundary not found in STATE.md → warn user, proceed without locking context (treat as greenfield)
   - If locking test file not found → warn user, proceed without locking context

Report:
```
Slice has wraps field: [{boundary_names}]

Locking context loaded:
  {boundary_name}: {N} locking tests, [WRAPPED] contract found
  {boundary_name}: {N} locking tests, [WRAPPED] contract found

This is a locking-to-integration transition. After verification:
  - Locking tests will be deleted
  - [WRAPPED] tags will be removed from contracts
  - STATE.md boundary status will change to "refactored"
```

### 3. Check for Previous Attempt

```bash
ls tests/integration/{slice-id}.test.* 2>/dev/null
ls .greenlight/.continue-here.md 2>/dev/null
```

If tests exist from a previous attempt → ask:
```
Previous tests found for this slice. Options:
1) Continue from existing tests (skip test writing)
2) Start fresh (rewrite tests)
```

If .continue-here.md exists → suggest `/gl:resume` instead.

### 4. Update State

Update `.greenlight/STATE.md`:
- Current slice → this slice
- Step → "tests"
- Last activity → today's date

---

## Checkpoint: Create Rollback Tag

After pre-flight completes, create a checkpoint tag under the greenlight/checkpoint/ namespace so the implementer can roll back if needed.

The tag name follows the pattern: `greenlight/checkpoint/{slice_id}`

```bash
# Idempotent: remove any pre-existing tag for this slice before re-creating
git tag -d greenlight/checkpoint/{slice_id} 2>/dev/null || true

# Create fresh checkpoint tag at the current HEAD
git tag greenlight/checkpoint/{slice_id}
```

Report:
```
Checkpoint created: greenlight/checkpoint/{slice_id}
```

This tag is passed to the implementer as context and cleaned up on successful completion.

---

## Step 1: Write Tests (Agent A — Fresh Context)

The test writer NEVER sees implementation code. Send ONLY contracts, stack info, and existing test patterns.

**If slice has `wraps` field:** Include locking test NAMES (not source code) as additional context. The test writer uses these to ensure integration tests cover at least all locked behaviours (superset requirement).

```
Task(prompt="
Read agents/gl-test-writer.md
Read CLAUDE.md

<slice>
ID: {slice_id}
Name: {slice_name}
Description: {what user can do after this}
</slice>

<contracts>
{FULL contract definitions for this slice — types, interfaces, error states, invariants}
{If wraps slice: include [WRAPPED] contracts for the wrapped boundaries}
</contracts>

<stack>
Language: {language}
Test framework: {test_framework}
Assertion library: {assertion_library}
</stack>

<existing_tests>
{summary of test files from prior slices — patterns, setup approach, NOT implementation}
</existing_tests>

<test_fixtures>
{existing factory functions from tests/fixtures/}
</test_fixtures>

{IF WRAPS SLICE — include this additional block:}
<locked_behaviours>
This slice refactors wrapped boundaries. The following locking test NAMES
describe behaviours that are currently locked. Your integration tests MUST
cover at least all of these behaviours (superset requirement).

Boundary: {boundary_name}
Locking tests:
  - [LOCK] {test name 1}
  - [LOCK] {test name 2}
  - [LOCK] {test name 3}
  ...

NOTE: You are receiving test NAMES only, not source code. This preserves
isolation. Use these names to understand what behaviours exist, then write
integration tests that cover all of them plus any new behaviours from the
contracts.
</locked_behaviours>

Write integration tests for this slice. Create:
1. tests/integration/{slice-id}.test.{ext}
2. New factories in tests/fixtures/ if needed

Do NOT write implementation code.
", subagent_type="gl-test-writer", model="{resolved_model.test_writer}", description="Write tests for slice {slice_id}")
```

### Validate Test Writer Output

After the test writer returns:
1. Verify test file exists: `tests/integration/{slice-id}.test.{ext}`
2. Count tests written
3. Check against contract coverage (every contract should have tests)

If test file missing or empty → report error, do not proceed.

---

## Step 2: Confirm Tests Fail

Run the new tests:

```bash
{config.test.command} {config.test.filter_flag} {slice-id}
```

**Expected: ALL new tests FAIL.**

| Result | Action |
|--------|--------|
| All fail | Good — proceed to implementation |
| Some pass | Investigate — tests may be trivial or there's leftover implementation. If trivial, flag for test writer improvement. If leftover, proceed. |
| All pass | Nothing to implement. Verify this is expected. If yes, skip to verification. |
| Test runner errors (syntax, import) | Fix test infrastructure issues (NOT test logic). Re-run. |

Report to user:
```
Tests written: {N}
All failing as expected

Tests cover:
- {contract 1}: {N} tests (success, validation, errors)
- {contract 2}: {N} tests (success, validation)

Proceeding to implementation...
```

Commit tests:
```bash
git add tests/integration/{slice-id}.test.*
git add tests/fixtures/*.* # if new/modified
git commit -m "test({slice-id}): add integration tests

- {N} tests covering {contracts}
- All tests currently failing (pre-implementation)
"
```

---

## Step 3: Implement (Agent B — Fresh Context)

The implementer gets test names but NOT test code. This is critical for TDD integrity.

**Extract test names from the test file:**
```bash
grep -n "it('\|it(\"" tests/integration/{slice-id}.test.* | sed 's/.*it(//' | sed "s/[',\"].*//"
```

```
Task(prompt="
Read agents/gl-implementer.md
Read CLAUDE.md
Read references/deviation-rules.md

<slice>
ID: {slice_id}
Name: {slice_name}
</slice>

<contracts>
{FULL contract definitions}
</contracts>

<test_expectations>
{test names/descriptions ONLY — extracted above}
</test_expectations>

<existing_code>
{relevant source files from prior slices — read and include}
</existing_code>

<stack>
Language: {language}
Framework: {framework}
Database: {database}
Patterns: {patterns from prior slices}
</stack>

<checkpoint>
checkpoint_tag: greenlight/checkpoint/{slice_id}
Rollback: orchestrator restores from checkpoint_tag on CIRCUIT BREAK
</checkpoint>

<test_command>
{config.test.command} {config.test.filter_flag} {slice-id}
</test_command>

<full_test_command>
{config.test.command}
</full_test_command>

Implement production code to make all tests pass.
Run tests after each implementation group.
Run full suite at the end — no regressions allowed.
Commit with conventional format.
", subagent_type="gl-implementer", model="{resolved_model.implementer}", description="Implement slice {slice_id}")
```

### Handle Implementation Result

**Success (all tests pass):**
- Update STATE.md step → "security" (if security_scan enabled) or "verifying"
- Proceed to security scan or verification

**Partial failure (some tests fail after max retries):**
```
Implementation incomplete: {N}/{total} tests passing

Failing:
- {test name}: {failure reason}

Options:
1) Spawn fresh implementer with failure context
2) Spawn debugger to investigate
3) Pause and review manually (/gl:pause)
```

Max 3 implementation attempts. After 3 failures → CIRCUIT BREAK.

**CIRCUIT BREAK — implementation could not be completed after 3 attempts:**

```
CIRCUIT BREAK

Slice {slice_id} implementation failed after {N} attempts.

Failing tests:
- {test name}: {failure reason}

Recovery options:
1) guidance + retry — provide additional guidance, roll back to the checkpoint, and spawn a fresh implementer with the guidance as additional context

2) debugger — spawn gl-debugger to investigate the root cause before retrying

3) Pause — stop here and review manually (/gl:pause)

Which option?
```

On retry (option 1): roll back to the checkpoint tag (see checkpoint protocol below), reset attempt counters, and spawn a fresh implementer with the guidance as additional context.
On debugger (option 2): spawn gl-debugger with failure context, then resume from the implementer step.
On pause (option 3): save state and halt.

**Architectural stop (Rule 4 deviation):**
Present the stop to the user using checkpoint protocol:
```
DECISION NEEDED

The implementer encountered an architectural issue:
{description from implementer}

Impact: {what it affects}

Options:
A) {option from implementer}
B) {option from implementer}
C) Pause and think about it

Which approach?
```

After decision → spawn fresh implementer with the decision as context.

---

## Step 4: Verify All Tests Green

Run FULL test suite (not just this slice):

```bash
{config.test.command}
```

| Result | Action |
|--------|--------|
| All pass | Proceed to security scan |
| This slice's tests fail | Spawn fresh implementer with failure output (max 3 attempts) |
| Other slice's tests fail | Regression detected. Spawn fresh implementer to fix WITHOUT modifying other tests |
| Locking tests fail (wraps slice) | Regression: existing locked behaviour broken. Implementer MUST fix — locking tests serve as guardrails during refactoring |
| Infrastructure error | Fix infrastructure, re-run |

**For wraps slices:** Both locking tests AND new integration tests must pass at this point. Locking tests are the safety net — they prove the refactored code still does what the original code did. They are only deleted AFTER verification succeeds in Step 6.

---

## Step 5: Security Scan (Agent C — Fresh Context)

**Skip if `config.workflow.security_scan` is false.**

```bash
# Get the diff for this slice
git diff HEAD~{N}..HEAD  # N = commits in this slice
```

```
Task(prompt="
Read agents/gl-security.md
Read CLAUDE.md

<slice>
ID: {slice_id}
Name: {slice_name}
</slice>

<contracts>
{contracts with security requirements highlighted}
</contracts>

<diff>
{git diff output}
</diff>

<files_changed>
{list of new/modified files with paths}
</files_changed>

Review this slice for security vulnerabilities.
For each issue found, write a FAILING test in tests/security/{slice-id}-security.test.{ext}

{IF WRAPS SLICE — include this additional block:}
<known_security_issues>
This slice refactors a previously wrapped boundary. The following security
issues were documented during wrapping (from the [WRAPPED] contract's
Security section):

{list of known issues from wrapped contract Security section}

Check whether these known issues have been addressed by the refactoring.
- NEW vulnerabilities introduced by the refactoring: write FAILING tests (normal behaviour)
- Pre-existing issues that persist unchanged: FLAG but do NOT block (they were known before)
- Pre-existing issues that have been fixed: note as resolved
</known_security_issues>
", subagent_type="gl-security", model="{resolved_model.security}", description="Security scan for slice {slice_id}")
```

### Handle Security Results

**No issues found:**
```
Security scan: PASS
No vulnerabilities detected in {N} files reviewed.
```

**Issues found:**
```
Security scan found {N} issues:
- [{severity}] {description}
- [{severity}] {description}

{N} security tests written. Running fixes...
```

If security tests were written:
1. Commit the security tests
2. Spawn implementer to make security tests pass
3. Run full suite — all tests must be green (functional + security)
4. Max 2 security fix attempts. After 2 failures → pause.

---

## Step 6: Verification (Agent D — Fresh Context)

**This step validates the slice actually delivers what was promised.**

```
Task(prompt="
Read agents/gl-verifier.md
Read references/verification-patterns.md

<slice>
ID: {slice_id}
Name: {slice_name}
Description: {what user can do after this}
</slice>

<contracts>
{all contracts for this slice}
</contracts>

<test_results>
{output from final test run}
</test_results>

<files_changed>
{all files created/modified in this slice}
</files_changed>

<mode>slice</mode>

Verify this slice satisfies its contracts. Check test coverage, implementation substance, and wiring.
", subagent_type="gl-verifier", model="{resolved_model.verifier}", description="Verify slice {slice_id}")
```

### Handle Verification Result

**PASS:**
Proceed to completion.

**PASS with warnings:**
Log warnings, proceed to completion. Warnings don't block.

**FAIL:**
```
Verification failed:
{issues from verifier report}

Fixing {N} issues...
```

Route each failure to the appropriate fix:
- Missing test coverage → flag for next slice or add to quick task backlog
- Stub detected → spawn fresh implementer
- Dead code (not wired) → spawn fresh implementer to wire it
- Standards violation → spawn fresh implementer to refactor

After fixes → re-run verification. Max 2 verification cycles.

---

## Step 6a: Locking-to-Integration Transition (wraps slices only)

**Skip if slice has no `wraps` field.**

After verification succeeds and both locking tests AND integration tests pass:

### 1. Delete Locking Tests

```bash
# For each boundary in the wraps field:
rm tests/locking/{boundary-name}.test.{ext}

# If tests/locking/ is now empty:
rmdir tests/locking/ 2>/dev/null
```

Print: `Locking tests removed: tests/locking/{boundary-name}.test.{ext}`

### 2. Remove [WRAPPED] Tag from Contracts

For each wrapped boundary in CONTRACTS.md:
- Find the contract heading: `### Contract: {BoundaryName} [WRAPPED]`
- Remove the `[WRAPPED]` tag: `### Contract: {BoundaryName}`
- Remove the `**Source:**`, `**Wrapped on:**`, and `**Locking tests:**` metadata lines
- Change `**Slice:** wrappable` to `**Slice:** {current_slice_id}`

Print: `[WRAPPED] tag removed from contract: {BoundaryName}`

### 3. Update STATE.md Wrapped Boundaries

For each boundary in the wraps field:
- Change status from `wrapped` to `refactored`
- Optionally add note: `(replaced by slice {slice_id})`

### 4. Run Full Test Suite Again

Confirm all tests still pass after locking test deletion:

```bash
{config.test.command}
```

If tests fail after locking test deletion → something is wrong. Do NOT proceed. Report error.

### 5. Commit Transition

```bash
git add tests/locking/ .greenlight/CONTRACTS.md .greenlight/STATE.md
git commit -m "refactor({slice_id}): transition {boundary_names} from locking to integration

- Locking tests deleted: {list}
- [WRAPPED] tags removed from contracts: {list}
- STATE.md boundaries marked as refactored
"
```

Report:
```
Locking-to-integration transition complete:
  {boundary_name}: locking tests deleted, [WRAPPED] tag removed, status → refactored
```

---

## Step 6b: Verification Tier Gate

**This gate runs after Step 6 (Verification) and Step 6a (Locking-to-Integration Transition). It is blocking — the pipeline does not continue to Step 7 until the gate passes. This applies even in yolo mode. Acceptance checkpoints always pause, regardless of workflow settings.**

### 1. Read Verification Tiers

For each contract in this slice, read the `**Verification:**` field.

- Valid values: `auto` or `verify`
- Missing field: default to `verify` (safe default)
- If `visual_checkpoint` is set in config.json and is true → log a deprecation warning: `"visual_checkpoint is deprecated. Use **Verification: verify** in your contracts instead."` Treat the effective tier as if `verify` is set for backward compatibility.

### 2. Compute Effective Tier

Apply the rule: **verify > auto** (highest wins).

If any contract has `**Verification:** verify`, the effective tier is `verify`.
If all contracts have `**Verification:** auto`, the effective tier is `auto`.

### 3. Execute Based on Effective Tier

**If effective tier is `auto`:**

```
Log: "Verification tier: auto. Skipping acceptance checkpoint."
```

Proceed to Step 7.

**If effective tier is `verify`:**

Aggregate acceptance criteria from all contracts. Aggregate steps from all contracts. Present the checkpoint below and wait for a human response. The gate is blocking.

### Checkpoint Format

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

### Response Handling

- `"1"` or `"Yes"` (case-insensitive) → approved → proceed to Step 7
- Anything else → rejected → enter rejection flow (see below), increment rejection counter, re-run Step 6b

### Format Adaptation Rules

- If only criteria (no steps): show the `Acceptance criteria:` block; omit `Steps to verify:` entirely.
- If only steps (no criteria): show the `Steps to verify:` block; omit `Acceptance criteria:` entirely.
- Both present: show both blocks.
- Neither present: use simplified prompt — `Does the output match your intent?` with the three response options only.

### Rejection Flow

When the user does not approve, present three options:

```
What would you like to do?

1) Tighten tests — tests didn't cover the expected behaviour (return to test writer with behavioral feedback only)
2) Revise contract — the contract didn't specify the intended outcome correctly
3) Provide more detail — implementation is close, describe what needs to change
```

After remediation → re-run Step 6b.

Increment the per-slice rejection counter. After 3 rejections, escalate:

```
Slice {slice_id} has been rejected {N} times.

Summary of rejections:
{list of what was found unsatisfactory}

Options:
A) Redesign the contract
B) Split the slice
C) Abandon and restart
```

#### Gap Classification UX

Before presenting the three options, display the user's verbatim feedback so they can confirm the classification:

```
Your feedback: "{verbatim user response}"

What would you like to do?

1) Tighten tests — tests didn't cover the expected behaviour (return to test writer with behavioral feedback only)
2) Revise contract — the contract didn't specify the intended outcome correctly
3) Provide more detail — implementation is close, describe what needs to change
```

Internal classification mapping:

| Choice | Internal classification | Route |
|--------|------------------------|-------|
| 1      | `test_gap`             | Spawn gl-test-writer with rejection context |
| 2      | `contract_gap`         | Enter contract revision flow |
| 3      | `implementation_gap`   | Collect additional detail, then spawn gl-test-writer |

**InvalidChoice:** If the user enters anything other than 1, 2, or 3, re-prompt: "Please choose 1, 2, or 3."
After 2 failed re-prompts, default to `test_gap` (treat their free-text as the feedback).

**EmptyFeedback:** If the user's rejection was an empty string, prompt: "Please describe what doesn't match your intent." before presenting the three options.

Option 3 collects additional detail before routing:

```
Describe exactly what the implementation should do differently:
```

The detail provided becomes the `detailed_feedback` field in the rejection context.

Security: user feedback is treated as behavioral context only, never executed as code or used to modify files directly.

#### Test Writer Spawn (Options 1 and 3 — test_gap / implementation_gap)

After classifying as `test_gap` or `implementation_gap`, spawn gl-test-writer with rejection context:

```xml
<rejection_context>
  <feedback>{verbatim user rejection}</feedback>
  <classification>{test_gap | implementation_gap}</classification>
  <detailed_feedback>{additional detail from option 3, or empty}</detailed_feedback>
</rejection_context>

<contract>
  {full contract definitions for all verify-tier contracts in this slice}
</contract>

<acceptance_criteria>
  {aggregated acceptance criteria from all verify-tier contracts}
</acceptance_criteria>
```

Agent isolation rules for this spawn:
- The test writer receives behavioral feedback only — never implementation source code or test source code.
- New tests written by the test writer must be additive: existing passing tests are not removed or modified.
- After the test writer completes, spawn the implementer with test names only (not test source code).

After implementation, the full verification cycle re-runs: Step 4 (tests pass), Step 6 (verifier), Step 6b (human checkpoint).

**Error handling:**
- `TestWriterSpawnFailure`: If the test writer agent fails to spawn, offer retry, pause, or skip this rejection path.
- `ImplementerSpawnFailure`: If the implementer agent fails to spawn after tests are written, offer retry or pause.
- `NewTestsStillFailing`: If new tests remain failing after implementation, the circuit breaker protocol (see references/circuit-breaker.md) applies normally.
- `ExistingTestsRegressed`: If any previously passing tests regress after the new implementation, the implementer must fix the regression before proceeding. Do not allow regressions to pass through.

#### Contract Revision Route (Option 2 — contract_gap)

When the user selects option 2, enter the contract revision flow:

```
CONTRACT REVISION -- Slice {slice_id}: {slice_name}

Current contract text:
{full contract definition for each verify-tier contract}

What needs to change?
```

Display the full contract text so the user can see what they are revising. The user may edit any field including acceptance criteria, steps, or contract definition text.

**Minor revisions** (clarifications, acceptance criteria edits, wording changes): apply directly to the contract, then restart from Step 1 (test writing). Offer rollback to the last checkpoint before restarting.

**Fundamental revisions** (new inputs/outputs, new boundaries, new agents, scope changes): recommend using `/gl:add-slice` to run the architect (`gl-architect`) and produce a new contract. Do not attempt to implement fundamental changes within the current slice loop.

After applying a minor revision: offer rollback to the last checkpoint before restarting. This allows the user to recover the current implementation if the revision turns out to be wrong.

Contract revision increments the rejection counter by 1 (same as all other rejection paths).

**EmptyRevision:** If the user provides no revision description, re-prompt: "Please describe what the contract should say instead."

---

## Step 7: Generate Summary and Update Documentation

After verification passes, generate a summary and update project documentation.

### 7a: Summary Generation (C-41)

Spawn a Task with fresh context to write the slice summary:

```
Task(prompt="
Collect and document the following information for slice {slice_id}:

<slice_data>
ID: {slice_id}
Name: {slice_name}
Description: {what user can do after this}
Contracts satisfied: {list}
Test count: {N}
Security test count: {M}
Files created/modified: {list}
Verification status: PASS
</slice_data>

Write a summary to `.greenlight/summaries/{slice-id}-SUMMARY.md` in product language (not implementation language).

Check if architecture changed (new boundaries, integrations, or patterns). If so, note this in the summary.

Summary failure does not block the pipeline.
", subagent_type="gl-summarizer", model="{resolved_model.summarizer}", description="Generate summary for slice {slice_id}")
```

If Task fails, log warning and continue. Summary generation is non-blocking.

### 7b: Decision Aggregation (C-44)

After verification, aggregate decision notes from all agents:

1. Collect decision notes from test writer, implementer, security, and verifier outputs
2. Filter for meaningful decisions (architectural choices, tradeoffs, security decisions)
3. Format as DECISIONS.md table rows:
   - Columns: #, Decision, Context, Chosen, Rejected, Date, Source
   - Source format: `slice:{slice-id}`
4. append to DECISIONS.md (append-only, never modify existing rows)

If DECISIONS.md doesn't exist, create it with header row.

Decision aggregation failure does not block the slice.

### 7c: ROADMAP.md Update (C-45)

If ROADMAP.md exists:

1. Read ROADMAP.md
2. Find the slice row (search for slice ID)
3. Update: Status=complete, Tests={N}, Completed={today's date}, Key Decision={if any}
4. If slice row not found, log warning and skip (don't block)

If architecture changed (noted in summary):
- Update ROADMAP.md diagram if needed

If ROADMAP.md doesn't exist, skip with warning.

---

## Step 9: Visual Checkpoint (if applicable)

**Skip if `config.workflow.visual_checkpoint` is false.**

Only trigger if the slice has user-facing UI components (check contracts for UI/component boundaries).

Follow `references/checkpoint-protocol.md`:

```
VISUAL CHECK

What was built: {slice_name}

How to verify:
1. {command to start — e.g., npm run dev}
2. Navigate to {URL}
3. {what to look for — specific elements, behaviours}

Type "approved" or describe issues.
```

If issues → spawn debugger to investigate, then re-implement.

---

## Step 10: Complete

All tests green (functional + security), verification passed.

### Cleanup Checkpoint Tag

Remove the checkpoint tag now that the slice has completed successfully. The tag was used during implementation to allow rollback via `git checkout greenlight/checkpoint/{slice_id} -- .` if a CIRCUIT BREAK occurred.

```bash
# Best-effort cleanup — ignore error if tag was already removed
git tag -d greenlight/checkpoint/{slice_id} 2>/dev/null || true
```

### Update STATE.md

- Slice status → "complete"
- Test counts updated
- Security test counts updated
- Progress bar updated
- Step → next action

### Final Report

```
Slice {slice_id} complete: "{slice_name}"

Tests: {N} passing ({functional} functional, {security} security)
Files: {list of new/modified}
Contracts satisfied: {list}
Deviations: {count} ({types})
Verification: PASS

Progress: [{progress_bar}] {done}/{total} slices
```

### Suggest Next

Check GRAPH.json for what's now available:

```
Next actions:
- /gl:slice {next_id} — "{next_name}" (dependencies met)
- /gl:slice {parallel_id} — "{parallel_name}" (can run in parallel)
```

If all slices complete:
```
All slices complete! Run /gl:ship for final audit.
```

---

## Error Recovery

### Agent Spawn Failure

If any agent fails to spawn or returns an error:
1. Log the error
2. Retry once with same context
3. If retry fails → report to user and pause

### Context Overflow

If an agent reports context degradation (unlikely with proper slice sizing):
1. The agent should have stopped gracefully
2. Check what was completed
3. Split remaining work if needed
4. Spawn fresh agent for remaining work

### State Corruption

If STATE.md is inconsistent with actual test results:
1. Run full test suite
2. Rebuild state from test results + GRAPH.json
3. Log the discrepancy

Always trust test results over STATE.md when they disagree.
