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

---

# Brownfield MVP Contracts

> **Scope:** `/gl:assess`, `/gl:wrap`, brownfield-aware updates to existing commands
> **Deliverables:** Markdown prompt files (agents + commands), NOT Go code
> **Date:** 2026-02-08
> **Design Reference:** DESIGN.md sections 1-9, Technical Decisions TD-1 through TD-12

---

## Brownfield Contract Index

| # | Contract | Boundary | Slice |
|---|----------|----------|-------|
| C-17 | AssessorAgentBehaviour | /gl:assess orchestrator -> gl-assessor agent | S-08 |
| C-18 | AssessOrchestration | User -> /gl:assess command | S-08 |
| C-19 | AssessOutput | gl-assessor agent -> Filesystem (ASSESS.md) | S-08 |
| C-20 | AssessSecuritySpawn | /gl:assess orchestrator -> gl-security agent | S-08 |
| C-21 | WrapperAgentBehaviour | /gl:wrap orchestrator -> gl-wrapper agent | S-09 |
| C-22 | WrapOrchestration | User -> /gl:wrap command | S-09 |
| C-23 | WrapContractExtraction | gl-wrapper agent -> Filesystem (CONTRACTS.md) | S-09 |
| C-24 | WrapLockingTests | gl-wrapper agent -> Filesystem (tests/locking/) | S-09 |
| C-25 | WrapSecurityBaseline | /gl:wrap orchestrator -> gl-security agent | S-09 |
| C-26 | WrapStateTracking | /gl:wrap orchestrator -> Filesystem (STATE.md) | S-09 |
| C-27 | StatusBrownfieldDisplay | /gl:status command -> User (wrapped boundaries) | S-10 |
| C-28 | HelpBrownfieldSection | /gl:help command -> User (brownfield commands) | S-10 |
| C-29 | SettingsBrownfieldAgents | /gl:settings command -> User (assessor, wrapper) | S-10 |
| C-30 | SliceWrapsField | /gl:slice command -> Locking-to-integration transition | S-11 |
| C-31 | ArchitectWrappedContracts | gl-architect agent -> CONTRACTS.md ([WRAPPED] awareness) | S-11 |
| C-32 | TestWriterLockingAwareness | gl-test-writer agent -> Locking test name context | S-11 |
| C-33 | ManifestBrownfieldUpdate | Go CLI -> Manifest (4 new file paths) | S-12 |
| C-34 | ConfigProfileBrownfield | config.json template -> Profile (assessor, wrapper entries) | S-12 |
| C-35 | CLAUDEmdIsolationUpdate | CLAUDE.md -> Agent Isolation Rules (gl-assessor, gl-wrapper rows) | S-12 |

---

## S-08: Codebase Assessment (/gl:assess)

*User Action: "User can assess an existing codebase for gaps, risks, and untested boundaries"*

### C-17: AssessorAgentBehaviour

```
Contract: AssessorAgentBehaviour
Boundary: /gl:assess orchestrator -> gl-assessor agent (markdown prompt file)
Slice: S-08 (Codebase Assessment)
Design refs: FR-1, FR-2, FR-3, FR-4, FR-5, NFR-1, TD-3, TD-5, TD-9

AGENT DEFINITION: src/agents/gl-assessor.md

Input (context passed via Task spawn):
  - .greenlight/codebase/ docs (ARCHITECTURE.md, STRUCTURE.md from /gl:map)
  - .greenlight/config.json (stack, test commands, project directories)
  - CLAUDE.md engineering standards (gap comparison baseline)
  - Source code (read access to entire codebase)
  - Test files (read access)

Output (files written directly by agent):
  - .greenlight/ASSESS.md (structured assessment following schema in DESIGN.md 4.1)

Tools: Read, Bash, Glob, Grep
Model: Resolved from config.json profiles (default: sonnet in balanced profile)

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | NoCodabaseDocs | /gl:map has not been run; .greenlight/codebase/ is missing or empty | Warn user results will be shallow. Proceed with direct codebase scanning. Recommend running /gl:map first. (FR-1.3) |
  | NoConfig | .greenlight/config.json does not exist | Orchestrator catches this BEFORE spawning agent. Agent is never spawned. (FR-1.2) |
  | NoCoverageCommand | config.test.coverage_command is not configured | Note in ASSESS.md that line/branch coverage percentages are unavailable. File mapping still runs. (FR-2.6) |
  | SecurityAgentFailure | gl-security agent fails or returns no output | Continue assessment without security findings. Note in ASSESS.md: "Security scan not performed: {reason}". (NFR-8) |
  | ContextBudgetExceeded | Codebase too large for single agent context | Split analysis by directory/module. Assess each partition separately and aggregate findings into single ASSESS.md. (NFR-1) |

Invariants:
  - Agent is entirely read-only: MUST NOT modify any source code, test files, or config files
  - The ONLY file the agent writes is .greenlight/ASSESS.md
  - Agent reads CLAUDE.md standards and compares codebase against each section (FR-5.1)
  - Test coverage detection uses file mapping (always) plus coverage command (only if configured) per TD-3
  - Boundaries are classified as explicit/implicit/none with source file and line range (FR-3.2, FR-3.5)
  - Wrap recommendations are grouped by priority tier: Critical, High, Medium per TD-5
  - Each recommended boundary includes: name, type, contract status, test status, estimated complexity, risk level (FR-6.3)
  - Agent MUST complete within 50% context window (NFR-1)
  - Assessment output follows the exact ASSESS.md schema defined in DESIGN.md section 4.1

Security:
  - Agent CANNOT modify any code (read-only analytical agent)
  - Agent isolation: Can See (codebase docs, test results, standards), Cannot Do (modify any code)
  - Per FR-20.1: gl-assessor row in CLAUDE.md isolation table

Dependencies: None (assess is always available, FR-1.4)
```

### C-18: AssessOrchestration

