# S-38: Status Command

## What Changed
Users can now run `greenlight status` from the terminal to see project progress without needing Claude. The command reads all slice files and the dependency graph, then displays a formatted report with progress bar, running/ready/blocked slices, and test counts. A `--compact` flag outputs a single line suitable for tmux status bars.

## Contracts Satisfied
- C-98: RunStatus
- C-99: RunStatusCompact

## Test Coverage
- 24 integration tests using temp directories and t.Chdir for filesystem isolation
- Default mode: progress bar, running with step, ready/blocked with deps, test sums, error handling
- Compact mode: single-line output, tmux resilience (placeholder on error), format validation

## Files
- internal/cmd/status.go (modified — replaced stub with full implementation)
- internal/cmd/status_test.go (new — 24 tests)

## Decisions
- 20-character ASCII progress bar using `#` (filled) and `.` (empty)
- Graceful degradation when GRAPH.json missing: shows progress and tests, warns about dependency info
- Compact mode always returns exit 0, even on error, for tmux status bar resilience
- Blocked slices show specific unmet dependency IDs in parentheses
