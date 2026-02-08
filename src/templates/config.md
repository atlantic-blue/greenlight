# Config Template

Template for `.greenlight/config.json` — project configuration.

## Schema

```json
{
  "version": "1.0.0",
  "mode": "interactive",
  "model_profile": "balanced",
  "model_overrides": {},
  "profiles": {
    "quality": {
      "architect": "opus",
      "designer": "opus",
      "test_writer": "opus",
      "implementer": "opus",
      "security": "opus",
      "debugger": "opus",
      "verifier": "opus",
      "codebase_mapper": "opus"
    },
    "balanced": {
      "architect": "opus",
      "designer": "opus",
      "test_writer": "sonnet",
      "implementer": "sonnet",
      "security": "sonnet",
      "debugger": "sonnet",
      "verifier": "sonnet",
      "codebase_mapper": "sonnet"
    },
    "budget": {
      "architect": "sonnet",
      "designer": "sonnet",
      "test_writer": "sonnet",
      "implementer": "sonnet",
      "security": "haiku",
      "debugger": "sonnet",
      "verifier": "haiku",
      "codebase_mapper": "haiku"
    }
  },
  "workflow": {
    "security_scan": true,
    "visual_checkpoint": true,
    "auto_parallel": true,
    "max_implementation_retries": 3,
    "max_security_retries": 2,
    "run_full_suite_after_slice": true
  },
  "test": {
    "command": "npm test",
    "filter_flag": "--filter",
    "coverage_command": "npm test -- --coverage",
    "security_filter": "security"
  },
  "project": {
    "name": "",
    "stack": "",
    "src_dir": "src",
    "test_dir": "tests"
  }
}
```

## Model Resolution

Every command that spawns an agent resolves the model at runtime:

```
1. Check model_overrides[agent_name] — if set, use it
2. Else check profiles[model_profile][agent_name] — use profile default
3. Else fall back to "sonnet"
```

Example: profile is "balanced", overrides has `"security": "opus"`:
- architect → opus (from balanced profile)
- security → opus (from override, not profile's sonnet)
- implementer → sonnet (from balanced profile)

The `model_overrides` object only contains agents the user has explicitly changed. An empty `model_overrides` means all agents follow the profile.

## Field Definitions

### mode

| Value | Behaviour |
|-------|-----------|
| `interactive` | Pause at visual checkpoints, confirm slice start |
| `yolo` | Skip visual checkpoints, auto-start slices. Decision and external action checkpoints still pause. |

### model_profile

Selects which profile to use as the base model assignment.

| Profile | When to use |
|---------|-------------|
| quality | High-stakes projects, complex domains, when correctness matters most |
| balanced | Default. Good trade-off between cost and quality for most projects |
| budget | Prototypes, learning projects, cost-sensitive work |

### model_overrides

Per-agent overrides that take precedence over the active profile. Only include agents you want to change.

```json
"model_overrides": {
  "security": "opus",
  "implementer": "opus"
}
```

Set an agent to `null` or remove the key to revert to the profile default.

### profiles

The three built-in profiles. Each maps agent names to model identifiers.

Agent names: `architect`, `designer`, `test_writer`, `implementer`, `security`, `debugger`, `verifier`, `codebase_mapper`.

Model identifiers: `opus`, `sonnet`, `haiku`.

### workflow

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| security_scan | boolean | true | Run security agent after each slice |
| visual_checkpoint | boolean | true | Pause for visual verification on UI slices |
| auto_parallel | boolean | true | Suggest parallel slices when available |
| max_implementation_retries | number | 3 | Max attempts to make tests pass |
| max_security_retries | number | 2 | Max attempts to fix security issues |
| run_full_suite_after_slice | boolean | true | Run all tests (not just current slice) after implementation |

### test

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| command | string | "npm test" | Base test command |
| filter_flag | string | "--filter" | Flag to run specific test file |
| coverage_command | string | "npm test -- --coverage" | Command for coverage report |
| security_filter | string | "security" | Filter pattern for security tests |

### project

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| name | string | "" | Project name (set by /gl:init) |
| stack | string | "" | Stack description (set by /gl:init) |
| src_dir | string | "src" | Source code directory |
| test_dir | string | "tests" | Test directory |

## Validation Rules

1. `version` must be a valid semver string
2. `mode` must be "interactive" or "yolo"
3. `model_profile` must be "quality", "balanced", or "budget"
4. `model_overrides.*` values must be "opus", "sonnet", or "haiku"
5. `workflow.*` booleans must be true/false
6. `workflow.max_implementation_retries` must be 1-5
7. `workflow.max_security_retries` must be 1-3
8. `test.command` must be a non-empty string
9. `project.src_dir` and `project.test_dir` must be valid directory names

## Defaults

If config.json is missing or a field is absent, use the balanced profile defaults. Never fail because a config field is missing — always fall back to defaults.

## Reading Config

Every command that needs config should:

```
1. Read .greenlight/config.json
2. Merge with defaults (config values override defaults)
3. Resolve models (overrides > profile > sonnet)
4. Use
```

If config is missing entirely, use balanced profile defaults for all agents.
