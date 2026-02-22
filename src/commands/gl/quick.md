---
name: gl:quick
description: Ad-hoc task (bug fix, small feature, config) with test-first guarantees
allowed-tools: [Read, Write, Bash, Glob, Grep, Task, AskUserQuestion]
---

# Greenlight: Quick Task

Fast mode for ad-hoc work. Still test-first.

**Read CLAUDE.md.** Standards apply even in quick mode.
**Read .greenlight/STATE.md** for current project state.
**Read .greenlight/config.json** for test commands.

## Model Resolution

Before spawning any agent, resolve its model from `.greenlight/config.json`:

1. Check `model_overrides[agent_name]` — if set, use it
2. Else check `profiles[model_profile][agent_name]` — use profile default
3. Else fall back to `sonnet`

Agents spawned by this command: `debugger`.

## Flow

Ask: "What do you need to do?"

Classify and route:

### Bug Fix
1. Spawn debugger agent to investigate and reproduce
2. Debugger writes failing test capturing the bug
3. Debugger fixes the code to make test pass
4. Run full suite — no regressions
5. Commit: `fix: {description}`

```
Task(prompt="
Read agents/gl-debugger.md
Read CLAUDE.md
Read references/deviation-rules.md

<bug>
{user's description of the bug}
</bug>

<test_command>
{config.test.command}
</test_command>

<project_context>
{relevant files, recent changes}
</project_context>

Investigate, write failing test, fix.
", subagent_type="gl-debugger", model="{resolved_model.debugger}", description="Debug: {short description}")
```

### Small Feature (< 1 hour)
1. Define contract for what's being added (you do this, not the architect)
2. Write tests from contract (keep it simple — no separate agent for small features)
3. Implement until green
4. Quick security scan on diff (if `config.workflow.security_scan` is true)
5. Commit: `feat: {description}`

**When to refuse:** If the feature touches multiple existing contracts or needs new database tables → it's not a quick task:
```
This looks bigger than a quick task — touches {N} contracts / {M} boundaries.
Run /gl:add-slice to add it to the graph, then /gl:slice to build it properly.
```

### Config / Chore / Refactor
1. Run existing tests — confirm all green
2. Make the change
3. Run tests again — confirm nothing broke
4. Commit: `chore:` / `refactor:` / `docs:`

### Exploration / Spike
1. Create branch: `git checkout -b spike/{description}`
2. Prototype freely — standards relaxed, no tests required
3. Report findings — what worked, what didn't, what you'd recommend
4. **Don't merge** — spikes are throwaway. If the spike proves viable, create a proper slice:
```
Spike complete. Findings: {summary}

If you want to build this properly:
1. /gl:add-slice — define contracts for what the spike proved
2. /gl:slice {N} — build it with TDD
```

## Tracking

After completion, log in `.greenlight/QUICK.md` (create if doesn't exist):

```markdown
# Quick Tasks

| # | Type | Description | Date | Commit | Tests Added |
|---|------|-------------|------|--------|-------------|
| 1 | fix  | Login timeout on slow networks | 2026-02-07 | abc1234 | +2 |
| 2 | feat | Rate limit header in responses | 2026-02-07 | def5678 | +3 |
```

Update STATE.md test summary with new test counts.

## Quick Summary Generation (C-43)

After quick task completes, generate a summary (non-blocking):

### Summary Generation

Spawn a Task with fresh context to write the quick summary:

```
Task(prompt="
Collect and document the following information for quick task:

<quick_data>
Type: {bug_fix | small_feature | config | refactor}
Description: {user's description}
timestamp: {ISO8601 timestamp of completion}
Tests added: {N}
Files modified: {list}
Commit: {commit_hash}
Decision made: {yes/no}
</quick_data>

Write a summary to `.greenlight/summaries/quick-{timestamp}-SUMMARY.md`.

Summary failure does not block quick task completion.
", subagent_type="gl-summarizer", model="{resolved_model.summarizer}", description="Generate quick task summary")
```

If Task fails, log warning and continue. Summary generation is non-blocking.

### DECISIONS.md Append

If a decision was made during the quick task:

1. Format as DECISIONS.md table row
2. append to DECISIONS.md (create with header if doesn't exist)
3. Source format: `quick:{timestamp}`

Decision append is non-blocking. If it fails, log warning and continue.

---

## File-Per-Slice State Integration (C-84)

This section documents how /gl:quick adapts its state writes when the project uses file-per-slice format (C-80). Legacy format behaviour is completely unchanged — if the state format is legacy, all writes go to STATE.md as before.

### State Format Detection

Detect state format (C-80) before performing any state reads or writes:

- If file-per-slice: update test counts in only the relevant slice file from `.greenlight/slices/`, regenerate STATE.md (D-34)
- If legacy: update STATE.md as before (no change)

### File-Per-Slice Write Path

When using file-per-slice format:

1. Identify the relevant slice file at `.greenlight/slices/{id}.md`
2. Update test counts in only the relevant slice file — not all files, only the specific slice being updated
3. Regenerate STATE.md (D-34) after writing to the slice file

Only the relevant slice file is updated. /gl:quick does not touch other slice files.

### Legacy Fallback

If format is legacy: update STATE.md as before. Legacy format behaviour is completely unchanged — no change.

### Error Handling

| Error State | When | Behaviour |
|-------------|------|-----------|
| FormatDetectionFailure | Cannot determine state format | Report error. Suggest running /gl:init |
| SliceFileNotFound | Target slice file does not exist in `.greenlight/slices/` | Report error. Cannot update test counts without slice file |
| RegenerationFailure | STATE.md regeneration fails | Warn but continue. Slice file is still correct |

### Crash Safety (NFR-4)

All state writes use write-to-temp-then-rename (atomic writes) to prevent corruption:

- Slice file writes: write to temp file, then rename to target path (atomic)
- STATE.md regeneration: write to temp file, then rename to STATE.md (atomic)

### Invariants

- Only the relevant slice file is updated (not all slice files)
- STATE.md is regenerated after every state write (D-34)
- Legacy format behaviour is completely unchanged
- Crash safety via write-to-temp-then-rename on all writes (NFR-4)
