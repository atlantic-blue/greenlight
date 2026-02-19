package installer_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// S-24 Tests: Rejection Flow
// These tests verify that src/commands/gl/slice.md and src/references/verification-tiers.md
// contain the full rejection flow documentation as defined in contracts C-67, C-68, and C-69.

// =============================================================================
// C-67 Tests: RejectionClassification
// Gap classification UX — presented to the user after rejection
// Verifies content in slice.md (Step 6b) and verification-tiers.md (Rejection Flow)
// =============================================================================

func TestSliceMd_RejectionFlowPresentsGapClassificationWithThreeOptions(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	hasOption1 := strings.Contains(sectionContent, "1)") || strings.Contains(sectionContent, "1.")
	hasOption2 := strings.Contains(sectionContent, "2)") || strings.Contains(sectionContent, "2.")
	hasOption3 := strings.Contains(sectionContent, "3)") || strings.Contains(sectionContent, "3.")

	if !hasOption1 || !hasOption2 || !hasOption3 {
		t.Error("slice.md Step 6b rejection flow missing three numbered options (1, 2, 3) for gap classification")
	}
}

func TestSliceMd_RejectionFlowOptionsAreInCorrectOrder(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	tightenPos := strings.Index(sectionContent, "ighten")
	revisePos := strings.Index(sectionContent, "evise")
	detailPos := strings.Index(sectionContent, "etail")

	if tightenPos == -1 {
		t.Fatal("slice.md Step 6b rejection flow missing 'tighten' option")
	}
	if revisePos == -1 {
		t.Fatal("slice.md Step 6b rejection flow missing 'revise' option")
	}
	if detailPos == -1 {
		t.Fatal("slice.md Step 6b rejection flow missing 'detail' option")
	}

	if tightenPos >= revisePos {
		t.Errorf("slice.md rejection flow: 'tighten tests' (pos %d) must appear before 'revise contract' (pos %d)", tightenPos, revisePos)
	}
	if revisePos >= detailPos {
		t.Errorf("slice.md rejection flow: 'revise contract' (pos %d) must appear before 'provide more detail' (pos %d)", revisePos, detailPos)
	}
}

func TestSliceMd_RejectionFlowOption1IncludesRoutingConsequence(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// Option 1 must mention the routing consequence: test writer with behavioral feedback
	hasTestWriterRoute := strings.Contains(sectionContent, "test writer") ||
		strings.Contains(sectionContent, "gl-test-writer") ||
		strings.Contains(sectionContent, "return to test writer")

	if !hasTestWriterRoute {
		t.Error("slice.md Step 6b rejection option 1 missing routing consequence — must document that it routes to test writer with behavioral feedback")
	}
}

func TestSliceMd_RejectionFlowOption2IncludesRoutingConsequence(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// Option 2 must mention the routing consequence: contract update before retrying
	hasContractRoute := strings.Contains(sectionContent, "contract") &&
		(strings.Contains(sectionContent, "update") ||
			strings.Contains(sectionContent, "revise") ||
			strings.Contains(sectionContent, "restart"))

	if !hasContractRoute {
		t.Error("slice.md Step 6b rejection option 2 missing routing consequence — must document that it routes to contract revision before retrying")
	}
}

func TestSliceMd_RejectionFlowOption3IncludesRoutingConsequence(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// Option 3 must mention the routing consequence: implementer with additional context
	hasImplementerRoute := strings.Contains(sectionContent, "implementer") ||
		strings.Contains(sectionContent, "fresh implementer") ||
		strings.Contains(sectionContent, "additional context") ||
		strings.Contains(sectionContent, "guidance")

	if !hasImplementerRoute {
		t.Error("slice.md Step 6b rejection option 3 missing routing consequence — must document that it routes to a fresh implementer with additional context")
	}
}

func TestSliceMd_RejectionFlowShowsVerbatimUserFeedback(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// The classification prompt must show the user's verbatim feedback
	hasVerbatimFeedback := strings.Contains(sectionContent, "verbatim") ||
		strings.Contains(sectionContent, "Your feedback") ||
		strings.Contains(sectionContent, "feedback:") ||
		strings.Contains(sectionContent, "user's response") ||
		strings.Contains(sectionContent, "your response")

	if !hasVerbatimFeedback {
		t.Error("slice.md Step 6b rejection flow missing display of user's verbatim feedback in the classification prompt")
	}
}

