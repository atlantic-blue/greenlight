package frontmatter_test

import (
	"errors"
	"testing"

	"github.com/atlantic-blue/greenlight/internal/frontmatter"
)

// ----------------------------------------------------------------------------
// Parse — happy path
// ----------------------------------------------------------------------------

func TestParse_SingleKeyValue(t *testing.T) {
	content := "---\nid: S-35\n---\n"

	fields, body, err := frontmatter.Parse(content)

	if err != nil {
		t.Fatalf("Parse() unexpected error: %v", err)
	}
	if fields["id"] != "S-35" {
		t.Errorf("Parse() fields[\"id\"] = %q, want %q", fields["id"], "S-35")
	}
	if body != "" {
		t.Errorf("Parse() body = %q, want empty string", body)
	}
}

func TestParse_MultipleKeyValuePairs(t *testing.T) {
	content := "---\nid: S-35\nname: Frontmatter Parser\nstatus: complete\n---\n"

	fields, _, err := frontmatter.Parse(content)

	if err != nil {
		t.Fatalf("Parse() unexpected error: %v", err)
	}

	expected := map[string]string{
		"id":     "S-35",
		"name":   "Frontmatter Parser",
		"status": "complete",
	}

	for key, want := range expected {
		if fields[key] != want {
			t.Errorf("Parse() fields[%q] = %q, want %q", key, fields[key], want)
		}
	}

	if len(fields) != len(expected) {
		t.Errorf("Parse() len(fields) = %d, want %d", len(fields), len(expected))
	}
}

func TestParse_ValueContainingColonSplitsOnFirstColonOnly(t *testing.T) {
	content := "---\nurl: https://example.com/path\n---\n"

	fields, _, err := frontmatter.Parse(content)

	if err != nil {
		t.Fatalf("Parse() unexpected error: %v", err)
	}
	if fields["url"] != "https://example.com/path" {
		t.Errorf("Parse() fields[\"url\"] = %q, want %q", fields["url"], "https://example.com/path")
	}
}

func TestParse_EmptyValueIsValid(t *testing.T) {
	content := "---\nkey:\n---\n"

	fields, _, err := frontmatter.Parse(content)

	if err != nil {
		t.Fatalf("Parse() unexpected error: %v", err)
	}
	value, exists := fields["key"]
	if !exists {
		t.Fatal("Parse() expected key to exist in fields map")
	}
	if value != "" {
		t.Errorf("Parse() fields[\"key\"] = %q, want empty string", value)
	}
}

func TestParse_TrimsWhitespaceFromKeysAndValues(t *testing.T) {
	content := "---\n  id  :  S-35  \n---\n"

	fields, _, err := frontmatter.Parse(content)

	if err != nil {
		t.Fatalf("Parse() unexpected error: %v", err)
	}
	if fields["id"] != "S-35" {
		t.Errorf("Parse() fields[\"id\"] = %q, want %q", fields["id"], "S-35")
	}
	if _, hasSpacedKey := fields["  id  "]; hasSpacedKey {
		t.Error("Parse() should trim keys: found key with surrounding whitespace")
	}
}

func TestParse_BodyPreservesOriginalFormatting(t *testing.T) {
	body := "\n## Description\n\nThis is the body.\n  Indented line.\n"
	content := "---\nid: S-35\n---\n" + body

	_, parsedBody, err := frontmatter.Parse(content)

	if err != nil {
		t.Fatalf("Parse() unexpected error: %v", err)
	}
	if parsedBody != body {
		t.Errorf("Parse() body = %q, want %q", parsedBody, body)
	}
}

func TestParse_EmptyFrontmatterReturnsEmptyMapAndEmptyBody(t *testing.T) {
	content := "---\n---\n"

	fields, body, err := frontmatter.Parse(content)

	if err != nil {
		t.Fatalf("Parse() unexpected error: %v", err)
	}
	if len(fields) != 0 {
		t.Errorf("Parse() len(fields) = %d, want 0", len(fields))
	}
	if body != "" {
		t.Errorf("Parse() body = %q, want empty string", body)
	}
}

