package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/atlantic-blue/greenlight/internal/process"
	"github.com/atlantic-blue/greenlight/internal/state"
	"github.com/atlantic-blue/greenlight/internal/tmux"
)

const configPath = ".greenlight/config.json"

const defaultMaxWindows = 4

const defaultTmuxSessionPrefix = "gl"

const defaultWatchInterval = 30

// sliceConfig holds the parsed structure of .greenlight/config.json.
type sliceConfig struct {
	Parallel struct {
		ClaudeFlags           []string `json:"claude_flags"`
		TmuxSessionPrefix     string   `json:"tmux_session_prefix"`
		WatchIntervalSeconds  int      `json:"watch_interval_seconds"`
	} `json:"parallel"`
}

// sliceFlags holds the parsed command-line flags for the slice subcommand.
type sliceFlags struct {
	dryRun     bool
	sequential bool
	watch      bool
	max        int
}

// RunSlice handles the "slice" subcommand.
//
// With a slice ID argument: validates it against GRAPH.json, then runs it
// headlessly via process.SpawnClaude (shell context) or prints info (inside Claude).
//
// Without a slice ID: auto-detects ready slices via state.FindReadySlices and
// runs in parallel (tmux), sequential fallback, or inside-Claude single mode.
func RunSlice(args []string, stdout io.Writer) int {
	if _, statError := os.Stat(greenlightDir); os.IsNotExist(statError) {
		fmt.Fprintln(stdout, "error: not a greenlight project. Run 'gl init' to set up this project.")
		return 1
	}

	flags, remaining := parseSliceFlags(args)
	sliceID := extractSliceID(remaining)

	config := loadSliceConfig()
	executionContext := state.DetectContext()

	if sliceID == "" {
		return runAutoDetect(flags, config, executionContext, stdout)
	}

	return runNamedSlice(sliceID, flags, config, executionContext, stdout)
}

// parseSliceFlags extracts known flags from args and returns the flags struct
// plus the remaining positional arguments.
func parseSliceFlags(args []string) (sliceFlags, []string) {
	flags := sliceFlags{
		max: defaultMaxWindows,
	}
	var remaining []string

	for index := 0; index < len(args); index++ {
		arg := args[index]
		switch {
		case arg == "--dry-run":
			flags.dryRun = true
		case arg == "--sequential":
			flags.sequential = true
		case arg == "--watch":
			flags.watch = true
		case arg == "--max":
			if index+1 < len(args) {
				index++
				if value, convertError := strconv.Atoi(args[index]); convertError == nil {
					flags.max = value
				}
			}
		case strings.HasPrefix(arg, "--max="):
			valueStr := strings.TrimPrefix(arg, "--max=")
			if value, convertError := strconv.Atoi(valueStr); convertError == nil {
				flags.max = value
			}
		default:
			remaining = append(remaining, arg)
		}
	}

	return flags, remaining
}

// extractSliceID returns the first non-flag positional argument as the slice ID.
func extractSliceID(args []string) string {
	for _, arg := range args {
		if !strings.HasPrefix(arg, "--") {
			return arg
		}
	}
	return ""
}

// loadSliceConfig reads .greenlight/config.json and returns parsed config.
// Returns an empty config on any error (graceful degradation).
func loadSliceConfig() sliceConfig {
	var config sliceConfig
	data, readError := os.ReadFile(configPath)
	if readError != nil {
		return config
	}

	if unmarshalError := json.Unmarshal(data, &config); unmarshalError != nil {
		return config
	}

	return config
}

// tmuxSessionPrefix returns the configured prefix or the default "gl".
func tmuxSessionPrefix(config sliceConfig) string {
	if config.Parallel.TmuxSessionPrefix != "" {
		return config.Parallel.TmuxSessionPrefix
	}
	return defaultTmuxSessionPrefix
}

// watchIntervalSeconds returns the configured interval or the default.
func watchIntervalSeconds(config sliceConfig) int {
	if config.Parallel.WatchIntervalSeconds > 0 {
		return config.Parallel.WatchIntervalSeconds
	}
	return defaultWatchInterval
}

// projectName returns the base name of the current working directory.
func projectName() string {
	cwd, cwdError := os.Getwd()
	if cwdError != nil {
		return "project"
	}
	return filepath.Base(cwd)
}

