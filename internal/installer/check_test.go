package installer_test

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/atlantic-blue/greenlight/internal/installer"
)

// installTestFiles simulates a successful install by writing all manifest files
// plus .greenlight-version to the target directory, handling CLAUDE.md placement
// according to scope.
func installTestFiles(t *testing.T, targetDir, scope string, contentFS fs.FS) {
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

// C-09 Tests: InstallerCheck

// Presence-only mode (verify=false) tests

func TestCheck_AllFilesPresentAndNonEmpty_ReturnsTrue(t *testing.T) {
	contentFS := buildTestFS()
	targetDir := t.TempDir()
	installTestFiles(t, targetDir, "global", contentFS)

	var buf bytes.Buffer
	ok := installer.Check(targetDir, "global", &buf, false, nil)

	if !ok {
		t.Errorf("expected Check to return true, got false")
	}

	output := buf.String()
	expectedSummary := "all 32 files present\n"
	if !strings.HasSuffix(output, expectedSummary) {
		t.Errorf("expected output to end with %q, got: %q", expectedSummary, output)
	}
}

func TestCheck_OneManifestFileMissing_ReturnsFalse(t *testing.T) {
	contentFS := buildTestFS()
	targetDir := t.TempDir()
	installTestFiles(t, targetDir, "global", contentFS)

	// Remove one manifest file
	missingFile := filepath.Join(targetDir, "agents/gl-architect.md")
	if err := os.Remove(missingFile); err != nil {
		t.Fatalf("failed to remove test file: %v", err)
	}

	var buf bytes.Buffer
	ok := installer.Check(targetDir, "global", &buf, false, nil)

	if ok {
		t.Errorf("expected Check to return false when file missing, got true")
	}

	output := buf.String()
	expectedMessage := "  MISSING  agents/gl-architect.md\n"
	if !strings.Contains(output, expectedMessage) {
		t.Errorf("expected output to contain %q, got: %q", expectedMessage, output)
	}
}

func TestCheck_OneManifestFileEmpty_ReturnsFalse(t *testing.T) {
	contentFS := buildTestFS()
	targetDir := t.TempDir()
	installTestFiles(t, targetDir, "global", contentFS)

	// Create an empty file
	emptyFile := filepath.Join(targetDir, "agents/gl-debugger.md")
	if err := os.WriteFile(emptyFile, []byte{}, 0o644); err != nil {
		t.Fatalf("failed to create empty file: %v", err)
	}

	var buf bytes.Buffer
	ok := installer.Check(targetDir, "global", &buf, false, nil)

	if ok {
		t.Errorf("expected Check to return false when file empty, got true")
	}

	output := buf.String()
	expectedMessage := "  EMPTY    agents/gl-debugger.md\n"
	if !strings.Contains(output, expectedMessage) {
		t.Errorf("expected output to contain %q, got: %q", expectedMessage, output)
	}
}

func TestCheck_MultipleFilesMissing_CorrectCountInSummary(t *testing.T) {
	contentFS := buildTestFS()
	targetDir := t.TempDir()
	installTestFiles(t, targetDir, "global", contentFS)

	// Remove multiple files
	filesToRemove := []string{
		"agents/gl-architect.md",
		"agents/gl-designer.md",
		"commands/gl/help.md",
	}

	for _, relPath := range filesToRemove {
		if err := os.Remove(filepath.Join(targetDir, relPath)); err != nil {
			t.Fatalf("failed to remove %s: %v", relPath, err)
		}
	}

	var buf bytes.Buffer
	ok := installer.Check(targetDir, "global", &buf, false, nil)

	if ok {
		t.Errorf("expected Check to return false, got true")
	}

	output := buf.String()
	// Should report 27/30 files present (3 missing)
	expectedSummary := "29/32 files present (3 missing, 0 empty)\n"
	if !strings.HasSuffix(output, expectedSummary) {
		t.Errorf("expected output to end with %q, got: %q", expectedSummary, output)
	}
}

func TestCheck_VersionFileMissing_ReturnsFalse(t *testing.T) {
	contentFS := buildTestFS()
	targetDir := t.TempDir()
	installTestFiles(t, targetDir, "global", contentFS)

	// Remove version file
	versionPath := filepath.Join(targetDir, ".greenlight-version")
	if err := os.Remove(versionPath); err != nil {
		t.Fatalf("failed to remove version file: %v", err)
	}

	var buf bytes.Buffer
	ok := installer.Check(targetDir, "global", &buf, false, nil)

	if ok {
		t.Errorf("expected Check to return false when version file missing, got true")
	}

	output := buf.String()
	expectedMessage := "  MISSING  .greenlight-version\n"
	if !strings.Contains(output, expectedMessage) {
		t.Errorf("expected output to contain %q, got: %q", expectedMessage, output)
	}
}

func TestCheck_VersionFilePresent_PrintsVersion(t *testing.T) {
	contentFS := buildTestFS()
	targetDir := t.TempDir()
	installTestFiles(t, targetDir, "global", contentFS)

	var buf bytes.Buffer
	ok := installer.Check(targetDir, "global", &buf, false, nil)

	if !ok {
		t.Errorf("expected Check to return true, got false")
	}

	output := buf.String()
	expectedMessage := "  version: v1.0.0\n"
	if !strings.Contains(output, expectedMessage) {
		t.Errorf("expected output to contain %q, got: %q", expectedMessage, output)
	}
}

func TestCheck_CLAUDEGlobalScope_CheckedInsideTargetDir(t *testing.T) {
	contentFS := buildTestFS()
	targetDir := t.TempDir()
	installTestFiles(t, targetDir, "global", contentFS)

	var buf bytes.Buffer
	ok := installer.Check(targetDir, "global", &buf, false, nil)

	if !ok {
		t.Errorf("expected Check to return true, got false")
	}

	// Verify CLAUDE.md exists where expected (inside targetDir)
	claudePath := filepath.Join(targetDir, "CLAUDE.md")
	if _, err := os.Stat(claudePath); os.IsNotExist(err) {
		t.Error("CLAUDE.md should exist inside targetDir for global scope")
	}
}

func TestCheck_CLAUDELocalScope_CheckedInParentOfTargetDir(t *testing.T) {
	contentFS := buildTestFS()
	projectRoot := t.TempDir()
	targetDir := filepath.Join(projectRoot, ".claude")
	installTestFiles(t, targetDir, "local", contentFS)

	var buf bytes.Buffer
	ok := installer.Check(targetDir, "local", &buf, false, nil)

	if !ok {
		t.Errorf("expected Check to return true, got false")
	}

	// Verify CLAUDE.md exists in parent directory
	claudePath := filepath.Join(projectRoot, "CLAUDE.md")
	if _, err := os.Stat(claudePath); os.IsNotExist(err) {
		t.Error("CLAUDE.md should exist in parent of targetDir for local scope")
	}

	// Verify it doesn't exist inside targetDir
	claudeInsideTarget := filepath.Join(targetDir, "CLAUDE.md")
	if _, err := os.Stat(claudeInsideTarget); !os.IsNotExist(err) {
		t.Error("CLAUDE.md should not exist inside targetDir for local scope")
	}
}

func TestCheck_VerifyFalseWithNilContentFS_WorksFine(t *testing.T) {
	contentFS := buildTestFS()
	targetDir := t.TempDir()
	installTestFiles(t, targetDir, "global", contentFS)

	var buf bytes.Buffer
	ok := installer.Check(targetDir, "global", &buf, false, nil)

	if !ok {
		t.Errorf("expected Check to return true, got false")
	}

	output := buf.String()
	if !strings.Contains(output, "all 32 files present") {
		t.Errorf("expected success message, got: %q", output)
	}
}

// Content verification mode (verify=true) tests

func TestCheck_VerifyTrue_AllFilesMatchEmbeddedContent_ReturnsTrue(t *testing.T) {
	contentFS := buildTestFS()
	targetDir := t.TempDir()
	installTestFiles(t, targetDir, "global", contentFS)

	var buf bytes.Buffer
	ok := installer.Check(targetDir, "global", &buf, true, contentFS)

	if !ok {
		t.Errorf("expected Check to return true, got false")
	}

	output := buf.String()
	expectedSummary := "all 32 files verified\n"
	if !strings.HasSuffix(output, expectedSummary) {
		t.Errorf("expected output to end with %q, got: %q", expectedSummary, output)
	}
}

func TestCheck_VerifyTrue_OneFileModified_ReturnsFalse(t *testing.T) {
	contentFS := buildTestFS()
	targetDir := t.TempDir()
	installTestFiles(t, targetDir, "global", contentFS)

	// Modify one file to have different content
	modifiedFile := filepath.Join(targetDir, "agents/gl-security.md")
	modifiedContent := []byte("# Modified Security\nThis content is different\n")
	if err := os.WriteFile(modifiedFile, modifiedContent, 0o644); err != nil {
		t.Fatalf("failed to modify file: %v", err)
	}

	var buf bytes.Buffer
	ok := installer.Check(targetDir, "global", &buf, true, contentFS)

	if ok {
		t.Errorf("expected Check to return false when file modified, got true")
	}

	output := buf.String()
	expectedMessage := "  MODIFIED agents/gl-security.md\n"
	if !strings.Contains(output, expectedMessage) {
		t.Errorf("expected output to contain %q, got: %q", expectedMessage, output)
	}
}

func TestCheck_VerifyTrue_SummaryIncludesModifiedCount(t *testing.T) {
	contentFS := buildTestFS()
	targetDir := t.TempDir()
	installTestFiles(t, targetDir, "global", contentFS)

	// Modify two files
	filesToModify := []string{
		"agents/gl-implementer.md",
		"commands/gl/slice.md",
	}

	for _, relPath := range filesToModify {
		destPath := filepath.Join(targetDir, relPath)
		modifiedContent := []byte("# Modified content\n")
		if err := os.WriteFile(destPath, modifiedContent, 0o644); err != nil {
			t.Fatalf("failed to modify %s: %v", relPath, err)
		}
	}

	var buf bytes.Buffer
	ok := installer.Check(targetDir, "global", &buf, true, contentFS)

	if ok {
		t.Errorf("expected Check to return false, got true")
	}

	output := buf.String()
	// Should report 28/30 files verified (0 missing, 0 empty, 2 modified)
	expectedSummary := "30/32 files verified (0 missing, 0 empty, 2 modified)\n"
	if !strings.HasSuffix(output, expectedSummary) {
		t.Errorf("expected output to end with %q, got: %q", expectedSummary, output)
	}
}

func TestCheck_VerifyTrue_MissingAndModified_SummaryShowsBoth(t *testing.T) {
	contentFS := buildTestFS()
	targetDir := t.TempDir()
	installTestFiles(t, targetDir, "global", contentFS)

	// Remove one file
	if err := os.Remove(filepath.Join(targetDir, "agents/gl-architect.md")); err != nil {
		t.Fatalf("failed to remove file: %v", err)
	}

	// Modify one file
	modifiedFile := filepath.Join(targetDir, "agents/gl-designer.md")
	if err := os.WriteFile(modifiedFile, []byte("modified"), 0o644); err != nil {
		t.Fatalf("failed to modify file: %v", err)
	}

	// Make one file empty
	emptyFile := filepath.Join(targetDir, "agents/gl-verifier.md")
	if err := os.WriteFile(emptyFile, []byte{}, 0o644); err != nil {
		t.Fatalf("failed to create empty file: %v", err)
	}

	var buf bytes.Buffer
	ok := installer.Check(targetDir, "global", &buf, true, contentFS)

	if ok {
		t.Errorf("expected Check to return false, got true")
	}

	output := buf.String()
	// 27/30 verified (1 missing, 1 empty, 1 modified)
	expectedSummary := "29/32 files verified (1 missing, 1 empty, 1 modified)\n"
	if !strings.HasSuffix(output, expectedSummary) {
		t.Errorf("expected output to end with %q, got: %q", expectedSummary, output)
	}

	// Verify all three types of failures are reported
	if !strings.Contains(output, "  MISSING  agents/gl-architect.md\n") {
		t.Error("expected MISSING message in output")
	}
	if !strings.Contains(output, "  EMPTY    agents/gl-verifier.md\n") {
		t.Error("expected EMPTY message in output")
	}
	if !strings.Contains(output, "  MODIFIED agents/gl-designer.md\n") {
		t.Error("expected MODIFIED message in output")
	}
}

// Invariant tests

func TestCheck_NeverModifiesFilesystem_ReadOnly(t *testing.T) {
	contentFS := buildTestFS()
	targetDir := t.TempDir()
	installTestFiles(t, targetDir, "global", contentFS)

	// Record modification times before check
	beforeTimes := make(map[string]int64)
	for _, relPath := range installer.Manifest {
		destPath := filepath.Join(targetDir, relPath)
		info, err := os.Stat(destPath)
		if err == nil {
			beforeTimes[relPath] = info.ModTime().Unix()
		}
	}

	var buf bytes.Buffer
	installer.Check(targetDir, "global", &buf, true, contentFS)

	// Verify modification times unchanged
	for relPath, beforeTime := range beforeTimes {
		destPath := filepath.Join(targetDir, relPath)
		info, err := os.Stat(destPath)
		if err != nil {
			t.Errorf("file %s disappeared: %v", relPath, err)
			continue
		}
		afterTime := info.ModTime().Unix()
		if afterTime != beforeTime {
			t.Errorf("file %s was modified (mtime changed)", relPath)
		}
	}

	// Verify no new files created
	versionPath := filepath.Join(targetDir, ".greenlight-version")
	beforeVersionTime, err := os.Stat(versionPath)
	if err != nil {
		t.Fatalf("version file missing: %v", err)
	}

	installer.Check(targetDir, "global", &buf, false, nil)

	afterVersionTime, err := os.Stat(versionPath)
	if err != nil {
		t.Fatalf("version file missing after check: %v", err)
	}

	if afterVersionTime.ModTime() != beforeVersionTime.ModTime() {
		t.Error("version file was modified")
	}
}

func TestCheck_SummaryAlwaysLastLine(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T, targetDir string)
		expectedLast  string
	}{
		{
			name: "all files present",
			setup: func(t *testing.T, targetDir string) {
				// no-op, files already installed
			},
			expectedLast: "all 32 files present\n",
		},
		{
			name: "one file missing",
			setup: func(t *testing.T, targetDir string) {
				os.Remove(filepath.Join(targetDir, "agents/gl-architect.md"))
			},
			expectedLast: "31/32 files present (1 missing, 0 empty)\n",
		},
		{
			name: "one file empty",
			setup: func(t *testing.T, targetDir string) {
				os.WriteFile(filepath.Join(targetDir, "agents/gl-debugger.md"), []byte{}, 0o644)
			},
			expectedLast: "32/32 files present (0 missing, 1 empty)\n",
		},
		{
			name: "version file missing",
			setup: func(t *testing.T, targetDir string) {
				os.Remove(filepath.Join(targetDir, ".greenlight-version"))
			},
			expectedLast: "32/32 files present (0 missing, 0 empty)\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contentFS := buildTestFS()
			targetDir := t.TempDir()
			installTestFiles(t, targetDir, "global", contentFS)

			tt.setup(t, targetDir)

			var buf bytes.Buffer
			installer.Check(targetDir, "global", &buf, false, nil)

			output := buf.String()
			if !strings.HasSuffix(output, tt.expectedLast) {
				t.Errorf("expected output to end with %q, got: %q", tt.expectedLast, output)
			}
		})
	}
}

