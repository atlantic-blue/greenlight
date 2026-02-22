package installer_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// C-83 Tests: StatusSliceAggregation
// C-84 Tests: SupportingCommandStateAdaptation
// These tests verify that the 6 supporting command markdown files contain all
// required instructions as defined in contracts C-83 and C-84 for slice S-31.

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func readStatusMd(t *testing.T) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/status.md"))
	if err != nil {
		t.Fatalf("failed to read src/commands/gl/status.md: %v", err)
	}
	return string(content)
}

func readPauseMd(t *testing.T) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/pause.md"))
	if err != nil {
		t.Fatalf("failed to read src/commands/gl/pause.md: %v", err)
	}
	return string(content)
}

func readResumeMd(t *testing.T) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/resume.md"))
	if err != nil {
		t.Fatalf("failed to read src/commands/gl/resume.md: %v", err)
	}
	return string(content)
}

func readShipMd(t *testing.T) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/ship.md"))
	if err != nil {
		t.Fatalf("failed to read src/commands/gl/ship.md: %v", err)
	}
	return string(content)
}

func readAddSliceMd(t *testing.T) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/add-slice.md"))
	if err != nil {
		t.Fatalf("failed to read src/commands/gl/add-slice.md: %v", err)
	}
	return string(content)
}

func readQuickMd(t *testing.T) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/quick.md"))
	if err != nil {
		t.Fatalf("failed to read src/commands/gl/quick.md: %v", err)
	}
	return string(content)
}

// ---------------------------------------------------------------------------
// C-83: StatusSliceAggregation — status.md
// ---------------------------------------------------------------------------

// State detection

func TestStatusMd_C83_ContainsStateFormatDetection(t *testing.T) {
	doc := readStatusMd(t)
	hasDetection := strings.Contains(doc, "C-80") ||
		strings.Contains(doc, "state format") ||
		strings.Contains(doc, "State Format") ||
		strings.Contains(doc, "detect")
	if !hasDetection {
		t.Error("status.md missing state format detection instruction (C-80 reference)")
	}
}

func TestStatusMd_C83_ContainsSlicesDirectoryReference(t *testing.T) {
	doc := readStatusMd(t)
	if !strings.Contains(doc, ".greenlight/slices/") {
		t.Error("status.md missing .greenlight/slices/ directory reference")
	}
}

// File-per-slice read (behaviour 2a)

func TestStatusMd_C83_ContainsReadAllMdFilesInstruction(t *testing.T) {
	doc := readStatusMd(t)
	hasReadAll := strings.Contains(doc, "all .md files") ||
		strings.Contains(doc, "all slice files") ||
		strings.Contains(doc, "Read all") ||
		strings.Contains(doc, "read all")
	if !hasReadAll {
		t.Error("status.md missing instruction to read all .md files from slices/ directory")
	}
}

// Parse frontmatter (behaviour 2b)

func TestStatusMd_C83_ContainsParseFrontmatterInstruction(t *testing.T) {
	doc := readStatusMd(t)
	hasFrontmatter := strings.Contains(doc, "frontmatter") ||
		strings.Contains(doc, "Frontmatter") ||
		strings.Contains(doc, "front matter") ||
		strings.Contains(doc, "front-matter")
	if !hasFrontmatter {
		t.Error("status.md missing instruction to parse frontmatter from each slice file")
	}
}

// Compute slice table (behaviour 2c)

func TestStatusMd_C83_ContainsSliceTableWithIDField(t *testing.T) {
	doc := readStatusMd(t)
	hasID := strings.Contains(doc, "ID") || strings.Contains(doc, " id ")
	if !hasID {
		t.Error("status.md missing slice table ID field instruction")
	}
}

func TestStatusMd_C83_ContainsSliceTableWithNameField(t *testing.T) {
	doc := readStatusMd(t)
	hasName := strings.Contains(doc, "Name") || strings.Contains(doc, "name")
	if !hasName {
		t.Error("status.md missing slice table Name field instruction")
	}
}

func TestStatusMd_C83_ContainsSliceTableWithStatusField(t *testing.T) {
	doc := readStatusMd(t)
	hasStatus := strings.Contains(doc, "Status") || strings.Contains(doc, "status")
	if !hasStatus {
		t.Error("status.md missing slice table Status field instruction")
	}
}

func TestStatusMd_C83_ContainsSliceTableWithTestsField(t *testing.T) {
	doc := readStatusMd(t)
	hasTests := strings.Contains(doc, "Tests") || strings.Contains(doc, "tests")
	if !hasTests {
		t.Error("status.md missing slice table Tests field instruction")
	}
}

// Compute progress (behaviour 2c)