// runNamedSlice runs a specific slice by ID.
// Validates the ID against GRAPH.json, then dispatches to the appropriate handler.
func runNamedSlice(
	sliceID string,
	flags sliceFlags,
	config sliceConfig,
	executionContext state.ExecutionContext,
	stdout io.Writer,
) int {
	graph, graphError := state.ReadGraph(graphPath)
	if graphError != nil {
		if errors.Is(graphError, state.ErrFileNotFound) {
			fmt.Fprintf(stdout, "error: unknown slice ID %q. Run 'gl status' to see available slices.\n", sliceID)
			return 1
		}
		fmt.Fprintf(stdout, "error: could not read GRAPH.json: %v\n", graphError)
		return 1
	}

	if _, exists := graph.Slices[sliceID]; !exists {
		fmt.Fprintf(stdout, "error: unknown slice ID %q. Run 'gl status' to see available slices.\n", sliceID)
		return 1
	}

	if flags.dryRun {
		printDryRun(sliceID, config, stdout)
		return 0
	}

	if executionContext.InsideClaude {
		printInsideClaudeInfo(sliceID, graph, stdout)
		return 0
	}

	return spawnSlice(sliceID, config, stdout)
}

// runAutoDetect handles the no-slice-ID case: finds ready slices and dispatches
// to parallel, sequential, or single-slice mode depending on context.
func runAutoDetect(
	flags sliceFlags,
	config sliceConfig,
	executionContext state.ExecutionContext,
	stdout io.Writer,
) int {
	slices, graph, loadError := loadProjectData()
	if loadError != nil {
		fmt.Fprintf(stdout, "error: could not read project state: %v\n", loadError)
		return 1
	}

	if graph == nil {
		fmt.Fprintln(stdout, "error: GRAPH.json missing — cannot auto-detect ready slices.")
		return 1
	}

	// Enhanced dry-run (C-114): categorised output always takes precedence over
	// normal dispatch. This must be checked before the "0 ready" early return.
	if flags.dryRun {
		printEnhancedDryRun(slices, graph, flags, config, executionContext, stdout)
		return 0
	}

	// Watch mode (C-113): check for immediate termination before any dispatch.
	if flags.watch {
		return runWatch(slices, graph, flags, config, executionContext, stdout)
	}

	readySlices := state.FindReadySlices(slices, graph)

	if len(readySlices) == 0 {
		fmt.Fprintln(stdout, "No ready slices: all pending slices are blocked waiting on dependencies.")
		return 0
	}

	// Inside Claude: always single slice, never parallel.
	if executionContext.InsideClaude {
		targetSlice := readySlices[0]
		if len(readySlices) > 1 {
			remaining := len(readySlices) - 1
			fmt.Fprintf(stdout, "hint: %d more ready slice(s) available. Use --max or parallel mode (S-44) to run concurrently.\n", remaining)
		}
		return runNamedSlice(targetSlice.ID, flags, config, executionContext, stdout)
	}

	// --sequential flag forces sequential mode regardless of tmux availability.
	if flags.sequential {
		return runSequential(readySlices, flags, config, executionContext, stdout)
	}

	// tmux unavailable: fall back to sequential.
	if !tmux.IsAvailable() {
		return runSequential(readySlices, flags, config, executionContext, stdout)
	}

	// 2+ ready slices and tmux available: parallel mode.
	if len(readySlices) >= 2 {
		return runParallel(readySlices, flags, config, stdout)
	}

	// Single ready slice: run directly (no parallel tmux session needed).
	targetSlice := readySlices[0]
	return runNamedSlice(targetSlice.ID, flags, config, executionContext, stdout)
}

