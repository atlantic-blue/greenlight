---
name: gl:wrap
description: Wrap an existing boundary with contracts and locking tests
allowed-tools: [Read, Write, Bash, Glob, Grep, Task]
---

# Greenlight: Wrap Boundary

Extracts contracts from existing code and writes locking tests to lock in current behaviour. This prepares boundaries for safe refactoring.

**Read first:**
- .greenlight/config.json (REQUIRED)
- .greenlight/ASSESS.md (optional, used for recommendations)
- .greenlight/CONTRACTS.md
- .greenlight/STATE.md
- CLAUDE.md

## Model Resolution

Wrap uses TWO models:
1. **Wrapper agent:** Analyses code, extracts contracts, writes tests (complex work)
2. **Security agent:** Documents security baseline (lightweight review)

```python
def resolve_wrapper_model(boundary_complexity: str, config: dict) -> str:
    """
    Complexity determined by:
    - File count
    - Total LOC
    - External dependency count
    - Code complexity (nesting, conditionals)
    """
    if boundary_complexity == "high":
        return config.get("defaultModel") or "claude-opus-4-6"
    elif boundary_complexity == "medium":
        return config.get("defaultModel") or "claude-sonnet-4-5"
    else:  # low
        return "claude-sonnet-4-5"

def resolve_security_model(config: dict) -> str:
    """Security in wrap is document-only, use fast model"""
    return "claude-sonnet-4-5"
```

Model selection shown to user before spawning agents.

## Pre-flight Checks

### 1. Validate Environment
```bash
# Config must exist
if not exists(".greenlight/config.json"):
    error("No config.json found. Run /gl:init first.")

# CONTRACTS.md must exist (even if empty)
if not exists(".greenlight/CONTRACTS.md"):
    create(".greenlight/CONTRACTS.md", template="# Contracts\n\n")

# STATE.md must exist
if not exists(".greenlight/STATE.md"):
    error("No STATE.md found. Run /gl:init first.")
```

### 2. Read Assessment (Optional)
```bash
if exists(".greenlight/ASSESS.md"):
    assessment = read(".greenlight/ASSESS.md")
    parse_boundaries(assessment)
else:
    assessment = None
    # Will scan codebase for boundaries
```

### 3. Read Existing Contracts
```bash
contracts = read(".greenlight/CONTRACTS.md")
existing_wrapped = extract_wrapped_boundaries(contracts)
```

## Interactive Boundary Selection

### Display Candidates

**If ASSESS.md exists:**
Show boundaries from assessment, grouped by priority:

```
Reading assessment... (.greenlight/ASSESS.md)

Recommended boundaries to wrap:

CRITICAL (untested, external, high complexity)
  1. auth — 3 files, 12 functions, no tests
  2. payments — 2 files, 8 functions, external API

HIGH (untested, external)
  3. users — 4 files, 15 functions, partial tests

MEDIUM (internal, untested)
  4. notifications — 2 files, 6 functions, internal queue

Already wrapped: {list of wrapped boundaries}

Which boundary to wrap? [1-4 / list / scan / cancel]
```

**If ASSESS.md missing:**
```
No assessment found.

Options:
  1. scan — Scan codebase for boundaries (recommended)
  2. manual — Enter boundary files manually
  3. cancel

Choice? [1-3]
```

If user picks "scan":
```python
# Use Glob to find potential boundaries
boundaries = scan_for_boundaries(
    patterns=["src/**/*.ts", "lib/**/*.py", "app/**/*.rb"],
    indicators=["export", "public", "API", "Controller", "Service"]
)

# Present top candidates
present_boundary_candidates(boundaries, limit=10)
```

If user picks "manual":
```
Enter boundary files (one per line, empty line to finish):
> src/auth/login.ts
> src/auth/middleware.ts
>

Boundary name? > auth
```

### User Selection
```
Which boundary to wrap? > 1
```

Store selection: `boundary_name`, `boundary_files[]`

## Step 1: Analyse Boundary

Before spawning the wrapper agent, do quick analysis:

```python
files = read_files(boundary_files)
total_loc = sum(count_lines(f) for f in files)
file_count = len(files)

complexity = estimate_complexity(files)  # low | medium | high

print(f"""
Analysing {boundary_name} boundary...
  Files: {file_count}
  Total LOC: {total_loc}
  Complexity: {complexity}
""")

# Check if too large for single agent session
if total_loc > 3000 or file_count > 15:
    print(f"""
⚠️  Large boundary detected.

This boundary may exceed context budget for single wrap session.

Options:
  1. proceed — Attempt wrap (agent will split if needed)
  2. split — Manually split into smaller boundaries
  3. cancel

Choice? [1-3]
""")

    if choice == "split":
        print("Suggest logical splits based on:")
        suggest_splits(files)
        exit()
```