```
Contract: AssessOrchestration
Boundary: User -> /gl:assess command (markdown prompt file)
Slice: S-08 (Codebase Assessment)
Design refs: FR-1, FR-4.1, FR-6.4, FR-6.5, NFR-3, NFR-6

COMMAND DEFINITION: src/commands/gl/assess.md

Input:
  - User invokes /gl:assess (no arguments)
  - .greenlight/config.json MUST exist

Output:
  - .greenlight/ASSESS.md written (created or overwritten per NFR-3)
  - Commit with message: "docs: greenlight codebase assessment" (FR-6.4)
  - Summary report displayed to user (FR-6.5)

Orchestration steps (follows existing Greenlight orchestration pattern from DESIGN.md 3.4):
  1. Read .greenlight/config.json for project context and model resolution
  2. Read .greenlight/codebase/ docs if they exist (from /gl:map)
  3. Resolve gl-assessor model from config.json profiles
  4. Spawn gl-assessor agent via Task with structured XML context blocks
  5. Spawn gl-security agent in "full-audit" mode (FR-4.1) -- may run in parallel
  6. Aggregate security findings into ASSESS.md if security agent succeeds
  7. Verify ASSESS.md exists and is non-empty
  8. Commit ASSESS.md with conventional format
  9. Display summary report to user
  10. Recommend next action: /gl:wrap or /gl:design

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | NoConfigJson | .greenlight/config.json does not exist | Print "No config found. Run /gl:init first." and stop |
  | AssessorFailure | gl-assessor agent fails to produce ASSESS.md | Report error to user. Do not commit. Suggest retrying |
  | SecurityFailure | gl-security agent fails | Continue without security findings. Warn user. Note in ASSESS.md (NFR-8) |

Invariants:
  - /gl:assess is idempotent: running multiple times overwrites ASSESS.md (NFR-3)
  - /gl:assess is read-only except for writing ASSESS.md (NFR-6)
  - No other command is required to have been run first (except /gl:init for config.json) (FR-1.4)
  - /gl:map is recommended but not required (FR-1.3)
  - User interaction is minimal: mostly automated with progress reports (DESIGN.md 5.1)
  - Summary report format follows DESIGN.md section 5.1 exactly

Security:
  - Command does not modify production code
  - Security agent runs in full-audit mode (same as /gl:ship)
  - Security findings are documented, not enforced (DESIGN.md 6.1)

Dependencies: None
```

### C-19: AssessOutput

```
Contract: AssessOutput
Boundary: gl-assessor agent -> Filesystem (.greenlight/ASSESS.md)
Slice: S-08 (Codebase Assessment)
Design refs: FR-6, DESIGN.md 4.1

FILE SPECIFICATION: .greenlight/ASSESS.md

Output schema (mandatory sections):
  - Summary table: source files, test files, file coverage, line coverage (or "not configured"),
    boundaries identified, explicit/implicit/none counts, security findings, standards compliance
  - Test Coverage: by-module table (module, source files, test files, coverage, status) +
    untested files table (file, type, risk, recommended priority)
  - Contract Inventory: boundaries table (boundary, type, contract status, location, tests) +
    summary by status (explicit/implicit/none counts and percentages)
  - Risk Assessment: security findings table (severity, category, location, description) +
    fragile areas table (file, concern, severity, detail) +
    tech debt table (file, type, detail)
  - Architecture Gaps: standards compliance table (per CLAUDE.md section: pass/partial/fail + key gaps) +
    specific violations table (standard, violation, location, severity)
  - Wrap Recommendations: three tiers (Critical, High, Medium) each with
    boundary, type, rationale, estimated complexity

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | EmptyAssessment | No source files found in project | Write ASSESS.md with zeroed summary. Note: "No source files detected in configured directories" |
  | PartialData | Some analysis phases fail | Write ASSESS.md with available data. Mark failed sections as "Analysis unavailable: {reason}" |

Invariants:
  - ASSESS.md follows the exact schema from DESIGN.md section 4.1
  - Every section is present even if empty (no omitted sections)
  - Generated date is included at the top
  - Project name and stack come from config.json
  - File paths in tables are relative to project root
  - Severity classifications are consistent: CRITICAL, HIGH, MEDIUM, LOW
  - Module test status classifications are consistent: tested (>50%), partial (1-50%), untested (0%) per FR-2.7
  - Boundary contract status classifications are consistent: explicit, implicit, none per FR-3.2

Dependencies: C-17 (agent must run to produce this output)
```

### C-20: AssessSecuritySpawn

```
Contract: AssessSecuritySpawn
Boundary: /gl:assess orchestrator -> gl-security agent
Slice: S-08 (Codebase Assessment)
Design refs: FR-4.1, FR-4.6, NFR-8, DESIGN.md 6.1

SECURITY AGENT MODE: full-audit

Input (context passed to gl-security):
  - Mode: "full-audit"
  - Scope: entire codebase
  - Source code files
  - CLAUDE.md security standards

Output:
  - Security findings list with: severity, category, file location, description
  - Findings are returned to orchestrator/assessor for inclusion in ASSESS.md

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | SecurityAgentFailed | Agent spawn fails or returns error | Command continues without security findings. Warn user. ASSESS.md notes "Security scan not performed" (NFR-8) |
  | NoFindings | Agent completes but finds no issues | Report "No security vulnerabilities detected" in ASSESS.md |

Invariants:
  - Security agent in assess is DOCUMENT ONLY -- no failing tests written (DESIGN.md 6.1)
  - Security agent uses same vulnerability checklist as /gl:slice and /gl:ship
  - Security findings are merged into ASSESS.md under "Security Findings" section (FR-4.6)
  - Agent failure does not block assessment completion (NFR-8)

Dependencies: C-17 (assessment context), C-18 (orchestrator spawns security agent)
```

---

## S-09: Boundary Wrapping (/gl:wrap)

*User Action: "User can wrap existing code in contracts and locking tests without rewriting it"*

### C-21: WrapperAgentBehaviour

