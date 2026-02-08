# DESIGN.md -- Greenlight CLI Stabilisation

> **Project:** Greenlight
> **Scope:** Stabilise existing CLI with tests and contracts. No new features.
> **Stack:** Go 1.24, stdlib only, zero external dependencies.
> **Date:** 2026-02-08

---

## 1. Requirements

### 1.1 Functional Requirements

#### FR-1: Install Globally (`greenlight install --global`)

| ID | Requirement |
|----|-------------|
| FR-1.1 | Copy all 41 manifest files from embedded FS to `~/.claude/<relPath>` |
| FR-1.2 | Create intermediate directories (`agents/`, `commands/gl/`, `references/`, `templates/`) with `0o755` permissions |
| FR-1.3 | Write files with `0o644` permissions |
| FR-1.4 | Place CLAUDE.md at `~/.claude/CLAUDE.md` (inside target directory) |
| FR-1.5 | Write `.greenlight-version` file containing version, commit, and build date on separate lines |
| FR-1.6 | Print progress line per installed file |
| FR-1.7 | Print summary line on completion: `greenlight installed to <dir>` |
| FR-1.8 | Return exit code 0 on success, 1 on any error |

#### FR-2: Install Locally (`greenlight install --local [--on-conflict=keep|replace|append]`)

| ID | Requirement |
|----|-------------|
| FR-2.1 | Copy all manifest files to `./.claude/<relPath>` |
| FR-2.2 | Place CLAUDE.md asymmetrically at `./CLAUDE.md` (project root, not `.claude/`) |
| FR-2.3 | If CLAUDE.md does not exist at destination, write it directly regardless of strategy |
| FR-2.4 | `--on-conflict=keep` (default): write greenlight version as `CLAUDE_GREENLIGHT.md` alongside existing file |
| FR-2.5 | `--on-conflict=replace`: back up existing to `CLAUDE.md.backup`, overwrite with greenlight content |
| FR-2.6 | `--on-conflict=append`: append greenlight content to existing file with newline separator |
| FR-2.7 | **Return error for invalid `--on-conflict` values** (do not silently default) |
| FR-2.8 | Return exit code 0 on success, 1 on any error |

#### FR-3: Check Installation (`greenlight check --global|--local [--verify]`)

| ID | Requirement |
|----|-------------|
| FR-3.1 | Iterate manifest, stat each expected file path |
| FR-3.2 | Report per-file status: `MISSING`, `EMPTY`, or `ERROR` |
| FR-3.3 | Check `.greenlight-version` exists and print the version from it |
| FR-3.4 | CLAUDE.md path resolution mirrors install (global: inside targetDir, local: project root) |
| FR-3.5 | Print summary: "all N files present" or "X/N files present (Y missing, Z empty)" |
| FR-3.6 | Return exit code 0 if all files present and non-empty, 1 otherwise |
| FR-3.7 | **`--verify` flag**: when set, compare file content against embedded source (hash comparison); report `MODIFIED` for mismatches |
| FR-3.8 | Without `--verify`, check presence and non-emptiness only (current behaviour) |

#### FR-4: Uninstall (`greenlight uninstall --global|--local`)

| ID | Requirement |
|----|-------------|
| FR-4.1 | Remove all manifest files from targetDir (except CLAUDE.md) |
| FR-4.2 | Skip missing files without error (idempotent removal) |
| FR-4.3 | Remove `.greenlight-version` file |
| FR-4.4 | **Remove conflict artifacts**: `CLAUDE_GREENLIGHT.md` and `CLAUDE.md.backup` if present, printing each removal |
| FR-4.5 | Clean up empty directories deepest-first: `commands/gl`, `commands`, `agents`, `references`, `templates` |
| FR-4.6 | Print progress per removed file |
| FR-4.7 | Print summary: `greenlight uninstalled from <dir>` |
| FR-4.8 | Return exit code 0 on success, 1 on filesystem errors (not counting missing files) |

#### FR-5: Show Version (`greenlight version`)

| ID | Requirement |
|----|-------------|
| FR-5.1 | Print `greenlight <version> (commit: <hash>, built: <date>)` to stdout |
| FR-5.2 | Values set via ldflags at build time; defaults: Version="dev", GitCommit="unknown", BuildDate="unknown" |
| FR-5.3 | Return exit code 0 always |

