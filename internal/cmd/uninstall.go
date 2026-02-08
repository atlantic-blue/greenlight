package cmd

import (
	"fmt"
	"io"

	"github.com/atlantic-blue/greenlight/internal/installer"
)

// RunUninstall handles the "uninstall" subcommand.
func RunUninstall(args []string, stdout io.Writer) int {
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

	if err := installer.Uninstall(targetDir, stdout); err != nil {
		fmt.Fprintf(stdout, "error: %v\n", err)
		return 1
	}
	return 0
}
