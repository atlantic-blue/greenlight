# S-29: Init Command and State Detection

## What Changed
The `/gl:init` command now creates file-per-slice state alongside the legacy STATE.md. New projects get individual slice files in `.greenlight/slices/`, a `project-state.json` for non-slice metadata, and a regenerated STATE.md marked as generated. A state detection flow tells all commands which format is active.

## User Impact
- New projects initialised with `/gl:init` use the parallel-safe file-per-slice format from day one
- Each slice gets its own state file — no more write conflicts between concurrent sessions
- Existing projects with legacy STATE.md continue to work unchanged
- All commands follow a consistent detection flow: check `slices/` first, then `STATE.md`, then suggest `/gl:init`

## Contracts Satisfied
- C-78: InitSliceDirectory (verify)
- C-79: InitProjectState (auto)
- C-80: InitStateDetection (verify)

## Tests
- 43 passing (0 security)
- Coverage: all 3 contracts fully covered

## Files
- `src/commands/gl/init.md` (modified — +183 lines)
- `internal/installer/init_state_detection_test.go` (new — 43 tests)

## Architecture
No architecture changes. Extends the existing init command with new output phases.
