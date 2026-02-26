package installer_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// C-87 Tests: CLAUDEmdStateFormatRule
// These tests verify that src/CLAUDE.md contains the state format awareness rule
// as defined in contract C-87.

// C-88 Tests: StateTemplateDocUpdate
// These tests verify that src/templates/state.md documents both state formats
// as defined in contract C-88.

// C-89 Tests: CheckpointProtocolStateUpdate
// These tests verify that src/references/checkpoint-protocol.md contains
// slice file references as defined in contract C-89.

// ---------------------------------------------------------------------------
// Reader helpers
// ---------------------------------------------------------------------------

func readClaudeMdSrc(t *testing.T) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/CLAUDE.md"))
	if err != nil {
		t.Fatalf("failed to read src/CLAUDE.md: %v", err)
	}
	return string(content)
}

func readStateTemplateMd(t *testing.T) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/state.md"))
	if err != nil {
		t.Fatalf("failed to read src/templates/state.md: %v", err)
	}
	return string(content)
}

func readCheckpointProtocolMd(t *testing.T) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/checkpoint-protocol.md"))
	if err != nil {
		t.Fatalf("failed to read src/references/checkpoint-protocol.md: %v", err)
	}
	return string(content)
}

// ---------------------------------------------------------------------------
// C-87: CLAUDEmdStateFormatRule — src/CLAUDE.md
// ---------------------------------------------------------------------------

func TestClaudeMdSrc_C87_ContainsStateFormatHeading(t *testing.T) {
	doc := readClaudeMdSrc(t)

	if !strings.Contains(doc, "### State Format") {
		t.Error("src/CLAUDE.md missing '### State Format' heading required by C-87")
	}
}

func TestClaudeMdSrc_C87_ContainsGreenlightSlicesReference(t *testing.T) {
	doc := readClaudeMdSrc(t)

	if !strings.Contains(doc, ".greenlight/slices/") {
		t.Error("src/CLAUDE.md missing '.greenlight/slices/' reference required by C-87")
	}
}

func TestClaudeMdSrc_C87_ContainsSourceOfTruthForSliceFiles(t *testing.T) {
	doc := readClaudeMdSrc(t)

	if !strings.Contains(doc, "source of truth") {
		t.Error("src/CLAUDE.md missing 'source of truth' reference for slice files required by C-87")
	}
}

func TestClaudeMdSrc_C87_ContainsGeneratedOrDoNotWriteForStateMd(t *testing.T) {
	doc := readClaudeMdSrc(t)

	if !strings.Contains(doc, "generated") && !strings.Contains(doc, "do not write") {
		t.Error("src/CLAUDE.md missing 'generated' or 'do not write' guidance for STATE.md in file-per-slice mode required by C-87")
	}
}

func TestClaudeMdSrc_C87_ContainsStateFormatMdCrossReference(t *testing.T) {
	doc := readClaudeMdSrc(t)

	if !strings.Contains(doc, "state-format.md") && !strings.Contains(doc, "references/state-format.md") {
		t.Error("src/CLAUDE.md missing 'state-format.md' cross-reference required by C-87")
	}
}

func TestClaudeMdSrc_C87_ContainsHardRuleToCheckSlicesFirst(t *testing.T) {
	doc := readClaudeMdSrc(t)

	hasCheck := strings.Contains(doc, "check") || strings.Contains(doc, "Check")
	hasSlices := strings.Contains(doc, "slices/")

	if !hasCheck || !hasSlices {
		t.Error("src/CLAUDE.md missing hard rule to check slices/ before reading STATE.md required by C-87")
	}
}

func TestClaudeMdSrc_C87_RuleAppliesToAllAgents(t *testing.T) {
	doc := readClaudeMdSrc(t)

	if !strings.Contains(doc, "all agents") && !strings.Contains(doc, "agents") {
		t.Error("src/CLAUDE.md missing statement that state format rule applies to agents required by C-87")
	}
}

// ---------------------------------------------------------------------------
// C-88: StateTemplateDocUpdate — src/templates/state.md
// ---------------------------------------------------------------------------

func TestStateTemplateMd_C88_ContainsStateFormatDetectionSection(t *testing.T) {
	doc := readStateTemplateMd(t)

	if !strings.Contains(doc, "State Format Detection") && !strings.Contains(doc, "format detection") {
		t.Error("src/templates/state.md missing 'State Format Detection' section required by C-88")
	}
}

func TestStateTemplateMd_C88_ContainsFilePerSliceReference(t *testing.T) {
	doc := readStateTemplateMd(t)

	if !strings.Contains(doc, "file-per-slice") {
		t.Error("src/templates/state.md missing 'file-per-slice' format reference required by C-88")
	}
}

