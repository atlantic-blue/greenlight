package cli

import (
	"fmt"
	"io"
	"io/fs"

	"github.com/atlantic-blue/greenlight/internal/cmd"
)

// Run dispatches to the appropriate subcommand based on args.
// contentFS provides the embedded source content.
func Run(args []string, contentFS fs.FS, stdout io.Writer) int {

	if len(args) == 0 {
		printUsage(stdout)
		return 0
	}

	switch args[0] {
	case "install":
		return cmd.RunInstall(args[1:], contentFS, stdout)
	case "uninstall":
		return cmd.RunUninstall(args[1:], stdout)
	case "check":
		return cmd.RunCheck(args[1:], contentFS, stdout)
	case "version":
		return cmd.RunVersion(stdout)
	case "status":
		return cmd.RunStatus(args[1:], stdout)
	case "slice":
		return cmd.RunSlice(args[1:], stdout)
	case "init":
		return cmd.RunInit(args[1:], stdout)
	case "design":
		return cmd.RunDesign(args[1:], stdout)
	case "roadmap":
		return cmd.RunRoadmap(args[1:], stdout)
	case "changelog":
		return cmd.RunChangelog(args[1:], stdout)
	case "help", "--help", "-h":
		return cmd.RunHelp(args[1:], stdout)
	default:
		fmt.Fprintf(stdout, "unknown command: %s\n\n", args[0])
		printUsage(stdout)
		return 1
	}
}

func printUsage(w io.Writer) {
	fmt.Fprint(w, `
   ╔══════════════════════════════════════════════════════════╗
   ║                                                          ║
   ║	 ▗▄▄▖▗▄▄▖ ▗▄▄▄▖▗▄▄▄▖▗▖  ▗▖ ▗▖   ▗▄▄▄▖ ▗▄▄▖▗▖ ▗▖▗▄▄▄▖  ║
   ║	▐▌   ▐▌ ▐▌▐▌   ▐▌   ▐▛▚▖▐▌ ▐▌     █  ▐▌   ▐▌ ▐▌  █    ║
   ║	▐▌▝▜▌▐▛▀▚▖▐▛▀▀▘▐▛▀▀▘▐▌ ▝▜▌ ▐▌     █  ▐▌▝▜▌▐▛▀▜▌  █    ║
   ║	▝▚▄▞▘▐▌ ▐▌▐▙▄▄▖▐▙▄▄▖▐▌  ▐▌ ▐▙▄▄▖▗▄█▄▖▝▚▄▞▘▐▌ ▐▌  █    ║
   ║                                                          ║
   ╚══════════════════════════════════════════════════════════╝
                  Tests are the source of truth.
                        Green means done.


Usage: greenlight <command> [flags]

Project lifecycle:
  init        Initialise a new greenlight project
  design      Run the design phase for a feature
  roadmap     View or update the project roadmap

Building:
  slice       Run a vertical slice end-to-end

State & progress:
  status      Show current project status
  changelog   View or generate the changelog

Admin:
  install     Install greenlight files
  uninstall   Remove greenlight files
  check       Verify installation
  version     Show version information
  help        Show this help

Install flags:
  --global                Install to ~/.claude/
  --local                 Install to ./.claude/
  --on-conflict=<mode>    Handle existing CLAUDE.md: keep (default), replace, append

Examples:
  greenlight install --global
  greenlight install --local --on-conflict=replace
  greenlight uninstall --global
  greenlight check --local
`)
}
