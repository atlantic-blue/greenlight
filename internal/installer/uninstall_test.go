package installer_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/atlantic-blue/greenlight/internal/installer"
)

// installAllFiles simulates a complete install by writing all manifest files
// + version file to targetDir, and CLAUDE.md to the appropriate location based on scope.
func installAllFiles(t *testing.T, targetDir, scope string) {
	t.Helper()

	// Create all manifest files except CLAUDE.md
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

	// Create CLAUDE.md in appropriate location based on scope
	var claudeMdPath string
	if scope == "global" {
		claudeMdPath = filepath.Join(targetDir, "CLAUDE.md")
	} else if scope == "local" {
		// For local scope, CLAUDE.md goes in parent of targetDir
		parent := filepath.Dir(targetDir)
		if targetDir == ".claude" {
			parent = "."
		}
		claudeMdPath = filepath.Join(parent, "CLAUDE.md")
	}

	if err := os.WriteFile(claudeMdPath, []byte("user content"), 0644); err != nil {
		t.Fatalf("writing CLAUDE.md: %v", err)
	}
}

// Test C-10.1: Removes all manifest files except CLAUDE.md
func TestUninstall_RemovesAllManifestFilesExceptClaudeMd(t *testing.T) {
	targetDir := t.TempDir()
	installAllFiles(t, targetDir, "global")

	var stdout bytes.Buffer
	err := installer.Uninstall(targetDir, "global", &stdout)
	if err != nil {
		t.Fatalf("Uninstall failed: %v", err)
	}

	// Verify all manifest files except CLAUDE.md are removed
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
	}

	for _, relPath := range manifestFiles {
		fullPath := filepath.Join(targetDir, relPath)
		if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
			t.Errorf("file %s should have been removed but still exists", relPath)
		}
	}
}

// Test C-10.2: CLAUDE.md is NEVER removed (NFR-3)
func TestUninstall_ClaudeMdNeverRemoved(t *testing.T) {
	targetDir := t.TempDir()
	installAllFiles(t, targetDir, "global")

	var stdout bytes.Buffer
	err := installer.Uninstall(targetDir, "global", &stdout)
	if err != nil {
		t.Fatalf("Uninstall failed: %v", err)
	}

	claudeMdPath := filepath.Join(targetDir, "CLAUDE.md")
	data, err := os.ReadFile(claudeMdPath)
	if err != nil {
		t.Fatalf("CLAUDE.md should exist after uninstall: %v", err)
	}

	if string(data) != "user content" {
		t.Errorf("CLAUDE.md content changed, got %q, want %q", string(data), "user content")
	}
}

// Test C-10.3: Removes .greenlight-version file
func TestUninstall_RemovesVersionFile(t *testing.T) {
	targetDir := t.TempDir()
	installAllFiles(t, targetDir, "global")

	var stdout bytes.Buffer
	err := installer.Uninstall(targetDir, "global", &stdout)
	if err != nil {
		t.Fatalf("Uninstall failed: %v", err)
	}

	versionPath := filepath.Join(targetDir, ".greenlight-version")
	if _, err := os.Stat(versionPath); !os.IsNotExist(err) {
		t.Errorf(".greenlight-version should have been removed")
	}
}

// Test C-10.4: Missing files are skipped without error (idempotent)
func TestUninstall_MissingFilesSkippedWithoutError(t *testing.T) {
	targetDir := t.TempDir()

	// Create only a subset of files
	someFiles := []string{
		"agents/gl-architect.md",
		"commands/gl/help.md",
		"references/checkpoint-protocol.md",
	}

	for _, relPath := range someFiles {
		fullPath := filepath.Join(targetDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("creating directory for %s: %v", relPath, err)
		}
		if err := os.WriteFile(fullPath, []byte("test content"), 0644); err != nil {
			t.Fatalf("writing %s: %v", relPath, err)
		}
	}

	var stdout bytes.Buffer
	err := installer.Uninstall(targetDir, "global", &stdout)
	if err != nil {
		t.Fatalf("Uninstall should not error on missing files: %v", err)
	}

	// Verify the files that existed were removed
	for _, relPath := range someFiles {
		fullPath := filepath.Join(targetDir, relPath)
		if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
			t.Errorf("file %s should have been removed", relPath)
		}
	}
}

