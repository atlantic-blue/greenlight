# S-24: Rejection Flow

## What Changed
The /gl:slice pipeline's Step 6b rejection flow now includes detailed gap classification UX, test writer spawn with rejection context, and contract revision routing. When a user rejects a verification checkpoint, they choose between tightening tests, revising the contract, or providing more detail — each route has full error handling, agent isolation rules, and re-verification cycles.

## Contracts Satisfied
- C-67: RejectionClassification — gap classification UX with verbatim feedback display, 3 options mapping to test_gap/contract_gap/implementation_gap, InvalidChoice/EmptyFeedback error handling, default to test_gap after 2 failed re-prompts
- C-68: RejectionToTestWriter — test writer spawn with rejection_context XML (feedback, classification, detailed_feedback), contract and acceptance_criteria blocks, additive tests, implementer receives test names only, full re-verification cycle, 4 error states handled
- C-69: RejectionToContractRevision — CONTRACT REVISION header with full contract display, minor vs fundamental revision paths, /gl:add-slice recommendation, restart from Step 1, rollback offer, EmptyRevision handling

## Test Coverage
- 48 tests (16 for C-67, 18 for C-68, 14 for C-69)
- 0 security tests (documentation-only changes)

## Files Created/Modified
- `src/commands/gl/slice.md` — expanded rejection flow within Step 6b (+101 lines)
- `src/references/verification-tiers.md` — added classification mapping table and implementer isolation constraint (+14 lines)

## Architecture Impact
Additive change to existing Step 6b rejection flow. No new boundaries or integrations. Expands the placeholder rejection options into full routing with error handling.

## Verification
- Tier: verify (human approval obtained)
- All 48 tests passing, full suite green
