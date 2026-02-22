package installer_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// C-85 / C-86 Tests: MigrateStateCommand
// These tests verify that src/commands/gl/migrate-state.md exists and contains
// all required instructions as defined in contracts C-85 and C-86.

func readMigrateStateMd(t *testing.T) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/migrate-state.md"))
	if err != nil {
		t.Fatalf("failed to read src/commands/gl/migrate-state.md: %v", err)
	}
	return string(content)
}

// --- Command metadata ---

func TestMigrateStateMd_C85_FileExists(t *testing.T) {
	path := filepath.Join(projectRoot(), "src/commands/gl/migrate-state.md")
	_, err := os.Stat(path)
	if err != nil {
		t.Fatalf("migrate-state.md does not exist at %s: %v", path, err)
	}
}

func TestMigrateStateMd_C85_ContainsCommandName(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasName := strings.Contains(doc, "gl:migrate-state") || strings.Contains(doc, "migrate-state")
	if !hasName {
		t.Error("migrate-state.md missing command name 'gl:migrate-state' or 'migrate-state'")
	}
}

func TestMigrateStateMd_C85_ContainsDescription(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasDescription := strings.Contains(doc, "description") || strings.Contains(doc, "Description") ||
		strings.Contains(doc, "Converts") || strings.Contains(doc, "converts") ||
		strings.Contains(doc, "migrate") || strings.Contains(doc, "Migrate")
	if !hasDescription {
		t.Error("migrate-state.md missing a description of the command's purpose")
	}
}

// --- Pre-checks: step 1 (STATE.md must exist) ---

func TestMigrateStateMd_C85_VerifiesStateMdExists(t *testing.T) {
	doc := readMigrateStateMd(t)
	if !strings.Contains(doc, "STATE.md") {
		t.Error("migrate-state.md missing reference to STATE.md — must verify its existence")
	}
}

func TestMigrateStateMd_C85_ErrorMessageForMissingStateMd(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasNoStateMdError := strings.Contains(doc, "No STATE.md found") ||
		strings.Contains(doc, "Nothing to migrate")
	if !hasNoStateMdError {
		t.Error("migrate-state.md missing error message for absent STATE.md ('No STATE.md found' or 'Nothing to migrate')")
	}
}

// --- Pre-checks: step 2 (slices/ must not exist) ---

func TestMigrateStateMd_C85_VerifiesSlicesDirAbsent(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasSlicesCheck := strings.Contains(doc, "slices/") || strings.Contains(doc, ".greenlight/slices")
	if !hasSlicesCheck {
		t.Error("migrate-state.md missing check that .greenlight/slices/ does not already exist")
	}
}

func TestMigrateStateMd_C85_ErrorMessageForAlreadyMigrated(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasAlreadyMigratedError := strings.Contains(doc, "Already using file-per-slice") ||
		strings.Contains(doc, "Nothing to migrate") ||
		strings.Contains(doc, "AlreadyMigrated")
	if !hasAlreadyMigratedError {
		t.Error("migrate-state.md missing error for already-migrated state ('Already using file-per-slice' or 'Nothing to migrate')")
	}
}

// --- Parsing: step 3 ---

func TestMigrateStateMd_C85_InstructsParseStateMd(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasParseInstruction := strings.Contains(doc, "Parse STATE.md") ||
		strings.Contains(doc, "parse STATE.md") ||
		strings.Contains(doc, "Parse the STATE.md") ||
		strings.Contains(doc, "read STATE.md") ||
		strings.Contains(doc, "Read STATE.md")
	if !hasParseInstruction {
		t.Error("migrate-state.md missing instruction to parse STATE.md")
	}
}

func TestMigrateStateMd_C85_ExtractsSliceTable(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasTableExtraction := strings.Contains(doc, "slice table") ||
		strings.Contains(doc, "table rows") ||
		strings.Contains(doc, "Extract slice") ||
		strings.Contains(doc, "extract slice")
	if !hasTableExtraction {
		t.Error("migrate-state.md missing instruction to extract slice table rows from STATE.md")
	}
}