```
Contract: WrapperAgentBehaviour
Boundary: /gl:wrap orchestrator -> gl-wrapper agent (markdown prompt file)
Slice: S-09 (Boundary Wrapping)
Design refs: FR-9, FR-10, FR-11, NFR-2, NFR-5, TD-2, TD-6, TD-7, TD-8, TD-10

AGENT DEFINITION: src/agents/gl-wrapper.md

Input (context passed via Task spawn):
  - Selected boundary name and files
  - .greenlight/config.json (test commands, stack)
  - .greenlight/codebase/ docs (codebase understanding)
  - .greenlight/CONTRACTS.md (existing contracts -- do not duplicate)
  - CLAUDE.md engineering standards
  - Implementation source code for the selected boundary (DELIBERATE EXCEPTION per TD-2)

Output (files written directly by agent):
  - Extracted contracts presented to user for confirmation
  - Approved contracts written to .greenlight/CONTRACTS.md with [WRAPPED] tag
  - Locking tests written to tests/locking/{boundary-name}.test.{ext}

Tools: Read, Write, Bash, Glob, Grep
Model: Resolved from config.json profiles (default: sonnet in balanced profile)

ISOLATION EXCEPTION (TD-2):
  gl-wrapper is a deliberate exception to agent isolation. It sees implementation code
  AND writes locking tests. This is necessary because locking tests must verify what
  code currently does, not what contracts say it should do. Exception scope: ONLY
  applies to tests in tests/locking/. When a boundary is later refactored via /gl:slice,
  locking tests are deleted and proper integration tests are written under strict isolation.

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | BoundaryTooLarge | Selected boundary exceeds 50% context budget | Suggest splitting into sub-boundaries. Present sub-boundary options to user for selection. (TD-8, NFR-2) |
  | ContractRejected | User rejects extracted contracts during review | Allow user to edit or re-extract. Do not write rejected contracts. (FR-9.3) |
  | LockingTestFailure | Locking test fails against existing code | Fix the test (not the code). Re-run. Max 3 fix-rerun cycles. (FR-11.2, FR-11.3, FR-11.4) |
  | NonDeterministicBehaviour | Test fails due to timestamps, random IDs, env values | Auto-handle: freeze time, use matchers, mock env. After 2 failed attempts, document in ASSESS.md as non-deterministic and skip. (TD-7, FR-10.6, FR-10.7) |
  | ExistingContracts | CONTRACTS.md already has contracts for this boundary | Do not duplicate. Warn user. Ask whether to overwrite or skip. (FR-7.5) |
  | MaxRetriesExceeded | 3 test-fix-rerun cycles exhausted | Escalate to user with failure details. Do not commit partial results. (FR-11.4) |

Invariants:
  - MUST NEVER modify production source code (NFR-5)
  - Only writes: contracts (CONTRACTS.md) and locking tests (tests/locking/)
  - Extracted contracts are DESCRIPTIVE (what code DOES), not prescriptive (what it SHOULD do) (FR-9.5)
  - Contracts follow gl-architect.md format: input, output, errors, invariants, security, dependencies (FR-9.7)
  - [WRAPPED] tag includes Source, Wrapped on date, and Locking tests path (FR-9.6)
  - Locking tests MUST pass against existing code without any source code changes (FR-10.2)
  - Locking tests go in tests/locking/{boundary-name}.test.{ext} -- one file per boundary (TD-6, FR-10.3)
  - Test names use [LOCK] prefix: "[LOCK] should return user object when valid email provided" (FR-10.8)
  - Tests cover both happy paths and observable error paths (FR-10.4)
  - Tests do NOT test implementation details -- test observable behaviour at the boundary (FR-10.5)
  - Agent MUST complete one boundary wrap within 50% context window (NFR-2)
  - Only one boundary wrapped per invocation (FR-8.5)

Security:
  - Agent CANNOT modify production source code
  - Isolation exception is scoped to tests/locking/ only
  - Per FR-20.2: gl-wrapper row in CLAUDE.md isolation table

Dependencies: None (wrap works independently)
```

### C-22: WrapOrchestration

