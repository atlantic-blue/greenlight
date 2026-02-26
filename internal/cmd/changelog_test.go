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

// writeSummaryFile writes content to a file path, failing the test on any error.
func writeSummaryFile(t *testing.T, path string, content string) {
	t.Helper()
	if writeError := os.WriteFile(path, []byte(content), 0o644); writeError != nil {
		t.Fatalf("writeSummaryFile(%q) error: %v", path, writeError)
	}
}

// summaryEntry describes a single file to place in the summaries directory.
type summaryEntry struct {
	filename string
	content  string
}

// setupChangelogProject creates a temp directory with a .greenlight/summaries/
// directory populated with the given entries, then calls t.Chdir(tmpDir).
//
// Pass nil or an empty slice to create the summaries dir but leave it empty.
// Pass createSummariesDir=false to omit the summaries directory entirely (only
// .greenlight/ is created).
func setupChangelogProject(t *testing.T, entries []summaryEntry, createSummariesDir bool) string {
	t.Helper()

	tmpDir := t.TempDir()

	greenlightDir := filepath.Join(tmpDir, ".greenlight")
	if mkdirError := os.MkdirAll(greenlightDir, 0o755); mkdirError != nil {
		t.Fatalf("setupChangelogProject: failed to create .greenlight dir: %v", mkdirError)
	}

	if createSummariesDir {
		summariesDir := filepath.Join(greenlightDir, "summaries")
		if mkdirError := os.MkdirAll(summariesDir, 0o755); mkdirError != nil {
			t.Fatalf("setupChangelogProject: failed to create summaries dir: %v", mkdirError)
		}

		for _, entry := range entries {
			filePath := filepath.Join(summariesDir, entry.filename)
			writeSummaryFile(t, filePath, entry.content)
		}
	}

	t.Chdir(tmpDir)
	return tmpDir
}

// ----------------------------------------------------------------------------
// C-102 — RunChangelog: happy paths
// ----------------------------------------------------------------------------

