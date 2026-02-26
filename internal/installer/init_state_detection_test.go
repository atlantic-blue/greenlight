package installer_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// C-78 Tests: InitSliceDirectory
// C-79 Tests: InitProjectState
// C-80 Tests: InitStateDetection
// These tests verify that src/commands/gl/init.md contains all required instructions
// as defined in contracts C-78, C-79, and C-80 for slice S-29.
// Verification: verify (C-78, C-80), auto (C-79)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func readInitMd(t *testing.T) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/init.md"))
	if err != nil {
		t.Fatalf("failed to read src/commands/gl/init.md: %v", err)
	}
	return string(content)
}

// ---------------------------------------------------------------------------
// C-78: InitSliceDirectory
// ---------------------------------------------------------------------------

// Success cases

func TestInitMd_C78_ContainsSlicesDirectoryCreationInstruction(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, ".greenlight/slices/") {
		t.Error("init.md missing instruction to create .greenlight/slices/ directory")
	}
}

func TestInitMd_C78_ContainsIndividualSliceFileCreationInstruction(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "GRAPH.json") {
		t.Error("init.md missing instruction to create individual slice files from GRAPH.json")
	}
}

func TestInitMd_C78_ContainsFrontmatterIdField(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "id") {
		t.Error("init.md missing frontmatter 'id' field reference for slice files")
	}
}

func TestInitMd_C78_ContainsFrontmatterStatusField(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "status") {
		t.Error("init.md missing frontmatter 'status' field reference for slice files")
	}
}

func TestInitMd_C78_ContainsFrontmatterStepField(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "step") {
		t.Error("init.md missing frontmatter 'step' field reference for slice files")
	}
}

func TestInitMd_C78_ContainsFrontmatterMilestoneField(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "milestone") {
		t.Error("init.md missing frontmatter 'milestone' field reference for slice files")
	}
}

func TestInitMd_C78_ContainsFrontmatterStartedField(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "started") {
		t.Error("init.md missing frontmatter 'started' field reference for slice files")
	}
}

func TestInitMd_C78_ContainsFrontmatterUpdatedField(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "updated") {
		t.Error("init.md missing frontmatter 'updated' field reference for slice files")
	}
}

func TestInitMd_C78_ContainsFrontmatterTestsField(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "tests") {
		t.Error("init.md missing frontmatter 'tests' field reference for slice files")
	}
}

func TestInitMd_C78_ContainsFrontmatterSecurityTestsField(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "security_tests") {
		t.Error("init.md missing frontmatter 'security_tests' field reference for slice files")
	}
}

func TestInitMd_C78_ContainsFrontmatterSessionField(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "session") {
		t.Error("init.md missing frontmatter 'session' field reference for slice files")
	}
}

func TestInitMd_C78_ContainsFrontmatterDepsField(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "deps") {
		t.Error("init.md missing frontmatter 'deps' field reference for slice files")
	}
}

func TestInitMd_C78_ContainsSliceIdValidationPattern(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "S-{digits}") && !strings.Contains(doc, "S-\\d") && !strings.Contains(doc, "S-[0-9]") {
		t.Error("init.md missing slice ID validation pattern (S-{digits}) instruction")
	}
}

func TestInitMd_C78_ContainsCrashSafetyInstruction(t *testing.T) {
	doc := readInitMd(t)
	hasTempRename := strings.Contains(doc, "write-to-temp-then-rename") || strings.Contains(doc, "temp-then-rename")
	hasAtomic := strings.Contains(doc, "atomic")
	if !hasTempRename && !hasAtomic {
		t.Error("init.md missing crash safety instruction (write-to-temp-then-rename or atomic write)")
	}
}

func TestInitMd_C78_ContainsStateMdGenerationWithGeneratedHeader(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "GENERATED") {
		t.Error("init.md missing STATE.md generation instruction with GENERATED header comment")
	}
}

// Error handling

func TestInitMd_C78_ContainsSlicesDirExistsError(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "SlicesDirExists") {
		t.Error("init.md missing SlicesDirExists error handling instruction")
	}
}

func TestInitMd_C78_ContainsInvalidSliceIdInGraphError(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "InvalidSliceIdInGraph") {
		t.Error("init.md missing InvalidSliceIdInGraph error handling instruction")
	}
}

func TestInitMd_C78_ContainsDirectoryCreateFailureError(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "DirectoryCreateFailure") {
		t.Error("init.md missing DirectoryCreateFailure error handling instruction")
	}
}

func TestInitMd_C78_ContainsFileWriteFailureError(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "FileWriteFailure") {
		t.Error("init.md missing FileWriteFailure error handling instruction")
	}
}

// Permissions

func TestInitMd_C78_ContainsDirectoryPermission0o755(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "0o755") && !strings.Contains(doc, "0755") {
		t.Error("init.md missing directory permission specification (0o755)")
	}
}

func TestInitMd_C78_ContainsFilePermission0o644(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "0o644") && !strings.Contains(doc, "0644") {
		t.Error("init.md missing file permission specification (0o644)")
	}
}

// Invariants

