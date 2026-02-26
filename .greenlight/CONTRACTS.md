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
| C-62 | ContractSchemaExtension | gl-architect.md -> Contract format (verification fields) | S-22 |
| C-63 | VerifierTierAwareness | gl-verifier.md -> Verification report (tier reporting) | S-22 |
| C-64 | VerificationTiersProtocol | Reference doc -> Orchestrator + agents (verification tier rules) | S-23 |
| C-65 | VerificationTierGate | /gl:slice orchestrator -> Verification tier gate (Step 6b) | S-23 |
| C-66 | VerifyCheckpointPresentation | /gl:slice orchestrator -> User (acceptance checkpoint) | S-23 |
| C-67 | RejectionClassification | /gl:slice orchestrator -> User (gap classification UX) | S-24 |
| C-68 | RejectionToTestWriter | /gl:slice orchestrator -> gl-test-writer (rejection feedback) | S-24 |
| C-69 | RejectionToContractRevision | /gl:slice orchestrator -> User (contract revision route) | S-24 |
| C-70 | RejectionCounter | /gl:slice orchestrator -> Per-slice rejection tracking | S-25 |
| C-71 | RejectionEscalation | /gl:slice orchestrator -> User (escalation at 3 rejections) | S-25 |
| C-72 | CLAUDEmdVerificationTierRule | CLAUDE.md -> All agents (verification tier hard rule) | S-26 |
| C-73 | CheckpointProtocolAcceptanceType | checkpoint-protocol.md -> Acceptance checkpoint type | S-26 |
| C-74 | ManifestVerificationTiersUpdate | Go CLI -> Manifest (1 new file path) | S-26 |
| C-75 | ArchitectTierGuidance | gl-architect.md -> Tier selection guidance and acceptance criteria generation | S-27 |
| C-91 | FrontmatterParse | Frontmatter parser -> Slice state files (parse flat YAML frontmatter) | S-35 |
| C-92 | FrontmatterWrite | Frontmatter writer -> Slice state files (write flat YAML frontmatter) | S-35 |
| C-93 | StateReadSlices | State reader -> Filesystem (read all slice frontmatter into structs) | S-36 |
| C-94 | StateReadGraph | State reader -> Filesystem (read GRAPH.json dependency data) | S-36 |
| C-95 | StateFindReadySlices | State reader -> Ready slice computation (pending + deps complete) | S-36 |
| C-96 | StateDetectContext | State reader -> Environment (detect shell vs Claude context) | S-36 |
| C-97 | CLIDispatchExtension | User -> CLI dispatcher (new subcommands) | S-37 |
| C-98 | RunStatus | CLI -> Command handler (status with progress display) | S-38 |
| C-99 | RunStatusCompact | CLI -> Command handler (compact one-liner for tmux status bar) | S-38 |
| C-100 | RunHelp | CLI -> Command handler (context-aware help listing) | S-39 |
| C-101 | RunRoadmap | CLI -> Command handler (display ROADMAP.md) | S-40 |
| C-102 | RunChangelog | CLI -> Command handler (display changelog from summaries) | S-40 |
| C-103 | ProcessSpawnClaude | Process spawner -> os/exec (spawn Claude with configurable flags) | S-41 |
| C-104 | ProcessSpawnInteractive | Process spawner -> os/exec (launch interactive Claude session) | S-41 |
| C-105 | RunSliceSingle | CLI -> Command handler (run single slice headlessly) | S-42 |
| C-106 | RunSliceAutoDetect | CLI -> Command handler (auto-detect ready slices, run one) | S-42 |
| C-107 | TmuxIsAvailable | tmux manager -> os/exec (check tmux availability) | S-43 |
| C-108 | TmuxNewSession | tmux manager -> os/exec (create named tmux session) | S-43 |
| C-109 | TmuxAddWindow | tmux manager -> os/exec (add window to existing session) | S-43 |
| C-110 | TmuxAttachSession | tmux manager -> os/exec (attach to existing session) | S-43 |
| C-111 | RunSliceParallel | CLI -> Command handler (parallel slice execution via tmux) | S-44 |
| C-112 | RunSliceSequentialFallback | CLI -> Command handler (sequential fallback when no tmux) | S-44 |
| C-113 | RunSliceWatch | CLI -> Command handler (watch mode poll loop) | S-45 |
| C-114 | RunSliceDryRun | CLI -> Command handler (preview mode without execution) | S-45 |
| C-115 | RunInit | CLI -> Command handler (launch interactive init session) | S-46 |
| C-116 | RunDesign | CLI -> Command handler (launch interactive design session) | S-46 |

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

---

# Documentation & Roadmap Contracts

> **Scope:** `/gl:roadmap`, `/gl:changelog`, auto-summaries, decision log, living architecture diagram
> **Deliverables:** Markdown prompt files (commands + updates to existing commands/agents), NOT Go code
> **Date:** 2026-02-09
> **Design Reference:** DESIGN.md sections 1.4-1.6, Technical Decisions TD-13 through TD-21, FRs 23-30

---

## Documentation Contract Index

| # | Contract | Boundary | Slice |
|---|----------|----------|-------|
| C-36 | DesignRoadmapProduction | /gl:design orchestrator -> Filesystem (ROADMAP.md) | S-13 |
| C-37 | DesignDecisionsSeeding | /gl:design orchestrator -> Filesystem (DECISIONS.md) | S-13 |
| C-38 | ManifestDocumentationUpdate | Go CLI -> Manifest (2 new file paths) | S-13 |
| C-39 | HelpInsightSection | /gl:help command -> User (INSIGHT section) | S-13 |
| C-40 | StatusDocumentationReference | /gl:status command -> User (roadmap/changelog reference) | S-13 |
| C-41 | SliceSummaryGeneration | /gl:slice orchestrator -> Summary Task -> Filesystem | S-14 |
| C-42 | WrapSummaryGeneration | /gl:wrap orchestrator -> Summary Task -> Filesystem | S-14 |
| C-43 | QuickSummaryGeneration | /gl:quick orchestrator -> Summary Task -> Filesystem | S-14 |
| C-44 | DecisionAggregation | /gl:slice orchestrator -> Filesystem (DECISIONS.md) | S-14 |
| C-45 | RoadmapAutoUpdate | Orchestrators (slice/wrap) -> Filesystem (ROADMAP.md) | S-14 |
| C-46 | RoadmapDisplay | User -> /gl:roadmap command (display) | S-15 |
| C-47 | RoadmapMilestonePlanning | User -> /gl:roadmap milestone -> gl-designer | S-15 |
| C-48 | RoadmapMilestoneArchive | User -> /gl:roadmap archive -> Filesystem (ROADMAP.md) | S-15 |
| C-49 | ChangelogDisplay | User -> /gl:changelog command (display) | S-16 |
| C-50 | ChangelogFiltering | /gl:changelog command -> Summary filtering (milestone, date) | S-16 |
| C-51 | BrownfieldDesignContext | /gl:design orchestrator -> gl-designer (brownfield context blocks) | S-17 |
| C-52 | BrownfieldRoadmapContext | /gl:roadmap milestone -> gl-designer (assessment + wrap progress) | S-17 |
| C-53 | DesignerBrownfieldAwareness | gl-designer agent -> Brownfield-aware design (risk tiers, [WRAPPED] tags) | S-17 |

---

## S-13: Documentation Infrastructure and Design Update

*User Action: "User can start a project with a product roadmap and decision log from day one"*

### C-36: DesignRoadmapProduction

```
Contract: DesignRoadmapProduction
Boundary: /gl:design orchestrator -> Filesystem (.greenlight/ROADMAP.md)
Slice: S-13 (Documentation Infrastructure and Design Update)
Design refs: FR-19.4, FR-23.1, FR-23.2, FR-23.3, FR-23.4, FR-23.6, FR-29.1, FR-29.2, TD-15, TD-19

COMMAND UPDATE: src/commands/gl/design.md

Behaviour (after design approval):
  1. Produce .greenlight/ROADMAP.md containing:
     - Project overview (name, updated date)
     - Architecture Diagram section with Mermaid diagram (FR-29.1, FR-29.2)
     - Initial milestone section with: name, goal, status (active), slice table (FR-23.3, FR-23.4)
     - Wrap progress section (if wrapped boundaries exist in STATE.md) (FR-23.6)
     - Empty Archived Milestones section
  2. Slice table columns: Slice, Description, Status, Tests, Completed, Key Decision (FR-23.4)
  3. All slices from GRAPH.json are listed in the initial milestone's slice table with status=pending
  4. Commit ROADMAP.md as part of design commit

Input (additional context for ROADMAP.md generation):
  - GRAPH.json (slices for milestone table)
  - DESIGN.md (architecture for diagram, project overview)
  - STATE.md (wrapped boundaries, if any)
  - config.json (project name)

Output:
  - .greenlight/ROADMAP.md following schema in DESIGN.md section 4.5

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | NoGraphJson | GRAPH.json does not exist after design approval | Create ROADMAP.md with empty slice table. Note: "No slices defined yet" |
  | DiagramGenerationFailure | Cannot produce Mermaid diagram from design | Create ROADMAP.md with placeholder diagram: "Architecture diagram pending" |

Invariants:
  - ROADMAP.md is produced as part of /gl:design, not as a separate step
  - Architecture diagram uses Mermaid format (TD-19), text-based, no images (FR-29.6)
  - Initial milestone is always named from the design session's scope
  - Every slice in GRAPH.json appears in the milestone slice table
  - ROADMAP.md follows the exact schema in DESIGN.md section 4.5
  - ROADMAP.md is created, never appended to, during /gl:design (fresh start)
  - When invoked via /gl:roadmap milestone (FR-19.6), ROADMAP.md is appended to, not overwritten
  - Wrap progress section is present only if STATE.md has wrapped boundaries

Security:
  - No sensitive data in ROADMAP.md (no credentials, API keys, PII)
  - ROADMAP.md does not include security finding details, only wrap progress counts (DESIGN.md 6.4)

Dependencies: None (design.md update is standalone)
```

### C-37: DesignDecisionsSeeding

```
Contract: DesignDecisionsSeeding
Boundary: /gl:design orchestrator -> Filesystem (.greenlight/DECISIONS.md)
Slice: S-13 (Documentation Infrastructure and Design Update)
Design refs: FR-19.5, FR-28.1, FR-28.2, FR-28.4, FR-28.6, FR-28.7, TD-18

COMMAND UPDATE: src/commands/gl/design.md

Behaviour (after design approval):
  1. Produce .greenlight/DECISIONS.md containing:
     - Header with project name
     - Decision log table with columns: #, Decision, Context, Chosen, Rejected, Date, Source
     - All major technical decisions from the design session as initial entries
  2. Decisions are numbered sequentially: D-1, D-2, D-3... (FR-28.6)
  3. Source column for all initial entries is "design" (FR-28.3)
  4. Commit DECISIONS.md as part of design commit

Input (context for seeding):
  - DESIGN.md technical decisions table (the seed data)
  - Design session outputs (additional decisions made during design)

Output:
  - .greenlight/DECISIONS.md following schema in DESIGN.md section 4.6

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | NoDesignDecisions | Design session produced no explicit decisions | Create DECISIONS.md with header and empty table. Note: "No decisions recorded during design" |

Invariants:
  - DECISIONS.md is seeded from DESIGN.md technical decisions, not duplicated (FR-28.7, TD-18)
  - DECISIONS.md is append-only after creation (NFR-11)
  - Numbering is sequential across the entire file regardless of source (FR-28.6)
  - Source values follow the schema: design, milestone, slice:{id}, quick, wrap:{boundary} (FR-28.3)
  - Date is the date of the design session
  - DECISIONS.md is created, never appended to, during initial /gl:design (fresh start)
  - When invoked via /gl:roadmap milestone (FR-19.6), new decisions are appended to existing DECISIONS.md

Security:
  - No sensitive data in DECISIONS.md
  - Security-related decisions may be recorded (e.g., "chose argon2 over bcrypt") but without exposing implementation details (DESIGN.md 6.4)

Dependencies: None (design.md update is standalone)
```

### C-38: ManifestDocumentationUpdate

```go
// Contract: ManifestDocumentationUpdate
// Boundary: Go CLI -> Manifest (2 new file paths for documentation commands)
// Slice: S-13 (Documentation Infrastructure and Design Update)
// Design refs: DESIGN.md 7.3, UD-10, UD-17
//
// FILE: internal/installer/installer.go
//
// Change: Add 2 new entries to Manifest slice (in addition to 4 brownfield entries from C-33)
//
// New entries (inserted in alphabetical order within commands/gl/ section):
//   "commands/gl/changelog.md"    // NEW -- /gl:changelog command
//   "commands/gl/roadmap.md"      // NEW -- /gl:roadmap command
//
// Updated Manifest (32 entries, up from 30 after brownfield, 26 original):
//   "agents/gl-architect.md"
//   "agents/gl-assessor.md"
//   "agents/gl-codebase-mapper.md"
//   "agents/gl-debugger.md"
//   "agents/gl-designer.md"
//   "agents/gl-implementer.md"
//   "agents/gl-security.md"
//   "agents/gl-test-writer.md"
//   "agents/gl-verifier.md"
//   "agents/gl-wrapper.md"
//   "commands/gl/add-slice.md"
//   "commands/gl/assess.md"
//   "commands/gl/changelog.md"     <-- NEW
//   "commands/gl/design.md"
//   "commands/gl/help.md"
//   "commands/gl/init.md"
//   "commands/gl/map.md"
//   "commands/gl/pause.md"
//   "commands/gl/quick.md"
//   "commands/gl/resume.md"
//   "commands/gl/roadmap.md"       <-- NEW
//   "commands/gl/settings.md"
//   "commands/gl/ship.md"
//   "commands/gl/slice.md"
//   "commands/gl/status.md"
//   "commands/gl/wrap.md"
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
// - Entries within commands/gl/ section are alphabetically ordered
// - go:embed directive uses wildcards -- no main.go change needed
// - Manifest count increases from 30 (post-brownfield) to 32
// - All existing tests that validate manifest count must be updated to expect 32
// - This change is additive to C-33 (brownfield manifest update). Both changes apply.
//
// Dependencies: C-33 (brownfield manifest update must be applied first or simultaneously)
```

### C-39: HelpInsightSection

```
Contract: HelpInsightSection
Boundary: /gl:help command -> User (INSIGHT section and updated FLOW line)
Slice: S-13 (Documentation Infrastructure and Design Update)
Design refs: FR-17.3, FR-17.4, FR-17.5, DESIGN.md 5.5

COMMAND UPDATE: src/commands/gl/help.md

Output (updated help display additions):
  - INSIGHT section inserted between MONITOR and SHIP sections (FR-17.3)
  - Commands listed in INSIGHT:
      /gl:roadmap           Product roadmap + milestones
      /gl:changelog         Human-readable changelog from summaries
    (FR-17.4)
  - FLOW line updated to include documentation steps (FR-17.5):
      map? -> assess? -> init -> design (ROADMAP, DECISIONS) -> wrap? ->
      slice 1 (summary) -> ... -> ship -> roadmap milestone -> ...
  - Three-views tagline added:
      Three views: /gl:status (machine), /gl:roadmap (product), /gl:changelog (history)
  - BUILD section /gl:slice description updated to include "-> summary" step:
      /gl:slice <N>         TDD loop: test -> implement ->
                            security -> verify -> commit -> summary

Exact output (from DESIGN.md 5.5):
  INSIGHT
    /gl:roadmap           Product roadmap + milestones
    /gl:changelog         Human-readable changelog from summaries

Errors: None (help always succeeds)

Invariants:
  - INSIGHT section is always present in help output (not conditional)
  - INSIGHT appears between MONITOR and SHIP sections (FR-17.3)
  - FLOW line includes documentation steps: design (ROADMAP, DECISIONS), slice (summary) (FR-17.5)
  - Three-views tagline is always present
  - Does not remove or modify any existing help sections
  - This change is additive to C-28 (BROWNFIELD section). Both changes apply.

Dependencies: C-28 (brownfield help section must be applied first or simultaneously)
```

### C-40: StatusDocumentationReference

```
Contract: StatusDocumentationReference
Boundary: /gl:status command -> User (roadmap and changelog reference line)
Slice: S-13 (Documentation Infrastructure and Design Update)
Design refs: FR-16.4, DESIGN.md 5.6

COMMAND UPDATE: src/commands/gl/status.md

Output (additional line at bottom of status display):
  Product view: /gl:roadmap | History: /gl:changelog

  Displayed unconditionally after the "Next:" recommendation.
  Provides discoverability for the human-readable views.

Errors: None (status always succeeds; this is a static line)

Invariants:
  - Reference line is always displayed, even if ROADMAP.md does not exist yet (FR-16.4)
  - Line appears after the "Next:" recommendation, at the very bottom of status output
  - Does not modify any existing status display content
  - This change is additive to C-27 (brownfield status display). Both changes apply.

Dependencies: C-27 (brownfield status display must be applied first or simultaneously)
```

---

## S-14: Auto-Summaries and Decision Aggregation

*User Action: "User can see what was built and why after each slice, wrap, or quick task"*

### C-41: SliceSummaryGeneration

```
Contract: SliceSummaryGeneration
Boundary: /gl:slice orchestrator -> Summary Task -> Filesystem
Slice: S-14 (Auto-Summaries and Decision Aggregation)
Design refs: FR-15.7, FR-15.10, FR-25.1, FR-25.2, FR-25.3, FR-25.4, FR-25.5, FR-25.6, FR-25.7, FR-29.3, FR-29.4, FR-29.5, TD-16, NFR-9, NFR-12

COMMAND UPDATE: src/commands/gl/slice.md

Behaviour (added after verification succeeds in existing /gl:slice pipeline):
  1. Collect structured data from the completed slice:
     - Slice ID, slice name, milestone (if any)
     - Contracts satisfied (names from GRAPH.json)
     - Test count and pass/fail results
     - Key files changed (from git diff --stat)
     - Deviation log entries (if any)
     - Security results summary
     - Decision notes from each agent's output (see C-44)
  2. Spawn a Task with fresh context (TD-16, NFR-9)
  3. Pass structured data to the Task via XML context blocks
  4. Task writes .greenlight/summaries/{slice-id}-SUMMARY.md (FR-25.1)
  5. Task checks if architecture changed: new service, new external integration,
     new database table, new endpoint group (FR-29.3)
  6. If architecture changed, Task updates Architecture Diagram section in ROADMAP.md (FR-29.4)
  7. If architecture did NOT change, diagram is NOT modified (FR-29.5)
  8. After Task completes, orchestrator updates ROADMAP.md:
     mark slice as complete, add completion date, test count, key decision (FR-25.7)

Input (structured data passed to Task):
  - slice_id: string
  - slice_name: string
  - milestone: string (optional)
  - contracts_satisfied: string[] (contract names)
  - test_count: number
  - test_results: "all passing" | "{N} failing"
  - key_files: string[] (from git diff --stat)
  - deviations: string[] (deviation log entries, or empty)
  - security_summary: string (security scan results)
  - decisions: { decision: string, chosen: string, context: string }[] (from agents)
  - architecture_context: string (current Mermaid diagram from ROADMAP.md)

Output:
  - .greenlight/summaries/{slice-id}-SUMMARY.md (following schema in DESIGN.md 4.7)

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | TaskSpawnFailure | Task cannot be spawned | Slice is still considered complete. Warn user: "Summary generation failed. Run /gl:changelog to check for gaps." (NFR-12) |
  | TaskWriteFailure | Task fails to write summary file | Same as above. Warn user. Do not retry (NFR-12) |
  | RoadmapUpdateFailure | ROADMAP.md update fails (file missing or write error) | Warn user. Slice is still considered complete. Summary file may exist without roadmap update |
  | ArchitectureDiagramFailure | Task cannot determine if architecture changed | Do not modify diagram. Warn user. Slice is still complete |
  | NoSummariesDir | .greenlight/summaries/ directory does not exist | Create directory with 0o755 permissions before writing |

Invariants:
  - Summary generation is mandatory -- every completed slice triggers it (FR-25.5)
  - Summary failure does NOT block the TDD pipeline (NFR-12)
  - Summary is written in product language, not implementation language (FR-25.4)
  - Summary is NOT over-templated -- Task writes natural-language informed by structured data (FR-25.6)
  - Task receives structured data, it does NOT discover data by reading files (NFR-9)
  - Task MUST complete within a single invocation (NFR-9)
  - Architecture diagram is only updated if architecture actually changed (FR-29.4, FR-29.5)
  - Architecture diagram remains text-based Mermaid (FR-29.6, TD-19)
  - ROADMAP.md is updated after summary is written (FR-25.7)
  - Summary file naming: .greenlight/summaries/{slice-id}-SUMMARY.md (DESIGN.md 4.7)

Security:
  - Summaries include security results summary, not full details (DESIGN.md 6.4)
  - No sensitive data (credentials, API keys, PII) appears in summaries

Dependencies: C-36 (ROADMAP.md must exist for update; created by /gl:design)
```

### C-42: WrapSummaryGeneration

```
Contract: WrapSummaryGeneration
Boundary: /gl:wrap orchestrator -> Summary Task -> Filesystem
Slice: S-14 (Auto-Summaries and Decision Aggregation)
Design refs: FR-14.4, FR-14.5, FR-26.1, FR-26.2, FR-26.3, FR-26.4, TD-14, TD-16, NFR-9, NFR-12

COMMAND UPDATE: src/commands/gl/wrap.md

Behaviour (added after successful wrap commit in existing /gl:wrap pipeline):
  1. Collect structured data from the completed wrap:
     - Boundary name
     - Contracts extracted (count and names)
     - Locking tests written (count)
     - Known security issues (count and severities)
     - Files covered (file paths in the boundary)
  2. Spawn a Task with fresh context (TD-16, NFR-9)
  3. Task writes .greenlight/summaries/{boundary-name}-wrap-SUMMARY.md (FR-26.1)
  4. After Task completes, orchestrator updates ROADMAP.md wrap progress
     section if ROADMAP.md exists (FR-14.5)

Input (structured data passed to Task):
  - boundary_name: string
  - contracts_extracted: { name: string, boundary: string }[]
  - contracts_count: number
  - locking_tests_count: number
  - security_issues: { count: number, severities: string }
  - files_covered: string[] (file paths)

Output:
  - .greenlight/summaries/{boundary-name}-wrap-SUMMARY.md

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | TaskSpawnFailure | Task cannot be spawned | Wrap is still considered complete. Warn user (NFR-12) |
  | TaskWriteFailure | Task fails to write summary file | Warn user. Wrap is still complete (NFR-12) |
  | RoadmapMissing | ROADMAP.md does not exist for wrap progress update | Skip roadmap update. No error |
  | NoSummariesDir | .greenlight/summaries/ does not exist | Create directory before writing |

Invariants:
  - Wrap summary is mandatory -- every completed wrap triggers it (TD-14)
  - Summary failure does NOT block wrap completion (NFR-12)
  - Summary is in product language (FR-26.3): "Locked the authentication boundary..."
  - Wrap summaries appear in /gl:changelog output alongside slice summaries (FR-26.4)
  - Task receives structured data, does NOT discover by reading files (NFR-9)
  - Summary file naming: .greenlight/summaries/{boundary-name}-wrap-SUMMARY.md (DESIGN.md 4.7)
  - ROADMAP.md wrap progress update is best-effort (skip if ROADMAP.md missing)

Security:
  - Security issues referenced by count and severity only, not full details
  - No sensitive data in wrap summaries

Dependencies: C-36 (ROADMAP.md for wrap progress update; optional, skip if missing)
```

### C-43: QuickSummaryGeneration

```
Contract: QuickSummaryGeneration
Boundary: /gl:quick orchestrator -> Summary Task -> Filesystem
Slice: S-14 (Auto-Summaries and Decision Aggregation)
Design refs: FR-27.1, FR-27.2, FR-27.3, FR-27.4, TD-16, NFR-9, NFR-12

COMMAND UPDATE: src/commands/gl/quick.md

Behaviour (added after /gl:quick completes):
  1. Collect structured data from the completed quick task:
     - Timestamp (ISO 8601)
     - Task description (what was done)
     - Test count and results
     - Key files changed (from git diff --stat)
     - Whether a decision was involved
     - Associated milestone (via user confirmation, if applicable)
  2. Spawn a Task with fresh context (TD-16, NFR-9)
  3. Task writes .greenlight/summaries/quick-{timestamp}-SUMMARY.md (FR-27.1)
  4. If the quick task involved a decision, append to DECISIONS.md (FR-27.3)
  5. If the quick task is associated with a milestone, update ROADMAP.md (FR-27.4)

Input (structured data passed to Task):
  - timestamp: string (ISO 8601, used in filename)
  - description: string (what was done)
  - test_count: number
  - test_results: string
  - key_files: string[] (from git diff --stat)
  - decision: { decision: string, chosen: string, context: string } | null

Output:
  - .greenlight/summaries/quick-{timestamp}-SUMMARY.md
  - .greenlight/DECISIONS.md (appended, if decision was made)
  - .greenlight/ROADMAP.md (updated, if associated with milestone)

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | TaskSpawnFailure | Task cannot be spawned | Quick is still considered complete. Warn user (NFR-12) |
  | DecisionAppendFailure | Cannot append to DECISIONS.md | Warn user. Decision is lost. Quick is still complete |
  | NoDecisionsFile | DECISIONS.md does not exist | Create DECISIONS.md with header and the decision entry |
  | NoSummariesDir | .greenlight/summaries/ does not exist | Create directory before writing |

Invariants:
  - Quick summaries follow the same format as slice summaries but are more concise (FR-27.2)
  - Summary failure does NOT block quick task completion (NFR-12)
  - Decision appending follows the same rules as slice decisions: sequential numbering, source="quick" (FR-28.3)
  - DECISIONS.md is append-only (NFR-11)
  - Milestone association is optional and confirmed by user (FR-27.4)
  - Summary file naming: .greenlight/summaries/quick-{timestamp}-SUMMARY.md (DESIGN.md 4.7)
  - Timestamp format in filename: ISO 8601 date-time with hyphens replacing colons

Security:
  - No sensitive data in quick summaries
  - Decision entries do not expose security implementation details

Dependencies: C-37 (DECISIONS.md schema; optional, created if missing)
```

### C-44: DecisionAggregation

