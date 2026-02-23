# S-35: Frontmatter Parser

**Milestone:** cli-orchestrator
**Completed:** 2026-02-23
**Tests:** 22 (0 security)

## What Was Built

A Go package (`internal/frontmatter`) that parses and writes flat key-value YAML frontmatter from `.greenlight/slices/*.md` files. This is the foundation for all CLI state reading.

## Contracts Satisfied

- **C-91 FrontmatterParse:** Parses content with `---` delimiters into a `map[string]string` of fields and a body string. Three error sentinels for missing, unclosed, and invalid frontmatter.
- **C-92 FrontmatterWrite:** Serializes fields and body back into frontmatter format with sorted keys for deterministic output. Roundtrip-safe with Parse.

## Key Decisions

- Line-by-line string splitting (D-43) — no YAML library needed
- Split on first colon only — preserves URLs and timestamps with colons in values
- Whitespace-only lines between delimiters are skipped silently
- Keys written in alphabetical order for deterministic output

## Files

- `internal/frontmatter/frontmatter.go` (new — 119 lines)
- `internal/frontmatter/frontmatter_test.go` (new — 22 tests)
