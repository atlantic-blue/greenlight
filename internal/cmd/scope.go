package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/atlantic-blue/greenlight/internal/installer"
)

// ParseScope extracts --global or --local from args and returns the scope
// along with any remaining arguments.
func ParseScope(args []string) (scope string, remaining []string, err error) {
	for _, arg := range args {
		switch arg {
		case "--global":
			if scope != "" {
				return "", nil, fmt.Errorf("cannot specify both --global and --local")
			}
			scope = "global"
		case "--local":
			if scope != "" {
				return "", nil, fmt.Errorf("cannot specify both --global and --local")
			}
			scope = "local"
		default:
			remaining = append(remaining, arg)
		}
	}
	if scope == "" {
		return "", nil, fmt.Errorf("must specify --global or --local")
	}
	return scope, remaining, nil
}

// ResolveDir returns the target directory for the given scope.
// Global: ~/.claude   Local: ./.claude
func ResolveDir(scope string) (string, error) {
	switch scope {
	case "global":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot determine home directory: %w", err)
		}
		return filepath.Join(home, ".claude"), nil
	case "local":
		return ".claude", nil
	default:
		return "", fmt.Errorf("unknown scope: %s", scope)
	}
}

// ParseConflictStrategy extracts --on-conflict=<strategy> from args.
// Returns the strategy and remaining args. Defaults to "keep".
func ParseConflictStrategy(args []string) (installer.ConflictStrategy, []string) {
	strategy := installer.ConflictKeep
	var remaining []string
	for _, arg := range args {
		if len(arg) > 14 && arg[:14] == "--on-conflict=" {
			val := arg[14:]
			switch installer.ConflictStrategy(val) {
			case installer.ConflictKeep, installer.ConflictReplace, installer.ConflictAppend:
				strategy = installer.ConflictStrategy(val)
			}
		} else {
			remaining = append(remaining, arg)
		}
	}
	return strategy, remaining
}
