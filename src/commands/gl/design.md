---
name: gl:design
description: Run an interactive system design session. Produces DESIGN.md with requirements, architecture, and technical decisions.
allowed-tools: [Read, Write, Bash, Glob, Grep, Task, AskUserQuestion, WebSearch, WebFetch]
---

# Greenlight: System Design

You are the Greenlight orchestrator for the design phase. Spawn the designer agent to run an interactive system design session that bridges the init interview and contract generation.

**Read CLAUDE.md first.** Internalise the engineering standards.

## Prerequisites

Check that `/gl:init` has been run:

```bash
cat .greenlight/config.json 2>/dev/null
```

**If no config exists:**
```
No project found. Run /gl:init first to set up the project.
```
Stop here.

**If config exists:** Read the config and any existing interview context.

## Check for Existing Design

```bash
cat .greenlight/DESIGN.md 2>/dev/null
```

**If DESIGN.md exists:** Ask the user: "A design already exists. Do you want to revise it or start fresh?"
- **Revise:** Pass existing DESIGN.md to the designer as context
- **Start fresh:** Proceed without it

## Model Resolution

Before spawning the designer, resolve its model from `.greenlight/config.json`:

1. Check `model_overrides["designer"]` — if set, use it
2. Else check `profiles[model_profile]["designer"]` — use profile default
3. Else fall back to `sonnet`

## Gather Context

Read available project context:

```bash
cat .greenlight/config.json
cat .greenlight/INTERVIEW.md 2>/dev/null
```

If brownfield (existing codebase was mapped):
```bash
ls .greenlight/codebase/ 2>/dev/null
```

Build the context block from what's available.

## Spawn Designer

```
Task(prompt="
Read agents/gl-designer.md
Read CLAUDE.md

<project_context>
{project name, value prop, users, MVP scope from config/interview}
</project_context>

<stack>
{chosen stack from config}
</stack>

<existing_code>
{if brownfield: summary from .greenlight/codebase/ docs. Otherwise: 'Greenfield project'}
</existing_code>

<existing_design>
{if revising: contents of current DESIGN.md. Otherwise: 'No existing design'}
</existing_design>

Run the full design session:
1. Deep dive on requirements (functional, non-functional, constraints, out of scope)
2. Research genuine unknowns (use WebSearch when needed)
3. Propose solution (architecture, data model, API surface, security, deployment)
4. Discuss gray areas with the user
5. Write .greenlight/DESIGN.md
6. Get user approval

Do NOT produce contracts or dependency graphs. That's the architect's job.
", subagent_type="gl-designer", model="{resolved_model.designer}", description="System design session")
```

## After Design

When the designer returns with an approved DESIGN.md, report:

```
Design session complete.

{summary from designer}

Next step: Run /gl:init to generate contracts from this design.
The architect will read DESIGN.md and produce typed contracts
and a dependency graph.
```
