package cmd

import (
	"fmt"
	"io"
)

// RunStatus handles the "status" subcommand.
func RunStatus(args []string, stdout io.Writer) int {
	fmt.Fprintln(stdout, "status: not implemented yet")
	return 0
}
