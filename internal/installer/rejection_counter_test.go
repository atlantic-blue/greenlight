package installer_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// S-25 Tests: Rejection Counter
// These tests verify that src/commands/gl/slice.md and src/references/verification-tiers.md
// contain the full rejection counter and escalation documentation as defined in
// contracts C-70 (RejectionCounter) and C-71 (RejectionEscalation).

// =============================================================================
// C-70 Tests: RejectionCounter
// Per-slice rejection tracking in-memory state
// Verifies content in slice.md (Step 6b) and verification-tiers.md (Rejection Counter)
// =============================================================================

func TestSliceMd_RejectionCounterInitializesStateAtStartOfStep6b(t *testing.T) {
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

	// The rejection state must be initialized with slice_id, rejection_count, and rejection_log
	hasSliceId := strings.Contains(sectionContent, "slice_id")
	hasRejectionCount := strings.Contains(sectionContent, "rejection_count")
	hasRejectionLog := strings.Contains(sectionContent, "rejection_log")

	if !hasSliceId {
		t.Error("slice.md Step 6b missing 'slice_id' field in rejection state initialization")
	}
	if !hasRejectionCount {
		t.Error("slice.md Step 6b missing 'rejection_count' field in rejection state initialization")
	}
	if !hasRejectionLog {
		t.Error("slice.md Step 6b missing 'rejection_log' field in rejection state initialization")
	}
}

func TestSliceMd_RejectionCounterInitializesRejectionCountToZero(t *testing.T) {
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

	// Initial rejection_count must be 0
	hasZeroInitialCount := strings.Contains(sectionContent, "rejection_count: 0") ||
		(strings.Contains(sectionContent, "rejection_count") && strings.Contains(sectionContent, ": 0"))

	if !hasZeroInitialCount {
		t.Error("slice.md Step 6b missing rejection_count initialised to 0 in rejection state")
	}
}

func TestSliceMd_RejectionCounterIncrementsOnEachNonApprovedResponse(t *testing.T) {
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

	// Counter must increment on every non-approved response
	hasIncrement := strings.Contains(sectionContent, "Increment") ||
		strings.Contains(sectionContent, "increment") ||
		strings.Contains(sectionContent, "rejection_count") && strings.Contains(sectionContent, "+1") ||
		strings.Contains(sectionContent, "rejection_count") && strings.Contains(sectionContent, "increment")

	if !hasIncrement {
		t.Error("slice.md Step 6b missing documentation that rejection_count increments on each non-approved response")
	}
}

func TestSliceMd_RejectionCounterLogEntryIncludesAttemptField(t *testing.T) {
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

	// Each rejection log entry must include an attempt field
	hasAttemptField := strings.Contains(sectionContent, "attempt:")

	if !hasAttemptField {
		t.Error("slice.md Step 6b missing 'attempt:' field in rejection log entry structure")
	}
}

func TestSliceMd_RejectionCounterLogEntryIncludesFeedbackField(t *testing.T) {
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

	// Each rejection log entry must capture the user's verbatim feedback
	hasFeedbackField := strings.Contains(sectionContent, "feedback:")

	if !hasFeedbackField {
		t.Error("slice.md Step 6b missing 'feedback:' field in rejection log entry structure")
	}
}

func TestSliceMd_RejectionCounterLogEntryIncludesClassificationField(t *testing.T) {
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

	// Each rejection log entry must include a classification field
	hasClassificationField := strings.Contains(sectionContent, "classification:")

	if !hasClassificationField {
		t.Error("slice.md Step 6b missing 'classification:' field in rejection log entry structure")
	}
}

func TestSliceMd_RejectionCounterLogEntryIncludesActionTakenField(t *testing.T) {
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

	// Each rejection log entry must include an action_taken field
	hasActionTakenField := strings.Contains(sectionContent, "action_taken:")

	if !hasActionTakenField {
		t.Error("slice.md Step 6b missing 'action_taken:' field in rejection log entry structure")
	}
}

func TestSliceMd_RejectionCounterPersistsAcrossRejectionLoopsWithinSingleExecution(t *testing.T) {
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

	// Counter must persist across rejection loops within a single /gl:slice execution
	hasPersistenceNote := strings.Contains(sectionContent, "persists") ||
		strings.Contains(sectionContent, "persist") ||
		strings.Contains(sectionContent, "Counter persists") ||
		strings.Contains(sectionContent, "counter persists")

	if !hasPersistenceNote {
		t.Error("slice.md Step 6b missing documentation that rejection counter persists across rejection loops within a single /gl:slice execution")
	}
}

