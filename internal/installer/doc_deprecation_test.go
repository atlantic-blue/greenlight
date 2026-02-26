package installer_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/atlantic-blue/greenlight/internal/installer"
)

// S-26 Tests: Documentation and Deprecation
// Covers C-72 (CLAUDEmdVerificationTierRule), C-73 (CheckpointProtocolAcceptanceType),
// and C-74 (ManifestVerificationTiersUpdate).

// =============================================================================
// C-72 Tests: CLAUDEmdVerificationTierRule
// Verifies that src/CLAUDE.md contains the Verification Tiers subsection
// positioned correctly within Code Quality Constraints, with 4 bullet points
// and the required content per contract C-72.
// =============================================================================

func TestCLAUDEmd_ContainsVerificationTiersSubsection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/CLAUDE.md"))
	if err != nil {
		t.Fatalf("failed to read src/CLAUDE.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "Verification Tiers") {
		t.Error("src/CLAUDE.md missing '### Verification Tiers' subsection")
	}
}

func TestCLAUDEmd_VerificationTiersMentionsAutoTier(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/CLAUDE.md"))
	if err != nil {
		t.Fatalf("failed to read src/CLAUDE.md: %v", err)
	}

	doc := string(content)

	tiersPos := strings.Index(doc, "Verification Tiers")
	if tiersPos == -1 {
		t.Fatal("src/CLAUDE.md missing 'Verification Tiers' subsection")
	}

	sectionContent := doc[tiersPos:]

	if !strings.Contains(sectionContent, "auto") {
		t.Error("src/CLAUDE.md Verification Tiers section missing 'auto' tier mention")
	}
}

func TestCLAUDEmd_VerificationTiersMentionsVerifyTier(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/CLAUDE.md"))
	if err != nil {
		t.Fatalf("failed to read src/CLAUDE.md: %v", err)
	}

	doc := string(content)

	tiersPos := strings.Index(doc, "Verification Tiers")
	if tiersPos == -1 {
		t.Fatal("src/CLAUDE.md missing 'Verification Tiers' subsection")
	}

	sectionContent := doc[tiersPos:]

	if !strings.Contains(sectionContent, "verify") {
		t.Error("src/CLAUDE.md Verification Tiers section missing 'verify' tier mention")
	}
}

func TestCLAUDEmd_VerificationTiersMentionsDefault(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/CLAUDE.md"))
	if err != nil {
		t.Fatalf("failed to read src/CLAUDE.md: %v", err)
	}

	doc := string(content)

	tiersPos := strings.Index(doc, "Verification Tiers")
	if tiersPos == -1 {
		t.Fatal("src/CLAUDE.md missing 'Verification Tiers' subsection")
	}

	sectionContent := doc[tiersPos:]

	if !strings.Contains(sectionContent, "default") && !strings.Contains(sectionContent, "Default") {
		t.Error("src/CLAUDE.md Verification Tiers section missing explicit statement of the default tier")
	}
}

func TestCLAUDEmd_VerificationTiersReferencesVerificationTiersMd(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/CLAUDE.md"))
	if err != nil {
		t.Fatalf("failed to read src/CLAUDE.md: %v", err)
	}

	doc := string(content)

	tiersPos := strings.Index(doc, "Verification Tiers")
	if tiersPos == -1 {
		t.Fatal("src/CLAUDE.md missing 'Verification Tiers' subsection")
	}

	sectionContent := doc[tiersPos:]

	if !strings.Contains(sectionContent, "references/verification-tiers.md") {
		t.Error("src/CLAUDE.md Verification Tiers section missing reference to 'references/verification-tiers.md'")
	}
}

func TestCLAUDEmd_VerificationTiersMentionsTestWriterRejectionRouting(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/CLAUDE.md"))
	if err != nil {
		t.Fatalf("failed to read src/CLAUDE.md: %v", err)
	}

	doc := string(content)

	tiersPos := strings.Index(doc, "Verification Tiers")
	if tiersPos == -1 {
		t.Fatal("src/CLAUDE.md missing 'Verification Tiers' subsection")
	}

	sectionContent := doc[tiersPos:]

	hasTestWriterRouting := strings.Contains(sectionContent, "test writer") ||
		strings.Contains(sectionContent, "gl-test-writer")

	if !hasTestWriterRouting {
		t.Error("src/CLAUDE.md Verification Tiers section missing 'test writer' rejection routing reference")
	}
}

