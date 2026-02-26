# Slice State Template

Template for `.greenlight/slices/{slice-id}.md` — per-slice state file.

This template is read-only at runtime. Commands read this template to understand the schema and write slice state files. The template itself is never modified during execution.

---

## Schema Definition

Each slice state file begins with YAML frontmatter between `---` delimiters. All fields are flat key-value pairs (no nesting, per D-30).

### Fields

| Field | Type | Required/Optional | Description |
|-------|------|-------------------|-------------|
| `id` | string | required | Slice identifier. Must match `S-{N}` or `S-{N}.{N}` |
| `status` | enum | required | Current lifecycle status (8 values, closed enum) |
| `step` | enum | required | Current TDD loop step (7 values, closed enum) |
| `milestone` | string | required | Milestone slug this slice belongs to |
| `started` | ISO date | optional | Date when work began (`YYYY-MM-DD`) |
| `updated` | ISO timestamp | required | Last modification timestamp (`YYYY-MM-DDTHH:MM:SSZ`) |
| `tests` | integer | required | Count of passing functional tests |
| `security_tests` | integer | required | Count of passing security tests |
| `session` | string | optional | Advisory session identifier (ISO timestamp + hyphen + random alphanumeric suffix) |
| `deps` | comma-separated | optional | Dependency slice IDs (e.g., `S-1, S-2`) |

All field names are lowercase with underscores (snake_case). No camelCase, no nesting.

---

## Status Lifecycle

The `status` field tracks where a slice is in the overall pipeline. The `step` field tracks where work is within the current active phase.

### Status Values (closed enum — only these 8 values are valid)

| Status | Meaning |
|--------|---------|
| `pending` | Not started, dependencies may not be met |
| `ready` | Dependencies met, can be claimed |
| `tests` | Test writer is working or tests are written |
| `implementing` | Implementer is working |
| `security` | Security scan in progress |
| `fixing` | Fixing security or test failures |
| `verifying` | Verifier reviewing results |
| `complete` | All tests green, security clean, accepted |

### Step Values (closed enum — only these 7 values are valid)

| Step | Meaning |
|------|---------|
| `none` | No active step (pending or ready) |
| `tests` | Writing tests |
| `implementing` | Implementing code |
| `security` | Running security scan |
| `fixing` | Fixing failures |
| `verifying` | Verifying output |
| `complete` | Slice complete |

### Transition Rules

Status transitions follow the TDD loop:

```
pending → ready → tests → implementing → security → fixing → verifying → complete
```

Allowed transitions:

- `pending` → `ready` (when all deps are complete)
- `ready` → `tests` (when claimed by test writer)
- `tests` → `implementing` (when tests are written and failing)
- `implementing` → `security` (when tests pass)
- `security` → `fixing` (when security issues found)
- `security` → `verifying` (when security is clean)
- `fixing` → `implementing` (when fixes need re-implementation)
- `fixing` → `verifying` (when fixes pass security)
- `verifying` → `complete` (when verifier approves)
- `verifying` → `fixing` (when verifier rejects)

Backward transitions are only allowed through `fixing`. No other backward transitions are valid.

---

## Session Tracking

The `session` field is advisory only (per D-33). It is not a lock. Multiple agents may write to the same slice; the session field provides a hint, not a guarantee.

### Format

Session identifiers use ISO timestamp + hyphen + random alphanumeric suffix:

```
2024-01-15T10:30:00Z-x7k2m9
```

The suffix is at least 6 random alphanumeric characters. This ensures uniqueness across concurrent agents.

### Lifecycle

- Set on claim: when an agent claims a slice, it writes its session identifier to the `session` field.
- Cleared on completion: when a slice reaches `complete` status, the session field is cleared (empty string or omitted).

### Advisory Semantics

The session field does not prevent concurrent access. Commands warn before claiming a slice that has a non-empty session field, but they do not block. If two agents claim the same slice simultaneously, the last write wins. Session tracking is a coordination hint, not a distributed lock.

---

## Body Sections

After the frontmatter, the body of each slice state file follows this structure:

```markdown
# {slice-id}: {slice-name}

## Why

Reason this slice exists — the user value it delivers.

## What

Deliverables: what code, files, or changes this slice produces.

## Dependencies

Which slices must be complete before this one can start.
List slice IDs or "none".

## Contracts

Which contracts (C-{N}) define the behaviour this slice implements.

## Decisions

Architectural or design decisions made during this slice.
Populated during implementation.

## Files

Files created or modified by this slice.
Populated during implementation.
```