func TestMigrateStateMd_C85_TableRowsIncludeIdNameStatus(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasID := strings.Contains(doc, "ID") || strings.Contains(doc, " id")
	hasName := strings.Contains(doc, "Name") || strings.Contains(doc, "name")
	hasStatus := strings.Contains(doc, "Status") || strings.Contains(doc, "status")
	if !hasID || !hasName || !hasStatus {
		t.Error("migrate-state.md must reference ID, Name, and Status fields from slice table rows")
	}
}

func TestMigrateStateMd_C85_ExtractsCurrentSection(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasCurrentSection := strings.Contains(doc, "Current section") ||
		strings.Contains(doc, "Current Section") ||
		strings.Contains(doc, "Extract Current") ||
		strings.Contains(doc, "extract Current") ||
		strings.Contains(doc, "## Current")
	if !hasCurrentSection {
		t.Error("migrate-state.md missing instruction to extract the Current section from STATE.md")
	}
}

func TestMigrateStateMd_C85_ExtractsOverviewSection(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasOverview := strings.Contains(doc, "Overview") || strings.Contains(doc, "overview")
	if !hasOverview {
		t.Error("migrate-state.md missing instruction to extract Overview section (value prop, stack, mode)")
	}
}

func TestMigrateStateMd_C85_ExtractsSessionSection(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasSession := strings.Contains(doc, "Session") || strings.Contains(doc, "session")
	if !hasSession {
		t.Error("migrate-state.md missing instruction to extract Session section from STATE.md")
	}
}

func TestMigrateStateMd_C85_ExtractsBlockersSection(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasBlockers := strings.Contains(doc, "Blockers") || strings.Contains(doc, "blockers")
	if !hasBlockers {
		t.Error("migrate-state.md missing instruction to extract Blockers section from STATE.md")
	}
}

func TestMigrateStateMd_C85_ExtractsDecisionsSection(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasDecisions := strings.Contains(doc, "Decisions") || strings.Contains(doc, "decisions")
	if !hasDecisions {
		t.Error("migrate-state.md missing instruction to extract Decisions section from STATE.md")
	}
}

// --- Directory and file creation: steps 4-5 ---

func TestMigrateStateMd_C85_CreatesSlicesDirectory(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasCreateDir := strings.Contains(doc, "Create .greenlight/slices") ||
		strings.Contains(doc, "create .greenlight/slices") ||
		strings.Contains(doc, "mkdir") ||
		strings.Contains(doc, "MkdirAll") ||
		strings.Contains(doc, "slices/ directory")
	if !hasCreateDir {
		t.Error("migrate-state.md missing instruction to create .greenlight/slices/ directory")
	}
}

func TestMigrateStateMd_C85_DirectoryPermissions(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasPermissions := strings.Contains(doc, "0o755") || strings.Contains(doc, "0755") ||
		strings.Contains(doc, "755")
	if !hasPermissions {
		t.Error("migrate-state.md missing directory permissions (0o755 or 755)")
	}
}

func TestMigrateStateMd_C85_CreatesIndividualSliceFiles(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasSliceFiles := strings.Contains(doc, "{id}.md") ||
		strings.Contains(doc, "slices/{id}") ||
		strings.Contains(doc, "slices/S-") ||
		strings.Contains(doc, "individual slice file") ||
		strings.Contains(doc, "per-slice file")
	if !hasSliceFiles {
		t.Error("migrate-state.md missing instruction to create individual .greenlight/slices/{id}.md files")
	}
}

func TestMigrateStateMd_C85_ValidatesSliceIdFormat(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasValidation := strings.Contains(doc, "S-{digits}") ||
		strings.Contains(doc, "S-XX") ||
		strings.Contains(doc, "Validate slice ID") ||
		strings.Contains(doc, "validate slice ID") ||
		strings.Contains(doc, "slice ID format") ||
		strings.Contains(doc, "ID format")
	if !hasValidation {
		t.Error("migrate-state.md missing slice ID format validation (S-{digits} or equivalent)")
	}
}

