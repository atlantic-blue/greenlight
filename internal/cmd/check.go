package cmd

import (
	"fmt"
	"io"
	"io/fs"

	"github.com/atlantic-blue/greenlight/internal/installer"
)

// RunCheck handles the "check" subcommand.
func RunCheck(args []string, contentFS fs.FS, stdout io.Writer) int {
	scope, remaining, err := ParseScope(args)
	if err != nil {
		fmt.Fprintf(stdout, "error: %v\n", err)
		return 1
	}

	verify, _ := ParseVerifyFlag(remaining)

	targetDir, err := ResolveDir(scope)
	if err != nil {
		fmt.Fprintf(stdout, "error: %v\n", err)
		return 1
	}

	if !installer.Check(targetDir, scope, stdout, verify, contentFS) {
		return 1
	}
	return 0
}
