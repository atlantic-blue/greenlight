# S-22: Schema Extension

## What Changed
The architect and verifier agents now understand verification tiers. Contracts can specify whether a slice needs human approval (`verify`) or can proceed automatically (`auto`) after tests pass.

## Contracts Satisfied
- C-62: ContractSchemaExtension — three new optional fields (Verification, Acceptance Criteria, Steps) added to the architect's contract format template
- C-63: VerifierTierAwareness — verifier reports effective tier, per-contract breakdown, and warnings in its verification output

## Test Coverage
- 24 tests (14 for C-62, 10 for C-63)
- 0 security tests (documentation-only changes)

## Files Modified
- `src/agents/gl-architect.md` — added verification tier fields to `<contract_format>` section
- `src/agents/gl-verifier.md` — added Verification Tier section to `<output_format>` report

## Architecture Impact
None — additive changes to existing agent definitions. No new boundaries or integrations.

## Verification
- Tier: auto (infrastructure/internal plumbing)
- All 24 tests passing, full suite green (480 tests)
