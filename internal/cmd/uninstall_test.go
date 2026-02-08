package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/atlantic-blue/greenlight/internal/cmd"
)

// setupLocalUninstallTest creates a temp directory, changes to it, creates .claude with files,
// and returns cleanup function to restore original directory.
func setupLocalUninstallTest(t *testing.T) (targetDir string, cleanup func()) {
	t.Helper()

	// Remember original directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getting working directory: %v", err)
	}

	// Create and change to temp directory
	tempRoot := t.TempDir()
	if err := os.Chdir(tempRoot); err != nil {
		t.Fatalf("changing to temp directory: %v", err)
	}

	// Create .claude directory
	targetDir = filepath.Join(tempRoot, ".claude")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		t.Fatalf("creating .claude directory: %v", err)
	}

	// Create some manifest files
	manifestFiles := []string{
		"agents/gl-architect.md",
		"commands/gl/help.md",
		"references/checkpoint-protocol.md",
	}

	for _, relPath := range manifestFiles {
		fullPath := filepath.Join(targetDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("creating directory for %s: %v", relPath, err)
		}
		if err := os.WriteFile(fullPath, []byte("test content"), 0644); err != nil {
			t.Fatalf("writing %s: %v", relPath, err)
		}
	}

	// Create .greenlight-version
	versionPath := filepath.Join(targetDir, ".greenlight-version")
	if err := os.WriteFile(versionPath, []byte("1.0.0"), 0644); err != nil {
		t.Fatalf("writing version file: %v", err)
	}

	// Create CLAUDE.md in parent (for local scope)
	claudeMdPath := filepath.Join(tempRoot, "CLAUDE.md")
	if err := os.WriteFile(claudeMdPath, []byte("user content"), 0644); err != nil {
		t.Fatalf("writing CLAUDE.md: %v", err)
	}

	cleanup = func() {
		os.Chdir(originalDir)
	}

	return targetDir, cleanup
}

// setupGlobalUninstallTest creates files in a temp directory simulating global install.
func setupGlobalUninstallTest(t *testing.T) string {
	t.Helper()

	targetDir := t.TempDir()

	// Create some manifest files
	manifestFiles := []string{
		"agents/gl-architect.md",
		"commands/gl/help.md",
		"references/checkpoint-protocol.md",
	}

	for _, relPath := range manifestFiles {
		fullPath := filepath.Join(targetDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("creating directory for %s: %v", relPath, err)
		}
		if err := os.WriteFile(fullPath, []byte("test content"), 0644); err != nil {
			t.Fatalf("writing %s: %v", relPath, err)
		}
	}

	// Create .greenlight-version
	versionPath := filepath.Join(targetDir, ".greenlight-version")
	if err := os.WriteFile(versionPath, []byte("1.0.0"), 0644); err != nil {
		t.Fatalf("writing version file: %v", err)
	}

	// Create CLAUDE.md (for global scope, it's in targetDir)
	claudeMdPath := filepath.Join(targetDir, "CLAUDE.md")
	if err := os.WriteFile(claudeMdPath, []byte("user content"), 0644); err != nil {
		t.Fatalf("writing CLAUDE.md: %v", err)
	}

	return targetDir
}

// Test C-14.1: Returns 0 on success with --local
func TestRunUninstall_ReturnsZeroOnSuccess_Local(t *testing.T) {
	targetDir, cleanup := setupLocalUninstallTest(t)
	defer cleanup()

	var stdout bytes.Buffer
	exitCode := cmd.RunUninstall([]string{"--local"}, &stdout)

	if exitCode != 0 {
		t.Errorf("exit code = %d, want 0\noutput: %s", exitCode, stdout.String())
	}

	// Verify files were removed
	testFile := filepath.Join(targetDir, "agents", "gl-architect.md")
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Errorf("file should have been removed")
	}
}

// Test C-14.2: Returns 1 when no scope flag
func TestRunUninstall_ReturnsOneWhenNoScopeFlag(t *testing.T) {
	var stdout bytes.Buffer
	exitCode := cmd.RunUninstall([]string{}, &stdout)

	if exitCode != 1 {
		t.Errorf("exit code = %d, want 1", exitCode)
	}
}

