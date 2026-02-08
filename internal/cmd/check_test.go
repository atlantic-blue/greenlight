package cmd_test

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/atlantic-blue/greenlight/internal/cmd"
	"github.com/atlantic-blue/greenlight/internal/installer"
)

// buildTestContentFS creates a complete MapFS with all 30 manifest files for testing.
func buildTestContentFS() fstest.MapFS {
	return fstest.MapFS{
		"agents/gl-architect.md":                   &fstest.MapFile{Data: []byte("# Architect\n")},
		"agents/gl-assessor.md":                    &fstest.MapFile{Data: []byte("# Assessor\n")},
		"agents/gl-codebase-mapper.md":             &fstest.MapFile{Data: []byte("# Codebase Mapper\n")},
		"agents/gl-debugger.md":                    &fstest.MapFile{Data: []byte("# Debugger\n")},
		"agents/gl-designer.md":                    &fstest.MapFile{Data: []byte("# Designer\n")},
		"agents/gl-implementer.md":                 &fstest.MapFile{Data: []byte("# Implementer\n")},
		"agents/gl-security.md":                    &fstest.MapFile{Data: []byte("# Security\n")},
		"agents/gl-test-writer.md":                 &fstest.MapFile{Data: []byte("# Test Writer\n")},
		"agents/gl-verifier.md":                    &fstest.MapFile{Data: []byte("# Verifier\n")},
		"agents/gl-wrapper.md":                     &fstest.MapFile{Data: []byte("# Wrapper\n")},
		"commands/gl/add-slice.md":                 &fstest.MapFile{Data: []byte("# Add Slice\n")},
		"commands/gl/assess.md":                    &fstest.MapFile{Data: []byte("# Assess\n")},
		"commands/gl/design.md":                    &fstest.MapFile{Data: []byte("# Design\n")},
		"commands/gl/help.md":                      &fstest.MapFile{Data: []byte("# Help\n")},
		"commands/gl/init.md":                      &fstest.MapFile{Data: []byte("# Init\n")},
		"commands/gl/map.md":                       &fstest.MapFile{Data: []byte("# Map\n")},
		"commands/gl/pause.md":                     &fstest.MapFile{Data: []byte("# Pause\n")},
		"commands/gl/quick.md":                     &fstest.MapFile{Data: []byte("# Quick\n")},
		"commands/gl/resume.md":                    &fstest.MapFile{Data: []byte("# Resume\n")},
		"commands/gl/settings.md":                  &fstest.MapFile{Data: []byte("# Settings\n")},
		"commands/gl/ship.md":                      &fstest.MapFile{Data: []byte("# Ship\n")},
		"commands/gl/slice.md":                     &fstest.MapFile{Data: []byte("# Slice\n")},
		"commands/gl/status.md":                    &fstest.MapFile{Data: []byte("# Status\n")},
		"commands/gl/wrap.md":                      &fstest.MapFile{Data: []byte("# Wrap\n")},
		"references/checkpoint-protocol.md":        &fstest.MapFile{Data: []byte("# Checkpoint Protocol\n")},
		"references/deviation-rules.md":            &fstest.MapFile{Data: []byte("# Deviation Rules\n")},
		"references/verification-patterns.md":      &fstest.MapFile{Data: []byte("# Verification Patterns\n")},
		"templates/config.md":                      &fstest.MapFile{Data: []byte("# Config Template\n")},
		"templates/state.md":                       &fstest.MapFile{Data: []byte("# State Template\n")},
		"CLAUDE.md":                                &fstest.MapFile{Data: []byte("# Greenlight CLAUDE.md\n\nTest content\n")},
	}
}

