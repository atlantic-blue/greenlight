package cmd

import (
	"fmt"
	"io"
)

// RunSlice handles the "slice" subcommand.
func RunSlice(args []string, stdout io.Writer) int {
	fmt.Fprintln(stdout, "slice: not implemented yet")
	return 0
}
