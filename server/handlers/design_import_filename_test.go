package handlers

import (
	"net/http"
	"testing"
)

// The import names a file fetched from a URL after the filename parameter of
// its Content-Disposition, or after the URL path when the header has none.
func TestGetFileNameFromResponse(t *testing.T) {
	const fileURL = "https://example.com/designs/from-url.yaml"

	for _, tc := range []struct {
		name, disposition, want string
	}{
		{"filename and filename*", `inline; filename="design.yaml"; filename*=UTF-8''design.yaml`, "design.yaml"},
		{"parameter after filename", `attachment; filename="design.yaml"; size=2048`, "design.yaml"},
		{"filename only", `attachment; filename="design.yaml"`, "design.yaml"},
		{"encoded filename*", `attachment; filename*=UTF-8''my%20design.yaml`, "my design.yaml"},
		{"no header", "", "from-url.yaml"},
		{"no filename parameter", "inline", "from-url.yaml"},
		{"empty filename", `attachment; filename=`, "from-url.yaml"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resp := &http.Response{Header: http.Header{}}
			if tc.disposition != "" {
				resp.Header.Set("Content-Disposition", tc.disposition)
			}

			if got := getFileNameFromResponse(resp, fileURL); got != tc.want {
				t.Errorf("Content-Disposition %q: file name %q, want %q", tc.disposition, got, tc.want)
			}
		})
	}
}
