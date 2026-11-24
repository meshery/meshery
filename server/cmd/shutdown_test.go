//go:build !windows

package main

import (
	"os"
	"syscall"
	"testing"
	"time"
)

// TestSIGTERMStartsGracefulShutdown sends the process the signal Kubernetes
// stops a Pod with. Untrapped, SIGTERM kills this test binary outright, so a
// regression fails the package rather than this test alone.
func TestSIGTERMStartsGracefulShutdown(t *testing.T) {
	for _, sig := range []syscall.Signal{syscall.SIGTERM, syscall.SIGINT} {
		t.Run(sig.String(), func(t *testing.T) {
			shutdown, stop := notifyShutdown()
			defer stop()

			if err := syscall.Kill(os.Getpid(), sig); err != nil {
				t.Fatalf("failed to signal the test process: %v", err)
			}
			select {
			case got := <-shutdown:
				if got != sig {
					t.Fatalf("received %v, want %v", got, sig)
				}
			case <-time.After(5 * time.Second):
				t.Fatalf("%v did not start a graceful shutdown", sig)
			}
		})
	}
}
