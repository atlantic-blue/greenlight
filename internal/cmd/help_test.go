package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/atlantic-blue/greenlight/internal/cmd"
)

// ----------------------------------------------------------------------------
// C-100 — RunHelp: command listing (always present)
// ----------------------------------------------------------------------------

// TestRunHelp_ShowsAllCommandCategories verifies that the output contains all
// four category headings regardless of whether a .greenlight/ project exists.
func TestRunHelp_ShowsAllCommandCategories(t *testing.T) {
	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	var buf bytes.Buffer
	cmd.RunHelp([]string{}, &buf)

	output := buf.String()

	categories := []string{
		"lifecycle",
		"building",
		"state",
		"admin",
	}

	for _, category := range categories {
		if !strings.Contains(strings.ToLower(output), category) {
			t.Errorf("output missing category %q:\n%s", category, output)
		}
	}
}

// TestRunHelp_ShowsLifecycleCommands verifies that init, design, and roadmap
// appear under the project lifecycle category.
func TestRunHelp_ShowsLifecycleCommands(t *testing.T) {
	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	var buf bytes.Buffer
	cmd.RunHelp([]string{}, &buf)

	output := buf.String()

	lifecycleCommands := []string{"init", "design", "roadmap"}
	for _, command := range lifecycleCommands {
		if !strings.Contains(output, command) {
			t.Errorf("output missing lifecycle command %q:\n%s", command, output)
		}
	}
}

// TestRunHelp_ShowsBuildingCommands verifies that slice appears under the
// building category.
func TestRunHelp_ShowsBuildingCommands(t *testing.T) {
	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	var buf bytes.Buffer
	cmd.RunHelp([]string{}, &buf)

	output := buf.String()

	if !strings.Contains(output, "slice") {
		t.Errorf("output missing building command 'slice':\n%s", output)
	}
}

// TestRunHelp_ShowsStateCommands verifies that status and changelog appear
// under the state & progress category.
func TestRunHelp_ShowsStateCommands(t *testing.T) {
	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	var buf bytes.Buffer
	cmd.RunHelp([]string{}, &buf)

	output := buf.String()

	stateCommands := []string{"status", "changelog"}
	for _, command := range stateCommands {
		if !strings.Contains(output, command) {
			t.Errorf("output missing state command %q:\n%s", command, output)
		}
	}
}

// TestRunHelp_ShowsAdminCommands verifies that install, uninstall, check,
// version, and help all appear under the admin category.
func TestRunHelp_ShowsAdminCommands(t *testing.T) {
	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	var buf bytes.Buffer
	cmd.RunHelp([]string{}, &buf)

	output := buf.String()

	adminCommands := []string{"install", "uninstall", "check", "version", "help"}
	for _, command := range adminCommands {
		if !strings.Contains(output, command) {
			t.Errorf("output missing admin command %q:\n%s", command, output)
		}
	}
}

// ----------------------------------------------------------------------------
// C-100 — RunHelp: exit code invariant
// ----------------------------------------------------------------------------