func TestParse_WhitespaceOnlyLinesBetweenDelimitersAreSkipped(t *testing.T) {
	content := "---\nid: S-35\n   \n\nstatus: pending\n---\n"

	fields, _, err := frontmatter.Parse(content)

	if err != nil {
		t.Fatalf("Parse() unexpected error: %v", err)
	}
	if fields["id"] != "S-35" {
		t.Errorf("Parse() fields[\"id\"] = %q, want %q", fields["id"], "S-35")
	}
	if fields["status"] != "pending" {
		t.Errorf("Parse() fields[\"status\"] = %q, want %q", fields["status"], "pending")
	}
	if len(fields) != 2 {
		t.Errorf("Parse() len(fields) = %d, want 2 (whitespace-only lines should be skipped)", len(fields))
	}
}

// ----------------------------------------------------------------------------
// Parse — error cases
// ----------------------------------------------------------------------------

func TestParse_EmptyContentReturnsErrNoFrontmatter(t *testing.T) {
	_, _, err := frontmatter.Parse("")

	if !errors.Is(err, frontmatter.ErrNoFrontmatter) {
		t.Errorf("Parse(\"\") error = %v, want %v", err, frontmatter.ErrNoFrontmatter)
	}
}

func TestParse_ContentWithNoOpeningDelimiterReturnsErrNoFrontmatter(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "plain text with no delimiters",
			content: "id: S-35\nstatus: pending\n",
		},
		{
			name:    "body content only",
			content: "## Description\n\nThis is a document with no frontmatter.\n",
		},
		{
			name:    "closing delimiter without opening",
			content: "id: S-35\n---\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := frontmatter.Parse(tt.content)

			if !errors.Is(err, frontmatter.ErrNoFrontmatter) {
				t.Errorf("Parse() error = %v, want %v", err, frontmatter.ErrNoFrontmatter)
			}
		})
	}
}

func TestParse_UnclosedFrontmatterReturnsErrUnclosedFrontmatter(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "opening delimiter only with content",
			content: "---\nid: S-35\nstatus: pending\n",
		},
		{
			name:    "opening delimiter only with no content",
			content: "---\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := frontmatter.Parse(tt.content)

			if !errors.Is(err, frontmatter.ErrUnclosedFrontmatter) {
				t.Errorf("Parse() error = %v, want %v", err, frontmatter.ErrUnclosedFrontmatter)
			}
		})
	}
}

func TestParse_LineWithNoColonSeparatorReturnsErrInvalidLine(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "line is a bare word with no colon",
			content: "---\nstatus\n---\n",
		},
		{
			name:    "line looks like a list item",
			content: "---\n- item\n---\n",
		},
		{
			name:    "line is a comment-like string",
			content: "---\n# comment\n---\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := frontmatter.Parse(tt.content)

			if !errors.Is(err, frontmatter.ErrInvalidLine) {
				t.Errorf("Parse() error = %v, want %v", err, frontmatter.ErrInvalidLine)
			}
		})
	}
}

// ----------------------------------------------------------------------------
// Write — happy path
// ----------------------------------------------------------------------------

func TestWrite_ProducesKeyValuePairsInSortedOrder(t *testing.T) {
	fields := map[string]string{
		"status": "pending",
		"id":     "S-35",
		"name":   "Frontmatter Parser",
	}

	result := frontmatter.Write(fields, "")

	// Keys must appear in alphabetical order: id, name, status
	want := "---\nid: S-35\nname: Frontmatter Parser\nstatus: pending\n---\n"
	if result != want {
		t.Errorf("Write() = %q, want %q", result, want)
	}
}

func TestWrite_EmptyFieldsMapProducesMinimalDelimiters(t *testing.T) {
	result := frontmatter.Write(map[string]string{}, "")

	want := "---\n---\n"
	if result != want {
		t.Errorf("Write() with empty fields = %q, want %q", result, want)
	}
}

func TestWrite_EmptyBodyProducesNoTrailingContent(t *testing.T) {
	fields := map[string]string{"id": "S-35"}

	result := frontmatter.Write(fields, "")

	want := "---\nid: S-35\n---\n"
	if result != want {
		t.Errorf("Write() with empty body = %q, want %q", result, want)
	}
}

