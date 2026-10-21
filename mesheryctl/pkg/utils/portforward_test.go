package utils

import (
	"errors"
	"testing"
)

func TestPortForwardURLUsesBoundLocalPort(t *testing.T) {
	pf := &PortForward{host: "localhost", localPort: 49152}
	if got, want := pf.URLFor(""), "http://localhost:49152"; got != want {
		t.Fatalf("URLFor(\"\") = %q, want %q", got, want)
	}
	if got := pf.LocalPort(); got != 49152 {
		t.Fatalf("LocalPort() = %d, want 49152", got)
	}
}

func TestIsPortBusyError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil"},
		{name: "unix address in use", err: errors.New("bind: address already in use"), want: true},
		{name: "windows address in use", err: errors.New("Only one usage of each socket address is normally permitted"), want: true},
		{name: "client-go listener detail", err: errors.New("unable to listen on any of the requested ports: [{9081 8080}]: Unable to listen on port 9081: Listeners failed to create with the following errors: [unable to create listener: Error listen tcp4 127.0.0.1:9081: bind: address already in use] for meshery/meshery-abc"), want: true},
		{name: "generic listener failure", err: errors.New("unable to listen on any of the requested ports: [{9081 8080}]"), want: false},
		{name: "permission denied", err: errors.New("unable to create listener: Error listen tcp4 127.0.0.1:9081: bind: permission denied")},
		{name: "connection refused", err: errors.New("connection refused")},
		{name: "no running pods", err: errors.New("no running pods found for meshery")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPortBusyError(tt.err); got != tt.want {
				t.Fatalf("IsPortBusyError(%v) = %t, want %t", tt.err, got, tt.want)
			}
		})
	}
}
