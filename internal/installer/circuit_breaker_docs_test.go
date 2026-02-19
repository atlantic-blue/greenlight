package installer_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// C-54 Tests: CircuitBreakerProtocol
// These tests verify that src/references/circuit-breaker.md exists and contains
// all required sections as defined in contract C-54.

func TestCircuitBreakerMd_Exists(t *testing.T) {
	path := filepath.Join(projectRoot(), "src/references/circuit-breaker.md")
	_, err := os.Stat(path)
	if err != nil {
		t.Fatalf("circuit-breaker.md does not exist at %s: %v", path, err)
	}
}

func TestCircuitBreakerMd_ContainsAttemptTrackingSection(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/circuit-breaker.md"))
	if err != nil {
		t.Fatalf("failed to read circuit-breaker.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "attempt_count") {
		t.Error("circuit-breaker.md missing attempt_count field in attempt tracking state")
	}
	if !strings.Contains(doc, "slice_id") {
		t.Error("circuit-breaker.md missing slice_id field in attempt tracking state")
	}
	if !strings.Contains(doc, "test_name") {
		t.Error("circuit-breaker.md missing test_name field in attempt tracking state")
	}
	if !strings.Contains(doc, "files_touched") {
		t.Error("circuit-breaker.md missing files_touched field in attempt tracking state")
	}
	if !strings.Contains(doc, "last_error") {
		t.Error("circuit-breaker.md missing last_error field in attempt tracking state")
	}
	if !strings.Contains(doc, "checkpoint_tag") {
		t.Error("circuit-breaker.md missing checkpoint_tag field in attempt tracking state")
	}
	if !strings.Contains(doc, "total_failed_attempts") {
		t.Error("circuit-breaker.md missing total_failed_attempts slice-level accumulator")
	}
}

func TestCircuitBreakerMd_ContainsPerTestTripThreshold(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/circuit-breaker.md"))
	if err != nil {
		t.Fatalf("failed to read circuit-breaker.md: %v", err)
	}

	doc := string(content)

	// The per-test threshold is 3 — this is mandatory and not configurable
	if !strings.Contains(doc, "3") {
		t.Error("circuit-breaker.md missing per-test trip threshold of 3")
	}

	// Must describe what counts as an attempt: a complete cycle
	if !strings.Contains(doc, "attempt") {
		t.Error("circuit-breaker.md missing description of what constitutes an attempt")
	}

	// Counter increments only on test failure, not infrastructure errors
	if !strings.Contains(doc, "infrastructure") {
		t.Error("circuit-breaker.md missing distinction between test failures and infrastructure errors")
	}
}

func TestCircuitBreakerMd_ContainsSliceLevelCeiling(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/circuit-breaker.md"))
	if err != nil {
		t.Fatalf("failed to read circuit-breaker.md: %v", err)
	}

	doc := string(content)

	// The slice-level ceiling is 7
	if !strings.Contains(doc, "7") {
		t.Error("circuit-breaker.md missing slice-level ceiling value of 7")
	}

	if !strings.Contains(doc, "total_failed_attempts") {
		t.Error("circuit-breaker.md missing total_failed_attempts for slice-level ceiling check")
	}
}

func TestCircuitBreakerMd_ContainsDiagnosticReportFormat(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/circuit-breaker.md"))
	if err != nil {
		t.Fatalf("failed to read circuit-breaker.md: %v", err)
	}

	doc := string(content)

	// All 8 required fields from the Structured Diagnostic Report Format (C-54)
	requiredFields := []string{
		"test_expectation",
		"actual_error",
		"attempt_log",
		"cumulative_files_modified",
		"scope_violations",
		"best_hypothesis",
		"specific_question",
		"recovery_options",
	}

	for _, field := range requiredFields {
		if !strings.Contains(doc, field) {
			t.Errorf("circuit-breaker.md diagnostic report missing required field: %s", field)
		}
	}
}

func TestCircuitBreakerMd_ContainsScopeLockProtocol(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/circuit-breaker.md"))
	if err != nil {
		t.Fatalf("failed to read circuit-breaker.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "scope") {
		t.Error("circuit-breaker.md missing scope lock protocol")
	}

	// Must reference GRAPH.json as part of scope inference
	if !strings.Contains(doc, "GRAPH.json") {
		t.Error("circuit-breaker.md missing GRAPH.json reference in scope lock protocol")
	}

	// Unjustifiable out-of-scope modification counts as a failed attempt
	if !strings.Contains(doc, "justification") {
		t.Error("circuit-breaker.md missing justification requirement for out-of-scope modifications")
	}
}

func TestCircuitBreakerMd_ContainsCounterResetProtocol(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/circuit-breaker.md"))
	if err != nil {
		t.Fatalf("failed to read circuit-breaker.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "reset") {
		t.Error("circuit-breaker.md missing counter reset protocol")
	}

	// Must mention rollback to checkpoint
	if !strings.Contains(doc, "rollback") {
		t.Error("circuit-breaker.md missing rollback instruction in counter reset protocol")
	}

	// Fresh implementer with clean state
	if !strings.Contains(doc, "what was tried") {
		t.Error("circuit-breaker.md missing 'what was tried' context handoff in reset protocol")
	}
}

func TestCircuitBreakerMd_ContainsDeviationRulesAdditive(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/circuit-breaker.md"))
	if err != nil {
		t.Fatalf("failed to read circuit-breaker.md: %v", err)
	}

	doc := string(content)

	// Must be additive to deviation rules, not a replacement
	if !strings.Contains(doc, "deviation") {
		t.Error("circuit-breaker.md missing reference to deviation rules")
	}

	// Deviation Rule 4 (ARCH-STOP) takes priority
	if !strings.Contains(doc, "Rule 4") {
		t.Error("circuit-breaker.md missing explicit statement that Deviation Rule 4 takes priority")
	}

	// Circuit breaker handles unproductive work, not unplanned work
	if !strings.Contains(doc, "additive") {
		t.Error("circuit-breaker.md missing statement that the protocol is additive to deviation rules")
	}
}

// C-55 Tests: ImplementerCircuitBreaker
// These tests verify that src/agents/gl-implementer.md has been updated with
// the circuit breaker integration as defined in contract C-55.

func TestImplementerMd_ReferencesCircuitBreaker(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-implementer.md"))
	if err != nil {
		t.Fatalf("failed to read gl-implementer.md: %v", err)
	}

	doc := string(content)

	if !strings.Contains(doc, "circuit-breaker.md") {
		t.Error("gl-implementer.md must reference references/circuit-breaker.md")
	}
}

func TestImplementerMd_ContainsScopeCheckBeforeModification(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-implementer.md"))
	if err != nil {
		t.Fatalf("failed to read gl-implementer.md: %v", err)
	}

	doc := string(content)

	// Scope check must happen before every file modification
	if !strings.Contains(doc, "scope") {
		t.Error("gl-implementer.md missing scope check instruction")
	}

	// Must mention the check happens before file modification
	if !strings.Contains(doc, "Before") && !strings.Contains(doc, "before") {
		t.Error("gl-implementer.md missing instruction that scope check occurs before file modification")
	}
}

func TestImplementerMd_ContainsAttemptTracking(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-implementer.md"))
	if err != nil {
		t.Fatalf("failed to read gl-implementer.md: %v", err)
	}

	doc := string(content)

	// Must instruct implementer to maintain structured attempt state with the
	// specific field names defined in the circuit breaker protocol (C-54/C-55)
	if !strings.Contains(doc, "attempt_count") {
		t.Error("gl-implementer.md missing attempt_count state field — must track per-test attempt counts")
	}

	// Must reference the slice-level accumulator
	if !strings.Contains(doc, "total_failed_attempts") {
		t.Error("gl-implementer.md missing total_failed_attempts — must track slice-level attempt ceiling")
	}
}

func TestImplementerMd_ContainsDiagnosticReportProduction(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-implementer.md"))
	if err != nil {
		t.Fatalf("failed to read gl-implementer.md: %v", err)
	}

	doc := string(content)

	// Must instruct implementer to produce a diagnostic report on circuit trip
	if !strings.Contains(doc, "Diagnostic") && !strings.Contains(doc, "diagnostic") {
		t.Error("gl-implementer.md missing instruction to produce diagnostic report on circuit trip")
	}

	// Must instruct implementer to STOP after producing the report
	if !strings.Contains(doc, "STOP") {
		t.Error("gl-implementer.md missing STOP instruction when circuit trips")
	}
}

func TestImplementerMd_ContainsInfrastructureErrorDistinction(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-implementer.md"))
	if err != nil {
		t.Fatalf("failed to read gl-implementer.md: %v", err)
	}

	doc := string(content)

	// Must distinguish infrastructure errors from test failures
	if !strings.Contains(doc, "infrastructure") && !strings.Contains(doc, "Infrastructure") {
		t.Error("gl-implementer.md missing distinction between infrastructure errors and test failures")
	}
}

func TestImplementerMd_DoesNotContainOldErrorRecovery(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/agents/gl-implementer.md"))
	if err != nil {
		t.Fatalf("failed to read gl-implementer.md: %v", err)
	}

	doc := string(content)

	// The old "Step 5: Know When to Stop" section must be replaced by the circuit breaker protocol
	if strings.Contains(doc, "Step 5: Know When to Stop") {
		t.Error("gl-implementer.md still contains old error recovery section 'Step 5: Know When to Stop' — must be replaced by circuit breaker integration")
	}
}

// C-56 Tests: ScopeLockProtocol
// These tests verify that circuit-breaker.md contains the complete scope lock
// protocol with inference rules, justification format, and files_in_scope override.

func TestCircuitBreakerMd_ScopeLockHasInferenceRules(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/circuit-breaker.md"))
	if err != nil {
		t.Fatalf("failed to read circuit-breaker.md: %v", err)
	}

	doc := string(content)

	// Scope inference rules in priority order per C-56:
	// 1. Explicit override: files_in_scope
	// 2. Inferred from contracts: package names, file paths, GRAPH.json packages/deliverables
	// 3. Fallback: slice's packages/deliverables from GRAPH.json

	if !strings.Contains(doc, "infer") && !strings.Contains(doc, "Infer") {
		t.Error("circuit-breaker.md missing scope inference rules")
	}

	if !strings.Contains(doc, "contracts") && !strings.Contains(doc, "contract") {
		t.Error("circuit-breaker.md missing contracts as a source for scope inference")
	}

	if !strings.Contains(doc, "packages") && !strings.Contains(doc, "deliverables") {
		t.Error("circuit-breaker.md missing packages/deliverables as fallback scope inference source")
	}
}

func TestCircuitBreakerMd_ScopeLockHasJustificationFormat(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/circuit-breaker.md"))
	if err != nil {
		t.Fatalf("failed to read circuit-breaker.md: %v", err)
	}

	doc := string(content)

	// Justification format must include all four required elements per C-56:
	// file path, failing test name, reason, relationship to contract boundary

	if !strings.Contains(doc, "justification") && !strings.Contains(doc, "Justification") {
		t.Error("circuit-breaker.md missing justification format for out-of-scope modifications")
	}

	if !strings.Contains(doc, "failing test") && !strings.Contains(doc, "Failing test") {
		t.Error("circuit-breaker.md justification format missing failing test name requirement")
	}

	if !strings.Contains(doc, "Reason") || !strings.Contains(doc, "Relationship") {
		t.Error("circuit-breaker.md justification format missing Reason and Relationship fields")
	}
}

func TestCircuitBreakerMd_ScopeLockHasFilesInScopeOverride(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(projectRoot(), "src/references/circuit-breaker.md"))
	if err != nil {
		t.Fatalf("failed to read circuit-breaker.md: %v", err)
	}

	doc := string(content)

	// files_in_scope is the optional explicit override with highest priority per C-56
	if !strings.Contains(doc, "files_in_scope") {
		t.Error("circuit-breaker.md missing files_in_scope optional override field")
	}
}