```
Contract: DecisionAggregation
Boundary: /gl:slice orchestrator -> Filesystem (.greenlight/DECISIONS.md)
Slice: S-14 (Auto-Summaries and Decision Aggregation)
Design refs: FR-15.9, FR-28.5, FR-28.6, TD-17, NFR-11

COMMAND UPDATE: src/commands/gl/slice.md

Behaviour (added after verification succeeds in existing /gl:slice pipeline):
  1. Collect decision notes from each agent's output during the slice:
     - Test writer: decisions about test patterns, fixtures, strategies
     - Implementer: decisions about algorithms, libraries, patterns
     - Security: decisions about security approaches
     - Verifier: observations about notable design choices
  2. Filter for meaningful decisions (not every implementation choice is a decision)
  3. Format each decision as a DECISIONS.md table row:
     - # : next sequential number (D-{N})
     - Decision: what was decided
     - Context: why this decision was needed
     - Chosen: what was chosen
     - Rejected: what was considered and rejected (or "-")
     - Date: date of slice completion
     - Source: "slice:{slice-id}" (FR-28.3)
  4. Append to .greenlight/DECISIONS.md

Input:
  - Agent outputs containing decision notes
  - Current DECISIONS.md (to determine next sequential number)
  - Slice ID (for source column)

Output:
  - .greenlight/DECISIONS.md (appended with new decision entries)

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | NoDecisionsFile | DECISIONS.md does not exist | Create DECISIONS.md with header and the decision entries |
  | NoDecisions | Agents produced no meaningful decisions | Do not append. No error |
  | NumberingConflict | Cannot determine next sequential number | Read DECISIONS.md, find highest D-{N}, increment |

Invariants:
  - DECISIONS.md is append-only (NFR-11) -- orchestrator never modifies existing entries
  - Decision numbering is sequential across the entire file (FR-28.6)
  - Source format for slices: "slice:{slice-id}" (FR-28.3)
  - Each agent notes decisions in its output; orchestrator aggregates (TD-17)
  - Decision aggregation failure does NOT block slice completion
  - Not every implementation choice is a decision -- filter for meaningful architectural
    or pattern choices that a future reader would want to know about
  - Decisions are captured while agent context is fresh (TD-17)

Security:
  - Security-related decisions are recorded without exposing implementation details
  - No sensitive data in decision entries

Dependencies: C-37 (DECISIONS.md schema; optional, created if missing)
```

### C-45: RoadmapAutoUpdate

```
Contract: RoadmapAutoUpdate
Boundary: Orchestrators (slice/wrap) -> Filesystem (.greenlight/ROADMAP.md)
Slice: S-14 (Auto-Summaries and Decision Aggregation)
Design refs: FR-15.8, FR-14.5, FR-23.5, FR-25.7

COMMAND UPDATES: src/commands/gl/slice.md, src/commands/gl/wrap.md

Behaviour for /gl:slice (after summary generation):
  1. Read ROADMAP.md
  2. Find the slice row in the current milestone's slice table
  3. Update the row:
     - Status: "complete"
     - Tests: {N} (test count from verification)
     - Completed: {YYYY-MM-DD}
     - Key Decision: {brief summary of most significant decision, or "-"}
  4. Write updated ROADMAP.md

Behaviour for /gl:wrap (after wrap summary generation):
  1. Read ROADMAP.md (if it exists)
  2. Find or create the Wrap Progress section
  3. Add or update the wrapped boundary row:
     - Boundary: {name}
     - Status: wrapped
     - Locking Tests: {N}
     - Known Issues: {N}
  4. Write updated ROADMAP.md

Input:
  - .greenlight/ROADMAP.md (current state)
  - Slice completion data (for slice updates): slice ID, test count, completion date, key decision
  - Wrap completion data (for wrap updates): boundary name, locking test count, known issue count

Output:
  - .greenlight/ROADMAP.md (updated in place)

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | NoRoadmap | ROADMAP.md does not exist | Skip update. Warn user: "ROADMAP.md not found. Run /gl:design to create it." |
  | SliceNotFound | Slice ID not found in any milestone table | Append slice to the active milestone table |
  | ParseFailure | Cannot parse ROADMAP.md structure | Skip update. Warn user. Do not corrupt existing file |

Invariants:
  - ROADMAP.md updates are best-effort -- failure does not block slice or wrap completion
  - Slice rows are updated in place (matched by slice ID), not appended as duplicates
  - Wrap progress rows are updated in place (matched by boundary name)
  - ROADMAP.md structure is preserved during updates (no sections removed or reordered)
  - Updates happen after summary generation, as the final documentation step
  - When ROADMAP.md does not exist, no roadmap update occurs (no error, just a warning)

Security:
  - ROADMAP.md does not include security details, only known issue counts
  - No sensitive data written to ROADMAP.md

Dependencies: C-36 (ROADMAP.md created by /gl:design; optional for wrap updates)
```

---

## S-15: Roadmap Command (/gl:roadmap)

*User Action: "User can view product roadmap and plan new milestones"*

### C-46: RoadmapDisplay

```
Contract: RoadmapDisplay
Boundary: User -> /gl:roadmap command (read-only display)
Slice: S-15 (Roadmap Command)
Design refs: FR-24.1, FR-24.8, NFR-10

COMMAND DEFINITION: src/commands/gl/roadmap.md

Input:
  - User invokes /gl:roadmap (no arguments)
  - .greenlight/ROADMAP.md must exist

Output:
  - Display the contents of ROADMAP.md to the user
  - No files modified (read-only) (NFR-10)

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | NoRoadmap | .greenlight/ROADMAP.md does not exist | Print "No roadmap found. Run /gl:design to create one." and stop |
  | EmptyRoadmap | ROADMAP.md exists but is empty | Print "ROADMAP.md is empty. Run /gl:design to populate it." and stop |

Invariants:
  - /gl:roadmap (no arguments) is strictly read-only (NFR-10)
  - No files are created, modified, or deleted
  - Displays the full contents of ROADMAP.md including architecture diagram, milestone tables, and archived milestones
  - config.json is read for project context (FR-24.8)

Security:
  - Read-only operation, no risk of data modification

Dependencies: C-36 (ROADMAP.md must exist; created by /gl:design)
```

### C-47: RoadmapMilestonePlanning

```
Contract: RoadmapMilestonePlanning
Boundary: User -> /gl:roadmap milestone -> gl-designer (scoped design session)
Slice: S-15 (Roadmap Command)
Design refs: FR-19.6, FR-24.2, FR-24.3, FR-24.4, FR-24.5, FR-24.8, TD-13, TD-21

COMMAND DEFINITION: src/commands/gl/roadmap.md (milestone sub-command)

Input:
  - User invokes /gl:roadmap milestone
  - .greenlight/ROADMAP.md must exist
  - .greenlight/config.json must exist

Orchestration steps:
  1. Read .greenlight/config.json for project context and model resolution (FR-24.8)
  2. Read .greenlight/ROADMAP.md (current milestones, architecture) (FR-24.3)
  3. Read .greenlight/DESIGN.md (existing design context) (FR-24.3)
  4. Read .greenlight/CONTRACTS.md (existing contracts)
  5. Read .greenlight/STATE.md (wrapped boundaries, slice status)
  6. Display current milestone status to user
  7. Resolve gl-designer model from config.json profiles
  8. Spawn gl-designer via Task with milestone scope (TD-13):
     - Receives: ROADMAP.md, DESIGN.md, CONTRACTS.md, STATE.md
     - Skips: init interview, stack decisions (FR-24.4)
     - Runs: lighter design session (goal, user actions, constraints)
     - Produces: new slices with milestone field (TD-21)
  9. Append new milestone section to ROADMAP.md (FR-24.5)
  10. Append new slices to GRAPH.json with milestone field (FR-24.5, TD-21)
  11. Append design decisions to DECISIONS.md (FR-24.5)
  12. Commit with message: "docs: plan milestone {milestone-name}"

Output:
  - .greenlight/ROADMAP.md (new milestone section appended)
  - .greenlight/GRAPH.json (new slices appended with milestone field)
  - .greenlight/DECISIONS.md (new design decisions appended)

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | NoRoadmap | ROADMAP.md does not exist | Print "No roadmap found. Run /gl:design first." and stop |
  | NoConfig | config.json does not exist | Print "No config found. Run /gl:init first." and stop |
  | DesignerFailure | gl-designer agent fails to produce output | Report error to user. Do not commit. Suggest retrying |
  | NoDesignMd | DESIGN.md does not exist | Print "No design found. Run /gl:design first." and stop |

Invariants:
  - Milestone planning spawns gl-designer with milestone scope (TD-13)
  - Designer session is lighter: no init interview, no stack decisions (FR-24.4)
  - New slices include milestone field matching the new milestone name (TD-21)
  - New slices are appended to GRAPH.json, not replacing existing slices
  - Existing milestones in ROADMAP.md are preserved (NFR-10)
  - Decisions are appended to DECISIONS.md with source="milestone" (FR-28.3)
  - DECISIONS.md is append-only (NFR-11)
  - A slice belongs to at most one milestone (TD-21)
  - Commit format: "docs: plan milestone {milestone-name}"

Security:
  - No sensitive data in milestone planning
  - Same security constraints as /gl:design

Dependencies: C-36 (ROADMAP.md must exist), C-37 (DECISIONS.md must exist)
```

### C-48: RoadmapMilestoneArchive

```
Contract: RoadmapMilestoneArchive
Boundary: User -> /gl:roadmap archive -> Filesystem (.greenlight/ROADMAP.md)
Slice: S-15 (Roadmap Command)
Design refs: FR-24.6, FR-24.7

COMMAND DEFINITION: src/commands/gl/roadmap.md (archive sub-command)

Input:
  - User invokes /gl:roadmap archive
  - .greenlight/ROADMAP.md must exist with at least one completed milestone

Orchestration steps:
  1. Read .greenlight/ROADMAP.md
  2. Identify completed milestones (all slices in status=complete)
  3. Present completed milestones to user for selection
  4. User picks milestone to archive
  5. Compress the milestone:
     - Move from active section to Archived Milestones section
     - Format: "{milestone-name} -- completed {date}"
     - Summary: "{N} slices, {N} tests. {one-line summary of what was achieved.}" (FR-24.7)
  6. Write updated ROADMAP.md
  7. Commit with message: "docs: archive milestone {milestone-name}"

Output:
  - .greenlight/ROADMAP.md (milestone moved to Archived section, compressed)

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | NoRoadmap | ROADMAP.md does not exist | Print "No roadmap found. Run /gl:design first." and stop |
  | NoCompletedMilestones | No milestones have all slices complete | Print "No completed milestones available for archiving." and stop |
  | ArchiveFailure | Cannot parse or update ROADMAP.md | Report error. Do not corrupt file. Suggest manual archiving |

Invariants:
  - Only completed milestones (all slices done) can be archived
  - Archived milestones are compressed: name, completion date, slice count, test count, one-line summary (FR-24.7)
  - Active and planning milestones are never archived
  - Archiving is a one-way operation (archived milestones are not restored)
  - Archived Milestones section is at the bottom of ROADMAP.md
  - Commit format: "docs: archive milestone {milestone-name}"

Security:
  - No sensitive data involved in archiving

Dependencies: C-36 (ROADMAP.md must exist)
```

---

## S-16: Changelog Command (/gl:changelog)

*User Action: "User can see a human-readable changelog of everything that was built"*

### C-49: ChangelogDisplay

```
Contract: ChangelogDisplay
Boundary: User -> /gl:changelog command (read-only display)
Slice: S-16 (Changelog Command)
Design refs: FR-30.1, FR-30.4, FR-30.5, FR-30.6, FR-30.7, TD-20

COMMAND DEFINITION: src/commands/gl/changelog.md

Input:
  - User invokes /gl:changelog (no arguments)
  - .greenlight/summaries/ directory must exist with at least one summary file
  - .greenlight/config.json for project context (FR-30.7)

Output:
  - Display a formatted changelog to the user (read-only, no files written) (FR-30.6)
  - Format (from DESIGN.md 5.4):
      CHANGELOG -- {project name}

      {date}  {type}:{name}     {one-line summary}    {N} tests
      {date}  {type}:{name}     {one-line summary}    {N} tests

      {N} entries ({N} slices, {N} wraps, {N} quick)

Behaviour:
  1. Read .greenlight/config.json for project name
  2. Scan .greenlight/summaries/ directory for all summary files
  3. Parse each summary file to extract: date, type (slice/wrap/quick), name, one-line summary, test count
  4. Sort chronologically, newest first (FR-30.4)
  5. Format and display

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | NoSummariesDir | .greenlight/summaries/ does not exist | Print "No summaries found. Summaries are generated after each /gl:slice, /gl:wrap, or /gl:quick." and stop |
  | EmptySummariesDir | summaries/ exists but contains no summary files | Print "No summaries found yet. Complete a slice, wrap, or quick task to generate summaries." and stop |
  | NoConfig | config.json does not exist | Use "Unknown Project" as project name. Continue |
  | MalformedSummary | A summary file cannot be parsed | Skip that entry. Warn user: "Could not parse {filename}" |

Invariants:
  - /gl:changelog is strictly read-only -- no files are written (FR-30.6, TD-20)
  - Changelog is formatted chronologically, newest first (FR-30.4)
  - Each entry includes: date, type:name, one-line summary, test count (FR-30.5)
  - Type values: "slice", "wrap", "quick"
  - Summary files are identified by naming convention:
    {slice-id}-SUMMARY.md, {boundary}-wrap-SUMMARY.md, quick-{timestamp}-SUMMARY.md
  - Changelog reads from summaries/ directory only (TD-20)
  - Unparseable summary files are skipped, not fatal
  - Entry count summary is always displayed at the bottom

Security:
  - Read-only operation, no risk of data modification
  - No sensitive data exposed in changelog output

Dependencies: None (command can be built and tested independently; reads from summaries/ directory)
```

### C-50: ChangelogFiltering

```
Contract: ChangelogFiltering
Boundary: /gl:changelog command -> Summary filtering (milestone and date filters)
Slice: S-16 (Changelog Command)
Design refs: FR-30.2, FR-30.3

COMMAND DEFINITION: src/commands/gl/changelog.md (milestone and since sub-commands)

Input (milestone filter):
  - User invokes /gl:changelog milestone {name}
  - {name} is a milestone name matching the milestone field in GRAPH.json slices

Input (date filter):
  - User invokes /gl:changelog since {DATE}
  - {DATE} is an ISO 8601 date (YYYY-MM-DD)

Behaviour for milestone filter:
  1. Read GRAPH.json to find all slices with milestone={name}
  2. Filter summaries to include only those matching the milestone's slice IDs
  3. Include wrap summaries if the wrap boundary is associated with the milestone
  4. Display filtered changelog with header: "CHANGELOG -- {milestone-name}"

Behaviour for date filter:
  1. Parse {DATE} as ISO 8601 date
  2. Filter summaries to include only those with completion date >= {DATE}
  3. Display filtered changelog with header: "CHANGELOG -- since {DATE}"

Output:
  - Filtered changelog display (same format as unfiltered, but with subset of entries)
  - No files written (read-only)

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | UnknownMilestone | Milestone name not found in GRAPH.json | Print "Milestone '{name}' not found in GRAPH.json." and stop |
  | InvalidDate | Date argument is not valid ISO 8601 | Print "Invalid date format: {DATE}. Use YYYY-MM-DD." and stop |
  | NoMatchingEntries | Filter matches zero summaries | Print "No changelog entries found for {filter}." and stop |
  | NoGraphJson | GRAPH.json not found (milestone filter only) | Print "No GRAPH.json found. Cannot filter by milestone." and stop |

Invariants:
  - Filtered changelog uses the same display format as unfiltered (FR-30.4, FR-30.5)
  - Milestone filter matches by slice ID, not by text search
  - Date filter uses >= comparison (inclusive of the given date)
  - Both filters are read-only -- no files written (FR-30.6)
  - Date parsing accepts YYYY-MM-DD format only
  - Milestone names are case-sensitive, matching GRAPH.json exactly

Security:
  - Read-only operation
  - No sensitive data exposed

Dependencies: C-49 (changelog display base; same command file, different sub-command)
```

---

## Cross-Cutting: Summary Task Specification

```
Contract: SummaryTaskSpecification
Boundary: Orchestrators (slice/wrap/quick) -> Task (summary generation)
Not a separate slice -- referenced by S-14 (C-41, C-42, C-43)
Design refs: FR-25.2, FR-25.3, FR-25.4, FR-25.6, TD-16, NFR-9, NFR-12

TASK SPECIFICATION (not a separate agent definition):

The summary Task is a fresh-context Task invocation that receives structured data and
writes a markdown file. It is spawned by the orchestrator, not by a user command.

Input (XML context blocks passed to Task):
  <summary_context>
    <type>{slice|wrap|quick}</type>
    <id>{slice-id or boundary-name or timestamp}</id>
    <name>{human-readable name}</name>
    <milestone>{milestone name, if any}</milestone>
    <date>{YYYY-MM-DD}</date>
  </summary_context>

  <results>
    <contracts>{list of contracts satisfied}</contracts>
    <tests>{count and pass/fail}</tests>
    <files>{git diff --stat output}</files>
    <deviations>{deviation log entries}</deviations>
    <security>{security scan summary}</security>
    <decisions>{decision notes from agents}</decisions>
  </results>

  <architecture_context>
    {current Mermaid diagram from ROADMAP.md, for architecture change detection}
  </architecture_context>

Output:
  - One markdown file in .greenlight/summaries/ following naming convention
  - Updated ROADMAP.md architecture diagram if architecture changed

Task behaviour:
  - Write a natural-language summary in product language (FR-25.4)
  - Summary is NOT over-templated -- natural language informed by structured data (FR-25.6)
  - Follow the summary schema in DESIGN.md section 4.7
  - For slices: check if architecture changed and update ROADMAP.md diagram if yes
  - For wraps: describe what was locked in product language (FR-26.3)
  - For quick: more concise than slice summaries (FR-27.2)
  - Complete within a single Task invocation (NFR-9)

Invariants:
  - Task receives ALL data it needs via context blocks -- it does NOT read files to discover data (NFR-9)
  - Task is spawned with fresh context to avoid quality degradation (TD-16)
  - Task failure does NOT block the parent pipeline (NFR-12)
  - No sensitive data in summaries (DESIGN.md 6.4)
  - This is NOT a separate agent definition -- it is a Task call within existing orchestrators
```

---

## Cross-Cutting: GRAPH.json milestone Field

```
Contract: GraphJsonMilestoneField
Boundary: GRAPH.json -> Slice objects (optional milestone field)
Not a separate slice -- referenced by S-13 (design produces it), S-15 (roadmap milestone adds it)
Design refs: TD-21, DESIGN.md 4.4

FIELD SPECIFICATION:

  "milestone": "{milestone-name}"   // optional string field on slice objects

Rules:
  - Optional field on slice objects in GRAPH.json
  - String type matching a milestone name in ROADMAP.md
  - When not specified, slice belongs to the initial/default milestone (from first /gl:design)
  - A slice belongs to at most one milestone
  - Created by /gl:roadmap milestone when spawning a scoped design session
  - Used by /gl:changelog milestone {name} to filter summaries
  - Does NOT affect dependency resolution or wave ordering

Invariants:
  - Field is optional -- all existing slices work without it
  - No impact on build order or dependency graph
  - Milestone names are simple strings, not structured objects
```

---

## S-17: Brownfield-Roadmap Integration

### C-51: BrownfieldDesignContext

```
Contract: BrownfieldDesignContext
Boundary: /gl:design orchestrator -> gl-designer (brownfield context blocks)
Slice: S-17

design.md passes brownfield context to the designer agent:
  - Reads ASSESS.md (conditional, 2>/dev/null)
  - Reads CONTRACTS.md (conditional, for [WRAPPED] tags)
  - Reads STATE.md (conditional, for wrap progress)
  - Passes <existing_assessment> context block
  - Passes <existing_contracts> context block
  - Passes <existing_state> context block

Invariants:
  - All reads are conditional (2>/dev/null) -- greenfield projects work without them
  - Context blocks default to 'No assessment yet' / 'No contracts yet' / 'No state yet'
  - Existing design.md functionality unchanged for greenfield projects
```

### C-52: BrownfieldRoadmapContext

```
Contract: BrownfieldRoadmapContext
Boundary: /gl:roadmap milestone -> gl-designer (assessment + wrap progress)
Slice: S-17

roadmap.md milestone planning passes brownfield context to designer:
  - Reads ASSESS.md in Gather Context section (conditional, 2>/dev/null)
  - Passes <existing_assessment> in Task spawn block
  - Passes <wrap_progress> in Task spawn block (from STATE.md Wrapped Boundaries)

Invariants:
  - ASSESS.md read is conditional -- projects without assessments work fine
  - <existing_assessment> appears after Spawn gl-designer heading
  - <wrap_progress> appears after Spawn gl-designer heading
  - Existing milestone planning functionality unchanged for non-brownfield projects
```

### C-53: DesignerBrownfieldAwareness

```
Contract: DesignerBrownfieldAwareness
Boundary: gl-designer agent -> Brownfield-aware design (risk tiers, [WRAPPED] tags)
Slice: S-17

gl-designer.md handles brownfield context:
  - Documents <existing_assessment>, <existing_contracts>, <existing_state> in context_protocol
  - References [WRAPPED] tag for boundaries with locking tests
  - Supports milestone_planning session mode (skip init phases, focus on milestone scope)
  - References risk tiers (Critical/High/Medium) for slice prioritization
  - References wrap progress / wrapped boundaries for milestone ordering
  - Output checklist includes brownfield-specific items

Invariants:
  - Brownfield context is in context_protocol section (ordering)
  - All brownfield handling is conditional -- greenfield projects unaffected
  - milestone_planning mode skips Phase 1-2 and Phase 4
```

---

## Updated User Action Mapping

| User Action | Slice(s) | Contracts | Enabled By |
|-------------|----------|-----------|------------|
| 1. User can assess an existing codebase for gaps, risks, and untested boundaries | S-08 | C-17, C-18, C-19, C-20 | /gl:assess command + gl-assessor agent |
| 2. User can wrap existing code in contracts and locking tests without rewriting it | S-09 | C-21, C-22, C-23, C-24, C-25, C-26 | /gl:wrap command + gl-wrapper agent |
| 3. User can see brownfield progress alongside greenfield slices | S-10 | C-27, C-28, C-29 | /gl:status, /gl:help, /gl:settings updates |
| 4. User can refactor a wrapped boundary through the normal TDD loop | S-11 | C-30, C-31, C-32 | /gl:slice wraps field + architect/test-writer updates |
| 5. User can start a project with a product roadmap and decision log from day one | S-13 | C-36, C-37, C-38, C-39, C-40 | /gl:design update + infrastructure |
| 6. User can see what was built and why after each slice, wrap, or quick task | S-14 | C-41, C-42, C-43, C-44, C-45 | Auto-summaries + decision aggregation |
| 7. User can view product roadmap and plan new milestones | S-15 | C-46, C-47, C-48 | /gl:roadmap command |
| 8. User can see a human-readable changelog of everything that was built | S-16 | C-49, C-50 | /gl:changelog command |
| 9. Brownfield context informs design and milestone planning | S-17 | C-51, C-52, C-53 | /gl:design + /gl:roadmap milestone + gl-designer updates |

---

---

# Circuit Breaker Module Contracts

> **Scope:** Attempt tracking, structured diagnostics, scope lock, manual override (/gl-debug), rollback via git tags
> **Deliverables:** Markdown prompt files (reference doc + command + agent/command updates), NOT Go code
> **Date:** 2026-02-18
> **Design Reference:** Circuit Breaker System Design — FR-1 through FR-8

---

## Circuit Breaker Contract Index

| # | Contract | Boundary | Slice |
|---|----------|----------|-------|
| C-54 | CircuitBreakerProtocol | Reference doc -> All agents (circuit breaker rules) | S-18 |
| C-55 | ImplementerCircuitBreaker | gl-implementer agent -> Circuit breaker protocol (error recovery rewrite) | S-18 |
| C-56 | ScopeLockProtocol | gl-implementer agent -> Contract-inferred scope (justify-or-stop) | S-18 |
| C-57 | SliceCheckpointTags | /gl:slice orchestrator -> Git (lightweight checkpoint tags) | S-19 |
| C-58 | SliceRollbackIntegration | /gl:slice orchestrator -> Rollback on circuit break | S-19 |
| C-59 | DebugCommand | User -> /gl-debug command (manual diagnostic override) | S-20 |
| C-60 | CLAUDEmdCircuitBreakerRule | CLAUDE.md -> All agents (5-line hard rule) | S-21 |
| C-61 | ManifestCircuitBreakerUpdate | Go CLI -> Manifest (2 new file paths) | S-21 |

---

## S-18: Circuit Breaker Protocol and Implementer Integration

*User Actions:*
- *1. Implementer automatically stops after 3 failed attempts on a test (instead of looping endlessly)*
- *2. User receives a structured diagnostic report when the circuit trips*
- *3. Implementer justifies every out-of-scope file modification before making it*

### C-54: CircuitBreakerProtocol

```
Contract: CircuitBreakerProtocol
Boundary: Reference doc -> All agents that read it (gl-implementer, /gl:slice, /gl-debug)
Slice: S-18 (Circuit Breaker Protocol and Implementer Integration)
Design refs: FR-1, FR-2, FR-3, FR-4, FR-5

FILE SPECIFICATION: src/references/circuit-breaker.md (~180 lines)

This is the authoritative protocol document. All circuit breaker behaviour
is defined here and referenced by agents and commands.

Output (mandatory sections in the reference doc):

  1. Attempt Tracking State
     - State schema per test per slice:
       - slice_id: string
       - test_name: string
       - attempt_count: number (0-3)
       - files_touched_per_attempt: string[][] (array of arrays)
       - description_per_attempt: string[]
       - last_error_per_attempt: string[]
       - checkpoint_tag: string (greenlight/checkpoint/{slice_id})
     - Slice-level accumulator:
       - total_failed_attempts: number (0-7)
     - State is maintained in-memory by the implementer agent across attempts
       within a single spawn

  2. Per-Test Trip Threshold
     - After 3 failed attempts on any single test, circuit trips (FR-2)
     - Mandatory, not configurable
     - Each "attempt" = one complete cycle of: read error, hypothesize fix,
       modify files, run test, observe result
     - Attempt counter increments only on test FAILURE, not on infrastructure
       errors (syntax, import, missing dep)

  3. Slice-Level Ceiling
     - If total_failed_attempts across ALL tests in a slice exceeds 7,
       circuit trips regardless of per-test counts (FR-3)
     - This catches the pattern: 2 failures on test A, 2 on test B,
       2 on test C, 2 on test D = 8 total, trips even though no single
       test hit 3

  4. Structured Diagnostic Report Format
     - When circuit trips, produce a markdown report with these fields (FR-4):
       a. test_expectation: what the test expects (from test name + contract)
       b. actual_error: exact error output from the last test run (verbatim)
       c. attempt_log: table with columns: Attempt, Hypothesis, Files Touched, Result
       d. cumulative_files_modified: deduplicated list of all files touched across all attempts
       e. scope_violations: list of any out-of-scope file modifications with justifications
       f. best_hypothesis: the implementer's best guess at what's wrong
       g. specific_question: one concrete question for the human (not "what should I do?")
       h. recovery_options: numbered list including rollback command

  5. Scope Lock Protocol
     - Before modifying any file, check if within inferred scope (FR-5)
     - Inferred scope = files referenced by contracts for the current slice,
       plus files in packages listed in GRAPH.json slice definition
     - Optional override: files_in_scope field in GRAPH.json slice object
     - Out-of-scope modification requires justification:
       - Justification must reference the specific failing test
       - Justification must explain why the out-of-scope file must change
     - Unjustifiable out-of-scope modification = scope violation = counts as
       a failed attempt
     - Scope violations are tracked and reported in the diagnostic report

  6. Counter Reset Protocol
     - Triggered when human provides input after a circuit break (FR-8)
     - Per-test counters reset to 0
     - Slice accumulator resets to 0
     - Rollback to checkpoint tag: git checkout greenlight/checkpoint/{slice_id}
     - Fresh implementer spawned with:
       - Original contracts and test expectations
       - User's guidance/input
       - Summary of "what was tried" (from diagnostic report)
       - Clean codebase state (from checkpoint)

  7. Additive to Deviation Rules
     - Circuit breaker protocol is additive to deviation-rules.md
     - Deviation Rule 4 (ARCH-STOP) still takes priority over circuit breaker
     - If an architectural stop is needed, report it immediately — do not
       count it as a failed attempt

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | NoCheckpointTag | Checkpoint tag does not exist when rollback requested | Warn user. Cannot roll back. Start from current state |
  | ScopeInferenceFailure | Cannot infer scope from contracts/GRAPH.json | Default to all files in the slice's packages list. Warn user scope lock is degraded |
  | StateCorruption | Attempt state becomes inconsistent | Reset counters, warn user, continue from current state |

Invariants:
  - Per-test threshold is always 3 (not configurable)
  - Slice-level ceiling is always 7 (not configurable)
  - Diagnostic report is ALWAYS produced when circuit trips (never skipped)
  - Scope lock applies to EVERY file modification (no exceptions)
  - Circuit breaker does NOT modify test files (agent isolation preserved)
  - Protocol is additive to deviation rules (does not replace them)
  - Attempt counters count test FAILURES only, not infrastructure errors
  - The reference doc is read-only at runtime (no agent writes to it)

Security:
  - Error output in diagnostic reports may contain file paths but NOT
    credentials, tokens, or PII
  - Checkpoint tags are lightweight git tags (no signed tags, no GPG)

Dependencies: None (self-contained reference document)
```

