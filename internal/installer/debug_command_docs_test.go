package installer_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// C-59 Tests: DebugCommand
// These tests verify that src/commands/gl/debug.md exists and contains
// all required sections as defined in contract C-59.

func TestDebugMd_Exists(t *testing.T) {
	path := filepath.Join(projectRoot(), "src/commands/gl/debug.md")
	_, err := os.Stat(path)
	if err != nil {
		t.Fatalf("debug.md does not exist at %s: %v", path, err)
	}
}

func TestDebugMd_ContainsFrontmatter(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/debug.md"))
	if err != nil {
		t.Fatalf("failed to read debug.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "gl-debug") && !strings.Contains(doc, "gl:debug") {
		t.Error("debug.md missing command name 'gl-debug' or 'gl:debug' in frontmatter")
	}

	if !strings.Contains(doc, "diagnostic") {
		t.Error("debug.md missing 'diagnostic' in description")
	}
}

func TestDebugMd_ContainsSliceTargeting(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/debug.md"))
	if err != nil {
		t.Fatalf("failed to read debug.md: %v", err)
	}

	doc := string(content)

	// Must support optional slice_id argument
	if !strings.Contains(doc, "slice_id") {
		t.Error("debug.md missing 'slice_id' — must support optional slice argument")
	}

	// Must read STATE.md for current slice when no argument given
	if !strings.Contains(doc, "STATE.md") {
		t.Error("debug.md missing 'STATE.md' reference — must determine target slice from state")
	}

	// Must handle no active slice case
	if !strings.Contains(doc, "No active slice") || !strings.Contains(doc, "no active slice") {
		hasNoSliceHandling := strings.Contains(doc, "No active slice") ||
			strings.Contains(doc, "no active slice") ||
			strings.Contains(doc, "NoActiveSlice")
		if !hasNoSliceHandling {
			t.Error("debug.md missing 'No active slice' error handling")
		}
	}
}

func TestDebugMd_ContainsDiagnosticContextGathering(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/debug.md"))
	if err != nil {
		t.Fatalf("failed to read debug.md: %v", err)
	}

	doc := string(content)

	// Must run test suite
	if !strings.Contains(doc, "test") {
		t.Error("debug.md missing test execution — must run test suite for current state")
	}

	// Must gather git context
	if !strings.Contains(doc, "git log") {
		t.Error("debug.md missing 'git log' — must gather recent git activity")
	}

	if !strings.Contains(doc, "git diff") {
		t.Error("debug.md missing 'git diff' — must gather uncommitted changes")
	}

	// Must check for checkpoint tag
	if !strings.Contains(doc, "greenlight/checkpoint") {
		t.Error("debug.md missing 'greenlight/checkpoint' — must check for checkpoint tag")
	}

	// Must read contracts
	if !strings.Contains(doc, "CONTRACTS.md") || !strings.Contains(doc, "contracts") {
		hasContracts := strings.Contains(doc, "CONTRACTS.md") ||
			strings.Contains(doc, "contracts")
		if !hasContracts {
			t.Error("debug.md missing contract reading — must read contracts for the slice")
		}
	}
}

func TestDebugMd_ContainsStructuredDiagnosticReport(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/debug.md"))
	if err != nil {
		t.Fatalf("failed to read debug.md: %v", err)
	}

	doc := string(content)

	// Must produce a DIAGNOSTIC REPORT
	if !strings.Contains(doc, "DIAGNOSTIC REPORT") && !strings.Contains(doc, "Diagnostic Report") {
		t.Error("debug.md missing 'DIAGNOSTIC REPORT' heading in output format")
	}

	// Report must contain required sections
	requiredSections := []string{
		"Current State",
		"Test Results",
		"Recovery Options",
	}

	for _, section := range requiredSections {
		if !strings.Contains(doc, section) {
			t.Errorf("debug.md diagnostic report missing required section: %s", section)
		}
	}

	// Must include failing test details
	if !strings.Contains(doc, "Failing") && !strings.Contains(doc, "failing") {
		t.Error("debug.md missing failing test details in diagnostic report")
	}
}

func TestDebugMd_ContainsRecoveryOptions(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/debug.md"))
	if err != nil {
		t.Fatalf("failed to read debug.md: %v", err)
	}

	doc := string(content)

	// Must list recovery options for the user
	if !strings.Contains(doc, "checkpoint") {
		t.Error("debug.md missing checkpoint reference in recovery options")
	}

	if !strings.Contains(doc, "pause") || !strings.Contains(doc, "Pause") {
		hasPause := strings.Contains(doc, "pause") || strings.Contains(doc, "Pause")
		if !hasPause {
			t.Error("debug.md missing 'pause' recovery option")
		}
	}
}

func TestDebugMd_IsReadOnly(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/debug.md"))
	if err != nil {
		t.Fatalf("failed to read debug.md: %v", err)
	}

	doc := string(content)

	// Must explicitly state read-only nature
	hasReadOnly := strings.Contains(doc, "read-only") ||
		strings.Contains(doc, "Read-only") ||
		strings.Contains(doc, "no files written") ||
		strings.Contains(doc, "No files written")

	if !hasReadOnly {
		t.Error("debug.md missing read-only declaration — command must not write files or modify state")
	}

	// Must not spawn subagents
	hasNoSubagents := strings.Contains(doc, "no subagent") ||
		strings.Contains(doc, "not spawn") ||
		strings.Contains(doc, "Does not spawn") ||
		strings.Contains(doc, "direct read")

	if !hasNoSubagents {
		t.Error("debug.md missing no-subagent declaration — command must not spawn subagents")
	}
}

func TestDebugMd_ContainsNoConfigHandling(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/commands/gl/debug.md"))
	if err != nil {
		t.Fatalf("failed to read debug.md: %v", err)
	}

	doc := string(content)

	// Must handle missing config.json
	if !strings.Contains(doc, "config.json") {
		t.Error("debug.md missing 'config.json' reference — must handle NoConfig error state")
	}

	// Must handle partial diagnostics (graceful degradation)
	hasGracefulDegradation := strings.Contains(doc, "partial") ||
		strings.Contains(doc, "Continue") ||
		strings.Contains(doc, "continue")

	if !hasGracefulDegradation {
		t.Error("debug.md missing graceful degradation — must continue with partial diagnostic when some context is unavailable")
	}
}
