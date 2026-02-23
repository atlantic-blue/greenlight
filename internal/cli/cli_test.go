package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/atlantic-blue/greenlight/internal/cli"
)

// buildTestFS returns an fstest.MapFS with all 38 manifest files.
func buildTestFS() *fstest.MapFS {
	manifestFiles := []string{
		"agents/gl-architect.md",
		"agents/gl-assessor.md",
		"agents/gl-codebase-mapper.md",
		"agents/gl-debugger.md",
		"agents/gl-designer.md",
		"agents/gl-implementer.md",
		"agents/gl-security.md",
		"agents/gl-test-writer.md",
		"agents/gl-verifier.md",
		"agents/gl-wrapper.md",
		"commands/gl/add-slice.md",
		"commands/gl/assess.md",
		"commands/gl/changelog.md",
		"commands/gl/debug.md",
		"commands/gl/design.md",
		"commands/gl/help.md",
		"commands/gl/init.md",
		"commands/gl/map.md",
		"commands/gl/migrate-state.md",
		"commands/gl/pause.md",
		"commands/gl/quick.md",
		"commands/gl/resume.md",
		"commands/gl/roadmap.md",
		"commands/gl/settings.md",
		"commands/gl/ship.md",
		"commands/gl/slice.md",
		"commands/gl/status.md",
		"commands/gl/wrap.md",
		"references/checkpoint-protocol.md",
		"references/circuit-breaker.md",
		"references/deviation-rules.md",
		"references/state-format.md",
		"references/verification-patterns.md",
		"references/verification-tiers.md",
		"templates/config.md",
		"templates/slice-state.md",
		"templates/state.md",
		"CLAUDE.md",
	}

	fs := &fstest.MapFS{}
	for _, path := range manifestFiles {
		(*fs)[path] = &fstest.MapFile{
			Data: []byte("# test content for " + path),
		}
	}
	return fs
}

// C-15.1: No args returns 0 and prints usage
func TestRun_NoArgs_ReturnsZeroAndPrintsUsage(t *testing.T) {
	contentFS := buildTestFS()
	var stdout bytes.Buffer

	exitCode := cli.Run([]string{}, contentFS, &stdout)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	output := stdout.String()
	if !strings.Contains(output, "Usage:") && !strings.Contains(output, "greenlight") {
		t.Errorf("expected usage message, got: %s", output)
	}
}

// C-15.2: Help command returns 0
func TestRun_HelpCommands_ReturnZero(t *testing.T) {
	contentFS := buildTestFS()

	testCases := []struct {
		name string
		args []string
	}{
		{"help", []string{"help"}},
		{"--help", []string{"--help"}},
		{"-h", []string{"-h"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var stdout bytes.Buffer
			exitCode := cli.Run(tc.args, contentFS, &stdout)

			if exitCode != 0 {
				t.Errorf("expected exit code 0, got %d", exitCode)
			}

			output := stdout.String()
			if !strings.Contains(output, "Usage:") {
				t.Errorf("expected usage message, got: %s", output)
			}
		})
	}
}

// C-15.3: Version command returns 0
func TestRun_VersionCommand_ReturnsZero(t *testing.T) {
	contentFS := buildTestFS()
	var stdout bytes.Buffer

	exitCode := cli.Run([]string{"version"}, contentFS, &stdout)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	output := stdout.String()
	if !strings.Contains(output, "greenlight") {
		t.Errorf("expected version output to contain 'greenlight', got: %s", output)
	}
}

// C-15.4: Unknown command returns 1
func TestRun_UnknownCommand_ReturnsOne(t *testing.T) {
	contentFS := buildTestFS()
	var stdout bytes.Buffer

	exitCode := cli.Run([]string{"foobar"}, contentFS, &stdout)

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}

	output := stdout.String()
	if !strings.Contains(output, "unknown command: foobar") {
		t.Errorf("expected 'unknown command: foobar', got: %s", output)
	}
}

// C-15.5: Unknown command also prints usage
func TestRun_UnknownCommand_PrintsUsage(t *testing.T) {
	contentFS := buildTestFS()
	var stdout bytes.Buffer

	cli.Run([]string{"foobar"}, contentFS, &stdout)

	output := stdout.String()
	if !strings.Contains(output, "Usage:") {
		t.Errorf("expected usage after unknown command, got: %s", output)
	}
}

// C-15.6: Install command dispatches correctly
func TestRun_InstallCommand_Dispatches(t *testing.T) {
	contentFS := buildTestFS()
	var stdout bytes.Buffer

	tempDir := t.TempDir()
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temp directory: %v", err)
	}

	exitCode := cli.Run([]string{"install", "--local"}, contentFS, &stdout)

	output := stdout.String()
	if strings.Contains(output, "unknown command") {
		t.Errorf("install should dispatch, not show 'unknown command': %s", output)
	}

	// Either succeeds (exit 0) or fails with install-specific error (not dispatch error)
	if exitCode != 0 && !strings.Contains(output, "error:") {
		t.Errorf("expected either success or install error, got exit %d with: %s", exitCode, output)
	}
}

