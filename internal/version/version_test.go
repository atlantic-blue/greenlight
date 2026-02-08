package version_test

import (
	"testing"

	"github.com/atlantic-blue/greenlight/internal/version"
)

// TestVersionDefaults_AllNonEmpty verifies that all version variables
// have non-empty default values when not set via ldflags.
func TestVersionDefaults_AllNonEmpty(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "Version has default value",
			value:    version.Version,
			expected: "dev",
		},
		{
			name:     "GitCommit has default value",
			value:    version.GitCommit,
			expected: "unknown",
		},
		{
			name:     "BuildDate has default value",
			value:    version.BuildDate,
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == "" {
				t.Errorf("expected non-empty value, got empty string")
			}

			// Note: We can't assert exact default values in this test
			// because ldflags may have overridden them at build time.
			// The contract guarantees non-empty strings, which we verify above.
		})
	}
}

// TestVersionVariables_AreStrings verifies that all exported version
// variables are accessible and of type string.
func TestVersionVariables_AreStrings(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "Version is accessible",
			value: version.Version,
		},
		{
			name:  "GitCommit is accessible",
			value: version.GitCommit,
		},
		{
			name:  "BuildDate is accessible",
			value: version.BuildDate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Type assertion happens at compile time.
			// If this test compiles, the variables are strings.
			_ = tt.value
		})
	}
}
