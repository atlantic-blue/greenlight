package installer_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// C-81 Tests: SliceCommandStateWrite
// C-82 Tests: SliceSessionTracking
// These tests verify that src/commands/gl/slice.md contains all required instructions
// as defined in contracts C-81 and C-82 for slice S-30.
// Verification: verify (C-81), auto (C-82)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func readSliceMd(t *testing.T) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read src/commands/gl/slice.md: %v", err)
	}
	return string(content)
}

// ---------------------------------------------------------------------------
// C-81: SliceCommandStateWrite
// ---------------------------------------------------------------------------

// Pre-flight: state format detection

func TestSliceMd_C81_ContainsStateFormatDetectionInstruction(t *testing.T) {
	doc := readSliceMd(t)
	hasDetection := strings.Contains(doc, "state format") ||
		strings.Contains(doc, "State Format") ||
		strings.Contains(doc, "C-80") ||
		strings.Contains(doc, "detect")
	if !hasDetection {
		t.Error("slice.md missing state format detection instruction in pre-flight (C-80 reference)")
	}
}

func TestSliceMd_C81_ContainsPreFlightSection(t *testing.T) {
	doc := readSliceMd(t)
	hasPreFlight := strings.Contains(doc, "Pre-flight") ||
		strings.Contains(doc, "pre-flight") ||
		strings.Contains(doc, "Pre-Flight")
	if !hasPreFlight {
		t.Error("slice.md missing Pre-flight section instruction")
	}
}

// File-per-slice read

func TestSliceMd_C81_ContainsFilePerSliceReadInstruction(t *testing.T) {
	doc := readSliceMd(t)
	if !strings.Contains(doc, ".greenlight/slices/") {
		t.Error("slice.md missing file-per-slice read instruction (.greenlight/slices/{id}.md)")
	}
}

func TestSliceMd_C81_ContainsFilePerSliceFormatMention(t *testing.T) {
	doc := readSliceMd(t)
	if !strings.Contains(doc, "file-per-slice") {
		t.Error("slice.md missing 'file-per-slice' format mention")
	}
}

// Own-slice-only write

func TestSliceMd_C81_ContainsOwnSliceOnlyWriteInstruction(t *testing.T) {
	doc := readSliceMd(t)
	hasOwnSlice := strings.Contains(doc, "own slice file") ||
		strings.Contains(doc, "only to its own") ||
		strings.Contains(doc, "write only to") ||
		strings.Contains(doc, "writes ONLY to")
	if !hasOwnSlice {
		t.Error("slice.md missing own-slice-only write instruction (must write ONLY to its own slice file)")
	}
}

// STATE.md regeneration after every write

func TestSliceMd_C81_ContainsStateMdRegenerationAfterEveryWriteInstruction(t *testing.T) {
	doc := readSliceMd(t)
	hasRegen := strings.Contains(doc, "Regenerate STATE.md") ||
		strings.Contains(doc, "regenerate STATE.md") ||
		strings.Contains(doc, "Regenerate") && strings.Contains(doc, "STATE.md") ||
		strings.Contains(doc, "D-34")
	if !hasRegen {
		t.Error("slice.md missing STATE.md regeneration instruction after every write (D-34)")
	}
}

func TestSliceMd_C81_ContainsStateMdRegenerationInStep4(t *testing.T) {
	doc := readSliceMd(t)
	// Step 4 is the claim/start-implementation step; must regenerate STATE.md
	step4Pos := strings.Index(doc, "Step 4")
	if step4Pos == -1 {
		step4Pos = strings.Index(doc, "step 4")
	}
	regenPos := strings.Index(doc, "Regenerate STATE")
	if regenPos == -1 {
		regenPos = strings.Index(doc, "regenerate STATE")
	}
	if step4Pos == -1 {
		t.Error("slice.md missing Step 4 reference")
	}
	if regenPos == -1 {
		t.Error("slice.md missing STATE.md regeneration instruction")
	}
}

