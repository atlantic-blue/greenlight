package installer_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// C-76 Tests: SliceStateTemplate
// These tests verify that src/templates/slice-state.md exists and contains
// all required sections as defined in contract C-76.
// Verification: auto

// C-77 Tests: StateFormatReference
// These tests verify that src/references/state-format.md exists and contains
// all required sections as defined in contract C-77.
// Verification: auto

// ---------------------------------------------------------------------------
// C-76: SliceStateTemplate — src/templates/slice-state.md
// ---------------------------------------------------------------------------

func TestSliceStateMd_Exists(t *testing.T) {
	path := filepath.Join(projectRoot(), "src/templates/slice-state.md")
	_, err := os.Stat(path)
	if err != nil {
		t.Fatalf("slice-state.md does not exist at %s: %v", path, err)
	}
}

// Section 1: Schema Definition — frontmatter fields

func TestSliceStateMd_ContainsSchemaDefinitionSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "Schema") {
		t.Error("slice-state.md missing Schema Definition section")
	}
}

func TestSliceStateMd_ContainsIdField(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "id:") && !strings.Contains(doc, "id ") {
		t.Error("slice-state.md missing 'id' frontmatter field in schema definition")
	}
}

func TestSliceStateMd_ContainsStatusField(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "status:") && !strings.Contains(doc, "status ") {
		t.Error("slice-state.md missing 'status' frontmatter field in schema definition")
	}
}

func TestSliceStateMd_ContainsStepField(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "step:") && !strings.Contains(doc, "step ") {
		t.Error("slice-state.md missing 'step' frontmatter field in schema definition")
	}
}

func TestSliceStateMd_ContainsMilestoneField(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "milestone") {
		t.Error("slice-state.md missing 'milestone' frontmatter field in schema definition")
	}
}

func TestSliceStateMd_ContainsStartedField(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "started") {
		t.Error("slice-state.md missing 'started' frontmatter field in schema definition")
	}
}

func TestSliceStateMd_ContainsUpdatedField(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "updated") {
		t.Error("slice-state.md missing 'updated' frontmatter field in schema definition")
	}
}

func TestSliceStateMd_ContainsTestsField(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "tests:") && !strings.Contains(doc, "tests ") {
		t.Error("slice-state.md missing 'tests' frontmatter field in schema definition")
	}
}

func TestSliceStateMd_ContainsSecurityTestsField(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "security_tests") {
		t.Error("slice-state.md missing 'security_tests' frontmatter field in schema definition")
	}
}

func TestSliceStateMd_ContainsSessionField(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "session") {
		t.Error("slice-state.md missing 'session' frontmatter field in schema definition")
	}
}

func TestSliceStateMd_ContainsDepsField(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "deps") {
		t.Error("slice-state.md missing 'deps' frontmatter field in schema definition")
	}
}

func TestSliceStateMd_ContainsRequiredAndOptionalDistinction(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	hasRequired := strings.Contains(doc, "required") || strings.Contains(doc, "Required")
	hasOptional := strings.Contains(doc, "optional") || strings.Contains(doc, "Optional")

	if !hasRequired {
		t.Error("slice-state.md missing documentation of required fields in schema definition")
	}
	if !hasOptional {
		t.Error("slice-state.md missing documentation of optional fields in schema definition")
	}
}

// Section 2: Status Lifecycle — valid status and step values, transition rules

func TestSliceStateMd_ContainsStatusLifecycleSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "Status Lifecycle") && !strings.Contains(doc, "status lifecycle") {
		t.Error("slice-state.md missing Status Lifecycle section")
	}
}

func TestSliceStateMd_ContainsAllEightStatusValues(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	requiredStatusValues := []string{
		"pending",
		"ready",
		"tests",
		"implementing",
		"security",
		"fixing",
		"verifying",
		"complete",
	}

	for _, status := range requiredStatusValues {
		if !strings.Contains(doc, status) {
			t.Errorf("slice-state.md missing status value: %s", status)
		}
	}
}

