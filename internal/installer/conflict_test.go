package installer

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// C-05 Tests: ConflictStrategy Type

func TestConflictStrategy_Constants(t *testing.T) {
	tests := []struct {
		name     string
		value    ConflictStrategy
		expected string
	}{
		{
			name:     "ConflictKeep has correct value",
			value:    ConflictKeep,
			expected: "keep",
		},
		{
			name:     "ConflictReplace has correct value",
			value:    ConflictReplace,
			expected: "replace",
		},
		{
			name:     "ConflictAppend has correct value",
			value:    ConflictAppend,
			expected: "append",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, string(tt.value))
			}
		})
	}
}

func TestConflictStrategy_TypeCheck(t *testing.T) {
	// Compile-time type check
	var _ ConflictStrategy = ConflictKeep
	var _ ConflictStrategy = ConflictReplace
	var _ ConflictStrategy = ConflictAppend

	// Type can be assigned and compared
	var strategy ConflictStrategy = ConflictKeep
	if strategy != ConflictKeep {
		t.Error("type comparison failed")
	}
}

func TestConflictStrategy_Distinct(t *testing.T) {
	strategies := []ConflictStrategy{ConflictKeep, ConflictReplace, ConflictAppend}

	// Verify all three are distinct
	for i, s1 := range strategies {
		for j, s2 := range strategies {
			if i != j && s1 == s2 {
				t.Errorf("strategies at index %d and %d are not distinct: %q == %q", i, j, s1, s2)
			}
		}
	}
}

// C-06 Tests: HandleConflict

func TestHandleConflict_NoExistingFile(t *testing.T) {
	srcData := []byte("# Greenlight CLAUDE.md\n\nTest content\n")

	tests := []struct {
		name     string
		strategy ConflictStrategy
	}{
		{
			name:     "no existing file with keep strategy",
			strategy: ConflictKeep,
		},
		{
			name:     "no existing file with replace strategy",
			strategy: ConflictReplace,
		},
		{
			name:     "no existing file with append strategy",
			strategy: ConflictAppend,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			destPath := filepath.Join(tempDir, "CLAUDE.md")
			var buf bytes.Buffer

			err := handleConflict(destPath, srcData, tt.strategy, &buf)
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}

			// Verify file was created
			if _, err := os.Stat(destPath); os.IsNotExist(err) {
				t.Fatal("expected file to be created but it does not exist")
			}

			// Verify content
			content, err := os.ReadFile(destPath)
			if err != nil {
				t.Fatalf("failed to read created file: %v", err)
			}
			if !bytes.Equal(content, srcData) {
				t.Errorf("expected content %q, got %q", srcData, content)
			}

			// Verify no extra files were created
			entries, err := os.ReadDir(tempDir)
			if err != nil {
				t.Fatalf("failed to read temp dir: %v", err)
			}
			if len(entries) != 1 {
				t.Errorf("expected 1 file, got %d files", len(entries))
			}
		})
	}
}

func TestHandleConflict_KeepStrategy(t *testing.T) {
	existingContent := []byte("# Existing CLAUDE.md\n\nOriginal content\n")
	srcData := []byte("# Greenlight CLAUDE.md\n\nNew content\n")

	tempDir := t.TempDir()
	destPath := filepath.Join(tempDir, "CLAUDE.md")
	greenlightPath := filepath.Join(tempDir, "CLAUDE_GREENLIGHT.md")
	var buf bytes.Buffer

	// Create existing file
	if err := os.WriteFile(destPath, existingContent, 0o644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	// Execute
	err := handleConflict(destPath, srcData, ConflictKeep, &buf)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify existing file is untouched
	content, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("failed to read existing file: %v", err)
	}
	if !bytes.Equal(content, existingContent) {
		t.Errorf("existing file was modified. expected %q, got %q", existingContent, content)
	}

	// Verify CLAUDE_GREENLIGHT.md was created with srcData
	greenlightContent, err := os.ReadFile(greenlightPath)
	if err != nil {
		t.Fatalf("CLAUDE_GREENLIGHT.md was not created: %v", err)
	}
	if !bytes.Equal(greenlightContent, srcData) {
		t.Errorf("CLAUDE_GREENLIGHT.md has wrong content. expected %q, got %q", srcData, greenlightContent)
	}

	// Verify output message
	output := buf.String()
	expectedMessage := "  existing CLAUDE.md kept; greenlight version saved as CLAUDE_GREENLIGHT.md\n"
	if output != expectedMessage {
		t.Errorf("expected output %q, got %q", expectedMessage, output)
	}
}

