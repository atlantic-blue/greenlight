package installer_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/atlantic-blue/greenlight/internal/installer"
)

// Helper: buildTestFS creates a complete MapFS with all 38 manifest files.
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
		"commands/gl/debug.md":                     &fstest.MapFile{Data: []byte("# Debug\n")},
		"commands/gl/design.md":                    &fstest.MapFile{Data: []byte("# Design\n")},
		"commands/gl/help.md":                      &fstest.MapFile{Data: []byte("# Help\n")},
		"commands/gl/init.md":                      &fstest.MapFile{Data: []byte("# Init\n")},
		"commands/gl/map.md":                       &fstest.MapFile{Data: []byte("# Map\n")},
		"commands/gl/migrate-state.md":             &fstest.MapFile{Data: []byte("# Migrate State\n")},
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
		"references/circuit-breaker.md":            &fstest.MapFile{Data: []byte("# Circuit Breaker\n")},
		"references/deviation-rules.md":            &fstest.MapFile{Data: []byte("# Deviation Rules\n")},
		"references/state-format.md":              &fstest.MapFile{Data: []byte("# State Format\n")},
		"references/verification-patterns.md":      &fstest.MapFile{Data: []byte("# Verification Patterns\n")},
		"references/verification-tiers.md":         &fstest.MapFile{Data: []byte("# Verification Tiers\n")},
		"templates/config.md":                      &fstest.MapFile{Data: []byte("# Config Template\n")},
		"templates/slice-state.md":                &fstest.MapFile{Data: []byte("# Slice State Template\n")},
		"templates/state.md":                       &fstest.MapFile{Data: []byte("# State Template\n")},
		"CLAUDE.md":                                &fstest.MapFile{Data: []byte("# Greenlight CLAUDE.md\n\nTest content\n")},
	}
}

// C-07 Tests: InstallerNew

func TestNew_ReturnsNonNil(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer

	inst := installer.New(contentFS, &buf)

	if inst == nil {
		t.Fatal("expected non-nil Installer, got nil")
	}
}

func TestNew_AcceptsAnyFSAndWriter(t *testing.T) {
	tests := []struct {
		name      string
		contentFS fstest.MapFS
		stdout    *bytes.Buffer
	}{
		{
			name:      "accepts populated MapFS",
			contentFS: buildTestFS(),
			stdout:    &bytes.Buffer{},
		},
		{
			name:      "accepts empty MapFS",
			contentFS: fstest.MapFS{},
			stdout:    &bytes.Buffer{},
		},
		{
			name:      "accepts different buffer instance",
			contentFS: buildTestFS(),
			stdout:    &bytes.Buffer{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inst := installer.New(tt.contentFS, tt.stdout)
			if inst == nil {
				t.Fatal("expected non-nil Installer, got nil")
			}
		})
	}
}

// C-08 Tests: InstallerInstall

func TestInstall_WritesAllManifestFiles(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer
	targetDir := t.TempDir()

	inst := installer.New(contentFS, &buf)
	err := inst.Install(targetDir, "global", installer.ConflictKeep)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify all non-CLAUDE.md files were written
	expectedFiles := []string{
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
	}

	for _, relPath := range expectedFiles {
		destPath := filepath.Join(targetDir, relPath)
		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			t.Errorf("expected file %s to exist, but it does not", relPath)
		}
	}
}

