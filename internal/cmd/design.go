package cmd

import (
	"fmt"
	"io"
)

// RunDesign handles the "design" subcommand.
func RunDesign(args []string, stdout io.Writer) int {
	fmt.Fprintln(stdout, "design: not implemented yet")
	return 0
}
