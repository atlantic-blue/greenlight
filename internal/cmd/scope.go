package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/atlantic-blue/greenlight/internal/installer"
)

var ErrBothScopes = errors.New("cannot specify both --global and --local")
var ErrNoScope = errors.New("must specify --global or --local")
var ErrUnknownScope = errors.New("unknown scope")
var ErrInvalidConflictStrategy = errors.New("invalid --on-conflict value")

// ParseScope extracts --global or --local from args and returns the scope
// along with any remaining arguments.
func ParseScope(args []string) (scope string, remaining []string, err error) {
	remaining = []string{}
	for _, arg := range args {
		switch arg {
		case "--global":
			if scope != "" {
				return "", nil, ErrBothScopes
			}
			scope = "global"
		case "--local":
			if scope != "" {
				return "", nil, ErrBothScopes
			}
			scope = "local"
		default:
			remaining = append(remaining, arg)
		}
	}
	if scope == "" {
		return "", nil, ErrNoScope
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
		return "", ErrUnknownScope
	}
}

// ParseConflictStrategy extracts --on-conflict=<strategy> from args.
// Returns the strategy, remaining args, and error. Defaults to "keep".
func ParseConflictStrategy(args []string) (installer.ConflictStrategy, []string, error) {
	strategy := installer.ConflictKeep
	remaining := []string{}
	for _, arg := range args {
		if len(arg) >= 14 && arg[:14] == "--on-conflict=" {
			val := arg[14:]
			switch installer.ConflictStrategy(val) {
			case installer.ConflictKeep, installer.ConflictReplace, installer.ConflictAppend:
				strategy = installer.ConflictStrategy(val)
			default:
				return "", nil, ErrInvalidConflictStrategy
			}
		} else {
			remaining = append(remaining, arg)
		}
	}
	return strategy, remaining, nil
}

// ParseVerifyFlag extracts --verify from args.
// Returns true if --verify is present, false otherwise, along with remaining args.
func ParseVerifyFlag(args []string) (verify bool, remaining []string) {
	remaining = []string{}
	for _, arg := range args {
		if arg == "--verify" {
			verify = true
		} else {
			remaining = append(remaining, arg)
		}
	}
	return verify, remaining
}
