package models

import (
	"strings"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/models/events"
)

const (
	// RetryEventScopeControllersHelper keys MesheryControllersHelper retry events.
	RetryEventScopeControllersHelper = "controllers-helper"
	// RetryEventScopeConnectionStatus keys connection status update failures.
	RetryEventScopeConnectionStatus = "connection-status"

	retryEventKeySep = "\x1f"
)

// RetryEventSuppressor records error events so an identical repeat for the same
// connection is emitted once per TTL. ClearConnection forgets that connection's
// keys in one scope when the operation later succeeds. Admit and ClearConnection
// are safe for concurrent callers.
type RetryEventSuppressor struct {
	mu   sync.Mutex
	ttl  time.Duration
	now  func() time.Time
	seen map[string]time.Time
}

// DefaultRetryEventSuppressor is the process-wide suppressor for controller
// reconcile and connection status retries.
var DefaultRetryEventSuppressor = NewRetryEventSuppressor(time.Hour)

// NewRetryEventSuppressor returns a suppressor with the given TTL. A
// non-positive TTL uses one hour.
func NewRetryEventSuppressor(ttl time.Duration) *RetryEventSuppressor {
	if ttl <= 0 {
		ttl = time.Hour
	}
	return &RetryEventSuppressor{
		ttl:  ttl,
		now:  time.Now,
		seen: make(map[string]time.Time),
	}
}

func retryConnectionID(connectionID string) bool {
	return connectionID != "" && connectionID != uuid.Nil.String()
}

func retryEventMapKey(scope, connectionID, category, action, description string) string {
	return strings.Join([]string{scope, connectionID, category, action, description}, retryEventKeySep)
}

// Admit reports whether this error event should be emitted. The first
// observation of a connection, scope and event key is admitted. An identical
// repeat inside the TTL is not. A different description, category or action is
// a different key.
func (s *RetryEventSuppressor) Admit(scope, connectionID, category, action, description string) bool {
	if s == nil || !retryConnectionID(connectionID) {
		return true
	}
	key := retryEventMapKey(scope, connectionID, category, action, description)
	now := s.now()

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.seen == nil {
		s.seen = make(map[string]time.Time)
	}
	if at, ok := s.seen[key]; ok && now.Sub(at) < s.ttl {
		return false
	}
	s.seen[key] = now
	return true
}

// ClearConnection forgets every suppressed key for one connection in one scope
// so the next failure of that operation is reported immediately.
func (s *RetryEventSuppressor) ClearConnection(scope, connectionID string) {
	if s == nil || !retryConnectionID(connectionID) {
		return
	}
	prefix := scope + retryEventKeySep + connectionID + retryEventKeySep

	s.mu.Lock()
	defer s.mu.Unlock()
	for key := range s.seen {
		if strings.HasPrefix(key, prefix) {
			delete(s.seen, key)
		}
	}
}

// SuppressRepeatConnectionEvent reports whether severity is a failure that a
// reconcile retry should not emit again while it is unchanged.
func SuppressRepeatConnectionEvent(severity events.EventSeverity) bool {
	return severity == events.Error || severity == events.Warning
}

// ClearsConnectionEventSuppression reports whether severity is a success that
// should forget prior failures for the same operation.
func ClearsConnectionEventSuppression(severity events.EventSeverity) bool {
	return severity == events.Informational || severity == events.Success
}