```
Contract: WrapOrchestration
Boundary: User -> /gl:wrap command (markdown prompt file)
Slice: S-09 (Boundary Wrapping)
Design refs: FR-7, FR-8, FR-11, FR-12, FR-13, FR-14, NFR-4, NFR-5, DESIGN.md 5.2

COMMAND DEFINITION: src/commands/gl/wrap.md

Input:
  - User invokes /gl:wrap (no arguments -- interactive boundary selection)
  - .greenlight/config.json MUST exist

Output:
  - .greenlight/CONTRACTS.md updated with [WRAPPED] contracts
  - tests/locking/{boundary-name}.test.{ext} created
  - .greenlight/STATE.md Wrapped Boundaries section updated
  - .greenlight/ASSESS.md updated with known issues (if file exists)
  - Commit with message: "test(wrap): lock {boundary-name}" (FR-13.2)
  - Summary report displayed to user
  - Next action recommendation (FR-14)

Orchestration steps:
  1. Read .greenlight/config.json for project context and model resolution
  2. Read .greenlight/ASSESS.md if it exists (for prioritized boundary list) (FR-7.1)
  3. Read .greenlight/CONTRACTS.md if it exists (to avoid duplicates) (FR-7.5)
  4. Read .greenlight/codebase/ docs if they exist (FR-7.4)
  5. Present boundary candidates to user with priority tiers (FR-8.1, FR-8.2)
  6. User picks boundary to wrap (FR-8.3)
  7. Show estimated complexity: file count, function count, dependency count (FR-8.4)
  8. Resolve gl-wrapper model from config.json profiles
  9. Spawn gl-wrapper agent via Task with selected boundary context
  10. Agent reads implementation, extracts contracts, presents for user confirmation (FR-9.3)
  11. Agent writes locking tests, runs them (FR-11.1)
  12. Agent fixes failing tests (max 3 cycles) (FR-11.3, FR-11.4)
  13. Run full test suite to ensure no regressions (FR-11.6)
  14. Spawn gl-security agent in "slice" mode scoped to boundary files (FR-12.1)
  15. Document security findings (do not write failing tests) (FR-12.2, FR-12.4)
  16. Commit locking tests and contracts atomically (FR-13.1)
  17. Update STATE.md Wrapped Boundaries section (FR-13.4)
  18. Display wrap progress and suggest next action (FR-14)

User interaction flow (from DESIGN.md 5.2):
  1. Present boundary candidates with priority (from ASSESS.md or fresh scan)
  2. User picks boundary to wrap
  3. Wrapper extracts contracts, presents for user confirmation
  4. Wrapper writes locking tests, runs them
  5. Report results, suggest next action

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | NoConfigJson | .greenlight/config.json does not exist | Print "No config found. Run /gl:init first." and stop |
  | AlreadyWrapped | Selected boundary already has [WRAPPED] contracts | Warn user. Ask: overwrite existing locking tests? (NFR-4) |
  | WrapperFailure | gl-wrapper agent fails or produces no output | Report error. Do not commit. Suggest retrying |
  | AllTestsFailing | Locking tests fail after 3 fix cycles | Escalate to user with failure details. Do not commit |
  | FullSuiteRegression | Existing tests fail after locking tests added | Report regression. Do not commit. Ask user to investigate |
  | SecurityFailure | gl-security agent fails | Continue without security findings. Warn user. (NFR-8) |

Invariants:
  - /gl:wrap NEVER modifies production source code (NFR-5)
  - Only one boundary wrapped per invocation (FR-8.5)
  - User MUST confirm extracted contracts before they are written (FR-9.3)
  - Locking tests and contracts are committed atomically (FR-13.1)
  - Commit format: "test(wrap): lock {boundary-name}" with body listing counts (FR-13.2, FR-13.3)
  - Works without ASSESS.md -- user can choose what to wrap manually (FR-7.2)
  - Works with any stack Greenlight supports (NFR-7)

Security:
  - Command does not modify production code
  - Security agent runs in slice mode, document-only (no failing tests) (DESIGN.md 6.2)
  - Security issues recorded as known issues in STATE.md (FR-12.3)

Dependencies: None (wrap is always available after /gl:init)
```

### C-23: WrapContractExtraction

```
Contract: WrapContractExtraction
Boundary: gl-wrapper agent -> Filesystem (.greenlight/CONTRACTS.md with [WRAPPED] tag)
Slice: S-09 (Boundary Wrapping)
Design refs: FR-9, DESIGN.md 4.2

FILE SPECIFICATION: [WRAPPED] contract entries in .greenlight/CONTRACTS.md

Output format per wrapped contract:
  ### Contract: {BoundaryName} [WRAPPED]

  **Source:** `{file}:{start_line}-{end_line}`
  **Wrapped on:** {YYYY-MM-DD}
  **Locking tests:** `tests/locking/{boundary-name}.test.{ext}`

  **Boundary:** {what talks to what}
  **Slice:** wrappable (available for refactoring via /gl:slice with wraps field)

  **Input:** {inferred input type/interface in source language}
  **Output:** {inferred output type/interface in source language}
  **Errors:** table of observed error conditions
  **Invariants:** observed invariants from existing code behaviour
  **Security:** known issues from security baseline, or "none identified"
  **Dependencies:** other contracts this uses

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | DuplicateContract | Contract name already exists in CONTRACTS.md (non-wrapped) | Append boundary suffix to name. Warn user |
  | DuplicateWrapped | [WRAPPED] contract already exists for this boundary | Overwrite if user confirmed (NFR-4). Preserve original Wrapped on date |
  | ParseFailure | Cannot infer types from implementation code | Write contract with best-effort types. Add note: "Types inferred from runtime observation" |

Invariants:
  - Contracts are appended to existing CONTRACTS.md, never overwrite stabilisation contracts (C-01 through C-16)
  - [WRAPPED] tag is always present on wrapped contracts
  - Source field records exact file path and line range (FR-9.6)
  - Wrapped on date is the date of wrapping (FR-9.6)
  - Locking tests field points to the actual test file path (FR-9.6)
  - Contracts are descriptive (what code DOES), not prescriptive (FR-9.5)
  - Contract format follows gl-architect.md standard (FR-9.7)
  - User confirms contracts before they are written (FR-9.3)

Dependencies: C-21 (wrapper agent produces contracts), C-22 (orchestrator manages user confirmation)
```

### C-24: WrapLockingTests

```
Contract: WrapLockingTests
Boundary: gl-wrapper agent -> Filesystem (tests/locking/{boundary-name}.test.{ext})
Slice: S-09 (Boundary Wrapping)
Design refs: FR-10, FR-11, TD-6, TD-7

FILE SPECIFICATION: tests/locking/{boundary-name}.test.{ext}

Output:
  - One test file per boundary at tests/locking/{boundary-name}.test.{ext}
  - Extension matches project stack: .go for Go, .test.ts for TS, .test.py for Python, etc.
  - Test names prefixed with [LOCK]: "[LOCK] should return user object when valid email provided"

Test characteristics:
  - Verify EXISTING behaviour (locking tests, not specification tests) (FR-10.1)
  - MUST pass against existing code without any source code changes (FR-10.2)
  - Test both happy paths and observable error paths (FR-10.4)
  - Do NOT test implementation details; test observable behaviour at boundary (FR-10.5)
  - Handle non-determinism automatically: timestamps (freeze time), random IDs (matchers),
    environment-dependent values (mock env) (FR-10.6)

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | TestFailure | Locking test fails against existing code | Fix the TEST (never the code). Re-run. Max 3 cycles. (FR-11.2, FR-11.3, FR-11.4) |
  | NonDeterministic | Test fails due to non-determinism after 2 fix attempts | Document specific assertion in ASSESS.md as non-deterministic. Skip that test with note. (TD-7, FR-10.7) |
  | InfrastructureError | Test runner fails (syntax, import, missing deps) | Fix infrastructure issue (not test logic). Re-run. |
  | FullSuiteRegression | Full test suite fails after locking tests pass | Report regression. Do not commit. Escalate to user. (FR-11.6) |

Invariants:
  - All locking tests MUST pass before commit (FR-11.2)
  - Full test suite MUST pass after locking tests pass (FR-11.6)
  - Maximum 3 test-fix-rerun cycles per boundary (FR-11.4)
  - Test file location: tests/locking/{boundary-name}.test.{ext} (TD-6)
  - One file per boundary (TD-6)
  - Test names use [LOCK] prefix (FR-10.8)
  - Each test is independent -- no shared mutable state between tests
  - Test results reported: total tests, passing, skipped due to non-determinism (FR-11.5)

Dependencies: C-21 (wrapper agent writes tests), C-23 (contracts guide what to test)
```

