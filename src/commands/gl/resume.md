---
name: gl:resume
description: Resume work from previous session with full context restoration
allowed-tools: [Read, Write, Bash, Glob, Grep, Task, AskUserQuestion]
---

# Greenlight: Resume

Restore project state and continue from where you left off.

**Read first:**
- `CLAUDE.md` — engineering standards
- `.greenlight/STATE.md` — project state
- `.greenlight/config.json` — settings

## Step 1: Load State

Check for handoff file:

```bash
cat .greenlight/.continue-here.md 2>/dev/null
cat .greenlight/STATE.md 2>/dev/null
```

| Scenario | Action |
|----------|--------|
| `.continue-here.md` exists | Read and internalise: slice, step, test status, decisions, blockers, next action |
| Missing but `STATE.md` exists | Check for incomplete slices. Run tests to determine actual state |
| Neither exists | No project found. Suggest `/gl:init` |

## Step 2: Verify Current State

```bash
# Run tests to confirm actual status
{config.test.command} 2>&1 || true

# Check git status
git status --short
git log --oneline -5
```

Compare test results with handoff file. If they disagree, trust the test results.

Report discrepancies:
```
Note: Handoff says {N} passing, but test run shows {M} passing.
{possible explanation — e.g., "code may have been modified since pause"}
```

## Step 3: Present Status

```
┌─────────────────────────────────────────────────┐
│  GREENLIGHT: RESUMING                           │
├─────────────────────────────────────────────────┤
│                                                 │
│  Slice: {id} — {name}                           │
│  Step:  {tests | implementing | security | ...} │
│  Tests: {N} passing  {N} failing                │
│                                                 │
│  Last session: {date from handoff}              │
│  Next action: {from handoff or inferred}        │
│                                                 │
│  Blockers: {from handoff or "none"}             │
└─────────────────────────────────────────────────┘
```

## Step 4: Route to Next Action

Based on where we paused:

| Paused At | Resume With |
|-----------|-------------|
| Tests written, not implementing | Spawn implementer for this slice |
| Mid-implementation, some tests passing | Spawn fresh implementer with remaining failing test names |
| Implementation complete, no security scan | Run security scan |
| Security issues found, not fixed | Spawn implementer to fix security tests |
| Security fixed, not verified | Run verifier |
| Verification failed | Address verification issues |
| Visual checkpoint pending | Present checkpoint to user |

Load the context the next agent needs from the handoff file's "Context needed" section.

## Step 5: Clean Up

After routing to the next action, remove the handoff file:

```bash
rm .greenlight/.continue-here.md
git add .greenlight/.continue-here.md
git commit -m "chore({slice_id}): resumed from pause"
```

Update STATE.md session section:
```markdown
## Session
Last session: {now}
Resume file: None
```

---

## File-Per-Slice State Integration (C-84)

This section documents how /gl:resume adapts its state reads when the project uses file-per-slice format (C-80). Legacy format behaviour is completely unchanged — if the state format is legacy, all reads go to STATE.md as before.

### State Format Detection

Detect state format (C-80) before performing any state reads:

- If file-per-slice: read slice files from `.greenlight/slices/` to determine resumable state, read resume context from `project-state.json`
- If legacy: read STATE.md as before (no change)

### File-Per-Slice Read Path

When using file-per-slice format:

1. Read all slice files from `.greenlight/slices/` to identify in-progress or paused slices
2. Read resume context from `project-state.json` (slice, step, paused_at, next action, context needed)
3. If resume context exists in `project-state.json`, restore from it; otherwise infer state from slice files

### Legacy Fallback

If format is legacy: read STATE.md as before. Legacy format behaviour is completely unchanged — no change.

### Error Handling

| Error State | When | Behaviour |
|-------------|------|-----------|
| FormatDetectionFailure | Cannot determine state format | Report error. Suggest running /gl:init |
| SliceFileNotFound | Referenced slice file does not exist in `.greenlight/slices/` | Report error. Cannot resume without slice file |
| RegenerationFailure | STATE.md regeneration fails | Warn but continue. Slice files are still correct |

### Invariants

- Resume context from `project-state.json` takes precedence over inferred state from slice files
- Legacy format behaviour is completely unchanged
- Slice files in `.greenlight/slices/` are the source of truth for slice state