func TestHandleConflict_ReplaceStrategy(t *testing.T) {
	existingContent := []byte("# Existing CLAUDE.md\n\nOriginal content\n")
	srcData := []byte("# Greenlight CLAUDE.md\n\nNew content\n")

	tempDir := t.TempDir()
	destPath := filepath.Join(tempDir, "CLAUDE.md")
	backupPath := destPath + ".backup"
	var buf bytes.Buffer

	// Create existing file
	if err := os.WriteFile(destPath, existingContent, 0o644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	// Execute
	err := handleConflict(destPath, srcData, ConflictReplace, &buf)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify backup was created with original content
	backupContent, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("backup file was not created: %v", err)
	}
	if !bytes.Equal(backupContent, existingContent) {
		t.Errorf("backup has wrong content. expected %q, got %q", existingContent, backupContent)
	}

	// Verify destPath was overwritten with srcData
	content, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("failed to read destination file: %v", err)
	}
	if !bytes.Equal(content, srcData) {
		t.Errorf("destination file has wrong content. expected %q, got %q", srcData, content)
	}

	// Verify output message
	output := buf.String()
	expectedMessage := "  existing CLAUDE.md backed up to " + backupPath + "\n"
	if output != expectedMessage {
		t.Errorf("expected output %q, got %q", expectedMessage, output)
	}
}

func TestHandleConflict_AppendStrategy_ExistingEndsWithNewline(t *testing.T) {
	existingContent := []byte("# Existing CLAUDE.md\n\nOriginal content\n")
	srcData := []byte("# Greenlight CLAUDE.md\n\nNew content\n")

	tempDir := t.TempDir()
	destPath := filepath.Join(tempDir, "CLAUDE.md")
	var buf bytes.Buffer

	// Create existing file
	if err := os.WriteFile(destPath, existingContent, 0o644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	// Execute
	err := handleConflict(destPath, srcData, ConflictAppend, &buf)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify combined content
	content, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("failed to read destination file: %v", err)
	}

	expectedContent := append(existingContent, srcData...)
	if !bytes.Equal(content, expectedContent) {
		t.Errorf("expected content %q, got %q", expectedContent, content)
	}

	// Verify no double newline at join point
	contentStr := string(content)
	if strings.Contains(contentStr, "\n\n\n") {
		t.Error("found triple newline in content, indicating double newline at join point")
	}

	// Verify output message
	output := buf.String()
	expectedMessage := "  greenlight content appended to existing CLAUDE.md\n"
	if output != expectedMessage {
		t.Errorf("expected output %q, got %q", expectedMessage, output)
	}
}

func TestHandleConflict_AppendStrategy_ExistingNoNewline(t *testing.T) {
	existingContent := []byte("# Existing CLAUDE.md\n\nOriginal content")
	srcData := []byte("# Greenlight CLAUDE.md\n\nNew content\n")

	tempDir := t.TempDir()
	destPath := filepath.Join(tempDir, "CLAUDE.md")
	var buf bytes.Buffer

	// Create existing file
	if err := os.WriteFile(destPath, existingContent, 0o644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	// Execute
	err := handleConflict(destPath, srcData, ConflictAppend, &buf)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify combined content has newline inserted
	content, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("failed to read destination file: %v", err)
	}

	expectedContent := append(existingContent, '\n')
	expectedContent = append(expectedContent, srcData...)

	if !bytes.Equal(content, expectedContent) {
		t.Errorf("expected content %q, got %q", expectedContent, content)
	}

	// Verify newline was inserted at join point
	contentStr := string(content)
	if !strings.Contains(contentStr, "Original content\n# Greenlight CLAUDE.md") {
		t.Error("newline was not inserted at join point")
	}

	// Verify no double newline at join point
	if strings.Contains(contentStr, "Original content\n\n# Greenlight CLAUDE.md") {
		t.Error("double newline at join point")
	}

	// Verify output message
	output := buf.String()
	expectedMessage := "  greenlight content appended to existing CLAUDE.md\n"
	if output != expectedMessage {
		t.Errorf("expected output %q, got %q", expectedMessage, output)
	}
}

