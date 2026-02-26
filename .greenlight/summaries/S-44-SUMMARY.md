# S-44: Parallel Slice Execution

## What Changed
Users can now run multiple ready slices concurrently via tmux sessions, with automatic fallback to sequential mode when tmux is unavailable or when the `--sequential` flag is set.

## Contracts Satisfied
- C-111: RunSliceParallel — parallel execution via tmux with configurable max windows
- C-112: RunSliceSequentialFallback — sequential fallback when tmux unavailable or forced

## Key Behaviours
- Auto-detect mode routes to parallel when 2+ slices are ready and tmux is available
- `--max N` limits the number of concurrent tmux windows (default: 4)
- `--sequential` flag forces sequential mode regardless of tmux availability
- Inside Claude context: always single slice, never parallel, with hint about remaining ready slices
- tmux session naming follows `{prefix}-{project}` pattern from config (default prefix: `gl`)
- Config-driven claude flags passed to each spawned window
- `--dry-run` prints the execution plan without spawning any processes
- Graceful fallback to sequential on tmux session creation failure

## Tests
- 38 integration tests covering parallel dry-run, sequential fallback, inside Claude invariants, max flag, single slice boundary, and error handling

## Files Modified
- internal/cmd/slice_cmd.go — extended with parallel/sequential routing, tmux session management, --max flag parsing

## Decisions
- Default max windows set to 4 to balance parallelism with system resources
- Session naming uses project directory basename for disambiguation
