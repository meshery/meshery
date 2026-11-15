package utils

import "testing"

func TestParseURLGithubPreservesBlobDirectory(t *testing.T) {
	for _, tc := range []struct{ in, base, path string }{
		{"https://github.com/acme/project/blob/main/config/blob/design.yaml", "https://github.com/acme/project/main", "config/blob/design.yaml"},
		{"https://raw.githubusercontent.com/acme/project/main/config/blob/design.yaml", "https://raw.githubusercontent.com/acme/project/main/config/blob/design.yaml", ""},
	} {
		base, path, err := ParseURLGithub(tc.in)
		if err != nil || base != tc.base || path != tc.path {
			t.Errorf("ParseURLGithub(%q) = %q, %q, %v", tc.in, base, path, err)
		}
	}
}