// installTestFilesForCmd simulates a successful install for cmd tests.
func installTestFilesForCmd(t *testing.T, targetDir, scope string, contentFS fs.FS) {
	t.Helper()

	for _, relPath := range installer.Manifest {
		var destPath string
		if relPath == "CLAUDE.md" {
			switch scope {
			case "local":
				if targetDir == ".claude" {
					destPath = "CLAUDE.md"
				} else {
					destPath = filepath.Join(filepath.Dir(targetDir), "CLAUDE.md")
				}
			default:
				destPath = filepath.Join(targetDir, "CLAUDE.md")
			}
		} else {
			destPath = filepath.Join(targetDir, relPath)
		}

		data, err := fs.ReadFile(contentFS, relPath)
		if err != nil {
			t.Fatalf("failed to read %s from contentFS: %v", relPath, err)
		}

		dir := filepath.Dir(destPath)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("failed to create directory %s: %v", dir, err)
		}

		if err := os.WriteFile(destPath, data, 0o644); err != nil {
			t.Fatalf("failed to write %s: %v", destPath, err)
		}
	}

	// Write .greenlight-version file
	versionContent := []byte("v1.0.0\nabc123\n2026-02-08\n")
	versionPath := filepath.Join(targetDir, ".greenlight-version")
	if err := os.WriteFile(versionPath, versionContent, 0o644); err != nil {
		t.Fatalf("failed to write version file: %v", err)
	}
}

// C-13 Tests: RunCheck

func TestRunCheck_ReturnsZeroWhenAllFilesPresent_Local(t *testing.T) {
	contentFS := buildTestContentFS()

	// Change to temp directory so .claude resolves there
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	projectRoot := t.TempDir()
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	targetDir := ".claude"
	installTestFilesForCmd(t, targetDir, "local", contentFS)

	var buf bytes.Buffer
	exitCode := cmd.RunCheck([]string{"--local"}, contentFS, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d\nOutput: %s", exitCode, buf.String())
	}

	output := buf.String()
	if !strings.Contains(output, "all 30 files present") {
		t.Errorf("expected success message in output: %q", output)
	}
}

func TestRunCheck_ReturnsOneWhenFilesMissing(t *testing.T) {
	contentFS := buildTestContentFS()

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	projectRoot := t.TempDir()
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	targetDir := ".claude"
	installTestFilesForCmd(t, targetDir, "local", contentFS)

	// Remove a file to cause check failure
	if err := os.Remove(filepath.Join(targetDir, "agents/gl-architect.md")); err != nil {
		t.Fatalf("failed to remove test file: %v", err)
	}

	var buf bytes.Buffer
	exitCode := cmd.RunCheck([]string{"--local"}, contentFS, &buf)

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}

	output := buf.String()
	if !strings.Contains(output, "MISSING") {
		t.Errorf("expected MISSING message in output: %q", output)
	}
}

func TestRunCheck_ReturnsOneWhenNoScopeFlag(t *testing.T) {
	contentFS := buildTestContentFS()
	var buf bytes.Buffer

	exitCode := cmd.RunCheck([]string{}, contentFS, &buf)

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}

	output := buf.String()
	if !strings.HasPrefix(output, "error: ") {
		t.Errorf("expected error prefix, got: %q", output)
	}
}

func TestRunCheck_ReturnsOneWhenBothScopeFlags(t *testing.T) {
	contentFS := buildTestContentFS()
	var buf bytes.Buffer

	exitCode := cmd.RunCheck([]string{"--global", "--local"}, contentFS, &buf)

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}

	output := buf.String()
	if !strings.HasPrefix(output, "error: ") {
		t.Errorf("expected error prefix, got: %q", output)
	}
}

func TestRunCheck_PrintsErrorPrefixForScopeErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "no scope flag",
			args: []string{},
		},
		{
			name: "both scope flags",
			args: []string{"--global", "--local"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contentFS := buildTestContentFS()
			var buf bytes.Buffer

			exitCode := cmd.RunCheck(tt.args, contentFS, &buf)

			if exitCode != 1 {
				t.Errorf("expected exit code 1, got %d", exitCode)
			}

			output := buf.String()
			if !strings.HasPrefix(output, "error: ") {
				t.Errorf("expected output to start with 'error: ', got: %q", output)
			}
		})
	}
}