#### FR-6: CLI Dispatcher (cross-cutting)

| ID | Requirement |
|----|-------------|
| FR-6.1 | No args: print usage, exit 0 |
| FR-6.2 | Known command (`install`, `uninstall`, `check`, `version`): dispatch to handler |
| FR-6.3 | `help`, `--help`, `-h`: print usage, exit 0 |
| FR-6.4 | Unknown command: print `unknown command: <cmd>`, print usage, exit 1 |
| FR-6.5 | **Accept `io.Writer` parameter** for testability: `Run(args, contentFS, stdout)` |

#### FR-7: Flag Parsing (cross-cutting)

| ID | Requirement |
|----|-------------|
| FR-7.1 | Require exactly one of `--global` or `--local`; reject both or neither |
| FR-7.2 | Parse `--on-conflict=<value>` with valid values: `keep`, `replace`, `append` |
| FR-7.3 | **Return error for invalid `--on-conflict` values** with the invalid value in the message |

#### FR-8: Entry Point Correction

| ID | Requirement |
|----|-------------|
| FR-8.1 | **Print `fs.Sub` error to stderr before `os.Exit(1)`** (currently exits silently) |

### 1.2 Non-Functional Requirements

| ID | Category | Requirement |
|----|----------|-------------|
| NFR-1 | Correctness | Every command produces the documented exit code for success and failure cases |
| NFR-2 | Error handling | All errors are reported to the user with a clear message; no silent failures |
| NFR-3 | Safety | Uninstall never removes CLAUDE.md (may contain user content) |
| NFR-4 | Safety | Replace strategy creates backup before overwriting |
| NFR-5 | Idempotency | Install can be run multiple times; files are overwritten with embedded content |
| NFR-6 | Idempotency | Uninstall can be run multiple times; missing files are skipped without error |
| NFR-7 | Idempotency | Check can be run any number of times with no side effects (read-only) |
| NFR-8 | Testability | All packages testable via dependency injection (fs.FS, io.Writer, temp dirs) |
| NFR-9 | Zero dependencies | No external Go modules. stdlib only |

### 1.3 Constraints

| Constraint | Detail |
|------------|--------|
| Go version | 1.24 |
| Dependencies | stdlib only, zero external modules |
| CI | GoReleaser + semantic-release + lefthook + commitlint |
| Distribution | npm wrapper package (`npx greenlight-cc install`) |
| Embedding | Content embedded via `go:embed` from `src/` directory |
| Commits | Conventional commits enforced by commitlint |
| Build | Version, commit, date injected via ldflags |

### 1.4 Out of Scope

- New CLI commands (upgrade, doctor, etc.)
- Plugin system
- Homebrew tap or curl install script
- Usage analytics
- Coloured output or TUI
- Configuration file for greenlight itself
- Content validation beyond hash comparison

---

## 2. Technical Decisions

| # | Decision | Choice | Rationale |
|---|----------|--------|-----------|
| TD-1 | Invalid `--on-conflict` handling | **Return error** | Silent defaults hide typos. `--on-conflict=replce` would silently keep instead of replacing. Strict validation prevents data-loss surprises. |
| TD-2 | `cli.Run` stdout injection | **Add `io.Writer` parameter** | `Run(args, contentFS, stdout)` -- minimal signature change, enables full CLI dispatch testing without capturing os.Stdout. All subcommands already accept `io.Writer`. |
| TD-3 | Uninstall conflict artifact cleanup | **Remove + print** | `CLAUDE_GREENLIGHT.md` and `CLAUDE.md.backup` are greenlight-created files. Leaving them behind is a leak. Print each removal for transparency. |
| TD-4 | Check content verification | **`--verify` flag** | Presence-only is the right default (fast, no embedded FS needed). Content comparison is opt-in for diagnosing version mismatches. Requires `contentFS` to be passed to `RunCheck`. |

---

## 3. Architecture

### 3.1 Existing Architecture (as-is)

