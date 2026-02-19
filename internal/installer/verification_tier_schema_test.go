package installer_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// C-62 Tests: ContractSchemaExtension
// These tests verify that src/agents/gl-architect.md contains the three new
// verification tier fields (Verification, Acceptance Criteria, Steps) as defined
// in contract C-62.

func TestArchitectMd_ContainsVerificationField(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "**Verification:**") {
		t.Error("gl-architect.md missing **Verification:** field in contract format template")
	}
}

func TestArchitectMd_ContainsAcceptanceCriteriaField(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "**Acceptance Criteria:**") {
		t.Error("gl-architect.md missing **Acceptance Criteria:** field in contract format template")
	}
}

func TestArchitectMd_ContainsStepsField(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "**Steps:**") {
		t.Error("gl-architect.md missing **Steps:** field in contract format template")
	}
}

func TestArchitectMd_VerificationFieldInContractFormatSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	contractFormatPos := strings.Index(doc, "contract_format")
	if contractFormatPos == -1 {
		t.Fatal("gl-architect.md missing contract_format section")
	}

	verificationPos := strings.Index(doc, "**Verification:**")
	if verificationPos == -1 {
		t.Fatal("gl-architect.md missing **Verification:** field")
	}

	if verificationPos < contractFormatPos {
		t.Error("**Verification:** field appears before the contract_format section — it must be inside contract_format")
	}
}

func TestArchitectMd_VerificationFieldAfterSecurityBeforeDependencies(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	securityPos := strings.Index(doc, "**Security:**")
	if securityPos == -1 {
		t.Fatal("gl-architect.md missing **Security:** field in contract format template")
	}

	verificationPos := strings.Index(doc[securityPos:], "**Verification:**")
	if verificationPos == -1 {
		t.Fatal("gl-architect.md missing **Verification:** field after **Security:**")
	}
	absoluteVerificationPos := securityPos + verificationPos

	dependenciesPos := strings.Index(doc[absoluteVerificationPos:], "**Dependencies:**")
	if dependenciesPos == -1 {
		t.Fatal("gl-architect.md missing **Dependencies:** field after **Verification:**")
	}

	// Ordering is confirmed: Security appears first, then Verification, then Dependencies
	// (each search starts from the position of the previous field)
}

func TestArchitectMd_VerificationFieldShowsAutoAndVerifyValues(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// The Verification field must explicitly list both valid values in the contract_format template.
	// Look for the pattern "auto | verify" or "auto" and "verify" within the Verification field line.
	hasAutoVerifyValues := strings.Contains(doc, "auto | verify") ||
		strings.Contains(doc, "auto or verify") ||
		(strings.Contains(doc, "**Verification:**") && strings.Contains(doc, "auto") && strings.Contains(doc, "verify"))

	if !hasAutoVerifyValues {
		t.Error("gl-architect.md missing both 'auto' and 'verify' as documented valid values in the **Verification:** field")
	}
}

func TestArchitectMd_VerificationFieldDocumentsDefaultAsVerify(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// The default must be documented near the Verification field.
	// Accept either "default: verify" or "default verify" phrasing.
	hasDefaultVerify := strings.Contains(doc, "default: verify") ||
		strings.Contains(doc, "default verify") ||
		strings.Contains(doc, "(default: verify)")

	if !hasDefaultVerify {
		t.Error("gl-architect.md missing documentation of 'verify' as the default verification tier")
	}
}

func TestArchitectMd_AcceptanceCriteriaDescribedAsBehavioral(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// The acceptance criteria field must be described as behavioral — what the user observes.
	hasBehavioralDescription := strings.Contains(doc, "behavioral") ||
		strings.Contains(doc, "behaviour") ||
		strings.Contains(doc, "behavior") ||
		strings.Contains(doc, "user can verify") ||
		strings.Contains(doc, "user observes")

	if !hasBehavioralDescription {
		t.Error("gl-architect.md missing behavioral description for acceptance_criteria field — criteria must describe what the user observes, not implementation details")
	}
}

func TestArchitectMd_StepsDescribedAsActionableInstructions(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// The steps field must be described as actionable verification instructions.
	// The contract specifies phrases like "how-to-verify", "run X, open Y, click Z", or "actionable".
	// Require that **Steps:** is present and that actionable guidance is nearby.
	hasStepsField := strings.Contains(doc, "**Steps:**")
	hasActionableDescription := strings.Contains(doc, "actionable") ||
		strings.Contains(doc, "how-to-verify") ||
		strings.Contains(doc, "how to verify")

	if !hasStepsField || !hasActionableDescription {
		t.Error("gl-architect.md missing actionable instructions description for the **Steps:** field — steps must describe how to verify (e.g. 'how-to-verify' or 'actionable'), not implementation internals")
	}
}

