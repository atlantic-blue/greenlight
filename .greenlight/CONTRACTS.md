# CONTRACTS.md -- Greenlight CLI Typed Contracts

> **Project:** Greenlight
> **Scope:** All boundaries in the MVP CLI
> **Stack:** Go 1.24, stdlib only
> **Date:** 2026-02-08

---

## Contract Index

| # | Contract | Boundary | Slice |
|---|----------|----------|-------|
| C-01 | VersionInfo | Command handler -> Version package | S-01 |
| C-02 | ParseScope | Command handler -> Flag parsing | S-02 |
| C-03 | ResolveDir | Command handler -> Flag parsing | S-02 |
| C-04 | ParseConflictStrategy | Command handler -> Flag parsing | S-02 |
| C-05 | ConflictStrategy | Installer -> Conflict handler | S-03 |
| C-06 | HandleConflict | Installer -> Conflict handler | S-03 |
| C-07 | InstallerNew | Command handler -> Installer | S-04 |
| C-08 | InstallerInstall | Command handler -> Installer | S-04 |
| C-09 | InstallerCheck | Command handler -> Installer | S-05 |
| C-10 | InstallerUninstall | Command handler -> Installer | S-06 |
| C-11 | RunVersion | CLI -> Command handler | S-01 |
| C-12 | RunInstall | CLI -> Command handler | S-04 |
| C-13 | RunCheck | CLI -> Command handler | S-05 |
| C-14 | RunUninstall | CLI -> Command handler | S-06 |
| C-15 | CLIRun | User -> CLI dispatcher | S-07 |
| C-16 | EntryPoint | OS -> main | S-07 |

---

## S-01: Version

### C-01: VersionInfo

```go
// Contract: VersionInfo
// Boundary: Command handler -> Version package
// Slice: S-01 (Version)

// Package-level variables set via ldflags at build time.
// No function call -- consumers read these directly.

// Exported state:
//   version.Version   string  // semantic version, default "dev"
//   version.GitCommit string  // git short SHA, default "unknown"
//   version.BuildDate string  // ISO 8601 date, default "unknown"

// Errors: none (variables always have values)

// Invariants:
// - All three variables are non-empty strings at all times
// - Default values are "dev", "unknown", "unknown" respectively
// - Values are overridden only via ldflags at compile time
```

### C-11: RunVersion

```go
// Contract: RunVersion
// Boundary: CLI dispatcher -> Command handler
// Slice: S-01 (Version)

// Input
//   w io.Writer  // destination for output

// Output
//   int  // exit code

// Signature:
//   func RunVersion(w io.Writer) int

// Behaviour:
// - Writes "greenlight <Version> (commit: <GitCommit>, built: <BuildDate>)\n" to w
// - Returns 0 always

// Errors: none (always succeeds)

// Invariants:
// - Output format is exactly: "greenlight %s (commit: %s, built: %s)\n"
// - Exit code is always 0
```

---

## S-02: Flag Parsing

### C-02: ParseScope

```go
// Contract: ParseScope
// Boundary: Command handler -> Flag parsing (pure function)
// Slice: S-02 (Flag Parsing)

// Input
//   args []string  // raw CLI args after the subcommand name

// Output
//   scope     string    // "global" or "local"
//   remaining []string  // args not consumed by scope parsing
//   err       error     // non-nil on invalid input

// Signature:
//   func ParseScope(args []string) (scope string, remaining []string, err error)

// Errors:
// - ErrBothScopes: when both --global and --local are present
//     message: "cannot specify both --global and --local"
// - ErrNoScope: when neither --global nor --local is present
//     message: "must specify --global or --local"

// Invariants:
// - scope is exactly "global" or "local" when err is nil
// - remaining never contains "--global" or "--local"
// - remaining preserves order of non-scope args
// - err is non-nil if and only if scope is ""
```

### C-03: ResolveDir

