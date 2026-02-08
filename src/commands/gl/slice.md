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

From CONTRACTS.md:
- Full contract definitions for this slice's contracts
- Input/output types, error states, invariants, security requirements

From config.json:
- Model assignments for each agent
- Workflow toggles (security_scan, visual_checkpoint, etc.)
- Test commands

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

## Step 1: Write Tests (Agent A — Fresh Context)

The test writer NEVER sees implementation code. Send ONLY contracts, stack info, and existing test patterns.

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

Max 3 implementation attempts. After 3 failures → pause and ask user.

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
| Infrastructure error | Fix infrastructure, re-run |

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

## Step 7: Visual Checkpoint (if applicable)

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

## Step 8: Complete

All tests green (functional + security), verification passed.

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
