package shell

import (
	"strings"
	"testing"
)

func TestGenerateZsh(t *testing.T) {
	script := GenerateZsh()

	// Check for required elements
	requiredStrings := []string{
		"wt()",
		"git rev-parse --show-toplevel",
		".wt.yaml",
		"command wt",
		"create)",
		"exit)",
		"cd \"$target\"",
	}

	for _, s := range requiredStrings {
		if !strings.Contains(script, s) {
			t.Errorf("zsh script missing required string: %q", s)
		}
	}
}

func TestGenerateBash(t *testing.T) {
	script := GenerateBash()

	// Check for required elements
	requiredStrings := []string{
		"wt()",
		"git rev-parse --show-toplevel",
		".wt.yaml",
		"command wt",
		"create)",
		"exit)",
	}

	for _, s := range requiredStrings {
		if !strings.Contains(script, s) {
			t.Errorf("bash script missing required string: %q", s)
		}
	}
}

func TestGenerateFish(t *testing.T) {
	script := GenerateFish()

	// Check for required elements
	requiredStrings := []string{
		"function wt",
		"git rev-parse --show-toplevel",
		".wt.yaml",
		"command wt",
		"case create",
		"case exit",
	}

	for _, s := range requiredStrings {
		if !strings.Contains(script, s) {
			t.Errorf("fish script missing required string: %q", s)
		}
	}
}

func TestGenerate(t *testing.T) {
	tests := []struct {
		shell   string
		wantErr bool
	}{
		{"zsh", false},
		{"bash", false},
		{"fish", false},
		{"invalid", true},
		{"", true},
	}

	for _, tt := range tests {
		t.Run(tt.shell, func(t *testing.T) {
			script, err := Generate(tt.shell)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if script == "" {
				t.Error("expected non-empty script")
			}
		})
	}
}

func TestShellScriptsHaveComments(t *testing.T) {
	shells := []string{"zsh", "bash", "fish"}

	for _, sh := range shells {
		t.Run(sh, func(t *testing.T) {
			script, _ := Generate(sh)
			if !strings.HasPrefix(script, "#") {
				t.Error("script should start with a comment")
			}
		})
	}
}

func TestShellScriptsHandleWorktreeDetection(t *testing.T) {
	// All shell scripts should handle the case where we're in a worktree
	// and need to find the main repo's .wt.yaml
	shells := []string{"zsh", "bash", "fish"}

	for _, sh := range shells {
		t.Run(sh, func(t *testing.T) {
			script, _ := Generate(sh)
			// Should check for gitdir in .git file (worktree indicator)
			if !strings.Contains(script, "gitdir") {
				t.Error("script should handle worktree detection via gitdir")
			}
		})
	}
}
