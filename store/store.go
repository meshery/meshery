// Package store provides methods for interacting
// with a thread safe global store.
package store

// globalStore is intentionally a global variable so
// that the life time of the variable is through
// out the program lifecycle
var globalStore *threadSafeStore

// Initialize initializes the global store
//
// It will skip initialization if the
// store has already been initialized once
func Initialize() {
	if globalStore == nil {
		globalStore = newThreadSafeStore()
	}
}

// Set will set the key value pairs in the global store
//
// If the key already exists then it will OVERWRITE
// those values
func Set(key string, value interface{}) {
	globalStore.Lock()
	defer globalStore.Unlock()

	globalStore.store[key] = value
}

// Get will get the value corresponding to the given key
func Get(key string) (interface{}, bool) {
	globalStore.RLock()
	defer globalStore.RUnlock()

	val, ok := globalStore.store[key]
	return val, ok
}