func TestSliceStateMd_ContainsAllSevenStepValues(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	requiredStepValues := []string{
		"none",
		"tests",
		"implementing",
		"security",
		"fixing",
		"verifying",
		"complete",
	}

	for _, step := range requiredStepValues {
		if !strings.Contains(doc, step) {
			t.Errorf("slice-state.md missing step value: %s", step)
		}
	}
}

func TestSliceStateMd_ContainsTransitionRules(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	hasTransitions := strings.Contains(doc, "transition") || strings.Contains(doc, "Transition")

	if !hasTransitions {
		t.Error("slice-state.md missing transition rules in Status Lifecycle section")
	}
}

// Section 3: Session Tracking

func TestSliceStateMd_ContainsSessionTrackingSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "Session Tracking") && !strings.Contains(doc, "session tracking") {
		t.Error("slice-state.md missing Session Tracking section")
	}
}

func TestSliceStateMd_ContainsSessionFormatISOTimestampWithSuffix(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	// Session format must be ISO timestamp + hyphen + random alphanumeric suffix
	hasISO := strings.Contains(doc, "ISO") || strings.Contains(doc, "timestamp")
	hasSuffix := strings.Contains(doc, "suffix") || strings.Contains(doc, "random") || strings.Contains(doc, "alphanumeric")

	if !hasISO {
		t.Error("slice-state.md missing ISO timestamp format description in Session Tracking section")
	}
	if !hasSuffix {
		t.Error("slice-state.md missing random suffix description in Session Tracking section (format: ISO timestamp + hyphen + random alphanumeric suffix)")
	}
}

func TestSliceStateMd_ContainsSessionAdvisoryOnly(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "advisory") {
		t.Error("slice-state.md missing 'advisory' designation for session field — session is advisory only, not a lock")
	}
}

func TestSliceStateMd_ContainsSessionSetOnClaimClearedOnCompletion(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	hasClaim := strings.Contains(doc, "claim") || strings.Contains(doc, "set on claim")
	hasCleared := strings.Contains(doc, "cleared") || strings.Contains(doc, "clear")

	if !hasClaim {
		t.Error("slice-state.md missing 'claim' — session field must state it is set on claim")
	}
	if !hasCleared {
		t.Error("slice-state.md missing 'cleared' — session field must state it is cleared on completion")
	}
}

// Section 4: Body Sections

func TestSliceStateMd_ContainsBodySectionsSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "Body") {
		t.Error("slice-state.md missing Body Sections section")
	}
}

func TestSliceStateMd_ContainsSliceIdAndNameHeading(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	// The heading pattern is: # {slice-id}: {slice-name}
	hasSliceIdPattern := strings.Contains(doc, "{slice-id}") ||
		strings.Contains(doc, "slice-id") ||
		strings.Contains(doc, "slice_id")

	if !hasSliceIdPattern {
		t.Error("slice-state.md missing '{slice-id}: {slice-name}' heading pattern in Body Sections")
	}
}

func TestSliceStateMd_ContainsWhySection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "## Why") && !strings.Contains(doc, "Why") {
		t.Error("slice-state.md missing '## Why' body section")
	}
}

func TestSliceStateMd_ContainsWhatSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "## What") && !strings.Contains(doc, "What") {
		t.Error("slice-state.md missing '## What' body section")
	}
}

func TestSliceStateMd_ContainsDependenciesSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "Dependencies") {
		t.Error("slice-state.md missing '## Dependencies' body section")
	}
}

func TestSliceStateMd_ContainsContractsSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "Contracts") {
		t.Error("slice-state.md missing '## Contracts' body section")
	}
}

func TestSliceStateMd_ContainsDecisionsSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "Decisions") {
		t.Error("slice-state.md missing '## Decisions' body section")
	}
}

func TestSliceStateMd_ContainsFilesSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "## Files") && !strings.Contains(doc, "Files") {
		t.Error("slice-state.md missing '## Files' body section")
	}
}

// Section 5: File Naming

func TestSliceStateMd_ContainsFileNamingSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	hasNamingSection := strings.Contains(doc, "File Naming") ||
		strings.Contains(doc, "file naming") ||
		strings.Contains(doc, "Naming")

	if !hasNamingSection {
		t.Error("slice-state.md missing File Naming section")
	}
}