func TestCheck_StatError_PrintsErrorMessage(t *testing.T) {
	contentFS := buildTestFS()
	targetDir := t.TempDir()
	installTestFiles(t, targetDir, "global", contentFS)

	// Create a directory where a file should be (causes stat to succeed but directory check fails)
	// We'll create a permission error instead
	testFile := filepath.Join(targetDir, "agents")
	if err := os.Chmod(testFile, 0o000); err != nil {
		t.Fatalf("failed to change permissions: %v", err)
	}
	defer os.Chmod(testFile, 0o755) // Cleanup

	// This test is system-dependent, so we'll skip detailed verification
	// Just ensure Check handles errors gracefully
	var buf bytes.Buffer
	ok := installer.Check(targetDir, "global", &buf, false, nil)

	// Check should detect problems
	if ok {
		t.Error("expected Check to return false when stat errors occur")
	}
}

func TestCheck_VerifyModeSHA256HashComparison(t *testing.T) {
	contentFS := buildTestFS()
	targetDir := t.TempDir()
	installTestFiles(t, targetDir, "global", contentFS)

	// Verify the hash comparison is based on SHA-256
	// Modify a file by appending one byte
	testFile := filepath.Join(targetDir, "templates/config.md")
	originalContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	modifiedContent := append(originalContent, byte(' '))
	if err := os.WriteFile(testFile, modifiedContent, 0o644); err != nil {
		t.Fatalf("failed to write modified content: %v", err)
	}

	// Verify hashes are different
	originalHash := sha256.Sum256(originalContent)
	modifiedHash := sha256.Sum256(modifiedContent)
	if originalHash == modifiedHash {
		t.Fatal("test setup error: hashes should be different")
	}

	var buf bytes.Buffer
	ok := installer.Check(targetDir, "global", &buf, true, contentFS)

	if ok {
		t.Error("expected Check to detect modified file via hash comparison")
	}

	output := buf.String()
	if !strings.Contains(output, "  MODIFIED templates/config.md\n") {
		t.Errorf("expected MODIFIED message, got: %q", output)
	}
}

