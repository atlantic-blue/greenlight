package cmd_test

import (
	"path/filepath"
	"testing"

	"github.com/atlantic-blue/greenlight/internal/cmd"
	"github.com/atlantic-blue/greenlight/internal/installer"
)

func TestParseScope(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		wantScope     string
		wantRemaining []string
		wantErr       error
	}{
		{
			name:          "global flag only",
			args:          []string{"--global"},
			wantScope:     "global",
			wantRemaining: []string{},
			wantErr:       nil,
		},
		{
			name:          "local flag only",
			args:          []string{"--local"},
			wantScope:     "local",
			wantRemaining: []string{},
			wantErr:       nil,
		},
		{
			name:          "global with other args",
			args:          []string{"--global", "path/to/file"},
			wantScope:     "global",
			wantRemaining: []string{"path/to/file"},
			wantErr:       nil,
		},
		{
			name:          "local with other args",
			args:          []string{"--local", "path/to/file"},
			wantScope:     "local",
			wantRemaining: []string{"path/to/file"},
			wantErr:       nil,
		},
		{
			name:          "global with multiple remaining args preserves order",
			args:          []string{"--global", "arg1", "arg2", "arg3"},
			wantScope:     "global",
			wantRemaining: []string{"arg1", "arg2", "arg3"},
			wantErr:       nil,
		},
		{
			name:          "local with multiple remaining args preserves order",
			args:          []string{"--local", "--on-conflict=keep", "file.txt"},
			wantScope:     "local",
			wantRemaining: []string{"--on-conflict=keep", "file.txt"},
			wantErr:       nil,
		},
		{
			name:          "global flag in middle of args",
			args:          []string{"arg1", "--global", "arg2"},
			wantScope:     "global",
			wantRemaining: []string{"arg1", "arg2"},
			wantErr:       nil,
		},
		{
			name:          "local flag at end of args",
			args:          []string{"arg1", "arg2", "--local"},
			wantScope:     "local",
			wantRemaining: []string{"arg1", "arg2"},
			wantErr:       nil,
		},
		{
			name:          "both global and local flags",
			args:          []string{"--global", "--local"},
			wantScope:     "",
			wantRemaining: nil,
			wantErr:       cmd.ErrBothScopes,
		},
		{
			name:          "both flags in reverse order",
			args:          []string{"--local", "--global"},
			wantScope:     "",
			wantRemaining: nil,
			wantErr:       cmd.ErrBothScopes,
		},
		{
			name:          "both flags with other args",
			args:          []string{"--global", "file.txt", "--local"},
			wantScope:     "",
			wantRemaining: nil,
			wantErr:       cmd.ErrBothScopes,
		},
		{
			name:          "neither flag present",
			args:          []string{"path/to/file"},
			wantScope:     "",
			wantRemaining: nil,
			wantErr:       cmd.ErrNoScope,
		},
		{
			name:          "empty args",
			args:          []string{},
			wantScope:     "",
			wantRemaining: nil,
			wantErr:       cmd.ErrNoScope,
		},
		{
			name:          "only other flags",
			args:          []string{"--on-conflict=keep", "--verify"},
			wantScope:     "",
			wantRemaining: nil,
			wantErr:       cmd.ErrNoScope,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope, remaining, err := cmd.ParseScope(tt.args)

			if err != tt.wantErr {
				t.Errorf("ParseScope() error = %v, wantErr %v", err, tt.wantErr)
			}

			if scope != tt.wantScope {
				t.Errorf("ParseScope() scope = %q, want %q", scope, tt.wantScope)
			}

			if !equalStringSlices(remaining, tt.wantRemaining) {
				t.Errorf("ParseScope() remaining = %v, want %v", remaining, tt.wantRemaining)
			}

			// Verify invariant: remaining never contains scope flags
			for _, arg := range remaining {
				if arg == "--global" || arg == "--local" {
					t.Errorf("ParseScope() remaining contains scope flag: %q", arg)
				}
			}

			// Verify invariant: err is non-nil if and only if scope is empty
			if (err != nil) != (scope == "") {
				t.Errorf("ParseScope() invariant violated: err=%v but scope=%q", err, scope)
			}

			// Verify invariant: scope is exactly "global" or "local" when err is nil
			if err == nil && scope != "global" && scope != "local" {
				t.Errorf("ParseScope() invalid scope value when err is nil: %q", scope)
			}
		})
	}
}