func TestStatusMd_C83_ContainsProgressComputation(t *testing.T) {
	doc := readStatusMd(t)
	hasProgress := strings.Contains(doc, "progress") ||
		strings.Contains(doc, "Progress") ||
		strings.Contains(doc, "done/total") ||
		strings.Contains(doc, "complete count")
	if !hasProgress {
		t.Error("status.md missing progress computation instruction (done/total slices)")
	}
}

// Compute current (behaviour 2c)

func TestStatusMd_C83_ContainsCurrentSlicesComputation(t *testing.T) {
	doc := readStatusMd(t)
	hasCurrent := strings.Contains(doc, "Current") ||
		strings.Contains(doc, "current") ||
		strings.Contains(doc, "in-progress") ||
		strings.Contains(doc, "in progress")
	if !hasCurrent {
		t.Error("status.md missing current (in-progress) slices computation instruction")
	}
}

// Compute test summary (behaviour 2c)

func TestStatusMd_C83_ContainsTestSummaryComputation(t *testing.T) {
	doc := readStatusMd(t)
	hasTestSummary := strings.Contains(doc, "Test Summary") ||
		strings.Contains(doc, "test summary") ||
		strings.Contains(doc, "sum of tests") ||
		strings.Contains(doc, "security_tests")
	if !hasTestSummary {
		t.Error("status.md missing test summary computation instruction (sum of tests from all files)")
	}
}

// Project state (behaviour 2d)

func TestStatusMd_C83_ContainsProjectStateJsonRead(t *testing.T) {
	doc := readStatusMd(t)
	if !strings.Contains(doc, "project-state.json") {
		t.Error("status.md missing project-state.json read instruction")
	}
}

func TestStatusMd_C83_ContainsOverviewSessionBlockersRead(t *testing.T) {
	doc := readStatusMd(t)
	hasOverview := strings.Contains(doc, "overview") || strings.Contains(doc, "Overview")
	hasSession := strings.Contains(doc, "session") || strings.Contains(doc, "Session")
	hasBlockers := strings.Contains(doc, "blockers") || strings.Contains(doc, "Blockers")
	if !hasOverview {
		t.Error("status.md missing overview field read from project-state.json")
	}
	if !hasSession {
		t.Error("status.md missing session field read from project-state.json")
	}
	if !hasBlockers {
		t.Error("status.md missing blockers field read from project-state.json")
	}
}

// Display summary (behaviour 2e)

func TestStatusMd_C83_ContainsDisplaySummaryInstruction(t *testing.T) {
	doc := readStatusMd(t)
	hasDisplay := strings.Contains(doc, "Display") ||
		strings.Contains(doc, "display") ||
		strings.Contains(doc, "summary")
	if !hasDisplay {
		t.Error("status.md missing instruction to display computed summary to user")
	}
}

// Regenerate STATE.md (behaviour 2f)

func TestStatusMd_C83_ContainsStateMdRegenerationInstruction(t *testing.T) {
	doc := readStatusMd(t)
	hasRegen := strings.Contains(doc, "Regenerate STATE.md") ||
		strings.Contains(doc, "regenerate STATE.md") ||
		(strings.Contains(doc, "Regenerate") && strings.Contains(doc, "STATE.md")) ||
		strings.Contains(doc, "D-34")
	if !hasRegen {
		t.Error("status.md missing STATE.md regeneration instruction (D-34)")
	}
}

func TestStatusMd_C83_ContainsRegenerationAfterDisplayInstruction(t *testing.T) {
	doc := readStatusMd(t)
	displayPos := strings.Index(doc, "Display")
	if displayPos == -1 {
		displayPos = strings.Index(doc, "display")
	}
	regenPos := strings.Index(doc, "Regenerate STATE")
	if regenPos == -1 {
		regenPos = strings.Index(doc, "regenerate STATE")
	}
	if regenPos == -1 {
		regenPos = strings.Index(doc, "D-34")
	}
	if displayPos == -1 {
		t.Error("status.md missing display instruction")
	}
	if regenPos == -1 {
		t.Error("status.md missing STATE.md regeneration instruction")
	}
	if displayPos != -1 && regenPos != -1 && regenPos <= displayPos {
		t.Error("status.md STATE.md regeneration must come after display (regeneration happens after display)")
	}
}

// Legacy fallback (behaviour 3)

func TestStatusMd_C83_ContainsLegacyFallback(t *testing.T) {
	doc := readStatusMd(t)
	hasLegacy := strings.Contains(doc, "legacy") ||
		strings.Contains(doc, "Legacy") ||
		strings.Contains(doc, "as before") ||
		strings.Contains(doc, "backward")
	if !hasLegacy {
		t.Error("status.md missing legacy format fallback instruction (STATE.md as before)")
	}
}