// Test C-10.5: Uninstall is idempotent (running twice produces same result)
func TestUninstall_Idempotent(t *testing.T) {
	targetDir := t.TempDir()
	installAllFiles(t, targetDir, "global")

	// First uninstall
	var stdout1 bytes.Buffer
	err := installer.Uninstall(targetDir, "global", &stdout1)
	if err != nil {
		t.Fatalf("First uninstall failed: %v", err)
	}

	// Second uninstall
	var stdout2 bytes.Buffer
	err = installer.Uninstall(targetDir, "global", &stdout2)
	if err != nil {
		t.Fatalf("Second uninstall should succeed (idempotent): %v", err)
	}

	// CLAUDE.md should still exist
	claudeMdPath := filepath.Join(targetDir, "CLAUDE.md")
	if _, err := os.Stat(claudeMdPath); err != nil {
		t.Errorf("CLAUDE.md should still exist after second uninstall: %v", err)
	}
}

// Test C-10.6: Prints "removed <relPath>" for each removed file
func TestUninstall_PrintsRemovedForEachFile(t *testing.T) {
	targetDir := t.TempDir()
	installAllFiles(t, targetDir, "global")

	var stdout bytes.Buffer
	err := installer.Uninstall(targetDir, "global", &stdout)
	if err != nil {
		t.Fatalf("Uninstall failed: %v", err)
	}

	output := stdout.String()

	// Check for some representative files
	expectedRemovals := []string{
		"  removed agents/gl-architect.md\n",
		"  removed commands/gl/help.md\n",
		"  removed references/checkpoint-protocol.md\n",
		"  removed templates/config.md\n",
	}

	for _, expected := range expectedRemovals {
		if !strings.Contains(output, expected) {
			t.Errorf("output missing %q", expected)
		}
	}

	// CLAUDE.md should NOT be in removed list
	if strings.Contains(output, "removed CLAUDE.md\n") && !strings.Contains(output, "removed CLAUDE_GREENLIGHT.md\n") {
		t.Errorf("output should not contain 'removed CLAUDE.md'")
	}
}

// Test C-10.7: Prints "greenlight uninstalled from <targetDir>" at end
func TestUninstall_PrintsUninstallMessage(t *testing.T) {
	targetDir := t.TempDir()
	installAllFiles(t, targetDir, "global")

	var stdout bytes.Buffer
	err := installer.Uninstall(targetDir, "global", &stdout)
	if err != nil {
		t.Fatalf("Uninstall failed: %v", err)
	}

	output := stdout.String()
	expectedMsg := "greenlight uninstalled from " + targetDir + "\n"
	if !strings.Contains(output, expectedMsg) {
		t.Errorf("output missing final message, got:\n%s\nwant substring: %s", output, expectedMsg)
	}
}

// Test C-10.8: Empty directories are cleaned up deepest-first
func TestUninstall_CleansUpEmptyDirectories(t *testing.T) {
	targetDir := t.TempDir()
	installAllFiles(t, targetDir, "global")

	var stdout bytes.Buffer
	err := installer.Uninstall(targetDir, "global", &stdout)
	if err != nil {
		t.Fatalf("Uninstall failed: %v", err)
	}

	// Verify directories are cleaned up
	dirsToCleanup := []string{
		filepath.Join(targetDir, "commands", "gl"),
		filepath.Join(targetDir, "commands"),
		filepath.Join(targetDir, "agents"),
		filepath.Join(targetDir, "references"),
		filepath.Join(targetDir, "templates"),
	}

	for _, dir := range dirsToCleanup {
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			t.Errorf("empty directory %s should have been cleaned up", dir)
		}
	}
}

