---
name: gl:roadmap
description: Display project roadmap, plan milestones, and archive completed work
argument-hint: "[milestone|archive]"
allowed-tools: [Read, Write, Bash, Glob, Grep, Task, AskUserQuestion]
---

# Greenlight: Roadmap Management

Manage your project roadmap across three modes: display, milestone planning, and archiving.

**Read first:**
- `CLAUDE.md` — engineering standards
- `.greenlight/config.json` — project context and settings
- `.greenlight/ROADMAP.md` — product roadmap (read-only in display mode)

## Mode Selection

This command operates in three modes based on arguments:

1. **Display** (no arguments): `/gl:roadmap` — read-only display of ROADMAP.md
2. **Milestone Planning**: `/gl:roadmap milestone` — plan a new milestone with gl-designer
3. **Archive**: `/gl:roadmap archive` — archive completed milestones

---

## Mode 1: Display (Read-Only)

**Usage:** `/gl:roadmap`

Display the current project roadmap without making any modifications. This is a read-only operation. Do not modify any files. Do not create new files. Do not write anything.

### Prerequisites

Check that the project has been initialized:

```bash
cat .greenlight/config.json 2>/dev/null
```

**If no config found:**
```
No config found. Run /gl:init first to set up the project.
```
Stop here.

### Load Roadmap

```bash
cat .greenlight/ROADMAP.md 2>/dev/null
```

**Error handling:**

- **If ROADMAP.md does not exist:**
  ```
  No roadmap found. Run /gl:design to create one.
  ```
  Stop here.

- **If ROADMAP.md is empty:**
  ```
  ROADMAP.md is empty. Run /gl:design to populate it.
  ```
  Stop here.

### Read Project Context

Read project name and context from config.json to provide context for the roadmap display:

```bash
cat .greenlight/config.json
```

### Display Components

Present the following sections from ROADMAP.md:

1. **Architecture Diagram** — Mermaid diagram showing system architecture and component relationships
2. **Milestone Tables** — Tables linking slices to product milestones with status tracking
3. **Archived Milestones** — Previously completed and compressed milestone records

Report to user:
```
Project: {project_name}

ROADMAP.md contents:

{full contents of ROADMAP.md}
```

This is a read-only display. No modifications are made to any files.

**Design prerequisite:** If the user wants to modify the roadmap structure, suggest running `/gl:design` to revise the design and regenerate ROADMAP.md.

---

## Mode 2: Milestone Planning

**Usage:** `/gl:roadmap milestone`

Plan a new milestone by spawning gl-designer for a lighter design session focused on the milestone scope.

### Prerequisites

Verify project context exists:

```bash
cat .greenlight/config.json 2>/dev/null
cat .greenlight/ROADMAP.md 2>/dev/null
cat .greenlight/DESIGN.md 2>/dev/null
cat .greenlight/CONTRACTS.md 2>/dev/null
cat .greenlight/STATE.md 2>/dev/null
```

**Error handling:**

- **If no config found:** "No config found. Run /gl:init first."
- **If no roadmap found:** "No roadmap found. Run /gl:design to create one."
- **If no design found:** "No design found. Run /gl:design first."

Stop if any prerequisite is missing.

### Model Resolution

Before spawning gl-designer, resolve its model from `.greenlight/config.json`:

1. Check `model_overrides["designer"]` — if set, use it
2. Else check `profiles[model_profile]["designer"]` — use profile default
3. Else fall back to `sonnet`

### Gather Context

Read existing project context for the designer:

```bash
cat .greenlight/DESIGN.md
cat .greenlight/CONTRACTS.md
cat .greenlight/STATE.md
cat .greenlight/GRAPH.json 2>/dev/null
cat .greenlight/ASSESS.md 2>/dev/null
```

Build the context block from what's available.

### Spawn gl-designer for Milestone Planning

This is a lighter, scoped design session. Skip the init interview and skip stack decisions (already established). Focus on milestone-specific requirements and slicing.

```
Task(prompt="
Read agents/gl-designer.md
Read CLAUDE.md

<session_mode>
milestone_planning
</session_mode>

<existing_design>
{contents of DESIGN.md}
</existing_design>

<existing_contracts>
{contents of CONTRACTS.md}
</existing_contracts>

<existing_state>
{contents of STATE.md}
</existing_state>

<existing_graph>
{contents of GRAPH.json if exists, otherwise 'No graph yet'}
</existing_graph>

<existing_assessment>
{contents of ASSESS.md if exists, otherwise 'No assessment available'}
</existing_assessment>

<wrap_progress>
{Extract "Wrapped Boundaries" section from STATE.md if present. Otherwise: 'No wrapped boundaries yet'}
</wrap_progress>

<project_context>
{project name, stack, architecture from config.json}
</project_context>

This is a lighter session for milestone planning. Skip the init interview. Skip stack decisions (already established).

Run a scoped design session focused on a new milestone:
1. Ask user for milestone name and scope
2. Define requirements specific to this milestone
3. Design new slices for the milestone
4. Each new slice must include a milestone field: 'milestone: {milestone_name}'
5. Append new slices to GRAPH.json (preserve existing slices)
6. Append any new decisions to DECISIONS.md with source='milestone'
7. Get user approval

Do NOT regenerate the entire design. This is additive only.
", subagent_type="gl-designer", model="{resolved_model.designer}", description="Plan milestone")
```

