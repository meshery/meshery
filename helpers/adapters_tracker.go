package helpers

import (
	"context"
	"sync"
)

type AdaptersTracker struct {
	adapters     map[string]struct{}
	adaptersLock *sync.Mutex
}

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

func (a *AdaptersTracker) AddAdapter(ctx context.Context, adapterURL string) {
	a.adaptersLock.Lock()
	defer a.adaptersLock.Unlock()
	a.adapters[adapterURL] = struct{}{}
}

func (a *AdaptersTracker) RemoveAdapter(ctx context.Context, adapterURL string) {
	a.adaptersLock.Lock()
	defer a.adaptersLock.Unlock()
	delete(a.adapters, adapterURL)
}

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
