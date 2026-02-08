# Greenlight

**TDD-first development system for Claude Code.**

Tests are the source of truth. Green means done. Security is built in, not bolted on.

## What This Is

Greenlight is a set of Claude Code slash commands that enforce test-driven development when building with AI. Instead of asking Claude to write code and then review its own work, Greenlight separates concerns: one agent writes tests from contracts, another implements until tests pass, a third scans for security vulnerabilities, a fourth verifies contract satisfaction, and the test runner is the only judge.

## How It Works

```
/gl:map             Analyse existing codebase (brownfield only)
/gl:init            Define contracts, dependency graph, scaffold
/gl:slice 1         Tests → implement → security → verify → commit
/gl:slice 2         Repeat for each vertical slice
/gl:status          Real progress (tests passing, not Claude's opinion)
/gl:ship            Full security audit + deploy readiness
```

### The TDD Loop

Every slice follows the same cycle:

1. **Contract** — typed interfaces define what the system does (WHAT, never HOW)
2. **Tests** (Agent A) — integration tests written from contracts, never sees implementation
3. **Implementation** (Agent B) — fresh context, makes tests pass, never modifies tests
4. **Security Scan** (Agent C) — reviews diff, writes failing tests for vulnerabilities
5. **Security Fix** (Agent B) — makes security tests pass too
6. **Verification** (Agent D) — goal-backward analysis confirms contracts are satisfied
7. **Green** — test runner is the single source of truth

Agents A, B, C, and D never see each other's work. This prevents Claude from gaming its own tests.

### Vertical Slices

Greenlight builds thin paths through the entire stack that deliver user value:

- Slice 1: "A user can register with email" (API + DB + response)
- Slice 2: "A user can log in and get a token" (depends on Slice 1)
- Slice 3: "A user can see their dashboard" (depends on Slice 2)

Each slice is independently testable, committable, and deployable.

### Agent Isolation

| Agent | Can See | Cannot See | Cannot Do |
|-------|---------|------------|-----------|
| Architect | Requirements | Implementation | Write production code |
| Test Writer | Contracts | Implementation | Write production code |
| Implementer | Contracts, test names | Test source code | Modify test files |
| Security | Diffs, contracts | Test details | Fix production code |
| Verifier | Everything | N/A | Modify any code |
| Debugger | Everything | N/A | Modify tests (without approval) |

### Engineering Standards

`CLAUDE.md` enforces production-grade quality on every line:

- Explicit error handling, typed everything, functions under 30 lines
- Parameterised queries, structured logging, security defaults
- Context degradation awareness — agents stay under 50% context usage
- Deviation rules — auto-fix bugs, stop for architectural decisions

### Security Built In

Security isn't a phase at the end. Every slice gets a security scan. Every vulnerability becomes a failing test. Every fix is verified by the test runner. `/gl:ship` runs a full codebase audit before deployment.

## Installation

```bash
# Global (all projects)
bash install.sh --global

# Local (current project only)
bash install.sh --local

# Check installation
bash install.sh --check

# Uninstall
bash install.sh --global --uninstall
```

Verify with `/gl:help` inside Claude Code.

## Commands

### Setup
| Command | Description |
|---------|-------------|
| `/gl:init` | Interview, contracts, dependency graph, scaffold |
| `/gl:map` | Analyse existing codebase (brownfield) |
| `/gl:settings` | Configure model profiles, mode, workflow |

### Build
| Command | Description |
|---------|-------------|
| `/gl:slice <N>` | TDD loop: test, implement, security, verify, commit |
| `/gl:quick` | Ad-hoc tasks with test guarantees |
| `/gl:add-slice` | Add new slices to the dependency graph |

### Monitor
| Command | Description |
|---------|-------------|
| `/gl:status` | Progress dashboard from test results |
| `/gl:pause` | Save state for next session |
| `/gl:resume` | Restore and continue |

### Ship
| Command | Description |
|---------|-------------|
| `/gl:ship` | Full security audit + deploy readiness |

## Agents

| Agent | Role |
|-------|------|
| `gl-architect` | Produces typed contracts and dependency graphs |
| `gl-test-writer` | Writes tests from contracts (never sees implementation) |
| `gl-implementer` | Makes tests pass (never modifies tests) |
| `gl-security` | Reviews diffs, writes security tests (never fixes code) |
| `gl-verifier` | Goal-backward contract satisfaction verification |
| `gl-debugger` | Scientific method debugging with hypothesis testing |
| `gl-codebase-mapper` | Analyses existing code for brownfield projects |

## Configuration

### Model Profiles

| Profile | Architect | Tests | Implement | Security | Verify | Debug | Map |
|---------|-----------|-------|-----------|----------|--------|-------|-----|
| quality | opus | opus | opus | opus | opus | opus | opus |
| balanced | opus | sonnet | sonnet | sonnet | sonnet | sonnet | sonnet |
| budget | sonnet | sonnet | sonnet | haiku | haiku | sonnet | haiku |

### Modes

- **interactive** — confirm each step (default)
- **yolo** — auto-approve visual checkpoints; decisions and external actions still pause

### Workflow Toggles

- `security_scan` — per-slice security agent (default: on)
- `visual_checkpoint` — pause for UI verification (default: on)
- `auto_parallel` — suggest parallel slices (default: on)

## Project Structure

```
.claude/
  commands/gl/          Slash commands (11 commands)
  agents/               Agent definitions (7 agents)
  references/           Shared protocols (3 docs)
  templates/            Schema templates (2 docs)

.greenlight/            Project state (created by /gl:init)
  STATE.md              Slice progress tracker
  GRAPH.json            Dependency DAG
  CONTRACTS.md          Contract definitions
  QUICK.md              Ad-hoc task log
  config.json           Settings
  codebase/             Brownfield analysis (7 docs)

CLAUDE.md               Engineering standards
```

## Credits

Greenlight was inspired by [Get Shit Done](https://github.com/glittercowboy/get-shit-done). GSD's innovations in context engineering, subagent orchestration, and session management provided the foundation. Greenlight builds on these with TDD-first agent isolation, contract-driven planning, and security-as-tests.

## License

MIT
