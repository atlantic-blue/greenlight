# S-27: Architect Integration

## What Changed
The architect agent now has explicit guidance for setting verification tiers on contracts. A new "Verification Tier Selection" section in gl-architect.md tells the architect when to use `auto` vs `verify`, how to write acceptance criteria (behavioral, present tense, specific), and how to write verification steps (action verbs, commands/URLs). The output checklist enforces tier presence on every contract.

## Contracts Satisfied
- C-75: ArchitectTierGuidance — Verification Tier Selection section added to gl-architect.md with when-to-use-auto, when-to-use-verify, writing acceptance criteria, writing steps subsections. Output checklist extended with 4 new items. Error states documented (MissingTierOnContract defaults to verify, TooManyCriteria suggests splitting). Non-prescriptive: architect can override with good reasoning.

## Test Coverage
- 60 tests (all for C-75)
- 0 security tests (agent guidance update)

## Files Created/Modified
- `src/agents/gl-architect.md` — Verification Tier Selection section + output checklist additions

## Architecture Impact
No new boundaries. Agent guidance update only.

## Verification
- Tier: auto (agent guidance, no user-visible behaviour)
- All 60 tests passing, full suite green (710 total)