func TestMigrateStateMd_C85_SliceFilesHaveFrontmatterIdStatusTests(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasFrontmatter := strings.Contains(doc, "frontmatter") || strings.Contains(doc, "Frontmatter")
	hasID := strings.Contains(doc, "id:") || strings.Contains(doc, " id")
	hasStatus := strings.Contains(doc, "status:") || strings.Contains(doc, "status")
	hasTests := strings.Contains(doc, "tests:") || strings.Contains(doc, "tests")
	if !hasFrontmatter || !hasID || !hasStatus || !hasTests {
		t.Error("migrate-state.md missing frontmatter with id, status, and tests fields for slice files")
	}
}

func TestMigrateStateMd_C85_SliceFilesHaveFrontmatterStepMilestoneSession(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasStep := strings.Contains(doc, "step:") || strings.Contains(doc, "step")
	hasMilestone := strings.Contains(doc, "milestone:") || strings.Contains(doc, "milestone")
	hasSession := strings.Contains(doc, "session:") || strings.Contains(doc, "session")
	if !hasStep || !hasMilestone || !hasSession {
		t.Error("migrate-state.md missing frontmatter fields step, milestone, and session for slice files")
	}
}

func TestMigrateStateMd_C85_UsesWriteToTempThenRename(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasAtomic := strings.Contains(doc, "write-to-temp-then-rename") ||
		strings.Contains(doc, "temp-then-rename") ||
		strings.Contains(doc, "NFR-4") ||
		strings.Contains(doc, "atomic write") ||
		strings.Contains(doc, "atomic") ||
		strings.Contains(doc, "temp file") ||
		strings.Contains(doc, "os.Rename")
	if !hasAtomic {
		t.Error("migrate-state.md missing write-to-temp-then-rename / NFR-4 / atomic write instruction for slice files")
	}
}

func TestMigrateStateMd_C85_SliceFilePermissions(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasFilePermissions := strings.Contains(doc, "0o644") || strings.Contains(doc, "0644") ||
		strings.Contains(doc, "644")
	if !hasFilePermissions {
		t.Error("migrate-state.md missing file permissions (0o644 or 644) for slice files")
	}
}

// --- Project state: step 6 ---

func TestMigrateStateMd_C85_CreatesProjectStateJson(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasProjectState := strings.Contains(doc, "project-state.json") ||
		strings.Contains(doc, "project_state.json")
	if !hasProjectState {
		t.Error("migrate-state.md missing instruction to create .greenlight/project-state.json")
	}
}

func TestMigrateStateMd_C85_ProjectStateJsonContainsNonSliceSections(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasOverview := strings.Contains(doc, "overview") || strings.Contains(doc, "Overview")
	hasSession := strings.Contains(doc, "session") || strings.Contains(doc, "Session")
	hasBlockers := strings.Contains(doc, "blockers") || strings.Contains(doc, "Blockers")
	if !hasOverview || !hasSession || !hasBlockers {
		t.Error("migrate-state.md missing instruction that project-state.json contains overview, session, and blockers")
	}
}

// --- Backup: step 7 ---

func TestMigrateStateMd_C85_RenamesStateMdToBackup(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasBackup := strings.Contains(doc, "STATE.md.backup") || strings.Contains(doc, "backup")
	if !hasBackup {
		t.Error("migrate-state.md missing instruction to rename STATE.md to STATE.md.backup")
	}
}

func TestMigrateStateMd_C85_OriginalPreservedAsBackup(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasPreserved := strings.Contains(doc, "preserved") || strings.Contains(doc, "backup") ||
		strings.Contains(doc, "original")
	if !hasPreserved {
		t.Error("migrate-state.md missing statement that the original STATE.md is preserved as a backup")
	}
}

// --- Regeneration: step 8 ---

func TestMigrateStateMd_C85_GeneratesNewStateMd(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasGenerate := strings.Contains(doc, "Generate new STATE.md") ||
		strings.Contains(doc, "generate new STATE.md") ||
		strings.Contains(doc, "generated STATE.md") ||
		strings.Contains(doc, "new STATE.md")
	if !hasGenerate {
		t.Error("migrate-state.md missing instruction to generate a new STATE.md in generated format")
	}
}

