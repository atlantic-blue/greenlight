package process

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// LookPath is the PATH lookup function used to find the "claude" binary.
// It is a package-level variable so tests can replace it with a stub.
var LookPath = exec.LookPath

// Sentinel errors returned by this package.
var (
	ErrClaudeNotFound = errors.New("claude binary not found in PATH")
	ErrEmptyPrompt    = errors.New("prompt must not be empty")
	ErrStartFailure   = errors.New("failed to start claude process")
)

// SpawnOptions configures a headless (non-interactive) claude invocation.
type SpawnOptions struct {
	Prompt string
	Flags  []string
	Dir    string
	Stdout io.Writer
	Stderr io.Writer
}

// InteractiveOptions configures an interactive claude invocation.
type InteractiveOptions struct {
	Prompt string
	Flags  []string
	Dir    string
}

// verifyClaudeInPath checks that the "claude" binary is available in PATH.
// It returns ErrClaudeNotFound when the binary cannot be located.
func verifyClaudeInPath() error {
	_, lookupError := LookPath("claude")
	if lookupError != nil {
		return ErrClaudeNotFound
	}
	return nil
}

// filterDangerousFlags removes "--dangerously-skip-permissions" from the
// provided flags slice. This is a security invariant for interactive mode.
func filterDangerousFlags(flags []string) []string {
	filtered := make([]string, 0, len(flags))
	for _, flag := range flags {
		if flag == "--dangerously-skip-permissions" {
			continue
		}
		filtered = append(filtered, flag)
	}
	return filtered
}

// BuildClaudeCmd constructs an exec.Cmd for a headless claude invocation but
// does not start it. It returns ErrClaudeNotFound if "claude" is absent from
// PATH and ErrEmptyPrompt if opts.Prompt is empty.
func BuildClaudeCmd(opts SpawnOptions) (*exec.Cmd, error) {
	if pathError := verifyClaudeInPath(); pathError != nil {
		return nil, pathError
	}

	if opts.Prompt == "" {
		return nil, ErrEmptyPrompt
	}

	args := make([]string, 0, 2+len(opts.Flags))
	args = append(args, "-p", opts.Prompt)
	args = append(args, opts.Flags...)

	command := exec.Command("claude", args...)
	command.Dir = opts.Dir
	command.Stdout = opts.Stdout
	command.Stderr = opts.Stderr

	return command, nil
}

// SpawnClaude builds and starts a headless claude process. The process is
// started but not waited on — callers are responsible for calling cmd.Wait.
// It returns ErrClaudeNotFound, ErrEmptyPrompt, or a wrapped ErrStartFailure.
func SpawnClaude(opts SpawnOptions) (*exec.Cmd, error) {
	command, buildError := BuildClaudeCmd(opts)
	if buildError != nil {
		return nil, buildError
	}

	if startError := command.Start(); startError != nil {
		return nil, fmt.Errorf("%w: %w", ErrStartFailure, startError)
	}

	return command, nil
}

// BuildInteractiveCmd constructs an exec.Cmd for an interactive claude session
// but does not start it. The "-p" flag is added only when opts.Prompt is
// non-empty. "--dangerously-skip-permissions" is always stripped from Flags,
// regardless of what the caller passes. Returns ErrClaudeNotFound if "claude"
// is absent from PATH.
func BuildInteractiveCmd(opts InteractiveOptions) (*exec.Cmd, error) {
	if pathError := verifyClaudeInPath(); pathError != nil {
		return nil, pathError
	}

	safeFlags := filterDangerousFlags(opts.Flags)

	args := make([]string, 0, 2+len(safeFlags))
	if opts.Prompt != "" {
		args = append(args, "-p", opts.Prompt)
	}
	args = append(args, safeFlags...)

	command := exec.Command("claude", args...)
	command.Dir = opts.Dir

	return command, nil
}

// SpawnInteractive runs an interactive claude session and blocks until the
// process exits. Stdin, Stdout, and Stderr are connected to the terminal.
// Returns ErrClaudeNotFound if "claude" is absent from PATH, or any error
// returned by the process on exit.
func SpawnInteractive(opts InteractiveOptions) error {
	command, buildError := BuildInteractiveCmd(opts)
	if buildError != nil {
		return buildError
	}

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}