func TestSliceMd_RejectionFlowOption3CollectsAdditionalDetail(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// Option 3 must collect additional detail before routing
	hasDetailCollection := strings.Contains(sectionContent, "detailed") ||
		strings.Contains(sectionContent, "describe exactly") ||
		strings.Contains(sectionContent, "exactly what") ||
		strings.Contains(sectionContent, "additional detail") ||
		strings.Contains(sectionContent, "more detail")

	if !hasDetailCollection {
		t.Error("slice.md Step 6b rejection option 3 missing additional detail collection — must prompt user to describe exactly what they expected before routing")
	}
}

func TestSliceMd_RejectionFlowInvalidChoiceReprompts(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// InvalidChoice error: re-prompt when user enters something other than 1/2/3
	hasReprompt := strings.Contains(sectionContent, "Please choose 1, 2, or 3") ||
		strings.Contains(sectionContent, "re-prompt") ||
		strings.Contains(sectionContent, "reprompt") ||
		strings.Contains(sectionContent, "invalid choice") ||
		strings.Contains(sectionContent, "InvalidChoice")

	if !hasReprompt {
		t.Error("slice.md Step 6b rejection flow missing InvalidChoice error handling — must re-prompt when user enters something other than 1, 2, or 3")
	}
}

func TestSliceMd_RejectionFlowDefaultsToTestGapAfterFailedReprompts(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// After 2 re-prompts, default to test_gap with free-text as feedback
	hasDefaultAfterReprompts := strings.Contains(sectionContent, "test_gap") ||
		strings.Contains(sectionContent, "default to") ||
		strings.Contains(sectionContent, "defaults to") ||
		(strings.Contains(sectionContent, "re-prompt") && strings.Contains(sectionContent, "tighten"))

	if !hasDefaultAfterReprompts {
		t.Error("slice.md Step 6b rejection flow missing default behaviour after failed re-prompts — must default to test_gap (tighten tests) after 2 re-prompts")
	}
}

func TestSliceMd_RejectionFlowEmptyFeedbackHandled(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// EmptyFeedback: prompt for feedback when user rejection was an empty string
	hasEmptyFeedbackHandling := strings.Contains(sectionContent, "EmptyFeedback") ||
		strings.Contains(sectionContent, "empty") ||
		strings.Contains(sectionContent, "describe what doesn") ||
		strings.Contains(sectionContent, "Please describe") ||
		strings.Contains(sectionContent, "no feedback")

	if !hasEmptyFeedbackHandling {
		t.Error("slice.md Step 6b rejection flow missing EmptyFeedback handling — must prompt for feedback when user rejection was an empty string")
	}
}

func TestSliceMd_RejectionFlowOptionsPhrasedAsActions(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// Options must be phrased as actions (verbs), not categories (nouns)
	// e.g. "Tighten tests" not "Test gap", "Revise contract" not "Contract gap"
	hasTightenVerb := strings.Contains(sectionContent, "Tighten") ||
		strings.Contains(sectionContent, "tighten")
	hasReviseVerb := strings.Contains(sectionContent, "Revise") ||
		strings.Contains(sectionContent, "revise")
	hasProvideVerb := strings.Contains(sectionContent, "Provide") ||
		strings.Contains(sectionContent, "provide") ||
		strings.Contains(sectionContent, "describe")

	if !hasTightenVerb {
		t.Error("slice.md Step 6b rejection option 1 must be phrased as an action (e.g. 'Tighten tests'), not a category")
	}
	if !hasReviseVerb {
		t.Error("slice.md Step 6b rejection option 2 must be phrased as an action (e.g. 'Revise contract'), not a category")
	}
	if !hasProvideVerb {
		t.Error("slice.md Step 6b rejection option 3 must be phrased as an action (e.g. 'Provide more detail'), not a category")
	}
}

