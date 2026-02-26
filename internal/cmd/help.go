package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/atlantic-blue/greenlight/internal/state"
	"github.com/atlantic-blue/greenlight/internal/version"
)

// RunHelp handles the "help" subcommand.
//
// Always prints the command listing grouped by category, then appends a
// project state summary if .greenlight/ exists, or a suggestion to run
// 'gl init' if it does not. Always returns 0.
func RunHelp(args []string, stdout io.Writer) int {
	writeCommandListing(stdout)
	writeProjectSummary(stdout)
	return 0
}

// writeCommandListing prints the grouped command reference to the writer.
func writeCommandListing(stdout io.Writer) {
	fmt.Fprintf(stdout, "gl %s\n\n", version.Version)
	fmt.Fprint(stdout, `Usage: gl <command> [flags]

Project lifecycle:
  init        Initialise a new greenlight project
  design      Run the design phase for a feature
  roadmap     View or update the project roadmap

Building:
  slice       Run a vertical slice end-to-end

State & progress:
  status      Show current project status
  changelog   View or generate the changelog

Admin:
  install     Install greenlight files
  uninstall   Remove greenlight files
  check       Verify installation
  version     Show version information
  help        Show this help

`)
}

// writeProjectSummary appends a one-line project state summary when
// .greenlight/ exists, or a hint to run 'gl init' when it does not.
// All state-read errors are silently ignored — the summary is best-effort.
func writeProjectSummary(stdout io.Writer) {
	if _, statError := os.Stat(greenlightDir); os.IsNotExist(statError) {
		fmt.Fprintln(stdout, "Run 'gl init' to start a new project.")
		return
	}

	total, complete, ready := readProjectCounts()
	fmt.Fprintf(stdout, "Current project: %d slices, %d complete, %d ready\n", total, complete, ready)
}

// readProjectCounts reads slice data from disk and returns the total, complete,
// and ready counts. On any read error the counts default to zero.
func readProjectCounts() (total, complete, ready int) {
	slices, readError := state.ReadSlices(slicesDir)
	if readError != nil {
		return 0, 0, 0
	}

	total = len(slices)
	for _, slice := range slices {
		if slice.Status == "complete" {
			complete++
		}
	}

	graph, graphError := state.ReadGraph(graphPath)
	if graphError != nil {
		return total, complete, 0
	}

	readySlices := state.FindReadySlices(slices, graph)
	ready = len(readySlices)

	return total, complete, ready
}
