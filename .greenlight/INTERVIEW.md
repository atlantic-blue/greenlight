# Project Interview

## Value Proposition
TDD-first development system for Claude Code — a Go CLI that installs isolated agents, slash commands, and engineering standards that enforce test-driven development with agent isolation.

## Users
Developers using Claude Code who want contract-driven, TDD-enforced development with security built into every slice.

## MVP Scope
Stabilise the existing CLI with tests and contracts. No new features.

1. Install globally — copy agents, commands, templates, references, CLAUDE.md to ~/.claude/
2. Install locally — copy to .claude/ with CLAUDE.md conflict handling (keep/append/replace)
3. Check installation — verify installed files exist and match embedded content
4. Uninstall — remove installed files cleanly
5. Show version — print version info

## Stack
Go 1.24 (stdlib only, no external dependencies)

## Constraints
- Existing CI: GoReleaser + semantic-release + lefthook + commitlint
- npm wrapper package exists for `npx greenlight-cc install`
- Content embedded via `go:embed` from src/ directory
- Conventional commits enforced
- Binary already built and released

## Deferred Ideas
- `greenlight upgrade` — in-place upgrade preserving user config
- Homebrew tap
- Curl install script
- `greenlight doctor` — diagnose Claude Code config issues
- Plugin system for custom agents and commands
- Local-only usage analytics (opt-in)