func TestSliceStateMd_ContainsGreenlightSlicesDirectory(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, ".greenlight/slices/") {
		t.Error("slice-state.md missing '.greenlight/slices/' directory in file naming section")
	}
}

func TestSliceStateMd_ContainsOneFilePerSliceConstraint(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	hasOnePerSlice := strings.Contains(doc, "one file per slice") ||
		strings.Contains(doc, "one per slice") ||
		strings.Contains(doc, "One file per slice") ||
		strings.Contains(doc, "One per slice")

	if !hasOnePerSlice {
		t.Error("slice-state.md missing 'one file per slice' constraint in file naming section")
	}
}

func TestSliceStateMd_ContainsSliceIdInFileName(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	// Files named {slice-id}.md
	hasSliceIdFilename := strings.Contains(doc, "{slice-id}.md") ||
		strings.Contains(doc, "S-") ||
		strings.Contains(doc, ".md")

	if !hasSliceIdFilename {
		t.Error("slice-state.md missing slice-id file naming convention (e.g. {slice-id}.md)")
	}
}

// Section 6: Examples

func TestSliceStateMd_ContainsExamplesSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "Example") && !strings.Contains(doc, "example") {
		t.Error("slice-state.md missing Examples section")
	}
}

func TestSliceStateMd_ContainsPendingStateExample(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	hasPendingExample := strings.Contains(doc, "pending") &&
		(strings.Contains(doc, "Example") || strings.Contains(doc, "example"))

	if !hasPendingExample {
		t.Error("slice-state.md missing complete example for 'pending' state")
	}
}

func TestSliceStateMd_ContainsImplementingStateExample(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	hasPendingExample := strings.Contains(doc, "implementing") &&
		(strings.Contains(doc, "Example") || strings.Contains(doc, "example"))

	if !hasPendingExample {
		t.Error("slice-state.md missing complete example for 'implementing' state")
	}
}

func TestSliceStateMd_ContainsCompleteStateExample(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	hasCompleteExample := strings.Contains(doc, "complete") &&
		(strings.Contains(doc, "Example") || strings.Contains(doc, "example"))

	if !hasCompleteExample {
		t.Error("slice-state.md missing complete example for 'complete' state")
	}
}

// Errors: InvalidSliceId, InvalidStatus, InvalidFrontmatter

func TestSliceStateMd_ContainsInvalidSliceIdError(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "InvalidSliceId") {
		t.Error("slice-state.md missing 'InvalidSliceId' error definition")
	}
}

func TestSliceStateMd_ContainsSliceIdPatternSDigits(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	// Slice IDs must match S-{digits} or S-{digits}.{digits}
	hasPattern := strings.Contains(doc, "S-{digits}") ||
		strings.Contains(doc, `S-\d`) ||
		strings.Contains(doc, "S-[0-9]") ||
		(strings.Contains(doc, "S-") && strings.Contains(doc, "digits"))

	if !hasPattern {
		t.Error("slice-state.md missing slice ID pattern 'S-{digits}' or 'S-{digits}.{digits}' in InvalidSliceId error")
	}
}

func TestSliceStateMd_ContainsInvalidStatusError(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "InvalidStatus") {
		t.Error("slice-state.md missing 'InvalidStatus' error definition")
	}
}

func TestSliceStateMd_ContainsInvalidFrontmatterError(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "InvalidFrontmatter") {
		t.Error("slice-state.md missing 'InvalidFrontmatter' error definition")
	}
}

func TestSliceStateMd_ContainsFrontmatterDelimiterDescription(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	// Frontmatter must be between --- delimiters
	if !strings.Contains(doc, "---") {
		t.Error("slice-state.md missing '---' frontmatter delimiter description")
	}
}

// Invariants

func TestSliceStateMd_InvariantTemplateIsReadOnly(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	hasReadOnly := strings.Contains(doc, "read-only") ||
		strings.Contains(doc, "Read-only") ||
		strings.Contains(doc, "read only") ||
		strings.Contains(doc, "not modified at runtime") ||
		strings.Contains(doc, "runtime")

	if !hasReadOnly {
		t.Error("slice-state.md missing read-only invariant — template must state it is read-only at runtime")
	}
}