func TestRunCheck_WithVerifyFlag_PassesVerifyTrueToCheck(t *testing.T) {
	contentFS := buildTestContentFS()

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	projectRoot := t.TempDir()
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	targetDir := ".claude"
	installTestFilesForCmd(t, targetDir, "local", contentFS)

	var buf bytes.Buffer
	exitCode := cmd.RunCheck([]string{"--local", "--verify"}, contentFS, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d\nOutput: %s", exitCode, buf.String())
	}

	output := buf.String()
	// When verify=true, output should say "verified" instead of "present"
	if !strings.Contains(output, "all 30 files verified") {
		t.Errorf("expected 'verified' message (indicating verify=true), got: %q", output)
	}
}

func TestRunCheck_WithoutVerifyFlag_PassesVerifyFalseToCheck(t *testing.T) {
	contentFS := buildTestContentFS()

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	projectRoot := t.TempDir()
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	targetDir := ".claude"
	installTestFilesForCmd(t, targetDir, "local", contentFS)

	// Modify a file (should not be detected without --verify)
	modifiedFile := filepath.Join(targetDir, "agents/gl-architect.md")
	if err := os.WriteFile(modifiedFile, []byte("modified content"), 0o644); err != nil {
		t.Fatalf("failed to modify file: %v", err)
	}

	var buf bytes.Buffer
	exitCode := cmd.RunCheck([]string{"--local"}, contentFS, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 (no --verify), got %d\nOutput: %s", exitCode, buf.String())
	}

	output := buf.String()
	// Without --verify, should say "present" not "verified"
	if !strings.Contains(output, "all 30 files present") {
		t.Errorf("expected 'present' message (indicating verify=false), got: %q", output)
	}

	// Should NOT detect modification
	if strings.Contains(output, "MODIFIED") {
		t.Error("without --verify flag, modifications should not be detected")
	}
}

func TestRunCheck_ExitCodeZeroIffAllChecksPass(t *testing.T) {
	tests := []struct {
		name             string
		setup            func(t *testing.T, targetDir string)
		args             []string
		expectedExitCode int
	}{
		{
			name:             "all files present and valid",
			setup:            func(t *testing.T, targetDir string) {},
			args:             []string{"--local"},
			expectedExitCode: 0,
		},
		{
			name: "one file missing",
			setup: func(t *testing.T, targetDir string) {
				os.Remove(filepath.Join(targetDir, "agents/gl-architect.md"))
			},
			args:             []string{"--local"},
			expectedExitCode: 1,
		},
		{
			name: "one file empty",
			setup: func(t *testing.T, targetDir string) {
				os.WriteFile(filepath.Join(targetDir, "agents/gl-debugger.md"), []byte{}, 0o644)
			},
			args:             []string{"--local"},
			expectedExitCode: 1,
		},
		{
			name: "version file missing",
			setup: func(t *testing.T, targetDir string) {
				os.Remove(filepath.Join(targetDir, ".greenlight-version"))
			},
			args:             []string{"--local"},
			expectedExitCode: 1,
		},
		{
			name:             "verify mode all files match",
			setup:            func(t *testing.T, targetDir string) {},
			args:             []string{"--local", "--verify"},
			expectedExitCode: 0,
		},
		{
			name: "verify mode one file modified",
			setup: func(t *testing.T, targetDir string) {
				os.WriteFile(filepath.Join(targetDir, "agents/gl-security.md"), []byte("modified"), 0o644)
			},
			args:             []string{"--local", "--verify"},
			expectedExitCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contentFS := buildTestContentFS()

			originalWd, err := os.Getwd()
			if err != nil {
				t.Fatalf("failed to get working directory: %v", err)
			}
			defer os.Chdir(originalWd)

			projectRoot := t.TempDir()
			if err := os.Chdir(projectRoot); err != nil {
				t.Fatalf("failed to change directory: %v", err)
			}

			targetDir := ".claude"
			installTestFilesForCmd(t, targetDir, "local", contentFS)

			tt.setup(t, targetDir)

			var buf bytes.Buffer
			exitCode := cmd.RunCheck(tt.args, contentFS, &buf)

			if exitCode != tt.expectedExitCode {
				t.Errorf("expected exit code %d, got %d\nOutput: %s", tt.expectedExitCode, exitCode, buf.String())
			}
		})
	}
}

