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

// buildTestFS returns an fstest.MapFS with all 26 manifest files.
func buildTestFS() *fstest.MapFS {
	manifestFiles := []string{
		"agents/gl-architect.md",
		"agents/gl-codebase-mapper.md",
		"agents/gl-debugger.md",
		"agents/gl-designer.md",
		"agents/gl-implementer.md",
		"agents/gl-security.md",
		"agents/gl-test-writer.md",
		"agents/gl-verifier.md",
		"commands/gl/add-slice.md",
		"commands/gl/design.md",
		"commands/gl/help.md",
		"commands/gl/init.md",
		"commands/gl/map.md",
		"commands/gl/pause.md",
		"commands/gl/quick.md",
		"commands/gl/resume.md",
		"commands/gl/settings.md",
		"commands/gl/ship.md",
		"commands/gl/slice.md",
		"commands/gl/status.md",
		"references/checkpoint-protocol.md",
		"references/deviation-rules.md",
		"references/verification-patterns.md",
		"templates/config.md",
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
		// Check that at least one manifest file was installed
		claudeDir := filepath.Join(tempDir, ".claude")
		manifestPath := filepath.Join(claudeDir, "CLAUDE.md")
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