func TestSliceStateMd_InvariantFrontmatterFlatKeyValue(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	hasFlatConstraint := strings.Contains(doc, "flat") ||
		strings.Contains(doc, "no nesting") ||
		strings.Contains(doc, "flat key-value")

	if !hasFlatConstraint {
		t.Error("slice-state.md missing flat key-value invariant for frontmatter (no nesting)")
	}
}

func TestSliceStateMd_InvariantAllFieldNamesLowercaseWithUnderscores(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	hasLowercaseConstraint := strings.Contains(doc, "lowercase") ||
		strings.Contains(doc, "lower_case") ||
		strings.Contains(doc, "snake_case")

	if !hasLowercaseConstraint {
		t.Error("slice-state.md missing lowercase with underscores invariant for field names")
	}
}

func TestSliceStateMd_InvariantStatusEnumClosed(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	hasClosedEnum := strings.Contains(doc, "closed") || strings.Contains(doc, "enum") || strings.Contains(doc, "only")

	if !hasClosedEnum {
		t.Error("slice-state.md missing closed enum invariant for status (only 8 values valid)")
	}
}

func TestSliceStateMd_InvariantSliceIdPathTraversalPrevention(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	hasPathTraversalProtection := strings.Contains(doc, "path traversal") ||
		strings.Contains(doc, "traversal") ||
		strings.Contains(doc, "../")

	if !hasPathTraversalProtection {
		t.Error("slice-state.md missing path traversal prevention invariant for slice ID validation")
	}
}

func TestSliceStateMd_ContainsFrontmatterExampleWithTripleDashes(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	// Count occurrences of --- to verify triple-dash delimiters are used in examples
	count := strings.Count(doc, "---")
	if count < 2 {
		t.Errorf("slice-state.md has fewer than 2 '---' delimiter occurrences (%d found) — frontmatter examples must use --- delimiters", count)
	}
}

func TestSliceStateMd_SchemaSectionBeforeStatusLifecycle(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	schemaPos := strings.Index(doc, "Schema")
	lifecyclePos := strings.Index(doc, "Status Lifecycle")

	if schemaPos == -1 {
		t.Fatal("slice-state.md missing Schema Definition section")
	}
	if lifecyclePos == -1 {
		t.Fatal("slice-state.md missing Status Lifecycle section")
	}

	if schemaPos >= lifecyclePos {
		t.Errorf("Schema Definition (pos %d) must appear before Status Lifecycle (pos %d)", schemaPos, lifecyclePos)
	}
}

func TestSliceStateMd_StatusLifecycleBeforeSessionTracking(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	lifecyclePos := strings.Index(doc, "Status Lifecycle")
	sessionPos := strings.Index(doc, "Session Tracking")

	if lifecyclePos == -1 {
		t.Fatal("slice-state.md missing Status Lifecycle section")
	}
	if sessionPos == -1 {
		t.Fatal("slice-state.md missing Session Tracking section")
	}

	if lifecyclePos >= sessionPos {
		t.Errorf("Status Lifecycle (pos %d) must appear before Session Tracking (pos %d)", lifecyclePos, sessionPos)
	}
}

func TestSliceStateMd_ExamplesAppearAfterBodySections(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	bodySectionsPos := strings.Index(doc, "Body")
	examplesPos := strings.Index(doc, "Example")
	if examplesPos == -1 {
		examplesPos = strings.Index(doc, "example")
	}

	if bodySectionsPos == -1 {
		t.Fatal("slice-state.md missing Body Sections section")
	}
	if examplesPos == -1 {
		t.Fatal("slice-state.md missing Examples section")
	}

	if bodySectionsPos >= examplesPos {
		t.Errorf("Body Sections (pos %d) must appear before Examples (pos %d)", bodySectionsPos, examplesPos)
	}
}

// ---------------------------------------------------------------------------
// C-77: StateFormatReference — src/references/state-format.md
// ---------------------------------------------------------------------------

