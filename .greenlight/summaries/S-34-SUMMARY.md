# S-34: Manifest and Integration

## What Changed
Three new file paths added to the Go CLI manifest in installer.go. The manifest grows from 35 to 38 entries, enabling the CLI to install parallel state files alongside existing content. All test helpers (buildTestFS) and assertions across 6 test files updated for the new count.

## User Impact
- `/gl:migrate-state` command is now installed by the CLI
- State format reference (`references/state-format.md`) is installed for all agents
- Slice state template (`templates/slice-state.md`) is installed for new projects
- All three parallel state features are installable via `greenlight install`

## Contracts Satisfied
- C-90: ManifestParallelStateUpdate (auto)

## Tests
- 11 passing (0 security)
- Coverage: count, presence, ordering, regression, disk existence

## Files
- `internal/installer/installer.go` (modified — +3 manifest entries)
- `internal/installer/manifest_parallel_state_test.go` (new — 11 tests)
- `internal/installer/installer_test.go` (modified — buildTestFS + count updates)
- `internal/installer/check_test.go` (modified — summary string updates 35→38)
- `internal/installer/circuit_breaker_infra_test.go` (modified — count update)
- `internal/installer/doc_deprecation_test.go` (modified — count update)
- `internal/cmd/check_test.go` (modified — buildTestContentFS + summary updates)
- `internal/cmd/install_test.go` (modified — buildTestFS update)
- `internal/cli/cli_test.go` (modified — buildTestFS update)

## Architecture
No architecture changes. Manifest-only update enabling parallel state file installation.
