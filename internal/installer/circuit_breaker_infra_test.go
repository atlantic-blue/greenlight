package installer_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/atlantic-blue/greenlight/internal/installer"
)

// C-60 Tests: CLAUDEmdCircuitBreakerRule
// These tests verify that src/CLAUDE.md contains the circuit breaker
// hard rule subsection as defined in contract C-60.

func TestCLAUDEmd_ContainsCircuitBreakerSubsection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/CLAUDE.md"))
	if err != nil {
		t.Fatalf("failed to read CLAUDE.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "### Circuit Breaker") {
		t.Error("CLAUDE.md missing '### Circuit Breaker' subsection heading")
	}
}

func TestCLAUDEmd_CircuitBreakerAfterTestingBeforeLogging(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/CLAUDE.md"))
	if err != nil {
		t.Fatalf("failed to read CLAUDE.md: %v", err)
	}

	doc := string(content)

	testingPos := strings.Index(doc, "### Testing")
	circuitBreakerPos := strings.Index(doc, "### Circuit Breaker")
	loggingPos := strings.Index(doc, "### Logging & Observability")

	if testingPos == -1 {
		t.Fatal("CLAUDE.md missing '### Testing' section")
	}
	if circuitBreakerPos == -1 {
		t.Fatal("CLAUDE.md missing '### Circuit Breaker' section")
	}
	if loggingPos == -1 {
		t.Fatal("CLAUDE.md missing '### Logging & Observability' section")
	}

	if circuitBreakerPos <= testingPos {
		t.Errorf("'### Circuit Breaker' (pos %d) must appear AFTER '### Testing' (pos %d)", circuitBreakerPos, testingPos)
	}
	if circuitBreakerPos >= loggingPos {
		t.Errorf("'### Circuit Breaker' (pos %d) must appear BEFORE '### Logging & Observability' (pos %d)", circuitBreakerPos, loggingPos)
	}
}

func TestCLAUDEmd_CircuitBreakerContainsThresholds(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/CLAUDE.md"))
	if err != nil {
		t.Fatalf("failed to read CLAUDE.md: %v", err)
	}

	doc := string(content)

	// Find the circuit breaker section
	cbStart := strings.Index(doc, "### Circuit Breaker")
	if cbStart == -1 {
		t.Fatal("CLAUDE.md missing '### Circuit Breaker' section")
	}

	// Find the next section heading after Circuit Breaker
	nextSection := strings.Index(doc[cbStart+1:], "### ")
	var cbSection string
	if nextSection == -1 {
		cbSection = doc[cbStart:]
	} else {
		cbSection = doc[cbStart : cbStart+1+nextSection]
	}

	// Must state the per-test threshold of 3
	if !strings.Contains(cbSection, "3") {
		t.Error("CLAUDE.md Circuit Breaker section missing per-test threshold of 3")
	}

	// Must state the slice-level ceiling of 7
	if !strings.Contains(cbSection, "7") {
		t.Error("CLAUDE.md Circuit Breaker section missing slice-level ceiling of 7")
	}

	// Must reference scope verification
	if !strings.Contains(cbSection, "scope") {
		t.Error("CLAUDE.md Circuit Breaker section missing scope verification rule")
	}

	// Must reference the full protocol document
	if !strings.Contains(cbSection, "circuit-breaker.md") {
		t.Error("CLAUDE.md Circuit Breaker section missing reference to references/circuit-breaker.md")
	}
}

func TestCLAUDEmd_CircuitBreakerUsesImperatives(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/CLAUDE.md"))
	if err != nil {
		t.Fatalf("failed to read CLAUDE.md: %v", err)
	}

	doc := string(content)

	cbStart := strings.Index(doc, "### Circuit Breaker")
	if cbStart == -1 {
		t.Fatal("CLAUDE.md missing '### Circuit Breaker' section")
	}

	nextSection := strings.Index(doc[cbStart+1:], "### ")
	var cbSection string
	if nextSection == -1 {
		cbSection = doc[cbStart:]
	} else {
		cbSection = doc[cbStart : cbStart+1+nextSection]
	}

	// Must use imperative language — STOP is a hard rule
	if !strings.Contains(cbSection, "STOP") {
		t.Error("CLAUDE.md Circuit Breaker section missing imperative 'STOP' — must be a hard rule, not a recommendation")
	}
}

// C-61 Tests: ManifestCircuitBreakerUpdate
// These tests verify that the Go manifest includes the 2 new circuit breaker
// file paths and the total count is updated as defined in contract C-61.

func TestManifest_ContainsDebugCommand(t *testing.T) {
	found := false
	for _, entry := range installer.Manifest {
		if entry == "commands/gl/debug.md" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Manifest missing 'commands/gl/debug.md' entry")
	}
}

func TestManifest_ContainsCircuitBreakerReference(t *testing.T) {
	found := false
	for _, entry := range installer.Manifest {
		if entry == "references/circuit-breaker.md" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Manifest missing 'references/circuit-breaker.md' entry")
	}
}

func TestManifest_Contains34Entries(t *testing.T) {
	if len(installer.Manifest) != 38 {
		t.Errorf("expected 38 manifest entries, got %d", len(installer.Manifest))
	}
}

func TestManifest_CLAUDEmdIsLastEntry(t *testing.T) {
	if len(installer.Manifest) == 0 {
		t.Fatal("Manifest is empty")
	}

	lastEntry := installer.Manifest[len(installer.Manifest)-1]
	if lastEntry != "CLAUDE.md" {
		t.Errorf("CLAUDE.md must be the last manifest entry, got %q", lastEntry)
	}
}