```go
// Contract: ResolveDir
// Boundary: Command handler -> Flag parsing (OS interaction for global)
// Slice: S-02 (Flag Parsing)

// Input
//   scope string  // "global" or "local"

// Output
//   dir string  // resolved target directory path
//   err error   // non-nil if resolution fails

// Signature:
//   func ResolveDir(scope string) (string, error)

// Behaviour:
// - "global" -> filepath.Join(os.UserHomeDir(), ".claude")
// - "local"  -> ".claude"

// Errors:
// - ErrHomeDirUnavailable: when os.UserHomeDir() fails (global scope only)
//     message: "cannot determine home directory: <underlying>"
// - ErrUnknownScope: when scope is not "global" or "local"
//     message: "unknown scope: <value>"

// Invariants:
// - For "local", the return value is always literally ".claude"
// - For "global", the return value is always an absolute path ending in ".claude"
// - err is non-nil if and only if dir is ""
```

### C-04: ParseConflictStrategy

```go
// Contract: ParseConflictStrategy
// Boundary: Command handler -> Flag parsing (pure function)
// Slice: S-02 (Flag Parsing)

// Input
//   args []string  // raw CLI args (may contain --on-conflict=<value>)

// Output
//   strategy  ConflictStrategy  // resolved strategy
//   remaining []string          // args not consumed by strategy parsing
//   err       error             // non-nil for invalid strategy value

// Corrected signature (TD-1):
//   func ParseConflictStrategy(args []string) (ConflictStrategy, []string, error)

// NOTE: Current code returns (ConflictStrategy, []string) — silently ignores
// invalid values. The corrected signature adds an error return per TD-1/UD-1.

// Behaviour:
// - Extracts first arg matching "--on-conflict=<value>"
// - Valid values: "keep", "replace", "append"
// - When flag is absent, defaults to ConflictKeep ("keep")
// - When flag has an invalid value, returns error (does NOT default silently)

// Errors:
// - ErrInvalidConflictStrategy: when --on-conflict value is not keep/replace/append
//     message: "invalid --on-conflict value: <value> (valid: keep, replace, append)"

// Invariants:
// - strategy is one of ConflictKeep, ConflictReplace, ConflictAppend when err is nil
// - remaining never contains "--on-conflict=..." args
// - When flag is absent entirely, strategy is ConflictKeep and err is nil
// - When flag has invalid value, err is non-nil (strict, no silent default)
```

---

## S-03: Conflict Handling

### C-05: ConflictStrategy Type

```go
// Contract: ConflictStrategy
// Boundary: Installer <-> Conflict handler (shared type)
// Slice: S-03 (Conflict Handling)

// Type definition:
//   type ConflictStrategy string

// Valid values:
//   ConflictKeep    ConflictStrategy = "keep"
//   ConflictReplace ConflictStrategy = "replace"
//   ConflictAppend  ConflictStrategy = "append"

// Invariants:
// - Only these three values are valid
// - Used as enum; compared by value equality
```

### C-06: HandleConflict

```go
// Contract: HandleConflict
// Boundary: Installer -> Conflict handler (filesystem side effect)
// Slice: S-03 (Conflict Handling)

// Input
//   destPath string           // absolute or relative path to target CLAUDE.md
//   srcData  []byte           // greenlight CLAUDE.md content from embedded FS
//   strategy ConflictStrategy // how to resolve if destPath already exists
//   w        io.Writer        // progress output destination

// Output
//   err error  // non-nil on filesystem failure or unknown strategy

// Signature:
//   func handleConflict(destPath string, srcData []byte, strategy ConflictStrategy, w io.Writer) error

// Behaviour by case:
//
// Case 1: destPath does not exist
//   - Write srcData to destPath with 0o644 permissions
//   - Strategy is irrelevant
//
// Case 2: destPath exists, strategy=keep
//   - Leave existing file untouched
//   - Write srcData to CLAUDE_GREENLIGHT.md in same directory
//   - Print: "  existing CLAUDE.md kept; greenlight version saved as CLAUDE_GREENLIGHT.md\n"
//
// Case 3: destPath exists, strategy=replace
//   - Write existing content to destPath + ".backup" (e.g. CLAUDE.md.backup)
//   - Overwrite destPath with srcData
//   - Print: "  existing CLAUDE.md backed up to <backupPath>\n"
//
// Case 4: destPath exists, strategy=append
//   - Append srcData to existing content
//   - If existing does not end with newline, insert one before appending
//   - Overwrite destPath with combined content
//   - Print: "  greenlight content appended to existing CLAUDE.md\n"
//
// Case 5: unknown strategy
//   - Return error

// Errors:
// - ErrUnknownStrategy: when strategy is not keep/replace/append
//     message: "unknown conflict strategy: <value>"
// - ErrFilesystemWrite: when os.WriteFile or os.MkdirAll fails
//     wraps underlying OS error
// - ErrFilesystemRead: when os.ReadFile fails (and file exists)
//     wraps underlying OS error
// - ErrBackupCreation: when writing backup file fails (replace strategy)
//     message: "creating backup: <underlying>"

// Invariants:
// - When destPath does not exist, file is written regardless of strategy
// - Keep strategy never modifies the existing CLAUDE.md
// - Replace strategy always creates a backup before overwriting
// - Append strategy never produces double newlines at the join point
// - Directories are created with 0o755 if needed
// - Files are written with 0o644
```