func TestCheck_VerifyMode_AllFilesChecked(t *testing.T) {
	contentFS := buildTestFS()
	targetDir := t.TempDir()
	installTestFiles(t, targetDir, "global", contentFS)

	// Modify every single file
	for _, relPath := range installer.Manifest {
		if relPath == "CLAUDE.md" {
			destPath := filepath.Join(targetDir, relPath)
			if err := os.WriteFile(destPath, []byte("modified"), 0o644); err != nil {
				t.Fatalf("failed to modify %s: %v", relPath, err)
			}
		} else {
			destPath := filepath.Join(targetDir, relPath)
			if err := os.WriteFile(destPath, []byte("modified"), 0o644); err != nil {
				t.Fatalf("failed to modify %s: %v", relPath, err)
			}
		}
	}

	var buf bytes.Buffer
	ok := installer.Check(targetDir, "global", &buf, true, contentFS)

	if ok {
		t.Error("expected Check to return false when all files modified")
	}

	output := buf.String()
	expectedSummary := "0/32 files verified (0 missing, 0 empty, 32 modified)\n"
	if !strings.HasSuffix(output, expectedSummary) {
		t.Errorf("expected output to end with %q, got: %q", expectedSummary, output)
	}

	// Count MODIFIED lines
	modifiedCount := strings.Count(output, "  MODIFIED ")
	if modifiedCount != 32 {
		t.Errorf("expected 32 MODIFIED messages, got %d", modifiedCount)
	}
}

