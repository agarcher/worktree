package commands

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/agarcher/wt/internal/git"
)

func TestParseHookKeyValues(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantPairs []KeyValue
		wantRaw   []string
	}{
		{
			name:      "empty input",
			input:     "",
			wantPairs: nil,
			wantRaw:   nil,
		},
		{
			name:  "single key-value",
			input: "Port: 3000",
			wantPairs: []KeyValue{
				{Key: "Port", Value: "3000"},
			},
			wantRaw: nil,
		},
		{
			name:  "multiple key-values",
			input: "Port: 3000\nDatabase: dev_db\nAPI: http://localhost:3000",
			wantPairs: []KeyValue{
				{Key: "Port", Value: "3000"},
				{Key: "Database", Value: "dev_db"},
				{Key: "API", Value: "http://localhost:3000"},
			},
			wantRaw: nil,
		},
		{
			name:  "key with spaces",
			input: "API Endpoint: http://localhost:3000/api",
			wantPairs: []KeyValue{
				{Key: "API Endpoint", Value: "http://localhost:3000/api"},
			},
			wantRaw: nil,
		},
		{
			name:  "value with colons",
			input: "URL: http://localhost:3000",
			wantPairs: []KeyValue{
				{Key: "URL", Value: "http://localhost:3000"},
			},
			wantRaw: nil,
		},
		{
			name:      "raw text without colon",
			input:     "This is just some text",
			wantPairs: nil,
			wantRaw:   []string{"This is just some text"},
		},
		{
			name:  "mixed key-value and raw",
			input: "Port: 3000\nSome informational message\nDatabase: dev",
			wantPairs: []KeyValue{
				{Key: "Port", Value: "3000"},
				{Key: "Database", Value: "dev"},
			},
			wantRaw: []string{"Some informational message"},
		},
		{
			name:      "empty lines ignored",
			input:     "\n\nPort: 3000\n\n",
			wantPairs: []KeyValue{{Key: "Port", Value: "3000"}},
			wantRaw:   nil,
		},
		{
			name:      "whitespace trimmed",
			input:     "  Port  :  3000  ",
			wantPairs: []KeyValue{{Key: "Port", Value: "3000"}},
			wantRaw:   nil,
		},
		{
			name:      "key with no value treated as raw",
			input:     "SomeKey:",
			wantPairs: nil,
			wantRaw:   []string{"SomeKey:"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pairs, raw := ParseHookKeyValues(tt.input)

			// Check pairs
			if len(pairs) != len(tt.wantPairs) {
				t.Errorf("got %d pairs, want %d", len(pairs), len(tt.wantPairs))
			} else {
				for i, p := range pairs {
					if p.Key != tt.wantPairs[i].Key || p.Value != tt.wantPairs[i].Value {
						t.Errorf("pair %d: got {%q, %q}, want {%q, %q}",
							i, p.Key, p.Value, tt.wantPairs[i].Key, tt.wantPairs[i].Value)
					}
				}
			}

			// Check raw lines
			if len(raw) != len(tt.wantRaw) {
				t.Errorf("got %d raw lines, want %d", len(raw), len(tt.wantRaw))
			} else {
				for i, r := range raw {
					if r != tt.wantRaw[i] {
						t.Errorf("raw %d: got %q, want %q", i, r, tt.wantRaw[i])
					}
				}
			}
		})
	}
}

func TestPrintVerboseWorktree(t *testing.T) {
	tests := []struct {
		name     string
		info     VerboseInfo
		contains []string
	}{
		{
			name: "basic info",
			info: VerboseInfo{
				Name:          "feature-test",
				Branch:        "feature-test",
				Index:         1,
				CurrentMarker: "* ",
				Status:        &git.WorktreeStatus{},
			},
			contains: []string{"* feature-test", "Branch:", "feature-test", "Index:", "1"},
		},
		{
			name: "with status",
			info: VerboseInfo{
				Name:          "feature-test",
				Branch:        "feature-test",
				CurrentMarker: "  ",
				Status: &git.WorktreeStatus{
					CommitsAhead:          3,
					HasUncommittedChanges: true,
				},
			},
			contains: []string{"feature-test", "Status:", "↑3", "dirty"},
		},
		{
			name: "with created date",
			info: VerboseInfo{
				Name:          "feature-test",
				Branch:        "feature-test",
				CurrentMarker: "  ",
				CreatedAt:     time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
				Status:        &git.WorktreeStatus{},
			},
			contains: []string{"Created:", "2025-01-10", "ago)"},
		},
		{
			name: "with hook key-value output",
			info: VerboseInfo{
				Name:          "feature-test",
				Branch:        "feature-test",
				CurrentMarker: "  ",
				Status:        &git.WorktreeStatus{},
				HookOutput:    "Port: 3001\nDatabase: dev_test",
			},
			contains: []string{"Port:", "3001", "Database:", "dev_test"},
		},
		{
			name: "with hook raw output",
			info: VerboseInfo{
				Name:          "feature-test",
				Branch:        "feature-test",
				CurrentMarker: "  ",
				Status:        &git.WorktreeStatus{},
				HookOutput:    "This is informational text",
			},
			contains: []string{"This is informational text"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			PrintVerboseWorktree(&buf, tt.info)
			output := buf.String()

			for _, s := range tt.contains {
				if !strings.Contains(output, s) {
					t.Errorf("output missing %q\nGot:\n%s", s, output)
				}
			}
		})
	}
}

func TestFormatAge(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{
			name:     "less than an hour",
			duration: 30 * time.Minute,
			want:     "less than an hour",
		},
		{
			name:     "1 hour",
			duration: 1 * time.Hour,
			want:     "1 hour",
		},
		{
			name:     "multiple hours",
			duration: 5 * time.Hour,
			want:     "5 hours",
		},
		{
			name:     "1 day",
			duration: 24 * time.Hour,
			want:     "1 day",
		},
		{
			name:     "multiple days",
			duration: 3 * 24 * time.Hour,
			want:     "3 days",
		},
		{
			name:     "1 week",
			duration: 7 * 24 * time.Hour,
			want:     "1 week",
		},
		{
			name:     "multiple weeks",
			duration: 21 * 24 * time.Hour,
			want:     "3 weeks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatAge(tt.duration)
			if got != tt.want {
				t.Errorf("formatAge(%v) = %q, want %q", tt.duration, got, tt.want)
			}
		})
	}
}

func TestKeyValueAlignment(t *testing.T) {
	// Verify that keys are properly aligned when there are hook outputs
	info := VerboseInfo{
		Name:          "test",
		Branch:        "test-branch",
		Index:         1,
		CurrentMarker: "  ",
		Status:        &git.WorktreeStatus{},
		HookOutput:    "API Endpoint: http://localhost:3000\nPort: 3000",
	}

	var buf bytes.Buffer
	PrintVerboseWorktree(&buf, info)
	output := buf.String()

	// All colons should be at the same position (after max key width)
	lines := strings.Split(output, "\n")
	var colonPositions []int
	for _, line := range lines {
		if strings.Contains(line, ":") && strings.HasPrefix(strings.TrimSpace(line), "  ") == false {
			// Skip the header line
			continue
		}
		if idx := strings.Index(line, ":"); idx > 0 {
			colonPositions = append(colonPositions, idx)
		}
	}

	// All colons should be aligned (at the same position)
	if len(colonPositions) > 1 {
		first := colonPositions[0]
		for i, pos := range colonPositions[1:] {
			if pos != first {
				t.Errorf("colon position mismatch: line 0 at %d, line %d at %d", first, i+1, pos)
			}
		}
	}
}
