package installer_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// C-64 Tests: VerificationTiersProtocol
// These tests verify that src/references/verification-tiers.md exists and contains
// all required sections as defined in contract C-64.

func TestVerificationTiersMd_Exists(t *testing.T) {
	path := filepath.Join(projectRoot(), "src/references/verification-tiers.md")
	_, err := os.Stat(path)
	if err != nil {
		t.Fatalf("verification-tiers.md does not exist at %s: %v", path, err)
	}
}

func TestVerificationTiersMd_ContainsTierDefinitionsSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-tiers.md"))
	if err != nil {
		t.Fatalf("failed to read verification-tiers.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "Tier Definitions") {
		t.Error("verification-tiers.md missing 'Tier Definitions' section")
	}
}

func TestVerificationTiersMd_TierDefinitionsMentionsAutoTier(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-tiers.md"))
	if err != nil {
		t.Fatalf("failed to read verification-tiers.md: %v", err)
	}

	doc := string(content)

	tierDefinitionsPos := strings.Index(doc, "Tier Definitions")
	if tierDefinitionsPos == -1 {
		t.Fatal("verification-tiers.md missing 'Tier Definitions' section")
	}

	sectionContent := doc[tierDefinitionsPos:]

	if !strings.Contains(sectionContent, "auto") {
		t.Error("verification-tiers.md 'Tier Definitions' section missing 'auto' tier definition")
	}
}

func TestVerificationTiersMd_TierDefinitionsMentionsVerifyTier(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-tiers.md"))
	if err != nil {
		t.Fatalf("failed to read verification-tiers.md: %v", err)
	}

	doc := string(content)

	tierDefinitionsPos := strings.Index(doc, "Tier Definitions")
	if tierDefinitionsPos == -1 {
		t.Fatal("verification-tiers.md missing 'Tier Definitions' section")
	}

	sectionContent := doc[tierDefinitionsPos:]

	if !strings.Contains(sectionContent, "verify") {
		t.Error("verification-tiers.md 'Tier Definitions' section missing 'verify' tier definition")
	}
}

func TestVerificationTiersMd_DocumentsOnlyTwoTiers(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-tiers.md"))
	if err != nil {
		t.Fatalf("failed to read verification-tiers.md: %v", err)
	}

	doc := string(content)

	// The contract invariant states: Two tiers only: auto and verify (no third tier).
	// The document must explicitly state only two tiers exist.
	hasTwoTiersOnly := strings.Contains(doc, "two tiers") ||
		strings.Contains(doc, "Two tiers") ||
		strings.Contains(doc, "only two") ||
		strings.Contains(doc, "only: auto") ||
		strings.Contains(doc, "auto and verify")

	if !hasTwoTiersOnly {
		t.Error("verification-tiers.md missing documentation that only two tiers exist (auto and verify)")
	}
}

func TestVerificationTiersMd_DocumentsDefaultAsVerify(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-tiers.md"))
	if err != nil {
		t.Fatalf("failed to read verification-tiers.md: %v", err)
	}

	doc := string(content)

	hasDefaultVerify := strings.Contains(doc, "default: verify") ||
		strings.Contains(doc, "default is verify") ||
		strings.Contains(doc, "Default: verify") ||
		strings.Contains(doc, "Default is verify") ||
		strings.Contains(doc, "(default: verify)") ||
		strings.Contains(doc, "defaults to verify")

	if !hasDefaultVerify {
		t.Error("verification-tiers.md missing documentation that the default tier is 'verify'")
	}
}

func TestVerificationTiersMd_ContainsTierResolutionSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-tiers.md"))
	if err != nil {
		t.Fatalf("failed to read verification-tiers.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "Tier Resolution") {
		t.Error("verification-tiers.md missing 'Tier Resolution' section")
	}
}

