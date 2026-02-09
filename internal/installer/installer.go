package installer

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/atlantic-blue/greenlight/internal/version"
)

// Manifest lists all files that greenlight installs (relative to the content FS root).
var Manifest = []string{
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
	"references/deviation-rules.md",
	"references/verification-patterns.md",
	"templates/config.md",
	"templates/state.md",
	"CLAUDE.md",
}

// Installer copies embedded content to the filesystem.
type Installer struct {
	contentFS fs.FS
	stdout    io.Writer
}

// New creates an Installer backed by the given content filesystem.
func New(contentFS fs.FS, stdout io.Writer) *Installer {
	return &Installer{contentFS: contentFS, stdout: stdout}
}

// Install copies all manifest files to targetDir, handling CLAUDE.md
// according to the conflict strategy. scope is "global" or "local".
func (inst *Installer) Install(targetDir, scope string, strategy ConflictStrategy) error {
	for _, relPath := range Manifest {
		if relPath == "CLAUDE.md" {
			if err := inst.installCLAUDE(targetDir, scope, strategy); err != nil {
				return fmt.Errorf("installing CLAUDE.md: %w", err)
			}
			continue
		}
		destPath := filepath.Join(targetDir, relPath)
		if err := inst.copyFile(relPath, destPath); err != nil {
			return fmt.Errorf("installing %s: %w", relPath, err)
		}
		fmt.Fprintf(inst.stdout, "  installed %s\n", relPath)
	}

	if err := inst.writeVersionFile(targetDir); err != nil {
		return fmt.Errorf("writing version file: %w", err)
	}

	fmt.Fprintf(inst.stdout, "greenlight installed to %s\n", targetDir)
	return nil
}

// installCLAUDE handles the asymmetric CLAUDE.md placement:
//   - global: ~/.claude/CLAUDE.md
//   - local: ./CLAUDE.md (project root, one level above .claude/)
func (inst *Installer) installCLAUDE(targetDir, scope string, strategy ConflictStrategy) error {
	var destPath string
	switch scope {
	case "global":
		destPath = filepath.Join(targetDir, "CLAUDE.md")
	case "local":
		// targetDir is .claude, CLAUDE.md goes to parent (project root)
		destPath = filepath.Join(filepath.Dir(targetDir), "CLAUDE.md")
		if targetDir == ".claude" {
			destPath = "CLAUDE.md"
		}
	default:
		destPath = filepath.Join(targetDir, "CLAUDE.md")
	}

	srcData, err := fs.ReadFile(inst.contentFS, "CLAUDE.md")
	if err != nil {
		return err
	}

	if err := handleConflict(destPath, srcData, strategy, inst.stdout); err != nil {
		return err
	}
	fmt.Fprintf(inst.stdout, "  installed CLAUDE.md -> %s\n", destPath)
	return nil
}

// copyFile reads a file from the embedded FS and writes it to destPath.
func (inst *Installer) copyFile(srcPath, destPath string) error {
	data, err := fs.ReadFile(inst.contentFS, srcPath)
	if err != nil {
		return err
	}

	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	return os.WriteFile(destPath, data, 0o644)
}

// writeVersionFile writes a .greenlight-version file to targetDir.
func (inst *Installer) writeVersionFile(targetDir string) error {
	content := fmt.Sprintf("%s\n%s\n%s\n", version.Version, version.GitCommit, version.BuildDate)
	return os.WriteFile(filepath.Join(targetDir, ".greenlight-version"), []byte(content), 0o644)
}

// Uninstall removes greenlight-managed files from targetDir.
// It only removes files listed in the manifest plus the version file.
func Uninstall(targetDir, scope string, stdout io.Writer) error {
	for _, relPath := range Manifest {
		if relPath == "CLAUDE.md" {
			// Don't remove CLAUDE.md — it may have user content
			continue
		}
		destPath := filepath.Join(targetDir, relPath)
		if err := os.Remove(destPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("removing %s: %w", relPath, err)
		}
		fmt.Fprintf(stdout, "  removed %s\n", relPath)
	}

	// Remove version file
	versionPath := filepath.Join(targetDir, ".greenlight-version")
	os.Remove(versionPath)

	// Remove conflict artifacts using asymmetric path resolution
	var artifactDir string
	if scope == "global" {
		artifactDir = targetDir
	} else if scope == "local" {
		// For local scope, artifacts are in parent of targetDir
		if targetDir == ".claude" {
			artifactDir = "."
		} else {
			artifactDir = filepath.Dir(targetDir)
		}
	}

	// Remove CLAUDE_GREENLIGHT.md if present
	greenlightPath := filepath.Join(artifactDir, "CLAUDE_GREENLIGHT.md")
	if err := os.Remove(greenlightPath); err == nil {
		fmt.Fprintf(stdout, "  removed CLAUDE_GREENLIGHT.md\n")
	}

	// Remove CLAUDE.md.backup if present
	backupPath := filepath.Join(artifactDir, "CLAUDE.md.backup")
	if err := os.Remove(backupPath); err == nil {
		fmt.Fprintf(stdout, "  removed CLAUDE.md.backup\n")
	}

	// Clean up empty directories (deepest first)
	cleanEmptyDirs(targetDir, "commands/gl")
	cleanEmptyDirs(targetDir, "commands")
	cleanEmptyDirs(targetDir, "agents")
	cleanEmptyDirs(targetDir, "references")
	cleanEmptyDirs(targetDir, "templates")

	fmt.Fprintf(stdout, "greenlight uninstalled from %s\n", targetDir)
	return nil
}

