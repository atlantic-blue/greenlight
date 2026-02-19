package installer_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// C-75 Tests: ArchitectTierGuidance
// These tests verify that src/agents/gl-architect.md contains the Verification
// Tier Selection guidance section and updated output checklist as defined in
// contract C-75 (slice S-27: Architect Integration).

// --- Tier Selection Guidance section ---

func TestArchitectMd_ContainsVerificationTierSelectionSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "Verification Tier Selection") {
		t.Error("gl-architect.md missing '## Verification Tier Selection' section — contract C-75 requires this guidance section")
	}
}

func TestArchitectMd_TierSelectionDefaultIsVerify(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// The guidance must state that verify is the default tier.
	// Contract says: "Default: verify."
	hasVerifyDefault := strings.Contains(doc, "Default: verify") ||
		strings.Contains(doc, "default: verify") ||
		strings.Contains(doc, "**Default: verify**")

	if !hasVerifyDefault {
		t.Error("gl-architect.md missing 'Default: verify' statement in Verification Tier Selection section — contract C-75 requires verify to be the stated default")
	}
}

func TestArchitectMd_TierSelectionWhenInDoubtUseVerify(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract states: "When in doubt, use verify."
	hasWhenInDoubt := strings.Contains(doc, "When in doubt, use verify") ||
		strings.Contains(doc, "when in doubt, use verify") ||
		strings.Contains(doc, "in doubt") && strings.Contains(doc, "use verify")

	if !hasWhenInDoubt {
		t.Error("gl-architect.md missing 'When in doubt, use verify' guidance — contract C-75 requires this safe-default statement")
	}
}

func TestArchitectMd_TierSelectionContainsWhenToUseAutoSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	hasWhenToUseAuto := strings.Contains(doc, "When to use auto") ||
		strings.Contains(doc, "**When to use auto:**")

	if !hasWhenToUseAuto {
		t.Error("gl-architect.md missing 'When to use auto' subsection in Verification Tier Selection — contract C-75 requires this subsection")
	}
}

func TestArchitectMd_TierSelectionAutoExamplesIncludeInfrastructure(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract lists infrastructure contracts as an auto example.
	hasInfrastructureExample := strings.Contains(doc, "Infrastructure") ||
		strings.Contains(doc, "infrastructure")

	if !hasInfrastructureExample {
		t.Error("gl-architect.md missing infrastructure example in 'When to use auto' guidance — contract C-75 lists infrastructure contracts as an auto-tier example")
	}
}

func TestArchitectMd_TierSelectionAutoExamplesIncludeInternalPlumbing(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract lists "Internal plumbing" as an auto example.
	hasPlumbingExample := strings.Contains(doc, "plumbing") ||
		strings.Contains(doc, "internal plumbing")

	if !hasPlumbingExample {
		t.Error("gl-architect.md missing 'plumbing' example in 'When to use auto' guidance — contract C-75 lists internal plumbing as an auto-tier example")
	}
}

func TestArchitectMd_TierSelectionAutoExamplesIncludeSchemaOrTypeDefinitions(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract lists schema/type definitions with no user-visible behaviour as an auto example.
	hasSchemaExample := strings.Contains(doc, "Schema") ||
		strings.Contains(doc, "schema") ||
		strings.Contains(doc, "type definition") ||
		strings.Contains(doc, "Type definition")

	if !hasSchemaExample {
		t.Error("gl-architect.md missing schema/type definitions example in 'When to use auto' guidance — contract C-75 lists schema/type definitions as an auto-tier example")
	}
}

func TestArchitectMd_TierSelectionAutoExamplesIncludeCICD(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract lists "Build tooling, CI/CD configuration" as an auto example.
	hasCICDExample := strings.Contains(doc, "CI/CD") ||
		strings.Contains(doc, "ci/cd") ||
		strings.Contains(doc, "Build tooling") ||
		strings.Contains(doc, "build tooling")

	if !hasCICDExample {
		t.Error("gl-architect.md missing CI/CD / build tooling example in 'When to use auto' guidance — contract C-75 lists CI/CD configuration as an auto-tier example")
	}
}