---

## S-04: Install

### C-07: InstallerNew

```go
// Contract: InstallerNew
// Boundary: Command handler -> Installer (constructor)
// Slice: S-04 (Install)

// Input
//   contentFS fs.FS     // embedded content filesystem (after fs.Sub)
//   stdout    io.Writer  // progress output destination

// Output
//   *Installer  // configured installer instance

// Signature:
//   func New(contentFS fs.FS, stdout io.Writer) *Installer

// Errors: none (constructor always succeeds)

// Invariants:
// - Returned Installer is non-nil
// - Installer holds references to contentFS and stdout (no copies)
```

### C-08: InstallerInstall

```go
// Contract: InstallerInstall
// Boundary: Command handler -> Installer (filesystem side effect)
// Slice: S-04 (Install)

// Input (method on *Installer)
//   targetDir string           // destination directory (e.g. "~/.claude" or ".claude")
//   scope     string           // "global" or "local" (affects CLAUDE.md placement)
//   strategy  ConflictStrategy // CLAUDE.md conflict resolution

// Output
//   err error  // non-nil on any failure

// Signature:
//   func (inst *Installer) Install(targetDir, scope string, strategy ConflictStrategy) error

// Behaviour:
// 1. For each file in Manifest (except CLAUDE.md):
//    a. Read file from contentFS
//    b. Create destination directory with 0o755 if needed
//    c. Write file to targetDir/<relPath> with 0o644
//    d. Print "  installed <relPath>\n"
// 2. For CLAUDE.md:
//    a. Resolve destination path based on scope:
//       - global: targetDir/CLAUDE.md
//       - local: parent of targetDir/CLAUDE.md (project root)
//    b. Delegate to handleConflict with strategy
//    c. Print "  installed CLAUDE.md -> <destPath>\n"
// 3. Write .greenlight-version file to targetDir containing:
//    version\ngitcommit\nbuilddate\n
// 4. Print "greenlight installed to <targetDir>\n"

// Errors:
// - ErrReadEmbeddedFile: when fs.ReadFile fails for a manifest file
//     message: "installing <relPath>: <underlying>"
// - ErrCreateDirectory: when os.MkdirAll fails
//     message: "installing <relPath>: <underlying>"
// - ErrWriteFile: when os.WriteFile fails
//     message: "installing <relPath>: <underlying>"
// - ErrCLAUDEConflict: when handleConflict fails
//     message: "installing CLAUDE.md: <underlying>"
// - ErrWriteVersionFile: when writing .greenlight-version fails
//     message: "writing version file: <underlying>"

// Invariants:
// - All Manifest files are written before the version file
// - CLAUDE.md is placed asymmetrically: global=inside targetDir, local=parent of targetDir
// - Operation stops on first error (no partial-then-continue)
// - File permissions: directories 0o755, files 0o644
// - Version file format: three lines, each newline-terminated
// - Install is idempotent: running twice produces same result
```

### C-12: RunInstall