func TestResolveDir(t *testing.T) {
	tests := []struct {
		name        string
		scope       string
		wantSuffix  string // for global, we check it ends with this
		wantLiteral string // for local, we check exact match
		wantErr     error
	}{
		{
			name:       "local scope returns literal .claude",
			scope:      "local",
			wantSuffix: "",
			wantLiteral: ".claude",
			wantErr:    nil,
		},
		{
			name:        "global scope returns absolute path ending in .claude",
			scope:       "global",
			wantSuffix:  ".claude",
			wantLiteral: "",
			wantErr:     nil,
		},
		{
			name:        "unknown scope returns error",
			scope:       "unknown",
			wantSuffix:  "",
			wantLiteral: "",
			wantErr:     cmd.ErrUnknownScope,
		},
		{
			name:        "empty scope returns error",
			scope:       "",
			wantSuffix:  "",
			wantLiteral: "",
			wantErr:     cmd.ErrUnknownScope,
		},
		{
			name:        "invalid scope value returns error",
			scope:       "workspace",
			wantSuffix:  "",
			wantLiteral: "",
			wantErr:     cmd.ErrUnknownScope,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := cmd.ResolveDir(tt.scope)

			if err != tt.wantErr {
				t.Errorf("ResolveDir() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify invariant: err is non-nil if and only if dir is empty
			if (err != nil) != (dir == "") {
				t.Errorf("ResolveDir() invariant violated: err=%v but dir=%q", err, dir)
			}

			if tt.wantLiteral != "" {
				// Local scope: check exact match
				if dir != tt.wantLiteral {
					t.Errorf("ResolveDir() dir = %q, want %q", dir, tt.wantLiteral)
				}
			}

			if tt.wantSuffix != "" {
				// Global scope: check it's absolute and ends with .claude
				if !filepath.IsAbs(dir) {
					t.Errorf("ResolveDir() global scope should return absolute path, got %q", dir)
				}
				if filepath.Base(dir) != tt.wantSuffix {
					t.Errorf("ResolveDir() dir should end with %q, got %q", tt.wantSuffix, dir)
				}
			}
		})
	}
}

func TestParseConflictStrategy(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		wantStrategy  installer.ConflictStrategy
		wantRemaining []string
		wantErr       error
	}{
		{
			name:          "keep strategy",
			args:          []string{"--on-conflict=keep"},
			wantStrategy:  installer.ConflictKeep,
			wantRemaining: []string{},
			wantErr:       nil,
		},
		{
			name:          "replace strategy",
			args:          []string{"--on-conflict=replace"},
			wantStrategy:  installer.ConflictReplace,
			wantRemaining: []string{},
			wantErr:       nil,
		},
		{
			name:          "append strategy",
			args:          []string{"--on-conflict=append"},
			wantStrategy:  installer.ConflictAppend,
			wantRemaining: []string{},
			wantErr:       nil,
		},
		{
			name:          "strategy with other args",
			args:          []string{"--on-conflict=keep", "path/to/file"},
			wantStrategy:  installer.ConflictKeep,
			wantRemaining: []string{"path/to/file"},
			wantErr:       nil,
		},
		{
			name:          "strategy in middle of args",
			args:          []string{"arg1", "--on-conflict=replace", "arg2"},
			wantStrategy:  installer.ConflictReplace,
			wantRemaining: []string{"arg1", "arg2"},
			wantErr:       nil,
		},
		{
			name:          "multiple remaining args preserves order",
			args:          []string{"--on-conflict=append", "file1", "file2", "file3"},
			wantStrategy:  installer.ConflictAppend,
			wantRemaining: []string{"file1", "file2", "file3"},
			wantErr:       nil,
		},
		{
			name:          "flag absent defaults to keep",
			args:          []string{"path/to/file"},
			wantStrategy:  installer.ConflictKeep,
			wantRemaining: []string{"path/to/file"},
			wantErr:       nil,
		},
		{
			name:          "empty args defaults to keep",
			args:          []string{},
			wantStrategy:  installer.ConflictKeep,
			wantRemaining: []string{},
			wantErr:       nil,
		},
		{
			name:          "only other flags defaults to keep",
			args:          []string{"--global", "--verify"},
			wantStrategy:  installer.ConflictKeep,
			wantRemaining: []string{"--global", "--verify"},
			wantErr:       nil,
		},
		{
			name:          "invalid strategy value returns error",
			args:          []string{"--on-conflict=invalid"},
			wantStrategy:  "",
			wantRemaining: nil,
			wantErr:       cmd.ErrInvalidConflictStrategy,
		},
		{
			name:          "invalid strategy merge",
			args:          []string{"--on-conflict=merge"},
			wantStrategy:  "",
			wantRemaining: nil,
			wantErr:       cmd.ErrInvalidConflictStrategy,
		},
		{
			name:          "invalid strategy overwrite",
			args:          []string{"--on-conflict=overwrite"},
			wantStrategy:  "",
			wantRemaining: nil,
			wantErr:       cmd.ErrInvalidConflictStrategy,
		},
		{
			name:          "empty strategy value returns error",
			args:          []string{"--on-conflict="},
			wantStrategy:  "",
			wantRemaining: nil,
			wantErr:       cmd.ErrInvalidConflictStrategy,
		},
		{
			name:          "invalid strategy with other args",
			args:          []string{"--on-conflict=bad", "file.txt"},
			wantStrategy:  "",
			wantRemaining: nil,
			wantErr:       cmd.ErrInvalidConflictStrategy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy, remaining, err := cmd.ParseConflictStrategy(tt.args)

			if err != tt.wantErr {
				t.Errorf("ParseConflictStrategy() error = %v, wantErr %v", err, tt.wantErr)
			}

			if strategy != tt.wantStrategy {
				t.Errorf("ParseConflictStrategy() strategy = %q, want %q", strategy, tt.wantStrategy)
			}

			if !equalStringSlices(remaining, tt.wantRemaining) {
				t.Errorf("ParseConflictStrategy() remaining = %v, want %v", remaining, tt.wantRemaining)
			}

			// Verify invariant: remaining never contains --on-conflict args
			for _, arg := range remaining {
				if len(arg) >= 14 && arg[:14] == "--on-conflict=" {
					t.Errorf("ParseConflictStrategy() remaining contains --on-conflict flag: %q", arg)
				}
			}

			// Verify invariant: strategy is valid when err is nil
			if err == nil {
				validStrategies := []installer.ConflictStrategy{
					installer.ConflictKeep,
					installer.ConflictReplace,
					installer.ConflictAppend,
				}
				valid := false
				for _, vs := range validStrategies {
					if strategy == vs {
						valid = true
						break
					}
				}
				if !valid {
					t.Errorf("ParseConflictStrategy() invalid strategy when err is nil: %q", strategy)
				}
			}
		})
	}
}

// Helper function to compare string slices
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
