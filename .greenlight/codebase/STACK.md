# Stack

Technical stack for the Greenlight project.

## Language

**Go 1.24**
- Evidence: `/go.mod` (line 3)
- Standard library only, zero external dependencies
- Embedding support via `embed.FS` for bundling markdown content

## Build Tools

**GoReleaser v2**
- Evidence: `/.goreleaser.yaml`, `/.github/workflows/release.yml` (line 45-48)
- Cross-platform builds: Linux & macOS, amd64 & arm64
- Binary name: `greenlight`
- Build flags:
  - `-s -w` (strip debug info)
  - Version injection via ldflags to `internal/version` package
- Output: tar.gz archives with checksums

**Lefthook**
- Evidence: `/.lefthook.yml`
- Git hooks manager
- Enforces conventional commits via `commit-msg` hook
- Regex validation: `^(feat|fix|docs|chore|refactor|test|ci)(\(.+\))?!?: .+`

## Release Automation

**semantic-release v4**
- Evidence: `/.releaserc.json`, `/.github/workflows/release.yml` (line 22-29)
- Automated versioning from conventional commits
- Branch: `main`
- Plugins:
  - `@semantic-release/commit-analyzer`
  - `@semantic-release/release-notes-generator`
  - `@semantic-release/changelog`
  - `@semantic-release/git`
  - `@semantic-release/github`
- Generates `CHANGELOG.md` and GitHub releases

**commitlint**
- Evidence: `/.commitlintrc.json`, `/.github/workflows/commitlint.yml`
- Config: `@commitlint/config-conventional`
- CI validation via `wagoid/commitlint-github-action@v6`

## Distribution

**npm (wrapper package)**
- Evidence: `/npm/package.json`, `/npm/bin/index.js`
- Package name: `greenlight-cc`
- Node requirement: `>=16.0.0`
- Purpose: Downloads and runs platform-specific Go binary from GitHub releases
- Runtime binary resolution: macOS/Linux, x64/arm64
- Version synced with Go binary via CI pipeline

**GitHub Releases**
- Evidence: `/.github/workflows/release.yml` (line 31-51)
- Platform: GitHub Actions (ubuntu-latest)
- Artifacts: 4 binaries (darwin_amd64, darwin_arm64, linux_amd64, linux_arm64)
- Release mode: `append` (preserves existing release assets)

## CI/CD

**GitHub Actions**
- Evidence: `/.github/workflows/`
- Workflows:
  1. `release.yml`: 3-stage pipeline (semantic-release → goreleaser → npm-publish)
  2. `commitlint.yml`: Validates commit messages on push/PR to main
- Trigger: Push to `main` branch
- Secrets required: `GITHUB_TOKEN`, `NPM_TOKEN`

## Embedded Content

**go:embed**
- Evidence: `/main.go` (line 11-12)
- Embedded paths:
  - `src/agents/*.md` (8 agent definitions)
  - `src/commands/gl/*.md` (11 command definitions)
  - `src/references/*.md` (3 reference docs)
  - `src/templates/*.md` (2 template files)
  - `src/CLAUDE.md` (main config)
- Total manifest: 25 files embedded at compile time
- Accessed via `fs.Sub(embeddedContent, "src")` in runtime

## File Structure

```
greenlight/
├── main.go                  # Entry point, embed declarations
├── internal/
│   ├── cli/                 # Command dispatcher
│   ├── cmd/                 # Subcommand handlers (install, uninstall, check, version)
│   ├── installer/           # File copy, conflict resolution, verification
│   └── version/             # Version strings (injected at build time)
├── src/                     # Content embedded into binary
│   ├── agents/              # Agent prompts
│   ├── commands/gl/         # Command documentation
│   ├── references/          # Protocol docs
│   ├── templates/           # Config templates
│   └── CLAUDE.md            # Global instructions
└── npm/                     # NPM wrapper package
    ├── package.json
    └── bin/index.js         # Binary downloader/runner
```

## No Test Framework

No test files or test framework detected in codebase.
- Evidence: No `*_test.go` files found, no test configuration
