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

// testSlice describes a single slice file to be created in the temp project.
type testSlice struct {
	id            string
	status        string
	step          string
	milestone     string
	started       string
	updated       string
	tests         int
	securityTests int
	deps          string
}

// sliceFrontmatter renders a slice's frontmatter as a string.
func sliceFrontmatter(slice testSlice) string {
	var builder strings.Builder
	builder.WriteString("---\n")
	builder.WriteString("id: " + slice.id + "\n")
	builder.WriteString("status: " + slice.status + "\n")
	builder.WriteString("step: " + slice.step + "\n")
	builder.WriteString("milestone: " + slice.milestone + "\n")
	builder.WriteString("started: " + slice.started + "\n")
	builder.WriteString("updated: " + slice.updated + "\n")

	testsStr := "0"
	if slice.tests > 0 {
		testsStr = itoa(slice.tests)
	}
	builder.WriteString("tests: " + testsStr + "\n")

	secStr := "0"
	if slice.securityTests > 0 {
		secStr = itoa(slice.securityTests)
	}
	builder.WriteString("security_tests: " + secStr + "\n")

	builder.WriteString("session:\n")
	builder.WriteString("deps: " + slice.deps + "\n")
	builder.WriteString("---\n")
	return builder.String()
}

// itoa converts an int to its decimal string representation without importing strconv.
func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	digits := []byte{}
	for value > 0 {
		digits = append([]byte{byte('0' + value%10)}, digits...)
		value /= 10
	}
	return string(digits)
}

// writeStatusFile writes content to a file path, failing the test on error.
func writeStatusFile(t *testing.T, path string, content string) {
	t.Helper()
	if writeError := os.WriteFile(path, []byte(content), 0o644); writeError != nil {
		t.Fatalf("writeStatusFile(%q) error: %v", path, writeError)
	}
}

// setupTestProject creates a temp directory mimicking a greenlight project:
//
//	tmpDir/
//	  .greenlight/
//	    slices/
//	      S-01.md  (with frontmatter)
//	      ...
//	    GRAPH.json
//
// It calls t.Chdir(tmpDir) so that RunStatus can locate .greenlight/ from the
// working directory. The original directory is automatically restored when the
// test completes.
//
// If graphJSON is empty, GRAPH.json is NOT written (simulates missing file).
func setupTestProject(t *testing.T, slices []testSlice, graphJSON string) string {
	t.Helper()

	tmpDir := t.TempDir()

	slicesDir := filepath.Join(tmpDir, ".greenlight", "slices")
	if mkdirError := os.MkdirAll(slicesDir, 0o755); mkdirError != nil {
		t.Fatalf("setupTestProject: failed to create slices dir: %v", mkdirError)
	}

	for _, slice := range slices {
		filename := slice.id + ".md"
		path := filepath.Join(slicesDir, filename)
		writeStatusFile(t, path, sliceFrontmatter(slice))
	}

	if graphJSON != "" {
		graphPath := filepath.Join(tmpDir, ".greenlight", "GRAPH.json")
		writeStatusFile(t, graphPath, graphJSON)
	}

	t.Chdir(tmpDir)
	return tmpDir
}

// minimalGraph builds a GRAPH.json with simple slice entries and no edges.
func minimalGraph(sliceIDs []string) string {
	var builder strings.Builder
	builder.WriteString(`{"slices":{`)
	for index, id := range sliceIDs {
		if index > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(`"`)
		builder.WriteString(id)
		builder.WriteString(`":{"name":"`)
		builder.WriteString(id)
		builder.WriteString(` Slice","depends_on":[],"wave":1,"contracts":[]}`)
	}
	builder.WriteString(`},"edges":[]}`)
	return builder.String()
}

