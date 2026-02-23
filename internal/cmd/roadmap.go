package cmd

import (
	"fmt"
	"io"
)

// RunRoadmap handles the "roadmap" subcommand.
func RunRoadmap(args []string, stdout io.Writer) int {
	fmt.Fprintln(stdout, "roadmap: not implemented yet")
	return 0
}
