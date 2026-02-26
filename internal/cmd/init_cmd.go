package cmd

import (
	"errors"
	"fmt"
	"io"

	"github.com/atlantic-blue/greenlight/internal/process"
	"github.com/atlantic-blue/greenlight/internal/state"
)

// RunInit handles the "init" subcommand.
//
// Inside Claude context: prints instructions for /gl:init skill, returns 0.
// Shell context: prints launching message and spawns an interactive Claude
// session with the /gl:init prompt, blocking until the session ends.
//
// RunInit does NOT require .greenlight/ to exist — init is meant to create it.
func RunInit(args []string, stdout io.Writer) int {
	executionContext := state.DetectContext()

	if executionContext.InsideClaude {
		fmt.Fprintln(stdout, "Run the /gl:init skill to initialise this project.")
		return 0
	}

	fmt.Fprintln(stdout, "Launching Greenlight init...")

	spawnError := process.SpawnInteractive(process.InteractiveOptions{
		Prompt: "/gl:init",
	})

	if spawnError != nil {
		if errors.Is(spawnError, process.ErrClaudeNotFound) {
			fmt.Fprintln(stdout, "Error: claude binary not found in PATH. Install claude to use this command.")
			return 1
		}
		fmt.Fprintf(stdout, "Error: failed to launch interactive session: %v\n", spawnError)
		return 1
	}

	return 0
}