// graphWithDeps builds a GRAPH.json where depMap maps a slice ID to its
// dependency IDs.
func graphWithDeps(entries map[string][]string) string {
	var builder strings.Builder
	builder.WriteString(`{"slices":{`)
	first := true
	for id, deps := range entries {
		if !first {
			builder.WriteString(",")
		}
		first = false
		builder.WriteString(`"`)
		builder.WriteString(id)
		builder.WriteString(`":{"name":"`)
		builder.WriteString(id)
		builder.WriteString(` Slice","depends_on":[`)
		for dIndex, dep := range deps {
			if dIndex > 0 {
				builder.WriteString(",")
			}
			builder.WriteString(`"`)
			builder.WriteString(dep)
			builder.WriteString(`"`)
		}
		builder.WriteString(`],"wave":1,"contracts":[]}`)
	}
	builder.WriteString(`},"edges":[]}`)
	return builder.String()
}

// ----------------------------------------------------------------------------
// C-98 — RunStatus default mode: happy paths
// ----------------------------------------------------------------------------

// TestRunStatus_MixedStatusesShowsAllSections verifies that a project with
// complete, in_progress, pending (ready), and pending (blocked) slices produces
// output containing Progress, Running, Ready, Blocked, and Tests sections.
func TestRunStatus_MixedStatusesShowsAllSections(t *testing.T) {
	slices := []testSlice{
		{id: "S-01", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z",
			tests: 10, securityTests: 2},
		{id: "S-02", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z",
			tests: 8, securityTests: 1},
		{id: "S-03", status: "in_progress", step: "implementing", milestone: "core",
			started: "2026-01-10", updated: "2026-01-10T00:00:00Z",
			tests: 5, securityTests: 0},
		{id: "S-04", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: "", tests: 0, securityTests: 0},
		{id: "S-05", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: "", tests: 0, securityTests: 0,
			deps: "S-03"},
	}

	graphJSON := graphWithDeps(map[string][]string{
		"S-01": {},
		"S-02": {},
		"S-03": {"S-01", "S-02"},
		"S-04": {},
		"S-05": {"S-03"},
	})

	setupTestProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunStatus([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	output := buf.String()

	if !strings.Contains(output, "Progress:") {
		t.Errorf("output missing 'Progress:' section:\n%s", output)
	}
	if !strings.Contains(output, "Running:") {
		t.Errorf("output missing 'Running:' section:\n%s", output)
	}
	if !strings.Contains(output, "Ready:") {
		t.Errorf("output missing 'Ready:' section:\n%s", output)
	}
	if !strings.Contains(output, "Blocked:") {
		t.Errorf("output missing 'Blocked:' section:\n%s", output)
	}
	if !strings.Contains(output, "Tests:") {
		t.Errorf("output missing 'Tests:' section:\n%s", output)
	}
}

// TestRunStatus_ProgressBarShowsCorrectRatio verifies that the progress bar
// reflects the completed/total ratio using ASCII fill characters.
// 2 complete out of 4 total = 50% fill.
func TestRunStatus_ProgressBarShowsCorrectRatio(t *testing.T) {
	slices := []testSlice{
		{id: "S-01", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z"},
		{id: "S-02", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z"},
		{id: "S-03", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
		{id: "S-04", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}

	graphJSON := minimalGraph([]string{"S-01", "S-02", "S-03", "S-04"})
	setupTestProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunStatus([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	output := buf.String()

	// Progress line must contain the count "2/4"
	if !strings.Contains(output, "2/4") {
		t.Errorf("expected '2/4' in progress output, got:\n%s", output)
	}

	// Progress bar must use only ASCII characters
	if strings.ContainsAny(output, "\x1b\x1B") {
		t.Errorf("output contains ANSI escape codes; progress bar must be ASCII-only")
	}
}

// TestRunStatus_RunningSlicesShowCurrentStep verifies that in_progress slices
// appear in the Running section with their current step in parentheses.
func TestRunStatus_RunningSlicesShowCurrentStep(t *testing.T) {
	slices := []testSlice{
		{id: "S-01", status: "in_progress", step: "implementing", milestone: "core",
			started: "2026-01-10", updated: "2026-01-10T00:00:00Z"},
		{id: "S-02", status: "in_progress", step: "testing", milestone: "core",
			started: "2026-01-11", updated: "2026-01-11T00:00:00Z"},
	}

	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupTestProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunStatus([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	output := buf.String()

	if !strings.Contains(output, "S-01") {
		t.Errorf("expected S-01 in running output:\n%s", output)
	}
	if !strings.Contains(output, "implementing") {
		t.Errorf("expected step 'implementing' in running output:\n%s", output)
	}
	if !strings.Contains(output, "S-02") {
		t.Errorf("expected S-02 in running output:\n%s", output)
	}
	if !strings.Contains(output, "testing") {
		t.Errorf("expected step 'testing' in running output:\n%s", output)
	}
}

// TestRunStatus_BlockedSlicesShowUnmetDependencies verifies that pending slices
// whose dependencies are not complete appear in the Blocked section with their
// unmet dependency IDs listed.
func TestRunStatus_BlockedSlicesShowUnmetDependencies(t *testing.T) {
	slices := []testSlice{
		{id: "S-01", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
		{id: "S-02", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: "", deps: "S-01"},
		{id: "S-03", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: "", deps: "S-01, S-02"},
	}

	graphJSON := graphWithDeps(map[string][]string{
		"S-01": {},
		"S-02": {"S-01"},
		"S-03": {"S-01", "S-02"},
	})
	setupTestProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunStatus([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	output := buf.String()

	if !strings.Contains(output, "Blocked:") {
		t.Errorf("expected Blocked section in output:\n%s", output)
	}

	// S-03 is blocked on both S-01 and S-02 (both pending)
	if !strings.Contains(output, "S-03") {
		t.Errorf("expected S-03 in blocked section:\n%s", output)
	}
}

// TestRunStatus_TestCountsSummedAcrossAllSlices verifies that the Tests section
// displays the sum of all slice test counts and security test counts.
func TestRunStatus_TestCountsSummedAcrossAllSlices(t *testing.T) {
	slices := []testSlice{
		{id: "S-01", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z",
			tests: 100, securityTests: 10},
		{id: "S-02", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z",
			tests: 200, securityTests: 20},
		{id: "S-03", status: "in_progress", step: "implementing", milestone: "core",
			started: "2026-01-10", updated: "2026-01-10T00:00:00Z",
			tests: 50, securityTests: 5},
	}

	graphJSON := minimalGraph([]string{"S-01", "S-02", "S-03"})
	setupTestProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunStatus([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	output := buf.String()

	// Total tests: 100 + 200 + 50 = 350
	if !strings.Contains(output, "350") {
		t.Errorf("expected total test count 350 in output:\n%s", output)
	}

	// Total security tests: 10 + 20 + 5 = 35
	if !strings.Contains(output, "35") {
		t.Errorf("expected total security test count 35 in output:\n%s", output)
	}
}

// TestRunStatus_OutputContainsProgressLineWithASCIIBar verifies that the output
// contains a "Progress:" label and an ASCII-only progress bar (using characters
// like '#', '.', '[', ']').
func TestRunStatus_OutputContainsProgressLineWithASCIIBar(t *testing.T) {
	slices := []testSlice{
		{id: "S-01", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z"},
		{id: "S-02", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}

	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupTestProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunStatus([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	output := buf.String()

	if !strings.Contains(output, "Progress:") {
		t.Errorf("expected 'Progress:' label in output:\n%s", output)
	}

	// Verify the bar uses ASCII bracket characters only, no ANSI codes
	if strings.ContainsAny(output, "\x1b\x1B") {
		t.Errorf("progress bar contains ANSI escape codes; must be ASCII-only:\n%s", output)
	}

	// Verify at least one bar character is present ([ or # or . or ])
	hasBarCharacter := strings.ContainsAny(output, "[#.]")
	if !hasBarCharacter {
		t.Errorf("expected ASCII bar characters ([, #, ., ]) in progress output:\n%s", output)
	}
}

// TestRunStatus_ZeroStateHandledGracefully verifies that a project with no
// started slices (all pending) still produces valid output without panicking.
func TestRunStatus_ZeroStateHandledGracefully(t *testing.T) {
	slices := []testSlice{
		{id: "S-01", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
		{id: "S-02", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}

	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupTestProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunStatus([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 for zero-state project, got %d", exitCode)
	}

	output := buf.String()

	// Must contain progress even at zero
	if !strings.Contains(output, "0/2") {
		t.Errorf("expected '0/2' in zero-state output:\n%s", output)
	}
}

// TestRunStatus_AllCompleteShowsFullProgress verifies that when all slices are
// complete, the count reflects total/total.
func TestRunStatus_AllCompleteShowsFullProgress(t *testing.T) {
	slices := []testSlice{
		{id: "S-01", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z",
			tests: 5, securityTests: 1},
		{id: "S-02", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-02", updated: "2026-01-02T00:00:00Z",
			tests: 3, securityTests: 0},
	}

	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupTestProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunStatus([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	output := buf.String()

	if !strings.Contains(output, "2/2") {
		t.Errorf("expected '2/2' in fully-complete output:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-98 — RunStatus default mode: error cases
// ----------------------------------------------------------------------------

// TestRunStatus_NoGreenlightDir_ReturnsExitCode1 verifies that when there is no
// .greenlight/ directory in the working directory or any parent, RunStatus
// prints the expected error message and returns exit code 1.
func TestRunStatus_NoGreenlightDir_ReturnsExitCode1(t *testing.T) {
	// Chdir to a temp dir that has NO .greenlight/ directory.
	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	var buf bytes.Buffer
	exitCode := cmd.RunStatus([]string{}, &buf)

	if exitCode != 1 {
		t.Errorf("expected exit code 1 when .greenlight/ not found, got %d", exitCode)
	}

	output := buf.String()
	if !strings.Contains(output, "greenlight") {
		t.Errorf("expected error message mentioning greenlight, got:\n%s", output)
	}
}

// TestRunStatus_NoGreenlightDir_PrintsInitHint verifies that the error output
// for a missing .greenlight/ directory suggests running 'greenlight init'.
func TestRunStatus_NoGreenlightDir_PrintsInitHint(t *testing.T) {
	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	var buf bytes.Buffer
	cmd.RunStatus([]string{}, &buf)

	output := buf.String()
	if !strings.Contains(output, "init") {
		t.Errorf("expected hint to run init, got:\n%s", output)
	}
}

// TestRunStatus_EmptySlicesDir_ReturnsExitCode1 verifies that when
// .greenlight/slices/ exists but contains no .md files, RunStatus prints the
// expected message and returns exit code 1.
func TestRunStatus_EmptySlicesDir_ReturnsExitCode1(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .greenlight/slices/ but leave it empty.
	slicesDir := filepath.Join(tmpDir, ".greenlight", "slices")
	if mkdirError := os.MkdirAll(slicesDir, 0o755); mkdirError != nil {
		t.Fatalf("failed to create slices dir: %v", mkdirError)
	}

	// Write a GRAPH.json so that failure is purely due to empty slices dir.
	graphPath := filepath.Join(tmpDir, ".greenlight", "GRAPH.json")
	writeStatusFile(t, graphPath, `{"slices":{},"edges":[]}`)

	t.Chdir(tmpDir)

	var buf bytes.Buffer
	exitCode := cmd.RunStatus([]string{}, &buf)

	if exitCode != 1 {
		t.Errorf("expected exit code 1 for empty slices directory, got %d", exitCode)
	}

	output := buf.String()
	if !strings.Contains(output, "slice") {
		t.Errorf("expected error message mentioning slices, got:\n%s", output)
	}
}

// TestRunStatus_EmptySlicesDir_PrintsSetupHint verifies that the error output
// for an empty slices directory hints at running init.
func TestRunStatus_EmptySlicesDir_PrintsSetupHint(t *testing.T) {
	tmpDir := t.TempDir()

	slicesDir := filepath.Join(tmpDir, ".greenlight", "slices")
	if mkdirError := os.MkdirAll(slicesDir, 0o755); mkdirError != nil {
		t.Fatalf("failed to create slices dir: %v", mkdirError)
	}

	graphPath := filepath.Join(tmpDir, ".greenlight", "GRAPH.json")
	writeStatusFile(t, graphPath, `{"slices":{},"edges":[]}`)

	t.Chdir(tmpDir)

	var buf bytes.Buffer
	cmd.RunStatus([]string{}, &buf)

	output := buf.String()
	if !strings.Contains(output, "init") {
		t.Errorf("expected hint to run init in empty-slices error, got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-98 — RunStatus default mode: graceful degradation
// ----------------------------------------------------------------------------

// TestRunStatus_MissingGraphJSON_ShowsStatusWithoutDependencyInfo verifies
// that when GRAPH.json is absent, RunStatus still shows slice status but does
// not return exit code 1. It should warn the user that dependency info is
// unavailable.
func TestRunStatus_MissingGraphJSON_ShowsStatusWithoutDependencyInfo(t *testing.T) {
	slices := []testSlice{
		{id: "S-01", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z",
			tests: 10, securityTests: 1},
		{id: "S-02", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}

	// Pass empty graphJSON so setupTestProject skips writing GRAPH.json.
	setupTestProject(t, slices, "")

	var buf bytes.Buffer
	exitCode := cmd.RunStatus([]string{}, &buf)

	// Exit code 0: missing graph is degraded mode, not a fatal error.
	if exitCode != 0 {
		t.Errorf("expected exit code 0 for missing GRAPH.json (degraded mode), got %d", exitCode)
	}

	output := buf.String()

	// Still shows Progress section.
	if !strings.Contains(output, "Progress:") {
		t.Errorf("expected Progress section even when GRAPH.json missing:\n%s", output)
	}

	// Should warn the user that graph/dependency info is unavailable.
	lowerOutput := strings.ToLower(output)
	hasWarning := strings.Contains(lowerOutput, "graph") ||
		strings.Contains(lowerOutput, "depend") ||
		strings.Contains(lowerOutput, "warn") ||
		strings.Contains(lowerOutput, "unavailable") ||
		strings.Contains(lowerOutput, "missing")

	if !hasWarning {
		t.Errorf("expected warning about missing graph/dependency info:\n%s", output)
	}
}

// TestRunStatus_MissingGraphJSON_StillShowsSliceCount verifies that even with
// no GRAPH.json, the progress count still reflects the actual slices found.
func TestRunStatus_MissingGraphJSON_StillShowsSliceCount(t *testing.T) {
	slices := []testSlice{
		{id: "S-01", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z"},
		{id: "S-02", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-02", updated: "2026-01-02T00:00:00Z"},
		{id: "S-03", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}

	setupTestProject(t, slices, "")

	var buf bytes.Buffer
	exitCode := cmd.RunStatus([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	output := buf.String()

	// 2 complete of 3 total.
	if !strings.Contains(output, "2/3") {
		t.Errorf("expected '2/3' in output even when GRAPH.json missing:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-99 — RunStatus --compact flag: happy paths
// ----------------------------------------------------------------------------

// TestRunStatus_Compact_OutputsSingleLine verifies that --compact mode outputs
// a single line in the format "{N}/{M} done | {K} running".
func TestRunStatus_Compact_OutputsSingleLine(t *testing.T) {
	slices := []testSlice{
		{id: "S-01", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z"},
		{id: "S-02", status: "in_progress", step: "implementing", milestone: "core",
			started: "2026-01-10", updated: "2026-01-10T00:00:00Z"},
		{id: "S-03", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}

	graphJSON := minimalGraph([]string{"S-01", "S-02", "S-03"})
	setupTestProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunStatus([]string{"--compact"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	output := buf.String()

	// Strip the trailing newline; remaining content must be a single line.
	trimmed := strings.TrimRight(output, "\n")
	if strings.Contains(trimmed, "\n") {
		t.Errorf("--compact output must be a single line, got:\n%q", output)
	}

	// Must contain the "done" marker.
	if !strings.Contains(output, "done") {
		t.Errorf("--compact output missing 'done' marker:\n%q", output)
	}

	// Must contain the "running" marker.
	if !strings.Contains(output, "running") {
		t.Errorf("--compact output missing 'running' marker:\n%q", output)
	}

	// Must contain the pipe separator.
	if !strings.Contains(output, "|") {
		t.Errorf("--compact output missing '|' separator:\n%q", output)
	}
}

// TestRunStatus_Compact_ShowsCorrectCounts verifies the exact counts in compact
// output: 1 complete, 1 in_progress, 1 pending → "1/3 done | 1 running".
func TestRunStatus_Compact_ShowsCorrectCounts(t *testing.T) {
	slices := []testSlice{
		{id: "S-01", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z"},
		{id: "S-02", status: "in_progress", step: "implementing", milestone: "core",
			started: "2026-01-10", updated: "2026-01-10T00:00:00Z"},
		{id: "S-03", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}

	graphJSON := minimalGraph([]string{"S-01", "S-02", "S-03"})
	setupTestProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunStatus([]string{"--compact"}, &buf)

	output := buf.String()

	if !strings.Contains(output, "1/3") {
		t.Errorf("expected '1/3' in compact output, got:\n%q", output)
	}

	if !strings.Contains(output, "1 running") {
		t.Errorf("expected '1 running' in compact output, got:\n%q", output)
	}
}

// TestRunStatus_Compact_AllComplete verifies the compact format when all slices
// are complete: "4/4 done | 0 running".
func TestRunStatus_Compact_AllComplete(t *testing.T) {
	slices := []testSlice{
		{id: "S-01", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z"},
		{id: "S-02", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-02", updated: "2026-01-02T00:00:00Z"},
		{id: "S-03", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-03", updated: "2026-01-03T00:00:00Z"},
		{id: "S-04", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-04", updated: "2026-01-04T00:00:00Z"},
	}

	graphJSON := minimalGraph([]string{"S-01", "S-02", "S-03", "S-04"})
	setupTestProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunStatus([]string{"--compact"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	output := buf.String()

	if !strings.Contains(output, "4/4") {
		t.Errorf("expected '4/4' in all-complete compact output, got:\n%q", output)
	}

	if !strings.Contains(output, "0 running") {
		t.Errorf("expected '0 running' in all-complete compact output, got:\n%q", output)
	}
}

// TestRunStatus_Compact_FormatHasNoExtraNewlines verifies that --compact output
// ends with exactly one newline and contains no embedded newlines.
func TestRunStatus_Compact_FormatHasNoExtraNewlines(t *testing.T) {
	slices := []testSlice{
		{id: "S-01", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z"},
	}

	graphJSON := minimalGraph([]string{"S-01"})
	setupTestProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunStatus([]string{"--compact"}, &buf)

	output := buf.String()

	// Must end with exactly one newline.
	if !strings.HasSuffix(output, "\n") {
		t.Errorf("compact output must end with a newline:\n%q", output)
	}

	// The content before the trailing newline must contain no other newlines.
	body := strings.TrimRight(output, "\n")
	if strings.Contains(body, "\n") {
		t.Errorf("compact output contains embedded newlines:\n%q", output)
	}
}

// TestRunStatus_Compact_NoColorCodes verifies that --compact output contains no
// ANSI escape sequences, making it safe for use in tmux status bars.
func TestRunStatus_Compact_NoColorCodes(t *testing.T) {
	slices := []testSlice{
		{id: "S-01", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z"},
		{id: "S-02", status: "in_progress", step: "implementing", milestone: "core",
			started: "2026-01-10", updated: "2026-01-10T00:00:00Z"},
	}

	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupTestProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunStatus([]string{"--compact"}, &buf)

	output := buf.String()

	if strings.ContainsAny(output, "\x1b\x1B") {
		t.Errorf("--compact output must have no ANSI color codes (tmux compatible):\n%q", output)
	}
}

// ----------------------------------------------------------------------------
// C-99 — RunStatus --compact flag: error / resilience cases
// ----------------------------------------------------------------------------

// TestRunStatus_Compact_NoGreenlightDir_OutputsPlaceholder verifies that when
// .greenlight/ is not found and --compact is set, the output is the placeholder
// string "? slices | ? running" and the exit code is 0 (tmux bar resilience).
func TestRunStatus_Compact_NoGreenlightDir_OutputsPlaceholder(t *testing.T) {
	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	var buf bytes.Buffer
	exitCode := cmd.RunStatus([]string{"--compact"}, &buf)

	// Contract C-99: DataReadError → exit code 0 for tmux resilience.
	if exitCode != 0 {
		t.Errorf("expected exit code 0 for --compact with no .greenlight/ (tmux resilience), got %d", exitCode)
	}

	output := buf.String()

	// Must output the placeholder format.
	if !strings.Contains(output, "?") {
		t.Errorf("expected placeholder '?' in compact error output, got:\n%q", output)
	}
}

// TestRunStatus_Compact_NoGreenlightDir_PlaceholderIsSingleLine verifies that
// even in error mode, --compact output is a single line.
func TestRunStatus_Compact_NoGreenlightDir_PlaceholderIsSingleLine(t *testing.T) {
	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	var buf bytes.Buffer
	cmd.RunStatus([]string{"--compact"}, &buf)

	output := buf.String()

	// Must be a single line even in error/placeholder mode.
	trimmed := strings.TrimRight(output, "\n")
	if strings.Contains(trimmed, "\n") {
		t.Errorf("--compact placeholder must be a single line, got:\n%q", output)
	}
}

// ----------------------------------------------------------------------------
// Output written to provided io.Writer
// ----------------------------------------------------------------------------

// TestRunStatus_WritesToProvidedWriter verifies that all RunStatus output goes
// to the provided io.Writer and not to stdout/stderr.
func TestRunStatus_WritesToProvidedWriter(t *testing.T) {
	slices := []testSlice{
		{id: "S-01", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z"},
	}

	graphJSON := minimalGraph([]string{"S-01"})
	setupTestProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunStatus([]string{}, &buf)

	if buf.Len() == 0 {
		t.Error("expected output to be written to the provided writer, got empty buffer")
	}
}

// TestRunStatus_Compact_WritesToProvidedWriter verifies the same for --compact.
func TestRunStatus_Compact_WritesToProvidedWriter(t *testing.T) {
	slices := []testSlice{
		{id: "S-01", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z"},
	}

	graphJSON := minimalGraph([]string{"S-01"})
	setupTestProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunStatus([]string{"--compact"}, &buf)

	if buf.Len() == 0 {
		t.Error("expected --compact output to be written to the provided writer, got empty buffer")
	}
}

// ----------------------------------------------------------------------------
// ReadySlices section
// ----------------------------------------------------------------------------

// TestRunStatus_ReadySectionListsUnblockedPendingSlices verifies that pending
// slices whose graph dependencies are all satisfied appear in the Ready section.
func TestRunStatus_ReadySectionListsUnblockedPendingSlices(t *testing.T) {
	slices := []testSlice{
		{id: "S-01", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z"},
		{id: "S-02", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: "", deps: "S-01"},
	}

	graphJSON := graphWithDeps(map[string][]string{
		"S-01": {},
		"S-02": {"S-01"},
	})
	setupTestProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunStatus([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	output := buf.String()

	if !strings.Contains(output, "Ready:") {
		t.Errorf("expected Ready section in output:\n%s", output)
	}

	// S-02 depends on S-01 which is complete, so S-02 should be ready.
	if !strings.Contains(output, "S-02") {
		t.Errorf("expected S-02 listed as ready:\n%s", output)
	}
}