### C-55: ImplementerCircuitBreaker

```
Contract: ImplementerCircuitBreaker
Boundary: gl-implementer agent -> Circuit breaker protocol (error recovery section rewrite)
Slice: S-18 (Circuit Breaker Protocol and Implementer Integration)
Design refs: FR-1, FR-2, FR-3, FR-4, FR-5, FR-8

AGENT UPDATE: src/agents/gl-implementer.md — REWRITE of <error_recovery> section (~40 lines)

The existing <error_recovery> section (lines 145-182) is replaced entirely.
The new section integrates the circuit breaker protocol from
references/circuit-breaker.md into the implementer's error handling flow.

Input (additional context from orchestrator):
  - checkpoint_tag: string (greenlight/checkpoint/{slice_id})
  - files_in_scope: string[] (inferred from contracts + GRAPH.json, or explicit override)
  - what_was_tried: string (from previous diagnostic report, if counter was reset)
  - user_guidance: string (from human input after circuit break, if any)

Behaviour (replaces existing error_recovery):

  1. Maintain Attempt State
     - Track per-test attempt count (starts at 0)
     - Track slice-level total failed attempts (starts at 0)
     - Track files touched per attempt
     - Track hypothesis and result per attempt

  2. Before Every File Modification — Scope Check
     - Determine if file is in scope (from files_in_scope)
     - If in scope: proceed
     - If out of scope: generate justification tied to failing test
     - If justification is valid (references specific test, explains why): proceed,
       log as scope deviation
     - If justification is not valid: do NOT modify file, count as failed attempt,
       log scope violation

  3. On Test Failure — Increment and Check
     - Increment per-test attempt count
     - Increment slice-level total
     - If per-test count >= 3: TRIP — produce diagnostic, stop
     - If slice total > 7: TRIP — produce diagnostic, stop
     - If neither threshold reached: read error carefully, hypothesize,
       attempt targeted fix

  4. On Circuit Trip — Produce Diagnostic Report
     - Generate structured diagnostic per C-54 section 4
     - Report to orchestrator with full diagnostic
     - STOP implementation — do not attempt further fixes
     - Include rollback command: git tag -d greenlight/checkpoint/{slice_id} &&
       git checkout greenlight/checkpoint/{slice_id}

  5. On Infrastructure Error (syntax, import, missing dep)
     - Fix infrastructure issue (this is Rule 3: Unblock from deviation rules)
     - Do NOT increment attempt counter
     - Re-run test
     - If same test failure after infrastructure fix: NOW increment counter

Output (structured diagnostic when circuit trips):
  - Diagnostic report in markdown format per C-54 section 4
  - Report is returned to orchestrator (not written to file)

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | PerTestTrip | 3 failures on single test | Produce diagnostic for that test. Stop |
  | SliceCeilingTrip | >7 total failures across all tests | Produce diagnostic covering all failing tests. Stop |
  | ScopeViolation | Unjustifiable out-of-scope modification | Count as failed attempt. Log violation. Continue if under threshold |
  | InfrastructureError | Syntax/import/dep error (not test logic) | Fix without incrementing counter. If persists, treat as test failure |

Invariants:
  - Implementer reads references/circuit-breaker.md at start (added to "Read first" list)
  - Attempt counters are maintained in-memory for the duration of the agent spawn
  - Scope check happens BEFORE every file write (not after)
  - Infrastructure errors are distinguished from test failures (no counter increment)
  - Diagnostic report is structured (not free-form prose)
  - The implementer NEVER modifies test files (existing prohibition preserved)
  - The implementer NEVER disables or skips tests (existing prohibition preserved)
  - "What was tried" context from previous attempts is used to avoid repeating
    the same fix strategy

Security:
  - Diagnostic reports must not include credentials, tokens, or PII
  - Scope lock prevents unauthorized file modifications

Dependencies: C-54 (circuit breaker protocol must be defined first)
```

### C-56: ScopeLockProtocol

```
Contract: ScopeLockProtocol
Boundary: gl-implementer agent -> Contract-inferred scope (justify-or-stop)
Slice: S-18 (Circuit Breaker Protocol and Implementer Integration)
Design refs: FR-5

SCOPE LOCK SPECIFICATION (within circuit-breaker.md and gl-implementer.md)

Scope inference rules (in priority order):

  1. Explicit override: If GRAPH.json slice object has files_in_scope field,
     use that list exclusively
     ```json
     {
       "id": "S-XX",
       "files_in_scope": ["internal/auth/", "internal/middleware/auth.go"]
     }
     ```

  2. Inferred from contracts: Parse contract definitions for the current slice.
     Extract:
     - Package names from contract boundary descriptions
     - File paths from contract FILE SPECIFICATION fields
     - File paths from GRAPH.json slice "packages" or "deliverables" fields

  3. Fallback: If neither explicit nor inferred scope is available,
     use the slice's "packages" or "deliverables" field from GRAPH.json
     as the scope boundary

Input:
  - Current slice contracts (from CONTRACTS.md)
  - Current slice definition (from GRAPH.json)
  - File path being modified

Output:
  - is_in_scope: boolean
  - If out of scope: justification_required: boolean (always true)

Justification format:
  ```
  SCOPE JUSTIFICATION
  File: {file_path}
  Failing test: {test_name}
  Reason: {why this file must change to make the test pass}
  Relationship: {how this file relates to the contract boundary}
  ```

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | NoScopeData | No contracts, no GRAPH.json packages, no files_in_scope | All files are in scope (scope lock disabled). Warn in diagnostic |
  | AmbiguousScope | Multiple interpretations possible | Use the union of all interpretations (broader scope) |

Invariants:
  - Scope check runs BEFORE every file modification
  - In-scope files are never blocked (no justification needed)
  - Out-of-scope files ALWAYS require justification (no exceptions)
  - Justification must reference a specific failing test by name
  - Unjustifiable modification = scope violation = failed attempt
  - files_in_scope field in GRAPH.json is optional (not required)
  - Scope lock is additive to existing agent isolation rules
    (implementer still cannot modify test files, regardless of scope)
  - Paths in files_in_scope support both file paths and directory paths
    (directory = all files recursively within)

Dependencies: C-54 (protocol defines scope lock rules), C-55 (implementer enforces them)
```

---

## S-19: Slice Checkpoint Tags and Rollback Integration

*User Actions:*
- *5. Implementer rolls back to clean checkpoint state after human provides input*

### C-57: SliceCheckpointTags

```
Contract: SliceCheckpointTags
Boundary: /gl:slice orchestrator -> Git (lightweight checkpoint tags)
Slice: S-19 (Slice Checkpoint Tags and Rollback Integration)
Design refs: FR-7, FR-8

COMMAND UPDATE: src/commands/gl/slice.md — Modify Step 3 (add checkpoint tag creation)

Behaviour (added to Step 3 "Check for Previous Attempt" in existing /gl:slice):

  After pre-flight validation passes and before Step 1 (Write Tests):

  1. Create lightweight checkpoint tag:
     ```bash
     git tag greenlight/checkpoint/{slice_id}
     ```

  2. If tag already exists (from a previous attempt):
     ```bash
     # Remove old tag, create fresh one at current HEAD
     git tag -d greenlight/checkpoint/{slice_id} 2>/dev/null
     git tag greenlight/checkpoint/{slice_id}
     ```

  3. Report:
     ```
     Checkpoint created: greenlight/checkpoint/{slice_id}
     ```

  4. Pass checkpoint_tag to implementer context:
     ```xml
     <checkpoint>
     Tag: greenlight/checkpoint/{slice_id}
     Rollback: git checkout greenlight/checkpoint/{slice_id}
     </checkpoint>
     ```

Tag naming convention:
  - Format: greenlight/checkpoint/{slice_id}
  - Example: greenlight/checkpoint/S-18
  - Lightweight tags only (no annotated/signed tags)

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | GitTagFailure | git tag command fails | Warn user. Proceed without checkpoint (rollback unavailable). Do not block slice |
  | DirtyWorkingTree | Uncommitted changes when creating tag | Warn user: "Uncommitted changes detected. Checkpoint may not represent a clean state." Proceed anyway |

Invariants:
  - Checkpoint tag is created BEFORE any test writing or implementation
  - Tag points to HEAD at the moment pre-flight completes
  - Tag is lightweight (not annotated, not signed)
  - Old tag for same slice is always replaced (idempotent)
  - Tag creation failure does not block the slice pipeline
  - Tag is passed to implementer as part of context

Dependencies: None (git tag is a standard operation)
```

### C-58: SliceRollbackIntegration

```
Contract: SliceRollbackIntegration
Boundary: /gl:slice orchestrator -> Rollback on circuit break
Slice: S-19 (Slice Checkpoint Tags and Rollback Integration)
Design refs: FR-7, FR-8

COMMAND UPDATE: src/commands/gl/slice.md — Modify Step 3 (handle circuit break)
and Step 10 (cleanup tag on success)

Behaviour when implementer reports circuit break (from C-55):

  1. Receive structured diagnostic from implementer
  2. Present diagnostic report to user:
     ```
     CIRCUIT BREAK -- Slice {slice_id}

     {formatted diagnostic report from C-54 section 4}

     Recovery options:
     1) Provide guidance and retry (rollback to checkpoint, fresh implementer)
     2) Spawn debugger to investigate (/gl-debug)
     3) Pause and review manually (/gl:pause)

     Checkpoint: greenlight/checkpoint/{slice_id}
     Rollback command: git checkout greenlight/checkpoint/{slice_id}
     ```

  3. If user chooses option 1 (guidance + retry):
     a. Collect user guidance
     b. Roll back to checkpoint:
        ```bash
        git checkout greenlight/checkpoint/{slice_id} -- .
        ```
     c. Reset attempt counters (per C-54 section 6)
     d. Spawn fresh implementer with:
        - Original contracts and test expectations
        - User guidance
        - "What was tried" summary from diagnostic report
        - Clean codebase state (from rollback)
     e. Create new checkpoint tag at current HEAD

  4. If user chooses option 2 (debugger):
     a. Pass diagnostic report to /gl-debug (see C-59)

  5. If user chooses option 3 (pause):
     a. Save diagnostic to .greenlight/.continue-here.md
     b. Update STATE.md with circuit break info

Behaviour on slice completion (successful):

  1. After Step 10 (Complete) — clean up checkpoint tag:
     ```bash
     git tag -d greenlight/checkpoint/{slice_id} 2>/dev/null
     ```
  2. Report: `Checkpoint tag cleaned up: greenlight/checkpoint/{slice_id}`

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | RollbackFailure | git checkout fails (conflicts, missing tag) | Report error to user. Suggest manual recovery. Do not retry automatically |
  | NoCheckpointTag | Tag does not exist when rollback requested | Warn user: "No checkpoint tag found. Cannot roll back. Starting from current state." Proceed with retry from current state |
  | TagCleanupFailure | git tag -d fails on completion | Ignore. Stale tag is harmless |

Invariants:
  - Rollback uses `git checkout {tag} -- .` (restore working tree, not detach HEAD)
  - Counter reset happens AFTER rollback, BEFORE fresh implementer spawn
  - Fresh implementer always receives "what was tried" context
  - User MUST provide guidance before retry (not an automatic retry)
  - Tag cleanup on success is best-effort (failure does not block completion)
  - Circuit break presentation uses the structured diagnostic format (not free-form)
  - Existing /gl:slice error recovery (agent spawn failure, context overflow,
    state corruption) is preserved — circuit break is an ADDITIONAL error path

Security:
  - Diagnostic report presented to user may contain file paths and error messages
    but NOT credentials or tokens
  - Rollback does not discard committed changes (only working tree changes)

Dependencies: C-57 (checkpoint tags must be created first), C-55 (implementer produces diagnostic)
```

---

## S-20: Debug Command (/gl-debug)

*User Actions:*
- *4. User can force a diagnostic at any time with /gl-debug (manual pull cord)*

### C-59: DebugCommand

```
Contract: DebugCommand
Boundary: User -> /gl-debug command (manual diagnostic override)
Slice: S-20 (Debug Command)
Design refs: FR-6

COMMAND DEFINITION: src/commands/gl/debug.md (~80 lines)

Input:
  - User invokes /gl-debug (no arguments required)
  - Optional: /gl-debug {slice_id} (to specify which slice to diagnose)

Context read:
  - .greenlight/STATE.md (current slice, step, last activity)
  - .greenlight/GRAPH.json (slice definition, contracts)
  - .greenlight/CONTRACTS.md (contract definitions for current slice)
  - .greenlight/config.json (project context)
  - Test results from latest run (if available)
  - Git log for recent commits
  - Git diff for uncommitted changes

Behaviour:

  1. Determine target slice:
     - If slice_id argument provided: use that
     - If STATE.md has current slice: use that
     - If neither: report "No active slice found. Specify a slice: /gl-debug {slice_id}"

  2. Gather diagnostic context:
     a. Read current test results:
        ```bash
        {config.test.command} {config.test.filter_flag} {slice_id} 2>&1
        ```
     b. Read recent git activity:
        ```bash
        git log --oneline -10
        git diff --stat
        ```
     c. Read STATE.md for step and progress
     d. Read contracts for the slice
     e. Check for checkpoint tag:
        ```bash
        git tag -l "greenlight/checkpoint/{slice_id}"
        ```

  3. Produce structured diagnostic report:
     ```
     DIAGNOSTIC REPORT -- Slice {slice_id}: {slice_name}
     Generated: {timestamp}

     ## Current State
     Step: {step from STATE.md}
     Last activity: {date}
     Checkpoint: {tag exists? tag name : "none"}

     ## Test Results
     Total: {N}
     Passing: {N}
     Failing: {N}

     {for each failing test:}
     ### {test_name}
     Expected: {inferred from contract + test name}
     Actual: {exact error output}
     Contracts: {which contract(s) this test verifies}

     ## Recent Changes
     {git log --oneline -10}

     ## Uncommitted Changes
     {git diff --stat}

     ## Files in Scope
     {inferred scope from contracts/GRAPH.json}

     ## Recovery Options
     1) Resume implementation (/gl:slice {slice_id})
     2) Roll back to checkpoint: git checkout greenlight/checkpoint/{slice_id} -- .
     3) Pause for manual investigation (/gl:pause)
     4) Spawn fresh implementer with guidance

     ## Specific Question
     {auto-generated: the most likely root cause based on failing tests and recent changes}
     ```

  4. Present report to user (display only, no files written)

Output:
  - Structured diagnostic report displayed to user
  - No files written (read-only command)
  - Does not modify STATE.md or any project state

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | NoActiveSlice | No current slice in STATE.md and no argument provided | Print "No active slice. Usage: /gl-debug {slice_id}" and stop |
  | NoConfig | config.json does not exist | Print "No config found. Run /gl:init first." and stop |
  | TestRunFailure | Test command fails to execute | Include error in diagnostic: "Test command failed: {error}". Continue with partial diagnostic |
  | NoContracts | CONTRACTS.md missing or no contracts for slice | Include warning in diagnostic: "No contracts found for slice {id}". Continue with partial diagnostic |

Invariants:
  - /gl-debug is strictly read-only (no files written, no state modified)
  - /gl-debug can be run at ANY time, not just during circuit break
  - /gl-debug does not control the pipeline (it diagnoses, it does not resume/retry)
  - Diagnostic report follows the same structured format as circuit break diagnostics (C-54)
  - /gl-debug works without a checkpoint tag (tag presence is informational)
  - /gl-debug runs test suite to get current state (always fresh data)
  - Report is structured for future pause/resume integration
  - Does not spawn any subagents (direct read + display)

Security:
  - Read-only operation
  - Error output may contain file paths but NOT credentials or tokens
  - Does not expose test source code (displays test names and error output only)

Dependencies: None (standalone command, works independently)
```

---

## S-21: Circuit Breaker Infrastructure Integration

*User Actions:*
- Supports all 5 user actions (infrastructure enabling layer)

### C-60: CLAUDEmdCircuitBreakerRule

```
Contract: CLAUDEmdCircuitBreakerRule
Boundary: CLAUDE.md -> All agents (5-line hard rule in standards)
Slice: S-21 (Circuit Breaker Infrastructure Integration)
Design refs: FR-2, FR-3

FILE UPDATE: src/CLAUDE.md

Location: Insert as a new subsection within "Code Quality Constraints" section,
after "Testing" and before "Logging & Observability".

Content (exactly 5 lines, hard rule):
  ### Circuit Breaker
  - After 3 failed attempts on any single test, STOP and produce a structured diagnostic report
  - After 7 total failed attempts across all tests in a slice, STOP regardless of per-test counts
  - Before modifying any file, verify it is within inferred scope from contracts; justify out-of-scope changes
  - Full protocol: `references/circuit-breaker.md`

Errors: None (static content update)

Invariants:
  - Rule is exactly 5 lines (header + 4 bullet points)
  - Rule references the full protocol in references/circuit-breaker.md
  - Rule is placed within Code Quality Constraints section
  - Existing CLAUDE.md sections unchanged
  - This is a hard rule, not a recommendation — phrased as imperatives (STOP, verify, justify)
  - The 3-per-test and 7-per-slice thresholds are stated explicitly in CLAUDE.md
    (agents read CLAUDE.md first, before reading reference docs)

Dependencies: C-54 (references/circuit-breaker.md must be defined; file created in S-18)
```

### C-61: ManifestCircuitBreakerUpdate

```go
// Contract: ManifestCircuitBreakerUpdate
// Boundary: Go CLI -> Manifest (2 new file paths for circuit breaker)
// Slice: S-21 (Circuit Breaker Infrastructure Integration)
//
// FILE: internal/installer/installer.go
//
// Change: Add 2 new entries to Manifest slice
//
// New entries (inserted in alphabetical order within their sections):
//   "commands/gl/debug.md"               // NEW -- /gl-debug command
//   "references/circuit-breaker.md"      // NEW -- circuit breaker protocol
//
// Updated Manifest (34 entries, up from 32):
//   "agents/gl-architect.md"
//   "agents/gl-assessor.md"
//   "agents/gl-codebase-mapper.md"
//   "agents/gl-debugger.md"
//   "agents/gl-designer.md"
//   "agents/gl-implementer.md"
//   "agents/gl-security.md"
//   "agents/gl-test-writer.md"
//   "agents/gl-verifier.md"
//   "agents/gl-wrapper.md"
//   "commands/gl/add-slice.md"
//   "commands/gl/assess.md"
//   "commands/gl/changelog.md"
//   "commands/gl/debug.md"              <-- NEW
//   "commands/gl/design.md"
//   "commands/gl/help.md"
//   "commands/gl/init.md"
//   "commands/gl/map.md"
//   "commands/gl/pause.md"
//   "commands/gl/quick.md"
//   "commands/gl/resume.md"
//   "commands/gl/roadmap.md"
//   "commands/gl/settings.md"
//   "commands/gl/ship.md"
//   "commands/gl/slice.md"
//   "commands/gl/status.md"
//   "commands/gl/wrap.md"
//   "references/checkpoint-protocol.md"
//   "references/circuit-breaker.md"     <-- NEW
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
// - Entries within each section (agents/, commands/gl/, references/) are
//   alphabetically ordered
// - go:embed directive in main.go already uses wildcards
//   (src/commands/gl/*.md, src/references/*.md) so new .md files in
//   those directories are automatically embedded -- no main.go change needed
// - Manifest count increases from 32 to 34
// - All existing tests that validate manifest count must be updated to expect 34
// - This change is additive to C-33 and C-38 (previous manifest updates)
//
// Dependencies: C-33, C-38 (previous manifest updates must be applied first or simultaneously)
```

---

## Cross-Cutting: GRAPH.json files_in_scope Field

```
Contract: GraphJsonFilesInScopeField
Boundary: GRAPH.json -> Slice objects (optional files_in_scope field)
Not a separate slice -- referenced by S-18 (scope lock uses it)
Design ref: FR-5

FIELD SPECIFICATION:

  "files_in_scope": ["path/to/file.go", "path/to/dir/"]   // optional string array

Rules:
  - Optional field on slice objects in GRAPH.json
  - String array of file paths and/or directory paths (relative to project root)
  - Directory paths (ending in /) include all files recursively within
  - When not specified, scope is inferred from contracts and slice packages/deliverables
  - When specified, overrides inferred scope entirely
  - Used by implementer's scope lock (C-55, C-56) to validate file modifications
  - Does NOT affect dependency resolution, wave ordering, or any other GRAPH.json feature

Invariants:
  - Field is optional -- all existing slices work without it
  - No impact on build order or dependency graph
  - Paths are forward-slash separated (consistent with GRAPH.json conventions)
  - Empty array means "no files in scope" (all modifications require justification)
```

---

## Updated User Action Mapping (Circuit Breaker)

| User Action | Slice(s) | Contracts | Enabled By |
|-------------|----------|-----------|------------|
| 1. Implementer automatically stops after 3 failed attempts on a test | S-18, S-21 | C-54, C-55, C-60 | Circuit breaker protocol + implementer rewrite + CLAUDE.md rule |
| 2. User receives a structured diagnostic report when the circuit trips | S-18, S-19 | C-54, C-55, C-58 | Diagnostic format in protocol + implementer produces it + orchestrator presents it |
| 3. Implementer justifies every out-of-scope file modification | S-18 | C-54, C-55, C-56 | Scope lock protocol + implementer enforcement |
| 4. User can force a diagnostic at any time with /gl-debug | S-20, S-21 | C-59, C-61 | Debug command + manifest entry |
| 5. Implementer rolls back to clean checkpoint state after human provides input | S-19 | C-57, C-58 | Checkpoint tags + rollback integration |

---

## Verification Tiers Milestone

---

## S-22: Schema Extension

*User Actions:*
- *1. Architect can set a verification tier (auto/verify) on each contract, with verify as the default*

### C-62: ContractSchemaExtension

```
Contract: ContractSchemaExtension
Boundary: gl-architect.md -> Contract format (three new optional fields)
Slice: S-22 (Schema Extension)
Design refs: FR-7, DESIGN.md 4.6

FILE UPDATE: src/agents/gl-architect.md — Extend <contract_format> section

Three new optional fields added to the contract format template,
after the Security section and before the Dependencies line:

Content (added to the per-contract template in <contract_format>):

  **Verification:** verify
  **Acceptance Criteria:**
  - {behavioral criterion the user can verify}
  - {another criterion}

  **Steps:**
  - {step to verify, when how-to-verify is not obvious}
  - {another step}

Field rules:
  - verification: Optional. Values: "auto", "verify". Default: "verify".
    - "auto": slice proceeds directly from verification to summary/docs
      after tests pass. No human checkpoint. Use for infrastructure,
      config, and internal plumbing contracts.
    - "verify": slice presents an acceptance checkpoint after tests pass.
      Human must approve before slice completes. Use for everything else.
  - acceptance_criteria: Optional list under verify tier. Behavioral
    statements the user can check. Each criterion is a testable assertion
    about what the user should observe.
  - steps: Optional list under verify tier. How-to-verify instructions
    when the verification process is not obvious. Include commands to run,
    URLs to visit, or actions to perform.
  - If verification is "verify" and both acceptance_criteria and steps
    are empty, emit a warning: "Contract {name} has verify tier but no
    acceptance criteria or steps. Consider adding at least one."
  - If verification is "auto", acceptance_criteria and steps are ignored
    (present but not surfaced in the checkpoint).

Output (updated contract format template in gl-architect.md):

  The <contract_format> section gains three new lines in the per-contract
  template, positioned after **Security:** and before **Dependencies:**:

  **Verification:** auto | verify (default: verify)
  **Acceptance Criteria:**
  - [behavioral criterion the user can verify]

  **Steps:**
  - [step to verify the feature, when how-to-verify is not obvious]

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | InvalidTierValue | verification field has value other than auto/verify | Reject contract. Error: "Invalid verification tier: {value}. Must be auto or verify." |
  | EmptyVerifyCriteria | tier is verify but both acceptance_criteria and steps are empty | Warn (not error): "Contract {name} has verify tier but no acceptance criteria or steps." |

Invariants:
  - Default tier is always "verify" (safe default per TD-2 in DESIGN.md)
  - Existing contracts without verification field default to "verify"
  - The three fields are optional -- contracts missing them are valid
  - Field names are exactly: Verification, Acceptance Criteria, Steps
  - acceptance_criteria items are behavioral (what the user observes),
    not implementation (how the code works)
  - steps items are actionable instructions (run X, open Y, click Z),
    not descriptions of internal behaviour
  - Fields are positioned after Security and before Dependencies in
    the contract template (consistent ordering)

Security:
  - No security impact. Fields are metadata on the contract format.

Verification: auto
Dependencies: None (this is the foundation contract for the milestone)
```

### C-63: VerifierTierAwareness

