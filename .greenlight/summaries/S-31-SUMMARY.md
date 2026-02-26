# S-31: Supporting Command Updates

## What Changed
Six supporting commands now detect the state format and use file-per-slice reads/writes when available. `/gl:status` computes a summary from individual slice files. `/gl:pause` writes to the slice file and project-state.json. `/gl:resume` reads from slice files and project-state.json. `/gl:ship` reads all slice files to check completeness. `/gl:add-slice` creates new slice files with ID validation. `/gl:quick` updates test counts in the relevant slice file only. All commands regenerate STATE.md after writes and fall back to legacy format seamlessly.

## User Impact
- All 6 commands work correctly in both file-per-slice and legacy formats
- `/gl:status` shows real-time state computed from individual slice files (never cached)
- Corrupt slice files are skipped gracefully — other slices still display
- Legacy projects continue to work exactly as before — no forced migration
- Crash safety via write-to-temp-then-rename on all write operations

## Contracts Satisfied
- C-83: StatusSliceAggregation (verify)
- C-84: SupportingCommandStateAdaptation (verify)

## Tests
- 74 passing (0 security)
- Coverage: all contract elements fully covered across 6 command files

## Files
- `src/commands/gl/status.md` (modified — +49 lines)
- `src/commands/gl/pause.md` (modified — +49 lines)
- `src/commands/gl/resume.md` (modified — +39 lines)
- `src/commands/gl/ship.md` (modified — +39 lines)
- `src/commands/gl/add-slice.md` (modified — +49 lines)
- `src/commands/gl/quick.md` (modified — +49 lines)
- `internal/installer/supporting_command_state_test.go` (new — 74 tests)

## Architecture
No architecture changes. Extends 6 existing commands with format-aware state read/write paths.
