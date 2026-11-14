package model

import "testing"

// TestBrokerHostPort covers the endpoint shapes that previously reached
// strings.Split(endpoint, ":")[1] and panicked with an index out of range,
// crashing the server because SubscribeToBroker runs in a goroutine with no
// recover. brokerHostPort must report ok=false for them instead of panicking.
func TestBrokerHostPort(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		wantOK   bool
		wantHost string
		wantPort int32
	}{
		{name: "empty endpoint (broker timed out with no external endpoint)", endpoint: "", wantOK: false},
		{name: "bare host without port", endpoint: "meshery-broker.meshery.svc", wantOK: false},
		{name: "bare ipv4 without port", endpoint: "10.0.0.5", wantOK: false},
		{name: "non numeric port", endpoint: "meshery-broker:nats", wantOK: false},
		{name: "host and port", endpoint: "meshery-broker.meshery.svc:4222", wantOK: true, wantHost: "meshery-broker.meshery.svc", wantPort: 4222},
		{name: "ipv4 and port", endpoint: "10.0.0.5:4222", wantOK: true, wantHost: "10.0.0.5", wantPort: 4222},
		{name: "ipv6 and port", endpoint: "[::1]:4222", wantOK: true, wantHost: "::1", wantPort: 4222},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := brokerHostPort(tt.endpoint)
			if ok != tt.wantOK {
				t.Fatalf("brokerHostPort(%q) ok = %v, want %v", tt.endpoint, ok, tt.wantOK)
			}
			if !tt.wantOK {
				return
			}
			if got.Address != tt.wantHost || got.Port != tt.wantPort {
				t.Fatalf("brokerHostPort(%q) = {%q, %d}, want {%q, %d}", tt.endpoint, got.Address, got.Port, tt.wantHost, tt.wantPort)
			}
		})
	}
}