func TestVerificationTiersMd_TierResolutionDocumentsVerifyWinsOverAuto(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-tiers.md"))
	if err != nil {
		t.Fatalf("failed to read verification-tiers.md: %v", err)
	}

	doc := string(content)

	tierResolutionPos := strings.Index(doc, "Tier Resolution")
	if tierResolutionPos == -1 {
		t.Fatal("verification-tiers.md missing 'Tier Resolution' section")
	}

	sectionContent := doc[tierResolutionPos:]

	hasVerifyWins := strings.Contains(sectionContent, "verify > auto") ||
		strings.Contains(sectionContent, "highest wins") ||
		strings.Contains(sectionContent, "verify wins") ||
		strings.Contains(sectionContent, "verify takes precedence")

	if !hasVerifyWins {
		t.Error("verification-tiers.md 'Tier Resolution' section missing rule that verify takes precedence over auto (verify > auto)")
	}
}

func TestVerificationTiersMd_TierResolutionDocumentsOneCheckpointPerSlice(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-tiers.md"))
	if err != nil {
		t.Fatalf("failed to read verification-tiers.md: %v", err)
	}

	doc := string(content)

	tierResolutionPos := strings.Index(doc, "Tier Resolution")
	if tierResolutionPos == -1 {
		t.Fatal("verification-tiers.md missing 'Tier Resolution' section")
	}

	sectionContent := doc[tierResolutionPos:]

	hasOneCheckpoint := strings.Contains(sectionContent, "one checkpoint") ||
		strings.Contains(sectionContent, "single checkpoint") ||
		strings.Contains(sectionContent, "one checkpoint per slice") ||
		strings.Contains(sectionContent, "per-slice")

	if !hasOneCheckpoint {
		t.Error("verification-tiers.md 'Tier Resolution' section missing documentation of one checkpoint per slice")
	}
}

func TestVerificationTiersMd_ContainsVerifyCheckpointFormatSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-tiers.md"))
	if err != nil {
		t.Fatalf("failed to read verification-tiers.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "Verify Checkpoint Format") {
		t.Error("verification-tiers.md missing 'Verify Checkpoint Format' section")
	}
}

func TestVerificationTiersMd_CheckpointFormatContainsAllTestsPassingHeader(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-tiers.md"))
	if err != nil {
		t.Fatalf("failed to read verification-tiers.md: %v", err)
	}

	doc := string(content)

	checkpointFormatPos := strings.Index(doc, "Verify Checkpoint Format")
	if checkpointFormatPos == -1 {
		t.Fatal("verification-tiers.md missing 'Verify Checkpoint Format' section")
	}

	sectionContent := doc[checkpointFormatPos:]

	if !strings.Contains(sectionContent, "ALL TESTS PASSING") {
		t.Error("verification-tiers.md 'Verify Checkpoint Format' section missing 'ALL TESTS PASSING' header")
	}
}

func TestVerificationTiersMd_ContainsRejectionFlowSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-tiers.md"))
	if err != nil {
		t.Fatalf("failed to read verification-tiers.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "Rejection Flow") {
		t.Error("verification-tiers.md missing 'Rejection Flow' section")
	}
}

func TestVerificationTiersMd_RejectionFlowContainsThreeOptions(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-tiers.md"))
	if err != nil {
		t.Fatalf("failed to read verification-tiers.md: %v", err)
	}

	doc := string(content)

	rejectionFlowPos := strings.Index(doc, "Rejection Flow")
	if rejectionFlowPos == -1 {
		t.Fatal("verification-tiers.md missing 'Rejection Flow' section")
	}

	sectionContent := doc[rejectionFlowPos:]

	hasTightenTests := strings.Contains(sectionContent, "tighten tests") ||
		strings.Contains(sectionContent, "Tighten tests")
	hasReviseContract := strings.Contains(sectionContent, "revise contract") ||
		strings.Contains(sectionContent, "Revise contract")
	hasProvideDetail := strings.Contains(sectionContent, "provide more detail") ||
		strings.Contains(sectionContent, "more detail")

	if !hasTightenTests {
		t.Error("verification-tiers.md 'Rejection Flow' section missing 'tighten tests' option")
	}
	if !hasReviseContract {
		t.Error("verification-tiers.md 'Rejection Flow' section missing 'revise contract' option")
	}
	if !hasProvideDetail {
		t.Error("verification-tiers.md 'Rejection Flow' section missing 'provide more detail' option")
	}
}