func TestSliceMd_RejectionCounterResetsOnNewGlSliceExecution(t *testing.T) {
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

	// Counter must reset to 0 when /gl:slice is re-invoked (new execution)
	hasResetOnNewExecution := strings.Contains(sectionContent, "resets") ||
		strings.Contains(sectionContent, "reset") ||
		strings.Contains(sectionContent, "new execution") ||
		strings.Contains(sectionContent, "re-invoked")

	if !hasResetOnNewExecution {
		t.Error("slice.md Step 6b missing documentation that rejection counter resets to 0 on new /gl:slice execution")
	}
}

func TestSliceMd_RejectionCounterEscalationThresholdIsExactlyThree(t *testing.T) {
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

	// Escalation must trigger at exactly 3 rejections
	hasThresholdAt3 := strings.Contains(sectionContent, ">= 3") ||
		strings.Contains(sectionContent, "3 rejections") ||
		strings.Contains(sectionContent, "after 3") ||
		strings.Contains(sectionContent, "reaches 3") ||
		strings.Contains(sectionContent, "at 3")

	if !hasThresholdAt3 {
		t.Error("slice.md Step 6b missing documentation that escalation triggers at exactly 3 rejections")
	}
}

func TestSliceMd_RejectionCounterIsPerSliceNotPerContract(t *testing.T) {
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

	// Counter must be per-slice, not per-contract
	hasPerSliceNote := strings.Contains(sectionContent, "per-slice") ||
		strings.Contains(sectionContent, "per slice") ||
		strings.Contains(sectionContent, "not per-contract") ||
		strings.Contains(sectionContent, "not per contract")

	if !hasPerSliceNote {
		t.Error("slice.md Step 6b missing documentation that rejection counter is per-slice, not per-contract")
	}
}

func TestSliceMd_RejectionCounterDoesNotInteractWithCircuitBreaker(t *testing.T) {
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

	// Counter must not interact with circuit breaker attempt counters (different concerns)
	hasIsolationNote := strings.Contains(sectionContent, "circuit breaker") ||
		strings.Contains(sectionContent, "circuit-breaker") ||
		strings.Contains(sectionContent, "different concerns") ||
		strings.Contains(sectionContent, "separate") && strings.Contains(sectionContent, "counter")

	if !hasIsolationNote {
		t.Error("slice.md Step 6b missing documentation that rejection counter does not interact with circuit breaker attempt counters")
	}
}

func TestSliceMd_RejectionCounterContractRevisionDoesNotResetCounter(t *testing.T) {
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

	// A contract revision (option 2) that restarts from Step 1 must NOT reset the rejection counter
	hasNoResetOnContractRevision := strings.Contains(sectionContent, "does NOT reset") ||
		strings.Contains(sectionContent, "does not reset") ||
		strings.Contains(sectionContent, "counter persists") ||
		strings.Contains(sectionContent, "not reset the counter") ||
		strings.Contains(sectionContent, "counter is not reset")

	if !hasNoResetOnContractRevision {
		t.Error("slice.md Step 6b missing documentation that contract revision (option 2) does NOT reset the rejection counter")
	}
}

func TestSliceMd_RejectionCounterCounterOverflowHandled(t *testing.T) {
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

	// CounterOverflow error: counter exceeds 3 — trigger escalation immediately and log warning
	hasCounterOverflowHandling := strings.Contains(sectionContent, "CounterOverflow") ||
		strings.Contains(sectionContent, "counter overflow") ||
		strings.Contains(sectionContent, "exceeds 3") ||
		strings.Contains(sectionContent, "overflow")

	if !hasCounterOverflowHandling {
		t.Error("slice.md Step 6b missing CounterOverflow error handling — must trigger escalation immediately and log warning when counter exceeds 3")
	}
}

func TestSliceMd_RejectionCounterLogCorruptionHandled(t *testing.T) {
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

	// LogCorruption error: rejection log becomes inconsistent — reset log, preserve count, warn user
	hasLogCorruptionHandling := strings.Contains(sectionContent, "LogCorruption") ||
		strings.Contains(sectionContent, "log corruption") ||
		strings.Contains(sectionContent, "inconsistent") ||
		strings.Contains(sectionContent, "corrupted")

	if !hasLogCorruptionHandling {
		t.Error("slice.md Step 6b missing LogCorruption error handling — must reset log, preserve count, and warn user when log is inconsistent")
	}
}

