# Greenlight Engineering Standards

You are operating under the Greenlight system — a TDD-first, contract-driven development framework.
Every line of code you produce must be verifiable by automated tests. If it can't be tested, it shouldn't exist.

## Core Principles

1. **Contracts before code.** Every feature starts as a typed interface, API schema, or database model. These are the source of truth.
2. **Tests before implementation.** Integration tests are written from contracts before any business logic. Implementation means making tests green.
3. **Vertical slices over horizontal layers.** Build thin paths through the entire stack that deliver user value. Never build "the auth layer" — build "a user can register and log in."
4. **Prove it works, don't say it works.** Verification is `npm test` or `pytest`, never self-review.
5. **Fresh context per agent.** Keep orchestrators thin. Heavy work happens in subagents with fresh 200k context windows.

## Context Degradation Awareness

Claude's quality degrades predictably as context fills:

| Context Usage | Quality | Behaviour |
|---------------|---------|-----------|
| 0-30% | PEAK | Thorough, comprehensive, creative |
| 30-50% | GOOD | Solid, reliable work |
| 50-70% | DEGRADING | Starts cutting corners, less thorough |
| 70%+ | POOR | Rushed, minimal, "completion mode" |

**Rule:** Subagent tasks must complete within ~50% context. If a slice is too large, split it. Orchestrators stay under 30% — they route, they don't build.

## Agent Isolation Rules

These are non-negotiable boundaries between agents:

| Agent | Can See | Cannot See | Cannot Do |
|-------|---------|------------|-----------|
| gl-architect | Requirements, constraints | Implementation code | Write production code |
| gl-test-writer | Contracts, existing test patterns | Implementation code | Write production code |
| gl-implementer | Contracts, test names (not code) | Test source code | Modify test files |
| gl-security | Diffs, contracts | Test implementation details | Fix production code |
| gl-verifier | Everything | N/A | Modify any code |
| gl-debugger | Everything | N/A | Modify tests (without approval) |
| gl-assessor | Codebase docs, test results, standards | N/A (read-only analytical agent) | Modify any code |
| gl-wrapper | Implementation code, existing tests | N/A | Modify production code (only writes contracts and locking tests) |

**Wrapper isolation exception:** gl-wrapper is a deliberate exception. It sees implementation code AND writes locking tests. This is necessary because locking tests must verify what code currently does, not what contracts say it should do. This exception is scoped: only applies to tests in `tests/locking/`. When a boundary is later refactored via /gl:slice, locking tests are deleted and proper integration tests are written under strict isolation.

Violating isolation defeats the purpose of TDD. If an agent sees its own tests, it tests its implementation. If a test writer sees implementation, it tests internals instead of behaviour.

## Code Quality Constraints

These apply to ALL code regardless of language or framework.

### Error Handling
- Every function that can fail MUST have explicit error handling
- No empty catch blocks. Ever. Log, rethrow, or handle — pick one
- No swallowing errors with generic try/catch around entire functions
- Use typed/specific errors, not generic Error("something went wrong")
- Async functions must have error boundaries — no unhandled promise rejections
- Return early on failure conditions, don't nest happy path inside error checks

### Naming
- Functions describe what they DO: `validateUserEmail()` not `process()` or `handle()`
- Booleans read as questions: `isValid`, `hasPermission`, `canRetry`
- Constants are SCREAMING_SNAKE for true constants, not for config
- No abbreviations unless universally understood (id, url, api). Write `response` not `res`, `request` not `req`, `error` not `err`
- File names match the primary export

### Functions
- Single responsibility. If you need "and" to describe it, split it
- Maximum 30 lines per function. If longer, extract
- Maximum 3 parameters. Use an options object beyond that
- No boolean parameters — use named options or separate functions
- Pure functions where possible. Side effects should be explicit and isolated
- Guard clauses at the top, happy path at the bottom

### Security
- Never log sensitive data (passwords, tokens, PII, API keys)
- Validate ALL external input — user input, API responses, file contents
- Use parameterised queries. No string concatenation for SQL. No exceptions
- Secrets come from environment variables, never hardcoded, never in git
- Set explicit CORS, CSP, and security headers. Don't disable them to "fix" a bug
- Hash passwords with bcrypt/argon2, never MD5/SHA for passwords
- Rate limiting on all public endpoints
- HTTPS only. No mixed content. HSTS headers

### API Design
- Consistent response envelope: `{ data, error, meta }`
- HTTP status codes mean what they mean: 400 is client error, 500 is server error, don't return 200 with `{ error: true }`
- Pagination from day one on list endpoints. Default limit, max limit
- Validate request bodies against schemas, reject unknown fields
- Version APIs from the start: `/v1/` prefix

### Database
- Every table has: `id` (UUID preferred), `created_at`, `updated_at`
- Migrations are forward-only with explicit rollback scripts
- Indexes on foreign keys and any column used in WHERE clauses
- No `SELECT *` — explicitly name columns
- Connection pooling. Always
- Soft delete with `deleted_at` unless there's a legal reason to hard delete

### Testing
- Integration tests prove vertical slices work end-to-end
- Unit tests for complex business logic, algorithms, and edge cases
- Test names describe behaviour: `should reject registration with duplicate email`
- No testing implementation details — test behaviour and outputs
- Each test is independent. No shared mutable state between tests
- Use factories/fixtures for test data, not copy-pasted objects
- Test the sad path: timeout, invalid input, missing data, concurrent access

### Circuit Breaker
- After 3 failed attempts on any single test, STOP and produce a structured diagnostic report
- After 7 total failed attempts across all tests in a slice, STOP regardless of per-test counts
- Before modifying any file, verify it is within inferred scope from contracts; justify out-of-scope changes
- Full protocol: `references/circuit-breaker.md`

### Logging & Observability
- Structured logging (JSON) in production, human-readable in dev
- Every request gets a correlation ID
- Log at boundaries: incoming requests, outgoing calls, errors
- Log levels mean something: ERROR = needs attention, WARN = degraded, INFO = business events, DEBUG = dev only
- No `console.log` in production code. Use a proper logger

### File & Project Structure
- Group by feature/domain, not by type. `users/controller.ts` not `controllers/users.ts`
- One export per file for classes and components
- Index files only for public API re-exports
- Keep dependencies explicit — no circular imports
- Config in one place, loaded once, validated on startup

### Git
- Atomic commits: one logical change per commit
- Conventional commits: `feat:`, `fix:`, `test:`, `refactor:`, `docs:`, `chore:`
- Never commit generated files, build artifacts, or environment files
- Never use `git add .` or `git add -A` — stage specific files

### Performance
- Measure before optimising. Don't guess at bottlenecks
- N+1 queries are bugs. Use eager loading or batch queries
- Set timeouts on ALL external calls — HTTP, database, queues
- Paginate. Stream large datasets. Never load unbounded data into memory

## Deviation Rules

When agents discover unplanned work during execution, follow the deviation protocol in `references/deviation-rules.md`. Summary:

1. **Auto-fix bugs** — fix immediately, track
2. **Auto-add critical missing functionality** — add immediately, track
3. **Auto-fix blocking issues** — fix to unblock, track
4. **STOP for architectural changes** — report and wait for decision

## What NOT To Do

- Don't add comments that restate the code
- Don't create abstractions for things that happen once. Three uses = extract pattern
- Don't use `any` type in TypeScript
- Don't disable linter rules. Fix the code instead
- Don't mock everything in tests. Mock external boundaries only
- Don't optimise prematurely. Make it work, make it right, make it fast — in that order
- Don't build for scale you don't have
