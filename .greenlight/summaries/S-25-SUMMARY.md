# S-25: Rejection Counter

## What Changed
The /gl:slice pipeline's Step 6b now tracks rejections per slice with a structured counter and log. When a slice is rejected 3 times, the orchestrator escalates with full rejection history and three options: re-scope the contract, pair with the user for detailed guidance, or skip verification with an acknowledged mismatch.

## Contracts Satisfied
- C-70: RejectionCounter — YAML state initialization (slice_id, rejection_count, rejection_log), log entry structure (attempt, feedback, classification, action_taken), counter persistence within execution, reset on new execution, contract revision does NOT reset, CounterOverflow/LogCorruption error handling, in-memory only
- C-71: RejectionEscalation — ESCALATION header with slice name, rejection history display, 3 options (re-scope resets + restarts, pair resets + collects guidance, skip marks auto-verified + logs acknowledgment), InvalidEscalationChoice/EmptyRejectionLog error handling, follows checkpoint patterns

## Test Coverage
- 53 tests (30 for C-70, 23 for C-71)
- 0 security tests (documentation-only changes)

## Files Created/Modified
- `src/commands/gl/slice.md` — added rejection counter state, log structure, escalation format within Step 6b
- `src/references/verification-tiers.md` — expanded section 5 with YAML state, log entries, escalation format, error handling

## Architecture Impact
Additive change to existing Step 6b and verification-tiers.md section 5. No new boundaries or integrations.

## Verification
- Tier: verify (human approval obtained)
- All 53 tests passing, full suite green