func TestArchitectMd_TierSelectionContainsWhenToUseVerifySection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	hasWhenToUseVerify := strings.Contains(doc, "When to use verify") ||
		strings.Contains(doc, "**When to use verify:**")

	if !hasWhenToUseVerify {
		t.Error("gl-architect.md missing 'When to use verify' subsection in Verification Tier Selection — contract C-75 requires this subsection")
	}
}

func TestArchitectMd_TierSelectionVerifyExamplesIncludeUserVisibleBehaviour(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract lists "Any contract with user-visible behaviour" as a verify example.
	hasUserVisibleExample := strings.Contains(doc, "user-visible") ||
		strings.Contains(doc, "user visible")

	if !hasUserVisibleExample {
		t.Error("gl-architect.md missing 'user-visible behaviour' example in 'When to use verify' guidance — contract C-75 lists user-visible behaviour as a verify-tier example")
	}
}

func TestArchitectMd_TierSelectionVerifyExamplesIncludeUIComponents(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract lists "UI components, page layouts, visual output" as a verify example.
	hasUIExample := strings.Contains(doc, "UI components") ||
		strings.Contains(doc, "ui components") ||
		strings.Contains(doc, "page layout") ||
		strings.Contains(doc, "visual output")

	if !hasUIExample {
		t.Error("gl-architect.md missing UI components / page layout example in 'When to use verify' guidance — contract C-75 lists UI components as a verify-tier example")
	}
}

func TestArchitectMd_TierSelectionVerifyExamplesIncludeAPIEndpoints(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract lists "API endpoints where response format matters to the user" as a verify example.
	hasAPIExample := strings.Contains(doc, "API endpoint") ||
		strings.Contains(doc, "api endpoint") ||
		strings.Contains(doc, "response format")

	if !hasAPIExample {
		t.Error("gl-architect.md missing API endpoints example in 'When to use verify' guidance — contract C-75 lists API endpoints as a verify-tier example")
	}
}

func TestArchitectMd_TierSelectionVerifyExamplesIncludeBusinessLogic(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract lists "Business logic where intent may differ from specification" as a verify example.
	hasBusinessLogicExample := strings.Contains(doc, "Business logic") ||
		strings.Contains(doc, "business logic")

	if !hasBusinessLogicExample {
		t.Error("gl-architect.md missing business logic example in 'When to use verify' guidance — contract C-75 lists business logic as a verify-tier example")
	}
}

func TestArchitectMd_TierSelectionVerifyIsTheSafeDefault(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract says: "When you are uncertain (verify is the safe default)"
	hasSafeDefaultStatement := strings.Contains(doc, "safe default") ||
		strings.Contains(doc, "verify is the safe default")

	if !hasSafeDefaultStatement {
		t.Error("gl-architect.md missing 'safe default' statement in Verification Tier Selection — contract C-75 requires verify to be called the safe default when uncertain")
	}
}

func TestArchitectMd_TierSelectionContainsWritingAcceptanceCriteriaSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	hasAcceptanceCriteriaSection := strings.Contains(doc, "Writing acceptance criteria") ||
		strings.Contains(doc, "**Writing acceptance criteria:**")

	if !hasAcceptanceCriteriaSection {
		t.Error("gl-architect.md missing 'Writing acceptance criteria' subsection in Verification Tier Selection — contract C-75 requires this guidance subsection")
	}
}

func TestArchitectMd_AcceptanceCriteriaGuidanceRequiresBehavioralStatements(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract: "Each criterion is a behavioral statement the user can observe"
	hasBehavioralStatement := strings.Contains(doc, "behavioral statement") ||
		strings.Contains(doc, "behavioural statement") ||
		strings.Contains(doc, "user can observe")

	if !hasBehavioralStatement {
		t.Error("gl-architect.md missing 'behavioral statement the user can observe' guidance in Writing acceptance criteria — contract C-75 requires criteria to describe observable behaviour")
	}
}

func TestArchitectMd_AcceptanceCriteriaGuidanceRequiresPresentTense(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract: "Use present tense: 'User sees X', 'Page displays Y', 'API returns Z'"
	hasPresentTenseGuidance := strings.Contains(doc, "present tense") ||
		strings.Contains(doc, "Present tense")

	if !hasPresentTenseGuidance {
		t.Error("gl-architect.md missing 'present tense' guidance in Writing acceptance criteria — contract C-75 requires criteria to use present tense")
	}
}