```
main.go                              entry point, go:embed, fs.Sub, cli.Run()
  |
internal/cli/cli.go                  command dispatch (switch on args[0]), usage printing
  |
internal/cmd/                        one file per subcommand
  install.go                         RunInstall(args, contentFS, stdout) -> installer.Install
  uninstall.go                       RunUninstall(args, stdout) -> installer.Uninstall
  check.go                           RunCheck(args, stdout) -> installer.Check
  version.go                         RunVersion(stdout) -> version.Version
  scope.go                           ParseScope, ResolveDir, ParseConflictStrategy
  |
internal/installer/                  core file operations
  installer.go                       Installer struct, Install, Uninstall, Check, Manifest
  conflict.go                        ConflictStrategy type, handleConflict
  |
internal/version/                    build-time variables
  version.go                         Version, GitCommit, BuildDate (ldflags)
```

**Import graph:** `main -> cli -> cmd -> installer, version` (no circular deps).

**Design patterns:**
- **Command**: one `Run*` function per subcommand
- **Strategy**: `ConflictStrategy` enum with `handleConflict` dispatcher
- **Dependency Injection**: `io.Writer` for output, `fs.FS` for embedded content

### 3.2 Corrections

| File | Issue | Correction |
|------|-------|------------|
| `main.go` | `fs.Sub` error causes `os.Exit(1)` with no message | Print error to `os.Stderr` before exiting |
| `internal/cli/cli.go` | `os.Stdout` hardcoded in `Run()` | Add `io.Writer` parameter: `Run(args []string, contentFS fs.FS, stdout io.Writer) int` |
| `internal/cmd/scope.go` | `ParseConflictStrategy` silently ignores invalid values | Return `error` for unrecognised strategy values. Signature becomes `ParseConflictStrategy(args) (ConflictStrategy, []string, error)` |
| `internal/cmd/check.go` | `RunCheck` does not receive `contentFS` | Add `contentFS fs.FS` parameter for `--verify` support |
| `internal/installer/installer.go` | `Uninstall` does not clean conflict artifacts | Add removal of `CLAUDE_GREENLIGHT.md` and `CLAUDE.md.backup` with output |
| `internal/installer/installer.go` | `Check` does not support content verification | Add `verify bool` and `contentFS fs.FS` parameters; when true, compare SHA-256 hashes |
| `internal/installer/installer.go` | Version file removal in `Uninstall` swallows error | Check and report error from `os.Remove` on version file |

### 3.3 Dependency Injection Seams

| Seam | Interface | Test Double |
|------|-----------|-------------|
| Embedded content | `fs.FS` | `testing/fstest.MapFS` |
| Output | `io.Writer` | `bytes.Buffer` |
| Target filesystem | Real OS calls | `t.TempDir()` (real temp dirs, auto-cleanup) |
| Home directory | `os.UserHomeDir()` | Set `$HOME` env var in test |

The filesystem is not mocked. This is a file installer -- its job is writing files. Tests use real temp directories.

---

## 4. Data Model

### 4.1 Manifest

The manifest is a hardcoded `[]string` of 41 relative file paths in `internal/installer/installer.go`. Each path is relative to the embedded FS root (which is `src/` after `fs.Sub`).

```
agents/gl-architect.md
agents/gl-codebase-mapper.md
agents/gl-debugger.md
agents/gl-designer.md
agents/gl-implementer.md
agents/gl-security.md
agents/gl-test-writer.md
agents/gl-verifier.md
commands/gl/add-slice.md
commands/gl/design.md
commands/gl/help.md
commands/gl/init.md
commands/gl/map.md
commands/gl/pause.md
commands/gl/quick.md
commands/gl/resume.md
commands/gl/settings.md
commands/gl/ship.md
commands/gl/slice.md
commands/gl/status.md
references/checkpoint-protocol.md
references/deviation-rules.md
references/verification-patterns.md
templates/config.md
templates/state.md
CLAUDE.md
```

**Directories created during install:**
- `agents/`
- `commands/gl/`
- `references/`
- `templates/`

### 4.2 Version File (`.greenlight-version`)

**Location:** `<targetDir>/.greenlight-version`

**Format:** Three lines, newline-terminated:
```
<version>
<git-commit>
<build-date>
```

**Example:**
```
1.2.0
a1b2c3d
2026-02-08T12:00:00Z
```

### 4.3 Conflict Artifacts