// TestRunHelp_AlwaysReturnsExitCode0 verifies that RunHelp returns 0 when
// called outside any greenlight project (no .greenlight/ directory).
func TestRunHelp_AlwaysReturnsExitCode0(t *testing.T) {
	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	var buf bytes.Buffer
	exitCode := cmd.RunHelp([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 outside a project, got %d", exitCode)
	}
}

// TestRunHelp_ReturnsExitCode0WithProject verifies that RunHelp returns 0 when
// called inside a greenlight project (with .greenlight/ directory).
func TestRunHelp_ReturnsExitCode0WithProject(t *testing.T) {
	slices := []testSlice{
		{id: "S-01", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z",
			tests: 5, securityTests: 0},
	}
	setupTestProject(t, slices, minimalGraph([]string{"S-01"}))

	var buf bytes.Buffer
	exitCode := cmd.RunHelp([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 inside a project, got %d", exitCode)
	}
}

// ----------------------------------------------------------------------------
// C-100 — RunHelp: project state detection
// ----------------------------------------------------------------------------

// TestRunHelp_NoGreenlightDir_SuggestsInit verifies that when no .greenlight/
// directory is present, the output contains a suggestion to run 'gl init'.
func TestRunHelp_NoGreenlightDir_SuggestsInit(t *testing.T) {
	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	var buf bytes.Buffer
	cmd.RunHelp([]string{}, &buf)

	output := buf.String()

	if !strings.Contains(output, "gl init") {
		t.Errorf("expected output to suggest 'gl init' when no .greenlight/ dir:\n%s", output)
	}
}

// TestRunHelp_WithProject_ShowsSliceCounts verifies that with 3 slices (2
// complete, 1 pending), the state summary includes "3 slices" and "2 complete".
func TestRunHelp_WithProject_ShowsSliceCounts(t *testing.T) {
	slices := []testSlice{
		{id: "S-01", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z"},
		{id: "S-02", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-02", updated: "2026-01-02T00:00:00Z"},
		{id: "S-03", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	setupTestProject(t, slices, minimalGraph([]string{"S-01", "S-02", "S-03"}))

	var buf bytes.Buffer
	cmd.RunHelp([]string{}, &buf)

	output := buf.String()

	if !strings.Contains(output, "3 slices") {
		t.Errorf("expected '3 slices' in state summary:\n%s", output)
	}

	if !strings.Contains(output, "2 complete") {
		t.Errorf("expected '2 complete' in state summary:\n%s", output)
	}
}

// TestRunHelp_WithProject_ShowsReadyCount verifies that with pending slices
// whose dependencies are met, the state summary includes a ready count.
func TestRunHelp_WithProject_ShowsReadyCount(t *testing.T) {
	slices := []testSlice{
		{id: "S-01", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z"},
		{id: "S-02", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: "", deps: "S-01"},
		{id: "S-03", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: "", deps: "S-01"},
	}
	graphJSON := graphWithDeps(map[string][]string{
		"S-01": {},
		"S-02": {"S-01"},
		"S-03": {"S-01"},
	})
	setupTestProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunHelp([]string{}, &buf)

	output := buf.String()

	// S-02 and S-03 are ready because S-01 is complete; expect "2 ready".
	if !strings.Contains(output, "ready") {
		t.Errorf("expected a ready count in the state summary:\n%s", output)
	}

	if !strings.Contains(output, "2 ready") {
		t.Errorf("expected '2 ready' in state summary (S-02 and S-03 both unblocked):\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-100 — RunHelp: error resilience
// ----------------------------------------------------------------------------

// TestRunHelp_CorruptSliceFiles_StillShowsCommands verifies that malformed
// slice frontmatter does not prevent the command listing from appearing.
func TestRunHelp_CorruptSliceFiles_StillShowsCommands(t *testing.T) {
	tmpDir := t.TempDir()

	slicesDir := filepath.Join(tmpDir, ".greenlight", "slices")
	if mkdirError := os.MkdirAll(slicesDir, 0o755); mkdirError != nil {
		t.Fatalf("failed to create slices dir: %v", mkdirError)
	}

	// Write a slice file with invalid/corrupt frontmatter.
	corruptPath := filepath.Join(slicesDir, "S-01.md")
	if writeError := os.WriteFile(corruptPath, []byte("not valid frontmatter at all\n{{{"), 0o644); writeError != nil {
		t.Fatalf("failed to write corrupt slice file: %v", writeError)
	}

	t.Chdir(tmpDir)

	var buf bytes.Buffer
	exitCode := cmd.RunHelp([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 even with corrupt slice files, got %d", exitCode)
	}

	output := buf.String()

	// Command listing must always be present regardless of state read errors.
	commands := []string{"init", "slice", "status", "help"}
	for _, command := range commands {
		if !strings.Contains(output, command) {
			t.Errorf("output missing command %q despite corrupt slice file:\n%s", command, output)
		}
	}
}

// TestRunHelp_EmptySlicesDir_StillShowsCommands verifies that when
// .greenlight/slices/ exists but is empty, commands are still shown.
func TestRunHelp_EmptySlicesDir_StillShowsCommands(t *testing.T) {
	tmpDir := t.TempDir()

	slicesDir := filepath.Join(tmpDir, ".greenlight", "slices")
	if mkdirError := os.MkdirAll(slicesDir, 0o755); mkdirError != nil {
		t.Fatalf("failed to create slices dir: %v", mkdirError)
	}

	t.Chdir(tmpDir)

	var buf bytes.Buffer
	exitCode := cmd.RunHelp([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 with empty slices dir, got %d", exitCode)
	}

	output := buf.String()

	commands := []string{"init", "slice", "status", "help"}
	for _, command := range commands {
		if !strings.Contains(output, command) {
			t.Errorf("output missing command %q with empty slices dir:\n%s", command, output)
		}
	}
}

// ----------------------------------------------------------------------------
// C-100 — RunHelp: output writer contract
// ----------------------------------------------------------------------------

// TestRunHelp_WritesToProvidedWriter verifies that all RunHelp output is
// directed to the provided io.Writer and not written to os.Stdout directly.
func TestRunHelp_WritesToProvidedWriter(t *testing.T) {
	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	var buf bytes.Buffer
	cmd.RunHelp([]string{}, &buf)

	if buf.Len() == 0 {
		t.Error("expected output to be written to the provided writer, got empty buffer")
	}
}