func TestArchitectMd_AcceptanceCriteriaGuidanceIncludesUserSeesExample(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract gives "User sees X" as an example of present-tense criteria.
	hasUserSeesExample := strings.Contains(doc, "User sees") ||
		strings.Contains(doc, "user sees")

	if !hasUserSeesExample {
		t.Error("gl-architect.md missing 'User sees X' example in Writing acceptance criteria — contract C-75 lists this as a present-tense example")
	}
}

func TestArchitectMd_AcceptanceCriteriaGuidanceIncludesPageDisplaysExample(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract gives "Page displays Y" as an example of present-tense criteria.
	hasPageDisplaysExample := strings.Contains(doc, "Page displays") ||
		strings.Contains(doc, "page displays")

	if !hasPageDisplaysExample {
		t.Error("gl-architect.md missing 'Page displays Y' example in Writing acceptance criteria — contract C-75 lists this as a present-tense example")
	}
}

func TestArchitectMd_AcceptanceCriteriaGuidanceRequiresSpecificNotVague(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract: "Be specific: 'Cards render in a 3-column grid' not 'Layout looks correct'"
	hasSpecificGuidance := strings.Contains(doc, "Be specific") ||
		strings.Contains(doc, "be specific") ||
		strings.Contains(doc, "specific:")

	if !hasSpecificGuidance {
		t.Error("gl-architect.md missing 'Be specific' guidance in Writing acceptance criteria — contract C-75 requires specific (not vague) criteria")
	}
}

func TestArchitectMd_AcceptanceCriteriaGuidanceIncludesColumnGridExample(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract contrasts "3-column grid" (specific) vs "looks correct" (vague).
	hasSpecificVsVagueContrast := strings.Contains(doc, "3-column") ||
		strings.Contains(doc, "column grid")

	if !hasSpecificVsVagueContrast {
		t.Error("gl-architect.md missing 3-column grid example in Writing acceptance criteria — contract C-75 uses this to illustrate specific vs vague criteria")
	}
}

func TestArchitectMd_AcceptanceCriteriaGuidanceMentions2To5CriteriaGuideline(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract: "2-5 criteria per contract (more than 5 suggests the contract is too large)"
	has2To5Guideline := strings.Contains(doc, "2-5") ||
		strings.Contains(doc, "2–5") ||
		strings.Contains(doc, "two to five") ||
		strings.Contains(doc, "more than 5")

	if !has2To5Guideline {
		t.Error("gl-architect.md missing 2-5 criteria guideline in Writing acceptance criteria — contract C-75 specifies 2-5 as the recommended range per contract")
	}
}

func TestArchitectMd_AcceptanceCriteriaGuidanceMentionsNegativeCriteria(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract: "Include negative criteria when relevant: 'No error messages appear'"
	hasNegativeCriteriaGuidance := strings.Contains(doc, "negative criteria") ||
		strings.Contains(doc, "No error messages")

	if !hasNegativeCriteriaGuidance {
		t.Error("gl-architect.md missing negative criteria guidance in Writing acceptance criteria — contract C-75 requires mention of negative criteria (e.g. 'No error messages appear')")
	}
}

func TestArchitectMd_TierSelectionContainsWritingStepsSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	hasWritingStepsSection := strings.Contains(doc, "Writing steps") ||
		strings.Contains(doc, "**Writing steps:**")

	if !hasWritingStepsSection {
		t.Error("gl-architect.md missing 'Writing steps' subsection in Verification Tier Selection — contract C-75 requires this guidance subsection")
	}
}

func TestArchitectMd_WritingStepsGuidanceRequiresActionVerbs(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract: "Start each step with an action verb: 'Run...', 'Open...', 'Click...'"
	hasActionVerbGuidance := strings.Contains(doc, "action verb") ||
		strings.Contains(doc, "Start each step")

	if !hasActionVerbGuidance {
		t.Error("gl-architect.md missing action verb guidance in Writing steps — contract C-75 requires steps to start with action verbs (Run, Open, Click)")
	}
}

