package utils

import (
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
