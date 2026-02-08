package cli

import (
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/atlantic-blue/greenlight/internal/cmd"
)

// Run dispatches to the appropriate subcommand based on args.
// contentFS provides the embedded source content.
func Run(args []string, contentFS fs.FS) int {
	stdout := os.Stdout

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
		return cmd.RunCheck(args[1:], stdout)
	case "version":
		return cmd.RunVersion(stdout)
	case "help", "--help", "-h":
		printUsage(stdout)
		return 0
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

Commands:
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
