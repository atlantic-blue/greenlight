# S-45: Watch Mode

## What Changed
Users can now fire-and-forget to drain the dependency graph. Watch mode auto-launches new slices as slots free up, and the enhanced dry-run shows a categorized view of project state.

## Contracts Satisfied
- C-113: RunSliceWatch — poll loop with auto-refill, immediate termination when no work remains
- C-114: RunSliceDryRun — enhanced dry-run with Ready/Running/Blocked/Would launch categories

## Key Behaviours
- `gl slice --watch` enters a poll loop that detects completions and launches new slices
- Watch mode terminates when no running and no ready slices remain, printing a summary with done count and test totals
- `gl slice --dry-run` now shows categorized output: Ready, Running (with step), Blocked (with unmet deps), Would launch
- "Would launch" respects --max cap
- Watch interval configurable via `parallel.watch_interval_seconds` in config (default 30s)
- Inside Claude context: no poll loop, delegates to single-slice behaviour
- Immediate termination path for all-complete and all-blocked scenarios

## Tests
- 50 integration tests covering enhanced dry-run categories, watch termination, flag parsing, inside Claude, config, and error handling

## Files Modified
- internal/cmd/slice_cmd.go — added watch flag, enhanced dry-run, watch loop, termination summary, blocked/running helpers