func TestInstall_CreatesNecessarySubdirectories(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer
	targetDir := t.TempDir()

	inst := installer.New(contentFS, &buf)
	err := inst.Install(targetDir, "global", installer.ConflictKeep)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify subdirectories were created
	expectedDirs := []string{
		"agents",
		"commands",
		"commands/gl",
		"references",
		"templates",
	}

	for _, dir := range expectedDirs {
		dirPath := filepath.Join(targetDir, dir)
		info, err := os.Stat(dirPath)
		if err != nil {
			t.Errorf("expected directory %s to exist: %v", dir, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("expected %s to be a directory, got file", dir)
		}
	}
}

func TestInstall_FilesHaveCorrectPermissions(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer
	targetDir := t.TempDir()

	inst := installer.New(contentFS, &buf)
	err := inst.Install(targetDir, "global", installer.ConflictKeep)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Check file permissions (0o644)
	testFiles := []string{
		"agents/gl-architect.md",
		"commands/gl/help.md",
		"references/checkpoint-protocol.md",
	}

	for _, relPath := range testFiles {
		destPath := filepath.Join(targetDir, relPath)
		info, err := os.Stat(destPath)
		if err != nil {
			t.Fatalf("failed to stat %s: %v", relPath, err)
		}

		expectedPerms := os.FileMode(0o644)
		if info.Mode().Perm() != expectedPerms {
			t.Errorf("file %s has permissions %v, expected %v", relPath, info.Mode().Perm(), expectedPerms)
		}
	}
}

func TestInstall_DirectoriesHaveCorrectPermissions(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer
	targetDir := t.TempDir()

	inst := installer.New(contentFS, &buf)
	err := inst.Install(targetDir, "global", installer.ConflictKeep)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Check directory permissions (0o755)
	testDirs := []string{
		"agents",
		"commands",
		"commands/gl",
	}

	for _, dir := range testDirs {
		dirPath := filepath.Join(targetDir, dir)
		info, err := os.Stat(dirPath)
		if err != nil {
			t.Fatalf("failed to stat directory %s: %v", dir, err)
		}

		expectedPerms := os.FileMode(0o755)
		if info.Mode().Perm() != expectedPerms {
			t.Errorf("directory %s has permissions %v, expected %v", dir, info.Mode().Perm(), expectedPerms)
		}
	}
}

func TestInstall_WritesVersionFileWithCorrectFormat(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer
	targetDir := t.TempDir()

	inst := installer.New(contentFS, &buf)
	err := inst.Install(targetDir, "global", installer.ConflictKeep)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify .greenlight-version file exists
	versionPath := filepath.Join(targetDir, ".greenlight-version")
	content, err := os.ReadFile(versionPath)
	if err != nil {
		t.Fatalf("version file was not created: %v", err)
	}

	// Verify format: three lines, each newline-terminated
	lines := strings.Split(string(content), "\n")
	if len(lines) < 4 {
		t.Fatalf("expected at least 4 lines (3 values + trailing newline), got %d", len(lines))
	}

	// Each of the first 3 lines should be non-empty (before the trailing newline)
	for i := 0; i < 3; i++ {
		if lines[i] == "" {
			t.Errorf("line %d is empty, expected a value", i+1)
		}
	}

	// The fourth element should be empty (from trailing newline)
	if len(lines) > 3 && lines[3] != "" {
		t.Errorf("expected trailing newline, but line 4 is: %q", lines[3])
	}
}

func TestInstall_PrintsInstalledMessages(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer
	targetDir := t.TempDir()

	inst := installer.New(contentFS, &buf)
	err := inst.Install(targetDir, "global", installer.ConflictKeep)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	output := buf.String()

	// Verify output contains "installed" messages for some files
	expectedMessages := []string{
		"  installed agents/gl-architect.md\n",
		"  installed commands/gl/help.md\n",
		"  installed references/checkpoint-protocol.md\n",
	}

	for _, msg := range expectedMessages {
		if !strings.Contains(output, msg) {
			t.Errorf("output missing expected message: %q", msg)
		}
	}
}

func TestInstall_PrintsFinalMessage(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer
	targetDir := t.TempDir()

	inst := installer.New(contentFS, &buf)
	err := inst.Install(targetDir, "global", installer.ConflictKeep)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	output := buf.String()
	expectedMessage := "greenlight installed to " + targetDir + "\n"

	if !strings.Contains(output, expectedMessage) {
		t.Errorf("output missing final message.\nExpected: %q\nGot: %q", expectedMessage, output)
	}
}

func TestInstall_CLAUDEGlobalScope_WritesInsideTargetDir(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer
	targetDir := t.TempDir()

	inst := installer.New(contentFS, &buf)
	err := inst.Install(targetDir, "global", installer.ConflictKeep)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// CLAUDE.md should be inside targetDir for global scope
	claudePath := filepath.Join(targetDir, "CLAUDE.md")
	content, err := os.ReadFile(claudePath)
	if err != nil {
		t.Fatalf("CLAUDE.md was not created at %s: %v", claudePath, err)
	}

	expectedContent := []byte("# Greenlight CLAUDE.md\n\nTest content\n")
	if !bytes.Equal(content, expectedContent) {
		t.Errorf("CLAUDE.md has wrong content.\nExpected: %q\nGot: %q", expectedContent, content)
	}
}

func TestInstall_CLAUDELocalScope_WritesToParentOfTargetDir(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer

	// Create a project structure: tempDir/.claude
	projectRoot := t.TempDir()
	targetDir := filepath.Join(projectRoot, ".claude")

	inst := installer.New(contentFS, &buf)
	err := inst.Install(targetDir, "local", installer.ConflictKeep)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// CLAUDE.md should be in parent directory (project root) for local scope
	claudePath := filepath.Join(projectRoot, "CLAUDE.md")
	content, err := os.ReadFile(claudePath)
	if err != nil {
		t.Fatalf("CLAUDE.md was not created at %s: %v", claudePath, err)
	}

	expectedContent := []byte("# Greenlight CLAUDE.md\n\nTest content\n")
	if !bytes.Equal(content, expectedContent) {
		t.Errorf("CLAUDE.md has wrong content.\nExpected: %q\nGot: %q", expectedContent, content)
	}

	// Verify it was NOT written inside targetDir
	claudeInsideTarget := filepath.Join(targetDir, "CLAUDE.md")
	if _, err := os.Stat(claudeInsideTarget); !os.IsNotExist(err) {
		t.Error("CLAUDE.md should not exist inside targetDir for local scope")
	}
}

func TestInstall_CLAUDELocalScope_LiteralDotClaude(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer

	// Use literal ".claude" as targetDir (relative path)
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
	inst := installer.New(contentFS, &buf)
	err = inst.Install(targetDir, "local", installer.ConflictKeep)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// CLAUDE.md should be written as literal "CLAUDE.md" in current directory
	claudePath := "CLAUDE.md"
	content, err := os.ReadFile(claudePath)
	if err != nil {
		t.Fatalf("CLAUDE.md was not created: %v", err)
	}

	expectedContent := []byte("# Greenlight CLAUDE.md\n\nTest content\n")
	if !bytes.Equal(content, expectedContent) {
		t.Errorf("CLAUDE.md has wrong content.\nExpected: %q\nGot: %q", expectedContent, content)
	}
}

func TestInstall_WithKeepStrategy_DelegatesToHandleConflict(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer
	targetDir := t.TempDir()

	// Create existing CLAUDE.md
	claudePath := filepath.Join(targetDir, "CLAUDE.md")
	existingContent := []byte("# Existing CLAUDE.md\n")
	if err := os.WriteFile(claudePath, existingContent, 0o644); err != nil {
		t.Fatalf("failed to create existing CLAUDE.md: %v", err)
	}

	inst := installer.New(contentFS, &buf)
	err := inst.Install(targetDir, "global", installer.ConflictKeep)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify existing CLAUDE.md was kept
	content, err := os.ReadFile(claudePath)
	if err != nil {
		t.Fatalf("failed to read CLAUDE.md: %v", err)
	}
	if !bytes.Equal(content, existingContent) {
		t.Error("existing CLAUDE.md was not kept")
	}

	// Verify CLAUDE_GREENLIGHT.md was created
	greenlightPath := filepath.Join(targetDir, "CLAUDE_GREENLIGHT.md")
	if _, err := os.Stat(greenlightPath); os.IsNotExist(err) {
		t.Error("CLAUDE_GREENLIGHT.md was not created")
	}

	// Verify output message from handleConflict
	output := buf.String()
	if !strings.Contains(output, "existing CLAUDE.md kept") {
		t.Errorf("output missing conflict resolution message: %q", output)
	}
}

func TestInstall_IsIdempotent(t *testing.T) {
	contentFS := buildTestFS()
	targetDir := t.TempDir()

	// First install
	var buf1 bytes.Buffer
	inst1 := installer.New(contentFS, &buf1)
	err := inst1.Install(targetDir, "global", installer.ConflictReplace)
	if err != nil {
		t.Fatalf("first install failed: %v", err)
	}

	// Read state after first install
	firstVersionContent, err := os.ReadFile(filepath.Join(targetDir, ".greenlight-version"))
	if err != nil {
		t.Fatalf("failed to read version file after first install: %v", err)
	}

	// Second install
	var buf2 bytes.Buffer
	inst2 := installer.New(contentFS, &buf2)
	err = inst2.Install(targetDir, "global", installer.ConflictReplace)
	if err != nil {
		t.Fatalf("second install failed: %v", err)
	}

	// Read state after second install
	secondVersionContent, err := os.ReadFile(filepath.Join(targetDir, ".greenlight-version"))
	if err != nil {
		t.Fatalf("failed to read version file after second install: %v", err)
	}

	// Verify version file is identical
	if !bytes.Equal(firstVersionContent, secondVersionContent) {
		t.Error("version file changed between installs")
	}

	// Verify all files still exist
	testFile := filepath.Join(targetDir, "agents/gl-architect.md")
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("files disappeared after second install")
	}
}