Each section is separated by a blank line. The heading uses `# {slice-id}: {slice-name}` where `{slice-id}` is the ID from frontmatter (e.g., `S-12`) and `{slice-name}` is a short descriptive title.

---

## File Naming

Slice state files are stored in `.greenlight/slices/`. One file per slice. Files are named `{slice-id}.md`:

```
.greenlight/slices/S-1.md
.greenlight/slices/S-2.md
.greenlight/slices/S-12.md
.greenlight/slices/S-3.1.md
```

The slice ID in the filename must exactly match the `id` field in the frontmatter. IDs must match `S-{digits}` or `S-{digits}.{digits}`. This validation prevents path traversal attacks — an ID like `S-../../etc/passwd` is rejected before any file operation.

---

## Errors

### InvalidSliceId

```
Invalid slice ID: {id}. Must match S-{N} or S-{N}.{N}
```

Raised when the `id` field does not match the pattern `S-{digits}` or `S-{digits}.{digits}`. The pattern allows `S-1`, `S-28`, `S-3.1`, `S-12.4`. It rejects anything with path separators, dots in unexpected positions, or non-numeric segments. This prevents path traversal via `../` in slice IDs.

### InvalidStatus

```
Invalid status: {value}. Valid values: pending, ready, tests, implementing, security, fixing, verifying, complete
```

Raised when the `status` field contains a value not in the closed enum of 8 valid values.

### InvalidFrontmatter

```
Invalid frontmatter format. Expected flat key: value pairs between --- delimiters
```

Raised when the file does not begin with `---`, does not contain a closing `---`, or contains nested YAML structures. Only flat key: value pairs are valid.

---

## Invariants

1. **Template read-only at runtime.** This file is never modified by commands. It is a schema reference only.
2. **Frontmatter flat key-value only.** No nested objects, no arrays, no multi-line values. All values are scalars.
3. **All field names lowercase with underscores.** `security_tests` not `securityTests` or `SecurityTests`.
4. **Status enum closed.** Only the 8 defined status values are valid. No extensions without updating this template.
5. **Step enum closed.** Only the 7 defined step values are valid.
6. **Session format.** Session identifiers must be ISO timestamp + hyphen + random alphanumeric suffix.
7. **Slice ID path traversal prevention.** Slice IDs are validated against `S-{digits}` or `S-{digits}.{digits}` before use in file paths. IDs containing `../` or other path traversal sequences (e.g., `S-../../etc/passwd`) are rejected with InvalidSliceId.

---

## Examples

### Example: Pending State

A slice that has not yet been claimed:

```markdown
---
id: S-5
status: pending
step: none
milestone: auth
updated: 2024-01-15T08:00:00Z
tests: 0
security_tests: 0
---

# S-5: User Login

## Why

Users need to authenticate to access protected resources.

## What

Implement POST /v1/auth/login endpoint that validates credentials and returns a JWT.

## Dependencies

S-4 (User Registration must be complete before login can be tested)

## Contracts

C-8, C-9

## Decisions

None yet.

## Files

None yet.
```

### Example: Implementing State

A slice actively being implemented:

```markdown
---
id: S-5
status: implementing
step: implementing
milestone: auth
started: 2024-01-15
updated: 2024-01-15T10:45:00Z
tests: 8
security_tests: 0
session: 2024-01-15T10:30:00Z-x7k2m9
deps: S-4
---

# S-5: User Login

## Why

Users need to authenticate to access protected resources.

## What

Implement POST /v1/auth/login endpoint that validates credentials and returns a JWT.

## Dependencies

S-4

## Contracts

C-8, C-9

## Decisions

JWT expiry set to 24h per C-8 field `expires_in`.

## Files

- src/auth/login.ts (created)
- src/auth/jwt.ts (modified)
```

### Example: Complete State

A slice that has finished successfully:

```markdown
---
id: S-5
status: complete
step: complete
milestone: auth
started: 2024-01-15
updated: 2024-01-16T14:20:00Z
tests: 12
security_tests: 3
deps: S-4
---

# S-5: User Login

## Why

Users need to authenticate to access protected resources.

## What

Implemented POST /v1/auth/login with credential validation, JWT issuance, and rate limiting.

## Dependencies

S-4

## Contracts

C-8, C-9

## Decisions

JWT expiry set to 24h. Rate limit: 10 attempts per minute per IP.

## Files

- src/auth/login.ts (created, 48 lines)
- src/auth/jwt.ts (modified, added sign helper)
- src/auth/rate-limit.ts (created, 22 lines)
```