func TestSliceMd_RejectionCounterLogPreservesChronologicalOrder(t *testing.T) {
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

	// Rejection log must preserve chronological order
	hasChronologicalOrder := strings.Contains(sectionContent, "chronological") ||
		strings.Contains(sectionContent, "order") && strings.Contains(sectionContent, "log") ||
		strings.Contains(sectionContent, "Append") ||
		strings.Contains(sectionContent, "append")

	if !hasChronologicalOrder {
		t.Error("slice.md Step 6b missing documentation that rejection log preserves chronological order")
	}
}

func TestSliceMd_RejectionCounterLogIncludesVerbatimUserFeedback(t *testing.T) {
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

	// Rejection log must include the user's verbatim feedback
	hasVerbatimFeedbackInLog := strings.Contains(sectionContent, "verbatim") ||
		strings.Contains(sectionContent, "verbatim response") ||
		strings.Contains(sectionContent, "user's verbatim") ||
		strings.Contains(sectionContent, "user verbatim")

	if !hasVerbatimFeedbackInLog {
		t.Error("slice.md Step 6b missing documentation that rejection log includes the user's verbatim feedback")
	}
}

func TestSliceMd_RejectionCounterIsInMemoryOnlyNotPersistedToDisk(t *testing.T) {
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

	// Rejection log must be in-memory only — not written to disk
	hasInMemoryNote := strings.Contains(sectionContent, "in-memory") ||
		strings.Contains(sectionContent, "in memory") ||
		strings.Contains(sectionContent, "not persisted") ||
		strings.Contains(sectionContent, "not written to disk")

	if !hasInMemoryNote {
		t.Error("slice.md Step 6b missing documentation that rejection log is in-memory only (not persisted to disk)")
	}
}

// =============================================================================
// C-71 Tests: RejectionEscalation — slice.md
// Escalation format at 3 rejections — verifies content in slice.md (Step 6b)
// =============================================================================

func TestSliceMd_EscalationHeaderIncludesSliceName(t *testing.T) {
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

	// ESCALATION header must be present with the slice_name
	hasEscalationHeader := strings.Contains(sectionContent, "ESCALATION") ||
		strings.Contains(sectionContent, "Escalation")

	if !hasEscalationHeader {
		t.Error("slice.md Step 6b missing ESCALATION header format for the escalation prompt")
	}
}

func TestSliceMd_EscalationMessageStatesRejectedThreeTimes(t *testing.T) {
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

	// Escalation message must tell the user the slice has been rejected 3 times
	hasRejected3Times := strings.Contains(sectionContent, "rejected 3 times") ||
		strings.Contains(sectionContent, "3 times") ||
		strings.Contains(sectionContent, "3 rejections")

	if !hasRejected3Times {
		t.Error("slice.md Step 6b escalation missing 'rejected 3 times' message")
	}
}

func TestSliceMd_EscalationDisplaysRejectionHistoryWithFeedbackAndActions(t *testing.T) {
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

	// Escalation must display rejection history with both feedback and action taken for each entry
	hasRejectionHistory := strings.Contains(sectionContent, "Rejection history") ||
		strings.Contains(sectionContent, "rejection history") ||
		strings.Contains(sectionContent, "history:")

	if !hasRejectionHistory {
		t.Error("slice.md Step 6b escalation missing rejection history display showing feedback and action taken for each rejection")
	}
}

func TestSliceMd_EscalationPresentsTresEscalationOptions(t *testing.T) {
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

	// Escalation must present exactly three options: re-scope, pair, skip verification
	hasRescope := strings.Contains(sectionContent, "Re-scope") ||
		strings.Contains(sectionContent, "re-scope") ||
		strings.Contains(sectionContent, "rescope")
	hasPair := strings.Contains(sectionContent, "Pair") ||
		strings.Contains(sectionContent, "pair")
	hasSkip := strings.Contains(sectionContent, "Skip verification") ||
		strings.Contains(sectionContent, "skip verification") ||
		strings.Contains(sectionContent, "Skip") && strings.Contains(sectionContent, "verification")

	if !hasRescope {
		t.Error("slice.md Step 6b escalation missing 're-scope' option (option 1)")
	}
	if !hasPair {
		t.Error("slice.md Step 6b escalation missing 'pair' option (option 2)")
	}
	if !hasSkip {
		t.Error("slice.md Step 6b escalation missing 'skip verification' option (option 3)")
	}
}

