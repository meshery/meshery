package models

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRetryEventSuppressorAdmitRepeatTTLAndClear(t *testing.T) {
	current := time.Unix(1_700_000_000, 0)
	suppressor := NewRetryEventSuppressor(time.Hour)
	suppressor.now = func() time.Time { return current }

	const (
		scope       = RetryEventScopeControllersHelper
		connection  = "11111111-1111-1111-1111-111111111111"
		category    = "connection"
		action      = "update"
		description = "Failed to deploy Meshery Operator"
	)

	if !suppressor.Admit(scope, connection, category, action, description) {
		t.Fatal("first emit should be admitted")
	}
	if suppressor.Admit(scope, connection, category, action, description) {
		t.Fatal("identical repeat inside the TTL should be suppressed")
	}
	if !suppressor.Admit(scope, connection, category, action, "Failed to redeploy MeshSync") {
		t.Fatal("a different description is a different event and should be admitted")
	}

	current = current.Add(time.Hour)
	if !suppressor.Admit(scope, connection, category, action, description) {
		t.Fatal("the same event should be admitted again once the TTL has elapsed")
	}
	if suppressor.Admit(scope, connection, category, action, description) {
		t.Fatal("repeat after the refreshed TTL should be suppressed")
	}

	suppressor.ClearConnection(scope, connection)
	if !suppressor.Admit(scope, connection, category, action, description) {
		t.Fatal("a success that clears the key should admit the next failure immediately")
	}
	if !suppressor.Admit(scope, "22222222-2222-2222-2222-222222222222", category, action, description) {
		t.Fatal("clearing one connection should not suppress a different connection")
	}
}

func TestRetryEventSuppressorConcurrentAdmit(t *testing.T) {
	suppressor := NewRetryEventSuppressor(time.Hour)
	var admitted atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if suppressor.Admit(RetryEventScopeConnectionStatus, "33333333-3333-3333-3333-333333333333", "connection", "update", "Operation succeeded but failed to update the status of the connection to connected.") {
				admitted.Add(1)
			}
		}()
	}
	wg.Wait()
	if got := admitted.Load(); got != 1 {
		t.Fatalf("expected exactly one concurrent admit, got %d", got)
	}
}

func TestRetryEventSuppressorIgnoresEmptyConnection(t *testing.T) {
	suppressor := NewRetryEventSuppressor(time.Hour)
	if !suppressor.Admit(RetryEventScopeConnectionStatus, "", "connection", "update", "Failed to persist") {
		t.Fatal("an event with no connection id should still be emitted")
	}
	if !suppressor.Admit(RetryEventScopeConnectionStatus, "", "connection", "update", "Failed to persist") {
		t.Fatal("events with no connection id should not suppress each other")
	}
}
