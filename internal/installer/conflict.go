package installer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ConflictStrategy determines how to handle existing CLAUDE.md files.
type ConflictStrategy string

const (
	ConflictKeep    ConflictStrategy = "keep"
	ConflictReplace ConflictStrategy = "replace"
	ConflictAppend  ConflictStrategy = "append"
)

// handleConflict writes srcData to destPath, resolving conflicts with any
// existing file according to the given strategy.
func handleConflict(destPath string, srcData []byte, strategy ConflictStrategy, w io.Writer) error {
	dir := filepath.Dir(destPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	existing, err := os.ReadFile(destPath)
	if os.IsNotExist(err) {
		// No conflict — just write the file
		return os.WriteFile(destPath, srcData, 0o644)
	}
	if err != nil {
		return err
	}

	switch strategy {
	case ConflictKeep:
		// Save greenlight version alongside with a different name
		altPath := filepath.Join(filepath.Dir(destPath), "CLAUDE_GREENLIGHT.md")
		if err := os.WriteFile(altPath, srcData, 0o644); err != nil {
			return err
		}
		fmt.Fprintf(w, "  existing CLAUDE.md kept; greenlight version saved as CLAUDE_GREENLIGHT.md\n")
		return nil

	case ConflictReplace:
		// Backup existing, then overwrite
		backupPath := destPath + ".backup"
		if err := os.WriteFile(backupPath, existing, 0o644); err != nil {
			return fmt.Errorf("creating backup: %w", err)
		}
		fmt.Fprintf(w, "  existing CLAUDE.md backed up to %s\n", backupPath)
		return os.WriteFile(destPath, srcData, 0o644)

	case ConflictAppend:
		// Append greenlight content to existing file
		combined := make([]byte, 0, len(existing)+1+len(srcData))
		combined = append(combined, existing...)
		if len(existing) > 0 && existing[len(existing)-1] != '\n' {
			combined = append(combined, '\n')
		}
		combined = append(combined, srcData...)
		fmt.Fprintf(w, "  greenlight content appended to existing CLAUDE.md\n")
		return os.WriteFile(destPath, combined, 0o644)

	default:
		return fmt.Errorf("unknown conflict strategy: %s", strategy)
	}
}
