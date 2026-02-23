package state_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/atlantic-blue/greenlight/internal/state"
)

// ----------------------------------------------------------------------------
// ReadSlices — happy path
// ----------------------------------------------------------------------------

func TestReadSlices_SingleSliceFile(t *testing.T) {
	dir := t.TempDir()
	content := "---\nid: S-01\nstatus: pending\nstep: test-writer\nmilestone: core\nstarted: 2026-01-01\nupdated: 2026-01-01T10:00:00Z\ntests: 5\nsecurity_tests: 2\nsession: abc123\ndeps:\n---\n"
	writeFile(t, filepath.Join(dir, "S-01.md"), content)

	slices, err := state.ReadSlices(dir)

	if err != nil {
		t.Fatalf("ReadSlices() unexpected error: %v", err)
	}
	if len(slices) != 1 {
		t.Fatalf("ReadSlices() len = %d, want 1", len(slices))
	}

	got := slices[0]
	if got.ID != "S-01" {
		t.Errorf("ReadSlices() ID = %q, want %q", got.ID, "S-01")
	}
	if got.Status != "pending" {
		t.Errorf("ReadSlices() Status = %q, want %q", got.Status, "pending")
	}
	if got.Step != "test-writer" {
		t.Errorf("ReadSlices() Step = %q, want %q", got.Step, "test-writer")
	}
	if got.Milestone != "core" {
		t.Errorf("ReadSlices() Milestone = %q, want %q", got.Milestone, "core")
	}
	if got.Started != "2026-01-01" {
		t.Errorf("ReadSlices() Started = %q, want %q", got.Started, "2026-01-01")
	}
	if got.Updated != "2026-01-01T10:00:00Z" {
		t.Errorf("ReadSlices() Updated = %q, want %q", got.Updated, "2026-01-01T10:00:00Z")
	}
	if got.Tests != 5 {
		t.Errorf("ReadSlices() Tests = %d, want 5", got.Tests)
	}
	if got.SecurityTests != 2 {
		t.Errorf("ReadSlices() SecurityTests = %d, want 2", got.SecurityTests)
	}
	if got.Session != "abc123" {
		t.Errorf("ReadSlices() Session = %q, want %q", got.Session, "abc123")
	}
	if len(got.Deps) != 0 {
		t.Errorf("ReadSlices() Deps = %v, want empty slice", got.Deps)
	}
}

func TestReadSlices_MultipleSlicesAreSortedByID(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "S-03.md"), "---\nid: S-03\nstatus: pending\n---\n")
	writeFile(t, filepath.Join(dir, "S-01.md"), "---\nid: S-01\nstatus: complete\n---\n")
	writeFile(t, filepath.Join(dir, "S-02.md"), "---\nid: S-02\nstatus: in_progress\n---\n")

	slices, err := state.ReadSlices(dir)

	if err != nil {
		t.Fatalf("ReadSlices() unexpected error: %v", err)
	}
	if len(slices) != 3 {
		t.Fatalf("ReadSlices() len = %d, want 3", len(slices))
	}
	if slices[0].ID != "S-01" {
		t.Errorf("ReadSlices() slices[0].ID = %q, want %q", slices[0].ID, "S-01")
	}
	if slices[1].ID != "S-02" {
		t.Errorf("ReadSlices() slices[1].ID = %q, want %q", slices[1].ID, "S-02")
	}
	if slices[2].ID != "S-03" {
		t.Errorf("ReadSlices() slices[2].ID = %q, want %q", slices[2].ID, "S-03")
	}
}

func TestReadSlices_DepsAreParsedFromCommaSeparatedString(t *testing.T) {
	dir := t.TempDir()
	content := "---\nid: S-04\nstatus: pending\ndeps: S-01, S-02, S-03\n---\n"
	writeFile(t, filepath.Join(dir, "S-04.md"), content)

	slices, err := state.ReadSlices(dir)

	if err != nil {
		t.Fatalf("ReadSlices() unexpected error: %v", err)
	}
	if len(slices) != 1 {
		t.Fatalf("ReadSlices() len = %d, want 1", len(slices))
	}

	deps := slices[0].Deps
	if len(deps) != 3 {
		t.Fatalf("ReadSlices() Deps len = %d, want 3", len(deps))
	}
	if deps[0] != "S-01" {
		t.Errorf("ReadSlices() Deps[0] = %q, want %q", deps[0], "S-01")
	}
	if deps[1] != "S-02" {
		t.Errorf("ReadSlices() Deps[1] = %q, want %q", deps[1], "S-02")
	}
	if deps[2] != "S-03" {
		t.Errorf("ReadSlices() Deps[2] = %q, want %q", deps[2], "S-03")
	}
}

