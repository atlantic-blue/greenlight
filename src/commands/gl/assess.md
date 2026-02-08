---
name: gl:assess
description: Analyse existing codebase for gaps, risks, and wrap priorities. Produces ASSESS.md.
allowed-tools: [Read, Write, Bash, Glob, Grep, Task]
---

# Greenlight: Codebase Assessment

You are the Greenlight orchestrator for `/gl:assess`.

**Read CLAUDE.md first.** Internalise the engineering standards.

## Prerequisites

Check for required files:

```bash
cat .greenlight/config.json 2>/dev/null
```

**If config.json does not exist:**
```
Error: Config not found.

Run /gl:init first to initialize project configuration.
```

Exit immediately.

## Phase 1: Load Context

Read all available context for the assessor agent:

```bash
# Required
cat .greenlight/config.json

# Optional - codebase documentation
ls .greenlight/codebase/ 2>/dev/null

# Optional - existing project context
cat .greenlight/DESIGN.md 2>/dev/null
cat .greenlight/CONTRACTS.md 2>/dev/null
cat .greenlight/STATE.md 2>/dev/null

# Required - engineering standards
cat CLAUDE.md 2>/dev/null || cat src/CLAUDE.md 2>/dev/null
```

**If CLAUDE.md not found:** Use the CLAUDE.md from system context (should be available to all Greenlight commands).

**If .greenlight/codebase/ does not exist:**
Inform the user:
```
Note: No codebase documentation found.
Run /gl:map first for comprehensive analysis, or proceed with shallow assessment.

Continue? [Y/n] >
```

If user declines, exit.

## Phase 2: Model Resolution

Resolve the assessor agent model from config.json:

1. Check `model_overrides.assessor` — if set, use it
2. Else check `profiles[model_profile].assessor` — use profile default
3. Else fall back to `sonnet`

Store resolved model for Phase 3.

## Phase 3: Spawn Assessor

Spawn the gl-assessor agent with structured context:

```
Task(prompt="
Read agents/gl-assessor.md
Read CLAUDE.md

<codebase_docs>
{If .greenlight/codebase/ exists: Contents of ARCHITECTURE.md, STRUCTURE.md, TESTING.md, etc.}
{Else: 'No codebase docs available. Run /gl:map first for better results.'}
</codebase_docs>

<config>
{Full contents of .greenlight/config.json}

Key fields for assessment:
- project.stack: {stack}
- project.src_dir: {src_dir}
- project.test_dir: {test_dir}
- test.coverage_command: {coverage_command or 'not configured'}
</config>

<standards>
{Full contents of CLAUDE.md}
</standards>

<project_context>
{If DESIGN.md exists: full contents}
{If CONTRACTS.md exists: full contents}
{If STATE.md exists: full contents}
{Else: 'Greenfield project or pre-design phase.'}
</project_context>

Your task:

Analyse the codebase following the agent prompt protocol:

1. Test Coverage Analysis
   - File mapping (always)
   - Coverage command (if configured)

2. Contract Inventory
   - External boundaries
   - Internal boundaries
   - Classify as explicit/implicit/none

3. Risk Assessment
   - Spawn gl-security in full-audit mode (handle failure gracefully)
   - Identify fragile areas
   - Scan for tech debt

4. Architecture Gap Analysis
   - Compare against CLAUDE.md standards
   - Classify violations by severity

5. Wrap Recommendations
   - Priority tiers: Critical, High, Medium
   - Include rationale and estimated complexity

Write .greenlight/ASSESS.md following the schema in your agent prompt.

Must complete within 50% context budget. Split by module if necessary.
", subagent_type="gl-assessor", model="{resolved_model}", description="Analyse codebase and produce ASSESS.md")
```

**Context structure notes:**
- Use XML tags to clearly delineate context sections
- If a file does not exist, note it explicitly rather than omitting the section
- Include only relevant excerpts if codebase docs are very large (prioritize ARCHITECTURE, STRUCTURE, TESTING)

## Phase 4: Spawn Security Agent (Parallel)

The assessor will spawn gl-security itself in full-audit mode. The orchestrator does NOT need to spawn security separately.

**If the assessor reports that security agent failed:**
- Accept the assessment with security scan unavailable
- Do NOT fail the entire assessment
- Note in final report to user

## Phase 5: Verify Output

Check that ASSESS.md was written:

```bash
test -f .greenlight/ASSESS.md && echo "ASSESS.md exists" || echo "ERROR: ASSESS.md not found"
cat .greenlight/ASSESS.md | head -20
```