func TestVerificationTiersMd_GapClassificationPresentsThreeOptionsInOrder(t *testing.T) {
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

	tightenPos := strings.Index(sectionContent, "ighten")
	revisePos := strings.Index(sectionContent, "evise")
	detailPos := strings.Index(sectionContent, "etail")

	if tightenPos == -1 {
		t.Fatal("verification-tiers.md Rejection Flow missing 'tighten' option")
	}
	if revisePos == -1 {
		t.Fatal("verification-tiers.md Rejection Flow missing 'revise' option")
	}
	if detailPos == -1 {
		t.Fatal("verification-tiers.md Rejection Flow missing 'detail' option")
	}

	if tightenPos >= revisePos {
		t.Errorf("verification-tiers.md Rejection Flow: 'tighten' (pos %d) must appear before 'revise' (pos %d)", tightenPos, revisePos)
	}
	if revisePos >= detailPos {
		t.Errorf("verification-tiers.md Rejection Flow: 'revise' (pos %d) must appear before 'detail' (pos %d)", revisePos, detailPos)
	}
}

func TestVerificationTiersMd_RejectionFlowOption3CollectsDetailBeforeRouting(t *testing.T) {
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

	// Option 3 must collect additional detail before routing to the implementer or test writer
	hasDetailCollection := strings.Contains(sectionContent, "detail") &&
		(strings.Contains(sectionContent, "before") ||
			strings.Contains(sectionContent, "collect") ||
			strings.Contains(sectionContent, "describe") ||
			strings.Contains(sectionContent, "additional"))

	if !hasDetailCollection {
		t.Error("verification-tiers.md Rejection Flow option 3 must document that additional detail is collected before routing")
	}
}

func TestVerificationTiersMd_RejectionFlowOptionsIncludeRoutingConsequences(t *testing.T) {
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

	// Each option must show where it routes (transparency)
	hasTestWriterRoute := strings.Contains(sectionContent, "test writer") ||
		strings.Contains(sectionContent, "gl-test-writer")
	hasContractRoute := strings.Contains(sectionContent, "contract") &&
		(strings.Contains(sectionContent, "update") || strings.Contains(sectionContent, "restart") || strings.Contains(sectionContent, "revise"))
	hasImplementerRoute := strings.Contains(sectionContent, "implementer") ||
		strings.Contains(sectionContent, "guidance")

	if !hasTestWriterRoute {
		t.Error("verification-tiers.md Rejection Flow option 1 missing routing consequence to test writer")
	}
	if !hasContractRoute {
		t.Error("verification-tiers.md Rejection Flow option 2 missing routing consequence to contract revision")
	}
	if !hasImplementerRoute {
		t.Error("verification-tiers.md Rejection Flow option 3 missing routing consequence to implementer with guidance")
	}
}

func TestVerificationTiersMd_RejectionFlowMapsChoicesToInternalClassification(t *testing.T) {
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

	// Internal classification mappings must be documented
	hasTestGap := strings.Contains(sectionContent, "test_gap") ||
		strings.Contains(sectionContent, "test gap")
	hasContractGap := strings.Contains(sectionContent, "contract_gap") ||
		strings.Contains(sectionContent, "contract gap")
	hasImplementationGap := strings.Contains(sectionContent, "implementation_gap") ||
		strings.Contains(sectionContent, "implementation gap")

	if !hasTestGap {
		t.Error("verification-tiers.md Rejection Flow missing 'test_gap' internal classification mapping for option 1")
	}
	if !hasContractGap {
		t.Error("verification-tiers.md Rejection Flow missing 'contract_gap' internal classification mapping for option 2")
	}
	if !hasImplementationGap {
		t.Error("verification-tiers.md Rejection Flow missing 'implementation_gap' internal classification mapping for option 3")
	}
}

// =============================================================================
// C-68 Tests: RejectionToTestWriter
// Test writer spawn context — after test_gap or implementation_gap classification
// Verifies content in slice.md (Step 6b Rejection Flow) and verification-tiers.md
// =============================================================================

func TestSliceMd_RejectionFlowTestWriterSpawnIncludesRejectionContextXML(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// Test writer spawn must include rejection_context XML block
	hasRejectionContextXML := strings.Contains(sectionContent, "rejection_context") ||
		strings.Contains(sectionContent, "<rejection") ||
		strings.Contains(sectionContent, "rejection context")

	if !hasRejectionContextXML {
		t.Error("slice.md Step 6b rejection flow missing rejection_context XML block for test writer spawn")
	}
}

