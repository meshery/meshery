package store

import "sync"

// threadSafeStore is a wrapper around golang native
// hashmap
// this wrapper helps adding mutex to the map
type threadSafeStore struct {
	store map[string]map[string]Value
	sync.RWMutex
}

// Value is an interface that all the values MUST satisfy which
// can be stored in the key value pair store
type Value interface {
	SetID(id string)
	GetID() string
}

// newThreadSafeStore creates a new threadSafeStore
// and a returns a pointer to it
func newThreadSafeStore() *threadSafeStore {
	return &threadSafeStore{
		store: make(map[string]map[string]Value),
	}
}