func TestArchitectMd_WritingStepsGuidanceIncludesRunExample(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract lists "Run..." as an example action verb.
	hasRunExample := strings.Contains(doc, "Run...")  ||
		strings.Contains(doc, "\"Run\"") ||
		strings.Contains(doc, "`Run`")

	if !hasRunExample {
		t.Error("gl-architect.md missing 'Run' example in Writing steps action verbs — contract C-75 lists 'Run...' as an example")
	}
}

func TestArchitectMd_WritingStepsGuidanceIncludesOpenExample(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract lists "Open..." as an example action verb.
	hasOpenExample := strings.Contains(doc, "Open...") ||
		strings.Contains(doc, "\"Open\"") ||
		strings.Contains(doc, "`Open`")

	if !hasOpenExample {
		t.Error("gl-architect.md missing 'Open' example in Writing steps action verbs — contract C-75 lists 'Open...' as an example")
	}
}

func TestArchitectMd_WritingStepsGuidanceIncludesClickExample(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract lists "Click..." as an example action verb.
	hasClickExample := strings.Contains(doc, "Click...") ||
		strings.Contains(doc, "\"Click\"") ||
		strings.Contains(doc, "`Click`")

	if !hasClickExample {
		t.Error("gl-architect.md missing 'Click' example in Writing steps action verbs — contract C-75 lists 'Click...' as an example")
	}
}

func TestArchitectMd_WritingStepsGuidanceMentionsCommandsAndURLs(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract: "Include commands, URLs, or navigation paths"
	hasCommandsURLsGuidance := strings.Contains(doc, "commands, URLs") ||
		strings.Contains(doc, "commands, urls") ||
		strings.Contains(doc, "navigation paths") ||
		(strings.Contains(doc, "commands") && strings.Contains(doc, "URLs"))

	if !hasCommandsURLsGuidance {
		t.Error("gl-architect.md missing commands/URLs guidance in Writing steps — contract C-75 requires steps to include commands, URLs, or navigation paths")
	}
}

func TestArchitectMd_WritingStepsGuidanceStatesStepsAreOptional(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract: "Steps are optional -- omit when criteria are self-explanatory"
	hasStepsOptionalGuidance := strings.Contains(doc, "Steps are optional") ||
		strings.Contains(doc, "steps are optional") ||
		strings.Contains(doc, "omit when criteria are self-explanatory") ||
		strings.Contains(doc, "self-explanatory")

	if !hasStepsOptionalGuidance {
		t.Error("gl-architect.md missing 'Steps are optional' guidance — contract C-75 states steps can be omitted when criteria are self-explanatory")
	}
}

// --- Output Checklist additions ---

func TestArchitectMd_ChecklistItemEveryContractHasVerificationTier(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract: "[ ] Every contract has a verification tier (auto or verify)"
	hasEveryContractTierItem := strings.Contains(doc, "Every contract has a verification tier") ||
		strings.Contains(doc, "every contract has a verification tier")

	if !hasEveryContractTierItem {
		t.Error("gl-architect.md missing output checklist item for 'Every contract has a verification tier' — contract C-75 requires this checklist addition")
	}
}

func TestArchitectMd_ChecklistItemVerifyTierContractsHaveCriterionOrStep(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract: "[ ] verify-tier contracts have at least one acceptance criterion or step"
	hasVerifyTierCriterionItem := strings.Contains(doc, "verify-tier contracts have at least one acceptance criterion") ||
		strings.Contains(doc, "verify-tier contracts have at least one") ||
		(strings.Contains(doc, "verify-tier") && strings.Contains(doc, "acceptance criterion"))

	if !hasVerifyTierCriterionItem {
		t.Error("gl-architect.md missing output checklist item for 'verify-tier contracts have at least one acceptance criterion or step' — contract C-75 requires this checklist addition")
	}
}

func TestArchitectMd_ChecklistItemAutoTierContractsHaveReasonForSkipping(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract: "[ ] auto-tier contracts have a clear reason for skipping human verification"
	hasAutoTierReasonItem := strings.Contains(doc, "auto-tier contracts have a clear reason") ||
		strings.Contains(doc, "auto-tier contracts have") && strings.Contains(doc, "skipping human verification") ||
		strings.Contains(doc, "reason for skipping human verification")

	if !hasAutoTierReasonItem {
		t.Error("gl-architect.md missing output checklist item for 'auto-tier contracts have a clear reason for skipping human verification' — contract C-75 requires this checklist addition")
	}
}