func TestWrite_AppendBodyAfterClosingDelimiter(t *testing.T) {
	fields := map[string]string{"id": "S-35"}
	body := "## Description\n\nBody content here.\n"

	result := frontmatter.Write(fields, body)

	want := "---\nid: S-35\n---\n## Description\n\nBody content here.\n"
	if result != want {
		t.Errorf("Write() with body = %q, want %q", result, want)
	}
}

func TestWrite_DoesNotAddTrailingNewlineToBody(t *testing.T) {
	fields := map[string]string{"id": "S-35"}
	body := "body without trailing newline"

	result := frontmatter.Write(fields, body)

	want := "---\nid: S-35\n---\nbody without trailing newline"
	if result != want {
		t.Errorf("Write() should not add trailing newline to body, got %q, want %q", result, want)
	}
}

// ----------------------------------------------------------------------------
// Roundtrip invariants
// ----------------------------------------------------------------------------

func TestRoundtrip_WriteOutputCanBeParsedByParse(t *testing.T) {
	originalFields := map[string]string{
		"id":     "S-35",
		"name":   "Frontmatter Parser",
		"status": "pending",
	}
	originalBody := "## Description\n\nSlice body content.\n"

	written := frontmatter.Write(originalFields, originalBody)
	parsedFields, parsedBody, err := frontmatter.Parse(written)

	if err != nil {
		t.Fatalf("Parse(Write()) unexpected error: %v", err)
	}

	for key, want := range originalFields {
		if parsedFields[key] != want {
			t.Errorf("roundtrip: fields[%q] = %q, want %q", key, parsedFields[key], want)
		}
	}

	if len(parsedFields) != len(originalFields) {
		t.Errorf("roundtrip: len(fields) = %d, want %d", len(parsedFields), len(originalFields))
	}

	if parsedBody != originalBody {
		t.Errorf("roundtrip: body = %q, want %q", parsedBody, originalBody)
	}
}

func TestRoundtrip_WithRealisticSliceStateData(t *testing.T) {
	// Realistic data matching .greenlight/slices/*.md frontmatter
	originalFields := map[string]string{
		"id":          "S-35",
		"name":        "Frontmatter Parser",
		"status":      "in-progress",
		"slice":       "S-35",
		"description": "Parse and write flat key-value YAML frontmatter",
		"test-writer": "complete",
		"implementer": "pending",
		"verifier":    "pending",
		"url":         "https://github.com/atlantic-blue/greenlight/issues/35",
	}
	originalBody := "## Contracts\n\n- C-91: FrontmatterParse\n- C-92: FrontmatterWrite\n"

	written := frontmatter.Write(originalFields, originalBody)
	parsedFields, parsedBody, err := frontmatter.Parse(written)

	if err != nil {
		t.Fatalf("Parse(Write()) with realistic data unexpected error: %v", err)
	}

	for key, want := range originalFields {
		if parsedFields[key] != want {
			t.Errorf("roundtrip realistic: fields[%q] = %q, want %q", key, parsedFields[key], want)
		}
	}

	if parsedBody != originalBody {
		t.Errorf("roundtrip realistic: body = %q, want %q", parsedBody, originalBody)
	}
}

func TestRoundtrip_EmptyFrontmatterRoundtrips(t *testing.T) {
	written := frontmatter.Write(map[string]string{}, "")
	fields, body, err := frontmatter.Parse(written)

	if err != nil {
		t.Fatalf("Parse(Write(empty)) unexpected error: %v", err)
	}
	if len(fields) != 0 {
		t.Errorf("roundtrip empty: len(fields) = %d, want 0", len(fields))
	}
	if body != "" {
		t.Errorf("roundtrip empty: body = %q, want empty string", body)
	}
}

// ----------------------------------------------------------------------------
// Write — determinism invariant
// ----------------------------------------------------------------------------

func TestWrite_OutputIsDeterministicAcrossMultipleCalls(t *testing.T) {
	fields := map[string]string{
		"status": "pending",
		"id":     "S-35",
		"name":   "Frontmatter Parser",
		"slice":  "S-35",
	}

	first := frontmatter.Write(fields, "body")
	second := frontmatter.Write(fields, "body")
	third := frontmatter.Write(fields, "body")

	if first != second || second != third {
		t.Errorf("Write() is not deterministic:\n  first  = %q\n  second = %q\n  third  = %q", first, second, third)
	}
}
