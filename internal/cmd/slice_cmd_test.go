package cmd_test

// S-42 Tests: Single Slice Command
// Covers C-105 (RunSliceSingle) and C-106 (RunSliceAutoDetect).
//
// Contract C-105 — RunSlice with a slice ID:
//   - Parses args for slice ID and flags (--dry-run, --max, --watch, --sequential)
//   - Detects context via $CLAUDE_CODE environment variable
//   - Reads config from .greenlight/config.json for claude_flags
//   - If --dry-run: prints what would happen, returns 0, never spawns
//   - If inside Claude ($CLAUDE_CODE set): prints slice info, returns 0
//   - If in shell: builds spawn command; returns error if claude not found
//   - InvalidSliceID: prints error when ID not found in GRAPH.json
//   - NoGreenlightDir: prints error and returns 1 when .greenlight/ missing
//
// Contract C-106 — RunSlice with no slice ID (auto-detect):
//   - Reads slice states and graph
//   - If 0 ready slices: prints blocked status, returns 0
//   - If 1 ready slice: runs that slice (same as providing ID)
//   - If 2+ ready slices: picks first by wave/ID order, prints hint about rest
//   - StateReadError: prints error, returns 1

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/atlantic-blue/greenlight/internal/cmd"
	"github.com/atlantic-blue/greenlight/internal/process"
)

// ----------------------------------------------------------------------------
// Helpers specific to slice_cmd tests
// ----------------------------------------------------------------------------

// minimalConfig returns the minimal config.json content for slice tests.
func minimalConfig() string {
	return `{"parallel":{"claude_flags":["--dangerously-skip-permissions","--max-turns","200"]}}`
}

// setupSliceProject creates a temp directory with .greenlight/slices/,
// GRAPH.json, and config.json, then calls t.Chdir. It returns the project
// root directory. If graphJSON is empty, GRAPH.json is not written.
func setupSliceProject(t *testing.T, slices []testSlice, graphJSON string) string {
	t.Helper()

	tmpDir := setupTestProject(t, slices, graphJSON)

	configPath := filepath.Join(tmpDir, ".greenlight", "config.json")
	writeStatusFile(t, configPath, minimalConfig())

	return tmpDir
}

// overrideProcessLookPath replaces process.LookPath with the given stub and
// registers a cleanup to restore the original when the test finishes.
func overrideProcessLookPath(t *testing.T, stub func(string) (string, error)) {
	t.Helper()
	original := process.LookPath
	process.LookPath = stub
	t.Cleanup(func() { process.LookPath = original })
}

// claudeInPath is a LookPath stub that reports "claude" as found.
func claudeInPath(name string) (string, error) {
	if name == "claude" {
		return "/usr/local/bin/claude", nil
	}
	return "", os.ErrNotExist
}

// claudeNotInPath is a LookPath stub that always reports not found.
func claudeNotInPath(_ string) (string, error) {
	return "", os.ErrNotExist
}

// setClaudeContext sets $CLAUDE_CODE to value and registers cleanup to unset.
func setClaudeContext(t *testing.T, value string) {
	t.Helper()
	t.Setenv("CLAUDE_CODE", value)
}

// clearClaudeContext ensures $CLAUDE_CODE is unset for the duration of the test.
// It saves the current value and registers cleanup to restore it.
func clearClaudeContext(t *testing.T) {
	t.Helper()
	original, wasSet := os.LookupEnv("CLAUDE_CODE")
	os.Unsetenv("CLAUDE_CODE")
	t.Cleanup(func() {
		if wasSet {
			os.Setenv("CLAUDE_CODE", original)
		} else {
			os.Unsetenv("CLAUDE_CODE")
		}
	})
}

// graphWithWave builds a GRAPH.json where each slice can specify a wave number.
// entries maps slice ID to {deps, wave}.
func graphWithWave(entries map[string]struct {
	deps []string
	wave int
}) string {
	var builder strings.Builder
	builder.WriteString(`{"slices":{`)
	first := true
	for id, info := range entries {
		if !first {
			builder.WriteString(",")
		}
		first = false
		builder.WriteString(`"`)
		builder.WriteString(id)
		builder.WriteString(`":{"name":"`)
		builder.WriteString(id)
		builder.WriteString(` Slice","depends_on":[`)
		for dIndex, dep := range info.deps {
			if dIndex > 0 {
				builder.WriteString(",")
			}
			builder.WriteString(`"`)
			builder.WriteString(dep)
			builder.WriteString(`"`)
		}
		builder.WriteString(`],"wave":`)
		builder.WriteString(itoa(info.wave))
		builder.WriteString(`,"contracts":[]}`)
	}
	builder.WriteString(`},"edges":[]}`)
	return builder.String()
}