func TestSliceMd_EscalationOption1ReScopeResetsCounterAndRestartSlice(t *testing.T) {
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

	// Option 1 (re-scope): must reset rejection counter and restart slice from scratch
	hasRescopeResetAndRestart := (strings.Contains(sectionContent, "re-scope") || strings.Contains(sectionContent, "Re-scope")) &&
		(strings.Contains(sectionContent, "Reset") || strings.Contains(sectionContent, "reset")) &&
		(strings.Contains(sectionContent, "restart") || strings.Contains(sectionContent, "Restart") || strings.Contains(sectionContent, "scratch"))

	if !hasRescopeResetAndRestart {
		t.Error("slice.md Step 6b escalation option 1 (re-scope) missing documentation that it resets the rejection counter and restarts the slice from scratch")
	}
}

func TestSliceMd_EscalationOption2PairCollectsGuidanceAndSpawnsTestWriter(t *testing.T) {
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

	// Option 2 (pair): must collect step-by-step guidance from user and spawn test writer
	hasPairGuidanceAndTestWriter := (strings.Contains(sectionContent, "Pair") || strings.Contains(sectionContent, "pair")) &&
		(strings.Contains(sectionContent, "guidance") || strings.Contains(sectionContent, "step-by-step")) &&
		(strings.Contains(sectionContent, "test writer") || strings.Contains(sectionContent, "gl-test-writer"))

	if !hasPairGuidanceAndTestWriter {
		t.Error("slice.md Step 6b escalation option 2 (pair) missing documentation that it collects detailed guidance and spawns the test writer")
	}
}

func TestSliceMd_EscalationOption2PairResetsCounter(t *testing.T) {
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

	// Option 2 (pair) must also reset the rejection counter
	hasPairResetsCounter := (strings.Contains(sectionContent, "Pair") || strings.Contains(sectionContent, "pair")) &&
		(strings.Contains(sectionContent, "Reset") || strings.Contains(sectionContent, "reset"))

	if !hasPairResetsCounter {
		t.Error("slice.md Step 6b escalation option 2 (pair) missing documentation that it resets the rejection counter")
	}
}

func TestSliceMd_EscalationOption3SkipMarksSliceAsAutoVerified(t *testing.T) {
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

	// Option 3 (skip): must mark the slice as auto-verified and proceed
	hasSkipMarksAuto := (strings.Contains(sectionContent, "Skip") || strings.Contains(sectionContent, "skip")) &&
		(strings.Contains(sectionContent, "auto") || strings.Contains(sectionContent, "auto-verified"))

	if !hasSkipMarksAuto {
		t.Error("slice.md Step 6b escalation option 3 (skip) missing documentation that it marks the slice effective tier as 'auto' and proceeds")
	}
}

func TestSliceMd_EscalationOption3SkipCreatesExplicitLogEntry(t *testing.T) {
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

	// Option 3 (skip): must create an explicit log entry for auditability
	hasSkipLogEntry := strings.Contains(sectionContent, "Verification skipped") ||
		strings.Contains(sectionContent, "verification skipped") ||
		strings.Contains(sectionContent, "Mismatch acknowledged") ||
		strings.Contains(sectionContent, "mismatch acknowledged") ||
		strings.Contains(sectionContent, "acknowledged") && strings.Contains(sectionContent, "deferred")

	if !hasSkipLogEntry {
		t.Error("slice.md Step 6b escalation option 3 (skip) missing documentation for explicit log entry acknowledging verification was skipped")
	}
}

func TestSliceMd_EscalationOption3SkipDoesNotResetCounter(t *testing.T) {
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

	// Skip option exits the loop — it does NOT reset the counter
	// The documentation must distinguish: re-scope and pair reset, skip does not
	hasRescopeResetsNote := (strings.Contains(sectionContent, "re-scope") || strings.Contains(sectionContent, "Re-scope") ||
		strings.Contains(sectionContent, "rescope")) && strings.Contains(sectionContent, "reset")

	hasSkipExitsNote := (strings.Contains(sectionContent, "Skip") || strings.Contains(sectionContent, "skip")) &&
		(strings.Contains(sectionContent, "proceed") || strings.Contains(sectionContent, "exits") || strings.Contains(sectionContent, "Step 7"))

	if !hasRescopeResetsNote {
		t.Error("slice.md Step 6b escalation missing documentation that re-scope (and pair) reset the counter but skip does not")
	}
	if !hasSkipExitsNote {
		t.Error("slice.md Step 6b escalation option 3 (skip) missing documentation that it exits the loop and proceeds to Step 7 without resetting the counter")
	}
}

