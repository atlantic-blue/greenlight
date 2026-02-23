package cmd

import (
	"fmt"
	"io"
)

// RunChangelog handles the "changelog" subcommand.
func RunChangelog(args []string, stdout io.Writer) int {
	fmt.Fprintln(stdout, "changelog: not implemented yet")
	return 0
}