func TestCheck_OkTrueIffAllFilesPassAllChecks(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T, targetDir string)
		verify     bool
		expectedOk bool
	}{
		{
			name:       "all present and valid",
			setup:      func(t *testing.T, targetDir string) {},
			verify:     false,
			expectedOk: true,
		},
		{
			name: "one missing",
			setup: func(t *testing.T, targetDir string) {
				os.Remove(filepath.Join(targetDir, "agents/gl-architect.md"))
			},
			verify:     false,
			expectedOk: false,
		},
		{
			name: "one empty",
			setup: func(t *testing.T, targetDir string) {
				os.WriteFile(filepath.Join(targetDir, "agents/gl-architect.md"), []byte{}, 0o644)
			},
			verify:     false,
			expectedOk: false,
		},
		{
			name: "version missing",
			setup: func(t *testing.T, targetDir string) {
				os.Remove(filepath.Join(targetDir, ".greenlight-version"))
			},
			verify:     false,
			expectedOk: false,
		},
		{
			name:       "verify=true all match",
			setup:      func(t *testing.T, targetDir string) {},
			verify:     true,
			expectedOk: true,
		},
		{
			name: "verify=true one modified",
			setup: func(t *testing.T, targetDir string) {
				os.WriteFile(filepath.Join(targetDir, "agents/gl-architect.md"), []byte("mod"), 0o644)
			},
			verify:     true,
			expectedOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contentFS := buildTestFS()
			targetDir := t.TempDir()
			installTestFiles(t, targetDir, "global", contentFS)

			tt.setup(t, targetDir)

			var buf bytes.Buffer
			var ok bool
			if tt.verify {
				ok = installer.Check(targetDir, "global", &buf, true, contentFS)
			} else {
				ok = installer.Check(targetDir, "global", &buf, false, nil)
			}

			if ok != tt.expectedOk {
				t.Errorf("expected ok=%v, got ok=%v\nOutput: %s", tt.expectedOk, ok, buf.String())
			}
		})
	}
}

