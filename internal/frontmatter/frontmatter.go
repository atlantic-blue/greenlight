package frontmatter

import (
	"errors"
	"sort"
	"strings"
)

var (
	ErrNoFrontmatter      = errors.New("no frontmatter delimiter found")
	ErrUnclosedFrontmatter = errors.New("frontmatter opening delimiter has no closing delimiter")
	ErrInvalidLine        = errors.New("frontmatter line has no ':' separator")
)

// Parse extracts frontmatter key-value pairs and the body from content.
// Frontmatter must be delimited by "---" on the first non-empty line and a
// second "---" closing delimiter. Returns the parsed fields, remaining body,
// and any error.
func Parse(content string) (map[string]string, string, error) {
	lines := strings.Split(content, "\n")

	openIndex, found := findOpeningDelimiter(lines)
	if !found {
		return nil, "", ErrNoFrontmatter
	}

	closeIndex, found := findClosingDelimiter(lines, openIndex+1)
	if !found {
		return nil, "", ErrUnclosedFrontmatter
	}

	fields, parseError := parseFields(lines[openIndex+1 : closeIndex])
	if parseError != nil {
		return nil, "", parseError
	}

	body := strings.Join(lines[closeIndex+1:], "\n")
	return fields, body, nil
}

// Write serialises frontmatter fields and a body string into a frontmatter
// document. Keys are written in sorted order for deterministic output.
func Write(fields map[string]string, body string) string {
	var builder strings.Builder

	builder.WriteString("---\n")

	keys := sortedKeys(fields)
	for _, key := range keys {
		builder.WriteString(key)
		builder.WriteString(": ")
		builder.WriteString(fields[key])
		builder.WriteString("\n")
	}

	builder.WriteString("---\n")
	builder.WriteString(body)

	return builder.String()
}

// findOpeningDelimiter returns the index of the first non-empty line that is
// exactly "---", or (0, false) if not found.
func findOpeningDelimiter(lines []string) (int, bool) {
	for index, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if line == "---" {
			return index, true
		}
		return 0, false
	}
	return 0, false
}

// findClosingDelimiter searches for "---" starting from startIndex.
func findClosingDelimiter(lines []string, startIndex int) (int, bool) {
	for index := startIndex; index < len(lines); index++ {
		if lines[index] == "---" {
			return index, true
		}
	}
	return 0, false
}

// parseFields converts lines between delimiters into a key-value map.
// Lines containing only whitespace are skipped. Returns ErrInvalidLine if any
// line has no ":" separator.
func parseFields(lines []string) (map[string]string, error) {
	fields := make(map[string]string)

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		colonIndex := strings.Index(line, ":")
		if colonIndex < 0 {
			return nil, ErrInvalidLine
		}

		key := strings.TrimSpace(line[:colonIndex])
		value := strings.TrimSpace(line[colonIndex+1:])
		fields[key] = value
	}

	return fields, nil
}

// sortedKeys returns the keys of a map in sorted order.
func sortedKeys(fields map[string]string) []string {
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
