# S-46: Interactive Commands

## What Changed
Users can now launch interactive Claude sessions from the terminal for project initialization (`gl init`) and system design (`gl design`).

## Contracts Satisfied
- C-115: RunInit — launches interactive Claude session with /gl:init skill
- C-116: RunDesign — launches interactive Claude session with /gl:design skill, requires existing .greenlight/ directory

## Key Behaviours
- Inside Claude context: prints skill instructions without spawning another Claude
- Shell context: prints launching message and spawns interactive session, blocking until completion
- RunInit does not require .greenlight/ directory (init creates it)
- RunDesign requires .greenlight/ directory and prints helpful error if missing
- Claude not in PATH: prints install instructions and returns exit code 1
- --dangerously-skip-permissions is always stripped from interactive sessions (enforced by process package)

## Tests
- 46 integration tests covering context detection, shell spawn, error handling, args behaviour, and security invariants

## Files Modified
- internal/cmd/init_cmd.go — replaced stub with full implementation
- internal/cmd/design.go — replaced stub with full implementation
