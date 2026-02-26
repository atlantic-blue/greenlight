package cmd

import (
	"fmt"
	"io"
	"os"
)

const (
	greenlightDir = ".greenlight"
	roadmapPath   = ".greenlight/ROADMAP.md"
)

// RunRoadmap handles the "roadmap" subcommand.
// It reads .greenlight/ROADMAP.md and prints the contents verbatim to stdout.
// Returns 1 if .greenlight/ does not exist or ROADMAP.md is missing.
func RunRoadmap(args []string, stdout io.Writer) int {
	if _, statError := os.Stat(greenlightDir); os.IsNotExist(statError) {
		fmt.Fprintln(stdout, "error: .greenlight/ directory not found. Run 'greenlight init' to set up this greenlight project.")
		return 1
	}

	contents, readError := os.ReadFile(roadmapPath)
	if readError != nil {
		fmt.Fprintln(stdout, "error: ROADMAP.md not found. Use 'greenlight design' to create a roadmap for this project.")
		return 1
	}

	fmt.Fprint(stdout, string(contents))
	return 0
}