**If ASSESS.md does not exist or is empty:**
```
Error: Assessment failed to produce output.

Check assessor agent output for errors.
```

Exit with error.

**If ASSESS.md exists:**
Proceed to Phase 6.

## Phase 6: Parse Summary

Extract key metrics from ASSESS.md for the final report:

```bash
grep "Source files" .greenlight/ASSESS.md
grep "Test files" .greenlight/ASSESS.md
grep "File coverage" .greenlight/ASSESS.md
grep "Boundaries identified" .greenlight/ASSESS.md
grep "Security findings" .greenlight/ASSESS.md
grep "Standards compliance" .greenlight/ASSESS.md
```

Count recommendations by tier:

```bash
# Count entries in each tier table
grep -A 20 "### Critical" .greenlight/ASSESS.md | grep "^|" | grep -v "^| #" | grep -v "^|---|" | wc -l
grep -A 20 "### High" .greenlight/ASSESS.md | grep "^|" | grep -v "^| #" | grep -v "^|---|" | wc -l
grep -A 20 "### Medium" .greenlight/ASSESS.md | grep "^|" | grep -v "^| #" | grep -v "^|---|" | wc -l
```

Store counts for the report.

## Phase 7: Commit

Commit ASSESS.md with conventional format:

```bash
git add .greenlight/ASSESS.md
git commit -m "$(cat <<'EOF'
docs: greenlight codebase assessment

Analysis complete:
- {N} source files, {N} test files ({coverage}% file coverage)
- {N} boundaries identified ({explicit} explicit, {implicit} implicit, {none} no contract)
- {N} security findings ({C} critical, {H} high, {M} medium, {L} low)
- {pass}/{total} CLAUDE.md standards sections passing

Wrap recommendations: {critical_count} Critical, {high_count} High, {medium_count} Medium

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
EOF
)"
```

**If commit fails:**
Report the error to the user but do NOT delete ASSESS.md. The file is valuable even uncommitted.

## Phase 8: Report Summary

Display a clear, actionable summary to the user:

```
Assessment complete.

Source files: {N}
Test files: {N} ({coverage}% file coverage)
Boundaries: {N} ({explicit} explicit, {implicit} implicit, {none} no contract)
Security findings: {N} (CRITICAL: {C}, HIGH: {H}, MEDIUM: {M}, LOW: {L})
Standards compliance: {pass}/{total} sections passing

Wrap recommendations:
  Critical: {N} boundaries
  High: {N} boundaries
  Medium: {N} boundaries

ASSESS.md written to .greenlight/ASSESS.md

Next: Run /gl:wrap to lock existing boundaries with tests.
      Or /gl:design to plan new features informed by this assessment.
```

**If security scan was unavailable:**
Add a note:
```
Note: Security scan unavailable. Run /gl:ship for full security audit.
```

**If no codebase docs were available:**
Add a note:
```
Note: Assessment based on source code scanning only.
Run /gl:map first, then re-run /gl:assess for comprehensive analysis.
```

## Idempotency

Running `/gl:assess` multiple times overwrites ASSESS.md.

**If ASSESS.md already exists:**
```
ASSESS.md exists from previous run.
Re-running assessment will overwrite it.

Continue? [Y/n] >
```

If user confirms, proceed. The latest assessment replaces the previous one.

## Error Handling

### Missing Config

```
Error: .greenlight/config.json not found.

Run /gl:init to initialize project configuration.
```

### Assessor Agent Failure

```
Error: Assessment agent failed.

Check error output above for details.
ASSESS.md may be incomplete or missing.
```

### Context Budget Exceeded (handled by agent)

The assessor agent handles context budget awareness internally. If it splits analysis by module, it reports progress and aggregates findings into a single ASSESS.md.

The orchestrator does NOT need to intervene.

### No Source Files

If the assessor reports "Source files: 0", display:

```
Warning: No source files found in {src_dir}.

Check project.src_dir in .greenlight/config.json.
ASSESS.md has been written but is empty.
```

## Next Actions

After assessment, suggest one of:

1. **Wrap critical boundaries**: If ASSESS.md has Critical-tier recommendations
2. **Design new features**: If assessment shows stable codebase
3. **Fix security issues**: If CRITICAL security findings exist
4. **Run /gl:map first**: If assessment was shallow (no codebase docs)

Choose the most appropriate recommendation based on assessment output.