func TestMigrateStateMd_C85_GeneratedStateMdHasHeader(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasGeneratedHeader := strings.Contains(doc, "GENERATED") ||
		strings.Contains(doc, "generated") ||
		strings.Contains(doc, "header comment") ||
		strings.Contains(doc, "do not edit")
	if !hasGeneratedHeader {
		t.Error("migrate-state.md missing GENERATED header comment for the new STATE.md")
	}
}

// --- Report: step 9 ---

func TestMigrateStateMd_C85_ReportsSuccessToUser(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasReport := strings.Contains(doc, "Migrated") ||
		strings.Contains(doc, "migrated") ||
		strings.Contains(doc, "slices") ||
		strings.Contains(doc, "file-per-slice")
	if !hasReport {
		t.Error("migrate-state.md missing success report instruction ('Migrated N slices' or equivalent)")
	}
}

func TestMigrateStateMd_C85_ReportMentionsBackupLocation(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasBackupInReport := strings.Contains(doc, "STATE.md.backup") ||
		strings.Contains(doc, "Backup:")
	if !hasBackupInReport {
		t.Error("migrate-state.md missing backup location in success report")
	}
}

// --- Error handling ---

func TestMigrateStateMd_C85_HandlesNoStateMdError(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasNoStateMd := strings.Contains(doc, "NoStateMd") ||
		strings.Contains(doc, "No STATE.md") ||
		strings.Contains(doc, "STATE.md does not exist") ||
		strings.Contains(doc, "STATE.md not found")
	if !hasNoStateMd {
		t.Error("migrate-state.md missing NoStateMd error handling")
	}
}

func TestMigrateStateMd_C85_HandlesAlreadyMigratedError(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasAlreadyMigrated := strings.Contains(doc, "AlreadyMigrated") ||
		strings.Contains(doc, "Already using file-per-slice") ||
		strings.Contains(doc, "already migrated")
	if !hasAlreadyMigrated {
		t.Error("migrate-state.md missing AlreadyMigrated error handling")
	}
}

func TestMigrateStateMd_C85_HandlesParseFailureError(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasParseFailure := strings.Contains(doc, "ParseFailure") ||
		strings.Contains(doc, "parse failure") ||
		strings.Contains(doc, "cannot parse") ||
		strings.Contains(doc, "failed to parse")
	if !hasParseFailure {
		t.Error("migrate-state.md missing ParseFailure error handling")
	}
}

func TestMigrateStateMd_C85_HandlesInvalidSliceIdError(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasInvalidSliceId := strings.Contains(doc, "InvalidSliceId") ||
		strings.Contains(doc, "invalid slice ID") ||
		strings.Contains(doc, "invalid ID") ||
		strings.Contains(doc, "bad slice ID")
	if !hasInvalidSliceId {
		t.Error("migrate-state.md missing InvalidSliceId error handling")
	}
}

func TestMigrateStateMd_C85_InvalidSliceIdSkipsAndWarns(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasSkipWarn := strings.Contains(doc, "skip") || strings.Contains(doc, "Skip") ||
		strings.Contains(doc, "warn") || strings.Contains(doc, "Warn")
	if !hasSkipWarn {
		t.Error("migrate-state.md missing skip-and-warn behaviour for invalid slice IDs")
	}
}

func TestMigrateStateMd_C85_HandlesPartialWriteFailure(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasPartialWriteFailure := strings.Contains(doc, "PartialWriteFailure") ||
		strings.Contains(doc, "partial write") ||
		strings.Contains(doc, "partial failure") ||
		strings.Contains(doc, "abort")
	if !hasPartialWriteFailure {
		t.Error("migrate-state.md missing PartialWriteFailure error handling (abort, remove, restore)")
	}
}

func TestMigrateStateMd_C85_HandlesBackupRenameFailure(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasBackupRenameFailure := strings.Contains(doc, "BackupRenameFailure") ||
		strings.Contains(doc, "backup rename") ||
		strings.Contains(doc, "cannot rename") ||
		strings.Contains(doc, "rename failure")
	if !hasBackupRenameFailure {
		t.Error("migrate-state.md missing BackupRenameFailure error handling")
	}
}

// --- Invariants ---

func TestMigrateStateMd_C85_MigrationIsOneWay(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasOneWay := strings.Contains(doc, "one-way") || strings.Contains(doc, "One-way") ||
		strings.Contains(doc, "irreversible") || strings.Contains(doc, "one way")
	if !hasOneWay {
		t.Error("migrate-state.md missing one-way migration invariant")
	}
}