func TestArchitectMd_ChecklistItemAcceptanceCriteriaAreBehavioral(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract: "[ ] Acceptance criteria are behavioral (what user observes), not implementation"
	hasBehavioralCriteriaItem := strings.Contains(doc, "Acceptance criteria are behavioral") ||
		strings.Contains(doc, "acceptance criteria are behavioral") ||
		strings.Contains(doc, "Acceptance criteria are behavioural")

	if !hasBehavioralCriteriaItem {
		t.Error("gl-architect.md missing output checklist item for 'Acceptance criteria are behavioral (what user observes)' — contract C-75 requires this checklist addition")
	}
}

func TestArchitectMd_ChecklistContainsVerificationTierItemsInOutputChecklistSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// The checklist items must appear within the output_checklist section, not somewhere else.
	checklistPos := strings.Index(doc, "output_checklist")
	if checklistPos == -1 {
		t.Fatal("gl-architect.md missing output_checklist section")
	}

	verificationTierItemPos := strings.Index(doc[checklistPos:], "Every contract has a verification tier")
	if verificationTierItemPos == -1 {
		t.Error("gl-architect.md output_checklist section missing 'Every contract has a verification tier' item — contract C-75 requires the checklist additions to appear inside output_checklist")
	}
}

// --- Error handling ---

func TestArchitectMd_MissingTierOnContractErrorDocumented(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract error: MissingTierOnContract — defaults to verify, checklist catches this as a warning.
	hasMissingTierError := strings.Contains(doc, "MissingTierOnContract") ||
		strings.Contains(doc, "missing tier") ||
		strings.Contains(doc, "contract without verification")

	if !hasMissingTierError {
		t.Error("gl-architect.md missing MissingTierOnContract error documentation — contract C-75 requires this error state to be documented")
	}
}

func TestArchitectMd_MissingTierDefaultsToVerify(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract: MissingTierOnContract behaviour = "Defaults to verify. Output checklist catches this as a warning"
	hasMissingTierDefaultsToVerify := strings.Contains(doc, "Defaults to verify") ||
		strings.Contains(doc, "defaults to verify")

	if !hasMissingTierDefaultsToVerify {
		t.Error("gl-architect.md missing documentation that MissingTierOnContract defaults to verify — contract C-75 requires this default behaviour to be stated")
	}
}

func TestArchitectMd_MissingTierCaughtByChecklist(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract: "Output checklist catches this as a warning"
	hasChecklistCatchesWarning := strings.Contains(doc, "checklist catches") ||
		strings.Contains(doc, "Output checklist catches") ||
		strings.Contains(doc, "checklist catches this as a warning")

	if !hasChecklistCatchesWarning {
		t.Error("gl-architect.md missing documentation that the output checklist catches MissingTierOnContract as a warning — contract C-75 requires this")
	}
}

func TestArchitectMd_TooManyCriteriaErrorDocumented(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract error: TooManyCriteria — suggest splitting the contract. Not blocking.
	hasTooManyCriteriaError := strings.Contains(doc, "TooManyCriteria") ||
		strings.Contains(doc, "more than 5 acceptance criteria") ||
		strings.Contains(doc, "too many criteria")

	if !hasTooManyCriteriaError {
		t.Error("gl-architect.md missing TooManyCriteria error documentation — contract C-75 requires this error state to be documented")
	}
}

func TestArchitectMd_TooManyCriteriaSuggestsSplitting(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract: TooManyCriteria behaviour = "Suggest splitting the contract. Not blocking -- just a guideline"
	hasSplitSuggestion := strings.Contains(doc, "Suggest splitting") ||
		strings.Contains(doc, "suggest splitting") ||
		strings.Contains(doc, "splitting the contract")

	if !hasSplitSuggestion {
		t.Error("gl-architect.md missing 'suggest splitting the contract' guidance for TooManyCriteria — contract C-75 requires this non-blocking suggestion")
	}
}