func TestStatusMd_C83_ContainsLegacyUnchangedBehaviourInstruction(t *testing.T) {
	doc := readStatusMd(t)
	hasUnchanged := strings.Contains(doc, "unchanged") ||
		strings.Contains(doc, "no change") ||
		strings.Contains(doc, "as before")
	if !hasUnchanged {
		t.Error("status.md missing instruction that legacy format display is completely unchanged")
	}
}

// Error handling: NoSliceFiles

func TestStatusMd_C83_ContainsNoSliceFilesError(t *testing.T) {
	doc := readStatusMd(t)
	hasNoSliceFiles := strings.Contains(doc, "NoSliceFiles") ||
		strings.Contains(doc, "no .md files") ||
		strings.Contains(doc, "empty summary") ||
		(strings.Contains(doc, "no slice") && strings.Contains(doc, "files"))
	if !hasNoSliceFiles {
		t.Error("status.md missing NoSliceFiles error handling (display empty summary when no .md files)")
	}
}

// Error handling: CorruptSliceFile

func TestStatusMd_C83_ContainsCorruptSliceFileError(t *testing.T) {
	doc := readStatusMd(t)
	hasCorrupt := strings.Contains(doc, "CorruptSliceFile") ||
		strings.Contains(doc, "corrupt") ||
		strings.Contains(doc, "Corrupt") ||
		strings.Contains(doc, "invalid frontmatter") ||
		strings.Contains(doc, "invalid")
	if !hasCorrupt {
		t.Error("status.md missing CorruptSliceFile error handling (skip and warn on invalid frontmatter)")
	}
}

func TestStatusMd_C83_ContainsCorruptFileSkipAndWarnInstruction(t *testing.T) {
	doc := readStatusMd(t)
	hasSkip := strings.Contains(doc, "Skip") ||
		strings.Contains(doc, "skip") ||
		strings.Contains(doc, "Warn") ||
		strings.Contains(doc, "warn")
	if !hasSkip {
		t.Error("status.md missing instruction to skip corrupt slice file and warn")
	}
}

// Error handling: ProjectStateReadFailure

func TestStatusMd_C83_ContainsProjectStateReadFailureError(t *testing.T) {
	doc := readStatusMd(t)
	hasFailure := strings.Contains(doc, "ProjectStateReadFailure") ||
		strings.Contains(doc, "project-state") && strings.Contains(doc, "warn") ||
		strings.Contains(doc, "without overview") ||
		strings.Contains(doc, "cannot read project-state") ||
		strings.Contains(doc, "Cannot read project-state")
	if !hasFailure {
		t.Error("status.md missing ProjectStateReadFailure error handling (display without overview, warn)")
	}
}

// Invariants: fresh data

func TestStatusMd_C83_ContainsFreshDataInvariant(t *testing.T) {
	doc := readStatusMd(t)
	hasFresh := strings.Contains(doc, "never cached") ||
		strings.Contains(doc, "fresh") ||
		strings.Contains(doc, "every invocation") ||
		strings.Contains(doc, "every /gl:status")
	if !hasFresh {
		t.Error("status.md missing fresh data invariant (status computed fresh on every invocation, never cached)")
	}
}

// Invariants: slice table sorted by ID

func TestStatusMd_C83_ContainsSliceTableSortedByIdInvariant(t *testing.T) {
	doc := readStatusMd(t)
	hasSorted := strings.Contains(doc, "sorted") ||
		strings.Contains(doc, "ascending") ||
		strings.Contains(doc, "sort") ||
		strings.Contains(doc, "order")
	if !hasSorted {
		t.Error("status.md missing slice table sorted by ID (ascending) invariant")
	}
}

// ---------------------------------------------------------------------------
// C-84: SupportingCommandStateAdaptation — shared checks for all 5 commands
// ---------------------------------------------------------------------------

// /gl:pause — state detection

func TestPauseMd_C84_ContainsStateFormatDetection(t *testing.T) {
	doc := readPauseMd(t)
	hasDetection := strings.Contains(doc, "C-80") ||
		strings.Contains(doc, "state format") ||
		strings.Contains(doc, "State Format") ||
		strings.Contains(doc, "detect")
	if !hasDetection {
		t.Error("pause.md missing state format detection instruction (C-80 reference)")
	}
}

func TestPauseMd_C84_ContainsFilePerSliceReference(t *testing.T) {
	doc := readPauseMd(t)
	if !strings.Contains(doc, ".greenlight/slices/") {
		t.Error("pause.md missing .greenlight/slices/ file-per-slice reference")
	}
}