// cleanEmptyDirs removes a directory if it's empty.
func cleanEmptyDirs(base, sub string) {
	dir := filepath.Join(base, sub)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	if len(entries) == 0 {
		os.Remove(dir)
	}
}

// Check verifies that all expected files are present and non-empty in targetDir.
// When verify=true, it also checks that file contents match contentFS using SHA-256.
func Check(targetDir, scope string, stdout io.Writer, verify bool, contentFS fs.FS) (ok bool) {
	ok = true
	missing := 0
	empty := 0
	modified := 0

	for _, relPath := range Manifest {
		var destPath string
		if relPath == "CLAUDE.md" {
			switch scope {
			case "local":
				if targetDir == ".claude" {
					destPath = "CLAUDE.md"
				} else {
					destPath = filepath.Join(filepath.Dir(targetDir), "CLAUDE.md")
				}
			default:
				destPath = filepath.Join(targetDir, "CLAUDE.md")
			}
		} else {
			destPath = filepath.Join(targetDir, relPath)
		}

		info, err := os.Stat(destPath)
		if os.IsNotExist(err) {
			fmt.Fprintf(stdout, "  MISSING  %s\n", relPath)
			missing++
			ok = false
			continue
		}
		if err != nil {
			fmt.Fprintf(stdout, "  ERROR    %s: %v\n", relPath, err)
			ok = false
			continue
		}
		if info.Size() == 0 {
			fmt.Fprintf(stdout, "  EMPTY    %s\n", relPath)
			empty++
			ok = false
			continue
		}

		// If verify mode is enabled, check content hash
		if verify {
			if contentFS == nil {
				// Cannot verify without contentFS
				fmt.Fprintf(stdout, "  MODIFIED %s\n", relPath)
				modified++
				ok = false
				continue
			}

			// Read file from disk
			diskData, err := os.ReadFile(destPath)
			if err != nil {
				fmt.Fprintf(stdout, "  ERROR    %s: %v\n", relPath, err)
				ok = false
				continue
			}

			// Read expected content from contentFS
			expectedData, err := fs.ReadFile(contentFS, relPath)
			if err != nil {
				// If file doesn't exist in contentFS, treat as modified
				fmt.Fprintf(stdout, "  MODIFIED %s\n", relPath)
				modified++
				ok = false
				continue
			}

			// Compare SHA-256 hashes
			diskHash := sha256.Sum256(diskData)
			expectedHash := sha256.Sum256(expectedData)
			if diskHash != expectedHash {
				fmt.Fprintf(stdout, "  MODIFIED %s\n", relPath)
				modified++
				ok = false
				continue
			}
		}
	}

	// Check version file
	versionPath := filepath.Join(targetDir, ".greenlight-version")
	if _, err := os.Stat(versionPath); os.IsNotExist(err) {
		fmt.Fprintf(stdout, "  MISSING  .greenlight-version\n")
		// In presence mode, missing version file fails the check
		// In verify mode, version file is informational only
		if !verify {
			ok = false
		}
	} else {
		data, _ := os.ReadFile(versionPath)
		parts := strings.SplitN(string(data), "\n", 2)
		if len(parts) > 0 {
			fmt.Fprintf(stdout, "  version: %s\n", strings.TrimSpace(parts[0]))
		}
	}

	total := len(Manifest)
	verified := total - missing - empty - modified
	if verify {
		if ok {
			fmt.Fprintf(stdout, "all %d files verified\n", total)
		} else {
			fmt.Fprintf(stdout, "%d/%d files verified (%d missing, %d empty, %d modified)\n", verified, total, missing, empty, modified)
		}
	} else {
		present := total - missing
		if ok {
			fmt.Fprintf(stdout, "all %d files present\n", total)
		} else {
			fmt.Fprintf(stdout, "%d/%d files present (%d missing, %d empty)\n", present, total, missing, empty)
		}
	}
	return ok
}