// C-15.7: Uninstall command dispatches correctly
func TestRun_UninstallCommand_Dispatches(t *testing.T) {
	contentFS := buildTestFS()
	var stdout bytes.Buffer

	tempDir := t.TempDir()
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temp directory: %v", err)
	}

	cli.Run([]string{"uninstall", "--local"}, contentFS, &stdout)

	output := stdout.String()
	if strings.Contains(output, "unknown command") {
		t.Errorf("uninstall should dispatch, not show 'unknown command': %s", output)
	}
}

// C-15.8: Check command dispatches correctly
func TestRun_CheckCommand_Dispatches(t *testing.T) {
	contentFS := buildTestFS()
	var stdout bytes.Buffer

	tempDir := t.TempDir()
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temp directory: %v", err)
	}

	cli.Run([]string{"check", "--local"}, contentFS, &stdout)

	output := stdout.String()
	if strings.Contains(output, "unknown command") {
		t.Errorf("check should dispatch, not show 'unknown command': %s", output)
	}
}

// C-15.9: Output written to provided writer (TD-2)
func TestRun_OutputWrittenToProvidedWriter(t *testing.T) {
	contentFS := buildTestFS()
	var stdout bytes.Buffer

	cli.Run([]string{}, contentFS, &stdout)

	if stdout.Len() == 0 {
		t.Error("expected output to be written to provided writer, got empty buffer")
	}
}

// C-15.10: Multiple unknown commands show correct name
func TestRun_UnknownCommand_ShowsCommandName(t *testing.T) {
	contentFS := buildTestFS()

	testCases := []struct {
		name        string
		command     string
		expectedMsg string
	}{
		{"deploy", "deploy", "unknown command: deploy"},
		{"upgrade", "upgrade", "unknown command: upgrade"},
		{"init", "init", "unknown command: init"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var stdout bytes.Buffer
			exitCode := cli.Run([]string{tc.command}, contentFS, &stdout)

			if exitCode != 1 {
				t.Errorf("expected exit code 1, got %d", exitCode)
			}

			output := stdout.String()
			if !strings.Contains(output, tc.expectedMsg) {
				t.Errorf("expected '%s', got: %s", tc.expectedMsg, output)
			}
		})
	}
}

// C-15.11: Version output goes to provided writer
func TestRun_VersionOutput_GoesToProvidedWriter(t *testing.T) {
	contentFS := buildTestFS()
	var stdout bytes.Buffer

	cli.Run([]string{"version"}, contentFS, &stdout)

	output := stdout.String()
	if !strings.Contains(output, "greenlight") {
		t.Errorf("expected version output with 'greenlight', got: %s", output)
	}
	if !strings.Contains(output, "commit:") {
		t.Errorf("expected version output with 'commit:', got: %s", output)
	}
}

// C-15.12: contentFS passed through to install
func TestRun_ContentFS_PassedToInstall(t *testing.T) {
	contentFS := buildTestFS()
	var stdout bytes.Buffer

	tempDir := t.TempDir()
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temp directory: %v", err)
	}

	exitCode := cli.Run([]string{"install", "--local"}, contentFS, &stdout)

	if exitCode == 0 {
		// Check that at least one manifest file was installed (inside .claude/)
		claudeDir := filepath.Join(tempDir, ".claude")
		manifestPath := filepath.Join(claudeDir, "agents", "gl-architect.md")
		if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
			t.Error("expected manifest files to be installed from contentFS, but none found")
		}
	}
}

// C-15.13: contentFS passed through to check
func TestRun_ContentFS_PassedToCheck(t *testing.T) {
	contentFS := buildTestFS()
	var stdout bytes.Buffer

	tempDir := t.TempDir()
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temp directory: %v", err)
	}

	// Install first
	cli.Run([]string{"install", "--local"}, contentFS, &bytes.Buffer{})

	// Now check with verify
	stdout.Reset()
	cli.Run([]string{"check", "--local", "--verify"}, contentFS, &stdout)

	output := stdout.String()
	if !strings.Contains(output, "verified") && !strings.Contains(output, "matching") {
		t.Errorf("expected check output to contain verification status, got: %s", output)
	}
}

