package cmd

import (
	"fmt"
	"io"
)

// RunInit handles the "init" subcommand.
func RunInit(args []string, stdout io.Writer) int {
	fmt.Fprintln(stdout, "init: not implemented yet")
	return 0
}