func TestSliceMd_RejectionFlowTestWriterSpawnContextIncludesFeedbackField(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// The rejection_context must include a feedback field with verbatim rejection
	hasFeedbackField := strings.Contains(sectionContent, "<feedback>") ||
		strings.Contains(sectionContent, "feedback:") ||
		(strings.Contains(sectionContent, "verbatim") && strings.Contains(sectionContent, "rejection"))

	if !hasFeedbackField {
		t.Error("slice.md Step 6b rejection context for test writer spawn missing feedback field (verbatim user rejection)")
	}
}

func TestSliceMd_RejectionFlowTestWriterSpawnContextIncludesClassificationField(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// The rejection_context must include a classification field
	hasClassificationField := strings.Contains(sectionContent, "<classification>") ||
		strings.Contains(sectionContent, "classification:")

	if !hasClassificationField {
		t.Error("slice.md Step 6b rejection context for test writer spawn missing classification field (test_gap or implementation_gap)")
	}
}

func TestSliceMd_RejectionFlowTestWriterSpawnContextIncludesDetailedFeedbackField(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// The rejection_context must include a detailed_feedback field (for option 3)
	hasDetailedFeedbackField := strings.Contains(sectionContent, "<detailed_feedback>") ||
		strings.Contains(sectionContent, "detailed_feedback") ||
		strings.Contains(sectionContent, "detailed feedback")

	if !hasDetailedFeedbackField {
		t.Error("slice.md Step 6b rejection context for test writer spawn missing detailed_feedback field (additional detail from user option 3)")
	}
}

func TestSliceMd_RejectionFlowTestWriterSpawnContextIncludesContractBlock(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// Test writer spawn must include the full contract definition(s)
	hasContractBlock := strings.Contains(sectionContent, "<contract>") ||
		strings.Contains(sectionContent, "contract definition") ||
		strings.Contains(sectionContent, "full contract")

	if !hasContractBlock {
		t.Error("slice.md Step 6b rejection flow test writer spawn missing contract block — test writer must receive full contract definitions")
	}
}

func TestSliceMd_RejectionFlowTestWriterSpawnContextIncludesAcceptanceCriteriaBlock(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// Test writer spawn must include acceptance_criteria block
	hasAcceptanceCriteriaBlock := strings.Contains(sectionContent, "<acceptance_criteria>") ||
		strings.Contains(sectionContent, "acceptance_criteria") ||
		strings.Contains(sectionContent, "acceptance criteria")

	if !hasAcceptanceCriteriaBlock {
		t.Error("slice.md Step 6b rejection flow test writer spawn missing acceptance_criteria block")
	}
}

func TestSliceMd_RejectionFlowTestWriterReceivesBehavioralFeedbackOnly(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// Agent isolation: test writer must NOT receive implementation code
	hasBehavioralOnlyConstraint := strings.Contains(sectionContent, "behavioral feedback") ||
		strings.Contains(sectionContent, "behavioural feedback") ||
		strings.Contains(sectionContent, "behavioral feedback only") ||
		(strings.Contains(sectionContent, "test writer") && strings.Contains(sectionContent, "not") && strings.Contains(sectionContent, "implementation"))

	if !hasBehavioralOnlyConstraint {
		t.Error("slice.md Step 6b rejection flow missing agent isolation constraint — test writer must receive behavioral feedback only, never implementation code")
	}
}

func TestSliceMd_RejectionFlowNewTestsAreAdditive(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// New tests must be additive — existing passing tests must not be removed or modified
	hasAdditivConstraint := strings.Contains(sectionContent, "additive") ||
		strings.Contains(sectionContent, "existing tests") ||
		strings.Contains(sectionContent, "do not remove") ||
		strings.Contains(sectionContent, "passing tests")

	if !hasAdditivConstraint {
		t.Error("slice.md Step 6b rejection flow missing additive constraint — new tests must not remove or modify existing passing tests")
	}
}

func TestSliceMd_RejectionFlowImplementerSpawnedAfterTestWriter(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	testWriterPos := strings.Index(sectionContent, "test writer")
	if testWriterPos == -1 {
		testWriterPos = strings.Index(sectionContent, "gl-test-writer")
	}
	implementerPos := strings.LastIndex(sectionContent, "implementer")

	if testWriterPos == -1 {
		t.Fatal("slice.md Step 6b rejection flow missing test writer reference")
	}
	if implementerPos == -1 {
		t.Fatal("slice.md Step 6b rejection flow missing implementer reference after test writer")
	}

	if testWriterPos >= implementerPos {
		t.Errorf("slice.md Step 6b rejection flow: implementer (pos %d) must be spawned after test writer (pos %d)", implementerPos, testWriterPos)
	}
}

