# S-32: Migration Command

## What Changed
A new `/gl:migrate-state` command converts legacy STATE.md-based projects to the file-per-slice format. Migration is explicit, one-way, and all-or-nothing. The original STATE.md is always preserved as a backup. After migration, all commands automatically use file-per-slice format via state detection.

## User Impact
- Existing projects can safely migrate to file-per-slice state with a single command
- Original STATE.md preserved as STATE.md.backup (no data loss)
- Migration is atomic: either fully complete or fully rolled back
- After migration, concurrent sessions can work on different slices without state corruption
- No auto-migration — users choose when to migrate (D-32)

## Contracts Satisfied
- C-85: MigrateStateCommand (verify)
- C-86: MigrateStateBackup (auto)

## Tests
- 55 passing (0 security)
- Coverage: all contract elements fully covered

## Files
- `src/commands/gl/migrate-state.md` (new — ~287 lines)
- `internal/installer/migrate_state_command_test.go` (new — 55 tests)

## Architecture
No architecture changes. Adds a new command that bridges the legacy and file-per-slice state formats.