func TestMigrateStateMd_C85_MigrationIsAllOrNothing(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasAllOrNothing := strings.Contains(doc, "all-or-nothing") ||
		strings.Contains(doc, "All-or-nothing") ||
		strings.Contains(doc, "atomic") ||
		strings.Contains(doc, "all or nothing")
	if !hasAllOrNothing {
		t.Error("migrate-state.md missing all-or-nothing migration invariant")
	}
}

func TestMigrateStateMd_C85_ExplicitOnlyNoAutoMigration(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasExplicitOnly := strings.Contains(doc, "Explicit only") ||
		strings.Contains(doc, "explicit only") ||
		strings.Contains(doc, "no auto-migration") ||
		strings.Contains(doc, "D-32") ||
		strings.Contains(doc, "no automatic")
	if !hasExplicitOnly {
		t.Error("migrate-state.md missing explicit-only / no-auto-migration invariant (D-32)")
	}
}

func TestMigrateStateMd_C85_NoDualWrite(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasNoDualWrite := strings.Contains(doc, "No dual-write") ||
		strings.Contains(doc, "no dual-write") ||
		strings.Contains(doc, "dual-write") ||
		strings.Contains(doc, "D-38")
	if !hasNoDualWrite {
		t.Error("migrate-state.md missing no-dual-write invariant (D-38)")
	}
}

func TestMigrateStateMd_C85_PostMigrationDetectionWorks(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasDetection := strings.Contains(doc, "D-31") ||
		strings.Contains(doc, "detection") ||
		strings.Contains(doc, "file-per-slice automatically") ||
		strings.Contains(doc, "commands use file-per-slice")
	if !hasDetection {
		t.Error("migrate-state.md missing post-migration detection invariant (D-31)")
	}
}

// --- Security ---

func TestMigrateStateMd_C85_SliceIdValidationPreventsPathTraversal(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasPathTraversalPrevention := strings.Contains(doc, "path traversal") ||
		strings.Contains(doc, "Path traversal") ||
		strings.Contains(doc, "traversal") ||
		strings.Contains(doc, "sanitise") ||
		strings.Contains(doc, "sanitize") ||
		strings.Contains(doc, "validate") && strings.Contains(doc, "slice ID")
	if !hasPathTraversalPrevention {
		t.Error("migrate-state.md missing slice ID validation for path traversal prevention")
	}
}

func TestMigrateStateMd_C85_DoesNotModifySourceCode(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasNoSourceModify := strings.Contains(doc, "Does not modify") ||
		strings.Contains(doc, "does not modify") ||
		strings.Contains(doc, "No source") ||
		strings.Contains(doc, "no source") ||
		strings.Contains(doc, "source code") ||
		strings.Contains(doc, "test files")
	if !hasNoSourceModify {
		t.Error("migrate-state.md missing statement that it does not modify source code, test files, or config")
	}
}

// --- C-86 Tests: MigrateStateBackup ---
// Ordering, atomicity, cleanup, and crash safety.

func TestMigrateStateMd_C86_SliceFilesWrittenBeforeBackup(t *testing.T) {
	doc := readMigrateStateMd(t)
	// Check that slice files and project-state.json appear before the backup step
	slicePos := strings.Index(doc, "slices/")
	if slicePos == -1 {
		slicePos = strings.Index(doc, "slice file")
	}
	backupPos := strings.Index(doc, "STATE.md.backup")
	if backupPos == -1 {
		backupPos = strings.Index(doc, "backup")
	}
	if slicePos == -1 || backupPos == -1 {
		t.Error("migrate-state.md missing slice file creation or backup step references")
		return
	}
	if slicePos >= backupPos {
		t.Error("migrate-state.md slice file creation should be described before the backup step")
	}
}