### C-25: WrapSecurityBaseline

```
Contract: WrapSecurityBaseline
Boundary: /gl:wrap orchestrator -> gl-security agent
Slice: S-09 (Boundary Wrapping)
Design refs: FR-12, NFR-8, DESIGN.md 6.2

SECURITY AGENT MODE: slice (scoped to boundary files, document only)

Input (context passed to gl-security):
  - Mode: "slice" (diff-only, scoped to boundary)
  - Scope: files belonging to the wrapped boundary
  - Source code for boundary files
  - CLAUDE.md security standards

Output:
  - Security findings list with: severity, category, file location, description
  - Findings documented in ASSESS.md (or created if it doesn't exist) (FR-12.2)
  - Findings recorded as known issues in STATE.md Wrapped Boundaries table (FR-12.3)

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | SecurityAgentFailed | Agent spawn fails or returns error | Continue wrap without security findings. Warn user. (NFR-8) |
  | NoFindings | Agent completes but finds no issues | Record "0 known issues" in STATE.md |

Invariants:
  - Security agent MUST NOT write failing tests during wrap (FR-12.4)
  - Security issues are DOCUMENTED, not enforced (DESIGN.md 6.2)
  - Known issues are recorded in STATE.md Wrapped Boundaries table (FR-12.3)
  - Agent failure does not block wrap completion (NFR-8)
  - When boundary is later refactored via /gl:slice, security agent runs in normal mode (DESIGN.md 6.3)

Dependencies: C-21 (boundary must be wrapped first), C-22 (orchestrator spawns security agent)
```

### C-26: WrapStateTracking

```
Contract: WrapStateTracking
Boundary: /gl:wrap orchestrator -> Filesystem (.greenlight/STATE.md)
Slice: S-09 (Boundary Wrapping)
Design refs: FR-13.4, DESIGN.md 4.3

FILE SPECIFICATION: Wrapped Boundaries section in .greenlight/STATE.md

Output (appended section in STATE.md):
  ## Wrapped Boundaries

  | Boundary | Contracts | Locking Tests | Known Issues | Status |
  |----------|-----------|---------------|--------------|--------|
  | {name}   | {N}       | {N}           | {N}          | wrapped |

  Wrap progress: {N}/{M} boundaries wrapped

Status values:
  - wrapped: locking tests in place, contracts extracted, existing behaviour locked
  - refactored: replaced by integration tests via /gl:slice. Locking tests deleted.

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | NoStateFile | STATE.md does not exist | Create STATE.md with Wrapped Boundaries section only |
  | SectionMissing | STATE.md exists but has no Wrapped Boundaries section | Append section to existing STATE.md |
  | BoundaryExists | Boundary already listed in table | Update existing row (overwrite counts, preserve original wrap date) |

Invariants:
  - Wrapped Boundaries section counts toward STATE.md's 80-line budget (DESIGN.md 4.3)
  - Wrap progress denominator (M) comes from ASSESS.md priority list count, or is omitted if no ASSESS.md (FR-14.2)
  - When a boundary status changes to "refactored", it can be compressed to a summary line (DESIGN.md 4.3)
  - Section is placed below existing Slices section

Dependencies: C-22 (orchestrator updates STATE.md after successful wrap)
```

---

## S-10: Brownfield-Aware Command Updates

*User Action: "User can see brownfield progress alongside greenfield slices"*

### C-27: StatusBrownfieldDisplay

```
Contract: StatusBrownfieldDisplay
Boundary: /gl:status command -> User (wrapped boundaries display)
Slice: S-10 (Brownfield-Aware Command Updates)
Design refs: FR-16, DESIGN.md 5.4

COMMAND UPDATE: src/commands/gl/status.md

Input:
  - .greenlight/STATE.md (check for Wrapped Boundaries section)

Output (additional display section, shown if wrapped boundaries exist):
  Wrapped Boundaries:
    {name}                 wrapped   {N} locking tests  {N} known issues
    {name}                 wrapped   {N} locking tests  {N} known issues
    {name}                 refactored (replaced by slice {N})

  Wrap: {N}/{M} boundaries wrapped, {R} refactored

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | NoWrappedBoundaries | STATE.md has no Wrapped Boundaries section | Do not display section. Standard status output only |
  | EmptyTable | Wrapped Boundaries section exists but table is empty | Do not display section |

Invariants:
  - Wrapped Boundaries table is ONLY displayed if any wrapped boundaries exist in STATE.md (FR-16.1)
  - Table columns: Boundary, Contracts, Locking Tests (count), Known Issues (count), Status (FR-16.2)
  - Status values: wrapped, refactored (FR-16.3)
  - Section appears after the slice status table, before the Next recommendation
  - Does not break existing status display for greenfield-only projects

Dependencies: C-26 (STATE.md must have Wrapped Boundaries section to display)
```

### C-28: HelpBrownfieldSection