func TestStateTemplateMd_C88_ContainsLegacyFormatReference(t *testing.T) {
	doc := readStateTemplateMd(t)

	if !strings.Contains(doc, "legacy") {
		t.Error("src/templates/state.md missing 'legacy' format reference required by C-88")
	}
}

func TestStateTemplateMd_C88_ContainsGeneratedReference(t *testing.T) {
	doc := readStateTemplateMd(t)

	if !strings.Contains(doc, "generated") && !strings.Contains(doc, "Generated") {
		t.Error("src/templates/state.md missing 'generated' reference for STATE.md required by C-88")
	}
}

func TestStateTemplateMd_C88_ContainsHeaderCommentOrDoNotEdit(t *testing.T) {
	doc := readStateTemplateMd(t)

	if !strings.Contains(doc, "GENERATED") && !strings.Contains(doc, "do not edit") && !strings.Contains(doc, "overwritten") {
		t.Error("src/templates/state.md missing header comment or 'do not edit'/'overwritten' reference required by C-88")
	}
}

func TestStateTemplateMd_C88_ContainsMigrateStateReference(t *testing.T) {
	doc := readStateTemplateMd(t)

	if !strings.Contains(doc, "migrate-state") && !strings.Contains(doc, "/gl:migrate-state") {
		t.Error("src/templates/state.md missing 'migrate-state' migration reference required by C-88")
	}
}

func TestStateTemplateMd_C88_ContainsExplicitOrOneWayMigrationDescription(t *testing.T) {
	doc := readStateTemplateMd(t)

	if !strings.Contains(doc, "explicit") && !strings.Contains(doc, "one-way") {
		t.Error("src/templates/state.md missing 'explicit' or 'one-way' migration description required by C-88")
	}
}

func TestStateTemplateMd_C88_ContainsStateFormatMdCrossReference(t *testing.T) {
	doc := readStateTemplateMd(t)

	if !strings.Contains(doc, "state-format.md") {
		t.Error("src/templates/state.md missing 'state-format.md' cross-reference required by C-88")
	}
}

func TestStateTemplateMd_C88_ContainsSliceStateMdCrossReference(t *testing.T) {
	doc := readStateTemplateMd(t)

	if !strings.Contains(doc, "slice-state.md") {
		t.Error("src/templates/state.md missing 'slice-state.md' cross-reference required by C-88")
	}
}

func TestStateTemplateMd_C88_DirectsNewProjectsToFilePerSlice(t *testing.T) {
	doc := readStateTemplateMd(t)

	if !strings.Contains(doc, "new project") && !strings.Contains(doc, "recommended") {
		t.Error("src/templates/state.md missing guidance directing new projects to file-per-slice format required by C-88")
	}
}

// ---------------------------------------------------------------------------
// C-89: CheckpointProtocolStateUpdate — src/references/checkpoint-protocol.md
// ---------------------------------------------------------------------------

func TestCheckpointProtocolMd_C89_ContainsSlicesDirectoryOrSliceFileReference(t *testing.T) {
	doc := readCheckpointProtocolMd(t)

	if !strings.Contains(doc, ".greenlight/slices/") && !strings.Contains(doc, "slice file") {
		t.Error("src/references/checkpoint-protocol.md missing '.greenlight/slices/' or 'slice file' reference required by C-89")
	}
}

func TestCheckpointProtocolMd_C89_ContainsFilePerSliceReference(t *testing.T) {
	doc := readCheckpointProtocolMd(t)

	if !strings.Contains(doc, "file-per-slice") {
		t.Error("src/references/checkpoint-protocol.md missing 'file-per-slice' reference required by C-89")
	}
}

func TestCheckpointProtocolMd_C89_ContainsStateFormatDetectionReference(t *testing.T) {
	doc := readCheckpointProtocolMd(t)

	if !strings.Contains(doc, "state format") && !strings.Contains(doc, "state-format.md") {
		t.Error("src/references/checkpoint-protocol.md missing 'state format' or 'state-format.md' detection reference required by C-89")
	}
}

func TestCheckpointProtocolMd_C89_WorksWithBothFormats(t *testing.T) {
	doc := readCheckpointProtocolMd(t)

	if !strings.Contains(doc, "both formats") && !strings.Contains(doc, "legacy") {
		t.Error("src/references/checkpoint-protocol.md missing 'both formats' or 'legacy' reference — must work with both formats as required by C-89")
	}
}

func TestCheckpointProtocolMd_C89_ExistingContentPreserved(t *testing.T) {
	doc := readCheckpointProtocolMd(t)

	if !strings.Contains(doc, "checkpoint") && !strings.Contains(doc, "rollback") {
		t.Error("src/references/checkpoint-protocol.md existing content appears missing — 'checkpoint' or 'rollback' term not found; existing protocol logic must be preserved as required by C-89")
	}
}
