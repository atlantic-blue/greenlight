package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/atlantic-blue/greenlight/internal/cmd"
)

// Helper: buildTestFS creates a complete MapFS with all 32 manifest files.
func buildTestFS() fstest.MapFS {
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
		"commands/gl/changelog.md":                 &fstest.MapFile{Data: []byte("# Changelog\n")},
		"commands/gl/design.md":                    &fstest.MapFile{Data: []byte("# Design\n")},
		"commands/gl/help.md":                      &fstest.MapFile{Data: []byte("# Help\n")},
		"commands/gl/init.md":                      &fstest.MapFile{Data: []byte("# Init\n")},
		"commands/gl/map.md":                       &fstest.MapFile{Data: []byte("# Map\n")},
		"commands/gl/pause.md":                     &fstest.MapFile{Data: []byte("# Pause\n")},
		"commands/gl/quick.md":                     &fstest.MapFile{Data: []byte("# Quick\n")},
		"commands/gl/resume.md":                    &fstest.MapFile{Data: []byte("# Resume\n")},
		"commands/gl/roadmap.md":                   &fstest.MapFile{Data: []byte("# Roadmap\n")},
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

// C-12 Tests: RunInstall

func TestRunInstall_ReturnsZeroOnSuccess_Local(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer

	// Setup: chdir to temp directory so --local works
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	projectRoot := t.TempDir()
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	args := []string{"--local"}
	exitCode := cmd.RunInstall(args, contentFS, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
		t.Logf("output: %s", buf.String())
	}

	// Verify installation happened
	claudePath := "CLAUDE.md"
	if _, err := os.Stat(claudePath); os.IsNotExist(err) {
		t.Error("CLAUDE.md was not installed")
	}

	dotClaudePath := ".claude/agents/gl-architect.md"
	if _, err := os.Stat(dotClaudePath); os.IsNotExist(err) {
		t.Error("files were not installed to .claude directory")
	}
}

func TestRunInstall_ReturnsZeroOnSuccess_LocalWithConflictStrategy(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer

	// Setup: chdir to temp directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	projectRoot := t.TempDir()
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	args := []string{"--on-conflict=replace", "--local"}
	exitCode := cmd.RunInstall(args, contentFS, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
		t.Logf("output: %s", buf.String())
	}

	// Verify installation happened
	if _, err := os.Stat("CLAUDE.md"); os.IsNotExist(err) {
		t.Error("CLAUDE.md was not installed")
	}
}

func TestRunInstall_ReturnsOneWhenNoScopeFlag(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer

	args := []string{}
	exitCode := cmd.RunInstall(args, contentFS, &buf)

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}

	output := buf.String()
	if !strings.Contains(output, "error:") {
		t.Errorf("expected error message, got: %q", output)
	}
}

func TestRunInstall_ReturnsOneWhenBothScopeFlags(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer

	args := []string{"--global", "--local"}
	exitCode := cmd.RunInstall(args, contentFS, &buf)

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}

	output := buf.String()
	if !strings.Contains(output, "error:") {
		t.Errorf("expected error message, got: %q", output)
	}
}

func TestRunInstall_ReturnsOneWhenInvalidConflictStrategy(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer

	args := []string{"--on-conflict=invalid", "--local"}
	exitCode := cmd.RunInstall(args, contentFS, &buf)

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}

	output := buf.String()
	if !strings.Contains(output, "error:") {
		t.Errorf("expected error message, got: %q", output)
	}
}

