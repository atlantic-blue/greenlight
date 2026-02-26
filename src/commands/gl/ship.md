---
name: gl:ship
description: Final verification — all tests green, security audit, deployment readiness
allowed-tools: [Read, Bash, Glob, Grep, Task, AskUserQuestion]
---

# Greenlight: Ship

All slices should be green. Final verification before deploy.

**Read first:**
- `.greenlight/STATE.md` — verify all slices complete
- `.greenlight/config.json` — test commands, project config

## Model Resolution

Before spawning any agent, resolve its model from `.greenlight/config.json`:

1. Check `model_overrides[agent_name]` — if set, use it
2. Else check `profiles[model_profile][agent_name]` — use profile default
3. Else fall back to `sonnet`

Agents spawned by this command: `security`, `verifier`.

## Pre-Check

Verify all slices are complete:

```bash
cat .greenlight/STATE.md
cat .greenlight/GRAPH.json
```

If any slice is not "complete" → stop:
```
Cannot ship: {N} slices incomplete
- Slice {id}: {name} — status: {status}

Complete all slices first, then run /gl:ship
```

## Step 1: Clean Build

Start fresh to catch any "works on my machine" issues:

```bash
# Remove build artifacts and cached dependencies
rm -rf node_modules dist build .cache  # adapt to stack
# Clean install from lock file
npm ci  # or pip install -r requirements.txt, go mod download, etc.
```

```bash
# Lint
npm run lint 2>&1
```

```bash
# Full test suite
npm test 2>&1
```

```bash
# Production build
npm run build 2>&1
```

If any step fails → stop, report the failure, suggest fix. Don't skip steps.

## Step 2: Full Security Audit

Spawn security agent in full-audit mode against entire codebase:

```
Task(prompt="
Read agents/gl-security.md
Read CLAUDE.md

<mode>full-audit</mode>

<codebase_summary>
{list all source files in src/ with line counts}
</codebase_summary>

<contracts>
{all contracts from .greenlight/CONTRACTS.md}
</contracts>

Perform full security audit:
1. Review all source files against the full vulnerability checklist
2. Check cross-cutting concerns: auth consistency, data flow between slices
3. Run dependency audit

Also check:
- npm audit / pip audit (dependency vulnerabilities)
- Hardcoded secrets patterns in all files
- .env.example exists and matches required vars
- .gitignore includes sensitive files
- No secrets in git history

Write any new security tests to tests/security/audit.test.{ext}
", subagent_type="gl-security", model="{resolved_model.security}", description="Full security audit")
```

If new security tests generated → spawn implementer to fix → re-run ALL tests.

## Step 3: Full Verification

Spawn verifier in ship mode:

```
Task(prompt="
Read agents/gl-verifier.md
Read references/verification-patterns.md

<contracts>
{all contracts}
</contracts>

<test_results>
{output from full test suite run}
</test_results>

<files_changed>
{all source files}
</files_changed>

<mode>ship</mode>

Full project verification:
1. Every contract has test coverage
2. No stubs in production code
3. All modules wired and reachable
4. Cross-slice integration verified
", subagent_type="gl-verifier", model="{resolved_model.verifier}", description="Full project verification")
```

## Step 4: Automated Checklist

```bash
# No console.log in production
grep -rn "console\.log" src/ --include="*.ts" --include="*.js" | grep -v "logger" | head -20 || echo "CLEAN"

# No 'any' types (TypeScript)
grep -rn ": any\| as any" src/ --include="*.ts" | head -20 || echo "CLEAN"

# No TODO/FIXME/HACK
grep -rn "TODO\|FIXME\|HACK\|XXX" src/ | head -20 || echo "CLEAN"

# Secrets scan
grep -rn "password\s*=\|secret\s*=\|api_key\s*=\|apikey\s*=" src/ --include="*.ts" --include="*.js" --include="*.py" -i | head -10 || echo "CLEAN"

# .env.example exists
[ -f .env.example ] && echo "OK: .env.example" || echo "MISSING: .env.example"

# README exists
[ -f README.md ] && echo "OK: README.md" || echo "MISSING: README.md"

# .gitignore covers essentials
for pattern in node_modules .env dist coverage; do
  grep -q "$pattern" .gitignore 2>/dev/null && echo "OK: .gitignore has $pattern" || echo "MISSING: .gitignore needs $pattern"
done

# Lock file committed
ls package-lock.json yarn.lock pnpm-lock.yaml 2>/dev/null && echo "OK: lock file" || echo "MISSING: lock file"
```

