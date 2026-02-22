# S-33: Documentation Updates

## What Changed
Three documentation files updated with state format awareness. CLAUDE.md now has a hard rule requiring all agents to check `.greenlight/slices/` before reading STATE.md. The state template documents both formats side by side with migration instructions. The checkpoint protocol references slice files for state context.

## User Impact
- All agents are now aware of the file-per-slice state format
- Documentation guides new projects toward file-per-slice format
- Migration path clearly documented via /gl:migrate-state
- Checkpoint save/restore works correctly with both state formats

## Contracts Satisfied
- C-87: CLAUDEmdStateFormatRule (auto)
- C-88: StateTemplateDocUpdate (auto)
- C-89: CheckpointProtocolStateUpdate (auto)

## Tests
- 22 passing (0 security)
- Coverage: all contract elements fully covered

## Files
- `src/CLAUDE.md` (modified — +5 lines)
- `src/templates/state.md` (modified — +33 lines)
- `src/references/checkpoint-protocol.md` (modified — +10 lines)
- `internal/installer/documentation_state_update_test.go` (new — 22 tests)

## Architecture
No architecture changes. Documentation-only updates.