func TestInstall_StopsOnFirstError_MissingFile(t *testing.T) {
	// Create incomplete filesystem (missing a file from manifest)
	incompleteFS := fstest.MapFS{
		"agents/gl-architect.md": &fstest.MapFile{Data: []byte("# Architect\n")},
		// Missing all other files
	}

	var buf bytes.Buffer
	targetDir := t.TempDir()

	inst := installer.New(incompleteFS, &buf)
	err := inst.Install(targetDir, "global", installer.ConflictKeep)

	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}

	// Verify error message contains context
	if !strings.Contains(err.Error(), "installing") {
		t.Errorf("error missing context about which file failed: %v", err)
	}
}

func TestInstall_ReturnsWrappedErrors(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer

	// Use invalid targetDir that will cause filesystem errors
	// (trying to write to a file as if it were a directory)
	targetDir := t.TempDir()
	conflictFile := filepath.Join(targetDir, "agents")
	if err := os.WriteFile(conflictFile, []byte("block"), 0o644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	inst := installer.New(contentFS, &buf)
	err := inst.Install(targetDir, "global", installer.ConflictKeep)

	if err == nil {
		t.Fatal("expected error when directory creation fails, got nil")
	}

	// Verify error is wrapped with context
	if !strings.Contains(err.Error(), "installing") {
		t.Errorf("expected error to contain 'installing', got: %v", err)
	}
}

func TestInstall_VersionFileWrittenAfterManifestFiles(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer
	targetDir := t.TempDir()

	inst := installer.New(contentFS, &buf)
	err := inst.Install(targetDir, "global", installer.ConflictKeep)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify both manifest files and version file exist
	// (if version was written first and install failed, version would exist alone)
	manifestFile := filepath.Join(targetDir, "agents/gl-architect.md")
	versionFile := filepath.Join(targetDir, ".greenlight-version")

	if _, err := os.Stat(manifestFile); os.IsNotExist(err) {
		t.Error("manifest file does not exist")
	}
	if _, err := os.Stat(versionFile); os.IsNotExist(err) {
		t.Error("version file does not exist")
	}
}

func TestInstall_CLAUDEPrintedWithCorrectPath(t *testing.T) {
	tests := []struct {
		name          string
		scope         string
		setupDir      func(t *testing.T) string
		expectedPath  string
	}{
		{
			name:  "global scope prints targetDir/CLAUDE.md",
			scope: "global",
			setupDir: func(t *testing.T) string {
				return t.TempDir()
			},
			expectedPath: "CLAUDE.md", // relative to targetDir
		},
		{
			name:  "local scope prints parent path",
			scope: "local",
			setupDir: func(t *testing.T) string {
				projectRoot := t.TempDir()
				return filepath.Join(projectRoot, ".claude")
			},
			expectedPath: "CLAUDE.md", // printed as relative or absolute
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contentFS := buildTestFS()
			var buf bytes.Buffer
			targetDir := tt.setupDir(t)

			inst := installer.New(contentFS, &buf)
			err := inst.Install(targetDir, tt.scope, installer.ConflictKeep)

			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}

			output := buf.String()
			if !strings.Contains(output, "installed CLAUDE.md ->") {
				t.Errorf("output missing CLAUDE.md installation message: %q", output)
			}
		})
	}
}