// Test C-10.9: Non-empty directories are NOT cleaned up
func TestUninstall_DoesNotRemoveNonEmptyDirectories(t *testing.T) {
	targetDir := t.TempDir()
	installAllFiles(t, targetDir, "global")

	// Add a user file to agents directory
	userFile := filepath.Join(targetDir, "agents", "my-custom-agent.md")
	if err := os.WriteFile(userFile, []byte("user content"), 0644); err != nil {
		t.Fatalf("writing user file: %v", err)
	}

	var stdout bytes.Buffer
	err := installer.Uninstall(targetDir, "global", &stdout)
	if err != nil {
		t.Fatalf("Uninstall failed: %v", err)
	}

	// Verify agents directory still exists because of user file
	agentsDir := filepath.Join(targetDir, "agents")
	if _, err := os.Stat(agentsDir); os.IsNotExist(err) {
		t.Errorf("non-empty agents directory should not have been removed")
	}

	// Verify user file still exists
	if _, err := os.Stat(userFile); err != nil {
		t.Errorf("user file should still exist: %v", err)
	}
}

// Test C-10.10: Removes CLAUDE_GREENLIGHT.md conflict artifact (global scope)
func TestUninstall_RemovesClaudeGreenlightMd_GlobalScope(t *testing.T) {
	targetDir := t.TempDir()
	installAllFiles(t, targetDir, "global")

	// Create conflict artifact for global scope (in targetDir)
	conflictPath := filepath.Join(targetDir, "CLAUDE_GREENLIGHT.md")
	if err := os.WriteFile(conflictPath, []byte("greenlight content"), 0644); err != nil {
		t.Fatalf("writing conflict artifact: %v", err)
	}

	var stdout bytes.Buffer
	err := installer.Uninstall(targetDir, "global", &stdout)
	if err != nil {
		t.Fatalf("Uninstall failed: %v", err)
	}

	if _, err := os.Stat(conflictPath); !os.IsNotExist(err) {
		t.Errorf("CLAUDE_GREENLIGHT.md should have been removed")
	}
}

// Test C-10.11: Removes CLAUDE.md.backup conflict artifact (global scope)
func TestUninstall_RemovesClaudeMdBackup_GlobalScope(t *testing.T) {
	targetDir := t.TempDir()
	installAllFiles(t, targetDir, "global")

	// Create backup artifact for global scope (in targetDir)
	backupPath := filepath.Join(targetDir, "CLAUDE.md.backup")
	if err := os.WriteFile(backupPath, []byte("backup content"), 0644); err != nil {
		t.Fatalf("writing backup artifact: %v", err)
	}

	var stdout bytes.Buffer
	err := installer.Uninstall(targetDir, "global", &stdout)
	if err != nil {
		t.Fatalf("Uninstall failed: %v", err)
	}

	if _, err := os.Stat(backupPath); !os.IsNotExist(err) {
		t.Errorf("CLAUDE.md.backup should have been removed")
	}
}

// Test C-10.12: Removes conflict artifacts in correct location for local scope
func TestUninstall_RemovesConflictArtifacts_LocalScope(t *testing.T) {
	// Create temp dir structure: tempDir/.claude
	tempRoot := t.TempDir()
	targetDir := filepath.Join(tempRoot, ".claude")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		t.Fatalf("creating .claude dir: %v", err)
	}

	installAllFiles(t, targetDir, "local")

	// For local scope, conflict artifacts are in parent of targetDir (tempRoot)
	greenlightPath := filepath.Join(tempRoot, "CLAUDE_GREENLIGHT.md")
	backupPath := filepath.Join(tempRoot, "CLAUDE.md.backup")

	if err := os.WriteFile(greenlightPath, []byte("greenlight content"), 0644); err != nil {
		t.Fatalf("writing CLAUDE_GREENLIGHT.md: %v", err)
	}
	if err := os.WriteFile(backupPath, []byte("backup content"), 0644); err != nil {
		t.Fatalf("writing CLAUDE.md.backup: %v", err)
	}

	var stdout bytes.Buffer
	err := installer.Uninstall(targetDir, "local", &stdout)
	if err != nil {
		t.Fatalf("Uninstall failed: %v", err)
	}

	// Verify both artifacts removed from parent directory
	if _, err := os.Stat(greenlightPath); !os.IsNotExist(err) {
		t.Errorf("CLAUDE_GREENLIGHT.md should have been removed from parent dir")
	}
	if _, err := os.Stat(backupPath); !os.IsNotExist(err) {
		t.Errorf("CLAUDE.md.backup should have been removed from parent dir")
	}
}

