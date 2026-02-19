# Circuit Breaker Protocol

Shared protocol for detecting and halting unproductive implementation work.

All implementing agents (gl-implementer, gl-debugger) follow these rules. This protocol is **additive** to `references/deviation-rules.md` — it does not replace it. **Rule 4 (ARCH-STOP) from the deviation rules takes priority** over circuit breaker handling. The circuit breaker handles unproductive work (spinning on the same failure); deviation Rule 4 handles unplanned architectural scope.

---

## Section 1: Attempt Tracking State

Each implementing agent maintains a structured state object per test, updated after every implementation attempt:

```yaml
slice_id: S-{N}                    # Slice identifier from orchestrator
test_name: "{test function name}"  # Exact test name being attempted
attempt_count: 0                   # Number of attempts on this specific test (resets on success)
files_touched: []                  # List of files modified during attempts on this test
description: ""                    # Human-readable summary of what was attempted
last_error: ""                     # Exact error output from last test run
checkpoint_tag: ""                 # Git tag or commit SHA at start of this test's attempts
total_failed_attempts: 0           # Slice-level accumulator — counts all failures across all tests
```

### Field Definitions

- **slice_id**: Identifies which slice is being implemented. Used in diagnostic reports.
- **test_name**: The exact test name as reported by the test runner. Used to scope attempt counting.
- **attempt_count**: Increments once per complete attempt cycle on a specific test. Resets to 0 on test success. Does NOT increment for infrastructure errors.
- **files_touched**: Accumulates every file path modified during attempts on this test. Used to assess scope drift.
- **description**: One-sentence summary of what strategy was tried in the last attempt.
- **last_error**: Copy-paste of the actual test runner output for the failing assertion. Never paraphrased.
- **checkpoint_tag**: The git tag or commit hash captured before any modifications for this test's attempts. Used for rollback.
- **total_failed_attempts**: Slice-level accumulator. Increments each time any test attempt fails. Never resets within a slice execution. Used to enforce the slice-level ceiling.

---

## Section 2: Per-Test Trip Threshold

**Threshold: 3 failures per test.**

This is mandatory and not configurable. After 3 failed attempt cycles on the same test, the circuit breaker trips for that test.

### What Counts as an Attempt

An attempt is one complete cycle:
1. Implement or modify code
2. Run the test command
3. Observe the result

If the test fails, that is one failed attempt. Increment `attempt_count` by 1 and `total_failed_attempts` by 1.

### Infrastructure Errors Are Not Attempts

An infrastructure error does NOT count as a failed attempt. Infrastructure errors include:
- Build failures caused by syntax errors in unrelated files
- Missing environment variables needed for test runner setup
- Network timeouts when fetching dependencies
- Filesystem permission errors
- Test runner crashes unrelated to assertion failures

If the test runner fails to execute the test (as opposed to executing and asserting failure), do not increment either counter. Fix the infrastructure issue under deviation Rule 3 (auto-fix blocking issues), then re-run. Only increment counters when a test actually runs and fails its assertions.

### On Threshold Trip

When `attempt_count` reaches 3 for a specific test:
1. Stop attempting that test
2. Produce a Structured Diagnostic Report (see Section 4)
3. STOP and report to the orchestrator

---

## Section 3: Slice-Level Ceiling

**Ceiling: 7 total failures across all tests.**

When `total_failed_attempts` reaches 7, the circuit breaker trips at the slice level regardless of which tests are failing or how many tests have been attempted.

This prevents a scenario where an implementer makes slow progress (e.g., 2 failed attempts per test across 4 tests) without ever hitting the per-test threshold of 3.

### On Ceiling Trip

When `total_failed_attempts` reaches 7:
1. Stop all implementation work immediately
2. Produce a Structured Diagnostic Report for each currently-failing test
3. STOP and report to the orchestrator with all diagnostic reports

---

## Section 4: Structured Diagnostic Report Format

When a circuit trips (per-test at 3, or slice-level at 7), produce a diagnostic report with these 8 required fields:

```yaml
test_expectation: |
  What the test expects — describe the expected behaviour, not the assertion syntax.
  Example: "endpoint returns 409 when email already exists in database"

actual_error: |
  Exact copy of the test runner output. Do not paraphrase. Include the assertion
  failure message, stack trace if present, and exit code.

attempt_log:
  - attempt: 1
    strategy: "Added duplicate-email check before insert"
    files: ["src/users/service.ts"]
    result: "Still 500 — constraint violation thrown before check"
  - attempt: 2
    strategy: "Moved check to repository layer"
    files: ["src/users/repository.ts"]
    result: "409 returned but wrong error message format"
  - attempt: 3
    strategy: "Unified error message format via error factory"
    files: ["src/users/errors.ts", "src/users/repository.ts"]
    result: "Test expects 'email_already_exists' key, implementation returns 'duplicate_email'"

cumulative_files_modified:
  - src/users/service.ts
  - src/users/repository.ts
  - src/users/errors.ts

scope_violations:
  - file: src/shared/errors.ts
    justification: "Modified shared error factory to support new code — outside slice boundary"

best_hypothesis: |
  Best current theory on root cause, based on all evidence gathered.
  Example: "Contract C-12 specifies error key 'email_already_exists' but the test
  fixture was built against a draft that used 'duplicate_email'. The contract is
  ambiguous about which value is canonical."

specific_question: |
  The single most useful question that, if answered, would unblock implementation.
  Example: "Does contract C-12 intend 'email_already_exists' or 'duplicate_email'
  as the error key? The test and the prose description disagree."

recovery_options:
  - option: "Clarify contract C-12 and rerun with fresh implementer"
    risk: "Low — just needs contract update"
  - option: "Accept 'duplicate_email' as canonical and update the test"
    risk: "Medium — test writer decision required"
  - option: "Use counter reset protocol and retry with different approach"
    risk: "Low if hypothesis about key name is correct"
```