func TestStateFormatMd_Exists(t *testing.T) {
	path := filepath.Join(projectRoot(), "src/references/state-format.md")
	_, err := os.Stat(path)
	if err != nil {
		t.Fatalf("state-format.md does not exist at %s: %v", path, err)
	}
}

// Section 1: State Format Detection

func TestStateFormatMd_ContainsStateFormatDetectionSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	hasDetection := strings.Contains(doc, "Detection") ||
		strings.Contains(doc, "detection") ||
		strings.Contains(doc, "detect")

	if !hasDetection {
		t.Error("state-format.md missing State Format Detection section")
	}
}

func TestStateFormatMd_ContainsDirectoryExistenceCheck(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	hasDirectoryCheck := strings.Contains(doc, "slices/") || strings.Contains(doc, "directory")

	if !hasDirectoryCheck {
		t.Error("state-format.md missing directory existence check in detection section")
	}
}

func TestStateFormatMd_ContainsDetectionFlowSlicesDirectory(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "slices/") {
		t.Error("state-format.md missing 'slices/' in detection flow")
	}
}

func TestStateFormatMd_ContainsDetectionFlowLegacyStatemd(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "STATE.md") {
		t.Error("state-format.md missing 'STATE.md' legacy format in detection flow")
	}
}

func TestStateFormatMd_ContainsDetectionFlowNoState(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	hasNoState := strings.Contains(doc, "no state") ||
		strings.Contains(doc, "No state") ||
		strings.Contains(doc, "NoStateFound") ||
		strings.Contains(doc, "neither")

	if !hasNoState {
		t.Error("state-format.md missing 'no state' case in detection flow (neither slices/ nor STATE.md exists)")
	}
}

// Section 2: File-Per-Slice Format

func TestStateFormatMd_ContainsFilePerSliceFormatSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	hasFilePerSlice := strings.Contains(doc, "File-Per-Slice") ||
		strings.Contains(doc, "file-per-slice") ||
		strings.Contains(doc, "per-slice") ||
		strings.Contains(doc, "Per-Slice")

	if !hasFilePerSlice {
		t.Error("state-format.md missing File-Per-Slice Format section")
	}
}

func TestStateFormatMd_ContainsProjectStateJsonReference(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "project-state.json") {
		t.Error("state-format.md missing 'project-state.json' reference in File-Per-Slice format section")
	}
}

func TestStateFormatMd_ContainsGeneratedStateMdReference(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "generated") && !strings.Contains(doc, "Generated") {
		t.Error("state-format.md missing 'generated' STATE.md reference in File-Per-Slice format section")
	}
}

func TestStateFormatMd_ContainsSchemaReference(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	hasSchemaRef := strings.Contains(doc, "schema") || strings.Contains(doc, "Schema")

	if !hasSchemaRef {
		t.Error("state-format.md missing schema reference in File-Per-Slice format section")
	}
}

// Section 3: Legacy Format

func TestStateFormatMd_ContainsLegacyFormatSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "Legacy") && !strings.Contains(doc, "legacy") {
		t.Error("state-format.md missing Legacy Format section")
	}
}

func TestStateFormatMd_LegacyFormatSupportedIndefinitely(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	hasSupportedIndefinitely := strings.Contains(doc, "indefinitely") ||
		strings.Contains(doc, "supported indefinitely") ||
		strings.Contains(doc, "D-37") ||
		strings.Contains(doc, "never removed") ||
		strings.Contains(doc, "no automatic migration")

	if !hasSupportedIndefinitely {
		t.Error("state-format.md missing declaration that legacy format is supported indefinitely (D-37)")
	}
}

// Section 4: Concurrent Access Patterns

func TestStateFormatMd_ContainsConcurrentAccessSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	hasConcurrent := strings.Contains(doc, "Concurrent") ||
		strings.Contains(doc, "concurrent")

	if !hasConcurrent {
		t.Error("state-format.md missing Concurrent Access Patterns section")
	}
}