// Test C-10.13: Missing conflict artifacts are skipped without error
func TestUninstall_MissingConflictArtifactsSkipped(t *testing.T) {
	targetDir := t.TempDir()
	installAllFiles(t, targetDir, "global")

	// Do NOT create conflict artifacts - they're missing

	var stdout bytes.Buffer
	err := installer.Uninstall(targetDir, "global", &stdout)
	if err != nil {
		t.Fatalf("Uninstall should not error on missing conflict artifacts: %v", err)
	}
}

// Test C-10.14: Returns wrapped error on filesystem failure
func TestUninstall_ReturnsErrorOnFilesystemFailure(t *testing.T) {
	targetDir := t.TempDir()

	// Create a file we'll make read-only directory around
	agentFile := filepath.Join(targetDir, "agents", "gl-architect.md")
	if err := os.MkdirAll(filepath.Dir(agentFile), 0755); err != nil {
		t.Fatalf("creating directory: %v", err)
	}
	if err := os.WriteFile(agentFile, []byte("test"), 0644); err != nil {
		t.Fatalf("writing file: %v", err)
	}

	// Make parent directory read-only to cause removal failure
	agentsDir := filepath.Join(targetDir, "agents")
	if err := os.Chmod(agentsDir, 0555); err != nil {
		t.Fatalf("making directory read-only: %v", err)
	}
	defer os.Chmod(agentsDir, 0755) // cleanup

	var stdout bytes.Buffer
	err := installer.Uninstall(targetDir, "global", &stdout)
	if err == nil {
		t.Fatalf("Uninstall should return error when file removal fails")
	}

	// Error should mention the file path
	if !strings.Contains(err.Error(), "gl-architect.md") && !strings.Contains(err.Error(), "removing") {
		t.Errorf("error should mention file or removal operation, got: %v", err)
	}
}

// Test C-10.15: Prints "removed CLAUDE_GREENLIGHT.md" when artifact removed
func TestUninstall_PrintsRemovedClaudeGreenlight(t *testing.T) {
	targetDir := t.TempDir()
	installAllFiles(t, targetDir, "global")

	conflictPath := filepath.Join(targetDir, "CLAUDE_GREENLIGHT.md")
	if err := os.WriteFile(conflictPath, []byte("greenlight content"), 0644); err != nil {
		t.Fatalf("writing conflict artifact: %v", err)
	}

	var stdout bytes.Buffer
	err := installer.Uninstall(targetDir, "global", &stdout)
	if err != nil {
		t.Fatalf("Uninstall failed: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "  removed CLAUDE_GREENLIGHT.md\n") {
		t.Errorf("output should contain 'removed CLAUDE_GREENLIGHT.md', got:\n%s", output)
	}
}

// Test C-10.16: Prints "removed CLAUDE.md.backup" when artifact removed
func TestUninstall_PrintsRemovedClaudeBackup(t *testing.T) {
	targetDir := t.TempDir()
	installAllFiles(t, targetDir, "global")

	backupPath := filepath.Join(targetDir, "CLAUDE.md.backup")
	if err := os.WriteFile(backupPath, []byte("backup content"), 0644); err != nil {
		t.Fatalf("writing backup artifact: %v", err)
	}

	var stdout bytes.Buffer
	err := installer.Uninstall(targetDir, "global", &stdout)
	if err != nil {
		t.Fatalf("Uninstall failed: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "  removed CLAUDE.md.backup\n") {
		t.Errorf("output should contain 'removed CLAUDE.md.backup', got:\n%s", output)
	}
}