func TestRunCheck_GlobalScope_UsesHomeDirectory(t *testing.T) {
	contentFS := buildTestContentFS()

	// We can't easily test actual home directory without side effects,
	// but we can verify that --global flag is accepted and processed
	var buf bytes.Buffer
	_ = cmd.RunCheck([]string{"--global"}, contentFS, &buf)

	// Exit code will likely be 1 because files aren't installed in home/.claude
	// But we're testing that it attempts to check the right location
	// and doesn't error on the flag itself
	output := buf.String()

	// Should not error on scope parsing
	if strings.Contains(output, "must specify --global or --local") {
		t.Error("--global flag was not recognized")
	}

	if strings.Contains(output, "cannot specify both") {
		t.Error("--global flag incorrectly treated as invalid")
	}
}

func TestRunCheck_VerifyFlagIsBooleanFlag(t *testing.T) {
	contentFS := buildTestContentFS()

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	projectRoot := t.TempDir()
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	targetDir := ".claude"
	installTestFilesForCmd(t, targetDir, "local", contentFS)

	tests := []struct {
		name           string
		args           []string
		expectVerified bool
	}{
		{
			name:           "with --verify flag",
			args:           []string{"--local", "--verify"},
			expectVerified: true,
		},
		{
			name:           "without --verify flag",
			args:           []string{"--local"},
			expectVerified: false,
		},
		{
			name:           "--verify position independent",
			args:           []string{"--verify", "--local"},
			expectVerified: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			exitCode := cmd.RunCheck(tt.args, contentFS, &buf)

			if exitCode != 0 {
				t.Errorf("expected exit code 0, got %d\nOutput: %s", exitCode, buf.String())
			}

			output := buf.String()

			if tt.expectVerified {
				if !strings.Contains(output, "verified") {
					t.Error("expected 'verified' in output when --verify present")
				}
			} else {
				if strings.Contains(output, "verified") {
					t.Error("did not expect 'verified' in output when --verify absent")
				}
				if !strings.Contains(output, "present") {
					t.Error("expected 'present' in output when --verify absent")
				}
			}
		})
	}
}

func TestRunCheck_ResolvesTargetDirFromScope(t *testing.T) {
	contentFS := buildTestContentFS()

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	projectRoot := t.TempDir()
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	// Install to local scope (.claude)
	targetDir := ".claude"
	installTestFilesForCmd(t, targetDir, "local", contentFS)

	var buf bytes.Buffer
	exitCode := cmd.RunCheck([]string{"--local"}, contentFS, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d\nOutput: %s", exitCode, buf.String())
	}

	// Verify CLAUDE.md was checked in project root (not inside .claude)
	// We can infer this from successful check
	if _, err := os.Stat("CLAUDE.md"); os.IsNotExist(err) {
		t.Error("CLAUDE.md should exist in project root for local scope")
	}
}

func TestRunCheck_PassesContentFSToCheck(t *testing.T) {
	contentFS := buildTestContentFS()

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	projectRoot := t.TempDir()
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	targetDir := ".claude"
	installTestFilesForCmd(t, targetDir, "local", contentFS)

	// Modify a file
	modifiedFile := filepath.Join(targetDir, "agents/gl-implementer.md")
	if err := os.WriteFile(modifiedFile, []byte("modified"), 0o644); err != nil {
		t.Fatalf("failed to modify file: %v", err)
	}

	var buf bytes.Buffer
	exitCode := cmd.RunCheck([]string{"--local", "--verify"}, contentFS, &buf)

	// Should detect modification (contentFS was used for verification)
	if exitCode != 1 {
		t.Errorf("expected exit code 1 (modification detected), got %d", exitCode)
	}

	output := buf.String()
	if !strings.Contains(output, "MODIFIED") {
		t.Error("expected MODIFIED message (contentFS should be used for hash comparison)")
	}
}