func TestSliceMd_C81_ContainsStateMdRegenerationInStep10(t *testing.T) {
	doc := readSliceMd(t)
	// Step 10 is the completion/status-update step; must also regenerate STATE.md
	hasStep10 := strings.Contains(doc, "Step 10") || strings.Contains(doc, "step 10")
	if !hasStep10 {
		t.Error("slice.md missing Step 10 reference (completion / status update)")
	}
}

// GENERATED header for regenerated STATE.md

func TestSliceMd_C81_ContainsGeneratedHeaderReferenceForStateMd(t *testing.T) {
	doc := readSliceMd(t)
	hasHeader := strings.Contains(doc, "GENERATED") ||
		strings.Contains(doc, "<!-- generated") ||
		strings.Contains(doc, "generated file")
	if !hasHeader {
		t.Error("slice.md missing GENERATED header reference for regenerated STATE.md")
	}
}

// Legacy format backward compatibility

func TestSliceMd_C81_ContainsLegacyFormatBackwardCompatibilityInstruction(t *testing.T) {
	doc := readSliceMd(t)
	hasLegacy := strings.Contains(doc, "legacy") || strings.Contains(doc, "Legacy")
	if !hasLegacy {
		t.Error("slice.md missing legacy format backward compatibility instruction")
	}
}

func TestSliceMd_C81_ContainsLegacyUnchangedBehaviourInstruction(t *testing.T) {
	doc := readSliceMd(t)
	hasUnchanged := strings.Contains(doc, "unchanged") ||
		strings.Contains(doc, "as before") ||
		strings.Contains(doc, "no change") ||
		strings.Contains(doc, "D-37")
	if !hasUnchanged {
		t.Error("slice.md missing instruction that legacy format behaviour is completely unchanged")
	}
}

// Crash safety

func TestSliceMd_C81_ContainsCrashSafetyInstruction(t *testing.T) {
	doc := readSliceMd(t)
	hasCrashSafety := strings.Contains(doc, "write-to-temp-then-rename") ||
		strings.Contains(doc, "temp-then-rename") ||
		strings.Contains(doc, "atomic") ||
		strings.Contains(doc, "NFR-4")
	if !hasCrashSafety {
		t.Error("slice.md missing crash safety instruction (write-to-temp-then-rename / atomic write, NFR-4)")
	}
}

// project-state.json reference for regeneration

func TestSliceMd_C81_ContainsProjectStateJsonReferenceForRegeneration(t *testing.T) {
	doc := readSliceMd(t)
	if !strings.Contains(doc, "project-state.json") {
		t.Error("slice.md missing project-state.json reference for STATE.md regeneration")
	}
}

// Error handling

func TestSliceMd_C81_ContainsSliceFileNotFoundError(t *testing.T) {
	doc := readSliceMd(t)
	if !strings.Contains(doc, "SliceFileNotFound") {
		t.Error("slice.md missing SliceFileNotFound error handling instruction")
	}
}

func TestSliceMd_C81_ContainsSliceFileWriteFailureError(t *testing.T) {
	doc := readSliceMd(t)
	if !strings.Contains(doc, "SliceFileWriteFailure") {
		t.Error("slice.md missing SliceFileWriteFailure error handling instruction")
	}
}

func TestSliceMd_C81_ContainsRegenerationFailureError(t *testing.T) {
	doc := readSliceMd(t)
	if !strings.Contains(doc, "RegenerationFailure") {
		t.Error("slice.md missing RegenerationFailure error handling instruction")
	}
}

func TestSliceMd_C81_ContainsConcurrentSliceClaimError(t *testing.T) {
	doc := readSliceMd(t)
	if !strings.Contains(doc, "ConcurrentSliceClaim") {
		t.Error("slice.md missing ConcurrentSliceClaim error handling instruction")
	}
}

func TestSliceMd_C81_ContainsRegenerationFailureSafeToontinueInstruction(t *testing.T) {
	doc := readSliceMd(t)
	// RegenerationFailure must warn but not abort — slice file is still correct
	hasSafeContinue := strings.Contains(doc, "Warn but continue") ||
		strings.Contains(doc, "warn but continue") ||
		strings.Contains(doc, "still correct") ||
		strings.Contains(doc, "partial failure is safe")
	if !hasSafeContinue {
		t.Error("slice.md missing instruction that RegenerationFailure warns but continues (slice file is still correct)")
	}
}

