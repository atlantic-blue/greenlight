package version

// Set via ldflags at build time:
//
//	go build -ldflags "-X github.com/atlantic-blue/greenlight/internal/version.Version=1.0.0"
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)
