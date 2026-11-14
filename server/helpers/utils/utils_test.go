package utils

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSplitAndTrim(t *testing.T) {
	tests := []struct {
		name string
		in   string
		sep  string
		want []string
	}{
		{
			name: "empty input returns nil",
			in:   "",
			sep:  ",",
			want: nil,
		},
		{
			name: "single value, no separator present",
			in:   "https://cloud.meshery.io",
			sep:  ",",
			want: []string{"https://cloud.meshery.io"},
		},
		{
			name: "two comma-separated values",
			in:   "https://cloud.meshery.io,https://cloud.acme.io",
			sep:  ",",
			want: []string{"https://cloud.meshery.io", "https://cloud.acme.io"},
		},
		{
			name: "trims whitespace around entries",
			in:   "  https://a.example  ,\thttps://b.example\n",
			sep:  ",",
			want: []string{"https://a.example", "https://b.example"},
		},
		{
			name: "drops empty entries from trailing or doubled separators",
			in:   "https://a.example,,https://b.example,",
			sep:  ",",
			want: []string{"https://a.example", "https://b.example"},
		},
		{
			name: "whitespace-only entries are dropped",
			in:   "https://a.example,   ,https://b.example",
			sep:  ",",
			want: []string{"https://a.example", "https://b.example"},
		},
		{
			name: "space-delimited ADAPTER_URLS form",
			in:   "meshery-istio:10000 meshery-linkerd:10001 meshery-consul:10002",
			sep:  ", \t\n\r",
			want: []string{"meshery-istio:10000", "meshery-linkerd:10001", "meshery-consul:10002"},
		},
		{
			name: "mixed comma and whitespace delimiters",
			in:   "a,b c\td\ne",
			sep:  ", \t\n\r",
			want: []string{"a", "b", "c", "d", "e"},
		},
		{
			name: "runs of mixed delimiters collapse",
			in:   "a , ,\tb",
			sep:  ", \t",
			want: []string{"a", "b"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := SplitAndTrim(tc.in, tc.sep)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("SplitAndTrim(%q, %q) = %#v, want %#v", tc.in, tc.sep, got, tc.want)
			}
		})
	}
}

func TestWriteSVGHelper_TraversalGuard(t *testing.T) {
	// Resolve ../../ui inside a temporary tree instead of the source checkout.
	workDir := filepath.Join(t.TempDir(), "helpers", "utils")
	if err := os.MkdirAll(workDir, 0755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(workDir)

	uiSVGMutex.Lock()
	savedPaths := UISVGPaths
	UISVGPaths = make([]string, 0)
	uiSVGMutex.Unlock()
	t.Cleanup(func() {
		uiSVGMutex.Lock()
		UISVGPaths = savedPaths
		uiSVGMutex.Unlock()
	})

	tests := []struct {
		name     string
		dirname  string
		filename string
		wantSafe bool
	}{
		{
			name:     "valid dirname and filename",
			dirname:  "test-model",
			filename: "test-kind",
			wantSafe: true,
		},
		{
			name:     "traversal dirname relative parent",
			dirname:  "../outside",
			filename: "test-kind",
			wantSafe: false,
		},
		{
			name:     "traversal dirname deep dotdot",
			dirname:  "../../../../etc",
			filename: "test-kind",
			wantSafe: false,
		},
		{
			name:     "traversal dirname single dot",
			dirname:  ".",
			filename: "test-kind",
			wantSafe: false,
		},
		{
			name:     "traversal dirname double dot",
			dirname:  "..",
			filename: "test-kind",
			wantSafe: false,
		},
		{
			name:     "traversal filename dotdot",
			dirname:  "test-model",
			filename: "../evil",
			wantSafe: false,
		},
		{
			name:     "traversal filename nested path",
			dirname:  "test-model",
			filename: "foo/bar",
			wantSafe: false,
		},
		{
			name:     "traversal filename windows backslash",
			dirname:  "test-model",
			filename: `..\evil`,
			wantSafe: false,
		},
		{
			name:     "empty dirname",
			dirname:  "",
			filename: "test-kind",
			wantSafe: false,
		},
		{
			name:     "empty filename",
			dirname:  "test-model",
			filename: "",
			wantSafe: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			color, white, complete := writeSVGHelper("<svg></svg>", "<svg></svg>", "<svg></svg>", tc.dirname, tc.filename)
			if tc.wantSafe {
				if color == "" || white == "" || complete == "" {
					t.Errorf("writeSVGHelper(%q, %q) returned empty paths for safe input: color=%q, white=%q, complete=%q", tc.dirname, tc.filename, color, white, complete)
				}
			} else {
				if color != "" || white != "" || complete != "" {
					t.Errorf("writeSVGHelper(%q, %q) returned non-empty paths for unsafe input: color=%q, white=%q, complete=%q", tc.dirname, tc.filename, color, white, complete)
				}
			}
		})
	}
}