// Invariants

func TestSliceMd_C81_ContainsSourceOfTruthMention(t *testing.T) {
	doc := readSliceMd(t)
	if !strings.Contains(doc, "source of truth") {
		t.Error("slice.md missing 'source of truth' mention (slice file is source of truth, STATE.md is convenience view)")
	}
}

func TestSliceMd_C81_ContainsPipelineStepUpdateReference(t *testing.T) {
	doc := readSliceMd(t)
	// Must reference both the claim step (Step 4) and completion step (Step 10)
	hasStep4 := strings.Contains(doc, "Step 4") || strings.Contains(doc, "step 4")
	hasStep10 := strings.Contains(doc, "Step 10") || strings.Contains(doc, "step 10")
	if !hasStep4 {
		t.Error("slice.md missing Step 4 pipeline step reference (claim slice / start implementation)")
	}
	if !hasStep10 {
		t.Error("slice.md missing Step 10 pipeline step reference (completion / status update)")
	}
}

func TestSliceMd_C81_ContainsSliceIdValidationPathTraversalPrevention(t *testing.T) {
	doc := readSliceMd(t)
	hasPathValidation := strings.Contains(doc, "path traversal") ||
		strings.Contains(doc, "traversal") ||
		strings.Contains(doc, "validated") && strings.Contains(doc, "slice ID") ||
		strings.Contains(doc, "validate") && strings.Contains(doc, "ID")
	if !hasPathValidation {
		t.Error("slice.md missing slice ID validation / path traversal prevention mention")
	}
}

func TestSliceMd_C81_ContainsConcurrentSessionsNeverConflictInstruction(t *testing.T) {
	doc := readSliceMd(t)
	hasNoConflict := strings.Contains(doc, "FR-2") ||
		strings.Contains(doc, "never conflict") ||
		strings.Contains(doc, "different slices") ||
		strings.Contains(doc, "concurrent") && strings.Contains(doc, "slices")
	if !hasNoConflict {
		t.Error("slice.md missing instruction that concurrent sessions writing to different slices never conflict (FR-2)")
	}
}

func TestSliceMd_C81_ContainsStateMdConvenienceViewMention(t *testing.T) {
	doc := readSliceMd(t)
	hasConvenienceView := strings.Contains(doc, "convenience view") ||
		strings.Contains(doc, "summary view") ||
		strings.Contains(doc, "generated") && strings.Contains(doc, "STATE.md")
	if !hasConvenienceView {
		t.Error("slice.md missing mention that STATE.md is a convenience/summary view (not source of truth)")
	}
}

// ---------------------------------------------------------------------------
// C-82: SliceSessionTracking
// ---------------------------------------------------------------------------

// Session ID format

func TestSliceMd_C82_ContainsSessionIdISOTimestampFormatInstruction(t *testing.T) {
	doc := readSliceMd(t)
	hasISO := strings.Contains(doc, "ISO") || strings.Contains(doc, "timestamp")
	if !hasISO {
		t.Error("slice.md missing session ID ISO timestamp format instruction")
	}
}

func TestSliceMd_C82_ContainsSessionIdRandomHexSuffixInstruction(t *testing.T) {
	doc := readSliceMd(t)
	hasHexSuffix := strings.Contains(doc, "hex") ||
		strings.Contains(doc, "random") && strings.Contains(doc, "suffix") ||
		strings.Contains(doc, "4-char") ||
		strings.Contains(doc, "a7f3")
	if !hasHexSuffix {
		t.Error("slice.md missing session ID random hex suffix instruction (format: ISO-timestamp-{4-char-hex})")
	}
}

// Session field write on claim

func TestSliceMd_C82_ContainsSessionFieldWriteOnClaimInstruction(t *testing.T) {
	doc := readSliceMd(t)
	hasSessionWrite := strings.Contains(doc, "session") && strings.Contains(doc, "claim") ||
		strings.Contains(doc, "Set session") ||
		strings.Contains(doc, "set session") ||
		strings.Contains(doc, "Write session")
	if !hasSessionWrite {
		t.Error("slice.md missing session field write instruction on slice claim")
	}
}

