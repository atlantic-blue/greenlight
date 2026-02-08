package cmd_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/atlantic-blue/greenlight/internal/cmd"
)

// TestRunVersion_ReturnsZero verifies that RunVersion always
// returns exit code 0.
func TestRunVersion_ReturnsZero(t *testing.T) {
	var buf bytes.Buffer
	exitCode := cmd.RunVersion(&buf)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
}

// TestRunVersion_OutputFormat verifies that RunVersion writes
// the correct format to the provided io.Writer.
func TestRunVersion_OutputFormat(t *testing.T) {
	var buf bytes.Buffer
	cmd.RunVersion(&buf)

	output := buf.String()

	// Expected format: "greenlight <Version> (commit: <GitCommit>, built: <BuildDate>)\n"
	if !strings.HasPrefix(output, "greenlight ") {
		t.Errorf("output should start with 'greenlight ', got: %q", output)
	}

	if !strings.Contains(output, "(commit: ") {
		t.Errorf("output should contain '(commit: ', got: %q", output)
	}

	if !strings.Contains(output, ", built: ") {
		t.Errorf("output should contain ', built: ', got: %q", output)
	}

	if !strings.HasSuffix(output, ")\n") {
		t.Errorf("output should end with ')\\n', got: %q", output)
	}
}

// TestRunVersion_WritesToProvidedWriter verifies that RunVersion
// writes to the provided io.Writer and not elsewhere.
func TestRunVersion_WritesToProvidedWriter(t *testing.T) {
	var buf bytes.Buffer
	cmd.RunVersion(&buf)

	if buf.Len() == 0 {
		t.Error("expected output to be written to provided writer, got empty buffer")
	}

	output := buf.String()
	if output == "" {
		t.Error("expected non-empty output")
	}
}

// TestRunVersion_OutputContainsVersionInfo verifies that the output
// includes all three version components.
func TestRunVersion_OutputContainsVersionInfo(t *testing.T) {
	var buf bytes.Buffer
	cmd.RunVersion(&buf)

	output := buf.String()

	// The output must contain the version pattern
	// Format: "greenlight <Version> (commit: <GitCommit>, built: <BuildDate>)\n"
	parts := strings.Split(output, " ")

	if len(parts) < 5 {
		t.Errorf("expected at least 5 parts in output, got %d: %q", len(parts), output)
	}

	// Verify structure
	if parts[0] != "greenlight" {
		t.Errorf("first part should be 'greenlight', got: %q", parts[0])
	}

	// Second part should be the version (non-empty)
	if parts[1] == "" {
		t.Error("version part should not be empty")
	}

	// Remaining parts should contain commit and built info
	hasCommit := false
	hasBuilt := false

	for _, part := range parts {
		if strings.Contains(part, "commit:") {
			hasCommit = true
		}
		if strings.Contains(part, "built:") {
			hasBuilt = true
		}
	}

	if !hasCommit {
		t.Errorf("output should contain 'commit:', got: %q", output)
	}

	if !hasBuilt {
		t.Errorf("output should contain 'built:', got: %q", output)
	}
}

// TestRunVersion_MultipleCallsProduceSameOutput verifies that
// RunVersion is idempotent and produces consistent output.
func TestRunVersion_MultipleCallsProduceSameOutput(t *testing.T) {
	var buf1 bytes.Buffer
	var buf2 bytes.Buffer

	exitCode1 := cmd.RunVersion(&buf1)
	exitCode2 := cmd.RunVersion(&buf2)

	if exitCode1 != exitCode2 {
		t.Errorf("exit codes differ: %d vs %d", exitCode1, exitCode2)
	}

	output1 := buf1.String()
	output2 := buf2.String()

	if output1 != output2 {
		t.Errorf("outputs differ:\n  first:  %q\n  second: %q", output1, output2)
	}
}
