package utils

import (
	"errors"
	"fmt"
	"k8s.io/apimachinery/pkg/util/httpstream"
	"net"
	"net/http"
	"testing"
	"time"
)

type stubDialer struct{}

func (stubDialer) Dial(...string) (httpstream.Connection, string, error) {
	return stubConnection{make(chan bool)}, "portforward.k8s.io", nil
}

type stubConnection struct{ closed chan bool }

func (c stubConnection) CreateStream(http.Header) (httpstream.Stream, error) {
	return nil, fmt.Errorf("unexpected stream")
}
func (c stubConnection) Close() error                       { return nil }
func (c stubConnection) CloseChan() <-chan bool             { return c.closed }
func (c stubConnection) SetIdleTimeout(time.Duration)       {}
func (c stubConnection) RemoveStreams(...httpstream.Stream) {}

func TestPortForwardListeners(t *testing.T) {
	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	newForward := func(port int) *PortForward {
		return &PortForward{host: "127.0.0.1", namespace: "meshery", podName: "meshery-test", localPort: port, remotePort: 8080, stopCh: make(chan struct{}), readyCh: make(chan struct{}), dialer: stubDialer{}}
	}
	busy := newForward(ln.Addr().(*net.TCPAddr).Port)
	if err := busy.Init(); !IsPortBusyError(err) {
		t.Fatalf("Init() error = %v, want busy port", err)
	}
	free := newForward(0)
	if err := free.Init(); err != nil {
		t.Fatal(err)
	}
	defer free.Stop()
	if free.LocalPort() == 0 {
		t.Fatal("Init did not resolve ephemeral port")
	}
	if got, want := free.URLFor(""), fmt.Sprintf("http://127.0.0.1:%d", free.LocalPort()); got != want {
		t.Fatalf("URLFor = %q, want %q", got, want)
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
