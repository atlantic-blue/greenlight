package main

import (
	"embed"
	"io/fs"
	"os"

	"github.com/atlantic-blue/greenlight/internal/cli"
)

//go:embed src/agents/*.md src/commands/gl/*.md src/references/*.md src/templates/*.md src/CLAUDE.md
var embeddedContent embed.FS

func main() {
	contentFS, err := fs.Sub(embeddedContent, "src")
	if err != nil {
		os.Exit(1)
	}
	os.Exit(cli.Run(os.Args[1:], contentFS))
}