func TestSliceMd_EscalationInvalidChoiceReprompts(t *testing.T) {
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

	// InvalidEscalationChoice: re-prompt when user enters something other than 1/2/3
	hasEscalationReprompt := strings.Contains(sectionContent, "InvalidEscalationChoice") ||
		strings.Contains(sectionContent, "Please choose 1, 2, or 3") ||
		strings.Contains(sectionContent, "1/2/3") && strings.Contains(sectionContent, "re-prompt")

	if !hasEscalationReprompt {
		t.Error("slice.md Step 6b escalation missing InvalidEscalationChoice error handling — must re-prompt 'Please choose 1, 2, or 3' on invalid input")
	}
}

func TestSliceMd_EscalationHandlesEmptyRejectionLog(t *testing.T) {
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

	// EmptyRejectionLog: escalation triggered but rejection log is empty — display without history, warn in log
	hasEmptyLogHandling := strings.Contains(sectionContent, "EmptyRejectionLog") ||
		strings.Contains(sectionContent, "empty rejection log") ||
		strings.Contains(sectionContent, "empty log") ||
		strings.Contains(sectionContent, "without history")

	if !hasEmptyLogHandling {
		t.Error("slice.md Step 6b escalation missing EmptyRejectionLog error handling — must display escalation without history and warn when log is empty")
	}
}

func TestSliceMd_EscalationFollowsExistingCheckpointPatterns(t *testing.T) {
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

	// Escalation format must follow the existing checkpoint patterns (numbered options with question)
	hasCheckpointPattern := (strings.Contains(sectionContent, "Which option?") ||
		strings.Contains(sectionContent, "which option") ||
		strings.Contains(sectionContent, "(1/2/3)")) &&
		(strings.Contains(sectionContent, "1)") || strings.Contains(sectionContent, "1."))

	if !hasCheckpointPattern {
		t.Error("slice.md Step 6b escalation must follow existing checkpoint format patterns with numbered options and a question prompt (e.g. 'Which option? (1/2/3)')")
	}
}

func TestSliceMd_EscalationThresholdMatchesCircuitBreakerThreshold(t *testing.T) {
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

	// Escalation threshold of 3 must match circuit breaker per-test threshold (system-wide consistency)
	hasThreshold3 := strings.Contains(sectionContent, "3") &&
		(strings.Contains(sectionContent, "circuit breaker") || strings.Contains(sectionContent, "consistent") ||
			strings.Contains(sectionContent, "threshold"))

	if !hasThreshold3 {
		t.Error("slice.md Step 6b escalation missing documentation that escalation threshold of 3 matches the circuit breaker per-test threshold for system-wide consistency")
	}
}

// =============================================================================
// C-70 Tests: RejectionCounter — verification-tiers.md
// Verifies the Rejection Counter section in verification-tiers.md is expanded
// with the YAML state format, log structure, and counter invariants
// =============================================================================

func TestVerificationTiersMd_RejectionCounterContainsYAMLStateFormat(t *testing.T) {
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

	// Section must contain the YAML state format with all three fields
	hasSliceId := strings.Contains(sectionContent, "slice_id")
	hasRejectionCount := strings.Contains(sectionContent, "rejection_count")
	hasRejectionLog := strings.Contains(sectionContent, "rejection_log")

	if !hasSliceId {
		t.Error("verification-tiers.md Rejection Counter section missing 'slice_id' in YAML state format")
	}
	if !hasRejectionCount {
		t.Error("verification-tiers.md Rejection Counter section missing 'rejection_count' in YAML state format")
	}
	if !hasRejectionLog {
		t.Error("verification-tiers.md Rejection Counter section missing 'rejection_log' in YAML state format")
	}
}

func TestVerificationTiersMd_RejectionCounterContainsLogEntryStructure(t *testing.T) {
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

	// Section must contain the log entry structure with attempt, feedback, classification, action_taken
	hasAttempt := strings.Contains(sectionContent, "attempt:")
	hasFeedback := strings.Contains(sectionContent, "feedback:")
	hasClassification := strings.Contains(sectionContent, "classification:")
	hasActionTaken := strings.Contains(sectionContent, "action_taken:")

	if !hasAttempt {
		t.Error("verification-tiers.md Rejection Counter section missing 'attempt:' field in log entry structure")
	}
	if !hasFeedback {
		t.Error("verification-tiers.md Rejection Counter section missing 'feedback:' field in log entry structure")
	}
	if !hasClassification {
		t.Error("verification-tiers.md Rejection Counter section missing 'classification:' field in log entry structure")
	}
	if !hasActionTaken {
		t.Error("verification-tiers.md Rejection Counter section missing 'action_taken:' field in log entry structure")
	}
}

