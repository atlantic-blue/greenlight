---
name: gl-debug
description: Force a structured diagnostic report for the current or specified slice. Read-only — no files written, no state modified.
argument-hint: "[slice_id]"
allowed-tools: [Read, Bash, Glob, Grep]
---

# Greenlight: Debug (Manual Diagnostic)

Produce a structured diagnostic report for any slice at any time. This is a manual pull cord — use it when you want to understand the current state of a slice without modifying anything.

**This command is read-only.** No files written, no state modified. Does not spawn subagents — direct read and display only.

## Prerequisites

Check that the project has been initialized:

```bash
cat .greenlight/config.json 2>/dev/null
```

If config.json does not exist: print "No config found. Run /gl:init first." and stop.

## Determine Target Slice

1. If `slice_id` argument provided: use that
2. If `.greenlight/STATE.md` has a current slice: use that
3. If neither: print "No active slice found. Specify a slice: /gl-debug {slice_id}" and stop

```bash
cat .greenlight/STATE.md
```

Extract the current slice ID and name from the "Current" section.

## Gather Diagnostic Context

Collect all available context. If any step fails, continue with partial diagnostic — do not abort.

### a. Run Test Suite

```bash
{config.test.command} {config.test.filter_flag} {slice_id} 2>&1
```

If the test command fails to execute, include the error in the report: "Test command failed: {error}". Continue with remaining context.

### b. Read Recent Git Activity

```bash
git log --oneline -10
git diff --stat
```

### c. Read State

```bash
cat .greenlight/STATE.md
```

Extract step, progress, and last activity.

### d. Read Contracts

```bash
cat .greenlight/CONTRACTS.md
```

Find and extract contracts for the target slice. If CONTRACTS.md is missing or no contracts found for the slice, include warning: "No contracts found for slice {id}". Continue with partial diagnostic.

### e. Check for Checkpoint Tag

```bash
git tag -l "greenlight/checkpoint/{slice_id}"
```

Report whether a checkpoint tag exists for this slice.

### f. Read Slice Definition

```bash
cat .greenlight/GRAPH.json
```

Extract slice definition: deliverables, packages, dependencies.

## Produce Diagnostic Report

Format and display the following structured report:

```
DIAGNOSTIC REPORT -- Slice {slice_id}: {slice_name}
Generated: {timestamp}

## Current State
Step: {step from STATE.md}
Last activity: {date}
Checkpoint: {tag name if exists, "none" otherwise}

## Test Results
Total: {N}
Passing: {N}
Failing: {N}

{for each failing test:}
### {test_name}
Expected: {inferred from contract + test name}
Actual: {exact error output}
Contracts: {which contract(s) this test covers}

## Recent Changes
{git log --oneline -10}

## Uncommitted Changes
{git diff --stat}

## Files in Scope
{deliverables and packages from GRAPH.json for this slice}

## Recovery Options
1) Resume implementation (/gl:slice {slice_id})
2) Roll back to checkpoint: git checkout greenlight/checkpoint/{slice_id} -- .
3) Pause for manual investigation (/gl:pause)
4) Spawn fresh implementer with guidance

## Specific Question
{auto-generated: the most likely root cause based on failing tests and recent changes}
```

Display the report to the user. No files are written.

## Invariants

- /gl-debug is strictly read-only (no files written, no state modified)
- /gl-debug can be run at ANY time, not just during circuit break
- /gl-debug does not control the pipeline (it diagnoses, it does not resume or retry)
- Diagnostic report follows the same structured format as circuit break diagnostics (C-54)
- /gl-debug works without a checkpoint tag (tag presence is informational only)
- /gl-debug runs the test suite to get current state (always fresh data)
- Report is structured for future pause/resume integration
- Does not spawn any subagents (direct read + display)