func TestHandleConflict_UnknownStrategy(t *testing.T) {
	srcData := []byte("# Greenlight CLAUDE.md\n")
	existingContent := []byte("# Existing\n")

	tempDir := t.TempDir()
	destPath := filepath.Join(tempDir, "CLAUDE.md")
	var buf bytes.Buffer

	// Create existing file so the strategy switch is reached
	if err := os.WriteFile(destPath, existingContent, 0o644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	unknownStrategy := ConflictStrategy("merge")

	err := handleConflict(destPath, srcData, unknownStrategy, &buf)
	if err == nil {
		t.Fatal("expected error for unknown strategy, got nil")
	}

	// Verify error message
	expectedError := "unknown conflict strategy: merge"
	if err.Error() != expectedError {
		t.Errorf("expected error %q, got %q", expectedError, err.Error())
	}
}

func TestHandleConflict_DirectoryCreation(t *testing.T) {
	srcData := []byte("# Greenlight CLAUDE.md\n")

	tempDir := t.TempDir()
	nestedPath := filepath.Join(tempDir, "nested", "deeper", "CLAUDE.md")
	var buf bytes.Buffer

	err := handleConflict(nestedPath, srcData, ConflictKeep, &buf)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify file was created
	content, err := os.ReadFile(nestedPath)
	if err != nil {
		t.Fatalf("file was not created: %v", err)
	}
	if !bytes.Equal(content, srcData) {
		t.Errorf("expected content %q, got %q", srcData, content)
	}

	// Verify directories were created
	dirPath := filepath.Dir(nestedPath)
	info, err := os.Stat(dirPath)
	if err != nil {
		t.Fatalf("directories were not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected directory, got file")
	}

	// Verify directory permissions
	expectedDirPerms := os.FileMode(0o755)
	if info.Mode().Perm() != expectedDirPerms {
		t.Errorf("expected directory permissions %v, got %v", expectedDirPerms, info.Mode().Perm())
	}
}

func TestHandleConflict_FilePermissions(t *testing.T) {
	srcData := []byte("# Greenlight CLAUDE.md\n")

	tests := []struct {
		name     string
		strategy ConflictStrategy
		setup    func(destPath string) error
		checkPath string
	}{
		{
			name:     "new file has 0o644 permissions",
			strategy: ConflictKeep,
			setup:    nil,
			checkPath: "CLAUDE.md",
		},
		{
			name:     "keep strategy - CLAUDE_GREENLIGHT.md has 0o644 permissions",
			strategy: ConflictKeep,
			setup: func(destPath string) error {
				return os.WriteFile(destPath, []byte("existing\n"), 0o644)
			},
			checkPath: "CLAUDE_GREENLIGHT.md",
		},
		{
			name:     "replace strategy - file has 0o644 permissions",
			strategy: ConflictReplace,
			setup: func(destPath string) error {
				return os.WriteFile(destPath, []byte("existing\n"), 0o644)
			},
			checkPath: "CLAUDE.md",
		},
		{
			name:     "replace strategy - backup has 0o644 permissions",
			strategy: ConflictReplace,
			setup: func(destPath string) error {
				return os.WriteFile(destPath, []byte("existing\n"), 0o644)
			},
			checkPath: "CLAUDE.md.backup",
		},
		{
			name:     "append strategy - file has 0o644 permissions",
			strategy: ConflictAppend,
			setup: func(destPath string) error {
				return os.WriteFile(destPath, []byte("existing\n"), 0o644)
			},
			checkPath: "CLAUDE.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			destPath := filepath.Join(tempDir, "CLAUDE.md")
			var buf bytes.Buffer

			if tt.setup != nil {
				if err := tt.setup(destPath); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			err := handleConflict(destPath, srcData, tt.strategy, &buf)
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}

			// Check permissions
			checkPath := filepath.Join(tempDir, tt.checkPath)
			info, err := os.Stat(checkPath)
			if err != nil {
				t.Fatalf("failed to stat %s: %v", checkPath, err)
			}

			expectedPerms := os.FileMode(0o644)
			if info.Mode().Perm() != expectedPerms {
				t.Errorf("expected permissions %v, got %v", expectedPerms, info.Mode().Perm())
			}
		})
	}
}

