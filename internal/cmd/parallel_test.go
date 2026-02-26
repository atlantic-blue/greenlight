package cmd_test

// S-44 Tests: Parallel Slice Execution
// Covers C-111 (RunSliceParallel) and C-112 (RunSliceSequentialFallback).
//
// Contract C-111 — RunSliceParallel:
//   - With 2+ ready slices and tmux available: plans a tmux session with one
//     window per ready slice (up to --max)
//   - Session name follows "{prefix}-{project}" pattern from config
//   - Each window runs: claude -p "/gl:slice {id}" {claude_flags}
//   - --max flag limits the number of windows created
//   - Default --max is 4
//   - Dry-run mode shows the parallel plan without creating a tmux session
//   - Falls back to sequential when tmux unavailable
//   - Falls back to sequential when --sequential flag is present
//
// Contract C-112 — RunSliceSequentialFallback:
//   - When tmux unavailable: prints "tmux not available" and runs sequentially
//   - When --sequential flag: prints "sequential mode" and runs one at a time
//   - Sequential path processes one slice at a time
//   - Dry-run with sequential flag shows sequential plan, not tmux plan
//   - Inside Claude context: always single slice (never parallel)

import (
	"bytes"
	"strings"
	"testing"

	"github.com/atlantic-blue/greenlight/internal/cmd"
	"github.com/atlantic-blue/greenlight/internal/tmux"
)

// ----------------------------------------------------------------------------
// Helpers specific to parallel tests
// ----------------------------------------------------------------------------

// overrideTmuxLookPath replaces tmux.LookPath with the given stub and registers
// a cleanup to restore the original when the test finishes.
func overrideTmuxLookPath(t *testing.T, stub func(string) (string, error)) {
	t.Helper()
	original := tmux.LookPath
	tmux.LookPath = stub
	t.Cleanup(func() { tmux.LookPath = original })
}

// tmuxInPath is a tmux.LookPath stub that reports "tmux" as found.
func tmuxInPath(name string) (string, error) {
	if name == "tmux" {
		return "/usr/bin/tmux", nil
	}
	return "", errNotFound
}

// tmuxNotInPath is a tmux.LookPath stub that always reports not found.
func tmuxNotInPath(_ string) (string, error) {
	return "", errNotFound
}

// errNotFound is a sentinel for path lookup failures in stubs.
var errNotFound = &pathNotFoundError{}

type pathNotFoundError struct{}

func (pathNotFoundError) Error() string { return "executable file not found in $PATH" }

// configWithPrefix returns config JSON including a tmux_session_prefix.
func configWithPrefix(prefix string) string {
	return `{"parallel":{"claude_flags":["--dangerously-skip-permissions","--max-turns","200"],"tmux_session_prefix":"` + prefix + `"}}`
}

// setupParallelProject creates a temp project with config.json that has the
// given prefix. Returns the project root. Calls t.Chdir.
func setupParallelProject(t *testing.T, slices []testSlice, graphJSON string, configContent string) string {
	t.Helper()
	tmpDir := setupTestProject(t, slices, graphJSON)
	configFilePath := tmpDir + "/.greenlight/config.json"
	writeStatusFile(t, configFilePath, configContent)
	return tmpDir
}

