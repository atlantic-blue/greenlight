# S-43: tmux Manager

## What Changed
The CLI can now manage tmux sessions for parallel slice execution. IsAvailable checks whether tmux is in PATH. NewSession creates a detached tmux session with a named window and working directory. AddWindow adds a window to an existing session with a command. AttachSession attaches the current terminal to a session with stdin/stdout/stderr passthrough.

## Contracts Satisfied
- C-107: TmuxIsAvailable
- C-108: TmuxNewSession
- C-109: TmuxAddWindow
- C-110: TmuxAttachSession

## Test Coverage
- 18 integration tests using LookPath dependency injection for PATH verification
- IsAvailable: true when found, false when missing, bool-only return (3 tests)
- BuildNewSessionCmd: detached flag, session name, window name, working directory, command string, PATH error (6 tests)
- NewSession: PATH error propagation (1 test)
- BuildAddWindowCmd: target session, window name, command, PATH error (4 tests)
- BuildAttachCmd: target session, attach-session subcommand, PATH error (3 tests)
- Sentinel errors: non-nil and mutually distinct (1 test)

## Files
- internal/tmux/tmux.go (new — IsAvailable, NewSession, AddWindow, AttachSession, Build* variants)
- internal/tmux/tmux_test.go (new — 18 tests)

## Decisions
- LookPath as package-level variable for testability (same pattern as internal/process)
- Build* functions exposed separately for testing command construction without process execution
- All arguments passed as separate exec.Command args (no shell concatenation) to prevent injection

## Architecture Impact
New package: internal/tmux. Imported by internal/cmd for parallel slice execution via tmux sessions.