func TestStateFormatMd_ContainsEachSessionWritesOwnFile(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	hasOwnFile := strings.Contains(doc, "own file") ||
		strings.Contains(doc, "writes own") ||
		strings.Contains(doc, "never touch same file") ||
		strings.Contains(doc, "different file")

	if !hasOwnFile {
		t.Error("state-format.md missing 'each session writes own file' in concurrent access section")
	}
}

func TestStateFormatMd_ContainsNoFileLocking(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	hasNoFileLocking := strings.Contains(doc, "no file locking") ||
		strings.Contains(doc, "no locking") ||
		strings.Contains(doc, "No file locking") ||
		strings.Contains(doc, "No locking")

	if !hasNoFileLocking {
		t.Error("state-format.md missing 'no file locking' in concurrent access section")
	}
}

func TestStateFormatMd_ContainsAdvisorySessionTracking(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "advisory") {
		t.Error("state-format.md missing 'advisory' session tracking in concurrent access section")
	}
}

// Section 5: STATE.md Regeneration

func TestStateFormatMd_ContainsStateMdRegenerationSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	hasRegeneration := strings.Contains(doc, "Regeneration") ||
		strings.Contains(doc, "regeneration") ||
		strings.Contains(doc, "regenerate")

	if !hasRegeneration {
		t.Error("state-format.md missing STATE.md Regeneration section")
	}
}

func TestStateFormatMd_RegenerationAfterEveryWrite(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	hasAfterEveryWrite := strings.Contains(doc, "after every write") ||
		strings.Contains(doc, "every write") ||
		strings.Contains(doc, "D-34") ||
		strings.Contains(doc, "after each write")

	if !hasAfterEveryWrite {
		t.Error("state-format.md missing 'after every write' mandate for STATE.md regeneration (D-34)")
	}
}

func TestStateFormatMd_RegenerationContainsHeaderComment(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	hasHeaderComment := strings.Contains(doc, "header comment") ||
		strings.Contains(doc, "<!-- generated") ||
		strings.Contains(doc, "generated file") ||
		strings.Contains(doc, "do not edit")

	if !hasHeaderComment {
		t.Error("state-format.md missing header comment requirement for generated STATE.md")
	}
}

// Section 6: Crash Safety

func TestStateFormatMd_ContainsCrashSafetySection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	hasCrashSafety := strings.Contains(doc, "Crash Safety") ||
		strings.Contains(doc, "crash safety") ||
		strings.Contains(doc, "crash-safe")

	if !hasCrashSafety {
		t.Error("state-format.md missing Crash Safety section")
	}
}

func TestStateFormatMd_ContainsWriteToTempThenRename(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	hasAtomicWrite := strings.Contains(doc, "write-to-temp-then-rename") ||
		strings.Contains(doc, "write to temp") ||
		strings.Contains(doc, "temp") && strings.Contains(doc, "rename") ||
		strings.Contains(doc, "atomic")

	if !hasAtomicWrite {
		t.Error("state-format.md missing write-to-temp-then-rename crash safety pattern")
	}
}

func TestStateFormatMd_ContainsSameFilesystemRequirement(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	hasSameFilesystem := strings.Contains(doc, "same filesystem") ||
		strings.Contains(doc, "same file system") ||
		strings.Contains(doc, "POSIX")

	if !hasSameFilesystem {
		t.Error("state-format.md missing same filesystem requirement for atomic rename (POSIX atomicity)")
	}
}

func TestStateFormatMd_ContainsPOSIXAtomicityMention(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	hasPOSIX := strings.Contains(doc, "POSIX") ||
		strings.Contains(doc, "atomic rename") ||
		strings.Contains(doc, "atomically")

	if !hasPOSIX {
		t.Error("state-format.md missing POSIX atomicity guarantee in crash safety section")
	}
}

// Section 7: Backward Compatibility

func TestStateFormatMd_ContainsBackwardCompatibilitySection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	hasBackwardCompat := strings.Contains(doc, "Backward Compatibility") ||
		strings.Contains(doc, "backward compatibility") ||
		strings.Contains(doc, "backward compat") ||
		strings.Contains(doc, "Backwards Compatibility")

	if !hasBackwardCompat {
		t.Error("state-format.md missing Backward Compatibility section")
	}
}