| Strategy | Artifact Created | Location |
|----------|-----------------|----------|
| `keep` | `CLAUDE_GREENLIGHT.md` | Same directory as CLAUDE.md |
| `replace` | `CLAUDE.md.backup` | Same directory as CLAUDE.md |
| `append` | None (modifies in-place) | -- |

### 4.4 CLAUDE.md Placement

| Scope | CLAUDE.md destination | Other files destination |
|-------|----------------------|------------------------|
| `--global` | `~/.claude/CLAUDE.md` | `~/.claude/<relPath>` |
| `--local` | `./CLAUDE.md` (project root) | `./.claude/<relPath>` |

This asymmetry is intentional: Claude Code reads `CLAUDE.md` from the project root for local projects, and from `~/.claude/CLAUDE.md` for global configuration.

---

## 5. API Surface

### 5.1 Commands

```
greenlight install --global [--on-conflict=keep|replace|append]
greenlight install --local  [--on-conflict=keep|replace|append]
greenlight uninstall --global
greenlight uninstall --local
greenlight check --global [--verify]
greenlight check --local  [--verify]
greenlight version
greenlight help
greenlight --help
greenlight -h
greenlight              (no args -- prints help)
```

### 5.2 Flags

| Flag | Commands | Values | Default | Description |
|------|----------|--------|---------|-------------|
| `--global` | install, uninstall, check | -- | -- | Target `~/.claude/` |
| `--local` | install, uninstall, check | -- | -- | Target `./.claude/` |
| `--on-conflict` | install | `keep`, `replace`, `append` | `keep` | CLAUDE.md conflict resolution strategy |
| `--verify` | check | -- | off | Compare file content against embedded source |

### 5.3 Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (bad args, filesystem error, check failure) |

### 5.4 Output Format

All output is plain text to stdout. Error messages are prefixed with `error:`.

**Install output:**
```
  installed agents/gl-architect.md
  installed agents/gl-codebase-mapper.md
  ...
  installed CLAUDE.md -> ./CLAUDE.md
greenlight installed to .claude
```

**Check output (passing):**
```
  version: 1.2.0
all 26 files present
```

**Check output (failing):**
```
  MISSING  agents/gl-architect.md
  EMPTY    references/deviation-rules.md
  version: 1.2.0
24/26 files present (1 missing, 1 empty)
```

**Check output with `--verify` (content mismatch):**
```
  MODIFIED agents/gl-architect.md
  version: 1.2.0
25/26 files verified (0 missing, 0 empty, 1 modified)
```

**Uninstall output:**
```
  removed agents/gl-architect.md
  ...
  removed CLAUDE_GREENLIGHT.md
greenlight uninstalled from .claude
```

**Error output:**
```
error: must specify --global or --local
error: invalid --on-conflict value: banana (valid: keep, replace, append)
```

---

## 6. Security

| Area | Approach |
|------|----------|
| File permissions | Directories: `0o755`. Files: `0o644`. No executable bits. |
| Path safety | All paths are constructed from hardcoded manifest entries or resolved via `os.UserHomeDir()`. No user-supplied paths beyond `--global`/`--local`. |
| No secrets | No credentials, tokens, or API keys in the codebase. Version info is public. |
| No network | Zero network calls. All content is embedded in the binary. |
| CLAUDE.md safety | Uninstall never removes CLAUDE.md. Replace strategy creates backup. Keep strategy preserves original. |
| Symlink safety | `os.WriteFile` and `os.Stat` follow symlinks. This is acceptable -- the tool writes to directories the user controls. No symlink attack surface beyond what the OS provides. |

---

## 7. Deployment

### 7.1 Build

```
go build -ldflags "-X .../version.Version=1.2.0 -X .../version.GitCommit=abc123 -X .../version.BuildDate=2026-02-08" -o greenlight
```

### 7.2 CI Pipeline (existing)

| Tool | Purpose |
|------|---------|
| GoReleaser | Cross-platform binary builds, GitHub releases |
| semantic-release | Version bumping from conventional commits |
| lefthook | Git hooks (pre-commit, commit-msg) |
| commitlint | Conventional commit message enforcement |

### 7.3 Distribution (existing)

| Channel | Method |
|---------|--------|
| npm | `npx greenlight-cc install` (wrapper package) |
| GitHub Releases | Binary downloads via GoReleaser |

### 7.4 Test Integration

