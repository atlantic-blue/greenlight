package cmd

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/atlantic-blue/greenlight/internal/state"
)

const (
	progressBarWidth = 20
	slicesDir        = ".greenlight/slices"
	graphPath        = ".greenlight/GRAPH.json"
)

// RunStatus handles the "status" subcommand.
//
// Default mode: prints a multi-line formatted report with Progress, Running,
// Ready, Blocked, and Tests sections.
//
// Compact mode (--compact flag): prints a single line suitable for tmux status
// bars: "{N}/{M} done | {K} running". On error it prints "? slices | ? running"
// and returns 0 (tmux bar resilience).
func RunStatus(args []string, stdout io.Writer) int {
	compact := parseCompactFlag(args)

	slices, graph, loadError := loadProjectData()
	if loadError != nil {
		return handleLoadError(loadError, compact, stdout)
	}

	stats := computeStats(slices)

	if compact {
		writeCompactOutput(stats, stdout)
		return 0
	}

	writeFullOutput(slices, graph, stats, stdout)
	return 0
}

// parseCompactFlag checks whether --compact is present in args.
func parseCompactFlag(args []string) bool {
	for _, arg := range args {
		if arg == "--compact" {
			return true
		}
	}
	return false
}

// projectData bundles the loaded slices and optional graph.
type projectData struct {
	slices []state.SliceInfo
	graph  *state.Graph
}

// loadProjectData reads slices from disk and optionally the graph.
// Returns an error if the slices directory or slice files are missing.
// A missing GRAPH.json is not fatal — graph will be nil in that case.
func loadProjectData() ([]state.SliceInfo, *state.Graph, error) {
	slices, readError := state.ReadSlices(slicesDir)
	if readError != nil {
		return nil, nil, readError
	}

	graph, graphError := state.ReadGraph(graphPath)
	if graphError != nil {
		// Missing graph is degraded mode, not fatal.
		return slices, nil, nil
	}

	return slices, graph, nil
}

// sliceStats holds aggregated metrics over all loaded slices.
type sliceStats struct {
	total         int
	complete      int
	inProgress    int
	pending       int
	totalTests    int
	totalSecTests int
	running       []state.SliceInfo
}

// computeStats aggregates counts from all slices.
func computeStats(slices []state.SliceInfo) sliceStats {
	stats := sliceStats{total: len(slices)}
	for _, slice := range slices {
		stats.totalTests += slice.Tests
		stats.totalSecTests += slice.SecurityTests
		switch slice.Status {
		case "complete":
			stats.complete++
		case "in_progress":
			stats.inProgress++
			stats.running = append(stats.running, slice)
		case "pending":
			stats.pending++
		}
	}
	return stats
}

// handleLoadError writes the appropriate error message and returns the exit code.
func handleLoadError(loadError error, compact bool, stdout io.Writer) int {
	if compact {
		fmt.Fprintln(stdout, "? slices | ? running")
		return 0
	}

	if errors.Is(loadError, state.ErrDirNotFound) {
		fmt.Fprintln(stdout, "error: .greenlight/ directory not found. Run 'greenlight init' to set up this project.")
		return 1
	}

	if errors.Is(loadError, state.ErrNoSliceFiles) {
		fmt.Fprintln(stdout, "error: no slice files found. Run 'greenlight init' to populate the slices directory.")
		return 1
	}

	fmt.Fprintf(stdout, "error: %v\n", loadError)
	return 1
}

// writeCompactOutput writes the single-line compact format to stdout.
func writeCompactOutput(stats sliceStats, stdout io.Writer) {
	fmt.Fprintf(stdout, "%d/%d done | %d running\n", stats.complete, stats.total, stats.inProgress)
}

// writeFullOutput writes the multi-line status report.
func writeFullOutput(slices []state.SliceInfo, graph *state.Graph, stats sliceStats, stdout io.Writer) {
	writeProgressSection(stats, stdout)
	writeRunningSection(stats.running, stdout)
	writeReadyAndBlockedSections(slices, graph, stdout)
	writeTestsSection(stats, stdout)

	if graph == nil {
		fmt.Fprintln(stdout, "warn: GRAPH.json missing — dependency info unavailable. Run 'greenlight init' to generate it.")
	}
}

// writeProgressSection prints the ASCII progress bar line.
func writeProgressSection(stats sliceStats, stdout io.Writer) {
	bar := buildProgressBar(stats.complete, stats.total)
	fmt.Fprintf(stdout, "Progress: %s %d/%d\n", bar, stats.complete, stats.total)
}