func TestInitMd_C78_ContainsNoAutoMigrateLegacyFormatInstruction(t *testing.T) {
	doc := readInitMd(t)
	hasNoMigrate := strings.Contains(doc, "auto-migrate") || strings.Contains(doc, "does NOT auto-migrate") || strings.Contains(doc, "not auto-migrate")
	if !hasNoMigrate {
		t.Error("init.md missing instruction that legacy format is not auto-migrated")
	}
}

// ---------------------------------------------------------------------------
// C-79: InitProjectState
// ---------------------------------------------------------------------------

// Success cases

func TestInitMd_C79_ContainsProjectStateJsonCreationInstruction(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "project-state.json") {
		t.Error("init.md missing instruction to create .greenlight/project-state.json")
	}
}

func TestInitMd_C79_ContainsOverviewSchemaSection(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "overview") {
		t.Error("init.md missing 'overview' key in project-state.json schema description")
	}
}

func TestInitMd_C79_ContainsSessionSchemaSection(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "session") {
		t.Error("init.md missing 'session' key in project-state.json schema description")
	}
}

func TestInitMd_C79_ContainsBlockersSchemaSection(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "blockers") {
		t.Error("init.md missing 'blockers' key in project-state.json schema description")
	}
}

func TestInitMd_C79_ContainsValuePropField(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "value_prop") {
		t.Error("init.md missing 'value_prop' field instruction for project-state.json overview")
	}
}

func TestInitMd_C79_ContainsStackField(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "stack") {
		t.Error("init.md missing 'stack' field instruction for project-state.json overview")
	}
}

func TestInitMd_C79_ContainsModeField(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "mode") {
		t.Error("init.md missing 'mode' field instruction for project-state.json overview")
	}
}

func TestInitMd_C79_ContainsLastSessionField(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "last_session") {
		t.Error("init.md missing 'last_session' field instruction for project-state.json session")
	}
}

func TestInitMd_C79_ContainsResumeFileNullInstruction(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "resume_file") {
		t.Error("init.md missing 'resume_file' field instruction for project-state.json session")
	}
}

func TestInitMd_C79_ContainsBlockersEmptyArrayInstruction(t *testing.T) {
	doc := readInitMd(t)
	hasEmptyArray := strings.Contains(doc, `"blockers": []`) || strings.Contains(doc, "blockers array") || strings.Contains(doc, "empty array")
	if !hasEmptyArray {
		t.Error("init.md missing instruction that blockers initialises as an empty array")
	}
}

// Error handling

func TestInitMd_C79_ContainsProjectStateExistsError(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "ProjectStateExists") {
		t.Error("init.md missing ProjectStateExists error handling instruction")
	}
}

func TestInitMd_C79_ContainsWriteFailureError(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "WriteFailure") {
		t.Error("init.md missing WriteFailure error handling instruction for project-state.json")
	}
}

func TestInitMd_C79_ContainsMissingDesignContextError(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "MissingDesignContext") {
		t.Error("init.md missing MissingDesignContext error handling instruction")
	}
}

// Invariants

func TestInitMd_C79_ContainsDefaultModeYolo(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "yolo") {
		t.Error("init.md missing instruction that overview.mode defaults to 'yolo'")
	}
}

// ---------------------------------------------------------------------------
// C-80: InitStateDetection
// ---------------------------------------------------------------------------

// Success cases

func TestInitMd_C80_ContainsStateDetectionFlowInstruction(t *testing.T) {
	doc := readInitMd(t)
	hasDetection := strings.Contains(doc, "state detection") || strings.Contains(doc, "State Detection") || strings.Contains(doc, "detection flow")
	if !hasDetection {
		t.Error("init.md missing state detection flow instruction")
	}
}

func TestInitMd_C80_ContainsSlicesDirectoryCheckAsPrimarySignal(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, ".greenlight/slices/") {
		t.Error("init.md missing .greenlight/slices/ directory check as primary detection signal")
	}
}

func TestInitMd_C80_ContainsStateMdFallbackCheck(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "STATE.md") {
		t.Error("init.md missing STATE.md fallback check as secondary detection signal")
	}
}

func TestInitMd_C80_ContainsNoStateFoundHandling(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "NoStateFound") {
		t.Error("init.md missing NoStateFound error handling instruction")
	}
}

func TestInitMd_C80_ContainsFilePerSliceFormatMention(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "file-per-slice") {
		t.Error("init.md missing 'file-per-slice' format mention")
	}
}

func TestInitMd_C80_ContainsLegacyFormatMention(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "legacy") {
		t.Error("init.md missing 'legacy' format mention")
	}
}

func TestInitMd_C80_ContainsReadOnlyDetectionInstruction(t *testing.T) {
	doc := readInitMd(t)
	hasReadOnly := strings.Contains(doc, "read-only") || strings.Contains(doc, "does not modify") || strings.Contains(doc, "Detection is read-only")
	if !hasReadOnly {
		t.Error("init.md missing instruction that state detection is read-only (no files modified)")
	}
}

func TestInitMd_C80_ContainsStateFormatMdReference(t *testing.T) {
	doc := readInitMd(t)
	if !strings.Contains(doc, "state-format.md") {
		t.Error("init.md missing reference to state-format.md as source of truth for detection logic")
	}
}
