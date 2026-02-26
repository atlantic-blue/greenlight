---
name: gl:pause
description: Create handoff file when stopping work mid-slice
allowed-tools: [Read, Write, Bash, Glob, Grep]
---

# Greenlight: Pause

Create `.greenlight/.continue-here.md` to preserve work state across sessions.

**Read:**
- `.greenlight/STATE.md` — current position
- `.greenlight/config.json` — test commands

## Gather State

```bash
# Current test status
{config.test.command} 2>&1 || true

# Uncommitted changes
git status --short

# Recent commits this session
git log --oneline -10
```

1. Current slice and step (tests written? implementing? security scan?)
2. Work completed this session
3. Work remaining
4. Decisions made and rationale
5. Blockers or issues
6. Files modified but not committed
7. Test status (which pass, which fail)

## Write Handoff

```markdown
---
slice: {slice_id}
slice_name: {name}
step: {pending | tests | implementing | security | fixing | verifying}
paused_at: {ISO timestamp}
---

## Current State
{Where exactly are we? What was the last completed action?}

## Test Status
- Passing: {N}
- Failing: {N}
- Failing tests:
  - {test name}: {failure reason}

## Completed This Session
- [x] {completed item}
- [x] {completed item}
- [ ] {incomplete item}

## Remaining Work
1. {specific next action}
2. {what else needs to happen}

## Decisions Made
- {decision}: {rationale}

## Deviations
- [{type}] {description}

## Blockers
- {blocker}: {status/workaround}

## Next Action
Start with: {specific first thing to do when resuming}
Agent to spawn: {gl-implementer | gl-security | gl-verifier | none}
Context needed: {what the agent needs to know}
```

## Update STATE.md

Update the Session section:
```markdown
## Session
Last session: {now}
Resume file: .greenlight/.continue-here.md
```

## Commit

```bash
git add .greenlight/.continue-here.md
git add .greenlight/STATE.md
git commit -m "wip({slice_id}): paused at {step}

Tests: {N} passing, {N} failing
Resume with /gl:resume
"
```

## Confirm

```
Paused

Slice: {slice_id} — {slice_name}
Step: {step}
Tests: {N} passing, {N} failing

Resume with: /gl:resume
```

---

## File-Per-Slice State Integration (C-84)

This section documents how /gl:pause adapts its state writes when the project uses file-per-slice format (C-80). Legacy format behaviour is completely unchanged — if the state format is legacy, all writes go to STATE.md as before.

### State Format Detection

Detect state format (C-80) before performing any state reads or writes:

- If file-per-slice: write pause state to own slice file, write resume context to `project-state.json`, regenerate STATE.md (D-34)
- If legacy: write to STATE.md as before (no change)

### File-Per-Slice Write Path

When using file-per-slice format:

1. Write pause state (step, paused_at, session context) to the own slice file at `.greenlight/slices/{id}.md` — only writes to its own slice file, never to another slice's file
2. Write resume context to `project-state.json` (slice, step, paused_at, next action, context needed)
3. Regenerate STATE.md (D-34) after writing to the slice file

### Legacy Fallback

If format is legacy: write to STATE.md as before. Legacy format behaviour is completely unchanged — no change.

### Error Handling

| Error State | When | Behaviour |
|-------------|------|-----------|
| FormatDetectionFailure | Cannot determine state format | Report error. Suggest running /gl:init |
| SliceFileNotFound | Target slice file does not exist in `.greenlight/slices/` | Create it. Warn |
| RegenerationFailure | STATE.md regeneration fails | Warn but continue. Slice file and project-state.json are still correct |

### Crash Safety (NFR-4)

All state writes use write-to-temp-then-rename (atomic writes) to prevent corruption:

- Slice file writes: write to temp file, then rename to target path (atomic)
- `project-state.json` writes: write to temp file, then rename (atomic)
- STATE.md regeneration: write to temp file, then rename to STATE.md (atomic)

### Invariants

- /gl:pause writes ONLY to its own slice file (never to another slice's file)
- Resume context is always written to `project-state.json` in file-per-slice mode
- STATE.md is regenerated after every state write (D-34)
- Legacy format behaviour is completely unchanged
- Crash safety via write-to-temp-then-rename on all writes (NFR-4)