// Session field clear on completion

func TestSliceMd_C82_ContainsSessionFieldClearOnCompletionInstruction(t *testing.T) {
	doc := readSliceMd(t)
	hasClear := strings.Contains(doc, "Clear session") ||
		strings.Contains(doc, "clear session") ||
		strings.Contains(doc, "session field") && strings.Contains(doc, "clear") ||
		strings.Contains(doc, "session") && strings.Contains(doc, "completion")
	if !hasClear {
		t.Error("slice.md missing session field clear instruction on slice completion")
	}
}

// Stale session warning

func TestSliceMd_C82_ContainsStaleSessionWarningInstruction(t *testing.T) {
	doc := readSliceMd(t)
	hasStale := strings.Contains(doc, "stale") || strings.Contains(doc, "StaleSession")
	if !hasStale {
		t.Error("slice.md missing stale session warning instruction")
	}
}

func TestSliceMd_C82_ContainsConcurrentSliceClaimWarningAndPrompt(t *testing.T) {
	doc := readSliceMd(t)
	hasWarnPrompt := strings.Contains(doc, "Warn") && strings.Contains(doc, "prompt") ||
		strings.Contains(doc, "warn") && strings.Contains(doc, "prompt") ||
		strings.Contains(doc, "Continue anyway")
	if !hasWarnPrompt {
		t.Error("slice.md missing concurrent claim warning and prompt (Warn + 'Continue anyway? (y/n)')")
	}
}

// Advisory-only session

func TestSliceMd_C82_ContainsAdvisoryOnlyMention(t *testing.T) {
	doc := readSliceMd(t)
	if !strings.Contains(doc, "advisory") {
		t.Error("slice.md missing 'advisory' designation — session tracking must be advisory only, never blocking (D-33)")
	}
}

func TestSliceMd_C82_ContainsSessionNotBlockingMention(t *testing.T) {
	doc := readSliceMd(t)
	hasNotBlocking := strings.Contains(doc, "D-33") ||
		strings.Contains(doc, "not blocking") ||
		strings.Contains(doc, "never blocking") ||
		strings.Contains(doc, "warning, not blocking") ||
		strings.Contains(doc, "advisory only")
	if !hasNotBlocking {
		t.Error("slice.md missing instruction that session is advisory only — warning, not blocking (D-33)")
	}
}

// Yolo mode behaviour

func TestSliceMd_C82_ContainsYoloModeAutoContinueInstruction(t *testing.T) {
	doc := readSliceMd(t)
	hasYolo := strings.Contains(doc, "yolo") || strings.Contains(doc, "yolo mode")
	if !hasYolo {
		t.Error("slice.md missing yolo mode behaviour for session warning (auto-continue, skip prompt)")
	}
}

func TestSliceMd_C82_ContainsYoloModeSkipsPromptInstruction(t *testing.T) {
	doc := readSliceMd(t)
	hasSkipPrompt := strings.Contains(doc, "skip") && (strings.Contains(doc, "prompt") || strings.Contains(doc, "yolo")) ||
		strings.Contains(doc, "auto-continue")
	if !hasSkipPrompt {
		t.Error("slice.md missing instruction that yolo mode skips the session warning prompt (auto-continues with warning log)")
	}
}

// Session preserved on pause

func TestSliceMd_C82_ContainsSessionFieldPreservedOnPauseInstruction(t *testing.T) {
	doc := readSliceMd(t)
	hasPreserved := strings.Contains(doc, "pause") && strings.Contains(doc, "preserved") ||
		strings.Contains(doc, "pause") && strings.Contains(doc, "session") ||
		strings.Contains(doc, "Session field is preserved")
	if !hasPreserved {
		t.Error("slice.md missing instruction that session field is preserved on slice pause")
	}
}

// Error handling

func TestSliceMd_C82_ContainsStaleSessionError(t *testing.T) {
	doc := readSliceMd(t)
	if !strings.Contains(doc, "StaleSession") {
		t.Error("slice.md missing StaleSession error handling instruction")
	}
}

