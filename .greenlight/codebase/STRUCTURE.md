# Structure

## Directory Tree

```
/Users/juliantellez/github.com/atlantic-blue/greenlight/
├── main.go                    # Entry point, go:embed directive
├── go.mod                     # Go module definition (no external deps)
├── internal/                  # Non-importable implementation packages
│   ├── cli/
│   │   └── cli.go            # Command dispatcher, usage printer
│   ├── cmd/
│   │   ├── check.go          # "check" subcommand
│   │   ├── install.go        # "install" subcommand
│   │   ├── scope.go          # Scope/flag parsing utilities
│   │   ├── uninstall.go      # "uninstall" subcommand
│   │   └── version.go        # "version" subcommand
│   ├── installer/
│   │   ├── conflict.go       # Conflict resolution strategies
│   │   └── installer.go      # Core install/uninstall/check logic
│   └── version/
│       └── version.go        # Build-time version variables
├── src/                       # Embedded content (not in binary source)
│   ├── agents/                # Agent definition files
│   │   ├── gl-architect.md
│   │   ├── gl-codebase-mapper.md
│   │   ├── gl-debugger.md
│   │   ├── gl-designer.md
│   │   ├── gl-implementer.md
│   │   ├── gl-security.md
│   │   ├── gl-test-writer.md
│   │   └── gl-verifier.md
│   ├── commands/gl/           # Command definition files
│   │   ├── add-slice.md
│   │   ├── design.md
│   │   ├── help.md
│   │   ├── init.md
│   │   ├── map.md
│   │   ├── pause.md
│   │   ├── quick.md
│   │   ├── resume.md
│   │   ├── settings.md
│   │   ├── ship.md
│   │   ├── slice.md
│   │   └── status.md
│   ├── references/            # Reference documentation
│   │   ├── checkpoint-protocol.md
│   │   ├── deviation-rules.md
│   │   └── verification-patterns.md
│   ├── templates/             # Template files
│   │   ├── config.md
│   │   └── state.md
│   ├── CLAUDE.md             # Main framework instructions
│   └── README.md             # Content documentation
├── .greenlight/               # Target directory for local installs
│   └── codebase/             # Codebase documentation
│       ├── ARCHITECTURE.md   # This document's sibling
│       └── STRUCTURE.md      # This document
├── .github/                   # GitHub Actions workflows
├── README.md                  # Project documentation
├── LICENSE                    # MIT license
└── CHANGELOG.md              # Release history
```

## Grouping Pattern

**Hybrid: By feature for commands, by type for infrastructure**

### By feature (commands)
Each subcommand gets its own file in `internal/cmd/`:
- `install.go` contains only install logic
- `check.go` contains only check logic
- `uninstall.go` contains only uninstall logic
- `version.go` contains only version logic

Shared utilities (scope parsing, directory resolution) live in `scope.go` within the same package.

### By type (infrastructure)
Infrastructure packages group related functionality:
- `cli/` contains CLI dispatch and UI (usage text)
- `installer/` contains file installation logic and conflict handling
- `version/` contains version constants

## Naming Conventions

### Files
- **One export per file**: `installer.go` exports `Installer` struct
- **Files match primary concept**: `conflict.go` handles conflict resolution
- **Subcommands named after command**: `install.go` for `install` command

### Functions
- **Public command handlers**: `RunInstall`, `RunCheck`, `RunUninstall`, `RunVersion`
  - Pattern: `Run<Command>`
  - Consistent signature: `(args []string, ...) int`
  - Return exit code

- **Public utilities**: `ParseScope`, `ResolveDir`, `ParseConflictStrategy`
  - Verb + noun pattern
  - Parse functions return `(value, remaining, error)`

- **Private helpers**: `printUsage`, `cleanEmptyDirs`, `handleConflict`
  - Lowercase first letter
  - Verb + object pattern

### Types
- **Structs**: PascalCase (`Installer`)
- **Constants**: PascalCase with prefix (`ConflictKeep`, `ConflictReplace`)
- **String enums**: Type alias + const block
  ```go
  type ConflictStrategy string
  const (
      ConflictKeep    ConflictStrategy = "keep"
      ConflictReplace ConflictStrategy = "replace"
  )
  ```

### Variables
- **Package-level**: camelCase (`embeddedContent`, `Manifest`)
- **Exported constants**: PascalCase (`Version`, `GitCommit`)
- **Local variables**: camelCase, full words (`targetDir`, `contentFS`)

## Entry Points

### Binary entry point
`/Users/juliantellez/github.com/atlantic-blue/greenlight/main.go`