func TestCLAUDEmd_VerificationTiersSectionPositionedAfterCircuitBreaker(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/CLAUDE.md"))
	if err != nil {
		t.Fatalf("failed to read src/CLAUDE.md: %v", err)
	}

	doc := string(content)

	circuitBreakerPos := strings.Index(doc, "Circuit Breaker")
	if circuitBreakerPos == -1 {
		t.Fatal("src/CLAUDE.md missing 'Circuit Breaker' section")
	}

	verificationTiersPos := strings.Index(doc, "Verification Tiers")
	if verificationTiersPos == -1 {
		t.Fatal("src/CLAUDE.md missing 'Verification Tiers' section")
	}

	if verificationTiersPos <= circuitBreakerPos {
		t.Errorf("Verification Tiers (pos %d) must appear AFTER Circuit Breaker (pos %d)", verificationTiersPos, circuitBreakerPos)
	}
}

func TestCLAUDEmd_VerificationTiersSectionPositionedBeforeLoggingAndObservability(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/CLAUDE.md"))
	if err != nil {
		t.Fatalf("failed to read src/CLAUDE.md: %v", err)
	}

	doc := string(content)

	verificationTiersPos := strings.Index(doc, "Verification Tiers")
	if verificationTiersPos == -1 {
		t.Fatal("src/CLAUDE.md missing 'Verification Tiers' section")
	}

	loggingPos := strings.Index(doc, "Logging & Observability")
	if loggingPos == -1 {
		t.Fatal("src/CLAUDE.md missing 'Logging & Observability' section")
	}

	if verificationTiersPos >= loggingPos {
		t.Errorf("Verification Tiers (pos %d) must appear BEFORE Logging & Observability (pos %d)", verificationTiersPos, loggingPos)
	}
}

func TestCLAUDEmd_VerificationTiersSectionHasExactlyFourBulletPoints(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/CLAUDE.md"))
	if err != nil {
		t.Fatalf("failed to read src/CLAUDE.md: %v", err)
	}

	doc := string(content)

	tiersPos := strings.Index(doc, "Verification Tiers")
	if tiersPos == -1 {
		t.Fatal("src/CLAUDE.md missing 'Verification Tiers' subsection")
	}

	// Find the next subsection heading ("###") after Verification Tiers to isolate the section
	sectionContent := doc[tiersPos:]
	nextSectionPos := strings.Index(sectionContent[3:], "### ")
	var tiersSection string
	if nextSectionPos == -1 {
		tiersSection = sectionContent
	} else {
		tiersSection = sectionContent[:nextSectionPos+3]
	}

	// Count lines starting with "- " (bullet points)
	lines := strings.Split(tiersSection, "\n")
	bulletCount := 0
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "- ") {
			bulletCount++
		}
	}

	if bulletCount != 4 {
		t.Errorf("src/CLAUDE.md Verification Tiers section must have exactly 4 bullet points, found %d", bulletCount)
	}
}

func TestCLAUDEmd_VerificationTiersIsPhrasedAsHardRule(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/CLAUDE.md"))
	if err != nil {
		t.Fatalf("failed to read src/CLAUDE.md: %v", err)
	}

	doc := string(content)

	tiersPos := strings.Index(doc, "Verification Tiers")
	if tiersPos == -1 {
		t.Fatal("src/CLAUDE.md missing 'Verification Tiers' subsection")
	}

	sectionContent := doc[tiersPos:]

	// Hard rule phrasing uses imperative verbs or definitive statements, not "should" or "recommended"
	// The contract states "This is a hard rule, not a recommendation -- phrased as imperatives"
	hasImperativeOrDefinitive := strings.Contains(sectionContent, "Every contract") ||
		strings.Contains(sectionContent, "determines") ||
		strings.Contains(sectionContent, "routes") ||
		strings.Contains(sectionContent, "Full protocol")

	if !hasImperativeOrDefinitive {
		t.Error("src/CLAUDE.md Verification Tiers section must be phrased as a hard rule using imperative or definitive language, not recommendations")
	}
}