```
Contract: VerifierTierAwareness
Boundary: gl-verifier.md -> Verification report (tier reporting addition)
Slice: S-22 (Schema Extension)
Design refs: FR-1, FR-4, DESIGN.md 6.6

FILE UPDATE: src/agents/gl-verifier.md — Add tier awareness to verification output

The verifier agent gains awareness of verification tier fields in contracts.
This is informational only -- the verifier reports tier status, it does not
enforce the verification gate. The orchestrator enforces the gate.

Behaviour (additions to verifier report):

  1. For each contract in the slice, read the verification field
     (default: "verify" if absent)

  2. Compute effective tier for the slice:
     - verify > auto (highest tier wins)
     - If any contract has tier "verify", effective tier is "verify"
     - If all contracts have tier "auto", effective tier is "auto"

  3. Include in verification report:
     ```
     ## Verification Tier
     Effective tier: {verify|auto}
     Per-contract tiers:
       - {contract_name}: {tier} {criteria_count} criteria, {steps_count} steps
       - {contract_name}: {tier} {criteria_count} criteria, {steps_count} steps

     Warnings:
       - {contract_name}: verify tier with no acceptance criteria or steps
     ```

  4. Flag contracts with verify tier but empty acceptance_criteria and
     empty steps: include in Warnings subsection.

Input (additional context verifier reads):
  - verification field from each contract (default: "verify")
  - acceptance_criteria list from each contract (default: empty)
  - steps list from each contract (default: empty)

Output (additions to existing verification report):
  - Effective tier for the slice (string: "verify" or "auto")
  - Per-contract tier breakdown with criteria/steps counts
  - Warnings for verify contracts missing both criteria and steps

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | MissingVerificationField | Contract has no verification field | Default to "verify". Note in report: "defaulted to verify" |
  | InvalidTierInContract | Contract has unrecognised verification value | Report as warning: "Unknown tier '{value}' on {contract}, treating as verify" |

Invariants:
  - Verifier ONLY reports tier information -- does not enforce the gate
  - The orchestrator reads the effective tier from the verifier report
    to decide whether to run Step 6b
  - Effective tier computation is deterministic: verify > auto
  - Missing verification field always defaults to "verify"
  - Warnings are informational, not blocking (verifier still passes)
  - The verifier report format is additive -- existing sections unchanged
  - Tier reporting section appears after existing report sections

Security:
  - No security impact. Informational metadata in verification report.

Verification: auto
Dependencies: C-62 (contract format must include verification fields first)
```

---

## S-23: Verification Gate

*User Actions:*
- *2. After tests pass, a slice with verify tier presents acceptance criteria and optional steps -- user must approve before slice completes*

### C-64: VerificationTiersProtocol

```
Contract: VerificationTiersProtocol
Boundary: Reference doc -> Orchestrator + agents (verification tier protocol)
Slice: S-23 (Verification Gate)
Design refs: FR-1, FR-2, FR-3, FR-4, DESIGN.md 4.2, 4.3

FILE SPECIFICATION: src/references/verification-tiers.md (~130 lines)

This is the authoritative protocol document for verification tiers.
All verification tier behaviour is defined here and referenced by
/gl:slice and agent markdown files.

Output (mandatory sections in the reference doc):

  1. Tier Definitions
     - Two tiers: auto and verify
     - auto: tests pass -> slice proceeds to summary/docs
     - verify: tests pass -> human acceptance checkpoint
     - Default: verify (safe default)
     - Tier is set per-contract in CONTRACTS.md verification field

  2. Tier Resolution
     - Per-slice resolution: highest tier wins (verify > auto)
     - If any contract in the slice has tier verify, effective is verify
     - Acceptance criteria aggregated from all verify contracts
     - Steps aggregated from all verify contracts
     - One checkpoint per slice (not per contract)

  3. Verify Checkpoint Format
     - Header: "ALL TESTS PASSING -- Slice {id}: {name}"
     - Body: aggregated acceptance_criteria as checklist
     - Body: aggregated steps as numbered list
     - Prompt: three options (approve, reject with description, partial)
     - Format adapts: criteria only, steps only, both, or neither
     - If neither criteria nor steps exist: simple "Does the output
       match your intent?" prompt

  4. Rejection Flow
     - Non-"approved" response triggers classification
     - Three options presented to user (implicit gap classification):
       1. Tighten the tests (test gap)
       2. Revise the contract (contract gap)
       3. Provide more detail (implementation gap)
     - Routing by choice (see C-67, C-68, C-69)

  5. Rejection Counter
     - Per-slice count, increments on each non-"approved" response
     - Escalation at 3 rejections (see C-70, C-71)
     - Counter is maintained in orchestrator context for the /gl:slice
       execution lifetime

  6. Agent Isolation in Rejection Loop
     - Test writer receives: verbatim user feedback, contract,
       acceptance criteria. No implementation code.
     - Implementer receives: test names only. No test source code.
     - Existing isolation rules are preserved.

  7. Backward Compatibility
     - Contracts without verification field default to verify
     - config.workflow.visual_checkpoint is deprecated
     - If visual_checkpoint is true, log warning: "visual_checkpoint
       is deprecated. Verification tiers in contracts supersede it."

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | NoContractsForSlice | Slice has no contracts in CONTRACTS.md | Skip verification gate entirely. Warn: "No contracts found for slice {id}. Skipping verification gate." |
  | AllContractsAuto | All contracts have tier auto | Skip checkpoint. Proceed to Step 7. Log: "All contracts are tier auto. Skipping acceptance checkpoint." |

Invariants:
  - Reference doc is read-only at runtime (no agent writes to it)
  - Two tiers only: auto and verify (no third tier)
  - Default is always verify (safe default)
  - Rejection routing always goes through test writer first (TDD-correct)
  - Escalation threshold is always 3 (not configurable)
  - Acceptance checkpoints always pause, even in yolo mode
  - Protocol is additive to existing /gl:slice pipeline
  - Protocol does not interact with circuit breaker (different pipeline steps)

Security:
  - User feedback in rejection flow may describe application behaviour
    but MUST NOT include credentials, tokens, or PII
  - Feedback is passed to test writer as behavioral description only

Verification: auto
Dependencies: None (self-contained reference document)
```

### C-65: VerificationTierGate

```
Contract: VerificationTierGate
Boundary: /gl:slice orchestrator -> Verification tier gate (new Step 6b)
Slice: S-23 (Verification Gate)
Design refs: FR-2, FR-3, FR-4, DESIGN.md 4.2, 6.1

COMMAND UPDATE: src/commands/gl/slice.md — Add Step 6b after Step 6/6a

Step 6b: Verification Tier Gate

This step is inserted after Step 6 (verification passes) and Step 6a
(locking-to-integration transition, if applicable). It reads the
effective verification tier and either skips to Step 7 or presents
an acceptance checkpoint.

Behaviour:

  1. Read verification tier from each contract in the slice
     (from CONTRACTS.md, default: "verify" if field absent)

  2. Compute effective tier:
     - verify > auto (highest wins)
     - If any contract has tier "verify": effective = verify
     - If all contracts have tier "auto": effective = auto

  3. If effective tier is "auto":
     - Log: "Verification tier: auto. Skipping acceptance checkpoint."
     - Proceed to Step 7 (summary/docs)

  4. If effective tier is "verify":
     - Aggregate acceptance_criteria from all verify-tier contracts
     - Aggregate steps from all verify-tier contracts
     - Present Verify Checkpoint (see C-66)
     - Wait for user response
     - Handle response:
       a. "Yes" / approved: proceed to Step 7
       b. Anything else: enter rejection flow (C-67)
          - Increment rejection counter (C-70)
          - Route based on classification (C-68 or C-69)
          - After rejection handling completes, re-run Step 6b

  5. Check for deprecated visual_checkpoint config:
     - Read config.workflow.visual_checkpoint
     - If true, log warning: "visual_checkpoint is deprecated.
       Verification tiers in contracts supersede it. See
       references/verification-tiers.md"
     - Do NOT execute the old visual checkpoint logic

Input:
  - Slice contracts (from CONTRACTS.md)
  - Verifier report (from Step 6, includes effective tier)
  - config.workflow.visual_checkpoint (for deprecation check)

Output:
  - Gate passes: proceed to Step 7
  - Gate triggers checkpoint: user approval required
  - Gate loops: rejection -> test writer -> implementer -> re-verify -> Step 6b

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | ContractReadFailure | Cannot parse contracts from CONTRACTS.md | Warn user. Default to verify tier (safe default). Present generic checkpoint |
  | UserResponseTimeout | User does not respond to checkpoint | Wait indefinitely (checkpoint is blocking by design) |
  | RejectionLoopFailure | Test writer or implementer spawn fails during rejection loop | Report error. Offer: retry, pause, skip verification |

Invariants:
  - Step 6b runs AFTER Step 6 and Step 6a (never before)
  - Step 6b runs BEFORE Step 7 (always)
  - The gate is blocking: /gl:slice cannot proceed past Step 6b without
    either auto tier or user approval
  - Acceptance checkpoints pause even in yolo mode (per DESIGN.md 6.3)
  - After rejection loop completes (new tests pass), Step 6b re-runs
    from the beginning (re-reads tiers, re-aggregates criteria)
  - visual_checkpoint deprecation warning is logged once per slice execution
  - Step 9 (existing visual checkpoint) becomes a no-op with deprecation
    message when verification tiers are active

Security:
  - No new security surface. Checkpoint is displayed in terminal.
  - User feedback is handled per agent isolation rules.

Verification: verify
Acceptance Criteria:
- After tests pass on a verify-tier slice, the orchestrator presents acceptance criteria to the user
- After tests pass on an auto-tier slice, the orchestrator skips directly to summary/docs
- When the user approves, the slice proceeds to Step 7
- When visual_checkpoint is true in config, a deprecation warning is logged

Steps:
- Run /gl:slice on a slice with at least one verify-tier contract
- Observe that after Step 6 verification, a checkpoint is presented with acceptance criteria
- Type "Yes" to approve and observe the slice proceeds to Step 7
- Run /gl:slice on a slice with all auto-tier contracts
- Observe that no checkpoint is presented and the slice proceeds directly to Step 7

Dependencies: C-62 (contracts must have verification fields), C-63 (verifier reports tier), C-64 (protocol defines gate behaviour)
```

### C-66: VerifyCheckpointPresentation

```
Contract: VerifyCheckpointPresentation
Boundary: /gl:slice orchestrator -> User (acceptance checkpoint display)
Slice: S-23 (Verification Gate)
Design refs: FR-3, FR-4, DESIGN.md 4.3

CHECKPOINT FORMAT: Presented to user during Step 6b

Format:

  ALL TESTS PASSING -- Slice {slice_id}: {slice_name}

  Please verify the output matches your intent.

  Acceptance criteria:
    [ ] {criterion 1 from contract A}
    [ ] {criterion 2 from contract A}
    [ ] {criterion 3 from contract B}

  Steps to verify:
    1. {step 1 from contract A}
    2. {step 2 from contract B}

  Does this match what you intended?
    1) Yes -- mark complete and continue
    2) No -- I'll describe what's wrong
    3) Partially -- some criteria met, I'll describe the gaps

Input:
  - slice_id: string
  - slice_name: string
  - aggregated acceptance_criteria: string[] (from all verify-tier contracts)
  - aggregated steps: string[] (from all verify-tier contracts)

Output:
  - user_response: string (one of: "1"/"Yes", "2"/"No", "3"/"Partially",
    or free-text rejection description)

Format adaptation rules:
  - If only criteria exist (no steps): show criteria section, omit steps section
  - If only steps exist (no criteria): show steps section, omit criteria section
  - If both exist: show both sections
  - If neither exists: show simplified prompt:
    "ALL TESTS PASSING -- Slice {id}: {name}\n\n
     Does the output match your intent?\n
     1) Yes -- mark complete\n
     2) No -- I'll describe what's wrong"

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | EmptySliceId | slice_id is empty | Use "unknown" as placeholder. Warn in log |
  | EmptySliceName | slice_name is empty | Use slice_id as name fallback |

Invariants:
  - Checkpoint format follows checkpoint-protocol.md patterns (NFR-4)
  - Criteria are presented as unchecked checkboxes ([ ])
  - Steps are presented as numbered list
  - Three response options are always present (unless simplified prompt)
  - User response "1" or "Yes" (case-insensitive) means approved
  - Any other response means rejection
  - Criteria ordering: contracts appear in CONTRACTS.md order
  - Steps ordering: contracts appear in CONTRACTS.md order
  - No criteria or steps are duplicated in aggregation
  - Checkpoint header always includes "ALL TESTS PASSING"

Security:
  - Checkpoint is displayed in terminal only (not logged to files)
  - Acceptance criteria may describe application behaviour but
    MUST NOT include credentials or secrets

Verification: auto
Dependencies: C-64 (protocol defines checkpoint format), C-65 (gate triggers checkpoint)
```

---

## S-24: Rejection Flow

*User Actions:*
- *3. When user rejects, feedback routes through the test writer (TDD-correct) to produce new tests that catch the mismatch*

### C-67: RejectionClassification

```
Contract: RejectionClassification
Boundary: /gl:slice orchestrator -> User (gap classification UX)
Slice: S-24 (Rejection Flow)
Design refs: FR-5, TD-7, DESIGN.md 4.4

GAP CLASSIFICATION: Presented to user after rejection

When the user responds with anything other than "Yes"/"1" to the
acceptance checkpoint, the orchestrator captures their feedback and
presents classification options.

Behaviour:

  1. Capture verbatim user feedback (their rejection response)

  2. Present classification:
     ```
     Your feedback: "{user's verbatim response}"

     How should we address this?

     1) Tighten the tests -- the tests aren't specific enough to catch
        this mismatch
        (routes to: test writer adds more precise assertions, then
         implementer passes them)

     2) Revise the contract -- the contract doesn't capture what I
        actually want
        (routes to: you update the contract, then the slice restarts)

     3) Provide more detail -- I'll describe exactly what I expect
        (routes to: test writer uses your detail to write targeted
         tests, then implementer passes them)

     Which option? (1/2/3)
     ```

  3. Map user choice to internal classification:
     | Choice | Internal Classification | Route |
     |--------|------------------------|-------|
     | 1 | test_gap | C-68 (spawn test writer with feedback) |
     | 2 | contract_gap | C-69 (user revises contract) |
     | 3 | implementation_gap | C-68 (spawn test writer with detail) |

  4. If choice is 3, collect additional detail:
     ```
     Please describe exactly what you expected:
     ```
     Capture response as detailed_feedback.

Input:
  - user_rejection: string (verbatim rejection from checkpoint)

Output:
  - classification: "test_gap" | "contract_gap" | "implementation_gap"
  - feedback: string (original rejection text)
  - detailed_feedback: string (additional detail, only for option 3)

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | InvalidChoice | User enters something other than 1/2/3 | Re-prompt: "Please choose 1, 2, or 3." Allow free-text after 2 re-prompts (default to test_gap with free-text as feedback) |
  | EmptyFeedback | User rejection was empty string | Prompt for feedback: "Please describe what doesn't match your intent." |

Invariants:
  - Classification options are always presented in the same order (1/2/3)
  - User's verbatim feedback is preserved exactly as typed
  - Option descriptions include the routing consequence (transparency)
  - Users do not need to understand Greenlight's internal taxonomy (TD-7)
  - Options are phrased as actions, not categories
  - Choice 3 always collects additional detail before routing
  - Default after failed re-prompts is test_gap (safest route -- adds tests)

Security:
  - User feedback may describe application behaviour
  - Feedback MUST NOT be used to execute commands or modify files directly
  - Feedback is passed as text context to test writer (no code execution)

Verification: auto
Dependencies: C-65 (gate must trigger rejection flow), C-66 (checkpoint must present options)
```

### C-68: RejectionToTestWriter

```
Contract: RejectionToTestWriter
Boundary: /gl:slice orchestrator -> gl-test-writer agent (rejection feedback routing)
Slice: S-24 (Rejection Flow)
Design refs: FR-5, TD-4, TD-8, NFR-3, DESIGN.md 4.4

TEST WRITER SPAWN: After test_gap or implementation_gap classification

When the user's rejection is classified as test_gap (option 1) or
implementation_gap (option 3), the orchestrator spawns the test writer
with behavioral feedback to write tighter tests.

Behaviour:

  1. Prepare context for test writer spawn:
     ```xml
     <rejection_context>
     <feedback>{user's verbatim rejection feedback}</feedback>
     <classification>{test_gap | implementation_gap}</classification>
     <detailed_feedback>{additional detail from user, if option 3}</detailed_feedback>
     </rejection_context>

     <contract>
     {full contract definition(s) for the verify-tier contracts in this slice}
     </contract>

     <acceptance_criteria>
     {aggregated acceptance criteria that the user was reviewing}
     </acceptance_criteria>
     ```

  2. Test writer writes new or tightened tests that:
     - Assert the specific behaviour the user described
     - Are behavioral (test what the user observes, not implementation)
     - Are additive (do not remove existing passing tests)

  3. After test writer completes, spawn implementer:
     - Implementer receives new test names (not test source code)
     - Implementer makes new tests pass
     - Existing tests must continue passing

  4. After implementer completes, re-run full verification:
     - Step 4 (run tests) to confirm all pass
     - Step 6 (verifier) to confirm contract coverage
     - Step 6b (verification gate) to re-present checkpoint

Input:
  - rejection_feedback: string (user's verbatim words)
  - classification: "test_gap" | "implementation_gap"
  - detailed_feedback: string (empty for test_gap, user's detail for implementation_gap)
  - contracts: string[] (full contract definitions for verify-tier contracts)
  - acceptance_criteria: string[] (aggregated criteria from contracts)

Output:
  - New or modified test files
  - Implementation changes to pass new tests
  - Re-verification cycle

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | TestWriterSpawnFailure | Test writer agent fails to spawn or errors | Report to user. Offer: retry, pause, skip verification |
  | ImplementerSpawnFailure | Implementer agent fails to spawn or errors | Report to user. Offer: retry, pause, skip verification |
  | NewTestsStillFailing | New tests fail after implementer completes | Enter normal circuit breaker flow (C-55). If circuit trips, present diagnostic + rollback |
  | ExistingTestsRegressed | Previously passing tests now fail | Report regression. Implementer must fix regressions before proceeding |

Invariants:
  - Test writer receives ONLY behavioral feedback -- no implementation code (NFR-3)
  - Test writer receives the contract and acceptance criteria -- not test source
  - Implementer receives test names only -- not test source code (existing isolation)
  - New tests are additive -- existing passing tests are not removed or modified
  - The full verification cycle re-runs after rejection handling (not partial)
  - Rejection handling integrates with existing circuit breaker protocol
  - Test writer feedback is the user's verbatim words (no AI summarisation)

Security:
  - Agent isolation is preserved throughout the rejection loop
  - Test writer does not see implementation code
  - Implementer does not see test source code
  - User feedback is treated as behavioral context, not executable content

Verification: verify
Acceptance Criteria:
- When user rejects with "tighten tests", the test writer is spawned with feedback and contract
- Test writer produces new tests that assert the user's described behaviour
- Implementer makes the new tests pass without breaking existing tests
- After implementation, the verification gate re-presents the checkpoint

Steps:
- Reject a slice checkpoint and choose option 1 (tighten tests)
- Observe test writer spawns with the rejection feedback
- Observe implementer spawns after test writer completes
- Observe verification gate re-runs after implementation

Dependencies: C-65 (gate triggers rejection), C-67 (classification routes to test writer)
```

### C-69: RejectionToContractRevision

```
Contract: RejectionToContractRevision
Boundary: /gl:slice orchestrator -> User (contract revision route)
Slice: S-24 (Rejection Flow)
Design refs: FR-5, DESIGN.md 4.4

CONTRACT REVISION: After contract_gap classification

When the user's rejection is classified as contract_gap (option 2),
the orchestrator presents the current contract for user revision and
restarts the slice.

Behaviour:

  1. Present the current contract(s) to the user:
     ```
     CONTRACT REVISION -- Slice {slice_id}: {slice_name}

     The following contracts define this slice's behaviour.
     Edit the acceptance criteria, contract definition, or both.

     Current contract(s):
     ---
     {full contract text for each verify-tier contract}
     ---

     What needs to change?
     ```

  2. Capture user's revision description

  3. Apply revision:
     - If user provides specific acceptance criteria changes: update
       CONTRACTS.md acceptance_criteria and/or steps fields
     - If user describes contract-level changes: flag for architect
       re-engagement (present recommendation to run /gl:add-slice
       with revision context)
     - If changes are minor (criteria wording): apply directly and
       restart from Step 1 (test writing)

  4. Restart slice:
     - Re-run from Step 1 (write tests) with updated contracts
     - Previous implementation is discarded (rollback to checkpoint if available)

Input:
  - contracts: string[] (full verify-tier contract definitions)
  - slice_id: string
  - slice_name: string

Output:
  - Updated contracts in CONTRACTS.md (if criteria/steps revision)
  - Slice restart from Step 1
  - Or: recommendation to re-run /gl:add-slice (if fundamental revision)

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | EmptyRevision | User provides no revision description | Re-prompt: "Please describe what the contract should say instead." |
  | ContractsReadFailure | Cannot read contracts from CONTRACTS.md | Report error. Cannot proceed with revision. Suggest manual edit |
  | ContractsWriteFailure | Cannot write updated contracts to CONTRACTS.md | Report error. Display proposed changes for manual application |

Invariants:
  - Contract revision is always user-driven (no AI auto-revision)
  - Minor revisions (acceptance criteria wording) can be applied in-place
  - Fundamental revisions (input/output changes, new boundaries) require
    architect re-engagement via /gl:add-slice
  - Slice restarts from Step 1 after revision (clean TDD loop)
  - Rollback to checkpoint is offered before restart if tag exists
  - User sees the full contract text before making changes
  - Contract revision increments the rejection counter (C-70)

Security:
  - User revision text is treated as contract metadata, not executable code
  - CONTRACTS.md modifications follow existing write patterns

Verification: auto
Dependencies: C-65 (gate triggers rejection), C-67 (classification routes to contract revision)
```

---

## S-25: Rejection Counter

*User Actions:*
- *4. After 3 rejections on a slice, escalation triggers with options: re-scope, pair, or skip*

### C-70: RejectionCounter

```
Contract: RejectionCounter
Boundary: /gl:slice orchestrator -> Per-slice rejection tracking (in-memory state)
Slice: S-25 (Rejection Counter)
Design refs: FR-6, TD-6, DESIGN.md 4.5

REJECTION COUNTER: Per-slice tracking within /gl:slice execution

The orchestrator tracks rejections per slice. The counter increments
on every non-"approved" response to the acceptance checkpoint,
regardless of classification choice.

Behaviour:

  1. Initialise rejection state at start of Step 6b:
     ```yaml
     slice_id: S-{N}
     rejection_count: 0
     rejection_log: []
     ```

  2. On each rejection (any response other than "1"/"Yes"):
     a. Increment rejection_count
     b. Append to rejection_log:
        ```yaml
        - attempt: {rejection_count}
          feedback: "{user's verbatim response}"
          classification: "{test_gap|contract_gap|implementation_gap}"
          action_taken: "{description of routing action}"
        ```
     c. Check escalation threshold: if rejection_count >= 3, trigger
        escalation (C-71)

  3. Counter persists across rejection loops within a single
     /gl:slice execution. If /gl:slice is re-invoked (new execution),
     counter resets to 0.

Input:
  - rejection event: { feedback: string, classification: string, action: string }

Output:
  - rejection_count: number (0-3+)
  - rejection_log: array of rejection entries
  - escalation_triggered: boolean

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | CounterOverflow | Counter exceeds 3 (should not happen due to escalation) | Trigger escalation immediately. Log warning |
  | LogCorruption | Rejection log becomes inconsistent | Reset log. Preserve count. Warn user |

Invariants:
  - Counter increments on EVERY non-approved response (including contract revisions)
  - Counter is per-slice, not per-contract
  - Counter resets to 0 on new /gl:slice execution (not persisted to disk)
  - Escalation triggers at exactly 3 (not before, not after)
  - Rejection log preserves chronological order
  - Rejection log includes the user's verbatim feedback
  - Counter does not interact with circuit breaker attempt counters
    (different concerns, different pipeline steps)
  - A contract revision (option 2) that restarts the slice from Step 1
    does NOT reset the rejection counter (the counter persists for the
    entire /gl:slice execution)

Security:
  - Rejection log is in-memory only (not written to disk)
  - User feedback in log may contain behavioral descriptions

Verification: auto
Dependencies: C-65 (gate triggers rejection which increments counter)
```

### C-71: RejectionEscalation

```
Contract: RejectionEscalation
Boundary: /gl:slice orchestrator -> User (escalation at 3 rejections)
Slice: S-25 (Rejection Counter)
Design refs: FR-6, TD-6, DESIGN.md 4.5

ESCALATION FORMAT: Presented to user when rejection_count reaches 3

Behaviour:

  1. When rejection_count reaches 3, present escalation:
     ```
     ESCALATION: {slice_name}

     This slice has been rejected 3 times. The verification criteria
     may not match what the contracts and tests can deliver.

     Rejection history:
     1. "{feedback 1}" -> {action taken 1}
     2. "{feedback 2}" -> {action taken 2}
     3. "{feedback 3}" -> {action taken 3}

     Options:
     1) Re-scope -- the contract is fundamentally wrong. Revise
        contracts and restart from scratch.
     2) Pair -- provide detailed, step-by-step guidance for exactly
        what you want.
     3) Skip verification -- mark this slice as auto-verified and
        proceed. (The mismatch is acknowledged but deferred.)

     Which option? (1/2/3)
     ```

  2. Route based on choice:
     | Choice | Action |
     |--------|--------|
     | 1 (re-scope) | Present full contract revision. Recommend /gl:add-slice with revision context. Reset rejection counter. Restart slice from scratch |
     | 2 (pair) | Collect step-by-step guidance from user. Pass as detailed context to test writer. Spawn test writer + implementer. Reset rejection counter. Re-run verification gate |
     | 3 (skip) | Mark slice effective tier as "auto" for this execution. Log: "Verification skipped after 3 rejections. Mismatch acknowledged." Proceed to Step 7 |

Input:
  - rejection_log: array of { feedback, classification, action_taken }
  - slice_id: string
  - slice_name: string

Output:
  - User choice: "re-scope" | "pair" | "skip"
  - Corresponding routing action

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | InvalidEscalationChoice | User enters something other than 1/2/3 | Re-prompt: "Please choose 1, 2, or 3." |
  | EmptyRejectionLog | Escalation triggered but rejection log is empty | Display escalation without history. Warn in log |

Invariants:
  - Escalation triggers at exactly 3 rejections (not configurable)
  - Rejection history is always displayed (full transparency)
  - "Skip verification" is an explicit user choice with logged acknowledgment
  - "Re-scope" resets the rejection counter (fresh start)
  - "Pair" resets the rejection counter (user is actively guiding)
  - "Skip" does not reset the counter (it exits the loop)
  - Escalation threshold of 3 matches circuit breaker per-test threshold
    (consistent system-wide limits)
  - Escalation format follows existing checkpoint patterns
  - After escalation, the slice either restarts (re-scope), continues
    with more tests (pair), or completes with acknowledged mismatch (skip)

Security:
  - Rejection history displayed to user only (not persisted to files)
  - "Skip" option creates an explicit log entry (auditability)

Verification: verify
Acceptance Criteria:
- After 3 rejections on a slice, escalation options are presented
- Re-scope option resets rejection counter and restarts the slice
- Pair option collects detailed guidance and spawns test writer
- Skip option marks the slice as auto-verified and proceeds

Steps:
- Reject a verify-tier slice three times
- Observe the escalation prompt with rejection history
- Choose option 3 (skip) and observe the slice proceeds with a logged acknowledgment

Dependencies: C-70 (counter must track rejections to trigger escalation)
```

