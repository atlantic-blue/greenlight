# Testing Status

## Summary

**Test coverage: 0%**

The greenlight codebase currently has **zero automated tests**. No test files exist.

## Test Framework

**None configured.**

Standard Go testing infrastructure is available but not used:
- Go 1.24 includes `testing` package
- No test files found matching `*_test.go` pattern
- No test configuration files (no `go.mod` test dependencies)

## Test File Search Results

```bash
# Search performed across entire project
find /Users/juliantellez/github.com/atlantic-blue/greenlight -name "*_test.go"
# Result: No files found
```

## Missing Test Coverage

All production code is untested. Critical areas without tests:

### 1. CLI Command Routing
**File**: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/cli/cli.go`

Untested scenarios:
- Command dispatch logic (switch statement, lines 22-38)
- Unknown command handling
- Help flag variations
- Empty args handling

### 2. Flag Parsing
**File**: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/cmd/scope.go`

Untested scenarios:
- `ParseScope()`: Conflicting flags (--global and --local together)
- `ParseScope()`: Missing scope flags
- `ResolveDir()`: Home directory resolution failure
- `ParseConflictStrategy()`: Invalid strategy values
- Edge case: `--on-conflict=` with empty value

### 3. File Installation Logic
**File**: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/installer.go`

Untested scenarios:
- `Install()`: Manifest file iteration and copying
- `Install()`: Directory creation (MkdirAll failure cases)
- `installCLAUDE()`: Asymmetric placement (global vs local)
- `copyFile()`: Embedded FS read failures
- `Uninstall()`: Partial failure (some files don't exist)
- `Check()`: File existence and size validation
- Version file writing and parsing

### 4. Conflict Resolution
**File**: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/conflict.go`

Untested scenarios:
- `handleConflict()`: All three strategies (keep, replace, append)
- Edge case: Existing file without trailing newline (append mode)
- Edge case: Backup file creation failure (replace mode)
- Edge case: Permission denied when creating alternate file (keep mode)

### 5. Main Entry Point
**File**: `/Users/juliantellez/github.com/atlantic-blue/greenlight/main.go`

Untested scenarios:
- Embedded FS subdirectory extraction
- Exit code propagation from CLI

### 6. Subcommands
**Files**: `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/cmd/*.go`

Untested scenarios:
- `RunInstall()`: Error handling and integration
- `RunUninstall()`: Cleanup and directory removal
- `RunCheck()`: Verification logic
- `RunVersion()`: Output formatting

## Test Patterns Observed

**None.** No existing tests to establish patterns.

## How To Run Tests (If They Existed)

Standard Go testing commands would apply:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run tests in specific package
go test ./internal/installer