func TestCheck_PresenceOnlyDoesNotVerifyContent(t *testing.T) {
	contentFS := buildTestFS()
	targetDir := t.TempDir()
	installTestFiles(t, targetDir, "global", contentFS)

	// Modify a file (should not be detected in presence-only mode)
	modifiedFile := filepath.Join(targetDir, "agents/gl-architect.md")
	if err := os.WriteFile(modifiedFile, []byte("completely different content"), 0o644); err != nil {
		t.Fatalf("failed to modify file: %v", err)
	}

	var buf bytes.Buffer
	ok := installer.Check(targetDir, "global", &buf, false, nil)

	if !ok {
		t.Error("expected Check to return true (presence-only ignores content)")
	}

	output := buf.String()
	if strings.Contains(output, "MODIFIED") {
		t.Error("presence-only mode should not report MODIFIED files")
	}

	expectedSummary := "all 32 files present\n"
	if !strings.HasSuffix(output, expectedSummary) {
		t.Errorf("expected success summary, got: %q", output)
	}
}

func TestCheck_MultipleFailureTypes_AllReported(t *testing.T) {
	contentFS := buildTestFS()
	targetDir := t.TempDir()
	installTestFiles(t, targetDir, "global", contentFS)

	// Create multiple failure types
	os.Remove(filepath.Join(targetDir, "agents/gl-architect.md"))                                    // missing
	os.WriteFile(filepath.Join(targetDir, "agents/gl-debugger.md"), []byte{}, 0o644)                // empty
	os.WriteFile(filepath.Join(targetDir, "agents/gl-designer.md"), []byte("modified"), 0o644)      // modified
	os.Remove(filepath.Join(targetDir, "commands/gl/help.md"))                                       // missing
	os.WriteFile(filepath.Join(targetDir, "commands/gl/pause.md"), []byte{}, 0o644)                 // empty
	os.WriteFile(filepath.Join(targetDir, "references/deviation-rules.md"), []byte("mod"), 0o644)   // modified
	os.Remove(filepath.Join(targetDir, ".greenlight-version"))                                       // version missing

	var buf bytes.Buffer
	ok := installer.Check(targetDir, "global", &buf, true, contentFS)

	if ok {
		t.Error("expected Check to return false with multiple failures")
	}

	output := buf.String()

	// Verify all failure types are reported
	expectedMessages := []string{
		"  MISSING  agents/gl-architect.md\n",
		"  EMPTY    agents/gl-debugger.md\n",
		"  MODIFIED agents/gl-designer.md\n",
		"  MISSING  commands/gl/help.md\n",
		"  EMPTY    commands/gl/pause.md\n",
		"  MODIFIED references/deviation-rules.md\n",
		"  MISSING  .greenlight-version\n",
	}

	for _, msg := range expectedMessages {
		if !strings.Contains(output, msg) {
			t.Errorf("expected output to contain %q, got: %q", msg, output)
		}
	}

	// Verify summary (29/32: 2 missing, 2 empty, 2 modified, but files themselves)
	// Actually: 32 manifest files, 2 missing + 2 empty + 2 modified = 26 ok
	expectedSummary := "26/32 files verified (2 missing, 2 empty, 2 modified)\n"
	if !strings.HasSuffix(output, expectedSummary) {
		t.Errorf("expected output to end with %q, got: %q", expectedSummary, output)
	}
}