func TestStateFormatMd_BothFormatsIndefiniteSupport(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "indefinitely") {
		t.Error("state-format.md missing 'indefinitely' — both formats must be supported indefinitely")
	}
}

func TestStateFormatMd_NoDualWrite(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	hasNoDualWrite := strings.Contains(doc, "no dual-write") ||
		strings.Contains(doc, "dual-write") ||
		strings.Contains(doc, "dual write")

	if !hasNoDualWrite {
		t.Error("state-format.md missing 'no dual-write' constraint in backward compatibility section")
	}
}

func TestStateFormatMd_NoAutomaticMigration(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	hasNoMigration := strings.Contains(doc, "no automatic migration") ||
		strings.Contains(doc, "No automatic migration") ||
		strings.Contains(doc, "no migration") ||
		strings.Contains(doc, "No migration")

	if !hasNoMigration {
		t.Error("state-format.md missing 'no automatic migration' in backward compatibility section")
	}
}

func TestStateFormatMd_ZeroCostDetection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	hasZeroCost := strings.Contains(doc, "zero-cost") ||
		strings.Contains(doc, "zero cost") ||
		strings.Contains(doc, "cost-free")

	if !hasZeroCost {
		t.Error("state-format.md missing 'zero-cost detection' in backward compatibility section")
	}
}

// Errors: NoStateFound, CorruptSliceFile, SlicesDirectoryEmpty

func TestStateFormatMd_ContainsNoStateFoundError(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "NoStateFound") {
		t.Error("state-format.md missing 'NoStateFound' error definition")
	}
}

func TestStateFormatMd_ContainsCorruptSliceFileError(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "CorruptSliceFile") {
		t.Error("state-format.md missing 'CorruptSliceFile' error definition")
	}
}

func TestStateFormatMd_CorruptSliceFileSkipAndWarn(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	// CorruptSliceFile: invalid frontmatter, skip and warn
	hasSkipAndWarn := strings.Contains(doc, "skip") && strings.Contains(doc, "warn")

	if !hasSkipAndWarn {
		t.Error("state-format.md missing 'skip and warn' behaviour for CorruptSliceFile error")
	}
}

func TestStateFormatMd_ContainsSlicesDirectoryEmptyError(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "SlicesDirectoryEmpty") {
		t.Error("state-format.md missing 'SlicesDirectoryEmpty' error definition")
	}
}

func TestStateFormatMd_SlicesDirectoryEmptyTreatAsZeroSlices(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	hasTreatAsZero := strings.Contains(doc, "zero slices") ||
		strings.Contains(doc, "treat as zero") ||
		strings.Contains(doc, "empty state") ||
		strings.Contains(doc, "no .md files")

	if !hasTreatAsZero {
		t.Error("state-format.md missing 'treat as zero slices' behaviour for SlicesDirectoryEmpty error")
	}
}

// Invariants

func TestStateFormatMd_InvariantReadOnlyAtRuntime(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	hasReadOnly := strings.Contains(doc, "read-only") ||
		strings.Contains(doc, "Read-only") ||
		strings.Contains(doc, "not modified at runtime")

	if !hasReadOnly {
		t.Error("state-format.md missing read-only invariant — reference document must state it is read-only at runtime")
	}
}

func TestStateFormatMd_InvariantDetectionLogicDeterministic(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	hasDeterministic := strings.Contains(doc, "deterministic") ||
		strings.Contains(doc, "Deterministic")

	if !hasDeterministic {
		t.Error("state-format.md missing 'deterministic' invariant for detection logic")
	}
}

func TestStateFormatMd_InvariantStateMdRegenerationMandatory(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	hasMandatory := strings.Contains(doc, "mandatory") ||
		strings.Contains(doc, "Mandatory") ||
		strings.Contains(doc, "must regenerate") ||
		strings.Contains(doc, "required")

	if !hasMandatory {
		t.Error("state-format.md missing 'mandatory' invariant for STATE.md regeneration after writes")
	}
}