// buildProgressBar returns an ASCII progress bar string like [####....].
func buildProgressBar(complete, total int) string {
	if total == 0 {
		return "[" + strings.Repeat(".", progressBarWidth) + "]"
	}

	filled := (complete * progressBarWidth) / total
	if filled > progressBarWidth {
		filled = progressBarWidth
	}

	return "[" + strings.Repeat("#", filled) + strings.Repeat(".", progressBarWidth-filled) + "]"
}

// writeRunningSection prints all in_progress slices with their current step.
func writeRunningSection(running []state.SliceInfo, stdout io.Writer) {
	if len(running) == 0 {
		fmt.Fprintln(stdout, "Running:  —")
		return
	}

	parts := make([]string, 0, len(running))
	for _, slice := range running {
		parts = append(parts, slice.ID+" ("+slice.Step+")")
	}
	fmt.Fprintln(stdout, "Running:  "+strings.Join(parts, ", "))
}

// writeReadyAndBlockedSections computes ready vs blocked slices and prints both.
func writeReadyAndBlockedSections(slices []state.SliceInfo, graph *state.Graph, stdout io.Writer) {
	if graph == nil {
		fmt.Fprintln(stdout, "Ready:    (dependency info unavailable)")
		fmt.Fprintln(stdout, "Blocked:  (dependency info unavailable)")
		return
	}

	readySlices := state.FindReadySlices(slices, graph)
	blockedSlices := findBlockedSlices(slices, graph, readySlices)

	writeReadySection(readySlices, stdout)
	writeBlockedSection(blockedSlices, graph, slices, stdout)
}

// writeReadySection prints pending slices that are unblocked.
func writeReadySection(ready []state.SliceInfo, stdout io.Writer) {
	if len(ready) == 0 {
		fmt.Fprintln(stdout, "Ready:    —")
		return
	}

	ids := make([]string, 0, len(ready))
	for _, slice := range ready {
		ids = append(ids, slice.ID)
	}
	fmt.Fprintln(stdout, "Ready:    "+strings.Join(ids, ", "))
}

// writeBlockedSection prints pending slices that have unmet dependencies.
func writeBlockedSection(blocked []state.SliceInfo, graph *state.Graph, allSlices []state.SliceInfo, stdout io.Writer) {
	if len(blocked) == 0 {
		fmt.Fprintln(stdout, "Blocked:  —")
		return
	}

	statusByID := buildStatusMapFromSlices(allSlices)
	parts := make([]string, 0, len(blocked))
	for _, slice := range blocked {
		unmet := findUnmetDeps(slice.ID, graph, statusByID)
		parts = append(parts, slice.ID+" (needs "+strings.Join(unmet, ", ")+")")
	}
	fmt.Fprintln(stdout, "Blocked:  "+strings.Join(parts, ", "))
}

// findBlockedSlices returns pending slices that are NOT in the ready set.
func findBlockedSlices(slices []state.SliceInfo, graph *state.Graph, ready []state.SliceInfo) []state.SliceInfo {
	readySet := make(map[string]bool, len(ready))
	for _, slice := range ready {
		readySet[slice.ID] = true
	}

	var blocked []state.SliceInfo
	for _, slice := range slices {
		if slice.Status == "pending" && !readySet[slice.ID] {
			blocked = append(blocked, slice)
		}
	}
	return blocked
}

// findUnmetDeps returns the dependency IDs that are not yet complete.
func findUnmetDeps(sliceID string, graph *state.Graph, statusByID map[string]string) []string {
	graphSlice, exists := graph.Slices[sliceID]
	if !exists {
		return nil
	}

	var unmet []string
	for _, depID := range graphSlice.DependsOn {
		if statusByID[depID] != "complete" {
			unmet = append(unmet, depID)
		}
	}
	return unmet
}

// buildStatusMapFromSlices creates a map of slice ID to status.
func buildStatusMapFromSlices(slices []state.SliceInfo) map[string]string {
	statusByID := make(map[string]string, len(slices))
	for _, slice := range slices {
		statusByID[slice.ID] = slice.Status
	}
	return statusByID
}

// writeTestsSection prints the summed test counts.
func writeTestsSection(stats sliceStats, stdout io.Writer) {
	fmt.Fprintf(stdout, "Tests:    %d total (%d security)\n", stats.totalTests, stats.totalSecTests)
}
