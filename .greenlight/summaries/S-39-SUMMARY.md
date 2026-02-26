# S-39: Help Command

## What Changed
Users can now run `greenlight help` from the terminal to see all available commands grouped by category. Inside a greenlight project, the output appends a state summary showing slice count, complete count, and ready count. Outside a project, it suggests running `gl init`.

## Contracts Satisfied
- C-100: RunHelp

## Test Coverage
- 13 integration tests using temp directories and t.Chdir for filesystem isolation
- Command listing: all 4 categories present, all individual commands listed
- Exit code: always returns 0 regardless of project presence or state errors
- Project detection: state summary with counts, init suggestion when no project
- Error resilience: corrupt slice files and empty slices dir don't prevent command listing
- Writer contract: all output directed to provided io.Writer

## Files
- internal/cmd/help.go (new — RunHelp with grouped command listing and project state summary)
- internal/cmd/help_test.go (new — 13 tests)
- internal/cli/cli.go (modified — help dispatch updated to call cmd.RunHelp)

## Decisions
- State summary is best-effort: read errors silently return zero counts
- Context detection uses .greenlight/ directory existence, not $CLAUDE_CODE env var
- Ready count requires both slices and GRAPH.json; if graph missing, ready defaults to 0
