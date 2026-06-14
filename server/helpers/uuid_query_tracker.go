package helpers

import (
	"context"
	"sync"
)

// UUIDQueryTracker tracks queries for a load test UUID
type UUIDQueryTracker struct {
	queries map[string]map[string]bool
	qLock   *sync.Mutex
}

// NewUUIDQueryTracker creates a new instance of UUIDQueryTracker
func NewUUIDQueryTracker() *UUIDQueryTracker {
	return &UUIDQueryTracker{
		queries: map[string]map[string]bool{},
		qLock:   &sync.Mutex{},
	}
}

// AddOrFlagQuery either adds a new query or flags an existing one
func (a *UUIDQueryTracker) AddOrFlagQuery(_ context.Context, uuid, query string, flag bool) {
	a.qLock.Lock()
	defer a.qLock.Unlock()
	queries, ok := a.queries[uuid]
	if !ok {
		queries = map[string]bool{}
	}
	queries[query] = flag
	a.queries[uuid] = queries
}

// RemoveUUID removes an existing UUID from the collection
func (a *UUIDQueryTracker) RemoveUUID(_ context.Context, uuid string) {
	a.qLock.Lock()
	defer a.qLock.Unlock()
	delete(a.queries, uuid)
}

// GetQueriesForUUID retrieves queries for UUID
func (a *UUIDQueryTracker) GetQueriesForUUID(_ context.Context, uuid string) map[string]bool {
	a.qLock.Lock()
	defer a.qLock.Unlock()
	queries, ok := a.queries[uuid]
	if !ok {
		return map[string]bool{}
	}
	return queries
}