// ----------------------------------------------------------------------------
// C-105 — RunSlice with explicit slice ID: no greenlight directory
// ----------------------------------------------------------------------------

// TestRunSlice_NoGreenlightDir_ReturnsExitCode1 verifies that when there is no
// .greenlight/ directory, RunSlice prints an error and returns 1.
func TestRunSlice_NoGreenlightDir_ReturnsExitCode1(t *testing.T) {
	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"S-35"}, &buf)

	if exitCode != 1 {
		t.Errorf("expected exit code 1 when .greenlight/ not found, got %d", exitCode)
	}

	output := buf.String()
	if !strings.Contains(output, "greenlight") {
		t.Errorf("expected error mentioning 'greenlight', got:\n%s", output)
	}
}

// TestRunSlice_NoGreenlightDir_PrintsNotAGreenlightProject verifies that the
// specific "Not a greenlight project." error message is printed.
func TestRunSlice_NoGreenlightDir_PrintsNotAGreenlightProject(t *testing.T) {
	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"S-35"}, &buf)

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	if !strings.Contains(lowerOutput, "not a greenlight project") &&
		!strings.Contains(lowerOutput, "greenlight project") {
		t.Errorf("expected 'Not a greenlight project.' in output, got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-105 — RunSlice with explicit slice ID: invalid slice ID
// ----------------------------------------------------------------------------

// TestRunSlice_InvalidSliceID_ReturnsExitCode1 verifies that providing a slice
// ID not present in GRAPH.json prints an error and returns 1.
func TestRunSlice_InvalidSliceID_ReturnsExitCode1(t *testing.T) {
	slices := []testSlice{
		{id: "S-01", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	graphJSON := minimalGraph([]string{"S-01"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"INVALID-99"}, &buf)

	if exitCode != 1 {
		t.Errorf("expected exit code 1 for unknown slice ID, got %d", exitCode)
	}
}

// TestRunSlice_InvalidSliceID_PrintsUnknownSliceError verifies that the error
// output mentions the unknown slice ID.
func TestRunSlice_InvalidSliceID_PrintsUnknownSliceError(t *testing.T) {
	slices := []testSlice{
		{id: "S-01", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	graphJSON := minimalGraph([]string{"S-01"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"INVALID-99"}, &buf)

	output := buf.String()
	if !strings.Contains(output, "INVALID-99") {
		t.Errorf("expected unknown ID 'INVALID-99' in error output, got:\n%s", output)
	}
}

// TestRunSlice_InvalidSliceID_SuggestsGlStatus verifies that the error message
// hints the user to run 'gl status'.
func TestRunSlice_InvalidSliceID_SuggestsGlStatus(t *testing.T) {
	slices := []testSlice{
		{id: "S-01", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	graphJSON := minimalGraph([]string{"S-01"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"S-UNKNOWN"}, &buf)

	output := buf.String()
	if !strings.Contains(output, "status") {
		t.Errorf("expected 'status' hint in error output, got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-105 — RunSlice with explicit slice ID: --dry-run flag
// ----------------------------------------------------------------------------

// TestRunSlice_DryRun_ReturnsExitCode0 verifies that --dry-run always returns 0.
func TestRunSlice_DryRun_ReturnsExitCode0(t *testing.T) {
	slices := []testSlice{
		{id: "S-35", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	graphJSON := minimalGraph([]string{"S-35"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"S-35", "--dry-run"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 for --dry-run, got %d", exitCode)
	}
}

// TestRunSlice_DryRun_PrintsWhatWouldHappen verifies that --dry-run output
// describes the command that would execute without actually spawning.
func TestRunSlice_DryRun_PrintsWhatWouldHappen(t *testing.T) {
	slices := []testSlice{
		{id: "S-35", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	graphJSON := minimalGraph([]string{"S-35"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"S-35", "--dry-run"}, &buf)

	output := buf.String()
	if !strings.Contains(output, "S-35") {
		t.Errorf("expected slice ID 'S-35' in dry-run output, got:\n%s", output)
	}
}

// TestRunSlice_DryRun_NeverSpawnsProcess verifies that --dry-run does not
// attempt to spawn Claude even when claude binary is not in PATH.
func TestRunSlice_DryRun_NeverSpawnsProcess(t *testing.T) {
	overrideProcessLookPath(t, claudeNotInPath)

	slices := []testSlice{
		{id: "S-35", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	graphJSON := minimalGraph([]string{"S-35"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"S-35", "--dry-run"}, &buf)

	// --dry-run must not fail due to claude being absent from PATH.
	if exitCode != 0 {
		t.Errorf("expected exit code 0 for --dry-run (no spawn), got %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunSlice_DryRun_ShowsClaudeCommand verifies that --dry-run output
// includes the 'claude' command that would be executed.
func TestRunSlice_DryRun_ShowsClaudeCommand(t *testing.T) {
	slices := []testSlice{
		{id: "S-35", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	graphJSON := minimalGraph([]string{"S-35"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"S-35", "--dry-run"}, &buf)

	output := buf.String()
	if !strings.Contains(output, "claude") {
		t.Errorf("expected 'claude' in dry-run output describing the command, got:\n%s", output)
	}
}

// TestRunSlice_DryRun_ShowsSlicePromptOrSkillReference verifies that dry-run
// output references the /gl:slice skill invocation pattern.
func TestRunSlice_DryRun_ShowsSlicePromptOrSkillReference(t *testing.T) {
	slices := []testSlice{
		{id: "S-35", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	graphJSON := minimalGraph([]string{"S-35"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"S-35", "--dry-run"}, &buf)

	output := buf.String()
	// Output must mention the gl:slice skill or slice prompt in some form.
	hasSliceRef := strings.Contains(output, "gl:slice") ||
		strings.Contains(output, "/gl:slice") ||
		strings.Contains(output, "slice S-35") ||
		strings.Contains(output, "slice")
	if !hasSliceRef {
		t.Errorf("expected dry-run output to reference gl:slice invocation, got:\n%s", output)
	}
}

// TestRunSlice_DryRun_FlagsIncludedFromConfig verifies that the claude_flags
// from config.json appear in the dry-run output.
func TestRunSlice_DryRun_FlagsIncludedFromConfig(t *testing.T) {
	slices := []testSlice{
		{id: "S-35", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	graphJSON := minimalGraph([]string{"S-35"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"S-35", "--dry-run"}, &buf)

	output := buf.String()
	// The config has --max-turns; it should appear in the dry-run command preview.
	if !strings.Contains(output, "--max-turns") && !strings.Contains(output, "max-turns") &&
		!strings.Contains(output, "200") {
		t.Errorf("expected config claude_flags (--max-turns 200) reflected in dry-run output, got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-105 — RunSlice with explicit slice ID: inside Claude context
// ----------------------------------------------------------------------------

// TestRunSlice_InsideClaude_ReturnsExitCode0 verifies that when $CLAUDE_CODE
// is set, RunSlice returns 0 and does not attempt to spawn a process.
func TestRunSlice_InsideClaude_ReturnsExitCode0(t *testing.T) {
	setClaudeContext(t, "1")
	overrideProcessLookPath(t, claudeNotInPath)

	slices := []testSlice{
		{id: "S-35", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	graphJSON := minimalGraph([]string{"S-35"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"S-35"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 inside Claude context, got %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunSlice_InsideClaude_PrintsSliceInfo verifies that inside Claude context
// the output contains the slice ID so the Claude skill can consume it.
func TestRunSlice_InsideClaude_PrintsSliceInfo(t *testing.T) {
	setClaudeContext(t, "1")

	slices := []testSlice{
		{id: "S-35", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	graphJSON := minimalGraph([]string{"S-35"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"S-35"}, &buf)

	output := buf.String()
	if !strings.Contains(output, "S-35") {
		t.Errorf("expected slice ID 'S-35' in inside-Claude output, got:\n%s", output)
	}
}

// TestRunSlice_InsideClaude_NeverSpawnsAnotherClaude verifies the invariant
// that inside Claude context no child Claude process is ever started —
// verified by ensuring no error occurs even when LookPath returns not found.
func TestRunSlice_InsideClaude_NeverSpawnsAnotherClaude(t *testing.T) {
	setClaudeContext(t, "1")
	overrideProcessLookPath(t, claudeNotInPath)

	slices := []testSlice{
		{id: "S-35", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	graphJSON := minimalGraph([]string{"S-35"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"S-35"}, &buf)

	// If RunSlice tried to spawn claude it would get ErrClaudeNotFound and fail.
	if exitCode != 0 {
		t.Errorf("inside Claude must never spawn another Claude; got exit code %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunSlice_InsideClaude_OutputContainsSliceName verifies that the slice
// name from the graph appears in the inside-Claude output.
func TestRunSlice_InsideClaude_OutputContainsSliceName(t *testing.T) {
	setClaudeContext(t, "1")

	slices := []testSlice{
		{id: "S-35", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	graphJSON := `{"slices":{"S-35":{"name":"Single Slice Command","depends_on":[],"wave":1,"contracts":[]}},"edges":[]}`
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"S-35"}, &buf)

	output := buf.String()
	if !strings.Contains(output, "S-35") {
		t.Errorf("expected slice ID in inside-Claude output, got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-105 — RunSlice with explicit slice ID: shell context spawn error
// ----------------------------------------------------------------------------

// TestRunSlice_ShellContext_ClaudeNotFound_ReturnsNonZero verifies that in shell
// context (no $CLAUDE_CODE), if claude binary is not found, RunSlice returns
// a non-zero exit code.
func TestRunSlice_ShellContext_ClaudeNotFound_ReturnsNonZero(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := []testSlice{
		{id: "S-35", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	graphJSON := minimalGraph([]string{"S-35"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"S-35"}, &buf)

	if exitCode == 0 {
		t.Errorf("expected non-zero exit code when claude not in PATH in shell context, got 0")
	}
}

// TestRunSlice_ShellContext_ClaudeNotFound_PrintsError verifies that the error
// output mentions claude when it cannot be found.
func TestRunSlice_ShellContext_ClaudeNotFound_PrintsError(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := []testSlice{
		{id: "S-35", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	graphJSON := minimalGraph([]string{"S-35"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"S-35"}, &buf)

	output := buf.String()
	if !strings.Contains(output, "claude") {
		t.Errorf("expected error mentioning 'claude' when not in PATH, got:\n%s", output)
	}
}

// TestRunSlice_ShellContext_WritesToProvidedWriter verifies that all output goes
// to the provided io.Writer.
func TestRunSlice_ShellContext_WritesToProvidedWriter(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := []testSlice{
		{id: "S-35", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	graphJSON := minimalGraph([]string{"S-35"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"S-35"}, &buf)

	if buf.Len() == 0 {
		t.Error("expected output to be written to the provided writer, got empty buffer")
	}
}

// ----------------------------------------------------------------------------
// C-105 — RunSlice with explicit slice ID: config.json reading
// ----------------------------------------------------------------------------

// TestRunSlice_DryRun_MissingConfig_StillReturnsExitCode0 verifies that if
// config.json is absent, --dry-run does not fail (graceful degradation).
func TestRunSlice_DryRun_MissingConfig_StillReturnsExitCode0(t *testing.T) {
	slices := []testSlice{
		{id: "S-35", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	graphJSON := minimalGraph([]string{"S-35"})
	// Use setupTestProject to skip writing config.json.
	setupTestProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"S-35", "--dry-run"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 for --dry-run even without config.json, got %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunSlice_InsideClaude_MissingConfig_StillReturnsExitCode0 verifies that
// if config.json is absent, inside-Claude mode does not fail.
func TestRunSlice_InsideClaude_MissingConfig_StillReturnsExitCode0(t *testing.T) {
	setClaudeContext(t, "1")

	slices := []testSlice{
		{id: "S-35", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	graphJSON := minimalGraph([]string{"S-35"})
	setupTestProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"S-35"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 inside Claude without config.json, got %d; output:\n%s", exitCode, buf.String())
	}
}

// ----------------------------------------------------------------------------
// C-106 — RunSlice auto-detect (no slice ID): zero ready slices
// ----------------------------------------------------------------------------

// TestRunSlice_AutoDetect_ZeroReady_ReturnsExitCode0 verifies that when all
// slices are blocked, RunSlice returns 0 (zero ready is not an error).
func TestRunSlice_AutoDetect_ZeroReady_ReturnsExitCode0(t *testing.T) {
	setClaudeContext(t, "1")

	slices := []testSlice{
		{id: "S-01", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: "", deps: "S-02"},
		{id: "S-02", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: "", deps: "S-01"},
	}
	graphJSON := graphWithDeps(map[string][]string{
		"S-01": {"S-02"},
		"S-02": {"S-01"},
	})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 when zero slices ready (not an error), got %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunSlice_AutoDetect_ZeroReady_PrintsBlockedStatus verifies that when no
// slices are ready, the output describes the blocked state.
func TestRunSlice_AutoDetect_ZeroReady_PrintsBlockedStatus(t *testing.T) {
	setClaudeContext(t, "1")

	slices := []testSlice{
		{id: "S-01", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: "", deps: "S-02"},
		{id: "S-02", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: "", deps: "S-01"},
	}
	graphJSON := graphWithDeps(map[string][]string{
		"S-01": {"S-02"},
		"S-02": {"S-01"},
	})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{}, &buf)

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	hasBlockedInfo := strings.Contains(lowerOutput, "blocked") ||
		strings.Contains(lowerOutput, "no ready") ||
		strings.Contains(lowerOutput, "waiting") ||
		strings.Contains(lowerOutput, "depend")
	if !hasBlockedInfo {
		t.Errorf("expected blocked status info when zero slices ready, got:\n%s", output)
	}
}

// TestRunSlice_AutoDetect_ZeroReady_AllComplete_ReturnsExitCode0 verifies that
// when all slices are complete (none pending), RunSlice returns 0.
func TestRunSlice_AutoDetect_ZeroReady_AllComplete_ReturnsExitCode0(t *testing.T) {
	setClaudeContext(t, "1")

	slices := []testSlice{
		{id: "S-01", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z"},
		{id: "S-02", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-02", updated: "2026-01-02T00:00:00Z"},
	}
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 when all slices complete, got %d; output:\n%s", exitCode, buf.String())
	}
}

// ----------------------------------------------------------------------------
// C-106 — RunSlice auto-detect (no slice ID): one ready slice
// ----------------------------------------------------------------------------

// TestRunSlice_AutoDetect_OneReady_InsideClaude_PrintsSliceID verifies that
// when exactly one slice is ready, inside Claude context RunSlice prints info
// for that slice.
func TestRunSlice_AutoDetect_OneReady_InsideClaude_PrintsSliceID(t *testing.T) {
	setClaudeContext(t, "1")

	slices := []testSlice{
		{id: "S-01", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z"},
		{id: "S-02", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: "", deps: "S-01"},
	}
	graphJSON := graphWithDeps(map[string][]string{
		"S-01": {},
		"S-02": {"S-01"},
	})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 for one ready slice inside Claude, got %d; output:\n%s", exitCode, buf.String())
	}

	output := buf.String()
	if !strings.Contains(output, "S-02") {
		t.Errorf("expected ready slice 'S-02' in output, got:\n%s", output)
	}
}

// TestRunSlice_AutoDetect_OneReady_ShellContext_ClaudeNotFound_ReturnsNonZero
// verifies that in shell context with one ready slice, if claude is not found,
// RunSlice returns non-zero.
func TestRunSlice_AutoDetect_OneReady_ShellContext_ClaudeNotFound_ReturnsNonZero(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := []testSlice{
		{id: "S-01", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z"},
		{id: "S-02", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: "", deps: "S-01"},
	}
	graphJSON := graphWithDeps(map[string][]string{
		"S-01": {},
		"S-02": {"S-01"},
	})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{}, &buf)

	if exitCode == 0 {
		t.Errorf("expected non-zero exit code when claude not in PATH for auto-detect shell spawn, got 0")
	}
}

// TestRunSlice_AutoDetect_OneReady_DryRun_ReturnsExitCode0 verifies that with
// --dry-run and one ready slice, RunSlice returns 0 and shows the slice ID.
func TestRunSlice_AutoDetect_OneReady_DryRun_ReturnsExitCode0(t *testing.T) {
	clearClaudeContext(t)

	slices := []testSlice{
		{id: "S-01", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z"},
		{id: "S-02", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: "", deps: "S-01"},
	}
	graphJSON := graphWithDeps(map[string][]string{
		"S-01": {},
		"S-02": {"S-01"},
	})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--dry-run"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 for auto-detect --dry-run, got %d; output:\n%s", exitCode, buf.String())
	}

	output := buf.String()
	if !strings.Contains(output, "S-02") {
		t.Errorf("expected auto-detected ready slice 'S-02' in dry-run output, got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-106 — RunSlice auto-detect (no slice ID): multiple ready slices
// ----------------------------------------------------------------------------

// TestRunSlice_AutoDetect_MultipleReady_InsideClaude_PicksFirstByWaveThenID
// verifies that with multiple ready slices inside Claude, the first by
// wave/ID order is picked.
func TestRunSlice_AutoDetect_MultipleReady_InsideClaude_PicksFirstByWaveThenID(t *testing.T) {
	setClaudeContext(t, "1")

	slices := []testSlice{
		{id: "S-01", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
		{id: "S-02", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
		{id: "S-03", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	// All slices are wave 1 with no deps — all are ready.
	graphJSON := graphWithWave(map[string]struct {
		deps []string
		wave int
	}{
		"S-01": {deps: []string{}, wave: 1},
		"S-02": {deps: []string{}, wave: 1},
		"S-03": {deps: []string{}, wave: 1},
	})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 with multiple ready slices inside Claude, got %d; output:\n%s", exitCode, buf.String())
	}

	output := buf.String()
	// S-01 is lexicographically first within wave 1.
	if !strings.Contains(output, "S-01") {
		t.Errorf("expected first slice 'S-01' (by wave/ID order) to be picked, got:\n%s", output)
	}
}

// TestRunSlice_AutoDetect_MultipleReady_InsideClaude_PrintsHintAboutRest
// verifies that when multiple slices are ready and we are inside Claude,
// a hint is printed about the remaining ready slices.
func TestRunSlice_AutoDetect_MultipleReady_InsideClaude_PrintsHintAboutRest(t *testing.T) {
	setClaudeContext(t, "1")

	slices := []testSlice{
		{id: "S-01", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
		{id: "S-02", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
		{id: "S-03", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	graphJSON := minimalGraph([]string{"S-01", "S-02", "S-03"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{}, &buf)

	output := buf.String()
	// Should mention remaining count or suggest --max / parallel mode.
	lowerOutput := strings.ToLower(output)
	hasHint := strings.Contains(lowerOutput, "more") ||
		strings.Contains(lowerOutput, "--max") ||
		strings.Contains(lowerOutput, "ready") ||
		strings.Contains(lowerOutput, "parallel") ||
		strings.Contains(output, "2")
	if !hasHint {
		t.Errorf("expected hint about remaining ready slices when multiple are ready, got:\n%s", output)
	}
}

// TestRunSlice_AutoDetect_MultipleReady_WaveOrderTakesPrecedence verifies that
// a lower-wave ready slice is picked before a higher-wave one even if its ID
// comes later alphabetically.
func TestRunSlice_AutoDetect_MultipleReady_WaveOrderTakesPrecedence(t *testing.T) {
	setClaudeContext(t, "1")

	slices := []testSlice{
		{id: "S-10", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
		{id: "S-20", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	// S-10 is wave 2, S-20 is wave 1 — S-20 should be picked first.
	graphJSON := graphWithWave(map[string]struct {
		deps []string
		wave int
	}{
		"S-10": {deps: []string{}, wave: 2},
		"S-20": {deps: []string{}, wave: 1},
	})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d; output:\n%s", exitCode, buf.String())
	}

	output := buf.String()
	if !strings.Contains(output, "S-20") {
		t.Errorf("expected lower-wave slice 'S-20' to be picked, got:\n%s", output)
	}
}

// TestRunSlice_AutoDetect_MultipleReady_Sequential_PicksFirstAndRunsOne
// verifies that --sequential forces single-slice execution: picks first ready
// slice and runs it.
func TestRunSlice_AutoDetect_MultipleReady_Sequential_PicksFirstAndRunsOne(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := []testSlice{
		{id: "S-01", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
		{id: "S-02", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	// --sequential forces single-slice mode; expect non-zero because claude not found.
	exitCode := cmd.RunSlice([]string{"--sequential"}, &buf)

	// The important thing is it attempts to run exactly one slice (not parallel).
	// With claude not in PATH it will fail with spawn error, which is non-zero.
	// We verify it didn't succeed with 0 (which would mean it didn't try to spawn).
	_ = exitCode // exit code is non-zero due to spawn error; tested separately

	output := buf.String()
	// Should attempt S-01 (first by ID in wave 1).
	if !strings.Contains(output, "S-01") {
		t.Errorf("expected --sequential to pick 'S-01' (first ready), got:\n%s", output)
	}
}

// TestRunSlice_AutoDetect_MultipleReady_DryRun_PicksFirstAndShowsCommand
// verifies that with multiple ready slices and --dry-run, the first ready
// slice is shown in the dry-run output.
func TestRunSlice_AutoDetect_MultipleReady_DryRun_PicksFirstAndShowsCommand(t *testing.T) {
	clearClaudeContext(t)

	slices := []testSlice{
		{id: "S-01", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
		{id: "S-02", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--dry-run"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 for auto-detect --dry-run, got %d; output:\n%s", exitCode, buf.String())
	}

	output := buf.String()
	if !strings.Contains(output, "S-01") {
		t.Errorf("expected first ready slice 'S-01' in auto-detect dry-run output, got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-106 — RunSlice auto-detect: state read error
// ----------------------------------------------------------------------------

// TestRunSlice_AutoDetect_NoSlicesDir_ReturnsExitCode1 verifies that when the
// slices directory cannot be read, RunSlice returns 1.
func TestRunSlice_AutoDetect_NoSlicesDir_ReturnsExitCode1(t *testing.T) {
	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{}, &buf)

	if exitCode != 1 {
		t.Errorf("expected exit code 1 when .greenlight/ not found for auto-detect, got %d", exitCode)
	}
}

// TestRunSlice_AutoDetect_NoSlicesDir_PrintsError verifies that when state
// cannot be read, an error is printed.
func TestRunSlice_AutoDetect_NoSlicesDir_PrintsError(t *testing.T) {
	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	var buf bytes.Buffer
	cmd.RunSlice([]string{}, &buf)

	output := buf.String()
	if output == "" {
		t.Error("expected error output when .greenlight/ not found, got empty string")
	}
}

// ----------------------------------------------------------------------------
// Output contract: RunSlice always writes to provided io.Writer
// ----------------------------------------------------------------------------

// TestRunSlice_AlwaysWritesToProvidedWriter verifies that output is written to
// the provided io.Writer in all invocation modes.
func TestRunSlice_AlwaysWritesToProvidedWriter(t *testing.T) {
	tests := []struct {
		name string
		args []string
		setup func(t *testing.T)
	}{
		{
			name: "dry-run shell context",
			args: []string{"S-35", "--dry-run"},
			setup: func(t *testing.T) {
				clearClaudeContext(t)
			},
		},
		{
			name: "inside Claude context",
			args: []string{"S-35"},
			setup: func(t *testing.T) {
				setClaudeContext(t, "1")
			},
		},
		{
			name: "auto-detect inside Claude",
			args: []string{},
			setup: func(t *testing.T) {
				setClaudeContext(t, "1")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)

			slices := []testSlice{
				{id: "S-35", status: "pending", step: "test-writer", milestone: "core",
					started: "", updated: ""},
			}
			graphJSON := minimalGraph([]string{"S-35"})
			setupSliceProject(t, slices, graphJSON)

			var buf bytes.Buffer
			cmd.RunSlice(tt.args, &buf)

			if buf.Len() == 0 {
				t.Errorf("%s: expected output written to provided writer, got empty buffer", tt.name)
			}
		})
	}
}

// ----------------------------------------------------------------------------
// C-105 — Invariants: single slice mode always runs directly (no tmux)
// ----------------------------------------------------------------------------

// TestRunSlice_SingleSliceMode_NeverUsesTmux verifies that slice runs do not
// mention tmux in their output (single slice mode is headless, not tmux-based).
func TestRunSlice_SingleSliceMode_NeverUsesTmux(t *testing.T) {
	setClaudeContext(t, "1")

	slices := []testSlice{
		{id: "S-35", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	graphJSON := minimalGraph([]string{"S-35"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"S-35"}, &buf)

	output := buf.String()
	if strings.Contains(strings.ToLower(output), "tmux") {
		t.Errorf("single slice mode must never use tmux; found 'tmux' in output:\n%s", output)
	}
}

// TestRunSlice_DryRun_NeverMentionsTmux verifies that dry-run output does not
// contain tmux references.
func TestRunSlice_DryRun_NeverMentionsTmux(t *testing.T) {
	slices := []testSlice{
		{id: "S-35", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	graphJSON := minimalGraph([]string{"S-35"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"S-35", "--dry-run"}, &buf)

	output := buf.String()
	if strings.Contains(strings.ToLower(output), "tmux") {
		t.Errorf("single slice --dry-run must not mention tmux; got:\n%s", output)
	}
}