func TestReadSlices_EmptyDepsFieldYieldsEmptySlice(t *testing.T) {
	dir := t.TempDir()
	content := "---\nid: S-01\nstatus: pending\ndeps:\n---\n"
	writeFile(t, filepath.Join(dir, "S-01.md"), content)

	slices, err := state.ReadSlices(dir)

	if err != nil {
		t.Fatalf("ReadSlices() unexpected error: %v", err)
	}
	if len(slices[0].Deps) != 0 {
		t.Errorf("ReadSlices() Deps = %v, want empty slice", slices[0].Deps)
	}
}

func TestReadSlices_SingleDepWithNoCommaIsParsedCorrectly(t *testing.T) {
	dir := t.TempDir()
	content := "---\nid: S-02\nstatus: pending\ndeps: S-01\n---\n"
	writeFile(t, filepath.Join(dir, "S-02.md"), content)

	slices, err := state.ReadSlices(dir)

	if err != nil {
		t.Fatalf("ReadSlices() unexpected error: %v", err)
	}
	if len(slices[0].Deps) != 1 {
		t.Fatalf("ReadSlices() Deps len = %d, want 1", len(slices[0].Deps))
	}
	if slices[0].Deps[0] != "S-01" {
		t.Errorf("ReadSlices() Deps[0] = %q, want %q", slices[0].Deps[0], "S-01")
	}
}

func TestReadSlices_NonNumericTestCountDefaultsToZero(t *testing.T) {
	dir := t.TempDir()
	content := "---\nid: S-01\nstatus: pending\ntests: not-a-number\nsecurity_tests: also-bad\n---\n"
	writeFile(t, filepath.Join(dir, "S-01.md"), content)

	slices, err := state.ReadSlices(dir)

	if err != nil {
		t.Fatalf("ReadSlices() unexpected error: %v", err)
	}
	if slices[0].Tests != 0 {
		t.Errorf("ReadSlices() Tests = %d, want 0 for non-numeric value", slices[0].Tests)
	}
	if slices[0].SecurityTests != 0 {
		t.Errorf("ReadSlices() SecurityTests = %d, want 0 for non-numeric value", slices[0].SecurityTests)
	}
}

func TestReadSlices_MissingTestCountFieldDefaultsToZero(t *testing.T) {
	dir := t.TempDir()
	content := "---\nid: S-01\nstatus: pending\n---\n"
	writeFile(t, filepath.Join(dir, "S-01.md"), content)

	slices, err := state.ReadSlices(dir)

	if err != nil {
		t.Fatalf("ReadSlices() unexpected error: %v", err)
	}
	if slices[0].Tests != 0 {
		t.Errorf("ReadSlices() Tests = %d, want 0 when field is absent", slices[0].Tests)
	}
	if slices[0].SecurityTests != 0 {
		t.Errorf("ReadSlices() SecurityTests = %d, want 0 when field is absent", slices[0].SecurityTests)
	}
}

func TestReadSlices_UnknownStatusIsPreservedAsIs(t *testing.T) {
	dir := t.TempDir()
	content := "---\nid: S-01\nstatus: weird-custom-value\n---\n"
	writeFile(t, filepath.Join(dir, "S-01.md"), content)

	slices, err := state.ReadSlices(dir)

	if err != nil {
		t.Fatalf("ReadSlices() unexpected error: %v", err)
	}
	if slices[0].Status != "weird-custom-value" {
		t.Errorf("ReadSlices() Status = %q, want %q", slices[0].Status, "weird-custom-value")
	}
}

