package helpers

import (
	"context"
	"sync"
)

// AdaptersTracker is used to hold the list of known adapters
type AdaptersTracker struct {
	adapters     map[string]struct{}
	adaptersLock *sync.Mutex
}

// NewAdaptersTracker returns an instance of AdaptersTracker
func NewAdaptersTracker(adapterURLs []string) *AdaptersTracker {
	initialAdapters := map[string]struct{}{}
	for _, u := range adapterURLs {
		initialAdapters[u] = struct{}{}
	}
	a := &AdaptersTracker{
		adapters:     initialAdapters,
		adaptersLock: &sync.Mutex{},
	}

	return a
}

// AddAdapter is used to add new adapters to the collection
func (a *AdaptersTracker) AddAdapter(ctx context.Context, adapterURL string) {
	a.adaptersLock.Lock()
	defer a.adaptersLock.Unlock()
	a.adapters[adapterURL] = struct{}{}
}

// RemoveAdapter is used to remove existing adapters from the collection
func (a *AdaptersTracker) RemoveAdapter(ctx context.Context, adapterURL string) {
	a.adaptersLock.Lock()
	defer a.adaptersLock.Unlock()
	delete(a.adapters, adapterURL)
}

// GetAdapters returns the list of existing adapters
func (a *AdaptersTracker) GetAdapters(ctx context.Context) []string {
	a.adaptersLock.Lock()
	defer a.adaptersLock.Unlock()

	ad := make([]string, len(a.adapters))
	c := 0
	for x := range a.adapters {
		ad[c] = x
		c++
	}
	return ad
}
