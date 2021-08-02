package store

import "sync"

// threadSafeStore is a wrapper around golang native
// hashmap
// this wrapper helps adding mutex to the map
type threadSafeStore struct {
	store map[string]interface{}
	sync.RWMutex
}

// newThreadSafeStore creates a new threadSafeStore
// and a returns a pointer to it
func newThreadSafeStore() *threadSafeStore {
	return &threadSafeStore{
		store: make(map[string]interface{}),
	}
}