func TestReadSlices_NonMdFilesInDirectoryAreIgnored(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "S-01.md"), "---\nid: S-01\nstatus: pending\n---\n")
	writeFile(t, filepath.Join(dir, "notes.txt"), "some notes")
	writeFile(t, filepath.Join(dir, "config.json"), `{"key":"value"}`)

	slices, err := state.ReadSlices(dir)

	if err != nil {
		t.Fatalf("ReadSlices() unexpected error: %v", err)
	}
	if len(slices) != 1 {
		t.Errorf("ReadSlices() len = %d, want 1 (non-.md files should be ignored)", len(slices))
	}
}

// ----------------------------------------------------------------------------
// ReadSlices — error cases
// ----------------------------------------------------------------------------

func TestReadSlices_DirectoryDoesNotExistReturnsErrDirNotFound(t *testing.T) {
	_, err := state.ReadSlices("/nonexistent/path/that/does/not/exist")

	if !errors.Is(err, state.ErrDirNotFound) {
		t.Errorf("ReadSlices() error = %v, want %v", err, state.ErrDirNotFound)
	}
}

func TestReadSlices_EmptyDirectoryReturnsErrNoSliceFiles(t *testing.T) {
	dir := t.TempDir()

	_, err := state.ReadSlices(dir)

	if !errors.Is(err, state.ErrNoSliceFiles) {
		t.Errorf("ReadSlices() error = %v, want %v", err, state.ErrNoSliceFiles)
	}
}

func TestReadSlices_DirectoryWithOnlyNonMdFilesReturnsErrNoSliceFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "notes.txt"), "some notes")
	writeFile(t, filepath.Join(dir, "data.json"), `{}`)

	_, err := state.ReadSlices(dir)

	if !errors.Is(err, state.ErrNoSliceFiles) {
		t.Errorf("ReadSlices() error = %v, want %v", err, state.ErrNoSliceFiles)
	}
}

func TestReadSlices_FileWithInvalidFrontmatterReturnsErrParseFailure(t *testing.T) {
	dir := t.TempDir()
	// No opening delimiter — frontmatter.Parse will fail
	writeFile(t, filepath.Join(dir, "S-01.md"), "id: S-01\nstatus: pending\n")

	_, err := state.ReadSlices(dir)

	if !errors.Is(err, state.ErrParseFailure) {
		t.Errorf("ReadSlices() error = %v, want %v", err, state.ErrParseFailure)
	}
}

func TestReadSlices_ErrParseFailureIncludesFilename(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "S-99.md"), "id: S-99\nstatus: pending\n")

	_, err := state.ReadSlices(dir)

	if err == nil {
		t.Fatal("ReadSlices() expected error, got nil")
	}
	errMsg := err.Error()
	if !containsString(errMsg, "S-99.md") {
		t.Errorf("ReadSlices() error %q does not include filename %q", errMsg, "S-99.md")
	}
}

// ----------------------------------------------------------------------------
// ReadGraph — happy path
// ----------------------------------------------------------------------------