### After Milestone Planning

When the designer returns, verify outputs:

```bash
cat .greenlight/GRAPH.json 2>/dev/null
grep -i "milestone:" .greenlight/GRAPH.json
```

Verify that new slices have been added to GRAPH.json with milestone field populated.

### Read DESIGN.md for Verification

After the designer completes, read DESIGN.md to verify milestone context was captured:

```bash
cat .greenlight/DESIGN.md
```

### Read CONTRACTS.md for Milestone Contracts

Read CONTRACTS.md to verify new contracts were added for the milestone slices:

```bash
cat .greenlight/CONTRACTS.md
```

### Read STATE.md for Milestone State

Read STATE.md to verify state tracking includes the new milestone:

```bash
cat .greenlight/STATE.md
```

### Commit Milestone Plan

```bash
git add .greenlight/GRAPH.json .greenlight/DECISIONS.md
git commit -m "docs: plan milestone {milestone_name}

- Added {N} slices to GRAPH.json
- Decisions appended to DECISIONS.md
"
```

Report:
```
Milestone planned: {milestone_name}

Added {N} slices to GRAPH.json with milestone field
Decisions appended to DECISIONS.md (source=milestone)

Next: Run /gl:slice {first_slice_id} to start building
```

---

## Mode 3: Archive Completed Milestone

**Usage:** `/gl:roadmap archive`

Archive a completed milestone by identifying finished slices, compressing to summary format, and moving to the Archived Milestones section.

### Prerequisites

```bash
cat .greenlight/ROADMAP.md 2>/dev/null
cat .greenlight/STATE.md 2>/dev/null
cat .greenlight/GRAPH.json 2>/dev/null
```

**Error handling:**

- **If no roadmap found (NoRoadmap):** "No roadmap found. Cannot archive without a roadmap."

Stop if ROADMAP.md doesn't exist.

### Identify Completed Milestones

Read STATE.md and GRAPH.json to identify milestones where all slices are complete:

```bash
cat .greenlight/STATE.md
cat .greenlight/GRAPH.json
```

Logic:
1. Extract unique milestone names from GRAPH.json (from milestone field)
2. For each milestone, check if all slices with that milestone are marked "complete" in STATE.md
3. Build list of fully completed milestones

**If no completed milestones (NoCompletedMilestones):**
```
No completed milestones found. All active milestones have pending slices.
```
Stop here.

### User Selection

If multiple completed milestones exist, ask user to select which one to archive. Present options and let user pick or choose which milestone to archive:

```
Completed milestones available for archiving:

1. {milestone_1} — {N} slices complete
2. {milestone_2} — {M} slices complete

Which milestone do you want to archive? (Enter number or name)
```

### Compress Milestone

For the selected milestone, compress to one-line summary format:

```
[{milestone_name}] — {N} slices, {key feature summary}, completed {date}
```

### Update ROADMAP.md

1. Read current ROADMAP.md
2. Find the milestone table rows for the archived milestone
3. Move them to the "Archived Milestones" section
4. Convert to compressed summary format
5. Write updated ROADMAP.md

```bash
# Read current content
cat .greenlight/ROADMAP.md

# After processing, write updated version with archived section
```

Ensure the "Archived Milestones" section exists. If not, create it at the end of ROADMAP.md.

### Commit Archive

```bash
git add .greenlight/ROADMAP.md
git commit -m "docs: archive milestone {milestone_name}

- Moved {N} slices to Archived Milestones section
- Compressed to summary format
"
```

Report:
```
Milestone archived: {milestone_name}

{N} slices moved to Archived Milestones section
Format: one-line summary with completion date
```

---

## Error Recovery

### Missing Prerequisites

If required files are missing at any stage, report which files are needed and which command to run:

```
Missing prerequisites:
- config.json: Run /gl:init
- DESIGN.md: Run /gl:design
- ROADMAP.md: Run /gl:design
```

### Designer Spawn Failure

If gl-designer fails to spawn or returns an error during milestone planning:
1. Log the error
2. Retry once with same context
3. If retry fails → report to user and stop

### Empty or Invalid GRAPH.json

If GRAPH.json is empty or invalid during milestone planning:
- Log warning
- Proceed with designer — it will create the initial graph structure
