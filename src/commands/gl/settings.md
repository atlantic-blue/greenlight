---
name: gl:settings
description: Configure model profiles, per-agent overrides, mode, and workflow options
allowed-tools: [Read, Write, Bash, AskUserQuestion]
---

# Greenlight: Settings

View and modify `.greenlight/config.json`.

**Read templates/config.md** for schema reference.

## Parse Arguments

If the user passed arguments, route directly:

- `/gl:settings` (no args) → [Display Current](#display-current)
- `/gl:settings profile <name>` → [Switch Profile](#switch-profile)
- `/gl:settings model <agent> <model>` → [Override Agent Model](#override-agent-model)
- `/gl:settings model <agent> reset` → [Reset Agent Override](#reset-agent-override)
- `/gl:settings mode <interactive|yolo>` → [Switch Mode](#switch-mode)

## Display Current

```bash
cat .greenlight/config.json 2>/dev/null || echo "No config found. Run /gl:init first."
```

If no config, stop here.

Resolve the effective model for each agent:
1. Check `model_overrides[agent]` — if set, use it (mark as "override")
2. Else check `profiles[model_profile][agent]` — use it (mark as "profile")
3. Else `sonnet` (mark as "default")

Display:

```
┌─────────────────────────────────────────────────────┐
│  GREENLIGHT SETTINGS                                │
├─────────────────────────────────────────────────────┤
│                                                     │
│  Profile: {model_profile}                           │
│  Mode: {mode}                                       │
│                                                     │
│  MODELS (override > profile > default)              │
│                                                     │
│  architect        {model}    ({source})             │
│  designer         {model}    ({source})             │
│  test_writer      {model}    ({source})             │
│  implementer      {model}    ({source})             │
│  security         {model}    ({source})             │
│  debugger         {model}    ({source})             │
│  verifier         {model}    ({source})             │
│  codebase_mapper  {model}    ({source})             │
│  assessor         {model}    ({source})             │
│  wrapper          {model}    ({source})             │
│                                                     │
│  To change:                                         │
│  /gl:settings profile <name>                        │
│  /gl:settings model <agent> <model>                 │
│  /gl:settings model <agent> reset                   │
│  /gl:settings mode <interactive|yolo>               │
└─────────────────────────────────────────────────────┘
```

Where `{source}` is one of:
- `profile` — from the active profile
- `override` — user set this explicitly
- `default` — fallback, no profile or override matched

## Switch Profile

Valid profiles: `quality`, `balanced`, `budget`.

Show what will change:

```
Switching to {new_profile}:

  architect        {old} → {new}
  designer         {old} → {new}
  test_writer      {old} → {new}
  implementer      {old} → {new}
  security         {old} → {new}
  debugger         {old} → {new}
  verifier         {old} → {new}
  codebase_mapper  {old} → {new}
  assessor         {old} → {new}
  wrapper          {old} → {new}

Note: {N} agent overrides are still active and take precedence.
{list any overrides}

Apply? [y/N]
```

If confirmed, update `model_profile` in config.json. Do NOT clear existing `model_overrides` — they take precedence by design.

## Override Agent Model

Valid agents: `architect`, `designer`, `test_writer`, `implementer`, `security`, `debugger`, `verifier`, `codebase_mapper`, `assessor`, `wrapper`.

Valid models: `opus`, `sonnet`, `haiku`.

```
Setting {agent} model override: {model}

  {agent}: {old_effective} → {new_model} (override)

Note: This overrides the {profile} profile default of {profile_default}.
To revert: /gl:settings model {agent} reset

Apply? [y/N]
```

If confirmed, set `model_overrides[agent] = model` in config.json.

## Reset Agent Override

Remove a per-agent override so the agent reverts to the profile default.

```
Removing {agent} override.

  {agent}: {current_override} → {profile_default} (profile)

Apply? [y/N]
```

If confirmed, remove the key from `model_overrides` in config.json.

## Switch Mode

Valid modes: `interactive`, `yolo`.

```
Switching to {mode} mode.
{if yolo: "Visual checkpoints will be auto-approved. Decision and external action checkpoints still pause."}

Apply? [y/N]
```

If confirmed, update `mode` in config.json.

## Apply Changes

After any change:

1. Read current config.json
2. Apply the change
3. Write updated config.json
4. Confirm:

```
Settings updated:
  {field}: {old} → {new}
```