func TestVerificationTiersMd_RejectionCounterDocumentsEscalationTriggerAtExactlyThree(t *testing.T) {
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

	// Escalation must trigger at exactly 3 — not before, not after
	hasExactThreshold := strings.Contains(sectionContent, ">= 3") ||
		strings.Contains(sectionContent, "exactly 3") ||
		strings.Contains(sectionContent, "reaches 3") ||
		strings.Contains(sectionContent, "at 3") ||
		strings.Contains(sectionContent, "after 3")

	if !hasExactThreshold {
		t.Error("verification-tiers.md Rejection Counter section missing documentation that escalation triggers at exactly 3 (not before, not after)")
	}
}

func TestVerificationTiersMd_RejectionCounterDocumentsCounterPersistedAcrossLoops(t *testing.T) {
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

	// Counter persists across rejection loops within a single /gl:slice execution
	hasPersistenceAcrossLoops := strings.Contains(sectionContent, "persists") ||
		strings.Contains(sectionContent, "persist") ||
		strings.Contains(sectionContent, "within a single")

	if !hasPersistenceAcrossLoops {
		t.Error("verification-tiers.md Rejection Counter section missing documentation that counter persists across rejection loops within a single execution")
	}
}

func TestVerificationTiersMd_RejectionCounterDocumentsCounterResetsOnNewExecution(t *testing.T) {
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

	// Counter resets to 0 on new /gl:slice execution (not persisted to disk)
	hasResetOnNewExecution := strings.Contains(sectionContent, "resets") ||
		strings.Contains(sectionContent, "reset") ||
		strings.Contains(sectionContent, "new execution") ||
		strings.Contains(sectionContent, "re-invoked")

	if !hasResetOnNewExecution {
		t.Error("verification-tiers.md Rejection Counter section missing documentation that counter resets to 0 on new /gl:slice execution")
	}
}

func TestVerificationTiersMd_RejectionCounterDocumentsCounterOverflowHandling(t *testing.T) {
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

	// CounterOverflow error handling must be documented
	hasCounterOverflowHandling := strings.Contains(sectionContent, "CounterOverflow") ||
		strings.Contains(sectionContent, "counter overflow") ||
		strings.Contains(sectionContent, "exceeds 3") ||
		strings.Contains(sectionContent, "overflow")

	if !hasCounterOverflowHandling {
		t.Error("verification-tiers.md Rejection Counter section missing CounterOverflow error handling documentation")
	}
}

func TestVerificationTiersMd_RejectionCounterDocumentsLogCorruptionHandling(t *testing.T) {
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

	// LogCorruption error handling must be documented
	hasLogCorruptionHandling := strings.Contains(sectionContent, "LogCorruption") ||
		strings.Contains(sectionContent, "log corruption") ||
		strings.Contains(sectionContent, "inconsistent") ||
		strings.Contains(sectionContent, "corrupted")

	if !hasLogCorruptionHandling {
		t.Error("verification-tiers.md Rejection Counter section missing LogCorruption error handling documentation")
	}
}

func TestVerificationTiersMd_RejectionCounterDocumentsInMemoryOnlyConstraint(t *testing.T) {
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

	// Rejection log is in-memory only — not written to disk
	hasInMemoryConstraint := strings.Contains(sectionContent, "in-memory") ||
		strings.Contains(sectionContent, "in memory") ||
		strings.Contains(sectionContent, "not written to disk") ||
		strings.Contains(sectionContent, "not persisted")

	if !hasInMemoryConstraint {
		t.Error("verification-tiers.md Rejection Counter section missing documentation that rejection log is in-memory only (not written to disk)")
	}
}

func TestVerificationTiersMd_RejectionCounterDocumentsContractRevisionDoesNotResetCounter(t *testing.T) {
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

	// Contract revision restarting from Step 1 must NOT reset the rejection counter
	hasNoResetInvariant := strings.Contains(sectionContent, "does NOT reset") ||
		strings.Contains(sectionContent, "does not reset") ||
		strings.Contains(sectionContent, "counter persists") ||
		strings.Contains(sectionContent, "not reset the counter") ||
		strings.Contains(sectionContent, "counter is not reset")

	if !hasNoResetInvariant {
		t.Error("verification-tiers.md Rejection Counter section missing invariant that contract revision (restarting from Step 1) does NOT reset the rejection counter")
	}
}

// =============================================================================
// C-71 Tests: RejectionEscalation — verification-tiers.md
// Verifies escalation details in verification-tiers.md (Rejection Counter section)
// =============================================================================

