package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/atlantic-blue/greenlight/internal/process"
	"github.com/atlantic-blue/greenlight/internal/state"
)

const configPath = ".greenlight/config.json"

// sliceConfig holds the parsed structure of .greenlight/config.json.
type sliceConfig struct {
	Parallel struct {
		ClaudeFlags []string `json:"claude_flags"`
	} `json:"parallel"`
}

// sliceFlags holds the parsed command-line flags for the slice subcommand.
type sliceFlags struct {
	dryRun     bool
	sequential bool
}

// RunSlice handles the "slice" subcommand.
//
// With a slice ID argument: validates it against GRAPH.json, then runs it
// headlessly via process.SpawnClaude (shell context) or prints info (inside Claude).
//
// Without a slice ID: auto-detects ready slices via state.FindReadySlices and
// runs the first one ordered by wave then ID.
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
	var flags sliceFlags
	var remaining []string

	for _, arg := range args {
		switch arg {
		case "--dry-run":
			flags.dryRun = true
		case "--sequential":
			flags.sequential = true
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

// runAutoDetect handles the no-slice-ID case: finds ready slices and runs first.
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

	readySlices := state.FindReadySlices(slices, graph)

	if len(readySlices) == 0 {
		fmt.Fprintln(stdout, "No ready slices: all pending slices are blocked waiting on dependencies.")
		return 0
	}

	targetSlice := readySlices[0]

	if len(readySlices) > 1 {
		remaining := len(readySlices) - 1
		fmt.Fprintf(stdout, "hint: %d more ready slice(s) available. Use --max or parallel mode (S-44) to run concurrently.\n", remaining)
	}

	return runNamedSlice(targetSlice.ID, flags, config, executionContext, stdout)
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
