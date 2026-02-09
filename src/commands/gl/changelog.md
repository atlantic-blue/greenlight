---
name: gl:changelog
description: Display project changelog from summaries directory, with milestone and date filtering
argument-hint: "[milestone <name>|since <YYYY-MM-DD>]"
allowed-tools: [Read, Glob, Bash]
---

# /gl:changelog

Displays a chronological changelog of completed work by reading summary files from the `.greenlight/summaries/` directory.

**This is a read-only command. No files are written.**

## Core Behaviour

The changelog command scans the `summaries/` directory for summary files and displays them in reverse chronological order (newest first).

### Summary File Naming Conventions

- Slices: `{slice-id}-SUMMARY.md`
- Wraps: `{boundary}-wrap-SUMMARY.md`
- Quick fixes: `quick-{timestamp}-SUMMARY.md`

### Parsing Summary Files

Each SUMMARY.md file is parsed to extract:
- **Date** — completion date
- **Type** — one of: `slice`, `wrap`, or `quick`
- **Name** — the slice/wrap/quick identifier
- **One-line summary** — brief description of the work
- **Test count** — number of tests added or passing

### Output Format

```
CHANGELOG -- {Project Name}

{date} {type}:{name} — {one-line summary} ({N} tests)
{date} {type}:{name} — {one-line summary} ({N} tests)
...

{N} entries ({N} slices, {N} wraps, {N} quick)
```

### Project Name Resolution

The project name is read from `config.json`. If the config file is not found or cannot be read, display:

```
CHANGELOG -- Unknown Project
```

Do not fail the command if `config.json` is missing — this is a graceful degradation.

### Error Handling

**No summaries directory:**
- If `.greenlight/summaries/` does not exist, display: `No summaries found`

**Empty summaries directory:**
- If the directory exists but contains no summary files, display: `No summaries found yet — empty summaries directory`

**Malformed summary files:**
- If a SUMMARY.md file cannot be parsed or is malformed, skip it and continue processing other files
- Log: `Could not parse {filename} — malformed summary, skipping`

## Subcommand: milestone

Filter the changelog to show only entries associated with a specific milestone.

### Usage

```
/gl:changelog milestone <milestone-name>
```

### Behaviour

1. Read `GRAPH.json` to find all slices associated with the specified milestone name
2. Filter changelog entries to include only those slices whose name matches the milestone
3. Display filtered results with header: `CHANGELOG -- {Project Name} (milestone: {milestone-name})`
4. Show entry count summary at bottom

**This filter operation is read-only. No files are written.**

### Error Handling

**NoGraphJson:**
- If `GRAPH.json` is missing or cannot be read, display: `Cannot filter by milestone — GRAPH.json not found`

**UnknownMilestone:**
- If the milestone name is not found in GRAPH.json, display: `Milestone '{name}' not found`

**No matching entries:**
- If no changelog entries match the milestone, display: `No changelog entries found for milestone '{name}'`

## Subcommand: since

Filter the changelog to show only entries since a specified date.

### Usage

```
/gl:changelog since <YYYY-MM-DD>
```

### Behaviour

1. Parse the date argument in ISO 8601 format: `YYYY-MM-DD`
2. Filter changelog entries to include only those with a date >= the specified date (inclusive)
3. Display filtered results with header: `CHANGELOG -- {Project Name} (since {date})`
4. Show entry count summary at bottom
5. Most recent entries appear first (chronological sorting, newest first)

**This filter operation is read-only. No files are written.**

### Date Filtering Rules

- Date format must be `YYYY-MM-DD`
- Filter is inclusive: entries matching the exact date are included
- Entries are sorted chronologically after filtering

### Error Handling

**InvalidDate:**
- If the date cannot be parsed or is not in `YYYY-MM-DD` format, display: `Invalid date format — use YYYY-MM-DD`

**No matching entries:**
- If no entries exist on or after the specified date, display: `No entries found since {date}`

## Implementation Notes

- Always scan the entire `summaries/` directory
- Parse each SUMMARY.md file to extract metadata
- Sort entries by date (most recent first) before displaying
- Count entries by type (slice, wrap, quick) for the summary line
- Handle missing or malformed files gracefully — do not fail the entire command
- All filtering operations are read-only — no files are modified or written