---

## S-26: Documentation and Deprecation

*User Actions:*
- *5. Slice with auto tier behaves exactly as today -- no regression (documentation confirms)*
- *Supports all user actions (infrastructure enabling layer)*

### C-72: CLAUDEmdVerificationTierRule

```
Contract: CLAUDEmdVerificationTierRule
Boundary: CLAUDE.md -> All agents (verification tier hard rule in standards)
Slice: S-26 (Documentation and Deprecation)
Design refs: DESIGN.md 4.7

FILE UPDATE: src/CLAUDE.md

Location: Insert as a new subsection within "Code Quality Constraints" section,
after "Circuit Breaker" and before "Logging & Observability".

Content (exactly 4 lines, hard rule):
  ### Verification Tiers
  - Every contract has a verification tier: `auto` or `verify` (default)
  - After tests pass and verifier approves, the tier gate determines if human acceptance is required
  - Rejection feedback routes to the test writer first -- if the implementation is wrong, the tests weren't tight enough
  - Full protocol: `references/verification-tiers.md`

Errors: None (static content update)

Invariants:
  - Rule is exactly 5 lines (header + 4 bullet points)
  - Rule references the full protocol in references/verification-tiers.md
  - Rule is placed within Code Quality Constraints section, after Circuit Breaker
  - Existing CLAUDE.md sections unchanged
  - This is a hard rule, not a recommendation -- phrased as imperatives
  - The default tier (verify) is stated explicitly
  - The TDD-correct rejection routing is stated explicitly (test writer first)

Security:
  - No security impact. Static content update.

Verification: auto
Dependencies: C-64 (references/verification-tiers.md must be defined)
```

### C-73: CheckpointProtocolAcceptanceType

```
Contract: CheckpointProtocolAcceptanceType
Boundary: checkpoint-protocol.md -> Acceptance checkpoint type (new row + Visual deprecation)
Slice: S-26 (Documentation and Deprecation)
Design refs: NFR-4, DESIGN.md 6.3

FILE UPDATE: src/references/checkpoint-protocol.md — Add Acceptance type, deprecate Visual

Changes to the checkpoint type table:

  1. Deprecate Visual checkpoint type:
     - Strike through or mark as deprecated: "~~Visual~~ Deprecated -- use verify tier"
     - Existing description preserved for reference

  2. Add Acceptance checkpoint type:
     | Checkpoint Type | Trigger | When to Pause |
     |-----------------|---------|---------------|
     | Acceptance | Slice verification tier is verify | always (even in yolo mode) |

  3. Add Acceptance checkpoint description:
     - Trigger: Step 6b verification tier gate resolves to verify
     - Content: aggregated acceptance criteria + steps from contracts
     - User action: approve, reject with feedback, or partial
     - Rejection routes: test writer (TDD-correct) or contract revision
     - Escalation: after 3 rejections

  4. Update mode table (if exists) to show Acceptance checkpoints
     always pause, regardless of mode (interactive, yolo, CI).

Additional changes:

  5. Update src/references/verification-patterns.md:
     - Add cross-reference: "For human acceptance verification (per-contract
       tiers, rejection flows, escalation), see references/verification-tiers.md"
     - Insert after existing content, before any closing section

  6. Update src/templates/config.md:
     - Add deprecation note on visual_checkpoint field:
       "Deprecated: visual_checkpoint is superseded by verification tiers in
        contracts. See references/verification-tiers.md. Field is preserved
        for backward compatibility but ignored when verification tiers are
        active."

  7. Update Step 9 of src/commands/gl/slice.md:
     - Replace visual checkpoint logic with deprecation warning:
       "If config.workflow.visual_checkpoint is true, log: 'visual_checkpoint
        is deprecated. Verification tiers in contracts supersede it.
        See references/verification-tiers.md'"
     - Do not execute visual checkpoint logic (it is subsumed by Step 6b)

Errors: None (static content updates)

Invariants:
  - Acceptance checkpoints ALWAYS pause, even in yolo mode
  - Visual checkpoint type is deprecated, not removed (backward compatibility)
  - Cross-reference in verification-patterns.md is informational only
  - Config template deprecation note is informational only
  - Step 9 becomes a no-op with deprecation warning (not removed from pipeline)
  - Existing checkpoint types (Decision, External Action, Circuit Break)
    are unchanged

Security:
  - No security impact. Documentation updates only.

Verification: auto
Dependencies: C-64 (verification-tiers.md must be defined), C-65 (Step 6b must be defined)
```

### C-74: ManifestVerificationTiersUpdate

```go
// Contract: ManifestVerificationTiersUpdate
// Boundary: Go CLI -> Manifest (1 new file path for verification tiers)
// Slice: S-26 (Documentation and Deprecation)
//
// FILE: internal/installer/installer.go
//
// Change: Add 1 new entry to Manifest slice
//
// New entry (inserted in alphabetical order within references/ section):
//   "references/verification-tiers.md"    // NEW -- verification tiers protocol
//
// Updated Manifest (35 entries, up from 34):
//   "agents/gl-architect.md"
//   "agents/gl-assessor.md"
//   "agents/gl-codebase-mapper.md"
//   "agents/gl-debugger.md"
//   "agents/gl-designer.md"
//   "agents/gl-implementer.md"
//   "agents/gl-security.md"
//   "agents/gl-test-writer.md"
//   "agents/gl-verifier.md"
//   "agents/gl-wrapper.md"
//   "commands/gl/add-slice.md"
//   "commands/gl/assess.md"
//   "commands/gl/changelog.md"
//   "commands/gl/debug.md"
//   "commands/gl/design.md"
//   "commands/gl/help.md"
//   "commands/gl/init.md"
//   "commands/gl/map.md"
//   "commands/gl/pause.md"
//   "commands/gl/quick.md"
//   "commands/gl/resume.md"
//   "commands/gl/roadmap.md"
//   "commands/gl/settings.md"
//   "commands/gl/ship.md"
//   "commands/gl/slice.md"
//   "commands/gl/status.md"
//   "commands/gl/wrap.md"
//   "references/checkpoint-protocol.md"
//   "references/circuit-breaker.md"
//   "references/deviation-rules.md"
//   "references/verification-patterns.md"
//   "references/verification-tiers.md"    <-- NEW
//   "templates/config.md"
//   "templates/state.md"
//   "CLAUDE.md"
//
// Errors: none (compile-time constant)
//
// Invariants:
// - CLAUDE.md remains the LAST entry
// - Entries within each section (agents/, commands/gl/, references/) are
//   alphabetically ordered
// - go:embed directive in main.go already uses wildcards
//   (src/references/*.md) so new .md files in references/ are
//   automatically embedded -- no main.go change needed
// - Manifest count increases from 34 to 35
// - All existing tests that validate manifest count must be updated to expect 35
// - This change is additive to C-33, C-38, and C-61 (previous manifest updates)
//
// Verification: auto
// Dependencies: C-61 (previous manifest update must be applied first or simultaneously)
```

---

## S-27: Architect Integration

*User Actions:*
- *1. Architect can set a verification tier (auto/verify) on each contract, with verify as the default*

### C-75: ArchitectTierGuidance

```
Contract: ArchitectTierGuidance
Boundary: gl-architect.md -> Tier selection guidance and acceptance criteria generation
Slice: S-27 (Architect Integration)
Design refs: FR-7, TD-2, TD-3, DESIGN.md build order step 6

FILE UPDATE: src/agents/gl-architect.md — Add tier selection guidance

Two additions to the architect agent:

1. Tier Selection Guidance (new section or addition to <rules>):

   ## Verification Tier Selection

   Every contract you produce should include a verification tier.

   **Default: verify.** When in doubt, use verify. The cost of an
   unnecessary human checkpoint is low (user types "approved"). The
   cost of a missing checkpoint is a completed slice that doesn't
   match intent.

   **When to use auto:**
   - Infrastructure contracts (manifest updates, config changes)
   - Internal plumbing (agent file updates, reference doc updates)
   - Schema/type definitions with no user-visible behaviour
   - Build tooling, CI/CD configuration
   - Contracts where "tests pass" fully captures correctness

   **When to use verify:**
   - Any contract with user-visible behaviour
   - UI components, page layouts, visual output
   - API endpoints where response format matters to the user
   - Business logic where intent may differ from specification
   - Any contract where "tests pass" does NOT fully capture correctness
   - When you are uncertain (verify is the safe default)

   **Writing acceptance criteria:**
   - Each criterion is a behavioral statement the user can observe
   - Use present tense: "User sees X", "Page displays Y", "API returns Z"
   - Be specific: "Cards render in a 3-column grid" not "Layout looks correct"
   - Include negative criteria when relevant: "No error messages appear"
   - 2-5 criteria per contract (more than 5 suggests the contract is too large)

   **Writing steps:**
   - Include when how-to-verify is not obvious
   - Start each step with an action verb: "Run...", "Open...", "Click..."
   - Include commands, URLs, or navigation paths
   - Steps are optional -- omit when criteria are self-explanatory

2. Output Checklist Update (addition to <output_checklist>):

   Add to the checklist:
   - [ ] Every contract has a verification tier (auto or verify)
   - [ ] verify-tier contracts have at least one acceptance criterion or step
   - [ ] auto-tier contracts have a clear reason for skipping human verification
   - [ ] Acceptance criteria are behavioral (what user observes), not implementation

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | MissingTierOnContract | Architect produces contract without verification field | Defaults to verify. Output checklist catches this as a warning |
  | TooManyCriteria | Contract has more than 5 acceptance criteria | Suggest splitting the contract. Not blocking -- just a guideline |

Invariants:
  - Every contract produced by the architect includes a verification field
  - Default is always verify (safe default)
  - Acceptance criteria are behavioral, not implementation
  - Steps are actionable, not descriptive
  - Auto tier requires justification (why tests alone capture correctness)
  - Architect output checklist enforces verification field presence
  - Guidance is non-prescriptive: the architect can override with
    good reasoning (auto for a UI contract with thorough tests, verify
    for an infrastructure contract with user-visible side effects)

Security:
  - No security impact. Agent guidance update.

Verification: auto
Dependencies: C-62 (contract format must include verification fields)
```

---

## Updated User Action Mapping (Verification Tiers)

| User Action | Slice(s) | Contracts | Enabled By |
|-------------|----------|-----------|------------|
| 1. Architect can set a verification tier (auto/verify) on each contract, with verify as the default | S-22, S-27 | C-62, C-63, C-75 | Contract schema extension + architect tier guidance |
| 2. After tests pass, a slice with verify tier presents acceptance criteria and optional steps -- user must approve before slice completes | S-23, S-26 | C-64, C-65, C-66, C-73 | Verification gate + protocol + checkpoint type |
| 3. When user rejects, feedback routes through the test writer (TDD-correct) to produce new tests that catch the mismatch | S-24 | C-67, C-68, C-69 | Rejection classification + test writer routing + contract revision |
| 4. After 3 rejections on a slice, escalation triggers with options: re-scope, pair, or skip | S-25 | C-70, C-71 | Rejection counter + escalation |
| 5. Slice with auto tier behaves exactly as today -- no regression | S-23, S-26 | C-65, C-72, C-73 | Verification gate auto path + CLAUDE.md rule + checkpoint deprecation |

---

# Parallel State Milestone Contracts

> **Scope:** Fix concurrent session state corruption by replacing single STATE.md with per-slice state files
> **Deliverables:** 3 new markdown files, 11 modified markdown files, 3 Go manifest entries
> **Date:** 2026-02-22
> **Design Reference:** DESIGN.md (parallel-state): FR-1 through FR-10, NFR-1 through NFR-4, D-30 through D-38

---

## Parallel State Contract Index

| # | Contract | Boundary | Slice |
|---|----------|----------|-------|
| C-76 | SliceStateTemplate | Template file -> Agents (slice state schema and lifecycle) | S-28 |
| C-77 | StateFormatReference | Reference doc -> All commands (state detection, concurrency, backward compat) | S-28 |
| C-78 | InitSliceDirectory | /gl:init orchestrator -> Filesystem (slices/ directory + slice files) | S-29 |
| C-79 | InitProjectState | /gl:init orchestrator -> Filesystem (project-state.json) | S-29 |
| C-80 | InitStateDetection | All commands -> State detection logic (format resolution) | S-29 |
| C-81 | SliceCommandStateWrite | /gl:slice orchestrator -> Filesystem (own slice file + STATE.md regeneration) | S-30 |
| C-82 | SliceSessionTracking | /gl:slice orchestrator -> Slice file frontmatter (advisory session field) | S-30 |
| C-83 | StatusSliceAggregation | /gl:status command -> Filesystem (read all slice files, compute summary) | S-31 |
| C-84 | SupportingCommandStateAdaptation | Supporting commands -> State detection + slice file reads/writes | S-31 |
| C-85 | MigrateStateCommand | User -> /gl:migrate-state command (legacy STATE.md to file-per-slice) | S-32 |
| C-86 | MigrateStateBackup | /gl:migrate-state -> Filesystem (backup + atomic creation) | S-32 |
| C-87 | CLAUDEmdStateFormatRule | CLAUDE.md -> All agents (state format awareness hard rule) | S-33 |
| C-88 | StateTemplateDocUpdate | templates/state.md -> Both formats documented | S-33 |
| C-89 | CheckpointProtocolStateUpdate | checkpoint-protocol.md -> Slice file references for state context | S-33 |
| C-90 | ManifestParallelStateUpdate | Go CLI -> Manifest (3 new file paths) | S-34 |

---

## S-28: Slice State Foundation Documents

*User Actions:*
- *1. Run multiple /gl:slice sessions in parallel without state corruption (foundation)*
- *2. See accurate slice status across all concurrent sessions (foundation)*

### C-76: SliceStateTemplate

```
Contract: SliceStateTemplate
Boundary: Template file -> Agents (slice state file schema and lifecycle)
Slice: S-28 (Slice State Foundation Documents)
Design refs: FR-1, FR-7, FR-9, D-30, D-33, D-36, DESIGN.md 4.2, 4.3

FILE SPECIFICATION: src/templates/slice-state.md (~120 lines)

This is the authoritative template defining the schema, lifecycle, and
examples for per-slice state files stored in .greenlight/slices/{id}.md.
Agents reference this template when creating or updating slice state files.

Output (mandatory sections in the template):

  1. Schema Definition
     - Frontmatter field definitions (flat key-value between --- delimiters)
     - Fields: id, status, step, milestone, started, updated, tests,
       security_tests, session, deps
     - Field types, allowed values, required/optional status
     - Per D-30: flat key-value only, no nesting

  2. Status Lifecycle
     - Valid status values: pending, ready, tests, implementing, security,
       fixing, verifying, complete
     - Valid step values: none, tests, implementing, security, fixing,
       verifying, complete
     - Transition rules: which status transitions are valid
     - Each status maps to a step in the TDD loop

  3. Session Tracking
     - Advisory session field format: ISO timestamp + random suffix
       (e.g., 2026-02-22T14:00:00Z-a7f3)
     - Per D-33: advisory only, not blocking
     - Session field set on slice claim, cleared on completion
     - Other sessions can read this field to detect active work

  4. Body Sections
     - Heading: # {slice-id}: {slice-name}
     - ## Why (rationale for the slice)
     - ## What (what the slice delivers)
     - ## Dependencies (slice dependencies with status)
     - ## Contracts (contract references for the slice)
     - ## Decisions (runtime decisions made during implementation)
     - ## Files (files created or modified by the slice)

  5. File Naming
     - Per D-36: files named {slice-id}.md (e.g., S-28.md)
     - Directory: .greenlight/slices/
     - One file per slice, one writer per file at a time

  6. Examples
     - Complete example of a new slice file (status: pending)
     - Complete example of an in-progress slice file (status: implementing)
     - Complete example of a completed slice file (status: complete)

Input: None (static template, read by agents)

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | InvalidSliceId | Slice ID does not match S-{digits} or S-{digits}.{digits} pattern | Reject file creation. Error: "Invalid slice ID: {id}. Must match S-{N} or S-{N}.{N}" |
  | InvalidStatus | Status field has unrecognised value | Reject update. Error: "Invalid status: {value}. Valid values: pending, ready, tests, implementing, security, fixing, verifying, complete" |
  | InvalidFrontmatter | Frontmatter is not flat key-value between --- delimiters | Reject parse. Error: "Invalid frontmatter format. Expected flat key: value pairs between --- delimiters" |

Invariants:
  - Template is read-only at runtime (no agent writes to it)
  - Frontmatter is flat key-value only (no nesting, no arrays in values except deps which is comma-separated)
  - All field names are lowercase with underscores
  - Status enum is closed: only the 8 listed values are valid
  - Step enum is closed: only the 7 listed values are valid
  - Session field format is exactly: ISO timestamp + hyphen + random alphanumeric suffix
  - Slice ID validation prevents path traversal (S-../../etc/passwd is rejected)
  - Body sections are human-readable documentation, not machine-parsed state
  - Template matches DESIGN.md 4.2 schema exactly

Security:
  - Slice ID validation prevents path traversal attacks
  - No sensitive data in slice files (status, test counts, file lists only)
  - Session field is advisory only (no access control)

Verification: auto
Dependencies: None (this is the foundation document for the milestone)
```

### C-77: StateFormatReference

```
Contract: StateFormatReference
Boundary: Reference doc -> All commands (state detection, concurrency, backward compat)
Slice: S-28 (Slice State Foundation Documents)
Design refs: FR-2, FR-5, FR-8, FR-10, NFR-1, NFR-4, D-31, D-34, D-37, DESIGN.md 4.5, 4.6

FILE SPECIFICATION: src/references/state-format.md (~100 lines)

This is the authoritative reference document for state format detection,
migration protocol, backward compatibility rules, concurrent access
patterns, STATE.md regeneration, and advisory session tracking. All
commands reference this document before reading or writing state.

Output (mandatory sections in the reference doc):

  1. State Format Detection
     - Per D-31: check .greenlight/slices/ directory existence
     - Detection flow:
       a. If .greenlight/slices/ exists -> file-per-slice format
       b. Else if .greenlight/STATE.md exists -> legacy format
       c. Else -> no state, suggest /gl:init
     - Every command that reads state MUST follow this detection flow
     - Detection is a single directory existence check (no file parsing)

  2. File-Per-Slice Format
     - Directory: .greenlight/slices/
     - Files: {slice-id}.md (one per slice)
     - Schema: defined in templates/slice-state.md
     - Project state: .greenlight/project-state.json
     - Summary view: .greenlight/STATE.md (generated, not source of truth)

  3. Legacy Format
     - Single file: .greenlight/STATE.md (source of truth)
     - Per D-37: supported indefinitely
     - Commands work identically regardless of format
     - Migration via /gl:migrate-state (D-32: explicit, not automatic)

  4. Concurrent Access Patterns
     - Each session writes only to its own slice file
     - No shared mutable file for slice state
     - Per NFR-1: no file locking required
     - Advisory session tracking: session field in frontmatter
     - Warn before claiming a slice with active session field

  5. STATE.md Regeneration
     - Per D-34: regenerated after every state write operation
     - Generated STATE.md has header comment:
       <!-- GENERATED by greenlight -- source of truth is .greenlight/slices/*.md -->
     - Format: overview from project-state.json + slice table from all
       slice files + current section + test summary + blockers
     - Matches DESIGN.md 4.6 format exactly

  6. Crash Safety
     - Per NFR-4: write-to-temp-then-rename pattern
     - Temp files in .greenlight/slices/ (same filesystem for atomic rename)
     - POSIX rename atomicity guarantees

  7. Backward Compatibility
     - Per D-37: both formats supported indefinitely
     - Per D-38: no dual-write period
     - Detection is zero-cost (single directory check)

Input: None (static reference, read by commands and agents)

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | NoStateFound | Neither slices/ nor STATE.md exists | Return NoStateError. Commands should suggest /gl:init |
  | CorruptSliceFile | A slice file has invalid frontmatter | Skip the corrupt file. Warn: "Skipping corrupt slice file: {filename}. Run /gl:init to recreate." Include remaining valid files in aggregation |
  | SlicesDirectoryEmpty | slices/ exists but contains no .md files | Treat as file-per-slice format with zero slices. Commands can proceed (e.g., /gl:add-slice creates the first file) |

Invariants:
  - Reference doc is read-only at runtime (no agent writes to it)
  - Detection logic is deterministic: same filesystem state always produces same result
  - Legacy format support is indefinite (D-37)
  - No automatic migration (D-32)
  - STATE.md regeneration is mandatory after every write in file-per-slice mode (D-34)
  - Concurrent sessions never touch the same slice file (design guarantee)
  - Crash safety uses write-to-temp-then-rename (NFR-4)
  - Detection and format rules are the single source of truth for all commands

Security:
  - No security impact. Static reference document.
  - Slice ID validation rules referenced from C-76 apply here too.

Verification: auto
Dependencies: None (self-contained reference document)
```

---

## S-29: Init Command and State Detection

*User Actions:*
- *1. Run multiple /gl:slice sessions in parallel without state corruption (new projects use file-per-slice from day one)*

### C-78: InitSliceDirectory

```
Contract: InitSliceDirectory
Boundary: /gl:init orchestrator -> Filesystem (slices/ directory + individual slice files)
Slice: S-29 (Init Command and State Detection)
Design refs: FR-1, FR-2, D-36, DESIGN.md 4.1, 4.2, 5.3

COMMAND UPDATE: src/commands/gl/init.md — Phase 6: create slices/ and slice files

In file-per-slice mode, /gl:init creates the .greenlight/slices/
directory and writes individual slice state files instead of the
monolithic STATE.md. Each slice defined in GRAPH.json gets its own
file in .greenlight/slices/{id}.md.

Behaviour:

  1. Create .greenlight/slices/ directory (permissions 0o755)

  2. For each slice in GRAPH.json:
     a. Validate slice ID matches S-{digits} pattern
     b. Create .greenlight/slices/{id}.md with:
        - Frontmatter: id, status (pending), step (none), milestone,
          started (empty), updated (current ISO timestamp), tests (0),
          security_tests (0), session (empty), deps (from GRAPH.json)
        - Body: heading with slice name, empty body sections
     c. Use write-to-temp-then-rename for crash safety (NFR-4)

  3. Generate .greenlight/STATE.md as summary view:
     - Header comment: <!-- GENERATED by greenlight -->
     - Computed from all slice files per DESIGN.md 4.6 format

Input:
  - GRAPH.json (slice definitions with IDs, names, dependencies)
  - Milestone name from design phase

Output:
  - .greenlight/slices/ directory created
  - One .md file per slice in .greenlight/slices/
  - .greenlight/STATE.md generated as summary view

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | SlicesDirExists | .greenlight/slices/ already exists | Warn user. Offer to overwrite or skip. Do not silently overwrite |
  | InvalidSliceIdInGraph | GRAPH.json contains slice ID not matching S-{digits} | Skip invalid slice. Warn: "Skipping invalid slice ID: {id}" |
  | DirectoryCreateFailure | Cannot create .greenlight/slices/ directory | Report error. Abort init. Suggest checking permissions |
  | FileWriteFailure | Cannot write a slice file | Report error for that file. Continue with remaining slices. Warn user of partial state |

Invariants:
  - Directory permissions are 0o755, file permissions are 0o644
  - Each slice file follows the schema in templates/slice-state.md (C-76)
  - Slice IDs are validated before file creation (path traversal prevention)
  - All slice files are created atomically (write-to-temp-then-rename)
  - STATE.md is generated after all slice files are written
  - Existing /gl:init behaviour for non-state operations is unchanged
  - If /gl:init was already run with legacy format, this does NOT auto-migrate
    (user must run /gl:migrate-state per D-32)

Security:
  - Slice ID validation prevents path traversal (C-76 rules)
  - File permissions follow existing conventions
  - No sensitive data in created files

Verification: verify
Acceptance Criteria:
- After /gl:init, .greenlight/slices/ directory exists with one .md file per slice from GRAPH.json
- Each slice file has valid frontmatter with status pending and step none
- STATE.md is generated with the GENERATED header comment and a table of all slices
- Existing /gl:init functionality (config, CONTRACTS.md, etc.) is unchanged

Steps:
- Run /gl:init on a new project with a GRAPH.json containing 3+ slices
- Verify .greenlight/slices/ directory exists
- Open any slice file and verify frontmatter matches templates/slice-state.md schema
- Open STATE.md and verify it has the GENERATED comment header

Dependencies: C-76 (slice state schema must be defined), C-77 (state format rules must be defined)
```

### C-79: InitProjectState

```
Contract: InitProjectState
Boundary: /gl:init orchestrator -> Filesystem (project-state.json)
Slice: S-29 (Init Command and State Detection)
Design refs: FR-4, D-35, DESIGN.md 4.4

COMMAND UPDATE: src/commands/gl/init.md — Create project-state.json

/gl:init creates .greenlight/project-state.json to store non-slice state:
project overview, session metadata, and active blockers.

Behaviour:

  1. Create .greenlight/project-state.json with:
     ```json
     {
       "overview": {
         "value_prop": "{from design phase}",
         "stack": "{from design phase}",
         "mode": "{from config.json or default yolo}"
       },
       "session": {
         "last_session": "{current ISO timestamp}",
         "resume_file": null
       },
       "blockers": []
     }
     ```

  2. Overview fields populated from design phase context
  3. Session initialised with current timestamp
  4. Blockers array starts empty

Input:
  - Design phase context (value prop, stack, mode)
  - Current timestamp

Output:
  - .greenlight/project-state.json created

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | ProjectStateExists | .greenlight/project-state.json already exists | Overwrite with fresh state. Warn user: "Overwriting existing project-state.json" |
  | WriteFailure | Cannot write project-state.json | Report error. Abort. Suggest checking permissions |
  | MissingDesignContext | Design phase did not provide value_prop or stack | Use placeholder values: value_prop "TBD", stack "TBD". Warn user |

Invariants:
  - project-state.json follows DESIGN.md 4.4 schema exactly
  - JSON is valid and parseable
  - overview.mode defaults to "yolo" if not specified
  - session.last_session is always a valid ISO timestamp
  - blockers is always an array (never null)
  - File permissions are 0o644

Security:
  - No sensitive data in project-state.json
  - No secrets, tokens, or PII

Verification: auto
Dependencies: None
```