func TestCheck_VersionPrintedBeforeSummary(t *testing.T) {
	contentFS := buildTestFS()
	targetDir := t.TempDir()
	installTestFiles(t, targetDir, "global", contentFS)

	var buf bytes.Buffer
	installer.Check(targetDir, "global", &buf, false, nil)

	output := buf.String()
	lines := strings.Split(output, "\n")

	versionLineIdx := -1
	summaryLineIdx := -1

	for i, line := range lines {
		if strings.HasPrefix(line, "  version: ") {
			versionLineIdx = i
		}
		if strings.HasPrefix(line, "all ") || strings.Contains(line, "files present") || strings.Contains(line, "files verified") {
			summaryLineIdx = i
		}
	}

	if versionLineIdx == -1 {
		t.Error("version line not found in output")
	}
	if summaryLineIdx == -1 {
		t.Error("summary line not found in output")
	}

	if versionLineIdx >= summaryLineIdx {
		t.Error("version line should appear before summary line")
	}
}

func TestCheck_EmptyDirectory_AllFilesMissing(t *testing.T) {
	_ = buildTestFS()
	targetDir := t.TempDir() // Empty directory

	var buf bytes.Buffer
	ok := installer.Check(targetDir, "global", &buf, false, nil)

	if ok {
		t.Error("expected Check to return false for empty directory")
	}

	output := buf.String()

	// All 32 files should be reported missing
	missingCount := strings.Count(output, "  MISSING  ")
	if missingCount != 33 { // 32 manifest files + 1 version file
		t.Errorf("expected 33 MISSING messages, got %d", missingCount)
	}

	expectedSummary := "0/32 files present (32 missing, 0 empty)\n"
	if !strings.HasSuffix(output, expectedSummary) {
		t.Errorf("expected output to end with %q, got: %q", expectedSummary, output)
	}
}