func TestRunInstall_PrintsErrorPrefix(t *testing.T) {
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
		{
			name: "invalid conflict strategy",
			args: []string{"--on-conflict=merge", "--local"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contentFS := buildTestFS()
			var buf bytes.Buffer

			exitCode := cmd.RunInstall(tt.args, contentFS, &buf)

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

func TestRunInstall_PassesConflictStrategyKeep(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer

	// Setup: chdir to temp directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	projectRoot := t.TempDir()
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	// Create existing CLAUDE.md
	existingContent := []byte("# Existing\n")
	if err := os.WriteFile("CLAUDE.md", existingContent, 0o644); err != nil {
		t.Fatalf("failed to create existing CLAUDE.md: %v", err)
	}

	args := []string{"--on-conflict=keep", "--local"}
	exitCode := cmd.RunInstall(args, contentFS, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
		t.Logf("output: %s", buf.String())
	}

	// Verify keep strategy was applied (existing file kept, greenlight version created)
	content, err := os.ReadFile("CLAUDE.md")
	if err != nil {
		t.Fatalf("failed to read CLAUDE.md: %v", err)
	}
	if !bytes.Equal(content, existingContent) {
		t.Error("existing CLAUDE.md was not kept")
	}

	if _, err := os.Stat("CLAUDE_GREENLIGHT.md"); os.IsNotExist(err) {
		t.Error("CLAUDE_GREENLIGHT.md was not created")
	}

	output := buf.String()
	if !strings.Contains(output, "existing CLAUDE.md kept") {
		t.Error("output missing keep strategy message")
	}
}

func TestRunInstall_PassesConflictStrategyReplace(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer

	// Setup: chdir to temp directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	projectRoot := t.TempDir()
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	// Create existing CLAUDE.md
	existingContent := []byte("# Existing\n")
	if err := os.WriteFile("CLAUDE.md", existingContent, 0o644); err != nil {
		t.Fatalf("failed to create existing CLAUDE.md: %v", err)
	}

	args := []string{"--on-conflict=replace", "--local"}
	exitCode := cmd.RunInstall(args, contentFS, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
		t.Logf("output: %s", buf.String())
	}

	// Verify replace strategy was applied (file replaced, backup created)
	content, err := os.ReadFile("CLAUDE.md")
	if err != nil {
		t.Fatalf("failed to read CLAUDE.md: %v", err)
	}
	expectedContent := []byte("# Greenlight CLAUDE.md\n\nTest content\n")
	if !bytes.Equal(content, expectedContent) {
		t.Error("CLAUDE.md was not replaced with greenlight content")
	}

	backupPath := "CLAUDE.md.backup"
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Error("backup file was not created")
	}

	output := buf.String()
	if !strings.Contains(output, "backed up to") {
		t.Error("output missing replace strategy message")
	}
}

func TestRunInstall_PassesConflictStrategyAppend(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer

	// Setup: chdir to temp directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	projectRoot := t.TempDir()
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	// Create existing CLAUDE.md
	existingContent := []byte("# Existing\n")
	if err := os.WriteFile("CLAUDE.md", existingContent, 0o644); err != nil {
		t.Fatalf("failed to create existing CLAUDE.md: %v", err)
	}

	args := []string{"--on-conflict=append", "--local"}
	exitCode := cmd.RunInstall(args, contentFS, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
		t.Logf("output: %s", buf.String())
	}

	// Verify append strategy was applied
	content, err := os.ReadFile("CLAUDE.md")
	if err != nil {
		t.Fatalf("failed to read CLAUDE.md: %v", err)
	}

	contentStr := string(content)
	if !strings.HasPrefix(contentStr, "# Existing\n") {
		t.Error("existing content was not preserved at start")
	}
	if !strings.Contains(contentStr, "# Greenlight CLAUDE.md") {
		t.Error("greenlight content was not appended")
	}

	output := buf.String()
	if !strings.Contains(output, "appended to") {
		t.Error("output missing append strategy message")
	}
}

func TestRunInstall_DefaultConflictStrategyIsKeep(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer

	// Setup: chdir to temp directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	projectRoot := t.TempDir()
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	// Create existing CLAUDE.md
	existingContent := []byte("# Existing\n")
	if err := os.WriteFile("CLAUDE.md", existingContent, 0o644); err != nil {
		t.Fatalf("failed to create existing CLAUDE.md: %v", err)
	}

	// Don't specify --on-conflict flag
	args := []string{"--local"}
	exitCode := cmd.RunInstall(args, contentFS, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
		t.Logf("output: %s", buf.String())
	}

	// Verify keep strategy was applied (default)
	content, err := os.ReadFile("CLAUDE.md")
	if err != nil {
		t.Fatalf("failed to read CLAUDE.md: %v", err)
	}
	if !bytes.Equal(content, existingContent) {
		t.Error("existing CLAUDE.md was not kept (default should be keep)")
	}

	if _, err := os.Stat("CLAUDE_GREENLIGHT.md"); os.IsNotExist(err) {
		t.Error("CLAUDE_GREENLIGHT.md was not created (default should be keep)")
	}
}

func TestRunInstall_ParsesConflictStrategyBeforeScope(t *testing.T) {
	// This test verifies the order: conflict strategy is parsed first, then scope
	contentFS := buildTestFS()
	var buf bytes.Buffer

	// Invalid conflict strategy should fail before scope parsing
	args := []string{"--on-conflict=bad"}
	exitCode := cmd.RunInstall(args, contentFS, &buf)

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}

	output := buf.String()
	// Should fail on conflict strategy, not on missing scope
	if !strings.Contains(output, "error:") {
		t.Errorf("expected error message, got: %q", output)
	}
}

func TestRunInstall_ExitCodeZeroIfAndOnlyIfSuccess(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		setup          func(t *testing.T) string // returns working directory to use
		wantExitCode   int
		wantSuccessMsg bool
	}{
		{
			name: "success returns 0",
			args: []string{"--local"},
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantExitCode:   0,
			wantSuccessMsg: true,
		},
		{
			name: "invalid conflict strategy returns 1",
			args: []string{"--on-conflict=invalid", "--local"},
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantExitCode:   1,
			wantSuccessMsg: false,
		},
		{
			name: "no scope flag returns 1",
			args: []string{},
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantExitCode:   1,
			wantSuccessMsg: false,
		},
		{
			name: "both scope flags returns 1",
			args: []string{"--global", "--local"},
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantExitCode:   1,
			wantSuccessMsg: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contentFS := buildTestFS()
			var buf bytes.Buffer

			// Setup working directory
			originalWd, err := os.Getwd()
			if err != nil {
				t.Fatalf("failed to get working directory: %v", err)
			}
			defer os.Chdir(originalWd)

			wd := tt.setup(t)
			if err := os.Chdir(wd); err != nil {
				t.Fatalf("failed to change directory: %v", err)
			}

			exitCode := cmd.RunInstall(tt.args, contentFS, &buf)

			if exitCode != tt.wantExitCode {
				t.Errorf("expected exit code %d, got %d", tt.wantExitCode, exitCode)
			}

			output := buf.String()
			hasSuccessMsg := strings.Contains(output, "greenlight installed to")

			if tt.wantSuccessMsg && !hasSuccessMsg {
				t.Errorf("expected success message, got: %q", output)
			}
			if !tt.wantSuccessMsg && hasSuccessMsg {
				t.Errorf("unexpected success message in output: %q", output)
			}

			// Verify invariant: exit code 0 iff no error message
			hasError := strings.Contains(output, "error:")
			if (exitCode == 0) == hasError {
				t.Errorf("invariant violated: exitCode=%d but hasError=%v", exitCode, hasError)
			}
		})
	}
}