func TestVerificationTiersMd_RejectionCounterContainsEscalationFormat(t *testing.T) {
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

	// Escalation format must appear in the Rejection Counter section
	hasEscalationFormat := strings.Contains(sectionContent, "ESCALATION") ||
		strings.Contains(sectionContent, "Escalation")

	if !hasEscalationFormat {
		t.Error("verification-tiers.md Rejection Counter section missing ESCALATION format documentation")
	}
}

func TestVerificationTiersMd_RejectionCounterEscalationIncludesThreeOptions(t *testing.T) {
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

	// Escalation must include three options: re-scope, pair, skip
	hasRescope := strings.Contains(sectionContent, "Re-scope") ||
		strings.Contains(sectionContent, "re-scope") ||
		strings.Contains(sectionContent, "rescope")
	hasPair := strings.Contains(sectionContent, "Pair") ||
		strings.Contains(sectionContent, "pair")
	hasSkip := strings.Contains(sectionContent, "Skip") ||
		strings.Contains(sectionContent, "skip")

	if !hasRescope {
		t.Error("verification-tiers.md Rejection Counter section missing 're-scope' escalation option")
	}
	if !hasPair {
		t.Error("verification-tiers.md Rejection Counter section missing 'pair' escalation option")
	}
	if !hasSkip {
		t.Error("verification-tiers.md Rejection Counter section missing 'skip' escalation option")
	}
}

func TestVerificationTiersMd_RejectionCounterEscalationRoutingForAllThreeOptions(t *testing.T) {
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

	// Each escalation option must document its routing action
	hasRescopeRouting := (strings.Contains(sectionContent, "re-scope") || strings.Contains(sectionContent, "Re-scope")) &&
		(strings.Contains(sectionContent, "restart") || strings.Contains(sectionContent, "scratch"))
	hasPairRouting := (strings.Contains(sectionContent, "pair") || strings.Contains(sectionContent, "Pair")) &&
		(strings.Contains(sectionContent, "guidance") || strings.Contains(sectionContent, "test writer"))
	hasSkipRouting := (strings.Contains(sectionContent, "skip") || strings.Contains(sectionContent, "Skip")) &&
		(strings.Contains(sectionContent, "auto") || strings.Contains(sectionContent, "proceed"))

	if !hasRescopeRouting {
		t.Error("verification-tiers.md Rejection Counter section missing routing action for re-scope option (restart from scratch)")
	}
	if !hasPairRouting {
		t.Error("verification-tiers.md Rejection Counter section missing routing action for pair option (collect guidance and spawn test writer)")
	}
	if !hasSkipRouting {
		t.Error("verification-tiers.md Rejection Counter section missing routing action for skip option (mark as auto, proceed)")
	}
}

func TestVerificationTiersMd_RejectionCounterEscalationInvalidChoiceHandled(t *testing.T) {
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

	// InvalidEscalationChoice error handling must be documented
	hasInvalidChoiceHandling := strings.Contains(sectionContent, "InvalidEscalationChoice") ||
		strings.Contains(sectionContent, "Please choose 1, 2, or 3") ||
		strings.Contains(sectionContent, "invalid choice") ||
		strings.Contains(sectionContent, "re-prompt")

	if !hasInvalidChoiceHandling {
		t.Error("verification-tiers.md Rejection Counter section missing InvalidEscalationChoice error handling (re-prompt on invalid input)")
	}
}

func TestVerificationTiersMd_RejectionCounterEscalationEmptyLogHandled(t *testing.T) {
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

	// EmptyRejectionLog error handling must be documented
	hasEmptyLogHandling := strings.Contains(sectionContent, "EmptyRejectionLog") ||
		strings.Contains(sectionContent, "empty rejection log") ||
		strings.Contains(sectionContent, "empty log") ||
		strings.Contains(sectionContent, "without history")

	if !hasEmptyLogHandling {
		t.Error("verification-tiers.md Rejection Counter section missing EmptyRejectionLog error handling documentation")
	}
}

func TestVerificationTiersMd_RejectionCounterEscalationSkipCreatesExplicitLogEntry(t *testing.T) {
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

	// Skip option must create an explicit log entry for auditability
	hasSkipLogEntry := strings.Contains(sectionContent, "Verification skipped") ||
		strings.Contains(sectionContent, "verification skipped") ||
		strings.Contains(sectionContent, "Mismatch acknowledged") ||
		strings.Contains(sectionContent, "mismatch acknowledged") ||
		strings.Contains(sectionContent, "acknowledged")

	if !hasSkipLogEntry {
		t.Error("verification-tiers.md Rejection Counter section missing documentation that skip option creates an explicit log entry (auditability)")
	}
}