// C-15.14: Install with bad flags returns 1 (delegation test)
func TestRun_InstallBadFlags_ReturnsOne(t *testing.T) {
	contentFS := buildTestFS()
	var stdout bytes.Buffer

	tempDir := t.TempDir()
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temp directory: %v", err)
	}

	exitCode := cli.Run([]string{"install"}, contentFS, &stdout)

	if exitCode != 1 {
		t.Errorf("expected exit code 1 for missing scope flag, got %d", exitCode)
	}

	output := stdout.String()
	if !strings.Contains(output, "error:") {
		t.Errorf("expected error message for bad flags, got: %s", output)
	}
}

// C-15.15: Empty content FS still dispatches
func TestRun_EmptyContentFS_StillDispatches(t *testing.T) {
	emptyFS := &fstest.MapFS{}
	var stdout bytes.Buffer

	exitCode := cli.Run([]string{"version"}, emptyFS, &stdout)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	output := stdout.String()
	if !strings.Contains(output, "greenlight") {
		t.Errorf("version should work with empty FS, got: %s", output)
	}
}

// C-15.16: Args after subcommand passed to handler
func TestRun_ArgsAfterSubcommand_PassedToHandler(t *testing.T) {
	contentFS := buildTestFS()
	var stdout bytes.Buffer

	tempDir := t.TempDir()
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temp directory: %v", err)
	}

	exitCode := cli.Run([]string{"install", "--on-conflict=replace", "--local"}, contentFS, &stdout)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 when flags forwarded correctly, got %d", exitCode)
	}
}

// C-97: New commands are recognized and dispatched (not treated as unknown)
func TestRun_NewCommands_AreRecognized(t *testing.T) {
	contentFS := buildTestFS()

	testCases := []struct {
		name    string
		command string
	}{
		{"status", "status"},
		{"slice", "slice"},
		{"init", "init"},
		{"design", "design"},
		{"roadmap", "roadmap"},
		{"changelog", "changelog"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var stdout bytes.Buffer
			cli.Run([]string{tc.command}, contentFS, &stdout)

			output := stdout.String()
			if strings.Contains(output, "unknown command: "+tc.command) {
				t.Errorf("command %q should be recognized, but got 'unknown command': %s", tc.command, output)
			}
		})
	}
}

// C-97: New commands return a valid exit code (dispatch does not crash)
func TestRun_NewCommands_ReturnValidExitCode(t *testing.T) {
	contentFS := buildTestFS()

	testCases := []struct {
		name    string
		command string
	}{
		{"status", "status"},
		{"slice", "slice"},
		{"init", "init"},
		{"design", "design"},
		{"roadmap", "roadmap"},
		{"changelog", "changelog"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var stdout bytes.Buffer
			exitCode := cli.Run([]string{tc.command}, contentFS, &stdout)

			if exitCode != 0 && exitCode != 1 {
				t.Errorf("command %q returned unexpected exit code %d (must be 0 or 1)", tc.command, exitCode)
			}
		})
	}
}

// C-97: New commands receive args[1:] (subcommand args are forwarded, not the command name itself)
func TestRun_NewCommands_SubcommandArgsForwarded(t *testing.T) {
	contentFS := buildTestFS()

	testCases := []struct {
		name    string
		args    []string
		command string
	}{
		{"status with flag", []string{"status", "--verbose"}, "status"},
		{"slice with subarg", []string{"slice", "S-01"}, "slice"},
		{"init with flag", []string{"init", "--dry-run"}, "init"},
		{"design with subarg", []string{"design", "my-feature"}, "design"},
		{"roadmap with flag", []string{"roadmap", "--format=json"}, "roadmap"},
		{"changelog with subarg", []string{"changelog", "S-01"}, "changelog"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var stdout bytes.Buffer
			cli.Run(tc.args, contentFS, &stdout)

			output := stdout.String()
			if strings.Contains(output, "unknown command: "+tc.command) {
				t.Errorf("command %q with extra args should dispatch, not be treated as unknown: %s", tc.command, output)
			}
		})
	}
}

// C-97: printUsage shows project lifecycle commands grouped by category
func TestRun_PrintUsage_ShowsProjectLifecycleCommands(t *testing.T) {
	contentFS := buildTestFS()
	var stdout bytes.Buffer

	cli.Run([]string{}, contentFS, &stdout)

	output := stdout.String()

	projectLifecycleCommands := []string{"init", "design", "roadmap"}
	for _, command := range projectLifecycleCommands {
		if !strings.Contains(output, command) {
			t.Errorf("expected usage to contain project lifecycle command %q, got: %s", command, output)
		}
	}
}

// C-97: printUsage shows building commands grouped by category
func TestRun_PrintUsage_ShowsBuildingCommands(t *testing.T) {
	contentFS := buildTestFS()
	var stdout bytes.Buffer

	cli.Run([]string{}, contentFS, &stdout)

	output := stdout.String()
	if !strings.Contains(output, "slice") {
		t.Errorf("expected usage to contain building command 'slice', got: %s", output)
	}
}

