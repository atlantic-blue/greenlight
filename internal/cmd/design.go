package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/atlantic-blue/greenlight/internal/process"
	"github.com/atlantic-blue/greenlight/internal/state"
)

// RunDesign handles the "design" subcommand.
//
// Inside Claude context: prints instructions for /gl:design skill, returns 0.
// Shell context: verifies .greenlight/ exists, then prints a launching message
// and spawns an interactive Claude session with the /gl:design prompt, blocking
// until the session ends.
func RunDesign(args []string, stdout io.Writer) int {
	executionContext := state.DetectContext()

	if executionContext.InsideClaude {
		fmt.Fprintln(stdout, "Run the /gl:design skill to design slices for this project.")
		return 0
	}

	if _, statError := os.Stat(greenlightDir); os.IsNotExist(statError) {
		fmt.Fprintln(stdout, "Not a greenlight project. Run 'gl init' first.")
		return 1
	}

	fmt.Fprintln(stdout, "Launching Greenlight design session...")

	spawnError := process.SpawnInteractive(process.InteractiveOptions{
		Prompt: "/gl:design",
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