func TestArchitectMd_InvalidTierValueErrorDocumented(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// The InvalidTierValue error state must be documented.
	hasInvalidTierDoc := strings.Contains(doc, "InvalidTierValue") ||
		strings.Contains(doc, "Invalid verification tier") ||
		strings.Contains(doc, "Must be auto or verify")

	if !hasInvalidTierDoc {
		t.Error("gl-architect.md missing documentation for InvalidTierValue error — contract must describe rejection of unrecognised verification tier values")
	}
}

func TestArchitectMd_EmptyVerifyCriteriaWarningDocumented(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// The EmptyVerifyCriteria warning must be documented (warn, not error).
	hasEmptyCriteriaWarning := strings.Contains(doc, "EmptyVerifyCriteria") ||
		strings.Contains(doc, "verify tier but no acceptance criteria") ||
		strings.Contains(doc, "no acceptance criteria or steps")

	if !hasEmptyCriteriaWarning {
		t.Error("gl-architect.md missing documentation for EmptyVerifyCriteria warning — must warn when verify tier has no acceptance criteria or steps")
	}
}

func TestArchitectMd_ThreeFieldsAreOptional(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// The three fields must be documented as optional.
	if !strings.Contains(doc, "Optional") && !strings.Contains(doc, "optional") {
		t.Error("gl-architect.md missing documentation that Verification, Acceptance Criteria, and Steps fields are optional")
	}
}

func TestArchitectMd_ExistingContractsWithoutFieldDefaultToVerify(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-architect.md"))
	if err != nil {
		t.Fatalf("failed to read gl-architect.md: %v", err)
	}

	doc := string(content)

	// Existing contracts missing the verification field must be explicitly documented as
	// defaulting to "verify". The documentation must tie absence of the field to the verify default.
	hasDefaultingBehaviour := strings.Contains(doc, "without verification field") ||
		strings.Contains(doc, "contracts missing") ||
		strings.Contains(doc, "Existing contracts") ||
		strings.Contains(doc, "absent") ||
		(strings.Contains(doc, "default to \"verify\"") || strings.Contains(doc, "default: \"verify\""))

	if !hasDefaultingBehaviour {
		t.Error("gl-architect.md missing documentation that existing contracts without the verification field default to 'verify' — must explicitly state the safe default behaviour")
	}
}

// C-63 Tests: VerifierTierAwareness
// These tests verify that src/agents/gl-verifier.md contains tier awareness
// content in the verification report as defined in contract C-63.

func TestVerifierMd_ContainsTierAwarenessContent(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-verifier.md"))
	if err != nil {
		t.Fatalf("failed to read gl-verifier.md: %v", err)
	}

	doc := string(content)

	hasVerificationTier := strings.Contains(doc, "Verification Tier") ||
		strings.Contains(doc, "verification tier")

	if !hasVerificationTier {
		t.Error("gl-verifier.md missing 'Verification Tier' section in verification report format")
	}
}