// runWatch implements C-113: watch mode loop.
// On initial check, if no running and no ready slices, terminates immediately.
// Inside Claude context, delegates to normal auto-detect (no poll loop).
func runWatch(
	slices []state.SliceInfo,
	graph *state.Graph,
	flags sliceFlags,
	config sliceConfig,
	executionContext state.ExecutionContext,
	stdout io.Writer,
) int {
	// Inside Claude: delegate to normal single-slice behaviour without poll loop.
	if executionContext.InsideClaude {
		readySlices := state.FindReadySlices(slices, graph)
		if len(readySlices) == 0 {
			printWatchSummary(slices, stdout)
			return 0
		}
		targetSlice := readySlices[0]
		if len(readySlices) > 1 {
			remaining := len(readySlices) - 1
			fmt.Fprintf(stdout, "hint: %d more ready slice(s) available. Use --max or parallel mode (S-44) to run concurrently.\n", remaining)
		}
		return runNamedSlice(targetSlice.ID, flags, config, executionContext, stdout)
	}

	readySlices := state.FindReadySlices(slices, graph)
	runningSlices := filterRunning(slices)

	// Immediate termination: no work in progress and nothing ready to launch.
	if len(readySlices) == 0 && len(runningSlices) == 0 {
		printWatchSummary(slices, stdout)
		return 0
	}

	// Initial launch: start up to --max ready slices.
	planned := readySlices
	if len(planned) > flags.max {
		planned = planned[:flags.max]
	}

	if flags.sequential || !tmux.IsAvailable() {
		if len(planned) > 0 {
			return spawnSlice(planned[0].ID, config, stdout)
		}
		printWatchSummary(slices, stdout)
		return 0
	}

	if len(planned) >= 2 {
		return executeTmuxSession(planned, flags.max, config, stdout)
	}

	if len(planned) == 1 {
		return spawnSlice(planned[0].ID, config, stdout)
	}

	printWatchSummary(slices, stdout)
	return 0
}

// printWatchSummary prints a termination summary for watch mode.
// When pending slices are blocked (0 ready, 0 running), it also notes the
// blocked state so callers can diagnose dependency deadlocks.
func printWatchSummary(slices []state.SliceInfo, stdout io.Writer) {
	total := len(slices)
	doneCount := 0
	pendingCount := 0
	totalTests := 0
	for _, slice := range slices {
		switch slice.Status {
		case "complete":
			doneCount++
			totalTests += slice.Tests
		case "pending":
			pendingCount++
		}
	}
	fmt.Fprintf(stdout, "Watch complete: %d/%d slices done, %d tests\n", doneCount, total, totalTests)
	if pendingCount > 0 {
		fmt.Fprintf(stdout, "No ready slices: %d pending slice(s) are blocked waiting on dependencies.\n", pendingCount)
	}
}

// filterRunning returns slices with status "in_progress".
func filterRunning(slices []state.SliceInfo) []state.SliceInfo {
	var running []state.SliceInfo
	for _, slice := range slices {
		if slice.Status == "in_progress" {
			running = append(running, slice)
		}
	}
	return running
}

// printEnhancedDryRun implements C-114: categorised dry-run output.
// Shows Ready/Running/Blocked/Would-launch categories and includes the
// parallel or sequential plan detail so existing tests continue to pass.
func printEnhancedDryRun(
	slices []state.SliceInfo,
	graph *state.Graph,
	flags sliceFlags,
	config sliceConfig,
	executionContext state.ExecutionContext,
	stdout io.Writer,
) {
	readySlices := state.FindReadySlices(slices, graph)
	runningSlices := filterRunning(slices)
	blockedSlices := findBlockedSlices(slices, graph, readySlices)

	// Limit ready display to --max (to keep "Would launch" and "Ready" aligned).
	displayReady := readySlices
	if len(displayReady) > flags.max {
		displayReady = displayReady[:flags.max]
	}

	// Print Ready category.
	if len(displayReady) == 0 {
		fmt.Fprintln(stdout, "Ready (0):    —")
	} else {
		ids := make([]string, 0, len(displayReady))
		for _, slice := range displayReady {
			ids = append(ids, slice.ID)
		}
		fmt.Fprintf(stdout, "Ready (%d):    %s\n", len(displayReady), strings.Join(ids, ", "))
	}

	// Print Running category.
	if len(runningSlices) == 0 {
		fmt.Fprintln(stdout, "Running (0):  —")
	} else {
		parts := make([]string, 0, len(runningSlices))
		for _, slice := range runningSlices {
			parts = append(parts, slice.ID+" ("+slice.Step+")")
		}
		fmt.Fprintf(stdout, "Running (%d):  %s\n", len(runningSlices), strings.Join(parts, ", "))
	}

	// Print Blocked category.
	statusByID := buildStatusMapFromSlices(slices)
	if len(blockedSlices) == 0 {
		fmt.Fprintln(stdout, "Blocked (0):  —")
	} else {
		parts := make([]string, 0, len(blockedSlices))
		for _, slice := range blockedSlices {
			unmet := findUnmetDeps(slice.ID, graph, statusByID)
			if len(unmet) > 0 {
				parts = append(parts, slice.ID+" (needs "+strings.Join(unmet, ", ")+")")
			} else {
				parts = append(parts, slice.ID)
			}
		}
		fmt.Fprintf(stdout, "Blocked (%d):  %s\n", len(blockedSlices), strings.Join(parts, ", "))
	}

	// Print Would launch section.
	if len(displayReady) == 0 {
		fmt.Fprintln(stdout, "Would launch: none")
		return
	}

	// Determine launch mode and print details.
	if !executionContext.InsideClaude && !flags.sequential && tmux.IsAvailable() && len(displayReady) >= 2 {
		printEnhancedDryRunParallel(displayReady, flags.max, config, stdout)
	} else {
		printEnhancedDryRunSequential(displayReady, flags, stdout)
	}
}

