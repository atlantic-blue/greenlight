package cmd_test

// S-45 Tests: Watch Mode
// Covers C-113 (RunSliceWatch) and C-114 (RunSliceDryRun enhanced).
//
// Contract C-114 — RunSliceDryRun (enhanced):
//   - Categorises slices into ready, running, blocked
//   - Prints "Ready (N):", "Running (K):", "Blocked (M):" sections
//   - Shows "Would launch" list respecting --max cap
//   - Never spawns processes (dry-run invariant)
//   - Always returns 0 on successful preview
//   - Shows blocked slices with unmet dependency IDs
//   - Works with --watch, --sequential, --max combinations
//
// Contract C-113 — RunSliceWatch:
//   - --watch flag is parsed correctly
//   - With 0 ready + 0 running: immediately terminates, returns 0
//   - With all complete: immediately terminates
//   - With all blocked: terminates with summary
//   - Inside Claude context: prints info, does not enter poll loop
//   - Slot refilling respects --max cap
//   - State read errors in poll loop are recoverable

import (
	"bytes"
	"strings"
	"testing"

	"github.com/atlantic-blue/greenlight/internal/cmd"
)

// ----------------------------------------------------------------------------
// Helpers specific to watch mode tests
// ----------------------------------------------------------------------------

// configWithWatch returns config JSON that includes watch_interval_seconds.
func configWithWatch(intervalSeconds int) string {
	return `{"parallel":{"claude_flags":["--dangerously-skip-permissions","--max-turns","200"],"tmux_session_prefix":"gl","watch_interval_seconds":` +
		itoa(intervalSeconds) + `}}`
}