All 8 fields are required. Do not omit fields. If a field is not applicable, state why explicitly.

---

## Section 5: Scope Lock Protocol

Before modifying any file, the implementing agent must verify that the file is within scope for the current slice.

### Scope Inference Rules

Scope is determined in priority order:

1. **Explicit override (highest priority):** If the orchestrator provides a `files_in_scope` list, only those files may be modified without justification. All other files require the justification format below.

2. **Inferred from contracts:** Files referenced by path or package name in the slice's contracts are in scope. If a contract says "create `src/users/service.ts`", that file is in scope.

3. **Inferred from GRAPH.json packages and deliverables:** The slice's entry in `GRAPH.json` lists `packages` and `deliverables`. All files within those packages or matching those deliverable paths are in scope. Files outside these boundaries require justification.

### Scope Check Before File Modification

Before modifying any file, verify:
1. Is this file listed in `files_in_scope`? → Proceed.
2. Is this file referenced in the slice's contracts? → Proceed.
3. Is this file within the slice's `packages` or `deliverables` in GRAPH.json? → Proceed.
4. None of the above? → You must provide justification before proceeding.

### Justification Format

When modifying a file outside inferred scope, record the justification in the attempt tracking state:

```
File: src/shared/errors.ts
Failing test: TestUserRegistration_ReturnsConflictOnDuplicateEmail
Reason: Error factory function needed to standardise error key format required by the contract
Relationship: The contract (C-12) references the error shape; the factory is the only way to produce that shape without duplicating logic
```

The justification must include:
- **File**: Absolute or repo-relative path of the out-of-scope file
- **Failing test**: Exact name of the test that requires this modification
- **Reason**: Why modifying this file is necessary to make the test pass
- **Relationship**: How this file relates to the slice's contract boundary

An unjustifiable out-of-scope modification (one where no honest justification can be written) must not be made. If you cannot write a justification, the modification is out of scope. Apply deviation Rule 4 (ARCH-STOP) and report to the orchestrator.

---

## Section 6: Counter Reset Protocol

The counter reset protocol is used when the orchestrator decides to retry a slice after a circuit trip, rather than escalating.

### Reset Procedure

1. **Reset `attempt_count` to 0** for all tests. The per-test counter starts fresh.
2. **Reset `total_failed_attempts` to 0**. The slice ceiling starts fresh.
3. **Rollback to checkpoint**: Use the `checkpoint_tag` to restore the codebase to the state before the failed attempts. This ensures the fresh implementer starts from a known-good baseline, not a partially-modified state.
4. **Spawn a fresh implementer** with the same slice inputs, plus a `<prior_attempts>` block containing:
   - The diagnostic reports from the tripped circuit
   - A summary of what was tried ("what was tried") in each diagnostic report
   - The specific questions raised in each diagnostic report
   - Any guidance from the orchestrator or user following the circuit trip

### What the Fresh Implementer Receives

The fresh implementer sees the same inputs as the original, plus:

```xml
<prior_attempts>
Circuit tripped after {N} total failed attempts.

For test "{test_name}":
- what was tried: {attempt_log summary}
- hypothesis: {best_hypothesis}
- guidance: {orchestrator or user guidance}
</prior_attempts>
```

### What the Fresh Implementer Must Not Do

- Do not repeat the same strategies documented in the `attempt_log`
- Do not assume prior modifications are still in the codebase (rollback happened)
- Do not increment attempt counters for infrastructure fixes performed before first test run

---

## Section 7: Relationship to Deviation Rules

This protocol is **additive** to `references/deviation-rules.md`. The deviation rules and the circuit breaker serve different purposes:

| Protocol | Handles |
|----------|---------|
| Deviation rules | Unplanned work discovered during execution |
| Circuit breaker | Unproductive repetition on planned work |

**Rule 4 (ARCH-STOP) from the deviation rules takes priority over the circuit breaker.** If a test cannot pass because it requires an architectural change (new table, new service, breaking API change), apply Rule 4 immediately — do not exhaust 3 attempts first. The circuit breaker is for cases where the work is within scope but the implementation approach is not converging.

### Combined Priority Order

```
Deviation Rule 4 (ARCH-STOP)
  > Circuit breaker per-test threshold (3 failures)
  > Circuit breaker slice ceiling (7 total_failed_attempts)
  > Deviation Rules 1-3 (auto-fix)
```

When multiple rules could apply:
- If the fix requires architectural change → Rule 4 immediately
- If a test has failed 3 times → circuit breaker, STOP
- If `total_failed_attempts` reaches 7 → circuit breaker, STOP
- If the issue is a bug, missing critical feature, or blocker → fix automatically under Rules 1-3
