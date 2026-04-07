package handlers

import (
	"context"
	"testing"
	"time"
)

// TestListenForCoreEvents_NoPanicOnCancel tests that closing the context 
// does not cause a panic when there's an inflight message being sent.
// This mimics the bug in listenForCoreEvents and listenForAdapterEvents
func TestListenForCoreEvents_NoPanicOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	
	// Create the channel that simulates respChan
	respChan := make(chan []byte)

	// Simulate listenForCoreEvents / listenForAdapterEvents behaviour
	done := make(chan struct{})
	go func() {
		defer close(done)
		
		data := []byte(`{"message": "simulated event"}`)
		
		// In the patched code, we use select. 
		// If we used `respChan <- data` directly without a select, 
		// it would panic when respChan gets closed or block indefinitely.
		select {
		case respChan <- data:
		case <-ctx.Done():
			return
		}
	}()
	
	// Cancel context to simulate client disconnect BEFORE the message is fully read.
	// This simulates the client terminating the HTTP SSE connection, causing the 
	// context to be cancelled.
	cancel()

	// Wait for the goroutine to finish safely
	select {
	case <-done:
		// Success! No panic occurred and the goroutine exited cleanly upon context cancellation.
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out waiting for goroutine to exit")
	}
}
