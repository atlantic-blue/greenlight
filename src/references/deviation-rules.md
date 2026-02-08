# Deviation Rules

Shared protocol for handling unplanned work discovered during execution.

All implementing agents (gl-implementer, gl-debugger) follow these rules. Orchestrators (gl:slice, gl:quick) reference them to validate agent behaviour.

## The Four Rules

### Rule 1: Auto-fix Bugs (ALWAYS)

**Trigger:** Code doesn't work as intended — broken behaviour, errors, incorrect output.

**Action:** Fix immediately, track for summary.

**Examples:**
- Wrong query returns incorrect data
- Logic errors in conditionals
- Type errors blocking compilation
- Null pointer exceptions on valid input
- Broken validation regex
- Off-by-one errors

**Process:**
1. Fix the bug inline
2. Run affected tests to verify fix
3. Continue current task
4. Track: `[BUG-FIX] {file}: {description}`

**No user permission needed.** Bugs must be fixed for correct operation.

### Rule 2: Auto-add Critical Missing Functionality (ALWAYS)

**Trigger:** Code is missing essential features for correctness or security.

**Action:** Add immediately, track for summary.

**Examples:**
- Missing input validation on user-facing endpoint
- No error handling on async operation
- Missing null/undefined checks on external data
- No auth check on protected route
- Missing CSRF protection on state-changing endpoint
- No rate limiting on auth endpoints
- Missing error logging for debugging
- No timeout on external HTTP calls

**Process:**
1. Add the missing functionality
2. Run tests to verify it works
3. Continue current task
4. Track: `[CRITICAL-ADD] {file}: {description}`

**No user permission needed.** These are requirements for basic correctness and security.

### Rule 3: Auto-fix Blocking Issues (ALWAYS)

**Trigger:** Something prevents completing the current task.

**Action:** Fix immediately to unblock, track for summary.

**Examples:**
- Missing dependency not in package.json
- Wrong types blocking TypeScript compilation
- Broken import paths from file moves
- Missing environment variable for test setup
- Build configuration error
- Test fixture missing required field
- Circular dependency preventing module load

**Process:**
1. Fix the blocking issue
2. Verify the current task can now proceed
3. Continue
4. Track: `[UNBLOCK] {file}: {description}`

**No user permission needed.** Can't complete task without fixing blocker.

### Rule 4: STOP for Architectural Changes (ALWAYS)

**Trigger:** Fix requires significant structural modification beyond the current slice.

**Action:** STOP execution, report to orchestrator, wait for decision.

**Examples:**
- Adding a new database table or collection
- Major schema changes (column type change, new required field on existing model)
- Switching libraries (e.g., replacing date-fns with dayjs)
- Changing authentication approach
- Adding new infrastructure (cache layer, message queue)
- Breaking changes to existing API contracts
- Modifying contracts from completed slices
- Adding new external service dependency

**Process:**
1. STOP current task immediately
2. Report to orchestrator:
   - What you found
   - Why the current approach doesn't work
   - Proposed architectural change
   - Impact on existing slices and contracts
   - Alternatives considered
3. Wait for user decision
4. Track: `[ARCH-STOP] {description}`

**User decision required.** These affect system design beyond the current slice.

## Priority Order

When multiple rules could apply, use the highest-priority match:

```
Rule 4 (STOP) > Rule 1 (Bug) = Rule 2 (Critical) = Rule 3 (Blocker)
```

1. If Rule 4 applies → STOP immediately (architectural decision needed)
2. If Rules 1-3 apply → fix automatically, track for summary
3. If genuinely unsure which rule → apply Rule 4 (better to ask than break things)

## Deviation Summary Format

At the end of every task, include a deviation summary:

```markdown
## Deviations

| Type | File | Description |
|------|------|-------------|
| BUG-FIX | src/auth/validate.ts | Fixed regex that didn't handle + in emails |
| CRITICAL-ADD | src/api/users.ts | Added input validation on POST /users |
| UNBLOCK | package.json | Added missing zod dependency |

Total: 3 deviations (0 architectural stops)
```

If no deviations: `No deviations from plan.`

## Commit Protocol for Deviations

- Rules 1-3: Include deviation in the task's commit message body
- Rule 4: Do NOT commit. Report and wait.

```bash
git commit -m "feat(user-registration): implement CreateUser endpoint

- [BUG-FIX] Fixed email regex to handle + addresses
- [CRITICAL-ADD] Added rate limiting on registration endpoint
"
```