Single responsibility: Embed content and invoke CLI.

### CLI entry point
`/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/cli/cli.go:14`

Function: `Run(args []string, contentFS fs.FS) int`

Takes parsed args (without binary name) and embedded content filesystem.

### Command entry points
Each in `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/cmd/`:
- `RunInstall(args []string, contentFS fs.FS, stdout io.Writer) int`
- `RunUninstall(args []string, stdout io.Writer) int`
- `RunCheck(args []string, stdout io.Writer) int`
- `RunVersion(stdout io.Writer) int`

## Config File Locations

**No runtime config files**

Greenlight has no configuration files of its own. It installs configuration files for Claude Code:
- Global: `~/.claude/CLAUDE.md`
- Local: `./CLAUDE.md` (project root)

**Version tracking**
After installation, creates `.greenlight-version` in target directory:
- Global: `~/.claude/.greenlight-version`
- Local: `./.claude/.greenlight-version`

Format:
```
<version>
<commit>
<build-date>
```

## Test Locations

**No tests present**

As of this codebase snapshot, no test files exist. If added, Go convention would place them:
- `internal/cli/cli_test.go`
- `internal/cmd/install_test.go`
- `internal/installer/installer_test.go`

Test files sit alongside source files, using `_test.go` suffix.

## Package Organization

### `internal/` boundary
All implementation code lives under `internal/`, which prevents external imports. This is a single-binary tool, not a library.

### Package responsibilities

**`internal/cli`**
- Command-line interface coordination
- Subcommand routing
- Usage text and help output

**`internal/cmd`**
- Individual subcommand implementations
- Flag parsing
- Scope resolution
- Exit code logic

**`internal/installer`**
- File installation mechanics
- Manifest management
- Conflict resolution strategies
- Integrity checking (check command)
- Cleanup (uninstall command)

**`internal/version`**
- Build-time version information
- No logic, only variables set via ldflags

## Import Graph

```
main.go
 └─> internal/cli
      └─> internal/cmd
           ├─> internal/installer
           └─> internal/version
```

**No circular dependencies**

Dependency flow is strictly downward:
1. `main` → `cli`
2. `cli` → `cmd`
3. `cmd` → `installer` and `version`

No package imports anything above it in the hierarchy.

## Embedded Content Structure

Content lives in `/Users/juliantellez/github.com/atlantic-blue/greenlight/src/` and mirrors the installed directory structure:

**Agents** (`agents/`)
Eight agent definitions, each a markdown file describing a specialized agent role in the Greenlight framework.

**Commands** (`commands/gl/`)
Twelve command definitions for the `gl` slash command in Claude Code. Each defines a workflow step.

**References** (`references/`)
Three reference documents:
- Checkpoint protocol (saving progress)
- Deviation rules (handling unexpected work)
- Verification patterns (testing strategies)

**Templates** (`templates/`)
Two template files:
- `config.md`: Project configuration structure
- `state.md`: Project state tracking structure

**Root** (`CLAUDE.md`)
Main framework instructions loaded by Claude Code.

## File Permissions

**Consistent permissions across all operations**

Directories: `0o755` (rwxr-xr-x)
- Owner can read, write, execute
- Group and others can read and execute

Files: `0o644` (rw-r--r--)
- Owner can read and write
- Group and others can read only

Set in:
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/installer.go:118`
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/installer.go:122`
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/conflict.go:24`

## Manifest

**Central source of truth**: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/installer.go:15`

Hardcoded string slice containing all installed file paths (relative to content FS root).

Used by:
- `Install()`: Iterates manifest to copy files
- `Uninstall()`: Iterates manifest to remove files
- `Check()`: Iterates manifest to verify files

**Current count**: 41 files total

Adding new content requires:
1. Add markdown file to `src/` directory
2. Add path to manifest slice
3. Update `go:embed` directive in `main.go` if adding new subdirectory

## Build Artifacts

**Generated at build time**:
- Binary: `greenlight` (single executable)
- Version file: `.greenlight-version` (written during install)

**Not in source control**:
- `greenlight` binary (ignored)
- `dist/` (build output directory)

## Documentation

**User-facing**:
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/README.md`: Installation and usage
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/src/README.md`: Embedded content overview

**Developer-facing**:
- This document (`STRUCTURE.md`)
- `ARCHITECTURE.md` (sibling document)

**No inline code comments**

Code is self-documenting through function names and structure. Comments appear only for:
- Package-level documentation
- Build directive explanation (`//go:embed`)
- Exported function signatures (godoc format)