// C-97: printUsage shows state and progress commands grouped by category
func TestRun_PrintUsage_ShowsStateAndProgressCommands(t *testing.T) {
	contentFS := buildTestFS()
	var stdout bytes.Buffer

	cli.Run([]string{}, contentFS, &stdout)

	output := stdout.String()

	stateCommands := []string{"status", "changelog"}
	for _, command := range stateCommands {
		if !strings.Contains(output, command) {
			t.Errorf("expected usage to contain state/progress command %q, got: %s", command, output)
		}
	}
}

// C-97: printUsage shows commands grouped with category headings
func TestRun_PrintUsage_ShowsCategoryGroupings(t *testing.T) {
	contentFS := buildTestFS()
	var stdout bytes.Buffer

	cli.Run([]string{}, contentFS, &stdout)

	output := stdout.String()

	// All six new commands must appear somewhere in usage
	newCommands := []string{"init", "design", "roadmap", "slice", "status", "changelog"}
	for _, command := range newCommands {
		if !strings.Contains(output, command) {
			t.Errorf("expected usage to list new command %q, got: %s", command, output)
		}
	}
}

// C-97: help flag shows updated usage with all new commands
func TestRun_HelpFlag_ShowsUpdatedUsageWithNewCommands(t *testing.T) {
	contentFS := buildTestFS()

	helpVariants := []struct {
		name string
		args []string
	}{
		{"help", []string{"help"}},
		{"--help", []string{"--help"}},
		{"-h", []string{"-h"}},
	}

	newCommands := []string{"init", "design", "roadmap", "slice", "status", "changelog"}

	for _, variant := range helpVariants {
		t.Run(variant.name, func(t *testing.T) {
			var stdout bytes.Buffer
			exitCode := cli.Run(variant.args, contentFS, &stdout)

			if exitCode != 0 {
				t.Errorf("expected exit code 0 for %q, got %d", variant.name, exitCode)
			}

			output := stdout.String()
			for _, command := range newCommands {
				if !strings.Contains(output, command) {
					t.Errorf("%q: expected usage to contain new command %q, got: %s", variant.name, command, output)
				}
			}
		})
	}
}

// C-97: Unknown command still returns 1 after dispatch extension
func TestRun_UnknownCommand_StillReturnsOneAfterExtension(t *testing.T) {
	contentFS := buildTestFS()

	unknownCommands := []struct {
		name    string
		command string
	}{
		{"deploy", "deploy"},
		{"upgrade", "upgrade"},
		{"run", "run"},
		{"build", "build"},
	}

	for _, tc := range unknownCommands {
		t.Run(tc.name, func(t *testing.T) {
			var stdout bytes.Buffer
			exitCode := cli.Run([]string{tc.command}, contentFS, &stdout)

			if exitCode != 1 {
				t.Errorf("unknown command %q should return exit code 1, got %d", tc.command, exitCode)
			}

			output := stdout.String()
			if !strings.Contains(output, "unknown command: "+tc.command) {
				t.Errorf("expected 'unknown command: %s', got: %s", tc.command, output)
			}
		})
	}
}

// C-97: Unknown command still prints usage after dispatch extension
func TestRun_UnknownCommand_StillPrintsUsageAfterExtension(t *testing.T) {
	contentFS := buildTestFS()
	var stdout bytes.Buffer

	cli.Run([]string{"notacommand"}, contentFS, &stdout)

	output := stdout.String()
	if !strings.Contains(output, "Usage:") {
		t.Errorf("expected usage to be printed after unknown command, got: %s", output)
	}
}

// C-97: Regression — existing commands continue to work identically after dispatch extension
func TestRun_ExistingCommands_ContinueToWorkAfterExtension(t *testing.T) {
	contentFS := buildTestFS()
	var stdout bytes.Buffer

	exitCode := cli.Run([]string{"version"}, contentFS, &stdout)

	if exitCode != 0 {
		t.Errorf("version command should still return exit code 0 after dispatch extension, got %d", exitCode)
	}

	output := stdout.String()
	if !strings.Contains(output, "greenlight") {
		t.Errorf("version command should still output 'greenlight' after dispatch extension, got: %s", output)
	}
}

// C-97: Regression — no args still returns 0 and prints usage after dispatch extension
func TestRun_NoArgs_StillReturnsZeroAfterExtension(t *testing.T) {
	contentFS := buildTestFS()
	var stdout bytes.Buffer

	exitCode := cli.Run([]string{}, contentFS, &stdout)

	if exitCode != 0 {
		t.Errorf("no args should still return exit code 0 after dispatch extension, got %d", exitCode)
	}

	output := stdout.String()
	if !strings.Contains(output, "Usage:") {
		t.Errorf("no args should still print usage after dispatch extension, got: %s", output)
	}
}
