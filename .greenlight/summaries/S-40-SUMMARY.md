# S-40: Roadmap and Changelog Commands

## What Changed
Users can now run `greenlight roadmap` and `greenlight changelog` from the terminal to view project documentation without needing Claude. The roadmap command reads and prints `.greenlight/ROADMAP.md` verbatim. The changelog command reads all `.md` files from `.greenlight/summaries/`, sorts them chronologically by filename, and prints each entry separated by `---`.

## Contracts Satisfied
- C-101: RunRoadmap
- C-102: RunChangelog

## Test Coverage
- 25 integration tests (9 roadmap + 16 changelog) using temp directories and t.Chdir for filesystem isolation
- Roadmap: verbatim output, multi-line preservation, missing directory/file error messages, writer contract
- Changelog: sorting order, separator placement (between entries, not after last), empty/missing states, writer contract

## Files
- internal/cmd/roadmap.go (modified — replaced stub with full implementation)
- internal/cmd/roadmap_test.go (new — 9 tests)
- internal/cmd/changelog.go (modified — replaced stub with full implementation)
- internal/cmd/changelog_test.go (new — 16 tests)

## Decisions
- Roadmap prints file content verbatim with no post-processing
- Missing summaries directory is not an error (exit 0 with "No changelog entries yet.")
- Individual file read errors in changelog are skipped with a warning, remaining entries still print
- Separator `---` appears between entries only, never after the final entry
