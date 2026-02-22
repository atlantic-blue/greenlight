---
name: gl:status
description: Show project progress from test results and contract coverage
allowed-tools: [Read, Bash, Glob, Grep]
---

# Greenlight: Status

Show the user exactly where things stand using measured data — not opinions.

**Read:**
- `.greenlight/STATE.md`
- `.greenlight/GRAPH.json`
- `.greenlight/config.json`

## Gather Data

```bash
# Run test suite, capture results
{config.test.command} 2>&1 || true

# Count test files
find tests/integration -name "*.test.*" 2>/dev/null | wc -l
find tests/security -name "*.test.*" 2>/dev/null | wc -l

# Git status
git log --oneline -5
git status --short
```

## Display

Parse STATE.md to check for wrapped boundaries. Count:
- Total boundaries with `status: wrapped` or `status: refactored`
- Locking tests count from `locking_tests: [...]` array length
- Known issues count from `known_issues: [...]` array length

**If no wrapped boundaries exist:** Show greenfield-only display:

```
┌───────────────────────────────────────────────────────────┐
│  GREENLIGHT STATUS                                        │
├───────────────────────────────────────────────────────────┤
│                                                           │
│  Tests:  {pass} passing  {fail} failing  {skip} skipped   │
│  Slices: {done}/{total}  [{progress_bar}]                 │
│                                                           │
│  1. {name}                complete  {N} tests ({S} sec)   │
│  2. {name}                complete  {N} tests ({S} sec)   │
│  3. {name}                failing   {N}/{M} passing       │
│  4. {name}                blocked   (needs 3)             │
│  5. {name}                ready                           │
│                                                           │
│  Next: /gl:slice 3 (fix failing)                          │
│    or: /gl:slice 5 (no dependencies)                      │
│                                                           │
│  Last commit: {hash} {msg} ({time ago})                   │
│  Uncommitted: {Y/N}                                       │
│  Mode: {interactive/yolo}                                 │
│                                                           │
│  Product view: /gl:roadmap | History: /gl:changelog       │
└───────────────────────────────────────────────────────────┘
```

**If wrapped boundaries exist:** Show brownfield display:

```
┌───────────────────────────────────────────────────────────┐
│  GREENLIGHT STATUS                                        │
├───────────────────────────────────────────────────────────┤
│                                                           │
│  Tests:  {pass} passing  {fail} failing  {skip} skipped   │
│  Slices: {done}/{total}  [{progress_bar}]                 │
│                                                           │
│  1. {name}                complete  {N} tests ({S} sec)   │
│  2. {name}                failing   {N}/{M} passing       │
│  3. {name}                blocked   (needs 2)             │
│                                                           │
│  Wrapped Boundaries:                                      │
│    auth                   wrapped   12 locking tests      │
│                           2 known issues                  │
│    payments               wrapped    8 locking tests      │
│                           0 known issues                  │
│    users                  refactored (replaced by slice 1)│
│                                                           │
│  Wrap: 2/6 boundaries wrapped, 1 refactored               │
│                                                           │
│  Next: /gl:slice 2 (fix failing)                          │
│    or: /gl:wrap (wrap next boundary)                      │
│                                                           │
│  Last commit: {hash} {msg} ({time ago})                   │
│  Uncommitted: {Y/N}                                       │
│  Mode: {interactive/yolo}                                 │
│                                                           │
│  Product view: /gl:roadmap | History: /gl:changelog       │
└───────────────────────────────────────────────────────────┘
```

## Intelligence

Route based on current state:

| Situation | Recommendation |
|-----------|---------------|
| Failing tests | Fix first: `/gl:slice {N}` or `/gl:quick` |
| Wrapped boundaries with issues | "Boundary {name} has {N} known issues — consider `/gl:slice` to refactor" |
| Unwrapped boundaries exist | "Run `/gl:wrap` to wrap {N} remaining boundaries" |
| Multiple slices ready | Suggest parallel: "Slices {A} and {B} can run simultaneously" |
| All slices complete | `/gl:ship` |
| Uncommitted changes | Warn — may need to commit or stash |
| No tests exist | `/gl:init` |
| Continue file exists | `/gl:resume` |
| No .greenlight/ dir | `/gl:init` to get started |

No narrative. The table IS the status. One-line recommendation only.

---

## File-Per-Slice State Integration (C-83)

This section documents how /gl:status adapts its state reads when the project uses file-per-slice format (C-80). Legacy format behaviour is completely unchanged — if the state format is legacy, all reads go to STATE.md as before.

### State Format Detection

Detect state format (C-80) before performing any state reads:

- If file-per-slice: read all slice files from `.greenlight/slices/` and aggregate
- If legacy: read STATE.md as before (no change)

### File-Per-Slice Read Path

When using file-per-slice format:

1. Read all .md files from `.greenlight/slices/`
2. Parse frontmatter from each slice file to extract: ID, Name, Status, Tests, security_tests, deps
3. Compute the slice table: ID | Name | Status | Tests | Security | Deps (sorted by slice ID in ascending order)
4. Compute progress: done/total slices (count of slices with status complete vs total)
5. Compute current: identify all in-progress slices
6. Compute Test Summary: sum of tests and security_tests across all slice files
7. Read `project-state.json` for overview, session, and blockers
8. Display the computed summary to the user
9. Regenerate STATE.md (D-34) after display

STATE.md regeneration must occur after display so the user always sees fresh data first.

### Legacy Fallback

If format is legacy: read STATE.md as before. Legacy format display is completely unchanged — no change to existing behaviour.

### Error Handling

| Error State | When | Behaviour |
|-------------|------|-----------|
| FormatDetectionFailure | Cannot determine state format | Report error. Suggest running /gl:init |
| NoSliceFiles | No .md files found in `.greenlight/slices/` | Display empty summary. Report "no slices found" |
| CorruptSliceFile | Slice file has invalid frontmatter or is malformed | Skip and warn. Continue with remaining files |
| ProjectStateReadFailure | Cannot read `project-state.json` | Warn but continue. Display summary without overview |

### Invariants

- Status is computed fresh on every invocation — data is never cached
- Slice table is sorted by slice ID in ascending order
- Legacy format behaviour is completely unchanged
- Regeneration of STATE.md happens after display to the user (display first, then regenerate)
