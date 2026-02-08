# Architecture

## Pattern

**CLI tool with embedded content installer**

Greenlight is a single-binary CLI that installs TDD-first development framework files into Claude Code's configuration directories. It embeds markdown documentation at build time and copies it to either global (`~/.claude/`) or local (`./.claude/`) locations.

## Execution Flow

```
main.go
  ├─> Embed content via go:embed directive (src/**/*.md)
  ├─> Create fs.Sub from embedded content
  └─> cli.Run(args, contentFS)
       ├─> Parse subcommand from args[0]
       └─> Dispatch to appropriate cmd.Run* function
            ├─> install: Parse flags → Resolve target directory → installer.Install()
            ├─> uninstall: Parse flags → Resolve target directory → installer.Uninstall()
            ├─> check: Parse flags → Resolve target directory → installer.Check()
            └─> version: version.Version → Print to stdout
```

**Entry point**: `/Users/juliantellez/github.com/atlantic-blue/greenlight/main.go`

1. Embeds all markdown content from `src/` at compile time using `go:embed`
2. Creates a virtual filesystem from embedded content
3. Delegates to `cli.Run()` with remaining args and the content filesystem

**CLI dispatcher**: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/cli/cli.go`

1. Takes args and content filesystem
2. Switches on first argument to determine subcommand
3. Routes to appropriate command handler in `internal/cmd/`
4. Returns exit code (0 for success, 1 for failure)

**Command handlers**: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/cmd/*.go`

Each command:
1. Parses flags from args (scope, conflict strategy)
2. Resolves target directory based on scope
3. Calls installer functions
4. Returns exit code

## Key Design Patterns

### Command Pattern
Each subcommand has its own `Run*` function in `internal/cmd/`. Commands are stateless functions that take args, dependencies, and an output writer.

### Strategy Pattern
Conflict resolution for existing `CLAUDE.md` files uses three strategies:
- `keep`: Save greenlight version as `CLAUDE_GREENLIGHT.md`
- `replace`: Backup existing, overwrite with greenlight version
- `append`: Concatenate greenlight content to existing file

Implemented in `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/conflict.go`

### Dependency Injection
All functions take explicit dependencies:
- `io.Writer` for output (enables testing without stdout)
- `fs.FS` for embedded content (abstraction over file system)
- No global state

### Manifest-Based Installation
`installer.Manifest` is a hardcoded list of all files to install. The installer iterates this list rather than walking the embedded filesystem. This makes the installation predictable and auditable.

Location: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/installer.go`

## Error Handling Strategy

**Explicit error returns throughout**

Every function that can fail returns `error` as the last return value. Errors are wrapped with context using `fmt.Errorf()` with `%w` verb to preserve error chains.

Pattern:
```go
if err := doThing(); err != nil {
    return fmt.Errorf("context about what failed: %w", err)
}
```

**Early exit on error**

Commands return exit code 1 immediately on error. No error recovery or retry logic. Errors are written to stdout with `fmt.Fprintf()`.

**File operation errors**

File operations distinguish between "file not found" (expected in some flows) and other errors:
```go
if os.IsNotExist(err) {
    // Handle missing file case
}
if err != nil {
    // Handle other errors
}
```

Used in:
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/conflict.go` (handleConflict)
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/installer.go` (Check, Uninstall)

## Embedded Content (go:embed)

**Directive location**: `/Users/juliantellez/github.com/atlantic-blue/greenlight/main.go:11`

```go
//go:embed src/agents/*.md src/commands/gl/*.md src/references/*.md src/templates/*.md src/CLAUDE.md
var embeddedContent embed.FS
```

**What gets embedded**:
- All agent definitions (`src/agents/*.md`)
- All command definitions (`src/commands/gl/*.md`)
- Reference documentation (`src/references/*.md`)
- Template files (`src/templates/*.md`)
- Main CLAUDE.md config file

**How it's used**:
1. At build time, Go compiler includes file contents in binary
2. At runtime, `fs.Sub(embeddedContent, "src")` creates a filesystem rooted at `src/`
3. Installer reads files using `fs.ReadFile(contentFS, relPath)`
4. Files are written to target directory with `os.WriteFile()`

**Asymmetric CLAUDE.md placement**:
- Global install: `~/.claude/CLAUDE.md`
- Local install: `./CLAUDE.md` (project root, NOT inside `.claude/`)

This asymmetry is intentional. Local CLAUDE.md files sit at project root so Claude Code finds them before descending into `.claude/` subdirectory.

## Version Information

**Build-time injection**: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/version/version.go`

Three variables set via ldflags during build:
- `Version`: Semantic version (e.g., "1.0.0")
- `GitCommit`: Short commit hash
- `BuildDate`: ISO 8601 timestamp

Default to "dev", "unknown", "unknown" if not set.

**Version file tracking**: `.greenlight-version`

After installation, greenlight writes a version file containing these three values. Used by `check` command to display installed version.

## No External Dependencies

The codebase has zero external dependencies beyond the Go standard library. Only imports:
- `embed`
- `fmt`
- `io`
- `io/fs`
- `os`
- `path/filepath`
- `strings`

This keeps the binary small, compilation fast, and eliminates supply chain risk.

## Scope Resolution

**Global vs Local**: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/cmd/scope.go`

Two installation scopes:
- `--global`: Installs to `~/.claude/` (user home directory)
- `--local`: Installs to `./.claude/` (current working directory)

Scope is required for all commands except `version` and `help`. Commands error if both `--global` and `--local` are provided.

Directory resolution:
- Global: `os.UserHomeDir()` + `.claude`
- Local: `.claude` (relative to CWD)

## File Permissions

All files written with consistent permissions:
- Directories: `0o755` (rwxr-xr-x)
- Files: `0o644` (rw-r--r--)

No executable files are created. All content is markdown documentation.

## Uninstall Behavior

**Selective removal**: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/installer.go:133`

Uninstall removes only files in the manifest EXCEPT `CLAUDE.md`. This file may contain user content, so it's left in place.

After removing files, uninstall cleans up empty directories in deepest-first order:
1. `commands/gl/`
2. `commands/`
3. `agents/`
4. `references/`
5. `templates/`

If any directory contains files, it's left in place.

## Check Verification

**Integrity check**: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/installer.go:174`

The `check` command verifies:
1. All manifest files exist
2. All files are non-empty (size > 0)
3. Version file exists and is readable

Returns detailed report:
- `MISSING` for absent files
- `EMPTY` for zero-byte files
- `ERROR` for files that exist but can't be read

Exit code 0 if all checks pass, 1 otherwise.
