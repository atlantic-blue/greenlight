# Code Conventions

This document describes the code patterns, style, and conventions found in the greenlight codebase.

## Language & Version
- **Go 1.24**
- Single module: `github.com/atlantic-blue/greenlight`
- No external dependencies (uses only Go standard library)

## Project Structure

```
/Users/juliantellez/github.com/atlantic-blue/greenlight/
├── main.go                    # Entry point with embedded content
├── internal/
│   ├── cli/                   # Command dispatcher
│   ├── cmd/                   # Subcommand implementations
│   ├── installer/             # File installation logic
│   └── version/               # Build-time version info
└── src/                       # Embedded markdown content
```

**Pattern**: Domain-based organization. Files grouped by functional area, not by technical layer.

Evidence:
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/cli/cli.go`
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/cmd/version.go`
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/installer.go`

## Naming Conventions

### Packages
- **All lowercase, single word** where possible: `cli`, `cmd`, `installer`, `version`
- No underscores or hyphens in package names

### Functions
- **Exported functions**: PascalCase (`RunInstall`, `ParseScope`, `ResolveDir`)
- **Unexported functions**: camelCase (`printUsage`, `handleConflict`, `cleanEmptyDirs`, `copyFile`)
- **Naming pattern**: Verb-noun or action-based (`RunCheck`, `ParseConflictStrategy`, `writeVersionFile`)

### Variables
- **Exported constants**: PascalCase (`ConflictKeep`, `ConflictReplace`, `ConflictAppend`)
- **Exported variables**: PascalCase (`Manifest`, `Version`, `GitCommit`, `BuildDate`)
- **Local variables**: camelCase (`targetDir`, `destPath`, `srcData`)
- **Single-letter names**: Used sparingly for common cases (`w` for `io.Writer`, `err` for errors)

### Types
- **Exported types**: PascalCase (`ConflictStrategy`, `Installer`)
- **Descriptive names**: Type name indicates purpose (`ConflictStrategy` not `Strategy`)

Evidence:
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/conflict.go` (lines 11-17)
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/installer.go` (lines 44-48)

## Code Style

### Import Organization
- **Standard library only** (no third-party dependencies)
- Imports grouped in single block
- No blank lines between imports (all standard library)
- Alphabetically sorted

Example from `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/installer.go`:
```go
import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/atlantic-blue/greenlight/internal/version"
)
```

Pattern: Standard library first, then internal packages separated by blank line.

### Function Signatures
- **Explicit return types** always declared
- **Error handling**: Functions that can fail return `error` or `(result, error)`
- **Context injection**: `io.Writer` passed as parameter for testability (not direct `fmt.Println`)
- **Simple returns**: Return type `int` for exit codes, `bool` for success/failure checks

Examples:
```go
func Run(args []string, contentFS fs.FS) int
func ParseScope(args []string) (scope string, remaining []string, err error)
func Check(targetDir, scope string, stdout io.Writer) (ok bool)
```

Evidence:
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/cli/cli.go` (line 14)
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/cmd/scope.go` (line 13)

### Error Handling

**Pattern 1: Immediate exit on error (CLI commands)**
```go
if err != nil {
    fmt.Fprintf(stdout, "error: %v\n", err)
    return 1
}
```

Evidence: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/cmd/check.go` (lines 14-16)

**Pattern 2: Wrapped errors with context**
```go
if err := inst.Install(targetDir, scope, strategy); err != nil {
    return fmt.Errorf("installing %s: %w", relPath, err)
}
```

Evidence: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/installer.go` (lines 66-68)

**Pattern 3: Silent error suppression for cleanup operations**
```go
os.Remove(versionPath)  // No error check for cleanup
```

Evidence: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/installer.go` (line 148)

**Pattern 4: Check specific error types**
```go
if os.IsNotExist(err) {
    // No conflict — just write the file
    return os.WriteFile(destPath, srcData, 0o644)
}
```

Evidence: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/conflict.go` (lines 30-33)

**Anti-pattern found**: Silent error suppression in cleanup
- Location: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/installer.go` (line 148)
- Issue: `os.Remove(versionPath)` doesn't check errors
- Rationale: Acceptable for cleanup operations that shouldn't block program flow

### Exit Codes
- **0**: Success
- **1**: Error (all error conditions)
- Return codes propagate from subcommands to `main()`

