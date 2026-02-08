# Integrations

External services and dependencies for the Greenlight project.

## APIs Consumed

**GitHub Releases API**
- Evidence: `/npm/bin/index.js` (line 38)
- Purpose: NPM wrapper downloads platform-specific binaries
- URL pattern: `https://github.com/atlantic-blue/greenlight/releases/download/v{VERSION}/greenlight_{VERSION}_{platform}_{arch}.tar.gz`
- Used at runtime when executing `greenlight-cc` npm package
- No authentication required (public releases)

**GitHub API (via GoReleaser)**
- Evidence: `/.goreleaser.yaml` (line 33-35)
- Purpose: Publishes release artifacts and checksums
- Owner: `atlantic-blue`
- Repo: `greenlight`
- Authentication: `GITHUB_TOKEN` (provided by Actions)
- Triggered during CI release workflow

**GitHub API (via semantic-release)**
- Evidence: `/.releaserc.json` (line 14)
- Purpose: Creates GitHub releases with generated notes
- Plugin: `@semantic-release/github`
- Authentication: `GITHUB_TOKEN` (provided by Actions)
- Triggered on push to `main` branch

## Third-Party SDKs

**None in Go binary**
- Evidence: `/go.mod` (no dependencies listed)
- Pure Go standard library implementation

**Node.js Built-in Modules (npm wrapper only)**
- Evidence: `/npm/bin/index.js`
- Modules used:
  - `os`: Platform/architecture detection
  - `fs`: File system operations
  - `path`: Path manipulation
  - `https`: Binary download
  - `zlib`: Gunzip decompression
  - `child_process`: Spawn Go binary
- No external npm dependencies

## Databases

**None**
- No database connections or ORMs detected
- Greenlight operates on filesystem only

## File Storage

**Local Filesystem**
- Evidence: `/internal/installer/installer.go`
- Operations:
  - Write: Copies embedded content to `~/.claude/` (global) or `./.claude/` (local)
  - Read: Verifies installation via `Check()` function
  - Delete: Removes managed files via `Uninstall()`
- File locations:
  - Global: `~/.claude/` + manifest files
  - Local: `./.claude/` + manifest files
  - Special case: `CLAUDE.md` goes to `~/.claude/CLAUDE.md` (global) or `./CLAUDE.md` (local project root)
- Manifest: 25 files tracked (agents, commands, references, templates, CLAUDE.md)
- Version tracking: `.greenlight-version` file contains version/commit/builddate

**Temporary Storage**
- Evidence: `/npm/bin/index.js` (line 114)
- Purpose: NPM wrapper downloads binary to temp dir, executes it, then cleans up
- Path: `os.tmpdir()/greenlight-*`
- Lifecycle: Created on run, deleted on exit

## CI/CD Platforms

**GitHub Actions**
- Evidence: `/.github/workflows/release.yml`, `/.github/workflows/commitlint.yml`
- Runners: `ubuntu-latest`
- Actions used:
  - `actions/checkout@v4`
  - `actions/setup-go@v5`
  - `actions/setup-node@v4`
  - `cycjimmy/semantic-release-action@v4`
  - `goreleaser/goreleaser-action@v6`
  - `wagoid/commitlint-github-action@v6`
- Secrets:
  - `GITHUB_TOKEN`: GitHub API access (auto-provided)
  - `NPM_TOKEN`: NPM publishing (user-provided)

**NPM Registry**
- Evidence: `/npm/package.json`, `/.github/workflows/release.yml` (line 69-73)
- Package: `greenlight-cc`
- Registry: `https://registry.npmjs.org`
- Authentication: `NODE_AUTH_TOKEN` from `NPM_TOKEN` secret
- Publish triggered after successful GoReleaser build

## External Content Sources

**None at Runtime**
- All content embedded at compile time via `go:embed`
- No runtime fetching of prompts, templates, or configuration
- Self-contained binary with zero external dependencies