func TestReadGraph_ValidGraphWithSlicesAndEdges(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "GRAPH.json")
	content := `{
		"slices": {
			"S-01": {
				"id": "S-01",
				"name": "Version",
				"depends_on": [],
				"wave": 1,
				"contracts": ["C-01"]
			},
			"S-02": {
				"id": "S-02",
				"name": "Flag Parsing",
				"depends_on": ["S-01"],
				"wave": 2,
				"contracts": ["C-02", "C-03"]
			}
		},
		"edges": [
			{"from": "S-02", "to": "S-01", "reason": "S-02 requires version infrastructure"}
		]
	}`
	writeFile(t, path, content)

	graph, err := state.ReadGraph(path)

	if err != nil {
		t.Fatalf("ReadGraph() unexpected error: %v", err)
	}
	if graph == nil {
		t.Fatal("ReadGraph() returned nil graph")
	}
	if len(graph.Slices) != 2 {
		t.Errorf("ReadGraph() len(Slices) = %d, want 2", len(graph.Slices))
	}

	s01, exists := graph.Slices["S-01"]
	if !exists {
		t.Fatal("ReadGraph() Slices[\"S-01\"] missing")
	}
	if s01.ID != "S-01" {
		t.Errorf("ReadGraph() Slices[\"S-01\"].ID = %q, want %q", s01.ID, "S-01")
	}
	if s01.Name != "Version" {
		t.Errorf("ReadGraph() Slices[\"S-01\"].Name = %q, want %q", s01.Name, "Version")
	}
	if s01.Wave != 1 {
		t.Errorf("ReadGraph() Slices[\"S-01\"].Wave = %d, want 1", s01.Wave)
	}
	if len(s01.Contracts) != 1 || s01.Contracts[0] != "C-01" {
		t.Errorf("ReadGraph() Slices[\"S-01\"].Contracts = %v, want [C-01]", s01.Contracts)
	}

	s02, exists := graph.Slices["S-02"]
	if !exists {
		t.Fatal("ReadGraph() Slices[\"S-02\"] missing")
	}
	if len(s02.DependsOn) != 1 || s02.DependsOn[0] != "S-01" {
		t.Errorf("ReadGraph() Slices[\"S-02\"].DependsOn = %v, want [S-01]", s02.DependsOn)
	}

	if len(graph.Edges) != 1 {
		t.Fatalf("ReadGraph() len(Edges) = %d, want 1", len(graph.Edges))
	}
	edge := graph.Edges[0]
	if edge.From != "S-02" {
		t.Errorf("ReadGraph() Edges[0].From = %q, want %q", edge.From, "S-02")
	}
	if edge.To != "S-01" {
		t.Errorf("ReadGraph() Edges[0].To = %q, want %q", edge.To, "S-01")
	}
	if edge.Reason != "S-02 requires version infrastructure" {
		t.Errorf("ReadGraph() Edges[0].Reason = %q, want %q", edge.Reason, "S-02 requires version infrastructure")
	}
}

func TestReadGraph_UnknownFieldsAreIgnored(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "GRAPH.json")
	content := `{
		"slices": {
			"S-01": {
				"id": "S-01",
				"name": "Version",
				"depends_on": [],
				"wave": 1,
				"contracts": [],
				"future_field": "some future value",
				"another_unknown": 42
			}
		},
		"edges": [],
		"project": "greenlight",
		"description": "future top-level field"
	}`
	writeFile(t, path, content)

	graph, err := state.ReadGraph(path)

	if err != nil {
		t.Fatalf("ReadGraph() unexpected error with unknown fields: %v", err)
	}
	if len(graph.Slices) != 1 {
		t.Errorf("ReadGraph() len(Slices) = %d, want 1", len(graph.Slices))
	}
}

func TestReadGraph_EmptyEdgesArrayIsValid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "GRAPH.json")
	content := `{"slices": {"S-01": {"id": "S-01", "name": "Version", "depends_on": [], "wave": 1, "contracts": []}}, "edges": []}`
	writeFile(t, path, content)

	graph, err := state.ReadGraph(path)

	if err != nil {
		t.Fatalf("ReadGraph() unexpected error: %v", err)
	}
	if len(graph.Edges) != 0 {
		t.Errorf("ReadGraph() len(Edges) = %d, want 0", len(graph.Edges))
	}
}

func TestReadGraph_MissingOptionalFieldsDefaultToZeroValues(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "GRAPH.json")
	// Slice entry with only id, everything else absent
	content := `{"slices": {"S-01": {"id": "S-01"}}}`
	writeFile(t, path, content)

	graph, err := state.ReadGraph(path)

	if err != nil {
		t.Fatalf("ReadGraph() unexpected error: %v", err)
	}
	s01 := graph.Slices["S-01"]
	if s01.Wave != 0 {
		t.Errorf("ReadGraph() Slices[\"S-01\"].Wave = %d, want 0 (zero value)", s01.Wave)
	}
	if s01.Name != "" {
		t.Errorf("ReadGraph() Slices[\"S-01\"].Name = %q, want empty string (zero value)", s01.Name)
	}
	if len(s01.DependsOn) != 0 {
		t.Errorf("ReadGraph() Slices[\"S-01\"].DependsOn = %v, want nil/empty", s01.DependsOn)
	}
	if len(s01.Contracts) != 0 {
		t.Errorf("ReadGraph() Slices[\"S-01\"].Contracts = %v, want nil/empty", s01.Contracts)
	}
}

