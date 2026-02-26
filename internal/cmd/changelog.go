package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const summariesDir = ".greenlight/summaries"

// RunChangelog handles the "changelog" subcommand.
// It reads all .md files from .greenlight/summaries/, sorts them by filename
// ascending, and prints each one separated by "---".
// Returns 1 if .greenlight/ does not exist.
// Returns 0 with "No changelog entries yet." if summaries/ is absent or empty.
func RunChangelog(args []string, stdout io.Writer) int {
	if _, statError := os.Stat(greenlightDir); os.IsNotExist(statError) {
		fmt.Fprintln(stdout, "error: .greenlight/ directory not found. Run 'greenlight init' to set up this greenlight project.")
		return 1
	}

	if _, statError := os.Stat(summariesDir); os.IsNotExist(statError) {
		fmt.Fprintln(stdout, "No changelog entries yet.")
		return 0
	}

	entries, readDirError := os.ReadDir(summariesDir)
	if readDirError != nil {
		fmt.Fprintf(stdout, "warn: could not read summaries directory: %v\n", readDirError)
		fmt.Fprintln(stdout, "No changelog entries yet.")
		return 0
	}

	mdFiles := collectMarkdownFilenames(entries)
	sort.Strings(mdFiles)

	if len(mdFiles) == 0 {
		fmt.Fprintln(stdout, "No changelog entries yet.")
		return 0
	}

	printChangelogEntries(mdFiles, stdout)
	return 0
}

// collectMarkdownFilenames filters directory entries to .md files only.
func collectMarkdownFilenames(entries []os.DirEntry) []string {
	filenames := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".md") {
			filenames = append(filenames, entry.Name())
		}
	}
	return filenames
}

// printChangelogEntries prints each summary file content to stdout, separated
// by "---" between entries (not after the last one). Individual read errors
// cause the file to be skipped with a warning.
func printChangelogEntries(filenames []string, stdout io.Writer) {
	printedCount := 0
	for _, filename := range filenames {
		filePath := filepath.Join(summariesDir, filename)
		contents, readError := os.ReadFile(filePath)
		if readError != nil {
			fmt.Fprintf(stdout, "warn: could not read %s: %v\n", filename, readError)
			continue
		}

		if printedCount > 0 {
			fmt.Fprintln(stdout, "---")
		}

		fmt.Fprint(stdout, string(contents))
		printedCount++
	}
}