### CI Portability Check

Tests must pass on CI, not just locally. Scan for patterns that break on other machines:

```bash
# Hardcoded absolute paths in test files (breaks CI)
grep -rn "/Users/\|/home/\|C:\\\\Users\\\\" . --include="*_test.go" --include="*.test.*" --include="*.spec.*" | head -20 || echo "CLEAN"

# Hardcoded localhost ports without env override
grep -rn "localhost:[0-9]\|127\.0\.0\.1:[0-9]" . --include="*_test.go" --include="*.test.*" --include="*.spec.*" | grep -v "// ci-ok" | head -10 || echo "CLEAN"

# OS-specific path separators in test assertions
grep -rn '\\\\' . --include="*_test.go" --include="*.test.*" | grep -v "regexp\|regex\|escape" | head -10 || echo "CLEAN"
```

**If hardcoded paths found:** NO-GO. Tests will pass locally but fail in CI. Fix by resolving paths relative to the test file location (e.g. `runtime.Caller(0)` in Go, `__dirname` in Node, `Path(__file__)` in Python).

**If CI workflow exists**, verify tests actually run there:
```bash
# Check CI config exists and runs tests
ls .github/workflows/*.yml .gitlab-ci.yml Jenkinsfile 2>/dev/null
```

## Step 5: Report

```
┌─────────────────────────────────────────────────┐
│  GREENLIGHT: SHIP CHECK                         │
├─────────────────────────────────────────────────┤
│                                                 │
│  Build:        {result}                    {ok} │
│  Lint:         {result}                    {ok} │
│  Tests:        {N} passing                 {ok} │
│  CI Portable:  {result}                    {ok} │
│  Security:     {N} tests, {N} audit issues {ok} │
│  Verification: {result}                    {ok} │
│  Dependencies: {N} vulnerabilities         {ok} │
│  Secrets:      {result}                    {ok} │
│  Contracts:    {N}/{N} covered             {ok} │
│  Env vars:     {result}                    {ok} │
│  Docs:         {result}                    {ok} │
│                                                 │
│  VERDICT: GREENLIGHT / NO-GO                    │
└─────────────────────────────────────────────────┘
```

### GREENLIGHT (all checks pass)

```
GREENLIGHT — ready to ship

  git tag -a v{version} -m "{project name} MVP"
  git push origin main --tags

Slices: {N}/{N} complete
Tests: {N} passing ({functional} functional, {security} security)
```

### NO-GO (any check fails)

```
NO-GO — {N} issues blocking shipment

{list each issue with category, description, and suggested fix}

Fix issues and re-run /gl:ship
```

Prioritise fixes:
1. CRITICAL security issues first
2. Failing tests
3. Missing contract coverage
4. Everything else

---

## File-Per-Slice State Integration (C-84)

This section documents how /gl:ship adapts its state reads when the project uses file-per-slice format (C-80). Legacy format behaviour is completely unchanged — if the state format is legacy, all reads go to STATE.md as before.

### State Format Detection

Detect state format (C-80) before performing any state reads:

- If file-per-slice: read all slice files from `.greenlight/slices/` for the completeness pre-check
- If legacy: read STATE.md as before (no change)

### File-Per-Slice Read Path

When using file-per-slice format:

1. Read all slice files from `.greenlight/slices/`
2. Pre-check: all slice files must have status complete — if any slice file does not have status complete, stop and report which slices are incomplete
3. If any incomplete slices found: stop and report which slices are not complete before proceeding

### Legacy Fallback

If format is legacy: read STATE.md as before. Legacy format behaviour is completely unchanged — no change.

### Error Handling

| Error State | When | Behaviour |
|-------------|------|-----------|
| FormatDetectionFailure | Cannot determine state format | Report error. Suggest running /gl:init |
| SliceFileNotFound | A referenced slice file does not exist in `.greenlight/slices/` | Report error. Cannot verify completeness without slice file |
| RegenerationFailure | STATE.md regeneration fails | Warn but continue. Slice files are still correct |

### Invariants

- All slice files must be read before the ship check proceeds
- Legacy format behaviour is completely unchanged
- Slice files in `.greenlight/slices/` are the source of truth for completeness verification