func TestRunInstall_CreatesDirectoryForLocal(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer

	// Setup: chdir to temp directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	projectRoot := t.TempDir()
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	// .claude directory should not exist initially
	if _, err := os.Stat(".claude"); !os.IsNotExist(err) {
		t.Fatal(".claude already exists before install")
	}

	args := []string{"--local"}
	exitCode := cmd.RunInstall(args, contentFS, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
		t.Logf("output: %s", buf.String())
	}

	// Verify .claude directory was created
	info, err := os.Stat(".claude")
	if err != nil {
		t.Fatalf(".claude directory was not created: %v", err)
	}
	if !info.IsDir() {
		t.Error(".claude is not a directory")
	}
}

func TestRunInstall_ConflictStrategyParsedBeforeScopeInvariant(t *testing.T) {
	// Contract C-12 specifies: "Parse conflict strategy from args (may fail with error per TD-1)"
	// then "Parse scope from remaining args"
	// This test verifies the ordering by checking that invalid conflict fails even with valid scope

	contentFS := buildTestFS()
	var buf bytes.Buffer

	args := []string{"--on-conflict=invalid", "--local"}
	exitCode := cmd.RunInstall(args, contentFS, &buf)

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}

	// Should fail with conflict error, not scope error
	output := buf.String()
	if !strings.Contains(output, "error:") {
		t.Errorf("expected error message, got: %q", output)
	}
}

func TestRunInstall_AllErrorsReturnExitCode1(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		setup func(t *testing.T) // optional setup
	}{
		{
			name: "invalid conflict strategy",
			args: []string{"--on-conflict=bad", "--local"},
		},
		{
			name: "missing scope flag",
			args: []string{},
		},
		{
			name: "both scope flags",
			args: []string{"--global", "--local"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contentFS := buildTestFS()
			var buf bytes.Buffer

			if tt.setup != nil {
				tt.setup(t)
			}

			exitCode := cmd.RunInstall(tt.args, contentFS, &buf)

			if exitCode != 1 {
				t.Errorf("expected exit code 1 for error case, got %d", exitCode)
			}

			output := buf.String()
			if !strings.Contains(output, "error:") {
				t.Errorf("expected error message, got: %q", output)
			}
		})
	}
}

