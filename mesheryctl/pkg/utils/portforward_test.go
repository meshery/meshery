package utils

import (
	"errors"
	"testing"
	"time"
)

// waitForStop waits for the port-forward to report that it has stopped, failing
// the test if it never does. A hang here is the bug this file guards against.
func waitForStop(t *testing.T, pf *PortForward) {
	t.Helper()

	select {
	case <-pf.GetStop():
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for the port-forward to stop")
	}
}

func newTestPortForward() *PortForward {
	return &PortForward{
		host:       "localhost",
		namespace:  "meshery",
		podName:    "meshery-0",
		localPort:  9081,
		remotePort: 8080,
		stopCh:     make(chan struct{}, 1),
		readyCh:    make(chan struct{}),
	}
}

// The forwarder fails before it ever signals readiness, so Init reports the
// error to its caller.
func TestPortForwardInitFailsBeforeReady(t *testing.T) {
	pf := newTestPortForward()
	wantErr := errors.New("error upgrading connection: dial tcp: connection refused")

	err := pf.start(func() error {
		return wantErr
	})

	if !errors.Is(err, wantErr) {
		t.Fatalf("start() = %v, want %v", err, wantErr)
	}

	waitForStop(t, pf)

	if !errors.Is(pf.Err(), wantErr) {
		t.Fatalf("Err() = %v, want %v", pf.Err(), wantErr)
	}
}

// client-go closes Ready and only then waits on the connection, so a dropped
// tunnel surfaces after Init has already returned nil. The caller learns about
// it through GetStop() and Err() rather than blocking forever.
func TestPortForwardReportsTunnelDropAfterReady(t *testing.T) {
	pf := newTestPortForward()
	wantErr := errors.New("lost connection to pod")
	dropTunnel := make(chan struct{})

	err := pf.start(func() error {
		close(pf.readyCh)
		<-dropTunnel
		return wantErr
	})
	if err != nil {
		t.Fatalf("start() = %v, want nil once the tunnel is ready", err)
	}

	select {
	case <-pf.GetStop():
		t.Fatal("port-forward reported a stop while the tunnel was still up")
	default:
	}

	close(dropTunnel)
	waitForStop(t, pf)

	if !errors.Is(pf.Err(), wantErr) {
		t.Fatalf("Err() = %v, want %v", pf.Err(), wantErr)
	}
}

// A tunnel torn down by Stop() is not a failure, so Err() stays nil.
func TestPortForwardStopIsNotAnError(t *testing.T) {
	pf := newTestPortForward()

	err := pf.start(func() error {
		close(pf.readyCh)
		<-pf.stopCh
		return nil
	})
	if err != nil {
		t.Fatalf("start() = %v, want nil", err)
	}

	pf.Stop()
	waitForStop(t, pf)

	if pf.Err() != nil {
		t.Fatalf("Err() = %v, want nil after a clean stop", pf.Err())
	}
}

// Stop() races with the forwarder goroutine closing the same channel, and
// callers are told to call it even after an error, so it has to be repeatable.
func TestPortForwardStopIsIdempotent(t *testing.T) {
	pf := newTestPortForward()

	err := pf.start(func() error {
		return errors.New("unable to listen on any of the requested ports")
	})
	if err == nil {
		t.Fatal("start() = nil, want an error")
	}

	waitForStop(t, pf)

	// The forwarder goroutine has already closed stopCh; these must not panic.
	pf.Stop()
	pf.Stop()
}