func TestCLAUDEmd_VerificationTiersSectionIsWithinCodeQualityConstraints(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/CLAUDE.md"))
	if err != nil {
		t.Fatalf("failed to read src/CLAUDE.md: %v", err)
	}

	doc := string(content)

	codeQualityPos := strings.Index(doc, "Code Quality Constraints")
	if codeQualityPos == -1 {
		t.Fatal("src/CLAUDE.md missing 'Code Quality Constraints' section")
	}

	verificationTiersPos := strings.Index(doc, "Verification Tiers")
	if verificationTiersPos == -1 {
		t.Fatal("src/CLAUDE.md missing 'Verification Tiers' subsection")
	}

	// Verification Tiers must appear after Code Quality Constraints heading
	if verificationTiersPos <= codeQualityPos {
		t.Errorf("Verification Tiers (pos %d) must appear within 'Code Quality Constraints' (pos %d)", verificationTiersPos, codeQualityPos)
	}

	// Ensure it's within the Code Quality Constraints section (before a top-level ## heading that follows it)
	deviationRulesPos := strings.Index(doc, "## Deviation Rules")
	if deviationRulesPos != -1 && verificationTiersPos >= deviationRulesPos {
		t.Errorf("Verification Tiers (pos %d) must be within 'Code Quality Constraints', but appears after 'Deviation Rules' (pos %d)", verificationTiersPos, deviationRulesPos)
	}
}

func TestCLAUDEmd_ExistingSectionsUnchanged(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/CLAUDE.md"))
	if err != nil {
		t.Fatalf("failed to read src/CLAUDE.md: %v", err)
	}

	doc := string(content)

	// Verify key existing sections are still present and unchanged
	existingSections := []string{
		"Core Principles",
		"Agent Isolation Rules",
		"Code Quality Constraints",
		"Error Handling",
		"Naming",
		"Functions",
		"Security",
		"API Design",
		"Database",
		"Testing",
		"Circuit Breaker",
		"Logging & Observability",
		"File & Project Structure",
		"Git",
		"Performance",
		"Deviation Rules",
		"What NOT To Do",
	}

	for _, section := range existingSections {
		if !strings.Contains(doc, section) {
			t.Errorf("src/CLAUDE.md missing existing section: %q — existing sections must remain unchanged", section)
		}
	}
}

// =============================================================================
// C-73 Tests: CheckpointProtocolAcceptanceType
// Verifies checkpoint-protocol.md, verification-patterns.md, config.md, and
// slice.md Step 9 are updated per contract C-73.
// =============================================================================

// --- checkpoint-protocol.md ---