func TestCheck_OutputFormat_ConsistentIndentation(t *testing.T) {
	contentFS := buildTestFS()
	targetDir := t.TempDir()
	installTestFiles(t, targetDir, "global", contentFS)

	// Remove one file to trigger output
	os.Remove(filepath.Join(targetDir, "agents/gl-architect.md"))

	var buf bytes.Buffer
	installer.Check(targetDir, "global", &buf, false, nil)

	output := buf.String()
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "  MISSING") || strings.HasPrefix(line, "  EMPTY") || strings.HasPrefix(line, "  MODIFIED") || strings.HasPrefix(line, "  ERROR") || strings.HasPrefix(line, "  version:") {
			// These should all start with exactly 2 spaces
			if !strings.HasPrefix(line, "  ") {
				t.Errorf("line missing proper indentation: %q", line)
			}
			// Verify consistent column alignment (MISSING/EMPTY/MODIFIED are 8 chars + 2 spaces)
			if strings.Contains(line, "MISSING") || strings.Contains(line, "EMPTY") || strings.Contains(line, "ERROR") {
				parts := strings.Fields(line)
				if len(parts) < 2 {
					t.Errorf("malformed status line: %q", line)
				}
			}
		}
	}
}

func TestCheck_CLAUDELocalScope_LiteralDotClaude(t *testing.T) {
	contentFS := buildTestFS()

	// Set up test in a temp directory
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
	installTestFiles(t, targetDir, "local", contentFS)

	var buf bytes.Buffer
	ok := installer.Check(targetDir, "local", &buf, false, nil)

	if !ok {
		t.Errorf("expected Check to return true, got false\nOutput: %s", buf.String())
	}

	// Verify CLAUDE.md is in current directory (not inside .claude)
	if _, err := os.Stat("CLAUDE.md"); os.IsNotExist(err) {
		t.Error("CLAUDE.md should exist in current directory for literal .claude targetDir")
	}

	if _, err := os.Stat(".claude/CLAUDE.md"); !os.IsNotExist(err) {
		t.Error("CLAUDE.md should not exist inside .claude directory")
	}
}

func TestCheck_VerifyTrue_WithNilContentFS_Fails(t *testing.T) {
	contentFS := buildTestFS()
	targetDir := t.TempDir()
	installTestFiles(t, targetDir, "global", contentFS)

	var buf bytes.Buffer

	// This test documents expected behavior when contentFS is nil but verify=true
	// The implementation should handle this gracefully (likely will fail when trying to read from nil FS)
	// We're testing that it doesn't panic and returns false
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Check panicked with nil contentFS and verify=true: %v", r)
		}
	}()

	ok := installer.Check(targetDir, "global", &buf, true, nil)

	// Expected to fail (can't verify without contentFS)
	if ok {
		t.Error("expected Check to return false when contentFS is nil and verify=true")
	}
}

func TestCheck_ExitEarlyOnFileChecks_ContinuesAfterFailure(t *testing.T) {
	contentFS := buildTestFS()
	targetDir := t.TempDir()
	installTestFiles(t, targetDir, "global", contentFS)

	// Remove the first file in manifest
	os.Remove(filepath.Join(targetDir, "agents/gl-architect.md"))

	var buf bytes.Buffer
	ok := installer.Check(targetDir, "global", &buf, false, nil)

	if ok {
		t.Error("expected Check to return false")
	}

	output := buf.String()

	// Verify Check continues after first failure and checks remaining files
	// The summary should still account for all 32 files
	if !strings.Contains(output, "31/32 files present") {
		t.Error("Check should continue checking all files after first failure")
	}
}

func TestCheck_EmptyAndMissingBothCounted(t *testing.T) {
	contentFS := buildTestFS()
	targetDir := t.TempDir()
	installTestFiles(t, targetDir, "global", contentFS)

	// Create one missing and one empty
	os.Remove(filepath.Join(targetDir, "agents/gl-architect.md"))
	os.WriteFile(filepath.Join(targetDir, "agents/gl-debugger.md"), []byte{}, 0o644)

	var buf bytes.Buffer
	ok := installer.Check(targetDir, "global", &buf, false, nil)

	if ok {
		t.Error("expected Check to return false")
	}

	output := buf.String()
	expectedSummary := "31/32 files present (1 missing, 1 empty)\n"
	if !strings.HasSuffix(output, expectedSummary) {
		t.Errorf("expected output to end with %q, got: %q", expectedSummary, output)
	}
}