func TestHandleConflict_EmptyExistingFile(t *testing.T) {
	emptyContent := []byte("")
	srcData := []byte("# Greenlight CLAUDE.md\n\nNew content\n")

	tempDir := t.TempDir()
	destPath := filepath.Join(tempDir, "CLAUDE.md")
	var buf bytes.Buffer

	// Create empty existing file
	if err := os.WriteFile(destPath, emptyContent, 0o644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	// Test append strategy with empty existing file
	err := handleConflict(destPath, srcData, ConflictAppend, &buf)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify content
	content, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("failed to read destination file: %v", err)
	}

	// Empty file (len 0) skips the newline guard (len > 0 check),
	// so srcData is appended directly
	if !bytes.Equal(content, srcData) {
		t.Errorf("expected content %q, got %q", srcData, content)
	}
}

func TestHandleConflict_LargeFile(t *testing.T) {
	// Test with larger files to ensure no buffer issues
	existingContent := bytes.Repeat([]byte("existing line\n"), 1000)
	srcData := bytes.Repeat([]byte("new line\n"), 1000)

	tempDir := t.TempDir()
	destPath := filepath.Join(tempDir, "CLAUDE.md")
	var buf bytes.Buffer

	// Create existing file
	if err := os.WriteFile(destPath, existingContent, 0o644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	// Test append strategy
	err := handleConflict(destPath, srcData, ConflictAppend, &buf)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify combined content
	content, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("failed to read destination file: %v", err)
	}

	expectedContent := append(existingContent, srcData...)
	if !bytes.Equal(content, expectedContent) {
		t.Error("content mismatch for large file")
	}
}

func TestHandleConflict_OutputWriter(t *testing.T) {
	// Verify that different writers work correctly
	existingContent := []byte("existing\n")
	srcData := []byte("new\n")

	tests := []struct {
		name     string
		strategy ConflictStrategy
		expectOutput bool
	}{
		{
			name:     "keep strategy writes to output",
			strategy: ConflictKeep,
			expectOutput: true,
		},
		{
			name:     "replace strategy writes to output",
			strategy: ConflictReplace,
			expectOutput: true,
		},
		{
			name:     "append strategy writes to output",
			strategy: ConflictAppend,
			expectOutput: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			destPath := filepath.Join(tempDir, "CLAUDE.md")
			var buf bytes.Buffer

			// Create existing file
			if err := os.WriteFile(destPath, existingContent, 0o644); err != nil {
				t.Fatalf("failed to create existing file: %v", err)
			}

			err := handleConflict(destPath, srcData, tt.strategy, &buf)
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}

			if tt.expectOutput && buf.Len() == 0 {
				t.Error("expected output but got none")
			}
		})
	}
}

func TestHandleConflict_AbsoluteAndRelativePaths(t *testing.T) {
	srcData := []byte("# Greenlight CLAUDE.md\n")

	tests := []struct {
		name     string
		pathFunc func(tempDir string) string
	}{
		{
			name: "absolute path",
			pathFunc: func(tempDir string) string {
				return filepath.Join(tempDir, "CLAUDE.md")
			},
		},
		{
			name: "relative path",
			pathFunc: func(tempDir string) string {
				return "CLAUDE.md"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			// Change to temp directory for relative path test
			if tt.name == "relative path" {
				originalWd, err := os.Getwd()
				if err != nil {
					t.Fatalf("failed to get working directory: %v", err)
				}
				defer os.Chdir(originalWd)

				if err := os.Chdir(tempDir); err != nil {
					t.Fatalf("failed to change directory: %v", err)
				}
			}

			destPath := tt.pathFunc(tempDir)
			var buf bytes.Buffer

			err := handleConflict(destPath, srcData, ConflictKeep, &buf)
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}

			// Verify file was created
			var checkPath string
			if filepath.IsAbs(destPath) {
				checkPath = destPath
			} else {
				checkPath = filepath.Join(tempDir, destPath)
			}

			content, err := os.ReadFile(checkPath)
			if err != nil {
				t.Fatalf("file was not created: %v", err)
			}
			if !bytes.Equal(content, srcData) {
				t.Errorf("expected content %q, got %q", srcData, content)
			}
		})
	}
}
