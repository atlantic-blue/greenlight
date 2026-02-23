# S-41: Process Spawner

## What Changed
The CLI can now spawn Claude processes for both headless and interactive modes. SpawnClaude starts a non-blocking headless process with a prompt flag, while SpawnInteractive launches a blocking terminal session. Both verify that the claude binary is in PATH before attempting to start. Interactive mode enforces a security invariant: --dangerously-skip-permissions is always stripped from flags.

## Contracts Satisfied
- C-103: ProcessSpawnClaude
- C-104: ProcessSpawnInteractive

## Test Coverage
- 15 integration tests using LookPath dependency injection for PATH verification
- BuildClaudeCmd: prompt flag, additional flags, working directory, stdout/stderr wiring, empty prompt error, PATH error
- BuildInteractiveCmd: optional prompt, dangerous flag stripping (4 sub-cases), working directory, PATH error
- Sentinel errors: non-nil and mutually distinct
- SpawnClaude: PATH error propagation

## Files
- internal/process/process.go (new — SpawnClaude, SpawnInteractive, BuildClaudeCmd, BuildInteractiveCmd)
- internal/process/process_test.go (new — 15 tests)

## Decisions
- LookPath as package-level variable for testability (no interface needed for single function)
- BuildClaudeCmd/BuildInteractiveCmd exposed as separate functions for testing command construction without process execution
- Interactive mode always strips --dangerously-skip-permissions even if caller explicitly passes it

## Architecture Impact
New package: internal/process. Imported by internal/cmd for slice, init, and design commands.
