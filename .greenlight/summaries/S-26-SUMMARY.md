# S-26: Documentation and Deprecation

## What Changed
Verification tiers are now documented across all Greenlight standards and infrastructure. CLAUDE.md has a 5-line hard rule. The checkpoint protocol recognizes Acceptance as a new checkpoint type (always pauses, even in yolo). Visual checkpoint is deprecated. The manifest includes verification-tiers.md.

## Contracts Satisfied
- C-72: CLAUDEmdVerificationTierRule — 5-line subsection in Code Quality Constraints after Circuit Breaker, before Logging & Observability. States default tier (verify), test writer rejection routing, and full protocol reference.
- C-73: CheckpointProtocolAcceptanceType — Acceptance checkpoint type added to checkpoint-protocol.md (always pauses, verify tier trigger, Step 6b reference). Visual marked deprecated. Cross-reference in verification-patterns.md. Deprecation note in config template. Step 9 is now a no-op with deprecation warning.
- C-74: ManifestVerificationTiersUpdate — references/verification-tiers.md added to manifest (34 → 35 entries). Alphabetical ordering maintained. CLAUDE.md remains last.

## Test Coverage
- 31 tests (11 for C-72, 12 for C-73, 8 for C-74)
- 0 security tests (documentation-only changes)

## Files Created/Modified
- `src/CLAUDE.md` — Verification Tiers subsection added
- `src/references/checkpoint-protocol.md` — Acceptance type, Visual deprecated
- `src/references/verification-patterns.md` — cross-reference added
- `src/templates/config.md` — visual_checkpoint deprecation note
- `src/commands/gl/slice.md` — Step 9 deprecation warning
- `internal/installer/installer.go` — manifest +1 entry
- Test fixtures updated across 6 test files (count 34 → 35)

## Architecture Impact
Infrastructure update. Manifest grows from 34 to 35 entries. No new boundaries.

## Verification
- Tier: auto (all contracts are infrastructure)
- All 31 tests passing, full suite green (650 total)