### C-80: InitStateDetection

```
Contract: InitStateDetection
Boundary: All commands -> State detection logic (format resolution)
Slice: S-29 (Init Command and State Detection)
Design refs: FR-5, FR-8, D-31, DESIGN.md 4.5

STATE DETECTION: Logic documented in references/state-format.md,
implemented by every command that reads state

Every state-reading command must follow the detection flow before
accessing state. This contract defines the detection behaviour that
all commands implement consistently.

Behaviour:

  1. Check if .greenlight/slices/ directory exists:
     - Yes -> file-per-slice format. Read individual slice files.
     - No -> continue to step 2.

  2. Check if .greenlight/STATE.md exists:
     - Yes -> legacy format. Parse STATE.md as source of truth.
     - No -> continue to step 3.

  3. No state found:
     - Return NoStateError
     - Commands should display: "No project state found. Run /gl:init to get started."

  Commands that implement this detection:
  - /gl:slice (pre-flight, Step 4, Step 10)
  - /gl:status (full read)
  - /gl:pause (current slice lookup)
  - /gl:resume (full state read)
  - /gl:ship (pre-check)
  - /gl:add-slice (slice list)
  - /gl:quick (test summary)

Input:
  - Filesystem state (.greenlight/ directory contents)

Output:
  - format: "file-per-slice" | "legacy" | "none"
  - state_data: parsed state (format-dependent)

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | NoStateFound | Neither slices/ nor STATE.md exists | Display suggestion to run /gl:init. Do not crash |
  | SlicesDirUnreadable | slices/ exists but cannot be read | Report permission error. Suggest checking directory permissions |
  | LegacyStateMalformed | STATE.md exists but cannot be parsed | Report parse error. Suggest running /gl:migrate-state or re-running /gl:init |

Invariants:
  - Detection logic is identical across all commands (single source of truth in C-77)
  - Directory existence check is the primary signal (D-31)
  - Detection always succeeds or produces a clear error (no ambiguous states)
  - Detection does not modify any files (read-only operation)
  - Detection result determines which read/write paths the command follows
  - Legacy format commands work exactly as before (no regression)
  - File-per-slice commands read individual files instead of monolithic STATE.md

Security:
  - Detection is read-only. No files modified.
  - Directory existence check does not follow symlinks (prevents symlink attacks)

Verification: verify
Acceptance Criteria:
- Commands detect file-per-slice format when .greenlight/slices/ directory exists
- Commands detect legacy format when only .greenlight/STATE.md exists
- Commands display helpful message when no state exists
- Existing commands work identically with legacy STATE.md format (no regression)

Steps:
- Create a project with file-per-slice format (/gl:init). Run /gl:status and verify it reads from slices/
- Create a project with legacy STATE.md. Run /gl:status and verify it reads from STATE.md
- Remove both slices/ and STATE.md. Run /gl:status and verify it suggests /gl:init

Dependencies: C-77 (detection logic defined in state-format.md)
```

---

## S-30: Slice Command State Write

*User Actions:*
- *1. Run multiple /gl:slice sessions in parallel without state corruption*

### C-81: SliceCommandStateWrite

```
Contract: SliceCommandStateWrite
Boundary: /gl:slice orchestrator -> Filesystem (own slice file + STATE.md regeneration)
Slice: S-30 (Slice Command State Write)
Design refs: FR-1, FR-2, FR-10, D-34, DESIGN.md 4.8, 5.2

COMMAND UPDATE: src/commands/gl/slice.md — Read/write own slice file, regenerate STATE.md

/gl:slice is the primary command affected by this change. It must:
1. Detect state format (C-80)
2. Read from its own slice file (not monolithic STATE.md)
3. Write only to its own slice file (not shared STATE.md)
4. Regenerate STATE.md as summary view after every write

Changes to /gl:slice pipeline:

  Pre-flight (before Step 1):
    - Detect state format (C-80)
    - If file-per-slice: read target slice file from .greenlight/slices/{id}.md
    - If legacy: read STATE.md as before (no change)

  Step 4 (claim slice / start implementation):
    - If file-per-slice: write status + step + session to own slice file
    - Set session field (C-82)
    - Regenerate STATE.md (D-34)
    - If legacy: write to STATE.md as before (no change)

  Step 10 (completion / status update):
    - If file-per-slice: write final status + test counts to own slice file
    - Clear session field
    - Regenerate STATE.md (D-34)
    - If legacy: write to STATE.md as before (no change)

  STATE.md regeneration:
    - Read all slice files from .greenlight/slices/
    - Read project-state.json for overview, session, blockers
    - Write STATE.md in generated format (DESIGN.md 4.6)
    - Header: <!-- GENERATED by greenlight -- source of truth is .greenlight/slices/*.md -->

  Crash safety:
    - Slice file writes use write-to-temp-then-rename (NFR-4)
    - STATE.md regeneration uses write-to-temp-then-rename
    - If regeneration fails, slice file is still correct (partial failure is safe)

Input:
  - Slice ID (from orchestrator context)
  - State format (from detection)
  - Slice state data (status, step, tests, session, etc.)

Output:
  - Updated .greenlight/slices/{id}.md
  - Regenerated .greenlight/STATE.md
  - Or: updated legacy .greenlight/STATE.md (if legacy format)

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | SliceFileNotFound | Target slice file does not exist in slices/ | Create it (first write). Warn: "Creating new slice file for {id}" |
  | SliceFileWriteFailure | Cannot write to slice file | Report error. Do not proceed to next step. Suggest checking permissions |
  | RegenerationFailure | Cannot regenerate STATE.md | Warn but continue. Slice file is still correct. STATE.md will be regenerated on next write |
  | ConcurrentSliceClaim | Target slice has active session field from another session | Warn: "Slice {id} appears to be in progress (session: {session}). Continue anyway? (y/n)" |

Invariants:
  - /gl:slice writes ONLY to its own slice file (never to another slice's file)
  - STATE.md is regenerated after every state write (D-34)
  - Legacy format behaviour is completely unchanged (D-37)
  - Crash safety via write-to-temp-then-rename on all writes (NFR-4)
  - Concurrent sessions writing to different slices never conflict (FR-2)
  - Session field is advisory only -- warning, not blocking (D-33)
  - Slice file is the source of truth, STATE.md is a convenience view
  - Pipeline step ordering is unchanged (only read/write targets change)

Security:
  - Slice ID validated before file path construction (path traversal prevention)
  - Advisory session warning prevents accidental duplicate work
  - No sensitive data written to slice files

Verification: verify
Acceptance Criteria:
- Running /gl:slice on a slice writes status updates to .greenlight/slices/{id}.md only
- Running two /gl:slice sessions on different slices simultaneously does not corrupt either slice's state
- STATE.md is regenerated after each write and reflects all slice states accurately
- Running /gl:slice on a legacy-format project works exactly as before

Steps:
- Run /gl:slice on a file-per-slice project. After Step 4, verify .greenlight/slices/{id}.md has status implementing
- Open STATE.md and verify it shows the slice as in-progress
- In a second terminal, run /gl:slice on a different slice. Verify neither slice file is corrupted
- Run /gl:slice on a legacy project and verify STATE.md is written directly (no slices/ involved)

Dependencies: C-76 (slice file schema), C-77 (state format rules), C-80 (state detection)
```

### C-82: SliceSessionTracking

```
Contract: SliceSessionTracking
Boundary: /gl:slice orchestrator -> Slice file frontmatter (advisory session field)
Slice: S-30 (Slice Command State Write)
Design refs: FR-9, D-33, DESIGN.md 4.3

SESSION TRACKING: Advisory field in slice file frontmatter

When a session starts working on a slice, it writes a session identifier
to the slice file's frontmatter. This allows other sessions to detect
active work and warn before claiming the same slice.

Behaviour:

  1. On slice claim (Step 4):
     a. Generate session ID: {ISO timestamp}-{random 4-char hex suffix}
        Example: 2026-02-22T14:00:00Z-a7f3
     b. Write session field to slice file frontmatter
     c. If slice already has a session field with a value:
        - Warn: "Slice {id} appears to be in progress (session: {session}).
          Another terminal may be working on this slice."
        - Prompt: "Continue anyway? (y/n)"
        - If yes: overwrite session field with new session ID
        - If no: abort slice execution

  2. On slice completion (Step 10):
     a. Clear session field (set to empty string)
     b. Write updated frontmatter to slice file

  3. On slice pause:
     a. Session field is preserved (slice is still claimed)
     b. /gl:pause writes resume context to project-state.json

  4. On unexpected exit (crash/ctrl-c):
     a. Session field remains in frontmatter (stale session)
     b. Next session to claim this slice will see the stale session
        and receive a warning (advisory only, not blocking)

Input:
  - Current ISO timestamp
  - Random hex suffix (4 characters)
  - Existing session field value (if any)

Output:
  - session field written to slice file frontmatter
  - Warning displayed to user (if existing session detected)

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | StaleSession | Session field has value but no active process | Warning only. User can override. No automatic cleanup |
  | SessionGenerationFailure | Cannot generate random suffix | Use timestamp-only as session ID. Warn: "Session ID generated without random suffix" |

Invariants:
  - Session tracking is advisory only (D-33) -- never blocking
  - Session ID format is exactly: ISO-timestamp-{4-char-hex}
  - Session field is set BEFORE any agent work begins (earliest possible signal)
  - Session field is cleared ONLY on successful completion
  - Stale sessions are warned about but never automatically cleaned up
  - Session tracking only applies to file-per-slice format (not legacy)
  - Other sessions can read session fields from any slice file in .greenlight/slices/
  - Warning prompt is skipped in yolo mode (auto-continue with warning log)

Security:
  - Session ID contains no sensitive information (timestamp + random hex)
  - Session tracking does not provide access control (advisory only)
  - No PID, username, or machine identifier in session ID

Verification: auto
Dependencies: C-76 (session field defined in slice state schema), C-81 (slice writes include session updates)
```

---

## S-31: Supporting Command Updates

*User Actions:*
- *2. See accurate slice status across all concurrent sessions*

### C-83: StatusSliceAggregation

```
Contract: StatusSliceAggregation
Boundary: /gl:status command -> Filesystem (read all slice files, compute summary)
Slice: S-31 (Supporting Command Updates)
Design refs: FR-3, FR-10, D-34, DESIGN.md 4.6, 4.8

COMMAND UPDATE: src/commands/gl/status.md — Read from slices/ directory, compute summary

/gl:status reads all slice files from .greenlight/slices/ and computes
a summary view. In file-per-slice mode, it reads individual files
instead of parsing the monolithic STATE.md.

Behaviour:

  1. Detect state format (C-80)

  2. If file-per-slice:
     a. Read all .md files from .greenlight/slices/
     b. Parse frontmatter from each file (flat key-value)
     c. Compute summary:
        - Slice table: ID, Name, Status, Tests, Security, Deps
        - Progress: done/total slices (complete count / total count)
        - Current: slices where status is not pending and not complete
        - Test Summary: sum of tests and security_tests from all files
     d. Read project-state.json for overview, session, blockers
     e. Display computed summary to user
     f. Regenerate STATE.md (D-34)

  3. If legacy:
     a. Read and display STATE.md as before (no change)

Input:
  - .greenlight/slices/*.md (all slice files)
  - .greenlight/project-state.json (overview, session, blockers)

Output:
  - Summary display to user
  - Regenerated .greenlight/STATE.md (in file-per-slice mode)

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | NoSliceFiles | slices/ exists but has no .md files | Display empty summary: "No slices found. Run /gl:init or /gl:add-slice" |
  | CorruptSliceFile | A slice file has invalid frontmatter | Skip that file. Warn: "Skipping corrupt file: {filename}". Include remaining files in summary |
  | ProjectStateReadFailure | Cannot read project-state.json | Display slice summary without overview/session/blockers. Warn user |

Invariants:
  - Status display is computed from individual slice files (never cached)
  - All slice files are read on every /gl:status invocation (fresh data)
  - Performance: reading 50+ files completes in under 1 second (NFR-2)
  - Corrupt files are skipped, not blocking (graceful degradation)
  - Legacy format display is completely unchanged
  - STATE.md regeneration happens after display (even if display fails)
  - Slice table is sorted by slice ID (ascending)

Security:
  - Read-only operation. No files modified except STATE.md regeneration.
  - No sensitive data displayed.

Verification: verify
Acceptance Criteria:
- /gl:status displays a summary table computed from all files in .greenlight/slices/
- Status reflects real-time slice state (each file read fresh on invocation)
- /gl:status on a legacy project displays STATE.md as before (no regression)
- Corrupt slice files are skipped with a warning, not blocking

Steps:
- Run /gl:status on a file-per-slice project with multiple slices in various states
- Verify the summary table matches the frontmatter in individual slice files
- Manually corrupt one slice file's frontmatter. Run /gl:status. Verify warning appears and other slices display correctly

Dependencies: C-77 (state format rules), C-80 (state detection)
```

### C-84: SupportingCommandStateAdaptation

```
Contract: SupportingCommandStateAdaptation
Boundary: Supporting commands -> State detection + slice file reads/writes
Slice: S-31 (Supporting Command Updates)
Design refs: FR-8, D-34, DESIGN.md 4.8, 5.2

COMMAND UPDATES: 5 supporting commands adapted for file-per-slice format

Each supporting command gains state detection (C-80) and format-aware
read/write behaviour. All changes follow the same pattern: detect
format, use appropriate read/write paths, regenerate STATE.md after writes.

Commands updated:

  1. /gl:pause (src/commands/gl/pause.md)
     - Detect format. If file-per-slice:
       - Write pause state to own slice file (status unchanged, session preserved)
       - Write resume context to project-state.json (session.resume_file)
       - Regenerate STATE.md
     - If legacy: write to STATE.md as before

  2. /gl:resume (src/commands/gl/resume.md)
     - Detect format. If file-per-slice:
       - Read all slice files to determine resumable state
       - Read project-state.json for resume context
       - Resume from slice file + project-state.json
     - If legacy: read STATE.md as before

  3. /gl:ship (src/commands/gl/ship.md)
     - Detect format. If file-per-slice:
       - Read all slice files. Pre-check: all must have status complete
       - If any non-complete slice found: report which slices are incomplete
     - If legacy: read STATE.md as before

  4. /gl:add-slice (src/commands/gl/add-slice.md)
     - Detect format. If file-per-slice:
       - Create new slice file in .greenlight/slices/{id}.md
       - Regenerate STATE.md
     - If legacy: update STATE.md as before

  5. /gl:quick (src/commands/gl/quick.md)
     - Detect format. If file-per-slice:
       - Update test counts in relevant slice file
       - Regenerate STATE.md
     - If legacy: update STATE.md as before

Input (per command):
  - State format (from detection)
  - Command-specific data (slice ID, test counts, etc.)

Output (per command):
  - Updated slice file(s) and/or project-state.json
  - Regenerated STATE.md (for write operations)
  - Or: updated legacy STATE.md (if legacy format)

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | FormatDetectionFailure | Cannot determine state format | Report error. Suggest running /gl:init |
  | SliceFileNotFound | Referenced slice file does not exist | Create it if write operation, report missing if read operation |
  | ProjectStateNotFound | project-state.json missing in file-per-slice mode | Create default project-state.json. Warn user |
  | RegenerationFailure | STATE.md regeneration fails | Warn but continue. Slice files are still correct |

Invariants:
  - Every command follows the same detection flow (C-80)
  - Write commands regenerate STATE.md after writes (D-34)
  - Legacy format behaviour is completely unchanged for all 5 commands
  - Each command writes only to its own concern (no cross-slice writes)
  - /gl:add-slice validates new slice ID before creating file (C-76 rules)
  - /gl:ship reads ALL slice files (pre-check is comprehensive)
  - /gl:quick updates only the relevant slice file (not all files)
  - All write operations use write-to-temp-then-rename (NFR-4)
  - Command behaviour is identical regardless of format (same user experience)

Security:
  - Slice ID validation on /gl:add-slice prevents path traversal
  - No sensitive data in any state files
  - All commands are non-destructive (no data loss on format switch)

Verification: verify
Acceptance Criteria:
- /gl:pause on a file-per-slice project writes pause context to project-state.json, not STATE.md
- /gl:resume on a file-per-slice project reads from slice files and project-state.json
- /gl:ship on a file-per-slice project checks all individual slice files for completeness
- /gl:add-slice creates a new file in .greenlight/slices/ and regenerates STATE.md
- /gl:quick updates test counts in the relevant slice file and regenerates STATE.md
- All 5 commands work exactly as before on legacy-format projects

Steps:
- Run /gl:add-slice on a file-per-slice project. Verify new file created in slices/ and STATE.md updated
- Run /gl:ship on a file-per-slice project with one incomplete slice. Verify it reports the incomplete slice

Dependencies: C-77 (state format rules), C-78 (slices/ directory exists), C-80 (state detection)
```

---

## S-32: Migration Command

*User Actions:*
- *3. Migrate existing projects to the new state format*

### C-85: MigrateStateCommand

```
Contract: MigrateStateCommand
Boundary: User -> /gl:migrate-state command (legacy STATE.md to file-per-slice conversion)
Slice: S-32 (Migration Command)
Design refs: FR-6, D-32, D-38, DESIGN.md 4.7, 5.1

COMMAND SPECIFICATION: src/commands/gl/migrate-state.md (~80 lines)

/gl:migrate-state converts an existing STATE.md-based project to the
file-per-slice format. Migration is explicit (D-32), one-way, and
all-or-nothing. No dual-write period (D-38).

Behaviour (migration flow per DESIGN.md 4.7):

  1. Verify .greenlight/STATE.md exists
     - If not: error "No STATE.md found. Nothing to migrate."

  2. Verify .greenlight/slices/ does NOT exist
     - If exists: error "Already using file-per-slice format. Nothing to migrate."

  3. Parse STATE.md:
     a. Extract slice table rows (ID, Name, Status, Tests, Security, Deps)
     b. Extract Current section (active slice, step)
     c. Extract Decisions section
     d. Extract Blockers section
     e. Extract Session section
     f. Extract Overview section (value prop, stack, mode)

  4. Create .greenlight/slices/ directory (0o755)

  5. For each slice row:
     a. Validate slice ID format (S-{digits})
     b. Create .greenlight/slices/{id}.md with:
        - Frontmatter from table data (id, status, tests, security_tests, deps)
        - step: from Current section if this is the active slice, else "none"
        - milestone: from project context or "unknown"
        - started: inferred from status (non-pending -> current date)
        - updated: current ISO timestamp
        - session: empty (no active session after migration)
     c. Body: heading with slice name, minimal body sections
     d. Use write-to-temp-then-rename (NFR-4)

  6. Create .greenlight/project-state.json from non-slice sections:
     - overview from Overview section
     - session from Session section
     - blockers from Blockers section

  7. Rename .greenlight/STATE.md to .greenlight/STATE.md.backup

  8. Generate new .greenlight/STATE.md (generated format with header comment)

  9. Report: "Migrated {N} slices to file-per-slice format. Backup: STATE.md.backup"

Input:
  - .greenlight/STATE.md (legacy format, source of truth)

Output:
  - .greenlight/slices/ directory with individual slice files
  - .greenlight/project-state.json
  - .greenlight/STATE.md.backup (original STATE.md)
  - .greenlight/STATE.md (regenerated, generated format)
  - Summary report to user

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | NoStateMd | .greenlight/STATE.md does not exist | Error: "No STATE.md found. Nothing to migrate." Stop |
  | AlreadyMigrated | .greenlight/slices/ directory already exists | Error: "Already using file-per-slice format. Nothing to migrate." Stop |
  | ParseFailure | Cannot parse STATE.md table or sections | Error: "Failed to parse STATE.md: {detail}. Manual migration may be needed." Stop |
  | InvalidSliceId | Parsed slice ID does not match expected pattern | Skip that slice. Warn: "Skipping invalid slice ID: {id}" |
  | PartialWriteFailure | Some slice files fail to write | Abort migration. Remove partially created slices/ directory. Restore STATE.md from backup. Error: "Migration failed. Original STATE.md preserved." |
  | BackupRenameFailure | Cannot rename STATE.md to backup | Abort migration. Remove slices/ directory. Error: "Cannot create backup. Migration aborted." |

Invariants:
  - Migration is one-way (file-per-slice to legacy is not supported)
  - Migration is all-or-nothing (partial migration is rolled back)
  - Original STATE.md is ALWAYS preserved as STATE.md.backup
  - Migration is explicit only (D-32: no auto-migration on access)
  - No dual-write period (D-38)
  - Slice IDs are validated before file creation (path traversal prevention)
  - Migration creates the same state as /gl:init would for the same slices
  - After migration, all commands use file-per-slice format automatically (D-31 detection)
  - Migration can be run at any time (not tied to a specific command or phase)

Security:
  - Slice ID validation prevents path traversal during file creation
  - Backup preserves the original STATE.md (no data loss)
  - Migration does not modify any source code, test files, or config files
  - File permissions follow conventions: directories 0o755, files 0o644

Verification: verify
Acceptance Criteria:
- Running /gl:migrate-state on a legacy project creates .greenlight/slices/ with one file per slice
- STATE.md.backup contains the original STATE.md content exactly
- New STATE.md has the GENERATED header comment
- project-state.json contains overview, session, and blockers from original STATE.md
- After migration, /gl:status reads from slices/ (not the old STATE.md)

Steps:
- Create a project with legacy STATE.md (or use an existing one)
- Run /gl:migrate-state
- Verify .greenlight/slices/ contains files matching the slice table from original STATE.md
- Verify STATE.md.backup exists and matches the original
- Run /gl:status and verify it displays data from the new format
- Verify /gl:migrate-state on an already-migrated project reports "Already using file-per-slice format"

Dependencies: C-76 (slice file schema), C-77 (state format rules), C-79 (project-state.json schema)
```

### C-86: MigrateStateBackup

```
Contract: MigrateStateBackup
Boundary: /gl:migrate-state -> Filesystem (backup creation + atomic directory creation)
Slice: S-32 (Migration Command)
Design refs: FR-6, NFR-4, DESIGN.md 4.7 steps 7-8

BACKUP AND ATOMICITY: Migration safety guarantees

The migration creates a backup before any destructive operation and
ensures the slices/ directory is created atomically.

Behaviour:

  1. All slice files are written to .greenlight/slices/ first
  2. project-state.json is written
  3. Only AFTER all files are successfully written:
     a. Rename STATE.md to STATE.md.backup (os.Rename, atomic on POSIX)
     b. Generate new STATE.md in generated format
  4. If any step in 1-2 fails:
     a. Remove the partially created slices/ directory and all contents
     b. STATE.md is untouched (no rename occurred yet)
     c. Report failure to user

  Cleanup on failure:
  - os.RemoveAll(".greenlight/slices/") removes partial directory
  - STATE.md was not renamed yet, so original is preserved
  - project-state.json removal (if partially written)

Input:
  - Parsed STATE.md data (from C-85)
  - Filesystem operations (directory create, file write, rename)

Output:
  - Atomic migration: either fully complete or fully rolled back
  - STATE.md.backup preserved

Errors:
  | Error State | When | Behaviour |
  |-------------|------|-----------|
  | CleanupFailure | Cannot remove partially created slices/ directory | Report error. Manual cleanup instructions: "Remove .greenlight/slices/ manually. Your STATE.md is preserved." |
  | BackupExists | STATE.md.backup already exists | Rename to STATE.md.backup.{timestamp}. Warn user |

Invariants:
  - Backup is created AFTER slice files are written (ensures backup is only needed if migration succeeded)
  - STATE.md rename (to backup) happens LAST (ordering guarantees rollback safety)
  - If migration fails at any point before rename, STATE.md is untouched
  - Multiple backups are preserved (timestamped suffixes) -- no data loss
  - write-to-temp-then-rename for individual slice files (NFR-4)
  - Entire migration is designed for crash safety: the worst case is a stale slices/ directory alongside an intact STATE.md

Security:
  - Backup file contains no sensitive data (same content as STATE.md)
  - No permission escalation during migration
  - Cleanup removes only files in .greenlight/ (no traversal)

Verification: auto
Dependencies: C-85 (migration command defines the flow that this contract ensures is safe)
```

---

## S-33: Documentation Updates

*User Actions:*
- *1. Run multiple /gl:slice sessions in parallel without state corruption (documentation confirms)*
- *2. See accurate slice status across all concurrent sessions (documentation confirms)*

### C-87: CLAUDEmdStateFormatRule

```
Contract: CLAUDEmdStateFormatRule
Boundary: CLAUDE.md -> All agents (state format awareness hard rule in standards)
Slice: S-33 (Documentation Updates)
Design refs: DESIGN.md 5.2 (src/CLAUDE.md change)

FILE UPDATE: src/CLAUDE.md

Location: Add state format awareness. Insert within an appropriate
section (after existing rules, before "What NOT To Do").

Content (approximately 5 lines, hard rule):
  ### State Format
  - Check `.greenlight/slices/` before reading STATE.md directly
  - If `slices/` exists, individual slice files are the source of truth
  - STATE.md is generated output in file-per-slice mode -- do not write to it directly
  - Full protocol: `references/state-format.md`

Errors: None (static content update)

Invariants:
  - Rule is concise (header + 4 bullet points)
  - Rule references the full protocol in references/state-format.md
  - Existing CLAUDE.md sections unchanged
  - This is a hard rule: agents MUST check slices/ before reading STATE.md
  - Rule applies to all agents, not just specific ones

Security:
  - No security impact. Static content update.

Verification: auto
Dependencies: C-77 (references/state-format.md must be defined)
```

### C-88: StateTemplateDocUpdate

```
Contract: StateTemplateDocUpdate
Boundary: templates/state.md -> Both formats documented
Slice: S-33 (Documentation Updates)
Design refs: DESIGN.md 5.2 (src/templates/state.md change)

FILE UPDATE: src/templates/state.md — Document both state formats

The existing templates/state.md documents the STATE.md format. This
update adds documentation for:
1. The file-per-slice format as the recommended format for new projects
2. The generated nature of STATE.md in file-per-slice mode
3. Migration instructions for existing projects
4. Cross-reference to references/state-format.md and templates/slice-state.md

Content additions (~30 lines):

  Section: State Format Detection
    - Explain the two formats (file-per-slice and legacy)
    - Explain how commands detect which format is active

  Section: Generated STATE.md
    - Explain that STATE.md is generated output in file-per-slice mode
    - Show the generated header comment
    - Explain that changes to STATE.md will be overwritten

  Section: Migration
    - Reference /gl:migrate-state command
    - Note: migration is explicit, one-way, with backup

Errors: None (static content update)

Invariants:
  - Existing STATE.md template content is preserved (additive update)
  - Both formats are documented side by side
  - Migration instructions reference /gl:migrate-state (not manual steps)
  - Cross-references to state-format.md and slice-state.md are included
  - New projects are directed to file-per-slice format

Security:
  - No security impact. Documentation update.

Verification: auto
Dependencies: C-76 (slice state template must be defined), C-77 (state format reference must be defined)
```