func TestPauseMd_C84_ContainsLegacyFallback(t *testing.T) {
	doc := readPauseMd(t)
	hasLegacy := strings.Contains(doc, "legacy") ||
		strings.Contains(doc, "Legacy") ||
		strings.Contains(doc, "as before") ||
		strings.Contains(doc, "backward")
	if !hasLegacy {
		t.Error("pause.md missing legacy format fallback instruction (STATE.md as before)")
	}
}

// /gl:pause — write to own slice file

func TestPauseMd_C84_ContainsWriteToOwnSliceFile(t *testing.T) {
	doc := readPauseMd(t)
	hasOwnSlice := strings.Contains(doc, "own slice file") ||
		strings.Contains(doc, "slice file") ||
		strings.Contains(doc, "Write pause") ||
		strings.Contains(doc, "write pause")
	if !hasOwnSlice {
		t.Error("pause.md missing instruction to write pause state to own slice file")
	}
}

// /gl:pause — write resume context to project-state.json

func TestPauseMd_C84_ContainsWriteResumeContextToProjectStateJson(t *testing.T) {
	doc := readPauseMd(t)
	hasResumeContext := strings.Contains(doc, "resume context") ||
		strings.Contains(doc, "Resume context") ||
		(strings.Contains(doc, "project-state.json") && strings.Contains(doc, "resume"))
	if !hasResumeContext {
		t.Error("pause.md missing instruction to write resume context to project-state.json")
	}
}

func TestPauseMd_C84_ContainsProjectStateJsonReference(t *testing.T) {
	doc := readPauseMd(t)
	if !strings.Contains(doc, "project-state.json") {
		t.Error("pause.md missing project-state.json reference")
	}
}

// /gl:pause — regenerate STATE.md

func TestPauseMd_C84_ContainsStateMdRegeneration(t *testing.T) {
	doc := readPauseMd(t)
	hasRegen := strings.Contains(doc, "Regenerate STATE.md") ||
		strings.Contains(doc, "regenerate STATE.md") ||
		(strings.Contains(doc, "Regenerate") && strings.Contains(doc, "STATE.md")) ||
		strings.Contains(doc, "D-34")
	if !hasRegen {
		t.Error("pause.md missing STATE.md regeneration instruction (D-34)")
	}
}

// /gl:pause — error handling

func TestPauseMd_C84_ContainsFormatDetectionFailureError(t *testing.T) {
	doc := readPauseMd(t)
	hasFailure := strings.Contains(doc, "FormatDetectionFailure") ||
		(strings.Contains(doc, "format") && strings.Contains(doc, "failure")) ||
		strings.Contains(doc, "/gl:init")
	if !hasFailure {
		t.Error("pause.md missing FormatDetectionFailure error handling")
	}
}

func TestPauseMd_C84_ContainsSliceFileNotFoundError(t *testing.T) {
	doc := readPauseMd(t)
	hasNotFound := strings.Contains(doc, "SliceFileNotFound") ||
		(strings.Contains(doc, "slice file") && strings.Contains(doc, "not found")) ||
		(strings.Contains(doc, "slice file") && strings.Contains(doc, "does not exist"))
	if !hasNotFound {
		t.Error("pause.md missing SliceFileNotFound error handling")
	}
}

func TestPauseMd_C84_ContainsRegenerationFailureWarnInstruction(t *testing.T) {
	doc := readPauseMd(t)
	hasWarnContinue := strings.Contains(doc, "RegenerationFailure") ||
		strings.Contains(doc, "Warn but continue") ||
		strings.Contains(doc, "warn but continue") ||
		(strings.Contains(doc, "warn") && strings.Contains(doc, "continue"))
	if !hasWarnContinue {
		t.Error("pause.md missing RegenerationFailure warn-but-continue instruction")
	}
}

// /gl:pause — write-to-temp-then-rename invariant

func TestPauseMd_C84_ContainsWriteToTempThenRenameInvariant(t *testing.T) {
	doc := readPauseMd(t)
	hasAtomic := strings.Contains(doc, "write-to-temp-then-rename") ||
		strings.Contains(doc, "temp-then-rename") ||
		strings.Contains(doc, "atomic") ||
		strings.Contains(doc, "NFR-4") ||
		strings.Contains(doc, "temp") && strings.Contains(doc, "rename")
	if !hasAtomic {
		t.Error("pause.md missing write-to-temp-then-rename atomic write invariant (NFR-4)")
	}
}

// ---------------------------------------------------------------------------
// C-84: SupportingCommandStateAdaptation — /gl:resume
// ---------------------------------------------------------------------------

func TestResumeMd_C84_ContainsStateFormatDetection(t *testing.T) {
	doc := readResumeMd(t)
	hasDetection := strings.Contains(doc, "C-80") ||
		strings.Contains(doc, "state format") ||
		strings.Contains(doc, "State Format") ||
		strings.Contains(doc, "detect")
	if !hasDetection {
		t.Error("resume.md missing state format detection instruction (C-80 reference)")
	}
}