// ----------------------------------------------------------------------------
// ReadGraph — error cases
// ----------------------------------------------------------------------------

func TestReadGraph_FileDoesNotExistReturnsErrFileNotFound(t *testing.T) {
	_, err := state.ReadGraph("/nonexistent/GRAPH.json")

	if !errors.Is(err, state.ErrFileNotFound) {
		t.Errorf("ReadGraph() error = %v, want %v", err, state.ErrFileNotFound)
	}
}

func TestReadGraph_InvalidJSONReturnsErrInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "GRAPH.json")
	writeFile(t, path, `{this is not valid json`)

	_, err := state.ReadGraph(path)

	if !errors.Is(err, state.ErrInvalidJSON) {
		t.Errorf("ReadGraph() error = %v, want %v", err, state.ErrInvalidJSON)
	}
}

func TestReadGraph_ValidJSONWithNoSlicesFieldReturnsErrMissingSlices(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "GRAPH.json")
	writeFile(t, path, `{"edges": [], "project": "greenlight"}`)

	_, err := state.ReadGraph(path)

	if !errors.Is(err, state.ErrMissingSlices) {
		t.Errorf("ReadGraph() error = %v, want %v", err, state.ErrMissingSlices)
	}
}

func TestReadGraph_EmptyFileReturnsErrInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "GRAPH.json")
	writeFile(t, path, "")

	_, err := state.ReadGraph(path)

	if !errors.Is(err, state.ErrInvalidJSON) {
		t.Errorf("ReadGraph() error = %v, want %v", err, state.ErrInvalidJSON)
	}
}

// ----------------------------------------------------------------------------
// FindReadySlices — happy path
// ----------------------------------------------------------------------------

func TestFindReadySlices_PendingSliceWithNoDepsIsReady(t *testing.T) {
	slices := []state.SliceInfo{
		{ID: "S-01", Status: "pending", Deps: []string{}},
	}
	graph := &state.Graph{
		Slices: map[string]state.GraphSlice{},
		Edges:  []state.Edge{},
	}

	ready := state.FindReadySlices(slices, graph)

	if len(ready) != 1 {
		t.Fatalf("FindReadySlices() len = %d, want 1", len(ready))
	}
	if ready[0].ID != "S-01" {
		t.Errorf("FindReadySlices() ready[0].ID = %q, want %q", ready[0].ID, "S-01")
	}
}

func TestFindReadySlices_PendingSliceWithAllDepsCompleteIsReady(t *testing.T) {
	slices := []state.SliceInfo{
		{ID: "S-01", Status: "complete", Deps: []string{}},
		{ID: "S-02", Status: "complete", Deps: []string{"S-01"}},
		{ID: "S-03", Status: "pending", Deps: []string{"S-01", "S-02"}},
	}
	graph := &state.Graph{
		Slices: map[string]state.GraphSlice{
			"S-03": {ID: "S-03", DependsOn: []string{"S-01", "S-02"}, Wave: 2},
		},
		Edges: []state.Edge{},
	}

	ready := state.FindReadySlices(slices, graph)

	if len(ready) != 1 {
		t.Fatalf("FindReadySlices() len = %d, want 1", len(ready))
	}
	if ready[0].ID != "S-03" {
		t.Errorf("FindReadySlices() ready[0].ID = %q, want %q", ready[0].ID, "S-03")
	}
}

