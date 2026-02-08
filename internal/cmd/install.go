package cmd

import (
	"fmt"
	"io"
	"io/fs"

	"github.com/atlantic-blue/greenlight/internal/installer"
)

// RunInstall handles the "install" subcommand.
func RunInstall(args []string, contentFS fs.FS, stdout io.Writer) int {
	strategy, args, err := ParseConflictStrategy(args)
	if err != nil {
		fmt.Fprintf(stdout, "error: %v\n", err)
		return 1
	}

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

	inst := installer.New(contentFS, stdout)
	if err := inst.Install(targetDir, scope, strategy); err != nil {
		fmt.Fprintf(stdout, "error: %v\n", err)
		return 1
	}
	return 0
}
