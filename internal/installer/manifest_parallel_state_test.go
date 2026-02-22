package installer_test

// S-34 Tests: Manifest and Integration
// Covers C-90 (ManifestParallelStateUpdate).
//
// Contract: 3 new entries added to installer.Manifest, growing it from 35 to 38.
// New entries:
//   "commands/gl/migrate-state.md"
//   "references/state-format.md"
//   "templates/slice-state.md"
//
// Invariants:
//   - CLAUDE.md remains the last entry
//   - All sections are alphabetically ordered
//   - All 35 previously-existing entries are retained

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/atlantic-blue/greenlight/internal/installer"
)

// TestManifest_C90_CountIs38 verifies that the manifest grows from 35 to 38 entries
// after the 3 new parallel-state files are added.
func TestManifest_C90_CountIs38(t *testing.T) {
	got := len(installer.Manifest)
	if got != 38 {
		t.Errorf("installer.Manifest must have 38 entries (was 35, +3 new), got %d", got)
	}
}

// TestManifest_C90_ContainsMigrateStateMd verifies the new migration command file
// is present in the manifest.
func TestManifest_C90_ContainsMigrateStateMd(t *testing.T) {
	for _, entry := range installer.Manifest {
		if entry == "commands/gl/migrate-state.md" {
			return
		}
	}
	t.Error("installer.Manifest missing new entry: \"commands/gl/migrate-state.md\"")
}

// TestManifest_C90_ContainsStateFormatMd verifies the new state-format reference file
// is present in the manifest.
func TestManifest_C90_ContainsStateFormatMd(t *testing.T) {
	for _, entry := range installer.Manifest {
		if entry == "references/state-format.md" {
			return
		}
	}
	t.Error("installer.Manifest missing new entry: \"references/state-format.md\"")
}

// TestManifest_C90_ContainsSliceStateMd verifies the new slice-state template file
// is present in the manifest.
func TestManifest_C90_ContainsSliceStateMd(t *testing.T) {
	for _, entry := range installer.Manifest {
		if entry == "templates/slice-state.md" {
			return
		}
	}
	t.Error("installer.Manifest missing new entry: \"templates/slice-state.md\"")
}

// TestManifest_C90_CLAUDEmdLast verifies the invariant that CLAUDE.md is always
// the final entry in the manifest, regardless of how many entries are added.
func TestManifest_C90_CLAUDEmdLast(t *testing.T) {
	if len(installer.Manifest) == 0 {
		t.Fatal("installer.Manifest is empty")
	}

	lastEntry := installer.Manifest[len(installer.Manifest)-1]
	if lastEntry != "CLAUDE.md" {
		t.Errorf("installer.Manifest last entry must be \"CLAUDE.md\", got %q", lastEntry)
	}
}

// TestManifest_C90_MigrateStateAlphabeticalInCommands verifies that
// "commands/gl/migrate-state.md" is positioned alphabetically between
// "commands/gl/map.md" and "commands/gl/pause.md".
func TestManifest_C90_MigrateStateAlphabeticalInCommands(t *testing.T) {
	var commandEntries []string
	for _, entry := range installer.Manifest {
		if strings.HasPrefix(entry, "commands/gl/") {
			commandEntries = append(commandEntries, entry)
		}
	}

	mapIdx := -1
	migrateStateIdx := -1
	pauseIdx := -1

	for i, entry := range commandEntries {
		switch entry {
		case "commands/gl/map.md":
			mapIdx = i
		case "commands/gl/migrate-state.md":
			migrateStateIdx = i
		case "commands/gl/pause.md":
			pauseIdx = i
		}
	}

	if mapIdx == -1 {
		t.Fatal("installer.Manifest missing \"commands/gl/map.md\"")
	}
	if migrateStateIdx == -1 {
		t.Fatal("installer.Manifest missing \"commands/gl/migrate-state.md\"")
	}
	if pauseIdx == -1 {
		t.Fatal("installer.Manifest missing \"commands/gl/pause.md\"")
	}

	if migrateStateIdx <= mapIdx {
		t.Errorf("\"commands/gl/migrate-state.md\" (idx %d) must appear AFTER \"commands/gl/map.md\" (idx %d)", migrateStateIdx, mapIdx)
	}
	if migrateStateIdx >= pauseIdx {
		t.Errorf("\"commands/gl/migrate-state.md\" (idx %d) must appear BEFORE \"commands/gl/pause.md\" (idx %d)", migrateStateIdx, pauseIdx)
	}
}

// TestManifest_C90_StateFormatAlphabeticalInReferences verifies that
// "references/state-format.md" is positioned alphabetically between
// "references/deviation-rules.md" and "references/verification-patterns.md".
func TestManifest_C90_StateFormatAlphabeticalInReferences(t *testing.T) {
	var referenceEntries []string
	for _, entry := range installer.Manifest {
		if strings.HasPrefix(entry, "references/") {
			referenceEntries = append(referenceEntries, entry)
		}
	}

	deviationRulesIdx := -1
	stateFormatIdx := -1
	verificationPatternsIdx := -1

	for i, entry := range referenceEntries {
		switch entry {
		case "references/deviation-rules.md":
			deviationRulesIdx = i
		case "references/state-format.md":
			stateFormatIdx = i
		case "references/verification-patterns.md":
			verificationPatternsIdx = i
		}
	}

	if deviationRulesIdx == -1 {
		t.Fatal("installer.Manifest missing \"references/deviation-rules.md\"")
	}
	if stateFormatIdx == -1 {
		t.Fatal("installer.Manifest missing \"references/state-format.md\"")
	}
	if verificationPatternsIdx == -1 {
		t.Fatal("installer.Manifest missing \"references/verification-patterns.md\"")
	}

	if stateFormatIdx <= deviationRulesIdx {
		t.Errorf("\"references/state-format.md\" (idx %d) must appear AFTER \"references/deviation-rules.md\" (idx %d)", stateFormatIdx, deviationRulesIdx)
	}
	if stateFormatIdx >= verificationPatternsIdx {
		t.Errorf("\"references/state-format.md\" (idx %d) must appear BEFORE \"references/verification-patterns.md\" (idx %d)", stateFormatIdx, verificationPatternsIdx)
	}
}