func TestResumeMd_C84_ContainsFilePerSliceReference(t *testing.T) {
	doc := readResumeMd(t)
	if !strings.Contains(doc, ".greenlight/slices/") {
		t.Error("resume.md missing .greenlight/slices/ file-per-slice reference")
	}
}

func TestResumeMd_C84_ContainsLegacyFallback(t *testing.T) {
	doc := readResumeMd(t)
	hasLegacy := strings.Contains(doc, "legacy") ||
		strings.Contains(doc, "Legacy") ||
		strings.Contains(doc, "as before") ||
		strings.Contains(doc, "backward")
	if !hasLegacy {
		t.Error("resume.md missing legacy format fallback instruction (STATE.md as before)")
	}
}

// /gl:resume — read slice files

func TestResumeMd_C84_ContainsReadSliceFilesInstruction(t *testing.T) {
	doc := readResumeMd(t)
	hasRead := strings.Contains(doc, "Read all slice files") ||
		strings.Contains(doc, "read all slice files") ||
		strings.Contains(doc, "all slice files") ||
		strings.Contains(doc, "slice files")
	if !hasRead {
		t.Error("resume.md missing instruction to read slice files to determine resumable state")
	}
}

// /gl:resume — read project-state.json for resume context

func TestResumeMd_C84_ContainsReadResumeContextFromProjectStateJson(t *testing.T) {
	doc := readResumeMd(t)
	hasResumeContext := strings.Contains(doc, "resume context") ||
		strings.Contains(doc, "Resume context") ||
		(strings.Contains(doc, "project-state.json") && strings.Contains(doc, "resume"))
	if !hasResumeContext {
		t.Error("resume.md missing instruction to read resume context from project-state.json")
	}
}

func TestResumeMd_C84_ContainsProjectStateJsonReference(t *testing.T) {
	doc := readResumeMd(t)
	if !strings.Contains(doc, "project-state.json") {
		t.Error("resume.md missing project-state.json reference")
	}
}

// /gl:resume — error handling

func TestResumeMd_C84_ContainsFormatDetectionFailureError(t *testing.T) {
	doc := readResumeMd(t)
	hasFailure := strings.Contains(doc, "FormatDetectionFailure") ||
		(strings.Contains(doc, "format") && strings.Contains(doc, "failure")) ||
		strings.Contains(doc, "/gl:init")
	if !hasFailure {
		t.Error("resume.md missing FormatDetectionFailure error handling")
	}
}

func TestResumeMd_C84_ContainsSliceFileNotFoundError(t *testing.T) {
	doc := readResumeMd(t)
	hasNotFound := strings.Contains(doc, "SliceFileNotFound") ||
		(strings.Contains(doc, "slice file") && strings.Contains(doc, "not found")) ||
		(strings.Contains(doc, "slice file") && strings.Contains(doc, "does not exist"))
	if !hasNotFound {
		t.Error("resume.md missing SliceFileNotFound error handling")
	}
}

func TestResumeMd_C84_ContainsRegenerationFailureWarnInstruction(t *testing.T) {
	doc := readResumeMd(t)
	hasWarnContinue := strings.Contains(doc, "RegenerationFailure") ||
		strings.Contains(doc, "Warn but continue") ||
		strings.Contains(doc, "warn but continue") ||
		(strings.Contains(doc, "warn") && strings.Contains(doc, "continue"))
	if !hasWarnContinue {
		t.Error("resume.md missing RegenerationFailure warn-but-continue instruction")
	}
}

// ---------------------------------------------------------------------------
// C-84: SupportingCommandStateAdaptation — /gl:ship
// ---------------------------------------------------------------------------

func TestShipMd_C84_ContainsStateFormatDetection(t *testing.T) {
	doc := readShipMd(t)
	hasDetection := strings.Contains(doc, "C-80") ||
		strings.Contains(doc, "state format") ||
		strings.Contains(doc, "State Format") ||
		strings.Contains(doc, "detect")
	if !hasDetection {
		t.Error("ship.md missing state format detection instruction (C-80 reference)")
	}
}

func TestShipMd_C84_ContainsFilePerSliceReference(t *testing.T) {
	doc := readShipMd(t)
	if !strings.Contains(doc, ".greenlight/slices/") {
		t.Error("ship.md missing .greenlight/slices/ file-per-slice reference")
	}
}

func TestShipMd_C84_ContainsLegacyFallback(t *testing.T) {
	doc := readShipMd(t)
	hasLegacy := strings.Contains(doc, "legacy") ||
		strings.Contains(doc, "Legacy") ||
		strings.Contains(doc, "as before") ||
		strings.Contains(doc, "backward")
	if !hasLegacy {
		t.Error("ship.md missing legacy format fallback instruction (STATE.md as before)")
	}
}

