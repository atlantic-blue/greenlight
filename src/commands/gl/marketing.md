---
name: gl:marketing
description: Research-driven marketing planning — competitor analysis, channel strategy, revenue milestones, and task tracking
argument-hint: "<init|research|plan|status|ask|refresh> [args]"
allowed-tools: [Read, Write, Bash, Glob, Grep, Task, AskUserQuestion, WebSearch, WebFetch]
---

# Greenlight: Marketing

Research-driven marketing planning. Everything is researched, nothing is invented.

This command answers: "how do we reach customers and generate revenue?" It is not about CI, deployment, code quality, or software readiness — those are `/gl:ship`, `/gl:slice`, and `/gl:status`.

**Read first:**
- `CLAUDE.md` — engineering standards
- `.greenlight/config.json` — project config
- `.greenlight/DESIGN.md` — product context
- `.greenlight/MARKETING.md` — marketing state (if exists)
- `.greenlight/MARKETING_RESEARCH.md` — research findings (if exists)

## Prerequisites

```bash
cat .greenlight/config.json 2>/dev/null
```

**If no config exists:**
```
No project found. Run /gl:init first to set up the project.
```
Stop here.

## Subcommand Routing

Parse the argument to determine which subcommand to run:

- `/gl:marketing init` → [Init](#init)
- `/gl:marketing research` → [Research](#research)
- `/gl:marketing plan` → [Plan](#plan)
- `/gl:marketing status` → [Status](#status)
- `/gl:marketing ask "<question>"` → [Ask](#ask)
- `/gl:marketing refresh` → [Refresh](#refresh)
- `/gl:marketing` (no args) → Show subcommand help:

```
/gl:marketing — Research-driven marketing planning

Subcommands:
  init                Interview about commercial context → MARKETING.md
  research            Deep market research → MARKETING_RESEARCH.md
  plan                Prioritised 30-day task plan from research findings
  status              Progress against milestones and revenue target
  ask "<question>"    Ask a grounded question with full marketing context
  refresh             Re-run research to update stale data

Start with: /gl:marketing init
```

---

## Init

Interview the user about the product's commercial context. Ask only what is needed to produce a useful marketing plan.

### Interview Topics

Cover these areas conversationally — not as a questionnaire:

1. **Revenue target** — monthly net amount after taxes and platform fees, currency
2. **Target audience** — who has the problem, how they currently solve it
3. **Competitor landscape** — what exists, what it costs, why it falls short
4. **Positioning** — what makes this product the right choice
5. **Channels already tried** — and their results so far
6. **Budget** — available for paid acquisition per month
7. **Team** — size and hours available per week for marketing

### Write MARKETING.md

On completion, write `.greenlight/MARKETING.md`:

```markdown
# Marketing

## Revenue Target
- monthly_net: {amount}
- currency: {GBP|USD|EUR}
- current_monthly_net: 0
- last_updated: {YYYY-MM-DD}

## Product
- name: {from DESIGN.md or user}
- listing_url: {app store or product URL}
- price: {amount}
- pricing_model: {one-time | subscription | freemium}
- current_review_count: 0
- current_rating: 0

## Team
- hours_per_week_available: {N}
- paid_acquisition_budget_per_month: {amount}

## Audience
- primary: {who}
- problem_they_have: {what}
- how_they_solve_it_today: {what}
- language_they_use: {exact phrases from forums/reviews if known}

## Positioning
- statement: {one sentence}
- key_differentiator: {what}

## Competitors
- name: {competitor 1}
  url: {URL}
  pricing: {pricing}
  weaknesses: {known weaknesses}
  last_researched: {YYYY-MM-DD}

## Channels
- name: {channel} status: {active|planned|paused}
  notes: {details}
  result_so_far: {metrics if any}

## Milestones
(populated by /gl:marketing plan)

## Tasks
(populated by /gl:marketing plan and /gl:marketing ask)

## Metrics
- app_store_reviews: 0
- monthly_downloads: 0
- monthly_revenue: 0
- conversion_rate: 0
- cpa: 0
- notes:
```

**If MARKETING.md already exists:** Read it first, ask the user what changed, update only the fields that changed.

### Trigger Research

After writing MARKETING.md, immediately trigger the research phase to validate and enrich what the user said with real market data:

```
MARKETING.md written.

Starting research phase to validate assumptions with real market data...
```

Proceed to [Research](#research).

---

## Research

The most important phase. Must run before any plan is produced.

### Prerequisites

```bash
cat .greenlight/MARKETING.md 2>/dev/null
```

**If MARKETING.md does not exist:**
```
No marketing context found. Run /gl:marketing init first.
```
Stop here.

### Gather Context

```bash
cat .greenlight/config.json
cat .greenlight/DESIGN.md 2>/dev/null
cat .greenlight/MARKETING.md
```

### Model Resolution

Resolve the marketing-researcher agent model from `.greenlight/config.json`:

1. Check `model_overrides["marketing_researcher"]` — if set, use it
2. Else check `profiles[model_profile]["marketing_researcher"]` — use profile default
3. Else fall back to `opus`

The researcher defaults to opus because research quality directly determines plan quality. Bad research produces bad plans.

### Spawn Marketing Researcher

```
Task(prompt="
Read agents/marketing-researcher.md
Read CLAUDE.md

<product_context>
{product name, value prop, stack, users from DESIGN.md}
</product_context>

<marketing_context>
{full contents of MARKETING.md — revenue target, audience, positioning, competitors, channels}
</marketing_context>

<existing_research>
{if MARKETING_RESEARCH.md exists: full contents. Otherwise: 'No existing research'}
</existing_research>

Conduct deep market research. For each competitor in MARKETING.md, research:
1. Current pricing (fetch the actual pricing page)
2. Feature set (free tier vs paid)
3. App Store or product listing (rating, review count, sentiment)
4. Estimated revenue if publicly available
5. Weaknesses from user reviews (search for negative reviews explicitly)
6. Distribution channels they use

Also research:
7. Market size and demand validation (search volume, community discussions, demand trend)
8. Pricing benchmarks for comparable products in this category
9. Channel effectiveness (case studies, CPA benchmarks, organic traction examples)
10. Realistic revenue benchmarks for solo-founder apps in this category

Write findings to .greenlight/MARKETING_RESEARCH.md following the schema in your agent prompt.

Every claim must have a source URL. Every estimate must be labelled as such.
Flag contradictions between user assumptions and research findings.
", subagent_type="general-purpose", model="{resolved_model.marketing_researcher}", description="Market research")
```

### After Research

Verify output exists:

```bash
cat .greenlight/MARKETING_RESEARCH.md 2>/dev/null
```

**If MARKETING_RESEARCH.md does not exist or is empty:** Report failure. Research agent did not produce output.

**If output exists:** Report summary:

```
Research complete.

Key findings:
- Competitors analysed: {N}
- Pricing range: {min}-{max}
- Market demand: {growing|flat|declining}
- Assumptions validated: {N}
- Assumptions contradicted: {N}
- Data gaps: {N}

Full report: .greenlight/MARKETING_RESEARCH.md

Next: /gl:marketing plan to generate a prioritised task plan from these findings.
```

---

## Plan

Read `.greenlight/MARKETING.md` and `.greenlight/MARKETING_RESEARCH.md`.

### Prerequisites

```bash
cat .greenlight/MARKETING.md 2>/dev/null
cat .greenlight/MARKETING_RESEARCH.md 2>/dev/null
```

**If MARKETING_RESEARCH.md does not exist:**
```
No research found. Running /gl:marketing research first — plans without research produce generic advice.
```
Run [Research](#research) first. Do not produce a plan without research.

### Model Resolution

Resolve the marketing agent model from `.greenlight/config.json`:

1. Check `model_overrides["marketing"]` — if set, use it
2. Else check `profiles[model_profile]["marketing"]` — use profile default
3. Else fall back to `opus`

### Spawn Marketing Agent

```
Task(prompt="
Read agents/marketing.md
Read CLAUDE.md

<product_context>
{product name, value prop from DESIGN.md}
</product_context>

<marketing_context>
{full contents of MARKETING.md}
</marketing_context>

<research>
{full contents of MARKETING_RESEARCH.md}
</research>

Produce a prioritised 30-day task plan grounded entirely in the research findings.

The plan must:
- Contain specific, actionable tasks — no generic advice
- Cite the research finding that justifies each task
- Tie every task to a channel and a specific revenue impact estimate
- Be ordered by highest expected return on time invested
- Include success criteria for each task that are measurable
- Be realistic for the team size and hours stated in MARKETING.md
- Flag dependencies between tasks

Structure as milestones:
1. Foundation — what must be true before acquisition starts
2. First traction — first meaningful revenue signal
3. Repeatable channel — one channel proven and scalable
4. Revenue target — monthly net target achieved consistently

For each milestone:
- What success looks like (specific, measurable)
- Tasks to get there (with research citations)
- Realistic timeframe based on research benchmarks
- What failure looks like and how to diagnose it early

Write the milestones and tasks to .greenlight/MARKETING.md under the Milestones and Tasks sections.
", subagent_type="general-purpose", model="{resolved_model.marketing}", description="Marketing plan")
```

### GitHub Projects Integration

After the marketing agent completes:

```bash
cat .greenlight/config.json | grep github_project_url
```

**If `github_project_url` is configured:**

Create issues in the linked GitHub Project for each task:

```bash
# For each task in the plan:
# 1. Check for existing issue with same title
gh issue list --label marketing --search "{task_title}" --json title --jq '.[].title' 2>/dev/null

# 2. If no duplicate, create issue
gh issue create \
  --title "{task_title}" \
  --label "marketing" \
  --body "{task_description with research justification, success criteria, estimated time, budget required}"
```

- Label every marketing issue with `marketing`
- Never create duplicate issues — check for existing issues with the same title
- Never touch issues without the `marketing` label

**If `github_project_url` is not configured:**
```
Tasks written to MARKETING.md.

Tip: Configure github_project_url in .greenlight/config.json to sync tasks
to a GitHub Project for tracking.
```

### Report

```
Plan complete.

Milestones:
  1. Foundation: {summary} ({timeframe})
  2. First traction: {summary} ({timeframe})
  3. Repeatable channel: {summary} ({timeframe})
  4. Revenue target: {summary} ({timeframe})

Tasks created: {N}
{if github: "GitHub issues created: {N}"}

Next: /gl:marketing status to track progress
      /gl:marketing ask "what should I work on this week?"
```

---

## Status

Read current marketing task status and display progress against milestones.

### Gather Data

```bash
cat .greenlight/MARKETING.md 2>/dev/null
cat .greenlight/MARKETING_RESEARCH.md 2>/dev/null
cat .greenlight/config.json
```

**If MARKETING.md does not exist:**
```
No marketing context found. Run /gl:marketing init first.
```
Stop here.

### GitHub Projects Read (if configured)

```bash
# If github_project_url is set:
gh issue list --label marketing --json number,title,state,assignees --jq '.[] | "\(.number)\t\(.state)\t\(.title)\t\(.assignees | map(.login) | join(","))"' 2>/dev/null
```

### Display

```
┌───────────────────────────────────────────────────────────┐
│  MARKETING STATUS                                          │
├───────────────────────────────────────────────────────────┤
│                                                            │
│  Revenue: {current}/{target} {currency}/month              │
│  Milestone: {N}. {name} [{progress}]                       │
│                                                            │
│  Tasks:                                                    │
│    Done:       {N}                                         │
│    In progress: {N}                                        │
│    Remaining:  {N}                                         │
│    Blocked:    {N}                                         │
│                                                            │
│  Next 3:                                                   │
│    1. {task} ({channel}, {est. time})                      │
│    2. {task} ({channel}, {est. time})                      │
│    3. {task} ({channel}, {est. time})                      │
│                                                            │
│  Research: {age} days old {fresh|stale — refresh due}      │
│  Last updated: {date}                                      │
│                                                            │
│  /gl:marketing ask "what should I work on this week?"      │
│  /gl:marketing refresh  (if research > 30 days old)        │
└───────────────────────────────────────────────────────────┘
```

### Intelligence

| Situation | Recommendation |
|-----------|---------------|
| No plan exists | `/gl:marketing plan` |
| Research > 30 days old | `/gl:marketing refresh` — competitor pricing may have changed |
| Tasks overdue or blocked | List them with suggested unblock action |
| Revenue at 0, plan exists | Focus on foundation milestone tasks |
| Revenue growing | Show growth rate and projected time to target |
| All milestone 1 tasks done | "Foundation complete. Move to milestone 2: first traction" |

---

## Ask

Answer a specific question with full marketing context.

### Usage

`/gl:marketing ask "what should natalia work on this week?"`

### Prerequisites

```bash
cat .greenlight/MARKETING.md 2>/dev/null
cat .greenlight/MARKETING_RESEARCH.md 2>/dev/null
cat .greenlight/DESIGN.md 2>/dev/null
cat .greenlight/config.json
```

**If MARKETING.md does not exist:**
```
No marketing context found. Run /gl:marketing init first.
```
Stop here.

### Model Resolution

Resolve the marketing agent model (same as plan).

### Spawn Marketing Agent

```
Task(prompt="
Read agents/marketing.md
Read CLAUDE.md

<product_context>
{product name, value prop from DESIGN.md}
</product_context>

<marketing_context>
{full contents of MARKETING.md}
</marketing_context>

<research>
{full contents of MARKETING_RESEARCH.md, or 'No research yet — flag this in your answer'}
</research>

<question>
{user's question}
</question>

Answer the question grounded in the research findings.

If the question requires information not in the research, say so explicitly
and offer to research it. Do not invent data to fill gaps.

If the answer generates actionable tasks and the user confirms them,
write them to .greenlight/MARKETING.md under the Tasks section.
", subagent_type="general-purpose", model="{resolved_model.marketing}", description="Marketing: {short question}")
```

### Task Creation from Ask

If the marketing agent produces tasks and the user confirms:

1. Append tasks to MARKETING.md
2. If `github_project_url` is configured, create GitHub issues (same dedup logic as plan)

---

## Refresh

Re-run the research phase to update MARKETING_RESEARCH.md with current data.

### Prerequisites

```bash
cat .greenlight/MARKETING.md 2>/dev/null
cat .greenlight/MARKETING_RESEARCH.md 2>/dev/null
```

**If MARKETING.md does not exist:**
```
No marketing context found. Run /gl:marketing init first.
```
Stop here.

### Check Staleness

Parse `Last updated:` from MARKETING_RESEARCH.md header.

```
Research last updated: {date} ({N} days ago)
{if N <= 30: "Research is still fresh. Refresh anyway? [y/N]"}
{if N > 30: "Research is stale. Refreshing..."}
```

### Run Research

Execute the same research flow as [Research](#research), passing the existing MARKETING_RESEARCH.md as context so the researcher can identify what changed.

### After Refresh

Compare old and new findings. Flag changes:

```
Research refreshed.

Changes detected:
- {competitor} pricing changed: {old} → {new}
- New competitor found: {name}
- Market demand trend: {old} → {new}

{if changes affect plan: "Current plan may need updating. Run /gl:marketing plan to regenerate."}
{if no changes: "No significant changes. Current plan remains valid."}
```

---

## GitHub Projects Integration

This integration is optional. It activates only when `github_project_url` is present in `.greenlight/config.json`.

### Config Schema Addition

The following fields are recognised in `.greenlight/config.json`:

```json
{
  "github_project_url": "https://github.com/orgs/{org}/projects/{N}",
  "marketing_revenue_target": 1000,
  "marketing_revenue_currency": "GBP"
}
```

### When Creating Issues

- Use `gh` CLI (already available in Greenlight environments)
- Label every marketing issue with `marketing`
- Title = task title from the plan
- Body = task description, research justification, success criteria, estimated time, budget required
- Never create duplicate issues — check for existing issues with the same title before creating
- Never touch issues without the `marketing` label

### When Reading Issues

- Fetch all issues labelled `marketing`
- Return them grouped by state (open/closed)
- Include issue number, title, and assignee if set

### Authentication

- Uses the existing `gh` CLI authentication (user must be logged in via `gh auth login`)
- If `gh` is not authenticated, report clearly:
  ```
  GitHub CLI not authenticated. Run: gh auth login
  ```

### Error Handling

- If the GitHub API returns an error, fall back to MARKETING.md gracefully
- Inform the user clearly what failed and how to fix it
- Never silently swallow errors

---

## Error Handling

| Error | When | Response |
|-------|------|----------|
| NoConfig | config.json missing | "Run /gl:init first" |
| NoMarketing | MARKETING.md missing | "Run /gl:marketing init first" |
| NoResearch | Plan requested without research | Run research first, refuse to plan without it |
| ResearchAgentFailure | Research agent fails | Report error, suggest retrying |
| PlanAgentFailure | Plan agent fails | Report error, suggest retrying |
| GitHubAuthFailure | gh not authenticated | "Run: gh auth login" |
| GitHubAPIError | Issue creation fails | Fall back to MARKETING.md, report error |
| StaleResearch | Research > 30 days old | Warn on status, suggest /gl:marketing refresh |