func TestVerificationTiersMd_ContainsRejectionCounterSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-tiers.md"))
	if err != nil {
		t.Fatalf("failed to read verification-tiers.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "Rejection Counter") {
		t.Error("verification-tiers.md missing 'Rejection Counter' section")
	}
}

func TestVerificationTiersMd_RejectionCounterDocumentsEscalationAtThree(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-tiers.md"))
	if err != nil {
		t.Fatalf("failed to read verification-tiers.md: %v", err)
	}

	doc := string(content)

	rejectionCounterPos := strings.Index(doc, "Rejection Counter")
	if rejectionCounterPos == -1 {
		t.Fatal("verification-tiers.md missing 'Rejection Counter' section")
	}

	sectionContent := doc[rejectionCounterPos:]

	hasEscalationAt3 := strings.Contains(sectionContent, "escalation at 3") ||
		strings.Contains(sectionContent, "Escalation at 3") ||
		strings.Contains(sectionContent, "escalate at 3") ||
		strings.Contains(sectionContent, "3 rejections") ||
		strings.Contains(sectionContent, "after 3")

	if !hasEscalationAt3 {
		t.Error("verification-tiers.md 'Rejection Counter' section missing escalation threshold of 3")
	}
}

func TestVerificationTiersMd_RejectionCounterDocumentsPerSliceCount(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-tiers.md"))
	if err != nil {
		t.Fatalf("failed to read verification-tiers.md: %v", err)
	}

	doc := string(content)

	rejectionCounterPos := strings.Index(doc, "Rejection Counter")
	if rejectionCounterPos == -1 {
		t.Fatal("verification-tiers.md missing 'Rejection Counter' section")
	}

	sectionContent := doc[rejectionCounterPos:]

	hasPerSlice := strings.Contains(sectionContent, "per-slice") ||
		strings.Contains(sectionContent, "per slice") ||
		strings.Contains(sectionContent, "Per-slice") ||
		strings.Contains(sectionContent, "Per slice")

	if !hasPerSlice {
		t.Error("verification-tiers.md 'Rejection Counter' section missing documentation that counter is per-slice")
	}
}

func TestVerificationTiersMd_ContainsAgentIsolationSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-tiers.md"))
	if err != nil {
		t.Fatalf("failed to read verification-tiers.md: %v", err)
	}

	doc := string(content)

	hasAgentIsolation := strings.Contains(doc, "Agent Isolation") ||
		strings.Contains(doc, "agent isolation")

	if !hasAgentIsolation {
		t.Error("verification-tiers.md missing 'Agent Isolation' section")
	}
}

func TestVerificationTiersMd_AgentIsolationMentionsBehavioralFeedbackOnly(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-tiers.md"))
	if err != nil {
		t.Fatalf("failed to read verification-tiers.md: %v", err)
	}

	doc := string(content)

	agentIsolationPos := strings.Index(doc, "Agent Isolation")
	if agentIsolationPos == -1 {
		t.Fatal("verification-tiers.md missing 'Agent Isolation' section")
	}

	sectionContent := doc[agentIsolationPos:]

	hasBehavioralFeedback := strings.Contains(sectionContent, "behavioral feedback") ||
		strings.Contains(sectionContent, "behavioural feedback") ||
		strings.Contains(sectionContent, "behavioral feedback only") ||
		strings.Contains(sectionContent, "behaviour only") ||
		strings.Contains(sectionContent, "behavioral only")

	if !hasBehavioralFeedback {
		t.Error("verification-tiers.md 'Agent Isolation' section missing documentation that test writer receives behavioral feedback only (not implementation code)")
	}
}

func TestVerificationTiersMd_ContainsBackwardCompatibilitySection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-tiers.md"))
	if err != nil {
		t.Fatalf("failed to read verification-tiers.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "Backward Compatibility") {
		t.Error("verification-tiers.md missing 'Backward Compatibility' section")
	}
}