// /gl:ship — read ALL slice files

func TestShipMd_C84_ContainsReadAllSliceFilesInstruction(t *testing.T) {
	doc := readShipMd(t)
	hasReadAll := strings.Contains(doc, "all slice files") ||
		strings.Contains(doc, "Read all slice") ||
		strings.Contains(doc, "read all slice") ||
		strings.Contains(doc, "all .md files")
	if !hasReadAll {
		t.Error("ship.md missing instruction to read ALL slice files")
	}
}

// /gl:ship — completeness pre-check

func TestShipMd_C84_ContainsCompletenessPreCheck(t *testing.T) {
	doc := readShipMd(t)
	hasPreCheck := strings.Contains(doc, "Pre-check") ||
		strings.Contains(doc, "pre-check") ||
		strings.Contains(doc, "all must have status complete") ||
		strings.Contains(doc, "all slices complete") ||
		strings.Contains(doc, "all slices must") ||
		(strings.Contains(doc, "complete") && strings.Contains(doc, "pre"))
	if !hasPreCheck {
		t.Error("ship.md missing completeness pre-check instruction (all slices must be complete)")
	}
}

func TestShipMd_C84_ContainsIncompleteSlicesReport(t *testing.T) {
	doc := readShipMd(t)
	hasIncomplete := strings.Contains(doc, "incomplete") ||
		strings.Contains(doc, "Incomplete") ||
		strings.Contains(doc, "non-complete") ||
		strings.Contains(doc, "not complete") ||
		strings.Contains(doc, "which slices")
	if !hasIncomplete {
		t.Error("ship.md missing instruction to report which slices are incomplete")
	}
}

// /gl:ship — error handling

func TestShipMd_C84_ContainsFormatDetectionFailureError(t *testing.T) {
	doc := readShipMd(t)
	hasFailure := strings.Contains(doc, "FormatDetectionFailure") ||
		(strings.Contains(doc, "format") && strings.Contains(doc, "failure")) ||
		strings.Contains(doc, "/gl:init")
	if !hasFailure {
		t.Error("ship.md missing FormatDetectionFailure error handling")
	}
}

func TestShipMd_C84_ContainsSliceFileNotFoundError(t *testing.T) {
	doc := readShipMd(t)
	hasNotFound := strings.Contains(doc, "SliceFileNotFound") ||
		(strings.Contains(doc, "slice file") && strings.Contains(doc, "not found")) ||
		(strings.Contains(doc, "slice file") && strings.Contains(doc, "does not exist"))
	if !hasNotFound {
		t.Error("ship.md missing SliceFileNotFound error handling")
	}
}

func TestShipMd_C84_ContainsRegenerationFailureWarnInstruction(t *testing.T) {
	doc := readShipMd(t)
	hasWarnContinue := strings.Contains(doc, "RegenerationFailure") ||
		strings.Contains(doc, "Warn but continue") ||
		strings.Contains(doc, "warn but continue") ||
		(strings.Contains(doc, "warn") && strings.Contains(doc, "continue"))
	if !hasWarnContinue {
		t.Error("ship.md missing RegenerationFailure warn-but-continue instruction")
	}
}

// ---------------------------------------------------------------------------
// C-84: SupportingCommandStateAdaptation — /gl:add-slice
// ---------------------------------------------------------------------------

func TestAddSliceMd_C84_ContainsStateFormatDetection(t *testing.T) {
	doc := readAddSliceMd(t)
	hasDetection := strings.Contains(doc, "C-80") ||
		strings.Contains(doc, "state format") ||
		strings.Contains(doc, "State Format") ||
		strings.Contains(doc, "detect")
	if !hasDetection {
		t.Error("add-slice.md missing state format detection instruction (C-80 reference)")
	}
}

func TestAddSliceMd_C84_ContainsFilePerSliceReference(t *testing.T) {
	doc := readAddSliceMd(t)
	if !strings.Contains(doc, ".greenlight/slices/") {
		t.Error("add-slice.md missing .greenlight/slices/ file-per-slice reference")
	}
}

func TestAddSliceMd_C84_ContainsLegacyFallback(t *testing.T) {
	doc := readAddSliceMd(t)
	hasLegacy := strings.Contains(doc, "legacy") ||
		strings.Contains(doc, "Legacy") ||
		strings.Contains(doc, "as before") ||
		strings.Contains(doc, "backward")
	if !hasLegacy {
		t.Error("add-slice.md missing legacy format fallback instruction (STATE.md as before)")
	}
}

// /gl:add-slice — create new slice file in slices/