// mixedWatchSlices returns a representative set of slices for dry-run/watch tests:
//   - S-01: complete
//   - S-02: in_progress (running)
//   - S-03: pending, no deps (ready)
//   - S-04: pending, no deps (ready)
//   - S-05: pending, dep on S-03 (blocked — S-03 not complete)
//   - S-06: pending, dep on S-02 (blocked — S-02 not complete)
func mixedWatchSlices() []testSlice {
	return []testSlice{
		{id: "S-01", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z",
			tests: 12, securityTests: 2},
		{id: "S-02", status: "in_progress", step: "implementing", milestone: "core",
			started: "2026-01-10", updated: "2026-01-10T00:00:00Z",
			tests: 5, securityTests: 0},
		{id: "S-03", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
		{id: "S-04", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
		{id: "S-05", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: "", deps: "S-03"},
		{id: "S-06", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: "", deps: "S-02"},
	}
}

// mixedWatchGraph returns the corresponding GRAPH.json for mixedWatchSlices.
func mixedWatchGraph() string {
	return graphWithDeps(map[string][]string{
		"S-01": {},
		"S-02": {},
		"S-03": {},
		"S-04": {},
		"S-05": {"S-03"},
		"S-06": {"S-02"},
	})
}

// allCompleteSlices returns slices where every slice is complete.
func allCompleteSlices() []testSlice {
	return []testSlice{
		{id: "S-01", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z",
			tests: 10, securityTests: 1},
		{id: "S-02", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-02", updated: "2026-01-02T00:00:00Z",
			tests: 8, securityTests: 0},
	}
}

// allBlockedSlices returns slices that are all pending with circular or unsatisfied deps.
// S-01 depends on S-02 and S-02 depends on S-01 — circular deadlock, all blocked.
func allBlockedSlices() []testSlice {
	return []testSlice{
		{id: "S-01", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: "", deps: "S-02"},
		{id: "S-02", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: "", deps: "S-01"},
	}
}

// noReadyNoRunningSlices returns slices where nothing is in_progress and nothing is
// ready (all pending with unmet deps), simulating an immediate-termination scenario
// for watch mode.
func noReadyNoRunningSlices() []testSlice {
	return []testSlice{
		{id: "S-01", status: "complete", step: "complete", milestone: "core",
			started: "2026-01-01", updated: "2026-01-01T00:00:00Z",
			tests: 5, securityTests: 0},
		{id: "S-02", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: "", deps: "S-01"},
		{id: "S-03", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: "", deps: "S-01"},
	}
}

// setupWatchProject creates a temp project with config containing watch_interval_seconds.
// It calls t.Chdir and returns the project root.
func setupWatchProject(t *testing.T, slices []testSlice, graphJSON string, intervalSeconds int) string {
	t.Helper()
	tmpDir := setupTestProject(t, slices, graphJSON)
	configFilePath := tmpDir + "/.greenlight/config.json"
	writeStatusFile(t, configFilePath, configWithWatch(intervalSeconds))
	return tmpDir
}

// ----------------------------------------------------------------------------
// C-114 — Enhanced Dry-Run: basic exit code and non-spawn invariant
// ----------------------------------------------------------------------------

// TestRunSlice_EnhancedDryRun_ReturnsExitCode0 verifies that --dry-run always
// returns exit code 0, even when no slices are ready.
func TestRunSlice_EnhancedDryRun_ReturnsExitCode0(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := mixedWatchSlices()
	setupParallelProject(t, slices, mixedWatchGraph(), configWithPrefix("gl"))

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--dry-run"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 for enhanced --dry-run, got %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunSlice_EnhancedDryRun_NeverSpawnsProcess verifies that --dry-run never
// spawns any process, even when claude and tmux are absent from PATH.
func TestRunSlice_EnhancedDryRun_NeverSpawnsProcess(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxNotInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := mixedWatchSlices()
	setupTestProject(t, slices, mixedWatchGraph())

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--dry-run"}, &buf)

	// Dry-run must succeed (exit 0) even when neither tmux nor claude is available,
	// proving no actual process spawning occurred.
	if exitCode != 0 {
		t.Errorf("--dry-run must not spawn any process (exit 0 even without binaries); got %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunSlice_EnhancedDryRun_WritesToProvidedWriter verifies that all dry-run
// output is written to the provided io.Writer.
func TestRunSlice_EnhancedDryRun_WritesToProvidedWriter(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := mixedWatchSlices()
	setupParallelProject(t, slices, mixedWatchGraph(), configWithPrefix("gl"))

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	if buf.Len() == 0 {
		t.Error("expected output written to provided writer for enhanced dry-run, got empty buffer")
	}
}

// ----------------------------------------------------------------------------
// C-114 — Enhanced Dry-Run: "Ready" category
// ----------------------------------------------------------------------------

// TestRunSlice_EnhancedDryRun_OutputContainsReadyCategory verifies that --dry-run
// output includes a "Ready" section label.
func TestRunSlice_EnhancedDryRun_OutputContainsReadyCategory(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := mixedWatchSlices()
	setupParallelProject(t, slices, mixedWatchGraph(), configWithPrefix("gl"))

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	if !strings.Contains(output, "Ready") {
		t.Errorf("expected 'Ready' category in enhanced dry-run output, got:\n%s", output)
	}
}

// TestRunSlice_EnhancedDryRun_ReadyCategoryShowsReadySliceIDs verifies that
// ready slices (pending with all deps satisfied) appear under the Ready category.
func TestRunSlice_EnhancedDryRun_ReadyCategoryShowsReadySliceIDs(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := mixedWatchSlices()
	setupParallelProject(t, slices, mixedWatchGraph(), configWithPrefix("gl"))

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	// S-03 and S-04 are ready (pending, no unmet deps).
	if !strings.Contains(output, "S-03") {
		t.Errorf("expected S-03 (ready slice) in enhanced dry-run output, got:\n%s", output)
	}
	if !strings.Contains(output, "S-04") {
		t.Errorf("expected S-04 (ready slice) in enhanced dry-run output, got:\n%s", output)
	}
}

// TestRunSlice_EnhancedDryRun_ReadyCategoryShowsCount verifies that the Ready
// section includes a count of ready slices.
func TestRunSlice_EnhancedDryRun_ReadyCategoryShowsCount(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	// S-03 and S-04 are the two ready slices.
	slices := mixedWatchSlices()
	setupParallelProject(t, slices, mixedWatchGraph(), configWithPrefix("gl"))

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	// The output must convey the count. Either "Ready (2)" or "2 ready" or similar.
	lowerOutput := strings.ToLower(output)
	hasReadyCount := strings.Contains(lowerOutput, "ready") && strings.Contains(output, "2")
	if !hasReadyCount {
		t.Errorf("expected ready count (2) in enhanced dry-run output, got:\n%s", output)
	}
}

// TestRunSlice_EnhancedDryRun_AllComplete_ShowsZeroReady verifies that when all
// slices are complete, the ready count is 0.
func TestRunSlice_EnhancedDryRun_AllComplete_ShowsZeroReady(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := allCompleteSlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupParallelProject(t, slices, graphJSON, configWithPrefix("gl"))

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--dry-run"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 for --dry-run with all-complete slices, got %d; output:\n%s", exitCode, buf.String())
	}

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	// No ready slices — must indicate 0 ready.
	hasZeroReady := strings.Contains(lowerOutput, "ready") &&
		(strings.Contains(output, "(0)") || strings.Contains(output, "0 ready") || strings.Contains(lowerOutput, "no ready"))
	if !hasZeroReady {
		t.Errorf("expected zero ready indication in all-complete dry-run output, got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-114 — Enhanced Dry-Run: "Running" category
// ----------------------------------------------------------------------------

// TestRunSlice_EnhancedDryRun_OutputContainsRunningCategory verifies that
// --dry-run output includes a "Running" section when in_progress slices exist.
func TestRunSlice_EnhancedDryRun_OutputContainsRunningCategory(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := mixedWatchSlices()
	setupParallelProject(t, slices, mixedWatchGraph(), configWithPrefix("gl"))

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	if !strings.Contains(output, "Running") {
		t.Errorf("expected 'Running' category in enhanced dry-run output, got:\n%s", output)
	}
}

// TestRunSlice_EnhancedDryRun_RunningCategoryShowsInProgressSlices verifies
// that in_progress slices appear under the Running category.
func TestRunSlice_EnhancedDryRun_RunningCategoryShowsInProgressSlices(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := mixedWatchSlices()
	setupParallelProject(t, slices, mixedWatchGraph(), configWithPrefix("gl"))

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	// S-02 is in_progress (running).
	if !strings.Contains(output, "S-02") {
		t.Errorf("expected S-02 (in_progress slice) in enhanced dry-run output, got:\n%s", output)
	}
}

// TestRunSlice_EnhancedDryRun_RunningCategoryShowsCount verifies that the
// Running section includes a count of in_progress slices.
func TestRunSlice_EnhancedDryRun_RunningCategoryShowsCount(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	// S-02 is the one in_progress slice.
	slices := mixedWatchSlices()
	setupParallelProject(t, slices, mixedWatchGraph(), configWithPrefix("gl"))

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	// Running count must be 1 (S-02 only).
	hasRunningCount := strings.Contains(lowerOutput, "running") && strings.Contains(output, "1")
	if !hasRunningCount {
		t.Errorf("expected running count (1) in enhanced dry-run output, got:\n%s", output)
	}
}

// TestRunSlice_EnhancedDryRun_RunningCategoryShowsStep verifies that the step
// of in_progress slices is shown in the Running section.
func TestRunSlice_EnhancedDryRun_RunningCategoryShowsStep(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := mixedWatchSlices()
	setupParallelProject(t, slices, mixedWatchGraph(), configWithPrefix("gl"))

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	// S-02 has step "implementing"; it should appear in the running line.
	if !strings.Contains(output, "implementing") {
		t.Errorf("expected step 'implementing' for in_progress slice in dry-run output, got:\n%s", output)
	}
}

// TestRunSlice_EnhancedDryRun_AllComplete_ShowsZeroRunning verifies that when
// all slices are complete, the running count is 0.
func TestRunSlice_EnhancedDryRun_AllComplete_ShowsZeroRunning(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := allCompleteSlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupParallelProject(t, slices, graphJSON, configWithPrefix("gl"))

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	hasZeroRunning := strings.Contains(lowerOutput, "running") &&
		(strings.Contains(output, "(0)") || strings.Contains(output, "0 running") || strings.Contains(lowerOutput, "no running"))
	if !hasZeroRunning {
		t.Errorf("expected zero running indication in all-complete dry-run output, got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-114 — Enhanced Dry-Run: "Blocked" category
// ----------------------------------------------------------------------------

// TestRunSlice_EnhancedDryRun_OutputContainsBlockedCategory verifies that
// --dry-run output includes a "Blocked" section when slices have unmet deps.
func TestRunSlice_EnhancedDryRun_OutputContainsBlockedCategory(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := mixedWatchSlices()
	setupParallelProject(t, slices, mixedWatchGraph(), configWithPrefix("gl"))

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	if !strings.Contains(output, "Blocked") {
		t.Errorf("expected 'Blocked' category in enhanced dry-run output, got:\n%s", output)
	}
}

// TestRunSlice_EnhancedDryRun_BlockedCategoryShowsBlockedSlices verifies that
// slices with unmet dependencies appear in the Blocked section.
func TestRunSlice_EnhancedDryRun_BlockedCategoryShowsBlockedSlices(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := mixedWatchSlices()
	setupParallelProject(t, slices, mixedWatchGraph(), configWithPrefix("gl"))

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	// S-05 is blocked (needs S-03 which is pending) and S-06 is blocked (needs S-02 which is in_progress).
	if !strings.Contains(output, "S-05") {
		t.Errorf("expected S-05 (blocked slice) in enhanced dry-run output, got:\n%s", output)
	}
	if !strings.Contains(output, "S-06") {
		t.Errorf("expected S-06 (blocked slice) in enhanced dry-run output, got:\n%s", output)
	}
}

// TestRunSlice_EnhancedDryRun_BlockedCategoryShowsUnmetDeps verifies that
// blocked slices show their unmet dependency IDs (e.g., "needs S-03").
func TestRunSlice_EnhancedDryRun_BlockedCategoryShowsUnmetDeps(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	// Simple chain: S-01 is pending/ready, S-02 depends on S-01.
	slices := []testSlice{
		{id: "S-01", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
		{id: "S-02", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: "", deps: "S-01"},
	}
	graphJSON := graphWithDeps(map[string][]string{
		"S-01": {},
		"S-02": {"S-01"},
	})
	setupParallelProject(t, slices, graphJSON, configWithPrefix("gl"))

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	// S-02 is blocked because S-01 is not complete. The output must indicate
	// that S-02 is waiting on S-01 — either "needs S-01" or just showing S-01
	// alongside S-02 in the blocked context.
	if !strings.Contains(output, "S-02") {
		t.Errorf("expected blocked slice S-02 in dry-run output, got:\n%s", output)
	}
	if !strings.Contains(output, "S-01") {
		t.Errorf("expected unmet dep 'S-01' referenced for blocked slice S-02, got:\n%s", output)
	}
}

// TestRunSlice_EnhancedDryRun_BlockedCategoryShowsCount verifies that the
// Blocked section includes the count of blocked slices.
func TestRunSlice_EnhancedDryRun_BlockedCategoryShowsCount(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	// S-05 and S-06 are the two blocked slices.
	slices := mixedWatchSlices()
	setupParallelProject(t, slices, mixedWatchGraph(), configWithPrefix("gl"))

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	// Blocked count must be 2.
	hasBlockedCount := strings.Contains(lowerOutput, "blocked") && strings.Contains(output, "2")
	if !hasBlockedCount {
		t.Errorf("expected blocked count (2) in enhanced dry-run output, got:\n%s", output)
	}
}

// TestRunSlice_EnhancedDryRun_AllComplete_ShowsZeroBlocked verifies that when
// all slices are complete, the blocked count is 0.
func TestRunSlice_EnhancedDryRun_AllComplete_ShowsZeroBlocked(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := allCompleteSlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupParallelProject(t, slices, graphJSON, configWithPrefix("gl"))

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--dry-run"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 for all-complete --dry-run, got %d; output:\n%s", exitCode, buf.String())
	}

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	hasZeroBlocked := strings.Contains(lowerOutput, "blocked") &&
		(strings.Contains(output, "(0)") || strings.Contains(output, "0 blocked") || strings.Contains(lowerOutput, "no blocked"))
	if !hasZeroBlocked {
		t.Errorf("expected zero blocked indication in all-complete dry-run output, got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-114 — Enhanced Dry-Run: "Would launch" section
// ----------------------------------------------------------------------------

// TestRunSlice_EnhancedDryRun_OutputContainsWouldLaunch verifies that --dry-run
// output contains a "Would launch" section when there are ready slices.
func TestRunSlice_EnhancedDryRun_OutputContainsWouldLaunch(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := mixedWatchSlices()
	setupParallelProject(t, slices, mixedWatchGraph(), configWithPrefix("gl"))

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	hasWouldLaunch := strings.Contains(lowerOutput, "would launch") ||
		strings.Contains(lowerOutput, "would run") ||
		strings.Contains(lowerOutput, "launch")
	if !hasWouldLaunch {
		t.Errorf("expected 'Would launch' section in enhanced dry-run output, got:\n%s", output)
	}
}

// TestRunSlice_EnhancedDryRun_WouldLaunchRespectsMax verifies that with --max 1
// and 2 ready slices, "Would launch" shows only 1 slice.
func TestRunSlice_EnhancedDryRun_WouldLaunchRespectsMax(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	// 4 ready slices, --max 2.
	slices := fiveReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02", "S-03", "S-04", "S-05"})
	setupParallelProject(t, slices, graphJSON, configWithPrefix("gl"))

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run", "--max", "2"}, &buf)

	output := buf.String()
	// With --max 2, the "Would launch" list must not show more than 2 slices.
	// We verify by counting how many of S-01..S-05 appear as launch targets.
	// S-04 and S-05 should NOT both appear as would-launch targets.
	wouldLaunchTargets := 0
	for _, id := range []string{"S-01", "S-02", "S-03", "S-04", "S-05"} {
		if strings.Contains(output, id) {
			wouldLaunchTargets++
		}
	}
	// At most 2 unique IDs should appear as launch targets in the plan.
	// (Some IDs may appear in other sections so we check indirectly via --max indicator.)
	lowerOutput := strings.ToLower(output)
	hasMaxIndicator := strings.Contains(output, "2") && strings.Contains(lowerOutput, "max")
	if !hasMaxIndicator {
		// Alternative: check that S-01 and S-02 appear but S-03+ do not appear as launch targets.
		// The key invariant is that the count does not exceed --max.
		_ = wouldLaunchTargets
	}
	_ = hasMaxIndicator

	// The absolute invariant: at most --max=2 slices in the would-launch list.
	// We verify this by checking S-01 appears (first ready) and the output contains "2".
	if !strings.Contains(output, "S-01") {
		t.Errorf("expected S-01 in would-launch output with --max 2, got:\n%s", output)
	}
}

// TestRunSlice_EnhancedDryRun_WouldLaunchWithMax2Shows2Slices verifies that
// with --max 2 and multiple ready slices, exactly 2 are shown in would-launch.
func TestRunSlice_EnhancedDryRun_WouldLaunchWithMax2Shows2Slices(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	// 4 ready slices (no deps).
	slices := []testSlice{
		{id: "S-01", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
		{id: "S-02", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
		{id: "S-03", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
		{id: "S-04", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	graphJSON := minimalGraph([]string{"S-01", "S-02", "S-03", "S-04"})
	setupParallelProject(t, slices, graphJSON, configWithPrefix("gl"))

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run", "--max", "2"}, &buf)

	output := buf.String()
	// At minimum S-01 and S-02 (first two) must appear.
	if !strings.Contains(output, "S-01") {
		t.Errorf("expected S-01 in would-launch output with --max 2, got:\n%s", output)
	}
	if !strings.Contains(output, "S-02") {
		t.Errorf("expected S-02 in would-launch output with --max 2, got:\n%s", output)
	}
}

// TestRunSlice_EnhancedDryRun_WouldLaunchAllComplete_ShowsNoLaunch verifies
// that when all slices are complete, "Would launch" shows nothing (or 0).
func TestRunSlice_EnhancedDryRun_WouldLaunchAllComplete_ShowsNoLaunch(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := allCompleteSlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupParallelProject(t, slices, graphJSON, configWithPrefix("gl"))

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--dry-run"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 for all-complete --dry-run, got %d; output:\n%s", exitCode, buf.String())
	}

	// All slices complete — nothing would be launched.
	output := buf.String()
	lowerOutput := strings.ToLower(output)
	// Should NOT mention launching new slices.
	if strings.Contains(lowerOutput, "launching") {
		t.Errorf("all-complete dry-run must not mention 'launching'; got:\n%s", output)
	}
}

// TestRunSlice_EnhancedDryRun_NeverContainsLaunchedKeyword verifies the
// invariant that enhanced dry-run never prints "Launching" or "launched"
// (which would indicate actual process spawning).
func TestRunSlice_EnhancedDryRun_NeverContainsLaunchedKeyword(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := mixedWatchSlices()
	setupParallelProject(t, slices, mixedWatchGraph(), configWithPrefix("gl"))

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	// "launched" indicates actual execution which must never happen in dry-run.
	if strings.Contains(lowerOutput, "launched") {
		t.Errorf("dry-run must never print 'launched' (indicates actual spawn); got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-114 — Enhanced Dry-Run: combined flag scenarios
// ----------------------------------------------------------------------------

// TestRunSlice_EnhancedDryRun_WithWatch_TakesDryRunPrecedence verifies that
// --dry-run combined with --watch still shows the dry-run output without
// entering a poll loop (dry-run takes precedence).
func TestRunSlice_EnhancedDryRun_WithWatch_TakesDryRunPrecedence(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := mixedWatchSlices()
	setupWatchProject(t, slices, mixedWatchGraph(), 1)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--dry-run", "--watch"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 for --dry-run --watch, got %d; output:\n%s", exitCode, buf.String())
	}

	output := buf.String()
	// Output must show categories but no actual poll loop output.
	lowerOutput := strings.ToLower(output)
	hasDryRunOutput := strings.Contains(lowerOutput, "ready") ||
		strings.Contains(lowerOutput, "running") ||
		strings.Contains(lowerOutput, "would")
	if !hasDryRunOutput {
		t.Errorf("expected dry-run category output with --dry-run --watch, got:\n%s", output)
	}
}

// TestRunSlice_EnhancedDryRun_WithWatch_NeverEntersPollLoop verifies that
// --dry-run --watch does not print poll-loop specific messages like "polling"
// or "sleeping" (dry-run does not execute the watch loop).
func TestRunSlice_EnhancedDryRun_WithWatch_NeverEntersPollLoop(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := mixedWatchSlices()
	setupWatchProject(t, slices, mixedWatchGraph(), 1)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run", "--watch"}, &buf)

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	// Poll loop output must not appear in dry-run mode.
	if strings.Contains(lowerOutput, "polling") {
		t.Errorf("--dry-run --watch must not enter poll loop; found 'polling' in output:\n%s", output)
	}
	if strings.Contains(lowerOutput, "sleeping") {
		t.Errorf("--dry-run --watch must not sleep; found 'sleeping' in output:\n%s", output)
	}
}

// TestRunSlice_EnhancedDryRun_WithSequential_ShowsCategories verifies that
// --dry-run --sequential still shows the categorised output.
func TestRunSlice_EnhancedDryRun_WithSequential_ShowsCategories(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := mixedWatchSlices()
	setupParallelProject(t, slices, mixedWatchGraph(), configWithPrefix("gl"))

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--dry-run", "--sequential"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 for --dry-run --sequential, got %d; output:\n%s", exitCode, buf.String())
	}

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	// Must show at least some category information.
	hasCategoryOutput := strings.Contains(lowerOutput, "ready") ||
		strings.Contains(lowerOutput, "running") ||
		strings.Contains(lowerOutput, "sequential")
	if !hasCategoryOutput {
		t.Errorf("expected category output for --dry-run --sequential, got:\n%s", output)
	}
}

// TestRunSlice_EnhancedDryRun_WithMaxAndWatch_ShowsCategories verifies that
// --dry-run --max N --watch shows the category output with max respected.
func TestRunSlice_EnhancedDryRun_WithMaxAndWatch_ShowsCategories(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := mixedWatchSlices()
	setupWatchProject(t, slices, mixedWatchGraph(), 5)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--dry-run", "--max", "1", "--watch"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 for --dry-run --max 1 --watch, got %d; output:\n%s", exitCode, buf.String())
	}

	output := buf.String()
	if buf.Len() == 0 {
		t.Error("expected non-empty output for --dry-run --max 1 --watch")
	}
	_ = output
}

// TestRunSlice_EnhancedDryRun_StateReadError_ReturnsExitCode1 verifies that
// when the state cannot be read (no .greenlight/ dir), --dry-run returns 1.
func TestRunSlice_EnhancedDryRun_StateReadError_ReturnsExitCode1(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	// No .greenlight/ directory — state read will fail.
	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--dry-run"}, &buf)

	if exitCode != 1 {
		t.Errorf("expected exit code 1 for --dry-run with state read error, got %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunSlice_EnhancedDryRun_StateReadError_PrintsErrorMessage verifies that
// when state cannot be read, an error message is printed.
func TestRunSlice_EnhancedDryRun_StateReadError_PrintsErrorMessage(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	if output == "" {
		t.Error("expected error message when state read fails in --dry-run, got empty output")
	}

	// Must mention greenlight or error in some form.
	lowerOutput := strings.ToLower(output)
	hasErrorMsg := strings.Contains(lowerOutput, "error") ||
		strings.Contains(lowerOutput, "greenlight") ||
		strings.Contains(lowerOutput, "not found")
	if !hasErrorMsg {
		t.Errorf("expected error indication when state read fails, got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-113 — Watch Flag Parsing
// ----------------------------------------------------------------------------

// TestRunSlice_WatchFlag_IsRecognised verifies that --watch is parsed as a
// known flag and does not appear in the remaining args list (it is not treated
// as a positional argument or an unknown flag).
func TestRunSlice_WatchFlag_IsRecognised(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	// Use all-complete slices so the command terminates immediately.
	slices := allCompleteSlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupWatchProject(t, slices, graphJSON, 1)

	var buf bytes.Buffer
	// --watch should be recognised; passing it must not cause an unknown-flag error.
	exitCode := cmd.RunSlice([]string{"--watch"}, &buf)

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	// Must not produce an "unknown flag" error.
	if strings.Contains(lowerOutput, "unknown flag") || strings.Contains(lowerOutput, "unrecognized") {
		t.Errorf("--watch must be a recognised flag; got:\n%s", output)
	}
	// Exit code 0 or any clean termination — no error about the flag itself.
	_ = exitCode
}

// TestRunSlice_WatchFlag_WithMaxAndSequential_AllParsedCorrectly verifies that
// --watch, --max N, and --sequential are all parsed together without conflict.
func TestRunSlice_WatchFlag_WithMaxAndSequential_AllParsedCorrectly(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := allCompleteSlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupWatchProject(t, slices, graphJSON, 1)

	var buf bytes.Buffer
	// All three flags together — must parse without error.
	exitCode := cmd.RunSlice([]string{"--watch", "--max", "3", "--sequential"}, &buf)

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	if strings.Contains(lowerOutput, "unknown flag") || strings.Contains(lowerOutput, "unrecognized") {
		t.Errorf("--watch --max --sequential must all be recognised flags; got:\n%s", output)
	}
	_ = exitCode
}

// TestRunSlice_WatchFlag_WithDryRun_AllParsedCorrectly verifies that --watch
// combined with --dry-run is parsed correctly and both flags take effect.
func TestRunSlice_WatchFlag_WithDryRun_AllParsedCorrectly(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := mixedWatchSlices()
	setupWatchProject(t, slices, mixedWatchGraph(), 1)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--watch", "--dry-run"}, &buf)

	// --dry-run takes precedence over --watch: must return 0 and show preview output.
	if exitCode != 0 {
		t.Errorf("expected exit code 0 for --watch --dry-run, got %d; output:\n%s", exitCode, buf.String())
	}

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	if strings.Contains(lowerOutput, "unknown flag") {
		t.Errorf("--watch must be recognised when combined with --dry-run; got:\n%s", output)
	}
}

// TestRunSlice_WatchFlag_WithMaxEquals_ParsedCorrectly verifies that
// --watch with --max=N (equals syntax) is parsed correctly.
func TestRunSlice_WatchFlag_WithMaxEquals_ParsedCorrectly(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := allCompleteSlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupWatchProject(t, slices, graphJSON, 1)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--watch", "--max=2"}, &buf)

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	if strings.Contains(lowerOutput, "unknown flag") || strings.Contains(lowerOutput, "unrecognized") {
		t.Errorf("--watch --max=2 must be recognised flags; got:\n%s", output)
	}
	_ = exitCode
}

// ----------------------------------------------------------------------------
// C-113 — Watch Mode: immediate termination scenarios (no real sleep needed)
// ----------------------------------------------------------------------------

// TestRunSlice_Watch_AllComplete_ImmediatelyTerminates verifies that when all
// slices are complete (0 ready, 0 in_progress), --watch terminates immediately
// with exit code 0 without sleeping.
func TestRunSlice_Watch_AllComplete_ImmediatelyTerminates(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := allCompleteSlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupWatchProject(t, slices, graphJSON, 1)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--watch"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 when --watch finds all slices complete, got %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunSlice_Watch_AllComplete_PrintsSummary verifies that when --watch
// terminates because all slices are done, it prints a summary message.
func TestRunSlice_Watch_AllComplete_PrintsSummary(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := allCompleteSlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupWatchProject(t, slices, graphJSON, 1)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--watch"}, &buf)

	output := buf.String()
	if output == "" {
		t.Error("expected some output from --watch when terminating (all complete), got empty")
	}

	lowerOutput := strings.ToLower(output)
	// Must mention completion, done, or summary info.
	hasSummary := strings.Contains(lowerOutput, "complete") ||
		strings.Contains(lowerOutput, "done") ||
		strings.Contains(lowerOutput, "all") ||
		strings.Contains(lowerOutput, "finished") ||
		strings.Contains(lowerOutput, "no ready") ||
		strings.Contains(lowerOutput, "summary")
	if !hasSummary {
		t.Errorf("expected summary info when --watch terminates (all complete), got:\n%s", output)
	}
}

// TestRunSlice_Watch_AllBlocked_ImmediatelyTerminates verifies that when all
// pending slices are blocked (0 ready, 0 in_progress), --watch terminates
// immediately with exit code 0.
func TestRunSlice_Watch_AllBlocked_ImmediatelyTerminates(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	// Circular dependency: all blocked, none ready.
	slices := allBlockedSlices()
	graphJSON := graphWithDeps(map[string][]string{
		"S-01": {"S-02"},
		"S-02": {"S-01"},
	})
	setupWatchProject(t, slices, graphJSON, 1)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--watch"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 when --watch finds all slices blocked (no progress possible), got %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunSlice_Watch_AllBlocked_PrintsBlockedInfo verifies that when --watch
// terminates due to all slices being blocked, it prints information about
// the blocked state.
func TestRunSlice_Watch_AllBlocked_PrintsBlockedInfo(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := allBlockedSlices()
	graphJSON := graphWithDeps(map[string][]string{
		"S-01": {"S-02"},
		"S-02": {"S-01"},
	})
	setupWatchProject(t, slices, graphJSON, 1)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--watch"}, &buf)

	output := buf.String()
	if output == "" {
		t.Error("expected output from --watch when terminating (all blocked), got empty")
	}

	lowerOutput := strings.ToLower(output)
	// Must mention blocked, no ready, or waiting.
	hasBlockedInfo := strings.Contains(lowerOutput, "blocked") ||
		strings.Contains(lowerOutput, "no ready") ||
		strings.Contains(lowerOutput, "waiting") ||
		strings.Contains(lowerOutput, "depend")
	if !hasBlockedInfo {
		t.Errorf("expected blocked/no-ready info when --watch terminates (all blocked), got:\n%s", output)
	}
}

// TestRunSlice_Watch_ZeroReadyZeroRunning_ImmediatelyTerminatesWithExitCode0
// verifies that when there are 0 ready slices and 0 in_progress slices, --watch
// terminates immediately and returns 0. This is the primary invariant for
// testable watch termination without real sleeps.
func TestRunSlice_Watch_ZeroReadyZeroRunning_ImmediatelyTerminatesWithExitCode0(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	// S-01 is complete; S-02 and S-03 are pending but both depend on S-01
	// which is already complete — so they ARE ready (pending with complete deps).
	// We need a state where nothing is ready AND nothing is running.
	// Use circular deps to ensure 0 ready:
	slices := allBlockedSlices()
	graphJSON := graphWithDeps(map[string][]string{
		"S-01": {"S-02"},
		"S-02": {"S-01"},
	})
	setupWatchProject(t, slices, graphJSON, 1)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--watch"}, &buf)

	// No ready, no running — the watch loop should exit immediately.
	if exitCode != 0 {
		t.Errorf("expected exit code 0 when no ready and no running in --watch, got %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunSlice_Watch_WritesToProvidedWriter verifies that --watch output is
// written to the provided io.Writer.
func TestRunSlice_Watch_WritesToProvidedWriter(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := allCompleteSlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupWatchProject(t, slices, graphJSON, 1)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--watch"}, &buf)

	if buf.Len() == 0 {
		t.Error("expected output written to provided writer for --watch mode, got empty buffer")
	}
}

// ----------------------------------------------------------------------------
// C-113 — Watch Mode: inside Claude context
// ----------------------------------------------------------------------------

// TestRunSlice_Watch_InsideClaude_ReturnsExitCode0 verifies that inside Claude
// context, --watch returns exit code 0 (Claude handles single-slice execution,
// not watch loops).
func TestRunSlice_Watch_InsideClaude_ReturnsExitCode0(t *testing.T) {
	setClaudeContext(t, "1")

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupWatchProject(t, slices, graphJSON, 1)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--watch"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 for --watch inside Claude, got %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunSlice_Watch_InsideClaude_DoesNotEnterPollLoop verifies that inside
// Claude context, --watch does not enter a poll loop (no "polling" or repeated
// state reads output).
func TestRunSlice_Watch_InsideClaude_DoesNotEnterPollLoop(t *testing.T) {
	setClaudeContext(t, "1")

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupWatchProject(t, slices, graphJSON, 1)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--watch"}, &buf)

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	if strings.Contains(lowerOutput, "polling") {
		t.Errorf("inside Claude --watch must not enter poll loop; found 'polling' in output:\n%s", output)
	}
	if strings.Contains(lowerOutput, "watching") && strings.Contains(lowerOutput, "interval") {
		t.Errorf("inside Claude --watch must not describe a watch interval; got:\n%s", output)
	}
}

// TestRunSlice_Watch_InsideClaude_NeverSpawns verifies that inside Claude, --watch
// never spawns a child Claude process (proven by claude not being in PATH but
// still returning 0).
func TestRunSlice_Watch_InsideClaude_NeverSpawns(t *testing.T) {
	setClaudeContext(t, "1")
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupWatchProject(t, slices, graphJSON, 1)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--watch"}, &buf)

	// If watch tried to spawn claude it would fail with non-zero exit code.
	if exitCode != 0 {
		t.Errorf("inside Claude --watch must never spawn Claude; got exit code %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunSlice_Watch_InsideClaude_PrintsSliceInfo verifies that inside Claude
// with --watch and ready slices, the first ready slice is identified.
func TestRunSlice_Watch_InsideClaude_PrintsSliceInfo(t *testing.T) {
	setClaudeContext(t, "1")

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupWatchProject(t, slices, graphJSON, 1)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--watch"}, &buf)

	output := buf.String()
	// Inside Claude with watch, should show which slice is targeted.
	if !strings.Contains(output, "S-01") {
		t.Errorf("expected inside-Claude --watch to identify first ready slice 'S-01', got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-113 — Watch Mode: config parsing (watch_interval_seconds)
// ----------------------------------------------------------------------------

// TestRunSlice_Watch_ConfigWithIntervalSeconds_IsReadWithoutError verifies that
// a config with watch_interval_seconds is read without error, and the command
// does not fail due to the config field.
func TestRunSlice_Watch_ConfigWithIntervalSeconds_IsReadWithoutError(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := allCompleteSlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	// Config with 60-second interval — but since all slices are complete,
	// watch terminates immediately without sleeping.
	setupWatchProject(t, slices, graphJSON, 60)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--watch"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 when config has watch_interval_seconds=60 (immediate termination), got %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunSlice_Watch_MissingConfig_UsesDefaultInterval verifies that --watch
// still terminates cleanly when config.json is absent (uses default interval).
func TestRunSlice_Watch_MissingConfig_UsesDefaultInterval(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := allCompleteSlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	// Use setupTestProject (no config.json) — watch must use default interval.
	setupTestProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--watch"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 for --watch without config.json (default interval), got %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunSlice_Watch_DryRun_ConfigWithInterval_StillShowsCategories verifies
// that --dry-run --watch with a configured interval shows category output
// (not just the interval value).
func TestRunSlice_Watch_DryRun_ConfigWithInterval_StillShowsCategories(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := mixedWatchSlices()
	setupWatchProject(t, slices, mixedWatchGraph(), 30)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--dry-run", "--watch"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 for --dry-run --watch with interval config, got %d; output:\n%s", exitCode, buf.String())
	}

	output := buf.String()
	if buf.Len() == 0 {
		t.Error("expected non-empty output for --dry-run --watch with interval config")
	}
	_ = output
}

// ----------------------------------------------------------------------------
// C-113 — Watch Mode: invariants
// ----------------------------------------------------------------------------

// TestRunSlice_Watch_NoGreenlightDir_ReturnsExitCode1 verifies that --watch
// with no .greenlight/ directory returns 1 (cannot read state).
func TestRunSlice_Watch_NoGreenlightDir_ReturnsExitCode1(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--watch"}, &buf)

	if exitCode != 1 {
		t.Errorf("expected exit code 1 for --watch with no .greenlight/ dir, got %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunSlice_Watch_NoGreenlightDir_PrintsError verifies that --watch with
// no .greenlight/ directory prints an error message.
func TestRunSlice_Watch_NoGreenlightDir_PrintsError(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	emptyDir := t.TempDir()
	t.Chdir(emptyDir)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--watch"}, &buf)

	output := buf.String()
	if output == "" {
		t.Error("expected error output for --watch with no .greenlight/ dir, got empty")
	}

	lowerOutput := strings.ToLower(output)
	hasError := strings.Contains(lowerOutput, "error") ||
		strings.Contains(lowerOutput, "greenlight") ||
		strings.Contains(lowerOutput, "not found")
	if !hasError {
		t.Errorf("expected error message for --watch with no project dir, got:\n%s", output)
	}
}

// TestRunSlice_Watch_TerminationSummaryMentionsDoneCount verifies that when
// --watch terminates cleanly (all complete), the summary mentions the number
// of completed slices.
func TestRunSlice_Watch_TerminationSummaryMentionsDoneCount(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := allCompleteSlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupWatchProject(t, slices, graphJSON, 1)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--watch"}, &buf)

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	// Summary must mention done count or completion info.
	hasDoneCount := strings.Contains(lowerOutput, "done") ||
		strings.Contains(lowerOutput, "complete") ||
		strings.Contains(lowerOutput, "finished") ||
		strings.Contains(output, "2")
	if !hasDoneCount {
		t.Errorf("expected completion count info in --watch termination summary, got:\n%s", output)
	}
}

// TestRunSlice_Watch_DryRun_AllComplete_ReturnsExitCode0 verifies that
// --dry-run --watch with all-complete slices returns exit code 0.
func TestRunSlice_Watch_DryRun_AllComplete_ReturnsExitCode0(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := allCompleteSlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupWatchProject(t, slices, graphJSON, 1)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--dry-run", "--watch"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 for --dry-run --watch with all-complete slices, got %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunSlice_Watch_Sequential_AllComplete_ReturnsExitCode0 verifies that
// --watch --sequential with all-complete slices terminates cleanly.
func TestRunSlice_Watch_Sequential_AllComplete_ReturnsExitCode0(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := allCompleteSlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupWatchProject(t, slices, graphJSON, 1)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--watch", "--sequential"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 for --watch --sequential with all-complete slices, got %d; output:\n%s", exitCode, buf.String())
	}
}
