# S-37: CLI Dispatch Extension

## What Changed
The CLI now recognizes six new commands: status, slice, init, design, roadmap, and changelog. Each dispatches to a stub handler that prints "not implemented yet" until real slices provide the implementation. The help output groups commands into four categories: Project lifecycle, Building, State & progress, and Admin.

## Contracts Satisfied
- C-97: CLIDispatchExtension

## Test Coverage
- 12 new integration tests covering dispatch routing, exit codes, argument forwarding, and categorized usage output
- All existing commands continue to work unchanged

## Files
- internal/cli/cli.go (modified — 6 new switch cases, updated printUsage)
- internal/cli/cli_test.go (modified — 12 new tests)
- internal/cmd/status.go (new — stub handler)
- internal/cmd/slice_cmd.go (new — stub handler)
- internal/cmd/init_cmd.go (new — stub handler)
- internal/cmd/design.go (new — stub handler)
- internal/cmd/roadmap.go (new — stub handler)
- internal/cmd/changelog.go (new — stub handler)

## Decisions
- New commands don't accept contentFS since they operate on .greenlight/ state, not embedded content
- Stub handlers return exit code 0 with "not implemented yet" message
- Usage organized into 4 categories matching the design spec