Tests run via `go test ./...` in CI. The `config.json` already specifies:
```json
{
  "test": {
    "command": "go test ./...",
    "filter_flag": "-run",
    "coverage_command": "go test -cover ./..."
  }
}
```

---

## 8. Test Architecture

### 8.1 Strategy

**Real filesystem, in-memory content, captured output.**

Each test:
1. Creates a `testing/fstest.MapFS` with fake embedded content (simulating `go:embed` FS)
2. Creates a `t.TempDir()` for the target directory (real filesystem, auto-cleanup)
3. Captures output via `bytes.Buffer`
4. Calls the function under test
5. Asserts: exit code, file existence, file content, output messages

### 8.2 Test File Placement

Co-located with source (standard Go convention):

```
internal/cli/cli_test.go
internal/cmd/install_test.go
internal/cmd/check_test.go
internal/cmd/uninstall_test.go
internal/cmd/version_test.go
internal/cmd/scope_test.go
internal/installer/installer_test.go
internal/installer/conflict_test.go
```

### 8.3 Mocking Approach

| Dependency | Mock? | Technique |
|------------|-------|-----------|
| Embedded content (`fs.FS`) | Yes | `testing/fstest.MapFS` |
| Output (`io.Writer`) | Yes | `bytes.Buffer` |
| Target filesystem | No | `t.TempDir()` (real temp directory) |
| Home directory | Override | `t.Setenv("HOME", tempDir)` |
| Build-time version vars | Override | Set `version.Version` etc. in test setup, restore in cleanup |

### 8.4 Test Categories

#### Package: `internal/installer` (Integration)

| Test | Behaviour |
|------|-----------|
| `TestInstall_WritesAllManifestFiles` | All 41 files written to target dir with correct content |
| `TestInstall_CreatesDirectories` | Intermediate dirs created with correct permissions |
| `TestInstall_WritesVersionFile` | `.greenlight-version` contains version, commit, date |
| `TestInstall_CLAUDEmd_GlobalPlacement` | CLAUDE.md written inside target dir |
| `TestInstall_CLAUDEmd_LocalPlacement` | CLAUDE.md written to parent of target dir |
| `TestInstall_Idempotent` | Running install twice succeeds, files have correct content |
| `TestUninstall_RemovesManifestFiles` | All manifest files removed (except CLAUDE.md) |
| `TestUninstall_SkipsMissingFiles` | No error when files already absent |
| `TestUninstall_RemovesConflictArtifacts` | Removes CLAUDE_GREENLIGHT.md and CLAUDE.md.backup |
| `TestUninstall_CleansEmptyDirs` | Empty dirs removed deepest-first |
| `TestUninstall_PreservesCLAUDEmd` | CLAUDE.md not removed |
| `TestUninstall_RemovesVersionFile` | .greenlight-version removed |
| `TestCheck_AllPresent` | Returns true, prints "all N files present" |
| `TestCheck_MissingFile` | Returns false, prints MISSING line |
| `TestCheck_EmptyFile` | Returns false, prints EMPTY line |
| `TestCheck_PrintsVersion` | Prints version from .greenlight-version |
| `TestCheck_Verify_ContentMatch` | With --verify, all files match embedded content |
| `TestCheck_Verify_ContentMismatch` | With --verify, modified file reported as MODIFIED |

#### Package: `internal/installer` (Unit -- conflict.go)

| Test | Behaviour |
|------|-----------|
| `TestHandleConflict_NoExistingFile` | Writes file directly regardless of strategy |
| `TestHandleConflict_Keep` | Preserves existing, writes CLAUDE_GREENLIGHT.md |
| `TestHandleConflict_Replace` | Backs up existing, overwrites with source |
| `TestHandleConflict_Append` | Appends source to existing with newline |
| `TestHandleConflict_Append_ExistingEndsWithNewline` | No double newline |
| `TestHandleConflict_UnknownStrategy` | Returns error |

#### Package: `internal/cmd` (Unit -- scope.go)

