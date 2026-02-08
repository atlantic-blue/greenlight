# Codebase Concerns

**Generated:** 2026-02-08
**Project:** greenlight
**Language:** Go 1.24 (primary), JavaScript (installer/tooling)

---

## Summary

This document identifies technical debt, security concerns, missing error handling, hardcoded values, and other issues requiring attention in the Greenlight codebase.

**Overall Health:** GOOD
**Critical Issues:** 0
**High Priority:** 2
**Medium Priority:** 4
**Low Priority:** 3

---

## HIGH Priority

### 1. Missing Test Coverage for All Production Code

**Priority:** HIGH
**Type:** Missing Tests
**Impact:** Cannot verify correctness, violates TDD-first principles

**Evidence:**
- No `*_test.go` files found in any package
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/cli/cli.go` (38 lines, no tests)
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/cmd/*.go` (multiple files, no tests)
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/installer.go` (238 lines, no tests)
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/conflict.go` (72 lines, no tests)
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/main.go` (no tests)

**Rationale:**
Greenlight is a TDD-first framework that enforces "tests are the source of truth" for users, yet the Greenlight codebase itself has zero test coverage. This is a credibility issue. The installer, conflict resolution, manifest validation, and file operations all have zero automated verification.

**Recommendation:**
Add integration tests for:
- Installation (global/local, conflict strategies)
- Uninstallation (cleanup verification)
- Check command (validation logic)
- Conflict handling (keep/replace/append strategies)
- Manifest file installation

---

### 2. Silent Error Swallowing in Main Entry Point

**Priority:** HIGH
**Type:** Error Handling
**Impact:** Installation failures exit silently without user feedback

**Location:** `/Users/juliantellez/github.com/atlantic-blue/greenlight/main.go:15-18`

```go
func main() {
    contentFS, err := fs.Sub(embeddedContent, "src")
    if err != nil {
        os.Exit(1)  // No error message to user
    }
    os.Exit(cli.Run(os.Args[1:], contentFS))
}
```

**Issue:**
If `fs.Sub` fails (corrupted binary, missing embedded files), the program exits with code 1 but prints nothing to stderr. Users see the command fail with no explanation.

**Recommendation:**
```go
if err != nil {
    fmt.Fprintf(os.Stderr, "error: failed to load embedded content: %v\n", err)
    os.Exit(1)
}
```

---

## MEDIUM Priority

### 3. Hardcoded Directory Names

**Priority:** MEDIUM
**Type:** Hardcoded Values
**Impact:** Limits flexibility, makes testing harder

**Evidence:**
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/cmd/scope.go:45` - `.claude` hardcoded
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/installer.go:218` - `.greenlight-version` hardcoded
- Manifest directories hardcoded in slice (lines 15-42)

**Issue:**
Directory names like `.claude`, `.greenlight`, and file names like `.greenlight-version` are hardcoded string literals scattered across multiple files. Changes require hunting down all references.

**Recommendation:**
Define constants in a single location:
```go
const (
    ClaudeDir = ".claude"
    GreenlightDir = ".greenlight"
    VersionFile = ".greenlight-version"
)
```

---

### 4. Large Functions Without Decomposition

**Priority:** MEDIUM
**Type:** Code Complexity
**Impact:** Hard to test, maintain, and reason about

**Evidence:**
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/gsd/bin/install.js:1069-1278` - `install()` function is 209 lines
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/gsd/bin/install.js:774-966` - `uninstall()` function is 192 lines
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/npm/bin/index.js:63-110` - `extractTarGz()` manual tar parsing is 47 lines

**Issue:**
Functions exceed 30-line guideline from `CLAUDE.md` engineering standards. Large functions mix concerns: validation, file operations, user messaging, and state management in the same scope.

**Recommendation:**
Extract subfunctions:
- `install()` → `installCommands()`, `installAgents()`, `installHooks()`, `configureSettings()`
- `uninstall()` → `removeCommands()`, `removeAgents()`, `cleanSettings()`
- `extractTarGz()` → `parseTarHeader()`, `extractFile()`

---

### 5. No Input Validation in Installer Commands

**Priority:** MEDIUM
**Type:** Security / Validation
**Impact:** Edge cases may cause unexpected behavior

**Evidence:**
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/cmd/install.go:12-32` - No validation of `args` array
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/cmd/scope.go:13-34` - Loops over `args` without bounds checking
- No validation that `targetDir` is a safe path (no parent traversal check)

**Issue:**
Functions process command-line arguments without validating:
- Argument count (too few/too many)
- Path safety (prevent `--global ../../../../etc/passwd`)
- Conflicting flags parsed in sequence without early exit

**Recommendation:**
Add validation:
```go
func ResolveDir(scope string) (string, error) {
    switch scope {
    case "global":
        home, err := os.UserHomeDir()
        if err != nil {
            return "", err
        }
        dir := filepath.Join(home, ".claude")
        // Validate dir is under home
        if !strings.HasPrefix(dir, home) {
            return "", fmt.Errorf("invalid path")
        }
        return dir, nil
    // ...
}
```