### C-89: CheckpointProtocolStateUpdate

```
Contract: CheckpointProtocolStateUpdate
Boundary: checkpoint-protocol.md -> Slice file references for state context
Slice: S-33 (Documentation Updates)
Design refs: DESIGN.md 5.2 (src/references/checkpoint-protocol.md change)

FILE UPDATE: src/references/checkpoint-protocol.md — Reference slice files for state context

The checkpoint protocol references STATE.md for state context during
checkpoint save/restore. This update adds awareness that state context
may come from individual slice files instead.

Content changes (~10 lines):

  Update state context references:
    - Where checkpoint save reads "current slice" from STATE.md,
      add: "In file-per-slice mode, read from .greenlight/slices/{id}.md"
    - Where checkpoint restore writes state context,
      add: "In file-per-slice mode, write to individual slice file"
    - Add note: "State format detection (references/state-format.md)
      determines which path to use"

Errors: None (static content update)

Invariants:
  - Existing checkpoint protocol logic is unchanged
  - Additions are conditional on state format (additive, not replacing)
  - Cross-reference to state-format.md for detection logic
  - Checkpoint save/restore works with both formats

Security:
  - No security impact. Documentation update.

Verification: auto
Dependencies: C-77 (state format reference must be defined)
```

---

## S-34: Manifest and Integration

*User Actions:*
- *1. Run multiple /gl:slice sessions in parallel without state corruption (files installed)*
- *2. See accurate slice status across all concurrent sessions (files installed)*
- *3. Migrate existing projects to the new state format (command installed)*

### C-90: ManifestParallelStateUpdate

```go
// Contract: ManifestParallelStateUpdate
// Boundary: Go CLI -> Manifest (3 new file paths for parallel state)
// Slice: S-34 (Manifest and Integration)
//
// FILE: internal/installer/installer.go
//
// Change: Add 3 new entries to Manifest slice
//
// New entries (inserted in alphabetical order within their sections):
//   "commands/gl/migrate-state.md"    // NEW -- migration command
//   "references/state-format.md"      // NEW -- state format protocol
//   "templates/slice-state.md"        // NEW -- slice state template
//
// Updated Manifest (38 entries, up from 35):
//   "agents/gl-architect.md"
//   "agents/gl-assessor.md"
//   "agents/gl-codebase-mapper.md"
//   "agents/gl-debugger.md"
//   "agents/gl-designer.md"
//   "agents/gl-implementer.md"
//   "agents/gl-security.md"
//   "agents/gl-test-writer.md"
//   "agents/gl-verifier.md"
//   "agents/gl-wrapper.md"
//   "commands/gl/add-slice.md"
//   "commands/gl/assess.md"
//   "commands/gl/changelog.md"
//   "commands/gl/debug.md"
//   "commands/gl/design.md"
//   "commands/gl/help.md"
//   "commands/gl/init.md"
//   "commands/gl/map.md"
//   "commands/gl/migrate-state.md"    <-- NEW
//   "commands/gl/pause.md"
//   "commands/gl/quick.md"
//   "commands/gl/resume.md"
//   "commands/gl/roadmap.md"
//   "commands/gl/settings.md"
//   "commands/gl/ship.md"
//   "commands/gl/slice.md"
//   "commands/gl/status.md"
//   "commands/gl/wrap.md"
//   "references/checkpoint-protocol.md"
//   "references/circuit-breaker.md"
//   "references/deviation-rules.md"
//   "references/state-format.md"      <-- NEW
//   "references/verification-patterns.md"
//   "references/verification-tiers.md"
//   "templates/config.md"
//   "templates/slice-state.md"        <-- NEW
//   "templates/state.md"
//   "CLAUDE.md"
//
// Errors: none (compile-time constant)
//
// Invariants:
// - CLAUDE.md remains the LAST entry
// - Entries within each section (agents/, commands/gl/, references/, templates/)
//   are alphabetically ordered
// - go:embed directive in main.go already uses wildcards
//   (src/templates/*.md, src/references/*.md, src/commands/gl/*.md)
//   so new .md files are automatically embedded -- no main.go change needed
// - Manifest count increases from 35 to 38
// - All existing tests that validate manifest count must be updated to expect 38
// - This change is additive to C-74 (previous manifest update)
//
// Verification: auto
// Dependencies: C-74 (previous manifest update must be applied first or simultaneously)
```

---

---

## Milestone: CLI Orchestrator

> **Scope:** Extend the Go CLI binary into a full orchestrator with parallel execution
> **Stack:** Go 1.24, stdlib only
> **Date:** 2026-02-23
> **Slices:** S-35 through S-46
> **Contracts:** C-91 through C-116

---

## S-35: Frontmatter Parser

*User Actions:*
- *Supports all user actions (foundation for reading slice state files from Go)*

### C-91: FrontmatterParse

```go
// Contract: FrontmatterParse
// Boundary: Frontmatter parser -> Slice state files (parse flat YAML frontmatter)
// Slice: S-35 (Frontmatter Parser)
// Design refs: FR-9, D-43, DESIGN.md 4.2
//
// FILE: internal/frontmatter/frontmatter.go
//
// Package frontmatter provides a simple key-value parser for flat YAML
// frontmatter delimited by "---" lines. No nesting. No arrays. No YAML
// library. Line-by-line string splitting only (D-43).
//
// Signature:
//   func Parse(content string) (map[string]string, string, error)
//
// Input:
//   content string  // full file content including "---" delimiters
//
// Output:
//   fields  map[string]string  // key-value pairs from frontmatter
//   body    string             // everything after the closing "---"
//   err     error              // nil on success
//
// Behaviour:
//   1. Find opening "---" (must be first non-empty line)
//   2. Find closing "---"
//   3. Between delimiters, split each line on first ":"
//   4. Trim whitespace from keys and values
//   5. Return map of key-value pairs and remaining body
//
// Errors:
//   | Error | When |
//   |-------|------|
//   | ErrNoFrontmatter | Content has no opening "---" delimiter |
//   | ErrUnclosedFrontmatter | Opening "---" found but no closing "---" |
//   | ErrInvalidLine | Line between delimiters has no ":" separator |
//
// Invariants:
//   - Empty values are valid (key with no value after colon)
//   - Keys are trimmed of whitespace
//   - Values are trimmed of whitespace
//   - Body preserves original formatting (no trimming)
//   - Empty content returns ErrNoFrontmatter
//   - Content with only "---\n---" returns empty map and empty body
//   - Values containing ":" are preserved (split on first ":" only)
//   - Lines containing only whitespace between delimiters are skipped
//
// Verification: auto
// Dependencies: none
```

### C-92: FrontmatterWrite

```go
// Contract: FrontmatterWrite
// Boundary: Frontmatter writer -> Slice state files (write flat YAML frontmatter)
// Slice: S-35 (Frontmatter Parser)
// Design refs: FR-9, D-43, DESIGN.md 4.2
//
// FILE: internal/frontmatter/frontmatter.go
//
// Signature:
//   func Write(fields map[string]string, body string) string
//
// Input:
//   fields  map[string]string  // key-value pairs for frontmatter
//   body    string             // content after the closing "---"
//
// Output:
//   content string  // complete file content with frontmatter + body
//
// Behaviour:
//   1. Write opening "---\n"
//   2. Write each key-value pair as "key: value\n" in sorted key order
//   3. Write closing "---\n"
//   4. Append body
//
// Errors: none (pure function, always succeeds)
//
// Invariants:
//   - Keys are written in sorted (alphabetical) order for deterministic output
//   - Output of Write can be parsed back by Parse (roundtrip)
//   - Empty fields map produces "---\n---\n" + body
//   - Empty body produces frontmatter with no trailing content
//   - No trailing newline is added to body (preserves original)
//
// Verification: auto
// Dependencies: none
```

---

## S-36: State Reader

*User Actions:*
- *1. User can see project status from the terminal without Claude (gl status)*
- *Supports all user actions that need slice state or dependency data*

### C-93: StateReadSlices

```go
// Contract: StateReadSlices
// Boundary: State reader -> Filesystem (read all slice frontmatter into structs)
// Slice: S-36 (State Reader)
// Design refs: FR-10, DESIGN.md 4.4
//
// FILE: internal/state/state.go
//
// Signature:
//   func ReadSlices(dir string) ([]SliceInfo, error)
//
// Types:
//   type SliceInfo struct {
//       ID           string   // e.g. "S-35"
//       Status       string   // "pending", "in_progress", "complete", "failed"
//       Step         string   // current step within the slice
//       Milestone    string   // milestone name
//       Started      string   // ISO date
//       Updated      string   // ISO datetime
//       Tests        int      // test count
//       SecurityTests int     // security test count
//       Session      string   // advisory session ID
//       Deps         []string // dependency slice IDs
//   }
//
// Input:
//   dir string  // path to .greenlight/slices/ directory
//
// Output:
//   slices []SliceInfo  // parsed slice info, sorted by ID
//   err    error        // nil on success
//
// Behaviour:
//   1. Read all .md files from dir
//   2. Parse frontmatter from each file using frontmatter.Parse
//   3. Convert frontmatter fields to SliceInfo struct
//   4. Sort by ID (ascending)
//   5. Return slice list
//
// Errors:
//   | Error | When |
//   |-------|------|
//   | ErrDirNotFound | Provided directory does not exist |
//   | ErrNoSliceFiles | Directory exists but contains no .md files |
//   | ErrParseFailure | A slice file has invalid frontmatter (includes filename in error) |
//
// Invariants:
//   - Returned slices are always sorted by ID ascending
//   - Tests and SecurityTests default to 0 if missing or non-numeric
//   - Deps field is split on "," with whitespace trimmed; empty string yields empty slice
//   - Invalid files are reported individually (not batch failure)
//   - Status values outside the known set are preserved as-is (no validation)
//
// Verification: auto
// Dependencies: C-91 (FrontmatterParse)
```

### C-94: StateReadGraph

```go
// Contract: StateReadGraph
// Boundary: State reader -> Filesystem (read GRAPH.json dependency data)
// Slice: S-36 (State Reader)
// Design refs: FR-10, DESIGN.md 4.4
//
// FILE: internal/state/state.go
//
// Signature:
//   func ReadGraph(path string) (*Graph, error)
//
// Types:
//   type Graph struct {
//       Slices map[string]GraphSlice  // keyed by slice ID
//       Edges  []Edge                 // dependency edges
//   }
//
//   type GraphSlice struct {
//       ID          string   // e.g. "S-35"
//       Name        string   // human-readable name
//       DependsOn   []string // slice IDs this depends on
//       Wave        int      // wave assignment
//       Contracts   []string // contract IDs
//   }
//
//   type Edge struct {
//       From   string // dependent slice
//       To     string // dependency slice
//       Reason string // why the dependency exists
//   }
//
// Input:
//   path string  // path to .greenlight/GRAPH.json
//
// Output:
//   graph *Graph  // parsed graph data
//   err   error   // nil on success
//
// Behaviour:
//   1. Read file at path
//   2. Parse JSON into Graph struct
//   3. Return parsed graph
//
// Errors:
//   | Error | When |
//   |-------|------|
//   | ErrFileNotFound | GRAPH.json does not exist at path |
//   | ErrInvalidJSON | File content is not valid JSON |
//   | ErrMissingSlices | JSON is valid but has no "slices" field |
//
// Invariants:
//   - Graph always has at least a Slices map (possibly empty)
//   - Missing optional fields default to zero values
//   - Unknown JSON fields are ignored (forward compatibility)
//   - File is read once per call (no caching)
//
// Verification: auto
// Dependencies: none
```

### C-95: StateFindReadySlices

```go
// Contract: StateFindReadySlices
// Boundary: State reader -> Ready slice computation (pending + deps complete)
// Slice: S-36 (State Reader)
// Design refs: FR-10, DESIGN.md 4.4
//
// FILE: internal/state/state.go
//
// Signature:
//   func FindReadySlices(slices []SliceInfo, graph *Graph) []SliceInfo
//
// Input:
//   slices []SliceInfo  // current slice states from ReadSlices
//   graph  *Graph       // dependency graph from ReadGraph
//
// Output:
//   ready []SliceInfo  // slices that are ready to build, sorted by wave then ID
//
// Behaviour:
//   1. For each slice where status == "pending":
//      a. Look up dependencies in graph
//      b. Check all dependency slices have status == "complete"
//      c. If all deps complete, include in ready list
//   2. Sort by wave (ascending), then by ID (ascending) within wave
//   3. Return ready list
//
// Errors: none (pure computation, always returns a result)
//
// Invariants:
//   - A slice with no dependencies and status "pending" is always ready
//   - A slice with status "in_progress", "complete", or "failed" is never ready
//   - Ready slices are sorted by wave first, then by ID
//   - If graph has no entry for a slice, it is treated as having no dependencies
//   - Empty input returns empty output
//   - Result is deterministic for the same input
//
// Verification: auto
// Dependencies: C-93 (StateReadSlices), C-94 (StateReadGraph)
```

### C-96: StateDetectContext

```go
// Contract: StateDetectContext
// Boundary: State reader -> Environment (detect shell vs Claude context)
// Slice: S-36 (State Reader)
// Design refs: FR-3, D-44, DESIGN.md 4.10
//
// FILE: internal/state/state.go
//
// Signature:
//   func DetectContext() ExecutionContext
//
// Types:
//   type ExecutionContext struct {
//       InsideClaude bool   // true if $CLAUDE_CODE env var is set
//       ClaudeValue  string // value of $CLAUDE_CODE (empty if not set)
//   }
//
// Input: none (reads from environment)
//
// Output:
//   ctx ExecutionContext  // execution context information
//
// Behaviour:
//   1. Read $CLAUDE_CODE environment variable
//   2. If set (any non-empty value): InsideClaude = true
//   3. If unset or empty: InsideClaude = false
//
// Errors: none (environment read always succeeds)
//
// Invariants:
//   - Detection is based solely on $CLAUDE_CODE env var (D-44)
//   - Any non-empty value means "inside Claude"
//   - Empty string and unset both mean "in shell"
//   - Function is side-effect free (only reads environment)
//
// Verification: auto
// Dependencies: none
```

---

## S-37: CLI Dispatch Extension

*User Actions:*
- *7. User can get help with context-aware command listing (gl help)*
- *Supports all new CLI commands (dispatch routing)*

### C-97: CLIDispatchExtension

```go
// Contract: CLIDispatchExtension
// Boundary: User -> CLI dispatcher (new subcommands added to Run switch)
// Slice: S-37 (CLI Dispatch Extension)
// Design refs: FR-1, D-46, DESIGN.md 4.1
//
// FILE: internal/cli/cli.go
//
// Change: Extend the existing Run() switch statement with new cases.
//         Update printUsage() with new command listing.
//
// Existing switch cases (unchanged):
//   "install"   -> cmd.RunInstall
//   "uninstall" -> cmd.RunUninstall
//   "check"     -> cmd.RunCheck
//   "version"   -> cmd.RunVersion
//   "help"      -> printUsage
//
// New switch cases:
//   "status"    -> cmd.RunStatus(args[1:], stdout)
//   "slice"     -> cmd.RunSlice(args[1:], stdout)
//   "init"      -> cmd.RunInit(args[1:], stdout)
//   "design"    -> cmd.RunDesign(args[1:], stdout)
//   "roadmap"   -> cmd.RunRoadmap(args[1:], stdout)
//   "changelog" -> cmd.RunChangelog(args[1:], stdout)
//
// Updated printUsage includes all commands grouped by category:
//   Project lifecycle: init, design, roadmap
//   Building: slice
//   State & progress: status, changelog
//   Admin: install, uninstall, check, version, help
//
// Signature (unchanged):
//   func Run(args []string, contentFS fs.FS, stdout io.Writer) int
//
// Input: unchanged
// Output: unchanged (exit code int)
//
// Errors:
//   | Error | When |
//   |-------|------|
//   | UnknownCommand | args[0] not in any switch case | prints error + usage, returns 1 |
//
// Invariants:
//   - All existing commands continue to work identically (no regression)
//   - New commands follow the same pattern: return cmd.RunXxx(args[1:], stdout)
//   - printUsage() shows all commands grouped by category
//   - help, --help, -h all print the updated usage
//   - Unknown commands still print error message + usage + return 1
//   - New command handlers receive args[1:] (subcommand args only)
//
// Verification: verify
// Acceptance Criteria:
// - Running `gl` with no args shows updated help with all new commands grouped by category
// - Running `gl unknowncmd` prints error message and updated usage
// - All existing commands (install, check, uninstall, version) still work
// - New commands (status, slice, init, design, roadmap, changelog) dispatch correctly
//
// Dependencies: none (dispatch only -- handlers are stubbed until their slices)
```

---

## S-38: Status Command

*User Actions:*
- *1. User can see project status from the terminal without Claude (gl status)*

### C-98: RunStatus

```go
// Contract: RunStatus
// Boundary: CLI -> Command handler (status with progress display)
// Slice: S-38 (Status Command)
// Design refs: FR-7, FR-10, DESIGN.md 4.3, 7
//
// FILE: internal/cmd/status.go
//
// Signature:
//   func RunStatus(args []string, stdout io.Writer) int
//
// Input:
//   args   []string   // subcommand args (may contain --compact)
//   stdout io.Writer  // output destination
//
// Output:
//   exit code int  // 0 on success, 1 on error
//
// Printed output (default mode):
//   Progress: [########..........] 18/36 slices
//   Running:  S-28 (implementing), S-30 (testing)
//   Ready:    S-35, S-36, S-37
//   Blocked:  S-42 (needs S-36, S-41)
//   Tests:    1,247 passing, 134 security
//
// Behaviour:
//   1. Locate .greenlight/ directory (current dir or parent)
//   2. Read all slice files via state.ReadSlices
//   3. Read GRAPH.json via state.ReadGraph
//   4. Compute ready slices via state.FindReadySlices
//   5. Compute summary statistics:
//      - Total slices, complete count, in_progress count
//      - Ready slices (pending with all deps complete)
//      - Blocked slices (pending with incomplete deps)
//      - Total tests, total security tests
//   6. Format and print to stdout
//
// Errors:
//   | Error | When |
//   |-------|------|
//   | NoGreenlightDir | No .greenlight/ directory found | prints "Not a greenlight project. Run 'gl init' first.", returns 1 |
//   | NoSliceFiles | .greenlight/slices/ has no .md files | prints "No slices found. Run 'gl init' to set up.", returns 1 |
//   | GraphReadError | Cannot read GRAPH.json | prints status without dependency info, warns user |
//
// Invariants:
//   - Status is computed fresh on every call (no caching)
//   - Blocked slices show which dependencies are unmet
//   - Output is human-readable (not JSON)
//   - Progress bar uses ASCII characters only
//   - Test counts are summed from all slice files
//   - Zero-state (no slices) is handled gracefully
//   - Runs without Claude process (local command)
//
// Security:
//   - Read-only operation
//   - No sensitive data displayed
//
// Verification: verify
// Acceptance Criteria:
// - Running `gl status` displays progress bar, running/ready/blocked slices, and test counts
// - Progress bar accurately reflects complete/total ratio
// - Blocked slices show their unmet dependency IDs
// - Running `gl status` outside a greenlight project prints an error message
//
// Steps:
// - Run `gl status` in a project with mixed slice states
// - Verify counts match the frontmatter in .greenlight/slices/*.md files
//
// Dependencies: C-93 (StateReadSlices), C-94 (StateReadGraph), C-95 (StateFindReadySlices)
```

### C-99: RunStatusCompact

```go
// Contract: RunStatusCompact
// Boundary: CLI -> Command handler (compact one-liner for tmux status bar)
// Slice: S-38 (Status Command)
// Design refs: DESIGN.md 4.9
//
// FILE: internal/cmd/status.go
//
// Signature: handled by RunStatus when --compact flag is present
//
// Input:
//   args contains "--compact"
//
// Output:
//   Single line to stdout: "18/36 done | 4 running"
//   exit code 0
//
// Behaviour:
//   1. Same data read as RunStatus (slices + graph)
//   2. Output single line: "{complete}/{total} done | {in_progress} running"
//   3. No progress bar, no details, no newlines beyond the final one
//
// Errors:
//   | Error | When |
//   |-------|------|
//   | DataReadError | Cannot read slice files | prints "? slices | ? running", returns 0 (status bar must not break) |
//
// Invariants:
//   - Output is always a single line
//   - On error, outputs placeholder rather than failing (tmux status bar resilience)
//   - Format is exactly: "{N}/{M} done | {K} running"
//   - No color codes or special characters (tmux compatible)
//
// Verification: auto
// Dependencies: C-93 (StateReadSlices), C-98 (RunStatus shares data reading logic)
```

---

## S-39: Help Command

*User Actions:*
- *7. User can get help with context-aware command listing (gl help)*

### C-100: RunHelp

```go
// Contract: RunHelp
// Boundary: CLI -> Command handler (context-aware help listing)
// Slice: S-39 (Help Command)
// Design refs: FR-7, DESIGN.md 4.1
//
// FILE: internal/cmd/help.go
//
// Signature:
//   func RunHelp(args []string, stdout io.Writer) int
//
// Input:
//   args   []string   // subcommand args (unused for now)
//   stdout io.Writer  // output destination
//
// Output:
//   exit code int  // always 0
//
// Printed output:
//   All available commands with descriptions, grouped by category.
//   If inside a greenlight project, appends current state summary.
//
// Behaviour:
//   1. Print command listing grouped by category:
//      - Project lifecycle: init, design, roadmap
//      - Building: slice
//      - State & progress: status, changelog
//      - Admin: install, uninstall, check, version, help
//   2. Detect if .greenlight/ directory exists
//   3. If exists: read slice count, complete count, ready count
//   4. Append state summary: "Current project: {N} slices, {K} complete, {R} ready"
//   5. If not exists: append "Run 'gl init' to start a new project"
//
// Errors:
//   | Error | When |
//   |-------|------|
//   | StateReadError | Cannot read slice state | Print help without state summary. Do not fail. |
//
// Invariants:
//   - Help always returns 0 (never fails)
//   - Command listing is always printed regardless of state read errors
//   - State summary is best-effort (errors are silently ignored)
//   - Context detection uses directory existence, not $CLAUDE_CODE
//   - Runs without Claude process (local command)
//
// Verification: verify
// Acceptance Criteria:
// - Running `gl help` displays all commands grouped by category
// - Inside a greenlight project, help appends current project state summary
// - Outside a greenlight project, help suggests running `gl init`
//
// Dependencies: C-93 (StateReadSlices, for state summary)
```

---

## S-40: Roadmap and Changelog Commands

*User Actions:*
- *6. User can see roadmap and changelog from the terminal (gl roadmap, gl changelog)*

### C-101: RunRoadmap

```go
// Contract: RunRoadmap
// Boundary: CLI -> Command handler (display ROADMAP.md)
// Slice: S-40 (Roadmap and Changelog Commands)
// Design refs: FR-7, DESIGN.md 4.3
//
// FILE: internal/cmd/roadmap.go
//
// Signature:
//   func RunRoadmap(args []string, stdout io.Writer) int
//
// Input:
//   args   []string   // subcommand args (unused)
//   stdout io.Writer  // output destination
//
// Output:
//   exit code int  // 0 on success, 1 on error
//
// Behaviour:
//   1. Locate .greenlight/ directory
//   2. Read .greenlight/ROADMAP.md
//   3. Print contents to stdout
//
// Errors:
//   | Error | When |
//   |-------|------|
//   | NoGreenlightDir | No .greenlight/ directory | prints "Not a greenlight project.", returns 1 |
//   | NoRoadmap | ROADMAP.md does not exist | prints "No roadmap found. Run 'gl design' first.", returns 1 |
//   | ReadError | Cannot read file | prints error message, returns 1 |
//
// Invariants:
//   - Read-only operation (no files modified)
//   - Contents printed verbatim (no formatting transformation)
//   - Runs without Claude process (local command)
//   - Returns 0 only on successful read and print
//
// Verification: verify
// Acceptance Criteria:
// - Running `gl roadmap` displays the contents of .greenlight/ROADMAP.md
// - Running `gl roadmap` with no ROADMAP.md prints a helpful error
//
// Dependencies: none
```

### C-102: RunChangelog

```go
// Contract: RunChangelog
// Boundary: CLI -> Command handler (display changelog from summaries)
// Slice: S-40 (Roadmap and Changelog Commands)
// Design refs: FR-7, DESIGN.md 4.3
//
// FILE: internal/cmd/changelog.go
//
// Signature:
//   func RunChangelog(args []string, stdout io.Writer) int
//
// Input:
//   args   []string   // subcommand args (unused)
//   stdout io.Writer  // output destination
//
// Output:
//   exit code int  // 0 on success, 1 on error
//
// Behaviour:
//   1. Locate .greenlight/ directory
//   2. Read all .md files from .greenlight/summaries/
//   3. Sort by filename (chronological -- filenames use date prefix)
//   4. Print each summary separated by "---"
//
// Errors:
//   | Error | When |
//   |-------|------|
//   | NoGreenlightDir | No .greenlight/ directory | prints "Not a greenlight project.", returns 1 |
//   | NoSummaries | summaries/ directory missing or empty | prints "No changelog entries yet.", returns 0 |
//   | ReadError | Cannot read a summary file | skip file, warn, continue with remaining |
//
// Invariants:
//   - Read-only operation (no files modified)
//   - Summaries are sorted by filename ascending (oldest first)
//   - Missing summaries directory is not an error (empty changelog)
//   - Individual file read errors do not block other files
//   - Runs without Claude process (local command)
//
// Verification: verify
// Acceptance Criteria:
// - Running `gl changelog` displays all summary files in chronological order
// - Running `gl changelog` with no summaries prints "No changelog entries yet."
//
// Dependencies: none
```

---

## S-41: Process Spawner

*User Actions:*
- *Supports user actions 2-5 (any command that launches Claude)*

### C-103: ProcessSpawnClaude

```go
// Contract: ProcessSpawnClaude
// Boundary: Process spawner -> os/exec (spawn Claude with configurable flags)
// Slice: S-41 (Process Spawner)
// Design refs: FR-2, D-42, DESIGN.md 4.3, 4.4
//
// FILE: internal/process/process.go
//
// Signature:
//   func SpawnClaude(opts SpawnOptions) (*exec.Cmd, error)
//
// Types:
//   type SpawnOptions struct {
//       Prompt       string   // prompt text for -p flag (headless mode)
//       Flags        []string // additional CLI flags (e.g. --dangerously-skip-permissions)
//       Dir          string   // working directory for the process
//       Stdout       io.Writer // stdout destination
//       Stderr       io.Writer // stderr destination
//   }
//
// Input:
//   opts SpawnOptions  // configuration for the Claude process
//
// Output:
//   cmd *exec.Cmd  // the started process (caller manages lifecycle)
//   err error      // nil if process started successfully
//
// Behaviour:
//   1. Verify "claude" is in PATH
//   2. Build command: claude -p "{prompt}" {flags...}
//   3. Set working directory to opts.Dir
//   4. Connect stdout and stderr
//   5. Start the process (cmd.Start)
//   6. Return the running command (caller calls cmd.Wait)
//
// Errors:
//   | Error | When |
//   |-------|------|
//   | ErrClaudeNotFound | "claude" binary not in PATH | "claude not found. Install Claude Code: https://claude.ai/code" |
//   | ErrStartFailure | Process failed to start | wraps underlying os error |
//   | ErrEmptyPrompt | Prompt is empty string | "prompt is required for headless mode" |
//
// Invariants:
//   - Process is started but not waited on (caller manages lifecycle)
//   - Prompt is required (non-empty) for headless mode
//   - "claude" binary is verified before attempting to start
//   - Flags are passed as-is (no validation by process package)
//   - Dir defaults to current directory if empty
//   - No environment variables are modified (inherits current env)
//
// Verification: auto
// Dependencies: none
```