```go
// Contract: RunInstall
// Boundary: CLI dispatcher -> Command handler
// Slice: S-04 (Install)

// Input
//   args      []string  // CLI args after "install" (flags)
//   contentFS fs.FS     // embedded content filesystem
//   stdout    io.Writer // output destination

// Output
//   int  // exit code: 0 success, 1 failure

// Signature:
//   func RunInstall(args []string, contentFS fs.FS, stdout io.Writer) int

// Behaviour:
// 1. Parse conflict strategy from args (may fail with error per TD-1)
// 2. Parse scope from remaining args
// 3. Resolve target directory from scope
// 4. Create Installer, call Install
// 5. On any error: print "error: <message>\n" to stdout, return 1

// Errors (surfaced as exit code 1):
// - Invalid --on-conflict value
// - Missing or conflicting scope flags
// - Home directory resolution failure
// - Any filesystem error during install

// Invariants:
// - Exit code is 0 if and only if install completed successfully
// - All error messages are prefixed with "error: "
// - Conflict strategy is parsed before scope
```

---

## S-05: Check

### C-09: InstallerCheck

```go
// Contract: InstallerCheck
// Boundary: Command handler -> Installer (read-only filesystem operation)
// Slice: S-05 (Check)

// Corrected signature (TD-4):
//   func Check(targetDir, scope string, stdout io.Writer, verify bool, contentFS fs.FS) bool

// NOTE: Current code is Check(targetDir, scope string, stdout io.Writer) bool.
// The corrected signature adds verify and contentFS parameters per TD-4/UD-4.

// Input
//   targetDir string    // directory to check
//   scope     string    // "global" or "local" (affects CLAUDE.md path resolution)
//   stdout    io.Writer // output destination
//   verify    bool      // when true, compare content hashes against embedded source
//   contentFS fs.FS     // embedded content (required when verify=true, may be nil otherwise)

// Output
//   ok bool  // true if all files pass checks

// Behaviour (presence-only mode, verify=false):
// 1. For each file in Manifest:
//    a. Resolve path (CLAUDE.md uses asymmetric placement)
//    b. Stat the file
//    c. If missing: print "  MISSING  <relPath>\n", set ok=false
//    d. If stat error: print "  ERROR    <relPath>: <err>\n", set ok=false
//    e. If empty (size 0): print "  EMPTY    <relPath>\n", set ok=false
// 2. Check .greenlight-version:
//    a. If missing: print "  MISSING  .greenlight-version\n", set ok=false
//    b. If present: print "  version: <first line>\n"
// 3. Print summary:
//    a. All present: "all <N> files present\n"
//    b. Failures: "<present>/<total> files present (<missing> missing, <empty> empty)\n"

// Behaviour (content verification mode, verify=true):
// 1-2. Same as above, plus:
//    c-extra. If file exists and non-empty: compare SHA-256 hash against embedded source
//    d-extra. If hash mismatch: print "  MODIFIED <relPath>\n", set ok=false
// 3. Print summary:
//    a. All verified: "all <N> files verified\n"
//    b. Failures: "<ok>/<total> files verified (<missing> missing, <empty> empty, <modified> modified)\n"

// Errors: none (returns bool, prints diagnostics)

// Invariants:
// - Check is read-only: never modifies filesystem
// - CLAUDE.md path resolution matches Install placement logic
// - ok is true if and only if every manifest file is present and non-empty
//   (and content-matches when verify=true)
// - Summary line always printed as last output
// - verify=false does not require contentFS (may be nil)
```

### C-13: RunCheck