func TestFindReadySlices_ReadySlicesAreSortedByWaveThenID(t *testing.T) {
	slices := []state.SliceInfo{
		{ID: "S-04", Status: "pending", Deps: []string{}},
		{ID: "S-02", Status: "pending", Deps: []string{}},
		{ID: "S-03", Status: "pending", Deps: []string{}},
		{ID: "S-01", Status: "pending", Deps: []string{}},
	}
	graph := &state.Graph{
		Slices: map[string]state.GraphSlice{
			"S-01": {ID: "S-01", Wave: 2},
			"S-02": {ID: "S-02", Wave: 1},
			"S-03": {ID: "S-03", Wave: 2},
			"S-04": {ID: "S-04", Wave: 1},
		},
		Edges: []state.Edge{},
	}

	ready := state.FindReadySlices(slices, graph)

	if len(ready) != 4 {
		t.Fatalf("FindReadySlices() len = %d, want 4", len(ready))
	}
	// Wave 1 comes first, sorted by ID within wave
	if ready[0].ID != "S-02" {
		t.Errorf("FindReadySlices() ready[0].ID = %q, want %q (wave 1, first ID)", ready[0].ID, "S-02")
	}
	if ready[1].ID != "S-04" {
		t.Errorf("FindReadySlices() ready[1].ID = %q, want %q (wave 1, second ID)", ready[1].ID, "S-04")
	}
	// Wave 2 comes second, sorted by ID within wave
	if ready[2].ID != "S-01" {
		t.Errorf("FindReadySlices() ready[2].ID = %q, want %q (wave 2, first ID)", ready[2].ID, "S-01")
	}
	if ready[3].ID != "S-03" {
		t.Errorf("FindReadySlices() ready[3].ID = %q, want %q (wave 2, second ID)", ready[3].ID, "S-03")
	}
}

func TestFindReadySlices_SliceNotInGraphIsTreatedAsNoDeps(t *testing.T) {
	slices := []state.SliceInfo{
		{ID: "S-01", Status: "pending", Deps: []string{}},
	}
	// Graph has no entry for S-01
	graph := &state.Graph{
		Slices: map[string]state.GraphSlice{},
		Edges:  []state.Edge{},
	}

	ready := state.FindReadySlices(slices, graph)

	if len(ready) != 1 {
		t.Fatalf("FindReadySlices() len = %d, want 1 (no graph entry = no deps)", len(ready))
	}
	if ready[0].ID != "S-01" {
		t.Errorf("FindReadySlices() ready[0].ID = %q, want %q", ready[0].ID, "S-01")
	}
}

func TestFindReadySlices_OutputIsDeterministicForSameInput(t *testing.T) {
	slices := []state.SliceInfo{
		{ID: "S-03", Status: "pending", Deps: []string{}},
		{ID: "S-01", Status: "pending", Deps: []string{}},
		{ID: "S-02", Status: "pending", Deps: []string{}},
	}
	graph := &state.Graph{
		Slices: map[string]state.GraphSlice{
			"S-01": {ID: "S-01", Wave: 1},
			"S-02": {ID: "S-02", Wave: 1},
			"S-03": {ID: "S-03", Wave: 1},
		},
		Edges: []state.Edge{},
	}

	first := state.FindReadySlices(slices, graph)
	second := state.FindReadySlices(slices, graph)

	if len(first) != len(second) {
		t.Fatalf("FindReadySlices() not deterministic: len first=%d, second=%d", len(first), len(second))
	}
	for i := range first {
		if first[i].ID != second[i].ID {
			t.Errorf("FindReadySlices() not deterministic: first[%d].ID=%q, second[%d].ID=%q", i, first[i].ID, i, second[i].ID)
		}
	}
}

// ----------------------------------------------------------------------------
// FindReadySlices — sad path
// ----------------------------------------------------------------------------

func TestFindReadySlices_PendingSliceWithIncompleteDepsIsNotReady(t *testing.T) {
	slices := []state.SliceInfo{
		{ID: "S-01", Status: "pending", Deps: []string{}},
		{ID: "S-02", Status: "pending", Deps: []string{"S-01"}},
	}
	graph := &state.Graph{
		Slices: map[string]state.GraphSlice{
			"S-02": {ID: "S-02", DependsOn: []string{"S-01"}, Wave: 2},
		},
		Edges: []state.Edge{},
	}

	ready := state.FindReadySlices(slices, graph)

	// Only S-01 is ready; S-02 depends on S-01 which is still pending
	if len(ready) != 1 {
		t.Fatalf("FindReadySlices() len = %d, want 1", len(ready))
	}
	if ready[0].ID != "S-01" {
		t.Errorf("FindReadySlices() ready[0].ID = %q, want %q", ready[0].ID, "S-01")
	}
}