// printEnhancedDryRunParallel prints the would-launch section for parallel mode.
func printEnhancedDryRunParallel(
	planned []state.SliceInfo,
	max int,
	config sliceConfig,
	stdout io.Writer,
) {
	prefix := tmuxSessionPrefix(config)
	project := projectName()
	sessionName := fmt.Sprintf("%s-%s", prefix, project)
	claudeFlags := strings.Join(config.Parallel.ClaudeFlags, " ")

	ids := make([]string, 0, len(planned))
	for _, slice := range planned {
		ids = append(ids, slice.ID)
	}
	fmt.Fprintf(stdout, "Would launch: %s (up to --max %d)\n", strings.Join(ids, ", "), max)
	fmt.Fprintf(stdout, "dry-run: parallel mode — tmux session %q with %d window(s) (max: %d)\n",
		sessionName, len(planned), max)

	for _, slice := range planned {
		command := fmt.Sprintf("claude -p '/gl:slice %s' %s", slice.ID, claudeFlags)
		fmt.Fprintf(stdout, "  window %s: %s\n", slice.ID, command)
	}
}

// printEnhancedDryRunSequential prints the would-launch section for sequential mode.
func printEnhancedDryRunSequential(
	planned []state.SliceInfo,
	flags sliceFlags,
	stdout io.Writer,
) {
	if len(planned) == 0 {
		fmt.Fprintln(stdout, "Would launch: none")
		return
	}

	targetSlice := planned[0]

	if flags.sequential {
		fmt.Fprintf(stdout, "Would launch: %s (sequential)\n", targetSlice.ID)
		fmt.Fprintf(stdout, "dry-run: sequential mode (--sequential flag set)\n")
	} else {
		fmt.Fprintf(stdout, "Would launch: %s (sequential)\n", targetSlice.ID)
		fmt.Fprintf(stdout, "dry-run: sequential mode (tmux not available)\n")
	}
	fmt.Fprintf(stdout, "would run: slice %s\n", targetSlice.ID)
}

// runParallel handles parallel execution via tmux sessions.
// In dry-run mode it prints the plan without spawning anything.
func runParallel(
	readySlices []state.SliceInfo,
	flags sliceFlags,
	config sliceConfig,
	stdout io.Writer,
) int {
	// Limit to max windows.
	planned := readySlices
	if len(planned) > flags.max {
		planned = planned[:flags.max]
	}

	if flags.dryRun {
		printDryRunParallel(planned, flags.max, config, stdout)
		return 0
	}

	return executeTmuxSession(planned, flags.max, config, stdout)
}

// printDryRunParallel writes the parallel tmux plan to stdout without spawning.
func printDryRunParallel(planned []state.SliceInfo, max int, config sliceConfig, stdout io.Writer) {
	prefix := tmuxSessionPrefix(config)
	project := projectName()
	sessionName := fmt.Sprintf("%s-%s", prefix, project)
	claudeFlags := strings.Join(config.Parallel.ClaudeFlags, " ")

	fmt.Fprintf(stdout, "dry-run: parallel mode — tmux session %q with %d window(s) (max: %d)\n",
		sessionName, len(planned), max)

	for _, slice := range planned {
		command := fmt.Sprintf("claude -p '/gl:slice %s' %s", slice.ID, claudeFlags)
		fmt.Fprintf(stdout, "  window %s: %s\n", slice.ID, command)
	}
}

