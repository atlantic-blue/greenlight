---
name: gl:map
description: Analyse existing codebase with parallel agents before initializing
argument-hint: "[optional: specific area, e.g., 'api' or 'auth']"
allowed-tools: [Read, Write, Bash, Glob, Grep, Task]
---

# Greenlight: Map Codebase

Analyse existing codebase using parallel mapper agents. Run this before `/gl:init` on brownfield projects.

## Pre-check

```bash
# Verify there's actually code to map
find . -name "*.ts" -o -name "*.js" -o -name "*.py" -o -name "*.swift" -o -name "*.go" -o -name "*.rs" 2>/dev/null | grep -v node_modules | grep -v .git | head -20
```

If no source files found → nothing to map. Suggest `/gl:init` for greenfield.

## Create Output Directory

```bash
mkdir -p .greenlight/codebase
```

## Model Resolution

Before spawning any agent, resolve its model from `.greenlight/config.json`:

1. Check `model_overrides["codebase_mapper"]` — if set, use it
2. Else check `profiles[model_profile]["codebase_mapper"]` — use profile default
3. Else fall back to `sonnet`

If no config exists yet (mapping before init), use `sonnet`.

## Spawn 4 Parallel Mappers

Each mapper writes directly to `.greenlight/codebase/`. Orchestrator receives confirmation only — NOT document contents (saves context).

```
# Agent 1: tech focus → STACK.md, INTEGRATIONS.md
Task(prompt="
Read agents/gl-codebase-mapper.md
Focus: tech
Write to .greenlight/codebase/STACK.md and .greenlight/codebase/INTEGRATIONS.md
Project root: {pwd}
", subagent_type="gl-codebase-mapper", model="{resolved_model.codebase_mapper}", description="Map: tech stack")

# Agent 2: arch focus → ARCHITECTURE.md, STRUCTURE.md
Task(prompt="
Read agents/gl-codebase-mapper.md
Focus: arch
Write to .greenlight/codebase/ARCHITECTURE.md and .greenlight/codebase/STRUCTURE.md
Project root: {pwd}
", subagent_type="gl-codebase-mapper", model="{resolved_model.codebase_mapper}", description="Map: architecture")

# Agent 3: quality focus → CONVENTIONS.md, TESTING.md
Task(prompt="
Read agents/gl-codebase-mapper.md
Focus: quality
Write to .greenlight/codebase/CONVENTIONS.md and .greenlight/codebase/TESTING.md
Project root: {pwd}
", subagent_type="gl-codebase-mapper", model="{resolved_model.codebase_mapper}", description="Map: quality")

# Agent 4: concerns focus → CONCERNS.md
Task(prompt="
Read agents/gl-codebase-mapper.md
Focus: concerns
Write to .greenlight/codebase/CONCERNS.md
Project root: {pwd}
", subagent_type="gl-codebase-mapper", model="{resolved_model.codebase_mapper}", description="Map: concerns")
```

**Launch all 4 in parallel** — they write to different files and don't depend on each other.

## Verify

```bash
echo "=== Codebase Map ==="
for f in .greenlight/codebase/*.md; do
  echo "$(basename $f): $(wc -l < "$f") lines"
done
echo "===================="
```

All 7 documents should exist with content. If any are empty or missing, report which mapper failed.

## Handle Critical Findings

Check CONCERNS.md for critical findings:

```bash
grep -i "CRITICAL" .greenlight/codebase/CONCERNS.md 2>/dev/null
```

If critical findings exist (hardcoded secrets, SQL injection, etc.):
```
CRITICAL findings in codebase:
{list from CONCERNS.md}

Address these before proceeding with /gl:init.
```

## Commit

```bash
git add .greenlight/codebase/
git commit -m "docs: greenlight codebase map

Documents: $(ls .greenlight/codebase/*.md | wc -l | tr -d ' ')
"
```

## Next

```
Codebase mapped

Documents: 7
Location: .greenlight/codebase/

{if critical findings: "Address critical findings first."}
{otherwise: "Run /gl:init to set up contracts and slices."}
The architect will use this map to understand your existing code.
```
