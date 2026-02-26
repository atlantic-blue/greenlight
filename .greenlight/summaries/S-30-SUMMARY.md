# S-30: Slice Command State Write

## What Changed
The `/gl:slice` command now reads from and writes to individual slice files instead of the monolithic STATE.md. Each session writes only to its own slice's file, eliminating the concurrent write conflicts that caused state corruption. STATE.md is regenerated as a summary view after every write. Advisory session tracking warns when two sessions attempt to claim the same slice.

## User Impact
- Running `/gl:slice` on different slices in parallel no longer risks state corruption
- Each slice session owns its file — writes never touch other slices' state
- STATE.md stays up to date as a convenience view, regenerated automatically
- Session tracking warns (but never blocks) when a slice appears to be claimed by another session
- Legacy projects continue to work exactly as before — no forced migration

## Contracts Satisfied
- C-81: SliceCommandStateWrite (verify)
- C-82: SliceSessionTracking (auto)

## Tests
- 42 passing (0 security)
- Coverage: all contract elements fully covered

## Files
- `src/commands/gl/slice.md` (modified — +91 lines)
- `internal/installer/slice_command_state_test.go` (new — 42 tests)

## Architecture
No architecture changes. Extends the existing /gl:slice pipeline with format-aware read/write paths.