func TestAddSliceMd_C84_ContainsCreateSliceFileInstruction(t *testing.T) {
	doc := readAddSliceMd(t)
	hasCreate := strings.Contains(doc, "Create new slice file") ||
		strings.Contains(doc, "create new slice file") ||
		strings.Contains(doc, "Create slice file") ||
		strings.Contains(doc, "new slice file") ||
		(strings.Contains(doc, "create") && strings.Contains(doc, "slices/"))
	if !hasCreate {
		t.Error("add-slice.md missing instruction to create new slice file in .greenlight/slices/")
	}
}

// /gl:add-slice — validate slice ID

func TestAddSliceMd_C84_ContainsSliceIdValidation(t *testing.T) {
	doc := readAddSliceMd(t)
	hasValidation := strings.Contains(doc, "validate") ||
		strings.Contains(doc, "Validate") ||
		strings.Contains(doc, "validation") ||
		strings.Contains(doc, "Validation")
	if !hasValidation {
		t.Error("add-slice.md missing slice ID validation instruction")
	}
}

func TestAddSliceMd_C84_ContainsPathTraversalPrevention(t *testing.T) {
	doc := readAddSliceMd(t)
	hasPathTraversal := strings.Contains(doc, "path traversal") ||
		strings.Contains(doc, "traversal") ||
		strings.Contains(doc, "path validation") ||
		(strings.Contains(doc, "validate") && strings.Contains(doc, "ID"))
	if !hasPathTraversal {
		t.Error("add-slice.md missing path traversal prevention / slice ID validation instruction")
	}
}

// /gl:add-slice — regenerate STATE.md

func TestAddSliceMd_C84_ContainsStateMdRegeneration(t *testing.T) {
	doc := readAddSliceMd(t)
	hasRegen := strings.Contains(doc, "Regenerate STATE.md") ||
		strings.Contains(doc, "regenerate STATE.md") ||
		(strings.Contains(doc, "Regenerate") && strings.Contains(doc, "STATE.md")) ||
		strings.Contains(doc, "D-34")
	if !hasRegen {
		t.Error("add-slice.md missing STATE.md regeneration instruction (D-34)")
	}
}

// /gl:add-slice — error handling

func TestAddSliceMd_C84_ContainsFormatDetectionFailureError(t *testing.T) {
	doc := readAddSliceMd(t)
	hasFailure := strings.Contains(doc, "FormatDetectionFailure") ||
		(strings.Contains(doc, "format") && strings.Contains(doc, "failure")) ||
		strings.Contains(doc, "/gl:init")
	if !hasFailure {
		t.Error("add-slice.md missing FormatDetectionFailure error handling")
	}
}

func TestAddSliceMd_C84_ContainsSliceFileNotFoundError(t *testing.T) {
	doc := readAddSliceMd(t)
	hasNotFound := strings.Contains(doc, "SliceFileNotFound") ||
		(strings.Contains(doc, "slice file") && strings.Contains(doc, "not found")) ||
		(strings.Contains(doc, "does not exist") && strings.Contains(doc, "slice"))
	if !hasNotFound {
		t.Error("add-slice.md missing SliceFileNotFound error handling")
	}
}

func TestAddSliceMd_C84_ContainsRegenerationFailureWarnInstruction(t *testing.T) {
	doc := readAddSliceMd(t)
	hasWarnContinue := strings.Contains(doc, "RegenerationFailure") ||
		strings.Contains(doc, "Warn but continue") ||
		strings.Contains(doc, "warn but continue") ||
		(strings.Contains(doc, "warn") && strings.Contains(doc, "continue"))
	if !hasWarnContinue {
		t.Error("add-slice.md missing RegenerationFailure warn-but-continue instruction")
	}
}

// /gl:add-slice — write-to-temp-then-rename invariant

func TestAddSliceMd_C84_ContainsWriteToTempThenRenameInvariant(t *testing.T) {
	doc := readAddSliceMd(t)
	hasAtomic := strings.Contains(doc, "write-to-temp-then-rename") ||
		strings.Contains(doc, "temp-then-rename") ||
		strings.Contains(doc, "atomic") ||
		strings.Contains(doc, "NFR-4") ||
		(strings.Contains(doc, "temp") && strings.Contains(doc, "rename"))
	if !hasAtomic {
		t.Error("add-slice.md missing write-to-temp-then-rename atomic write invariant (NFR-4)")
	}
}

// ---------------------------------------------------------------------------
// C-84: SupportingCommandStateAdaptation — /gl:quick
// ---------------------------------------------------------------------------

func TestQuickMd_C84_ContainsStateFormatDetection(t *testing.T) {
	doc := readQuickMd(t)
	hasDetection := strings.Contains(doc, "C-80") ||
		strings.Contains(doc, "state format") ||
		strings.Contains(doc, "State Format") ||
		strings.Contains(doc, "detect")
	if !hasDetection {
		t.Error("quick.md missing state format detection instruction (C-80 reference)")
	}
}

