# Verification Patterns

How Greenlight verifies that slices actually work — not that agents said they work.

## Core Principle

**Tests are the only judge.** Claude reviewing Claude is not verification. The test runner is the single source of truth. A slice is done when all tests pass, not when an agent reports success.

**See also:** `references/verification-tiers.md` — defines the `auto` and `verify` tier system that determines whether a human acceptance checkpoint is required after tests pass.

## Three Levels of Verification

### Level 1: Existence

The artifact exists and is non-empty.

```bash
# File exists
[ -f src/users/service.ts ] && echo "exists" || echo "MISSING"

# File has content (not just imports/boilerplate)
wc -l src/users/service.ts  # should be > 10 lines for real implementation
```

**Catches:** Missing files, forgotten implementations, empty stubs.

### Level 2: Substantive

The artifact contains real implementation, not stubs or placeholders.

**Stub Detection Patterns:**

```bash
# Generic returns (function exists but does nothing useful)
grep -n "return \[\]" src/    # empty array returns
grep -n "return {}" src/      # empty object returns
grep -n "return null" src/    # null returns (when contract says otherwise)
grep -n "return undefined" src/
grep -n "// TODO" src/        # placeholder comments
grep -n "throw new Error('Not implemented')" src/

# Framework-specific stubs
# React: Component renders nothing meaningful
grep -n "return null" src/ --include="*.tsx"
grep -n "return <></>" src/ --include="*.tsx"
grep -n "return <div />" src/ --include="*.tsx"

# API: Handler that doesn't process anything
grep -n "res.json({})" src/
grep -n "res.send()" src/
grep -n "res.status(200).end()" src/

# Database: Schema exists but queries are hardcoded
grep -n "SELECT 1" src/
grep -n "INSERT INTO.*VALUES.*null" src/
```

**Catches:** Pass-through functions, empty handlers, placeholder returns.

### Level 3: Wired

The artifact is connected to the rest of the system — it's reachable, called, and integrated.

**Wiring Checks:**

```bash
# Component is imported somewhere
grep -rn "import.*UserService" src/ --include="*.ts"

# Route is registered
grep -rn "router\.\(get\|post\|put\|delete\).*users" src/

# Database model is used in a query
grep -rn "UserModel\.\(find\|create\|update\|delete\)" src/

# Environment variable is read
grep -rn "process.env.DATABASE_URL" src/

# Middleware is applied
grep -rn "app.use.*auth" src/
```

**Catches:** Orphaned modules, dead code, implemented-but-not-connected features.

## Contract Coverage Verification

Every contract in CONTRACTS.md must have:

1. **At least one integration test** that exercises the contract's happy path
2. **At least one test per error state** defined in the contract
3. **Production code** that the tests call (not mocked away)

```bash
# For each contract, verify test exists
# Contract: CreateUser → tests/integration/user-registration.test.ts should contain "CreateUser" or "create user"
grep -rn "CreateUser\|create.*user\|registration" tests/integration/

# Verify production code exists
grep -rn "export.*createUser\|export.*CreateUser" src/
```

## Test Quality Verification

Not all tests are equal. Verify tests are meaningful:

### Tests That Prove Nothing

```javascript
// BAD: Tests implementation, not behaviour
it('should call database.insert', () => {
  expect(mockDb.insert).toHaveBeenCalled()  // proves nothing about correctness
})

// BAD: Tautological test
it('should return what it returns', () => {
  const result = fn()
  expect(result).toBeDefined()  // everything is "defined"
})

// BAD: Tests mock, not real system
it('should create user', () => {
  mockService.createUser.mockResolvedValue(mockUser)
  const result = await mockService.createUser(input)
  expect(result).toEqual(mockUser)  // you're testing the mock
})
```

### Tests That Prove Behaviour

```javascript
// GOOD: Tests actual contract behaviour
it('should return 201 and user object when registration succeeds', async () => {
  const response = await api.post('/v1/users', validUserInput)
  expect(response.status).toBe(201)
  expect(response.body.data).toMatchObject({
    email: validUserInput.email,
    name: validUserInput.name
  })
  expect(response.body.data.password).toBeUndefined()  // never expose
})

// GOOD: Tests error contract
it('should return 409 when email already exists', async () => {
  await api.post('/v1/users', validUserInput)  // create first
  const response = await api.post('/v1/users', validUserInput)  // duplicate
  expect(response.status).toBe(409)
  expect(response.body.error.code).toBe('EMAIL_EXISTS')
})
```

## Regression Detection

After every implementation step, verify no existing tests broke:

```bash
# Run FULL suite, not just current slice
npm test 2>&1

# Check exit code
echo $?  # 0 = all pass, non-zero = regression

# If regression detected:
# 1. Identify which tests broke
# 2. Check if they belong to a different slice
# 3. If yes → the implementation has a cross-slice side effect
# 4. Fix without modifying tests from other slices
```

## Security Verification

Security tests follow the same pattern but target vulnerability scenarios:

```bash
# Run security tests specifically
npm test -- --filter security 2>&1

# Check for known vulnerability patterns in new code
# (gl-security agent does this, but orchestrator spot-checks)
grep -rn "eval(" src/ --include="*.ts" --include="*.js"
grep -rn "innerHTML" src/ --include="*.tsx" --include="*.jsx"
grep -rn "dangerouslySetInnerHTML" src/ --include="*.tsx"
grep -rn "exec(" src/ --include="*.ts" --include="*.js"
```

## Verification Timing

| Event | What to Verify |
|-------|---------------|
| After test writer finishes | All new tests FAIL (confirms they're testing something real) |
| After implementer finishes | All tests PASS (functional + security from prior slices) |
| After security scan | New security tests FAIL (confirms vulnerabilities exist) |
| After security fix | ALL tests PASS (functional + security, old + new) |
| Before slice marked complete | Full suite green, no regressions, no stubs |
| Before ship | Full suite + security audit + dependency audit |

## Automated Verification Script

Orchestrators can run this quick health check:

```bash
#!/bin/bash
# Quick verification: does this slice look complete?

echo "=== Greenlight Verification ==="

# 1. All tests pass
echo -n "Tests: "
npm test --silent 2>&1 && echo "PASS" || echo "FAIL"

# 2. No stubs in src/
echo -n "Stubs: "
STUBS=$(grep -rn "TODO\|Not implemented\|FIXME\|HACK" src/ 2>/dev/null | wc -l | tr -d ' ')
[ "$STUBS" -eq 0 ] && echo "CLEAN" || echo "FOUND ($STUBS)"

# 3. No console.log in production
echo -n "Console: "
LOGS=$(grep -rn "console\.log" src/ --include="*.ts" --include="*.js" 2>/dev/null | grep -v "logger" | wc -l | tr -d ' ')
[ "$LOGS" -eq 0 ] && echo "CLEAN" || echo "FOUND ($LOGS)"

# 4. Lint passes
echo -n "Lint: "
npm run lint --silent 2>&1 && echo "PASS" || echo "FAIL"

echo "==========================="
```