```go
// Contract: RunCheck
// Boundary: CLI dispatcher -> Command handler
// Slice: S-05 (Check)

// Corrected signature (TD-4):
//   func RunCheck(args []string, contentFS fs.FS, stdout io.Writer) int

// NOTE: Current code is RunCheck(args []string, stdout io.Writer) int.
// The corrected signature adds contentFS per TD-4/UD-4 for --verify support.

// Input
//   args      []string  // CLI args after "check" (flags)
//   contentFS fs.FS     // embedded content filesystem (for --verify)
//   stdout    io.Writer // output destination

// Output
//   int  // exit code: 0 all checks pass, 1 any failure

// Behaviour:
// 1. Parse scope from args
// 2. Parse --verify flag from remaining args
// 3. Resolve target directory from scope
// 4. Call Check with verify flag and contentFS
// 5. Return 0 if Check returns true, 1 otherwise

// Errors (surfaced as exit code 1):
// - Missing or conflicting scope flags
// - Home directory resolution failure
// - Any check failure (missing, empty, or modified files)

// Invariants:
// - Exit code is 0 if and only if all files pass all checks
// - All error messages are prefixed with "error: "
// - --verify is a boolean flag (presence = true, absence = false)
```

---

## S-06: Uninstall

### C-10: InstallerUninstall

```go
// Contract: InstallerUninstall
// Boundary: Command handler -> Installer (filesystem side effect)
// Slice: S-06 (Uninstall)

// Corrected signature (TD-3):
//   func Uninstall(targetDir, scope string, stdout io.Writer) error

// NOTE: Current code is Uninstall(targetDir string, stdout io.Writer) error.
// The corrected signature adds scope for CLAUDE.md artifact path resolution.

// Input
//   targetDir string    // directory to uninstall from
//   scope     string    // "global" or "local" (for conflict artifact location)
//   stdout    io.Writer // output destination

// Output
//   err error  // non-nil on filesystem errors (not counting missing files)

// Behaviour:
// 1. For each file in Manifest (except CLAUDE.md):
//    a. Attempt to remove targetDir/<relPath>
//    b. If file missing: skip without error (idempotent)
//    c. If other error: return error
//    d. Print "  removed <relPath>\n"
// 2. Remove .greenlight-version:
//    a. Attempt removal, report error if non-NotExist error
// 3. Remove conflict artifacts (TD-3):
//    a. Resolve CLAUDE.md directory (same asymmetric logic as install)
//    b. Remove CLAUDE_GREENLIGHT.md if present, print "  removed CLAUDE_GREENLIGHT.md\n"
//    c. Remove CLAUDE.md.backup if present, print "  removed CLAUDE.md.backup\n"
// 4. Clean up empty directories deepest-first:
//    - commands/gl, commands, agents, references, templates
// 5. Print "greenlight uninstalled from <targetDir>\n"

// Errors:
// - ErrRemoveFile: when os.Remove fails with a non-NotExist error
//     message: "removing <relPath>: <underlying>"
// - ErrRemoveVersionFile: when removing .greenlight-version fails (non-NotExist)

// Invariants:
// - CLAUDE.md is NEVER removed (may contain user content) -- NFR-3
// - Missing files are skipped without error (idempotent) -- NFR-6
// - Conflict artifacts (CLAUDE_GREENLIGHT.md, CLAUDE.md.backup) ARE removed -- TD-3
// - Empty directories are cleaned up deepest-first
// - Uninstall is idempotent: running twice produces same result
```

### C-14: RunUninstall

```go
// Contract: RunUninstall
// Boundary: CLI dispatcher -> Command handler
// Slice: S-06 (Uninstall)

// Input
//   args   []string  // CLI args after "uninstall" (flags)
//   stdout io.Writer // output destination

// Output
//   int  // exit code: 0 success, 1 failure

// Signature:
//   func RunUninstall(args []string, stdout io.Writer) int

// Behaviour:
// 1. Parse scope from args
// 2. Resolve target directory from scope
// 3. Call Uninstall with targetDir, scope, stdout
// 4. On error: print "error: <message>\n", return 1

// Errors (surfaced as exit code 1):
// - Missing or conflicting scope flags
// - Home directory resolution failure
// - Any filesystem error during uninstall

// Invariants:
// - Exit code is 0 if and only if uninstall completed without filesystem errors
// - All error messages are prefixed with "error: "
```

---

## S-07: CLI Dispatch

### C-15: CLIRun