// TestManifest_C90_SliceStateAlphabeticalInTemplates verifies that
// "templates/slice-state.md" is positioned alphabetically between
// "templates/config.md" and "templates/state.md".
func TestManifest_C90_SliceStateAlphabeticalInTemplates(t *testing.T) {
	var templateEntries []string
	for _, entry := range installer.Manifest {
		if strings.HasPrefix(entry, "templates/") {
			templateEntries = append(templateEntries, entry)
		}
	}

	configIdx := -1
	sliceStateIdx := -1
	stateIdx := -1

	for i, entry := range templateEntries {
		switch entry {
		case "templates/config.md":
			configIdx = i
		case "templates/slice-state.md":
			sliceStateIdx = i
		case "templates/state.md":
			stateIdx = i
		}
	}

	if configIdx == -1 {
		t.Fatal("installer.Manifest missing \"templates/config.md\"")
	}
	if sliceStateIdx == -1 {
		t.Fatal("installer.Manifest missing \"templates/slice-state.md\"")
	}
	if stateIdx == -1 {
		t.Fatal("installer.Manifest missing \"templates/state.md\"")
	}

	if sliceStateIdx <= configIdx {
		t.Errorf("\"templates/slice-state.md\" (idx %d) must appear AFTER \"templates/config.md\" (idx %d)", sliceStateIdx, configIdx)
	}
	if sliceStateIdx >= stateIdx {
		t.Errorf("\"templates/slice-state.md\" (idx %d) must appear BEFORE \"templates/state.md\" (idx %d)", sliceStateIdx, stateIdx)
	}
}

// TestManifest_C90_ContainsAllPrevious35Entries verifies that none of the 35 previously-existing
// entries were removed when the 3 new entries were added.
func TestManifest_C90_ContainsAllPrevious35Entries(t *testing.T) {
	previous35 := []string{
		"agents/gl-architect.md",
		"agents/gl-assessor.md",
		"agents/gl-codebase-mapper.md",
		"agents/gl-debugger.md",
		"agents/gl-designer.md",
		"agents/gl-implementer.md",
		"agents/gl-security.md",
		"agents/gl-test-writer.md",
		"agents/gl-verifier.md",
		"agents/gl-wrapper.md",
		"commands/gl/add-slice.md",
		"commands/gl/assess.md",
		"commands/gl/changelog.md",
		"commands/gl/debug.md",
		"commands/gl/design.md",
		"commands/gl/help.md",
		"commands/gl/init.md",
		"commands/gl/map.md",
		"commands/gl/pause.md",
		"commands/gl/quick.md",
		"commands/gl/resume.md",
		"commands/gl/roadmap.md",
		"commands/gl/settings.md",
		"commands/gl/ship.md",
		"commands/gl/slice.md",
		"commands/gl/status.md",
		"commands/gl/wrap.md",
		"references/checkpoint-protocol.md",
		"references/circuit-breaker.md",
		"references/deviation-rules.md",
		"references/verification-patterns.md",
		"references/verification-tiers.md",
		"templates/config.md",
		"templates/state.md",
		"CLAUDE.md",
	}

	manifestSet := make(map[string]bool, len(installer.Manifest))
	for _, entry := range installer.Manifest {
		manifestSet[entry] = true
	}

	for _, expected := range previous35 {
		if !manifestSet[expected] {
			t.Errorf("installer.Manifest missing previously-existing entry: %q", expected)
		}
	}
}

// TestManifest_C90_NewFilesExistOnDisk verifies that the 3 new source files
// actually exist in the src/ directory alongside all existing content.
// The go:embed wildcards in main.go will pick them up automatically.
func TestManifest_C90_NewFilesExistOnDisk(t *testing.T) {
	root := projectRoot()

	newFiles := []string{
		"src/commands/gl/migrate-state.md",
		"src/references/state-format.md",
		"src/templates/slice-state.md",
	}

	for _, relPath := range newFiles {
		fullPath := filepath.Join(root, relPath)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("new source file does not exist on disk: %s", fullPath)
		}
	}
}

// TestManifest_C90_SectionsAlphabetical verifies that within each section
// (agents/, commands/gl/, references/, templates/), entries are sorted
// alphabetically. This is a structural invariant of the manifest.
func TestManifest_C90_SectionsAlphabetical(t *testing.T) {
	sections := map[string][]string{
		"agents/":       {},
		"commands/gl/":  {},
		"references/":   {},
		"templates/":    {},
	}

	for _, entry := range installer.Manifest {
		for prefix := range sections {
			if strings.HasPrefix(entry, prefix) {
				sections[prefix] = append(sections[prefix], entry)
				break
			}
		}
	}

	for prefix, entries := range sections {
		if len(entries) == 0 {
			t.Errorf("no entries found for section %q in manifest", prefix)
			continue
		}

		sorted := make([]string, len(entries))
		copy(sorted, entries)
		sort.Strings(sorted)

		for i, entry := range entries {
			if entry != sorted[i] {
				t.Errorf("section %q is not alphabetically sorted: at index %d expected %q, got %q", prefix, i, sorted[i], entry)
			}
		}
	}
}
