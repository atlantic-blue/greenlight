# S-36: State Reader

## What Changed
The CLI can now read project state from the filesystem. A new `internal/state` package provides four functions: read all slice frontmatter into structured data, read the dependency graph from GRAPH.json, compute which slices are ready to build (pending with all dependencies complete), and detect whether the CLI is running inside Claude or from the user's terminal.

## Contracts Satisfied
- C-93: StateReadSlices
- C-94: StateReadGraph
- C-95: StateFindReadySlices
- C-96: StateDetectContext

## Test Coverage
- 41 integration tests covering all 4 contracts
- ReadSlices: 14 tests (file reading, deps parsing, defaults, error cases)
- ReadGraph: 8 tests (JSON parsing, field handling, error cases)
- FindReadySlices: 14 tests (dependency resolution, sorting, edge cases)
- DetectContext: 5 tests (environment variable detection)

## Files
- internal/state/state.go (new — 285 lines)
- internal/state/state_test.go (new — 41 tests)

## Decisions
- Used os.ReadDir + os.ReadFile for filesystem access (stdlib only)
- json.RawMessage intermediate struct to distinguish missing "slices" field from invalid JSON
- Status lookup map for O(1) dependency checking in FindReadySlices
- Slices not in graph treated as having no dependencies (always ready if pending)
