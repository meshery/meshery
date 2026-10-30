package design

import (
	"flag"
	"testing"
)

var (
	update                 = flag.Bool("update", false, "update golden files")
	invalidFilePath        = "/invalid/path/design.yaml"
	validDesignSourceTypes = []string{"Helm Chart", "Kubernetes Manifest", "Docker Compose", "Meshery Design"}
)

func resetVariables() {
	skipSave = false
	designFile = ""
	file = ""
	designDeployFlags = cmdDesignDeployFlags{}
	applyContexts = nil
	undeployContexts = nil
}

func TestAppendContextsQuery(t *testing.T) {
	tests := []struct {
		name        string
		rawURL      string
		contexts    []string
		expected    string
		expectErr   bool
		expectedErr string
	}{
		{
			name:     "empty contexts defaults to all",
			rawURL:   "api/pattern/deploy",
			contexts: nil,
			expected: "api/pattern/deploy?contexts=all",
		},
		{
			name:     "empty slice contexts defaults to all",
			rawURL:   "http://localhost:9081/api/pattern/deploy",
			contexts: []string{},
			expected: "http://localhost:9081/api/pattern/deploy?contexts=all",
		},
		{
			name:     "whitespace contexts defaults to all",
			rawURL:   "api/pattern/deploy",
			contexts: []string{"   ", "", "\t", "\n"},
			expected: "api/pattern/deploy?contexts=all",
		},
		{
			name:     "pflag empty default representation defaults to all",
			rawURL:   "api/pattern/deploy",
			contexts: []string{"[]"},
			expected: "api/pattern/deploy?contexts=all",
		},
		{
			name:     "all alone defaults to all",
			rawURL:   "api/pattern/deploy",
			contexts: []string{"all"},
			expected: "api/pattern/deploy?contexts=all",
		},
		{
			name:     "all case-insensitive and repeated deduplicates to single all",
			rawURL:   "api/pattern/deploy",
			contexts: []string{"ALL", "all"},
			expected: "api/pattern/deploy?contexts=all",
		},
		{
			name:        "all combined with context ID prefix returns error",
			rawURL:      "api/pattern/deploy",
			contexts:    []string{"all", "ctx-1"},
			expectErr:   true,
			expectedErr: "cannot combine 'all' with specific context IDs",
		},
		{
			name:        "all combined with context ID suffix returns error",
			rawURL:      "api/pattern/deploy",
			contexts:    []string{"ctx-1", "all"},
			expectErr:   true,
			expectedErr: "cannot combine 'all' with specific context IDs",
		},
		{
			name:        "all case-insensitive combined with multiple context IDs returns error",
			rawURL:      "api/pattern/deploy",
			contexts:    []string{"ctx-1", "ALL", "ctx-2"},
			expectErr:   true,
			expectedErr: "cannot combine 'all' with specific context IDs",
		},
		{
			name:        "all combined with context ID in single comma-separated string returns error",
			rawURL:      "api/pattern/deploy",
			contexts:    []string{"all,ctx-1"},
			expectErr:   true,
			expectedErr: "cannot combine 'all' with specific context IDs",
		},
		{
			name:     "single context",
			rawURL:   "api/pattern/deploy",
			contexts: []string{"ctx-1"},
			expected: "api/pattern/deploy?contexts=ctx-1",
		},
		{
			name:     "multiple explicit contexts in separate elements",
			rawURL:   "http://localhost:9081/api/pattern/deploy",
			contexts: []string{"ctx-1", "ctx-2"},
			expected: "http://localhost:9081/api/pattern/deploy?contexts=ctx-1&contexts=ctx-2",
		},
		{
			name:     "multiple explicit contexts in comma-separated string",
			rawURL:   "http://localhost:9081/api/pattern/deploy",
			contexts: []string{"ctx-1,ctx-2"},
			expected: "http://localhost:9081/api/pattern/deploy?contexts=ctx-1&contexts=ctx-2",
		},
		{
			name:     "multiple explicit contexts with duplicates deduplicates",
			rawURL:   "http://localhost:9081/api/pattern/deploy",
			contexts: []string{"ctx-1", "ctx-1", "ctx-2"},
			expected: "http://localhost:9081/api/pattern/deploy?contexts=ctx-1&contexts=ctx-2",
		},
		{
			name:     "multiple explicit contexts with comma duplicates deduplicates",
			rawURL:   "http://localhost:9081/api/pattern/deploy",
			contexts: []string{"ctx-1,ctx-2", "ctx-1"},
			expected: "http://localhost:9081/api/pattern/deploy?contexts=ctx-1&contexts=ctx-2",
		},
		{
			name:     "preserves existing query parameters and replaces stale contexts",
			rawURL:   "api/pattern/deploy?contexts=old&foo=bar",
			contexts: []string{"ctx-1"},
			expected: "api/pattern/deploy?contexts=ctx-1&foo=bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := appendContextsQuery(tt.rawURL, tt.contexts)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.expectedErr)
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("error = %q, want %q", err.Error(), tt.expectedErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("appendContextsQuery() = %q, want %q", got, tt.expected)
			}
		})
	}
}
