package proxy

import "testing"

// A logout must not be lost when a login notification is still queued in
// the auth channel: the queued value is stale state, and the listener must
// receive the latest state instead.
func TestSetToken_ReplacesQueuedAuthState(t *testing.T) {
	p := &Proxy{authWsChan: make(chan bool, 1)}
	// Simulate a login notification the websocket listener has not read yet.
	p.authWsChan <- true

	p.setToken("")

	select {
	case got := <-p.authWsChan:
		if got {
			t.Fatalf("expected queued auth state false after logout, got true")
		}
	default:
		t.Fatalf("expected a queued auth state after setToken, channel empty")
	}
}

// The same rule in the other direction: a login replaces a queued logout.
func TestSetToken_ReplacesQueuedLogout(t *testing.T) {
	p := &Proxy{authWsChan: make(chan bool, 1)}
	p.authWsChan <- false

	p.setToken("abc")

	select {
	case got := <-p.authWsChan:
		if !got {
			t.Fatalf("expected queued auth state true after login, got false")
		}
	default:
		t.Fatalf("expected a queued auth state after setToken, channel empty")
	}
}
