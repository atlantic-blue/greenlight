package installer_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// C-57 Tests: SliceCheckpointTags
// These tests verify that src/commands/gl/slice.md contains the required
// checkpoint tag creation behaviour as defined in contract C-57.

func TestSliceMd_ContainsCheckpointTagCreation(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "greenlight/checkpoint/") {
		t.Error("slice.md missing checkpoint tag namespace 'greenlight/checkpoint/'")
	}

	if !strings.Contains(doc, "git tag") {
		t.Error("slice.md missing 'git tag' command for checkpoint tag creation")
	}
}

func TestSliceMd_ContainsCheckpointBeforeImplementation(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	checkpointPos := strings.Index(doc, "greenlight/checkpoint/")
	if checkpointPos == -1 {
		t.Fatal("slice.md missing checkpoint tag namespace 'greenlight/checkpoint/'")
	}

	// Step 1 (Write Tests) is the earliest agent step; implementation comes after.
	// The checkpoint must appear before "Step 1: Write Tests" which is the first agent invocation.
	stepOnePos := strings.Index(doc, "Step 1: Write Tests")
	if stepOnePos == -1 {
		t.Fatal("slice.md missing 'Step 1: Write Tests' section")
	}

	if checkpointPos >= stepOnePos {
		t.Errorf(
			"checkpoint tag creation (pos %d) must appear BEFORE 'Step 1: Write Tests' (pos %d) — checkpoint must be created before any test writing or implementation",
			checkpointPos,
			stepOnePos,
		)
	}
}

func TestSliceMd_ContainsCheckpointIdempotent(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	// Idempotent behaviour requires deleting a pre-existing tag before re-creating it.
	// The checkpoint section must describe removing the old tag with 'git tag -d'.
	checkpointPos := strings.Index(doc, "greenlight/checkpoint/")
	if checkpointPos == -1 {
		t.Fatal("slice.md missing checkpoint tag namespace 'greenlight/checkpoint/'")
	}

	// Find the end of the checkpoint section: the next top-level step heading.
	stepOnePos := strings.Index(doc, "Step 1: Write Tests")
	if stepOnePos == -1 {
		t.Fatal("slice.md missing 'Step 1: Write Tests' section")
	}

	checkpointSection := doc[checkpointPos:stepOnePos]

	if !strings.Contains(checkpointSection, "git tag -d") {
		t.Error("slice.md checkpoint section missing 'git tag -d' for idempotent tag replacement — old tag for same slice must be removed before creating a fresh one")
	}
}

func TestSliceMd_ContainsCheckpointContextForImplementer(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	// The implementer Task block must carry the checkpoint tag as an XML context block.
	// Per C-57: "Pass checkpoint_tag to implementer context as XML block."
	// Locate the implementer Task block and confirm a checkpoint XML block exists within it.
	implementerTaskStart := strings.Index(doc, "Step 3: Implement")
	if implementerTaskStart == -1 {
		t.Fatal("slice.md missing 'Step 3: Implement' section")
	}

	// The implementer section runs until the next top-level step.
	stepFourPos := strings.Index(doc[implementerTaskStart:], "Step 4:")
	if stepFourPos == -1 {
		t.Fatal("slice.md missing 'Step 4:' section after implementer step")
	}

	implementerSection := doc[implementerTaskStart : implementerTaskStart+stepFourPos]

	// Must contain an XML context block for the checkpoint tag, not just the word "checkpoint"
	// from an unrelated reference (e.g. references/checkpoint-protocol.md).
	hasCheckpointXMLBlock := strings.Contains(implementerSection, "<checkpoint") ||
		strings.Contains(implementerSection, "checkpoint_tag")

	if !hasCheckpointXMLBlock {
		t.Error("slice.md implementer Task block missing checkpoint XML context — checkpoint tag must be passed as '<checkpoint>' XML block or 'checkpoint_tag' field to the implementer, per C-57")
	}
}

// C-58 Tests: SliceRollbackIntegration
// These tests verify that src/commands/gl/slice.md contains the required
// circuit break handling and rollback behaviour as defined in contract C-58.

func TestSliceMd_ContainsCircuitBreakHandling(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	hasUpperCase := strings.Contains(doc, "CIRCUIT BREAK")
	hasLowerCase := strings.Contains(doc, "circuit break")

	if !hasUpperCase && !hasLowerCase {
		t.Error("slice.md missing circuit break handling — must contain 'CIRCUIT BREAK' or 'circuit break'")
	}
}

func TestSliceMd_ContainsRollbackCommand(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "git checkout") {
		t.Error("slice.md missing 'git checkout' rollback command")
	}

	if !strings.Contains(doc, "greenlight/checkpoint") {
		t.Error("slice.md missing 'greenlight/checkpoint' reference in rollback command")
	}

	// Both must appear in the same section — find "git checkout" and confirm
	// "greenlight/checkpoint" appears nearby (within the same rollback context).
	gitCheckoutPos := strings.Index(doc, "git checkout")
	checkpointRefPos := strings.LastIndex(doc, "greenlight/checkpoint")

	if gitCheckoutPos == -1 || checkpointRefPos == -1 {
		t.Fatal("could not locate both 'git checkout' and 'greenlight/checkpoint' in slice.md")
	}

	// The checkpoint reference used in rollback must appear after the first checkpoint
	// creation and within a reasonable proximity to the git checkout command.
	// Allow up to 2000 characters distance as a generous window for the same section.
	distance := checkpointRefPos - gitCheckoutPos
	if distance < 0 {
		distance = gitCheckoutPos - checkpointRefPos
	}
	if distance > 2000 {
		t.Errorf(
			"'git checkout' (pos %d) and 'greenlight/checkpoint' (pos %d) are too far apart (%d chars) — they should appear together in the rollback section",
			gitCheckoutPos,
			checkpointRefPos,
			distance,
		)
	}
}

func TestSliceMd_ContainsRecoveryOptions(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	// All three recovery options must be present in the circuit break section.
	// Per C-58: guidance + retry, spawn debugger, pause.
	if !strings.Contains(doc, "guidance") {
		t.Error("slice.md missing 'guidance' recovery option in circuit break handling")
	}

	if !strings.Contains(doc, "retry") {
		t.Error("slice.md missing 'retry' recovery option in circuit break handling")
	}

	if !strings.Contains(doc, "debugger") {
		t.Error("slice.md missing 'debugger' recovery option in circuit break handling")
	}

	if !strings.Contains(doc, "pause") && !strings.Contains(doc, "Pause") {
		t.Error("slice.md missing 'pause' recovery option in circuit break handling")
	}
}

func TestSliceMd_ContainsTagCleanupOnCompletion(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/slice.md"))
	if err != nil {
		t.Fatalf("failed to read slice.md: %v", err)
	}

	doc := string(content)

	// Locate the completion step.
	stepTenPos := strings.Index(doc, "Step 10: Complete")
	if stepTenPos == -1 {
		t.Fatal("slice.md missing 'Step 10: Complete' section")
	}

	completionSection := doc[stepTenPos:]

	// Tag cleanup on success is best-effort (C-58 invariant) but must be documented.
	if !strings.Contains(completionSection, "git tag -d") {
		t.Error("slice.md Step 10 (Complete) missing 'git tag -d' — checkpoint tag must be cleaned up on successful slice completion")
	}

	if !strings.Contains(completionSection, "greenlight/checkpoint") {
		t.Error("slice.md Step 10 (Complete) missing 'greenlight/checkpoint' reference — cleanup must target the checkpoint tag for this slice")
	}
}
