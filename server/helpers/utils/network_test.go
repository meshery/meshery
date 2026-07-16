package utils

import (
	"testing"
)

func TestValidateURLForOutboundRequest(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"loopback hostname", "http://localhost/api", true},
		{"loopback IP", "http://127.0.0.1/", true},
		{"private 10.x", "http://10.0.0.1/", true},
		{"private 192.168.x", "http://192.168.1.1/", true},
		{"private 172.16.x", "http://172.16.0.1/", true},
		{"AWS metadata", "http://169.254.169.254/latest/meta-data/", true},
		{"bad scheme", "ftp://example.com", true},
		{"no host", "http://", true},
		{"valid http", "http://93.184.216.34/api", false},
		{"valid https", "https://1.1.1.1/packages", false},
		{"unspecified IPv4", "http://0.0.0.0/", true},
		{"unspecified IPv6", "http://[::]/", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURLForOutboundRequest(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURLForOutboundRequest(%q) error = %v, wantErr = %v", tt.url, err, tt.wantErr)
			}
		})
	}
}
