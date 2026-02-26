package process_test

// S-41 Tests: Process Spawner
// Covers C-103 (ProcessSpawnClaude) and C-104 (ProcessSpawnInteractive).
//
// Contract C-103 — SpawnClaude / BuildClaudeCmd:
//   - Verifies "claude" is in PATH before attempting to start
//   - Builds command: claude -p "{prompt}" {flags...}
//   - Sets working directory to opts.Dir
//   - Connects cmd.Stdout and cmd.Stderr to provided writers
//   - Returns ErrEmptyPrompt when prompt is empty
//   - Returns ErrClaudeNotFound when "claude" is not in PATH
//   - Process is started but not waited on (SpawnClaude)
//
// Contract C-104 — SpawnInteractive / BuildInteractiveCmd:
//   - Verifies "claude" is in PATH before attempting to start
//   - Builds command: claude {flags...}
//   - Adds -p "{prompt}" only when prompt is non-empty
//   - NEVER adds --dangerously-skip-permissions
//   - Returns ErrClaudeNotFound when "claude" is not in PATH

import (
	"bytes"
	"errors"
	"os/exec"
	"testing"

	"github.com/atlantic-blue/greenlight/internal/process"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// claudeFoundLookPath returns a LookPath stub that reports "claude" as found.
func claudeFoundLookPath(name string) (string, error) {
	if name == "claude" {
		return "/usr/local/bin/claude", nil
	}
	return "", exec.ErrNotFound
}

// claudeNotFoundLookPath returns a LookPath stub that always reports not found.
func claudeNotFoundLookPath(_ string) (string, error) {
	return "", exec.ErrNotFound
}

// overrideLookPath replaces process.LookPath with the given stub and registers
// a cleanup to restore the original value when the test finishes.
func overrideLookPath(t *testing.T, stub func(string) (string, error)) {
	t.Helper()
	original := process.LookPath
	process.LookPath = stub
	t.Cleanup(func() { process.LookPath = original })
}

// ---------------------------------------------------------------------------
// C-103: BuildClaudeCmd
// ---------------------------------------------------------------------------

// TestBuildClaudeCmd_IncludesPromptFlag verifies that the built command's Args
// contain the "-p" flag immediately followed by the provided prompt text.
func TestBuildClaudeCmd_IncludesPromptFlag(t *testing.T) {
	overrideLookPath(t, claudeFoundLookPath)

	const prompt = "write a hello world program"
	cmd, err := process.BuildClaudeCmd(process.SpawnOptions{
		Prompt: prompt,
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	foundFlag := false
	foundPrompt := false
	for i, arg := range cmd.Args {
		if arg == "-p" {
			foundFlag = true
			if i+1 < len(cmd.Args) && cmd.Args[i+1] == prompt {
				foundPrompt = true
			}
		}
	}

	if !foundFlag {
		t.Errorf("cmd.Args does not contain \"-p\" flag; got: %v", cmd.Args)
	}
	if !foundPrompt {
		t.Errorf("cmd.Args does not contain prompt %q after \"-p\" flag; got: %v", prompt, cmd.Args)
	}
}

// TestBuildClaudeCmd_IncludesAdditionalFlags verifies that extra flags provided
// via SpawnOptions.Flags are present in the built command's Args.
func TestBuildClaudeCmd_IncludesAdditionalFlags(t *testing.T) {
	overrideLookPath(t, claudeFoundLookPath)

	extraFlags := []string{"--output-format", "json", "--verbose"}
	cmd, err := process.BuildClaudeCmd(process.SpawnOptions{
		Prompt: "list files",
		Flags:  extraFlags,
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	for _, flag := range extraFlags {
		found := false
		for _, arg := range cmd.Args {
			if arg == flag {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("cmd.Args does not contain extra flag %q; got: %v", flag, cmd.Args)
		}
	}
}

// TestBuildClaudeCmd_SetsWorkingDirectory verifies that cmd.Dir is set to
// the value provided in SpawnOptions.Dir.
func TestBuildClaudeCmd_SetsWorkingDirectory(t *testing.T) {
	overrideLookPath(t, claudeFoundLookPath)

	const workDir = "/tmp/myproject"
	cmd, err := process.BuildClaudeCmd(process.SpawnOptions{
		Prompt: "analyse code",
		Dir:    workDir,
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cmd.Dir != workDir {
		t.Errorf("cmd.Dir = %q; want %q", cmd.Dir, workDir)
	}
}

// TestBuildClaudeCmd_EmptyDirDefaultsToEmpty verifies that when SpawnOptions.Dir
// is an empty string, cmd.Dir is also an empty string (inheriting the caller's
// working directory at process start time, per os/exec semantics).
func TestBuildClaudeCmd_EmptyDirDefaultsToEmpty(t *testing.T) {
	overrideLookPath(t, claudeFoundLookPath)

	cmd, err := process.BuildClaudeCmd(process.SpawnOptions{
		Prompt: "hello",
		Dir:    "",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cmd.Dir != "" {
		t.Errorf("cmd.Dir = %q; want empty string when opts.Dir is empty", cmd.Dir)
	}
}

// TestBuildClaudeCmd_ConnectsStdoutStderr verifies that cmd.Stdout and cmd.Stderr
// are set to the writers provided in SpawnOptions.
func TestBuildClaudeCmd_ConnectsStdoutStderr(t *testing.T) {
	overrideLookPath(t, claudeFoundLookPath)

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	cmd, err := process.BuildClaudeCmd(process.SpawnOptions{
		Prompt: "run something",
		Stdout: &stdoutBuf,
		Stderr: &stderrBuf,
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cmd.Stdout != &stdoutBuf {
		t.Error("cmd.Stdout was not set to the provided stdout writer")
	}
	if cmd.Stderr != &stderrBuf {
		t.Error("cmd.Stderr was not set to the provided stderr writer")
	}
}

// TestBuildClaudeCmd_EmptyPromptReturnsError verifies that BuildClaudeCmd
// returns ErrEmptyPrompt when SpawnOptions.Prompt is an empty string.
func TestBuildClaudeCmd_EmptyPromptReturnsError(t *testing.T) {
	overrideLookPath(t, claudeFoundLookPath)

	_, err := process.BuildClaudeCmd(process.SpawnOptions{
		Prompt: "",
	})
	if err == nil {
		t.Fatal("expected ErrEmptyPrompt, got nil")
	}

	if !errors.Is(err, process.ErrEmptyPrompt) {
		t.Errorf("expected errors.Is(err, process.ErrEmptyPrompt); got: %v", err)
	}
}

// TestBuildClaudeCmd_ClaudeNotInPath_ReturnsError verifies that BuildClaudeCmd
// returns ErrClaudeNotFound when the "claude" binary is not present in PATH.
func TestBuildClaudeCmd_ClaudeNotInPath_ReturnsError(t *testing.T) {
	overrideLookPath(t, claudeNotFoundLookPath)

	_, err := process.BuildClaudeCmd(process.SpawnOptions{
		Prompt: "some prompt",
	})
	if err == nil {
		t.Fatal("expected ErrClaudeNotFound, got nil")
	}

	if !errors.Is(err, process.ErrClaudeNotFound) {
		t.Errorf("expected errors.Is(err, process.ErrClaudeNotFound); got: %v", err)
	}
}

// TestSpawnClaude_ClaudeNotInPath_ReturnsError verifies that SpawnClaude
// returns ErrClaudeNotFound (propagated from the PATH check) before ever
// attempting to start a process when "claude" is absent from PATH.
func TestSpawnClaude_ClaudeNotInPath_ReturnsError(t *testing.T) {
	overrideLookPath(t, claudeNotFoundLookPath)

	_, err := process.SpawnClaude(process.SpawnOptions{
		Prompt: "some prompt",
	})
	if err == nil {
		t.Fatal("expected ErrClaudeNotFound, got nil")
	}

	if !errors.Is(err, process.ErrClaudeNotFound) {
		t.Errorf("expected errors.Is(err, process.ErrClaudeNotFound); got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// C-104: BuildInteractiveCmd
// ---------------------------------------------------------------------------

// TestBuildInteractiveCmd_NoPrompt_NoPromptFlag verifies that when
// InteractiveOptions.Prompt is empty, cmd.Args does not include "-p".
func TestBuildInteractiveCmd_NoPrompt_NoPromptFlag(t *testing.T) {
	overrideLookPath(t, claudeFoundLookPath)

	cmd, err := process.BuildInteractiveCmd(process.InteractiveOptions{
		Prompt: "",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	for _, arg := range cmd.Args {
		if arg == "-p" {
			t.Errorf("cmd.Args must not contain \"-p\" when prompt is empty; got: %v", cmd.Args)
			break
		}
	}
}

// TestBuildInteractiveCmd_WithPrompt_IncludesPromptFlag verifies that when
// InteractiveOptions.Prompt is non-empty, cmd.Args includes "-p" followed by
// the prompt text.
func TestBuildInteractiveCmd_WithPrompt_IncludesPromptFlag(t *testing.T) {
	overrideLookPath(t, claudeFoundLookPath)

	const prompt = "start the design session"
	cmd, err := process.BuildInteractiveCmd(process.InteractiveOptions{
		Prompt: prompt,
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	foundFlag := false
	foundPrompt := false
	for i, arg := range cmd.Args {
		if arg == "-p" {
			foundFlag = true
			if i+1 < len(cmd.Args) && cmd.Args[i+1] == prompt {
				foundPrompt = true
			}
		}
	}

	if !foundFlag {
		t.Errorf("cmd.Args does not contain \"-p\" flag; got: %v", cmd.Args)
	}
	if !foundPrompt {
		t.Errorf("cmd.Args does not contain prompt %q after \"-p\" flag; got: %v", prompt, cmd.Args)
	}
}

// TestBuildInteractiveCmd_NeverIncludesDangerousFlag verifies the invariant that
// --dangerously-skip-permissions is never present in the interactive command's
// Args, regardless of any provided flags or prompt.
func TestBuildInteractiveCmd_NeverIncludesDangerousFlag(t *testing.T) {
	overrideLookPath(t, claudeFoundLookPath)

	tests := []struct {
		name string
		opts process.InteractiveOptions
	}{
		{
			name: "no prompt, no extra flags",
			opts: process.InteractiveOptions{},
		},
		{
			name: "with prompt",
			opts: process.InteractiveOptions{Prompt: "hello"},
		},
		{
			name: "caller attempts to inject dangerous flag",
			opts: process.InteractiveOptions{
				Flags: []string{"--dangerously-skip-permissions"},
			},
		},
		{
			name: "caller attempts to inject dangerous flag with prompt",
			opts: process.InteractiveOptions{
				Prompt: "hello",
				Flags:  []string{"--dangerously-skip-permissions"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := process.BuildInteractiveCmd(tt.opts)
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}

			for _, arg := range cmd.Args {
				if arg == "--dangerously-skip-permissions" {
					t.Errorf("cmd.Args must NEVER contain \"--dangerously-skip-permissions\"; got: %v", cmd.Args)
					break
				}
			}
		})
	}
}

// TestBuildInteractiveCmd_SetsWorkingDirectory verifies that cmd.Dir is set to
// the value provided in InteractiveOptions.Dir.
func TestBuildInteractiveCmd_SetsWorkingDirectory(t *testing.T) {
	overrideLookPath(t, claudeFoundLookPath)

	const workDir = "/home/user/project"
	cmd, err := process.BuildInteractiveCmd(process.InteractiveOptions{
		Dir: workDir,
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cmd.Dir != workDir {
		t.Errorf("cmd.Dir = %q; want %q", cmd.Dir, workDir)
	}
}

// TestBuildInteractiveCmd_IncludesAdditionalFlags verifies that extra flags
// provided via InteractiveOptions.Flags appear in the built command's Args
// (unless they are --dangerously-skip-permissions, which is always stripped).
func TestBuildInteractiveCmd_IncludesAdditionalFlags(t *testing.T) {
	overrideLookPath(t, claudeFoundLookPath)

	extraFlags := []string{"--model", "claude-opus-4-6"}
	cmd, err := process.BuildInteractiveCmd(process.InteractiveOptions{
		Flags: extraFlags,
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	for _, flag := range extraFlags {
		found := false
		for _, arg := range cmd.Args {
			if arg == flag {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("cmd.Args does not contain extra flag %q; got: %v", flag, cmd.Args)
		}
	}
}

// TestBuildInteractiveCmd_ClaudeNotInPath_ReturnsError verifies that
// BuildInteractiveCmd returns ErrClaudeNotFound when "claude" is absent
// from PATH.
func TestBuildInteractiveCmd_ClaudeNotInPath_ReturnsError(t *testing.T) {
	overrideLookPath(t, claudeNotFoundLookPath)

	_, err := process.BuildInteractiveCmd(process.InteractiveOptions{})
	if err == nil {
		t.Fatal("expected ErrClaudeNotFound, got nil")
	}

	if !errors.Is(err, process.ErrClaudeNotFound) {
		t.Errorf("expected errors.Is(err, process.ErrClaudeNotFound); got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Sentinel errors
// ---------------------------------------------------------------------------

// TestSentinelErrors_AreDefined verifies that ErrClaudeNotFound, ErrEmptyPrompt,
// and ErrStartFailure are all non-nil and are distinct from one another, so
// callers can distinguish error cases with errors.Is.
func TestSentinelErrors_AreDefined(t *testing.T) {
	sentinels := []struct {
		name  string
		value error
	}{
		{"ErrClaudeNotFound", process.ErrClaudeNotFound},
		{"ErrEmptyPrompt", process.ErrEmptyPrompt},
		{"ErrStartFailure", process.ErrStartFailure},
	}

	// All sentinels must be non-nil.
	for _, s := range sentinels {
		if s.value == nil {
			t.Errorf("process.%s must not be nil", s.name)
		}
	}

	// All sentinels must be distinct from each other.
	for i, s1 := range sentinels {
		for j, s2 := range sentinels {
			if i == j {
				continue
			}
			if errors.Is(s1.value, s2.value) {
				t.Errorf("process.%s and process.%s must be distinct sentinel errors", s1.name, s2.name)
			}
		}
	}
}