func TestSliceMd_RejectionFlowImplementerReceivesTestNamesOnly(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// Implementer must receive test names only — not test source code
	hasTestNamesOnlyConstraint := strings.Contains(sectionContent, "test names") ||
		strings.Contains(sectionContent, "names only") ||
		strings.Contains(sectionContent, "not test source") ||
		strings.Contains(sectionContent, "not test code")

	if !hasTestNamesOnlyConstraint {
		t.Error("slice.md Step 6b rejection flow missing implementer isolation constraint — implementer must receive test names only, not test source code")
	}
}

func TestSliceMd_RejectionFlowFullVerificationCycleRerunsAfterImplementation(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// Full verification cycle must re-run: Step 4 (tests), Step 6 (verifier), Step 6b (checkpoint)
	hasFullVerificationRerun := strings.Contains(sectionContent, "re-run") ||
		strings.Contains(sectionContent, "rerun") ||
		strings.Contains(sectionContent, "re-present") ||
		strings.Contains(sectionContent, "full verification") ||
		strings.Contains(sectionContent, "re-run Step 6b") ||
		strings.Contains(sectionContent, "run Step 6b")

	if !hasFullVerificationRerun {
		t.Error("slice.md Step 6b rejection flow missing full verification cycle re-run — Step 4, Step 6, and Step 6b must all re-run after implementation")
	}
}

func TestSliceMd_RejectionFlowHandlesTestWriterSpawnFailure(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// TestWriterSpawnFailure error must be handled
	hasSpawnFailureHandling := strings.Contains(sectionContent, "TestWriterSpawnFailure") ||
		strings.Contains(sectionContent, "spawn failure") ||
		strings.Contains(sectionContent, "fails to spawn") ||
		strings.Contains(sectionContent, "agent fails") ||
		(strings.Contains(sectionContent, "retry") && strings.Contains(sectionContent, "pause"))

	if !hasSpawnFailureHandling {
		t.Error("slice.md Step 6b rejection flow missing TestWriterSpawnFailure error handling — must offer retry, pause, or skip when test writer spawn fails")
	}
}

func TestSliceMd_RejectionFlowHandlesExistingTestsRegression(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// ExistingTestsRegressed must be handled: implementer must fix regressions before proceeding
	hasRegressionHandling := strings.Contains(sectionContent, "ExistingTestsRegressed") ||
		strings.Contains(sectionContent, "regression") ||
		strings.Contains(sectionContent, "previously passing") ||
		strings.Contains(sectionContent, "regressed")

	if !hasRegressionHandling {
		t.Error("slice.md Step 6b rejection flow missing ExistingTestsRegressed handling — implementer must fix regressions in previously passing tests before proceeding")
	}
}

func TestSliceMd_RejectionFlowIntegratesWithCircuitBreaker(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// Rejection flow must integrate with the existing circuit breaker protocol
	hasCircuitBreakerIntegration := strings.Contains(sectionContent, "circuit breaker") ||
		strings.Contains(sectionContent, "circuit-breaker") ||
		strings.Contains(sectionContent, "CIRCUIT BREAK") ||
		strings.Contains(sectionContent, "C-55")

	if !hasCircuitBreakerIntegration {
		t.Error("slice.md Step 6b rejection flow missing circuit breaker integration — new tests failing after implementation must enter normal circuit breaker flow")
	}
}

func TestVerificationTiersMd_AgentIsolationDocumentsTestWriterReceivesBehavioralFeedback(t *testing.T) {
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

	// Test writer must receive behavioral feedback only — never implementation source or test source
	hasBehavioralOnlyConstraint := strings.Contains(sectionContent, "behavioral feedback only") ||
		strings.Contains(sectionContent, "behavioural feedback only") ||
		strings.Contains(sectionContent, "behavioral feedback") ||
		strings.Contains(sectionContent, "behavioural feedback")

	if !hasBehavioralOnlyConstraint {
		t.Error("verification-tiers.md Agent Isolation section missing constraint that test writer receives behavioral feedback only")
	}
}

