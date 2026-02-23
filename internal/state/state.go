package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/atlantic-blue/greenlight/internal/frontmatter"
)

// Sentinel errors for ReadSlices.
var (
	ErrDirNotFound  = errors.New("directory not found")
	ErrNoSliceFiles = errors.New("no .md slice files found in directory")
	ErrParseFailure = errors.New("failed to parse frontmatter")
)

// Sentinel errors for ReadGraph.
var (
	ErrFileNotFound  = errors.New("file not found")
	ErrInvalidJSON   = errors.New("invalid JSON")
	ErrMissingSlices = errors.New("graph JSON is missing required 'slices' field")
)

// SliceInfo holds the parsed frontmatter fields from a slice .md file.
type SliceInfo struct {
	ID            string
	Status        string
	Step          string
	Milestone     string
	Started       string
	Updated       string
	Session       string
	Tests         int
	SecurityTests int
	Deps          []string
}

// GraphSlice represents a single slice node in the dependency graph.
type GraphSlice struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	DependsOn []string `json:"depends_on"`
	Wave      int      `json:"wave"`
	Contracts []string `json:"contracts"`
}

// Edge represents a directed dependency edge between two slices.
type Edge struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Reason string `json:"reason"`
}

// Graph holds the full dependency graph parsed from GRAPH.json.
type Graph struct {
	Slices map[string]GraphSlice
	Edges  []Edge
}

// ExecutionContext describes the runtime environment.
type ExecutionContext struct {
	InsideClaude bool
	ClaudeValue  string
}

// ReadSlices reads all .md files from dir, parses frontmatter from each, and
// returns a slice of SliceInfo sorted by ID ascending.
func ReadSlices(dir string) ([]SliceInfo, error) {
	entries, readError := os.ReadDir(dir)
	if readError != nil {
		return nil, fmt.Errorf("%w: %s", ErrDirNotFound, dir)
	}

	var mdFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			mdFiles = append(mdFiles, entry.Name())
		}
	}

	if len(mdFiles) == 0 {
		return nil, ErrNoSliceFiles
	}

	slices := make([]SliceInfo, 0, len(mdFiles))
	for _, filename := range mdFiles {
		fullPath := dir + "/" + filename
		raw, readFileError := os.ReadFile(fullPath)
		if readFileError != nil {
			return nil, fmt.Errorf("%w: %s", ErrParseFailure, filename)
		}

		fields, _, parseError := frontmatter.Parse(string(raw))
		if parseError != nil {
			return nil, fmt.Errorf("%w: %s", ErrParseFailure, filename)
		}

		info := buildSliceInfo(fields)
		slices = append(slices, info)
	}

	sort.Slice(slices, func(i, j int) bool {
		return slices[i].ID < slices[j].ID
	})

	return slices, nil
}

// buildSliceInfo converts parsed frontmatter fields into a SliceInfo struct.
func buildSliceInfo(fields map[string]string) SliceInfo {
	info := SliceInfo{
		ID:        fields["id"],
		Status:    fields["status"],
		Step:      fields["step"],
		Milestone: fields["milestone"],
		Started:   fields["started"],
		Updated:   fields["updated"],
		Session:   fields["session"],
	}

	if testsStr, exists := fields["tests"]; exists {
		if value, convertError := strconv.Atoi(testsStr); convertError == nil {
			info.Tests = value
		}
	}

	if secStr, exists := fields["security_tests"]; exists {
		if value, convertError := strconv.Atoi(secStr); convertError == nil {
			info.SecurityTests = value
		}
	}

	info.Deps = parseDeps(fields["deps"])
	return info
}

// parseDeps splits a comma-separated deps string into a slice of trimmed IDs.
// An empty or whitespace-only string returns an empty (non-nil) slice.
func parseDeps(raw string) []string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return []string{}
	}

	parts := strings.Split(trimmed, ",")
	deps := make([]string, 0, len(parts))
	for _, part := range parts {
		dep := strings.TrimSpace(part)
		if dep != "" {
			deps = append(deps, dep)
		}
	}
	return deps
}

// rawGraph is used internally to detect whether "slices" key was present in JSON.
type rawGraph struct {
	Slices json.RawMessage `json:"slices"`
	Edges  []Edge          `json:"edges"`
}

// ReadGraph reads the GRAPH.json file at path and parses it into a Graph.
func ReadGraph(path string) (*Graph, error) {
	data, readError := os.ReadFile(path)
	if readError != nil {
		return nil, fmt.Errorf("%w: %s", ErrFileNotFound, path)
	}

	var raw rawGraph
	if unmarshalError := json.Unmarshal(data, &raw); unmarshalError != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidJSON, unmarshalError.Error())
	}

	if raw.Slices == nil {
		return nil, ErrMissingSlices
	}

	var slicesMap map[string]GraphSlice
	if unmarshalError := json.Unmarshal(raw.Slices, &slicesMap); unmarshalError != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidJSON, unmarshalError.Error())
	}

	edges := raw.Edges
	if edges == nil {
		edges = []Edge{}
	}

	return &Graph{
		Slices: slicesMap,
		Edges:  edges,
	}, nil
}

// FindReadySlices returns the subset of pending slices whose graph dependencies
// are all complete. Slices not in the graph are treated as having no deps.
// The result is sorted by wave ascending, then ID ascending.
func FindReadySlices(slices []SliceInfo, graph *Graph) []SliceInfo {
	if len(slices) == 0 {
		return []SliceInfo{}
	}

	statusByID := buildStatusMap(slices)

	var ready []SliceInfo
	for _, slice := range slices {
		if slice.Status != "pending" {
			continue
		}

		if allDepsComplete(slice.ID, graph, statusByID) {
			ready = append(ready, slice)
		}
	}

	sortByWaveThenID(ready, graph)
	return ready
}

// buildStatusMap creates a map of slice ID to status for fast lookup.
func buildStatusMap(slices []SliceInfo) map[string]string {
	statusByID := make(map[string]string, len(slices))
	for _, slice := range slices {
		statusByID[slice.ID] = slice.Status
	}
	return statusByID
}

// allDepsComplete returns true if every dependency in the graph for the given
// slice ID has status "complete".
func allDepsComplete(sliceID string, graph *Graph, statusByID map[string]string) bool {
	graphSlice, exists := graph.Slices[sliceID]
	if !exists {
		return true
	}

	for _, depID := range graphSlice.DependsOn {
		if statusByID[depID] != "complete" {
			return false
		}
	}
	return true
}

// sortByWaveThenID sorts a slice of SliceInfo by wave ascending, then ID ascending.
func sortByWaveThenID(slices []SliceInfo, graph *Graph) {
	sort.Slice(slices, func(i, j int) bool {
		waveI := waveOf(slices[i].ID, graph)
		waveJ := waveOf(slices[j].ID, graph)
		if waveI != waveJ {
			return waveI < waveJ
		}
		return slices[i].ID < slices[j].ID
	})
}

// waveOf returns the wave number for a slice from the graph, defaulting to 0.
func waveOf(sliceID string, graph *Graph) int {
	if graphSlice, exists := graph.Slices[sliceID]; exists {
		return graphSlice.Wave
	}
	return 0
}

// DetectContext reads the $CLAUDE_CODE environment variable and returns an
// ExecutionContext indicating whether the process is running inside Claude Code.
func DetectContext() ExecutionContext {
	value, set := os.LookupEnv("CLAUDE_CODE")
	if !set || value == "" {
		return ExecutionContext{
			InsideClaude: false,
			ClaudeValue:  "",
		}
	}

	return ExecutionContext{
		InsideClaude: true,
		ClaudeValue:  value,
	}
}