```
Contract: HelpBrownfieldSection
Boundary: /gl:help command -> User (brownfield commands in help output)
Slice: S-10 (Brownfield-Aware Command Updates)
Design refs: FR-17

COMMAND UPDATE: src/commands/gl/help.md

Output (updated help display):
  - BROWNFIELD section inserted between SETUP and BUILD sections (FR-17.1)
  - Commands listed: /gl:assess (Gap analysis + risk assessment), /gl:wrap (Extract contracts + locking tests) (FR-17.2)
  - FLOW line updated to: map? -> assess? -> init -> design -> wrap? -> slice 1 -> ... -> ship (FR-17.3)

Exact output (from DESIGN.md 5.3):
  BROWNFIELD
    /gl:assess            Gap analysis + risk assessment -> ASSESS.md
    /gl:wrap              Extract contracts + locking tests

  FLOW
    map? -> assess? -> init -> design -> wrap? -> slice 1 -> ... -> ship

Errors: None (help always succeeds)

Invariants:
  - BROWNFIELD section is always present in help output (not conditional)
  - BROWNFIELD appears between SETUP and BUILD sections (FR-17.1)
  - FLOW line includes assess? and wrap? as optional steps (FR-17.3)
  - Does not remove or modify any existing help sections

Dependencies: None
```

### C-29: SettingsBrownfieldAgents

```
Contract: SettingsBrownfieldAgents
Boundary: /gl:settings command -> User (assessor and wrapper agent display)
Slice: S-10 (Brownfield-Aware Command Updates)
Design refs: FR-18

COMMAND UPDATE: src/commands/gl/settings.md

Output (updated settings display):
  - assessor and wrapper agent models shown in settings table (FR-18.1)
  - Valid agents list includes assessor and wrapper (FR-18.2)

Updated MODELS section in settings display:
  architect        {model}    ({source})
  designer         {model}    ({source})
  test_writer      {model}    ({source})
  implementer      {model}    ({source})
  security         {model}    ({source})
  debugger         {model}    ({source})
  verifier         {model}    ({source})
  codebase_mapper  {model}    ({source})
  assessor         {model}    ({source})     <-- NEW
  wrapper          {model}    ({source})     <-- NEW

Updated valid agents list:
  architect, designer, test_writer, implementer, security, debugger,
  verifier, codebase_mapper, assessor, wrapper

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | MissingProfileEntries | Config.json profiles lack assessor/wrapper entries | Fall back to sonnet for assessor, sonnet for wrapper (same as other missing agents) |

Invariants:
  - assessor and wrapper appear in the agent model display table (FR-18.1)
  - assessor and wrapper are accepted as valid agent names for model override commands (FR-18.2)
  - Model resolution follows same chain: override > profile > fallback (sonnet)
  - Does not break existing settings display for agents without brownfield config

Dependencies: C-34 (config.json profiles must include assessor and wrapper entries)
```

---

## S-11: Locking-to-Integration Transition

*User Action: "User can refactor a wrapped boundary through the normal TDD loop"*

### C-30: SliceWrapsField

```
Contract: SliceWrapsField
Boundary: /gl:slice command -> Locking-to-integration transition
Slice: S-11 (Locking-to-Integration Transition)
Design refs: FR-15, DESIGN.md 4.4, DESIGN.md 6.3

COMMAND UPDATE: src/commands/gl/slice.md

Input (additional context when slice has wraps field):
  - GRAPH.json slice object with optional "wraps" field (array of boundary names)
  - STATE.md Wrapped Boundaries section
  - Existing locking tests in tests/locking/{boundary-name}.test.{ext}
  - [WRAPPED] contracts in CONTRACTS.md

Behaviour when slice has wraps field:
  1. Pre-flight: Read locking test file names for wrapped boundaries (FR-15.2)
  2. Extract locking test NAMES (not source code) for context (FR-15.3)
  3. Pass locking test names to test writer as "existing locked behaviours" context
  4. Normal TDD loop proceeds (test writer, implementer, security, verifier)
  5. After verification succeeds: confirm both locking tests AND integration tests pass
  6. Delete locking tests from tests/locking/ for the wrapped boundary (FR-15.4)
  7. Remove [WRAPPED] tag from corresponding contracts in CONTRACTS.md (FR-15.5)
  8. Update STATE.md Wrapped Boundaries status to "refactored" (FR-15.6)

Behaviour when slice has NO wraps field:
  - No change from existing /gl:slice behaviour

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | MissingLockingTests | wraps field references boundary but tests/locking/ file not found | Warn user. Proceed without locking test context. Treat as greenfield slice |
  | MissingWrappedContract | wraps field references boundary but no [WRAPPED] contract found | Warn user. Proceed without wrapped contract context |
  | LockingTestsStillFailing | After implementation, locking tests fail | Regression: existing behaviour broken. Implementer must fix. Locking tests serve as guardrail during refactoring |
  | TransitionFailure | Integration tests pass but locking test deletion or tag removal fails | Report error. Manual cleanup may be needed. Do not mark as refactored |

Invariants:
  - wraps field is optional on GRAPH.json slice objects (FR-15.1, TD-12)
  - wraps field is an array of boundary names matching STATE.md Wrapped Boundaries table entries (DESIGN.md 4.4)
  - Test writer receives locking test NAMES only (not source code) (FR-15.3, FR-22.2)
  - Integration tests MUST be a superset of locked behaviours (FR-22.3)
  - Locking tests are deleted ONLY after verification succeeds (FR-15.4)
  - [WRAPPED] tag removal happens ONLY after locking tests are deleted (FR-15.5)
  - STATE.md boundary status changes to "refactored" ONLY after full transition (FR-15.6)
  - A slice can wrap multiple boundaries if closely related (DESIGN.md 4.4)
  - During refactoring, BOTH locking tests AND integration tests must pass (DESIGN.md 6.3)
  - Security agent checks whether known issues from wrapped contract have been addressed (DESIGN.md 6.3)
  - Pre-existing security issues that persist do not block the slice (DESIGN.md 6.3)
  - The wraps field does NOT create a dependency on the boundary being wrapped first (DESIGN.md 4.4)

Security:
  - Security agent runs in normal mode during refactoring (writes failing tests for NEW vulnerabilities)
  - Known issues from wrapped contract are checked but do not block if pre-existing (DESIGN.md 6.3)

Dependencies: C-21 (wrapped boundary must exist), C-24 (locking tests must exist), C-32 (test writer needs locking test awareness)
```

