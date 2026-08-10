package helpers

import (
	"testing"
)

func TestParseAdapterLocation(t *testing.T) {
	tests := []struct {
		name      string
		location  string
		wantHost  string
		wantPort  string
		expectErr bool
	}{
		{"Standard gRPC host:port", "localhost:10000", "localhost", "10000", false},
		{"HTTP scheme with port", "http://localhost:11434", "localhost", "11434", false},
		{"HTTPS scheme with port", "https://meshery.local:8080", "meshery.local", "8080", false},
		{"IPv6 with port", "[::1]:10000", "::1", "10000", false},
		{"HTTP IPv6 with port", "http://[::1]:11434", "::1", "11434", false},
		{"Host without port", "localhost", "localhost", "", false},
		{"Empty location", "", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHost, gotPort, err := ParseAdapterLocation(tt.location)
			if (err != nil) != tt.expectErr {
				t.Fatalf("ParseAdapterLocation() error = %v, expectErr %v", err, tt.expectErr)
			}
			if gotHost != tt.wantHost {
				t.Errorf("ParseAdapterLocation() gotHost = %v, want %v", gotHost, tt.wantHost)
			}
			if gotPort != tt.wantPort {
				t.Errorf("ParseAdapterLocation() gotPort = %v, want %v", gotPort, tt.wantPort)
			}
		})
	}
}