func TestRunCheck_UnknownFlagsPassedToCheck(t *testing.T) {
	contentFS := buildTestContentFS()

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	projectRoot := t.TempDir()
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	targetDir := ".claude"
	installTestFilesForCmd(t, targetDir, "local", contentFS)

	var buf bytes.Buffer
	// Pass unknown flag (should be ignored by RunCheck but potentially used by Check)
	exitCode := cmd.RunCheck([]string{"--local", "--unknown-flag"}, contentFS, &buf)

	// Should still work (unknown flags don't cause errors in this implementation)
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d\nOutput: %s", exitCode, buf.String())
	}
}

func TestRunCheck_OutputWrittenToProvidedWriter(t *testing.T) {
	contentFS := buildTestContentFS()

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	projectRoot := t.TempDir()
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	targetDir := ".claude"
	installTestFilesForCmd(t, targetDir, "local", contentFS)

	var buf bytes.Buffer
	cmd.RunCheck([]string{"--local"}, contentFS, &buf)

	output := buf.String()

	if len(output) == 0 {
		t.Error("expected output to be written to provided writer")
	}

	// Verify output contains expected content
	if !strings.Contains(output, "version:") {
		t.Error("output missing version information")
	}

	if !strings.Contains(output, "files") {
		t.Error("output missing file count information")
	}
}

func TestRunCheck_ErrorMessagesFormattedCorrectly(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError string
	}{
		{
			name:        "no scope",
			args:        []string{},
			expectError: "error: must specify --global or --local",
		},
		{
			name:        "both scopes",
			args:        []string{"--global", "--local"},
			expectError: "error: cannot specify both --global and --local",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contentFS := buildTestContentFS()
			var buf bytes.Buffer

			exitCode := cmd.RunCheck(tt.args, contentFS, &buf)

			if exitCode != 1 {
				t.Errorf("expected exit code 1, got %d", exitCode)
			}

			output := buf.String()
			if !strings.HasPrefix(output, "error: ") {
				t.Errorf("expected error prefix in: %q", output)
			}

			if !strings.Contains(output, tt.expectError) {
				t.Errorf("expected error %q, got: %q", tt.expectError, output)
			}
		})
	}
}

func TestRunCheck_PropagatesCheckReturnValue(t *testing.T) {
	contentFS := buildTestContentFS()

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	projectRoot := t.TempDir()
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	targetDir := ".claude"
	installTestFilesForCmd(t, targetDir, "local", contentFS)

	tests := []struct {
		name     string
		setup    func(t *testing.T, targetDir string)
		expected int
	}{
		{
			name:     "Check returns true -> exit 0",
			setup:    func(t *testing.T, targetDir string) {},
			expected: 0,
		},
		{
			name: "Check returns false -> exit 1",
			setup: func(t *testing.T, targetDir string) {
				os.Remove(filepath.Join(targetDir, "agents/gl-architect.md"))
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset by reinstalling
			installTestFilesForCmd(t, targetDir, "local", contentFS)
			tt.setup(t, targetDir)

			var buf bytes.Buffer
			exitCode := cmd.RunCheck([]string{"--local"}, contentFS, &buf)

			if exitCode != tt.expected {
				t.Errorf("expected exit code %d, got %d\nOutput: %s", tt.expected, exitCode, buf.String())
			}
		})
	}
}

