package proxy

import (
	"net/url"
	"strings"
	"testing"
)

func TestRedactURL(t *testing.T) {
	cases := []struct{ raw, want string }{
		{"/api/system/version", "/api/system/version"},
		{"/api/x?page=2", "/api/x?page=2"},
		{"/token/store?token=secret", "/token/store?token=REDACTED"},
		{"/api/x?page=2&token=secret", "/api/x?page=2&token=REDACTED"},
		// url.Query drops these malformed pairs, which would hide the token
		// from a check on the parsed values.
		{"/token/store?token=secret;extra=1", "/token/store?REDACTED"},
		{"/token/store?token=secret%", "/token/store?REDACTED"},
	}
	for _, c := range cases {
		u, err := url.Parse(c.raw)
		if err != nil {
			t.Fatalf("parse %q: %v", c.raw, err)
		}
		got := redactURL(u)
		if got != c.want {
			t.Errorf("redactURL(%q) = %q, want %q", c.raw, got, c.want)
		}
		if strings.Contains(got, "secret") {
			t.Errorf("redactURL(%q) leaked the token: %q", c.raw, got)
		}
	}
}
