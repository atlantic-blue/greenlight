package cmd_test

// S-46 Tests: Interactive Commands
// Covers C-115 (RunInit) and C-116 (RunDesign).
//
// Contract C-115 — RunInit:
//   - Detects context via state.DetectContext() ($CLAUDE_CODE env var)
//   - Inside Claude: prints /gl:init skill instructions, returns 0, never spawns
//   - Shell context, claude in PATH: prints "Launching" message, attempts
//     interactive spawn (which will fail without a real binary), returns non-zero
//   - Shell context, claude NOT in PATH: prints install instructions, returns 1
//   - Does NOT require .greenlight/ directory (unlike design)
//   - Args are unused: extra args do not affect behaviour
//   - NEVER uses --dangerously-skip-permissions in interactive mode
//
// Contract C-116 — RunDesign:
//   - Detects context via state.DetectContext() ($CLAUDE_CODE env var)
//   - Inside Claude: prints /gl:design skill instructions, returns 0, never spawns
//   - Shell context, .greenlight/ exists, claude in PATH: prints "Launching" message
//   - Shell context, no .greenlight/ directory: prints "Not a greenlight project", returns 1
//   - Shell context, .greenlight/ exists, claude NOT in PATH: prints install instructions, returns 1
//   - Args are unused: extra args do not affect behaviour
//   - NEVER uses --dangerously-skip-permissions in interactive mode

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
// Helpers specific to interactive command tests
// ----------------------------------------------------------------------------

// setupEmptyProject creates a temp directory with a .greenlight/ subdirectory
// and calls t.Chdir so that RunDesign can locate it from the working directory.
// It returns the project root path.
func setupEmptyProject(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	greenlightDir := filepath.Join(tmpDir, ".greenlight")
	if mkdirError := os.MkdirAll(greenlightDir, 0o755); mkdirError != nil {
		t.Fatalf("setupEmptyProject: failed to create .greenlight dir: %v", mkdirError)
	}
	t.Chdir(tmpDir)
	return tmpDir
}

// setupBareDir creates a temp directory with NO .greenlight/ subdirectory and
// calls t.Chdir. This simulates a directory that is not a greenlight project.
func setupBareDir(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	t.Chdir(tmpDir)
	return tmpDir
}

// ----------------------------------------------------------------------------
// C-115 — RunInit: inside Claude context
// ----------------------------------------------------------------------------