// Test C-14.3: Returns 1 when both scope flags
func TestRunUninstall_ReturnsOneWhenBothScopeFlags(t *testing.T) {
	var stdout bytes.Buffer
	exitCode := cmd.RunUninstall([]string{"--local", "--global"}, &stdout)

	if exitCode != 1 {
		t.Errorf("exit code = %d, want 1", exitCode)
	}
}

// Test C-14.4: Prints "error: " prefix on failure
func TestRunUninstall_PrintsErrorPrefix(t *testing.T) {
	var stdout bytes.Buffer
	cmd.RunUninstall([]string{}, &stdout)

	output := stdout.String()
	if !strings.HasPrefix(output, "error: ") {
		t.Errorf("output should start with 'error: ', got: %s", output)
	}
}

// Test C-14.5: Exit code 0 iff uninstall completed successfully
func TestRunUninstall_ExitCodeZeroOnlyOnSuccess(t *testing.T) {
	tests := []struct {
		name           string
		setupFunc      func(t *testing.T) (args []string, cleanup func())
		wantExitCode   int
		wantSuccess    bool
	}{
		{
			name: "success with local scope",
			setupFunc: func(t *testing.T) ([]string, func()) {
				_, cleanup := setupLocalUninstallTest(t)
				return []string{"--local"}, cleanup
			},
			wantExitCode: 0,
			wantSuccess:  true,
		},
		{
			name: "failure with no scope",
			setupFunc: func(t *testing.T) ([]string, func()) {
				return []string{}, func() {}
			},
			wantExitCode: 1,
			wantSuccess:  false,
		},
		{
			name: "failure with both scopes",
			setupFunc: func(t *testing.T) ([]string, func()) {
				return []string{"--local", "--global"}, func() {}
			},
			wantExitCode: 1,
			wantSuccess:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, cleanup := tt.setupFunc(t)
			defer cleanup()

			var stdout bytes.Buffer
			exitCode := cmd.RunUninstall(args, &stdout)

			if exitCode != tt.wantExitCode {
				t.Errorf("exit code = %d, want %d", exitCode, tt.wantExitCode)
			}

			output := stdout.String()
			hasError := strings.HasPrefix(output, "error: ")

			if tt.wantSuccess && hasError {
				t.Errorf("expected success but got error: %s", output)
			}
			if !tt.wantSuccess && !hasError {
				t.Errorf("expected error but got: %s", output)
			}
		})
	}
}

// Test C-14.6: Passes scope to Uninstall
func TestRunUninstall_PassesScopeToUninstall(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		setupFunc func(t *testing.T) (string, func())
	}{
		{
			name: "local scope",
			args: []string{"--local"},
			setupFunc: func(t *testing.T) (string, func()) {
				targetDir, cleanup := setupLocalUninstallTest(t)
				return targetDir, cleanup
			},
		},
		{
			name: "global scope",
			args: []string{"--global"},
			setupFunc: func(t *testing.T) (string, func()) {
				// For global scope, we need to set HOME
				targetDir := setupGlobalUninstallTest(t)

				originalHome := os.Getenv("HOME")
				os.Setenv("HOME", filepath.Dir(targetDir))

				cleanup := func() {
					os.Setenv("HOME", originalHome)
				}

				// Rename targetDir to be ~/.claude for the test
				homeDir := filepath.Dir(targetDir)
				claudeDir := filepath.Join(homeDir, ".claude")
				if err := os.Rename(targetDir, claudeDir); err != nil {
					t.Fatalf("renaming to .claude: %v", err)
				}

				return claudeDir, cleanup
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetDir, cleanup := tt.setupFunc(t)
			defer cleanup()

			var stdout bytes.Buffer
			exitCode := cmd.RunUninstall(tt.args, &stdout)

			if exitCode != 0 {
				t.Fatalf("uninstall failed with exit code %d: %s", exitCode, stdout.String())
			}

			// Verify files were removed from the correct target directory
			testFile := filepath.Join(targetDir, "agents", "gl-architect.md")
			if _, err := os.Stat(testFile); !os.IsNotExist(err) {
				t.Errorf("file should have been removed from %s", targetDir)
			}
		})
	}
}

