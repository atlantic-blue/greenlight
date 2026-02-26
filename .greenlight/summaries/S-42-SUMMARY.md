# S-42: Single Slice Command

## What Changed
Users can now run slices from the terminal. `gl slice S-35` spawns a headless Claude session for that slice. `gl slice` with no ID auto-detects the first ready slice by wave/ID order. `gl slice --dry-run` previews the command without spawning. Inside Claude, the command outputs slice info instead of spawning another process.

## Contracts Satisfied
- C-105: RunSliceSingle
- C-106: RunSliceAutoDetect

## Test Coverage
- 36 tests covering both named-slice and auto-detect modes
- No greenlight directory: exit 1 with error (2 tests)
- Invalid slice ID: exit 1 with error mentioning ID and suggesting gl status (3 tests)
- Dry-run: prints command without spawning, shows config flags (6 tests)
- Inside Claude: outputs slice info, never spawns another Claude (4 tests)
- Shell context spawn error: non-zero exit, error message (3 tests)
- Config resilience: graceful degradation without config.json (2 tests)
- Auto-detect zero ready: returns 0, prints blocked status (3 tests)
- Auto-detect one ready: picks and runs it (3 tests)
- Auto-detect multiple ready: wave order, hints about remaining, sequential flag (5 tests)
- Auto-detect state errors: exit 1 on missing state (2 tests)
- Output invariants: writer contract, no tmux mentions (3 tests)

## Files
- internal/cmd/slice_cmd.go (modified — full implementation replacing stub)
- internal/cmd/slice_cmd_test.go (new — 36 tests)

## Decisions
- Config flags read from .greenlight/config.json parallel.claude_flags, graceful degradation on missing config
- Auto-detect with 2+ ready slices picks first by wave/ID order (parallel deferred to S-44)
- Inside Claude: outputs info for skill consumption, never spawns nested process
- Slice ID validated against GRAPH.json before any execution

## Architecture Impact
internal/cmd now imports internal/process and internal/state. RunSlice is the first command that bridges state reading with process spawning, establishing the pattern for S-44 (parallel) and S-45 (watch).