func TestVerificationTiersMd_BackwardCompatibilityMentionsVisualCheckpointDeprecation(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-tiers.md"))
	if err != nil {
		t.Fatalf("failed to read verification-tiers.md: %v", err)
	}

	doc := string(content)

	backwardCompatPos := strings.Index(doc, "Backward Compatibility")
	if backwardCompatPos == -1 {
		t.Fatal("verification-tiers.md missing 'Backward Compatibility' section")
	}

	sectionContent := doc[backwardCompatPos:]

	hasVisualCheckpointDeprecation := strings.Contains(sectionContent, "visual_checkpoint") &&
		(strings.Contains(sectionContent, "deprecated") || strings.Contains(sectionContent, "deprecat"))

	if !hasVisualCheckpointDeprecation {
		t.Error("verification-tiers.md 'Backward Compatibility' section missing documentation that visual_checkpoint is deprecated")
	}
}

func TestVerificationTiersMd_DocumentsCheckpointsPauseInYoloMode(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-tiers.md"))
	if err != nil {
		t.Fatalf("failed to read verification-tiers.md: %v", err)
	}

	doc := string(content)

	hasPausesInYolo := strings.Contains(doc, "yolo") &&
		(strings.Contains(doc, "pause") || strings.Contains(doc, "blocking"))

	if !hasPausesInYolo {
		t.Error("verification-tiers.md missing documentation that acceptance checkpoints pause even in yolo mode")
	}
}

// C-65 Tests: VerificationTierGate
// These tests verify that src/commands/gl/slice.md contains Step 6b
// as defined in contract C-65.

func TestSliceMd_ContainsStep6b(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "Step 6b") {
		t.Error("slice.md missing 'Step 6b' section — verification tier gate is required")
	}
}

func TestSliceMd_Step6bPositionedAfterStep6AndBeforeStep7(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6Pos := strings.Index(doc, "Step 6:")
	if step6Pos == -1 {
		t.Fatal("slice.md missing 'Step 6:' section")
	}

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	step7Pos := strings.Index(doc, "Step 7:")
	if step7Pos == -1 {
		t.Fatal("slice.md missing 'Step 7:' section")
	}

	if step6bPos <= step6Pos {
		t.Errorf("Step 6b (pos %d) must appear AFTER Step 6 (pos %d)", step6bPos, step6Pos)
	}

	if step6bPos >= step7Pos {
		t.Errorf("Step 6b (pos %d) must appear BEFORE Step 7 (pos %d)", step6bPos, step7Pos)
	}
}

func TestSliceMd_Step6bContainsTierResolutionLogic(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	step7Pos := strings.Index(doc, "Step 7:")
	if step7Pos == -1 {
		t.Fatal("slice.md missing 'Step 7:' section")
	}

	step6bSection := doc[step6bPos:step7Pos]

	hasTierResolution := strings.Contains(step6bSection, "verify > auto") ||
		strings.Contains(step6bSection, "highest wins") ||
		strings.Contains(step6bSection, "effective tier") ||
		strings.Contains(step6bSection, "Effective tier")

	if !hasTierResolution {
		t.Error("slice.md Step 6b section missing tier resolution logic (verify > auto or equivalent)")
	}
}

func TestSliceMd_Step6bContainsAutoTierSkipBehavior(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	step7Pos := strings.Index(doc, "Step 7:")
	if step7Pos == -1 {
		t.Fatal("slice.md missing 'Step 7:' section")
	}

	step6bSection := doc[step6bPos:step7Pos]

	if !strings.Contains(step6bSection, "Skipping acceptance checkpoint") {
		t.Error("slice.md Step 6b section missing 'Skipping acceptance checkpoint' message for auto tier")
	}
}

func TestSliceMd_Step6bContainsVerifyTierCheckpointBehavior(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	step7Pos := strings.Index(doc, "Step 7:")
	if step7Pos == -1 {
		t.Fatal("slice.md missing 'Step 7:' section")
	}

	step6bSection := doc[step6bPos:step7Pos]

	hasVerifyCheckpoint := strings.Contains(step6bSection, "checkpoint") &&
		strings.Contains(step6bSection, "verify")

	if !hasVerifyCheckpoint {
		t.Error("slice.md Step 6b section missing verify-tier checkpoint behavior")
	}
}