func TestArchitectMd_TooManyCriteriaIsNotBlocking(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract: TooManyCriteria is "Not blocking -- just a guideline"
	hasNotBlockingStatement := strings.Contains(doc, "Not blocking") ||
		strings.Contains(doc, "not blocking") ||
		strings.Contains(doc, "just a guideline") ||
		strings.Contains(doc, "non-blocking")

	if !hasNotBlockingStatement {
		t.Error("gl-architect.md missing 'not blocking' documentation for TooManyCriteria — contract C-75 requires TooManyCriteria to be a non-blocking guideline, not an error")
	}
}

// --- Invariants / cross-cutting ---

func TestArchitectMd_GuidanceIsNonPrescriptive(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract invariant: "Guidance is non-prescriptive: the architect can override with good reasoning"
	hasNonPrescriptiveStatement := strings.Contains(doc, "non-prescriptive") ||
		strings.Contains(doc, "can override") ||
		strings.Contains(doc, "architect can override")

	if !hasNonPrescriptiveStatement {
		t.Error("gl-architect.md missing non-prescriptive guidance statement — contract C-75 invariant requires the architect to be able to override tier selection with good reasoning")
	}
}

func TestArchitectMd_TierSelectionSectionIsInArchitectFile(t *testing.T) {
	// Verify the guidance is in gl-architect.md and NOT in another agent file.
	architectContent, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	if !strings.Contains(string(architectContent), "Verification Tier Selection") {
		t.Error("Verification Tier Selection section must appear in gl-architect.md — contract C-75 boundary is gl-architect.md")
	}
}

func TestArchitectMd_TierSelectionSectionNotInTestWriterFile(t *testing.T) {
	// The Verification Tier Selection section belongs in gl-architect.md only.
	// gl-test-writer.md should not contain this architect-specific guidance section.
	testWriterContent, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-test-writer.md"))
	if err != nil {
		t.Fatalf("failed to read gl-test-writer.md: %v", err)
	}

	// The full heading should not appear in the test writer agent file.
	if strings.Contains(string(testWriterContent), "## Verification Tier Selection") {
		t.Error("'## Verification Tier Selection' section found in gl-test-writer.md — this architect-specific guidance must only live in gl-architect.md")
	}
}

func TestArchitectMd_VerificationTierSelectionAppearsAfterContractFormatSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// The Verification Tier Selection section is a new section added to the architect agent.
	// It must appear somewhere in the document (already tested), and it should be an additive
	// section — appearing after the contract format section that it extends.
	contractFormatPos := strings.Index(doc, "contract_format")
	if contractFormatPos == -1 {
		t.Fatal("gl-architect.md missing contract_format section")
	}

	tierSelectionPos := strings.Index(doc, "Verification Tier Selection")
	if tierSelectionPos == -1 {
		t.Fatal("gl-architect.md missing 'Verification Tier Selection' section")
	}

	if tierSelectionPos < contractFormatPos {
		t.Error("'Verification Tier Selection' section appears before contract_format in gl-architect.md — it should be an additive section that follows the contract format definition")
	}
}

func TestArchitectMd_EveryContractMustIncludeVerificationFieldInvariant(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract invariant: "Every contract produced by the architect includes a verification field"
	hasEveryContractInvariant := strings.Contains(doc, "Every contract produced by the architect") ||
		strings.Contains(doc, "every contract produced by the architect") ||
		strings.Contains(doc, "Every contract has a verification tier")

	if !hasEveryContractInvariant {
		t.Error("gl-architect.md missing invariant that every contract produced by the architect includes a verification field — contract C-75 requires this invariant to be documented")
	}
}

func TestArchitectMd_AutoTierRequiresJustification(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Contract invariant: "Auto tier requires justification (why tests alone capture correctness)"
	hasAutoJustificationRequirement := strings.Contains(doc, "Auto tier requires justification") ||
		strings.Contains(doc, "auto-tier contracts have a clear reason") ||
		strings.Contains(doc, "tests alone capture correctness") ||
		strings.Contains(doc, "tests pass") && strings.Contains(doc, "fully captures correctness")

	if !hasAutoJustificationRequirement {
		t.Error("gl-architect.md missing auto tier justification requirement — contract C-75 invariant states auto tier requires justification for why tests alone capture correctness")
	}
}