### C-31: ArchitectWrappedContracts

```
Contract: ArchitectWrappedContracts
Boundary: gl-architect agent -> CONTRACTS.md ([WRAPPED] contract awareness)
Slice: S-11 (Locking-to-Integration Transition)
Design refs: FR-21

AGENT UPDATE: src/agents/gl-architect.md

Behaviour when [WRAPPED] contracts exist in CONTRACTS.md:
  - Recognise [WRAPPED] contracts as existing boundaries (FR-21.1)
  - Do NOT redefine wrapped contracts; treat them as given (FR-21.2)
  - When adding new slices, CAN reference wrapped contracts as dependencies (FR-21.3)
  - When a slice's wraps field targets a wrapped boundary, plan the contract transition:
    wrapped contract becomes a proper contract (tag removed) (FR-21.4)
  - In GRAPH.json, CAN add wraps field to new slices targeting wrapped boundaries

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | WrappedContractConflict | New greenfield contract has same name as existing [WRAPPED] contract | Use the wraps field to plan a transition slice. Do not create duplicate |

Invariants:
  - Architect NEVER redefines a [WRAPPED] contract (FR-21.2)
  - [WRAPPED] contracts are treated as existing, immutable boundaries
  - New slices CAN depend on wrapped contracts
  - wraps field in GRAPH.json is the mechanism for planning boundary refactoring (TD-12)
  - Contract transition lifecycle: [WRAPPED] -> wraps field in slice -> refactored (tag removed)

Dependencies: C-23 (wrapped contracts exist in CONTRACTS.md)
```

### C-32: TestWriterLockingAwareness

```
Contract: TestWriterLockingAwareness
Boundary: gl-test-writer agent -> Locking test name context
Slice: S-11 (Locking-to-Integration Transition)
Design refs: FR-22

AGENT UPDATE: src/agents/gl-test-writer.md

Behaviour when writing tests for a slice that wraps a boundary:
  - Check for existing locking tests in tests/locking/ (FR-22.1)
  - Receive locking test NAMES (not source code) as context (FR-22.2)
  - Use test names to understand what behaviours are already locked
  - Integration tests MUST cover at least all behaviours that locking tests covered (superset) (FR-22.3)
  - Integration tests go in tests/integration/ as normal (NOT in tests/locking/)

Behaviour when slice does NOT wrap a boundary:
  - No change from existing gl-test-writer behaviour

Input (additional context for wraps slices):
  - Locking test names/descriptions extracted from test file
  - Format: list of "[LOCK] {description}" strings
  - NOT the test source code (FR-22.4)

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | NoLockingTests | Slice wraps a boundary but no locking tests found | Warn orchestrator. Write integration tests from contracts only (standard behaviour) |
  | IncompleteCoverage | Integration tests do not cover all locked behaviours | Orchestrator flags during verification. Test writer may be respawned |

Invariants:
  - Test writer MUST NOT see locking test implementation (only names/descriptions) (FR-22.4)
  - Integration tests are a superset of locked behaviours (FR-22.3)
  - Standard agent isolation applies: test writer cannot see implementation code
  - Locking test names are informational context, not a test specification
  - New integration tests follow the same patterns as greenfield slices

Dependencies: C-30 (slice orchestrator provides locking test names as context)
```

---

## S-12: Infrastructure and Config Updates

*Supports all user actions (infrastructure enabling layer)*

### C-33: ManifestBrownfieldUpdate

```go
// Contract: ManifestBrownfieldUpdate
// Boundary: Go CLI -> Manifest (4 new file paths)
// Slice: S-12 (Infrastructure and Config Updates)
// Design refs: DESIGN.md 7.3, UD-9, UD-11
//
// FILE: internal/installer/installer.go
//
// Change: Add 4 new entries to Manifest slice
//
// New entries (inserted in alphabetical order within their section):
//   "agents/gl-assessor.md"     // NEW -- brownfield assessment agent
//   "agents/gl-wrapper.md"      // NEW -- brownfield wrapping agent
//   "commands/gl/assess.md"     // NEW -- /gl:assess command
//   "commands/gl/wrap.md"       // NEW -- /gl:wrap command
//
// Updated Manifest (30 entries, up from 26):
//   "agents/gl-architect.md"
//   "agents/gl-assessor.md"      <-- NEW
//   "agents/gl-codebase-mapper.md"
//   "agents/gl-debugger.md"
//   "agents/gl-designer.md"
//   "agents/gl-implementer.md"
//   "agents/gl-security.md"
//   "agents/gl-test-writer.md"
//   "agents/gl-verifier.md"
//   "agents/gl-wrapper.md"       <-- NEW
//   "commands/gl/add-slice.md"
//   "commands/gl/assess.md"      <-- NEW
//   "commands/gl/design.md"
//   "commands/gl/help.md"
//   "commands/gl/init.md"
//   "commands/gl/map.md"
//   "commands/gl/pause.md"
//   "commands/gl/quick.md"
//   "commands/gl/resume.md"
//   "commands/gl/settings.md"
//   "commands/gl/ship.md"
//   "commands/gl/slice.md"
//   "commands/gl/status.md"
//   "commands/gl/wrap.md"        <-- NEW
//   "references/checkpoint-protocol.md"
//   "references/deviation-rules.md"
//   "references/verification-patterns.md"
//   "templates/config.md"
//   "templates/state.md"
//   "CLAUDE.md"
//
// Errors: none (compile-time constant)
//
// Invariants:
// - CLAUDE.md remains the LAST entry
// - Entries within each section (agents/, commands/gl/, etc.) are alphabetically ordered
// - go:embed directive in main.go already uses wildcards (src/agents/*.md, src/commands/gl/*.md)
//   so new .md files in those directories are automatically embedded -- no main.go change needed
// - Manifest count increases from 26 to 30
// - All existing tests that validate manifest count must be updated to expect 30
//
// Dependencies: None (this is a Go code change, independent of markdown files)
```