| Test | Behaviour |
|------|-----------|
| `TestParseScope_Global` | Returns "global" |
| `TestParseScope_Local` | Returns "local" |
| `TestParseScope_Both` | Returns error |
| `TestParseScope_Neither` | Returns error |
| `TestParseScope_RemainingArgs` | Unrecognised args returned in remaining |
| `TestResolveDir_Global` | Returns ~/.claude |
| `TestResolveDir_Local` | Returns .claude |
| `TestResolveDir_Unknown` | Returns error |
| `TestParseConflictStrategy_Valid` | Each valid value parsed correctly |
| `TestParseConflictStrategy_Invalid` | Returns error with invalid value in message |
| `TestParseConflictStrategy_Default` | Returns "keep" when flag absent |

#### Package: `internal/cmd` (Integration -- subcommand handlers)

| Test | Behaviour |
|------|-----------|
| `TestRunInstall_Success` | Exit 0, files written |
| `TestRunInstall_MissingScope` | Exit 1, error message |
| `TestRunInstall_InvalidConflictStrategy` | Exit 1, error message |
| `TestRunUninstall_Success` | Exit 0, files removed |
| `TestRunUninstall_MissingScope` | Exit 1, error message |
| `TestRunCheck_AllPresent` | Exit 0 |
| `TestRunCheck_MissingFiles` | Exit 1 |
| `TestRunCheck_WithVerify` | Exit 0 when content matches, exit 1 when modified |
| `TestRunVersion_Output` | Correct format string |

#### Package: `internal/cli` (Integration)

| Test | Behaviour |
|------|-----------|
| `TestRun_NoArgs` | Prints usage, exit 0 |
| `TestRun_Help` | Prints usage, exit 0 |
| `TestRun_UnknownCommand` | Prints error + usage, exit 1 |
| `TestRun_InstallDispatch` | Dispatches to install handler |
| `TestRun_VersionDispatch` | Dispatches to version handler |

### 8.5 Test Helpers

| Helper | Purpose |
|--------|---------|
| `newTestFS(files map[string]string) fstest.MapFS` | Build an in-memory FS from a map of path -> content |
| `assertFileContent(t, path, expected string)` | Read file, compare content, fail with diff |
| `assertFileExists(t, path string)` | Stat file, fail if missing |
| `assertFileNotExists(t, path string)` | Stat file, fail if present |

These helpers are defined in `_test.go` files within each package (not shared across packages) to avoid test coupling.

---

## 9. Deferred

Items explicitly out of scope for MVP. Captured for future consideration.

| Item | Source | Notes |
|------|--------|-------|
| `greenlight upgrade` | Interview | In-place upgrade preserving user config |
| Homebrew tap | Interview | Alternative distribution channel |
| Curl install script | Interview | Alternative distribution channel |
| `greenlight doctor` | Interview | Diagnose Claude Code config issues |
| Plugin system | Interview | Custom agents and commands |
| Local-only usage analytics | Interview | Opt-in telemetry |
| Atomic install (rollback on partial failure) | Design | Currently partial failure leaves files. Could use temp dir + rename. Low priority. |
| Input validation on unknown CLI flags | Design | Unknown flags are silently ignored. Could warn. |
| Content integrity on install (hash check) | Design | Install could verify written content matches source. Belt-and-suspenders. |
| Coloured terminal output | Design | Could use ANSI codes for status indicators. User preference. |

---

## 10. User Decisions

Decisions locked during the design session. These are final and should not be revisited without explicit discussion.

| # | Gray Area | Decision | Rationale |
|---|-----------|----------|-----------|
| UD-1 | Invalid `--on-conflict` values | **Return error** (strict) | Silent defaults hide typos. `--on-conflict=replce` silently keeping instead of replacing is a data-loss risk. Error message includes the invalid value and lists valid options. |
| UD-2 | `cli.Run` stdout injection | **Add `io.Writer` parameter** | `Run(args, contentFS, stdout)`. Minimal signature change. Enables CLI dispatch testing. All subcommands already accept `io.Writer`. |
| UD-3 | Uninstall conflict artifact cleanup | **Remove + print** (transparent) | `CLAUDE_GREENLIGHT.md` and `CLAUDE.md.backup` are greenlight-created files. Leaving them after uninstall is a leak. Each removal is printed for transparency. |
| UD-4 | Check content verification | **`--verify` flag** (opt-in) | Presence-only is the right default (fast, no embedded FS needed). Content hash comparison via `--verify` for diagnosing version mismatches. Requires passing `contentFS` to check path. |
