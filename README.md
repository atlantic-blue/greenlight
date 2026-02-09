# Greenlight

[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux-blue.svg)]()
[![Architecture](https://img.shields.io/badge/arch-amd64%20%7C%20arm64-blue.svg)]()

TDD-first development system for [Claude Code](https://docs.anthropic.com/en/docs/claude-code). Tests are the source of truth. Green means done. Security is built in, not bolted on.

<img width="613" height="434" alt="Screenshot 2026-02-08 at 08 52 11" src="https://github.com/user-attachments/assets/145e1553-7b4b-48d1-a77d-a2052c60d3d1" />

<img width="1470" height="671" alt="Screenshot 2026-02-08 at 08 50 32" src="https://github.com/user-attachments/assets/5a33ab58-e576-4a80-8a88-614b2e45ea46" />

## Quick Start

```bash
npx greenlight-cc install --global
```

Or with Go:

```bash
go install github.com/atlantic-blue/greenlight@latest
greenlight install --global
```

Then inside Claude Code, run `/gl:init` to get started.

## Why Greenlight

Claude Code is powerful but permissive. it writes code, reviews its own work, and moves on. There's no external verification. Greenlight fixes this by enforcing **agent isolation**: one agent writes tests from contracts, another implements until tests pass, a third scans for vulnerabilities, and the test runner is the only judge. No agent ever sees its own tests.

This matters because:

- **AI reviews its own code poorly.** Separating test writing from implementation prevents Claude from gaming its own tests.
- **TDD works better with AI.** Contracts define WHAT, agents figure out HOW, and the test suite proves it works.
- **Security can't be an afterthought.** Every slice gets a security scan. Every vulnerability becomes a failing test.

## What Greenlight Is

A set of Claude Code slash commands, agents, and engineering standards that enforce test-driven development:

- **10 agents** with strict isolation boundaries (designer, architect, test writer, implementer, security, verifier, debugger, codebase mapper, assessor, wrapper)
- **16 slash commands** (`/gl:init`, `/gl:design`, `/gl:slice`, `/gl:ship`, etc.) that orchestrate the workflow
- **Engineering standards** (`CLAUDE.md`) covering error handling, naming, security, API design, testing, and more
- **Context degradation awareness** agents stay under 50% context usage to maintain quality

## What Greenlight Is Not

- Not a framework, library, or runtime dependency. It installs config files and does nothing at build time
- Not a replacement for your test framework. It orchestrates Claude Code to use whatever test runner your project already has
- Not opinionated about language or stack. The contracts and standards are language-agnostic

## Install

### From GitHub Releases (recommended)

Download the latest binary for your platform from [Releases](https://github.com/atlantic-blue/greenlight/releases), extract it, and add it to your PATH.

```bash
# macOS / Linux
tar xzf greenlight_*.tar.gz
sudo mv greenlight /usr/local/bin/
```

### From Source

Requires [Go 1.24+](https://go.dev/dl/).

```bash
go install github.com/atlantic-blue/greenlight@latest
```

### Usage

```bash
# Install globally (all projects)
greenlight install --global

# Install locally (current project only)
greenlight install --local

# Verify installation
greenlight check --global
greenlight check --local

# If you already have a CLAUDE.md
greenlight install --local --on-conflict=keep      # save as CLAUDE_GREENLIGHT.md (default)
greenlight install --local --on-conflict=append     # append to existing
greenlight install --local --on-conflict=replace    # backup + overwrite

# Remove
greenlight uninstall --global
greenlight uninstall --local
```

Verify with `/gl:help` inside Claude Code.

## How It Works

### The TDD Loop

Every slice follows the same cycle:

1. **Contract** typed interfaces define what the system does (WHAT, never HOW)
2. **Tests** (Agent A) integration tests written from contracts, never sees implementation
3. **Implementation** (Agent B) fresh context, makes tests pass, never modifies tests
4. **Security Scan** (Agent C) reviews diff, writes failing tests for vulnerabilities
5. **Security Fix** (Agent B) makes security tests pass too
6. **Verification** (Agent D) goal-backward analysis confirms contracts are satisfied
7. **Green** test runner is the single source of truth

### Vertical Slices

Greenlight builds thin paths through the entire stack that deliver user value:

- Slice 1: "A user can register with email" (API + DB + response)
- Slice 2: "A user can log in and get a token" (depends on Slice 1)
- Slice 3: "A user can see their dashboard" (depends on Slice 2)

Each slice is independently testable, committable, and deployable.

### Agent Isolation

| Agent | Can See | Cannot See | Cannot Do |
|-------|---------|------------|-----------|
| Designer | Interview context, requirements | Implementation, contracts | Write code or contracts |
| Architect | Requirements, DESIGN.md | Implementation | Write production code |
| Test Writer | Contracts | Implementation | Write production code |
| Implementer | Contracts, test names | Test source code | Modify test files |
| Security | Diffs, contracts | Test details | Fix production code |
| Verifier | Everything | N/A | Modify any code |
| Debugger | Everything | N/A | Modify tests (without approval) |
| Assessor | Codebase docs, test results | N/A (read-only) | Modify any code |
| Wrapper | Implementation code, existing tests | N/A | Modify production code (only writes locking tests) |

## Commands

### Setup
| Command | Description |
|---------|-------------|
| `/gl:init` | Brief interview, project config, scaffold |
| `/gl:design` | System design session: requirements, research, architecture |
| `/gl:map` | Analyse existing codebase (brownfield) |
| `/gl:settings` | Configure model profiles, mode, workflow |

### Brownfield
| Command | Description |
|---------|-------------|
| `/gl:assess` | Gap analysis and risk assessment for existing codebases |
| `/gl:wrap` | Extract contracts and locking tests from existing boundaries |

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
| `/gl:roadmap` | Display roadmap, plan milestones, archive completed work |
| `/gl:changelog` | Chronological changelog from summaries, with filtering |
| `/gl:pause` | Save state for next session |
| `/gl:resume` | Restore and continue |

### Ship
| Command | Description |
|---------|-------------|
| `/gl:ship` | Full security audit + deploy readiness |

## Typical Flow

### Greenfield
1. `/gl:init` — brief interview, scaffold project config
2. `/gl:design` — system design: requirements, research, architecture, contracts
3. `/gl:slice 1` — TDD loop: test, implement, security scan, verify, commit
4. `/gl:slice 2` ... `/gl:slice N` — repeat for each slice
5. `/gl:ship` — full security audit + deploy readiness

### Brownfield (existing codebase)
1. `/gl:map` — analyse existing codebase structure
2. `/gl:assess` — gap analysis, identify risks and missing coverage
3. `/gl:wrap` — extract contracts and locking tests from existing boundaries
4. `/gl:init` — brief interview, scaffold project config
5. `/gl:design` — system design incorporating existing architecture
6. `/gl:slice 1` ... `/gl:slice N` — build new features with TDD
7. `/gl:ship` — full security audit + deploy readiness

### Ongoing
- `/gl:quick` — ad-hoc bug fixes and small features (still test-first)
- `/gl:roadmap milestone` — plan next milestone, add slices to the graph
- `/gl:changelog` — view completed work history
- `/gl:roadmap archive` — archive completed milestones

## Configuration

### Model Profiles

Three built-in profiles control which Claude model each agent uses. Choose based on your priorities:

| Profile | Designer | Architect | Tests | Implement | Security | Verify | Debug | Map |
|---------|----------|-----------|-------|-----------|----------|--------|-------|-----|
| quality | opus | opus | opus | opus | opus | opus | opus | opus |
| balanced | opus | opus | sonnet | sonnet | sonnet | sonnet | sonnet | sonnet |
| budget | sonnet | sonnet | sonnet | sonnet | haiku | sonnet | haiku | haiku |

```bash
/gl:settings profile balanced        # switch profile
/gl:settings model security opus     # override one agent
/gl:settings model security reset    # revert to profile default
```

Per-agent overrides take precedence over the profile. Resolution order: `model_overrides[agent]` > `profiles[profile][agent]` > `sonnet`.

### Why These Defaults

**Designer and architect default to opus.** These agents make decisions that cascade through everything downstream. A bad architectural choice or missed requirement means every slice built on top of it is wrong. The cost of opus here is small compared to reworking multiple slices.

**Implementation agents default to sonnet.** The test writer, implementer, and debugger follow contracts and standards. Their work is constrained by inputs (contracts) and verified by outputs (test runner). Sonnet handles "make this test pass" and "find the bug" reliably. The TDD loop catches quality issues regardless of model.

**Security defaults to sonnet but consider upgrading.** Most security checks are mechanical: SQL injection, missing auth, XSS. Sonnet handles these well. But if your project handles financial data, health records, or PII, `opus` is worth it for catching subtle auth bypass and business logic vulnerabilities.

**Codebase mapper defaults to sonnet.** It reads and documents existing code, not making decisions. Sonnet is thorough enough for this.

This is about cost control, not capability hiding. Every agent works with any model. The profiles just set sensible defaults for common priorities.

### Modes

- **interactive** confirm each step (default)
- **yolo** auto-approve visual checkpoints; decisions and external actions still pause

## Project Structure

After `greenlight install` and `/gl:init`:

```
.claude/
  commands/gl/          Slash commands (16 commands)
  agents/               Agent definitions (10 agents)
  references/           Shared protocols (3 docs)
  templates/            Schema templates (2 docs)

.greenlight/            Project state (created during setup)
  INTERVIEW.md          Brief interview context (from /gl:init)
  DESIGN.md             System design document (from /gl:design)
  CONTRACTS.md          Typed contracts (from architect)
  GRAPH.json            Dependency DAG
  STATE.md              Slice progress tracker
  ROADMAP.md            Product roadmap with milestone tracking
  DECISIONS.md          Decision log with source tracing
  QUICK.md              Ad-hoc task history
  summaries/            Per-slice, per-wrap, and quick task summaries
  config.json           Settings

CLAUDE.md               Engineering standards
```

## Roadmap

- [ ] `greenlight upgrade` — in-place upgrade preserving user config
- [ ] Homebrew tap (`brew install atlantic-blue/tap/greenlight`)
- [ ] Curl install script for environments without Go
- [ ] `greenlight doctor` — diagnose common Claude Code configuration issues
- [ ] Plugin system for custom agents and commands
- [ ] Telemetry-free usage analytics (opt-in, local-only)

## Why We Built This

Greenlight started from working extensively with [Get Shit Done (GSD)](https://github.com/glittercowboy/get-shit-done). GSD pioneered important ideas, context engineering for LLMs, subagent orchestration with fresh context windows, session persistence across resets, deviation handling, and structured project execution with Claude Code. We learned a lot from it and credit it as the foundation that made Greenlight possible.

But as we used GSD on real projects, we kept running into a set of structural problems that couldn't be solved by configuration or tweaking, they required a different architecture.

### What We Wanted to Fix

**Verification depends on judgment, not proof.** GSD's verification step spawns another agent to review the implementation. This is better than no review, but it's fundamentally the same model as the one that produced the code. It shares Claude's reasoning patterns, blind spots, and tendency to accept plausible-looking output. When the implementer introduces a subtle bug, the verifier is likely to miss it for the same reasons. There's no external oracle. We wanted verification that doesn't depend on any agent's opinion.

**Agent boundaries are structural but not informational.** GSD deserves credit for using separate subagents with fresh context windows. That's a genuine innovation for managing context degradation. But there's no enforced boundary on what information flows between agents. The executor can see tests. The verifier sees the implementation it's checking. This means there's nothing preventing Claude from writing code that satisfies the specific test assertions it can read rather than the underlying requirements. The separation between "what should it do" and "does it do it" can collapse without either agent being aware.

**Security lives at the end of the pipeline.** GSD treats security as part of verification. A review step after implementation is complete. Vulnerabilities found late are expensive because you've already built on top of the insecure code. More importantly, a post-hoc security review tends to focus on what's obviously wrong rather than what's subtly missing, because the reviewer is anchored by working code.

**Phases encourage horizontal construction.** GSD's phase model naturally leads to building in layers. Database schema in one phase, API in another, frontend in a third. This means nothing works end-to-end until late in the project. Integration issues surface late. You can't deploy, demo, or validate incrementally. And because later phases depend on assumptions made in earlier ones, errors compound.

**Context degradation goes unmanaged.** Past 50–60% of Claude's context window, output quality degrades measurably. Shorter implementations, fewer edge cases handled, shallower reasoning. GSD's subagent architecture helps by giving each agent a fresh window, but there's no explicit budget for how much context a task should consume. Long phases can still push agents past the degradation threshold, producing worse code at the end than at the start.

### What Greenlight Does Differently

**The test runner is the only judge.** No agent decides whether the code is correct `npm test` or `pytest` does. Tests are written from typed contracts by an agent that never sees implementation. Implementation is done by an agent that never sees test source code. If the tests pass and the contracts are satisfied, the slice is done. If not, it isn't. Verification is mechanical, not judgmental.

**Agent isolation is enforced by information boundaries.** The test writer receives contracts and produces tests. The implementer receives contracts and test names, not test code. The security agent receives diffs and produces failing tests, it cannot fix code. Each agent operates in a fresh context window with only the information it needs. This isn't a convention that agents are asked to follow; it's a constraint on what they're given.

**Security is a failing test, not a review.** Every vertical slice gets a security scan. Every vulnerability the security agent identifies becomes a failing test case. The implementer makes that test pass through the same TDD loop as everything else. Security isn't a phase or a gate, it's a red test that needs to go green, verified by the same test runner that checks everything else.

**Vertical slices from day one.** Each slice is a thin path through the entire stack that delivers one piece of user value. "A user can register with email" touches the API, database, validation, and response; all working together, all testable end to end. You can deploy it, demo it, and build the next slice on top of it independently. A dependency graph determines build order, and slices without dependencies can run in parallel.

**Explicit context budgets.** Agents are designed to complete their work within 50% of the context window. If a slice is too large for that budget, it gets split before work begins. Orchestrators stay under 30% they route work to subagents, they don't do the work themselves. This keeps output quality consistent from first slice to last.

### Standing on GSD's Shoulders

To be clear: Greenlight wouldn't exist without GSD. The ideas of using markdown files as agent prompts, orchestrating subagents through Claude Code's Task tool, managing state across sessions with handoff files, and building structured deviation rules for when agents encounter unexpected problems. These all came from GSD, and they're good ideas. Greenlight's contribution is layering TDD, contract first design, agent isolation, and integrated security on top of that foundation.

GSD is optimised for speed and shipping. Greenlight is optimised for confidence and correctness. Different tools for different priorities.

## Contributing

This project uses [Conventional Commits](https://www.conventionalcommits.org/). All commit messages must follow the format:

```
<type>(<optional scope>): <description>
```

Types: `feat`, `fix`, `docs`, `chore`, `refactor`, `test`, `ci`

Examples:
- `feat: add user authentication`
- `fix(installer): handle nil pointer`
- `feat!: redesign config format` (breaking change)

Commits are validated locally via [lefthook](https://github.com/evilmartians/lefthook) and in CI via [commitlint](https://commitlint.js.org/).

```bash
brew install lefthook
lefthook install
```

## License

[MIT](LICENSE)