---

### 6. Missing go.sum Dependency Lock File

**Priority:** MEDIUM
**Type:** Dependency Management
**Impact:** Builds are not reproducible

**Evidence:**
- `go.mod` exists at `/Users/juliantellez/github.com/atlantic-blue/greenlight/go.mod`
- `go.sum` does NOT exist (glob returned no results)

**Issue:**
Without `go.sum`, Go cannot verify checksums of dependencies. This is a security and reproducibility issue. CI/CD builds may pull different versions of transitive dependencies.

**Recommendation:**
Run `go mod tidy` to generate `go.sum` and commit it.

---

## LOW Priority

### 7. Magic Numbers in Installer

**Priority:** LOW
**Type:** Code Quality
**Impact:** Readability

**Evidence:**
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/conflict.go:24` - `0o755` directory permissions
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/conflict.go:32` - `0o644` file permissions
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/installer.go:118` - `0o755`, `0o644` repeated

**Issue:**
File permission values appear as magic numbers. While Unix developers recognize `0o644`, they're unexplained.

**Recommendation:**
Use named constants:
```go
const (
    DirPerm  = 0o755 // Owner: rwx, Group: r-x, Other: r-x
    FilePerm = 0o644 // Owner: rw-, Group: r--, Other: r--
)
```

---

### 8. JavaScript Installer Has Deeply Nested Logic

**Priority:** LOW
**Type:** Code Complexity
**Impact:** Maintainability

**Evidence:**
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/gsd/bin/install.js` is 1,529 lines with:
  - 4-5 levels of nesting in `install()` function
  - Mixed concerns (CLI parsing, file operations, user prompts, config patching)
  - Large switch statements (lines 440-542 for frontmatter conversion)

**Issue:**
While functional, the installer mixes framework detection (Claude/OpenCode/Gemini), file transformations, user interaction, and configuration updates in a single file.

**Recommendation:**
Extract modules:
- `cli.js` - argument parsing, prompts
- `runtimes.js` - runtime-specific config (getGlobalDir, getDirName)
- `transforms.js` - frontmatter conversions
- `installer.js` - orchestration

---

### 9. No Explicit Cleanup on Installer Errors

**Priority:** LOW
**Type:** Error Handling
**Impact:** Partial installations may leave files

**Evidence:**
- `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/installer/installer.go:57-78` - If loop fails mid-iteration, earlier files remain
- No transaction-like rollback on partial install failure

**Issue:**
If installation fails (disk full, permission denied) partway through the manifest, already-copied files remain. Subsequent `install` attempts may see stale files.

**Recommendation:**
Use a staging directory:
1. Copy all files to `/tmp/greenlight-staging-{uuid}/`
2. Validate completeness
3. Atomic move to target directory
4. On error, delete staging directory

---

## Anti-Patterns Observed

### Inconsistent Error Handling Between Go and JavaScript

- **Go code:** Returns `error` values, callers handle with `if err != nil`
- **JavaScript code:** Mixes `try/catch`, error callbacks, and synchronous `fs` calls without consistent error propagation

**Impact:** Different reliability characteristics between the Go binary (robust) and the JS installer (fragile to environment issues).

---

## Code Smells

### 1. String Prefix Parsing Instead of Flag Library

**Location:** `/Users/juliantellez/github.com/atlantic-blue/greenlight/internal/cmd/scope.go:53-70`

```go
func ParseConflictStrategy(args []string) (installer.ConflictStrategy, []string) {
    // ...
    for _, arg := range args {
        if len(arg) > 14 && arg[:14] == "--on-conflict=" {
            val := arg[14:]
            // ...
```

**Issue:** Manual string slicing is error-prone. Go's `flag` package or a library like `cobra` would handle this cleanly.

---

### 2. Embedding Entire `src/` Directory

**Location:** `/Users/juliantellez/github.com/atlantic-blue/greenlight/main.go:11`

```go
//go:embed src/agents/*.md src/commands/gl/*.md src/references/*.md src/templates/*.md src/CLAUDE.md
var embeddedContent embed.FS
```

**Issue:** Embeds only specific patterns. If a new directory is added to `src/`, the embed directive must be updated manually. Easy to forget.

**Recommendation:** `//go:embed src` to embed the entire tree, or use a manifest file.

---

## Positive Observations

1. **No Secrets in Code:** Grep for `password|secret|token` found only documentation and template examples. No hardcoded credentials.

2. **Consistent Code Style:** All Go files are `gofmt`-compliant (verified with `gofmt -l`).

3. **Clear Separation of Concerns:** Go packages are well-organized: `cli/`, `cmd/`, `installer/`, `version/`.

4. **No Dead Code Detected:** All defined functions are used (based on cross-referencing imports and call sites).

5. **Explicit Version Management:** Version injection via ldflags (`-X github.com/.../version.Version=...`) is a best practice.

---

## TODO/FIXME/HACK Comments

None found in production Go code.

JavaScript installer and template files contain TODO patterns, but these are:
- Examples in documentation (e.g., `gsd/get-shit-done/templates/verification-report.md:69`)
- Comments about detecting TODOs in user code (not actual TODOs in the installer itself)

