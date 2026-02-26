package tmux_test

// S-43 Tests: Tmux Manager
// Covers C-107 (TmuxIsAvailable), C-108 (TmuxNewSession / BuildNewSessionCmd),
// C-109 (TmuxAddWindow / BuildAddWindowCmd), and C-110 (TmuxAttachSession / BuildAttachCmd).
//
// Contract C-107 — IsAvailable:
//   - Returns true when tmux binary is in PATH
//   - Returns false when tmux binary is not in PATH
//   - Never returns an error; always returns a plain bool
//
// Contract C-108 — NewSession / BuildNewSessionCmd:
//   - Built command includes "-d" (detached flag)
//   - Built command includes "-s" followed by session name
//   - Built command includes "-n" followed by window name
//   - Built command includes "-c" followed by working directory
//   - Built command includes the initial command string
//   - Returns ErrTmuxNotFound when tmux is absent from PATH
//   - NewSession propagates ErrTmuxNotFound before attempting to start a process
//
// Contract C-109 — AddWindow / BuildAddWindowCmd:
//   - Built command includes "-t" followed by session name
//   - Built command includes "-n" followed by window name
//   - Built command includes the command string
//   - Returns ErrTmuxNotFound when tmux is absent from PATH
//
// Contract C-110 — AttachSession / BuildAttachCmd:
//   - Built command includes "-t" followed by session name
//   - Args[1] is "attach-session"
//   - Returns ErrTmuxNotFound when tmux is absent from PATH