func TestFindReadySlices_InProgressSliceIsNeverReady(t *testing.T) {
	slices := []state.SliceInfo{
		{ID: "S-01", Status: "in_progress", Deps: []string{}},
	}
	graph := &state.Graph{
		Slices: map[string]state.GraphSlice{},
		Edges:  []state.Edge{},
	}

	ready := state.FindReadySlices(slices, graph)

	if len(ready) != 0 {
		t.Errorf("FindReadySlices() len = %d, want 0 (in_progress is never ready)", len(ready))
	}
}

func TestFindReadySlices_CompleteSliceIsNeverReady(t *testing.T) {
	slices := []state.SliceInfo{
		{ID: "S-01", Status: "complete", Deps: []string{}},
	}
	graph := &state.Graph{
		Slices: map[string]state.GraphSlice{},
		Edges:  []state.Edge{},
	}

	ready := state.FindReadySlices(slices, graph)

	if len(ready) != 0 {
		t.Errorf("FindReadySlices() len = %d, want 0 (complete is never ready)", len(ready))
	}
}

func TestFindReadySlices_FailedSliceIsNeverReady(t *testing.T) {
	slices := []state.SliceInfo{
		{ID: "S-01", Status: "failed", Deps: []string{}},
	}
	graph := &state.Graph{
		Slices: map[string]state.GraphSlice{},
		Edges:  []state.Edge{},
	}

	ready := state.FindReadySlices(slices, graph)

	if len(ready) != 0 {
		t.Errorf("FindReadySlices() len = %d, want 0 (failed is never ready)", len(ready))
	}
}

func TestFindReadySlices_SliceWithFailedDepIsNotReady(t *testing.T) {
	slices := []state.SliceInfo{
		{ID: "S-01", Status: "failed", Deps: []string{}},
		{ID: "S-02", Status: "pending", Deps: []string{"S-01"}},
	}
	graph := &state.Graph{
		Slices: map[string]state.GraphSlice{
			"S-02": {ID: "S-02", DependsOn: []string{"S-01"}, Wave: 2},
		},
		Edges: []state.Edge{},
	}

	ready := state.FindReadySlices(slices, graph)

	if len(ready) != 0 {
		t.Errorf("FindReadySlices() len = %d, want 0 (dep is failed, not complete)", len(ready))
	}
}

func TestFindReadySlices_SliceWithInProgressDepIsNotReady(t *testing.T) {
	slices := []state.SliceInfo{
		{ID: "S-01", Status: "in_progress", Deps: []string{}},
		{ID: "S-02", Status: "pending", Deps: []string{"S-01"}},
	}
	graph := &state.Graph{
		Slices: map[string]state.GraphSlice{
			"S-02": {ID: "S-02", DependsOn: []string{"S-01"}, Wave: 2},
		},
		Edges: []state.Edge{},
	}

	ready := state.FindReadySlices(slices, graph)

	if len(ready) != 0 {
		t.Errorf("FindReadySlices() len = %d, want 0 (dep is in_progress, not complete)", len(ready))
	}
}

// ----------------------------------------------------------------------------
// FindReadySlices — edge cases
// ----------------------------------------------------------------------------

func TestFindReadySlices_EmptyInputReturnsEmptyOutput(t *testing.T) {
	graph := &state.Graph{
		Slices: map[string]state.GraphSlice{},
		Edges:  []state.Edge{},
	}

	ready := state.FindReadySlices([]state.SliceInfo{}, graph)

	if len(ready) != 0 {
		t.Errorf("FindReadySlices() len = %d, want 0 for empty input", len(ready))
	}
}

func TestFindReadySlices_NilSlicesInputReturnsEmptyOutput(t *testing.T) {
	graph := &state.Graph{
		Slices: map[string]state.GraphSlice{},
		Edges:  []state.Edge{},
	}

	ready := state.FindReadySlices(nil, graph)

	if len(ready) != 0 {
		t.Errorf("FindReadySlices() len = %d, want 0 for nil input", len(ready))
	}
}