Resolve model:
```python
wrapper_model = resolve_wrapper_model(complexity, config)
security_model = resolve_security_model(config)

print(f"""
Model selection:
  Wrapper: {wrapper_model}
  Security: {security_model}

Scope fits within context budget. Proceeding.
""")
```

## Step 2: Spawn Wrapper Agent

```python
Task(
    prompt=f"""
<boundary>
  <name>{boundary_name}</name>
  <files>
    {'\n    '.join(f'<file>{f}</file>' for f in boundary_files)}
  </files>
  <config>{config_contents}</config>
  <existing_contracts>{contracts_contents}</existing_contracts>
  <codebase_docs>
    {readme_contents}
    {architecture_contents if exists else ''}
  </codebase_docs>
  <claude_standards>{claude_md_contents}</claude_standards>
</boundary>

Extract contracts from this boundary and write locking tests.

Follow your agent prompt (gl-wrapper.md) precisely:
1. Analyse boundary and report findings
2. Extract contracts (descriptive, not prescriptive)
3. Wait for confirmation
4. Write locking tests to tests/locking/{boundary_name}.test.{ext}
5. Run tests, fix if needed (max 3 cycles)
6. Run full test suite
7. Report completion

CRITICAL: You are in READ-ONLY mode for production source code.
NEVER modify source code to make tests pass.
""",
    subagent_type="gl-wrapper",
    model=wrapper_model,
    description=f"Wrapping {boundary_name} boundary"
)
```

### Agent Interaction Protocol

The wrapper agent will pause at key points for orchestrator/user input:

**After analysis:**
```
Agent: "Ready to extract contracts."
Orchestrator: Confirms or adjusts scope
```

**After contract extraction:**
```
Agent: "Accept these contracts? Awaiting confirmation."
Orchestrator: Shows contracts to user, waits for y/N/edit
  - y: proceed
  - N: abort
  - edit: user provides feedback, agent revises
```

**After test writing:**
```
Agent: "Running tests..."
Orchestrator: Monitors, no action needed unless failure
```

**After completion:**
```
Agent: "Wrap complete: {boundary_name}"
Orchestrator: Proceeds to security baseline
```

## Step 3: Security Baseline (Document Only)

After wrapper completes successfully:

```python
print("\nSecurity baseline...\n")

Task(
    prompt=f"""
You are reviewing the {boundary_name} boundary as part of a WRAP operation.

**Mode:** DOCUMENT ONLY (no failing tests)

Scope:
{'\n'.join(f'  - {f}' for f in boundary_files)}

Review for:
- Authentication/authorization issues
- Input validation gaps
- Secrets in code
- SQL injection vectors
- XSS vulnerabilities
- Rate limiting absence
- Logging of sensitive data

Document findings in ASSESS.md under "Security Baseline: {boundary_name}".
Include severity (HIGH/MEDIUM/LOW) and location.

Do NOT write failing tests. Do NOT block the wrap.

If you encounter errors, document what you could review and note limitations.
""",
    subagent_type="gl-security",
    model=security_model,
    mode="slice",  # Document only
    description=f"Security baseline for {boundary_name}"
)
```

**If security agent fails:**
```python
print("""
  ⚠️  Security baseline could not complete.
  Proceeding without security documentation.
  Recommend manual security review.
""")
# Continue — security failure doesn't block wrap
```

**If security agent succeeds:**
```python
findings = parse_security_findings()
print(f"  {len(findings)} issues documented")
for finding in findings[:3]:  # Show first 3
    print(f"    - {finding.severity}: {finding.summary}")
```

## Step 4: Commit Atomically

Wrap changes are committed atomically:

```python
# Files that should be staged:
staged_files = [
    ".greenlight/CONTRACTS.md",  # New contracts appended
    f"tests/locking/{boundary_name}.test.{ext}",  # New locking tests
    ".greenlight/STATE.md",  # Updated wrap progress
    ".greenlight/ASSESS.md"  # Updated with security findings (if successful)
]

# Verify all expected files exist
for file in staged_files:
    if not exists(file):
        error(f"Expected file missing: {file}. Wrap incomplete.")

# Stage specific files (never use git add -A)
git_add_files = []
for file in staged_files:
    if git_status_shows_changes(file):
        git_add_files.append(file)

if not git_add_files:
    print("No changes to commit (wrap produced no modifications).")
    exit()

# Commit
git add {' '.join(git_add_files)}

git commit -m "$(cat <<'EOF'
test(wrap): lock {boundary_name}

Extracted {N} contract(s) and wrote {M} locking tests.
All tests passing.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
EOF
)"

print(f"\nCommitted: test(wrap): lock {boundary_name}\n")
```

