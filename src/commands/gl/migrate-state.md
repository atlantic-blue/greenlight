---
name: gl:migrate-state
description: Converts legacy STATE.md-based projects to file-per-slice format. Migrate STATE.md content into individual slice files under .greenlight/slices/.
allowed-tools: [Read, Write, Bash, Glob, Grep]
---

# Greenlight: Migrate State

Converts a legacy `STATE.md` project to the file-per-slice format. This is a one-way, all-or-nothing migration — it cannot be undone automatically. Run it only when explicitly requested (D-32: Explicit only — no auto-migration).

## When to Use

Use this command when:
- The project has a `.greenlight/STATE.md` file in legacy format
- The project does NOT yet have a `.greenlight/slices/` directory
- The user explicitly requests migration to file-per-slice format

Do NOT run this automatically. No auto-migration. This command is invoked only on explicit user request (D-32).

---

## Flow

### Step 1 — Verify STATE.md Exists

Check that `.greenlight/STATE.md` exists in the project root.

```
if STATE.md does not exist:
    print: "No STATE.md found. Nothing to migrate."
    exit with error: NoStateMd
```

### Step 2 — Verify slices/ Directory Does Not Exist

Check that `.greenlight/slices/` does not already exist.

```
if .greenlight/slices/ exists:
    print: "Already using file-per-slice format. Nothing to migrate."
    exit with error: AlreadyMigrated
```

If `AlreadyMigrated`, STATE.md is untouched and the process stops cleanly.

### Step 3 — Parse STATE.md

Read STATE.md and extract all structured content:

1. **Extract slice table** — Parse table rows from the `## Current` section. For each row, extract:
   - `ID` — e.g. `S-1`, `S-42`
   - `Name` — the slice name
   - `Status` — e.g. `done`, `in-progress`, `pending`
   - `Tests` — test status (pass/fail/pending)
   - `Security` — security review status
   - `Deps` — slice dependencies

2. **Extract Current Section** — The full `## Current` section content (active slice details)

3. **Extract Overview** — The `## Overview` section (value prop, stack, mode, project description)

4. **Extract Session** — The `## Session` section (current session notes and context)

5. **Extract Blockers** — The `## Blockers` section (active blockers list)

6. **Extract Decisions** — The `## Decisions` section (architectural decisions log)

If parsing fails and the structure cannot be extracted:
```
exit with error: ParseFailure
message: "failed to parse STATE.md — cannot extract required sections"
```

Do not proceed if ParseFailure occurs. STATE.md remains intact.

### Step 4 — Create .greenlight/slices/ Directory

Create the slices/ directory with permissions 0o755:

```bash
mkdir -p .greenlight/slices/   # permissions: 755
```

Equivalent to `os.MkdirAll(".greenlight/slices", 0o755)`.

### Step 5 — Write Individual Slice Files

For each slice row extracted from the slice table:

#### 5a — Validate Slice ID Format

Validate that the ID matches the pattern `S-{digits}` (e.g. `S-1`, `S-32`, `S-100`).

- ID must match: `^S-\d+$`
- Validate slice ID to prevent path traversal — the ID must not contain `..`, `/`, or other path-manipulation characters
- If the ID format is invalid: warn and skip that row (do not abort entire migration)
  - Error type: `InvalidSliceId`
  - Behaviour: log a warning, skip the invalid slice, continue with remaining slices

#### 5b — Create Slice File

For each valid slice, create `.greenlight/slices/{id}.md` with YAML frontmatter:

```markdown
---
id: S-{n}
name: {slice name}
status: {status}
tests: {test status}
security: {security status}
deps: [{dep list}]
step: {step number or empty}
milestone: {milestone or empty}
session: {session ref or empty}
---

# S-{n}: {Slice Name}

{any additional slice content from STATE.md}
```

The frontmatter must include: `id`, `name`, `status`, `tests`, `security`, `deps`, `step`, `milestone`, `session`.

#### 5c — Write Atomically (NFR-4)

Use write-to-temp-then-rename (atomic write) for each slice file:

1. Write content to a temp file (e.g. `.greenlight/slices/{id}.md.tmp`)
2. Use `os.Rename(tempPath, targetPath)` to atomically move it to the final path
3. Set file permissions to 0o644

This is the NFR-4 atomic write pattern. It ensures crash safety: if the process is interrupted, no partial file is left at the target path.

If any slice file write fails:
```
exit with error: PartialWriteFailure
abort migration — do not proceed to backup step
clean up: remove the .greenlight/slices/ directory (RemoveAll)
```

### Step 6 — Create project-state.json

Create `.greenlight/project-state.json` containing the non-slice sections extracted from STATE.md:

```json
{
  "overview": "{content from ## Overview section}",
  "session": "{content from ## Session section}",
  "blockers": "{content from ## Blockers section}",
  "decisions": "{content from ## Decisions section}"
}
```

Write this file atomically (write-to-temp-then-rename, NFR-4) with permissions 0o644.

### Step 7 — Backup STATE.md (Only AFTER All Files Succeed)

This step happens AFTER all files have been written successfully. STATE.md is untouched until this point.

Rename `STATE.md` to `STATE.md.backup`:

```
os.Rename(".greenlight/STATE.md", ".greenlight/STATE.md.backup")
```

**BackupExists handling:** If `STATE.md.backup` already exists, use a timestamp suffix to avoid overwriting an existing backup:

```
STATE.md.backup.{timestamp}
```

Example: `STATE.md.backup.20260222T143000`

Multiple backups are preserved with timestamped suffixes. No prior backup is ever overwritten.

If the backup rename fails:
```
exit with error: BackupRenameFailure
message: "cannot rename STATE.md to STATE.md.backup"
```

The original STATE.md remains intact if BackupRenameFailure occurs.

### Step 8 — Generate New STATE.md

Generate a new `STATE.md` that is driven by the slice files. This is a generated STATE.md — it must include a header comment indicating it is auto-generated and should not be edited directly:

```markdown
<!-- GENERATED — do not edit by hand. This file is auto-generated from .greenlight/slices/. Run /gl:status to regenerate. -->

# Project State
...
```

The new STATE.md is assembled by reading all slice files and regenerating the slice table. This is the standard generated STATE.md format used by all post-migration commands (D-31: after migration, all commands use file-per-slice automatically via detection).

Write the new STATE.md atomically (write-to-temp-then-rename, NFR-4).

### Step 9 — Report

Print a success report to the user:

```
Migrated {N} slices to file-per-slice format. Backup: STATE.md.backup
```

Example:
```
Migrated 12 slices to file-per-slice format. Backup: STATE.md.backup
```

Include any warnings about skipped slices (InvalidSliceId).

---

## Error Reference

| Error Code | Condition | Behaviour |
|------------|-----------|-----------|
| `NoStateMd` | STATE.md does not exist | Print "No STATE.md found. Nothing to migrate." Stop. |
| `AlreadyMigrated` | .greenlight/slices/ already exists | Print "Already using file-per-slice format. Nothing to migrate." Stop. |
| `ParseFailure` | failed to parse STATE.md — cannot extract sections | Report error. STATE.md untouched. Stop. |
| `InvalidSliceId` | Slice row has invalid ID format (not S-{digits}) | Warn and skip that slice. Continue migration. |
| `PartialWriteFailure` | A slice file write fails mid-migration | Abort. Remove partial .greenlight/slices/ directory. STATE.md intact. |
| `BackupRenameFailure` | Cannot rename STATE.md to STATE.md.backup | Report error. All slice files written. STATE.md still at original path. |
| `CleanupFailure` | RemoveAll of partial slices/ fails after PartialWriteFailure | Report error. Provide manual cleanup instructions: `rm -rf .greenlight/slices/` |

---

## Failure Recovery and Cleanup

### On PartialWriteFailure

If writing slice files fails partway through:

1. Abort immediately — do not proceed to backup or regeneration
2. Attempt cleanup: `os.RemoveAll(".greenlight/slices/")` to remove the partial slices/ directory
3. If cleanup succeeds: STATE.md is intact, project is unchanged
4. If cleanup fails (CleanupFailure): report the failure and provide manual removal instructions:
   ```
   Migration failed. Partial files may remain at .greenlight/slices/.
   Run: rm -rf .greenlight/slices/
   STATE.md is intact and unchanged.
   ```

### Crash Safety

All slice file writes use write-to-temp-then-rename (NFR-4 atomic write pattern). In the worst case — a crash or power loss mid-migration — the outcome is:

- Stale `.greenlight/slices/` files may exist alongside an intact `STATE.md`
- The project can recover: delete `.greenlight/slices/` and re-run `/gl:migrate-state`
- STATE.md is never corrupted or lost before the explicit backup rename in Step 7

---

## Invariants

| Invariant | Description |
|-----------|-------------|
| One-way | Migration is one-way. There is no automatic reverse. Backup is preserved for manual recovery. |
| All-or-nothing | The migration is atomic. If any step fails before the backup rename, STATE.md remains untouched and intact. |
| Explicit only (D-32) | No auto-migration. This command runs only when explicitly invoked. No automatic migration happens in other commands. |
| No dual-write (D-38) | After migration, writes go only to slice files. STATE.md is generated, never written to directly. No dual-write to both formats. |
| Post-migration detection (D-31) | After migration, all commands detect file-per-slice format automatically via presence of .greenlight/slices/ and switch to the new read/write path. |
| Slice files written first | All slice files and project-state.json are written BEFORE STATE.md is renamed to backup. Only AFTER all files succeed does the backup rename happen. |

---

## Security

- **Path traversal prevention:** Validate slice ID (must match `^S-\d+$`) before constructing any file path. Reject IDs containing `..`, `/`, or non-digit characters after `S-`. This prevents a malformed STATE.md from writing files outside `.greenlight/slices/`.
- **File permissions:** Directories created with 0o755. Files written with 0o644.
- **Does not modify source code:** This command operates exclusively on `.greenlight/` files. It does not modify source code, test files, configuration files, or any files outside `.greenlight/`.
- **No sensitive data logged:** Slice IDs, names, and statuses are project metadata — no passwords, tokens, or PII are processed.

---

## Scope

This command reads and writes only within `.greenlight/`:

- Reads: `.greenlight/STATE.md`
- Creates: `.greenlight/slices/{id}.md` (one per slice), `.greenlight/project-state.json`
- Renames: `.greenlight/STATE.md` → `.greenlight/STATE.md.backup` (or `.backup.{timestamp}`)
- Generates: `.greenlight/STATE.md` (new, GENERATED format)

Does not modify: source code, test files, `CONTRACTS.md`, `GRAPH.json`, `config.json`, or any file outside `.greenlight/`.