func TestCheckpointProtocolMd_ContainsAcceptanceCheckpointType(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/checkpoint-protocol.md"))
	if err != nil {
		t.Fatalf("failed to read checkpoint-protocol.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "Acceptance") {
		t.Error("checkpoint-protocol.md missing 'Acceptance' checkpoint type")
	}
}

func TestCheckpointProtocolMd_VisualCheckpointMarkedAsDeprecated(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/checkpoint-protocol.md"))
	if err != nil {
		t.Fatalf("failed to read checkpoint-protocol.md: %v", err)
	}

	doc := string(content)

	// Visual checkpoint must exist but be marked deprecated, not removed
	hasVisual := strings.Contains(doc, "Visual") || strings.Contains(doc, "visual")
	if !hasVisual {
		t.Error("checkpoint-protocol.md must still contain 'Visual' checkpoint type (deprecated, not removed)")
	}

	hasDeprecatedVisual := (strings.Contains(doc, "Visual") || strings.Contains(doc, "visual")) &&
		(strings.Contains(doc, "deprecated") || strings.Contains(doc, "Deprecated"))

	if !hasDeprecatedVisual {
		t.Error("checkpoint-protocol.md Visual checkpoint type must be marked as deprecated")
	}
}

func TestCheckpointProtocolMd_AcceptanceCheckpointAlwaysPauses(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/checkpoint-protocol.md"))
	if err != nil {
		t.Fatalf("failed to read checkpoint-protocol.md: %v", err)
	}

	doc := string(content)

	acceptancePos := strings.Index(doc, "Acceptance")
	if acceptancePos == -1 {
		t.Fatal("checkpoint-protocol.md missing 'Acceptance' checkpoint type")
	}

	sectionContent := doc[acceptancePos:]

	// Acceptance checkpoint must always pause — even in yolo mode
	hasPausesAlways := strings.Contains(sectionContent, "always") &&
		(strings.Contains(sectionContent, "pause") || strings.Contains(sectionContent, "Pause"))

	if !hasPausesAlways {
		t.Error("checkpoint-protocol.md Acceptance checkpoint section must state it always pauses")
	}
}

func TestCheckpointProtocolMd_AcceptancePausesEvenInYoloMode(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/checkpoint-protocol.md"))
	if err != nil {
		t.Fatalf("failed to read checkpoint-protocol.md: %v", err)
	}

	doc := string(content)

	// Must explicitly document that acceptance checkpoints pause even in yolo mode
	hasPausesInYolo := strings.Contains(doc, "yolo") &&
		(strings.Contains(doc, "Acceptance") || strings.Contains(doc, "acceptance")) &&
		(strings.Contains(doc, "pause") || strings.Contains(doc, "always"))

	if !hasPausesInYolo {
		t.Error("checkpoint-protocol.md must document that Acceptance checkpoints pause even in yolo mode")
	}
}

func TestCheckpointProtocolMd_AcceptanceCheckpointTriggerIsVerifyTier(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/checkpoint-protocol.md"))
	if err != nil {
		t.Fatalf("failed to read checkpoint-protocol.md: %v", err)
	}

	doc := string(content)

	acceptancePos := strings.Index(doc, "Acceptance")
	if acceptancePos == -1 {
		t.Fatal("checkpoint-protocol.md missing 'Acceptance' checkpoint type")
	}

	sectionContent := doc[acceptancePos:]

	hasTierTrigger := strings.Contains(sectionContent, "verify") &&
		(strings.Contains(sectionContent, "tier") || strings.Contains(sectionContent, "Tier"))

	if !hasTierTrigger {
		t.Error("checkpoint-protocol.md Acceptance checkpoint section must state the trigger is the verify tier")
	}
}

func TestCheckpointProtocolMd_AcceptanceDescriptionReferencesStep6b(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/checkpoint-protocol.md"))
	if err != nil {
		t.Fatalf("failed to read checkpoint-protocol.md: %v", err)
	}

	doc := string(content)

	acceptancePos := strings.Index(doc, "Acceptance")
	if acceptancePos == -1 {
		t.Fatal("checkpoint-protocol.md missing 'Acceptance' checkpoint type")
	}

	sectionContent := doc[acceptancePos:]

	if !strings.Contains(sectionContent, "Step 6b") {
		t.Error("checkpoint-protocol.md Acceptance checkpoint description must reference Step 6b")
	}
}

func TestCheckpointProtocolMd_ModeTableShowsAcceptanceAlwaysPauses(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/checkpoint-protocol.md"))
	if err != nil {
		t.Fatalf("failed to read checkpoint-protocol.md: %v", err)
	}

	doc := string(content)

	// Mode table must include Acceptance row that shows it always pauses
	modeTablePos := strings.Index(doc, "Mode")
	if modeTablePos == -1 {
		t.Fatal("checkpoint-protocol.md missing mode table")
	}

	modeTableSection := doc[modeTablePos:]

	hasAcceptanceInTable := strings.Contains(modeTableSection, "Acceptance") &&
		(strings.Contains(modeTableSection, "Pause") || strings.Contains(modeTableSection, "pause") || strings.Contains(modeTableSection, "always"))

	if !hasAcceptanceInTable {
		t.Error("checkpoint-protocol.md mode table must include an Acceptance row showing it always pauses")
	}
}

func TestCheckpointProtocolMd_ExistingCheckpointTypesUnchanged(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/checkpoint-protocol.md"))
	if err != nil {
		t.Fatalf("failed to read checkpoint-protocol.md: %v", err)
	}

	doc := string(content)

	// Decision and External Action checkpoint types must remain unchanged
	if !strings.Contains(doc, "Decision") {
		t.Error("checkpoint-protocol.md must retain existing 'Decision' checkpoint type (unchanged)")
	}

	if !strings.Contains(doc, "External Action") && !strings.Contains(doc, "External action") {
		t.Error("checkpoint-protocol.md must retain existing 'External Action' checkpoint type (unchanged)")
	}
}

// --- verification-patterns.md ---

func TestVerificationPatternsMd_ContainsCrossReferenceToVerificationTiersMd(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-patterns.md"))
	if err != nil {
		t.Fatalf("failed to read verification-patterns.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "verification-tiers.md") {
		t.Error("verification-patterns.md must contain a cross-reference to 'verification-tiers.md'")
	}
}

// --- templates/config.md ---

func TestConfigTemplateMd_ContainsVisualCheckpointDeprecationNote(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/config.md"))
	if err != nil {
		t.Fatalf("failed to read templates/config.md: %v", err)
	}

	doc := string(content)

	hasDeprecationNote := strings.Contains(doc, "visual_checkpoint") &&
		(strings.Contains(doc, "deprecated") || strings.Contains(doc, "Deprecated"))

	if !hasDeprecationNote {
		t.Error("templates/config.md must contain a deprecation note for 'visual_checkpoint'")
	}
}

// --- slice.md Step 9 ---

func TestSliceMd_Step9ContainsVisualCheckpointDeprecationWarning(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step9Pos := strings.Index(doc, "Step 9")
	if step9Pos == -1 {
		t.Fatal("slice.md missing 'Step 9' section")
	}

	step9Section := doc[step9Pos:]

	// Step 9 must include a deprecation warning for visual_checkpoint
	hasDeprecationWarning := strings.Contains(step9Section, "visual_checkpoint") &&
		(strings.Contains(step9Section, "deprecated") || strings.Contains(step9Section, "Deprecated") ||
			strings.Contains(step9Section, "warning") || strings.Contains(step9Section, "Warning"))

	if !hasDeprecationWarning {
		t.Error("slice.md Step 9 must contain a deprecation warning for 'visual_checkpoint'")
	}
}

func TestSliceMd_Step9IsNoOpWithDeprecationWarning(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step9Pos := strings.Index(doc, "Step 9")
	if step9Pos == -1 {
		t.Fatal("slice.md missing 'Step 9' section")
	}

	// Find the boundary of Step 9 (next step heading)
	sectionAfterStep9 := doc[step9Pos:]
	step10Pos := strings.Index(sectionAfterStep9, "Step 10")
	var step9Section string
	if step10Pos == -1 {
		step9Section = sectionAfterStep9
	} else {
		step9Section = sectionAfterStep9[:step10Pos]
	}

	// Step 9 must be a no-op — it should indicate the step is superseded or a no-op
	isNoOp := strings.Contains(step9Section, "no-op") ||
		strings.Contains(step9Section, "no op") ||
		strings.Contains(step9Section, "superseded") ||
		strings.Contains(step9Section, "skip") ||
		strings.Contains(step9Section, "Skip") ||
		(strings.Contains(step9Section, "deprecated") && strings.Contains(step9Section, "Step 6b"))

	if !isNoOp {
		t.Error("slice.md Step 9 must be a no-op with deprecation warning (referencing Step 6b as the replacement)")
	}
}

// =============================================================================
// C-74 Tests: ManifestVerificationTiersUpdate
// Verifies that installer.Manifest contains "references/verification-tiers.md"
// with a total count of 38 entries, alphabetically ordered, and CLAUDE.md last.
// =============================================================================

func TestManifest_CountIs35(t *testing.T) {
	got := len(installer.Manifest)
	if got != 38 {
		t.Errorf("installer.Manifest must have 38 entries, got %d", got)
	}
}

func TestManifest_ContainsVerificationTiersMd(t *testing.T) {
	for _, entry := range installer.Manifest {
		if entry == "references/verification-tiers.md" {
			return
		}
	}
	t.Error("installer.Manifest missing 'references/verification-tiers.md'")
}

func TestManifest_VerificationTiersMdIsAlphabeticallyAfterVerificationPatternsMd(t *testing.T) {
	verificationPatternsIdx := -1
	verificationTiersIdx := -1

	for i, entry := range installer.Manifest {
		if entry == "references/verification-patterns.md" {
			verificationPatternsIdx = i
		}
		if entry == "references/verification-tiers.md" {
			verificationTiersIdx = i
		}
	}

	if verificationPatternsIdx == -1 {
		t.Fatal("installer.Manifest missing 'references/verification-patterns.md'")
	}
	if verificationTiersIdx == -1 {
		t.Fatal("installer.Manifest missing 'references/verification-tiers.md'")
	}

	if verificationTiersIdx <= verificationPatternsIdx {
		t.Errorf("'references/verification-tiers.md' (idx %d) must appear AFTER 'references/verification-patterns.md' (idx %d) for alphabetical ordering", verificationTiersIdx, verificationPatternsIdx)
	}
}

func TestManifest_CLAUDEmdRemainsLastEntry(t *testing.T) {
	if len(installer.Manifest) == 0 {
		t.Fatal("installer.Manifest is empty")
	}

	lastEntry := installer.Manifest[len(installer.Manifest)-1]
	if lastEntry != "CLAUDE.md" {
		t.Errorf("installer.Manifest last entry must be 'CLAUDE.md', got %q", lastEntry)
	}
}

func TestManifest_VerificationTiersMdIsWithinReferencesSection(t *testing.T) {
	tiersIdx := -1
	for i, entry := range installer.Manifest {
		if entry == "references/verification-tiers.md" {
			tiersIdx = i
			break
		}
	}

	if tiersIdx == -1 {
		t.Fatal("installer.Manifest missing 'references/verification-tiers.md'")
	}

	if !strings.HasPrefix(installer.Manifest[tiersIdx], "references/") {
		t.Errorf("'references/verification-tiers.md' must be in the references/ section, got: %q", installer.Manifest[tiersIdx])
	}
}

func TestManifest_VerificationTiersMdEntryFileIsNonEmpty(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-tiers.md"))
	if err != nil {
		t.Fatalf("failed to read src/references/verification-tiers.md: %v", err)
	}

	if len(content) == 0 {
		t.Error("src/references/verification-tiers.md must be non-empty")
	}
}

func TestManifest_AllExistingEntriesRetained(t *testing.T) {
	// All 35 previously-existing entries must still be present in the updated manifest of 38
	expectedExisting := []string{
		"agents/gl-architect.md",
		"agents/gl-assessor.md",
		"agents/gl-codebase-mapper.md",
		"agents/gl-debugger.md",
		"agents/gl-designer.md",
		"agents/gl-implementer.md",
		"agents/gl-security.md",
		"agents/gl-test-writer.md",
		"agents/gl-verifier.md",
		"agents/gl-wrapper.md",
		"commands/gl/add-slice.md",
		"commands/gl/assess.md",
		"commands/gl/changelog.md",
		"commands/gl/debug.md",
		"commands/gl/design.md",
		"commands/gl/help.md",
		"commands/gl/init.md",
		"commands/gl/map.md",
		"commands/gl/pause.md",
		"commands/gl/quick.md",
		"commands/gl/resume.md",
		"commands/gl/roadmap.md",
		"commands/gl/settings.md",
		"commands/gl/ship.md",
		"commands/gl/slice.md",
		"commands/gl/status.md",
		"commands/gl/wrap.md",
		"references/checkpoint-protocol.md",
		"references/circuit-breaker.md",
		"references/deviation-rules.md",
		"references/verification-patterns.md",
		"templates/config.md",
		"templates/state.md",
		"CLAUDE.md",
	}

	manifestSet := make(map[string]bool, len(installer.Manifest))
	for _, entry := range installer.Manifest {
		manifestSet[entry] = true
	}

	for _, expected := range expectedExisting {
		if !manifestSet[expected] {
			t.Errorf("installer.Manifest missing previously-existing entry: %q", expected)
		}
	}
}
