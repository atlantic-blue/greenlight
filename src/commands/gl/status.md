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