// TestRunInit_InsideClaude_ReturnsExitCode0 verifies that when $CLAUDE_CODE is
// set, RunInit returns 0 and does not attempt to spawn a process.
func TestRunInit_InsideClaude_ReturnsExitCode0(t *testing.T) {
	setClaudeContext(t, "1")
	overrideProcessLookPath(t, claudeNotInPath)

	setupBareDir(t)

	var buf bytes.Buffer
	exitCode := cmd.RunInit([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 inside Claude context for RunInit, got %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunInit_InsideClaude_OutputContainsGlInit verifies that inside Claude
// context, RunInit prints instructions referencing the /gl:init skill.
func TestRunInit_InsideClaude_OutputContainsGlInit(t *testing.T) {
	setClaudeContext(t, "1")

	setupBareDir(t)

	var buf bytes.Buffer
	cmd.RunInit([]string{}, &buf)

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	hasInitRef := strings.Contains(lowerOutput, "gl:init") ||
		strings.Contains(lowerOutput, "/gl:init") ||
		strings.Contains(lowerOutput, "init")
	if !hasInitRef {
		t.Errorf("expected output to reference gl:init skill inside Claude context, got:\n%s", output)
	}
}

// TestRunInit_InsideClaude_OutputDoesNotContainLaunching verifies that inside
// Claude context, RunInit never prints a "Launching" message (no spawn attempt).
func TestRunInit_InsideClaude_OutputDoesNotContainLaunching(t *testing.T) {
	setClaudeContext(t, "1")

	setupBareDir(t)

	var buf bytes.Buffer
	cmd.RunInit([]string{}, &buf)

	output := buf.String()
	if strings.Contains(strings.ToLower(output), "launching") {
		t.Errorf("inside Claude context RunInit must not print 'Launching'; got:\n%s", output)
	}
}

// TestRunInit_InsideClaude_NeverSpawnsProcess verifies the invariant that inside
// Claude context no child Claude process is ever started — proven by confirming
// no error occurs even when LookPath reports claude as absent.
func TestRunInit_InsideClaude_NeverSpawnsProcess(t *testing.T) {
	setClaudeContext(t, "1")
	overrideProcessLookPath(t, claudeNotInPath)

	setupBareDir(t)

	var buf bytes.Buffer
	exitCode := cmd.RunInit([]string{}, &buf)

	// If RunInit tried to spawn claude it would get ErrClaudeNotFound and fail.
	if exitCode != 0 {
		t.Errorf("inside Claude RunInit must never spawn another Claude; got exit code %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunInit_InsideClaude_WritesToProvidedWriter verifies that inside Claude
// context, RunInit writes its output to the provided io.Writer.
func TestRunInit_InsideClaude_WritesToProvidedWriter(t *testing.T) {
	setClaudeContext(t, "1")

	setupBareDir(t)

	var buf bytes.Buffer
	cmd.RunInit([]string{}, &buf)

	if buf.Len() == 0 {
		t.Error("expected output written to provided writer inside Claude context for RunInit, got empty buffer")
	}
}

// ----------------------------------------------------------------------------
// C-115 — RunInit: shell context, claude in PATH
// ----------------------------------------------------------------------------

// TestRunInit_ShellContext_ClaudeInPath_PrintsLaunching verifies that in shell
// context with claude available, RunInit prints a "Launching" message before
// attempting the interactive session.
func TestRunInit_ShellContext_ClaudeInPath_PrintsLaunching(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeInPath)

	setupBareDir(t)

	var buf bytes.Buffer
	cmd.RunInit([]string{}, &buf)

	output := buf.String()
	if !strings.Contains(strings.ToLower(output), "launching") {
		t.Errorf("expected 'Launching' message in shell context for RunInit, got:\n%s", output)
	}
}

// TestRunInit_ShellContext_ClaudeInPath_ReturnsNonZeroOnSpawnFailure verifies
// that in shell context where the fake claude binary cannot actually run, RunInit
// returns a non-zero exit code (the spawn attempt fails at os.exec level).
func TestRunInit_ShellContext_ClaudeInPath_ReturnsNonZeroOnSpawnFailure(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeInPath)

	setupBareDir(t)

	var buf bytes.Buffer
	exitCode := cmd.RunInit([]string{}, &buf)

	// LookPath says claude exists but the binary at /usr/local/bin/claude is
	// not a real executable in the test environment, so exec.Cmd.Run() will fail.
	if exitCode == 0 {
		t.Errorf("expected non-zero exit code when claude binary cannot actually run, got 0; output:\n%s", buf.String())
	}
}

// TestRunInit_ShellContext_ClaudeInPath_WritesToProvidedWriter verifies that
// RunInit writes output to the provided io.Writer even when the spawn fails.
func TestRunInit_ShellContext_ClaudeInPath_WritesToProvidedWriter(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeInPath)

	setupBareDir(t)

	var buf bytes.Buffer
	cmd.RunInit([]string{}, &buf)

	if buf.Len() == 0 {
		t.Error("expected output written to provided writer in shell context for RunInit, got empty buffer")
	}
}

// ----------------------------------------------------------------------------
// C-115 — RunInit: shell context, claude NOT in PATH
// ----------------------------------------------------------------------------

// TestRunInit_ShellContext_ClaudeNotFound_ReturnsExitCode1 verifies that in
// shell context, when the claude binary is absent from PATH, RunInit returns 1.
func TestRunInit_ShellContext_ClaudeNotFound_ReturnsExitCode1(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeNotInPath)

	setupBareDir(t)

	var buf bytes.Buffer
	exitCode := cmd.RunInit([]string{}, &buf)

	if exitCode != 1 {
		t.Errorf("expected exit code 1 when claude not in PATH for RunInit, got %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunInit_ShellContext_ClaudeNotFound_PrintsInstallInstructions verifies
// that when claude is absent from PATH, RunInit prints install instructions
// mentioning "claude" or "install".
func TestRunInit_ShellContext_ClaudeNotFound_PrintsInstallInstructions(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeNotInPath)

	setupBareDir(t)

	var buf bytes.Buffer
	cmd.RunInit([]string{}, &buf)

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	hasInstallRef := strings.Contains(lowerOutput, "install") ||
		strings.Contains(lowerOutput, "claude") ||
		strings.Contains(lowerOutput, "not found")
	if !hasInstallRef {
		t.Errorf("expected install instructions when claude not in PATH for RunInit, got:\n%s", output)
	}
}

// TestRunInit_ShellContext_ClaudeNotFound_PrintsClaudeReference verifies that
// the error output for missing claude binary specifically mentions "claude".
func TestRunInit_ShellContext_ClaudeNotFound_PrintsClaudeReference(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeNotInPath)

	setupBareDir(t)

	var buf bytes.Buffer
	cmd.RunInit([]string{}, &buf)

	output := buf.String()
	if !strings.Contains(strings.ToLower(output), "claude") {
		t.Errorf("expected error mentioning 'claude' when not in PATH for RunInit, got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-115 — RunInit: does NOT require .greenlight/ directory
// ----------------------------------------------------------------------------

// TestRunInit_ShellContext_NoGreenlightDir_ClaudeNotFound_ReturnsExitCode1
// verifies that RunInit in a directory without .greenlight/ still fails only
// because claude is absent — not because of a missing project directory.
func TestRunInit_ShellContext_NoGreenlightDir_ClaudeNotFound_ReturnsExitCode1(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeNotInPath)

	bareDir := t.TempDir()
	t.Chdir(bareDir)

	var buf bytes.Buffer
	exitCode := cmd.RunInit([]string{}, &buf)

	// Exit code 1 is expected because claude is not found, NOT because of missing .greenlight/.
	if exitCode != 1 {
		t.Errorf("expected exit code 1 (claude not found) even without .greenlight/, got %d; output:\n%s", exitCode, buf.String())
	}

	output := buf.String()
	// The error must mention claude (not a missing project error).
	if !strings.Contains(strings.ToLower(output), "claude") {
		t.Errorf("RunInit error in bare dir must mention 'claude', not a project-dir error; got:\n%s", output)
	}
}

// TestRunInit_InsideClaude_NoGreenlightDir_ReturnsExitCode0 verifies that inside
// Claude context, RunInit returns 0 even when .greenlight/ does not exist.
// Init is meant to CREATE the project structure, so it must not require it.
func TestRunInit_InsideClaude_NoGreenlightDir_ReturnsExitCode0(t *testing.T) {
	setClaudeContext(t, "1")

	bareDir := t.TempDir()
	t.Chdir(bareDir)

	var buf bytes.Buffer
	exitCode := cmd.RunInit([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("RunInit inside Claude must succeed even without .greenlight/ dir, got exit code %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunInit_ShellContext_NoGreenlightDir_ClaudeInPath_PrintsLaunching verifies
// that RunInit with claude in PATH (but no .greenlight/) still prints "Launching"
// — it does not refuse to run because the project directory is absent.
func TestRunInit_ShellContext_NoGreenlightDir_ClaudeInPath_PrintsLaunching(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeInPath)

	bareDir := t.TempDir()
	t.Chdir(bareDir)

	var buf bytes.Buffer
	cmd.RunInit([]string{}, &buf)

	output := buf.String()
	// Must attempt launch (print "Launching"), not refuse due to missing .greenlight/.
	if !strings.Contains(strings.ToLower(output), "launching") {
		t.Errorf("RunInit must attempt launch without .greenlight/ dir; expected 'Launching', got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-115 — RunInit: args are unused
// ----------------------------------------------------------------------------

// TestRunInit_InsideClaude_ExtraArgsIgnored_ReturnsExitCode0 verifies that
// extra positional arguments do not affect the behaviour of RunInit.
func TestRunInit_InsideClaude_ExtraArgsIgnored_ReturnsExitCode0(t *testing.T) {
	setClaudeContext(t, "1")

	setupBareDir(t)

	var buf bytes.Buffer
	exitCode := cmd.RunInit([]string{"--some-flag", "extra-arg", "another"}, &buf)

	if exitCode != 0 {
		t.Errorf("extra args must be ignored by RunInit inside Claude, got exit code %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunInit_ShellContext_ExtraArgsIgnored_ClaudeNotFound_ReturnsExitCode1
// verifies that extra args do not change the exit code or error path.
func TestRunInit_ShellContext_ExtraArgsIgnored_ClaudeNotFound_ReturnsExitCode1(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeNotInPath)

	setupBareDir(t)

	var buf bytes.Buffer
	exitCode := cmd.RunInit([]string{"--verbose", "foo", "bar"}, &buf)

	if exitCode != 1 {
		t.Errorf("extra args must be ignored by RunInit; expected exit code 1 (claude not found), got %d; output:\n%s", exitCode, buf.String())
	}
}

// ----------------------------------------------------------------------------
// C-115 — RunInit: invariant — no --dangerously-skip-permissions
// ----------------------------------------------------------------------------

// TestRunInit_InsideClaude_NoDangerousSkipPermissions verifies that the output
// for inside-Claude mode never references --dangerously-skip-permissions.
func TestRunInit_InsideClaude_NoDangerousSkipPermissions(t *testing.T) {
	setClaudeContext(t, "1")

	setupBareDir(t)

	var buf bytes.Buffer
	cmd.RunInit([]string{}, &buf)

	output := buf.String()
	if strings.Contains(output, "--dangerously-skip-permissions") {
		t.Errorf("RunInit inside Claude must never reference --dangerously-skip-permissions; got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-116 — RunDesign: inside Claude context
// ----------------------------------------------------------------------------

// TestRunDesign_InsideClaude_ReturnsExitCode0 verifies that when $CLAUDE_CODE is
// set, RunDesign returns 0 and does not attempt to spawn a process.
func TestRunDesign_InsideClaude_ReturnsExitCode0(t *testing.T) {
	setClaudeContext(t, "1")
	overrideProcessLookPath(t, claudeNotInPath)

	setupBareDir(t)

	var buf bytes.Buffer
	exitCode := cmd.RunDesign([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 inside Claude context for RunDesign, got %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunDesign_InsideClaude_OutputContainsGlDesign verifies that inside Claude
// context, RunDesign prints instructions referencing the /gl:design skill.
func TestRunDesign_InsideClaude_OutputContainsGlDesign(t *testing.T) {
	setClaudeContext(t, "1")

	setupBareDir(t)

	var buf bytes.Buffer
	cmd.RunDesign([]string{}, &buf)

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	hasDesignRef := strings.Contains(lowerOutput, "gl:design") ||
		strings.Contains(lowerOutput, "/gl:design") ||
		strings.Contains(lowerOutput, "design")
	if !hasDesignRef {
		t.Errorf("expected output to reference gl:design skill inside Claude context, got:\n%s", output)
	}
}

// TestRunDesign_InsideClaude_OutputDoesNotContainLaunching verifies that inside
// Claude context, RunDesign never prints a "Launching" message (no spawn attempt).
func TestRunDesign_InsideClaude_OutputDoesNotContainLaunching(t *testing.T) {
	setClaudeContext(t, "1")

	setupBareDir(t)

	var buf bytes.Buffer
	cmd.RunDesign([]string{}, &buf)

	output := buf.String()
	if strings.Contains(strings.ToLower(output), "launching") {
		t.Errorf("inside Claude context RunDesign must not print 'Launching'; got:\n%s", output)
	}
}

// TestRunDesign_InsideClaude_NeverSpawnsProcess verifies the invariant that
// inside Claude context no child Claude process is ever started.
func TestRunDesign_InsideClaude_NeverSpawnsProcess(t *testing.T) {
	setClaudeContext(t, "1")
	overrideProcessLookPath(t, claudeNotInPath)

	setupBareDir(t)

	var buf bytes.Buffer
	exitCode := cmd.RunDesign([]string{}, &buf)

	// If RunDesign tried to spawn claude it would get ErrClaudeNotFound and fail.
	if exitCode != 0 {
		t.Errorf("inside Claude RunDesign must never spawn another Claude; got exit code %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunDesign_InsideClaude_WritesToProvidedWriter verifies that inside Claude
// context, RunDesign writes its output to the provided io.Writer.
func TestRunDesign_InsideClaude_WritesToProvidedWriter(t *testing.T) {
	setClaudeContext(t, "1")

	setupBareDir(t)

	var buf bytes.Buffer
	cmd.RunDesign([]string{}, &buf)

	if buf.Len() == 0 {
		t.Error("expected output written to provided writer inside Claude context for RunDesign, got empty buffer")
	}
}

// ----------------------------------------------------------------------------
// C-116 — RunDesign: shell context, .greenlight/ exists, claude in PATH
// ----------------------------------------------------------------------------

// TestRunDesign_ShellContext_GreenlightExists_ClaudeInPath_PrintsLaunching
// verifies that in shell context with .greenlight/ present and claude available,
// RunDesign prints a "Launching" message before attempting the interactive session.
func TestRunDesign_ShellContext_GreenlightExists_ClaudeInPath_PrintsLaunching(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeInPath)

	setupEmptyProject(t)

	var buf bytes.Buffer
	cmd.RunDesign([]string{}, &buf)

	output := buf.String()
	if !strings.Contains(strings.ToLower(output), "launching") {
		t.Errorf("expected 'Launching' message in shell context for RunDesign, got:\n%s", output)
	}
}

// TestRunDesign_ShellContext_GreenlightExists_ClaudeInPath_ReturnsNonZeroOnSpawnFailure
// verifies that in shell context where the fake claude binary cannot actually run,
// RunDesign returns a non-zero exit code.
func TestRunDesign_ShellContext_GreenlightExists_ClaudeInPath_ReturnsNonZeroOnSpawnFailure(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeInPath)

	setupEmptyProject(t)

	var buf bytes.Buffer
	exitCode := cmd.RunDesign([]string{}, &buf)

	// LookPath says claude exists but the binary cannot actually run in tests.
	if exitCode == 0 {
		t.Errorf("expected non-zero exit code when claude binary cannot actually run for RunDesign, got 0; output:\n%s", buf.String())
	}
}

// TestRunDesign_ShellContext_GreenlightExists_ClaudeInPath_WritesToProvidedWriter
// verifies that RunDesign writes output to the provided io.Writer even when spawn fails.
func TestRunDesign_ShellContext_GreenlightExists_ClaudeInPath_WritesToProvidedWriter(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeInPath)

	setupEmptyProject(t)

	var buf bytes.Buffer
	cmd.RunDesign([]string{}, &buf)

	if buf.Len() == 0 {
		t.Error("expected output written to provided writer in shell context for RunDesign, got empty buffer")
	}
}

// ----------------------------------------------------------------------------
// C-116 — RunDesign: shell context, no .greenlight/ directory
// ----------------------------------------------------------------------------

// TestRunDesign_ShellContext_NoGreenlightDir_ReturnsExitCode1 verifies that in
// shell context without .greenlight/, RunDesign prints an error and returns 1.
func TestRunDesign_ShellContext_NoGreenlightDir_ReturnsExitCode1(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeInPath)

	setupBareDir(t)

	var buf bytes.Buffer
	exitCode := cmd.RunDesign([]string{}, &buf)

	if exitCode != 1 {
		t.Errorf("expected exit code 1 when .greenlight/ missing for RunDesign, got %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunDesign_ShellContext_NoGreenlightDir_PrintsNotAGreenlightProject verifies
// that the error message for a missing .greenlight/ directory references
// "greenlight project" or "gl init".
func TestRunDesign_ShellContext_NoGreenlightDir_PrintsNotAGreenlightProject(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeInPath)

	setupBareDir(t)

	var buf bytes.Buffer
	cmd.RunDesign([]string{}, &buf)

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	hasProjectError := strings.Contains(lowerOutput, "not a greenlight project") ||
		strings.Contains(lowerOutput, "greenlight project") ||
		strings.Contains(lowerOutput, "gl init") ||
		strings.Contains(lowerOutput, "run 'gl init'")
	if !hasProjectError {
		t.Errorf("expected 'Not a greenlight project' error for missing .greenlight/, got:\n%s", output)
	}
}

// TestRunDesign_ShellContext_NoGreenlightDir_PrintsInitHint verifies that the
// error output for a missing .greenlight/ suggests running 'gl init'.
func TestRunDesign_ShellContext_NoGreenlightDir_PrintsInitHint(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeInPath)

	setupBareDir(t)

	var buf bytes.Buffer
	cmd.RunDesign([]string{}, &buf)

	output := buf.String()
	if !strings.Contains(strings.ToLower(output), "init") {
		t.Errorf("expected hint to run 'gl init' when .greenlight/ missing for RunDesign, got:\n%s", output)
	}
}

// TestRunDesign_ShellContext_NoGreenlightDir_DoesNotPrintLaunching verifies that
// RunDesign does not attempt to launch claude when the project directory is missing.
func TestRunDesign_ShellContext_NoGreenlightDir_DoesNotPrintLaunching(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeInPath)

	setupBareDir(t)

	var buf bytes.Buffer
	cmd.RunDesign([]string{}, &buf)

	output := buf.String()
	if strings.Contains(strings.ToLower(output), "launching") {
		t.Errorf("RunDesign must not print 'Launching' when .greenlight/ is missing; got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-116 — RunDesign: shell context, .greenlight/ exists, claude NOT in PATH
// ----------------------------------------------------------------------------

// TestRunDesign_ShellContext_GreenlightExists_ClaudeNotFound_ReturnsExitCode1
// verifies that in shell context with .greenlight/ present but claude absent
// from PATH, RunDesign prints install instructions and returns 1.
func TestRunDesign_ShellContext_GreenlightExists_ClaudeNotFound_ReturnsExitCode1(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeNotInPath)

	setupEmptyProject(t)

	var buf bytes.Buffer
	exitCode := cmd.RunDesign([]string{}, &buf)

	if exitCode != 1 {
		t.Errorf("expected exit code 1 when claude not in PATH for RunDesign, got %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunDesign_ShellContext_GreenlightExists_ClaudeNotFound_PrintsInstallInstructions
// verifies that when claude is absent from PATH, RunDesign prints install
// instructions that mention "claude" or "install".
func TestRunDesign_ShellContext_GreenlightExists_ClaudeNotFound_PrintsInstallInstructions(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeNotInPath)

	setupEmptyProject(t)

	var buf bytes.Buffer
	cmd.RunDesign([]string{}, &buf)

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	hasInstallRef := strings.Contains(lowerOutput, "install") ||
		strings.Contains(lowerOutput, "claude") ||
		strings.Contains(lowerOutput, "not found")
	if !hasInstallRef {
		t.Errorf("expected install instructions when claude not in PATH for RunDesign, got:\n%s", output)
	}
}

// TestRunDesign_ShellContext_GreenlightExists_ClaudeNotFound_PrintsClaudeReference
// verifies that the error output for a missing claude binary specifically mentions "claude".
func TestRunDesign_ShellContext_GreenlightExists_ClaudeNotFound_PrintsClaudeReference(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeNotInPath)

	setupEmptyProject(t)

	var buf bytes.Buffer
	cmd.RunDesign([]string{}, &buf)

	output := buf.String()
	if !strings.Contains(strings.ToLower(output), "claude") {
		t.Errorf("expected error mentioning 'claude' when not in PATH for RunDesign, got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-116 — RunDesign: args are unused
// ----------------------------------------------------------------------------

// TestRunDesign_InsideClaude_ExtraArgsIgnored_ReturnsExitCode0 verifies that
// extra positional arguments do not affect the behaviour of RunDesign.
func TestRunDesign_InsideClaude_ExtraArgsIgnored_ReturnsExitCode0(t *testing.T) {
	setClaudeContext(t, "1")

	setupBareDir(t)

	var buf bytes.Buffer
	exitCode := cmd.RunDesign([]string{"--some-flag", "extra-arg", "another"}, &buf)

	if exitCode != 0 {
		t.Errorf("extra args must be ignored by RunDesign inside Claude, got exit code %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunDesign_ShellContext_ExtraArgsIgnored_NoGreenlightDir_ReturnsExitCode1
// verifies that extra args do not change the exit code when .greenlight/ is absent.
func TestRunDesign_ShellContext_ExtraArgsIgnored_NoGreenlightDir_ReturnsExitCode1(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeInPath)

	setupBareDir(t)

	var buf bytes.Buffer
	exitCode := cmd.RunDesign([]string{"--verbose", "foo", "bar"}, &buf)

	if exitCode != 1 {
		t.Errorf("extra args must be ignored by RunDesign; expected exit code 1 (no .greenlight/), got %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunDesign_ShellContext_ExtraArgsIgnored_ClaudeNotFound_ReturnsExitCode1
// verifies that extra args do not affect the error path when claude is absent.
func TestRunDesign_ShellContext_ExtraArgsIgnored_ClaudeNotFound_ReturnsExitCode1(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeNotInPath)

	setupEmptyProject(t)

	var buf bytes.Buffer
	exitCode := cmd.RunDesign([]string{"--verbose", "foo", "bar"}, &buf)

	if exitCode != 1 {
		t.Errorf("extra args must be ignored by RunDesign; expected exit code 1 (claude not found), got %d; output:\n%s", exitCode, buf.String())
	}
}

// ----------------------------------------------------------------------------
// C-116 — RunDesign: invariant — no --dangerously-skip-permissions
// ----------------------------------------------------------------------------

// TestRunDesign_InsideClaude_NoDangerousSkipPermissions verifies that the output
// for inside-Claude mode never references --dangerously-skip-permissions.
func TestRunDesign_InsideClaude_NoDangerousSkipPermissions(t *testing.T) {
	setClaudeContext(t, "1")

	setupBareDir(t)

	var buf bytes.Buffer
	cmd.RunDesign([]string{}, &buf)

	output := buf.String()
	if strings.Contains(output, "--dangerously-skip-permissions") {
		t.Errorf("RunDesign inside Claude must never reference --dangerously-skip-permissions; got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// Invariants: inside Claude context never spawns (both commands)
// ----------------------------------------------------------------------------

// TestRunInit_InsideClaude_OutputNeverContainsSpawn verifies that inside Claude
// context, RunInit output contains neither "spawn" nor "Launching".
func TestRunInit_InsideClaude_OutputNeverContainsSpawn(t *testing.T) {
	setClaudeContext(t, "1")

	setupBareDir(t)

	var buf bytes.Buffer
	cmd.RunInit([]string{}, &buf)

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	if strings.Contains(lowerOutput, "spawn") {
		t.Errorf("inside Claude RunInit output must not contain 'spawn'; got:\n%s", output)
	}
	if strings.Contains(lowerOutput, "launching") {
		t.Errorf("inside Claude RunInit output must not contain 'Launching'; got:\n%s", output)
	}
}

// TestRunDesign_InsideClaude_OutputNeverContainsSpawn verifies that inside Claude
// context, RunDesign output contains neither "spawn" nor "Launching".
func TestRunDesign_InsideClaude_OutputNeverContainsSpawn(t *testing.T) {
	setClaudeContext(t, "1")

	setupBareDir(t)

	var buf bytes.Buffer
	cmd.RunDesign([]string{}, &buf)

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	if strings.Contains(lowerOutput, "spawn") {
		t.Errorf("inside Claude RunDesign output must not contain 'spawn'; got:\n%s", output)
	}
	if strings.Contains(lowerOutput, "launching") {
		t.Errorf("inside Claude RunDesign output must not contain 'Launching'; got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// Invariants: neither command reads slice state files
// ----------------------------------------------------------------------------

// TestRunInit_InsideClaude_DoesNotRequireSlicesDir verifies that RunInit does
// not read .greenlight/slices/ — no slice state is needed for init.
func TestRunInit_InsideClaude_DoesNotRequireSlicesDir(t *testing.T) {
	setClaudeContext(t, "1")

	// No slices/ subdirectory, no GRAPH.json — pure empty project root.
	tmpDir := t.TempDir()
	t.Chdir(tmpDir)

	var buf bytes.Buffer
	exitCode := cmd.RunInit([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("RunInit inside Claude must not require slices dir, got exit code %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunDesign_InsideClaude_DoesNotRequireSlicesDir verifies that RunDesign in
// inside-Claude mode does not require .greenlight/slices/ to exist — it only
// routes to the /gl:design skill.
func TestRunDesign_InsideClaude_DoesNotRequireSlicesDir(t *testing.T) {
	setClaudeContext(t, "1")

	// .greenlight/ exists but slices/ does not — pure routing should still work.
	tmpDir := t.TempDir()
	greenlightDir := filepath.Join(tmpDir, ".greenlight")
	if mkdirError := os.MkdirAll(greenlightDir, 0o755); mkdirError != nil {
		t.Fatalf("failed to create .greenlight dir: %v", mkdirError)
	}
	t.Chdir(tmpDir)

	var buf bytes.Buffer
	exitCode := cmd.RunDesign([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("RunDesign inside Claude must not require slices dir, got exit code %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunInit_ShellContext_NoGreenlightDir_OutputDoesNotMentionSlices verifies
// that when RunInit fails because claude is absent, it does not mention slice
// state — it is purely a routing + spawn command.
func TestRunInit_ShellContext_NoGreenlightDir_OutputDoesNotMentionSlices(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeNotInPath)

	bareDir := t.TempDir()
	t.Chdir(bareDir)

	var buf bytes.Buffer
	cmd.RunInit([]string{}, &buf)

	output := buf.String()
	// The output must not mention slices or state files — init is a pure spawn
	// command that does not read project state.
	if strings.Contains(strings.ToLower(output), "slices") {
		t.Errorf("RunInit must not reference slice state files in its output; got:\n%s", output)
	}
}

// TestRunDesign_ShellContext_GreenlightExists_ClaudeNotFound_OutputDoesNotMentionSlices
// verifies that when RunDesign fails because claude is absent, it does not mention
// slice state — it is purely a routing + spawn command.
func TestRunDesign_ShellContext_GreenlightExists_ClaudeNotFound_OutputDoesNotMentionSlices(t *testing.T) {
	clearClaudeContext(t)
	overrideProcessLookPath(t, claudeNotInPath)

	setupEmptyProject(t)

	var buf bytes.Buffer
	cmd.RunDesign([]string{}, &buf)

	output := buf.String()
	// The output must not mention slices or state files — design spawn is a
	// pure routing command, it does not read project slice state.
	if strings.Contains(strings.ToLower(output), "slices") {
		t.Errorf("RunDesign must not reference slice state files in its output; got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// Invariants: process.BuildInteractiveCmd never includes --dangerously-skip-permissions
// ----------------------------------------------------------------------------

// TestBuildInteractiveCmd_StripsDesignDangerousFlag verifies that
// process.BuildInteractiveCmd always strips --dangerously-skip-permissions even
// when explicitly passed in flags. This is the security invariant for both
// RunInit and RunDesign interactive sessions.
func TestBuildInteractiveCmd_StripsDesignDangerousFlag(t *testing.T) {
	overrideProcessLookPath(t, claudeInPath)

	command, buildError := process.BuildInteractiveCmd(process.InteractiveOptions{
		Prompt: "/gl:design",
		Flags:  []string{"--dangerously-skip-permissions", "--max-turns", "10"},
	})

	if buildError != nil {
		t.Fatalf("BuildInteractiveCmd returned unexpected error: %v", buildError)
	}

	for _, arg := range command.Args {
		if arg == "--dangerously-skip-permissions" {
			t.Errorf("BuildInteractiveCmd must strip --dangerously-skip-permissions; found in args: %v", command.Args)
		}
	}
}

// TestBuildInteractiveCmd_StripsInitDangerousFlag mirrors the above test for
// the /gl:init prompt, ensuring the invariant holds for init sessions too.
func TestBuildInteractiveCmd_StripsInitDangerousFlag(t *testing.T) {
	overrideProcessLookPath(t, claudeInPath)

	command, buildError := process.BuildInteractiveCmd(process.InteractiveOptions{
		Prompt: "/gl:init",
		Flags:  []string{"--dangerously-skip-permissions"},
	})

	if buildError != nil {
		t.Fatalf("BuildInteractiveCmd returned unexpected error: %v", buildError)
	}

	for _, arg := range command.Args {
		if arg == "--dangerously-skip-permissions" {
			t.Errorf("BuildInteractiveCmd must strip --dangerously-skip-permissions for init sessions; found in args: %v", command.Args)
		}
	}
}

// TestBuildInteractiveCmd_IncludesPromptWhenNonEmpty verifies that when a
// non-empty prompt is provided, it appears in the command arguments as -p <prompt>.
func TestBuildInteractiveCmd_IncludesPromptWhenNonEmpty(t *testing.T) {
	overrideProcessLookPath(t, claudeInPath)

	command, buildError := process.BuildInteractiveCmd(process.InteractiveOptions{
		Prompt: "/gl:init",
	})

	if buildError != nil {
		t.Fatalf("BuildInteractiveCmd returned unexpected error: %v", buildError)
	}

	args := command.Args
	foundPromptFlag := false
	for index, arg := range args {
		if arg == "-p" && index+1 < len(args) && args[index+1] == "/gl:init" {
			foundPromptFlag = true
			break
		}
	}

	if !foundPromptFlag {
		t.Errorf("expected '-p /gl:init' in BuildInteractiveCmd args, got: %v", args)
	}
}

// TestBuildInteractiveCmd_ClaudeNotFound_ReturnsError verifies that
// BuildInteractiveCmd returns ErrClaudeNotFound when claude is absent from PATH.
func TestBuildInteractiveCmd_ClaudeNotFound_ReturnsError(t *testing.T) {
	overrideProcessLookPath(t, claudeNotInPath)

	_, buildError := process.BuildInteractiveCmd(process.InteractiveOptions{
		Prompt: "/gl:init",
	})

	if buildError == nil {
		t.Error("expected error from BuildInteractiveCmd when claude not in PATH, got nil")
	}
}
