package tmux

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// LookPath is the PATH lookup function used to find the "tmux" binary.
// It is a package-level variable so tests can replace it with a stub.
var LookPath = exec.LookPath

// Sentinel errors returned by this package.
var (
	ErrTmuxNotFound    = errors.New("tmux not found in PATH")
	ErrSessionExists   = errors.New("tmux session already exists")
	ErrCreateFailed    = errors.New("failed to create tmux session")
	ErrSessionNotFound = errors.New("tmux session not found")
	ErrAddWindowFailed = errors.New("failed to add tmux window")
	ErrAttachFailed    = errors.New("failed to attach to tmux session")
)

// SessionOptions configures a new tmux session.
type SessionOptions struct {
	Name    string
	Dir     string
	Command string
	Window  string
}

// verifyTmuxInPath checks that the "tmux" binary is available in PATH.
// It returns ErrTmuxNotFound when the binary cannot be located.
func verifyTmuxInPath() error {
	_, lookupError := LookPath("tmux")
	if lookupError != nil {
		return ErrTmuxNotFound
	}
	return nil
}

// IsAvailable reports whether tmux is available in the system PATH.
func IsAvailable() bool {
	_, lookupError := LookPath("tmux")
	return lookupError == nil
}

// BuildNewSessionCmd constructs an exec.Cmd for creating a new detached tmux
// session but does not start it. It returns ErrTmuxNotFound if "tmux" is
// absent from PATH.
func BuildNewSessionCmd(opts SessionOptions) (*exec.Cmd, error) {
	if pathError := verifyTmuxInPath(); pathError != nil {
		return nil, pathError
	}

	command := exec.Command(
		"tmux", "new-session",
		"-d",
		"-s", opts.Name,
		"-n", opts.Window,
		"-c", opts.Dir,
		opts.Command,
	)

	return command, nil
}

// NewSession creates a new detached tmux session using the given options.
// It returns ErrTmuxNotFound if tmux is absent from PATH, or a wrapped
// ErrCreateFailed if the tmux command itself fails.
func NewSession(opts SessionOptions) error {
	command, buildError := BuildNewSessionCmd(opts)
	if buildError != nil {
		return buildError
	}

	if runError := command.Run(); runError != nil {
		return fmt.Errorf("%w: %w", ErrCreateFailed, runError)
	}

	return nil
}

// BuildAddWindowCmd constructs an exec.Cmd for adding a new window to an
// existing tmux session but does not start it. It returns ErrTmuxNotFound if
// "tmux" is absent from PATH.
func BuildAddWindowCmd(session, name, command string) (*exec.Cmd, error) {
	if pathError := verifyTmuxInPath(); pathError != nil {
		return nil, pathError
	}

	cmd := exec.Command(
		"tmux", "new-window",
		"-t", session,
		"-n", name,
		command,
	)

	return cmd, nil
}

// AddWindow adds a new window to an existing tmux session and runs the given
// command inside it. It returns ErrTmuxNotFound if tmux is absent from PATH,
// or a wrapped ErrAddWindowFailed if the tmux command itself fails.
func AddWindow(session, name, command string) error {
	cmd, buildError := BuildAddWindowCmd(session, name, command)
	if buildError != nil {
		return buildError
	}

	if runError := cmd.Run(); runError != nil {
		return fmt.Errorf("%w: %w", ErrAddWindowFailed, runError)
	}

	return nil
}

// BuildAttachCmd constructs an exec.Cmd for attaching to an existing tmux
// session but does not start it. It returns ErrTmuxNotFound if "tmux" is
// absent from PATH.
func BuildAttachCmd(session string) (*exec.Cmd, error) {
	if pathError := verifyTmuxInPath(); pathError != nil {
		return nil, pathError
	}

	command := exec.Command("tmux", "attach-session", "-t", session)

	return command, nil
}

// AttachSession attaches the current terminal to an existing tmux session and
// blocks until the session ends. Stdin, Stdout, and Stderr are connected to
// the terminal. It returns ErrTmuxNotFound if tmux is absent from PATH, or a
// wrapped ErrAttachFailed if the tmux command itself fails.
func AttachSession(session string) error {
	command, buildError := BuildAttachCmd(session)
	if buildError != nil {
		return buildError
	}

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if runError := command.Run(); runError != nil {
		return fmt.Errorf("%w: %w", ErrAttachFailed, runError)
	}

	return nil
}
