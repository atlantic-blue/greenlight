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