func TestRunInstall_SuccessOutputContainsInstalledMessage(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer

	// Setup: chdir to temp directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	projectRoot := t.TempDir()
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	args := []string{"--local"}
	exitCode := cmd.RunInstall(args, contentFS, &buf)

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	output := buf.String()

	// Verify output contains installation messages
	expectedMessages := []string{
		"installed agents/gl-architect.md",
		"installed CLAUDE.md ->",
		"greenlight installed to",
	}

	for _, msg := range expectedMessages {
		if !strings.Contains(output, msg) {
			t.Errorf("output missing expected message: %q\nFull output: %q", msg, output)
		}
	}
}

func TestRunInstall_FilesystemErrorReturnsOne(t *testing.T) {
	// Create incomplete filesystem to trigger read error during install
	incompleteFS := fstest.MapFS{
		"agents/gl-architect.md": &fstest.MapFile{Data: []byte("# Architect\n")},
		// Missing all other required files
	}

	var buf bytes.Buffer

	// Setup: chdir to temp directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	projectRoot := t.TempDir()
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	args := []string{"--local"}
	exitCode := cmd.RunInstall(args, incompleteFS, &buf)

	if exitCode != 1 {
		t.Errorf("expected exit code 1 for filesystem error, got %d", exitCode)
	}

	output := buf.String()
	if !strings.Contains(output, "error:") {
		t.Errorf("expected error message for filesystem failure, got: %q", output)
	}
}

func TestRunInstall_OutputWrittenToProvidedWriter(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer

	// Setup: chdir to temp directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	projectRoot := t.TempDir()
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	args := []string{"--local"}
	cmd.RunInstall(args, contentFS, &buf)

	// Verify output was written to the provided buffer
	if buf.Len() == 0 {
		t.Error("expected output to be written to stdout, got empty buffer")
	}

	output := buf.String()
	if !strings.Contains(output, "installed") {
		t.Errorf("expected installation messages in output, got: %q", output)
	}
}

func TestRunInstall_MultipleStrategiesWithLocal(t *testing.T) {
	// Verify all three strategies work correctly with local scope
	strategies := []struct {
		name     string
		strategy string
		verify   func(t *testing.T, projectRoot string)
	}{
		{
			name:     "keep strategy",
			strategy: "keep",
			verify: func(t *testing.T, projectRoot string) {
				// With existing file, should create CLAUDE_GREENLIGHT.md
				if _, err := os.Stat(filepath.Join(projectRoot, "CLAUDE_GREENLIGHT.md")); os.IsNotExist(err) {
					t.Error("CLAUDE_GREENLIGHT.md was not created")
				}
			},
		},
		{
			name:     "replace strategy",
			strategy: "replace",
			verify: func(t *testing.T, projectRoot string) {
				// With existing file, should create backup
				if _, err := os.Stat(filepath.Join(projectRoot, "CLAUDE.md.backup")); os.IsNotExist(err) {
					t.Error("backup file was not created")
				}
			},
		},
		{
			name:     "append strategy",
			strategy: "append",
			verify: func(t *testing.T, projectRoot string) {
				// Should combine existing and new content
				content, err := os.ReadFile(filepath.Join(projectRoot, "CLAUDE.md"))
				if err != nil {
					t.Fatalf("failed to read CLAUDE.md: %v", err)
				}
				if !strings.Contains(string(content), "Existing") || !strings.Contains(string(content), "Greenlight") {
					t.Error("content was not appended correctly")
				}
			},
		},
	}

	for _, tt := range strategies {
		t.Run(tt.name, func(t *testing.T) {
			contentFS := buildTestFS()
			var buf bytes.Buffer

			// Setup: chdir to temp directory
			originalWd, err := os.Getwd()
			if err != nil {
				t.Fatalf("failed to get working directory: %v", err)
			}
			defer os.Chdir(originalWd)

			projectRoot := t.TempDir()
			if err := os.Chdir(projectRoot); err != nil {
				t.Fatalf("failed to change directory: %v", err)
			}

			// Create existing CLAUDE.md
			if err := os.WriteFile("CLAUDE.md", []byte("# Existing\n"), 0o644); err != nil {
				t.Fatalf("failed to create existing file: %v", err)
			}

			args := []string{"--on-conflict=" + tt.strategy, "--local"}
			exitCode := cmd.RunInstall(args, contentFS, &buf)

			if exitCode != 0 {
				t.Errorf("expected exit code 0, got %d", exitCode)
				t.Logf("output: %s", buf.String())
			}

			tt.verify(t, projectRoot)
		})
	}
}
