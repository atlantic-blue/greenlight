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

## Intelligence

Route based on current state:

| Situation | Recommendation |
|-----------|---------------|
| Failing tests | Fix first: `/gl:slice {N}` or `/gl:quick` |
| Multiple slices ready | Suggest parallel: "Slices {A} and {B} can run simultaneously" |
| All slices complete | `/gl:ship` |
| Uncommitted changes | Warn — may need to commit or stash |
| No tests exist | `/gl:init` |
| Continue file exists | `/gl:resume` |
| No .greenlight/ dir | `/gl:init` to get started |

No narrative. The table IS the status. One-line recommendation only.