func TestVerificationTiersMd_AgentIsolationDocumentsTestWriterNotReceivingImplementationCode(t *testing.T) {
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

	// Test writer must not receive implementation code
	hasNoImplementationCodeConstraint := strings.Contains(sectionContent, "does NOT receive") ||
		strings.Contains(sectionContent, "does not receive") ||
		strings.Contains(sectionContent, "never implementation") ||
		strings.Contains(sectionContent, "no implementation") ||
		(strings.Contains(sectionContent, "NOT") && strings.Contains(sectionContent, "implementation"))

	if !hasNoImplementationCodeConstraint {
		t.Error("verification-tiers.md Agent Isolation section missing constraint that test writer does NOT receive implementation code")
	}
}

func TestVerificationTiersMd_AgentIsolationDocumentsImplementerReceivesTestNamesOnly(t *testing.T) {
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

	// Implementer must receive test names only — not test source code
	hasTestNamesOnlyConstraint := strings.Contains(sectionContent, "test names") ||
		strings.Contains(sectionContent, "names only") ||
		strings.Contains(sectionContent, "not test source") ||
		(strings.Contains(sectionContent, "implementer") && strings.Contains(sectionContent, "names"))

	if !hasTestNamesOnlyConstraint {
		t.Error("verification-tiers.md Agent Isolation section missing constraint that implementer receives test names only, not test source code")
	}
}

func TestVerificationTiersMd_AgentIsolationDocumentsOrchestratorFiltersFeedback(t *testing.T) {
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

	// Orchestrator is responsible for filtering feedback before passing to test writer
	hasOrchestratorFilterConstraint := strings.Contains(sectionContent, "orchestrator") &&
		(strings.Contains(sectionContent, "filter") ||
			strings.Contains(sectionContent, "responsible") ||
			strings.Contains(sectionContent, "filtering"))

	if !hasOrchestratorFilterConstraint {
		t.Error("verification-tiers.md Agent Isolation section missing statement that orchestrator is responsible for filtering feedback before passing to gl-test-writer")
	}
}

// =============================================================================
// C-69 Tests: RejectionToContractRevision
// Contract revision route — after contract_gap classification
// Verifies content in slice.md (Step 6b) and verification-tiers.md (Rejection Flow)
// =============================================================================

func TestSliceMd_RejectionFlowContractRevisionShowsContractRevisionHeader(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// CONTRACT REVISION header must be shown with slice_id and slice_name
	hasContractRevisionHeader := strings.Contains(sectionContent, "CONTRACT REVISION") ||
		strings.Contains(sectionContent, "Contract Revision") ||
		strings.Contains(sectionContent, "contract revision")

	if !hasContractRevisionHeader {
		t.Error("slice.md Step 6b rejection flow option 2 missing CONTRACT REVISION header format")
	}
}

func TestSliceMd_RejectionFlowContractRevisionDisplaysFullContractText(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// User must see full contract text before making revision
	hasFullContractDisplay := strings.Contains(sectionContent, "full contract") ||
		strings.Contains(sectionContent, "current contract") ||
		strings.Contains(sectionContent, "contract text") ||
		strings.Contains(sectionContent, "contract definition")

	if !hasFullContractDisplay {
		t.Error("slice.md Step 6b rejection flow contract revision missing display of full contract text — user must see the full contract before making changes")
	}
}

func TestSliceMd_RejectionFlowContractRevisionAllowsAcceptanceCriteriaEdits(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// User must be able to edit acceptance criteria
	hasAcceptanceCriteriaEdit := strings.Contains(sectionContent, "acceptance criteria") &&
		(strings.Contains(sectionContent, "edit") ||
			strings.Contains(sectionContent, "update") ||
			strings.Contains(sectionContent, "change") ||
			strings.Contains(sectionContent, "revise"))

	if !hasAcceptanceCriteriaEdit {
		t.Error("slice.md Step 6b rejection flow contract revision missing acceptance criteria edit capability")
	}
}