**If commit fails:**
```python
print("""
❌ Commit failed.

Wrap completed but changes not committed.
Manually review and commit:
  git status
  git diff
  git add {files}
  git commit -m "test(wrap): lock {boundary_name}"
""")
# Don't exit — user can recover
```

## Step 5: Update STATE.md

After successful commit:

```python
# Read current STATE.md
state = read(".greenlight/STATE.md")

# Check if "Wrapped Boundaries" section exists
if "## Wrapped Boundaries" not in state:
    # Add section after Slices section
    wrapped_section = """
## Wrapped Boundaries

| Boundary | Contracts | Locking Tests | Known Issues | Status |
|----------|-----------|---------------|--------------|--------|

Wrap progress: 0/{total_boundaries} boundaries wrapped
"""
    state = insert_after(state, "## Slices", wrapped_section)

# Add this boundary to table
boundary_row = f"| {boundary_name} | {contract_count} | {test_count} | {issue_count} | wrapped |"
state = append_to_table(state, "Wrapped Boundaries", boundary_row)

# Update progress counter
wrapped_count = count_wrapped_boundaries(state)
total_boundaries = estimate_total_boundaries()  # From ASSESS.md or scan
state = update_wrap_progress(state, wrapped_count, total_boundaries)

write(".greenlight/STATE.md", state)
```

## Step 6: Display Summary

```python
print(f"""
╔══════════════════════════════════════════════════════════════
║ Wrap Complete: {boundary_name}
╚══════════════════════════════════════════════════════════════

Contracts: {contract_count} written to CONTRACTS.md
Locking tests: tests/locking/{boundary_name}.test.{ext} ({test_count} tests)
Security: {issue_count} known issues documented

All {test_count} locking tests passing.
Full suite: {total_test_count} tests passing (no regressions).

Wrap progress: {wrapped_count}/{total_boundaries} boundaries wrapped

Next steps:
  /gl:wrap — Wrap next boundary (recommended: {next_boundary})
  /gl:design — Design features using wrapped boundaries
  /gl:slice 1 — Start building new functionality

This boundary is now available for safe refactoring via /gl:slice
using the 'wraps' field.
""")
```

## Error Handling

| Error | When | Response |
|-------|------|----------|
| NoConfig | config.json missing | "Run /gl:init first" |
| NoStateFile | STATE.md missing | "Run /gl:init first" |
| BoundaryTooLarge | Agent reports excessive size | Show agent's split suggestions, offer to retry with subset |
| ExistingContracts | Boundary already wrapped | "Boundary '{name}' already has contracts. Overwrite? [y/N]" |
| LockingTestFailure | Tests won't pass after 3 cycles | "Cannot lock {boundary}. Manual review needed. See tests/locking/{boundary}.test.{ext}" |
| TestRegressions | Full suite fails after wrap | "REGRESSION: Locking tests caused failures in {failing_tests}. Recommend rollback." |
| SecurityFailure | Security agent fails | Warn, continue without security docs |
| CommitFailure | Git commit fails | Show manual recovery steps, don't exit |
| AgentContextExceeded | Agent hits context limit mid-wrap | "Partial wrap completed. Resume with remaining files? [y/N]" |

## Working Without ASSESS.md

If no assessment exists, wrap still works:

1. **Scan mode:** Use Glob/Grep to find boundary candidates
   - Look for exports, public methods, API routes, controllers
   - Group by directory/module
   - Rank by lack of tests (grep for test file existence)

2. **Manual mode:** User provides file list and boundary name

3. **No recommendations:** Can't prioritize by risk, but wrapping still valuable

## Context Budget Management

Orchestrator is thin (stays under 30%). Heavy work in wrapper agent.

If wrapper reports BoundaryTooLarge:
```python
print(f"""
⚠️  Boundary too large for single session.

Agent suggests splitting into:
{agent_split_suggestions}

Options:
  1. Split and wrap separately
  2. Wrap subset now (specify files)
  3. Cancel

Choice? [1-3]
""")
```

## State Tracking

After wrap completes, STATE.md includes:

```markdown
## Wrapped Boundaries

| Boundary | Contracts | Locking Tests | Known Issues | Status |
|----------|-----------|---------------|--------------|--------|
| auth | 4 | 12 | 2 | wrapped |
| payments | 2 | 8 | 0 | wrapped |

Wrap progress: 2/6 boundaries wrapped
```