func TestFindReadySlices_OnlyPendingSlicesWithAllCompleteDepsAreIncluded(t *testing.T) {
	slices := []state.SliceInfo{
		{ID: "S-01", Status: "complete", Deps: []string{}},
		{ID: "S-02", Status: "in_progress", Deps: []string{"S-01"}},
		{ID: "S-03", Status: "failed", Deps: []string{"S-01"}},
		{ID: "S-04", Status: "pending", Deps: []string{"S-01"}},
		{ID: "S-05", Status: "pending", Deps: []string{"S-02"}},
	}
	graph := &state.Graph{
		Slices: map[string]state.GraphSlice{
			"S-04": {ID: "S-04", DependsOn: []string{"S-01"}, Wave: 2},
			"S-05": {ID: "S-05", DependsOn: []string{"S-02"}, Wave: 3},
		},
		Edges: []state.Edge{},
	}

	ready := state.FindReadySlices(slices, graph)

	// Only S-04 is ready: pending, dep S-01 is complete
	// S-05 is pending but dep S-02 is in_progress, not complete
	if len(ready) != 1 {
		t.Fatalf("FindReadySlices() len = %d, want 1", len(ready))
	}
	if ready[0].ID != "S-04" {
		t.Errorf("FindReadySlices() ready[0].ID = %q, want %q", ready[0].ID, "S-04")
	}
}

// ----------------------------------------------------------------------------
// DetectContext — happy path
// ----------------------------------------------------------------------------

func TestDetectContext_ClaudeCodeSetReturnsInsideClaudeTrue(t *testing.T) {
	t.Setenv("CLAUDE_CODE", "1")

	ctx := state.DetectContext()

	if !ctx.InsideClaude {
		t.Errorf("DetectContext() InsideClaude = false, want true when CLAUDE_CODE is set")
	}
}

func TestDetectContext_ClaudeCodeValueIsPreserved(t *testing.T) {
	t.Setenv("CLAUDE_CODE", "claude-sonnet-4-6")

	ctx := state.DetectContext()

	if ctx.ClaudeValue != "claude-sonnet-4-6" {
		t.Errorf("DetectContext() ClaudeValue = %q, want %q", ctx.ClaudeValue, "claude-sonnet-4-6")
	}
}

func TestDetectContext_ClaudeCodeSetToArbitraryNonEmptyValueIsInsideClaude(t *testing.T) {
	t.Setenv("CLAUDE_CODE", "any-non-empty-value")

	ctx := state.DetectContext()

	if !ctx.InsideClaude {
		t.Errorf("DetectContext() InsideClaude = false, want true for any non-empty CLAUDE_CODE value")
	}
	if ctx.ClaudeValue != "any-non-empty-value" {
		t.Errorf("DetectContext() ClaudeValue = %q, want %q", ctx.ClaudeValue, "any-non-empty-value")
	}
}

// ----------------------------------------------------------------------------
// DetectContext — sad path
// ----------------------------------------------------------------------------

func TestDetectContext_ClaudeCodeUnsetReturnsInsideClaudeFalse(t *testing.T) {
	// Ensure the variable is not set for this test
	t.Setenv("CLAUDE_CODE", "")
	os.Unsetenv("CLAUDE_CODE")

	ctx := state.DetectContext()

	if ctx.InsideClaude {
		t.Errorf("DetectContext() InsideClaude = true, want false when CLAUDE_CODE is unset")
	}
	if ctx.ClaudeValue != "" {
		t.Errorf("DetectContext() ClaudeValue = %q, want empty string when CLAUDE_CODE is unset", ctx.ClaudeValue)
	}
}

func TestDetectContext_ClaudeCodeEmptyStringReturnsInsideClaudeFalse(t *testing.T) {
	t.Setenv("CLAUDE_CODE", "")

	ctx := state.DetectContext()

	if ctx.InsideClaude {
		t.Errorf("DetectContext() InsideClaude = true, want false when CLAUDE_CODE is empty string")
	}
	if ctx.ClaudeValue != "" {
		t.Errorf("DetectContext() ClaudeValue = %q, want empty string", ctx.ClaudeValue)
	}
}

// ----------------------------------------------------------------------------
// helpers
// ----------------------------------------------------------------------------

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writeFile(%q) error: %v", path, err)
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
