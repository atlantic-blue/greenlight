package cmd

import (
	"fmt"
	"io"

	"github.com/atlantic-blue/greenlight/internal/installer"
)

// RunCheck handles the "check" subcommand.
func RunCheck(args []string, stdout io.Writer) int {
	scope, _, err := ParseScope(args)
	if err != nil {
		fmt.Fprintf(stdout, "error: %v\n", err)
		return 1
	}

	targetDir, err := ResolveDir(scope)
	if err != nil {
		fmt.Fprintf(stdout, "error: %v\n", err)
		return 1
	}

	if !installer.Check(targetDir, scope, stdout) {
		return 1
	}
	return 0
}