```go
// Contract: CLIRun
// Boundary: User -> CLI dispatcher
// Slice: S-07 (CLI Dispatch)

// Corrected signature (TD-2):
//   func Run(args []string, contentFS fs.FS, stdout io.Writer) int

// NOTE: Current code is Run(args []string, contentFS fs.FS) int with hardcoded
// os.Stdout. The corrected signature adds io.Writer per TD-2/UD-2.

// Input
//   args      []string  // os.Args[1:] — subcommand and flags
//   contentFS fs.FS     // embedded content filesystem (after fs.Sub)
//   stdout    io.Writer // all output destination

// Output
//   int  // exit code: 0 success, 1 failure

// Behaviour:
// - No args:                print usage, return 0
// - "help"|"--help"|"-h":   print usage, return 0
// - "install":              dispatch to RunInstall(args[1:], contentFS, stdout)
// - "uninstall":            dispatch to RunUninstall(args[1:], stdout)
// - "check":                dispatch to RunCheck(args[1:], contentFS, stdout)
// - "version":              dispatch to RunVersion(stdout)
// - anything else:          print "unknown command: <cmd>\n\n", print usage, return 1

// Errors: none directly (delegates error handling to subcommands)

// Invariants:
// - Exactly one subcommand handler is called per invocation
// - Unknown commands print the invalid command name in the error message
// - Usage is printed for: no args, help variants, unknown commands
// - The return value is always the exit code from the dispatched handler
//   (or 0 for help, or 1 for unknown command)
// - stdout parameter is passed through to all subcommand handlers (TD-2)
// - contentFS is passed to install and check handlers
```

### C-16: EntryPoint

```go
// Contract: EntryPoint
// Boundary: OS -> main (process entry)
// Slice: S-07 (CLI Dispatch)

// Corrected behaviour (FR-8.1):
//   Print fs.Sub error to stderr before os.Exit(1)

// Behaviour:
// 1. Compute contentFS via fs.Sub(embeddedContent, "src")
// 2. If fs.Sub fails: print error to os.Stderr, os.Exit(1)
// 3. Call cli.Run(os.Args[1:], contentFS, os.Stdout)
// 4. os.Exit with the returned exit code

// Invariants:
// - fs.Sub error is printed to stderr (not swallowed) -- FR-8.1
// - os.Exit is called with the return value of cli.Run
// - os.Stdout is passed as the io.Writer to cli.Run (TD-2)
```

---

## Cross-Cutting: Manifest

```go
// Contract: Manifest
// Boundary: Installer internal (shared across Install, Check, Uninstall)
// Not a separate slice -- used by S-04, S-05, S-06

// Type:
//   var Manifest []string

// Content: 26 relative paths (from embedded FS root after fs.Sub)
// See DESIGN.md section 4.1 for full listing.

// Invariants:
// - Manifest is the single source of truth for which files are managed
// - All paths are forward-slash separated (Go fs.FS convention)
// - CLAUDE.md is always the last entry
// - Manifest is never modified at runtime (read-only)
// - Install writes exactly these files (plus .greenlight-version)
// - Check verifies exactly these files (plus .greenlight-version)
// - Uninstall removes exactly these files minus CLAUDE.md (plus .greenlight-version, plus conflict artifacts)
```

---

## Summary of Corrections Encoded in Contracts

| Contract | Current Signature | Corrected Signature | Design Decision |
|----------|------------------|---------------------|-----------------|
| C-04 ParseConflictStrategy | `([]string) (ConflictStrategy, []string)` | `([]string) (ConflictStrategy, []string, error)` | TD-1/UD-1 |
| C-09 InstallerCheck | `(string, string, io.Writer) bool` | `(string, string, io.Writer, bool, fs.FS) bool` | TD-4/UD-4 |
| C-10 InstallerUninstall | `(string, io.Writer) error` | `(string, string, io.Writer) error` | TD-3/UD-3 |
| C-13 RunCheck | `([]string, io.Writer) int` | `([]string, fs.FS, io.Writer) int` | TD-4/UD-4 |
| C-15 CLIRun | `([]string, fs.FS) int` | `([]string, fs.FS, io.Writer) int` | TD-2/UD-2 |
| C-16 EntryPoint | Silent exit on fs.Sub error | Print to stderr before exit | FR-8.1 |