// TestRunChangelog_MultipleSummariesPrintedInFilenameOrder verifies that when
// multiple summary files exist, RunChangelog prints them in ascending filename
// order, separated by "---", and returns exit code 0.
func TestRunChangelog_MultipleSummariesPrintedInFilenameOrder(t *testing.T) {
	entries := []summaryEntry{
		{filename: "S-01-SUMMARY.md", content: "## S-01: Bootstrap\nSlice one complete.\n"},
		{filename: "S-02-SUMMARY.md", content: "## S-02: Database\nSlice two complete.\n"},
		{filename: "S-03-SUMMARY.md", content: "## S-03: Auth\nSlice three complete.\n"},
	}

	setupChangelogProject(t, entries, true)

	var buf bytes.Buffer
	exitCode := cmd.RunChangelog([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	output := buf.String()

	// All three summaries must be present.
	if !strings.Contains(output, "S-01: Bootstrap") {
		t.Errorf("expected S-01 content in output:\n%s", output)
	}
	if !strings.Contains(output, "S-02: Database") {
		t.Errorf("expected S-02 content in output:\n%s", output)
	}
	if !strings.Contains(output, "S-03: Auth") {
		t.Errorf("expected S-03 content in output:\n%s", output)
	}

	// Summaries must be separated by "---".
	if !strings.Contains(output, "---") {
		t.Errorf("expected '---' separator between summaries:\n%s", output)
	}
}

// TestRunChangelog_SingleSummaryPrintedWithoutSeparator verifies that a single
// summary is printed without a "---" separator and the content is present.
func TestRunChangelog_SingleSummaryPrintedWithoutSeparator(t *testing.T) {
	entries := []summaryEntry{
		{filename: "S-05-SUMMARY.md", content: "## S-05: Single Entry\nOnly one summary.\n"},
	}

	setupChangelogProject(t, entries, true)

	var buf bytes.Buffer
	exitCode := cmd.RunChangelog([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	output := buf.String()

	if !strings.Contains(output, "S-05: Single Entry") {
		t.Errorf("expected single summary content in output:\n%s", output)
	}

	// A single entry should not be followed by a "---" separator.
	if strings.Contains(output, "---") {
		t.Errorf("single summary should not contain '---' separator:\n%s", output)
	}
}

// TestRunChangelog_SummariesSortedAscendingByFilename verifies that summaries
// are output in strict ascending filename order (oldest first), regardless of
// the order they were written to disk.
func TestRunChangelog_SummariesSortedAscendingByFilename(t *testing.T) {
	// Write entries in reverse order; output must be S-10 before S-22 before S-34.
	entries := []summaryEntry{
		{filename: "S-34-SUMMARY.md", content: "Content of S-34\n"},
		{filename: "S-10-SUMMARY.md", content: "Content of S-10\n"},
		{filename: "S-22-SUMMARY.md", content: "Content of S-22\n"},
	}

	setupChangelogProject(t, entries, true)

	var buf bytes.Buffer
	exitCode := cmd.RunChangelog([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	output := buf.String()

	posS10 := strings.Index(output, "Content of S-10")
	posS22 := strings.Index(output, "Content of S-22")
	posS34 := strings.Index(output, "Content of S-34")

	if posS10 < 0 || posS22 < 0 || posS34 < 0 {
		t.Fatalf("not all summaries found in output:\n%s", output)
	}

	if !(posS10 < posS22 && posS22 < posS34) {
		t.Errorf("summaries not in ascending filename order: S-10 at %d, S-22 at %d, S-34 at %d\n%s",
			posS10, posS22, posS34, output)
	}
}

// TestRunChangelog_ExitCode0OnSuccess verifies that a successful run with
// multiple summaries returns exit code 0.
func TestRunChangelog_ExitCode0OnSuccess(t *testing.T) {
	entries := []summaryEntry{
		{filename: "S-01-SUMMARY.md", content: "First summary.\n"},
		{filename: "S-02-SUMMARY.md", content: "Second summary.\n"},
	}

	setupChangelogProject(t, entries, true)

	var buf bytes.Buffer
	exitCode := cmd.RunChangelog([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
}

// ----------------------------------------------------------------------------
// C-102 — RunChangelog: error cases
// ----------------------------------------------------------------------------

// TestRunChangelog_NoGreenlightDir_MentionsGreenlightProject verifies that when
// there is no .greenlight/ directory, the output mentions "greenlight project"
// and the exit code is 1.
func TestRunChangelog_NoGreenlightDir_MentionsGreenlightProject(t *testing.T) {
	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	var buf bytes.Buffer
	exitCode := cmd.RunChangelog([]string{}, &buf)

	if exitCode != 1 {
		t.Errorf("expected exit code 1 when .greenlight/ not found, got %d", exitCode)
	}

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	if !strings.Contains(lowerOutput, "greenlight") {
		t.Errorf("expected error message mentioning 'greenlight', got:\n%s", output)
	}
}

// TestRunChangelog_NoGreenlightDir_ReturnsExitCode1 verifies the exact exit
// code when .greenlight/ is absent.
func TestRunChangelog_NoGreenlightDir_ReturnsExitCode1(t *testing.T) {
	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	var buf bytes.Buffer
	exitCode := cmd.RunChangelog([]string{}, &buf)

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
}

// ----------------------------------------------------------------------------
// C-102 — RunChangelog: "No changelog entries yet." cases (exit 0)
// ----------------------------------------------------------------------------

// TestRunChangelog_NoSummariesDir_PrintsNoEntriesMessage verifies that when
// .greenlight/ exists but the summaries/ subdirectory is absent, RunChangelog
// prints "No changelog entries yet." and returns exit code 0.
func TestRunChangelog_NoSummariesDir_PrintsNoEntriesMessage(t *testing.T) {
	// createSummariesDir=false: only .greenlight/ is created, no summaries/.
	setupChangelogProject(t, nil, false)

	var buf bytes.Buffer
	exitCode := cmd.RunChangelog([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 when summaries dir absent, got %d", exitCode)
	}

	output := buf.String()
	if !strings.Contains(output, "No changelog entries yet.") {
		t.Errorf("expected 'No changelog entries yet.' in output, got:\n%s", output)
	}
}

// TestRunChangelog_EmptySummariesDir_PrintsNoEntriesMessage verifies that when
// .greenlight/summaries/ exists but contains no .md files, RunChangelog prints
// "No changelog entries yet." and returns exit code 0.
func TestRunChangelog_EmptySummariesDir_PrintsNoEntriesMessage(t *testing.T) {
	// createSummariesDir=true, entries=nil: directory exists but is empty.
	setupChangelogProject(t, nil, true)

	var buf bytes.Buffer
	exitCode := cmd.RunChangelog([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 for empty summaries dir, got %d", exitCode)
	}

	output := buf.String()
	if !strings.Contains(output, "No changelog entries yet.") {
		t.Errorf("expected 'No changelog entries yet.' in output, got:\n%s", output)
	}
}

// TestRunChangelog_NoSummariesDir_ExitCode0 verifies the exact exit code when
// the summaries directory is absent (not an error condition per C-102).
func TestRunChangelog_NoSummariesDir_ExitCode0(t *testing.T) {
	setupChangelogProject(t, nil, false)

	var buf bytes.Buffer
	exitCode := cmd.RunChangelog([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 (not an error) for missing summaries dir, got %d", exitCode)
	}
}

// TestRunChangelog_EmptySummariesDir_ExitCode0 verifies the exact exit code
// when the summaries directory is present but empty.
func TestRunChangelog_EmptySummariesDir_ExitCode0(t *testing.T) {
	setupChangelogProject(t, nil, true)

	var buf bytes.Buffer
	exitCode := cmd.RunChangelog([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 (not an error) for empty summaries dir, got %d", exitCode)
	}
}

// ----------------------------------------------------------------------------
// C-102 — RunChangelog: output writer contract
// ----------------------------------------------------------------------------

// TestRunChangelog_WritesToProvidedWriter verifies that all output from
// RunChangelog is directed to the provided io.Writer and not to os.Stdout.
func TestRunChangelog_WritesToProvidedWriter(t *testing.T) {
	entries := []summaryEntry{
		{filename: "S-01-SUMMARY.md", content: "First summary.\n"},
	}

	setupChangelogProject(t, entries, true)

	var buf bytes.Buffer
	cmd.RunChangelog([]string{}, &buf)

	if buf.Len() == 0 {
		t.Error("expected output to be written to the provided writer, got empty buffer")
	}
}

// TestRunChangelog_ErrorWritesToProvidedWriter verifies that error output is
// written to the provided io.Writer and not to os.Stderr.
func TestRunChangelog_ErrorWritesToProvidedWriter(t *testing.T) {
	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	var buf bytes.Buffer
	cmd.RunChangelog([]string{}, &buf)

	if buf.Len() == 0 {
		t.Error("expected error output to be written to the provided writer, got empty buffer")
	}
}

// TestRunChangelog_NoEntriesWritesToProvidedWriter verifies that the
// "No changelog entries yet." message is written to the provided writer.
func TestRunChangelog_NoEntriesWritesToProvidedWriter(t *testing.T) {
	setupChangelogProject(t, nil, true)

	var buf bytes.Buffer
	cmd.RunChangelog([]string{}, &buf)

	if buf.Len() == 0 {
		t.Error("expected 'No changelog entries yet.' to be written to the provided writer, got empty buffer")
	}
}

// ----------------------------------------------------------------------------
// C-102 — RunChangelog: separator placement
// ----------------------------------------------------------------------------

// TestRunChangelog_SeparatorBetweenEntriesNotAtEnd verifies that "---" appears
// between entries but not appended after the final entry.
func TestRunChangelog_SeparatorBetweenEntriesNotAtEnd(t *testing.T) {
	entries := []summaryEntry{
		{filename: "S-01-SUMMARY.md", content: "First.\n"},
		{filename: "S-02-SUMMARY.md", content: "Second.\n"},
	}

	setupChangelogProject(t, entries, true)

	var buf bytes.Buffer
	exitCode := cmd.RunChangelog([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	output := buf.String()

	// "---" must appear between the two entries.
	if !strings.Contains(output, "---") {
		t.Errorf("expected '---' separator between entries:\n%s", output)
	}

	// The separator must appear before the final entry's content.
	posSep := strings.Index(output, "---")
	posSecond := strings.Index(output, "Second.")
	if posSep >= posSecond {
		t.Errorf("separator must appear before the second entry: sep at %d, second at %d\n%s",
			posSep, posSecond, output)
	}
}

// TestRunChangelog_TwoSummariesExactlyOneSeparator verifies that two summaries
// produce exactly one "---" separator.
func TestRunChangelog_TwoSummariesExactlyOneSeparator(t *testing.T) {
	entries := []summaryEntry{
		{filename: "S-01-SUMMARY.md", content: "Alpha.\n"},
		{filename: "S-02-SUMMARY.md", content: "Beta.\n"},
	}

	setupChangelogProject(t, entries, true)

	var buf bytes.Buffer
	cmd.RunChangelog([]string{}, &buf)

	output := buf.String()

	count := strings.Count(output, "---")
	if count != 1 {
		t.Errorf("expected exactly 1 '---' separator for 2 entries, got %d:\n%s", count, output)
	}
}

// TestRunChangelog_ThreeSummariesExactlyTwoSeparators verifies that three
// summaries produce exactly two "---" separators.
func TestRunChangelog_ThreeSummariesExactlyTwoSeparators(t *testing.T) {
	entries := []summaryEntry{
		{filename: "S-01-SUMMARY.md", content: "Alpha.\n"},
		{filename: "S-02-SUMMARY.md", content: "Beta.\n"},
		{filename: "S-03-SUMMARY.md", content: "Gamma.\n"},
	}

	setupChangelogProject(t, entries, true)

	var buf bytes.Buffer
	cmd.RunChangelog([]string{}, &buf)

	output := buf.String()

	count := strings.Count(output, "---")
	if count != 2 {
		t.Errorf("expected exactly 2 '---' separators for 3 entries, got %d:\n%s", count, output)
	}
}