func TestQuickMd_C84_ContainsFilePerSliceReference(t *testing.T) {
	doc := readQuickMd(t)
	if !strings.Contains(doc, ".greenlight/slices/") {
		t.Error("quick.md missing .greenlight/slices/ file-per-slice reference")
	}
}

func TestQuickMd_C84_ContainsLegacyFallback(t *testing.T) {
	doc := readQuickMd(t)
	hasLegacy := strings.Contains(doc, "legacy") ||
		strings.Contains(doc, "Legacy") ||
		strings.Contains(doc, "as before") ||
		strings.Contains(doc, "backward")
	if !hasLegacy {
		t.Error("quick.md missing legacy format fallback instruction (STATE.md as before)")
	}
}

// /gl:quick — update test counts in slice file

func TestQuickMd_C84_ContainsUpdateTestCountsInstruction(t *testing.T) {
	doc := readQuickMd(t)
	hasUpdate := strings.Contains(doc, "Update test counts") ||
		strings.Contains(doc, "update test counts") ||
		strings.Contains(doc, "test counts") ||
		(strings.Contains(doc, "update") && strings.Contains(doc, "slice file") && strings.Contains(doc, "test"))
	if !hasUpdate {
		t.Error("quick.md missing instruction to update test counts in the relevant slice file")
	}
}

func TestQuickMd_C84_ContainsRelevantSliceFileOnlyUpdate(t *testing.T) {
	doc := readQuickMd(t)
	hasRelevant := strings.Contains(doc, "relevant slice file") ||
		strings.Contains(doc, "only the relevant") ||
		strings.Contains(doc, "relevant slice") ||
		strings.Contains(doc, "own slice")
	if !hasRelevant {
		t.Error("quick.md missing instruction that only the relevant slice file is updated (not all files)")
	}
}

// /gl:quick — regenerate STATE.md

func TestQuickMd_C84_ContainsStateMdRegeneration(t *testing.T) {
	doc := readQuickMd(t)
	hasRegen := strings.Contains(doc, "Regenerate STATE.md") ||
		strings.Contains(doc, "regenerate STATE.md") ||
		(strings.Contains(doc, "Regenerate") && strings.Contains(doc, "STATE.md")) ||
		strings.Contains(doc, "D-34")
	if !hasRegen {
		t.Error("quick.md missing STATE.md regeneration instruction (D-34)")
	}
}

// /gl:quick — error handling

func TestQuickMd_C84_ContainsFormatDetectionFailureError(t *testing.T) {
	doc := readQuickMd(t)
	hasFailure := strings.Contains(doc, "FormatDetectionFailure") ||
		(strings.Contains(doc, "format") && strings.Contains(doc, "failure")) ||
		strings.Contains(doc, "/gl:init")
	if !hasFailure {
		t.Error("quick.md missing FormatDetectionFailure error handling")
	}
}

func TestQuickMd_C84_ContainsSliceFileNotFoundError(t *testing.T) {
	doc := readQuickMd(t)
	hasNotFound := strings.Contains(doc, "SliceFileNotFound") ||
		(strings.Contains(doc, "slice file") && strings.Contains(doc, "not found")) ||
		(strings.Contains(doc, "slice file") && strings.Contains(doc, "does not exist"))
	if !hasNotFound {
		t.Error("quick.md missing SliceFileNotFound error handling")
	}
}

func TestQuickMd_C84_ContainsRegenerationFailureWarnInstruction(t *testing.T) {
	doc := readQuickMd(t)
	hasWarnContinue := strings.Contains(doc, "RegenerationFailure") ||
		strings.Contains(doc, "Warn but continue") ||
		strings.Contains(doc, "warn but continue") ||
		(strings.Contains(doc, "warn") && strings.Contains(doc, "continue"))
	if !hasWarnContinue {
		t.Error("quick.md missing RegenerationFailure warn-but-continue instruction")
	}
}

// /gl:quick — write-to-temp-then-rename invariant

func TestQuickMd_C84_ContainsWriteToTempThenRenameInvariant(t *testing.T) {
	doc := readQuickMd(t)
	hasAtomic := strings.Contains(doc, "write-to-temp-then-rename") ||
		strings.Contains(doc, "temp-then-rename") ||
		strings.Contains(doc, "atomic") ||
		strings.Contains(doc, "NFR-4") ||
		(strings.Contains(doc, "temp") && strings.Contains(doc, "rename"))
	if !hasAtomic {
		t.Error("quick.md missing write-to-temp-then-rename atomic write invariant (NFR-4)")
	}
}