---

## Missing Features (Based on README Roadmap)

From `/Users/juliantellez/github.com/atlantic-blue/greenlight/README.md:227-235`:

- [ ] `greenlight upgrade` command (not implemented)
- [ ] `greenlight doctor` command (not implemented)
- [ ] Homebrew tap (not set up)
- [ ] Curl install script (not available)
- [ ] Plugin system (not implemented)

**Status:** These are documented future work, not current concerns.

---

## Dependency Analysis

### Go Dependencies

```
/Users/juliantellez/github.com/atlantic-blue/greenlight/go.mod:
module github.com/atlantic-blue/greenlight

go 1.24
```

**Zero external dependencies.** All code uses Go stdlib only. This is excellent for:
- Security surface (no transitive CVEs)
- Build reproducibility
- Binary size

### JavaScript Dependencies

**gsd/package.json:**
```json
"devDependencies": {
  "esbuild": "^0.24.0"
}
```

**npm/package.json:**
- Zero dependencies

**Analysis:** Minimal JS dependencies. Only bundler for hooks. No runtime dependencies in the npm wrapper.

---

## Configuration Files

### `.commitlintrc.json`

Extends `@commitlint/config-conventional`. Enforces conventional commits. Validated by lefthook pre-commit hook.

**Concern:** Dependency on `@commitlint/config-conventional` not listed in any `package.json`. Likely installed globally or in CI only. Not a project dependency.

### `.releaserc.json`

Semantic-release config for automated releases. Uses:
- `@semantic-release/commit-analyzer`
- `@semantic-release/release-notes-generator`
- `@semantic-release/changelog`
- `@semantic-release/git`
- `@semantic-release/github`

**Concern:** These plugins are not listed in any `package.json` in the repo. Implies CI environment must install them separately. Documented in CI config, but no local `package.json` for developers.

---

## Structural Observations

### Directory Layout

```
greenlight/
├── .github/          # CI workflows
├── .greenlight/      # Project-local state
├── gsd/              # Legacy "Get Shit Done" installer (1,529 lines JS)
├── internal/         # Go source (cli, cmd, installer, version)
├── npm/              # NPM wrapper binary (147 lines JS)
├── prototype/        # Prototype docs (appears unused in production)
├── src/              # Embedded markdown files (agents, commands, refs, templates)
└── main.go           # Entry point
```

**Concern:** Three installation paths:
1. Go binary (`greenlight install`)
2. NPM wrapper (`npx greenlight-cc install`)
3. GSD installer (`npx get-shit-done-cc`)

This is confusing. README says "Quick Start: `npx greenlight-cc install --global`" but the `gsd/` directory is a separate 1,500-line installer for a different product (Get Shit Done).

**Recommendation:** Clarify in README:
- `greenlight-cc` is the Greenlight installer
- `get-shit-done-cc` is a separate product
- Why both exist in the same repo

---

## Files That Should Not Be Committed

### `greenlight` Binary

**Location:** `/Users/juliantellez/github.com/atlantic-blue/greenlight/greenlight` (2.6 MB)

**Issue:** The compiled binary is committed to git. This is an anti-pattern:
- Bloats repo size
- Different builds for different platforms (this is `darwin/amd64` or `darwin/arm64`)
- Should be in `.gitignore`

**Evidence from `.gitignore`:**
```
greenlight
dist/
npm/node_modules/
```

The first line SHOULD exclude it, but the file is committed anyway (seen in `ls -la` output). This suggests it was force-added (`git add -f`) at some point.

**Recommendation:**
```bash
git rm greenlight
git commit -m "chore: remove committed binary"
```

---

## Summary Table

| Category | Count | Files Affected |
|----------|-------|----------------|
| Missing tests | 10 | All Go files |
| Silent error exits | 1 | main.go |
| Hardcoded values | 5+ | scope.go, installer.go |
| Large functions | 3 | install.js, extractTarGz |
| Missing validation | 3 | install.go, scope.go |
| Missing go.sum | 1 | Root |
| Magic numbers | 6 | conflict.go, installer.go |
| Deep nesting | 1 | install.js |
| No cleanup on error | 1 | installer.go |
| Committed binary | 1 | greenlight |

**Total Issues:** 10 distinct concerns

---

## Recommendations by Priority

**HIGH (Do First):**
1. Add test suite for Go packages (installer, conflict, cli, cmd)
2. Fix silent error in main.go (add stderr message)

**MEDIUM (Do Soon):**
3. Extract hardcoded directory names to constants
4. Refactor large functions in JS installer (install, uninstall)
5. Add path validation in ResolveDir
6. Generate and commit go.sum

**LOW (Do Eventually):**
7. Replace magic numbers with named constants
8. Modularize JS installer into separate files
9. Add atomic install with staging directory

**Cleanup:**
10. Remove committed `greenlight` binary from git
11. Clarify relationship between `greenlight-cc` and `get-shit-done-cc` in README

---

**End of Report**