func TestSliceMd_Step6bContainsVisualCheckpointDeprecationHandling(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	step7Pos := strings.Index(doc, "Step 7:")
	if step7Pos == -1 {
		t.Fatal("slice.md missing 'Step 7:' section")
	}

	step6bSection := doc[step6bPos:step7Pos]

	hasDeprecation := strings.Contains(step6bSection, "visual_checkpoint") &&
		(strings.Contains(step6bSection, "deprecated") || strings.Contains(step6bSection, "deprecat") || strings.Contains(step6bSection, "warning"))

	if !hasDeprecation {
		t.Error("slice.md Step 6b section missing visual_checkpoint deprecation warning handling")
	}
}

func TestSliceMd_Step6bDocumentsGateIsBlocking(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	step7Pos := strings.Index(doc, "Step 7:")
	if step7Pos == -1 {
		t.Fatal("slice.md missing 'Step 7:' section")
	}

	step6bSection := doc[step6bPos:step7Pos]

	hasBlockingBehavior := strings.Contains(step6bSection, "blocking") ||
		strings.Contains(step6bSection, "pause") ||
		strings.Contains(step6bSection, "Wait") ||
		strings.Contains(step6bSection, "wait")

	if !hasBlockingBehavior {
		t.Error("slice.md Step 6b section missing documentation that the gate is blocking")
	}
}

func TestSliceMd_Step6bDocumentsPausesInYoloMode(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	step7Pos := strings.Index(doc, "Step 7:")
	if step7Pos == -1 {
		t.Fatal("slice.md missing 'Step 7:' section")
	}

	step6bSection := doc[step6bPos:step7Pos]

	hasPausesInYolo := strings.Contains(step6bSection, "yolo") &&
		(strings.Contains(step6bSection, "pause") || strings.Contains(step6bSection, "blocking") || strings.Contains(step6bSection, "even"))

	if !hasPausesInYolo {
		t.Error("slice.md Step 6b section missing documentation that acceptance checkpoints pause even in yolo mode")
	}
}

func TestSliceMd_Step6bDocumentsRejectionFlowOnNonApprovedResponse(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	step7Pos := strings.Index(doc, "Step 7:")
	if step7Pos == -1 {
		t.Fatal("slice.md missing 'Step 7:' section")
	}

	step6bSection := doc[step6bPos:step7Pos]

	hasRejectionFlow := strings.Contains(step6bSection, "rejection") ||
		strings.Contains(step6bSection, "Rejection") ||
		strings.Contains(step6bSection, "rejected") ||
		strings.Contains(step6bSection, "reject")

	if !hasRejectionFlow {
		t.Error("slice.md Step 6b section missing rejection flow documentation for non-approved responses")
	}
}

// C-66 Tests: VerifyCheckpointPresentation
// These tests verify that src/commands/gl/slice.md contains the acceptance checkpoint
// format as defined in contract C-66.

func TestSliceMd_ContainsAllTestsPassingCheckpointHeader(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "ALL TESTS PASSING") {
		t.Error("slice.md missing 'ALL TESTS PASSING' checkpoint header format as required by C-66")
	}
}

func TestSliceMd_AllTestsPassingHeaderAppearsInStep6bSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionAfterStep6b := doc[step6bPos:]

	if !strings.Contains(sectionAfterStep6b, "ALL TESTS PASSING") {
		t.Error("slice.md 'ALL TESTS PASSING' checkpoint header must appear in or after Step 6b section")
	}
}

func TestSliceMd_CheckpointContainsAcceptanceCriteriaChecklist(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionAfterStep6b := doc[step6bPos:]

	// Acceptance criteria must be presented as unchecked checkboxes [ ]
	hasCheckboxFormat := strings.Contains(sectionAfterStep6b, "[ ]")
	hasCriteriaLabel := strings.Contains(sectionAfterStep6b, "Acceptance criteria") ||
		strings.Contains(sectionAfterStep6b, "acceptance criteria")

	if !hasCheckboxFormat {
		t.Error("slice.md missing unchecked checkbox format '[ ]' for acceptance criteria in checkpoint")
	}
	if !hasCriteriaLabel {
		t.Error("slice.md checkpoint missing 'Acceptance criteria' label")
	}
}

