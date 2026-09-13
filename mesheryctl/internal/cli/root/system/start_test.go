package system

import (
	"testing"
)

// TestParseContextEndpoint guards against #21696: configureDockerServices split
// the context endpoint on ":" and resolveDockerEndpoint then indexed the result
// at [1] without checking its length, panicking with "index out of range [1]
// with length 1" whenever the endpoint had no port (e.g. "localhost") or was
// empty. A malformed endpoint must now surface as an actionable error.
func TestParseContextEndpoint(t *testing.T) {
	tests := []struct {
		name        string
		endpoint    string
		wantAddress string
		wantPort    string
		wantErr     bool
	}{
		{
			name:        "scheme, host and port",
			endpoint:    "http://localhost:9081",
			wantAddress: "http://localhost",
			wantPort:    "9081",
		},
		{
			name:        "host and port, no scheme",
			endpoint:    "localhost:9081",
			wantAddress: "localhost",
			wantPort:    "9081",
		},
		{
			name:        "remote host with port",
			endpoint:    "http://example.com:9081",
			wantAddress: "http://example.com",
			wantPort:    "9081",
		},
		{
			name:        "loopback IP with port",
			endpoint:    "http://127.0.0.1:9081",
			wantAddress: "http://127.0.0.1",
			wantPort:    "9081",
		},
		{
			name:        "surrounding whitespace is tolerated",
			endpoint:    "  http://localhost:9081  ",
			wantAddress: "http://localhost",
			wantPort:    "9081",
		},
		// The cases below used to panic rather than return an error.
		{
			name:     "empty endpoint",
			endpoint: "",
			wantErr:  true,
		},
		{
			name:     "bare host, no scheme or port",
			endpoint: "localhost",
			wantErr:  true,
		},
		{
			name:     "scheme and host, no port",
			endpoint: "http://localhost",
			wantErr:  true,
		},
		{
			name:     "whitespace only",
			endpoint: "   ",
			wantErr:  true,
		},
		// net.SplitHostPort accepts any string after the colon, so the port is
		// validated separately as a usable TCP port.
		{
			name:     "non-numeric port",
			endpoint: "http://localhost:abcd",
			wantErr:  true,
		},
		{
			name:     "zero port",
			endpoint: "http://localhost:0",
			wantErr:  true,
		},
		{
			name:     "port above the valid range",
			endpoint: "http://localhost:65536",
			wantErr:  true,
		},
		{
			name:     "negative port",
			endpoint: "localhost:-1",
			wantErr:  true,
		},
		{
			name:        "highest valid port",
			endpoint:    "http://localhost:65535",
			wantAddress: "http://localhost",
			wantPort:    "65535",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address, port, err := parseContextEndpoint(tt.endpoint)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("parseContextEndpoint(%q): expected an error, got address %q and port %q", tt.endpoint, address, port)
				}
				return
			}

			if err != nil {
				t.Fatalf("parseContextEndpoint(%q): unexpected error: %v", tt.endpoint, err)
			}
			if address != tt.wantAddress {
				t.Errorf("parseContextEndpoint(%q): address = %q, want %q", tt.endpoint, address, tt.wantAddress)
			}
			if port != tt.wantPort {
				t.Errorf("parseContextEndpoint(%q): port = %q, want %q", tt.endpoint, port, tt.wantPort)
			}
		})
	}
}