Evidence: `/Users/juliantellez/github.com/atlantic-blue/greenlight/main.go` (line 19)

### File Permissions
- **Directories**: `0o755` (rwxr-xr-x)
- **Files**: `0o644` (rw-r--r--)
- Octal notation used consistently

Evidence: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/installer.go` (lines 118, 122)

## Common Patterns

### Command Dispatch Pattern
Central dispatcher with switch statement routing to subcommand handlers:

```go
switch args[0] {
case "install":
    return cmd.RunInstall(args[1:], contentFS, stdout)
case "uninstall":
    return cmd.RunUninstall(args[1:], stdout)
// ...
}
```

Evidence: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/cli/cli.go` (lines 22-38)

### Flag Parsing Pattern
Manual flag parsing using string prefix matching (no `flag` package):

```go
for _, arg := range args {
    if len(arg) > 14 && arg[:14] == "--on-conflict=" {
        val := arg[14:]
        // ...
    }
}
```

Evidence: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/cmd/scope.go` (lines 58-68)

**Rationale**: Simple, explicit control flow for minimal flag requirements.

### Embedded Filesystem Pattern
Content embedded at build time using `//go:embed`:

```go
//go:embed src/agents/*.md src/commands/gl/*.md src/references/*.md src/templates/*.md src/CLAUDE.md
var embeddedContent embed.FS
```

Evidence: `/Users/juliantellez/github.com/atlantic-blue/greenlight/main.go` (lines 11-12)

### Manifest-Driven Installation
File list defined as global variable, iterated for install/uninstall/check:

```go
var Manifest = []string{
    "agents/gl-architect.md",
    "agents/gl-codebase-mapper.md",
    // ...
}
```

Evidence: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/installer.go` (lines 15-42)

### Dependency Injection
Writers injected for testability:

```go
func New(contentFS fs.FS, stdout io.Writer) *Installer {
    return &Installer{contentFS: contentFS, stdout: stdout}
}
```

Evidence: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/installer.go` (lines 51-53)

### String-Based Enums
Type aliases with const values:

```go
type ConflictStrategy string

const (
    ConflictKeep    ConflictStrategy = "keep"
    ConflictReplace ConflictStrategy = "replace"
    ConflictAppend  ConflictStrategy = "append"
)
```

Evidence: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/conflict.go` (lines 10-17)

## Documentation

### Comments
- **Package-level comments**: Present on all packages
- **Function comments**: Present on exported functions
- **Comment style**: Complete sentences, starts with function name
- **Inline comments**: Used sparingly for complex logic

Example:
```go
// RunVersion prints version information to w.
func RunVersion(w io.Writer) int {
```

Evidence: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/cmd/version.go` (lines 10-11)

### Version Information
Build-time injection using `-ldflags`:

```go
// Set via ldflags at build time:
//
//	go build -ldflags "-X github.com/atlantic-blue/greenlight/internal/version.Version=1.0.0"
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)
```

Evidence: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/version/version.go` (lines 3-10)

## Anti-Patterns Found

### 1. Silent Main Function Error Handling
**Location**: `/Users/juliantellez/github.com/atlantic-blue/greenlight/main.go` (lines 15-18)

```go
contentFS, err := fs.Sub(embeddedContent, "src")
if err != nil {
    os.Exit(1)  // No error message logged
}
```

**Issue**: Fails silently without indicating what went wrong.

### 2. Magic Numbers in String Parsing
**Location**: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/cmd/scope.go` (lines 59-60)

```go
if len(arg) > 14 && arg[:14] == "--on-conflict=" {
```

**Issue**: Hard-coded length `14` makes code brittle. Should use `strings.HasPrefix()` or `len("--on-conflict=")`.

## Summary

The greenlight codebase follows clean Go idioms with a focus on simplicity:

- **Zero external dependencies** (standard library only)
- **Manual flag parsing** (no framework overhead)
- **Dependency injection** for testability (io.Writer parameters)
- **Explicit error handling** with wrapped context
- **Domain-driven structure** (not layered architecture)
- **Embedded content** for single-binary distribution

The code prioritizes readability and simplicity over abstraction. Error handling is explicit and consistent. The CLI design uses exit codes idiomatically. File permissions are appropriate for a system tool.