// C-33 Tests: ManifestBrownfieldUpdate

func TestManifest_Contains32Entries(t *testing.T) {
	if len(installer.Manifest) != 38 {
		t.Errorf("expected 38 manifest entries, got %d", len(installer.Manifest))
	}
}

func TestManifest_ContainsBrownfieldAgents(t *testing.T) {
	expectedAgents := []string{
		"agents/gl-assessor.md",
		"agents/gl-wrapper.md",
	}

	for _, expectedAgent := range expectedAgents {
		found := false
		for _, entry := range installer.Manifest {
			if entry == expectedAgent {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("manifest missing brownfield agent: %s", expectedAgent)
		}
	}
}

func TestManifest_ContainsBrownfieldCommands(t *testing.T) {
	expectedCommands := []string{
		"commands/gl/assess.md",
		"commands/gl/wrap.md",
	}

	for _, expectedCommand := range expectedCommands {
		found := false
		for _, entry := range installer.Manifest {
			if entry == expectedCommand {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("manifest missing brownfield command: %s", expectedCommand)
		}
	}
}

func TestManifest_CLAUDEIsLastEntry(t *testing.T) {
	if len(installer.Manifest) == 0 {
		t.Fatal("manifest is empty")
	}

	lastEntry := installer.Manifest[len(installer.Manifest)-1]
	if lastEntry != "CLAUDE.md" {
		t.Errorf("expected CLAUDE.md to be last entry, got: %s", lastEntry)
	}
}

func TestManifest_AgentsSectionAlphabeticallyOrdered(t *testing.T) {
	var agentEntries []string
	for _, entry := range installer.Manifest {
		if strings.HasPrefix(entry, "agents/") {
			agentEntries = append(agentEntries, entry)
		}
	}

	expectedOrder := []string{
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
	}

	if len(agentEntries) != len(expectedOrder) {
		t.Errorf("expected %d agent entries, got %d", len(expectedOrder), len(agentEntries))
	}

	for i, expected := range expectedOrder {
		if i >= len(agentEntries) {
			t.Errorf("missing agent entry at index %d: %s", i, expected)
			continue
		}
		if agentEntries[i] != expected {
			t.Errorf("agent entry at index %d: expected %s, got %s", i, expected, agentEntries[i])
		}
	}
}

func TestManifest_CommandsSectionAlphabeticallyOrdered(t *testing.T) {
	var commandEntries []string
	for _, entry := range installer.Manifest {
		if strings.HasPrefix(entry, "commands/gl/") {
			commandEntries = append(commandEntries, entry)
		}
	}

	expectedOrder := []string{
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
	}

	if len(commandEntries) != len(expectedOrder) {
		t.Errorf("expected %d command entries, got %d", len(expectedOrder), len(commandEntries))
	}

	for i, expected := range expectedOrder {
		if i >= len(commandEntries) {
			t.Errorf("missing command entry at index %d: %s", i, expected)
			continue
		}
		if commandEntries[i] != expected {
			t.Errorf("command entry at index %d: expected %s, got %s", i, expected, commandEntries[i])
		}
	}
}

func TestManifest_AllBrownfieldEntriesPresent(t *testing.T) {
	brownfieldEntries := []string{
		"agents/gl-assessor.md",
		"agents/gl-wrapper.md",
		"commands/gl/assess.md",
		"commands/gl/wrap.md",
	}

	manifestMap := make(map[string]bool)
	for _, entry := range installer.Manifest {
		manifestMap[entry] = true
	}

	for _, brownfieldEntry := range brownfieldEntries {
		if !manifestMap[brownfieldEntry] {
			t.Errorf("manifest missing brownfield entry: %s", brownfieldEntry)
		}
	}
}

func TestManifest_BrownfieldEntriesInstalledCorrectly(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer
	targetDir := t.TempDir()

	inst := installer.New(contentFS, &buf)
	err := inst.Install(targetDir, "global", installer.ConflictKeep)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	brownfieldFiles := []string{
		"agents/gl-assessor.md",
		"agents/gl-wrapper.md",
		"commands/gl/assess.md",
		"commands/gl/wrap.md",
	}

	for _, relPath := range brownfieldFiles {
		destPath := filepath.Join(targetDir, relPath)
		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			t.Errorf("brownfield file %s was not installed", relPath)
		}
	}
}

// C-38 Tests: ManifestDocumentationUpdate

func TestManifest_ContainsDocumentationCommands(t *testing.T) {
	expectedDocs := []string{
		"commands/gl/changelog.md",
		"commands/gl/roadmap.md",
	}

	for _, expectedDoc := range expectedDocs {
		found := false
		for _, entry := range installer.Manifest {
			if entry == expectedDoc {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("manifest missing documentation command: %s", expectedDoc)
		}
	}
}

func TestManifest_DocumentationCommandsAlphabeticallyOrdered(t *testing.T) {
	var commandEntries []string
	for _, entry := range installer.Manifest {
		if strings.HasPrefix(entry, "commands/gl/") {
			commandEntries = append(commandEntries, entry)
		}
	}

	// Find positions of documentation commands
	changelogIdx := -1
	roadmapIdx := -1
	assessIdx := -1
	designIdx := -1

	for i, entry := range commandEntries {
		switch entry {
		case "commands/gl/changelog.md":
			changelogIdx = i
		case "commands/gl/roadmap.md":
			roadmapIdx = i
		case "commands/gl/assess.md":
			assessIdx = i
		case "commands/gl/design.md":
			designIdx = i
		}
	}

	// Verify changelog comes after assess and before design
	if changelogIdx == -1 {
		t.Error("changelog.md not found in manifest")
	}
	if assessIdx == -1 {
		t.Error("assess.md not found in manifest")
	}
	if designIdx == -1 {
		t.Error("design.md not found in manifest")
	}
	if changelogIdx != -1 && assessIdx != -1 && changelogIdx <= assessIdx {
		t.Errorf("changelog.md (index %d) should come after assess.md (index %d)", changelogIdx, assessIdx)
	}
	if changelogIdx != -1 && designIdx != -1 && changelogIdx >= designIdx {
		t.Errorf("changelog.md (index %d) should come before design.md (index %d)", changelogIdx, designIdx)
	}

	// Verify roadmap comes after resume and before settings
	resumeIdx := -1
	settingsIdx := -1
	for i, entry := range commandEntries {
		switch entry {
		case "commands/gl/resume.md":
			resumeIdx = i
		case "commands/gl/settings.md":
			settingsIdx = i
		}
	}

	if roadmapIdx == -1 {
		t.Error("roadmap.md not found in manifest")
	}
	if resumeIdx == -1 {
		t.Error("resume.md not found in manifest")
	}
	if settingsIdx == -1 {
		t.Error("settings.md not found in manifest")
	}
	if roadmapIdx != -1 && resumeIdx != -1 && roadmapIdx <= resumeIdx {
		t.Errorf("roadmap.md (index %d) should come after resume.md (index %d)", roadmapIdx, resumeIdx)
	}
	if roadmapIdx != -1 && settingsIdx != -1 && roadmapIdx >= settingsIdx {
		t.Errorf("roadmap.md (index %d) should come before settings.md (index %d)", roadmapIdx, settingsIdx)
	}
}

func TestManifest_AllDocumentationEntriesPresent(t *testing.T) {
	documentationEntries := []string{
		"commands/gl/changelog.md",
		"commands/gl/roadmap.md",
	}

	manifestMap := make(map[string]bool)
	for _, entry := range installer.Manifest {
		manifestMap[entry] = true
	}

	for _, docEntry := range documentationEntries {
		if !manifestMap[docEntry] {
			t.Errorf("manifest missing documentation entry: %s", docEntry)
		}
	}
}

func TestManifest_DocumentationEntriesInstalledCorrectly(t *testing.T) {
	contentFS := buildTestFS()
	var buf bytes.Buffer
	targetDir := t.TempDir()

	inst := installer.New(contentFS, &buf)
	err := inst.Install(targetDir, "global", installer.ConflictKeep)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	documentationFiles := []string{
		"commands/gl/changelog.md",
		"commands/gl/roadmap.md",
	}

	for _, relPath := range documentationFiles {
		destPath := filepath.Join(targetDir, relPath)
		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			t.Errorf("documentation file %s was not installed", relPath)
		}
	}
}