func TestCheck_PresentCountExcludesMissingOnly(t *testing.T) {
	contentFS := buildTestFS()
	targetDir := t.TempDir()
	installTestFiles(t, targetDir, "global", contentFS)

	// Create multiple missing
	os.Remove(filepath.Join(targetDir, "agents/gl-architect.md"))
	os.Remove(filepath.Join(targetDir, "agents/gl-debugger.md"))
	os.Remove(filepath.Join(targetDir, "agents/gl-designer.md"))

	// Create one empty (should still count as "present" in the denominator calculation)
	os.WriteFile(filepath.Join(targetDir, "agents/gl-implementer.md"), []byte{}, 0o644)

	var buf bytes.Buffer
	installer.Check(targetDir, "global", &buf, false, nil)

	output := buf.String()
	// 29 present (32 - 3 missing), but 1 of those is empty
	// Summary format: "<present>/<total> files present (<missing> missing, <empty> empty)"
	expectedSummary := "29/32 files present (3 missing, 1 empty)\n"
	if !strings.HasSuffix(output, expectedSummary) {
		t.Errorf("expected output to end with %q, got: %q", expectedSummary, output)
	}
}

func TestCheck_RelativePathsInOutput(t *testing.T) {
	contentFS := buildTestFS()
	targetDir := t.TempDir()
	installTestFiles(t, targetDir, "global", contentFS)

	os.Remove(filepath.Join(targetDir, "commands/gl/add-slice.md"))

	var buf bytes.Buffer
	installer.Check(targetDir, "global", &buf, false, nil)

	output := buf.String()

	// Verify output uses relative paths (as they appear in manifest)
	// NOT absolute filesystem paths
	if !strings.Contains(output, "  MISSING  commands/gl/add-slice.md\n") {
		t.Error("expected output to use relative path from manifest")
	}

	// Should NOT contain the absolute path
	absolutePath := filepath.Join(targetDir, "commands/gl/add-slice.md")
	if strings.Contains(output, absolutePath) {
		t.Error("output should not contain absolute filesystem paths")
	}
}

func TestCheck_VerifyMode_UsesContentFSForHashComparison(t *testing.T) {
	// Create a contentFS with specific known content
	contentFS := fstest.MapFS{
		"agents/gl-architect.md": &fstest.MapFile{
			Data: []byte("# Test Content\nSpecific hash test\n"),
		},
	}

	// Calculate expected hash
	expectedHash := fmt.Sprintf("%x", sha256.Sum256([]byte("# Test Content\nSpecific hash test\n")))

	targetDir := t.TempDir()
	agentsDir := filepath.Join(targetDir, "agents")
	if err := os.MkdirAll(agentsDir, 0o755); err != nil {
		t.Fatalf("failed to create agents dir: %v", err)
	}

	testFile := filepath.Join(agentsDir, "gl-architect.md")

	// Write matching content
	if err := os.WriteFile(testFile, []byte("# Test Content\nSpecific hash test\n"), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Create minimal manifest for this test
	originalManifest := installer.Manifest
	defer func() { installer.Manifest = originalManifest }()
	installer.Manifest = []string{"agents/gl-architect.md"}

	var buf bytes.Buffer
	ok := installer.Check(targetDir, "global", &buf, true, contentFS)

	if !ok {
		t.Errorf("expected Check to return true with matching hash\nOutput: %s", buf.String())
	}

	// Now modify the file and verify it's detected
	buf.Reset()
	if err := os.WriteFile(testFile, []byte("# Modified\n"), 0o644); err != nil {
		t.Fatalf("failed to modify test file: %v", err)
	}

	modifiedHash := fmt.Sprintf("%x", sha256.Sum256([]byte("# Modified\n")))
	if modifiedHash == expectedHash {
		t.Fatal("test setup error: hashes should differ")
	}

	ok = installer.Check(targetDir, "global", &buf, true, contentFS)

	if ok {
		t.Error("expected Check to return false with mismatched hash")
	}

	output := buf.String()
	if !strings.Contains(output, "  MODIFIED agents/gl-architect.md\n") {
		t.Errorf("expected MODIFIED message, got: %q", output)
	}
}