func TestMigrateStateMd_C86_StateMdRenameHappensLast(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasRenameLastStatement := strings.Contains(doc, "rename happens last") ||
		strings.Contains(doc, "last step") ||
		strings.Contains(doc, "AFTER all files") ||
		strings.Contains(doc, "after all files") ||
		strings.Contains(doc, "Only AFTER") ||
		strings.Contains(doc, "only after")
	if !hasRenameLastStatement {
		t.Error("migrate-state.md missing statement that STATE.md rename happens last / after all files succeed")
	}
}

func TestMigrateStateMd_C86_StateMdUntouchedOnFailureBeforeRename(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasUntouched := strings.Contains(doc, "STATE.md untouched") ||
		strings.Contains(doc, "untouched") ||
		strings.Contains(doc, "STATE.md intact") ||
		strings.Contains(doc, "intact")
	if !hasUntouched {
		t.Error("migrate-state.md missing statement that STATE.md remains untouched if failure occurs before rename")
	}
}

func TestMigrateStateMd_C86_RemovesPartialSlicesDirOnFailure(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasRemove := strings.Contains(doc, "RemoveAll") ||
		strings.Contains(doc, "remove partial") ||
		strings.Contains(doc, "Remove partial") ||
		strings.Contains(doc, "remove slices/") ||
		strings.Contains(doc, "clean up") ||
		strings.Contains(doc, "cleanup")
	if !hasRemove {
		t.Error("migrate-state.md missing instruction to remove partial slices/ directory on failure")
	}
}

func TestMigrateStateMd_C86_CleanupFailureReportsManualInstructions(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasCleanupFailure := strings.Contains(doc, "CleanupFailure") ||
		strings.Contains(doc, "cleanup failure") ||
		strings.Contains(doc, "manual cleanup") ||
		strings.Contains(doc, "manual removal")
	if !hasCleanupFailure {
		t.Error("migrate-state.md missing CleanupFailure handling with manual cleanup instructions")
	}
}

func TestMigrateStateMd_C86_HandlesBackupAlreadyExists(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasBackupExists := strings.Contains(doc, "BackupExists") ||
		strings.Contains(doc, "backup already exists") ||
		strings.Contains(doc, "STATE.md.backup already") ||
		strings.Contains(doc, "already exists")
	if !hasBackupExists {
		t.Error("migrate-state.md missing BackupExists handling (timestamp suffix fallback)")
	}
}

func TestMigrateStateMd_C86_TimestampSuffixForMultipleBackups(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasTimestamp := strings.Contains(doc, "timestamp") ||
		strings.Contains(doc, "STATE.md.backup.{timestamp}") ||
		strings.Contains(doc, ".backup.")
	if !hasTimestamp {
		t.Error("migrate-state.md missing timestamped suffix for multiple backups (STATE.md.backup.{timestamp})")
	}
}

func TestMigrateStateMd_C86_MultipleBackupsPreserved(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasMultipleBackups := strings.Contains(doc, "Multiple backups") ||
		strings.Contains(doc, "multiple backups") ||
		strings.Contains(doc, "preserved") ||
		strings.Contains(doc, "timestamped")
	if !hasMultipleBackups {
		t.Error("migrate-state.md missing statement that multiple backups are preserved with timestamped suffixes")
	}
}

func TestMigrateStateMd_C86_WriteToTempThenRenameForSliceFiles(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasAtomicWrite := strings.Contains(doc, "write-to-temp-then-rename") ||
		strings.Contains(doc, "NFR-4") ||
		strings.Contains(doc, "temp-then-rename") ||
		strings.Contains(doc, "atomic") ||
		strings.Contains(doc, "os.Rename")
	if !hasAtomicWrite {
		t.Error("migrate-state.md missing write-to-temp-then-rename (NFR-4) for individual slice files")
	}
}

func TestMigrateStateMd_C86_CrashSafetyDescription(t *testing.T) {
	doc := readMigrateStateMd(t)
	hasCrashSafety := strings.Contains(doc, "crash") ||
		strings.Contains(doc, "Crash") ||
		strings.Contains(doc, "worst case") ||
		strings.Contains(doc, "worst-case") ||
		strings.Contains(doc, "power loss") ||
		strings.Contains(doc, "stale slices")
	if !hasCrashSafety {
		t.Error("migrate-state.md missing crash safety description (worst case scenario with stale slices/ alongside intact STATE.md)")
	}
}