func TestVerificationTiersMd_RejectionCounterEscalationReScopeAndPairResetCounterSkipDoesNot(t *testing.T) {
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

	// Re-scope and pair reset the counter; skip does not
	hasResetForRescopeAndPair := strings.Contains(sectionContent, "Reset") ||
		strings.Contains(sectionContent, "reset")

	if !hasResetForRescopeAndPair {
		t.Error("verification-tiers.md Rejection Counter section missing documentation that re-scope and pair reset the rejection counter (skip does not)")
	}
}

// =============================================================================
// Cross-cutting invariant tests
// These verify invariants that span both contracts (C-70 and C-71) across files
// =============================================================================

func TestSliceMd_RejectionCounterTrackingAppearsWithinStep6bSection(t *testing.T) {
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

	// Counter tracking content must be within Step 6b (before Step 7)
	step6bSection := doc[step6bPos:step7Pos]

	hasCounterInStep6b := strings.Contains(step6bSection, "rejection_count") ||
		strings.Contains(step6bSection, "rejection counter") ||
		strings.Contains(step6bSection, "Rejection counter")

	if !hasCounterInStep6b {
		t.Error("slice.md rejection counter tracking content must appear within Step 6b section (before Step 7)")
	}
}

func TestSliceMd_EscalationContentAppearsWithinStep6bSection(t *testing.T) {
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

	// Escalation content must be within Step 6b (before Step 7)
	step6bSection := doc[step6bPos:step7Pos]

	hasEscalationInStep6b := strings.Contains(step6bSection, "ESCALATION") ||
		strings.Contains(step6bSection, "Escalation")

	if !hasEscalationInStep6b {
		t.Error("slice.md escalation content must appear within Step 6b section (before Step 7)")
	}
}

func TestVerificationTiersMd_RejectionCounterSectionContainsStateStructureAndEscalationDetails(t *testing.T) {
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

	// Section must contain both state structure (YAML) AND escalation details
	hasStateStructure := strings.Contains(sectionContent, "slice_id") &&
		strings.Contains(sectionContent, "rejection_count") &&
		strings.Contains(sectionContent, "rejection_log")
	hasEscalationDetails := strings.Contains(sectionContent, "ESCALATION") ||
		strings.Contains(sectionContent, "Escalation")

	if !hasStateStructure {
		t.Error("verification-tiers.md Rejection Counter section missing YAML state structure (slice_id, rejection_count, rejection_log)")
	}
	if !hasEscalationDetails {
		t.Error("verification-tiers.md Rejection Counter section missing escalation details")
	}
}

func TestBothFiles_EscalationContentPresentInBothSliceMdAndVerificationTiersMd(t *testing.T) {
	sliceContent, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	tiersContent, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-tiers.md"))
	if err != nil {
		t.Fatalf("failed to read verification-tiers.md: %v", err)
	}

	sliceDoc := string(sliceContent)
	tiersDoc := string(tiersContent)

	// Escalation content must appear in both files
	sliceHasEscalation := strings.Contains(sliceDoc, "ESCALATION") ||
		strings.Contains(sliceDoc, "Escalation")
	tiersHasEscalation := strings.Contains(tiersDoc, "ESCALATION") ||
		strings.Contains(tiersDoc, "Escalation")

	if !sliceHasEscalation {
		t.Error("slice.md missing escalation content — escalation must be documented in both slice.md and verification-tiers.md")
	}
	if !tiersHasEscalation {
		t.Error("verification-tiers.md missing escalation content — escalation must be documented in both slice.md and verification-tiers.md")
	}
}

func TestBothFiles_RejectionCounterStateFormatPresentInBothFiles(t *testing.T) {
	sliceContent, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	tiersContent, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/verification-tiers.md"))
	if err != nil {
		t.Fatalf("failed to read verification-tiers.md: %v", err)
	}

	sliceDoc := string(sliceContent)
	tiersDoc := string(tiersContent)

	// Rejection counter YAML state format must appear in both files
	sliceHasStateFormat := strings.Contains(sliceDoc, "rejection_count") &&
		strings.Contains(sliceDoc, "rejection_log")
	tiersHasStateFormat := strings.Contains(tiersDoc, "rejection_count") &&
		strings.Contains(tiersDoc, "rejection_log")

	if !sliceHasStateFormat {
		t.Error("slice.md missing rejection counter state format (rejection_count, rejection_log)")
	}
	if !tiersHasStateFormat {
		t.Error("verification-tiers.md missing rejection counter state format (rejection_count, rejection_log)")
	}
}