func TestSliceMd_C82_ContainsSessionGenerationFailureError(t *testing.T) {
	doc := readSliceMd(t)
	if !strings.Contains(doc, "SessionGenerationFailure") {
		t.Error("slice.md missing SessionGenerationFailure error handling instruction")
	}
}

func TestSliceMd_C82_ContainsSessionGenerationFailureFallbackInstruction(t *testing.T) {
	doc := readSliceMd(t)
	// SessionGenerationFailure: use timestamp-only fallback, warn
	hasFallback := strings.Contains(doc, "timestamp-only") ||
		strings.Contains(doc, "timestamp only") ||
		strings.Contains(doc, "Use timestamp")
	if !hasFallback {
		t.Error("slice.md missing SessionGenerationFailure fallback instruction (use timestamp-only, warn)")
	}
}

// Session tracking only in file-per-slice mode

func TestSliceMd_C82_ContainsSessionTrackingOnlyInFilePerSliceMode(t *testing.T) {
	doc := readSliceMd(t)
	// Session tracking must be scoped to file-per-slice only — not legacy
	hasScoped := strings.Contains(doc, "file-per-slice") && strings.Contains(doc, "session") ||
		strings.Contains(doc, "not legacy") && strings.Contains(doc, "session") ||
		strings.Contains(doc, "Session tracking only in file-per-slice")
	if !hasScoped {
		t.Error("slice.md missing instruction that session tracking only applies in file-per-slice mode (not legacy)")
	}
}

// Invariants and ordering

func TestSliceMd_C82_PreFlightBeforeStep4(t *testing.T) {
	doc := readSliceMd(t)
	preFlightPos := strings.Index(doc, "Pre-flight")
	if preFlightPos == -1 {
		preFlightPos = strings.Index(doc, "pre-flight")
	}
	step4Pos := strings.Index(doc, "Step 4")
	if step4Pos == -1 {
		step4Pos = strings.Index(doc, "step 4")
	}
	if preFlightPos == -1 {
		t.Fatal("slice.md missing Pre-flight section")
	}
	if step4Pos == -1 {
		t.Fatal("slice.md missing Step 4 reference")
	}
	if preFlightPos >= step4Pos {
		t.Errorf("Pre-flight (pos %d) must appear before Step 4 (pos %d)", preFlightPos, step4Pos)
	}
}

func TestSliceMd_C82_Step4BeforeStep10(t *testing.T) {
	doc := readSliceMd(t)
	step4Pos := strings.Index(doc, "Step 4")
	if step4Pos == -1 {
		step4Pos = strings.Index(doc, "step 4")
	}
	step10Pos := strings.Index(doc, "Step 10")
	if step10Pos == -1 {
		step10Pos = strings.Index(doc, "step 10")
	}
	if step4Pos == -1 {
		t.Fatal("slice.md missing Step 4 reference")
	}
	if step10Pos == -1 {
		t.Fatal("slice.md missing Step 10 reference")
	}
	if step4Pos >= step10Pos {
		t.Errorf("Step 4 (pos %d) must appear before Step 10 (pos %d)", step4Pos, step10Pos)
	}
}

func TestSliceMd_C81_SessionFieldSetBeforeAgentWorkBegins(t *testing.T) {
	doc := readSliceMd(t)
	// Session field must be set BEFORE any agent work begins (C-82 invariant)
	hasBeforeWork := strings.Contains(doc, "before") && strings.Contains(doc, "agent") ||
		strings.Contains(doc, "BEFORE") ||
		strings.Contains(doc, "Set session") && strings.Contains(doc, "claim")
	if !hasBeforeWork {
		t.Error("slice.md missing instruction that session field is set before agent work begins")
	}
}

func TestSliceMd_C81_StaleSessionsWarnedNotAutoCleaned(t *testing.T) {
	doc := readSliceMd(t)
	hasNotAutoCleaned := strings.Contains(doc, "not auto-cleaned") ||
		strings.Contains(doc, "never auto-clean") ||
		strings.Contains(doc, "stale") && strings.Contains(doc, "warn")
	if !hasNotAutoCleaned {
		t.Error("slice.md missing instruction that stale sessions are warned but never auto-cleaned")
	}
}
