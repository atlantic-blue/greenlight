---
name: gl:add-slice
description: Add a new vertical slice to the dependency graph
allowed-tools: [Read, Write, Bash, Glob, Grep, Task, AskUserQuestion]
---

# Greenlight: Add Slice

Add a new vertical slice to an existing project.

**Read first:**
- `CLAUDE.md` — engineering standards
- `.greenlight/STATE.md` — current state
- `.greenlight/GRAPH.json` — existing dependency graph
- `.greenlight/CONTRACTS.md` — existing contracts
- `.greenlight/config.json` — settings

## Flow

### 1. Understand the New Slice

Ask: "What should the user be able to do?"

Extract a clear user action: "User can {verb} {object}" — e.g., "User can reset their password via email."

### 2. Check Existing Contracts

Does this new slice:
- Need entirely new contracts? (new boundary)
- Extend existing contracts? (new method on existing service)
- Reuse existing contracts? (new frontend using existing API)

If extending existing contracts → flag which completed slices might be affected:
```
This slice extends the UserService contract (used in slices 1, 2).
Changes won't break existing tests, but I'll note the dependency.
```

### 3. Brief Discussion (if needed)

For simple additions, skip discussion. For ambiguous ones, ask 2-3 clarifying questions:
- Where does this fit in the user flow?
- Any new boundaries (new external service, new DB table)?
- Dependencies on existing slices?

### 4. Spawn Architect

```
Task(prompt="
Read agents/gl-architect.md
Read CLAUDE.md

<existing_contracts>
{current CONTRACTS.md content}
</existing_contracts>

<existing_graph>
{current GRAPH.json content}
</existing_graph>

<new_slice>
User action: {what user can do}
Description: {details from discussion}
</new_slice>

<decisions>
{any decisions from discussion}
</decisions>

Produce:
1. New/modified contracts for this slice
2. Updated GRAPH.json with new slice added (next available ID, correct dependencies)
3. Flag any existing slices affected by contract changes
", subagent_type="gl-architect", model="{config.models.architect}", description="Add slice: {name}")
```

### 5. User Review

Present the new contracts and updated graph. Show:
- What's new
- What's changed (if any existing contracts modified)
- Where it fits in the dependency graph

### 6. Apply

- Update `.greenlight/CONTRACTS.md` with new/modified contracts
- Update `.greenlight/GRAPH.json` with new slice
- Update `.greenlight/STATE.md` with new slice row (status: pending)

### 7. Commit

```bash
git add .greenlight/CONTRACTS.md .greenlight/GRAPH.json .greenlight/STATE.md
git commit -m "docs: add slice {id} — {name}

Contracts: {new contracts}
Dependencies: {deps or 'none'}
"
```

### 8. Next

```
Slice added: {id} — {name}

New contracts: {list}
Dependencies: {list or 'none'}
Status: {ready | blocked by {deps}}

{if ready: "Run /gl:slice {id} to build it."}
{if blocked: "Complete {deps} first, then /gl:slice {id}."}
```