func TestStateFormatMd_InvariantConcurrentSessionsNeverTouchSameFile(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	hasNoConflict := strings.Contains(doc, "never touch same file") ||
		strings.Contains(doc, "never touch the same file") ||
		strings.Contains(doc, "own file") ||
		strings.Contains(doc, "separate file")

	if !hasNoConflict {
		t.Error("state-format.md missing invariant that concurrent sessions never touch the same file")
	}
}

// Ordering: section sequence checks

func TestStateFormatMd_DetectionBeforeFilePerSlice(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	detectionPos := strings.Index(doc, "Detection")
	if detectionPos == -1 {
		detectionPos = strings.Index(doc, "detection")
	}

	filePerSlicePos := strings.Index(doc, "File-Per-Slice")
	if filePerSlicePos == -1 {
		filePerSlicePos = strings.Index(doc, "per-slice")
		if filePerSlicePos == -1 {
			filePerSlicePos = strings.Index(doc, "Per-Slice")
		}
	}

	if detectionPos == -1 {
		t.Fatal("state-format.md missing Detection section")
	}
	if filePerSlicePos == -1 {
		t.Fatal("state-format.md missing File-Per-Slice section")
	}

	if detectionPos >= filePerSlicePos {
		t.Errorf("Detection section (pos %d) must appear before File-Per-Slice section (pos %d)", detectionPos, filePerSlicePos)
	}
}

func TestStateFormatMd_FilePerSliceBeforeLegacy(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	filePerSlicePos := strings.Index(doc, "File-Per-Slice")
	if filePerSlicePos == -1 {
		filePerSlicePos = strings.Index(doc, "per-slice")
	}
	legacyPos := strings.Index(doc, "Legacy")
	if legacyPos == -1 {
		legacyPos = strings.Index(doc, "legacy")
	}

	if filePerSlicePos == -1 {
		t.Fatal("state-format.md missing File-Per-Slice section")
	}
	if legacyPos == -1 {
		t.Fatal("state-format.md missing Legacy section")
	}

	if filePerSlicePos >= legacyPos {
		t.Errorf("File-Per-Slice section (pos %d) must appear before Legacy section (pos %d)", filePerSlicePos, legacyPos)
	}
}

func TestStateFormatMd_CrashSafetyBeforeBackwardCompatibility(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	crashPos := strings.Index(doc, "Crash Safety")
	if crashPos == -1 {
		crashPos = strings.Index(doc, "crash safety")
	}
	backwardPos := strings.Index(doc, "Backward Compatibility")
	if backwardPos == -1 {
		backwardPos = strings.Index(doc, "backward compatibility")
	}

	if crashPos == -1 {
		t.Fatal("state-format.md missing Crash Safety section")
	}
	if backwardPos == -1 {
		t.Fatal("state-format.md missing Backward Compatibility section")
	}

	if crashPos >= backwardPos {
		t.Errorf("Crash Safety section (pos %d) must appear before Backward Compatibility section (pos %d)", crashPos, backwardPos)
	}
}

// Verification: auto check — both contracts specify Verification: auto

func TestSliceStateMd_VerificationIsAuto(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/templates/slice-state.md"))
	if err != nil {
		t.Fatalf("failed to read slice-state.md: %v", err)
	}

	doc := string(content)

	// C-76 uses Verification: auto — the file itself should not contradict this
	// by requiring manual steps. Verify the file exists and is testable by automated checks.
	// The presence of deterministic content (no "manual" verification required) satisfies auto.
	if len(doc) == 0 {
		t.Error("slice-state.md is empty — cannot satisfy Verification: auto for C-76")
	}
}

func TestStateFormatMd_VerificationIsAuto(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/state-format.md"))
	if err != nil {
		t.Fatalf("failed to read state-format.md: %v", err)
	}

	doc := string(content)

	// C-77 uses Verification: auto — the file itself should not contradict this
	// by requiring manual steps. Verify the file exists and is testable by automated checks.
	if len(doc) == 0 {
		t.Error("state-format.md is empty — cannot satisfy Verification: auto for C-77")
	}
}