func TestSliceMd_RejectionFlowFundamentalRevisionRecommendsAddSlice(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// Fundamental contract changes must recommend /gl:add-slice
	hasFundamentalRevisionPath := strings.Contains(sectionContent, "gl:add-slice") ||
		strings.Contains(sectionContent, "/gl:add-slice") ||
		strings.Contains(sectionContent, "add-slice") ||
		strings.Contains(sectionContent, "architect") ||
		strings.Contains(sectionContent, "fundamental")

	if !hasFundamentalRevisionPath {
		t.Error("slice.md Step 6b rejection flow contract revision missing recommendation to re-run /gl:add-slice for fundamental revisions (input/output changes, new boundaries)")
	}
}

func TestSliceMd_RejectionFlowContractRevisionRestartsSliceFromStep1(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// After contract revision, slice must restart from Step 1 (test writing)
	hasSliceRestart := strings.Contains(sectionContent, "restart") ||
		strings.Contains(sectionContent, "re-run from Step 1") ||
		strings.Contains(sectionContent, "restart from Step 1") ||
		strings.Contains(sectionContent, "Step 1") ||
		strings.Contains(sectionContent, "from the beginning")

	if !hasSliceRestart {
		t.Error("slice.md Step 6b rejection flow contract revision missing slice restart from Step 1 — after revision the full TDD loop must re-run")
	}
}

func TestSliceMd_RejectionFlowContractRevisionOffersRollbackBeforeRestart(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// Rollback to checkpoint must be offered before restart
	hasRollbackOffer := strings.Contains(sectionContent, "rollback") ||
		strings.Contains(sectionContent, "roll back") ||
		strings.Contains(sectionContent, "checkpoint") ||
		strings.Contains(sectionContent, "checkpoint tag")

	if !hasRollbackOffer {
		t.Error("slice.md Step 6b rejection flow contract revision missing rollback offer — rollback to checkpoint must be offered before restarting")
	}
}

func TestSliceMd_RejectionFlowEmptyRevisionReprompts(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// EmptyRevision error: re-prompt when user provides no revision description
	hasEmptyRevisionHandling := strings.Contains(sectionContent, "EmptyRevision") ||
		strings.Contains(sectionContent, "no revision") ||
		strings.Contains(sectionContent, "empty revision") ||
		strings.Contains(sectionContent, "Please describe what the contract") ||
		strings.Contains(sectionContent, "what the contract should")

	if !hasEmptyRevisionHandling {
		t.Error("slice.md Step 6b rejection flow contract revision missing EmptyRevision error handling — must re-prompt when user provides no revision description")
	}
}

func TestSliceMd_RejectionFlowContractRevisionIncrementsRejectionCounter(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// Contract revision must increment the rejection counter
	hasCounterIncrement := strings.Contains(sectionContent, "rejection counter") ||
		strings.Contains(sectionContent, "increment") ||
		strings.Contains(sectionContent, "counter")

	if !hasCounterIncrement {
		t.Error("slice.md Step 6b rejection flow missing rejection counter increment documentation — all rejection paths (including contract revision) must increment the counter")
	}
}

func TestVerificationTiersMd_RejectionFlowDocumentsContractRevisionWithSliceRestart(t *testing.T) {
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

	// Option 2 must document that contract revision leads to a slice restart
	hasContractRevisionRestart := strings.Contains(sectionContent, "restart") ||
		strings.Contains(sectionContent, "re-run") ||
		strings.Contains(sectionContent, "retrying") ||
		strings.Contains(sectionContent, "before retrying")

	if !hasContractRevisionRestart {
		t.Error("verification-tiers.md Rejection Flow option 2 missing documentation that contract revision leads to slice restart from Step 1")
	}
}

func TestVerificationTiersMd_RejectionFlowDocumentsRejectionCounterIncrementForContractRevision(t *testing.T) {
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

	// Rejection counter must increment on every rejection path, including contract revision
	hasAllPathsIncrement := strings.Contains(sectionContent, "each rejection") ||
		strings.Contains(sectionContent, "every rejection") ||
		strings.Contains(sectionContent, "Incremented by 1") ||
		strings.Contains(sectionContent, "incremented") ||
		strings.Contains(sectionContent, "increment")

	if !hasAllPathsIncrement {
		t.Error("verification-tiers.md Rejection Counter section missing documentation that the counter increments on all rejection paths including contract revision")
	}
}

// =============================================================================
// Cross-cutting invariant tests
// These verify invariants that span all three contracts (C-67, C-68, C-69)
// =============================================================================

