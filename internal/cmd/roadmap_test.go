package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/atlantic-blue/greenlight/internal/cmd"
)

// ----------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------

// writeRoadmapFile writes content to a file, failing the test on any error.
func writeRoadmapFile(t *testing.T, path string, content string) {
	t.Helper()
	if writeError := os.WriteFile(path, []byte(content), 0o644); writeError != nil {
		t.Fatalf("writeRoadmapFile(%q) error: %v", path, writeError)
	}
}

// setupRoadmapProject creates a temp directory with a .greenlight/ directory
// and optionally a ROADMAP.md inside it, then calls t.Chdir(tmpDir) so that
// RunRoadmap can locate .greenlight/ from the working directory.
//
// Pass an empty roadmapContent to skip writing ROADMAP.md (simulates missing
// file).
func setupRoadmapProject(t *testing.T, roadmapContent string) string {
	t.Helper()

	tmpDir := t.TempDir()

	greenlightDir := filepath.Join(tmpDir, ".greenlight")
	if mkdirError := os.MkdirAll(greenlightDir, 0o755); mkdirError != nil {
		t.Fatalf("setupRoadmapProject: failed to create .greenlight dir: %v", mkdirError)
	}

	if roadmapContent != "" {
		roadmapPath := filepath.Join(greenlightDir, "ROADMAP.md")
		writeRoadmapFile(t, roadmapPath, roadmapContent)
	}

	t.Chdir(tmpDir)
	return tmpDir
}

// ----------------------------------------------------------------------------
// C-101 — RunRoadmap: happy paths
// ----------------------------------------------------------------------------

// TestRunRoadmap_PrintsContentsVerbatim verifies that when ROADMAP.md exists,
// RunRoadmap prints its exact contents to the writer and returns exit code 0.
func TestRunRoadmap_PrintsContentsVerbatim(t *testing.T) {
	const roadmapContent = "# Roadmap\n\n- Milestone 1\n- Milestone 2\n"
	setupRoadmapProject(t, roadmapContent)

	var buf bytes.Buffer
	exitCode := cmd.RunRoadmap([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	output := buf.String()
	if output != roadmapContent {
		t.Errorf("expected verbatim roadmap contents\nwant: %q\ngot:  %q", roadmapContent, output)
	}
}

// TestRunRoadmap_MultiLineContentPreservedExactly verifies that multi-line
// content including blank lines, headers, and indentation is output unchanged.
func TestRunRoadmap_MultiLineContentPreservedExactly(t *testing.T) {
	const roadmapContent = `# Project Roadmap

## Milestone 1: Foundation
- S-01: Bootstrap project
- S-02: Database schema

## Milestone 2: Core Features
- S-10: User registration
- S-11: Authentication

## Milestone 3: Polish
- S-20: Dashboard UI
`
	setupRoadmapProject(t, roadmapContent)

	var buf bytes.Buffer
	exitCode := cmd.RunRoadmap([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	output := buf.String()
	if output != roadmapContent {
		t.Errorf("multi-line roadmap content was not preserved exactly\nwant: %q\ngot:  %q", roadmapContent, output)
	}
}

// TestRunRoadmap_ExitCode0OnSuccess verifies that a successful read returns 0.
func TestRunRoadmap_ExitCode0OnSuccess(t *testing.T) {
	setupRoadmapProject(t, "# Roadmap\n")

	var buf bytes.Buffer
	exitCode := cmd.RunRoadmap([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 on success, got %d", exitCode)
	}
}

// ----------------------------------------------------------------------------
// C-101 — RunRoadmap: error cases
// ----------------------------------------------------------------------------

// TestRunRoadmap_NoGreenlightDir_MentionsGreenlightProject verifies that when
// there is no .greenlight/ directory, the output mentions "greenlight project"
// and the exit code is 1.
func TestRunRoadmap_NoGreenlightDir_MentionsGreenlightProject(t *testing.T) {
	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	var buf bytes.Buffer
	exitCode := cmd.RunRoadmap([]string{}, &buf)

	if exitCode != 1 {
		t.Errorf("expected exit code 1 when .greenlight/ not found, got %d", exitCode)
	}

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	if !strings.Contains(lowerOutput, "greenlight") {
		t.Errorf("expected error message mentioning 'greenlight', got:\n%s", output)
	}
}

// TestRunRoadmap_NoGreenlightDir_ReturnsExitCode1 verifies the exact exit code
// when .greenlight/ is absent.
func TestRunRoadmap_NoGreenlightDir_ReturnsExitCode1(t *testing.T) {
	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	var buf bytes.Buffer
	exitCode := cmd.RunRoadmap([]string{}, &buf)

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
}

// TestRunRoadmap_MissingROADMAPFile_MentionsDesign verifies that when
// .greenlight/ exists but ROADMAP.md is absent, the error message mentions
// "design" (hinting where to find roadmap documentation).
func TestRunRoadmap_MissingROADMAPFile_MentionsDesign(t *testing.T) {
	// Pass empty content so setupRoadmapProject skips writing ROADMAP.md.
	setupRoadmapProject(t, "")

	var buf bytes.Buffer
	exitCode := cmd.RunRoadmap([]string{}, &buf)

	if exitCode != 1 {
		t.Errorf("expected exit code 1 when ROADMAP.md not found, got %d", exitCode)
	}

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	if !strings.Contains(lowerOutput, "design") {
		t.Errorf("expected error message mentioning 'design', got:\n%s", output)
	}
}

// TestRunRoadmap_MissingROADMAPFile_ReturnsExitCode1 verifies the exact exit
// code when .greenlight/ exists but ROADMAP.md is absent.
func TestRunRoadmap_MissingROADMAPFile_ReturnsExitCode1(t *testing.T) {
	setupRoadmapProject(t, "")

	var buf bytes.Buffer
	exitCode := cmd.RunRoadmap([]string{}, &buf)

	if exitCode != 1 {
		t.Errorf("expected exit code 1 for missing ROADMAP.md, got %d", exitCode)
	}
}

// ----------------------------------------------------------------------------
// C-101 — RunRoadmap: output writer contract
// ----------------------------------------------------------------------------

// TestRunRoadmap_WritesToProvidedWriter verifies that all output from RunRoadmap
// is directed to the provided io.Writer and not to os.Stdout.
func TestRunRoadmap_WritesToProvidedWriter(t *testing.T) {
	const roadmapContent = "# Roadmap\n"
	setupRoadmapProject(t, roadmapContent)

	var buf bytes.Buffer
	cmd.RunRoadmap([]string{}, &buf)

	if buf.Len() == 0 {
		t.Error("expected output to be written to the provided writer, got empty buffer")
	}
}

// TestRunRoadmap_ErrorWritesToProvidedWriter verifies that error messages are
// also written to the provided io.Writer and not to os.Stderr.
func TestRunRoadmap_ErrorWritesToProvidedWriter(t *testing.T) {
	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	var buf bytes.Buffer
	cmd.RunRoadmap([]string{}, &buf)

	if buf.Len() == 0 {
		t.Error("expected error output to be written to the provided writer, got empty buffer")
	}
}