func TestSliceMd_CheckpointContainsStepsAsNumberedList(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionAfterStep6b := doc[step6bPos:]

	hasStepsLabel := strings.Contains(sectionAfterStep6b, "Steps to verify") ||
		strings.Contains(sectionAfterStep6b, "steps to verify")

	if !hasStepsLabel {
		t.Error("slice.md checkpoint missing 'Steps to verify' label for numbered verification steps")
	}
}

func TestSliceMd_CheckpointContainsThreeResponseOptions(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionAfterStep6b := doc[step6bPos:]

	hasYesOption := strings.Contains(sectionAfterStep6b, "Yes") &&
		strings.Contains(sectionAfterStep6b, "mark complete")
	hasNoOption := strings.Contains(sectionAfterStep6b, "No") &&
		(strings.Contains(sectionAfterStep6b, "describe") || strings.Contains(sectionAfterStep6b, "wrong"))
	hasPartialOption := strings.Contains(sectionAfterStep6b, "Partially") ||
		strings.Contains(sectionAfterStep6b, "partially")

	if !hasYesOption {
		t.Error("slice.md checkpoint missing 'Yes -- mark complete' response option")
	}
	if !hasNoOption {
		t.Error("slice.md checkpoint missing 'No -- describe what's wrong' response option")
	}
	if !hasPartialOption {
		t.Error("slice.md checkpoint missing 'Partially' response option")
	}
}

func TestSliceMd_CheckpointDocumentsFormatAdaptationRules(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionAfterStep6b := doc[step6bPos:]

	// Format adaptation: criteria only, steps only, both, neither
	hasAdaptationRules := strings.Contains(sectionAfterStep6b, "only criteria") ||
		strings.Contains(sectionAfterStep6b, "criteria only") ||
		strings.Contains(sectionAfterStep6b, "only steps") ||
		strings.Contains(sectionAfterStep6b, "steps only") ||
		strings.Contains(sectionAfterStep6b, "If only criteria") ||
		strings.Contains(sectionAfterStep6b, "If only steps") ||
		(strings.Contains(sectionAfterStep6b, "omit") && strings.Contains(sectionAfterStep6b, "criteria"))

	if !hasAdaptationRules {
		t.Error("slice.md checkpoint missing format adaptation rules (criteria only, steps only, both, neither cases)")
	}
}

func TestSliceMd_CheckpointDocumentsNeitherCaseSimplifiedPrompt(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionAfterStep6b := doc[step6bPos:]

	// When neither criteria nor steps exist, a simplified prompt is shown
	hasNeitherCase := strings.Contains(sectionAfterStep6b, "neither") ||
		strings.Contains(sectionAfterStep6b, "simplified") ||
		strings.Contains(sectionAfterStep6b, "Does the output match")

	if !hasNeitherCase {
		t.Error("slice.md checkpoint missing 'neither criteria nor steps' case with simplified prompt")
	}
}

func TestSliceMd_CheckpointDocumentsYesOrOneAsApprovedResponse(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionAfterStep6b := doc[step6bPos:]

	// Contract C-66: "1" or "Yes" (case-insensitive) means approved
	hasApprovedResponseRule := strings.Contains(sectionAfterStep6b, "\"Yes\"") ||
		strings.Contains(sectionAfterStep6b, "\"1\"") ||
		strings.Contains(sectionAfterStep6b, "case-insensitive") ||
		(strings.Contains(sectionAfterStep6b, "1") && strings.Contains(sectionAfterStep6b, "approved"))

	if !hasApprovedResponseRule {
		t.Error("slice.md checkpoint missing documentation that '1' or 'Yes' (case-insensitive) means approved")
	}
}