# Run specific test
go test -run TestInstall ./internal/installer
```

## Recommended Test Structure

Based on codebase structure, tests should follow:

```
/Users/juliantellez/github.com/atlantic-blue/greenlight/
├── internal/
│   ├── cli/
│   │   ├── cli.go
│   │   └── cli_test.go          # MISSING
│   ├── cmd/
│   │   ├── version.go
│   │   ├── version_test.go      # MISSING
│   │   ├── check.go
│   │   ├── check_test.go        # MISSING
│   │   ├── install.go
│   │   ├── install_test.go      # MISSING
│   │   ├── uninstall.go
│   │   ├── uninstall_test.go    # MISSING
│   │   ├── scope.go
│   │   └── scope_test.go        # MISSING
│   ├── installer/
│   │   ├── installer.go
│   │   ├── installer_test.go    # MISSING
│   │   ├── conflict.go
│   │   └── conflict_test.go     # MISSING
│   └── version/
│       ├── version.go
│       └── version_test.go      # MISSING (likely not needed - just vars)
└── main_test.go                 # MISSING
```

## Test Types Needed

### Unit Tests
All packages need unit tests:

1. **installer package**
   - Test conflict strategies independently
   - Test file copying with mock filesystem
   - Test manifest iteration
   - Test version file format

2. **cmd package**
   - Test flag parsing with various inputs
   - Test scope resolution
   - Test directory resolution edge cases

3. **cli package**
   - Test command routing
   - Test help text output
   - Test unknown command handling

### Integration Tests
End-to-end scenarios needed:

1. **Install flow**
   - Global install to temp directory
   - Local install to temp directory
   - Install with each conflict strategy
   - Verify all manifest files created
   - Verify version file created

2. **Uninstall flow**
   - Remove all files except CLAUDE.md
   - Clean up empty directories
   - Handle missing files gracefully

3. **Check flow**
   - Detect missing files
   - Detect empty files
   - Verify version file parsing

### Table-Driven Tests
Suitable for:
- Flag parsing (`ParseScope`, `ParseConflictStrategy`)
- Conflict strategies (`handleConflict`)
- Command dispatch (`cli.Run`)

Example structure:
```go
func TestParseScope(t *testing.T) {
    tests := []struct {
        name      string
        args      []string
        wantScope string
        wantErr   bool
    }{
        {"global flag", []string{"--global"}, "global", false},
        {"local flag", []string{"--local"}, "local", false},
        {"both flags", []string{"--global", "--local"}, "", true},
        {"no flags", []string{}, "", true},
    }
    // ...
}
```

## Testability Analysis

### Well-Designed for Testing
1. **Dependency injection**: `io.Writer` parameters allow output capture
2. **Return exit codes**: Functions return `int` for verification
3. **Filesystem abstraction**: Uses `fs.FS` interface (mockable)
4. **Pure functions**: Many functions have no side effects beyond filesystem

Evidence:
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/installer.go` (line 51): `New(contentFS fs.FS, stdout io.Writer)`
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/cli/cli.go` (line 14): `Run(args []string, contentFS fs.FS) int`

### Testing Challenges
1. **Filesystem operations**: Heavy use of `os` package (need temp directories or mocking)
2. **Global state**: `os.UserHomeDir()` reads from environment
3. **Embedded content**: Main function embeds real files (need test fixtures)

## Test Utilities Needed

To effectively test this codebase:

1. **Filesystem helpers**
   ```go
   // Create temporary test directory
   // Populate with test files
   // Clean up after test
   ```

2. **Output capture**
   ```go
   // Already possible via io.Writer injection
   var buf bytes.Buffer
   exitCode := cmd.RunVersion(&buf)
   output := buf.String()
   ```

3. **Mock embedded FS**
   ```go
   // Use testing/fstest.MapFS
   testFS := fstest.MapFS{
       "CLAUDE.md": {Data: []byte("test content")},
   }
   ```

## Build Integration

No test automation configured:
- No CI test runs in `.github/workflows` (not checked, but no tests exist)
- No pre-commit test hooks in `.lefthook.yml` (exists but contents not read)
- No coverage requirements
- No test timeouts

## Recommendations

### Priority 1: Core Logic Tests
Start with pure functions and flag parsing:
1. `ParseScope()` - straightforward, no dependencies
2. `ParseConflictStrategy()` - straightforward, no dependencies
3. `handleConflict()` - needs temp files, critical logic

### Priority 2: Integration Tests
End-to-end flows with temporary directories:
1. Full install/uninstall/check cycle
2. Conflict strategy verification
3. Error scenarios (permissions, missing directories)

### Priority 3: CLI Tests
Command routing and output:
1. Command dispatch
2. Help text
3. Error messages

### Test Commands

```bash
# Run tests
go test ./...

# Run tests with coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run tests with race detector
go test -race ./...

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test ./internal/installer

# Run specific test function
go test -run TestInstall ./internal/installer

# Benchmark tests (none exist yet)
go test -bench=. ./...
```

## Coverage Goals

Based on Greenlight engineering standards:

- **Target**: 80%+ coverage for business logic
- **Critical paths**: 100% coverage for install/uninstall/conflict handling
- **CLI routing**: 100% coverage (all commands reachable)
- **Error paths**: All error conditions should be tested

## Current State: RED

In TDD terms, this codebase is in a permanent RED state:
- No tests exist
- No way to verify behavior beyond manual testing
- No regression protection
- Cannot confidently refactor

**Next step**: Add tests before making any changes to production code.