func TestSliceMd_RejectionFlowAppearsWithinStep6bSection(t *testing.T) {
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

	// Rejection flow content must be within Step 6b, not after Step 7
	step6bSection := doc[step6bPos:step7Pos]

	hasRejectionFlowInSection := strings.Contains(step6bSection, "rejection") ||
		strings.Contains(step6bSection, "Rejection")

	if !hasRejectionFlowInSection {
		t.Error("slice.md Step 6b rejection flow content must appear within the Step 6b section (before Step 7)")
	}
}

func TestVerificationTiersMd_RejectionFlowSectionAppearsBeforeRejectionCounter(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-tiers.md"))
	if err != nil {
		t.Fatalf("failed to read verification-tiers.md: %v", err)
	}

	doc := string(content)

	rejectionFlowPos := strings.Index(doc, "Rejection Flow")
	if rejectionFlowPos == -1 {
		t.Fatal("verification-tiers.md missing 'Rejection Flow' section")
	}

	rejectionCounterPos := strings.Index(doc, "Rejection Counter")
	if rejectionCounterPos == -1 {
		t.Fatal("verification-tiers.md missing 'Rejection Counter' section")
	}

	if rejectionFlowPos >= rejectionCounterPos {
		t.Errorf("verification-tiers.md 'Rejection Flow' section (pos %d) must appear before 'Rejection Counter' section (pos %d)", rejectionFlowPos, rejectionCounterPos)
	}
}

func TestVerificationTiersMd_RejectionFlowSectionAppearsBeforeAgentIsolation(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-tiers.md"))
	if err != nil {
		t.Fatalf("failed to read verification-tiers.md: %v", err)
	}

	doc := string(content)

	rejectionFlowPos := strings.Index(doc, "Rejection Flow")
	if rejectionFlowPos == -1 {
		t.Fatal("verification-tiers.md missing 'Rejection Flow' section")
	}

	agentIsolationPos := strings.Index(doc, "Agent Isolation")
	if agentIsolationPos == -1 {
		t.Fatal("verification-tiers.md missing 'Agent Isolation' section")
	}

	if rejectionFlowPos >= agentIsolationPos {
		t.Errorf("verification-tiers.md 'Rejection Flow' section (pos %d) must appear before 'Agent Isolation' section (pos %d)", rejectionFlowPos, agentIsolationPos)
	}
}

func TestVerificationTiersMd_SecurityConstraintFeedbackNotExecutedAsCode(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-tiers.md"))
	if err != nil {
		t.Fatalf("failed to read verification-tiers.md: %v", err)
	}

	doc := string(content)

	// Security invariant: user feedback must not be used to execute commands or modify files directly
	hasFeedbackSecurityConstraint := strings.Contains(doc, "not.*execute") ||
		strings.Contains(doc, "text context") ||
		strings.Contains(doc, "no code execution") ||
		strings.Contains(doc, "treated as") ||
		strings.Contains(doc, "not used to execute") ||
		strings.Contains(doc, "context, not executable") ||
		strings.Contains(doc, "behavioral context")

	if !hasFeedbackSecurityConstraint {
		t.Error("verification-tiers.md missing security constraint that user feedback is treated as text context only, never executed as code or used to modify files directly")
	}
}

func TestSliceMd_RejectionFlowUserFeedbackNotExecutedAsCode(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	step6bPos := strings.Index(doc, "Step 6b")
	if step6bPos == -1 {
		t.Fatal("slice.md missing 'Step 6b' section")
	}

	sectionContent := doc[step6bPos:]

	// Security: user feedback must be treated as text context only
	hasFeedbackSecurityConstraint := strings.Contains(sectionContent, "behavioral feedback") ||
		strings.Contains(sectionContent, "text context") ||
		strings.Contains(sectionContent, "not executable") ||
		strings.Contains(sectionContent, "no code execution") ||
		strings.Contains(sectionContent, "behavioral context") ||
		(strings.Contains(sectionContent, "feedback") && strings.Contains(sectionContent, "context"))

	if !hasFeedbackSecurityConstraint {
		t.Error("slice.md Step 6b rejection flow missing security constraint — user feedback must be treated as behavioral text context, never executed as code or used to modify files directly")
	}
}