// Test C-14.7: Output written to provided writer
func TestRunUninstall_OutputWrittenToProvidedWriter(t *testing.T) {
	_, cleanup := setupLocalUninstallTest(t)
	defer cleanup()

	var stdout bytes.Buffer
	cmd.RunUninstall([]string{"--local"}, &stdout)

	output := stdout.String()
	if output == "" {
		t.Errorf("no output written to provided writer")
	}

	// Should contain uninstall message
	if !strings.Contains(output, "greenlight uninstalled from") {
		t.Errorf("output should contain uninstall message, got: %s", output)
	}
}

// Test C-14.8: Returns 1 on filesystem error during uninstall
func TestRunUninstall_ReturnsOneOnFilesystemError(t *testing.T) {
	targetDir, cleanup := setupLocalUninstallTest(t)
	defer cleanup()

	// Make a directory read-only to cause uninstall to fail
	agentsDir := filepath.Join(targetDir, "agents")
	if err := os.Chmod(agentsDir, 0555); err != nil {
		t.Fatalf("making directory read-only: %v", err)
	}
	defer os.Chmod(agentsDir, 0755)

	var stdout bytes.Buffer
	exitCode := cmd.RunUninstall([]string{"--local"}, &stdout)

	if exitCode != 1 {
		t.Errorf("exit code = %d, want 1 on filesystem error", exitCode)
	}

	output := stdout.String()
	if !strings.HasPrefix(output, "error: ") {
		t.Errorf("output should start with 'error: ' on failure, got: %s", output)
	}
}

// Test C-14.9: Invalid scope flag returns 1
func TestRunUninstall_InvalidScopeFlagReturnsOne(t *testing.T) {
	var stdout bytes.Buffer
	exitCode := cmd.RunUninstall([]string{"--invalid"}, &stdout)

	if exitCode != 1 {
		t.Errorf("exit code = %d, want 1 for invalid flag", exitCode)
	}

	output := stdout.String()
	if !strings.HasPrefix(output, "error: ") {
		t.Errorf("output should start with 'error: ' on invalid flag, got: %s", output)
	}
}

// Test C-14.10: Multiple --local flags treated as single flag
func TestRunUninstall_MultipleLocalFlagsTreatedAsSingle(t *testing.T) {
	_, cleanup := setupLocalUninstallTest(t)
	defer cleanup()

	var stdout bytes.Buffer
	exitCode := cmd.RunUninstall([]string{"--local", "--local"}, &stdout)

	// Should succeed - multiple same flags should be idempotent
	if exitCode != 0 {
		t.Errorf("exit code = %d, want 0 for multiple --local flags", exitCode)
	}
}

// Test C-14.11: Error message includes context
func TestRunUninstall_ErrorMessageIncludesContext(t *testing.T) {
	var stdout bytes.Buffer
	cmd.RunUninstall([]string{}, &stdout)

	output := stdout.String()

	// Error message should provide some context about what went wrong
	if len(output) < 20 {
		t.Errorf("error message too short, should include context: %s", output)
	}
}

// Test C-14.12: Success with --global flag
func TestRunUninstall_SuccessWithGlobalFlag(t *testing.T) {
	targetDir := setupGlobalUninstallTest(t)

	// Set HOME to parent of targetDir
	originalHome := os.Getenv("HOME")
	homeDir := filepath.Dir(targetDir)
	os.Setenv("HOME", homeDir)
	defer os.Setenv("HOME", originalHome)

	// Rename targetDir to be ~/.claude
	claudeDir := filepath.Join(homeDir, ".claude")
	if err := os.Rename(targetDir, claudeDir); err != nil {
		t.Fatalf("renaming to .claude: %v", err)
	}

	var stdout bytes.Buffer
	exitCode := cmd.RunUninstall([]string{"--global"}, &stdout)

	if exitCode != 0 {
		t.Errorf("exit code = %d, want 0 for --global\noutput: %s", exitCode, stdout.String())
	}

	// Verify files were removed
	testFile := filepath.Join(claudeDir, "agents", "gl-architect.md")
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Errorf("file should have been removed from global location")
	}
}