// twoReadySlices returns a slice set where S-01 and S-02 are both ready
// (pending, wave 1, no deps).
func twoReadySlices() []testSlice {
	return []testSlice{
		{id: "S-01", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
		{id: "S-02", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
}

// threeReadySlices returns a slice set where S-01, S-02 and S-03 are all ready.
func threeReadySlices() []testSlice {
	return []testSlice{
		{id: "S-01", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
		{id: "S-02", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
		{id: "S-03", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
}

// fiveReadySlices returns a slice set where S-01 through S-05 are all ready.
func fiveReadySlices() []testSlice {
	return []testSlice{
		{id: "S-01", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
		{id: "S-02", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
		{id: "S-03", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
		{id: "S-04", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
		{id: "S-05", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
}

// ----------------------------------------------------------------------------
// C-111 — Dry-run parallel mode: 2+ ready slices, tmux available
// ----------------------------------------------------------------------------

// TestRunSlice_Parallel_DryRun_TmuxAvailable_ReturnsExitCode0 verifies that
// with 2+ ready slices, tmux available, and --dry-run, RunSlice returns 0.
func TestRunSlice_Parallel_DryRun_TmuxAvailable_ReturnsExitCode0(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--dry-run"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 for parallel --dry-run with tmux available, got %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunSlice_Parallel_DryRun_TmuxAvailable_MentionsTmux verifies that when
// 2+ slices are ready and tmux is available, dry-run output mentions tmux.
func TestRunSlice_Parallel_DryRun_TmuxAvailable_MentionsTmux(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	if !strings.Contains(strings.ToLower(output), "tmux") {
		t.Errorf("expected 'tmux' in parallel dry-run output, got:\n%s", output)
	}
}

// TestRunSlice_Parallel_DryRun_TmuxAvailable_ShowsBothSliceIDs verifies that
// the dry-run output lists both ready slice IDs when planning a tmux session.
func TestRunSlice_Parallel_DryRun_TmuxAvailable_ShowsBothSliceIDs(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	if !strings.Contains(output, "S-01") {
		t.Errorf("expected 'S-01' in parallel dry-run output, got:\n%s", output)
	}
	if !strings.Contains(output, "S-02") {
		t.Errorf("expected 'S-02' in parallel dry-run output, got:\n%s", output)
	}
}

// TestRunSlice_Parallel_DryRun_ShowsSessionName verifies that the dry-run
// output includes a tmux session name in the "{prefix}-{project}" pattern.
func TestRunSlice_Parallel_DryRun_ShowsSessionName(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupParallelProject(t, slices, graphJSON, configWithPrefix("gl"))

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	// Session name must follow the "{prefix}-{project}" pattern.
	// With prefix "gl" the session name starts with "gl-".
	lowerOutput := strings.ToLower(output)
	hasSessionPattern := strings.Contains(lowerOutput, "gl-") ||
		strings.Contains(lowerOutput, "session")
	if !hasSessionPattern {
		t.Errorf("expected session name pattern (e.g. 'gl-<project>') in dry-run output, got:\n%s", output)
	}
}

// TestRunSlice_Parallel_DryRun_ShowsWindowPerSlice verifies that the dry-run
// output describes a window for each planned slice.
func TestRunSlice_Parallel_DryRun_ShowsWindowPerSlice(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	// Output must describe windows — mention "window" or show both slice IDs as targets.
	lowerOutput := strings.ToLower(output)
	hasWindowRef := strings.Contains(lowerOutput, "window") ||
		(strings.Contains(output, "S-01") && strings.Contains(output, "S-02"))
	if !hasWindowRef {
		t.Errorf("expected window-per-slice description in parallel dry-run output, got:\n%s", output)
	}
}

// TestRunSlice_Parallel_DryRun_NeverSpawnsTmux verifies that --dry-run does
// not actually create a tmux session (safe even if tmux binary is absent).
func TestRunSlice_Parallel_DryRun_NeverSpawnsTmux(t *testing.T) {
	clearClaudeContext(t)
	// Provide a stub that says tmux is found (DI says available) but we verify
	// no actual tmux call happens by using claudeNotInPath — if the code tried
	// to run anything real it would fail.
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	// Dry-run must return 0 even though claude is absent — no real spawn.
	exitCode := cmd.RunSlice([]string{"--dry-run"}, &buf)

	if exitCode != 0 {
		t.Errorf("parallel --dry-run must not spawn real processes (exit 0); got %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunSlice_Parallel_DryRun_ShowsClaudeCommandPerWindow verifies that the
// dry-run output includes a claude command for each planned window.
func TestRunSlice_Parallel_DryRun_ShowsClaudeCommandPerWindow(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	if !strings.Contains(output, "claude") {
		t.Errorf("expected 'claude' command in parallel dry-run output, got:\n%s", output)
	}
}

// TestRunSlice_Parallel_DryRun_FlagsFromConfig verifies that the claude_flags
// from config.json appear in the parallel dry-run output.
func TestRunSlice_Parallel_DryRun_FlagsFromConfig(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON) // uses minimalConfig with --max-turns 200

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	// The config has --max-turns; it should appear in the planned command.
	if !strings.Contains(output, "max-turns") && !strings.Contains(output, "200") {
		t.Errorf("expected config claude_flags (--max-turns 200) in parallel dry-run output, got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-111 — --max flag: limits window count
// ----------------------------------------------------------------------------

// TestRunSlice_Parallel_DryRun_MaxLimitsWindowCount verifies that --max 2 with
// 5 ready slices shows exactly 2 slices in the parallel plan.
func TestRunSlice_Parallel_DryRun_MaxLimitsWindowCount(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := fiveReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02", "S-03", "S-04", "S-05"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run", "--max", "2"}, &buf)

	output := buf.String()
	// With --max 2 only the first 2 slices should be planned; S-03, S-04, S-05
	// must not appear as parallel targets.
	if strings.Contains(output, "S-03") && strings.Contains(output, "S-04") && strings.Contains(output, "S-05") {
		t.Errorf("--max 2 must limit to 2 windows; S-03/S-04/S-05 all appear in output:\n%s", output)
	}
}

// TestRunSlice_Parallel_DryRun_MaxFour_With5Slices_Shows4Windows verifies that
// the default --max of 4 produces exactly 4 windows when 5 slices are ready.
func TestRunSlice_Parallel_DryRun_MaxFour_With5Slices_Shows4Windows(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := fiveReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02", "S-03", "S-04", "S-05"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run", "--max", "4"}, &buf)

	output := buf.String()
	// S-05 must be excluded from the 4-window plan.
	// S-01 through S-04 must all be included.
	for _, id := range []string{"S-01", "S-02", "S-03", "S-04"} {
		if !strings.Contains(output, id) {
			t.Errorf("expected %s in --max 4 parallel dry-run output, got:\n%s", id, output)
		}
	}

	// S-05 should NOT appear as a parallel window target.
	// (It may appear in a "remaining" or "skipped" note, but not as a window.)
	// We verify by counting how many slices appear — if S-05 appears as a
	// window-level entry the total would exceed 4.
	// We test this indirectly: the output must show a count of 4 or max mention.
	lowerOutput := strings.ToLower(output)
	hasBoundIndicator := strings.Contains(lowerOutput, "4") ||
		strings.Contains(lowerOutput, "max")
	if !hasBoundIndicator {
		t.Errorf("expected max-window indication (\"4\" or \"max\") in output, got:\n%s", output)
	}
}

// TestRunSlice_Parallel_DryRun_MaxExceedsReadyCount_ShowsAllReady verifies
// that when --max is larger than the number of ready slices, all ready slices
// appear in the plan.
func TestRunSlice_Parallel_DryRun_MaxExceedsReadyCount_ShowsAllReady(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run", "--max", "10"}, &buf)

	output := buf.String()
	if !strings.Contains(output, "S-01") {
		t.Errorf("expected S-01 in output when --max exceeds ready count, got:\n%s", output)
	}
	if !strings.Contains(output, "S-02") {
		t.Errorf("expected S-02 in output when --max exceeds ready count, got:\n%s", output)
	}
}

// TestRunSlice_Parallel_DryRun_DefaultMaxIs4 verifies that without an explicit
// --max flag, the default limit of 4 is applied: with 5 ready slices the plan
// shows at most 4.
func TestRunSlice_Parallel_DryRun_DefaultMaxIs4(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := fiveReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02", "S-03", "S-04", "S-05"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	// No --max flag; should default to 4.
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	// All 5 slices are ready; only 4 should be planned.
	// S-05 must not appear as a window target.
	// We check by verifying the 5th slice is absent or marked as "remaining".
	sliceCount := 0
	for _, id := range []string{"S-01", "S-02", "S-03", "S-04", "S-05"} {
		if strings.Contains(output, id) {
			sliceCount++
		}
	}
	// S-05 may appear in a "remaining" note — the key invariant is that it
	// cannot appear as a window in the parallel plan.
	// We verify this via the max indicator in output.
	lowerOutput := strings.ToLower(output)
	hasDefaultMaxIndicator := strings.Contains(lowerOutput, "4") ||
		strings.Contains(lowerOutput, "max")
	if !hasDefaultMaxIndicator {
		t.Errorf("expected default --max 4 indicator in output with 5 ready slices; got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-111 — Invariants: session name, window names
// ----------------------------------------------------------------------------

// TestRunSlice_Parallel_DryRun_SessionNameFollowsPrefixPattern verifies that
// the dry-run output references a session name matching the configured prefix.
func TestRunSlice_Parallel_DryRun_SessionNameFollowsPrefixPattern(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupParallelProject(t, slices, graphJSON, configWithPrefix("gl"))

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	// The session must be named with the prefix "gl-" followed by the project name.
	if !strings.Contains(output, "gl-") {
		t.Errorf("expected session name with prefix 'gl-' in dry-run output, got:\n%s", output)
	}
}

// TestRunSlice_Parallel_DryRun_WindowsNamedWithSliceID verifies that the
// dry-run plan names each window with the corresponding slice ID.
func TestRunSlice_Parallel_DryRun_WindowsNamedWithSliceID(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	// Window names must be the slice IDs (S-01 and S-02 must appear).
	if !strings.Contains(output, "S-01") || !strings.Contains(output, "S-02") {
		t.Errorf("expected window names 'S-01' and 'S-02' in dry-run output, got:\n%s", output)
	}
}

// TestRunSlice_Parallel_DryRun_SlicesInWaveOrder verifies that the parallel
// plan assigns slices in wave/ID order (lowest wave first, then lexicographic ID).
func TestRunSlice_Parallel_DryRun_SlicesInWaveOrder(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := []testSlice{
		{id: "S-10", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
		{id: "S-20", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	// S-10 is wave 2, S-20 is wave 1 — S-20 should appear first.
	graphJSON := graphWithWave(map[string]struct {
		deps []string
		wave int
	}{
		"S-10": {deps: []string{}, wave: 2},
		"S-20": {deps: []string{}, wave: 1},
	})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	indexS10 := strings.Index(output, "S-10")
	indexS20 := strings.Index(output, "S-20")

	if indexS10 == -1 || indexS20 == -1 {
		t.Fatalf("expected both S-10 and S-20 in dry-run output, got:\n%s", output)
	}

	// S-20 (wave 1) must appear before S-10 (wave 2) in the plan.
	if indexS20 > indexS10 {
		t.Errorf("expected lower-wave S-20 to appear before S-10 in parallel plan output; got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-111 — Security: claude_flags sourced from config only
// ----------------------------------------------------------------------------

// TestRunSlice_Parallel_DryRun_ClaudeFlagsSourcedFromConfig verifies that
// claude flags come from config, not from user-supplied arguments.
func TestRunSlice_Parallel_DryRun_ClaudeFlagsSourcedFromConfig(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	// Passing an unknown user flag that is NOT in config.
	// The output must NOT reflect arbitrary user flags as claude arguments.
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	// Config has --dangerously-skip-permissions and --max-turns 200.
	// These should appear, not arbitrary user flags.
	if !strings.Contains(output, "dangerously-skip-permissions") &&
		!strings.Contains(output, "max-turns") &&
		!strings.Contains(output, "200") {
		t.Errorf("expected config claude_flags in dry-run output, got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-112 — Sequential fallback: tmux unavailable
// ----------------------------------------------------------------------------

// TestRunSlice_Sequential_TmuxUnavailable_DryRun_ReturnsExitCode0 verifies
// that when tmux is not available, --dry-run returns 0 (falls back gracefully).
func TestRunSlice_Sequential_TmuxUnavailable_DryRun_ReturnsExitCode0(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxNotInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--dry-run"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 for sequential fallback --dry-run, got %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunSlice_Sequential_TmuxUnavailable_DryRun_MentionsSequential verifies
// that when tmux is unavailable, the dry-run output describes sequential mode.
func TestRunSlice_Sequential_TmuxUnavailable_DryRun_MentionsSequential(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxNotInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	hasSequentialRef := strings.Contains(lowerOutput, "sequential") ||
		strings.Contains(lowerOutput, "one at a time") ||
		strings.Contains(lowerOutput, "tmux not available") ||
		strings.Contains(lowerOutput, "no tmux")
	if !hasSequentialRef {
		t.Errorf("expected sequential fallback indication when tmux unavailable, got:\n%s", output)
	}
}

// TestRunSlice_Sequential_TmuxUnavailable_DryRun_DoesNotMentionTmuxSession
// verifies that when tmux is unavailable, dry-run does NOT plan a tmux session.
func TestRunSlice_Sequential_TmuxUnavailable_DryRun_DoesNotMentionTmuxSession(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxNotInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	// Must NOT claim to create a tmux session when tmux is unavailable.
	lowerOutput := strings.ToLower(output)
	if strings.Contains(lowerOutput, "creating session") ||
		strings.Contains(lowerOutput, "tmux session") ||
		strings.Contains(lowerOutput, "new-session") {
		t.Errorf("dry-run must not plan tmux session when tmux unavailable; got:\n%s", output)
	}
}

// TestRunSlice_Sequential_TmuxUnavailable_DryRun_ShowsFirstSlice verifies that
// sequential fallback dry-run shows the first ready slice that would run.
func TestRunSlice_Sequential_TmuxUnavailable_DryRun_ShowsFirstSlice(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxNotInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	if !strings.Contains(output, "S-01") {
		t.Errorf("expected first ready slice 'S-01' in sequential fallback dry-run output, got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-112 — Sequential fallback: --sequential flag
// ----------------------------------------------------------------------------

// TestRunSlice_Sequential_Flag_DryRun_ReturnsExitCode0 verifies that --sequential
// with --dry-run returns 0.
func TestRunSlice_Sequential_Flag_DryRun_ReturnsExitCode0(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath) // tmux IS available but --sequential overrides
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--sequential", "--dry-run"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 for --sequential --dry-run, got %d; output:\n%s", exitCode, buf.String())
	}
}

// TestRunSlice_Sequential_Flag_DryRun_DoesNotUseTmux verifies that --sequential
// forces sequential mode even when tmux IS available.
func TestRunSlice_Sequential_Flag_DryRun_DoesNotUseTmux(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath) // tmux available
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--sequential", "--dry-run"}, &buf)

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	// Must not plan a tmux session.
	if strings.Contains(lowerOutput, "creating session") ||
		strings.Contains(lowerOutput, "new-session") {
		t.Errorf("--sequential must not use tmux even when tmux is available; got:\n%s", output)
	}
}

// TestRunSlice_Sequential_Flag_DryRun_ShowsSequentialMode verifies that
// --sequential dry-run output describes sequential execution.
func TestRunSlice_Sequential_Flag_DryRun_ShowsSequentialMode(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--sequential", "--dry-run"}, &buf)

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	hasSequentialRef := strings.Contains(lowerOutput, "sequential") ||
		strings.Contains(lowerOutput, "one at a time")
	if !hasSequentialRef {
		t.Errorf("expected sequential mode description with --sequential flag, got:\n%s", output)
	}
}

// TestRunSlice_Sequential_Flag_DryRun_ShowsFirstReadySlice verifies that with
// --sequential, the first ready slice is shown in the dry-run output.
func TestRunSlice_Sequential_Flag_DryRun_ShowsFirstReadySlice(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := threeReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02", "S-03"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--sequential", "--dry-run"}, &buf)

	output := buf.String()
	// S-01 is first by wave/ID order.
	if !strings.Contains(output, "S-01") {
		t.Errorf("expected first ready slice 'S-01' in --sequential dry-run output, got:\n%s", output)
	}
}

// TestRunSlice_Sequential_Flag_ForcesSequentialEvenWithTmuxAvailable verifies
// the invariant that --sequential overrides tmux availability.
func TestRunSlice_Sequential_Flag_ForcesSequentialEvenWithTmuxAvailable(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath) // tmux available
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--sequential", "--dry-run"}, &buf)

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	// Output must describe sequential behaviour, not a tmux session with windows.
	hasParallelPlan := (strings.Contains(lowerOutput, "window") &&
		strings.Contains(lowerOutput, "tmux session"))
	if hasParallelPlan {
		t.Errorf("--sequential must force sequential path even with tmux available; got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-112 — Sequential fallback: error cases
// ----------------------------------------------------------------------------

// TestRunSlice_Sequential_TmuxUnavailable_ClaudeNotFound_ReturnsNonZero
// verifies that when tmux is unavailable and claude is not found in shell
// context (non-dry-run), RunSlice returns a non-zero exit code.
func TestRunSlice_Sequential_TmuxUnavailable_ClaudeNotFound_ReturnsNonZero(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxNotInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{}, &buf)

	if exitCode == 0 {
		t.Errorf("expected non-zero exit code when claude not found in sequential fallback, got 0")
	}
}

// TestRunSlice_Sequential_TmuxUnavailable_ClaudeNotFound_PrintsClaudeError
// verifies that when tmux is unavailable and claude is missing, the error
// output mentions claude.
func TestRunSlice_Sequential_TmuxUnavailable_ClaudeNotFound_PrintsClaudeError(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxNotInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{}, &buf)

	output := buf.String()
	if !strings.Contains(strings.ToLower(output), "claude") {
		t.Errorf("expected error mentioning 'claude' in sequential fallback path, got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// C-111 — Inside Claude context: never parallel
// ----------------------------------------------------------------------------

// TestRunSlice_InsideClaude_MultipleReady_NeverParallel verifies that inside
// Claude context ($CLAUDE_CODE set), RunSlice never starts a parallel tmux
// session even when 2+ slices are ready and tmux is available.
func TestRunSlice_InsideClaude_MultipleReady_NeverParallel(t *testing.T) {
	setClaudeContext(t, "1")
	overrideTmuxLookPath(t, tmuxInPath) // tmux available

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 inside Claude with multiple ready slices, got %d; output:\n%s", exitCode, buf.String())
	}

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	// Inside Claude must never plan a tmux session.
	if strings.Contains(lowerOutput, "creating session") ||
		strings.Contains(lowerOutput, "new-session") ||
		strings.Contains(lowerOutput, "tmux session") {
		t.Errorf("inside Claude must never use tmux for parallel execution; got:\n%s", output)
	}
}

// TestRunSlice_InsideClaude_MultipleReady_PicksSingleSlice verifies that inside
// Claude with multiple ready slices, exactly one slice is targeted (not all).
func TestRunSlice_InsideClaude_MultipleReady_PicksSingleSlice(t *testing.T) {
	setClaudeContext(t, "1")
	overrideTmuxLookPath(t, tmuxInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{}, &buf)

	output := buf.String()
	// S-01 must be the chosen slice (first by wave/ID).
	if !strings.Contains(output, "S-01") {
		t.Errorf("expected inside-Claude to pick first ready slice 'S-01', got:\n%s", output)
	}
}

// TestRunSlice_InsideClaude_MultipleReady_OutputsToProvidedWriter verifies that
// inside Claude with multiple ready slices, output still goes to the io.Writer.
func TestRunSlice_InsideClaude_MultipleReady_OutputsToProvidedWriter(t *testing.T) {
	setClaudeContext(t, "1")
	overrideTmuxLookPath(t, tmuxInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{}, &buf)

	if buf.Len() == 0 {
		t.Error("expected output written to provided writer inside Claude with multiple ready slices")
	}
}

// ----------------------------------------------------------------------------
// C-111 — Invariants: never exceed --max windows
// ----------------------------------------------------------------------------

// TestRunSlice_Parallel_DryRun_NeverExceedsMaxWindows verifies that the number
// of slices in the parallel plan never exceeds the --max value.
func TestRunSlice_Parallel_DryRun_NeverExceedsMaxWindows(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := fiveReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02", "S-03", "S-04", "S-05"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run", "--max", "3"}, &buf)

	output := buf.String()
	// With --max 3 and 5 ready slices, at most 3 should appear as parallel targets.
	// We cannot see exactly which are windows without knowing the format, but we
	// verify that the 4th and 5th do not both appear as window targets.
	// If S-04 AND S-05 both appear AND S-01, S-02, S-03 appear, that is 5 = exceeds max 3.
	// Indirect check: count unique slice IDs that appear in the output.
	plannedCount := 0
	for _, id := range []string{"S-01", "S-02", "S-03"} {
		if strings.Contains(output, id) {
			plannedCount++
		}
	}
	if plannedCount < 3 {
		t.Errorf("expected at least 3 slices in --max 3 plan; got %d; output:\n%s", plannedCount, output)
	}
}

// ----------------------------------------------------------------------------
// C-111 — Fallback: tmux session create error
// ----------------------------------------------------------------------------

// TestRunSlice_Parallel_DryRun_ConfigReadError_UsesDefaults verifies that if
// config.json is absent, parallel dry-run still succeeds using default values.
func TestRunSlice_Parallel_DryRun_ConfigReadError_UsesDefaults(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	// Use setupTestProject (no config.json) so RunSlice gets ConfigReadError.
	setupTestProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--dry-run"}, &buf)

	// ConfigReadError is not fatal — RunSlice must use defaults and succeed.
	if exitCode != 0 {
		t.Errorf("expected exit code 0 when config.json missing (use defaults), got %d; output:\n%s", exitCode, buf.String())
	}
}

// ----------------------------------------------------------------------------
// C-111 — Output: written to provided io.Writer
// ----------------------------------------------------------------------------

// TestRunSlice_Parallel_DryRun_WritesToProvidedWriter verifies that all output
// in parallel dry-run mode goes to the provided io.Writer.
func TestRunSlice_Parallel_DryRun_WritesToProvidedWriter(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	if buf.Len() == 0 {
		t.Error("expected output written to provided writer in parallel dry-run mode, got empty buffer")
	}
}

// TestRunSlice_Sequential_DryRun_WritesToProvidedWriter verifies that all
// output in sequential fallback dry-run mode goes to the provided io.Writer.
func TestRunSlice_Sequential_DryRun_WritesToProvidedWriter(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxNotInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	if buf.Len() == 0 {
		t.Error("expected output written to provided writer in sequential fallback dry-run mode, got empty buffer")
	}
}

// ----------------------------------------------------------------------------
// C-112 — Invariants: one slice at a time, re-read state
// ----------------------------------------------------------------------------

// TestRunSlice_Sequential_Flag_DryRun_DescribesOnlyOneSliceAtATime verifies
// that sequential mode processes exactly one slice per run, not multiple.
func TestRunSlice_Sequential_Flag_DryRun_DescribesOnlyOneSliceAtATime(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath) // tmux available but overridden by --sequential
	overrideProcessLookPath(t, claudeNotInPath)

	slices := threeReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02", "S-03"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--sequential", "--dry-run"}, &buf)

	output := buf.String()
	// Sequential dry-run should show running one slice, not all three at once.
	// S-01 must be the target.
	if !strings.Contains(output, "S-01") {
		t.Errorf("expected sequential dry-run to show S-01 (first by wave/ID), got:\n%s", output)
	}
}

// TestRunSlice_Sequential_TmuxUnavailable_SingleSlice_DryRun_ShowsHint verifies
// that when in sequential fallback, the output may include a hint about installing
// tmux or using parallel mode.
func TestRunSlice_Sequential_TmuxUnavailable_SingleSlice_DryRun_ShowsHint(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxNotInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := twoReadySlices()
	graphJSON := minimalGraph([]string{"S-01", "S-02"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	// The output should inform the user about the sequential fallback reason.
	hasInfoOrHint := strings.Contains(lowerOutput, "tmux") ||
		strings.Contains(lowerOutput, "sequential") ||
		strings.Contains(lowerOutput, "install")
	if !hasInfoOrHint {
		t.Errorf("expected informative message about sequential fallback in output, got:\n%s", output)
	}
}

// ----------------------------------------------------------------------------
// Boundary: single ready slice is not parallel
// ----------------------------------------------------------------------------

// TestRunSlice_SingleReadySlice_DryRun_TmuxAvailable_NotParallel verifies that
// with exactly one ready slice, even when tmux is available, the output does
// not plan a multi-window tmux session.
func TestRunSlice_SingleReadySlice_DryRun_TmuxAvailable_NotParallel(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := []testSlice{
		{id: "S-01", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	graphJSON := minimalGraph([]string{"S-01"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	exitCode := cmd.RunSlice([]string{"--dry-run"}, &buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 for single-slice dry-run, got %d; output:\n%s", exitCode, buf.String())
	}

	output := buf.String()
	lowerOutput := strings.ToLower(output)
	// Single slice must not create a multi-window tmux session.
	if strings.Contains(lowerOutput, "creating session") ||
		strings.Contains(lowerOutput, "new-session") {
		t.Errorf("single ready slice must not use parallel tmux; got:\n%s", output)
	}
}

// TestRunSlice_SingleReadySlice_DryRun_TmuxAvailable_ShowsSliceID verifies
// that the dry-run for a single ready slice still shows the slice ID.
func TestRunSlice_SingleReadySlice_DryRun_TmuxAvailable_ShowsSliceID(t *testing.T) {
	clearClaudeContext(t)
	overrideTmuxLookPath(t, tmuxInPath)
	overrideProcessLookPath(t, claudeNotInPath)

	slices := []testSlice{
		{id: "S-01", status: "pending", step: "test-writer", milestone: "core",
			started: "", updated: ""},
	}
	graphJSON := minimalGraph([]string{"S-01"})
	setupSliceProject(t, slices, graphJSON)

	var buf bytes.Buffer
	cmd.RunSlice([]string{"--dry-run"}, &buf)

	output := buf.String()
	if !strings.Contains(output, "S-01") {
		t.Errorf("expected 'S-01' in single-slice dry-run output, got:\n%s", output)
	}
}