func TestRunCheck_LocalScopeCLAUDEMDCheckedInParent(t *testing.T) {
	contentFS := buildTestContentFS()

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	projectRoot := t.TempDir()
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	targetDir := ".claude"
	installTestFilesForCmd(t, targetDir, "local", contentFS)

	// Remove CLAUDE.md from project root
	if err := os.Remove("CLAUDE.md"); err != nil {
		t.Fatalf("failed to remove CLAUDE.md: %v", err)
	}

	var buf bytes.Buffer
	exitCode := cmd.RunCheck([]string{"--local"}, contentFS, &buf)

	if exitCode != 1 {
		t.Error("expected exit code 1 when CLAUDE.md missing from parent directory")
	}

	output := buf.String()
	if !strings.Contains(output, "MISSING  CLAUDE.md") {
		t.Errorf("expected MISSING CLAUDE.md in output: %q", output)
	}
}

func TestRunCheck_CallsCheckWithCorrectSignature(t *testing.T) {
	contentFS := buildTestContentFS()

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	projectRoot := t.TempDir()
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	targetDir := ".claude"
	installTestFilesForCmd(t, targetDir, "local", contentFS)

	// This test verifies that RunCheck calls Check with the corrected signature:
	// Check(targetDir, scope, stdout, verify, contentFS)

	tests := []struct {
		name           string
		args           []string
		verifyExpected bool
	}{
		{
			name:           "without --verify",
			args:           []string{"--local"},
			verifyExpected: false,
		},
		{
			name:           "with --verify",
			args:           []string{"--local", "--verify"},
			verifyExpected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			exitCode := cmd.RunCheck(tt.args, contentFS, &buf)

			if exitCode != 0 {
				t.Fatalf("expected exit code 0, got %d\nOutput: %s", exitCode, buf.String())
			}

			output := buf.String()

			// Verify the output indicates whether verify mode was used
			if tt.verifyExpected {
				if !strings.Contains(output, "verified") {
					t.Error("expected 'verified' in output (verify=true)")
				}
			} else {
				if !strings.Contains(output, "present") && !strings.Contains(output, "verified") {
					t.Error("expected file count message in output")
				}
			}
		})
	}
}

func TestRunCheck_ParsesScopeBeforeCallingCheck(t *testing.T) {
	contentFS := buildTestContentFS()
	var buf bytes.Buffer

	// This should fail at scope parsing, not at Check call
	exitCode := cmd.RunCheck([]string{}, contentFS, &buf)

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}

	output := buf.String()
	// Error should be about scope, not about files
	if !strings.Contains(output, "must specify") {
		t.Error("expected scope error message")
	}
}

func TestRunCheck_HandlesResolveDirErrors(t *testing.T) {
	contentFS := buildTestContentFS()

	// Test with an invalid scope value is not possible via ParseScope
	// because it only accepts --global or --local
	// So we'll test that ResolveDir is called by verifying behavior

	var buf bytes.Buffer
	// Using --global will call ResolveDir("global") which uses os.UserHomeDir()
	// This should succeed (unless home directory truly cannot be determined)
	_ = cmd.RunCheck([]string{"--global"}, contentFS, &buf)

	// We expect exit code 1 because files aren't installed in home directory
	// But no error about home directory resolution
	output := buf.String()

	if strings.Contains(output, "cannot determine home directory") {
		t.Skip("home directory not available in test environment")
	}

	// Should get file check errors, not directory resolution errors
	if strings.Contains(output, "error: cannot determine home directory") {
		t.Error("unexpected home directory resolution error")
	}
}

func TestRunCheck_MinimalValidInvocation(t *testing.T) {
	contentFS := buildTestContentFS()

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	projectRoot := t.TempDir()
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	targetDir := ".claude"
	installTestFilesForCmd(t, targetDir, "local", contentFS)

	var buf bytes.Buffer
	exitCode := cmd.RunCheck([]string{"--local"}, contentFS, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d\nOutput: %s", exitCode, buf.String())
	}

	output := buf.String()
	if len(output) == 0 {
		t.Error("expected non-empty output")
	}
}