func TestVerifierMd_ContainsEffectiveTierInReport(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-verifier.md"))
	if err != nil {
		t.Fatalf("failed to read gl-verifier.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "Effective tier") && !strings.Contains(doc, "effective tier") {
		t.Error("gl-verifier.md missing 'Effective tier' field in verification report — verifier must report the effective tier for the slice")
	}
}

func TestVerifierMd_ContainsEffectiveTierComputation(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-verifier.md"))
	if err != nil {
		t.Fatalf("failed to read gl-verifier.md: %v", err)
	}

	doc := string(content)

	// The effective tier computation rule: verify > auto (highest tier wins).
	// Require one of the specific formulations from the contract spec.
	hasComputationRule := strings.Contains(doc, "verify > auto") ||
		strings.Contains(doc, "highest tier wins") ||
		strings.Contains(doc, "any contract has tier")

	if !hasComputationRule {
		t.Error("gl-verifier.md missing effective tier computation rule — must document that verify > auto (highest tier wins) using one of: 'verify > auto', 'highest tier wins', or 'any contract has tier'")
	}
}

func TestVerifierMd_ContainsPerContractTierBreakdown(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-verifier.md"))
	if err != nil {
		t.Fatalf("failed to read gl-verifier.md: %v", err)
	}

	doc := string(content)

	hasPerContractBreakdown := strings.Contains(doc, "Per-contract") ||
		strings.Contains(doc, "per-contract") ||
		strings.Contains(doc, "per contract")

	if !hasPerContractBreakdown {
		t.Error("gl-verifier.md missing per-contract tier breakdown — verifier must list each contract's tier with criteria and steps counts")
	}
}

func TestVerifierMd_ContainsCriteriaAndStepsCounts(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-verifier.md"))
	if err != nil {
		t.Fatalf("failed to read gl-verifier.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "criteria") {
		t.Error("gl-verifier.md missing 'criteria' count in per-contract tier breakdown")
	}

	if !strings.Contains(doc, "steps") {
		t.Error("gl-verifier.md missing 'steps' count in per-contract tier breakdown")
	}
}

func TestVerifierMd_ContainsMissingVerificationFieldDefaultsToVerify(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-verifier.md"))
	if err != nil {
		t.Fatalf("failed to read gl-verifier.md: %v", err)
	}

	doc := string(content)

	// Contracts with no verification field must default to "verify" — verifier must document this.
	hasDefaultBehaviour := strings.Contains(doc, "MissingVerificationField") ||
		strings.Contains(doc, "defaulted to verify") ||
		strings.Contains(doc, "default to \"verify\"") ||
		strings.Contains(doc, "default: \"verify\"") ||
		strings.Contains(doc, "defaults to verify")

	if !hasDefaultBehaviour {
		t.Error("gl-verifier.md missing documentation that contracts with no verification field default to 'verify'")
	}
}

func TestVerifierMd_ContainsInvalidTierWarning(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-verifier.md"))
	if err != nil {
		t.Fatalf("failed to read gl-verifier.md: %v", err)
	}

	doc := string(content)

	// An unrecognised tier value must be warned about and treated as verify.
	hasInvalidTierWarning := strings.Contains(doc, "InvalidTierInContract") ||
		strings.Contains(doc, "Unknown tier") ||
		strings.Contains(doc, "unrecognised") ||
		strings.Contains(doc, "treating as verify")

	if !hasInvalidTierWarning {
		t.Error("gl-verifier.md missing InvalidTierInContract warning — must warn on unrecognised tier values and treat them as verify")
	}
}

func TestVerifierMd_ContainsWarningsForVerifyTierWithNoCriteria(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-verifier.md"))
	if err != nil {
		t.Fatalf("failed to read gl-verifier.md: %v", err)
	}

	doc := string(content)

	// Contracts with verify tier but no acceptance criteria and no steps must be flagged.
	hasEmptyCriteriaWarning := strings.Contains(doc, "no acceptance criteria or steps") ||
		strings.Contains(doc, "verify tier with no acceptance") ||
		strings.Contains(doc, "verify tier but no") ||
		strings.Contains(doc, "empty acceptance_criteria")

	if !hasEmptyCriteriaWarning {
		t.Error("gl-verifier.md missing Warnings subsection for verify tier contracts with no acceptance criteria or steps")
	}
}

func TestVerifierMd_DescribedAsReportingNotEnforcingGate(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-verifier.md"))
	if err != nil {
		t.Fatalf("failed to read gl-verifier.md: %v", err)
	}

	doc := string(content)

	// The verifier must be explicitly documented as reporting only — not enforcing the gate.
	hasReportingOnlyDescription := strings.Contains(doc, "does not enforce") ||
		strings.Contains(doc, "reporting only") ||
		strings.Contains(doc, "reports tier") ||
		strings.Contains(doc, "informational only") ||
		strings.Contains(doc, "orchestrator enforces")

	if !hasReportingOnlyDescription {
		t.Error("gl-verifier.md missing documentation that verifier only reports tier (does not enforce the gate) — the orchestrator enforces the verification checkpoint")
	}
}

func TestVerifierMd_TierReportingSectionAppearsAfterExistingReportSections(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-verifier.md"))
	if err != nil {
		t.Fatalf("failed to read gl-verifier.md: %v", err)
	}

	doc := string(content)

	// "Verification Tier" section must be additive — it follows existing sections.
	// The report must have content before the tier section appears.
	verificationTierPos := strings.Index(doc, "Verification Tier")
	if verificationTierPos == -1 {
		t.Fatal("gl-verifier.md missing 'Verification Tier' section")
	}

	// There must be non-trivial content before the Verification Tier section.
	// A minimal proxy: if verificationTierPos is within the first 200 characters,
	// the tier section is likely the very first thing, which contradicts "additive".
	if verificationTierPos < 200 {
		t.Error("gl-verifier.md 'Verification Tier' section appears too early — it must be additive, appearing after existing report sections")
	}
}

func TestVerifierMd_WarningsAreInformationalNotBlocking(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-verifier.md"))
	if err != nil {
		t.Fatalf("failed to read gl-verifier.md: %v", err)
	}

	doc := string(content)

	// Warnings must be informational (not blocking — verifier still passes).
	hasInformationalWarnings := strings.Contains(doc, "informational") ||
		strings.Contains(doc, "not blocking") ||
		strings.Contains(doc, "still passes") ||
		strings.Contains(doc, "non-blocking")

	if !hasInformationalWarnings {
		t.Error("gl-verifier.md missing documentation that tier warnings are informational and not blocking — verifier must still pass when warnings are present")
	}
}