### C-104: ProcessSpawnInteractive

```go
// Contract: ProcessSpawnInteractive
// Boundary: Process spawner -> os/exec (launch interactive Claude session)
// Slice: S-41 (Process Spawner)
// Design refs: FR-2, FR-8, D-42, DESIGN.md 4.3
//
// FILE: internal/process/process.go
//
// Signature:
//   func SpawnInteractive(opts InteractiveOptions) error
//
// Types:
//   type InteractiveOptions struct {
//       Prompt   string   // initial prompt (optional, for skill loading)
//       Flags    []string // additional CLI flags
//       Dir      string   // working directory
//   }
//
// Input:
//   opts InteractiveOptions  // configuration for interactive session
//
// Output:
//   err error  // nil if process completed successfully
//
// Behaviour:
//   1. Verify "claude" is in PATH
//   2. Build command: claude {flags...}
//   3. If Prompt is non-empty: add -p "{prompt}" but NOT --dangerously-skip-permissions
//   4. Connect stdin, stdout, stderr to os.Stdin/os.Stdout/os.Stderr
//   5. Run the process (cmd.Run -- blocks until user exits Claude)
//   6. Return the exit error (nil on clean exit)
//
// Errors:
//   | Error | When |
//   |-------|------|
//   | ErrClaudeNotFound | "claude" binary not in PATH | same message as SpawnClaude |
//   | ErrProcessFailed | Process exited with non-zero code | wraps exec.ExitError |
//
// Invariants:
//   - Interactive mode NEVER uses --dangerously-skip-permissions
//   - stdin is connected to terminal (user can interact)
//   - Prompt is optional (unlike headless mode)
//   - Function blocks until Claude session ends
//   - No environment variables are modified
//
// Verification: auto
// Dependencies: none
```

---

## S-42: Single Slice Command

*User Actions:*
- *2. User can run a single slice headlessly from the terminal (gl slice S-35)*

### C-105: RunSliceSingle

```go
// Contract: RunSliceSingle
// Boundary: CLI -> Command handler (run single slice headlessly)
// Slice: S-42 (Single Slice Command)
// Design refs: FR-2, FR-3, DESIGN.md 4.4, 4.5
//
// FILE: internal/cmd/slice.go
//
// Signature:
//   func RunSlice(args []string, stdout io.Writer) int
//
// Input:
//   args   []string   // e.g. ["S-35"] or ["S-35", "--dry-run"]
//   stdout io.Writer  // output destination
//
// Output:
//   exit code int  // 0 on success, non-zero on failure
//
// Behaviour (when slice ID is provided):
//   1. Parse args for slice ID and flags (--dry-run, --max, --watch, --sequential)
//   2. Detect context via state.DetectContext()
//   3. Read config from .greenlight/config.json for claude_flags
//   4. If --dry-run: print what would happen, return 0
//   5. If inside Claude ($CLAUDE_CODE set):
//      - Print slice info for Claude skill to consume
//      - Return 0
//   6. If in shell:
//      - Spawn claude headlessly: claude -p "/gl:slice {id}" {claude_flags}
//      - Wait for completion
//      - Return claude's exit code
//
// Errors:
//   | Error | When |
//   |-------|------|
//   | InvalidSliceID | Provided ID not found in GRAPH.json | "Unknown slice: {id}. Run 'gl status' to see available slices." |
//   | ClaudeSpawnError | Failed to start Claude process | prints error from process.SpawnClaude |
//   | NoGreenlightDir | No .greenlight/ directory | "Not a greenlight project." |
//
// Invariants:
//   - Single slice mode always runs directly (no tmux)
//   - Inside Claude context: outputs info, never spawns another Claude
//   - claude_flags come from config.json parallel.claude_flags
//   - Exit code reflects Claude process exit code
//   - --dry-run never spawns a process
//
// Security:
//   - claude_flags are read from project config, not user input
//
// Verification: verify
// Acceptance Criteria:
// - Running `gl slice S-35` spawns a headless Claude session for that slice
// - Running `gl slice S-35 --dry-run` prints what would happen without spawning
// - Running `gl slice INVALID` prints an error with the unknown ID
//
// Steps:
// - Run `gl slice S-35 --dry-run` and verify output shows the command that would execute
//
// Dependencies: C-96 (StateDetectContext), C-94 (StateReadGraph), C-103 (ProcessSpawnClaude)
```

### C-106: RunSliceAutoDetect

```go
// Contract: RunSliceAutoDetect
// Boundary: CLI -> Command handler (auto-detect ready slices, run one)
// Slice: S-42 (Single Slice Command)
// Design refs: FR-2, FR-10, DESIGN.md 4.4
//
// FILE: internal/cmd/slice.go
//
// Signature: handled by RunSlice when no slice ID is provided
//
// Input:
//   args contains no slice ID (e.g. [] or ["--dry-run"])
//
// Output:
//   exit code int  // 0 on success, 1 on error
//
// Behaviour (when no slice ID):
//   1. Read slice states and graph
//   2. Find ready slices via state.FindReadySlices
//   3. If 0 ready:
//      - Print blocked status: which slices are blocked and why
//      - Return 0
//   4. If 1 ready:
//      - Run that slice directly (same as providing ID)
//   5. If 2+ ready and single-slice context (inside Claude or --sequential):
//      - Pick first ready slice by wave/ID order
//      - Run it
//      - Print hint: "{N} more slices ready. Run 'gl slice --max {N}'"
//   6. If 2+ ready and parallel possible:
//      - Defer to parallel execution (C-111)
//
// Errors:
//   | Error | When |
//   |-------|------|
//   | StateReadError | Cannot read slice state or graph | prints error, returns 1 |
//
// Invariants:
//   - Auto-detect always picks by wave order, then ID order
//   - Inside Claude: always runs exactly one slice, hints about the rest
//   - Zero ready slices is not an error (returns 0 with information)
//   - Blocked status shows which deps are unmet for each blocked slice
//
// Verification: verify
// Acceptance Criteria:
// - Running `gl slice` with no args auto-detects and runs the first ready slice
// - When no slices are ready, prints blocked status with unmet dependencies
// - When multiple slices are ready, prints count and suggests parallel mode
//
// Dependencies: C-93 (StateReadSlices), C-94 (StateReadGraph), C-95 (StateFindReadySlices), C-105 (RunSliceSingle)
```

---

## S-43: tmux Manager

*User Actions:*
- *3. User can run multiple slices in parallel via tmux (gl slice --max 4)*

### C-107: TmuxIsAvailable

```go
// Contract: TmuxIsAvailable
// Boundary: tmux manager -> os/exec (check tmux availability)
// Slice: S-43 (tmux Manager)
// Design refs: D-41, DESIGN.md 4.7
//
// FILE: internal/tmux/tmux.go
//
// Signature:
//   func IsAvailable() bool
//
// Input: none
//
// Output:
//   available bool  // true if tmux is installed and executable
//
// Behaviour:
//   1. Run exec.LookPath("tmux")
//   2. Return true if found, false otherwise
//
// Errors: none (returns false on any error)
//
// Invariants:
//   - Never returns an error (boolean only)
//   - Checks PATH lookup, does not verify tmux version
//   - Does not start a tmux process
//   - Result may change between calls (if tmux installed/uninstalled)
//
// Verification: auto
// Dependencies: none
```

### C-108: TmuxNewSession

```go
// Contract: TmuxNewSession
// Boundary: tmux manager -> os/exec (create named tmux session)
// Slice: S-43 (tmux Manager)
// Design refs: D-41, DESIGN.md 4.8
//
// FILE: internal/tmux/tmux.go
//
// Signature:
//   func NewSession(opts SessionOptions) error
//
// Types:
//   type SessionOptions struct {
//       Name    string // session name (e.g. "gl-greenlight")
//       Dir     string // working directory for the session
//       Command string // initial command to run in first window
//       Window  string // name for the first window
//   }
//
// Input:
//   opts SessionOptions  // session configuration
//
// Output:
//   err error  // nil if session created successfully
//
// Behaviour:
//   1. Run: tmux new-session -d -s {name} -n {window} -c {dir} {command}
//   2. -d flag creates session detached (does not attach immediately)
//
// Errors:
//   | Error | When |
//   |-------|------|
//   | ErrSessionExists | Session with this name already exists | "tmux session '{name}' already exists" |
//   | ErrTmuxNotFound | tmux not in PATH | "tmux not found" |
//   | ErrCreateFailed | tmux command failed | wraps underlying error |
//
// Invariants:
//   - Session is created detached (caller attaches separately)
//   - Session name must be non-empty
//   - Command runs in the first window of the session
//   - Working directory is set for the session (inherited by windows)
//
// Verification: auto
// Dependencies: C-107 (TmuxIsAvailable)
```

### C-109: TmuxAddWindow

```go
// Contract: TmuxAddWindow
// Boundary: tmux manager -> os/exec (add window to existing session)
// Slice: S-43 (tmux Manager)
// Design refs: D-41, DESIGN.md 4.8
//
// FILE: internal/tmux/tmux.go
//
// Signature:
//   func AddWindow(session string, name string, command string) error
//
// Input:
//   session string  // existing session name
//   name    string  // window name (e.g. "S-35")
//   command string  // command to run in the new window
//
// Output:
//   err error  // nil if window added successfully
//
// Behaviour:
//   1. Run: tmux new-window -t {session} -n {name} {command}
//
// Errors:
//   | Error | When |
//   |-------|------|
//   | ErrSessionNotFound | Target session does not exist | "tmux session '{session}' not found" |
//   | ErrAddWindowFailed | tmux command failed | wraps underlying error |
//
// Invariants:
//   - Window is added to an existing session (session must exist)
//   - Window name is set for identification
//   - Command starts immediately in the new window
//
// Verification: auto
// Dependencies: C-108 (TmuxNewSession must exist to add windows to)
```

### C-110: TmuxAttachSession

```go
// Contract: TmuxAttachSession
// Boundary: tmux manager -> os/exec (attach to existing session)
// Slice: S-43 (tmux Manager)
// Design refs: D-41, DESIGN.md 4.8
//
// FILE: internal/tmux/tmux.go
//
// Signature:
//   func AttachSession(session string) error
//
// Input:
//   session string  // session name to attach to
//
// Output:
//   err error  // nil if attached and then detached cleanly
//
// Behaviour:
//   1. Run: tmux attach-session -t {session}
//   2. Blocks until user detaches or session ends
//   3. Connects stdin/stdout/stderr to terminal
//
// Errors:
//   | Error | When |
//   |-------|------|
//   | ErrSessionNotFound | Target session does not exist | "tmux session '{session}' not found" |
//   | ErrAttachFailed | tmux command failed | wraps underlying error |
//
// Invariants:
//   - Function blocks until user detaches (Ctrl-B D) or session terminates
//   - stdin is connected to terminal (user can interact with tmux)
//   - Attaching to an already-attached session is valid (tmux handles it)
//
// Verification: auto
// Dependencies: C-108 (TmuxNewSession must create the session first)
```

---

## S-44: Parallel Slice Execution

*User Actions:*
- *3. User can run multiple slices in parallel via tmux (gl slice --max 4)*

### C-111: RunSliceParallel

```go
// Contract: RunSliceParallel
// Boundary: CLI -> Command handler (parallel slice execution via tmux)
// Slice: S-44 (Parallel Slice Execution)
// Design refs: FR-4, D-41, DESIGN.md 4.4, 4.8
//
// FILE: internal/cmd/slice.go
//
// Signature: handled by RunSlice when 2+ slices ready and tmux available
//
// Input:
//   ready  []SliceInfo  // 2+ ready slices
//   max    int          // from --max flag (default 4)
//
// Output:
//   exit code int  // 0 on success, 1 on error
//
// Behaviour:
//   1. Read config for parallel settings (claude_flags, tmux_session_prefix)
//   2. Create tmux session named "{prefix}-{project}" (e.g. "gl-greenlight")
//      - First window runs first ready slice
//   3. Add windows for remaining ready slices (up to --max)
//      - Each window runs: claude -p "/gl:slice {id}" {claude_flags}
//   4. Set tmux layout to tiled
//   5. Attach to session (blocks until user detaches)
//   6. Return 0
//
// Errors:
//   | Error | When |
//   |-------|------|
//   | TmuxSessionCreateError | Cannot create tmux session | prints error, falls back to sequential |
//   | TmuxWindowError | Cannot add window | skip that slice, continue with others |
//   | ConfigReadError | Cannot read config.json | use defaults (max=4, standard flags) |
//
// Invariants:
//   - Never exceeds --max windows
//   - Session name follows "{prefix}-{project}" pattern from config
//   - Falls back to sequential mode on tmux errors (not a fatal error)
//   - Each window is named with the slice ID (e.g. "S-35")
//   - Slices are assigned to windows in wave/ID order
//   - tmux tiled layout is applied after all windows are added
//
// Security:
//   - claude_flags from config only (not from user args)
//
// Verification: verify
// Acceptance Criteria:
// - Running `gl slice --max 4` with 5 ready slices creates a tmux session with 4 windows
// - Each tmux window runs a headless Claude session for one slice
// - Session name follows the configured prefix pattern
// - If tmux is unavailable, falls back to sequential execution with a hint
//
// Steps:
// - Run `gl slice --max 2 --dry-run` and verify it shows 2 slices would be launched in tmux
//
// Dependencies: C-95 (StateFindReadySlices), C-103 (ProcessSpawnClaude), C-108 (TmuxNewSession), C-109 (TmuxAddWindow), C-110 (TmuxAttachSession)
```

### C-112: RunSliceSequentialFallback

```go
// Contract: RunSliceSequentialFallback
// Boundary: CLI -> Command handler (sequential fallback when no tmux)
// Slice: S-44 (Parallel Slice Execution)
// Design refs: FR-6, DESIGN.md 4.7
//
// FILE: internal/cmd/slice.go
//
// Signature: handled by RunSlice when 2+ slices ready but no tmux or --sequential
//
// Input:
//   ready []SliceInfo  // 2+ ready slices
//
// Output:
//   exit code int  // 0 if all succeeded, 1 if any failed
//
// Behaviour:
//   1. Print: "tmux not available, running sequentially"
//      (or "sequential mode" if --sequential flag)
//   2. Pick first ready slice by wave/ID order
//   3. Run it via SpawnClaude (blocks until complete)
//   4. Re-read slice states (state may have changed)
//   5. Find new ready slices
//   6. If more ready: run next one
//   7. Repeat until no more ready slices
//   8. Print summary at end
//
// Errors:
//   | Error | When |
//   |-------|------|
//   | SliceFailure | A slice exits non-zero | log failure, continue with next ready slice |
//   | StateReadError | Cannot re-read state after completion | stop sequential loop, report error |
//
// Invariants:
//   - State is re-read after every slice completion (picks up new ready slices)
//   - Failed slices do not block subsequent slices (unless they are dependencies)
//   - Sequential mode processes one slice at a time
//   - Summary shows how many completed, how many failed, how many remaining
//   - --sequential flag forces this path even when tmux is available
//
// Verification: verify
// Acceptance Criteria:
// - Running `gl slice --sequential` with 3 ready slices runs them one at a time
// - After each slice completes, state is re-read and next ready slice is picked
// - A failing slice does not prevent other independent slices from running
//
// Dependencies: C-95 (StateFindReadySlices), C-103 (ProcessSpawnClaude), C-105 (RunSliceSingle)
```

---

## S-45: Watch Mode

*User Actions:*
- *4. User can fire-and-forget to drain the dependency graph (gl slice --watch)*

### C-113: RunSliceWatch

```go
// Contract: RunSliceWatch
// Boundary: CLI -> Command handler (watch mode poll loop)
// Slice: S-45 (Watch Mode)
// Design refs: FR-5, D-45, DESIGN.md 4.6
//
// FILE: internal/cmd/slice.go
//
// Signature: handled by RunSlice when --watch flag is present
//
// Input:
//   args contains "--watch"
//   max int          // from --max flag (default 4)
//   interval int     // from config parallel.watch_interval_seconds (default 30)
//
// Output:
//   exit code int  // 0 when all done or all blocked, 1 on fatal error
//
// Behaviour:
//   1. Initial launch: find ready slices, launch up to --max via tmux (or sequential)
//   2. Enter poll loop:
//      a. Sleep for watch_interval_seconds
//      b. Re-read all slice states
//      c. Detect completed slices since last poll:
//         - Log: "S-{id} complete ({N} tests)"
//      d. Count running (in_progress) and ready (pending + deps complete)
//      e. Calculate available slots: max - running
//      f. For each available slot (up to ready count):
//         - Launch new tmux window or sequential process
//         - Log: "launched S-{id} ({name})"
//      g. If no running AND no ready:
//         - Print summary: total done, total tests, remaining blocked
//         - Exit loop, return 0
//   3. Handle SIGINT/SIGTERM: print summary and exit cleanly
//
// Errors:
//   | Error | When |
//   |-------|------|
//   | StateReadError | Cannot read state during poll | log warning, retry next interval |
//   | LaunchError | Cannot spawn new slice process | log error, skip slot, try next interval |
//
// Invariants:
//   - Poll interval is configurable via config (D-45, default 30s)
//   - Watch mode auto-terminates when no in_progress and no ready slices remain
//   - Completed slices are logged with test counts
//   - Slot refilling respects --max cap
//   - SIGINT triggers clean shutdown (not abrupt kill)
//   - State read errors in poll loop are recoverable (retry next interval)
//   - Watch mode works with both tmux and sequential modes
//
// Security:
//   - No external network access
//   - Poll reads only local files
//
// Verification: verify
// Acceptance Criteria:
// - Running `gl slice --watch` drains the dependency graph by auto-launching new slices
// - Completed slices are logged with test counts as they finish
// - Watch mode terminates when all slices are done or blocked
// - New slices are launched as slots become available
//
// Steps:
// - Run `gl slice --watch --dry-run` to verify poll behaviour description
// - Run `gl slice --watch --max 2` on a project with 4 independent slices
// - Observe that at most 2 run simultaneously, and new ones launch as slots free
//
// Dependencies: C-95 (StateFindReadySlices), C-111 (RunSliceParallel), C-112 (RunSliceSequentialFallback)
```

### C-114: RunSliceDryRun

```go
// Contract: RunSliceDryRun
// Boundary: CLI -> Command handler (preview mode without execution)
// Slice: S-45 (Watch Mode)
// Design refs: DESIGN.md 7
//
// FILE: internal/cmd/slice.go
//
// Signature: handled by RunSlice when --dry-run flag is present
//
// Input:
//   args contains "--dry-run"
//
// Output:
//   exit code 0 (always succeeds)
//
// Printed output:
//   Ready ({N}):   S-35, S-36, S-43
//   Running ({K}): S-28 (implementing)
//   Blocked ({M}): S-42 (needs S-36, S-41), S-44 (needs S-42, S-43)
//   Would launch:  S-35, S-36, S-43 (up to --max)
//
// Behaviour:
//   1. Read slice states and graph
//   2. Categorize slices: ready, running, blocked
//   3. Print summary of what would happen
//   4. Do NOT spawn any processes
//
// Errors:
//   | Error | When |
//   |-------|------|
//   | StateReadError | Cannot read state | prints error, returns 1 |
//
// Invariants:
//   - Never spawns a process or modifies any state
//   - Shows blocked slices with their unmet dependency IDs
//   - "Would launch" respects --max cap
//   - Always returns 0 on successful preview
//   - Works with all combinations of --watch, --max, --sequential
//
// Verification: verify
// Acceptance Criteria:
// - Running `gl slice --dry-run` shows ready, running, blocked categories
// - "Would launch" list respects the --max cap
// - No processes are spawned in dry-run mode
//
// Dependencies: C-93 (StateReadSlices), C-94 (StateReadGraph), C-95 (StateFindReadySlices)
```

---

## S-46: Interactive Commands

*User Actions:*
- *5. User can launch interactive Claude sessions (gl init, gl design)*

### C-115: RunInit

```go
// Contract: RunInit
// Boundary: CLI -> Command handler (launch interactive init session)
// Slice: S-46 (Interactive Commands)
// Design refs: FR-2, FR-8, DESIGN.md 4.3
//
// FILE: internal/cmd/init.go
//
// Signature:
//   func RunInit(args []string, stdout io.Writer) int
//
// Input:
//   args   []string   // subcommand args (unused)
//   stdout io.Writer  // output destination
//
// Output:
//   exit code int  // 0 on success, non-zero on error
//
// Behaviour:
//   1. Detect context via state.DetectContext()
//   2. If inside Claude:
//      - Print instructions for the /gl:init skill
//      - Return 0 (Claude will handle it)
//   3. If in shell:
//      - Print "Launching Greenlight init..."
//      - Launch interactive Claude session with /gl:init prompt
//      - Block until session ends
//      - Return Claude's exit code
//
// Errors:
//   | Error | When |
//   |-------|------|
//   | ClaudeNotFound | claude binary not in PATH | prints install instructions, returns 1 |
//   | LaunchError | Failed to start interactive session | prints error, returns 1 |
//
// Invariants:
//   - Interactive mode: NEVER uses --dangerously-skip-permissions
//   - Inside Claude: prints info, never spawns another Claude
//   - Blocks until user exits the interactive session
//   - No state files are read (init creates them)
//
// Verification: verify
// Acceptance Criteria:
// - Running `gl init` from the shell launches an interactive Claude session
// - The session loads with the /gl:init skill prompt
// - No --dangerously-skip-permissions flag is used
//
// Dependencies: C-96 (StateDetectContext), C-104 (ProcessSpawnInteractive)
```

### C-116: RunDesign

```go
// Contract: RunDesign
// Boundary: CLI -> Command handler (launch interactive design session)
// Slice: S-46 (Interactive Commands)
// Design refs: FR-2, FR-8, DESIGN.md 4.3
//
// FILE: internal/cmd/design.go
//
// Signature:
//   func RunDesign(args []string, stdout io.Writer) int
//
// Input:
//   args   []string   // subcommand args (unused)
//   stdout io.Writer  // output destination
//
// Output:
//   exit code int  // 0 on success, non-zero on error
//
// Behaviour:
//   1. Detect context via state.DetectContext()
//   2. If inside Claude:
//      - Print instructions for the /gl:design skill
//      - Return 0
//   3. If in shell:
//      - Verify .greenlight/ directory exists (design requires existing project)
//      - Print "Launching Greenlight design session..."
//      - Launch interactive Claude session with /gl:design prompt
//      - Block until session ends
//      - Return Claude's exit code
//
// Errors:
//   | Error | When |
//   |-------|------|
//   | NoGreenlightDir | No .greenlight/ directory | "Not a greenlight project. Run 'gl init' first.", returns 1 |
//   | ClaudeNotFound | claude binary not in PATH | prints install instructions, returns 1 |
//   | LaunchError | Failed to start interactive session | prints error, returns 1 |
//
// Invariants:
//   - Interactive mode: NEVER uses --dangerously-skip-permissions
//   - Requires existing .greenlight/ directory (unlike init)
//   - Inside Claude: prints info, never spawns another Claude
//   - Blocks until user exits the interactive session
//
// Verification: verify
// Acceptance Criteria:
// - Running `gl design` from the shell launches an interactive Claude session
// - Running `gl design` outside a greenlight project prints an error
// - The session loads with the /gl:design skill prompt
//
// Dependencies: C-96 (StateDetectContext), C-104 (ProcessSpawnInteractive)
```

---

## Updated User Action Mapping (Parallel State)

| User Action | Slice(s) | Contracts | Enabled By |
|-------------|----------|-----------|------------|
| 1. Run multiple /gl:slice sessions in parallel without state corruption | S-28, S-29, S-30, S-31, S-33, S-34 | C-76, C-77, C-78, C-79, C-80, C-81, C-82, C-84, C-87, C-90 | Foundation docs + init + slice writes + supporting commands + CLAUDE.md rule + manifest |
| 2. See accurate slice status across all concurrent sessions | S-28, S-29, S-31, S-33, S-34 | C-76, C-77, C-80, C-83, C-84, C-87, C-90 | Foundation docs + detection + status aggregation + supporting commands + CLAUDE.md rule + manifest |
| 3. Migrate existing projects to the new state format | S-28, S-32, S-34 | C-76, C-77, C-85, C-86, C-90 | Foundation docs + migration command + backup safety + manifest |

---

## Updated User Action Mapping (CLI Orchestrator)

| User Action | Slice(s) | Contracts | Enabled By |
|-------------|----------|-----------|------------|
| 1. User can see project status from the terminal without Claude (gl status) | S-35, S-36, S-37, S-38 | C-91, C-92, C-93, C-94, C-95, C-97, C-98, C-99 | Frontmatter + state + dispatch + status command |
| 2. User can run a single slice headlessly from the terminal (gl slice S-35) | S-35, S-36, S-37, S-41, S-42 | C-91, C-93, C-94, C-95, C-96, C-97, C-103, C-105, C-106 | Frontmatter + state + dispatch + process + slice command |
| 3. User can run multiple slices in parallel via tmux (gl slice --max 4) | S-35, S-36, S-37, S-41, S-42, S-43, S-44 | C-91, C-93, C-94, C-95, C-96, C-97, C-103, C-105, C-107, C-108, C-109, C-110, C-111, C-112 | All of UA-2 + tmux + parallel execution |
| 4. User can fire-and-forget to drain the dependency graph (gl slice --watch) | S-35, S-36, S-37, S-41, S-42, S-43, S-44, S-45 | C-91, C-93, C-94, C-95, C-97, C-103, C-107, C-108, C-109, C-111, C-113, C-114 | All of UA-3 + watch mode |
| 5. User can launch interactive Claude sessions (gl init, gl design) | S-36, S-37, S-41, S-46 | C-96, C-97, C-104, C-115, C-116 | State detection + dispatch + interactive process + init/design commands |
| 6. User can see roadmap and changelog from the terminal (gl roadmap, gl changelog) | S-37, S-40 | C-97, C-101, C-102 | Dispatch + roadmap/changelog commands |
| 7. User can get help with context-aware command listing (gl help) | S-36, S-37, S-39 | C-93, C-97, C-100 | State + dispatch + help command |