### C-34: ConfigProfileBrownfield

```
Contract: ConfigProfileBrownfield
Boundary: config.json template -> Profile (assessor and wrapper agent entries)
Slice: S-12 (Infrastructure and Config Updates)
Design refs: DESIGN.md 4.5, TD-9, TD-10

FILE UPDATES:
  - src/templates/config.md (profile schema documentation)
  - src/commands/gl/init.md (default config.json generation)

Profile additions:
  quality:
    assessor: "opus"
    wrapper: "opus"
  balanced:
    assessor: "sonnet"
    wrapper: "sonnet"
  budget:
    assessor: "haiku"
    wrapper: "sonnet"

Agent names list update:
  From: architect, designer, test_writer, implementer, security, debugger, verifier, codebase_mapper
  To:   architect, designer, test_writer, implementer, security, debugger, verifier, codebase_mapper, assessor, wrapper

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | MissingEntries | Existing config.json lacks assessor/wrapper in profiles | Model resolution falls back to sonnet (same as any missing agent) |

Invariants:
  - gl-assessor defaults to sonnet in balanced profile (TD-9)
  - gl-wrapper defaults to sonnet in balanced profile (TD-10)
  - Model resolution chain unchanged: override > profile > sonnet fallback
  - Existing config.json files without assessor/wrapper entries still work (graceful degradation)
  - No Go CLI code changes needed for config.json -- it's read by Claude Code agents, not parsed by Go

Dependencies: None
```

### C-35: CLAUDEmdIsolationUpdate

```
Contract: CLAUDEmdIsolationUpdate
Boundary: CLAUDE.md -> Agent Isolation Rules (gl-assessor and gl-wrapper rows)
Slice: S-12 (Infrastructure and Config Updates)
Design refs: FR-20

FILE UPDATE: src/CLAUDE.md

Agent Isolation Rules table additions:
  | gl-assessor | codebase docs, test results, standards | N/A (read-only analytical agent) | modify any code |
  | gl-wrapper  | implementation code, existing tests | N/A | modify production code (only writes contracts and locking tests) |

Exception note (added below Agent Isolation Rules table, clearly marked):
  "gl-wrapper is a deliberate exception. It sees implementation code AND writes locking
  tests. This is necessary because locking tests must verify what code currently does,
  not what contracts say it should do. This exception is scoped: only applies to tests
  in tests/locking/. When a boundary is later refactored via /gl:slice, locking tests
  are deleted and proper integration tests are written under strict isolation." (FR-20.3)

Errors: None (static content update)

Invariants:
  - gl-assessor row present in isolation table (FR-20.1)
  - gl-wrapper row present in isolation table (FR-20.2)
  - Exception note is clearly marked and scoped (FR-20.3)
  - Existing agent isolation rules unchanged
  - Exception note explains: why (locking tests need code), scope (tests/locking/ only),
    lifecycle (deleted when refactored)

Dependencies: None
```

---

## Cross-Cutting: Agent Behaviour Rules for [WRAPPED] Contracts

```
Contract: WrappedContractAgentBehaviours
Boundary: All agents -> [WRAPPED] contracts in CONTRACTS.md
Not a separate slice -- referenced by S-08, S-09, S-11
Design refs: DESIGN.md 4.2 (agent behaviour rules table)

Agent-specific behaviours with [WRAPPED] contracts:

  | Agent | Behaviour with [WRAPPED] contracts |
  |-------|-----------------------------------|
  | gl-architect | Do NOT redefine. Reference as existing. Can add wraps field to new slices (C-31) |
  | gl-test-writer | Check for locking tests. When slice has wraps field, receive locking test NAMES as context. Integration tests must be superset (C-32) |
  | gl-implementer | Build on top of existing code. Use wrapped contract as constraint. Existing code is starting point for refactoring |
  | gl-security | Note known issues from wrapped contract. Check if issues persist after refactoring. Pre-existing issues do not block |
  | gl-verifier | Verify locking tests are removed after successful refactoring. Verify [WRAPPED] tag is removed. Verify STATE.md boundary status updated |
  | gl-assessor | Read-only. Can reference wrapped boundaries in gap analysis |
  | gl-wrapper | Creates [WRAPPED] contracts. Does not interact with existing [WRAPPED] contracts for different boundaries |

Lifecycle:
  [WRAPPED] contract created by /gl:wrap
    -> Slice with wraps field targets this boundary (/gl:slice)
    -> Test writer receives locking test names as context
    -> Integration tests written (superset of locked behaviours)
    -> Implementer refactors code, making integration tests pass
    -> Verification: both locking tests AND integration tests pass
    -> Locking tests deleted from tests/locking/
    -> [WRAPPED] tag removed from contract
    -> STATE.md boundary status -> refactored

Invariants:
  - A [WRAPPED] contract is never modified directly by any agent except during transition
  - The transition from [WRAPPED] to regular contract only happens through /gl:slice with wraps field
  - All agents can read [WRAPPED] contracts; only the transition process can remove the tag
```

---

## User Action Mapping

| User Action | Slice(s) | Contracts | Enabled By |
|-------------|----------|-----------|------------|
| 1. User can assess an existing codebase for gaps, risks, and untested boundaries | S-08 | C-17, C-18, C-19, C-20 | /gl:assess command + gl-assessor agent |
| 2. User can wrap existing code in contracts and locking tests without rewriting it | S-09 | C-21, C-22, C-23, C-24, C-25, C-26 | /gl:wrap command + gl-wrapper agent |
| 3. User can see brownfield progress alongside greenfield slices | S-10 | C-27, C-28, C-29 | /gl:status, /gl:help, /gl:settings updates |
| 4. User can refactor a wrapped boundary through the normal TDD loop | S-11 | C-30, C-31, C-32 | /gl:slice wraps field + architect/test-writer updates |