import (
	"errors"
	"os/exec"
	"testing"

	"github.com/atlantic-blue/greenlight/internal/tmux"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// tmuxFoundLookPath returns a LookPath stub that reports "tmux" as found.
func tmuxFoundLookPath(name string) (string, error) {
	if name == "tmux" {
		return "/usr/bin/tmux", nil
	}
	return "", exec.ErrNotFound
}

// tmuxNotFoundLookPath returns a LookPath stub that always reports not found.
func tmuxNotFoundLookPath(_ string) (string, error) {
	return "", exec.ErrNotFound
}

// overrideLookPath replaces tmux.LookPath with the given stub and registers
// a cleanup to restore the original value when the test finishes.
func overrideLookPath(t *testing.T, stub func(string) (string, error)) {
	t.Helper()
	original := tmux.LookPath
	tmux.LookPath = stub
	t.Cleanup(func() { tmux.LookPath = original })
}

// containsArg reports whether slice contains the given value.
func containsArg(args []string, value string) bool {
	for _, arg := range args {
		if arg == value {
			return true
		}
	}
	return false
}

// argAfter returns the element immediately following needle in args, and
// whether needle was found with a successor element.
func argAfter(args []string, needle string) (string, bool) {
	for index, arg := range args {
		if arg == needle && index+1 < len(args) {
			return args[index+1], true
		}
	}
	return "", false
}

// ---------------------------------------------------------------------------
// C-107: IsAvailable
// ---------------------------------------------------------------------------

// TestIsAvailable_TmuxInPath_ReturnsTrue verifies that IsAvailable returns
// true when the tmux binary is present in PATH.
func TestIsAvailable_TmuxInPath_ReturnsTrue(t *testing.T) {
	overrideLookPath(t, tmuxFoundLookPath)

	if !tmux.IsAvailable() {
		t.Error("IsAvailable() = false; want true when tmux is in PATH")
	}
}

// TestIsAvailable_TmuxNotInPath_ReturnsFalse verifies that IsAvailable returns
// false when the tmux binary is absent from PATH.
func TestIsAvailable_TmuxNotInPath_ReturnsFalse(t *testing.T) {
	overrideLookPath(t, tmuxNotFoundLookPath)

	if tmux.IsAvailable() {
		t.Error("IsAvailable() = true; want false when tmux is not in PATH")
	}
}

// TestIsAvailable_NeverReturnsError verifies the behavioural contract that
// IsAvailable always returns a plain bool in both the found and not-found cases,
// exercising both branches to confirm no panic or unexpected state.
func TestIsAvailable_NeverReturnsError(t *testing.T) {
	cases := []struct {
		name string
		stub func(string) (string, error)
		want bool
	}{
		{"tmux found", tmuxFoundLookPath, true},
		{"tmux not found", tmuxNotFoundLookPath, false},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			overrideLookPath(t, testCase.stub)
			got := tmux.IsAvailable()
			if got != testCase.want {
				t.Errorf("IsAvailable() = %v; want %v", got, testCase.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// C-108: BuildNewSessionCmd
// ---------------------------------------------------------------------------

// TestBuildNewSessionCmd_IncludesDetachedFlag verifies that the built command's
// Args contain the "-d" flag, which runs the session in the background.
func TestBuildNewSessionCmd_IncludesDetachedFlag(t *testing.T) {
	overrideLookPath(t, tmuxFoundLookPath)

	cmd, err := tmux.BuildNewSessionCmd(tmux.SessionOptions{
		Name:    "my-session",
		Window:  "main",
		Dir:     "/tmp",
		Command: "bash",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if !containsArg(cmd.Args, "-d") {
		t.Errorf("cmd.Args does not contain \"-d\" detached flag; got: %v", cmd.Args)
	}
}

// TestBuildNewSessionCmd_IncludesSessionName verifies that the built command's
// Args contain "-s" followed immediately by the session name.
func TestBuildNewSessionCmd_IncludesSessionName(t *testing.T) {
	overrideLookPath(t, tmuxFoundLookPath)

	const sessionName = "my-session"
	cmd, err := tmux.BuildNewSessionCmd(tmux.SessionOptions{
		Name:    sessionName,
		Window:  "main",
		Dir:     "/tmp",
		Command: "bash",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	got, found := argAfter(cmd.Args, "-s")
	if !found {
		t.Fatalf("cmd.Args does not contain \"-s\" flag; got: %v", cmd.Args)
	}
	if got != sessionName {
		t.Errorf("arg after \"-s\" = %q; want %q", got, sessionName)
	}
}

// TestBuildNewSessionCmd_IncludesWindowName verifies that the built command's
// Args contain "-n" followed immediately by the window name.
func TestBuildNewSessionCmd_IncludesWindowName(t *testing.T) {
	overrideLookPath(t, tmuxFoundLookPath)

	const windowName = "editor"
	cmd, err := tmux.BuildNewSessionCmd(tmux.SessionOptions{
		Name:    "my-session",
		Window:  windowName,
		Dir:     "/tmp",
		Command: "vim",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	got, found := argAfter(cmd.Args, "-n")
	if !found {
		t.Fatalf("cmd.Args does not contain \"-n\" flag; got: %v", cmd.Args)
	}
	if got != windowName {
		t.Errorf("arg after \"-n\" = %q; want %q", got, windowName)
	}
}

// TestBuildNewSessionCmd_IncludesWorkingDir verifies that the built command's
// Args contain "-c" followed immediately by the working directory path.
func TestBuildNewSessionCmd_IncludesWorkingDir(t *testing.T) {
	overrideLookPath(t, tmuxFoundLookPath)

	const workDir = "/home/user/project"
	cmd, err := tmux.BuildNewSessionCmd(tmux.SessionOptions{
		Name:    "my-session",
		Window:  "main",
		Dir:     workDir,
		Command: "bash",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	got, found := argAfter(cmd.Args, "-c")
	if !found {
		t.Fatalf("cmd.Args does not contain \"-c\" flag; got: %v", cmd.Args)
	}
	if got != workDir {
		t.Errorf("arg after \"-c\" = %q; want %q", got, workDir)
	}
}

// TestBuildNewSessionCmd_IncludesCommand verifies that the built command's
// Args contain the initial command string to be executed in the session.
func TestBuildNewSessionCmd_IncludesCommand(t *testing.T) {
	overrideLookPath(t, tmuxFoundLookPath)

	const initialCommand = "htop"
	cmd, err := tmux.BuildNewSessionCmd(tmux.SessionOptions{
		Name:    "my-session",
		Window:  "monitor",
		Dir:     "/tmp",
		Command: initialCommand,
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if !containsArg(cmd.Args, initialCommand) {
		t.Errorf("cmd.Args does not contain command %q; got: %v", initialCommand, cmd.Args)
	}
}

// TestBuildNewSessionCmd_TmuxNotInPath_ReturnsError verifies that
// BuildNewSessionCmd returns ErrTmuxNotFound when tmux is absent from PATH.
func TestBuildNewSessionCmd_TmuxNotInPath_ReturnsError(t *testing.T) {
	overrideLookPath(t, tmuxNotFoundLookPath)

	_, err := tmux.BuildNewSessionCmd(tmux.SessionOptions{
		Name:    "my-session",
		Window:  "main",
		Dir:     "/tmp",
		Command: "bash",
	})
	if err == nil {
		t.Fatal("expected ErrTmuxNotFound, got nil")
	}

	if !errors.Is(err, tmux.ErrTmuxNotFound) {
		t.Errorf("expected errors.Is(err, tmux.ErrTmuxNotFound); got: %v", err)
	}
}

// TestNewSession_TmuxNotInPath_ReturnsError verifies that NewSession propagates
// ErrTmuxNotFound (from the PATH check) before ever attempting to start tmux.
func TestNewSession_TmuxNotInPath_ReturnsError(t *testing.T) {
	overrideLookPath(t, tmuxNotFoundLookPath)

	err := tmux.NewSession(tmux.SessionOptions{
		Name:    "my-session",
		Window:  "main",
		Dir:     "/tmp",
		Command: "bash",
	})
	if err == nil {
		t.Fatal("expected ErrTmuxNotFound, got nil")
	}

	if !errors.Is(err, tmux.ErrTmuxNotFound) {
		t.Errorf("expected errors.Is(err, tmux.ErrTmuxNotFound); got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// C-109: BuildAddWindowCmd
// ---------------------------------------------------------------------------

// TestBuildAddWindowCmd_IncludesTargetSession verifies that the built command's
// Args contain "-t" followed immediately by the session name.
func TestBuildAddWindowCmd_IncludesTargetSession(t *testing.T) {
	overrideLookPath(t, tmuxFoundLookPath)

	const session = "my-session"
	cmd, err := tmux.BuildAddWindowCmd(session, "logs", "tail -f /var/log/syslog")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	got, found := argAfter(cmd.Args, "-t")
	if !found {
		t.Fatalf("cmd.Args does not contain \"-t\" flag; got: %v", cmd.Args)
	}
	if got != session {
		t.Errorf("arg after \"-t\" = %q; want %q", got, session)
	}
}

// TestBuildAddWindowCmd_IncludesWindowName verifies that the built command's
// Args contain "-n" followed immediately by the window name.
func TestBuildAddWindowCmd_IncludesWindowName(t *testing.T) {
	overrideLookPath(t, tmuxFoundLookPath)

	const windowName = "logs"
	cmd, err := tmux.BuildAddWindowCmd("my-session", windowName, "tail -f /var/log/syslog")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	got, found := argAfter(cmd.Args, "-n")
	if !found {
		t.Fatalf("cmd.Args does not contain \"-n\" flag; got: %v", cmd.Args)
	}
	if got != windowName {
		t.Errorf("arg after \"-n\" = %q; want %q", got, windowName)
	}
}

// TestBuildAddWindowCmd_IncludesCommand verifies that the built command's
// Args contain the command string to be run in the new window.
func TestBuildAddWindowCmd_IncludesCommand(t *testing.T) {
	overrideLookPath(t, tmuxFoundLookPath)

	const windowCommand = "tail -f /var/log/syslog"
	cmd, err := tmux.BuildAddWindowCmd("my-session", "logs", windowCommand)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if !containsArg(cmd.Args, windowCommand) {
		t.Errorf("cmd.Args does not contain command %q; got: %v", windowCommand, cmd.Args)
	}
}

// TestBuildAddWindowCmd_TmuxNotInPath_ReturnsError verifies that
// BuildAddWindowCmd returns ErrTmuxNotFound when tmux is absent from PATH.
func TestBuildAddWindowCmd_TmuxNotInPath_ReturnsError(t *testing.T) {
	overrideLookPath(t, tmuxNotFoundLookPath)

	_, err := tmux.BuildAddWindowCmd("my-session", "logs", "tail -f /var/log/syslog")
	if err == nil {
		t.Fatal("expected ErrTmuxNotFound, got nil")
	}

	if !errors.Is(err, tmux.ErrTmuxNotFound) {
		t.Errorf("expected errors.Is(err, tmux.ErrTmuxNotFound); got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// C-110: BuildAttachCmd
// ---------------------------------------------------------------------------

// TestBuildAttachCmd_IncludesTargetSession verifies that the built command's
// Args contain "-t" followed immediately by the session name.
func TestBuildAttachCmd_IncludesTargetSession(t *testing.T) {
	overrideLookPath(t, tmuxFoundLookPath)

	const session = "my-session"
	cmd, err := tmux.BuildAttachCmd(session)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	got, found := argAfter(cmd.Args, "-t")
	if !found {
		t.Fatalf("cmd.Args does not contain \"-t\" flag; got: %v", cmd.Args)
	}
	if got != session {
		t.Errorf("arg after \"-t\" = %q; want %q", got, session)
	}
}

// TestBuildAttachCmd_UsesAttachSession verifies that the second element of
// cmd.Args (Args[1]) is "attach-session", which is the tmux subcommand.
func TestBuildAttachCmd_UsesAttachSession(t *testing.T) {
	overrideLookPath(t, tmuxFoundLookPath)

	cmd, err := tmux.BuildAttachCmd("my-session")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(cmd.Args) < 2 {
		t.Fatalf("cmd.Args has fewer than 2 elements; got: %v", cmd.Args)
	}
	if cmd.Args[1] != "attach-session" {
		t.Errorf("cmd.Args[1] = %q; want \"attach-session\"", cmd.Args[1])
	}
}

// TestBuildAttachCmd_TmuxNotInPath_ReturnsError verifies that BuildAttachCmd
// returns ErrTmuxNotFound when tmux is absent from PATH.
func TestBuildAttachCmd_TmuxNotInPath_ReturnsError(t *testing.T) {
	overrideLookPath(t, tmuxNotFoundLookPath)

	_, err := tmux.BuildAttachCmd("my-session")
	if err == nil {
		t.Fatal("expected ErrTmuxNotFound, got nil")
	}

	if !errors.Is(err, tmux.ErrTmuxNotFound) {
		t.Errorf("expected errors.Is(err, tmux.ErrTmuxNotFound); got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Sentinel errors
// ---------------------------------------------------------------------------

// TestSentinelErrors_AreDefined verifies that all sentinel errors exported by
// the tmux package are non-nil and are distinct from one another, so callers
// can reliably distinguish error cases with errors.Is.
func TestSentinelErrors_AreDefined(t *testing.T) {
	sentinels := []struct {
		name  string
		value error
	}{
		{"ErrTmuxNotFound", tmux.ErrTmuxNotFound},
		{"ErrSessionExists", tmux.ErrSessionExists},
		{"ErrCreateFailed", tmux.ErrCreateFailed},
		{"ErrSessionNotFound", tmux.ErrSessionNotFound},
		{"ErrAddWindowFailed", tmux.ErrAddWindowFailed},
		{"ErrAttachFailed", tmux.ErrAttachFailed},
	}

	// All sentinels must be non-nil.
	for _, sentinel := range sentinels {
		if sentinel.value == nil {
			t.Errorf("tmux.%s must not be nil", sentinel.name)
		}
	}

	// All sentinels must be distinct from each other.
	for index1, sentinel1 := range sentinels {
		for index2, sentinel2 := range sentinels {
			if index1 == index2 {
				continue
			}
			if errors.Is(sentinel1.value, sentinel2.value) {
				t.Errorf("tmux.%s and tmux.%s must be distinct sentinel errors", sentinel1.name, sentinel2.name)
			}
		}
	}
}