// executeTmuxSession creates the tmux session and attaches to it.
// Falls back to sequential on any tmux error.
func executeTmuxSession(
	planned []state.SliceInfo,
	max int,
	config sliceConfig,
	stdout io.Writer,
) int {
	prefix := tmuxSessionPrefix(config)
	project := projectName()
	sessionName := fmt.Sprintf("%s-%s", prefix, project)

	firstSlice := planned[0]
	claudeFlags := strings.Join(config.Parallel.ClaudeFlags, " ")
	firstCommand := fmt.Sprintf("claude -p '/gl:slice %s' %s", firstSlice.ID, claudeFlags)

	createError := tmux.NewSession(tmux.SessionOptions{
		Name:    sessionName,
		Window:  firstSlice.ID,
		Command: firstCommand,
	})
	if createError != nil {
		fmt.Fprintf(stdout, "warn: tmux session create failed (%v), falling back to sequential\n", createError)
		return spawnSlice(firstSlice.ID, config, stdout)
	}

	for _, slice := range planned[1:] {
		windowCommand := fmt.Sprintf("claude -p '/gl:slice %s' %s", slice.ID, claudeFlags)
		addError := tmux.AddWindow(sessionName, slice.ID, windowCommand)
		if addError != nil {
			fmt.Fprintf(stdout, "warn: failed to add window for %s: %v\n", slice.ID, addError)
		}
	}

	_ = max // max already applied when building planned slice

	attachError := tmux.AttachSession(sessionName)
	if attachError != nil {
		fmt.Fprintf(stdout, "warn: could not attach to tmux session %q: %v\n", sessionName, attachError)
		return 1
	}

	return 0
}

// runSequential handles sequential fallback execution.
// In dry-run mode it prints the sequential plan without spawning anything.
func runSequential(
	readySlices []state.SliceInfo,
	flags sliceFlags,
	config sliceConfig,
	executionContext state.ExecutionContext,
	stdout io.Writer,
) int {
	if flags.dryRun {
		printDryRunSequential(readySlices, flags.sequential, stdout)
		return 0
	}

	// Non-dry-run sequential: run the first ready slice.
	targetSlice := readySlices[0]
	return spawnSlice(targetSlice.ID, config, stdout)
}

// printDryRunSequential writes the sequential plan to stdout without spawning.
func printDryRunSequential(readySlices []state.SliceInfo, isSequentialFlag bool, stdout io.Writer) {
	if isSequentialFlag {
		fmt.Fprintln(stdout, "dry-run: sequential mode (--sequential flag set)")
	} else {
		fmt.Fprintln(stdout, "dry-run: sequential mode (tmux not available)")
	}

	if len(readySlices) > 0 {
		targetSlice := readySlices[0]
		fmt.Fprintf(stdout, "would run: slice %s\n", targetSlice.ID)
	}
}

// printDryRun writes what would happen without spawning any process.
func printDryRun(sliceID string, config sliceConfig, stdout io.Writer) {
	flags := strings.Join(config.Parallel.ClaudeFlags, " ")
	fmt.Fprintf(stdout, "dry-run: would execute slice %s\n", sliceID)
	fmt.Fprintf(stdout, "command: claude -p '/gl:slice %s' %s\n", sliceID, flags)
}

// printInsideClaudeInfo writes slice info to stdout for the Claude skill to consume.
func printInsideClaudeInfo(sliceID string, graph *state.Graph, stdout io.Writer) {
	graphSlice := graph.Slices[sliceID]
	fmt.Fprintf(stdout, "slice: %s — %s\n", sliceID, graphSlice.Name)
	fmt.Fprintf(stdout, "run: /gl:slice %s\n", sliceID)
}

// spawnSlice builds and starts a headless claude process for the given slice.
func spawnSlice(sliceID string, config sliceConfig, stdout io.Writer) int {
	prompt := fmt.Sprintf("/gl:slice %s", sliceID)

	spawnedCmd, spawnError := process.SpawnClaude(process.SpawnOptions{
		Prompt: prompt,
		Flags:  config.Parallel.ClaudeFlags,
		Stdout: stdout,
		Stderr: stdout,
	})

	if spawnError != nil {
		fmt.Fprintf(stdout, "error: could not spawn claude for slice %s: %v\n", sliceID, spawnError)
		return 1
	}

	waitError := spawnedCmd.Wait()
	if waitError != nil {
		return 1
	}

	return 0
}
