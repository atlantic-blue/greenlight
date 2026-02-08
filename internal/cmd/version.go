package cmd

import (
	"fmt"
	"io"

	"github.com/atlantic-blue/greenlight/internal/version"
)

// RunVersion prints version information to w.
func RunVersion(w io.Writer) int {
	fmt.Fprintf(w, "greenlight %s (commit: %s, built: %s)\n",
		version.Version, version.GitCommit, version.BuildDate)
	return 0
}
