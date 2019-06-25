package helpers

import (
	"context"
	"sync"
)

type UUIDQueryTracker struct {
	queries map[string]map[string]bool
	qLock   *sync.Mutex
}

func NewUUIDQueryTracker() *UUIDQueryTracker {
	return &UUIDQueryTracker{
		queries: map[string]map[string]bool{},
		qLock:   &sync.Mutex{},
	}
}

func (a *UUIDQueryTracker) AddOrFlagQuery(ctx context.Context, uuid, query string, flag bool) {
	a.qLock.Lock()
	defer a.qLock.Unlock()
	queries, ok := a.queries[uuid]
	if !ok {
		queries = map[string]bool{}
	}
	queries[query] = flag
	a.queries[uuid] = queries
}

func (a *UUIDQueryTracker) RemoveUUID(ctx context.Context, uuid string) {
	a.qLock.Lock()
	defer a.qLock.Unlock()
	delete(a.queries, uuid)
}

func (a *UUIDQueryTracker) GetQueriesForUUID(ctx context.Context, uuid string) map[string]bool {
	a.qLock.Lock()
	defer a.qLock.Unlock()
	queries, ok := a.queries[uuid]
	if !ok {
		return map[string]bool{}
	}
	return queries
}