Status values:
- **wrapped:** Initial wrap complete, available for refactoring
- **refactored:** Later refactored via /gl:slice (wraps field used)

When a slice refactors a wrapped boundary, status updates to "refactored".

## Recommendations After Wrap

```python
def recommend_next(state: dict, assessment: dict) -> str:
    """Recommend what to do after successful wrap"""

    unwrapped = get_unwrapped_boundaries(state, assessment)

    if unwrapped:
        # Prioritize by risk from assessment
        next_boundary = unwrapped[0]
        return f"/gl:wrap — wrap next boundary (recommended: {next_boundary})"
    else:
        # All boundaries wrapped
        return "/gl:design — design features using wrapped contracts"
```

## Full Flow Example

```
$ /gl:wrap

Reading config... (.greenlight/config.json)
Reading assessment... (.greenlight/ASSESS.md)

Recommended boundaries to wrap:

CRITICAL
  1. auth — 3 files, 12 functions, no tests, external dependency
  2. payments — 2 files, 8 functions, external API

HIGH
  3. users — 4 files, 15 functions, partial tests

Already wrapped: []

Which boundary to wrap? [1-3 / list / scan / cancel] > 1

Analysing auth boundary...
  Files: 3
  Total LOC: 450
  Complexity: medium

Model selection:
  Wrapper: claude-sonnet-4-5
  Security: claude-sonnet-4-5

Scope fits within context budget. Proceeding.

[Spawning gl-wrapper agent...]

═══════════════════════════════════════════════════════════════
Boundary Analysis: auth

Files analysed: 3
- src/auth/login.ts (180 lines)
- src/auth/middleware.ts (150 lines)
- src/auth/tokens.ts (120 lines)

Identified contracts: 2
1. AuthenticateUser — Login with email/password, returns JWT
2. AuthorizeRequest — Middleware verifies JWT, attaches user

Entry points: 2
- authenticateUser(email, password)
- authorizeRequest(req, res, next)

External dependencies:
- Database (users table)
- JWT library

Complexity estimate: medium
Context usage: ~22%

Non-determinism detected:
- JWT tokens include timestamp and random jti claim
- createdAt timestamps

Existing contracts check: none found

Ready to extract contracts.
═══════════════════════════════════════════════════════════════

Proceeding with extraction...

═══════════════════════════════════════════════════════════════
Extracted Contracts

### Contract: AuthenticateUser [WRAPPED]

**Source:** `src/auth/login.ts:15-45`
**Wrapped on:** 2026-02-08
**Locking tests:** `tests/locking/auth.test.ts`

**Boundary:** HTTP Controller → Authentication Service → Database
**Slice:** wrappable (available for refactoring via /gl:slice with wraps field)

**Input:**
```typescript
{
  email: string
  password: string
}
```

**Output:**
```typescript
{
  token: string
  expiresIn: number
  user: {
    id: string
    email: string
    role: string
  }
}
```

**Errors:**
| Error | Status | When |
|-------|--------|------|
| "Invalid credentials" | 401 | Email not found OR password mismatch (no distinction) |
| "Account locked" | 403 | failedLoginAttempts >= 5 |
| "Database error" | 500 | Database connection fails |

**Invariants:**
- Email always converted to lowercase before lookup
- Password never logged or returned
- JWT expires in 3600 seconds (hardcoded)
- Token includes user.id, user.email, user.role claims

**Security:**
- Known issues:
  - HIGH: No rate limiting (vulnerable to brute force)
  - MEDIUM: Account lockout is permanent (no unlock mechanism)
  - LOW: Password appears in debug logs on validation error

**Dependencies:** none

---

### Contract: AuthorizeRequest [WRAPPED]

**Source:** `src/auth/middleware.ts:10-35`
**Wrapped on:** 2026-02-08
**Locking tests:** `tests/locking/auth.test.ts`

**Boundary:** Express Middleware → JWT Verification → Request Context
**Slice:** wrappable (available for refactoring via /gl:slice with wraps field)

**Input:**
```typescript
{
  headers: {
    authorization?: string  // Format: "Bearer {token}"
  }
}
```

**Output:**
```typescript
// Mutates request object:
req.user = {
  id: string
  email: string
  role: string
}
// Then calls next()
```

**Errors:**
| Error | Status | When |
|-------|--------|------|
| "No token provided" | 401 | Authorization header missing |
| "Invalid token format" | 401 | Authorization header doesn't start with "Bearer " |
| "Token expired" | 401 | JWT exp claim < current time |
| "Invalid token" | 401 | JWT signature verification fails |

**Invariants:**
- Always calls next() on success
- Always calls next(error) on failure
- Never throws synchronously
- User object structure matches JWT claims exactly

**Security:**
- Known issues:
  - LOW: No token revocation (logout doesn't invalidate tokens)

**Dependencies:** AuthenticateUser (consumes tokens it generates)

═══════════════════════════════════════════════════════════════

Accept these contracts? [y/N/edit] > y

Writing locking tests... tests/locking/auth.test.ts

═══════════════════════════════════════════════════════════════
Locking Tests Written

File: tests/locking/auth.test.ts
Test count: 12

Coverage:
- ✅ Happy path: Valid credentials return token
- ✅ Happy path: Valid token authorizes request
- ✅ Error case: Invalid password returns 401
- ✅ Error case: Unknown email returns 401 (same error)
- ✅ Error case: Locked account returns 403
- ✅ Error case: Missing authorization header returns 401
- ✅ Error case: Malformed token returns 401
- ✅ Error case: Expired token returns 401
- ✅ Edge case: Email case insensitivity
- ✅ Edge case: Token format (Bearer prefix required)
- ✅ Non-determinism: JWT token format verified, not exact value
- ✅ Non-determinism: Timestamps within valid range

Running tests...
═══════════════════════════════════════════════════════════════

Locking Tests: PASSED

All 12 tests passing.

Next: Run full test suite to check for regressions.

═══════════════════════════════════════════════════════════════
Full Suite Check

Running: npm test

  [LOCK] Authentication Boundary
    successful authentication
      ✓ [LOCK] should return token for valid credentials
      ✓ [LOCK] should return user object without password
    failed authentication
      ✓ [LOCK] should return 401 for invalid password
      ✓ [LOCK] should return 401 for unknown email
      ✓ [LOCK] should return 403 for locked account
    authorization
      ✓ [LOCK] should attach user to request with valid token
      ✓ [LOCK] should return 401 for missing token
      ✓ [LOCK] should return 401 for malformed token
      ✓ [LOCK] should return 401 for expired token
    edge cases
      ✓ [LOCK] should handle email case insensitivity
      ✓ [LOCK] should require Bearer prefix in authorization
      ✓ [LOCK] should generate JWT with expected claims

  Users API
    ✓ should create user
    ✓ should list users

  14 passing (1.2s)

No regressions detected.
═══════════════════════════════════════════════════════════════

Wrap complete: auth

Contracts written: 2
Locking tests: tests/locking/auth.test.ts (12 tests)
Status: ✅ All tests passing

Contracts appended to: .greenlight/CONTRACTS.md
Ready for: Atomic commit

Known issues: 3 (1 HIGH, 1 MEDIUM, 1 LOW)

Security baseline...

  3 issues documented (1 HIGH, 1 MEDIUM, 1 LOW)
    - HIGH: No rate limiting (vulnerable to brute force)
    - MEDIUM: Account lockout is permanent
    - LOW: Password appears in debug logs

Committed: test(wrap): lock auth

╔══════════════════════════════════════════════════════════════
║ Wrap Complete: auth
╚══════════════════════════════════════════════════════════════

Contracts: 2 written to CONTRACTS.md
Locking tests: tests/locking/auth.test.ts (12 tests)
Security: 3 known issues documented

All 12 locking tests passing.
Full suite: 14 tests passing (no regressions).

Wrap progress: 1/6 boundaries wrapped

Next steps:
  /gl:wrap — Wrap next boundary (recommended: payments)
  /gl:design — Design features using wrapped boundaries
  /gl:slice 1 — Start building new functionality

This boundary is now available for safe refactoring via /gl:slice
using the 'wraps' field.
```

## Integration with Other Commands

### With /gl:assess
- Assessment output drives wrap recommendations
- Prioritizes by risk: untested + external + complex

### With /gl:design
- Wrapped boundaries appear in CONTRACTS.md
- Designer can reference [WRAPPED] contracts when planning slices

### With /gl:slice
- Slices can refactor wrapped boundaries using `wraps` field
- When refactored, STATE.md status updates to "refactored"
- Original locking tests must still pass (or be explicitly updated with user approval)

### With /gl:security
- Security agent runs in "slice" mode during wrap (document only)
- Later runs in normal mode when slice refactors wrapped boundary

## Summary

Wrap is a preparation command. It:
1. Documents existing behaviour as [WRAPPED] contracts
2. Locks that behaviour with locking tests
3. Creates safety net for refactoring
4. Identifies security issues early
5. Makes boundaries available to /gl:design and /gl:slice

Wrap never modifies production code. It only adds documentation and tests.
