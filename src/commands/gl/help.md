---
name: gl:help
description: Show all Greenlight commands and current state
allowed-tools: [Read, Bash, Glob]
---

# Greenlight: Help

```
┌──────────────────────────────────────────────────────────────┐
│  GREENLIGHT v1.0 — TDD-first development for Claude Code    │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  SETUP                                                       │
│  /gl:init              Brief interview + project config      │
│  /gl:design            System design session → DESIGN.md     │
│  /gl:map               Analyse existing codebase first       │
│  /gl:settings          Configure models, mode, options       │
│                                                              │
│  BUILD                                                       │
│  /gl:slice <N>         TDD loop: test → implement →          │
│                        security → verify → commit            │
│  /gl:quick             Ad-hoc task with test guarantees      │
│  /gl:add-slice         Add new slice to graph                │
│                                                              │
│  MONITOR                                                     │
│  /gl:status            Real progress from test results       │
│  /gl:pause             Save state for next session           │
│  /gl:resume            Restore and continue                  │
│                                                              │
│  SHIP                                                        │
│  /gl:ship              Full audit + deploy readiness         │
│                                                              │
│  FLOW                                                        │
│  map? → init → design → slice 1 → slice 2 → ... → ship      │
│                                                              │
│  Tests are the source of truth. Green means done.            │
│  Security is built in, not bolted on.                        │
└──────────────────────────────────────────────────────────────┘
```

## Context-Aware Help

Check for project state:

```bash
cat .greenlight/STATE.md 2>/dev/null
```

**If project exists:** Show current progress inline:
```
Current: Slice {N} — {name} ({step})
Tests: {pass} passing, {fail} failing
Next: /gl:slice {N}
```

**If no project:**
```
No project found. Run /gl:init to get started.
For existing codebases, run /gl:map first.
```
